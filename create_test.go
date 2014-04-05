package capn_test

import (
	"fmt"
	"testing"

	air "github.com/glycerine/go-capnproto/aircraftlib"
	cv "github.com/smartystreets/goconvey/convey"
)

func TestCreationOfZDate(t *testing.T) {
	const n = 1
	packed := false
	seg, _ := zdateFilledSegment(n, packed)
	text := CapnpDecodeSegment(seg, "", "aircraftlib/aircraft.capnp", "Z")

	//expectedText := `(year = 2004, month = 12, day = 7)`
	expectedText := `(zdatevec = [(year = 2004, month = 12, day = 7)])`

	cv.Convey("Given a go-capnproto created Zdate", t, func() {
		cv.Convey("When we decode it with capnp", func() {
			cv.Convey(fmt.Sprintf("Then we should get the expected text '%s'", expectedText), func() {
				cv.So(text, cv.ShouldEqual, expectedText)
			})
		})
	})
}

func TestCreationOfManyZDate(t *testing.T) {
	const n = 10
	packed := false
	seg, _ := zdateFilledSegment(n, packed)
	text := CapnpDecodeSegment(seg, "", "aircraftlib/aircraft.capnp", "Z")

	expectedText := `(zdatevec = [(year = 2004, month = 12, day = 7), (year = 2005, month = 12, day = 7), (year = 2006, month = 12, day = 7), (year = 2007, month = 12, day = 7), (year = 2008, month = 12, day = 7), (year = 2009, month = 12, day = 7), (year = 2010, month = 12, day = 7), (year = 2011, month = 12, day = 7), (year = 2012, month = 12, day = 7), (year = 2013, month = 12, day = 7)])`

	cv.Convey("Given a go-capnproto created segment with 10 Zdate", t, func() {
		cv.Convey("When we decode it with capnp", func() {
			cv.Convey(fmt.Sprintf("Then we should get the expected text '%s'", expectedText), func() {
				cv.So(text, cv.ShouldEqual, expectedText)
			})
		})
	})
}

func TestCreationOfManyZDatePacked(t *testing.T) {
	const n = 10
	packed := true
	seg, _ := zdateFilledSegment(n, packed)
	text := CapnpDecodeSegment(seg, "", "aircraftlib/aircraft.capnp", "Z")

	expectedText := `(zdatevec = [(year = 2004, month = 12, day = 7), (year = 2005, month = 12, day = 7), (year = 2006, month = 12, day = 7), (year = 2007, month = 12, day = 7), (year = 2008, month = 12, day = 7), (year = 2009, month = 12, day = 7), (year = 2010, month = 12, day = 7), (year = 2011, month = 12, day = 7), (year = 2012, month = 12, day = 7), (year = 2013, month = 12, day = 7)])`

	cv.Convey("Given a go-capnproto created a PACKED segment with 10 Zdate", t, func() {
		cv.Convey("When we decode it with capnp", func() {
			cv.Convey(fmt.Sprintf("Then we should get the expected text '%s'", expectedText), func() {
				cv.So(text, cv.ShouldEqual, expectedText)
			})
		})
	})
}

func TestSegmentWriteToPackedOfManyZDatePacked(t *testing.T) {
	const n = 10
	packed := true
	_, byteSlice := zdateFilledSegment(n, packed)

	// check the packing-- is it wrong?
	text := CapnpDecodeBuf(byteSlice, "", "", "Z", true)

	expectedText := `(zdatevec = [(year = 2004, month = 12, day = 7), (year = 2005, month = 12, day = 7), (year = 2006, month = 12, day = 7), (year = 2007, month = 12, day = 7), (year = 2008, month = 12, day = 7), (year = 2009, month = 12, day = 7), (year = 2010, month = 12, day = 7), (year = 2011, month = 12, day = 7), (year = 2012, month = 12, day = 7), (year = 2013, month = 12, day = 7)])`

	cv.Convey("Given a go-capnproto write packed with WriteToPacked() with 10 Zdate", t, func() {
		cv.Convey("When we decode it with capnp", func() {
			cv.Convey(fmt.Sprintf("Then we should get the expected text '%s'", expectedText), func() {
				cv.So(text, cv.ShouldEqual, expectedText)
			})
		})
	})
}

/// now for Zdata (not Zdate)

func TestCreationOfZData(t *testing.T) {
	const n = 20
	seg, _ := zdataFilledSegment(n)
	text := CapnpDecodeSegment(seg, "", "aircraftlib/aircraft.capnp", "Z")

	expectedText := `(zdata = (data = "\x00\x01\x02\x03\x04\x05\x06\a\b\t\n\v\f\r\x0e\x0f\x10\x11\x12\x13"))`

	cv.Convey("Given a go-capnproto created Zdata DATA element with n=20", t, func() {
		cv.Convey("When we decode it with capnp", func() {
			cv.Convey(fmt.Sprintf("Then we should get the expected text '%s'", expectedText), func() {
				cv.So(text, cv.ShouldEqual, expectedText)
			})
			cv.Convey("And our data should contain Z_ZDATA with contents 0,1,2,...,n", func() {
				z := air.ReadRootZ(seg)
				cv.So(z.Which(), cv.ShouldEqual, air.Z_ZDATA)

				var data []byte = z.Zdata().Data()
				cv.So(len(data), cv.ShouldEqual, n)
				for i := range data {
					cv.So(data[i], cv.ShouldEqual, i)
				}

			})
		})
	})

}
