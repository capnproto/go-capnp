package bbr

import (
	"context"
	"fmt"
	"time"
)

// A Snapshot captures the internal state of a limiter at the time it was called, along
// with the current time. This is useful for debugging.
type Snapshot struct {
	lim Limiter
	now time.Time
}

// Take a snapshot of the limiter.
func SnapshotLimiter(l *Limiter) Snapshot {
	var ret Snapshot
	l.whilePaused(func() {
		ret.lim = *l
		ret.now = l.clock.Now()
		ret.lim.btlBwFilter = l.btlBwFilter.snapshot()
		ret.lim.rtPropFilter = l.rtPropFilter.snapshot()
		ret.lim.state = l.state.snapshot()
	})
	return ret
}

// A SnapshottingLimiter is a wrapper around Limiter which takes snapshots on each
// operation and passes them to a callback.
type SnapshottingLimiter struct {
	Limiter        *Limiter       // The limiter to wrap
	RecordSnapshot func(Snapshot) // The callback to pass snapshots to
}

func (l SnapshottingLimiter) StartMessage(ctx context.Context, size uint64) (gotResponse func(), err error) {
	l.snapshot()
	r, err := l.Limiter.StartMessage(ctx, size)
	l.snapshot()
	return func() {
		l.snapshot()
		r()
		l.snapshot()
	}, err
}

func (l SnapshottingLimiter) Release() {
	l.Limiter.Release()
}

func (l SnapshottingLimiter) snapshot() {
	l.RecordSnapshot(SnapshotLimiter(l.Limiter))
}

type o map[string]any
type a []any

// Returns a value that can be formatted as json using the encoding/json package, which
// captures the data in the snapshot. The exact format of this is unspecified, as it
// exposes implementation details of the Limiter type.
func (s Snapshot) Json() any {
	lim := &s.lim
	ret := o{
		"now":             s.now,
		"cwndGain":        lim.cwndGain,
		"pacingGain":      lim.pacingGain,
		"btlBw":           lim.btlBwFilter.Estimate,
		"rtProp":          lim.rtPropFilter.Estimate,
		"nextSendTime":    lim.nextSendTime,
		"sent":            lim.sent,
		"delivered":       lim.delivered,
		"deliveredTime":   lim.deliveredTime,
		"inflight":        lim.inflight(),
		"bdp":             lim.computeBDP(),
		"appLimitedUntil": lim.appLimitedUntil,
		"state": o{
			"type":  fmt.Sprintf("%T", lim.state),
			"value": fmt.Sprintf("%v", lim.state), // TODO: better formatting of state?
		},
	}

	bwhead, bwtail := lim.btlBwFilter.q.Items()
	bwsamples := []float64{}
	for _, v := range bwhead {
		bwsamples = append(bwsamples, float64(v))
	}
	for _, v := range bwtail {
		bwsamples = append(bwsamples, float64(v))
	}

	rthead, rttail := lim.rtPropFilter.q.Items()
	rtSamples := a{}
	for _, v := range rthead {
		rtSamples = append(rtSamples, o{
			"now": v.now,
			"rtt": int64(v.rtt),
		})
	}
	for _, v := range rttail {
		rtSamples = append(rtSamples, o{
			"now": v.now,
			"rtt": int64(v.rtt),
		})
	}

	ret["samples"] = o{
		"btlBw":  bwsamples,
		"rtProp": rtSamples,
	}
	return ret
}
