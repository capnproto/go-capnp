package pogs

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"

	"capnproto.org/go/capnp/v3"
	air "capnproto.org/go/capnp/v3/internal/aircraftlib"
	"github.com/kylelemons/godebug/pretty"
)

type Z struct {
	Which air.Z_Which

	F64 float64
	F32 float32

	I64 int64
	I32 int32
	I16 int16
	I8  int8

	U64 uint64
	U32 uint32
	U16 uint16
	U8  uint8

	Bool bool
	Text string
	Blob []byte

	F64vec []float64
	F32vec []float32

	I64vec []int64
	I8vec  []int8

	U64vec []uint64
	U8vec  []uint8

	Boolvec []bool
	Datavec [][]byte
	Textvec []string

	Zvec    []*Z
	Zvecvec [][]*Z

	Planebase *PlaneBase
	Airport   air.Airport

	Grp *ZGroup

	Echo   air.Echo
	Echoes []air.Echo

	AnyPtr        capnp.Ptr
	AnyStruct     capnp.Struct
	AnyList       capnp.List
	AnyCapability capnp.Client
}

type PlaneBase struct {
	Name     string
	Homes    []air.Airport
	Rating   int64
	CanFly   bool
	Capacity int64
	MaxSpeed float64
}

func (p *PlaneBase) equal(q *PlaneBase) bool {
	if p == nil && q == nil {
		return true
	}
	if (p == nil) != (q == nil) {
		return false
	}
	if len(p.Homes) != len(q.Homes) {
		return false
	}
	for i := range p.Homes {
		if p.Homes[i] != q.Homes[i] {
			return false
		}
	}
	return p.Name == q.Name &&
		p.Rating == q.Rating &&
		p.CanFly == q.CanFly &&
		p.Capacity == q.Capacity &&
		p.MaxSpeed == q.MaxSpeed
}

type ZGroup struct {
	First  uint64
	Second uint64
}

var goodTests = []Z{
	{Which: air.Z_Which_f64, F64: 3.5},
	{Which: air.Z_Which_f32, F32: 3.5},
	{Which: air.Z_Which_i64, I64: -123},
	{Which: air.Z_Which_i32, I32: -123},
	{Which: air.Z_Which_i16, I16: -123},
	{Which: air.Z_Which_i8, I8: -123},
	{Which: air.Z_Which_u64, U64: 123},
	{Which: air.Z_Which_u32, U32: 123},
	{Which: air.Z_Which_u16, U16: 123},
	{Which: air.Z_Which_u8, U8: 123},
	{Which: air.Z_Which_bool, Bool: true},
	{Which: air.Z_Which_bool, Bool: false},
	{Which: air.Z_Which_text, Text: "Hello, World!"},
	{Which: air.Z_Which_blob, Blob: nil},
	{Which: air.Z_Which_blob, Blob: []byte{}},
	{Which: air.Z_Which_blob, Blob: []byte("Hello, World!")},
	{Which: air.Z_Which_f64vec, F64vec: nil},
	{Which: air.Z_Which_f64vec, F64vec: []float64{-2.0, 4.5}},
	{Which: air.Z_Which_f32vec, F32vec: nil},
	{Which: air.Z_Which_f32vec, F32vec: []float32{-2.0, 4.5}},
	{Which: air.Z_Which_i64vec, I64vec: nil},
	{Which: air.Z_Which_i64vec, I64vec: []int64{-123, 0, 123}},
	{Which: air.Z_Which_i8vec, I8vec: nil},
	{Which: air.Z_Which_i8vec, I8vec: []int8{-123, 0, 123}},
	{Which: air.Z_Which_u64vec, U64vec: nil},
	{Which: air.Z_Which_u64vec, U64vec: []uint64{0, 123}},
	{Which: air.Z_Which_u8vec, U8vec: nil},
	{Which: air.Z_Which_u8vec, U8vec: []uint8{0, 123}},
	{Which: air.Z_Which_boolvec, Boolvec: nil},
	{Which: air.Z_Which_boolvec, Boolvec: []bool{false, true, false}},
	{Which: air.Z_Which_datavec, Datavec: nil},
	{Which: air.Z_Which_datavec, Datavec: [][]byte{[]byte("hi"), []byte("bye")}},
	{Which: air.Z_Which_datavec, Datavec: [][]byte{nil, nil, nil}},
	{Which: air.Z_Which_textvec, Textvec: nil},
	{Which: air.Z_Which_textvec, Textvec: []string{"John", "Paul", "George", "Ringo"}},
	{Which: air.Z_Which_textvec, Textvec: []string{"", "", ""}},
	{Which: air.Z_Which_zvec, Zvec: []*Z{
		{Which: air.Z_Which_i64, I64: -123},
		{Which: air.Z_Which_text, Text: "Hi"},
	}},
	{Which: air.Z_Which_zvecvec, Zvecvec: [][]*Z{
		{
			{Which: air.Z_Which_i64, I64: 1},
			{Which: air.Z_Which_i64, I64: 2},
		},
		{
			{Which: air.Z_Which_i64, I64: 3},
			{Which: air.Z_Which_i64, I64: 4},
		},
	}},
	{Which: air.Z_Which_planebase, Planebase: nil},
	{Which: air.Z_Which_planebase, Planebase: &PlaneBase{
		Name:     "Boeing",
		Homes:    []air.Airport{air.Airport_lax, air.Airport_dfw},
		Rating:   123,
		CanFly:   true,
		Capacity: 100,
		MaxSpeed: 9001.0,
	}},
	{Which: air.Z_Which_airport, Airport: air.Airport_lax},
	{Which: air.Z_Which_grp, Grp: &ZGroup{First: 123, Second: 456}},
	{Which: air.Z_Which_echo, Echo: air.Echo{}},
	{Which: air.Z_Which_echo, Echo: air.Echo(capnp.ErrorClient(errors.New("boo")))},
	{Which: air.Z_Which_echoes, Echoes: []air.Echo{
		{},
		air.Echo(capnp.ErrorClient(errors.New("boo"))),
		{},
		air.Echo(capnp.ErrorClient(errors.New("boo"))),
		{},
	}},
	{Which: air.Z_Which_anyPtr, AnyPtr: capnp.Ptr{}},
	{Which: air.Z_Which_anyPtr, AnyPtr: newTestStruct().ToPtr()},
	{Which: air.Z_Which_anyPtr, AnyPtr: newTestList().ToPtr()},
	{Which: air.Z_Which_anyPtr, AnyPtr: newTestInterface().ToPtr()},
	{Which: air.Z_Which_anyStruct, AnyStruct: capnp.Struct{}},
	{Which: air.Z_Which_anyStruct, AnyStruct: newTestStruct()},
	{Which: air.Z_Which_anyList, AnyList: capnp.List{}},
	{Which: air.Z_Which_anyList, AnyList: newTestList()},
	{Which: air.Z_Which_anyCapability, AnyCapability: capnp.Client{}},
	{Which: air.Z_Which_anyCapability, AnyCapability: capnp.ErrorClient(errors.New("boo"))},
}

