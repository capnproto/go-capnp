package flowcontrol

import (
	"context"
	"sync"

	"capnproto.org/go/capnp/v3/internal/chanmutex"
)

type FlowLimiter interface {
	StartMessage(ctx context.Context, size uint64) (gotResponse func(), err error)
}

var (
	NopLimiter FlowLimiter = nopLimiter{}
)

type nopLimiter struct{}

func (nopLimiter) StartMessage(context.Context, uint64) (func(), error) {
	return func() {}, nil
}

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
