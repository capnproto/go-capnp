package capnp_test

import (
	"fmt"
	"testing"

	cv "github.com/smartystreets/goconvey/convey"
	"zombiezen.com/go/capnproto"
	air "zombiezen.com/go/capnproto/internal/aircraftlib"
)

func TestTextAndListTextContaintingEmptyStruct(t *testing.T) {

	emptyZjobBytes := CapnpEncode("()", "Zjob")

	cv.Convey("Given a simple struct message Zjob containing a string and a list of string (all empty)", t, func() {
		cv.Convey("then the go-capnproto serialization should match the capnp c++ serialization", func() {
			ShowBytes(emptyZjobBytes, 10)

			msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)
			air.NewRootZjob(seg)

			buf, err := msg.Marshal()
			cv.So(err, cv.ShouldEqual, nil)

			cv.So(buf, cv.ShouldResemble, emptyZjobBytes)
		})
	})
}

func TestTextContaintingStruct(t *testing.T) {

	zjobBytes := CapnpEncode(`(cmd = "abc")`, "Zjob")

	cv.Convey("Given a simple struct message Zjob containing a string 'abc' and a list of string (empty)", t, func() {
		cv.Convey("then the go-capnproto serialization should match the capnp c++ serialization", func() {

			msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)
			zjob, err := air.NewRootZjob(seg)
			cv.So(err, cv.ShouldEqual, nil)
			zjob.SetCmd("abc")

			act, err := msg.Marshal()
			cv.So(err, cv.ShouldEqual, nil)

			fmt.Printf("          actual:\n")
			ShowBytes(act, 10)

			fmt.Printf("\n\n          expected:\n")
			ShowBytes(zjobBytes, 10)

			cv.So(act, cv.ShouldResemble, zjobBytes)
		})
	})
}

func TestTextListContaintingStruct(t *testing.T) {

	zjobBytes := CapnpEncode(`(args = ["xyz"])`, "Zjob")

	cv.Convey("Given a simple struct message Zjob containing an unset string and a list of string ('xyz' as the only element)", t, func() {
		cv.Convey("then the go-capnproto serialization should match the capnp c++ serialization", func() {

			msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)
			zjob, err := air.NewRootZjob(seg)
			cv.So(err, cv.ShouldEqual, nil)
			tl, err := capnp.NewTextList(seg, 1)
			cv.So(err, cv.ShouldEqual, nil)
			tl.Set(0, "xyz")
			zjob.SetArgs(tl)

			act, err := msg.Marshal()
			cv.So(err, cv.ShouldEqual, nil)

			fmt.Printf("          actual:\n")
			ShowBytes(act, 10)

			fmt.Printf("expected:\n")
			ShowBytes(zjobBytes, 10)

			cv.So(act, cv.ShouldResemble, zjobBytes)
		})
	})
}

func TestTextAndTextListContaintingStruct(t *testing.T) {

	zjobBytes := CapnpEncode(`(cmd = "abc", args = ["xyz"])`, "Zjob")

	cv.Convey("Given a simple struct message Zjob containing a string (cmd='abc') and a list of string (args=['xyz'])", t, func() {
		cv.Convey("then the go-capnproto serialization should match the capnp c++ serialization", func() {

			msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)
			zjob, err := air.NewRootZjob(seg)
			cv.So(err, cv.ShouldEqual, nil)
			zjob.SetCmd("abc")
			tl, err := capnp.NewTextList(seg, 1)
			cv.So(err, cv.ShouldEqual, nil)
			tl.Set(0, "xyz")
			zjob.SetArgs(tl)

			act, err := msg.Marshal()
			cv.So(err, cv.ShouldEqual, nil)

			fmt.Printf("          actual:\n")
			ShowBytes(act, 10)

			fmt.Printf("expected:\n")
			ShowBytes(zjobBytes, 10)

			cv.So(act, cv.ShouldResemble, zjobBytes)
		})
	})
}

