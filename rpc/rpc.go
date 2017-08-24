// Package rpc implements the Cap'n Proto RPC protocol.
package rpc // import "zombiezen.com/go/capnproto2/rpc"

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"zombiezen.com/go/capnproto2"
	rpccp "zombiezen.com/go/capnproto2/std/capnp/rpc"
)

// A Conn is a connection to another Cap'n Proto vat.
// It is safe to use from multiple goroutines.
type Conn struct {
	// mu protects all fields in the Conn.  However, mu should not be held
	// while making calls that take indeterminate time (I/O or application
	// code).  Condition channels protect operations on any field that
	// take an indeterminate amount of time.  Thus, critical sections
	// involving mu are quite short, while still ensuring mutually
	// exclusive access to resources.
	mu sync.Mutex

	closed     bool
	recvCloser interface {
		CloseRecv() error
	}

	// send is protected by sending, a condition channel.
	// sending is non-nil if a send.NewMessage is being attempted, and
	// the channel is closed when the message is finished.
	sending <-chan struct{}
	send    Sender
}

// Options specifies optional parameters for creating a Conn.
type Options struct {
}

// NewConn creates a new connection that communications on a given
// transport.  Closing the connection will close the transport.
// Passing nil for opts is the same as passing the zero value.
func NewConn(t Transport, opts *Options) *Conn {
	return &Conn{
		send:       t,
		recvCloser: t,
	}
}

// Bootstrap returns the remote vat's bootstrap interface.
func (c *Conn) Bootstrap(ctx context.Context) *capnp.Client {
	c.mu.Lock()
	for {
		if c.closed {
			c.mu.Unlock()
			// TODO(someday): classify as disconnected
			return capnp.ErrorClient(errors.New("rpc bootstrap: connection closed"))
		}
		s := c.sending
		if s == nil {
			break
		}
		c.mu.Unlock()
		select {
		case <-s:
		case <-ctx.Done():
			return capnp.ErrorClient(fmt.Errorf("rpc bootstrap: %v", ctx.Err()))
		}
		c.mu.Lock()
	}
	sending := make(chan struct{})
	c.sending = sending
	c.mu.Unlock()
	defer func() {
		c.mu.Lock()
		close(sending)
		c.sending = nil
		c.mu.Unlock()
	}()

	msg, send, cancel, err := c.send.NewMessage(ctx)
	if err != nil {
		return capnp.ErrorClient(fmt.Errorf("rpc bootstrap: create message: %v", err))
	}
	boot, err := msg.NewBootstrap()
	if err != nil {
		cancel()
		return capnp.ErrorClient(fmt.Errorf("rpc bootstrap: create message: %v", err))
	}
	// TODO(soon): allocate an ID
	boot.SetQuestionId(0)
	if err := send(); err != nil {
		return capnp.ErrorClient(fmt.Errorf("rpc bootstrap: send message: %v", err))
	}
	return nil
}

// Close sends an abort to the remote vat and closes the underlying
// transport.
func (c *Conn) Close() error {
	// Mark closed and stop receiving messages.
	c.mu.Lock()
	c.closed = true
	c.mu.Unlock()
	rerr := c.recvCloser.CloseRecv()

	// Close Sender after all sends are finished.
	c.mu.Lock()
	for {
		s := c.sending
		if s == nil {
			break
		}
		c.mu.Unlock()
		<-s
		c.mu.Lock()
	}
	c.mu.Unlock()

	// Send abort message (ignoring error).
	{
		msg, send, cancel, err := c.send.NewMessage(context.Background())
		if err != nil {
			goto closeSend
		}
		abort, err := msg.NewAbort()
		if err != nil {
			cancel()
			goto closeSend
		}
		// TODO(soon): allocate an ID
		abort.SetType(rpccp.Exception_Type_failed)
		if err := abort.SetReason("connection closed"); err != nil {
			cancel()
			goto closeSend
		}
		send()
	}
closeSend:
	serr := c.send.CloseSend()

	if rerr != nil {
		return fmt.Errorf("rpc: close transport: %v", rerr)
	}
	if serr != nil {
		return fmt.Errorf("rpc: close transport: %v", serr)
	}
	return nil
}
