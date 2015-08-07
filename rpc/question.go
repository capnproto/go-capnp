package rpc

import (
	"sync"

	"golang.org/x/net/context"
	"zombiezen.com/go/capnproto"
	"zombiezen.com/go/capnproto/internal/fulfiller"
	"zombiezen.com/go/capnproto/internal/queue"
	"zombiezen.com/go/capnproto/rpc/rpccapnp"
)

type questionTable struct {
	tab []*question
	gen idgen

	manager *manager
	calls   chan<- *appCall
	cancels chan<- *question
}

// new creates a new question with an unassigned ID.
func (qt *questionTable) new(ctx context.Context, method *capnp.Method) *question {
	id := questionID(qt.gen.next())
	q := &question{
		ctx:      ctx,
		method:   method,
		manager:  qt.manager,
		calls:    qt.calls,
		cancels:  qt.cancels,
		resolved: make(chan struct{}),
		id:       id,
	}
	// TODO(light): populate paramCaps
	if int(id) == len(qt.tab) {
		qt.tab = append(qt.tab, q)
	} else {
		qt.tab[id] = q
	}
	return q
}

func (qt *questionTable) get(id questionID) *question {
	var q *question
	if int(id) < len(qt.tab) {
		q = qt.tab[id]
	}
	return q
}

func (qt *questionTable) pop(id questionID) *question {
	var q *question
	if int(id) < len(qt.tab) {
		q = qt.tab[id]
		qt.tab[id] = nil
		qt.gen.remove(uint32(id))
	}
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
	obj     capnp.Object
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

// fulfill is called to resolve a question succesfully and returns the disembargoes.
// It must be called from the coordinate goroutine.
func (q *question) fulfill(obj capnp.Object, makeDisembargo func() (embargoID, embargo)) []rpccapnp.Message {
	q.mu.Lock()
	if q.state != questionInProgress {
		q.mu.Unlock()
		panic("question.fulfill called more than once")
	}
	ctab := obj.Segment.Message.CapTable()
	visited := make([]bool, len(ctab))
	msgs := make([]rpccapnp.Message, 0, len(q.derived))
	for _, d := range q.derived {
		in := capnp.TransformObject(obj, d)
		if in.Type() != capnp.TypeInterface {
			continue
		}
		client := extractRPCClient(in.ToInterface().Client())
		if ic, ok := client.(*importClient); ok && ic.manager == q.manager {
			// Imported from remote vat.  Don't need to disembargo.
			continue
		}
		if cn := in.ToInterface().Capability(); !visited[cn] {
			id, e := makeDisembargo()
			ctab[cn] = newEmbargoClient(q.manager, ctab[cn], e)
			m := newDisembargoMessage(nil, rpccapnp.Disembargo_context_Which_senderLoopback, id)
			mt := rpccapnp.NewMessageTarget(m.Segment)
			pa := rpccapnp.NewPromisedAnswer(m.Segment)
			pa.SetQuestionId(uint32(q.id))
			transformToPromisedAnswer(m.Segment, pa, d)
			mt.SetPromisedAnswer(pa)
			m.Disembargo().SetTarget(mt)
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

func (q *question) peek() (id questionID, obj capnp.Object, err error, ok bool) {
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
	return obj.ToStruct(), err
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
	c := capnp.TransformObject(obj, transform).ToInterface().Client()
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

	mu sync.RWMutex
	q  queue.Queue
}

func newEmbargoClient(manager *manager, client capnp.Client, e embargo) *embargoClient {
	ec := &embargoClient{
		manager: manager,
		client:  client,
		embargo: e,
	}
	ec.q.Init(make(ecallList, callQueueSize), 0)
	go ec.flushQueue()
	return ec
}

func (ec *embargoClient) push(cl *capnp.Call) capnp.Answer {
	f := new(fulfiller.Fulfiller)
	cl = cl.Copy(nil)
	if ok := ec.q.Push(ecall{cl, f}); !ok {
		return capnp.ErrorAnswer(errQueueFull)
	}
	return f
}

func (ec *embargoClient) peek() ecall {
	if ec.q.Len() == 0 {
		return ecall{}
	}
	return ec.q.Peek().(ecall)
}

func (ec *embargoClient) pop() ecall {
	if ec.q.Len() == 0 {
		return ecall{}
	}
	return ec.q.Pop().(ecall)
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
	for {
		c := ec.pop()
		if c.call == nil {
			break
		}
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
	ec.mu.RLock()
	c := ec.peek()
	ec.mu.RUnlock()
	for c.call != nil {
		ans := ec.client.Call(c.call)
		go joinFulfiller(c.f, ans)
		ec.mu.Lock()
		ec.pop()
		c = ec.peek()
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

func (el ecallList) At(i int) interface{} {
	return el[i]
}

func (el ecallList) Set(i int, x interface{}) {
	if x == nil {
		el[i] = ecall{}
	} else {
		el[i] = x.(ecall)
	}
}
