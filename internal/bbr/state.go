package bbr

import (
	"math"
	"time"
)

type state interface {
	initialize(mgr *Manager)
	preAck(mgr *Manager, pm packetMeta, now time.Time)
	postAck(mgr *Manager, pm packetMeta, now time.Time)
}

type startupState struct {
	prevBtlBwEstimate int64
	plateuRounds      int
}

func (s *startupState) initialize(mgr *Manager) {
	mgr.cwndGain = 2 / math.Ln2
	mgr.pacingGain = 2 / math.Ln2
}

func (s *startupState) preAck(mgr *Manager, p packetMeta, now time.Time) {
}

func (s *startupState) postAck(mgr *Manager, p packetMeta, now time.Time) {
	newBtlBwEstimate := mgr.btlBwFilter.Estimate
	if s.prevBtlBwEstimate == 0 {
		// This is our first sample.
		s.prevBtlBwEstimate = newBtlBwEstimate
		return
	}
	if float64(newBtlBwEstimate) < 1.25*float64(s.prevBtlBwEstimate) {
		s.plateuRounds++
	} else {
		s.plateuRounds = 0
	}
	if s.plateuRounds >= 3 {
		mgr.changeState(&drainState{})
	}
}

type drainState struct {
}

func (s *drainState) initialize(mgr *Manager) {
	mgr.pacingGain = 1 / mgr.pacingGain
}

func (s *drainState) preAck(mgr *Manager, p packetMeta, now time.Time) {
}

func (s *drainState) postAck(mgr *Manager, p packetMeta, now time.Time) {
	// TODO: check if mgr.inflight == estimated BDP, and if so switch to probeBW.
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
