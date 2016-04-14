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
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"zombiezen.com/go/capnproto2"
	"zombiezen.com/go/capnproto2/internal/schema"
)

var (
	genPromises = flag.Bool("promises", true, "generate code for promises")
)

const (
	go_capnproto_import = "zombiezen.com/go/capnproto2"
	server_import       = go_capnproto_import + "/server"
	context_import      = "golang.org/x/net/context"
)

type nodeMap map[uint64]*node

func (m nodeMap) mustFind(id uint64) *node {
	n := m[id]
	assert(n != nil, "could not find node %#x\n", id)
	return n
}

type generator struct {
	buf     bytes.Buffer
	fileID  uint64
	nodes   nodeMap
	imports imports
	data    staticData
}

func newGenerator(fileID uint64, nodes nodeMap) *generator {
	g := &generator{
		fileID: fileID,
		nodes:  nodes,
	}
	g.imports.init()
	g.data.init(fileID)
	return g
}

// Basename returns the name of the schema file with the directory name removed.
func (g *generator) Basename() string {
	f := g.nodes.mustFind(g.fileID)
	n, _ := f.DisplayName()
	return filepath.Base(n)
}

func (g *generator) Imports() *imports {
	return &g.imports
}

func (g *generator) Capnp() string {
	return g.imports.Capnp()
}

// generate produces unformatted Go source code from the nodes defined in it.
func (g *generator) generate(pkg string) []byte {
	var out bytes.Buffer
	fmt.Fprintf(&out, "package %s\n\n", pkg)
	out.WriteString("// AUTO GENERATED - DO NOT EDIT\n\n")
	out.WriteString("import (\n")
	for _, imp := range g.imports.usedImports() {
		fmt.Fprintf(&out, "%v\n", imp)
	}
	out.WriteString(")\n")
	out.Write(g.buf.Bytes())
	if len(g.data.buf) > 0 {
		fmt.Fprintf(&out, "var %s = []byte{", g.data.name)
		for i, b := range g.data.buf {
			if i%8 == 0 {
				out.WriteByte('\n')
			}
			fmt.Fprintf(&out, "%d,", b)
		}
		fmt.Fprintf(&out, "\n}\n")
	}
	return out.Bytes()
}

func (g *generator) remoteScope(n, from *node) string {
	displayName, _ := n.DisplayName()
	fromDisplayName, _ := from.DisplayName()
	assert(n.pkg != "", "missing package declaration for %s", displayName)
	assert(n.imp != "", "missing import declaration for %s", displayName)
	assert(from.imp != "", "missing import declaration for %s", fromDisplayName)

	if n.imp == from.imp {
		return ""
	} else {
		name := g.imports.add(importSpec{path: n.imp, name: n.pkg})
		return name + "."
	}
}

func (g *generator) RemoteNew(n, from *node) string {
	return g.remoteScope(n, from) + "New" + n.Name
}

func (g *generator) RemoteName(n, from *node) string {
	return g.remoteScope(n, from) + n.Name
}

var templates = template.New("").Funcs(template.FuncMap{
	"title": strings.Title,
	"hasDiscriminant": func(f field) bool {
		return f.DiscriminantValue() != schema.Field_noDiscriminant
	},
	"discriminantOffset": func(n *node) uint32 {
		return n.StructGroup().DiscriminantOffset() * 2
	},
})

type node struct {
	schema.Node
	pkg   string
	imp   string
	nodes []*node
	Name  string
}

type field struct {
	schema.Field
	Name string
}

func assert(chk bool, format string, a ...interface{}) {
	if !chk {
		panic(assertionError(fmt.Sprintf(format, a...)))
	}
}

func assertSuccess(err error) {
	if err != nil {
		panic(assertionError(err.Error() + "\n"))
	}
}

type assertionError string

