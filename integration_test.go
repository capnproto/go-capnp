package capnp_test

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"reflect"
	"strings"
	"testing"

	cv "github.com/smartystreets/goconvey/convey"
	"zombiezen.com/go/capnproto2"
	air "zombiezen.com/go/capnproto2/internal/aircraftlib"
)

type bitListTest struct {
	list []bool
	text string
}

var bitListTests = []bitListTest{
	{
		[]bool{true, false, true},
		"(boolvec = [true, false, true])\n",
	},
	{
		[]bool{false},
		"(boolvec = [false])\n",
	},
	{
		[]bool{true},
		"(boolvec = [true])\n",
	},
	{
		[]bool{false, true},
		"(boolvec = [false, true])\n",
	},
	{
		[]bool{true, true},
		"(boolvec = [true, true])\n",
	},
	{
		[]bool{false, false, true},
		"(boolvec = [false, false, true])\n",
	},
	{
		[]bool{true, false, true, false, true},
		"(boolvec = [true, false, true, false, true])\n",
	},
	{
		[]bool{
			false, false, false, false, false, false, false, false,
			false, false, false, false, false, false, false, false,
			false, false, false, false, false, false, false, false,
			false, false, false, false, false, false, false, false,
			false, false, false, false, false, false, false, false,
			false, false, false, false, false, false, false, false,
			false, false, false, false, false, false, false, false,
			false, false, false, false, false, false, false, false,
			true, true,
		},
		"(boolvec = [false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, false, true, true])\n",
	},
}

func (blt bitListTest) makeMessage() (*capnp.Message, error) {
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		return nil, err
	}
	z, err := air.NewRootZ(seg)
	if err != nil {
		return nil, err
	}
	list, err := capnp.NewBitList(seg, int32(len(blt.list)))
	if err != nil {
		return nil, err
	}
	for i := range blt.list {
		list.Set(i, blt.list[i])
	}
	if err := z.SetBoolvec(list); err != nil {
		return nil, err
	}
	return msg, nil
}

func TestBitList(t *testing.T) {
	for _, test := range bitListTests {
		msg, err := test.makeMessage()
		if err != nil {
			t.Errorf("%v: make message: %v", test.list, err)
			continue
		}

		z, err := air.ReadRootZ(msg)
		if err != nil {
			t.Errorf("%v: read root Z: %v", test.list, err)
			continue
		}
		if w := z.Which(); w != air.Z_Which_boolvec {
			t.Errorf("%v: root.Which() = %v; want boolvec", test.list, w)
			continue
		}
		list, err := z.Boolvec()
		if err != nil {
			t.Errorf("%v: read Z.boolvec: %v", test.list, err)
			continue
		}
		if n := list.Len(); n != len(test.list) {
			t.Errorf("%v: len(Z.boolvec) = %d; want %d", test.list, n, len(test.list))
			continue
		}
		for i := range test.list {
			if li := list.At(i); li != test.list[i] {
				t.Errorf("%v: Z.boolvec[%d] = %t; want %t", test.list, i, li, test.list[i])
			}
		}
	}
}

func TestBitList_Decode(t *testing.T) {
	// TODO(light): skip if tool not present
	for _, test := range bitListTests {
		msg, err := test.makeMessage()
		if err != nil {
			t.Errorf("%v: make message: %v", test.list, err)
			continue
		}
		seg, _ := msg.Segment(0)
		text := CapnpDecodeSegment(seg, "", schemaPath, "Z")
		if want := strings.TrimSpace(test.text); text != want {
			t.Errorf("%v: capnp decode = %q; want %q", test.list, text, want)
		}
	}
}

func TestCreationOfZDate(t *testing.T) {
	const n = 1
	packed := false
	seg, _ := zdateFilledSegment(n, packed)
	text := CapnpDecodeSegment(seg, "", schemaPath, "Z")

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
	text := CapnpDecodeSegment(seg, "", schemaPath, "Z")

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
	text := CapnpDecodeSegment(seg, "", schemaPath, "Z")

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
	text := CapnpDecodeSegment(seg, "", schemaPath, "Z")

	expectedText := `(zdata = (data = "\x00\x01\x02\x03\x04\x05\x06\a\b\t\n\v\f\r\x0e\x0f\x10\x11\x12\x13"))`

	cv.Convey("Given a go-capnproto created Zdata DATA element with n=20", t, func() {
		cv.Convey("When we decode it with capnp", func() {
			cv.Convey(fmt.Sprintf("Then we should get the expected text '%s'", expectedText), func() {
				cv.So(text, cv.ShouldEqual, expectedText)
			})
			cv.Convey("And our data should contain Z_ZDATA with contents 0,1,2,...,n", func() {
				z, err := air.ReadRootZ(seg.Message())
				cv.So(err, cv.ShouldEqual, nil)
				cv.So(z.Which(), cv.ShouldEqual, air.Z_Which_zdata)

				zdata, err := z.Zdata()
				cv.So(err, cv.ShouldEqual, nil)
				data, err := zdata.Data()
				cv.So(err, cv.ShouldEqual, nil)
				cv.So(len(data), cv.ShouldEqual, n)
				for i := range data {
					cv.So(data[i], cv.ShouldEqual, i)
				}

			})
		})
	})

}

