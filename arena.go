package capnp

import (
	"errors"

	"capnproto.org/go/capnp/v3/exc"
)

// An Arena loads and allocates segments for a Message.
type Arena interface {
	// NumSegments returns the number of segments in the arena.
	// This must not be larger than 1<<32.
	NumSegments() int64

	// Data loads the data for the segment with the given ID.  IDs are in
	// the range [0, NumSegments()).
	// must be tightly packed in the range [0, NumSegments()).
	//
	// TODO: remove in favor of Segment(x).Data().
	// Deprecated.
	Data(id SegmentID) ([]byte, error)

	// Segment returns the segment identified with the specified id. This
	// may return nil if the segment with the specified ID does not exist.
	Segment(id SegmentID) *Segment

	// Allocate selects a segment to place a new object in, creating a
	// segment or growing the capacity of a previously loaded segment if
	// necessary.  If Allocate does not return an error, then the
	// difference of the capacity and the length of the returned slice
	// must be at least minsz.  Some allocators may specifically choose
	// to grow the passed seg (if non nil).
	//
	// If Allocate creates a new segment, the ID must be one larger than
	// the last segment's ID or zero if it is the first segment.
	//
	// If Allocate returns an previously loaded segment's ID, then the
	// arena is responsible for preserving the existing data.
	Allocate(sz Size, msg *Message, seg *Segment) (*Segment, address, error)

	// Release all resources associated with the Arena. Callers MUST NOT
	// use the Arena after it has been released.
	//
	// Calling Release() is OPTIONAL, but may reduce allocations.
	//
	// Implementations MAY use Release() as a signal to return resources
	// to free lists, or otherwise reuse the Arena.   However, they MUST
	// NOT assume Release() will be called.
	Release()
}

// SingleSegmentArena is an Arena implementation that stores message data
// in a continguous slice.  Allocation is performed by first allocating a
// new slice and copying existing data. SingleSegment arena does not fail
// unless the caller attempts to access another segment.
func SingleSegment(b []byte) Arena {
	var alloc allocator = (*bufferPoolAllocator)(nil)
	if b != nil {
		// When b is specified, do not return the buffer to any
		// caches, because we don't know where the caller got the
		// buffer from.
		alloc = simpleAllocator{}
	}
	return arena{
		alloc: alloc,
		segs:  &singleSegmentList{data: b},
	}
}

// MultiSegment returns a new arena that allocates new segments when
// they are full.  b MAY be nil.  Callers MAY use b to populate the
// buffer for reading or to reserve memory of a specific size.
func MultiSegment(b [][]byte) Arena {
	var alloc allocator = (*bufferPoolAllocator)(nil)
	var segs []Segment
	if b != nil {
		// When b is specified, do not return the buffer to any
		// caches, because we don't know where the caller got the
		// buffer from.
		alloc = simpleAllocator{}
		segs = make([]Segment, len(b))
		for i := range b {
			segs[i] = Segment{id: SegmentID(i), data: b[i]}
		}
	}
	return arena{
		alloc: alloc,
		segs:  &multiSegmentList{segs: segs},
	}
}

// demuxArena demuxes a byte slice (that contains data for a list of
// segments identified on the header) into an appropriate arena.
func demuxArena(hdr streamHeader, data []byte) (Arena, error) {
	maxSeg := hdr.maxSegment()
	if int64(maxSeg) > int64(maxInt-1) {
		return arena{}, errors.New("number of segments overflows int")
	}

	if maxSeg == 0 && len(data) == 0 {
		return SingleSegment(nil), nil
	}

	segBufs := make([][]byte, maxSeg+1)
	off := 0
	for i := range segBufs {
		sz, err := hdr.segmentSize(SegmentID(i))
		if err != nil {
			return arena{}, exc.WrapError("decode", err)
		}
		segBufs[i] = data[off : off+int(sz)]
		off += int(sz)
	}

	return MultiSegment(segBufs), nil
}

