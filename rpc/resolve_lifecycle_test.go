package rpc

import (
	"context"
	"encoding/binary"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"capnproto.org/go/capnp/v3"
	transportpkg "capnproto.org/go/capnp/v3/rpc/transport"
	"capnproto.org/go/capnp/v3/server"
	rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"
)

const (
	resolvePromiseID     = importID(7)
	resolveReplacementID = importID(8)
)

func newResolveLifecycleConn(t *testing.T) (*Conn, Transport) {
	return newResolveLifecycleConnWithOptions(t, nil)
}

func newResolveLifecycleConnWithOptions(t *testing.T, opts *Options) (*Conn, Transport) {
	t.Helper()

	left, right := transportpkg.NewPipe(16)
	conn := NewConn(NewTransport(left), opts)
	peer := NewTransport(right)
	t.Cleanup(func() {
		if err := conn.Close(); err != nil {
			t.Error("close connection:", err)
		}
		if err := peer.Close(); err != nil {
			t.Error("close peer transport:", err)
		}
	})
	return conn, peer
}

func addResolvePromise(t *testing.T, conn *Conn) capnp.Client {
	t.Helper()

	var promise capnp.Client
	conn.withLocked(func(c *lockedConn) {
		promise = c.addImport(resolvePromiseID, true)
	})
	t.Cleanup(promise.Release)
	return promise
}

func newResolveMessage(
	t *testing.T,
	promiseID importID,
	build func(rpccp.Resolve) error,
) *countingIncomingMessage {
	t.Helper()

	_, seg := capnp.NewSingleSegmentMessage(nil)
	message, err := rpccp.NewRootMessage(seg)
	if err != nil {
		t.Fatal(err)
	}
	resolve, err := message.NewResolve()
	if err != nil {
		t.Fatal(err)
	}
	resolve.SetPromiseId(uint32(promiseID))
	if err := build(resolve); err != nil {
		t.Fatal(err)
	}
	return &countingIncomingMessage{message: message}
}

func resolveToSenderHosted(t *testing.T, promiseID, replacementID importID) *countingIncomingMessage {
	t.Helper()
	return newResolveMessage(t, promiseID, func(resolve rpccp.Resolve) error {
		desc, err := resolve.NewCap()
		if err == nil {
			desc.SetSenderHosted(uint32(replacementID))
		}
		return err
	})
}

func resolveToSenderPromise(t *testing.T, promiseID, replacementID importID) *countingIncomingMessage {
	t.Helper()
	return newResolveMessage(t, promiseID, func(resolve rpccp.Resolve) error {
		desc, err := resolve.NewCap()
		if err == nil {
			desc.SetSenderPromise(uint32(replacementID))
		}
		return err
	})
}

func resolveToThirdPartyHosted(t *testing.T, promiseID, vineID importID) *countingIncomingMessage {
	t.Helper()
	return newResolveMessage(t, promiseID, func(resolve rpccp.Resolve) error {
		desc, err := resolve.NewCap()
		if err != nil {
			return err
		}
		thirdParty, err := desc.NewThirdPartyHosted()
		if err == nil {
			thirdParty.SetVineId(uint32(vineID))
		}
		return err
	})
}

func resolveToReceiverHosted(t *testing.T, promiseID importID, exportID exportID) *countingIncomingMessage {
	t.Helper()
	return newResolveMessage(t, promiseID, func(resolve rpccp.Resolve) error {
		desc, err := resolve.NewCap()
		if err == nil {
			desc.SetReceiverHosted(uint32(exportID))
		}
		return err
	})
}

func resolveToReceiverAnswer(t *testing.T, promiseID importID, answerID answerID) *countingIncomingMessage {
	t.Helper()
	return newResolveMessage(t, promiseID, func(resolve rpccp.Resolve) error {
		desc, err := resolve.NewCap()
		if err != nil {
			return err
		}
		answer, err := desc.NewReceiverAnswer()
		if err == nil {
			answer.SetQuestionId(uint32(answerID))
		}
		return err
	})
}

