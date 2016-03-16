//go:generate bash -c "capnp compile -o- schema.capnp | capnpc-go -promises=false"

/*
capnpc-go is the Cap'n proto code generator for Go.  It reads a
CodeGeneratorRequest from stdin and for a file foo.capnp it writes
foo.capnp.go.  This is usually invoked from `capnp compile -ogo`.

See https://capnproto.org/otherlang.html#how-to-write-compiler-plugins
for more details.
*/
package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"zombiezen.com/go/capnproto2"
)

var (
	genPromises = flag.Bool("promises", true, "generate code for promises")
)

const (
	go_capnproto_import = "zombiezen.com/go/capnproto2"
	server_import       = go_capnproto_import + "/server"
	context_import      = "golang.org/x/net/context"
)

var (
	g_nodes   = make(map[uint64]*node)
	g_imports imports
	g_segment []byte
	g_bufname string
)

type imports struct {
	specs []importSpec
	used  map[string]bool // keyed on import path
}

func (i *imports) init() {
	i.specs = nil
	i.used = make(map[string]bool)

	i.reserve(importSpec{path: go_capnproto_import, name: "capnp"})
	i.reserve(importSpec{path: server_import, name: "server"})
	i.reserve(importSpec{path: context_import, name: "context"})

	i.reserve(importSpec{path: "bufio", name: "bufio"})
	i.reserve(importSpec{path: "bytes", name: "bytes"})
	i.reserve(importSpec{path: "io", name: "io"})
	i.reserve(importSpec{path: "math", name: "math"})
	i.reserve(importSpec{path: "strconv", name: "strconv"})
}

func (i *imports) capnp() string {
	return i.add(importSpec{path: go_capnproto_import, name: "capnp"})
}

func (i *imports) server() string {
	return i.add(importSpec{path: server_import, name: "server"})
}

func (i *imports) context() string {
	return i.add(importSpec{path: context_import, name: "context"})
}

func (i *imports) math() string {
	return i.add(importSpec{path: "math", name: "math"})
}

func (i *imports) strconv() string {
	return i.add(importSpec{path: "strconv", name: "strconv"})
}

func (i *imports) usedImports() []importSpec {
	specs := make([]importSpec, 0, len(i.specs))
	for _, s := range i.specs {
		if i.used[s.path] {
			specs = append(specs, s)
		}
	}
	return specs
}

func (i *imports) byPath(path string) (spec importSpec, ok bool) {
	for _, spec = range i.specs {
		if spec.path == path {
			return spec, true
		}
	}
	return importSpec{}, false
}

func (i *imports) byName(name string) (spec importSpec, ok bool) {
	for _, spec = range i.specs {
		if spec.name == name {
			return spec, true
		}
	}
	return importSpec{}, false
}

func (i *imports) add(spec importSpec) (name string) {
	name = i.reserve(spec)
	i.used[spec.path] = true
	return name
}

// reserve adds an import spec without marking it as used.
func (i *imports) reserve(spec importSpec) (name string) {
	if ispec, ok := i.byPath(spec.path); ok {
		return ispec.name
	}
	if spec.name == "" {
		spec.name = pkgFromImport(spec.path)
	}
	if _, found := i.byName(spec.name); found {
		for base, n := spec.name, uint64(2); ; n++ {
			spec.name = base + strconv.FormatUint(n, 10)
			if _, found = i.byName(spec.name); !found {
				break
			}
		}
	}
	i.specs = append(i.specs, spec)
	return spec.name
}

func pkgFromImport(path string) string {
	if i := strings.LastIndex(path, "/"); i != -1 {
		path = path[i+1:]
	}
	p := []rune(path)
	n := 0
	for _, r := range p {
		if isIdent(r) {
			p[n] = r
			n++
		}
	}
	if n == 0 || !isLower(p[0]) {
		return "pkg" + string(p[:n])
	}
	return string(p[:n])
}

func isLower(r rune) bool {
	return 'a' <= r && r <= 'z' || r == '_'
}

func isIdent(r rune) bool {
	return isLower(r) || 'A' <= r && r <= 'Z' || r >= 0x80 && unicode.IsLetter(r)
}

type importSpec struct {
	path string
	name string
}

func (spec importSpec) String() string {
	if spec.name == "" {
		return strconv.Quote(spec.path)
	}
	return spec.name + " " + strconv.Quote(spec.path)
}

type node struct {
	Node
	pkg   string
	imp   string
	nodes []*node
	Name  string
}

type field struct {
	Field
	Name string
}

func assert(chk bool, format string, a ...interface{}) {
	if !chk {
		panic(assertionError(fmt.Sprintf(format, a...)))
	}
}

