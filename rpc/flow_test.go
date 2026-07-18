package rpc

// tests for streaming/flow control

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/flowcontrol"
	"capnproto.org/go/capnp/v3/rpc/internal/testcapnp"
	"capnproto.org/go/capnp/v3/rpc/transport"
)

// measureTransport is a wrapper around another transport, and measures the
// total size of all messages received from RecvMessage(), but not released.
// It tracks the current and all-time-maximum of this value.
type measuringTransport struct {
	Transport

	mu              sync.Mutex
	inUse, maxInUse uint64
	changed         chan struct{}
}

func (t *measuringTransport) RecvMessage() (transport.IncomingMessage, error) {
	inMsg, err := t.Transport.RecvMessage()
	if err != nil {
		return inMsg, err
	}

	size, err := inMsg.Message().Message().TotalSize()
	if err != nil {
		return inMsg, err
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	if t.inUse += size; t.inUse > t.maxInUse {
		t.maxInUse = t.inUse
	}
	select {
	case t.changed <- struct{}{}:
	default:
	}

	return releaseHook{
		t:               t,
		size:            size,
		IncomingMessage: inMsg,
	}, nil
}

type releaseHook struct {
	t    *measuringTransport
	size uint64
	transport.IncomingMessage
}

func (rh releaseHook) Release() {
	rh.IncomingMessage.Release()

	rh.t.mu.Lock()
	rh.t.inUse -= rh.size
	rh.t.mu.Unlock()
}

func (t *measuringTransport) waitForInUse(ctx context.Context, min uint64) bool {
	for {
		t.mu.Lock()
		inUse := t.inUse
		t.mu.Unlock()
		if inUse >= min {
			return true
		}

		select {
		case <-ctx.Done():
			return false
		case <-t.changed:
		}
	}
}

func (t *measuringTransport) maxInUseSnapshot() uint64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.maxInUse
}

// Test that attaching a fixed-size FlowLimiter results in actually limiting the
// flow of messages. The server holds accepted calls until the test releases it,
// which fills the client's limiter without depending on scheduler timing.
func TestFixedFlowLimit(t *testing.T) {
	t.Parallel()

	limit := int64(1 << 20)
	const callCount = 1024

	clientConn, serverConn := net.Pipe()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	serverTrans := &measuringTransport{
		Transport: NewStreamTransport(serverConn),
		changed:   make(chan struct{}, 1),
	}
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
		conn := NewConn(serverTrans, &Options{
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
	capnp.Client(client).SetFlowLimiter(flowcontrol.NewFixedLimiter(limit))

	sent := make(chan error, 1)
	go func() {
		for i := 0; i < callCount; i++ {
			if err := client.Push(ctx, func(p testcapnp.StreamTest_push_Params) error {
				return p.SetData(data)
			}); err != nil {
				sent <- err
				return
			}
		}
		sent <- nil
	}()

	if !serverTrans.waitForInUse(ctx, uint64(limit-(limit/10))) {
		t.Fatal("server did not receive enough held calls to exercise the limiter")
	}
	select {
	case err := <-sent:
		t.Fatalf("all calls were sent while the server held responses: %v", err)
	default:
	}

	// The server cannot respond until this gate is opened. Once it is, each
	// response returns capacity to the limiter and lets the client proceed.
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

	maxInUse := serverTrans.maxInUseSnapshot()
	// We check that the max bytes in flight never exceeded the limit by more than 10%.
	// We leave a little headroom to allow for things like:
	//
	// * space taken up by non-calls, like Finish messages and such
	// * descrepencies due to the rpc system choosing to add metadata after the
	//   size of the call has already been measured.
	//
	// Note that this is theoretically a per-client limit, not a per-connection limit,
	// but this test only has one client. TODO: we could enhance this by having a test
	// with multiple clients, and examining the sum of their uses.
	if maxInUse >= uint64(limit+(limit/10)) {
		t.Fatalf("flow control didn't limit flow enough: max in use = %d, limit = %d", maxInUse, limit)
	}

	// Let's also make sure that we aren't too far *under* the limit; if we've written the test
	// right, the client should be much faster than the server, so we should get close to it.
	if maxInUse <= uint64(limit-(limit/10)) {
		t.Fatalf("too little flow: max in use = %d, limit = %d", maxInUse, limit)
	}
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
