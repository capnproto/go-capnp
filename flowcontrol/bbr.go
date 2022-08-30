package flowcontrol

import (
	"capnproto.org/go/capnp/v3/internal/bbr"
	"capnproto.org/go/capnp/v3/internal/clock"
)

// NewBBR returns a new FlowLimiter that adaptively determines the correct rate,
// using the BBR algorithm described at: https://queue.acm.org/detail.cfm?id=3022184
func NewBBR() FlowLimiter {
	return bbr.NewLimiter(clock.System)
}
