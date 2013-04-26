package main

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func pfxname(pfx, name string) string {
	r, rn := utf8.DecodeRuneInString(name)
	if unicode.IsUpper(r) {
		p, pn := utf8.DecodeRuneInString(pfx)
		return string(unicode.ToUpper(p)) + pfx[pn:] + name
	} else {
		return pfx + string(unicode.ToUpper(r)) + name[rn:]
	}
}

func (p *file) resolveGoTypes() {
	for _, t := range p.types {
		switch t.typ {
		case structType, interfaceType:
			t.name = strings.Replace(t.name, "·", "_", -1)
		case enumType, unionType:
			t.name = strings.Replace(t.name, "·", "_", -1)
			t.enumPrefix = t.name + "_"
		case methodType:
			t.name = "args_" + strings.Replace(t.name, "·", "_", -1)
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
		case listType:
			// We cap the max list depth at one for the go types
			if t.listType.typ == listType {
				t.name = "[]C.Pointer"
			} else {
				t.name = "[]" + t.listType.name
			}
		case bitsetType:
			t.name = "C.Bitset"
		default:
			panic("unhandled")
		}
	}
}

func (v *value) GoString(p *file) string {
	t := v.typ
	switch t.typ {
	case bitsetType:
		return sprintf("C.ToBitset(C.Must(%s.ReadRoot(%d)))", p.bufName, v.Marshal(p))

	case stringType:
		return strconv.Quote(v.string)

	case listType:
		// Data fields with a string value
		if v.tok == '"' && t.listType.typ == uint8Type {
			return "[]byte( " + strconv.Quote(v.string) + ")"
		}

		switch t.listType.typ {
		case voidType:
			return sprintf("make([]struct{}, %d)", len(v.array))

		case listType:
			out := "[]C.Pointer{"
			for i, v := range v.array {
				if i > 0 {
					out += ", "
				}
				out += sprintf("C.Must(%s.ReadRoot(%d))", p.bufName, v.Marshal(p))
			}
			out += "}"
			return out

		default:
			out := t.name + "{"
			for i, v := range v.array {
				if i > 0 {
					out += ", "
				}
				out += v.GoString(p)
			}
			out += "}"
			return out
		}

	case structType:
		return sprintf("%s{Ptr: C.Must(%s.ReadRoot(%d))}", t.name, p.bufName, v.Marshal(p))

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

	case int8Type, uint8Type, int16Type, uint16Type,
		int32Type, uint32Type, int64Type, uint64Type:
		return sprintf("%v", v.num)

	default:
		panic("unhandled")
	}
}

