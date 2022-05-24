package rpc

import (
	"context"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/internal/mpsc"
	"capnproto.org/go/capnp/v3/rpc/internal/proc"
	rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"
)

type sender struct {
	tx    *mpsc.Tx[sendJob]
	trans Transport
}

func spawnSender(ctx context.Context, trans Transport) proc.Handle[*sender] {
	q := mpsc.New[sendJob]()
	return proc.Spawn(
		ctx,
		&sender{
			tx:    &q.Tx,
			trans: trans,
		},
		func(ctx context.Context, self proc.Self[*sender]) {
			rx := &q.Rx
			for {
				job, err := rx.Recv(ctx)
				if err != nil {
					break
				}
				job.onSent(job.send())
			}
			self.BeginShutdown()
			for {
				job, ok := rx.TryRecv()
				if !ok {
					break
				}
				job.onSent(ErrConnClosed)
			}
		},
	)

}

type sendJob struct {
	send   func() error
	onSent func(error)
}

type sendArgs struct {
	Msg     rpccp.Message
	Send    func(onSent func(err error))
	Release capnp.ReleaseFunc
}

func startMessage(ctx context.Context, h proc.Handle[*sender]) (sendArgs, error) {
	var (
		msg     rpccp.Message
		send    func() error
		release capnp.ReleaseFunc
		err     error
	)
	ok := h.WithLive(func(s *sender) {
		msg, send, release, err = s.trans.NewMessage(ctx)
	})
	if !ok {
		return sendArgs{}, ErrConnClosed
	}
	if err != nil {
		return sendArgs{}, err
	}
	return sendArgs{
		Msg:     msg,
		Release: release,
		Send: func(onSent func(error)) {
			ok := h.WithLive(func(s *sender) {
				s.tx.Send(sendJob{
					send:   send,
					onSent: onSent,
				})
			})
			if !ok {
				onSent(ErrConnClosed)
			}
		},
	}, nil
}
