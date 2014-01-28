package main

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"

	C "github.com/glycerine/go-capnproto"
)

var (
	go_capnproto_import = "github.com/glycerine/go-capnproto"
	fprintf             = fmt.Fprintf
	sprintf             = fmt.Sprintf
	title               = strings.Title
)

var g_nodes = make(map[uint64]*node)
var g_imported map[string]bool
var g_segment *C.Segment
var g_bufname string

type node struct {
	Node
	pkg   string
	imp   string
	nodes []*node
	name  string
}

func assert(chk bool, format string, a ...interface{}) {
	if !chk {
		panic(sprintf(format, a...))
		os.Exit(1)
	}
}

func copyData(obj C.Object) int {
	r, off, err := g_segment.NewRoot()
	assert(err == nil, "%v\n", err)
	err = r.Set(0, obj)
	assert(err == nil, "%v\n", err)
	return off
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
		g_imported[n.imp] = true
		return n.pkg + "."
	}
}

func (n *node) remoteName(from *node) string {
	return n.remoteScope(from) + n.name
}

func (n *node) resolveName(base, name string, file *node) {
	if base != "" {
		n.name = base + title(name)
	} else {
		n.name = title(name)
	}

	n.pkg = file.pkg
	n.imp = file.imp

	if n.Which() != NODE_STRUCT || !n.Struct().IsGroup() {
		file.nodes = append(file.nodes, n)
	}

	for _, nn := range n.NestedNodes().ToArray() {
		if ni := g_nodes[nn.Id()]; ni != nil {
			ni.resolveName(n.name, nn.Name(), file)
		}
	}

	if n.Which() == NODE_STRUCT {
		for _, f := range n.Struct().Fields().ToArray() {
			if f.Which() == FIELD_GROUP {
				findNode(f.Group().TypeId()).resolveName(n.name, f.Name(), file)
			}
		}
	}
}

type enumval struct {
	Enumerant
	val int
	tag string
}

func (n *node) defineEnum(w io.Writer) {
	for _, a := range n.Annotations().ToArray() {
		if a.Id() == C.Doc {
			fprintf(w, "// %s\n", a.Value().Text())
		}
	}
	fprintf(w, "type %s uint16\n", n.name)

	if es := n.Enum().Enumerants(); es.Len() > 0 {
		fprintf(w, "const (\n")

		ev := make([]enumval, es.Len())
		for i := 0; i < es.Len(); i++ {
			e := es.At(i)

			t := e.Name()
			for _, an := range e.Annotations().ToArray() {
				if an.Id() == C.Tag {
					t = an.Value().Text()
				} else if an.Id() == C.Notag {
					t = ""
				}
			}
			ev[e.CodeOrder()] = enumval{e, i, t}
		}

		fprintf(w, "%s_%s %s = %d\n", strings.ToUpper(n.name), strings.ToUpper(ev[0].Name()), n.name, ev[0].val)

		for _, e := range ev[1:] {
			fprintf(w, "%s_%s = %d\n", strings.ToUpper(n.name), strings.ToUpper(e.Name()), e.val)
		}

		fprintf(w, ")\n")

		fprintf(w, "func (c %s) String() string {\n", n.name)
		fprintf(w, "switch c {\n")
		for _, e := range ev {
			if e.tag != "" {
				fprintf(w, "case %s_%s: return \"%s\"\n", strings.ToUpper(n.name), strings.ToUpper(e.Name()), e.tag)
			}
		}
		fprintf(w, "default: return \"\"\n")
		fprintf(w, "}\n}\n")
	}

	fprintf(w, "type %s_List C.PointerList\n", n.name)
	fprintf(w, "func New%sList(s *C.Segment, sz int) %s_List { return %s_List(s.NewUInt16List(sz)) }\n", n.name, n.name, n.name)
	fprintf(w, "func (s %s_List) Len() int { return C.UInt16List(s).Len() }\n", n.name)
	fprintf(w, "func (s %s_List) At(i int) %s { return %s(C.UInt16List(s).At(i)) }\n", n.name, n.name, n.name)
	fprintf(w, "func (s %s_List) ToArray() []%s { return *(*[]%s)(unsafe.Pointer(C.UInt16List(s).ToEnumArray())) }\n", n.name, n.name, n.name)

	g_imported["unsafe"] = true
}

