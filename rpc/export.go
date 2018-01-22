package rpc

import (
	"context"

	"zombiezen.com/go/capnproto2"
	rpccp "zombiezen.com/go/capnproto2/std/capnp/rpc"
)

// An exportID is an index into the exports table.
type exportID uint32

// expent is an entry in a Conn's export table.
type expent struct {
	client   *capnp.Client
	wireRefs uint32
}

// findExport returns the export entry with the given ID or nil if
// couldn't be found.
func (c *Conn) findExport(id exportID) *expent {
	if int64(id) >= int64(len(c.exports)) {
		return nil
	}
	return c.exports[id] // might be nil
}

// releaseExport decreases the number of wire references to an export
// by a given number.  If the export's reference count reaches zero,
// then releaseExport will pop export from the table and return the
// export's client.  The caller must be holding onto c.mu, and the
// caller is responsible for releasing the client once the caller is no
// longer holding onto c.mu.
func (c *Conn) releaseExport(id exportID, count uint32) (*capnp.Client, error) {
	ent := c.findExport(id)
	if ent == nil {
		return nil, errorf("unknown export ID %d", id)
	}
	switch {
	case count == ent.wireRefs:
		client := ent.client
		c.exports[id] = nil
		c.exportID.remove(uint32(id))
		return client, nil
	case count > ent.wireRefs:
		return nil, errorf("export ID %d released too many references", id)
	default:
		ent.wireRefs -= count
		return nil, nil
	}
}

func (c *Conn) releaseExports(refs map[exportID]uint32) (releaseList, error) {
	n := len(refs)
	var rl releaseList
	var firstErr error
	for id, count := range refs {
		client, err := c.releaseExport(id, count)
		if err != nil {
			if firstErr == nil {
				firstErr = err
			}
			n--
			continue
		}
		if client == nil {
			n--
			continue
		}
		if rl == nil {
			rl = make(releaseList, 0, n)
		}
		rl = append(rl, client)
		n--
	}
	return rl, firstErr
}

// sendCap writes a capability descriptor, returning an export ID if
// this vat is hosting the capability.  The caller must be holding
// onto c.mu.
func (c *Conn) sendCap(d rpccp.CapDescriptor, client *capnp.Client, state capnp.ClientState) (_ exportID, isExport bool) {
	if !client.IsValid() {
		d.SetNone()
		return 0, false
	}

	if ic, ok := state.Brand.Value.(*importClient); ok && ic.c == c {
		if ent := c.imports[ic.id]; ent != nil && ent.generation == ic.generation {
			d.SetReceiverHosted(uint32(ic.id))
			return 0, false
		}
	}
	// TODO(someday): Check for unresolved client for senderPromise.
	// TODO(someday): Check for pipeline client on question for receiverAnswer.

	// Default to sender-hosted (export).
	for id, ent := range c.exports {
		if ent.client.IsSame(client) {
			ent.wireRefs++
			d.SetSenderHosted(uint32(id))
			return exportID(id), true
		}
	}
	ee := &expent{
		client:   client.AddRef(),
		wireRefs: 1,
	}
	id := exportID(c.exportID.next())
	if int64(id) == int64(len(c.exports)) {
		c.exports = append(c.exports, ee)
	} else {
		c.exports[id] = ee
	}
	d.SetSenderHosted(uint32(id))
	return id, true
}

// fillPayloadCapTable adds descriptors of payload's message's
// capabilities into payload's capability table and returns the
// reference counts added to the exports table.
//
// The caller must be holding onto c.mu.
func (c *Conn) fillPayloadCapTable(payload rpccp.Payload, clients []*capnp.Client, states []capnp.ClientState) (map[exportID]uint32, error) {
	if len(clients) != len(states) {
		panic("states slice must be same size as cap table")
	}
	if !payload.IsValid() || len(clients) == 0 {
		return nil, nil
	}
	list, err := payload.NewCapTable(int32(len(clients)))
	if err != nil {
		return nil, errorf("payload capability table: %v", err)
	}
	var refs map[exportID]uint32
	for i, client := range clients {
		id, isExport := c.sendCap(list.At(i), client, states[i])
		if !isExport {
			continue
		}
		if refs == nil {
			refs = make(map[exportID]uint32, len(clients)-i)
		}
		refs[id]++
	}
	return refs, nil
}

