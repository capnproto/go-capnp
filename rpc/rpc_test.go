package rpc_test

import (
	"errors"
	"io"
	"testing"

	"golang.org/x/net/context"
	"zombiezen.com/go/capnproto"
	zrpc "zombiezen.com/go/capnproto/rpc"
	"zombiezen.com/go/capnproto/rpc/internal/rpc"
)

const (
	interfaceID       uint64 = 0xa7317bd7216570aa
	methodID          uint16 = 9
	bootstrapExportID uint32 = 84
)

func newTestConn(t *testing.T, options ...zrpc.ConnOption) (*zrpc.Conn, rwcPipe) {
	p, q := newPipe()
	c := zrpc.NewConn(p, options...)
	return c, q
}

func TestBootstrap(t *testing.T) {
	ctx := context.Background()
	conn, p := newTestConn(t)
	defer conn.Close()
	defer p.Close()

	readBootstrap(t, ctx, conn, p)
}

func readBootstrap(t *testing.T, ctx context.Context, conn *zrpc.Conn, p rwcPipe) (client capnp.Client, questionID uint32) {
	clientCh := make(chan capnp.Client, 1)
	clientCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		clientCh <- conn.Bootstrap(clientCtx)
	}()

	msg, err := readMessage(p)
	if err != nil {
		t.Fatal("Read Bootstrap failed:", err)
	}
	if msg.Which() != rpc.Message_Which_bootstrap {
		t.Fatalf("Received %v message from bootstrap, want Message_Which_bootstrap", msg.Which())
	}
	questionID = msg.Bootstrap().QuestionId()
	// If this deadlocks, then Bootstrap isn't using a promised client.
	client = <-clientCh
	if client == nil {
		t.Fatal("Bootstrap client is nil")
	}
	return
}

func TestBootstrapFulfilled(t *testing.T) {
	ctx := context.Background()
	conn, p := newTestConn(t)
	defer conn.Close()
	defer p.Close()

	bootstrapAndFulfill(t, ctx, conn, p)
}

func bootstrapAndFulfill(t *testing.T, ctx context.Context, conn *zrpc.Conn, p rwcPipe) capnp.Client {
	client, bootstrapID := readBootstrap(t, ctx, conn, p)

	err := writeMessage(p, func(msg rpc.Message) {
		ret := rpc.NewReturn(msg.Segment)
		ret.SetAnswerId(bootstrapID)
		payload := rpc.NewPayload(msg.Segment)
		payload.SetContent(capnp.Object(msg.Segment.NewInterface(0)))
		capTable := rpc.NewCapDescriptor_List(msg.Segment, 1)
		capTable.At(0).SetSenderHosted(bootstrapExportID)
		payload.SetCapTable(capTable)
		ret.SetResults(payload)
		msg.SetReturn(ret)
	})
	if err != nil {
		t.Fatal("error writing Return:", err)
	}

	if finish, err := readMessage(p); err != nil {
		t.Fatal("error reading Finish:", err)
	} else if finish.Which() != rpc.Message_Which_finish {
		t.Fatalf("message sent is %v; want Message_Which_finish", finish.Which())
	} else {
		if id := finish.Finish().QuestionId(); id != bootstrapID {
			t.Fatalf("finish question ID is %d; want %d", id, bootstrapID)
		}
		if finish.Finish().ReleaseResultCaps() {
			t.Fatal("finish released bootstrap capability")
		}
	}
	return client
}

