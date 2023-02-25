// Package bufferpool supports object pooling for byte buffers.
package bufferpool

import "sync"

const (
	minSize  = 1024
	nBuckets = 11 // max 1Mb
)

// A default global pool.
var Default Pool

// A pool of buffers, in variable sizes.
type Pool struct {
	buckets [nBuckets]sync.Pool
}

// Get a buffer of the given length. Its capacity may be larger than the
// requested size.
func (p *Pool) Get(size int) []byte {
	for i := range p.buckets {
		if capacity := minSize << i; capacity >= size {
			if item := p.buckets[i].Get(); item != nil {
				return item.([]byte)[:size]
			}

			return make([]byte, size, capacity)
		}
	}

	return make([]byte, size)
}

// Return a buffer to the pool. Zeros the slice (but not the full backing array)
// before making it available for future use.
func (p *Pool) Put(buf []byte) {
	for i := range buf {
		buf[i] = 0
	}

	capacity := cap(buf)
	for i := range p.buckets {
		bucket := (minSize << i)
		next := (minSize<<i + 1)

		if bucket <= capacity && capacity < next {
			p.buckets[i].Put(buf[:capacity])
			return
		}
	}
}
