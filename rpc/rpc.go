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
	recvDone   <-chan struct{}
	recvCloser interface {
		CloseRecv() error
	}

	// send is protected by sending, a condition channel.
	// sending is non-nil if a send.NewMessage is being attempted, and
	// the channel is closed when the message is finished.
	sending <-chan struct{}
	send    Sender

	nextQuestionID uint32
}

// Options specifies optional parameters for creating a Conn.
type Options struct {
}

// NewConn creates a new connection that communications on a given
// transport.  Closing the connection will close the transport.
// Passing nil for opts is the same as passing the zero value.
func NewConn(t Transport, opts *Options) *Conn {
	done := make(chan struct{})
	c := &Conn{
		send:       t,
		recvCloser: t,
		recvDone:   done,
	}
	go func() {
		defer close(done)
		c.receive(context.Background(), t)
	}()
	return c
}

// newMessage creates a new outgoing message.  The caller must be
// holding onto c.mu, but if newMessage does not return an error, then
// c.mu will be released.  Once send or cancel is called, then c.mu
// will be reacquired.
func (c *Conn) newMessage(ctx context.Context) (_ rpccp.Message, send func() error, cancel func(), _ error) {
	for {
		if c.closed {
			// TODO(someday): classify as disconnected
			return rpccp.Message{}, nil, nil, errors.New("connection closed")
		}
		s := c.sending
		if s == nil {
			break
		}
		c.mu.Unlock()
		select {
		case <-s:
		case <-ctx.Done():
			c.mu.Lock()
			return rpccp.Message{}, nil, nil, ctx.Err()
		}
		c.mu.Lock()
	}
	sending := make(chan struct{})
	c.sending = sending
	c.mu.Unlock()
	msg, tsend, tcancel, err := c.send.NewMessage(ctx)
	if err != nil {
		c.mu.Lock()
		close(sending)
		c.sending = nil
		return rpccp.Message{}, nil, nil, err
	}
	return msg, func() error {
			err := tsend()
			c.mu.Lock()
			close(sending)
			c.sending = nil
			return err
		}, func() {
			tcancel()
			c.mu.Lock()
			close(sending)
			c.sending = nil
		}, nil
}

// Bootstrap returns the remote vat's bootstrap interface.
func (c *Conn) Bootstrap(ctx context.Context) *capnp.Client {
	defer c.mu.Unlock()
	c.mu.Lock()
	id := c.nextQuestionID
	c.nextQuestionID++
	msg, send, cancel, err := c.newMessage(ctx)
	if err != nil {
		return capnp.ErrorClient(fmt.Errorf("rpc bootstrap: create message: %v", err))
	}
	boot, err := msg.NewBootstrap()
	if err != nil {
		cancel()
		return capnp.ErrorClient(fmt.Errorf("rpc bootstrap: create message: %v", err))
	}
	boot.SetQuestionId(id)
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
	rerr := c.recvCloser.CloseRecv() // will wait on recvDone at the end

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

	<-c.recvDone
	if rerr != nil {
		return fmt.Errorf("rpc: close transport: %v", rerr)
	}
	if serr != nil {
		return fmt.Errorf("rpc: close transport: %v", serr)
	}
	return nil
}

// receive receives and dispatches messages coming from r.
// It is intended to run in its own goroutine.
func (c *Conn) receive(ctx context.Context, r Receiver) {
	for {
		recv, releaseRecv, err := r.RecvMessage(ctx)
		if err != nil {
			// TODO(soon): log error
			return
		}
		switch recv.Which() {
		case rpccp.Message_Which_unimplemented:
			// no-op for now to avoid feedback loop
		case rpccp.Message_Which_return:
			ret, err := recv.Return()
			if err != nil {
				// TODO(soon): log error
				continue
			}
			err = c.handleReturn(ctx, ret, releaseRecv)
			if err != nil {
				// TODO(soon): log error
				continue
			}
		default:
			err := c.handleUnimplemented(ctx, recv)
			releaseRecv()
			if err != nil {
				// TODO(soon): log error
				continue
			}
		}
	}
}

func (c *Conn) handleReturn(ctx context.Context, ret rpccp.Return, releaseRet capnp.ReleaseFunc) error {
	defer releaseRet()
	defer c.mu.Unlock()
	c.mu.Lock()
	// TODO(soon): disconnect if return ID not in questions table.
	msg, send, cancel, err := c.newMessage(ctx)
	if err != nil {
		return err
	}
	fin, err := msg.NewFinish()
	if err != nil {
		cancel()
		return err
	}
	fin.SetQuestionId(ret.AnswerId())
	fin.SetReleaseResultCaps(false)
	if err := send(); err != nil {
		return err
	}
	return nil
}

func (c *Conn) handleUnimplemented(ctx context.Context, recv rpccp.Message) error {
	defer c.mu.Unlock()
	c.mu.Lock()
	msg, send, cancel, err := c.newMessage(ctx)
	if err != nil {
		return err
	}
	if err := msg.SetUnimplemented(recv); err != nil {
		cancel()
		return nil
	}
	err = send()
	return err
}
