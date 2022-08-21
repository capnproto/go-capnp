package bbr

import (
	"context"
	"math"
	"time"

	"capnproto.org/go/capnp/v3/internal/clock"
)

// A packetMeta contains metadata about a packet that was sent.
type packetMeta struct {
	SendTime time.Time // The time at which the packet was sent.
	Size     int64     // The size of the packet.

	// Whether the connection flow was app-limited when this packet
	// was sent:
	AppLimited bool

	// value of corresponding fields in Limiter when this packet was sent:
	Delivered     int64
	DeliveredTime time.Time
}

type sendRequest struct {
	size      int64
	replyChan chan<- packetMeta
}

func (l *Limiter) StartMessage(ctx context.Context, size uint64) (gotResponse func(), err error) {
	if size > math.MaxInt64 {
		panic("TODO: overflow")
	}
	replyChan := make(chan packetMeta)
	select {
	case <-ctx.Done():
		return func() {}, ctx.Err()
	case <-l.ctx.Done():
		return func() {}, l.ctx.Err()
	case l.chSend <- sendRequest{size: int64(size), replyChan: replyChan}:
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

		bdp := l.btlBwFilter.Estimate * l.rtPropFilter.Estimate.Nanoseconds()

		if bdp > 0 && l.inflight() >= int64(l.cwndGain*float64(bdp)) {
			// We're at our threshold; wait for an ack,
			// but ignore other signals.
			//
			// Note: if bdp <= 0, that means we overflowed, so BDP is
			// estimated to be *extremely* large -- larger any possible
			// value of inflight(). In reality, this probably just means
			// we don't have data yet, but the result is the same: just
			// send.
		} else if l.isAppLimited() {
			// Last run through we didn't have enough
			// to send, so we should just wait until we do;
			// don't watch the timer, but do watch the send
			// queue:
			sendReqs = l.chSend
		} else {
			// App is sending fast enough for congestion
			// control to be active; monitor the timer
			// and wait until it fires before servicing
			// a request:
			timeToSend = l.timer.Chan()
		}

		select {
		case <-ctx.Done():
			return
		case p := <-l.chAck:
			l.onAck(p)
		case <-timeToSend:
			l.trySend(ctx)
		case req := <-sendReqs:
			if !l.timer.Stop() {
				<-l.timer.Chan()
			}

			now := l.clock.Now()
			l.doSend(now, req)
		}
	}
}

func (l *Limiter) isAppLimited() bool {
	return l.appLimitedUntil > 0
}

func (l *Limiter) trySend(ctx context.Context) {
	select {
	case <-ctx.Done():
	case req := <-l.chSend:
		now := l.clock.Now()
		l.doSend(now, req)
	default:
		l.appLimitedUntil = l.inflight()
	}
}

func (l *Limiter) doSend(now time.Time, req sendRequest) {
	p := packetMeta{
		Size:          req.size,
		AppLimited:    l.isAppLimited(),
		SendTime:      now,
		Delivered:     l.delivered,
		DeliveredTime: l.deliveredTime,
	}
	nextSleep := time.Duration(float64(req.size) / (l.pacingGain * float64(l.btlBwFilter.Estimate)))
	l.timer.Reset(nextSleep)
	req.replyChan <- p
	l.nextSendTime = now.Add(nextSleep)
	l.sent += req.size
}

// A Limiter manages a BBR flow.
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
	// TODO: maybe use a larger unit of time, to save bits for the
	// actual data.
	rtPropFilter rtPropFilter
	btlBwFilter  btlBwFilter

	// cwndGain (Congestion WiNDow; standard terminology) is a factor
	// by which we may exceed BDP. i.e. if inflight() is greater than
	// estimated BDP * cwndGain, attempts to send will block.
	//
	// TODO(perf): The BBR paper says that the reason this is
	// needed (i.e. isn't always just 1) is to deal with delayed and
	// aggregated ACKs -- a feature of some TCP implementations, but
	// not something that any capnp implementation does. So do we
	// actually need this? would dropping it be wrothwhile for
	// the latency improvement?
	cwndGain float64

	// pacingGain is the factor by which our sending rate changes
	// with each message; if pacingGain > 1, we will increase our
	// sending rate, if pacingGain < 1, we will decrease it.
	pacingGain float64

	// Earliest time at which it is appropriate to send.
	nextSendTime time.Time

	sent          int64     // Total data sent
	delivered     int64     // Total data delivered & ACKed
	deliveredTime time.Time // Time of the last ACK we received.

	// If appLimitedUntil is not zero, it indicates that inflight()
	// was limited to the specified value *not* because our congestion
	// control logic decidecd that we should wait, but because the app
	// didn't have any more data to send.
	appLimitedUntil int64

	// The state the flow is in
	state state

	// A clock, for measuring the current time.
	clock clock.Clock

	// A timer, to notify us when it's time to send a packet.
	timer clock.Timer
}

// inflight returns the total bytes in-flight
func (l *Limiter) inflight() int64 {
	return l.sent - l.delivered
}

func NewLimiter(clock clock.Clock) Limiter {
	ctx, cancel := context.WithCancel(context.Background())
	l := Limiter{
		ctx:    ctx,
		cancel: cancel,

		chSend: make(chan sendRequest),
		chAck:  make(chan packetMeta),

		rtPropFilter: newRtPropFilter(),
		btlBwFilter:  newBtlBwFilter(),
		clock:        clock,

		// Set the timer to go off for the first time immediately:
		timer: clock.NewTimer(0),
	}
	l.changeState(&startupState{})
	go l.run(ctx)
	return l
}

// onAck should be invoked on each packetMeta when the acknowledgement for
// that packet is received.
func (l *Limiter) onAck(p packetMeta) {
	now := l.clock.Now()
	l.state.preAck(l, p, now)

	rtt := now.Sub(p.SendTime)

	l.rtPropFilter.AddSample(rtPropSample{
		now: now,
		rtt: rtt,
	})

	l.delivered += p.Size
	l.deliveredTime = now

	deliveryRate := (l.delivered - p.Delivered) /
		(l.deliveredTime.Sub(p.DeliveredTime)).Nanoseconds()

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
