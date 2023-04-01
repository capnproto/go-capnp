package rc

import (
	"testing"

	"github.com/tj/assert"
)

func TestRef(t *testing.T) {
	released := false
	release := func() {
		released = true
	}
	value := 4

	first := NewRef(value, release)
	assert.Equal(t, value, *first.Value())
	second := first.AddRef()
	assert.Equal(t, value, *second.Value())
	first.Release()
	assert.False(t, released)
	assert.Equal(t, value, *second.Value())
	assert.Panics(t, func() {
		first.Value()
	})
	assert.Panics(t, func() {
		first.Release()
	})
	assert.False(t, released)
	second.Release()
	assert.True(t, released)
}
