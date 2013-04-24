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
	sprintf             = fmt.Sprintf
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
	case 'u':
		return t.unum
	case 'i':
		return t.inum
	case 'f':
		return t.fnum
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
	unionType
)

func (t *typ) String() string {
	return t.name
}

type typ struct {
	typ        typeType
	name       string
	comment    string
	enumPrefix string
	fields     []*field
	sortFields ordinalFields
	dataSize   int
	ptrSize    int
	listType   *typ
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
	union   *field
}

type value struct {
	typ *typ
	tok    rune
	bool	bool
	name   string
	string string
	float   float64
	num    int64
	fields []*value
	dataPtr C.Pointer
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
		return v

	case '0':
		if inum, err := strconv.ParseInt(tok.str, 0, 64); err == nil {
			return &value{typ: 'i', num: inum, float: float64(inum)}
		}

		if unum, err := strconv.ParseUint(tok.str, 0, 64); err == nil {
			return &value{typ: 'u', num: int64(unum), float: float64(unum)}
		}

		if fnum, err := strconv.ParseFloat(tok.str, 64); err == nil {
			return &value{typ: 'f', float: fnum}
		}

		panic(fmt.Errorf("can't parse %s as a number", tok))

	case '\'':
		r, sz := utf8.DecodeRuneInString(t.str)
		if r == utf8.RuneError || sz != len(t.str) {
			panic(fmt.Errorf("can't parse %s as a character", tok))
		}
		return &value{typ: 'i', inum: r}

	case '"', 'a':
		return &value{
			typ:    tok.typ,
			string: tok.str,
		}
	}

	panic(fmt.Errorf("unexpected token parsing value %s", tok))
}

