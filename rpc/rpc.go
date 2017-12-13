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

	// mu protects all the following fields in the Conn.  mu should not be
	// held while making calls that take indeterminate time (I/O or
	// application code).  Condition channels protect operations on any
	// field that take an indeterminate amount of time.  Thus, critical
	// sections involving mu are quite short, while still ensuring
	// mutually exclusive access to resources.
	mu sync.Mutex

	bgctx    context.Context
	bgcancel context.CancelFunc
	bgtasks  sync.WaitGroup
	closed   bool          // set when Close() is called, used to distinguish user Closing multiple times
	shut     chan struct{} // closed when shutdown() returns

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

	// openSends is added to before creating a new message (while holding
	// c.mu) and subtracted from after releasing the message.
	openSends sync.WaitGroup

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
	c.bgtasks.Add(1)
	go func() {
		abortErr := c.receive(c.bgctx, t)
		c.bgtasks.Done()

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
	defer c.mu.Unlock()
	c.mu.Lock()
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
		return capnp.ErrorClient(annotate(err).errorf("bootstrap"))
	}
	ctx, cancel := context.WithCancel(ctx)
	q := c.newQuestion(ctx, id, capnp.Method{})
	bc, cp := capnp.NewPromisedClient(bootstrapClient{
		c:      q.p.Answer().Client().AddRef(),
		cancel: cancel,
	})
	q.bootstrapPromise = cp // safe to write because we're still holding c.mu
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

