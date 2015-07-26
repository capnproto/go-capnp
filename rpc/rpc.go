// Package rpc implements the Cap'n Proto RPC protocol.
package rpc

import (
	"fmt"
	"log"
	"sync"

	"golang.org/x/net/context"
	"zombiezen.com/go/capnproto"
	"zombiezen.com/go/capnproto/rpc/rpccapnp"
)

// Note on concurrency:
// Each connection has two primary goroutines: the reader and the task
// queue.  User code -- like client calls -- should be executed in a
// separate goroutine so that the connection can still be used.
// Table entries (i.e. imports, exports, questions, and answers) should
// only be modified in the task goroutine, but can be read from any
// goroutine.  The reader goroutine should not block on a send.

// A Conn is a connection to another Cap'n Proto vat.
// It is safe to use from multiple goroutines.
type Conn struct {
	transport Transport
	main      capnp.Client

	questions questionTable
	answers   answerTable
	imports   importTable
	exports   exportTable

	manager manager
	tasks   chan task
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
func NewConn(t Transport, options ...ConnOption) *Conn {
	conn := &Conn{
		transport: t,
		tasks:     make(chan task),
	}
	conn.manager.init()
	for _, o := range options {
		o(conn)
	}
	conn.manager.do(conn.dispatchRecv)
	conn.manager.do(conn.dispatchTasks)
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
	// Stop helper goroutines.
	if !c.manager.shutdown(ErrConnClosed) {
		return ErrConnClosed
	}
	// Hang up.
	// TODO(light): add timeout to write.
	ctx := context.Background()
	s := capnp.NewBuffer(nil)
	n := rpccapnp.NewRootMessage(s)
	e := rpccapnp.NewException(s)
	toException(e, errShutdown)
	n.SetAbort(e)
	werr := c.transport.SendMessage(ctx, n)
	cerr := c.transport.Close()
	if werr != nil {
		return werr
	}
	if cerr != nil {
		return cerr
	}
	return nil
}

// A task is a function scheduled to run on a connection.
type task struct {
	f    func() error
	done chan error
}

// do runs a function in the connection's task queue.
func (c *Conn) do(ctx context.Context, f func() error) error {
	t := task{f, make(chan error, 1)}
	select {
	case c.tasks <- t:
	case <-ctx.Done():
		return ctx.Err()
	case <-c.manager.finish:
		return c.manager.err()
	}
	select {
	case err := <-t.done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-c.manager.finish:
		return c.manager.err()
	}
}

// dispatchTasks runs in its own goroutine.
func (c *Conn) dispatchTasks() {
	for {
		var t task
		select {
		case t = <-c.tasks:
		case <-c.manager.finish:
			return
		}
		err := t.f()
		t.done <- err
	}
}

// Bootstrap returns the receiver's main interface.
func (c *Conn) Bootstrap(ctx context.Context) capnp.Client {
	// TODO(light): don't block
	var q *question
	err := c.do(ctx, func() error {
		q = c.questions.new(c, ctx, nil)
		s := capnp.NewBuffer(nil)
		m := rpccapnp.NewRootMessage(s)
		b := rpccapnp.NewBootstrap(s)
		b.SetQuestionId(uint32(q.id))
		m.SetBootstrap(b)
		return c.transport.SendMessage(ctx, m)
	})
	var ans capnp.Answer
	if err != nil {
		ans = capnp.ErrorAnswer(err)
	} else {
		ans = q
	}
	return capnp.NewPipeline(ans).Client()
}

// handleMessage is run in the tasks goroutine.  It will block the
// reader goroutine until it sends to readContinue.  If m is
// needed beyond the function's lifetime, then it should be copied.
func (c *Conn) handleMessage(m rpccapnp.Message, readContinue chan<- struct{}) {
	switch m.Which() {
	case rpccapnp.Message_Which_unimplemented:
		// no-op for now to avoid feedback loop
	case rpccapnp.Message_Which_abort:
		a := Abort{copyRPCMessage(m).Abort()}
		readContinue <- struct{}{}
		log.Print(a)
		c.manager.shutdown(a)
	case rpccapnp.Message_Which_return:
		mm := copyRPCMessage(m)
		readContinue <- struct{}{}
		if err := c.handleReturn(mm.Return()); err != nil {
			log.Println("rpc: handle return:", err)
		}
	case rpccapnp.Message_Which_finish:
		// TODO(light): what if answers never had this ID?
		// TODO(light): return if cancelled
		// TODO(light): release
		id := answerID(m.Finish().QuestionId())
		readContinue <- struct{}{}
		c.answers.pop(id)
	case rpccapnp.Message_Which_bootstrap:
		id := answerID(m.Bootstrap().QuestionId())
		readContinue <- struct{}{}
		if err := c.handleBootstrap(id); err != nil {
			log.Println("rpc: handle bootstrap:", err)
		}
	case rpccapnp.Message_Which_call:
		mm := copyRPCMessage(m)
		readContinue <- struct{}{}
		if err := c.handleCall(mm); err != nil {
			log.Println("rpc: handle call:", err)
		}
	case rpccapnp.Message_Which_release:
		id := exportID(m.Release().Id())
		refs := int(m.Release().ReferenceCount())
		readContinue <- struct{}{}
		c.exports.release(id, refs)
	default:
		log.Printf("rpc: received unimplemented message, which = %v", m.Which())
		um := newUnimplementedMessage(nil, m)
		readContinue <- struct{}{}
		if err := c.transport.SendMessage(c.manager.context(), um); err != nil {
			log.Println("rpc: writing unimplemented:", err)
		}
	}
}

func copyMessage(msg capnp.Message) capnp.Message {
	n := 0
	for {
		if _, err := msg.Lookup(uint32(n)); err != nil {
			break
		}
		n++
	}
	segments := make([][]byte, n)
	for i := range segments {
		s, err := msg.Lookup(uint32(i))
		if err != nil {
			panic(err)
		}
		segments[i] = make([]byte, len(s.Data))
		copy(segments[i], s.Data)
	}
	return capnp.NewMultiBuffer(segments).Message
}

func copyRPCMessage(m rpccapnp.Message) rpccapnp.Message {
	mm := copyMessage(m.Segment.Message)
	seg, err := mm.Lookup(0)
	if err != nil {
		panic(err)
	}
	return rpccapnp.ReadRootMessage(seg)
}

// newContext creates a new context for a received message sequence.
func (c *Conn) newContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(c.manager.context())
}

