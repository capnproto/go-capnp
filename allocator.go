package capnp

import (
	"errors"

	"capnproto.org/go/capnp/v3/exp/bufferpool"
)

// allocator defines methods for an allocator for the basic arena
// implementation: the allocator defines the strategy of how to grow and
// release memory slices that back segments.
type allocator interface {
	// Grow an array. When called, it has already been determined that
	// the slice must be grown. Implementations must copy the contents of
	// the existing byte slice to the new one.
	//
	// The returned slice may be have a capacity larger than strcictly
	// required to store an addition sz bytes. In this case, the length
	// must be len(b)+sz, while cap may be any amount >= than that.
	Grow(b []byte, totalMsgSize int64, sz Size) ([]byte, error)

	// Release signals that the passed byte slice won't be used by the
	// arena anymore.
	Release(b []byte)
}

// bufferPoolAllocator allocates buffers from a buffer pool. If nil, then
// the global default buffer pool is used.
type bufferPoolAllocator bufferpool.Pool

func (bpa *bufferPoolAllocator) Grow(b []byte, totalMsgSize int64, sz Size) ([]byte, error) {
	pool := (*bufferpool.Pool)(bpa)
	if pool == nil {
		pool = &bufferpool.Default
	}

	inc, err := nextAlloc(totalMsgSize, int64(maxAllocSize()), sz)
	if err != nil {
		return nil, err
	}
	nb := pool.Get(len(b) + int(inc))

	// When there was a prior buffer, copy the contents to the new buffer,
	// clear the old buffer and return the old buffer to the pool to be
	// reused.
	if b != nil {
		copy(nb[:cap(nb)], b[:cap(b)])
		for i := range b {
			b[i] = 0
		}
		pool.Put(b)
	}

	return nb, nil
}

func (bpa *bufferPoolAllocator) Release(b []byte) {
	if b == nil {
		panic("nil buffer passed to release")
	}
	pool := (*bufferpool.Pool)(bpa)
	if pool == nil {
		pool = &bufferpool.Default
	}
	pool.Put(b)
}

// simpleAllocator allocates buffers without any caching or reuse, using the
// standard memory management functions.
//
// This allocator is concurent safe for use across multiple arenas.
type simpleAllocator struct{}

func (_ simpleAllocator) Grow(b []byte, totalMsgSize int64, sz Size) ([]byte, error) {
	inc, err := nextAlloc(totalMsgSize, int64(maxAllocSize()), sz)
	if err != nil {
		return nil, err
	}
	return append(b, make([]byte, inc)...), nil
}

func (_ simpleAllocator) Release(b []byte) {
	// Nothing to do. The runtime GC will reclaim it.
}

// segmentList defines the operations needed for a container of segments.
type segmentList interface {
	// NumSegments must return the number of segments that exist.
	NumSegments() int

	// SegmentFor returns a segment on which to store sz bytes. This may be
	// a new or an existing segment.
	SegmentFor(sz Size) (*Segment, error)

	// Segment returns the specified segment. Returns nil if the segment
	// does not exist.
	Segment(id SegmentID) *Segment

	// Reset clears the list of segments and sets all segments to point to
	// a nil message and data slice.
	Reset()
}

// singleSegmentList is a segment list that only stores a single segment.
type singleSegmentList Segment

func (ssl *singleSegmentList) NumSegments() int                    { return 1 }
func (ssl *singleSegmentList) SegmentFor(_ Size) (*Segment, error) { return (*Segment)(ssl), nil }
func (ssl *singleSegmentList) Reset() {
	ssl.data = nil
	ssl.msg = nil
}
func (ssl *singleSegmentList) Segment(id SegmentID) *Segment {
	if id == 0 {
		return (*Segment)(ssl)
	}
	return nil
}

// multiSegmentList is a segment list that stores segments in a byte slice.
//
// New segments are allocated if none of the existing segments has enough
// capacity for new data.
type multiSegmentList struct {
	segs []Segment
}

func (msl *multiSegmentList) NumSegments() int {
	return len(msl.segs)
}

