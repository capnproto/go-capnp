package rpc

import (
	"errors"
	"sync"

	"golang.org/x/net/context"
	"zombiezen.com/go/capnproto"
	"zombiezen.com/go/capnproto/rpc/rpccapnp"
)

// callQueueSize is the maximum number of calls that can be queued per answer or client.
// TODO(light): make this a ConnOption
const callQueueSize = 64

type answerTable struct {
	tab     map[answerID]*answer
	manager *manager
	returns chan<- *outgoingReturn
}

func (at *answerTable) get(id answerID) *answer {
	var a *answer
	if at.tab != nil {
		a = at.tab[id]
	}
	return a
}

// insert creates a new question with the given ID, returning nil
// if the ID is already in use.
func (at *answerTable) insert(id answerID, cancel context.CancelFunc) *answer {
	if at.tab == nil {
		at.tab = make(map[answerID]*answer)
	}
	var a *answer
	if _, ok := at.tab[id]; !ok {
		a = &answer{
			id:       id,
			cancel:   cancel,
			manager:  at.manager,
			returns:  at.returns,
			resolved: make(chan struct{}),
			queue:    make([]pcall, 0, callQueueSize),
		}
		at.tab[id] = a
	}
	return a
}

func (at *answerTable) pop(id answerID) *answer {
	var a *answer
	if at.tab != nil {
		a = at.tab[id]
		delete(at.tab, id)
	}
	return a
}

type answer struct {
	id         answerID
	cancel     context.CancelFunc
	resultCaps []exportID
	manager    *manager
	returns    chan<- *outgoingReturn
	resolved   chan struct{}

	mu    sync.RWMutex
	obj   capnp.Object
	err   error
	done  bool
	queue []pcall
}

// fulfill is called to resolve an answer succesfully.
// It must be called from the coordinate goroutine.
func (a *answer) fulfill(obj capnp.Object) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.done {
		panic("answer.fulfill called more than once")
	}
	a.obj, a.done = obj, true
	// TODO(light): populate resultCaps
	queues := a.emptyQueue(obj)
	ctab := obj.Segment.Message.CapTable()
	for capIdx, q := range queues {
		ctab[capIdx] = newQueueClient(ctab[capIdx], q)
	}
	close(a.resolved)
}

// reject is called to resolve an answer with failure.
// It must be called from the coordinate goroutine.
func (a *answer) reject(err error) {
	if err == nil {
		panic("answer.reject called with nil")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.done {
		panic("answer.reject called more than once")
	}
	a.err, a.done = err, true
	for i := range a.queue {
		// TODO(light): DO NOT SUBMIT this needs returns
		a.queue[i].a.reject(err)
		a.queue[i] = pcall{}
	}
	close(a.resolved)
}

// emptyQueue splits the queue by which capability it targets
// and drops any invalid calls.  Once this function returns, a.queue
// will be nil.
func (a *answer) emptyQueue(obj capnp.Object) map[uint32][]qcall {
	qs := make(map[uint32][]qcall, len(a.queue))
	for i, pc := range a.queue {
		c := capnp.TransformObject(obj, pc.transform)
		if c.Type() != capnp.TypeInterface {
			// TODO(light): DO NOT SUBMIT this needs returns
			pc.a.reject(capnp.ErrNullClient)
			continue
		}
		cn := c.ToInterface().Capability()
		if qs[cn] == nil {
			qs[cn] = make([]qcall, 0, len(a.queue)-i)
		}
		qs[cn] = append(qs[cn], pc.qcall)
	}
	a.queue = nil
	return qs
}

func (a *answer) peek() (obj capnp.Object, err error, ok bool) {
	a.mu.RLock()
	obj, err, ok = a.obj, a.err, a.done
	a.mu.RUnlock()
	return
}

// queueCall is called from the coordinate goroutine to add a call to
// the queue.
func (a *answer) queueCall(result *answer, transform []capnp.PipelineOp, call *capnp.Call) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.done {
		panic("answer.queueCall called on resolved answer")
	}
	if len(a.queue) == cap(a.queue) {
		return errQueueFull
	}
	a.queue = append(a.queue, pcall{
		transform: transform,
		qcall: qcall{
			a:    result,
			call: call,
		},
	})
	return nil
}

// queueDisembargo is called from the coordinate goroutine to add a
// disembargo message to the queue.
func (a *answer) queueDisembargo(transform []capnp.PipelineOp, id embargoID, target rpccapnp.MessageTarget) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.done {
		panic("answer.queueDisembargo called on resolved answer")
	}
	if len(a.queue) == cap(a.queue) {
		return errQueueFull
	}
	a.queue = append(a.queue, pcall{
		transform: transform,
		qcall: qcall{
			embargoID:     id,
			embargoTarget: target,
		},
	})
	return nil
}

// joinAnswer resolves an RPC answer by waiting on a generic answer.
// It waits until the generic answer is finished, so it should be run
// in its own goroutine.
func joinAnswer(a *answer, ca capnp.Answer) {
	s, err := ca.Struct()
	r := &outgoingReturn{
		a:   a,
		obj: capnp.Object(s),
		err: err,
	}
	select {
	case a.returns <- r:
	case <-a.manager.finish:
	}
}

