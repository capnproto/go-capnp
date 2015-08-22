package packed

import (
	"bytes"
	"encoding/hex"
	"io"
	"testing"
)

var compressionTests = []struct {
	original   []byte
	compressed []byte
}{
	{
		[]byte{},
		[]byte{},
	},
	{
		[]byte{0, 0, 0, 0, 0, 0, 0, 0},
		[]byte{0, 0},
	},
	{
		[]byte{0, 0, 12, 0, 0, 34, 0, 0},
		[]byte{0x24, 12, 34},
	},
	{
		[]byte{
			0x08, 0x00, 0x00, 0x00, 0x03, 0x00, 0x02, 0x00,
			0x19, 0x00, 0x00, 0x00, 0xaa, 0x01, 0x00, 0x00,
		},
		[]byte{0x51, 0x08, 0x03, 0x02, 0x31, 0x19, 0xaa, 0x01},
	},
	{
		[]byte{0x8, 0, 0, 0, 0x3, 0, 0x2, 0, 0x19, 0, 0, 0, 0xaa, 0x1, 0, 0},
		[]byte{0x51, 0x08, 0x03, 0x02, 0x31, 0x19, 0xaa, 0x01},
	},
	{
		[]byte{
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
		},
		[]byte{0x00, 0x03},
	},
	{
		[]byte{
			0x8a, 0x8a, 0x8a, 0x8a, 0x8a, 0x8a, 0x8a, 0x8a,
			0x8a, 0x8a, 0x8a, 0x8a, 0x8a, 0x8a, 0x8a, 0x8a,
			0x8a, 0x8a, 0x8a, 0x8a, 0x8a, 0x8a, 0x8a, 0x8a,
			0x8a, 0x8a, 0x8a, 0x8a, 0x8a, 0x8a, 0x8a, 0x8a,
		},
		[]byte{
			0xff,
			0x8a, 0x8a, 0x8a, 0x8a, 0x8a, 0x8a, 0x8a, 0x8a,
			0x03,
			0x8a, 0x8a, 0x8a, 0x8a, 0x8a, 0x8a, 0x8a, 0x8a,
			0x8a, 0x8a, 0x8a, 0x8a, 0x8a, 0x8a, 0x8a, 0x8a,
			0x8a, 0x8a, 0x8a, 0x8a, 0x8a, 0x8a, 0x8a, 0x8a,
		},
	},
	{
		[]byte{1, 3, 2, 4, 5, 7, 6, 8},
		[]byte{0xff, 1, 3, 2, 4, 5, 7, 6, 8, 0},
	},
	{
		[]byte{0, 0, 0, 0, 0, 0, 0, 0, 1, 3, 2, 4, 5, 7, 6, 8},
		[]byte{0, 0, 0xff, 1, 3, 2, 4, 5, 7, 6, 8, 0},
	},
	{
		[]byte{0, 0, 12, 0, 0, 34, 0, 0, 1, 3, 2, 4, 5, 7, 6, 8},
		[]byte{0x24, 12, 34, 0xff, 1, 3, 2, 4, 5, 7, 6, 8, 0},
	},
	{
		[]byte{1, 3, 2, 4, 5, 7, 6, 8, 8, 6, 7, 4, 5, 2, 3, 1},
		[]byte{0xff, 1, 3, 2, 4, 5, 7, 6, 8, 1, 8, 6, 7, 4, 5, 2, 3, 1},
	},
	{
		[]byte{
			1, 2, 3, 4, 5, 6, 7, 8,
			1, 2, 3, 4, 5, 6, 7, 8,
			1, 2, 3, 4, 5, 6, 7, 8,
			1, 2, 3, 4, 5, 6, 7, 8,
			0, 2, 4, 0, 9, 0, 5, 1,
		},
		[]byte{
			0xff, 1, 2, 3, 4, 5, 6, 7, 8,
			3,
			1, 2, 3, 4, 5, 6, 7, 8,
			1, 2, 3, 4, 5, 6, 7, 8,
			1, 2, 3, 4, 5, 6, 7, 8,
			0xd6, 2, 4, 9, 5, 1,
		},
	},
	{
		[]byte{
			1, 2, 3, 4, 5, 6, 7, 8,
			1, 2, 3, 4, 5, 6, 7, 8,
			6, 2, 4, 3, 9, 0, 5, 1,
			1, 2, 3, 4, 5, 6, 7, 8,
			0, 2, 4, 0, 9, 0, 5, 1,
		},
		[]byte{
			0xff, 1, 2, 3, 4, 5, 6, 7, 8,
			3,
			1, 2, 3, 4, 5, 6, 7, 8,
			6, 2, 4, 3, 9, 0, 5, 1,
			1, 2, 3, 4, 5, 6, 7, 8,
			0xd6, 2, 4, 9, 5, 1,
		},
	},
	{
		[]byte{
			8, 0, 100, 6, 0, 1, 1, 2,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 0, 0, 0, 0, 0, 0,
			0, 0, 1, 0, 2, 0, 3, 1,
		},
		[]byte{
			0xed, 8, 100, 6, 1, 1, 2,
			0, 2,
			0xd4, 1, 2, 3, 1,
		},
	},
	{
		[]byte{
			0x0, 0x0, 0x0, 0x0, 0x5, 0x0, 0x0, 0x0,
			0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x1, 0x0,
			0x25, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
			0x1, 0x0, 0x0, 0x0, 0xc, 0x0, 0x0, 0x0,
			0xd4, 0x7, 0xc, 0x7, 0x0, 0x0, 0x0, 0x0,
		},
		[]byte{
			0x10, 0x5,
			0x50, 0x2, 0x1,
			0x1, 0x25,
			0x0, 0x0,
			0x11, 0x1, 0xc,
			0xf, 0xd4, 0x7, 0xc, 0x7,
		},
	},
}

