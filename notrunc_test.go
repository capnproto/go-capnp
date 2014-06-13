package capn_test

import (
	"fmt"
	"testing"

	capn "github.com/glycerine/go-capnproto"
	air "github.com/glycerine/go-capnproto/aircraftlib"
	cv "github.com/smartystreets/goconvey/convey"
)

func TestDataVersioningAvoidsUnnecessaryTruncation(t *testing.T) {

	expFull := CapnpEncode("(val = 9, duo = 8, ptr1 = (val = 77), ptr2 = (val = 55))", "VerTwoDataTwoPtr")
	//expEmpty := CapnpEncode("()", "VerEmpty")

	cv.Convey("Given a struct with 0 ptr fields, and a newer version of the struct with two data and two pointer fields", t, func() {
		cv.Convey("then old code expecting the smaller struct but reading the newer-bigger struct should not truncate it if it doesn't have to (e.g. not assigning into a composite list), and should preserve all data when re-serializing it.", func() {

			seg := capn.NewBuffer(nil)
			scratch := capn.NewBuffer(nil)

			big := air.NewRootVerTwoDataTwoPtr(seg)
			one := air.NewVerOneData(scratch)
			one.SetVal(77)
			two := air.NewVerOneData(scratch)
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
			weThinkEmptyButActuallyFull := air.ReadRootVerEmpty(seg)

			freshSeg := capn.NewBuffer(nil)
			wrapEmpty := air.NewRootWrapEmpty(freshSeg)

			// here is the critical step, this should not truncate:
			wrapEmpty.SetMightNotBeReallyEmpty(weThinkEmptyButActuallyFull)

			// now verify:
			freshBytes := ShowSeg("\n\n          after wrapEmpty.SetMightNotBeReallyEmpty(weThinkEmptyButActuallyFull), segment freshSeg is:", freshSeg)

			reseg, _, err := capn.ReadFromMemoryZeroCopy(freshBytes)
			if err != nil {
				panic(err)
			}
			ShowSeg("      after re-reading freshBytes, segment reseg is:", reseg)
			fmt.Printf("freshBytes decoded by capnp as Wrap2x2: '%s'\n", string(CapnpDecode(freshBytes, "Wrap2x2")))

			wrap22 := air.ReadRootWrap2x2plus(reseg)
			notEmpty := wrap22.MightNotBeReallyEmpty()
			val := notEmpty.Val()
			cv.So(val, cv.ShouldEqual, 9)
			duo := notEmpty.Duo()
			cv.So(duo, cv.ShouldEqual, 8)
			ptr1 := notEmpty.Ptr1()
			ptr2 := notEmpty.Ptr2()
			cv.So(ptr1.Val(), cv.ShouldEqual, 77)
			cv.So(ptr2.Val(), cv.ShouldEqual, 55)
			// Tre should get the default, as it was never set
			cv.So(notEmpty.Tre(), cv.ShouldEqual, 0)
			// same for Lst3
			cv.So(notEmpty.Lst3().Len(), cv.ShouldEqual, 0)
		})
	})
}
