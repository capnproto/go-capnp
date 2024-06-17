package capnp_test

import (
	"encoding/binary"
	"encoding/hex"
	"math"
	"strings"
	"testing"
	"time"
	"unsafe"

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
	c.SetBirthDayp(a.BirthDay.UnixNano())
	c.SetSiblingsp(int32(a.Siblings))
	c.SetSpousep(a.Spouse)
	c.SetMoneyp(a.Money)      // c.SetMoney(a.Money)
	x.fieldName.Set(a.Name)   // c.FlatSetName(a.Name) // c.SetName(a.Name)
	x.fieldPhone.Set(a.Phone) // c.FlatSetPhone(a.Phone) // c.SetPhone(a.Phone)

	return x.arena.Data(0)
}

func (x *CapNProtoSerializer) Unmarshal(d []byte, i interface{}) error {
	a := i.(*goserbench.SmallStruct)

	x.arena.ReplaceBuffer(d)

	var err error

	c := x.c
	a.Name, err = c.GetNameSuperUnsafe() //c.GetName() // c.Name()
	if err != nil {
		return err
	}
	a.BirthDay = time.Unix(0, c.GetBirthDay())
	a.Phone, err = c.GetPhoneSuperUnsafe() // c.GetPhone() // c.Phone()
	if err != nil {
		return err
	}
	a.Siblings = int(c.GetSiblings())
	a.Spouse = c.GetSpouse()
	a.Money = c.GetMoney()
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
	if err := a.SetPhone(strings.Repeat("a", goserbench.MaxSmallStructPhoneSize)); err != nil {
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

func BenchmarkSetUint(b *testing.B) {
	b.Run("baseline", func(b *testing.B) {
		buf := make([]byte, 1024)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			binary.LittleEndian.PutUint64(buf[48:], uint64(i))
		}
	})
	b.Run("unsafe", func(b *testing.B) {
		buf := make([]byte, 1024)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			*(*uint64)(unsafe.Pointer(&buf[48])) = uint64(i)
		}

	})
	b.Run("capnp", func(b *testing.B) {
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

		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			a.SetBirthDayp(int64(i))
		}
	})
}

func BenchmarkSetFloat(b *testing.B) {
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

	b.Run("baseline", func(b *testing.B) {
		buf := make([]byte, 1024)
		var f float64
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			f += float64(i)
			binary.LittleEndian.PutUint64(buf[48:], math.Float64bits(f))
		}

		_ = hex.EncodeToString(buf)
	})
	b.Run("capnp", func(b *testing.B) {
		b.ResetTimer()

		var f float64
		for i := 0; i < b.N; i++ {
			f += float64(i)
			a.SetMoneyp(f)
		}

		_ = hex.EncodeToString(arena.Segment(0).Data())
	})
	b.Run("uint", func(b *testing.B) {
		b.ResetTimer()

		var f float64
		for i := 0; i < b.N; i++ {
			f += float64(i)
			a.SetBirthDayp(int64(math.Float64bits(f)))
		}

		_ = hex.EncodeToString(arena.Segment(0).Data())
	})

}

func BenchmarkSetBool(b *testing.B) {
	buf := make([]byte, 1024)
	b.Run("readandmask", func(b *testing.B) {
		var bl bool
		for i := 0; i < b.N; i++ {
			//addr, off := i%len(buf), i%8
			addr, off := 12, 0

			bl = !bl
			aux := buf[addr]
			if bl {
				aux |= (1 << off)
			} else {
				aux &^= (1 << off)
			}
			buf[addr] = aux
		}
	})
	b.Run("directly", func(b *testing.B) {
		var bl bool
		for i := 0; i < b.N; i++ {
			addr, off := i%len(buf), i%8

			bl = !bl
			if bl {
				buf[addr] |= (1 << off)
			} else {
				buf[addr] &^= (1 << off)
			}
		}
	})
	b.Run("xor", func(b *testing.B) {
		var bl bool
		for i := 0; i < b.N; i++ {
			addr, off := i%len(buf), i%8

			bl = !bl

			aux := byte(1 << off)
			if bl {
				aux ^= 0xff
			}
			buf[addr] ^= aux
		}
	})

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

	b.Run("capnp", func(b *testing.B) {
		var bl bool
		for i := 0; i < b.N; i++ {
			bl = !bl
			a.SetSpousep(bl)
		}
	})

}

func sliceAt32(b []byte) []byte {
	return b[32:]
}

func BenchmarkGetUint64(b *testing.B) {
	buf := make([]byte, 1024)
	const target = int64(0x17010666)
	binary.LittleEndian.PutUint64(buf[32:], uint64(target))
	b.Run("baseline", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v := binary.LittleEndian.Uint64(buf[32:])
			if int64(v) != target {
				panic("boo")
			}
		}
	})

	b.Run("func", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v := binary.LittleEndian.Uint64(sliceAt32(buf))
			if int64(v) != target {
				panic("boo")
			}
		}
	})

	arena := &capnp.SimpleSingleSegmentArena{}
	var msg capnp.Message
	msg.ResetForRead(arena)
	msg.ResetReadLimit(1 << 31)

	a, err := aircraftlib.AllocateNewRootBenchmark(&msg)
	if err != nil {
		panic(err)
	}

	a.SetBirthDay(int64(target))

	b.Run("capnp", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v := a.GetBirthDay()
			if v != target {
				panic("boo")
			}
		}
	})
}

func BenchmarkGetText(b *testing.B) {
	const target = "i am the walrus!"
	buf := make([]byte, 1024)
	copy(buf[32:], target)
	b.Run("baseline", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tb := buf[32 : 32+16]
			v := *(*string)(unsafe.Pointer(&tb))
			if v != target {
				b.Fatalf("v %q", v)
			}
		}
	})

	arena := &capnp.SimpleSingleSegmentArena{}
	var msg capnp.Message
	msg.ResetForRead(arena)
	msg.ResetReadLimit(1 << 31)

	a, err := aircraftlib.AllocateNewRootBenchmark(&msg)
	if err != nil {
		panic(err)
	}

	err = a.SetName(target)
	if err != nil {
		b.Fatal(err)
	}

	b.Run("capnp", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v, err := a.GetName()
			if err != nil {
				b.Fatal(err)
			}
			if v != target {
				b.Fatalf("v %q", v)
			}
		}
	})

	b.Run("superunsafe", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v, err := a.GetNameSuperUnsafe()
			if err != nil {
				b.Fatal(err)
			}
			if v != target {
				b.Fatalf("v %q", v)
			}
		}
	})

}
