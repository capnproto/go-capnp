package main

import (
	"fmt"
	C "github.com/jmckaskill/go-capnproto"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type CFile struct {
	buf       *C.Segment
	constants []*field
	types     []*typ
	vals      []string
	nextval   int
}

func (p *CFile) resolveTypes() {
	for _, t := range p.types {
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
		case methodType:
			t.name = strings.Replace(t.name, "·", "_", -1) + "_args"
		case returnType:
			t.name = strings.Replace(t.name, "·", "_", -1) + "_ret"
		case voidType:
			t.name = "void"
		case boolType:
			t.name = "int"
		case int8Type:
			t.name = "int8_t"
		case int16Type:
			t.name = "int16_t"
		case int32Type:
			t.name = "int32_t"
		case int64Type:
			t.name = "int64_t"
		case uint8Type:
			t.name = "uint8_t"
		case uint16Type:
			t.name = "uint16_t"
		case uint32Type:
			t.name = "uint32_t"
		case uint64Type:
			t.name = "uint64_t"
		case float32Type:
			t.name = "float"
		case float64Type:
			t.name = "double"
		case stringType:
			t.name = "struct capn_text"
		case dataType:
			t.name = "struct capn_data"
		case listType:
			switch t.listType.typ {
			case boolType:
				t.name = "struct capn_list1"
			case int8Type, uint8Type:
				t.name = "struct capn_list8"
			case int16Type, uint16Type, enumType:
				t.name = "struct capn_list16"
			case int32Type, uint32Type, float32Type:
				t.name = "struct capn_list32"
			case int64Type, uint64Type, float64Type:
				t.name = "struct capn_list64"
			default:
				t.name = "struct capn_ptr"
			}
		default:
			panic("unhandled")
		}
	}
}

func (c *CFile) declareInterface(t *typ) {
	out("\nstruct %s_vt {\n", t.name)
	out("\tvoid (*marshal)(const struct %s *, struct capn_ptr*, int /*off*/);\n", t.name)
	for _, method := range t.fields {
		mt := method.typ
		mr := method.ret

		if method.comment != "" {
			out("\t/* %s */\n", method.comment)
		}

		switch len(mr.fields) {
		case 0:
			out("\tvoid")
		case 1:
			switch rt := mr.fields[0].typ; rt.typ {
			case enumType:
				out("\tenum %s", rt.name)
			case structType:
				out("\tstruct %s_ptr", rt.name)
			case interfaceType:
				out("\tstruct %s", rt.name)
			default:
				out("\t%s", rt.name)
			}
		default:
			out("\tstruct %s", mr.name)
		}

		out(" (*%s)(struct %s*", method.name, t.name)
		for _, a := range mt.fields {
			switch a.typ.typ {
			case enumType:
				out(", enum %s %s", a.typ.name, a.name)
			case listType:
				out(", %s *%s", a.typ.name, a.name)
			case structType:
				out(", struct %s *%s", a.typ.name, a.name)
			case interfaceType:
				out(", struct %s %s", a.typ.name, a.name)
			default:
				out(", %s %s", a.typ.name, a.name)
			}
		}
		out(");\n")
	}
	out("};\n")

	out("\nstruct %s {\n", t.name)
	out("\tconst struct %s_vt *vt;\n", t.name)
	out("\tstruct capn_ptr p;\n")
	out("\tvoid *u;\n")
	out("};\n")
}

func (c *CFile) writeStructMember(tab string, f *field) {
	if f.comment != "" {
		out("%s/* %s */\n", tab, f.comment)
	}

	switch f.typ.typ {
	case structType:
		out("%sstruct %s_ptr %s;\n", tab, f.typ.name, f.name)
	case interfaceType:
		out("%sstruct %s %s;\n", tab, f.typ.name, f.name)
	case boolType:
		out("%sunsigned int %s : 1;\n", tab, f.name)
	case enumType:
		out("%senum %s %s;\n", tab, f.typ.name, f.name)
	case voidType:
	default:
		out("%s%s %s;\n", tab, f.typ.name, f.name)
	}
}