func TestCallOnPromisedAnswer(t *testing.T) {
	ctx := context.Background()
	conn, p := newTestConn(t)
	defer conn.Close()
	defer p.Close()
	client, bootstrapID := readBootstrap(t, ctx, conn, p)

	readDone := startReadMessage(p)
	client.Call(&capnp.Call{
		Ctx: ctx,
		Method: capnp.Method{
			InterfaceID: interfaceID,
			MethodID:    methodID,
		},
		ParamsSize: capnp.ObjectSize{DataSize: 8},
		ParamsFunc: func(s capnp.Struct) { s.Set64(0, 42) },
	})
	read := <-readDone

	if read.err != nil {
		t.Fatal("Reading failed:", read.err)
	}
	if read.msg.Which() == rpc.Message_Which_call {
		if target := read.msg.Call().Target(); target.Which() == rpc.MessageTarget_Which_promisedAnswer {
			if qid := target.PromisedAnswer().QuestionId(); qid != bootstrapID {
				t.Errorf("Target question ID = %d; want %d", qid, bootstrapID)
			}
			// TODO(light): allow no-ops
			if xform := target.PromisedAnswer().Transform(); xform.Len() != 0 {
				t.Error("Target transform is non-empty")
			}
		} else {
			t.Errorf("Target is %v, want MessageTarget_Which_promisedAnswer", target.Which())
		}
		if id := read.msg.Call().InterfaceId(); id != interfaceID {
			t.Errorf("Interface ID = %x; want %x", id, interfaceID)
		}
		if id := read.msg.Call().MethodId(); id != methodID {
			t.Errorf("Method ID = %d; want %d", id, methodID)
		}
		params := read.msg.Call().Params()
		if x := params.Content().ToStruct().Get64(0); x != 42 {
			t.Errorf("Params content value = %d; want %d", x, 42)
		}
		if sendResultsTo := read.msg.Call().SendResultsTo().Which(); sendResultsTo != rpc.Call_sendResultsTo_Which_caller {
			t.Errorf("Send results to %v; want caller", sendResultsTo)
		}
	} else {
		t.Errorf("Conn sent %v message, want Message_Which_call", read.msg.Which())
	}
}

func TestCallOnExportId(t *testing.T) {
	ctx := context.Background()
	conn, p := newTestConn(t)
	defer conn.Close()
	defer p.Close()
	client := bootstrapAndFulfill(t, ctx, conn, p)

	readDone := startReadMessage(p)
	client.Call(&capnp.Call{
		Ctx: ctx,
		Method: capnp.Method{
			InterfaceID: interfaceID,
			MethodID:    methodID,
		},
		ParamsSize: capnp.ObjectSize{DataSize: 8},
		ParamsFunc: func(s capnp.Struct) { s.Set64(0, 42) },
	})
	read := <-readDone

	if read.err != nil {
		t.Fatal("Reading failed:", read.err)
	}
	if read.msg.Which() == rpc.Message_Which_call {
		if target := read.msg.Call().Target(); target.Which() == rpc.MessageTarget_Which_importedCap {
			if id := target.ImportedCap(); id != bootstrapExportID {
				t.Errorf("Target imported cap = %d; want %d", id, bootstrapExportID)
			}
		} else {
			t.Errorf("Target is %v, want MessageTarget_Which_importedCap", target.Which())
		}
		if id := read.msg.Call().InterfaceId(); id != interfaceID {
			t.Errorf("Interface ID = %x; want %x", id, interfaceID)
		}
		if id := read.msg.Call().MethodId(); id != methodID {
			t.Errorf("Method ID = %d; want %d", id, methodID)
		}
		params := read.msg.Call().Params()
		if x := params.Content().ToStruct().Get64(0); x != 42 {
			t.Errorf("Params content value = %d; want %d", x, 42)
		}
		if sendResultsTo := read.msg.Call().SendResultsTo().Which(); sendResultsTo != rpc.Call_sendResultsTo_Which_caller {
			t.Errorf("Send results to %v; want caller", sendResultsTo)
		}
	} else {
		t.Errorf("Conn sent %v message, want Message_Which_call", read.msg.Which())
	}
}

func TestMainInterface(t *testing.T) {
	main := mockClient()
	conn, p := newTestConn(t, zrpc.MainInterface(main))
	defer conn.Close()
	defer p.Close()

	bootstrapRoundtrip(t, p)
}

