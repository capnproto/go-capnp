// Package rpc implements the Cap'n Proto RPC protocol.
package rpc

import (
	"bytes"
	"io"
	"log"
	"sync"
	"time"

	"golang.org/x/net/context"
	"zombiezen.com/go/capnproto"
	"zombiezen.com/go/capnproto/rpc/internal/rpc"
)

// A Conn is a connection to another Cap'n Proto vat.
type Conn struct {
	c    io.ReadWriteCloser
	main capnp.Client

	questions questionTable
	answers   answerTable
	imports   map[importID]struct{}
	exports   exportTable

	manager manager
	writes  chan write
	calls   chan *call
}

// A ConnOption is an option for opening a connection.
type ConnOption func(*Conn)

// MainInterface specifies that the connection should use client when
// receiving bootstrap messages.  By default, all bootstrap messages will
// fail.
func MainInterface(client capnp.Client) ConnOption {
	return func(c *Conn) {
		c.main = client
	}
}

// NewConn creates a new connection that communicates on c.
// Closing the connection will cause c to be closed.
func NewConn(c io.ReadWriteCloser, options ...ConnOption) *Conn {
	conn := &Conn{
		c:      c,
		writes: make(chan write),
		calls:  make(chan *call),
	}
	conn.manager.finish = make(chan struct{})
	for _, o := range options {
		o(conn)
	}
	conn.manager.do(conn.dispatchReads)
	conn.manager.do(conn.dispatchWrites)
	conn.manager.do(conn.dispatchCalls)
	return conn
}

// Wait waits until the connection is closed or aborted by the remote vat.
// Wait will always return an error, usually ErrConnClosed or of type Abort.
func (c *Conn) Wait() error {
	<-c.manager.finish
	return c.manager.err()
}

// Close closes the connection.
func (c *Conn) Close() error {
	// Hang up.
	s := capnp.NewBuffer(nil)
	n := rpc.NewRootMessage(s)
	e := rpc.NewException(s)
	toException(e, errShutdown)
	n.SetAbort(e)
	werr := c.write(context.Background(), n)

	// Stop helper goroutines.
	if !c.manager.shutdown(ErrConnClosed) {
		return werr
	}
	cerr := c.c.Close()
	if werr != nil {
		return werr
	}
	if cerr != nil {
		return cerr
	}
	return nil
}

// Bootstrap returns the receiver's main interface.
func (c *Conn) Bootstrap(ctx context.Context) capnp.Client {
	q := c.questions.new(c, ctx, nil)
	s := capnp.NewBuffer(nil)
	m := rpc.NewRootMessage(s)
	b := rpc.NewBootstrap(s)
	b.SetQuestionId(uint32(q.id))
	m.SetBootstrap(b)
	var ans capnp.Answer
	if err := c.write(ctx, m); err != nil {
		ans = capnp.ErrorAnswer(err)
	} else {
		ans = q
	}
	return capnp.NewPipeline(ans).Client()
}

// handleMessage is run from the reader goroutine.
func (c *Conn) handleMessage(m rpc.Message) {
	switch m.Which() {
	case rpc.Message_Which_unimplemented:
		// no-op for now to avoid feedback loop
	case rpc.Message_Which_abort:
		a := Abort{m.Abort()}
		log.Print(a)
		c.manager.shutdown(a)
	case rpc.Message_Which_return:
		id := questionID(m.Return().AnswerId())
		q := c.questions.get(id)
		go c.handleReturn(id, q, m.Return())
	case rpc.Message_Which_finish:
		// TODO(light): what if answers never had this ID?
		// TODO(light): return if cancelled
		// TODO(light): release
		c.answers.pop(answerID(m.Finish().QuestionId()))
	case rpc.Message_Which_bootstrap:
		ctx, cancel := c.newContext()
		id := answerID(m.Bootstrap().QuestionId())
		a := c.answers.insert(id, cancel)
		go func() {
			m := c.handleBootstrap(ctx, id, a)
			if err := c.write(ctx, m); err != nil {
				log.Println("rpc: bootstrap:", err)
			}
		}()
	case rpc.Message_Which_call:
		ctx, cancel := c.newContext()
		id := answerID(m.Call().QuestionId())
		a := c.answers.insert(id, cancel)
		target := c.resolveTarget(m.Call().Target())
		go func() {
			hold := make(chan struct{})
			m := c.handleCall(ctx, id, a, target, m.Call(), hold)
			err := c.write(ctx, m)
			close(hold)
			if err != nil {
				log.Println("rpc: call:", err)
			}
		}()
	default:
		log.Printf("rpc: received unimplemented message, which = %v", m.Which())
		ctx, _ := c.newContext()
		go func() {
			err := c.write(ctx, newUnimplementedMessage(nil, m))
			if err != nil {
				log.Println("rpc: writing unimplemented:", err)
			}
		}()
	}
}

