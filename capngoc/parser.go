package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf16"
	"unicode/utf8"
)

const importPath = "github.com/jmckaskill/go-capnproto"

var (
	errUnfinishedString = errors.New("unfinished string")
	errInvalidUnicode   = errors.New("invalid unicode")
)

func strerr(err error) {
	if err == io.EOF || err == io.ErrUnexpectedEOF {
		panic(errUnfinishedString)
	} else {
		panic(err)
	}
}

const eofToken rune = 'E'

type token struct {
	typ rune
	str string
}

func (t token) String() string {
	switch t.typ {
	case eofToken:
		return "EOF"
	case '@':
		return "@" + t.str
	case '"':
		return strconv.Quote(t.str)
	case '\'':
		return "'" + strings.Trim(strconv.Quote(t.str), "\"") + "'"
	case '0', 'a':
		return t.str
	case '#':
		return "# " + t.str
	}

	return string(t.typ)
}

type file struct {
	in        *bufio.Reader
	finished  bool
	line      int
	types     []*typ
	constants []*field
}

// next returns the next token. The returned rune indicates which type of
// token is returned. Valid tokens are:
// @ ordinal
// 0 number constant
// :;=,{}() standalone symbol
// " string constant - returned string has all string escapes pre-processed
// ' character constant - returned string has all string escapes pre-processed
// # comment - has whitespace stripped
// a standalane word, symbol, atom, etc
// eofToken end of file
func (p *file) rawnext() token {
	var ch rune
	var ret []rune
	var err error

	if p.finished {
		return token{eofToken, ""}
	}

	// Strip whitespace
	for {
		ch, _, err = p.in.ReadRune()
		if err == io.EOF {
			return token{eofToken, ""}
		} else if err != nil {
			panic(err)
		}
		if ch == '\n' {
			p.line++
		}
		if !unicode.IsSpace(ch) {
			break
		}
	}

	switch ch {
	case '@':
		// ordinal
		for {
			ch, _, err = p.in.ReadRune()
			if err == io.EOF {
				p.finished = true
				break
			} else if err != nil {
				panic(err)
			}

			if unicode.IsSpace(ch) {
				p.in.UnreadRune()
				break
			}

			ret = append(ret, ch)
		}

		// if we have a zero length number, the number parse will fail
		return token{'@', string(ret)}

	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0', '.', '-', '+':
		// number
		for {
			if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '.' && ch != '-' && ch != '+' {
				p.in.UnreadRune()
				break
			}

			ret = append(ret, ch)

			ch, _, err = p.in.ReadRune()
			if err == io.EOF {
				p.finished = true
				break
			} else if err != nil {
				panic(err)
			}
		}

		return token{'0', string(ret)}

	case ':', ';', '=', ',', '{', '}', '(', ')', '[', ']':
		// single char symbol
		return token{ch, ""}

	case '"', '\'':
		// string
		finish := ch
		for {
			ch, _, err = p.in.ReadRune()
			if err != nil {
				strerr(err)
			}

			switch ch {
			case '\\':
				ch, _, err = p.in.ReadRune()
				if err != nil {
					strerr(err)
				}

				// List of string espaces should match latest C/C++ standards
				switch ch {
				case 'a':
					ret = append(ret, '\a')
				case 'b':
					ret = append(ret, '\b')
				case 'f':
					ret = append(ret, '\f')
				case 'n':
					ret = append(ret, '\n')
				case 'r':
					ret = append(ret, '\r')
				case 't':
					ret = append(ret, '\t')
				case 'v':
					ret = append(ret, '\v')
				case '\'':
					ret = append(ret, '\'')
				case '"':
					ret = append(ret, '"')
				case '\\':
					ret = append(ret, '\\')
				case '?':
					ret = append(ret, '?')
				case 'x':
					// hex
					var h [2]byte
					_, err := io.ReadFull(p.in, h[:])
					if err != nil {
						strerr(err)
					}

					val, err := strconv.ParseUint(string(h[:]), 16, 8)
					if err != nil {
						panic(fmt.Errorf("error parsing hex escape: %v", err))
					}

					ret = append(ret, rune(val))

				case 'u':
					// unicode
					var h [4]byte
					_, err := io.ReadFull(p.in, h[:])
					if err != nil {
						strerr(err)
					}

					val, err := strconv.ParseUint(string(h[:]), 16, 16)
					if err != nil {
						strerr(fmt.Errorf("error parsing unicode espace: %v", err))
					}

					if utf16.IsSurrogate(rune(val)) {
						var h2 [6]byte
						_, err := io.ReadFull(p.in, h2[:])
						if err != nil {
							strerr(err)
						}
						if h2[0] != '\\' || h2[1] != 'u' {
							strerr(errInvalidUnicode)
						}
						val2, err := strconv.ParseUint(string(h2[:]), 16, 16)
						if err != nil {
							strerr(fmt.Errorf("error parsing unicode espace: %v", err))
						}
						ret = append(ret, utf16.DecodeRune(rune(val), rune(val2)))
					} else {
						ret = append(ret, rune(val))
					}

				case 'U':
					// long unicode
					var h [8]byte
					_, err := io.ReadFull(p.in, h[:])
					if err != nil {
						strerr(err)
					}

					val, err := strconv.ParseUint(string(h[:]), 16, 21)
					if err != nil {
						panic(fmt.Errorf("error parsing unicode escape: %v", err))
					}

					ret = append(ret, rune(val))

				case '0', '1', '2', '3', '4', '5', '6', '7':
					// octal
					oct := []rune{ch}
					for len(oct) < 3 {
						ch, _, err = p.in.ReadRune()
						if err != nil {
							strerr(err)
						} else if ch < '0' || ch > '7' {
							p.in.UnreadRune()
							break
						}
						oct = append(oct, ch)
					}
					val, err := strconv.ParseUint(string(oct), 8, 8)
					if err != nil {
						panic(fmt.Errorf("error parsing octal escape: %v", err))
					}
					ret = append(ret, rune(val))

				default:
					panic(fmt.Errorf("unexpected string espace \\%c", ch))
				}
			case finish:
				return token{finish, string(ret)}
			case '\n':
				ret = append(ret, ch)
				p.line++
			default:
				ret = append(ret, ch)
			}
		}

	case '#':
		// comment to the end of the line
		for {
			ch, _, err = p.in.ReadRune()
			if err == io.EOF {
				p.finished = true
				break
			} else if err != nil {
				panic(err)
			} else if ch == '\n' {
				p.line++
				break
			}

			ret = append(ret, ch)
		}

		return token{'#', strings.TrimSpace(string(ret))}
	}

	if unicode.IsLetter(ch) || ch == '_' {
		// word, symbolic name, etc
		ret = append(ret, ch)
		for {
			ch, _, err = p.in.ReadRune()
			if err == io.EOF {
				p.finished = true
				break
			} else if err != nil {
				panic(err)
			} else if !unicode.IsLetter(ch) && !unicode.IsDigit(ch) && ch != '.' && ch != '_' {
				p.in.UnreadRune()
				break
			}

			ret = append(ret, ch)
		}

		return token{'a', string(ret)}
	}

	panic(fmt.Errorf("unexpected char %c (%x)", ch, ch))
}

