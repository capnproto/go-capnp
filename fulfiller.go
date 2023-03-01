package capnp

// A Fulfiller supplies a value for a pending promise.
type Fulfiller[T any] interface {
	// Fulfill supplies the value for the corresponding
	// Promise
	Fulfill(T)

	// Reject rejects the corresponding promise, with
	// the specified error.
	Reject(error)
}
