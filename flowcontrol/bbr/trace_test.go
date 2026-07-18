package bbr

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"testing"
	"time"
)

func (s Snapshot) report(t *testing.T) {
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
	diff := actual - expected
	if math.Abs(diff) <= math.Abs(expected*tolerance) {
		return nil
	}

	return fmt.Errorf(
		"%v = %v not within tolerance of %v of expected value %v (min %v, max %v). Actual error is %v",
		msg, actual, tolerance, expected, min, max, diff/expected)
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
		if link.bandwidth > 0 {
			if link.bandwidth < btlBw || btlBw == 0 {
				btlBw = link.bandwidth
			}
			// Bandwidth limited links only introduce meaningful delays when bandwidth limited, so the
			// minimum such delay will be however long it takes to transfer the smallest packet:
			btlProp := float64(minPacketBytes) / float64(link.bandwidth)
			rtProp += time.Duration(btlProp)
		}

		if link.delay > 0 {
			delayBw := bytesPerNs(avgPacketSize) / bytesPerNs(link.delay.Nanoseconds())
			if delayBw < btlBw || btlBw == 0 {
				btlBw = delayBw
			}
		}
	}

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
				{delay: 5 * time.Millisecond},
				{bandwidth: 1000 * bytesPerSecond},
			},
			packetSizes: []uint64{10},
			rtProp:      15 * time.Millisecond,
			btlBw:       1000 * bytesPerSecond,
		},
		{
			path: []testLink{
				{delay: 50 * time.Millisecond},
				{bandwidth: 1000 * bytesPerSecond},
			},
			packetSizes: []uint64{10},
			rtProp:      60 * time.Millisecond,
			btlBw:       200 * bytesPerSecond,
		},
		{
			path: []testLink{
				{delay: 100 * time.Millisecond},
				{bandwidth: 1000 * bytesPerSecond},
			},
			packetSizes: []uint64{10},
			rtProp:      110 * time.Millisecond,
			btlBw:       100 * bytesPerSecond,
		},
		{
			path: []testLink{
				{delay: 50 * time.Millisecond},
				{bandwidth: 1000 * bytesPerSecond},
			},
			packetSizes: []uint64{5},
			rtProp:      55 * time.Millisecond,
			btlBw:       100 * bytesPerSecond,
		},
	}
	for i, c := range cases {
		t.Run(fmt.Sprintf("Case %v", i), func(t *testing.T) {
			rtProp, btlBw := trueValues(c.path, c.packetSizes)
			err := withinTolerance(float64(rtProp), float64(c.rtProp), 0.01, "rtProp")
			if err != nil {
				t.Error(err)
			}
			err = withinTolerance(float64(btlBw), float64(c.btlBw), 0.01, "btlBw")
			if err != nil {
				t.Error(err)
			}
		})
	}
}

type tolerances struct {
	rtProp, btlBw float64
}

// estimatesCorrect checks that the snapshot's estimates are close to the true values. If not,
// the error describes the descrepancy.
func estimatesCorrect(path []testLink, packetSizes []uint64, tolerances tolerances, snapshot Snapshot) error {
	rtProp, btlBw := trueValues(path, packetSizes)
	estRtProp := snapshot.lim.rtPropFilter.Estimate
	estBtlBw := snapshot.lim.btlBwFilter.Estimate
	errRtProp := withinTolerance(float64(estRtProp), float64(rtProp), tolerances.rtProp, "rtProp")
	errBtlBw := withinTolerance(float64(estBtlBw), float64(btlBw), tolerances.btlBw, "btlBw")
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
	if os.Getenv("GATHER_DATA") != "1" {
		t.Skip("Not generating data for graphs.")
	}
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
					trace = append(trace, s.Json())
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
	cases := []struct {
		name             string
		path             []testLink
		packets          []uint64
		checkConvergence bool
	}{
		{
			name: "low bandwidth",
			path: []testLink{
				{delay: 50 * time.Millisecond},
				{bandwidth: 1000 * bytesPerSecond},
			},
			packets:          repeat(40, []uint64{50}),
			checkConvergence: false,
		},
		{
			name: "high bandwidth",
			path: []testLink{
				{delay: 50 * time.Millisecond},
				{bandwidth: 1e6 * bytesPerSecond},
			},
			packets:          repeat(20, []uint64{1, 4900, 5000, 5000, 5000}),
			checkConvergence: true,
		},
	}

	tolerances := tolerances{
		rtProp: 0.25,
		btlBw:  0.15,
	}

	for i, c := range cases {
		t.Run(fmt.Sprintf("Case %v: %s", i, c.name), func(t *testing.T) {
			snapshots := runTrace(c.path, c.packets)
			assertTraceTransitions(t, snapshots)
			if c.checkConvergence {
				assertTraceConverges(t, c.path, c.packets, tolerances, snapshots)
			}
			assertTraceRepeatable(t, c.path, c.packets, snapshots)
		})
	}
}

func assertTraceTransitions(t *testing.T, snapshots []Snapshot) {
	t.Helper()
	states := make([]string, 0, len(snapshots))
	for _, s := range snapshots {
		switch s.lim.state.(type) {
		case *startupState:
			states = append(states, "startup")
		case *drainState:
			states = append(states, "drain")
		case *probeBWState:
			states = append(states, "probeBW")
		}
	}
	want := []string{"startup", "drain", "probeBW"}
	pos := 0
	for _, state := range states {
		if state == want[pos] {
			pos++
			if pos == len(want) {
				return
			}
		}
	}
	t.Fatalf("state transition sequence %v does not contain %v", states, want)
}

func assertTraceConverges(t *testing.T, path []testLink, packets []uint64, tolerances tolerances, snapshots []Snapshot) {
	t.Helper()
	for _, s := range snapshots {
		if err := estimatesCorrect(path, packets, tolerances, s); err == nil {
			return
		}
	}
	t.Fatal("estimates never converged within tolerance")
}

func assertTraceRepeatable(t *testing.T, path []testLink, packets []uint64, want []Snapshot) {
	t.Helper()
	assertTraceTransitions(t, want)
	for i := 0; i < 3; i++ {
		assertTraceTransitions(t, runTrace(path, packets))
	}
}

func snapshotJSON(snapshots []Snapshot) []any {
	ret := make([]any, 0, len(snapshots))
	for _, s := range snapshots {
		ret = append(ret, s.Json())
	}
	return ret
}

func TestProbeBWGainCycle(t *testing.T) {
	clock := newSimClock(sampleStartTime)
	lim := &Limiter{
		clock:        clock,
		rtPropFilter: rtPropFilter{Estimate: time.Millisecond},
		randInt:      func() int { return 6 },
	}
	state := &probeBWState{}
	state.initialize(lim)
	if lim.pacingGain != 1.25 {
		t.Fatalf("initial ProbeBW pacing gain = %v; want 1.25", lim.pacingGain)
	}
	for _, want := range []float64{0.75, 1, 1} {
		now := state.nextPacingGainChange.Add(time.Nanosecond)
		clock.AdvanceTo(now)
		state.postAck(lim, packetMeta{}, now)
		if lim.pacingGain != want {
			t.Fatalf("ProbeBW pacing gain = %v; want %v", lim.pacingGain, want)
		}
	}
}

func runTrace(path []testLink, packetSizes []uint64) []Snapshot {
	sim := newSimulator(panicTestingT{}, path...)
	return sim.run(packetSizes)
}

// panicTestingT is used by runTrace, whose caller owns test reporting.
type panicTestingT struct{}

func (panicTestingT) Helper()           {}
func (panicTestingT) Fatal(args ...any) { panic(fmt.Sprint(args...)) }
