package rpc

import (
	"context"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/internal/mpsc"
	"capnproto.org/go/capnp/v3/rpc/internal/proc"
	rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"
)

type sender struct {
	proc  proc.Handle
	trans Transport
	tx    *mpsc.Tx[sendJob]
}

func spawnSender(ctx context.Context, trans Transport) *sender {
	q := mpsc.New[sendJob]()
	rx := &q.Rx
	return &sender{
		tx:    &q.Tx,
		trans: trans,
		proc: proc.Spawn(ctx, func(ctx context.Context, self proc.Self) {
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
		}),
	}

}

// Stop the sender. Returns
func (s *sender) Stop() {
	s.proc.Cancel()
}

func (s *sender) Done() <-chan struct{} {
	return s.proc.Done()
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

func (s *sender) StartMessage(ctx context.Context) (sendArgs, error) {
	var (
		msg     rpccp.Message
		send    func() error
		release capnp.ReleaseFunc
		err     error
	)
	ok := s.proc.WithLive(func() {
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
			ok := s.proc.WithLive(func() {
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
