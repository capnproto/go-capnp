package main

import (
	"fmt"
	C "github.com/jmckaskill/go-capnproto"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

const importPath = "github.com/jmckaskill/go-capnproto"

type GoFile struct {
	buf       *C.Segment
	bufName   string
	constants []*field
	types     []*typ
	vals      []string
	valbase   string
}

func pfxname(pfx, name string) string {
	r, rn := utf8.DecodeRuneInString(name)
	if unicode.IsUpper(r) {
		p, pn := utf8.DecodeRuneInString(pfx)
		return string(unicode.ToUpper(p)) + pfx[pn:] + name
	} else {
		return pfx + string(unicode.ToUpper(r)) + name[rn:]
	}
}

func (g *GoFile) resolveTypes() {
	for _, t := range g.types {
		switch t.typ {
		case structType, interfaceType:
			t.name = strings.Replace(t.name, "·", "_", -1)
		case enumType:
			t.name = strings.Replace(t.name, "·", "_", -1)
			t.enumPrefix = t.name + "_"
		case unionType:
			findex := strings.LastIndex(t.name, "·")
			t.name = strings.Replace(t.name, "·", "_", -1)
			t.enumPrefix = t.name[:findex] + "_"
		case methodType, returnType:
			/* not used directly */
		case voidType:
			t.name = "struct{}"
		case boolType:
			t.name = "bool"
		case int8Type:
			t.name = "int8"
		case int16Type:
			t.name = "int16"
		case int32Type:
			t.name = "int32"
		case int64Type:
			t.name = "int64"
		case uint8Type:
			t.name = "uint8"
		case uint16Type:
			t.name = "uint16"
		case uint32Type:
			t.name = "uint32"
		case uint64Type:
			t.name = "uint64"
		case float32Type:
			t.name = "float32"
		case float64Type:
			t.name = "float64"
		case stringType:
			t.name = "string"
		case dataType:
			t.name = "[]byte"
		case listType:
			// We cap the max list depth at one for the go types
			switch t.listType.typ {
			case listType, stringType, dataType:
				t.name = "[]C.Pointer"
			case boolType:
				t.name = "C.Bitset"
			default:
				t.name = "[]" + t.listType.name
			}
		default:
			panic("unhandled")
		}
	}
}

func (g *GoFile) GoString(v *value, symbol string) string {
	t := v.typ

	switch t.typ {
	case stringType:
		return strconv.Quote(v.string)
	case enumType:
		return sprintf("%s(%d)", t.name, v.num)
	case boolType:
		if v.bool {
			return "true"
		} else {
			return "false"
		}
	case float32Type, float64Type:
		return sprintf("%v", v.float)
	case uint8Type, uint16Type, uint32Type, uint64Type:
		return sprintf("%d", uint64(v.num))
	case int8Type, int16Type, int32Type, int64Type:
		return sprintf("%d", v.num)
	}

	if v.symbol != "" {
		return v.symbol
	}

	s := ""

	switch t.typ {
	case dataType:
		if v.tok == '"' {
			s = "[]byte(" + strconv.Quote(v.string) + ")"
		} else {
			s = "[]byte{"
			for i, v := range v.array {
				if i > 0 {
					s += ", "
				}
				s += g.GoString(v, "")
			}
			s += "}"
		}

	case listType:
		switch t.listType.typ {
		case voidType:
			s = sprintf("make([]struct{}, %d)", len(v.array))
		case boolType:
			s = sprintf("%s.ReadBitset(%d, C.Bitset{})", g.bufName, v.Marshal(g.buf))
		case listType, stringType, dataType:
			s = "[]C.Pointer{"
			for i, v := range v.array {
				if i > 0 {
					s += ", "
				}
				s += sprintf("%s.ReadPtr(%d)", g.bufName, v.Marshal(g.buf))
			}
			s += "}"
		default:
			s = t.name + "{"
			for i, v := range v.array {
				if i > 0 {
					s += ", "
				}
				s += g.GoString(v, "")
			}
			s += "}"
		}

	case structType:
		s = sprintf("%s{%s.ReadPtr(%d)}", t.name, g.bufName, v.Marshal(g.buf))
	default:
		panic("unhandled")
	}

	if symbol != "" {
		v.symbol = symbol
		return s
	}

	v.symbol = sprintf("%s_%d", g.valbase, len(g.vals))
	g.vals = append(g.vals, s)
	return v.symbol

	panic("unhandled")
}

func (g *GoFile) getter(f *field, ptr string) string {
	xor := ""
	def := "nil"
	if f.value != nil {
		def = g.GoString(f.value, "")
		xor = " ^ " + def
	}

	switch f.typ.typ {
	case boolType:
		if f.value != nil && f.value.bool {
			return sprintf("!%s.ReadStruct1(%d)", ptr, f.offset)
		} else {
			return sprintf("%s.ReadStruct1(%d)", ptr, f.offset)
		}
	case int8Type:
		return sprintf("int8(%s.ReadStruct8(%d))%s", ptr, f.offset, xor)
	case uint8Type:
		return sprintf("%s.ReadStruct8(%d)%s", ptr, f.offset, xor)
	case int16Type:
		return sprintf("int16(%s.ReadStruct16(%d))%s", ptr, f.offset, xor)
	case uint16Type:
		return sprintf("%s.ReadStruct16(%d)%s", ptr, f.offset, xor)
	case int32Type:
		return sprintf("int32(%s.ReadStruct32(%d))%s", ptr, f.offset, xor)
	case uint32Type:
		return sprintf("%s.ReadStruct32(%d)%s", ptr, f.offset, xor)
	case int64Type:
		return sprintf("int64(%s.ReadStruct64(%d))%s", ptr, f.offset, xor)
	case uint64Type:
		return sprintf("%s.ReadStruct64(%d)%s", ptr, f.offset, xor)
	case enumType, unionType:
		if f.value != nil {
			return sprintf("%s(%s.ReadStruct16(%d)^%d)", f.typ.name, ptr, f.offset, f.value.num)
		} else {
			return sprintf("%s(%s.ReadStruct16(%d))", f.typ.name, ptr, f.offset)
		}

	case float32Type:
		if f.value != nil {
			return sprintf("%s.ReadStructF32(%d, %s)", ptr, f.offset, def)
		} else {
			return sprintf("M.Float32frombits(%s.ReadStruct32(%d))", ptr, f.offset)
		}

	case float64Type:
		if f.value != nil {
			return sprintf("%s.ReadStructF64(%d, %s)", ptr, f.offset, def)
		} else {
			return sprintf("M.Float64frombits(%s.ReadStruct64(%d))", ptr, f.offset)
		}

	case dataType:
		return sprintf("%s.ReadData(%d, %s)", ptr, f.offset, def)

	case stringType:
		if f.value != nil {
			return sprintf("%s.ReadString(%d, %s)", ptr, f.offset, def)
		} else {
			return sprintf("%s.ReadString(%d, \"\")", ptr, f.offset)
		}

	case structType:
		if f.value != nil {
			return sprintf("%s{%s.ReadStruct(%d, %s.P)}", f.typ.name, ptr, f.offset, def)
		} else {
			return sprintf("%s{%s.ReadPtr(%d)}", f.typ.name, ptr, f.offset)
		}

	case interfaceType:
		return sprintf("%s{%s.ReadPtr(%d)}", pfxname("remote", f.typ.name), ptr, f.offset)

	case listType:
		switch lt := f.typ.listType; lt.typ {
		case voidType:
			return sprintf("%s.ReadVoidList(%d, %s)", ptr, f.offset, def)
		case boolType:
			if f.value != nil {
				return sprintf("%s.ReadBitset(%d, %s)", ptr, f.offset, def)
			} else {
				return sprintf("%s.ReadBitset(%d, C.Bitset{})", ptr, f.offset)
			}
		case int8Type:
			return sprintf("%s.ReadI8List(%d, %s)", ptr, f.offset, def)
		case uint8Type:
			return sprintf("%s.ReadU8List(%d, %s)", ptr, f.offset, def)
		case int16Type:
			return sprintf("%s.ReadI16List(%d, %s)", ptr, f.offset, def)
		case uint16Type:
			return sprintf("%s.ReadU16List(%d, %s)", ptr, f.offset, def)
		case int32Type:
			return sprintf("%s.ReadI32List(%d, %s)", ptr, f.offset, def)
		case uint32Type:
			return sprintf("%s.ReadU32List(%d, %s)", ptr, f.offset, def)
		case int64Type:
			return sprintf("%s.ReadI64List(%d, %s)", ptr, f.offset, def)
		case uint64Type:
			return sprintf("%s.ReadU64List(%d, %s)", ptr, f.offset, def)
		case float32Type:
			return sprintf("%s.ReadF32List(%d, %s)", ptr, f.offset, def)
		case float64Type:
			return sprintf("%s.ReadF64List(%d, %s)", ptr, f.offset, def)
		case enumType, structType, interfaceType:
			return sprintf("%sList(%s, %d, %s)", pfxname("read", lt.name), ptr, f.offset, def)
		case stringType, listType, dataType:
			return sprintf("%s.ReadPointerList(%d, %s)", ptr, f.offset, def)
		}
	}

	panic("unhandled")
}

func (g *GoFile) setter(f *field, ptr, arg string) string {
	def := "nil"
	if f.value != nil {
		switch f.typ.typ {
		case int8Type, uint8Type, int16Type, uint16Type, int32Type, uint32Type, int64Type, uint64Type:
			arg = sprintf("%s ^ %s", arg, g.GoString(f.value, ""))
		case boolType:
			if f.value.bool {
				arg = sprintf("!%s", arg)
			}
		case stringType, dataType, listType, structType, float32Type, float64Type:
			def = g.GoString(f.value, "")
		case enumType:
			// handled below
		default:
			panic("unhandled")
		}
	}

	switch f.typ.typ {
	case boolType:
		return sprintf("%s.WriteStruct1(%d, %s)", ptr, f.offset, arg)
	case int8Type:
		return sprintf("%s.WriteStruct8(%d, uint8(%s))", ptr, f.offset, arg)
	case uint8Type:
		return sprintf("%s.WriteStruct8(%d, %s)", ptr, f.offset, arg)
	case int16Type:
		return sprintf("%s.WriteStruct16(%d, uint16(%s))", ptr, f.offset, arg)
	case uint16Type:
		return sprintf("%s.WriteStruct16(%d, %s)", ptr, f.offset, arg)
	case enumType:
		if f.value != nil {
			return sprintf("%s.WriteStruct16(%d, uint16(%s)^%d)", ptr, f.offset, arg, f.value.num)
		} else {
			return sprintf("%s.WriteStruct16(%d, uint16(%s))", ptr, f.offset, arg)
		}
	case int32Type:
		return sprintf("%s.WriteStruct32(%d, uint32(%s))", ptr, f.offset, arg)
	case uint32Type:
		return sprintf("%s.WriteStruct32(%d, %s)", ptr, f.offset, arg)
	case int64Type:
		return sprintf("%s.WriteStruct64(%d, uint64(%s))", ptr, f.offset, arg)
	case uint64Type:
		return sprintf("%s.WriteStruct64(%d, %s)", ptr, f.offset, arg)
	case float32Type:
		if f.value != nil {
			return sprintf("%s.WriteStructF32(%d, %s, %s)", ptr, f.offset, arg, def)
		} else {
			return sprintf("%s.WriteStruct32(%d, M.Float32bits(%s))", ptr, f.offset, arg)
		}
	case float64Type:
		if f.value != nil {
			return sprintf("%s.WriteStructF64(%d, %s, %s)", ptr, f.offset, arg, def)
		} else {
			return sprintf("%s.WriteStruct64(%d, M.Float64bits(%s))", ptr, f.offset, arg)
		}
	case stringType:
		if f.value != nil {
			return sprintf("%s.WriteString(%d, %s, %s)", ptr, f.offset, arg, def)
		} else {
			return sprintf("%s.WriteString(%d, %s, \"\")", ptr, f.offset, arg)
		}
	case dataType:
		return sprintf("%s.WriteU8List(%d, %s, %s)", ptr, f.offset, arg, def)
	case interfaceType:
		return sprintf("%s.MarshalCaptain(%s, %d)", arg, ptr, f.offset)
	case structType:
		if f.value != nil {
			return sprintf("%s.WriteStruct(%d, %s.P, %s.P)", ptr, f.offset, arg, def)
		} else {
			return sprintf("%s.WritePtr(%d, %s.P)", ptr, f.offset, arg)
		}
	case listType:
		lt := f.typ.listType
		switch lt.typ {
		case voidType:
			return sprintf("%s.WriteVoidList(%d, %s, %s)", ptr, f.offset, arg, def)
		case boolType:
			if f.value != nil {
				return sprintf("%s.WriteBitset(%d, %s, %s)", ptr, f.offset, arg, def)
			} else {
				return sprintf("%s.WriteBitset(%d, %s, C.Bitset{})", ptr, f.offset, arg)
			}
		case int8Type:
			return sprintf("%s.WriteI8List(%d, %s, %s)", ptr, f.offset, arg, def)
		case uint8Type:
			return sprintf("%s.WriteU8List(%d, %s, %s)", ptr, f.offset, arg, def)
		case int16Type:
			return sprintf("%s.WriteI16List(%d, %s, %s)", ptr, f.offset, arg, def)
		case uint16Type:
			return sprintf("%s.WriteU16List(%d, %s, %s)", ptr, f.offset, arg, def)
		case int32Type:
			return sprintf("%s.WriteI32List(%d, %s, %s)", ptr, f.offset, arg, def)
		case uint32Type:
			return sprintf("%s.WriteU32List(%d, %s, %s)", ptr, f.offset, arg, def)
		case int64Type:
			return sprintf("%s.WriteI64List(%d, %s, %s)", ptr, f.offset, arg, def)
		case uint64Type:
			return sprintf("%s.WriteU64List(%d, %s, %s)", ptr, f.offset, arg, def)
		case float32Type:
			return sprintf("%s.WriteF32List(%d, %s, %s)", ptr, f.offset, arg, def)
		case float64Type:
			return sprintf("%s.WriteF64List(%d, %s, %s)", ptr, f.offset, arg, def)
		case enumType, structType, interfaceType:
			return sprintf("%sList(%s, %d, %s, %s)", pfxname("write", lt.name), ptr, f.offset, arg, def)
		case stringType, listType, dataType:
			return sprintf("%s.WritePointerList(%d, %s, %s)", ptr, f.offset, arg, def)
		default:
			panic("unhandled")
		}

	}

	panic("unhandled")
}

func (g *GoFile) writeStruct(t *typ) {
	out("type %s struct { P C.Pointer }\n", t.name)

	out("func %s(seg *C.Segment) (%s, error) {\n", pfxname("new", t.name), t.name)
	out("p, err := seg.NewStruct(%d, %d)\n", t.dataSize, t.ptrSize)
	out("return %s{p}, err\n", t.name)
	out("}\n")

	out("func %sList(seg *C.Segment, sz int) (C.Pointer, error) {\n", pfxname("new", t.name))
	out("return seg.NewList(%d, %d, sz)\n", t.dataSize, t.ptrSize)
	out("}\n")

	out("func (p %s) MarshalCaptain(r C.Pointer, i int) error {\n", t.name)
	out("return r.WritePtr(i, p.P)\n")
	out("}\n")

	out("\n")

	for _, f := range t.fields {
		if f.typ.typ == voidType {
			continue
		}

		if len(f.comment) > 0 {
			out("/* %s */\n", f.comment)
		}

		if f.union != nil {
			out("func (p %s) %s() (ret %s) {\n", t.name, f.name, f.typ.name)
			out("if p.P.ReadStruct16(%d) != %d {\n", f.union.offset, f.ordinal)
			if f.value != nil {
				out("return %s\n", g.GoString(f.value, ""))
			} else {
				out("return\n")
			}
			out("}\n")
			out("return %s\n", g.getter(f, "p.P"))
			out("}\n")
		} else {
			out("func (p %s) %s() %s { return %s }\n",
				t.name, f.name, f.typ.name,
				g.getter(f, "p.P"))
		}
	}

	out("\n")

	for _, f := range t.fields {
		if f.typ.typ != unionType && f.typ.typ != voidType {
			if len(f.comment) > 0 {
				out("/* %s */\n", f.comment)
			}

			if f.union != nil {
				out("func (p %s) %s(v %s) error {", t.name, pfxname("set", f.name), f.typ.name)
				out("if err := p.P.WriteStruct16(%d, %d); err != nil {\n", f.union.offset, f.ordinal)
				out("return err\n")
				out("}\n")
				out("return %s", g.setter(f, "p.P", "v"))
				out("}\n")
			} else {
				out("func (p %s) %s(v %s) error { return %s }\n",
					t.name, pfxname("set", f.name), f.typ.name,
					g.setter(f, "p.P", "v"))
			}
		}
	}

	out("\n")
}

func (g *GoFile) writeInterface(t *typ) {
	out("type %s interface {\n", t.name)
	out("C.Marshaller\n")

	for _, method := range t.fields {
		// Interface method declaration
		mt := method.typ
		mr := method.ret

		if len(method.comment) > 0 {
			out("/* %s */\n", method.comment)
		}

		out("%s(", method.name)
		for ai, a := range mt.fields {
			if ai > 0 {
				out(", ")
			}
			out("%s %s", a.name, a.typ.name)
		}

		switch len(mr.fields) {
		case 0:
			out(")\n")
		case 1:
			out(") %s\n", mr.fields[0].typ.name)
		default:
			out(") (")
			for ai, a := range mr.fields {
				if ai > 0 {
					out(", ")
				}
				out("%s %s", a.name, a.typ.name)
			}
			out(")\n")
		}
	}

	out("}\n")

	out("\n")

	out("type %s struct { P C.Pointer }\n", pfxname("remote", t.name))
	out("func (p %s) MarshalCaptain(r C.Pointer, i int) error {\n", pfxname("remote", t.name))
	out("return r.WritePtr(i, p.P)\n")
	out("}\n")

	for _, method := range t.fields {
		// Remote interface method
		mt := method.typ
		mr := method.ret

		out("func (p %s) %s(", pfxname("remote", t.name), method.name)
		for ai, a := range mt.fields {
			if ai > 0 {
				out(", ")
			}
			out("a%d %s", ai, a.typ.name)
		}
		out(") ")

		switch len(mr.fields) {
		case 0:
			out("{\n")
		default:
			out("(")
			for ai, a := range mr.fields {
				if ai > 0 {
					out(", ")
				}
				out("r%d %s", ai, a.typ.name)
			}
			out(") ")
			out("{\n")
			for ai, a := range mr.fields {
				if a.value != nil {
					out("r%d = %s\n", ai, g.GoString(a.value, ""))
				}
			}
		}

		out("c, err := p.P.Segment.Session.NewCall()\n")
		out("if err != nil { return }\n")
		out("c.Message.SetObject(p.P)\n")
		out("c.Message.SetMethod(%d)\n", method.ordinal)

		if len(mt.fields) > 0 {
			out("args, _ := c.Message.P.Segment.NewStruct(%d, %d)\n", mt.dataSize, mt.ptrSize)
			for ai, a := range mt.fields {
				out("%s\n", g.setter(a, "args", sprintf("a%d", ai)))
			}
			out("c.Message.SetArguments(args)\n")
		}

		switch len(mr.fields) {
		case 0:
			out("c.Send(c)\n")
		default:
			out("if c.Send(c) != nil { return }\n")
			out("return ")
			for ai, a := range mr.fields {
				if ai > 0 {
					out(", ")
				}
				out("%s", g.getter(a, "c.Reply"))
			}
			out("\n")
		}

		out("}\n")
	}

	out("\n")

	out("func %s(p %s, in, out C.Message) error {\n", pfxname("dispatch", t.name), t.name)
	out("switch (in.Method()) {\n")

	for _, method := range t.fields {
		// Local method dispatch
		mt := method.typ
		mr := method.ret
		out("case %d:\n", method.ordinal)
		out("a := in.Arguments()\n")

		if len(mr.fields) > 0 {
			for ai := range mr.fields {
				if ai > 0 {
					out(", ")
				}
				out("r%d", ai)
			}
			out(" := ")
		}

		out("p.%s(", method.name)
		for ai, a := range mt.fields {
			if ai > 0 {
				out(", ")
			}
			out("%s", g.getter(a, "a"))
		}
		out(")\n")

		switch len(mr.fields) {
		case 0:
			out("return nil\n")
		default:
			out("ret, err := out.P.Segment.NewStruct(%d, %d)\n", mr.dataSize, mr.ptrSize)
			out("if err != nil { return err }\n")
			for ai, a := range mr.fields {
				out("if err := %s; err != nil { return err }\n", g.setter(a, "ret", sprintf("r%d", ai)))
			}
			out("return out.SetArguments(ret)\n")
		}
	}

	out("}\n")
	out("return C.ErrInvalidInterface\n")
	out("}\n")

	out("\n")
}

func (p *GoFile) writeEnum(t *typ) {
	out("type %s uint16\n", t.name)
	out("const (\n")

	for _, f := range t.fields {
		if len(f.comment) > 0 {
			out("/* %s */\n", f.comment)
		}
		out("%s%s %s = %d\n", t.enumPrefix, f.name, t.name, f.ordinal)
	}
	out(")\n")

	out("\n")
}

var ReadEnumFunc = strings.TrimLeft(`
func %sList(p C.Pointer, i int, def []%s) []%s {
	if m := p.ReadPtr(i); m.Type() == C.List {
		r := make([]%s, m.Size())
		for i := range r {
			r[i] = %s(m.Read16(i))
		}
		return r
	}
	return def
}
`, "\r\n")

var WriteEnumFunc = strings.TrimLeft(`
func %sList(p C.Pointer, i int, v, def []%s) (err error) {
	var m C.Pointer
	if !C.SliceEqual(v, def) {
		if m, err = p.Segment.NewList(16, 0, len(v)); err != nil {
			return err
		}
		d := p.Data()
		for i, u := range v {
			putLittle16(d[2*i:], uint16(u))
		}
	}
	return p.WritePtr(i, m)
}
`, "\r\n")

var ReadPtrFunc = strings.TrimLeft(`
func %sList(p C.Pointer, i int, def []%s) []%s {
	if m := p.ReadPtr(i); m.Type() == C.List {
		r := make([]%s, m.Size())
		for i := range r {
			r[i] = %s{m.ReadPtr(i)}
		}
		return r
	}
	return def
}
`, "\r\n")

var WritePtrFunc = strings.TrimLeft(`
func %sList(p C.Pointer, i int, v, def []%s) (err error) {
	var m C.Pointer
	if !C.SliceEqual(v, def) {
		if m, err = p.Segment.NewPointerList(len(v)); err != nil {
			return err
		}
		for i, u := range v {
			if err := u.MarshalCaptain(m, i); err != nil {
				return err
			}
		}
	}
	return p.WritePtr(i, m)
}
`, "\r\n")

func (g *GoFile) writeListFuncs(t *typ) {
	switch t.typ {
	case enumType:
		out(ReadEnumFunc, pfxname("read", t.name), t.name, t.name, t.name, t.name)
		out(WriteEnumFunc, pfxname("write", t.name), t.name)
	case structType:
		out(ReadPtrFunc, pfxname("read", t.name), t.name, t.name, t.name, t.name)
		out(WritePtrFunc, pfxname("write", t.name), t.name)
	case interfaceType:
		out(ReadPtrFunc, pfxname("read", t.name), t.name, t.name, t.name, pfxname("Remote", t.name))
		out(WritePtrFunc, pfxname("write", t.name), t.name)
	}
}

func (p *file) writeGo(name string) {
	base := []rune{}
	for _, r := range filepath.Base(name) {
		if !unicode.IsNumber(r) && !unicode.IsLetter(r) && r != '_' {
			r = '_'
		}
		base = append(base, r)
	}

	g := &GoFile{
		constants: p.constants,
		types:     p.types,
		bufName:   string(base),
		buf:       C.NewBuffer(nil),
		valbase:   string(base),
	}

	g.resolveTypes()
	f, err := os.Create(name + ".go")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	currentOutput = f
	defer f.Close()

	s := `package %s
import (
	"encoding/binary"
	C %s
	M "math"
)

var (
	putLittle16 = binary.LittleEndian.PutUint16
	_ = M.Float32bits

`

	out(s, *pkg, strconv.Quote(importPath))

	for _, c := range g.constants {
		if len(c.comment) > 0 {
			out("/* %s */\n", c.comment)
		}

		// For many types the inferred type is correct so don't output
		// the type on the left
		switch c.typ.typ {
		case boolType, float64Type, stringType, dataType, listType, structType, enumType:
			out("%s = %s\n", c.name, g.GoString(c.value, c.name))

		case int8Type, uint8Type, int16Type, uint16Type,
			int32Type, uint32Type, int64Type, uint64Type, float32Type:
			out("%s %s = %s\n", c.name, c.typ.name, g.GoString(c.value, c.name))

		case voidType:

		default:
			panic("unhandled")
		}
	}

	out(")\n")

	for _, t := range g.types {
		if len(t.comment) > 0 {
			out("/* %s */\n", t.comment)
		}

		switch t.typ {
		case enumType:
			g.writeEnum(t)
			g.writeListFuncs(t)
		case unionType:
			g.writeEnum(t)
		case structType:
			g.writeStruct(t)
			g.writeListFuncs(t)
		case interfaceType:
			g.writeInterface(t)
			g.writeListFuncs(t)
		}
	}

	out("var (\n")
	for i, v := range g.vals {
		out("%s_%d = %s\n", g.valbase, i, v)
	}

	out("%s = C.NewBuffer([]byte{", g.bufName)
	for i, b := range g.buf.Data {
		if i%8 == 0 {
			out("\n")
		}
		out(" %d,", b)
	}
	out("\n}).Root()\n")
	out(")\n")
}
