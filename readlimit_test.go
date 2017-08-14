package capnp

import "testing"

func TestMessage_canRead(t *testing.T) {
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
			name: "read a word with default limit",
			calls: []canReadCall{
				{8, true},
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
		m := &Message{TraverseLimit: test.init}
		for i, c := range test.calls {
			ok := m.canRead(c.sz)
			if ok != c.ok {
				// TODO(light): show previous calls
				t.Errorf("in %s, calls[%d] ok = %t; want %t", test.name, i, ok, c.ok)
			}
		}
	}
}

func TestMessage_ResetReadLimit(t *testing.T) {
	{
		m := &Message{TraverseLimit: 42}
		t.Log("   m := &Message{TraverseLimit: 42}")
		ok := m.canRead(42)
		t.Logf("   m.canRead(42) -> %t", ok)
		m.ResetReadLimit(8)
		t.Log("   m.ResetReadLimit(8)")
		if m.canRead(8) {
			t.Log("   m.canRead(8) -> true")
		} else {
			t.Error("!! m.canRead(8) -> false; want true")
		}
	}
	t.Log()
	{
		m := &Message{TraverseLimit: 42}
		t.Log("   m := &Message{TraverseLimit: 42}")
		ok := m.canRead(40)
		t.Logf("   m.canRead(40) -> %t", ok)
		m.ResetReadLimit(8)
		t.Log("   m.ResetReadLimit(8)")
		if m.canRead(9) {
			t.Error("!! m.canRead(9) -> true; want false")
		} else {
			t.Log("   m.canRead(9) -> false")
		}
	}
	t.Log()
	{
		m := &Message{TraverseLimit: 42}
		t.Log("   m := &Message{TraverseLimit: 42}")
		ok := m.canRead(42)
		t.Logf("   m.canRead(42) -> %t", ok)
		m.ResetReadLimit(8)
		t.Log("   m.ResetReadLimit(8)")
		if m.canRead(8) {
			t.Log("   m.canRead(8) -> true")
		} else {
			t.Error("!! m.canRead(8) -> false; want true")
		}
	}
	t.Log()
	{
		m := &Message{TraverseLimit: 42}
		t.Log("   m := &Message{TraverseLimit: 42}")
		ok := m.canRead(40)
		t.Logf("   m.canRead(40) -> %t", ok)
		m.ResetReadLimit(8)
		t.Log("   m.ResetReadLimit(8)")
		if m.canRead(9) {
			t.Error("!! m.canRead(9) -> true; want false")
		} else {
			t.Log("   m.canRead(9) -> false")
		}
	}
	t.Log()
	{
		m := new(Message)
		t.Log("   m := new(Message)")
		m.ResetReadLimit(0)
		t.Log("   m.ResetReadLimit(0)")
		if !m.canRead(0) {
			t.Error("!! m.canRead(0) -> false; want true")
		} else {
			t.Log("   m.canRead(0) -> true")
		}
	}
	t.Log()
	{
		m := new(Message)
		t.Log("   m := new(Message)")
		m.ResetReadLimit(0)
		t.Log("   m.ResetReadLimit(0)")
		if m.canRead(1) {
			t.Error("!! m.canRead(1) -> true; want false")
		} else {
			t.Log("   m.canRead(1) -> false")
		}
	}
}

func TestMessage_Unread(t *testing.T) {
	{
		m := &Message{TraverseLimit: 42}
		t.Log("   m := &Message{TraverseLimit: 42}")
		ok := m.canRead(42)
		t.Logf("   m.canRead(42) -> %t", ok)
		m.Unread(8)
		t.Log("   m.Unread(8)")
		if m.canRead(8) {
			t.Log("   m.canRead(8) -> true")
		} else {
			t.Error("!! m.canRead(8) -> false; want true")
		}
	}
	t.Log()
	{
		m := &Message{TraverseLimit: 42}
		t.Log("   m := &Message{TraverseLimit: 42}")
		ok := m.canRead(40)
		t.Logf("   m.canRead(40) -> %t", ok)
		m.Unread(8)
		t.Log("   m.Unread(8)")
		if m.canRead(9) {
			t.Log("   m.canRead(9) -> true")
		} else {
			t.Error("!! m.canRead(9) -> false; want true")
		}
	}
}
