// Package rpc implements the Cap'n Proto RPC protocol.
package rpc // import "zombiezen.com/go/capnproto2/rpc"

import (
	"context"
	"fmt"
	"sync"

	"zombiezen.com/go/capnproto2"
	"zombiezen.com/go/capnproto2/internal/errors"
	rpccp "zombiezen.com/go/capnproto2/std/capnp/rpc"
)

// A Conn is a connection to another Cap'n Proto vat.
// It is safe to use from multiple goroutines.
type Conn struct {
	bootstrap *capnp.Client
	reporter  ErrorReporter

	// bgctx is a Context that is canceled when shutdown starts.
	bgctx context.Context

	// bgcancel cancels bgctx.  It must only be called while holding mu.
	bgcancel context.CancelFunc

	// tasks block shutdown.
	tasks sync.WaitGroup

	// mu protects all the following fields in the Conn.  mu should not be
	// held while making calls that take indeterminate time (I/O or
	// application code).  Condition channels protect operations on any
	// field that take an indeterminate amount of time.  Thus, critical
	// sections involving mu are quite short, while still ensuring
	// mutually exclusive access to resources.
	mu sync.Mutex

	closed bool          // set when Close() is called, used to distinguish user Closing multiple times
	shut   chan struct{} // closed when shutdown() returns

	recvCloser interface {
		CloseRecv() error
	}

	// sendCond is non-nil if an operation involving sender is in
	// progress, and the channel is closed when the operation is finished.
	// This is referred to as the "sender lock".  See newMessage to start
	// an operation on sender.
	sendCond chan struct{}

	// sender is the send half of the transport.  sender is protected by
	// sendCond.  sender should not be used after bgctx is canceled (Close
	// is starting).
	sender Sender

	// Tables
	questions  []*question
	questionID idgen
	answers    map[answerID]*answer
	exports    []*expent
	exportID   idgen
	imports    map[importID]*impent
}

// Options specifies optional parameters for creating a Conn.
type Options struct {
	// BootstrapClient is the capability that will be returned to the
	// remote peer when receiving a Bootstrap message.  NewConn "steals"
	// this reference: it will release the client when the connection is
	// closed.
	BootstrapClient *capnp.Client

	// ErrorReporter will be called upon when errors occur while the Conn
	// is receiving messages from the remote vat.
	ErrorReporter ErrorReporter
}

// A type that implements ErrorReporter can receive errors from a Conn.
// ReportError should be quick to return and should not use the Conn
// that it is attached to.
type ErrorReporter interface {
	ReportError(error)
}

// NewConn creates a new connection that communications on a given
// transport.  Closing the connection will close the transport.
// Passing nil for opts is the same as passing the zero value.
//
// Once a connection is created, it will immediately start receiving
// requests from the transport.
func NewConn(t Transport, opts *Options) *Conn {
	bgctx, bgcancel := context.WithCancel(context.Background())
	c := &Conn{
		sender:     t,
		recvCloser: t,
		shut:       make(chan struct{}),
		bgctx:      bgctx,
		bgcancel:   bgcancel,
	}
	if opts != nil {
		c.bootstrap = opts.BootstrapClient
		c.reporter = opts.ErrorReporter
	}
	c.tasks.Add(1)
	go func() {
		abortErr := c.receive(c.bgctx, t)
		c.tasks.Done()

		c.mu.Lock()
		select {
		case <-c.bgctx.Done():
			c.mu.Unlock()
		default:
			if abortErr != nil {
				c.report(abortErr)
			}
			// shutdown unlocks c.mu.
			if err := c.shutdown(abortErr); err != nil {
				c.report(err)
			}
		}
	}()
	return c
}

