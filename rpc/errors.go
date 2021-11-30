package rpc

import (
	goerr "errors"
	"fmt"

	"capnproto.org/go/capnp/v3/internal/errors"
)

const prefix = "rpc"

var (
	// Base errors
	ErrConnClosed        = goerr.New("connection closed")
	ErrNotACapability    = goerr.New("not a capability")
	ErrCapTablePopulated = goerr.New("capability table already populated")

	// RPC exceptions
	ExcClosed        = disconnected(ErrConnClosed)
	ExcAlreadyClosed = failed(goerr.New("close on closed connection"))
)

func failedf(format string, args ...interface{}) errors.Error {
	return failed(fmt.Errorf(format, args...))
}

func failed(err error) errors.Error {
	return exception(errors.Failed, err)
}

func disconnectedf(format string, args ...interface{}) errors.Error {
	return disconnected(fmt.Errorf(format, args...))
}

func disconnected(err error) errors.Error {
	return exception(errors.Disconnected, err)
}

func unimplementedf(format string, args ...interface{}) errors.Error {
	return unimplemented(fmt.Errorf(format, args...))
}

func unimplemented(err error) errors.Error {
	return exception(errors.Unimplemented, err)
}

func annotate(err error, msg string) error {
	return errors.Annotate(prefix, msg, err)
}

func exception(t errors.Type, err error) errors.Error {
	return errors.Error{Type: t, Prefix: prefix, Cause: err}
}
