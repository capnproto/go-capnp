package pogs

import (
	"errors"
	"fmt"
	"math"
	"reflect"

	"zombiezen.com/go/capnproto2"
	"zombiezen.com/go/capnproto2/schemas"
	"zombiezen.com/go/capnproto2/std/capnp/schema"
)

// Extract copies s into val, a pointer to a Go struct.
func Extract(val interface{}, typeID uint64, s capnp.Struct) error {
	e := &extracter{reg: &schemas.DefaultRegistry}
	err := e.extractStruct(reflect.ValueOf(val), typeID, s)
	if err != nil {
		return fmt.Errorf("pogs: extract @%#x: %v", typeID, err)
	}
	return nil
}

type extracter struct {
	reg   *schemas.Registry
	nodes map[uint64]schema.Node
}

func (e *extracter) extractStruct(val reflect.Value, typeID uint64, s capnp.Struct) error {
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("can't extract struct into %v", val.Kind())
	}
	if !val.CanSet() {
		return errors.New("can't modify struct, did you pass in a pointer to your struct?")
	}
	n, err := e.findNode(typeID)
	if err != nil {
		return err
	}
	if !n.IsValid() || n.Which() != schema.Node_Which_structNode {
		return fmt.Errorf("cannot find struct type %#x", typeID)
	}
	var discriminant uint16
	var hasWhichField bool
	if n.StructNode().DiscriminantCount() > 0 {
		discriminant = s.Uint16(capnp.DataOffset(n.StructNode().DiscriminantOffset() * 2))
		f := val.FieldByName("Which")
		if f.IsValid() && f.Kind() == reflect.Uint16 {
			hasWhichField = true
			f.SetUint(uint64(discriminant))
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
			if !hasWhichField {
				dn, _ := n.DisplayNameBytes()
				dn = dn[n.DisplayNamePrefixLength():]
				return fmt.Errorf("can't extract %s into %v: has union field but no Which field", dn, val.Type())
			}
			if dv != discriminant {
				continue
			}
		}
		if err := e.extractField(vf, s, f); err != nil {
			return err
		}
	}
	return nil
}

func (e *extracter) extractField(val reflect.Value, s capnp.Struct, f schema.Field) error {
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
		return fmt.Errorf("extract field %s: default value is a %v, want %v", name, dv.Which(), typ.Which())
	}
	if !isTypeMatch(val.Type(), typ) {
		name, _ := f.NameBytes()
		return fmt.Errorf("can't extract field %s of type %v into a Go %v", name, typ.Which(), val.Type())
	}
	switch typ.Which() {
	case schema.Type_Which_int64:
		v := s.Uint64(capnp.DataOffset(f.Slot().Offset() * 8))
		d := uint64(dv.Int64())
		val.SetInt(int64(v ^ d))
	case schema.Type_Which_float64:
		v := s.Uint64(capnp.DataOffset(f.Slot().Offset() * 8))
		d := uint64(math.Float64bits(dv.Float64()))
		val.SetFloat(math.Float64frombits(v ^ d))
	default:
		return fmt.Errorf("unknown field type %v", typ.Which())
	}
	return nil
}

var typeMap = map[schema.Type_Which]reflect.Kind{
	schema.Type_Which_bool:    reflect.Bool,
	schema.Type_Which_int64:   reflect.Int64,
	schema.Type_Which_float64: reflect.Float64,
	schema.Type_Which_enum:    reflect.Uint16,
}

func isTypeMatch(r reflect.Type, s schema.Type) bool {
	return typeMap[s.Which()] == r.Kind()
}

func fieldName(s string) string {
	if len(s) == 0 {
		return ""
	}
	// TODO(light): check it's lowercase.
	x := s[0] - 'a' + 'A'
	return string(x) + s[1:]
}

func (e *extracter) findNode(id uint64) (schema.Node, error) {
	if n := e.nodes[id]; n.IsValid() {
		return n, nil
	}
	data, err := e.reg.Find(id)
	if err != nil {
		return schema.Node{}, err
	}
	msg, err := capnp.Unmarshal(data)
	if err != nil {
		return schema.Node{}, err
	}
	req, err := schema.ReadRootCodeGeneratorRequest(msg)
	if err != nil {
		return schema.Node{}, err
	}
	nodes, err := req.Nodes()
	if err != nil {
		return schema.Node{}, err
	}
	if e.nodes == nil {
		e.nodes = make(map[uint64]schema.Node)
	}
	for i := 0; i < nodes.Len(); i++ {
		n := nodes.At(i)
		e.nodes[n.Id()] = n
	}
	return e.nodes[id], nil
}
