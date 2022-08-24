package bbr

import (
	"context"
	"time"

	"capnproto.org/go/capnp/v3/internal/clock"
	"capnproto.org/go/capnp/v3/internal/mpsc"
)

type testPacket struct {
	packetMeta
	gotResponse func()
}

type testLink struct {
	bandwidth bytesPerNs
	delay     time.Duration
}

func processAcks(ctx context.Context, rx *mpsc.Rx[testPacket]) {
	for {
		p, err := rx.Recv(ctx)
		if err != nil {
			return
		}
		p.gotResponse()
	}
}

func startTestPath(ctx context.Context, clock clock.Clock, links ...testLink) (*mpsc.Tx[testPacket], *mpsc.Rx[testPacket]) {
	q := mpsc.New[testPacket]()

	initTx := &q.Tx

	for _, l := range links {
		rx := &q.Rx
		q = mpsc.New[testPacket]()
		tx := &q.Tx
		go l.run(ctx, clock, rx, tx)
	}

	return initTx, &q.Rx
}

func (l *testLink) run(ctx context.Context, clock clock.Clock, rx *mpsc.Rx[testPacket], tx *mpsc.Tx[testPacket]) {
	timer := clock.NewTimer(0)
	<-timer.Chan()
	for {

		// Wait for a packet to arrive:
		p, err := rx.Recv(ctx)
		if err != nil {
			return
		}

		var delay time.Duration
		if l.bandwidth > 0 {
			// We're the bottleneck; take an appropriate amount of time
			// to process the packet, based on its size.
			delay = time.Duration(p.Size / int64(l.bandwidth))
		} else {
			// We're not the bottleneck; just add our constant delay.
			delay = l.delay
		}

		// Wait until the right time, and then pass it along:
		timer.Reset(delay)
		select {
		case <-ctx.Done():
			return
		case <-timer.Chan():
			tx.Send(p)
		}
	}
}
