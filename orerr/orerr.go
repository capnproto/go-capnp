// Package orerr provides a type OrErr for working with (value, error)
// pairs.
package orerr

// An OrErr[V] bundles a value of type V with a possibly-nil error. This
// allows the user to pass around a single value instead of dealing
// with tuples.
//
// Values can be extracted with Get; callers should check the error
// as with any other method that returns an error.
type OrErr[V any] struct {
	value V
	err   error
}

// New bundles a value and possible error in an OrErr.
func New[V any](value V, err error) OrErr[V] {
	return OrErr[V]{
		value: value,
		err:   err,
	}
}

// Get gets the components of the T.
func (t OrErr[V]) Get() (V, error) {
	return t.value, t.err
}

// Err returns the error from the OrErr.
func (t OrErr[V]) Err() error {
	return t.err
}
