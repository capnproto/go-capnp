package rpc

import (
	"context"

	"zombiezen.com/go/capnproto2"
	"zombiezen.com/go/capnproto2/internal/errors"
	rpccp "zombiezen.com/go/capnproto2/std/capnp/rpc"
)

// An answerID is an index into the answers table.
type answerID uint32

// answer is an entry in a Conn's answer table.
type answer struct {
	id      answerID
	ret     rpccp.Return // might fail in newAnswer
	results rpccp.Payload
	s       sendSession // might fail in newAnswer

	// All fields below are protected by s.c.mu.

	// state is a bitmask of which events have occurred in answer's
	// lifetime: 1 for return sent, 2 for finish received.
	// state is protected by conn.mu.
	state uint8

	cancel context.CancelFunc
}

// newAnswer adds an entry to the answers table and creates a new return
// message.  newAnswer may return both an answer and an error.
// The caller must be holding onto c.mu.
func (c *Conn) newAnswer(ctx context.Context, id answerID, cancel context.CancelFunc) (*answer, error) {
	if c.answers == nil {
		c.answers = make(map[answerID]*answer)
	} else if c.answers[id] != nil {
		// TODO(soon): abort
		return nil, errorf("answer ID %d reused", id)
	}
	ans := &answer{
		id:     id,
		cancel: cancel,
	}
	c.answers[id] = ans
	var err error
	ans.s, err = c.startSend(ctx)
	if err != nil {
		ans.s = sendSession{}
		return ans, err
	}
	ans.ret, err = ans.s.msg.NewReturn()
	if err != nil {
		ans.s.finish()
		ans.s = sendSession{}
		return ans, errorf("create return: %v", err)
	}
	ans.ret.SetAnswerId(uint32(id))
	ans.ret.SetReleaseParamCaps(false)
	ans.s.releaseSender()
	return ans, nil
}

// AllocResults allocates the results struct.
func (ans *answer) AllocResults(sz capnp.ObjectSize) (capnp.Struct, error) {
	var err error
	ans.results, err = ans.ret.NewResults()
	if err != nil {
		return capnp.Struct{}, errorf("alloc results: %v", err)
	}
	s, err := capnp.NewStruct(ans.results.Segment(), sz)
	if err != nil {
		return capnp.Struct{}, errorf("alloc results: %v", err)
	}
	if err := ans.results.SetContent(s.ToPtr()); err != nil {
		return capnp.Struct{}, errorf("alloc results: %v", err)
	}
	return s, nil
}

// setBootstrap sets the results to an interface pointer, stealing the
// reference.
func (ans *answer) setBootstrap(c *capnp.Client) error {
	var err error
	ans.results, err = ans.ret.NewResults()
	if err != nil {
		return errorf("alloc bootstrap results: %v", err)
	}
	capID := ans.results.Message().AddCap(c)
	iface := capnp.NewInterface(ans.results.Segment(), capID)
	if err := ans.results.SetContent(iface.ToPtr()); err != nil {
		return errorf("alloc bootstrap results: %v", err)
	}
	return nil
}

// Return sends the return message.  The caller must NOT be holding onto
// ans.s.c.mu or the sender lock.
func (ans *answer) Return(e error) {
	defer ans.s.c.mu.Unlock()
	ans.s.c.mu.Lock()
	ans.lockedReturn(e)
}

// lockedReturn sends the return message.  The caller must be holding
// onto ans.s.c.mu.
func (ans *answer) lockedReturn(e error) {
	if e == nil {
		if err := ans.s.c.fillPayloadCapTable(ans.results); err != nil {
			ans.s.c.report(annotate(err).errorf("send return"))
			// Continue.  Don't fail to send return if cap table isn't fully filled.
		}
	} else {
		exc, err := ans.ret.NewException()
		if err != nil {
			ans.s.acquireSender()
			ans.s.finish()
			ans.s.c.reportf("send exception: %v", err)
			return
		}
		exc.SetType(rpccp.Exception_Type(errors.TypeOf(e)))
		if err := exc.SetReason(e.Error()); err != nil {
			ans.s.acquireSender()
			ans.s.finish()
			ans.s.c.reportf("send exception: %v", err)
			return
		}
	}
	ans.s.acquireSender()
	err := ans.s.send()
	ans.s.finish()
	if err != nil {
		ans.s.c.reportf("send return: %v", err)
	}

	ans.state |= 1
	if !ans.isDone() {
		return
	}
	delete(ans.s.c.answers, ans.id)
	// TODO(soon): release result caps (while not holding c.mu)
}

// isDone reports whether the answer should be removed from the answers
// table.
func (ans *answer) isDone() bool {
	return ans.state&3 == 3
}