func TestInterfaceSet(t *testing.T) {
	cl := air.Echo{Client: capnp.ErrorClient(errors.New("foo"))}
	_, s, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatal(err)
	}
	base, err := air.NewRootEchoBase(s)
	if err != nil {
		t.Fatal(err)
	}

	base.SetEcho(cl)

	if base.Echo() != cl {
		t.Errorf("base.Echo() = %#v; want %#v", base.Echo(), cl)
	}
}

func TestInterfaceSetNull(t *testing.T) {
	cl := air.Echo{Client: capnp.ErrorClient(errors.New("foo"))}
	msg, s, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatal(err)
	}
	base, err := air.NewRootEchoBase(s)
	if err != nil {
		t.Fatal(err)
	}
	base.SetEcho(cl)

	base.SetEcho(air.Echo{})

	if e := base.Echo().Client; e != nil {
		t.Errorf("base.Echo() = %#v; want nil", e)
	}
	if len(msg.CapTable) != 1 {
		t.Errorf("msg.CapTable = %#v; want len = 1", msg.CapTable)
	}
}

func TestInterfaceCopyToOtherMessage(t *testing.T) {
	cl := air.Echo{Client: capnp.ErrorClient(errors.New("foo"))}
	_, s1, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatal(err)
	}
	base1, err := air.NewRootEchoBase(s1)
	if err != nil {
		t.Fatal(err)
	}
	if err := base1.SetEcho(cl); err != nil {
		t.Fatal(err)
	}

	_, s2, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatal(err)
	}
	hoth2, err := air.NewRootHoth(s2)
	if err != nil {
		t.Fatal(err)
	}
	if err := hoth2.SetBase(base1); err != nil {
		t.Fatal(err)
	}

	if base, err := hoth2.Base(); err != nil {
		t.Errorf("hoth2.Base() error: %v", err)
	} else if base.Echo() != cl {
		t.Errorf("hoth2.Base().Echo() = %#v; want %#v", base.Echo(), cl)
	}
	tab2 := s2.Message().CapTable
	if len(tab2) == 1 {
		if tab2[0] != cl.Client {
			t.Errorf("s2.Message().CapTable[0] = %#v; want %#v", tab2[0], cl.Client)
		}
	} else {
		t.Errorf("len(s2.Message().CapTable) = %d; want 1", len(tab2))
	}
}

func TestInterfaceCopyToOtherMessageWithCaps(t *testing.T) {
	cl := air.Echo{Client: capnp.ErrorClient(errors.New("foo"))}
	_, s1, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatal(err)
	}
	base1, err := air.NewRootEchoBase(s1)
	if err != nil {
		t.Fatal(err)
	}
	if err := base1.SetEcho(cl); err != nil {
		t.Fatal(err)
	}

	_, s2, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatal(err)
	}
	s2.Message().AddCap(nil)
	hoth2, err := air.NewRootHoth(s2)
	if err != nil {
		t.Fatal(err)
	}
	if err := hoth2.SetBase(base1); err != nil {
		t.Fatal(err)
	}

	if base, err := hoth2.Base(); err != nil {
		t.Errorf("hoth2.Base() error: %v", err)
	} else if base.Echo() != cl {
		t.Errorf("hoth2.Base().Echo() = %#v; want %#v", base.Echo(), cl)
	}
	tab2 := s2.Message().CapTable
	if len(tab2) != 2 {
		t.Errorf("len(s2.Message().CapTable) = %d; want 2", len(tab2))
	}
}

// demonstrate and test serialization to List(List(Struct(List))), nested lists.