func (f *field) writeGetter(p *file, ptr string) {
	xor := ""
	def := ""
	if f.value != nil {
		def = f.value.GoString(p)
		xor = " ^ " + def
	}

	switch f.typ.typ {
	case voidType:
		// nothing to do
	case boolType:
		def := 0
		if f.value != nil && f.value.bool {
			def = 1
		}
		out("return (C.ReadUInt8(%s, %d) & %d) != %d\n", ptr, f.offset/8, f.offset%8, def)
	case int8Type:
		out("return int8(C.ReadUInt8(%s, %d))%s\n", ptr, f.offset/8, xor)
	case uint8Type:
		out("return C.ReadUInt8(%s, %d)%s\n", ptr, f.offset/8, xor)
	case int16Type:
		out("return int16(C.ReadUInt16(%s, %d))%s\n", ptr, f.offset/8, xor)
	case uint16Type:
		out("return C.ReadUInt16(%s, %d)%s\n", ptr, f.offset/8, xor)
	case int32Type:
		out("return int32(C.ReadUInt32(%s, %d))%s\n", ptr, f.offset/8, xor)
	case uint32Type:
		out("return C.ReadUInt32(%s, %d)%s\n", ptr, f.offset/8, xor)
	case int64Type:
		out("return int64(C.ReadUInt64(%s, %d))%s", ptr, f.offset/8, xor)
	case uint64Type:
		out("return C.ReadUInt64(%s, %d)%s", ptr, f.offset/8, xor)
	case enumType, unionType:
		out("return %s(C.ReadUInt16(%s, %d))%s", f.typ.name, ptr, f.offset/8, xor)

	case float32Type:
		if f.value != nil {
			out(`u := C.ReadUInt32(%s, %d)
			u ^= M.Float32bits(%s)
			return M.Float32frombits(u)
			`, ptr, f.offset/8, def)
		} else {
			out("return M.Float32frombits(C.ReadUInt32(%s, %d))\n", ptr, f.offset/8)
		}

	case float64Type:
		if f.value != nil {
			out(`u := C.ReadUInt64(%s, %d)
			u ^= M.Float64bits(%s)
			return M.Float64frombits(u)
			`, ptr, f.offset/8, def)
		} else {
			out("return M.Float64frombits(C.ReadUInt64(%s, %d))\n", ptr, f.offset/8)
		}

	case stringType:
		out("return C.ToString(%s, %s)\n", ptr, def)

	case structType:
		if f.value != nil {
			out("if %s == nil { return %s }\n", ptr, def)
		}
		out("return %s{Ptr: %s}\n", f.typ.name, ptr)

	case interfaceType:
		out("return _%s_remote{Ptr: %s}\n", f.typ.name, ptr)

	case bitsetType:
		if f.value != nil {
			out(`ret := C.ToBitset(%s)
			if ret == nil {
				ret = %s
			}
			return ret
			`, ptr, def)
		} else {
			out("return C.ToBitset(%s)\n", ptr)
		}

	case listType:
		ret := "return"
		if f.value != nil {
			ret = "ret :="
		}
		switch f.typ.listType.typ {
		case voidType:
			out("%s C.ToVoidList(%s)\n", ret, ptr)
		case int8Type:
			out("%s C.ToInt8List(%s)\n", ret, ptr)
		case uint8Type:
			out("%s C.ToUInt8List(%s)\n", ret, ptr)
		case int16Type:
			out("%s C.ToInt16List(%s)\n", ret, ptr)
		case uint16Type:
			out("%s C.ToUInt16List(%s)\n", ret, ptr)
		case int32Type:
			out("%s C.ToInt32List(%s)\n", ret, ptr)
		case uint32Type:
			out("%s C.ToUInt32List(%s)\n", ret, ptr)
		case int64Type:
			out("%s C.ToInt64List(%s)\n", ret, ptr)
		case uint64Type:
			out("%s C.ToUInt64List(%s)\n", ret, ptr)
		case enumType:
			// the binary layout matches, but the type does not
			out(`u16 := C.ToUInt16List(%s)
			ret := %s(nil)
			pret := (*reflect.SliceHeader)(unsafe.Pointer(&ret))
			*pret = *(*reflect.SliceHeader)(unsafe.Pointer(&u16))
			`, ptr, f.typ.name)
		case float32Type:
			out("%s C.ToFloat32List(%s)\n", ret, ptr)
		case float64Type:
			out("%s C.ToFloat64List(%s)\n", ret, ptr)
		case stringType:
			out("%s C.ToStringList(%s)\n", ret, ptr)
		case bitsetType:
			out("%s C.ToBitsetList(%s)\n", ret, ptr)
		case structType:
			// the binary layout matches, but the type does not
			out(`list := C.ToPointerList(%s)
			ret := %s(nil)
			pret := (*reflect.SliceHeader)(unsafe.Pointer(&ret))
			*pret = *(*reflect.SliceHeader)(unsafe.Pointer(&list))
			`, ptr, f.typ.name)
		case interfaceType:
			out(`ptrs := C.ToPointerList(%s)
			ret := make(%s, len(ptrs))
			for i := range ptrs {
				ret[i] = _%s_remote{Ptr: ptr}
			}
			`, ptr, f.typ.name, f.typ.listType.name)
		case listType:
			out("%s C.ToPointerList(%s)\n", ret, ptr)
		default:
			panic("unhandled")
		}

		if f.value != nil {
			out("if ret == nil { return %s }\nreturn ret\n", def)
		} else {
			switch f.typ.listType.typ {
			case enumType, structType, interfaceType:
				out("return ret\n")
			}
		}

	default:
		panic("unhandled")
	}

}