func resolveToException(t *testing.T, promiseID importID) *countingIncomingMessage {
	t.Helper()
	return newResolveMessage(t, promiseID, func(resolve rpccp.Resolve) error {
		ex, err := resolve.NewException()
		if err == nil {
			ex.SetType(rpccp.Exception_Type_failed)
			err = ex.SetReason("resolution failed")
		}
		return err
	})
}

func queueResolveMarker(t *testing.T, conn *Conn, questionID uint32) {
	t.Helper()

	done := make(chan error, 1)
	conn.withLocked(func(c *lockedConn) {
		c.sendMessage(context.Background(), func(message rpccp.Message) error {
			bootstrap, err := message.NewBootstrap()
			if err == nil {
				bootstrap.SetQuestionId(questionID)
			}
			return err
		}, func(err error) {
			done <- err
		})
	})
	select {
	case err := <-done:
		if err != nil {
			t.Fatal("send marker:", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timed out sending marker")
	}
}

func requireResolveMarker(t *testing.T, peer Transport, questionID uint32) {
	t.Helper()

	in := recvPeerMessage(t, peer)
	defer in.Release()
	if got := in.Message().Which(); got != rpccp.Message_Which_bootstrap {
		t.Fatalf("next message = %v; want bootstrap marker", got)
	}
	bootstrap, err := in.Message().Bootstrap()
	if err != nil {
		t.Fatal(err)
	}
	if got := bootstrap.QuestionId(); got != questionID {
		t.Fatalf("marker question ID = %d; want %d", got, questionID)
	}
}

func requireResolveRelease(
	t *testing.T,
	peer Transport,
	id importID,
	referenceCount uint32,
) {
	t.Helper()

	in := recvPeerMessage(t, peer)
	defer in.Release()
	if got := in.Message().Which(); got != rpccp.Message_Which_release {
		t.Fatalf("next message = %v; want release", got)
	}
	release, err := in.Message().Release()
	if err != nil {
		t.Fatal(err)
	}
	if got := importID(release.Id()); got != id {
		t.Fatalf("release import ID = %d; want %d", got, id)
	}
	if got := release.ReferenceCount(); got != referenceCount {
		t.Fatalf("release reference count = %d; want %d", got, referenceCount)
	}
}

type resolveShutdownFunc func()

func (f resolveShutdownFunc) Shutdown() {
	f()
}

type resolveTestNetwork struct{}

func (resolveTestNetwork) LocalID() PeerID             { return PeerID{} }
func (resolveTestNetwork) Dial(PeerID) (*Conn, error)  { return nil, nil }
func (resolveTestNetwork) Serve(context.Context) error { return nil }

func newTrackedResolveClient() (capnp.Client, <-chan struct{}) {
	shutdown := make(chan struct{})
	client := capnp.NewClient(server.New(nil, nil, resolveShutdownFunc(func() {
		close(shutdown)
	})))
	return client, shutdown
}

func addReflectedExport(
	t *testing.T,
	conn *Conn,
) (exportID, func(), <-chan struct{}) {
	t.Helper()

	client, shutdown := newTrackedResolveClient()
	snapshot := client.Snapshot()
	var id exportID
	conn.withLocked(func(c *lockedConn) {
		id = c.lk.exports.Add(&expent{
			snapshot: snapshot,
			wireRefs: 1,
			cancel:   func() {},
		})
	})

	var once sync.Once
	release := func() {
		once.Do(func() {
			var ent *expent
			conn.withLocked(func(c *lockedConn) {
				ent, _ = c.lk.exports.Remove(id)
			})
			if ent != nil {
				ent.snapshot.Release()
			}
			client.Release()
		})
	}
	t.Cleanup(release)
	return id, release, shutdown
}

func addReflectedAnswer(
	t *testing.T,
	conn *Conn,
) (answerID, func(), <-chan struct{}) {
	t.Helper()

	client, shutdown := newTrackedResolveClient()
	snapshot := client.Snapshot()

	message, seg := capnp.NewSingleSegmentMessage(nil)
	results, err := rpccp.NewRootPayload(seg)
	if err != nil {
		t.Fatal(err)
	}
	capID := message.CapTable().Add(client.AddRef())
	if err := results.SetContent(capnp.NewInterface(seg, capID).ToPtr()); err != nil {
		t.Fatal(err)
	}

	const id = answerID(43)
	conn.withLocked(func(c *lockedConn) {
		if !c.lk.answers.Create(id, &ansent{
			returner: ansReturner{
				results:         results,
				resultsCapTable: []capnp.ClientSnapshot{snapshot},
			},
		}) {
			t.Fatal("create reflected answer")
		}
	})

	var once sync.Once
	release := func() {
		once.Do(func() {
			var answer *ansent
			conn.withLocked(func(c *lockedConn) {
				answer, _ = c.lk.answers.Remove(id)
			})
			if answer != nil {
				for i := range answer.returner.resultsCapTable {
					answer.returner.resultsCapTable[i].Release()
				}
			}
			message.Release()
			client.Release()
		})
	}
	t.Cleanup(release)
	return id, release, shutdown
}

func addResolveAnswerReturningClient(
	t *testing.T,
	conn *Conn,
	client capnp.Client,
) answerID {
	t.Helper()

	message, seg := capnp.NewSingleSegmentMessage(nil)
	results, err := rpccp.NewRootPayload(seg)
	if err != nil {
		t.Fatal(err)
	}
	capID := message.CapTable().Add(client.AddRef())
	if err := results.SetContent(capnp.NewInterface(seg, capID).ToPtr()); err != nil {
		t.Fatal(err)
	}
	snapshot := client.Snapshot()

	const id = answerID(44)
	conn.withLocked(func(c *lockedConn) {
		if !c.lk.answers.Create(id, &ansent{
			returner: ansReturner{
				results:         results,
				resultsCapTable: []capnp.ClientSnapshot{snapshot},
			},
		}) {
			t.Fatal("create self-referencing answer")
		}
	})
	t.Cleanup(func() {
		var answer *ansent
		conn.withLocked(func(c *lockedConn) {
			answer, _ = c.lk.answers.Remove(id)
		})
		if answer != nil {
			for i := range answer.returner.resultsCapTable {
				answer.returner.resultsCapTable[i].Release()
			}
		}
		message.Release()
	})
	return id
}

func requireResolveClientShutdown(t *testing.T, shutdown <-chan struct{}) {
	t.Helper()

	select {
	case <-shutdown:
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for reflected replacement shutdown")
	}
}

func TestHandleResolveFulfillsKnownPromiseAndReleasesMessage(t *testing.T) {
	conn, _ := newResolveLifecycleConn(t)
	promise := addResolvePromise(t, conn)
	in := resolveToSenderHosted(t, resolvePromiseID, resolveReplacementID)

	if err := conn.handleResolve(context.Background(), in); err != nil {
		t.Fatal("handleResolve:", err)
	}
	if got := atomic.LoadInt32(&in.releases); got != 1 {
		t.Fatalf("incoming message releases = %d; want 1", got)
	}
	if err := promise.Resolve(context.Background()); err != nil {
		t.Fatal("resolve imported promise:", err)
	}
	conn.withLocked(func(c *lockedConn) {
		if _, ok := c.lk.imports.Find(resolveReplacementID); !ok {
			t.Error("resolved capability was released instead of transferred to the promise")
		}
	})
}

func TestHandleResolveRejectsKnownPromiseAndReleasesMessage(t *testing.T) {
	conn, _ := newResolveLifecycleConn(t)
	promise := addResolvePromise(t, conn)
	in := resolveToException(t, resolvePromiseID)

	if err := conn.handleResolve(context.Background(), in); err != nil {
		t.Fatal("handleResolve:", err)
	}
	if got := atomic.LoadInt32(&in.releases); got != 1 {
		t.Fatalf("incoming message releases = %d; want 1", got)
	}
	if err := promise.Resolve(context.Background()); err != nil {
		t.Fatal("resolve imported promise:", err)
	}
	snapshot := promise.Snapshot()
	resolutionErr, ok := snapshot.Brand().Value.(error)
	snapshot.Release()
	if !ok || !strings.Contains(resolutionErr.Error(), "resolution failed") {
		t.Fatalf("resolved promise brand = %v; want resolution failure", resolutionErr)
	}
}

func TestHandleResolveUnknownPromiseReleasesImportedReplacement(t *testing.T) {
	tests := []struct {
		name string
		new  func(*testing.T, importID, importID) *countingIncomingMessage
	}{
		{name: "senderHosted", new: resolveToSenderHosted},
		{name: "senderPromise", new: resolveToSenderPromise},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			conn, peer := newResolveLifecycleConn(t)
			in := test.new(t, resolvePromiseID, resolveReplacementID)

			if err := conn.handleResolve(context.Background(), in); err != nil {
				t.Fatal("handleResolve:", err)
			}
			if got := atomic.LoadInt32(&in.releases); got != 1 {
				t.Fatalf("incoming message releases = %d; want 1", got)
			}
			conn.withLocked(func(c *lockedConn) {
				if _, ok := c.lk.imports.Find(resolveReplacementID); ok {
					t.Error("replacement capability remains imported after late Resolve")
				}
			})
			select {
			case <-conn.Done():
				t.Error("late Resolve closed the connection")
			default:
			}

			const markerQuestionID = 91
			queueResolveMarker(t, conn, markerQuestionID)
			requireResolveRelease(t, peer, resolveReplacementID, 1)
			requireResolveMarker(t, peer, markerQuestionID)
		})
	}
}

