package util

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChkfatal(t *testing.T) {
	assert.NotPanics(t, func() {
		Chkfatal(nil)
	}, "Chkfatal does not panic on a nil error.")

	assert.Panics(t, func() {
		Chkfatal(errors.New("Some error."))
	}, "Chkfatal panics on a non-nil error.")
}
