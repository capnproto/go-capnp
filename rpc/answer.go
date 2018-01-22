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
	// c and id must be set before any answer methods are called.
	c  *Conn
	id answerID

	// cancel cancels the Context used in the received method call.
	// May be nil.
	cancel context.CancelFunc

	// ret is the outgoing Return struct.  ret is valid iff there was no
	// error creating the message.  If ret is invalid, then this answer
	// entry is a placeholder until the remote vat cancels the call.
	ret rpccp.Return

	// sendMsg sends the return message.  The caller must be holding onto
	// the sender lock but not ans.c.mu.
	sendMsg func() error

	// releaseMsg releases the return message.  The caller must be holding
	// onto the sender lock but not ans.c.mu.
	releaseMsg capnp.ReleaseFunc

	// results is the memoized answer to ret.Results().
	// Set by AllocResults and setBootstrap, but contents can only be read
	// if flags has resultsReady but not finishReceived set.
	results rpccp.Payload

	// All fields below are protected by s.c.mu.

	// flags is a bitmask of events that have occurred in an answer's
	// lifetime.
	flags answerFlags

	// resultCapTable is the CapTable for results.  It is not kept in the
	// results message because CapTable cannot be used once results are
	// sent.  However, the capabilities need to be retained for promised
	// answer targets.
	resultCapTable []*capnp.Client

	// exportRefs is the number of references to exports placed in the
	// results.
	exportRefs map[exportID]uint32

	// pcall is the PipelineCaller returned by RecvCall.  It will be set
	// to nil once results are ready.
	pcall capnp.PipelineCaller

	// pcalls is added to for every pending RecvCall and subtracted from
	// for every RecvCall return (delivery acknowledgement).  This is used
	// to satisfy the Returner.Return contract.
	pcalls sync.WaitGroup

	// err is the error passed to (*answer).sendException or from creating
	// the Return message.  Can only be read after resultsReady is set in
	// flags.
	err error
}

type answerFlags uint8

const (
	returnSent answerFlags = 1 << iota
	finishReceived
	resultsReady
	releaseResultCapsFlag
)

// errorAnswer returns a placeholder answer with an error result already set.
func errorAnswer(c *Conn, id answerID, err error) *answer {
	return &answer{
		c:     c,
		id:    id,
		err:   err,
		flags: resultsReady | returnSent,
	}
}

// newReturn creates a new Return message.  The caller must be holding
// onto the sender lock but not c.mu.
func (c *Conn) newReturn(ctx context.Context) (rpccp.Return, func() error, capnp.ReleaseFunc, error) {
	msg, send, release, err := c.transport.NewMessage(ctx)
	if err != nil {
		return rpccp.Return{}, nil, nil, errorf("create return: %v", err)
	}
	ret, err := msg.NewReturn()
	if err != nil {
		release()
		return rpccp.Return{}, nil, nil, errorf("create return: %v", err)
	}
	return ret, send, release, nil
}

// setPipelineCaller sets ans.pcall to pcall if the answer has not
// already returned.  The caller MUST NOT be holding onto ans.c.mu
// or the sender lock.
func (ans *answer) setPipelineCaller(pcall capnp.PipelineCaller) {
	ans.c.mu.Lock()
	if ans.flags&resultsReady == 0 {
		ans.pcall = pcall
	}
	ans.c.mu.Unlock()
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
	if ans.ret.HasResults() || len(ans.ret.Message().CapTable) > 0 || len(ans.resultCapTable) > 0 {
		panic("setBootstrap called after creating results")
	}
	// Add the capability to the table early to avoid leaks if setBootstrap fails.
	ans.resultCapTable = []*capnp.Client{c}

	var err error
	ans.results, err = ans.ret.NewResults()
	if err != nil {
		return errorf("alloc bootstrap results: %v", err)
	}
	iface := capnp.NewInterface(ans.results.Segment(), 0)
	if err := ans.results.SetContent(iface.ToPtr()); err != nil {
		return errorf("alloc bootstrap results: %v", err)
	}
	return nil
}

