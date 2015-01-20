package capn_test

import (
	"testing"

	"github.com/glycerine/go-capnproto"
	air "github.com/glycerine/go-capnproto/aircraftlib"
)

func TestPipelineOpString(t *testing.T) {
	tests := []struct {
		op capn.PipelineOp
		s  string
	}{
		{
			capn.PipelineOp{Field: 4},
			"get field 4",
		},
		{
			capn.PipelineOp{Field: 4, DefaultSegment: capn.NewBuffer(nil), DefaultOffset: 0},
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
	s := capn.NewBuffer(nil)
	root := air.NewRootStackingRoot(s)
	a := air.NewStackingA(s)
	a.SetNum(1)
	root.SetA(a) // assumed to be pointer index 1
	b := air.NewStackingB(s)
	b.SetNum(2)
	a.SetB(b)

	d := capn.NewBuffer(nil)
	da := air.NewRootStackingA(d)
	da.SetNum(56)
	db := air.NewStackingB(d)
	db.SetNum(78)
	da.SetB(db)

	tests := []struct {
		p         capn.Object
		transform []capn.PipelineOp
		out       capn.Object
	}{
		{
			capn.Object(root),
			nil,
			capn.Object(root),
		},
		{
			capn.Object(root),
			[]capn.PipelineOp{},
			capn.Object(root),
		},
		{
			capn.Object(root),
			[]capn.PipelineOp{
				{Field: 0},
			},
			capn.Object{},
		},
		{
			capn.Object(root),
			[]capn.PipelineOp{
				{Field: 0, DefaultSegment: d, DefaultOffset: 0},
			},
			capn.Object(da),
		},
		{
			capn.Object(root),
			[]capn.PipelineOp{
				{Field: 1},
			},
			capn.Object(a),
		},
		{
			capn.Object(root),
			[]capn.PipelineOp{
				{Field: 1, DefaultSegment: d, DefaultOffset: 0},
			},
			capn.Object(a),
		},
		{
			capn.Object(root),
			[]capn.PipelineOp{
				{Field: 1},
				{Field: 0},
			},
			capn.Object(b),
		},
		{
			capn.Object(root),
			[]capn.PipelineOp{
				{Field: 0},
				{Field: 0},
			},
			capn.Object{},
		},
		{
			capn.Object(root),
			[]capn.PipelineOp{
				{Field: 0, DefaultSegment: d, DefaultOffset: 0},
				{Field: 0},
			},
			capn.Object(db),
		},
		{
			capn.Object(root),
			[]capn.PipelineOp{
				{Field: 0},
				{Field: 0, DefaultSegment: d, DefaultOffset: 0},
			},
			capn.Object(da),
		},
		{
			capn.Object(root),
			[]capn.PipelineOp{
				{Field: 0, DefaultSegment: d, DefaultOffset: 0},
				{Field: 1, DefaultSegment: d, DefaultOffset: 0},
			},
			capn.Object(da),
		},
	}

	for _, test := range tests {
		out := capn.TransformObject(test.p, test.transform)
		if out != test.out {
			t.Errorf("TransformObject(%+v, %v) = %+v; want %+v", test.p, test.transform, out, test.out)
		}
	}
}
