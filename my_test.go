package capnp_test

import (
	"testing"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/internal/aircraftlib"
)

func BenchmarkSetText05(b *testing.B) {
	var msg capnp.Message
	arena := &capnp.SimpleSingleSegmentArena{}

	msg.ResetForRead(arena)
	seg, err := msg.Segment(0)
	if err != nil {
		b.Fatal(err)
	}
	tx, err := aircraftlib.NewBenchmarkA(seg)
	if err != nil {
		b.Fatal(err)
	}

	err = tx.SetName("my own descr")
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		msg.ResetForRead(arena)
		seg, _ := msg.Segment(0)

		tx, err := aircraftlib.NewBenchmarkA(seg)
		if err != nil {
			b.Fatal(err)
		}

		err = tx.SetName("my own descr")
		if err != nil {
			b.Fatal(err)
		}
	}

	// b.Log(arena.String())

}

func BenchmarkSetInt(b *testing.B) {
	var msg capnp.Message
	arena := &capnp.SimpleSingleSegmentArena{}

	msg.ResetForRead(arena)
	seg, err := msg.Segment(0)
	if err != nil {
		b.Fatal(err)
	}
	tx, err := aircraftlib.NewBenchmarkA(seg)
	if err != nil {
		b.Fatal(err)
	}

	tx.SetBirthDay(0x20010101)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tx.SetBirthDay(0x20010101 + int64(i))
	}

	// b.Log(arena.String())

}
func BenchmarkSetTextBaselineCopyClear(b *testing.B) {
	var pt *[]byte
	buf := make([]byte, 1024)
	pt = &buf

	b.ResetTimer()
	src := "my own descr"
	off := 48
	for i := 0; i < b.N; i++ {
		s := (*pt)[off : off+len(src)+1]
		n := copy(s, []byte(src))
		for i := n; i < len(s); i++ {
			s[i] = 0
		}
	}
}
