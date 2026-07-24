package bbr

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

	"capnproto.org/go/capnp/v3/exp/clock"
	"capnproto.org/go/capnp/v3/flowcontrol"

	"github.com/stretchr/testify/assert"
)

func TestGateStateAbortRemovesReservationAndPreservesSuffix(t *testing.T) {
	var g gateState
	a := g.commit(10)
	b := g.commit(20)

	assert.Equal(t, gateActionUnsend, g.complete(a, flowcontrol.MessageOutcomeAbortedBeforeEnqueue, nil))
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

	assert.Equal(t, gateActionNone, g.complete(a, flowcontrol.MessageOutcomeFatal, first))
	assert.Equal(t, gateActionNone, g.complete(b, flowcontrol.MessageOutcomeFailedAfterEnqueue, errors.New("second")))
	assert.ErrorIs(t, g.poison, first)
}

func TestGateStateTerminalCompletionIsIdempotent(t *testing.T) {
	var g gateState
	a := g.commit(10)

	assert.Equal(t, gateActionUnsend, g.complete(a, flowcontrol.MessageOutcomeAbortedBeforeEnqueue, nil))
	assert.Equal(t, gateActionNone, g.complete(a, flowcontrol.MessageOutcomeAbortedBeforeEnqueue, nil))
	assert.Empty(t, g.reservations)
}

func TestGateStateContradictoryTerminalOutcomeIsIgnored(t *testing.T) {
	var g gateState
	a := g.commit(10)

	assert.Equal(t, gateActionAck, g.complete(a, flowcontrol.MessageOutcomeSucceeded, nil))
	assert.Equal(t, gateActionNone, g.complete(a, flowcontrol.MessageOutcomeFatal, errors.New("late")))
	assert.NoError(t, g.poison)
	assert.Equal(t, gateActionNone, g.complete(a, flowcontrol.MessageOutcomeSucceeded, nil))
}

func TestGateStateNilAndPoisonedCompletionAreIdempotent(t *testing.T) {
	var g gateState
	assert.Equal(t, gateActionNone, g.complete(nil, flowcontrol.MessageOutcomeFatal, errors.New("ignored")))

	a := g.commit(10)
	assert.Equal(t, gateActionNone, g.complete(a, flowcontrol.MessageOutcomeFatal, nil))
	assert.EqualError(t, g.poison, "bbr: gate-next poisoned")
	assert.Equal(t, gateActionNone, g.complete(a, flowcontrol.MessageOutcomeFatal, errors.New("late")))
	assert.EqualError(t, g.poison, "bbr: gate-next poisoned")
}

