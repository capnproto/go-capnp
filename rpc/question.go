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
	id        questionID
	bootstrap bool
	conn      *Conn

	p       *capnp.Promise
	release capnp.ReleaseFunc  // written before resolving p
	done    context.CancelFunc // called after resolving p

	// state is a bitmask of which events have occurred in question's
	// lifetime: 1 for return received, 2 for sending finish.
	// state is protected by conn.mu.
	state uint8
}

// newQuestion adds a new question to c's table.  The caller must be
// holding onto c.mu.
func (c *Conn) newQuestion(ctx context.Context, id questionID, method capnp.Method, bootstrap bool) *question {
	ctx, cancel := context.WithCancel(ctx)
	q := &question{
		id:        id,
		conn:      c,
		done:      cancel,
		bootstrap: bootstrap,
	}
	q.p = capnp.NewPromise(method, q)
	if int(id) == len(c.questions) {
		c.questions = append(c.questions, q)
	} else {
		c.questions[id] = q
	}
	c.runBackground(func(bgctx context.Context) {
		var rejectErr error
		select {
		case <-ctx.Done():
			rejectErr = ctx.Err()
		case <-bgctx.Done():
			rejectErr = bgctx.Err()
			q.done()
		}
		c.mu.Lock()
		if q.sentFinish() {
			c.mu.Unlock()
			return
		}
		q.state |= 2 // sending finish
		q.release = func() {}
		select {
		case <-bgctx.Done():
		default:
			err := c.sendMessage(bgctx, func(msg rpccp.Message) error {
				fin, err := msg.NewFinish()
				if err != nil {
					return err
				}
				fin.SetQuestionId(uint32(q.id))
				fin.SetReleaseResultCaps(true)
				return nil
			})
			if err != nil {
				c.report(annotate(err).errorf("send finish"))
			}
		}
		c.mu.Unlock()
		q.p.Reject(rejectErr)
		if q.bootstrap {
			q.p.ReleaseClients()
		}
	})
	return q
}

func (q *question) PipelineSend(ctx context.Context, transform []capnp.PipelineOp, s capnp.Send) (*capnp.Answer, capnp.ReleaseFunc) {
	defer q.conn.mu.Unlock()
	q.conn.mu.Lock()
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
		return capnp.ErrorAnswer(s.Method, annotate(err).errorf("send to promised answer")), func() {}
	}
	q2 := q.conn.newQuestion(ctx, id, s.Method, false)
	ans := q2.p.Answer()
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

func (q *question) sentFinish() bool {
	return q.state&2 != 0
}