func newTestStruct() capnp.Struct {
	_, seg, _ := capnp.NewMessage(capnp.SingleSegment(nil))
	s, _ := capnp.NewRootStruct(seg, capnp.ObjectSize{DataSize: 8})
	s.SetUint32(0, 0xdeadbeef)
	return s
}

func newTestList() capnp.List {
	_, seg, _ := capnp.NewMessage(capnp.SingleSegment(nil))
	l, _ := capnp.NewInt32List(seg, 3)
	l.Set(0, 123)
	l.Set(1, 456)
	l.Set(2, 789)
	return capnp.List(l)
}

func newTestInterface() capnp.Interface {
	msg, seg, _ := capnp.NewMessage(capnp.SingleSegment(nil))
	id := msg.CapTable().Add(capnp.ErrorClient(errors.New("boo")))
	return capnp.NewInterface(seg, id)
}

func TestExtract(t *testing.T) {
	for _, test := range goodTests {
		_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
		if err != nil {
			t.Errorf("NewMessage for %s: %v", zpretty.Sprint(test), err)
			continue
		}
		z, err := air.NewRootZ(seg)
		if err != nil {
			t.Errorf("NewRootZ for %s: %v", zpretty.Sprint(test), err)
			continue
		}
		if err := zfill(z, &test); err != nil {
			t.Errorf("zfill for %s: %v", zpretty.Sprint(test), err)
			continue
		}
		out := new(Z)
		if err := Extract(out, air.Z_TypeID, capnp.Struct(z)); err != nil {
			t.Errorf("Extract(%v) error: %v", z, err)
		}
		if !test.equal(out) {
			t.Errorf("Extract(%v) produced %s; want %s", z, zpretty.Sprint(out), zpretty.Sprint(test))
		}
	}
}

func TestInsert(t *testing.T) {
	for _, test := range goodTests {
		_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
		if err != nil {
			t.Errorf("NewMessage for %s: %v", zpretty.Sprint(test), err)
			continue
		}
		z, err := air.NewRootZ(seg)
		if err != nil {
			t.Errorf("NewRootZ for %s: %v", zpretty.Sprint(test), err)
			continue
		}
		err = Insert(air.Z_TypeID, capnp.Struct(z), &test)
		if err != nil {
			t.Errorf("Insert(%s) error: %v", zpretty.Sprint(test), err)
		}
		if equal, err := zequal(&test, z); err != nil {
			t.Errorf("Insert(%s) compare err: %v", zpretty.Sprint(test), err)
		} else if !equal {
			t.Errorf("Insert(%s) produced %v", zpretty.Sprint(test), z)
		}
	}
}

