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
	case schema.Type_Which_int64:
		v := uint64(val.Int())
		d := uint64(dv.Int64())
		s.SetUint64(capnp.DataOffset(f.Slot().Offset()*8), v^d)
	case schema.Type_Which_float64:
		v := math.Float64bits(val.Float())
		d := uint64(math.Float64bits(dv.Float64()))
		s.SetUint64(capnp.DataOffset(f.Slot().Offset()*8), v^d)
	default:
		return fmt.Errorf("unknown field type %v", typ.Which())
	}
	return nil
}