// extractCapTable reads the state of all the capabilities in a
// message's capability table and sets the message's capability table
// to nil.  The caller must not be holding onto any locks, since this
// function can call application code (ClientHook.Brand).
func extractCapTable(msg *capnp.Message) ([]*capnp.Client, []capnp.ClientState) {
	if len(msg.CapTable) == 0 {
		msg.CapTable = nil // in case msg.CapTable is a 0-length slice
		return nil, nil
	}
	ctab := msg.CapTable
	msg.CapTable = nil
	states := make([]capnp.ClientState, len(ctab))
	for i, c := range ctab {
		states[i] = c.State()
	}
	return ctab, states
}

type embargoID uint32

type embargo struct {
	c      *capnp.Client
	p      *capnp.ClientPromise
	lifted chan struct{}
}

// embargo creates a new embargoed client, stealing the reference.
//
// The caller must be holding onto c.mu.
func (c *Conn) embargo(client *capnp.Client) (embargoID, *capnp.Client) {
	id := embargoID(c.embargoID.next())
	e := &embargo{
		c:      client,
		lifted: make(chan struct{}),
	}
	if int64(id) == int64(len(c.embargoes)) {
		c.embargoes = append(c.embargoes, e)
	} else {
		c.embargoes[id] = e
	}
	var c2 *capnp.Client
	c2, c.embargoes[id].p = capnp.NewPromisedClient(c.embargoes[id])
	return id, c2
}

// findEmbargo returns the embargo entry with the given ID or nil if
// couldn't be found.
func (c *Conn) findEmbargo(id embargoID) *embargo {
	if int64(id) >= int64(len(c.embargoes)) {
		return nil
	}
	return c.embargoes[id] // might be nil
}

// lift disembargoes the client.  It must be called only once.
func (e *embargo) lift() {
	close(e.lifted)
	e.p.Fulfill(e.c)
}

func (e *embargo) Send(ctx context.Context, s capnp.Send) (*capnp.Answer, capnp.ReleaseFunc) {
	select {
	case <-e.lifted:
		return e.c.SendCall(ctx, s)
	case <-ctx.Done():
		return capnp.ErrorAnswer(s.Method, ctx.Err()), func() {}
	}
}

func (e *embargo) Recv(ctx context.Context, r capnp.Recv) capnp.PipelineCaller {
	select {
	case <-e.lifted:
		return e.c.RecvCall(ctx, r)
	case <-ctx.Done():
		r.Reject(ctx.Err())
		return nil
	}
}

func (e *embargo) Brand() capnp.Brand {
	return capnp.Brand{}
}

func (e *embargo) Shutdown() {
	e.c.Release()
}

// senderLoopback holds the salient information for a sender-loopback
// Disembargo message.
type senderLoopback struct {
	id        embargoID
	question  questionID
	transform []capnp.PipelineOp
}

func (sl *senderLoopback) buildDisembargo(msg rpccp.Message) error {
	d, err := msg.NewDisembargo()
	if err != nil {
		return errorf("build disembargo: %v", err)
	}
	tgt, err := d.NewTarget()
	if err != nil {
		return errorf("build disembargo: %v", err)
	}
	pa, err := tgt.NewPromisedAnswer()
	if err != nil {
		return errorf("build disembargo: %v", err)
	}
	oplist, err := pa.NewTransform(int32(len(sl.transform)))
	if err != nil {
		return errorf("build disembargo: %v", err)
	}

	d.Context().SetSenderLoopback(uint32(sl.id))
	pa.SetQuestionId(uint32(sl.question))
	for i, op := range sl.transform {
		oplist.At(i).SetGetPointerField(op.Field)
	}
	return nil
}

type releaseList []*capnp.Client

func (rl releaseList) release() {
	for _, c := range rl {
		c.Release()
	}
	for i := range rl {
		rl[i] = nil
	}
}
