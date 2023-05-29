package rpc_test

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/exc"
	"capnproto.org/go/capnp/v3/pogs"
	"capnproto.org/go/capnp/v3/rpc"
	testcp "capnproto.org/go/capnp/v3/rpc/internal/testcapnp"
	"capnproto.org/go/capnp/v3/rpc/transport"
	"capnproto.org/go/capnp/v3/schemas"
	"capnproto.org/go/capnp/v3/server"
	rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func init() {
	testcp.RegisterSchema(schemas.DefaultRegistry)
	rpccp.RegisterSchema(schemas.DefaultRegistry)
}

const (
	interfaceID       uint64 = 0xa7317bd7216570aa
	methodID          uint16 = 9
	bootstrapExportID uint32 = 84
)

func TestMain(m *testing.M) {
	var (
		mu     sync.Mutex
		leaked bool
	)
	capnp.SetClientLeakFunc(func(msg string) {
		mu.Lock()
		leaked = true
		fmt.Fprintln(os.Stderr, "LEAK:", msg)
		mu.Unlock()
	})
	status := m.Run()
	runtime.GC() // try to trigger any finalizers
	mu.Lock()
	if status == 0 && leaked {
		os.Exit(1)
	}
	os.Exit(status)
}

// TestSendAbort calls Close on a new connection, verifying that it
// sends an Abort message and it reports no errors.  Level 0 requirement.
func TestSendAbort(t *testing.T) {
	t.Parallel()
	t.Helper()

	t.Run("ReceiverListening", func(t *testing.T) {
		t.Parallel()

		left, right := net.Pipe()
		p1, p2 := transport.NewStream(left), transport.NewStream(right)
		defer p2.Close()

		conn := rpc.NewConn(p1, &rpc.Options{
			ErrorReporter: testErrorReporter{tb: t, fail: true},
			// Give it plenty of time to actually send the message;
			// otherwise we might time out and close the connection first.
			// "plenty of time" here really means defer to the test suite's
			// timeout.
			AbortTimeout: time.Duration(math.MaxInt64),
		})

		ctx := context.Background()
		select {
		case <-conn.Done():
			t.Error("conn.Done closed before Close")
		default:
		}

		go func() {
			if err := conn.Close(); err != nil {
				t.Error("conn.Close():", err)
			}
			select {
			case <-conn.Done():
			default:
				t.Error("conn.Done open after Close")
			}
		}()

		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if rmsg.Which != rpccp.Message_Which_abort {
			t.Fatalf("Received %v message; want abort", rmsg.Which)
		}
		if rmsg.Abort == nil {
			t.Error("Received null abort message")
		} else if rmsg.Abort.Type != rpccp.Exception_Type_failed {
			t.Errorf("Received exception type %v; want failed", rmsg.Abort.Type)
		}
	})
	t.Run("ReceiverNotListening", func(t *testing.T) {
		t.Parallel()

		p1, p2 := net.Pipe()
		defer p2.Close()
		conn := rpc.NewConn(transport.NewStream(p1), &rpc.Options{
			ErrorReporter: testErrorReporter{tb: t, fail: true},
		})

		// Should have a timeout.
		if err := conn.Close(); err != nil {
			t.Error("conn.Close():", err)
		}
	})
}

