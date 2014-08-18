package capn_test

import (
	"bytes"
	"fmt"
	"testing"

	capn "github.com/glycerine/go-capnproto"
	air "github.com/glycerine/go-capnproto/aircraftlib"
	cv "github.com/smartystreets/goconvey/convey"
)

func ValAtBit(value int64, bitPosition uint) bool {
	return (int64(1)<<bitPosition)&value != 0
}

func TestValAtBit(t *testing.T) {
	cv.Convey("ValAtBit should return the value of bit i", t, func() {
		two_to_62 := int64(2) << 61
		//fmt.Printf("pow(2,62) = %x\n", two_to_62)

		cv.So(ValAtBit(0, 0), cv.ShouldEqual, false)

		cv.So(ValAtBit(1, 0), cv.ShouldEqual, true)

		cv.So(ValAtBit(2, 1), cv.ShouldEqual, true)
		cv.So(ValAtBit(2, 0), cv.ShouldEqual, false)

		cv.So(ValAtBit(3, 2), cv.ShouldEqual, false)
		cv.So(ValAtBit(3, 1), cv.ShouldEqual, true)
		cv.So(ValAtBit(3, 0), cv.ShouldEqual, true)

		cv.So(ValAtBit(4, 3), cv.ShouldEqual, false)
		cv.So(ValAtBit(4, 2), cv.ShouldEqual, true)
		cv.So(ValAtBit(4, 1), cv.ShouldEqual, false)
		cv.So(ValAtBit(4, 0), cv.ShouldEqual, false)

		cv.So(ValAtBit(5, 3), cv.ShouldEqual, false)
		cv.So(ValAtBit(5, 2), cv.ShouldEqual, true)
		cv.So(ValAtBit(5, 1), cv.ShouldEqual, false)
		cv.So(ValAtBit(5, 0), cv.ShouldEqual, true)

		cv.So(ValAtBit(6, 3), cv.ShouldEqual, false)
		cv.So(ValAtBit(6, 2), cv.ShouldEqual, true)
		cv.So(ValAtBit(6, 1), cv.ShouldEqual, true)
		cv.So(ValAtBit(6, 0), cv.ShouldEqual, false)

		cv.So(ValAtBit(7, 3), cv.ShouldEqual, false)
		cv.So(ValAtBit(7, 2), cv.ShouldEqual, true)
		cv.So(ValAtBit(7, 1), cv.ShouldEqual, true)
		cv.So(ValAtBit(7, 0), cv.ShouldEqual, true)

		cv.So(ValAtBit(8, 3), cv.ShouldEqual, true)
		cv.So(ValAtBit(8, 2), cv.ShouldEqual, false)
		cv.So(ValAtBit(8, 1), cv.ShouldEqual, false)
		cv.So(ValAtBit(8, 0), cv.ShouldEqual, false)

		cv.So(ValAtBit(two_to_62, 62), cv.ShouldEqual, true)
		cv.So(ValAtBit(two_to_62, 2), cv.ShouldEqual, false)
		cv.So(ValAtBit(two_to_62, 1), cv.ShouldEqual, false)
		cv.So(ValAtBit(two_to_62, 0), cv.ShouldEqual, false)

		cv.So(ValAtBit(9, 3), cv.ShouldEqual, true)
		cv.So(ValAtBit(9, 2), cv.ShouldEqual, false)
		cv.So(ValAtBit(9, 1), cv.ShouldEqual, false)
		cv.So(ValAtBit(9, 0), cv.ShouldEqual, true)
	})
}

func zboolvec_value_FilledSegment(value int64, elementCount uint) (*capn.Segment, []byte) {
	seg := capn.NewBuffer(nil)
	z := air.NewRootZ(seg)
	list := seg.NewBitList(int(elementCount))
	if value > 0 {
		for i := uint(0); i < elementCount; i++ {
			list.Set(int(i), ValAtBit(value, i))
		}
	}
	z.SetBoolvec(list)

	buf := bytes.Buffer{}
	seg.WriteTo(&buf)
	return seg, buf.Bytes()
}

func TestBitList(t *testing.T) {
	seg, _ := zboolvec_value_FilledSegment(5, 3)
	text := CapnpDecodeSegment(seg, "", "aircraftlib/aircraft.capnp", "Z")

	expectedText := `(boolvec = [true, false, true])`

	cv.Convey("Given a go-capnproto created List(Bool) Z::boolvec with bool values [true, false, true]", t, func() {
		cv.Convey("When we decode it with capnp", func() {
			cv.Convey(fmt.Sprintf("Then we should get the expected text '%s'", expectedText), func() {
				cv.So(text, cv.ShouldEqual, expectedText)
			})
			cv.Convey("And our data should contain Z_BOOLVEC with contents true, false, true", func() {
				z := air.ReadRootZ(seg)
				cv.So(z.Which(), cv.ShouldEqual, air.Z_BOOLVEC)

				var bitlist = z.Boolvec()
				cv.So(bitlist.Len(), cv.ShouldEqual, 3)
				cv.So(bitlist.At(0), cv.ShouldEqual, true)
				cv.So(bitlist.At(1), cv.ShouldEqual, false)
				cv.So(bitlist.At(2), cv.ShouldEqual, true)
			})
		})
	})

}

