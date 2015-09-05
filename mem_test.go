package capnp

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
)

func TestNewMessage(t *testing.T) {
	tests := []struct {
		arena Arena
		fails bool
	}{
		{arena: SingleSegment(nil)},
		{arena: MultiSegment(nil)},
		{arena: readOnlyArena{SingleSegment(make([]byte, 0, 7))}, fails: true},
		{arena: readOnlyArena{SingleSegment(make([]byte, 0, 8))}},
		{arena: MultiSegment(nil)},
		{arena: MultiSegment([][]byte{make([]byte, 8)}), fails: true},
		{arena: MultiSegment([][]byte{incrementingData(8)}), fails: true},
		// This is somewhat arbitrary, but more than one segment = data.
		// This restriction may be lifted if it's not useful.
		{arena: MultiSegment([][]byte{make([]byte, 0, 16), make([]byte, 0)}), fails: true},
	}
	for _, test := range tests {
		msg, seg, err := NewMessage(test.arena)
		if err != nil {
			if !test.fails {
				t.Errorf("NewMessage(%v) failed unexpectedly: %v", test.arena, err)
			}
			continue
		}
		if test.fails {
			t.Errorf("NewMessage(%v) succeeded; want error", test.arena)
			continue
		}
		if n := msg.NumSegments(); n != 1 {
			t.Errorf("NewMessage(%v).NumSegments() = %d; want 1", test.arena, n)
		}
		if seg.ID() != 0 {
			t.Errorf("NewMessage(%v) segment.ID() = %d; want 0", test.arena, seg.ID())
		}
		if len(seg.Data()) != 8 {
			t.Errorf("NewMessage(%v) segment.Data() = % 02x; want length 8", test.arena, seg.Data())
		}
	}
}

func TestAlloc(t *testing.T) {
	type allocTest struct {
		name string

		seg  *Segment
		size Size

		allocID SegmentID
		addr    Address
	}
	var tests []allocTest

	{
		msg := &Message{Arena: SingleSegment(nil)}
		seg, err := msg.Segment(0)
		if err != nil {
			t.Fatal(err)
		}
		tests = append(tests, allocTest{
			name:    "empty alloc in empty segment",
			seg:     seg,
			size:    0,
			allocID: 0,
			addr:    0,
		})
	}
	{
		msg := &Message{Arena: SingleSegment(nil)}
		seg, err := msg.Segment(0)
		if err != nil {
			t.Fatal(err)
		}
		tests = append(tests, allocTest{
			name:    "alloc in empty segment",
			seg:     seg,
			size:    8,
			allocID: 0,
			addr:    0,
		})
	}
	{
		msg := &Message{Arena: MultiSegment([][]byte{
			incrementingData(24)[:8],
			incrementingData(24)[:8],
			incrementingData(24)[:8],
		})}
		seg, err := msg.Segment(1)
		if err != nil {
			t.Fatal(err)
		}
		tests = append(tests, allocTest{
			name:    "prefers given segment",
			seg:     seg,
			size:    16,
			allocID: 1,
			addr:    8,
		})
	}
	{
		msg := &Message{Arena: MultiSegment([][]byte{
			incrementingData(24)[:8],
			incrementingData(24),
		})}
		seg, err := msg.Segment(1)
		if err != nil {
			t.Fatal(err)
		}
		tests = append(tests, allocTest{
			name:    "given segment full with another available",
			seg:     seg,
			size:    16,
			allocID: 0,
			addr:    8,
		})
	}
	{
		msg := &Message{Arena: MultiSegment([][]byte{
			incrementingData(24),
			incrementingData(24),
		})}
		seg, err := msg.Segment(1)
		if err != nil {
			t.Fatal(err)
		}
		tests = append(tests, allocTest{
			name:    "given segment full and no others available",
			seg:     seg,
			size:    16,
			allocID: 2,
			addr:    0,
		})
	}

	for i, test := range tests {
		seg, addr, err := alloc(test.seg, test.size)
		if err != nil {
			t.Errorf("tests[%d] - %s: alloc(..., %d) error: %v", i, test.name, test.size, err)
			continue
		}
		if seg.ID() != test.allocID {
			t.Errorf("tests[%d] - %s: alloc(..., %d) returned segment %d; want segment %d", i, test.name, test.size, seg.ID(), test.allocID)
		}
		if addr != test.addr {
			t.Errorf("tests[%d] - %s: alloc(..., %d) returned address %v; want address %v", i, test.name, test.size, addr, test.addr)
		}
		if !seg.regionInBounds(addr, test.size) {
			t.Errorf("tests[%d] - %s: alloc(..., %d) returned address %v, which is not in bounds (len(seg.data) == %d)", i, test.name, test.size, addr, len(seg.Data()))
		} else if data := seg.slice(addr, test.size); !isZeroFilled(data) {
			t.Errorf("tests[%d] - %s: alloc(..., %d) region has data % 02x; want zero-filled", i, test.name, test.size, data)
		}
	}
}