func (ae assertionError) Error() string {
	return string(ae)
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

func parseAnnotations(list schema.Annotation_List) *annotations {
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

// resolveName is called as part of building up a node map to populate the name field of n.
func resolveName(nodes nodeMap, n *node, base, name string, file *node) {
	na, _ := n.Annotations()
	name = parseAnnotations(na).Rename(name)
	if base == "" {
		n.Name = strings.Title(name)
	} else {
		n.Name = base + "_" + name
	}
	n.pkg = file.pkg
	n.imp = file.imp

	if n.Which() != schema.Node_Which_structGroup || !n.StructGroup().IsGroup() {
		file.nodes = append(file.nodes, n)
	}

	nnodes, _ := n.NestedNodes()
	for i := 0; i < nnodes.Len(); i++ {
		nn := nnodes.At(i)
		if ni := nodes[nn.Id()]; ni != nil {
			nname, _ := nn.Name()
			resolveName(nodes, ni, n.Name, nname, file)
		}
	}

	if n.Which() == schema.Node_Which_structGroup {
		fields, _ := n.StructGroup().Fields()
		for i := 0; i < fields.Len(); i++ {
			f := fields.At(i)
			if f.Which() == schema.Field_Which_group {
				fa, _ := f.Annotations()
				fname, _ := f.Name()
				fname = parseAnnotations(fa).Rename(fname)
				resolveName(nodes, nodes.mustFind(f.Group().TypeId()), n.Name, fname, file)
			}
		}
	} else if n.Which() == schema.Node_Which_interface {
		m, _ := n.Interface().Methods()
		for i := 0; i < m.Len(); i++ {
			mm := m.At(i)
			mname, _ := mm.Name()
			mann, _ := mm.Annotations()
			mname = parseAnnotations(mann).Rename(mname)
			base := n.Name + "_" + mname
			if p := nodes.mustFind(mm.ParamStructType()); p.ScopeId() == 0 {
				resolveName(nodes, p, base, "Params", file)
			}
			if r := nodes.mustFind(mm.ResultStructType()); r.ScopeId() == 0 {
				resolveName(nodes, r, base, "Results", file)
			}
		}
	}
}

type enumval struct {
	schema.Enumerant
	Name   string
	Val    int
	Tag    string
	parent *node
}

func makeEnumval(enum *node, i int, e schema.Enumerant) enumval {
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

func (g *generator) defineEnum(n *node) {
	es, _ := n.Enum().Enumerants()
	ev := make([]enumval, es.Len())
	for i := 0; i < es.Len(); i++ {
		e := es.At(i)
		ev[e.CodeOrder()] = makeEnumval(n, i, e)
	}
	nann, _ := n.Annotations()
	templates.ExecuteTemplate(&g.buf, "enum", enumParams{
		G:           g,
		Node:        n,
		Annotations: parseAnnotations(nann),
		EnumValues:  ev,
	})
}

// Value formats a value from a schema (like a field default) as Go source.
func (g *generator) Value(n *node, t schema.Type, v schema.Value) string {
	switch t.Which() {
	case schema.Type_Which_void:
		return "struct{}{}"

	case schema.Type_Which_interface:
		// The only statically representable interface value is null.
		return g.imports.Capnp() + ".Client(nil)"

	case schema.Type_Which_bool:
		assert(v.Which() == schema.Value_Which_bool, "expected bool value")
		if v.Bool() {
			return "true"
		} else {
			return "false"
		}

	case schema.Type_Which_uint8, schema.Type_Which_uint16, schema.Type_Which_uint32, schema.Type_Which_uint64:
		return fmt.Sprintf("uint%d(%d)", intbits(t.Which()), uintValue(t, v))

	case schema.Type_Which_int8, schema.Type_Which_int16, schema.Type_Which_int32, schema.Type_Which_int64:
		return fmt.Sprintf("int%d(%d)", intbits(t.Which()), intValue(t, v))

	case schema.Type_Which_float32:
		assert(v.Which() == schema.Value_Which_float32, "expected float32 value")
		return fmt.Sprintf("%s.Float32frombits(0x%x)", g.imports.Math(), math.Float32bits(v.Float32()))

	case schema.Type_Which_float64:
		assert(v.Which() == schema.Value_Which_float64, "expected float64 value")
		return fmt.Sprintf("%s.Float64frombits(0x%x)", g.imports.Math(), math.Float64bits(v.Float64()))

	case schema.Type_Which_text:
		assert(v.Which() == schema.Value_Which_text, "expected text value")
		text, _ := v.Text()
		return strconv.Quote(text)

	case schema.Type_Which_data:
		assert(v.Which() == schema.Value_Which_data, "expected data value")
		var buf bytes.Buffer
		buf.WriteString("[]byte{")
		data, _ := v.Data()
		for i, b := range data {
			if i > 0 {
				buf.WriteString(", ")
			}
			fmt.Fprintf(&buf, "%d", b)
		}
		buf.WriteString("}")
		return buf.String()

	case schema.Type_Which_enum:
		assert(v.Which() == schema.Value_Which_enum, "expected enum value")
		en := g.nodes.mustFind(t.Enum().TypeId())
		assert(en.Which() == schema.Node_Which_enum, "expected enum type ID")
		enums, _ := en.Enum().Enumerants()
		if val := int(v.Enum()); val >= enums.Len() {
			return fmt.Sprintf("%s(%d)", g.RemoteName(en, n), val)
		} else {
			ev := makeEnumval(en, val, enums.At(val))
			return g.remoteScope(en, n) + ev.FullName()
		}

	case schema.Type_Which_structGroup:
		assert(v.Which() == schema.Value_Which_structField, "expected struct value")
		data, _ := v.StructFieldPtr()
		var buf bytes.Buffer
		templates.ExecuteTemplate(&buf, "structValue", structValueTemplateParams{
			G:     g,
			Node:  n,
			Typ:   g.nodes.mustFind(t.StructGroup().TypeId()),
			Value: g.data.copyData(data),
		})
		return buf.String()

	case schema.Type_Which_anyPointer:
		assert(v.Which() == schema.Value_Which_anyPointer, "expected pointer value")
		data, _ := v.AnyPointerPtr()
		var buf bytes.Buffer
		templates.ExecuteTemplate(&buf, "pointerValue", structValueTemplateParams{
			G:     g,
			Value: g.data.copyData(data),
		})
		return buf.String()

	case schema.Type_Which_list:
		assert(v.Which() == schema.Value_Which_list, "expected list value")
		data, _ := v.ListPtr()
		var buf bytes.Buffer
		templates.ExecuteTemplate(&buf, "listValue", listValueTemplateParams{
			G:     g,
			Typ:   g.fieldType(n, t, new(annotations)),
			Value: g.data.copyData(data),
		})
		return buf.String()
	default:
		panic(assertionError("unreachable"))
	}
}

func (g *generator) defineAnnotation(n *node) {
	templates.ExecuteTemplate(&g.buf, "annotation", annotationParams{
		G:    g,
		Node: n,
	})
}

func isGoConstType(t schema.Type) bool {
	w := t.Which()
	return w == schema.Type_Which_bool ||
		w == schema.Type_Which_int8 ||
		w == schema.Type_Which_uint8 ||
		w == schema.Type_Which_int16 ||
		w == schema.Type_Which_uint16 ||
		w == schema.Type_Which_int32 ||
		w == schema.Type_Which_uint32 ||
		w == schema.Type_Which_int64 ||
		w == schema.Type_Which_uint64 ||
		w == schema.Type_Which_text ||
		w == schema.Type_Which_enum
}

func (g *generator) defineConstNodes(nodes []*node) {
	constNodes := make([]*node, 0, len(nodes))
	for _, n := range nodes {
		if n.Which() != schema.Node_Which_const {
			continue
		}
		t, _ := n.Const().Type()
		if isGoConstType(t) {
			constNodes = append(constNodes, n)
		}
	}
	nc := len(constNodes)
	for _, n := range nodes {
		if n.Which() != schema.Node_Which_const {
			continue
		}
		t, _ := n.Const().Type()
		if !isGoConstType(t) {
			constNodes = append(constNodes, n)
		}
	}
	err := templates.ExecuteTemplate(&g.buf, "constants", constantsParams{
		G:      g,
		Consts: constNodes[:nc],
		Vars:   constNodes[nc:],
	})
	assertSuccess(err)
}

func (g *generator) defineField(n *node, f field) {
	fann, _ := f.Annotations()
	ann := parseAnnotations(fann)
	t, _ := f.Slot().Type()
	def, _ := f.Slot().DefaultValue()
	params := structFieldParams{
		G:           g,
		Node:        n,
		Field:       f,
		Annotations: ann,
		FieldType:   g.fieldType(n, t, ann),
	}
	switch t.Which() {
	case schema.Type_Which_void:
		templates.ExecuteTemplate(&g.buf, "structVoidField", params)
	case schema.Type_Which_bool:
		assert(def.Which() == schema.Value_Which_void || def.Which() == schema.Value_Which_bool, "expected bool default")
		templates.ExecuteTemplate(&g.buf, "structBoolField", structBoolFieldParams{
			structFieldParams: params,
			Default:           def.Which() == schema.Value_Which_bool && def.Bool(),
		})

	case schema.Type_Which_uint8, schema.Type_Which_uint16, schema.Type_Which_uint32, schema.Type_Which_uint64:
		templates.ExecuteTemplate(&g.buf, "structUintField", structUintFieldParams{
			structFieldParams: params,
			Bits:              intbits(t.Which()),
			Default:           uintFieldDefault(t, def),
		})

	case schema.Type_Which_int8, schema.Type_Which_int16, schema.Type_Which_int32, schema.Type_Which_int64:
		templates.ExecuteTemplate(&g.buf, "structIntField", structIntFieldParams{
			structUintFieldParams: structUintFieldParams{
				structFieldParams: params,
				Bits:              intbits(t.Which()),
				Default:           uint64(intFieldDefaultMask(t, def)),
			},
		})

	case schema.Type_Which_enum:
		assert(def.Which() == schema.Value_Which_void || def.Which() == schema.Value_Which_enum, "expected enum default")
		ni := g.nodes.mustFind(t.Enum().TypeId())
		var d uint64
		if def.Which() == schema.Value_Which_enum {
			d = uint64(def.Enum())
		}
		templates.ExecuteTemplate(&g.buf, "structIntField", structIntFieldParams{
			structUintFieldParams: structUintFieldParams{
				structFieldParams: params,
				Bits:              16,
				Default:           d,
			},
			EnumName: g.RemoteName(ni, n),
		})
	case schema.Type_Which_float32:
		assert(def.Which() == schema.Value_Which_void || def.Which() == schema.Value_Which_float32, "expected float32 default")
		var d uint64
		if def.Which() == schema.Value_Which_float32 && def.Float32() != 0 {
			d = uint64(math.Float32bits(def.Float32()))
		}
		templates.ExecuteTemplate(&g.buf, "structFloatField", structUintFieldParams{
			structFieldParams: params,
			Bits:              32,
			Default:           d,
		})

	case schema.Type_Which_float64:
		assert(def.Which() == schema.Value_Which_void || def.Which() == schema.Value_Which_float64, "expected float64 default")
		var d uint64
		if def.Which() == schema.Value_Which_float64 && def.Float64() != 0 {
			d = math.Float64bits(def.Float64())
		}
		templates.ExecuteTemplate(&g.buf, "structFloatField", structUintFieldParams{
			structFieldParams: params,
			Bits:              64,
			Default:           d,
		})

	case schema.Type_Which_text:
		assert(def.Which() == schema.Value_Which_void || def.Which() == schema.Value_Which_text, "expected text default")
		var d string
		if def.Which() == schema.Value_Which_text {
			d, _ = def.Text()
		}
		templates.ExecuteTemplate(&g.buf, "structTextField", structTextFieldParams{
			structFieldParams: params,
			Default:           d,
		})

	case schema.Type_Which_data:
		assert(def.Which() == schema.Value_Which_void || def.Which() == schema.Value_Which_data, "expected data default")
		var d []byte
		if def.Which() == schema.Value_Which_data {
			d, _ = def.Data()
		}
		templates.ExecuteTemplate(&g.buf, "structDataField", structDataFieldParams{
			structFieldParams: params,
			Default:           d,
		})

	case schema.Type_Which_structGroup:
		assert(def.Which() == schema.Value_Which_void || def.Which() == schema.Value_Which_structField, "expected struct default")
		var defref staticDataRef
		if def.Which() == schema.Value_Which_structField {
			if sf, _ := def.StructFieldPtr(); sf.IsValid() {
				defref = g.data.copyData(sf)
			}
		}
		templates.ExecuteTemplate(&g.buf, "structStructField", structObjectFieldParams{
			structFieldParams: params,
			TypeNode:          g.nodes.mustFind(t.StructGroup().TypeId()),
			Default:           defref,
		})

	case schema.Type_Which_anyPointer:
		assert(def.Which() == schema.Value_Which_void || def.Which() == schema.Value_Which_anyPointer, "expected object default")
		var defref staticDataRef
		if def.Which() == schema.Value_Which_anyPointer {
			if p, _ := def.AnyPointerPtr(); p.IsValid() {
				defref = g.data.copyData(p)
			}
		}
		templates.ExecuteTemplate(&g.buf, "structPointerField", structObjectFieldParams{
			structFieldParams: params,
			Default:           defref,
		})

	case schema.Type_Which_list:
		assert(def.Which() == schema.Value_Which_void || def.Which() == schema.Value_Which_list, "expected list default")
		var defref staticDataRef
		if def.Which() == schema.Value_Which_list {
			if l, _ := def.ListPtr(); l.IsValid() {
				defref = g.data.copyData(l)
			}
		}
		templates.ExecuteTemplate(&g.buf, "structListField", structObjectFieldParams{
			structFieldParams: params,
			Default:           defref,
		})

	case schema.Type_Which_interface:
		templates.ExecuteTemplate(&g.buf, "structInterfaceField", params)
	}
}

func (g *generator) fieldType(n *node, t schema.Type, ann *annotations) string {
	switch t.Which() {
	case schema.Type_Which_bool:
		return "bool"
	case schema.Type_Which_int8:
		return "int8"
	case schema.Type_Which_int16:
		return "int16"
	case schema.Type_Which_int32:
		return "int32"
	case schema.Type_Which_int64:
		return "int64"
	case schema.Type_Which_uint8:
		return "uint8"
	case schema.Type_Which_uint16:
		return "uint16"
	case schema.Type_Which_uint32:
		return "uint32"
	case schema.Type_Which_uint64:
		return "uint64"
	case schema.Type_Which_float32:
		return "float32"
	case schema.Type_Which_float64:
		return "float64"
	case schema.Type_Which_text:
		return "string"
	case schema.Type_Which_data:
		return "[]byte"
	case schema.Type_Which_enum:
		ni := g.nodes.mustFind(t.Enum().TypeId())
		return g.RemoteName(ni, n)
	case schema.Type_Which_structGroup:
		ni := g.nodes.mustFind(t.StructGroup().TypeId())
		return g.RemoteName(ni, n)
	case schema.Type_Which_interface:
		ni := g.nodes.mustFind(t.Interface().TypeId())
		return g.RemoteName(ni, n)
	case schema.Type_Which_anyPointer:
		return g.imports.Capnp() + ".Pointer"
	case schema.Type_Which_list:
		switch lt, _ := t.List().ElementType(); lt.Which() {
		case schema.Type_Which_void:
			return g.imports.Capnp() + ".VoidList"
		case schema.Type_Which_bool:
			return g.imports.Capnp() + ".BitList"
		case schema.Type_Which_int8:
			return g.imports.Capnp() + ".Int8List"
		case schema.Type_Which_uint8:
			return g.imports.Capnp() + ".UInt8List"
		case schema.Type_Which_int16:
			return g.imports.Capnp() + ".Int16List"
		case schema.Type_Which_uint16:
			return g.imports.Capnp() + ".UInt16List"
		case schema.Type_Which_int32:
			return g.imports.Capnp() + ".Int32List"
		case schema.Type_Which_uint32:
			return g.imports.Capnp() + ".UInt32List"
		case schema.Type_Which_int64:
			return g.imports.Capnp() + ".Int64List"
		case schema.Type_Which_uint64:
			return g.imports.Capnp() + ".UInt64List"
		case schema.Type_Which_float32:
			return g.imports.Capnp() + ".Float32List"
		case schema.Type_Which_float64:
			return g.imports.Capnp() + ".Float64List"
		case schema.Type_Which_text:
			return g.imports.Capnp() + ".TextList"
		case schema.Type_Which_data:
			return g.imports.Capnp() + ".DataList"
		case schema.Type_Which_enum:
			ni := g.nodes.mustFind(lt.Enum().TypeId())
			return g.RemoteName(ni, n) + "_List"
		case schema.Type_Which_structGroup:
			ni := g.nodes.mustFind(lt.StructGroup().TypeId())
			return g.RemoteName(ni, n) + "_List"
		case schema.Type_Which_anyPointer, schema.Type_Which_list, schema.Type_Which_interface:
			return g.imports.Capnp() + ".PointerList"
		}
	}
	return ""
}

func intFieldDefaultMask(t schema.Type, def schema.Value) uint64 {
	if def.Which() == schema.Value_Which_void {
		return 0
	}
	v := intValue(t, def)
	mask := uint64(1)<<intbits(t.Which()) - 1
	return uint64(v) & mask
}

func intValue(t schema.Type, v schema.Value) int64 {
	switch t.Which() {
	case schema.Type_Which_int8:
		assert(v.Which() == schema.Value_Which_int8, "expected int8 value")
		return int64(v.Int8())
	case schema.Type_Which_int16:
		assert(v.Which() == schema.Value_Which_int16, "expected int16 value")
		return int64(v.Int16())
	case schema.Type_Which_int32:
		assert(v.Which() == schema.Value_Which_int32, "expected int32 value")
		return int64(v.Int32())
	case schema.Type_Which_int64:
		assert(v.Which() == schema.Value_Which_int64, "expected int64 value")
		return v.Int64()
	}
	panic("unreachable")
}

func uintFieldDefault(t schema.Type, def schema.Value) uint64 {
	if def.Which() == schema.Value_Which_void {
		return 0
	}
	return uintValue(t, def)
}

func uintValue(t schema.Type, v schema.Value) uint64 {
	switch t.Which() {
	case schema.Type_Which_uint8:
		assert(v.Which() == schema.Value_Which_uint8, "expected uint8 value")
		return uint64(v.Uint8())
	case schema.Type_Which_uint16:
		assert(v.Which() == schema.Value_Which_uint16, "expected uint16 value")
		return uint64(v.Uint16())
	case schema.Type_Which_uint32:
		assert(v.Which() == schema.Value_Which_uint32, "expected uint32 value")
		return uint64(v.Uint32())
	case schema.Type_Which_uint64:
		assert(v.Which() == schema.Value_Which_uint64, "expected uint64 value")
		return v.Uint64()
	}
	panic("unreachable")
}

func intbits(t schema.Type_Which) uint {
	switch t {
	case schema.Type_Which_uint8, schema.Type_Which_int8:
		return 8
	case schema.Type_Which_uint16, schema.Type_Which_int16:
		return 16
	case schema.Type_Which_uint32, schema.Type_Which_int32:
		return 32
	case schema.Type_Which_uint64, schema.Type_Which_int64:
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

func (g *generator) defineStructTypes(n, baseNode *node) {
	assert(n.Which() == schema.Node_Which_structGroup, "invalid struct node")

	nann, _ := n.Annotations()
	ann := parseAnnotations(nann)
	templates.ExecuteTemplate(&g.buf, "structTypes", structTypesParams{
		G:           g,
		Node:        n,
		Annotations: ann,
		BaseNode:    baseNode,
	})

	for _, f := range n.codeOrderFields() {
		if f.Which() == schema.Field_Which_group {
			g.defineStructTypes(g.nodes.mustFind(f.Group().TypeId()), baseNode)
		}
	}
}

func (g *generator) defineStructEnums(n *node) {
	assert(n.Which() == schema.Node_Which_structGroup, "invalid struct node")
	fields := n.codeOrderFields()
	members := make([]field, 0, len(fields))
	es := make(enumString, 0, len(fields))
	for _, f := range fields {
		if f.DiscriminantValue() != schema.Field_noDiscriminant {
			members = append(members, f)
			es = append(es, f.Name)
		}
	}
	if n.StructGroup().DiscriminantCount() > 0 {
		templates.ExecuteTemplate(&g.buf, "structEnums", structEnumsParams{
			G:          g,
			Node:       n,
			Fields:     members,
			EnumString: es,
		})
	}
	for _, f := range fields {
		if f.Which() == schema.Field_Which_group {
			g.defineStructEnums(g.nodes.mustFind(f.Group().TypeId()))
		}
	}
}

func (g *generator) defineStructFuncs(n *node) {
	assert(n.Which() == schema.Node_Which_structGroup, "invalid struct node")

	templates.ExecuteTemplate(&g.buf, "structFuncs", structFuncsParams{
		G:    g,
		Node: n,
	})

	for _, f := range n.codeOrderFields() {
		switch f.Which() {
		case schema.Field_Which_slot:
			g.defineField(n, f)
		case schema.Field_Which_group:
			grp := g.nodes.mustFind(f.Group().TypeId())
			templates.ExecuteTemplate(&g.buf, "structGroup", structGroupParams{
				G:     g,
				Node:  n,
				Group: grp,
				Field: f,
			})
			g.defineStructFuncs(grp)
		}
	}
}

func (g *generator) ObjectSize(n *node) string {
	assert(n.Which() == schema.Node_Which_structGroup, "ObjectSize for invalid struct node")
	return fmt.Sprintf("%s.ObjectSize{DataSize: %d, PointerCount: %d}", g.imports.Capnp(), int(n.StructGroup().DataWordCount())*8, n.StructGroup().PointerCount())
}

func (g *generator) defineNewStructFunc(n *node) {
	assert(n.Which() == schema.Node_Which_structGroup, "invalid struct node")

	templates.ExecuteTemplate(&g.buf, "newStructFunc", newStructParams{
		G:    g,
		Node: n,
	})
}

func (g *generator) defineStructList(n *node) {
	assert(n.Which() == schema.Node_Which_structGroup, "invalid struct node")

	templates.ExecuteTemplate(&g.buf, "structList", structListParams{
		G:    g,
		Node: n,
	})
}

func (g *generator) defineStructPromise(n *node) {
	templates.ExecuteTemplate(&g.buf, "promise", promiseTemplateParams{
		G:      g,
		Node:   n,
		Fields: n.codeOrderFields(),
	})

	for _, f := range n.codeOrderFields() {
		switch f.Which() {
		case schema.Field_Which_slot:
			t, _ := f.Slot().Type()
			if tw := t.Which(); tw == schema.Type_Which_structGroup || tw == schema.Type_Which_interface || tw == schema.Type_Which_anyPointer {
				g.definePromiseField(n, f)
			}
		case schema.Field_Which_group:
			grp := g.nodes.mustFind(f.Group().TypeId())
			templates.ExecuteTemplate(&g.buf, "promiseGroup", promiseGroupTemplateParams{
				G:     g,
				Node:  n,
				Field: f,
				Group: grp,
			})
			g.defineStructPromise(grp)
		}
	}
}

func (g *generator) definePromiseField(n *node, f field) {
	slot := f.Slot()
	switch t, _ := slot.Type(); t.Which() {
	case schema.Type_Which_structGroup:
		ni := g.nodes.mustFind(t.StructGroup().TypeId())
		params := promiseFieldStructTemplateParams{
			G:      g,
			Node:   n,
			Field:  f,
			Struct: ni,
		}
		if def, _ := slot.DefaultValue(); def.Which() == schema.Value_Which_structField {
			if sf, _ := def.StructFieldPtr(); sf.IsValid() {
				params.Default = g.data.copyData(sf)
			}
		}
		templates.ExecuteTemplate(&g.buf, "promiseFieldStruct", params)
	case schema.Type_Which_anyPointer:
		templates.ExecuteTemplate(&g.buf, "promiseFieldAnyPointer", promiseFieldAnyPointerTemplateParams{
			G:     g,
			Node:  n,
			Field: f,
		})
	case schema.Type_Which_interface:
		templates.ExecuteTemplate(&g.buf, "promiseFieldInterface", promiseFieldInterfaceTemplateParams{
			G:         g,
			Node:      n,
			Field:     f,
			Interface: g.nodes.mustFind(t.Interface().TypeId()),
		})
	}
}

type interfaceMethod struct {
	schema.Method
	Interface    *node
	ID           int
	Name         string
	OriginalName string
	Params       *node
	Results      *node
}

func methodSet(methods []interfaceMethod, n *node, nodes nodeMap) []interfaceMethod {
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
			Params:       nodes.mustFind(m.ParamStructType()),
			Results:      nodes.mustFind(m.ResultStructType()),
		})
	}
	// TODO(light): sort added methods by code order

	supers, _ := n.Interface().Superclasses()
	for i := 0; i < supers.Len(); i++ {
		s := supers.At(i)
		methods = methodSet(methods, nodes.mustFind(s.Id()), nodes)
	}
	return methods
}

func (g *generator) defineInterfaceClient(n *node) {
	m := methodSet(nil, n, g.nodes)
	nann, _ := n.Annotations()
	templates.ExecuteTemplate(&g.buf, "interfaceClient", interfaceClientTemplateParams{
		G:           g,
		Node:        n,
		Annotations: parseAnnotations(nann),
		Methods:     m,
	})
}

func (g *generator) defineInterfaceServer(n *node) {
	m := methodSet(nil, n, g.nodes)
	nann, _ := n.Annotations()
	templates.ExecuteTemplate(&g.buf, "interfaceServer", interfaceServerTemplateParams{
		G:           g,
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

func generateFile(reqf schema.CodeGeneratorRequest_RequestedFile, nodes nodeMap) (generr error) {
	defer func() {
		e := recover()
		if ae, ok := e.(assertionError); ok {
			generr = ae
		} else if e != nil {
			panic(e)
		}
	}()

	f := nodes.mustFind(reqf.Id())
	g := newGenerator(f.Id(), nodes)

	for _, n := range f.nodes {
		if n.Which() == schema.Node_Which_annotation {
			g.defineAnnotation(n)
		}
	}

	g.defineConstNodes(f.nodes)

	for _, n := range f.nodes {
		switch n.Which() {
		case schema.Node_Which_enum:
			g.defineEnum(n)
		case schema.Node_Which_structGroup:
			if !n.StructGroup().IsGroup() {
				g.defineStructTypes(n, n)
				g.defineStructEnums(n)
				g.defineNewStructFunc(n)
				g.defineStructFuncs(n)
				g.defineStructList(n)
				if *genPromises {
					g.defineStructPromise(n)
				}
			}
		case schema.Node_Which_interface:
			g.defineInterfaceClient(n)
			g.defineInterfaceServer(n)
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

	unformatted := g.generate(f.pkg)
	formatted, err := format.Source(unformatted)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't format generated code:", err)
		formatted = unformatted
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

	req, err := schema.ReadRootCodeGeneratorRequest(msg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "capnpc-go: Reading input:", err)
		os.Exit(1)
	}
	allfiles := []*node{}

	reqNodes, _ := req.Nodes()
	nodes := make(nodeMap, reqNodes.Len())
	for i := 0; i < reqNodes.Len(); i++ {
		ni := reqNodes.At(i)
		n := &node{Node: ni}
		nodes[n.Id()] = n

		if n.Which() == schema.Node_Which_file {
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
			if ni := nodes[nn.Id()]; ni != nil {
				nname, _ := nn.Name()
				resolveName(nodes, ni, "", nname, f)
			}
		}
	}

	success := true
	reqFiles, _ := req.RequestedFiles()
	for i := 0; i < reqFiles.Len(); i++ {
		reqf := reqFiles.At(i)
		fname, _ := reqf.Filename()
		err := generateFile(reqf, nodes)
		if err != nil {
			fmt.Fprintf(os.Stderr, "capnpc-go: Generating Go for %s failed: %v\n", fname, err)
			success = false
		}
	}
	if !success {
		os.Exit(1)
	}
}
