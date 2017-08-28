package rpc_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"zombiezen.com/go/capnproto2"
	"zombiezen.com/go/capnproto2/rpc"
	rpccapnp "zombiezen.com/go/capnproto2/std/capnp/rpc"
)

type pipe struct {
	r  <-chan pipeMsg
	rc chan struct{} // close to hang up reads

	w  chan<- pipeMsg
	wc <-chan struct{} // closed when writes are no longer listened to
}

type pipeMsg struct {
	msg     rpccapnp.Message
	release capnp.ReleaseFunc
}

func newPipe(n int) (p1, p2 *pipe) {
	ch1 := make(chan pipeMsg, n)
	ch2 := make(chan pipeMsg, n)
	close1 := make(chan struct{})
	close2 := make(chan struct{})
	return &pipe{r: ch1, w: ch2, rc: close1, wc: close2},
		&pipe{r: ch2, w: ch1, rc: close2, wc: close1}
}

func (p *pipe) NewMessage(ctx context.Context) (_ rpccapnp.Message, send func() error, cancel func(), _ error) {
	msg, seg, _ := capnp.NewMessage(capnp.MultiSegment(nil))
	rmsg, _ := rpccapnp.NewRootMessage(seg)
	send = func() error {
		pm := pipeMsg{rmsg, func() { msg.Reset(nil) }}
		select {
		case p.w <- pm:
			return nil
		case <-p.wc:
			p.w = nil
			return errors.New("rpc pipe: send on closed pipe")
		case <-ctx.Done():
			return fmt.Errorf("rpc pipe: %v", ctx.Err())
		}
	}
	cancel = func() {
		msg.Reset(nil)
	}
	return rmsg, send, cancel, nil
}

func (p *pipe) CloseSend() error {
	close(p.w)
	return nil
}

func (p *pipe) RecvMessage(ctx context.Context) (rpccapnp.Message, capnp.ReleaseFunc, error) {
	select {
	case pm, ok := <-p.r:
		if !ok {
			return rpccapnp.Message{}, nil, errors.New("rpc pipe: receive on closed pipe")
		}
		return pm.msg, pm.release, nil
	case <-p.rc:
		return rpccapnp.Message{}, nil, errors.New("rpc pipe: receive interrupted by close")
	case <-ctx.Done():
		return rpccapnp.Message{}, nil, fmt.Errorf("rpc pipe: %v", ctx.Err())
	}
}

func (p *pipe) CloseRecv() error {
	close(p.rc)
	for {
		select {
		case _, ok := <-p.r:
			if !ok {
				return nil
			}
		default:
			return nil
		}
	}
}

func TestPipeTransport(t *testing.T) {
	testTransport(t, func() (t1, t2 rpc.Transport, err error) {
		p1, p2 := newPipe(1)
		return p1, p2, nil
	})
}
