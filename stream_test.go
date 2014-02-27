package capn_test

import (
	"bytes"
	"flag"
	capn "github.com/jmckaskill/go-capnproto"
	"io"
	"io/ioutil"
	"testing"
)

var benchForever bool

func init() {
	flag.BoolVar(&benchForever, "bench.forever", false, "benchmark forever")
	flag.Parse()
}

func TestReadFromStream(t *testing.T) {
	const n = 10
	r := zdateReader(n, false)
	for i := 0; i < n; i++ {
		s, err := capn.ReadFromStream(r, nil)
		if err != nil {
			t.Fatalf("ReadFromStream: %v", err)
		}
		m := ReadRootZdate(s)
		js, err := m.MarshalJSON()
		if err != nil {
			t.Fatalf("MarshalJSON: %v", err)
		}
		t.Logf("%s", string(js))
	}
}

func TestDecompressorZdate(t *testing.T) {
	const n = 10

	r := zdateReader(n, false)
	expected, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}

	r = zdateReader(n, true)
	actual, err := ioutil.ReadAll(capn.NewDecompressor(r))
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}

	if !bytes.Equal(expected, actual) {
		t.Fatal("decompressor failed")
	}
}

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
}

func TestCompressor(t *testing.T) {
	for i, test := range compressionTests {
		var buf bytes.Buffer
		c := capn.NewCompressor(&buf)
		c.Write(test.original)
		if !bytes.Equal(test.compressed, buf.Bytes()) {
			t.Errorf("test:%d: failed", i)
		}
	}
}

func TestDecompressor(t *testing.T) {
	for i, test := range compressionTests {
		for readSize := 1; readSize <= 8+2*len(test.original); readSize++ {
			r := bytes.NewReader(test.compressed)
			d := capn.NewDecompressor(r)
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

func TestReadFromPackedStream(t *testing.T) {
	const n = 10

	r := zdateReader(n, true)
	for i := 0; i < n; i++ {
		s, err := capn.ReadFromPackedStream(r, nil)
		if err != nil {
			t.Fatalf("ReadFromPackedStream: %v", err)
		}
		m := ReadRootZdate(s)
		js, err := m.MarshalJSON()
		if err != nil {
			t.Fatalf("MarshalJSON: %v", err)
		}
		t.Logf("%s", string(js))
	}
}

func BenchmarkCompressor(b *testing.B) {
	r := zdateReader(100, false)
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		panic(err)
	}
	c := capn.NewCompressor(ioutil.Discard)
	b.SetBytes(int64(len(buf)))

	b.ResetTimer()
	for i := 0; i < b.N || benchForever; i++ {
		c.Write(buf)
	}
}

func BenchmarkDecompressor(b *testing.B) {
	r := zdateReader(100, true)
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		b.Fatalf("%v", err)
	}
	b.SetBytes(int64(len(buf)))

	// determine buffer size to read all the data
	// in a single call
	r.Seek(0, 0)
	d := capn.NewDecompressor(r)
	buf, err = ioutil.ReadAll(d)
	if err != nil {
		b.Fatalf("%v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N || benchForever; i++ {
		r.Seek(0, 0)
		n, err := d.Read(buf)
		if err != nil {
			b.Fatalf("%v", err)
		}
		if n != len(buf) {
			b.Fatal("short read")
		}
	}
}
