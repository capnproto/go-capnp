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
	// Close.
	Close() error
}

// StreamTransport serializes and deserializes unpacked Cap'n Proto
// messages on a byte stream.  StreamTransport adds no buffering beyond
// what its underlying stream has.
type StreamTransport struct {
	cr ctxReader
	wc io.WriteCloser

	mu       sync.RWMutex
	interval time.Duration
	closed   bool
}

// NewStreamTransport creates a new transport that reads and writes to rwc.
// Closing the transport will close rwc.
//
// If rwc has SetReadDeadline or SetWriteDeadline methods, they will be
// used to handle Context cancellation and deadlines.
func NewStreamTransport(rwc io.ReadWriteCloser) *StreamTransport {
	return &StreamTransport{
		cr:       ctxReader{r: rwc},
		wc:       rwc,
		interval: 500 * time.Millisecond,
	}
}

// NewMessage allocates a new message to be sent.
//
// It is safe to call NewMessage concurrently with RecvMessage.
func (s *StreamTransport) NewMessage(ctx context.Context) (_ rpccp.Message, send func() error, release capnp.ReleaseFunc, _ error) {
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
		select {
		case <-ctx.Done():
			return errors.New(errors.Failed, "rpc stream transport", "send: "+ctx.Err().Error())
		default:
		}
		b, err := msg.Marshal()
		if err != nil {
			return errors.New(errors.Failed, "rpc stream transport", "send: "+err.Error())
		}
		s.mu.RLock()
		interval := s.interval
		s.mu.RUnlock()
		_, err = writeCtx(ctx, s.wc, b, interval)
		if err != nil {
			return errors.New(errors.Failed, "rpc stream transport", "send: "+err.Error())
		}
		return nil
	}
	release = func() {
		msg.Reset(nil)
	}
	return rmsg, send, release, nil
}

// RecvMessage reads the next message from the underlying reader.
//
// It is safe to call RecvMessage concurrently with NewMessage.
func (s *StreamTransport) RecvMessage(ctx context.Context) (rpccp.Message, capnp.ReleaseFunc, error) {
	s.mu.RLock()
	s.cr.interval = s.interval
	s.mu.RUnlock()
	s.cr.ctx = ctx
	msg, err := capnp.NewDecoder(&s.cr).Decode()
	if err != nil {
		return rpccp.Message{}, nil, errors.New(errors.Failed, "rpc stream transport", "receive: "+err.Error())
	}
	rmsg, err := rpccp.ReadRootMessage(msg)
	if err != nil {
		return rpccp.Message{}, nil, errors.New(errors.Failed, "rpc stream transport", "receive: "+err.Error())
	}
	return rmsg, func() { msg.Reset(nil) }, nil
}

// SetInterruptInterval sets the frequency at which reads and writes are
// woken up to check for cancellation.
func (s *StreamTransport) SetInterruptInterval(d time.Duration) {
	s.mu.Lock()
	s.interval = d
	s.mu.Unlock()
}

// Close closes the underlying ReadWriteCloser.
func (s *StreamTransport) Close() error {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return errors.New(errors.Disconnected, "rpc stream transport", "already closed")
	}
	s.closed = true
	s.mu.Unlock()
	err := s.wc.Close()
	s.cr.wait()
	if err != nil {
		return errors.New(errors.Failed, "rpc stream transport", "close: "+err.Error())
	}
	return nil
}

// ctxReader adds timeouts and cancellation to a reader.
type ctxReader struct {
	r        io.Reader
	interval time.Duration
	ctx      context.Context // set to change Context

	// internal state
	result chan readResult
	pos, n int
	err    error
	buf    [1024]byte
}

type readResult struct {
	n   int
	err error
}

