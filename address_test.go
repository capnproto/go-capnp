package capnp

import (
	"testing"
)

func TestAddressElement(t *testing.T) {
	tests := []struct {
		a   Address
		i   int32
		sz  Size
		out Address
	}{
		{0, 0, 0, 0},
		{0, 1, 0, 0},
		{0, 1, 8, 8},
		{0, 2, 8, 16},
		{24, 1, 0, 24},
		{24, 1, 8, 32},
		{24, 2, 8, 40},
	}
	for _, test := range tests {
		if out := test.a.element(test.i, test.sz); out != test.out {
			t.Errorf("%#v.element(%d, %d) = %#v; want %#v", test.a, test.i, test.sz, out, test.out)
		}
	}
}
