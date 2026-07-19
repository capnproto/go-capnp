package rpc

import (
	"context"

	"capnproto.org/go/capnp/v3"
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

// flags.Contains(flag) Returns true iff flags contains flag, which must
// be a single flag.
func (flags questionFlags) Contains(flag questionFlags) bool {
	return flags&flag != 0
}

// newQuestion adds a new question to c's table.
func (c *lockedConn) newQuestion(method capnp.Method) *question {
	q := &question{
		c:             (*Conn)(c),
		release:       func() {},
		finishMsgSend: make(chan struct{}),
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

	q.c.withLocked(func(c *lockedConn) {
		// Promise already fulfilled?
		if q.flags.Contains(finished) {
			return
		}
		q.flags |= finished
		q.release = func() {}

		c.sendMessage(c.bgctx, func(m rpccp.Message) error {
			fin, err := m.NewFinish()
			if err != nil {
				return err
			}
			fin.SetQuestionId(uint32(q.id))
			fin.SetReleaseResultCaps(true)
			return nil
		}, func(err error) {
			if err == nil {
				syncutil.With(&q.c.lk, func() { q.flags |= finishSent })
			} else if q.c.bgctx.Err() == nil {
				q.c.er.ReportError(rpcerr.Annotate(err, "send finish"))
			}
			close(q.finishMsgSend)
			q.p.Reject(rejectErr)
		})
	})
}

func (q *question) PipelineSend(ctx context.Context, transform []capnp.PipelineOp, s capnp.Send) (*capnp.Answer, capnp.ReleaseFunc) {
	return withLockedConn2(q.c, func(c *lockedConn) (*capnp.Answer, capnp.ReleaseFunc) {
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
	c.sendMessage(ctx, func(m rpccp.Message) error {
		return c.newCallMessage(m, q.id, s, populateTarget)
	}, func(err error) {
		if err != nil {
			syncutil.With(&q.c.lk, func() { q.c.lk.questions.Remove(q.id) })
			q.p.Reject(rpcerr.WrapFailed("send message", err))
			return
		}

		q.c.tasks.Add(1)
		go func() {
			defer q.c.tasks.Done()
			q.handleCancel(ctx)
		}()
	})

	ans := q.p.Answer()
	return ans, func() {
		<-ans.Done()
		q.p.ReleaseClients()
		q.release()
	}
}

func (c *lockedConn) newCallMessage(msg rpccp.Message, qid questionID, s capnp.Send, populateTarget func(rpccp.MessageTarget) error) error {
	call, err := msg.NewCall()
	if err != nil {
		return rpcerr.WrapFailed("build call message", err)
	}
	call.SetQuestionId(uint32(qid))
	call.SetInterfaceId(s.Method.InterfaceID)
	call.SetMethodId(s.Method.MethodID)

	target, err := call.NewTarget()
	if err != nil {
		return rpcerr.WrapFailed("build call message", err)
	}
	if err := populateTarget(target); err != nil {
		return err
	}

	payload, err := call.NewParams()
	if err != nil {
		return rpcerr.WrapFailed("build call message", err)
	}
	args, err := capnp.NewStruct(payload.Segment(), s.ArgsSize)
	if err != nil {
		return rpcerr.WrapFailed("build call message", err)
	}
	if err := payload.SetContent(args.ToPtr()); err != nil {
		return rpcerr.WrapFailed("build call message", err)
	}

	if s.PlaceArgs == nil {
		return nil
	}
	if err := s.PlaceArgs(args); err != nil {
		return rpcerr.WrapFailed("place arguments", err)
	}
	// TODO(soon): save param refs
	_, err = c.fillPayloadCapTable(payload)

	if err != nil {
		return rpcerr.Annotate(err, "build call message")
	}

	return err
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