// Bootstrap returns the remote vat's bootstrap interface.  This creates
// a new client that the caller is responsible for releasing.
func (c *Conn) Bootstrap(ctx context.Context) *capnp.Client {
	c.mu.Lock()
	if !c.startTask() {
		c.mu.Unlock()
		return capnp.ErrorClient(disconnected("connection closed"))
	}
	defer c.tasks.Done()
	id := questionID(c.questionID.next())
	err := c.sendMessage(ctx, func(msg rpccp.Message) error {
		boot, err := msg.NewBootstrap()
		if err != nil {
			return err
		}
		boot.SetQuestionId(uint32(id))
		return nil
	})
	if err != nil {
		c.questionID.remove(uint32(id))
		c.mu.Unlock()
		return capnp.ErrorClient(annotate(err).errorf("bootstrap"))
	}
	ctx, cancel := context.WithCancel(ctx)
	q := c.newQuestion(ctx, id, capnp.Method{})
	bc, cp := capnp.NewPromisedClient(bootstrapClient{
		c:      q.p.Answer().Client().AddRef(),
		cancel: cancel,
	})
	q.bootstrapPromise = cp // safe to write because we're still holding c.mu
	c.mu.Unlock()
	return bc
}

type bootstrapClient struct {
	c      *capnp.Client
	cancel context.CancelFunc
}

func (bc bootstrapClient) Send(ctx context.Context, s capnp.Send) (*capnp.Answer, capnp.ReleaseFunc) {
	return bc.c.SendCall(ctx, s)
}

func (bc bootstrapClient) Recv(ctx context.Context, r capnp.Recv) capnp.PipelineCaller {
	return bc.c.RecvCall(ctx, r)
}

func (bc bootstrapClient) Brand() interface{} {
	return bc.c.Brand()
}

func (bc bootstrapClient) Shutdown() {
	bc.cancel()
	bc.c.Release()
}

// Close sends an abort to the remote vat and closes the underlying
// transport.
func (c *Conn) Close() error {
	c.mu.Lock()
	if c.closed {
		return fail("close on closed connection")
	}
	c.closed = true
	select {
	case <-c.bgctx.Done():
		c.mu.Unlock()
		<-c.shut
		return nil
	default:
		// shutdown unlocks c.mu.
		return c.shutdown(errors.New(errors.Failed, "", "connection closed"))
	}
}

// Done returns a channel that is closed after the connection is
// shut down.
func (c *Conn) Done() <-chan struct{} {
	return c.shut
}

// shutdown tears down the connection and transport, optionally sending
// an abort message before closing.  The caller must be holding onto
// c.mu, although it will be released while shutting down, and c.bgctx
// must not be Done.
func (c *Conn) shutdown(abortErr error) error {
	defer close(c.shut)

	// Cancel all work.
	c.bgcancel()
	for _, a := range c.answers {
		if a != nil {
			a.cancel()
		}
	}
	c.mu.Unlock()
	rerr := c.recvCloser.CloseRecv()

	// Wait for work to stop.
	c.tasks.Wait()

	// Clear all tables, releasing exported clients and unfinished answers.
	c.mu.Lock()
	exports := c.exports
	answers := c.answers
	c.imports = nil
	c.exports = nil
	c.questions = nil
	c.answers = nil
	c.mu.Unlock()

	c.bootstrap.Release()
	c.bootstrap = nil
	for _, e := range exports {
		if e != nil {
			e.client.Release()
		}
	}
	for _, a := range answers {
		if a != nil {
			// Because shutdown is now the only task running, no need to
			// acquire sender lock.
			a.s.releaseMsg()
		}
	}

	// Send abort message (ignoring error).
	if abortErr != nil {
		// TODO(soon): add timeout
		msg, send, release, err := c.sender.NewMessage(context.Background())
		if err != nil {
			goto closeSend
		}
		abort, err := msg.NewAbort()
		if err != nil {
			release()
			goto closeSend
		}
		abort.SetType(rpccp.Exception_Type(errors.TypeOf(abortErr)))
		if err := abort.SetReason(abortErr.Error()); err != nil {
			release()
			goto closeSend
		}
		send()
		release()
	}
closeSend:
	serr := c.sender.CloseSend()

	if rerr != nil {
		return errorf("close transport: %v", rerr)
	}
	if serr != nil {
		return errorf("close transport: %v", serr)
	}
	return nil
}

