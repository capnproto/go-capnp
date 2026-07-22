package flowcontrol

import (
	"context"
)

var (
	// A flow limiter which does not actually limit anything; messages will be
	// sent as fast as possible.
	NopLimiter GateNextFlowLimiter = nopLimiter{}
)

type nopLimiter struct{}

func (nopLimiter) StartMessage(context.Context, uint64) (func(), error) {
	return func() {}, nil
}

func (nopLimiter) Release() {}

func (nopLimiter) GateNext() GateNextController { return nopGate{} }

type nopGate struct{}

func (nopGate) CommitMessage(uint64) (func(context.Context) error, func(MessageOutcomeKind, error)) {
	return func(context.Context) error { return nil }, func(MessageOutcomeKind, error) {}
}

func (nopGate) Poison(error) {}
