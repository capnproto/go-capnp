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
	if r == nil || r.state != gateReservationProvisional {
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