func (p *file) next() token {
	for {
		tok := p.rawnext()
		if tok.typ != '#' {
			return tok
		}
	}
}

type typeType int

const (
	structType typeType = iota
	enumType
	interfaceType
	voidType
	boolType
	int8Type
	int16Type
	int32Type
	int64Type
	uint8Type
	uint16Type
	uint32Type
	uint64Type
	float32Type
	float64Type
	stringType
	listType
	bitsetType
)

func (t *typ) String() string {
	return t.name
}

type typ struct {
	typ      typeType
	name     string
	comment  string
	fields   []*field
	dataSize int
	ptrSize  int
	listType *typ
}

type field struct {
	name    string
	comment string
	ordinal int
	typestr string
	typ     *typ
	args    *typ
	value   *value
	offset  int
}

type value struct {
	typ    rune
	name   string
	string string
	fields []*value
}

func (p *file) parseValue(tok token) *value {
	switch tok.typ {
	case '(':
		// struct
		v := &value{typ: '('}
		for {
			tok := p.next()
			if tok.typ == ')' {
				break
			} else if tok.typ != 'a' {
				panic(fmt.Errorf("expected struct field name got %s", tok))
			}
			p.expect('=', "=")

			value := p.parseValue(p.next())
			value.name = tok.str
			v.fields = append(v.fields, value)

			tok = p.next()
			if tok.typ == ')' {
				break
			} else if tok.typ != ',' {
				panic(fmt.Errorf("expected , got %s", tok))
			}
		}
		fmt.Printf("have value %#v\n", v)
		return v

	case '[':
		// array
		v := &value{typ: '['}
		for {
			tok := p.next()
			if tok.typ == ']' {
				break
			}

			value := p.parseValue(tok)
			v.fields = append(v.fields, value)

			tok = p.next()
			if tok.typ == ']' {
				break
			} else if tok.typ != ',' {
				panic(fmt.Errorf("expected , got %s", tok))
			}
		}
		fmt.Printf("have value %#v\n", v)
		return v

	case '0', '\'', '"', 'a':
		v := &value{
			typ:    tok.typ,
			string: tok.str,
		}
		fmt.Printf("have value %#v\n", v)
		return v
	}

	panic(fmt.Errorf("unexpected token parsing value %s", tok))
}

func (p *file) parseTypeName() string {
	tok := p.next()
	if tok.typ != 'a' {
		panic(fmt.Errorf("expected type name got %s", tok))
	}

	switch tok.str {
	case "Void":
		return "struct{}"
	case "Bool":
		return "bool"
	case "Data":
		return "[]uint8"
	case "Text":
		return "string"
	case "Int8":
		return "int8"
	case "UInt8":
		return "uint8"
	case "Int16":
		return "int16"
	case "UInt16":
		return "uint16"
	case "Int32":
		return "int32"
	case "UInt32":
		return "uint32"
	case "Int64":
		return "int64"
	case "UInt64":
		return "uint64"
	case "Float32":
		return "float32"
	case "Float64":
		return "float64"
	case "List":
		tok = p.next()
		if tok.typ != '(' {
			panic(fmt.Errorf("malformed list type - expected ( got %s", tok))
		}

		inner := p.parseTypeName()

		tok = p.next()
		if tok.typ != ')' {
			panic(fmt.Errorf("malformed list type - expected ) got %s", tok))
		}

		if inner == "bool" {
			return "capnproto.Bitset"
		} else {
			return "[]" + inner
		}

	default:
		return strings.Replace(tok.str, ".", "·", -1)
	}

}

func (p *file) parseOrdinal() int {
	tok := p.next()
	if tok.typ != '@' {
		panic(fmt.Errorf("expected ordinal got %s", tok))
	}

	ord, err := strconv.ParseUint(tok.str, 10, 16)
	if err != nil {
		panic(fmt.Errorf("error parsing %s: %v", tok, err))
	}

	return int(ord)
}