// newContext creates a new context for a received message sequence.
func (c *Conn) newContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(context.Background())
}

func newUnimplementedMessage(buf []byte, m rpc.Message) rpc.Message {
	s := capnp.NewBuffer(buf)
	n := rpc.NewRootMessage(s)
	n.SetUnimplemented(m)
	return n
}

// sendCall is run from the calls goroutine.
func (c *Conn) sendCall(cl *call) capnp.Answer {
	q := c.questions.new(c, cl.Ctx, &cl.Method)
	hold := make(chan struct{})
	msg := c.newCallMessage(nil, q.id, cl, hold)
	err := c.write(cl.Ctx, msg)
	close(hold)
	if err != nil {
		return capnp.ErrorAnswer(err)
	}
	return q
}

func (c *Conn) newCallMessage(buf []byte, id questionID, cl *call, hold <-chan struct{}) rpc.Message {
	s := capnp.NewBuffer(buf)
	msg := rpc.NewRootMessage(s)

	msgCall := rpc.NewCall(s)
	msgCall.SetQuestionId(uint32(id))
	msgCall.SetInterfaceId(cl.Method.InterfaceID)
	msgCall.SetMethodId(cl.Method.MethodID)

	target := rpc.NewMessageTarget(s)
	if cl.transform != nil {
		a := rpc.NewPromisedAnswer(s)
		a.SetQuestionId(uint32(cl.questionID))
		transformToPromisedAnswer(s, a, cl.transform)
		target.SetPromisedAnswer(a)
	} else {
		target.SetImportedCap(uint32(cl.importID))
	}
	msgCall.SetTarget(target)

	payload := rpc.NewPayload(s)
	params := cl.PlaceParams(s)
	payload.SetContent(capnp.Object(params))
	payload.SetCapTable(c.makeCapTable(s, hold))
	msgCall.SetParams(payload)

	msg.SetCall(msgCall)
	return msg
}

// handleReturn is run in its own goroutine.
func (c *Conn) handleReturn(id questionID, q *question, m rpc.Return) {
	if q == nil {
		log.Printf("rpc: received return for unknown question id=%d", id)
		return
	}
	releaseResultCaps := true
	switch m.Which() {
	case rpc.Return_Which_results:
		releaseResultCaps = false
		c.populateMessageCapTable(m.Results())
		q.resolve(capnp.ImmediateAnswer(m.Results().Content()))
	case rpc.Return_Which_exception:
		e := error(Exception{m.Exception()})
		if q.method != nil {
			e = &capnp.MethodError{
				Method: q.method,
				Err:    e,
			}
		} else {
			e = bootstrapError{e}
		}
		q.resolve(capnp.ErrorAnswer(e))
	case rpc.Return_Which_canceled:
		q.resolve(capnp.ErrorAnswer(errCallCanceled))
		// Don't send another finish message.
		return
	default:
		s := capnp.NewBuffer(nil)
		mm := rpc.NewRootMessage(s)
		mm.SetUnimplemented(rpc.ReadRootMessage(m.Segment))
		if err := c.write(context.Background(), mm); err != nil {
			log.Println("rpc: failed to write unimplemented return:", err)
		}
		return
	}
	fin := newFinishMessage(nil, id, releaseResultCaps)
	if err := c.write(context.Background(), fin); err != nil {
		log.Printf("rpc: failed to write finish for ID=%d: %v", id, err)
	}
	c.questions.remove(id)
}

