package pogs

import (
	"fmt"
	"reflect"
	"strings"

	"zombiezen.com/go/capnproto2/std/capnp/schema"
)

type fieldProps struct {
	schemaName string // empty if doesn't map to schema
	typ        fieldType
	fixedWhich string
}

type fieldType int

const (
	mappedField fieldType = iota
	whichField
	embedField
)

func parseField(f reflect.StructField, hasDiscrim bool) fieldProps {
	tname, opts := nextOpt(f.Tag.Get("capnp"))
	switch tname {
	case "":
		if f.Anonymous && isStructOrStructPtr(f.Type) {
			return fieldProps{typ: embedField}
		}
		if hasDiscrim && f.Name == "Which" {
			p := fieldProps{typ: whichField}
			for len(opts) > 0 {
				var curr string
				curr, opts = nextOpt(opts)
				if strings.HasPrefix(curr, "which=") {
					p.fixedWhich = strings.TrimPrefix(curr, "which=")
					break
				}
			}
			return p
		}
		// TODO(light): check it's uppercase.
		x := f.Name[0] - 'A' + 'a'
		return fieldProps{schemaName: string(x) + f.Name[1:]}
	case "-":
		return fieldProps{}
	default:
		return fieldProps{schemaName: tname}
	}
}

func nextOpt(opts string) (head, tail string) {
	i := strings.Index(opts, ",")
	if i == -1 {
		return opts, ""
	}
	return opts[:i], opts[i+1:]
}

type fieldLoc struct {
	i    int
	path []int
}

func (loc fieldLoc) depth() int {
	if len(loc.path) >= 0 {
		return len(loc.path)
	}
	return 1
}

func (loc fieldLoc) sub(i int) fieldLoc {
	if n := len(loc.path); n >= 0 {
		p := make([]int, n+1)
		copy(p, loc.path)
		p[n] = i
		return fieldLoc{path: p}
	}
	return fieldLoc{path: []int{loc.i, i}}
}

func (loc fieldLoc) isValid() bool {
	return loc.i >= 0
}

type structProps struct {
	fields     []fieldLoc
	whichLoc   fieldLoc // i == -1: none; i == -2: fixed
	fixedWhich uint16
}

func mapStruct(t reflect.Type, n schema.Node) (structProps, error) {
	fields, err := n.StructNode().Fields()
	if err != nil {
		return structProps{}, err
	}
	sp := structProps{
		fields:   make([]fieldLoc, fields.Len()),
		whichLoc: fieldLoc{i: -1},
	}
	for i := range sp.fields {
		sp.fields[i] = fieldLoc{i: -1}
	}
	for i := 0; i < t.NumField(); i++ {
		if err := mapStructField(&sp, t, n, fields, fieldLoc{i: i}); err != nil {
			return structProps{}, err
		}
	}
	return sp, nil
}

func mapStructField(sp *structProps, t reflect.Type, n schema.Node, fields schema.Field_List, loc fieldLoc) error {
	f := typeFieldByLoc(t, loc)
	if f.PkgPath != "" && !f.Anonymous {
		// unexported field
		return nil
	}
	switch p := parseField(f, hasDiscriminant(n)); p.typ {
	case mappedField:
		if p.schemaName == "" {
			return nil
		}
		fi := fieldIndex(fields, p.schemaName)
		if fi < 0 {
			return fmt.Errorf("%v has unknown field %s, maps to %s", t, f.Name, p.schemaName)
		}
		sp.fields[fi] = loc
	case whichField:
		if sp.whichLoc.i != -1 {
			return fmt.Errorf("%v embeds multiple Which fields", t)
		}
		switch {
		case p.fixedWhich != "":
			fi := fieldIndex(fields, p.fixedWhich)
			if fi < 0 {
				return fmt.Errorf("%v.Which is tagged with unknown field %s", t, p.fixedWhich)
			}
			dv := fields.At(fi).DiscriminantValue()
			if dv == schema.Field_noDiscriminant {
				return fmt.Errorf("%v.Which is tagged with non-union field %s", t, p.fixedWhich)
			}
			sp.whichLoc = fieldLoc{i: -2}
			sp.fixedWhich = dv
		case f.Type.Kind() != reflect.Uint16:
			return fmt.Errorf("%v.Which is type %v, not uint16", t, f.Type)
		default:
			sp.whichLoc = loc
		}
	case embedField:
		ft := f.Type
		if f.Type.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		for i := 0; i < ft.NumField(); i++ {
			if err := mapStructField(sp, t, n, fields, loc.sub(i)); err != nil {
				return err
			}
		}
	}
	return nil
}