func (f *field) writeDataSetter(p *file, ptr, ret string) {
	arg := "v"

	switch f.typ.typ {
	case float32Type:
		arg = "M.Float32bits(v)"
	case float64Type:
		arg = "M.Float64bits(v)"
	}

	if f.value != nil {
		switch f.typ.typ {
		case enumType, int8Type, uint8Type, int16Type, uint16Type, int32Type, uint32Type, int64Type, uint64Type:
			arg = sprintf("%s ^ %s", arg, f.value.GoString(p))
		case float32Type:
			arg = sprintf("%s ^ M.Float32bits(%s)", arg, f.value.GoString(p))
		case float64Type:
			arg = sprintf("%s ^ M.Float64bits(%s)", arg, f.value.GoString(p))
		case boolType:
			if f.value.bool {
				arg = sprintf("!%s", arg)
			}
		case voidType:
			// nothing to do
		default:
			panic("unhandled")
		}
	}

	switch f.typ.typ {
	case voidType:
		// nothing to do ...
	case boolType:
		out("%s C.WriteBool(%s, %d, %s)", ret, ptr, f.offset, arg)
	case int8Type:
		out("%s C.WriteUInt8(%s, %d, uint8(%s))", ret, ptr, f.offset/8, arg)
	case uint8Type:
		out("%s C.WriteUInt8(%s, %d, %s)", ret, ptr, f.offset/8, arg)
	case int16Type, enumType:
		out("%s C.WriteUInt16(%s, %d, uint16(%s))", ret, ptr, f.offset/8, arg)
	case uint16Type:
		out("%s C.WriteUInt16(%s, %d, %s)", ret, ptr, f.offset/8, arg)
	case int32Type:
		out("%s C.WriteUInt32(%s, %d, uint32(%s))", ret, ptr, f.offset/8, arg)
	case uint32Type, float32Type:
		out("%s C.WriteUInt32(%s, %d, %s)", ret, ptr, f.offset/8, arg)
	case int64Type:
		out("%s C.WriteUInt64(%s, %d, uint64(%s))", ret, ptr, f.offset/8, arg)
	case uint64Type, float64Type:
		out("%s C.WriteUInt64(%s, %d, %s)", ret, ptr, f.offset/8, arg)
	default:
		panic("unhandled")
	}
}

func (f *field) writeNewPointer(new, ret string) {
	switch f.typ.typ {
	case stringType:
		out("data, err := C.NewString(%s, v)\n", new)
	case interfaceType:
		out("data, err := v.MarshalCaptain(%s)\n", new)
	case bitsetType:
		out("data, err := C.NewBitset(%s, v)\n", new)
	case listType:
		switch f.typ.listType.typ {
		case voidType:
			out("data, err := C.NewList(%s, C.VoidList, len(v))\n", new)
		case int8Type:
			out("data, err := C.NewInt8List(%s, v)\n", new)
		case uint8Type:
			out("data, err := C.NewUInt8List(%s, v)\n", new)
		case int16Type:
			out("data, err := C.NewInt16List(%s, v)\n", new)
		case uint16Type:
			out("data, err := C.NewUInt16List(%s, v)\n", new)
		case int32Type:
			out("data, err := C.NewInt32List(%s, v)\n", new)
		case uint32Type:
			out("data, err := C.NewUInt32List(%s, v)\n", new)
		case int64Type:
			out("data, err := C.NewInt64List(%s, v)\n", new)
		case uint64Type:
			out("data, err := C.NewUInt64List(%s, v)\n", new)
		case float32Type:
			out("data, err := C.NewFloat32List(%s, v)\n", new)
		case float64Type:
			out("data, err := C.NewFloat64List(%s, v)\n", new)
		case enumType:
			// the binary layout matches, but the type does not
			out(`u16 := []uint16(nil)
			pu16 := (*reflect.SliceHeader)(unsafe.Pointer(&u16))
			*pu16 = *(*reflect.SliceHeader)(unsafe.Pointer(&v))
			data, err := C.NewUInt16List(%s, u16)
			`, new)
		case stringType:
			out("data, err := C.NewStringList(%s, v)\n", new)
		case bitsetType:
			out("data, err := C.NewBitsetList(%s, v)\n", new)
		case structType:
			// the binary layout matches, but the type does not
			out(`ptrs := []C.Pointer(nil)
			pptrs := (*reflect.SliceHeader)(unsafe.Pointer(&ptrs))
			*pptrs = *(*reflect.SliceHeader)(unsafe.Pointer(&v))
			data, err := C.NewList(%s, C.PointerList, len(ptrs))
			if err != nil { %s }
			err = data.WritePtrs(0, ptrs)
			`, new)
		case interfaceType:
			out(`data, err := C.NewList(%s, C.PointerList)
			if err != nil { %s }
			for i, iface := range v {
				cookie, err := iface.MarshalCaptain(data.New)
				if err != nil { %s }
				err = data.WritePtrs(i, []C.Pointer{cookie})
				if err != nil { %s }
			}
			`, new, ret, ret, ret)
		case listType:
			out(`data, err := C.NewList(%s, C.PointerList)
			if err != nil { %s }
			err = data.WritePtrs(0, v)
			`, new, ret)
		default:
			panic("unhandled")
		}

	default:
		panic("unhandled")
	}
}

