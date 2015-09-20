// Package rpc implements the Cap'n Proto RPC protocol.
package rpc // import "zombiezen.com/go/capnproto2/rpc"

import (
	"fmt"
	"log"

	"golang.org/x/net/context"
	"zombiezen.com/go/capnproto2"
	"zombiezen.com/go/capnproto2/rpc/rpccapnp"
)

// A Conn is a connection to another Cap'n Proto vat.
// It is safe to use from multiple goroutines.
type Conn struct {
	transport Transport
	main      capnp.Client

	manager     manager
	in          <-chan rpccapnp.Message
	out         chan<- rpccapnp.Message
	calls       chan *appCall
	cancels     <-chan *question
	releases    chan *outgoingRelease
	returns     <-chan *outgoingReturn
	queueCloses <-chan queueClientClose

	// Mutable state. Only accessed from coordinate goroutine.
	questions questionTable
	answers   answerTable
	imports   importTable
	exports   exportTable
	embargoes embargoTable
}

type connParams struct {
	main           capnp.Client
	sendBufferSize int
}

// A ConnOption is an option for opening a connection.
type ConnOption struct {
	f func(*connParams)
}

// MainInterface specifies that the connection should use client when
// receiving bootstrap messages.  By default, all bootstrap messages will
// fail.
func MainInterface(client capnp.Client) ConnOption {
	return ConnOption{func(c *connParams) {
		c.main = client
	}}
}

// SendBufferSize sets the number of outgoing messages to buffer on the
// connection.  This is in addition to whatever buffering the connection's
// transport performs.
func SendBufferSize(numMsgs int) ConnOption {
	return ConnOption{func(c *connParams) {
		c.sendBufferSize = numMsgs
	}}
}

// NewConn creates a new connection that communicates on c.
// Closing the connection will cause c to be closed.
func NewConn(t Transport, options ...ConnOption) *Conn {
	conn := &Conn{transport: t}
	p := &connParams{
		sendBufferSize: 4,
	}
	conn.manager.init()
	for _, o := range options {
		o.f(p)
	}
	conn.main = p.main
	i := make(chan rpccapnp.Message)
	o := make(chan rpccapnp.Message, p.sendBufferSize)
	calls := make(chan *appCall)
	cancels := make(chan *question)
	rets := make(chan *outgoingReturn)
	queueCloses := make(chan queueClientClose)
	releases := make(chan *outgoingRelease)
	conn.in, conn.out = i, o
	conn.calls = calls
	conn.cancels = cancels
	conn.releases = releases
	conn.returns = rets
	conn.queueCloses = queueCloses
	conn.questions.manager = &conn.manager
	conn.questions.calls = calls
	conn.questions.cancels = cancels
	conn.answers.manager = &conn.manager
	conn.answers.out = o
	conn.answers.returns = rets
	conn.answers.queueCloses = queueCloses
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
	n := newAbortMessage(nil, errShutdown)
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
			ans, err := c.handleCall(ac)
			if err == nil {
				ac.achan <- ans
			} else {
				log.Println("rpc: failed to handle call:", err)
				ac.achan <- capnp.ErrorAnswer(err)
			}
		case q := <-c.cancels:
			c.handleCancel(q)
		case r := <-c.releases:
			r.echan <- c.handleRelease(r.id)
		case r := <-c.returns:
			c.handleReturn(r)
		case qcc := <-c.queueCloses:
			c.handleQueueClose(qcc)
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
		case <-ctx.Done():
			return capnp.ErrorClient(ctx.Err())
		case <-c.manager.finish:
			return capnp.ErrorClient(c.manager.err())
		}
	case <-ctx.Done():
		return capnp.ErrorClient(ctx.Err())
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
		ma, err := m.Abort()
		if err != nil {
			log.Println("rpc: decode abort:", err)
			// Keep going, since we're trying to abort anyway.
		}
		a := Abort{ma}
		log.Print(a)
		c.manager.shutdown(a)
	case rpccapnp.Message_Which_return:
		if err := c.handleReturnMessage(m); err != nil {
			log.Println("rpc: handle return:", err)
		}
	case rpccapnp.Message_Which_finish:
		// TODO(light): what if answers never had this ID?
		// TODO(light): return if cancelled
		mfin, err := m.Finish()
		if err != nil {
			log.Println("rpc: decode finish:", err)
			return
		}
		id := answerID(mfin.QuestionId())
		a := c.answers.pop(id)
		a.cancel()
		if mfin.ReleaseResultCaps() {
			c.exports.releaseList(a.resultCaps)
		}
	case rpccapnp.Message_Which_bootstrap:
		boot, err := m.Bootstrap()
		if err != nil {
			log.Println("rpc: decode bootstrap:", err)
			return
		}
		id := answerID(boot.QuestionId())
		if err := c.handleBootstrapMessage(id); err != nil {
			log.Println("rpc: handle bootstrap:", err)
		}
	case rpccapnp.Message_Which_call:
		if err := c.handleCallMessage(m); err != nil {
			log.Println("rpc: handle call:", err)
		}
	case rpccapnp.Message_Which_release:
		rel, err := m.Release()
		if err != nil {
			log.Println("rpc: decode release:", err)
			return
		}
		id := exportID(rel.Id())
		refs := int(rel.ReferenceCount())
		c.exports.release(id, refs)
	case rpccapnp.Message_Which_disembargo:
		if err := c.handleDisembargoMessage(m); err != nil {
			// Any failure in a disembargo is a protocol violation.
			c.abort(err)
		}
	default:
		log.Printf("rpc: received unimplemented message, which = %v", m.Which())
		um := newUnimplementedMessage(nil, m)
		c.sendMessage(um)
	}
}

