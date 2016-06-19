package pogs

import (
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

	Airport air.Airport
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
	case air.Z_Which_airport:
		c.SetAirport(g.Airport)
	default:
		return fmt.Errorf("zfill: unknown type: %v", g.Which)
	}
	return nil
}

func TestExtract(t *testing.T) {
	tests := []Z{
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
		{Which: air.Z_Which_airport, Airport: air.Airport_lax},
	}
	for _, test := range tests {
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
	tests := []Z{
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
		{Which: air.Z_Which_airport, Airport: air.Airport_lax},
	}
	for _, test := range tests {
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
