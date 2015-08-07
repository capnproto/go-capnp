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
			capnp.PipelineOp{Field: 4, DefaultSegment: capnp.NewBuffer(nil), DefaultOffset: 0},
			"get field 4 with default",
		},
	}
	for _, test := range tests {
		if s := test.op.String(); s != test.s {
			t.Errorf("%#v.String() = %q; want %q", test.op, s, test.s)
		}
	}
}

func TestTransformObject(t *testing.T) {
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
		p         capnp.Object
		transform []capnp.PipelineOp
		out       capnp.Object
	}{
		{
			capnp.Object(root),
			nil,
			capnp.Object(root),
		},
		{
			capnp.Object(root),
			[]capnp.PipelineOp{},
			capnp.Object(root),
		},
		{
			capnp.Object(root),
			[]capnp.PipelineOp{
				{Field: 0},
			},
			capnp.Object{},
		},
		{
			capnp.Object(root),
			[]capnp.PipelineOp{
				{Field: 0, DefaultSegment: d, DefaultOffset: 0},
			},
			capnp.Object(da),
		},
		{
			capnp.Object(root),
			[]capnp.PipelineOp{
				{Field: 1},
			},
			capnp.Object(a),
		},
		{
			capnp.Object(root),
			[]capnp.PipelineOp{
				{Field: 1, DefaultSegment: d, DefaultOffset: 0},
			},
			capnp.Object(a),
		},
		{
			capnp.Object(root),
			[]capnp.PipelineOp{
				{Field: 1},
				{Field: 0},
			},
			capnp.Object(b),
		},
		{
			capnp.Object(root),
			[]capnp.PipelineOp{
				{Field: 0},
				{Field: 0},
			},
			capnp.Object{},
		},
		{
			capnp.Object(root),
			[]capnp.PipelineOp{
				{Field: 0, DefaultSegment: d, DefaultOffset: 0},
				{Field: 0},
			},
			capnp.Object(db),
		},
		{
			capnp.Object(root),
			[]capnp.PipelineOp{
				{Field: 0},
				{Field: 0, DefaultSegment: d, DefaultOffset: 0},
			},
			capnp.Object(da),
		},
		{
			capnp.Object(root),
			[]capnp.PipelineOp{
				{Field: 0, DefaultSegment: d, DefaultOffset: 0},
				{Field: 1, DefaultSegment: d, DefaultOffset: 0},
			},
			capnp.Object(da),
		},
	}

	for _, test := range tests {
		out := capnp.TransformObject(test.p, test.transform)
		if out != test.out {
			t.Errorf("TransformObject(%+v, %v) = %+v; want %+v", test.p, test.transform, out, test.out)
		}
	}
}
