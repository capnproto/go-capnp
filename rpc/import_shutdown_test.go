package rpc

import (
	"context"
	"sync/atomic"
	"testing"
	"testing/synctest"

	"capnproto.org/go/capnp/v3"
	transportpkg "capnproto.org/go/capnp/v3/rpc/transport"
	rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"
)

type trackedImportResolver struct {
	conn          *Conn
	inner         capnp.Resolver[capnp.Client]
	rejectEntered chan<- error
	allowReject   <-chan struct{}
	fulfills      atomic.Int32
	rejects       atomic.Int32
}

func (r *trackedImportResolver) Fulfill(client capnp.Client) {
	r.fulfills.Add(1)
	r.inner.Fulfill(client)
}

func (r *trackedImportResolver) Reject(err error) {
	// Resolver callbacks may reenter the connection.  This acquisition
	// deadlocks if shutdown invokes Reject while holding Conn.lk.
	r.conn.withLocked(func(*lockedConn) {})
	r.rejects.Add(1)
	r.rejectEntered <- err
	<-r.allowReject
	r.inner.Reject(err)
}

func TestImportedPromiseShutdownRejectsBeforeDone(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		left, right := transportpkg.NewPipe(16)
		conn := NewConn(NewTransport(left), nil)
		peer := NewTransport(right)
		defer func() {
			if err := conn.Close(); err != nil {
				t.Error("close connection:", err)
			}
		}()

		rejectEntered := make(chan error, 2)
		allowReject := make(chan struct{})
		tracker := &trackedImportResolver{
			conn:          conn,
			rejectEntered: rejectEntered,
			allowReject:   allowReject,
		}

		const id = importID(7)
		var client capnp.Client
		conn.withLocked(func(c *lockedConn) {
			client = c.addImport(id, true)
			imp, ok := c.lk.imports.Find(id)
			if !ok {
				t.Fatal("promise import was not created")
			}
			tracker.inner = imp.takeResolver()
			imp.resolver = tracker
		})

		answer, releaseAnswer := client.SendCall(context.Background(), capnp.Send{
			Method: capnp.Method{InterfaceID: 1, MethodID: 2},
		})
		defer releaseAnswer()
		defer client.Release()

		in, err := peer.RecvMessage()
		if err != nil {
			t.Fatal("receive queued promise call:", err)
		}
		call, err := in.Message().Call()
		if err != nil {
			t.Fatal("read queued promise call:", err)
		}
		target, err := call.Target()
		if err != nil {
			t.Fatal("read queued promise target:", err)
		}
		if target.Which() != rpccp.MessageTarget_Which_importedCap ||
			importID(target.ImportedCap()) != id {
			t.Fatalf("queued promise target = %v; want imported capability %d", target, id)
		}
		in.Release()

		resolveResult := make(chan error, 1)
		go func() {
			resolveResult <- client.Resolve(context.Background())
		}()

		if err := peer.Close(); err != nil {
			t.Fatal("close scripted peer:", err)
		}

		rejectErr := <-rejectEntered
		if !capnp.IsDisconnected(rejectErr) {
			t.Fatalf("promise rejection = %v; want disconnected error", rejectErr)
		}
		select {
		case <-conn.Done():
			t.Fatal("Conn.Done closed before imported promise rejection completed")
		default:
		}

		close(allowReject)
		synctest.Wait()

		if got := tracker.rejects.Load(); got != 1 {
			t.Fatalf("promise Reject calls = %d; want 1", got)
		}
		if got := tracker.fulfills.Load(); got != 0 {
			t.Fatalf("promise Fulfill calls = %d; want 0", got)
		}
		select {
		case <-conn.Done():
		default:
			t.Fatal("Conn.Done remains open after imported promise rejection completed")
		}
		conn.withLocked(func(c *lockedConn) {
			if _, ok := c.lk.imports.Find(id); ok {
				t.Error("promise import remains registered after shutdown")
			}
		})

		if err := <-resolveResult; err != nil {
			t.Fatalf("Client.Resolve error = %v; want successful resolution to an error client", err)
		}
		snapshot := client.Snapshot()
		resolutionErr, ok := snapshot.Brand().Value.(error)
		snapshot.Release()
		if !ok || !capnp.IsDisconnected(resolutionErr) {
			t.Fatalf("resolved promise brand = %v; want disconnected error", resolutionErr)
		}
		if _, err := answer.Struct(); !capnp.IsDisconnected(err) {
			t.Fatalf("queued call error = %v; want disconnected error", err)
		}

		if err := conn.Close(); err != nil {
			t.Fatal("second connection close:", err)
		}
		if got := tracker.rejects.Load(); got != 1 {
			t.Fatalf("promise Reject calls after second Close = %d; want 1", got)
		}
	})
}