type assertionError string

func (ae assertionError) Error() string {
	return string(ae)
}

func copyData(obj capnp.Pointer) staticDataRef {
	m, _, err := capnp.NewMessage(capnp.SingleSegment(nil))
	assert(err == nil, "%v\n", err)
	err = m.SetRoot(obj)
	assert(err == nil, "%v\n", err)
	data, err := m.Marshal()
	assert(err == nil, "%v\n", err)
	var ref staticDataRef
	ref.Start = len(g_segment)
	g_segment = append(g_segment, data...)
	ref.End = len(g_segment)
	return ref
}

type staticDataRef struct {
	Start, End int
}

func (ref staticDataRef) IsValid() bool {
	return ref.Start < ref.End
}

func (ref staticDataRef) String() string {
	return fmt.Sprintf("%s[%d:%d]", g_bufname, ref.Start, ref.End)
}

// Tag types
const (
	defaultTag = iota
	noTag
	customTag
)

type annotations struct {
	Doc       string
	Package   string
	Import    string
	TagType   int
	CustomTag string
	Name      string
}

func parseAnnotations(list Annotation_List) *annotations {
	ann := new(annotations)
	for i, n := 0, list.Len(); i < n; i++ {
		a := list.At(i)
		val, _ := a.Value()
		text, _ := val.Text()
		switch a.Id() {
		case capnp.Doc:
			ann.Doc = text
		case capnp.Package:
			ann.Package = text
		case capnp.Import:
			ann.Import = text
		case capnp.Tag:
			ann.TagType = customTag
			ann.CustomTag = text
		case capnp.Notag:
			ann.TagType = noTag
		case capnp.Name:
			ann.Name = text
		}
	}
	return ann
}

// Tag returns the string value that an enumerant value called name should have.
// An empty string indicates that this enumerant value has no tag.
func (ann *annotations) Tag(name string) string {
	switch ann.TagType {
	case noTag:
		return ""
	case customTag:
		return ann.CustomTag
	case defaultTag:
		fallthrough
	default:
		return name
	}
}

// Rename returns the overridden name from the annotations or the given name
// if no annotation was found.
func (ann *annotations) Rename(given string) string {
	if ann.Name == "" {
		return given
	}
	return ann.Name
}

func findNode(id uint64) *node {
	n := g_nodes[id]
	assert(n != nil, "could not find node 0x%x\n", id)
	return n
}

func (n *node) remoteScope(from *node) string {
	displayName, _ := n.DisplayName()
	fromDisplayName, _ := from.DisplayName()
	assert(n.pkg != "", "missing package declaration for %s", displayName)
	assert(n.imp != "", "missing import declaration for %s", displayName)
	assert(from.imp != "", "missing import declaration for %s", fromDisplayName)

	if n.imp == from.imp {
		return ""
	} else {
		name := g_imports.add(importSpec{path: n.imp, name: n.pkg})
		return name + "."
	}
}

func (n *node) RemoteNew(from *node) string {
	return n.remoteScope(from) + "New" + n.Name
}

func (n *node) RemoteName(from *node) string {
	return n.remoteScope(from) + n.Name
}

func (n *node) resolveName(base, name string, file *node) {
	na, _ := n.Annotations()
	name = parseAnnotations(na).Rename(name)
	if base == "" {
		n.Name = strings.Title(name)
	} else {
		n.Name = base + "_" + name
	}
	n.pkg = file.pkg
	n.imp = file.imp

	if n.Which() != Node_Which_structGroup || !n.StructGroup().IsGroup() {
		file.nodes = append(file.nodes, n)
	}

	nnodes, _ := n.NestedNodes()
	for i := 0; i < nnodes.Len(); i++ {
		nn := nnodes.At(i)
		if ni := g_nodes[nn.Id()]; ni != nil {
			nname, _ := nn.Name()
			ni.resolveName(n.Name, nname, file)
		}
	}

	if n.Which() == Node_Which_structGroup {
		fields, _ := n.StructGroup().Fields()
		for i := 0; i < fields.Len(); i++ {
			f := fields.At(i)
			if f.Which() == Field_Which_group {
				fa, _ := f.Annotations()
				fname, _ := f.Name()
				fname = parseAnnotations(fa).Rename(fname)
				findNode(f.Group().TypeId()).resolveName(n.Name, fname, file)
			}
		}
	} else if n.Which() == Node_Which_interface {
		m, _ := n.Interface().Methods()
		for i := 0; i < m.Len(); i++ {
			mm := m.At(i)
			mname, _ := mm.Name()
			mann, _ := mm.Annotations()
			mname = parseAnnotations(mann).Rename(mname)
			base := n.Name + "_" + mname
			if p := findNode(mm.ParamStructType()); p.ScopeId() == 0 {
				p.resolveName(base, "Params", file)
			}
			if r := findNode(mm.ResultStructType()); r.ScopeId() == 0 {
				r.resolveName(base, "Results", file)
			}
		}
	}
}