// fieldBySchemaName returns the field for the given name.
// Returns an invalid value if the field was not found or it is
// contained inside a nil anonymous struct pointer.
func (sp structProps) fieldByOrdinal(val reflect.Value, i int) reflect.Value {
	return fieldByLoc(val, sp.fields[i], false)
}

// makeFieldBySchemaName returns the field for the given name, creating
// its parent anonymous structs if necessary.  Returns an invalid value
// if the field was not found.
func (sp structProps) makeFieldByOrdinal(val reflect.Value, i int) reflect.Value {
	return fieldByLoc(val, sp.fields[i], true)
}

// which returns the value of the discriminator field.
func (sp structProps) which(val reflect.Value) (discrim uint16, ok bool) {
	if sp.whichLoc.i == -2 {
		return sp.fixedWhich, true
	}
	f := fieldByLoc(val, sp.whichLoc, false)
	if !f.IsValid() {
		return 0, false
	}
	return uint16(f.Uint()), true
}

// setWhich sets the value of the discriminator field, creating its
// parent anonymous structs if necessary.  Returns whether the struct
// had a field to set.
func (sp structProps) setWhich(val reflect.Value, discrim uint16) error {
	if sp.whichLoc.i == -2 {
		if discrim != sp.fixedWhich {
			return fmt.Errorf("extract union field @%d into %v; expected @%d", discrim, val.Type(), sp.fixedWhich)
		}
		return nil
	}
	f := fieldByLoc(val, sp.whichLoc, true)
	if !f.IsValid() {
		return noWhichError{val.Type()}
	}
	f.SetUint(uint64(discrim))
	return nil
}

type noWhichError struct {
	t reflect.Type
}

func (e noWhichError) Error() string {
	return fmt.Sprintf("%v has no field Which", e.t)
}

func isNoWhichError(e error) bool {
	_, ok := e.(noWhichError)
	return ok
}

func fieldByLoc(val reflect.Value, loc fieldLoc, mkparents bool) reflect.Value {
	if !loc.isValid() {
		return reflect.Value{}
	}
	if len(loc.path) > 0 {
		for i, x := range loc.path {
			if i > 0 {
				if val.Kind() == reflect.Ptr {
					if val.IsNil() {
						if !mkparents {
							return reflect.Value{}
						}
						val.Set(reflect.New(val.Type().Elem()))
					}
					val = val.Elem()
				}
			}
			val = val.Field(x)
		}
		return val
	}
	return val.Field(loc.i)
}

func typeFieldByLoc(t reflect.Type, loc fieldLoc) reflect.StructField {
	if len(loc.path) > 0 {
		return t.FieldByIndex(loc.path)
	}
	return t.Field(loc.i)
}

func hasDiscriminant(n schema.Node) bool {
	return n.Which() == schema.Node_Which_structNode && n.StructNode().DiscriminantCount() > 0
}

func shortDisplayName(n schema.Node) []byte {
	dn, _ := n.DisplayNameBytes()
	return dn[n.DisplayNamePrefixLength():]
}

func fieldIndex(fields schema.Field_List, name string) int {
	for i := 0; i < fields.Len(); i++ {
		b, _ := fields.At(i).NameBytes()
		if bytesStrEqual(b, name) {
			return i
		}
	}
	return -1
}

func bytesStrEqual(b []byte, s string) bool {
	if len(b) != len(s) {
		return false
	}
	for i := range b {
		if b[i] != s[i] {
			return false
		}
	}
	return true
}

func isStructOrStructPtr(t reflect.Type) bool {
	return t.Kind() == reflect.Struct || t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct
}