func TestZserverWithOneFullJob(t *testing.T) {

	exp := CapnpEncode(`(waitingjobs = [(cmd = "abc", args = ["xyz"])])`, "Zserver")

	cv.Convey("Given an Zserver with one empty job", t, func() {
		cv.Convey("then the go-capnproto serialization should match the capnp c++ serialization", func() {

			msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)
			_, scratch, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)

			server, err := air.NewRootZserver(seg)
			cv.So(err, cv.ShouldEqual, nil)

			joblist, err := air.NewZjob_List(seg, 1)
			cv.So(err, cv.ShouldEqual, nil)

			zjob, err := air.NewZjob(scratch)
			cv.So(err, cv.ShouldEqual, nil)
			zjob.SetCmd("abc")
			tl, err := capnp.NewTextList(scratch, 1)
			cv.So(err, cv.ShouldEqual, nil)
			tl.Set(0, "xyz")
			zjob.SetArgs(tl)

			joblist.Set(0, zjob)

			server.SetWaitingjobs(joblist)

			act, err := msg.Marshal()
			cv.So(err, cv.ShouldEqual, nil)

			fmt.Printf("          actual:\n")
			ShowBytes(act, 10)
			fmt.Printf("act decoded by capnp: '%s'\n", string(CapnpDecode(act, "Zserver")))
			save(act, "myact")

			fmt.Printf("expected:\n")
			ShowBytes(exp, 10)
			fmt.Printf("exp decoded by capnp: '%s'\n", string(CapnpDecode(exp, "Zserver")))
			save(exp, "myexp")

			cv.So(act, cv.ShouldResemble, exp)
		})
	})
}

func TestZserverWithAccessors(t *testing.T) {

	exp := CapnpEncode(`(waitingjobs = [(cmd = "abc"), (cmd = "xyz")])`, "Zserver")

	cv.Convey("Given an Zserver with a custom list", t, func() {
		cv.Convey("then all the accessors should work as expected", func() {

			msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)
			_, scratch, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)

			server, err := air.NewRootZserver(seg)
			cv.So(err, cv.ShouldEqual, nil)

			joblist, err := air.NewZjob_List(seg, 2)
			cv.So(err, cv.ShouldEqual, nil)

			// .Set(int, item)
			zjob, err := air.NewZjob(scratch)
			cv.So(err, cv.ShouldEqual, nil)
			zjob.SetCmd("abc")
			joblist.Set(0, zjob)

			zjob, err = air.NewZjob(scratch)
			cv.So(err, cv.ShouldEqual, nil)
			zjob.SetCmd("xyz")
			joblist.Set(1, zjob)

			// .Len()
			cv.So(joblist.Len(), cv.ShouldEqual, 2)

			// .At(int)
			cmd := func(i int) string {
				s, err := joblist.At(i).Cmd()
				cv.So(err, cv.ShouldEqual, nil)
				return s
			}
			cv.So(cmd(0), cv.ShouldEqual, "abc")
			cv.So(cmd(1), cv.ShouldEqual, "xyz")

			server.SetWaitingjobs(joblist)

			act, err := msg.Marshal()
			cv.So(err, cv.ShouldEqual, nil)

			fmt.Printf("          actual:\n")
			ShowBytes(act, 10)
			fmt.Printf("act decoded by capnp: '%s'\n", string(CapnpDecode(act, "Zserver")))
			save(act, "myact")

			fmt.Printf("expected:\n")
			ShowBytes(exp, 10)
			fmt.Printf("exp decoded by capnp: '%s'\n", string(CapnpDecode(exp, "Zserver")))
			save(exp, "myexp")

			cv.So(act, cv.ShouldResemble, exp)
		})
	})
}

func TestEnumFromString(t *testing.T) {
	cv.Convey("Given an enum tag string matching a constant", t, func() {
		cv.Convey("FromString should return the corresponding matching constant value", func() {
			cv.So(air.AirportFromString("jfk"), cv.ShouldEqual, air.Airport_jfk)
		})
	})
	cv.Convey("Given an enum tag string that does not match a constant", t, func() {
		cv.Convey("FromString should return 0", func() {
			cv.So(air.AirportFromString("notEverMatching"), cv.ShouldEqual, 0)
		})
	})
}

func TestSetObjectBetweenSegments(t *testing.T) {

	exp := CapnpEncode(`(counter = (size = 9))`, "Bag")

	cv.Convey("Given an Counter in one segment and a Bag in another", t, func() {
		cv.Convey("we should be able to copy from one segment to the other with SetCounter() on a Bag", func() {

			msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)
			_, scratch, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)

			// in seg
			segbag, err := air.NewRootBag(seg)
			cv.So(err, cv.ShouldEqual, nil)

			// in scratch
			xc, err := air.NewRootCounter(scratch)
			cv.So(err, cv.ShouldEqual, nil)
			xc.SetSize(9)

			// copy from scratch to seg
			err = segbag.SetCounter(xc)
			cv.So(err, cv.ShouldEqual, nil)

			act, err := msg.Marshal()
			cv.So(err, cv.ShouldEqual, nil)

			fmt.Printf("          actual:\n")
			ShowBytes(act, 10)
			fmt.Printf("act decoded by capnp: '%s'\n", string(CapnpDecode(act, "Bag")))
			save(act, "myact")

			fmt.Printf("expected:\n")
			ShowBytes(exp, 10)
			fmt.Printf("exp decoded by capnp: '%s'\n", string(CapnpDecode(exp, "Bag")))
			save(exp, "myexp")

			cv.So(act, cv.ShouldResemble, exp)
		})
	})
}