func TestWriteBitList0(t *testing.T) {
	seg, _ := zboolvec_value_FilledSegment(0, 1)
	cv.Convey("Given a go-capnproto created List(Bool) Z::boolvec with bool values [false]", t, func() {
		cv.Convey("Decoding it with c++ capnp should yield the expected text", func() {
			cv.So(CapnpDecodeSegment(seg, "", "aircraftlib/aircraft.capnp", "Z"), cv.ShouldEqual, `(boolvec = [false])`)
		})
	})

	cv.Convey("And we should be able to read back what we wrote", t, func() {
		z := air.ReadRootZ(seg)
		cv.So(z.Which(), cv.ShouldEqual, air.Z_BOOLVEC)

		var bitlist = z.Boolvec()
		cv.So(bitlist.Len(), cv.ShouldEqual, 1)
		cv.So(bitlist.At(0), cv.ShouldEqual, false)
	})
}

func TestWriteBitList1(t *testing.T) {
	seg, _ := zboolvec_value_FilledSegment(1, 1)
	cv.Convey("Given a go-capnproto created List(Bool) Z::boolvec with bool values [true]", t, func() {
		cv.Convey("Decoding it with c++ capnp should yield the expected text", func() {
			cv.So(CapnpDecodeSegment(seg, "", "aircraftlib/aircraft.capnp", "Z"), cv.ShouldEqual, `(boolvec = [true])`)
		})
	})

	cv.Convey("And we should be able to read back what we wrote", t, func() {
		z := air.ReadRootZ(seg)
		cv.So(z.Which(), cv.ShouldEqual, air.Z_BOOLVEC)

		var bitlist = z.Boolvec()
		cv.So(bitlist.Len(), cv.ShouldEqual, 1)
		cv.So(bitlist.At(0), cv.ShouldEqual, true)
	})

}

func TestWriteBitList2(t *testing.T) {
	seg, _ := zboolvec_value_FilledSegment(2, 2)
	//seg, by := zboolvec_value_FilledSegment(2, 2)
	//ShowBytes(by, 0)
	cv.Convey("Given a go-capnproto created List(Bool) Z::boolvec with bool values [false, true]", t, func() {
		cv.Convey("Decoding it with c++ capnp should yield the expected text", func() {
			cv.So(CapnpDecodeSegment(seg, "", "aircraftlib/aircraft.capnp", "Z"), cv.ShouldEqual, `(boolvec = [false, true])`)
		})
	})

	cv.Convey("And we should be able to read back what we wrote", t, func() {
		z := air.ReadRootZ(seg)
		cv.So(z.Which(), cv.ShouldEqual, air.Z_BOOLVEC)

		var bitlist = z.Boolvec()
		cv.So(bitlist.Len(), cv.ShouldEqual, 2)
		cv.So(bitlist.At(0), cv.ShouldEqual, false)
		cv.So(bitlist.At(1), cv.ShouldEqual, true)
	})
}

func TestWriteBitList3(t *testing.T) {
	seg, _ := zboolvec_value_FilledSegment(3, 2)
	cv.Convey("Given a go-capnproto created List(Bool) Z::boolvec with bool values [true, true]", t, func() {
		cv.Convey("Decoding it with c++ capnp should yield the expected text", func() {
			cv.So(CapnpDecodeSegment(seg, "", "aircraftlib/aircraft.capnp", "Z"), cv.ShouldEqual, `(boolvec = [true, true])`)
		})
	})

	cv.Convey("And we should be able to read back what we wrote", t, func() {
		z := air.ReadRootZ(seg)
		cv.So(z.Which(), cv.ShouldEqual, air.Z_BOOLVEC)

		var bitlist = z.Boolvec()
		cv.So(bitlist.Len(), cv.ShouldEqual, 2)
		cv.So(bitlist.At(0), cv.ShouldEqual, true)
		cv.So(bitlist.At(1), cv.ShouldEqual, true)
	})

}

func TestWriteBitList4(t *testing.T) {
	seg, _ := zboolvec_value_FilledSegment(4, 3)
	cv.Convey("Given a go-capnproto created List(Bool) Z::boolvec with bool values [false, false, true]", t, func() {
		cv.Convey("Decoding it with c++ capnp should yield the expected text", func() {
			cv.So(CapnpDecodeSegment(seg, "", "aircraftlib/aircraft.capnp", "Z"), cv.ShouldEqual, `(boolvec = [false, false, true])`)
		})
	})

	cv.Convey("And we should be able to read back what we wrote", t, func() {
		z := air.ReadRootZ(seg)
		cv.So(z.Which(), cv.ShouldEqual, air.Z_BOOLVEC)

		var bitlist = z.Boolvec()
		cv.So(bitlist.Len(), cv.ShouldEqual, 3)
		cv.So(bitlist.At(0), cv.ShouldEqual, false)
		cv.So(bitlist.At(1), cv.ShouldEqual, false)
		cv.So(bitlist.At(2), cv.ShouldEqual, true)
	})
}

