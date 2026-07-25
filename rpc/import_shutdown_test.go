package rpc

import (
	"context"
	"strings"
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

type terminalCountingImportResolver struct {
	inner    capnp.Resolver[capnp.Client]
	fulfills atomic.Int32
	rejects  atomic.Int32
}

func (r *terminalCountingImportResolver) Fulfill(client capnp.Client) {
	r.fulfills.Add(1)
	r.inner.Fulfill(client)
}

func (r *terminalCountingImportResolver) Reject(err error) {
	r.rejects.Add(1)
	r.inner.Reject(err)
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

func TestImportedPromiseShutdownRejectsRetainedSnapshot(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		left, right := transportpkg.NewPipe(16)
		conn := NewConn(NewTransport(left), nil)
		peer := NewTransport(right)
		defer peer.Close()

		const id = importID(8)
		var client capnp.Client
		conn.withLocked(func(c *lockedConn) {
			client = c.addImport(id, true)
		})

		snapshot := client.Snapshot()
		defer snapshot.Release()
		client.Release()

		conn.withLocked(func(c *lockedConn) {
			imp, ok := c.lk.imports.Find(id)
			if !ok {
				t.Fatal("promise import was removed while its snapshot remained live")
			}
			if live, ok := imp.wc.AddRef(); ok {
				live.Release()
				t.Fatal("promise import retained a live client cursor")
			}
		})

		resolveCtx, cancelResolve := context.WithCancel(context.Background())
		defer cancelResolve()
		resolveResult := make(chan error, 1)
		go func() {
			resolveResult <- snapshot.Resolve(resolveCtx)
		}()

		closeResult := make(chan error, 1)
		go func() {
			closeResult <- conn.Close()
		}()
		synctest.Wait()

		if err := <-closeResult; err != nil {
			t.Fatal("close connection:", err)
		}
		select {
		case err := <-resolveResult:
			if err != nil {
				t.Fatalf("snapshot Resolve error = %v; want nil", err)
			}
		default:
			cancelResolve()
			synctest.Wait()
			<-resolveResult
			t.Fatal("snapshot remained unresolved after connection shutdown")
		}

		resolutionErr, ok := snapshot.Brand().Value.(error)
		if !ok || !capnp.IsDisconnected(resolutionErr) {
			t.Fatalf("resolved snapshot brand = %v; want disconnected error", resolutionErr)
		}
	})
}

func TestImportedPromiseShutdownRejectsSnapshotAcrossTransientGeneration(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		left, right := transportpkg.NewPipe(16)
		conn := NewConn(NewTransport(left), nil)
		peer := NewTransport(right)
		defer peer.Close()

		const id = importID(9)
		var first capnp.Client
		conn.withLocked(func(c *lockedConn) {
			first = c.addImport(id, true)
		})
		snapshot := first.Snapshot()
		defer snapshot.Release()
		first.Release()

		var current capnp.Client
		conn.withLocked(func(c *lockedConn) {
			current = c.addImport(id, true)
		})
		current.Release()

		conn.withLocked(func(c *lockedConn) {
			imp, ok := c.lk.imports.Find(id)
			if !ok {
				t.Fatal("transient current generation removed import retained by an older snapshot")
			}
			if imp.liveHooks != 1 {
				t.Fatalf("live import hooks = %d; want 1 retained snapshot hook", imp.liveHooks)
			}
		})

		if err := conn.Close(); err != nil {
			t.Fatal("close connection:", err)
		}
		if !snapshot.IsResolved() {
			t.Fatal("snapshot remained unresolved after connection shutdown")
		}
		if err := snapshot.Resolve(context.Background()); err != nil {
			t.Fatal("resolve retained snapshot:", err)
		}
		resolutionErr, ok := snapshot.Brand().Value.(error)
		if !ok || !capnp.IsDisconnected(resolutionErr) {
			t.Fatalf("resolved snapshot brand = %v; want disconnected error", resolutionErr)
		}
	})
}

func TestImportedPromiseResolveRaceShutdownHasSingleTerminalCallback(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		left, right := transportpkg.NewPipe(16)
		conn := NewConn(NewTransport(left), nil)
		peer := NewTransport(right)
		defer peer.Close()

		const id = importID(10)
		var client capnp.Client
		tracker := &terminalCountingImportResolver{}
		conn.withLocked(func(c *lockedConn) {
			client = c.addImport(id, true)
			imp, ok := c.lk.imports.Find(id)
			if !ok {
				t.Fatal("promise import was not created")
			}
			tracker.inner = imp.takeResolver()
			imp.resolver = tracker
		})
		snapshot := client.Snapshot()
		defer snapshot.Release()
		client.Release()

		in := resolveToException(t, id)
		start := make(chan struct{})
		resolveResult := make(chan error, 1)
		closeResult := make(chan error, 1)
		go func() {
			<-start
			resolveResult <- conn.handleResolve(context.Background(), in)
		}()
		go func() {
			<-start
			closeResult <- conn.Close()
		}()
		close(start)
		synctest.Wait()

		if err := <-resolveResult; err != nil {
			t.Fatal("handle Resolve:", err)
		}
		if err := <-closeResult; err != nil {
			t.Fatal("close connection:", err)
		}
		if got := tracker.fulfills.Load(); got != 0 {
			t.Fatalf("promise Fulfill calls = %d; want 0", got)
		}
		if got := tracker.rejects.Load(); got != 1 {
			t.Fatalf("promise Reject calls = %d; want 1", got)
		}
		if !snapshot.IsResolved() {
			t.Fatal("snapshot remained unresolved after Resolve/Close race")
		}
		if err := snapshot.Resolve(context.Background()); err != nil {
			t.Fatal("resolve retained snapshot:", err)
		}
		resolutionErr, ok := snapshot.Brand().Value.(error)
		if !ok || (!capnp.IsDisconnected(resolutionErr) &&
			!strings.Contains(resolutionErr.Error(), "resolution failed")) {
			t.Fatalf("resolved snapshot brand = %v; want Resolve or shutdown rejection", resolutionErr)
		}
	})
}
