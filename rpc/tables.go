package rpc

import (
	"log"

	"golang.org/x/net/context"
	"zombiezen.com/go/capnproto"
	"zombiezen.com/go/capnproto/rpc/internal/refcount"
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
	conn      *Conn
	method    *capnp.Method // nil if this is bootstrap
	ctx       context.Context
	paramCaps []exportID

	canceled bool // only accessed from tasks goroutine

	fulfiller
	id questionID // protected by fulfiller.mu
}

func (q *question) PipelineCall(transform []capnp.PipelineOp, ccall *capnp.Call) capnp.Answer {
	q.init()
	if a := q.peek(); a != nil {
		return a.PipelineCall(transform, ccall)
	}
	if transform == nil {
		transform = []capnp.PipelineOp{}
	}
	var qa, sent capnp.Answer
	// TODO(light): queue without blocking?
	err := q.conn.do(ccall.Ctx, func() error {
		var id questionID
		qa, id = q.promiseInfo()
		if qa != nil {
			// Answered while scheduling.
			// Don't block tasks on sending a pipeline call.
			return nil
		}
		sent = q.conn.sendCall(&call{
			Call:       ccall,
			questionID: id,
			transform:  transform,
		})
		return nil
	})
	if err != nil {
		return capnp.ErrorAnswer(err)
	}
	if qa != nil {
		return qa.PipelineCall(transform, ccall)
	}
	return sent
}

// promiseInfo returns the underlying answer if q is resolved,
// otherwise a question ID.
func (q *question) promiseInfo() (ans capnp.Answer, id questionID) {
	q.mu.RLock()
	ans, id = q.answer, q.id
	q.mu.RUnlock()
	if ans != nil {
		id = 0
	}
	return ans, id
}

type questionTable struct {
	tab []*question
	gen idgen
}

// new creates a new question with an unassigned ID.
func (qt *questionTable) new(conn *Conn, ctx context.Context, method *capnp.Method) *question {
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
	if done := q.ctx.Done(); done != nil {
		go func() {
			select {
			case <-done:
				conn.sendCancel(q)
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

type answer struct {
	id         answerID
	cancel     context.CancelFunc
	resultCaps []exportID
	fulfiller
}

type answerTable struct {
	tab map[answerID]*answer
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
		a = &answer{id: id, cancel: cancel}
		a.init()
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

// impent is an entry in the import table.
type impent struct {
	rc   *refcount.RefCount
	refs int
}

type importTable struct {
	tab map[importID]*impent
}

// addRef increases the counter of the times the import ID was sent to this vat.
func (it *importTable) addRef(c *Conn, id importID) capnp.Client {
	if it.tab == nil {
		it.tab = make(map[importID]*impent)
	}
	ent := it.tab[id]
	var ref capnp.Client
	if ent == nil {
		client := &importClient{c: c, id: id}
		var rc *refcount.RefCount
		rc, ref = refcount.New(client)
		ent = &impent{rc: rc, refs: 0}
		it.tab[id] = ent
	}
	if ref == nil {
		ref = ent.rc.Ref()
	}
	ent.refs++
	return ref
}

// pop removes the import ID and returns the number of times the import ID was sent to this vat.
func (it *importTable) pop(id importID) (refs int) {
	if it.tab != nil {
		if ent := it.tab[id]; ent != nil {
			refs = ent.refs
		}
		delete(it.tab, id)
	}
	return
}

type export struct {
	id     exportID
	client capnp.Client

	// for use by the table only
	refs int
}

type exportTable struct {
	tab []*export
	gen idgen
}

func (et *exportTable) get(id exportID) *export {
	var e *export
	if int(id) < len(et.tab) {
		e = et.tab[id]
	}
	return e
}

// add ensures that the client is present in the table, returning its ID.
// If the client is already in the table, the previous ID is returned.
func (et *exportTable) add(client capnp.Client) exportID {
	for i, e := range et.tab {
		if e != nil && e.client == client {
			e.refs++
			return exportID(i)
		}
	}
	id := exportID(et.gen.next())
	export := &export{
		id:     id,
		client: client,
		refs:   1,
	}
	if int(id) == len(et.tab) {
		et.tab = append(et.tab, export)
	} else {
		et.tab[id] = export
	}
	return id
}

func (et *exportTable) release(id exportID, refs int) {
	if int(id) >= len(et.tab) {
		return
	}
	e := et.tab[id]
	if e == nil {
		return
	}
	e.refs -= refs
	if e.refs > 0 {
		return
	}
	if e.refs < 0 {
		log.Printf("rpc: warning: export %v has negative refcount (%d)", id, e.refs)
	}
	if err := e.client.Close(); err != nil {
		log.Printf("rpc: export %v close: %v", id, err)
	}
	et.tab[id] = nil
	et.gen.remove(uint32(id))
}

// releaseList decrements the reference count of each of the given exports by 1.
func (et *exportTable) releaseList(ids []exportID) {
	for _, id := range ids {
		et.release(id, 1)
	}
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
