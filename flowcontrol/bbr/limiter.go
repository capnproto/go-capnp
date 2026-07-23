package bbr

import (
	"context"
	"math"
	"math/rand"
	"time"

	"capnproto.org/go/capnp/v3/exp/clock"
)

// A packetMeta contains metadata about a packet that was sent.
type packetMeta struct {
	SendTime time.Time // The time at which the packet was sent.
	Size     uint64    // The size of the packet.

	// Whether the connection flow was app-limited when this packet
	// was sent:
	AppLimited bool

	// value of corresponding fields in Limiter when this packet was sent:
	Delivered     uint64
	DeliveredTime time.Time
}

type sendRequest struct {
	size      uint64
	replyChan chan<- packetMeta
}

type limiterWait struct {
	acceptSend  bool
	hasDeadline bool
	deadline    time.Time
}

type limiterEventKind uint8

const (
	limiterSendEvent limiterEventKind = iota
	limiterAckEvent
	limiterPacingEvent
)

type limiterEvent struct {
	kind limiterEventKind
	size uint64
	ack  packetMeta

	// sendPending is only meaningful for limiterPacingEvent. A pacing
	// deadline atomically consumes a waiting send when one is available,
	// matching the non-blocking receive in the production actor.
	sendPending bool
}

type limiterStepResult struct {
	sent   bool
	packet packetMeta
}

func (l *Limiter) StartMessage(ctx context.Context, size uint64) (gotResponse func(), err error) {
	replyChan := make(chan packetMeta)
	select {
	case <-ctx.Done():
		return func() {}, ctx.Err()
	case <-l.ctx.Done():
		return func() {}, l.ctx.Err()
	case l.chSend <- sendRequest{size: size, replyChan: replyChan}:
		pm := <-replyChan
		return func() {
			select {
			case <-l.ctx.Done():
			case l.chAck <- pm:
			}
		}, nil
	}
}

func (l *Limiter) Release() {
	l.cancel()
}

// Compute the bandwidth delay product, in bytes.
func (l *Limiter) computeBDP() float64 {
	bandwidth := l.btlBwFilter.Estimate
	delay := l.rtPropFilter.Estimate
	return float64(bandwidth) * float64(delay.Nanoseconds())
}

// prepareStep determines what the limiter actor should wait for next. It and
// step form the synchronous limiter state machine. They must be called either
// by run, or by the sole owner of an unstarted limiter in a test; they are not
// safe to call concurrently with a running limiter actor.
func (l *Limiter) prepareStep(now time.Time) limiterWait {
	bdp := l.computeBDP()

	if l.inflight() == 0 {
		// With nothing on the wire there is no ACK that can unblock an
		// inaccurate initial pacing estimate, so always accept a send.
		return limiterWait{acceptSend: true}
	}
	if l.packetsInflight >= l.maxPacketsInflight ||
		(bdp > 0 && float64(l.inflight()) >= l.cwndGain*bdp) {
		return limiterWait{}
	}
	if l.pacingReady {
		return limiterWait{acceptSend: true}
	}
	if now.After(l.nextSendTime) {
		l.appLimitedUntil = l.inflight()
		return limiterWait{acceptSend: true}
	}
	return limiterWait{hasDeadline: true, deadline: l.nextSendTime}
}

// step applies one explicit ACK, send, or pacing-deadline event. Pause and
// cancellation do not alter congestion-control state and remain concerns of
// the production channel adapter in run.
func (l *Limiter) step(now time.Time, event limiterEvent) limiterStepResult {
	switch event.kind {
	case limiterSendEvent:
		return limiterStepResult{sent: true, packet: l.doSendAt(now, event.size)}
	case limiterAckEvent:
		l.onAckAt(now, event.ack)
		return limiterStepResult{}
	case limiterPacingEvent:
		if event.sendPending {
			return limiterStepResult{sent: true, packet: l.doSendAt(now, event.size)}
		}
		// Remember that this pacing opportunity has already fired. A real
		// clock advances past a zero-duration timer; a virtual clock does
		// not, so the latch prevents repeatedly firing at the same instant.
		l.appLimitedUntil = l.inflight()
		l.pacingReady = true
		return limiterStepResult{}
	default:
		panic("bbr: invalid limiter event")
	}
}

