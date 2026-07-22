package rpc

import (
	"context"
	"errors"
	"sync"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/flowcontrol"
	"capnproto.org/go/capnp/v3/internal/syncutil"
	rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"
)

// A questionID is an index into the questions table.
type questionID uint32

type question struct {
	c  *Conn
	id questionID

	p       *capnp.Promise
	release capnp.ReleaseFunc // written before resolving p

	// Protected by c.lk:
	owner          questionTerminalOwner
	callSend       questionCallSendState
	callSendDone   chan struct{}
	callSendErr    error
	drainComplete  bool
	drainDone      chan struct{}
	returnReceived bool
	noFinishNeeded bool
	finish         questionFinishState
	called         [][]capnp.PipelineOp // paths to called clients
}

type questionTerminalOwner uint8

const (
	questionOwnerNone questionTerminalOwner = iota
	questionOwnerReturn
	questionOwnerCancel
	questionOwnerSendFailure
)

type questionFinishState uint8

const (
	questionFinishNotQueued questionFinishState = iota
	questionFinishQueued
	questionFinishSent
	questionFinishFailed
	questionFinishSuppressed
)

type questionCallSendState uint8

const (
	questionCallSendPending questionCallSendState = iota
	questionCallSendSucceeded
	questionCallSendFailed
)

// newQuestion adds a new question to c's table.
func (c *lockedConn) newQuestion(method capnp.Method) *question {
	q := &question{
		c:            (*Conn)(c),
		release:      func() {},
		callSendDone: make(chan struct{}),
		drainDone:    make(chan struct{}),
	}
	q.id = c.lk.questions.Add(q)
	q.p = capnp.NewPromise(method, q, nil) // TODO(someday): customize error message for bootstrap
	c.setAnswerQuestion(q.p.Answer(), q)
	return q
}

func (c *lockedConn) getAnswerQuestion(ans *capnp.Answer) (*question, bool) {
	m := ans.Metadata()
	m.Lock()
	defer m.Unlock()
	q, ok := m.Get(questionKey{(*Conn)(c)})
	if !ok {
		return nil, false
	}
	return q.(*question), true
}

func (c *lockedConn) setAnswerQuestion(ans *capnp.Answer, q *question) {
	m := ans.Metadata()
	syncutil.With(m, func() {
		m.Put(questionKey{(*Conn)(c)}, q)
	})
}

type questionKey struct {
	conn *Conn
}

// handleCancel rejects the question's promise upon cancelation of its
// Context.
//
// The caller MUST NOT hold q.c.lk.
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

	claimed := false
	q.c.withLocked(func(c *lockedConn) {
		if q.owner == questionOwnerNone {
			q.owner = questionOwnerCancel
			q.release = func() {}
			claimed = true
		}
	})
	if !claimed {
		return
	}

	// Reject drains every pipeline call that selected q before cancellation
	// claimed terminal ownership. No Finish may be queued before it returns.
	q.p.Reject(rejectErr)
	q.c.withLocked(func(c *lockedConn) {
		q.drainComplete = true
		q.queueFinish(c, true)
		q.signalDrainDone()
	})
}

func (q *question) signalDrainDone() {
	if q.drainDone != nil {
		close(q.drainDone)
	}
}

