package rpc

import (
	"sync"

	"golang.org/x/net/context"
	"zombiezen.com/go/capnproto2"
	"zombiezen.com/go/capnproto2/internal/fulfiller"
	"zombiezen.com/go/capnproto2/internal/queue"
	rpccapnp "zombiezen.com/go/capnproto2/std/capnp/rpc"
)

// newQuestion creates a new question with an unassigned ID.
func (c *Conn) newQuestion(ctx context.Context, method *capnp.Method) *question {
	id := questionID(c.questionID.next())
	q := &question{
		ctx:      ctx,
		method:   method,
		manager:  &c.manager,
		calls:    c.calls,
		cancels:  c.cancels,
		resolved: make(chan struct{}),
		id:       id,
	}
	// TODO(light): populate paramCaps
	if int(id) == len(c.questions) {
		c.questions = append(c.questions, q)
	} else {
		c.questions[id] = q
	}
	return q
}

func (c *Conn) findQuestion(id questionID) *question {
	if int(id) >= len(c.questions) {
		return nil
	}
	return c.questions[id]
}

func (c *Conn) popQuestion(id questionID) *question {
	q := c.findQuestion(id)
	if q == nil {
		return nil
	}
	c.questions[id] = nil
	c.questionID.remove(uint32(id))
	return q
}

type question struct {
	ctx       context.Context
	method    *capnp.Method // nil if this is bootstrap
	paramCaps []exportID
	calls     chan<- *appCall
	cancels   chan<- *question
	manager   *manager
	resolved  chan struct{}

	// Fields below are protected by mu.
	mu      sync.RWMutex
	id      questionID
	obj     capnp.Ptr
	err     error
	state   questionState
	derived [][]capnp.PipelineOp
}

type questionState uint8

// Question states
const (
	questionInProgress questionState = iota
	questionResolved
	questionCanceled
)

// start signals that the question has been sent.
func (q *question) start() {
	go func() {
		select {
		case <-q.resolved:
		case <-q.ctx.Done():
			select {
			case q.cancels <- q:
			case <-q.resolved:
			case <-q.manager.finish:
			}
		case <-q.manager.finish:
			// TODO(light): connection should reject all questions on shutdown.
		}
	}()
}

// fulfill is called to resolve a question successfully and returns the disembargoes.
// It must be called from the coordinate goroutine.
func (q *question) fulfill(obj capnp.Ptr, makeDisembargo func() (embargoID, embargo)) []rpccapnp.Message {
	q.mu.Lock()
	if q.state != questionInProgress {
		q.mu.Unlock()
		panic("question.fulfill called more than once")
	}
	ctab := obj.Segment().Message().CapTable
	visited := make([]bool, len(ctab))
	msgs := make([]rpccapnp.Message, 0, len(q.derived))
	for _, d := range q.derived {
		tgt, err := capnp.TransformPtr(obj, d)
		if err != nil {
			continue
		}
		in := tgt.Interface()
		if !in.IsValid() {
			continue
		}
		client := extractRPCClient(in.Client())
		if ic, ok := client.(*importClient); ok && ic.manager == q.manager {
			// Imported from remote vat.  Don't need to disembargo.
			continue
		}
		if cn := in.Capability(); !visited[cn] {
			id, e := makeDisembargo()
			ctab[cn] = newEmbargoClient(q.manager, ctab[cn], e)
			m := newDisembargoMessage(nil, rpccapnp.Disembargo_context_Which_senderLoopback, id)
			dis, _ := m.Disembargo()
			mt, _ := dis.NewTarget()
			pa, _ := mt.NewPromisedAnswer()
			pa.SetQuestionId(uint32(q.id))
			transformToPromisedAnswer(m.Segment(), pa, d)
			mt.SetPromisedAnswer(pa)
			msgs = append(msgs, m)
			visited[cn] = true
		}
	}
	q.obj, q.state = obj, questionResolved
	close(q.resolved)
	q.mu.Unlock()
	return msgs
}

// reject is called to resolve a question with failure.
// It must be called from the coordinate goroutine.
func (q *question) reject(state questionState, err error) {
	if err == nil {
		panic("question.reject called with nil")
	}
	q.mu.Lock()
	if q.state != questionInProgress {
		q.mu.Unlock()
		panic("question.reject called more than once")
	}
	q.err, q.state = err, state
	close(q.resolved)
	q.mu.Unlock()
}

func (q *question) peek() (id questionID, obj capnp.Ptr, err error, ok bool) {
	q.mu.RLock()
	id, obj, err, ok = q.id, q.obj, q.err, q.state != questionInProgress
	q.mu.RUnlock()
	return
}

func (q *question) addPromise(transform []capnp.PipelineOp) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for _, d := range q.derived {
		if transformsEqual(transform, d) {
			return
		}
	}
	q.derived = append(q.derived, transform)
}