func (n *node) writeValue(w io.Writer, t Type, v Value) {
	switch t.Which() {
	case TYPE_VOID, TYPE_INTERFACE:
		fprintf(w, "C.Void{}")

	case TYPE_BOOL:
		assert(v.Which() == VALUE_BOOL, "expected bool value")
		if v.Bool() {
			fprintf(w, "true")
		} else {
			fprintf(w, "false")
		}

	case TYPE_INT8:
		assert(v.Which() == VALUE_INT8, "expected int8 value")
		fprintf(w, "int8(%d)", v.Int8())

	case TYPE_UINT8:
		assert(v.Which() == VALUE_UINT8, "expected uint8 value")
		fprintf(w, "uint8(%d)", v.Uint8())

	case TYPE_INT16:
		assert(v.Which() == VALUE_INT16, "expected int16 value")
		fprintf(w, "int16(%d)", v.Int16())

	case TYPE_UINT16:
		assert(v.Which() == VALUE_UINT16, "expected uint16 value")
		fprintf(w, "uint16(%d)", v.Uint16())

	case TYPE_INT32:
		assert(v.Which() == VALUE_INT32, "expected int32 value")
		fprintf(w, "int32(%d)", v.Int32())

	case TYPE_UINT32:
		assert(v.Which() == VALUE_UINT32, "expected uint32 value")
		fprintf(w, "uint32(%d)", v.Uint32())

	case TYPE_INT64:
		assert(v.Which() == VALUE_INT64, "expected int64 value")
		fprintf(w, "int64(%d)", v.Int64())

	case TYPE_UINT64:
		assert(v.Which() == VALUE_UINT64, "expected uint64 value")
		fprintf(w, "uint64(%d)", v.Uint64())

	case TYPE_FLOAT32:
		assert(v.Which() == VALUE_FLOAT32, "expected float32 value")
		fprintf(w, "math.Float32frombits(0x%x)", math.Float32bits(v.Float32()))
		g_imported["math"] = true

	case TYPE_FLOAT64:
		assert(v.Which() == VALUE_FLOAT64, "expected float64 value")
		fprintf(w, "math.Float64frombits(0x%x)", math.Float64bits(v.Float64()))
		g_imported["math"] = true

	case TYPE_TEXT:
		assert(v.Which() == VALUE_TEXT, "expected text value")
		fprintf(w, "%s", strconv.Quote(v.Text()))

	case TYPE_DATA:
		assert(v.Which() == VALUE_DATA, "expected data value")
		fprintf(w, "[]byte{")
		for i, b := range v.Data() {
			if i > 0 {
				fprintf(w, ", ")
			}
			fprintf(w, "%d", b)
		}
		fprintf(w, "}")

	case TYPE_ENUM:
		assert(v.Which() == VALUE_ENUM, "expected enum value")
		en := findNode(t.Enum().TypeId())
		assert(en.Which() == NODE_ENUM, "expected enum type ID")
		ev := en.Enum().Enumerants()
		if val := int(v.Enum()); val >= ev.Len() {
			fprintf(w, "%s(%d)", en.remoteName(n), val)
		} else {
			fprintf(w, "%s%s", en.remoteScope(n), ev.At(val).Name())
		}

	case TYPE_STRUCT:
		fprintf(w, "%s(%s.Root(%d))", findNode(t.Struct().TypeId()).remoteName(n), g_bufname, copyData(v.Struct()))

	case TYPE_OBJECT:
		fprintf(w, "%s.Root(%d)", g_bufname, copyData(v.Object()))

	case TYPE_LIST:
		assert(v.Which() == VALUE_LIST, "expected list value")

		switch lt := t.List().ElementType(); lt.Which() {
		case TYPE_VOID, TYPE_INTERFACE:
			fprintf(w, "make([]C.Void, %d)", v.List().ToVoidList().Len())
		case TYPE_BOOL:
			fprintf(w, "C.BitList(%s.Root(%d))", g_bufname, copyData(v.List()))
		case TYPE_INT8:
			fprintf(w, "C.Int8List(%s.Root(%d))", g_bufname, copyData(v.List()))
		case TYPE_UINT8:
			fprintf(w, "C.UInt8List(%s.Root(%d))", g_bufname, copyData(v.List()))
		case TYPE_INT16:
			fprintf(w, "C.Int16List(%s.Root(%d))", g_bufname, copyData(v.List()))
		case TYPE_UINT16:
			fprintf(w, "C.UInt16List(%s.Root(%d))", g_bufname, copyData(v.List()))
		case TYPE_INT32:
			fprintf(w, "C.Int32List(%s.Root(%d))", g_bufname, copyData(v.List()))
		case TYPE_UINT32:
			fprintf(w, "C.UInt32List(%s.Root(%d))", g_bufname, copyData(v.List()))
		case TYPE_FLOAT32:
			fprintf(w, "C.Float32List(%s.Root(%d))", g_bufname, copyData(v.List()))
		case TYPE_INT64:
			fprintf(w, "C.Int64List(%s.Root(%d))", g_bufname, copyData(v.List()))
		case TYPE_UINT64:
			fprintf(w, "C.UInt64List(%s.Root(%d))", g_bufname, copyData(v.List()))
		case TYPE_FLOAT64:
			fprintf(w, "C.Float64List(%s.Root(%d))", g_bufname, copyData(v.List()))
		case TYPE_TEXT:
			fprintf(w, "C.TextList(%s.Root(%d))", g_bufname, copyData(v.List()))
		case TYPE_DATA:
			fprintf(w, "C.DataList(%s.Root(%d))", g_bufname, copyData(v.List()))
		case TYPE_ENUM:
			fprintf(w, "%s_List(%s.Root(%d))", findNode(lt.Enum().TypeId()).remoteName(n), g_bufname, copyData(v.List()))
		case TYPE_STRUCT:
			fprintf(w, "%s_List(%s.Root(%d))", findNode(lt.Struct().TypeId()).remoteName(n), g_bufname, copyData(v.List()))
		case TYPE_LIST, TYPE_OBJECT:
			fprintf(w, "C.PointerList(%s.Root(%d))", g_bufname, copyData(v.List()))
		}
	}
}

