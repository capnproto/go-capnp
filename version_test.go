package capnp_test

import (
	"fmt"
	"testing"

	cv "github.com/smartystreets/goconvey/convey"
	"zombiezen.com/go/capnproto"
	air "zombiezen.com/go/capnproto/internal/aircraftlib"
)

func TestV0ListofEmptyShouldMatchCapnp(t *testing.T) {

	exp := CapnpEncode("(mylist = [(),()])", "HoldsVerEmptyList")

	cv.Convey("Given an empty struct with 0 data/0 ptr fields", t, func() {
		cv.Convey("then a list of 2 empty structs should match the capnp representation", func() {

			msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)
			holder, err := air.NewRootHoldsVerEmptyList(seg)
			cv.So(err, cv.ShouldEqual, nil)

			elist, err := air.NewVerEmpty_List(seg, 2)
			cv.So(err, cv.ShouldEqual, nil)

			ShowSeg("          pre SetMylist, segment seg is:", seg)

			fmt.Printf("Then we do the SetMylist():\n")
			holder.SetMylist(elist)

			// save
			act, err := msg.Marshal()
			cv.So(err, cv.ShouldEqual, nil)
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

			msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)
			_, scratch, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)
			holder, err := air.NewRootHoldsVerTwoDataList(seg)
			cv.So(err, cv.ShouldEqual, nil)

			twolist, err := air.NewVerTwoData_List(scratch, 2)
			cv.So(err, cv.ShouldEqual, nil)

			d0 := twolist.At(0)
			d0.SetVal(27)
			d0.SetDuo(26)
			d1 := twolist.At(1)
			d1.SetVal(42)
			d1.SetDuo(41)

			holder.SetMylist(twolist)

			ShowSeg("     before serializing out, segment scratch is:", scratch)
			ShowSeg("     before serializing out, segment seg is:", seg)

			// serialize out
			segbytes, err := msg.Marshal()
			cv.So(err, cv.ShouldEqual, nil)

			// and read-back in using smaller expectations
			remsg, err := capnp.Unmarshal(segbytes)
			cv.So(err, cv.ShouldEqual, nil)
			reseg, err := remsg.Segment(0)
			cv.So(err, cv.ShouldEqual, nil)
			ShowSeg("      after re-reading segbytes, segment reseg is:", reseg)
			fmt.Printf("segbytes decoded by capnp as HoldsVerEmptyList: '%s'\n", string(CapnpDecode(segbytes, "HoldsVerEmptyList")))
			fmt.Printf("segbytes decoded by capnp as HoldsVerOneDataList: '%s'\n", string(CapnpDecode(segbytes, "HoldsVerOneDataList")))
			fmt.Printf("segbytes decoded by capnp as HoldsVerTwoDataList: '%s'\n", string(CapnpDecode(segbytes, "HoldsVerTwoDataList")))

			reHolder, err := air.ReadRootHoldsVerEmptyList(remsg)
			cv.So(err, cv.ShouldEqual, nil)
			elist, err := reHolder.Mylist()
			cv.So(err, cv.ShouldEqual, nil)
			lene := elist.Len()
			cv.So(lene, cv.ShouldEqual, 2)

			reHolder1, err := air.ReadRootHoldsVerOneDataList(remsg)
			cv.So(err, cv.ShouldEqual, nil)
			onelist, err := reHolder1.Mylist()
			cv.So(err, cv.ShouldEqual, nil)
			lenone := onelist.Len()
			cv.So(lenone, cv.ShouldEqual, 2)

			for i := 0; i < 2; i++ {
				ele := onelist.At(i)
				val := ele.Val()
				cv.So(val, cv.ShouldEqual, twolist.At(i).Val())
			}

			reHolder2, err := air.ReadRootHoldsVerTwoDataList(remsg)
			cv.So(err, cv.ShouldEqual, nil)
			twolist2, err := reHolder2.Mylist()
			cv.So(err, cv.ShouldEqual, nil)
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

			msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)
			_, scratch, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)

			emptyholder, err := air.NewRootHoldsVerEmptyList(seg)
			cv.So(err, cv.ShouldEqual, nil)
			elist, err := air.NewVerEmpty_List(scratch, 2)
			cv.So(err, cv.ShouldEqual, nil)
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
			segbytes, err := msg.Marshal()
			cv.So(err, cv.ShouldEqual, nil)

			remsg, err := capnp.Unmarshal(segbytes)
			cv.So(err, cv.ShouldEqual, nil)
			reseg, err := remsg.Segment(0)
			cv.So(err, cv.ShouldEqual, nil)
			ShowSeg("      after re-reading segbytes, segment reseg is:", reseg)
			fmt.Printf("segbytes decoded by capnp as HoldsVerOneDataList: '%s'\n", string(CapnpDecode(segbytes, "HoldsVerOneDataList")))

			reHolder, err := air.ReadRootHoldsVerOneDataList(remsg)
			cv.So(err, cv.ShouldEqual, nil)
			onelist, err := reHolder.Mylist()
			cv.So(err, cv.ShouldEqual, nil)
			lenone := onelist.Len()
			cv.So(lenone, cv.ShouldEqual, 2)
			for i := 0; i < 2; i++ {
				ele := onelist.At(i)
				val := ele.Val()
				cv.So(val, cv.ShouldEqual, 0)
			}

			reHolder2, err := air.ReadRootHoldsVerTwoDataList(remsg)
			cv.So(err, cv.ShouldEqual, nil)
			twolist, err := reHolder2.Mylist()
			cv.So(err, cv.ShouldEqual, nil)
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

			msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)
			_, scratch, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)

			emptyholder, err := air.NewRootHoldsVerEmptyList(seg)
			cv.So(err, cv.ShouldEqual, nil)
			elist, err := air.NewVerEmpty_List(scratch, 2)
			cv.So(err, cv.ShouldEqual, nil)
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
			segbytes, err := msg.Marshal()
			cv.So(err, cv.ShouldEqual, nil)

			remsg, err := capnp.Unmarshal(segbytes)
			cv.So(err, cv.ShouldEqual, nil)
			reseg, err := remsg.Segment(0)
			cv.So(err, cv.ShouldEqual, nil)
			ShowSeg("      after re-reading segbytes, segment reseg is:", reseg)
			fmt.Printf("segbytes decoded by capnp as HoldsVerOneDataList: '%s'\n", string(CapnpDecode(segbytes, "HoldsVerOneDataList")))

			reHolder, err := air.ReadRootHoldsVerTwoTwoList(remsg)
			cv.So(err, cv.ShouldEqual, nil)
			list22, err := reHolder.Mylist()
			cv.So(err, cv.ShouldEqual, nil)
			len22 := list22.Len()
			cv.So(len22, cv.ShouldEqual, 2)
			for i := 0; i < 2; i++ {
				ele := list22.At(i)
				val := ele.Val()
				cv.So(val, cv.ShouldEqual, 0)
				duo := ele.Duo()
				cv.So(duo, cv.ShouldEqual, 0)
				ptr1, err := ele.Ptr1()
				cv.So(err, cv.ShouldEqual, nil)
				ptr2, err := ele.Ptr2()
				cv.So(err, cv.ShouldEqual, nil)
				fmt.Printf("ptr1 = %#v\n", ptr1)
				cv.So(ptr1.Segment(), cv.ShouldEqual, nil)
				fmt.Printf("ptr2 = %#v\n", ptr2)
				cv.So(ptr2.Segment(), cv.ShouldEqual, nil)
			}

		})
	})
}