// shutdown tears down the connection and transport, optionally sending
// an abort message before closing.  The caller must be holding onto
// c.mu, although it will be released while shutting down, and c.bgctx
// must not be Done.
func (c *Conn) shutdown(abortErr error) error {
	// Mark closed and stop receiving messages.
	defer close(c.shut)
	c.bgcancel()
	c.mu.Unlock()
	rerr := c.recvCloser.CloseRecv()

	// Wait for other tasks to stop.
	c.bgtasks.Wait()

	// Release exported clients and ongoing answers.
	c.mu.Lock()
	exports := c.exports
	answers := c.answers
	c.imports = nil
	c.exports = nil
	c.questions = nil
	c.answers = nil
	for _, a := range answers {
		if a != nil {
			a.cancel()
			a.s.acquireSender()
			a.s.finish()
			// TODO(soon): release result caps (while not holding c.mu)
		}
	}
	c.mu.Unlock()
	c.bootstrap.Release()
	c.bootstrap = nil
	for _, e := range exports {
		if e != nil {
			e.client.Release()
		}
	}

	// Wait for all other sends to finish.
	c.openSends.Wait()

	// Send abort message (ignoring error).  No locking needed, since
	// c.startSend will always return errors after closing bgctx.
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
	defer c.mu.Unlock()
	c.mu.Lock()
	if c.answers[id] != nil {
		return errorf("incoming bootstrap: answer ID %d reused", id)
	}
	ans, err := c.newAnswer(ctx, id, func() {})
	if err != nil {
		c.report(annotate(err).errorf("incoming bootstrap"))
		return nil
	}
	if c.bootstrap.IsValid() {
		err := ans.setBootstrap(c.bootstrap.AddRef())
		ans.lockedReturn(err)
	} else {
		ans.lockedReturn(errors.New(errors.Failed, "", "vat does not expose a public/bootstrap interface"))
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
		return nil
	}
	tgt, err := call.Target()
	if err != nil {
		retErr := errorf("incoming call: read target: %v", err)
		ans.lockedReturn(retErr)
		c.mu.Unlock()
		releaseCall()
		c.report(retErr)
		return nil
	}
	argsPayload, err := call.Params()
	if err != nil {
		retErr := errorf("incoming call: read params: %v", err)
		ans.lockedReturn(retErr)
		c.mu.Unlock()
		releaseCall()
		c.report(retErr)
		return nil
	}
	args, err := c.recvPayload(argsPayload)
	if err != nil {
		retErr := annotate(err).errorf("incoming call: read params")
		ans.lockedReturn(retErr)
		c.mu.Unlock()
		releaseCall()
		c.report(retErr)
		return nil
	}
	method := capnp.Method{
		InterfaceID: call.InterfaceId(),
		MethodID:    call.MethodId(),
	}
	switch tgt.Which() {
	case rpccp.MessageTarget_Which_importedCap:
		id := exportID(tgt.ImportedCap())
		ent := c.findExport(id)
		if ent == nil {
			c.mu.Unlock()
			releaseCall()
			return errorf("incoming call: unknown export ID %d", id)
		}
		c.mu.Unlock()
		pcall := ent.client.RecvCall(callCtx, capnp.Recv{
			Args:        args.Struct(),
			Method:      method,
			ReleaseArgs: releaseCall,
			Returner:    ans,
		})
		// Place PipelineCaller into answer.  Since the receive goroutine is
		// the only one that uses answer.pcall, it's fine that there's a
		// time gap for this being set.
		ans.setPipelineCaller(pcall)
		return nil
	case rpccp.MessageTarget_Which_promisedAnswer:
		pa, err := tgt.PromisedAnswer()
		if err != nil {
			retErr := errorf("incoming call: read target answer: %v", err)
			ans.lockedReturn(retErr)
			c.mu.Unlock()
			releaseCall()
			c.report(retErr)
			return nil
		}
		tgtID := answerID(pa.QuestionId())
		tgtAns := c.answers[tgtID]
		if tgtAns == nil {
			c.mu.Unlock()
			releaseCall()
			return errorf("incoming call: unknown target answer ID %d", tgtID)
		}
		opList, err := pa.Transform()
		if err != nil {
			retErr := errorf("incoming call: read target transform: %v", err)
			ans.lockedReturn(retErr)
			c.mu.Unlock()
			releaseCall()
			c.report(retErr)
			return nil
		}
		xform, err := parseTransform(opList)
		if err != nil {
			retErr := annotate(err).errorf("incoming call: read target transform")
			ans.lockedReturn(retErr)
			c.mu.Unlock()
			releaseCall()
			c.report(retErr)
			return nil
		}
		if tgtAns.state&4 != 0 {
			// Results ready.
			if tgtAns.err != nil {
				ans.lockedReturn(tgtAns.err)
				c.mu.Unlock()
				releaseCall()
				return nil
			}
			// tgtAns.results is guaranteed to stay alive because it hasn't
			// received finish yet (it would have been deleted from the
			// answers table), and it can't receive a finish because this is
			// happening on the receive goroutine.
			content, err := tgtAns.results.Content()
			if err != nil {
				retErr := errorf("incoming call: read results from target answer: %v")
				ans.lockedReturn(retErr)
				c.mu.Unlock()
				releaseCall()
				c.report(retErr)
				return nil
			}
			sub, err := capnp.Transform(content, xform)
			if err != nil {
				// Not reporting, as this is the caller's fault.
				ans.lockedReturn(err)
				c.mu.Unlock()
				releaseCall()
				return nil
			}
			tgt := sub.Interface().Client()
			c.mu.Unlock()
			pcall := tgt.RecvCall(callCtx, capnp.Recv{
				Args:        args.Struct(),
				Method:      method,
				ReleaseArgs: releaseCall,
				Returner:    ans,
			})
			ans.setPipelineCaller(pcall)
		} else {
			// Results not ready, use pipeline caller.
			tgtAns.pcalls.Add(1)
			tgt := tgtAns.pcall
			c.mu.Unlock()
			pcall := tgt.PipelineRecv(callCtx, xform, capnp.Recv{
				Args:        args.Struct(),
				Method:      method,
				ReleaseArgs: releaseCall,
				Returner:    ans,
			})
			tgtAns.pcalls.Done()
			ans.setPipelineCaller(pcall)
		}
		return nil
	default:
		retErr := unimplementedf("incoming call: unknown message target %v", tgt.Which())
		ans.lockedReturn(retErr)
		c.mu.Unlock()
		releaseCall()
		c.report(retErr)
		return nil
	}
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
	q := c.questions[qid]
	if q == nil {
		c.mu.Unlock()
		releaseRet()
		return errorf("incoming return: question %d does not exist", qid)
	}
	defer c.mu.Unlock()
	sentFinish := q.sentFinish()
	q.state |= 3 // return received and sending finish
	c.questions[qid] = nil

	pr := c.parseReturn(ret)
	if pr.parseFailed {
		c.report(annotate(pr.err).errorf("incoming return"))
	}
	if pr.err == nil {
		if q.bootstrapPromise != nil {
			q.release = func() {}
		} else {
			q.release = releaseRet
		}
		c.mu.Unlock()
		q.p.Fulfill(pr.result)
		if q.bootstrapPromise != nil {
			q.bootstrapPromise.Fulfill(q.p.Answer().Client())
			q.p.ReleaseClients()
			releaseRet()
		}
		q.done()
		c.mu.Lock()
	} else {
		// TODO(someday): send unimplemented message back to remote if
		// pr.unimplemented == true.
		q.release = func() {}
		c.mu.Unlock()
		q.p.Reject(pr.err)
		if q.bootstrapPromise != nil {
			q.bootstrapPromise.Fulfill(q.p.Answer().Client())
			q.p.ReleaseClients()
		}
		releaseRet()
		q.done()
		c.mu.Lock()
	}
	if !sentFinish {
		err := c.sendMessage(ctx, func(msg rpccp.Message) error {
			fin, err := msg.NewFinish()
			if err != nil {
				return err
			}
			fin.SetQuestionId(uint32(qid))
			fin.SetReleaseResultCaps(false)
			return nil
		})
		c.questionID.remove(uint32(qid))
		if err != nil {
			c.report(annotate(err).errorf("incoming return: send finish"))
			return nil
		}
	} else {
		c.questionID.remove(uint32(qid))
	}
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
	defer c.mu.Unlock()
	c.mu.Lock()
	ans := c.answers[id]
	if ans == nil {
		return errorf("incoming finish: unknown answer ID %d", id)
	}
	ans.state |= 2
	ans.cancel()
	if !ans.isDone() {
		// Not returned yet.
		// TODO(soon): record releaseResultCaps
		return nil
	}
	ans.s.acquireSender()
	ans.s.finish()
	delete(c.answers, id)
	if releaseResultCaps {
		// TODO(soon): release result caps (while not holding c.mu)
	}
	return nil
}

