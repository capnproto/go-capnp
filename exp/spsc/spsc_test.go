package spsc

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test that filling/overflowing the internal items queue doesn't block the sender.
func TestFillItemsNonBlock(t *testing.T) {
	t.Parallel()

	q := New[int]()

	for i := 0; i < itemsBuffer+1; i++ {
		q.Send(i)
	}
}

// Try filling the queue, then draining it with TryRecv().
func TestFillThenTryDrain(t *testing.T) {
	t.Parallel()

	q := New[int]()

	limit := itemsBuffer + 1

	for i := 0; i < limit; i++ {
		q.Send(i)
	}

	for i := 0; i < limit; i++ {
		v, ok := q.TryRecv()
		assert.True(t, ok)
		assert.Equal(t, i, v)
	}
	_, ok := q.TryRecv()
	assert.False(t, ok)
}

// Try filling the queue, then draining it with Recv().
func TestFillThenDrain(t *testing.T) {
	t.Parallel()

	q := New[int]()

	limit := itemsBuffer + 1

	for i := 0; i < limit; i++ {
		q.Send(i)
	}

	ctx := context.Background()
	for i := 0; i < limit; i++ {
		v, err := q.Recv(ctx)
		assert.Nil(t, err)
		assert.Equal(t, i, v)
	}
	ctx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel()
	_, err := q.Recv(ctx)
	assert.NotNil(t, err)
	assert.ErrorIs(t, err, ctx.Err())
}

// Test that we get io.EOF after closing the tx.
func TestClose(t *testing.T) {
	t.Parallel()

	// try it with various numbers of items in the queue:
	cases := []struct {
		size int
		desc string
	}{
		{0, "empty"},
		{1, "one element"},
		{itemsBuffer, "full buffer"},
		{itemsBuffer + 1, "overfull buffer"},
	}

	for _, c := range cases {
		t.Run("size: "+c.desc, func(t *testing.T) {
			q := New[int]()

			for i := 0; i < c.size; i++ {
				q.Send(i)
			}
			q.Close()

			ctx := context.Background()
			for i := 0; i < c.size; i++ {
				v, err := q.Recv(ctx)
				assert.NoError(t, err)
				assert.Equal(t, i, v)
			}

			_, err := q.Recv(ctx)
			assert.Equal(t, err, io.EOF)
		})
	}
}
