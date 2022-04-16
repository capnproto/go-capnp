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
	"capnproto.org/go/capnp/v3/server"
	rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"
)

// measureTransport is a wrapper around another transport, and measures the
// total size of all messages received from RecvMessage(), but not released.
// It tracks the current and all-time-maximum of this value.
type measuringTransport struct {
	Transport

	mu              sync.Mutex
	inUse, maxInUse uint64
}

func (t *measuringTransport) RecvMessage(ctx context.Context) (rpccp.Message, capnp.ReleaseFunc, error) {
	msg, release, err := t.Transport.RecvMessage(ctx)
	if err != nil {
		return msg, release, err
	}

	size, err := msg.Struct.Message().TotalSize()
	if err != nil {
		return msg, release, err
	}

	t.mu.Lock()
	t.inUse += size
	if t.inUse > t.maxInUse {
		t.maxInUse = t.inUse
	}
	t.mu.Unlock()

	oldRelease := release
	release = capnp.ReleaseFunc(func() {
		oldRelease()
		t.mu.Lock()
		defer t.mu.Unlock()
		t.inUse -= size
	})
	return msg, release, err
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

	limit := uint64(1 << 20)

	clientConn, serverConn := net.Pipe()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	serverTrans := &measuringTransport{Transport: NewStreamTransport(serverConn)}
	done := make(chan struct{})
	go func() {
		// Server

		bootstrap := testcapnp.StreamTest_ServerToClient(
			slowStreamTestServer{},
			&server.Policy{
				// Crank this way up, so we don't block on the server side.
				// Exact value is somewhat arbitrary, but big enough that
				// we should always hit the flow limit *long* before we start
				// blocking server side.
				//
				// In the long term, this will be unbounded anyway, but
				// for now we have to do this.
				MaxConcurrentCalls: int(limit * 3),
			})
		conn := NewConn(serverTrans, &Options{
			BootstrapClient: bootstrap.Client,
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

		client := testcapnp.StreamTest{Client: conn.Bootstrap(ctx)}
		defer client.Release()

		// Make a decently sized payload, so we can expect the size of the
		// parameters to dominate the size of rpc messages:
		data := make([]byte, 2048)

		// Rig up the flow control, then send calls as fast as the limiter will
		// let us:
		client.Client.SetFlowLimiter(flowcontrol.NewFixedLimiter(limit))
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
	p.Ack()
	// Take a while processing this, so calls can build up
	time.Sleep(200 * time.Millisecond)
	return nil
}
