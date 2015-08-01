package rpc

import (
	"sync"

	"golang.org/x/net/context"
	"zombiezen.com/go/capnproto"
)

type questionTable struct {
	tab []*question
	gen idgen
}

// new creates a new question with an unassigned ID.
func (qt *questionTable) new(conn *Conn, ctx context.Context, method *capnp.Method) *question {
	id := questionID(qt.gen.next())
	q := &question{
		id:       id,
		conn:     conn,
		ctx:      ctx,
		method:   method,
		resolved: make(chan struct{}),
	}
	// TODO(light): populate paramCaps
	if int(id) == len(qt.tab) {
		qt.tab = append(qt.tab, q)
	} else {
		qt.tab[id] = q
	}
	if done := q.ctx.Done(); done != nil {
		go func() {
			select {
			case <-done:
				conn.cancelQuestion(q)
			case <-q.resolved:
			}
		}()
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
	conn      *Conn
	method    *capnp.Method // nil if this is bootstrap
	ctx       context.Context
	paramCaps []exportID
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

func (q *question) peekState() questionState {
	q.mu.RLock()
	state := q.state
	q.mu.RUnlock()
	return state
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
	var (
		obj  capnp.Object
		err  error
		sent capnp.Answer
	)
	err = q.conn.do(ccall.Ctx, func() error {
		q.mu.Lock()
		if q.state != questionInProgress {
			// Answered while scheduling.
			// Don't block tasks on sending a pipeline call.
			obj, err = q.obj, q.err
			q.mu.Unlock()
			return nil
		}
		q.addPromise(transform)
		q.mu.Unlock()
		sent = q.conn.sendCall(&call{
			Call:       ccall,
			questionID: q.id,
			transform:  transform,
		})
		return nil
	})
	if err != nil {
		return capnp.ErrorAnswer(err)
	}
	if sent == nil {
		return pcall(obj, err)
	}
	return sent
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