func (n *node) defineAnnotation(w io.Writer) {
	fprintf(w, "const %s = uint64(0x%x)\n", n.name, n.Id())
}

func constIsVar(n *node) bool {
	switch n.Const().Type().Which() {
	case TYPE_BOOL, TYPE_INT8, TYPE_UINT8, TYPE_INT16,
		TYPE_UINT16, TYPE_INT32, TYPE_UINT32, TYPE_INT64,
		TYPE_UINT64, TYPE_TEXT, TYPE_ENUM:
		return false
	default:
		return true
	}
}

func defineConstNodes(w io.Writer, nodes []*node) {

	any := false

	for _, n := range nodes {
		if n.Which() == NODE_CONST && !constIsVar(n) {
			if !any {
				fprintf(w, "const (\n")
				any = true
			}
			fprintf(w, "%s = ", n.name)
			n.writeValue(w, n.Const().Type(), n.Const().Value())
			fprintf(w, "\n")
		}
	}

	if any {
		fprintf(w, ")\n")
	}

	any = false

	for _, n := range nodes {
		if n.Which() == NODE_CONST && constIsVar(n) {
			if !any {
				fprintf(w, "var (\n")
				any = true
			}
			fprintf(w, "%s = ", n.name)
			n.writeValue(w, n.Const().Type(), n.Const().Value())
			fprintf(w, "\n")
		}
	}

	if any {
		fprintf(w, ")\n")
	}
}