type enumval struct {
	Enumerant
	Name   string
	Val    int
	Tag    string
	parent *node
}

func makeEnumval(enum *node, i int, e Enumerant) enumval {
	eann, _ := e.Annotations()
	ann := parseAnnotations(eann)
	name, _ := e.Name()
	name = ann.Rename(name)
	t := ann.Tag(name)
	return enumval{e, name, i, t, enum}
}

func (e *enumval) FullName() string {
	return e.parent.Name + "_" + e.Name
}

func (n *node) defineEnum(w io.Writer) {
	es, _ := n.Enum().Enumerants()
	ev := make([]enumval, es.Len())
	for i := 0; i < es.Len(); i++ {
		e := es.At(i)
		ev[e.CodeOrder()] = makeEnumval(n, i, e)
	}
	nann, _ := n.Annotations()
	templates.ExecuteTemplate(w, "enum", enumParams{
		Node:        n,
		Annotations: parseAnnotations(nann),
		EnumValues:  ev,
	})
}

func (n *node) writeValue(w io.Writer, t Type, v Value) {
	switch t.Which() {
	case Type_Which_void:
		fmt.Fprintf(w, "struct{}{}")

	case Type_Which_interface:
		// The only statically representable interface value is null.
		fmt.Fprintf(w, "%s.Client(nil)", g_imports.capnp())

	case Type_Which_bool:
		assert(v.Which() == Value_Which_bool, "expected bool value")
		if v.Bool() {
			fmt.Fprint(w, "true")
		} else {
			fmt.Fprint(w, "false")
		}

	case Type_Which_uint8, Type_Which_uint16, Type_Which_uint32, Type_Which_uint64:
		fmt.Fprintf(w, "uint%d(%d)", intbits(t.Which()), uintValue(t, v))

	case Type_Which_int8, Type_Which_int16, Type_Which_int32, Type_Which_int64:
		fmt.Fprintf(w, "int%d(%d)", intbits(t.Which()), intValue(t, v))

	case Type_Which_float32:
		assert(v.Which() == Value_Which_float32, "expected float32 value")
		fmt.Fprintf(w, "%s.Float32frombits(0x%x)", g_imports.math(), math.Float32bits(v.Float32()))

	case Type_Which_float64:
		assert(v.Which() == Value_Which_float64, "expected float64 value")
		fmt.Fprintf(w, "%s.Float64frombits(0x%x)", g_imports.math(), math.Float64bits(v.Float64()))

	case Type_Which_text:
		assert(v.Which() == Value_Which_text, "expected text value")
		text, _ := v.Text()
		fmt.Fprintf(w, "%q", text)

	case Type_Which_data:
		assert(v.Which() == Value_Which_data, "expected data value")
		fmt.Fprint(w, "[]byte{")
		data, _ := v.Data()
		for i, b := range data {
			if i > 0 {
				fmt.Fprint(w, ", ")
			}
			fmt.Fprintf(w, "%d", b)
		}
		fmt.Fprint(w, "}")

	case Type_Which_enum:
		assert(v.Which() == Value_Which_enum, "expected enum value")
		en := findNode(t.Enum().TypeId())
		assert(en.Which() == Node_Which_enum, "expected enum type ID")
		enums, _ := en.Enum().Enumerants()
		if val := int(v.Enum()); val >= enums.Len() {
			fmt.Fprintf(w, "%s(%d)", en.RemoteName(n), val)
		} else {
			ev := makeEnumval(en, val, enums.At(val))
			fmt.Fprintf(w, "%s%s", en.remoteScope(n), ev.FullName())
		}

	case Type_Which_structGroup:
		assert(v.Which() == Value_Which_structField, "expected struct value")
		c := g_imports.capnp()
		data, _ := v.StructField()
		fmt.Fprintf(w, "%s{Struct: %s.ToStruct(%s.MustUnmarshalRoot(%v))}", findNode(t.StructGroup().TypeId()).RemoteName(n), c, c, copyData(data))

	case Type_Which_anyPointer:
		assert(v.Which() == Value_Which_anyPointer, "expected pointer value")
		data, _ := v.AnyPointer()
		fmt.Fprintf(w, "%s.MustUnmarshalRoot(%v)", g_imports.capnp(), copyData(data))

	case Type_Which_list:
		assert(v.Which() == Value_Which_list, "expected list value")
		c := g_imports.capnp()
		typ := n.fieldType(t, new(annotations))
		data, _ := v.List()
		fmt.Fprintf(w, "%s{List: %s.ToList(%s.MustUnmarshalRoot(%v))}", typ, c, c, copyData(data))
	}
}