func (l *Limiter) run(ctx context.Context) {
	defer close(l.done)
	for {
		var (
			sendTimer  clock.Timer
			timeToSend <-chan time.Time
			sendReqs   <-chan sendRequest
		)

		now := l.clock.Now()
		wait := l.prepareStep(now)
		if wait.acceptSend {
			sendReqs = l.chSend
		} else if wait.hasDeadline {
			sleep := wait.deadline.Sub(now)
			sendTimer = l.clock.NewTimer(sleep)
			timeToSend = sendTimer.Chan()
		}

		var (
			event       limiterEvent
			eventTime   time.Time
			handleEvent bool
			replyChan   chan<- packetMeta
			cancelled   bool
			paused      bool
		)
		select {
		case <-ctx.Done():
			cancelled = true
		case p := <-l.chAck:
			eventTime = l.clock.Now()
			event = limiterEvent{kind: limiterAckEvent, ack: p}
			handleEvent = true
		case <-timeToSend:
			event = limiterEvent{kind: limiterPacingEvent}
			select {
			case req := <-l.chSend:
				eventTime = l.clock.Now()
				event.size = req.size
				event.sendPending = true
				replyChan = req.replyChan
			default:
			}
			handleEvent = true
		case req := <-sendReqs:
			eventTime = l.clock.Now()
			event = limiterEvent{kind: limiterSendEvent, size: req.size}
			replyChan = req.replyChan
			handleEvent = true
		case gateEvent := <-l.chGate:
			l.handleGateEvent(gateEvent)
		case <-l.chPause:
			paused = true
		}
		if sendTimer != nil {
			sendTimer.Stop()
		}
		if cancelled {
			return
		}
		if paused {
			l.chPause <- struct{}{}
			continue
		}
		if handleEvent {
			result := l.step(eventTime, event)
			if result.sent {
				replyChan <- result.packet
			}
		}
	}
}

func (l *Limiter) doSendAt(now time.Time, size uint64) packetMeta {
	l.pacingReady = false
	l.packetsInflight++
	p := packetMeta{
		Size:          size,
		AppLimited:    l.appLimitedUntil > 0,
		SendTime:      now,
		Delivered:     l.delivered,
		DeliveredTime: l.deliveredTime,
	}
	l.nextSendTime = now.Add(time.Duration(
		float64(size) / (l.pacingGain * float64(l.btlBwFilter.Estimate)),
	))
	l.sent += size
	return p
}

// A Limiter implements flowcontrol.FlowLimiter using the BBR algorithm.
type Limiter struct {
	// When ctx is canceled, the Limiter's goroutine will exit.
	// Senders must monitor this context when sending, to detect
	// if the manager has shut down.
	ctx context.Context

	// cancels ctx
	cancel context.CancelFunc

	// When a response to a packet comes in, the original packetMeta
	// should be sent on this channel.
	chAck chan packetMeta

	// Used to request permission to send data. This channel is
	// unbuffered, and the manager's goroutine will only receive
	// when it is appropriate to send a packet. Once the manager
	// goroutine receives from this channel, it will promptly and
	// unconditionally send the corresponding packetMeta on the
	// request's replyChan. The sending goroutine must immediately
	// read back this response.
	chSend chan sendRequest

	// Filters for estimating the round trip propagation time and
	// bottleneck bandwidth, respectively.
	//
	// The bottlneck bandwidth is measured in units of
	// bytes/time.Duration(1) i.e. bytes per nanosecond.
	rtPropFilter rtPropFilter
	btlBwFilter  btlBwFilter

	// cwndGain (Congestion WiNDow; standard terminology) is a factor
	// by which we may exceed BDP. i.e. if inflight() is greater than
	// estimated BDP * cwndGain, attempts to send will block.
	cwndGain float64

	// pacingGain is the factor by which our sending rate changes
	// with each message; if pacingGain > 1, we will increase our
	// sending rate, if pacingGain < 1, we will decrease it.
	pacingGain float64

	// Earliest time at which it is appropriate to send.
	nextSendTime time.Time

	// pacingReady records that a pacing timer fired while no sender was
	// waiting. It is explicit because virtual time does not advance between
	// firing a zero-duration timer and preparing the next actor iteration.
	pacingReady bool

	sent          uint64    // Total data sent
	delivered     uint64    // Total data delivered & ACKed
	deliveredTime time.Time // Time of the last ACK we received.

	packetsInflight    uint64 // The number of packets currently in-flight.
	maxPacketsInflight uint64 // The maximum allowable packets in-flight.

	// If appLimitedUntil is not zero, it indicates that inflight()
	// was limited to the specified value *not* because our congestion
	// control logic decided that we should wait, but because the app
	// didn't have any more data to send.
	appLimitedUntil uint64

	// The state the flow is in
	state state

	// A clock, for measuring the current time.
	clock clock.Clock

	// This channel is used for testing; the whilePaused() method needs it.
	chPause chan struct{}

	// Closed after the actor exits. This makes snapshots after Release safe:
	// once done is closed, no goroutine can still mutate limiter state.
	done chan struct{}

	// Private test seam for deterministic ProbeBW initialization.
	randInt func() int

	// Private actor-owned GateNext reservation ledger. It is not exposed until
	// the capnp admission surface is complete.
	gate   gateState
	chGate chan gateEvent
}