func (n *node) defineField(w io.Writer, f Field) {
	t := f.Slot().Type()
	def := f.Slot().DefaultValue()
	off := f.Slot().Offset()

	if t.Which() == TYPE_VOID || t.Which() == TYPE_INTERFACE {
		return
	}

	settag := ""
	if f.DiscriminantValue() != 0xFFFF {
		settag = sprintf(" C.Struct(s).Set16(%d, %d);", n.Struct().DiscriminantOffset()*2, f.DiscriminantValue())
	}

	var g, s bytes.Buffer

	for _, a := range f.Annotations().ToArray() {
		if a.Id() == C.Doc {
			fprintf(&g, "// %s\n", a.Value().Text())
		}
	}
	fprintf(&g, "func (s %s) %s() ", n.name, title(f.Name()))
	fprintf(&s, "func (s %s) Set%s", n.name, title(f.Name()))

	switch t.Which() {
	case TYPE_BOOL:
		assert(def.Which() == VALUE_VOID || def.Which() == VALUE_BOOL, "expected bool default")
		if def.Which() == VALUE_BOOL && def.Bool() {
			fprintf(&g, "bool { return !C.Struct(s).Get1(%d) }\n", off)
			fprintf(&s, "(v bool) {%s C.Struct(s).Set1(%d, !v) }\n", settag, off)
		} else {
			fprintf(&g, "bool { return C.Struct(s).Get1(%d) }\n", off)
			fprintf(&s, "(v bool) {%s C.Struct(s).Set1(%d, v) }\n", settag, off)
		}

	case TYPE_INT8:
		assert(def.Which() == VALUE_VOID || def.Which() == VALUE_INT8, "expected int8 default")
		if def.Which() == VALUE_INT8 && def.Int8() != 0 {
			fprintf(&g, "int8 { return int8(C.Struct(s).Get8(%d)) ^ %d }\n", off, def.Int8())
			fprintf(&s, "(v int8) {%s C.Struct(s).Set8(%d, uint8(v^%d)) }\n", settag, off, def.Int8())
		} else {
			fprintf(&g, "int8 { return int8(C.Struct(s).Get8(%d)) }\n", off)
			fprintf(&s, "(v int8) {%s C.Struct(s).Set8(%d, uint8(v)) }\n", settag, off)
		}

	case TYPE_UINT8:
		assert(def.Which() == VALUE_VOID || def.Which() == VALUE_UINT8, "expected uint8 default")
		if def.Which() == VALUE_UINT8 && def.Uint8() != 0 {
			fprintf(&g, "uint8 { return C.Struct(s).Get8(%d) ^ %d }\n", off, def.Uint8())
			fprintf(&s, "(v uint8) {%s C.Struct(s).Set8(%d, v^%d) }\n", settag, off, def.Uint8())
		} else {
			fprintf(&g, "uint8 { return C.Struct(s).Get8(%d) }\n", off)
			fprintf(&s, "(v uint8) {%s C.Struct(s).Set8(%d, v) }\n", settag, off)
		}

	case TYPE_INT16:
		assert(def.Which() == VALUE_VOID || def.Which() == VALUE_INT16, "expected int16 default")
		if def.Which() == VALUE_INT16 && def.Int16() != 0 {
			fprintf(&g, "int16 { return int16(C.Struct(s).Get16(%d)) ^ %d }\n", off*2, def.Int16())
			fprintf(&s, "(v int16) {%s C.Struct(s).Set16(%d, uint16(v^%d)) }\n", settag, off*2, def.Int16())
		} else {
			fprintf(&g, "int16 { return int16(C.Struct(s).Get16(%d)) }\n", off*2)
			fprintf(&s, "(v int16) {%s C.Struct(s).Set16(%d, uint16(v)) }\n", settag, off*2)
		}

	case TYPE_UINT16:
		assert(def.Which() == VALUE_VOID || def.Which() == VALUE_UINT16, "expected uint16 default")
		if def.Which() == VALUE_UINT16 && def.Uint16() != 0 {
			fprintf(&g, "uint16 { return C.Struct(s).Get16(%d) ^ %d }\n", off*2, def.Uint16())
			fprintf(&s, "(v uint16) {%s C.Struct(s).Set16(%d, v^%d) }\n", settag, off*2, def.Uint16())
		} else {
			fprintf(&g, "uint16 { return C.Struct(s).Get16(%d) }\n", off*2)
			fprintf(&s, "(v uint16) {%s C.Struct(s).Set16(%d, v) }\n", settag, off*2)
		}

	case TYPE_INT32:
		assert(def.Which() == VALUE_VOID || def.Which() == VALUE_INT32, "expected int32 default")
		if def.Which() == VALUE_INT32 && def.Int32() != 0 {
			fprintf(&g, "int32 { return int32(C.Struct(s).Get32(%d)) ^ %d }\n", off*4, def.Int32())
			fprintf(&s, "(v int32) {%s C.Struct(s).Set32(%d, uint32(v^%d)) }\n", settag, off*4, def.Int32())
		} else {
			fprintf(&g, "int32 { return int32(C.Struct(s).Get32(%d)) }\n", off*4)
			fprintf(&s, "(v int32) {%s C.Struct(s).Set32(%d, uint32(v)) }\n", settag, off*4)
		}

	case TYPE_UINT32:
		assert(def.Which() == VALUE_VOID || def.Which() == VALUE_UINT32, "expected uint32 default")
		if def.Which() == VALUE_UINT32 && def.Uint32() != 0 {
			fprintf(&g, "uint32 { return C.Struct(s).Get32(%d) ^ %d }\n", off*4, def.Uint32())
			fprintf(&s, "(v uint32) {%s C.Struct(s).Set32(%d, v^%d) }\n", settag, off*4, def.Uint32())
		} else {
			fprintf(&g, "uint32 { return C.Struct(s).Get32(%d) }\n", off*4)
			fprintf(&s, "(v uint32) {%s C.Struct(s).Set32(%d, v) }\n", settag, off*4)
		}

	case TYPE_INT64:
		assert(def.Which() == VALUE_VOID || def.Which() == VALUE_INT64, "expected int64 default")
		if def.Which() == VALUE_INT64 && def.Int64() != 0 {
			fprintf(&g, "int64 { return int64(C.Struct(s).Get64(%d)) ^ %d }\n", off*8, def.Int64())
			fprintf(&s, "(v int64) {%s C.Struct(s).Set64(%d, uint64(v^%d)) }\n", settag, off*8, def.Int64())
		} else {
			fprintf(&g, "int64 { return int64(C.Struct(s).Get64(%d)) }\n", off*8)
			fprintf(&s, "(v int64) {%s C.Struct(s).Set64(%d, uint64(v)) }\n", settag, off*8)
		}

	case TYPE_UINT64:
		assert(def.Which() == VALUE_VOID || def.Which() == VALUE_UINT64, "expected uint64 default")
		if def.Which() == VALUE_UINT64 && def.Uint64() != 0 {
			fprintf(&g, "uint64 { return C.Struct(s).Get64(%d) ^ %d }\n", off*8, def.Uint64())
			fprintf(&s, "(v uint64) {%s C.Struct(s).Set64(%d, v^%d) }\n", settag, off*8, def.Uint64())
		} else {
			fprintf(&g, "uint64 { return C.Struct(s).Get64(%d) }\n", off*8)
			fprintf(&s, "(v uint64) {%s C.Struct(s).Set64(%d, v) }\n", settag, off*8)
		}

	case TYPE_FLOAT32:
		assert(def.Which() == VALUE_VOID || def.Which() == VALUE_FLOAT32, "expected float32 default")
		if def.Which() == VALUE_FLOAT32 && def.Float32() != 0 {
			fprintf(&g, "float32 { return math.Float32frombits(C.Struct(s).Get32(%d) ^ 0x%x) }\n", off*4, math.Float32bits(def.Float32()))
			fprintf(&s, "(v float32) {%s C.Struct(s).Set32(%d, math.Float32bits(v) ^ 0x%x) }\n", settag, off*4, math.Float32bits(def.Float32()))
		} else {
			fprintf(&g, "float32 { return math.Float32frombits(C.Struct(s).Get32(%d)) }\n", off*4)
			fprintf(&s, "(v float32) {%s C.Struct(s).Set32(%d, math.Float32bits(v)) }\n", settag, off*4)
		}
		g_imported["math"] = true

	case TYPE_FLOAT64:
		assert(def.Which() == VALUE_VOID || def.Which() == VALUE_FLOAT64, "expected float64 default")
		if def.Which() == VALUE_FLOAT64 && def.Float64() != 0 {
			fprintf(&g, "float64 { return math.Float64frombits(C.Struct(s).Get64(%d) ^ 0x%x) }\n", off*8, math.Float64bits(def.Float64()))
			fprintf(&s, "(v float64) {%s C.Struct(s).Set64(%d, math.Float64bits(v) ^ 0x%x) }\n", settag, off*8, math.Float64bits(def.Float64()))
		} else {
			fprintf(&g, "float64 { return math.Float64frombits(C.Struct(s).Get64(%d)) }\n", off*8)
			fprintf(&s, "(v float64) {%s C.Struct(s).Set64(%d, math.Float64bits(v)) }\n", settag, off*8)
		}
		g_imported["math"] = true

	case TYPE_TEXT:
		assert(def.Which() == VALUE_VOID || def.Which() == VALUE_TEXT, "expected text default")
		if def.Which() == VALUE_TEXT && def.Text() != "" {
			fprintf(&g, "string { return C.Struct(s).GetObject(%d).ToTextDefault(%s) }\n", off, strconv.Quote(def.Text()))
		} else {
			fprintf(&g, "string { return C.Struct(s).GetObject(%d).ToText() }\n", off)
		}
		fprintf(&s, "(v string) {%s C.Struct(s).SetObject(%d, s.Segment.NewText(v)) }\n", settag, off)

	case TYPE_DATA:
		assert(def.Which() == VALUE_VOID || def.Which() == VALUE_DATA, "expected data default")
		if def.Which() == VALUE_DATA && len(def.Data()) > 0 {
			dstr := "[]byte{"
			for i, b := range def.Data() {
				if i > 0 {
					dstr += ", "
				}
				dstr += sprintf("%d", b)
			}
			dstr += "}"
			fprintf(&g, "[]byte { return C.Struct(s).GetObject(%d).ToDataDefault(%s) }\n", off, dstr)
		} else {
			fprintf(&g, "[]byte { return C.Struct(s).GetObject(%d).ToData() }\n", off)
		}
		fprintf(&s, "(v []byte) {%s C.Struct(s).SetObject(%d, s.Segment.NewData(v)) }\n", settag, off)

	case TYPE_ENUM:
		ni := findNode(t.Enum().TypeId())
		assert(def.Which() == VALUE_VOID || def.Which() == VALUE_ENUM, "expected enum default")
		if def.Which() == VALUE_ENUM && def.Enum() != 0 {
			fprintf(&g, "%s { return %s(C.Struct(s).Get16(%d) ^ %d) }\n", ni.remoteName(n), ni.remoteName(n), off*2, def.Enum())
			fprintf(&s, "(v %s) {%s C.Struct(s).Set16(%d, uint16(v)^%d) }\n", ni.remoteName(n), settag, off*2, def.Uint16())
		} else {
			fprintf(&g, "%s { return %s(C.Struct(s).Get16(%d)) }\n", ni.remoteName(n), ni.remoteName(n), off*2)
			fprintf(&s, "(v %s) {%s C.Struct(s).Set16(%d, uint16(v)) }\n", ni.remoteName(n), settag, off*2)
		}

	case TYPE_STRUCT:
		ni := findNode(t.Struct().TypeId())
		assert(def.Which() == VALUE_VOID || def.Which() == VALUE_STRUCT, "expected struct default")
		if def.Which() == VALUE_STRUCT && def.Struct().HasData() {
			fprintf(&g, "%s { return %s(C.Struct(s).GetObject(%d).ToStructDefault(%s, %d)) }\n",
				ni.remoteName(n), ni.remoteName(n), off, g_bufname, copyData(def.Struct()))
		} else {
			fprintf(&g, "%s { return %s(C.Struct(s).GetObject(%d).ToStruct()) }\n",
				ni.remoteName(n), ni.remoteName(n), off)
		}
		fprintf(&s, "(v %s) {%s C.Struct(s).SetObject(%d, C.Object(v)) }\n", ni.remoteName(n), settag, off)

	case TYPE_OBJECT:
		assert(def.Which() == VALUE_VOID || def.Which() == VALUE_OBJECT, "expected object default")
		if def.Which() == VALUE_OBJECT && def.Object().HasData() {
			fprintf(&g, "C.Object { return C.Struct(s).GetObject(%d).ToObjectDefault(%s, %d) }\n",
				off, g_bufname, copyData(def.Object()))
		} else {
			fprintf(&g, "C.Object { return C.Struct(s).GetObject(%d) }\n", off)
		}
		fprintf(&s, "(v C.Object) {%s C.Struct(s).SetObject(%d, v) }\n", settag, off)

	case TYPE_LIST:
		assert(def.Which() == VALUE_VOID || def.Which() == VALUE_LIST, "expected list default")

		typ := ""

		switch lt := t.List().ElementType(); lt.Which() {
		case TYPE_VOID, TYPE_INTERFACE:
			typ = "C.VoidList"
		case TYPE_BOOL:
			typ = "C.BitList"
		case TYPE_INT8:
			typ = "C.Int8List"
		case TYPE_UINT8:
			typ = "C.UInt8List"
		case TYPE_INT16:
			typ = "C.Int16List"
		case TYPE_UINT16:
			typ = "C.UInt16List"
		case TYPE_INT32:
			typ = "C.Int32List"
		case TYPE_UINT32:
			typ = "C.UInt32List"
		case TYPE_INT64:
			typ = "C.Int64List"
		case TYPE_UINT64:
			typ = "C.UInt64List"
		case TYPE_FLOAT32:
			typ = "C.Float32List"
		case TYPE_FLOAT64:
			typ = "C.Float64List"
		case TYPE_TEXT:
			typ = "C.TextList"
		case TYPE_DATA:
			typ = "C.DataList"
		case TYPE_ENUM:
			ni := findNode(lt.Enum().TypeId())
			typ = sprintf("%s_List", ni.remoteName(n))
		case TYPE_STRUCT:
			ni := findNode(lt.Struct().TypeId())
			typ = sprintf("%s_List", ni.remoteName(n))
		case TYPE_OBJECT, TYPE_LIST:
			typ = "C.PointerList"
		}

		ldef := C.Object{}
		if def.Which() == VALUE_LIST {
			ldef = def.List()
		}

		if ldef.HasData() {
			fprintf(&g, "%s { return %s(C.Struct(s).GetObject(%d).ToListDefault(%s, %d)) }\n",
				typ, typ, off, g_bufname, copyData(ldef))
		} else {
			fprintf(&g, "%s { return %s(C.Struct(s).GetObject(%d)) }\n",
				typ, typ, off)
		}

		fprintf(&s, "(v %s) {%s C.Struct(s).SetObject(%d, C.Object(v)) }\n", typ, settag, off)
	}

	w.Write(g.Bytes())
	w.Write(s.Bytes())
}

