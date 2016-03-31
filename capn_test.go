package capnp

import (
	"bytes"
	"fmt"
	"testing"
)

func TestSegmentInBounds(t *testing.T) {
	tests := []struct {
		n    int
		addr Address
		ok   bool
	}{
		{0, 0, false},
		{0, 1, false},
		{0, 2, false},
		{1, 0, true},
		{1, 1, false},
		{1, 2, false},
		{2, 0, true},
		{2, 1, true},
		{2, 2, false},
	}
	for _, test := range tests {
		seg := &Segment{data: make([]byte, test.n)}
		if ok := seg.inBounds(test.addr); ok != test.ok {
			t.Errorf("&Segment{data: make([]byte, %d)}.inBounds(%#v) = %t; want %t", test.n, test.addr, ok, test.ok)
		}
	}
}

func TestSegmentRegionInBounds(t *testing.T) {
	tests := []struct {
		n    int
		addr Address
		sz   Size
		ok   bool
	}{
		{0, 0, 0, true}, // zero-sized region <= len is okay
		{0, 0, 1, false},
		{0, 1, 0, false},
		{0, 1, 1, false},
		{1, 0, 0, true},
		{1, 0, 1, true},
		{1, 1, 0, true},
		{1, 1, 1, false},
		{2, 0, 0, true},
		{2, 0, 1, true},
		{2, 0, 2, true},
		{2, 0, 3, false},
		{2, 1, 0, true},
		{2, 1, 1, true},
		{2, 1, 2, false},
		{2, 1, 3, false},
		{2, 2, 0, true},
		{2, 2, 1, false},
	}
	for _, test := range tests {
		seg := &Segment{data: make([]byte, test.n)}
		if ok := seg.regionInBounds(test.addr, test.sz); ok != test.ok {
			t.Errorf("&Segment{data: make([]byte, %d)}.regionInBounds(%#v, %#v) = %t; want %t", test.n, test.addr, test.sz, ok, test.ok)
		}
	}
}

func TestSegmentReadUint8(t *testing.T) {
	tests := []struct {
		data   []byte
		addr   Address
		val    uint8
		panics bool
	}{
		{data: []byte{}, addr: 0, panics: true},
		{data: []byte{42}, addr: 0, val: 42},
		{data: []byte{42}, addr: 1, panics: true},
		{data: []byte{1, 42, 2}, addr: 0, val: 1},
		{data: []byte{1, 42, 2}, addr: 1, val: 42},
		{data: []byte{1, 42, 2}, addr: 2, val: 2},
		{data: []byte{1, 42, 2}, addr: 3, panics: true},
	}
	for _, test := range tests {
		seg := &Segment{data: test.data}
		var val uint8
		err := catchPanic(func() {
			val = seg.readUint8(test.addr)
		})
		if err != nil {
			if !test.panics {
				t.Errorf("&Segment{data: % x}.readUint8(%v) unexpected panic: %v", test.data, test.addr, err)
			}
			continue
		}
		if test.panics {
			t.Errorf("&Segment{data: % x}.readUint8(%v) did not panic as expected", test.data, test.addr)
			continue
		}
		if val != test.val {
			t.Errorf("&Segment{data: % x}.readUint8(%v) = %#x; want %#x", test.data, test.addr, val, test.val)
		}
	}
}

func TestSegmentReadUint16(t *testing.T) {
	tests := []struct {
		data   []byte
		addr   Address
		val    uint16
		panics bool
	}{
		{data: []byte{}, addr: 0, panics: true},
		{data: []byte{0x00}, addr: 0, panics: true},
		{data: []byte{0x00, 0x00}, addr: 0, val: 0},
		{data: []byte{0x01, 0x00}, addr: 0, val: 1},
		{data: []byte{0x34, 0x12}, addr: 0, val: 0x1234},
		{data: []byte{0x34, 0x12, 0x56}, addr: 0, val: 0x1234},
		{data: []byte{0x34, 0x12, 0x56}, addr: 1, val: 0x5612},
		{data: []byte{0x34, 0x12, 0x56}, addr: 2, panics: true},
	}
	for _, test := range tests {
		seg := &Segment{data: test.data}
		var val uint16
		err := catchPanic(func() {
			val = seg.readUint16(test.addr)
		})
		if err != nil {
			if !test.panics {
				t.Errorf("&Segment{data: % x}.readUint16(%v) unexpected panic: %v", test.data, test.addr, err)
			}
			continue
		}
		if test.panics {
			t.Errorf("&Segment{data: % x}.readUint16(%v) did not panic as expected", test.data, test.addr)
			continue
		}
		if val != test.val {
			t.Errorf("&Segment{data: % x}.readUint16(%v) = %#x; want %#x", test.data, test.addr, val, test.val)
		}
	}
}

