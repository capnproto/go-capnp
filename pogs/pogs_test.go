package pogs

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"zombiezen.com/go/capnproto2"
	air "zombiezen.com/go/capnproto2/internal/aircraftlib"
)

const zTypeID = 0xea26e9973bd6a0d9

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

	Planebase *PlaneBase
	Airport   air.Airport
}

type PlaneBase struct {
	Name string
	// TODO(light): Homes []air.Airport
	Rating   int64
	CanFly   bool
	Capacity int64
	MaxSpeed float64
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
	{Which: air.Z_Which_textvec, Textvec: nil},
	{Which: air.Z_Which_textvec, Textvec: []string{"John", "Paul", "George", "Ringo"}},
	{Which: air.Z_Which_planebase, Planebase: nil},
	{Which: air.Z_Which_planebase, Planebase: &PlaneBase{
		Name:     "Boeing",
		Rating:   123,
		CanFly:   true,
		Capacity: 100,
		MaxSpeed: 9001.0,
	}},
	{Which: air.Z_Which_airport, Airport: air.Airport_lax},
}

func TestExtract(t *testing.T) {
	for _, test := range goodTests {
		_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
		if err != nil {
			t.Errorf("NewMessage for %+v: %v", test, err)
			continue
		}
		z, err := air.NewRootZ(seg)
		if err != nil {
			t.Errorf("NewRootZ for %+v: %v", test, err)
			continue
		}
		if err := zfill(z, &test); err != nil {
			t.Errorf("zfill for %+v: %v", test, err)
			continue
		}
		out := new(Z)
		if err := Extract(out, zTypeID, z.Struct); err != nil {
			t.Errorf("Extract(%v) error: %v", z, err)
		}
		if !reflect.DeepEqual(out, &test) {
			t.Errorf("Extract(%v) produced %+v; want %+v", z, out, &test)
		}
	}
}