func TestHandleResolveUnknownPromiseExceptionIsNoOp(t *testing.T) {
	conn, peer := newResolveLifecycleConn(t)
	in := resolveToException(t, resolvePromiseID)

	if err := conn.handleResolve(context.Background(), in); err != nil {
		t.Fatal("handleResolve:", err)
	}
	if got := atomic.LoadInt32(&in.releases); got != 1 {
		t.Fatalf("incoming message releases = %d; want 1", got)
	}
	select {
	case <-conn.Done():
		t.Error("late exception Resolve closed the connection")
	default:
	}

	const markerQuestionID = 92
	queueResolveMarker(t, conn, markerQuestionID)
	requireResolveMarker(t, peer, markerQuestionID)
}

func TestHandleResolveUnknownPromiseReleasesReflectedReplacementLocally(t *testing.T) {
	tests := []struct {
		name  string
		setup func(*testing.T, *Conn) (
			*countingIncomingMessage,
			func(),
			<-chan struct{},
		)
	}{
		{
			name: "receiverHosted",
			setup: func(t *testing.T, conn *Conn) (
				*countingIncomingMessage,
				func(),
				<-chan struct{},
			) {
				id, release, shutdown := addReflectedExport(t, conn)
				return resolveToReceiverHosted(t, resolvePromiseID, id), release, shutdown
			},
		},
		{
			name: "receiverAnswer",
			setup: func(t *testing.T, conn *Conn) (
				*countingIncomingMessage,
				func(),
				<-chan struct{},
			) {
				id, release, shutdown := addReflectedAnswer(t, conn)
				return resolveToReceiverAnswer(t, resolvePromiseID, id), release, shutdown
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			conn, peer := newResolveLifecycleConn(t)
			in, releaseOwner, shutdown := test.setup(t, conn)

			if err := conn.handleResolve(context.Background(), in); err != nil {
				t.Fatal("handleResolve:", err)
			}
			if got := atomic.LoadInt32(&in.releases); got != 1 {
				t.Fatalf("incoming message releases = %d; want 1", got)
			}
			conn.withLocked(func(c *lockedConn) {
				embargoes := 0
				c.lk.embargoes.Range(func(_ embargoID, _ *embargo) bool {
					embargoes++
					return true
				})
				if embargoes != 0 {
					t.Errorf("late reflected Resolve created %d embargoes; want 0", embargoes)
				}
			})

			// All cleanup triggered by handleResolve is queued before this
			// marker. Seeing the marker first proves there was no Release or
			// Disembargo for the reflected replacement.
			const markerQuestionID = 93
			queueResolveMarker(t, conn, markerQuestionID)
			requireResolveMarker(t, peer, markerQuestionID)

			// Drop every owner that existed before handleResolve. The server
			// shuts down only if the replacement reference materialized by
			// recvCap was also released.
			releaseOwner()
			requireResolveClientShutdown(t, shutdown)
		})
	}
}