// Read reads into p.  It makes a best effort to respect the Done signal
// in cr.ctx.
func (cr *ctxReader) Read(p []byte) (int, error) {
	if cr.pos < cr.n {
		// Buffered from previous read.
		n := copy(p, cr.buf[cr.pos:cr.n])
		cr.pos += n
		if cr.pos == cr.n && cr.err != nil {
			err := cr.err
			cr.err = nil
			return n, err
		}
		return n, nil
	}
	if cr.result != nil {
		// Read in progress.
		select {
		case r := <-cr.result:
			cr.result = nil
			cr.n = r.n
			cr.pos = copy(p, cr.buf[:cr.n])
			if cr.pos == cr.n && r.err != nil {
				return cr.pos, r.err
			}
			cr.err = r.err
			return cr.pos, nil
		case <-cr.ctx.Done():
			return 0, cr.ctx.Err()
		}
	}
	select {
	case <-cr.ctx.Done():
		return 0, cr.ctx.Err()
	default:
	}
	rd, ok := cr.r.(interface {
		SetReadDeadline(time.Time) error
	})
	if !ok {
		return cr.leakyRead(p)
	}
	deadline, hasDeadline := cr.ctx.Deadline()
	if err := rd.SetReadDeadline(nextDeadline(cr.interval, deadline, hasDeadline)); err != nil {
		return cr.leakyRead(p)
	}
	for {
		select {
		case <-cr.ctx.Done():
			return 0, cr.ctx.Err()
		default:
		}
		n, err := cr.r.Read(p)
		if isTimeout(err) {
			err = nil
		}
		if n > 0 || err != nil {
			return n, err
		}
		rd.SetReadDeadline(nextDeadline(cr.interval, deadline, hasDeadline))
	}
}

// leakyRead reads from the underlying reader in a separate goroutine.
// If the Context is Done before the read completes, then the goroutine
// will stay alive until cr.wait() is called.
func (cr *ctxReader) leakyRead(p []byte) (int, error) {
	cr.result = make(chan readResult)
	max := len(p)
	if max > len(cr.buf) {
		max = len(cr.buf)
	}
	go func() {
		n, err := cr.r.Read(cr.buf[:max])
		cr.result <- readResult{n, err}
	}()
	select {
	case r := <-cr.result:
		cr.result = nil
		copy(p, cr.buf[:r.n])
		return r.n, r.err
	case <-cr.ctx.Done():
		return 0, cr.ctx.Err()
	}
}

// wait waits until any goroutine started by leakyRead finishes.
func (cr *ctxReader) wait() {
	if cr.result == nil {
		return
	}
	r := <-cr.result
	cr.result = nil
	cr.pos, cr.n = 0, r.n
	cr.err = r.err
}

// writeCtx writes bytes to a writer while making a best effort to
// respect the Done signal of the Context.  However, once any bytes have
// been written to w, writeCtx will ignore the Done signal to avoid
// partial writes.
func writeCtx(ctx context.Context, w io.Writer, b []byte, interval time.Duration) (int, error) {
	select {
	case <-ctx.Done():
		// Early cancel.
		return 0, ctx.Err()
	default:
	}
	// Check for timeout support.
	wd, ok := w.(interface {
		SetWriteDeadline(time.Time) error
	})
	if !ok {
		return w.Write(b)
	}
	deadline, hasDeadline := ctx.Deadline()
	if err := wd.SetWriteDeadline(nextDeadline(interval, deadline, hasDeadline)); err != nil {
		return w.Write(b)
	}
	// Poll for cancel while we haven't written anything.
	n := 0
	for n == 0 {
		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
		}
		var err error
		n, err = w.Write(b)
		if err == nil || (err != nil && !isTimeout(err)) {
			return n, err
		}
		wd.SetWriteDeadline(nextDeadline(interval, deadline, hasDeadline))
	}
	// Data has been written.  Block until finished, since partial writes
	// are guaranteed protocol violations.
	wd.SetWriteDeadline(time.Time{})
	nn, err := w.Write(b[n:])
	n += nn
	return n, err
}

func nextDeadline(interval time.Duration, deadline time.Time, hasDeadline bool) time.Time {
	d := time.Now().Add(interval)
	if hasDeadline && d.After(deadline) {
		return deadline
	}
	return d
}

func isTimeout(e error) bool {
	te, ok := e.(interface {
		Timeout() bool
	})
	return ok && te.Timeout()
}
