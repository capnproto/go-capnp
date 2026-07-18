package bbr

import (
	"context"
	"sort"
	"sync"
	"time"

	"capnproto.org/go/capnp/v3/exp/clock"
)

// simClock is a clock whose timers are fired explicitly by simulator.  It
// deliberately does not try to run goroutines: the simulator chooses both
// virtual-time ordering and the next timer to fire.
type simClock struct {
	mu      sync.Mutex
	now     time.Time
	nextID  uint64
	timers  []*simTimer
	changed chan struct{}
}

type simTimer struct {
	clock    *simClock
	deadline time.Time
	id       uint64
	active   bool
	ch       chan time.Time
}

func newSimClock(now time.Time) *simClock {
	return &simClock{now: now, changed: make(chan struct{}, 1)}
}

func (c *simClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *simClock) NewTimer(d time.Duration) clock.Timer {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.nextID++
	t := &simTimer{
		clock:    c,
		deadline: c.now.Add(d),
		id:       c.nextID,
		active:   true,
		ch:       make(chan time.Time, 1),
	}
	c.timers = append(c.timers, t)
	c.notifyChanged()
	return t
}

func (c *simClock) notifyChanged() {
	select {
	case c.changed <- struct{}{}:
	default:
	}
}

func (c *simClock) AdvanceTo(now time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if now.Before(c.now) {
		panic("simulated clock moved backwards")
	}
	c.now = now
}

func (c *simClock) nextTimer() *simTimer {
	c.mu.Lock()
	defer c.mu.Unlock()
	var next *simTimer
	for _, t := range c.timers {
		if !t.active || (next != nil && (next.deadline.Before(t.deadline) || next.deadline.Equal(t.deadline) && next.id < t.id)) {
			continue
		}
		next = t
	}
	return next
}

func (c *simClock) fire(t *simTimer) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if !t.active || t.deadline.After(c.now) {
		return false
	}
	t.active = false
	t.ch <- c.now
	return true
}

func (t *simTimer) Chan() <-chan time.Time { return t.ch }

func (t *simTimer) Reset(d time.Duration) {
	t.clock.mu.Lock()
	defer t.clock.mu.Unlock()
	t.deadline = t.clock.now.Add(d)
	t.active = true
	t.clock.notifyChanged()
}

func (t *simTimer) Stop() bool {
	t.clock.mu.Lock()
	defer t.clock.mu.Unlock()
	wasActive := t.active
	t.active = false
	return wasActive
}

type simPacket struct {
	deadline time.Time
	sequence uint64
	packet   testPacket
}

// simPath reproduces testLink's fixed-delay and serial bandwidth-limited
// links without worker goroutines.
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

type simulator struct {
	t       testingT
	clock   *simClock
	path    *simPath
	lim     *Limiter
	events  chan struct{}
	results []Snapshot
}

// testingT keeps the simulator useful to the small path tests without making
// it responsible for reporting failures from a goroutine.
type testingT interface {
	Helper()
	Fatal(...any)
}

func newSimulator(t testingT, links ...testLink) *simulator {
	t.Helper()
	events := make(chan struct{}, 1)
	clock := newSimClock(sampleStartTime)
	s := &simulator{
		t:      t,
		clock:  clock,
		path:   newSimPath(clock, links...),
		events: events,
	}
	s.lim = newLimiter(clock, limiterOptions{
		randInt: func() int { return 6 }, // begin ProbeBW at the 1.25 gain phase
		afterEvent: func() {
			events <- struct{}{}
		},
	})
	return s
}

func (s *simulator) close() { s.lim.Release() }

func (s *simulator) snapshot() {
	s.results = append(s.results, SnapshotLimiter(s.lim))
}

func (s *simulator) settle() { _ = SnapshotLimiter(s.lim) }

func (s *simulator) waitEvent() { <-s.events }

func (s *simulator) drainEvents() {
	for {
		select {
		case <-s.events:
		default:
			return
		}
	}
}

func (s *simulator) step() bool {
	packet, hasPacket := s.path.nextPacket()
	timer := s.clock.nextTimer()
	if !hasPacket && timer == nil {
		return false
	}
	if hasPacket && (timer == nil || packet.deadline.Before(timer.deadline) || packet.deadline.Equal(timer.deadline)) {
		s.clock.AdvanceTo(packet.deadline)
		s.path.popPacket().packet.gotResponse()
		s.waitEvent()
		s.snapshot()
		s.drainEvents()
		return true
	}
	s.clock.AdvanceTo(timer.deadline)
	if !s.clock.fire(timer) {
		return true
	}
	s.waitEvent()
	s.settle()
	s.snapshot()
	s.drainEvents()
	return true
}

func (s *simulator) send(size uint64) {
	result := make(chan testPacket, 1)
	ctx := context.Background()
	go func() {
		gotResponse, err := s.lim.StartMessage(ctx, size)
		if err != nil {
			panic(err)
		}
		result <- testPacket{size: size, gotResponse: gotResponse}
	}()

	for {
		select {
		case packet := <-result:
			s.path.send(packet)
			s.settle()
			s.snapshot()
			s.drainEvents()
			return
		default:
		}
		if s.step() {
			continue
		}
		select {
		case packet := <-result:
			s.path.send(packet)
			s.settle()
			s.snapshot()
			s.drainEvents()
			return
		case <-s.events:
		case <-s.clock.changed:
		}
	}
}

func (s *simulator) drain() {
	for s.step() {
	}
}

func (s *simulator) run(packetSizes []uint64) []Snapshot {
	defer s.close()
	for _, size := range packetSizes {
		s.send(size)
	}
	s.drain()
	return s.results
}
