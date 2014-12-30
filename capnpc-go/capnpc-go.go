package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	C "github.com/glycerine/go-capnproto"
)

const (
	go_capnproto_import = "github.com/glycerine/go-capnproto"
	context_import      = "golang.org/x/net/context"
)

var (
	g_nodes   = make(map[uint64]*node)
	g_imports imports
	g_segment *C.Segment
	g_bufname string
)

type imports struct {
	specs []importSpec
	used  map[string]bool // keyed on import path
}

func (i *imports) init() {
	i.specs = nil
	i.used = make(map[string]bool)

	i.reserve(importSpec{path: go_capnproto_import, name: "C"})
	i.reserve(importSpec{path: context_import, name: "context"})

	i.reserve(importSpec{path: "bufio", name: "bufio"})
	i.reserve(importSpec{path: "bytes", name: "bytes"})
	i.reserve(importSpec{path: "io", name: "io"})
	i.reserve(importSpec{path: "json", name: "encoding/json"})
	i.reserve(importSpec{path: "math", name: "math"})
}

func (i *imports) capn() string {
	return i.add(importSpec{path: go_capnproto_import, name: "C"})
}

func (i *imports) context() string {
	return i.add(importSpec{path: context_import, name: "context"})
}

func (i *imports) math() string {
	return i.add(importSpec{path: "math", name: "math"})
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

func assert(chk bool, format string, a ...interface{}) {
	if !chk {
		panic(assertionError(fmt.Sprintf(format, a...)))
	}
}

type assertionError string

func (ae assertionError) Error() string {
	return string(ae)
}

func copyData(obj C.Object) int {
	r, off, err := g_segment.NewRoot()
	assert(err == nil, "%v\n", err)
	err = r.Set(0, obj)
	assert(err == nil, "%v\n", err)
	return off
}

// Tag types
const (
	defaultTag = iota
	noTag
	customTag
)

type annotations struct {
	Doc        string
	CustomType string
	Package    string
	Import     string
	TagType    int
	CustomTag  string
}

func parseAnnotations(list Annotation_List) *annotations {
	ann := new(annotations)
	for i, n := 0, list.Len(); i < n; i++ {
		a := list.At(i)
		switch a.Id() {
		case C.Doc:
			ann.Doc = a.Value().Text()
		case C.Customtype:
			ann.CustomType = a.Value().Text()
		case C.Package:
			ann.Package = a.Value().Text()
		case C.Import:
			ann.Import = a.Value().Text()
		case C.Tag:
			ann.TagType = customTag
			ann.CustomTag = a.Value().Text()
		case C.Notag:
			ann.TagType = noTag
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

func findNode(id uint64) *node {
	n := g_nodes[id]
	assert(n != nil, "could not find node 0x%x\n", id)
	return n
}

func (n *node) remoteScope(from *node) string {
	assert(n.pkg != "", "missing package declaration for %s", n.DisplayName())
	assert(n.imp != "", "missing import declaration for %s", n.DisplayName())
	assert(from.imp != "", "missing import declaration for %s", from.DisplayName())

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
	if base == "" {
		n.Name = strings.Title(name)
	} else {
		n.Name = base + "_" + name
	}
	n.pkg = file.pkg
	n.imp = file.imp

	if n.Which() != Node_Which_struct || !n.Struct().IsGroup() {
		file.nodes = append(file.nodes, n)
	}

	for i := 0; i < n.NestedNodes().Len(); i++ {
		nn := n.NestedNodes().At(i)
		if ni := g_nodes[nn.Id()]; ni != nil {
			ni.resolveName(n.Name, nn.Name(), file)
		}
	}

	if n.Which() == Node_Which_struct {
		for i := 0; i < n.Struct().Fields().Len(); i++ {
			f := n.Struct().Fields().At(i)
			if f.Which() == Field_Which_group {
				findNode(f.Group().TypeId()).resolveName(n.Name, f.Name(), file)
			}
		}
	} else if n.Which() == Node_Which_interface {
		for i, m := 0, n.Interface().Methods(); i < m.Len(); i++ {
			mm := m.At(i)
			base := n.Name + "_" + mm.Name()
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
	val    int
	tag    string
	parent *node
}

func (e *enumval) fullName() string {
	return e.parent.Name + "_" + e.Name()
}

func (n *node) defineEnum(w io.Writer) {
	ann := parseAnnotations(n.Annotations())
	if ann.Doc != "" {
		fmt.Fprintf(w, "// %s\n", ann.Doc)
	}
	fmt.Fprintf(w, "type %s uint16\n", n.Name)

	if es := n.Enum().Enumerants(); es.Len() > 0 {
		fmt.Fprintf(w, "const (\n")

		ev := make([]enumval, es.Len())
		for i := 0; i < es.Len(); i++ {
			e := es.At(i)
			ann := parseAnnotations(e.Annotations())
			t := ann.Tag(e.Name())
			ev[e.CodeOrder()] = enumval{e, i, t, n}
		}

		// not an iota, so type has to go on each line
		for _, e := range ev {
			fmt.Fprintf(w, "%s %s = %d\n", e.fullName(), n.Name, e.val)
		}

		fmt.Fprintf(w, ")\n")

		fmt.Fprintf(w, "func (c %s) String() string {\n", n.Name)
		fmt.Fprintf(w, "switch c {\n")
		for _, e := range ev {
			if e.tag != "" {
				fmt.Fprintf(w, "case %s: return \"%s\"\n", e.fullName(), e.tag)
			}
		}
		fmt.Fprintf(w, "default: return \"\"\n")
		fmt.Fprintf(w, "}\n}\n\n")

		fmt.Fprintf(w, "func %sFromString(c string) %s {\n", n.Name, n.Name)
		fmt.Fprintf(w, "switch c {\n")
		for _, e := range ev {
			if e.tag != "" {
				fmt.Fprintf(w, "case \"%s\": return %s\n", e.tag, e.fullName())
			}
		}
		fmt.Fprintf(w, "default: return 0\n")
		fmt.Fprintf(w, "}\n}\n")
	}

	c := g_imports.capn()
	fmt.Fprintf(w, "type %s_List %s.PointerList\n", n.Name, c)
	fmt.Fprintf(w, "func New%[1]s_List(s *%[2]s.Segment, sz int) %[1]s_List { return %[1]s_List(s.NewUInt16List(sz)) }\n", n.Name, c)
	fmt.Fprintf(w, "func (s %s_List) Len() int { return %s.UInt16List(s).Len() }\n", n.Name, c)
	fmt.Fprintf(w, "func (s %[1]s_List) At(i int) %[1]s { return %[1]s(%[2]s.UInt16List(s).At(i)) }\n", n.Name, c)
}

func (n *node) writeValue(w io.Writer, t Type, v Value) {
	switch t.Which() {
	case Type_Which_void:
		fmt.Fprintf(w, "%s.Void{}", g_imports.capn())

	case Type_Which_interface:
		// The only statically representable interface value is null.
		fmt.Fprintf(w, "%s.Promise(nil)", g_imports.capn())

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
		fmt.Fprintf(w, "%q", v.Text())

	case Type_Which_data:
		assert(v.Which() == Value_Which_data, "expected data value")
		fmt.Fprint(w, "[]byte{")
		for i, b := range v.Data() {
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
		enums := en.Enum().Enumerants()
		if val := int(v.Enum()); val >= enums.Len() {
			fmt.Fprintf(w, "%s(%d)", en.RemoteName(n), val)
		} else {
			ev := enumval{
				Enumerant: enums.At(val),
				parent:    en,
			}
			fmt.Fprintf(w, "%s%s", en.remoteScope(n), ev.fullName())
		}

	case Type_Which_struct:
		fmt.Fprintf(w, "%s(%s.Root(%d))", findNode(t.Struct().TypeId()).RemoteName(n), g_bufname, copyData(v.Struct()))

	case Type_Which_anyPointer:
		fmt.Fprintf(w, "%s.Root(%d)", g_bufname, copyData(v.AnyPointer()))

	case Type_Which_list:
		assert(v.Which() == Value_Which_list, "expected list value")

		if lt := t.List().ElementType(); lt.Which() == Type_Which_void {
			fmt.Fprintf(w, "make([]%s.Void, %d)", g_imports.capn(), v.List().ToVoidList().Len())
		} else {
			typ := n.fieldType(t, new(annotations))
			fmt.Fprintf(w, "%s(%s.Root(%d))", typ, g_bufname, copyData(v.List()))
		}
	}
}

func (n *node) defineAnnotation(w io.Writer) {
	fmt.Fprintf(w, "const %s = uint64(0x%x)\n", n.Name, n.Id())
}

func constIsVar(n *node) bool {
	switch n.Const().Type().Which() {
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
			n.writeValue(w, n.Const().Type(), n.Const().Value())
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
			n.writeValue(w, n.Const().Type(), n.Const().Value())
			fmt.Fprintf(w, "\n")
		}
	}

	if any {
		fmt.Fprintf(w, ")\n")
	}
}

func (n *node) defineField(w io.Writer, f Field) {
	t := f.Slot().Type()
	def := f.Slot().DefaultValue()
	off := f.Slot().Offset()

	var g, s bytes.Buffer

	fieldTitle := strings.Title(f.Name())
	settag := ""
	if f.DiscriminantValue() != Field_noDiscriminant {
		settag = fmt.Sprintf(" %s.Struct(s).Set16(%d, %d);", g_imports.capn(), n.Struct().DiscriminantOffset()*2, f.DiscriminantValue())
	}
	if t.Which() == Type_Which_void {
		if f.DiscriminantValue() != Field_noDiscriminant {
			fmt.Fprintf(&s, "func (s %s) Set%s() {%s }\n", n.Name, fieldTitle, settag)
			w.Write(s.Bytes())
		}
		return
	}

	ann := parseAnnotations(f.Annotations())
	if ann.Doc != "" {
		fmt.Fprintf(&g, "// %s\n", ann.Doc)
	}
	typ := n.fieldType(t, ann)
	if typ == "" {
		return
	}
	fmt.Fprintf(&g, "func (s %s) %s() %s { ", n.Name, fieldTitle, typ)
	fmt.Fprintf(&s, "func (s %s) Set%s(v %s) {%s ", n.Name, fieldTitle, typ, settag)

	switch t.Which() {
	case Type_Which_bool:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_bool, "expected bool default")
		c := g_imports.capn()
		if def.Which() == Value_Which_bool && def.Bool() {
			fmt.Fprintf(&g, "return !%s.Struct(s).Get1(%d) }\n", c, off)
			fmt.Fprintf(&s, "%s.Struct(s).Set1(%d, !v) }\n", c, off)
		} else {
			fmt.Fprintf(&g, "return %s.Struct(s).Get1(%d) }\n", c, off)
			fmt.Fprintf(&s, "%s.Struct(s).Set1(%d, v) }\n", c, off)
		}

	case Type_Which_uint8, Type_Which_uint16, Type_Which_uint32, Type_Which_uint64:
		n.defineUintField(&g, &s, intbits(t.Which()), off, uintFieldDefault(t, def))

	case Type_Which_int8, Type_Which_int16, Type_Which_int32, Type_Which_int64:
		n.defineIntField(&g, &s, intbits(t.Which()), off, intFieldDefault(t, def), "")

	case Type_Which_enum:
		ni := findNode(t.Enum().TypeId())
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_enum, "expected enum default")
		var d int64
		if def.Which() == Value_Which_enum {
			d = int64(def.Enum())
		}
		n.defineIntField(&g, &s, 16, off, d, ni.RemoteName(n))

	case Type_Which_float32:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_float32, "expected float32 default")
		c, m := g_imports.capn(), g_imports.math()
		fmt.Fprintf(&g, "return %s.Float32frombits(%s.Struct(s).Get32(%d)", m, c, off*4)
		fmt.Fprintf(&s, "%s.Struct(s).Set32(%d, %s.Float32bits(v)", c, off*4, m)
		if def.Which() == Value_Which_float64 && def.Float64() != 0 {
			bits := math.Float32bits(def.Float32())
			fmt.Fprintf(&g, " ^ 0x%x) }\n", bits)
			fmt.Fprintf(&s, " ^ 0x%x) }\n", bits)
		} else {
			fmt.Fprintf(&g, ") }\n")
			fmt.Fprintf(&s, ") }\n")
		}

	case Type_Which_float64:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_float64, "expected float64 default")
		c, m := g_imports.capn(), g_imports.math()
		fmt.Fprintf(&g, "return %s.Float64frombits(%s.Struct(s).Get64(%d)", m, c, off*8)
		fmt.Fprintf(&s, "%s.Struct(s).Set64(%d, %s.Float64bits(v)", c, off*8, m)
		if def.Which() == Value_Which_float64 && def.Float64() != 0 {
			bits := math.Float64bits(def.Float64())
			fmt.Fprintf(&g, " ^ 0x%x) }\n", bits)
			fmt.Fprintf(&s, " ^ 0x%x) }\n", bits)
		} else {
			fmt.Fprintf(&g, ") }\n")
			fmt.Fprintf(&s, ") }\n")
		}

	case Type_Which_text:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_text, "expected text default")
		c := g_imports.capn()
		fmt.Fprintf(&g, "return %s.Struct(s).GetObject(%d)", c, off)
		if def.Which() == Value_Which_text && def.Text() != "" {
			fmt.Fprintf(&g, ".ToTextDefault(%q) }\n", def.Text())
		} else {
			fmt.Fprintf(&g, ".ToText() }\n")
		}
		fmt.Fprintf(&s, "%s.Struct(s).SetObject(%d, s.Segment.NewText(v)) }\n", c, off)

	case Type_Which_data:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_data, "expected data default")
		c := g_imports.capn()
		if typ != "[]byte" {
			fmt.Fprintf(&g, "return %s(", typ)
		} else {
			fmt.Fprint(&g, "return ")
		}
		fmt.Fprintf(&g, "%s.Struct(s).GetObject(%d)", c, off)
		if def.Which() == Value_Which_data && len(def.Data()) > 0 {
			fmt.Fprint(&g, ".ToDataDefault([]byte{")
			for i, b := range def.Data() {
				if i > 0 {
					fmt.Fprint(&g, ", ")
				}
				fmt.Fprintf(&g, "%d", b)
			}
			fmt.Fprint(&g, "})")
		} else {
			fmt.Fprint(&g, ".ToData()")
		}
		if typ != "[]byte" {
			fmt.Fprint(&g, ")")
		}
		fmt.Fprint(&g, " }\n")

		fmt.Fprintf(&s, "%s.Struct(s).SetObject(%d, ", c, off)
		if ann.CustomType != "" {
			fmt.Fprintf(&s, "s.Segment.NewData([]byte(v))) }\n")
		} else {
			fmt.Fprintf(&s, "s.Segment.NewData(v)) }\n")
		}

	case Type_Which_struct:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_struct, "expected struct default")
		c := g_imports.capn()
		fmt.Fprintf(&g, "return %s(%s.Struct(s).GetObject(%d)", typ, c, off)
		if def.Which() == Value_Which_struct && def.Struct().HasData() {
			fmt.Fprintf(&g, ".ToStructDefault(%s, %d)) }\n", g_bufname, copyData(def.Struct()))
		} else {
			fmt.Fprint(&g, ".ToStruct()) }\n")
		}
		fmt.Fprintf(&s, "%s.Struct(s).SetObject(%d, %s.Object(v)) }\n", c, off, c)

	case Type_Which_anyPointer:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_anyPointer, "expected object default")
		c := g_imports.capn()
		fmt.Fprintf(&g, "return %s.Struct(s).GetObject(%d)", c, off)
		if def.Which() == Value_Which_anyPointer && def.AnyPointer().HasData() {
			fmt.Fprintf(&g, ".ToObjectDefault(%s, %d) }\n", g_bufname, copyData(def.AnyPointer()))
		} else {
			fmt.Fprint(&g, " }\n")
		}
		fmt.Fprintf(&s, "%s.Struct(s).SetObject(%d, v) }\n", c, off)

	case Type_Which_list:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_list, "expected list default")
		c := g_imports.capn()
		var ldef C.Object
		if def.Which() == Value_Which_list {
			ldef = def.List()
		}
		fmt.Fprintf(&g, "return %s(%s.Struct(s).GetObject(%d)", typ, c, off)
		if ldef.HasData() {
			fmt.Fprintf(&g, ".ToListDefault(%s, %d)) }\n", g_bufname, copyData(ldef))
		} else {
			fmt.Fprint(&g, ") }\n")
		}
		fmt.Fprintf(&s, "%s.Struct(s).SetObject(%d, %s.Object(v)) }\n", c, off, c)

	case Type_Which_interface:
		c := g_imports.capn()
		ni := findNode(t.Interface().TypeId())
		fmt.Fprintf(&g, "return %s(%s.Struct(s).GetObject(%d).ToInterface().Client()) }\n", ni.RemoteNew(n), c, off)
		fmt.Fprintf(&s, "ci := s.Segment.Message.AddCap(v.GenericClient()); %[1]s.Struct(s).SetObject(%[2]d, %[1]s.Object(s.Segment.NewInterface(ci))) }\n", c, off)
	}

	w.Write(g.Bytes())
	w.Write(s.Bytes())
}