func newUnimplementedMessage(buf []byte, m rpccapnp.Message) rpccapnp.Message {
	s := capnp.NewBuffer(buf)
	n := rpccapnp.NewRootMessage(s)
	n.SetUnimplemented(m)
	return n
}

func (c *Conn) sendCall(cl *call) capnp.Answer {
	var q *question
	err := c.do(cl.Ctx, func() error {
		q = c.questions.new(c, cl.Ctx, &cl.Method)
		msg := c.newCallMessage(nil, q.id, cl)
		return c.transport.SendMessage(cl.Ctx, msg)
	})
	if err != nil {
		return capnp.ErrorAnswer(err)
	}
	return q
}

func (c *Conn) newCallMessage(buf []byte, id questionID, cl *call) rpccapnp.Message {
	s := capnp.NewBuffer(buf)
	msg := rpccapnp.NewRootMessage(s)

	msgCall := rpccapnp.NewCall(s)
	msgCall.SetQuestionId(uint32(id))
	msgCall.SetInterfaceId(cl.Method.InterfaceID)
	msgCall.SetMethodId(cl.Method.MethodID)

	target := rpccapnp.NewMessageTarget(s)
	if cl.transform != nil {
		a := rpccapnp.NewPromisedAnswer(s)
		a.SetQuestionId(uint32(cl.questionID))
		transformToPromisedAnswer(s, a, cl.transform)
		target.SetPromisedAnswer(a)
	} else {
		target.SetImportedCap(uint32(cl.importID))
	}
	msgCall.SetTarget(target)

	payload := rpccapnp.NewPayload(s)
	params := cl.PlaceParams(s)
	payload.SetContent(capnp.Object(params))
	payload.SetCapTable(c.makeCapTable(s))
	msgCall.SetParams(payload)

	msg.SetCall(msgCall)
	return msg
}

