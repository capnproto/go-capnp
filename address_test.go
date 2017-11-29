package capnp

import (
	"testing"
)

func TestAddressAddSize(t *testing.T) {
	tests := []struct {
		a   address
		sz  Size
		out address
		ok  bool
	}{
		{0, 0, 0, true},
		{0, 1, 1, true},
		{1, 0, 1, true},
		{1, 1, 2, true},
		{0xfffffff7, 0, 0xfffffff7, true},
		{0, 0xfffffff7, 0xfffffff7, true},
		{0xfffffff7, 1, 0xfffffff8, true},
		{1, 0xfffffff7, 0xfffffff8, true},
		{0xfffffff8, 0, 0xfffffff8, true},
		{0, 0xfffffff8, 0xfffffff8, true},
		{0xffffffff, 0, 0, false},
		{0, 0xffffffff, 0, false},
		{0xfffffff8, 0xfffffff8, 0, false},
		{0xffffffff, 0xfffffff8, 0, false},
		{0xfffffff8, 0xffffffff, 0, false},
		{0xffffffff, 0xffffffff, 0, false},
	}
	for _, test := range tests {
		out, ok := test.a.addSize(test.sz)
		if ok != test.ok || (ok && out != test.out) {
			t.Errorf("%#v.addSize(%#v) = %v, %t; want %v, %t", test.a, test.sz, out, ok, test.out, test.ok)
		}
	}
}

func TestAddressElement(t *testing.T) {
	tests := []struct {
		a   address
		i   int32
		sz  Size
		out address
		ok  bool
	}{
		{0, 0, 0, 0, true},
		{0, 1, 0, 0, true},
		{0, 1, 1, 1, true},
		{0, -1, 1, 0, false},
		{1, -1, 1, 0, true},
		{8, -1, 8, 0, true},
		{24, -1, 8, 16, true},
		{0, 1, 8, 8, true},
		{0, 2, 8, 16, true},
		{24, 1, 0, 24, true},
		{24, 1, 8, 32, true},
		{24, 2, 8, 40, true},
		{0, 0x1fffffff, 8, 0xfffffff8, true},
		{1, 0x1fffffff, 8, 0, false},
		{0, 0x7fffffff, 3, 0, false},
		{0xffffffff, 0x7fffffff, 0xffffffff, 0, false},
		{0xffffffff, -0x80000000, 0xffffffff, 0, false},
		{0, 0x7fffffff, 0xffffffff, 0, false},
		{0, -0x80000000, 0xffffffff, 0, false},
	}
	for _, test := range tests {
		out, ok := test.a.element(test.i, test.sz)
		if ok != test.ok || (ok && out != test.out) {
			t.Errorf("%#v.element(%d, %d) = %#v, %t; want %#v, %t", test.a, test.i, test.sz, out, ok, test.out, test.ok)
		}
	}
}

func TestSizeTimes(t *testing.T) {
	tests := []struct {
		sz  Size
		n   int32
		out Size
		ok  bool
	}{
		{0, 0, 0, true},
		{0, 0x7fffffff, 0, true},
		{0, -0x80000000, 0, true},
		{8, 0x1fffffff, 0xfffffff8, true},
		{8, 0x20000000, 0, false},
		{0x7ffffffc, 2, 0xfffffff8, true},
		{0xfffffff8, 0, 0, true},
		{0xfffffff9, 0, 0, true},
		{0xffffffff, 0, 0, true},
		{0xfffffff8, 1, 0xfffffff8, true},
		{0xfffffff9, 1, 0, false},
		{0xffffffff, 1, 0, false},
	}
	for _, test := range tests {
		out, ok := test.sz.times(test.n)
		if ok != test.ok || (ok && out != test.out) {
			t.Errorf("%#v.times(%d) = %#v, %t; want %#v, %t", test.sz, test.n, out, ok, test.out, test.ok)
		}
	}
}
