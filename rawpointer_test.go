package capnp

import (
	"testing"
)

func TestRawStructPointer(t *testing.T) {
	tests := []struct {
		ptr    rawPointer
		offset pointerOffset
		size   ObjectSize
	}{
		{0x0000000000000000, 0, ObjectSize{}},
		{0x0000000000000004, 1, ObjectSize{}},
		{0x0000000000000008, 2, ObjectSize{}},
		{0x000000000000000c, 3, ObjectSize{}},
		{0x0000000000000010, 4, ObjectSize{}},
		{0x0403020100000000, 0, ObjectSize{DataSize: 0x0201 * 8, PointerCount: 0x0403}},
		{0xffffffff00000000, 0, ObjectSize{DataSize: 0xffff * 8, PointerCount: 0xffff}},
		{0xfffffffffffffff8, -2, ObjectSize{DataSize: 0xffff * 8, PointerCount: 0xffff}},
		{0x00000000fffffffc, -1, ObjectSize{}},
		{0x04030201fffffffc, -1, ObjectSize{DataSize: 0x0201 * 8, PointerCount: 0x0403}},
		{0xfffffffffffffffc, -1, ObjectSize{DataSize: 0xffff * 8, PointerCount: 0xffff}},
	}
	for _, test := range tests {
		if typ := test.ptr.pointerType(); typ != structPointer {
			t.Errorf("rawPointer(%#016x).pointerType() = %d; want %d", uint64(test.ptr), typ, structPointer)
		}
		if offset := test.ptr.offset(); offset != test.offset {
			t.Errorf("rawPointer(%#016x).offset() = %d; want %d", uint64(test.ptr), offset, test.offset)
		}
		if size := test.ptr.structSize(); size != test.size {
			t.Errorf("rawPointer(%#016x).structSize() = %d; want %d", uint64(test.ptr), size, test.size)
		}
	}
}

func TestRawListPointer(t *testing.T) {
	tests := []struct {
		ptr    rawPointer
		offset pointerOffset
		lt     int
		n      int32
	}{
		{0x0000000000000001, 0, voidList, 0},
		{0x0000000000000005, 1, voidList, 0},
		{0x0000000000000009, 2, voidList, 0},
		{0x000000000000000d, 3, voidList, 0},
		{0x0000000100000001, 0, bit1List, 0},
		{0x0000000200000001, 0, byte1List, 0},
		{0x0000000300000001, 0, byte2List, 0},
		{0x0000000400000001, 0, byte4List, 0},
		{0x0000000500000001, 0, byte8List, 0},
		{0x0000000600000001, 0, pointerList, 0},
		{0x0000000700000001, 0, compositeList, 0},
		{0xfffffff8fffffffd, -1, voidList, 0x1fffffff},
		{0xfffffff9fffffffd, -1, bit1List, 0x1fffffff},
		{0xfffffffafffffffd, -1, byte1List, 0x1fffffff},
		{0xfffffffbfffffffd, -1, byte2List, 0x1fffffff},
		{0xfffffffcfffffffd, -1, byte4List, 0x1fffffff},
		{0xfffffffdfffffffd, -1, byte8List, 0x1fffffff},
		{0xfffffffefffffffd, -1, pointerList, 0x1fffffff},
		{0xfffffffffffffff9, -2, compositeList, 0x1fffffff},
		{0xfffffffffffffffd, -1, compositeList, 0x1fffffff},
	}
	for _, test := range tests {
		if typ := test.ptr.pointerType(); typ != listPointer {
			t.Errorf("rawPointer(%#016x).pointerType() = %d; want %d", uint64(test.ptr), typ, listPointer)
		}
		if offset := test.ptr.offset(); offset != test.offset {
			t.Errorf("rawPointer(%#016x).offset() = %d; want %d", uint64(test.ptr), offset, test.offset)
		}
		if lt := test.ptr.listType(); lt != test.lt {
			t.Errorf("rawPointer(%#016x).listType() = %d; want %d", uint64(test.ptr), lt, test.lt)
		}
		if n := test.ptr.numListElements(); n != test.n {
			t.Errorf("rawPointer(%#016x).numListElements() = %d; want %d", uint64(test.ptr), n, test.n)
		}
	}
}

func TestRawOtherPointer(t *testing.T) {
	tests := []struct {
		ptr rawPointer
		typ uint32
		cap CapabilityID
	}{
		{0x0000000000000003, 0, 0},
		{0x0000000000000007, 1, 0},
		{0x000000000000000b, 2, 0},
		{0x000000000000000f, 3, 0},
		{0xfffffffffffffffb, 0x3ffffffe, 0xffffffff},
		{0xffffffffffffffff, 0x3fffffff, 0xffffffff},
	}
	for _, test := range tests {
		if typ := test.ptr.pointerType(); typ != otherPointer {
			t.Errorf("rawPointer(%#016x).pointerType() = %d; want %d", uint64(test.ptr), typ, otherPointer)
		}
		if typ := test.ptr.otherPointerType(); typ != test.typ {
			t.Errorf("rawPointer(%#016x).otherPointerType() = %d; want %d", uint64(test.ptr), typ, test.typ)
		}
		if cap := test.ptr.capabilityIndex(); cap != test.cap {
			t.Errorf("rawPointer(%#016x).capabilityIndex() = %d; want %d", uint64(test.ptr), cap, test.cap)
		}
	}
}

func TestRawFarPointer(t *testing.T) {
	tests := []struct {
		ptr  rawPointer
		typ  int
		addr Address
		seg  SegmentID
	}{
		{0x0000000000000002, farPointer, 0, 0},
		{0x0000000000000006, doubleFarPointer, 0, 0},
		{0x000000000000000a, farPointer, 8, 0},
		{0x000000000000000e, doubleFarPointer, 8, 0},
		{0xfffffffffffffffa, farPointer, 0xfffffff8, 0xffffffff},
		{0xfffffffffffffffe, doubleFarPointer, 0xfffffff8, 0xffffffff},
	}
	for _, test := range tests {
		if typ := test.ptr.pointerType(); typ != test.typ {
			t.Errorf("rawPointer(%#016x).pointerType() = %d; want %d", uint64(test.ptr), typ, test.typ)
		}
		if addr := test.ptr.farAddress(); addr != test.addr {
			t.Errorf("rawPointer(%#016x).farAddress() = %v; want %v", uint64(test.ptr), addr, test.addr)
		}
		if seg := test.ptr.farSegment(); seg != test.seg {
			t.Errorf("rawPointer(%#016x).farSegment() = %d; want %d", uint64(test.ptr), seg, test.seg)
		}
	}
}