func TestGateStateInvalidOutcomePoisonsReservation(t *testing.T) {
	var g gateState
	a := g.commit(10)

	assert.Equal(t, gateActionNone, g.complete(a, flowcontrol.MessageOutcomeUnknown, nil))
	assert.Equal(t, gateReservationPoisoned, a.state)
	assert.EqualError(t, g.poison, "bbr: invalid gate-next terminal outcome")
	assert.Equal(t, gateActionNone, g.complete(a, flowcontrol.MessageOutcomeAbortedBeforeEnqueue, nil))
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
	action, err := lim.gateComplete(a, flowcontrol.MessageOutcomeAbortedBeforeEnqueue, nil)
	assert.NoError(t, err)
	assert.Equal(t, gateActionUnsend, action)

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
	action, err := lim.gateComplete(a, flowcontrol.MessageOutcomeFatal, poison)
	assert.NoError(t, err)
	assert.Equal(t, gateActionNone, action)
	_, err = lim.gateCommit(30)
	assert.ErrorIs(t, err, poison)
	action, err = lim.gateComplete(b, flowcontrol.MessageOutcomeAbortedBeforeEnqueue, nil)
	assert.NoError(t, err)
	assert.Equal(t, gateActionNone, action)
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

func TestGateWaitResultGrantWinsCanceledCaller(t *testing.T) {
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
	granted, err := lim.gateGrant(r)
	if !assert.NoError(t, err) || !assert.True(t, granted) {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	assert.NoError(t, lim.gateWaitResult(ctx, a))
}

func TestGateWaitRegistrationRejectsDuplicateAndPoisonedReservation(t *testing.T) {
	lim := NewLimiter(nil)
	defer lim.Release()
	r, err := lim.gateCommit(10)
	if !assert.NoError(t, err) {
		return
	}
	_, err = lim.gateStartWait(context.Background(), r)
	assert.NoError(t, err)
	_, err = lim.gateStartWait(context.Background(), r)
	assert.EqualError(t, err, "bbr: gate-next permission already has a waiter")

	terminal, err := lim.gateCommit(20)
	if !assert.NoError(t, err) {
		return
	}
	fatal := errors.New("fatal")
	_, err = lim.gateComplete(terminal, flowcontrol.MessageOutcomeFatal, fatal)
	assert.NoError(t, err)
	_, err = lim.gateStartWait(context.Background(), terminal)
	assert.ErrorIs(t, err, fatal)
}

func TestGateWaitRemainsUsableAfterAcknowledgement(t *testing.T) {
	lim := NewLimiter(nil)
	defer lim.Release()
	r, err := lim.gateCommit(10)
	if !assert.NoError(t, err) {
		return
	}
	_, err = lim.gateComplete(r, flowcontrol.MessageOutcomeSucceeded, nil)
	if !assert.NoError(t, err) {
		return
	}

	first, err := lim.gateStartWait(context.Background(), r)
	if !assert.NoError(t, err) {
		return
	}
	assert.NoError(t, <-first.result)
	granted, err := lim.gateGrant(r)
	assert.NoError(t, err)
	assert.False(t, granted)

	// Ticket forwarding may invoke the same published permission again after
	// another caller won the grant. It must reuse that grant, not consume a
	// second admission opportunity.
	retry, err := lim.gateStartWait(context.Background(), r)
	if !assert.NoError(t, err) {
		return
	}
	assert.NoError(t, <-retry.result)
	granted, err = lim.gateGrant(r)
	assert.NoError(t, err)
	assert.False(t, granted)
}

func TestGateControllerIsStable(t *testing.T) {
	lim := NewLimiter(nil)
	defer lim.Release()

	assert.Same(t, lim.GateNext(), lim.GateNext())
}

func TestGatePendingGrantClosesWindowAndConvertsAtCommit(t *testing.T) {
	now := time.Unix(1, 0)
	lim := newLimiterState(clock.NewManual(now), limiterOptions{})
	lim.maxPacketsInflight = 1
	first := applyGateEvent(lim, now, gateEvent{kind: gateCommitEvent, size: 10}).reservation
	waiter := newGateWaitAttempt(first)
	registered := applyGateEvent(lim, now, gateEvent{kind: gateWaitEvent, waiter: waiter})
	if !assert.NoError(t, registered.err) {
		return
	}
	assert.False(t, lim.prepareStep(now).acceptSend)

	ack := applyGateEvent(lim, now.Add(time.Millisecond), gateEvent{
		kind:        gateCompleteEvent,
		reservation: first,
		outcome:     flowcontrol.MessageOutcomeSucceeded,
	})
	assert.Equal(t, gateActionAck, ack.action)
	assert.True(t, lim.prepareStep(now.Add(time.Millisecond)).acceptSend)
	assert.True(t, lim.gate.grantNextWait())
	assert.NoError(t, <-waiter.result)
	assert.Len(t, lim.gate.pendingGrants, 1)
	assert.False(t, lim.prepareStep(now.Add(time.Millisecond)).acceptSend)

	second := applyGateEvent(lim, now.Add(time.Millisecond), gateEvent{
		kind: gateCommitEvent,
		size: 20,
	}).reservation
	assert.NotNil(t, second)
	assert.Empty(t, lim.gate.pendingGrants)
	assert.True(t, first.successorUsed)
	assert.Equal(t, uint64(20), lim.inflight())
	assert.Equal(t, uint64(1), lim.packetsInflight)
}

func TestGateProductionGrantWaitsForOpenWindow(t *testing.T) {
	lim := NewLimiter(nil)
	defer lim.Release()
	lim.whilePaused(func() {
		lim.maxPacketsInflight = 1
	})
	r, err := lim.gateCommit(10)
	if !assert.NoError(t, err) {
		return
	}
	waiter, err := lim.gateStartWait(context.Background(), r)
	if !assert.NoError(t, err) {
		return
	}
	select {
	case err := <-waiter.result:
		t.Fatalf("gate granted with a full packet window: %v", err)
	default:
	}

	_, err = lim.gateComplete(r, flowcontrol.MessageOutcomeSucceeded, nil)
	assert.NoError(t, err)
	assert.NoError(t, <-waiter.result)
}

func TestGateDefinitelyUnsentReplaysAndWakesSuccessor(t *testing.T) {
	lim := NewLimiter(nil)
	defer lim.Release()
	lim.whilePaused(func() {
		lim.maxPacketsInflight = 1
	})
	r, err := lim.gateCommit(10)
	if !assert.NoError(t, err) {
		return
	}
	waiter, err := lim.gateStartWait(context.Background(), r)
	if !assert.NoError(t, err) {
		return
	}
	select {
	case err := <-waiter.result:
		t.Fatalf("gate granted with a full packet window: %v", err)
	default:
	}

	action, err := lim.gateComplete(r, flowcontrol.MessageOutcomeAbortedBeforeEnqueue, nil)
	assert.NoError(t, err)
	assert.Equal(t, gateActionUnsend, action)
	assert.NoError(t, <-waiter.result)

	next, err := lim.gateCommit(20)
	assert.NoError(t, err)
	assert.NotNil(t, next)
	lim.whilePaused(func() {
		assert.Empty(t, lim.gate.pendingGrants)
		assert.Equal(t, uint64(20), lim.inflight())
		assert.Equal(t, uint64(1), lim.packetsInflight)
	})
}

func TestGateControllerWaitSurvivesCompletionBeforeRegistration(t *testing.T) {
	lim := NewLimiter(nil)
	defer lim.Release()
	controller := lim.GateNext()
	waitNext, complete := controller.CommitMessage(10)
	complete(flowcontrol.MessageOutcomeSucceeded, nil)

	assert.NoError(t, waitNext(context.Background()))
	assert.NoError(t, waitNext(context.Background()))

	_, completeNext := controller.CommitMessage(20)
	completeNext(flowcontrol.MessageOutcomeAbortedBeforeEnqueue, nil)
}

func TestGateGrantRejectsMissingAndPoisonedWaiter(t *testing.T) {
	lim := NewLimiter(nil)
	defer lim.Release()
	r, err := lim.gateCommit(10)
	if !assert.NoError(t, err) {
		return
	}
	granted, err := lim.gateGrant(r)
	assert.NoError(t, err)
	assert.False(t, granted)
	assert.NoError(t, lim.gatePoison(errors.New("poison")))
	granted, err = lim.gateGrant(r)
	assert.NoError(t, err)
	assert.False(t, granted)
}

func TestGateCommitChargesAndSuccessAcknowledgesBBR(t *testing.T) {
	now := time.Unix(1, 0)
	lim := newLimiterState(clock.NewManual(now), limiterOptions{})
	commit := applyGateEvent(lim, now, gateEvent{kind: gateCommitEvent, size: 10})
	r := commit.reservation
	if !assert.NotNil(t, r) {
		return
	}
	assert.Equal(t, uint64(10), r.packet.Size)
	assert.Equal(t, now, r.packet.SendTime)
	assert.Equal(t, uint64(10), lim.inflight())
	assert.Equal(t, uint64(1), lim.packetsInflight)

	ack := applyGateEvent(lim, now.Add(time.Millisecond), gateEvent{
		kind:        gateCompleteEvent,
		reservation: r,
		outcome:     flowcontrol.MessageOutcomeSucceeded,
	})
	assert.Equal(t, gateActionAck, ack.action)
	assert.Zero(t, lim.inflight())
	assert.Zero(t, lim.packetsInflight)
	assert.Equal(t, uint64(10), lim.delivered)
}

func TestGateAbortReversesOnlyOnWireAccounting(t *testing.T) {
	now := time.Unix(1, 0)
	lim := newLimiterState(clock.NewManual(now), limiterOptions{})
	a := applyGateEvent(lim, now, gateEvent{kind: gateCommitEvent, size: 10}).reservation
	b := applyGateEvent(lim, now, gateEvent{kind: gateCommitEvent, size: 20}).reservation
	nextSendTime := lim.nextSendTime
	pacingReady := lim.pacingReady
	delivered := lim.delivered
	stateType := fmt.Sprintf("%T", lim.state)

	abort := applyGateEvent(lim, now, gateEvent{
		kind:        gateCompleteEvent,
		reservation: a,
		outcome:     flowcontrol.MessageOutcomeAbortedBeforeEnqueue,
	})
	assert.Equal(t, gateActionUnsend, abort.action)
	assert.Equal(t, uint64(20), lim.sent)
	assert.Equal(t, uint64(1), lim.packetsInflight)
	assert.Equal(t, uint64(20), lim.inflight())
	assert.Equal(t, nextSendTime, lim.nextSendTime)
	assert.Equal(t, pacingReady, lim.pacingReady)
	assert.Equal(t, delivered, lim.delivered)
	assert.Equal(t, stateType, fmt.Sprintf("%T", lim.state))
	assert.Same(t, b, lim.gate.reservations[0])
}

func TestGateAbortOfLaterReservationDoesNotBreakEarlierAck(t *testing.T) {
	now := time.Unix(1, 0)
	lim := newLimiterState(clock.NewManual(now), limiterOptions{})
	a := applyGateEvent(lim, now, gateEvent{kind: gateCommitEvent, size: 10}).reservation
	b := applyGateEvent(lim, now, gateEvent{kind: gateCommitEvent, size: 20}).reservation

	abort := applyGateEvent(lim, now, gateEvent{
		kind:        gateCompleteEvent,
		reservation: b,
		outcome:     flowcontrol.MessageOutcomeAbortedBeforeEnqueue,
	})
	assert.Equal(t, gateActionUnsend, abort.action)
	ack := applyGateEvent(lim, now.Add(time.Millisecond), gateEvent{
		kind:        gateCompleteEvent,
		reservation: a,
		outcome:     flowcontrol.MessageOutcomeSucceeded,
	})
	assert.Equal(t, gateActionAck, ack.action)
	assert.Zero(t, lim.inflight())
	assert.Zero(t, lim.packetsInflight)
	assert.Equal(t, uint64(10), lim.delivered)
}

func TestGateFatalDoesNotAcknowledgeOrUnsendBBR(t *testing.T) {
	now := time.Unix(1, 0)
	lim := newLimiterState(clock.NewManual(now), limiterOptions{})
	r := applyGateEvent(lim, now, gateEvent{kind: gateCommitEvent, size: 10}).reservation
	action := applyGateEvent(lim, now, gateEvent{
		kind:        gateCompleteEvent,
		reservation: r,
		outcome:     flowcontrol.MessageOutcomeFatal,
		err:         errors.New("fatal"),
	}).action
	assert.Equal(t, gateActionNone, action)
	assert.Equal(t, uint64(10), lim.inflight())
	assert.Equal(t, uint64(1), lim.packetsInflight)
	assert.Zero(t, lim.delivered)
}

func applyGateEvent(lim *Limiter, now time.Time, event gateEvent) gateEventResult {
	event.reply = make(chan gateEventResult, 1)
	lim.handleGateEvent(now, event)
	return <-event.reply
}
