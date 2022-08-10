package bbr

// A ring buffer/circular queue. The queue always contains ringBufferSize
// elements. The zero value is a buffer containing all zero values of T.
type ringBuffer[T any] struct {
	elts [ringBufferSize]T // The elements in the ring buffer
	next int               // The index of the head of the queue
}

// ringBufferSize is the max number of elements that can be stored
// in a ringBuffer. The BBR paper suggests somewhere between 6 and
// 10. Other than that constraint, the value is chosen arbitrarily.
const ringBufferSize = 6

// Peek returns the element at the head of the queue.
func (r *ringBuffer[T]) Peek() T {
	return r.elts[r.next]
}

// Shift adds v to the end of the queue, and removes and returns
// the value at the start of the queue.
func (r *ringBuffer[T]) Shift(v T) T {
	ret := r.Peek()
	r.elts[r.next] = v
	r.next++
	r.next %= ringBufferSize
	return ret
}
