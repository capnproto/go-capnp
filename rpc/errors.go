package rpc

import (
	goerr "errors"
	"fmt"

	"capnproto.org/go/capnp/v3/exc"
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

func failedf(format string, args ...interface{}) exc.Exception {
	return failed(fmt.Errorf(format, args...))
}

func failed(err error) exc.Exception {
	return exception(exc.Failed, err)
}

func disconnectedf(format string, args ...interface{}) exc.Exception {
	return disconnected(fmt.Errorf(format, args...))
}

func disconnected(err error) exc.Exception {
	return exception(exc.Disconnected, err)
}

func unimplementedf(format string, args ...interface{}) exc.Exception {
	return unimplemented(fmt.Errorf(format, args...))
}

func unimplemented(err error) exc.Exception {
	return exception(exc.Unimplemented, err)
}

func annotate(err error, msg string) error {
	return exc.Annotate(prefix, msg, err)
}

func exception(t exc.Type, err error) exc.Exception {
	return exc.Exception{Type: t, Prefix: prefix, Cause: err}
}

func annotatef(err error, format string, args ...interface{}) error {
	return exc.Annotate("rpc", fmt.Sprintf(format, args...), err)
}

type annotatingErrReporter struct {
	ErrorReporter
}

func (er annotatingErrReporter) ReportError(err error) {
	if er.ErrorReporter != nil && err != nil {
		er.ErrorReporter.ReportError(err)
	}
}

func (er annotatingErrReporter) reportf(format string, args ...interface{}) {
	er.ReportError(fmt.Errorf(format, args...))
}

func (er annotatingErrReporter) annotatef(err error, format string, args ...interface{}) {
	er.ReportError(annotatef(err, format, args...))
}
