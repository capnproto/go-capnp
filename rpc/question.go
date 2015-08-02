package rpc

import (
	"sync"

	"golang.org/x/net/context"
	"zombiezen.com/go/capnproto"
)

type questionTable struct {
	tab []*question
	gen idgen

	manager *manager
	calls   chan<- *appCall
	out     chan<- outgoingMessage
}

// new creates a new question with an unassigned ID.
func (qt *questionTable) new(ctx context.Context, method *capnp.Method) *question {
	id := questionID(qt.gen.next())
	q := &question{
		ctx:      ctx,
		method:   method,
		manager:  qt.manager,
		calls:    qt.calls,
		out:      qt.out,
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
	out       chan<- outgoingMessage
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
			if q.reject(questionCanceled, q.ctx.Err()) {
				// TODO(light): timeout?
				msg := newFinishMessage(nil, q.id, true /* release */)
				select {
				case q.out <- outgoingMessage{q.manager.context(), msg}:
				case <-q.manager.finish:
				}
			}
		case <-q.manager.finish:
			q.reject(questionCanceled, q.manager.err())
		}
	}()
}

func (q *question) fulfill(obj capnp.Object) bool {
	q.mu.Lock()
	if q.state != questionInProgress {
		q.mu.Unlock()
		return false
	}
	close(q.resolved)
	q.obj = obj
	q.state = questionResolved
	// TODO(light): embargo clients and kick off loopback.
	q.mu.Unlock()
	return true
}

func (q *question) reject(state questionState, err error) bool {
	if err == nil {
		panic("question.reject called with nil")
	}
	q.mu.Lock()
	if q.state != questionInProgress {
		q.mu.Unlock()
		return false
	}
	close(q.resolved)
	q.err, q.state = err, state
	q.mu.Unlock()
	return true
}

func (q *question) peek() (id questionID, obj capnp.Object, err error, ok bool) {
	q.mu.RLock()
	id, obj, err, ok = q.id, q.obj, q.err, q.state != questionInProgress
	q.mu.RUnlock()
	return
}

func (q *question) addPromise(transform []capnp.PipelineOp) {
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
	pcall := func(obj capnp.Object, err error) capnp.Answer {
		if err != nil {
			return capnp.ErrorAnswer(err)
		}
		c := capnp.TransformObject(obj, transform).ToInterface().Client()
		if c == nil {
			return capnp.ErrorAnswer(capnp.ErrNullClient)
		}
		return c.Call(ccall)
	}
	if _, obj, err, ok := q.peek(); ok {
		return pcall(obj, err)
	}
	if transform == nil {
		transform = []capnp.PipelineOp{}
	}
	q.mu.Lock()
	if q.state != questionInProgress {
		// Answered while acquiring lock.
		obj, err := q.obj, q.err
		q.mu.Unlock()
		return pcall(obj, err)
	}
	q.addPromise(transform)
	qchan := make(chan *question, 1)
	ac := &appCall{
		Call:       ccall,
		kind:       appPipelineCall,
		qchan:      qchan,
		questionID: q.id,
		transform:  transform,
	}
	select {
	case q.calls <- ac:
	case <-ccall.Ctx.Done():
		q.mu.Unlock()
		return capnp.ErrorAnswer(ccall.Ctx.Err())
	case <-q.manager.finish:
		q.mu.Unlock()
		return capnp.ErrorAnswer(q.manager.err())
	}
	q.mu.Unlock()
	select {
	case q := <-qchan:
		return q
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
