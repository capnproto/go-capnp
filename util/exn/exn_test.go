package exn

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestTryReturn tests what happens when Try's callback returns without
// calling throw().
func TestTryReturn(t *testing.T) {
	v, err := Try(func(throw Thrower) int {
		return 1
	})
	assert.Nil(t, err)
	assert.Equal(t, 1, v, "Result should be what was returned")
}

// TestTryThrowNil tests what happens when Try's callback calls throw(nil).
func TestTryThrowNil(t *testing.T) {
	v, err := Try(func(throw Thrower) int {
		throw(nil)
		return 1
	})
	assert.Nil(t, err)
	assert.Equal(t, 1, v, "Result should be what was returned")
}

// TestTryThrowErr tests what happens when Try's callback calls throw
// with a non-nil error.
func TestTryThrowErr(t *testing.T) {
	someErr := errors.New("Some error")
	v, err := Try(func(throw Thrower) int {
		throw(someErr)
		return 1
	})
	assert.Equal(t, someErr, err, "Error should be what we threw")
	assert.Equal(t, 0, v, "Result should be the zero value.")
}

// TestTryOtherPanic verifies that if the callback to Try panics *directly*,
// the panic bubbles up, rather than being captured by Try.
func TestTryOtherPanic(t *testing.T) {
	assert.Panics(t, func() {
		Try(func(throw Thrower) int {
			panic("A real panic!")
		})
	})
}

// Like TestTryReturn, but for Try0.
func TestTry0Return(t *testing.T) {
	err := Try0(func(throw Thrower) {
	})
	assert.Nil(t, err)
}

// Like TestTryThrowNil, but for Try0.
func TestTry0ThrowNil(t *testing.T) {
	err := Try0(func(throw Thrower) {
		throw(nil)
	})
	assert.Nil(t, err)
}

// Like TestTryThrowErr, but for Try0.
func TestTry0ThrowErr(t *testing.T) {
	someErr := errors.New("Some error")
	err := Try0(func(throw Thrower) {
		throw(someErr)
	})
	assert.Equal(t, someErr, err, "Error should be what we threw")
}
