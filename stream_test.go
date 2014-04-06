package capn_test

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	capn "github.com/glycerine/go-capnproto"
	air "github.com/glycerine/go-capnproto/aircraftlib"
	cv "github.com/smartystreets/goconvey/convey"
)

var benchForever bool

func init() {
	flag.BoolVar(&benchForever, "bench.forever", false, "benchmark forever")
	flag.Parse()
}

func TestReadFromStream(t *testing.T) {
	const n = 10
	r := zdateReader(n, false)
	s, err := capn.ReadFromStream(r, nil)
	if err != nil {
		t.Fatalf("ReadFromStream: %v", err)
	}
	z := air.ReadRootZ(s)
	if z.Which() != air.Z_ZDATEVEC {
		panic("expected Z_ZDATEVEC in root Z of segment")
	}
	zdatelist := z.Zdatevec()

	if capn.JSON_enabled {
		for i := 0; i < n; i++ {
			zdate := zdatelist.At(i)
			js, err := zdate.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON: %v", err)
			}
			t.Logf("%s", string(js))
		}
	}
}

func TestReadFromStreamBackToBack(t *testing.T) {
	const n = 10

	r := zdateReaderNBackToBack(n, false)

	for i := 0; i < n; i++ {
		s, err := capn.ReadFromStream(r, nil)
		if err != nil {
			t.Fatalf("ReadFromStream: %v", err)
		}
		m := air.ReadRootZdate(s)
		if capn.JSON_enabled {
			js, err := m.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON: %v", err)
			}
			t.Logf("%s", string(js))
		}
	}

}

func TestDecompressorZdate1(t *testing.T) {
	const n = 1

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
		fmt.Printf("expected to get: '%s'\n actually observed instead: '%s'\n", expected, actual)
		t.Fatal("decompressor failed")
	}
}

func TestDecompressorUNPACKZdate2(t *testing.T) {
	const n = 2

	un := zdateReader(n, false)
	expected, err := ioutil.ReadAll(un)
	fmt.Printf("expected: '%#v' of len(%d)\n", expected, len(expected))
	// prints:
	// expected: []byte{
	// 0x0,  0x0, 0x0, 0x0, 0x7,  0x0, 0x0, 0x0,
	// 0x0,  0x0, 0x0, 0x0, 0x2,  0x0, 0x1, 0x0,
	// 0x25, 0x0, 0x0, 0x0, 0x0,  0x0, 0x0, 0x0,
	// 0x0,  0x0, 0x0, 0x0, 0x0,  0x0, 0x0, 0x0,
	// 0x1,  0x0, 0x0, 0x0, 0x14, 0x0, 0x0, 0x0,
	// 0xd4, 0x7, 0xc, 0x7, 0xd5, 0x7, 0xc, 0x7,
	// 0xd4, 0x7, 0xc, 0x7, 0x0,  0x0, 0x0, 0x0,
	// 0xd5, 0x7, 0xc, 0x7, 0x0,  0x0, 0x0, 0x0}  of len(64)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}

	//r = zdateReader(n, true)
	_, slicePacked := zdateFilledSegment(n, true)
	fmt.Printf("slicePacked = %#v\n", slicePacked)
	// prints slicePacked = []byte{
	// tag0x10, 0x7, tag:0x50, 0x2, 0x1, tag0x1, 0x25, tag0x0,
	// 0x0, tag0x11, 0x1, 0x14, tag0xff, 0xd4, 0x7, 0xc,
	// 0x7, 0xd5, 0x7, 0xc, 0x7, N:0x2, 0xd4, 0x7,
	// 0xc, 0x7, 0x0, 0x0, 0x0, 0x0, 0xd5, 0x7,
	// 0xc, 0x7, 0x0, 0x0, 0x0, 0x0}

	pa := bytes.NewReader(slicePacked)

	actual, err := ioutil.ReadAll(capn.NewDecompressor(pa))
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}

	cv.Convey("Given the []byte slice from a capnp conversion from packed to unpacked form of a two Zdate vector", t, func() {
		cv.Convey("When we use go-capnproto NewDecompressor", func() {
			cv.Convey("Then we should get the same unpacked bytes as capnp provides", func() {
				cv.So(len(actual), cv.ShouldResemble, len(expected))
				cv.So(actual, cv.ShouldResemble, expected)
			})
		})
	})
}