func (p *file) parseTypeName() string {
	tok := p.next()
	if tok.typ != 'a' {
		panic(fmt.Errorf("expected type name got %s", tok))
	}

	switch tok.str {
	case "Data":
		return "[]uint8"
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

		switch inner {
		case "bool":
			return "List(Bool)"
		case "Void":
			return "List(Void)"
		default:
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

	if field.value.typ == 'a' {
		c, err := p.findConstant(ns, field.value.name)
		if err != nil {
			panic(err)
		}
		field.value = c.value
	}

	p.constants = append(p.constants, field)
	return tok
}

func (p *file) parseEnum(ns string) {
	tok := p.expect('a', "enum name")

	t := &typ{
		typ:        enumType,
		name:       ns + tok.str,
		enumPrefix: ns + tok.str + "_",
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

		f.args = &typ{
			typ:  structType,
			name: sprintf("_%s_%s_args", t.name, f.name),
		}

		tok = p.next()
		for tok.typ != ')' {
			arg := &field{ordinal: len(f.args.fields)}

			// Name is optional ie method @0(:bool) :bool is valid
			if tok.typ == 'a' {
				arg.name = tok.str
				tok = p.next()
			} else {
				arg.name = sprintf("arg%d", len(f.args.fields))
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

			f.args.fields = append(f.args.fields, arg)

			if tok.typ == ')' {
				break
			} else if tok.typ != ',' {
				panic(fmt.Errorf("expected comma or ) got %s", tok))
			}

			tok = p.next()
		}

		p.types = append(p.types, f.args)

		tok = p.next()
		if tok.typ == ':' {
			f.typestr = p.parseTypeName()
			p.expect(';', "method terminator")
		} else if tok.typ != ';' {
			f.typestr = "Void"
			panic(fmt.Errorf("expected : or ; got %s", tok))
		}

		f.comment, tok = p.parseComment()

		t.fields = append(t.fields, f)
	}

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
	union := (*field)(nil)

	for tok.typ != '}' {
		if tok.typ != 'a' {
			panic(fmt.Errorf("expected field got %s", tok))
		}

		switch tok.str {
		case "interface":
			if union != nil {
				panic(fmt.Errorf("unexpected interface in union"))
			}
			p.parseInterface(ns)
			tok = p.next()
		case "struct":
			if union != nil {
				panic(fmt.Errorf("unexpected struct in union"))
			}
			p.parseStruct(ns)
			tok = p.next()
		case "const":
			if union != nil {
				panic(fmt.Errorf("unexpected const in union"))
			}
			tok = p.parseConst(ns)
		case "enum":
			if union != nil {
				panic(fmt.Errorf("unexpected enum in union"))
			}
			p.parseEnum(ns)
			tok = p.next()
		default:
			f := &field{
				name:  tok.str,
				union: union,
			}

			f.ordinal = p.parseOrdinal()

			p.expect(':', "type seperator")

			f.typestr = p.parseTypeName()

			switch f.typestr {
			case "union":
				if union != nil {
					panic(fmt.Errorf("unexpected union in union"))
				}
				p.expect('{', "union open brace {")
				f.typestr = ""
				f.typ = &typ{
					typ:  unionType,
					name: ns + f.name,
				}
				p.types = append(p.types, f.typ)
				union = f

			default:
				tok = p.next()
				switch tok.typ {
				case ';':
				case '=':
					f.value = p.parseValue(p.next())
					tok = p.next()
				default:
					panic(fmt.Errorf("expected field terminator ; got %s", tok))
				}
			}

			f.comment, tok = p.parseComment()

			t.fields = append(t.fields, f)

			if union != nil && f != union {
				if f.ordinal <= union.ordinal {
					panic(fmt.Errorf("union field %s has lower ordinal than the union tag %s", f, union))
				}

				union.typ.fields = append(union.typ.fields, f)
			}
		}

		if union != nil && tok.typ == '}' {
			union = nil
			tok = p.next()
		}
	}

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
	p.types = append(p.types, &typ{typ: voidType, name: "Void"})
	p.types = append(p.types, &typ{typ: boolType, name: "Bool"})
	p.types = append(p.types, &typ{typ: int8Type, name: "Int8"})
	p.types = append(p.types, &typ{typ: uint8Type, name: "UInt8"})
	p.types = append(p.types, &typ{typ: int16Type, name: "Int16"})
	p.types = append(p.types, &typ{typ: uint16Type, name: "UInt16"})
	p.types = append(p.types, &typ{typ: int32Type, name: "Int32"})
	p.types = append(p.types, &typ{typ: uint32Type, name: "UInt32"})
	p.types = append(p.types, &typ{typ: int64Type, name: "Int64"})
	p.types = append(p.types, &typ{typ: uint64Type, name: "UInt64"})
	p.types = append(p.types, &typ{typ: float32Type, name: "Float32"})
	p.types = append(p.types, &typ{typ: float64Type, name: "Float64"})
	p.types = append(p.types, &typ{typ: stringType, name: "Text"})
	p.types = append(p.types, &typ{typ: bitsetType, name: "List(Bool)"})
}

func (p *file) findConstant(ns string, name string) (*field, error) {
	if ns != "" && strings.Index(name, "·") < 0 {
		nsname := ns + "·" + name
		for _, c := range p.constants {
			if c.name == nsname {
				return c, nil
			}
		}
	}

	for _, c := range p.constants {
		if c.name == name {
			return c, nil
		}
	}

	return nil, fmt.Errorf("can't find constant %s", name)
}

func (p *file) doFindType(pfx int, name string) (*typ, error) {
	for i := len(p.types) - 1; i >= 0; i-- {
		t := p.types[i]
		if t.typ == unionType {
			continue
		}

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
		t, err := p.doFindType(pfx, name[:pfx]+ns.name+"·"+name[pfx:])
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
		if t.typ == unionType {
			continue
		}
		for _, f := range t.fields {
			f.typ, err = p.findType(t, f.typestr)
			if err != nil {
				return err
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

	for _, t := range p.types {
		switch *lang {
		case "go":
			switch t.typ {
			case structType, interfaceType:
				t.name = strings.Replace(t.name, "·", "_")
			case enumType, unionType:
				t.name = strings.Replace(t.name, "·", "_")
				t.enumPrefix = t.name + "_"
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

		case "c":
			switch t.typ {
			case structType, interfaceType:
				t.name = "struct " + strings.Replace(t.name, "·", "_")
			case enumType, unionType:
				base := strings.Replace(t.name, "·", "_")
				t.name = "enum " + base
				t.enumPrefix = base + "_"
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
				t.name = "char"
			case listType, bitsetType:
				t.name = "struct capn_ptr"
			default:
				panic("unhandled")
			}
		default:
			panic("unhandled language")
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

func (t *typ) isptr() bool {
	switch t.typ {
	case stringType, structType, interfaceType, listType, bitsetType:
		return true
	default:
		return false
	}
}

func (t *typ) datasize() int {
	switch t.typ {
	case boolType:
		return 1
	case int8Type, uint8Type:
		return 8
	case int16Type, uint16Type, enumType, unionType:
		return 16
	case int32Type, uint32Type, float32Type:
		return 32
	case int64Type, uint64Type, float64Type:
		return 64
	default:
		panic("unhandled")
	}
}

func (p *file) resolveOffsets() {
	for _, t := range p.types {
		if t.typ != structType {
			continue
		}

		t.sortFields := make(ordinalFields, len(t.fields))
		copy(t.sortFields, t.fields)
		sort.Sort(t.sortFields)

	next_field:
		for i, f := range t.sortFields {
			if f.ordinal != i {
				return fmt.Errorf("missing ordinal %d in type %s", i, t.name)
			}
			if f.typ.typ == voidType {
				continue
			}

			if f.union != nil {
				// want the last field with the same or bigger size
				for i := len(f.union.typ.fields) - 1; i >= 0; i-- {
					g := f.union.typ.fields[i]

					if f.ordinal <= g.ordinal || f.typ.isptr() != g.typ.isptr() || g.typ.typ == voidType {
						continue
					}

					if f.typ.isptr() {
						f.offset = g.offset
						continue next_field

					} else if f.typ.datasize() <= g.typ.datasize() {
						f.offset = g.offset
						continue next_field
					}
				}
			}

			if f.typ.isptr() {
				f.offset = t.ptrSize
				t.ptrSize++
			} else {
				sz := f.typ.datasize()
				f.offset = align(t.dataSize, sz)
				t.dataSize = f.offset + sz
			}

			if f.typ.typ == unionType {
				// Sort the union fields for the offset calculation
				fields := ordinalFields(f.typ.fields)
				sort.Sort(fields)
				f.typ.fields = []*field(fields)
			}

		}

		t.dataSize = align(t.dataSize, 64) / 64
	}
}

func (p *file) resolveValues() {
	for _, c := range p.constants {
		c.value.SetType(c.typ)
	}

	for _, t := range p.types {
		for _, f := range t.fields {
			if f.value != nil {
				f.value.SetType(f.typ)
			}
		}
	}
}

func (v *value) SetType(t *typ) {
	if v.typ != nil && t != v.typ {
		goto err
	} else if v.typ != nil {
		return
	}

	v.typ = t

	switch t.typ {
	case structType:
		if v.tok != '(' {
			goto err
		}

		for _, w := range v.fields {
			f := t.findField(w.name)
			if f != nil {
				goto err
			}

			w.SetType(f.typ)
		}

	case listType:
		if v.tok != '[' {
			goto err
		}

		for _, w := range v.fields {
			w.SetType(t.listType)
		}

	case enumType:
		if v.tok != 'a' {
			goto err
		}

		f := t.findField(v.str)
		if f != nil {
			goto err
		}

		v.num = f.ordinal

	case boolType:
		if v.tok != 'a' {
			goto err
		}

		switch v.string {
		case "true":
			v.bool = true
		case "false":
			v.bool = false
		default:
			goto err
		}

	case int8Type, int16Type, int32Type, int64Type
		uint8Type, uint16Type, uint32Type, uint64Type:
		if v.tok != 'i' {
			goto err
		}

	case float32Type, float64Type:
		if v.tok != 'i' && v.tok != 'f' {
			goto err
		}

	case stringType:
		if v.tok != '"' {
			goto err
		}
	}

	return

err:
	panic(fmt.Errorf("unexpected value %v with type %v", v, t))
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
	return nil
}

var pkg = flag.String("pkg", "main", "Package name to use with generated files")
var lang = flag.String("lang", "go", "Language to generate code for (c or go)")

func (v *value) MarshalCaptain(new C.NewFunc) (C.Pointer, error) {
	if v.dataPtr != nil {
		return v.dataPtr, nil
	}

	switch v.typ.typ {
	case structType:
		p, err := C.NewStruct(new, t.dataSize, t.ptrSize)
		if err != nil {
			return nil, err
		}

		for _, w := range v.fields {
			var err error

			def := v.typ.findField(w.name).value
			if def == nil {
				def = &value{}
			}

			switch w.typ.typ {
			case structType, stringType, listType, bitsetType:
				m, err := w.MarshalCaptain(new)
				if err != nil {
					return nil, err
				}
				if err := p.WritePtrs(w.offset, []C.Pointer{m}); err != nil {
					return nil, err
				}

			case boolType:
				u := w.bool != def.bool
				err = C.WriteBool(p, w.offset, u)

			case int8Type, uint8Type:
				u := uint8(int8(w.num)) ^ uint8(int8(def.num))
				err = C.WriteUInt8(p, w.offset/8, u)

			case int16Type, uint16Type, enumType:
				u := uint16(int16(w.num)) ^ uint16(int16(def.num))
				err = C.WriteUInt16(p, w.offset/8, u)

			case int32Type, uint32Type:
				u := uint32(int32(w.num)) ^ uint32(int32(def.num))
				err = C.WriteUInt32(p, w.offset/8, u)

			case int64Type, uint64Type:
				u := uint64(int64(w.num)) ^ uint64(int64(def.num))
				err = C.WriteUInt64(p, w.offset/8, u)

			case float32Type:
				u := math.Float32bits(float32(w.float)) ^ math.Float32bits(float32(def.float))
				err = C.WriteUInt32(p, w.offset/8, u)

			case float64Type:
				u := math.Float64bits(w.float) ^ math.Float64bits(def.float)
				err = C.WriteUInt64(p, w.offset/8, u)
			}

			if err != nil {
				return nil, err
			}
		}

	case stringType:
		return C.NewString(new, v.str)

	case bitsetType:
		set := C.MakeBitset(len(v.fields))
		for i, w := range v.fields {
			if w.bool {
				set.Set(i)
			}
		}
		return C.NewBitset(new, set)

	case listType:
		lt := t.listType

		if v.typ == '"' && lt.typ == uint8Type {
			return C.Must(C.NewUInt8List(new, []byte(v.string)))
		}

		v.Expect(t, "[")

		switch lt.typ {
		case voidType:
			return C.NewList(p.buf.New, C.VoidList, len(v.fields))

		case listType, bitsetType, stringType:
			p, err := C.NewPointerList(new)
			if err != nil {
				return nil, err
			}
			for i, w := v.fields {
				m, err := w.MarshalCaptain(p.New, lt))
				if err != nil {
					return nil, err
				}
				if err := p.WritePtrs(i, []C.Pointer{m}) {
					return nil, err
				}
			}
			return p, nil

		case structType:
			p, err := C.NewCompositeList(new, len(v.fields), lt.dataSize, lt.ptrSize)
			if err != nil {
				return err
			}
			for i, w := v.fields {
				m, err := w.MarshalCaptain(p.New, lt))
				if err != nil {
					return nil, err
				}
				if err := p.WritePtrs(i, []C.Pointer{m}) {
					return nil, err
				}
			}
			return p, nil

		case int8Type, uint8Type:
			d := make([]uint8, len(v.fields))
			for i, w := v.fields {
				d[i] = uint8(int8(w.inum))
			}
			return C.NewUInt8List(new, d)

		case int16Type, uint16Type, enumType:
			d := make([]uint16, len(v.fields))
			for i, w := v.fields {
				d[i] = uint16(int16(w.inum)))
			}
			return C.NewUInt16List(new, d)

		case int32Type, uint32Type:
			d := make([]uint32, len(v.fields))
			for i, w := v.fields {
				d[i] = uint32(int32(w.inum)))
			}
			return C.NewUInt32List(new, d)

		case int64Type, uint64Type:
			d := make([]uint64, len(v.fields))
			for i, w := v.fields {
				d[i] = uint64(int64(w.inum)))
			}
			return C.NewUInt64List(new, d)

		case float32List:
			d := make([]float32, len(v.fields))
			for i, w := v.fields {
				d[i] = float32(v.float)
			}
			return C.NewFloat32List(new, d)

		case float64List:
			d := make([]float64, len(v.fields))
			for i, w := v.fields {
				d[i] = v.float
			}
			return C.NewFloat32List(new, d)
		}
	}

	return nil, fmt.Errorf("unexpected value %v in %v", v, t)
}

func (v *value) Marshal(p *file) int {
	if v.dataPtr == nil {
		v.dataPtr = C.Must(v.MarshalCaptain(p.buf.New))
	}
	return v.dataPtr.Type().SegmentOffset()
}

func (v *value) GoString(p *file) string {
	switch t.typ {
	case bitsetType:
		return sprintf("C.ToBitset(C.Must(%s.ReadRoot(%d)))", p.bufName, v.Marshal(p))

	case stringType:
		return strconv.Quote(v.string)

	case listType:
		// Data fields with a string value
		if v.typ == '"' && t.listType.typ == uint8Type {
			return "[]byte( " + strconv.Quote(v.string) + ")"
		}

		switch t.listType {
		case voidType:
			return sprintf("make([]struct{}, %d)", len(v.fields))

		case listType:
			out := "[]C.Pointer{"
			for i, v := range v.fields {
				if i > 0 {
					out += ", "
				}
				out += sprintf("C.Must(%s.ReadRoot(%d))", p.bufName, v.Marshal(p))
			}
			out += "}"
			return out

		default:
			out := t.name + "{"
			for i, v := range v.fields {
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
			out("data, err := C.NewVoidList(%s, v)\n", new)
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
			data, err := C.NewPointerList(%s, ptrs)
			`, new)
		case interfaceType:
			out(`cookies := make([]C.Pointer, len(v))
			for i, iface := range v {
				var err error
				cookies[i], err = iface.MarshalCaptain(%s)
				if err != nil { %s }
			}
			data, err := C.NewPointerList(%s, cookies)
			`, new, ret, new)
		case listType:
			out("data, err := C.NewPointerList(%s, v)\n", new)
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
			f.writeGetter(out, "ptr")
		} else {
			f.writeGetter(out, "p.Ptr")
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
				f.writeNewPointer(out, "p.Ptr.New", "return err")
				out(`if err != nil { return err }
				return p.Ptr.WritePtrs(%d, []C.Pointer{data})
				`, f.offset)
			} else {
				f.writeDataSetter(out, "p.Ptr", "return")
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

	out(`type _%s_remote struct {
			Ptr C.Pointer
		}
		func (p _%s_remote) MarshalCaptain(new C.NewFunc) (C.Pointer, error) {
			return p.Ptr, nil
		}
		`, t.name, t.name)

	for _, f := range t.fields {
		if f.typ != nil {
			out("func getret_%s_%s(p C.Pointer) %s {\n", t.name, f.name, f.typ.name)
			f.writeGetter(out, "p")
			out("}\n")

			out("func setret_%s_%s(new C.NewFunc, v %s) (C.Pointer, error) {\n", t.name, f.name, f.typ.name)

			if f.typ.typ == structType {
				out("return v.Ptr, nil\n")
			} else if f.typ.isptr() {
				f.writeNewPointer(out, "new", "return nil, err")
				out("if err != nil { return nil, err }\n")
				out("return data, nil\n")
			} else {
				out(`data, err := C.NewStruct(new, 8, 0)
					if err != nil { return nil, err }
					`)
				f.writeDataSetter(out, "data", "if err :=")
				out(" != nil { return nil, err }\n")
				out("return data, nil\n")
			}
			out("}\n")
		}

		out("func (p _%s_remote) %s(", t.name, f.name)
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

		out(`var args _%s_%s_args
			args, err = new_%s_%s_args(p.Ptr.New)
			if err != nil { return }
			`, t.name, f.name, t.name, f.name)
		for ai, a := range f.args.fields {
			out("args.%s(a%d)\n", pfxname("set", a.name), ai)
		}

		if f.typ != nil {
			out(`var rets C.Pointer
				rets, err = p.Ptr.Call(%d, args.Ptr)
				if err != nil { return }
				ret = getret_%s_%s(rets)
				return
				`, f.ordinal, t.name, f.name)
		} else {
			out(`_, err = p.Ptr.Call(%d, args.Ptr)
				return
				`, f.ordinal)
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

	for _, f := range t.fields {
		out(`case %d:
					a := _%s_%s_args{Ptr: args}
					`, f.ordinal, t.name, f.name)

		if f.typ != nil {
			out("r, err := p.%s(", f.name)
		} else {
			out("err := p.%s(", f.name)
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
			out("return setret_%s_%s(retnew, r)\n", t.name, f.name)
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

func (p *file) writeGo() {
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
		if c.typ.typ == voidType {
			continue
		}
		if len(c.comment) > 0 {
			out("/* %s */\n", c.comment)
		}

		// For many types the inferred type is correct so don't output
		// the type on the left
		switch c.typ.typ {
		case boolType, float64Type, stringType, bitsetType, listType, structType, enumType:
			out("%s = %s\n", c.name, c.value.String(c.typ, false))

		case int8Type, uint8Type, int16Type, uint16Type,
			int32Type, uint32Type, int64Type, uint64Type, float32Type:
			out("%s %s = %s\n", c.name, c.typ.name, c.value.String(c.typ, false))

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
			p.writeStruct(t, out)
		case interfaceType:
			p.writeInterface(t, out)
		}
	}
}

var currentOutput io.Writer
func out(f string, args ...interface{}) {
	fmt.Fprintf(currentOutput, f, args...)
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
		currentOutput = writer
		p.write()
		writer.Flush()
		out.Close()
	}
}
