package rpc

import (
	"context"
	"fmt"
	"net"
	"strings"
)

// Serve serves a Cap'n Proto RPC to incoming connections
// Serve exits with the listener error
func Serve(lis net.Listener, options *Options) error {
	if options == nil { //|| !options.BootstrapClient.IsValid() {
		return fmt.Errorf("invalid options")
	}
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
			conn := NewConn(transport, options)

			select {
			case <-conn.Done():
				// Remote client connection closed
				return
			}
		}()
	}
}

// ListenAndServe opens a listener on the given address and serves a Cap'n Proto RPC to incoming connections
// If address starts with "unix:" it is considered a Unix Domain Socket path, otherwise a TCP address.
// Context can be used to stop listening.
func ListenAndServe(ctx context.Context, address string, options *Options) error {
	var listener net.Listener
	var err error

	if address == "" {
		return fmt.Errorf("missing address")
	}
	// UDS paths start with either '.' or '/'
	if strings.HasPrefix(address, "unix:") {
		listener, err = net.Listen("unix", address[5:])
	} else {
		listener, err = net.Listen("tcp", address)
	}
	if err == nil {
		// to close this listener, close the context
		go func() {
			select {
			case <-ctx.Done():
				_ = listener.Close()
				return
			}
		}()
		err = Serve(listener, options)
	}
	return err
}
