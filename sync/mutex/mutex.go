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

// Lock acquires the mutex and returns a reference to the value. Call Unlock()
// on the return value to release the lock.
//
// Where possible, you should prefer using Mutex.With and similar functions. If
// those are insufficient to handle your use case consider whether your locking
// scheme is too complicated before going ahead and using this.
func (m *Mutex[T]) Lock() *Locked[T] {
	m.mu.Lock()
	return &Locked[T]{mu: m}
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

// A Locked[T] is a reference to a value of type T which is guarded by a Mutex,
// which the caller has acquired.
type Locked[T any] struct {
	mu       *Mutex[T]
	unlocked bool
}

// Value returns a reference to the protected value. It must not be used after
// calling Unlock.
func (l *Locked[T]) Value() *T {
	if l.unlocked {
		panic("Called Locked.Value after Unlock.")
	}
	return &l.mu.val
}

// Unlock releases the mutex. Any references to the value obtained via Value
// must not be used after this is called.
func (l *Locked[T]) Unlock() {
	if l.unlocked {
		panic("Called Locked.Unlock twice.")
	}
	l.unlocked = true
	l.mu.mu.Unlock()
}