func TestSingleSegment(t *testing.T) {
	// fresh arena
	{
		arena := SingleSegment(nil)
		if n := arena.NumSegments(); n != 1 {
			t.Errorf("SingleSegment(nil).NumSegments() = %d; want 1", n)
		}
		data, err := arena.Data(0)
		if len(data) != 0 {
			t.Errorf("SingleSegment(nil).Data(0) = %#v; want nil", data)
		}
		if err != nil {
			t.Errorf("SingleSegment(nil).Data(0) error: %v", err)
		}
		_, err = arena.Data(1)
		if err == nil {
			t.Error("SingleSegment(nil).Data(1) succeeded; want error")
		}
	}

	// existing data
	{
		arena := SingleSegment(incrementingData(8))
		if n := arena.NumSegments(); n != 1 {
			t.Errorf("SingleSegment(incrementingData(8)).NumSegments() = %d; want 1", n)
		}
		data, err := arena.Data(0)
		if want := incrementingData(8); !bytes.Equal(data, want) {
			t.Errorf("SingleSegment(incrementingData(8)).Data(0) = %#v; want %#v", data, want)
		}
		if err != nil {
			t.Errorf("SingleSegment(incrementingData(8)).Data(0) error: %v", err)
		}
		_, err = arena.Data(1)
		if err == nil {
			t.Error("SingleSegment(incrementingData(8)).Data(1) succeeded; want error")
		}
	}
}

func TestSingleSegmentAllocate(t *testing.T) {
	tests := []arenaAllocTest{
		{
			name: "empty arena",
			init: func() (Arena, map[SegmentID]*Segment) {
				return SingleSegment(nil), nil
			},
			size: 8,
			id:   0,
			data: []byte{},
		},
		{
			name: "unloaded",
			init: func() (Arena, map[SegmentID]*Segment) {
				buf := incrementingData(24)
				return SingleSegment(buf[:16]), nil
			},
			size: 8,
			id:   0,
			data: incrementingData(16),
		},
		{
			name: "loaded",
			init: func() (Arena, map[SegmentID]*Segment) {
				buf := incrementingData(24)
				buf = buf[:16]
				segs := map[SegmentID]*Segment{
					0: &Segment{id: 0, data: buf},
				}
				return SingleSegment(buf), segs
			},
			size: 8,
			id:   0,
			data: incrementingData(16),
		},
		{
			name: "loaded changes length",
			init: func() (Arena, map[SegmentID]*Segment) {
				buf := incrementingData(32)
				segs := map[SegmentID]*Segment{
					0: &Segment{id: 0, data: buf[:24]},
				}
				return SingleSegment(buf[:16]), segs
			},
			size: 8,
			id:   0,
			data: incrementingData(24),
		},
		{
			name: "message-filled segment",
			init: func() (Arena, map[SegmentID]*Segment) {
				buf := incrementingData(24)
				segs := map[SegmentID]*Segment{
					0: &Segment{id: 0, data: buf},
				}
				return SingleSegment(buf[:16]), segs
			},
			size: 8,
			id:   0,
			data: incrementingData(24),
		},
	}
	for i := range tests {
		tests[i].run(t, i)
	}
}

