package bbr

import (
	"time"

	"capnproto.org/go/capnp/v3/internal/clock"
)

type packet struct {
	SendTime   time.Time
	Size       int64
	AppLimited bool
}

// A state tracks the state used by BBR to manage flow control.
type state struct {
	rtPropFilter rtPropFilter
	btlBwFilter  btlBwFilter

	delivered     int64
	deliveredTime time.Time

	appLimitedUntil int64

	clock clock.Clock
}

func (s *state) onAck(p packet) {
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