// receive receives and dispatches messages coming from r.  receive runs
// in a background goroutine.
//
// After receive returns, the connection is shut down.  If receive
// returns a non-nil error, it is sent to the remove vat as an abort.
func (c *Conn) receive(ctx context.Context, r Receiver) error {
	for {
		recv, releaseRecv, err := r.RecvMessage(ctx)
		if err != nil {
			return err
		}
		switch recv.Which() {
		case rpccp.Message_Which_unimplemented:
			// no-op for now to avoid feedback loop
		case rpccp.Message_Which_abort:
			exc, err := recv.Abort()
			if err != nil {
				releaseRecv()
				c.reportf("read abort: %v", err)
				return nil
			}
			reason, err := exc.Reason()
			if err != nil {
				releaseRecv()
				c.reportf("read abort reason: %v", err)
				return nil
			}
			ty := exc.Type()
			releaseRecv()
			c.report(errors.New(errors.Type(ty), "rpc", "remote abort: "+reason))
			return nil
		case rpccp.Message_Which_bootstrap:
			bootstrap, err := recv.Bootstrap()
			if err != nil {
				releaseRecv()
				c.reportf("read bootstrap: %v", err)
				continue
			}
			qid := answerID(bootstrap.QuestionId())
			releaseRecv()
			if err := c.handleBootstrap(ctx, qid); err != nil {
				return err
			}
		case rpccp.Message_Which_call:
			call, err := recv.Call()
			if err != nil {
				releaseRecv()
				c.reportf("read call: %v", err)
				continue
			}
			if err := c.handleCall(ctx, call, releaseRecv); err != nil {
				return err
			}
		case rpccp.Message_Which_return:
			ret, err := recv.Return()
			if err != nil {
				releaseRecv()
				c.reportf("read return: %v", err)
				continue
			}
			if err := c.handleReturn(ctx, ret, releaseRecv); err != nil {
				return err
			}
		case rpccp.Message_Which_finish:
			fin, err := recv.Finish()
			if err != nil {
				releaseRecv()
				c.reportf("read finish: %v", err)
				continue
			}
			qid := answerID(fin.QuestionId())
			releaseResultCaps := fin.ReleaseResultCaps()
			releaseRecv()
			if err := c.handleFinish(ctx, qid, releaseResultCaps); err != nil {
				return err
			}
		case rpccp.Message_Which_release:
			rel, err := recv.Release()
			if err != nil {
				releaseRecv()
				c.reportf("read release: %v", err)
				continue
			}
			id := exportID(rel.Id())
			count := rel.ReferenceCount()
			releaseRecv()
			if err := c.handleRelease(ctx, id, count); err != nil {
				return err
			}
		default:
			err := c.handleUnknownMessage(ctx, recv)
			releaseRecv()
			if err != nil {
				return err
			}
		}
	}
}

func (c *Conn) handleBootstrap(ctx context.Context, id answerID) error {
	c.mu.Lock()
	if c.answers[id] != nil {
		c.mu.Unlock()
		return errorf("incoming bootstrap: answer ID %d reused", id)
	}
	ans, err := c.newAnswer(ctx, id, func() {})
	if err != nil {
		c.mu.Unlock()
		c.report(annotate(err).errorf("incoming bootstrap"))
		return nil
	}
	if !c.bootstrap.IsValid() {
		ans.sendException(errors.New(errors.Failed, "", "vat does not expose a public/bootstrap interface"))
		c.mu.Unlock()
		return nil
	}
	if err := ans.setBootstrap(c.bootstrap.AddRef()); err != nil {
		ans.sendException(err)
		c.mu.Unlock()
		return nil
	}
	refs := ans.sendReturn()
	c.mu.Unlock()
	if len(refs) > 0 {
		// Answer cannot possibly encounter a Finish, since we still
		// haven't returned to receive().
		panic("unreachable")
	}
	return nil
}

