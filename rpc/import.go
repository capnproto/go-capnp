package rpc

import (
	"context"

	"capnproto.org/go/capnp/v3"
	rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"
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
//
// The caller must be holding onto c.mu.
func (c *Conn) addImport(id importID) *capnp.Client {
	if ent := c.imports[id]; ent != nil {
		ent.wireRefs++
		client, ok := ent.wc.AddRef()
		if !ok {
			ent.generation++
			client = capnp.NewClient(&importClient{
				c:          c,
				id:         id,
				generation: ent.generation,
			})
			ent.wc = client.WeakRef()
		}
		return client
	}
	client := capnp.NewClient(&importClient{
		c:  c,
		id: id,
	})
	c.imports[id] = &impent{
		wc:       client.WeakRef(),
		wireRefs: 1,
	}
	return client
}

// An importClient implements capnp.Client for a remote capability.
type importClient struct {
	c          *Conn
	id         importID
	generation uint64
}

func (ic *importClient) Send(ctx context.Context, s capnp.Send) (*capnp.Answer, capnp.ReleaseFunc) {
	// Acquire sender lock.
	ic.c.mu.Lock()
	if !ic.c.startTask() {
		ic.c.mu.Unlock()
		return capnp.ErrorAnswer(s.Method, ExcClosed), func() {}
	}
	defer ic.c.tasks.Done()
	ent := ic.c.imports[ic.id]
	if ent == nil || ic.generation != ent.generation {
		ic.c.mu.Unlock()
		return capnp.ErrorAnswer(s.Method, disconnectedf("send on closed import")), func() {}
	}
	if err := ic.c.tryLockSender(ctx); err != nil {
		ic.c.mu.Unlock()
		return capnp.ErrorAnswer(s.Method, err), func() {}
	}
	q := ic.c.newQuestion(s.Method)
	ic.c.mu.Unlock()

	// Create call message.
	msg, send, release, err := ic.c.transport.NewMessage(ctx)
	if err != nil {
		ic.c.mu.Lock()
		ic.c.questions[q.id] = nil
		ic.c.questionID.remove(uint32(q.id))
		ic.c.mu.Unlock()
		return capnp.ErrorAnswer(s.Method, failedf("create message: %w", err)), func() {}
	}
	ic.c.mu.Lock()
	ic.c.unlockSender() // Can't be holding either lock while calling PlaceArgs.
	ic.c.mu.Unlock()
	err = ic.c.newImportCallMessage(msg, ic.id, q.id, s)
	if err != nil {
		ic.c.mu.Lock()
		ic.c.questions[q.id] = nil
		ic.c.questionID.remove(uint32(q.id))
		ic.c.lockSender()
		ic.c.mu.Unlock()
		release()
		ic.c.mu.Lock()
		ic.c.unlockSender()
		ic.c.mu.Unlock()
		return capnp.ErrorAnswer(s.Method, err), func() {}
	}

	// Send call.
	ic.c.mu.Lock()
	ic.c.lockSender()
	ic.c.mu.Unlock()
	err = send()
	release()

	ic.c.mu.Lock()
	ic.c.unlockSender()
	if err != nil {
		ic.c.questions[q.id] = nil
		ic.c.questionID.remove(uint32(q.id))
		ic.c.mu.Unlock()
		return capnp.ErrorAnswer(s.Method, failedf("send message: %w", err)), func() {}
	}
	q.c.tasks.Add(1)
	go func() {
		defer q.c.tasks.Done()
		q.handleCancel(ctx)
	}()
	ic.c.mu.Unlock()

	ans := q.p.Answer()
	return ans, func() {
		<-ans.Done()
		q.p.ReleaseClients()
		q.release()
	}
}

// newImportCallMessage builds a Call message targeted to an import.
//
// The caller MUST NOT be holding onto c.mu or the sender lock.
func (c *Conn) newImportCallMessage(msg rpccp.Message, imp importID, qid questionID, s capnp.Send) error {
	call, err := msg.NewCall()
	if err != nil {
		return failedf("build call message: %w", err)
	}
	call.SetQuestionId(uint32(qid))
	call.SetInterfaceId(s.Method.InterfaceID)
	call.SetMethodId(s.Method.MethodID)
	target, err := call.NewTarget()
	if err != nil {
		return failedf("build call message: %w", err)
	}
	target.SetImportedCap(uint32(imp))
	payload, err := call.NewParams()
	if err != nil {
		return failedf("build call message: %w", err)
	}
	args, err := capnp.NewStruct(payload.Segment(), s.ArgsSize)
	if err != nil {
		return failedf("build call message: %w", err)
	}
	if err := payload.SetContent(args.ToPtr()); err != nil {
		return failedf("build call message: %w", err)
	}

	if s.PlaceArgs == nil {
		return nil
	}
	m := args.Message()
	if err := s.PlaceArgs(args); err != nil {
		for _, c := range m.CapTable {
			c.Release()
		}
		m.CapTable = nil
		return failedf("place arguments: %w", err)
	}
	clients := extractCapTable(m)
	c.mu.Lock()
	// TODO(soon): save param refs
	_, err = c.fillPayloadCapTable(payload, clients)
	c.mu.Unlock()
	releaseList(clients).release()
	if err != nil {
		return annotate(err, "build call message")
	}
	return nil
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

func (ic *importClient) Brand() capnp.Brand {
	return capnp.Brand{Value: ic}
}

func (ic *importClient) Shutdown() {
	ic.c.mu.Lock()
	if !ic.c.startTask() {
		ic.c.mu.Unlock()
		return
	}
	defer ic.c.tasks.Done()
	ent := ic.c.imports[ic.id]
	if ic.generation != ent.generation {
		// A new reference was added concurrently with the Shutdown.  See
		// impent.generation documentation for an explanation.
		ic.c.mu.Unlock()
		return
	}
	delete(ic.c.imports, ic.id)
	err := ic.c.sendMessage(context.Background(), func(msg rpccp.Message) error {
		rel, err := msg.NewRelease()
		if err != nil {
			return err
		}
		rel.SetId(uint32(ic.id))
		rel.SetReferenceCount(uint32(ent.wireRefs))
		return nil
	})
	ic.c.mu.Unlock()
	if err != nil {
		ic.c.report(annotate(err, "send release"))
	}
}