// nextAlloc computes how much more space to allocate given the number
// of bytes allocated in the entire message and the requested number of
// bytes.  It will always return a multiple of wordSize.  max must be a
// multiple of wordSize.  The sum of curr and the returned size will
// always be less than max.
func nextAlloc(curr int64, req Size) (int, error) {

	if req == 0 {
		return 0, nil
	}

	req = req.padToWord()
	totalWant := curr + int64(req)

	// Sanity checks.
	if req > maxAllocSize() {
		return 0, errors.New("alloc " + req.String() + ": too large")
	}
	if totalWant <= curr || totalWant > int64(maxAllocSize()) {
		return 0, errors.New("alloc " + req.String() + ": message size overflow")
	}

	if totalWant < 1024 {
		return int(req), nil
	}

	// doubleCurr is double the current total message size (padded to the
	// word sized).
	doubleCurr := (curr*2 + 7) &^ 7

	// The following attempts to amortize allocation costs across a wide
	// range of uses.
	switch {
	case totalWant < 1024*1024 && doubleCurr < 1024*1024:
		// When doubling the message size still keeps the message under
		// 1MiB, double the message size.
		return int(doubleCurr), nil

	default:
		// Otherwise, allocate the requested amount.
		return int(req), nil
	}
}

func hasCapacity(b []byte, sz Size) bool {
	return sz <= Size(cap(b)-len(b))
}

// ReadOnlySingleSegmentArena is a single segment arena backed by a byte slice
// that does not allow allocations.
type ReadOnlySingleSegmentArena Segment

func (a *ReadOnlySingleSegmentArena) NumSegments() int64 {
	return 1
}

func (a *ReadOnlySingleSegmentArena) Data(id SegmentID) ([]byte, error) {
	if id != 0 {
		return nil, errors.New("segment out of bounds")
	}
	return a.data, nil
}

// Segment returns the segment identified with the specified id. This
// may return nil if the segment with the specified ID does not exist.
func (a *ReadOnlySingleSegmentArena) Segment(id SegmentID) *Segment {
	if id > 0 {
		return nil
	}
	return (*Segment)(a)
}

func (a *ReadOnlySingleSegmentArena) Allocate(sz Size, msg *Message, seg *Segment) (*Segment, address, error) {
	return nil, 0, errors.New("ReadOnlySingleSegmentArena cannot allocate")
}

func (a *ReadOnlySingleSegmentArena) Release() {
	// This does nothing.
}

// UseBuffer switches the internal buffer to use the specified one.
func (a *ReadOnlySingleSegmentArena) UseBuffer(b []byte) {
	a.data = b
}

// SimpleSingleSegmentArena is an alternative implementation of a single
// segment arena that is not subject to the same legacy behavior as the the
// single segment arena initialized by SingleSegmentArena.
//
// This arena is not safe for concurrent access and holds onto the last
// allocated buffer for reuse after a call to Relese().
type SimpleSingleSegmentArena Segment

func (a *SimpleSingleSegmentArena) NumSegments() int64 { return 1 }

func (a *SimpleSingleSegmentArena) Data(id SegmentID) ([]byte, error) {
	if id != 0 {
		return nil, errors.New("segment out of bounds")
	}
	return a.data, nil
}

func (a *SimpleSingleSegmentArena) Segment(id SegmentID) *Segment {
	if id != 0 {
		return nil
	}
	return (*Segment)(a)
}

func (a *SimpleSingleSegmentArena) Allocate(sz Size, msg *Message, seg *Segment) (*Segment, address, error) {
	totalMsgSize := int64(len(a.data))
	inc, err := nextAlloc(totalMsgSize, sz)
	if err != nil {
		return nil, 0, err
	}

	addr := address(len(a.data))
	capNeeded := len(a.data) + inc
	if capNeeded <= cap(a.data) {
		a.data = a.data[:capNeeded]
	} else {
		a.data = append(a.data, make([]byte, inc)...)
	}
	return (*Segment)(a), addr, nil
}

func (a *SimpleSingleSegmentArena) Release() {
	for i := range a.data {
		a.data[i] = 0
	}
	a.data = a.data[:0]
}

// ReplaceBuffer replaces the internal buffer with a new buffer. This
// effectively resets a message to the one encoded by the passed buffer, which
// should be a single segment message.
func (a *SimpleSingleSegmentArena) ReplaceBuffer(b []byte) {
	a.data = b
}