func newUnimplementedMessage(buf []byte, m rpccapnp.Message) rpccapnp.Message {
	n := newMessage(buf)
	n.SetUnimplemented(m)
	return n
}

// handleCall is run from the coordinate goroutine to send a question to a remote vat.
func (c *Conn) handleCall(ac *appCall) (capnp.Answer, error) {
	if ac.kind == appPipelineCall && c.questions.get(ac.question.id) != ac.question {
		// Question has been finished.  The call should happen as if it is
		// back in application code.
		_, obj, err, done := ac.question.peek()
		if !done {
			panic("question popped but not done")
		}
		client := clientFromResolution(ac.transform, obj, err)
		return c.nestedCall(client, ac.Call), nil
	}
	q := c.questions.new(ac.Ctx, &ac.Method)
	if ac.kind == appPipelineCall {
		pq := c.questions.get(ac.question.id)
		pq.addPromise(ac.transform)
	}
	msg, err := c.newCallMessage(nil, q.id, ac)
	if err != nil {
		return nil, err
	}
	select {
	case c.out <- msg:
		q.start()
		return q, nil
	case <-ac.Ctx.Done():
		c.questions.pop(q.id)
		return nil, ac.Ctx.Err()
	case <-c.manager.finish:
		c.questions.pop(q.id)
		return nil, c.manager.err()
	}
}

