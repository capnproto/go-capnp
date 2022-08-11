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
}

// A state tracks the state used by BBR to manage flow control.
type state struct {
	// Filters for estimating the round trip propagation time and
	// bottleneck bandwidth, respectively:
	rtPropFilter minFilter[time.Duration]
	btlBwFilter  maxFilter[int64]

	delivered     int64
	deliveredTime time.Time

	appLimitedUntil int64

	clock clock.Clock
}

// onAck should be invoked on each packetMeta when the acknowledgement for
// that packet is received.
func (s *state) onAck(p packetMeta) {
	now := s.clock.Now()
	rtt := now.Sub(p.SendTime)

	s.rtPropFilter.AddSample(rtt)

	s.delivered += p.Size
	s.deliveredTime = now

	deliveryRate := TODO

	if deliveryRate > s.btlBwFilter.estimate || !p.AppLimited {
		s.btlBwFilter.AddSample(deliveryRate)
	}
	if s.appLimitedUntil > 0 {
		s.appLimitedUntil -= p.Size
	}
}
