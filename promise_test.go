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
			capnp.PipelineOp{Field: 4, DefaultSegment: capnp.NewBuffer(nil), DefaultAddress: 0},
			"get field 4 with default",
		},
	}
	for _, test := range tests {
		if s := test.op.String(); s != test.s {
			t.Errorf("%#v.String() = %q; want %q", test.op, s, test.s)
		}
	}
}

func TestTransformPointer(t *testing.T) {
	s := capnp.NewBuffer(nil)
	root := air.NewRootStackingRoot(s)
	a := air.NewStackingA(s)
	a.SetNum(1)
	root.SetA(a) // assumed to be pointer index 1
	b := air.NewStackingB(s)
	b.SetNum(2)
	a.SetB(b)

	d := capnp.NewBuffer(nil)
	da := air.NewRootStackingA(d)
	da.SetNum(56)
	db := air.NewStackingB(d)
	db.SetNum(78)
	da.SetB(db)

	tests := []struct {
		p         capnp.Pointer
		transform []capnp.PipelineOp
		out       capnp.Pointer
	}{
		{
			capnp.Pointer(root),
			nil,
			capnp.Pointer(root),
		},
		{
			capnp.Pointer(root),
			[]capnp.PipelineOp{},
			capnp.Pointer(root),
		},
		{
			capnp.Pointer(root),
			[]capnp.PipelineOp{
				{Field: 0},
			},
			capnp.Pointer{},
		},
		{
			capnp.Pointer(root),
			[]capnp.PipelineOp{
				{Field: 0, DefaultSegment: d, DefaultAddress: 0},
			},
			capnp.Pointer(da),
		},
		{
			capnp.Pointer(root),
			[]capnp.PipelineOp{
				{Field: 1},
			},
			capnp.Pointer(a),
		},
		{
			capnp.Pointer(root),
			[]capnp.PipelineOp{
				{Field: 1, DefaultSegment: d, DefaultAddress: 0},
			},
			capnp.Pointer(a),
		},
		{
			capnp.Pointer(root),
			[]capnp.PipelineOp{
				{Field: 1},
				{Field: 0},
			},
			capnp.Pointer(b),
		},
		{
			capnp.Pointer(root),
			[]capnp.PipelineOp{
				{Field: 0},
				{Field: 0},
			},
			capnp.Pointer{},
		},
		{
			capnp.Pointer(root),
			[]capnp.PipelineOp{
				{Field: 0, DefaultSegment: d, DefaultAddress: 0},
				{Field: 0},
			},
			capnp.Pointer(db),
		},
		{
			capnp.Pointer(root),
			[]capnp.PipelineOp{
				{Field: 0},
				{Field: 0, DefaultSegment: d, DefaultAddress: 0},
			},
			capnp.Pointer(da),
		},
		{
			capnp.Pointer(root),
			[]capnp.PipelineOp{
				{Field: 0, DefaultSegment: d, DefaultAddress: 0},
				{Field: 1, DefaultSegment: d, DefaultAddress: 0},
			},
			capnp.Pointer(da),
		},
	}

	for _, test := range tests {
		out := capnp.TransformPointer(test.p, test.transform)
		if out != test.out {
			t.Errorf("TransformPointer(%+v, %v) = %+v; want %+v", test.p, test.transform, out, test.out)
		}
	}
}
