package pogs

import (
	"fmt"
	"reflect"

	"zombiezen.com/go/capnproto2/std/capnp/schema"
)

type fieldKey struct {
	schemaName string // empty if doesn't map to schema
	which      bool
}

func (p fieldKey) isZero() bool {
	return p.schemaName == "" && p.which == false
}

func parseField(f reflect.StructField, hasDiscrim bool) fieldKey {
	if f.PkgPath != "" {
		// unexported field
		return fieldKey{}
	}
	switch tag := f.Tag.Get("capnp"); tag {
	case "":
		if hasDiscrim && f.Name == "Which" {
			return fieldKey{which: true}
		}
		// TODO(light): check it's uppercase.
		x := f.Name[0] - 'A' + 'a'
		return fieldKey{schemaName: string(x) + f.Name[1:]}
	case "-":
		return fieldKey{}
	default:
		return fieldKey{schemaName: tag}
	}
}

type fieldProps struct {
	i int
}

type structProps map[fieldKey]fieldProps

func mapStruct(t reflect.Type, hasDiscrim bool) (structProps, error) {
	m := make(structProps)
	for i := 0; i < t.NumField(); i++ {
		// TODO(light): anonymous fields
		f := t.Field(i)
		k := parseField(f, hasDiscrim)
		if k.which && f.Type.Kind() != reflect.Uint16 {
			return nil, fmt.Errorf("%v.Which is type %v, not uint16", t, f.Type)
		}
		if !k.isZero() {
			m[k] = fieldProps{i: i}
		}
	}
	return m, nil
}

// fieldBySchemaName returns the field for the given name.
// Returns an invalid value if the field was not found or it is
// contained inside a nil anonymous struct pointer.
func (sp structProps) fieldBySchemaName(val reflect.Value, name string) reflect.Value {
	return sp.field(val, fieldKey{schemaName: name}, false)
}

// makeFieldBySchemaName returns the field for the given name, creating
// its parent anonymous structs if necessary.  Returns an invalid value
// if the field was not found.
func (sp structProps) makeFieldBySchemaName(val reflect.Value, name string) reflect.Value {
	return sp.field(val, fieldKey{schemaName: name}, true)
}

// which returns the value of the discriminator field.
func (sp structProps) which(val reflect.Value) (discrim uint16, ok bool) {
	f := sp.field(val, fieldKey{which: true}, false)
	if !f.IsValid() {
		return 0, false
	}
	return uint16(f.Uint()), true
}

// setWhich sets the value of the discriminator field, creating its
// parent anonymous structs if necessary.  Returns whether the struct
// had a field to set.
func (sp structProps) setWhich(val reflect.Value, discrim uint16) bool {
	f := sp.field(val, fieldKey{which: true}, true)
	if !f.IsValid() {
		return false
	}
	f.SetUint(uint64(discrim))
	return true
}

func (sp structProps) field(val reflect.Value, k fieldKey, mkparents bool) reflect.Value {
	p, ok := sp[k]
	if !ok {
		return reflect.Value{}
	}
	// TODO(light): mkparents
	return val.Field(p.i)
}

func hasDiscriminant(n schema.Node) bool {
	return n.Which() == schema.Node_Which_structNode && n.StructNode().DiscriminantCount() > 0
}

func shortDisplayName(n schema.Node) []byte {
	dn, _ := n.DisplayNameBytes()
	return dn[n.DisplayNamePrefixLength():]
}
