package rpc

import (
	"bytes"
	"io"
	"time"

	"golang.org/x/net/context"
	"zombiezen.com/go/capnproto"
	"zombiezen.com/go/capnproto/rpc/rpccapnp"
)

// Transport is the interface that abstracts sending and receiving
// individual messages of the Cap'n Proto RPC protocol.
type Transport interface {
	// SendMessage sends msg.
	SendMessage(ctx context.Context, msg rpccapnp.Message) error

	// RecvMessage waits to receive a message and returns it.
	// Implementations may re-use buffers between calls, so the message is
	// only valid until the next call to RecvMessage.
	RecvMessage(ctx context.Context) (rpccapnp.Message, error)

	// Close releases any resources associated with the transport.
	Close() error
}

type streamTransport struct {
	rwc      io.ReadWriteCloser
	deadline writeDeadlineSetter
	rbuf     bytes.Buffer
	wbuf     bytes.Buffer
}

// StreamTransport creates a transport that sends and receives messages
// by serializing and deserializing unpacked Cap'n Proto messages.
// Closing the transport will close the underlying ReadWriteCloser.
func StreamTransport(rwc io.ReadWriteCloser) Transport {
	d, _ := rwc.(writeDeadlineSetter)
	s := &streamTransport{rwc: rwc, deadline: d}
	s.rbuf.Grow(4096)
	s.wbuf.Grow(4096)
	return s
}

func (s *streamTransport) SendMessage(ctx context.Context, msg rpccapnp.Message) error {
	s.wbuf.Reset()
	if _, err := msg.Segment.WriteTo(&s.wbuf); err != nil {
		return err
	}
	if d, ok := ctx.Deadline(); ok && s.deadline != nil {
		// TODO(light): log error
		s.deadline.SetWriteDeadline(d)
	}
	_, err := s.rwc.Write(s.wbuf.Bytes())
	return err
}

func (s *streamTransport) RecvMessage(ctx context.Context) (rpccapnp.Message, error) {
	seg, err := capnp.ReadFromStream(s.rwc, &s.rbuf)
	if err != nil {
		return rpccapnp.Message{}, err
	}
	msg := rpccapnp.ReadRootMessage(seg)
	return msg, nil
}

func (s *streamTransport) Close() error {
	return s.rwc.Close()
}

type writeDeadlineSetter interface {
	SetWriteDeadline(t time.Time) error
}
