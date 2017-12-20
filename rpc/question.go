package rpc

import (
	"context"
	"fmt"

	"zombiezen.com/go/capnproto2"
	rpccp "zombiezen.com/go/capnproto2/std/capnp/rpc"
)

// A questionID is an index into the questions table.
type questionID uint32

type question struct {
	id   questionID
	conn *Conn

	bootstrapPromise *capnp.ClientPromise

	p       *capnp.Promise
	release capnp.ReleaseFunc  // written before resolving p
	done    context.CancelFunc // called after resolving p

	// Protected by conn.mu:

	flags         questionFlags
	finishMsgSend chan struct{} // closed after attempting to send the Finish message
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
func (c *Conn) newQuestion(ctx context.Context, id questionID, method capnp.Method) *question {
	ctx, cancel := context.WithCancel(ctx)
	q := &question{
		id:            id,
		conn:          c,
		done:          cancel,
		finishMsgSend: make(chan struct{}),
	}
	q.p = capnp.NewPromise(method, q) // TODO(someday): customize error message for bootstrap
	if int(id) == len(c.questions) {
		c.questions = append(c.questions, q)
	} else {
		c.questions[id] = q
	}
	c.tasks.Add(1)
	go func() {
		defer c.tasks.Done()
		var rejectErr error
		select {
		case <-ctx.Done():
			rejectErr = ctx.Err()
		case <-c.bgctx.Done():
			rejectErr = disconnected("connection closed")
			q.done()
		}
		c.mu.Lock()
		if q.flags&finished != 0 {
			c.mu.Unlock()
			return
		}
		q.flags |= finished
		q.release = func() {}
		err := c.sendMessage(c.bgctx, func(msg rpccp.Message) error {
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
			case <-c.bgctx.Done():
			default:
				c.report(annotate(err).errorf("send finish"))
			}
		}
		close(q.finishMsgSend)
		c.mu.Unlock()
		q.p.Reject(rejectErr)
		if q.bootstrapPromise != nil {
			q.bootstrapPromise.Fulfill(q.p.Answer().Client())
			q.p.ReleaseClients()
		}
	}()
	return q
}

func (q *question) PipelineSend(ctx context.Context, transform []capnp.PipelineOp, s capnp.Send) (*capnp.Answer, capnp.ReleaseFunc) {
	q.conn.mu.Lock()
	if !q.conn.startTask() {
		q.conn.mu.Unlock()
		return capnp.ErrorAnswer(s.Method, disconnected("connection closed")), func() {}
	}
	defer q.conn.tasks.Done()
	id := questionID(q.conn.questionID.next())
	err := q.conn.sendMessage(ctx, func(msg rpccp.Message) error {
		call, err := msg.NewCall()
		if err != nil {
			return err
		}
		call.SetQuestionId(uint32(id))
		call.SetInterfaceId(s.Method.InterfaceID)
		call.SetMethodId(s.Method.MethodID)
		tgt, err := call.NewTarget()
		if err != nil {
			return err
		}
		pa, err := tgt.NewPromisedAnswer()
		if err != nil {
			return err
		}
		pa.SetQuestionId(uint32(q.id))
		xform, err := pa.NewTransform(int32(len(transform)))
		if err != nil {
			return err
		}
		for i, t := range transform {
			xform.At(i).SetGetPointerField(t.Field)
		}
		params, err := call.NewParams()
		if err != nil {
			return err
		}
		args, err := capnp.NewStruct(params.Segment(), s.ArgsSize)
		if err != nil {
			return err
		}
		if err := params.SetContent(args.ToPtr()); err != nil {
			return err
		}
		if err := s.PlaceArgs(args); err != nil {
			// Using fmt.Errorf to annotate to avoid stutter when we wrap the
			// sendMessage error.
			return fmt.Errorf("place args: %v", err)
		}
		// TODO(soon): fill in capability table
		return nil
	})
	if err != nil {
		q.conn.questionID.remove(uint32(id))
		q.conn.mu.Unlock()
		return capnp.ErrorAnswer(s.Method, annotate(err).errorf("send to promised answer")), func() {}
	}
	q2 := q.conn.newQuestion(ctx, id, s.Method)
	ans := q2.p.Answer()
	q.conn.mu.Unlock()
	return ans, func() {
		<-ans.Done()
		q2.p.ReleaseClients()
		q2.release()
	}
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
