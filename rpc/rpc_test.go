package rpc_test

import (
	"context"
	"testing"

	"zombiezen.com/go/capnproto2"
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

	// 1. Read bootstrap
	client := conn.Bootstrap(ctx)
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
	qid := rmsg.Bootstrap.QuestionID

	// 2. Write back a return
	msg, send, release, err := p2.NewMessage(ctx)
	if err != nil {
		t.Fatal("p2.NewMessage():", err)
	}
	iptr := capnp.NewInterface(msg.Segment(), 0)
	err = pogs.Insert(rpccapnp.Message_TypeID, msg.Struct, &rpcMessage{
		Which: rpccapnp.Message_Which_return,
		Return: &rpcReturn{
			AnswerID: qid,
			Which:    rpccapnp.Return_Which_results,
			Results: &rpcPayload{
				Content: iptr.ToPtr(),
				CapTable: []rpcCapDescriptor{
					{
						Which:        rpccapnp.CapDescriptor_Which_senderHosted,
						SenderHosted: bootstrapExportID,
					},
				},
			},
		},
	})
	if err != nil {
		release()
		t.Fatal("pogs.Insert(p2.NewMessage(), &rpcMessage{...}):", err)
	}
	err = send()
	release()
	if err != nil {
		t.Fatal("send():", err)
	}

	// 3. Read finish after client is resolved.
	if err := client.Resolve(ctx); err != nil {
		t.Error("client.Resolve:", err)
	}
	msg, release, err = p2.RecvMessage(ctx)
	if err != nil {
		t.Fatal("p2.RecvMessage:", err)
	}
	defer release()
	rmsg = rpcMessage{}
	if err := pogs.Extract(&rmsg, rpccapnp.Message_TypeID, msg.Struct); err != nil {
		t.Fatal("pogs.Extract(p2.RecvMessage(ctx)):", err)
	}
	if rmsg.Which != rpccapnp.Message_Which_finish {
		t.Fatalf("Received %v message; want finish", rmsg.Which)
	}
	if rmsg.Finish.QuestionID != qid {
		t.Errorf("Received finish for question %d; want %d", rmsg.Finish.QuestionID, qid)
	}
	if rmsg.Finish.ReleaseResultCaps {
		t.Error("Received finish that releases bootstrap result capabilities")
	}
}

type rpcMessage struct {
	Which         rpccapnp.Message_Which
	Unimplemented *rpcMessage
	Abort         *rpcException
	Bootstrap     *rpcBootstrap
	Return        *rpcReturn
	Finish        *rpcFinish
}

type rpcException struct {
	Reason string
	Type   rpccapnp.Exception_Type
}

type rpcBootstrap struct {
	QuestionID uint32 `capnp:"questionId"`
}

type rpcReturn struct {
	AnswerID         uint32 `capnp:"answerId"`
	ReleaseParamCaps bool

	Which                 rpccapnp.Return_Which
	Results               *rpcPayload
	Exception             *rpcException
	TakeFromOtherQuestion uint32
}

type rpcFinish struct {
	QuestionID        uint32 `capnp:"questionId"`
	ReleaseResultCaps bool
}

type rpcPayload struct {
	Content  capnp.Ptr
	CapTable []rpcCapDescriptor
}

type rpcCapDescriptor struct {
	Which          rpccapnp.CapDescriptor_Which
	SenderHosted   uint32
	SenderPromise  uint32
	ReceiverHosted uint32
}

type rpcPromisedAnswer struct {
	QuestionID uint32 `capnp:"questionId"`
	Transform  []rpcPromisedAnswerOp
}

type rpcPromisedAnswerOp struct {
	Which           rpccapnp.PromisedAnswer_Op_Which
	GetPointerField uint16
}