// joinFulfiller resolves a fulfiller by waiting on a generic answer.
// It waits until the generic answer is finished, so it should be run
// in its own goroutine.
func joinFulfiller(f *capnp.Fulfiller, ca capnp.Answer) {
	s, err := ca.Struct()
	if err != nil {
		f.Reject(err)
	} else {
		f.Fulfill(s)
	}
}

// outgoingReturn is a message sent to the coordinate goroutine to
// indicate that a call started by an answer has completed.  A simple
// message is insufficient, since the connection needs to populate the
// return message's capability table.
type outgoingReturn struct {
	a   *answer
	obj capnp.Object
	err error
}

type queueClient struct {
	client capnp.Client

	mu       sync.RWMutex
	queue    []qcall
	start, n int
}

func newQueueClient(client capnp.Client, queue []qcall) *queueClient {
	qc := &queueClient{client: client, queue: make([]qcall, callQueueSize)}
	qc.n = copy(qc.queue, queue)
	go qc.flushQueue()
	return qc
}

func (qc *queueClient) pushCall(cl *capnp.Call) capnp.Answer {
	if qc.n == len(qc.queue) {
		return capnp.ErrorAnswer(errQueueFull)
	}
	f := new(capnp.Fulfiller)
	i := (qc.start + qc.n) % len(qc.queue)
	qc.queue[i] = qcall{call: cl, f: f}
	qc.n++
	return f
}

func (qc *queueClient) pushEmbargo(id embargoID, tgt rpccapnp.MessageTarget) error {
	if qc.n == len(qc.queue) {
		return errQueueFull
	}
	i := (qc.start + qc.n) % len(qc.queue)
	qc.queue[i] = qcall{embargoID: id, embargoTarget: tgt}
	qc.n++
	return nil
}

func (qc *queueClient) pop() qcall {
	if qc.n == 0 {
		return qcall{}
	}
	c := qc.queue[qc.start]
	qc.queue[qc.start] = qcall{}
	qc.start = (qc.start + 1) % len(qc.queue)
	qc.n--
	return c
}

// flushQueue is run in its own goroutine.
func (qc *queueClient) flushQueue() {
	for {
		qc.mu.Lock()
		c := qc.pop()
		qc.mu.Unlock()
		if c.which() == qcallInvalid {
			return
		}
		qc.handle(&c)
	}
}

func (qc *queueClient) handle(c *qcall) {
	switch c.which() {
	case qcallRemoteCall:
		answer := qc.client.Call(c.call)
		go joinAnswer(c.a, answer)
	case qcallLocalCall:
		answer := qc.client.Call(c.call)
		go joinFulfiller(c.f, answer)
	case qcallDisembargo:
		// TODO(light): start disembargo
	}
}

func (qc *queueClient) Call(cl *capnp.Call) capnp.Answer {
	// Fast path: queue is flushed.
	qc.mu.RLock()
	n := qc.n
	qc.mu.RUnlock()
	if n == 0 {
		return qc.client.Call(cl)
	}

	// Add to queue.
	qc.mu.Lock()
	// Since we released the lock, check that the queue hasn't been flushed.
	if qc.n == 0 {
		qc.mu.Unlock()
		return qc.client.Call(cl)
	}
	ans := qc.pushCall(cl)
	qc.mu.Unlock()
	return ans
}

func (qc *queueClient) WrappedClient() capnp.Client {
	qc.mu.RLock()
	ok := qc.n == 0
	qc.mu.RUnlock()
	if !ok {
		return nil
	}
	return qc.client
}

func (qc *queueClient) Close() error {
	qc.mu.Lock()
	// reject all queued calls
	for {
		c := qc.pop()
		if w := c.which(); w == qcallRemoteCall {
			// TODO(light): DO NOT SUBMIT this needs returns
			c.a.reject(errQueueCallCancel)
		} else if w == qcallLocalCall {
			c.f.Reject(errQueueCallCancel)
		} else if w == qcallDisembargo {
			// TODO(light): close disembargo?
		} else {
			break
		}
	}
	qc.mu.Unlock()
	return qc.client.Close()
}

// pcall is a queued pipeline call.
type pcall struct {
	transform []capnp.PipelineOp
	qcall
}

// qcall is a queued call.
type qcall struct {
	// Calls
	a    *answer          // non-nil if remote call
	f    *capnp.Fulfiller // non-nil if local call
	call *capnp.Call

	// Disembargo
	embargoID     embargoID
	embargoTarget rpccapnp.MessageTarget
}

// Queued call types.
const (
	qcallInvalid = iota
	qcallRemoteCall
	qcallLocalCall
	qcallDisembargo
)

func (c *qcall) which() int {
	if c.a != nil {
		return qcallRemoteCall
	} else if c.f != nil {
		return qcallLocalCall
	} else if capnp.Object(c.embargoTarget).Type() != capnp.TypeNull {
		return qcallDisembargo
	} else {
		return qcallInvalid
	}
}

var (
	errQueueFull       = errors.New("rpc: pipeline queue full")
	errQueueCallCancel = errors.New("rpc: queued call canceled")
)