// sendCap writes a client descriptor.  The caller must be holding onto c.mu.
func (c *Conn) sendCap(d rpccp.CapDescriptor, client *capnp.Client) {
	if !client.IsValid() {
		d.SetNone()
		return
	}

	brand := client.Brand()
	if ic, ok := brand.(*importClient); ok && ic.conn == c {
		if ent := c.imports[ic.id]; ent != nil && ent.generation == ic.generation {
			d.SetReceiverHosted(uint32(ic.id))
			return
		}
	}
	// TODO(someday): Check for unresolved client for senderPromise.
	// TODO(someday): Check for pipeline client on question for receiverAnswer.

	// Default to sender-hosted (export).
	for i, ent := range c.exports {
		if ent.client.IsSame(client) {
			ent.wireRefs++
			d.SetSenderHosted(uint32(i))
			return
		}
	}
	c.exports = append(c.exports, &expent{
		client:   client.AddRef(),
		wireRefs: 1,
	})
	d.SetSenderHosted(uint32(len(c.exports) - 1))
}

// fillPayloadCapTable adds descriptors of payload's message's
// capabilities into payload's capability table.  The caller must be
// holding onto c.mu.
func (c *Conn) fillPayloadCapTable(payload rpccp.Payload) error {
	tab := payload.Message().CapTable
	if len(tab) == 0 {
		return nil
	}
	list, err := payload.NewCapTable(int32(len(tab)))
	if err != nil {
		return errorf("payload capability table: %v", err)
	}
	for i, client := range tab {
		c.sendCap(list.At(i), client)
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
	ent := c.findExport(id)
	if ent == nil {
		c.mu.Unlock()
		return errorf("incoming release: unknown export ID %d", id)
	}
	ent.wireRefs -= int(count)
	if ent.wireRefs < 0 {
		c.mu.Unlock()
		return errorf("incoming release: export ID %d released too many references", id)
	}
	if ent.wireRefs > 0 {
		c.mu.Unlock()
		return nil
	}
	client := ent.client
	c.exports[id] = nil
	c.exportID.remove(uint32(id))
	c.mu.Unlock()
	client.Release()
	return nil
}

func (c *Conn) handleUnknownMessage(ctx context.Context, recv rpccp.Message) error {
	defer c.mu.Unlock()
	c.mu.Lock()
	c.reportf("unknown message type %v from remote", recv.Which())
	err := c.sendMessage(ctx, func(msg rpccp.Message) error {
		return msg.SetUnimplemented(recv)
	})
	if err != nil {
		c.report(annotate(err).errorf("send unimplemented"))
		return nil
	}
	return nil
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
	c.openSends.Add(1)
	c.sendCond = make(chan struct{})
	c.mu.Unlock()
	msg, send, release, err := c.sender.NewMessage(ctx)
	if err != nil {
		c.mu.Lock()
		close(c.sendCond)
		c.sendCond = nil
		c.openSends.Done()
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
	s.c.openSends.Done()
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

// An exportID is an index into the exports table.
type exportID uint32

// expent is an entry in a Conn's export table.
type expent struct {
	client   *capnp.Client
	wireRefs int
}

func (c *Conn) findExport(id exportID) *expent {
	if int64(id) >= int64(len(c.exports)) {
		return nil
	}
	return c.exports[id]
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