func (p *file) parseComment() (string, token) {
	firstParagraphEnd := p.line + 1
	comment := ""

	tok := p.rawnext()

	for tok.typ == '#' {
		if p.line <= firstParagraphEnd {
			firstParagraphEnd++
			comment += tok.str + "\n"
		}

		tok = p.next()
	}

	return comment, tok
}

func (p *file) expect(typ rune, name string) token {
	tok := p.next()
	if tok.typ != typ {
		panic(fmt.Errorf("expected %s got %s", name, tok))
	}
	return tok
}

func (p *file) parseConst(ns string) token {
	tok := p.expect('a', "const name")

	field := &field{
		name: ns + tok.str,
	}

	p.expect(':', "type seperator colon :")
	field.typestr = p.parseTypeName()

	p.expect('=', "constant value =")
	field.value = p.parseValue(p.next())

	p.expect(';', "constant terminator ;")
	field.comment, tok = p.parseComment()

	p.constants = append(p.constants, field)
	return tok
}

func (p *file) parseEnum(ns string) {
	tok := p.expect('a', "enum name")

	t := &typ{
		typ:  enumType,
		name: ns + tok.str,
	}

	p.expect('{', "opening brace {")
	t.comment, tok = p.parseComment()

	for tok.typ != '}' {
		if tok.typ != 'a' {
			panic(fmt.Errorf("expected enum value name got", tok))
		}

		field := &field{
			name: tok.str,
		}

		field.ordinal = p.parseOrdinal()
		field.comment, tok = p.parseComment()

		t.fields = append(t.fields, field)
	}

	fmt.Printf("have enum %#v\n", t)
	p.types = append(p.types, t)
}

func (p *file) parseInterface(ns string) {
	tok := p.expect('a', "interface name")

	t := &typ{
		typ:  interfaceType,
		name: ns + tok.str,
	}

	p.expect('{', "opening brace {")
	t.comment, tok = p.parseComment()

	for tok.typ != '}' {
		if tok.typ != 'a' {
			panic(fmt.Errorf("expected method name got %s", tok))
		}

		f := &field{
			name: tok.str,
		}

		f.ordinal = p.parseOrdinal()
		p.expect('(', "arguments opening brace (")

		f.args = &typ{typ: structType}

		tok = p.next()
		for tok.typ != ')' {
			arg := &field{}

			// Name is optional ie method @0(:bool) :bool is valid
			if tok.typ == 'a' {
				arg.name = tok.str
				tok = p.next()
			}

			if tok.typ != ':' {
				panic(fmt.Errorf("expected type colon : got %s", tok))
			}

			arg.typestr = p.parseTypeName()

			// Can give a default value
			// method @0 (:bool = true) :bool
			tok = p.next()
			if tok.typ == '=' {
				arg.value = p.parseValue(p.next())
				tok = p.next()
			}

			if tok.typ == ')' {
				break
			} else if tok.typ != ',' {
				panic(fmt.Errorf("expected comma or ) got %s", tok))
			}

			arg.comment, tok = p.parseComment()

			fmt.Printf("have arg %#v\n", arg)
			f.args.fields = append(f.args.fields, arg)
		}

		tok = p.next();
		if tok.typ == ':' {
			f.typestr = p.parseTypeName()
			p.expect(';', "method terminator")
		} else if tok.typ != ';' {
			panic(fmt.Errorf("expected : or ; got %s", tok))
		}

		f.comment, tok = p.parseComment()

		fmt.Printf("have field %#v\n", f)
		t.fields = append(t.fields, f)
	}

	fmt.Printf("have interface %#v\n", t)
	p.types = append(p.types, t)
}

func (p *file) parseStruct(ns string) {
	tok := p.expect('a', "struct name")

	t := &typ{
		typ:  structType,
		name: ns + tok.str,
	}

	ns = t.name + "·"

	p.expect('{', "opening brace {")
	t.comment, tok = p.parseComment()

	for tok.typ != '}' {
		if tok.typ != 'a' {
			panic(fmt.Errorf("expected field got %s", tok))
		}

		switch tok.str {
		case "interface":
			p.parseInterface(ns)
			tok = p.next()
		case "struct":
			p.parseStruct(ns)
			tok = p.next()
		case "const":
			tok = p.parseConst(ns)
		case "enum":
			p.parseEnum(ns)
			tok = p.next()
		default:
			f := &field{
				name: tok.str,
			}

			f.ordinal = p.parseOrdinal()

			p.expect(':', "type seperator :")
			f.typestr = p.parseTypeName()

			tok = p.next()
			if tok.typ == '=' {
				f.value = p.parseValue(p.next())
				tok = p.next()
			} else if tok.typ != ';' {
				panic(fmt.Errorf("expected field terminator ; got %s", tok))
			}

			f.comment, tok = p.parseComment()

			fmt.Printf("have field %#v\n", f)
			t.fields = append(t.fields, f)
		}
	}

	fmt.Printf("have type %#v\n", t)
	p.types = append(p.types, t)
}

func (p *file) parse() (err error) {
	defer func() {
		if r, ok := recover().(error); ok {
			err = r
		}
	}()

	tok := p.next()

	for tok.typ != eofToken {
		switch tok.typ {
		case 'a':
			switch tok.str {
			case "struct":
				p.parseStruct("")
				tok = p.next()
			case "interface":
				p.parseInterface("")
				tok = p.next()
			case "enum":
				p.parseEnum("")
				tok = p.next()
			case "const":
				tok = p.parseConst("")
			default:
				panic(fmt.Errorf("unexpected token %s", tok))
			}
		default:
			panic(fmt.Errorf("unexpected token %s", tok))
		}
	}

	return nil
}