func (c *Conn) newCallMessage(buf []byte, id questionID, ac *appCall) (rpccapnp.Message, error) {
	msg := newMessage(buf)

	if ac.kind == appBootstrapCall {
		boot, _ := msg.NewBootstrap()
		boot.SetQuestionId(uint32(id))
		return msg, nil
	}

	msgCall, _ := msg.NewCall()
	msgCall.SetQuestionId(uint32(id))
	msgCall.SetInterfaceId(ac.Method.InterfaceID)
	msgCall.SetMethodId(ac.Method.MethodID)

	target, _ := msgCall.NewTarget()
	switch ac.kind {
	case appImportCall:
		target.SetImportedCap(uint32(ac.importID))
	case appPipelineCall:
		a, err := target.NewPromisedAnswer()
		if err != nil {
			return rpccapnp.Message{}, err
		}
		a.SetQuestionId(uint32(ac.question.id))
		err = transformToPromisedAnswer(a.Segment(), a, ac.transform)
		if err != nil {
			return rpccapnp.Message{}, err
		}
	default:
		panic("unknown call type")
	}

	payload, _ := msgCall.NewParams()
	params, err := ac.PlaceParams(payload.Segment())
	if err != nil {
		return rpccapnp.Message{}, err
	}
	if err := payload.SetContent(params); err != nil {
		return rpccapnp.Message{}, err
	}
	ctab, err := c.makeCapTable(payload.Segment())
	if err != nil {
		return rpccapnp.Message{}, err
	}
	if err := payload.SetCapTable(ctab); err != nil {
		return rpccapnp.Message{}, err
	}

	return msg, nil
}

func transformToPromisedAnswer(s *capnp.Segment, answer rpccapnp.PromisedAnswer, transform []capnp.PipelineOp) error {
	opList, err := rpccapnp.NewPromisedAnswer_Op_List(s, int32(len(transform)))
	if err != nil {
		return err
	}
	for i, op := range transform {
		opList.At(i).SetGetPointerField(uint16(op.Field))
	}
	err = answer.SetTransform(opList)
	return err
}

// handleCancel is called from the coordinate goroutine to handle a question's cancelation.
func (c *Conn) handleCancel(q *question) {
	q.reject(questionCanceled, q.ctx.Err())
	// TODO(light): timeout?
	msg := newFinishMessage(nil, q.id, true /* release */)
	c.sendMessage(msg)
}

// handleRelease is run in the coordinate goroutine to handle an import
// client's release request.  It sends a release message for an import ID.
func (c *Conn) handleRelease(id importID) error {
	i := c.imports.pop(id)
	if i == 0 {
		return nil
	}
	// TODO(light): deadline to close?
	msg := newMessage(nil)
	mr, err := msg.NewRelease()
	if err != nil {
		return err
	}
	mr.SetId(uint32(id))
	mr.SetReferenceCount(uint32(i))
	return c.sendMessage(msg)
}

// handleReturnMessage is run in the coordinate goroutine.
func (c *Conn) handleReturnMessage(m rpccapnp.Message) error {
	ret, err := m.Return()
	if err != nil {
		return err
	}
	id := questionID(ret.AnswerId())
	q := c.questions.pop(id)
	if q == nil {
		return fmt.Errorf("received return for unknown question id=%d", id)
	}
	if ret.ReleaseParamCaps() {
		c.exports.releaseList(q.paramCaps)
	}
	if _, _, _, resolved := q.peek(); resolved {
		// If the question was already resolved, that means it was canceled,
		// in which case we already sent the finish message.
		return nil
	}
	releaseResultCaps := true
	switch ret.Which() {
	case rpccapnp.Return_Which_results:
		releaseResultCaps = false
		results, err := ret.Results()
		if err != nil {
			return err
		}
		if err := c.populateMessageCapTable(results); err == errUnimplemented {
			um := newUnimplementedMessage(nil, m)
			c.sendMessage(um)
			return errUnimplemented
		} else if err != nil {
			c.abort(err)
			return err
		}
		content, err := results.Content()
		if err != nil {
			return err
		}
		disembargoes := q.fulfill(content, c.embargoes.new)
		for _, d := range disembargoes {
			if err := c.sendMessage(d); err != nil {
				// shutdown
				return nil
			}
		}
	case rpccapnp.Return_Which_exception:
		exc, err := ret.Exception()
		if err != nil {
			return err
		}
		e := error(Exception{exc})
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
		um := newUnimplementedMessage(nil, m)
		c.sendMessage(um)
		return errUnimplemented
	}
	fin := newFinishMessage(nil, id, releaseResultCaps)
	c.sendMessage(fin)
	return nil
}