func (p *file) writeStruct(t *typ) {
	out(`type %s struct {
		Ptr C.Pointer
	}
	`, t.name)

	out(`func %s(new C.NewFunc) (%s, error) {
		ptr, err := C.NewStruct(new, %d, %d)
		return %s{Ptr: ptr}, err
	}
	`, pfxname("new", t.name), t.name, t.dataSize, t.ptrSize, t.name)

	for _, f := range t.fields {
		if f.typ.typ == voidType {
			continue
		}

		if len(f.comment) > 0 {
			out("/* %s */\n", f.comment)
		}

		out("func (p %s) %s() %s {\n", t.name, f.name, f.typ.name)

		if f.typ.isptr() {
			out("ptr := C.ReadPtr(p.Ptr, %d)\n", f.offset)
			f.writeGetter(p, "ptr")
		} else {
			f.writeGetter(p, "p.Ptr")
		}

		out("}\n")

		if f.typ.typ != unionType {
			out("func (p %s) %s(v %s) error {\n", t.name, pfxname("set", f.name), f.typ.name)

			if f.union != nil {
				out(`if err := C.WriteUInt16(p.Ptr, %d, uint16(%s%s)); err != nil {
					return err
				}
				`, f.union.offset/8, f.union.typ.enumPrefix, f.name)
			}

			if f.typ.typ == structType {
				out("return p.Ptr.WritePtrs(%d, []C.Pointer{v.Ptr})\n", f.offset)
			} else if f.typ.isptr() {
				f.writeNewPointer("p.Ptr.New", "return err")
				out(`if err != nil { return err }
				return p.Ptr.WritePtrs(%d, []C.Pointer{data})
				`, f.offset)
			} else {
				f.writeDataSetter(p, "p.Ptr", "return")
				out("\n")
			}

			out("}\n")
		}
	}
}