// handleCallSend completes the initial Call send. The caller MUST NOT hold
// q.c.lk; it is invoked by the sender loop.
func (q *question) handleCallSend(ctx context.Context, outcome sendOutcome, annotation string) {
	var callErr error
	if outcome.err != nil {
		callErr = rpcerr.Annotate(outcome.err, annotation)
	} else if outcome.disposition != sendSucceeded {
		callErr = ExcClosed
	}
	q.c.withLocked(func(c *lockedConn) {
		if q.callSend != questionCallSendPending {
			return
		}
		if outcome.disposition == sendSucceeded {
			q.callSend = questionCallSendSucceeded
		} else {
			q.callSend = questionCallSendFailed
			q.callSendErr = callErr
		}
		close(q.callSendDone)
	})

	switch outcome.disposition {
	case sendSucceeded:
		startCancel := false
		q.c.withLocked(func(c *lockedConn) {
			if q.owner == questionOwnerNone {
				startCancel = c.startTask()
			}
		})
		if startCancel {
			go func() {
				defer q.c.tasks.Done()
				q.handleCancel(ctx)
			}()
		}

	case sendDefinitelyUnsent:
		reject := false
		q.c.withLocked(func(c *lockedConn) {
			if current, ok := c.lk.questions.Find(q.id); ok && current == q && q.owner == questionOwnerNone {
				q.owner = questionOwnerSendFailure
				reject = true
			}
		})
		if reject {
			// callSendDone wakes admitted old-route calls. They must yield
			// without enqueueing before rejection completes and the ID is freed.
			q.p.Reject(callErr)
			q.c.withLocked(func(c *lockedConn) {
				if current, ok := c.lk.questions.Find(q.id); ok && current == q {
					c.lk.questions.Remove(q.id)
				}
			})
		}

	case sendDeliveryAmbiguous:
		reject := false
		q.c.withLocked(func(c *lockedConn) {
			if current, ok := c.lk.questions.Find(q.id); ok && current == q && q.owner == questionOwnerNone {
				q.owner = questionOwnerSendFailure
				reject = true
			}
		})
		if reject {
			q.p.Reject(callErr)
		}

	case sendAbortedByShutdown:
		// callSendDone was closed above so admitted pipeline calls can yield;
		// shutdown owns Promise rejection and table cleanup.
	}
}

// queueFinish queues the question's single terminal Finish after Promise
// admission has drained. The caller MUST hold q.c.lk.
func (q *question) queueFinish(c *lockedConn, releaseResultCaps bool) {
	if !q.drainComplete || q.finish != questionFinishNotQueued {
		return
	}
	if c.lk.closing {
		q.finish = questionFinishFailed
		return
	}
	q.finish = questionFinishQueued
	c.sendMessageOutcome(c.bgctx, func(m rpccp.Message) error {
		fin, err := m.NewFinish()
		if err == nil {
			fin.SetQuestionId(uint32(q.id))
			fin.SetReleaseResultCaps(releaseResultCaps)
		}
		return err
	}, true, func(outcome sendOutcome) {
		q.c.withLocked(func(c *lockedConn) {
			if q.finish != questionFinishQueued {
				return
			}
			if outcome.disposition == sendSucceeded {
				q.finish = questionFinishSent
				if q.returnReceived {
					c.lk.questions.release(q.id)
				}
			} else {
				q.finish = questionFinishFailed
			}
		})
	})
}

func (q *question) PipelineSend(ctx context.Context, transform []capnp.PipelineOp, s capnp.Send) (*capnp.Answer, capnp.ReleaseFunc) {
	select {
	case <-q.callSendDone:
	case <-ctx.Done():
		return capnp.ErrorAnswer(s.Method, ctx.Err()), func() {}
	case <-q.c.bgctx.Done():
		return capnp.ErrorAnswer(s.Method, ExcClosed), func() {}
	}
	return withLockedConn2(q.c, func(c *lockedConn) (*capnp.Answer, capnp.ReleaseFunc) {
		if q.callSend != questionCallSendSucceeded {
			return capnp.ErrorAnswer(s.Method, q.callSendErr), func() {}
		}
		// Mark this transform as having been used for a call ASAP.
		// q's Return could be received while q2 is being sent.
		// Don't bother cleaning it up if the call fails because:
		// a) this may not have been the only call for the given transform,
		// b) the transform isn't guaranteed to be an import, and
		// c) the worst that happens is we trade bandwidth for code simplicity.
		q.mark(transform)
		return c.startCall(ctx, s, nil, func(target rpccp.MessageTarget) error {
			pa, err := target.NewPromisedAnswer()
			if err != nil {
				return rpcerr.WrapFailed("build call message", err)
			}
			pa.SetQuestionId(uint32(q.id))
			oplist, err := pa.NewTransform(int32(len(transform)))
			if err != nil {
				return rpcerr.WrapFailed("build call message", err)
			}
			for i, op := range transform {
				oplist.At(i).SetGetPointerField(op.Field)
			}
			return nil
		})
	})
}

