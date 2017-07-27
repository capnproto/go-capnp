// Package chanmu provides a mutex that can be used in a select.
package chanmu

import "context"

// Mutex is a mutex implemented in terms of a channel.  Receiving from
// a mutex acquires it; sending to a mutex releases it.
type Mutex chan struct{}

// New returns a new mutex.
func New() Mutex {
	mu := make(Mutex, 1)
	mu <- struct{}{}
	return mu
}

// Lock acquires mu.
func (mu Mutex) Lock() {
	if mu == nil {
		panic("Lock on nil Mutex")
	}
	<-mu
}

// TryLock acquires mu or returns an error if the Context finished first.
func (mu Mutex) TryLock(ctx context.Context) error {
	if mu == nil {
		panic("Lock on nil Mutex")
	}
	select {
	case <-mu:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Unlock releases mu.
func (mu Mutex) Unlock() {
	select {
	case mu <- struct{}{}:
	default:
		panic("double Unlock")
	}
}