func TestInsert_Size(t *testing.T) {
	const baseSize = 8
	tests := []struct {
		name string
		sz   capnp.ObjectSize
		z    Z
		ok   bool
	}{
		{
			name: "void into empty",
			z:    Z{Which: air.Z_Which_void},
		},
		{
			name: "void into 0-byte",
			sz:   capnp.ObjectSize{DataSize: baseSize},
			z:    Z{Which: air.Z_Which_void},
			ok:   true,
		},
		{
			name: "void into 1-byte",
			sz:   capnp.ObjectSize{DataSize: baseSize + 1},
			z:    Z{Which: air.Z_Which_void},
			ok:   true,
		},
		{
			name: "bool into empty",
			z:    Z{Which: air.Z_Which_bool, Bool: true},
		},
		{
			name: "bool into 0 byte",
			sz:   capnp.ObjectSize{DataSize: baseSize},
			z:    Z{Which: air.Z_Which_bool, Bool: true},
		},
		{
			name: "bool into 1 byte",
			sz:   capnp.ObjectSize{DataSize: baseSize + 1},
			z:    Z{Which: air.Z_Which_bool, Bool: true},
			ok:   true,
		},
		{
			name: "bool into 0 byte, 1-pointer",
			sz:   capnp.ObjectSize{DataSize: baseSize, PointerCount: 1},
			z:    Z{Which: air.Z_Which_bool, Bool: true},
		},
		{
			name: "int8 into 0-byte",
			sz:   capnp.ObjectSize{DataSize: baseSize},
			z:    Z{Which: air.Z_Which_i8, I8: 123},
		},
		{
			name: "int8 into 1-byte",
			sz:   capnp.ObjectSize{DataSize: baseSize + 1},
			z:    Z{Which: air.Z_Which_i8, I8: 123},
			ok:   true,
		},
		{
			name: "uint8 into 0-byte",
			sz:   capnp.ObjectSize{DataSize: baseSize},
			z:    Z{Which: air.Z_Which_u8, U8: 123},
		},
		{
			name: "uint8 into 1-byte",
			sz:   capnp.ObjectSize{DataSize: baseSize + 1},
			z:    Z{Which: air.Z_Which_u8, U8: 123},
			ok:   true,
		},
		{
			name: "int16 into 0-byte",
			sz:   capnp.ObjectSize{DataSize: baseSize},
			z:    Z{Which: air.Z_Which_i16, I16: 123},
		},
		{
			name: "int16 into 2-byte",
			sz:   capnp.ObjectSize{DataSize: baseSize + 2},
			z:    Z{Which: air.Z_Which_i16, I16: 123},
			ok:   true,
		},
		{
			name: "uint16 into 0-byte",
			sz:   capnp.ObjectSize{DataSize: baseSize},
			z:    Z{Which: air.Z_Which_u16, U16: 123},
		},
		{
			name: "uint16 into 2-byte",
			sz:   capnp.ObjectSize{DataSize: baseSize + 2},
			z:    Z{Which: air.Z_Which_u16, U16: 123},
			ok:   true,
		},
		{
			name: "enum into 0-byte",
			sz:   capnp.ObjectSize{DataSize: baseSize},
			z:    Z{Which: air.Z_Which_airport, Airport: air.Airport_jfk},
		},
		{
			name: "enum into 2-byte",
			sz:   capnp.ObjectSize{DataSize: baseSize + 2},
			z:    Z{Which: air.Z_Which_airport, Airport: air.Airport_jfk},
			ok:   true,
		},
		{
			name: "int32 into 0-byte",
			sz:   capnp.ObjectSize{DataSize: baseSize},
			z:    Z{Which: air.Z_Which_i32, I32: 123},
		},
		{
			name: "int32 into 4-byte",
			sz:   capnp.ObjectSize{DataSize: baseSize + 4},
			z:    Z{Which: air.Z_Which_i32, I32: 123},
			ok:   true,
		},
		{
			name: "uint32 into 0-byte",
			sz:   capnp.ObjectSize{DataSize: baseSize},
			z:    Z{Which: air.Z_Which_u32, U32: 123},
		},
		{
			name: "uint32 into 4-byte",
			sz:   capnp.ObjectSize{DataSize: baseSize + 4},
			z:    Z{Which: air.Z_Which_u32, U32: 123},
			ok:   true,
		},
		{
			name: "float32 into 0-byte",
			sz:   capnp.ObjectSize{DataSize: baseSize},
			z:    Z{Which: air.Z_Which_f32, F32: 123},
		},
		{
			name: "float32 into 4-byte",
			sz:   capnp.ObjectSize{DataSize: baseSize + 4},
			z:    Z{Which: air.Z_Which_f32, F32: 123},
			ok:   true,
		},
		{
			name: "int64 into 0-byte",
			sz:   capnp.ObjectSize{DataSize: baseSize},
			z:    Z{Which: air.Z_Which_i64, I64: 123},
		},
		{
			name: "int64 into 8-byte",
			sz:   capnp.ObjectSize{DataSize: baseSize + 8},
			z:    Z{Which: air.Z_Which_i64, I64: 123},
			ok:   true,
		},
		{
			name: "uint64 into 0-byte",
			sz:   capnp.ObjectSize{DataSize: baseSize},
			z:    Z{Which: air.Z_Which_u64, U64: 123},
		},
		{
			name: "uint64 into 8-byte",
			sz:   capnp.ObjectSize{DataSize: baseSize + 8},
			z:    Z{Which: air.Z_Which_u64, U64: 123},
			ok:   true,
		},
		{
			name: "float64 into 0-byte",
			sz:   capnp.ObjectSize{DataSize: baseSize},
			z:    Z{Which: air.Z_Which_f64, F64: 123},
		},
		{
			name: "float64 into 8-byte",
			sz:   capnp.ObjectSize{DataSize: baseSize + 8},
			z:    Z{Which: air.Z_Which_f64, F64: 123},
			ok:   true,
		},
		{
			name: "text into 0 pointer",
			sz:   capnp.ObjectSize{DataSize: baseSize, PointerCount: 0},
			z:    Z{Which: air.Z_Which_text, Text: "hi"},
		},
		{
			name: "text into 1 pointer",
			sz:   capnp.ObjectSize{DataSize: baseSize, PointerCount: 1},
			z:    Z{Which: air.Z_Which_text, Text: "hi"},
			ok:   true,
		},
		{
			name: "data into 0 pointer",
			sz:   capnp.ObjectSize{DataSize: baseSize, PointerCount: 0},
			z:    Z{Which: air.Z_Which_blob, Blob: []byte("hi")},
		},
		{
			name: "data into 1 pointer",
			sz:   capnp.ObjectSize{DataSize: baseSize, PointerCount: 1},
			z:    Z{Which: air.Z_Which_blob, Blob: []byte("hi")},
			ok:   true,
		},
		{
			name: "list into 0 pointer",
			sz:   capnp.ObjectSize{DataSize: baseSize, PointerCount: 0},
			z:    Z{Which: air.Z_Which_f64vec, F64vec: []float64{123}},
		},
		{
			name: "list into 1 pointer",
			sz:   capnp.ObjectSize{DataSize: baseSize, PointerCount: 1},
			z:    Z{Which: air.Z_Which_f64vec, F64vec: []float64{123}},
			ok:   true,
		},
		{
			name: "struct into 0 pointer",
			sz:   capnp.ObjectSize{DataSize: baseSize, PointerCount: 0},
			z:    Z{Which: air.Z_Which_planebase, Planebase: new(PlaneBase)},
		},
		{
			name: "struct into 1 pointer",
			sz:   capnp.ObjectSize{DataSize: baseSize, PointerCount: 1},
			z:    Z{Which: air.Z_Which_planebase, Planebase: new(PlaneBase)},
			ok:   true,
		},
	}
	for _, test := range tests {
		_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
		if err != nil {
			t.Errorf("%s: NewMessage: %v", test.name, err)
			continue
		}
		st, err := capnp.NewRootStruct(seg, test.sz)
		if err != nil {
			t.Errorf("%s: NewRootStruct(seg, %v): %v", test.name, test.sz, err)
			continue
		}
		err = Insert(air.Z_TypeID, st, &test.z)
		if test.ok && err != nil {
			t.Errorf("%s: Insert(%#x, capnp.NewStruct(seg, %v), %s) = %v; want nil", test.name, uint64(air.Z_TypeID), test.sz, zpretty.Sprint(test.z), err)
		}
		if !test.ok && err == nil {
			t.Errorf("%s: Insert(%#x, capnp.NewStruct(seg, %v), %s) = nil; want error about not fitting", test.name, uint64(air.Z_TypeID), test.sz, zpretty.Sprint(test.z))
		}
	}
}

