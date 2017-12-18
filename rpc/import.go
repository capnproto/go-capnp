package rpc

import (
	"context"
	"fmt"

	"zombiezen.com/go/capnproto2"
	rpccp "zombiezen.com/go/capnproto2/std/capnp/rpc"
)

// An importID is an index into the imports table.
type importID uint32

// impent is an entry in the import table.  All fields are protected by
// Conn.mu.
type impent struct {
	wc *capnp.WeakClient

	// wireRefs is the number of times that the importID has appeared in
	// messages received from the remote vat.  Used to populate the
	// Release.referenceCount field.
	wireRefs int

	// generation is a counter used to disambiguate the following
	// condition:
	//
	// 1) An import given to application code.
	// 2) A new reference to the import is received over the wire, while
	//    the application concurrently closes the import.
	//    importClient.Close is called, but the receive has the lock
	//    first.
	// 3) Conn.addImport attempts to return a weak client reference, but
	//    can't because it has been closed.  It creates a new client
	//    with a new importClient.
	// 4) importClient.Close attempts to remove the import from the import
	//    table.  This is the critical step: it needs to be informed that
	//    it should not do this because another client has been created.
	//    No release message should be sent.
	//
	// The generation counter solves this by amending steps 3 and 4.  When
	// a new importClient is created, generation is incremented.
	// When importClient.Close is called, then it must check that the
	// importClient's generation matches the entry's generation before
	// removing the entry from the table and sending a release message.
	generation uint64
}

// addImport returns a client that represents the given import,
// incrementing the number of references to this import from this vat.
// This is separate from the reference counting that capnp.Client does.
func (c *Conn) addImport(id importID) *capnp.Client {
	if c.imports == nil {
		c.imports = make(map[importID]*impent)
	} else if ent := c.imports[id]; ent != nil {
		ent.wireRefs++
		client, ok := ent.wc.AddRef()
		if !ok {
			ent.generation++
			client = capnp.NewClient(&importClient{
				id:         id,
				conn:       c,
				generation: ent.generation,
			})
			ent.wc = client.WeakRef()
		}
		return client
	}
	client := capnp.NewClient(&importClient{
		id:   id,
		conn: c,
	})
	c.imports[id] = &impent{
		wc:       client.WeakRef(),
		wireRefs: 1,
	}
	return client
}

// An importClient implements capnp.Client for a remote capability.
type importClient struct {
	id         importID
	conn       *Conn
	generation uint64
}

func (ic *importClient) Send(ctx context.Context, s capnp.Send) (*capnp.Answer, capnp.ReleaseFunc) {
	defer ic.conn.mu.Unlock()
	ic.conn.mu.Lock()
	ent := ic.conn.imports[ic.id]
	if ent == nil || ic.generation != ent.generation {
		return capnp.ErrorAnswer(s.Method, disconnected("send on closed import")), func() {}
	}
	id := questionID(ic.conn.questionID.next())
	err := ic.conn.sendMessage(ctx, func(msg rpccp.Message) error {
		msgCall, err := msg.NewCall()
		if err != nil {
			return err
		}
		msgCall.SetQuestionId(uint32(id))
		msgCall.SetInterfaceId(s.Method.InterfaceID)
		msgCall.SetMethodId(s.Method.MethodID)
		target, err := msgCall.NewTarget()
		if err != nil {
			return err
		}
		target.SetImportedCap(uint32(ic.id))
		payload, err := msgCall.NewParams()
		if err != nil {
			return err
		}
		args, err := capnp.NewStruct(payload.Segment(), s.ArgsSize)
		if err != nil {
			return err
		}
		if err := payload.SetContent(args.ToPtr()); err != nil {
			return err
		}
		if err := s.PlaceArgs(args); err != nil {
			// Using fmt.Errorf to annotate to avoid stutter when we wrap the
			// sendMessage error.
			return fmt.Errorf("place parameters: %v", err)
		}
		// TODO(soon): fill in capability table
		return nil
	})
	if err != nil {
		ic.conn.questionID.remove(uint32(id))
		return capnp.ErrorAnswer(s.Method, annotate(err).errorf("send to import")), func() {}
	}
	q := ic.conn.newQuestion(ctx, id, s.Method)
	ans := q.p.Answer()
	return ans, func() {
		<-ans.Done()
		q.p.ReleaseClients()
		q.release()
	}
}

func (ic *importClient) Recv(ctx context.Context, r capnp.Recv) capnp.PipelineCaller {
	ans, finish := ic.Send(ctx, capnp.Send{
		Method:   r.Method,
		ArgsSize: r.Args.Size(),
		PlaceArgs: func(s capnp.Struct) error {
			err := s.CopyFrom(r.Args)
			r.ReleaseArgs()
			return err
		},
	})
	r.ReleaseArgs()
	select {
	case <-ans.Done():
		returnAnswer(r.Returner, ans, finish)
		return nil
	default:
		go returnAnswer(r.Returner, ans, finish)
		return ans
	}
}

func returnAnswer(ret capnp.Returner, ans *capnp.Answer, finish func()) {
	defer finish()
	result, err := ans.Struct()
	if err != nil {
		ret.Return(err)
		return
	}
	recvResult, err := ret.AllocResults(result.Size())
	if err != nil {
		ret.Return(err)
		return
	}
	if err := recvResult.CopyFrom(result); err != nil {
		ret.Return(err)
		return
	}
	ret.Return(nil)
}

func (ic *importClient) Brand() interface{} {
	return ic
}

func (ic *importClient) Shutdown() {
	defer ic.conn.mu.Unlock()
	ic.conn.mu.Lock()
	ent := ic.conn.imports[ic.id]
	if ent == nil || ic.generation != ent.generation {
		// If ent == nil, then the Conn was already shut down.  Otherwise,
		// a new reference was added concurrently with the Shutdown.  See
		// impent.generation documentation for an explanation.
		return
	}
	delete(ic.conn.imports, ic.id)
	err := ic.conn.sendMessage(context.Background(), func(msg rpccp.Message) error {
		rel, err := msg.NewRelease()
		if err != nil {
			return err
		}
		rel.SetId(uint32(ic.id))
		rel.SetReferenceCount(uint32(ent.wireRefs))
		return nil
	})
	if err != nil {
		ic.conn.report(annotate(err).errorf("send release"))
		return
	}
}
