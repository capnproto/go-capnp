// Package thunk provides "thunks" (type T), which are wrappers around lazy or
// concurrent computations.
package thunk

import "sync"

// A T[A] wraps a value A which may not be available yet; goroutines may get the
// value by calling Force.
type T[A any] struct {
	once  sync.Once
	f     func() A
	value A
}

// Force waits for the result of the T to be available and returns it.
func (t *T[A]) Force() A {
	t.once.Do(func() {
		t.value = t.f()
		t.f = nil
	})
	return t.value
}

// Lazy returns a new T which, when forced, will return the value returned by f().
// f() will be invoked lazily, i.e. not until Force() is called for the first time.
func Lazy[A any](f func() A) *T[A] {
	return &T[A]{f: f}
}

// Go is like Lazy, except that f() is invoked immediately in a separate goroutine.
// Calls to Force will block until f() returns.
func Go[A any](f func() A) *T[A] {
	ret := Lazy(f)
	go ret.Force()
	return ret
}

// Promise returns a pair of a T and a function fulfill to supply the value of
// the T; calls to t.Force() will block until fulfill has been invoked, at which
// point they will return the value passed to fulfill. Fulfill must not be called
// more than once.
func Promise[A any]() (t *T[A], fulfill func(A)) {
	ch := make(chan A, 1)
	t = Lazy(func() A {
		return <-ch
	})
	return t, func(val A) {
		ch <- val
		close(ch)
	}
}
