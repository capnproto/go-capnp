// Package rpc implements the Cap'n Proto RPC protocol.
package rpc

import (
	"fmt"
	"log"

	"golang.org/x/net/context"
	"zombiezen.com/go/capnproto"
	"zombiezen.com/go/capnproto/rpc/rpccapnp"
)

// A Conn is a connection to another Cap'n Proto vat.
// It is safe to use from multiple goroutines.
type Conn struct {
	transport Transport
	main      capnp.Client

	manager    manager
	in         <-chan rpccapnp.Message
	out        chan<- outgoingMessage
	calls      chan *appCall
	cancels    <-chan *question
	releases   chan *outgoingRelease
	localCalls chan<- *localCall
	returns    <-chan *outgoingReturn

	// Mutable state. Only accessed from coordinate goroutine.
	questions questionTable
	answers   answerTable
	imports   importTable
	exports   exportTable
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
	conn := &Conn{transport: t}
	conn.manager.init()
	for _, o := range options {
		o(conn)
	}
	i := make(chan rpccapnp.Message)
	o := make(chan outgoingMessage)
	calls := make(chan *appCall)
	cancels := make(chan *question)
	lc := make(chan *localCall, 16)
	rets := make(chan *outgoingReturn)
	releases := make(chan *outgoingRelease)
	conn.in, conn.out = i, o
	conn.calls = calls
	conn.cancels = cancels
	conn.localCalls = lc
	conn.returns = rets
	conn.releases = releases
	conn.questions.manager = &conn.manager
	conn.questions.calls = calls
	conn.questions.cancels = cancels
	conn.answers.manager = &conn.manager
	conn.answers.returns = rets
	conn.imports.manager = &conn.manager
	conn.imports.calls = calls
	conn.imports.releases = releases

	conn.manager.do(conn.coordinate)
	conn.manager.do(func() {
		dispatchRecv(&conn.manager, t, i)
	})
	conn.manager.do(func() {
		dispatchSend(&conn.manager, t, o)
	})
	conn.manager.do(func() {
		dispatchLocalCalls(&conn.manager, lc)
	})
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

// coordinate runs in its own goroutine.
// It manages dispatching received messages and calls.
func (c *Conn) coordinate() {
	for {
		select {
		case m := <-c.in:
			c.handleMessage(m)
		case ac := <-c.calls:
			c.handleCall(ac)
		case q := <-c.cancels:
			c.handleCancel(q)
		case r := <-c.releases:
			r.echan <- c.handleRelease(r.id)
		case r := <-c.returns:
			c.handleReturn(r)
		case <-c.manager.finish:
			return
		}
	}
}

// Bootstrap returns the receiver's main interface.
func (c *Conn) Bootstrap(ctx context.Context) capnp.Client {
	// TODO(light): Create a client that returns immediately.
	ac, achan := newAppBootstrapCall(ctx)
	select {
	case c.calls <- ac:
		select {
		case a := <-achan:
			return capnp.NewPipeline(a).Client()
		case <-c.manager.finish:
			return capnp.ErrorClient(c.manager.err())
		}
	case <-c.manager.finish:
		return capnp.ErrorClient(c.manager.err())
	}
}

// handleMessage is run in the coordinate goroutine.
func (c *Conn) handleMessage(m rpccapnp.Message) {
	switch m.Which() {
	case rpccapnp.Message_Which_unimplemented:
		// no-op for now to avoid feedback loop
	case rpccapnp.Message_Which_abort:
		a := Abort{m.Abort()}
		log.Print(a)
		c.manager.shutdown(a)
	case rpccapnp.Message_Which_return:
		if err := c.handleReturnMessage(m.Return()); err != nil {
			log.Println("rpc: handle return:", err)
		}
	case rpccapnp.Message_Which_finish:
		// TODO(light): what if answers never had this ID?
		// TODO(light): return if cancelled
		id := answerID(m.Finish().QuestionId())
		releaseCaps := m.Finish().ReleaseResultCaps()
		a := c.answers.pop(id)
		a.cancel()
		if releaseCaps {
			c.exports.releaseList(a.resultCaps)
		}
	case rpccapnp.Message_Which_bootstrap:
		id := answerID(m.Bootstrap().QuestionId())
		if err := c.handleBootstrapMessage(id); err != nil {
			log.Println("rpc: handle bootstrap:", err)
		}
	case rpccapnp.Message_Which_call:
		if err := c.handleCallMessage(m); err != nil {
			log.Println("rpc: handle call:", err)
		}
	case rpccapnp.Message_Which_release:
		id := exportID(m.Release().Id())
		refs := int(m.Release().ReferenceCount())
		c.exports.release(id, refs)
	default:
		log.Printf("rpc: received unimplemented message, which = %v", m.Which())
		um := newUnimplementedMessage(nil, m)
		c.out <- outgoingMessage{c.manager.context(), um}
	}
}

func newUnimplementedMessage(buf []byte, m rpccapnp.Message) rpccapnp.Message {
	s := capnp.NewBuffer(buf)
	n := rpccapnp.NewRootMessage(s)
	n.SetUnimplemented(m)
	return n
}

// handleCall is run from the coordinate goroutine.
func (c *Conn) handleCall(ac *appCall) {
	if ac.kind == appPipelineCall && c.questions.get(ac.question.id) != ac.question {
		// Question has been finished.  The call should happen as if it is
		// back in application code.
		go func() {
			<-ac.question.resolved // not strictly necessary
			_, obj, err, _ := ac.question.peek()
			if err != nil {
				ac.achan <- capnp.ErrorAnswer(err)
				return
			}
			c := capnp.TransformObject(obj, ac.transform).ToInterface().Client()
			if c == nil {
				ac.achan <- capnp.ErrorAnswer(capnp.ErrNullClient)
				return
			}
			ac.achan <- c.Call(ac.Call)
		}()
		return
	}
	q := c.questions.new(ac.Ctx, &ac.Method)
	msg := c.newCallMessage(nil, q.id, ac)
	select {
	case c.out <- outgoingMessage{ac.Ctx, msg}:
		q.start()
		ac.achan <- q
	case <-ac.Ctx.Done():
		c.questions.pop(q.id)
		ac.achan <- capnp.ErrorAnswer(ac.Ctx.Err())
	case <-c.manager.finish:
		c.questions.pop(q.id)
		ac.achan <- capnp.ErrorAnswer(c.manager.err())
	}
}

func (c *Conn) newCallMessage(buf []byte, id questionID, ac *appCall) rpccapnp.Message {
	s := capnp.NewBuffer(buf)
	msg := rpccapnp.NewRootMessage(s)

	if ac.kind == appBootstrapCall {
		boot := rpccapnp.NewBootstrap(s)
		boot.SetQuestionId(uint32(id))
		msg.SetBootstrap(boot)
		return msg
	}

	msgCall := rpccapnp.NewCall(s)
	msgCall.SetQuestionId(uint32(id))
	msgCall.SetInterfaceId(ac.Method.InterfaceID)
	msgCall.SetMethodId(ac.Method.MethodID)

	target := rpccapnp.NewMessageTarget(s)
	switch ac.kind {
	case appImportCall:
		target.SetImportedCap(uint32(ac.importID))
	case appPipelineCall:
		a := rpccapnp.NewPromisedAnswer(s)
		a.SetQuestionId(uint32(ac.question.id))
		transformToPromisedAnswer(s, a, ac.transform)
		target.SetPromisedAnswer(a)
	default:
		panic("unknown call type")
	}
	msgCall.SetTarget(target)

	payload := rpccapnp.NewPayload(s)
	params := ac.PlaceParams(s)
	payload.SetContent(capnp.Object(params))
	payload.SetCapTable(c.makeCapTable(s))
	msgCall.SetParams(payload)

	msg.SetCall(msgCall)
	return msg
}

func transformToPromisedAnswer(s *capnp.Segment, answer rpccapnp.PromisedAnswer, transform []capnp.PipelineOp) {
	opList := rpccapnp.NewPromisedAnswer_Op_List(s, len(transform))
	for i, op := range transform {
		opList.At(i).SetGetPointerField(uint16(op.Field))
	}
	answer.SetTransform(opList)
}

// handleCancel is called from the coordinate goroutine to handle a question's cancelation.
func (c *Conn) handleCancel(q *question) {
	q.reject(questionCanceled, q.ctx.Err())
	// TODO(light): timeout?
	msg := newFinishMessage(nil, q.id, true /* release */)
	select {
	case c.out <- outgoingMessage{q.manager.context(), msg}:
	case <-c.manager.finish:
	}
}

// handleRelease is run in the coordinate goroutine to handle an import
// client's release request.  It sends a release message for an import ID.
func (c *Conn) handleRelease(id importID) error {
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
	c.out <- outgoingMessage{c.manager.context(), msg}
	return nil
}

// handleReturnMessage is run in the coordinate goroutine.
func (c *Conn) handleReturnMessage(m rpccapnp.Return) error {
	id := questionID(m.AnswerId())
	q := c.questions.pop(id)
	if q == nil {
		return fmt.Errorf("received return for unknown question id=%d", id)
	}
	if m.ReleaseParamCaps() {
		c.exports.releaseList(q.paramCaps)
	}
	if _, _, _, resolved := q.peek(); resolved {
		// If the question was already resolved, that means it was canceled,
		// in which case we already sent the finish message.
		return nil
	}
	releaseResultCaps := true
	switch m.Which() {
	case rpccapnp.Return_Which_results:
		releaseResultCaps = false
		c.populateMessageCapTable(m.Results())
		q.fulfill(m.Results().Content())
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
		q.reject(questionResolved, e)
	case rpccapnp.Return_Which_canceled:
		err := &questionError{
			id:     id,
			method: q.method,
			err:    fmt.Errorf("receiver reported canceled"),
		}
		log.Println(err)
		q.reject(questionResolved, err)
		return nil
	default:
		um := newUnimplementedMessage(nil, rpccapnp.ReadRootMessage(m.Segment))
		select {
		case c.out <- outgoingMessage{c.manager.context(), um}:
		case <-c.manager.finish:
		}
		return nil
	}
	fin := newFinishMessage(nil, id, releaseResultCaps)
	select {
	case c.out <- outgoingMessage{c.manager.context(), fin}:
	case <-c.manager.finish:
	}
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
			client := c.imports.addRef(id)
			msg.AddCap(client)
		case rpccapnp.CapDescriptor_Which_receiverHosted:
			id := exportID(desc.ReceiverHosted())
			e := c.exports.get(id)
			if e == nil {
				msg.AddCap(nil)
			} else {
				msg.AddCap(e.client)
			}
		// TODO(light): case rpccapnp.CapDescriptor_Which_receiverAnswer:
		default:
			// TODO(light): send unimplemented
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
		if isImportFromConn(client, c) {
			desc.SetReceiverHosted(uint32(client.id))
			return
		}
	case *capnp.PipelineClient:
		p := (*capnp.Pipeline)(client)
		q, ok := p.Answer().(*question)
		if !ok {
			goto fallback
		}
		id, obj, err, done := q.peek()
		if !done {
			if !isQuestionFromConn(q, c) {
				goto fallback
			}
			a := rpccapnp.NewPromisedAnswer(desc.Segment)
			a.SetQuestionId(uint32(id))
			transformToPromisedAnswer(desc.Segment, a, p.Transform())
			desc.SetReceiverAnswer(a)
			return
		}
		if err != nil {
			id := c.exports.add(capnp.ErrorClient(err))
			desc.SetSenderHosted(uint32(id))
			return
		}
		tclient := capnp.TransformObject(obj, p.Transform()).ToInterface().Client()
		if tclient == nil {
			id := c.exports.add(capnp.ErrorClient(capnp.ErrNullClient))
			desc.SetSenderHosted(uint32(id))
			return
		}
		c.descriptorForClient(desc, tclient)
		return
	}

	// Fallback: host and export ourselves.
fallback:
	id := c.exports.add(client)
	desc.SetSenderHosted(uint32(id))
}

