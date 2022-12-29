package bbr

import (
	"fmt"
	"math"
	"time"
)

// Size of the bottleneck bandwidth filter's window. The paper suggests
// 6-10, other than that this is arbitrary.
const btlBwFilterSize = 6

// Units in which we measure the bottleneck bandwidth. Also equivalent
// to GB/s.
type bytesPerNs float64

// Filter that estimates the bottleneck bandwidth.
type btlBwFilter struct {
	q        queue[bytesPerNs]
	Estimate bytesPerNs
}

func newBtlBwFilter() btlBwFilter {
	return btlBwFilter{
		q: *newQueue[bytesPerNs](btlBwFilterSize),

		// We set this to something that is only barely
		// non-zero, so it won't result in divide by
		// zero errors but also won't take precedence
		// over any actual data we receive.
		Estimate: 1e-10, // 1 byte per 10s
	}
}

func (f *btlBwFilter) AddSample(deliveryRate bytesPerNs) {
	if f.q.Len() == btlBwFilterSize {
		f.q.Pop()
	}
	f.q.Push(deliveryRate)
	f.Estimate = f.q.Fold(0, max[bytesPerNs])
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
			// Set this to a value that will immediately be superseded
			// as soon as we get a real sample.
			now: time.Unix(math.MinInt64, 0),
			rtt: math.MaxInt64,
		},
	}
}

type rtPropSample struct {
	rtt time.Duration
	now time.Time
}

func (s rtPropSample) String() string {
	return fmt.Sprintf("rtPropSample{rtt: %v, now = %v}", s.rtt, s.now)
}

func (f *rtPropFilter) AddSample(sample rtPropSample) {
	// We want to avoid an un-bounded growing queue for two reasons:
	//
	// 1. Space usage
	// 2. Recomputing the estimate is an O(n) operation, so we want to
	//    keep n small if we're doing that on each ack.
	//
	// We manage this as follows: rather than adding each sample to the
	// queue individually, we compute the minimum RTT for each 1-second
	// window, and add that aggregate value to the queue, which in turn
	// drops samples more than 30 seconds old. This gives the benefits of
	// the sliding window without needing to store each and every sample.

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

func (f btlBwFilter) snapshot() btlBwFilter {
	ret := f
	ret.q = f.q.snapshot()
	return ret
}

func (f rtPropFilter) snapshot() rtPropFilter {
	ret := f
	ret.q = f.q.snapshot()
	return ret
}

// min and max compute the minimum and maximum of two numbers, respectively.
// Presumably, as Go generics become more widely used, these will be dropped
// in favor of some standard library function.
//
// We could use a broader constraint here (constraints.Ordered), but we'd have
// to either define the alias ourselves or import the exp package, and we
// only actually use these at more specific types anyway.

func min[T ~int64](x, y T) T {
	if x < y {
		return x
	} else {
		return y
	}
}

func max[T ~float64](x, y T) T {
	if x > y {
		return x
	} else {
		return y
	}
}