// TestRecvAbort writes an abort message to a connection, waits for
// bootstrap resolution/disconnect (to acknowledge delivery), and then
// closes the connection, verifying that Close does not return an error.
// Level 0 requirement.
func TestRecvAbort(t *testing.T) {
	t.Parallel()

	left, right := transport.NewPipe(1)
	p1, p2 := rpc.NewTransport(left), rpc.NewTransport(right)
	defer p2.Close()

	conn := rpc.NewConn(p1, &rpc.Options{
		ErrorReporter: testErrorReporter{tb: t},
	})

	select {
	case <-conn.Done():
		t.Error("conn.Done closed before receiving abort")
	default:
	}
	err := sendMessage(context.Background(), p2, &rpcMessage{
		Which: rpccp.Message_Which_abort,
		Abort: &rpcException{
			Type:   rpccp.Exception_Type_failed,
			Reason: "over it",
		},
	})
	require.NoError(t, err, "must send 'failed' exception")

	ctx := context.Background()
	boot := conn.Bootstrap(ctx)
	defer boot.Release()

	err = boot.Resolve(ctx)
	require.NoError(t, err, "should resolve bootstrap capability")

	ans, releaseCall := boot.SendCall(context.Background(), capnp.Send{
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
	_, err = ans.Struct()
	releaseCall()
	require.True(t, exc.IsType(err, exc.Disconnected), "should be 'disconnected' exception")

	boot.Release()
	<-conn.Done()

	err = conn.Close()
	assert.NoError(t, err, "should close without error")
}

// TestSendBootstrapError calls Bootstrap, raises an exception, then
// makes an RPC on the client.  It checks to see that the RPC returns an
// error with the correct message.  Level 0 requirement.
func TestSendBootstrapError(t *testing.T) {
	t.Parallel()

	left, right := transport.NewPipe(1)
	p1, p2 := rpc.NewTransport(left), rpc.NewTransport(right)

	conn := rpc.NewConn(p1, &rpc.Options{
		ErrorReporter: testErrorReporter{tb: t},
	})
	defer finishTest(t, conn, p2)

	ctx := context.Background()

	// 1. Read bootstrap
	client := conn.Bootstrap(ctx)
	defer client.Release()
	var qid uint32
	{
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if rmsg.Which != rpccp.Message_Which_bootstrap {
			t.Fatalf("Received %v message; want bootstrap", rmsg.Which)
		}
		qid = rmsg.Bootstrap.QuestionID
	}

	// 2. Raise an exception.
	{
		msg := &rpcMessage{
			Which: rpccp.Message_Which_return,
			Return: &rpcReturn{
				AnswerID: qid,
				Which:    rpccp.Return_Which_exception,
				Exception: &rpcException{
					Type:   rpccp.Exception_Type_failed,
					Reason: "everything went wrong",
				},
			},
		}
		if err := sendMessage(ctx, p2, msg); err != nil {
			t.Fatal(err)
		}
	}

	// 3. Read finish after client is resolved.
	{
		if err := client.Resolve(ctx); err != nil {
			t.Error("client.Resolve:", err)
		}
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if rmsg.Which != rpccp.Message_Which_finish {
			t.Fatalf("Received %v message; want finish", rmsg.Which)
		}
		if rmsg.Finish.QuestionID != qid {
			t.Errorf("Received finish for question %d; want %d", rmsg.Finish.QuestionID, qid)
		}
		if rmsg.Finish.ReleaseResultCaps {
			t.Error("Received finish that releases bootstrap result capabilities")
		}
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
	_, err := ans.Struct()
	const want = "everything went wrong"
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Errorf("call on bootstrap error = %v; want to contain %q", err, want)
	}
}

// TestSendBootstrapCall calls Bootstrap, sends back an export, then
// makes an RPC on the returned capability.  It checks to see that the
// correct messages were sent on the wire and that the correct return
// value came back.  Level 0 requirement.
func TestSendBootstrapCall(t *testing.T) {
	t.Parallel()

	left, right := transport.NewPipe(1)
	p1, p2 := rpc.NewTransport(left), rpc.NewTransport(right)

	conn := rpc.NewConn(p1, &rpc.Options{
		ErrorReporter: testErrorReporter{tb: t},
	})
	defer finishTest(t, conn, p2)

	ctx := context.Background()

	// 1. Read bootstrap
	client := conn.Bootstrap(ctx)
	defer client.Release()
	if err := client.Resolve(canceledContext(ctx)); err == nil {
		t.Error("bootstrap client reports resolved before return")
	}
	var qid uint32
	{
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if rmsg.Which != rpccp.Message_Which_bootstrap {
			t.Fatalf("Received %v message; want bootstrap", rmsg.Which)
		}
		qid = rmsg.Bootstrap.QuestionID
	}

	// 2. Write back a return
	{
		outMsg, err := p2.NewMessage()
		if err != nil {
			t.Fatal("p2.NewMessage():", err)
		}
		iptr := capnp.NewInterface(outMsg.Message().Segment(), 0)
		err = pogs.Insert(rpccp.Message_TypeID, capnp.Struct(outMsg.Message()), &rpcMessage{
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
			outMsg.Release()
			t.Fatal("pogs.Insert(p2.NewMessage(), &rpcMessage{...}):", err)
		}
		err = outMsg.Send()
		outMsg.Release()
		if err != nil {
			t.Fatal("send():", err)
		}
	}

	// 3. Read finish after client is resolved.
	{
		if err := client.Resolve(ctx); err != nil {
			t.Error("client.Resolve:", err)
		}
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if rmsg.Which != rpccp.Message_Which_finish {
			t.Fatalf("Received %v message; want finish", rmsg.Which)
		}
		if rmsg.Finish.QuestionID != qid {
			t.Errorf("Received finish for question %d; want %d", rmsg.Finish.QuestionID, qid)
		}
		if rmsg.Finish.ReleaseResultCaps {
			t.Error("Received finish that releases bootstrap result capabilities")
		}
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
	{
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
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
	}

	// 5. Return a response.
	{
		outMsg, err := p2.NewMessage()
		if err != nil {
			t.Fatal("p2.NewMessage():", err)
		}
		resp, err := capnp.NewStruct(outMsg.Message().Segment(), capnp.ObjectSize{DataSize: 8})
		if err != nil {
			t.Fatal("capnp.NewStruct:", err)
		}
		resp.SetUint64(0, 0xdeadbeef)
		err = pogs.Insert(rpccp.Message_TypeID, capnp.Struct(outMsg.Message()), &rpcMessage{
			Which: rpccp.Message_Which_return,
			Return: &rpcReturn{
				AnswerID: qid,
				Which:    rpccp.Return_Which_results,
				Results:  &rpcPayload{Content: resp.ToPtr()},
			},
		})
		if err != nil {
			outMsg.Release()
			t.Fatal("pogs.Insert(p2.NewMessage(), &rpcMessage{...}):", err)
		}
		err = outMsg.Send()
		outMsg.Release()
		if err != nil {
			t.Fatal("send():", err)
		}
	}

	// 6. Read result from answer.
	{
		resp, err := ans.Struct()
		if err != nil {
			t.Error("ans.Struct():", err)
		} else if resp.Uint64(0) != 0xdeadbeef {
			t.Errorf("ans.Struct().Uint64(0) = %#x; want 0xdeadbeef", resp.Uint64(0))
		}
	}

	// 7. Read the finish
	{
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if rmsg.Which != rpccp.Message_Which_finish {
			t.Fatalf("Received %v message; want finish", rmsg.Which)
		}
		if rmsg.Finish.QuestionID != qid {
			t.Errorf("Received finish for question %d; want %d", rmsg.Finish.QuestionID, qid)
		}
		if rmsg.Finish.ReleaseResultCaps {
			t.Error("Received finish that releases call result capabilities")
		}
	}

	// 8. Release the client
	client.Release()
	{
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
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
}

// TestSendBootstrapCallException calls Bootstrap, sends back an export,
// then makes an RPC on the returned capability, which will raise an
// exception.  It checks to see that the correct messages were sent on
// the wire and that the error message is correct.  Level 0 requirement.
func TestSendBootstrapCallException(t *testing.T) {
	t.Parallel()

	left, right := transport.NewPipe(1)
	p1, p2 := rpc.NewTransport(left), rpc.NewTransport(right)

	conn := rpc.NewConn(p1, &rpc.Options{
		ErrorReporter: testErrorReporter{tb: t},
	})
	defer finishTest(t, conn, p2)

	ctx := context.Background()

	// 1. Read bootstrap
	client := conn.Bootstrap(ctx)
	defer client.Release()
	if err := client.Resolve(canceledContext(ctx)); err == nil {
		t.Error("bootstrap client reports resolved before return")
	}
	var qid uint32
	{
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if rmsg.Which != rpccp.Message_Which_bootstrap {
			t.Fatalf("Received %v message; want bootstrap", rmsg.Which)
		}
		qid = rmsg.Bootstrap.QuestionID
	}

	// 2. Write back a return
	{
		outMsg, err := p2.NewMessage()
		if err != nil {
			t.Fatal("p2.NewMessage():", err)
		}
		iptr := capnp.NewInterface(outMsg.Message().Segment(), 0)
		err = pogs.Insert(rpccp.Message_TypeID, capnp.Struct(outMsg.Message()), &rpcMessage{
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
			outMsg.Release()
			t.Fatal("pogs.Insert(p2.NewMessage(), &rpcMessage{...}):", err)
		}
		err = outMsg.Send()
		outMsg.Release()
		if err != nil {
			t.Fatal("send():", err)
		}
	}

	// 3. Read finish after client is resolved.
	{
		if err := client.Resolve(ctx); err != nil {
			t.Error("client.Resolve:", err)
		}
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if rmsg.Which != rpccp.Message_Which_finish {
			t.Fatalf("Received %v message; want finish", rmsg.Which)
		}
		if rmsg.Finish.QuestionID != qid {
			t.Errorf("Received finish for question %d; want %d", rmsg.Finish.QuestionID, qid)
		}
		if rmsg.Finish.ReleaseResultCaps {
			t.Error("Received finish that releases bootstrap result capabilities")
		}
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
	{
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
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
	}

	// 5. Raise an exception.
	{
		msg := &rpcMessage{
			Which: rpccp.Message_Which_return,
			Return: &rpcReturn{
				AnswerID: qid,
				Which:    rpccp.Return_Which_exception,
				Exception: &rpcException{
					Type:   rpccp.Exception_Type_failed,
					Reason: "everything went wrong",
				},
			},
		}
		if err := sendMessage(ctx, p2, msg); err != nil {
			t.Fatal(err)
		}
	}

	// 6. Read result from answer.
	{
		_, err := ans.Struct()
		const want = "everything went wrong"
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("ans.Struct() = _, %v; want error to contain %q", err, want)
		}
	}

	// 7. Read the finish
	{
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if rmsg.Which != rpccp.Message_Which_finish {
			t.Fatalf("Received %v message; want finish", rmsg.Which)
		}
		if rmsg.Finish.QuestionID != qid {
			t.Errorf("Received finish for question %d; want %d", rmsg.Finish.QuestionID, qid)
		}
		if rmsg.Finish.ReleaseResultCaps {
			t.Error("Received finish that releases call result capabilities")
		}
	}
}

// TestSendBootstrapPipelineCall calls Bootstrap and makes an RPC on the
// returned capability without resolving the client.  Level 0 requirement.
func TestSendBootstrapPipelineCall(t *testing.T) {
	t.Parallel()

	left, right := transport.NewPipe(1)
	p1, p2 := rpc.NewTransport(left), rpc.NewTransport(right)

	conn := rpc.NewConn(p1, &rpc.Options{
		ErrorReporter: testErrorReporter{tb: t},
	})
	defer finishTest(t, conn, p2)

	ctx := context.Background()

	// 1. Read bootstrap
	client := conn.Bootstrap(ctx)
	defer client.Release()
	if err := client.Resolve(canceledContext(ctx)); err == nil {
		t.Error("bootstrap client reports resolved before return")
	}
	var bootstrapQID uint32
	{
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if rmsg.Which != rpccp.Message_Which_bootstrap {
			t.Fatalf("Received %v message; want bootstrap", rmsg.Which)
		}
		bootstrapQID = rmsg.Bootstrap.QuestionID
	}

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
	var qid uint32
	{
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
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
		if rmsg.Call.Target.Which != rpccp.MessageTarget_Which_promisedAnswer {
			t.Errorf("call.target which = %v; want promisedAnswer", rmsg.Call.Target.Which)
		} else {
			if rmsg.Call.Target.PromisedAnswer.QuestionID != bootstrapQID {
				t.Errorf("call.target.promisedAnswer.questionID = %d; want %d", rmsg.Call.Target.PromisedAnswer.QuestionID, bootstrapQID)
			}
			if !rmsg.Call.Target.PromisedAnswer.transformEquals() {
				t.Errorf("call.target.promisedAnswer.transform = %v; want []", rmsg.Call.Target.PromisedAnswer.Transform)
			}
		}
	}

	// 3. Return a response.
	{
		outMsg, err := p2.NewMessage()
		if err != nil {
			t.Fatal("p2.NewMessage():", err)
		}
		resp, err := capnp.NewStruct(outMsg.Message().Segment(), capnp.ObjectSize{DataSize: 8})
		if err != nil {
			t.Fatal("capnp.NewStruct:", err)
		}
		resp.SetUint64(0, 0xdeadbeef)
		err = pogs.Insert(rpccp.Message_TypeID, capnp.Struct(outMsg.Message()), &rpcMessage{
			Which: rpccp.Message_Which_return,
			Return: &rpcReturn{
				AnswerID: qid,
				Which:    rpccp.Return_Which_results,
				Results:  &rpcPayload{Content: resp.ToPtr()},
			},
		})
		if err != nil {
			outMsg.Release()
			t.Fatal("pogs.Insert(p2.NewMessage(), &rpcMessage{...}):", err)
		}
		err = outMsg.Send()
		outMsg.Release()
		if err != nil {
			t.Fatal("send():", err)
		}
	}

	// 4. Read result from answer.
	{
		resp, err := ans.Struct()
		if err != nil {
			t.Error("ans.Struct():", err)
		} else if resp.Uint64(0) != 0xdeadbeef {
			t.Errorf("ans.Struct().Uint64(0) = %#x; want 0xdeadbeef", resp.Uint64(0))
		}
	}

	// 5. Read the finish
	{
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if rmsg.Which != rpccp.Message_Which_finish {
			t.Fatalf("Received %v message; want finish", rmsg.Which)
		}
		if rmsg.Finish.QuestionID != qid {
			t.Errorf("Received finish for question %d; want %d", rmsg.Finish.QuestionID, qid)
		}
		if rmsg.Finish.ReleaseResultCaps {
			t.Error("Received finish that releases call result capabilities")
		}
	}

	// 6. Send back a return for the bootstrap message:
	bootstrapExportID := uint32(99)
	{
		outMsg, err := p2.NewMessage()
		require.NoError(t, err)
		iface := capnp.NewInterface(outMsg.Message().Segment(), 0)
		require.NoError(t, pogs.Insert(rpccp.Message_TypeID, capnp.Struct(outMsg.Message()),
			&rpcMessage{
				Which: rpccp.Message_Which_return,
				Return: &rpcReturn{
					AnswerID: bootstrapQID,
					Which:    rpccp.Return_Which_results,
					Results: &rpcPayload{
						Content: iface.ToPtr(),
						CapTable: []rpcCapDescriptor{
							{
								Which:        rpccp.CapDescriptor_Which_senderHosted,
								SenderHosted: bootstrapExportID,
							},
						},
					},
				},
			},
		))
		require.NoError(t, outMsg.Send())
		outMsg.Release()
	}

	// 7. Read the finish:
	{
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if rmsg.Which != rpccp.Message_Which_finish {
			t.Fatalf("Received %v message; want finish", rmsg.Which)
		}
		if rmsg.Finish.QuestionID != bootstrapQID {
			t.Errorf("Received finish for question %d; want %d", rmsg.Finish.QuestionID, bootstrapQID)
		}
		require.False(
			t,
			rmsg.Finish.ReleaseResultCaps,
			"Received finish that releases bootstrap (should receive separate releasemessage)",
		)
	}

	// 8. Release the client, read the release message.
	client.Release()
	{
		rmsg, release, err := recvMessage(ctx, p2)
		require.NoError(t, err)
		defer release()
		require.Equal(t, rpccp.Message_Which_release, rmsg.Which)
		require.Equal(t, bootstrapExportID, rmsg.Release.ID)
		require.Equal(t, uint32(1), rmsg.Release.ReferenceCount)
	}
}

// TestRecvBootstrapError does not set Options.BootstrapClient and
// receives a Bootstrap message.  It checks that an exception was sent
// back.  Level 0 requirement.
func TestRecvBootstrapError(t *testing.T) {
	t.Parallel()

	left, right := transport.NewPipe(1)
	p1, p2 := rpc.NewTransport(left), rpc.NewTransport(right)

	conn := rpc.NewConn(p1, &rpc.Options{
		ErrorReporter: testErrorReporter{tb: t},
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
	defer func() {
		// 3. Write finish
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
	}()

	// 2. Read return
	rmsg, release, err := recvMessage(ctx, p2)
	if err != nil {
		t.Fatal(err)
	}
	defer release()
	if rmsg.Which != rpccp.Message_Which_return {
		t.Fatalf("Received %v message; want return", rmsg.Which)
	}
	if rmsg.Return.AnswerID != bootstrapQID {
		t.Errorf("return.answerId = %d; want %d", rmsg.Return.AnswerID, bootstrapQID)
	}
	if rmsg.Return.Which != rpccp.Return_Which_exception {
		t.Fatalf("return is %v; want exception", rmsg.Return.Which)
	}
	if rmsg.Return.Exception.Type != rpccp.Exception_Type_failed {
		// Exception type is not a strict requirement, but this seems the
		// most appropriate.
		t.Errorf("return.exception.type = %v; want failed", rmsg.Return.Exception.Type)
	}
}

// TestRecvBootstrapCall sets Options.BootstrapClient on NewConn,
// bootstraps, waits for a return, then sends a call to the RPC
// connection.  It checks that the correct messages were sent and that
// the return values are correct.  Level 0 requirement.
func TestRecvBootstrapCall(t *testing.T) {
	t.Parallel()

	srvShutdown := make(chan struct{})
	srv := newServer(
		func(ctx context.Context, call *server.Call) error {
			resp, err := call.AllocResults(capnp.ObjectSize{DataSize: 8})
			if err != nil {
				return err
			}
			resp.SetUint64(0, 0xdeadbeef|uint64(call.Args().Uint32(0))<<32)
			return nil
		},
		func() {
			close(srvShutdown)
		})
	left, right := transport.NewPipe(1)
	p1, p2 := rpc.NewTransport(left), rpc.NewTransport(right)

	conn := rpc.NewConn(p1, &rpc.Options{
		BootstrapClient: srv,
		ErrorReporter:   testErrorReporter{tb: t},
	})
	defer func() {
		finishTest(t, conn, p2)
		<-srvShutdown // Hangs if bootstrap client is never shut down.
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
	bootstrapImportID, err := recvBootstrapReturn(ctx, p2, bootstrapQID)
	if err != nil {
		t.Fatal(err)
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
		outMsg, err := p2.NewMessage()
		if err != nil {
			t.Fatal("p2.NewMessage():", err)
		}
		params, err := capnp.NewStruct(outMsg.Message().Segment(), capnp.ObjectSize{DataSize: 8})
		if err != nil {
			t.Fatal("capnp.NewStruct:", err)
		}
		params.SetUint32(0, 0x2a2b)
		err = pogs.Insert(rpccp.Message_TypeID, capnp.Struct(outMsg.Message()), &rpcMessage{
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
			outMsg.Release()
			t.Fatal("pogs.Insert(p2.NewMessage(), &rpcMessage{...}):", err)
		}
		err = outMsg.Send()
		outMsg.Release()
		if err != nil {
			t.Fatal("send():", err)
		}
	}

	// 5. Read return
	{
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
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

// TestRecvBootstrapCallException sets Options.BootstrapClient on
// NewConn, bootstraps, waits for a return, then sends a call to the RPC
// connection that will return an error.  It checks that the correct
// messages were sent and that the return values are correct.  Level 0
// requirement.
func TestRecvBootstrapCallException(t *testing.T) {
	t.Parallel()

	srv := newServer(func(ctx context.Context, call *server.Call) error {
		return errors.New("everything went wrong")
	}, nil)
	left, right := transport.NewPipe(1)
	p1, p2 := rpc.NewTransport(left), rpc.NewTransport(right)

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

	// 2. Read return
	var bootstrapImportID uint32
	bootstrapImportID, err := recvBootstrapReturn(ctx, p2, bootstrapQID)
	if err != nil {
		t.Fatal(err)
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
		outMsg, err := p2.NewMessage()
		if err != nil {
			t.Fatal("p2.NewMessage():", err)
		}
		params, err := capnp.NewStruct(outMsg.Message().Segment(), capnp.ObjectSize{DataSize: 8})
		if err != nil {
			t.Fatal("capnp.NewStruct:", err)
		}
		params.SetUint32(0, 0x2a2b)
		err = pogs.Insert(rpccp.Message_TypeID, capnp.Struct(outMsg.Message()), &rpcMessage{
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
			outMsg.Release()
			t.Fatal("pogs.Insert(p2.NewMessage(), &rpcMessage{...}):", err)
		}
		err = outMsg.Send()
		outMsg.Release()
		if err != nil {
			t.Fatal("send():", err)
		}
	}

	// 5. Read return
	{
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if rmsg.Which != rpccp.Message_Which_return {
			t.Fatalf("Received %v message; want return", rmsg.Which)
		}
		if rmsg.Return.AnswerID != callQID {
			t.Errorf("Received return for answer %d; want %d", rmsg.Return.AnswerID, callQID)
		}
		if rmsg.Return.Which != rpccp.Return_Which_exception {
			t.Fatalf("return which = %v; want results", rmsg.Return.Which)
		}
		if rmsg.Return.Exception.Type != rpccp.Exception_Type_failed {
			t.Errorf("return.exception.type = %v; want failed", rmsg.Return.Exception.Type)
		}
		const want = "everything went wrong"
		if !strings.Contains(rmsg.Return.Exception.Reason, want) {
			t.Errorf("return.exception.reason = %q; want to contain %q", rmsg.Return.Exception.Reason, want)
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
}

// TestRecvBootstrapPipelineCall sets Options.BootstrapClient on
// NewConn, bootstraps, waits for a return, then sends a call on the
// promised answer without sending a finish.  It checks that the correct
// messages were sent and that the return values are correct.  Level 0
// requirement.
func TestRecvBootstrapPipelineCall(t *testing.T) {
	t.Parallel()

	srvShutdown := make(chan struct{})
	srv := newServer(
		func(ctx context.Context, call *server.Call) error {
			resp, err := call.AllocResults(capnp.ObjectSize{DataSize: 8})
			if err != nil {
				return err
			}
			resp.SetUint64(0, 0xdeadbeef|uint64(call.Args().Uint32(0))<<32)
			return nil
		},
		func() {
			close(srvShutdown)
		})
	left, right := transport.NewPipe(1)
	p1, p2 := rpc.NewTransport(left), rpc.NewTransport(right)

	conn := rpc.NewConn(p1, &rpc.Options{
		BootstrapClient: srv,
		ErrorReporter:   testErrorReporter{tb: t},
	})
	defer func() {
		finishTest(t, conn, p2)
		<-srvShutdown // Will hang if closing does not shut down the client.
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
	_, err := recvBootstrapReturn(ctx, p2, bootstrapQID)
	if err != nil {
		t.Fatal(err)
	}

	// 3. Write call
	const callQID = 55
	{
		outMsg, err := p2.NewMessage()
		if err != nil {
			t.Fatal("p2.NewMessage():", err)
		}
		params, err := capnp.NewStruct(outMsg.Message().Segment(), capnp.ObjectSize{DataSize: 8})
		if err != nil {
			t.Fatal("capnp.NewStruct:", err)
		}
		params.SetUint32(0, 0x2a2b)
		err = pogs.Insert(rpccp.Message_TypeID, capnp.Struct(outMsg.Message()), &rpcMessage{
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
			outMsg.Release()
			t.Fatal("pogs.Insert(p2.NewMessage(), &rpcMessage{...}):", err)
		}
		err = outMsg.Send()
		outMsg.Release()
		if err != nil {
			t.Fatal("send():", err)
		}
	}

	// 4. Read return
	{
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
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

// TestDuplicateBootstrap calls Bootstrap twice on the same connection,
// and verifies that the results are the same.
func TestDuplicateBootstrap(t *testing.T) {
	t.Parallel()

	left, right := transport.NewPipe(1)
	p1, p2 := rpc.NewTransport(left), rpc.NewTransport(right)

	srv := newServer(func(ctx context.Context, call *server.Call) error {
		return nil
	}, nil)

	srvConn := rpc.NewConn(p1, &rpc.Options{
		BootstrapClient: srv,
		ErrorReporter:   testErrorReporter{tb: t},
	})
	defer srvConn.Close()

	clientConn := rpc.NewConn(p2, &rpc.Options{
		ErrorReporter: testErrorReporter{tb: t},
	})
	defer clientConn.Close()

	ctx := context.Background()
	bs1 := clientConn.Bootstrap(ctx)
	bs2 := clientConn.Bootstrap(ctx)
	assert.NoError(t, bs1.Resolve(ctx))
	assert.NoError(t, bs2.Resolve(ctx))
	assert.True(t, bs1.IsValid())
	assert.True(t, bs2.IsValid())
	if os.Getenv("FLAKY_TESTS") == "1" {
		// This is currently failing, see #523
		assert.True(t, bs1.IsSame(bs2))
	}

	bs1.Release()
	bs2.Release()
}

// TestUseConnAfterBootstrapError triggers a failed bootstrap call, and then
// verifies that the connection still works.
func TestUseConnAfterBootstrapError(t *testing.T) {
	t.Parallel()

	left, right := transport.NewPipe(1)
	p1, p2 := rpc.NewTransport(left), rpc.NewTransport(right)

	srv := newServer(func(ctx context.Context, call *server.Call) error {
		return nil
	}, nil)

	srvConn := rpc.NewConn(p1, &rpc.Options{
		BootstrapClient: srv,
		ErrorReporter:   testErrorReporter{tb: t},
	})
	defer srvConn.Close()

	clientConn := rpc.NewConn(p2, &rpc.Options{
		ErrorReporter: testErrorReporter{tb: t},
	})
	defer clientConn.Close()

	// srv -> client bootstrap should fail.
	ctx := context.Background()
	clientBs := srvConn.Bootstrap(ctx)
	assert.NoError(t, clientBs.Resolve(ctx))
	snapshotBs := clientBs.Snapshot()
	err, ok := snapshotBs.Brand().Value.(error)
	snapshotBs.Release()
	assert.True(t, ok)
	assert.NotNil(t, err)
	clientBs.Release()

	// client -> srv bootstrap should succeed.
	srvBs := clientConn.Bootstrap(ctx)
	defer srvBs.Release()
	assert.NoError(t, srvBs.Resolve(ctx))

	// ...and we should be able to make a call on it:
	ans, rel := srvBs.SendCall(ctx, capnp.Send{
		Method: capnp.Method{
			InterfaceID: interfaceID,
			MethodID:    methodID,
		},
		ArgsSize:  capnp.ObjectSize{},
		PlaceArgs: nil,
	})
	defer rel()
	_, err = ans.Struct()
	assert.NoError(t, err)
}

// TestCallOnClosedConn obtains the bootstrap capability, closes the
// connection, then attempts to make a call on the capability, verifying
// that the call returns a disconnected error.  Level 0 requirement.
func TestCallOnClosedConn(t *testing.T) {
	t.Parallel()

	left, right := transport.NewPipe(1)
	p1, p2 := rpc.NewTransport(left), rpc.NewTransport(right)

	defer p2.Close()
	conn := rpc.NewConn(p1, &rpc.Options{
		ErrorReporter: testErrorReporter{tb: t},
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
	defer client.Release()
	var qid uint32
	{
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if rmsg.Which != rpccp.Message_Which_bootstrap {
			t.Fatalf("Received %v message; want bootstrap", rmsg.Which)
		}
		qid = rmsg.Bootstrap.QuestionID
	}

	// 2. Write back a return
	{
		outMsg, err := p2.NewMessage()
		if err != nil {
			t.Fatal("p2.NewMessage():", err)
		}
		iptr := capnp.NewInterface(outMsg.Message().Segment(), 0)
		err = pogs.Insert(rpccp.Message_TypeID, capnp.Struct(outMsg.Message()), &rpcMessage{
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
			outMsg.Release()
			t.Fatal("pogs.Insert(p2.NewMessage(), &rpcMessage{...}):", err)
		}
		err = outMsg.Send()
		outMsg.Release()
		if err != nil {
			t.Fatal("send():", err)
		}
	}

	// 3. Read finish after client is resolved.
	{
		if err := client.Resolve(ctx); err != nil {
			t.Error("client.Resolve:", err)
		}
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if rmsg.Which != rpccp.Message_Which_finish {
			t.Fatalf("Received %v message; want finish", rmsg.Which)
		}
		if rmsg.Finish.QuestionID != qid {
			t.Errorf("Received finish for question %d; want %d", rmsg.Finish.QuestionID, qid)
		}
		if rmsg.Finish.ReleaseResultCaps {
			t.Error("Received finish that releases bootstrap result capabilities")
		}
	}

	// 4. Close the Conn.
	closed = true
	if err := conn.Close(); err != nil {
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
	_, err := ans.Struct()
	if !capnp.IsDisconnected(err) {
		t.Errorf("call after Close returned error: %v; want disconnected", err)
	}
}

// TestRecvCancel makes a call, sends a finish before it returns, then
// checks to see whether the call's Context was canceled and whether the
// capability the call returned is released.  Level 0 requirement.
func TestRecvCancel(t *testing.T) {
	t.Parallel()

	callCancel := make(chan struct{})
	retcapShutdown := make(chan struct{})
	srv := newServer(func(ctx context.Context, call *server.Call) error {
		// Wait until canceled
		call.Go()
		<-ctx.Done()
		close(callCancel)

		// Return a capability
		resp, err := call.AllocResults(capnp.ObjectSize{PointerCount: 1})
		if err != nil {
			t.Error("alloc results:", err)
			close(retcapShutdown)
			return err
		}
		retcap := newServer(nil, func() { close(retcapShutdown) })
		capID := resp.Message().CapTable().Add(retcap)
		if err := resp.SetPtr(0, capnp.NewInterface(resp.Segment(), capID).ToPtr()); err != nil {
			t.Error("set pointer:", err)
			return err
		}
		return nil
	}, nil)
	left, right := transport.NewPipe(1)
	p1, p2 := rpc.NewTransport(left), rpc.NewTransport(right)

	defer p2.Close()
	conn := rpc.NewConn(p1, &rpc.Options{
		BootstrapClient: srv,
		ErrorReporter:   testErrorReporter{tb: t},
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

	// 2. Write call
	const callQID = 55
	{
		outMsg, err := p2.NewMessage()
		if err != nil {
			t.Fatal("p2.NewMessage():", err)
		}
		err = pogs.Insert(rpccp.Message_TypeID, capnp.Struct(outMsg.Message()), &rpcMessage{
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
			},
		})
		if err != nil {
			outMsg.Release()
			t.Fatal("pogs.Insert(p2.NewMessage(), &rpcMessage{...}):", err)
		}
		err = outMsg.Send()
		outMsg.Release()
		if err != nil {
			t.Fatal("send():", err)
		}
	}

	// 3. Read bootstrap return
	{
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if rmsg.Which != rpccp.Message_Which_return {
			t.Fatalf("Received %v message; want return", rmsg.Which)
		}
		if rmsg.Return.AnswerID != bootstrapQID {
			t.Errorf("Received return for answer %d; want %d", rmsg.Return.AnswerID, bootstrapQID)
		}
	}

	// 4. Write bootstrap finish
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

	// 5. Write call cancel and verify that call's Context was canceled.
	select {
	case <-callCancel:
		t.Error("call context done before cancel written")
	default:
	}
	{
		msg := &rpcMessage{
			Which: rpccp.Message_Which_finish,
			Finish: &rpcFinish{
				QuestionID:        callQID,
				ReleaseResultCaps: true,
			},
		}
		if err := sendMessage(ctx, p2, msg); err != nil {
			t.Fatal(err)
		}
	}
	<-callCancel

	// 6. Read call return
	{
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if rmsg.Which != rpccp.Message_Which_return {
			t.Fatalf("Received %v message; want return", rmsg.Which)
		}
		if rmsg.Return.AnswerID != callQID {
			t.Errorf("Received return for answer %d; want %d", rmsg.Return.AnswerID, callQID)
		}
		// Don't care whether results, exception, or canceled.
	}

	// 7. Close the Conn
	closed = true
	if err := conn.Close(); err != nil {
		t.Error(err)
	}

	// 8. Verify that returned capability was shut down.
	// There's no guarantee exactly when the release/shutdown will happen,
	// but Close should trigger it. Otherwise, this will hang:
	<-retcapShutdown
}

// TestSendCancel makes a call, cancels the Context, then checks to
// see whether a finish message was sent.  Level 0 requirement.
func TestSendCancel(t *testing.T) {
	t.Parallel()

	left, right := transport.NewPipe(1)
	p1, p2 := rpc.NewTransport(left), rpc.NewTransport(right)

	conn := rpc.NewConn(p1, &rpc.Options{
		ErrorReporter: testErrorReporter{tb: t},
	})
	defer finishTest(t, conn, p2)
	ctx := context.Background()

	// 1. Read bootstrap.
	client := conn.Bootstrap(ctx)
	defer client.Release()
	var bootQID uint32
	{
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if rmsg.Which != rpccp.Message_Which_bootstrap {
			t.Fatalf("Received %v message; want bootstrap", rmsg.Which)
		}
		bootQID = rmsg.Bootstrap.QuestionID
	}

	// 2. Write back a return.
	{
		outMsg, err := p2.NewMessage()
		if err != nil {
			t.Fatal("p2.NewMessage():", err)
		}
		iptr := capnp.NewInterface(outMsg.Message().Segment(), 0)
		err = pogs.Insert(rpccp.Message_TypeID, capnp.Struct(outMsg.Message()), &rpcMessage{
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
			outMsg.Release()
			t.Fatal("pogs.Insert(p2.NewMessage(), &rpcMessage{...}):", err)
		}
		err = outMsg.Send()
		outMsg.Release()
		if err != nil {
			t.Fatal("send():", err)
		}
	}

	// 3. Read bootstrap finish.
	{
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if rmsg.Which != rpccp.Message_Which_finish {
			t.Fatalf("Received %v message; want finish", rmsg.Which)
		}
		if rmsg.Finish.QuestionID != bootQID {
			t.Errorf("Received finish for question %d; want %d", rmsg.Finish.QuestionID, bootQID)
		}
	}

	// 4. Make a call.
	callCtx, cancelCall := context.WithCancel(ctx)
	ans, releaseCall := client.SendCall(callCtx, capnp.Send{
		Method: capnp.Method{
			InterfaceID: interfaceID,
			MethodID:    methodID,
		},
	})
	defer releaseCall()
	var callQID uint32
	{
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if rmsg.Which != rpccp.Message_Which_call {
			t.Fatalf("Received %v message; want call", rmsg.Which)
		}
		callQID = rmsg.Call.QuestionID
		if rmsg.Call.InterfaceID != interfaceID {
			t.Errorf("call.interfaceId = %x; want %x", rmsg.Call.InterfaceID, interfaceID)
		}
		if rmsg.Call.MethodID != methodID {
			t.Errorf("call.methodId = %x; want %x", rmsg.Call.MethodID, methodID)
		}
	}

	// 5. Cancel the call.
	cancelCall()
	{
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if rmsg.Which != rpccp.Message_Which_finish {
			t.Fatalf("Received %v message; want finish", rmsg.Which)
		}
		if rmsg.Finish.QuestionID != callQID {
			t.Errorf("finish.questionId = %d; want %d", rmsg.Finish.QuestionID, callQID)
		}
		if !rmsg.Finish.ReleaseResultCaps {
			t.Error("finish.releaseResultCaps = false; want true")
		}
	}

	// 6. Verify that answer finishes without any other input.
	<-ans.Done()
	releaseCall()

	// 7. Write canceled return.
	{
		msg := &rpcMessage{
			Which: rpccp.Message_Which_return,
			Return: &rpcReturn{
				AnswerID:         callQID,
				ReleaseParamCaps: false,
				Which:            rpccp.Return_Which_canceled,
			},
		}
		if err := sendMessage(ctx, p2, msg); err != nil {
			t.Fatal(err)
		}
	}

	// 8. Release client (avoid filling pipe buffer).
	client.Release()
	{
		rmsg, release, err := recvMessage(ctx, p2)
		if err != nil {
			t.Fatal("recvMessage(ctx, p2):", err)
		}
		defer release()
		if rmsg.Which != rpccp.Message_Which_release {
			t.Fatalf("Received %v message; want release", rmsg.Which)
		}
		if rmsg.Release.ID != bootstrapExportID {
			t.Errorf("release.id = %d; want %d", rmsg.Release.ID, bootstrapExportID)
		}
		if rmsg.Release.ReferenceCount != 1 {
			t.Errorf("release.id = %d; want 1", rmsg.Release.ReferenceCount)
		}
	}
}

func TestHandleReturn_regression(t *testing.T) {
	t.Parallel()

	// Common setup for the tests below. Creates a connection with
	// a suitable bootstrap interface on one end, then passes the
	// other end of the connection to the callback.
	withConn := func(f func(*rpc.Conn)) {
		p1, p2 := transport.NewPipe(1)

		srv := testcp.PingPong_ServerToClient(pingPongServer{})
		conn1 := rpc.NewConn(rpc.NewTransport(p2), &rpc.Options{
			BootstrapClient: capnp.Client(srv),
		})
		defer conn1.Close()

		conn2 := rpc.NewConn(rpc.NewTransport(p1), nil)
		defer conn2.Close()
		f(conn2)
	}

	t.Run("MethodCallWithExpiredContext", func(t *testing.T) {
		withConn(func(conn *rpc.Conn) {
			pp := testcp.PingPong(conn.Bootstrap(context.Background()))
			defer pp.Release()

			// create an EXPIRED context
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			f, release := pp.EchoNum(ctx, func(ps testcp.PingPong_echoNum_Params) error {
				ps.SetN(42)
				return nil
			})
			defer release()

			_, err := f.Struct()
			assert.ErrorIs(t, err, ctx.Err())
		})
	})

	t.Run("BootstrapWithExpiredContext", func(t *testing.T) {
		withConn(func(conn *rpc.Conn) {
			// create an EXPIRED context
			ctx, cancel := context.WithCancel(context.Background())
			cancel()

			// NOTE: bootstrap with expired context
			pp := testcp.PingPong(conn.Bootstrap(ctx))
			defer pp.Release()

			f, release := pp.EchoNum(ctx, func(ps testcp.PingPong_echoNum_Params) error {
				ps.SetN(42)
				return nil
			})
			defer release()

			_, err := f.Struct()
			assert.ErrorIs(t, err, ctx.Err())
		})
	})
}

func TestPromisedCapability(t *testing.T) {
	t.Parallel()

	var g errgroup.Group
	for i := 0; i < 1024; i++ {
		g.Go(func() error {
			ppp := testcp.PingPongProvider_ServerToClient(pingPongProvider{})
			defer ppp.Release()

			f, release := ppp.PingPong(context.Background(), nil)
			defer release()

			pp := f.PingPong()
			return capnp.Client(pp).Resolve(context.Background())
		})
	}

	assert.NoError(t, g.Wait())
}

type pingPongProvider struct{}

func (pingPongProvider) PingPong(ctx context.Context, call testcp.PingPongProvider_pingPong) error {
	res, err := call.AllocResults()
	if err == nil {
		pp := testcp.PingPong_ServerToClient(pingPonger{})
		err = res.SetPingPong(pp)
	}

	return err
}

type pingPonger struct{}

func (pingPonger) EchoNum(ctx context.Context, call testcp.PingPong_echoNum) error {
	results, err := call.AllocResults()
	if err != nil {
		return err
	}
	results.SetN(call.Args().N())
	return nil
}

// finishTest drains both sides of a pipe and reports any errors to t.
func finishTest(t errorfer, conn *rpc.Conn, p2 rpc.Transport) {
	ctx, cancel := context.WithCancel(context.Background())
	drained := make(chan struct{})
	go func() {
		defer close(drained)
		for {
			m, release, err := recvMessage(ctx, p2)
			if err != nil {
				return
			}
			w := m.Which
			release()
			switch w {
			case rpccp.Message_Which_abort:
				return
			case rpccp.Message_Which_release, rpccp.Message_Which_finish:
				// Ignore clean-up messages.
			default:
				// Notify if test ignored a substantial message.
				t.Errorf("conn sent a %v message while finishing test", w)
			}
		}
	}()
	if err := conn.Close(); err != nil {
		t.Errorf("conn.Close(): %v", err)
	}
	cancel()
	<-drained
	if err := p2.Close(); err != nil {
		t.Errorf("p2.Close(): %v", err)
	}
}

type errorfer interface {
	Errorf(string, ...any)
}

func newServer(impl func(context.Context, *server.Call) error, shutdown shutdownFunc) capnp.Client {
	var methods []server.Method
	if impl != nil {
		methods = []server.Method{{
			Method: capnp.Method{
				InterfaceID: interfaceID,
				MethodID:    methodID,
			},
			Impl: impl,
		}}
	}
	return capnp.NewClient(server.New(methods, nil /* brand */, shutdown))
}

type shutdownFunc func()

func (f shutdownFunc) Shutdown() {
	if f == nil {
		return
	}
	f()
}

type rpcMessage struct {
	Which         rpccp.Message_Which
	Unimplemented *rpcMessage
	Abort         *rpcException
	Bootstrap     *rpcBootstrap
	Call          *rpcCall
	Return        *rpcReturn
	Finish        *rpcFinish
	Resolve       *rpcResolve
	Release       *rpcRelease
	Disembargo    *rpcDisembargo
}

func sendMessage(ctx context.Context, t rpc.Transport, msg *rpcMessage) error {
	outMsg, err := t.NewMessage()
	if err != nil {
		return fmt.Errorf("send message: %v", err)
	}
	defer outMsg.Release()
	if err := pogs.Insert(rpccp.Message_TypeID, capnp.Struct(outMsg.Message()), msg); err != nil {
		return fmt.Errorf("send message: %v", err)
	}
	if err := outMsg.Send(); err != nil {
		return fmt.Errorf("send message: %v", err)
	}
	return nil
}

func recvMessage(ctx context.Context, t rpc.Transport) (*rpcMessage, capnp.ReleaseFunc, error) {
	inMsg, err := t.RecvMessage()
	if err != nil {
		return nil, nil, err
	}
	r := new(rpcMessage)
	if err := pogs.Extract(r, rpccp.Message_TypeID, capnp.Struct(inMsg.Message())); err != nil {
		inMsg.Release()
		return nil, nil, fmt.Errorf("extract RPC message: %v", err)
	}
	if r.Which == rpccp.Message_Which_abort ||
		r.Which == rpccp.Message_Which_bootstrap ||
		r.Which == rpccp.Message_Which_finish ||
		r.Which == rpccp.Message_Which_resolve ||
		r.Which == rpccp.Message_Which_release ||
		r.Which == rpccp.Message_Which_disembargo {
		// These messages are guaranteed to not contain pointers back to
		// the original message, so we can release them early.
		inMsg.Release()
		return r, func() {}, nil
	}
	return r, inMsg.Release, nil
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

type rpcPromisedAnswer struct {
	QuestionID uint32 `capnp:"questionId"`
	Transform  []rpcPromisedAnswerOp
}

func (pa *rpcPromisedAnswer) transformEquals(path ...uint16) bool {
	for _, op := range pa.Transform {
		switch op.Which {
		case rpccp.PromisedAnswer_Op_Which_noop:
			// Skip.
		case rpccp.PromisedAnswer_Op_Which_getPointerField:
			if len(path) == 0 || path[0] != op.GetPointerField {
				return false
			}
			path = path[1:]
		default:
			return false
		}
	}
	return len(path) == 0
}

type rpcPromisedAnswerOp struct {
	Which           rpccp.PromisedAnswer_Op_Which
	GetPointerField uint16
}

func recvBootstrapReturn(ctx context.Context, t rpc.Transport, qid uint32) (uint32, error) {
	rmsg, release, err := recvMessage(ctx, t)
	if err != nil {
		return 0, fmt.Errorf("receive bootstrap: %v", err)
	}
	defer release()
	if rmsg.Which != rpccp.Message_Which_return {
		return 0, fmt.Errorf("received %v message; want return (for bootstrap)", rmsg.Which)
	}
	if rmsg.Return.AnswerID != qid {
		return 0, fmt.Errorf("received return for answer %d; want %d (bootstrap)", rmsg.Return.AnswerID, qid)
	}
	if rmsg.Return.Which != rpccp.Return_Which_results {
		return 0, fmt.Errorf("bootstrap return which = %v; want results", rmsg.Return.Which)
	}
	iface := rmsg.Return.Results.Content.Interface()
	if !iface.IsValid() {
		return 0, errors.New("parse bootstrap return: content is not an interface pointer")
	}
	ctab := rmsg.Return.Results.CapTable
	if iface.Capability() != 0 || len(ctab) != 1 {
		// This is a bit more restrictive than necessary, but we don't need
		// the flexibility.
		return 0, fmt.Errorf("parse bootstrap return: capability index, table length = %d, %d; want 0, 1", iface.Capability(), len(ctab))
	}
	if ctab[0].Which != rpccp.CapDescriptor_Which_senderHosted {
		return 0, fmt.Errorf("parse bootstrap return: received %v capability; want senderHosted", ctab[0].Which)
	}
	return ctab[0].SenderHosted, nil
}

func canceledContext(parent context.Context) context.Context {
	ctx, cancel := context.WithCancel(parent)
	cancel()
	return ctx
}

type testErrorReporter struct {
	tb interface {
		Log(...any)
		Fail()
	}
	fail bool
}

func (r testErrorReporter) ReportError(e error) {
	r.tb.Log("conn error:", e)
	if r.fail {
		r.tb.Fail()
	}
}