func isQuestionFromConn(q *question, c *Conn) bool {
	// TODO(light): ideally there would be better ways to check.
	return q.manager == &c.manager
}

func isImportFromConn(ic *importClient, c *Conn) bool {
	// TODO(light): ideally there would be better ways to check.
	return ic.manager == &c.manager
}

// handleBootstrapMessage is run in the coordinate goroutine to handle
// a received bootstrap message.
func (c *Conn) handleBootstrapMessage(id answerID) error {
	ctx, cancel := c.newContext()
	a := c.answers.insert(id, cancel)
	retmsg := newReturnMessage(id)
	send := func() error {
		select {
		case c.out <- outgoingMessage{ctx, retmsg}:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		case <-c.manager.finish:
			return c.manager.err()
		}
	}
	ret := retmsg.Return()
	if a == nil {
		// Question ID reused, error out.
		setReturnException(ret, errQuestionReused)
		return send()
	}
	if c.main == nil {
		e := setReturnException(ret, errNoMainInterface)
		a.reject(Exception{e})
		return send()
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
	a.fulfill(in)
	return send()
}

// handleCallMessage is run in the coordinate goroutine to handle a
// received call message.  It mutates the capability table of its
// parameter.
func (c *Conn) handleCallMessage(m rpccapnp.Message) error {
	ctx, cancel := c.newContext()
	id := answerID(m.Call().QuestionId())
	a := c.answers.insert(id, cancel)
	if a == nil {
		// Question ID reused, error out.
		ret := newReturnMessage(id)
		setReturnException(ret.Return(), errQuestionReused)
		select {
		case c.out <- outgoingMessage{ctx, ret}:
			return nil
		case <-c.manager.finish:
			return c.manager.err()
		}
	}
	fail := func(err error) error {
		ret := newReturnMessage(id)
		setReturnException(ret.Return(), err)
		a.reject(err)
		select {
		case c.out <- outgoingMessage{ctx, ret}:
			return nil
		case <-c.manager.finish:
			return c.manager.err()
		}
	}
	c.populateMessageCapTable(m.Call().Params())
	meth := capnp.Method{
		InterfaceID: m.Call().InterfaceId(),
		MethodID:    m.Call().MethodId(),
	}
	cl := &capnp.Call{
		Ctx:    ctx,
		Method: meth,
		Params: m.Call().Params().Content().ToStruct(),
	}
	lc, err := c.newLocalCall(a, m.Call().Target(), cl)
	if err != nil {
		return fail(err)
	}
	a.start()
	select {
	case c.localCalls <- lc:
		return nil
	case <-c.manager.finish:
		err := c.manager.err()
		go a.reject(err)
		return err
	}
}

// newContext creates a new context for a local call.
func (c *Conn) newContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(c.manager.context())
}

