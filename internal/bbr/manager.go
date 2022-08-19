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

	// value of corresponding fields in Manager when this packet was sent:
	Delivered     int64
	DeliveredTime time.Time
}

type sendRequest struct {
	size      int64
	replyChan chan<- packetMeta
}

func (m *Manager) StartMessage(ctx context.Context, size uint64) (gotResponse func(), err error) {
	if size > math.MaxInt64 {
		panic("TODO: overflow")
	}
	replyChan := make(chan packetMeta)
	select {
	case <-ctx.Done():
		return func() {}, ctx.Err()
	case <-m.ctx.Done():
		return func() {}, m.ctx.Err()
	case m.chSend <- sendRequest{size: int64(size), replyChan: replyChan}:
		pm := <-replyChan
		return func() {
			select {
			case <-m.ctx.Done():
			case m.chAck <- pm:
			}
		}, nil
	}
}

func (m *Manager) Release() {
	m.cancel()
}

func (m *Manager) run(ctx context.Context) {
	for {
		// These channels may or may not be nil, depending on what
		// we want to listen for:
		var (
			// Fires when we cross the threshold past m.nextSendTime.
			timeToSend <-chan time.Time

			// Fires when someone wants to send.
			sendReqs <-chan sendRequest
		)

		bdp := m.btlBwFilter.Estimate * m.rtPropFilter.Estimate.Nanoseconds()

		if m.inflight >= int64(m.cwndGain*float64(bdp)) {
			// We're at our threshold; wait for an ack,
			// but ignore other signals.
		} else if m.isAppLimited() {
			// Last run through we didn't have enough
			// to send, so we should just wait until we do;
			// don't watch the timer, but do watch the send
			// queue:
			sendReqs = m.chSend
		} else {
			// App is sending fast enough for congestion
			// control to be active; monitor the timer
			// and wait until it fires before servicing
			// a request:
			timeToSend = m.timer.Chan()
		}

		select {
		case <-ctx.Done():
			return
		case p := <-m.chAck:
			m.onAck(p)
		case <-timeToSend:
			m.trySend(ctx)
		case req := <-sendReqs:
			now := m.clock.Now()
			m.doSend(now, req)
		}
	}
}

func (m *Manager) isAppLimited() bool {
	return m.appLimitedUntil > 0
}

func (m *Manager) trySend(ctx context.Context) {
	select {
	case <-ctx.Done():
	case req := <-m.chSend:
		now := m.clock.Now()
		m.doSend(now, req)
	default:
		m.appLimitedUntil = m.inflight
	}
}

func (m *Manager) doSend(now time.Time, req sendRequest) {
	p := packetMeta{
		Size:          req.size,
		AppLimited:    m.isAppLimited(),
		SendTime:      now,
		Delivered:     m.delivered,
		DeliveredTime: m.deliveredTime,
	}
	if !m.timer.Stop() {
		<-m.timer.Chan()
	}
	nextSleep := time.Duration(float64(req.size) / (m.pacingGain * float64(m.btlBwFilter.Estimate)))
	m.timer.Reset(nextSleep)
	req.replyChan <- p
	m.nextSendTime = now.Add(nextSleep)
}

// A Manager manages a BBR flow.
type Manager struct {
	ctx    context.Context
	cancel context.CancelFunc
	chAck  chan packetMeta
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

	// Total bytes in-flight
	inflight int64

	// cwndGain (Congestion WiNDow; standard terminology) is a factor
	// by which we may exceed BDP. i.e. if inflight is greater than
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

	delivered     int64     // Total data delivered/ACKed
	deliveredTime time.Time // Time of the last ACK we received.

	// If appLimitedUntil is not zero, it indicates that inflight
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

func NewManager(clock clock.Clock) Manager {
	ctx, cancel := context.WithCancel(context.Background())
	m := Manager{
		ctx:    ctx,
		cancel: cancel,

		chSend: make(chan sendRequest),
		chAck:  make(chan packetMeta),

		rtPropFilter: newRtPropFilter(),
		btlBwFilter:  newBtlBwFilter(),
		clock:        clock,
		// TODO: timer.
	}
	m.changeState(&startupState{})
	return m
}

// onAck should be invoked on each packetMeta when the acknowledgement for
// that packet is received.
func (m *Manager) onAck(p packetMeta) {
	now := m.clock.Now()
	m.state.preAck(m, p, now)

	rtt := now.Sub(p.SendTime)

	m.rtPropFilter.AddSample(rtPropSample{
		now: now,
		rtt: rtt,
	})

	m.delivered += p.Size
	m.deliveredTime = now

	deliveryRate := (m.delivered - p.Delivered) /
		(m.deliveredTime.Sub(p.DeliveredTime)).Nanoseconds()

	if deliveryRate > m.btlBwFilter.Estimate || !p.AppLimited {
		m.btlBwFilter.AddSample(deliveryRate)
	}
	if m.appLimitedUntil > 0 {
		m.appLimitedUntil -= p.Size
	}

	m.state.postAck(m, p, now)
}

func (m *Manager) changeState(s state) {
	m.state = s
	m.state.initialize(m)
}
