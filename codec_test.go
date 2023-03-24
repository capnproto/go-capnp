package capnp

import (
	"bytes"
	"io"
	"testing"
)

func TestEncoder(t *testing.T) {
	t.Parallel()

	for i, test := range serializeTests {
		if test.decodeFails {
			continue
		}
		msg := &Message{Arena: test.arena()}
		var buf bytes.Buffer
		enc := NewEncoder(&buf)
		err := enc.Encode(msg)
		out := buf.Bytes()
		if err != nil {
			if !test.encodeFails {
				t.Errorf("serializeTests[%d] - %s: Encode error: %v", i, test.name, err)
			}
			continue
		}
		if test.encodeFails {
			t.Errorf("serializeTests[%d] - %s: Encode success; want error", i, test.name)
			continue
		}
		if !bytes.Equal(out, test.out) {
			t.Errorf("serializeTests[%d] - %s: Encode = % 02x; want % 02x", i, test.name, out, test.out)
		}
	}
}

func TestDecoder(t *testing.T) {
	t.Parallel()

	for i, test := range serializeTests {
		if test.encodeFails {
			continue
		}
		msg, err := NewDecoder(bytes.NewReader(test.out)).Decode()
		if err != nil {
			if !test.decodeFails {
				t.Errorf("serializeTests[%d] - %s: Decode error: %v", i, test.name, err)
			}
			if test.decodeError != nil && err != test.decodeError {
				t.Errorf("serializeTests[%d] - %s: Decode error: %v; want %v", i, test.name, err, test.decodeError)
			}
			continue
		}
		if test.decodeFails {
			t.Errorf("serializeTests[%d] - %s: Decode success; want error", i, test.name)
			continue
		}
		if msg.NumSegments() != int64(len(test.segs)) {
			t.Errorf("serializeTests[%d] - %s: Decode NumSegments() = %d; want %d", i, test.name, msg.NumSegments(), len(test.segs))
			continue
		}
		for j := range test.segs {
			seg, err := msg.Segment(SegmentID(j))
			if err != nil {
				t.Errorf("serializeTests[%d] - %s: Decode Segment(%d) error: %v", i, test.name, j, err)
				continue
			}
			if !bytes.Equal(seg.Data(), test.segs[j]) {
				t.Errorf("serializeTests[%d] - %s: Decode Segment(%d) = % 02x; want % 02x", i, test.name, j, seg.Data(), test.segs[j])
			}
		}
	}
}

func TestDecoder_MaxMessageSize(t *testing.T) {
	t.Parallel()

	zeroWord := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	tests := []struct {
		name    string
		maxSize uint64
		r       io.Reader
		ok      bool
	}{
		{
			name:    "header too big",
			maxSize: 15,
			r: bytes.NewReader([]byte{
				0x02, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
			}),
		},
		{
			name:    "header at limit",
			maxSize: 16,
			r: bytes.NewReader([]byte{
				0x02, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00,
			}),
			ok: true,
		},
		{
			name:    "body too large",
			maxSize: 64,
			r: io.MultiReader(
				bytes.NewReader([]byte{
					0x00, 0x00, 0x00, 0x00,
					0x09, 0x00, 0x00, 0x00,
				}),
				bytes.NewReader(bytes.Repeat(zeroWord, 9)),
			),
		},
		{
			name:    "body plus header too large",
			maxSize: 64,
			r: io.MultiReader(
				bytes.NewReader([]byte{
					0x00, 0x00, 0x00, 0x00,
					0x08, 0x00, 0x00, 0x00,
				}),
				bytes.NewReader(bytes.Repeat(zeroWord, 8)),
			),
		},
		{
			name:    "body plus header at limit",
			maxSize: 72,
			r: io.MultiReader(
				bytes.NewReader([]byte{
					0x00, 0x00, 0x00, 0x00,
					0x08, 0x00, 0x00, 0x00,
				}),
				bytes.NewReader(bytes.Repeat(zeroWord, 8)),
			),
			ok: true,
		},
	}
	for _, test := range tests {
		d := NewDecoder(test.r)
		d.MaxMessageSize = test.maxSize
		_, err := d.Decode()
		switch {
		case err != nil && test.ok:
			t.Errorf("%s test: Decode error: %v", test.name, err)
		case err == nil && !test.ok:
			t.Errorf("%s test: Decode success; want error", test.name)
		}
	}
}

// TestStreamHeaderPadding is a regression test for
// stream header padding.
//
// Encoder reuses a buffer for stream header marshalling,
// this test ensures that the padding is explicitly
// zeroed. This was not done in previous versions and
// resulted in the padding being garbage.
func TestStreamHeaderPadding(t *testing.T) {
	t.Parallel()

	msg := &Message{
		Arena: MultiSegment([][]byte{
			incrementingData(8),
			incrementingData(8),
			incrementingData(8),
		}),
	}
	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	err := enc.Encode(msg)
	buf.Reset()
	if err != nil {
		t.Fatalf("Encode error: %v", err)
	}
	msg = &Message{
		Arena: MultiSegment([][]byte{
			incrementingData(8),
			incrementingData(8),
		}),
	}
	err = enc.Encode(msg)
	out := buf.Bytes()
	if err != nil {
		t.Fatalf("Encode error: %v", err)
	}
	want := []byte{
		0x01, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x00, 0x00,
		0x01, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x01, 0x02, 0x03,
		0x04, 0x05, 0x06, 0x07,
		0x00, 0x01, 0x02, 0x03,
		0x04, 0x05, 0x06, 0x07,
	}
	if !bytes.Equal(out, want) {
		t.Errorf("Encode = % 02x; want % 02x", out, want)
	}
}
