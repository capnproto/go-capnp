package capnp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"testing"
)

func TestSegmentInBounds(t *testing.T) {
	tests := []struct {
		n    int
		addr address
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
		addr address
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
		addr   address
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
		addr   address
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
		addr   address
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
		addr   address
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
		addr   address
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
		addr   address
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
		addr   address
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
		addr   address
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

func TestSetPtrCopyListMember(t *testing.T) {
	_, seg, err := NewMessage(SingleSegment(nil))
	if err != nil {
		t.Fatal("NewMessage:", err)
	}
	root, err := NewRootStruct(seg, ObjectSize{PointerCount: 2})
	if err != nil {
		t.Fatal("NewRootStruct:", err)
	}
	plist, err := NewCompositeList(seg, ObjectSize{PointerCount: 1}, 1)
	if err != nil {
		t.Fatal("NewCompositeList:", err)
	}
	if err := root.SetPtr(0, plist.ToPtr()); err != nil {
		t.Fatal("root.SetPtr(0, plist):", err)
	}
	sub, err := NewStruct(seg, ObjectSize{DataSize: 8})
	if err != nil {
		t.Fatal("NewStruct:", err)
	}
	sub.SetUint64(0, 42)
	pl0 := plist.Struct(0)
	if err := pl0.SetPtr(0, sub.ToPtr()); err != nil {
		t.Fatal("pl0.SetPtr(0, sub.ToPtr()):", err)
	}

	if err := root.SetPtr(1, pl0.ToPtr()); err != nil {
		t.Error("root.SetPtr(1, pl0):", err)
	}

	p1, err := root.Ptr(1)
	if err != nil {
		t.Error("root.Ptr(1):", err)
	}
	s1 := p1.Struct()
	if !s1.IsValid() {
		t.Error("root.Ptr(1) is not a valid struct")
	}
	if SamePtr(s1.ToPtr(), pl0.ToPtr()) {
		t.Error("list member not copied; points to same object")
	}
	s1p0, err := s1.Ptr(0)
	if err != nil {
		t.Error("root.Ptr(1).Struct().Ptr(0):", err)
	}
	s1s0 := s1p0.Struct()
	if !s1s0.IsValid() {
		t.Error("root.Ptr(1).Struct().Ptr(0) is not a valid struct")
	}
	if SamePtr(s1s0.ToPtr(), sub.ToPtr()) {
		t.Error("sub-object not copied; points to same object")
	}
	if got := s1s0.Uint64(0); got != 42 {
		t.Errorf("sub-object data = %d; want 42", got)
	}
}

func TestSetPtrToZeroSizeStruct(t *testing.T) {
	_, seg, err := NewMessage(SingleSegment(nil))
	if err != nil {
		t.Fatal("NewMessage:", err)
	}
	root, err := NewRootStruct(seg, ObjectSize{PointerCount: 1})
	if err != nil {
		t.Fatal("NewRootStruct:", err)
	}
	sub, err := NewStruct(seg, ObjectSize{})
	if err != nil {
		t.Fatal("NewStruct:", err)
	}
	if err := root.SetPtr(0, sub.ToPtr()); err != nil {
		t.Fatal("root.SetPtr(0, sub.ToPtr()):", err)
	}
	ptrSlice := seg.Data()[root.off : root.off+8]
	want := []byte{0xfc, 0xff, 0xff, 0xff, 0, 0, 0, 0}
	if !bytes.Equal(ptrSlice, want) {
		t.Errorf("SetPtr wrote % 02x; want % 02x", ptrSlice, want)
	}
}

func TestReadFarPointers(t *testing.T) {
	msg := &Message{
		// an rpc.capnp Message
		Arena: MultiSegment([][]byte{
			// Segment 0
			{
				// Double-far pointer: segment 2, offset 0
				0x06, 0x00, 0x00, 0x00, 0x02, 0x00, 0x00, 0x00,
			},
			// Segment 1
			{
				// (Root) Struct data section
				0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				// Struct pointer section
				// Double-far pointer: segment 4, offset 0
				0x06, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00,
			},
			// Segment 2
			{
				// Far pointer landing pad: segment 1, offset 0
				0x02, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
				// Far pointer landing pad tag word: struct with 1 word data and 1 pointer
				0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x01, 0x00,
			},
			// Segment 3
			{
				// (Root>0) Struct data section
				0x00, 0x00, 0x00, 0x00, 0x09, 0x00, 0x00, 0x00,
				0xaa, 0x70, 0x65, 0x21, 0xd7, 0x7b, 0x31, 0xa7,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				// Struct pointer section
				// Far pointer: segment 4, offset 4
				0x22, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00,
				// Far pointer: segment 4, offset 7
				0x3a, 0x00, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00,
				// Null
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			// Segment 4
			{
				// Far pointer landing pad: segment 3, offset 0
				0x02, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00,
				// Far pointer landing pad tag word: struct with 3 word data and 3 pointer
				0x00, 0x00, 0x00, 0x00, 0x03, 0x00, 0x03, 0x00,
				// (Root>0>0) Struct data section
				0x54, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				// Struct pointer section
				// Null
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				// Far pointer landing pad: struct pointer: offset -3, 1 word data, 1 pointer
				0xf4, 0xff, 0xff, 0xff, 0x01, 0x00, 0x01, 0x00,
				// (Root>0>1) Struct pointer section
				// Struct pointer: offset 2, 1 word data
				0x08, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00,
				// Null
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				// Far pointer landing pad: struct pointer: offset -3, 2 pointers
				0xf4, 0xff, 0xff, 0xff, 0x00, 0x00, 0x02, 0x00,
				// (Root>0>1>0) Struct data section
				0x2a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		}),
	}
	rootp, err := msg.Root()
	if err != nil {
		t.Error("Root:", err)
	}
	root := rootp.Struct()
	if root.Uint16(0) != 2 {
		t.Errorf("root.Uint16(0) = %d; want 2", root.Uint16(0))
	}
	callp, err := root.Ptr(0)
	if err != nil {
		t.Error("root.Ptr(0):", err)
	}
	call := callp.Struct()
	targetp, err := call.Ptr(0)
	if err != nil {
		t.Error("root.Ptr(0).Ptr(0):", err)
	}
	if got := targetp.Struct().Uint32(0); got != 84 {
		t.Errorf("root.Ptr(0).Ptr(0).Uint32(0) = %d; want 84", got)
	}
	paramsp, err := call.Ptr(1)
	if err != nil {
		t.Error("root.Ptr(0).Ptr(1):", err)
	}
	contentp, err := paramsp.Struct().Ptr(0)
	if err != nil {
		t.Error("root.Ptr(0).Ptr(1).Ptr(0):", err)
	}
	if got := contentp.Struct().Uint64(0); got != 42 {
		t.Errorf("root.Ptr(0).Ptr(1).Ptr(0).Uint64(0) = %d; want 42", got)
	}
}

func TestWriteFarPointer(t *testing.T) {
	// TODO(someday): run same test with a two-word list

	msg := &Message{
		Arena: MultiSegment([][]byte{
			make([]byte, 8),
			make([]byte, 0, 24),
		}),
	}
	seg1, err := msg.Segment(1)
	if err != nil {
		t.Fatal("msg.Segment(1):", err)
	}
	s, err := NewStruct(seg1, ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		t.Fatal("NewStruct(msg.Segment(1), ObjectSize{8, 1}):", err)
	}
	if s.Segment() != seg1 {
		t.Fatalf("struct allocated in segment %d", s.Segment().ID())
	}
	if err := msg.SetRoot(s.ToPtr()); err != nil {
		t.Error("msg.SetRoot(...):", err)
	}
	seg0, err := msg.Segment(0)
	if err != nil {
		t.Fatal("msg.Segment(0):", err)
	}

	root := rawPointer(binary.LittleEndian.Uint64(seg0.Data()))
	if root.pointerType() != farPointer {
		t.Fatalf("root (%#016x) type = %v; want %v (farPointer)", root, root.pointerType(), farPointer)
	}
	if root.farSegment() != 1 {
		t.Fatalf("root points to segment %d; want 1", root.farSegment())
	}
	padAddr := root.farAddress()
	if padAddr > address(len(seg1.Data())-8) {
		t.Fatalf("root points to out of bounds address %v; size of segment is %d", padAddr, len(seg1.Data()))
	}

	pad := rawPointer(binary.LittleEndian.Uint64(seg1.Data()[padAddr:]))
	if pad.pointerType() != structPointer {
		t.Errorf("landing pad (%#016x) type = %v; want %v (structPointer)", pad, pad.pointerType(), structPointer)
	}
	if got, ok := pad.offset().resolve(padAddr + 8); !ok || got != s.off {
		t.Errorf("landing pad (%#016x @ %v) resolved address = %v, %t; want %v, true", pad, padAddr, got, ok, s.off)
	}
	if got, want := pad.structSize(), (ObjectSize{DataSize: 8, PointerCount: 1}); got != want {
		t.Errorf("landing pad (%#016x) struct size = %v; want %v", pad, got, want)
	}
}

func TestWriteDoubleFarPointer(t *testing.T) {
	// TODO(someday): run same test with a two-word list

	msg := &Message{
		Arena: MultiSegment([][]byte{
			make([]byte, 8),
			make([]byte, 0, 16),
		}),
	}
	seg1, err := msg.Segment(1)
	if err != nil {
		t.Fatal("msg.Segment(1):", err)
	}
	s, err := NewStruct(seg1, ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		t.Fatal("NewStruct(msg.Segment(1), ObjectSize{8, 1}):", err)
	}
	if s.Segment() != seg1 {
		t.Fatalf("struct allocated in segment %d", s.Segment().ID())
	}
	if err := msg.SetRoot(s.ToPtr()); err != nil {
		t.Error("msg.SetRoot(...):", err)
	}
	seg0, err := msg.Segment(0)
	if err != nil {
		t.Fatal("msg.Segment(0):", err)
	}

	root := rawPointer(binary.LittleEndian.Uint64(seg0.Data()))
	if root.pointerType() != doubleFarPointer {
		t.Fatalf("root (%#016x) type = %v; want %v (doubleFarPointer)", root, root.pointerType(), doubleFarPointer)
	}
	if root.farSegment() == 0 || root.farSegment() == 1 {
		t.Fatalf("root points to segment %d; want !=0,1", root.farSegment())
	}
	padSeg, err := msg.Segment(root.farSegment())
	if err != nil {
		t.Fatalf("msg.Segment(%d): %v", root.farSegment(), err)
	}
	padAddr := root.farAddress()
	if padAddr > address(len(padSeg.Data())-16) {
		t.Fatalf("root points to out of bounds address %v; size of segment is %d", padAddr, len(padSeg.Data()))
	}

	pad1 := rawPointer(binary.LittleEndian.Uint64(padSeg.Data()[padAddr:]))
	if pad1.pointerType() != farPointer {
		t.Errorf("landing pad pointer 1 (%#016x) type = %v; want %v (farPointer)", pad1, pad1.pointerType(), farPointer)
	}
	if pad1.farSegment() != 1 {
		t.Fatalf("landing pad pointer 1 (%#016x) points to segment %d; want 1", pad1, pad1.farSegment())
	}
	if pad1.farAddress() != s.off {
		t.Fatalf("landing pad pointer 1 (%#016x) points to address %v; want %v", pad1, pad1.farAddress(), s.off)
	}

	pad2 := rawPointer(binary.LittleEndian.Uint64(padSeg.Data()[padAddr+8:]))
	if pad2.pointerType() != structPointer {
		t.Errorf("landing pad pointer 2 (%#016x) type = %v; want %v (structPointer)", pad2, pad2.pointerType(), structPointer)
	}
	if pad2.offset() != 0 {
		t.Errorf("landing pad pointer 2 (%#016x) offset = %d; want 0", pad2, pad2.offset())
	}
	if got, want := pad2.structSize(), (ObjectSize{DataSize: 8, PointerCount: 1}); got != want {
		t.Errorf("landing pad pointer 2 (%#016x) struct size = %v; want %v", pad2, got, want)
	}
}

func TestSetInterfacePtr(t *testing.T) {
	t.Run("SameMessage", func(t *testing.T) {
		msg, seg, err := NewMessage(SingleSegment(nil))
		if err != nil {
			t.Fatal("NewMessage:", err)
		}
		msg.AddCap(nil) // just to make the capability ID below non-zero
		root, err := NewRootStruct(seg, ObjectSize{PointerCount: 2})
		if err != nil {
			t.Fatal("NewRootStruct:", err)
		}
		hook := new(dummyHook)
		client := NewClient(hook)
		id := msg.AddCap(client)
		iface := NewInterface(seg, id)
		defer func() {
			for _, c := range msg.CapTable {
				c.Release()
			}
			if hook.shutdowns == 0 {
				t.Error("client leaked")
			}
		}()

		if err := root.SetPtr(0, iface.ToPtr()); err != nil {
			t.Fatal("root.SetPtr(0, iface.ToPtr()):", err)
		}
		ptr0, err := root.Ptr(0)
		if err != nil {
			t.Fatal("root.Ptr(0):", err)
		}
		iface0 := ptr0.Interface()
		if !iface0.IsValid() {
			t.Error("root.Ptr(0) is not an interface")
		} else if iface0.Capability() != id {
			t.Errorf("root.Ptr(0).Interface().Capability() = %d; want %d", iface0.Capability(), id)
		}

		if err := root.SetPtr(1, iface.ToPtr()); err != nil {
			t.Fatal("root.SetPtr(1, iface.ToPtr()):", err)
		}
		ptr1, err := root.Ptr(1)
		if err != nil {
			t.Fatal("root.Ptr(1):", err)
		}
		iface1 := ptr1.Interface()
		if !iface1.IsValid() {
			t.Error("root.Ptr(1) is not an interface")
		} else if iface1.Capability() != id {
			t.Errorf("root.Ptr(1).Interface().Capability() = %d; want %d", iface1.Capability(), id)
		}
	})
	t.Run("DifferentMessages", func(t *testing.T) {
		msg1 := &Message{Arena: SingleSegment(nil)}
		seg1, err := msg1.Segment(0)
		if err != nil {
			t.Fatal("msg1.Segment(0):", err)
		}
		msg2, seg2, err := NewMessage(SingleSegment(nil))
		if err != nil {
			t.Fatal("NewMessage:", err)
		}
		root, err := NewRootStruct(seg2, ObjectSize{PointerCount: 1})
		if err != nil {
			t.Fatal("NewRootStruct:", err)
		}
		defer func() {
			for _, c := range msg1.CapTable {
				c.Release()
			}
			for _, c := range msg2.CapTable {
				c.Release()
			}
		}()

		hook := new(dummyHook)
		client := NewClient(hook)
		iface1 := NewInterface(seg1, msg1.AddCap(client))
		if err := root.SetPtr(0, iface1.ToPtr()); err != nil {
			t.Fatal("root.SetPtr(0, iface1.ToPtr()):", err)
		}
		ptr, err := root.Ptr(0)
		if err != nil {
			t.Fatal("root.Ptr(0):", err)
		}
		iface2 := ptr.Interface()
		if !iface2.Client().IsSame(iface1.Client()) {
			t.Errorf("root.Ptr(0).Interface().Client() = %v; want %v", iface2.Client(), iface1.Client())
		}
		for _, c := range msg1.CapTable {
			c.Release()
		}
		msg1.CapTable = nil
		if hook.shutdowns > 0 {
			t.Error("copying interface across messages did not add reference")
		}
		for _, c := range msg2.CapTable {
			c.Release()
		}
		msg2.CapTable = nil
		if hook.shutdowns == 0 {
			t.Error("client not shut down after releasing both message capability tables")
		}
	})
}

func TestEqual(t *testing.T) {
	msg, seg, _ := NewMessage(SingleSegment(nil))
	emptyStruct1, _ := NewStruct(seg, ObjectSize{})
	emptyStruct2, _ := NewStruct(seg, ObjectSize{})
	zeroStruct1, _ := NewStruct(seg, ObjectSize{DataSize: 8, PointerCount: 1})
	zeroStruct2, _ := NewStruct(seg, ObjectSize{DataSize: 8, PointerCount: 1})
	structA1, _ := NewStruct(seg, ObjectSize{DataSize: 16, PointerCount: 1})
	structA1.SetUint32(0, 0xdeadbeef)
	subA1, _ := NewStruct(seg, ObjectSize{DataSize: 8})
	subA1.SetUint32(0, 0x0cafefe0)
	structA1.SetPtr(0, subA1.ToPtr())
	structA2, _ := NewStruct(seg, ObjectSize{DataSize: 16, PointerCount: 1})
	structA2.SetUint32(0, 0xdeadbeef)
	subA2, _ := NewStruct(seg, ObjectSize{DataSize: 8})
	subA2.SetUint32(0, 0x0cafefe0)
	structA2.SetPtr(0, subA2.ToPtr())
	structB, _ := NewStruct(seg, ObjectSize{DataSize: 8, PointerCount: 2})
	structB.SetUint32(0, 0xdeadbeef)
	subB, _ := NewStruct(seg, ObjectSize{DataSize: 8})
	subB.SetUint32(0, 0x0cafefe0)
	structB.SetPtr(0, subB.ToPtr())
	structB.SetPtr(1, emptyStruct1.ToPtr())
	structC, _ := NewStruct(seg, ObjectSize{DataSize: 16, PointerCount: 1})
	structC.SetUint32(0, 0xfeed1234)
	subC, _ := NewStruct(seg, ObjectSize{DataSize: 8})
	subC.SetUint32(0, 0x0cafefe0)
	structC.SetPtr(0, subA2.ToPtr())
	structD, _ := NewStruct(seg, ObjectSize{DataSize: 16, PointerCount: 1})
	structD.SetUint32(0, 0xdeadbeef)
	subD, _ := NewStruct(seg, ObjectSize{DataSize: 8})
	subD.SetUint32(0, 0x12345678)
	structD.SetPtr(0, subD.ToPtr())
	emptyStructList1, _ := NewCompositeList(seg, ObjectSize{DataSize: 8, PointerCount: 1}, 0)
	emptyStructList2, _ := NewCompositeList(seg, ObjectSize{DataSize: 8, PointerCount: 1}, 0)
	emptyInt32List, _ := NewInt32List(seg, 0)
	emptyFloat32List, _ := NewFloat32List(seg, 0)
	emptyFloat64List, _ := NewFloat64List(seg, 0)
	list123Int, _ := NewInt32List(seg, 3)
	list123Int.Set(0, 1)
	list123Int.Set(1, 2)
	list123Int.Set(2, 3)
	list12Int, _ := NewInt32List(seg, 2)
	list12Int.Set(0, 1)
	list12Int.Set(1, 2)
	list456Int, _ := NewInt32List(seg, 3)
	list456Int.Set(0, 4)
	list456Int.Set(1, 5)
	list456Int.Set(2, 6)
	list123Struct, _ := NewCompositeList(seg, ObjectSize{DataSize: 8}, 3)
	list123Struct.Struct(0).SetUint32(0, 1)
	list123Struct.Struct(1).SetUint32(0, 2)
	list123Struct.Struct(2).SetUint32(0, 3)
	list12Struct, _ := NewCompositeList(seg, ObjectSize{DataSize: 8}, 2)
	list12Struct.Struct(0).SetUint32(0, 1)
	list12Struct.Struct(1).SetUint32(0, 2)
	list456Struct, _ := NewCompositeList(seg, ObjectSize{DataSize: 8}, 3)
	list456Struct.Struct(0).SetUint32(0, 4)
	list456Struct.Struct(1).SetUint32(0, 5)
	list456Struct.Struct(2).SetUint32(0, 6)
	plistA1, _ := NewPointerList(seg, 1)
	plistA1.Set(0, structA1.ToPtr())
	plistA2, _ := NewPointerList(seg, 1)
	plistA2.Set(0, structA2.ToPtr())
	plistB, _ := NewPointerList(seg, 1)
	plistB.Set(0, structB.ToPtr())
	ec := ErrorClient(errors.New("boo"))
	msg.CapTable = []*Client{
		0: ec,
		1: ec,
		2: ErrorClient(errors.New("another boo")),
		3: nil,
		4: nil,
	}
	iface1 := NewInterface(seg, 0)
	iface2 := NewInterface(seg, 1)
	ifaceAlt := NewInterface(seg, 2)
	ifaceMissing1 := NewInterface(seg, 3)
	ifaceMissing2 := NewInterface(seg, 4)
	ifaceOOB1 := NewInterface(seg, 5)
	ifaceOOB2 := NewInterface(seg, 6)

	tests := []struct {
		name   string
		p1, p2 Ptr
		equal  bool
	}{

		// Structs
		{"EmptyStruct_EmptyStruct", emptyStruct1.ToPtr(), emptyStruct2.ToPtr(), true},
		{"EmptyStruct_ZeroStruct", emptyStruct1.ToPtr(), zeroStruct2.ToPtr(), true},
		{"ZeroStruct_ZeroStruct", zeroStruct1.ToPtr(), zeroStruct2.ToPtr(), true},
		{"EmptyStruct_StructA", emptyStruct1.ToPtr(), structA1.ToPtr(), false},
		{"StructA_EmptyStruct", structA1.ToPtr(), emptyStruct1.ToPtr(), false},
		{"StructA1_StructA1", structA1.ToPtr(), structA1.ToPtr(), true},
		{"StructA1_StructA2", structA1.ToPtr(), structA2.ToPtr(), true},
		{"StructA2_StructA1", structA2.ToPtr(), structA1.ToPtr(), true},
		{"StructA_StructB", structA1.ToPtr(), structB.ToPtr(), false},
		{"StructB_StructA", structB.ToPtr(), structA1.ToPtr(), false},
		{"StructA_StructC", structA1.ToPtr(), structC.ToPtr(), false},
		{"StructC_StructA", structC.ToPtr(), structA1.ToPtr(), false},
		{"StructA_StructD", structA1.ToPtr(), structD.ToPtr(), false},
		{"StructD_StructA", structD.ToPtr(), structA1.ToPtr(), false},

		// Lists
		{"EmptyStructList_EmptyStructList", emptyStructList1.ToPtr(), emptyStructList2.ToPtr(), true},
		{"EmptyInt32List_EmptyFloat32List", emptyInt32List.ToPtr(), emptyFloat32List.ToPtr(), true}, // identical on wire
		{"EmptyInt32List_EmptyFloat64List", emptyInt32List.ToPtr(), emptyFloat64List.ToPtr(), false},
		{"List123Int_List456Int", list123Int.ToPtr(), list456Int.ToPtr(), false},
		{"List123Struct_List456Struct", list123Struct.ToPtr(), list456Struct.ToPtr(), false},
		{"List123Int_List123Struct", list123Int.ToPtr(), list123Struct.ToPtr(), true},
		{"List123Struct_List123Int", list123Struct.ToPtr(), list123Int.ToPtr(), true},
		{"List123Int_List12Int", list123Int.ToPtr(), list12Int.ToPtr(), false},
		{"List123Struct_List12Struct", list123Struct.ToPtr(), list12Struct.ToPtr(), false},
		{"PointerListA1_PointerListA2", plistA1.ToPtr(), plistA2.ToPtr(), true},
		{"PointerListA2_PointerListA1", plistA2.ToPtr(), plistA1.ToPtr(), true},
		{"PointerListA_PointerListB", plistA1.ToPtr(), plistB.ToPtr(), false},
		{"PointerListB_PointerListA", plistB.ToPtr(), plistA1.ToPtr(), false},

		// Interfaces
		{"InterfaceA1_InterfaceA1", iface1.ToPtr(), iface1.ToPtr(), true},
		{"InterfaceA1_InterfaceA2", iface1.ToPtr(), iface2.ToPtr(), true},
		{"InterfaceA_InterfaceB", iface1.ToPtr(), ifaceAlt.ToPtr(), false},
		{"InterfaceMissingCap_Null", ifaceMissing1.ToPtr(), Ptr{}, false},
		{"InterfaceMissingCap_InterfaceMissingCap", ifaceMissing1.ToPtr(), ifaceMissing2.ToPtr(), true},
		{"InterfaceOOB1_InterfaceOOB1", ifaceOOB1.ToPtr(), ifaceOOB1.ToPtr(), true},
		{"InterfaceOOB1_InterfaceOOB2", ifaceOOB1.ToPtr(), ifaceOOB2.ToPtr(), false},
		{"InterfaceOOB_InterfaceMissingCap", ifaceOOB1.ToPtr(), ifaceMissing1.ToPtr(), false},

		// Null
		{"Null_Null", Ptr{}, Ptr{}, true},
		{"EmptyStruct_Null", emptyStruct1.ToPtr(), Ptr{}, false},
		{"Null_EmptyStruct", Ptr{}, emptyStruct1.ToPtr(), false},
		{"Null_EmptyStructList", Ptr{}, emptyStructList1.ToPtr(), false},
		{"EmptyStructList_Null", emptyStructList1.ToPtr(), Ptr{}, false},
		{"Interface_Null", iface1.ToPtr(), Ptr{}, false},
		{"Null_Interface", Ptr{}, iface1.ToPtr(), false},

		// Misc combinations that shouldn't be equal.
		{"EmptyStruct_EmptyList", emptyStruct1.ToPtr(), emptyStructList1.ToPtr(), false},
		{"EmptyStruct_InterfaceA", emptyStruct1.ToPtr(), iface1.ToPtr(), false},
		{"EmptyList_InterfaceA", emptyStructList1.ToPtr(), iface1.ToPtr(), false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Equal(test.p1, test.p2)
			if err != nil {
				t.Fatal(err)
			}
			if got != test.equal {
				if got {
					t.Error("p1 equals p2; want not equal")
				} else {
					t.Error("p1 does not equal p2; want equal")
				}
			}
		})
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
