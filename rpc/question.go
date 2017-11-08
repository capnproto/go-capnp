package rpc

import (
	"context"
	"errors"
	"fmt"

	"zombiezen.com/go/capnproto2"
	rpccp "zombiezen.com/go/capnproto2/std/capnp/rpc"
)

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
			return fmt.Errorf("place args: %v", err)
		}
		// TODO(soon): fill in capability table
		return nil
	})
	if err != nil {
		q.conn.questionID.remove(uint32(id))
		return capnp.ErrorAnswer(fmt.Errorf("rpc: call promised answer: %v", err)), func() {}
	}
	q2 := q.conn.newQuestion(ctx, id, false)
	ans := q2.p.Answer()
	return ans, func() {
		<-ans.Done()
		q2.p.ReleaseClients()
		q2.release()
	}
}

func (q *question) PipelineRecv(ctx context.Context, transform []capnp.PipelineOp, r capnp.Recv) capnp.PipelineCaller {
	r.Reject(errors.New("TODO(soon)"))
	return nil
}

func (q *question) sentFinish() bool {
	return q.state&2 != 0
}
