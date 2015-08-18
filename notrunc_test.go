package capnp_test

import (
	"fmt"
	"testing"

	cv "github.com/smartystreets/goconvey/convey"
	"zombiezen.com/go/capnproto"
	air "zombiezen.com/go/capnproto/internal/aircraftlib"
)

func TestDataVersioningAvoidsUnnecessaryTruncation(t *testing.T) {

	expFull := CapnpEncode("(val = 9, duo = 8, ptr1 = (val = 77), ptr2 = (val = 55))", "VerTwoDataTwoPtr")
	//expEmpty := CapnpEncode("()", "VerEmpty")

	cv.Convey("Given a struct with 0 ptr fields, and a newer version of the struct with two data and two pointer fields", t, func() {
		cv.Convey("then old code expecting the smaller struct but reading the newer-bigger struct should not truncate it if it doesn't have to (e.g. not assigning into a composite list), and should preserve all data when re-serializing it.", func() {

			msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)
			_, scratch, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)

			big, err := air.NewRootVerTwoDataTwoPtr(seg)
			cv.So(err, cv.ShouldEqual, nil)
			one, err := air.NewVerOneData(scratch)
			cv.So(err, cv.ShouldEqual, nil)
			one.SetVal(77)
			two, err := air.NewVerOneData(scratch)
			cv.So(err, cv.ShouldEqual, nil)
			two.SetVal(55)
			big.SetVal(9)
			big.SetDuo(8)
			big.SetPtr1(one)
			big.SetPtr2(two)

			bigVerBytes := ShowSeg("\n\n   with our 2x2 new big struct, segment seg is:", seg)

			cv.So(bigVerBytes, cv.ShouldResemble, expFull)

			// now pretend to be an old client, reading and writing
			// expecting an empty struct, but full data should be preserved
			// and written, because we aren't writing into a cramped/
			// fixed-space composite-list space.

			// Before test, verify that if we force reading into text-form, we get
			// what we expect.
			actEmptyCap := string(CapnpDecode(bigVerBytes, "VerEmpty"))
			cv.So(actEmptyCap, cv.ShouldResemble, "()\n")

			// okay, now the actual test:
			weThinkEmptyButActuallyFull, err := air.ReadRootVerEmpty(msg)
			cv.So(err, cv.ShouldEqual, nil)

			_, freshSeg, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)
			wrapEmpty, err := air.NewRootWrapEmpty(freshSeg)
			cv.So(err, cv.ShouldEqual, nil)

			// here is the critical step, this should not truncate:
			wrapEmpty.SetMightNotBeReallyEmpty(weThinkEmptyButActuallyFull)

			// now verify:
			freshBytes := ShowSeg("\n\n          after wrapEmpty.SetMightNotBeReallyEmpty(weThinkEmptyButActuallyFull), segment freshSeg is:", freshSeg)

			remsg, err := capnp.Unmarshal(freshBytes)
			cv.So(err, cv.ShouldEqual, nil)
			reseg, err := remsg.Segment(0)
			cv.So(err, cv.ShouldEqual, nil)
			ShowSeg("      after re-reading freshBytes, segment reseg is:", reseg)
			fmt.Printf("freshBytes decoded by capnp as Wrap2x2: '%s'\n", string(CapnpDecode(freshBytes, "Wrap2x2")))

			wrap22, err := air.ReadRootWrap2x2plus(remsg)
			cv.So(err, cv.ShouldEqual, nil)
			notEmpty, err := wrap22.MightNotBeReallyEmpty()
			cv.So(err, cv.ShouldEqual, nil)
			val := notEmpty.Val()
			cv.So(val, cv.ShouldEqual, 9)
			duo := notEmpty.Duo()
			cv.So(duo, cv.ShouldEqual, 8)
			ptr1, err := notEmpty.Ptr1()
			cv.So(err, cv.ShouldEqual, nil)
			ptr2, err := notEmpty.Ptr2()
			cv.So(err, cv.ShouldEqual, nil)
			cv.So(ptr1.Val(), cv.ShouldEqual, 77)
			cv.So(ptr2.Val(), cv.ShouldEqual, 55)
			// Tre should get the default, as it was never set
			cv.So(notEmpty.Tre(), cv.ShouldEqual, 0)
			// same for Lst3
			lst3, err := notEmpty.Lst3()
			cv.So(err, cv.ShouldEqual, nil)
			cv.So(lst3.Len(), cv.ShouldEqual, 0)
		})
	})
}
