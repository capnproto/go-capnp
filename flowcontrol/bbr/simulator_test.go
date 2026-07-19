package bbr

import (
	"context"
	"sort"
	"sync"
	"testing"
	"time"

	"capnproto.org/go/capnp/v3/exp/clock"
)

// simClock advances only when the synchronous simulator tells it to. Its
// NewTimer method deliberately panics: limiter tests must drive prepareStep
// and step directly rather than accidentally starting the production actor.
type simClock struct {
	now time.Time
}

func newSimClock(now time.Time) *simClock { return &simClock{now: now} }

func (c *simClock) Now() time.Time { return c.now }

func (c *simClock) NewTimer(time.Duration) clock.Timer {
	panic("bbr: synchronous simulator must not create clock timers")
}

func (c *simClock) AdvanceTo(now time.Time) {
	if now.Before(c.now) {
		panic("simulated clock moved backwards")
	}
	c.now = now
}

type simPacket struct {
	deadline time.Time
	sequence uint64
	packet   testPacket
}

// simPath models overlapping fixed propagation delay and serial
// bandwidth-limited links without worker goroutines.
type simPath struct {
	clock         *simClock
	links         []testLink
	nextAvailable []time.Time
	sequence      uint64
	packets       []simPacket
}

func newSimPath(clock *simClock, links ...testLink) *simPath {
	return &simPath{
		clock:         clock,
		links:         links,
		nextAvailable: make([]time.Time, len(links)),
	}
}

func (p *simPath) send(packet testPacket) {
	deadline := p.clock.Now()
	for i, link := range p.links {
		if link.bandwidth == 0 {
			deadline = deadline.Add(link.delay)
			continue
		}
		if p.nextAvailable[i].After(deadline) {
			deadline = p.nextAvailable[i]
		}
		deadline = deadline.Add(time.Duration(float64(packet.size) / float64(link.bandwidth)))
		p.nextAvailable[i] = deadline
	}
	p.sequence++
	p.packets = append(p.packets, simPacket{deadline: deadline, sequence: p.sequence, packet: packet})
	sort.SliceStable(p.packets, func(i, j int) bool {
		if p.packets[i].deadline.Equal(p.packets[j].deadline) {
			return p.packets[i].sequence < p.packets[j].sequence
		}
		return p.packets[i].deadline.Before(p.packets[j].deadline)
	})
}

func (p *simPath) nextPacket() (simPacket, bool) {
	if len(p.packets) == 0 {
		return simPacket{}, false
	}
	return p.packets[0], true
}

func (p *simPath) popPacket() simPacket {
	packet := p.packets[0]
	p.packets = p.packets[1:]
	return packet
}

func (p *simPath) advance(d time.Duration) {
	p.clock.AdvanceTo(p.clock.Now().Add(d))
	for len(p.packets) > 0 && !p.packets[0].deadline.After(p.clock.Now()) {
		p.popPacket().packet.gotResponse()
	}
}

func TestLinkDelayOverlapsPackets(t *testing.T) {
	clock := newSimClock(sampleStartTime)
	path := newSimPath(clock, testLink{delay: 4 * time.Second})
	var arrivals []time.Time
	for i := 0; i < 2; i++ {
		path.send(testPacket{gotResponse: func() {
			arrivals = append(arrivals, clock.Now())
		}})
	}

	path.advance(4 * time.Second)
	if len(arrivals) != 2 {
		t.Fatalf("fixed-delay link delivered %d packets; want 2", len(arrivals))
	}
	if !arrivals[0].Equal(arrivals[1]) {
		t.Fatalf("fixed-delay packets arrived at %v and %v; propagation delay must overlap", arrivals[0], arrivals[1])
	}
}

type simulator struct {
	clock   *simClock
	path    *simPath
	lim     *Limiter
	wait    limiterWait
	results []Snapshot
}

func newSimulator(links ...testLink) *simulator {
	clock := newSimClock(sampleStartTime)
	s := &simulator{
		clock: clock,
		path:  newSimPath(clock, links...),
		lim: newLimiterState(clock, limiterOptions{
			randInt: func() int { return 6 }, // enter ProbeBW at the 1.25 phase
		}),
	}
	s.wait = s.lim.prepareStep(clock.Now())
	return s
}