func transformsEqual(t, u []capnp.PipelineOp) bool {
	if len(t) != len(u) {
		return false
	}
	for i := range t {
		if t[i].Field != u[i].Field {
			return false
		}
	}
	return true
}

func (q *question) Struct() (capnp.Struct, error) {
	<-q.resolved
	_, obj, err, _ := q.peek()
	return obj.Struct(), err
}

func (q *question) PipelineCall(transform []capnp.PipelineOp, ccall *capnp.Call) capnp.Answer {
	ac, achan := newAppPipelineCall(q, transform, ccall)
	select {
	case q.calls <- ac:
	case <-ccall.Ctx.Done():
		return capnp.ErrorAnswer(ccall.Ctx.Err())
	case <-q.manager.finish:
		return capnp.ErrorAnswer(q.manager.err())
	}
	select {
	case a := <-achan:
		return a
	case <-ccall.Ctx.Done():
		return capnp.ErrorAnswer(ccall.Ctx.Err())
	case <-q.manager.finish:
		return capnp.ErrorAnswer(q.manager.err())
	}
}

func (q *question) PipelineClose(transform []capnp.PipelineOp) error {
	<-q.resolved
	_, obj, err, _ := q.peek()
	if err != nil {
		return err
	}
	x, err := capnp.TransformPtr(obj, transform)
	if err != nil {
		return err
	}
	c := x.Interface().Client()
	if c == nil {
		return capnp.ErrNullClient
	}
	return c.Close()
}

// embargoClient is a client that waits until an embargo signal is
// received to deliver calls.
type embargoClient struct {
	manager *manager
	client  capnp.Client
	embargo embargo

	mu    sync.RWMutex
	q     queue.Queue
	calls ecallList
}

func newEmbargoClient(manager *manager, client capnp.Client, e embargo) *embargoClient {
	ec := &embargoClient{
		manager: manager,
		client:  client,
		embargo: e,
		calls:   make(ecallList, callQueueSize),
	}
	ec.q.Init(ec.calls, 0)
	go ec.flushQueue()
	return ec
}

func (ec *embargoClient) push(cl *capnp.Call) capnp.Answer {
	f := new(fulfiller.Fulfiller)
	cl, err := cl.Copy(nil)
	if err != nil {
		return capnp.ErrorAnswer(err)
	}
	i := ec.q.Push()
	if i == -1 {
		return capnp.ErrorAnswer(errQueueFull)
	}
	ec.calls[i] = ecall{cl, f}
	return f
}

func (ec *embargoClient) Call(cl *capnp.Call) capnp.Answer {
	// Fast path: queue is flushed.
	ec.mu.RLock()
	ok := ec.isPassthrough()
	ec.mu.RUnlock()
	if ok {
		return ec.client.Call(cl)
	}

	ec.mu.Lock()
	if ec.isPassthrough() {
		ec.mu.Unlock()
		return ec.client.Call(cl)
	}
	ans := ec.push(cl)
	ec.mu.Unlock()
	return ans
}

func (ec *embargoClient) WrappedClient() capnp.Client {
	ec.mu.RLock()
	ok := ec.isPassthrough()
	ec.mu.RUnlock()
	if !ok {
		return nil
	}
	return ec.client
}

func (ec *embargoClient) isPassthrough() bool {
	select {
	case <-ec.embargo:
	default:
		return false
	}
	return ec.q.Len() == 0
}

func (ec *embargoClient) Close() error {
	ec.mu.Lock()
	for ; ec.q.Len() > 0; ec.q.Pop() {
		c := ec.calls[ec.q.Front()]
		c.f.Reject(errQueueCallCancel)
	}
	ec.mu.Unlock()
	return ec.client.Close()
}

// flushQueue is run in its own goroutine.
func (ec *embargoClient) flushQueue() {
	select {
	case <-ec.embargo:
	case <-ec.manager.finish:
		return
	}
	var c ecall
	ec.mu.RLock()
	if i := ec.q.Front(); i != -1 {
		c = ec.calls[i]
	}
	ec.mu.RUnlock()
	for c.call != nil {
		ans := ec.client.Call(c.call)
		go joinFulfiller(c.f, ans)

		ec.mu.Lock()
		ec.q.Pop()
		if i := ec.q.Front(); i != -1 {
			c = ec.calls[i]
		} else {
			c = ecall{}
		}
		ec.mu.Unlock()
	}
}

type ecall struct {
	call *capnp.Call
	f    *fulfiller.Fulfiller
}

type ecallList []ecall

func (el ecallList) Len() int {
	return len(el)
}

func (el ecallList) Clear(i int) {
	el[i] = ecall{}
}
