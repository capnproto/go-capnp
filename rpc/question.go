package rpc

import (
	"context"

	"capnproto.org/go/capnp/v3"
	rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"
)

// A questionID is an index into the questions table.
type questionID uint32

type question struct {
	c  *Conn
	id questionID

	bootstrapPromise *capnp.ClientPromise

	p       *capnp.Promise
	release capnp.ReleaseFunc // written before resolving p

	// Protected by c.mu:

	flags         questionFlags
	finishMsgSend chan struct{}        // closed after attempting to send the Finish message
	called        [][]capnp.PipelineOp // paths to called clients
}

// questionFlags is a bitmask of which events have occurred in a question's
// lifetime.
type questionFlags uint8

const (
	// finished is set when the question's Context has been canceled or
	// its Return message has been received.  The codepath that sets this
	// flag is responsible for sending the Finish message.
	finished questionFlags = 1 << iota

	// finishSent indicates whether the Finish message was sent
	// successfully.  It is only valid to query after finishMsgSend is
	// closed.
	finishSent
)

// newQuestion adds a new question to c's table.  The caller must be
// holding onto c.mu.
func (c *Conn) newQuestion(method capnp.Method) *question {
	q := &question{
		c:             c,
		id:            questionID(c.questionID.next()),
		finishMsgSend: make(chan struct{}),
	}
	q.p = capnp.NewPromise(method, q) // TODO(someday): customize error message for bootstrap
	if int(q.id) == len(c.questions) {
		c.questions = append(c.questions, q)
	} else {
		c.questions[q.id] = q
	}
	return q
}

// handleCancel rejects the question's promise upon cancelation of its
// Context.
//
// The caller must not be holding onto q.c.mu or the sender lock.
func (q *question) handleCancel(ctx context.Context) {
	var rejectErr error
	select {
	case <-ctx.Done():
		rejectErr = ctx.Err()
	case <-q.c.bgctx.Done():
		rejectErr = ExcClosed
	case <-q.p.Answer().Done():
		return
	}

	q.c.mu.Lock()
	if q.flags&finished != 0 {
		// Promise already fulfilled.
		q.c.mu.Unlock()
		return
	}
	q.flags |= finished
	q.release = func() {}
	err := q.c.sendMessage(q.c.bgctx, func(msg rpccp.Message) error {
		fin, err := msg.NewFinish()
		if err != nil {
			return err
		}
		fin.SetQuestionId(uint32(q.id))
		fin.SetReleaseResultCaps(true)
		return nil
	})
	if err == nil {
		q.flags |= finishSent
	} else {
		select {
		case <-q.c.bgctx.Done():
		default:
			q.c.report(annotate(err, "send finish"))
		}
	}
	close(q.finishMsgSend)
	q.c.mu.Unlock()

	q.p.Reject(rejectErr)
	if q.bootstrapPromise != nil {
		q.bootstrapPromise.Fulfill(q.p.Answer().Client())
		q.p.ReleaseClients()
	}
}

func (q *question) PipelineSend(ctx context.Context, transform []capnp.PipelineOp, s capnp.Send) (*capnp.Answer, capnp.ReleaseFunc) {
	// Acquire sender lock.
	q.c.mu.Lock()
	if !q.c.startTask() {
		q.c.mu.Unlock()
		return capnp.ErrorAnswer(s.Method, ExcClosed), func() {}
	}
	defer q.c.tasks.Done()
	// Mark this transform as having been used for a call ASAP.
	// q's Return could be received while q2 is being sent.
	// Don't bother cleaning it up if the call fails because:
	// a) this may not have been the only call for the given transform,
	// b) the transform isn't guaranteed to be an import, and
	// c) the worst that happens is we trade bandwidth for code simplicity.
	q.mark(transform)
	if err := q.c.tryLockSender(ctx); err != nil {
		q.c.mu.Unlock()
		return capnp.ErrorAnswer(s.Method, err), func() {}
	}
	q2 := q.c.newQuestion(s.Method)
	q.c.mu.Unlock()

	// Create call message.
	msg, send, release, err := q.c.transport.NewMessage(ctx)
	if err != nil {
		q.c.mu.Lock()
		q.c.questions[q2.id] = nil
		q.c.questionID.remove(uint32(q2.id))
		q.c.mu.Unlock()
		return capnp.ErrorAnswer(s.Method, failedf("create message: %w", err)), func() {}
	}
	q.c.mu.Lock()
	q.c.unlockSender() // Can't be holding either lock while calling PlaceArgs.
	q.c.mu.Unlock()
	err = q.c.newPipelineCallMessage(msg, q.id, transform, q2.id, s)
	if err != nil {
		q.c.mu.Lock()
		q.c.questions[q2.id] = nil
		q.c.questionID.remove(uint32(q2.id))
		q.c.lockSender()
		q.c.mu.Unlock()
		release()
		q.c.mu.Lock()
		q.c.unlockSender()
		q.c.mu.Unlock()
		return capnp.ErrorAnswer(s.Method, err), func() {}
	}

	// Send call.
	q.c.mu.Lock()
	q.c.lockSender()
	q.c.mu.Unlock()
	err = send()
	release()

	q.c.mu.Lock()
	q.c.unlockSender()
	if err != nil {
		q.c.questions[q2.id] = nil
		q.c.questionID.remove(uint32(q2.id))
		q.c.mu.Unlock()
		return capnp.ErrorAnswer(s.Method, failedf("send message: %w", err)), func() {}
	}
	q2.c.tasks.Add(1)
	go func() {
		defer q2.c.tasks.Done()
		q2.handleCancel(ctx)
	}()
	q.c.mu.Unlock()

	ans := q2.p.Answer()
	return ans, func() {
		<-ans.Done()
		q2.p.ReleaseClients()
		q2.release()
	}
}

