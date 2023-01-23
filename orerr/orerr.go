// Package orerr provides a type T for working with (value, error)
// pairs.
package orerr

// A T[V] bundles a value of type V with a possibly-nil error. This
// allows the user to pass around a single value instead of dealing
// with tuples.
//
// Values can be extracted with Get; callers should check the error
// as with any other method that returns an error.
type T[V any] struct {
	value V
	err   error
}

// New bundles a value and possible error in a T.
func New[V any](value V, err error) T[V] {
	return T[V]{
		value: value,
		err:   err,
	}
}

// Get gets the components of the T.
func (t T[V]) Get() (V, error) {
	return t.value, t.err
}

// Err returns the error from the T.
func (t T[V]) Err() error {
	return t.err
}