func TestDecodeOnKnownWellPackedData(t *testing.T) {

	// length 17
	byteSliceIn := []byte{0x10, 0x5, 0x50, 0x2, 0x1, 0x1, 0x25, 0x0, 0x0, 0x11, 0x1, 0xc, 0xf, 0xd4, 0x7, 0xc, 0x7}
	fmt.Printf("len of byteSliceIn is %d\n", len(byteSliceIn))
	// annotated: byteSliceIn := []byte{tag:0x10, 0x5, tag:0x50, 0x2, 0x1, tag:0x1, 0x25, tag:0x0, 0x0, tag:0x11, 0x1, 0xc, tag:0xf, 0xd4, 0x7, 0xc, 0x7}

	// length 48
	expectedOut := []byte{0x0, 0x0, 0x0, 0x0, 0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x1, 0x0, 0x25, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0xc, 0x0, 0x0, 0x0, 0xd4, 0x7, 0xc, 0x7, 0x0, 0x0, 0x0, 0x0}

	r := bytes.NewReader(byteSliceIn)
	actual, err := ioutil.ReadAll(capn.NewDecompressor(r))
	if err != nil {
		panic(err)
	}

	cv.Convey("Given a known-to-be-correctly packed 17-byte long sequence for a ZdateVector holding a single Zdate", t, func() {
		cv.Convey("When we use go-capnproto NewDecompressor", func() {
			cv.Convey("Then we should get the same unpacked bytes as capnp provides", func() {
				cv.So(len(actual), cv.ShouldResemble, len(expectedOut))
				cv.So(actual, cv.ShouldResemble, expectedOut)
			})
		})
	})

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
		if i == 7 {
			fmt.Printf("at test 7\n")
		}
		c := capn.NewCompressor(&buf)
		c.Write(test.original)
		if !bytes.Equal(test.compressed, buf.Bytes()) {
			t.Errorf("test:%d: failed", i)
			fmt.Printf("   test.original = %#v\n test.compressed = %#v\n    buf.Bytes() =  %#v\n", test.original, test.compressed, buf.Bytes())
		}
	}
}

func TestCompressor7(t *testing.T) {
	i := 7
	test := compressionTests[i]
	var buf bytes.Buffer
	c := capn.NewCompressor(&buf)

	fmt.Printf("compressing test.original = %#v\n", test.original)

	c.Write(test.original)
	if !bytes.Equal(test.compressed, buf.Bytes()) {
		t.Errorf("test:%d: failed", i)
		fmt.Printf("   test.original = %#v\n test.compressed = %#v\n    buf.Bytes() =  %#v\n", test.original, test.compressed, buf.Bytes())
	}
}

func TestDecompressor(t *testing.T) {
	errCount := 0
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
				errCount++
				t.Errorf("test:%d readSize:%d expected %d bytes, got %d",
					i, readSize, len(test.original), len(actual))
				continue
			}

			if !bytes.Equal(test.original, actual) {
				errCount++
				t.Errorf("test:%d readSize:%d: bytes not equal", i, readSize)
			}
		}
	}
	if errCount == 0 {
		fmt.Printf("TestDecompressor() passed. (0 errors).\n")
	}

}

func TestDecompressorVerbosely(t *testing.T) {
	cv.Convey("Testing the go-capnproto Decompressor.Read() function:", t, func() {
		for _, test := range compressionTests {

			fmt.Printf("\n\nGiven compressed text '%#v'\n", test.compressed)

			//	fmt.Printf("   test.original = %#v\n test.compressed = %#v\n", test.original, test.compressed)
			for readSize := 1; readSize <= 8+2*len(test.original); readSize++ {
				fmt.Printf("\n  When we use go-capnproto NewDecompressor, with readSize: %d\n    Then we should get the original text back.", readSize)

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

				cv.So(len(actual), cv.ShouldResemble, len(test.original))
				if len(test.original) > 0 {
					cv.So(actual, cv.ShouldResemble, test.original)
				}
			} // end readSize loop
		}
		fmt.Printf("\n\n")
	})
}

func TestReadFromPackedStream(t *testing.T) {
	const n = 10

	r := zdateReaderNBackToBack(n, true)

	for i := 0; i < n; i++ {
		s, err := capn.ReadFromPackedStream(r, nil)
		if err != nil {
			t.Fatalf("ReadFromPackedStream: %v, i=%d", err, i)
		}
		m := air.ReadRootZdate(s)
		if capn.JSON_enabled {
			js, err := m.MarshalJSON()
			if err != nil {
				t.Fatalf("MarshalJSON: %v", err)
			}
			t.Logf("%s", string(js))
		}
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
