package rpc_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc"
	testcp "capnproto.org/go/capnp/v3/rpc/internal/testcapnp"
)

func TestServe(t *testing.T) {
	t.Parallel()
	t.Log("Opening listener")
	lis, err := net.Listen("tcp", ":0")
	defer lis.Close()
	assert.NoError(t, err)

	errChannel := make(chan error)
	go func() {
		err2 := rpc.Serve(lis, nil)
		t.Log("Serve has ended")
		errChannel <- err2
	}()

	time.Sleep(time.Second)

	t.Log("Closing listener")
	err = lis.Close()
	assert.NoError(t, err)
	select {
	case <-time.After(time.Second * 2):
		t.Fail()
	case err = <-errChannel:
		assert.ErrorIs(t, err, net.ErrClosed)
	}
}

// TestServeCapability serves the ping pong capability and tests
// if a client can successfuly receive served data.
func TestServeCapability(t *testing.T) {

	t.Parallel()
	ctx := context.Background()
	t.Log("Opening listener")
	lis, err := net.Listen("tcp", ":0")
	defer lis.Close()
	assert.NoError(t, err)

	srv := testcp.PingPong_ServerToClient(pingPongServer{})
	opts := &rpc.Options{
		BootstrapClient: capnp.Client(srv),
	}
	errChannel := make(chan error)
	go func() {
		err2 := rpc.Serve(lis, opts)
		t.Log("Serve has ended")
		errChannel <- err2
	}()

	// connect to the server and invoke the magic N method
	addr := lis.Addr().String()
	conn, err := net.Dial("tcp", addr)
	assert.NoError(t, err)
	transport := rpc.NewStreamTransport(conn)
	rpcConn := rpc.NewConn(transport, nil)
	ppClient := testcp.PingPong(rpcConn.Bootstrap(ctx))
	method, release := ppClient.EchoNum(ctx, func(ps testcp.PingPong_echoNum_Params) error {
		ps.SetN(42)
		return nil
	})
	defer release()
	resp, err := method.Struct()
	assert.NoError(t, err)
	numberN := resp.N()
	assert.Equal(t, int64(42), numberN)
	t.Logf("Received pingpong: N=%d", numberN)
	err = lis.Close()
	assert.NoError(t, err)
}

func TestListenAndServe(t *testing.T) {
	var err error
	t.Parallel()
	ctx, cancelFunc := context.WithCancel(context.Background())
	errChannel := make(chan error)

	go func() {
		t.Log("Starting ListenAndServe")
		err2 := rpc.ListenAndServe(ctx, "tcp", ":0", nil)
		errChannel <- err2
		t.Log("ListenAndServe has ended")
	}()

	time.Sleep(time.Second)
	t.Log("Closing context")

	cancelFunc()
	select {
	case <-time.After(time.Second * 2):
		t.Error("Cancelling context didn't end the listener")
	case err = <-errChannel:
		assert.ErrorIs(t, err, net.ErrClosed)
	}
}
