package bbr

import (
	"math"
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
	// TODO: check if lim.inflight == estimated BDP, and if so switch to probeBW.
	// Hm, do we actually want to ensure *equality*? will that converge. Need
	// to make sure.
}

/*
const (
	probeBWState stateName = iota
	probeRTTState
	startupState
	drainState
)
*/
