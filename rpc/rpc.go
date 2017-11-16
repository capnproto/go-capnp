// Package rpc implements the Cap'n Proto RPC protocol.
package rpc // import "zombiezen.com/go/capnproto2/rpc"

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"zombiezen.com/go/capnproto2"
	rpccp "zombiezen.com/go/capnproto2/std/capnp/rpc"
)

// A Conn is a connection to another Cap'n Proto vat.
// It is safe to use from multiple goroutines.
type Conn struct {
	bootstrap *capnp.Client

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
}

// NewConn creates a new connection that communications on a given
// transport.  Closing the connection will close the transport.
// Passing nil for opts is the same as passing the zero value.
func NewConn(t Transport, opts *Options) *Conn {
	bgctx, bgcancel := context.WithCancel(context.Background())
	c := &Conn{
		sender:     t,
		recvCloser: t,
		bgctx:      bgctx,
		bgcancel:   bgcancel,
	}
	if opts != nil {
		c.bootstrap = opts.BootstrapClient
	}
	c.runBackground(func(ctx context.Context) {
		c.receive(ctx, t)
	})
	return c
}

func (c *Conn) runBackground(f func(ctx context.Context)) {
	c.bgtasks.Add(1)
	go func() {
		defer c.bgtasks.Done()
		f(c.bgctx)
	}()
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
		return capnp.ErrorClient(fmt.Errorf("rpc bootstrap: %v", err))
	}
	ctx, cancel := context.WithCancel(ctx)
	q := c.newQuestion(ctx, id, true)
	// TODO(soon): tie finishing bootstrap to resolving the client
	return capnp.NewClient(bootstrapClient{
		c:      q.p.Answer().Client().AddRef(),
		cancel: cancel,
	})
}

// newQuestion adds a new question to c's table.  The caller must be
// holding onto c.mu.
func (c *Conn) newQuestion(ctx context.Context, id questionID, bootstrap bool) *question {
	ctx, cancel := context.WithCancel(ctx)
	q := &question{
		id:        id,
		conn:      c,
		done:      cancel,
		bootstrap: bootstrap,
	}
	q.p = capnp.NewPromise(q)
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
			// TODO(soon): log error
			c.sendMessage(bgctx, func(msg rpccp.Message) error {
				fin, err := msg.NewFinish()
				if err != nil {
					return err
				}
				fin.SetQuestionId(uint32(q.id))
				fin.SetReleaseResultCaps(true)
				return nil
			})
		}
		c.mu.Unlock()
		q.p.Reject(rejectErr)
		if q.bootstrap {
			q.p.ReleaseClients()
		}
	})
	return q
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
	// Mark closed and stop receiving messages.
	c.mu.Lock()
	c.bgcancel()
	c.mu.Unlock()
	rerr := c.recvCloser.CloseRecv()

	// Wait for all other sends to finish.
	c.openSends.Wait()

	// Send abort message (ignoring error).  No locking needed, since
	// c.startSend will always return errors after closing bgctx.
	{
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
		abort.SetType(rpccp.Exception_Type_failed)
		if err := abort.SetReason("connection closed"); err != nil {
			release()
			goto closeSend
		}
		send()
		release()
	}
closeSend:
	serr := c.sender.CloseSend()

	c.bgtasks.Wait()

	// Release exported clients.
	c.mu.Lock()
	c.imports = nil
	exports := c.exports
	c.exports = nil
	c.mu.Unlock()
	c.bootstrap.Release()
	c.bootstrap = nil
	for _, e := range exports {
		if e != nil {
			e.client.Release()
		}
	}

	if rerr != nil {
		return fmt.Errorf("rpc: close transport: %v", rerr)
	}
	if serr != nil {
		return fmt.Errorf("rpc: close transport: %v", serr)
	}
	return nil
}

