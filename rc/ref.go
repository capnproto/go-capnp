// Package rc provides reference-counted cells.
package rc

import "sync/atomic"

// A Ref is a reference to a refcounted cell containing a T.
// It must not be moved after it is first used (but it is ok
// to move the return values of NewRef/AddRef *before* using
// them). The zero value is an already-released reference to
// some non-existant cell; it is not useful.
type Ref[T any] struct {
	// Pointer to the actual data. When we release our
	// reference, we set this to nil so we can't access
	// it.
	cell *cell[T]
}

// Container for the actual data; Ref just points to this.
type cell[T any] struct {
	value    T      // The actual value that is stored.
	refcount int32  // The refernce count.
	release  func() // Function to call when refcount hits zero.
}

// NewRef returns a Ref pointing to value. When all references
// are released, the function release will be called.
func NewRef[T any](value T, release func()) *Ref[T] {
	return &Ref[T]{
		cell: &cell[T]{
			value:    value,
			refcount: 1,
			release:  release,
		},
	}
}

// AddRef returns a new reference to the same underlying data as
// the receiver. The references are not interchangable: to
// release the underlying data you must call Release on each
// Ref separately, and you cannot access the value through
// a released Ref even if you know there are other live references
// to it.
//
// Panics if this reference has already been released.
func (r *Ref[T]) AddRef() *Ref[T] {
	if r.cell == nil {
		panic("called AddRef() on already-released Ref.")
	}
	atomic.AddInt32(&r.cell.refcount, 1)
	return &Ref[T]{cell: r.cell}
}

// Release this reference to the value. If this is the last reference,
// this calls the release function that was passed to NewRef.
func (r *Ref[T]) Release() {
	val := atomic.AddInt32(&r.cell.refcount, -1)
	if val == 0 {
		r.cell.release()
	}
	r.cell = nil
}

// Return a pointer to the value. Panics if the reference has already
// been released.
func (r *Ref[T]) Value() *T {
	return &r.cell.value
}
