package bbr

import (
	"context"
	"errors"

	"capnproto.org/go/capnp/v3/flowcontrol"
)

// gateState is the actor-owned ledger for the GateNext adapter. It is kept
// separate from the channel adapter so deterministic tests can exercise every
// reservation transition without scheduling a goroutine. The Limiter actor is
// its sole production owner.
//
// A reservation is provisional until it is acknowledged. A definitely-unsent
// operation is removed from the active log and reports replay to the actor;
// surviving successors remain valid provisional entries. Fatal completion
// preserves the first poison error and makes the controller terminal.
type gateState struct {
	nextID       uint64
	reservations []*gateReservation
	waiters      []*gateWaitAttempt
	poison       error
}

type gateReservationState uint8

const (
	gateReservationProvisional gateReservationState = iota
	gateReservationAcknowledged
	gateReservationAborted
	gateReservationPoisoned
)

type gateReservation struct {
	id    uint64
	size  uint64
	state gateReservationState
}

type gateEventKind uint8

const (
	gateCommitEvent gateEventKind = iota
	gateCompleteEvent
	gatePoisonEvent
	gateWaitEvent
	gateCancelWaitEvent
	gateGrantWaitEvent
)

// gateEvent is private actor traffic. Its reply channel is buffered so the
// actor never waits for a caller that has observed limiter shutdown.
type gateEvent struct {
	kind        gateEventKind
	size        uint64
	reservation *gateReservation
	outcome     flowcontrol.MessageOutcomeKind
	err         error
	reply       chan gateEventResult
	waiter      *gateWaitAttempt
}

type gateEventResult struct {
	reservation *gateReservation
	replay      bool
	err         error
}

type gateWaitState uint8

const (
	gateWaitWaiting gateWaitState = iota
	gateWaitGranted
	gateWaitCanceled
	gateWaitFailed
)

// gateWaitAttempt represents one invocation of a retryable successor
// permission. The actor alone changes state or sends result.
type gateWaitAttempt struct {
	reservation *gateReservation
	state       gateWaitState
	result      chan error
}

func newGateWaitAttempt(r *gateReservation) *gateWaitAttempt {
	return &gateWaitAttempt{reservation: r, result: make(chan error, 1)}
}

func (l *Limiter) gateCommit(size uint64) (*gateReservation, error) {
	result, err := l.gateEvent(gateEvent{kind: gateCommitEvent, size: size})
	if err == nil {
		err = result.err
	}
	return result.reservation, err
}

func (l *Limiter) gateComplete(r *gateReservation, kind flowcontrol.MessageOutcomeKind, err error) (bool, error) {
	result, eventErr := l.gateEvent(gateEvent{
		kind:        gateCompleteEvent,
		reservation: r,
		outcome:     kind,
		err:         err,
	})
	return result.replay, eventErr
}

func (l *Limiter) gatePoison(err error) error {
	_, eventErr := l.gateEvent(gateEvent{kind: gatePoisonEvent, err: err})
	return eventErr
}

// gateWait waits for one successor permission. Cancellation is resolved by
// the actor: a canceled attempt is removed without consuming permission, while
// a grant that won the race is returned as success.
func (l *Limiter) gateWait(ctx context.Context, r *gateReservation) error {
	a, err := l.gateStartWait(ctx, r)
	if err != nil {
		return err
	}
	return l.gateWaitResult(ctx, a)
}

func (l *Limiter) gateStartWait(ctx context.Context, r *gateReservation) (*gateWaitAttempt, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	a := newGateWaitAttempt(r)
	event := gateEvent{kind: gateWaitEvent, waiter: a, reply: make(chan gateEventResult, 1)}
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-l.ctx.Done():
		return nil, l.ctx.Err()
	case l.chGate <- event:
	}
	select {
	case result := <-event.reply:
		if result.err != nil {
			return nil, result.err
		}
	case <-ctx.Done():
		return nil, l.gateCancelWait(a, ctx.Err())
	case <-l.ctx.Done():
		return nil, l.ctx.Err()
	}
	return a, nil
}

func (l *Limiter) gateWaitResult(ctx context.Context, a *gateWaitAttempt) error {
	select {
	case err := <-a.result:
		return err
	case <-ctx.Done():
		return l.gateCancelWait(a, ctx.Err())
	case <-l.ctx.Done():
		return l.ctx.Err()
	}
}

func (l *Limiter) gateCancelWait(a *gateWaitAttempt, cancelErr error) error {
	result, err := l.gateEvent(gateEvent{
		kind:   gateCancelWaitEvent,
		waiter: a,
		err:    cancelErr,
	})
	if err != nil {
		return err
	}
	return result.err
}

// gateGrant is a private deterministic test seam. Production admission will
// call the same actor transition once BBR's pacing/cwnd state permits it.
func (l *Limiter) gateGrant(r *gateReservation) (bool, error) {
	result, err := l.gateEvent(gateEvent{kind: gateGrantWaitEvent, reservation: r})
	return result.replay, err
}

func (l *Limiter) gateEvent(event gateEvent) (gateEventResult, error) {
	event.reply = make(chan gateEventResult, 1)
	select {
	case <-l.ctx.Done():
		return gateEventResult{}, l.ctx.Err()
	case l.chGate <- event:
	}
	select {
	case <-l.ctx.Done():
		return gateEventResult{}, l.ctx.Err()
	case result := <-event.reply:
		return result, nil
	}
}