func (c *CFile) declareStruct(t *typ) {
	out("\nstruct %s {\n", t.name)

	for _, f := range t.fields {
		if f.typ.typ == voidType || f.union != nil {
			continue
		}

		if f.typ.typ == unionType {
			out("\tenum %s %s_tag;\n", f.typ.name, f.name)
			out("\tunion {\n")
			for _, a := range f.typ.fields {
				c.writeStructMember("\t\t", a)
			}
			out("\t} %s;\n", f.name)
		} else {
			c.writeStructMember("\t", f)
		}
	}

	out("};\n")
}

func (c *CFile) declareStructFuncs(t *typ) {
	out("\nstruct %s_ptr {\n", t.name)
	out("\tstruct capn_ptr p;\n")
	out("};\n")
	out("\nstruct %s_ptr new_%s(struct capn_segment*);\n", t.name, t.name)
	out("struct capn_ptr new_%s_list(struct capn_segment*, int sz);\n", t.name)
	out("int read_%s(const struct %s_ptr*, struct %s*);\n", t.name, t.name, t.name)
	out("int write_%s(struct %s_ptr*, const struct %s*);\n", t.name, t.name, t.name)
}

func (c *CFile) readMember(tab string, f *field) {
	xor := ""
	def := "nil"
	if f.value != nil {
		def = c.CString(f.value, "")
		xor = " ^ " + def
	}

	mbr := sprintf("s->%s", f.name)
	if f.union != nil {
		mbr = sprintf("s->%s.%s", f.union.name, f.name)
	}

	switch f.typ.typ {
	case voidType:
		/* nothing to do */
	case boolType:
		if f.value != nil && f.value.bool {
			out("%s%s = (capn_get8(&p->p, %d) & %d) == 0;\n", tab, mbr, f.offset/8, 1<<uint(f.offset%8))
		} else {
			out("%s%s = (capn_get8(&p->p, %d) & %d) != 0;\n", tab, mbr, f.offset/8, 1<<uint(f.offset%8))
		}
	case int8Type:
		out("%s%s = ((int8_t) capn_get8(&p->p, %d))%s;\n", tab, mbr, f.offset/8, xor)
	case int16Type:
		out("%s%s = ((int16_t) capn_get16(&p->p, %d))%s;\n", tab, mbr, f.offset/8, xor)
	case enumType:
		out("%s%s = (enum %s) capn_get16(&p->p, %d)%s;\n", tab, mbr, f.typ.name, f.offset/8, xor)
	case int32Type:
		out("%s%s = ((int32_t) capn_get32(&p->p, %d))%s;\n", tab, mbr, f.offset/8, xor)
	case int64Type:
		out("%s%s = ((int64_t) capn_get64(&p->p, %d))%s;\n", tab, mbr, f.offset/8, xor)
	case uint8Type:
		out("%s%s = capn_get8(&p->p, %d)%s;\n", tab, mbr, f.offset/8, xor)
	case uint16Type:
		out("%s%s = capn_get16(&p->p, %d)%s;\n", tab, mbr, f.offset/8, xor)
	case uint32Type:
		out("%s%s = capn_get32(&p->p, %d)%s;\n", tab, mbr, f.offset/8, xor)
	case uint64Type:
		out("%s%s = capn_get64(&p->p, %d)%s;\n", tab, mbr, f.offset/8, xor)
	case float32Type:
		if f.value != nil {
			out("%s%s = capn_get_float(&p->p, %d, %s);\n", tab, mbr, f.offset/8, def)
		} else {
			out("%s%s = capn_get_float(&p->p, %d, 0.0f);\n", tab, mbr, f.offset/8)
		}
	case float64Type:
		if f.value != nil {
			out("%s%s = capn_get_double(&p->p, %d, %s);\n", tab, mbr, f.offset/8, def)
		} else {
			out("%s%s = capn_get_double(&p->p, %d, 0.0);\n", tab, mbr, f.offset/8)
		}
	case stringType:
		out("%s%s = capn_read_text(&p->p, %d);\n", tab, mbr, f.offset)
		if f.value != nil {
			out("%sif (!%s.str) {\n%s\t%s = %s;\n%s}\n", tab, mbr, tab, mbr, def, tab)
		}
	case dataType:
		out("%s%s = capn_read_data(&p->p, %d);\n", tab, mbr, f.offset)
		if f.value != nil {
			out("%sif (!%s.data) {\n%s\t%s = %s;\n%s}\n", tab, mbr, tab, mbr, def, tab)
		}
	case structType:
		out("%s%s.p = capn_read_ptr(&p->p, %d);\n", tab, mbr, f.offset)
		if f.value != nil {
			out("%sif (%s.p.type == CAPN_NULL) {\n%s\t%s = %s;\n%s}\n", tab, mbr, tab, mbr, def, tab)
		}
	case interfaceType:
		out("%s%s.vt = &%s_remote_vt;\n", tab, mbr, f.typ.name)
		out("%s%s.p = capn_read_ptr(&p->p, %d);\n", tab, mbr, f.offset)
	case listType:
		switch f.typ.listType.typ {
		case boolType, int8Type, uint8Type, int16Type, uint16Type,
			int32Type, uint32Type, int64Type, uint64Type, float32Type, float64Type, enumType:
			out("%s%s.p = capn_read_ptr(&p->p, %d);\n", tab, mbr, f.offset)
			if f.value != nil {
				out("%sif (%s.p.type == CAPN_NULL) {\n%s\t%s = %s;\n%s}\n", tab, mbr, tab, mbr, def, tab)
			}
		default:
			out("%s%s = capn_read_ptr(&p->p, %d);\n", tab, mbr, f.offset)
			if f.value != nil {
				out("%sif (%s.type == CAPN_NULL) {\n%s\t%s = %s;\n%s}\n", tab, mbr, tab, mbr, def, tab)
			}
		}
	default:
		panic("unhandled")
	}
}

