// Package exc provides an error type for capnp exceptions.
package exc

import (
	"errors"
	"fmt"
	"strconv"
)

// Exception is an error that designates a Cap'n Proto exception.
type Exception struct {
	Type   Type
	Prefix string
	Cause  error
}

// New creates a new error that formats as "<prefix>: <msg>".
// The type can be recovered using the TypeOf() function.
func New(typ Type, prefix, msg string) Exception {
	return Exception{typ, prefix, errors.New(msg)}
}

func (e Exception) Error() string {
	if e.Prefix == "" {
		return e.Cause.Error()
	}

	return fmt.Sprintf("%s: %v", e.Prefix, e.Cause)
}

func (e Exception) Unwrap() error { return e.Cause }

func (e Exception) GoString() string {
	return fmt.Sprintf("errors.Error{Type: %s, Prefix: %q, Cause: fmt.Errorf(%q)}",
		e.Type.GoString(),
		e.Prefix,
		e.Cause)
}

// Annotate is creates a new error that formats as "<prefix>: <msg>: <e>".
// If e.Prefix == prefix, the prefix will not be duplicated.
// The returned Error.Type == e.Type.
func (e Exception) Annotate(prefix, msg string) Exception {
	if prefix != e.Prefix {
		return Exception{e.Type, prefix, fmt.Errorf("%s: %w", msg, e)}
	}

	return Exception{e.Type, prefix, fmt.Errorf("%s: %w", msg, e.Cause)}
}

// Annotate creates a new error that formats as "<prefix>: <msg>: <err>".
// If err has the same prefix, then the prefix won't be duplicated.
// The returned error's type will match err's type.
func Annotate(prefix, msg string, err error) error {
	if err == nil {
		panic("Annotate on nil error") // TODO:  return nil?
	}

	if ce, ok := err.(Exception); ok {
		return ce.Annotate(prefix, msg)
	}

	return Exception{Failed, prefix, fmt.Errorf("%s: %w", msg, err)}
}

// TypeOf returns err's type if err was created by this package or
// Failed if it was not.
func TypeOf(err error) Type {
	ce, ok := err.(Exception)
	if !ok {
		return Failed
	}
	return ce.Type
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