func (n *node) fieldType(t Type, ann *annotations) string {
	customtype := ann.CustomType
	if customtype != "" {
		if i := strings.LastIndex(customtype, "."); i != -1 {
			pkg := g_imports.add(importSpec{path: customtype[:i]})
			customtype = pkg + "." + customtype[i+1:]
		}
	}
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
		if customtype != "" {
			return customtype
		}
		return "[]byte"
	case Type_Which_enum:
		ni := findNode(t.Enum().TypeId())
		return ni.RemoteName(n)
	case Type_Which_struct:
		ni := findNode(t.Struct().TypeId())
		return ni.RemoteName(n)
	case Type_Which_interface:
		ni := findNode(t.Interface().TypeId())
		return ni.RemoteName(n)
	case Type_Which_anyPointer:
		return g_imports.capn() + ".Object"
	case Type_Which_list:
		switch lt := t.List().ElementType(); lt.Which() {
		case Type_Which_void:
			return g_imports.capn() + ".VoidList"
		case Type_Which_bool:
			return g_imports.capn() + ".BitList"
		case Type_Which_int8:
			return g_imports.capn() + ".Int8List"
		case Type_Which_uint8:
			return g_imports.capn() + ".UInt8List"
		case Type_Which_int16:
			return g_imports.capn() + ".Int16List"
		case Type_Which_uint16:
			return g_imports.capn() + ".UInt16List"
		case Type_Which_int32:
			return g_imports.capn() + ".Int32List"
		case Type_Which_uint32:
			return g_imports.capn() + ".UInt32List"
		case Type_Which_int64:
			return g_imports.capn() + ".Int64List"
		case Type_Which_uint64:
			return g_imports.capn() + ".UInt64List"
		case Type_Which_float32:
			return g_imports.capn() + ".Float32List"
		case Type_Which_float64:
			return g_imports.capn() + ".Float64List"
		case Type_Which_text:
			return g_imports.capn() + ".TextList"
		case Type_Which_data:
			return g_imports.capn() + ".DataList"
		case Type_Which_enum:
			ni := findNode(lt.Enum().TypeId())
			return ni.RemoteName(n) + "_List"
		case Type_Which_struct:
			ni := findNode(lt.Struct().TypeId())
			return ni.RemoteName(n) + "_List"
		case Type_Which_anyPointer, Type_Which_list, Type_Which_interface:
			return g_imports.capn() + ".PointerList"
		}
	}
	return ""
}

