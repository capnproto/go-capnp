package rpc

import (
	"errors"
	"fmt"

	"capnproto.org/go/capnp/v3/exc"
)

var (
	rpcerr = exc.Annotator("rpc")

	// Base errors
	ErrConnClosed        = errors.New("connection closed")
	ErrNotACapability    = errors.New("not a capability")
	ErrCapTablePopulated = errors.New("capability table already populated")

	// RPC exceptions
	ExcClosed        = rpcerr.Disconnected(ErrConnClosed)
	ExcAlreadyClosed = rpcerr.Failed(errors.New("close on closed connection"))
)

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
	er.ReportError(rpcerr.Annotatef(err, format, args...))
}
