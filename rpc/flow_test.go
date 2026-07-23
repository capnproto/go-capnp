package rpc

// tests for streaming/flow control

import (
	"context"
	"errors"
	"net"
	"sync"
	"testing"
	"time"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/flowcontrol"
	"capnproto.org/go/capnp/v3/rpc/internal/testcapnp"
	"capnproto.org/go/capnp/v3/rpc/transport"
	serverpkg "capnproto.org/go/capnp/v3/server"
)

// observingLimiter reports when a call to StartMessage is waiting for a
// response. It lets the integration test observe the configured limiter
// directly instead of inferring its state from the receiver's message sizes.
type observingLimiter struct {
	flowcontrol.FlowLimiter

	mu                 sync.Mutex
	started, proceeded uint64
	changed            chan struct{}
}

type countingTransport struct {
	Transport

	mu          sync.Mutex
	newMessages int
}

func (t *countingTransport) NewMessage() (transport.OutgoingMessage, error) {
	msg, err := t.Transport.NewMessage()
	if err == nil {
		t.mu.Lock()
		t.newMessages++
		t.mu.Unlock()
	}
	return msg, err
}

func (t *countingTransport) count() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.newMessages
}

type blockingGateController struct {
	waitEntered chan struct{}
	grant       chan struct{}
	completions chan flowcontrol.MessageOutcomeKind
	poisons     chan error
}

func (g *blockingGateController) CommitMessage(uint64) (func(context.Context) error, func(flowcontrol.MessageOutcomeKind, error)) {
	return func(ctx context.Context) error {
			g.waitEntered <- struct{}{}
			select {
			case <-g.grant:
				return nil
			case <-ctx.Done():
				return ctx.Err()
			}
		}, func(kind flowcontrol.MessageOutcomeKind, _ error) {
			g.completions <- kind
		}
}

func (g *blockingGateController) Poison(err error) {
	g.poisons <- err
}

type blockingGateLimiter struct {
	gate *blockingGateController
}

func (*blockingGateLimiter) StartMessage(context.Context, uint64) (func(), error) {
	return func() {}, errors.New("legacy StartMessage called for prepared RPC send")
}

func (*blockingGateLimiter) Release() {}

func (l *blockingGateLimiter) GateNext() flowcontrol.GateNextController {
	return l.gate
}

var _ flowcontrol.GateNextFlowLimiter = (*blockingGateLimiter)(nil)

func newObservingLimiter(limit int64) *observingLimiter {
	return &observingLimiter{
		FlowLimiter: flowcontrol.NewFixedLimiter(limit),
		changed:     make(chan struct{}, 1),
	}
}

func (l *observingLimiter) StartMessage(ctx context.Context, size uint64) (func(), error) {
	l.mu.Lock()
	l.started++
	l.mu.Unlock()
	l.signalChanged()

	gotResponse, err := l.FlowLimiter.StartMessage(ctx, size)
	if err != nil {
		return gotResponse, err
	}

	l.mu.Lock()
	l.proceeded++
	l.mu.Unlock()
	l.signalChanged()
	return gotResponse, nil
}

func (l *observingLimiter) waitForBlocked(ctx context.Context) bool {
	for {
		l.mu.Lock()
		blocked := l.started > l.proceeded
		l.mu.Unlock()
		if blocked {
			return true
		}

		select {
		case <-ctx.Done():
			return false
		case <-l.changed:
		}
	}
}

func (l *observingLimiter) signalChanged() {
	select {
	case l.changed <- struct{}{}:
	default:
	}
}

// Test that attaching a fixed-size FlowLimiter results in actually limiting the
// flow of messages. The server holds accepted calls until the test releases it,
// which fills the client's limiter without depending on scheduler timing.
func TestFixedFlowLimit(t *testing.T) {
	// This test deliberately saturates an RPC connection and uses a short
	// deadline to detect a blocked client. Keep it serial so unrelated parallel
	// RPC tests cannot consume that deadline.

	limit := int64(1 << 20)

	clientConn, serverConn := net.Pipe()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	unblock := make(chan struct{})
	releaseServer := func() {
		select {
		case <-unblock:
		default:
			close(unblock)
		}
	}
	defer releaseServer()
	server := gatedStreamTestServer{unblock: unblock}
	done := make(chan struct{})
	go func() {
		// Server

		bootstrap := testcapnp.StreamTest_ServerToClient(server)
		conn := NewConn(NewStreamTransport(serverConn), &Options{
			BootstrapClient: capnp.Client(bootstrap),
		})
		defer conn.Close()
		<-ctx.Done()
		close(done)
	}()

	trans := NewStreamTransport(clientConn)
	conn := NewConn(trans, nil)
	defer conn.Close()

	client := testcapnp.StreamTest(conn.Bootstrap(ctx))
	defer client.Release()

	// Make a decently sized payload, so we can expect the size of the
	// parameters to dominate the size of rpc messages.
	data := make([]byte, 2048)
	limiter := newObservingLimiter(limit)
	capnp.Client(client).SetFlowLimiter(limiter)

	stop := make(chan struct{})
	sent := make(chan error, 1)
	go func() {
		for {
			select {
			case <-stop:
				sent <- nil
				return
			default:
			}
			if err := client.Push(ctx, func(p testcapnp.StreamTest_push_Params) error {
				return p.SetData(data)
			}); err != nil {
				sent <- err
				return
			}
		}
	}()

	if !limiter.waitForBlocked(ctx) {
		t.Fatal("flow limiter did not block while the server held responses")
	}
	select {
	case err := <-sent:
		t.Fatalf("all calls were sent while the server held responses: %v", err)
	default:
	}

	// The server cannot respond until this gate is opened. Once it is, each
	// response returns capacity to the limiter and lets the client proceed.
	close(stop)
	releaseServer()
	select {
	case err := <-sent:
		if err != nil {
			t.Fatal(err)
		}
	case <-ctx.Done():
		t.Fatal("client did not finish after the server released held calls")
	}

	cancel()
	<-done

}