func (n *node) defineUintField(g, s io.Writer, bits int, off uint32, def uint64) {
	c := g_imports.capn()
	off = off * uint32(bits/8)
	if def != 0 {
		fmt.Fprintf(g, "return %s.Struct(s).Get%d(%d) ^ %d }\n", c, bits, off, def)
		fmt.Fprintf(s, "%s.Struct(s).Set%d(%d, v^%d) }\n", c, bits, off, def)
	} else {
		fmt.Fprintf(g, "return %s.Struct(s).Get%d(%d) }\n", c, bits, off)
		fmt.Fprintf(s, "%s.Struct(s).Set%d(%d, v) }\n", c, bits, off)
	}
}

func (n *node) defineIntField(g, s io.Writer, bits int, off uint32, def int64, enum string) {
	c := g_imports.capn()
	off = off * uint32(bits/8)
	var rettype string
	if enum == "" {
		rettype = fmt.Sprintf("int%d", bits)
	} else {
		rettype = enum
	}

	if def != 0 {
		fmt.Fprintf(g, "return %s(%s.Struct(s).Get%d(%d)) ^ %d }\n", rettype, c, bits, off, def)
		fmt.Fprintf(s, "%s.Struct(s).Set%d(%d, uint%d(v^%d)) }\n", c, bits, off, bits, def)
	} else {
		fmt.Fprintf(g, "return %s(%s.Struct(s).Get%d(%d)) }\n", rettype, c, bits, off)
		fmt.Fprintf(s, "%s.Struct(s).Set%d(%d, uint%d(v)) }\n", c, bits, off, bits)
	}
}

