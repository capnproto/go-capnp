package bbr

import (
	"context"
	"testing"
	"time"

	"capnproto.org/go/capnp/v3/exp/clock"

	"github.com/stretchr/testify/assert"
)

func TestReleaseLimiter(t *testing.T) {
	lim := NewLimiter(clock.System)
	lim.Release()
	select {
	case <-lim.ctx.Done():
	case <-time.NewTimer(10 * time.Millisecond).C:
		t.Fatal("lim.Release did not cancel the context.")
	}

	_, err := lim.StartMessage(context.Background(), 1)
	assert.NotNil(t, err, "Error should be non-nil if the limiter has shut down.")
	assert.Equal(t, lim.ctx.Err(), err, "Error should be that of the limiter's context.")
}

func TestLimiterOneShot(t *testing.T) {
	lim := NewLimiter(clock.System)
	defer lim.Release()

	lim.whilePaused(func() {
		assert.Equal(t, lim.inflight(), uint64(0),
			"Limiter should start off with nothing in-flight.")
	})

	got1, err := lim.StartMessage(context.Background(), 7)
	assert.Nil(t, err, "StartMessage() failed.")

	lim.whilePaused(func() {
		assert.Equal(t, lim.inflight(), uint64(7),
			"Once we send the message, limiter should have that much data in-flight.")
	})

	got1()

	lim.whilePaused(func() {
		assert.Equal(t, lim.inflight(), uint64(0),
			"Once we receive the ack, in-flight data should be zero again.")
	})
}
