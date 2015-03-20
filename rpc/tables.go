package rpc

import (
	"sync"

	"golang.org/x/net/context"
	"zombiezen.com/go/capnproto"
)

// Table IDs
type (
	questionID uint32
	answerID   uint32
	exportID   uint32
	importID   uint32
	embargoID  uint32
)

type question struct {
	conn   *Conn
	method *capnp.Method // nil if this is bootstrap
	ctx    context.Context

	fulfiller
	// id should only be used if fulfiller isn't finished and while holding fulfiller's read lock.
	id questionID
}

func (q *question) PipelineCall(transform []capnp.PipelineOp, ccall *capnp.Call) capnp.Answer {
	q.init()
	q.mu.RLock()
	defer q.mu.RUnlock()

	if a := q.answer; a != nil {
		return a.PipelineCall(transform, ccall)
	}
	if transform == nil {
		transform = []capnp.PipelineOp{}
	}
	ready := make(chan capnp.Answer, 1)
	q.conn.calls <- &call{
		Call:       ccall,
		ready:      ready,
		questionID: q.id,
		transform:  transform,
	}
	return <-ready
}

// promiseInfo returns the underlying answer if q is resolved,
// otherwise a question ID that is valid until hold is closed.
// If hold is nil, then the returned ID is zero.
func (q *question) promiseInfo(hold <-chan struct{}) (ans capnp.Answer, id questionID) {
	q.mu.RLock()
	if q.answer != nil {
		ans = q.answer
		q.mu.RUnlock()
		return ans, 0
	}
	if hold == nil {
		q.mu.RUnlock()
		return nil, 0
	}
	go func() {
		<-hold
		q.mu.RUnlock()
	}()
	return nil, q.id
}

type questionTable struct {
	mu  sync.RWMutex
	tab []*question
	gen idgen
}

// new creates a new question with an unassigned ID.
func (qt *questionTable) new(conn *Conn, ctx context.Context, method *capnp.Method) *question {
	qt.mu.Lock()
	id := questionID(qt.gen.next())
	q := &question{
		id:     id,
		conn:   conn,
		ctx:    ctx,
		method: method,
	}
	q.init()
	if int(id) == len(qt.tab) {
		qt.tab = append(qt.tab, q)
	} else {
		qt.tab[id] = q
	}
	qt.mu.Unlock()
	return q
}

func (qt *questionTable) get(id questionID) *question {
	qt.mu.RLock()
	var q *question
	if int(id) < len(qt.tab) {
		q = qt.tab[id]
	}
	qt.mu.RUnlock()
	return q
}

func (qt *questionTable) remove(id questionID) bool {
	qt.mu.Lock()
	ok := int(id) < len(qt.tab) && qt.tab[id] != nil
	if ok {
		qt.tab[id] = nil
	}
	qt.gen.remove(uint32(id))
	qt.mu.Unlock()
	return ok
}

type answer struct {
	id     answerID
	cancel context.CancelFunc
	fulfiller
}

type answerTable struct {
	mu  sync.Mutex
	tab map[answerID]*answer
}

func (at *answerTable) get(id answerID) *answer {
	at.mu.Lock()
	var a *answer
	if at.tab != nil {
		a = at.tab[id]
	}
	at.mu.Unlock()
	return a
}

// insert creates a new question with the given ID, returning nil
// if the ID is already in use.
func (at *answerTable) insert(id answerID, cancel context.CancelFunc) *answer {
	at.mu.Lock()
	if at.tab == nil {
		at.tab = make(map[answerID]*answer)
	}
	var a *answer
	if _, ok := at.tab[id]; !ok {
		a = &answer{id: id, cancel: cancel}
		a.init()
		at.tab[id] = a
	}
	at.mu.Unlock()
	return a
}

func (at *answerTable) pop(id answerID) *answer {
	at.mu.Lock()
	var a *answer
	if at.tab != nil {
		a = at.tab[id]
		delete(at.tab, id)
	}
	at.mu.Unlock()
	return a
}

type export struct {
	id     exportID
	client capnp.Client
	refCount
}

type exportTable struct {
	mu  sync.RWMutex
	tab []*export
	gen idgen
}

func (et *exportTable) get(id exportID) *export {
	et.mu.RLock()
	var e *export
	if int(id) < len(et.tab) {
		e = et.tab[id]
	}
	et.mu.RUnlock()
	return e
}

// add puts client in the table with a new ID.
func (et *exportTable) add(client capnp.Client) exportID {
	et.mu.Lock()
	id := exportID(et.gen.next())
	export := &export{
		id:     id,
		client: client,
		refCount: refCount{
			refs:    1,
			release: func() {}, // TODO(light)
		},
	}
	if int(id) == len(et.tab) {
		et.tab = append(et.tab, export)
	} else {
		et.tab[id] = export
	}
	et.mu.Unlock()
	return id
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

type refCount struct {
	refs    int
	release func()
	mu      sync.Mutex
}

func (rc *refCount) inc(n int) {
	rc.mu.Lock()
	if rc.refs > 0 {
		rc.refs += n
	}
	rc.mu.Unlock()
}

func (rc *refCount) dec(n int) {
	rc.mu.Lock()
	if rc.refs > 0 {
		rc.refs -= n
		if rc.refs <= 0 {
			rc.release()
		}
	}
	rc.mu.Unlock()
}
