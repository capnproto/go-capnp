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
	unionType
)

func (t *typ) String() string {
	return t.name
}

type typ struct {
	typ        typeType
	name       string
	enumPrefix string
	comment    string
	fields     []*field
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

	case '0', '\'', '"', 'a':
		v := &value{
			typ:    tok.typ,
			string: tok.str,
		}
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
			return "C.Bitset"
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
			panic(fmt.Errorf("expected : or ; got %s", tok))
		}

		if f.typestr == "struct{}" {
			f.typestr = ""
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
					typ:        unionType,
					name:       sprintf("%s_%s", t.name, f.name),
					enumPrefix: sprintf("%s_", t.name),
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
	p.types = append(p.types, &typ{typ: bitsetType, name: "C.Bitset"})
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
			t.name = "[]C.Pointer"
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

func (p *file) resolveOffsets() error {
	for _, t := range p.types {
		if t.typ != structType {
			continue
		}

		fields := ordinalFields(t.fields)
		sort.Sort(fields)

	next_field:
		for i, f := range fields {
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
	return nil
}

var pkg = flag.String("pkg", "main", "Package name to use with generated files")

type printer func(string, ...interface{})

func (v *value) String(t *typ, marshalled bool) string {
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
		out := ""
		if marshalled {
			out += "C.Must(NewBitset(C.NewMemory, "
		}
		out += "C.Bitset{"
		for i, b := range set {
			if i > 0 {
				out += ", "
			}
			out += sprintf("%#02x", b)
		}
		out += "}"
		if marshalled {
			out += "))"
		}
		return out

	case stringType:
		if v.typ != '"' {
			panic(fmt.Errorf("unexpected value %v in string", v))
		}

		if marshalled {
			return sprintf("C.Must(C.NewString(C.NewMemory, %s))", strconv.Quote(v.string))
		} else {
			return sprintf("%s", strconv.Quote(v.string))
		}

	case listType:
		if v.typ == '"' && t.listType.typ == uint8Type {
			// Data fields with a string value
			if marshalled {
				return sprintf("C.Must(C.NewUInt8List(C.NewMemory, []byte(%s)))", strconv.Quote(v.string))
			} else {
				return sprintf("[]byte(%s)", strconv.Quote(v.string))
			}
		}

		if t.listType.typ == voidType {
			if marshalled {
				return sprintf("C.Must(C.NewVoidList(C.NewMemory, make([]struct{}, %d)))", len(v.fields))
			} else {
				return sprintf("make([]struct{}, %d)", len(v.fields))
			}
		}

		if v.typ != '[' {
			panic(fmt.Errorf("unexpected value %v in list", v))
		}

		listTypeName := t.name
		innerMarshalled := false

		if t.listType.typ == voidType {
			// otherwise we get []struct{}{nil}
			listTypeName = "[](struct{})"
		}

		out := ""
		if marshalled {
			switch t.listType.typ {
			case int8Type:
				out += "C.Must(C.NewInt8List(C.NewMemory, "
			case uint8Type:
				out += "C.Must(C.NewUInt8List(C.NewMemory, "
			case int16Type:
				out += "C.Must(C.NewInt16List(C.NewMemory, "
			case uint16Type:
				out += "C.Must(C.NewUInt16List(C.NewMemory, "
			case enumType:
				out += "C.Must(C.NewUInt16List(C.NewMemory, "
				listTypeName = "[]uint16"
				innerMarshalled = true
			case int32Type:
				out += "C.Must(C.NewInt32List(C.NewMemory, "
			case uint32Type:
				out += "C.Must(C.NewUInt32List(C.NewMemory, "
			case int64Type:
				out += "C.Must(C.NewInt64List(C.NewMemory, "
			case uint64Type:
				out += "C.Must(C.NewUInt64List(C.NewMemory, "
			case float32Type:
				out += "C.Must(C.NewFloat32List(C.NewMemory, "
			case float64Type:
				out += "C.Must(C.NewFloat64List(C.NewMemory, "
			case stringType:
				out += "C.Must(C.NewStringList(C.NewMemory, "
			case bitsetType:
				out += "C.Must(C.NewBitsetList(C.NewMemory, "
			case listType:
				out += "C.Must(C.NewPointerList(C.NewMemory, "
				listTypeName = "[]C.Pointer"
				innerMarshalled = true
			default:
				panic("unhandled")
			}

		} else if t.listType.typ == listType {
			listTypeName = "[]C.Pointer"
			innerMarshalled = true
		}

		out += sprintf("%s{", listTypeName)
		for i, v := range v.fields {
			if i > 0 {
				out += ", "
			}
			out += v.String(t.listType, innerMarshalled)
		}
		out += "}"

		if marshalled {
			out += "))"
		}
		return out

	case structType:
		if v.typ != '(' {
			panic(fmt.Errorf("unexpected value %v in struct", v))
		}
		out := sprintf("func() %s {\n", t.name)
		out += sprintf("p, _ := %s(C.NewMemory)\n", pfxname("new", t.name))
		for _, v := range v.fields {
			f := t.findField(v.name)
			if f == nil {
				panic(fmt.Errorf("can't find field %s in type %s", v.name, t.name))
			}
			if f.typ.typ != voidType {
				out += sprintf("p.%s(", pfxname("set", f.name))
				out += v.String(f.typ, false)
				out += ")\n"
			}
		}
		out += "return p\n"
		out += "}()"
		if marshalled {
			out += ".Ptr"
		}
		return out

	case enumType:
		if v.typ != 'a' {
			panic(fmt.Errorf("unexpected value %v in enum for %v", v, t))
		}

		// Can be either a constant (use v.string) or an enum value (use type_v.string)
		str := v.string
		if f := t.findField(str); f != nil {
			str = t.enumPrefix + v.string
		}

		if marshalled {
			return sprintf("uint16(%s)", str)
		} else {
			return str
		}

	case boolType, float32Type, float64Type,
		int8Type, uint8Type, int16Type, uint16Type,
		int32Type, uint32Type, int64Type, uint64Type:

		// number types, can be number or constants
		switch v.typ {
		case '0', 'a':
			return v.string
		default:
			panic(fmt.Errorf("unexpected value %v in number for %v", v, t))
		}

	default:
		panic("unhandled")
	}
}

func (f *field) writeGetterXorDefault(out printer) {
	if f.value != nil {
		out("^ %s", f.value.String(f.typ, false))
	}
	out("\n")
}

func (f *field) writeGetter(out printer, ptr string) {
	switch f.typ.typ {
	case voidType:
		// nothing to do
	case boolType:
		def := 0
		if f.value != nil && f.value.string == "true" {
			def = 1
		}
		out("return (C.ReadUInt8(%s, %d) & %d) != %d\n", ptr, f.offset/8, f.offset%8, def)
	case int8Type:
		out("return int8(C.ReadUInt8(%s, %d))", ptr, f.offset/8)
		f.writeGetterXorDefault(out)
	case uint8Type:
		out("return C.ReadUInt8(%s, %d)", ptr, f.offset/8)
		f.writeGetterXorDefault(out)
	case int16Type:
		out("return int16(C.ReadUInt16(%s, %d))", ptr, f.offset/8)
		f.writeGetterXorDefault(out)
	case uint16Type:
		out("return C.ReadUInt16(%s, %d)", ptr, f.offset/8)
		f.writeGetterXorDefault(out)
	case int32Type:
		out("return int32(C.ReadUInt32(%s, %d))", ptr, f.offset/8)
		f.writeGetterXorDefault(out)
	case uint32Type:
		out("return C.ReadUInt32(%s, %d)", ptr, f.offset/8)
		f.writeGetterXorDefault(out)
	case int64Type:
		out("return int64(C.ReadUInt64(%s, %d))", ptr, f.offset/8)
		f.writeGetterXorDefault(out)
	case uint64Type:
		out("return C.ReadUInt64(%s, %d)", ptr, f.offset/8)
		f.writeGetterXorDefault(out)
	case enumType, unionType:
		out("return %s(C.ReadUInt16(%s, %d))", f.typ.name, ptr, f.offset/8)
		f.writeGetterXorDefault(out)

	case float32Type:
		if f.value != nil {
			out(`u := C.ReadUInt32(%s, %d)
			u ^= M.Float32bits(%s)
			return M.Float32frombits(u)
			`, ptr, f.offset/8, f.value.string)
		} else {
			out("return M.Float32frombits(C.ReadUInt32(%s, %d))\n", ptr, f.offset/8)
		}

	case float64Type:
		if f.value != nil {
			out(`u := C.ReadUInt64(%s, %d)
			u ^= M.Float64bits(%s)
			return M.Float64frombits(u)
			`, ptr, f.offset/8, f.value.string)
		} else {
			out("return M.Float64frombits(C.ReadUInt64(%s, %d))\n", ptr, f.offset/8)
		}

	case stringType:
		def := ""
		if f.value != nil {
			def = f.value.string
		}
		out("return C.ToString(%s, %s)\n", ptr, strconv.Quote(def))

	case structType:
		if f.value != nil {
			out(`ret := %s{Ptr: %s}
			if ret.Ptr == nil {
				ret = %s
			}
			return ret
			`, f.typ.name, ptr, f.value.String(f.typ, false))
		} else {
			out("return %s{Ptr: %s}\n", f.typ.name, ptr)
		}

	case interfaceType:
		out("return _%s_remote{Ptr: %s}\n", f.typ.name, ptr)

	case bitsetType:
		if f.value != nil {
			out(`ret := C.ToBitset(%s)
			if ret == nil {
				ret = %s
			}
			return ret
			`, ptr, f.value.String(f.typ, false))
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
			out(`if ret == nil {
				ret = %s
			}
			return ret
			`, f.value.String(f.typ, false))
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

func (f *field) writeDataSetter(out printer, ptr, ret string) {
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
			arg = sprintf("%s ^ %s", arg, f.value.String(f.typ, false))
		case float32Type:
			arg = sprintf("%s ^ M.Float32bits(%s)", arg, f.value.String(f.typ, false))
		case float64Type:
			arg = sprintf("%s ^ M.Float64bits(%s)", arg, f.value.String(f.typ, false))
		case boolType:
			if f.value.string == "true" {
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
	case int64Type, float64Type:
		out("%s C.WriteUInt64(%s, %d, uint64(%s))", ret, ptr, f.offset/8, arg)
	case uint64Type:
		out("%s C.WriteUInt64(%s, %d, %s)", ret, ptr, f.offset/8, arg)
	default:
		panic("unhandled")
	}
}

func (f *field) writeNewPointer(out printer, new, ret string) {
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

func (p *file) writeStruct(t *typ, out printer) {
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

func (p *file) writeInterface(t *typ, out printer) {
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

func (p *file) write(out printer) {
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