// For testing purpoes; temporarily pauses the goroutine managing the limiter,
// and runs the callback. This lets us inspect the state of the limiter in
// tests without data races.
func (l *Limiter) whilePaused(f func()) {
	select {
	case l.chPause <- struct{}{}:
		defer func() {
			<-l.chPause
		}()
		f()
	case <-l.done:
		f()
	}
}

// inflight returns the total bytes in-flight
func (l *Limiter) inflight() uint64 {
	return l.sent - l.delivered
}

// NewLimiter returns a new BBR-based flow limiter. The clock is used to measure
// message resposne times. If nil is passed (typical, except for testing & debugging),
// the system clock will be used.
func NewLimiter(clk clock.Clock) *Limiter {
	return newLimiter(clk, limiterOptions{randInt: rand.Int})
}

type limiterOptions struct {
	randInt func() int
}

func newLimiter(clk clock.Clock, options limiterOptions) *Limiter {
	l := newLimiterState(clk, options)
	go l.run(l.ctx)
	return l
}

// newLimiterState constructs a limiter without starting its actor. Production
// must use newLimiter; deterministic tests may drive the synchronous state
// machine directly while retaining exclusive ownership.
func newLimiterState(clk clock.Clock, options limiterOptions) *Limiter {
	if clk == nil {
		clk = clock.System
	}
	if options.randInt == nil {
		options.randInt = rand.Int
	}
	now := clk.Now()
	ctx, cancel := context.WithCancel(context.Background())
	l := &Limiter{
		ctx:    ctx,
		cancel: cancel,

		chSend: make(chan sendRequest),
		chAck:  make(chan packetMeta),
		chGate: make(chan gateEvent),

		rtPropFilter: newRtPropFilter(),
		btlBwFilter:  newBtlBwFilter(),
		clock:        clk,

		nextSendTime:  now,
		deliveredTime: now,

		maxPacketsInflight: math.MaxUint64,

		chPause: make(chan struct{}),
		done:    make(chan struct{}),
		randInt: options.randInt,
	}
	l.changeState(&startupState{}, now)
	return l
}

// onAckAt should be invoked on each packetMeta when the acknowledgement for
// that packet is received.
func (l *Limiter) onAckAt(now time.Time, p packetMeta) {
	l.packetsInflight--
	rtt := now.Sub(p.SendTime)

	l.rtPropFilter.AddSample(rtPropSample{
		now: now,
		rtt: rtt,
	})

	l.delivered += p.Size
	l.deliveredTime = now

	deltaD := l.delivered - p.Delivered
	deltaT := l.deliveredTime.Sub(p.DeliveredTime)
	deliveryRate := bytesPerNs(float64(deltaD) / float64(deltaT.Nanoseconds()))

	if deliveryRate > l.btlBwFilter.Estimate || !p.AppLimited {
		l.btlBwFilter.AddSample(deliveryRate)
	}
	if l.appLimitedUntil > 0 {
		if p.Size >= l.appLimitedUntil {
			l.appLimitedUntil = 0
		} else {
			l.appLimitedUntil -= p.Size
		}
	}

	l.state.postAck(l, p, now)
}

func (l *Limiter) changeState(s state, now time.Time) {
	l.state = s
	l.state.initialize(l, now)
}