func TestGateNextRPCImportPreparesAfterGate(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	unblockServer := make(chan struct{})
	serverDone := make(chan struct{})
	go func() {
		defer close(serverDone)
		bootstrap := testcapnp.StreamTest_ServerToClient(gatedStreamTestServer{unblock: unblockServer})
		conn := NewConn(NewStreamTransport(serverConn), &Options{BootstrapClient: capnp.Client(bootstrap)})
		defer conn.Close()
		<-ctx.Done()
	}()

	trans := &countingTransport{Transport: NewStreamTransport(clientConn)}
	conn := NewConn(trans, nil)
	defer conn.Close()
	client := testcapnp.StreamTest(conn.Bootstrap(ctx))
	defer client.Release()
	if err := client.Resolve(ctx); err != nil {
		t.Fatal(err)
	}

	gate := &blockingGateController{
		waitEntered: make(chan struct{}, 1),
		grant:       make(chan struct{}),
		completions: make(chan flowcontrol.MessageOutcomeKind, 4),
		poisons:     make(chan error, 1),
	}
	limiter := &blockingGateLimiter{gate: gate}
	capnp.Client(client).SetFlowLimiter(limiter)

	if err := client.Push(ctx, func(p testcapnp.StreamTest_push_Params) error {
		return p.SetData([]byte{1})
	}); err != nil {
		t.Fatal(err)
	}
	before := trans.count()

	placed := make(chan struct{})
	secondDone := make(chan error, 1)
	go func() {
		secondDone <- client.Push(ctx, func(p testcapnp.StreamTest_push_Params) error {
			close(placed)
			return p.SetData([]byte{2})
		})
	}()
	<-gate.waitEntered
	select {
	case <-placed:
		t.Fatal("RPC PlaceArgs ran behind a closed GateNext permission")
	default:
	}
	if got := trans.count(); got != before {
		t.Fatalf("transport NewMessage count = %d; want %d while gate is closed", got, before)
	}

	close(gate.grant)
	if err := <-secondDone; err != nil {
		t.Fatal(err)
	}
	<-placed
	if got := trans.count(); got != before+1 {
		t.Fatalf("transport NewMessage count = %d; want %d after gate opened", got, before+1)
	}

	close(unblockServer)
	if err := capnp.Client(client).WaitStreaming(); err != nil {
		t.Fatal(err)
	}
	for range 2 {
		if kind := <-gate.completions; kind != flowcontrol.MessageOutcomeSucceeded {
			t.Fatalf("terminal outcome = %v; want succeeded", kind)
		}
	}
	select {
	case kind := <-gate.completions:
		t.Fatalf("duplicate terminal completion %v", kind)
	default:
	}
	select {
	case err := <-gate.poisons:
		t.Fatalf("unexpected controller poison: %v", err)
	default:
	}
	cancel()
	<-serverDone
}

