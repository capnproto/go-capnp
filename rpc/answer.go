package rpc

import (
	"context"
	"sync"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/internal/errors"
	"capnproto.org/go/capnp/v3/internal/syncutil"
	rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"
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

	// sendMsg sends the return message.  The caller MUST NOT hold ans.c.mu.
	sendMsg func() error

	// releaseMsg releases the return message.  The caller MUST NOT hold
	// ans.c.mu.
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
	// promise is a promise wrapping pcall. It will be resolved, and
	// subsequently set to nil, once the results are ready.
	promise *capnp.Promise

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

// newReturn creates a new Return message.
func (c *Conn) newReturn(ctx context.Context) (rpccp.Return, func() error, capnp.ReleaseFunc, error) {
	msg, send, release, err := c.transport.NewMessage(ctx)
	if err != nil {
		return rpccp.Return{}, nil, nil, failedf("create return: %w", err)
	}
	ret, err := msg.NewReturn()
	if err != nil {
		release()
		return rpccp.Return{}, nil, nil, failedf("create return: %w", err)
	}
	return ret, send, release, nil
}

// setPipelineCaller sets ans.pcall to pcall if the answer has not
// already returned.  The caller MUST NOT hold ans.c.mu.
//
// This also sets ans.promise to a new promise, wrapping pcall.
func (ans *answer) setPipelineCaller(m capnp.Method, pcall capnp.PipelineCaller) {
	syncutil.With(&ans.c.mu, func() {
		if ans.flags&resultsReady == 0 {
			ans.pcall = pcall
			ans.promise = capnp.NewPromise(m, pcall)
		}
	})
}

// AllocResults allocates the results struct.
func (ans *answer) AllocResults(sz capnp.ObjectSize) (capnp.Struct, error) {
	var err error
	ans.results, err = ans.ret.NewResults()
	if err != nil {
		return capnp.Struct{}, failedf("alloc results: %w", err)
	}
	s, err := capnp.NewStruct(ans.results.Segment(), sz)
	if err != nil {
		return capnp.Struct{}, failedf("alloc results: %w", err)
	}
	if err := ans.results.SetContent(s.ToPtr()); err != nil {
		return capnp.Struct{}, failedf("alloc results: %w", err)
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
		return failedf("alloc bootstrap results: %w", err)
	}
	iface := capnp.NewInterface(ans.results.Segment(), 0)
	if err := ans.results.SetContent(iface.ToPtr()); err != nil {
		return failedf("alloc bootstrap results: %w", err)
	}
	return nil
}

// Return sends the return message.
//
// The caller MUST NOT hold ans.c.mu.
func (ans *answer) Return(e error) {
	if ans.results.IsValid() {
		ans.resultCapTable = ans.results.Message().CapTable
	}
	ans.c.mu.Lock()
	if e != nil {
		rl := ans.sendException(e)
		ans.c.mu.Unlock()
		rl.release()
		ans.pcalls.Wait()
		ans.c.tasks.Done() // added by handleCall
		return
	}
	rl, err := ans.sendReturn()
	if err != nil {
		select {
		case <-ans.c.bgctx.Done():
		default:
			ans.c.tasks.Done() // added by handleCall
			if err := ans.c.shutdown(err); err != nil {
				ans.c.er.ReportError(err)
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
// The caller MUST hold ans.c.mu.  The result's capability table MUST
// have been extracted into ans.resultCapTable before calling sendReturn,
// and call only one of sendReturn or sendException.
func (ans *answer) sendReturn() (releaseList, error) {
	ans.pcall = nil
	ans.flags |= resultsReady

	var err error
	ans.exportRefs, err = ans.c.fillPayloadCapTable(ans.results, ans.resultCapTable)
	if err != nil {
		ans.c.er.annotatef(err, "send return")
		// Continue.  Don't fail to send return if cap table isn't fully filled.
	}

	select {
	case <-ans.c.bgctx.Done():
	default:
		fin := ans.flags&finishReceived != 0
		if ans.promise != nil {
			if fin {
				// Can't use ans.result after a finish, but it's
				// ok to return an error if the finish comes in
				// before the return. Possible enhancement: use
				// the cancel variant of return.
				ans.promise.Reject(failedf("Received finish before return"))
			} else {
				ans.promise.Resolve(ans.results.Content())
			}
			ans.promise = nil
		}
		ans.c.mu.Unlock()
		if err := ans.sendMsg(); err != nil {
			ans.c.er.reportf("send return: %w", err)
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
	syncutil.Without(&ans.c.mu, func() {
		ans.releaseMsg()
	})
	return rl, err
}

// sendException sends an exception on the answer's return message.
//
// The caller MUST hold ans.c.mu. The result's capability table MUST
// have been extracted into ans.resultCapTable before calling, and
// call only one of sendReturn or sendException.
func (ans *answer) sendException(e error) releaseList {
	ans.err = e
	ans.pcall = nil
	ans.flags |= resultsReady

	if ans.promise != nil {
		ans.promise.Reject(e)
		ans.promise = nil
	}

	select {
	case <-ans.c.bgctx.Done():
	default:
		// Send exception.
		fin := ans.flags&finishReceived != 0
		ans.c.mu.Unlock()
		if exc, err := ans.ret.NewException(); err != nil {
			ans.c.er.reportf("send exception: %w", err)
		} else {
			exc.SetType(rpccp.Exception_Type(errors.TypeOf(e)))
			if err := exc.SetReason(e.Error()); err != nil {
				ans.c.er.reportf("send exception: %w", err)
			} else if err := ans.sendMsg(); err != nil {
				ans.c.er.reportf("send return: %w", err)
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
	// destroy will never return an error because sendException does
	// create any exports.
	rl, _ := ans.destroy()
	syncutil.Without(&ans.c.mu, func() {
		ans.releaseMsg()
	})
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
