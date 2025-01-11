package bufferpool_test

import (
	"testing"

	"capnproto.org/go/capnp/v3/exp/bufferpool"
	"github.com/stretchr/testify/assert"
)

func TestInvariants(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		bufferpool.NewPool(1024, 1)
	}, "should panic when BucketCount cannot satisfy MinAlloc.")
}

func TestGet(t *testing.T) {
	t.Parallel()
	t.Helper()

	minAlloc, bucketCount := 64, 8
	pool := bufferpool.NewPool(minAlloc, bucketCount)

	t.Run("Size<MinAlloc", func(t *testing.T) {
		size := minAlloc / 2 // size < minAlloc
		assert.Len(t, pool.Get(size), size, "should return buffer with len=%d", size)
		assert.Equal(t, minAlloc, cap(pool.Get(size)), "should return buffer with cap=MinAlloc")
	})

	t.Run("Size=MinAlloc", func(t *testing.T) {
		assert.Len(t, pool.Get(minAlloc), minAlloc, "should return buffer with len=MinAlloc")
		assert.Equal(t, minAlloc, cap(pool.Get(minAlloc)), "should return buffer with cap=MinAlloc")
	})

	t.Run("Size>MinAlloc", func(t *testing.T) {
		size := minAlloc + 1
		nextAlloc := minAlloc << 1
		assert.Len(t, pool.Get(size), size, "should return buffer with len=%d", size)
		assert.Equal(t, nextAlloc, cap(pool.Get(size)), "should return buffer with cap=%d", nextAlloc)
	})

	t.Run("Size>MaxSize", func(t *testing.T) {
		size := 1<<(bucketCount+1) + 1
		assert.Len(t, pool.Get(size), size, "should return buffer with len=%d", size)
		assert.GreaterOrEqual(t, size, cap(pool.Get(size)), "should return buffer with cap>=%d", size)
	})
}

func TestPut(t *testing.T) {
	t.Parallel()

	pool := bufferpool.NewPool(1024, 20)

	buf := make([]byte, 1024, 1024)
	for i := range buf {
		buf[i] = byte(i)
	}
	buf = buf[:8]

	pool.Put(buf)
	buf = pool.Get(8)
	fullBuf := buf[:cap(buf)]

	assert.Len(t, fullBuf, 1024, "buf should have len and cap 1024")
	assert.Equal(t, make([]byte, 8), buf, "should zero first 8 bytes")
	assert.NotEqual(t, 0, fullBuf[9], "first byte after clearing should not be zero")
}

func BenchmarkPool(b *testing.B) {
	pool := bufferpool.NewPool(0, 0)
	const size = 32

	// Make cache hot.
	buf := pool.Get(size)
	pool.Put(buf)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf := pool.Get(size)
		pool.Put(buf)
	}
}
