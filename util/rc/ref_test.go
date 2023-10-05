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
	assert.Panics(t, func() {
		first.AddRef()
	}, "Trying to call AddRef() on a released ref should panic")
	assert.Panics(t, func() {
		first.Weak()
	}, "Trying to call Weak() on a released ref should panic")
	first.Release()
	assert.False(t, released, "Calling Release() twice should have no effect")
	second.Release()
	assert.True(t, released, "Releasing the second reference should drop the value")
}

func TestSteal(t *testing.T) {
	released := false
	release := func() {
		released = true
	}
	value := 4

	first := NewRef(value, release)
	second := first.Steal()
	assert.False(t, released, "steal should not reduce the refcount")
	assert.Panics(t, func() {
		first.Value()
	}, "receiver is invalid after steal")
	second.Release()
	assert.True(t, released, "releasing the new reference has an effect")
}

func TestWeakRef(t *testing.T) {
	released := false
	release := func() {
		released = true
	}
	value := 4

	first := NewRef(value, release)
	weak := first.Weak()

	second, ok := weak.AddRef()
	assert.True(t, ok, "WeakRef().AddRef() should succeed if the ref is live.")
	assert.Equal(t, value, *second.Value(), "The strong reference should have the correct value.")

	second.Release()
	assert.False(t, released, "Dropping the returned ref should keep the other ref alive")

	first.Release()
	assert.True(t, released, "Dropping the first ref should release the value")

	third, ok := weak.AddRef()
	assert.False(t, ok, "Creating a strong ref after the value is released should fail")
	assert.Nil(t, third, "The returned ref should be nil if creating a strong ref fails")
}
