package rpc

import (
	"context"
	"unsafe"

	"capnproto.org/go/capnp/v3/internal/mpsc"
	rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"
)

type preparer interface {
	Prepare(rpccp.Message) error
}

type sendQueue mpsc.Queue

func newSendQueue() *sendQueue { return (*sendQueue)(mpsc.New()) }

func (sq *sendQueue) Send(ctx context.Context, p preparer) (err error) {
	cherr := cherrPool.Get()
	defer cherrPool.Put(cherr)

	sq.SendAsync(ctx, p, cherr)

	select {
	case err = <-cherr:
	case <-ctx.Done():
		err = ctx.Err()
	}

	return
}

func (sq *sendQueue) SendAsync(ctx context.Context, p preparer, r ErrorReporter) {
	(*mpsc.Queue)(sq).Send(sendReq{
		er:       r,
		preparer: p,
	})
}

func (sq *sendQueue) Recv(ctx context.Context) preparer {
	v, err := (*mpsc.Queue)(sq).Recv(ctx)
	if err != nil {
		return sendFailure(err)
	}

	return *(*preparer)(unsafe.Pointer(&v))
}

type prepFunc func(rpccp.Message) error

func (prepare prepFunc) Prepare(m rpccp.Message) error { return prepare(m) }

func sendFailure(err error) prepFunc {
	return func(rpccp.Message) error { return err }
}

type sendReq struct {
	er ErrorReporter
	preparer
}

func (r sendReq) Prepare(msg rpccp.Message) (err error) {
	if err = r.preparer.Prepare(msg); r.er != nil {
		r.er.ReportError(err)
	}

	return
}

type errorReporterFunc func(error)

func (report errorReporterFunc) ReportError(err error) { report(err) }

type errChan chan error

func (cherr errChan) ReportError(err error) {
	cherr <- err // buffered
}

var cherrPool = make(errChanPool, 64)

type errChanPool chan errChan

func (pool errChanPool) Get() (cherr errChan) {
	select {
	case cherr = <-pool:
	default:
		cherr = make(errChan, 1)
	}

	return
}

func (pool errChanPool) Put(cherr errChan) {
	select {
	case _, ok := <-cherr:
		if !ok {
			return
		}
	default:
	}

	select {
	case pool <- cherr:
	default:
	}
}
