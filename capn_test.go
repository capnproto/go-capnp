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
	if s1.Segment() == pl0.Segment() && s1.Address() == pl0.Address() {
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
	if s1s0.Segment() == sub.Segment() && s1s0.Address() == sub.Address() {
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
	addr := root.Address()
	end, _ := addr.addSize(wordSize)
	ptrSlice := seg.Data()[addr:end]
	want := []byte{0xfc, 0xff, 0xff, 0xff, 0, 0, 0, 0}
	if !bytes.Equal(ptrSlice, want) {
		t.Errorf("SetPtr wrote % 02x; want % 02x", ptrSlice, want)
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