// A localCall is the Conn-internal type for queued calls.
type localCall struct {
	result *answer
	call   *capnp.Call

	// Imported cap
	client capnp.Client

	// Promised answer
	promise   *answer
	transform []capnp.PipelineOp
}

func (c *Conn) newLocalCall(resultAnswer *answer, mt rpccapnp.MessageTarget, cl *capnp.Call) (*localCall, error) {
	switch mt.Which() {
	case rpccapnp.MessageTarget_Which_importedCap:
		id := exportID(mt.ImportedCap())
		e := c.exports.get(id)
		if e == nil {
			return nil, errBadTarget
		}
		return &localCall{
			result: resultAnswer,
			call:   cl,
			client: e.client,
		}, nil
	case rpccapnp.MessageTarget_Which_promisedAnswer:
		id := answerID(mt.PromisedAnswer().QuestionId())
		if id == resultAnswer.id {
			// Grandfather paradox.
			return nil, errBadTarget
		}
		pa := c.answers.get(id)
		if pa == nil {
			return nil, errBadTarget
		}
		transform := promisedAnswerOpsToTransform(mt.PromisedAnswer().Transform())
		return &localCall{
			result:    resultAnswer,
			call:      cl,
			promise:   pa,
			transform: transform,
		}, nil
	default:
		// TODO(light): send unimplemented
		return nil, errBadTarget
	}
}