func (c *Conn) handleCall(ctx context.Context, call rpccp.Call, releaseCall capnp.ReleaseFunc) error {
	id := answerID(call.QuestionId())
	if call.SendResultsTo().Which() != rpccp.Call_sendResultsTo_Which_caller {
		// TODO(someday): handle SendResultsTo.yourself
		c.reportf("incoming call: results destination is not caller")
		c.mu.Lock()
		err := c.sendMessage(ctx, func(m rpccp.Message) error {
			mm, err := m.NewUnimplemented()
			if err != nil {
				return err
			}
			if err := mm.SetCall(call); err != nil {
				return err
			}
			return nil
		})
		c.mu.Unlock()
		releaseCall()
		if err != nil {
			c.report(annotate(err).errorf("incoming call: send unimplemented"))
		}
		return nil
	}

	c.mu.Lock()
	if c.answers[id] != nil {
		c.mu.Unlock()
		releaseCall()
		return errorf("incoming call: answer ID %d reused", id)
	}
	callCtx, cancel := context.WithCancel(c.bgctx)
	ans, err := c.newAnswer(ctx, id, cancel)
	if err != nil {
		c.mu.Unlock()
		releaseCall()
		c.report(annotate(err).errorf("incoming call"))
		cancel()
		return nil
	}
	var p parsedCall
	if err := c.parseCall(&p, call); err != nil {
		err = annotate(err).errorf("incoming call")
		ans.sendException(err)
		c.mu.Unlock()
		cancel()
		releaseCall()
		c.report(err)
		return nil
	}
	switch p.targetWhich {
	case rpccp.MessageTarget_Which_importedCap:
		ent := c.findExport(p.importedCap)
		if ent == nil {
			c.mu.Unlock()
			releaseCall()
			cancel()
			return errorf("incoming call: unknown export ID %d", id)
		}
		c.tasks.Add(1) // will be finished by answer.Return
		c.mu.Unlock()
		pcall := ent.client.RecvCall(callCtx, capnp.Recv{
			Args:        p.args,
			Method:      p.method,
			ReleaseArgs: releaseCall,
			Returner:    ans,
		})
		// Place PipelineCaller into answer.  Since the receive goroutine is
		// the only one that uses answer.pcall, it's fine that there's a
		// time gap for this being set.
		ans.setPipelineCaller(pcall)
		return nil
	case rpccp.MessageTarget_Which_promisedAnswer:
		tgtAns := c.answers[p.promisedAnswer]
		if tgtAns == nil {
			c.mu.Unlock()
			releaseCall()
			cancel()
			return errorf("incoming call: unknown target answer ID %d", p.promisedAnswer)
		}
		if tgtAns.flags&resultsReady != 0 {
			// Results ready.
			if tgtAns.err != nil {
				ans.sendException(tgtAns.err)
				c.mu.Unlock()
				releaseCall()
				cancel()
				return nil
			}
			// tgtAns.results is guaranteed to stay alive because it hasn't
			// received finish yet (it would have been deleted from the
			// answers table), and it can't receive a finish because this is
			// happening on the receive goroutine.
			content, err := tgtAns.results.Content()
			if err != nil {
				err = errorf("incoming call: read results from target answer: %v")
				ans.sendException(err)
				c.mu.Unlock()
				releaseCall()
				cancel()
				c.report(err)
				return nil
			}
			sub, err := capnp.Transform(content, p.transform)
			if err != nil {
				// Not reporting, as this is the caller's fault.
				ans.sendException(err)
				c.mu.Unlock()
				releaseCall()
				cancel()
				return nil
			}
			tgt := sub.Interface().Client()
			c.tasks.Add(1) // will be finished by answer.Return
			c.mu.Unlock()
			pcall := tgt.RecvCall(callCtx, capnp.Recv{
				Args:        p.args,
				Method:      p.method,
				ReleaseArgs: releaseCall,
				Returner:    ans,
			})
			ans.setPipelineCaller(pcall)
		} else {
			// Results not ready, use pipeline caller.
			tgtAns.pcalls.Add(1)
			tgt := tgtAns.pcall
			c.tasks.Add(1) // will be finished by answer.Return
			c.mu.Unlock()
			pcall := tgt.PipelineRecv(callCtx, p.transform, capnp.Recv{
				Args:        p.args,
				Method:      p.method,
				ReleaseArgs: releaseCall,
				Returner:    ans,
			})
			tgtAns.pcalls.Done()
			ans.setPipelineCaller(pcall)
		}
		return nil
	default:
		panic("unreachable")
	}
}

