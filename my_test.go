package capnp_test

import (
	"strings"
	"testing"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/internal/aircraftlib"

	"github.com/alecthomas/go_serialization_benchmarks/goserbench"
)

func BenchmarkSetText05(b *testing.B) {
	var msg capnp.Message
	arena := &capnp.SimpleSingleSegmentArena{}

	msg.ResetForRead(arena)
	seg, err := msg.Segment(0)
	if err != nil {
		b.Fatal(err)
	}
	tx, err := aircraftlib.NewBenchmarkA(seg)
	if err != nil {
		b.Fatal(err)
	}

	err = tx.SetName("my own descr")
	if err != nil {
		b.Fatal(err)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		msg.ResetForRead(arena)
		seg, _ := msg.Segment(0)

		tx, err := aircraftlib.NewBenchmarkA(seg)
		if err != nil {
			b.Fatal(err)
		}

		err = tx.SetName("my own descr")
		if err != nil {
			b.Fatal(err)
		}
	}

	// b.Log(arena.String())

}

func BenchmarkSetInt(b *testing.B) {
	var msg capnp.Message
	arena := &capnp.SimpleSingleSegmentArena{}

	msg.ResetForRead(arena)
	seg, err := msg.Segment(0)
	if err != nil {
		b.Fatal(err)
	}
	tx, err := aircraftlib.NewBenchmarkA(seg)
	if err != nil {
		b.Fatal(err)
	}

	tx.SetBirthDay(0x20010101)

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tx.SetBirthDay(0x20010101 + int64(i))
	}

	// b.Log(arena.String())

}
func BenchmarkSetTextBaselineCopyClear(b *testing.B) {
	var pt *[]byte
	buf := make([]byte, 1024)
	pt = &buf

	b.ResetTimer()
	src := "my own descr"
	off := 48
	for i := 0; i < b.N; i++ {
		s := (*pt)[off : off+len(src)+1]
		n := copy(s, []byte(src))
		for i := n; i < len(s); i++ {
			s[i] = 0
		}
	}
}

type CapNProtoSerializer struct {
	msg   capnp.Message
	arena *capnp.SimpleSingleSegmentArena
	c     *aircraftlib.BenchmarkA

	fieldName  *capnp.TextField
	fieldPhone *capnp.TextField
}

func (x *CapNProtoSerializer) Marshal(o interface{}) ([]byte, error) {
	a := o.(*goserbench.SmallStruct)
	/*
		x.msg.ResetForRead(x.arena)
		seg, err := x.msg.Segment(0)
		if err != nil {
			return nil, err
		}

		c, err := aircraftlib.AllocateNewRootBenchmark(&x.msg)
		if err != nil {
			return nil, err
		}
	*/

	c := x.c
	c.SetBirthDay(a.BirthDay.UnixNano())
	c.SetSiblingsp(int32(a.Siblings))
	c.SetSpousep(a.Spouse)
	c.SetMoneyp(a.Money)      // c.SetMoney(a.Money)
	x.fieldName.Set(a.Name)   // c.FlatSetName(a.Name) // c.SetName(a.Name)
	x.fieldPhone.Set(a.Phone) // c.FlatSetPhone(a.Phone) // c.SetPhone(a.Phone)

	return x.arena.Data(0)
}

func (x *CapNProtoSerializer) Unmarshal(d []byte, i interface{}) error {
	/*
		a := i.(*goserbench.SmallStruct)

		s, _, err := capn.ReadFromMemoryZeroCopy(d)
		if err != nil {
			return err

		}
		o := aircraftlib.ReadRootBenchmarkA(s)
		a.Name = o.Name()
		a.BirthDay = time.Unix(0, o.BirthDay())
		a.Phone = o.Phone()
		a.Siblings = int(o.Siblings())
		a.Spouse = o.Spouse()
		a.Money = o.Money()
		return nil
	*/
	return nil
}

func NewCapNProtoSerializer() goserbench.Serializer {
	arena := &capnp.SimpleSingleSegmentArena{}
	var msg capnp.Message
	msg.ResetForRead(arena)
	msg.ResetReadLimit(1 << 31)

	a, err := aircraftlib.AllocateNewRootBenchmark(&msg)
	if err != nil {
		panic(err)
	}
	if err := a.SetName(strings.Repeat("a", goserbench.MaxSmallStructNameSize)); err != nil {
		panic(err)
	}
	if err := a.SetPhone(strings.Repeat("a", goserbench.MaxSmallStructNameSize)); err != nil {
		panic(err)
	}

	fieldName, err := a.NameField()
	if err != nil {
		panic(err)
	}
	fieldPhone, err := a.PhoneField()
	if err != nil {
		panic(err)
	}

	return &CapNProtoSerializer{
		arena:      arena,
		c:          &a,
		fieldName:  &fieldName,
		fieldPhone: &fieldPhone,
	}

}

func BenchmarkGoserBench(b *testing.B) {
	b.Run("marshal", func(b *testing.B) {
		goserbench.BenchMarshalSmallStruct(b, NewCapNProtoSerializer())
	})

	b.Run("unmarshal", func(b *testing.B) {
		goserbench.BenchUnmarshalSmallStruct(b, NewCapNProtoSerializer(), false)
	})

}
