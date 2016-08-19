package rpc

import (
	"errors"
	"sync"

	"golang.org/x/net/context"
	"zombiezen.com/go/capnproto2"
	"zombiezen.com/go/capnproto2/internal/fulfiller"
	"zombiezen.com/go/capnproto2/internal/queue"
	rpccapnp "zombiezen.com/go/capnproto2/std/capnp/rpc"
)

// callQueueSize is the maximum number of calls that can be queued per answer or client.
// TODO(light): make this a ConnOption
const callQueueSize = 64

// insertAnswer creates a new answer with the given ID, returning nil
// if the ID is already in use.
func (c *Conn) insertAnswer(id answerID, cancel context.CancelFunc) *answer {
	if c.answers == nil {
		c.answers = make(map[answerID]*answer)
	} else if _, exists := c.answers[id]; exists {
		return nil
	}
	a := &answer{
		id:          id,
		cancel:      cancel,
		manager:     &c.manager,
		out:         c.out,
		returns:     c.returns,
		queueCloses: c.queueCloses,
		resolved:    make(chan struct{}),
		queue:       make([]pcall, 0, callQueueSize),
	}
	c.answers[id] = a
	return a
}

func (c *Conn) popAnswer(id answerID) *answer {
	if c.answers == nil {
		return nil
	}
	a := c.answers[id]
	delete(c.answers, id)
	return a
}

type answer struct {
	id          answerID
	cancel      context.CancelFunc
	resultCaps  []exportID
	manager     *manager
	out         chan<- rpccapnp.Message
	returns     chan<- *outgoingReturn
	queueCloses chan<- queueClientClose
	resolved    chan struct{}

	mu    sync.RWMutex
	obj   capnp.Ptr
	err   error
	done  bool
	queue []pcall
}

// fulfill is called to resolve an answer successfully and returns a list
// of return messages to send.
// It must be called from the coordinate goroutine.
func (a *answer) fulfill(msgs []rpccapnp.Message, obj capnp.Ptr, makeCapTable capTableMaker) []rpccapnp.Message {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.done {
		panic("answer.fulfill called more than once")
	}
	a.obj, a.done = obj, true
	// TODO(light): populate resultCaps

	retmsg := newReturnMessage(nil, a.id)
	ret, _ := retmsg.Return()
	payload, _ := ret.NewResults()
	payload.SetContentPtr(obj)
	payloadTab, err := makeCapTable(ret.Segment())
	if err != nil {
		// TODO(light): handle this more gracefully
		panic(err)
	}
	payload.SetCapTable(payloadTab)
	msgs = append(msgs, retmsg)

	queues, msgs := a.emptyQueue(msgs, obj)
	ctab := obj.Segment().Message().CapTable
	for capIdx, q := range queues {
		ctab[capIdx] = newQueueClient(a.manager, ctab[capIdx], q, a.out, a.queueCloses)
	}
	close(a.resolved)
	return msgs
}

