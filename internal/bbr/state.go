package bbr

// This file contains implementations of the various "states" the limiter
// can be in. This mostly corresponds to the paper's appendix.

import (
	"math"
	"math/rand"
	"time"
)

type state interface {
	initialize(lim *Limiter)
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

	lastRtPropChange time.Time
	rtProp           time.Duration
}

func (s *probeBWState) initialize(lim *Limiter) {
	lim.cwndGain = 2

	now := lim.clock.Now()
	s.rtProp = lim.rtPropFilter.Estimate
	s.lastRtPropChange = now

	// Randomly select an initial pacing gain; anything but the value
	// below 1 will do (see paper).
	s.pacingGainIndex = rand.Int() % (len(probeBWPacingGains) - 1)
	if s.pacingGainIndex == 1 {
		// Don't start with the 3/4.
		s.pacingGainIndex++
	}
	lim.pacingGain = probeBWPacingGains[s.pacingGainIndex]
	s.nextPacingGainChange = now.Add(s.rtProp)
}

func (s *probeBWState) postAck(lim *Limiter, p packetMeta, now time.Time) {
	rtProp := lim.rtPropFilter.Estimate
	if rtProp < s.rtProp {
		s.rtProp = rtProp
		s.lastRtPropChange = now
	}

	if now.Sub(s.lastRtPropChange) > 10*time.Second {
		// Been a while since we've measured rtProp; switch to probeRTT.
		lim.changeState(&probeRTTState{})
		return
	}

	if now.After(s.nextPacingGainChange) {
		s.pacingGainIndex++
		s.pacingGainIndex %= len(probeBWPacingGains)
		lim.pacingGain = probeBWPacingGains[s.pacingGainIndex]
		s.nextPacingGainChange = now.Add(rtProp)
	}
}

type probeRTTState struct {
	// The time at which we can exit this state.
	exitTime time.Time

	// The value of Limiter.sent when we entered this state.
	// Used to work out when we've seen a full round trip.
	initSent int64
}

func (s *probeRTTState) initialize(lim *Limiter) {
	// TODO: The paper says to send cwnd to "4 packets;" I don't know
	// exactly how to translate this to our setting...
	lim.cwndGain = 1
	// TODO: pacingGain?
	now := lim.clock.Now()

	s.exitTime = now.Add(200 * time.Millisecond)
	s.initSent = lim.sent
}

func (s *probeRTTState) postAck(lim *Limiter, p packetMeta, now time.Time) {
	// We can tell if we've seen a full round trip if the amount of data
	// *delivered* is greater than the amount that had been *sent* when
	// we started:
	afterRoundTrip := lim.delivered > s.initSent
	if afterRoundTrip && now.After(s.exitTime) {
		// TODO: figure out whether to switch into startup or probeBW.
	}
}