func bootstrapRoundtrip(t *testing.T, p io.ReadWriter) (importID, questionID uint32) {
	questionID = 54
	err := writeMessage(p, func(msg rpc.Message) {
		bootstrap := rpc.NewBootstrap(msg.Segment)
		bootstrap.SetQuestionId(questionID)
		msg.SetBootstrap(bootstrap)
	})
	if err != nil {
		t.Fatal("Write Bootstrap failed:", err)
	}
	msg, err := readMessage(p)
	if err != nil {
		t.Fatal("Read Bootstrap response failed:", err)
	}

	if msg.Which() != rpc.Message_Which_return {
		t.Fatalf("Conn sent %v message, want Message_Which_return", msg.Which())
	}
	if id := msg.Return().AnswerId(); id != questionID {
		t.Fatalf("msg.Return().AnswerId() = %d; want %d", id, questionID)
	}
	if msg.Return().Which() != rpc.Return_Which_results {
		t.Fatalf("msg.Return().Which() = %v; want Return_Which_results", msg.Return().Which())
	}
	payload := msg.Return().Results()
	if tp := payload.Content().Type(); tp != capnp.TypeInterface {
		t.Fatalf("Result payload contains a %v; want interface", tp)
	}
	capIdx := int(payload.Content().ToInterface().Capability())
	if n := payload.CapTable().Len(); capIdx >= n {
		t.Fatalf("Payload capTable has size %d, but capability index = %d", n, capIdx)
	}
	if cw := payload.CapTable().At(capIdx).Which(); cw != rpc.CapDescriptor_Which_senderHosted {
		t.Fatalf("Capability type is %d; want CapDescriptor_Which_senderHosted", cw)
	}
	return payload.CapTable().At(capIdx).SenderHosted(), questionID
}

func TestReceiveCallOnPromisedAnswer(t *testing.T) {
	const questionID = 999
	called := false
	main := stubClient(func(ctx context.Context, params capnp.Struct) (capnp.Struct, error) {
		s := capnp.NewBuffer(nil)
		result := s.NewRootStruct(capnp.ObjectSize{})
		called = true
		return result, nil
	})
	conn, p := newTestConn(t, zrpc.MainInterface(main))
	defer conn.Close()
	defer p.Close()
	_, bootqID := bootstrapRoundtrip(t, p)

	err := writeMessage(p, func(msg rpc.Message) {
		call := rpc.NewCall(msg.Segment)
		call.SetQuestionId(questionID)
		call.SetInterfaceId(interfaceID)
		call.SetMethodId(methodID)
		target := rpc.NewMessageTarget(msg.Segment)
		pa := rpc.NewPromisedAnswer(msg.Segment)
		pa.SetQuestionId(bootqID)
		target.SetPromisedAnswer(pa)
		call.SetTarget(target)
		payload := rpc.NewPayload(msg.Segment)
		payload.SetContent(capnp.Object(msg.Segment.NewStruct(capnp.ObjectSize{})))
		call.SetParams(payload)
		msg.SetCall(call)
	})
	if err != nil {
		t.Fatal("Call message failed:", err)
	}
	retmsg, err := readMessage(p)
	if err != nil {
		t.Fatal("Read Call return failed:", err)
	}

	if !called {
		t.Error("interface not called")
	}
	if retmsg.Which() != rpc.Message_Which_return {
		t.Fatalf("Return message is %v; want %v", retmsg.Which(), rpc.Message_Which_return)
	}
	ret := retmsg.Return()
	if id := ret.AnswerId(); id != questionID {
		t.Errorf("Return.answerId = %d; want %d", id, questionID)
	}
	if ret.Which() == rpc.Return_Which_results {
		// TODO(light)
	} else if ret.Which() == rpc.Return_Which_exception {
		t.Error("Return.exception:", ret.Exception().Reason())
	} else {
		t.Errorf("Return.Which() = %v; want %v", ret.Which(), rpc.Return_Which_results)
	}
}

