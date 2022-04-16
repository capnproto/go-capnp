package exc_test

import (
	"errors"
	"fmt"
	"testing"

	"capnproto.org/go/capnp/v3/exc"
	"github.com/stretchr/testify/assert"
)

func TestTypeOf(t *testing.T) {
	t.Parallel()

	tests := []struct {
		err  error
		want exc.Type
	}{
		{nil, exc.Failed},
		{errors.New("generic"), exc.Failed},
		{exc.New(exc.Failed, "capnp", "failed error"), exc.Failed},
		{exc.New(exc.Overloaded, "capnp", "overloaded error"), exc.Overloaded},
		{exc.New(exc.Disconnected, "capnp", "disconnected error"), exc.Disconnected},
		{exc.New(exc.Unimplemented, "capnp", "unimplemented error"), exc.Unimplemented},
	}
	for _, test := range tests {
		assert.Equal(t, test.want, exc.TypeOf(test.err))
	}
}

func TestIsType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		err error
		tpe exc.Type
	}{
		{
			err: exc.New(exc.Failed, "capnp", "failed error"),
			tpe: exc.Failed,
		},
		{
			err: exc.New(exc.Overloaded, "capnp", "overloaded error"),
			tpe: exc.Overloaded,
		},
		{
			err: exc.New(exc.Disconnected, "capnp", "disconnected error"),
			tpe: exc.Disconnected,
		},
		{
			err: exc.New(exc.Unimplemented, "capnp", "unimplemented error"),
			tpe: exc.Unimplemented,
		},
	}
	for i, test := range tests {
		err := fmt.Errorf("test: %w", test.err)
		assert.True(t, exc.IsType(err, test.tpe),
			"case %d should match exception type '%s'", i, test.tpe)
	}
}
