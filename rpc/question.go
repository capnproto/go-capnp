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

// fulfill is called to resolve a question succesfully.
// It must be called from the coordinate goroutine.
func (q *question) fulfill(obj capnp.Object) {
	q.mu.Lock()
	if q.state != questionInProgress {
		q.mu.Unlock()
		panic("question.fulfill called more than once")
	}
	q.obj, q.state = obj, questionResolved
	close(q.resolved)
	// TODO(light): return embargoes.
	q.mu.Unlock()
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
