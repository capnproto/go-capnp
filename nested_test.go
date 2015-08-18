package capnp_test

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"testing"

	cv "github.com/smartystreets/goconvey/convey"
	"zombiezen.com/go/capnproto"
	air "zombiezen.com/go/capnproto/internal/aircraftlib"
)

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
