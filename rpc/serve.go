package rpc

import (
	"context"
	"net"
)

// Serve serves a Cap'n Proto RPC to incoming connections
// Serve exits with the listener error
func Serve(lis net.Listener, opt *Options) error {
	// Accept incoming connections
	for {
		rwc, err := lis.Accept()
		if err != nil {
			return err
		}
		// For each new incoming connection, create a new RPC transport connection that will serve incoming RPC requests
		// rpc.Options will contain the bootstrap capability
		go func() {
			transport := NewStreamTransport(rwc)
			conn := NewConn(transport, opt)

			<-conn.Done()
			// Remote client connection closed
			return
		}()
	}
}

// ListenAndServe opens a listener on the given address and serves a Cap'n Proto RPC to incoming connections
// network and address are passed to net.Listen. Use network "unix" for Unix Domain Sockets
// and "tcp" for regular TCP connections.
func ListenAndServe(ctx context.Context, network, addr string, opt *Options) error {
	//var listener net.Listener
	//var err error

	//listener, err = net.Listen(network, address)
	listener, err := new(net.ListenConfig).Listen(ctx, network, addr)

	if err == nil {
		// to close this listener, close the context
		go func() {
			<-ctx.Done()
			_ = listener.Close()
		}()
		err = Serve(listener, opt)
	}
	return err
}
