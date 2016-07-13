package pogs

import "reflect"

type fieldProps struct {
	schemaName string // empty if doesn't map to schema
	which      bool
}

func (p fieldProps) isZero() bool {
	return p.schemaName == "" && p.which == false
}

func parseField(f reflect.StructField, hasDiscrim bool) fieldProps {
	if f.PkgPath != "" {
		// unexported field
		return fieldProps{}
	}
	tag := f.Tag.Get("capnp")
	switch tag {
	case "":
		if hasDiscrim && f.Name == "Which" {
			return fieldProps{which: true}
		}
		// TODO(light): check it's uppercase.
		x := f.Name[0] - 'A' + 'a'
		return fieldProps{schemaName: string(x) + f.Name[1:]}
	case "-":
		return fieldProps{}
	default:
		return fieldProps{schemaName: tag}
	}
}

func mapStruct(t reflect.Type, hasDiscrim bool) map[fieldProps][]int {
	m := make(map[fieldProps][]int)
	for i := 0; i < t.NumField(); i++ {
		// TODO(light): anonymous fields
		p := parseField(t.Field(i), hasDiscrim)
		if !p.isZero() {
			m[p] = []int{i}
		}
	}
	return m
}