func (c *CFile) writeMember(tab string, f *field) {
	xor := ""
	def := "nil"
	if f.value != nil {
		def = c.CString(f.value, "")
		xor = " ^ " + def
	}

	mbr := sprintf("s->%s", f.name)
	if f.union != nil {
		mbr = sprintf("s->%s.%s", f.union.name, f.name)
	}

	switch f.typ.typ {
	case voidType:
		/* nothing to do */
	case boolType:
		mask := 1 << uint(f.offset%8)
		if f.value != nil && f.value.bool {
			out("%scapn_set8(&p->p, %d, (capn_get8(&p->p, %d) & ~%d) & ~(%s << %d));\n", tab, f.offset/8, f.offset/8, mask, mbr, f.offset%8)
		} else {
			out("%scapn_set8(&p->p, %d, (capn_get8(&p->p, %d) & ~%d) | (%s << %d));\n", tab, f.offset/8, f.offset/8, mask, mbr, f.offset%8)
		}
	case int8Type:
		out("%scapn_set8(&p->p, %d, (uint8_t) (%s%s));\n", tab, f.offset/8, mbr, xor)
	case int16Type:
		out("%scapn_set16(&p->p, %d, (uint16_t) (%s%s));\n", tab, f.offset/8, mbr, xor)
	case enumType:
		out("%scapn_set16(&p->p, %d, (uint16_t) (%s%s));\n", tab, f.offset/8, mbr, xor)
	case int32Type:
		out("%scapn_set32(&p->p, %d, (uint32_t) (%s%s));\n", tab, f.offset/8, mbr, xor)
	case int64Type:
		out("%scapn_set64(&p->p, %d, (uint64_t) (%s%s));\n", tab, f.offset/8, mbr, xor)
	case uint8Type:
		out("%scapn_set8(&p->p, %d, %s%s);\n", tab, f.offset/8, mbr, xor)
	case uint16Type:
		out("%scapn_set16(&p->p, %d, %s%s);\n", tab, f.offset/8, mbr, xor)
	case uint32Type:
		out("%scapn_set32(&p->p, %d, %s%s);\n", tab, f.offset/8, mbr, xor)
	case uint64Type:
		out("%scapn_set64(&p->p, %d, %s%s);\n", tab, f.offset/8, mbr, xor)
	case float32Type:
		if f.value != nil {
			out("%scapn_set_float(&p->p, %d, %s, %s);\n", tab, f.offset/8, mbr, def)
		} else {
			out("%scapn_set_float(&p->p, %d, %s, 0.0f);\n", tab, f.offset/8, mbr)
		}
	case float64Type:
		if f.value != nil {
			out("%scapn_set_double(&p->p, %d, %s, %s);\n", tab, f.offset/8, mbr, def)
		} else {
			out("%scapn_set_double(&p->p, %d, %s, 0.0);\n", tab, f.offset/8, mbr)
		}
	case stringType:
		if f.value != nil {
			out("%sif (%s.str != %s.str || %s.size != %s.size) {\n", tab, mbr, def, mbr, def)
			out("%s\tcapn_write_text(&p->p, %d, %s);\n", tab, f.offset, mbr)
			out("%s} else {\n", tab)
			out("%s\tcapn_write_ptr(&p->p, %d, 0);\n", tab, f.offset)
			out("%s}\n", tab)
		} else {
			out("%scapn_write_text(&p->p, %d, %s);\n", tab, f.offset, mbr)
		}
	case dataType:
		if f.value != nil {
			out("%sif (%s.data != %s.data || %s.size != %s.size) {\n", tab, mbr, def, mbr, def)
			out("%s\tcapn_write_data(&p->p, %d, %s);\n", tab, f.offset, mbr)
			out("%s} else {\n", tab)
			out("%s\tcapn_write_ptr(&p->p, %d, 0);\n", tab, f.offset)
			out("%s}\n", tab)
		} else {
			out("%scapn_write_data(&p->p, %d, %s);\n", tab, f.offset, mbr)
		}
	case structType:
		if f.value != nil {
			out("%scapn_write_ptr(&p->p, %d, %s.p.data != %s.p.data ? &%s.p : 0);\n", tab, f.offset, mbr, def, mbr)
		} else {
			out("%scapn_write_ptr(&p->p, %d, &%s.p);\n", tab, f.offset, mbr)
		}
	case interfaceType:
		out("%s%s.vt->marshal(&%s, &p->p, %d);\n", tab, mbr, mbr, f.offset)
	case listType:
		switch f.typ.listType.typ {
		case boolType, int8Type, uint8Type, int16Type, uint16Type,
			int32Type, uint32Type, int64Type, uint64Type, float32Type, float64Type, enumType:
			if f.value != nil {
				out("%scapn_write_ptr(&p->p, %d, (%s.p.data != %s.p.data || %s.p.size != %s.p.size) ? &%s.p : 0);\n", tab, f.offset, mbr, def, mbr, def, mbr)
			} else {
				out("%scapn_write_ptr(&p->p, %d, &%s.p);\n", tab, f.offset, mbr)
			}
		default:
			if f.value != nil {
				out("%scapn_write_ptr(&p->p, %d, (%s.data != %s.data || %s.size != %s.size) ? &%s : 0);\n", tab, f.offset, mbr, def, mbr, def, mbr)
			} else {
				out("%scapn_write_ptr(&p->p, %d, &%s);\n", tab, f.offset, mbr)
			}
		}
	default:
		panic("unhandled")
	}
}