func (n *node) defineAnnotation(w io.Writer) {
	templates.ExecuteTemplate(w, "annotation", annotationParams{
		Node: n,
	})
}

func constIsVar(n *node) bool {
	t, _ := n.Const().Type()
	switch t.Which() {
	case Type_Which_bool, Type_Which_int8, Type_Which_uint8, Type_Which_int16,
		Type_Which_uint16, Type_Which_int32, Type_Which_uint32, Type_Which_int64,
		Type_Which_uint64, Type_Which_text, Type_Which_enum:
		return false
	default:
		return true
	}
}

func defineConstNodes(w io.Writer, nodes []*node) {

	any := false

	for _, n := range nodes {
		if n.Which() == Node_Which_const && !constIsVar(n) {
			if !any {
				fmt.Fprintf(w, "const (\n")
				any = true
			}
			fmt.Fprintf(w, "%s = ", n.Name)
			kt, _ := n.Const().Type()
			kv, _ := n.Const().Value()
			n.writeValue(w, kt, kv)
			fmt.Fprintf(w, "\n")
		}
	}

	if any {
		fmt.Fprintf(w, ")\n")
	}

	any = false

	for _, n := range nodes {
		if n.Which() == Node_Which_const && constIsVar(n) {
			if !any {
				fmt.Fprintf(w, "var (\n")
				any = true
			}
			fmt.Fprintf(w, "%s = ", n.Name)
			kt, _ := n.Const().Type()
			kv, _ := n.Const().Value()
			n.writeValue(w, kt, kv)
			fmt.Fprintf(w, "\n")
		}
	}

	if any {
		fmt.Fprintf(w, ")\n")
	}
}

func (n *node) defineField(w io.Writer, f field) {
	fann, _ := f.Annotations()
	ann := parseAnnotations(fann)
	t, _ := f.Slot().Type()
	def, _ := f.Slot().DefaultValue()
	params := structFieldParams{
		Node:        n,
		Field:       f,
		Annotations: ann,
		FieldType:   n.fieldType(t, ann),
	}
	switch t.Which() {
	case Type_Which_void:
		templates.ExecuteTemplate(w, "structVoidField", params)
	case Type_Which_bool:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_bool, "expected bool default")
		templates.ExecuteTemplate(w, "structBoolField", structBoolFieldParams{
			structFieldParams: params,
			Default:           def.Which() == Value_Which_bool && def.Bool(),
		})

	case Type_Which_uint8, Type_Which_uint16, Type_Which_uint32, Type_Which_uint64:
		templates.ExecuteTemplate(w, "structUintField", structUintFieldParams{
			structFieldParams: params,
			Bits:              intbits(t.Which()),
			Default:           uintFieldDefault(t, def),
		})

	case Type_Which_int8, Type_Which_int16, Type_Which_int32, Type_Which_int64:
		templates.ExecuteTemplate(w, "structIntField", structIntFieldParams{
			structUintFieldParams: structUintFieldParams{
				structFieldParams: params,
				Bits:              intbits(t.Which()),
				Default:           uint64(intFieldDefaultMask(t, def)),
			},
		})

	case Type_Which_enum:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_enum, "expected enum default")
		ni := findNode(t.Enum().TypeId())
		var d uint64
		if def.Which() == Value_Which_enum {
			d = uint64(def.Enum())
		}
		templates.ExecuteTemplate(w, "structIntField", structIntFieldParams{
			structUintFieldParams: structUintFieldParams{
				structFieldParams: params,
				Bits:              16,
				Default:           d,
			},
			EnumName: ni.RemoteName(n),
		})
	case Type_Which_float32:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_float32, "expected float32 default")
		var d uint64
		if def.Which() == Value_Which_float32 && def.Float32() != 0 {
			d = uint64(math.Float32bits(def.Float32()))
		}
		templates.ExecuteTemplate(w, "structFloatField", structUintFieldParams{
			structFieldParams: params,
			Bits:              32,
			Default:           d,
		})

	case Type_Which_float64:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_float64, "expected float64 default")
		var d uint64
		if def.Which() == Value_Which_float64 && def.Float64() != 0 {
			d = math.Float64bits(def.Float64())
		}
		templates.ExecuteTemplate(w, "structFloatField", structUintFieldParams{
			structFieldParams: params,
			Bits:              64,
			Default:           d,
		})

	case Type_Which_text:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_text, "expected text default")
		var d string
		if def.Which() == Value_Which_text {
			d, _ = def.Text()
		}
		templates.ExecuteTemplate(w, "structTextField", structTextFieldParams{
			structFieldParams: params,
			Default:           d,
		})

	case Type_Which_data:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_data, "expected data default")
		var d []byte
		if def.Which() == Value_Which_data {
			d, _ = def.Data()
		}
		templates.ExecuteTemplate(w, "structDataField", structDataFieldParams{
			structFieldParams: params,
			Default:           d,
		})

	case Type_Which_structGroup:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_structField, "expected struct default")
		var defref staticDataRef
		if def.Which() == Value_Which_structField {
			if sf, _ := def.StructField(); capnp.HasData(sf) {
				defref = copyData(sf)
			}
		}
		templates.ExecuteTemplate(w, "structStructField", structObjectFieldParams{
			structFieldParams: params,
			TypeNode:          findNode(t.StructGroup().TypeId()),
			Default:           defref,
		})

	case Type_Which_anyPointer:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_anyPointer, "expected object default")
		var defref staticDataRef
		if def.Which() == Value_Which_anyPointer {
			if p, _ := def.AnyPointer(); capnp.HasData(p) {
				defref = copyData(p)
			}
		}
		templates.ExecuteTemplate(w, "structPointerField", structObjectFieldParams{
			structFieldParams: params,
			Default:           defref,
		})

	case Type_Which_list:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_list, "expected list default")
		var defref staticDataRef
		if def.Which() == Value_Which_list {
			if l, _ := def.List(); capnp.HasData(l) {
				defref = copyData(l)
			}
		}
		templates.ExecuteTemplate(w, "structListField", structObjectFieldParams{
			structFieldParams: params,
			Default:           defref,
		})

	case Type_Which_interface:
		templates.ExecuteTemplate(w, "structInterfaceField", params)
	}
}