func TestMultiSegment(t *testing.T) {
	// fresh arena
	{
		arena := MultiSegment(nil)
		if n := arena.NumSegments(); n != 0 {
			t.Errorf("MultiSegment(nil).NumSegments() = %d; want 1", n)
		}
		_, err := arena.Data(0)
		if err == nil {
			t.Error("MultiSegment(nil).Data(0) succeeded; want error")
		}
	}

	// existing data
	{
		arena := MultiSegment([][]byte{incrementingData(8), incrementingData(24)})
		if n := arena.NumSegments(); n != 2 {
			t.Errorf("MultiSegment(...).NumSegments() = %d; want 2", n)
		}
		data, err := arena.Data(0)
		if want := incrementingData(8); !bytes.Equal(data, want) {
			t.Errorf("MultiSegment(...).Data(0) = %#v; want %#v", data, want)
		}
		if err != nil {
			t.Errorf("MultiSegment(...).Data(0) error: %v", err)
		}
		data, err = arena.Data(1)
		if want := incrementingData(24); !bytes.Equal(data, want) {
			t.Errorf("MultiSegment(...).Data(1) = %#v; want %#v", data, want)
		}
		if err != nil {
			t.Errorf("MultiSegment(...).Data(1) error: %v", err)
		}
		_, err = arena.Data(2)
		if err == nil {
			t.Error("MultiSegment(...).Data(2) succeeded; want error")
		}
	}
}

func TestMultiSegmentAllocate(t *testing.T) {
	tests := []arenaAllocTest{
		{
			name: "empty arena",
			init: func() (Arena, map[SegmentID]*Segment) {
				return MultiSegment(nil), nil
			},
			size: 8,
			id:   0,
			data: []byte{},
		},
		{
			name: "space in unloaded segment",
			init: func() (Arena, map[SegmentID]*Segment) {
				buf := incrementingData(24)
				return MultiSegment([][]byte{buf[:16]}), nil
			},
			size: 8,
			id:   0,
			data: incrementingData(16),
		},
		{
			name: "space in loaded segment",
			init: func() (Arena, map[SegmentID]*Segment) {
				buf := incrementingData(24)
				buf = buf[:16]
				segs := map[SegmentID]*Segment{
					0: &Segment{id: 0, data: buf},
				}
				return MultiSegment([][]byte{buf}), segs
			},
			size: 8,
			id:   0,
			data: incrementingData(16),
		},
		{
			name: "space in loaded segment changes length",
			init: func() (Arena, map[SegmentID]*Segment) {
				buf := incrementingData(32)
				segs := map[SegmentID]*Segment{
					0: &Segment{id: 0, data: buf[:24]},
				}
				return MultiSegment([][]byte{buf[:16]}), segs
			},
			size: 8,
			id:   0,
			data: incrementingData(24),
		},
		{
			name: "message-filled segment",
			init: func() (Arena, map[SegmentID]*Segment) {
				buf := incrementingData(24)
				segs := map[SegmentID]*Segment{
					0: &Segment{id: 0, data: buf},
				}
				return MultiSegment([][]byte{buf[:16]}), segs
			},
			size: 8,
			id:   1,
			data: []byte{},
		},
	}

	for i := range tests {
		tests[i].run(t, i)
	}
}

type arenaAllocTest struct {
	name string

	// Arrange
	init func() (Arena, map[SegmentID]*Segment)
	size Size

	// Assert
	id   SegmentID
	data []byte
}

func (test *arenaAllocTest) run(t *testing.T, i int) {
	arena, segs := test.init()
	id, data, err := arena.Allocate(test.size, segs)

	if err != nil {
		t.Errorf("tests[%d] - %s: Allocate error: %v", i, test.name, err)
		return
	}
	if id != test.id {
		t.Errorf("tests[%d] - %s: Allocate id = %d; want %d", i, test.name, id, test.id)
	}
	if !bytes.Equal(data, test.data) {
		t.Errorf("tests[%d] - %s: Allocate data = % 02x; want % 02x", i, test.name, data, test.data)
	}
	if Size(cap(data)-len(data)) < test.size {
		t.Errorf("tests[%d] - %s: Allocate len(data) = %d, cap(data) = %d; cap(data) should be at least %d", len(data), cap(data), Size(len(data))+test.size)
	}
}

func incrementingData(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(i % 256)
	}
	return b
}

func isZeroFilled(b []byte) bool {
	for _, bb := range b {
		if bb != 0 {
			return false
		}
	}
	return true
}

type readOnlyArena struct {
	Arena
}

func (ro readOnlyArena) String() string {
	return fmt.Sprintf("readOnlyArena{%v}", ro.Arena)
}

func (readOnlyArena) Allocate(sz Size, segs map[SegmentID]*Segment) (SegmentID, []byte, error) {
	return 0, nil, errReadOnlyArena
}

var errReadOnlyArena = errors.New("Allocate called on read-only arena")