// start with smaller Struct(List)
func Test001StructList(t *testing.T) {

	cv.Convey("Given type Nester1 struct { Strs []string } in go, where Nester1 is a struct, and a mirror/parallel capnp struct air.Nester1Capn { strs @0: List(Text); } defined in the aircraftlib schema", t, func() {
		cv.Convey("When we Save() Nester to capn and then Load() it back, the data should match, so that we have working Struct(List) serialization and deserializatoin in go-capnproto", func() {

			// Does Nester1 alone serialization and deser okay?
			rw := Nester1{Strs: []string{"xenophilia", "watchowski"}}

			var o bytes.Buffer
			rw.Save(&o)

			msg, err := capnp.Unmarshal(o.Bytes())
			cv.So(err, cv.ShouldEqual, nil)
			seg, err := msg.Segment(0)
			cv.So(err, cv.ShouldEqual, nil)

			text := CapnpDecodeSegment(seg, "", schemaPath, "Nester1Capn")
			if false {
				fmt.Printf("text = '%s'\n", text)
			}
			rw2 := &Nester1{}
			rw2.Load(&o)

			//fmt.Printf("rw = '%#v'\n", rw)
			//fmt.Printf("rw2 = '%#v'\n", rw2)

			same := reflect.DeepEqual(&rw, rw2)
			cv.So(same, cv.ShouldEqual, true)
		})
	})
}

func Test002ListListStructList(t *testing.T) {

	cv.Convey("Given type RWTest struct { NestMatrix [][]Nester1; } in go, where Nester1 is a struct, and a mirror/parallel capnp struct air.RWTestCapn { nestMatrix @0: List(List(Nester1Capn)); } defined in the aircraftlib schema", t, func() {
		cv.Convey("When we Save() RWTest to capn and then Load() it back, the data should match, so that we have working List(List(Struct)) serialization and deserializatoin in go-capnproto", func() {

			// full RWTest
			rw := RWTest{
				NestMatrix: [][]Nester1{
					[]Nester1{
						Nester1{Strs: []string{"z", "w"}},
						Nester1{Strs: []string{"q", "r"}},
					},
					[]Nester1{
						Nester1{Strs: []string{"zebra", "wally"}},
						Nester1{Strs: []string{"qubert", "rocks"}},
					},
				},
			}

			var o bytes.Buffer
			rw.Save(&o)

			msg, err := capnp.Unmarshal(o.Bytes())
			cv.So(err, cv.ShouldEqual, nil)
			seg, err := msg.Segment(0)
			cv.So(err, cv.ShouldEqual, nil)

			text := CapnpDecodeSegment(seg, "", schemaPath, "RWTestCapn")

			if false {
				fmt.Printf("text = '%s'\n", text)
			}

			rw2 := &RWTest{}
			rw2.Load(&o)

			//fmt.Printf("rw = '%#v'\n", rw)
			//fmt.Printf("rw2 = '%#v'\n", rw2)

			same := reflect.DeepEqual(&rw, rw2)
			cv.So(same, cv.ShouldEqual, true)
		})
	})
}

type Nester1 struct {
	Strs []string
}

type RWTest struct {
	NestMatrix [][]Nester1
}

func (s *Nester1) Save(w io.Writer) {
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		panic(err)
	}
	msg.SetRoot(Nester1GoToCapn(seg, s))
	data, err := msg.Marshal()
	if err != nil {
		panic(err)
	}
	w.Write(data)
}

func Nester1GoToCapn(seg *capnp.Segment, src *Nester1) air.Nester1Capn {
	//fmt.Printf("\n\n   Nester1GoToCapn sees seg = '%#v'\n", seg)
	dest, _ := air.NewNester1Capn(seg)

	mylist1, _ := capnp.NewTextList(seg, int32(len(src.Strs)))
	for i := range src.Strs {
		mylist1.Set(i, string(src.Strs[i]))
	}
	dest.SetStrs(mylist1)

	//fmt.Printf("after Nester1GoToCapn setting\n")
	return dest
}

func Nester1CapnToGo(src air.Nester1Capn, dest *Nester1) *Nester1 {
	if dest == nil {
		dest = &Nester1{}
	}
	srcStrs, _ := src.Strs()
	dest.Strs = make([]string, srcStrs.Len())
	for i := range dest.Strs {
		dest.Strs[i], _ = srcStrs.At(i)
	}

	return dest
}

func (s *Nester1) Load(r io.Reader) {
	capMsg, err := capnp.NewDecoder(r).Decode()
	if err != nil {
		panic(fmt.Errorf("capnp.ReadFromStream error: %s", err))
	}
	z, _ := air.ReadRootNester1Capn(capMsg)
	Nester1CapnToGo(z, s)
}