func (n *node) fieldType(t Type, ann *annotations) string {
	switch t.Which() {
	case Type_Which_bool:
		return "bool"
	case Type_Which_int8:
		return "int8"
	case Type_Which_int16:
		return "int16"
	case Type_Which_int32:
		return "int32"
	case Type_Which_int64:
		return "int64"
	case Type_Which_uint8:
		return "uint8"
	case Type_Which_uint16:
		return "uint16"
	case Type_Which_uint32:
		return "uint32"
	case Type_Which_uint64:
		return "uint64"
	case Type_Which_float32:
		return "float32"
	case Type_Which_float64:
		return "float64"
	case Type_Which_text:
		return "string"
	case Type_Which_data:
		return "[]byte"
	case Type_Which_enum:
		ni := findNode(t.Enum().TypeId())
		return ni.RemoteName(n)
	case Type_Which_structGroup:
		ni := findNode(t.StructGroup().TypeId())
		return ni.RemoteName(n)
	case Type_Which_interface:
		ni := findNode(t.Interface().TypeId())
		return ni.RemoteName(n)
	case Type_Which_anyPointer:
		return g_imports.capnp() + ".Pointer"
	case Type_Which_list:
		switch lt, _ := t.List().ElementType(); lt.Which() {
		case Type_Which_void:
			return g_imports.capnp() + ".VoidList"
		case Type_Which_bool:
			return g_imports.capnp() + ".BitList"
		case Type_Which_int8:
			return g_imports.capnp() + ".Int8List"
		case Type_Which_uint8:
			return g_imports.capnp() + ".UInt8List"
		case Type_Which_int16:
			return g_imports.capnp() + ".Int16List"
		case Type_Which_uint16:
			return g_imports.capnp() + ".UInt16List"
		case Type_Which_int32:
			return g_imports.capnp() + ".Int32List"
		case Type_Which_uint32:
			return g_imports.capnp() + ".UInt32List"
		case Type_Which_int64:
			return g_imports.capnp() + ".Int64List"
		case Type_Which_uint64:
			return g_imports.capnp() + ".UInt64List"
		case Type_Which_float32:
			return g_imports.capnp() + ".Float32List"
		case Type_Which_float64:
			return g_imports.capnp() + ".Float64List"
		case Type_Which_text:
			return g_imports.capnp() + ".TextList"
		case Type_Which_data:
			return g_imports.capnp() + ".DataList"
		case Type_Which_enum:
			ni := findNode(lt.Enum().TypeId())
			return ni.RemoteName(n) + "_List"
		case Type_Which_structGroup:
			ni := findNode(lt.StructGroup().TypeId())
			return ni.RemoteName(n) + "_List"
		case Type_Which_anyPointer, Type_Which_list, Type_Which_interface:
			return g_imports.capnp() + ".PointerList"
		}
	}
	return ""
}

