// Package mpsc implements a multiple-producer, single-consumer queue.
package mpsc

import (
	"capnproto.org/go/capnp/v3/internal/chanmutex"
	"context"
)

// A multiple-producer, single-consumer queue. Create one with New(),
// and send from many gorotuines with Tx.Send(). Only one gorotuine may
// call Rx.Recv().
type Queue struct {
	Tx
	Rx
}

// The receive end of a Queue.
type Rx struct {
	// The head of the list. If the list is empty, this will be
	// non-nil but have a locked mu field.
	head *node
}

// The send/transmit end of a Queue.
type Tx struct {
	// Mutex which must be held by senders. A goroutine must hold this
	// lock to manipulate `tail`.
	mu chanmutex.Mutex

	// Pointer to the tail of the list. This will have a locked mu,
	// and zero values for other fields.
	tail *node
}

// Alias for interface{}, the values in the queue. TODO: once Go
// supports generics, get rid of this and make Queue generic in the
// type of the values.
type Value interface{}

// A node in the linked linst that makes up the queue internally.
type node struct {
	// A mutex which guards the other fields in the node.
	// Nodes start out with this locked, and then we unlock it
	// after filling in the other fields.
	mu chanmutex.Mutex

	// The next node in the list, if any. Must be non-nil if
	// mu is unlocked:
	next *node

	// The value in this node:
	value Value
}

// Create a new node, with a locked mutex and zero values for
// the other fields.
func newNode() *node {
	return &node{mu: chanmutex.NewLocked()}
}

// Create a new, initially empty Queue.
func New() *Queue {
	node := newNode()
	return &Queue{
		Tx: Tx{
			tail: node,
			mu:   chanmutex.NewUnlocked(),
		},
		Rx: Rx{head: node},
	}
}

// Send a message on the queue.
func (tx *Tx) Send(v Value) {
	newTail := newNode()

	tx.mu.Lock()

	oldTail := tx.tail
	oldTail.next = newTail
	oldTail.value = v
	tx.tail = newTail
	oldTail.mu.Unlock()

	tx.mu.Unlock()
}

// Receive a message from the queue. Blocks if the queue is empty.
// If the context ends before the receive happens, this returns
// ctx.Err().
func (rx *Rx) Recv(ctx context.Context) (Value, error) {
	select {
	case <-rx.head.mu:
		return rx.doRecv(), nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Try to receive a message from the queue. If successful, ok will be true.
// If the queue is empty, this will return immediately with ok = false.
func (rx *Rx) TryRecv() (v Value, ok bool) {
	select {
	case <-rx.head.mu:
		return rx.doRecv(), true
	default:
		return nil, false
	}
}

// Helper for shared logic between Recv and TryRecv. Must be holding
// rx.head.mu.
func (rx *Rx) doRecv() Value {
	ret := rx.head.value
	rx.head = rx.head.next
	return ret
}