func TestHandleResolveInvalidDescriptorsReleaseMessage(t *testing.T) {
	tests := []struct {
		name      string
		networked bool
		build     func(rpccp.Resolve) error
		want      string
	}{
		{
			name: "none",
			build: func(resolve rpccp.Resolve) error {
				_, err := resolve.NewCap()
				return err
			},
			want: "cap descriptor is none",
		},
		{
			name: "unknown descriptor",
			build: func(resolve rpccp.Resolve) error {
				desc, err := resolve.NewCap()
				if err == nil {
					capnp.Struct(desc).SetUint16(0, 99)
				}
				return err
			},
			want: "unknown cap descriptor type",
		},
		{
			name: "unknown receiver-hosted target",
			build: func(resolve rpccp.Resolve) error {
				desc, err := resolve.NewCap()
				if err == nil {
					desc.SetReceiverHosted(99)
				}
				return err
			},
			want: "invalid export",
		},
		{
			name: "unknown receiver-answer target",
			build: func(resolve rpccp.Resolve) error {
				desc, err := resolve.NewCap()
				if err != nil {
					return err
				}
				answer, err := desc.NewReceiverAnswer()
				if err == nil {
					answer.SetQuestionId(99)
				}
				return err
			},
			want: "no such question id",
		},
		{
			name:      "networked third-party target",
			networked: true,
			build: func(resolve rpccp.Resolve) error {
				desc, err := resolve.NewCap()
				if err != nil {
					return err
				}
				thirdParty, err := desc.NewThirdPartyHosted()
				if err == nil {
					thirdParty.SetVineId(uint32(resolveReplacementID))
				}
				return err
			},
			want: "third-party handoff not implemented",
		},
		{
			name: "unknown resolve union",
			build: func(resolve rpccp.Resolve) error {
				capnp.Struct(resolve).SetUint16(4, 99)
				return nil
			},
			want: "unknown resolve type",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var opts *Options
			if test.networked {
				opts = &Options{Network: resolveTestNetwork{}}
			}
			conn, _ := newResolveLifecycleConnWithOptions(t, opts)
			in := newResolveMessage(t, resolvePromiseID, test.build)

			err := conn.handleResolve(context.Background(), in)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("handleResolve error = %v; want error containing %q", err, test.want)
			}
			if got := atomic.LoadInt32(&in.releases); got != 1 {
				t.Fatalf("incoming message releases = %d; want 1", got)
			}
		})
	}
}

