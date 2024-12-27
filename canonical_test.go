package capnp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCanonicalize(t *testing.T) {

	tests := []struct {
		name string
		f    func() Struct
		want []byte
	}{{
		name: "Struct{}",
		f:    func() Struct { return Struct{} },
		want: []byte{0, 0, 0, 0, 0, 0, 0, 0},
	}, {
		name: "empty struct",
		f: func() Struct {
			_, seg := NewSingleSegmentMessage(nil)
			s, _ := NewStruct(seg, ObjectSize{})
			return s
		},
		want: []byte{0xfc, 0xff, 0xff, 0xff, 0, 0, 0, 0},
	}, {
		name: "zero data, zero pointer struct",
		f: func() Struct {
			_, seg := NewSingleSegmentMessage(nil)
			s, _ := NewStruct(seg, ObjectSize{DataSize: 8, PointerCount: 1})
			return s
		},
		want: []byte{0xfc, 0xff, 0xff, 0xff, 0, 0, 0, 0},
	}, {
		name: "one word data struct",
		f: func() Struct {
			_, seg := NewSingleSegmentMessage(nil)
			s, _ := NewStruct(seg, ObjectSize{DataSize: 8, PointerCount: 1})
			s.SetUint16(0, 0xbeef)
			return s
		},
		want: []byte{
			0, 0, 0, 0, 1, 0, 0, 0,
			0xef, 0xbe, 0, 0, 0, 0, 0, 0,
		},
	}, {
		name: "two pointers to zero structs",
		f: func() Struct {
			_, seg := NewSingleSegmentMessage(nil)
			s, _ := NewStruct(seg, ObjectSize{PointerCount: 2})
			e1, _ := NewStruct(seg, ObjectSize{DataSize: 8})
			e2, _ := NewStruct(seg, ObjectSize{DataSize: 8})
			s.SetPtr(0, e1.ToPtr())
			s.SetPtr(1, e2.ToPtr())
			return s
		},
		want: []byte{
			0, 0, 0, 0, 0, 0, 2, 0,
			0xfc, 0xff, 0xff, 0xff, 0, 0, 0, 0,
			0xfc, 0xff, 0xff, 0xff, 0, 0, 0, 0,
		},
	}, {
		name: "pointer to interface",
		f: func() Struct {
			_, seg := NewSingleSegmentMessage(nil)
			s, _ := NewStruct(seg, ObjectSize{PointerCount: 2})
			iface := NewInterface(seg, 1)
			s.SetPtr(0, iface.ToPtr())
			return s
		},
		want: []byte{
			0, 0, 0, 0, 0, 0, 1, 0,
			3, 0, 0, 0, 1, 0, 0, 0,
		},
	}, {
		name: "int list",
		f: func() Struct {
			_, seg := NewSingleSegmentMessage(nil)
			s, _ := NewStruct(seg, ObjectSize{PointerCount: 1})
			l, _ := NewInt8List(seg, 5)
			s.SetPtr(0, l.ToPtr())
			l.Set(0, 1)
			l.Set(1, 2)
			l.Set(2, 3)
			l.Set(3, 4)
			l.Set(4, 5)
			return s
		},
		want: []byte{
			0, 0, 0, 0, 0, 0, 1, 0,
			0x01, 0, 0, 0, 0x2a, 0, 0, 0,
			1, 2, 3, 4, 5, 0, 0, 0,
		},
	}, {
		name: "zero int list",
		f: func() Struct {
			_, seg := NewSingleSegmentMessage(nil)
			s, _ := NewStruct(seg, ObjectSize{PointerCount: 1})
			l, _ := NewInt8List(seg, 5)
			s.SetPtr(0, l.ToPtr())
			return s

		},
		want: []byte{
			0, 0, 0, 0, 0, 0, 1, 0,
			0x01, 0, 0, 0, 0x2a, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
		},
	}, {
		name: "struct list",
		f: func() Struct {
			_, seg := NewSingleSegmentMessage(nil)
			s, _ := NewStruct(seg, ObjectSize{PointerCount: 1})
			l, _ := NewCompositeList(seg, ObjectSize{DataSize: 8, PointerCount: 1}, 2)
			s.SetPtr(0, l.ToPtr())
			l.Struct(0).SetUint64(0, 0xdeadbeef)
			txt, _ := NewText(seg, "xyzzy")
			l.Struct(1).SetPtr(0, txt.ToPtr())
			return s

		},
		want: []byte{
			0, 0, 0, 0, 0, 0, 1, 0,
			0x01, 0, 0, 0, 0x27, 0, 0, 0,
			0x08, 0, 0, 0, 1, 0, 1, 0,
			0xef, 0xbe, 0xad, 0xde, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0x01, 0, 0, 0, 0x32, 0, 0, 0,
			'x', 'y', 'z', 'z', 'y', 0, 0, 0,
		},
	}, {
		name: "zero struct list",
		f: func() Struct {
			_, seg := NewSingleSegmentMessage(nil)
			s, _ := NewStruct(seg, ObjectSize{PointerCount: 1})
			l, _ := NewCompositeList(seg, ObjectSize{DataSize: 16, PointerCount: 2}, 3)
			s.SetPtr(0, l.ToPtr())
			return s

		},
		want: []byte{
			0, 0, 0, 0, 0, 0, 1, 0,
			0x01, 0, 0, 0, 0x07, 0, 0, 0,
			0x0c, 0, 0, 0, 0, 0, 0, 0,
		},
	}, {
		name: "zero-length struct list",
		f: func() Struct {
			_, seg := NewSingleSegmentMessage(nil)
			s, _ := NewStruct(seg, ObjectSize{PointerCount: 1})
			l, _ := NewCompositeList(seg, ObjectSize{DataSize: 16, PointerCount: 2}, 0)
			s.SetPtr(0, l.ToPtr())
			return s

		},
		want: []byte{
			0, 0, 0, 0, 0, 0, 1, 0,
			0x01, 0, 0, 0, 0x07, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
		},
	}}

	for i := range tests {
		tc := tests[i]
		t.Run(tc.name, func(t *testing.T) {
			b, err := Canonicalize(tc.f())
			require.NoError(t, err)
			require.Equal(t, tc.want, b)
		})
	}

}
