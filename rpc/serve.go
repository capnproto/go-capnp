package rpc

import (
	"context"
	"errors"
	"net"

	"capnproto.org/go/capnp/v3"
)

// Serve serves a Cap'n Proto RPC to incoming connections.
//
// Serve will take ownership of bootstrapClient and release it after the listener closes.
//
// Serve exits with the listener error if the listener is closed by the owner.
func Serve(lis net.Listener, bootstrapClient capnp.Client) error {
	if !bootstrapClient.IsValid() {
		err := errors.New("BootstrapClient is not valid")
		return err
	}
	// Accept incoming connections
	defer bootstrapClient.Release()
	for {
		conn, err := lis.Accept()
		if err != nil {
			// Since we took ownership of the bootstrap client, release it after we're done.
			bootstrapClient.Release()
			return err
		}

		// the RPC connection takes ownership of the bootstrap interface and will release it when the connection
		// exits, so use AddRef to avoid releasing the provided bootstrap client capability.
		opts := Options{
			BootstrapClient: bootstrapClient.AddRef(),
		}
		// For each new incoming connection, create a new RPC transport connection that will serve incoming RPC requests
		transport := NewStreamTransport(conn)
		rpcConn := NewConn(transport, &opts)
		_ = rpcConn
	}
}

// ListenAndServe opens a listener on the given address and serves a Cap'n Proto RPC to incoming connections
//
// network and address are passed to net.Listen. Use network "unix" for Unix Domain Sockets
// and "tcp" for regular TCP IP4 or IP6 connections.
//
// ListenAndServe will take ownership of bootstrapClient and release it on exit.
func ListenAndServe(ctx context.Context, network, addr string, bootstrapClient capnp.Client) error {

	listener, err := net.Listen(network, addr)

	if err == nil {
		// to close this listener, close the context
		go func() {
			<-ctx.Done()
			_ = listener.Close()
		}()
		err = Serve(listener, bootstrapClient)
	}
	return err
}
