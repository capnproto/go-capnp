// Package maybe provides support for working with optional values.
package maybe

// A T[V] represents an optional value of type V. The zero value for T[V]
// is considered a "missing" value.
type T[V any] struct {
	value V
	ok    bool
}

// Create a new, non-empty T with value 'value'.
func New[V any](value V) T[V] {
	return T[V]{
		value: value,
		ok:    true,
	}
}

// Get the underlying value, if any. ok is true iff the value was present.
func (t T[V]) Get() (value V, ok bool) {
	return t.value, t.ok
}
