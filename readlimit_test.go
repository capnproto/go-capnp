package capnp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	t.Parallel()
	t.Helper()

	for _, tt := range []struct {
		limit            int
		initial, canRead Size
		readLim          uint64
	}{
		{
			limit:   42,
			initial: 42,
			readLim: 8,
			canRead: 8,
		},
		{
			limit:   42,
			initial: 40,
			readLim: 8,
			canRead: 9,
		},
		{},
		{canRead: 1},
	} {
		t.Run(fmt.Sprintf("TraverseLimit=%d", tt.limit), func(t *testing.T) {
			m := &Message{TraverseLimit: uint64(tt.limit)}
			require.True(t, m.canRead(tt.initial), "should be able to read %s bytes", tt.initial)

			m.ResetReadLimit(tt.readLim)
			if tt.readLim < uint64(tt.canRead) {
				assert.False(t, m.canRead(tt.canRead), "should fail to read %d bytes", tt.canRead)
			} else {
				assert.True(t, m.canRead(tt.canRead), "should succeed in reading %d bytes", tt.canRead)
			}
		})
	}
}

func TestMessage_Unread(t *testing.T) {
	t.Parallel()
	t.Helper()

	t.Run("UnreadFromTraverseLimit", func(t *testing.T) {
		t.Parallel()

		m := &Message{TraverseLimit: 42}
		require.True(t, m.canRead(42), "should be able to read up to TraverseLimit")

		m.Unread(8)
		assert.True(t, m.canRead(8), "should be able to read 8 bytes after unreading")
	})

	t.Run("UnreadBeforeTraverseLimit", func(t *testing.T) {
		t.Parallel()

		m := &Message{TraverseLimit: 42}
		require.True(t, m.canRead(40), "should be able to read fewer than 'TraverseLimit' bytes")

		m.Unread(8)
		assert.True(t, m.canRead(9), "should be able to read 9 bytes after unreading")
	})
}