func newFinishMessage(buf []byte, questionID questionID, release bool) rpc.Message {
	s := capnp.NewBuffer(buf)
	m := rpc.NewRootMessage(s)
	f := rpc.NewFinish(s)
	f.SetQuestionId(uint32(questionID))
	f.SetReleaseResultCaps(release)
	m.SetFinish(f)
	return m
}

// populateMessageCapTable converts the descriptors in the payload into
// clients and sets it on the message the payload is a part of.
func (c *Conn) populateMessageCapTable(payload rpc.Payload) {
	msg := payload.Segment.Message
	for i, n := 0, payload.CapTable().Len(); i < n; i++ {
		desc := payload.CapTable().At(i)
		switch desc.Which() {
		case rpc.CapDescriptor_Which_none:
			msg.AddCap(nil)
		case rpc.CapDescriptor_Which_senderHosted:
			// TODO(light): add import to table
			msg.AddCap(&importClient{c: c, id: importID(desc.SenderHosted())})
		case rpc.CapDescriptor_Which_receiverHosted:
			id := exportID(desc.ReceiverHosted())
			e := c.exports.get(id)
			if e == nil {
				msg.AddCap(nil)
			} else {
				msg.AddCap(e.client)
			}
		case rpc.CapDescriptor_Which_receiverAnswer:
			msg.AddCap(c.getPromisedAnswer(desc.ReceiverAnswer()))
		default:
			log.Println("rpc: unknown capability type", desc.Which())
			msg.AddCap(nil)
		}
	}
}

// makeCapTable converts the clients in the segment's message into capability descriptors.
// The hold channel should be closed when the descriptors have been written,
// since this blocks sending a Finish.
func (c *Conn) makeCapTable(s *capnp.Segment, hold <-chan struct{}) rpc.CapDescriptor_List {
	msgtab := s.Message.CapTable()
	t := rpc.NewCapDescriptor_List(s, len(msgtab))
	for i, client := range msgtab {
		desc := t.At(i)
		if client == nil {
			desc.SetNone()
			continue
		}
		c.descriptorForClient(desc, client, hold)
	}
	return t
}

func (c *Conn) descriptorForClient(desc rpc.CapDescriptor, client capnp.Client, hold <-chan struct{}) {
	if client == nil {
		id := c.exports.add(capnp.ErrorClient(capnp.ErrNullClient))
		desc.SetSenderHosted(uint32(id))
		return
	}
	switch client := client.(type) {
	case *importClient:
		if client.c == c {
			desc.SetReceiverHosted(uint32(client.id))
			return
		}
	case *capnp.PipelineClient:
		p := (*capnp.Pipeline)(client)
		if q, ok := p.Answer().(*question); ok {
			hold := hold // shadow intentional.
			if q.conn != c {
				hold = nil
			}
			ans, id := q.promiseInfo(hold)
			if ans == nil && q.conn == c {
				a := rpc.NewPromisedAnswer(desc.Segment)
				a.SetQuestionId(uint32(id))
				transformToPromisedAnswer(desc.Segment, a, p.Transform())
				desc.SetReceiverAnswer(a)
				return
			} else if ans != nil {
				s, err := ans.Struct()
				if err == nil {
					client := capnp.TransformObject(capnp.Object(s), p.Transform()).ToInterface().Client()
					if client != nil {
						c.descriptorForClient(desc, client, hold)
						return
					}
					err = capnp.ErrNullClient
				}
				id := c.exports.add(capnp.ErrorClient(err))
				desc.SetSenderHosted(uint32(id))
				return
			}
		}
	}

	// Fallback: host and export ourselves.
	id := c.exports.add(client)
	desc.SetSenderHosted(uint32(id))
}

