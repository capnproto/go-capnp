// Package tracing implements support for tracing the behavior of a FlowLimiter.
package tracing

import (
	context "context"
	"time"

	"capnproto.org/go/capnp/v3/flowcontrol"
)

var _ flowcontrol.FlowLimiter = &TraceLimiter{}

// A TraceLimiter wraps an underlying FlowLimiter, and records data about messages.
type TraceLimiter struct {
	// The underlying FlowLimiter
	underlying flowcontrol.FlowLimiter
	emitRecord func(TraceRecord)
}

// Return a new TraceLimiter, wrapping underlying. Each time one of the limiter's gotResponse()
// callbacks is invoked, emitRecord is called with data about the call.
func New(underlying flowcontrol.FlowLimiter, emitRecord func(TraceRecord)) *TraceLimiter {
	return &TraceLimiter{
		underlying: underlying,
		emitRecord: emitRecord,
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
	return func() {
		now := time.Now()
		r()
		record.ResponseAt = now
		l.emitRecord(record)
	}, nil
}

// Release releases the underlying flow limiter.
func (l *TraceLimiter) Release() {
	l.underlying.Release()
}