func newFinishMessage(buf []byte, questionID questionID, release bool) rpccapnp.Message {
	m := newMessage(buf)
	f, _ := m.NewFinish()
	f.SetQuestionId(uint32(questionID))
	f.SetReleaseResultCaps(release)
	return m
}

// populateMessageCapTable converts the descriptors in the payload into
// clients and sets it on the message the payload is a part of.
func (c *Conn) populateMessageCapTable(payload rpccapnp.Payload) error {
	msg := payload.Segment().Message()
	ctab, err := payload.CapTable()
	if err != nil {
		return err
	}
	for i, n := 0, ctab.Len(); i < n; i++ {
		desc := ctab.At(i)
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
				return fmt.Errorf("rpc: capability table references unknown export ID %d", id)
			}
			msg.AddCap(e.client)
		case rpccapnp.CapDescriptor_Which_receiverAnswer:
			recvAns, err := desc.ReceiverAnswer()
			if err != nil {
				return err
			}
			id := answerID(recvAns.QuestionId())
			a := c.answers.get(id)
			if a == nil {
				return fmt.Errorf("rpc: capability table references unknown answer ID %d", id)
			}
			recvTransform, err := recvAns.Transform()
			if err != nil {
				return err
			}
			transform := promisedAnswerOpsToTransform(recvTransform)
			msg.AddCap(a.pipelineClient(transform))
		default:
			log.Println("rpc: unknown capability type", desc.Which())
			return errUnimplemented
		}
	}
	return nil
}

// makeCapTable converts the clients in the segment's message into capability descriptors.
func (c *Conn) makeCapTable(s *capnp.Segment) (rpccapnp.CapDescriptor_List, error) {
	msgtab := s.Message().CapTable
	t, err := rpccapnp.NewCapDescriptor_List(s, int32(len(msgtab)))
	if err != nil {
		return rpccapnp.CapDescriptor_List{}, nil
	}
	for i, client := range msgtab {
		desc := t.At(i)
		if client == nil {
			desc.SetNone()
			continue
		}
		c.descriptorForClient(desc, client)
	}
	return t, nil
}

// handleBootstrapMessage is run in the coordinate goroutine to handle
// a received bootstrap message.
func (c *Conn) handleBootstrapMessage(id answerID) error {
	a := c.answers.insert(id, func() {})
	if a == nil {
		// Question ID reused, error out.
		retmsg := newReturnMessage(nil, id)
		r, _ := retmsg.Return()
		setReturnException(r, errQuestionReused)
		return c.sendMessage(retmsg)
	}
	msgs := make([]rpccapnp.Message, 0, 1)
	if c.main == nil {
		msgs = a.reject(msgs, errNoMainInterface)
		for _, m := range msgs {
			if err := c.sendMessage(m); err != nil {
				return err
			}
		}
		return nil
	}
	m := &capnp.Message{
		Arena:    capnp.SingleSegment(make([]byte, 0)),
		CapTable: []capnp.Client{c.main},
	}
	s, _ := m.Segment(0)
	in := capnp.NewInterface(s, 0)
	msgs = a.fulfill(msgs, in, c.makeCapTable)
	for _, m := range msgs {
		if err := c.sendMessage(m); err != nil {
			return err
		}
	}
	return nil
}

