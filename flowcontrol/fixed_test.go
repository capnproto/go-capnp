package flowcontrol

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFixed(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	lim := NewFixedLimiter(10)

	// Start a couple messages:
	got4, err := lim.StartMessage(ctx, 4)
	require.NoError(t, err, "Limiter returned an error")
	got6, err := lim.StartMessage(ctx, 6)
	require.NoError(t, err, "Limiter returned an error")

	// We're now exactly at the limit, so if we try again it should block:
	func() {
		ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
		defer cancel()

		_, err = lim.StartMessage(ctxTimeout, 1)
		assert.ErrorIs(t, err, context.DeadlineExceeded, "should return context error")
	}()

	// Ok, finish one of them and then it should go through again:
	got4()
	got1, err := lim.StartMessage(ctx, 1)
	require.NoError(t, err, "Limiter returned an error")

	// There are 10 - (6 + 1) = 3 bytes remaining. It should therefore block
	// if we ask for four:
	func() {
		ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
		defer cancel()

		_, err = lim.StartMessage(ctxTimeout, 4)
		assert.ErrorIs(t, err, context.DeadlineExceeded, "should return context error")
	}()
	got6()
	got1()
}

func TestFixeLimiterPanics(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		NewFixedLimiter(1024).StartMessage(context.Background(), 1025)
	}, "should panic if reservation would cause deadlock")
}
