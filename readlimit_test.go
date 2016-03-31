package capnp

import "testing"

func TestReadLimiter_canRead(t *testing.T) {
	t.Parallel()
	type canReadCall struct {
		sz Size
		ok bool
	}
	tests := []struct {
		name  string
		init  uint64
		calls []canReadCall
	}{
		{
			name: "can always read zero",
			init: 0,
			calls: []canReadCall{
				{0, true},
			},
		},
		{
			name: "can't read a byte when limit is zero",
			init: 0,
			calls: []canReadCall{
				{1, false},
			},
		},
		{
			name: "reading a word from a high limit is okay",
			init: 128,
			calls: []canReadCall{
				{8, true},
			},
		},
		{
			name: "reading a byte after depleting the limit fails",
			init: 8,
			calls: []canReadCall{
				{8, true},
				{1, false},
			},
		},
		{
			name: "reading a byte after hitting the limit fails",
			init: 8,
			calls: []canReadCall{
				{8, true},
				{1, false},
			},
		},
		{
			name: "reading a byte after hitting the limit in multiple calls fails",
			init: 8,
			calls: []canReadCall{
				{6, true},
				{2, true},
				{1, false},
			},
		},
	}
	for _, test := range tests {
		rl := &ReadLimiter{limit: test.init}
		for i, c := range test.calls {
			ok := rl.canRead(c.sz)
			if ok != c.ok {
				// TODO(light): show previous calls
				t.Errorf("in %s, calls[%d] ok = %t; want %t", test.name, i, ok, c.ok)
			}
		}
	}
}

func TestReadLimiter_Reset(t *testing.T) {
	{
		rl := &ReadLimiter{limit: 42}
		t.Log("   rl := &ReadLimiter{limit: 42}")
		ok := rl.canRead(42)
		t.Logf("   rl.canRead(42) -> %t", ok)
		rl.Reset(8)
		t.Log("   rl.Reset(8)")
		if rl.canRead(8) {
			t.Log("   rl.canRead(8) -> true")
		} else {
			t.Error("!! rl.canRead(8) -> false; want true")
		}
	}
	t.Log()
	{
		rl := &ReadLimiter{limit: 42}
		t.Log("   rl := &ReadLimiter{limit: 42}")
		ok := rl.canRead(40)
		t.Logf("   rl.canRead(40) -> %t", ok)
		rl.Reset(8)
		t.Log("   rl.Reset(8)")
		if rl.canRead(9) {
			t.Error("!! rl.canRead(9) -> true; want false")
		} else {
			t.Log("   rl.canRead(9) -> false")
		}
	}
}

func TestReadLimiter_Unread(t *testing.T) {
	{
		rl := &ReadLimiter{limit: 42}
		t.Log("   rl := &ReadLimiter{limit: 42}")
		ok := rl.canRead(42)
		t.Logf("   rl.canRead(42) -> %t", ok)
		rl.Unread(8)
		t.Log("   rl.Unread(8)")
		if rl.canRead(8) {
			t.Log("   rl.canRead(8) -> true")
		} else {
			t.Error("!! rl.canRead(8) -> false; want true")
		}
	}
	t.Log()
	{
		rl := &ReadLimiter{limit: 42}
		t.Log("   rl := &ReadLimiter{limit: 42}")
		ok := rl.canRead(40)
		t.Logf("   rl.canRead(40) -> %t", ok)
		rl.Unread(8)
		t.Log("   rl.Unread(8)")
		if rl.canRead(9) {
			t.Log("   rl.canRead(9) -> true")
		} else {
			t.Error("!! rl.canRead(9) -> false; want true")
		}
	}
}
