package rpc

import (
	"context"
	"io"
	"sync/atomic"
	"time"

	capnp "zombiezen.com/go/capnproto2"
	"zombiezen.com/go/capnproto2/internal/errors"
	rpccp "zombiezen.com/go/capnproto2/std/capnp/rpc"
)

// A Transport sends and receives Cap'n Proto RPC messages to and from
// another vat.
//
// It is safe to call NewMessage and its returned functions concurrently
// with RecvMessage.
type Transport interface {
	// NewMessage allocates a new message to be sent over the transport.
	// The caller must call the release function when it no longer needs
	// to reference the message.  Before releasing the message, send may be
	// called at most once to send the mssage, taking its cancelation and
	// deadline from ctx.
	//
	// Messages returned by NewMessage must have a nil CapTable.
	// The caller may modify the CapTable before sending it, but the
	// message's CapTable must be nil before it is sent or released.
	//
	// The Arena in the returned message should be fast at allocating new
	// segments.
	NewMessage(ctx context.Context) (_ rpccp.Message, send func() error, _ capnp.ReleaseFunc, _ error)

	// RecvMessage receives the next message sent from the remote vat.
	// The returned message is only valid until the release function is
	// called or Close is called.  The release function may be called
	// concurrently with RecvMessage or with any other release function
	// returned by RecvMessage.
	//
	// Messages returned by RecvMessage must have a nil CapTable.
	// The caller may modify the CapTable, but the message's CapTable must
	// be nil before it is released.
	//
	// The Arena in the returned message should not fetch segments lazily;
	// the Arena should be fast to access other segments.
	RecvMessage(ctx context.Context) (rpccp.Message, capnp.ReleaseFunc, error)

	// Close releases any resources associated with the transport.  All
	// messages created with NewMessage must be released before calling
	// Close.  It is not safe to call Close concurrently with any other
	// operations on the transport.
	Close() error
}

// A transport serializes and deserializes Cap'n Proto using a Codec.
// It adds no buffering beyond what is provided by the underlying
// byte transfer mechanism.
type transport struct {
	c      Codec
	closed bool
	err    errorValue
}

// NewTransport creates a new transport that uses the supplied codec
// to read and write messages across the wire.
func NewTransport(c Codec) Transport { return &transport{c: c} }

// NewStreamTransport creates a new transport that reads and writes to rwc.
// Closing the transport will close rwc.
//
// If rwc has SetReadDeadline or SetWriteDeadline methods, they will be
// used to handle Context cancellation and deadlines.  If rwc does not
// have these methods, then rwc.Close must be safe to call concurrently
// with rwc.Read.  Notably, this is not true of *os.File before Go 1.9
// (see https://golang.org/issue/7970).
func NewStreamTransport(rwc io.ReadWriteCloser) Transport {
	return NewTransport(newStreamCodec(rwc, basicEncoding{}))
}

// NewPackedStreamTransport creates a new transport that uses a packed
// encoding.
//
// See:  NewStreamTransport.
func NewPackedStreamTransport(rwc io.ReadWriteCloser) Transport {
	return NewTransport(newStreamCodec(rwc, packedEncoding{}))
}

// NewMessage allocates a new message to be sent.
//
// It is safe to call NewMessage concurrently with RecvMessage.
func (s *transport) NewMessage(ctx context.Context) (_ rpccp.Message, send func() error, release capnp.ReleaseFunc, _ error) {
	// Check if stream is broken
	if err := s.err.Load(); err != nil {
		return rpccp.Message{}, nil, nil, err
	}

	// TODO(soon): reuse memory
	msg, seg, err := capnp.NewMessage(capnp.MultiSegment(nil))
	if err != nil {
		return rpccp.Message{}, nil, nil, errors.New(errors.Failed, "rpc stream transport", "new message: "+err.Error())
	}
	rmsg, err := rpccp.NewRootMessage(seg)
	if err != nil {
		return rpccp.Message{}, nil, nil, errors.New(errors.Failed, "rpc stream transport", "new message: "+err.Error())
	}

	send = func() error {
		// context expired?
		if err := ctx.Err(); err != nil {
			return errors.New(errors.Failed, "rpc stream transport", "send: "+ctx.Err().Error())
		}

		// stream error?
		if err := s.err.Load(); err != nil {
			return err
		}

		// ok, go!
		if err = s.c.Encode(ctx, msg); err != nil {
			if _, ok := err.(partialWriteError); ok {
				s.err.Set(errors.New(errors.Disconnected, "rpc stream transport", "broken due to partial write"))
			}

			err = errors.New(errors.Failed, "rpc stream transport", "send: "+err.Error())
		}

		return err
	}

	return rmsg, send, func() { msg.Reset(nil) }, nil
}

// SetPartialWriteTimeout sets the timeout for completing the
// transmission of a partially sent message after the send is cancelled
// or interrupted for any future sends.  If not set, a reasonable
// non-zero value is used.
//
// Setting a shorter timeout may free up resources faster in the case of
// an unresponsive remote peer, but may also make the transport respond
// too aggressively to bursts of latency.
func (s *transport) SetPartialWriteTimeout(d time.Duration) {
	s.c.SetPartialWriteTimeout(d)
}

// RecvMessage reads the next message from the underlying reader.
//
// It is safe to call RecvMessage concurrently with NewMessage.
func (s *transport) RecvMessage(ctx context.Context) (rpccp.Message, capnp.ReleaseFunc, error) {
	if err := s.err.Load(); err != nil {
		return rpccp.Message{}, nil, err
	}

	msg, err := s.c.Decode(ctx)
	if err != nil {
		return rpccp.Message{}, nil, errors.New(errors.Failed, "rpc stream transport", "receive: "+err.Error())
	}
	rmsg, err := rpccp.ReadRootMessage(msg)
	if err != nil {
		return rpccp.Message{}, nil, errors.New(errors.Failed, "rpc stream transport", "receive: "+err.Error())
	}
	return rmsg, func() { msg.Reset(nil) }, nil
}

// Close closes the underlying ReadWriteCloser.  It is not safe to call
// Close concurrently with any other operations on the transport.
func (s *transport) Close() error {
	if s.closed {
		return errors.New(errors.Disconnected, "rpc stream transport", "already closed")
	}
	s.closed = true
	err := s.c.Close()
	if err != nil {
		return errors.New(errors.Failed, "rpc stream transport", "close: "+err.Error())
	}
	return nil
}

type partialWriteError struct{ error }

type errorValue atomic.Value

func (ev *errorValue) Load() error {
	if err := (*atomic.Value)(ev).Load(); err != nil {
		return err.(error)
	}

	return nil
}

func (ev *errorValue) Set(err error) {
	(*atomic.Value)(ev).Store(err)
}
