package rpc

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"capnproto.org/go/capnp/v3"
	rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"
)

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
