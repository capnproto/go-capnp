package bbr

import (
	"context"
	"testing"
	"time"

	"capnproto.org/go/capnp/v3/internal/clock"
	"github.com/stretchr/testify/assert"
)

func TestStartupOneAtATime10Sec(t *testing.T) {
	testStartupOneAtATime(t, 10*time.Second, 10, 1)
}

func TestStartupOneAtATime100Sec(t *testing.T) {
	t.Skip("FIXME: I'm not sure this test makes any sense.")
	testStartupOneAtATime(t, 100*time.Second, 10, 3)
}

func testStartupOneAtATime(t *testing.T, rtProp time.Duration, msgBytes uint64, rounds int) {
	t.Parallel()
	clock := clock.NewManual(sampleStartTime)
	lim := NewLimiter(clock)
	defer lim.Release()

	t.Log("Now: ", clock.Now())

	ctx := context.Background()

	lim.whilePaused(func() {
		t.Log("Next send time: ", lim.nextSendTime)
		_ = lim.state.(*startupState)
	})

	for i := 0; i < rounds; i++ {
		got, err := lim.StartMessage(ctx, msgBytes)
		assert.Nil(t, err)
		clock.Advance(rtProp)
		lim.whilePaused(func() {
			assert.Equal(t, lim.inflight(), uint64(10))
		})
		got()

		lim.whilePaused(func() {
			assert.Equal(t, lim.inflight(), uint64(0))
			s := lim.state.(*startupState)
			assert.Equal(t, lim.rtPropFilter.Estimate, rtProp)
			assert.Equal(t, s.plateauRounds, i)
			t.Log("Next send time: ", lim.nextSendTime)
		})
	}
	ctx, _ = context.WithTimeout(ctx, 10*time.Millisecond)
	_, err := lim.StartMessage(ctx, msgBytes)
	assert.NotNil(t, err, "StartMessage() should fail.")
	assert.Equal(t, ctx.Err(), err, "timeout should happen before send.")
}