func TestReceiveCallOnExport(t *testing.T) {
	const questionID = 999
	called := false
	main := stubClient(func(ctx context.Context, params capnp.Struct) (capnp.Struct, error) {
		s := capnp.NewBuffer(nil)
		result := s.NewRootStruct(capnp.ObjectSize{})
		called = true
		return result, nil
	})
	conn, p := newTestConn(t, zrpc.MainInterface(main))
	defer conn.Close()
	defer p.Close()
	importID := sendBootstrapAndFinish(t, p)

	err := writeMessage(p, func(msg rpc.Message) {
		call := rpc.NewCall(msg.Segment)
		call.SetQuestionId(questionID)
		call.SetInterfaceId(interfaceID)
		call.SetMethodId(methodID)
		target := rpc.NewMessageTarget(msg.Segment)
		target.SetImportedCap(importID)
		call.SetTarget(target)
		payload := rpc.NewPayload(msg.Segment)
		payload.SetContent(capnp.Object(msg.Segment.NewStruct(capnp.ObjectSize{})))
		call.SetParams(payload)
		msg.SetCall(call)
	})
	if err != nil {
		t.Fatal("Call message failed:", err)
	}
	retmsg, err := readMessage(p)
	if err != nil {
		t.Fatal("Read Call return failed:", err)
	}

	if !called {
		t.Error("interface not called")
	}
	if retmsg.Which() != rpc.Message_Which_return {
		t.Fatalf("Return message is %v; want %v", retmsg.Which(), rpc.Message_Which_return)
	}
	ret := retmsg.Return()
	if id := ret.AnswerId(); id != questionID {
		t.Errorf("Return.answerId = %d; want %d", id, questionID)
	}
	if ret.Which() == rpc.Return_Which_results {
		// TODO(light)
	} else if ret.Which() == rpc.Return_Which_exception {
		t.Error("Return.exception:", ret.Exception().Reason())
	} else {
		t.Errorf("Return.Which() = %v; want %v", ret.Which(), rpc.Return_Which_results)
	}
}

func sendBootstrapAndFinish(t *testing.T, p io.ReadWriter) (importID uint32) {
	importID, questionID := bootstrapRoundtrip(t, p)
	err := writeMessage(p, func(msg rpc.Message) {
		finish := rpc.NewFinish(msg.Segment)
		finish.SetQuestionId(questionID)
		finish.SetReleaseResultCaps(false)
		msg.SetFinish(finish)
	})
	if err != nil {
		t.Fatal("Write Bootstrap Finish failed:", err)
	}
	return importID
}

func readMessage(r io.Reader) (rpc.Message, error) {
	s, err := capnp.ReadFromStream(r, nil)
	var msg rpc.Message
	if err == nil {
		msg = rpc.ReadRootMessage(s)
	}
	return msg, err
}

func writeMessage(w io.Writer, f func(rpc.Message)) error {
	s := capnp.NewBuffer(nil)
	m := rpc.NewRootMessage(s)
	f(m)
	_, err := m.Segment.WriteTo(w)
	return err
}

func startReadMessage(r io.Reader) <-chan asyncRead {
	ch := make(chan asyncRead, 1)
	go func() {
		msg, err := readMessage(r)
		ch <- asyncRead{msg, err}
	}()
	return ch
}

type asyncRead struct {
	msg rpc.Message
	err error
}

func mockClient() capnp.Client {
	return capnp.ErrorClient(errors.New("mock client"))
}

type stubClient func(ctx context.Context, params capnp.Struct) (capnp.Struct, error)

func (stub stubClient) Call(call *capnp.Call) capnp.Answer {
	if call.Method.InterfaceID != interfaceID || call.Method.MethodID != methodID {
		return capnp.ErrorAnswer(errors.New("stub client method not implemented"))
	}
	s, err := stub(call.Ctx, call.PlaceParams(nil))
	if err != nil {
		return capnp.ErrorAnswer(err)
	}
	return capnp.ImmediateAnswer(capnp.Object(s))
}

func (stub stubClient) Close() error {
	return nil
}

// rwcPipe is a synchronous in-memory pipe that implements io.ReadWriteCloser.
type rwcPipe struct {
	*io.PipeReader
	*io.PipeWriter
}

// newPipe returns two pipes connected to each other.
func newPipe() (p, q rwcPipe) {
	pr, qw := io.Pipe()
	qr, pw := io.Pipe()
	p = rwcPipe{
		PipeReader: pr,
		PipeWriter: pw,
	}
	q = rwcPipe{
		PipeReader: qr,
		PipeWriter: qw,
	}
	return
}

func (p rwcPipe) Close() error {
	rerr := p.PipeReader.Close()
	werr := p.PipeWriter.Close()
	switch {
	case rerr != nil:
		return rerr
	case werr != nil:
		return werr
	default:
		return nil
	}
}

func (p rwcPipe) CloseWithError(err error) error {
	rerr := p.PipeReader.CloseWithError(err)
	werr := p.PipeWriter.CloseWithError(err)
	switch {
	case rerr != nil:
		return rerr
	case werr != nil:
		return werr
	default:
		return nil
	}
}
