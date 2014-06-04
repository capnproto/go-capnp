package capn_test

import (
	"bytes"
	"fmt"
	"testing"

	capn "github.com/glycerine/go-capnproto"
	air "github.com/glycerine/go-capnproto/aircraftlib"
	cv "github.com/smartystreets/goconvey/convey"
)

func TestV0ListofEmptyShouldMatchCapnp(t *testing.T) {

	exp := CapnpEncode("(mylist = [(),()])", "HoldsVerEmptyList")

	cv.Convey("Given an empty struct with 0 data/0 ptr fields", t, func() {
		cv.Convey("then a list of 2 empty structs should match the capnp representation", func() {

			seg := capn.NewBuffer(nil)
			scratch := capn.NewBuffer(nil)
			holder := air.NewRootHoldsVerEmptyList(seg)

			emptyBytes := ShowSeg("          after NewRootHoldsVerEmptyList(seg), segment seg is:", seg)

			elist := air.NewVerEmptyList(seg, 2)
			plist := capn.PointerList(elist)

			addList2bytes := ShowSeg("          pre NewVerEmpty(scratch), segment seg is:", seg)
			cv.So(emptyBytes, cv.ShouldResemble, addList2bytes)

			e0 := air.NewVerEmpty(scratch)
			e1 := air.NewVerEmpty(scratch)
			plist.Set(0, capn.Object(e0))
			plist.Set(1, capn.Object(e1))

			ShowSeg("          pre SetMylist, segment seg is:", seg)

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

func TestV1DataVersioningBiggerToEmpty(t *testing.T) {

	//expTwoSet := CapnpEncode("(mylist = [(val = 27, duo = 26),(val = 42, duo = 41)])", "HoldsVerTwoDataList")
	//expOneDataOneDefault := CapnpEncode("(mylist = [(val = 27, duo = 0),(val = 42, duo = 0)])", "HoldsVerTwoDataList")
	//expTwoEmpty := CapnpEncode("(mylist = [(),()])", "HoldsVerTwoDataList")

	//expEmpty := CapnpEncode("(mylist = [(),()])", "HoldsVerEmptyList")
	//expOne := CapnpEncode("(mylist = [(val = 27),(val = 42)])", "HoldsVerOneDataList")

	cv.Convey("Given a struct with 0 data/0 ptr fields, and a newer version of the struct with 2 data fields", t, func() {
		cv.Convey("then reading serialized bigger-struct-list into the smaller (empty or one data-member) list should work, truncating/ignoring the new fields", func() {

			seg := capn.NewBuffer(nil)
			scratch := capn.NewBuffer(nil)
			holder := air.NewRootHoldsVerTwoDataList(seg)

			twolist := air.NewVerTwoDataList(scratch, 2)
			plist := capn.PointerList(twolist)

			d0 := air.NewVerTwoData(scratch)
			d0.SetVal(27)
			d0.SetDuo(26)
			d1 := air.NewVerTwoData(scratch)
			d1.SetVal(42)
			d1.SetDuo(41)
			plist.Set(0, capn.Object(d0))
			plist.Set(1, capn.Object(d1))

			holder.SetMylist(twolist)

			ShowSeg("     before serializing out, segment scratch is:", scratch)
			ShowSeg("     before serializing out, segment seg is:", seg)

			// serialize out
			buf := bytes.Buffer{}
			seg.WriteTo(&buf)
			segbytes := buf.Bytes()

			// and read-back in using smaller expectations
			reseg, _, err := capn.ReadFromMemoryZeroCopy(segbytes)
			if err != nil {
				panic(err)
			}
			ShowSeg("      after re-reading segbytes, segment reseg is:", reseg)
			fmt.Printf("segbytes decoded by capnp as HoldsVerEmptyList: '%s'\n", string(CapnpDecode(segbytes, "HoldsVerEmptyList")))
			fmt.Printf("segbytes decoded by capnp as HoldsVerOneDataList: '%s'\n", string(CapnpDecode(segbytes, "HoldsVerOneDataList")))
			fmt.Printf("segbytes decoded by capnp as HoldsVerTwoDataList: '%s'\n", string(CapnpDecode(segbytes, "HoldsVerTwoDataList")))

			reHolder := air.ReadRootHoldsVerEmptyList(reseg)
			elist := reHolder.Mylist()
			lene := elist.Len()
			cv.So(lene, cv.ShouldEqual, 2)

			reHolder1 := air.ReadRootHoldsVerOneDataList(reseg)
			onelist := reHolder1.Mylist()
			lenone := onelist.Len()
			cv.So(lenone, cv.ShouldEqual, 2)

			for i := 0; i < 2; i++ {
				ele := onelist.At(i)
				val := ele.Val()
				cv.So(val, cv.ShouldEqual, twolist.At(i).Val())
			}

			reHolder2 := air.ReadRootHoldsVerTwoDataList(reseg)
			twolist2 := reHolder2.Mylist()
			lentwo2 := twolist2.Len()
			cv.So(lentwo2, cv.ShouldEqual, 2)

			for i := 0; i < 2; i++ {
				ele := twolist2.At(i)
				val := ele.Val()
				duo := ele.Duo()
				cv.So(val, cv.ShouldEqual, twolist.At(i).Val())
				cv.So(duo, cv.ShouldEqual, twolist.At(i).Duo())
			}

		})
	})
}

func TestV1DataVersioningEmptyToBigger(t *testing.T) {

	//expOneSet := CapnpEncode("(mylist = [(val = 27),(val = 42)])", "HoldsVerOneDataList")
	//expOneZeroed := CapnpEncode("(mylist = [(val = 0),(val = 0)])", "HoldsVerOneDataList")
	//expOneEmpty := CapnpEncode("(mylist = [(),()])", "HoldsVerOneDataList")
	expEmpty := CapnpEncode("(mylist = [(),()])", "HoldsVerEmptyList")

	cv.Convey("Given a struct with 0 data/0 ptr fields, and a newer version of the struct with 1 data fields", t, func() {
		cv.Convey("then reading from serialized form the small list into the bigger (one or two data values) list should work, getting default value 0 for val/duo.", func() {

			seg := capn.NewBuffer(nil)
			scratch := capn.NewBuffer(nil)

			emptyholder := air.NewRootHoldsVerEmptyList(seg)
			elist := air.NewVerEmptyList(scratch, 2)
			emptyholder.SetMylist(elist)

			actEmpty := ShowSeg("          after NewRootHoldsVerEmptyList(seg) and SetMylist(elist), segment seg is:", seg)
			actEmptyCap := string(CapnpDecode(actEmpty, "HoldsVerEmptyList"))
			expEmptyCap := string(CapnpDecode(expEmpty, "HoldsVerEmptyList"))
			cv.So(actEmptyCap, cv.ShouldResemble, expEmptyCap)

			fmt.Printf("\n actEmpty is \n")
			ShowBytes(actEmpty, 10)
			fmt.Printf("actEmpty decoded by capnp: '%s'\n", string(CapnpDecode(actEmpty, "HoldsVerEmptyList")))
			cv.So(actEmpty, cv.ShouldResemble, expEmpty)

			// seg is set, now read into bigger list
			buf := bytes.Buffer{}
			seg.WriteTo(&buf)
			segbytes := buf.Bytes()

			reseg, _, err := capn.ReadFromMemoryZeroCopy(segbytes)
			if err != nil {
				panic(err)
			}
			ShowSeg("      after re-reading segbytes, segment reseg is:", reseg)
			fmt.Printf("segbytes decoded by capnp as HoldsVerOneDataList: '%s'\n", string(CapnpDecode(segbytes, "HoldsVerOneDataList")))

			reHolder := air.ReadRootHoldsVerOneDataList(reseg)
			onelist := reHolder.Mylist()
			lenone := onelist.Len()
			cv.So(lenone, cv.ShouldEqual, 2)
			for i := 0; i < 2; i++ {
				ele := onelist.At(i)
				val := ele.Val()
				cv.So(val, cv.ShouldEqual, 0)
			}

			reHolder2 := air.ReadRootHoldsVerTwoDataList(reseg)
			twolist := reHolder2.Mylist()
			lentwo := twolist.Len()
			cv.So(lentwo, cv.ShouldEqual, 2)
			for i := 0; i < 2; i++ {
				ele := twolist.At(i)
				val := ele.Val()
				cv.So(val, cv.ShouldEqual, 0)
				duo := ele.Duo()
				cv.So(duo, cv.ShouldEqual, 0)
			}

		})
	})
}

func TestDataVersioningZeroPointersToMore(t *testing.T) {

	expEmpty := CapnpEncode("(mylist = [(),()])", "HoldsVerEmptyList")

	cv.Convey("Given a struct with 0 ptr fields, and a newer version of the struct with 1-2 pointer fields", t, func() {
		cv.Convey("then serializing the empty list and reading it back into 1 or 2 pointer fields should default initialize the pointer fields", func() {

			seg := capn.NewBuffer(nil)
			scratch := capn.NewBuffer(nil)

			emptyholder := air.NewRootHoldsVerEmptyList(seg)
			elist := air.NewVerEmptyList(scratch, 2)
			emptyholder.SetMylist(elist)

			actEmpty := ShowSeg("          after NewRootHoldsVerEmptyList(seg) and SetMylist(elist), segment seg is:", seg)
			actEmptyCap := string(CapnpDecode(actEmpty, "HoldsVerEmptyList"))
			expEmptyCap := string(CapnpDecode(expEmpty, "HoldsVerEmptyList"))
			cv.So(actEmptyCap, cv.ShouldResemble, expEmptyCap)

			fmt.Printf("\n actEmpty is \n")
			ShowBytes(actEmpty, 10)
			fmt.Printf("actEmpty decoded by capnp: '%s'\n", string(CapnpDecode(actEmpty, "HoldsVerEmptyList")))
			cv.So(actEmpty, cv.ShouldResemble, expEmpty)

			// seg is set, now read into bigger list
			buf := bytes.Buffer{}
			seg.WriteTo(&buf)
			segbytes := buf.Bytes()

			reseg, _, err := capn.ReadFromMemoryZeroCopy(segbytes)
			if err != nil {
				panic(err)
			}
			ShowSeg("      after re-reading segbytes, segment reseg is:", reseg)
			fmt.Printf("segbytes decoded by capnp as HoldsVerOneDataList: '%s'\n", string(CapnpDecode(segbytes, "HoldsVerOneDataList")))

			reHolder := air.ReadRootHoldsVerTwoTwoList(reseg)
			list22 := reHolder.Mylist()
			len22 := list22.Len()
			cv.So(len22, cv.ShouldEqual, 2)
			for i := 0; i < 2; i++ {
				ele := list22.At(i)
				val := ele.Val()
				cv.So(val, cv.ShouldEqual, 0)
				duo := ele.Duo()
				cv.So(duo, cv.ShouldEqual, 0)
				ptr1 := ele.Ptr1()
				ptr2 := ele.Ptr2()
				fmt.Printf("ptr1 = %#v\n", ptr1)
				cv.So(ptr1.Segment, cv.ShouldEqual, nil)
				fmt.Printf("ptr2 = %#v\n", ptr2)
				cv.So(ptr2.Segment, cv.ShouldEqual, nil)
			}

		})
	})
}

func TestDataVersioningZeroPointersToTwo(t *testing.T) {
	cv.Convey("Given a struct with 2 ptr fields, and another version of the struct with 0 or 1 pointer fields", t, func() {
		cv.Convey("then reading serialized bigger-struct-list into the smaller (empty or one data-pointer) list should work, truncating/ignoring the new fields", func() {

			seg := capn.NewBuffer(nil)
			scratch := capn.NewBuffer(nil)
			holder := air.NewRootHoldsVerTwoTwoList(seg)

			twolist := air.NewVerTwoDataTwoPtrList(scratch, 2)
			plist := capn.PointerList(twolist)

			d0 := air.NewVerTwoDataTwoPtr(scratch)
			d0.SetVal(27)
			d0.SetDuo(26)

			v1 := air.NewVerOneData(scratch)
			v1.SetVal(25)
			v2 := air.NewVerOneData(scratch)
			v2.SetVal(23)

			d0.SetPtr1(v1)
			d0.SetPtr2(v2)

			d1 := air.NewVerTwoDataTwoPtr(scratch)
			d1.SetVal(42)
			d1.SetDuo(41)

			w1 := air.NewVerOneData(scratch)
			w1.SetVal(40)
			w2 := air.NewVerOneData(scratch)
			w2.SetVal(38)

			d1.SetPtr1(w1)
			d1.SetPtr2(w2)

			plist.Set(0, capn.Object(d0))
			plist.Set(1, capn.Object(d1))

			holder.SetMylist(twolist)

			ShowSeg("     before serializing out, segment scratch is:", scratch)
			ShowSeg("     before serializing out, segment seg is:", seg)

			// serialize out
			buf := bytes.Buffer{}
			seg.WriteTo(&buf)
			segbytes := buf.Bytes()

			// and read-back in using smaller expectations
			reseg, _, err := capn.ReadFromMemoryZeroCopy(segbytes)
			if err != nil {
				panic(err)
			}
			ShowSeg("      after re-reading segbytes, segment reseg is:", reseg)
			fmt.Printf("segbytes decoded by capnp as HoldsVerEmptyList: '%s'\n", string(CapnpDecode(segbytes, "HoldsVerEmptyList")))
			fmt.Printf("segbytes decoded by capnp as HoldsVerOnePtrList: '%s'\n", string(CapnpDecode(segbytes, "HoldsVerOnePtrList")))
			fmt.Printf("segbytes decoded by capnp as HoldsVerTwoTwoList: '%s'\n", string(CapnpDecode(segbytes, "HoldsVerTwoTwoList")))

			reHolder := air.ReadRootHoldsVerEmptyList(reseg)
			elist := reHolder.Mylist()
			lene := elist.Len()
			cv.So(lene, cv.ShouldEqual, 2)

			reHolder1 := air.ReadRootHoldsVerOnePtrList(reseg)
			onelist := reHolder1.Mylist()
			lenone := onelist.Len()
			cv.So(lenone, cv.ShouldEqual, 2)

			for i := 0; i < 2; i++ {
				ele := onelist.At(i)
				ptr1 := ele.Ptr()
				cv.So(ptr1.Val(), cv.ShouldEqual, twolist.At(i).Ptr1().Val())
			}

			reHolder2 := air.ReadRootHoldsVerTwoTwoPlus(reseg)
			twolist2 := reHolder2.Mylist()
			lentwo2 := twolist2.Len()
			cv.So(lentwo2, cv.ShouldEqual, 2)

			for i := 0; i < 2; i++ {
				ele := twolist2.At(i)
				ptr1 := ele.Ptr1()
				ptr2 := ele.Ptr2()
				cv.So(ptr1.Val(), cv.ShouldEqual, twolist.At(i).Ptr1().Val())
				//cv.So(ptr1.Duo(), cv.ShouldEqual, twolist.At(i).Ptr1().Duo())
				cv.So(ptr2.Val(), cv.ShouldEqual, twolist.At(i).Ptr2().Val())
				//cv.So(ptr2.Duo(), cv.ShouldEqual, twolist.At(i).Ptr2().Duo())
				cv.So(ele.Tre(), cv.ShouldEqual, 0)
				cv.So(ele.Lst3().Len(), cv.ShouldEqual, 0)
			}

		})
	})
}
