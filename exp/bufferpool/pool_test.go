package bufferpool_test

import (
	"testing"

	"capnproto.org/go/capnp/v3/exp/bufferpool"
	"github.com/stretchr/testify/assert"
)

func TestInvariants(t *testing.T) {
	t.Parallel()

	pool := bufferpool.Pool{
		BucketCount: 1,
		MinAlloc:    1024,
	}

	assert.Panics(t, func() {
		pool.Get(32)
	}, "should panic when BucketCount cannot satisfy MinAlloc.")
}

func TestGet(t *testing.T) {
	t.Parallel()
	t.Helper()

	pool := bufferpool.Pool{
		BucketCount: 8,
		MinAlloc:    64,
	}

	t.Run("Size<MinAlloc", func(t *testing.T) {
		assert.Len(t, pool.Get(32), 32, "should return buffer with len=32")
		assert.Equal(t, 32, cap(pool.Get(32)), "should return buffer with cap=MinAlloc")
	})

	t.Run("Size=MinAlloc", func(t *testing.T) {
		assert.Len(t, pool.Get(pool.MinAlloc), pool.MinAlloc, "should return buffer with len=MinAlloc")
		assert.Equal(t, pool.MinAlloc, cap(pool.Get(pool.MinAlloc)), "should return buffer with cap=MinAlloc")
	})

	t.Run("Size>MinAlloc", func(t *testing.T) {
		assert.Len(t, pool.Get(33), 33, "should return buffer with len=33")
		assert.Equal(t, 128, cap(pool.Get(128)), "should return buffer with cap=128")
	})

	t.Run("Size>MaxSize", func(t *testing.T) {
		assert.Len(t, pool.Get(512), 512, "should return buffer with len=512")
		assert.GreaterOrEqual(t, 512, cap(pool.Get(512)), "should return buffer with cap>=512")
	})
}

func TestPut(t *testing.T) {
	t.Parallel()

	var pool bufferpool.Pool

	buf := make([]byte, 1024)
	for i := range buf {
		buf[i] = byte(i)
	}
	buf = buf[:8]

	pool.Put(buf)
	buf = pool.Get(8)

	assert.Equal(t, make([]byte, 8), buf, "should zero first 8 bytes")
}

func BenchmarkPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		buf := bufferpool.Default.Get(32)
		bufferpool.Default.Put(buf)
	}
}