func (p *file) writeInterface(t *typ) {
	out(`type %s interface {
		C.Marshaller
		`, t.name)

	for _, method := range t.fields {
		if len(method.comment) > 0 {
			out("/* %s */\n", method.comment)
		}

		out("%s(", method.name)
		for ai, a := range method.typ.fields {
			if ai > 0 {
				out(", ")
			}
			out("%s %s", a.name, a.typ.name)
		}
		out(") (")

		if method.typ.ret != nil {
			out("%s, ", method.typ.ret.typ.name)
		}

		out("error)\n")
	}

	out("}\n")

	out(`type _%s_remote struct {
			Ptr C.Pointer
		}
		func (p _%s_remote) MarshalCaptain(new C.NewFunc) (C.Pointer, error) {
			return p.Ptr, nil
		}
		`, t.name, t.name)

	for _, method := range t.fields {
		ret := method.typ.ret
		if ret != nil {
			out("func getret_%s_%s(p C.Pointer) %s {\n", t.name, method.name, ret.typ.name)
			ret.writeGetter(p, "p")
			out("}\n")

			out("func setret_%s_%s(new C.NewFunc, v %s) (C.Pointer, error) {\n", t.name, method.name, ret.typ.name)

			if ret.typ.typ == structType {
				out("return v.Ptr, nil\n")

			} else if ret.typ.isptr() {
				ret.writeNewPointer("new", "return nil, err")
				out("if err != nil { return nil, err }\n")
				out("return data, nil\n")

			} else {
				out(`data, err := C.NewStruct(new, 8, 0)
					if err != nil { return nil, err }
					`)
				ret.writeDataSetter(p, "data", "if err :=")
				out(" != nil { return nil, err }\n")
				out("return data, nil\n")
			}
			out("}\n")
		}

		out("func (p _%s_remote) %s(", t.name, method.name)
		for ai, a := range method.typ.fields {
			if ai > 0 {
				out(", ")
			}
			out("a%d %s", ai, a.typ.name)
		}
		out(") (")
		if ret != nil {
			out("ret %s, ", ret.typ.name)
		}
		out("err error) {\n")

		out(`var args _%s_%s_args
			args, err = new_%s_%s_args(p.Ptr.New)
			if err != nil { return }
			`, t.name, method.name, t.name, method.name)

		for ai, a := range method.typ.fields {
			out("args.%s(a%d)\n", pfxname("set", a.name), ai)
		}

		if ret != nil {
			out(`var rets C.Pointer
				rets, err = p.Ptr.Call(%d, args.Ptr)
				if err != nil { return }
				ret = getret_%s_%s(rets)
				return
				`, method.ordinal, t.name, method.name)
		} else {
			out(`_, err = p.Ptr.Call(%d, args.Ptr)
				return
				`, method.ordinal)
		}

		out("}\n")
	}

	out(`func %s(iface interface{}, method int, args C.Pointer, retnew C.NewFunc) (C.Pointer, error) {
		p, ok := iface.(%s)
		if !ok {
			return nil, C.ErrInvalidInterface
		}
		switch (method) {
		`, pfxname("dispatch", t.name), t.name)

	for _, method := range t.fields {
		out(`case %d:
			a := _%s_%s_args{Ptr: args}
			`, method.ordinal, t.name, method.name)

		if method.typ.ret != nil {
			out("r, err := p.%s(", method.name)
		} else {
			out("err := p.%s(", method.name)
		}
		for ai, a := range method.typ.fields {
			if ai > 0 {
				out(", ")
			}
			out("a.%s()", a.name)
		}
		out(")\n")
		out("if err != nil { return nil, err }\n")

		if method.typ.ret != nil {
			out("return setret_%s_%s(retnew, r)\n", t.name, method.name)
		} else {
			out("return nil, nil\n")
		}
	}

	out(`default:
		return nil, C.ErrInvalidInterface
		}
		}
		`)
}

func (p *file) writeGo(name string) {
	p.bufName = strings.Replace(name, ".", "_", -1)

	out(`package %s
	import (
		M "math"
		"reflect"
		"unsafe"
		C %s
	)

	var (
		_ = M.Float32bits
		_ = reflect.SliceHeader{}
		_ = unsafe.Pointer(nil)

	`, *pkg, strconv.Quote(importPath))

	for _, c := range p.constants {
		if len(c.comment) > 0 {
			out("/* %s */\n", c.comment)
		}

		// For many types the inferred type is correct so don't output
		// the type on the left
		switch c.typ.typ {
		case boolType, float64Type, stringType, bitsetType, listType, structType, enumType:
			out("%s = %s\n", c.name, c.value.GoString(p))

		case int8Type, uint8Type, int16Type, uint16Type,
			int32Type, uint32Type, int64Type, uint64Type, float32Type:
			out("%s %s = %s\n", c.name, c.typ.name, c.value.GoString(p))

		case voidType:
			// nothing to do

		default:
			panic("unhandled")
		}
	}

	out(")\n")

	for _, t := range p.types {
		if len(t.comment) > 0 {
			out("/* %s */\n", t.comment)
		}

		switch t.typ {
		case enumType, unionType:
			out(`type %s uint16
			const (
			`, t.name)

			for _, f := range t.fields {
				if len(f.comment) > 0 {
					out("/* %s */\n", f.comment)
				}
				out("%s%s %s = %d\n", t.enumPrefix, f.name, t.name, f.ordinal)
			}
			out(")\n")

		case structType:
			p.writeStruct(t)
		case interfaceType:
			p.writeInterface(t)
		}
	}

	out("var %s = C.NewBuffer([]byte{", p.bufName)
	for i, b := range p.buf.Bytes() {
		if i%16 == 0 {
			out("\n")
		}
		out(" %#02x,", b)
	}
	out("\n})\n")
}
