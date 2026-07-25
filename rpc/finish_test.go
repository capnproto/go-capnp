package rpc

import (
	"context"
	"strings"
	"sync/atomic"
	"testing"

	"capnproto.org/go/capnp/v3"
	transportpkg "capnproto.org/go/capnp/v3/rpc/transport"
	rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"
)

func newFinishMessage(t *testing.T, id answerID) *countingIncomingMessage {
	t.Helper()
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
	return &countingIncomingMessage{message: message}
}

func sendAnswerFinish(t *testing.T, peer Transport, id answerID) {
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
	finish.SetQuestionId(uint32(id))
	if err := out.Send(); err != nil {
		t.Fatal(err)
	}
}

func TestUnknownFinishKeepsConnectionUsable(t *testing.T) {
	left, right := transportpkg.NewPipe(4)
	conn := NewConn(NewTransport(left), nil)
	peer := NewTransport(right)
	t.Cleanup(func() {
		_ = conn.Close()
		_ = peer.Close()
	})

	const (
		unknownID   answerID = 87
		bootstrapID answerID = 88
		nextID      answerID = 89
	)
	requireReturn := func(id answerID) {
		in := recvPeerMessage(t, peer)
		defer in.Release()
		if got := in.Message().Which(); got != rpccp.Message_Which_return {
			t.Fatalf("message after Finish = %v; want Return", got)
		}
		ret, err := in.Message().Return()
		if err != nil {
			t.Fatal(err)
		}
		if got := answerID(ret.AnswerId()); got != id {
			t.Errorf("return answer ID = %d; want %d", got, id)
		}
	}

	sendAnswerFinish(t, peer, unknownID)
	sendBootstrapMarker(t, peer, uint32(bootstrapID))
	requireReturn(bootstrapID)

	// The first Finish retires the answer. The second is now unknown and
	// must not prevent the receive loop from serving another request.
	sendAnswerFinish(t, peer, bootstrapID)
	sendAnswerFinish(t, peer, bootstrapID)
	sendBootstrapMarker(t, peer, uint32(nextID))
	requireReturn(nextID)
}

func TestDuplicateFinishForLiveAnswerFails(t *testing.T) {
	const id answerID = 7
	conn := new(Conn)
	var cancels atomic.Int32
	ans := &ansent{
		cancel: func() { cancels.Add(1) },
		returner: ansReturner{
			c:  conn,
			id: id,
		},
	}
	if !conn.lk.answers.Create(id, ans) {
		t.Fatal("answer ID already present")
	}

	first := newFinishMessage(t, id)
	if err := conn.handleFinish(context.Background(), first); err != nil {
		t.Fatalf("first Finish: %v", err)
	}
	if got := atomic.LoadInt32(&first.releases); got != 1 {
		t.Errorf("first Finish release count = %d; want 1", got)
	}
	if !ans.flags.Contains(finishReceived) {
		t.Fatal("first Finish did not mark live answer")
	}

	second := newFinishMessage(t, id)
	err := conn.handleFinish(context.Background(), second)
	if err == nil || !strings.Contains(err.Error(), "already received finish") {
		t.Fatalf("second Finish error = %v; want duplicate Finish error", err)
	}
	if got := atomic.LoadInt32(&second.releases); got != 1 {
		t.Errorf("second Finish release count = %d; want 1", got)
	}
	if got := cancels.Load(); got != 1 {
		t.Errorf("answer canceled %d times; want 1", got)
	}
}
