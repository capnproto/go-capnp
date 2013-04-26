package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	C "github.com/jmckaskill/go-capnproto"
	"io"
	"math"
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
	buf       *C.Buffer
	bufName   string
}

type typeType int

const (
	structType typeType = iota
	enumType
	interfaceType
	methodType
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
	ret        *field
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
	value   *value
	offset  int
	union   *field
}

type value struct {
	typ     *typ
	tok     rune
	bool    bool
	string  string
	float   float64
	num     int64
	array   []*value
	members map[string]*value
	dataPtr C.Pointer
}

func (v *value) String() string {
	switch v.tok {
	case 'a':
		return sprintf("%s(%s)", v.typ, v.string)
	case '"':
		return sprintf("%s(%s)", v.typ, strconv.Quote(v.string))
	case '[':
		return sprintf("%s(%s)", v.typ, v.array)
	case '(':
		return sprintf("%s(%s)", v.typ, v.members)
	case 'b':
		return sprintf("%s(%s)", v.typ, v.bool)
	case 'v':
		return sprintf("%s(void)", v.typ)
	case 'f':
		return sprintf("%s(%s)", v.typ, v.float)
	case 'i':
		return sprintf("%s(%s)", v.typ, v.num)
	default:
		panic("unhandled")
	}
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
		v := &value{tok: '(', members: make(map[string]*value)}
		for {
			tok := p.next()
			if tok.typ == ')' {
				break
			} else if tok.typ != 'a' {
				panic(fmt.Errorf("expected struct field name got %s", tok))
			}
			p.expect('=', "=")

			v.members[tok.str] = p.parseValue(p.next())

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
		v := &value{tok: '['}
		for {
			tok := p.next()
			if tok.typ == ']' {
				break
			}

			value := p.parseValue(tok)
			v.array = append(v.array, value)

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
			return &value{tok: 'i', num: inum, float: float64(inum)}
		}

		if unum, err := strconv.ParseUint(tok.str, 0, 64); err == nil {
			return &value{tok: 'i', num: int64(unum), float: float64(unum)}
		}

		if fnum, err := strconv.ParseFloat(tok.str, 64); err == nil {
			return &value{tok: 'f', float: fnum}
		}

		panic(fmt.Errorf("can't parse %s as a number", tok))

	case '\'':
		r, sz := utf8.DecodeRuneInString(tok.str)
		if r == utf8.RuneError || sz != len(tok.str) {
			panic(fmt.Errorf("can't parse %s as a character", tok))
		}
		return &value{tok: 'i', num: int64(r)}

	case '"', 'a':
		return &value{
			tok:    tok.typ,
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
		return "[]UInt8"
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

		f.typ = &typ{
			typ:     methodType,
			name:    t.name + "·" + f.name,
			comment: t.comment,
		}

		tok = p.next()
		for tok.typ != ')' {
			arg := &field{ordinal: len(f.typ.fields)}

			// Name is optional ie method @0(:bool) :bool is valid
			if tok.typ == 'a' {
				arg.name = tok.str
				tok = p.next()
			} else {
				arg.name = sprintf("arg%d", len(f.typ.fields))
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

			f.typ.fields = append(f.typ.fields, arg)

			if tok.typ == ')' {
				break
			} else if tok.typ != ',' {
				panic(fmt.Errorf("expected comma or ) got %s", tok))
			}

			tok = p.next()
		}

		p.types = append(p.types, f.typ)

		tok = p.next()
		if tok.typ == ':' {
			if str := p.parseTypeName(); str != "Void" {
				f.typ.ret = &field{typestr: str}
			}
			p.expect(';', "method terminator")
		} else if tok.typ != ';' {
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

func (p *file) addBuiltins() {
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
	p.constants = append(p.constants, &field{typestr: "Void", name: "void", value: &value{tok: 'v'}})
	p.constants = append(p.constants, &field{typestr: "Bool", name: "true", value: &value{tok: 'b', bool: true}})
	p.constants = append(p.constants, &field{typestr: "Bool", name: "false", value: &value{tok: 'b', bool: false}})
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
		switch t.typ {
		case structType, methodType:
			for _, f := range t.fields {
				if f.typ == nil {
					f.typ, err = p.findType(t, f.typestr)
					if err != nil {
						return err
					}
				}
			}
			if t.ret != nil {
				t.ret.typ, err = p.findType(t, t.ret.typestr)
				if err != nil {
					return err
				}
			}
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

		t.sortFields = make(ordinalFields, len(t.fields))
		copy(t.sortFields, t.fields)
		sort.Sort(t.sortFields)

	next_field:
		for i, f := range t.sortFields {
			if f.ordinal != i {
				panic(fmt.Errorf("missing ordinal %d in type %s", i, t.name))
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

		t.dataSize = align(t.dataSize, 64) / 8
	}
}

func (p *file) resolveValues() {
	for _, c := range p.constants {
		c.value = c.value.SetType(p, c.typ)
	}

	for _, t := range p.types {
		for _, f := range t.fields {
			if f.value != nil {
				f.value = f.value.SetType(p, f.typ)
			}
		}
	}
}

func (v *value) SetType(p *file, t *typ) *value {
	if v.typ != nil && t != v.typ {
		goto err
	} else if v.typ != nil {
		return v
	}

	v.typ = t

	if v.tok == 'a' {
		c, err := p.findConstant(t.name, v.string)
		if err == nil {
			return c.value.SetType(p, t)
		}
	}

	switch t.typ {
	case structType:
		if v.tok != '(' {
			goto err
		}

		for name, w := range v.members {
			f := t.findField(name)
			if f == nil {
				panic(fmt.Errorf("type %s does not have a field %s", t.name, name))
			}

			v.members[name] = w.SetType(p, f.typ)
		}

	case listType:
		switch v.tok {
		case '[':
			for i, w := range v.array {
				v.array[i] = w.SetType(p, t.listType)
			}
		case '"':
			if t.listType.typ != uint8Type {
				goto err
			}
		default:
			goto err
		}

	case enumType:
		if v.tok != 'a' {
			goto err
		}

		f := t.findField(v.string)
		if f == nil {
			goto err
		}

		v.num = int64(f.ordinal)

	case boolType:
		if v.tok != 'b' {
			goto err
		}

	case int8Type, int16Type, int32Type, int64Type,
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

	case voidType:
		if v.tok != 'v' {
			goto err
		}

	default:
		panic("unhandled")
	}

	return v

err:
	panic(fmt.Errorf("unexpected value %v with type %v", v, t))
}

func (t *typ) findField(name string) *field {
	for _, f := range t.fields {
		if f.name == name {
			return f
		}
	}
	return nil
}

func (v *value) MarshalCaptain(new C.NewFunc) (C.Pointer, error) {
	if v.dataPtr != nil {
		return v.dataPtr, nil
	}

	t := v.typ
	switch t.typ {
	case structType:
		p, err := C.NewStruct(new, t.dataSize, t.ptrSize)
		if err != nil {
			return nil, err
		}

		for name, w := range v.members {
			var err error

			ft := t.findField(name)
			def := ft.value
			if def == nil {
				def = &value{}
			}
			println("write", ft.offset, t.dataSize, t.ptrSize)

			switch w.typ.typ {
			case structType, stringType, listType, bitsetType:
				m, err := w.MarshalCaptain(p.New)
				if err != nil {
					return nil, err
				}
				if err := p.WritePtrs(ft.offset, []C.Pointer{m}); err != nil {
					return nil, err
				}

			case boolType:
				u := w.bool != def.bool
				err = C.WriteBool(p, ft.offset, u)

			case int8Type, uint8Type:
				u := uint8(int8(w.num)) ^ uint8(int8(def.num))
				err = C.WriteUInt8(p, ft.offset/8, u)

			case int16Type, uint16Type, enumType:
				u := uint16(int16(w.num)) ^ uint16(int16(def.num))
				err = C.WriteUInt16(p, ft.offset/8, u)

			case int32Type, uint32Type:
				u := uint32(int32(w.num)) ^ uint32(int32(def.num))
				err = C.WriteUInt32(p, ft.offset/8, u)

			case int64Type, uint64Type:
				u := uint64(int64(w.num)) ^ uint64(int64(def.num))
				err = C.WriteUInt64(p, ft.offset/8, u)

			case float32Type:
				u := math.Float32bits(float32(w.float)) ^ math.Float32bits(float32(def.float))
				err = C.WriteUInt32(p, ft.offset/8, u)

			case float64Type:
				u := math.Float64bits(w.float) ^ math.Float64bits(def.float)
				err = C.WriteUInt64(p, ft.offset/8, u)
			}

			if err != nil {
				return nil, err
			}
		}

	case stringType:
		return C.NewString(new, v.string)

	case bitsetType:
		set := C.MakeBitset(len(v.array))
		for i, w := range v.array {
			if w.bool {
				set.Set(i)
			}
		}
		return C.NewBitset(new, set)

	case listType:
		lt := t.listType

		if v.tok == '"' && lt.typ == uint8Type {
			return C.NewUInt8List(new, []byte(v.string))
		}

		switch lt.typ {
		case voidType:
			return C.NewList(new, C.VoidList, len(v.array))

		case listType, bitsetType, stringType:
			p, err := C.NewList(new, C.PointerList, len(v.array))
			if err != nil {
				return nil, err
			}
			for i, w := range v.array {
				m, err := w.MarshalCaptain(p.New)
				if err != nil {
					return nil, err
				}
				if err := p.WritePtrs(i, []C.Pointer{m}); err != nil {
					return nil, err
				}
			}
			return p, nil

		case structType:
			p, err := C.NewCompositeList(new, len(v.array), lt.dataSize, lt.ptrSize)
			if err != nil {
				return nil, err
			}
			for i, w := range v.array {
				m, err := w.MarshalCaptain(p.New)
				if err != nil {
					return nil, err
				}
				if err := p.WritePtrs(i, []C.Pointer{m}); err != nil {
					return nil, err
				}
			}
			return p, nil

		case int8Type, uint8Type:
			d := make([]uint8, len(v.array))
			for i, w := range v.array {
				d[i] = uint8(int8(w.num))
			}
			return C.NewUInt8List(new, d)

		case int16Type, uint16Type, enumType:
			d := make([]uint16, len(v.array))
			for i, w := range v.array {
				d[i] = uint16(int16(w.num))
			}
			return C.NewUInt16List(new, d)

		case int32Type, uint32Type:
			d := make([]uint32, len(v.array))
			for i, w := range v.array {
				d[i] = uint32(int32(w.num))
			}
			return C.NewUInt32List(new, d)

		case int64Type, uint64Type:
			d := make([]uint64, len(v.array))
			for i, w := range v.array {
				d[i] = uint64(int64(w.num))
			}
			return C.NewUInt64List(new, d)

		case float32Type:
			d := make([]float32, len(v.array))
			for i, w := range v.array {
				d[i] = float32(w.float)
			}
			return C.NewFloat32List(new, d)

		case float64Type:
			d := make([]float64, len(v.array))
			for i, w := range v.array {
				d[i] = w.float
			}
			return C.NewFloat64List(new, d)

		default:
			panic("unhandled")
		}

	default:
		panic("unhandled")
	}

	return nil, fmt.Errorf("unexpected value %v in %v", v, t)
}

func (v *value) Marshal(p *file) int {
	if v.dataPtr == nil {
		v.dataPtr = C.Must(v.MarshalCaptain(p.buf.NewRoot))
	}
	return v.dataPtr.Type().SegmentOffset()
}

var currentOutput *bufio.Writer

func out(f string, args ...interface{}) {
	fmt.Fprintf(currentOutput, f, args...)
}

var pkg = flag.String("pkg", "main", "Package name to use with generated files")
var lang = flag.String("lang", "go", "Language to generate code for (c or go)")

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

		p.addBuiltins()

		err = p.parse()
		f.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s:%d: %v\n", name, p.line, err)
			os.Exit(1)
		}

		if err := p.resolveTypes(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		p.resolveOffsets()
		p.resolveValues()
		p.buf = C.NewBuffer(nil)

		switch *lang {
		case "go":
			p.resolveGoTypes()
			out, err := os.Create(name + ".go")
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			currentOutput = bufio.NewWriter(out)
			p.writeGo(name)
			currentOutput.Flush()
			out.Close()

		case "c":
			p.resolveCTypes()
			out, err := os.Create(name + ".h")
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			currentOutput = bufio.NewWriter(out)
			p.writeCHeader(name)
			currentOutput.Flush()
			out.Close()

			out, err = os.Create(name + ".c")
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			currentOutput = bufio.NewWriter(out)
			p.writeCSource(name)
			currentOutput.Flush()
			out.Close()
		}
	}
}
