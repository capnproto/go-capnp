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
	F64   float64
	I64   int64
}

func zequal(g *Z, c air.Z) (bool, error) {
	if g.Which != c.Which() {
		return false, nil
	}
	switch g.Which {
	case air.Z_Which_f64:
		return g.F64 == c.F64(), nil
	case air.Z_Which_i64:
		return g.I64 == c.I64(), nil
	default:
		return false, fmt.Errorf("zequal: unknown type: %v", g.Which)
	}
}

func zfill(c air.Z, g *Z) error {
	switch g.Which {
	case air.Z_Which_f64:
		c.SetF64(g.F64)
	case air.Z_Which_i64:
		c.SetI64(g.I64)
	default:
		return fmt.Errorf("zfill: unknown type: %v", g.Which)
	}
	return nil
}

func TestExtract(t *testing.T) {
	tests := []Z{
		{Which: air.Z_Which_f64, F64: 3.5},
		{Which: air.Z_Which_i64, I64: 123},
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
		{Which: air.Z_Which_i64, I64: 123},
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
			t.Errorf("Insert(%+v) compare err: %v", err)
		} else if !equal {
			t.Errorf("Insert(%+v) produced %v", test, z)
		}
	}
}
