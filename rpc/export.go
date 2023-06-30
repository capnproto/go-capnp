package rpc

import (
	"context"
	"errors"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/exc"
	"capnproto.org/go/capnp/v3/internal/str"
	"capnproto.org/go/capnp/v3/internal/syncutil"
	rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"
	"zenhack.net/go/util/deferred"
	"zenhack.net/go/util/maybe"
)

// An exportID is an index into the exports table.
type exportID uint32

// expent is an entry in a Conn's export table.
type expent struct {
	snapshot capnp.ClientSnapshot
	wireRefs uint32

	// Should be called when removing this entry from the exports table:
	cancel context.CancelFunc

	// If present, this export is a promise which resolved to some third
	// party capability, and the question corresponds to the provide message.
	// Note that this means the question will belong to a different connection
	// from the expent.
	provide maybe.Maybe[*question]
}

// A key for use in a client's Metadata, whose value is the export
// id (if any) via which the client is exposed on the given
// connection.
type exportIDKey struct {
	Conn *Conn
}

func (c *lockedConn) findExportID(m *capnp.Metadata) (_ exportID, ok bool) {
	maybeID, ok := m.Get(exportIDKey{(*Conn)(c)})
	if ok {
		return maybeID.(exportID), true
	}
	return 0, false
}

func (c *lockedConn) setExportID(m *capnp.Metadata, id exportID) {
	m.Put(exportIDKey{(*Conn)(c)}, id)
}

func (c *lockedConn) clearExportID(m *capnp.Metadata) {
	m.Delete(exportIDKey{(*Conn)(c)})
}

// findExport returns the export entry with the given ID or nil if
// couldn't be found. The caller must be holding c.mu
func (c *lockedConn) findExport(id exportID) *expent {
	if int64(id) >= int64(len(c.lk.exports)) {
		return nil
	}
	return c.lk.exports[id] // might be nil
}

// releaseExport decreases the number of wire references to an export
// by a given number.  If the export's reference count reaches zero,
// then releaseExport will pop export from the table and schedule further
// cleanup (such as releasing snaphost) via dq.
func (c *lockedConn) releaseExport(dq *deferred.Queue, id exportID, count uint32) error {
	ent := c.findExport(id)
	if ent == nil {
		return rpcerr.Failed(errors.New("unknown export ID " + str.Utod(id)))
	}
	switch {
	case count == ent.wireRefs:
		defer ent.cancel()
		snapshot := ent.snapshot
		c.lk.exports[id] = nil
		c.lk.exportID.remove(id)
		metadata := snapshot.Metadata()
		if metadata != nil {
			syncutil.With(metadata, func() {
				c.clearExportID(metadata)
			})
		}
		dq.Defer(snapshot.Release)
		return nil
	case count > ent.wireRefs:
		return rpcerr.Failed(errors.New("export ID " + str.Utod(id) + " released too many references"))
	default:
		ent.wireRefs -= count
		return nil
	}
}

