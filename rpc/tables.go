package rpc

import (
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

type importTable struct {
	tab      map[importID]*impent
	manager  *manager
	calls    chan<- *appCall
	releases chan<- *outgoingRelease
}

// addRef increases the counter of the times the import ID was sent to this vat.
func (it *importTable) addRef(id importID) capnp.Client {
	if it.tab == nil {
		it.tab = make(map[importID]*impent)
	}
	ent := it.tab[id]
	var ref capnp.Client
	if ent == nil {
		client := &importClient{
			id:       id,
			manager:  it.manager,
			calls:    it.calls,
			releases: it.releases,
		}
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

// An outgoingRelease is a message sent to the coordinate goroutine to
// indicate that an import should be released.
type outgoingRelease struct {
	id    importID
	echan chan<- error
}

// An importClient implements capnp.Client for a remote capability.
type importClient struct {
	id       importID
	manager  *manager
	calls    chan<- *appCall
	releases chan<- *outgoingRelease
}

func (ic *importClient) Call(cl *capnp.Call) capnp.Answer {
	// TODO(light): don't send if closed.
	ac, achan := newAppImportCall(ic.id, cl)
	select {
	case ic.calls <- ac:
		select {
		case a := <-achan:
			return a
		case <-ic.manager.finish:
			return capnp.ErrorAnswer(ic.manager.err())
		}
	case <-ic.manager.finish:
		return capnp.ErrorAnswer(ic.manager.err())
	}
}

func (ic *importClient) Close() error {
	echan := make(chan error, 1)
	r := &outgoingRelease{
		id:    ic.id,
		echan: echan,
	}
	select {
	case ic.releases <- r:
		select {
		case err := <-echan:
			return err
		case <-ic.manager.finish:
			return ic.manager.err()
		}
	case <-ic.manager.finish:
		return ic.manager.err()
	}
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

func (et *exportTable) releaseAll() {
	for id, e := range et.tab {
		if e == nil {
			continue
		}
		if err := e.client.Close(); err != nil {
			log.Printf("rpc: export %v close: %v", id, err)
		}
		et.tab[id] = nil
		et.gen.remove(uint32(id))
	}
}

// releaseList decrements the reference count of each of the given exports by 1.
func (et *exportTable) releaseList(ids []exportID) {
	for _, id := range ids {
		et.release(id, 1)
	}
}

type embargoTable struct {
	tab []chan<- struct{}
	gen idgen
}

type embargo <-chan struct{}

func (et *embargoTable) new() (embargoID, embargo) {
	id := embargoID(et.gen.next())
	e := make(chan struct{})
	if int(id) == len(et.tab) {
		et.tab = append(et.tab, e)
	} else {
		et.tab[id] = e
	}
	return id, e
}

func (et *embargoTable) disembargo(id embargoID) {
	if int(id) >= len(et.tab) {
		return
	}
	e := et.tab[id]
	if e == nil {
		return
	}
	close(e)
	et.tab[id] = nil
	et.gen.remove(uint32(id))
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