func TestHandleResolveParseFailureReleasesMessage(t *testing.T) {
	conn, _ := newResolveLifecycleConn(t)
	in := malformedResolveMessage(t)

	if err := conn.handleResolve(context.Background(), in); err != nil {
		t.Fatal("handleResolve:", err)
	}
	if got := atomic.LoadInt32(&in.releases); got != 1 {
		t.Fatalf("incoming message releases = %d; want 1", got)
	}
}

func malformedResolveMessage(t *testing.T) *countingIncomingMessage {
	t.Helper()

	valid := resolveToSenderHosted(t, resolvePromiseID, resolveReplacementID)
	wire, err := valid.message.Message().Marshal()
	if err != nil {
		t.Fatal(err)
	}
	valid.Release()

	// In a single-segment message, the root Message struct starts at word 1.
	// Its first pointer (word 2) is the Resolve struct. Point it far beyond
	// the segment so Message.Resolve reports a bounds error.
	if len(wire) < 32 {
		t.Fatalf("encoded resolve has %d bytes; want at least 32", len(wire))
	}
	const invalidResolvePointer = uint64(0x1fffffff)<<2 | uint64(1)<<32 | uint64(1)<<48
	binary.LittleEndian.PutUint64(wire[24:32], invalidResolvePointer)

	message, err := capnp.Unmarshal(wire)
	if err != nil {
		t.Fatal(err)
	}
	root, err := rpccp.ReadRootMessage(message)
	if err != nil {
		t.Fatal(err)
	}
	return &countingIncomingMessage{message: root}
}