// receive receives and dispatches messages coming from r.
// It is intended to run in its own goroutine.
func (c *Conn) receive(ctx context.Context, r Receiver) {
	for {
		recv, releaseRecv, err := r.RecvMessage(ctx)
		if err != nil {
			// TODO(soon): log error
			return
		}
		switch recv.Which() {
		case rpccp.Message_Which_unimplemented:
			// no-op for now to avoid feedback loop
		case rpccp.Message_Which_bootstrap:
			bootstrap, err := recv.Bootstrap()
			if err != nil {
				// TODO(soon): log error
				continue
			}
			qid := answerID(bootstrap.QuestionId())
			releaseRecv()
			if err := c.handleBootstrap(ctx, qid); err != nil {
				// TODO(soon): log error
				continue
			}
		case rpccp.Message_Which_call:
			call, err := recv.Call()
			if err != nil {
				// TODO(soon): log error
				continue
			}
			err = c.handleCall(ctx, call, releaseRecv)
			if err != nil {
				// TODO(soon): log error
				continue
			}
		case rpccp.Message_Which_return:
			ret, err := recv.Return()
			if err != nil {
				// TODO(soon): log error
				continue
			}
			err = c.handleReturn(ctx, ret, releaseRecv)
			if err != nil {
				// TODO(soon): log error
				continue
			}
		case rpccp.Message_Which_finish:
			fin, err := recv.Finish()
			if err != nil {
				// TODO(soon): log error
				continue
			}
			qid := answerID(fin.QuestionId())
			releaseResultCaps := fin.ReleaseResultCaps()
			releaseRecv()
			err = c.handleFinish(ctx, qid, releaseResultCaps)
			if err != nil {
				// TODO(soon): log error
				continue
			}
		case rpccp.Message_Which_release:
			rel, err := recv.Release()
			if err != nil {
				// TODO(soon): log error
				continue
			}
			id := exportID(rel.Id())
			count := rel.ReferenceCount()
			releaseRecv()
			if err := c.handleRelease(ctx, id, count); err != nil {
				// TODO(soon): log error
				continue
			}
		default:
			err := c.handleUnimplemented(ctx, recv)
			releaseRecv()
			if err != nil {
				// TODO(soon): log error
				continue
			}
		}
	}
}

func (c *Conn) handleBootstrap(ctx context.Context, id answerID) error {
	defer c.mu.Unlock()
	c.mu.Lock()
	ans, err := c.newAnswer(ctx, id, func() {})
	if err != nil {
		return err
	}
	if c.bootstrap.IsValid() {
		err := ans.setBootstrap(c.bootstrap.AddRef())
		ans.lockedReturn(err)
	} else {
		ans.lockedReturn(errors.New("vat does not expose a public/bootstrap interface"))
	}
	return nil
}

func (c *Conn) handleCall(ctx context.Context, call rpccp.Call, releaseCall capnp.ReleaseFunc) error {
	id := answerID(call.QuestionId())
	if call.SendResultsTo().Which() != rpccp.Call_sendResultsTo_Which_caller {
		// TODO(someday): handle yourself
		releaseCall()
		// TODO(someday): classify as unimplemented
		return errors.New("call results destination other than caller unimplemented")
	}
	tgt, err := call.Target()
	if err != nil {
		releaseCall()
		return fmt.Errorf("read target: %v", err)
	}
	argsPayload, err := call.Params()
	if err != nil {
		releaseCall()
		return fmt.Errorf("read params: %v", err)
	}

	c.mu.Lock()
	callCtx, cancel := context.WithCancel(c.bgctx)
	ans, err := c.newAnswer(ctx, id, cancel)
	if err != nil {
		c.mu.Unlock()
		return err
	}
	args, err := c.recvPayload(argsPayload)
	if err != nil {
		ans.lockedReturn(fmt.Errorf("read params: %v", err))
		c.mu.Unlock()
		return nil
	}
	c.mu.Unlock()
	// TODO(soon): store PipelineCaller in answer.
	c.targetClient(tgt).RecvCall(callCtx, capnp.Recv{
		Args: args.Struct(),
		Method: capnp.Method{
			InterfaceID: call.InterfaceId(),
			MethodID:    call.MethodId(),
		},
		ReleaseArgs: releaseCall,
		Returner:    ans,
	})
	return nil
}

func (c *Conn) handleReturn(ctx context.Context, ret rpccp.Return, releaseRet capnp.ReleaseFunc) error {
	c.mu.Lock()
	// TODO(soon): disconnect if return ID not in questions table.
	qid := questionID(ret.AnswerId())
	if uint32(qid) >= uint32(len(c.questions)) {
		c.mu.Unlock()
		releaseRet()
		return fmt.Errorf("rpc: receive return: question %d does not exist", qid)
	}
	q := c.questions[qid]
	if q == nil {
		c.mu.Unlock()
		releaseRet()
		return fmt.Errorf("rpc: receive return: question %d does not exist", qid)
	}
	defer c.mu.Unlock()
	sentFinish := q.sentFinish()
	q.state |= 3 // return received and sending finish
	c.questions[qid] = nil

	pr := c.parseReturn(ret)
	if pr.err == nil {
		if q.bootstrap {
			q.release = func() {}
		} else {
			q.release = releaseRet
		}
		c.mu.Unlock()
		q.p.Fulfill(pr.result)
		if q.bootstrap {
			q.p.ReleaseClients()
			releaseRet()
		}
		q.done()
		c.mu.Lock()
	} else {
		q.release = func() {}
		c.mu.Unlock()
		q.p.Reject(pr.err)
		if q.bootstrap {
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
			return fmt.Errorf("rpc: receive return: send finish: %v", err)
		}
	} else {
		c.questionID.remove(uint32(qid))
	}
	if pr.parseFailed {
		// TODO(soon): remove stutter of "rpc:" prefix
		return fmt.Errorf("rpc: receive return: %v", pr.err)
	}
	return nil
}

