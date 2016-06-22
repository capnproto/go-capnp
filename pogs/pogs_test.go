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

func zequal(g *Z, c air.Z) (bool, error) {
	if g.Which != c.Which() {
		return false, nil
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
		return nil
	case air.Z_Which_airport:
		c.SetAirport(g.Airport)
	default:
		return fmt.Errorf("zfill: unknown type: %v", g.Which)
	}
	return nil
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
	Which air.Z_Which
	Text  []byte
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