func (n *node) codeOrderFields() []Field {
	fields := n.Struct().Fields().ToArray()
	mbrs := make([]Field, len(fields))
	for _, f := range fields {
		mbrs[f.CodeOrder()] = f
	}
	return mbrs
}

func (n *node) defineStructTypes(w io.Writer, baseNode *node) {
	assert(n.Which() == NODE_STRUCT, "invalid struct node")

	for _, a := range n.Annotations().ToArray() {
		if a.Id() == C.Doc {
			fprintf(w, "// %s\n", a.Value().Text())
		}
	}
	if baseNode != nil {
		fprintf(w, "type %s %s\n", n.name, baseNode.name)
	} else {
		fprintf(w, "type %s C.Struct\n", n.name)
		baseNode = n
	}

	for _, f := range n.codeOrderFields() {
		if f.Which() == FIELD_GROUP {
			findNode(f.Group().TypeId()).defineStructTypes(w, baseNode)
		}
	}
}

func (n *node) defineStructEnums(w io.Writer) {
	assert(n.Which() == NODE_STRUCT, "invalid struct node")

	if n.Struct().DiscriminantCount() > 0 {
		fprintf(w, "type %s_Which uint16\n", n.name)
		fprintf(w, "const (\n")

		writtenType := false
		for _, f := range n.codeOrderFields() {
			if f.DiscriminantValue() == 0xFFFF {
				// Non-union member
			} else if !writtenType {
				fprintf(w, "%s_%s %s_Which = %d\n", strings.ToUpper(n.name), strings.ToUpper(f.Name()), n.name, f.DiscriminantValue())
				writtenType = true
			} else {
				fprintf(w, "%s_%s = %d\n", strings.ToUpper(n.name), strings.ToUpper(f.Name()), f.DiscriminantValue())
			}
		}
		fprintf(w, ")\n")
	}

	for _, f := range n.codeOrderFields() {
		if f.Which() == FIELD_GROUP {
			findNode(f.Group().TypeId()).defineStructEnums(w)
		}
	}
}

