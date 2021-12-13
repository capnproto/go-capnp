// Package flowcontrol provides support code for per-object flow control.
//
// This is most important for dealing with streaming interfaces; see
// https://capnproto.org/news/2020-04-23-capnproto-0.8.html#multi-stream-flow-control
// for a description of the general problem.
//
// The Go implementation's approach differs from that of the C++ implementation in that
// we don't treat the `stream` annotation specially; we instead do flow control on all
// objects. Calls to methods will transparently block for the appropriate amount of
// time, so it is safe to simply call rpc methods in a loop.
//
// To change the default flow control policy on a Client, call Client.SetFlowLimiter
// with the desired FlowLimiter.
package flowcontrol

import (
	"context"
	"sync"

	"capnproto.org/go/capnp/v3/internal/chanmutex"
)

// A `FlowLimiter` is used to manage flow control for a stream of messages.
type FlowLimiter interface {
	// StartMessage informs the flow limiter than the caller wants to
	// send a message of the specified size. It blocks until an appropriate
	// time to do so, or until the context is canceled. If the returned
	// error is nil, the caller should then proceed in sending the message
	// immediately, and it should arrange to call gotResponse() as soon as
	// a response is received.
	StartMessage(ctx context.Context, size uint64) (gotResponse func(), err error)
}

var (
	// A flow limiter which does not actually limit anything; messages will be
	// sent as fast as possible.
	NopLimiter FlowLimiter = nopLimiter{}
)

type nopLimiter struct{}

func (nopLimiter) StartMessage(context.Context, uint64) (func(), error) {
	return func() {}, nil
}

// Returns a FlowLimiter that enforces a fixed limit on the total size of
// outstanding messages.
func NewFixedLimiter(size uint64) FlowLimiter {
	return &fixedLimiter{
		total: size,
		avail: size,
	}
}

type fixedLimiter struct {
	mu           sync.Mutex
	total, avail uint64

	pending requestQueue
}

func (fl *fixedLimiter) StartMessage(ctx context.Context, size uint64) (gotResponse func(), err error) {
	gotResponse = fl.makeCallback(size)
	fl.mu.Lock()
	ready := fl.pending.put(size)
	fl.pumpQueue()
	fl.mu.Unlock()
	select {
	case <-ready:
		return gotResponse, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// must be holding mu
func (fl *fixedLimiter) pumpQueue() {
	next := fl.pending.peek()
	for next != nil && next.size <= fl.avail {
		fl.pending.get() // actually de-queue it.
		next.ready <- struct{}{}
		next = fl.pending.peek()
	}
}

func (fl *fixedLimiter) makeCallback(size uint64) func() {
	return func() {
		fl.mu.Lock()
		defer fl.mu.Unlock()
		fl.avail += size
		fl.pumpQueue()
	}
}

type requestQueue struct {
	head, tail *request
}

type request struct {
	ready chanmutex.Mutex
	size  uint64
	next  *request
}

func (q *requestQueue) peek() *request {
	return q.head
}

func (q *requestQueue) get() *request {
	ret := q.head
	if ret != nil {
		q.head = ret.next
	}
	if q.head == nil {
		q.tail = nil
	}
	return ret
}

func (q *requestQueue) put(size uint64) chanmutex.Mutex {
	req := &request{
		ready: chanmutex.NewUnlocked(),
		size:  size,
		next:  nil,
	}
	if q.tail == nil {
		q.tail = req
		q.head = req
	} else {
		q.tail.next = req
		q.tail = req
	}
	return req.ready
}
