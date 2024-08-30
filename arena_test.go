package capnp

import (
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

func segmentData(a Arena, id SegmentID) []byte {
	seg := a.Segment(id)
	if seg == nil {
		return nil
	}

	return seg.Data()
}

func TestSingleSegment(t *testing.T) {
	t.Parallel()
	t.Helper()

	t.Run("FreshArena", func(t *testing.T) {
		t.Parallel()

		arena := SingleSegment(nil)
		require.Equal(t, int64(1), arena.NumSegments())
		data0 := segmentData(arena, 0)
		require.Empty(t, data0)
		data1 := segmentData(arena, 1)
		require.Empty(t, data1)
	})

	t.Run("ExistingData", func(t *testing.T) {
		t.Parallel()

		arena := SingleSegment(incrementingData(8))
		require.Equal(t, int64(1), arena.NumSegments())
		data0 := segmentData(arena, 0)
		require.Equal(t, incrementingData(8), data0)
		data1 := segmentData(arena, 1)
		require.Empty(t, data1)
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
		require.Equal(t, int64(0), arena.NumSegments())
		data0 := segmentData(arena, 0)
		require.Empty(t, data0)
	})

	t.Run("ExistingData", func(t *testing.T) {
		t.Parallel()

		arena := MultiSegment([][]byte{incrementingData(8), incrementingData(24)})
		require.Equal(t, int64(2), arena.NumSegments())
		data0 := segmentData(arena, 0)
		require.Equal(t, incrementingData(8), data0)
		data1 := segmentData(arena, 1)
		require.Equal(t, incrementingData(24), data1)
		data2 := segmentData(arena, 2)
		require.Empty(t, data2)
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