func (n *node) defineStructFuncs(w io.Writer) {
	assert(n.Which() == NODE_STRUCT, "invalid struct node")

	if n.Struct().DiscriminantCount() > 0 {
		fprintf(w, "func (s %s) Which() %s_Which { return %s_Which(C.Struct(s).Get16(%d)) }\n",
			n.name, n.name, n.name, n.Struct().DiscriminantOffset()*2)
	}

	for _, f := range n.codeOrderFields() {
		switch f.Which() {
		case FIELD_SLOT:
			n.defineField(w, f)
		case FIELD_GROUP:
			g := findNode(f.Group().TypeId())
			fprintf(w, "func (s %s) %s() %s { return %s(s) }\n", n.name, title(f.Name()), g.name, g.name)
			if f.DiscriminantValue() != 0xFFFF {
				fprintf(w, "func (s %s) Set%s() { C.Struct(s).Set16(%d, %d) }\n", n.name, title(f.Name()), n.Struct().DiscriminantOffset()*2, f.DiscriminantValue())
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
	g_imported["io"] = true
	g_imported["bufio"] = true
	g_imported["bytes"] = true

	fprintf(w, "func (s %s) WriteJSON(w io.Writer) error {\n", n.name)
	fprintf(w, "b := bufio.NewWriter(w);")
	fprintf(w, "var err error;")
	fprintf(w, "var buf []byte;")
	fprintf(w, "_ = buf;")

	switch n.Which() {
	case NODE_ENUM:
		n.jsonEnum(w)
	case NODE_STRUCT:
		n.jsonStruct(w)
	}

	fprintf(w, "err = b.Flush(); return err\n};\n")

	fprintf(w, "func (s %s) MarshalJSON() ([]byte, error) {\n", n.name)
	fprintf(w, "b := bytes.Buffer{}; err := s.WriteJSON(&b); return b.Bytes(), err };")
}

func writeErrCheck(w io.Writer) {
	fprintf(w, "if err != nil { return err; };")
}

func (n *node) jsonEnum(w io.Writer) {
	g_imported["encoding/json"] = true
	fprintf(w, `buf, err = json.Marshal(s.String());`)
	writeErrCheck(w)
	fprintf(w, "_, err = b.Write(buf);")
	writeErrCheck(w)
}

// Write statements that will write a json struct
func (n *node) jsonStruct(w io.Writer) {
	fprintf(w, `err = b.WriteByte('{');`)
	writeErrCheck(w)
	for i, f := range n.codeOrderFields() {
		if f.DiscriminantValue() != 0xFFFF {
			enumname := fmt.Sprintf("%s_%s", strings.ToUpper(n.name), strings.ToUpper(f.Name()))
			fprintf(w, "if s.Which() == %s {", enumname)
		}
		if i != 0 {
			fprintf(w, `
				err = b.WriteByte(',');
			`)
			writeErrCheck(w)
		}
		fprintf(w, `_, err = b.WriteString("\"%s\":");`, f.Name())
		writeErrCheck(w)
		f.json(w)
		if f.DiscriminantValue() != 0xFFFF {
			fprintf(w, "};")
		}
	}
	fprintf(w, `err = b.WriteByte('}');`)
	writeErrCheck(w)
}

// This function writes statements that write the fields json representation to the bufio.
func (f *Field) json(w io.Writer) {
	switch f.Which() {
	case FIELD_SLOT:
		fs := f.Slot()
		fprintf(w, "{ s := s.%s(); ", title(f.Name()))
		fs.Type().json(w)
		fprintf(w, "}; ")
	case FIELD_GROUP:
		tid := f.Group().TypeId()
		n := findNode(tid)
		fprintf(w, "{ s := s.%s();", title(f.Name()))
		n.jsonStruct(w)
		fprintf(w, "};")
	}
}

func (t Type) json(w io.Writer) {
	switch t.Which() {
	case TYPE_UINT8, TYPE_UINT16, TYPE_UINT32, TYPE_UINT64,
		TYPE_INT8, TYPE_INT16, TYPE_INT32, TYPE_INT64,
		TYPE_FLOAT32, TYPE_FLOAT64, TYPE_BOOL, TYPE_TEXT, TYPE_DATA:
		g_imported["encoding/json"] = true
		fprintf(w, "buf, err = json.Marshal(s);")
		writeErrCheck(w)
		fprintf(w, "_, err = b.Write(buf);")
		writeErrCheck(w)
	case TYPE_ENUM, TYPE_STRUCT:
		// since we handle groups at the field level, only named struct types make it in here
		// so we can just call the named structs json dumper
		fprintf(w, "err = s.WriteJSON(b);")
		writeErrCheck(w)
	case TYPE_LIST:
		fprintf(w, "{ err = b.WriteByte('[');")
		writeErrCheck(w)
		fprintf(w, "for i, s := range s.ToArray() {")
		fprintf(w, `if i != 0 { _, err = b.WriteString(", "); };`)
		writeErrCheck(w)
		typ := t.List().ElementType()
		typ.json(w)
		fprintf(w, "}; err = b.WriteByte(']'); };")
		writeErrCheck(w)
	}
}

func (n *node) defineNewStructFunc(w io.Writer) {
	assert(n.Which() == NODE_STRUCT, "invalid struct node")

	fprintf(w, "func New%s(s *C.Segment) %s { return %s(s.NewStruct(%d, %d)) }\n",
		n.name, n.name, n.name, n.Struct().DataWordCount()*8, n.Struct().PointerCount())
	fprintf(w, "func NewRoot%s(s *C.Segment) %s { return %s(s.NewRootStruct(%d, %d)) }\n",
		n.name, n.name, n.name, n.Struct().DataWordCount()*8, n.Struct().PointerCount())
	fprintf(w, "func ReadRoot%s(s *C.Segment) %s { return %s(s.Root(0).ToStruct()) }\n",
		n.name, n.name, n.name)
}

func (n *node) defineStructList(w io.Writer) {
	assert(n.Which() == NODE_STRUCT, "invalid struct node")

	fprintf(w, "type %s_List C.PointerList\n", n.name)

	switch n.Struct().PreferredListEncoding() {
	case ELEMENTSIZE_EMPTY:
		fprintf(w, "func New%sList(s *C.Segment, sz int) %s_List { return %s_List(s.NewVoidList(sz)) }\n", n.name, n.name, n.name)
	case ELEMENTSIZE_BIT:
		fprintf(w, "func New%sList(s *C.Segment, sz int) %s_List { return %s_List(s.NewBitList(sz)) }\n", n.name, n.name, n.name)
	case ELEMENTSIZE_BYTE:
		fprintf(w, "func New%sList(s *C.Segment, sz int) %s_List { return %s_List(s.NewUInt8List(sz)) }\n", n.name, n.name, n.name)
	case ELEMENTSIZE_TWOBYTES:
		fprintf(w, "func New%sList(s *C.Segment, sz int) %s_List { return %s_List(s.NewUInt16List(sz)) }\n", n.name, n.name, n.name)
	case ELEMENTSIZE_FOURBYTES:
		fprintf(w, "func New%sList(s *C.Segment, sz int) %s_List { return %s_List(s.NewUInt32List(sz)) }\n", n.name, n.name, n.name)
	case ELEMENTSIZE_EIGHTBYTES:
		fprintf(w, "func New%sList(s *C.Segment, sz int) %s_List { return %s_List(s.NewUInt64List(sz)) }\n", n.name, n.name, n.name)
	default:
		fprintf(w, "func New%sList(s *C.Segment, sz int) %s_List { return %s_List(s.NewCompositeList(%d, %d, sz)) }\n",
			n.name, n.name, n.name, n.Struct().DataWordCount()*8, n.Struct().PointerCount())
	}

	fprintf(w, "func (s %s_List) Len() int { return C.PointerList(s).Len() }\n", n.name)
	fprintf(w, "func (s %s_List) At(i int) %s { return %s(C.PointerList(s).At(i).ToStruct()) }\n", n.name, n.name, n.name)
	fprintf(w, "func (s %s_List) ToArray() []%s { return *(*[]%s)(unsafe.Pointer(C.PointerList(s).ToArray())) }\n", n.name, n.name, n.name)

	g_imported["unsafe"] = true
}

func main() {
	s, err := C.ReadFromStream(os.Stdin, nil)
	assert(err == nil, "%v\n", err)

	req := ReadRootCodeGeneratorRequest(s)
	allfiles := []*node{}

	for _, ni := range req.Nodes().ToArray() {
		n := &node{Node: ni}
		g_nodes[n.Id()] = n

		if n.Which() == NODE_FILE {
			allfiles = append(allfiles, n)
		}
	}

	for _, f := range allfiles {
		for _, a := range f.Annotations().ToArray() {
			if v := a.Value(); v.Which() == VALUE_TEXT {
				switch a.Id() {
				case C.Package:
					f.pkg = v.Text()
				case C.Import:
					f.imp = v.Text()
				}
			}
		}

		for _, nn := range f.NestedNodes().ToArray() {
			if ni := g_nodes[nn.Id()]; ni != nil {
				ni.resolveName("", nn.Name(), f)
			}
		}
	}

	for _, reqf := range req.RequestedFiles().ToArray() {
		f := findNode(reqf.Id())
		buf := bytes.Buffer{}
		g_imported = make(map[string]bool)
		g_segment = C.NewBuffer([]byte{})
		g_bufname = sprintf("x_%x", f.Id())

		for _, n := range f.nodes {
			if n.Which() == NODE_ANNOTATION {
				n.defineAnnotation(&buf)
			}
		}

		defineConstNodes(&buf, f.nodes)

		for _, n := range f.nodes {
			switch n.Which() {
			case NODE_ANNOTATION:
			case NODE_ENUM:
				n.defineEnum(&buf)
				n.defineTypeJsonFuncs(&buf)
			case NODE_STRUCT:
				if !n.Struct().IsGroup() {
					n.defineStructTypes(&buf, nil)
					n.defineStructEnums(&buf)
					n.defineNewStructFunc(&buf)
					n.defineStructFuncs(&buf)
					n.defineTypeJsonFuncs(&buf)
					n.defineStructList(&buf)
				}
			}
		}

		assert(f.pkg != "", "missing package annotation for %s", reqf.Filename())

		file, err := os.Create(reqf.Filename() + ".go")
		assert(err == nil, "%v\n", err)
		fprintf(file, "package %s\n\n", f.pkg)
		fprintf(file, "// AUTO GENERATED - DO NOT EDIT\n\n")

		fprintf(file, "import (\n")
		fprintf(file, "C \"%s\"\n", go_capnproto_import)
		for imp := range g_imported {
			fprintf(file, "%s\n", strconv.Quote(imp))
		}
		fprintf(file, ")\n")

		file.Write(buf.Bytes())

		if len(g_segment.Data) > 0 {
			fprintf(file, "var %s = C.NewBuffer([]byte{", g_bufname)
			for i, b := range g_segment.Data {
				if i%8 == 0 {
					fprintf(file, "\n")
				}
				fprintf(file, "%d,", b)
			}
			fprintf(file, "\n})\n")
		}
		file.Close()

		cmd := exec.Command("gofmt", "-w", reqf.Filename()+".go")
		cmd.Stderr = os.Stderr
		err = cmd.Run()
		assert(err == nil, "%v\n", err)
	}
}