type parsedCall struct {
	method capnp.Method
	args   capnp.Struct

	targetWhich    rpccp.MessageTarget_Which
	importedCap    exportID
	promisedAnswer answerID
	transform      []capnp.PipelineOp
}

func (c *Conn) parseCall(p *parsedCall, call rpccp.Call) error {
	p.method = capnp.Method{
		InterfaceID: call.InterfaceId(),
		MethodID:    call.MethodId(),
	}
	payload, err := call.Params()
	if err != nil {
		return errorf("read params: %v", err)
	}
	ptr, err := c.recvPayload(payload)
	if err != nil {
		return annotate(err).errorf("read params")
	}
	p.args = ptr.Struct()
	tgt, err := call.Target()
	if err != nil {
		return errorf("read target: %v", err)
	}
	p.targetWhich = tgt.Which()
	switch p.targetWhich {
	case rpccp.MessageTarget_Which_importedCap:
		p.importedCap = exportID(tgt.ImportedCap())
	case rpccp.MessageTarget_Which_promisedAnswer:
		pa, err := tgt.PromisedAnswer()
		if err != nil {
			return errorf("read target answer: %v", err)
		}
		p.promisedAnswer = answerID(pa.QuestionId())
		opList, err := pa.Transform()
		if err != nil {
			return errorf("read target transform: %v", err)
		}
		p.transform, err = parseTransform(opList)
		if err != nil {
			return annotate(err).errorf("read target transform")
		}
	default:
		return unimplementedf("unknown message target %v", p.targetWhich)
	}
	return nil
}

func parseTransform(list rpccp.PromisedAnswer_Op_List) ([]capnp.PipelineOp, error) {
	ops := make([]capnp.PipelineOp, 0, list.Len())
	for i := 0; i < list.Len(); i++ {
		li := list.At(i)
		switch li.Which() {
		case rpccp.PromisedAnswer_Op_Which_noop:
			// do nothing
		case rpccp.PromisedAnswer_Op_Which_getPointerField:
			ops = append(ops, capnp.PipelineOp{Field: li.GetPointerField()})
		default:
			return nil, errorf("transform element %d: unknown type %v", i, li.Which())
		}
	}
	return ops, nil
}

