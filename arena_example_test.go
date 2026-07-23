package capnp_test

import (
	"errors"
	"fmt"

	"capnproto.org/go/capnp/v3"
)

var (
	errFixedArenaCapacity = errors.New("fixed arena: capacity exceeded")
	errFixedArenaSegment  = errors.New("fixed arena: segment is not associated with arena")
)

// fixedArena is a single-segment Arena backed by caller-owned storage.
// It never grows, copies, pools, or adds segments.
type fixedArena struct {
	seg *capnp.Segment
}

func newFixedArena(buffer []byte) fixedArena {
	// Starting with buffer[:0] keeps data word-aligned during normal message
	// construction, which passes word-padded allocation sizes to Arena.Allocate.
	return fixedArena{seg: capnp.NewSegment(0, buffer[:0])}
}

func (*fixedArena) NumSegments() int64 {
	return 1
}

func (a *fixedArena) Segment(id capnp.SegmentID) *capnp.Segment {
	if id == 0 {
		return a.seg
	}
	return nil
}

func (a *fixedArena) Allocate(minsz capnp.Size, msg *capnp.Message, preferred *capnp.Segment) (*capnp.Segment, capnp.Address, error) {
	if preferred == nil || preferred != a.seg {
		return nil, 0, errFixedArenaSegment
	}
	data := a.seg.Data()
	if minsz > capnp.Size(cap(data)-len(data)) {
		// Callers must Reset the Message before reusing the Arena after a
		// construction failure.
		return nil, 0, errFixedArenaCapacity
	}
	addr := capnp.Address(len(data))
	a.seg.SetData(data[:len(data)+int(minsz)])
	a.seg.BindTo(msg)
	return a.seg, addr, nil
}

func (a *fixedArena) Release() {
	// This makes the previous message unavailable without scrubbing buffer.
	// Callers that need to clear sensitive data must do so themselves.
	a.seg.SetData(a.seg.Data()[:0])
	a.seg.BindTo(nil)
}

func ExampleArena() {
	// The caller owns this buffer for the arena's entire lifetime.
	var buffer [64]byte
	arena := newFixedArena(buffer[:])
	msg, seg, err := capnp.NewMessage(&arena)
	if err != nil {
		panic(err)
	}

	root, err := capnp.NewRootStruct(seg, capnp.ObjectSize{PointerCount: 1})
	if err != nil {
		panic(err)
	}
	if err := root.SetText(0, "first"); err != nil {
		panic(err)
	}

	// Reset invalidates root and all other pointers from the first message.
	seg, err = msg.Reset(&arena)
	if err != nil {
		panic(err)
	}
	root, err = capnp.NewRootStruct(seg, capnp.ObjectSize{PointerCount: 1})
	if err != nil {
		panic(err)
	}
	if err := root.SetText(0, "second"); err != nil {
		panic(err)
	}
	p, err := root.Ptr(0)
	if err != nil {
		panic(err)
	}
	fmt.Println(p.Text())

	// Output:
	// second
}