// Return sends the return message.
//
// The caller must NOT be holding onto ans.c.mu or the sender lock.
func (ans *answer) Return(e error) {
	var cstates []capnp.ClientState
	if ans.results.IsValid() {
		ans.resultCapTable, cstates = extractCapTable(ans.results.Message())
	}
	ans.c.mu.Lock()
	ans.c.lockSender()
	if e != nil {
		rl := ans.sendException(e)
		ans.c.unlockSender()
		ans.c.mu.Unlock()
		rl.release()
		ans.pcalls.Wait()
		ans.c.tasks.Done() // added by handleCall
		return
	}
	rl, err := ans.sendReturn(cstates)
	ans.c.unlockSender()
	if err != nil {
		select {
		case <-ans.c.bgctx.Done():
		default:
			ans.c.tasks.Done() // added by handleCall
			if err := ans.c.shutdown(err); err != nil {
				ans.c.report(err)
			}
			// shutdown released c.mu
			rl.release()
			ans.pcalls.Wait()
			return
		}
	}
	ans.c.mu.Unlock()
	rl.release()
	ans.pcalls.Wait()
	ans.c.tasks.Done() // added by handleCall
}

// sendReturn sends the return message with results allocated by a
// previous call to AllocResults.  If the answer already received a
// Finish with releaseResultCaps set to true, then sendReturn returns
// the number of references to be subtracted from each export.
//
// The caller must be holding onto ans.c.mu and the sender lock.
// The result's capability table must have been extracted into
// ans.resultsCapTable before calling sendReturn. Only one of
// sendReturn or sendException should be called.
func (ans *answer) sendReturn(cstates []capnp.ClientState) (releaseList, error) {
	ans.pcall = nil
	ans.flags |= resultsReady
	var err error
	ans.exportRefs, err = ans.c.fillPayloadCapTable(ans.results, ans.resultCapTable, cstates)
	if err != nil {
		ans.c.report(annotate(err).errorf("send return"))
		// Continue.  Don't fail to send return if cap table isn't fully filled.
	}

	select {
	case <-ans.c.bgctx.Done():
	default:
		fin := ans.flags&finishReceived != 0
		ans.c.mu.Unlock()
		if err := ans.sendMsg(); err != nil {
			ans.c.reportf("send return: %v", err)
		}
		if fin {
			ans.releaseMsg()
			ans.c.mu.Lock()
			return ans.destroy()
		}
		ans.c.mu.Lock()
	}
	ans.flags |= returnSent
	if ans.flags&finishReceived == 0 {
		return nil, nil
	}
	rl, err := ans.destroy()
	ans.c.mu.Unlock()
	ans.releaseMsg()
	ans.c.mu.Lock()
	return rl, err
}

// sendException sends an exception on the answer's return message.
//
// The caller must be holding onto ans.c.mu and the sender lock.
// The result's capability table must have been extracted into
// ans.resultsCapTable before calling sendException. Only one of
// sendReturn or sendException should be called.
func (ans *answer) sendException(e error) releaseList {
	ans.err = e
	ans.pcall = nil
	ans.flags |= resultsReady

	select {
	case <-ans.c.bgctx.Done():
	default:
		// Send exception.
		fin := ans.flags&finishReceived != 0
		ans.c.mu.Unlock()
		if exc, err := ans.ret.NewException(); err != nil {
			ans.c.reportf("send exception: %v", err)
		} else {
			exc.SetType(rpccp.Exception_Type(errors.TypeOf(e)))
			if err := exc.SetReason(e.Error()); err != nil {
				ans.c.reportf("send exception: %v", err)
			} else if err := ans.sendMsg(); err != nil {
				ans.c.reportf("send return: %v", err)
			}
		}
		if fin {
			ans.releaseMsg()
			ans.c.mu.Lock()
			rl, _ := ans.destroy()
			return rl
		}
		ans.c.mu.Lock()
	}
	ans.flags |= returnSent
	if ans.flags&finishReceived == 0 {
		return nil
	}
	// destory will never return an error because sendException does
	// create any exports.
	rl, _ := ans.destroy()
	ans.c.mu.Unlock()
	ans.releaseMsg()
	ans.c.mu.Lock()
	return rl
}

// destroy removes the answer from the table and returns the clients to
// release.  The answer must have sent a return and received a finish.
// The caller must be holding onto ans.c.mu.
//
// shutdown has its own strategy for cleaning up an answer.
func (ans *answer) destroy() (releaseList, error) {
	delete(ans.c.answers, ans.id)
	rl := releaseList(ans.resultCapTable)
	if ans.flags&releaseResultCapsFlag == 0 || len(ans.exportRefs) == 0 {
		return rl, nil
	}
	exportReleases, err := ans.c.releaseExports(ans.exportRefs)
	return append(rl, exportReleases...), err
}
