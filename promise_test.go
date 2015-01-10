package capn_test

import (
	"testing"

	"github.com/glycerine/go-capnproto"
)

func TestPromiseOpString(t *testing.T) {
	tests := []struct {
		op capn.PromiseOp
		s  string
	}{
		{
			capn.PromiseOp{Field: 4},
			"get field 4",
		},
		{
			capn.PromiseOp{Field: 4, DefaultSegment: capn.NewBuffer(nil), DefaultOffset: 0},
			"get field 4 with default",
		},
	}
	for _, test := range tests {
		if s := test.op.String(); s != test.s {
			t.Errorf("%#v.String() = %q; want %q", test.op, s, test.s)
		}
	}
}
