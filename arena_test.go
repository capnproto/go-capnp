package capnp

import (
	"bytes"
	"errors"
	"testing"

	"capnproto.org/go/capnp/v3/exp/bufferpool"
	"github.com/stretchr/testify/require"
)

type arenaAllocTest struct {
	name string

	// Arrange
	init func() (Arena, map[SegmentID]*Segment)
	size Size

	// Assert
	id   SegmentID
	data []byte
}

func (test *arenaAllocTest) run(t *testing.T) {
	arena, _ := test.init()
	seg, addr, err := arena.Allocate(test.size, nil, nil)

	require.NoError(t, err, "Allocate error")
	require.Equal(t, test.id, seg.id)

	// Allocate() contract is that segment data starting at addr should
	// have anough room for test.size bytes.
	require.Less(t, int(addr), len(seg.data))

	data := seg.data[addr:]
	require.LessOrEqual(t, test.size, Size(cap(seg.data)))

	data = data[:test.size]
	require.Equal(t, test.data, data)
}

func incrementingData(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(i % 256)
	}
	return b
}

func segmentData(a Arena, id SegmentID) ([]byte, error) {
	seg := a.Segment(id)
	if seg == nil {
		return nil, errors.New("segment does not exist")
	}

	return seg.Data(), nil
}

func TestSingleSegment(t *testing.T) {
	t.Parallel()
	t.Helper()

	t.Run("FreshArena", func(t *testing.T) {
		t.Parallel()

		arena := SingleSegment(nil)
		if n := arena.NumSegments(); n != 1 {
			t.Errorf("SingleSegment(nil).NumSegments() = %d; want 1", n)
		}
		data, err := segmentData(arena, 0)
		if len(data) != 0 {
			t.Errorf("SingleSegment(nil).Data(0) = %#v; want nil", data)
		}
		if err != nil {
			t.Errorf("SingleSegment(nil).Data(0) error: %v", err)
		}
		_, err = segmentData(arena, 1)
		if err == nil {
			t.Error("SingleSegment(nil).Data(1) succeeded; want error")
		}
	})

	t.Run("ExistingData", func(t *testing.T) {
		t.Parallel()

		arena := SingleSegment(incrementingData(8))
		if n := arena.NumSegments(); n != 1 {
			t.Errorf("SingleSegment(incrementingData(8)).NumSegments() = %d; want 1", n)
		}
		data, err := segmentData(arena, 0)
		if want := incrementingData(8); !bytes.Equal(data, want) {
			t.Errorf("SingleSegment(incrementingData(8)).Data(0) = %#v; want %#v", data, want)
		}
		if err != nil {
			t.Errorf("SingleSegment(incrementingData(8)).Data(0) error: %v", err)
		}
		_, err = segmentData(arena, 1)
		if err == nil {
			t.Error("SingleSegment(incrementingData(8)).Data(1) succeeded; want error")
		}
	})
}

func TestSingleSegmentAllocate(t *testing.T) {
	t.Parallel()

	tests := []arenaAllocTest{
		{
			name: "empty arena",
			init: func() (Arena, map[SegmentID]*Segment) {
				return SingleSegment(nil), nil
			},
			size: 8,
			id:   0,
			data: []byte{7: 0},
		},
		{
			name: "unloaded",
			init: func() (Arena, map[SegmentID]*Segment) {
				buf := incrementingData(24)
				return SingleSegment(buf[:16]), nil
			},
			size: 8,
			id:   0,
			data: incrementingData(24)[16 : 16+8],
		},
		{
			name: "loaded",
			init: func() (Arena, map[SegmentID]*Segment) {
				buf := incrementingData(24)
				buf = buf[:16]
				segs := map[SegmentID]*Segment{
					0: {id: 0, data: buf},
				}
				return SingleSegment(buf), segs
			},
			size: 8,
			id:   0,
			data: incrementingData(24)[16 : 16+8],
		},
		{
			name: "loaded changes length",
			init: func() (Arena, map[SegmentID]*Segment) {
				buf := incrementingData(32)
				segs := map[SegmentID]*Segment{
					0: {id: 0, data: buf[:24]},
				}
				return SingleSegment(buf[:16]), segs
			},
			size: 8,
			id:   0,
			data: incrementingData(32)[16 : 16+8],
		},
		{
			name: "message-filled segment",
			init: func() (Arena, map[SegmentID]*Segment) {
				buf := incrementingData(24)
				segs := map[SegmentID]*Segment{
					0: {id: 0, data: buf},
				}
				return SingleSegment(buf[:16]), segs
			},
			size: 8,
			id:   0,
			data: incrementingData(24)[16 : 16+8],
		},
	}
	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, tc.run)
	}
}

