package rpc

import (
	"context"
	"sync"

	"zombiezen.com/go/capnproto2"
	"zombiezen.com/go/capnproto2/internal/errors"
	rpccp "zombiezen.com/go/capnproto2/std/capnp/rpc"
)

// An answerID is an index into the answers table.
type answerID uint32

// answer is an entry in a Conn's answer table.
type answer struct {
	// Fields set by newAnswer:

	id  answerID
	ret rpccp.Return // not set if newAnswer fails
	s   sendSession  // not set if newAnswer fails

	// All fields below are protected by s.c.mu.

	flags answerFlags

	// results is the memoized answer to ret.Results().
	// Set by AllocResults and setBootstrap, but contents cannot be
	// used by the rpc package until flags has resultsReady set.
	results rpccp.Payload

	// exportRefs is the number of references to exports placed in the
	// results.
	exportRefs map[exportID]uint32

	// cancel is the cancel function for the Context used in the received
	// method call.
	cancel context.CancelFunc

	// pcall is the PipelineCaller returned by RecvCall.  It will be set
	// to nil once results are ready.
	pcall capnp.PipelineCaller

	// pcalls is added to for every pending RecvCall and subtracted from
	// for every RecvCall return (delivery acknowledgement).  This is used
	// to satisfy the Returner.Return contract.
	pcalls sync.WaitGroup

	// err is the error passed to (*answer).sendException.
	err error
}

// answerFlags is a bitmask of events that have occurred in an answer's
// lifetime.
type answerFlags uint8

const (
	returnSent answerFlags = 1 << iota
	finishReceived
	resultsReady
	releaseResultCapsFlag
)

// newAnswer adds an entry to the answers table and creates a new return
// message.  newAnswer may return both an answer and an error.  Results
// should not be set on the answer if newAnswer returns a non-nil error.
// The caller must be holding onto c.mu.
func (c *Conn) newAnswer(ctx context.Context, id answerID, cancel context.CancelFunc) (*answer, error) {
	if c.answers == nil {
		c.answers = make(map[answerID]*answer)
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

// setPipelineCaller sets ans.pcall to pcall if the answer has not
// already returned.  The caller MUST NOT be holding onto ans.s.c.mu.
func (ans *answer) setPipelineCaller(pcall capnp.PipelineCaller) {
	ans.s.c.mu.Lock()
	if ans.flags&resultsReady == 0 {
		ans.pcall = pcall
	}
	ans.s.c.mu.Unlock()
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
	// Add the capability to the table early to avoid leaks if setBootstrap fails.
	capID := ans.ret.Message().AddCap(c)

	var err error
	ans.results, err = ans.ret.NewResults()
	if err != nil {
		return errorf("alloc bootstrap results: %v", err)
	}
	iface := capnp.NewInterface(ans.results.Segment(), capID)
	if err := ans.results.SetContent(iface.ToPtr()); err != nil {
		return errorf("alloc bootstrap results: %v", err)
	}
	return nil
}

// Return sends the return message.
//
// The caller must NOT be holding onto ans.s.c.mu or the sender lock.
func (ans *answer) Return(e error) {
	c := ans.s.c
	c.mu.Lock()
	if e != nil {
		ans.sendException(e)
		c.mu.Unlock()
		ans.pcalls.Wait()
		c.tasks.Done() // added by handleCall
		return
	}
	refs := ans.sendReturn()
	rl, err := c.releaseExports(refs)
	if err != nil {
		select {
		case <-ans.s.c.bgctx.Done():
		default:
			c.tasks.Done() // added by handleCall
			if err := c.shutdown(err); err != nil {
				c.report(err)
			}
			// shutdown released c.mu
			rl.release()
			ans.pcalls.Wait()
			return
		}
	}
	c.mu.Unlock()
	rl.release()
	ans.pcalls.Wait()
	c.tasks.Done() // added by handleCall
}

// sendReturn sends the return message with results allocated by a
// previous call to AllocResults.  If the answer already received a
// Finish with releaseResultCaps set to true, then sendReturn returns
// the number of references to be subtracted from each export.
//
// The caller must be holding onto ans.s.c.mu.  Only one of sendReturn
// or sendException should be called.
func (ans *answer) sendReturn() map[exportID]uint32 {
	ans.pcall = nil
	ans.flags |= resultsReady

	refs, err := ans.s.c.fillPayloadCapTable(ans.results)
	ans.exportRefs = refs
	if err != nil {
		ans.s.c.report(annotate(err).errorf("send return"))
		// Continue.  Don't fail to send return if cap table isn't fully filled.
	}

	// Send results.
	fin := ans.flags&finishReceived != 0
	ans.s.acquireSender()
	if err := ans.s.send(); err != nil {
		ans.s.c.reportf("send return: %v", err)
	}
	if !fin {
		ans.s.releaseSender()
		ans.flags |= returnSent
		if ans.flags&finishReceived == 0 {
			return nil
		}
		ans.s.acquireSender()
	}

	// Already received finish, delete answer.
	ans.s.finish()
	delete(ans.s.c.answers, ans.id)
	if ans.flags&releaseResultCapsFlag != 0 {
		return ans.exportRefs
	}
	return nil
}

// sendException sends an exception on the answer's return message.
//
// The caller must be holding onto ans.s.c.mu.  Only one of sendReturn
// or sendException should be called.
func (ans *answer) sendException(e error) {
	ans.err = e
	ans.pcall = nil
	ans.flags |= resultsReady

	exc, err := ans.ret.NewException()
	if err != nil {
		ans.s.acquireSender()
		ans.s.finish() // TODO(now): is this okay?
		ans.s.c.reportf("send exception: %v", err)
		return
	}
	exc.SetType(rpccp.Exception_Type(errors.TypeOf(e)))
	if err := exc.SetReason(e.Error()); err != nil {
		ans.s.acquireSender()
		ans.s.finish() // TODO(now): is this okay?
		ans.s.c.reportf("send exception: %v", err)
		return
	}

	// Send results.
	fin := ans.flags&finishReceived != 0
	ans.s.acquireSender()
	if err := ans.s.send(); err != nil {
		ans.s.c.reportf("send return: %v", err)
	}
	if !fin {
		ans.s.releaseSender()
		ans.flags |= returnSent
		if ans.flags&finishReceived == 0 {
			return
		}
		ans.s.acquireSender()
	}

	// Already received finish, delete answer.
	ans.s.finish()
	delete(ans.s.c.answers, ans.id)
}