func intFieldDefault(t Type, def Value) int64 {
	if def.Which() == Value_Which_void {
		return 0
	}
	return intValue(t, def)
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

func intbits(t Type_Which) int {
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

func (n *node) codeOrderFields() []Field {
	numFields := n.Struct().Fields().Len()
	mbrs := make([]Field, numFields)
	for i := 0; i < numFields; i++ {
		f := n.Struct().Fields().At(i)
		mbrs[f.CodeOrder()] = f
	}
	return mbrs
}

func (n *node) defineStructTypes(w io.Writer, baseNode *node) {
	assert(n.Which() == Node_Which_struct, "invalid struct node")

	ann := parseAnnotations(n.Annotations())
	if ann.Doc != "" {
		fmt.Fprintf(w, "// %s\n", ann.Doc)
	}
	if baseNode != nil {
		fmt.Fprintf(w, "type %s %s\n", n.Name, baseNode.Name)
	} else {
		fmt.Fprintf(w, "type %s %s.Struct\n", n.Name, g_imports.capn())
		baseNode = n
	}

	for _, f := range n.codeOrderFields() {
		if f.Which() == Field_Which_group {
			findNode(f.Group().TypeId()).defineStructTypes(w, baseNode)
		}
	}
}

func (n *node) defineStructEnums(w io.Writer) {
	assert(n.Which() == Node_Which_struct, "invalid struct node")

	if n.Struct().DiscriminantCount() > 0 {
		fmt.Fprintf(w, "type %s_Which uint16\n", n.Name)
		fmt.Fprintf(w, "const (\n")

		for _, f := range n.codeOrderFields() {
			if f.DiscriminantValue() == Field_noDiscriminant {
				// Non-union member
			} else {
				fmt.Fprintf(w, "%s_Which_%s %s_Which = %d\n", n.Name, f.Name(), n.Name, f.DiscriminantValue())
			}
		}
		fmt.Fprintf(w, ")\n")
	}

	for _, f := range n.codeOrderFields() {
		if f.Which() == Field_Which_group {
			findNode(f.Group().TypeId()).defineStructEnums(w)
		}
	}
}

func (n *node) defineStructFuncs(w io.Writer) {
	assert(n.Which() == Node_Which_struct, "invalid struct node")

	if n.Struct().DiscriminantCount() > 0 {
		fmt.Fprintf(w, "func (s %s) Which() %s_Which { return %s_Which(%s.Struct(s).Get16(%d)) }\n",
			n.Name, n.Name, n.Name, g_imports.capn(), n.Struct().DiscriminantOffset()*2)
	}

	for _, f := range n.codeOrderFields() {
		switch f.Which() {
		case Field_Which_slot:
			n.defineField(w, f)
		case Field_Which_group:
			g := findNode(f.Group().TypeId())
			fmt.Fprintf(w, "func (s %s) %s() %s { return %s(s) }\n", n.Name, strings.Title(f.Name()), g.Name, g.Name)
			if f.DiscriminantValue() != Field_noDiscriminant {
				fmt.Fprintf(w, "func (s %s) Set%s() { %s.Struct(s).Set16(%d, %d) }\n", n.Name, strings.Title(f.Name()), g_imports.capn(), n.Struct().DiscriminantOffset()*2, f.DiscriminantValue())
			}
			g.defineStructFuncs(w)
		}
	}
}

// This writes the WriteJSON function.
//
// This is an unusual interface, but it was chosen because the types in go-capnproto
// didn't match right to use the json.Marshaler interface.
// This function recurses through the type, writing statements that will dump json to a wire
// For all statements, the json encoder js and the bufio writer b will be in scope.
// The value will be in scope as s. Some features need to redefine s, like unions.
// In that case, Make a new block and redeclare s
func (n *node) defineTypeJsonFuncs(w io.Writer) {
	if C.JSON_enabled {
		ioImp := g_imports.add(importSpec{path: "io", name: "io"})
		bufioImp := g_imports.add(importSpec{path: "bufio", name: "bufio"})
		bytesImp := g_imports.add(importSpec{path: "bytes", name: "bytes"})

		fmt.Fprintf(w, "func (s %s) WriteJSON(w %s.Writer) error {\n", n.Name, ioImp)
		fmt.Fprintf(w, "b := %s.NewWriter(w);", bufioImp)
		fmt.Fprintf(w, "var err error;")
		fmt.Fprintf(w, "var buf []byte;")
		fmt.Fprintf(w, "_ = buf;")

		switch n.Which() {
		case Node_Which_enum:
			n.jsonEnum(w)
		case Node_Which_struct:
			n.jsonStruct(w)
		}

		fmt.Fprintf(w, "err = b.Flush(); return err\n};\n")

		fmt.Fprintf(w, "func (s %s) MarshalJSON() ([]byte, error) {\n", n.Name)
		fmt.Fprintf(w, "var b %s.Buffer; err := s.WriteJSON(&b); return b.Bytes(), err };", bytesImp)

	} else {
		fmt.Fprintf(w, "// capn.JSON_enabled == false so we stub MarshalJSON().")
		fmt.Fprintf(w, "\nfunc (s %s) MarshalJSON() (bs []byte, err error) { return }\n", n.Name)
	}
}

func writeErrCheck(w io.Writer) {
	fmt.Fprintf(w, "if err != nil { return err; };")
}

func (n *node) jsonEnum(w io.Writer) {
	json := g_imports.add(importSpec{path: "encoding/json", name: "json"})
	fmt.Fprintf(w, `buf, err = %s.Marshal(s.String());`, json)
	writeErrCheck(w)
	fmt.Fprintf(w, "_, err = b.Write(buf);")
	writeErrCheck(w)
}

// Write statements that will write a json struct
func (n *node) jsonStruct(w io.Writer) {
	fmt.Fprintf(w, `err = b.WriteByte('{');`)
	writeErrCheck(w)
	for i, f := range n.codeOrderFields() {
		if f.DiscriminantValue() != Field_noDiscriminant {
			enumname := n.Name + "_Which_" + f.Name()
			fmt.Fprintf(w, "if s.Which() == %s {", enumname)
		}
		if i != 0 {
			fmt.Fprintf(w, `
				err = b.WriteByte(',');
			`)
			writeErrCheck(w)
		}
		fmt.Fprintf(w, `_, err = b.WriteString("\"%s\":");`, f.Name())
		writeErrCheck(w)
		f.json(w)
		if f.DiscriminantValue() != Field_noDiscriminant {
			fmt.Fprintf(w, "};")
		}
	}
	fmt.Fprintf(w, `err = b.WriteByte('}');`)
	writeErrCheck(w)
}

// This function writes statements that write the fields json representation to the bufio.
func (f *Field) json(w io.Writer) {

	switch f.Which() {
	case Field_Which_slot:
		fs := f.Slot()
		// we don't generate setters for Void fields
		if fs.Type().Which() == Type_Which_void {
			fs.Type().json(w)
			return
		}
		fmt.Fprintf(w, "{ s := s.%s(); ", strings.Title(f.Name()))
		fs.Type().json(w)
		fmt.Fprintf(w, "}; ")
	case Field_Which_group:
		tid := f.Group().TypeId()
		n := findNode(tid)
		fmt.Fprintf(w, "{ s := s.%s();", strings.Title(f.Name()))

		n.jsonStruct(w)
		fmt.Fprintf(w, "};")
	}
}

func (t Type) json(w io.Writer) {
	switch t.Which() {
	case Type_Which_uint8, Type_Which_uint16, Type_Which_uint32, Type_Which_uint64,
		Type_Which_int8, Type_Which_int16, Type_Which_int32, Type_Which_int64,
		Type_Which_float32, Type_Which_float64, Type_Which_bool, Type_Which_text, Type_Which_data:
		json := g_imports.add(importSpec{path: "encoding/json", name: "json"})
		fmt.Fprintf(w, "buf, err = %s.Marshal(s);", json)
		writeErrCheck(w)
		fmt.Fprintf(w, "_, err = b.Write(buf);")
		writeErrCheck(w)
	case Type_Which_enum, Type_Which_struct:
		// since we handle groups at the field level, only named struct types make it in here
		// so we can just call the named structs json dumper
		fmt.Fprintf(w, "err = s.WriteJSON(b);")
		writeErrCheck(w)
	case Type_Which_list:
		typ := t.List().ElementType()
		which := typ.Which()
		if which == Type_Which_list || which == Type_Which_anyPointer {
			// untyped list, cant do anything but report
			// that a field existed.
			//
			// s will be unused in this case, so ignore
			fmt.Fprintf(w, `_ = s;`)
			fmt.Fprintf(w, `_, err = b.WriteString("\"untyped list\"");`)
			writeErrCheck(w)
			return
		}
		fmt.Fprintf(w, "{ err = b.WriteByte('[');")
		writeErrCheck(w)
		fmt.Fprintf(w, "for i := 0; i < s.Len(); i++ { s := s.At(i); ")
		fmt.Fprintf(w, `if i != 0 { _, err = b.WriteString(", "); };`)
		writeErrCheck(w)
		typ.json(w)
		fmt.Fprintf(w, "}; err = b.WriteByte(']'); };")
		writeErrCheck(w)
	case Type_Which_void:
		fmt.Fprintf(w, `_ = s;`)
		fmt.Fprintf(w, `_, err = b.WriteString("null");`)
		writeErrCheck(w)
	}
}

func (n *node) ObjectSize() string {
	assert(n.Which() == Node_Which_struct, "ObjectSize for invalid struct node")
	return fmt.Sprintf("%s.ObjectSize{DataSize: %d, PointerCount: %d}", g_imports.capn(), int(n.Struct().DataWordCount())*8, n.Struct().PointerCount())
}

func (n *node) defineNewStructFunc(w io.Writer) {
	assert(n.Which() == Node_Which_struct, "invalid struct node")

	os := n.ObjectSize()
	c := g_imports.capn()
	fmt.Fprintf(w, "func New%[1]s(s *%[2]s.Segment) %[1]s { return %[1]s(s.NewStruct(%[3]s)) }\n", n.Name, c, os)
	fmt.Fprintf(w, "func NewRoot%[1]s(s *%[2]s.Segment) %[1]s { return %[1]s(s.NewRootStruct(%[3]s)) }\n", n.Name, c, os)
	fmt.Fprintf(w, "func AutoNew%[1]s(s *%[2]s.Segment) %[1]s { return %[1]s(s.NewStructAR(%[3]s)) }\n", n.Name, c, os)
	fmt.Fprintf(w, "func ReadRoot%[1]s(s *%[2]s.Segment) %[1]s { return %[1]s(s.Root(0).ToStruct()) }\n", n.Name, c)
}

func (n *node) defineStructList(w io.Writer) {
	assert(n.Which() == Node_Which_struct, "invalid struct node")

	c := g_imports.capn()
	fmt.Fprintf(w, "type %s_List %s.PointerList\n", n.Name, c)

	fmt.Fprintf(w, "func New%s_List(s *%s.Segment, sz int) %s_List ", n.Name, c, n.Name)
	switch n.Struct().PreferredListEncoding() {
	case ElementSize_empty:
		fmt.Fprintf(w, "{ return %s_List(s.NewVoidList(sz)) }\n", n.Name)
	case ElementSize_bit:
		fmt.Fprintf(w, "{ return %s_List(s.NewBitList(sz)) }\n", n.Name)
	case ElementSize_byte:
		fmt.Fprintf(w, "{ return %s_List(s.NewUInt8List(sz)) }\n", n.Name)
	case ElementSize_twoBytes:
		fmt.Fprintf(w, "{ return %s_List(s.NewUInt16List(sz)) }\n", n.Name)
	case ElementSize_fourBytes:
		fmt.Fprintf(w, "{ return %s_List(s.NewUInt32List(sz)) }\n", n.Name)
	case ElementSize_eightBytes:
		fmt.Fprintf(w, "{ return %s_List(s.NewUInt64List(sz)) }\n", n.Name)
	default:
		fmt.Fprintf(w, "{ return %s_List(s.NewCompositeList(%s, sz)) }\n",
			n.Name, n.ObjectSize())
	}

	fmt.Fprintf(w, "func (s %s_List) Len() int { return %s.PointerList(s).Len() }\n", n.Name, c)
	fmt.Fprintf(w, "func (s %s_List) At(i int) %s { return %s(%s.PointerList(s).At(i).ToStruct()) }\n", n.Name, n.Name, n.Name, c)
	fmt.Fprintf(w, "func (s %s_List) Set(i int, item %s) { %s.PointerList(s).Set(i, %s.Object(item)) }\n", n.Name, n.Name, c, c)
}

func (n *node) defineStructPromise(w io.Writer) {
	templates.ExecuteTemplate(w, "promise", promiseTemplateParams{
		Node:   n,
		Fields: n.codeOrderFields(),
	})

	for _, f := range n.codeOrderFields() {
		switch f.Which() {
		case Field_Which_slot:
			if t := f.Slot().Type().Which(); t == Type_Which_struct || t == Type_Which_interface || t == Type_Which_anyPointer {
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

func (n *node) definePromiseField(w io.Writer, f Field) {
	slot := f.Slot()
	switch t := slot.Type(); t.Which() {
	case Type_Which_struct:
		ni := findNode(t.Struct().TypeId())
		params := promiseFieldStructTemplateParams{
			Node:   n,
			Field:  f,
			Struct: ni,
		}
		if def := slot.DefaultValue(); def.Which() == Value_Which_struct && def.Struct().HasData() {
			params.BufName = g_bufname
			params.DefaultOffset = copyData(def.Struct())
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
	Interface *node
	ID        int
	Params    *node
	Results   *node
}

func (n *node) methodSet(methods []interfaceMethod) []interfaceMethod {
	for i, ms := 0, n.Interface().Methods(); i < ms.Len(); i++ {
		m := ms.At(i)
		methods = append(methods, interfaceMethod{
			Method:    m,
			Interface: n,
			ID:        i,
			Params:    findNode(m.ParamStructType()),
			Results:   findNode(m.ResultStructType()),
		})
	}
	// TODO(light): sort added methods by code order

	for i, supers := 0, n.Interface().Superclasses(); i < supers.Len(); i++ {
		s := supers.At(i)
		methods = findNode(s.Id()).methodSet(methods)
	}
	return methods
}

func (n *node) defineInterfaceClient(w io.Writer) {
	m := n.methodSet(nil)
	templates.ExecuteTemplate(w, "interfaceClient", interfaceClientTemplateParams{
		Node:        n,
		Annotations: parseAnnotations(n.Annotations()),
		Methods:     m,
	})
}

func (n *node) defineInterfaceServer(w io.Writer) {
	m := n.methodSet(nil)
	templates.ExecuteTemplate(w, "interfaceServer", interfaceServerTemplateParams{
		Node:        n,
		Annotations: parseAnnotations(n.Annotations()),
		Methods:     m,
	})
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
	g_segment = C.NewBuffer([]byte{})
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
			n.defineTypeJsonFuncs(&buf)
		case Node_Which_struct:
			if !n.Struct().IsGroup() {
				n.defineStructTypes(&buf, nil)
				n.defineStructEnums(&buf)
				n.defineNewStructFunc(&buf)
				n.defineStructFuncs(&buf)
				n.defineTypeJsonFuncs(&buf)
				n.defineStructList(&buf)
				n.defineStructPromise(&buf)
			}
		case Node_Which_interface:
			n.defineInterfaceClient(&buf)
			n.defineInterfaceServer(&buf)
		}
	}

	if f.pkg == "" {
		return fmt.Errorf("missing package annotation for %s", reqf.Filename())
	}

	if dirPath, _ := filepath.Split(reqf.Filename()); dirPath != "" {
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
	if len(g_segment.Data) > 0 {
		fmt.Fprintf(&unformatted, "var %s = %s.NewBuffer([]byte{", g_bufname, g_imports.capn())
		for i, b := range g_segment.Data {
			if i%8 == 0 {
				fmt.Fprintf(&unformatted, "\n")
			}
			fmt.Fprintf(&unformatted, "%d,", b)
		}
		fmt.Fprintf(&unformatted, "\n})\n")
	}
	formatted, err := format.Source(unformatted.Bytes())
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't format generated code:", err)
		formatted = unformatted.Bytes()
	}

	file, err := os.Create(reqf.Filename() + ".go")
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
	s, err := C.ReadFromStream(os.Stdin, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, "capnpc-go: Reading input:", err)
		os.Exit(1)
	}

	req := ReadRootCodeGeneratorRequest(s)
	allfiles := []*node{}

	for i := 0; i < req.Nodes().Len(); i++ {
		ni := req.Nodes().At(i)
		n := &node{Node: ni}
		g_nodes[n.Id()] = n

		if n.Which() == Node_Which_file {
			allfiles = append(allfiles, n)
		}
	}

	for _, f := range allfiles {
		ann := parseAnnotations(f.Annotations())
		f.pkg = ann.Package
		f.imp = ann.Import
		for i := 0; i < f.NestedNodes().Len(); i++ {
			nn := f.NestedNodes().At(i)
			if ni := g_nodes[nn.Id()]; ni != nil {
				ni.resolveName("", nn.Name(), f)
			}
		}
	}

	success := true
	for i := 0; i < req.RequestedFiles().Len(); i++ {
		reqf := req.RequestedFiles().At(i)
		err := generateFile(reqf)
		if err != nil {
			fmt.Fprintf(os.Stderr, "capnpc-go: Generating Go for %s failed: %v\n", reqf.Filename(), err)
			success = false
		}
	}
	if !success {
		os.Exit(1)
	}
}
