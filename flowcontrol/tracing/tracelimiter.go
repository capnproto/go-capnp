// Package tracing implements support for tracing the behavior of a FlowLimiter.
package tracing

import (
	context "context"
	"sync"
	"time"

	"capnproto.org/go/capnp/v3/flowcontrol"
)

var _ flowcontrol.FlowLimiter = &TraceLimiter{}

// A TraceLimiter wraps an underlying FlowLimiter, and records data about messages.
type TraceLimiter struct {
	// The underlying FlowLimiter
	underlying flowcontrol.FlowLimiter
	mu         sync.Mutex
	records    []TraceRecord
}

// Return a new TraceLimiter, wrapping underlying.
func New(underlying flowcontrol.FlowLimiter) *TraceLimiter {
	return &TraceLimiter{
		underlying: underlying,
	}
}

// A TraceRecord records information about a message sent through the limiter.
type TraceRecord struct {
	Size       uint64    // The size of the message
	RequestAt  time.Time // The time at which StartMessage() was called.
	ProceedAt  time.Time // The time at which StartMessage() returned.
	ResponseAt time.Time // The time at which gotResponse() was called.
}

// StartMessage implements FlowLimiter.StartMessage for TraceLimiter.
func (l *TraceLimiter) StartMessage(ctx context.Context, size uint64) (gotResponse func(), err error) {
	record := TraceRecord{
		Size:      size,
		RequestAt: time.Now(),
	}
	r, err := l.underlying.StartMessage(ctx, size)
	if err != nil {
		return r, err
	}
	record.ProceedAt = time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()
	index := len(l.records)
	l.records = append(l.records, record)
	return func() {
		now := time.Now()
		r()
		l.mu.Lock()
		defer l.mu.Unlock()
		l.records[index].ResponseAt = now
	}, nil
}

// Release releases the underlying flow limiter.
func (l *TraceLimiter) Release() {
	l.underlying.Release()
}

// Records returns the records for all messages that have been sent using this limiter. It is not
// safe to use the return value while concurrently calling StartMessage() or gotResponse().
//
// If a message has been sent, but its gotResponse() callback has not yet been called, the ResponseAt
// value for that record will be the zero value.
func (l *TraceLimiter) Records() []TraceRecord {
	return l.records
}
