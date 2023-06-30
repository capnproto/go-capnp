package rpc

import (
	"errors"

	"capnproto.org/go/capnp/v3/exc"
)

var (
	rpcerr = exc.Annotator("rpc")

	// Base errors
	ErrConnClosed        = errors.New("connection closed")
	ErrNotACapability    = errors.New("not a capability")
	ErrCapTablePopulated = errors.New("capability table already populated")

	// RPC exceptions
	ExcClosed = rpcerr.Disconnected(ErrConnClosed)
)

type errReporter struct {
	Logger
}

func (er errReporter) ReportError(err error) {
	if er.Logger != nil && err != nil {
		er.Logger.Error(err.Error())
	}
}
