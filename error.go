package capnp

import (
	"fmt"

	"capnproto.org/go/capnp/v3/exc"
)

// Unimplemented returns an error that formats as the given text and
// will report true when passed to IsUnimplemented.
func Unimplemented(s string) error {
	return exc.New(exc.Unimplemented, "", s)
}

// IsUnimplemented reports whether e indicates that functionality is unimplemented.
func IsUnimplemented(e error) bool {
	return exc.TypeOf(e) == exc.Unimplemented
}

// Disconnected returns an error that formats as the given text and
// will report true when passed to IsDisconnected.
func Disconnected(s string) error {
	return exc.New(exc.Disconnected, "", s)
}

// IsDisconnected reports whether e indicates a failure due to loss of a necessary capability.
func IsDisconnected(e error) bool {
	return exc.TypeOf(e) == exc.Disconnected
}

func newError(msg string) error {
	return exc.New(exc.Failed, "capnp", msg)
}

func errorf(format string, args ...interface{}) error {
	return newError(fmt.Sprintf(format, args...))
}

type annotater struct {
	err error
}

func annotate(err error) annotater {
	return annotater{err}
}

func (a annotater) errorf(format string, args ...interface{}) error {
	return exc.Annotate("capnp", fmt.Sprintf(format, args...), a.err)
}
