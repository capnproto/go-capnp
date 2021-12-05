package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorString(t *testing.T) {
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
		got := New(test.typ, test.prefix, test.msg).Error()
		if got != test.want {
			t.Errorf("New(%#v, %q, %q).Error() = %q; want %q", test.typ, test.prefix, test.msg, got, test.want)
		}
	}
}

func TestTypeOf(t *testing.T) {
	tests := []struct {
		err  error
		want Type
	}{
		{nil, Failed},
		{errors.New("generic"), Failed},
		{New(Failed, "capnp", "failed error"), Failed},
		{New(Overloaded, "capnp", "overloaded error"), Overloaded},
		{New(Disconnected, "capnp", "disconnected error"), Disconnected},
		{New(Unimplemented, "capnp", "unimplemented error"), Unimplemented},
	}
	for _, test := range tests {
		if got := TypeOf(test.err); got != test.want {
			t.Errorf("TypeOf(%#v) = %#v; want %#v", test.err, got, test.want)
		}
	}
}

func TestAnnotate(t *testing.T) {
	tests := []struct {
		prefix string
		msg    string
		err    error

		want     string
		wantType Type
	}{
		{
			prefix: "prefix",
			msg:    "context",
			err:    nil,
		},
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
		got := Annotate(test.prefix, test.msg, test.err)
		if test.err == nil {
			assert.Nil(t, got)
			continue
		}

		if got.Error() != test.want {
			t.Errorf("Annotate(%q, %q, %#v).Error() = %q; %q", test.prefix, test.msg, test.err, got.Error(), test.want)
		}
		gotType := TypeOf(got)
		if gotType != test.wantType {
			t.Errorf("TypeOf(Annotate(%q, %q, %#v)) = %#v; %#v", test.prefix, test.msg, test.err, gotType, test.wantType)
		}
	}
}