func intFieldDefaultMask(t Type, def Value) uint64 {
	if def.Which() == Value_Which_void {
		return 0
	}
	v := intValue(t, def)
	mask := uint64(1)<<intbits(t.Which()) - 1
	return uint64(v) & mask
}

func intValue(t Type, v Value) int64 {
	switch t.Which() {
	case Type_Which_int8:
		assert(v.Which() == Value_Which_int8, "expected int8 value")
		return int64(v.Int8())
	case Type_Which_int16:
		assert(v.Which() == Value_Which_int16, "expected int16 value")
		return int64(v.Int16())
	case Type_Which_int32:
		assert(v.Which() == Value_Which_int32, "expected int32 value")
		return int64(v.Int32())
	case Type_Which_int64:
		assert(v.Which() == Value_Which_int64, "expected int64 value")
		return v.Int64()
	}
	panic("unreachable")
}

func uintFieldDefault(t Type, def Value) uint64 {
	if def.Which() == Value_Which_void {
		return 0
	}
	return uintValue(t, def)
}

func uintValue(t Type, v Value) uint64 {
	switch t.Which() {
	case Type_Which_uint8:
		assert(v.Which() == Value_Which_uint8, "expected uint8 value")
		return uint64(v.Uint8())
	case Type_Which_uint16:
		assert(v.Which() == Value_Which_uint16, "expected uint16 value")
		return uint64(v.Uint16())
	case Type_Which_uint32:
		assert(v.Which() == Value_Which_uint32, "expected uint32 value")
		return uint64(v.Uint32())
	case Type_Which_uint64:
		assert(v.Which() == Value_Which_uint64, "expected uint64 value")
		return v.Uint64()
	}
	panic("unreachable")
}

func intbits(t Type_Which) uint {
	switch t {
	case Type_Which_uint8, Type_Which_int8:
		return 8
	case Type_Which_uint16, Type_Which_int16:
		return 16
	case Type_Which_uint32, Type_Which_int32:
		return 32
	case Type_Which_uint64, Type_Which_int64:
		return 64
	}
	return 0
}

func (n *node) codeOrderFields() []field {
	fields, _ := n.StructGroup().Fields()
	numFields := fields.Len()
	mbrs := make([]field, numFields)
	for i := 0; i < numFields; i++ {
		f := fields.At(i)
		fann, _ := f.Annotations()
		fname, _ := f.Name()
		fname = parseAnnotations(fann).Rename(fname)
		mbrs[f.CodeOrder()] = field{Field: f, Name: fname}
	}
	return mbrs
}

func (n *node) defineStructTypes(w io.Writer, baseNode *node) {
	assert(n.Which() == Node_Which_structGroup, "invalid struct node")

	nann, _ := n.Annotations()
	ann := parseAnnotations(nann)
	templates.ExecuteTemplate(w, "structTypes", structTypesParams{
		Node:        n,
		Annotations: ann,
		BaseNode:    baseNode,
	})

	for _, f := range n.codeOrderFields() {
		if f.Which() == Field_Which_group {
			findNode(f.Group().TypeId()).defineStructTypes(w, baseNode)
		}
	}
}

func (n *node) defineStructEnums(w io.Writer) {
	assert(n.Which() == Node_Which_structGroup, "invalid struct node")
	fields := n.codeOrderFields()
	members := make([]field, 0, len(fields))
	es := make(enumString, 0, len(fields))
	for _, f := range fields {
		if f.DiscriminantValue() != Field_noDiscriminant {
			members = append(members, f)
			es = append(es, f.Name)
		}
	}
	if n.StructGroup().DiscriminantCount() > 0 {
		templates.ExecuteTemplate(w, "structEnums", structEnumsParams{
			Node:       n,
			Fields:     members,
			EnumString: es,
		})
	}
	for _, f := range fields {
		if f.Which() == Field_Which_group {
			findNode(f.Group().TypeId()).defineStructEnums(w)
		}
	}
}

func (n *node) defineStructFuncs(w io.Writer) {
	assert(n.Which() == Node_Which_structGroup, "invalid struct node")

	templates.ExecuteTemplate(w, "structFuncs", structFuncsParams{
		Node: n,
	})

	for _, f := range n.codeOrderFields() {
		switch f.Which() {
		case Field_Which_slot:
			n.defineField(w, f)
		case Field_Which_group:
			g := findNode(f.Group().TypeId())
			templates.ExecuteTemplate(w, "structGroup", structGroupParams{
				Node:  n,
				Group: g,
				Field: f,
			})
			g.defineStructFuncs(w)
		}
	}
}

