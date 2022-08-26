package bbr

import (
	"math"
	"math/rand"
	"time"
)

type state interface {
	initialize(lim *Limiter)
	preAck(lim *Limiter, pm packetMeta, now time.Time)
	postAck(lim *Limiter, pm packetMeta, now time.Time)
}

type startupState struct {
	prevBtlBwEstimate bytesPerNs
	plateuRounds      int
}

func (s *startupState) initialize(lim *Limiter) {
	lim.cwndGain = 2 / math.Ln2
	lim.pacingGain = 2 / math.Ln2
}

func (s *startupState) preAck(lim *Limiter, p packetMeta, now time.Time) {
}

func (s *startupState) postAck(lim *Limiter, p packetMeta, now time.Time) {
	newBtlBwEstimate := lim.btlBwFilter.Estimate
	if float64(newBtlBwEstimate) < 1.25*float64(s.prevBtlBwEstimate) {
		s.plateuRounds++
	} else {
		s.plateuRounds = 0
	}
	s.prevBtlBwEstimate = newBtlBwEstimate
	if s.plateuRounds >= 3 {
		lim.changeState(&drainState{})
	}
}

type drainState struct {
}

func (s *drainState) initialize(lim *Limiter) {
	lim.pacingGain = 1 / lim.pacingGain
}

func (s *drainState) preAck(lim *Limiter, p packetMeta, now time.Time) {
}

func (s *drainState) postAck(lim *Limiter, p packetMeta, now time.Time) {
	// XXX, do we actually want to ensure *equality*? will that always
	// converge? Appears to work in tests, but need to think about it
	// and make sure (our tests include messages of size 1; probably
	// we should test with values that won't divide evenly, but also
	// we should reason this out -- I suspect it will not always hold
	//
	// Maybe instead we should just check that it's within 5% or so
	// (and we haven't overshot it).
	if lim.inflight() == int64(lim.computeBDP()) {
		lim.changeState(&probeBWState{})
	}
}

var probeBWPacingGains = []float64{
	1.25,
	0.75,
	1,
	1,
	1,
	1,
	1,
	1,
}

type probeBWState struct {
	// current index into probeBWPacingGains
	pacingGainIndex int

	// Time at which we should rotate to a new pacing gain
	// (last change + rtProp):
	nextPacingGainChange time.Time
}

func (s *probeBWState) initialize(lim *Limiter) {
	lim.cwndGain = 2

	// Randomly select an initial pacing gain; anything but the value
	// below 1 will do (see paper).
	s.pacingGainIndex = rand.Int() % (len(probeBWPacingGains) - 1)
	if s.pacingGainIndex == 1 {
		// Don't start with the 3/4.
		s.pacingGainIndex++
	}
	lim.pacingGain = probeBWPacingGains[s.pacingGainIndex]
	s.nextPacingGainChange = lim.clock.Now().Add(lim.rtPropFilter.Estimate)
}

func (s *probeBWState) preAck(lim *Limiter, p packetMeta, now time.Time) {
}

func (s *probeBWState) postAck(lim *Limiter, p packetMeta, now time.Time) {
	if now.After(s.nextPacingGainChange) {
		s.pacingGainIndex++
		s.pacingGainIndex %= len(probeBWPacingGains)
		lim.pacingGain = probeBWPacingGains[s.pacingGainIndex]
		s.nextPacingGainChange = now.Add(lim.rtPropFilter.Estimate)
	}
	// TODO: check if we should enter probeRTT.
}
