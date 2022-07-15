package rpc_test

import (
	"context"
	"fmt"
	"testing"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/pogs"
	"capnproto.org/go/capnp/v3/rpc"
	"capnproto.org/go/capnp/v3/rpc/transport"
	"capnproto.org/go/capnp/v3/server"
	rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"
)

func TestSendDisembargo(t *testing.T) {
	t.Run("SendQueuedResultToCaller", func(t *testing.T) {
		testSendDisembargo(t, rpccp.Call_sendResultsTo_Which_caller)
	})
	t.Run("SendQueuedResultToYourself", func(t *testing.T) {
		t.Skip("TODO(soon): not implemented")
		testSendDisembargo(t, rpccp.Call_sendResultsTo_Which_yourself)
	})
}

// testSendDisembargo makes a call on the bootstrap capability with an
// export as a parameter, makes a pipelined call on the answer, writes a
// return with that export.  The Conn should send a disembargo, at which
// point a second call will be made on the answer.  The second call
// should not be delivered until after the disembargo loops back.
// Level 1 requirement.
func testSendDisembargo(t *testing.T, sendPrimeTo rpccp.Call_sendResultsTo_Which) {
	p1, p2 := transport.NewPipe(1)
	conn := rpc.NewConn(p1, &rpc.Options{
		ErrorReporter: testErrorReporter{tb: t},
	})
	defer finishTest(t, conn, p2)
	ctx := context.Background()

	// 1. Send bootstrap
	client := conn.Bootstrap(ctx)
	defer client.Release()
	var bootQID uint32
	{
		msg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if msg.Which != rpccp.Message_Which_bootstrap {
			t.Fatalf("Received %v message; want bootstrap", msg.Which)
		}
		bootQID = msg.Bootstrap.QuestionID
	}

	// 2. Return bootstrap
	{
		msg, send, release, err := p2.NewMessage(ctx)
		if err != nil {
			t.Fatal("p2.NewMessage():", err)
		}
		iptr := capnp.NewInterface(msg.Segment(), 0)
		err = pogs.Insert(rpccp.Message_TypeID, msg.Struct, &rpcMessage{
			Which: rpccp.Message_Which_return,
			Return: &rpcReturn{
				AnswerID: bootQID,
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
	}

	// 3. Read bootstrap finish
	{
		msg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if msg.Which != rpccp.Message_Which_finish {
			t.Fatalf("Received %v message; want finish", msg.Which)
		}
		if msg.Finish.QuestionID != bootQID {
			t.Errorf("Received finish for question %d; want %d", msg.Finish.QuestionID, bootQID)
		}
	}

	// 4. Make a call (call A) on bootstrap with export parameter
	callSeq := 1 // synchronized by server
	srv := newServer(func(ctx context.Context, call *server.Call) error {
		n := call.Args().Uint64(0)
		if int(n) != callSeq {
			err := fmt.Errorf("received call #%d on export when expecting call #%d", n, callSeq)
			t.Error("Server:", err)
			return err
		}
		callSeq++
		res, err := call.AllocResults(capnp.ObjectSize{DataSize: 8})
		if err != nil {
			return err
		}
		res.SetUint64(0, n)
		return nil
	}, nil)
	defer srv.Release()
	ctxA, cancelA := context.WithCancel(ctx)
	ansA, releaseCallA := client.SendCall(ctxA, capnp.Send{
		Method: capnp.Method{
			InterfaceID: interfaceID,
			MethodID:    methodID,
		},
		ArgsSize: capnp.ObjectSize{PointerCount: 1},
		PlaceArgs: func(s capnp.Struct) error {
			id := s.Message().AddCap(srv)
			ptr := capnp.NewInterface(s.Segment(), id).ToPtr()
			return s.SetPtr(0, ptr)
		},
	})
	defer func() {
		cancelA()
		releaseCallA()
	}()

	var qidA uint32
	var importID uint32
	{
		msg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if msg.Which != rpccp.Message_Which_call {
			t.Fatalf("Received %v message; want call", msg.Which)
		}
		qidA = msg.Call.QuestionID
		if msg.Call.InterfaceID != interfaceID {
			t.Errorf("call.interfaceId = %x; want %x", msg.Call.InterfaceID, interfaceID)
		}
		if msg.Call.MethodID != methodID {
			t.Errorf("call.methodId = %x; want %x", msg.Call.MethodID, methodID)
		}
		p, err := msg.Call.Params.Content.Struct().Ptr(0)
		if err != nil {
			t.Fatalf("call.params.content.ptr[0]: %v", err)
		}
		if !p.Interface().IsValid() {
			t.Fatal("call.params.content.ptr[0] is not an interface")
		}
		id := p.Interface().Capability()
		release()
		if int64(id) >= int64(len(msg.Call.Params.CapTable)) {
			t.Fatalf("call.params.content.ptr[0] refers to capability %d; table is size %d", id, len(msg.Call.Params.CapTable))
		}
		desc := msg.Call.Params.CapTable[id]
		if desc.Which != rpccp.CapDescriptor_Which_senderHosted {
			t.Fatalf("call.params.capTable[%d].Which = %v; want senderHosted", id, desc.Which)
		}
		importID = desc.SenderHosted
	}

	// 5. Make a pipelined call (call B) on answer
	sendCall := func(ctx context.Context, n uint64) (*capnp.Answer, capnp.ReleaseFunc) {
		transform := []capnp.PipelineOp{
			{Field: 0},
		}
		return ansA.PipelineSend(ctx, transform, capnp.Send{
			Method: capnp.Method{
				InterfaceID: interfaceID,
				MethodID:    methodID,
			},
			ArgsSize: capnp.ObjectSize{DataSize: 8},
			PlaceArgs: func(s capnp.Struct) error {
				s.SetUint64(0, n)
				return nil
			},
		})
	}
	ctxB, cancelB := context.WithCancel(ctx)
	ansB, releaseCallB := sendCall(ctxB, 1)
	defer func() {
		cancelB()
		releaseCallB()
	}()
	var qidB uint32
	var bParams capnp.Ptr
	{
		msg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if msg.Which != rpccp.Message_Which_call {
			t.Fatalf("Received %v message; want call", msg.Which)
		}
		qidB = msg.Call.QuestionID
		bParams = msg.Call.Params.Content
		if msg.Call.InterfaceID != interfaceID {
			t.Errorf("call.interfaceId = %x; want %x", msg.Call.InterfaceID, interfaceID)
		}
		if msg.Call.MethodID != methodID {
			t.Errorf("call.methodId = %x; want %x", msg.Call.MethodID, methodID)
		}
		if msg.Call.Target.Which != rpccp.MessageTarget_Which_promisedAnswer {
			t.Errorf("call.target is %v; want promisedAnswer", msg.Call.Target.Which)
		}
		if msg.Call.Target.PromisedAnswer.QuestionID != qidA {
			t.Errorf("call.target.promisedAnswer.questionID = %d; want %d (call A)", msg.Call.Target.PromisedAnswer.QuestionID, qidA)
		}
		if !msg.Call.Target.PromisedAnswer.transformEquals(0) {
			want := []rpcPromisedAnswerOp{
				{Which: rpccp.PromisedAnswer_Op_Which_getPointerField, GetPointerField: 0},
			}
			t.Errorf("call.target.promisedAnswer.transform = %+v; want %+v", msg.Call.Target.PromisedAnswer.Transform, want)
		}
	}

	// 6. Write return for call A with the import.
	{
		msg, send, release, err := p2.NewMessage(ctx)
		if err != nil {
			t.Fatal("p2.NewMessage():", err)
		}
		results, err := capnp.NewStruct(msg.Segment(), capnp.ObjectSize{PointerCount: 1})
		if err != nil {
			t.Fatal("capnp.NewStruct:", err)
		}
		if err := results.SetPtr(0, capnp.NewInterface(msg.Segment(), 0).ToPtr()); err != nil {
			t.Fatal("results.SetPtr:", err)
		}
		err = pogs.Insert(rpccp.Message_TypeID, msg.Struct, &rpcMessage{
			Which: rpccp.Message_Which_return,
			Return: &rpcReturn{
				AnswerID: bootQID,
				Which:    rpccp.Return_Which_results,
				Results: &rpcPayload{
					Content: results.ToPtr(),
					CapTable: []rpcCapDescriptor{
						{
							Which:        rpccp.CapDescriptor_Which_receiverHosted,
							SenderHosted: importID,
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
	}

	// 7. Read sender-loopback disembargo.
	var embargoID uint32
	{
		msg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if msg.Which != rpccp.Message_Which_disembargo {
			t.Fatalf("Received %v message; want disembargo", msg.Which)
		}
		if msg.Disembargo.Context.Which != rpccp.Disembargo_context_Which_senderLoopback {
			t.Fatalf("disembargo.context is %v; want senderLoopback", msg.Disembargo.Context.Which)
		}
		embargoID = msg.Disembargo.Context.SenderLoopback
		if msg.Disembargo.Target.Which != rpccp.MessageTarget_Which_promisedAnswer {
			t.Errorf("disembargo.target is %v; want promisedAnswer", msg.Disembargo.Target.Which)
		}
		if msg.Disembargo.Target.PromisedAnswer.QuestionID != qidA {
			t.Errorf("disembargo.target.promisedAnswer.questionId = %d; want %d (call A)", msg.Disembargo.Target.PromisedAnswer.QuestionID, qidA)
		}
		if !msg.Disembargo.Target.PromisedAnswer.transformEquals(0) {
			want := []rpcPromisedAnswerOp{
				{Which: rpccp.PromisedAnswer_Op_Which_getPointerField, GetPointerField: 0},
			}
			t.Errorf("disembargo.target.promisedAnswer.transform = %+v; want %+v", msg.Disembargo.Target.PromisedAnswer.Transform, want)
		}
	}

	// 8. Read call A finish.
	{
		msg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if msg.Which != rpccp.Message_Which_finish {
			t.Fatalf("Received %v message; want finish", msg.Which)
		}
		if msg.Finish.QuestionID != qidA {
			t.Errorf("finish.questionId = %d; want %d (call A)", msg.Finish.QuestionID, qidA)
		}
		if msg.Finish.ReleaseResultCaps {
			t.Error("finish.releaseResultCaps = true; want false")
		}
	}

	// 9. Make a pipelined call (call C) on answer.  Should not send on
	// wire, but might block until embargo is lifted..
	ctxC, cancelC := context.WithCancel(ctx)
	var (
		ansCReady    = make(chan struct{})
		ansC         *capnp.Answer
		releaseCallC capnp.ReleaseFunc
	)
	go func() {
		ansC, releaseCallC = sendCall(ctxC, 2)
		close(ansCReady)
	}()
	defer func() {
		cancelC()
		<-ansCReady
		releaseCallC()
	}()

	// 10. Echo call B back (call B', send results to yourself).
	const bPrimeAnswer = 909
	{
		err := sendMessage(ctx, p2, &rpcMessage{
			Which: rpccp.Message_Which_call,
			Call: &rpcCall{
				QuestionID:  bPrimeAnswer,
				InterfaceID: interfaceID,
				MethodID:    methodID,
				Target: rpcMessageTarget{
					Which:       rpccp.MessageTarget_Which_importedCap,
					ImportedCap: importID,
				},
				Params: rpcPayload{Content: bParams},
				SendResultsTo: rpcCallSendResultsTo{
					Which: sendPrimeTo,
				},
			},
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	// 11. Write receiver-loopback disembargo.
	// This is the earliest legal time to do so.
	{
		err := sendMessage(ctx, p2, &rpcMessage{
			Which: rpccp.Message_Which_disembargo,
			Disembargo: &rpcDisembargo{
				Context: rpcDisembargoContext{
					Which:            rpccp.Disembargo_context_Which_receiverLoopback,
					ReceiverLoopback: embargoID,
				},
				Target: rpcMessageTarget{
					Which:       rpccp.MessageTarget_Which_importedCap,
					ImportedCap: importID,
				},
			},
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	// 12. If sendPrimeTo == yourself, then write call B return (take from call B').
	if sendPrimeTo == rpccp.Call_sendResultsTo_Which_yourself {
		err := sendMessage(ctx, p2, &rpcMessage{
			Which: rpccp.Message_Which_return,
			Return: &rpcReturn{
				AnswerID:              qidB,
				Which:                 rpccp.Return_Which_takeFromOtherQuestion,
				TakeFromOtherQuestion: bPrimeAnswer,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	// 13. Read call B' return.  If sendPrimeTo == caller, then echo it back.
	{
		msg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if msg.Which != rpccp.Message_Which_return {
			t.Fatalf("Received %v message; want return", msg.Which)
		}
		if msg.Return.AnswerID != bPrimeAnswer {
			t.Errorf("return.answerId = %d; want %d (call B')", msg.Return.AnswerID, bPrimeAnswer)
		}
		switch sendPrimeTo {
		case rpccp.Call_sendResultsTo_Which_caller:
			if msg.Return.Which != rpccp.Return_Which_results {
				t.Fatalf("return is %v; want results", msg.Return.Which)
			}
			err := sendMessage(ctx, p2, &rpcMessage{
				Which: rpccp.Message_Which_return,
				Return: &rpcReturn{
					AnswerID: qidB,
					Which:    rpccp.Return_Which_results,
					Results: &rpcPayload{
						Content: msg.Return.Results.Content,
					},
				},
			})
			if err != nil {
				t.Fatal(err)
			}
		case rpccp.Call_sendResultsTo_Which_yourself:
			if msg.Return.Which != rpccp.Return_Which_resultsSentElsewhere {
				t.Errorf("return is %v; want resultsSentElsewhere", msg.Return.Which)
			}
		}
		release()
	}

	// 14. Read call B finish.  Must come after B' return, since otherwise
	// it would cancel B'.
	{
		msg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if msg.Which != rpccp.Message_Which_finish {
			t.Fatalf("Received %v message; want finish", msg.Which)
		}
		if msg.Finish.QuestionID != qidB {
			t.Errorf("finish.questionId = %d; want %d (call A)", msg.Finish.QuestionID, qidA)
		}
		if msg.Finish.ReleaseResultCaps {
			t.Error("finish.releaseResultCaps = true; want false")
		}
	}

	// 15. Write call B' finish.
	{
		err := sendMessage(ctx, p2, &rpcMessage{
			Which: rpccp.Message_Which_finish,
			Finish: &rpcFinish{
				QuestionID:        bPrimeAnswer,
				ReleaseResultCaps: false,
			},
		})
		if err != nil {
			t.Fatal(err)
		}
	}

	// 16. Check call sequence.
	b, err := ansB.Struct()
	if err != nil {
		t.Error("call B error:", err)
	} else if b.Uint64(0) != 1 {
		t.Errorf("call B result = %d; want 1", b.Uint64(0))
	}
	<-ansCReady
	c, err := ansC.Struct()
	if err != nil {
		t.Error("call C error:", err)
	} else if c.Uint64(0) != 2 {
		t.Errorf("call C result = %d; want 2", c.Uint64(0))
	}
}

// TestRecvDisembargo exposes a capability that echoes back received
// capabilities, writes a call to the conn's capability followed by two
// pipelined calls to the return value, reads the returns and sends a
// disembargo, verifies the call delivery order and the disembargo
// loopback.  Level 1 requirement.
func TestRecvDisembargo(t *testing.T) {
	srv := newServer(func(ctx context.Context, call *server.Call) error {
		in, err := call.Args().Ptr(0)
		if err != nil {
			return fmt.Errorf("capability arg: %v", err)
		}
		res, err := call.AllocResults(capnp.ObjectSize{PointerCount: 1})
		if err != nil {
			return err
		}
		if err := res.SetPtr(0, in.Interface().ToPtr()); err != nil {
			return fmt.Errorf("copy capability to result: %v", err)
		}
		return nil
	}, nil)
	p1, p2 := transport.NewPipe(2)
	conn := rpc.NewConn(p1, &rpc.Options{
		BootstrapClient: srv,
		ErrorReporter:   testErrorReporter{tb: t},
	})
	defer finishTest(t, conn, p2)
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

	// 2. Read bootstrap return.
	bootstrapImportID, err := recvBootstrapReturn(ctx, p2, bootstrapQID)
	if err != nil {
		t.Fatal(err)
	}

	// 3. Write bootstrap finish
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

	// 4. Write call with an export
	const callQID = 55
	const exportID = 888
	{
		msg, send, release, err := p2.NewMessage(ctx)
		if err != nil {
			t.Fatal("p2.NewMessage():", err)
		}
		params, err := capnp.NewStruct(msg.Segment(), capnp.ObjectSize{PointerCount: 1})
		if err != nil {
			release()
			t.Fatal("capnp.NewStruct:", err)
		}
		err = params.SetPtr(0, capnp.NewInterface(params.Segment(), 0).ToPtr())
		if err != nil {
			release()
			t.Fatal("capnp.NewStruct.SetPtr:", err)
		}
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
					CapTable: []rpcCapDescriptor{
						{
							Which:        rpccp.CapDescriptor_Which_senderHosted,
							SenderHosted: exportID,
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
	}

	// 5. Write pipelined call A.
	const callA_QID = 101
	{
		msg := &rpcMessage{
			Which: rpccp.Message_Which_call,
			Call: &rpcCall{
				QuestionID: callA_QID,
				Target: rpcMessageTarget{
					Which: rpccp.MessageTarget_Which_promisedAnswer,
					PromisedAnswer: &rpcPromisedAnswer{
						QuestionID: callQID,
						Transform: []rpcPromisedAnswerOp{
							{Which: rpccp.PromisedAnswer_Op_Which_getPointerField, GetPointerField: 0},
						},
					},
				},
				InterfaceID: 123,
				MethodID:    0,
			},
		}
		if err := sendMessage(ctx, p2, msg); err != nil {
			t.Fatal(err)
		}
	}

	// 6. Write pipelined call B.
	const callB_QID = 102
	{
		msg := &rpcMessage{
			Which: rpccp.Message_Which_call,
			Call: &rpcCall{
				QuestionID: callB_QID,
				Target: rpcMessageTarget{
					Which: rpccp.MessageTarget_Which_promisedAnswer,
					PromisedAnswer: &rpcPromisedAnswer{
						QuestionID: callQID,
						Transform: []rpcPromisedAnswerOp{
							{Which: rpccp.PromisedAnswer_Op_Which_getPointerField, GetPointerField: 0},
						},
					},
				},
				InterfaceID: 456,
				MethodID:    0,
			},
		}
		if err := sendMessage(ctx, p2, msg); err != nil {
			t.Fatal(err)
		}
	}

	// 7. Read calls and returns and disembargoes.
	//
	// Here's where things get tricky: Cap'n Proto does not guarantee
	// return order, just delivery order.
	//
	// At least at time of writing, the Go implementation delays sending
	// the first return until the pipelined calls are delivered, for
	// simplicity.  As far as I can tell, this is not a protocol
	// violation.

	calls := make(map[uint64]int)
	returns := make(map[uint32]int)
	disembargoTime := 0
	const embargoID = 9889
	isDone := func() bool {
		return calls[123] == 0 || calls[456] == 0 ||
			returns[callQID] == 0 || returns[callA_QID] == 0 || returns[callB_QID] == 0 ||
			disembargoTime == 0
	}
	for clock := 1; isDone(); clock++ {
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		switch rmsg.Which {
		case rpccp.Message_Which_return:
			if returns[rmsg.Return.AnswerID] != 0 {
				t.Errorf("multiple returns for %d", rmsg.Return.AnswerID)
				continue
			}
			returns[rmsg.Return.AnswerID] = clock
			if rmsg.Return.AnswerID == callQID {
				// Start disembargo
				msg := &rpcMessage{
					Which: rpccp.Message_Which_disembargo,
					Disembargo: &rpcDisembargo{
						Target: rpcMessageTarget{
							Which: rpccp.MessageTarget_Which_promisedAnswer,
							PromisedAnswer: &rpcPromisedAnswer{
								QuestionID: callQID,
								Transform: []rpcPromisedAnswerOp{
									{Which: rpccp.PromisedAnswer_Op_Which_getPointerField, GetPointerField: 0},
								},
							},
						},
						Context: rpcDisembargoContext{
							Which:          rpccp.Disembargo_context_Which_senderLoopback,
							SenderLoopback: embargoID,
						},
					},
				}
				if err := sendMessage(ctx, p2, msg); err != nil {
					t.Fatal(err)
				}
			}
			if rmsg.Return.Which == rpccp.Return_Which_exception {
				t.Fatalf("call %d returned exception: %s", rmsg.Return.AnswerID, rmsg.Return.Exception.Reason)
			}
			// Send finish
			{
				msg := &rpcMessage{
					Which: rpccp.Message_Which_finish,
					Finish: &rpcFinish{
						QuestionID:        rmsg.Return.AnswerID,
						ReleaseResultCaps: false,
					},
				}
				if err := sendMessage(ctx, p2, msg); err != nil {
					t.Fatal(err)
				}
			}
		case rpccp.Message_Which_call:
			if calls[rmsg.Call.InterfaceID] != 0 {
				t.Errorf("multiple calls for %d", rmsg.Call.InterfaceID)
				continue
			}
			calls[rmsg.Call.InterfaceID] = clock
			// Send return
			msg := &rpcMessage{
				Which: rpccp.Message_Which_return,
				Return: &rpcReturn{
					AnswerID:         rmsg.Call.QuestionID,
					ReleaseParamCaps: true,
					Which:            rpccp.Return_Which_results,
					Results:          &rpcPayload{},
				},
			}
			if err := sendMessage(ctx, p2, msg); err != nil {
				t.Fatal(err)
			}
		case rpccp.Message_Which_disembargo:
			if rmsg.Disembargo.Context.Which != rpccp.Disembargo_context_Which_receiverLoopback {
				t.Errorf("received disembargo in %v context; want receiverLoopback", rmsg.Disembargo.Context.Which)
				continue
			}
			if rmsg.Disembargo.Context.ReceiverLoopback != embargoID {
				t.Errorf("received disembargo for ID %d; want %d", rmsg.Disembargo.Context.ReceiverLoopback, embargoID)
				continue
			}
			if disembargoTime != 0 {
				t.Error("encountered multiple disembargoes for same ID")
				continue
			}
			disembargoTime = clock
			if rmsg.Disembargo.Target.Which != rpccp.MessageTarget_Which_importedCap {
				t.Errorf("received disembargo target = %v; want importedCap", rmsg.Disembargo.Target.Which)
				continue
			}
			if rmsg.Disembargo.Target.ImportedCap != exportID {
				t.Errorf("received disembargo target capability = %d; want %d", rmsg.Disembargo.Target.ImportedCap, exportID)
			}
		case rpccp.Message_Which_finish, rpccp.Message_Which_release:
			// Ignore. Not relevant to test.
		default:
			t.Errorf("don't know how to handle %v; skipping", rmsg.Which)
		}
	}

	if calls[callA_QID] > calls[callB_QID] {
		t.Error("call A delivered after call B")
	}
	if disembargoTime < calls[callA_QID] || disembargoTime < calls[callB_QID] {
		t.Error("disembargo delivered before calls")
	}
}

// TestIssue3 exposes a capability that makes a call to its received
// capability argument, acks the call, then waits on its return.  In
// earlier versions of go-capnproto, this would cause a deadlock.
// See https://github.com/capnproto/go-capnproto2/issues/3 for history.
// Level 1 requirement.
func TestIssue3(t *testing.T) {
	const callbackInterfaceID = interfaceID + 100
	const callbackMethodID = 9
	srv := newServer(func(ctx context.Context, call *server.Call) error {
		in, err := call.Args().Ptr(0)
		if err != nil {
			return fmt.Errorf("capability arg: %v", err)
		}
		ans, release := in.Interface().Client().SendCall(ctx, capnp.Send{
			Method: capnp.Method{
				InterfaceID: callbackInterfaceID,
				MethodID:    callbackMethodID,
			},
		})
		defer release()
		call.Ack()
		callbackResult, err := ans.Struct()
		if err != nil {
			return fmt.Errorf("callback: %v", err)
		}
		res, err := call.AllocResults(callbackResult.Size())
		if err != nil {
			return err
		}
		if err := res.CopyFrom(callbackResult); err != nil {
			return fmt.Errorf("copy to result: %v", err)
		}
		return nil
	}, nil)
	p1, p2 := transport.NewPipe(1)
	conn := rpc.NewConn(p1, &rpc.Options{
		BootstrapClient: srv,
		ErrorReporter:   testErrorReporter{tb: t},
	})
	defer finishTest(t, conn, p2)
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

	// 2. Read bootstrap return.
	bootstrapImportID, err := recvBootstrapReturn(ctx, p2, bootstrapQID)
	if err != nil {
		t.Fatal(err)
	}

	// 3. Write bootstrap finish
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

	// 4. Write call with an export
	const callQID = 55
	const exportID = 888
	{
		msg, send, release, err := p2.NewMessage(ctx)
		if err != nil {
			t.Fatal("p2.NewMessage():", err)
		}
		params, err := capnp.NewStruct(msg.Segment(), capnp.ObjectSize{PointerCount: 1})
		if err != nil {
			release()
			t.Fatal("capnp.NewStruct:", err)
		}
		err = params.SetPtr(0, capnp.NewInterface(params.Segment(), 0).ToPtr())
		if err != nil {
			release()
			t.Fatal("capnp.NewStruct.SetPtr:", err)
		}
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
					CapTable: []rpcCapDescriptor{
						{
							Which:        rpccp.CapDescriptor_Which_senderHosted,
							SenderHosted: exportID,
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
	}

	// 5. Read callback
	var callbackQID uint32
	{
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if err != nil {
			t.Fatal("pogs.Extract(p2.RecvMessage(ctx)):", err)
		}
		if rmsg.Which != rpccp.Message_Which_call {
			t.Fatalf("Received %v message; want call", rmsg.Which)
		}
		callbackQID = rmsg.Call.QuestionID
		if rmsg.Call.InterfaceID != callbackInterfaceID || rmsg.Call.MethodID != callbackMethodID {
			got := capnp.Method{InterfaceID: rmsg.Call.InterfaceID, MethodID: rmsg.Call.MethodID}
			want := capnp.Method{InterfaceID: callbackInterfaceID, MethodID: callbackMethodID}
			t.Fatalf("Received call for %v; want %v", got, want)
		}
		if rmsg.Call.Target.Which != rpccp.MessageTarget_Which_importedCap {
			t.Fatalf("Received call on %v; want importedCap", rmsg.Call.Target.Which)
		}
		if rmsg.Call.Target.ImportedCap != exportID {
			t.Fatalf("Received call on export %d; want %d", rmsg.Call.Target.ImportedCap, exportID)
		}
		release()
	}

	// 6. Write callback return.
	{
		msg, send, release, err := p2.NewMessage(ctx)
		if err != nil {
			t.Fatal("p2.NewMessage():", err)
		}
		results, err := capnp.NewStruct(msg.Segment(), capnp.ObjectSize{DataSize: 8})
		if err != nil {
			release()
			t.Fatal("capnp.NewStruct:", err)
		}
		results.SetUint64(0, 0xdeadbeef)
		err = pogs.Insert(rpccp.Message_TypeID, msg.Struct, &rpcMessage{
			Which: rpccp.Message_Which_return,
			Return: &rpcReturn{
				AnswerID:         callbackQID,
				ReleaseParamCaps: true,
				Which:            rpccp.Return_Which_results,
				Results:          &rpcPayload{Content: results.ToPtr()},
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

	// 7. Read callback finish and call return (and possibly a release).
	seenFinish := false
	seenReturn := false
	exportRefs := 1
	for !seenFinish || !seenReturn {
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		switch rmsg.Which {
		case rpccp.Message_Which_return:
			if seenReturn {
				t.Fatal("Received two returns")
			}
			seenReturn = true
			if rmsg.Return.AnswerID != callQID {
				t.Errorf("Received return for question %d; want %d", rmsg.Return.AnswerID, callQID)
			}
			switch rmsg.Return.Which {
			case rpccp.Return_Which_results:
				content := rmsg.Return.Results.Content.Struct()
				if content.Uint64(0) != 0xdeadbeef {
					t.Errorf("Call result = %#x; want 0xdeadbeef", content.Uint64(0))
				}
			case rpccp.Return_Which_exception:
				t.Errorf("Call exception: %s", rmsg.Return.Exception.Reason)
			default:
				t.Fatalf("Received %v return; want results", rmsg.Return.Which)
			}
		case rpccp.Message_Which_finish:
			if seenFinish {
				t.Fatal("Received two finishes")
			}
			seenFinish = true
			if rmsg.Finish.QuestionID != callbackQID {
				t.Errorf("Received finish for question %d; want %d", rmsg.Finish.QuestionID, callbackQID)
			}
		case rpccp.Message_Which_release:
			if rmsg.Release.ID != exportID {
				t.Errorf("Received release for export ID %d; want %d", rmsg.Release.ID, exportID)
			} else if rmsg.Release.ReferenceCount > uint32(exportRefs) {
				t.Errorf("Released %d references for export, but only have %d", rmsg.Release.ReferenceCount, exportRefs)
				exportRefs = 0
			} else {
				exportRefs -= int(rmsg.Release.ReferenceCount)
			}
		default:
			t.Fatalf("Received %v message; wanted return, finish, or release", rmsg.Which)
		}
		release()
	}

	// 8. Write finish
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
}

type rpcResolve struct {
	PromiseID uint32 `capnp:"promiseId"`
	Which     rpccp.Resolve_Which
	Cap       *rpcCapDescriptor
	Exception *rpcException
}

type rpcRelease struct {
	ID             uint32 `capnp:"id"`
	ReferenceCount uint32
}

type rpcDisembargo struct {
	Target  rpcMessageTarget
	Context rpcDisembargoContext
}

type rpcDisembargoContext struct {
	Which            rpccp.Disembargo_context_Which
	SenderLoopback   uint32
	ReceiverLoopback uint32
	Provide          uint32
}