func (n *node) ObjectSize() string {
	assert(n.Which() == Node_Which_structGroup, "ObjectSize for invalid struct node")
	return fmt.Sprintf("%s.ObjectSize{DataSize: %d, PointerCount: %d}", g_imports.capnp(), int(n.StructGroup().DataWordCount())*8, n.StructGroup().PointerCount())
}

func (n *node) defineNewStructFunc(w io.Writer) {
	assert(n.Which() == Node_Which_structGroup, "invalid struct node")

	templates.ExecuteTemplate(w, "newStructFunc", newStructParams{
		Node: n,
	})
}

func (n *node) defineStructList(w io.Writer) {
	assert(n.Which() == Node_Which_structGroup, "invalid struct node")

	templates.ExecuteTemplate(w, "structList", structListParams{
		Node: n,
	})
}

func (n *node) defineStructPromise(w io.Writer) {
	templates.ExecuteTemplate(w, "promise", promiseTemplateParams{
		Node:   n,
		Fields: n.codeOrderFields(),
	})

	for _, f := range n.codeOrderFields() {
		switch f.Which() {
		case Field_Which_slot:
			t, _ := f.Slot().Type()
			if tw := t.Which(); tw == Type_Which_structGroup || tw == Type_Which_interface || tw == Type_Which_anyPointer {
				n.definePromiseField(w, f)
			}
		case Field_Which_group:
			g := findNode(f.Group().TypeId())
			templates.ExecuteTemplate(w, "promiseGroup", promiseGroupTemplateParams{
				Node:  n,
				Field: f,
				Group: g,
			})
			g.defineStructPromise(w)
		}
	}
}

func (n *node) definePromiseField(w io.Writer, f field) {
	slot := f.Slot()
	switch t, _ := slot.Type(); t.Which() {
	case Type_Which_structGroup:
		ni := findNode(t.StructGroup().TypeId())
		params := promiseFieldStructTemplateParams{
			Node:   n,
			Field:  f,
			Struct: ni,
		}
		if def, _ := slot.DefaultValue(); def.Which() == Value_Which_structField {
			if sf, _ := def.StructField(); capnp.HasData(sf) {
				params.Default = copyData(sf)
			}
		}
		templates.ExecuteTemplate(w, "promiseFieldStruct", params)
	case Type_Which_anyPointer:
		templates.ExecuteTemplate(w, "promiseFieldAnyPointer", promiseFieldAnyPointerTemplateParams{
			Node:  n,
			Field: f,
		})
	case Type_Which_interface:
		templates.ExecuteTemplate(w, "promiseFieldInterface", promiseFieldInterfaceTemplateParams{
			Node:      n,
			Field:     f,
			Interface: findNode(t.Interface().TypeId()),
		})
	}
}

type interfaceMethod struct {
	Method
	Interface    *node
	ID           int
	Name         string
	OriginalName string
	Params       *node
	Results      *node
}

func (n *node) methodSet(methods []interfaceMethod) []interfaceMethod {
	ms, _ := n.Interface().Methods()
	for i := 0; i < ms.Len(); i++ {
		m := ms.At(i)
		mname, _ := m.Name()
		mann, _ := m.Annotations()
		methods = append(methods, interfaceMethod{
			Method:       m,
			Interface:    n,
			ID:           i,
			OriginalName: mname,
			Name:         parseAnnotations(mann).Rename(mname),
			Params:       findNode(m.ParamStructType()),
			Results:      findNode(m.ResultStructType()),
		})
	}
	// TODO(light): sort added methods by code order

	supers, _ := n.Interface().Superclasses()
	for i := 0; i < supers.Len(); i++ {
		s := supers.At(i)
		methods = findNode(s.Id()).methodSet(methods)
	}
	return methods
}

func (n *node) defineInterfaceClient(w io.Writer) {
	m := n.methodSet(nil)
	nann, _ := n.Annotations()
	templates.ExecuteTemplate(w, "interfaceClient", interfaceClientTemplateParams{
		Node:        n,
		Annotations: parseAnnotations(nann),
		Methods:     m,
	})
}

func (n *node) defineInterfaceServer(w io.Writer) {
	m := n.methodSet(nil)
	nann, _ := n.Annotations()
	templates.ExecuteTemplate(w, "interfaceServer", interfaceServerTemplateParams{
		Node:        n,
		Annotations: parseAnnotations(nann),
		Methods:     m,
	})
}

type enumString []string