func TestWriteBitList21(t *testing.T) {
	seg, _ := zboolvec_value_FilledSegment(21, 5)
	cv.Convey("Given a go-capnproto created List(Bool) Z::boolvec with bool values [true, false, true, false, true]", t, func() {
		cv.Convey("Decoding it with c++ capnp should yield the expected text", func() {
			cv.So(CapnpDecodeSegment(seg, "", "aircraftlib/aircraft.capnp", "Z"), cv.ShouldEqual, `(boolvec = [true, false, true, false, true])`)
		})
	})

	cv.Convey("And we should be able to read back what we wrote", t, func() {
		z := air.ReadRootZ(seg)
		cv.So(z.Which(), cv.ShouldEqual, air.Z_BOOLVEC)

		var bitlist = z.Boolvec()
		cv.So(bitlist.Len(), cv.ShouldEqual, 5)
		cv.So(bitlist.At(0), cv.ShouldEqual, true)
		cv.So(bitlist.At(1), cv.ShouldEqual, false)
		cv.So(bitlist.At(2), cv.ShouldEqual, true)
		cv.So(bitlist.At(3), cv.ShouldEqual, false)
		cv.So(bitlist.At(4), cv.ShouldEqual, true)
	})
}

func TestWriteBitListTwo64BitWords(t *testing.T) {

	seg := capn.NewBuffer(nil)
	z := air.NewRootZ(seg)
	list := seg.NewBitList(66)
	list.Set(64, true)
	list.Set(65, true)

	z.SetBoolvec(list)

	buf := bytes.Buffer{}
	seg.WriteTo(&buf)

	cv.Convey("Given a go-capnproto created List(Bool) Z::boolvec with bool values [true (+ 64 more times)]", t, func() {
		cv.Convey("Decoding it with c++ capnp should yield the expected text", func() {
			cv.So(CapnpDecodeSegment(seg, "", "aircraftlib/aircraft.capnp", "Z"), cv.ShouldEqual, `(boolvec = [false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true, true])`)
		})
	})

	cv.Convey("And we should be able to read back what we wrote", t, func() {
		z := air.ReadRootZ(seg)
		cv.So(z.Which(), cv.ShouldEqual, air.Z_BOOLVEC)

		var bitlist = z.Boolvec()
		cv.So(bitlist.Len(), cv.ShouldEqual, 66)

		for i := 0; i < 64; i++ {
			cv.So(bitlist.At(i), cv.ShouldEqual, false)
		}
		cv.So(bitlist.At(64), cv.ShouldEqual, true)
		cv.So(bitlist.At(65), cv.ShouldEqual, true)
	})
}

/*
expected binaries: see data/bvec*.bin

jaten@i7:~/cap:master$ echo "(boolvec = [true, false, true])" | capnp encode aircraftlib/aircraft.capnp Z > bvec5.bin
jaten@i7:~/cap:master$ echo "(boolvec = [true, true])" | capnp encode aircraftlib/aircraft.capnp Z > bvec3.bin
jaten@i7:~/cap:master$ echo "(boolvec = [true])" | capnp encode aircraftlib/aircraft.capnp Z > bvec1.bin
jaten@i7:~/cap:master$ echo "(boolvec = [false])" | capnp encode aircraftlib/aircraft.capnp Z > bvec0.bin

jaten@i7:~/cap:master$ xxd -g 1 -c 8 bvec0.bin
0000000: 00 00 00 00 05 00 00 00  ........
0000008: 00 00 00 00 02 00 01 00  ........
0000010: 27 00 00 00 00 00 00 00  '.......
0000018: 00 00 00 00 00 00 00 00  ........
0000020: 01 00 00 00 09 00 00 00  ........
0000028: 00 00 00 00 00 00 00 00  ........

jaten@i7:~/cap:master$ xxd -g 1 -c 8 bvec1.bin
0000000: 00 00 00 00 05 00 00 00  ........
0000008: 00 00 00 00 02 00 01 00  ........
0000010: 27 00 00 00 00 00 00 00  '.......
0000018: 00 00 00 00 00 00 00 00  ........
0000020: 01 00 00 00 09 00 00 00  ........
0000028: 01 00 00 00 00 00 00 00  ........

jaten@i7:~/cap:master$ xxd -g 1 -c 8 bvec3.bin
0000000: 00 00 00 00 05 00 00 00  ........
0000008: 00 00 00 00 02 00 01 00  ........
0000010: 27 00 00 00 00 00 00 00  '.......
0000018: 00 00 00 00 00 00 00 00  ........
0000020: 01 00 00 00 11 00 00 00  ........
0000028: 03 00 00 00 00 00 00 00  ........

jaten@i7:~/cap:master$ xxd -g 1 -c 8 bvec5.bin
0000000: 00 00 00 00 05 00 00 00  ........
0000008: 00 00 00 00 02 00 01 00  ........
0000010: 27 00 00 00 00 00 00 00  '.......
0000018: 00 00 00 00 00 00 00 00  ........
0000020: 01 00 00 00 19 00 00 00  ........
0000028: 05 00 00 00 00 00 00 00  ........

*/
