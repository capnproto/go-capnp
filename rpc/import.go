package rpc

import (
	"context"
	"errors"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/internal/str"
	rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"
)

// An importID is an index into the imports table.
type importID uint32

// impent is an entry in the import table.  All fields are protected by
// Conn.mu.
type impent struct {
	wc capnp.WeakClient

	// wireRefs is the number of times that the importID has appeared in
	// messages received from the remote vat.  Used to populate the
	// Release.referenceCount field.
	wireRefs int

	// generation identifies the newest importClient hook. Only that hook may
	// send calls or be reflected back to the peer. It is incremented when a
	// weak client cursor can no longer be upgraded and addImport creates a
	// replacement hook.
	generation uint64

	// liveHooks counts importClient hooks that have not called Shutdown.
	// Snapshots retain hooks independently of client cursors, so hooks from
	// older generations may remain live after generation advances. The import
	// entry and its promise resolver must remain registered until every hook
	// generation shuts down.
	liveHooks int

	// If resolver is non-nil, then this is a promise (received as
	// CapDescriptor_Which_senderPromise), and when a resolve message
	// arrives we should use this to fulfill the promise locally.
	resolver capnp.Resolver[capnp.Client]
}

// importPromiseResolvers settles every live local view of one unresolved
// import. A new client generation gets a new resolver, while snapshots may
// keep earlier generations alive after their client cursors are released.
type importPromiseResolvers []capnp.Resolver[capnp.Client]

func (rs importPromiseResolvers) Fulfill(client capnp.Client) {
	for _, resolver := range rs {
		resolver.Fulfill(client)
	}
}

func (rs importPromiseResolvers) Reject(err error) {
	for _, resolver := range rs {
		resolver.Reject(err)
	}
}

func joinImportPromiseResolvers(
	first, second capnp.Resolver[capnp.Client],
) capnp.Resolver[capnp.Client] {
	if resolvers, ok := first.(importPromiseResolvers); ok {
		return append(resolvers, second)
	}
	return importPromiseResolvers{first, second}
}

// takeResolver transfers terminal ownership of an imported promise.
// The caller must hold c.lk.
func (e *impent) takeResolver() capnp.Resolver[capnp.Client] {
	resolver := e.resolver
	e.resolver = nil
	return resolver
}

// addImport returns a client that represents the given import,
// incrementing the number of references to this import from this vat.
// This is separate from the reference counting that capnp.Client does.
//
// The caller must be holding onto c.mu.
func (c *lockedConn) addImport(id importID, isPromise bool) capnp.Client {
	if ent, _ := c.lk.imports.Find(id); ent != nil {
		ent.wireRefs++
		client, ok := ent.wc.AddRef()
		if !ok {
			ent.generation++
			ent.liveHooks++
			hook := &importClient{
				c:          (*Conn)(c),
				id:         id,
				generation: ent.generation,
			}
			if isPromise && ent.resolver != nil {
				var resolver capnp.Resolver[capnp.Client]
				client, resolver = capnp.NewPromisedClient(hook)
				ent.resolver = joinImportPromiseResolvers(ent.resolver, resolver)
			} else {
				client = capnp.NewClient(hook)
			}
			ent.wc = client.WeakRef()
		}
		return client
	}
	hook := &importClient{
		c:  (*Conn)(c),
		id: id,
	}
	var (
		client   capnp.Client
		resolver capnp.Resolver[capnp.Client]
	)
	if isPromise {
		client, resolver = capnp.NewPromisedClient(hook)
	} else {
		client = capnp.NewClient(hook)
	}
	c.lk.imports.Create(id, &impent{
		wc:        client.WeakRef(),
		wireRefs:  1,
		liveHooks: 1,
		resolver:  resolver,
	})
	return client
}

// An importClient implements capnp.Client for a remote capability.
type importClient struct {
	c          *Conn
	id         importID
	generation uint64
}

func (ic *importClient) String() string {
	return "importClient{c: 0x" + str.PtrToHex(ic.c) + ", id: " + str.Utod(ic.id) + "}"
}

func (ic *importClient) Send(ctx context.Context, s capnp.Send) (*capnp.Answer, capnp.ReleaseFunc) {
	return withLockedConn2(ic.c, func(c *lockedConn) (*capnp.Answer, capnp.ReleaseFunc) {
		return c.startCall(ctx, s, func() error {
			ent, _ := c.lk.imports.Find(ic.id)
			if ent == nil || ic.generation != ent.generation {
				return rpcerr.Disconnected(errors.New("send on closed import"))
			}
			return nil
		}, func(target rpccp.MessageTarget) error {
			target.SetImportedCap(uint32(ic.id))
			return nil
		})
	})
}

// PrepareSend lets Client's flow-control path reserve against the fully built
// RPC message before it is placed on the transport queue.
func (ic *importClient) PrepareSend(ctx context.Context, s capnp.Send) (capnp.PreparedSend, error) {
	return ic.c.prepareCall(ctx, s, func(c *lockedConn) error {
		ent, _ := c.lk.imports.Find(ic.id)
		if ent == nil || ic.generation != ent.generation {
			return rpcerr.Disconnected(errors.New("send on closed import"))
		}
		return nil
	}, func(target rpccp.MessageTarget) error {
		target.SetImportedCap(uint32(ic.id))
		return nil
	})
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
	defer ret.ReleaseResults()
	result, err := ans.Struct()
	if err != nil {
		ret.PrepareReturn(err)
		ret.Return()
		return
	}
	recvResult, err := ret.AllocResults(result.Size())
	if err != nil {
		ret.PrepareReturn(err)
		ret.Return()
		return
	}
	if err := recvResult.CopyFrom(result); err != nil {
		ret.PrepareReturn(err)
		ret.Return()
		return
	}
	ret.PrepareReturn(nil)
	ret.Return()
}

func (ic *importClient) Brand() capnp.Brand {
	return capnp.Brand{Value: ic}
}

func (ic *importClient) Shutdown() {
	ic.c.withLocked(func(c *lockedConn) {
		if !c.startTask() {
			return
		}
		defer c.tasks.Done()

		ent, ok := c.lk.imports.Find(ic.id)
		if !ok {
			return
		}
		if ent.liveHooks <= 0 {
			panic("rpc: import has no live hooks")
		}
		ent.liveHooks--
		if ent.liveHooks > 0 {
			return
		}
		c.lk.imports.Remove(ic.id)
		c.sendMessage(c.bgctx, func(msg rpccp.Message) error {
			rel, err := msg.NewRelease()
			if err == nil {
				rel.SetId(uint32(ic.id))
				rel.SetReferenceCount(uint32(ent.wireRefs))
			}
			return err
		}, func(err error) {
			if err != nil && !isFatalSendError(err) {
				ic.c.er.ReportError(rpcerr.Annotate(err, "send release"))
			}
		})
	})
}