func TestGateNextRPCPipelinePreparesAfterGate(t *testing.T) {
	const (
		parentInterface = 0xf001
		parentMethod    = 1
		childInterface  = 0xf002
		childMethod     = 2
	)
	clientConn, serverConn := net.Pipe()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	child := capnp.NewClient(serverpkg.New([]serverpkg.Method{{
		Method: capnp.Method{InterfaceID: childInterface, MethodID: childMethod},
		Impl: func(context.Context, *serverpkg.Call) error {
			return nil
		},
	}}, nil, nil))
	defer child.Release()
	unblockParent := make(chan struct{})
	parent := capnp.NewClient(serverpkg.New([]serverpkg.Method{{
		Method: capnp.Method{InterfaceID: parentInterface, MethodID: parentMethod},
		Impl: func(_ context.Context, call *serverpkg.Call) error {
			call.Go()
			<-unblockParent
			results, err := call.AllocResults(capnp.ObjectSize{PointerCount: 1})
			if err != nil {
				return err
			}
			id := results.Message().CapTable().Add(child)
			return results.SetPtr(0, capnp.NewInterface(results.Segment(), id).ToPtr())
		},
	}}, nil, nil))
	defer parent.Release()

	serverDone := make(chan struct{})
	go func() {
		defer close(serverDone)
		conn := NewConn(NewStreamTransport(serverConn), &Options{BootstrapClient: parent})
		defer conn.Close()
		<-ctx.Done()
	}()

	trans := &countingTransport{Transport: NewStreamTransport(clientConn)}
	conn := NewConn(trans, nil)
	defer conn.Close()
	client := conn.Bootstrap(ctx)
	defer client.Release()
	if err := client.Resolve(ctx); err != nil {
		t.Fatal(err)
	}

	gate := &blockingGateController{
		waitEntered: make(chan struct{}, 1),
		grant:       make(chan struct{}),
		completions: make(chan flowcontrol.MessageOutcomeKind, 4),
		poisons:     make(chan error, 1),
	}
	client.SetFlowLimiter(&blockingGateLimiter{gate: gate})

	parentAnswer, releaseParent := client.SendCall(ctx, capnp.Send{
		Method: capnp.Method{InterfaceID: parentInterface, MethodID: parentMethod},
	})
	defer releaseParent()
	before := trans.count()

	placed := make(chan struct{})
	type pipelineResult struct {
		answer  *capnp.Answer
		release capnp.ReleaseFunc
	}
	pipelineDone := make(chan pipelineResult, 1)
	go func() {
		answer, release := parentAnswer.PipelineSend(ctx, []capnp.PipelineOp{{Field: 0}}, capnp.Send{
			Method: capnp.Method{InterfaceID: childInterface, MethodID: childMethod},
			PlaceArgs: func(capnp.Struct) error {
				close(placed)
				return nil
			},
		})
		pipelineDone <- pipelineResult{answer: answer, release: release}
	}()
	<-gate.waitEntered
	select {
	case <-placed:
		t.Fatal("pipelined RPC PlaceArgs ran behind a closed GateNext permission")
	default:
	}
	if got := trans.count(); got != before {
		t.Fatalf("transport NewMessage count = %d; want %d while pipeline gate is closed", got, before)
	}

	close(gate.grant)
	pipelined := <-pipelineDone
	defer pipelined.release()
	<-placed
	if got := trans.count(); got != before+1 {
		t.Fatalf("transport NewMessage count = %d; want %d after pipeline gate opened", got, before+1)
	}

	close(unblockParent)
	if _, err := parentAnswer.Struct(); err != nil {
		t.Fatal(err)
	}
	if _, err := pipelined.answer.Struct(); err != nil {
		t.Fatal(err)
	}
	for range 2 {
		if kind := <-gate.completions; kind != flowcontrol.MessageOutcomeSucceeded {
			t.Fatalf("terminal outcome = %v; want succeeded", kind)
		}
	}
	select {
	case kind := <-gate.completions:
		t.Fatalf("duplicate terminal completion %v", kind)
	default:
	}
	select {
	case err := <-gate.poisons:
		t.Fatalf("unexpected controller poison: %v", err)
	default:
	}
	cancel()
	<-serverDone
}

