package bbr

import (
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

// A Manager manages a BBR flow.
type Manager struct {
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

	// A clock, for measuring the current time.
	clock clock.Clock
}

func (m *Manager) send(size int64) (_ packetMeta, ok bool) {
	now := m.clock.Now()
	bdp := m.btlBwFilter.Estimate * m.rtPropFilter.Estimate.Nanoseconds()
	if m.inflight >= int64(m.cwndGain*float64(bdp)) {
		return packetMeta{}, false
	}
	if now.After(m.nextSendTime) {
	}
	panic("TODO")
}

// onAck should be invoked on each packetMeta when the acknowledgement for
// that packet is received.
func (m *Manager) onAck(p packetMeta) {
	now := m.clock.Now()
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
}
