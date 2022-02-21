package exc

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnwrap(t *testing.T) {
	t.Parallel()

	var (
		errGeneric = errors.New("something went wrong")
		err        = Annotate("annotated", "test", errGeneric)
		exc        Exception
	)

	assert.EqualError(t, errors.Unwrap(err), "test: something went wrong")
	assert.ErrorIs(t, err, errGeneric)

	assert.ErrorAs(t, err, &exc)
	assert.Equal(t, "annotated", exc.Prefix)
	assert.EqualError(t, exc.Cause, "test: something went wrong")
}

func TestErrorString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		typ    Type
		prefix string
		msg    string
		want   string
	}{
		{Failed, "", "", ""},
		{Failed, "", "goofed", "goofed"},
		{Failed, "capnp", "goofed", "capnp: goofed"},
	}
	for _, test := range tests {
		err := New(test.typ, test.prefix, test.msg)
		assert.EqualError(t, err, test.want)
	}
}

func TestAnnotate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		prefix string
		msg    string
		err    error

		want     string
		wantType Type
	}{
		{
			msg:      "context",
			err:      errors.New("goofed"),
			want:     "context: goofed",
			wantType: Failed,
		},
		{
			msg:      "context",
			err:      New(Failed, "", "goofed"),
			want:     "context: goofed",
			wantType: Failed,
		},
		{
			msg:      "context",
			err:      New(Failed, "capnp", "goofed"),
			want:     "context: capnp: goofed",
			wantType: Failed,
		},
		{
			msg:      "context",
			err:      New(Unimplemented, "", "unimplemented"),
			want:     "context: unimplemented",
			wantType: Unimplemented,
		},
		{
			msg:      "context",
			err:      New(Unimplemented, "capnp", "unimplemented"),
			want:     "context: capnp: unimplemented",
			wantType: Unimplemented,
		},
		{
			prefix:   "capnp",
			msg:      "context",
			err:      errors.New("goofed"),
			want:     "capnp: context: goofed",
			wantType: Failed,
		},
		{
			prefix:   "capnp",
			msg:      "context",
			err:      New(Failed, "", "goofed"),
			want:     "capnp: context: goofed",
			wantType: Failed,
		},
		{
			prefix:   "capnp",
			msg:      "context",
			err:      New(Failed, "capnp", "goofed"),
			want:     "capnp: context: goofed",
			wantType: Failed,
		},
		{
			prefix:   "rpc",
			msg:      "context",
			err:      New(Failed, "capnp", "goofed"),
			want:     "rpc: context: capnp: goofed",
			wantType: Failed,
		},
	}
	for _, test := range tests {
		err := Annotate(test.prefix, test.msg, test.err)
		assert.EqualError(t, err, test.want)
		assert.Equal(t, test.wantType, TypeOf(err))
	}
}
