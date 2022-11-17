package bbr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueueEmpty(t *testing.T) {
	q := &queue[int]{}
	assert.True(t, q.Empty(), "The zero value is an empty queue")
	q.Push(0)
	assert.False(t, q.Empty(), "After adding a value, Empty() reports false")
	q.Pop()
	assert.True(t, q.Empty(), "After removing it, Empty() reports true again")
}

func TestQueuePeek(t *testing.T) {
	q := &queue[int]{}

	assert.Panics(t, func() { q.Peek() }, "Peek() panics on an empty queue.")

	q.Push(1)
	q.Push(2)

	assert.Equal(t, 2, q.Len(), "Len() returns the correct length")

	assert.Equal(t, 1, q.Peek(), "Peek() returns the first element in the queue")
	assert.Equal(t, 2, q.Len(), "Peek() does not change the length")
	assert.Equal(t, 1, q.Peek(), "Peek() returns the same element when called twice.")
}

func TestQueuePop(t *testing.T) {
	q := &queue[int]{}
	q.Push(1)
	q.Push(2)

	assert.Equal(t, 2, q.Len(), "initial length is correct")
	assert.Equal(t, 1, q.Pop(), "Pop() returns the first element")
	assert.Equal(t, 1, q.Len(), "length is updated")
	assert.Equal(t, 2, q.Pop(), "Pop() returns the next element")
	assert.Equal(t, 0, q.Len(), "length is updated")
	assert.True(t, q.Empty(), "Final queue is empty.")

	assert.Panics(t, func() { q.Pop() }, "Pop() panics on an emopty queue.")
}

func TestQueueItems(t *testing.T) {
	q := &queue[int]{}

	q.Push(1)
	q.Push(2)
	q.Push(3)

	concat := func(head, tail []int) []int {
		ret := make([]int, 0, len(head)+len(tail))
		return append(append(ret, head...), tail...)
	}

	assert.Equal(t, []int{1, 2, 3}, concat(q.Items()))

	q.Pop()
	q.Push(4)

	assert.Equal(t, []int{2, 3, 4}, concat(q.Items()))
}

func TestQueueFold(t *testing.T) {
	q := &queue[int]{}

	q.Push(1)
	q.Push(2)
	q.Push(3)

	assert.Equal(t, 6, q.Fold(0, func(x, y int) int {
		return x + y
	}))

	q.Pop()
	q.Push(4)

	assert.Equal(t, 9, q.Fold(0, func(x, y int) int {
		return x + y
	}))
}