func (c *CFile) defineStructFuncs(t *typ) {
	out("\nstruct %s_ptr new_%s(struct capn_segment* seg) {\n", t.name, t.name)
	out("\tstruct %s_ptr ret = {capn_new_struct(seg, %d, %d)};\n", t.name, t.dataSize/8, t.ptrSize)
	out("\treturn ret;\n")
	out("}\n")

	out("\nstruct capn_ptr new_%s_list(struct capn_segment *seg, int sz) {\n", t.name)
	out("\treturn capn_new_list(seg, sz, %d, %d);\n", t.dataSize/8, t.ptrSize)
	out("}\n")

	out("\nint read_%s(const struct %s_ptr *p, struct %s *s) {\n", t.name, t.name, t.name)
	out("\tif (p->p.type != CAPN_STRUCT)\n")
	out("\t\treturn -1;\n")

	for _, f := range t.fields {
		if f.typ.typ == unionType {
			out("\ts->%s_tag = (enum %s) capn_get16(&p->p, %d);\n", f.name, f.typ.name, f.offset/8)
			out("\tswitch (s->%s_tag) {\n", f.name)
			for _, uf := range f.typ.fields {
				out("\tcase %d:\n", uf.ordinal)
				c.readMember("\t\t", uf)
				out("\t\tbreak;\n")
			}
			out("\t}\n")

		} else if f.union == nil {
			c.readMember("\t", f)
		}
	}

	out("\treturn 0;\n")
	out("}\n")

	out("\nint write_%s(struct %s_ptr *p, const struct %s *s) {\n", t.name, t.name, t.name)
	out("\tif (p->p.type != CAPN_STRUCT)\n")
	out("\t\treturn -1;\n")

	for _, f := range t.fields {
		if f.typ.typ == unionType {
			out("\tcapn_set16(&p->p, %d, (uint16_t) s->%s_tag);\n", f.offset/8, f.name)
			out("\tswitch (s->%s_tag) {\n", f.name)
			for _, uf := range f.typ.fields {
				out("\tcase %d:\n", uf.ordinal)
				c.writeMember("\t\t", uf)
				out("\t\tbreak;\n")
			}
			out("\t}\n")

		} else if f.union == nil {
			c.writeMember("\t", f)
		}
	}

	out("\treturn 0;\n")
	out("}\n")
}