// release sends a release message over the connection.
func (c *Conn) release(id importID) error {
	ctx := c.manager.context()
	return c.do(ctx, func() error {
		i := c.imports.pop(id)
		if i == 0 {
			return nil
		}
		// TODO(light): deadline to close?
		s := capnp.NewBuffer(nil)
		msg := rpccapnp.NewRootMessage(s)
		mr := rpccapnp.NewRelease(s)
		mr.SetId(uint32(id))
		mr.SetReferenceCount(uint32(i))
		msg.SetRelease(mr)
		return c.transport.SendMessage(ctx, msg)
	})
}

// handleReturn is run in the tasks goroutine.
func (c *Conn) handleReturn(m rpccapnp.Return) error {
	id := questionID(m.AnswerId())
	q := c.questions.get(id)
	if q == nil {
		return fmt.Errorf("received return for unknown question id=%d", id)
	}
	releaseResultCaps := true
	switch m.Which() {
	case rpccapnp.Return_Which_results:
		releaseResultCaps = false
		c.populateMessageCapTable(m.Results())
		q.resolve(capnp.ImmediateAnswer(m.Results().Content()))
	case rpccapnp.Return_Which_exception:
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
	case rpccapnp.Return_Which_canceled:
		q.resolve(capnp.ErrorAnswer(errCallCanceled))
		// Don't send another finish message.
		return nil
	default:
		um := newUnimplementedMessage(nil, rpccapnp.ReadRootMessage(m.Segment))
		return c.transport.SendMessage(c.manager.context(), um)
	}
	fin := newFinishMessage(nil, id, releaseResultCaps)
	if err := c.transport.SendMessage(c.manager.context(), fin); err != nil {
		log.Printf("rpc: failed to write finish for %v ID=%d: %v", q.method, id, err)
	}
	c.questions.remove(id)
	return nil
}

func newFinishMessage(buf []byte, questionID questionID, release bool) rpccapnp.Message {
	s := capnp.NewBuffer(buf)
	m := rpccapnp.NewRootMessage(s)
	f := rpccapnp.NewFinish(s)
	f.SetQuestionId(uint32(questionID))
	f.SetReleaseResultCaps(release)
	m.SetFinish(f)
	return m
}

// populateMessageCapTable converts the descriptors in the payload into
// clients and sets it on the message the payload is a part of.
func (c *Conn) populateMessageCapTable(payload rpccapnp.Payload) {
	msg := payload.Segment.Message
	for i, n := 0, payload.CapTable().Len(); i < n; i++ {
		desc := payload.CapTable().At(i)
		switch desc.Which() {
		case rpccapnp.CapDescriptor_Which_none:
			msg.AddCap(nil)
		case rpccapnp.CapDescriptor_Which_senderHosted:
			id := importID(desc.SenderHosted())
			client := c.imports.addRef(c, id)
			msg.AddCap(client)
		case rpccapnp.CapDescriptor_Which_receiverHosted:
			id := exportID(desc.ReceiverHosted())
			e := c.exports.get(id)
			if e == nil {
				msg.AddCap(nil)
			} else {
				msg.AddCap(e.client)
			}
		case rpccapnp.CapDescriptor_Which_receiverAnswer:
			msg.AddCap(c.getPromisedAnswer(desc.ReceiverAnswer()))
		default:
			log.Println("rpc: unknown capability type", desc.Which())
			msg.AddCap(nil)
		}
	}
}

