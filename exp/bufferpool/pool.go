// Package bufferpool supports object pooling for byte buffers.
package bufferpool

import "sync"

const (
	nBuckets = 20
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
	for i := 0; i < nBuckets; i++ {
		capacity := 1 << i
		if capacity >= size {
			var ret []byte
			item := p.buckets[i].Get()
			if item == nil {
				ret = make([]byte, capacity)
			} else {
				ret = item.([]byte)
			}
			ret = ret[:size]
			return ret
		}
	}
	return make([]byte, size)
}

// Return a buffer to the pool. Zeros the slice (but not the full backing array)
// before making it available for future use.
func (p *Pool) Put(buf []byte) {
	for i := 0; i < len(buf); i++ {
		buf[i] = 0
	}

	capacity := cap(buf)
	for i := 0; i < nBuckets; i++ {
		if (1 << i) == capacity {
			p.buckets[i].Put(buf[:capacity])
			return
		}
	}
}