// handleBootstrap is run in its own goroutine.
func (c *Conn) handleBootstrap(ctx context.Context, id answerID, a *answer) rpc.Message {
	retmsg := newReturnMessage(id)
	ret := retmsg.Return()
	if a == nil {
		// Question ID reused, error out.
		setReturnException(ret, errQuestionReused)
		return retmsg
	}
	if c.main == nil {
		e := setReturnException(ret, errNoMainInterface)
		a.resolve(capnp.ErrorAnswer(Exception{e}))
		return retmsg
	}
	exportID := c.exports.add(c.main)
	retmsg.Segment.Message.AddCap(c.main)
	payload := rpc.NewPayload(retmsg.Segment)
	const capIndex = 0
	in := capnp.Object(retmsg.Segment.NewInterface(capIndex))
	payload.SetContent(in)
	ctab := rpc.NewCapDescriptor_List(retmsg.Segment, capIndex+1)
	ctab.At(capIndex).SetSenderHosted(uint32(exportID))
	payload.SetCapTable(ctab)
	ret.SetResults(payload)
	a.resolve(capnp.ImmediateAnswer(capnp.Object(in)))
	return retmsg
}

func (c *Conn) resolveTarget(mt rpc.MessageTarget) capnp.Client {
	switch mt.Which() {
	case rpc.MessageTarget_Which_importedCap:
		id := exportID(mt.ImportedCap())
		e := c.exports.get(id)
		if e == nil {
			return nil
		}
		return e.client
	case rpc.MessageTarget_Which_promisedAnswer:
		return c.getPromisedAnswer(mt.PromisedAnswer())
	default:
		return nil
	}
}

func (c *Conn) getPromisedAnswer(pa rpc.PromisedAnswer) capnp.Client {
	id := answerID(pa.QuestionId())
	a := c.answers.get(id)
	if a == nil {
		return nil
	}
	p := capnp.NewPipeline(a)
	for i := 0; i < pa.Transform().Len(); i++ {
		op := pa.Transform().At(i)
		switch op.Which() {
		case rpc.PromisedAnswer_Op_Which_getPointerField:
			p = p.GetPipeline(int(op.GetPointerField()))
		case rpc.PromisedAnswer_Op_Which_noop:
			fallthrough
		default:
			// do nothing
		}
	}
	return p.Client()
}

// handleCall is run in its own goroutine.
// hold is closed when the message has been written.
func (c *Conn) handleCall(ctx context.Context, id answerID, a *answer, target capnp.Client, call rpc.Call, hold <-chan struct{}) rpc.Message {
	retmsg := newReturnMessage(id)
	ret := retmsg.Return()
	if a == nil {
		// Question ID reused, error out.
		setReturnException(ret, errQuestionReused)
		return retmsg
	}
	if target == nil {
		setReturnException(ret, errBadTarget)
		return retmsg
	}
	params := call.Params()
	c.populateMessageCapTable(params)

	answer := target.Call(&capnp.Call{
		Ctx: ctx,
		Method: capnp.Method{
			InterfaceID: call.InterfaceId(),
			MethodID:    call.MethodId(),
		},
		Params: params.Content().ToStruct(),
	})
	// TODO(light): check to see if it's one of our answer types
	results, err := answer.Struct()

	if err != nil {
		e := setReturnException(ret, err)
		a.resolve(capnp.ErrorAnswer(Exception{e}))
		return retmsg
	}
	payload := rpc.NewPayload(retmsg.Segment)
	payload.SetContent(capnp.Object(results))
	payload.SetCapTable(c.makeCapTable(retmsg.Segment, hold))
	ret.SetResults(payload)
	a.resolve(capnp.ImmediateAnswer(capnp.Object(results)))
	return retmsg
}

func newReturnMessage(id answerID) rpc.Message {
	s := capnp.NewBuffer(nil)
	retmsg := rpc.NewRootMessage(s)
	ret := rpc.NewReturn(s)
	ret.SetAnswerId(uint32(id))
	retmsg.SetReturn(ret)
	return retmsg
}

func setReturnException(ret rpc.Return, err error) rpc.Exception {
	e := rpc.NewException(ret.Segment)
	toException(e, err)
	ret.SetException(e)
	return e
}

// dispatchReads runs in its own goroutine.
func (c *Conn) dispatchReads() {
	type read struct {
		m   rpc.Message
		err error
	}

	reads := make(chan read, 1)
	start := make(chan struct{})
	defer close(start)
	go func() {
		for range start {
			// TODO(light): don't allocate on every read
			s, err := capnp.ReadFromStream(c.c, nil)
			var m rpc.Message
			if err == nil {
				m = rpc.ReadRootMessage(s)
			}
			reads <- read{m, err}
		}
	}()

	start <- struct{}{}
	for {
		select {
		case r := <-reads:
			if r.err == nil {
				c.handleMessage(r.m)
			} else if isTemporaryError(r.err) {
				log.Println("rpc: read temporary error:", r.err)
			} else {
				c.manager.shutdown(r.err)
				return
			}
			start <- struct{}{}
		case <-c.manager.finish:
			return
		}
	}
}

