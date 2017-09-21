package rpc

import (
	"context"
	"errors"

	"zombiezen.com/go/capnproto2"
)

type question struct {
	id        questionID
	bootstrap bool
	conn      *Conn

	p       *capnp.Promise    // resolving must be done while holding conn.mu
	release capnp.ReleaseFunc // written before resolving p
}

func (q *question) PipelineSend(ctx context.Context, transform []capnp.PipelineOp, s capnp.Send) (*capnp.Answer, capnp.ReleaseFunc) {
	return capnp.ErrorAnswer(errors.New("TODO(soon)")), func() {}
}

func (q *question) PipelineRecv(ctx context.Context, transform []capnp.PipelineOp, r capnp.Recv) (*capnp.Answer, capnp.ReleaseFunc) {
	return capnp.ErrorAnswer(errors.New("TODO(soon)")), func() {}
}
