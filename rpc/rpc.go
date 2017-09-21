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
	// mu protects all fields in the Conn.  However, mu should not be held
	// while making calls that take indeterminate time (I/O or application
	// code).  Condition channels protect operations on any field that
	// take an indeterminate amount of time.  Thus, critical sections
	// involving mu are quite short, while still ensuring mutually
	// exclusive access to resources.
	mu sync.Mutex

	closed     bool
	recvDone   <-chan struct{}
	recvCloser interface {
		CloseRecv() error
	}

	// sender is protected by sendCond, a condition channel.
	// sendCond is non-nil if an operation involving sender is in
	// progress, and the channel is closed when the operation is finished.
	// Details of this are handled by acquireSender and releaseSender.
	sendCond chan struct{}
	sender   Sender

	questions  []*question
	questionID idgen
	imports    map[importID]*impent
}

// Options specifies optional parameters for creating a Conn.
type Options struct {
}

// NewConn creates a new connection that communications on a given
// transport.  Closing the connection will close the transport.
// Passing nil for opts is the same as passing the zero value.
func NewConn(t Transport, opts *Options) *Conn {
	done := make(chan struct{})
	c := &Conn{
		sender:     t,
		recvCloser: t,
		recvDone:   done,
	}
	go func() {
		defer close(done)
		c.receive(context.Background(), t)
	}()
	return c
}

// sendMessage creates a new message on the Sender, calls f to build it,
// and sends it if f does not return an error.  The caller must be
// holding onto c.mu.  However, while f is being called, it will not be
// holding onto c.mu.
func (c *Conn) sendMessage(ctx context.Context, f func(msg rpccp.Message) error) error {
	// Acquire send condition.
	for {
		if c.closed {
			// TODO(someday): classify as disconnected
			return errors.New("connection closed")
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
			return ctx.Err()
		}
		c.mu.Lock()
	}
	c.sendCond = make(chan struct{})
	c.mu.Unlock()
	defer func() {
		c.mu.Lock()
		close(c.sendCond)
		c.sendCond = nil
	}()

	// Build and send message.
	msg, send, release, err := c.sender.NewMessage(ctx)
	if err != nil {
		return fmt.Errorf("create message: %v", err)
	}
	defer release()
	if err := f(msg); err != nil {
		return fmt.Errorf("build message: %v", err)
	}
	if err := send(); err != nil {
		return fmt.Errorf("send message: %v", err)
	}
	return nil
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
	q := &question{
		id:        id,
		conn:      c,
		bootstrap: true,
	}
	if int(id) == len(c.questions) {
		c.questions = append(c.questions, q)
	} else {
		c.questions[id] = q
	}
	p := capnp.NewPromise(q)
	q.p = p
	return p.Answer().Client().AddRef()
}

// Close sends an abort to the remote vat and closes the underlying
// transport.
func (c *Conn) Close() error {
	// Mark closed and stop receiving messages.
	c.mu.Lock()
	c.closed = true
	c.mu.Unlock()
	rerr := c.recvCloser.CloseRecv() // will wait on recvDone at the end

	// Close Sender after all sends are finished.
	// TODO(soon): This assumes that there aren't messages retained past
	// the release of the send condition.
	c.mu.Lock()
	for {
		s := c.sendCond
		if s == nil {
			break
		}
		c.mu.Unlock()
		<-s
		c.mu.Lock()
	}
	c.mu.Unlock()

	// Send abort message (ignoring error).
	{
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

	<-c.recvDone
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
	c.questions[qid] = nil

	pr := c.parseReturn(ret)
	if pr.err == nil {
		c.mu.Unlock()
		q.p.Fulfill(pr.result)
		if q.bootstrap {
			q.p.ReleaseClients()
			releaseRet()
		} else {
			q.release = releaseRet
		}
		c.mu.Lock()
	} else {
		q.release = func() {}
		c.mu.Unlock()
		q.p.Reject(pr.err)
		if q.bootstrap {
			q.p.ReleaseClients()
		}
		releaseRet()
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
		return fmt.Errorf("rpc: receive return: send finish: %v", err)
	}
	c.questionID.remove(uint32(qid))
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
	ptab, err := payload.CapTable()
	if err != nil {
		return capnp.Ptr{}, fmt.Errorf("read payload: %v", err)
	}
	p, err := payload.Content()
	if err != nil {
		return capnp.Ptr{}, fmt.Errorf("read payload: %v", err)
	}
	mtab := make([]*capnp.Client, ptab.Len())
	for i := 0; i < ptab.Len(); i++ {
		mtab[i] = c.recvCap(ptab.At(i))
	}
	payload.Message().CapTable = mtab
	return p, nil
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