func (p *file) addBuiltinTypes() {
	p.types = append(p.types, &typ{typ: voidType, name: "struct{}"})
	p.types = append(p.types, &typ{typ: boolType, name: "bool"})
	p.types = append(p.types, &typ{typ: int8Type, name: "int8"})
	p.types = append(p.types, &typ{typ: uint8Type, name: "uint8"})
	p.types = append(p.types, &typ{typ: int16Type, name: "int16"})
	p.types = append(p.types, &typ{typ: uint16Type, name: "uint16"})
	p.types = append(p.types, &typ{typ: int32Type, name: "int32"})
	p.types = append(p.types, &typ{typ: uint32Type, name: "uint32"})
	p.types = append(p.types, &typ{typ: int64Type, name: "int64"})
	p.types = append(p.types, &typ{typ: uint64Type, name: "uint64"})
	p.types = append(p.types, &typ{typ: float32Type, name: "float32"})
	p.types = append(p.types, &typ{typ: float64Type, name: "float64"})
	p.types = append(p.types, &typ{typ: stringType, name: "string"})
	p.types = append(p.types, &typ{typ: bitsetType, name: "capnproto.Bitset"})
}

func (p *file) doFindType(pfx int, name string) (*typ, error) {
	for i := len(p.types)-1; i >= 0; i-- {
		t := p.types[i]

		fmt.Printf("findtype %d %s %s\n", pfx, name, t.name)

		for j := 0; j <= pfx; j += 2 {
			if name[j:] == t.name {
				// create the list types
				for j > 0 {
					j -= 2
					t = &typ{
						typ:      listType,
						name:     name[j:],
						listType: t,
					}
					p.types = append(p.types, t)
				}

				return t, nil
			}
		}
	}

	return nil, fmt.Errorf("can't find type %s", name)
}

func (p *file) findType(ns *typ, name string) (*typ, error) {
	// If the user specifies a fully qualified type then we use that.
	// Otherwise we look for both the type in the local namespace and the
	// root namespace.
	pfx := 0
	for name[pfx:pfx+2] == "[]" {
		pfx += 2
	}
	if ns != nil && strings.Index(name, "·") < 0 {
		t, err := p.doFindType(pfx, name[:pfx] + ns.name + "·" + name[pfx:])
		if err == nil {
			return t, nil
		}
	}

	return p.doFindType(pfx, name)
}

