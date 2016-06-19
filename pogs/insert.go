package pogs

import (
	"fmt"
	"math"
	"reflect"

	"zombiezen.com/go/capnproto2"
	"zombiezen.com/go/capnproto2/internal/nodemap"
	"zombiezen.com/go/capnproto2/std/capnp/schema"
)

// Insert copies val, a pointer to a Go struct, into s.
func Insert(typeID uint64, s capnp.Struct, val interface{}) error {
	ins := new(inserter)
	err := ins.insertStruct(typeID, s, reflect.ValueOf(val))
	if err != nil {
		return fmt.Errorf("pogs: insert @%#x: %v", typeID, err)
	}
	return nil
}

type inserter struct {
	nodes nodemap.Map
}

func (ins *inserter) insertStruct(typeID uint64, s capnp.Struct, val reflect.Value) error {
	if val.Kind() == reflect.Ptr {
		// TODO(light): ignore if nil?
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("can't insert %v into a struct", val.Kind())
	}
	n, err := ins.nodes.Find(typeID)
	if err != nil {
		return err
	}
	if !n.IsValid() || n.Which() != schema.Node_Which_structNode {
		return fmt.Errorf("cannot find struct type %#x", typeID)
	}
	var discriminant uint16
	var hasWhich bool
	if n.StructNode().DiscriminantCount() > 0 {
		f := val.FieldByName("Which")
		if f.IsValid() && f.Kind() == reflect.Uint16 {
			hasWhich = true
			discriminant = uint16(f.Uint())
			s.SetUint16(capnp.DataOffset(n.StructNode().DiscriminantOffset()*2), discriminant)
		}
	}
	fields, err := n.StructNode().Fields()
	if err != nil {
		return err
	}
	for i := 0; i < fields.Len(); i++ {
		f := fields.At(i)
		// TODO(light): groups
		if f.Which() != schema.Field_Which_slot {
			continue
		}
		sname, err := f.Name()
		if err != nil {
			return err
		}
		fname := fieldName(sname)
		vf := val.FieldByName(fname)
		if !vf.IsValid() {
			// Don't have a field for this.
			continue
		}
		if dv := f.DiscriminantValue(); dv != schema.Field_noDiscriminant {
			if !hasWhich {
				dn, _ := n.DisplayNameBytes()
				dn = dn[n.DisplayNamePrefixLength():]
				return fmt.Errorf("can't insert %s from %v: has union field %s but no Which field", dn, val.Type(), fname)
			}
			if dv != discriminant {
				continue
			}
		}
		if err := ins.insertField(s, f, vf); err != nil {
			return err
		}
	}
	return nil
}

func (ins *inserter) insertField(s capnp.Struct, f schema.Field, val reflect.Value) error {
	typ, err := f.Slot().Type()
	if err != nil {
		return err
	}
	dv, err := f.Slot().DefaultValue()
	if err != nil {
		return err
	}
	if dv.IsValid() && int(typ.Which()) != int(dv.Which()) {
		name, _ := f.NameBytes()
		return fmt.Errorf("insert field %s: default value is a %v, want %v", name, dv.Which(), typ.Which())
	}
	if !isTypeMatch(val.Type(), typ) {
		name, _ := f.NameBytes()
		return fmt.Errorf("can't insert field %s of type Go %v into a %v", name, val.Type(), typ.Which())
	}
	switch typ.Which() {
	case schema.Type_Which_bool:
		v := val.Bool()
		d := dv.Bool()
		s.SetBit(capnp.BitOffset(f.Slot().Offset()), v != d) // != acts as XOR
	case schema.Type_Which_int8:
		v := int8(val.Int())
		d := dv.Int8()
		s.SetUint8(capnp.DataOffset(f.Slot().Offset()), uint8(v^d))
	case schema.Type_Which_int16:
		v := int16(val.Int())
		d := dv.Int16()
		s.SetUint16(capnp.DataOffset(f.Slot().Offset()*2), uint16(v^d))
	case schema.Type_Which_int32:
		v := int32(val.Int())
		d := dv.Int32()
		s.SetUint32(capnp.DataOffset(f.Slot().Offset()*4), uint32(v^d))
	case schema.Type_Which_int64:
		v := val.Int()
		d := dv.Int64()
		s.SetUint64(capnp.DataOffset(f.Slot().Offset()*8), uint64(v^d))
	case schema.Type_Which_uint8:
		v := uint8(val.Uint())
		d := dv.Uint8()
		s.SetUint8(capnp.DataOffset(f.Slot().Offset()), v^d)
	case schema.Type_Which_uint16:
		v := uint16(val.Uint())
		d := dv.Uint16()
		s.SetUint16(capnp.DataOffset(f.Slot().Offset()*2), v^d)
	case schema.Type_Which_enum:
		v := uint16(val.Uint())
		d := dv.Enum()
		s.SetUint16(capnp.DataOffset(f.Slot().Offset()*2), v^d)
	case schema.Type_Which_uint32:
		v := uint32(val.Uint())
		d := dv.Uint32()
		s.SetUint32(capnp.DataOffset(f.Slot().Offset()*4), v^d)
	case schema.Type_Which_uint64:
		v := val.Uint()
		d := dv.Uint64()
		s.SetUint64(capnp.DataOffset(f.Slot().Offset()*8), v^d)
	case schema.Type_Which_float32:
		v := math.Float32bits(float32(val.Float()))
		d := math.Float32bits(dv.Float32())
		s.SetUint32(capnp.DataOffset(f.Slot().Offset()*4), v^d)
	case schema.Type_Which_float64:
		v := math.Float64bits(val.Float())
		d := uint64(math.Float64bits(dv.Float64()))
		s.SetUint64(capnp.DataOffset(f.Slot().Offset()*8), v^d)
	default:
		return fmt.Errorf("unknown field type %v", typ.Which())
	}
	return nil
}