// PreparePipelineSend preserves the same initial-call admission ordering as
// PipelineSend.  In particular, a pipelined Call is never constructed or
// enqueued until its parent Call has been accepted by the sender.
func (q *question) PreparePipelineSend(ctx context.Context, transform []capnp.PipelineOp, s capnp.Send) (capnp.PreparedSend, error) {
	select {
	case <-q.callSendDone:
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-q.c.bgctx.Done():
		return nil, ExcClosed
	}
	err := withLockedConn1(q.c, func(c *lockedConn) error {
		if q.callSend != questionCallSendSucceeded {
			return q.callSendErr
		}
		q.mark(transform)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return q.c.prepareCall(ctx, s, nil, func(target rpccp.MessageTarget) error {
		pa, err := target.NewPromisedAnswer()
		if err != nil {
			return rpcerr.WrapFailed("build call message", err)
		}
		pa.SetQuestionId(uint32(q.id))
		oplist, err := pa.NewTransform(int32(len(transform)))
		if err != nil {
			return rpcerr.WrapFailed("build call message", err)
		}
		for i, op := range transform {
			oplist.At(i).SetGetPointerField(op.Field)
		}
		return nil
	})
}

// startCall starts an outbound question.  populateTarget selects either an
// imported capability or a promised answer without duplicating the question
// lifecycle around those two target encodings.
func (c *lockedConn) startCall(ctx context.Context, s capnp.Send, preflight func() error, populateTarget func(rpccp.MessageTarget) error) (*capnp.Answer, capnp.ReleaseFunc) {
	if !c.startTask() {
		return capnp.ErrorAnswer(s.Method, ExcClosed), func() {}
	}
	defer c.tasks.Done()
	if preflight != nil {
		if err := preflight(); err != nil {
			return capnp.ErrorAnswer(s.Method, err), func() {}
		}
	}

	q := c.newQuestion(s.Method)
	c.sendMessageOutcome(ctx, func(m rpccp.Message) error {
		return c.newCallMessage(m, q.id, s, populateTarget)
	}, false, func(outcome sendOutcome) { q.handleCallSend(ctx, outcome, "send message") })

	ans := q.p.Answer()
	return ans, func() {
		<-ans.Done()
		q.p.ReleaseClients()
		q.release()
	}
}

// preparedCall is the RPC implementation of capnp.PreparedSend.  The Call and
// question exist before it is returned, but the transport message is not made
// visible to the sender until Commit succeeds.
type preparedCall struct {
	c         *Conn
	q         *question
	msg       *preparedMessage
	ctx       context.Context
	preflight func(*lockedConn) error
	payload   rpccp.Payload

	mu        sync.Mutex
	committed bool
	aborted   bool
}

func (c *Conn) prepareCall(ctx context.Context, s capnp.Send, preflight func(*lockedConn) error, populateTarget func(rpccp.MessageTarget) error) (*preparedCall, error) {
	var q *question
	var err error
	c.withLocked(func(lc *lockedConn) {
		if !lc.startTask() {
			err = ExcClosed
			return
		}
		defer lc.tasks.Done()
		if preflight != nil {
			err = preflight(lc)
			if err != nil {
				return
			}
		}
		q = lc.newQuestion(s.Method)
	})
	if err != nil {
		return nil, err
	}
	p := &preparedCall{c: c, q: q, ctx: ctx, preflight: preflight}
	p.msg = (*lockedConn)(c).prepareMessage(func(m rpccp.Message) error {
		payload, err := newCallMessageUnsealed(m, q.id, s, populateTarget)
		p.payload = payload
		return err
	})
	if p.msg.preErr != nil {
		p.Abort()
		return nil, p.msg.preErr
	}
	return p, nil
}

func (p *preparedCall) Size() uint64 { return p.msg.size }

func (p *preparedCall) Commit(terminal func(flowcontrol.MessageOutcomeKind, error)) (*capnp.Answer, capnp.ReleaseFunc, error) {
	p.mu.Lock()
	if p.committed || p.aborted {
		p.mu.Unlock()
		return nil, nil, errors.New("prepared RPC send already completed")
	}
	p.committed = true
	p.mu.Unlock()

	var committed bool
	var commitErr error
	p.c.withLocked(func(c *lockedConn) {
		if p.ctx.Err() != nil {
			commitErr = p.ctx.Err()
			return
		}
		if c.bgctx.Err() != nil {
			commitErr = ExcClosed
			return
		}
		if p.preflight != nil {
			if err := p.preflight(c); err != nil {
				commitErr = err
				return
			}
		}
		if _, err := c.fillPayloadCapTable(p.payload); err != nil {
			commitErr = rpcerr.Annotate(err, "build call message")
			return
		}
		committed = p.msg.commit(p.ctx, false, func(outcome sendOutcome) {
			p.q.handleCallSend(p.ctx, outcome, "send message")
			switch outcome.disposition {
			case sendSucceeded:
				go func() { <-p.q.p.Answer().Done(); terminal(flowcontrol.MessageOutcomeSucceeded, nil) }()
			case sendDefinitelyUnsent:
				terminal(flowcontrol.MessageOutcomeAbortedBeforeEnqueue, outcome.err)
			default:
				terminal(flowcontrol.MessageOutcomeFatal, outcome.err)
			}
		})
	})
	if !committed {
		p.msg.abort()
		p.removeQuestion()
		return nil, nil, commitErr
	}
	ans := p.q.p.Answer()
	return ans, func() {
		<-ans.Done()
		p.q.p.ReleaseClients()
		p.q.release()
	}, nil
}

func (p *preparedCall) removeQuestion() {
	p.c.withLocked(func(c *lockedConn) {
		if q, ok := c.lk.questions.Find(p.q.id); ok && q == p.q {
			c.lk.questions.Remove(p.q.id)
		}
	})
}

func (p *preparedCall) Abort() {
	p.mu.Lock()
	if p.committed || p.aborted {
		p.mu.Unlock()
		return
	}
	p.aborted = true
	p.mu.Unlock()
	p.msg.abort()
	p.removeQuestion()
}

func (c *lockedConn) newCallMessage(msg rpccp.Message, qid questionID, s capnp.Send, populateTarget func(rpccp.MessageTarget) error) error {
	payload, err := newCallMessageUnsealed(msg, qid, s, populateTarget)
	if err != nil {
		return err
	}
	if s.PlaceArgs == nil {
		return nil
	}
	_, err = c.fillPayloadCapTable(payload)
	if err != nil {
		return rpcerr.Annotate(err, "build call message")
	}
	return nil
}

func newCallMessageUnsealed(msg rpccp.Message, qid questionID, s capnp.Send, populateTarget func(rpccp.MessageTarget) error) (rpccp.Payload, error) {
	call, err := msg.NewCall()
	if err != nil {
		return rpccp.Payload{}, rpcerr.WrapFailed("build call message", err)
	}
	call.SetQuestionId(uint32(qid))
	call.SetInterfaceId(s.Method.InterfaceID)
	call.SetMethodId(s.Method.MethodID)

	target, err := call.NewTarget()
	if err != nil {
		return rpccp.Payload{}, rpcerr.WrapFailed("build call message", err)
	}
	if err := populateTarget(target); err != nil {
		return rpccp.Payload{}, err
	}

	payload, err := call.NewParams()
	if err != nil {
		return rpccp.Payload{}, rpcerr.WrapFailed("build call message", err)
	}
	args, err := capnp.NewStruct(payload.Segment(), s.ArgsSize)
	if err != nil {
		return rpccp.Payload{}, rpcerr.WrapFailed("build call message", err)
	}
	if err := payload.SetContent(args.ToPtr()); err != nil {
		return rpccp.Payload{}, rpcerr.WrapFailed("build call message", err)
	}

	if s.PlaceArgs == nil {
		return payload, nil
	}
	if err := s.PlaceArgs(args); err != nil {
		return rpccp.Payload{}, rpcerr.WrapFailed("place arguments", err)
	}
	return payload, nil
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
// questions sent.  The caller must be holding onto q.c.lk.
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

func (q *question) Reject(err error) {
	if q != nil && q.p != nil {
		q.p.Reject(err)
	}
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
