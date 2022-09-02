package bbr

import (
	"context"
	"testing"
	"time"

	"capnproto.org/go/capnp/v3/internal/clock"
	"github.com/stretchr/testify/assert"
)

type snapshot struct {
	lim Limiter
	now time.Time
}

func (s *snapshot) report(t *testing.T) {
	lim := &s.lim
	t.Logf("Limiter snapshot at %v: \n\n", s.now)
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
}

func withinTolerance(t *testing.T, actual, expected, tolerance float64) {
	min := expected * (1 - tolerance)
	max := expected * (1 + tolerance)
	assert.Greater(t, actual, min)
	assert.Less(t, actual, max)

}

// trueValues returns the true/expected values of rtProp and and btlBw for a given path and
// minimum packet size.
func trueValues(path []testLink, minPacketBytes uint64) (rtProp time.Duration, btlBw bytesPerNs) {
	for _, link := range path {
		rtProp += link.delay
		if link.bandwidth > btlBw {
			btlBw = link.bandwidth
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
		path           []testLink
		minPacketBytes uint64
		rtProp         time.Duration
		btlBw          bytesPerNs
	}{
		{
			path: []testLink{
				{delay: 50 * time.Millisecond},
				{bandwidth: 1000 * bytesPerSecond},
			},
			minPacketBytes: 10,
			rtProp:         60 * time.Millisecond,
			btlBw:          1000 * bytesPerSecond,
		},
		{
			path: []testLink{
				{delay: 100 * time.Millisecond},
				{bandwidth: 1000 * bytesPerSecond},
			},
			minPacketBytes: 10,
			rtProp:         110 * time.Millisecond,
			btlBw:          1000 * bytesPerSecond,
		},
		{
			path: []testLink{
				{delay: 50 * time.Millisecond},
				{bandwidth: 1000 * bytesPerSecond},
			},
			minPacketBytes: 5,
			rtProp:         55 * time.Millisecond,
			btlBw:          1000 * bytesPerSecond,
		},
	}
	for _, c := range cases {
		rtProp, btlBw := trueValues(c.path, c.minPacketBytes)
		withinTolerance(t, float64(c.rtProp), float64(rtProp), 0.01)
		withinTolerance(t, float64(c.btlBw), float64(btlBw), 0.01)
	}
}

// TODO: work this into a proper test of... something...
//
// Right now this just gives me a decent way of watching what happens when
// we push data into a stream continuously, but it's not really measuring
// any pass/fail criteria.
func TestTrace(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping long test due to -short")
	}
	sizes := []uint64{1, 5, 7, 10, 20}
	var packets []uint64

	for i := 0; i < 100; i++ {
		packets = append(packets, sizes...)
	}

	path := []testLink{
		{delay: 50 * time.Millisecond},
		{bandwidth: 1000 * bytesPerSecond},
	}

	snapshots := runTrace(path, packets)

	for _, s := range snapshots {
		s.report(t)
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