func TestObjectWithTextBetweenSegments(t *testing.T) {

	exp := CapnpEncode(`(counter = (size = 9, words = "hello"))`, "Bag")

	cv.Convey("Given an Counter in one segment and a Bag with text in another", t, func() {
		cv.Convey("we should be able to copy from one segment to the other with SetCounter() on a Bag", func() {

			msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)
			_, scratch, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)

			// in seg
			segbag, err := air.NewRootBag(seg)
			cv.So(err, cv.ShouldEqual, nil)

			// in scratch
			xc, err := air.NewRootCounter(scratch)
			cv.So(err, cv.ShouldEqual, nil)
			xc.SetSize(9)
			xc.SetWords("hello")

			// copy from scratch to seg
			err = segbag.SetCounter(xc)
			cv.So(err, cv.ShouldEqual, nil)

			act, err := msg.Marshal()
			cv.So(err, cv.ShouldEqual, nil)

			fmt.Printf("          actual:\n")
			ShowBytes(act, 10)
			fmt.Printf("act decoded by capnp: '%s'\n", string(CapnpDecode(act, "Bag")))
			save(act, "myact")

			fmt.Printf("expected:\n")
			ShowBytes(exp, 10)
			fmt.Printf("exp decoded by capnp: '%s'\n", string(CapnpDecode(exp, "Bag")))
			save(exp, "myexp")

			cv.So(act, cv.ShouldResemble, exp)
		})
	})
}

func TestObjectWithListOfTextBetweenSegments(t *testing.T) {

	exp := CapnpEncode(`(counter = (size = 9, wordlist = ["hello","bye"]))`, "Bag")

	cv.Convey("Given an Counter in one segment and a Bag with text in another", t, func() {
		cv.Convey("we should be able to copy from one segment to the other with SetCounter() on a Bag", func() {

			msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)
			scratchMsg, scratch, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)

			// in seg
			segbag, err := air.NewRootBag(seg)
			cv.So(err, cv.ShouldEqual, nil)

			// in scratch
			xc, err := air.NewRootCounter(scratch)
			cv.So(err, cv.ShouldEqual, nil)
			xc.SetSize(9)
			tl, err := capnp.NewTextList(scratch, 2)
			cv.So(err, cv.ShouldEqual, nil)
			tl.Set(0, "hello")
			tl.Set(1, "bye")
			err = xc.SetWordlist(tl)
			cv.So(err, cv.ShouldEqual, nil)

			x, err := scratchMsg.Marshal()
			cv.So(err, cv.ShouldEqual, nil)

			save(x, "myscratch")
			fmt.Printf("scratch segment (%p):\n", scratch)
			ShowBytes(x, 10)
			fmt.Printf("scratch segment (%p) with Counter decoded by capnp: '%s'\n", scratch, string(CapnpDecode(x, "Counter")))

			pre, err := msg.Marshal()
			cv.So(err, cv.ShouldEqual, nil)
			fmt.Printf("Bag only segment seg (%p), pre-transfer:\n", seg)
			ShowBytes(pre, 10)

			// now for the actual test:
			// copy from scratch to seg
			segbag.SetCounter(xc)

			act, err := msg.Marshal()
			cv.So(err, cv.ShouldEqual, nil)

			save(act, "myact")
			save(exp, "myexp")

			fmt.Printf("expected:\n")
			ShowBytes(exp, 10)
			fmt.Printf("exp decoded by capnp: '%s'\n", string(CapnpDecode(exp, "Bag")))

			fmt.Printf("          actual:\n")
			ShowBytes(act, 10)
			fmt.Printf("act decoded by capnp: '%s'\n", string(CapnpDecode(act, "Bag")))

			cv.So(act, cv.ShouldResemble, exp)
		})
	})
}

