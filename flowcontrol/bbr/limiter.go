package bbr

import (
	"context"
	"math"
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

func (l *Limiter) run(ctx context.Context) {
	for {

		// These channels may or may not be nil, depending on what
		// we want to listen for:
		var (
			// Fires when we cross the threshold past l.nextSendTime.
			timeToSend <-chan time.Time

			// Fires when someone wants to send.
			sendReqs <-chan sendRequest
		)

		bdp := l.computeBDP()

		now := l.clock.Now()

		if l.inflight() == 0 {
			// We can't let the general purpose logic handle this case,
			// because if we don't have good information yet nextSendTime
			// might be in the far future, and there is no ack on its way
			// to save us from our ignorance. Fortunately we always want
			// to send if there's nothing on the wire, so just do it.
			sendReqs = l.chSend
		} else if l.packetsInflight >= l.maxPacketsInflight ||
			(bdp > 0 && float64(l.inflight()) >= l.cwndGain*bdp) {
			// We're at our threshold; wait for an ack,
			// but ignore other signals.
			//
			// Note: bdp == 0 means we just don't have any data yet,
			// so we filter out that scenario.
		} else if now.After(l.nextSendTime) {
			// We're bottlnecked on the app, not the path;
			// record that we're app limited, and don't watch the timer,
			// but do watch the send queue:
			l.appLimitedUntil = l.inflight()
			sendReqs = l.chSend
		} else {
			// App is sending fast enough for congestion
			// control to be active; wait until nextSendTime
			// before servicing a request:
			sleep := l.nextSendTime.Sub(now)
			timeToSend = l.clock.NewTimer(sleep).Chan()
		}

		select {
		case <-ctx.Done():
			return
		case p := <-l.chAck:
			l.onAck(p)
		case <-timeToSend:
			l.trySend()
		case req := <-sendReqs:
			now := l.clock.Now()
			l.doSend(now, req)
		case <-l.chPause:
			l.chPause <- struct{}{}
		}
	}
}

func (l *Limiter) trySend() {
	select {
	case req := <-l.chSend:
		now := l.clock.Now()
		l.doSend(now, req)
	default:
		l.appLimitedUntil = l.inflight()
	}
}

func (l *Limiter) doSend(now time.Time, req sendRequest) {
	l.packetsInflight++
	p := packetMeta{
		Size:          req.size,
		AppLimited:    l.appLimitedUntil > 0,
		SendTime:      now,
		Delivered:     l.delivered,
		DeliveredTime: l.deliveredTime,
	}
	l.nextSendTime = now.Add(time.Duration(
		float64(req.size) / (l.pacingGain * float64(l.btlBwFilter.Estimate)),
	))
	l.sent += req.size
	req.replyChan <- p
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
}

// For testing purpoes; temporarily pauses the goroutine managing the limiter,
// and runs the callback. This lets us inspect the state of the limiter in
// tests without data races.
func (l *Limiter) whilePaused(f func()) {
	l.chPause <- struct{}{}
	defer func() {
		<-l.chPause
	}()
	f()
}

// inflight returns the total bytes in-flight
func (l *Limiter) inflight() uint64 {
	return l.sent - l.delivered
}

// NewLimiter returns a new BBR-based flow limiter. The clock is used to measure
// message resposne times. If nil is passed (typical, except for testing & debugging),
// the system clock will be used.
func NewLimiter(clk clock.Clock) *Limiter {
	if clk == nil {
		clk = clock.System
	}
	now := clk.Now()
	ctx, cancel := context.WithCancel(context.Background())
	l := &Limiter{
		ctx:    ctx,
		cancel: cancel,

		chSend: make(chan sendRequest),
		chAck:  make(chan packetMeta),

		rtPropFilter: newRtPropFilter(),
		btlBwFilter:  newBtlBwFilter(),
		clock:        clk,

		nextSendTime:  now,
		deliveredTime: now,

		maxPacketsInflight: math.MaxUint64,

		chPause: make(chan struct{}),
	}
	l.changeState(&startupState{})
	go l.run(ctx)
	return l
}

// onAck should be invoked on each packetMeta when the acknowledgement for
// that packet is received.
func (l *Limiter) onAck(p packetMeta) {
	l.packetsInflight--
	now := l.clock.Now()
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
		l.appLimitedUntil -= p.Size
	}

	l.state.postAck(l, p, now)
}

func (l *Limiter) changeState(s state) {
	l.state = s
	l.state.initialize(l)
}
