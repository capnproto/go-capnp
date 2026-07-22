package flowcontrol

import (
	"context"
	"fmt"

	"golang.org/x/sync/semaphore"
)

// Returns a FlowLimiter that enforces a fixed limit on the total size of
// outstanding messages.
func NewFixedLimiter(size int64) FlowLimiter {
	return &fixedLimiter{
		size: size,
		sem:  semaphore.NewWeighted(size),
	}
}

type fixedLimiter struct {
	size int64
	sem  *semaphore.Weighted
}

func (fl *fixedLimiter) StartMessage(ctx context.Context, size uint64) (gotResponse func(), err error) {
	if int64(size) > fl.size {
		return nil, fmt.Errorf("%w: %d bytes exceeds %d byte window", ErrMessageTooLarge, size, fl.size)
	}

	if err = fl.sem.Acquire(ctx, int64(size)); err == nil {
		gotResponse = func() { fl.sem.Release(int64(size)) }
	}

	return
}

func (fixedLimiter) Release() {}