func (lc *localCall) do() error {
	switch {
	case lc.client != nil:
		answer := lc.client.Call(lc.call)
		go joinAnswer(lc.result, answer)
		return nil
	case lc.promise != nil:
		if err := lc.promise.queueCall(lc.result, lc.transform, lc.call); err != nil {
			return err
		}
		return nil
	default:
		panic("bad localCall value")
	}
}

// dispatchLocalCalls is run in its own goroutine.
// It starts calls to local capabilities.  This occurs in a
// separate goroutine because some clients may take a long time to
// start their call or may try to send a call on the connection.
// By having a buffered queue, we can relieve these issues.
func dispatchLocalCalls(m *manager, calls <-chan *localCall) {
	for {
		select {
		case lc := <-calls:
			if err := lc.do(); err != nil {
				log.Printf("rpc: queuing %v call: %v", lc.call.Method, err)
			}
		case <-m.finish:
			return
		}
	}
}

func promisedAnswerOpsToTransform(list rpccapnp.PromisedAnswer_Op_List) []capnp.PipelineOp {
	n := list.Len()
	transform := make([]capnp.PipelineOp, 0, n)
	for i := 0; i < n; i++ {
		op := list.At(i)
		switch op.Which() {
		case rpccapnp.PromisedAnswer_Op_Which_getPointerField:
			transform = append(transform, capnp.PipelineOp{
				Field: int(op.GetPointerField()),
			})
		case rpccapnp.PromisedAnswer_Op_Which_noop:
			// no-op
		}
	}
	return transform
}