func (p *file) resolveTypes() error {
	var err error

	for _, c := range p.constants {
		c.typ, err = p.findType(nil, c.typestr)
		if err != nil {
			return err
		}
	}

	for _, t := range p.types {
		for _, f := range t.fields {
			fmt.Printf("resolve field %s\n", f.name)
			if f.typestr != "" {
				f.typ, err = p.findType(t, f.typestr)
				if err != nil {
					return err
				}
			}

			if f.args != nil {
				for _, arg := range f.args.fields {
					arg.typ, err = p.findType(t, arg.typestr)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	// We cap the max list depth at one for the go types
	for _, t := range p.types {
		if t.typ == listType && t.listType.typ == listType {
			t.name = "[]capnproto.Pointer"
		}
	}

	return nil
}

type ordinalFields []*field

func (o ordinalFields) Len() int      { return len(o) }
func (o ordinalFields) Swap(i, j int) { o[i], o[j] = o[j], o[i] }

func (o ordinalFields) Less(i, j int) bool {
	return o[i].ordinal < o[j].ordinal
}

/* This alignment thing is from ORBit2
 * Align a value upward to a boundary, expressed as a number of bytes.
 * E.g. align to an 8-byte boundary with argument of 8.
 */

/*
 *   (this + boundary - 1)
 *          &
 *    ~(boundary - 1)
 */

func align(val, align int) int {
	return (val + align - 1) &^ (align - 1)
}

func (p *file) resolveOffsets() error {
	for _, t := range p.types {
		if t.typ != structType {
			continue
		}

		fields := ordinalFields(t.fields)
		sort.Sort(fields)

		for i, f := range fields {
			if f.ordinal != i {
				return fmt.Errorf("missing ordinal %d in type %s", i, t.name)
			}

			switch f.typ.typ {
			case voidType:
				f.offset = t.dataSize

			case boolType:
				f.offset = t.dataSize
				t.dataSize++

			case int8Type, uint8Type:
				f.offset = align(t.dataSize, 8)
				t.dataSize = f.offset + 8

			case int16Type, uint16Type, enumType:
				f.offset = align(t.dataSize, 16)
				t.dataSize = f.offset + 16

			case int32Type, uint32Type, float32Type:
				f.offset = align(t.dataSize, 32)
				t.dataSize = f.offset + 32

			case int64Type, uint64Type, float64Type:
				f.offset = align(t.dataSize, 64)
				t.dataSize = f.offset + 64

			case stringType, structType, interfaceType, listType, bitsetType:
				f.offset = t.ptrSize
				t.ptrSize++

			default:
				panic("unhandled")
			}
		}

		t.dataSize = align(t.dataSize, 64) / 64
	}

	return nil
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

func (t *typ) findField(name string) *field {
	for _, f := range t.fields {
		if f.name == name {
			return f
		}
	}
	panic(fmt.Errorf("can't find field %s in type %s", name, t.name))
}

var pkg = flag.String("pkg", "main", "Package name to use with generated files")

type printer func(string, ...interface{})

func (v *value) write(t *typ, out printer, marshalled bool) {
	switch t.typ {
	case bitsetType:
		if v.typ != '[' {
			panic(fmt.Errorf("unexpected value %v in bitset", v))
		}
		set := make([]byte, (len(v.fields)+7)/8)
		for i, v := range v.fields {
			switch v.string {
			case "true":
				set[i/8] |= 1 << uint(i%8)
			case "false":
			default:
				panic(fmt.Errorf("unexpected value %v in bitset", v))
			}
		}
		if marshalled {
			out("NewBitset(capnproto.Memory, ")
		}
		out("capnproto.Bitset{")
		for i, b := range set {
			if i > 0 {
				out(", ")
			}
			out("%d", b)
		}
		out("}")
		if marshalled {
			out(")")
		}

	case voidType:
		out("struct{}")

	case enumType:
		if marshalled {
			out("uint16(%s)", v.string)
		} else {
			out("%s", v.string)
		}

	case stringType:
		if v.typ != '"' {
			panic(fmt.Errorf("unexpected value %v in string", v))
		}

		if marshalled {
			out("NewString(capnproto.NewMemory, %s)", strconv.Quote(v.string))
		} else {
			out("%s", strconv.Quote(v.string))
		}

	case listType:
		if v.typ == '"' && t.listType.typ == uint8Type {
			// Data fields with a string value
			if marshalled {
				out("NewUInt8List(capnproto.NewMemory, %s)", strconv.Quote(v.string))
			} else {
				out("[]byte(%s)", strconv.Quote(v.string))
			}
			return
		}

		if v.typ != '[' {
			panic(fmt.Errorf("unexpected value %v in list", v))
		}

		listTypeName := t.name
		innerMarshalled := false

		if marshalled {
			switch t.listType.typ {
			case int8Type:
				out("NewInt8List(capnproto.NewMemory, ")
			case uint8Type:
				out("NewUInt8List(capnproto.NewMemory, ")
			case int16Type:
				out("NewInt16List(capnproto.NewMemory, ")
			case uint16Type:
				out("NewUInt16List(capnproto.NewMemory, ")
			case enumType:
				out("NewUInt16List(capnproto.NewMemory, ")
				listTypeName = "[]uint16"
				innerMarshalled = true
			case int32Type:
				out("NewInt32List(capnproto.NewMemory, ")
			case uint32Type:
				out("NewUInt32List(capnproto.NewMemory, ")
			case int64Type:
				out("NewInt64List(capnproto.NewMemory, ")
			case uint64Type:
				out("NewUInt64List(capnproto.NewMemory, ")
			case float32Type:
				out("NewFloat32List(capnproto.NewMemory, ")
			case float64Type:
				out("NewFloat64List(capnproto.NewMemory, ")
			case voidType:
				out("NewVoidList(capnproto.NewMemory, ")
			case stringType:
				out("NewStringList(capnproto.NewMemory, ")
			case bitsetType:
				out("NewBitsetList(capnproto.NewMemory, ")
			case listType:
				out("NewPointerList(capnproto.NewMemory, ")
				listTypeName = "[]capnproto.Pointer"
				innerMarshalled = true
			default:
				panic("unhandled")
			}

		} else if t.listType.typ == listType{
			listTypeName = "[]capnproto.Pointer"
			innerMarshalled = true
		}

		out("%s{", listTypeName)
		for i, v := range v.fields {
			if i > 0 {
				out(", ")
			}
			v.write(t.listType, out, innerMarshalled)
		}
		out("}")

		if marshalled {
			out(")")
		}

	case structType:
		if v.typ != '(' {
			panic(fmt.Errorf("unexpected value %v in struct", v))
		}
		out("func() %s {\n", t.name)
		out("ptr := %s(capnproto.Memory)\n", pfxname("new", t.name))
		for _, v := range v.fields {
			f := t.findField(v.name)
			if f.typ.typ != voidType {
				out("ptr.%s(", pfxname("set", f.name))
				v.write(f.typ, out, false)
				out(")\n")
			}
		}
		out("return ptr\n")
		out("}()")
		if marshalled {
			out(".Ptr")
		}

	default:
		// number types, can be number or constants
		switch v.typ {
		case '0', 'a':
			out("%s", v.string)
		default:
			panic(fmt.Errorf("unexpected value %v in number for %v", v, t))
		}
	}
}

func (f *field) writeGetterXorDefault(out printer) {
	if f.value != nil {
		out("ret ^= %s\n", f.value.string)
	}
}

func (f *field) writeGetter(out printer, ptr string) {
	switch f.typ.typ {
	case voidType:
		out("ret = struct{}\n")
	case boolType:
		def := 0
		if f.value != nil && f.value.string == "true" {
			def = 1
		}
		out("ret = (capnproto.ReadUInt8(%s, %d) & %d) != %d\n", f.offset/8, f.offset%8, def)
	case int8Type:
		out("ret = int8(capnproto.ReadUInt8(%s, %d))\n", ptr, f.offset/8)
		f.writeGetterXorDefault(out)
	case uint8Type:
		out("ret = capnproto.ReadUInt8(%s, %d)\n", ptr, f.offset/8)
		f.writeGetterXorDefault(out)
	case int16Type:
		out("ret = int16(capnproto.ReadUInt16(%s, %d))\n", ptr, f.offset/8)
		f.writeGetterXorDefault(out)
	case uint16Type:
		out("ret = capnproto.ReadUInt16(%s, %d)\n", ptr, f.offset/8)
		f.writeGetterXorDefault(out)
	case int32Type:
		out("ret = capnproto.ReadInt32(%s, %d)\n", ptr, f.offset/8)
		f.writeGetterXorDefault(out)
	case uint32Type:
		out("ret = capnproto.ReadUInt32(%s, %d)\n", ptr, f.offset/8)
		f.writeGetterXorDefault(out)
	case int64Type:
		out("ret = capnproto.ReadInt64(%s, %d)\n", ptr, f.offset/8)
		f.writeGetterXorDefault(out)
	case uint64Type:
		out("ret = capnproto.ReadUInt64(%s, %d)\n", ptr, f.offset/8)
		f.writeGetterXorDefault(out)
	case enumType:
		out("ret = %s(capnproto.ReadUInt16(%s, %d))\n", f.typ.name, ptr, f.offset/8)
		f.writeGetterXorDefault(out)

	case float32Type:
		out("u := capnproto.ReadUint32(%s, %d)\n", ptr, f.offset/8)
		if f.value != nil {
			out("u ^= math.Float32bits(%s)\n", f.value.string)
		}
		out("ret = math.Float32frombits(u)\n")

	case float64Type:
		out("u := capnproto.ReadUint64(%s, %d)\n", ptr, f.offset/8)
		if f.value != nil {
			out("u ^= math.Float64bits(%s)\n", f.value.string)
		}
		out("ret = math.Float64frombits(u)\n")

	case stringType:
		def := ""
		if f.value != nil {
			def = f.value.string
		}
		out("ret = capnproto.ToString(%s, %s)\n", ptr, strconv.Quote(def))

	case structType:
		out("ret.Ptr = %s\n", ptr)
		if f.value != nil {
			out("if ret.Ptr == nil {\n")
			out("ret = ")
			f.value.write(f.typ, out, false)
			out("\n}\n")
		}

	case interfaceType:
		out("ret = remote·%s{Ptr: %s}\n", f.typ.name, ptr)

	case bitsetType:
		out("ret = capnproto.ToBitset(%s)\n", ptr)
		if f.value != nil {
			out("if ret == nil {\n")
			out("ret = ")
			f.value.write(f.typ, out, false)
			out("\n}\n")
		}

	case listType:
		switch f.typ.listType.typ {
		case voidType:
			out("ret = capnproto.ToVoidList(%s)\n", ptr)
		case int8Type:
			out("ret = capnproto.ToInt8List(%s)\n", ptr)
		case uint8Type:
			out("ret = capnproto.ToUInt8List(%s)\n", ptr)
		case int16Type:
			out("ret = capnproto.ToInt16List(%s)\n", ptr)
		case uint16Type:
			out("ret = capnproto.ToUInt16List(%s)\n", ptr)
		case int32Type:
			out("ret = capnproto.ToInt32List(%s)\n", ptr)
		case uint32Type:
			out("ret = capnproto.ToUInt32List(%s)\n", ptr)
		case int64Type:
			out("ret = capnproto.ToInt64List(%s)\n", ptr)
		case uint64Type:
			out("ret = capnproto.ToUInt64List(%s)\n", ptr)
		case enumType:
			out(`
			u16 := capnproto.ToUInt16List(%s)
			// the binary layout matches, but the type does not
			pret := (*reflect.SliceHeader)(unsafe.Pointer(&ret))
			pret = *(*reflect.SliceHeader)(unsafe.Pointer(&u16))
			`, ptr)
		case float32Type:
			out("ret = capnproto.ToFloat32List(%s)\n", ptr)
		case float64Type:
			out("ret = capnproto.ToFloat64List(%s)\n", ptr)
		case stringType:
			out("ret = capnproto.ToStringList(%s)\n", ptr)
		case bitsetType:
			out("ret = capnproto.ToBitsetList(%s)\n", ptr)
		case structType:
			out(`
			"list := capnproto.ToPointerList(%s)
			"// the binary layout matches, but the type does not
			"pret := (*reflect.SliceHeader)(unsafe.Pointer(&ret))
			"*pret = *(*reflect.SliceHeader)(unsafe.Pointer(&list))
			`, ptr)
		case interfaceType:
			out(`
			ptrs := capnproto.ToPointerList(%s)
			ret := make(%s, len(ptrs))
			for i := range ptrs {
				ret[i] = remote·%s{Ptr: ptr}
			}
			`, ptr, f.typ.name, f.typ.listType.name)
		case listType:
			out("ret = ToPointerList(%s)\n", ptr)
		default:
			panic("unhandled")
		}

		if f.value != nil {
			out("if ret == nil {\n")
			out("ret = ")
			f.value.write(f.typ, out, false)
			out("\n}\n")
		}

	default:
		panic("unhandled")
	}

}

func (f *field) writeDataSetter(out printer, ptr string) {
	switch f.typ.typ {
	case float32Type:
		out("u := math.Float32bits(v)\n")
	case float64Type:
		out("u := math.Float64bits(v)\n")
	}

	if f.value != nil {
		switch f.typ.typ {
		case boolType, enumType, int8Type, uint8Type, int16Type, uint16Type, int32Type, uint32Type, int64Type, uint64Type:
			out("v ^= %s\n", f.value.string)
		case float32Type:
			out("u ^= math.Float32frombits(%s)\n", f.value.string)
		case float64Type:
			out("u ^= math.Float64frombits(%s)\n", f.value.string)
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
		out("err = capnproto.WriteBool(%s, %d, v)\n", ptr, f.offset)
	case int8Type:
		out("err = capnproto.WriteUInt8(%s, %d, uint8(v))\n", ptr, f.offset/8)
	case uint8Type:
		out("err = capnproto.WriteUInt8(%s, %d, v)\n", ptr, f.offset/8)
	case int16Type:
		out("err = capnproto.WriteUInt16(%s, %d, uint16(v))\n", ptr, f.offset/8)
	case uint16Type:
		out("err = capnproto.WriteUInt16(%s, %d, v)\n", ptr, f.offset/8)
	case int32Type:
		out("err = capnproto.WriteInt32(%s, %d, uint32(v))\n", ptr, f.offset/8)
	case uint32Type:
		out("err = capnproto.WriteUInt32(%s, %d, v)\n", ptr, f.offset/8)
	case int64Type:
		out("err = capnproto.WriteInt64(%s, %d, uint64(v))\n", ptr, f.offset/8)
	case uint64Type:
		out("err = capnproto.WriteUInt64(%s, %d, v)\n", ptr, f.offset/8)
	case enumType:
		out("err = capnproto.WriteUInt16(%s, %d, uint16(v))\n", ptr, f.offset/8)
	case float32Type:
		out("err = capnproto.WriteUInt32(%s, %d, u)\n", ptr, f.offset/8)
	case float64Type:
		out("err = capnproto.WriteUInt64(%s, %d, u)\n", ptr, f.offset/8)
	default:
		panic("unhandled")
	}
}

func (f *field) writeNewPointer(out printer, new string) {
	switch f.typ.typ {
	case stringType:
		out(`
		data, err = capnproto.NewString(%s, v)
		if err != nil { return }
		`, new)
	case structType:
		out("data = v.Ptr\n")
	case interfaceType:
		out(`
		data, err = v.MarshalCaptain(%s)
		if err != nil { return }
		`, new)
	case bitsetType:
		out(`
		data, err = capnproto.NewBitset(%s, v)
		if err != nil { return }
		`, new)
	case listType:
		switch f.typ.listType.typ {
		case voidType:
			out(`
			data, err = capnproto.NewVoidList(%s, v)
			if err != nil { return }
			`, new)
		case int8Type:
			out(`
			data, err = capnproto.NewInt8List(%s, v)
			if err != nil { return }
			`, new)
		case uint8Type:
			out(`
			data, err = capnproto.NewUInt8List(%s, v)
			if err != nil { return }
			`, new)
		case int16Type:
			out(`
			data, err = capnproto.NewInt16List(%s, v)
			if err != nil { return }
			`, new)
		case uint16Type:
			out(`
			data, err = capnproto.NewUInt16List(%s, v)
			if err != nil { return }
			`, new)
		case int32Type:
			out(`
			data, err = capnproto.NewInt32List(%s, v)
			if err != nil { return }
			`, new)
		case uint32Type:
			out(`
			data, err = capnproto.NewUInt32List(%s, v)
			if err != nil { return }
			`, new)
		case int64Type:
			out(`
			data, err = capnproto.NewInt64List(%s, v)
			if err != nil { return }
			`, new)
		case uint64Type:
			out(`
			data, err = capnproto.NewUInt64List(%s, v)
			if err != nil { return }
			`, new)
		case float32Type:
			out(`
			data, err = capnproto.NewFloat32List(%s, v)
			if err != nil { return }
			`, new)
		case float64Type:
			out(`
			data, err = capnproto.NewFloat64List(%s, v)
			if err != nil { return }
			`, new)
		case enumType:
			out(`
			// the binary layout matches, but the type does not
			var u16 []uint16
			pu16 := (*reflect.SliceHeader)(unsafe.Pointer(&u16))
			*pu16 = *(*reflect.SliceHeader)(unsafe.Pointer(&v))
			data, err = capnproto.NewUInt16List(%s, u16)
			if err != nil { return }
			`, new)
		case stringType:
			out(`
			data, err = capnproto.NewStringList(%s, v)
			if err != nil { return }
			`, new)
		case bitsetType:
			out(`
			data, err = capnproto.NewBitsetList(%s, v)
			if err != nil { return }
			`, new)
		case structType:
			out(`
			// the binary layout matches, but the type does not
			var ptrs []Pointer
			pptrs := (*reflect.SliceHeader)(unsafe.Pointer(&ptrs))
			*pptrs = *(*reflect.SliceHeader)(unsafe.Pointer(&v))
			data, err = capnproto.NewPointerList(%s, ptrs)
			if err != nil { return }
			`, new)
		case interfaceType:
			out(`
			cookies := make([]Pointer, len(v))
			for i, iface := range v {
				cookies[i], err = iface.MarshalCaptain(%s)
				if err != nil { return }
			}
			data, err = capnproto.NewPointerList(%s, cookies)
			if err != nil { return }
			`, new)
		case listType:
			out(`
			data, err = capnproto.NewPointerList(%s, v)
			if err != nil { return }
			`, new)
		default:
			panic("unhandled")
		}

	default:
		panic("unhandled")
	}
}

func (p *file) writeStruct(t *typ, out printer) {
	out(`
	type %s struct {
		Ptr capnproto.Pointer
	}`, t.name)

	out(`
	func %s(new capnproto.NewFunc) (Pointer, error)) %s {
		return %s{Ptr: new(capnproto.MakeStruct(%d, %d))}
	}
	`, pfxname("new", t.name), t.name, t.name, t.dataSize, t.ptrSize)

	for _, f := range t.fields {
		if len(f.comment) > 0 {
			out("/* %s */\n", f.comment)
		}

		out("func (p %s) %s() (ret %s) {\n", t.name, f.name, f.typ.name)

		switch f.typ.typ {
		case stringType, structType, interfaceType, bitsetType, listType:
			out("ptr := capnproto.ReadPtr(p.Ptr, %d)\n", f.offset)
			f.writeGetter(out, "ptr")
		default:
			f.writeGetter(out, "p.Ptr")
		}

		out("return ret\n}\n\n")

		out("func (p %s) %s(v %s) (err error) {\n", t.name, pfxname("set", f.name), f.typ.name)

		switch f.typ.typ {
		case stringType, structType, interfaceType, bitsetType, listType:
			out("var data captain.Pointer\n")
			f.writeNewPointer(out, "p.Ptr.New")
			out("return p.Ptr.WritePtrs(%d, []Pointer{data})\n", f.offset)
		default:
			f.writeDataSetter(out, "p.Ptr")
			out("return\n")
		}

		out("}\n\n")
	}
}

func (p *file) write(out printer) {
	out(`
	package %s
	import (
		"math"
		"reflect"
		"unsafe"
		%s
	)
	`, *pkg, strconv.Quote(importPath))

	for _, c := range p.constants {
		if len(c.comment) > 0 {
			out("/* %s */\n", c.comment)
		}

		out("var %s %s = ", c.name, c.typ.name)
		c.value.write(c.typ, out, false)
		out("\n")

	}

	out("\n")

	for _, t := range p.types {
		if len(t.comment) > 0 {
			out("/* %s */\n", t.comment)
		}

		switch t.typ {
		case enumType:
			out("type %s uint16\n", t.name)
			out("const (\n")
			for _, f := range t.fields {
				if len(f.comment) > 0 {
					out("/* %s */\n", f.comment)
				}
				out("%s %s = %d\n", f.name, t.name, f.ordinal)
			}
			out(")\n\n")

		case structType:
			p.writeStruct(t, out)

		case interfaceType:
			out("type %s interface{\n", t.name)

			for _, f := range t.fields {
				if len(f.comment) > 0 {
					out("/* %s */\n", f.comment)
				}

				out("%s(", f.name)
				for ai, a := range f.args.fields {
					if ai > 0 {
						out(", ")
					}
					out("%s %s", a.name, a.typ.name)
				}
				out(") (")

				if f.typ != nil {
					out("%s, ", f.typ.name)
				}

				out("error)\n")
			}

			out("}\n")

			out(`type remote·%s struct {
				Ptr capnproto.Pointer
			}
			`, t.name)

			for _, f := range t.fields {
				f.args.name = fmt.Sprintf("args·%s·%s", t.name, f.name)
				p.writeStruct(f.args, out)

				if f.typ != nil {
					out("func getret·%s·%s(p capnproto.Pointer) %s {\n", t.name, f.name, f.typ.name)
					f.writeGetter(out, "p")
					out("}\n\n")

					out("func setret·%s·%s(new capnproto.NewFunc, v %s) (data capnproto.Pointer, err error) {\n", t.name, f.name, f.typ.name)

					switch f.typ.typ {
					case stringType, structType, interfaceType, bitsetType, listType:
						f.writeNewPointer(out, "new")
					default:
						out(`
						data, err := new(captain.MakeStruct(8, 0))
						if err != nil { return nil, err }
						`)
						f.writeDataSetter(out, "data")
					}
					out("return\n}\n\n")
				}

				out("func (p remote·%s) %s(", t.name, f.name)
				for ai, a := range f.args.fields {
					if ai > 0 {
						out(", ")
					}
					out("a%d %s", ai, a.typ.name)
				}
				out(") (")
				if f.typ != nil {
					out("ret %s, ", f.typ.name)
				}
				out("err error) {\n")

				out("args := newargs·%s·%s(p.New)\n", t.name, f.name)
				for ai, a := range f.args.fields {
					out("args.%s(a%d)\n", pfxname("set", a.name), ai)
				}

				out("args, err = p.Ptr.Call(%d, args)\n", f.ordinal)

				if f.typ != nil {
					out(`
					if err == nil {
						ret = getret·%s·%s(args)
					}
					`, t.name, f.name)
				}

				out("return\n}\n")
			}

			out(`
			func dispatch·%s(iface interface{}, method int, args capnproto.Pointer, retnew capnproto.NewFunc) (capnproto.Pointer, error) {
				p, ok := iface.(%s)
				if !ok {
					return capnproto.ErrInvalidInterface
				}
				switch (method) {
			`, t.name, t.name)

			for _, f := range t.fields {
				out(`
				case %d:
					a := args·%s·%s{Ptr: args}
				`, f.ordinal, t.name, f.name)

				if f.typ != nil {
					out("r, err := p.%s(")
				} else {
					out("p.%s(")
				}
				for ai, a := range f.args.fields {
					if ai > 0 {
						out(", ")
					}
					out("a.%s()", a.name)
				}
				out(")\n")
				out("if err != nil { return nil, err }\n")

				if f.typ != nil {
					out("return setret·%s·%s(retnew, r)\n")
				} else {
					out("return nil, nil\n")
				}
			}
		}
	}
}

func main() {
	flag.Parse()

	for _, name := range flag.Args() {
		f, err := os.Open(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
			continue
		}

		p := &file{
			in:   bufio.NewReader(f),
			line: 1,
		}

		err = p.parse()
		f.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s:%d: %v\n", name, p.line, err)
			continue
		}

		p.addBuiltinTypes()

		if err := p.resolveTypes(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		if err := p.resolveOffsets(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		name += ".go"
		out, err := os.Create(name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
			continue
		}

		writer := bufio.NewWriter(out)
		p.write(func(format string, args ...interface{}) {
			fmt.Fprintf(writer, format, args...)
		})
		writer.Flush()
		out.Close()
	}
}
