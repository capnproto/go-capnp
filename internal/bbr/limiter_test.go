package bbr

import (
	"context"
	"testing"
	"time"

	"capnproto.org/go/capnp/v3/internal/clock"

	"github.com/stretchr/testify/assert"
)

func TestReleaseLimiter(t *testing.T) {
	lim := NewLimiter(clock.System)
	lim.Release()
	select {
	case <-lim.ctx.Done():
	case <-time.NewTimer(10 * time.Millisecond).C:
		t.Fatal("lim.Release did not cancel the context.")
	}

	_, err := lim.StartMessage(context.Background(), 1)
	assert.NotNil(t, err, "Error should be non-nil if the limiter has shut down.")
	assert.Equal(t, lim.ctx.Err(), err, "Error should be that of the limiter's context.")
}

// TODO: refine this, be more specific than "Stuff" wrt. what's being tested.
func TestStuff(t *testing.T) {
	// Arbitrary starting time.
	clock := clock.NewManual(time.Unix(1e9, 0))

	lim := NewLimiter(clock)
	defer lim.Release()

	ch := make(chan func())

	goStart := func(size uint64) {
		go func() {
			gotResponse, err := lim.StartMessage(context.TODO(), size)
			if err != nil {
				panic(err)
			}
			ch <- gotResponse
		}()
	}

	goStart(1)

	got1 := <-ch
	got1()
}
