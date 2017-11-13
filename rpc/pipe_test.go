package rpc_test

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"testing"

	"zombiezen.com/go/capnproto2"
	"zombiezen.com/go/capnproto2/rpc"
	rpccp "zombiezen.com/go/capnproto2/std/capnp/rpc"
)

type pipe struct {
	r  <-chan pipeMsg
	rc chan struct{} // close to hang up reads

	w    chan<- pipeMsg
	wc   <-chan struct{} // closed when writes are no longer listened to
	msgs map[*newMessageCaller]struct{}
}

type pipeMsg struct {
	msg     rpccp.Message
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

func (p *pipe) NewMessage(ctx context.Context) (_ rpccp.Message, send func() error, release capnp.ReleaseFunc, _ error) {
	msg, seg, _ := capnp.NewMessage(capnp.MultiSegment(nil))
	rmsg, _ := rpccp.NewRootMessage(seg)
	_, file, line, _ := runtime.Caller(1)
	caller := &newMessageCaller{file, line}
	if p.msgs == nil {
		p.msgs = make(map[*newMessageCaller]struct{})
	}
	p.msgs[caller] = struct{}{}

	// Variables don't need to be synchronized, since a Sender must be
	// used by one goroutine.
	done, sent := false, false
	send = func() error {
		if done {
			panic("send after release")
		}
		if sent {
			panic("double send")
		}
		sent = true
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
	release = func() {
		if done {
			return
		}
		done = true
		delete(p.msgs, caller)
		if !sent {
			msg.Reset(nil)
		}
	}
	return rmsg, send, release, nil
}

type newMessageCaller struct {
	file string
	line int
}

func (p *pipe) CloseSend() error {
	if len(p.msgs) > 0 {
		var callers []byte
		for c := range p.msgs {
			if len(callers) > 0 {
				callers = append(callers, ", "...)
			}
			if c.file == "" && c.line == 0 {
				callers = append(callers, "<???>"...)
				continue
			}
			callers = append(callers, c.file...)
			callers = append(callers, ':')
			callers = strconv.AppendInt(callers, int64(c.line), 10)
		}
		panic("CloseSend called before releasing all messages.  Unreleased: " + string(callers))
	}
	close(p.w)
	return nil
}

func (p *pipe) RecvMessage(ctx context.Context) (rpccp.Message, capnp.ReleaseFunc, error) {
	select {
	case pm, ok := <-p.r:
		if !ok {
			return rpccp.Message{}, nil, errors.New("rpc pipe: receive on closed pipe")
		}
		return pm.msg, pm.release, nil
	case <-p.rc:
		return rpccp.Message{}, nil, errors.New("rpc pipe: receive interrupted by close")
	case <-ctx.Done():
		return rpccp.Message{}, nil, fmt.Errorf("rpc pipe: %v", ctx.Err())
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