func TestSetBetweenSegments(t *testing.T) {

	exp := CapnpEncode(`(counter = (size = 9, words = "abc", wordlist = ["hello","byenow"]))`, "Bag")

	cv.Convey("Given an struct with Text and List(Text) in one segment", t, func() {
		cv.Convey("assigning it to a struct in a different segment should recursively import", func() {

			msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)
			_, scratch, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)

			// in seg
			segbag, err := air.NewRootBag(seg)
			cv.So(err, cv.ShouldEqual, nil)

			// in scratch
			xc, err := air.NewRootCounter(scratch)
			cv.So(err, cv.ShouldEqual, nil)
			xc.SetSize(9)
			tl, err := capnp.NewTextList(scratch, 2)
			cv.So(err, cv.ShouldEqual, nil)
			tl.Set(0, "hello")
			tl.Set(1, "byenow")
			err = xc.SetWordlist(tl)
			cv.So(err, cv.ShouldEqual, nil)
			err = xc.SetWords("abc")
			cv.So(err, cv.ShouldEqual, nil)

			fmt.Printf("\n\n starting copy from scratch to seg \n\n")

			// copy from scratch to seg
			err = segbag.SetCounter(xc)
			cv.So(err, cv.ShouldEqual, nil)

			act, err := msg.Marshal()
			cv.So(err, cv.ShouldEqual, nil)

			fmt.Printf("          actual:\n")
			ShowBytes(act, 10)
			//fmt.Printf("act decoded by capnp: '%s'\n", string(CapnpDecode(act, "Bag")))
			save(act, "myact")

			fmt.Printf("expected:\n")
			ShowBytes(exp, 10)
			//fmt.Printf("exp decoded by capnp: '%s'\n", string(CapnpDecode(exp, "Bag")))
			save(exp, "myexp")

			cv.So(act, cv.ShouldResemble, exp)
		})
	})
}

func ShowSeg(msg string, seg *capnp.Segment) []byte {
	b, err := seg.Message().Marshal()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s\n", msg)
	ShowBytes(b, 10)
	return b
}

func TestZserverWithOneEmptyJob(t *testing.T) {

	exp := CapnpEncode(`(waitingjobs = [()])`, "Zserver")

	cv.Convey("Given an Zserver with one empty job", t, func() {
		cv.Convey("then the go-capnproto serialization should match the capnp c++ serialization", func() {

			msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)
			_, scratch, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)
			server, err := air.NewRootZserver(seg)
			cv.So(err, cv.ShouldEqual, nil)

			joblist, err := air.NewZjob_List(seg, 1)
			cv.So(err, cv.ShouldEqual, nil)

			ShowSeg("          pre NewZjob, segment seg is:", seg)

			zjob, err := air.NewZjob(scratch)
			cv.So(err, cv.ShouldEqual, nil)
			err = joblist.Set(0, zjob)
			cv.So(err, cv.ShouldEqual, nil)

			ShowSeg("          pre SetWaitingjobs, segment seg is:", seg)

			fmt.Printf("Then we do the SetWaitingjobs:\n")
			server.SetWaitingjobs(joblist)

			// save
			act, err := msg.Marshal()
			cv.So(err, cv.ShouldEqual, nil)
			save(act, "my.act.zserver")

			// show
			ShowSeg("          actual:\n", seg)

			fmt.Printf("act decoded by capnp: '%s'\n", string(CapnpDecode(act, "Zserver")))

			fmt.Printf("expected:\n")
			ShowBytes(exp, 10)
			fmt.Printf("exp decoded by capnp: '%s'\n", string(CapnpDecode(exp, "Zserver")))

			cv.So(act, cv.ShouldResemble, exp)
		})
	})
}

func TestDefaultStructField(t *testing.T) {
	cv.Convey("Given a new root StackingRoot", t, func() {
		cv.Convey("then the aWithDefault field should have a default", func() {
			_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)
			root, err := air.NewRootStackingRoot(seg)
			cv.So(err, cv.ShouldEqual, nil)

			a, err := root.AWithDefault()
			cv.So(err, cv.ShouldEqual, nil)
			cv.So(a.Num(), cv.ShouldEqual, 42)
		})
	})
}

func TestDataTextCopyOptimization(t *testing.T) {
	cv.Convey("Given a text list from a different segment", t, func() {
		cv.Convey("Adding it to a different segment shouldn't panic", func() {
			_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)
			_, seg2, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)
			root, err := air.NewRootNester1Capn(seg)
			cv.So(err, cv.ShouldEqual, nil)

			strsl, err := capnp.NewTextList(seg2, 256)
			cv.So(err, cv.ShouldEqual, nil)
			for i := 0; i < strsl.Len(); i++ {
				strsl.Set(i, "testess")
			}

			root.SetStrs(strsl)
		})
	})
}