func isTemporaryError(e error) bool {
	type temp interface {
		Temporary() bool
	}
	t, ok := e.(temp)
	return ok && t.Temporary()
}

// dispatchWrites runs in its own goroutine.
func (c *Conn) dispatchWrites() {
	var buf bytes.Buffer
	deadline, _ := c.c.(writeDeadlineSetter)
	for {
		select {
		case w := <-c.writes:
			buf.Reset()
			w.done <- c.doWrite(w, &buf, deadline)
		case <-c.manager.finish:
			return
		}
	}
}

// doWrite runs in the writer goroutine.
func (c *Conn) doWrite(w write, buf *bytes.Buffer, deadline writeDeadlineSetter) error {
	if _, err := w.seg.WriteTo(buf); err != nil {
		return err
	}
	if d, ok := w.ctx.Deadline(); ok && deadline != nil {
		// TODO(light): log error
		deadline.SetWriteDeadline(d)
	}
	_, err := c.c.Write(buf.Bytes())
	return err
}

func (c *Conn) write(ctx context.Context, m rpc.Message) error {
	w := makeWrite(ctx, m.Segment)
	select {
	case c.writes <- w:
	case <-ctx.Done():
		return ctx.Err()
	case <-c.manager.finish:
		return c.manager.err()
	}
	select {
	case err := <-w.done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-c.manager.finish:
		return c.manager.err()
	}
}

// dispatchCalls runs in its own goroutine.
func (c *Conn) dispatchCalls() {
	for {
		select {
		case cl := <-c.calls:
			cl.ready <- c.sendCall(cl)
		case <-c.manager.finish:
			return
		}
	}
}

type write struct {
	ctx  context.Context
	seg  *capnp.Segment
	done chan error
}

func makeWrite(ctx context.Context, s *capnp.Segment) write {
	return write{ctx, s, make(chan error, 1)}
}

type call struct {
	*capnp.Call
	ready chan capnp.Answer

	// If transform != nil, then this call is on a promised answer.
	// Otherwise, importID is used.
	questionID questionID
	transform  []capnp.PipelineOp
	importID   importID
}

type importClient struct {
	c  *Conn
	id importID
}

func (ic *importClient) Call(cl *capnp.Call) capnp.Answer {
	ready := make(chan capnp.Answer, 1)
	ic.c.calls <- &call{
		Call:     cl,
		importID: ic.id,
		ready:    ready,
	}
	return <-ready
}

func (ic *importClient) Close() error {
	// TODO(light): Send release message.
	return nil
}

type writeDeadlineSetter interface {
	SetWriteDeadline(t time.Time) error
}

// manager signals the running goroutines in a Conn.
type manager struct {
	finish chan struct{}
	wg     sync.WaitGroup

	mu   sync.RWMutex
	done bool
	e    error
}

// do starts a function in a new goroutine and will block shutdown
// until it has returned.  If the manager has already started shutdown,
// then it is a no-op.
func (m *manager) do(f func()) {
	m.mu.RLock()
	done := m.done
	if !done {
		m.wg.Add(1)
	}
	m.mu.RUnlock()
	if !done {
		go func() {
			defer m.wg.Done()
			f()
		}()
	}
}

// shutdown closes the finish channel and sets the error.
// The first call to shutdown returns true; subsequent calls are no-ops
// and return false.
func (m *manager) shutdown(e error) bool {
	m.mu.Lock()
	ok := !m.done
	if ok {
		close(m.finish)
		m.done = true
		m.e = e
	}
	m.mu.Unlock()
	if ok {
		m.wg.Wait()
	}
	return ok
}

// err returns the error passed to shutdown.
func (m *manager) err() error {
	m.mu.RLock()
	e := m.e
	m.mu.RUnlock()
	return e
}
