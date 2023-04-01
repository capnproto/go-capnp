package rc

import "sync/atomic"

type Ref[T any] struct {
	cell *cell[T]
}

type cell[T any] struct {
	value    T
	refcount int32
	release  func()
}

func NewRef[T any](value T, release func()) *Ref[T] {
	return &Ref[T]{
		cell: &cell[T]{
			value:    value,
			refcount: 1,
			release:  release,
		},
	}
}

func (r *Ref[T]) AddRef() *Ref[T] {
	if r.cell == nil {
		panic("called AddRef() on already-released Ref.")
	}
	atomic.AddInt32(&r.cell.refcount, 1)
	return &Ref[T]{cell: r.cell}
}

func (r *Ref[T]) Release() {
	val := atomic.AddInt32(&r.cell.refcount, -1)
	if val == 0 {
		r.cell.release()
	}
	r.cell = nil
}

func (r *Ref[T]) Value() *T {
	return &r.cell.value
}
