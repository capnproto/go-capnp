// Package proc assits managing the lifecycle of goroutines.
package proc

import (
	"context"
	"sync"
)

// Spawn a new process.
//
// A process can be in one of the following states:
//
// - Running
// - Stopping
// - Stopped
//
// ...starting in the Running state.
//
// f is a callback that implements the actual logic of the process. It will be
// passed a child context of ctx, and a Self. It should:
//
// 1. Run until its context is canceled is canceled (or its work is done).
// 2. invoke BeginShutdown() on the Self that was passed to it, which
//    transitions the process to the Stopping state.
// 3. Perform any shutdown logic that is needed.
// 4. Return. This transitions the process into the Stopped state.
//
// Returns a handle, which can be used by clients of the process to interact
// with its lifecycle.
func Spawn(ctx context.Context, f func(context.Context, Self)) Handle {
	ctx, cancel := context.WithCancel(ctx)
	proc := &proc{
		cancel: cancel,
		done:   make(chan struct{}),
	}
	handle := Handle{proc}
	self := Self{proc}
	go func() {
		defer func() {
			self.BeginShutdown()
			close(proc.done)
		}()
		f(ctx, self)
	}()
	return handle
}

// A Handle is used for interacting with a process.
type Handle struct {
	proc *proc
}

// Cancel cancels the process's context.
func (h Handle) Cancel() {
	h.proc.cancel()
}

// Done returns a channel that will be closed when the process transitions
// to the Stopped state.
func (h Handle) Done() <-chan struct{} {
	return h.proc.done
}

// WithLive attempts to invoke f while keeping the process in the Running state.
//
// If the process has already exited the Running state, WithLive returns false
// without invoking f.
//
// If the process is still in the Running state, WithLive invokes the callback,
// while preventing the process from entering the Stopping state, and then returns
// true. If the process calls Self.BeginShutdown, it will block until f returns.
func (h Handle) WithLive(f func()) (ok bool) {
	h.proc.mu.Lock()
	defer h.proc.mu.Unlock()
	if h.proc.shuttingDown {
		return false
	}
	f()
	return true
}

// A Self is passed to the process to help manage its own lifecycle.
type Self struct {
	proc *proc
}

// BeginShutdown transitions the process from the Running state to the Stopped
// state, waiting for any ongoing calls to Handle.WithLive to complete. Once
// this returns, any calls to WithLive will return false without executing
// their callback.
func (s Self) BeginShutdown() {
	s.proc.mu.Lock()
	defer s.proc.mu.Unlock()
	s.proc.shuttingDown = true
}

// Internal state of the process, shared by Self and Handle.
type proc struct {
	mu           sync.Mutex
	shuttingDown bool
	cancel       context.CancelFunc
	done         chan struct{}
}
