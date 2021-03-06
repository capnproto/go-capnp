package rpc

import (
	"bytes"
	"context"
	"io"
	"time"

	capnp "zombiezen.com/go/capnproto2"
)

// A Codec is responsible for encoding and decoding messages from
// a single logical stream.
type Codec interface {
	Encode(context.Context, *capnp.Message) error
	Decode(context.Context) (*capnp.Message, error)
	Close() error
}

// MessageConn represents a message-oriented connection.
type MessageConn interface {
	// NextReader returns the next data message received from the peer.
	//
	// There can be at most one open reader on a connection. Implementations MUST
	// discard the previous message if the application has not already consumed it.
	//
	// Once this method returns a non-nil error, all subsequent calls to this method
	// MUST return the same error.
	NextReader() (io.Reader, error)

	// NextWriter returns a writer for the next message to send.
	// The writer's Close method flushes the complete message to the network.
	//
	// There can be at most one open writer on a connection.  NextWriter MUST
	// close the previous writer if the application has not already done so.
	NextWriter() (io.WriteCloser, error)

	// Close the underlying network connections.  Implementations will vary in
	// whether and how they announce disconnection to the remote side.
	Close() error
}

// StreamCodec can encode/decode messages from a stream of bytes
type StreamCodec struct {
	r   *ctxReader
	dec *capnp.Decoder

	wc  *ctxWriteCloser
	enc *capnp.Encoder
}

// NewStreamCodec creates a StreamCodec with the specified encoding.
func NewStreamCodec(rwc io.ReadWriteCloser, e Encoding) *StreamCodec {
	c := &StreamCodec{
		r: &ctxReader{Reader: rwc},
		wc: &ctxWriteCloser{
			WriteCloser:         rwc,
			partialWriteTimeout: 30 * time.Second,
		},
	}

	c.dec = e.NewDecoder(c.r)
	c.enc = e.NewEncoder(c.wc)

	return c
}

// Encode the message.
func (c *StreamCodec) Encode(ctx context.Context, m *capnp.Message) error {
	c.wc.setWriteContext(ctx)
	return c.enc.Encode(m)
}

// Decode the next message from the stream.
func (c *StreamCodec) Decode(ctx context.Context) (*capnp.Message, error) {
	c.r.setReadContext(ctx)
	return c.dec.Decode()
}

// SetPartialWriteTimeout sets the timeout for completing the
// transmission of a partially sent message after the send is cancelled
// or interrupted for any future sends.  If not set, a reasonable
// non-zero value is used.
//
// Setting a shorter timeout may free up resources faster in the case of
// an unresponsive remote peer, but may also make the transport respond
// too aggressively to bursts of latency.
func (c *StreamCodec) SetPartialWriteTimeout(d time.Duration) {
	c.wc.partialWriteTimeout = d
}

// Close the underlying byte-stream.
func (c StreamCodec) Close() error {
	defer c.r.wait()

	return c.wc.Close()
}

// Encoding specifies the wire transmission format.
type Encoding interface {
	NewEncoder(io.Writer) *capnp.Encoder
	NewDecoder(io.Reader) *capnp.Decoder
}

// BasicEncoding is a capnp.NewEncoder/Decoder pair.
type BasicEncoding struct{}

// NewEncoder using capnp.NewEncoder.
func (BasicEncoding) NewEncoder(w io.Writer) *capnp.Encoder { return capnp.NewEncoder(w) }

// NewDecoder using capnp.NewDecoder.
func (BasicEncoding) NewDecoder(r io.Reader) *capnp.Decoder { return capnp.NewDecoder(r) }

// PackedEncoding is a capnp.NewPackedEncoder/Decoder pair.
type PackedEncoding struct{}

// NewEncoder using capnp.NewPackedEncoder.
func (PackedEncoding) NewEncoder(w io.Writer) *capnp.Encoder { return capnp.NewPackedEncoder(w) }

// NewDecoder using capnp.NewPackedDecoder.
func (PackedEncoding) NewDecoder(r io.Reader) *capnp.Decoder { return capnp.NewPackedDecoder(r) }

