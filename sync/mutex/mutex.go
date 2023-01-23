// Package mutex provides mutexes that wrap the thing they protect.
// This results in clearer code than using sync.Mutex directly.
package mutex

import "sync"

// A Mutex[T] wraps a T, protecting it with a mutex. It must not
// be moved after first use. The zero value is an unlocked mutex
// containing the zero value of T.
type Mutex[T any] struct {
	mu  sync.Mutex
	val T
}

// New returns a new mutex containing the value val.
func New[T any](val T) Mutex[T] {
	return Mutex[T]{val: val}
}

// With invokes the callback with exclusive access to the value of the mutex.
func (m *Mutex[T]) With(f func(*T)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	f(&m.val)
}

// With1(m, ...) is like m.With(...), except that the callback returns
// a value.
func With1[T, A any](m *Mutex[T], f func(*T) A) A {
	var ret A
	m.With(func(t *T) {
		ret = f(t)
	})
	return ret
}

// With2 is like With1, but the callback returns two values instead of one.
func With2[T, A, B any](m *Mutex[T], f func(*T) (A, B)) (A, B) {
	var (
		a A
		b B
	)
	m.With(func(t *T) {
		a, b = f(t)
	})
	return a, b
}
