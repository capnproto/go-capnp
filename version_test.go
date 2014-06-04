package capn_test

import (
	"bytes"
	"fmt"
	"testing"

	capn "github.com/glycerine/go-capnproto"
	air "github.com/glycerine/go-capnproto/aircraftlib"
	cv "github.com/smartystreets/goconvey/convey"
)

func TestV1DataVersioningZeroDataToOne(t *testing.T) {

	exp := CapnpEncode("(mylist = [(),()])", "HoldsVerEmptyList")

	cv.Convey("Given a struct with 0 data/0 ptr fields, and a newer version of the struct with 1 data fields", t, func() {
		cv.Convey("then setting the old into a new list should work, as should setting the new into the old", func() {

			seg := capn.NewBuffer(nil)
			scratch := capn.NewBuffer(nil)
			holder := air.NewRootHoldsVerEmptyList(seg)

			emptyBytes :=
				ShowSeg("          after NewRootHoldsVerEmptyList(seg), segment seg is:", seg)

			elist := air.NewVerEmptyList(seg, 2)
			plist := capn.PointerList(elist)

			addList2bytes :=
				ShowSeg("          pre NewVerEmpty(scratch), segment seg is:", seg)
			cv.So(emptyBytes, cv.ShouldResemble, addList2bytes)

			e0 := air.NewVerEmpty(scratch)
			e1 := air.NewVerEmpty(scratch)
			plist.Set(0, capn.Object(e0))
			plist.Set(1, capn.Object(e1))

			ShowSeg("          pre SetWaitingjobs, segment seg is:", seg)

			fmt.Printf("Then we do the SetMylist():\n")
			holder.SetMylist(elist)

			// save
			buf := bytes.Buffer{}
			seg.WriteTo(&buf)
			act := buf.Bytes()
			save(act, "my.act.holder.elist")

			// show
			ShowSeg("          actual:\n", seg)

			fmt.Printf("act decoded by capnp: '%s'\n", string(CapnpDecode(act, "HoldsVerEmptyList")))

			fmt.Printf("expected:\n")
			ShowBytes(exp, 10)
			fmt.Printf("exp decoded by capnp: '%s'\n", string(CapnpDecode(exp, "HoldsVerEmptyList")))

			cv.So(act, cv.ShouldResemble, exp)
		})
	})
}

/*
func TestDataVersioningZeroDataToTwo(t *testing.T) {

	cv.Convey("Given a struct with 0 data/0 ptr fields, and a newer version of the struct with 2 data/0 pointer fields", t, func() {
		cv.Convey("then setting the old into a new list should work, as should setting the new into the old", func() {
		})
	})
}

func TestDataVersioningZeroPointersToOne(t *testing.T) {

	cv.Convey("Given a struct with 0 data/0 ptr fields, and a newer version of the struct with 0 data/1 pointer fields", t, func() {
		cv.Convey("then setting the old into a new list should work, as should setting the new into the old", func() {
		})
	})
}

func TestDataVersioningZeroPointersToTwo(t *testing.T) {
	cv.Convey("Given a struct with 0 data/0 ptr fields, and a newer version of the struct with 0 data/2 pointer fields", t, func() {
		cv.Convey("then setting the old into a new list should work, as should setting the new into the old", func() {
		})
	})
}

func TestDataVersioningToTwoDataTwoPtr(t *testing.T) {
	cv.Convey("Given a struct with 0 data/0 ptr fields, and a newer version of the struct with 2 data/2 pointer fields", t, func() {
		cv.Convey("then setting the old into a new list should work, as should setting the new into the old", func() {
		})
	})
}
*/