type BytesZ struct {
	Which   air.Z_Which
	Text    []byte
	Textvec [][]byte
}

func TestExtract_StringBytes(t *testing.T) {
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatalf("NewMessage: %v", err)
	}
	z, err := air.NewRootZ(seg)
	if err != nil {
		t.Fatalf("NewRootZ: %v", err)
	}
	err = zfill(z, &Z{Which: air.Z_Which_text, Text: "Hello, World!"})
	if err != nil {
		t.Fatalf("zfill: %v", err)
	}
	out := new(BytesZ)
	if err := Extract(out, air.Z_TypeID, capnp.Struct(z)); err != nil {
		t.Errorf("Extract(%v) error: %v", z, err)
	}
	want := &BytesZ{Which: air.Z_Which_text, Text: []byte("Hello, World!")}
	if out.Which != want.Which || !bytes.Equal(out.Text, want.Text) {
		t.Errorf("Extract(%v) produced %s; want %s", z, zpretty.Sprint(out), zpretty.Sprint(want))
	}
}

func TestExtract_StringListBytes(t *testing.T) {
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatalf("NewMessage: %v", err)
	}
	z, err := air.NewRootZ(seg)
	if err != nil {
		t.Fatalf("NewRootZ: %v", err)
	}
	err = zfill(z, &Z{Which: air.Z_Which_textvec, Textvec: []string{"Holmes", "Watson"}})
	if err != nil {
		t.Fatalf("zfill: %v", err)
	}
	out := new(BytesZ)
	if err := Extract(out, air.Z_TypeID, capnp.Struct(z)); err != nil {
		t.Errorf("Extract(%v) error: %v", z, err)
	}
	want := &BytesZ{Which: air.Z_Which_textvec, Textvec: [][]byte{[]byte("Holmes"), []byte("Watson")}}
	eq := sliceeq(len(out.Textvec), len(want.Textvec), func(i int) bool {
		return bytes.Equal(out.Textvec[i], want.Textvec[i])
	})
	if out.Which != want.Which || !eq {
		t.Errorf("Extract(%v) produced %s; want %s", z, zpretty.Sprint(out), zpretty.Sprint(want))
	}
}

func TestInsert_StringBytes(t *testing.T) {
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatalf("NewMessage: %v", err)
	}
	z, err := air.NewRootZ(seg)
	if err != nil {
		t.Fatalf("NewRootZ: %v", err)
	}
	bz := &BytesZ{Which: air.Z_Which_text, Text: []byte("Hello, World!")}
	err = Insert(air.Z_TypeID, capnp.Struct(z), bz)
	if err != nil {
		t.Errorf("Insert(%s) error: %v", zpretty.Sprint(bz), err)
	}
	want := &Z{Which: air.Z_Which_text, Text: "Hello, World!"}
	if equal, err := zequal(want, z); err != nil {
		t.Errorf("Insert(%s) compare err: %v", zpretty.Sprint(bz), err)
	} else if !equal {
		t.Errorf("Insert(%s) produced %v", zpretty.Sprint(bz), z)
	}
}

func TestInsert_StringListBytes(t *testing.T) {
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatalf("NewMessage: %v", err)
	}
	z, err := air.NewRootZ(seg)
	if err != nil {
		t.Fatalf("NewRootZ: %v", err)
	}
	bz := &BytesZ{Which: air.Z_Which_textvec, Textvec: [][]byte{[]byte("Holmes"), []byte("Watson")}}
	err = Insert(air.Z_TypeID, capnp.Struct(z), bz)
	if err != nil {
		t.Errorf("Insert(%s) error: %v", zpretty.Sprint(bz), err)
	}
	want := &Z{Which: air.Z_Which_textvec, Textvec: []string{"Holmes", "Watson"}}
	if equal, err := zequal(want, z); err != nil {
		t.Errorf("Insert(%s) compare err: %v", zpretty.Sprint(bz), err)
	} else if !equal {
		t.Errorf("Insert(%s) produced %v", zpretty.Sprint(bz), z)
	}
}

