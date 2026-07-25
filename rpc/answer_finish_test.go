package rpc

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"capnproto.org/go/capnp/v3"
	rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"
)

func newAnswerBootstrapMessage(t *testing.T, id answerID) *countingIncomingMessage {
	t.Helper()

	_, seg := capnp.NewSingleSegmentMessage(nil)
	message, err := rpccp.NewRootMessage(seg)
	if err != nil {
		t.Fatal(err)
	}
	bootstrap, err := message.NewBootstrap()
	if err != nil {
		t.Fatal(err)
	}
	bootstrap.SetQuestionId(uint32(id))
	return &countingIncomingMessage{message: message}
}

func TestFinishDestroysErrorAnswerWithoutReturnMessage(t *testing.T) {
	const id answerID = 11
	conn := new(Conn)
	if !conn.lk.answers.Create(id, errorAnswer(conn, id, errors.New("return allocation failed"))) {
		t.Fatal("answer ID already present")
	}

	_, seg := capnp.NewSingleSegmentMessage(nil)
	message, err := rpccp.NewRootMessage(seg)
	if err != nil {
		t.Fatal(err)
	}
	finish, err := message.NewFinish()
	if err != nil {
		t.Fatal(err)
	}
	finish.SetQuestionId(uint32(id))
	in := &countingIncomingMessage{message: message}

	if err := conn.handleFinish(context.Background(), in); err != nil {
		t.Fatal("handle Finish:", err)
	}
	if got := atomic.LoadInt32(&in.releases); got != 1 {
		t.Fatalf("incoming message releases = %d; want 1", got)
	}
	if _, ok := conn.lk.answers.Find(id); ok {
		t.Fatal("error answer remains live after Finish")
	}
}

func TestBootstrapReuseAfterReturnAllocationFailure(t *testing.T) {
	const id answerID = 12
	conn := &Conn{
		transport: newMessageErrorTransport{err: errors.New("return allocation failed")},
	}
	if !conn.lk.answers.Create(id, errorAnswer(conn, id, errors.New("existing answer"))) {
		t.Fatal("answer ID already present")
	}
	in := newAnswerBootstrapMessage(t, id)

	err := conn.handleBootstrap(in)
	if err == nil || !strings.Contains(err.Error(), "answer ID 12 reused") {
		t.Fatalf("handle Bootstrap error = %v; want answer reuse error", err)
	}
	if got := atomic.LoadInt32(&in.releases); got != 1 {
		t.Fatalf("incoming message releases = %d; want 1", got)
	}
	if _, ok := conn.lk.answers.Find(id); !ok {
		t.Fatal("existing answer removed after rejected Bootstrap")
	}
}

func TestBootstrapReuseReleasesUnsentReturn(t *testing.T) {
	const id answerID = 13
	transport := newFailingSendTransport(nil)
	conn := &Conn{transport: transport}
	if !conn.lk.answers.Create(id, errorAnswer(conn, id, errors.New("existing answer"))) {
		t.Fatal("answer ID already present")
	}
	in := newAnswerBootstrapMessage(t, id)

	err := conn.handleBootstrap(in)
	if err == nil || !strings.Contains(err.Error(), "answer ID 13 reused") {
		t.Fatalf("handle Bootstrap error = %v; want answer reuse error", err)
	}
	if got := atomic.LoadInt32(&in.releases); got != 1 {
		t.Fatalf("incoming message releases = %d; want 1", got)
	}
	select {
	case <-transport.firstRelease:
	case <-time.After(2 * time.Second):
		t.Fatal("unsent Bootstrap Return was not released")
	}
}