// ctxReader adds timeouts and cancellation to a reader.
type ctxReader struct {
	io.Reader
	ctx context.Context // set to change Context

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

func (cr *ctxReader) setReadContext(ctx context.Context) { cr.ctx = ctx }

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
	// Check for early cancel.
	select {
	case <-cr.ctx.Done():
		return 0, cr.ctx.Err()
	default:
	}
	// Query timeout support.
	rd, ok := cr.Reader.(interface {
		SetReadDeadline(time.Time) error
	})
	if !ok {
		return cr.leakyRead(p)
	}
	if err := rd.SetReadDeadline(time.Now()); err != nil {
		return cr.leakyRead(p)
	}
	// Start separate goroutine to wait on Context.Done.
	if d, ok := cr.ctx.Deadline(); ok {
		rd.SetReadDeadline(d)
	} else {
		rd.SetReadDeadline(time.Time{})
	}
	readDone := make(chan struct{})
	listenDone := make(chan struct{})
	go func() {
		defer close(listenDone)
		select {
		case <-cr.ctx.Done():
			rd.SetReadDeadline(time.Now()) // interrupt read
		case <-readDone:
		}
	}()
	n, err := cr.Reader.Read(p)
	close(readDone)
	<-listenDone
	return n, err
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
		n, err := cr.Reader.Read(cr.buf[:max])
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

// wait until any goroutine started by leakyRead finishes.
func (cr *ctxReader) wait() {
	if cr.result == nil {
		return
	}
	r := <-cr.result
	cr.result = nil
	cr.pos, cr.n = 0, r.n
	cr.err = r.err
}

type ctxWriteCloser struct {
	io.WriteCloser
	ctx                 context.Context
	partialWriteTimeout time.Duration
}

// Write bytes to a writer while making a best effort to
// respect the Done signal of the Context.  However, if allowPartial is
// false, then once any bytes have been written to w, writeCtx will
// ignore the Done signal to avoid partial writes.
func (wc *ctxWriteCloser) Write(b []byte) (int, error) {
	n, err := wc.write(b)
	if n > 0 && n < len(b) {
		err = partialWriteError{err}
	}

	return n, err
}

func (wc *ctxWriteCloser) setWriteContext(ctx context.Context) { wc.ctx = ctx }

func (wc *ctxWriteCloser) write(b []byte) (int, error) {
	select {
	case <-wc.ctx.Done():
		// Early cancel.
		return 0, wc.ctx.Err()
	default:
	}
	// Check for timeout support.
	wd, ok := wc.WriteCloser.(interface {
		SetWriteDeadline(time.Time) error
	})
	if !ok {
		return wc.WriteCloser.Write(b)
	}
	if err := wd.SetWriteDeadline(time.Now()); err != nil {
		return wc.WriteCloser.Write(b)
	}
	// Start separate goroutine to wait on Context.Done.
	if d, ok := wc.ctx.Deadline(); ok {
		wd.SetWriteDeadline(d)
	} else {
		wd.SetWriteDeadline(time.Time{})
	}
	writeDone := make(chan struct{})
	listenDone := make(chan struct{})
	go func() {
		defer close(listenDone)
		select {
		case <-wc.ctx.Done():
			wd.SetWriteDeadline(time.Now()) // interrupt write
		case <-writeDone:
		}
	}()
	n, err := wc.WriteCloser.Write(b)
	close(writeDone)
	<-listenDone
	if wc.partialWriteTimeout <= 0 || n == 0 || !isTimeout(err) {
		return n, err
	}
	// Data has been written.  Block with extra partial timeout, since
	// partial writes are guaranteed protocol violations.
	wd.SetWriteDeadline(time.Now().Add(wc.partialWriteTimeout))
	nn, err := wc.WriteCloser.Write(b[n:])
	return n + nn, err
}

func isTimeout(e error) bool {
	te, ok := e.(interface {
		Timeout() bool
	})
	return ok && te.Timeout()
}

type messageCodec struct {
	c MessageConn
	e Encoding
}

func newMessageCodec(c MessageConn, e Encoding) *messageCodec {
	return &messageCodec{
		c: c,
		e: e,
	}
}

func (c messageCodec) Encode(ctx context.Context, m *capnp.Message) error {
	var buf bytes.Buffer
	if err := c.e.NewEncoder(&buf).Encode(m); err != nil {
		return err
	}

	w, err := c.c.NextWriter()
	if err != nil {
		return err
	}

	if _, err = io.Copy(w, &buf); err != nil {
		return err
	}

	// does the connection support write deadlines?
	if conn, ok := c.c.(interface{ SetWriteDeadline(time.Time) error }); ok {
		t, _ := ctx.Deadline()
		conn.SetWriteDeadline(t) // t defaults to time.Time{}, i.e. no deadline.
	}

	select {
	case <-ctx.Done():
		// The context may have been cancelled in the meantime.
		return ctx.Err()

	default:
		// Do not use defer to close the writer.
		// The io.WriteCloser produced by MessageCon waits until its Close method
		// is called before writing the message to the network.
		return w.Close()
	}
}

func (c messageCodec) Decode(context.Context) (*capnp.Message, error) {
	panic("NOT IMPLEMENTED")
}

func (c messageCodec) Close() error {
	return c.c.Close()
}
