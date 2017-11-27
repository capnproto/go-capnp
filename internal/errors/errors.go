// Package errors provides errors with codes and prefixes.
package errors

import "strconv"

// capnpError holds a Cap'n Proto exception.
type capnpError struct {
	typ    Type
	prefix string
	msg    string
}

// New creates a new error that formats as "<prefix>: <msg>".
// The type can be recovered using the TypeOf() function.
func New(typ Type, prefix, msg string) error {
	return &capnpError{typ, prefix, msg}
}

func (e *capnpError) Error() string {
	if e.prefix == "" {
		return e.msg
	}
	return e.prefix + ": " + e.msg
}

func (e *capnpError) GoString() string {
	return "errors.New(" + e.typ.GoString() + ", " + strconv.Quote(e.prefix) + ", " + strconv.Quote(e.msg) + ")"
}

// Annotate creates a new error that formats as "<prefix>: <msg>: <err>".
// If err has the same prefix, then the prefix won't be duplicated.
// The returned error's type will match err's type.
func Annotate(prefix, msg string, err error) error {
	if err == nil {
		panic("Annotate on nil error")
	}
	ce, ok := err.(*capnpError)
	if !ok {
		return &capnpError{Failed, prefix, msg + ": " + err.Error()}
	}
	if prefix != ce.prefix {
		return &capnpError{ce.typ, prefix, msg + ": " + err.Error()}
	}
	return &capnpError{ce.typ, prefix, msg + ": " + ce.msg}
}

// TypeOf returns err's type if err was created by this package or
// Failed if it was not.
func TypeOf(err error) Type {
	ce, ok := err.(*capnpError)
	if !ok {
		return Failed
	}
	return ce.typ
}

// Type indicates the type of error, mirroring those in rpc.capnp.
type Type int

// Error types.
const (
	Failed        Type = 0
	Overloaded    Type = 1
	Disconnected  Type = 2
	Unimplemented Type = 3
)

// String returns the lowercased Go constant name, or a string in the
// form "type(X)" where X is the value of typ for any unrecognized type.
func (typ Type) String() string {
	switch typ {
	case Failed:
		return "failed"
	case Overloaded:
		return "overloaded"
	case Disconnected:
		return "disconnected"
	case Unimplemented:
		return "unimplemented"
	default:
		var buf [26]byte
		s := append(buf[:0], "type("...)
		s = strconv.AppendInt(s, int64(typ), 10)
		s = append(s, ')')
		return string(s)
	}
}

// GoString returns the Go constant name, or a string in the form
// "Type(X)" where X is the value of typ for any unrecognized type.
func (typ Type) GoString() string {
	switch typ {
	case Failed:
		return "Failed"
	case Overloaded:
		return "Overloaded"
	case Disconnected:
		return "Disconnected"
	case Unimplemented:
		return "Unimplemented"
	default:
		var buf [26]byte
		s := append(buf[:0], "Type("...)
		s = strconv.AppendInt(s, int64(typ), 10)
		s = append(s, ')')
		return string(s)
	}
}
