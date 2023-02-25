package bufferpool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicBuffer(t *testing.T) {
	t.Parallel()

	assert.Len(t, Get(32), 32, "should return buffer with 32-byte length")
	assert.Equal(t, 1024, cap(Get(32)), "should return buffer with 1KiB capacity")

	assert.Len(t, Get(1025), 1025, "should return buffer with 1025-byte length")
	assert.Equal(t, 2048, cap(Get(1025)), "should return buffer with 2KiB-byte capacity")
}

func TestBucket(t *testing.T) {
	t.Parallel()
	t.Helper()

	t.Run("Fits", func(t *testing.T) {
		t.Parallel()

		// 0th bucket
		bucket := newBucket(32, 0)
		assert.True(t, bucket.Fits(32), "should fit 32-byte buffer")
		assert.True(t, bucket.Fits(8), "should fit buffer smaller than min. capacity")
		assert.False(t, bucket.Fits(33), "should not fit buffer larger than min. capacity")

		// 1st bucket
		bucket = newBucket(32, 1)
		assert.True(t, bucket.Fits(32<<1), "should fit 32-byte buffer")
		assert.True(t, bucket.Fits(8<<1), "should fit buffer smaller than min. capacity")
		assert.False(t, bucket.Fits(33<<1), "should not fit buffer larger than min. capacity")
	})

	t.Run("Get", func(t *testing.T) {
		t.Parallel()

		// 0th bucket
		bucket := newBucket(32, 0)
		assert.Len(t, bucket.Get(), 32, "should allocate new buffers with min. capacity")

		// 1st bucket
		bucket = newBucket(32, 1)
		assert.Len(t, bucket.Get(), 32<<1, "should allocate new buffers with min. capacity")
	})

	t.Run("TryPut", func(t *testing.T) {
		t.Parallel()

		// 0th bucket
		bucket := newBucket(32, 0)
		assert.True(t, bucket.TryPut(make([]byte, 32)), "should consume buffer with exactly min. capacity")
		assert.True(t, bucket.TryPut(make([]byte, 63)), "should consume buffer with cap < limit")
		assert.False(t, bucket.TryPut(make([]byte, 64)), "should not consume buffer with cap == limit")

		// 1st bucket
		bucket = newBucket(32, 1)
		assert.True(t, bucket.TryPut(make([]byte, 32<<1)), "should consume buffer with exactly min. capacity")
		assert.True(t, bucket.TryPut(make([]byte, 63<<1)), "should consume buffer with cap < limit")
		assert.False(t, bucket.TryPut(make([]byte, 64<<1)), "should not consume buffer with cap == limit")
	})
}
