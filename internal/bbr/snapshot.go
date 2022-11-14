package bbr

import (
	"fmt"
	"time"
)

type Snapshot struct {
	lim Limiter
	now time.Time
}

func SnapshotLimiter(l *Limiter) Snapshot {
	ret := *l
	now := ret.clock.Now()
	ret.btlBwFilter = l.btlBwFilter.snapshot()
	ret.rtPropFilter = l.rtPropFilter.snapshot()
	ret.state = l.state.snapshot()
	return Snapshot{
		lim: ret,
		now: now,
	}
}

type o map[string]any
type a []any

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
			"value": fmt.Sprintf("%v", lim.state),
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