// newPipelineCallMessage builds a Call message targeted to a promised answer..
//
// The caller MUST NOT be holding onto c.mu or the sender lock.
func (c *Conn) newPipelineCallMessage(msg rpccp.Message, tgt questionID, transform []capnp.PipelineOp, qid questionID, s capnp.Send) error {
	call, err := msg.NewCall()
	if err != nil {
		return failedf("build call message: %w", err)
	}
	call.SetQuestionId(uint32(qid))
	call.SetInterfaceId(s.Method.InterfaceID)
	call.SetMethodId(s.Method.MethodID)

	target, err := call.NewTarget()
	if err != nil {
		return failedf("build call message: %w", err)
	}
	pa, err := target.NewPromisedAnswer()
	if err != nil {
		return failedf("build call message: %w", err)
	}
	pa.SetQuestionId(uint32(tgt))
	oplist, err := pa.NewTransform(int32(len(transform)))
	if err != nil {
		return failedf("build call message: %w", err)
	}
	for i, op := range transform {
		oplist.At(i).SetGetPointerField(op.Field)
	}

	payload, err := call.NewParams()
	if err != nil {
		return failedf("build call message: %w", err)
	}
	args, err := capnp.NewStruct(payload.Segment(), s.ArgsSize)
	if err != nil {
		return failedf("build call message: %w", err)
	}
	if err := payload.SetContent(args.ToPtr()); err != nil {
		return failedf("build call message: %w", err)
	}

	if s.PlaceArgs == nil {
		return nil
	}
	m := args.Message()
	if err := s.PlaceArgs(args); err != nil {
		for _, c := range m.CapTable {
			c.Release()
		}
		m.CapTable = nil
		return failedf("place arguments: %w", err)
	}
	clients := extractCapTable(m)
	c.mu.Lock()
	// TODO(soon): save param refs
	_, err = c.fillPayloadCapTable(payload, clients)
	c.mu.Unlock()
	releaseList(clients).release()
	if err != nil {
		return annotate(err, "build call message")
	}
	return nil
}

func (q *question) PipelineRecv(ctx context.Context, transform []capnp.PipelineOp, r capnp.Recv) capnp.PipelineCaller {
	ans, finish := q.PipelineSend(ctx, transform, capnp.Send{
		Method:   r.Method,
		ArgsSize: r.Args.Size(),
		PlaceArgs: func(s capnp.Struct) error {
			err := s.CopyFrom(r.Args)
			r.ReleaseArgs()
			return err
		},
	})
	r.ReleaseArgs()
	select {
	case <-ans.Done():
		returnAnswer(r.Returner, ans, finish)
		return nil
	default:
		go returnAnswer(r.Returner, ans, finish)
		return ans
	}
}

// mark adds the promised answer transform to the set of pipelined
// questions sent.  The caller must be holding onto q.c.mu.
func (q *question) mark(xform []capnp.PipelineOp) {
	for _, x := range q.called {
		if transformsEqual(x, xform) {
			// Already in set.
			return
		}
	}
	// Add a copy (don't retain default values).
	xform2 := make([]capnp.PipelineOp, len(xform))
	for i := range xform {
		xform2[i].Field = xform[i].Field
	}
	q.called = append(q.called, xform2)
}

func transformsEqual(x1, x2 []capnp.PipelineOp) bool {
	if len(x1) != len(x2) {
		return false
	}
	for i := range x1 {
		if x1[i].Field != x2[i].Field {
			return false
		}
	}
	return true
}