// apply owns one atomic limiter transition and records its stable wait-state
// boundary. A sent packet is put on the path before the snapshot, but ACK
// delivery remains an explicit later event.
func (s *simulator) apply(event limiterEvent) {
	result := s.lim.step(s.clock.Now(), event)
	if result.sent {
		packet := result.packet
		s.path.send(testPacket{
			size: packet.Size,
			gotResponse: func() {
				s.apply(limiterEvent{kind: limiterAckEvent, ack: packet})
			},
		})
	}
	s.wait = s.lim.prepareStep(s.clock.Now())
	s.results = append(s.results, s.lim.snapshotAt(s.clock.Now()))
}

func (s *simulator) deliverNextAck() bool {
	packet, ok := s.path.nextPacket()
	if !ok {
		return false
	}
	s.clock.AdvanceTo(packet.deadline)
	s.path.popPacket().packet.gotResponse()
	return true
}

func (s *simulator) deliverDueAck() bool {
	packet, ok := s.path.nextPacket()
	if !ok || packet.deadline.After(s.clock.Now()) {
		return false
	}
	return s.deliverNextAck()
}

func (s *simulator) send(size uint64) {
	for {
		// At one virtual timestamp, ACKs are processed FIFO before a due
		// pacing event and before this application send.
		if s.deliverDueAck() {
			continue
		}
		if s.wait.acceptSend {
			s.apply(limiterEvent{kind: limiterSendEvent, size: size})
			return
		}

		packet, hasPacket := s.path.nextPacket()
		if s.wait.hasDeadline {
			// Equal-time ACKs win. The ACK transition may invalidate or
			// replace this pacing deadline before it can fire.
			if hasPacket && !s.wait.deadline.Before(packet.deadline) {
				s.deliverNextAck()
				continue
			}
			s.clock.AdvanceTo(s.wait.deadline)
			s.apply(limiterEvent{
				kind:        limiterPacingEvent,
				size:        size,
				sendPending: true,
			})
			return
		}

		if !hasPacket {
			panic("bbr: limiter is ACK-gated with no packet in flight")
		}
		s.deliverNextAck()
	}
}

func (s *simulator) drain() {
	for {
		packet, hasPacket := s.path.nextPacket()
		if !hasPacket {
			return
		}
		if !packet.deadline.After(s.clock.Now()) {
			s.deliverNextAck()
			continue
		}
		if s.wait.hasDeadline && s.wait.deadline.Before(packet.deadline) {
			s.clock.AdvanceTo(s.wait.deadline)
			s.apply(limiterEvent{kind: limiterPacingEvent})
			continue
		}
		// An ACK wins a tie with a pacing deadline.
		s.deliverNextAck()
	}
}

func (s *simulator) run(packetSizes []uint64) []Snapshot {
	defer s.lim.Release()
	for _, size := range packetSizes {
		s.send(size)
	}
	s.drain()
	return s.results
}

func TestLimiterPrepareStepGates(t *testing.T) {
	t.Run("zero inflight bypasses limits", func(t *testing.T) {
		clock := newSimClock(sampleStartTime)
		lim := newLimiterState(clock, limiterOptions{})
		lim.maxPacketsInflight = 0
		if wait := lim.prepareStep(clock.Now()); !wait.acceptSend {
			t.Fatal("zero-inflight limiter did not accept a send")
		}
	})

	t.Run("packet limit gates sends", func(t *testing.T) {
		clock := newSimClock(sampleStartTime)
		lim := newLimiterState(clock, limiterOptions{})
		lim.sent = 1
		lim.packetsInflight = 1
		lim.maxPacketsInflight = 1
		wait := lim.prepareStep(clock.Now())
		if wait.acceptSend || wait.hasDeadline {
			t.Fatal("packet-limited state should wait only for an ACK")
		}
	})

	t.Run("byte window uses greater-than-or-equal gate", func(t *testing.T) {
		clock := newSimClock(sampleStartTime)
		lim := newLimiterState(clock, limiterOptions{})
		lim.sent = 1
		lim.packetsInflight = 1
		lim.btlBwFilter.Estimate = 1 * bytesPerSecond
		lim.rtPropFilter.Estimate = time.Second
		lim.cwndGain = 1
		wait := lim.prepareStep(clock.Now())
		if wait.acceptSend || wait.hasDeadline {
			t.Fatal("limiter at its byte window should wait only for an ACK")
		}
	})

	t.Run("zero BDP does not gate pacing", func(t *testing.T) {
		clock := newSimClock(sampleStartTime)
		lim := newLimiterState(clock, limiterOptions{})
		lim.sent = 1
		lim.packetsInflight = 1
		lim.rtPropFilter.Estimate = 0
		if wait := lim.prepareStep(clock.Now()); !wait.hasDeadline {
			t.Fatal("zero BDP incorrectly activated the byte window gate")
		}
	})
}

