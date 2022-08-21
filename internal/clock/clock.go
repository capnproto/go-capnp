// Package clock provides an interface for clocks.
//
// This is useful for testing code which checks does things based on
// time; it allows advancing the clock as needed by the tests, rather
// than needing to race with real time.
package clock

import (
	"math"
	"sync"
	"time"

	"capnproto.org/go/capnp/v3/internal/syncutil"
)

// A Clock can measure time.
type Clock interface {
	Now() time.Time
	NewTimer(d time.Duration) Timer
}

type Timer interface {
	Chan() <-chan time.Time
	Reset(d time.Duration)
	Stop() bool
}

// System reads the current time from the system clock.
var System Clock = systemClock{}

type systemClock struct{}

func (systemClock) Now() time.Time {
	return time.Now()
}

func (systemClock) NewTimer(d time.Duration) Timer {
	return (*systemTimer)(time.NewTimer(d))
}

type systemTimer time.Timer

func (t *systemTimer) Chan() <-chan time.Time {
	return t.C
}

func (t *systemTimer) Reset(d time.Duration) {
	(*time.Timer)(t).Reset(d)
}

func (t *systemTimer) Stop() bool {
	return (*time.Timer)(t).Stop()
}

// A Manual is a clock which is stopped, and only advances when its Advance
// method is called.
type Manual struct {
	mu     sync.Mutex
	now    time.Time
	timers []*manualTimer
}

// Returns a new Manual clock, with the given initial time.
func NewManual(now time.Time) *Manual {
	return &Manual{now: now}
}

// Now returns the current time.
func (m *Manual) Now() (now time.Time) {
	syncutil.With(&m.mu, func() {
		now = m.now
	})
	return now
}

func (m *Manual) NewTimer(d time.Duration) Timer {
	var ret *manualTimer
	syncutil.With(&m.mu, func() {
		ret = &manualTimer{
			ch:       make(chan time.Time, 1),
			deadline: m.now.Add(d),
			clock:    m,
		}
		m.timers = append(m.timers, ret)
		if d <= 0 {
			ret.ch <- m.now
		}
	})
	return ret
}

// Advance advances the clock forward by the given duration.
func (m *Manual) Advance(d time.Duration) {
	syncutil.With(&m.mu, func() {
		before := m.now
		m.now = before.Add(d)

		for i := range m.timers {
			t := m.timers[i]

			if before.Before(t.deadline) && !m.now.Before(t.deadline) {
				t.ch <- m.now
			}
		}
	})
}

type manualTimer struct {
	ch       chan time.Time
	deadline time.Time
	clock    *Manual
}

func (t *manualTimer) Chan() <-chan time.Time {
	return t.ch
}

func (t *manualTimer) Reset(d time.Duration) {
	syncutil.With(&t.clock.mu, func() {
		t.deadline = t.clock.now.Add(d)
	})
}

func (t *manualTimer) Stop() bool {
	var wasActive bool
	syncutil.With(&t.clock.mu, func() {
		wasActive = t.clock.now.Before(t.deadline)
		t.deadline = time.Unix(math.MinInt64, 0)
	})
	return wasActive
}