func (s *RWTest) Save(w io.Writer) {
	msg, seg, _ := capnp.NewMessage(capnp.SingleSegment(nil))
	msg.SetRoot(RWTestGoToCapn(seg, s))
	data, _ := msg.Marshal()
	w.Write(data)
}

func (s *RWTest) Load(r io.Reader) {
	capMsg, err := capnp.NewDecoder(r).Decode()
	if err != nil {
		panic(fmt.Errorf("capnp.ReadFromStream error: %s", err))
	}
	z, _ := air.ReadRootRWTestCapn(capMsg)
	RWTestCapnToGo(z, s)
}

func RWTestCapnToGo(src air.RWTestCapn, dest *RWTest) *RWTest {
	if dest == nil {
		dest = &RWTest{}
	}
	var n int
	srcMatrix, _ := src.NestMatrix()
	// NestMatrix
	n = srcMatrix.Len()
	dest.NestMatrix = make([][]Nester1, n)
	for i := 0; i < n; i++ {
		sm, _ := srcMatrix.At(i)
		dest.NestMatrix[i] = Nester1CapnListToSliceNester1(air.Nester1Capn_List{List: capnp.ToList(sm)})
	}

	return dest
}

func RWTestGoToCapn(seg *capnp.Segment, src *RWTest) air.RWTestCapn {
	dest, err := air.NewRWTestCapn(seg)
	if err != nil {
		panic(err)
	}

	// NestMatrix -> Nester1Capn (go slice to capn list)
	if len(src.NestMatrix) > 0 {
		plist, err := capnp.NewPointerList(seg, int32(len(src.NestMatrix)))
		if err != nil {
			panic(err)
		}
		for i, ele := range src.NestMatrix {
			err := plist.Set(i, SliceNester1ToNester1CapnList(seg, ele))
			if err != nil {
				panic(err)
			}
		}
		dest.SetNestMatrix(plist)
	}

	return dest
}

func Nester1CapnListToSliceNester1(p air.Nester1Capn_List) []Nester1 {
	v := make([]Nester1, p.Len())
	for i := range v {
		Nester1CapnToGo(p.At(i), &v[i])
	}
	return v
}

func SliceNester1ToNester1CapnList(seg *capnp.Segment, m []Nester1) air.Nester1Capn_List {
	lst, err := air.NewNester1Capn_List(seg, int32(len(m)))
	if err != nil {
		panic(err)
	}
	for i := range m {
		err := lst.Set(i, Nester1GoToCapn(seg, &m[i]))
		if err != nil {
			panic(err)
		}
	}
	return lst
}

func SliceStringToTextList(seg *capnp.Segment, m []string) capnp.TextList {
	lst, err := capnp.NewTextList(seg, int32(len(m)))
	if err != nil {
		panic(err)
	}
	for i := range m {
		lst.Set(i, string(m[i]))
	}
	return lst
}

func TextListToSliceString(p capnp.TextList) []string {
	v := make([]string, p.Len())
	for i := range v {
		s, err := p.At(i)
		if err != nil {
			panic(err)
		}
		v[i] = s
	}
	return v
}

