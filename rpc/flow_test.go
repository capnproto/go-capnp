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
