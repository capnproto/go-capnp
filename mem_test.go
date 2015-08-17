// +build ignore

package capnp_test

import (
	"testing"

	"zombiezen.com/go/capnproto"
	air "zombiezen.com/go/capnproto/internal/aircraftlib"
)

func TestReadFromStream(t *testing.T) {
	const n = 10
	r := zdateReader(n, false)
	s, err := capnp.ReadFromStream(r, nil)
	if err != nil {
		t.Fatalf("ReadFromStream: %v", err)
	}
	z := air.ReadRootZ(s)
	if z.Which() != air.Z_Which_zdatevec {
		panic("expected Z_ZDATEVEC in root Z of segment")
	}
}

func TestReadFromStreamBackToBack(t *testing.T) {
	const n = 10

	r := zdateReaderNBackToBack(n, false)

	for i := 0; i < n; i++ {
		_, err := capnp.ReadFromStream(r, nil)
		if err != nil {
			t.Fatalf("ReadFromStream: %v", err)
		}
	}
}

func TestReadFromPackedStream(t *testing.T) {
	const n = 10

	r := zdateReaderNBackToBack(n, true)

	for i := 0; i < n; i++ {
		_, err := capnp.ReadFromPackedStream(r, nil)
		if err != nil {
			t.Fatalf("ReadFromPackedStream: %v, i=%d", err, i)
		}
	}
}
