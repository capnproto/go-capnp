package bbr

import (
	"math"
	"time"
)

// Size of the bottleneck bandwidth filter's window. The paper suggests
// 6-10, other than that this is arbitrary.
const btlBwFilterSize = 6

// Filter that estimates the bottleneck bandwidth.
type btlBwFilter struct {
	q        queue[int64]
	Estimate int64
}

func newBtlBwFilter() btlBwFilter {
	return btlBwFilter{
		q: *newQueue[int64](btlBwFilterSize),

		// We set this to something that is non-zero,
		// but otherwise won't take precedence over any actual
		// data we receive:
		Estimate: 1,
	}
}

func (f *btlBwFilter) AddSample(deliveryRate int64) {
	if f.q.Len() == btlBwFilterSize {
		f.q.Pop()
	}
	f.q.Push(deliveryRate)
	f.Estimate = f.q.Fold(0, max[int64])
}

// Filter that estimates the round-trip propagation time.
type rtPropFilter struct {
	q          queue[rtPropSample]
	nextSample rtPropSample
	Estimate   time.Duration
}

func newRtPropFilter() rtPropFilter {
	return rtPropFilter{
		nextSample: rtPropSample{
			// Set this to a value that will be removed immediately,
			// and will be lower priority than any other value, so
			// it will removed from the queue when a new sample is
			// added.
			//
			// XXX: this still could theoretically do the wrong thing
			// if somebody's clock is set very wrong. We should find
			// a better way to do this; maybe ask the user to supply
			// an initial estimate, before any samples are collected?
			now: time.Unix(math.MinInt64, 0),
			rtt: math.MaxInt64,
		},
	}
}

type rtPropSample struct {
	rtt time.Duration
	now time.Time
}

func (f *rtPropFilter) AddSample(sample rtPropSample) {
	// We want to avoid an un-bounded growing queue for two reasons:
	//
	// 1. Space usage
	// 2. Recomputing the estimate is an O(n) operation, so we want to
	//    keep n small if we're doing that on each ack.
	//
	// We manage this as follows: rather than adding each sample to the queue
	// individually, we coalesce all samples within a given 1 second interval
	// into a single slot in the queue, taking their minimum. Since we drop
	// samples that are more than 30 seconds old, this bounds the queue to 30
	// elements.
	//
	// This gives enough granularity to get the benefits of the sliding window
	// without needing to store each and every sample.

	if sample.now.Sub(f.nextSample.now) > time.Second {
		f.q.Push(f.nextSample)
		f.nextSample = sample
	} else {
		f.nextSample = rtPropSample{
			// We keep the old `now`, since it's used to determine if we
			// should shift to the next sample:
			now: f.nextSample.now,
			rtt: min(f.nextSample.rtt, sample.rtt),
		}
	}

	// Clear out any samples older than 30 seconds:
	for !f.q.Empty() && sample.now.Sub(f.q.Peek().now) > 30*time.Second {
		f.q.Pop()
	}

	f.Estimate = foldQueue(
		&f.q,
		f.nextSample.rtt,
		func(rtProp time.Duration, sample rtPropSample) time.Duration {
			return min(rtProp, sample.rtt)
		},
	)
}

// min and max compute the minimum and maximum of two numbers, respectively.
// Presumably, as Go generics become more widely used, these will be dropped
// in favor of some standard library function.
//
// We could use a broader constraint here (any numeric type), but we'd have
// to either define the alias ourselves or import the exp package, and we
// only actually use these at types covered by ~int64 anyway.

func min[T ~int64](x, y T) T {
	if x < y {
		return x
	} else {
		return y
	}
}

func max[T ~int64](x, y T) T {
	if x > y {
		return x
	} else {
		return y
	}
}