func TestSegmentReadUint32(t *testing.T) {
	tests := []struct {
		data   []byte
		addr   Address
		val    uint32
		panics bool
	}{
		{data: []byte{}, addr: 0, panics: true},
		{data: []byte{0x00}, addr: 0, panics: true},
		{data: []byte{0x00, 0x00}, addr: 0, panics: true},
		{data: []byte{0x00, 0x00, 0x00}, addr: 0, panics: true},
		{data: []byte{0x00, 0x00, 0x00, 0x00}, addr: 0, val: 0},
		{data: []byte{0x78, 0x56, 0x34, 0x12}, addr: 0, val: 0x12345678},
		{data: []byte{0xff, 0x78, 0x56, 0x34, 0x12, 0xff}, addr: 1, val: 0x12345678},
		{data: []byte{0xff, 0x78, 0x56, 0x34, 0x12, 0xff}, addr: 2, val: 0xff123456},
		{data: []byte{0xff, 0x78, 0x56, 0x34, 0x12, 0xff}, addr: 3, panics: true},
	}
	for _, test := range tests {
		seg := &Segment{data: test.data}
		var val uint32
		err := catchPanic(func() {
			val = seg.readUint32(test.addr)
		})
		if err != nil {
			if !test.panics {
				t.Errorf("&Segment{data: % x}.readUint32(%v) unexpected panic: %v", test.data, test.addr, err)
			}
			continue
		}
		if test.panics {
			t.Errorf("&Segment{data: % x}.readUint32(%v) did not panic as expected", test.data, test.addr)
			continue
		}
		if val != test.val {
			t.Errorf("&Segment{data: % x}.readUint32(%v) = %#x; want %#x", test.data, test.addr, val, test.val)
		}
	}
}

func TestSegmentReadUint64(t *testing.T) {
	tests := []struct {
		data   []byte
		addr   Address
		val    uint64
		panics bool
	}{
		{data: []byte{}, addr: 0, panics: true},
		{data: []byte{0x00}, addr: 0, panics: true},
		{data: []byte{0x00, 0x00}, addr: 0, panics: true},
		{data: []byte{0x00, 0x00, 0x00}, addr: 0, panics: true},
		{data: []byte{0x00, 0x00, 0x00, 0x00}, addr: 0, panics: true},
		{data: []byte{0x00, 0x00, 0x00, 0x00, 0x00}, addr: 0, panics: true},
		{data: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, addr: 0, panics: true},
		{data: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, addr: 0, panics: true},
		{data: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, addr: 0, val: 0},
		{data: []byte{0xef, 0xcd, 0xab, 0x89, 0x67, 0x45, 0x23, 0x01}, addr: 0, val: 0x0123456789abcdef},
		{data: []byte{0xff, 0xef, 0xcd, 0xab, 0x89, 0x67, 0x45, 0x23, 0x01, 0xff}, addr: 0, val: 0x23456789abcdefff},
		{data: []byte{0xff, 0xef, 0xcd, 0xab, 0x89, 0x67, 0x45, 0x23, 0x01, 0xff}, addr: 1, val: 0x0123456789abcdef},
		{data: []byte{0xff, 0xef, 0xcd, 0xab, 0x89, 0x67, 0x45, 0x23, 0x01, 0xff}, addr: 2, val: 0xff0123456789abcd},
		{data: []byte{0xff, 0xef, 0xcd, 0xab, 0x89, 0x67, 0x45, 0x23, 0x01, 0xff}, addr: 3, panics: true},
	}
	for _, test := range tests {
		seg := &Segment{data: test.data}
		var val uint64
		err := catchPanic(func() {
			val = seg.readUint64(test.addr)
		})
		if err != nil {
			if !test.panics {
				t.Errorf("&Segment{data: % x}.readUint64(%v) unexpected panic: %v", test.data, test.addr, err)
			}
			continue
		}
		if test.panics {
			t.Errorf("&Segment{data: % x}.readUint64(%v) did not panic as expected", test.data, test.addr)
			continue
		}
		if val != test.val {
			t.Errorf("&Segment{data: % x}.readUint64(%v) = %#x; want %#x", test.data, test.addr, val, test.val)
		}
	}
}