func TestDataVersioningAvoidsUnnecessaryTruncation(t *testing.T) {

	expFull := encodeTestMessage(t, "VerTwoDataTwoPtr", "(val = 9, duo = 8, ptr1 = (val = 77), ptr2 = (val = 55))", []byte{
		0, 0, 0, 0, 7, 0, 0, 0,
		0, 0, 0, 0, 2, 0, 2, 0,
		9, 0, 0, 0, 0, 0, 0, 0,
		8, 0, 0, 0, 0, 0, 0, 0,
		4, 0, 0, 0, 1, 0, 0, 0,
		4, 0, 0, 0, 1, 0, 0, 0,
		77, 0, 0, 0, 0, 0, 0, 0,
		55, 0, 0, 0, 0, 0, 0, 0,
	})

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

func TestTextAndListTextContaintingEmptyStruct(t *testing.T) {

	emptyZjobBytes := encodeTestMessage(t, "Zjob", "()", []byte{
		0, 0, 0, 0, 3, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 2, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	})

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

	zjobBytes := encodeTestMessage(t, "Zjob", `(cmd = "abc")`, []byte{
		0, 0, 0, 0, 4, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 2, 0,
		5, 0, 0, 0, 34, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		97, 98, 99, 0, 0, 0, 0, 0,
	})

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

	zjobBytes := encodeTestMessage(t, "Zjob", `(args = ["xyz"])`, []byte{
		0, 0, 0, 0, 5, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 2, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		1, 0, 0, 0, 14, 0, 0, 0,
		1, 0, 0, 0, 34, 0, 0, 0,
		120, 121, 122, 0, 0, 0, 0, 0,
	})

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

	zjobBytes := encodeTestMessage(t, "Zjob", `(cmd = "abc", args = ["xyz"])`, []byte{
		0, 0, 0, 0, 6, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 2, 0,
		5, 0, 0, 0, 34, 0, 0, 0,
		5, 0, 0, 0, 14, 0, 0, 0,
		97, 98, 99, 0, 0, 0, 0, 0,
		1, 0, 0, 0, 34, 0, 0, 0,
		120, 121, 122, 0, 0, 0, 0, 0,
	})

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

	exp := encodeTestMessage(t, "Zserver", `(waitingjobs = [(cmd = "abc", args = ["xyz"])])`, []byte{
		0, 0, 0, 0, 8, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 1, 0,
		1, 0, 0, 0, 23, 0, 0, 0,
		4, 0, 0, 0, 0, 0, 2, 0,
		5, 0, 0, 0, 34, 0, 0, 0,
		5, 0, 0, 0, 14, 0, 0, 0,
		97, 98, 99, 0, 0, 0, 0, 0,
		1, 0, 0, 0, 34, 0, 0, 0,
		120, 121, 122, 0, 0, 0, 0, 0,
	})

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

			fmt.Printf("expected:\n")
			ShowBytes(exp, 10)
			fmt.Printf("exp decoded by capnp: '%s'\n", string(CapnpDecode(exp, "Zserver")))

			cv.So(act, cv.ShouldResemble, exp)
		})
	})
}

func TestZserverWithAccessors(t *testing.T) {

	exp := encodeTestMessage(t, "Zserver", `(waitingjobs = [(cmd = "abc"), (cmd = "xyz")])`, []byte{
		0, 0, 0, 0, 9, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 1, 0,
		1, 0, 0, 0, 39, 0, 0, 0,
		8, 0, 0, 0, 0, 0, 2, 0,
		13, 0, 0, 0, 34, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		9, 0, 0, 0, 34, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		97, 98, 99, 0, 0, 0, 0, 0,
		120, 121, 122, 0, 0, 0, 0, 0,
	})

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

			fmt.Printf("expected:\n")
			ShowBytes(exp, 10)
			fmt.Printf("exp decoded by capnp: '%s'\n", string(CapnpDecode(exp, "Zserver")))

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

	exp := encodeTestMessage(t, "Bag", `(counter = (size = 9))`, []byte{
		0, 0, 0, 0, 5, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 1, 0,
		0, 0, 0, 0, 1, 0, 2, 0,
		9, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	})

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

			fmt.Printf("expected:\n")
			ShowBytes(exp, 10)
			fmt.Printf("exp decoded by capnp: '%s'\n", string(CapnpDecode(exp, "Bag")))

			cv.So(act, cv.ShouldResemble, exp)
		})
	})
}

func TestObjectWithTextBetweenSegments(t *testing.T) {

	exp := encodeTestMessage(t, "Bag", `(counter = (size = 9, words = "hello"))`, []byte{
		0, 0, 0, 0, 6, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 1, 0,
		0, 0, 0, 0, 1, 0, 2, 0,
		9, 0, 0, 0, 0, 0, 0, 0,
		5, 0, 0, 0, 50, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		104, 101, 108, 108, 111, 0, 0, 0,
	})

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

			fmt.Printf("expected:\n")
			ShowBytes(exp, 10)
			fmt.Printf("exp decoded by capnp: '%s'\n", string(CapnpDecode(exp, "Bag")))

			cv.So(act, cv.ShouldResemble, exp)
		})
	})
}

