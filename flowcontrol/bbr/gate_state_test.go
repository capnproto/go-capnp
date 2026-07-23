package bbr

import (
	"errors"
	"testing"

	"capnproto.org/go/capnp/v3/flowcontrol"

	"github.com/stretchr/testify/assert"
)

func TestGateStateAbortRemovesReservationAndPreservesSuffix(t *testing.T) {
	var g gateState
	a := g.commit(10)
	b := g.commit(20)

	assert.True(t, g.complete(a, flowcontrol.MessageOutcomeAbortedBeforeEnqueue, nil))
	if assert.Len(t, g.reservations, 1) {
		assert.Same(t, b, g.reservations[0])
		assert.Equal(t, gateReservationProvisional, b.state)
	}
}

func TestGateStateAcknowledgedHolesCompactOnlyAtPrefix(t *testing.T) {
	var g gateState
	a := g.commit(10)
	b := g.commit(20)
	c := g.commit(30)

	g.complete(b, flowcontrol.MessageOutcomeSucceeded, nil)
	if assert.Len(t, g.reservations, 3) {
		assert.Same(t, a, g.reservations[0])
		assert.Same(t, b, g.reservations[1])
		assert.Same(t, c, g.reservations[2])
	}
	g.complete(a, flowcontrol.MessageOutcomeSucceeded, nil)
	if assert.Len(t, g.reservations, 1) {
		assert.Same(t, c, g.reservations[0])
	}
}

func TestGateStateFatalPreservesFirstPoison(t *testing.T) {
	var g gateState
	a := g.commit(10)
	b := g.commit(20)
	first := errors.New("first")

	assert.False(t, g.complete(a, flowcontrol.MessageOutcomeFatal, first))
	assert.False(t, g.complete(b, flowcontrol.MessageOutcomeFailedAfterEnqueue, errors.New("second")))
	assert.ErrorIs(t, g.poison, first)
}

func TestGateStateTerminalCompletionIsIdempotent(t *testing.T) {
	var g gateState
	a := g.commit(10)

	assert.True(t, g.complete(a, flowcontrol.MessageOutcomeAbortedBeforeEnqueue, nil))
	assert.False(t, g.complete(a, flowcontrol.MessageOutcomeAbortedBeforeEnqueue, nil))
	assert.Empty(t, g.reservations)
}

func TestGateStateContradictoryTerminalOutcomeIsIgnored(t *testing.T) {
	var g gateState
	a := g.commit(10)

	assert.False(t, g.complete(a, flowcontrol.MessageOutcomeSucceeded, nil))
	assert.False(t, g.complete(a, flowcontrol.MessageOutcomeFatal, errors.New("late")))
	assert.NoError(t, g.poison)
	assert.False(t, g.complete(a, flowcontrol.MessageOutcomeSucceeded, nil))
}

func TestGateStateNilAndPoisonedCompletionAreIdempotent(t *testing.T) {
	var g gateState
	assert.False(t, g.complete(nil, flowcontrol.MessageOutcomeFatal, errors.New("ignored")))

	a := g.commit(10)
	assert.False(t, g.complete(a, flowcontrol.MessageOutcomeFatal, nil))
	assert.EqualError(t, g.poison, "bbr: gate-next poisoned")
	assert.False(t, g.complete(a, flowcontrol.MessageOutcomeFatal, errors.New("late")))
	assert.EqualError(t, g.poison, "bbr: gate-next poisoned")
}

func TestGateStateInvalidOutcomePoisonsReservation(t *testing.T) {
	var g gateState
	a := g.commit(10)

	assert.False(t, g.complete(a, flowcontrol.MessageOutcomeUnknown, nil))
	assert.Equal(t, gateReservationPoisoned, a.state)
	assert.EqualError(t, g.poison, "bbr: invalid gate-next terminal outcome")
	assert.False(t, g.complete(a, flowcontrol.MessageOutcomeAbortedBeforeEnqueue, nil))
	assert.Same(t, a, g.reservations[0])
}
