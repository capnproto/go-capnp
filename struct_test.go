package capnp_test

import (
	"testing"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/internal/aircraftlib"
)

// BenchmarkSetText benchmarks setting a single text field in a message.
func BenchmarkSetText(b *testing.B) {
	var msg capnp.Message
	var arena = &capnp.SimpleSingleSegmentArena{}
	arena.ReplaceBuffer(make([]byte, 0, 1024))

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// NOTE: Needs to be ResetForRead() because Reset() allocates
		// the root pointer. This is part of API madness.
		msg.ResetForRead(arena)

		a, err := aircraftlib.AllocateNewRootBenchmark(&msg)
		if err != nil {
			b.Fatal(err)
		}

		err = a.SetName("my name")
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSetTextFlat(b *testing.B) {
	var msg capnp.Message
	var arena = &capnp.SimpleSingleSegmentArena{}
	arena.ReplaceBuffer(make([]byte, 0, 1024))

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// NOTE: Needs to be ResetForRead() because Reset() allocates
		// the root pointer. This is part of API madness.
		msg.ResetForRead(arena)

		a, err := aircraftlib.AllocateNewRootBenchmark(&msg)
		if err != nil {
			b.Fatal(err)
		}

		err = a.FlatSetName("my name")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkSetTextUpdate benchmarks updating the text field in-place.
func BenchmarkSetTextUpdate(b *testing.B) {
	var msg capnp.Message
	var arena = &capnp.SimpleSingleSegmentArena{}
	arena.ReplaceBuffer(make([]byte, 0, 1024))

	// NOTE: Needs to be ResetForRead() because Reset() allocates
	// the root pointer. This is part of API madness.
	msg.ResetForRead(arena)

	a, err := aircraftlib.AllocateNewRootBenchmark(&msg)
	if err != nil {
		b.Fatal(err)
	}

	err = a.SetName("my name")
	if err != nil {
		b.Fatal(err)
	}

	// WHY?!?!?!?
	msg.ResetReadLimit(1 << 31)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		err := a.UpdateName("my name")
		if err != nil {
			b.Fatal(err)
		}
	}
}