func (msl *multiSegmentList) SegmentFor(sz Size) (*Segment, error) {
	var seg *Segment
	for i := range msl.segs {
		if hasCapacity(msl.segs[i].data, sz) {
			seg = &msl.segs[i]
			break
		}
	}
	if seg == nil {
		i := len(msl.segs)
		msl.segs = append(msl.segs, Segment{id: SegmentID(i)})
		seg = &msl.segs[i]
	}
	return seg, nil
}

func (msl *multiSegmentList) Segment(id SegmentID) *Segment {
	if int(id) < len(msl.segs) {
		return &msl.segs[int(id)]
	}
	return nil
}

func (msl *multiSegmentList) Reset() {
	for i := range msl.segs {
		msl.segs[i].data = nil
		msl.segs[i].msg = nil
	}
	msl.segs = msl.segs[:0]
}

// arena is an implementation of an Arena that offloads most of its work to an
// associated allocator and segment list.
type arena struct {
	alloc allocator
	segs  segmentList
}

func (a arena) Allocate(sz Size, msg *Message, seg *Segment) (*Segment, address, error) {
	// Determine total allocated amount in the arena.
	var total int64
	for i := 0; i < a.segs.NumSegments(); i++ {
		seg := a.segs.Segment(SegmentID(i))
		if seg == nil {
			return nil, 0, errors.New("segment out of bounds")
		}

		total += int64(len(seg.data))
		if total < 0 {
			return nil, 0, errors.New("overflow attempting to allocate")
		}
	}

	// Determine the slice that will receive new data. Reuse seg if it has
	// enough space for the data, otherwise ask the segment list for a
	// segment to store data in (which may or may not be the same segment).
	var b []byte
	needsClearing := false
	if seg == nil || !hasCapacity(seg.data, sz) {
		var err error

		// Determine the segment to allocate in.
		seg, err = a.segs.SegmentFor(sz)
		if err != nil {
			return nil, 0, err
		}

		b = seg.data
		if !hasCapacity(b, sz) {
			// Size or resize the data.
			b, err = a.alloc.Grow(seg.data, total, sz)
			if err != nil {
				return nil, 0, err
			}
		} else {
			needsClearing = true
		}
	} else {
		b = seg.data
		needsClearing = true
	}

	// The segment's full data is in b[0:], while the buffer requested by
	// the caller is in b[<prior len of buffer>:]. When this was a new
	// segment, the two will be the same.
	//
	// The starting address of the newly allocated space is the end of the
	// prior data.
	addr := address(len(seg.data))
	seg.data = b[:addr.addSizeUnchecked(sz)]
	seg.msg = msg

	// Clear the data after addr to ensure it is zero. The allocators
	// usually already return cleared data, but sometimes a buffer is
	// explicitly passed with left over data, so this ensures the memory
	// that is about to be used is in fact all zeroes.
	if needsClearing {
		// TODO: use clear() once go 1.21 is the minimum required
		// version.
		toClear := seg.data[addr:]
		for i := range toClear {
			toClear[i] = 0
		}
	}

	return seg, addr, nil
}

func (a arena) Release() {
	if a.alloc == nil && a.segs == nil {
		// Empty arena. Use sane defaults.
		a.alloc = (*bufferPoolAllocator)(nil)
		a.segs = &singleSegmentList{}
	}

	for i := 0; i < a.segs.NumSegments(); i++ {
		// Release segment data to the allocator.
		seg := a.segs.Segment(SegmentID(i))
		if seg.data != nil {
			a.alloc.Release(seg.data)
		}
	}

	// Reset list of segments.
	a.segs.Reset()
}

// NumSegments returns the number of segments in the arena.
func (a arena) NumSegments() int64 {
	return int64(a.segs.NumSegments())
}

// Data returns the data in the given segment or an error.
func (a arena) Data(id SegmentID) ([]byte, error) {
	seg := a.segs.Segment(id)
	if seg == nil {
		return nil, errors.New("segment out of bounds")
	}
	return seg.data, nil
}

func (a arena) Segment(id SegmentID) *Segment {
	return a.segs.Segment(id)
}
