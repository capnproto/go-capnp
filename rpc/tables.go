package rpc

import (
	"errors"
	"log"

	"zombiezen.com/go/capnproto2"
	"zombiezen.com/go/capnproto2/rpc/internal/refcount"
)

// Table IDs
type (
	questionID uint32
	answerID   uint32
	exportID   uint32
	importID   uint32
	embargoID  uint32
)

// impent is an entry in the import table.
type impent struct {
	rc   *refcount.RefCount
	refs int
}

// addImport increases the counter of the times the import ID was sent to this vat.
func (c *Conn) addImport(id importID) capnp.Client {
	if c.imports == nil {
		c.imports = make(map[importID]*impent)
	}
	ent := c.imports[id]
	var ref capnp.Client
	if ent == nil {
		client := &importClient{
			id:   id,
			conn: c,
		}
		var rc *refcount.RefCount
		rc, ref = refcount.New(client)
		ent = &impent{rc: rc, refs: 0}
		c.imports[id] = ent
	}
	if ref == nil {
		ref = ent.rc.Ref()
	}
	ent.refs++
	return ref
}

// popImport removes the import ID and returns the number of times the import ID was sent to this vat.
func (c *Conn) popImport(id importID) (refs int) {
	if c.imports == nil {
		return 0
	}
	ent := c.imports[id]
	if ent == nil {
		return 0
	}
	refs = ent.refs
	delete(c.imports, id)
	return refs
}

// An importClient implements capnp.Client for a remote capability.
type importClient struct {
	id     importID
	conn   *Conn
	closed bool // protected by conn.mu
}

func (ic *importClient) Call(cl *capnp.Call) capnp.Answer {
	select {
	case <-ic.conn.mu:
	case <-cl.Ctx.Done():
		return capnp.ErrorAnswer(cl.Ctx.Err())
	case <-ic.conn.manager.finish:
		return capnp.ErrorAnswer(ic.conn.manager.err())
	}
	ans := ic.lockedCall(cl)
	ic.conn.mu.Unlock()
	return ans
}

// lockedCall is equivalent to Call but assumes that the caller is
// already holding onto ic.conn.mu.
func (ic *importClient) lockedCall(cl *capnp.Call) capnp.Answer {
	if ic.closed {
		return capnp.ErrorAnswer(errImportClosed)
	}

	q := ic.conn.newQuestion(cl.Ctx, &cl.Method)
	msg := newMessage(nil)
	msgCall, _ := msg.NewCall()
	msgCall.SetQuestionId(uint32(q.id))
	msgCall.SetInterfaceId(cl.Method.InterfaceID)
	msgCall.SetMethodId(cl.Method.MethodID)
	target, _ := msgCall.NewTarget()
	target.SetImportedCap(uint32(ic.id))
	payload, _ := msgCall.NewParams()
	if err := ic.conn.fillParams(payload, cl); err != nil {
		ic.conn.popQuestion(q.id)
		return capnp.ErrorAnswer(err)
	}

	select {
	case ic.conn.out <- msg:
	case <-cl.Ctx.Done():
		ic.conn.popQuestion(q.id)
		return capnp.ErrorAnswer(cl.Ctx.Err())
	case <-ic.conn.manager.finish:
		ic.conn.popQuestion(q.id)
		return capnp.ErrorAnswer(ic.conn.manager.err())
	}
	q.start()
	return q
}

func (ic *importClient) Close() error {
	ic.conn.mu.Lock()
	closed := ic.closed
	var i int
	if !closed {
		i = ic.conn.popImport(ic.id)
		ic.closed = true
	}
	ic.conn.mu.Unlock()

	if closed {
		return errImportClosed
	}
	if i == 0 {
		return nil
	}
	msg := newMessage(nil)
	mr, err := msg.NewRelease()
	if err != nil {
		return err
	}
	mr.SetId(uint32(ic.id))
	mr.SetReferenceCount(uint32(i))
	select {
	case ic.conn.out <- msg:
		return nil
	case <-ic.conn.manager.finish:
		return ic.conn.manager.err()
	}
}

type export struct {
	id     exportID
	client capnp.Client

	// for use by the table only
	refs int
}

func (c *Conn) findExport(id exportID) *export {
	if int(id) >= len(c.exports) {
		return nil
	}
	return c.exports[id]
}

// addExport ensures that the client is present in the table, returning its ID.
// If the client is already in the table, the previous ID is returned.
func (c *Conn) addExport(client capnp.Client) exportID {
	for i, e := range c.exports {
		if e != nil && e.client == client {
			e.refs++
			return exportID(i)
		}
	}
	id := exportID(c.exportID.next())
	export := &export{
		id:     id,
		client: client,
		refs:   1,
	}
	if int(id) == len(c.exports) {
		c.exports = append(c.exports, export)
	} else {
		c.exports[id] = export
	}
	return id
}

func (c *Conn) releaseExport(id exportID, refs int) {
	e := c.findExport(id)
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
	c.exports[id] = nil
	c.exportID.remove(uint32(id))
}

func (c *Conn) releaseAllExports() {
	for id, e := range c.exports {
		if e == nil {
			continue
		}
		if err := e.client.Close(); err != nil {
			log.Printf("rpc: export %v close: %v", id, err)
		}
		c.exports[id] = nil
		c.exportID.remove(uint32(id))
	}
}

type embargo <-chan struct{}

func (c *Conn) newEmbargo() (embargoID, embargo) {
	id := embargoID(c.embargoID.next())
	e := make(chan struct{})
	if int(id) == len(c.embargoes) {
		c.embargoes = append(c.embargoes, e)
	} else {
		c.embargoes[id] = e
	}
	return id, e
}

func (c *Conn) disembargo(id embargoID) {
	if int(id) >= len(c.embargoes) {
		return
	}
	e := c.embargoes[id]
	if e == nil {
		return
	}
	close(e)
	c.embargoes[id] = nil
	c.embargoID.remove(uint32(id))
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

var errImportClosed = errors.New("rpc: call on closed import")
