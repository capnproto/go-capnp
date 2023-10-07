// Package thunk provides "thunks" (type T), which are wrappers around lazy or
// concurrent computations.
package thunk

import "sync"

// A Thunk[A] wraps a value A which may not be available yet; goroutines may get the
// value by calling Force.
type Thunk[A any] struct {
	once  sync.Once
	f     func() A
	value A
}

// Force waits for the result of the Thunk to be available and returns it.
func (t *Thunk[A]) Force() A {
	t.once.Do(func() {
		t.value = t.f()
		t.f = nil
	})
	return t.value
}

// Lazy returns a new Thunk which, when forced, will return the value returned by f().
// f() will be invoked lazily, i.e. not until Force() is called for the first time.
func Lazy[A any](f func() A) *Thunk[A] {
	return &Thunk[A]{f: f}
}

// Go is like Lazy, except that f() is invoked immediately in a separate goroutine.
// Calls to Force will block until f() returns.
func Go[A any](f func() A) *Thunk[A] {
	ret := Lazy(f)
	go ret.Force()
	return ret
}

// Ready returns a new Thunk which is already ready; when forced it will return
// value immediately.
func Ready[A any](value A) *Thunk[A] {
	var ret Thunk[A]
	ret.once.Do(func() {})
	ret.value = value
	return &ret
}

// Promise returns a pair of a Thunk and a function fulfill to supply the value of
// the Thunk; calls to t.Force() will block until fulfill has been invoked, at which
// point they will return the value passed to fulfill. Fulfill must not be called
// more than once.
func Promise[A any]() (t *Thunk[A], fulfill func(A)) {
	ch := make(chan A, 1)
	t = Lazy(func() A {
		return <-ch
	})
	return t, func(val A) {
		ch <- val
		close(ch)
	}
}