func TestHandleResolveAlreadyResolvedImportReleasesReplacement(t *testing.T) {
	conn, _ := newResolveLifecycleConn(t)
	var nonPromise capnp.Client
	conn.withLocked(func(c *lockedConn) {
		nonPromise = c.addImport(resolvePromiseID, false)
	})
	t.Cleanup(nonPromise.Release)
	in := resolveToSenderHosted(t, resolvePromiseID, resolveReplacementID)

	err := conn.handleResolve(context.Background(), in)
	if err == nil || !strings.Contains(err.Error(), "is not a promise") {
		t.Fatalf("handleResolve error = %v; want non-promise error", err)
	}
	if got := atomic.LoadInt32(&in.releases); got != 1 {
		t.Fatalf("incoming message releases = %d; want 1", got)
	}
	conn.withLocked(func(c *lockedConn) {
		if _, ok := c.lk.imports.Find(resolveReplacementID); ok {
			t.Error("replacement capability remains imported after invalid duplicate Resolve")
		}
	})
}

func TestHandleResolveRejectsSelfResolution(t *testing.T) {
	tests := []struct {
		name       string
		addPromise bool
		newResolve func(*testing.T, importID, importID) *countingIncomingMessage
	}{
		{
			name:       "unknown promise to sender promise",
			newResolve: resolveToSenderPromise,
		},
		{
			name:       "known promise to sender hosted",
			addPromise: true,
			newResolve: resolveToSenderHosted,
		},
		{
			name:       "known promise to sender promise",
			addPromise: true,
			newResolve: resolveToSenderPromise,
		},
		{
			name:       "known promise to third-party vine",
			addPromise: true,
			newResolve: resolveToThirdPartyHosted,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			conn, _ := newResolveLifecycleConn(t)
			if test.addPromise {
				addResolvePromise(t, conn)
			}
			in := test.newResolve(t, resolvePromiseID, resolvePromiseID)

			err := conn.handleResolve(context.Background(), in)
			if err == nil || !strings.Contains(err.Error(), "resolved to itself") {
				t.Fatalf("handleResolve error = %v; want self-resolution error", err)
			}
			if got := atomic.LoadInt32(&in.releases); got != 1 {
				t.Fatalf("incoming message releases = %d; want 1", got)
			}
			conn.withLocked(func(c *lockedConn) {
				ent, _ := c.lk.imports.Find(resolvePromiseID)
				if test.addPromise {
					if ent == nil || ent.wireRefs != 1 {
						t.Errorf("self-resolution import = %+v; want one unchanged wire reference", ent)
					}
				} else if ent != nil {
					t.Error("self-resolution created an import entry")
				}
			})
		})
	}
}

func TestHandleResolveRejectsReflectedSelfResolution(t *testing.T) {
	conn, _ := newResolveLifecycleConn(t)
	promise := addResolvePromise(t, conn)
	answerID := addResolveAnswerReturningClient(t, conn, promise)
	in := resolveToReceiverAnswer(t, resolvePromiseID, answerID)

	err := conn.handleResolve(context.Background(), in)
	if err == nil || !strings.Contains(err.Error(), "resolved to itself") {
		t.Fatalf("handleResolve error = %v; want self-resolution error", err)
	}
	if got := atomic.LoadInt32(&in.releases); got != 1 {
		t.Fatalf("incoming message releases = %d; want 1", got)
	}
	conn.withLocked(func(c *lockedConn) {
		ent, _ := c.lk.imports.Find(resolvePromiseID)
		if ent == nil || ent.wireRefs != 1 {
			t.Errorf("reflected self-resolution import = %+v; want one unchanged wire reference", ent)
		}
		embargoes := 0
		c.lk.embargoes.Range(func(_ embargoID, _ *embargo) bool {
			embargoes++
			return true
		})
		if embargoes != 0 {
			t.Errorf("reflected self-resolution created %d embargoes; want 0", embargoes)
		}
	})
}
