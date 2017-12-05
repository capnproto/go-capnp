package rpc_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"zombiezen.com/go/capnproto2"
	"zombiezen.com/go/capnproto2/pogs"
	"zombiezen.com/go/capnproto2/rpc"
	"zombiezen.com/go/capnproto2/server"
	rpccp "zombiezen.com/go/capnproto2/std/capnp/rpc"
)

const (
	interfaceID       uint64 = 0xa7317bd7216570aa
	methodID          uint16 = 9
	bootstrapExportID uint32 = 84
)

// TestCloseAbort calls Close on a new connection, verifying that it
// sends an Abort message.  Level 0 requirement.
func TestCloseAbort(t *testing.T) {
	p1, p2 := newPipe(1)
	defer p2.CloseSend()
	defer p2.CloseRecv()
	conn := rpc.NewConn(p1, &rpc.Options{
		ErrorReporter: testErrorReporter{t},
	})

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

// TestBootstrapCall calls Bootstrap, sends back an export, then makes
// an RPC on the returned capability.  It checks to see that the correct
// messages were sent on the wire and that the correct return value came
// back.  Level 0 requirement.
func TestBootstrapCall(t *testing.T) {
	p1, p2 := newPipe(1)
	defer p2.CloseSend()
	defer p2.CloseRecv()
	conn := rpc.NewConn(p1, &rpc.Options{
		ErrorReporter: testErrorReporter{t},
	})
	defer func() {
		if err := conn.Close(); err != nil {
			t.Error(err)
		}
	}()

	ctx := context.Background()

	// 1. Read bootstrap
	client := conn.Bootstrap(ctx)
	if err := client.Resolve(canceledContext(ctx)); err == nil {
		t.Error("bootstrap client reports resolved before return")
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

// TestBootstrapPipelineCall calls Bootstrap and makes an RPC on the
// returned capability without resolving the client.  Level 0 requirement.
func TestBootstrapPipelineCall(t *testing.T) {
	p1, p2 := newPipe(1)
	defer p2.CloseSend()
	defer p2.CloseRecv()
	conn := rpc.NewConn(p1, &rpc.Options{
		ErrorReporter: testErrorReporter{t},
	})
	defer func() {
		if err := conn.Close(); err != nil {
			t.Error(err)
		}
	}()

	ctx := context.Background()

	// 1. Read bootstrap
	client := conn.Bootstrap(ctx)
	if err := client.Resolve(canceledContext(ctx)); err == nil {
		t.Error("bootstrap client reports resolved before return")
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

// TestBootstrapClient sets Options.BootstrapClient on NewConn,
// bootstraps, waits for a return, then sends a call to the RPC
// connection.  It checks that the correct messages were sent and that
// the return values are correct.  Level 0 requirement.
func TestBootstrapClient(t *testing.T) {
	srvShutdown := make(chan struct{})
	srv := capnp.NewClient(server.New(
		[]server.Method{
			{
				Method: capnp.Method{
					InterfaceID: interfaceID,
					MethodID:    methodID,
				},
				Impl: func(ctx context.Context, call *server.Call) error {
					resp, err := call.AllocResults(capnp.ObjectSize{DataSize: 8})
					if err != nil {
						return err
					}
					resp.SetUint64(0, 0xdeadbeef|uint64(call.Args().Uint32(0))<<32)
					return nil
				},
			},
		},
		nil, /* brand */
		shutdownFunc(func() { close(srvShutdown) }),
		nil /* policy */))
	p1, p2 := newPipe(1)
	defer p2.CloseSend()
	defer p2.CloseRecv()
	conn := rpc.NewConn(p1, &rpc.Options{
		BootstrapClient: srv,
		ErrorReporter:   testErrorReporter{t},
	})
	defer func() {
		if err := conn.Close(); err != nil {
			t.Error(err)
		}
		select {
		case <-srvShutdown:
		default:
			t.Error("Bootstrap client still alive after Close returned")
		}
	}()

	ctx := context.Background()

	// 1. Write bootstrap
	const bootstrapQID = 54
	{
		msg := &rpcMessage{
			Which:     rpccp.Message_Which_bootstrap,
			Bootstrap: &rpcBootstrap{QuestionID: bootstrapQID},
		}
		if err := sendMessage(ctx, p2, msg); err != nil {
			t.Fatal(err)
		}
	}

	// 2. Read return
	var bootstrapImportID uint32
	{
		msg, release, err := p2.RecvMessage(ctx)
		if err != nil {
			t.Fatal("p2.RecvMessage:", err)
		}
		defer release()
		var rmsg rpcMessage
		if err := pogs.Extract(&rmsg, rpccp.Message_TypeID, msg.Struct); err != nil {
			t.Fatal("pogs.Extract(p2.RecvMessage(ctx)):", err)
		}
		if rmsg.Which != rpccp.Message_Which_return {
			t.Fatalf("Received %v message; want return", rmsg.Which)
		}
		if rmsg.Return.AnswerID != bootstrapQID {
			t.Errorf("Received return for answer %d; want %d", rmsg.Return.AnswerID, bootstrapQID)
		}
		if rmsg.Return.Which != rpccp.Return_Which_results {
			t.Fatalf("return which = %v; want results", rmsg.Return.Which)
		}
		desc, err := payloadCapability(rmsg.Return.Results)
		if err != nil {
			t.Fatal(err)
		}
		if desc.Which != rpccp.CapDescriptor_Which_senderHosted {
			t.Fatalf("Received %v capability for bootstrap; want senderHosted", desc.Which)
		}
		if len(rmsg.Return.Results.CapTable) > 1 {
			t.Errorf("Received bootstrap return with %d capability descriptors; want 1", len(rmsg.Return.Results.CapTable))
		}
		bootstrapImportID = desc.SenderHosted
	}

	// 3. Write finish
	{
		msg := &rpcMessage{
			Which: rpccp.Message_Which_finish,
			Finish: &rpcFinish{
				QuestionID:        bootstrapQID,
				ReleaseResultCaps: false,
			},
		}
		if err := sendMessage(ctx, p2, msg); err != nil {
			t.Fatal(err)
		}
	}

	// 4. Write call
	const callQID = 55
	{
		msg, send, release, err := p2.NewMessage(ctx)
		if err != nil {
			t.Fatal("p2.NewMessage():", err)
		}
		params, err := capnp.NewStruct(msg.Segment(), capnp.ObjectSize{DataSize: 8})
		if err != nil {
			t.Fatal("capnp.NewStruct:", err)
		}
		params.SetUint32(0, 0x2a2b)
		err = pogs.Insert(rpccp.Message_TypeID, msg.Struct, &rpcMessage{
			Which: rpccp.Message_Which_call,
			Call: &rpcCall{
				QuestionID: callQID,
				Target: rpcMessageTarget{
					Which:       rpccp.MessageTarget_Which_importedCap,
					ImportedCap: bootstrapImportID,
				},
				InterfaceID: interfaceID,
				MethodID:    methodID,
				Params: rpcPayload{
					Content: params.ToPtr(),
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
	}

	// 5. Read return
	{
		msg, release, err := p2.RecvMessage(ctx)
		if err != nil {
			t.Fatal("p2.RecvMessage:", err)
		}
		defer release()
		var rmsg rpcMessage
		if err := pogs.Extract(&rmsg, rpccp.Message_TypeID, msg.Struct); err != nil {
			t.Fatal("pogs.Extract(p2.RecvMessage(ctx)):", err)
		}
		if rmsg.Which != rpccp.Message_Which_return {
			t.Fatalf("Received %v message; want return", rmsg.Which)
		}
		if rmsg.Return.AnswerID != callQID {
			t.Errorf("Received return for answer %d; want %d", rmsg.Return.AnswerID, callQID)
		}
		if rmsg.Return.Which != rpccp.Return_Which_results {
			t.Fatalf("return which = %v; want results", rmsg.Return.Which)
		}
		result := rmsg.Return.Results.Content.Struct()
		if got, want := result.Uint64(0), uint64(0x00002a2bdeadbeef); got != want {
			t.Errorf("return results content = %#016x; want %#016x", got, want)
		}
	}

	// 6. Write finish
	{
		msg := &rpcMessage{
			Which: rpccp.Message_Which_finish,
			Finish: &rpcFinish{
				QuestionID:        callQID,
				ReleaseResultCaps: false,
			},
		}
		if err := sendMessage(ctx, p2, msg); err != nil {
			t.Fatal(err)
		}
	}

	// 7. Write release
	{
		msg := &rpcMessage{
			Which: rpccp.Message_Which_release,
			Release: &rpcRelease{
				ID:             bootstrapImportID,
				ReferenceCount: 1,
			},
		}
		if err := sendMessage(ctx, p2, msg); err != nil {
			t.Fatal(err)
		}
	}

	// 8. srv should not be released (Conn should be holding on).
	select {
	case <-srvShutdown:
		t.Error("Bootstrap client released after release message")
	default:
	}
}

// TestPromisedBootstrapAnswerCall sets Options.BootstrapClient on
// NewConn, bootstraps, waits for a return, then sends a call on the
// promised answer without sending a finish.  It checks that the correct
// messages were sent and that the return values are correct.  Level 0
// requirement.
func TestPromisedBootstrapAnswerCall(t *testing.T) {
	srvShutdown := make(chan struct{})
	srv := capnp.NewClient(server.New(
		[]server.Method{
			{
				Method: capnp.Method{
					InterfaceID: interfaceID,
					MethodID:    methodID,
				},
				Impl: func(ctx context.Context, call *server.Call) error {
					resp, err := call.AllocResults(capnp.ObjectSize{DataSize: 8})
					if err != nil {
						return err
					}
					resp.SetUint64(0, 0xdeadbeef|uint64(call.Args().Uint32(0))<<32)
					return nil
				},
			},
		},
		nil, /* brand */
		shutdownFunc(func() { close(srvShutdown) }),
		nil /* policy */))
	p1, p2 := newPipe(1)
	defer p2.CloseSend()
	defer p2.CloseRecv()
	conn := rpc.NewConn(p1, &rpc.Options{
		BootstrapClient: srv,
		ErrorReporter:   testErrorReporter{t},
	})
	defer func() {
		if err := conn.Close(); err != nil {
			t.Error(err)
		}
		select {
		case <-srvShutdown:
		default:
			t.Error("Bootstrap client still alive after Close returned")
		}
	}()
	ctx := context.Background()

	// 1. Write bootstrap
	const bootstrapQID = 54
	{
		msg := &rpcMessage{
			Which:     rpccp.Message_Which_bootstrap,
			Bootstrap: &rpcBootstrap{QuestionID: bootstrapQID},
		}
		if err := sendMessage(ctx, p2, msg); err != nil {
			t.Fatal(err)
		}
	}

	// 2. Read return
	{
		msg, release, err := p2.RecvMessage(ctx)
		if err != nil {
			t.Fatal("p2.RecvMessage:", err)
		}
		defer release()
		var rmsg rpcMessage
		if err := pogs.Extract(&rmsg, rpccp.Message_TypeID, msg.Struct); err != nil {
			t.Fatal("pogs.Extract(p2.RecvMessage(ctx)):", err)
		}
		if rmsg.Which != rpccp.Message_Which_return {
			t.Fatalf("Received %v message; want return", rmsg.Which)
		}
		if rmsg.Return.AnswerID != bootstrapQID {
			t.Errorf("Received return for answer %d; want %d", rmsg.Return.AnswerID, bootstrapQID)
		}
		if rmsg.Return.Which != rpccp.Return_Which_results {
			t.Fatalf("return which = %v; want results", rmsg.Return.Which)
		}
		desc, err := payloadCapability(rmsg.Return.Results)
		if err != nil {
			t.Fatal(err)
		}
		if desc.Which != rpccp.CapDescriptor_Which_senderHosted {
			t.Fatalf("Received %v capability for bootstrap; want senderHosted", desc.Which)
		}
		if len(rmsg.Return.Results.CapTable) > 1 {
			t.Errorf("Received bootstrap return with %d capability descriptors; want 1", len(rmsg.Return.Results.CapTable))
		}
	}

	// 3. Write call
	const callQID = 55
	{
		msg, send, release, err := p2.NewMessage(ctx)
		if err != nil {
			t.Fatal("p2.NewMessage():", err)
		}
		params, err := capnp.NewStruct(msg.Segment(), capnp.ObjectSize{DataSize: 8})
		if err != nil {
			t.Fatal("capnp.NewStruct:", err)
		}
		params.SetUint32(0, 0x2a2b)
		err = pogs.Insert(rpccp.Message_TypeID, msg.Struct, &rpcMessage{
			Which: rpccp.Message_Which_call,
			Call: &rpcCall{
				QuestionID: callQID,
				Target: rpcMessageTarget{
					Which: rpccp.MessageTarget_Which_promisedAnswer,
					PromisedAnswer: &rpcPromisedAnswer{
						QuestionID: bootstrapQID,
						Transform:  nil,
					},
				},
				InterfaceID: interfaceID,
				MethodID:    methodID,
				Params: rpcPayload{
					Content: params.ToPtr(),
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
	}

	// 4. Read return
	{
		msg, release, err := p2.RecvMessage(ctx)
		if err != nil {
			t.Fatal("p2.RecvMessage:", err)
		}
		defer release()
		var rmsg rpcMessage
		if err := pogs.Extract(&rmsg, rpccp.Message_TypeID, msg.Struct); err != nil {
			t.Fatal("pogs.Extract(p2.RecvMessage(ctx)):", err)
		}
		if rmsg.Which != rpccp.Message_Which_return {
			t.Fatalf("Received %v message; want return", rmsg.Which)
		}
		if rmsg.Return.AnswerID != callQID {
			t.Errorf("Received return for answer %d; want %d", rmsg.Return.AnswerID, callQID)
		}
		if rmsg.Return.Which != rpccp.Return_Which_results {
			if rmsg.Return.Which == rpccp.Return_Which_exception {
				t.Logf("returned exception = %q", rmsg.Return.Exception.Reason)
			}
			t.Fatalf("return which = %v; want results", rmsg.Return.Which)
		}
		result := rmsg.Return.Results.Content.Struct()
		if got, want := result.Uint64(0), uint64(0x00002a2bdeadbeef); got != want {
			t.Errorf("return results content = %#016x; want %#016x", got, want)
		}
	}
}

func TestCallOnClosedConn(t *testing.T) {
	p1, p2 := newPipe(1)
	defer p2.CloseSend()
	defer p2.CloseRecv()
	conn := rpc.NewConn(p1, &rpc.Options{
		ErrorReporter: testErrorReporter{t},
	})
	closed := false
	defer func() {
		if closed {
			return
		}
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

	// 4. Close the Conn.
	err = conn.Close()
	closed = true
	if err != nil {
		t.Error(err)
	}

	// 5. Make a call.
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
	_, err = ans.Struct()
	if !capnp.IsDisconnected(err) {
		t.Errorf("call after Close returned error: %v; want disconnected", err)
	}
}

type shutdownFunc func()

func (f shutdownFunc) Shutdown() { f() }

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

func sendMessage(ctx context.Context, t rpc.Transport, msg *rpcMessage) error {
	s, send, release, err := t.NewMessage(ctx)
	if err != nil {
		return fmt.Errorf("send message: %v", err)
	}
	defer release()
	if err := pogs.Insert(rpccp.Message_TypeID, s.Struct, msg); err != nil {
		return fmt.Errorf("send message: %v", err)
	}
	if err := send(); err != nil {
		return fmt.Errorf("send message: %v", err)
	}
	return nil
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
	ReceiverAnswer *rpcPromisedAnswer
}

// payloadCapability returns the capability descriptor pointed to by a
// payload's content.
func payloadCapability(payload *rpcPayload) (*rpcCapDescriptor, error) {
	iface := payload.Content.Interface()
	if !iface.IsValid() {
		return nil, errors.New("parse payload: content is not an interface pointer")
	}
	if int64(iface.Capability()) >= int64(len(payload.CapTable)) {
		return nil, fmt.Errorf("parse payload: content points to capability %d (table has %d entries)", iface.Capability(), len(payload.CapTable))
	}
	return &payload.CapTable[iface.Capability()], nil
}

type rpcPromisedAnswer struct {
	QuestionID uint32 `capnp:"questionId"`
	Transform  []rpcPromisedAnswerOp
}

type rpcPromisedAnswerOp struct {
	Which           rpccp.PromisedAnswer_Op_Which
	GetPointerField uint16
}

func canceledContext(parent context.Context) context.Context {
	ctx, cancel := context.WithCancel(parent)
	cancel()
	return ctx
}

type testErrorReporter struct {
	t *testing.T
}

func (r testErrorReporter) ReportError(e error) {
	r.t.Log("conn error:", e)
}
