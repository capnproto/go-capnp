package rpc

import (
	"sync"

	"golang.org/x/net/context"
)

// manager signals the running goroutines in a Conn.
// Since there is one manager per connection, it's also a way of
// identifying an object's origin.
type manager struct {
	finish chan struct{}
	wg     sync.WaitGroup
	ctx    context.Context

	mu   sync.RWMutex
	done bool
	e    error
}

func (m *manager) init() {
	m.finish = make(chan struct{})
	var cancel context.CancelFunc
	m.ctx, cancel = context.WithCancel(context.Background())
	go func() {
		<-m.finish
		cancel()
	}()
}

// context returns a context that is cancelled when the manager shuts down.
func (m *manager) context() context.Context {
	return m.ctx
}

// do starts a function in a new goroutine and will block shutdown
// until it has returned.  If the manager has already started shutdown,
// then it is a no-op.
func (m *manager) do(f func()) {
	m.mu.RLock()
	done := m.done
	if !done {
		m.wg.Add(1)
	}
	m.mu.RUnlock()
	if !done {
		go func() {
			defer m.wg.Done()
			f()
		}()
	}
}

// shutdown closes the finish channel and sets the error.  The first
// call to shutdown returns true; subsequent calls are no-ops and return
// false.  This will not wait for the manager's goroutines to finish.
func (m *manager) shutdown(e error) bool {
	m.mu.Lock()
	ok := !m.done
	if ok {
		close(m.finish)
		m.done = true
		m.e = e
	}
	m.mu.Unlock()
	return ok
}

// wait blocks until the manager is shut down and all of its goroutines
// are finished.
func (m *manager) wait() {
	<-m.finish
	m.wg.Wait()
}

// err returns the error passed to shutdown.
func (m *manager) err() error {
	m.mu.RLock()
	e := m.e
	m.mu.RUnlock()
	return e
}
