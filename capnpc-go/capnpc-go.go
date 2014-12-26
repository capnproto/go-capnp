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

	C "github.com/glycerine/go-capnproto"
)

const go_capnproto_import = "github.com/glycerine/go-capnproto"

var (
	g_nodes    = make(map[uint64]*node)
	g_imported map[string]bool
	g_segment  *C.Segment
	g_bufname  string
)

type node struct {
	Node
	pkg   string
	imp   string
	nodes []*node
	name  string
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
	if base == "" {
		n.name = strings.Title(name)
	} else {
		n.name = base + "_" + name
	}
	n.pkg = file.pkg
	n.imp = file.imp

	if n.Which() != Node_Which_struct || !n.Struct().IsGroup() {
		file.nodes = append(file.nodes, n)
	}

	for i := 0; i < n.NestedNodes().Len(); i++ {
		nn := n.NestedNodes().At(i)
		if ni := g_nodes[nn.Id()]; ni != nil {
			ni.resolveName(n.name, nn.Name(), file)
		}
	}

	if n.Which() == Node_Which_struct {
		for i := 0; i < n.Struct().Fields().Len(); i++ {
			f := n.Struct().Fields().At(i)
			if f.Which() == Field_Which_group {
				findNode(f.Group().TypeId()).resolveName(n.name, f.Name(), file)
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
	return e.parent.name + "_" + e.Name()
}

func (n *node) defineEnum(w io.Writer) {
	for i := 0; i < n.Annotations().Len(); i++ {
		a := n.Annotations().At(i)
		if a.Id() == C.Doc {
			fmt.Fprintf(w, "// %s\n", a.Value().Text())
		}
	}
	fmt.Fprintf(w, "type %s uint16\n", n.name)

	if es := n.Enum().Enumerants(); es.Len() > 0 {
		fmt.Fprintf(w, "const (\n")

		ev := make([]enumval, es.Len())
		for i := 0; i < es.Len(); i++ {
			e := es.At(i)

			t := e.Name()
			for i := 0; i < e.Annotations().Len(); i++ {
				an := e.Annotations().At(i)
				if an.Id() == C.Tag {
					t = an.Value().Text()
				} else if an.Id() == C.Notag {
					t = ""
				}
			}
			ev[e.CodeOrder()] = enumval{e, i, t, n}
		}

		// not an iota, so type has to go on each line
		for _, e := range ev {
			fmt.Fprintf(w, "%s %s = %d\n", e.fullName(), n.name, e.val)
		}

		fmt.Fprintf(w, ")\n")

		fmt.Fprintf(w, "func (c %s) String() string {\n", n.name)
		fmt.Fprintf(w, "switch c {\n")
		for _, e := range ev {
			if e.tag != "" {
				fmt.Fprintf(w, "case %s: return \"%s\"\n", e.fullName(), e.tag)
			}
		}
		fmt.Fprintf(w, "default: return \"\"\n")
		fmt.Fprintf(w, "}\n}\n\n")

		fmt.Fprintf(w, "func %sFromString(c string) %s {\n", n.name, n.name)
		fmt.Fprintf(w, "switch c {\n")
		for _, e := range ev {
			if e.tag != "" {
				fmt.Fprintf(w, "case \"%s\": return %s\n", e.tag, e.fullName())
			}
		}
		fmt.Fprintf(w, "default: return 0\n")
		fmt.Fprintf(w, "}\n}\n")
	}

	fmt.Fprintf(w, "type %s_List C.PointerList\n", n.name)
	fmt.Fprintf(w, "func New%s_List(s *C.Segment, sz int) %s_List { return %s_List(s.NewUInt16List(sz)) }\n", n.name, n.name, n.name)
	fmt.Fprintf(w, "func (s %s_List) Len() int { return C.UInt16List(s).Len() }\n", n.name)
	fmt.Fprintf(w, "func (s %s_List) At(i int) %s { return %s(C.UInt16List(s).At(i)) }\n", n.name, n.name, n.name)
}

func (n *node) writeValue(w io.Writer, t Type, v Value) {
	switch t.Which() {
	case Type_Which_void, Type_Which_interface:
		fmt.Fprintf(w, "C.Void{}")

	case Type_Which_bool:
		assert(v.Which() == Value_Which_bool, "expected bool value")
		if v.Bool() {
			fmt.Fprintf(w, "true")
		} else {
			fmt.Fprintf(w, "false")
		}

	case Type_Which_int8:
		assert(v.Which() == Value_Which_int8, "expected int8 value")
		fmt.Fprintf(w, "int8(%d)", v.Int8())

	case Type_Which_uint8:
		assert(v.Which() == Value_Which_uint8, "expected uint8 value")
		fmt.Fprintf(w, "uint8(%d)", v.Uint8())

	case Type_Which_int16:
		assert(v.Which() == Value_Which_int16, "expected int16 value")
		fmt.Fprintf(w, "int16(%d)", v.Int16())

	case Type_Which_uint16:
		assert(v.Which() == Value_Which_uint16, "expected uint16 value")
		fmt.Fprintf(w, "uint16(%d)", v.Uint16())

	case Type_Which_int32:
		assert(v.Which() == Value_Which_int32, "expected int32 value")
		fmt.Fprintf(w, "int32(%d)", v.Int32())

	case Type_Which_uint32:
		assert(v.Which() == Value_Which_uint32, "expected uint32 value")
		fmt.Fprintf(w, "uint32(%d)", v.Uint32())

	case Type_Which_int64:
		assert(v.Which() == Value_Which_int64, "expected int64 value")
		fmt.Fprintf(w, "int64(%d)", v.Int64())

	case Type_Which_uint64:
		assert(v.Which() == Value_Which_uint64, "expected uint64 value")
		fmt.Fprintf(w, "uint64(%d)", v.Uint64())

	case Type_Which_float32:
		assert(v.Which() == Value_Which_float32, "expected float32 value")
		fmt.Fprintf(w, "math.Float32frombits(0x%x)", math.Float32bits(v.Float32()))
		g_imported["math"] = true

	case Type_Which_float64:
		assert(v.Which() == Value_Which_float64, "expected float64 value")
		fmt.Fprintf(w, "math.Float64frombits(0x%x)", math.Float64bits(v.Float64()))
		g_imported["math"] = true

	case Type_Which_text:
		assert(v.Which() == Value_Which_text, "expected text value")
		fmt.Fprintf(w, "%s", strconv.Quote(v.Text()))

	case Type_Which_data:
		assert(v.Which() == Value_Which_data, "expected data value")
		fmt.Fprintf(w, "[]byte{")
		for i, b := range v.Data() {
			if i > 0 {
				fmt.Fprintf(w, ", ")
			}
			fmt.Fprintf(w, "%d", b)
		}
		fmt.Fprintf(w, "}")

	case Type_Which_enum:
		assert(v.Which() == Value_Which_enum, "expected enum value")
		en := findNode(t.Enum().TypeId())
		assert(en.Which() == Node_Which_enum, "expected enum type ID")
		ev := en.Enum().Enumerants()
		if val := int(v.Enum()); val >= ev.Len() {
			fmt.Fprintf(w, "%s(%d)", en.remoteName(n), val)
		} else {
			fmt.Fprintf(w, "%s%s", en.remoteScope(n), ev.At(val).Name())
		}

	case Type_Which_struct:
		fmt.Fprintf(w, "%s(%s.Root(%d))", findNode(t.Struct().TypeId()).remoteName(n), g_bufname, copyData(v.Struct()))

	case Type_Which_anyPointer:
		fmt.Fprintf(w, "%s.Root(%d)", g_bufname, copyData(v.AnyPointer()))

	case Type_Which_list:
		assert(v.Which() == Value_Which_list, "expected list value")

		switch lt := t.List().ElementType(); lt.Which() {
		case Type_Which_void, Type_Which_interface:
			fmt.Fprintf(w, "make([]C.Void, %d)", v.List().ToVoidList().Len())
		case Type_Which_bool:
			fmt.Fprintf(w, "C.BitList(%s.Root(%d))", g_bufname, copyData(v.List()))
		case Type_Which_int8:
			fmt.Fprintf(w, "C.Int8List(%s.Root(%d))", g_bufname, copyData(v.List()))
		case Type_Which_uint8:
			fmt.Fprintf(w, "C.UInt8List(%s.Root(%d))", g_bufname, copyData(v.List()))
		case Type_Which_int16:
			fmt.Fprintf(w, "C.Int16List(%s.Root(%d))", g_bufname, copyData(v.List()))
		case Type_Which_uint16:
			fmt.Fprintf(w, "C.UInt16List(%s.Root(%d))", g_bufname, copyData(v.List()))
		case Type_Which_int32:
			fmt.Fprintf(w, "C.Int32List(%s.Root(%d))", g_bufname, copyData(v.List()))
		case Type_Which_uint32:
			fmt.Fprintf(w, "C.UInt32List(%s.Root(%d))", g_bufname, copyData(v.List()))
		case Type_Which_float32:
			fmt.Fprintf(w, "C.Float32List(%s.Root(%d))", g_bufname, copyData(v.List()))
		case Type_Which_int64:
			fmt.Fprintf(w, "C.Int64List(%s.Root(%d))", g_bufname, copyData(v.List()))
		case Type_Which_uint64:
			fmt.Fprintf(w, "C.UInt64List(%s.Root(%d))", g_bufname, copyData(v.List()))
		case Type_Which_float64:
			fmt.Fprintf(w, "C.Float64List(%s.Root(%d))", g_bufname, copyData(v.List()))
		case Type_Which_text:
			fmt.Fprintf(w, "C.TextList(%s.Root(%d))", g_bufname, copyData(v.List()))
		case Type_Which_data:
			fmt.Fprintf(w, "C.DataList(%s.Root(%d))", g_bufname, copyData(v.List()))
		case Type_Which_enum:
			fmt.Fprintf(w, "%s_List(%s.Root(%d))", findNode(lt.Enum().TypeId()).remoteName(n), g_bufname, copyData(v.List()))
		case Type_Which_struct:
			fmt.Fprintf(w, "%s_List(%s.Root(%d))", findNode(lt.Struct().TypeId()).remoteName(n), g_bufname, copyData(v.List()))
		case Type_Which_list, Type_Which_anyPointer:
			fmt.Fprintf(w, "C.PointerList(%s.Root(%d))", g_bufname, copyData(v.List()))
		}
	}
}

func (n *node) defineAnnotation(w io.Writer) {
	fmt.Fprintf(w, "const %s = uint64(0x%x)\n", n.name, n.Id())
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
			fmt.Fprintf(w, "%s = ", n.name)
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
			fmt.Fprintf(w, "%s = ", n.name)
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

	if t.Which() == Type_Which_interface {
		return
	}

	var g, s bytes.Buffer

	settag := ""
	if f.DiscriminantValue() != Field_noDiscriminant {
		settag = fmt.Sprintf(" C.Struct(s).Set16(%d, %d);", n.Struct().DiscriminantOffset()*2, f.DiscriminantValue())
		if t.Which() == Type_Which_void {
			fmt.Fprintf(&s, "func (s %s) Set%s() {%s }\n", n.name, strings.Title(f.Name()), settag)
			w.Write(s.Bytes())
			return
		}
	} else if t.Which() == Type_Which_void {
		return
	}

	customtype := ""
	for i := 0; i < f.Annotations().Len(); i++ {
		a := f.Annotations().At(i)
		if a.Id() == C.Doc {
			fmt.Fprintf(&g, "// %s\n", a.Value().Text())
		}
		if a.Id() == C.Customtype {
			customtype = a.Value().Text()
			if i := strings.LastIndex(customtype, "."); i != -1 {
				g_imported[customtype[:i]] = true
			}
		}
	}
	fmt.Fprintf(&g, "func (s %s) %s() ", n.name, strings.Title(f.Name()))
	fmt.Fprintf(&s, "func (s %s) Set%s", n.name, strings.Title(f.Name()))

	switch t.Which() {
	case Type_Which_bool:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_bool, "expected bool default")
		if def.Which() == Value_Which_bool && def.Bool() {
			fmt.Fprintf(&g, "bool { return !C.Struct(s).Get1(%d) }\n", off)
			fmt.Fprintf(&s, "(v bool) {%s C.Struct(s).Set1(%d, !v) }\n", settag, off)
		} else {
			fmt.Fprintf(&g, "bool { return C.Struct(s).Get1(%d) }\n", off)
			fmt.Fprintf(&s, "(v bool) {%s C.Struct(s).Set1(%d, v) }\n", settag, off)
		}

	case Type_Which_int8:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_int8, "expected int8 default")
		if def.Which() == Value_Which_int8 && def.Int8() != 0 {
			fmt.Fprintf(&g, "int8 { return int8(C.Struct(s).Get8(%d)) ^ %d }\n", off, def.Int8())
			fmt.Fprintf(&s, "(v int8) {%s C.Struct(s).Set8(%d, uint8(v^%d)) }\n", settag, off, def.Int8())
		} else {
			fmt.Fprintf(&g, "int8 { return int8(C.Struct(s).Get8(%d)) }\n", off)
			fmt.Fprintf(&s, "(v int8) {%s C.Struct(s).Set8(%d, uint8(v)) }\n", settag, off)
		}

	case Type_Which_uint8:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_uint8, "expected uint8 default")
		if def.Which() == Value_Which_uint8 && def.Uint8() != 0 {
			fmt.Fprintf(&g, "uint8 { return C.Struct(s).Get8(%d) ^ %d }\n", off, def.Uint8())
			fmt.Fprintf(&s, "(v uint8) {%s C.Struct(s).Set8(%d, v^%d) }\n", settag, off, def.Uint8())
		} else {
			fmt.Fprintf(&g, "uint8 { return C.Struct(s).Get8(%d) }\n", off)
			fmt.Fprintf(&s, "(v uint8) {%s C.Struct(s).Set8(%d, v) }\n", settag, off)
		}

	case Type_Which_int16:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_int16, "expected int16 default")
		if def.Which() == Value_Which_int16 && def.Int16() != 0 {
			fmt.Fprintf(&g, "int16 { return int16(C.Struct(s).Get16(%d)) ^ %d }\n", off*2, def.Int16())
			fmt.Fprintf(&s, "(v int16) {%s C.Struct(s).Set16(%d, uint16(v^%d)) }\n", settag, off*2, def.Int16())
		} else {
			fmt.Fprintf(&g, "int16 { return int16(C.Struct(s).Get16(%d)) }\n", off*2)
			fmt.Fprintf(&s, "(v int16) {%s C.Struct(s).Set16(%d, uint16(v)) }\n", settag, off*2)
		}

	case Type_Which_uint16:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_uint16, "expected uint16 default")
		if def.Which() == Value_Which_uint16 && def.Uint16() != 0 {
			fmt.Fprintf(&g, "uint16 { return C.Struct(s).Get16(%d) ^ %d }\n", off*2, def.Uint16())
			fmt.Fprintf(&s, "(v uint16) {%s C.Struct(s).Set16(%d, v^%d) }\n", settag, off*2, def.Uint16())
		} else {
			fmt.Fprintf(&g, "uint16 { return C.Struct(s).Get16(%d) }\n", off*2)
			fmt.Fprintf(&s, "(v uint16) {%s C.Struct(s).Set16(%d, v) }\n", settag, off*2)
		}

	case Type_Which_int32:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_int32, "expected int32 default")
		if def.Which() == Value_Which_int32 && def.Int32() != 0 {
			fmt.Fprintf(&g, "int32 { return int32(C.Struct(s).Get32(%d)) ^ %d }\n", off*4, def.Int32())
			fmt.Fprintf(&s, "(v int32) {%s C.Struct(s).Set32(%d, uint32(v^%d)) }\n", settag, off*4, def.Int32())
		} else {
			fmt.Fprintf(&g, "int32 { return int32(C.Struct(s).Get32(%d)) }\n", off*4)
			fmt.Fprintf(&s, "(v int32) {%s C.Struct(s).Set32(%d, uint32(v)) }\n", settag, off*4)
		}

	case Type_Which_uint32:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_uint32, "expected uint32 default")
		if def.Which() == Value_Which_uint32 && def.Uint32() != 0 {
			fmt.Fprintf(&g, "uint32 { return C.Struct(s).Get32(%d) ^ %d }\n", off*4, def.Uint32())
			fmt.Fprintf(&s, "(v uint32) {%s C.Struct(s).Set32(%d, v^%d) }\n", settag, off*4, def.Uint32())
		} else {
			fmt.Fprintf(&g, "uint32 { return C.Struct(s).Get32(%d) }\n", off*4)
			fmt.Fprintf(&s, "(v uint32) {%s C.Struct(s).Set32(%d, v) }\n", settag, off*4)
		}

	case Type_Which_int64:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_int64, "expected int64 default")
		if def.Which() == Value_Which_int64 && def.Int64() != 0 {
			fmt.Fprintf(&g, "int64 { return int64(C.Struct(s).Get64(%d)) ^ %d }\n", off*8, def.Int64())
			fmt.Fprintf(&s, "(v int64) {%s C.Struct(s).Set64(%d, uint64(v^%d)) }\n", settag, off*8, def.Int64())
		} else {
			fmt.Fprintf(&g, "int64 { return int64(C.Struct(s).Get64(%d)) }\n", off*8)
			fmt.Fprintf(&s, "(v int64) {%s C.Struct(s).Set64(%d, uint64(v)) }\n", settag, off*8)
		}

	case Type_Which_uint64:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_uint64, "expected uint64 default")
		if def.Which() == Value_Which_uint64 && def.Uint64() != 0 {
			fmt.Fprintf(&g, "uint64 { return C.Struct(s).Get64(%d) ^ %d }\n", off*8, def.Uint64())
			fmt.Fprintf(&s, "(v uint64) {%s C.Struct(s).Set64(%d, v^%d) }\n", settag, off*8, def.Uint64())
		} else {
			fmt.Fprintf(&g, "uint64 { return C.Struct(s).Get64(%d) }\n", off*8)
			fmt.Fprintf(&s, "(v uint64) {%s C.Struct(s).Set64(%d, v) }\n", settag, off*8)
		}

	case Type_Which_float32:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_float32, "expected float32 default")
		if def.Which() == Value_Which_float32 && def.Float32() != 0 {
			fmt.Fprintf(&g, "float32 { return math.Float32frombits(C.Struct(s).Get32(%d) ^ 0x%x) }\n", off*4, math.Float32bits(def.Float32()))
			fmt.Fprintf(&s, "(v float32) {%s C.Struct(s).Set32(%d, math.Float32bits(v) ^ 0x%x) }\n", settag, off*4, math.Float32bits(def.Float32()))
		} else {
			fmt.Fprintf(&g, "float32 { return math.Float32frombits(C.Struct(s).Get32(%d)) }\n", off*4)
			fmt.Fprintf(&s, "(v float32) {%s C.Struct(s).Set32(%d, math.Float32bits(v)) }\n", settag, off*4)
		}
		g_imported["math"] = true

	case Type_Which_float64:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_float64, "expected float64 default")
		if def.Which() == Value_Which_float64 && def.Float64() != 0 {
			fmt.Fprintf(&g, "float64 { return math.Float64frombits(C.Struct(s).Get64(%d) ^ 0x%x) }\n", off*8, math.Float64bits(def.Float64()))
			fmt.Fprintf(&s, "(v float64) {%s C.Struct(s).Set64(%d, math.Float64bits(v) ^ 0x%x) }\n", settag, off*8, math.Float64bits(def.Float64()))
		} else {
			fmt.Fprintf(&g, "float64 { return math.Float64frombits(C.Struct(s).Get64(%d)) }\n", off*8)
			fmt.Fprintf(&s, "(v float64) {%s C.Struct(s).Set64(%d, math.Float64bits(v)) }\n", settag, off*8)
		}
		g_imported["math"] = true

	case Type_Which_text:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_text, "expected text default")
		if def.Which() == Value_Which_text && def.Text() != "" {
			fmt.Fprintf(&g, "string { return C.Struct(s).GetObject(%d).ToTextDefault(%s) }\n", off, strconv.Quote(def.Text()))
		} else {
			fmt.Fprintf(&g, "string { return C.Struct(s).GetObject(%d).ToText() }\n", off)
		}
		fmt.Fprintf(&s, "(v string) {%s C.Struct(s).SetObject(%d, s.Segment.NewText(v)) }\n", settag, off)

	case Type_Which_data:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_data, "expected data default")
		if def.Which() == Value_Which_data && len(def.Data()) > 0 {
			dstr := "[]byte{"
			for i, b := range def.Data() {
				if i > 0 {
					dstr += ", "
				}
				dstr += fmt.Sprintf("%d", b)
			}
			dstr += "}"
			if len(customtype) != 0 {
				fmt.Fprintf(&g, "%s { return %s(C.Struct(s).GetObject(%d).ToDataDefault(%s)) }\n", customtype, customtype, off, dstr)
			} else {
				fmt.Fprintf(&g, "[]byte { return C.Struct(s).GetObject(%d).ToDataDefault(%s) }\n", off, dstr)
			}
		} else {
			if len(customtype) != 0 {
				fmt.Fprintf(&g, "%s { return %s(C.Struct(s).GetObject(%d).ToData()) }\n", customtype, customtype, off)
			} else {
				fmt.Fprintf(&g, "[]byte { return C.Struct(s).GetObject(%d).ToData() }\n", off)
			}
		}
		if len(customtype) != 0 {
			fmt.Fprintf(&s, "(v %s) {%s C.Struct(s).SetObject(%d, s.Segment.NewData([]byte(v))) }\n", customtype, settag, off)
		} else {
			fmt.Fprintf(&s, "(v []byte) {%s C.Struct(s).SetObject(%d, s.Segment.NewData(v)) }\n", settag, off)
		}

	case Type_Which_enum:
		ni := findNode(t.Enum().TypeId())
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_enum, "expected enum default")
		if def.Which() == Value_Which_enum && def.Enum() != 0 {
			fmt.Fprintf(&g, "%s { return %s(C.Struct(s).Get16(%d) ^ %d) }\n", ni.remoteName(n), ni.remoteName(n), off*2, def.Enum())
			fmt.Fprintf(&s, "(v %s) {%s C.Struct(s).Set16(%d, uint16(v)^%d) }\n", ni.remoteName(n), settag, off*2, def.Uint16())
		} else {
			fmt.Fprintf(&g, "%s { return %s(C.Struct(s).Get16(%d)) }\n", ni.remoteName(n), ni.remoteName(n), off*2)
			fmt.Fprintf(&s, "(v %s) {%s C.Struct(s).Set16(%d, uint16(v)) }\n", ni.remoteName(n), settag, off*2)
		}

	case Type_Which_struct:
		ni := findNode(t.Struct().TypeId())
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_struct, "expected struct default")
		if def.Which() == Value_Which_struct && def.Struct().HasData() {
			fmt.Fprintf(&g, "%s { return %s(C.Struct(s).GetObject(%d).ToStructDefault(%s, %d)) }\n",
				ni.remoteName(n), ni.remoteName(n), off, g_bufname, copyData(def.Struct()))
		} else {
			fmt.Fprintf(&g, "%s { return %s(C.Struct(s).GetObject(%d).ToStruct()) }\n",
				ni.remoteName(n), ni.remoteName(n), off)
		}
		fmt.Fprintf(&s, "(v %s) {%s C.Struct(s).SetObject(%d, C.Object(v)) }\n", ni.remoteName(n), settag, off)

	case Type_Which_anyPointer:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_anyPointer, "expected object default")
		if def.Which() == Value_Which_anyPointer && def.AnyPointer().HasData() {
			fmt.Fprintf(&g, "C.Object { return C.Struct(s).GetObject(%d).ToObjectDefault(%s, %d) }\n",
				off, g_bufname, copyData(def.AnyPointer()))
		} else {
			fmt.Fprintf(&g, "C.Object { return C.Struct(s).GetObject(%d) }\n", off)
		}
		fmt.Fprintf(&s, "(v C.Object) {%s C.Struct(s).SetObject(%d, v) }\n", settag, off)

	case Type_Which_list:
		assert(def.Which() == Value_Which_void || def.Which() == Value_Which_list, "expected list default")

		typ := ""

		switch lt := t.List().ElementType(); lt.Which() {
		case Type_Which_void, Type_Which_interface:
			typ = "C.VoidList"
		case Type_Which_bool:
			typ = "C.BitList"
		case Type_Which_int8:
			typ = "C.Int8List"
		case Type_Which_uint8:
			typ = "C.UInt8List"
		case Type_Which_int16:
			typ = "C.Int16List"
		case Type_Which_uint16:
			typ = "C.UInt16List"
		case Type_Which_int32:
			typ = "C.Int32List"
		case Type_Which_uint32:
			typ = "C.UInt32List"
		case Type_Which_int64:
			typ = "C.Int64List"
		case Type_Which_uint64:
			typ = "C.UInt64List"
		case Type_Which_float32:
			typ = "C.Float32List"
		case Type_Which_float64:
			typ = "C.Float64List"
		case Type_Which_text:
			typ = "C.TextList"
		case Type_Which_data:
			typ = "C.DataList"
		case Type_Which_enum:
			ni := findNode(lt.Enum().TypeId())
			typ = fmt.Sprintf("%s_List", ni.remoteName(n))
		case Type_Which_struct:
			ni := findNode(lt.Struct().TypeId())
			typ = fmt.Sprintf("%s_List", ni.remoteName(n))
		case Type_Which_anyPointer, Type_Which_list:
			typ = "C.PointerList"
		}

		ldef := C.Object{}
		if def.Which() == Value_Which_list {
			ldef = def.List()
		}

		if ldef.HasData() {
			fmt.Fprintf(&g, "%s { return %s(C.Struct(s).GetObject(%d).ToListDefault(%s, %d)) }\n",
				typ, typ, off, g_bufname, copyData(ldef))
		} else {
			fmt.Fprintf(&g, "%s { return %s(C.Struct(s).GetObject(%d)) }\n",
				typ, typ, off)
		}

		fmt.Fprintf(&s, "(v %s) {%s C.Struct(s).SetObject(%d, C.Object(v)) }\n", typ, settag, off)
	}

	w.Write(g.Bytes())
	w.Write(s.Bytes())
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

	for i := 0; i < n.Annotations().Len(); i++ {
		a := n.Annotations().At(i)
		if a.Id() == C.Doc {
			fmt.Fprintf(w, "// %s\n", a.Value().Text())
		}
	}
	if baseNode != nil {
		fmt.Fprintf(w, "type %s %s\n", n.name, baseNode.name)
	} else {
		fmt.Fprintf(w, "type %s C.Struct\n", n.name)
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
		fmt.Fprintf(w, "type %s_Which uint16\n", n.name)
		fmt.Fprintf(w, "const (\n")

		for _, f := range n.codeOrderFields() {
			if f.DiscriminantValue() == Field_noDiscriminant {
				// Non-union member
			} else {
				fmt.Fprintf(w, "%s_Which_%s %s_Which = %d\n", n.name, f.Name(), n.name, f.DiscriminantValue())
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
		fmt.Fprintf(w, "func (s %s) Which() %s_Which { return %s_Which(C.Struct(s).Get16(%d)) }\n",
			n.name, n.name, n.name, n.Struct().DiscriminantOffset()*2)
	}

	for _, f := range n.codeOrderFields() {
		switch f.Which() {
		case Field_Which_slot:
			n.defineField(w, f)
		case Field_Which_group:
			g := findNode(f.Group().TypeId())
			fmt.Fprintf(w, "func (s %s) %s() %s { return %s(s) }\n", n.name, strings.Title(f.Name()), g.name, g.name)
			if f.DiscriminantValue() != Field_noDiscriminant {
				fmt.Fprintf(w, "func (s %s) Set%s() { C.Struct(s).Set16(%d, %d) }\n", n.name, strings.Title(f.Name()), n.Struct().DiscriminantOffset()*2, f.DiscriminantValue())
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
		g_imported["io"] = true
		g_imported["bufio"] = true
		g_imported["bytes"] = true

		fmt.Fprintf(w, "func (s %s) WriteJSON(w io.Writer) error {\n", n.name)
		fmt.Fprintf(w, "b := bufio.NewWriter(w);")
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

		fmt.Fprintf(w, "func (s %s) MarshalJSON() ([]byte, error) {\n", n.name)
		fmt.Fprintf(w, "b := bytes.Buffer{}; err := s.WriteJSON(&b); return b.Bytes(), err };")

	} else {
		fmt.Fprintf(w, "// capn.JSON_enabled == false so we stub MarshalJSON().")
		fmt.Fprintf(w, "\nfunc (s %s) MarshalJSON() (bs []byte, err error) { return } \n", n.name)
	}
}

func writeErrCheck(w io.Writer) {
	fmt.Fprintf(w, "if err != nil { return err; };")
}

func (n *node) jsonEnum(w io.Writer) {
	g_imported["encoding/json"] = true
	fmt.Fprintf(w, `buf, err = json.Marshal(s.String());`)
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
			enumname := n.name + "_Which_" + f.Name()
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
		g_imported["encoding/json"] = true
		fmt.Fprintf(w, "buf, err = json.Marshal(s);")
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

func (n *node) defineNewStructFunc(w io.Writer) {
	assert(n.Which() == Node_Which_struct, "invalid struct node")

	fmt.Fprintf(w, "func New%s(s *C.Segment) %s { return %s(s.NewStruct(C.ObjectSize{DataSize: %d, PointerCount: %d})) }\n",
		n.name, n.name, n.name, n.Struct().DataWordCount()*8, n.Struct().PointerCount())
	fmt.Fprintf(w, "func NewRoot%s(s *C.Segment) %s { return %s(s.NewRootStruct(C.ObjectSize{DataSize: %d, PointerCount: %d})) }\n",
		n.name, n.name, n.name, n.Struct().DataWordCount()*8, n.Struct().PointerCount())
	fmt.Fprintf(w, "func AutoNew%s(s *C.Segment) %s { return %s(s.NewStructAR(C.ObjectSize{DataSize: %d, PointerCount: %d})) }\n",
		n.name, n.name, n.name, n.Struct().DataWordCount()*8, n.Struct().PointerCount())
	fmt.Fprintf(w, "func ReadRoot%s(s *C.Segment) %s { return %s(s.Root(0).ToStruct()) }\n",
		n.name, n.name, n.name)
}

func (n *node) defineStructList(w io.Writer) {
	assert(n.Which() == Node_Which_struct, "invalid struct node")

	fmt.Fprintf(w, "type %s_List C.PointerList\n", n.name)

	fmt.Fprintf(w, "func New%s_List(s *C.Segment, sz int) %s_List ", n.name, n.name)
	switch n.Struct().PreferredListEncoding() {
	case ElementSize_empty:
		fmt.Fprintf(w, "{ return %s_List(s.NewVoidList(sz)) }\n", n.name)
	case ElementSize_bit:
		fmt.Fprintf(w, "{ return %s_List(s.NewBitList(sz)) }\n", n.name)
	case ElementSize_byte:
		fmt.Fprintf(w, "{ return %s_List(s.NewUInt8List(sz)) }\n", n.name)
	case ElementSize_twoBytes:
		fmt.Fprintf(w, "{ return %s_List(s.NewUInt16List(sz)) }\n", n.name)
	case ElementSize_fourBytes:
		fmt.Fprintf(w, "{ return %s_List(s.NewUInt32List(sz)) }\n", n.name)
	case ElementSize_eightBytes:
		fmt.Fprintf(w, "{ return %s_List(s.NewUInt64List(sz)) }\n", n.name)
	default:
		fmt.Fprintf(w, "{ return %s_List(s.NewCompositeList(C.ObjectSize{DataSize: %d, PointerCount: %d}, sz)) }\n",
			n.name, n.Struct().DataWordCount()*8, n.Struct().PointerCount())
	}

	fmt.Fprintf(w, "func (s %s_List) Len() int { return C.PointerList(s).Len() }\n", n.name)
	fmt.Fprintf(w, "func (s %s_List) At(i int) %s { return %s(C.PointerList(s).At(i).ToStruct()) }\n", n.name, n.name, n.name)
	fmt.Fprintf(w, "func (s %s_List) Set(i int, item %s) { C.PointerList(s).Set(i, C.Object(item)) }\n", n.name, n.name)
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
	g_imported = make(map[string]bool)
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
		case Node_Which_annotation:
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
			}
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
	fmt.Fprintf(&unformatted, "C %q\n", go_capnproto_import)
	for imp := range g_imported {
		fmt.Fprintf(&unformatted, "%q\n", imp)
	}
	fmt.Fprintf(&unformatted, ")\n")
	unformatted.Write(buf.Bytes())
	if len(g_segment.Data) > 0 {
		fmt.Fprintf(&unformatted, "var %s = C.NewBuffer([]byte{", g_bufname)
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
		for i := 0; i < f.Annotations().Len(); i++ {
			a := f.Annotations().At(i)
			if v := a.Value(); v.Which() == Value_Which_text {
				switch a.Id() {
				case C.Package:
					f.pkg = v.Text()
				case C.Import:
					f.imp = v.Text()
				}
			}
		}

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
