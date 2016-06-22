package pogs

import (
	"errors"
	"fmt"
	"math"
	"reflect"

	"zombiezen.com/go/capnproto2"
	"zombiezen.com/go/capnproto2/internal/nodemap"
	"zombiezen.com/go/capnproto2/std/capnp/schema"
)

// Extract copies s into val, a pointer to a Go struct.
func Extract(val interface{}, typeID uint64, s capnp.Struct) error {
	e := new(extracter)
	err := e.extractStruct(reflect.ValueOf(val), typeID, s)
	if err != nil {
		return fmt.Errorf("pogs: extract @%#x: %v", typeID, err)
	}
	return nil
}

type extracter struct {
	nodes nodemap.Map
}

func (e *extracter) extractStruct(val reflect.Value, typeID uint64, s capnp.Struct) error {
	if val.Kind() == reflect.Ptr {
		// TODO(light): create if nil
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("can't extract struct into %v", val.Kind())
	}
	if !val.CanSet() {
		return errors.New("can't modify struct, did you pass in a pointer to your struct?")
	}
	n, err := e.nodes.Find(typeID)
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
	case schema.Type_Which_bool:
		v := s.Bit(capnp.BitOffset(f.Slot().Offset()))
		d := dv.Bool()
		val.SetBool(v != d) // != acts as XOR
	case schema.Type_Which_int8:
		v := int8(s.Uint8(capnp.DataOffset(f.Slot().Offset())))
		d := dv.Int8()
		val.SetInt(int64(v ^ d))
	case schema.Type_Which_int16:
		v := int16(s.Uint16(capnp.DataOffset(f.Slot().Offset() * 2)))
		d := dv.Int16()
		val.SetInt(int64(v ^ d))
	case schema.Type_Which_int32:
		v := int32(s.Uint32(capnp.DataOffset(f.Slot().Offset() * 4)))
		d := dv.Int32()
		val.SetInt(int64(v ^ d))
	case schema.Type_Which_int64:
		v := int64(s.Uint64(capnp.DataOffset(f.Slot().Offset() * 8)))
		d := dv.Int64()
		val.SetInt(v ^ d)
	case schema.Type_Which_uint8:
		v := s.Uint8(capnp.DataOffset(f.Slot().Offset()))
		d := dv.Uint8()
		val.SetUint(uint64(v ^ d))
	case schema.Type_Which_uint16:
		v := s.Uint16(capnp.DataOffset(f.Slot().Offset() * 2))
		d := dv.Uint16()
		val.SetUint(uint64(v ^ d))
	case schema.Type_Which_enum:
		v := s.Uint16(capnp.DataOffset(f.Slot().Offset() * 2))
		d := dv.Enum()
		val.SetUint(uint64(v ^ d))
	case schema.Type_Which_uint32:
		v := s.Uint32(capnp.DataOffset(f.Slot().Offset() * 4))
		d := dv.Uint32()
		val.SetUint(uint64(v ^ d))
	case schema.Type_Which_uint64:
		v := s.Uint64(capnp.DataOffset(f.Slot().Offset() * 8))
		d := dv.Uint64()
		val.SetUint(v ^ d)
	case schema.Type_Which_float32:
		v := s.Uint32(capnp.DataOffset(f.Slot().Offset() * 4))
		d := math.Float32bits(dv.Float32())
		val.SetFloat(float64(math.Float32frombits(v ^ d)))
	case schema.Type_Which_float64:
		v := s.Uint64(capnp.DataOffset(f.Slot().Offset() * 8))
		d := math.Float64bits(dv.Float64())
		val.SetFloat(math.Float64frombits(v ^ d))
	case schema.Type_Which_text:
		p, err := s.Ptr(uint16(f.Slot().Offset()))
		if err != nil {
			return err
		}
		var b []byte
		if p.IsValid() {
			b = p.TextBytes()
		} else {
			b, _ = dv.TextBytes()
		}
		if val.Kind() == reflect.String {
			val.SetString(string(b))
		} else {
			// byte slice, as guaranteed by isTypeMatch
			val.SetBytes(b)
		}
	case schema.Type_Which_data:
		p, err := s.Ptr(uint16(f.Slot().Offset()))
		if err != nil {
			return err
		}
		var b []byte
		if p.IsValid() {
			b = p.Data()
		} else {
			b, _ = dv.Data()
		}
		val.SetBytes(b)
	case schema.Type_Which_structType:
		p, err := s.Ptr(uint16(f.Slot().Offset()))
		if err != nil {
			return err
		}
		ss := p.Struct()
		if !ss.IsValid() {
			p, _ = dv.StructValuePtr()
			ss = p.Struct()
		}
		if val.Kind() == reflect.Struct {
			return e.extractStruct(val, typ.StructType().TypeId(), ss)
		}
		// Pointer to struct otherwise.
		if !ss.IsValid() {
			val.Set(reflect.Zero(val.Type()))
			return nil
		}
		newval := reflect.New(val.Type().Elem())
		val.Set(newval)
		return e.extractStruct(newval, typ.StructType().TypeId(), ss)
	default:
		return fmt.Errorf("unknown field type %v", typ.Which())
	}
	return nil
}

var typeMap = map[schema.Type_Which]reflect.Kind{
	schema.Type_Which_bool:    reflect.Bool,
	schema.Type_Which_int8:    reflect.Int8,
	schema.Type_Which_int16:   reflect.Int16,
	schema.Type_Which_int32:   reflect.Int32,
	schema.Type_Which_int64:   reflect.Int64,
	schema.Type_Which_uint8:   reflect.Uint8,
	schema.Type_Which_uint16:  reflect.Uint16,
	schema.Type_Which_uint32:  reflect.Uint32,
	schema.Type_Which_uint64:  reflect.Uint64,
	schema.Type_Which_float32: reflect.Float32,
	schema.Type_Which_float64: reflect.Float64,
	schema.Type_Which_enum:    reflect.Uint16,
}

func isTypeMatch(r reflect.Type, s schema.Type) bool {
	switch s.Which() {
	case schema.Type_Which_text:
		return r.Kind() == reflect.String || r.Kind() == reflect.Slice && r.Elem().Kind() == reflect.Uint8
	case schema.Type_Which_data:
		return r.Kind() == reflect.Slice && r.Elem().Kind() == reflect.Uint8
	case schema.Type_Which_structType:
		return r.Kind() == reflect.Struct || r.Kind() == reflect.Ptr && r.Elem().Kind() == reflect.Struct
	}
	k, ok := typeMap[s.Which()]
	return ok && k == r.Kind()
}

func fieldName(s string) string {
	if len(s) == 0 {
		return ""
	}
	// TODO(light): check it's lowercase.
	x := s[0] - 'a' + 'A'
	return string(x) + s[1:]
}