// StructZ is a variant of Z that has direct structs instead of pointers.
type StructZ struct {
	Which     air.Z_Which
	Zvec      []Z
	Planebase PlaneBase
	Grp       ZGroup
}

func (z *StructZ) equal(y *StructZ) bool {
	if z.Which != y.Which {
		return false
	}
	switch z.Which {
	case air.Z_Which_zvec:
		return sliceeq(len(z.Zvec), len(y.Zvec), func(i int) bool {
			return z.Zvec[i].equal(&y.Zvec[i])
		})
	case air.Z_Which_planebase:
		return z.Planebase.equal(&y.Planebase)
	case air.Z_Which_grp:
		return z.Grp.First == y.Grp.First && z.Grp.Second == y.Grp.Second
	default:
		panic("unknown Z which")
	}
}

func TestExtract_StructNoPtr(t *testing.T) {
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatalf("NewMessage: %v", err)
	}
	z, err := air.NewRootZ(seg)
	if err != nil {
		t.Fatalf("NewRootZ: %v", err)
	}
	err = zfill(z, &Z{Which: air.Z_Which_planebase, Planebase: &PlaneBase{Name: "foo"}})
	if err != nil {
		t.Fatalf("zfill: %v", err)
	}
	out := new(StructZ)
	if err := Extract(out, air.Z_TypeID, capnp.Struct(z)); err != nil {
		t.Errorf("Extract(%v) error: %v", z, err)
	}
	want := &StructZ{Which: air.Z_Which_planebase, Planebase: PlaneBase{Name: "foo"}}
	if !out.equal(want) {
		t.Errorf("Extract(%v) produced %s; want %s", z, zpretty.Sprint(out), zpretty.Sprint(want))
	}
}

func TestExtract_StructListNoPtr(t *testing.T) {
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatalf("NewMessage: %v", err)
	}
	z, err := air.NewRootZ(seg)
	if err != nil {
		t.Fatalf("NewRootZ: %v", err)
	}
	err = zfill(z, &Z{Which: air.Z_Which_zvec, Zvec: []*Z{
		{Which: air.Z_Which_i64, I64: 123},
	}})
	if err != nil {
		t.Fatalf("zfill: %v", err)
	}
	out := new(StructZ)
	if err := Extract(out, air.Z_TypeID, capnp.Struct(z)); err != nil {
		t.Errorf("Extract(%v) error: %v", z, err)
	}
	want := &StructZ{Which: air.Z_Which_zvec, Zvec: []Z{
		{Which: air.Z_Which_i64, I64: 123},
	}}
	if !out.equal(want) {
		t.Errorf("Extract(%v) produced %s; want %s", z, zpretty.Sprint(out), zpretty.Sprint(want))
	}
}

func TestExtract_GroupNoPtr(t *testing.T) {
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatalf("NewMessage: %v", err)
	}
	z, err := air.NewRootZ(seg)
	if err != nil {
		t.Fatalf("NewRootZ: %v", err)
	}
	err = zfill(z, &Z{Which: air.Z_Which_grp, Grp: &ZGroup{First: 123, Second: 456}})
	if err != nil {
		t.Fatalf("zfill: %v", err)
	}
	out := new(StructZ)
	if err := Extract(out, air.Z_TypeID, capnp.Struct(z)); err != nil {
		t.Errorf("Extract(%v) error: %v", z, err)
	}
	want := &StructZ{Which: air.Z_Which_grp, Grp: ZGroup{First: 123, Second: 456}}
	if !out.equal(want) {
		t.Errorf("Extract(%v) produced %s; want %s", z, zpretty.Sprint(out), zpretty.Sprint(want))
	}
}

func TestInsert_StructNoPtr(t *testing.T) {
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatalf("NewMessage: %v", err)
	}
	z, err := air.NewRootZ(seg)
	if err != nil {
		t.Fatalf("NewRootZ: %v", err)
	}
	bz := &StructZ{Which: air.Z_Which_planebase, Planebase: PlaneBase{Name: "foo"}}
	err = Insert(air.Z_TypeID, capnp.Struct(z), bz)
	if err != nil {
		t.Errorf("Insert(%s) error: %v", zpretty.Sprint(bz), err)
	}
	want := &Z{Which: air.Z_Which_planebase, Planebase: &PlaneBase{Name: "foo"}}
	if equal, err := zequal(want, z); err != nil {
		t.Errorf("Insert(%s) compare err: %v", zpretty.Sprint(bz), err)
	} else if !equal {
		t.Errorf("Insert(%s) produced %v", zpretty.Sprint(bz), z)
	}
}

func TestInsert_StructListNoPtr(t *testing.T) {
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatalf("NewMessage: %v", err)
	}
	z, err := air.NewRootZ(seg)
	if err != nil {
		t.Fatalf("NewRootZ: %v", err)
	}
	bz := &StructZ{Which: air.Z_Which_zvec, Zvec: []Z{
		{Which: air.Z_Which_i64, I64: 123},
	}}
	err = Insert(air.Z_TypeID, capnp.Struct(z), bz)
	if err != nil {
		t.Errorf("Insert(%s) error: %v", zpretty.Sprint(bz), err)
	}
	want := &Z{Which: air.Z_Which_zvec, Zvec: []*Z{
		{Which: air.Z_Which_i64, I64: 123},
	}}
	if equal, err := zequal(want, z); err != nil {
		t.Errorf("Insert(%s) compare err: %v", zpretty.Sprint(bz), err)
	} else if !equal {
		t.Errorf("Insert(%s) produced %v", zpretty.Sprint(bz), z)
	}
}