func (c *Conn) handleReturn(ctx context.Context, ret rpccp.Return, releaseRet capnp.ReleaseFunc) error {
	c.mu.Lock()
	qid := questionID(ret.AnswerId())
	if uint32(qid) >= uint32(len(c.questions)) {
		c.mu.Unlock()
		releaseRet()
		return errorf("incoming return: question %d does not exist", qid)
	}
	// Pop the question from the table.  Receiving the Return message
	// will always remove the question from the table, because it's the
	// only time the remote vat will use it.
	q := c.questions[qid]
	c.questions[qid] = nil
	if q == nil {
		c.mu.Unlock()
		releaseRet()
		return errorf("incoming return: question %d does not exist", qid)
	}
	canceled := q.flags&finished != 0
	q.flags |= finished
	if canceled {
		// Wait for cancelation task to write the Finish message.  If the
		// Finish message could not be sent to the remote vat, we can't
		// reuse the ID.
		select {
		case <-q.finishMsgSend:
			if q.flags&finishSent != 0 {
				c.questionID.remove(uint32(qid))
			}
			c.mu.Unlock()
			releaseRet()
		default:
			c.mu.Unlock()
			releaseRet()
			<-q.finishMsgSend
			c.mu.Lock()
			if q.flags&finishSent != 0 {
				c.questionID.remove(uint32(qid))
			}
			c.mu.Unlock()
		}
		return nil
	}
	pr := c.parseReturn(ret)
	if pr.parseFailed {
		c.report(annotate(pr.err).errorf("incoming return"))
	}
	switch {
	case q.bootstrapPromise != nil && pr.err == nil:
		q.release = func() {}
		c.mu.Unlock()
		q.p.Fulfill(pr.result)
		q.bootstrapPromise.Fulfill(q.p.Answer().Client())
		q.p.ReleaseClients()
		releaseRet()
		q.done()
		c.mu.Lock()
	case q.bootstrapPromise != nil && pr.err != nil:
		// TODO(someday): send unimplemented message back to remote if
		// pr.unimplemented == true.
		q.release = func() {}
		c.mu.Unlock()
		q.p.Reject(pr.err)
		q.bootstrapPromise.Fulfill(q.p.Answer().Client())
		q.p.ReleaseClients()
		releaseRet()
		q.done()
		c.mu.Lock()
	case q.bootstrapPromise == nil && pr.err != nil:
		// TODO(someday): send unimplemented message back to remote if
		// pr.unimplemented == true.
		q.release = func() {}
		c.mu.Unlock()
		q.p.Reject(pr.err)
		releaseRet()
		q.done()
		c.mu.Lock()
	default:
		q.release = releaseRet
		c.mu.Unlock()
		q.p.Fulfill(pr.result)
		q.done()
		c.mu.Lock()
	}
	err := c.sendMessage(ctx, func(msg rpccp.Message) error {
		fin, err := msg.NewFinish()
		if err != nil {
			return err
		}
		fin.SetQuestionId(uint32(qid))
		fin.SetReleaseResultCaps(false)
		return nil
	})
	if err != nil {
		close(q.finishMsgSend)
		c.mu.Unlock()
		c.report(annotate(err).errorf("incoming return: send finish"))
		return nil
	}
	q.flags |= finishSent
	c.questionID.remove(uint32(qid))
	close(q.finishMsgSend)
	c.mu.Unlock()
	return nil
}

func (c *Conn) parseReturn(ret rpccp.Return) parsedReturn {
	switch ret.Which() {
	case rpccp.Return_Which_results:
		r, err := ret.Results()
		if err != nil {
			return parsedReturn{err: errorf("parse return: %v", err), parseFailed: true}
		}
		content, err := c.recvPayload(r)
		if err != nil {
			return parsedReturn{err: errorf("parse return: %v", err), parseFailed: true}
		}
		return parsedReturn{result: content}
	case rpccp.Return_Which_exception:
		exc, err := ret.Exception()
		if err != nil {
			return parsedReturn{err: errorf("parse return: %v", err), parseFailed: true}
		}
		reason, err := exc.Reason()
		if err != nil {
			return parsedReturn{err: errorf("parse return: %v", err), parseFailed: true}
		}
		return parsedReturn{err: errors.New(errors.Type(exc.Type()), "", reason)}
	default:
		w := ret.Which()
		return parsedReturn{err: errorf("parse return: unhandled type %v", w), parseFailed: true, unimplemented: true}
	}
}

type parsedReturn struct {
	result        capnp.Ptr
	err           error
	parseFailed   bool
	unimplemented bool
}

