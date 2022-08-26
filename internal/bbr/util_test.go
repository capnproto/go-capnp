package bbr

import (
	"context"
	"testing"
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
			delay = time.Duration(float64(p.Size) / float64(l.bandwidth))
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

// Fail the test immediately if done is receivable or closed, or becomes so within ~1ms.
// XXX: This is inherently racy. But for our purposes, we just want to give very cheap
// operations time to complete.
//
// If need be, we can probably afford to bump the threshold up to 10 or 100ms, as long
// as all tests that use it only do so a couple times and can be run in parallel.
func assertNotDone(t *testing.T, done <-chan struct{}) {
	select {
	case <-time.NewTimer(time.Millisecond).C:
	case <-done:
		t.Fatal("Packet should not have arrived yet.")
	}
}

func TestLinkDelay(t *testing.T) {
	// Try sending a packet through a path with one link of a given known delay;
	// make sure it arrives at the correct time.

	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	clock := clock.NewManual(time.Unix(1e9, 0))
	tx, rx := startTestPath(ctx, clock, testLink{delay: 4 * time.Second})
	go processAcks(ctx, rx)

	done := make(chan struct{})

	tx.Send(testPacket{
		gotResponse: func() {
			close(done)
		},
	})

	assertNotDone(t, done)
	clock.Advance(2 * time.Second)
	assertNotDone(t, done)

	clock.Advance(3 * time.Second) // This puts us at t = 5, after our delay.
	<-done                         // would deadlock if the packet still did't send.
}

const (
	bytesPerSecond bytesPerNs = 1e-9
)

func TestLinkBandwidth(t *testing.T) {
	// Try sending a packet through a path with a known bandwidth, and make sure
	// it takes the correct amount of time.

	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	clock := clock.NewManual(time.Unix(1e9, 0))

	tx, rx := startTestPath(ctx, clock, testLink{bandwidth: 10 * bytesPerSecond})
	go processAcks(ctx, rx)

	done := make(chan struct{})

	tx.Send(testPacket{
		packetMeta: packetMeta{
			Size: 25,
		},
		gotResponse: func() {
			close(done)
		},
	})

	assertNotDone(t, done)
	clock.Advance(2 * time.Second)
	assertNotDone(t, done)
	clock.Advance(1 * time.Second)
	<-done
}

func TestLinkBandwidthMultiPacket(t *testing.T) {
	// Try sending multiple packets through a path with a known bandwidth,
	// and make sure it takes the right amount of time.
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	clock := clock.NewManual(time.Unix(1e9, 0))

	tx, rx := startTestPath(ctx, clock, testLink{bandwidth: 10 * bytesPerSecond})
	go processAcks(ctx, rx)

	done1 := make(chan struct{})
	done2 := make(chan struct{})

	tx.Send(testPacket{
		packetMeta: packetMeta{
			Size: 25,
		},
		gotResponse: func() {
			close(done1)
		},
	})
	tx.Send(testPacket{
		packetMeta: packetMeta{
			Size: 30,
		},
		gotResponse: func() {
			close(done2)
		},
	})

	assertNotDone(t, done1)
	assertNotDone(t, done2)
	clock.Advance(2 * time.Second)

	assertNotDone(t, done1)
	assertNotDone(t, done2)

	clock.Advance(1 * time.Second)
	<-done1
	assertNotDone(t, done2)

	clock.Advance(3 * time.Second)
	<-done2
}
