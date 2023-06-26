package rpc

import (
	"sort"
	"testing"
)

func TestIDGen(t *testing.T) {
	t.Run("NoReplacement", func(t *testing.T) {
		var gen idgen[uint32]
		for i := uint32(0); i <= 128; i++ {
			got := gen.next()
			if got != i {
				t.Errorf("after %d calls, next() = %d; want %d", i, got, i)
			}
		}
	})
	t.Run("Replacement", func(t *testing.T) {
		var gen idgen[uint32]
		for i := 0; i < 64; i++ {
			gen.next()
		}
		gen.remove(42)
		gen.remove(10)
		if got, want := gen.next(), uint32(10); got != want {
			t.Errorf("next() #1 = %d; want %d", got, want)
		}
		if got, want := gen.next(), uint32(42); got != want {
			t.Errorf("next() #2 = %d; want %d", got, want)
		}
		if got, want := gen.next(), uint32(64); got != want {
			t.Errorf("next() #3 = %d; want %d", got, want)
		}
	})
}

func TestUintSet(t *testing.T) {
	tests := []struct {
		name     string
		init     func() uintSet
		contents []uint // sorted
	}{
		{
			name: "Empty",
			init: func() uintSet {
				return uintSet(nil)
			},
			contents: nil,
		},
		{
			name: "Add0",
			init: func() uintSet {
				var s uintSet
				s.add(0)
				return s
			},
			contents: []uint{0},
		},
		{
			name: "Add1",
			init: func() uintSet {
				var s uintSet
				s.add(1)
				return s
			},
			contents: []uint{1},
		},
		{
			name: "Add2",
			init: func() uintSet {
				var s uintSet
				s.add(2)
				return s
			},
			contents: []uint{2},
		},
		{
			name: "Add63",
			init: func() uintSet {
				var s uintSet
				s.add(63)
				return s
			},
			contents: []uint{63},
		},
		{
			name: "Add64",
			init: func() uintSet {
				var s uintSet
				s.add(64)
				return s
			},
			contents: []uint{64},
		},
		{
			name: "Add65",
			init: func() uintSet {
				var s uintSet
				s.add(65)
				return s
			},
			contents: []uint{65},
		},
		{
			name: "Add127",
			init: func() uintSet {
				var s uintSet
				s.add(127)
				return s
			},
			contents: []uint{127},
		},
		{
			name: "Add128",
			init: func() uintSet {
				var s uintSet
				s.add(128)
				return s
			},
			contents: []uint{128},
		},
		{
			name: "Add129",
			init: func() uintSet {
				var s uintSet
				s.add(129)
				return s
			},
			contents: []uint{129},
		},
		{
			name: "MultiAdd",
			init: func() uintSet {
				var s uintSet
				s.add(1)
				s.add(2)
				s.add(4)
				s.add(8)
				s.add(16)
				s.add(32)
				s.add(64)
				s.add(128)
				s.add(256)
				return s
			},
			contents: []uint{1, 2, 4, 8, 16, 32, 64, 128, 256},
		},
		{
			name: "MultiAddReverse",
			init: func() uintSet {
				var s uintSet
				s.add(256)
				s.add(128)
				s.add(64)
				s.add(32)
				s.add(16)
				s.add(8)
				s.add(4)
				s.add(2)
				s.add(1)
				return s
			},
			contents: []uint{1, 2, 4, 8, 16, 32, 64, 128, 256},
		},
		{
			name: "Remove",
			init: func() uintSet {
				var s uintSet
				s.add(1)
				s.add(2)
				s.add(4)
				s.add(8)
				s.add(16)
				s.add(32)
				s.add(33)
				s.add(64)
				s.add(128)
				s.add(256)
				s.add(255)
				s.remove(33)
				s.remove(255)
				s.remove(40) // not in set
				return s
			},
			contents: []uint{1, 2, 4, 8, 16, 32, 64, 128, 256},
		},
		{
			name: "RemoveEmpty",
			init: func() uintSet {
				var s uintSet
				s.remove(33)
				return s
			},
			contents: []uint{},
		},
		{
			name: "RemoveToEmpty",
			init: func() uintSet {
				var s uintSet
				s.add(1)
				s.remove(1)
				return s
			},
			contents: []uint{},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := test.init()

			// Exhaustive testing of the set is possible, but time-
			// consuming (potentially 2^64).  Instead, test from zero up to
			// a number a little past the maximum expected number.  As a
			// last sanity check, test large uints to root out any overflow
			// bugs.
			want := append([]uint(nil), test.contents...)
			sort.Slice(want, func(i, j int) bool {
				return want[i] < want[j]
			})
			end := uint(129)
			if len(want) > 0 {
				end += want[len(want)-1]
			}
			var got []uint
			for i, j := uint(0), 0; i < end; i++ {
				inSet := s.has(i)
				if inSet {
					got = append(got, i)
				}
				if j < len(want) && i == want[j] {
					j++
					if !inSet {
						t.Errorf("has(%d) = false; want true", i)
					}
				} else if inSet {
					t.Errorf("has(%d) = true; want false", i)
				}
			}
			const uintMax = ^uint(0)
			for i := uintMax; i >= uintMax-65 && i > end; i-- {
				if s.has(i) {
					t.Errorf("has(%d) = true; want false", i)
				}
			}
			if t.Failed() {
				t.Logf("set = %v; want %v", got, want)
			}

			// Test minimum.
			if got, ok := s.min(); len(want) > 0 && (!ok || got != want[0]) {
				t.Errorf("min() = %d, %t; want %d, true", got, ok, want[0])
			} else if len(want) == 0 && ok {
				t.Errorf("min() = %d, true; want _, false", got)
			}
		})
	}
}