func TestSegmentWriteUint8(t *testing.T) {
	tests := []struct {
		data   []byte
		addr   Address
		val    uint8
		out    []byte
		panics bool
	}{
		{
			data:   []byte{},
			addr:   0,
			val:    0,
			panics: true,
		},
		{
			data: []byte{1},
			addr: 0,
			val:  42,
			out:  []byte{42},
		},
		{
			data:   []byte{42},
			addr:   1,
			val:    1,
			panics: true,
		},
		{
			data: []byte{1, 2, 3},
			addr: 0,
			val:  0xff,
			out:  []byte{0xff, 2, 3},
		},
		{
			data: []byte{1, 2, 3},
			addr: 1,
			val:  0xff,
			out:  []byte{1, 0xff, 3},
		},
		{
			data: []byte{1, 2, 3},
			addr: 2,
			val:  0xff,
			out:  []byte{1, 2, 0xff},
		},
		{
			data:   []byte{1, 2, 3},
			addr:   3,
			val:    0xff,
			panics: true,
		},
	}
	for _, test := range tests {
		out := make([]byte, len(test.data))
		copy(out, test.data)
		seg := &Segment{data: out}
		err := catchPanic(func() {
			seg.writeUint8(test.addr, test.val)
		})
		if err != nil {
			if !test.panics {
				t.Errorf("&Segment{data: % x}.writeUint8(%v, %#x) unexpected panic: %v", test.data, test.addr, test.val, err)
			}
			continue
		}
		if test.panics {
			t.Errorf("&Segment{data: % x}.writeUint8(%v, %#x) did not panic as expected", test.data, test.addr, test.val)
			continue
		}
		if !bytes.Equal(out, test.out) {
			t.Errorf("data after &Segment{data: % x}.writeUint8(%v, %#x) = % x; want % x", test.data, test.addr, test.val, out, test.out)
		}
	}
}

func TestSegmentWriteUint16(t *testing.T) {
	tests := []struct {
		data   []byte
		addr   Address
		val    uint16
		out    []byte
		panics bool
	}{
		{
			data:   []byte{},
			addr:   0,
			val:    0,
			panics: true,
		},
		{
			data: []byte{1, 2, 3, 4},
			addr: 1,
			val:  0x1234,
			out:  []byte{1, 0x34, 0x12, 4},
		},
	}
	for _, test := range tests {
		out := make([]byte, len(test.data))
		copy(out, test.data)
		seg := &Segment{data: out}
		err := catchPanic(func() {
			seg.writeUint16(test.addr, test.val)
		})
		if err != nil {
			if !test.panics {
				t.Errorf("&Segment{data: % x}.writeUint16(%v, %#x) unexpected panic: %v", test.data, test.addr, test.val, err)
			}
			continue
		}
		if test.panics {
			t.Errorf("&Segment{data: % x}.writeUint16(%v, %#x) did not panic as expected", test.data, test.addr, test.val)
			continue
		}
		if !bytes.Equal(out, test.out) {
			t.Errorf("data after &Segment{data: % x}.writeUint16(%v, %#x) = % x; want % x", test.data, test.addr, test.val, out, test.out)
		}
	}
}

func TestSegmentWriteUint32(t *testing.T) {
	tests := []struct {
		data   []byte
		addr   Address
		val    uint32
		out    []byte
		panics bool
	}{
		{
			data:   []byte{},
			addr:   0,
			val:    0,
			panics: true,
		},
		{
			data: []byte{1, 2, 3, 4, 5, 6},
			addr: 1,
			val:  0x01234567,
			out:  []byte{1, 0x67, 0x45, 0x23, 0x01, 6},
		},
	}
	for _, test := range tests {
		out := make([]byte, len(test.data))
		copy(out, test.data)
		seg := &Segment{data: out}
		err := catchPanic(func() {
			seg.writeUint32(test.addr, test.val)
		})
		if err != nil {
			if !test.panics {
				t.Errorf("&Segment{data: % x}.writeUint32(%v, %#x) unexpected panic: %v", test.data, test.addr, test.val, err)
			}
			continue
		}
		if test.panics {
			t.Errorf("&Segment{data: % x}.writeUint32(%v, %#x) did not panic as expected", test.data, test.addr, test.val)
			continue
		}
		if !bytes.Equal(out, test.out) {
			t.Errorf("data after &Segment{data: % x}.writeUint32(%v, %#x) = % x; want % x", test.data, test.addr, test.val, out, test.out)
		}
	}
}