func (c *Conn) handleFinish(ctx context.Context, id answerID, releaseResultCaps bool) error {
	c.mu.Lock()
	ans := c.answers[id]
	if ans == nil {
		c.mu.Unlock()
		return errorf("incoming finish: unknown answer ID %d", id)
	}
	if ans.flags&finishReceived != 0 {
		c.mu.Unlock()
		return errorf("incoming finish: answer ID %d already received finish", id)
	}
	ans.flags |= finishReceived
	ans.cancel()
	if ans.flags&returnSent == 0 {
		if releaseResultCaps {
			ans.flags |= releaseResultCapsFlag
		}
		c.mu.Unlock()
		return nil
	}
	ans.s.acquireSender()
	ans.s.finish()
	delete(c.answers, id)
	if !releaseResultCaps {
		c.mu.Unlock()
		return nil
	}
	rl, err := c.releaseExports(ans.exportRefs)
	c.mu.Unlock()
	rl.release()
	if err != nil {
		return annotate(err).errorf("incoming finish: release result caps")
	}
	return nil
}

// recvCap materializes a client for a given descriptor.  If there is an
// error reading a descriptor, then the resulting client will return the
// error whenever it is called.  The caller must be holding onto c.mu.
func (c *Conn) recvCap(d rpccp.CapDescriptor) *capnp.Client {
	switch d.Which() {
	case rpccp.CapDescriptor_Which_none:
		return nil
	case rpccp.CapDescriptor_Which_senderHosted:
		id := importID(d.SenderHosted())
		return c.addImport(id)
	default:
		return capnp.ErrorClient(errorf("unknown CapDescriptor type %v", d.Which()))
	}
}

// recvPayload extracts the content pointer after populating the
// message's capability table.  The caller must be holding onto c.mu.
func (c *Conn) recvPayload(payload rpccp.Payload) (capnp.Ptr, error) {
	if payload.Message().CapTable != nil {
		// RecvMessage likely violated its invariant.
		return capnp.Ptr{}, fail("read payload: capability table already populated")
	}
	p, err := payload.Content()
	if err != nil {
		return capnp.Ptr{}, errorf("read payload: %v", err)
	}
	ptab, err := payload.CapTable()
	if err != nil {
		// Don't allow unreadable capability table to stop other results,
		// just present an empty capability table.
		c.reportf("read payload: capability table: %v", err)
		return p, nil
	}
	mtab := make([]*capnp.Client, ptab.Len())
	for i := 0; i < ptab.Len(); i++ {
		mtab[i] = c.recvCap(ptab.At(i))
	}
	payload.Message().CapTable = mtab
	return p, nil
}

func (c *Conn) handleRelease(ctx context.Context, id exportID, count uint32) error {
	c.mu.Lock()
	client, err := c.releaseExport(id, count)
	c.mu.Unlock()
	if err != nil {
		return annotate(err).errorf("incoming release")
	}
	client.Release() // no-ops for nil
	return nil
}

func (c *Conn) handleUnknownMessage(ctx context.Context, recv rpccp.Message) error {
	c.reportf("unknown message type %v from remote", recv.Which())
	c.mu.Lock()
	err := c.sendMessage(ctx, func(msg rpccp.Message) error {
		return msg.SetUnimplemented(recv)
	})
	c.mu.Unlock()
	if err != nil {
		c.report(annotate(err).errorf("send unimplemented"))
	}
	return nil
}

// startTask increments c.tasks if c is not shutting down.
// It returns whether c.tasks was incremented.
//
// The caller must be holding onto c.mu.
func (c *Conn) startTask() bool {
	select {
	case <-c.bgctx.Done():
		return false
	default:
		c.tasks.Add(1)
		return true
	}
}

// sendSession manages the lifecycle of an outbound message.
type sendSession struct {
	msg        rpccp.Message
	send       func() error // can be called directly
	c          *Conn
	releaseMsg capnp.ReleaseFunc
}

