package rpc

// tests for streaming/flow control

import (
	"context"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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

	return releaseHook{
		t:               t,
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
	rh.t.mu.Lock()
}

// Test that attaching a fixed-size FlowLimiter results in actually limiting the
// flow of messages. We do this by spawning a server that responds to calls slowly,
// and then send calls as fast as the limiter will let us. At the end, we check
// that the maximum size of outstanding messages looks right.
func TestFixedFlowLimit(t *testing.T) {
	if os.Getenv("FLAKY_TESTS") != "1" {
		t.Skip("Not running TestFlowFixedLimit, which is flaky. Set FLAKY_TESTS=1 to enable")
		// TODO: at some point make this test run robustly in CI;
		// it seems to work reliably when run locally for both @zenhack
		// and @lthibault
	}
	t.Parallel()

	limit := int64(1 << 20)

	clientConn, serverConn := net.Pipe()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	serverTrans := &measuringTransport{Transport: NewStreamTransport(serverConn)}
	done := make(chan struct{})
	go func() {
		// Server

		bootstrap := testcapnp.StreamTest_ServerToClient(slowStreamTestServer{})
		conn := NewConn(serverTrans, &Options{
			BootstrapClient: capnp.Client(bootstrap),
		})
		defer conn.Close()
		<-ctx.Done()
		close(done)
	}()

	func() {
		// Client

		trans := NewStreamTransport(clientConn)
		conn := NewConn(trans, nil)
		defer conn.Close()

		client := testcapnp.StreamTest(conn.Bootstrap(ctx))
		defer client.Release()

		// Make a decently sized payload, so we can expect the size of the
		// parameters to dominate the size of rpc messages:
		data := make([]byte, 2048)

		// Rig up the flow control, then send calls as fast as the limiter will
		// let us:
		capnp.Client(client).SetFlowLimiter(flowcontrol.NewFixedLimiter(limit))
		for ctx.Err() == nil {
			client.Push(ctx, func(p testcapnp.StreamTest_push_Params) error {
				return p.SetData(data)
			})
		}
	}()

	<-done
	// The server has exited, it is now safe to access the transport without syncrhonization.
	//
	// We check that the max bytes in filght never exceeded the limit by more than 10%.
	// We leave a little headroom to allow for things like:
	//
	// * space taken up by non-calls, like Finish messages and such
	// * descrepencies due to the rpc system choosing to add metadata after the
	//   size of the call has already been measured.
	//
	// Note that this is theoretically a per-client limit, not a per-connection limit,
	// but this test only has one client. TODO: we could enhance this by having a test
	// with multiple clients, and examining the sum of their uses.
	assert.Less(t, serverTrans.maxInUse, limit+(limit/10),
		"Flow control didn't limit flow enough")

	// Let's also make sure that we aren't too far *under* the limit; if we've written the test
	// right, the client should be much faster than the server, so we should get close to it.
	assert.Greater(t, serverTrans.maxInUse, limit-(limit/10),
		"To little flow; something is wrong.")
}

type slowStreamTestServer struct{}

func (slowStreamTestServer) Push(ctx context.Context, p testcapnp.StreamTest_push) error {
	p.Go()
	// Take a while processing this, so calls can build up
	time.Sleep(200 * time.Millisecond)
	return nil
}
