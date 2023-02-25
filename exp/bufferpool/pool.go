// Package bufferpool supports object pooling for byte buffers.
package bufferpool

import "sync"

const (
	defaultMinSize     = 1024
	defaultBucketCount = 10
)

var pool Pool = &BasicPool{}

// Set the global pool.  This can be used to override the default implementation
// to optimize for different allocation profiles. For example, applications that
// target embedded platforms may wish to lower the minimum buffer size from 1KiB
// to 512B, or less.
//
// The default implementation uses BasicPool, which is suitable for a wide range
// of applications.  Most applications will not benefit from using a custom pool,
// so be sure to profile.
func Set(p Pool) {
	pool = p
}

// Get a buffer of with the specified size from the global pool.  The
// resulting buffer's capacity MAY exceed size. The buffer SHOULD NOT
// be resized beyond its underlying capacity, or it will leak to the
// garbage collector.
func Get(size int) []byte {
	return pool.Get(size)[:size]
}

// Put returns the buffer to the global pool.  The slice will be zeroed
// before returning it to the pool.   Note that only the first len(buf)
// bytes are zeroed, and not the full backing array.
//
// TODO:  should we zero out the entire array?  Users may have resized
// the slice, so it may contain sensitive data.
func Put(buf []byte) {
	for i := range buf {
		buf[i] = 0
	}

	pool.Put(buf[:cap(buf)])
}

// Pool is a free list of []byte buffers of heterogeneous sizes.
// Implementations MUST be thread safe and MUST support efficient
// paging of buffers with arbitrary capacity.   Most applications
// SHOULD use BasicPool.
type Pool interface {
	// Get a buffer of len >= size.  The caller will resize the
	// buffer to size, so implementations need not bother.
	Get(size int) []byte

	// Put returns the buffer to the pool.   The slice will be
	// zeroed and resized to its underlying capacity before it
	// is passed to Put.
	Put([]byte)
}

// BasicPool is a general-purpose implementation of Pool, using
// sync.Pool. BasicPool maintains a set of N buckets containing
// buffers of exponentially-increasing size, starting from Min.
// Buffers whose capacity exceeds Min << N are ignored and left
// to the garbage-collector.
//
// The zero-value BasicPool is ready to use, defaulting to N=10
// and Min=1024 (1KiB - ~1MiB buffer sizes).  Most applications
// will not benefit from tuning these parameters.
//
// As a general rule, increasing Min reduces GC latency at the
// expense of increased memory usage.  Increasing N can reduce
// GC latency in applications that frequently allocate buffers
// of size >=1MiB.
type BasicPool struct {
	once    sync.Once
	Min, N  int
	buckets []bucket
}

func (p *BasicPool) Get(size int) []byte {
	p.init()

	for _, b := range p.buckets {
		if b.Fits(size) {
			return b.Get()
		}
	}

	return make([]byte, size)
}

func (p *BasicPool) Put(buf []byte) {
	p.init()

	for _, b := range p.buckets {
		if b.TryPut(buf) {
			return
		}
	}
}

func (p *BasicPool) init() {
	p.once.Do(func() {
		if p.Min <= 0 {
			p.Min = defaultMinSize
		}

		if p.N <= 0 {
			p.N = defaultBucketCount
		}

		p.buckets = make([]bucket, p.N)
		for i := range p.buckets {
			p.buckets[i] = newBucket(p.Min, i)
		}
	})
}

type bucket struct {
	minCap, capLimit int
	pool             *sync.Pool
}

func newBucket(minSize, i int) bucket {
	return bucket{
		minCap:   minSize << i,
		capLimit: minSize << (i + 1),
		pool: &sync.Pool{New: func() any {
			return make([]byte, minSize<<i)
		}},
	}
}

// Returns true if all buffers in the bucket are guaranteed to
// have a capacity >= size.
func (b bucket) Fits(size int) bool {
	return b.minCap >= size
}

// Get a buffer.  The buffer's capacity is guaranteed to be at
// least b.minCap, and smaller than b.capLimit.
func (b bucket) Get() []byte {
	return b.pool.Get().([]byte)
}

// TryPut returns the buffer to the underlying sync.Pool if the
// buffer's capacity satisfies the bucket's constraints.
//
// Returns true if the buffer was put back in the pool.
func (b bucket) TryPut(buf []byte) (ok bool) {
	if ok = b.minCap <= cap(buf) && b.capLimit > cap(buf); ok {
		b.pool.Put(buf)
	}

	return
}