// startSend creates a new outbound message.  The caller must be holding
// c.mu, but importantly, if newMessage does not return an error, then
// newMessage releases c.mu and acquires the sender lock.
//
// startSend will return an error if Close() has been called on c or ctx
// is canceled or reaches its deadline.
func (c *Conn) startSend(ctx context.Context) (sendSession, error) {
	for {
		select {
		case <-c.bgctx.Done():
			return sendSession{}, disconnected("connection closed")
		default:
		}
		s := c.sendCond
		if s == nil {
			break
		}
		c.mu.Unlock()
		select {
		case <-s:
		case <-ctx.Done():
			c.mu.Lock()
			return sendSession{}, ctx.Err()
		case <-c.bgctx.Done():
			c.mu.Lock()
			return sendSession{}, disconnected("connection closed")
		}
		c.mu.Lock()
	}
	c.sendCond = make(chan struct{})
	c.mu.Unlock()
	msg, send, release, err := c.sender.NewMessage(ctx)
	if err != nil {
		c.mu.Lock()
		close(c.sendCond)
		c.sendCond = nil
		return sendSession{}, errorf("create message: %v", err)
	}
	return sendSession{msg, send, c, release}, nil
}

// releaseSender trades the sender lock for s.c.mu.
func (s sendSession) releaseSender() {
	s.c.mu.Lock()
	close(s.c.sendCond)
	s.c.sendCond = nil
}

// acquireSender trades back s.c.mu for the sender lock after a call to
// releaseSender.
func (s sendSession) acquireSender() {
	s.c.sendCond = make(chan struct{})
	s.c.mu.Unlock()
}

// finish releases the message and unblocks Conn.Close completing,
// if it is currently being called.  The caller must be holding the
// sender lock, and it will be traded for s.c.mu.
func (s sendSession) finish() {
	s.releaseMsg()
	s.releaseSender()
}

// sendMessage creates a new message on the Sender, calls f to build it,
// and sends it if f does not return an error.  The caller must be
// holding onto c.mu.  While f is being called, it will be holding onto
// the sender lock, but not c.mu.
func (c *Conn) sendMessage(ctx context.Context, f func(msg rpccp.Message) error) error {
	s, err := c.startSend(ctx)
	if err != nil {
		return err
	}
	defer s.finish()
	if err := f(s.msg); err != nil {
		return errorf("build message: %v", err)
	}
	if err := s.send(); err != nil {
		return errorf("send message: %v", err)
	}
	return nil
}

// report sends an error to c's reporter.  The caller does not have to
// be holding c.mu.
func (c *Conn) report(err error) {
	if c.reporter == nil {
		return
	}
	c.reporter.ReportError(err)
}

// reportf formats an error and sends it to c's reporter.
func (c *Conn) reportf(format string, args ...interface{}) {
	if c.reporter == nil {
		return
	}
	c.reporter.ReportError(errorf(format, args...))
}

// idgen returns a sequence of monotonically increasing IDs with
// support for replacement.  The zero value is a generator that
// starts at zero.
type idgen struct {
	i    uint32
	free []uint32
}

func (gen *idgen) next() uint32 {
	if n := len(gen.free); n > 0 {
		i := gen.free[n-1]
		gen.free = gen.free[:n-1]
		return i
	}
	i := gen.i
	gen.i++
	return i
}

func (gen *idgen) remove(i uint32) {
	gen.free = append(gen.free, i)
}

func fail(msg string) error {
	return errors.New(errors.Failed, "rpc", msg)
}

func disconnected(msg string) error {
	return errors.New(errors.Disconnected, "rpc", msg)
}

func errorf(format string, args ...interface{}) error {
	return fail(fmt.Sprintf(format, args...))
}

func unimplementedf(format string, args ...interface{}) error {
	return errors.New(errors.Unimplemented, "rpc", fmt.Sprintf(format, args...))
}

type annotater struct {
	err error
}

func annotate(err error) annotater {
	return annotater{err}
}

func (a annotater) errorf(format string, args ...interface{}) error {
	return errors.Annotate("rpc", fmt.Sprintf(format, args...), a.err)
}