func (c *lockedConn) releaseExportRefs(dq *deferred.Queue, refs map[exportID]uint32) error {
	var firstErr error
	for id, count := range refs {
		err := c.releaseExport(dq, id, count)
		if err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// send3PHPromise begins the process of performing a third party handoff,
// passing srcSnapshot across c. srcSnapshot must point to an object across
// srcConn.
//
// This will store a senderPromise capability in d, which will later be
// resolved to a thirdPartyHosted cap by a separate goroutine.
//
// Returns the export ID for the promise.
func (c *lockedConn) send3PHPromise(
	d rpccp.CapDescriptor,
	srcConn *Conn,
	srcSnapshot capnp.ClientSnapshot,
	target parsedMessageTarget,
) exportID {
	if c.network != srcConn.network {
		panic("BUG: tried to do 3PH between different networks")
	}

	p, r := capnp.NewLocalPromise[capnp.Client]()
	defer p.Release()
	pSnapshot := p.Snapshot()
	r.Fulfill(srcSnapshot.Client()) // FIXME: this may allow path shortening...

	// TODO(cleanup): most of this is copypasta from sendExport; consider
	// ways to factor out the common bits.
	promiseID := c.lk.exportID.next()
	metadata := pSnapshot.Metadata()
	metadata.Lock()
	defer metadata.Unlock()
	c.setExportID(metadata, promiseID)
	ee := &expent{
		snapshot: pSnapshot,
		wireRefs: 1,
		cancel:   func() {},
	}
	c.insertExport(promiseID, ee)
	d.SetSenderPromise(uint32(promiseID))

	go func() {
		c := (*Conn)(c)
		defer srcSnapshot.Release()

		// TODO(cleanup): we should probably make the src/dest arguments
		// consistent across all 3PH code:
		introInfo, err := c.network.Introduce(srcConn, c)
		if err != nil {
			// TODO: report somehow; see above.
			return
		}

		// XXX: think about what we should be doing for contexts here:
		var (
			provideQ *question
			vine     *vine
		)
		ctx, cancel := context.WithCancel(srcConn.bgctx)
		srcConn.withLocked(func(c *lockedConn) {
			provideQ = c.newQuestion(capnp.Method{})
			provideQ.flags |= isProvide
			vine = newVine(srcSnapshot.AddRef(), cancel)
			c.sendMessage(c.bgctx, func(m rpccp.Message) error {
				provide, err := m.NewProvide()
				if err != nil {
					return err
				}
				provide.SetQuestionId(uint32(provideQ.id))
				if err = provide.SetRecipient(capnp.Ptr(introInfo.SendToProvider)); err != nil {
					return err
				}
				encodedTgt, err := provide.NewTarget()
				if err != nil {
					return err
				}
				if err = target.Encode(encodedTgt); err != nil {
					return err
				}
				return nil
			}, func(err error) {
				if err != nil {
					srcConn.withLocked(func(c *lockedConn) {
						c.lk.questionID.remove(provideQ.id)
					})
					return
				}
				go provideQ.handleCancel(ctx)
			})
		})
		unlockedConn := c
		var (
			vineID    exportID
			vineEntry *expent
		)
		c.withLocked(func(c *lockedConn) {
			c.sendMessage(c.bgctx, func(m rpccp.Message) error {
				if len(c.lk.exports) <= int(promiseID) || c.lk.exports[promiseID] != ee {
					// At some point the receiver lost interest in the cap.
					// Return an error to indicate we didn't send the resolve:
					return errReceiverLostInterest
				}
				// We have to set this before sending the provide, so we're ready
				// for a disembargo. It's okay to wait up until now though, since
				// the receiver shouldn't send one until it sees the resolution:
				c.lk.exports[promiseID].provide = maybe.New(provideQ)

				resolve, err := m.NewResolve()
				if err != nil {
					return err
				}
				resolve.SetPromiseId(uint32(promiseID))
				capDesc, err := resolve.NewCap()
				if err != nil {
					return err
				}
				thirdCapDesc, err := capDesc.NewThirdPartyHosted()
				if err != nil {
					return err
				}
				if err = thirdCapDesc.SetId(capnp.Ptr(introInfo.SendToRecipient)); err != nil {
					return err
				}

				vineID = c.lk.exportID.next()
				client := capnp.NewClient(vine)
				defer client.Release()

				c.insertExport(vineID, &expent{
					snapshot: client.Snapshot(),
					wireRefs: 1,
					cancel:   cancel,
				})
				vineEntry = c.lk.exports[vineID]

				thirdCapDesc.SetVineId(uint32(vineID))
				return nil
			}, func(err error) {
				if err == nil {
					return
				}
				if vineEntry == nil {
					vine.Shutdown()
				} else {
					dq := &deferred.Queue{}
					defer dq.Run()
					unlockedConn.withLocked(func(c *lockedConn) {
						c.releaseExport(dq, vineID, 1)
					})
				}
			})
		})
	}()
	return promiseID
}

var errReceiverLostInterest = errors.New("receiver lost interest in the resolution")

// sendCap writes a capability descriptor, returning an export ID if
// this vat is hosting the capability. Steals the snapshot.
func (c *lockedConn) sendCap(d rpccp.CapDescriptor, snapshot capnp.ClientSnapshot) (_ exportID, isExport bool, _ error) {
	if !snapshot.IsValid() {
		d.SetNone()
		return 0, false, nil
	}

	defer snapshot.Release()
	bv := snapshot.Brand().Value
	unlockedConn := (*Conn)(c)
	if ic, ok := bv.(*importClient); ok {
		if ic.c == unlockedConn {
			if ent := c.lk.imports[ic.id]; ent != nil && ent.generation == ic.generation {
				d.SetReceiverHosted(uint32(ic.id))
				return 0, false, nil
			}
		}
		if c.network != nil && c.network == ic.c.network {

			exportID := c.send3PHPromise(
				d, ic.c, snapshot.AddRef(),
				parsedMessageTarget{
					which:       rpccp.MessageTarget_Which_importedCap,
					importedCap: exportID(ic.id),
				},
			)
			return exportID, true, nil
		}
	} else if pc, ok := bv.(capnp.PipelineClient); ok {
		if q, ok := c.getAnswerQuestion(pc.Answer()); ok {
			if q.c == unlockedConn {
				pcTrans := pc.Transform()
				pa, err := d.NewReceiverAnswer()
				if err != nil {
					return 0, false, err
				}
				trans, err := pa.NewTransform(int32(len(pcTrans)))
				if err != nil {
					return 0, false, err
				}
				for i, op := range pcTrans {
					trans.At(i).SetGetPointerField(op.Field)
				}
				pa.SetQuestionId(uint32(q.id))
				return 0, false, nil
			}
			if c.network != nil && c.network == q.c.network {
				exportID := c.send3PHPromise(
					d, q.c, snapshot.AddRef(),
					parsedMessageTarget{
						which:          rpccp.MessageTarget_Which_promisedAnswer,
						transform:      pc.Transform(),
						promisedAnswer: answerID(q.id),
					},
				)
				return exportID, true, nil
			}
		}
	}

	// Default to export.
	return c.sendExport(d, snapshot), true, nil
}

func (c *lockedConn) insertExport(id exportID, ee *expent) {
	if int64(id) == int64(len(c.lk.exports)) {
		c.lk.exports = append(c.lk.exports, ee)
	} else {
		c.lk.exports[id] = ee
	}
}

// sendExport is a helper for sendCap that handles the export cases.
func (c *lockedConn) sendExport(d rpccp.CapDescriptor, snapshot capnp.ClientSnapshot) exportID {
	metadata := snapshot.Metadata()
	metadata.Lock()
	defer metadata.Unlock()
	id, ok := c.findExportID(metadata)
	var ee *expent
	if ok {
		ee = c.lk.exports[id]
		ee.wireRefs++
	} else {
		// Not already present; allocate an export id for it:
		ee = &expent{
			snapshot: snapshot.AddRef(),
			wireRefs: 1,
			cancel:   func() {},
		}
		id = c.lk.exportID.next()
		c.insertExport(id, ee)
		c.setExportID(metadata, id)
	}
	if ee.snapshot.IsPromise() {
		c.sendSenderPromise(id, d)
	} else {
		d.SetSenderHosted(uint32(id))
	}
	return id
}

// sendSenderPromise is a helper for sendExport that handles the senderPromise case.
func (c *lockedConn) sendSenderPromise(id exportID, d rpccp.CapDescriptor) {
	// Send a promise, wait for the resolution asynchronously, then send
	// a resolve message:
	ee := c.lk.exports[id]
	d.SetSenderPromise(uint32(id))
	ctx, cancel := context.WithCancel(c.bgctx)
	ee.cancel = cancel
	waitRef := ee.snapshot.AddRef()
	go func() {
		defer cancel()
		defer waitRef.Release()
		// Logically we don't hold the lock anymore; it's held by the
		// goroutine that spawned this one. So cast back to an unlocked
		// Conn before trying to use it again:
		unlockedConn := (*Conn)(c)

		waitErr := waitRef.Resolve1(ctx)
		unlockedConn.withLocked(func(c *lockedConn) {
			if len(c.lk.exports) <= int(id) || c.lk.exports[id] != ee {
				// Export was removed from the table at some point;
				// remote peer is uninterested in the resolution, so
				// drop the reference and we're done
				return
			}

			sendRef := waitRef.AddRef()
			var (
				resolvedID exportID
				isExport   bool
			)
			c.sendMessage(c.bgctx, func(m rpccp.Message) error {
				res, err := m.NewResolve()
				if err != nil {
					return err
				}
				res.SetPromiseId(uint32(id))
				if waitErr != nil {
					ex, err := res.NewException()
					if err != nil {
						return err
					}
					return ex.MarshalError(waitErr)
				}
				desc, err := res.NewCap()
				if err != nil {
					return err
				}
				resolvedID, isExport, err = c.sendCap(desc, sendRef)
				return err
			}, func(err error) {
				sendRef.Release()
				if err != nil && isExport {
					dq := &deferred.Queue{}
					defer dq.Run()
					// release 1 ref of the thing it resolved to.
					err := withLockedConn1(
						unlockedConn,
						func(c *lockedConn) error {
							return c.releaseExport(dq, resolvedID, 1)
						},
					)
					if err != nil {
						c.er.ReportError(
							exc.WrapError("releasing export due to failure to send resolve", err),
						)
					}
				}
			})
		})
	}()
}

// fillPayloadCapTable adds descriptors of payload's message's
// capabilities into payload's capability table and returns the
// reference counts that have been added to the exports table.
func (c *lockedConn) fillPayloadCapTable(payload rpccp.Payload) (map[exportID]uint32, error) {
	if !payload.IsValid() {
		return nil, nil
	}
	clients := payload.Message().CapTable()
	if clients.Len() == 0 {
		return nil, nil
	}
	list, err := payload.NewCapTable(int32(clients.Len()))
	if err != nil {
		return nil, rpcerr.WrapFailed("payload capability table", err)
	}
	var refs map[exportID]uint32
	for i := 0; i < clients.Len(); i++ {
		id, isExport, err := c.sendCap(list.At(i), clients.At(i).Snapshot())
		if err != nil {
			return nil, rpcerr.WrapFailed("Serializing capability", err)
		}
		if isExport {
			if refs == nil {
				refs = make(map[exportID]uint32, clients.Len()-i)
			}
			refs[id]++
		}
	}
	return refs, nil
}

type embargoID uint32

type embargo struct {
	result capnp.Ptr
	q      *capnp.AnswerQueue
}

func (e embargo) String() string {
	return "embargo{client: " +
		e.client().String() +
		", q: 0x" + str.PtrToHex(e.q) +
		"}"
}

// embargo creates a new embargoed client, stealing the reference.
//
// The caller must be holding onto c.mu.
func (c *lockedConn) embargo(client capnp.Client) (embargoID, capnp.Client) {
	id := c.lk.embargoID.next()
	e := newEmbargo(client)
	if int64(id) == int64(len(c.lk.embargoes)) {
		c.lk.embargoes = append(c.lk.embargoes, e)
	} else {
		c.lk.embargoes[id] = e
	}
	return id, capnp.NewClient(e)
}

// findEmbargo returns the embargo entry with the given ID or nil if
// couldn't be found.
func (c *lockedConn) findEmbargo(id embargoID) *embargo {
	if int64(id) >= int64(len(c.lk.embargoes)) {
		return nil
	}
	return c.lk.embargoes[id] // might be nil
}

func newEmbargo(client capnp.Client) *embargo {
	msg, seg := capnp.NewSingleSegmentMessage(nil)
	capID := msg.CapTable().Add(client)
	iface := capnp.NewInterface(seg, capID)
	return &embargo{
		result: iface.ToPtr(),
		q:      capnp.NewAnswerQueue(capnp.Method{}),
	}
}

// lift disembargoes the client.  It must be called only once.
func (e *embargo) lift() {
	e.q.Fulfill(e.result)
}

func (e *embargo) Send(ctx context.Context, s capnp.Send) (*capnp.Answer, capnp.ReleaseFunc) {
	return e.q.PipelineSend(ctx, nil, s)
}

func (e *embargo) Recv(ctx context.Context, r capnp.Recv) capnp.PipelineCaller {

	return e.q.PipelineRecv(ctx, nil, r)
}

func (e *embargo) Brand() capnp.Brand {
	return capnp.Brand{}
}

func (e *embargo) client() capnp.Client {
	return e.result.Interface().Client()
}

func (e *embargo) Shutdown() {
	e.client().Release()
}

// senderLoopback holds the salient information for a sender-loopback
// Disembargo message.
type senderLoopback struct {
	id     embargoID
	target parsedMessageTarget
}

func (sl *senderLoopback) buildDisembargo(msg rpccp.Message) error {
	d, err := msg.NewDisembargo()
	if err != nil {
		return rpcerr.WrapFailed("build disembargo", err)
	}
	d.Context().SetSenderLoopback(uint32(sl.id))
	tgt, err := d.NewTarget()
	if err != nil {
		return rpcerr.WrapFailed("build disembargo", err)
	}
	switch sl.target.which {
	case rpccp.MessageTarget_Which_promisedAnswer:
		pa, err := tgt.NewPromisedAnswer()
		if err != nil {
			return rpcerr.WrapFailed("build disembargo", err)
		}
		oplist, err := pa.NewTransform(int32(len(sl.target.transform)))
		if err != nil {
			return rpcerr.WrapFailed("build disembargo", err)
		}

		pa.SetQuestionId(uint32(sl.target.promisedAnswer))
		for i, op := range sl.target.transform {
			oplist.At(i).SetGetPointerField(op.Field)
		}
	case rpccp.MessageTarget_Which_importedCap:
		tgt.SetImportedCap(uint32(sl.target.importedCap))
	default:
		return errors.New("unknown variant for MessageTarget: " + str.Utod(sl.target.which))
	}
	return nil
}