// makeCapTable converts the clients in the segment's message into capability descriptors.
func (c *Conn) makeCapTable(s *capnp.Segment) rpccapnp.CapDescriptor_List {
	msgtab := s.Message.CapTable()
	t := rpccapnp.NewCapDescriptor_List(s, len(msgtab))
	for i, client := range msgtab {
		desc := t.At(i)
		if client == nil {
			desc.SetNone()
			continue
		}
		c.descriptorForClient(desc, client)
	}
	return t
}

func (c *Conn) descriptorForClient(desc rpccapnp.CapDescriptor, client capnp.Client) {
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
			ans, id := q.promiseInfo()
			if ans == nil && q.conn == c {
				a := rpccapnp.NewPromisedAnswer(desc.Segment)
				a.SetQuestionId(uint32(id))
				transformToPromisedAnswer(desc.Segment, a, p.Transform())
				desc.SetReceiverAnswer(a)
				return
			} else if ans != nil {
				s, err := ans.Struct()
				if err == nil {
					client := capnp.TransformObject(capnp.Object(s), p.Transform()).ToInterface().Client()
					if client != nil {
						c.descriptorForClient(desc, client)
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

// handleBootstrap is run in the tasks goroutine.
func (c *Conn) handleBootstrap(id answerID) error {
	ctx, cancel := c.newContext()
	a := c.answers.insert(id, cancel)
	retmsg := newReturnMessage(id)
	ret := retmsg.Return()
	if a == nil {
		// Question ID reused, error out.
		setReturnException(ret, errQuestionReused)
		return c.transport.SendMessage(ctx, retmsg)
	}
	if c.main == nil {
		e := setReturnException(ret, errNoMainInterface)
		a.resolve(capnp.ErrorAnswer(Exception{e}))
		return c.transport.SendMessage(ctx, retmsg)
	}
	exportID := c.exports.add(c.main)
	retmsg.Segment.Message.AddCap(c.main)
	payload := rpccapnp.NewPayload(retmsg.Segment)
	const capIndex = 0
	in := capnp.Object(retmsg.Segment.NewInterface(capIndex))
	payload.SetContent(in)
	ctab := rpccapnp.NewCapDescriptor_List(retmsg.Segment, capIndex+1)
	ctab.At(capIndex).SetSenderHosted(uint32(exportID))
	payload.SetCapTable(ctab)
	ret.SetResults(payload)
	a.resolve(capnp.ImmediateAnswer(capnp.Object(in)))
	return c.transport.SendMessage(ctx, retmsg)
}

func (c *Conn) resolveTarget(mt rpccapnp.MessageTarget) capnp.Client {
	switch mt.Which() {
	case rpccapnp.MessageTarget_Which_importedCap:
		id := exportID(mt.ImportedCap())
		e := c.exports.get(id)
		if e == nil {
			return nil
		}
		return e.client
	case rpccapnp.MessageTarget_Which_promisedAnswer:
		return c.getPromisedAnswer(mt.PromisedAnswer())
	default:
		return nil
	}
}

func (c *Conn) getPromisedAnswer(pa rpccapnp.PromisedAnswer) capnp.Client {
	id := answerID(pa.QuestionId())
	a := c.answers.get(id)
	if a == nil {
		return nil
	}
	p := capnp.NewPipeline(a)
	for i := 0; i < pa.Transform().Len(); i++ {
		op := pa.Transform().At(i)
		switch op.Which() {
		case rpccapnp.PromisedAnswer_Op_Which_getPointerField:
			p = p.GetPipeline(int(op.GetPointerField()))
		case rpccapnp.PromisedAnswer_Op_Which_noop:
			fallthrough
		default:
			// do nothing
		}
	}
	return p.Client()
}

// handleCall is run in the tasks goroutine.  It mutates the capability
// table of its parameter.
func (c *Conn) handleCall(m rpccapnp.Message) error {
	ctx, cancel := c.newContext()
	id := answerID(m.Call().QuestionId())
	retmsg := newReturnMessage(id)
	ret := retmsg.Return()
	target := c.resolveTarget(m.Call().Target())
	if target == nil {
		setReturnException(ret, errBadTarget)
		return c.transport.SendMessage(ctx, retmsg)
	}
	a := c.answers.insert(id, cancel)
	if a == nil {
		// Question ID reused, error out.
		setReturnException(ret, errQuestionReused)
		return c.transport.SendMessage(ctx, retmsg)
	}
	c.populateMessageCapTable(m.Call().Params())

	go func() {
		meth := capnp.Method{
			InterfaceID: m.Call().InterfaceId(),
			MethodID:    m.Call().MethodId(),
		}
		answer := target.Call(&capnp.Call{
			Ctx:    ctx,
			Method: meth,
			Params: m.Call().Params().Content().ToStruct(),
		})
		// TODO(light): check to see if it's one of our answer types
		results, rerr := answer.Struct()

		err := c.do(ctx, func() error {
			if rerr != nil {
				e := setReturnException(ret, rerr)
				a.resolve(capnp.ErrorAnswer(Exception{e}))
				return c.transport.SendMessage(ctx, retmsg)
			}
			payload := rpccapnp.NewPayload(retmsg.Segment)
			payload.SetContent(capnp.Object(results))
			payload.SetCapTable(c.makeCapTable(retmsg.Segment))
			ret.SetResults(payload)
			a.resolve(capnp.ImmediateAnswer(capnp.Object(results)))
			return c.transport.SendMessage(ctx, retmsg)
		})
		if err != nil {
			log.Printf("rpc: writing return from %v: %v", &meth, err)
		}
	}()
	return nil
}

func newReturnMessage(id answerID) rpccapnp.Message {
	s := capnp.NewBuffer(nil)
	retmsg := rpccapnp.NewRootMessage(s)
	ret := rpccapnp.NewReturn(s)
	ret.SetAnswerId(uint32(id))
	retmsg.SetReturn(ret)
	return retmsg
}

func setReturnException(ret rpccapnp.Return, err error) rpccapnp.Exception {
	e := rpccapnp.NewException(ret.Segment)
	toException(e, err)
	ret.SetException(e)
	return e
}

// dispatchRecv runs in its own goroutine.
func (c *Conn) dispatchRecv() {
	for {
		msg, err := c.transport.RecvMessage(c.manager.context())
		if err == nil {
			cont := make(chan struct{})
			f := func() error {
				c.handleMessage(msg, cont)
				close(cont)
				return nil
			}
			// Partial reimplementation of do(). We don't want to block on
			// function completion; we want to block on the reader signal.
			t := task{f, make(chan error, 1)}
			select {
			case c.tasks <- t:
			case <-c.manager.finish:
				return
			}
			select {
			case <-cont:
			case <-t.done:
				panic("handleMessage task completed before signaling reader continue")
			case <-c.manager.finish:
				return
			}
		} else if isTemporaryError(err) {
			log.Println("rpc: read temporary error:", err)
		} else {
			c.manager.shutdown(err)
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

type msgSend struct {
	ctx  context.Context
	msg  rpccapnp.Message
	done chan error
}

func makeMsgSend(ctx context.Context, msg rpccapnp.Message) msgSend {
	return msgSend{ctx, msg, make(chan error, 1)}
}

type call struct {
	*capnp.Call

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
	return ic.c.sendCall(&call{
		Call:     cl,
		importID: ic.id,
	})
}

func (ic *importClient) Close() error {
	return ic.c.release(ic.id)
}

// manager signals the running goroutines in a Conn.
type manager struct {
	finish chan struct{}
	wg     sync.WaitGroup
	ctx    context.Context

	mu   sync.RWMutex
	done bool
	e    error
}

func (m *manager) init() {
	m.finish = make(chan struct{})
	var cancel context.CancelFunc
	m.ctx, cancel = context.WithCancel(context.Background())
	go func() {
		<-m.finish
		cancel()
	}()
}

// context returns a context that is cancelled when the manager shuts down.
func (m *manager) context() context.Context {
	return m.ctx
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