// handleCallMessage is run in the coordinate goroutine to handle a
// received call message.  It mutates the capability table of its
// parameter.
func (c *Conn) handleCallMessage(m rpccapnp.Message) error {
	mcall, err := m.Call()
	if err != nil {
		return err
	}
	mt, err := mcall.Target()
	if err != nil {
		return err
	}
	if mt.Which() != rpccapnp.MessageTarget_Which_importedCap && mt.Which() != rpccapnp.MessageTarget_Which_promisedAnswer {
		um := newUnimplementedMessage(nil, m)
		return c.sendMessage(um)
	}
	mparams, err := mcall.Params()
	if err != nil {
		return err
	}
	if err := c.populateMessageCapTable(mparams); err == errUnimplemented {
		um := newUnimplementedMessage(nil, m)
		return c.sendMessage(um)
	} else if err != nil {
		c.abort(err)
		return err
	}
	ctx, cancel := c.newContext()
	id := answerID(mcall.QuestionId())
	a := c.answers.insert(id, cancel)
	if a == nil {
		// Question ID reused, error out.
		c.abort(errQuestionReused)
		return errQuestionReused
	}
	meth := capnp.Method{
		InterfaceID: mcall.InterfaceId(),
		MethodID:    mcall.MethodId(),
	}
	paramContent, err := mparams.Content()
	if err != nil {
		return err
	}
	cl := &capnp.Call{
		Ctx:    ctx,
		Method: meth,
		Params: capnp.ToStruct(paramContent),
	}
	if err := c.routeCallMessage(a, mt, cl); err != nil {
		msgs := a.reject(nil, err)
		for _, m := range msgs {
			if err := c.sendMessage(m); err != nil {
				return err
			}
		}
		return nil
	}
	return nil
}

func (c *Conn) routeCallMessage(result *answer, mt rpccapnp.MessageTarget, cl *capnp.Call) error {
	switch mt.Which() {
	case rpccapnp.MessageTarget_Which_importedCap:
		id := exportID(mt.ImportedCap())
		e := c.exports.get(id)
		if e == nil {
			return errBadTarget
		}
		answer := c.nestedCall(e.client, cl)
		go joinAnswer(result, answer)
	case rpccapnp.MessageTarget_Which_promisedAnswer:
		mpromise, err := mt.PromisedAnswer()
		if err != nil {
			return err
		}
		id := answerID(mpromise.QuestionId())
		if id == result.id {
			// Grandfather paradox.
			return errBadTarget
		}
		pa := c.answers.get(id)
		if pa == nil {
			return errBadTarget
		}
		mtrans, err := mpromise.Transform()
		if err != nil {
			return err
		}
		transform := promisedAnswerOpsToTransform(mtrans)
		if obj, err, done := pa.peek(); done {
			client := clientFromResolution(transform, obj, err)
			answer := c.nestedCall(client, cl)
			go joinAnswer(result, answer)
			return nil
		}
		if err := pa.queueCall(result, transform, cl); err != nil {
			return err
		}
	default:
		panic("unreachable")
	}
	return nil
}

func (c *Conn) handleDisembargoMessage(msg rpccapnp.Message) error {
	d, err := msg.Disembargo()
	if err != nil {
		return err
	}
	dtarget, err := d.Target()
	if err != nil {
		return err
	}
	switch d.Context().Which() {
	case rpccapnp.Disembargo_context_Which_senderLoopback:
		id := embargoID(d.Context().SenderLoopback())
		if dtarget.Which() != rpccapnp.MessageTarget_Which_promisedAnswer {
			return errDisembargoNonImport
		}
		dpa, err := dtarget.PromisedAnswer()
		if err != nil {
			return err
		}
		aid := answerID(dpa.QuestionId())
		a := c.answers.get(aid)
		if a == nil {
			return errDisembargoMissingAnswer
		}
		dtrans, err := dpa.Transform()
		if err != nil {
			return err
		}
		transform := promisedAnswerOpsToTransform(dtrans)
		queued, err := a.queueDisembargo(transform, id, dtarget)
		if err != nil {
			return err
		}
		if !queued {
			// There's nothing to embargo; everything's been delivered.
			resp := newDisembargoMessage(nil, rpccapnp.Disembargo_context_Which_receiverLoopback, id)
			rd, _ := resp.Disembargo()
			if err := rd.SetTarget(dtarget); err != nil {
				return err
			}
			c.sendMessage(resp)
		}
	case rpccapnp.Disembargo_context_Which_receiverLoopback:
		id := embargoID(d.Context().ReceiverLoopback())
		c.embargoes.disembargo(id)
	default:
		um := newUnimplementedMessage(nil, msg)
		c.sendMessage(um)
	}
	return nil
}

