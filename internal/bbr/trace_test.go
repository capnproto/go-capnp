package bbr

import (
	"context"
	"testing"
	"time"

	"capnproto.org/go/capnp/v3/internal/clock"
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
	head, tail := lim.btlBwFilter.q.Items()
	for _, v := range head {
		t.Logf("  %v\n", v)
	}
	for _, v := range tail {
		t.Logf("  %v\n", v)
	}
}

// TODO: work this into a proper test of... something...
//
// Right now this just gives me a decent way of watching what happens when
// we push data into a stream continuously, but it's not really measuring
// any pass/fail criteria.
func TestTrace(t *testing.T) {
	sizes := []uint64{1, 5, 7, 10, 20}
	var packets []uint64

	for i := 0; i < 100; i++ {
		packets = append(packets, sizes...)
	}

	path := []testLink{
		//{delay: 4 * time.Second},
		{bandwidth: 1000 * bytesPerSecond},
		//{delay: 2 * time.Second},
	}

	snapshots := runTrace(path, packets)

	for _, s := range snapshots {
		s.report(t)
	}
}

func runTrace(path []testLink, packetSizes []uint64) []snapshot {
	var results []snapshot
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//	clock := clock.NewManual(sampleStartTime)
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
				lim: *lim,
				now: clock.Now(),
			})
		})
	}
	return results
}
