package rpc_test

import (
	"bytes"
	"errors"
	"sync"
	"testing"

	"golang.org/x/net/context"
	"zombiezen.com/go/capnproto"
	"zombiezen.com/go/capnproto/rpc"
	"zombiezen.com/go/capnproto/rpc/rpccapnp"
)

const (
	interfaceID       uint64 = 0xa7317bd7216570aa
	methodID          uint16 = 9
	bootstrapExportID uint32 = 84
)

func newTestConn(t *testing.T, options ...rpc.ConnOption) (*rpc.Conn, rpc.Transport) {
	p, q := newPipe()
	c := rpc.NewConn(p, options...)
	return c, q
}

func TestBootstrap(t *testing.T) {
	ctx := context.Background()
	conn, p := newTestConn(t)
	defer conn.Close()
	defer p.Close()

	readBootstrap(t, ctx, conn, p)
}

func readBootstrap(t *testing.T, ctx context.Context, conn *rpc.Conn, p rpc.Transport) (client capnp.Client, questionID uint32) {
	clientCh := make(chan capnp.Client, 1)
	clientCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		clientCh <- conn.Bootstrap(clientCtx)
	}()

	msg, err := p.RecvMessage(ctx)
	if err != nil {
		t.Fatal("Read Bootstrap failed:", err)
	}
	if msg.Which() != rpccapnp.Message_Which_bootstrap {
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

func bootstrapAndFulfill(t *testing.T, ctx context.Context, conn *rpc.Conn, p rpc.Transport) capnp.Client {
	client, bootstrapID := readBootstrap(t, ctx, conn, p)

	err := sendMessage(ctx, p, func(msg rpccapnp.Message) {
		ret := rpccapnp.NewReturn(msg.Segment)
		ret.SetAnswerId(bootstrapID)
		payload := rpccapnp.NewPayload(msg.Segment)
		payload.SetContent(capnp.Object(msg.Segment.NewInterface(0)))
		capTable := rpccapnp.NewCapDescriptor_List(msg.Segment, 1)
		capTable.At(0).SetSenderHosted(bootstrapExportID)
		payload.SetCapTable(capTable)
		ret.SetResults(payload)
		msg.SetReturn(ret)
	})
	if err != nil {
		t.Fatal("error writing Return:", err)
	}

	if finish, err := p.RecvMessage(ctx); err != nil {
		t.Fatal("error reading Finish:", err)
	} else if finish.Which() != rpccapnp.Message_Which_finish {
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

	readDone := startRecvMessage(p)
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
	if read.msg.Which() == rpccapnp.Message_Which_call {
		if target := read.msg.Call().Target(); target.Which() == rpccapnp.MessageTarget_Which_promisedAnswer {
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
		if sendResultsTo := read.msg.Call().SendResultsTo().Which(); sendResultsTo != rpccapnp.Call_sendResultsTo_Which_caller {
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

	readDone := startRecvMessage(p)
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
	if read.msg.Which() == rpccapnp.Message_Which_call {
		if target := read.msg.Call().Target(); target.Which() == rpccapnp.MessageTarget_Which_importedCap {
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
		if sendResultsTo := read.msg.Call().SendResultsTo().Which(); sendResultsTo != rpccapnp.Call_sendResultsTo_Which_caller {
			t.Errorf("Send results to %v; want caller", sendResultsTo)
		}
	} else {
		t.Errorf("Conn sent %v message, want Message_Which_call", read.msg.Which())
	}
}

func TestMainInterface(t *testing.T) {
	main := mockClient()
	conn, p := newTestConn(t, rpc.MainInterface(main))
	defer conn.Close()
	defer p.Close()

	bootstrapRoundtrip(t, p)
}

func bootstrapRoundtrip(t *testing.T, p rpc.Transport) (importID, questionID uint32) {
	questionID = 54
	err := sendMessage(context.TODO(), p, func(msg rpccapnp.Message) {
		bootstrap := rpccapnp.NewBootstrap(msg.Segment)
		bootstrap.SetQuestionId(questionID)
		msg.SetBootstrap(bootstrap)
	})
	if err != nil {
		t.Fatal("Write Bootstrap failed:", err)
	}
	msg, err := p.RecvMessage(context.TODO())
	if err != nil {
		t.Fatal("Read Bootstrap response failed:", err)
	}

	if msg.Which() != rpccapnp.Message_Which_return {
		t.Fatalf("Conn sent %v message, want Message_Which_return", msg.Which())
	}
	if id := msg.Return().AnswerId(); id != questionID {
		t.Fatalf("msg.Return().AnswerId() = %d; want %d", id, questionID)
	}
	if msg.Return().Which() != rpccapnp.Return_Which_results {
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
	if cw := payload.CapTable().At(capIdx).Which(); cw != rpccapnp.CapDescriptor_Which_senderHosted {
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
	conn, p := newTestConn(t, rpc.MainInterface(main))
	defer conn.Close()
	defer p.Close()
	_, bootqID := bootstrapRoundtrip(t, p)

	err := sendMessage(context.TODO(), p, func(msg rpccapnp.Message) {
		call := rpccapnp.NewCall(msg.Segment)
		call.SetQuestionId(questionID)
		call.SetInterfaceId(interfaceID)
		call.SetMethodId(methodID)
		target := rpccapnp.NewMessageTarget(msg.Segment)
		pa := rpccapnp.NewPromisedAnswer(msg.Segment)
		pa.SetQuestionId(bootqID)
		target.SetPromisedAnswer(pa)
		call.SetTarget(target)
		payload := rpccapnp.NewPayload(msg.Segment)
		payload.SetContent(capnp.Object(msg.Segment.NewStruct(capnp.ObjectSize{})))
		call.SetParams(payload)
		msg.SetCall(call)
	})
	if err != nil {
		t.Fatal("Call message failed:", err)
	}
	retmsg, err := p.RecvMessage(context.TODO())
	if err != nil {
		t.Fatal("Read Call return failed:", err)
	}

	if !called {
		t.Error("interface not called")
	}
	if retmsg.Which() != rpccapnp.Message_Which_return {
		t.Fatalf("Return message is %v; want %v", retmsg.Which(), rpccapnp.Message_Which_return)
	}
	ret := retmsg.Return()
	if id := ret.AnswerId(); id != questionID {
		t.Errorf("Return.answerId = %d; want %d", id, questionID)
	}
	if ret.Which() == rpccapnp.Return_Which_results {
		// TODO(light)
	} else if ret.Which() == rpccapnp.Return_Which_exception {
		t.Error("Return.exception:", ret.Exception().Reason())
	} else {
		t.Errorf("Return.Which() = %v; want %v", ret.Which(), rpccapnp.Return_Which_results)
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
	conn, p := newTestConn(t, rpc.MainInterface(main))
	defer conn.Close()
	defer p.Close()
	importID := sendBootstrapAndFinish(t, p)

	err := sendMessage(context.TODO(), p, func(msg rpccapnp.Message) {
		call := rpccapnp.NewCall(msg.Segment)
		call.SetQuestionId(questionID)
		call.SetInterfaceId(interfaceID)
		call.SetMethodId(methodID)
		target := rpccapnp.NewMessageTarget(msg.Segment)
		target.SetImportedCap(importID)
		call.SetTarget(target)
		payload := rpccapnp.NewPayload(msg.Segment)
		payload.SetContent(capnp.Object(msg.Segment.NewStruct(capnp.ObjectSize{})))
		call.SetParams(payload)
		msg.SetCall(call)
	})
	if err != nil {
		t.Fatal("Call message failed:", err)
	}
	retmsg, err := p.RecvMessage(context.TODO())
	if err != nil {
		t.Fatal("Read Call return failed:", err)
	}

	if !called {
		t.Error("interface not called")
	}
	if retmsg.Which() != rpccapnp.Message_Which_return {
		t.Fatalf("Return message is %v; want %v", retmsg.Which(), rpccapnp.Message_Which_return)
	}
	ret := retmsg.Return()
	if id := ret.AnswerId(); id != questionID {
		t.Errorf("Return.answerId = %d; want %d", id, questionID)
	}
	if ret.Which() == rpccapnp.Return_Which_results {
		// TODO(light)
	} else if ret.Which() == rpccapnp.Return_Which_exception {
		t.Error("Return.exception:", ret.Exception().Reason())
	} else {
		t.Errorf("Return.Which() = %v; want %v", ret.Which(), rpccapnp.Return_Which_results)
	}
}

func sendBootstrapAndFinish(t *testing.T, p rpc.Transport) (importID uint32) {
	importID, questionID := bootstrapRoundtrip(t, p)
	err := sendMessage(context.TODO(), p, func(msg rpccapnp.Message) {
		finish := rpccapnp.NewFinish(msg.Segment)
		finish.SetQuestionId(questionID)
		finish.SetReleaseResultCaps(false)
		msg.SetFinish(finish)
	})
	if err != nil {
		t.Fatal("Write Bootstrap Finish failed:", err)
	}
	return importID
}

func sendMessage(ctx context.Context, t rpc.Transport, f func(rpccapnp.Message)) error {
	s := capnp.NewBuffer(nil)
	m := rpccapnp.NewRootMessage(s)
	f(m)
	return t.SendMessage(ctx, m)
}

func startRecvMessage(t rpc.Transport) <-chan asyncRecv {
	ch := make(chan asyncRecv, 1)
	go func() {
		msg, err := t.RecvMessage(context.TODO())
		ch <- asyncRecv{msg, err}
	}()
	return ch
}

type asyncRecv struct {
	msg rpccapnp.Message
	err error
}

func mockClient() capnp.Client {
	return capnp.ErrorClient(errMockClient)
}

type stubClient func(ctx context.Context, params capnp.Struct) (capnp.Struct, error)

func (stub stubClient) Call(call *capnp.Call) capnp.Answer {
	if call.Method.InterfaceID != interfaceID || call.Method.MethodID != methodID {
		return capnp.ErrorAnswer(errNotImplemented)
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

type pipeTransport struct {
	r      <-chan rpccapnp.Message
	w      chan<- rpccapnp.Message
	finish chan struct{}

	rbuf bytes.Buffer

	mu       sync.Mutex
	inflight int
	done     bool
}

func newPipe() (p, q rpc.Transport) {
	a, b := make(chan rpccapnp.Message), make(chan rpccapnp.Message)
	p = &pipeTransport{
		r:      a,
		w:      b,
		finish: make(chan struct{}),
	}
	q = &pipeTransport{
		r:      b,
		w:      a,
		finish: make(chan struct{}),
	}
	return
}

func (p *pipeTransport) SendMessage(ctx context.Context, msg rpccapnp.Message) error {
	if !p.startSend() {
		return errClosed
	}
	defer p.finishSend()

	buf := new(bytes.Buffer)
	_, err := msg.Segment.WriteTo(buf)
	if err != nil {
		return err
	}
	seg, _, err := capnp.ReadFromMemoryZeroCopy(buf.Bytes())
	if err != nil {
		return err
	}

	select {
	case p.w <- rpccapnp.ReadRootMessage(seg):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-p.finish:
		return errClosed
	}
}

func (p *pipeTransport) startSend() bool {
	p.mu.Lock()
	ok := !p.done
	if ok {
		p.inflight++
	}
	p.mu.Unlock()
	return ok
}

func (p *pipeTransport) finishSend() {
	p.mu.Lock()
	p.inflight--
	p.mu.Unlock()
}

func (p *pipeTransport) RecvMessage(ctx context.Context) (rpccapnp.Message, error) {
	// Scribble over shared buffer to test for race conditions.
	for b, i := p.rbuf.Bytes(), 0; i < len(b); i++ {
		b[i] = 0xff
	}
	p.rbuf.Reset()

	select {
	case msg, ok := <-p.r:
		if !ok {
			return rpccapnp.Message{}, errBrokenPipe
		}
		if _, err := msg.Segment.WriteTo(&p.rbuf); err != nil {
			return rpccapnp.Message{}, err
		}
		seg, _, err := capnp.ReadFromMemoryZeroCopy(p.rbuf.Bytes())
		if err != nil {
			return rpccapnp.Message{}, err
		}
		return rpccapnp.ReadRootMessage(seg), nil
	case <-ctx.Done():
		return rpccapnp.Message{}, ctx.Err()
	}
}

func (p *pipeTransport) Close() error {
	p.mu.Lock()
	done := p.done
	if !done {
		p.done = true
		close(p.finish)
		if p.inflight == 0 {
			close(p.w)
		}
	}
	p.mu.Unlock()
	if done {
		return errClosed
	}
	return nil
}

var (
	errBrokenPipe     = errors.New("rpc_test: broken pipe")
	errClosed         = errors.New("rpc_test: write to broken pipe")
	errMockClient     = errors.New("rpc_test: mock client")
	errNotImplemented = errors.New("rpc_test: stub client method not implemented")
)
