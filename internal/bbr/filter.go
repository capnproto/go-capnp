package bbr

import (
	"math"
	"time"
)

// A filter estimates a value based on a sliding window of past samples.
type filter[T any] struct {
	buf      ringBuffer[T]
	estimate T
}

// An rtPropFilter estimates the round-trip propagation time, by taking
// the minimum round trip time of the samples in its window.
type rtPropFilter filter[time.Duration]

func (f *rtPropFilter) AddSample(rtt time.Duration) {
	f.buf.Shift(rtt)

	// Update the estimate (minimum)
	f.estimate = time.Duration(math.MaxInt64)
	for _, v := range f.buf.elts {
		if v < f.estimate {
			f.estimate = v
		}
	}
}

// A btlBwFilter estimates the bottleneck bandwidth, as the maximum delivery
// rate over the samples in its window
// (where delivery rate = data deliverd/time elapsed).
type btlBwFilter filter[int64]

func (f *btlBwFilter) AddSample(s sample) {
	f.buf.Shift(s.DeliveryRate())

	// Update ht estimate (maximum)
	f.estimate = math.MinInt64
	for _, v := range f.buf.elts {
		if v > f.estimate {
			f.estimate = v
		}
	}
}
