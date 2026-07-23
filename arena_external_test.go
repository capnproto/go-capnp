package capnp_test

import (
	"bytes"
	"errors"
	"testing"

	"capnproto.org/go/capnp/v3"
	air "capnproto.org/go/capnp/v3/internal/aircraftlib"
)

const fixedArenaText = "fixed arena text"

var fixedArenaTextBytes = []byte(fixedArenaText)

type externalArena struct {
	segs []*capnp.Segment
}

func (a *externalArena) NumSegments() int64 {
	return int64(len(a.segs))
}

func (a *externalArena) Segment(id capnp.SegmentID) *capnp.Segment {
	if int(id) >= len(a.segs) {
		return nil
	}
	return a.segs[id]
}

func (a *externalArena) Allocate(minsz capnp.Size, msg *capnp.Message, preferred *capnp.Segment) (*capnp.Segment, capnp.Address, error) {
	if preferred != nil && cap(preferred.Data())-len(preferred.Data()) >= int(minsz) {
		addr := capnp.Address(len(preferred.Data()))
		preferred.SetData(preferred.Data()[:len(preferred.Data())+int(minsz)])
		preferred.BindTo(msg)
		return preferred, addr, nil
	}

	capacity := int(minsz) + 8
	if capacity < 16 {
		capacity = 16
	}
	seg := capnp.NewSegment(capnp.SegmentID(len(a.segs)), make([]byte, int(minsz), capacity))
	seg.BindTo(msg)
	a.segs = append(a.segs, seg)
	return seg, 0, nil
}

func (*externalArena) Release() {}

func TestExternalArena(t *testing.T) {
	arena := new(externalArena)
	msg, seg, err := capnp.NewMessage(arena)
	if err != nil {
		t.Fatal(err)
	}

	root, err := capnp.NewRootStruct(seg, capnp.ObjectSize{PointerCount: 1})
	if err != nil {
		t.Fatal(err)
	}
	child, err := capnp.NewStruct(root.Segment(), capnp.ObjectSize{DataSize: 8})
	if err != nil {
		t.Fatal(err)
	}
	child.SetUint64(0, 42)
	if err := root.SetPtr(0, child.ToPtr()); err != nil {
		t.Fatal(err)
	}

	data, err := msg.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := capnp.Unmarshal(data)
	if err != nil {
		t.Fatal(err)
	}
	ptr, err := decoded.Root()
	if err != nil {
		t.Fatal(err)
	}
	got, err := ptr.Struct().Ptr(0)
	if err != nil {
		t.Fatal(err)
	}
	if got.Struct().Uint64(0) != 42 {
		t.Errorf("decoded child data = %d; want 42", got.Struct().Uint64(0))
	}
}

func TestFixedArenaReset(t *testing.T) {
	var backing [128]byte
	arena := newFixedArena(backing[:])
	msg, seg, err := capnp.NewMessage(&arena)
	if err != nil {
		t.Fatal(err)
	}

	for _, text := range []string{"first", "second"} {
		root, err := air.NewRootHoldsText(seg)
		if err != nil {
			t.Fatal(err)
		}
		if err := root.SetTxt(text); err != nil {
			t.Fatal(err)
		}
		got, err := root.Txt()
		if err != nil {
			t.Fatal(err)
		}
		if got != text {
			t.Errorf("root.Txt() = %q; want %q", got, text)
		}
		if got, want := &seg.Data()[0], &backing[0]; got != want {
			t.Errorf("segment data starts at %p; want %p", got, want)
		}
		if got, want := cap(seg.Data()), cap(backing); got != want {
			t.Errorf("segment data capacity = %d; want %d", got, want)
		}

		seg, err = msg.Reset(&arena)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestFixedArenaErrorsAndRecovery(t *testing.T) {
	// HoldsText needs a root word (8 bytes), three pointers (24 bytes), and
	// an 8-byte text allocation. The short text fits exactly; the long one does not.
	var backing [40]byte
	arena := newFixedArena(backing[:])
	msg, seg, err := capnp.NewMessage(&arena)
	if err != nil {
		t.Fatal(err)
	}

	root, err := air.NewRootHoldsText(seg)
	if err != nil {
		t.Fatal(err)
	}
	if err := root.SetTxt("too long"); !errors.Is(err, errFixedArenaCapacity) {
		t.Fatalf("SetTxt() error = %v; want errors.Is(err, errFixedArenaCapacity)", err)
	}

	seg, err = msg.Reset(&arena)
	if err != nil {
		t.Fatal(err)
	}
	root, err = air.NewRootHoldsText(seg)
	if err != nil {
		t.Fatal(err)
	}
	if err := root.SetTxt("ok"); err != nil {
		t.Fatal(err)
	}
	if got, err := root.Txt(); err != nil || got != "ok" {
		t.Errorf("root.Txt() = %q, %v; want %q, nil", got, err, "ok")
	}

	if _, _, err := arena.Allocate(0, msg, nil); !errors.Is(err, errFixedArenaSegment) {
		t.Fatalf("Allocate(nil) error = %v; want errors.Is(err, errFixedArenaSegment)", err)
	}
	foreign := capnp.NewSegment(1, nil)
	if _, _, err := arena.Allocate(0, msg, foreign); !errors.Is(err, errFixedArenaSegment) {
		t.Fatalf("Allocate(foreign) error = %v; want errors.Is(err, errFixedArenaSegment)", err)
	}
}

func TestFixedArenaAllocs(t *testing.T) {
	var backing [128]byte
	arena := newFixedArena(backing[:])
	msg, _, err := capnp.NewMessage(&arena)
	if err != nil {
		t.Fatal(err)
	}

	seg, err := msg.Reset(&arena)
	if err != nil {
		t.Fatal(err)
	}
	root, err := air.NewRootHoldsText(seg)
	if err != nil {
		t.Fatal(err)
	}
	if err := root.SetTxt(fixedArenaText); err != nil {
		t.Fatal(err)
	}
	got, err := root.TxtBytes()
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, fixedArenaTextBytes) {
		t.Fatalf("root.TxtBytes() = %q; want %q", got, fixedArenaTextBytes)
	}

	allocs := testing.AllocsPerRun(100, func() {
		seg, err := msg.Reset(&arena)
		if err != nil {
			t.Fatal(err)
		}
		root, err := air.NewRootHoldsText(seg)
		if err != nil {
			t.Fatal(err)
		}
		if err := root.SetTxt(fixedArenaText); err != nil {
			t.Fatal(err)
		}
		got, err := root.TxtBytes()
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(got, fixedArenaTextBytes) {
			t.Fatalf("root.TxtBytes() = %q; want %q", got, fixedArenaTextBytes)
		}
	})
	if allocs != 0 {
		t.Errorf("allocations per run = %f; want 0", allocs)
	}
}

func BenchmarkFixedArena(b *testing.B) {
	var backing [128]byte
	arena := newFixedArena(backing[:])
	msg, _, err := capnp.NewMessage(&arena)
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()
	for range b.N {
		seg, err := msg.Reset(&arena)
		if err != nil {
			b.Fatal(err)
		}
		root, err := air.NewRootHoldsText(seg)
		if err != nil {
			b.Fatal(err)
		}
		if err := root.SetTxt(fixedArenaText); err != nil {
			b.Fatal(err)
		}
		got, err := root.TxtBytes()
		if err != nil {
			b.Fatal(err)
		}
		if !bytes.Equal(got, fixedArenaTextBytes) {
			b.Fatalf("root.TxtBytes() = %q; want %q", got, fixedArenaTextBytes)
		}
	}
}