func (es enumString) ValueString() string {
	return strings.Join([]string(es), "")
}

func (es enumString) SliceFor(i int) string {
	n := 0
	for _, v := range es[:i] {
		n += len(v)
	}
	return fmt.Sprintf("[%d:%d]", n, n+len(es[i]))
}

func generateFile(reqf CodeGeneratorRequest_RequestedFile) (generr error) {
	defer func() {
		e := recover()
		if ae, ok := e.(assertionError); ok {
			generr = ae
		} else if e != nil {
			panic(e)
		}
	}()

	f := findNode(reqf.Id())
	var buf bytes.Buffer
	g_imports.init()
	g_segment = make([]byte, 0, 4096)
	g_bufname = fmt.Sprintf("x_%x", f.Id())

	for _, n := range f.nodes {
		if n.Which() == Node_Which_annotation {
			n.defineAnnotation(&buf)
		}
	}

	defineConstNodes(&buf, f.nodes)

	for _, n := range f.nodes {
		switch n.Which() {
		case Node_Which_enum:
			n.defineEnum(&buf)
		case Node_Which_structGroup:
			if !n.StructGroup().IsGroup() {
				n.defineStructTypes(&buf, n)
				n.defineStructEnums(&buf)
				n.defineNewStructFunc(&buf)
				n.defineStructFuncs(&buf)
				n.defineStructList(&buf)
				if *genPromises {
					n.defineStructPromise(&buf)
				}
			}
		case Node_Which_interface:
			n.defineInterfaceClient(&buf)
			n.defineInterfaceServer(&buf)
		}
	}

	fname, _ := reqf.Filename()
	if f.pkg == "" {
		return fmt.Errorf("missing package annotation for %s", fname)
	}

	if dirPath, _ := filepath.Split(fname); dirPath != "" {
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	var unformatted bytes.Buffer
	fmt.Fprintf(&unformatted, "package %s\n\n", f.pkg)
	fmt.Fprintf(&unformatted, "// AUTO GENERATED - DO NOT EDIT\n\n")
	fmt.Fprintf(&unformatted, "import (\n")
	for _, imp := range g_imports.usedImports() {
		fmt.Fprintf(&unformatted, "%v\n", imp)
	}
	fmt.Fprintf(&unformatted, ")\n")
	unformatted.Write(buf.Bytes())
	if len(g_segment) > 0 {
		fmt.Fprintf(&unformatted, "var %s = []byte{", g_bufname)
		for i, b := range g_segment {
			if i%8 == 0 {
				fmt.Fprintf(&unformatted, "\n")
			}
			fmt.Fprintf(&unformatted, "%d,", b)
		}
		fmt.Fprintf(&unformatted, "\n}\n")
	}
	formatted, err := format.Source(unformatted.Bytes())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't format generated code:", err)
		formatted = unformatted.Bytes()
	}

	file, err := os.Create(fname + ".go")
	if err != nil {
		return err
	}
	_, werr := file.Write(formatted)
	cerr := file.Close()
	if werr != nil {
		return err
	}
	if cerr != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()

	msg, err := capnp.NewDecoder(os.Stdin).Decode()
	if err != nil {
		fmt.Fprintln(os.Stderr, "capnpc-go: Reading input:", err)
		os.Exit(1)
	}

	req, err := ReadRootCodeGeneratorRequest(msg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "capnpc-go: Reading input:", err)
		os.Exit(1)
	}
	allfiles := []*node{}

	nodes, _ := req.Nodes()
	for i := 0; i < nodes.Len(); i++ {
		ni := nodes.At(i)
		n := &node{Node: ni}
		g_nodes[n.Id()] = n

		if n.Which() == Node_Which_file {
			allfiles = append(allfiles, n)
		}
	}

	for _, f := range allfiles {
		fann, _ := f.Annotations()
		ann := parseAnnotations(fann)
		f.pkg = ann.Package
		f.imp = ann.Import
		nnodes, _ := f.NestedNodes()
		for i := 0; i < nnodes.Len(); i++ {
			nn := nnodes.At(i)
			if ni := g_nodes[nn.Id()]; ni != nil {
				nname, _ := nn.Name()
				ni.resolveName("", nname, f)
			}
		}
	}

	success := true
	reqFiles, _ := req.RequestedFiles()
	for i := 0; i < reqFiles.Len(); i++ {
		reqf := reqFiles.At(i)
		fname, _ := reqf.Filename()
		err := generateFile(reqf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "capnpc-go: Generating Go for %s failed: %v\n", fname, err)
			success = false
		}
	}
	if !success {
		os.Exit(1)
	}
}
