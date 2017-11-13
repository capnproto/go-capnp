package rpc_test

import (
	"context"
	"testing"

	"zombiezen.com/go/capnproto2"
	"zombiezen.com/go/capnproto2/pogs"
	"zombiezen.com/go/capnproto2/rpc"
	rpccp "zombiezen.com/go/capnproto2/std/capnp/rpc"
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
	if err := pogs.Extract(&rmsg, rpccp.Message_TypeID, msg.Struct); err != nil {
		t.Fatal("pogs.Extract(p2.RecvMessage(ctx)):", err)
	}
	if rmsg.Which != rpccp.Message_Which_abort {
		t.Fatalf("Received %v message; want abort", rmsg.Which)
	}
	if rmsg.Abort == nil {
		t.Error("Received null abort message")
	} else if rmsg.Abort.Type != rpccp.Exception_Type_failed {
		t.Errorf("Received exception type %v; want failed", rmsg.Abort.Type)
	}
}

func TestBootstrapCall(t *testing.T) {
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
	if err := pogs.Extract(&rmsg, rpccp.Message_TypeID, msg.Struct); err != nil {
		t.Fatal("pogs.Extract(p2.RecvMessage(ctx)):", err)
	}
	if rmsg.Which != rpccp.Message_Which_bootstrap {
		t.Fatalf("Received %v message; want bootstrap", rmsg.Which)
	}
	qid := rmsg.Bootstrap.QuestionID

	// 2. Write back a return
	msg, send, release, err := p2.NewMessage(ctx)
	if err != nil {
		t.Fatal("p2.NewMessage():", err)
	}
	iptr := capnp.NewInterface(msg.Segment(), 0)
	err = pogs.Insert(rpccp.Message_TypeID, msg.Struct, &rpcMessage{
		Which: rpccp.Message_Which_return,
		Return: &rpcReturn{
			AnswerID: qid,
			Which:    rpccp.Return_Which_results,
			Results: &rpcPayload{
				Content: iptr.ToPtr(),
				CapTable: []rpcCapDescriptor{
					{
						Which:        rpccp.CapDescriptor_Which_senderHosted,
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
	if err := pogs.Extract(&rmsg, rpccp.Message_TypeID, msg.Struct); err != nil {
		t.Fatal("pogs.Extract(p2.RecvMessage(ctx)):", err)
	}
	if rmsg.Which != rpccp.Message_Which_finish {
		t.Fatalf("Received %v message; want finish", rmsg.Which)
	}
	if rmsg.Finish.QuestionID != qid {
		t.Errorf("Received finish for question %d; want %d", rmsg.Finish.QuestionID, qid)
	}
	if rmsg.Finish.ReleaseResultCaps {
		t.Error("Received finish that releases bootstrap result capabilities")
	}

	// 4. Make a call.
	ans, release := client.SendCall(ctx, capnp.Send{
		Method: capnp.Method{
			InterfaceID: interfaceID,
			MethodID:    methodID,
		},
		ArgsSize: capnp.ObjectSize{DataSize: 8},
		PlaceArgs: func(s capnp.Struct) error {
			s.SetUint64(0, 42)
			return nil
		},
	})
	defer release()
	msg, release, err = p2.RecvMessage(ctx)
	if err != nil {
		t.Fatal("p2.RecvMessage:", err)
	}
	defer release()
	rmsg = rpcMessage{}
	if err := pogs.Extract(&rmsg, rpccp.Message_TypeID, msg.Struct); err != nil {
		t.Fatal("pogs.Extract(p2.RecvMessage(ctx)):", err)
	}
	if rmsg.Which != rpccp.Message_Which_call {
		t.Fatalf("Received %v message; want call", rmsg.Which)
	}
	qid = rmsg.Call.QuestionID
	if rmsg.Call.InterfaceID != interfaceID {
		t.Errorf("call.interfaceId = %x; want %x", rmsg.Call.InterfaceID, interfaceID)
	}
	if rmsg.Call.MethodID != methodID {
		t.Errorf("call.methodId = %x; want %x", rmsg.Call.MethodID, methodID)
	}
	if p := rmsg.Call.Params.Content.Struct(); p.Uint64(0) != 42 {
		t.Errorf("call.params.content = %d; want 42", p.Uint64(0))
	}
	if rmsg.Call.SendResultsTo.Which != rpccp.Call_sendResultsTo_Which_caller {
		t.Errorf("call.sentResultsTo which = %v; want caller", rmsg.Call.SendResultsTo.Which)
	}
	if rmsg.Call.Target.Which != rpccp.MessageTarget_Which_importedCap {
		t.Errorf("call.target which = %v; want importedCap", rmsg.Call.Target.Which)
	} else if rmsg.Call.Target.ImportedCap != bootstrapExportID {
		t.Errorf("call.target.importedCap = %d; want %d", rmsg.Call.Target.ImportedCap, bootstrapExportID)
	}

	// 5. Return a response.
	msg, send, release, err = p2.NewMessage(ctx)
	if err != nil {
		t.Fatal("p2.NewMessage():", err)
	}
	resp, err := capnp.NewStruct(msg.Segment(), capnp.ObjectSize{DataSize: 8})
	if err != nil {
		t.Fatal("capnp.NewStruct:", err)
	}
	resp.SetUint64(0, 0xdeadbeef)
	err = pogs.Insert(rpccp.Message_TypeID, msg.Struct, &rpcMessage{
		Which: rpccp.Message_Which_return,
		Return: &rpcReturn{
			AnswerID: qid,
			Which:    rpccp.Return_Which_results,
			Results:  &rpcPayload{Content: resp.ToPtr()},
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

	// 6. Read result from answer.
	resp, err = ans.Struct()
	if err != nil {
		t.Error("ans.Struct():", err)
	} else if resp.Uint64(0) != 0xdeadbeef {
		t.Errorf("ans.Struct().Uint64(0) = %#x; want 0xdeadbeef", resp.Uint64(0))
	}

	// 7. Read the finish
	msg, release, err = p2.RecvMessage(ctx)
	if err != nil {
		t.Fatal("p2.RecvMessage:", err)
	}
	defer release()
	rmsg = rpcMessage{}
	if err := pogs.Extract(&rmsg, rpccp.Message_TypeID, msg.Struct); err != nil {
		t.Fatal("pogs.Extract(p2.RecvMessage(ctx)):", err)
	}
	if rmsg.Which != rpccp.Message_Which_finish {
		t.Fatalf("Received %v message; want finish", rmsg.Which)
	}
	if rmsg.Finish.QuestionID != qid {
		t.Errorf("Received finish for question %d; want %d", rmsg.Finish.QuestionID, qid)
	}
	if rmsg.Finish.ReleaseResultCaps {
		t.Error("Received finish that releases call result capabilities")
	}

	// 8. Release the client
	client.Release()

	// 9. Read the release.
	msg, release, err = p2.RecvMessage(ctx)
	if err != nil {
		t.Fatal("p2.RecvMessage:", err)
	}
	defer release()
	rmsg = rpcMessage{}
	if err := pogs.Extract(&rmsg, rpccp.Message_TypeID, msg.Struct); err != nil {
		t.Fatal("pogs.Extract(p2.RecvMessage(ctx)):", err)
	}
	if rmsg.Which != rpccp.Message_Which_release {
		t.Fatalf("Received %v message; want release", rmsg.Which)
	}
	if rmsg.Release.ID != bootstrapExportID {
		t.Errorf("Received release for import %d; want %d", rmsg.Release.ID, bootstrapExportID)
	}
	if rmsg.Release.ReferenceCount != 1 {
		t.Errorf("Received release for %d references; want 1", rmsg.Release.ReferenceCount)
	}
}

func TestBootstrapPipelineCall(t *testing.T) {
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
	if err := pogs.Extract(&rmsg, rpccp.Message_TypeID, msg.Struct); err != nil {
		t.Fatal("pogs.Extract(p2.RecvMessage(ctx)):", err)
	}
	if rmsg.Which != rpccp.Message_Which_bootstrap {
		t.Fatalf("Received %v message; want bootstrap", rmsg.Which)
	}
	bootstrapQID := rmsg.Bootstrap.QuestionID

	// 2. Make a call.
	ans, release := client.SendCall(ctx, capnp.Send{
		Method: capnp.Method{
			InterfaceID: interfaceID,
			MethodID:    methodID,
		},
		ArgsSize: capnp.ObjectSize{DataSize: 8},
		PlaceArgs: func(s capnp.Struct) error {
			s.SetUint64(0, 42)
			return nil
		},
	})
	defer release()
	msg, release, err = p2.RecvMessage(ctx)
	if err != nil {
		t.Fatal("p2.RecvMessage:", err)
	}
	defer release()
	rmsg = rpcMessage{}
	if err := pogs.Extract(&rmsg, rpccp.Message_TypeID, msg.Struct); err != nil {
		t.Fatal("pogs.Extract(p2.RecvMessage(ctx)):", err)
	}
	if rmsg.Which != rpccp.Message_Which_call {
		t.Fatalf("Received %v message; want call", rmsg.Which)
	}
	qid := rmsg.Call.QuestionID
	if rmsg.Call.InterfaceID != interfaceID {
		t.Errorf("call.interfaceId = %x; want %x", rmsg.Call.InterfaceID, interfaceID)
	}
	if rmsg.Call.MethodID != methodID {
		t.Errorf("call.methodId = %x; want %x", rmsg.Call.MethodID, methodID)
	}
	if p := rmsg.Call.Params.Content.Struct(); p.Uint64(0) != 42 {
		t.Errorf("call.params.content = %d; want 42", p.Uint64(0))
	}
	if rmsg.Call.SendResultsTo.Which != rpccp.Call_sendResultsTo_Which_caller {
		t.Errorf("call.sentResultsTo which = %v; want caller", rmsg.Call.SendResultsTo.Which)
	}
	if rmsg.Call.Target.Which != rpccp.MessageTarget_Which_promisedAnswer {
		t.Errorf("call.target which = %v; want promisedAnswer", rmsg.Call.Target.Which)
	} else {
		if rmsg.Call.Target.PromisedAnswer.QuestionID != bootstrapQID {
			t.Errorf("call.target.promisedAnswer.questionID = %d; want %d", rmsg.Call.Target.PromisedAnswer.QuestionID, bootstrapQID)
		}
		if xform := rmsg.Call.Target.PromisedAnswer.Transform; len(xform) != 0 {
			t.Errorf("call.target.promisedAnswer.transform = %v; want []", xform)
		}
	}

	// 3. Return a response.
	msg, send, release, err := p2.NewMessage(ctx)
	if err != nil {
		t.Fatal("p2.NewMessage():", err)
	}
	resp, err := capnp.NewStruct(msg.Segment(), capnp.ObjectSize{DataSize: 8})
	if err != nil {
		t.Fatal("capnp.NewStruct:", err)
	}
	resp.SetUint64(0, 0xdeadbeef)
	err = pogs.Insert(rpccp.Message_TypeID, msg.Struct, &rpcMessage{
		Which: rpccp.Message_Which_return,
		Return: &rpcReturn{
			AnswerID: qid,
			Which:    rpccp.Return_Which_results,
			Results:  &rpcPayload{Content: resp.ToPtr()},
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

	// 4. Read result from answer.
	resp, err = ans.Struct()
	if err != nil {
		t.Error("ans.Struct():", err)
	} else if resp.Uint64(0) != 0xdeadbeef {
		t.Errorf("ans.Struct().Uint64(0) = %#x; want 0xdeadbeef", resp.Uint64(0))
	}

	// 5. Read the finish
	msg, release, err = p2.RecvMessage(ctx)
	if err != nil {
		t.Fatal("p2.RecvMessage:", err)
	}
	defer release()
	rmsg = rpcMessage{}
	if err := pogs.Extract(&rmsg, rpccp.Message_TypeID, msg.Struct); err != nil {
		t.Fatal("pogs.Extract(p2.RecvMessage(ctx)):", err)
	}
	if rmsg.Which != rpccp.Message_Which_finish {
		t.Fatalf("Received %v message; want finish", rmsg.Which)
	}
	if rmsg.Finish.QuestionID != qid {
		t.Errorf("Received finish for question %d; want %d", rmsg.Finish.QuestionID, qid)
	}
	if rmsg.Finish.ReleaseResultCaps {
		t.Error("Received finish that releases call result capabilities")
	}

	// 6. Release the client
	client.Release()

	// 7. Read the finish.
	msg, release, err = p2.RecvMessage(ctx)
	if err != nil {
		t.Fatal("p2.RecvMessage:", err)
	}
	defer release()
	rmsg = rpcMessage{}
	if err := pogs.Extract(&rmsg, rpccp.Message_TypeID, msg.Struct); err != nil {
		t.Fatal("pogs.Extract(p2.RecvMessage(ctx)):", err)
	}
	if rmsg.Which != rpccp.Message_Which_finish {
		t.Fatalf("Received %v message; want finish", rmsg.Which)
	}
	if rmsg.Finish.QuestionID != bootstrapQID {
		t.Errorf("Received finish for question %d; want %d", rmsg.Finish.QuestionID, bootstrapQID)
	}
	if !rmsg.Finish.ReleaseResultCaps {
		t.Error("Received finish that does not release bootstrap")
	}
}

type rpcMessage struct {
	Which         rpccp.Message_Which
	Unimplemented *rpcMessage
	Abort         *rpcException
	Bootstrap     *rpcBootstrap
	Call          *rpcCall
	Return        *rpcReturn
	Finish        *rpcFinish
	Release       *rpcRelease
}

type rpcException struct {
	Reason string
	Type   rpccp.Exception_Type
}

type rpcBootstrap struct {
	QuestionID uint32 `capnp:"questionId"`
}

type rpcCall struct {
	QuestionID              uint32 `capnp:"questionId"`
	Target                  rpcMessageTarget
	InterfaceID             uint64 `capnp:"interfaceId"`
	MethodID                uint16 `capnp:"methodId"`
	AllowThirdPartyTailCall bool
	Params                  rpcPayload
	SendResultsTo           rpcCallSendResultsTo
}

type rpcCallSendResultsTo struct {
	Which rpccp.Call_sendResultsTo_Which
}

type rpcReturn struct {
	AnswerID         uint32 `capnp:"answerId"`
	ReleaseParamCaps bool

	Which                 rpccp.Return_Which
	Results               *rpcPayload
	Exception             *rpcException
	TakeFromOtherQuestion uint32
}

type rpcFinish struct {
	QuestionID        uint32 `capnp:"questionId"`
	ReleaseResultCaps bool
}

type rpcRelease struct {
	ID             uint32 `capnp:"id"`
	ReferenceCount uint32
}

type rpcMessageTarget struct {
	Which          rpccp.MessageTarget_Which
	ImportedCap    uint32
	PromisedAnswer *rpcPromisedAnswer
}

type rpcPayload struct {
	Content  capnp.Ptr
	CapTable []rpcCapDescriptor
}

type rpcCapDescriptor struct {
	Which          rpccp.CapDescriptor_Which
	SenderHosted   uint32
	SenderPromise  uint32
	ReceiverHosted uint32
}

type rpcPromisedAnswer struct {
	QuestionID uint32 `capnp:"questionId"`
	Transform  []rpcPromisedAnswerOp
}

type rpcPromisedAnswerOp struct {
	Which           rpccp.PromisedAnswer_Op_Which
	GetPointerField uint16
}