func TestGateNextRPCPrecommitCancellationDoesNotPoisonStream(t *testing.T) {
	clientConn, serverConn := net.Pipe()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	unblockServer := make(chan struct{})
	serverDone := make(chan struct{})
	go func() {
		defer close(serverDone)
		bootstrap := testcapnp.StreamTest_ServerToClient(gatedStreamTestServer{unblock: unblockServer})
		conn := NewConn(NewStreamTransport(serverConn), &Options{BootstrapClient: capnp.Client(bootstrap)})
		defer conn.Close()
		<-ctx.Done()
	}()

	trans := &countingTransport{Transport: NewStreamTransport(clientConn)}
	conn := NewConn(trans, nil)
	defer conn.Close()
	client := testcapnp.StreamTest(conn.Bootstrap(ctx))
	defer client.Release()
	if err := client.Resolve(ctx); err != nil {
		t.Fatal(err)
	}

	gate := &blockingGateController{
		waitEntered: make(chan struct{}, 2),
		grant:       make(chan struct{}),
		completions: make(chan flowcontrol.MessageOutcomeKind, 4),
		poisons:     make(chan error, 1),
	}
	capnp.Client(client).SetFlowLimiter(&blockingGateLimiter{gate: gate})
	if err := client.Push(ctx, func(p testcapnp.StreamTest_push_Params) error {
		return p.SetData([]byte{1})
	}); err != nil {
		t.Fatal(err)
	}
	before := trans.count()

	callCtx, cancelCall := context.WithCancel(ctx)
	placedCanceled := make(chan struct{})
	canceledDone := make(chan error, 1)
	go func() {
		canceledDone <- client.Push(callCtx, func(p testcapnp.StreamTest_push_Params) error {
			close(placedCanceled)
			return p.SetData([]byte{2})
		})
	}()
	<-gate.waitEntered
	cancelCall()
	if err := <-canceledDone; !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled Push = %v; want context.Canceled", err)
	}
	select {
	case <-placedCanceled:
		t.Fatal("canceled RPC call reached PlaceArgs before admission")
	default:
	}
	if got := trans.count(); got != before {
		t.Fatalf("transport NewMessage count = %d; want %d after pre-commit cancellation", got, before)
	}
	select {
	case err := <-gate.poisons:
		t.Fatalf("pre-commit cancellation poisoned controller: %v", err)
	default:
	}

	placedSuccessor := make(chan struct{})
	successorDone := make(chan error, 1)
	go func() {
		successorDone <- client.Push(ctx, func(p testcapnp.StreamTest_push_Params) error {
			close(placedSuccessor)
			return p.SetData([]byte{3})
		})
	}()
	<-gate.waitEntered
	close(gate.grant)
	if err := <-successorDone; err != nil {
		t.Fatalf("successor Push after cancellation = %v", err)
	}
	<-placedSuccessor

	close(unblockServer)
	if err := capnp.Client(client).WaitStreaming(); err != nil {
		t.Fatalf("stream poisoned by pre-commit cancellation: %v", err)
	}
	for range 2 {
		if kind := <-gate.completions; kind != flowcontrol.MessageOutcomeSucceeded {
			t.Fatalf("terminal outcome = %v; want succeeded", kind)
		}
	}
	select {
	case kind := <-gate.completions:
		t.Fatalf("canceled call produced terminal completion %v", kind)
	default:
	}
	cancel()
	<-serverDone
}

func TestGateNextRPCApplicationExceptionIsSuccessfulTransportOutcome(t *testing.T) {
	const (
		interfaceID = 0xf003
		methodID    = 3
	)
	clientConn, serverConn := net.Pipe()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	applicationErr := errors.New("application failure")
	bootstrap := capnp.NewClient(serverpkg.New([]serverpkg.Method{{
		Method: capnp.Method{InterfaceID: interfaceID, MethodID: methodID},
		Impl: func(context.Context, *serverpkg.Call) error {
			return applicationErr
		},
	}}, nil, nil))
	defer bootstrap.Release()

	serverDone := make(chan struct{})
	go func() {
		defer close(serverDone)
		conn := NewConn(NewStreamTransport(serverConn), &Options{BootstrapClient: bootstrap})
		defer conn.Close()
		<-ctx.Done()
	}()
	conn := NewConn(NewStreamTransport(clientConn), nil)
	defer conn.Close()
	client := conn.Bootstrap(ctx)
	defer client.Release()
	if err := client.Resolve(ctx); err != nil {
		t.Fatal(err)
	}

	gate := &blockingGateController{
		waitEntered: make(chan struct{}, 1),
		grant:       make(chan struct{}),
		completions: make(chan flowcontrol.MessageOutcomeKind, 2),
		poisons:     make(chan error, 1),
	}
	client.SetFlowLimiter(&blockingGateLimiter{gate: gate})
	answer, release := client.SendCall(ctx, capnp.Send{
		Method: capnp.Method{InterfaceID: interfaceID, MethodID: methodID},
	})
	defer release()
	if _, err := answer.Struct(); err == nil {
		t.Fatal("application call unexpectedly succeeded")
	}
	if kind := <-gate.completions; kind != flowcontrol.MessageOutcomeSucceeded {
		t.Fatalf("terminal outcome = %v; want succeeded transport", kind)
	}
	select {
	case kind := <-gate.completions:
		t.Fatalf("duplicate terminal completion %v", kind)
	default:
	}
	select {
	case err := <-gate.poisons:
		t.Fatalf("application exception poisoned controller: %v", err)
	default:
	}
	cancel()
	<-serverDone
}

type gatedStreamTestServer struct {
	unblock chan struct{}
}

func (s gatedStreamTestServer) Push(ctx context.Context, p testcapnp.StreamTest_push) error {
	p.Go()
	<-s.unblock
	return nil
}

type slowStreamTestServer struct{}

func (slowStreamTestServer) Push(ctx context.Context, p testcapnp.StreamTest_push) error {
	p.Go()
	// Take a while processing this, so calls can build up.
	time.Sleep(200 * time.Millisecond)
	return nil
}
