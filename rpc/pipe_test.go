package rpc_test

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"testing"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc"
	rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"
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

	// Variables aren't synchronized because the Transport interface does
	// not require them to be.  Should trigger race detector.
	sent, sendDone, recvDone := false, false, false
	// Since refs is used by Sender and Receiver, then it must be synchronized.
	var (
		refsMu sync.Mutex
		refs   int = 1
	)
	send = func() error {
		if sendDone {
			panic("send after release")
		}
		if sent {
			panic("double send")
		}
		sent = true
		refsMu.Lock()
		refs++
		refsMu.Unlock()
		pm := pipeMsg{
			msg: rmsg,
			release: func() {
				if recvDone {
					return
				}
				recvDone = true
				refsMu.Lock()
				r := refs - 1
				refs = r
				refsMu.Unlock()
				if r == 0 {
					msg.Reset(nil)
				}
			},
		}
		select {
		case p.w <- pm:
			return nil
		case <-p.wc:
			p.w = nil
			refsMu.Lock()
			r := refs - 1
			refs = r
			refsMu.Unlock()
			if r == 0 {
				msg.Reset(nil)
			}
			return errors.New("rpc pipe: send on closed pipe")
		case <-ctx.Done():
			refsMu.Lock()
			r := refs - 1
			refs = r
			refsMu.Unlock()
			if r == 0 {
				msg.Reset(nil)
			}
			return fmt.Errorf("rpc pipe: %v", ctx.Err())
		}
	}
	release = func() {
		if sendDone {
			return
		}
		sendDone = true
		delete(p.msgs, caller)
		refsMu.Lock()
		r := refs - 1
		refs = r
		refsMu.Unlock()
		if r == 0 {
			msg.Reset(nil)
		}
	}
	return rmsg, send, release, nil
}

type newMessageCaller struct {
	file string
	line int
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

func (p *pipe) Close() error {
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
		panic("Close called before releasing all messages.  Unreleased: " + string(callers))
	}
	close(p.w)
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
