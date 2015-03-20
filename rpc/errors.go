package rpc

import (
	"errors"

	"zombiezen.com/go/capnproto/rpc/internal/rpc"
)

// An Exception is a Cap'n Proto RPC error.
type Exception struct {
	rpc.Exception
}

// Error returns the exception's reason.
func (e Exception) Error() string {
	return "rpc exception: " + e.Reason()
}

// An Abort is a hang-up by a remote vat.
type Abort Exception

// Error returns the exception's reason.
func (a Abort) Error() string {
	return "rpc: aborted by remote: " + a.Reason()
}

// toException sets fields on exc to match err.
func toException(exc rpc.Exception, err error) {
	if ee, ok := err.(Exception); ok {
		exc.SetReason(ee.Reason())
		exc.SetType(ee.Type())
		return
	}

	exc.SetReason(err.Error())
	exc.SetType(rpc.Exception_Type_failed)
}

// Errors
var (
	ErrConnClosed = errors.New("rpc: connection closed")
)

// Internal errors
var (
	errQuestionReused  = errors.New("rpc: question ID reused")
	errNoMainInterface = errors.New("rpc: no bootstrap interface")
	errBadTarget       = errors.New("rpc: target not found")
	errShutdown        = errors.New("rpc: shutdown")
	errCallCanceled    = errors.New("rpc: call canceled")
)

type bootstrapError struct {
	err error
}

func (e bootstrapError) Error() string {
	return "rpc bootstrap:" + e.err.Error()
}
