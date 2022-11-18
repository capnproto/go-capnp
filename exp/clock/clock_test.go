package clock

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimers(t *testing.T) {
	// arbitrary initial now:
	init := time.Unix(1e9, 0)

	clock := NewManual(init)

	timer := clock.NewTimer(5 * time.Second)

	select {
	case <-timer.Chan():
		t.Fatal("Timer went off too early.")
	default:
	}

	clock.Advance(1 * time.Second)

	select {
	case <-timer.Chan():
		t.Fatal("Timer went off too early.")
	default:
	}

	clock.Advance(10 * time.Second)

	select {
	case <-timer.Chan():
	default:
		t.Fatal("Timer didn't go off.")
	}

	assert.False(t, timer.Stop(), "Timer should have already been stopped.")

	assert.Equal(t, init.Add(11*time.Second), clock.Now(), "Time should be 11 seconds later.")

	timer.Reset(1 * time.Second)

	select {
	case <-timer.Chan():
		t.Fatal("Timer went off again too early.")
	default:
	}

	clock.Advance(2 * time.Second)

	select {
	case <-timer.Chan():
	default:
		t.Fatal("Timer didn't go off.")
	}
}

func TestImmediateTimer(t *testing.T) {
	clock := NewManual(time.Unix(1e9, 0))

	timer := clock.NewTimer(0)
	select {
	case <-timer.Chan():
	default:
		t.Fatal("Timer with duration of zero should go off immediately.")
	}

	timer = clock.NewTimer(-1)
	select {
	case <-timer.Chan():
	default:
		t.Fatal("Timer with negative duration should go off immediately.")
	}
}