func TestDataVersioningZeroPointersToTwo(t *testing.T) {
	cv.Convey("Given a struct with 2 ptr fields, and another version of the struct with 0 or 1 pointer fields", t, func() {
		cv.Convey("then reading serialized bigger-struct-list into the smaller (empty or one data-pointer) list should work, truncating/ignoring the new fields", func() {

			msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)
			_, scratch, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)
			holder, err := air.NewRootHoldsVerTwoTwoList(seg)
			cv.So(err, cv.ShouldEqual, nil)

			twolist, err := air.NewVerTwoDataTwoPtr_List(scratch, 2)
			cv.So(err, cv.ShouldEqual, nil)

			d0 := twolist.At(0)
			d0.SetVal(27)
			d0.SetDuo(26)

			v1, err := air.NewVerOneData(scratch)
			cv.So(err, cv.ShouldEqual, nil)
			v1.SetVal(25)
			v2, err := air.NewVerOneData(scratch)
			cv.So(err, cv.ShouldEqual, nil)
			v2.SetVal(23)

			d0.SetPtr1(v1)
			d0.SetPtr2(v2)

			d1 := twolist.At(1)
			d1.SetVal(42)
			d1.SetDuo(41)

			w1, err := air.NewVerOneData(scratch)
			cv.So(err, cv.ShouldEqual, nil)
			w1.SetVal(40)
			w2, err := air.NewVerOneData(scratch)
			cv.So(err, cv.ShouldEqual, nil)
			w2.SetVal(38)

			d1.SetPtr1(w1)
			d1.SetPtr2(w2)

			holder.SetMylist(twolist)

			ShowSeg("     before serializing out, segment scratch is:", scratch)
			ShowSeg("     before serializing out, segment seg is:", seg)

			// serialize out
			segbytes, err := msg.Marshal()
			cv.So(err, cv.ShouldEqual, nil)

			// and read-back in using smaller expectations
			remsg, err := capnp.Unmarshal(segbytes)
			cv.So(err, cv.ShouldEqual, nil)
			reseg, err := remsg.Segment(0)
			cv.So(err, cv.ShouldEqual, nil)
			ShowSeg("      after re-reading segbytes, segment reseg is:", reseg)
			fmt.Printf("segbytes decoded by capnp as HoldsVerEmptyList: '%s'\n", string(CapnpDecode(segbytes, "HoldsVerEmptyList")))
			fmt.Printf("segbytes decoded by capnp as HoldsVerOnePtrList: '%s'\n", string(CapnpDecode(segbytes, "HoldsVerOnePtrList")))
			fmt.Printf("segbytes decoded by capnp as HoldsVerTwoTwoList: '%s'\n", string(CapnpDecode(segbytes, "HoldsVerTwoTwoList")))

			reHolder, err := air.ReadRootHoldsVerEmptyList(remsg)
			elist, err := reHolder.Mylist()
			cv.So(err, cv.ShouldEqual, nil)
			lene := elist.Len()
			cv.So(lene, cv.ShouldEqual, 2)

			reHolder1, err := air.ReadRootHoldsVerOnePtrList(remsg)
			onelist, err := reHolder1.Mylist()
			cv.So(err, cv.ShouldEqual, nil)
			lenone := onelist.Len()
			cv.So(lenone, cv.ShouldEqual, 2)

			for i := 0; i < 2; i++ {
				ele := onelist.At(i)
				ptr1, err := ele.Ptr()
				cv.So(err, cv.ShouldEqual, nil)
				origPtr1, err := twolist.At(i).Ptr1()
				cv.So(err, cv.ShouldEqual, nil)
				cv.So(ptr1.Val(), cv.ShouldEqual, origPtr1.Val())
			}

			reHolder2, err := air.ReadRootHoldsVerTwoTwoPlus(remsg)
			cv.So(err, cv.ShouldEqual, nil)
			twolist2, err := reHolder2.Mylist()
			cv.So(err, cv.ShouldEqual, nil)
			lentwo2 := twolist2.Len()
			cv.So(lentwo2, cv.ShouldEqual, 2)

			for i := 0; i < 2; i++ {
				ele := twolist2.At(i)
				ptr1, err := ele.Ptr1()
				cv.So(err, cv.ShouldEqual, nil)
				ptr2, err := ele.Ptr2()
				cv.So(err, cv.ShouldEqual, nil)
				origPtr1, err := ele.Ptr1()
				cv.So(err, cv.ShouldEqual, nil)
				origPtr2, err := ele.Ptr2()
				cv.So(err, cv.ShouldEqual, nil)
				cv.So(ptr1.Val(), cv.ShouldEqual, origPtr1.Val())
				//cv.So(ptr1.Duo(), cv.ShouldEqual, twolist.At(i).Ptr1().Duo())
				cv.So(ptr2.Val(), cv.ShouldEqual, origPtr2.Val())
				//cv.So(ptr2.Duo(), cv.ShouldEqual, twolist.At(i).Ptr2().Duo())
				cv.So(ele.Tre(), cv.ShouldEqual, 0)
				lst3, err := ele.Lst3()
				cv.So(err, cv.ShouldEqual, nil)
				cv.So(lst3.Len(), cv.ShouldEqual, 0)
			}

		})
	})
}
