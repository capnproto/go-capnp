package rpc

import (
	"context"
	"io"
	"sync"
	"time"

	"zombiezen.com/go/capnproto2"
	"zombiezen.com/go/capnproto2/internal/errors"
	rpccp "zombiezen.com/go/capnproto2/std/capnp/rpc"
)

// A Sender delivers Cap'n Proto RPC messages to another vat.
type Sender interface {
	// NewMessage allocates a new message to be sent over the transport.
	// The caller must call the release function when it no longer needs
	// to reference the message.  Before release is called, send may be
	// called at most once to send the mssage, taking its cancelation and
	// deadline from ctx.
	//
	// Messages returned by NewMessage must have a nil CapTable.  release
	// must release all clients in the CapTable.
	//
	// The Arena in the returned message should be fast at allocating new
	// segments.
	NewMessage(ctx context.Context) (_ rpccp.Message, send func() error, _ capnp.ReleaseFunc, _ error)

	// CloseSend releases any resources associated with the sender.
	// All messages created with NewMessage must be released before
	// calling CloseSend.
	CloseSend() error
}

// A Receiver receives Cap'n Proto RPC messages from another vat.
type Receiver interface {
	// RecvMessage waits to receive a message and returns it.
	// The returned message is only valid until the release function is
	// called or CloseRecv is called.
	//
	// Messages returned by RecvMessage must have a nil CapTable.
	// The caller may mutate the CapTable.
	//
	// The Arena in the returned message should not fetch segments lazily;
	// the Arena should be fast to access other segments.
	RecvMessage(ctx context.Context) (rpccp.Message, capnp.ReleaseFunc, error)

	// CloseRecv releases any resources associated with the receiver and
	// ends any unfinished RecvMessage call.
	CloseRecv() error
}

// Transport is a grouping of the Sender and Receiver interfaces.
//
// It is safe to call Sender methods concurrently with Receiver methods.
type Transport interface {
	Sender
	Receiver
}

// StreamTransport serializes and deserializes unpacked Cap'n Proto
// messages on a byte stream.  StreamTransport adds no buffering beyond
// what its underlying stream has.
//
// Sender methods on StreamTransport cannot be called concurrently with
// each other and Receiver methods on StreamTransport cannot be called
// concurrently with each other.  However, it is safe to call Sender
// methods concurrently with Receiver methods.
type StreamTransport struct {
	// Send
	enc      *capnp.Encoder
	deadline writeDeadlineSetter
	// Receive
	recv Receiver
	// Close
	c  io.Closer
	cw writeCloser

	mu     sync.Mutex
	closes uint8
}

// NewStreamTransport creates a new transport that reads and writes to rwc.
// Closing the transport will close rwc.
//
// If rwc has a SetWriteDeadline method, it will be used when a message
// is sent.  If rwc has CloseRead/CloseWrite methods, those will be used
// during CloseRecv/CloseSend.  Regardless, Close will be called after
// CloseRecv and CloseSend have both been called.
func NewStreamTransport(rwc io.ReadWriteCloser) *StreamTransport {
	d, _ := rwc.(writeDeadlineSetter)
	cw, _ := rwc.(writeCloser)
	s := &StreamTransport{
		enc:      capnp.NewEncoder(rwc),
		deadline: d,
		c:        rwc,
		cw:       cw,
	}
	dec := capnp.NewDecoder(rwc)
	// TODO(someday): reuse buffer as long as release is called before next RecvMessage.
	if c, ok := rwc.(readCloser); ok {
		s.recv = closerReceiver{dec, c}
	} else {
		s.recv = signalReceiver{dec, make(chan struct{})}
	}
	return s
}

// NewMessage allocates a new message to be sent.  The send function may
// make multiple calls to Write on the underlying writer.
func (s *StreamTransport) NewMessage(ctx context.Context) (_ rpccp.Message, send func() error, release capnp.ReleaseFunc, _ error) {
	// TODO(soon): reuse memory
	msg, seg, _ := capnp.NewMessage(capnp.MultiSegment(nil))
	rmsg, _ := rpccp.NewRootMessage(seg)
	send = func() error {
		if s.deadline != nil {
			// TODO(someday): log errors
			if d, ok := ctx.Deadline(); ok {
				s.deadline.SetWriteDeadline(d)
			} else {
				s.deadline.SetWriteDeadline(time.Time{})
			}
		}
		return s.enc.Encode(msg)
	}
	release = func() {
		msg.Reset(nil)
	}
	return rmsg, send, release, nil
}

