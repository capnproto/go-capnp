package capn_test

import (
	"fmt"
	"testing"

	cv "github.com/smartystreets/goconvey/convey"
)

func TestCreationOfZDate(t *testing.T) {
	const n = 1
	packed := false
	seg, _ := zdateFilledSegment(n, packed)
	text := CapnpDecodeSegment(seg, "", "test.capnp", "Z")

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
	text := CapnpDecodeSegment(seg, "", "test.capnp", "Z")

	expectedText := `(zdatevec = [(year = 2004, month = 12, day = 7), (year = 2005, month = 12, day = 7), (year = 2006, month = 12, day = 7), (year = 2007, month = 12, day = 7), (year = 2008, month = 12, day = 7), (year = 2009, month = 12, day = 7), (year = 2010, month = 12, day = 7), (year = 2011, month = 12, day = 7), (year = 2012, month = 12, day = 7), (year = 2013, month = 12, day = 7)])`

	cv.Convey("Given a go-capnproto created segment with 10 Zdate", t, func() {
		cv.Convey("When we decode it with capnp", func() {
			cv.Convey(fmt.Sprintf("Then we should get the expected text '%s'", expectedText), func() {
				cv.So(text, cv.ShouldEqual, expectedText)
			})
		})
	})
}