// newDisembargoMessage creates a disembargo message.  Its target will be left blank.
func newDisembargoMessage(buf []byte, which rpccapnp.Disembargo_context_Which, id embargoID) rpccapnp.Message {
	msg := newMessage(buf)
	d, _ := msg.NewDisembargo()
	switch which {
	case rpccapnp.Disembargo_context_Which_senderLoopback:
		d.Context().SetSenderLoopback(uint32(id))
	case rpccapnp.Disembargo_context_Which_receiverLoopback:
		d.Context().SetReceiverLoopback(uint32(id))
	default:
		panic("unreachable")
	}
	return msg
}

// newContext creates a new context for a local call.
func (c *Conn) newContext() (context.Context, context.CancelFunc) {
	return context.WithCancel(c.manager.context())
}

func promisedAnswerOpsToTransform(list rpccapnp.PromisedAnswer_Op_List) []capnp.PipelineOp {
	n := list.Len()
	transform := make([]capnp.PipelineOp, 0, n)
	for i := 0; i < n; i++ {
		op := list.At(i)
		switch op.Which() {
		case rpccapnp.PromisedAnswer_Op_Which_getPointerField:
			transform = append(transform, capnp.PipelineOp{
				Field: op.GetPointerField(),
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
	msgs := make([]rpccapnp.Message, 0, 32)
	if r.err == nil {
		msgs = r.a.fulfill(msgs, r.obj, c.makeCapTable)
	} else {
		msgs = r.a.reject(msgs, r.err)
	}
	for _, m := range msgs {
		if err := c.sendMessage(m); err != nil {
			return
		}
	}
}

func (c *Conn) handleQueueClose(qcc queueClientClose) {
	msgs := make([]rpccapnp.Message, 0, 32)
	msgs = qcc.qc.rejectQueue(msgs)
	close(qcc.done)
	for _, m := range msgs {
		if err := c.sendMessage(m); err != nil {
			return
		}
	}
}

func (c *Conn) sendMessage(msg rpccapnp.Message) error {
	return sendMessage(&c.manager, c.out, msg)
}

func (c *Conn) abort(err error) {
	// TODO(light): ensure that the message is sent before shutting down?
	am := newAbortMessage(nil, err)
	c.sendMessage(am)
	c.manager.shutdown(err)
}

func newAbortMessage(buf []byte, err error) rpccapnp.Message {
	n := newMessage(buf)
	e, _ := n.NewAbort()
	toException(e, err)
	return n
}

func newReturnMessage(buf []byte, id answerID) rpccapnp.Message {
	retmsg := newMessage(buf)
	ret, _ := retmsg.NewReturn()
	ret.SetAnswerId(uint32(id))
	ret.SetReleaseParamCaps(false)
	return retmsg
}

func setReturnException(ret rpccapnp.Return, err error) rpccapnp.Exception {
	e, _ := rpccapnp.NewException(ret.Segment())
	toException(e, err)
	ret.SetException(e)
	return e
}

// clientFromResolution retrieves a client from a resolved question or
// answer by applying a transform.
func clientFromResolution(transform []capnp.PipelineOp, obj capnp.Pointer, err error) capnp.Client {
	if err != nil {
		return capnp.ErrorClient(err)
	}
	out, err := capnp.Transform(obj, transform)
	if err != nil {
		return capnp.ErrorClient(err)
	}
	c := capnp.ToInterface(out).Client()
	if c == nil {
		return capnp.ErrorClient(capnp.ErrNullClient)
	}
	return c
}

func newMessage(buf []byte) rpccapnp.Message {
	_, s, err := capnp.NewMessage(capnp.SingleSegment(buf))
	if err != nil {
		panic(err)
	}
	m, err := rpccapnp.NewRootMessage(s)
	if err != nil {
		panic(err)
	}
	return m
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