func TestMultiSegment(t *testing.T) {
	t.Parallel()
	t.Helper()

	t.Run("FreshArena", func(t *testing.T) {
		t.Parallel()

		arena := MultiSegment(nil)
		if n := arena.NumSegments(); n != 0 {
			t.Errorf("MultiSegment(nil).NumSegments() = %d; want 1", n)
		}
		_, err := segmentData(arena, 0)
		if err == nil {
			t.Error("MultiSegment(nil).Data(0) succeeded; want error")
		}
	})

	t.Run("ExistingData", func(t *testing.T) {
		t.Parallel()

		arena := MultiSegment([][]byte{incrementingData(8), incrementingData(24)})
		if n := arena.NumSegments(); n != 2 {
			t.Errorf("MultiSegment(...).NumSegments() = %d; want 2", n)
		}
		data, err := segmentData(arena, 0)
		if want := incrementingData(8); !bytes.Equal(data, want) {
			t.Errorf("MultiSegment(...).Data(0) = %#v; want %#v", data, want)
		}
		if err != nil {
			t.Errorf("MultiSegment(...).Data(0) error: %v", err)
		}
		data, err = segmentData(arena, 1)
		if want := incrementingData(24); !bytes.Equal(data, want) {
			t.Errorf("MultiSegment(...).Data(1) = %#v; want %#v", data, want)
		}
		if err != nil {
			t.Errorf("MultiSegment(...).Data(1) error: %v", err)
		}
		_, err = segmentData(arena, 2)
		if err == nil {
			t.Error("MultiSegment(...).Data(2) succeeded; want error")
		}
	})
}

func TestMultiSegmentAllocate(t *testing.T) {
	t.Parallel()

	tests := []arenaAllocTest{
		{
			name: "empty arena",
			init: func() (Arena, map[SegmentID]*Segment) {
				return MultiSegment(nil), nil
			},
			size: 8,
			id:   0,
			data: []byte{7: 0},
		},
		{
			name: "space in unloaded segment",
			init: func() (Arena, map[SegmentID]*Segment) {
				buf := incrementingData(24)
				return MultiSegment([][]byte{buf[:16]}), nil
			},
			size: 8,
			id:   0,
			data: incrementingData(24)[16 : 16+8],
		},
		{
			name: "space in loaded segment",
			init: func() (Arena, map[SegmentID]*Segment) {
				buf := incrementingData(24)
				buf = buf[:16]
				segs := map[SegmentID]*Segment{
					0: {id: 0, data: buf},
				}
				return MultiSegment([][]byte{buf}), segs
			},
			size: 8,
			id:   0,
			data: incrementingData(24)[16 : 16+8],
		},
		{
			name: "space in loaded segment changes length",
			init: func() (Arena, map[SegmentID]*Segment) {
				buf := incrementingData(32)
				segs := map[SegmentID]*Segment{
					0: {id: 0, data: buf[:24]},
				}
				return MultiSegment([][]byte{buf[:16]}), segs
			},
			size: 8,
			id:   0,
			data: incrementingData(24)[16 : 16+8],
		},
		{
			name: "first segment is filled",
			init: func() (Arena, map[SegmentID]*Segment) {
				buf := incrementingData(24)
				segs := map[SegmentID]*Segment{
					0: {id: 0, data: buf},
				}
				msa := MultiSegment([][]byte{buf[:16:16]})
				msa.bp = &bufferpool.Default
				return msa, segs
			},
			size: 8,
			id:   1,
			data: []byte{7: 0},
		},
	}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, tc.run)
	}
}