func TestInsert_GroupNoPtr(t *testing.T) {
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatalf("NewMessage: %v", err)
	}
	z, err := air.NewRootZ(seg)
	if err != nil {
		t.Fatalf("NewRootZ: %v", err)
	}
	bz := &StructZ{Which: air.Z_Which_grp, Grp: ZGroup{First: 123, Second: 456}}
	err = Insert(air.Z_TypeID, capnp.Struct(z), bz)
	if err != nil {
		t.Errorf("Insert(%s) error: %v", zpretty.Sprint(bz), err)
	}
	want := &Z{Which: air.Z_Which_grp, Grp: &ZGroup{First: 123, Second: 456}}
	if equal, err := zequal(want, z); err != nil {
		t.Errorf("Insert(%s) compare err: %v", zpretty.Sprint(bz), err)
	} else if !equal {
		t.Errorf("Insert(%s) produced %v", zpretty.Sprint(bz), z)
	}
}

// TagZ is a variant of Z that has tags.
type TagZ struct {
	Which   air.Z_Which
	Float64 float64 `capnp:"f64"`
	I64     int64   `capnp:"-"`
	U8      bool    `capnp:"bool"`
}

func TestExtract_Tags(t *testing.T) {
	tests := []struct {
		name string
		z    Z
		tagz TagZ
	}{
		{
			name: "renamed field",
			z:    Z{Which: air.Z_Which_f64, F64: 3.5},
			tagz: TagZ{Which: air.Z_Which_f64, Float64: 3.5},
		},
		{
			name: "omitted field",
			z:    Z{Which: air.Z_Which_i64, I64: 42},
			tagz: TagZ{Which: air.Z_Which_i64},
		},
		{
			name: "field with overlapping name",
			z:    Z{Which: air.Z_Which_bool, Bool: true},
			tagz: TagZ{Which: air.Z_Which_bool, U8: true},
		},
	}
	for _, test := range tests {
		_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
		if err != nil {
			t.Errorf("%s: NewMessage: %v", test.name, err)
			continue
		}
		z, err := air.NewRootZ(seg)
		if err != nil {
			t.Errorf("%s: NewRootZ: %v", test.name, err)
			continue
		}
		if err := zfill(z, &test.z); err != nil {
			t.Errorf("%s: zfill: %v", test.name, err)
			continue
		}
		out := new(TagZ)
		if err := Extract(out, air.Z_TypeID, capnp.Struct(z)); err != nil {
			t.Errorf("%s: Extract error: %v", test.name, err)
		}
		if *out != test.tagz {
			t.Errorf("%s: Extract produced %s; want %s", test.name, zpretty.Sprint(out), zpretty.Sprint(test.tagz))
		}
	}
}

func TestInsert_Tags(t *testing.T) {
	tests := []struct {
		name string
		tagz TagZ
		z    Z
	}{
		{
			name: "renamed field",
			tagz: TagZ{Which: air.Z_Which_f64, Float64: 3.5},
			z:    Z{Which: air.Z_Which_f64, F64: 3.5},
		},
		{
			name: "omitted field",
			tagz: TagZ{Which: air.Z_Which_i64, I64: 42},
			z:    Z{Which: air.Z_Which_i64, I64: 0},
		},
		{
			name: "field with overlapping name",
			tagz: TagZ{Which: air.Z_Which_bool, U8: true},
			z:    Z{Which: air.Z_Which_bool, Bool: true},
		},
	}
	for _, test := range tests {
		_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
		if err != nil {
			t.Errorf("%s: NewMessage: %v", test.name, err)
			continue
		}
		z, err := air.NewRootZ(seg)
		if err != nil {
			t.Errorf("%s: NewRootZ: %v", test.name, err)
			continue
		}
		err = Insert(air.Z_TypeID, capnp.Struct(z), &test.tagz)
		if err != nil {
			t.Errorf("%s: Insert(%s) error: %v", test.name, zpretty.Sprint(test.tagz), err)
		}
		if equal, err := zequal(&test.z, z); err != nil {
			t.Errorf("%s: Insert(%s) compare err: %v", test.name, zpretty.Sprint(test.tagz), err)
		} else if !equal {
			t.Errorf("%s: Insert(%s) produced %v", test.name, zpretty.Sprint(test.tagz), z)
		}
	}
}

type ZBool struct {
	Bool bool
}

func TestExtract_FixedUnion(t *testing.T) {
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatalf("NewMessage: %v", err)
	}
	z, err := air.NewRootZ(seg)
	if err != nil {
		t.Fatalf("NewRootZ: %v", err)
	}
	if err := zfill(z, &Z{Which: air.Z_Which_bool, Bool: true}); err != nil {
		t.Fatalf("zfill: %v", err)
	}
	out := new(ZBool)
	if err := Extract(out, air.Z_TypeID, capnp.Struct(z)); err != nil {
		t.Errorf("Extract error: %v", err)
	}
	if !out.Bool {
		t.Errorf("Extract produced %s; want %s", zpretty.Sprint(out), zpretty.Sprint(&ZBool{Bool: true}))
	}
}

func TestExtract_FixedUnionMismatch(t *testing.T) {
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatalf("NewMessage: %v", err)
	}
	z, err := air.NewRootZ(seg)
	if err != nil {
		t.Fatalf("NewRootZ: %v", err)
	}
	if err := zfill(z, &Z{Which: air.Z_Which_i64, I64: 42}); err != nil {
		t.Fatalf("zfill: %v", err)
	}
	out := new(ZBool)
	if err := Extract(out, air.Z_TypeID, capnp.Struct(z)); err == nil {
		t.Error("Extract did not return an error")
	}
}