func (c *CFile) declareEnum(t *typ) {
	out("\nenum %s {\n", t.name)
	for i, f := range t.fields {
		if i > 0 {
			out(",\n")
		}
		if f.comment != "" {
			out("/* %s */\n", f.comment)
		}

		out("\t%s%s = %d", t.enumPrefix, f.name, f.ordinal)
	}
	out("\n};\n")
}

var cheader = `#ifndef CAPN_%s
#define CAPN_%s
#include "../c/capn.h"
/* AUTOGENERATED - DO NOT EDIT */

`

func (c *CFile) writeHeader(name string) {
	f, err := os.Create(name + ".h")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	currentOutput = f
	defer f.Close()

	buf := []rune{}
	for _, r := range filepath.Base(name) {
		// For max compatibility restrict to ascii identifier characters
		if !('a' <= r && r <= 'z') && !('A' <= r && r <= 'Z') && r != '_' && !('0' <= r && r <= '0') {
			r = '_'
		}
		buf = append(buf, r)
	}

	hdr := strings.ToUpper(string(buf))
	out(cheader, hdr, hdr)

	for _, t := range c.types {
		switch t.typ {
		case enumType, unionType:
			out("enum %s;\n", t.name)
		case structType:
			out("struct %s;\n", t.name)
			out("struct %s_ptr;\n", t.name)
		case interfaceType:
			out("struct %s;\n", t.name)
			out("struct %s_vt;\n", t.name)
		case returnType:
			if len(t.fields) > 1 {
				out("struct %s;\n", t.name)
			}
		}
	}

	out("\n")

	for _, f := range c.constants {
		if f.comment != "" {
			out("/* %s */\n", f.comment)
		}

		switch f.typ.typ {
		case structType:
			out("extern const struct %s_ptr %s;\n", f.typ.name, f.name)
		case enumType:
			out("extern const enum %s %s;\n", f.typ.name, f.name)
		case int8Type, uint8Type, int16Type, uint16Type,
			int32Type, uint32Type, int64Type, uint64Type,
			float32Type, float64Type, boolType, stringType, dataType, listType:
			out("extern const %s %s;\n", f.typ.name, f.name)
		case voidType:
		default:
			panic("unhandled")
		}
	}

	for _, t := range c.types {
		if t.comment != "" {
			out("/* %s */\n", t.comment)
		}
		switch t.typ {
		case enumType, unionType:
			c.declareEnum(t)
		case structType:
			c.declareStructFuncs(t)
			c.declareStruct(t)
		case returnType:
			if len(t.fields) > 1 {
				c.declareStruct(t)
			}
		case interfaceType:
			c.declareInterface(t)
		}
	}

	out("#endif\n")
}

