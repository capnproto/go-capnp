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
	assert.Equal(t, value, *first.Value(),
		"Ref.Value() should return the value passed")
	second := first.AddRef()
	assert.Equal(t, value, *second.Value(),
		"second ref should have the same value as the first")
	first.Release()
	assert.False(t, released,
		"Releasing the first ref should keep the value alive")
	assert.Equal(t, value, *second.Value(),
		"Value should be the same after releasing the first ref")
	assert.Panics(t, func() {
		first.Value()
	}, "Trying to access the value via a released ref should panic, even if the value is still live.")
	first.Release()
	assert.False(t, released, "Calling Release() twice should have no effect")
	second.Release()
	assert.True(t, released, "Releasing the second reference should drop the value")
}
