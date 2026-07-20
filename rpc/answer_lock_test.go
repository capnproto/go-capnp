package rpc

import (
	"context"
	"errors"
	"testing"
	"time"

	"capnproto.org/go/capnp/v3"
	transportpkg "capnproto.org/go/capnp/v3/rpc/transport"
)

type blockingServerPipelineCaller struct {
	admitted chan struct{}
	release  <-chan struct{}
}

func (c *blockingServerPipelineCaller) PipelineSend(context.Context, []capnp.PipelineOp, capnp.Send) (*capnp.Answer, capnp.ReleaseFunc) {
	close(c.admitted)
	<-c.release
	return capnp.ErrorAnswer(capnp.Method{}, errors.New("test call")), func() {}
}

func (c *blockingServerPipelineCaller) PipelineRecv(context.Context, []capnp.PipelineOp, capnp.Recv) capnp.PipelineCaller {
	close(c.admitted)
	<-c.release
	return nil
}

type signalingRPCPtrResolver struct{ result chan<- error }

func (r signalingRPCPtrResolver) Fulfill(capnp.Ptr) { r.result <- nil }
func (r signalingRPCPtrResolver) Reject(err error)  { r.result <- err }

type blockedServerAnswer struct {
	conn           *Conn
	peer           Transport
	ans            *ansent
	release        chan struct{}
	callReturned   chan struct{}
	resolverResult chan error
	finishObserved chan struct{}
}

func newBlockedServerAnswer(t *testing.T) *blockedServerAnswer {
	t.Helper()
	left, right := transportpkg.NewPipe(4)
	f := &blockedServerAnswer{
		conn:           NewConn(NewTransport(left), nil),
		peer:           NewTransport(right),
		release:        make(chan struct{}),
		callReturned:   make(chan struct{}),
		resolverResult: make(chan error, 1),
		finishObserved: make(chan struct{}),
	}
	t.Cleanup(func() {
		_ = f.conn.Close()
		_ = f.peer.Close()
	})

	admitted := make(chan struct{})
	caller := &blockingServerPipelineCaller{admitted: admitted, release: f.release}
	f.conn.withLocked(func(c *lockedConn) {
		ret, send, releaser, err := f.conn.newReturn()
		if err != nil {
			t.Fatal(err)
		}
		f.ans = &ansent{
			returner: ansReturner{c: f.conn, id: 0, ret: ret, msgReleaser: releaser},
			pcall:    caller,
			sendMsg:  send,
			cancel:   func() { close(f.finishObserved) },
		}
		f.ans.promise = capnp.NewPromise(capnp.Method{}, caller, signalingRPCPtrResolver{result: f.resolverResult})
		if !c.lk.answers.Create(0, f.ans) {
			t.Fatal("answer ID already present")
		}
		f.conn.tasks.Add(1)
	})
	go func() {
		defer close(f.callReturned)
		_, _ = f.ans.promise.Answer().PipelineSend(context.Background(), nil, capnp.Send{Method: capnp.Method{}})
	}()
	<-admitted
	return f
}

func (f *blockedServerAnswer) startReturn(returnErr error) <-chan struct{} {
	f.ans.returner.PrepareReturn(returnErr)
	completed := make(chan struct{})
	go func() {
		f.ans.returner.Return()
		close(completed)
	}()
	return completed
}

func sendFinish(t *testing.T, peer Transport) {
	t.Helper()
	out, err := peer.NewMessage()
	if err != nil {
		t.Fatal(err)
	}
	defer out.Release()
	finish, err := out.Message().NewFinish()
	if err != nil {
		t.Fatal(err)
	}
	finish.SetQuestionId(0)
	if err := out.Send(); err != nil {
		t.Fatal(err)
	}
}

func TestServerAnswerResolvesPromiseOutsideConnLock(t *testing.T) {
	t.Run("fulfill", func(t *testing.T) {
		testServerAnswerResolvesPromiseOutsideConnLock(t, nil)
	})
	t.Run("reject", func(t *testing.T) {
		testServerAnswerResolvesPromiseOutsideConnLock(t, errors.New("return failed"))
	})
}

func testServerAnswerResolvesPromiseOutsideConnLock(t *testing.T, returnErr error) {
	f := newBlockedServerAnswer(t)
	returnCompleted := f.startReturn(returnErr)
	resolvedErr := <-f.resolverResult
	if returnErr == nil && resolvedErr != nil {
		t.Fatalf("Promise resolver rejected successful Return: %v", resolvedErr)
	}
	if returnErr != nil && !errors.Is(resolvedErr, returnErr) {
		t.Fatalf("Promise resolver error = %v; want %v", resolvedErr, returnErr)
	}

	lockObserved := make(chan struct{})
	go func() {
		f.conn.withLocked(func(c *lockedConn) {
			got, ok := c.lk.answers.Find(0)
			if !ok || got != f.ans || !got.flags.Contains(resultsReady) || got.promise != nil {
				t.Error("answer did not retain completion ownership while resolving its promise")
			}
		})
		close(lockObserved)
	}()
	select {
	case <-lockObserved:
	case <-time.After(time.Second):
		t.Fatal("Conn.lk held while Promise resolution drained admitted calls")
	}
	select {
	case <-returnCompleted:
		t.Fatal("Return completed before admitted pipeline call yielded")
	default:
	}

	close(f.release)
	<-f.callReturned
	<-returnCompleted
}

func TestServerAnswerFinishDuringPromiseDrain(t *testing.T) {
	for _, timing := range []string{"before claim", "after claim"} {
		t.Run(timing, func(t *testing.T) {
			f := newBlockedServerAnswer(t)
			if timing == "before claim" {
				sendFinish(t, f.peer)
				<-f.finishObserved
			}

			returnCompleted := f.startReturn(nil)
			resolvedErr := <-f.resolverResult
			if timing == "before claim" {
				if resolvedErr == nil {
					t.Fatal("Promise fulfilled after Finish won race; want rejection")
				}
			} else {
				if resolvedErr != nil {
					t.Fatalf("Promise rejected after Return claimed completion: %v", resolvedErr)
				}
				sendFinish(t, f.peer)
				<-f.finishObserved
			}

			f.conn.withLocked(func(c *lockedConn) {
				got, ok := c.lk.answers.Find(0)
				if !ok || got != f.ans || got.flags.Contains(returnSent) || !got.flags.Contains(finishReceived) {
					t.Error("Finish destroyed or corrupted answer before Return completed")
				}
			})
			close(f.release)
			<-f.callReturned
			<-returnCompleted
			f.conn.withLocked(func(c *lockedConn) {
				if _, ok := c.lk.answers.Find(0); ok {
					t.Error("answer retained after Return and Finish completed")
				}
			})
		})
	}
}

func TestServerAnswerShutdownWaitsForPromiseDrain(t *testing.T) {
	f := newBlockedServerAnswer(t)
	returnCompleted := f.startReturn(nil)
	if err := <-f.resolverResult; err != nil {
		t.Fatalf("Promise resolver error = %v; want nil", err)
	}

	closeReturned := make(chan error, 1)
	go func() { closeReturned <- f.conn.Close() }()
	<-f.conn.bgctx.Done()
	select {
	case <-returnCompleted:
		t.Fatal("Return completed before admitted pipeline call yielded")
	default:
	}
	select {
	case err := <-closeReturned:
		t.Fatalf("Close returned before Return released completion ownership: %v", err)
	default:
	}

	close(f.release)
	<-f.callReturned
	<-returnCompleted
	if err := <-closeReturned; err != nil {
		t.Fatal(err)
	}
}
