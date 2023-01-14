package flowcontrol

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFixed(t *testing.T) {
	ctx := context.Background()
	lim := NewFixedLimiter(10)

	// Start a couple messages:
	got4, err := lim.StartMessage(ctx, 4)
	assert.Nil(t, err, "Limiter returned an error")
	got6, err := lim.StartMessage(ctx, 6)
	assert.Nil(t, err, "Limiter returned an error")

	// We're now exactly at the limit, so if we try again it should block:
	func() {
		ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
		defer cancel()

		_, err = lim.StartMessage(ctxTimeout, 1)
		assert.NotNil(t, err, "Limiter didn't return an error")
		assert.Equal(t, err, ctxTimeout.Err(), "Error wasn't from the context")
	}()

	// Ok, finish one of them and then it should go through again:
	got4()
	got1, err := lim.StartMessage(ctx, 1)
	assert.Nil(t, err, "Limiter returned an error")

	// There are 10 - (6 + 1) = 3 bytes remaining. It should therefore block
	// if we ask for four:
	func() {
		ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
		defer cancel()

		_, err = lim.StartMessage(ctxTimeout, 4)
		assert.NotNil(t, err, "Limiter didn't return an error")
		assert.Equal(t, err, ctxTimeout.Err(), "Error wasn't from the context")
	}()
	got6()
	got1()

}

func TestInflight(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	lim := NewInflightMessageLimiter(2)

	// Start a couple messages:
	got1, err := lim.StartMessage(ctx, 4) // use arbitrary size to ensure we don't block
	require.NoError(t, err, "Limiter returned an error")
	got2, err := lim.StartMessage(ctx, 6)
	require.NoError(t, err, "Limiter returned an error")

	// We're now exactly at the limit, so if we try again it should block:
	func() {
		ctxTimeout, cancel := context.WithTimeout(ctx, 10*time.Millisecond)
		defer cancel()

		release, err := lim.StartMessage(ctxTimeout, 1)
		require.ErrorIs(t, err, context.DeadlineExceeded, "Limiter didn't return an error")
		assert.Nil(t, release, "should not return release function for failed call to StartMessage")
	}()

	// Ok, finish one of them and then it should go through again:
	got1()
	got1, err = lim.StartMessage(ctx, 1)
	assert.NoError(t, err, "Limiter returned an error")

	// Clean up
	got1()
	got2()
}
