// Package clock provides an interface for clocks.
// This is useful in particular for testing code which checks
// the current time.
package clock

import (
	"time"
)

// A Clock is can return the current time.
type Clock interface {
	Now() time.Time
}

// System reads the current time from the system clock.
var System Clock = systemClock{}

type systemClock struct{}

func (systemClock) Now() time.Time {
	return time.Now()
}

// A Manual stores the current time, which does not change unless
// the CurrentTime field is reassigned manually.
type Manual struct {
	CurrentTime time.Time
}

// Now returns m.CurrentTime.
func (m *Manual) Now() time.Time {
	return m.CurrentTime
}

// Advance advances m.CurrentTime forward by the given duration.
func (m *Manual) Advance(d time.Duration) {
	m.CurrentTime = m.CurrentTime.Add(d)
}
