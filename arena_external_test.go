package capnp_test

import (
	"testing"

	"capnproto.org/go/capnp/v3"
)

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
