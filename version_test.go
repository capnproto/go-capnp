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
		cv.Convey("then reading from serialized form the small list into the bigger (one data value) list should work, getting default value 0 for val.", func() {

			seg := capn.NewBuffer(nil)
			scratch := capn.NewBuffer(nil)
			//scratch2 := capn.NewBuffer(nil)

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

/*
func TestDataVersioningOneDataToTwo(t *testing.T) {

	expTwoSet := CapnpEncode("(mylist = [(val = 27, duo = 26),(val = 42, duo = 41)])", "HoldsVerTwoDataList")
	expOneDataOneDefault := CapnpEncode("(mylist = [(val = 27, duo = 0),(val = 42, duo = 0)])", "HoldsVerTwoDataList")
	expTwoEmpty := CapnpEncode("(mylist = [(),()])", "HoldsVerTwoDataList")

	cv.Convey("Given a smaller struct with 1 data/0 ptr fields, and a newer, bigger version of the struct with 2 data/0 pointer fields", t, func() {
		cv.Convey("then setting the small struct list into the a bigger struct list should work, with defaults for the 2nd bigger value", func() {

			seg := capn.NewBuffer(nil)
			scratch := capn.NewBuffer(nil)
			scratch2 := capn.NewBuffer(nil)
			holder := air.NewRootHoldsVerTwoDataList(seg)

			ShowSeg("          after NewRootHoldsVerTwoDataList(seg), segment seg is:", seg)

			twolist := air.NewVerTwoDataList(seg, 2)
			ShowSeg("          after air.NewVerTwoDataList(seg, 2), segment seg is:", seg)

			holder.SetMylist(twolist)

			actTwoEmpty := ShowSeg("     after NewVerTwoDataList(seg), segment seg is:", seg)

			fmt.Printf("\n expTwoEmpty is \n")
			ShowBytes(expTwoEmpty, 10)
			fmt.Printf("expTwoEmpty decoded by capnp: '%s'\n", string(CapnpDecode(expTwoEmpty, "HoldsVerTwoDataList")))
			cv.So(actTwoEmpty, cv.ShouldResemble, expTwoEmpty)

			plist := capn.PointerList(twolist)

			d0 := air.NewVerTwoData(scratch)
			d0.SetVal(27)
			d0.SetDuo(26)
			d1 := air.NewVerTwoData(scratch)
			d1.SetVal(42)
			d1.SetDuo(41)

			ShowSeg("     after setting values {27,26}, {42,41} on VerTwoData structs in scratch, segment scratch is:", scratch)

			// these are copies from scratch to seg
			plist.Set(0, capn.Object(d0))
			plist.Set(1, capn.Object(d1))

			ShowSeg("          pre SetMylist(twolist), segment scratch is:", scratch)
			myseg := ShowSeg("          pre SetMylist(twolist), segment seg is:", seg)
			fmt.Printf("seg decoded by capnp: '%s'\n", string(CapnpDecode(myseg, "HoldsVerTwoDataList")))

			fmt.Printf("Then we do the SetMylist():\n")
			holder.SetMylist(twolist)

			actTwoSet := ShowSeg("          post SetMylist(twolist), segment seg should have twolist in it: actTwoSet = ", seg)

			fmt.Printf("\n expTwoSet is:\n")
			ShowBytes(expTwoSet, 10)
			fmt.Printf("expTwoSet decoded by capnp: '%s'\n", string(CapnpDecode(expTwoSet, "HoldsVerTwoDataList")))

			cv.So(actTwoSet, cv.ShouldResemble, expTwoSet)

			// now we have the bigger list in seg.

			// the test is: assign the smaller list (from a totally seperate segment, scratch2) to the bigger. (Version checking).
			// the values of the bigger list should get zero-ed.

			onelist := air.NewVerOneDataList(scratch2, 2)

			p1list := capn.PointerList(onelist)

			a0 := air.NewVerOneData(scratch2)
			a0.SetVal(27)
			a1 := air.NewVerOneData(scratch2)
			a1.SetVal(42)

			p1list.Set(0, capn.Object(a0))
			p1list.Set(1, capn.Object(a1))

			ShowSeg("     after setting values {27}, {42} on VerOneData structs in scratch, segment scratch2 is:", scratch2)

			fmt.Printf("onelist = %#v\n", onelist)
			fmt.Printf("twolist = %#v\n", twolist)

			holder.SetMylist(air.VerTwoData_List(onelist)) // cast required to simulate data version skew
			act := ShowSeg("          post SetMylist(onelist), segment seg should have an the two-data list with default values for 2nd (duo values) in it:", seg)

			actOneDataOneDefault := string(CapnpDecode(act, "HoldsVerTwoDataList"))
			fmt.Printf("act decoded by capnp: '%s'\n", actOneDataOneDefault)

			fmt.Printf("expOneDataOneDefault / expected:\n")
			ShowBytes(expOneDataOneDefault, 10)
			expOneDataOneDefaultCap := string(CapnpDecode(expOneDataOneDefault, "HoldsVerTwoDataList"))
			fmt.Printf("expOneDataOneDefault decoded by capnp: '%s'\n", expOneDataOneDefaultCap)

			// can't compare binaries, because binaries actually aren't actually expected/necessarily identical here.
			// cv.So(act, cv.ShouldResemble, expOneDataOneDefault)
			// but we can compare re-textifications:

			cv.So(actOneDataOneDefault, cv.ShouldResemble, expOneDataOneDefaultCap)

		})
	})
}


func TestDataVersioningTwoDataToZero(t *testing.T) {

	expTwoSet := CapnpEncode("(mylist = [(val = 27, duo = 26),(val = 42, duo = 41)])", "HoldsVerTwoDataList")
	expOneDataOneDefault := CapnpEncode("(mylist = [(val = 27, duo = 0),(val = 42, duo = 0)])", "HoldsVerTwoDataList")

	cv.Convey("Given a smaller struct with 1 data/0 ptr fields, and a newer, bigger version of the struct with 2 data/0 pointer fields", t, func() {
		cv.Convey("then setting the bigger struct list into the a smaller struct list should work, with truncation of the 2nd values", func() {

			seg := capn.NewBuffer(nil)
			scratch := capn.NewBuffer(nil)
			holder := air.NewRootHoldsVerOneDataList(seg)

			emptyBytes := ShowSeg("          after NewRootHoldsVerEmptyList(seg), segment seg is:", seg)

			elist := air.NewVerEmptyList(seg, 2)
			plist := capn.PointerList(elist)

			addList2bytes := ShowSeg("          pre NewVerEmpty(scratch), segment seg is:", seg)
			cv.So(emptyBytes, cv.ShouldResemble, addList2bytes)

			d0 := air.NewVerTwoData(scratch)
			d0.SetVal(27)
			d0.SetDuo(26)
			d1 := air.NewVerTwoData(scratch)
			d1.SetVal(42)
			d1.SetDuo(41)

			plist.Set(0, capn.Object(d0))
			plist.Set(1, capn.Object(d1))

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

			fmt.Printf("act decoded by capnp: '%s'\n", string(CapnpDecode(act, "HoldsVerTwoDataList")))

			fmt.Printf("expEmpty / expected:\n")
			ShowBytes(expEmpty, 10)
			fmt.Printf("expEmpty decoded by capnp: '%s'\n", string(CapnpDecode(expEmpty, "HoldsVerTwoDataList")))

			fmt.Printf("expTwo decoded by capnp: '%s'\n", string(CapnpDecode(expTwo, "HoldsVerTwoDataList")))

			cv.So(act, cv.ShouldResemble, expEmpty)

		})
	})
}
*/

/*
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

/*
	onelist := air.NewVerOneDataList(seg, 2)
	ShowSeg("          after air.NewVerOneDataList(seg, 2), segment seg is:", seg)

	holder.SetMylist(onelist)

	actOneEmpty := ShowSeg("     after NewVerOneDataList(seg), segment seg is:", seg)

	fmt.Printf("\n expOneEmpty is \n")
	ShowBytes(expOneEmpty, 10)
	fmt.Printf("expOneEmpty decoded by capnp: '%s'\n", string(CapnpDecode(expOneEmpty, "HoldsVerOneDataList")))
	cv.So(actOneEmpty, cv.ShouldResemble, expOneEmpty)

	plist := capn.PointerList(onelist)

	d0 := air.NewVerOneData(scratch)
	d0.SetVal(27)
	d1 := air.NewVerOneData(scratch)
	d1.SetVal(42)

	ShowSeg("     after setting values 27, 42 on VerOneData structs in scratch, segment scratch is:", scratch)

	// these are copies from scratch to seg
	plist.Set(0, capn.Object(d0))
	plist.Set(1, capn.Object(d1))

	ShowSeg("          pre SetMylist(onelist), segment scratch is:", scratch)
	ShowSeg("          pre SetMylist(onelist), segment seg is:", seg)

	fmt.Printf("Then we do the SetMylist():\n")
	holder.SetMylist(onelist)

	actOneSet := ShowSeg("          post SetMylist(onelist), segment seg should have onelist in it: actOneSet = ", seg)

	fmt.Printf("\n expOneSet is:\n")
	ShowBytes(expOneSet, 10)
	fmt.Printf("expOneSet decoded by capnp: '%s'\n", string(CapnpDecode(expOneSet, "HoldsVerOneDataList")))

	cv.So(actOneSet, cv.ShouldResemble, expOneSet)

	// now we have the bigger list in seg.

	// the test is: read seg back into the smaller list
	// the values of the bigger list should get zero-ed.

	air.ReadRootHoldsVerOneDataList(seg)

	elist := air.NewVerEmptyList(scratch2, 2)
	holder.SetMylist(air.VerOneData_List(elist)) // cast required to simulate data version skew
	act := ShowSeg("          post SetMylist(elist), segment seg should have an empty list in it:", seg)

	actOneZeroedCap := string(CapnpDecode(act, "HoldsVerOneDataList"))
	fmt.Printf("act decoded by capnp: '%s'\n", actOneZeroedCap)

	fmt.Printf("expOneZeroed / expected:\n")
	ShowBytes(expOneZeroed, 10)
	expOneZeroedCap := string(CapnpDecode(expOneZeroed, "HoldsVerOneDataList"))
	fmt.Printf("expOneZeroed decoded by capnp: '%s'\n", expOneZeroedCap)

	// can't compare binaries, because binaries actually aren't actually expected/necessarily identical here.
	cv.So(act, cv.ShouldResemble, expOneZeroed)
	// but we can compare re-textifications:

	cv.So(actOneZeroedCap, cv.ShouldResemble, expOneZeroedCap)
*/
