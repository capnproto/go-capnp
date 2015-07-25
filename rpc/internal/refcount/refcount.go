// Package refcount implements a reference-counting client.
package refcount

import (
	"errors"
	"runtime"
	"sync"

	"zombiezen.com/go/capnproto"
)

// A RefCount will close its underlying client once all its references are closed.
type RefCount struct {
	Client capnp.Client

	mu   sync.Mutex
	refs int
}

// New creates a reference counter and the first client reference.
func New(c capnp.Client) (rc *RefCount, ref capnp.Client) {
	rc = &RefCount{Client: c}
	ref = rc.Ref()
	return
}

// Ref makes a new client reference.
func (rc *RefCount) Ref() capnp.Client {
	// TODO(light): what if someone calls Ref() after refs hits zero?
	rc.mu.Lock()
	rc.refs++
	rc.mu.Unlock()
	r := &ref{rc: rc}
	runtime.SetFinalizer(r, (*ref).Close)
	return r
}

func (rc *RefCount) call(cl *capnp.Call) capnp.Answer {
	// We lock here so that we can prevent the client from being closed
	// while we start the call.
	rc.mu.Lock()
	defer rc.mu.Unlock()
	if rc.refs <= 0 {
		return capnp.ErrorAnswer(errClosed)
	}
	return rc.Client.Call(cl)
}

// decref decreases the reference count by one, closing the Client if it reaches zero.
func (rc *RefCount) decref() error {
	shouldClose := false

	rc.mu.Lock()
	if rc.refs <= 0 {
		rc.mu.Unlock()
		return errClosed
	}
	rc.refs--
	if rc.refs == 0 {
		shouldClose = true
	}
	rc.mu.Unlock()

	if shouldClose {
		return rc.Client.Close()
	}
	return nil
}

var errClosed = errors.New("rpc: Close() called on closed client")

type ref struct {
	rc *RefCount
}

func (r *ref) Call(cl *capnp.Call) capnp.Answer {
	return r.rc.call(cl)
}

func (r *ref) Close() error {
	return r.rc.decref()
}
