// Package rpc implements the Cap'n Proto RPC protocol.
package rpc

import (
	"log"
	"sync"

	"golang.org/x/net/context"
	"zombiezen.com/go/capnproto"
	"zombiezen.com/go/capnproto/rpc/rpccapnp"
)

// A Conn is a connection to another Cap'n Proto vat.
type Conn struct {
	transport Transport
	main      capnp.Client

	questions questionTable
	answers   answerTable
	imports   map[importID]struct{}
	exports   exportTable

	manager manager
	sends   chan msgSend
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
func NewConn(t Transport, options ...ConnOption) *Conn {
	conn := &Conn{
		transport: t,
		sends:     make(chan msgSend),
		calls:     make(chan *call),
	}
	conn.manager.init()
	for _, o := range options {
		o(conn)
	}
	conn.manager.do(conn.dispatchRecv)
	conn.manager.do(conn.dispatchSend)
	conn.manager.do(conn.dispatchCall)
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
	n := rpccapnp.NewRootMessage(s)
	e := rpccapnp.NewException(s)
	toException(e, errShutdown)
	n.SetAbort(e)
	werr := c.send(context.Background(), n)

	// Stop helper goroutines.
	if !c.manager.shutdown(ErrConnClosed) {
		return werr
	}
	cerr := c.transport.Close()
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
	m := rpccapnp.NewRootMessage(s)
	b := rpccapnp.NewBootstrap(s)
	b.SetQuestionId(uint32(q.id))
	m.SetBootstrap(b)
	var ans capnp.Answer
	if err := c.send(ctx, m); err != nil {
		ans = capnp.ErrorAnswer(err)
	} else {
		ans = q
	}
	return capnp.NewPipeline(ans).Client()
}

// handleMessage is run from the reader goroutine.
func (c *Conn) handleMessage(m rpccapnp.Message) {
	switch m.Which() {
	case rpccapnp.Message_Which_unimplemented:
		// no-op for now to avoid feedback loop
	case rpccapnp.Message_Which_abort:
		a := Abort{m.Abort()}
		log.Print(a)
		c.manager.shutdown(a)
	case rpccapnp.Message_Which_return:
		id := questionID(m.Return().AnswerId())
		q := c.questions.get(id)
		go c.handleReturn(id, q, copyRPCMessage(m).Return())
	case rpccapnp.Message_Which_finish:
		// TODO(light): what if answers never had this ID?
		// TODO(light): return if cancelled
		// TODO(light): release
		c.answers.pop(answerID(m.Finish().QuestionId()))
	case rpccapnp.Message_Which_bootstrap:
		ctx, cancel := c.newContext()
		id := answerID(m.Bootstrap().QuestionId())
		a := c.answers.insert(id, cancel)
		go func() {
			m := c.handleBootstrap(ctx, id, a)
			if err := c.send(ctx, m); err != nil {
				log.Println("rpc: bootstrap:", err)
			}
		}()
	case rpccapnp.Message_Which_call:
		ctx, cancel := c.newContext()
		id := answerID(m.Call().QuestionId())
		a := c.answers.insert(id, cancel)
		target := c.resolveTarget(m.Call().Target())
		go func(mc rpccapnp.Call) {
			hold := make(chan struct{})
			m := c.handleCall(ctx, id, a, target, mc, hold)
			err := c.send(ctx, m)
			close(hold)
			if err != nil {
				log.Println("rpc: call:", err)
			}
		}(copyRPCMessage(m).Call())
	default:
		log.Printf("rpc: received unimplemented message, which = %v", m.Which())
		ctx, _ := c.newContext()
		go func(m rpccapnp.Message) {
			err := c.send(ctx, newUnimplementedMessage(nil, m))
			if err != nil {
				log.Println("rpc: writing unimplemented:", err)
			}
		}(copyRPCMessage(m))
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

// sendCall is run from the calls goroutine.
func (c *Conn) sendCall(cl *call) capnp.Answer {
	q := c.questions.new(c, cl.Ctx, &cl.Method)
	hold := make(chan struct{})
	msg := c.newCallMessage(nil, q.id, cl, hold)
	err := c.send(cl.Ctx, msg)
	close(hold)
	if err != nil {
		return capnp.ErrorAnswer(err)
	}
	return q
}

func (c *Conn) newCallMessage(buf []byte, id questionID, cl *call, hold <-chan struct{}) rpccapnp.Message {
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
	payload.SetCapTable(c.makeCapTable(s, hold))
	msgCall.SetParams(payload)

	msg.SetCall(msgCall)
	return msg
}

// handleReturn is run in its own goroutine.
func (c *Conn) handleReturn(id questionID, q *question, m rpccapnp.Return) {
	if q == nil {
		log.Printf("rpc: received return for unknown question id=%d", id)
		return
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
		return
	default:
		s := capnp.NewBuffer(nil)
		mm := rpccapnp.NewRootMessage(s)
		mm.SetUnimplemented(rpccapnp.ReadRootMessage(m.Segment))
		if err := c.send(c.manager.context(), mm); err != nil {
			log.Println("rpc: failed to write unimplemented return:", err)
		}
		return
	}
	fin := newFinishMessage(nil, id, releaseResultCaps)
	if err := c.send(c.manager.context(), fin); err != nil {
		log.Printf("rpc: failed to write finish for ID=%d: %v", id, err)
	}
	c.questions.remove(id)
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
			// TODO(light): add import to table
			msg.AddCap(&importClient{c: c, id: importID(desc.SenderHosted())})
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
// The hold channel should be closed when the descriptors have been written,
// since this blocks sending a Finish.
func (c *Conn) makeCapTable(s *capnp.Segment, hold <-chan struct{}) rpccapnp.CapDescriptor_List {
	msgtab := s.Message.CapTable()
	t := rpccapnp.NewCapDescriptor_List(s, len(msgtab))
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

func (c *Conn) descriptorForClient(desc rpccapnp.CapDescriptor, client capnp.Client, hold <-chan struct{}) {
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
func (c *Conn) handleBootstrap(ctx context.Context, id answerID, a *answer) rpccapnp.Message {
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
	payload := rpccapnp.NewPayload(retmsg.Segment)
	const capIndex = 0
	in := capnp.Object(retmsg.Segment.NewInterface(capIndex))
	payload.SetContent(in)
	ctab := rpccapnp.NewCapDescriptor_List(retmsg.Segment, capIndex+1)
	ctab.At(capIndex).SetSenderHosted(uint32(exportID))
	payload.SetCapTable(ctab)
	ret.SetResults(payload)
	a.resolve(capnp.ImmediateAnswer(capnp.Object(in)))
	return retmsg
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

// handleCall is run in its own goroutine.
// hold is closed when the message has been written.
func (c *Conn) handleCall(ctx context.Context, id answerID, a *answer, target capnp.Client, call rpccapnp.Call, hold <-chan struct{}) rpccapnp.Message {
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
	payload := rpccapnp.NewPayload(retmsg.Segment)
	payload.SetContent(capnp.Object(results))
	payload.SetCapTable(c.makeCapTable(retmsg.Segment, hold))
	ret.SetResults(payload)
	a.resolve(capnp.ImmediateAnswer(capnp.Object(results)))
	return retmsg
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
			c.handleMessage(msg)
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

// dispatchSend runs in its own goroutine.
func (c *Conn) dispatchSend() {
	finish := c.manager.context().Done()
	for {
		select {
		case s := <-c.sends:
			s.done <- c.transport.SendMessage(s.ctx, s.msg)
		case <-finish:
			return
		}
	}
}

func (c *Conn) send(ctx context.Context, m rpccapnp.Message) error {
	s := makeMsgSend(ctx, m)
	select {
	case c.sends <- s:
	case <-ctx.Done():
		return ctx.Err()
	case <-c.manager.finish:
		return c.manager.err()
	}
	select {
	case err := <-s.done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	case <-c.manager.finish:
		return c.manager.err()
	}
}

// dispatchCall runs in its own goroutine.
func (c *Conn) dispatchCall() {
	for {
		select {
		case cl := <-c.calls:
			cl.ready <- c.sendCall(cl)
		case <-c.manager.finish:
			return
		}
	}
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
