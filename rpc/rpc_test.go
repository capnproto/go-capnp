package rpc_test

import (
	"context"
	"testing"

	"zombiezen.com/go/capnproto2/pogs"
	"zombiezen.com/go/capnproto2/rpc"
	rpccapnp "zombiezen.com/go/capnproto2/std/capnp/rpc"
)

const (
	interfaceID       uint64 = 0xa7317bd7216570aa
	methodID          uint16 = 9
	bootstrapExportID uint32 = 84
)

func TestCloseAbort(t *testing.T) {
	p1, p2 := newPipe(1)
	defer p2.CloseSend()
	defer p2.CloseRecv()
	conn := rpc.NewConn(p1, nil)

	ctx := context.Background()
	if err := conn.Close(); err != nil {
		t.Error("conn.Close():", err)
	}
	msg, release, err := p2.RecvMessage(ctx)
	if err != nil {
		t.Fatal("p2.RecvMessage:", err)
	}
	defer release()
	var rmsg rpcMessage
	if err := pogs.Extract(&rmsg, rpccapnp.Message_TypeID, msg.Struct); err != nil {
		t.Fatal("pogs.Extract(p2.RecvMessage(ctx)):", err)
	}
	if rmsg.Which != rpccapnp.Message_Which_abort {
		t.Fatalf("Received %v message; want abort", rmsg.Which)
	}
	if rmsg.Abort == nil {
		t.Error("Received null abort message")
	} else if rmsg.Abort.Type != rpccapnp.Exception_Type_failed {
		t.Errorf("Received exception type %v; want failed", rmsg.Abort.Type)
	}
}

func TestBootstrap(t *testing.T) {
	p1, p2 := newPipe(1)
	defer p2.CloseSend()
	defer p2.CloseRecv()
	conn := rpc.NewConn(p1, nil)
	defer func() {
		if err := conn.Close(); err != nil {
			t.Error(err)
		}
	}()

	ctx := context.Background()
	conn.Bootstrap(ctx)
	msg, release, err := p2.RecvMessage(ctx)
	if err != nil {
		t.Fatal("p2.RecvMessage:", err)
	}
	defer release()
	var rmsg rpcMessage
	if err := pogs.Extract(&rmsg, rpccapnp.Message_TypeID, msg.Struct); err != nil {
		t.Fatal("pogs.Extract(p2.RecvMessage(ctx)):", err)
	}
	if rmsg.Which != rpccapnp.Message_Which_bootstrap {
		t.Fatalf("Received %v message; want bootstrap", rmsg.Which)
	}
}

type rpcMessage struct {
	Which         rpccapnp.Message_Which
	Unimplemented *rpcMessage
	Abort         *rpcException
	Bootstrap     *rpcBootstrap
}

type rpcException struct {
	Reason string
	Type   rpccapnp.Exception_Type
}

type rpcBootstrap struct {
	QuestionID uint32 `capnp:"questionId"`
}