func TestObjectWithListOfTextBetweenSegments(t *testing.T) {

	exp := encodeTestMessage(t, "Bag", `(counter = (size = 9, wordlist = ["hello","bye"]))`, []byte{
		0, 0, 0, 0, 9, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 1, 0,
		0, 0, 0, 0, 1, 0, 2, 0,
		9, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		1, 0, 0, 0, 22, 0, 0, 0,
		5, 0, 0, 0, 50, 0, 0, 0,
		5, 0, 0, 0, 34, 0, 0, 0,
		104, 101, 108, 108, 111, 0, 0, 0,
		98, 121, 101, 0, 0, 0, 0, 0,
	})

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

	exp := encodeTestMessage(t, "Bag", `(counter = (size = 9, words = "abc", wordlist = ["hello","byenow"]))`, []byte{
		0, 0, 0, 0, 10, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 1, 0,
		0, 0, 0, 0, 1, 0, 2, 0,
		9, 0, 0, 0, 0, 0, 0, 0,
		5, 0, 0, 0, 34, 0, 0, 0,
		5, 0, 0, 0, 22, 0, 0, 0,
		97, 98, 99, 0, 0, 0, 0, 0,
		5, 0, 0, 0, 50, 0, 0, 0,
		5, 0, 0, 0, 58, 0, 0, 0,
		104, 101, 108, 108, 111, 0, 0, 0,
		98, 121, 101, 110, 111, 119, 0, 0,
	})

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

			fmt.Printf("expected:\n")
			ShowBytes(exp, 10)
			//fmt.Printf("exp decoded by capnp: '%s'\n", string(CapnpDecode(exp, "Bag")))

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

	exp := encodeTestMessage(t, "Zserver", `(waitingjobs = [()])`, []byte{
		0, 0, 0, 0, 5, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 1, 0,
		1, 0, 0, 0, 23, 0, 0, 0,
		4, 0, 0, 0, 0, 0, 2, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	})

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

			act, err := msg.Marshal()
			cv.So(err, cv.ShouldEqual, nil)

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

// highlight how much faster text movement between segments
// is when special casing Text and Data
//
// run this test with capnp.go:1334-1341 commented in/out to compare.
//
func BenchmarkTextMovementBetweenSegments(b *testing.B) {

	buf := make([]byte, 1<<21)
	buf2 := make([]byte, 1<<21)

	text := make([]byte, 1<<20)
	for i := range text {
		text[i] = byte(65 + rand.Int()%26)
	}
	//stext := string(text)
	//fmt.Printf("text = %#v\n", stext)

	astr := make([]string, 1000)
	for i := range astr {
		astr[i] = string(text[i*1000 : (i+1)*1000])
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, seg, _ := capnp.NewMessage(capnp.SingleSegment(buf[:0]))
		_, scratch, _ := capnp.NewMessage(capnp.SingleSegment(buf2[:0]))

		ht, _ := air.NewRootHoldsText(seg)
		tl, _ := capnp.NewTextList(scratch, 1000)

		for j := 0; j < 1000; j++ {
			tl.Set(j, astr[j])
		}

		ht.SetLst(tl)

	}
}

func TestV0ListofEmptyShouldMatchCapnp(t *testing.T) {

	exp := encodeTestMessage(t, "HoldsVerEmptyList", "(mylist = [(),()])", []byte{
		0, 0, 0, 0, 3, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 1, 0,
		1, 0, 0, 0, 7, 0, 0, 0,
		8, 0, 0, 0, 0, 0, 0, 0,
	})

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

			act, err := msg.Marshal()
			cv.So(err, cv.ShouldEqual, nil)

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

	expEmpty := encodeTestMessage(t, "HoldsVerEmptyList", "(mylist = [(),()])", []byte{
		0, 0, 0, 0, 3, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 1, 0,
		1, 0, 0, 0, 7, 0, 0, 0,
		8, 0, 0, 0, 0, 0, 0, 0,
	})

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

	expEmpty := encodeTestMessage(t, "HoldsVerEmptyList", "(mylist = [(),()])", []byte{
		0, 0, 0, 0, 3, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 1, 0,
		1, 0, 0, 0, 7, 0, 0, 0,
		8, 0, 0, 0, 0, 0, 0, 0,
	})

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

func TestVoidUnionSetters(t *testing.T) {
	want := encodeTestMessage(t, "VoidUnion", "(b = void)", []byte{
		0, 0, 0, 0, 2, 0, 0, 0,
		0, 0, 0, 0, 1, 0, 0, 0,
		1, 0, 0, 0, 0, 0, 0, 0,
	})

	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatal(err)
	}
	voidUnion, err := air.NewRootVoidUnion(seg)
	if err != nil {
		t.Fatal(err)
	}
	voidUnion.SetB()

	act, err := msg.Marshal()
	if err != nil {
		t.Fatal("msg.Marshal():", err)
	}
	if !bytes.Equal(act, want) {
		t.Errorf("msg.Marshal() =\n%s\n; want:\n%s", hex.Dump(act), hex.Dump(want))
	}
}