// handleReturn is called from the coordinate goroutine to send an
// answer's return value over the transport.
func (c *Conn) handleReturn(r *outgoingReturn) {
	msg := newReturnMessage(r.id)
	if r.err == nil {
		payload := rpccapnp.NewPayload(msg.Segment)
		payload.SetContent(r.obj)
		payload.SetCapTable(c.makeCapTable(msg.Segment))
		msg.Return().SetResults(payload)
	} else {
		setReturnException(msg.Return(), r.err)
	}
	select {
	case c.out <- outgoingMessage{c.manager.context(), msg}:
	case <-c.manager.finish:
	}
}

func newReturnMessage(id answerID) rpccapnp.Message {
	s := capnp.NewBuffer(nil)
	retmsg := rpccapnp.NewRootMessage(s)
	ret := rpccapnp.NewReturn(s)
	ret.SetAnswerId(uint32(id))
	ret.SetReleaseParamCaps(false)
	retmsg.SetReturn(ret)
	return retmsg
}

func setReturnException(ret rpccapnp.Return, err error) rpccapnp.Exception {
	e := rpccapnp.NewException(ret.Segment)
	toException(e, err)
	ret.SetException(e)
	return e
}

// An appCall is a message sent to the coordinate goroutine to indicate
// that the application code wants to initiate an outgoing call.
type appCall struct {
	*capnp.Call
	kind  int
	achan chan<- capnp.Answer

	// Import calls
	importID importID

	// Pipeline calls
	question  *question
	transform []capnp.PipelineOp
}

func newAppImportCall(id importID, cl *capnp.Call) (*appCall, <-chan capnp.Answer) {
	achan := make(chan capnp.Answer, 1)
	return &appCall{
		Call:     cl,
		kind:     appImportCall,
		achan:    achan,
		importID: id,
	}, achan
}

func newAppPipelineCall(q *question, transform []capnp.PipelineOp, cl *capnp.Call) (*appCall, <-chan capnp.Answer) {
	achan := make(chan capnp.Answer, 1)
	return &appCall{
		Call:      cl,
		kind:      appPipelineCall,
		achan:     achan,
		question:  q,
		transform: transform,
	}, achan
}

func newAppBootstrapCall(ctx context.Context) (*appCall, <-chan capnp.Answer) {
	achan := make(chan capnp.Answer, 1)
	return &appCall{
		Call:  &capnp.Call{Ctx: ctx},
		kind:  appBootstrapCall,
		achan: achan,
	}, achan
}

// Kinds of application calls.
const (
	appImportCall = iota
	appPipelineCall
	appBootstrapCall
)
