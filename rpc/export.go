package rpc

import (
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
func (c *Conn) sendCap(d rpccp.CapDescriptor, client *capnp.Client) (_ exportID, isExport bool) {
	if !client.IsValid() {
		d.SetNone()
		return 0, false
	}

	brand := client.Brand()
	if ic, ok := brand.(*importClient); ok && ic.conn == c {
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
	id := exportID(len(c.exports))
	c.exports = append(c.exports, &expent{
		client:   client.AddRef(),
		wireRefs: 1,
	})
	d.SetSenderHosted(uint32(id))
	return id, true
}

// fillPayloadCapTable adds descriptors of payload's message's
// capabilities into payload's capability table and returns the
// reference counts added to the exports table.  The caller must be
// holding onto c.mu.
func (c *Conn) fillPayloadCapTable(payload rpccp.Payload) (map[exportID]uint32, error) {
	msg := payload.Message()
	if msg == nil || len(msg.CapTable) == 0 {
		return nil, nil
	}
	list, err := payload.NewCapTable(int32(len(msg.CapTable)))
	if err != nil {
		return nil, errorf("payload capability table: %v", err)
	}
	var refs map[exportID]uint32
	for i, client := range msg.CapTable {
		id, isExport := c.sendCap(list.At(i), client)
		if !isExport {
			continue
		}
		if refs == nil {
			refs = make(map[exportID]uint32, len(msg.CapTable)-i)
		}
		refs[id]++
	}
	return refs, nil
}

type releaseList []*capnp.Client

func (rl releaseList) release() {
	for i := range rl {
		rl[i].Release()
		rl[i] = nil
	}
}