func TestInsert_FixedUnion(t *testing.T) {
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatalf("NewMessage: %v", err)
	}
	z, err := air.NewRootZ(seg)
	if err != nil {
		t.Fatalf("NewRootZ: %v", err)
	}
	zb := &ZBool{Bool: true}
	err = Insert(air.Z_TypeID, capnp.Struct(z), zb)
	if err != nil {
		t.Errorf("Insert(%s) error: %v", zpretty.Sprint(zb), err)
	}
	want := &Z{Which: air.Z_Which_bool, Bool: true}
	if equal, err := zequal(want, z); err != nil {
		t.Errorf("Insert(%s) compare err: %v", zpretty.Sprint(zb), err)
	} else if !equal {
		t.Errorf("Insert(%s) produced %v", zpretty.Sprint(zb), z)
	}
}

type ZBoolU8 struct {
	Bool bool
	U8   uint8
}

func TestMissingWhich(t *testing.T) {
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatalf("NewMessage: %v", err)
	}
	z, err := air.NewRootZ(seg)
	if err != nil {
		t.Fatalf("NewRootZ: %v", err)
	}
	zz := &ZBoolU8{Bool: true, U8: 42}
	err = Insert(air.Z_TypeID, capnp.Struct(z), zz)
	if err == nil {
		t.Errorf("Insert(%s) did not return error", zpretty.Sprint(zz))
	}
	err = Extract(zz, air.Z_TypeID, capnp.Struct(z))
	if err == nil {
		t.Errorf("Extract(%v) did not return error", zz)
	}
}

type ZDateWithExtra struct {
	Year  int16
	Month uint8
	Day   uint8

	ExtraField uint16
}

func TestExtraFields(t *testing.T) {
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatalf("NewMessage: %v", err)
	}
	z, err := air.NewRootZdate(seg)
	if err != nil {
		t.Fatalf("NewRootZdate: %v", err)
	}
	zd := &ZDateWithExtra{ExtraField: 42}
	err = Insert(air.Zdate_TypeID, capnp.Struct(z), zd)
	if err == nil {
		t.Errorf("Insert(%s) did not return error", zpretty.Sprint(zd))
	} else if s := err.Error(); !strings.Contains(s, "ExtraField") {
		t.Errorf("Insert(%s): %v; want error about ExtraField", zpretty.Sprint(zd), err)
	}
	err = Extract(zd, air.Zdate_TypeID, capnp.Struct(z))
	if err == nil {
		t.Errorf("Extract(%v) did not return error", z)
	} else if s := err.Error(); !strings.Contains(s, "ExtraField") {
		t.Errorf("Extract(%v): %v; want error about ExtraField", z, err)
	}
	if zd.ExtraField != 42 {
		t.Errorf("zd.ExtraField modified to %d; want 42", zd.ExtraField)
	}
}

func zequal(g *Z, c air.Z) (bool, error) {
	if g.Which != c.Which() {
		return false, nil
	}
	_, seg, _ := capnp.NewMessage(capnp.SingleSegment(nil))
	d, _ := air.NewRootZ(seg)
	if err := zfill(d, g); err != nil {
		return false, err
	}
	return capnp.Equal(c.ToPtr(), d.ToPtr())
}

func (z *Z) equal(y *Z) bool {
	if z.Which != y.Which {
		return false
	}
	_, seg, _ := capnp.NewMessage(capnp.SingleSegment(nil))
	c, _ := air.NewZ(seg)
	if err := zfill(c, z); err != nil {
		return false
	}
	d, _ := air.NewZ(seg)
	if err := zfill(d, y); err != nil {
		return false
	}
	eq, _ := capnp.Equal(c.ToPtr(), d.ToPtr())
	return eq
}

