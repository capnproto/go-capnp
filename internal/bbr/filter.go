package bbr

import (
	"math"
)

// A filter estimates a value based on a sliding window of past samples.
//
// N.b. the constraint that makes the most sense here would be
// golang.org/x/exp/constraints.Ordered, but it seems silly to pull in
// the exp package given this type is un-exported and we only actually
// use it at types that happen to satisfy ~int64. TODO(cleanup): when the
// stdlib has such a constraint, use it here for clarity.
type filter[T ~int64] struct {
	buf      ringBuffer[T]
	estimate T
}

func (f *filter[T]) addSample(sample T) {
	f.buf.Shift(sample)
}

func (f *filter[T]) estimateMin() {
	f.estimate = T(math.MaxInt64)
	for _, v := range f.buf.elts {
		if v < f.estimate {
			f.estimate = v
		}
	}
}

func (f *filter[T]) estimateMax() {
	f.estimate = T(math.MinInt64)
	for _, v := range f.buf.elts {
		if v > f.estimate {
			f.estimate = v
		}
	}
}

type minFilter[T ~int64] filter[T]

func (mf *minFilter[T]) AddSample(sample T) {
	f := (*filter[T])(mf)
	f.addSample(sample)
	f.estimateMin()
}

type maxFilter[T ~int64] filter[T]

func (mf *maxFilter[T]) AddSample(sample T) {
	f := (*filter[T])(mf)
	f.addSample(sample)
	f.estimateMax()
}
