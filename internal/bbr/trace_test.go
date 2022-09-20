package bbr

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"testing"
	"time"

	"capnproto.org/go/capnp/v3/internal/clock"
)

type snapshot struct {
	lim Limiter
	now time.Time
}

type o map[string]any
type a []any

func (s *snapshot) reportJson() any {
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

func (s *snapshot) report(t *testing.T) {
	lim := &s.lim
	t.Logf("Limiter snapshot at %v: \n", s.now)
	t.Logf("cwndGain        = %v\n", lim.cwndGain)
	t.Logf("pacingGain      = %v\n", lim.pacingGain)
	t.Logf("btlBw           = %v\n", lim.btlBwFilter.Estimate)
	t.Logf("rtProp          = %v\n", lim.rtPropFilter.Estimate)
	t.Logf("nextSendTime    = %v\n", lim.nextSendTime)
	t.Logf("sent            = %v\n", lim.sent)
	t.Logf("delivered       = %v\n", lim.delivered)
	t.Logf("deliveredTime   = %v\n", lim.deliveredTime)
	t.Logf("inflight        = %v\n", lim.inflight())
	t.Logf("bdp             = %v\n", lim.computeBDP())
	t.Logf("appLimitedUntil = %v\n", lim.appLimitedUntil)
	t.Logf("state           = %T%v\n", lim.state, lim.state)

	t.Logf("btlBw samples:\n")
	bwhead, bwtail := lim.btlBwFilter.q.Items()
	for _, v := range bwhead {
		t.Logf("  %v\n", v)
	}
	for _, v := range bwtail {
		t.Logf("  %v\n", v)
	}

	t.Logf("rtProp samples:\n")
	rthead, rttail := lim.rtPropFilter.q.Items()
	for _, v := range rthead {
		t.Logf("  %v\n", v)
	}
	for _, v := range rttail {
		t.Logf("  %v\n", v)
	}

	t.Logf("\n\n")
}

func withinTolerance(actual, expected, tolerance float64, msg string) error {
	min := expected * (1 - tolerance)
	max := expected * (1 + tolerance)
	if actual < min {
		return fmt.Errorf("%v = %v not within tolerance (min: %v)", msg, actual, min)
	}
	if actual > max {
		return fmt.Errorf("%v = %v not within tolerance (max: %v)", msg, actual, max)
	}
	return nil
}

// trueValues returns the true/expected values of rtProp and and btlBw for a given path and
// set of packets
func trueValues(path []testLink, packetSizes []uint64) (rtProp time.Duration, btlBw bytesPerNs) {
	var totalData uint64
	var minPacketBytes uint64 = math.MaxUint64

	for _, size := range packetSizes {
		totalData += size
		if size < minPacketBytes {
			minPacketBytes = size
		}
	}

	avgPacketSize := totalData / uint64(len(packetSizes))

	for _, link := range path {
		rtProp += link.delay
		if link.bandwidth > btlBw {
			btlBw = link.bandwidth
		}

		if link.delay > 0 {
			delayBw := bytesPerNs(avgPacketSize) / bytesPerNs(link.delay.Nanoseconds())
			if delayBw > btlBw {
				btlBw = delayBw
			}
		}
	}

	// The bottleneck only introduces meaningful delays when bandwidth limited, so the
	// minimum such delay will be however long it takes to transfer the smallest packet:
	btlProp := float64(minPacketBytes) / float64(btlBw)
	rtProp += time.Duration(btlProp)
	return
}

func TestTrueValues(t *testing.T) {
	var cases = []struct {
		path        []testLink
		packetSizes []uint64
		rtProp      time.Duration
		btlBw       bytesPerNs
	}{
		{
			path: []testLink{
				{delay: 50 * time.Millisecond},
				{bandwidth: 1000 * bytesPerSecond},
			},
			packetSizes: []uint64{10},
			rtProp:      60 * time.Millisecond,
			btlBw:       1000 * bytesPerSecond,
		},
		{
			path: []testLink{
				{delay: 100 * time.Millisecond},
				{bandwidth: 1000 * bytesPerSecond},
			},
			packetSizes: []uint64{10},
			rtProp:      110 * time.Millisecond,
			btlBw:       1000 * bytesPerSecond,
		},
		{
			path: []testLink{
				{delay: 50 * time.Millisecond},
				{bandwidth: 1000 * bytesPerSecond},
			},
			packetSizes: []uint64{5},
			rtProp:      55 * time.Millisecond,
			btlBw:       1000 * bytesPerSecond,
		},
	}
	for _, c := range cases {
		rtProp, btlBw := trueValues(c.path, c.packetSizes)
		err := withinTolerance(float64(c.rtProp), float64(rtProp), 0.01, "rtProp")
		if err != nil {
			t.Error(err)
		}
		err = withinTolerance(float64(c.btlBw), float64(btlBw), 0.01, "btlBw")
		if err != nil {
			t.Error(err)
		}
	}
}

// estimatesCorrect checks that the snapshot's estimates are close to the true values. If not,
// the error describes the descrepancy.
func estimatesCorrect(path []testLink, packetSizes []uint64, snapshot snapshot) error {
	rtProp, btlBw := trueValues(path, packetSizes)
	estRtProp := snapshot.lim.rtPropFilter.Estimate
	estBtlBw := snapshot.lim.btlBwFilter.Estimate
	errRtProp := withinTolerance(float64(estRtProp), float64(rtProp), 0.05, "rtProp")
	errBtlBw := withinTolerance(float64(estBtlBw), float64(btlBw), 0.05, "btlBw")
	if errRtProp == nil {
		return errBtlBw
	}
	if errBtlBw == nil {
		return errRtProp
	}
	return fmt.Errorf("Estimates are incorrect:\n%v\n%v\n", errRtProp, errBtlBw)
}

func repeat[T any](count int, values []T) []T {
	var ret []T
	for i := 0; i < count; i++ {
		ret = append(ret, values...)
	}
	return ret
}

type ceArg struct {
	btlBw          bytesPerNs
	path           []testLink
	minPacketBytes uint64
}

func computeErr(expected, actual float64) float64 {
	return (actual - expected) / expected
}

func TestGraphs(t *testing.T) {
	gatherData(t)
}

func gatherData(t *testing.T) {
	delays := []time.Duration{}
	bandwidths := []bytesPerNs{}
	packets := [][]uint64{
		repeat(10, []uint64{1, 49, 50, 50, 50}),
	}

	for i := 0; i < 10; i++ {
		delays = append(delays, time.Duration(i+1)*5*time.Millisecond)
	}
	for i := 0; i < 10; i++ {
		delays = append(delays, time.Duration(i+1)*50*time.Millisecond)
	}
	for i := 0; i < 50; i++ {
		bw := 100 * (i + 1)
		bandwidths = append(bandwidths, bytesPerNs(bw)*bytesPerSecond)
	}

	data := []any{}

	for i, d := range delays {
		for j, b := range bandwidths {
			path := []testLink{
				{delay: d},
				{bandwidth: b},
			}
			for k, p := range packets {
				t.Logf("Case (%v,%v,%v): (%v, %v, %v)", i, j, k, d, b, p)
				rtProp, btlBw := trueValues(path, p)
				snapshots := runTrace(path, p)

				sample := o{
					"delay":     int64(d),
					"bandwidth": float64(b),
					"packets":   p,
					"trueValues": o{
						"rtProp": int64(rtProp),
						"btlBw":  float64(btlBw),
					},
				}
				trace := a{}
				for _, s := range snapshots {
					trace = append(trace, s.reportJson())
				}
				sample["trace"] = trace
				data = append(data, sample)
			}
		}
	}

	buf, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile("data.json", buf, 0600)
	if err != nil {
		panic(err)
	}
}

// Collect traces of various packet sequences being streamed over various paths,
// and check that they pass some sanity checks.
func TestTrace(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long test due to -short")
	}
	t.Parallel()

	cases := []struct {
		path    []testLink
		packets []uint64
	}{
		{
			path: []testLink{
				{delay: 50 * time.Millisecond},
				{bandwidth: 1000 * bytesPerSecond},
			},
			packets: repeat(10, []uint64{1, 49, 50, 50, 50}),
		},
		{
			path: []testLink{
				{delay: 50 * time.Millisecond},
				{bandwidth: 1e6 * bytesPerSecond},
			},
			packets: repeat(20, []uint64{1, 4900, 5000, 5000, 5000}),
		},
		{
			path: []testLink{
				{delay: 50 * time.Millisecond},
				{bandwidth: 1e6 * bytesPerSecond},
			},
			packets: repeat(100, []uint64{1, 49, 50, 50, 50}),
		},
		{
			path: []testLink{
				{delay: 5 * time.Millisecond},
				{bandwidth: 1e6 * bytesPerSecond},
			},
			packets: repeat(100, []uint64{1, 49, 50, 50, 50}),
		},
		{
			path: []testLink{
				{delay: 5 * time.Millisecond},
				{bandwidth: 1e6 * bytesPerSecond},
			},
			packets: repeat(100, []uint64{1, 49, 50, 50, 50}),
		},
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Case %v", i), func(t *testing.T) {
			snapshots := runTrace(c.path, c.packets)
			for _, s := range snapshots {
				s.report(t)
			}
			t.Run("At end", func(t *testing.T) {
				err := estimatesCorrect(c.path, c.packets, snapshots[len(snapshots)-1])
				if err != nil {
					t.Error(err)
				}
			})
			t.Run("After startup", func(t *testing.T) {
				for _, s := range snapshots {
					// Find the first snapshot after startup.
					if _, ok := s.lim.state.(*drainState); ok {
						err := estimatesCorrect(c.path, c.packets, s)
						if err != nil {
							t.Error(err)
						}
						return
					}
				}
			})
			t.Run("At some point", func(t *testing.T) {
				for _, s := range snapshots {
					err := estimatesCorrect(c.path, c.packets, s)
					if err == nil {
						return
					}
				}
				t.Fatal("Estimates are never correct.")
			})
		})
	}
}

func (l *Limiter) snapshot() Limiter {
	ret := *l
	ret.btlBwFilter = l.btlBwFilter.snapshot()
	ret.rtPropFilter = l.rtPropFilter.snapshot()
	ret.state = l.state.snapshot()
	return ret
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

func (q queue[T]) snapshot() queue[T] {
	ret := q
	ret.buf = make([]T, len(q.buf), cap(q.buf))
	copy(ret.buf, q.buf)
	return ret
}

func runTrace(path []testLink, packetSizes []uint64) []snapshot {
	var results []snapshot
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	clock := clock.System
	lim := NewLimiter(clock)
	defer lim.Release()

	tx, rx := startTestPath(ctx, clock, path...)
	go processAcks(ctx, rx)

	for _, size := range packetSizes {
		got, err := lim.StartMessage(ctx, size)
		if err != nil {
			panic(err)
		}
		tx.Send(testPacket{
			size:        size,
			gotResponse: got,
		})
		lim.whilePaused(func() {
			results = append(results, snapshot{
				lim: lim.snapshot(),
				now: clock.Now(),
			})
		})
	}
	return results
}