// reject is called to resolve an answer with failure and returns a list
// of return messages to send.
// It must be called from the coordinate goroutine.
func (a *answer) reject(msgs []rpccapnp.Message, err error) []rpccapnp.Message {
	if err == nil {
		panic("answer.reject called with nil")
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.done {
		panic("answer.reject called more than once")
	}
	a.err, a.done = err, true
	m := newReturnMessage(nil, a.id)
	mret, _ := m.Return()
	setReturnException(mret, err)
	msgs = append(msgs, m)
	for i := range a.queue {
		msgs = a.queue[i].a.reject(msgs, err)
		a.queue[i] = pcall{}
	}
	close(a.resolved)
	return msgs
}

// emptyQueue splits the queue by which capability it targets
// and drops any invalid calls.  Once this function returns, a.queue
// will be nil.
func (a *answer) emptyQueue(msgs []rpccapnp.Message, obj capnp.Ptr) (map[capnp.CapabilityID][]qcall, []rpccapnp.Message) {
	qs := make(map[capnp.CapabilityID][]qcall, len(a.queue))
	for i, pc := range a.queue {
		c, err := capnp.TransformPtr(obj, pc.transform)
		if err != nil {
			msgs = pc.a.reject(msgs, err)
			continue
		}
		ci := c.Interface()
		if !ci.IsValid() {
			msgs = pc.a.reject(msgs, capnp.ErrNullClient)
			continue
		}
		cn := ci.Capability()
		if qs[cn] == nil {
			qs[cn] = make([]qcall, 0, len(a.queue)-i)
		}
		qs[cn] = append(qs[cn], pc.qcall)
	}
	a.queue = nil
	return qs, msgs
}

func (a *answer) peek() (obj capnp.Ptr, err error, ok bool) {
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
	cc, err := call.Copy(nil)
	if err != nil {
		return err
	}
	a.queue = append(a.queue, pcall{
		transform: transform,
		qcall: qcall{
			a:    result,
			call: cc,
		},
	})
	return nil
}

// queueDisembargo is called from the coordinate goroutine to add a
// disembargo message to the queue.
func (a *answer) queueDisembargo(transform []capnp.PipelineOp, id embargoID, target rpccapnp.MessageTarget) (queued bool, err error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if !a.done {
		return false, errDisembargoOngoingAnswer
	}
	if a.err != nil {
		return false, errDisembargoNonImport
	}
	targetPtr, err := capnp.TransformPtr(a.obj, transform)
	if err != nil {
		return false, err
	}
	client := targetPtr.Interface().Client()
	qc, ok := client.(*queueClient)
	if !ok {
		// No need to embargo, disembargo immediately.
		return false, nil
	}
	if ic, ok := extractRPCClient(qc.client).(*importClient); !(ok && a.manager == ic.manager) {
		return false, errDisembargoNonImport
	}
	qc.mu.Lock()
	if !qc.isPassthrough() {
		err = qc.pushEmbargo(id, target)
		if err == nil {
			queued = true
		}
	}
	qc.mu.Unlock()
	return queued, err
}

func (a *answer) pipelineClient(transform []capnp.PipelineOp) capnp.Client {
	return &localAnswerClient{a: a, transform: transform}
}

// joinAnswer resolves an RPC answer by waiting on a generic answer.
// It waits until the generic answer is finished, so it should be run
// in its own goroutine.
func joinAnswer(a *answer, ca capnp.Answer) {
	s, err := ca.Struct()
	r := &outgoingReturn{
		a:   a,
		obj: s.ToPtr(),
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
func joinFulfiller(f *fulfiller.Fulfiller, ca capnp.Answer) {
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
	obj capnp.Ptr
	err error
}

type queueClient struct {
	manager *manager
	client  capnp.Client
	out     chan<- rpccapnp.Message
	closes  chan<- queueClientClose

	mu    sync.RWMutex
	q     queue.Queue
	calls qcallList
}

func newQueueClient(m *manager, client capnp.Client, queue []qcall, out chan<- rpccapnp.Message, closes chan<- queueClientClose) *queueClient {
	qc := &queueClient{
		manager: m,
		client:  client,
		out:     out,
		closes:  closes,
		calls:   make(qcallList, callQueueSize),
	}
	qc.q.Init(qc.calls, copy(qc.calls, queue))
	go qc.flushQueue()
	return qc
}

func (qc *queueClient) pushCall(cl *capnp.Call) capnp.Answer {
	f := new(fulfiller.Fulfiller)
	cl, err := cl.Copy(nil)
	if err != nil {
		return capnp.ErrorAnswer(err)
	}
	i := qc.q.Push()
	if i == -1 {
		return capnp.ErrorAnswer(errQueueFull)
	}
	qc.calls[i] = qcall{call: cl, f: f}
	return f
}

func (qc *queueClient) pushEmbargo(id embargoID, tgt rpccapnp.MessageTarget) error {
	i := qc.q.Push()
	if i == -1 {
		return errQueueFull
	}
	qc.calls[i] = qcall{embargoID: id, embargoTarget: tgt}
	return nil
}

// flushQueue is run in its own goroutine.
func (qc *queueClient) flushQueue() {
	var c qcall
	qc.mu.RLock()
	if i := qc.q.Front(); i != -1 {
		c = qc.calls[i]
	}
	qc.mu.RUnlock()
	for c.which() != qcallInvalid {
		qc.handle(&c)

		qc.mu.Lock()
		qc.q.Pop()
		if i := qc.q.Front(); i != -1 {
			c = qc.calls[i]
		} else {
			c = qcall{}
		}
		qc.mu.Unlock()
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
		msg := newDisembargoMessage(nil, rpccapnp.Disembargo_context_Which_receiverLoopback, c.embargoID)
		d, _ := msg.Disembargo()
		d.SetTarget(c.embargoTarget)
		sendMessage(qc.manager, qc.out, msg)
	}
}

func (qc *queueClient) isPassthrough() bool {
	return qc.q.Len() == 0
}

func (qc *queueClient) Call(cl *capnp.Call) capnp.Answer {
	// Fast path: queue is flushed.
	qc.mu.RLock()
	ok := qc.isPassthrough()
	qc.mu.RUnlock()
	if ok {
		return qc.client.Call(cl)
	}

	// Add to queue.
	qc.mu.Lock()
	// Since we released the lock, check that the queue hasn't been flushed.
	if qc.isPassthrough() {
		qc.mu.Unlock()
		return qc.client.Call(cl)
	}
	ans := qc.pushCall(cl)
	qc.mu.Unlock()
	return ans
}

func (qc *queueClient) WrappedClient() capnp.Client {
	qc.mu.RLock()
	ok := qc.isPassthrough()
	qc.mu.RUnlock()
	if !ok {
		return nil
	}
	return qc.client
}

func (qc *queueClient) Close() error {
	done := make(chan struct{})
	select {
	case qc.closes <- queueClientClose{qc, done}:
	case <-qc.manager.finish:
		return qc.manager.err()
	}
	select {
	case <-done:
	case <-qc.manager.finish:
		return qc.manager.err()
	}
	return qc.client.Close()
}

// rejectQueue is called from the coordinate goroutine to close out a queueClient.
func (qc *queueClient) rejectQueue(msgs []rpccapnp.Message) []rpccapnp.Message {
	qc.mu.Lock()
	for ; qc.q.Len() > 0; qc.q.Pop() {
		c := qc.calls[qc.q.Front()]
		if w := c.which(); w == qcallRemoteCall {
			msgs = c.a.reject(msgs, errQueueCallCancel)
		} else if w == qcallLocalCall {
			c.f.Reject(errQueueCallCancel)
		} else if w == qcallDisembargo {
			m := newDisembargoMessage(nil, rpccapnp.Disembargo_context_Which_receiverLoopback, c.embargoID)
			d, _ := m.Disembargo()
			d.SetTarget(c.embargoTarget)
			msgs = append(msgs, m)
		}
	}
	qc.mu.Unlock()
	return msgs
}

// queueClientClose is a message sent to the coordinate goroutine to
// handle rejecting a queue.
type queueClientClose struct {
	qc   *queueClient
	done chan<- struct{}
}

// pcall is a queued pipeline call.
type pcall struct {
	transform []capnp.PipelineOp
	qcall
}

// qcall is a queued call.
type qcall struct {
	// Calls
	a    *answer              // non-nil if remote call
	f    *fulfiller.Fulfiller // non-nil if local call
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
	switch {
	case c.a != nil:
		return qcallRemoteCall
	case c.f != nil:
		return qcallLocalCall
	case c.embargoTarget.IsValid():
		return qcallDisembargo
	default:
		return qcallInvalid
	}
}

type qcallList []qcall

func (ql qcallList) Len() int {
	return len(ql)
}

func (ql qcallList) Clear(i int) {
	ql[i] = qcall{}
}

// A localAnswerClient is used to provide a pipelined client of an answer.
type localAnswerClient struct {
	a         *answer
	transform []capnp.PipelineOp
}

func (lac *localAnswerClient) Call(call *capnp.Call) capnp.Answer {
	lac.a.mu.Lock()
	if lac.a.done {
		obj, err := lac.a.obj, lac.a.err
		lac.a.mu.Unlock()
		return clientFromResolution(lac.transform, obj, err).Call(call)
	}
	defer lac.a.mu.Unlock()
	if len(lac.a.queue) == cap(lac.a.queue) {
		return capnp.ErrorAnswer(errQueueFull)
	}
	f := new(fulfiller.Fulfiller)
	cc, err := call.Copy(nil)
	if err != nil {
		return capnp.ErrorAnswer(err)
	}
	lac.a.queue = append(lac.a.queue, pcall{
		transform: lac.transform,
		qcall: qcall{
			f:    f,
			call: cc,
		},
	})
	return f
}

func (lac *localAnswerClient) WrappedClient() capnp.Client {
	obj, err, ok := lac.a.peek()
	if !ok {
		return nil
	}
	return clientFromResolution(lac.transform, obj, err)
}

func (lac *localAnswerClient) Close() error {
	obj, err, ok := lac.a.peek()
	if !ok {
		return nil
	}
	client := clientFromResolution(lac.transform, obj, err)
	return client.Close()
}

// A capTableMaker converts the clients in a segment's message into capability descriptors.
type capTableMaker func(*capnp.Segment) (rpccapnp.CapDescriptor_List, error)

var (
	errQueueFull       = errors.New("rpc: pipeline queue full")
	errQueueCallCancel = errors.New("rpc: queued call canceled")

	errDisembargoOngoingAnswer = errors.New("rpc: disembargo attempted on in-progress answer")
	errDisembargoNonImport     = errors.New("rpc: disembargo attempted on non-import capability")
	errDisembargoMissingAnswer = errors.New("rpc: disembargo attempted on missing answer (finished too early?)")
)
