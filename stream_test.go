package capn_test

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	capn "github.com/glycerine/go-capnproto"
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
	z := ReadRootZ(s)
	if z.Which() != Z_ZDATEVEC {
		panic("expected Z_ZDATEVEC in root Z of segment")
	}
	zdatelist := z.Zdatevec()

	for i := 0; i < n; i++ {
		zdate := zdatelist.At(i)
		js, err := zdate.MarshalJSON()
		if err != nil {
			t.Fatalf("MarshalJSON: %v", err)
		}
		t.Logf("%s", string(js))
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

/*
func TestDecompressorZdate2(t *testing.T) {
	const n = 2

	r := zdateReader(n, false)
	expected, err := ioutil.ReadAll(r)
	fmt.Printf("expected: '%#v' of len(%d)\n", expected, len(expected))
	// prints:
	// expected: []byte{
	// 0x0, 0x0, 0x0, 0x0, 0x7, 0x0, 0x0, 0x0,
	// 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x1, 0x0,
	// 0x25, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	// 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0,
	// 0x1, 0x0, 0x0, 0x0, 0x14, 0x0, 0x0, 0x0,
	// 0xd4, 0x7, 0xc, 0x7, 0xd5, 0x7, 0xc, 0x7,
	// 0xd4, 0x7, 0xc, 0x7, 0x0, 0x0, 0x0, 0x0,
	// 0xd5, 0x7, 0xc, 0x7, 0x0, 0x0, 0x0, 0x0}  of len(64)
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

	r = bytes.NewReader(slicePacked)

	fnout := "zdate2.packed.dat"
	f, err := os.Create(fnout)
	if err != nil {
		panic(err)
	}
	nbytes, err := io.Copy(f, bytes.NewReader(slicePacked))
	if err != nil {
		panic(err)
	}
	if nbytes <= 0 {
		panic(fmt.Sprintf(fmt.Sprintf("no bytes written to '%s'", fnout)))
	}
	f.Close()
	fmt.Printf("TestDecompressorZdate2(): wroted packed bytes out to file '%s'\n", fnout)

	help := "run: cat zdate2.packed.dat | bin/decp | bin/enc > zdate2.unpacked.dat\n"
	fn := "zdate2.unpacked.dat"
	capnpUn, err := os.Open(fn)
	if err != nil {
		fmt.Printf(help)
		panic(err)
	}
	capnpUnBytes, err := ioutil.ReadAll(capnpUn)
	if err != nil {
		panic(err)
	}
	fmt.Printf("\nfrom %s: capnpUnBytes: '%#v' of len %d\n", fn, capnpUnBytes, len(capnpUnBytes))
	// capnpUnBytes is the expected values, *for real*
	// prints: from zdate2.unpacked.dat: capnpUnBytes:
	// from zdate2.unpacked.dat: capnpUnBytes: '[]byte{0x0, 0x0, 0x0, 0x0, 0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x1, 0x0, 0x25, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x14, 0x0, 0x0, 0x0, 0xd4, 0x7, 0xc, 0x7, 0xd5, 0x7, 0xc, 0x7}' of len 48

	// parsed into words:
	// 0x0,  0x0, 0x0, 0x0, 0x5,  0x0, 0x0, 0x0,
	// 0x0,  0x0, 0x0, 0x0, 0x2,  0x0, 0x1, 0x0,
	// 0x25, 0x0, 0x0, 0x0, 0x0,  0x0, 0x0, 0x0,
	// 0x0,  0x0, 0x0, 0x0, 0x0,  0x0, 0x0, 0x0,
	// 0x1,  0x0, 0x0, 0x0, 0x14, 0x0, 0x0, 0x0,
	// 0xd4, 0x7, 0xc, 0x7, 0xd5, 0x7, 0xc, 0x7}' of len 48

	actual, err := ioutil.ReadAll(capn.NewDecompressor(bytes.NewReader(slicePacked)))
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}

	cv.Convey("Given the []byte slice from a capnp conversion from packed to unpacked form of a two Zdate vector", t, func() {
		cv.Convey("When we use go-capnproto NewDecompressor", func() {
			cv.Convey("Then we should get the same unpacked bytes as capnp provides", func() {
				cv.So(len(actual), cv.ShouldResemble, len(capnpUnBytes))
				cv.So(actual, cv.ShouldResemble, capnpUnBytes)
			})
			cv.Convey("Then we expect that the go-capnproto should produce the same unpacked bytes as the capnp (c++ original)", func() {
				cv.So(len(expected), cv.ShouldResemble, len(capnpUnBytes))
				cv.So(expected, cv.ShouldResemble, capnpUnBytes)
			})
		})
	})

	if !bytes.Equal(expected, actual) {
		fmt.Printf("expected to get: '%#v' (len: %d)\n actually observed instead: '%#v' (len %d)\n", expected, len(expected), actual, len(actual))
		t.Fatal("decompressor failed")
	}
}
*/

/*
func TestDecompressorZdate2(t *testing.T) {
	const n = 1

	r := zdateReader(n, false)
	expected, err := ioutil.ReadAll(r)
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	fmt.Printf("expected is '%#v'\n", expected)
	// prints: expected is '[]byte{0x0, 0x0, 0x0, 0x0, 0x6, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x1, 0x0, 0x25, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0xc, 0x0, 0x0, 0x0, 0xd4, 0x7, 0xc, 0x7, 0x0, 0x0, 0x0, 0x0, 0xd4, 0x7, 0xc, 0x7, 0x0, 0x0, 0x0, 0x0}'

	eDecode := CapnpDecodeBuf(expected, "", "", "Z", false)
	fmt.Printf("expected decoded is: '%s'\n", eDecode)

	_, byteSlice := zdateFilledSegment(n, true)
	//r = zdateReader(n, true)

	// save byteSlice to file, seems capnp is crashing on it?
	f, err := os.Create("packed.byteslice.dat")
	if err != nil {
		panic(err)
	}
	f.Write(byteSlice)
	f.Close()

	fmt.Printf("byteSlice of packed one-Zdate ZdateVector is: '%#v' of len %d\n", byteSlice, len(byteSlice))
	// byteSlice of packed one-Zdate ZdateVector is: '[]byte{
	// 0x10, 0x6, 0x50, 0x2, 0x1, 0x1, 0x25, 0x0, 0x0, 0x11, 0x1, 0xc, 0xf, 0xd4, 0x7, 0xc, 0x7, 0xf, 0xd4, 0x7, 0xc, 0x7}' of len 22
	// note that it has the last 5 bytes duplicated again at the end, which is wrong packing compared to capnp encode --packed packing.

	// but the capnp compression of the same data is:
	// cat packed.byteslice.dat   | capnp decode --packed --short test.capnp Z | capnp encode --packed test.capnp Z  > capnp.packed.one.zdate.vector.dat
	capnpPacked, err := os.Open("capnp.packed.one.zdate.vector.dat")
	if err != nil {
		panic(err)
	}
	capnpPackedBytes, err := ioutil.ReadAll(capnpPacked)
	if err != nil {
		panic(err)
	}
	fmt.Printf("from capnpPacked, byteSlice of packed one-Zdate ZdateVector is: '%#v' of len %d\n", capnpPackedBytes, len(capnpPackedBytes))
	// prints from capnpPacked, byteSlice of packed one-Zdate ZdateVector is: '[]byte{
	// 0x10, 0x5, 0x50, 0x2, 0x1, 0x1, 0x25, 0x0, 0x0, 0x11, 0x1, 0xc, 0xf, 0xd4, 0x7, 0xc, 0x7}' of len 17

	// check the packing-- is it wrong? no looks okay on creation.
	actDecode := CapnpDecodeBuf(byteSlice, "", "", "Z", true)
	fmt.Printf("packed byteSlice decoded is: '%s'\n", actDecode)

	if actDecode != eDecode {
		msg := "actual actDecode does not match expected eDecode from packed.\n"
		fmt.Printf(msg)
		panic(msg)
	} else {
		fmt.Printf("actual actDecode matches expected eDecode from packed.\n")
	}

	fmt.Printf("\nbytesSlice is %#v\n", byteSlice)
	r = bytes.NewReader(byteSlice)
	actual, err := ioutil.ReadAll(capn.NewDecompressor(r))
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}

	fmt.Printf("actual = %#v\n", actual)

	f, err = os.Create("decoded.from.packed.dat")
	if err != nil {
		panic(err)
	}
	nbytes, err := io.Copy(f, bytes.NewReader(actual))
	if err != nil {
		panic(err)
	}
	if nbytes <= 0 {
		panic(fmt.Sprintf("no bytes written to decoded.from.packed.dat"))
	}
	f.Close()

	if !bytes.Equal(expected, actual) {
		fmt.Printf("expected to get   : '%s'\n actually observed: '%s'\n", expected, actual)
	}

	//fmt.Printf("actual decoded is: '%s'\n", CapnpDecodeBuf(actual, "", "", "Z", true)) // corruption detected here.
	// seeing src/kj/io.c++:40: requirement not met: expected n >= minBytes; Premature EOF
	// expected to get   : `\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00%\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\d324\x00\x00\x00\x00\d324\x00\x00\x00\x00`
	//  actually observed: `\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00%\x00\x00\x00\x00\x00\x00\x00(8B of 0 missing at (0-)byte 24)\x00\x00\x00\x00\x00\x00\d324\x00\x00\x00\x00\d324\x00\x00\x00\x00`

	if !bytes.Equal(expected, actual) {
		t.Fatal("decompressor failed")
	}
}
*/

func TestDecodeOnKnownWellPackedData(t *testing.T) {

	// generate expected from
	// cat capnp.packed.one.zdate.vector.dat | capnp decode --packed --short test.capnp Z | capnp encode test.capnp Z  > capnp.unpacked.one.zdate.vector.dat
	//capnpUnpacked, err := os.Open("capnp.unpacked.one.zdate.vector.dat")
	//if err != nil {
	//   panic(err)
	//}
	//capnpUnpackedBytes, err := ioutil.ReadAll(capnpUnpacked)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("capnpUnpacked generation cmd: cat capnp.packed.one.zdate.vector.dat | capnp decode --packed --short test.capnp Z | capnp encode test.capnp Z  > capnp.unpacked.one.zdate.vector.dat\n")
	//fmt.Printf("from capnpUnpacked, byteSlice of packed one-Zdate ZdateVector is: '%#v' of len %d\n", capnpUnpackedBytes, len(capnpUnpackedBytes))
	// prints from capnpUnpacked, byteSlice of packed one-Zdate ZdateVector is: '[]byte{0x0, 0x0, 0x0, 0x0, 0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x1, 0x0, 0x25, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0xc, 0x0, 0x0, 0x0, 0xd4, 0x7, 0xc, 0x7, 0x0, 0x0, 0x0, 0x0}' of len 48

	// length 17
	byteSliceIn := []byte{0x10, 0x5, 0x50, 0x2, 0x1, 0x1, 0x25, 0x0, 0x0, 0x11, 0x1, 0xc, 0xf, 0xd4, 0x7, 0xc, 0x7}
	fmt.Printf("len of byteSliceIn is %d\n", len(byteSliceIn))
	// annotated: byteSliceIn := []byte{tag:0x10, 0x5, tag:0x50, 0x2, 0x1, tag:0x1, 0x25, tag:0x0, 0x0, tag:0x11, 0x1, 0xc, tag:0xf, 0xd4, 0x7, 0xc, 0x7}

	// length 48
	expectedOut := []byte{0x0, 0x0, 0x0, 0x0, 0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x1, 0x0, 0x25, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0xc, 0x0, 0x0, 0x0, 0xd4, 0x7, 0xc, 0x7, 0x0, 0x0, 0x0, 0x0}

	//fmt.Printf("len of expectedOut is %d\n", len(expectedOut))

	//fmt.Printf("\nbytesSlice is %#v\n", byteSliceIn)
	r := bytes.NewReader(byteSliceIn)
	actual, err := ioutil.ReadAll(capn.NewDecompressor(r))
	if err != nil {
		panic(err)
	}

	// length 40, missing 8 bytes upon expansion: looks the the 0x00 0x00 isn't being expanded into a word of zeros properly
	//	fmt.Printf("len of actual is %d\n", len(actual))
	//	fmt.Printf("actual is %#v\n", actual)
	// actual is []byte{0x0, 0x0, 0x0, 0x0, 0x5, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0, 0x1, 0x0, 0x25, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0xc, 0x0, 0x0, 0x0, 0xd4, 0x7, 0xc, 0x7, 0x0, 0x0, 0x0, 0x0}

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

func TestDecompressor3(t *testing.T) {
	i := 3
	test := compressionTests[i]
	fmt.Printf("   test.original = %#v\n test.compressed = %#v\n", test.original, test.compressed)
	//	for readSize := 1; readSize <= 8+2*len(test.original); readSize++ {
	readSize := 1
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
	}

	if !bytes.Equal(test.original, actual) {
		t.Errorf("test:%d readSize:%d: bytes not equal", i, readSize)
	}
	//	}
	fmt.Printf("func TestDecompressor3(t *testing.T)")

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
