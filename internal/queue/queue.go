// Package queue implements a generic queue using a ring buffer.
package queue

// A Queue wraps an Interface to provide queue operations.
type Queue struct {
	q     Interface
	start int
	n     int
}

// New creates a new queue that starts with n elements.
func New(q Interface, n int) *Queue {
	qq := new(Queue)
	qq.Init(q, n)
	return qq
}

// Init initializes a queue.  The old queue is untouched.
func (q *Queue) Init(r Interface, n int) {
	q.q = r
	q.start, q.n = 0, n
}

// Len returns the length of the queue.  This is different from the
// underlying interface's length.
func (q *Queue) Len() int {
	return q.n
}

// Push pushes an element on the queue.  If the queue is full,
// Push returns false.  If x is nil, Push panics.
func (q *Queue) Push(x interface{}) bool {
	n := q.q.Len()
	if q.n >= n {
		return false
	}
	i := (q.start + q.n) % n
	q.q.Set(i, x)
	q.n++
	return true
}

// Peek returns the element at the front of the queue.
// If the queue is empty, Peek panics.
func (q *Queue) Peek() interface{} {
	if q.n == 0 {
		panic("Queue.Pop called on empty queue")
	}
	return q.q.At(q.start)
}

// Pop pops an element from the queue.
// If the queue is empty, Pop panics.
func (q *Queue) Pop() interface{} {
	x := q.Peek()
	q.q.Set(q.start, nil)
	q.start = (q.start + 1) % q.q.Len()
	q.n--
	return x
}

// A type implementing Interface can be used to store elements in a Queue.
type Interface interface {
	// Len returns the number of elements available.
	Len() int
	// At returns the element at i.
	At(i int) interface{}
	// Set sets the element at i to x.
	// If x is nil, that element should be cleared.
	Set(i int, x interface{})
}
