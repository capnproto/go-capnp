package bbr

import (
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
}

type gateEventResult struct {
	reservation *gateReservation
	replay      bool
	err         error
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
	case gatePoisonEvent:
		l.gate.poisonWith(event.err)
	default:
		l.gate.poisonWith(errors.New("bbr: invalid gate-next event"))
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
