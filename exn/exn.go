// Package exn provides an exception-like mechanism.
//
// It is a bit easier to reason about than standard exceptions, since
// not any code can throw, only code given access to a throw callback.
// See the docs for Try.
package exn

// exn is a wrapper type used to distinguish values this package passes
// to panic() from anything else that might be passed to panic by
// unrelated code.
type exn struct {
	err error
}

// An alias for the type of the callback provided by Try.
type Thrower = func(error)

// Try invokes f, which is a callback with type:
//
// func(throw Thrower) T
//
// Thrower is an alias for func(error).
//
// If f returns normally, Try returns the value f returned and a nil
// error.
//
// If f invokes throw with a non-nil error, f's execution is
// terminated, and Try returns the zero value for T and the
// error that was passed to throw.
//
// If f invokes throw with a nil argument, throw returns normally
// without terminating f. This is so that error handling code can
// pass errors to throw unconditionally, like:
//
// 	v, err := foo()
// 	throw(err)
//
// f must not store throw or otherwise cause it to be invoked after
// Try returns.
func Try[T any](f func(Thrower) T) (result T, err error) {
	throw := func(e error) {
		if e == nil {
			return
		}
		panic(exn{err: e})
	}
	finishedCall := false
	defer func() {
		if finishedCall {
			return
		}

		panicVal := recover()
		if e, ok := panicVal.(exn); ok {
			err = e.err
		} else {
			panic(panicVal)
		}
	}()
	result = f(throw)
	finishedCall = true
	return
}

// Try0 is like Try, but f does not return a value, and Try0 only returns
// an error.
func Try0(f func(Thrower)) error {
	_, err := Try(func(throw func(error)) struct{} {
		f(throw)
		return struct{}{}
	})
	return err
}

// Try2 is like Try, but with two values instead of 1.
func Try2[A, B any](f func(Thrower) (A, B)) (A, B, error) {
	r, err := Try(func(throw func(error)) pair[A, B] {
		a, b := f(throw)
		return pair[A, B]{a: a, b: b}
	})
	return r.a, r.b, err
}

// TODO: move into its own package.
type pair[A, B any] struct {
	a A
	b B
}