func (l *Limiter) handleGateEvent(event gateEvent) {
	var result gateEventResult
	switch event.kind {
	case gateCommitEvent:
		if l.gate.poison == nil {
			result.reservation = l.gate.commit(event.size)
		} else {
			result.err = l.gate.poison
		}
	case gateCompleteEvent:
		result.replay = l.gate.complete(event.reservation, event.outcome, event.err)
		l.gate.failWaiters()
	case gatePoisonEvent:
		l.gate.poisonWith(event.err)
		l.gate.failWaiters()
	case gateWaitEvent:
		result.err = l.gate.registerWait(event.waiter)
	case gateCancelWaitEvent:
		result.err = l.gate.cancelWait(event.waiter, event.err)
	case gateGrantWaitEvent:
		result.replay = l.gate.grantWait(event.reservation)
	default:
		l.gate.poisonWith(errors.New("bbr: invalid gate-next event"))
		l.gate.failWaiters()
	}
	event.reply <- result
}

func (g *gateState) commit(size uint64) *gateReservation {
	g.nextID++
	r := &gateReservation{id: g.nextID, size: size}
	g.reservations = append(g.reservations, r)
	return r
}

// complete records one terminal outcome. It returns true only for a
// definitely-unsent reservation, which tells the actor it must replay that
// operation before admitting another successor. Other provisional reservations
// remain in the ledger and must not be committed again.
func (g *gateState) complete(r *gateReservation, kind flowcontrol.MessageOutcomeKind, err error) (replay bool) {
	if g.poison != nil || r == nil || r.state != gateReservationProvisional {
		return false
	}
	switch kind {
	case flowcontrol.MessageOutcomeSucceeded:
		r.state = gateReservationAcknowledged
		g.compactAcknowledgedPrefix()
	case flowcontrol.MessageOutcomeAbortedBeforeEnqueue:
		r.state = gateReservationAborted
		g.remove(r)
		return true
	case flowcontrol.MessageOutcomeFatal, flowcontrol.MessageOutcomeFailedAfterEnqueue:
		r.state = gateReservationPoisoned
		g.poisonWith(err)
	default:
		r.state = gateReservationPoisoned
		g.poisonWith(errors.New("bbr: invalid gate-next terminal outcome"))
	}
	return false
}

func (g *gateState) poisonWith(err error) {
	if g.poison != nil {
		return
	}
	if err == nil {
		err = errors.New("bbr: gate-next poisoned")
	}
	g.poison = err
}

func (g *gateState) registerWait(a *gateWaitAttempt) error {
	if a == nil || a.reservation == nil {
		return errors.New("bbr: invalid gate-next wait")
	}
	if g.poison != nil {
		return g.poison
	}
	if a.reservation.state != gateReservationProvisional && a.reservation.state != gateReservationAborted {
		return errors.New("bbr: gate-next permission is no longer available")
	}
	for _, waiter := range g.waiters {
		if waiter.reservation == a.reservation && waiter.state == gateWaitWaiting {
			return errors.New("bbr: gate-next permission already has a waiter")
		}
	}
	g.waiters = append(g.waiters, a)
	return nil
}

func (g *gateState) cancelWait(a *gateWaitAttempt, err error) error {
	if a == nil {
		return errors.New("bbr: invalid gate-next wait")
	}
	switch a.state {
	case gateWaitWaiting:
		g.removeWait(a)
		a.state = gateWaitCanceled
		if err == nil {
			err = context.Canceled
		}
		a.result <- err
		return err
	case gateWaitGranted:
		return nil
	case gateWaitCanceled:
		return err
	case gateWaitFailed:
		if g.poison != nil {
			return g.poison
		}
		return errors.New("bbr: gate-next wait failed")
	default:
		panic("bbr: invalid gate-next wait state")
	}
}

// grantWait grants at most one retryable successor permission. The bool result
// means a waiter was granted; it does not request reservation replay.
func (g *gateState) grantWait(r *gateReservation) bool {
	if g.poison != nil || r == nil {
		return false
	}
	for _, a := range g.waiters {
		if a.reservation != r || a.state != gateWaitWaiting {
			continue
		}
		g.removeWait(a)
		a.state = gateWaitGranted
		a.result <- nil
		return true
	}
	return false
}

func (g *gateState) failWaiters() {
	if g.poison == nil {
		return
	}
	for len(g.waiters) > 0 {
		a := g.waiters[0]
		g.removeWait(a)
		if a.state != gateWaitWaiting {
			continue
		}
		a.state = gateWaitFailed
		a.result <- g.poison
	}
}

func (g *gateState) removeWait(a *gateWaitAttempt) {
	for i, waiter := range g.waiters {
		if waiter != a {
			continue
		}
		copy(g.waiters[i:], g.waiters[i+1:])
		g.waiters[len(g.waiters)-1] = nil
		g.waiters = g.waiters[:len(g.waiters)-1]
		return
	}
}

func (g *gateState) remove(r *gateReservation) {
	for i, candidate := range g.reservations {
		if candidate != r {
			continue
		}
		copy(g.reservations[i:], g.reservations[i+1:])
		g.reservations[len(g.reservations)-1] = nil
		g.reservations = g.reservations[:len(g.reservations)-1]
		return
	}
}

func (g *gateState) compactAcknowledgedPrefix() {
	for len(g.reservations) > 0 && g.reservations[0].state == gateReservationAcknowledged {
		g.reservations[0] = nil
		g.reservations = g.reservations[1:]
	}
}