func TestInsert(t *testing.T) {
	for _, test := range goodTests {
		_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
		if err != nil {
			t.Errorf("NewMessage for %+v: %v", test, err)
			continue
		}
		z, err := air.NewRootZ(seg)
		if err != nil {
			t.Errorf("NewRootZ for %+v: %v", test, err)
			continue
		}
		err = Insert(zTypeID, z.Struct, &test)
		if err != nil {
			t.Errorf("Insert(%+v) error: %v", test, err)
		}
		if equal, err := zequal(&test, z); err != nil {
			t.Errorf("Insert(%+v) compare err: %v", test, err)
		} else if !equal {
			t.Errorf("Insert(%+v) produced %v", test, z)
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
	if err := Extract(out, zTypeID, z.Struct); err != nil {
		t.Errorf("Extract(%v) error: %v", z, err)
	}
	want := &BytesZ{Which: air.Z_Which_text, Text: []byte("Hello, World!")}
	if !reflect.DeepEqual(out, want) {
		t.Errorf("Extract(%v) produced %+v; want %+v", z, out, want)
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
	if err := Extract(out, zTypeID, z.Struct); err != nil {
		t.Errorf("Extract(%v) error: %v", z, err)
	}
	want := &BytesZ{Which: air.Z_Which_textvec, Textvec: [][]byte{[]byte("Holmes"), []byte("Watson")}}
	if !reflect.DeepEqual(out, want) {
		t.Errorf("Extract(%v) produced %+v; want %+v", z, out, want)
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
	err = Insert(zTypeID, z.Struct, bz)
	if err != nil {
		t.Errorf("Insert(%+v) error: %v", bz, err)
	}
	want := &Z{Which: air.Z_Which_text, Text: "Hello, World!"}
	if equal, err := zequal(want, z); err != nil {
		t.Errorf("Insert(%+v) compare err: %v", bz, err)
	} else if !equal {
		t.Errorf("Insert(%+v) produced %v", bz, z)
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
	err = Insert(zTypeID, z.Struct, bz)
	if err != nil {
		t.Errorf("Insert(%+v) error: %v", bz, err)
	}
	want := &Z{Which: air.Z_Which_textvec, Textvec: []string{"Holmes", "Watson"}}
	if equal, err := zequal(want, z); err != nil {
		t.Errorf("Insert(%+v) compare err: %v", bz, err)
	} else if !equal {
		t.Errorf("Insert(%+v) produced %v", bz, z)
	}
}

func zequal(g *Z, c air.Z) (bool, error) {
	if g.Which != c.Which() {
		return false, nil
	}
	listeq := func(has bool, n int, l capnp.List, f func(i int) bool) bool {
		if has != l.IsValid() {
			return false
		}
		if !has {
			return true
		}
		if l.Len() != n {
			return false
		}
		for i := 0; i < l.Len(); i++ {
			if !f(i) {
				return false
			}
		}
		return true
	}
	switch g.Which {
	case air.Z_Which_f64:
		return g.F64 == c.F64(), nil
	case air.Z_Which_f32:
		return g.F32 == c.F32(), nil
	case air.Z_Which_i64:
		return g.I64 == c.I64(), nil
	case air.Z_Which_i32:
		return g.I32 == c.I32(), nil
	case air.Z_Which_i16:
		return g.I16 == c.I16(), nil
	case air.Z_Which_i8:
		return g.I8 == c.I8(), nil
	case air.Z_Which_u64:
		return g.U64 == c.U64(), nil
	case air.Z_Which_u32:
		return g.U32 == c.U32(), nil
	case air.Z_Which_u16:
		return g.U16 == c.U16(), nil
	case air.Z_Which_u8:
		return g.U8 == c.U8(), nil
	case air.Z_Which_bool:
		return g.Bool == c.Bool(), nil
	case air.Z_Which_text:
		text, err := c.Text()
		if err != nil {
			return false, err
		}
		return g.Text == text, nil
	case air.Z_Which_blob:
		blob, err := c.Blob()
		if err != nil {
			return false, err
		}
		return bytes.Equal(g.Blob, blob), nil
	case air.Z_Which_f64vec:
		fv, err := c.F64vec()
		if err != nil {
			return false, err
		}
		return listeq(g.F64vec != nil, len(g.F64vec), fv.List, func(i int) bool {
			return fv.At(i) == g.F64vec[i]
		}), nil
	case air.Z_Which_f32vec:
		fv, err := c.F32vec()
		if err != nil {
			return false, err
		}
		return listeq(g.F32vec != nil, len(g.F32vec), fv.List, func(i int) bool {
			return fv.At(i) == g.F32vec[i]
		}), nil
	case air.Z_Which_i64vec:
		iv, err := c.I64vec()
		if err != nil {
			return false, err
		}
		return listeq(g.I64vec != nil, len(g.I64vec), iv.List, func(i int) bool {
			return iv.At(i) == g.I64vec[i]
		}), nil
	case air.Z_Which_i8vec:
		iv, err := c.I8vec()
		if err != nil {
			return false, err
		}
		return listeq(g.I8vec != nil, len(g.I8vec), iv.List, func(i int) bool {
			return iv.At(i) == g.I8vec[i]
		}), nil
	case air.Z_Which_u64vec:
		uv, err := c.U64vec()
		if err != nil {
			return false, err
		}
		return listeq(g.U64vec != nil, len(g.U64vec), uv.List, func(i int) bool {
			return uv.At(i) == g.U64vec[i]
		}), nil
	case air.Z_Which_u8vec:
		uv, err := c.U8vec()
		if err != nil {
			return false, err
		}
		return listeq(g.U8vec != nil, len(g.U8vec), uv.List, func(i int) bool {
			return uv.At(i) == g.U8vec[i]
		}), nil
	case air.Z_Which_boolvec:
		bv, err := c.Boolvec()
		if err != nil {
			return false, err
		}
		return listeq(g.Boolvec != nil, len(g.Boolvec), bv.List, func(i int) bool {
			return bv.At(i) == g.Boolvec[i]
		}), nil
	case air.Z_Which_datavec:
		dv, err := c.Datavec()
		if err != nil {
			return false, err
		}
		return listeq(g.Datavec != nil, len(g.Datavec), dv.List, func(i int) bool {
			var di []byte
			di, err = dv.At(i)
			return err == nil && bytes.Equal(di, g.Datavec[i])
		}), err
	case air.Z_Which_textvec:
		tv, err := c.Textvec()
		if err != nil {
			return false, err
		}
		return listeq(g.Textvec != nil, len(g.Textvec), tv.List, func(i int) bool {
			var s string
			s, err = tv.At(i)
			return err == nil && s == g.Textvec[i]
		}), err
	case air.Z_Which_planebase:
		pb, err := c.Planebase()
		if err != nil {
			return false, err
		}
		if (g.Planebase != nil) != pb.IsValid() {
			return false, nil
		}
		if g.Planebase == nil {
			return true, nil
		}
		name, _ := pb.Name()
		return g.Planebase.Name == name && g.Planebase.Rating == pb.Rating() && g.Planebase.CanFly == pb.CanFly() && g.Planebase.Capacity == pb.Capacity() && g.Planebase.MaxSpeed == pb.MaxSpeed(), nil
	case air.Z_Which_airport:
		return g.Airport == c.Airport(), nil
	default:
		return false, fmt.Errorf("zequal: unknown type: %v", g.Which)
	}
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
		pb.SetRating(g.Planebase.Rating)
		pb.SetCanFly(g.Planebase.CanFly)
		pb.SetCapacity(g.Planebase.Capacity)
		pb.SetMaxSpeed(g.Planebase.MaxSpeed)
	case air.Z_Which_airport:
		c.SetAirport(g.Airport)
	default:
		return fmt.Errorf("zfill: unknown type: %v", g.Which)
	}
	return nil
}
