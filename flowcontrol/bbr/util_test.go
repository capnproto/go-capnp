package bbr

import (
	"testing"
	"time"
)

var sampleStartTime = time.Unix(1e9, 0)

type testPacket struct {
	size        uint64
	gotResponse func()
}

type testLink struct {
	bandwidth bytesPerNs
	delay     time.Duration
}

const bytesPerSecond bytesPerNs = 1e-9

func TestLinkDelay(t *testing.T) {
	clock := newSimClock(sampleStartTime)
	path := newSimPath(clock, testLink{delay: 4 * time.Second})
	done := false
	path.send(testPacket{gotResponse: func() { done = true }})

	path.advance(2 * time.Second)
	if done {
		t.Fatal("packet arrived before its link delay elapsed")
	}
	path.advance(3 * time.Second)
	if !done {
		t.Fatal("packet did not arrive after its link delay elapsed")
	}
}

func TestMultiLinkDelay(t *testing.T) {
	clock := newSimClock(sampleStartTime)
	path := newSimPath(clock, testLink{delay: 4 * time.Second}, testLink{delay: 2 * time.Second})
	done := false
	path.send(testPacket{gotResponse: func() { done = true }})

	path.advance(5 * time.Second)
	if done {
		t.Fatal("packet arrived before it traversed both links")
	}
	path.advance(1 * time.Second)
	if !done {
		t.Fatal("packet did not arrive after it traversed both links")
	}
}

func TestLinkBandwidth(t *testing.T) {
	clock := newSimClock(sampleStartTime)
	path := newSimPath(clock, testLink{bandwidth: 10 * bytesPerSecond})
	done := false
	path.send(testPacket{size: 25, gotResponse: func() { done = true }})

	path.advance(2 * time.Second)
	if done {
		t.Fatal("packet arrived before its serialization delay elapsed")
	}
	path.advance(time.Second)
	if !done {
		t.Fatal("packet did not arrive after its serialization delay elapsed")
	}
}

func TestLinkBandwidthMultiPacket(t *testing.T) {
	clock := newSimClock(sampleStartTime)
	path := newSimPath(clock, testLink{bandwidth: 10 * bytesPerSecond})
	done1, done2 := false, false
	path.send(testPacket{size: 25, gotResponse: func() { done1 = true }})
	path.send(testPacket{size: 30, gotResponse: func() { done2 = true }})

	path.advance(2 * time.Second)
	if done1 || done2 {
		t.Fatal("packets arrived before the first serialization delay elapsed")
	}
	path.advance(time.Second)
	if !done1 || done2 {
		t.Fatal("serial link did not deliver exactly the first packet")
	}
	path.advance(3 * time.Second)
	if !done2 {
		t.Fatal("serial link did not deliver the second packet")
	}
}