func TestLimiterPacingDeadline(t *testing.T) {
	newPacedLimiter := func(t *testing.T) (*simClock, *Limiter, limiterWait) {
		t.Helper()
		clock := newSimClock(sampleStartTime)
		lim := newLimiterState(clock, limiterOptions{randInt: func() int { return 6 }})
		if wait := lim.prepareStep(clock.Now()); !wait.acceptSend {
			t.Fatal("new limiter did not accept its first send")
		}
		result := lim.step(clock.Now(), limiterEvent{kind: limiterSendEvent, size: 1})
		if !result.sent {
			t.Fatal("first send was not applied")
		}
		wait := lim.prepareStep(clock.Now())
		if !wait.hasDeadline {
			t.Fatal("in-flight limiter did not expose a pacing deadline")
		}
		clock.AdvanceTo(wait.deadline)
		return clock, lim, wait
	}

	t.Run("pending send at equality is not app limited", func(t *testing.T) {
		clock, lim, _ := newPacedLimiter(t)
		wait := lim.prepareStep(clock.Now())
		if wait.acceptSend || !wait.hasDeadline || !wait.deadline.Equal(clock.Now()) {
			t.Fatal("exact deadline was not preserved as a pacing event")
		}
		result := lim.step(clock.Now(), limiterEvent{
			kind:        limiterPacingEvent,
			size:        1,
			sendPending: true,
		})
		if !result.sent || result.packet.AppLimited {
			t.Fatal("pending deadline send was not sent as pacing-limited traffic")
		}
		if lim.pacingReady {
			t.Fatal("send did not clear pacing-ready latch")
		}
	})

	t.Run("idle equality latches readiness once", func(t *testing.T) {
		clock, lim, _ := newPacedLimiter(t)
		if wait := lim.prepareStep(clock.Now()); wait.acceptSend || !wait.hasDeadline {
			t.Fatal("exact deadline did not require a pacing event")
		}
		lim.step(clock.Now(), limiterEvent{kind: limiterPacingEvent})
		if !lim.pacingReady || lim.appLimitedUntil != lim.inflight() {
			t.Fatal("idle pacing event did not atomically mark app-limited readiness")
		}
		maxPacketsInflight := lim.maxPacketsInflight
		lim.maxPacketsInflight = lim.packetsInflight
		if wait := lim.prepareStep(clock.Now()); wait.acceptSend || wait.hasDeadline {
			t.Fatal("pacing-ready latch bypassed the packet window gate")
		}
		lim.maxPacketsInflight = maxPacketsInflight
		for i := 0; i < 2; i++ {
			if wait := lim.prepareStep(clock.Now()); !wait.acceptSend || wait.hasDeadline {
				t.Fatal("latched pacing event re-created a zero-duration timer")
			}
		}
		result := lim.step(clock.Now(), limiterEvent{kind: limiterSendEvent, size: 1})
		if !result.sent || !result.packet.AppLimited || lim.pacingReady {
			t.Fatal("send after idle pacing did not consume app-limited latch")
		}
	})
}

