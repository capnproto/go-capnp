package rpc_test

import (
	"context"
	"fmt"
	"testing"

	"zombiezen.com/go/capnproto2"
	"zombiezen.com/go/capnproto2/pogs"
	"zombiezen.com/go/capnproto2/rpc"
	"zombiezen.com/go/capnproto2/server"
	rpccp "zombiezen.com/go/capnproto2/std/capnp/rpc"
)

// TestRecvDisembargo exposes a capability that echoes back received
// capabilities, writes a call to the conn's capability followed by two
// pipelined calls to the return value, reads the returns and sends a
// disembargo, verifies the call delivery order and the disembargo
// loopback.  Level 1 requirement.
func TestRecvDisembargo(t *testing.T) {
	srv := capnp.NewClient(server.New([]server.Method{
		{
			Method: capnp.Method{
				InterfaceID: interfaceID,
				MethodID:    methodID,
			},
			Impl: func(ctx context.Context, call *server.Call) error {
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
			},
		},
	}, nil /* brand */, nil /* shutdown */, nil /* policy */))
	p1, p2 := newPipe(2)
	defer p2.Close()
	conn := rpc.NewConn(p1, &rpc.Options{
		BootstrapClient: srv,
		ErrorReporter:   testErrorReporter{tb: t},
	})
	defer func() {
		if err := conn.Close(); err != nil {
			t.Error(err)
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
		msg, release, err := p2.RecvMessage(ctx)
		if err != nil {
			t.Fatal("p2.RecvMessage:", err)
		}
		var rmsg rpcMessage
		err = pogs.Extract(&rmsg, rpccp.Message_TypeID, msg.Struct)
		release()
		if err != nil {
			t.Fatal("pogs.Extract(p2.RecvMessage(ctx)):", err)
		}
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
				t.Errorf("received disembargo for ID %d; want %d", rmsg.Disembargo.Context.ReceiverLoopback)
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