func TestPack(t *testing.T) {
	for _, test := range compressionTests {
		orig := make([]byte, len(test.original))
		copy(orig, test.original)
		compressed := Pack(nil, orig)
		if !bytes.Equal(compressed, test.compressed) {
			t.Errorf("Pack(nil,\n%s\n) =\n%s\n; want\n%s", hex.Dump(test.original), hex.Dump(compressed), hex.Dump(test.compressed))
		}
	}
}

func TestReader(t *testing.T) {
	for i, test := range compressionTests {
		for readSize := 1; readSize <= 8+2*len(test.original); readSize++ {
			r := bytes.NewReader(test.compressed)
			d := NewReader(r)
			buf := make([]byte, readSize)
			var actual []byte
			for {
				n, err := d.Read(buf)
				actual = append(actual, buf[:n]...)
				if err != nil {
					if err == io.EOF {
						break
					}
					t.Fatalf("Read: %v", err)
				}
			}

			if len(test.original) != len(actual) {
				t.Errorf("test:%d readSize:%d expected %d bytes, got %d",
					i, readSize, len(test.original), len(actual))
				continue
			}

			if !bytes.Equal(test.original, actual) {
				t.Errorf("test:%d readSize:%d: bytes not equal", i, readSize)
			}
		}
	}
}

var result []byte

func BenchmarkPack(b *testing.B) {
	src := bytes.Repeat([]byte{
		8, 0, 100, 6, 0, 1, 1, 2,
		8, 0, 100, 6, 0, 1, 1, 2,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 1, 0, 2, 0, 3, 1,
		'H', 'e', 'l', 'l', 'o', ',', ' ', 'W',
		'o', 'r', 'l', 'd', '!', ' ', ' ', 'P',
		'a', 'd', ' ', 't', 'e', 'x', 't', '.',
	}, 128)
	dst := make([]byte, 0, len(src))
	b.SetBytes(int64(len(src)))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Pack(dst, src)
	}
	result = dst
}

func BenchmarkDecompressor(b *testing.B) {
	const multiplier = 128
	src := bytes.Repeat([]byte{
		0xa7, 8, 100, 6, 1, 1, 2,
		0xa7, 8, 100, 6, 1, 1, 2,
		0x00, 3,
		0x2a,
		0xff, 'H', 'e', 'l', 'l', 'o', ',', ' ', 'W',
		2,
		'o', 'r', 'l', 'd', '!', ' ', ' ', 'P',
		'a', 'd', ' ', 't', 'e', 'x', 't', '.',
	}, multiplier)
	b.SetBytes(int64(len(src)))
	r := bytes.NewReader(src)

	dst := make([]byte, 10*8*multiplier)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r.Seek(0, 0)
		d := NewReader(r)
		_, err := d.Read(dst)
		if err != nil {
			b.Fatal(err)
		}
	}
	result = dst
}
