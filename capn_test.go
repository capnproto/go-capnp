package capnp

import (
	"testing"
)

func TestMakeOffsetKey(t *testing.T) {
	seg42 := &Segment{id: 42}
	tests := []struct {
		p          Pointer
		id         SegmentID
		boff, bend int64
	}{
		{
			p: Struct{
				seg:  seg42,
				off:  0,
				size: ObjectSize{0, 0},
			},
			id:   42,
			boff: 0,
			bend: 0,
		},
		{
			p: Struct{
				seg:  seg42,
				off:  8,
				size: ObjectSize{0, 0},
			},
			id:   42,
			boff: 64,
			bend: 64,
		},
		{
			p: Struct{
				seg:  seg42,
				off:  8,
				size: ObjectSize{1, 0},
			},
			id:   42,
			boff: 64,
			bend: 72,
		},
		{
			p: Struct{
				seg:  seg42,
				off:  8,
				size: ObjectSize{0, 1},
			},
			id:   42,
			boff: 64,
			bend: 128,
		},
		{
			p: List{
				seg:    seg42,
				off:    0,
				size:   ObjectSize{},
				length: 0,
			},
			id:   42,
			boff: 0,
			bend: 0,
		},
		{
			p: List{
				seg:    seg42,
				off:    0,
				size:   ObjectSize{},
				length: 1,
			},
			id:   42,
			boff: 0,
			bend: 0,
		},
		{
			p: List{
				seg:    seg42,
				off:    0,
				size:   ObjectSize{0, 1},
				length: 1,
			},
			id:   42,
			boff: 0,
			bend: 64,
		},
		{
			p: List{
				seg:    seg42,
				off:    8,
				size:   ObjectSize{0, 1},
				length: 1,
			},
			id:   42,
			boff: 64,
			bend: 128,
		},
		{
			p: List{
				seg:    seg42,
				off:    8,
				size:   ObjectSize{0, 1},
				length: 1,
				flags:  isCompositeList,
			},
			id:   42,
			boff: 0,
			bend: 128,
		},
		{
			p: List{
				seg:    seg42,
				off:    8,
				size:   ObjectSize{0, 1},
				length: 2,
			},
			id:   42,
			boff: 64,
			bend: 192,
		},
	}
	for _, test := range tests {
		off := makeOffsetKey(test.p)
		if off.id != test.id || off.boff != test.boff || off.bend != test.bend {
			t.Errorf("makeOffsetKey(%#v) = offset{id: %d, boff: %d, bend: %d}; want offset{id: %d, boff: %d, bend: %d}", test.p, off.id, off.boff, off.bend, test.id, test.boff, test.bend)
		}
	}
}