func zfill(c air.Z, g *Z) error {
	switch g.Which {
	case air.Z_Which_f64:
		c.SetF64(g.F64)
	case air.Z_Which_f32:
		c.SetF32(g.F32)
	case air.Z_Which_i64:
		c.SetI64(g.I64)
	case air.Z_Which_i32:
		c.SetI32(g.I32)
	case air.Z_Which_i16:
		c.SetI16(g.I16)
	case air.Z_Which_i8:
		c.SetI8(g.I8)
	case air.Z_Which_u64:
		c.SetU64(g.U64)
	case air.Z_Which_u32:
		c.SetU32(g.U32)
	case air.Z_Which_u16:
		c.SetU16(g.U16)
	case air.Z_Which_u8:
		c.SetU8(g.U8)
	case air.Z_Which_bool:
		c.SetBool(g.Bool)
	case air.Z_Which_text:
		return c.SetText(g.Text)
	case air.Z_Which_blob:
		return c.SetBlob(g.Blob)
	case air.Z_Which_f64vec:
		if g.F64vec == nil {
			return c.SetF64vec(capnp.Float64List{})
		}
		fv, err := c.NewF64vec(int32(len(g.F64vec)))
		if err != nil {
			return err
		}
		for i, f := range g.F64vec {
			fv.Set(i, f)
		}
	case air.Z_Which_f32vec:
		if g.F32vec == nil {
			return c.SetF32vec(capnp.Float32List{})
		}
		fv, err := c.NewF32vec(int32(len(g.F32vec)))
		if err != nil {
			return err
		}
		for i, f := range g.F32vec {
			fv.Set(i, f)
		}
	case air.Z_Which_i64vec:
		if g.I64vec == nil {
			return c.SetI64vec(capnp.Int64List{})
		}
		iv, err := c.NewI64vec(int32(len(g.I64vec)))
		if err != nil {
			return err
		}
		for i, n := range g.I64vec {
			iv.Set(i, n)
		}
	case air.Z_Which_i8vec:
		if g.I8vec == nil {
			return c.SetI8vec(capnp.Int8List{})
		}
		iv, err := c.NewI8vec(int32(len(g.I8vec)))
		if err != nil {
			return err
		}
		for i, n := range g.I8vec {
			iv.Set(i, n)
		}
	case air.Z_Which_u64vec:
		if g.U64vec == nil {
			return c.SetU64vec(capnp.UInt64List{})
		}
		uv, err := c.NewU64vec(int32(len(g.U64vec)))
		if err != nil {
			return err
		}
		for i, n := range g.U64vec {
			uv.Set(i, n)
		}
	case air.Z_Which_u8vec:
		if g.U8vec == nil {
			return c.SetU8vec(capnp.UInt8List{})
		}
		uv, err := c.NewU8vec(int32(len(g.U8vec)))
		if err != nil {
			return err
		}
		for i, n := range g.U8vec {
			uv.Set(i, n)
		}
	case air.Z_Which_boolvec:
		if g.Boolvec == nil {
			return c.SetBoolvec(capnp.BitList{})
		}
		vec, err := c.NewBoolvec(int32(len(g.Boolvec)))
		if err != nil {
			return err
		}
		for i, v := range g.Boolvec {
			vec.Set(i, v)
		}
	case air.Z_Which_datavec:
		if g.Datavec == nil {
			return c.SetDatavec(capnp.DataList{})
		}
		vec, err := c.NewDatavec(int32(len(g.Datavec)))
		if err != nil {
			return err
		}
		for i, v := range g.Datavec {
			if err := vec.Set(i, v); err != nil {
				return err
			}
		}
	case air.Z_Which_textvec:
		if g.Textvec == nil {
			return c.SetTextvec(capnp.TextList{})
		}
		vec, err := c.NewTextvec(int32(len(g.Textvec)))
		if err != nil {
			return err
		}
		for i, v := range g.Textvec {
			if err := vec.Set(i, v); err != nil {
				return err
			}
		}
	case air.Z_Which_zvec:
		if g.Zvec == nil {
			return c.SetZvec(air.Z_List{})
		}
		vec, err := c.NewZvec(int32(len(g.Zvec)))
		if err != nil {
			return err
		}
		for i, z := range g.Zvec {
			if err := zfill(vec.At(i), z); err != nil {
				return err
			}
		}
	case air.Z_Which_zvecvec:
		if g.Zvecvec == nil {
			return c.SetZvecvec(capnp.PointerList{})
		}
		vv, err := c.NewZvecvec(int32(len(g.Zvecvec)))
		if err != nil {
			return err
		}
		for i, zz := range g.Zvecvec {
			v, err := air.NewZ_List(vv.Segment(), int32(len(zz)))
			if err != nil {
				return err
			}
			if err := vv.Set(i, v.ToPtr()); err != nil {
				return err
			}
			for j, z := range zz {
				if err := zfill(v.At(j), z); err != nil {
					return err
				}
			}
		}
	case air.Z_Which_planebase:
		if g.Planebase == nil {
			return c.SetPlanebase(air.PlaneBase{})
		}
		pb, err := c.NewPlanebase()
		if err != nil {
			return err
		}
		if err := pb.SetName(g.Planebase.Name); err != nil {
			return err
		}
		if g.Planebase.Homes != nil {
			homes, err := pb.NewHomes(int32(len(g.Planebase.Homes)))
			if err != nil {
				return err
			}
			for i := range g.Planebase.Homes {
				homes.Set(i, g.Planebase.Homes[i])
			}
		}
		pb.SetRating(g.Planebase.Rating)
		pb.SetCanFly(g.Planebase.CanFly)
		pb.SetCapacity(g.Planebase.Capacity)
		pb.SetMaxSpeed(g.Planebase.MaxSpeed)
	case air.Z_Which_airport:
		c.SetAirport(g.Airport)
	case air.Z_Which_grp:
		c.SetGrp()
		if g.Grp != nil {
			c.Grp().SetFirst(g.Grp.First)
			c.Grp().SetSecond(g.Grp.Second)
		}
	case air.Z_Which_echo:
		c.SetEcho(g.Echo)
	case air.Z_Which_echoes:
		e, err := c.NewEchoes(int32(len(g.Echoes)))
		if err != nil {
			return err
		}
		for i, ee := range g.Echoes {
			if !ee.IsValid() {
				continue
			}
			err := e.Set(i, ee)
			if err != nil {
				return err
			}
		}
	case air.Z_Which_anyPtr:
		return c.SetAnyPtr(g.AnyPtr)
	case air.Z_Which_anyStruct:
		return c.SetAnyStruct(g.AnyStruct)
	case air.Z_Which_anyList:
		return c.SetAnyList(g.AnyList)
	case air.Z_Which_anyCapability:
		return c.SetAnyCapability(g.AnyCapability)
	default:
		return fmt.Errorf("zfill: unknown type: %v", g.Which)
	}
	return nil
}

var zpretty = &pretty.Config{
	Compact:        true,
	SkipZeroFields: true,
}

func sliceeq(na, nb int, f func(i int) bool) bool {
	if na != nb {
		return false
	}
	for i := 0; i < na; i++ {
		if !f(i) {
			return false
		}
	}
	return true
}