// CloseSend calls CloseWrite, if present, on the underlying
// io.ReadWriteCloser.  If CloseRecv was called before this function,
// then the underlying io.ReadWriteCloser is closed.
func (s *StreamTransport) CloseSend() error {
	s.mu.Lock()
	if s.closes&1 == 1 {
		s.mu.Unlock()
		return errors.New(errors.Disconnected, "rpc stream transport", "send already closed")
	}
	s.closes |= 1
	done := s.closes == 3
	s.mu.Unlock()

	werr := s.cw.CloseWrite()
	if !done {
		return werr
	}
	cerr := s.c.Close()
	if cerr != nil {
		return cerr
	}
	if werr != nil {
		return werr
	}
	return nil
}

// RecvMessage reads the next message from the underlying reader.
// The cancelation and deadline from ctx is ignored, but RecvMessage
// will return early if CloseRecv is called.
func (s *StreamTransport) RecvMessage(ctx context.Context) (rpccp.Message, capnp.ReleaseFunc, error) {
	return s.recv.RecvMessage(ctx)
}

// CloseRecv calls CloseRead, if present, on the underlying
// io.ReadWriteCloser.  If CloseSend was called before this function,
// then the underlying io.ReadWriteCloser is closed.
func (s *StreamTransport) CloseRecv() error {
	s.mu.Lock()
	if s.closes&2 == 2 {
		s.mu.Unlock()
		return errors.New(errors.Disconnected, "rpc stream transport", "receive already closed")
	}
	s.closes |= 2
	done := s.closes == 3
	s.mu.Unlock()

	rerr := s.recv.CloseRecv()
	if !done {
		return rerr
	}
	cerr := s.c.Close()
	if cerr != nil {
		return cerr
	}
	if rerr != nil {
		return rerr
	}
	return nil
}

// closerReceiver receives messages from a decoder, relying on a
// readCloser to interrupt the underlying io.Reader.
type closerReceiver struct {
	dec    *capnp.Decoder
	closer readCloser
}

func (cr closerReceiver) RecvMessage(ctx context.Context) (rpccp.Message, capnp.ReleaseFunc, error) {
	msg, err := cr.dec.Decode()
	if err != nil {
		return rpccp.Message{}, nil, err
	}
	rmsg, err := rpccp.ReadRootMessage(msg)
	if err != nil {
		return rpccp.Message{}, nil, err
	}
	return rmsg, func() {
		msg.Reset(nil)
	}, nil
}

func (cr closerReceiver) CloseRecv() error {
	return cr.closer.CloseRead()
}

// signalReceiver receives messages from a decoder, abandoning a Decode
// once CloseRecv is called.  It is assumed that the caller will then
// eventually interrupt the read, usually by calling Close on the
// underlying io.ReadCloser.
type signalReceiver struct {
	dec   *capnp.Decoder
	close chan struct{}
}

func (sr signalReceiver) RecvMessage(ctx context.Context) (rpccp.Message, capnp.ReleaseFunc, error) {
	select {
	case <-sr.close:
		return rpccp.Message{}, nil, errors.New(errors.Disconnected, "rpc stream transport", "receive on closed receiver")
	default:
	}
	var msg *capnp.Message
	var err error
	read := make(chan struct{})
	go func() {
		msg, err = sr.dec.Decode()
		close(read)
	}()
	select {
	case <-read:
	case <-sr.close:
		return rpccp.Message{}, nil, errors.New(errors.Disconnected, "rpc stream transport", "receive on closed receiver")
	}
	if err != nil {
		return rpccp.Message{}, nil, err
	}
	rmsg, err := rpccp.ReadRootMessage(msg)
	if err != nil {
		return rpccp.Message{}, nil, err
	}
	return rmsg, func() {
		msg.Reset(nil)
	}, nil
}

func (sr signalReceiver) CloseRecv() error {
	close(sr.close)
	return nil
}

// Optional interfaces that io.ReadWriteClosers could implement.
// See net.TCPConn for docs.

type writeDeadlineSetter interface {
	SetWriteDeadline(t time.Time) error
}

type readCloser interface {
	CloseRead() error
}

type writeCloser interface {
	CloseWrite() error
}