func (c *CFile) CString(v *value, symbol string) string {
	t := v.typ

	switch t.typ {
	case voidType:
		return ""
	case enumType:
		return sprintf("((enum %s) %d)", t.name, v.num)
	case boolType:
		if v.bool {
			return "1"
		} else {
			return "0"
		}
	case float32Type:
		return sprintf("%vf", v.float)
	case float64Type:
		return sprintf("%v", v.float)
	case uint8Type:
		return sprintf("UINT8_C(%v)", uint64(v.num))
	case uint16Type:
		return sprintf("UINT16_C(%v)", uint64(v.num))
	case uint32Type:
		return sprintf("UINT32_C(%v)", uint64(v.num))
	case uint64Type:
		return sprintf("UINT64_C(%d)", uint64(v.num))
	case int8Type:
		return sprintf("INT8_C(%v)", v.num)
	case int16Type:
		return sprintf("INT16_C(%v)", v.num)
	case int32Type:
		return sprintf("INT32_C(%v)", v.num)
	case int64Type:
		return sprintf("INT64_C(%d)", v.num)
	}

	if v.symbol != "" {
		return v.symbol
	}

	v.symbol = symbol

	switch t.typ {
	case stringType:
		return sprintf("{%d,%s}", len(v.string), strconv.Quote(v.string))
	case dataType:
		str := v.string
		if v.tok != '"' {
			buf := []byte{}
			for _, v := range v.array {
				buf = append(buf, byte(v.num))
			}
			str = string(buf)
		}
		return sprintf("{%d,(uint8_t*)%s}", len(str), strconv.Quote(str))
	case structType:
		off := v.Marshal(c.buf)
		return sprintf("{{CAPN_STRUCT,0,(char*)capnbuf+%d,0,%d,%d}}",
			off*8, t.dataSize/8, t.ptrSize)
	case listType:
		lt := t.listType
		off := v.Marshal(c.buf)

		switch lt.typ {
		case boolType:
			return sprintf("{{CAPN_BIT_LIST,%d,(char*)capnbuf+%d}}", len(v.array), off*8)

		case structType, voidType:
			return sprintf("{CAPN_LIST,%d,(char*)capnbuf+%d,0,%d,%d}",
				len(v.array), off*8, lt.dataSize/8, lt.ptrSize)

		case int8Type, uint8Type, int16Type, uint16Type, enumType,
			int32Type, uint32Type, float32Type, int64Type, uint64Type, float64Type:
			return sprintf("{{CAPN_LIST,%d,(char*)capnbuf+%d,0,%d}}",
				len(v.array), off*8, lt.dataSize/8)

		case listType, stringType, dataType:
			return sprintf("{CAPN_PTR_LIST,%d,(char*)capnbuf+%d}", len(v.array), off*8)
		default:
			panic("unhandled")
		}

	case interfaceType:
		return ""
	default:
		println(t.typ)
		panic("unhandled")
	}
}

func (c *CFile) defineConstant(name string, v *value) {
	static := ""
	if name == "" {
		name = sprintf("val_%d", c.nextval)
		static = "static "
		c.nextval++
	}

	switch v.typ.typ {
	case structType:
		out("%sconst struct %s_ptr %s = %s;\n", static, v.typ.name, name, c.CString(v, name))
	case stringType:
		out("%sconst struct capn_text %s = %s;\n", static, name, c.CString(v, name))
	case dataType:
		out("%sconst struct capn_data %s = %s;\n", static, name, c.CString(v, name))
	case enumType:
		out("%sconst enum %s %s = %s;\n", static, v.typ.name, name, c.CString(v, name))
	case int8Type, uint8Type, int16Type, uint16Type,
		int32Type, uint32Type, int64Type, uint64Type,
		float32Type, float64Type, boolType, listType:
		out("%sconst %s %s = %s;\n", static, v.typ.name, name, c.CString(v, name))
	case voidType:
	default:
		panic("unhandled")
	}
}

func (c *CFile) writeSource(name string) {
	f, err := os.Create(name + ".c")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	currentOutput = f
	defer f.Close()

	out("#include \"%s.h\"\n", name)

	out("\n")
	out("static const uint8_t capnbuf[];\n")

	for _, f := range c.constants {
		c.defineConstant(f.name, f.value)
	}

	for _, t := range c.types {
		for _, f := range t.fields {
			if f.value != nil && f.value.symbol == "" {
				c.defineConstant("", f.value)
			}
		}
	}

	for _, t := range c.types {
		switch t.typ {
		case interfaceType:
			out("static const struct %s_vt %s_remote_vt;\n", t.name, t.name)
		}
	}

	for _, t := range c.types {
		switch t.typ {
		case structType:
			c.defineStructFuncs(t)
		case returnType:
		case interfaceType:
		}
	}

	out("\nstatic const uint8_t capnbuf[] = {")
	for i, b := range c.buf.Data {
		if i%8 == 0 {
			out("\n")
		}
		out(" %d,", b)
	}
	out("\n};\n")
}

func (p *file) writeC(name string) {
	c := &CFile{
		constants: p.constants,
		types:     p.types,
		buf:       C.NewBuffer(nil),
	}

	c.resolveTypes()
	c.writeHeader(name)
	c.writeSource(name)
}