func (c *Conn) parseReturn(ret rpccp.Return) parsedReturn {
	switch ret.Which() {
	case rpccp.Return_Which_results:
		r, err := ret.Results()
		if err != nil {
			return parsedReturn{err: fmt.Errorf("rpc: parse return: %v", err), parseFailed: true}
		}
		content, err := c.recvPayload(r)
		if err != nil {
			return parsedReturn{err: fmt.Errorf("rpc: parse return: %v", err), parseFailed: true}
		}
		return parsedReturn{result: content}
	case rpccp.Return_Which_exception:
		exc, err := ret.Exception()
		if err != nil {
			return parsedReturn{err: fmt.Errorf("rpc: parse return: %v", err), parseFailed: true}
		}
		reason, err := exc.Reason()
		if err != nil {
			return parsedReturn{err: fmt.Errorf("rpc: parse return: %v", err), parseFailed: true}
		}
		return parsedReturn{err: errors.New(reason)}
	default:
		w := ret.Which()
		// TODO(someday): send unimplemented message back to remote
		return parsedReturn{err: fmt.Errorf("rpc: parse return: unhandled type %v", w), parseFailed: true}
	}
}

type parsedReturn struct {
	result      capnp.Ptr
	err         error
	parseFailed bool
}

func (c *Conn) handleFinish(ctx context.Context, id answerID, releaseResultCaps bool) error {
	defer c.mu.Unlock()
	c.mu.Lock()
	ans := c.answers[id]
	if ans == nil {
		return fmt.Errorf("finish sent for unknown answer ID %d", id)
	}
	ans.state |= 2
	ans.cancel()
	if !ans.isDone() {
		// Not returned yet.
		// TODO(soon): record releaseResultCaps
		return nil
	}
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
		return fmt.Errorf("payload capability table: %v", err)
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
		return capnp.ErrorClient(fmt.Errorf("rpc: unknown CapDescriptor type %v", d.Which()))
	}
}

// recvPayload extracts the content pointer after populating the
// message's capability table.  The caller must be holding onto c.mu.
func (c *Conn) recvPayload(payload rpccp.Payload) (capnp.Ptr, error) {
	if payload.Message().CapTable != nil {
		// RecvMessage likely violated its invariant.
		return capnp.Ptr{}, errors.New("read payload: capability table already populated")
	}
	p, err := payload.Content()
	if err != nil {
		return capnp.Ptr{}, fmt.Errorf("read payload: %v", err)
	}
	ptab, err := payload.CapTable()
	if err != nil {
		// Don't allow unreadable capability table to stop other results,
		// just present an empty capability table.
		// TODO(soon): log errors
		return p, nil
	}
	mtab := make([]*capnp.Client, ptab.Len())
	for i := 0; i < ptab.Len(); i++ {
		mtab[i] = c.recvCap(ptab.At(i))
	}
	payload.Message().CapTable = mtab
	return p, nil
}

// targetClient resolves a message target into a client.  Any error
// encountered during resolution results in an error client.  The caller
// must be holding onto c.mu.
func (c *Conn) targetClient(tgt rpccp.MessageTarget) *capnp.Client {
	switch tgt.Which() {
	case rpccp.MessageTarget_Which_importedCap:
		id := exportID(tgt.ImportedCap())
		ent := c.findExport(id)
		if ent == nil {
			return capnp.ErrorClient(fmt.Errorf("unknown export ID %d", id))
		}
		return ent.client
	default:
		return capnp.ErrorClient(fmt.Errorf("unhandled message target %v", tgt.Which()))
	}
}

func (c *Conn) handleRelease(ctx context.Context, id exportID, count uint32) error {
	c.mu.Lock()
	ent := c.findExport(id)
	if ent == nil {
		c.mu.Unlock()
		return fmt.Errorf("unknown export ID %d", id)
	}
	ent.wireRefs -= int(count)
	// TODO(soon): log c.exports[id].wireRefs < 0
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

func (c *Conn) handleUnimplemented(ctx context.Context, recv rpccp.Message) error {
	defer c.mu.Unlock()
	c.mu.Lock()
	err := c.sendMessage(ctx, func(msg rpccp.Message) error {
		return msg.SetUnimplemented(recv)
	})
	if err != nil {
		return fmt.Errorf("rpc: send unimplemented: %v", err)
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
			// TODO(someday): classify as disconnected
			return sendSession{}, errors.New("connection closed")
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
			// TODO(someday): classify as disconnected
			return sendSession{}, errors.New("connection closed")
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
		return sendSession{}, fmt.Errorf("create message: %v", err)
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
		return fmt.Errorf("build message: %v", err)
	}
	if err := s.send(); err != nil {
		return fmt.Errorf("send message: %v", err)
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