func TestLimiterAppLimitedAccountingWithOutOfOrderACK(t *testing.T) {
	clock := newSimClock(sampleStartTime)
	lim := newLimiterState(clock, limiterOptions{})

	older := lim.step(clock.Now(), limiterEvent{kind: limiterSendEvent, size: 2}).packet
	lim.appLimitedUntil = lim.inflight()
	newer := lim.step(clock.Now(), limiterEvent{kind: limiterSendEvent, size: 100}).packet
	if !newer.AppLimited {
		t.Fatal("packet sent past the app-limited boundary was not marked app-limited")
	}

	clock.AdvanceTo(clock.Now().Add(time.Second))
	lim.step(clock.Now(), limiterEvent{kind: limiterAckEvent, ack: newer})
	if lim.appLimitedUntil != 0 {
		t.Fatalf("app-limited marker after larger out-of-order ACK = %d; want 0", lim.appLimitedUntil)
	}
	if lim.inflight() != older.Size || lim.packetsInflight != 1 {
		t.Fatalf("out-of-order ACK left %d bytes in %d packets; want %d byte in 1 packet", lim.inflight(), lim.packetsInflight, older.Size)
	}
}

type trackingClock struct {
	mu      sync.Mutex
	now     time.Time
	created chan *trackingTimer
}

func newTrackingClock(now time.Time) *trackingClock {
	return &trackingClock{now: now, created: make(chan *trackingTimer, 8)}
}

func (c *trackingClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *trackingClock) Advance(d time.Duration) {
	c.mu.Lock()
	c.now = c.now.Add(d)
	c.mu.Unlock()
}

func (c *trackingClock) NewTimer(time.Duration) clock.Timer {
	timer := &trackingTimer{ch: make(chan time.Time), stopped: make(chan struct{})}
	c.created <- timer
	return timer
}

type trackingTimer struct {
	ch      chan time.Time
	stopped chan struct{}
	once    sync.Once
}

func (t *trackingTimer) Chan() <-chan time.Time { return t.ch }
func (t *trackingTimer) Reset(time.Duration)    {}
func (t *trackingTimer) Stop() (stopped bool) {
	t.once.Do(func() {
		stopped = true
		close(t.stopped)
	})
	return stopped
}

func TestLimiterRunStopsAbandonedPacingTimer(t *testing.T) {
	for _, test := range []struct {
		name string
		act  func(*trackingClock, *Limiter, func())
	}{
		{
			name: "ACK",
			act: func(clock *trackingClock, _ *Limiter, gotResponse func()) {
				clock.Advance(time.Millisecond)
				gotResponse()
			},
		},
		{
			name: "pause",
			act: func(_ *trackingClock, lim *Limiter, _ func()) {
				_ = SnapshotLimiter(lim)
			},
		},
		{
			name: "cancellation",
			act: func(_ *trackingClock, lim *Limiter, _ func()) {
				lim.Release()
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			clock := newTrackingClock(sampleStartTime)
			lim := NewLimiter(clock)
			gotResponse, err := lim.StartMessage(context.Background(), 1)
			if err != nil {
				t.Fatal(err)
			}

			var timer *trackingTimer
			select {
			case timer = <-clock.created:
			case <-time.After(time.Second):
				t.Fatal("limiter did not create its pacing timer")
			}
			test.act(clock, lim, gotResponse)
			select {
			case <-timer.stopped:
			case <-time.After(time.Second):
				t.Fatal("abandoned pacing timer was not stopped")
			}

			lim.Release()
			select {
			case <-lim.done:
			case <-time.After(time.Second):
				t.Fatal("limiter actor did not exit after Release")
			}
			// Once done is closed there is no actor to pause; this must still
			// return a safe final snapshot rather than deadlocking.
			_ = SnapshotLimiter(lim)
		})
	}
}

func TestLimiterRunCompletesAcceptedSendAfterRelease(t *testing.T) {
	lim := NewLimiter(newTrackingClock(sampleStartTime))
	reply := make(chan packetMeta)
	request := sendRequest{size: 7, replyChan: reply}

	select {
	case lim.chSend <- request:
		// The unbuffered rendezvous proves the actor accepted the request.
	case <-time.After(time.Second):
		lim.Release()
		t.Fatal("limiter did not accept the send request")
	}
	lim.Release()

	select {
	case packet := <-reply:
		if packet.Size != request.size {
			t.Fatalf("accepted packet size = %d; want %d", packet.Size, request.size)
		}
	case <-time.After(time.Second):
		t.Fatal("Release interrupted an already-accepted send")
	}
	select {
	case <-lim.done:
	case <-time.After(time.Second):
		t.Fatal("limiter actor did not exit after completing the accepted send")
	}
}
