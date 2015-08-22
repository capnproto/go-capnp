package rpc

import (
	"bytes"
	"io"
	"log"
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

	enc  *capnp.Encoder
	dec  *capnp.Decoder
	wbuf bytes.Buffer
}

// StreamTransport creates a transport that sends and receives messages
// by serializing and deserializing unpacked Cap'n Proto messages.
// Closing the transport will close the underlying ReadWriteCloser.
func StreamTransport(rwc io.ReadWriteCloser) Transport {
	d, _ := rwc.(writeDeadlineSetter)
	s := &streamTransport{
		rwc:      rwc,
		deadline: d,
		dec:      capnp.NewDecoder(rwc),
	}
	s.wbuf.Grow(4096)
	s.enc = capnp.NewEncoder(&s.wbuf)
	return s
}

func (s *streamTransport) SendMessage(ctx context.Context, msg rpccapnp.Message) error {
	s.wbuf.Reset()
	if err := s.enc.Encode(msg.Segment().Message()); err != nil {
		return err
	}
	if s.deadline != nil {
		// TODO(light): log errors
		if d, ok := ctx.Deadline(); ok {
			s.deadline.SetWriteDeadline(d)
		} else {
			s.deadline.SetWriteDeadline(time.Time{})
		}
	}
	_, err := s.rwc.Write(s.wbuf.Bytes())
	return err
}

func (s *streamTransport) RecvMessage(ctx context.Context) (rpccapnp.Message, error) {
	var (
		msg *capnp.Message
		err error
	)
	read := make(chan struct{})
	go func() {
		msg, err = s.dec.Decode()
		close(read)
	}()
	select {
	case <-read:
	case <-ctx.Done():
		return rpccapnp.Message{}, ctx.Err()
	}
	if err != nil {
		return rpccapnp.Message{}, err
	}
	return rpccapnp.ReadRootMessage(msg)
}

func (s *streamTransport) Close() error {
	return s.rwc.Close()
}

type writeDeadlineSetter interface {
	SetWriteDeadline(t time.Time) error
}

// dispatchSend runs in its own goroutine and sends messages on a transport.
func dispatchSend(m *manager, transport Transport, msgs <-chan rpccapnp.Message) {
	for {
		select {
		case msg := <-msgs:
			err := transport.SendMessage(m.context(), msg)
			if err != nil {
				log.Printf("rpc: writing %v: %v", msg.Which(), err)
			}
		case <-m.finish:
			return
		}
	}
}

// sendMessage sends a message to out to be sent.  It returns an error
// if the manager finished.
func sendMessage(m *manager, out chan<- rpccapnp.Message, msg rpccapnp.Message) error {
	select {
	case out <- msg:
		return nil
	case <-m.finish:
		return m.err()
	}
}

// dispatchRecv runs in its own goroutine and receives messages from a transport.
func dispatchRecv(m *manager, transport Transport, msgs chan<- rpccapnp.Message) {
	for {
		msg, err := transport.RecvMessage(m.context())
		if err != nil {
			if isTemporaryError(err) {
				log.Println("rpc: read temporary error:", err)
				continue
			}
			m.shutdown(err)
			return
		}
		select {
		case msgs <- copyRPCMessage(msg):
		case <-m.finish:
			return
		}
	}
}

// copyMessage clones a Cap'n Proto buffer.
func copyMessage(msg *capnp.Message) *capnp.Message {
	n := msg.NumSegments()
	segments := make([][]byte, n)
	for i := range segments {
		s, err := msg.Segment(capnp.SegmentID(i))
		if err != nil {
			panic(err)
		}
		segments[i] = make([]byte, len(s.Data()))
		copy(segments[i], s.Data())
	}
	return &capnp.Message{Arena: capnp.MultiSegment(segments)}
}

// copyRPCMessage clones an RPC packet.
func copyRPCMessage(m rpccapnp.Message) rpccapnp.Message {
	mm := copyMessage(m.Segment().Message())
	rpcMsg, err := rpccapnp.ReadRootMessage(mm)
	if err != nil {
		panic(err)
	}
	return rpcMsg
}

// isTemporaryError reports whether e has a Temporary() method that
// returns true.
func isTemporaryError(e error) bool {
	type temp interface {
		Temporary() bool
	}
	t, ok := e.(temp)
	return ok && t.Temporary()
}
