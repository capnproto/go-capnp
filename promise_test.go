package capnp_test

import (
	"testing"

	"zombiezen.com/go/capnproto"
	air "zombiezen.com/go/capnproto/internal/aircraftlib"
)

func TestPipelineOpString(t *testing.T) {
	tests := []struct {
		op capnp.PipelineOp
		s  string
	}{
		{
			capnp.PipelineOp{Field: 4},
			"get field 4",
		},
		{
			capnp.PipelineOp{Field: 4, DefaultValue: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
			"get field 4 with default",
		},
	}
	for _, test := range tests {
		if s := test.op.String(); s != test.s {
			t.Errorf("%#v.String() = %q; want %q", test.op, s, test.s)
		}
	}
}

func mustMarshal(t *testing.T, msg *capnp.Message) []byte {
	data, err := msg.Marshal()
	if err != nil {
		t.Fatal("Marshal:", err)
	}
	return data
}

func TestTransform(t *testing.T) {
	_, s, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatal(err)
	}
	root, err := air.NewRootStackingRoot(s)
	if err != nil {
		t.Fatal(err)
	}
	a, err := air.NewStackingA(s)
	if err != nil {
		t.Fatal(err)
	}
	a.SetNum(1)
	root.SetA(a) // assumed to be pointer index 1
	b, err := air.NewStackingB(s)
	b.SetNum(2)
	a.SetB(b)

	dmsg, d, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatal(err)
	}
	da, err := air.NewRootStackingA(d)
	if err != nil {
		t.Fatal(err)
	}
	da.SetNum(56)
	db, err := air.NewStackingB(d)
	if err != nil {
		t.Fatal(err)
	}
	db.SetNum(78)
	da.SetB(db)

	tests := []struct {
		p         capnp.Pointer
		transform []capnp.PipelineOp
		out       capnp.Pointer
	}{
		{
			root,
			nil,
			root,
		},
		{
			root,
			[]capnp.PipelineOp{},
			root,
		},
		{
			root,
			[]capnp.PipelineOp{
				{Field: 0},
			},
			nil,
		},
		{
			root,
			[]capnp.PipelineOp{
				{Field: 0, DefaultValue: mustMarshal(t, dmsg)},
			},
			da,
		},
		{
			root,
			[]capnp.PipelineOp{
				{Field: 1},
			},
			a,
		},
		{
			root,
			[]capnp.PipelineOp{
				{Field: 1, DefaultValue: mustMarshal(t, dmsg)},
			},
			a,
		},
		{
			root,
			[]capnp.PipelineOp{
				{Field: 1},
				{Field: 0},
			},
			b,
		},
		{
			root,
			[]capnp.PipelineOp{
				{Field: 0},
				{Field: 0},
			},
			nil,
		},
		{
			root,
			[]capnp.PipelineOp{
				{Field: 0, DefaultValue: mustMarshal(t, dmsg)},
				{Field: 0},
			},
			db,
		},
		{
			root,
			[]capnp.PipelineOp{
				{Field: 0},
				{Field: 0, DefaultValue: mustMarshal(t, dmsg)},
			},
			da,
		},
		{
			root,
			[]capnp.PipelineOp{
				{Field: 0, DefaultValue: mustMarshal(t, dmsg)},
				{Field: 1, DefaultValue: mustMarshal(t, dmsg)},
			},
			da,
		},
	}

	for _, test := range tests {
		out, err := capnp.Transform(test.p, test.transform)
		if out != test.out {
			t.Errorf("Transform(%+v, %v) = %+v; want %+v", test.p, test.transform, out, test.out)
		}
		if err != nil {
			t.Errorf("Transform(%+v, %v) error: %v", test.p, test.transform, err)
		}
	}
}
