package capnp_test

import (
	"testing"

	"zombiezen.com/go/capnproto"
	air "zombiezen.com/go/capnproto/internal/aircraftlib"
)

func TestDecode(t *testing.T) {
	const n = 10
	r := zdateReader(n, false)
	msg, err := capnp.NewDecoder(r).Decode()
	if err != nil {
		t.Fatalf("Decode: %v", err)
	}
	z, err := air.ReadRootZ(msg)
	if err != nil {
		t.Fatalf("ReadRootZ: %v", err)
	}
	if z.Which() != air.Z_Which_zdatevec {
		panic("expected Z_ZDATEVEC in root Z of segment")
	}
}

func TestDecodeBackToBack(t *testing.T) {
	const n = 10

	r := zdateReaderNBackToBack(n, false)
	d := capnp.NewDecoder(r)

	for i := 0; i < n; i++ {
		_, err := d.Decode()
		if err != nil {
			t.Fatalf("Decode: %v", err)
		}
	}
}

func TestPackedDecode(t *testing.T) {
	const n = 10

	r := zdateReaderNBackToBack(n, true)
	d := capnp.NewPackedDecoder(r)

	for i := 0; i < n; i++ {
		_, err := d.Decode()
		if err != nil {
			t.Fatalf("Decode: %v, i=%d", err, i)
		}
	}
}
