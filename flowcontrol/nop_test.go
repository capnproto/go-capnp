package flowcontrol

import (
	"context"
	"testing"
)

func TestNopLimiterGateNext(t *testing.T) {
	controller := NopLimiter.GateNext()
	waitNext, complete := controller.CommitMessage(1024)
	if err := waitNext(context.Background()); err != nil {
		t.Fatalf("waitNext() = %v; want nil", err)
	}
	complete(MessageOutcomeSucceeded, nil)
	controller.Poison(context.Canceled)

	// Nop's controller is deliberately zero-cost even after Poison.
	waitNext, complete = controller.CommitMessage(1)
	if err := waitNext(context.Background()); err != nil {
		t.Fatalf("waitNext() after Poison = %v; want nil", err)
	}
	complete(MessageOutcomeFatal, context.Canceled)
}
