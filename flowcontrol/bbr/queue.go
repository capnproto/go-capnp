package bbr

// queue is an un-bounded queue, backed by a ring buffer. The zero value
// is an empty queue.
type queue[T any] struct {
	buf    []T // Elements of the queue.
	start  int // Index of the head of the queue.
	length int // Number of elements in the queue.
}

// newQueue returns a new empty queue whose internal buffer's size is based
// on sizeHint.
func newQueue[T any](sizeHint int) *queue[T] {
	return &queue[T]{buf: make([]T, sizeHint)}
}

// Len returns the number of elements in the queue.
func (q *queue[T]) Len() int {
	return q.length
}

// Empty reports whether the queue is empty.
func (q *queue[T]) Empty() bool {
	return q.Len() == 0
}

// mustNotBeEmpty panics if the queue is empty.
func (q *queue[T]) mustNotBeEmpty() {
	if q.Empty() {
		panic("empty queue")
	}
}

// Peek returns the element at the head of the queue, without removing it.
// panics if the queue is empty.
func (q *queue[T]) Peek() T {
	q.mustNotBeEmpty()
	return q.buf[q.start]
}

// Pop removes and returns the element at the head of the queue.
// panics if the queue is empty.
func (q *queue[T]) Pop() T {
	ret := q.Peek()
	q.start++
	q.start %= len(q.buf)
	q.length--
	return ret
}

// Push adds the value to the end of the queue.
func (q *queue[T]) Push(v T) {
	if q.Len() == len(q.buf) {
		q.grow()
	}
	i := (q.start + q.length) % len(q.buf)
	q.buf[i] = v
	q.length++
}

// grow increases the underlying storage of the queue.
func (q *queue[T]) grow() {
	newBuf := make([]T, 0, 2*len(q.buf)+1)
	head, tail := q.Items()
	newBuf = append(append(newBuf, head...), tail...)
	newBuf = newBuf[:cap(newBuf)]
	q.buf = newBuf
	q.start = 0
}

func (q queue[T]) snapshot() queue[T] {
	ret := q
	ret.buf = make([]T, len(q.buf), cap(q.buf))
	copy(ret.buf, q.buf)
	return ret
}

// Items returns all the items in the queue, in two slices.
// The odd choice to return two slices allows this method to be zero-copy.
// The items in the first slice come before the items in the second slice.
func (q *queue[T]) Items() (head, tail []T) {
	excess := (q.start + q.length) - len(q.buf)
	if excess > 0 {
		return q.buf[q.start:], q.buf[:excess]
	} else {
		return q.buf[q.start : q.start+q.length], nil
	}
}

// Fold combines all of the values in the queue into a single value.
func (q *queue[T]) Fold(init T, combine func(acc, item T) T) T {
	return foldQueue(q, init, combine)
}

// Like queue.Fold, but doesn't require the result/init to be the same as the
// element type.
func foldQueue[A, B any](q *queue[A], init B, combine func(acc B, item A) B) B {
	acc := init
	h, t := q.Items()
	for _, v := range h {
		acc = combine(acc, v)
	}
	for _, v := range t {
		acc = combine(acc, v)
	}
	return acc
}
