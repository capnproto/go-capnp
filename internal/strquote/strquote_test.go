package strquote

import (
	"bytes"
	"testing"
)

func TestAppend(t *testing.T) {
	tests := []struct {
		name   string
		in     string
		expect string
	}{
		{"printable", "Hello! No escaping.", `"Hello! No escaping."`},
		{"controls", "\a\b\f\n\r\t\v", `"\a\b\f\n\r\t\v"`},
		{"quotes", "\"'", `"\"\'"`},
		{"backslash", "\\", `"\\"`},
		{"binary", "\x00\x1f\x7f\xff", `"\x00\x1f\x7f\xff"`},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf []byte
			out := Append(buf, []byte(tc.in))
			if !bytes.Equal(out, []byte(tc.expect)) {
				t.Errorf("Append(%q) = %q; want %q", tc.in, out, tc.expect)
			}
		})
	}
}
