package capnp

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToListDefault(t *testing.T) {
	msg := &Message{Arena: SingleSegment([]byte{
		0, 0, 0, 0, 0, 0, 0, 0,
		42, 0, 0, 0, 0, 0, 0, 0,
	})}
	seg, err := msg.Segment(0)
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		ptr  Ptr
		def  []byte
		list List
	}{
		{Ptr{}, nil, List{}},
		{Struct{}.ToPtr(), nil, List{}},
		{Struct{seg: seg, off: 0, depthLimit: maxDepth}.ToPtr(), nil, List{}},
		{List{}.ToPtr(), nil, List{}},
		{
			ptr: List{
				seg:        seg,
				off:        8,
				length:     1,
				size:       ObjectSize{DataSize: 8},
				depthLimit: maxDepth,
			}.ToPtr(),
			list: List{
				seg:        seg,
				off:        8,
				length:     1,
				size:       ObjectSize{DataSize: 8},
				depthLimit: maxDepth,
			},
		},
	}

	for _, test := range tests {
		list, err := test.ptr.ListDefault(test.def)
		if err != nil {
			t.Errorf("%#v.ListDefault(% 02x) error: %v", test.ptr, test.def, err)
			continue
		}
		if !deepPointerEqual(list.ToPtr(), test.list.ToPtr()) {
			t.Errorf("%#v.ListDefault(% 02x) = %#v; want %#v", test.ptr, test.def, list, test.list)
		}
	}
}

func TestTextListBytesAt(t *testing.T) {
	msg := &Message{Arena: SingleSegment([]byte{
		0, 0, 0, 0, 0, 0, 0, 0,
		0x01, 0, 0, 0, 0x22, 0, 0, 0,
		'f', 'o', 'o', 0, 0, 0, 0, 0,
	})}
	seg, err := msg.Segment(0)
	if err != nil {
		t.Fatal(err)
	}
	list := TextList{
		seg:        seg,
		off:        8,
		length:     1,
		size:       ObjectSize{PointerCount: 1},
		depthLimit: maxDepth,
	}
	b, err := list.BytesAt(0)
	if err != nil {
		t.Errorf("list.BytesAt(0) error: %v", err)
	}
	if !bytes.Equal(b, []byte("foo")) {
		t.Errorf("list.BytesAt(0) = %q; want \"foo\"", b)
	}
}

func TestListRaw(t *testing.T) {
	_, seg, err := NewMessage(SingleSegment(nil))
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		list List
		raw  rawPointer
	}{
		{list: List{}, raw: 0},
		{
			list: List{seg: seg, length: 3, size: ObjectSize{}},
			raw:  0x0000001800000001,
		},
		{
			list: List{seg: seg, off: 24, length: 15, flags: isBitList},
			raw:  0x0000007900000001,
		},
		{
			list: List{seg: seg, off: 40, length: 15, size: ObjectSize{DataSize: 1}},
			raw:  0x0000007a00000001,
		},
		{
			list: List{seg: seg, off: 40, length: 15, size: ObjectSize{DataSize: 2}},
			raw:  0x0000007b00000001,
		},
		{
			list: List{seg: seg, off: 40, length: 15, size: ObjectSize{DataSize: 4}},
			raw:  0x0000007c00000001,
		},
		{
			list: List{seg: seg, off: 40, length: 15, size: ObjectSize{DataSize: 8}},
			raw:  0x0000007d00000001,
		},
		{
			list: List{seg: seg, off: 40, length: 15, size: ObjectSize{PointerCount: 1}},
			raw:  0x0000007e00000001,
		},
		{
			list: List{seg: seg, off: 40, length: 7, size: ObjectSize{DataSize: 16, PointerCount: 1}, flags: isCompositeList},
			raw:  0x000000af00000001,
		},
	}
	for _, test := range tests {
		if raw := test.list.raw(); raw != test.raw {
			t.Errorf("%+v.raw() = %#v; want %#v", test.list, raw, test.raw)
		}
	}
}

// TestListCastRegression is a regression test for a bug where, if a struct
// list whose elements have a non-empty data section was cast to a pointer
// list, the pointer would be read out of the data section of the relevant
// element, rather than the pointer section.
func TestListCastRegression(t *testing.T) {
	_, seg, err := NewMessage(SingleSegment(nil))
	assert.Nil(t, err)

	txt, err := NewText(seg, "Text")
	assert.Nil(t, err)

	i := 3

	l, err := NewCompositeList(seg, ObjectSize{DataSize: 16, PointerCount: 1}, 6)
	assert.Nil(t, err)
	strct := l.Struct(i)
	strct.SetPtr(0, txt.ToPtr())

	ptrList := PointerList(l)
	ptr, err := ptrList.At(i)
	assert.Nil(t, err)
	assert.Equal(t, ptr.Text(), "Text")
}
