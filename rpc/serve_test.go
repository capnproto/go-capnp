package rpc_test

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc"
	testcp "capnproto.org/go/capnp/v3/rpc/internal/testcapnp"
)

// Test connect/disconnect to a pingpong capability
func TestServe(t *testing.T) {
	t.Parallel()
	errChannel := make(chan error)

	t.Log("Opening listener")
	lis, err := net.Listen("tcp", ":0")
	assert.NoError(t, err)
	srv := testcp.PingPong_ServerToClient(pingPongServer{})
	bootstrapClient := capnp.Client(srv)

	go func() {
		err2 := rpc.Serve(lis, bootstrapClient)
		t.Log("Serve has ended")
		errChannel <- err2
	}()

	// Create the pingpong client using the server address and close it
	addr := lis.Addr().String()
	conn, err := net.Dial("tcp", addr)
	assert.NoError(t, err)
	transport := rpc.NewStreamTransport(conn)
	rpcConn := rpc.NewConn(transport, nil)
	err = rpcConn.Close()
	assert.NoError(t, err)

	// repeat to ensure that a second connection is allowed and doesn't mess
	// with releasing the bootstrap reference counting.
	conn, err = net.Dial("tcp", addr)
	assert.NoError(t, err)
	transport = rpc.NewStreamTransport(conn)
	rpcConn = rpc.NewConn(transport, nil)
	err = rpcConn.Close()
	assert.NoError(t, err)

	t.Log("Closing server listener")
	err = lis.Close()
	assert.NoError(t, err)
	err = <-errChannel // Will hang if the server does not return.
	// Expect that the server finished with the 'connection closed' error.
	assert.ErrorIs(t, err, net.ErrClosed)
	// Check that the bootstrap client was released by Serve.
	assert.False(t, bootstrapClient.IsValid(), "Serve did not release its bootstrap client")
}

// TestServeCapability serves the ping pong capability and tests
// if a client can successfully receive served data.
func TestServeCapability(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	// Start the pingpong server
	t.Log("Opening listener")
	lis, err := net.Listen("tcp", ":0")
	assert.NoError(t, err)
	srv := testcp.PingPong_ServerToClient(pingPongServer{})
	bootstrapClient := capnp.Client(srv)

	errChannel := make(chan error)
	go func() {
		err2 := rpc.Serve(lis, bootstrapClient)
		t.Log("Serve has ended")
		errChannel <- err2
	}()

	// Create the pingpong client using the server address
	addr := lis.Addr().String()
	conn, err := net.Dial("tcp", addr)
	assert.NoError(t, err)
	transport := rpc.NewStreamTransport(conn)
	rpcConn := rpc.NewConn(transport, nil)
	defer rpcConn.Close()
	ppClient := testcp.PingPong(rpcConn.Bootstrap(ctx))
	defer ppClient.Release()

	// Invoke the magic N method. If Serve works this should provide the magic number.
	method, releaseMethod := ppClient.EchoNum(ctx, func(ps testcp.PingPong_echoNum_Params) error {
		ps.SetN(42)
		return nil
	})
	defer releaseMethod()
	resp, err := method.Struct()
	assert.NoError(t, err)
	numberN := resp.N()
	assert.Equal(t, int64(42), numberN)
	t.Logf("Received pingpong: N=%d", numberN)

	// shutdown the server and verify that Serve exits
	err = lis.Close()
	assert.NoError(t, err)

	err = <-errChannel // Will hang if the sever does not return
	assert.ErrorIs(t, err, net.ErrClosed)

	assert.False(t, bootstrapClient.IsValid(), "server bootstrap client not released")
}

func TestListenAndServe(t *testing.T) {
	var err error
	t.Parallel()
	ctx, cancelFunc := context.WithCancel(context.Background())
	errChannel := make(chan error)

	// Provide a server that listens
	srv := testcp.PingPong_ServerToClient(pingPongServer{})
	bootstrapClient := capnp.Client(srv)
	go func() {
		t.Log("Starting ListenAndServe")
		err2 := rpc.ListenAndServe(ctx, "tcp", ":0", bootstrapClient)
		errChannel <- err2
	}()

	cancelFunc()
	err = <-errChannel // Will hang if server does not return.
	assert.ErrorIs(t, err, net.ErrClosed)
}