func TestSegmentWriteUint64(t *testing.T) {
	tests := []struct {
		data   []byte
		addr   Address
		val    uint64
		out    []byte
		panics bool
	}{
		{
			data:   []byte{},
			addr:   0,
			val:    0,
			panics: true,
		},
		{
			data: []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			addr: 1,
			val:  0x0123456789abcdef,
			out:  []byte{1, 0xef, 0xcd, 0xab, 0x89, 0x67, 0x45, 0x23, 0x01, 10},
		},
	}
	for _, test := range tests {
		out := make([]byte, len(test.data))
		copy(out, test.data)
		seg := &Segment{data: out}
		err := catchPanic(func() {
			seg.writeUint64(test.addr, test.val)
		})
		if err != nil {
			if !test.panics {
				t.Errorf("&Segment{data: % x}.writeUint64(%v, %#x) unexpected panic: %v", test.data, test.addr, test.val, err)
			}
			continue
		}
		if test.panics {
			t.Errorf("&Segment{data: % x}.writeUint64(%v, %#x) did not panic as expected", test.data, test.addr, test.val)
			continue
		}
		if !bytes.Equal(out, test.out) {
			t.Errorf("data after &Segment{data: % x}.writeUint64(%v, %#x) = % x; want % x", test.data, test.addr, test.val, out, test.out)
		}
	}
}

func TestMakeOffsetKey(t *testing.T) {
	seg42 := &Segment{id: 42}
	tests := []struct {
		p          Ptr
		id         SegmentID
		boff, bend int64
	}{
		{
			p: Struct{
				seg:  seg42,
				off:  0,
				size: ObjectSize{0, 0},
			}.ToPtr(),
			id:   42,
			boff: 0,
			bend: 0,
		},
		{
			p: Struct{
				seg:  seg42,
				off:  8,
				size: ObjectSize{0, 0},
			}.ToPtr(),
			id:   42,
			boff: 64,
			bend: 64,
		},
		{
			p: Struct{
				seg:  seg42,
				off:  8,
				size: ObjectSize{1, 0},
			}.ToPtr(),
			id:   42,
			boff: 64,
			bend: 72,
		},
		{
			p: Struct{
				seg:  seg42,
				off:  8,
				size: ObjectSize{0, 1},
			}.ToPtr(),
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
			}.ToPtr(),
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
			}.ToPtr(),
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
			}.ToPtr(),
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
			}.ToPtr(),
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
			}.ToPtr(),
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
			}.ToPtr(),
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

func catchPanic(f func()) (err error) {
	defer func() {
		pval := recover()
		if pval == nil {
			return
		}
		e, ok := pval.(error)
		if !ok {
			err = fmt.Errorf("non-error panic: %#v", pval)
			return
		}
		err = e
	}()
	f()
	return nil
}

func TestCompare(t *testing.T) {
	// Offsets are in ascending order.
	data := []offset{
		{id: 0, boff: 10},
		{id: 0, boff: 20},
		{id: 0, boff: 30},
		{id: 0, boff: 65535},
		{id: 0, boff: 65536},
		{id: 1, boff: 0},
		{id: 1, boff: 5},
		{id: 1, boff: 65536},
	}
	formatOffset := func(o offset) string {
		return fmt.Sprintf("{id: %d, boff: %d}", o.id, o.boff)
	}

	for i, curr := range data {
		for _, prev := range data[:i] {
			if v := compare(curr, prev); v <= 0 {
				t.Errorf("compare(%s, %s) = %d; want >0", formatOffset(curr), formatOffset(prev), v)
			}
		}
		if v := compare(curr, curr); v != 0 {
			t.Errorf("compare(%s, %s) = %d; want 0", formatOffset(curr), formatOffset(curr), v)
		}
		for _, next := range data[i+1:] {
			if v := compare(curr, next); v >= 0 {
				t.Errorf("compare(%s, %s) = %d; want <0", formatOffset(curr), formatOffset(next), v)
			}
		}
	}
}
