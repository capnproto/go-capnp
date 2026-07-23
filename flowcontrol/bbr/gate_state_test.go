package bbr

import (
	"context"
	"errors"
	"sync"
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

func TestGateEventsAreOwnedByLimiterActor(t *testing.T) {
	lim := NewLimiter(nil)
	defer lim.Release()

	a, err := lim.gateCommit(10)
	if !assert.NoError(t, err) {
		return
	}
	assert.NotNil(t, a)
	replay, err := lim.gateComplete(a, flowcontrol.MessageOutcomeAbortedBeforeEnqueue, nil)
	assert.NoError(t, err)
	assert.True(t, replay)

	lim.whilePaused(func() {
		assert.Empty(t, lim.gate.reservations)
	})
}

func TestGateEventsFailAfterLimiterRelease(t *testing.T) {
	lim := NewLimiter(nil)
	lim.Release()
	<-lim.done

	_, err := lim.gateCommit(10)
	assert.ErrorIs(t, err, context.Canceled)
}

func TestGateEventsPoisonRejectsCommitAndReplay(t *testing.T) {
	lim := NewLimiter(nil)
	defer lim.Release()
	a, err := lim.gateCommit(10)
	if !assert.NoError(t, err) {
		return
	}
	b, err := lim.gateCommit(20)
	if !assert.NoError(t, err) {
		return
	}
	poison := errors.New("poison")
	replay, err := lim.gateComplete(a, flowcontrol.MessageOutcomeFatal, poison)
	assert.NoError(t, err)
	assert.False(t, replay)
	_, err = lim.gateCommit(30)
	assert.ErrorIs(t, err, poison)
	replay, err = lim.gateComplete(b, flowcontrol.MessageOutcomeAbortedBeforeEnqueue, nil)
	assert.NoError(t, err)
	assert.False(t, replay)
}

func TestGateEventsSerializeConcurrentCommits(t *testing.T) {
	lim := NewLimiter(nil)
	defer lim.Release()
	const n = 16
	ids := make(chan uint64, n)
	errs := make(chan error, n)
	var wg sync.WaitGroup
	for range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			r, err := lim.gateCommit(1)
			if err == nil {
				ids <- r.id
			}
			errs <- err
		}()
	}
	wg.Wait()
	close(ids)
	close(errs)
	seen := make(map[uint64]struct{}, n)
	for err := range errs {
		assert.NoError(t, err)
	}
	for id := range ids {
		if _, ok := seen[id]; !assert.False(t, ok, "duplicate id %d", id) {
			continue
		}
		seen[id] = struct{}{}
	}
	assert.Len(t, seen, n)
}

func TestGateEventsInvalidKindPoisonsLedger(t *testing.T) {
	lim := NewLimiter(nil)
	defer lim.Release()
	_, err := lim.gateEvent(gateEvent{kind: gateEventKind(99)})
	assert.NoError(t, err)
	lim.whilePaused(func() {
		assert.EqualError(t, lim.gate.poison, "bbr: invalid gate-next event")
	})
}

func TestGateWaitCancellationDoesNotConsumePermission(t *testing.T) {
	lim := NewLimiter(nil)
	defer lim.Release()
	r, err := lim.gateCommit(10)
	if !assert.NoError(t, err) {
		return
	}
	a, err := lim.gateStartWait(context.Background(), r)
	if !assert.NoError(t, err) {
		return
	}
	cancelErr := errors.New("canceled")
	assert.ErrorIs(t, lim.gateCancelWait(a, cancelErr), cancelErr)
	assert.ErrorIs(t, <-a.result, cancelErr)

	retry, err := lim.gateStartWait(context.Background(), r)
	if !assert.NoError(t, err) {
		return
	}
	granted, err := lim.gateGrant(r)
	assert.NoError(t, err)
	assert.True(t, granted)
	assert.NoError(t, <-retry.result)
}

func TestGatePoisonWakesWaitingAttempt(t *testing.T) {
	lim := NewLimiter(nil)
	defer lim.Release()
	r, err := lim.gateCommit(10)
	if !assert.NoError(t, err) {
		return
	}
	a, err := lim.gateStartWait(context.Background(), r)
	if !assert.NoError(t, err) {
		return
	}
	poison := errors.New("poison")
	assert.NoError(t, lim.gatePoison(poison))
	assert.ErrorIs(t, <-a.result, poison)
}

func TestGateFatalCompletionWakesWaitingAttempt(t *testing.T) {
	lim := NewLimiter(nil)
	defer lim.Release()
	r, err := lim.gateCommit(10)
	if !assert.NoError(t, err) {
		return
	}
	a, err := lim.gateStartWait(context.Background(), r)
	if !assert.NoError(t, err) {
		return
	}
	poison := errors.New("fatal")
	_, err = lim.gateComplete(r, flowcontrol.MessageOutcomeFatal, poison)
	assert.NoError(t, err)
	assert.ErrorIs(t, <-a.result, poison)
}
