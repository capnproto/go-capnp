package bbr

import (
	"context"
	"testing"
	"time"

	"capnproto.org/go/capnp/v3/internal/clock"
)

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
