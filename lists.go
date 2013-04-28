package capn

import (
	"math"
	"reflect"
)

type Bitset struct {
	Data []byte
	Size int
}

func MakeBitset(size int) Bitset { return Bitset{make([]byte, (size+7)/8), size} }
func (b Bitset) Test(i int) bool { return b.Data[i/8]&(1<<uint(i%8)) != 0 }
func (b *Bitset) Reset(i int)    { b.Data[i/8] &^= 1 << uint(i%8) }
func (b *Bitset) Set(i int)      { b.Data[i/8] |= 1 << uint(i%8) }

func (p Pointer) WriteStructF32(off int, v, def float32) error {
	return p.WriteStruct32(off, math.Float32bits(v)^math.Float32bits(def))
}
func (p Pointer) WriteStructF64(off int, v, def float64) error {
	return p.WriteStruct64(off, math.Float64bits(v)^math.Float64bits(def))
}

func (p Pointer) ReadStructF32(off int, def float32) float32 {
	return math.Float32frombits(p.ReadStruct32(off) ^ math.Float32bits(def))
}
func (p Pointer) ReadStructF64(off int, def float64) float64 {
	return math.Float64frombits(p.ReadStruct64(off) ^ math.Float64bits(def))
}

func (p Pointer) ReadStruct(off int, def Pointer) Pointer {
	if m := p.ReadPtr(off); m.typ == Struct {
		return m
	}
	return def
}

func (p Pointer) WriteStruct(off int, v, def Pointer) error {
	if v == def {
		return p.WritePtr(off, Pointer{})
	} else {
		return p.WritePtr(off, v)
	}
}

func (p Pointer) ReadData(off int, def []byte) []byte {
	if m := p.ReadPtr(off); m.typ == List && m.dataBits == 8 {
		return m.Data()
	}
	return def
}

func (p Pointer) ReadPointerList(off int, def []Pointer) []Pointer {
	if m := p.ReadPtr(off); m.typ == List {
		ret := make([]Pointer, m.Size())
		for i := range ret {
			ret[i] = m.ReadPtr(i)
		}
		return ret
	}
	return def
}

func (p Pointer) WriteString(off int, v, def string) (err error) {
	var m Pointer
	if v != def {
		if m, err = p.Segment.NewList(8, 0, len(v)+1); err != nil {
			return err
		}
		copy(m.Data(), []byte(v))
	}
	return p.WritePtr(off, m)
}
func (p Pointer) ReadString(off int, def string) string {
	if m := p.ReadPtr(off); m.typ == List && m.dataBits == 8 {
		d := m.Data()
		if len(d) > 0 && d[len(d)-1] == 0 {
			return string(d)
		}
	}
	return def
}

func (p Pointer) WriteVoidList(off int, v, def []struct{}) (err error) {
	var m Pointer
	if len(v) != len(def) {
		if m, err = p.Segment.NewList(0, 0, len(v)); err != nil {
			return err
		}
	}
	return p.WritePtr(off, m)
}

func SliceEqual(a, b interface{}) bool {
	if (a == nil) != (b == nil) {
		return false
	}
	u := reflect.ValueOf(a)
	v := reflect.ValueOf(b)
	if u.Len() != v.Len() {
		return false
	}
	for i := 0; i < u.Len(); i++ {
		if u.Index(i) != v.Index(i) {
			return false
		}
	}
	return true
}

func (p Pointer) WritePointerList(off int, v, def []Pointer) (err error) {
	var m Pointer
	if !SliceEqual(v, def) {
		if m, err = p.Segment.NewPointerList(len(v)); err != nil {
			return err
		}
		for i, u := range v {
			if err := m.WritePtr(i, u); err != nil {
				return err
			}
		}
	}
	return p.WritePtr(off, m)
}

func (p Pointer) WriteBitset(off int, v, def Bitset) (err error) {
	var m Pointer
	if v.Size != def.Size || !SliceEqual(v.Data, def.Data) {
		if m, err = p.Segment.NewList(1, 0, v.Size); err != nil {
			return err
		}
		copy(m.Data(), v.Data)
	}
	return p.WritePtr(off, m)
}

func (p Pointer) WriteU8List(off int, v, def []uint8) (err error) {
	var m Pointer
	if !SliceEqual(v, def) {
		if m, err = p.Segment.NewList(8, 0, len(v)); err != nil {
			return err
		}
		copy(m.Data(), v)
	}
	return p.WritePtr(off, m)
}
func (p Pointer) WriteU16List(off int, v, def []uint16) (err error) {
	var m Pointer
	if !SliceEqual(v, def) {
		if m, err = p.Segment.NewList(16, 0, len(v)); err != nil {
			return err
		}
		d := m.Data()
		for i, u := range v {
			putLittle16(d[16/8*i:], u)
		}
	}
	return p.WritePtr(off, m)
}
func (p Pointer) WriteU32List(off int, v, def []uint32) (err error) {
	var m Pointer
	if !SliceEqual(v, def) {
		if m, err = p.Segment.NewList(32, 0, len(v)); err != nil {
			return err
		}
		d := m.Data()
		for i, u := range v {
			putLittle32(d[32/8*i:], u)
		}
	}
	return p.WritePtr(off, m)
}
func (p Pointer) WriteU64List(off int, v, def []uint64) (err error) {
	var m Pointer
	if !SliceEqual(v, def) {
		if m, err = p.Segment.NewList(64, 0, len(v)); err != nil {
			return err
		}
		d := m.Data()
		for i, u := range v {
			putLittle64(d[64/8*i:], u)
		}
	}
	return p.WritePtr(off, m)
}
func (p Pointer) WriteI8List(off int, v, def []int8) (err error) {
	var m Pointer
	if !SliceEqual(v, def) {
		if m, err = p.Segment.NewList(8, 0, len(v)); err != nil {
			return err
		}
		d := m.Data()
		for i, u := range v {
			d[i] = uint8(u)
		}
	}
	return p.WritePtr(off, m)
}
func (p Pointer) WriteI16List(off int, v, def []int16) (err error) {
	var m Pointer
	if !SliceEqual(v, def) {
		if m, err = p.Segment.NewList(16, 0, len(v)); err != nil {
			return err
		}
		d := m.Data()
		for i, u := range v {
			putLittle16(d[16/8*i:], uint16(u))
		}
	}
	return p.WritePtr(off, m)
}
func (p Pointer) WriteI32List(off int, v, def []int32) (err error) {
	var m Pointer
	if !SliceEqual(v, def) {
		if m, err = p.Segment.NewList(32, 0, len(v)); err != nil {
			return err
		}
		d := m.Data()
		for i, u := range v {
			putLittle32(d[32/8*i:], uint32(u))
		}
	}
	return p.WritePtr(off, m)
}
func (p Pointer) WriteI64List(off int, v, def []int64) (err error) {
	var m Pointer
	if !SliceEqual(v, def) {
		if m, err = p.Segment.NewList(64, 0, len(v)); err != nil {
			return err
		}
		d := m.Data()
		for i, u := range v {
			putLittle64(d[64/8*i:], uint64(u))
		}
	}
	return p.WritePtr(off, m)
}
func (p Pointer) WriteF32List(off int, v, def []float32) (err error) {
	var m Pointer
	if !SliceEqual(v, def) {
		if m, err = p.Segment.NewList(32, 0, len(v)); err != nil {
			return err
		}
		d := m.Data()
		for i, u := range v {
			putLittle32(d[32/8*i:], math.Float32bits(u))
		}
	}
	return p.WritePtr(off, m)
}
func (p Pointer) WriteF64List(off int, v, def []float64) (err error) {
	var m Pointer
	if !SliceEqual(v, def) {
		if m, err = p.Segment.NewList(64, 0, len(v)); err != nil {
			return err
		}
		d := m.Data()
		for i, u := range v {
			putLittle64(d[64/8*i:], math.Float64bits(u))
		}
	}
	return p.WritePtr(off, m)
}

func (p Pointer) ReadVoidList(off int, def []struct{}) []struct{} {
	if m := p.ReadPtr(off); m.typ == List {
		return make([]struct{}, m.Size())
	}
	return def
}

func (p Pointer) ReadBitset(off int, def Bitset) Bitset {
	if m := p.ReadPtr(off); m.typ == List && m.dataBits >= 1 {
		r := MakeBitset(m.Size())
		for i := 0; i < r.Size; i++ {
			if m.Read1(i) {
				r.Set(i)
			}
		}
		return r
	}
	return def
}

func (p Pointer) ReadU8List(off int, def []uint8) []uint8 {
	if m := p.ReadPtr(off); m.typ != Struct {
		r := make([]uint8, m.Size())
		for i := range r {
			r[i] = m.Read8(i)
		}
		return r
	}
	return def
}
func (p Pointer) ReadU16List(off int, def []uint16) []uint16 {
	if m := p.ReadPtr(off); m.typ != Struct {
		r := make([]uint16, m.Size())
		for i := range r {
			r[i] = m.Read16(i)
		}
		return r
	}
	return def
}
func (p Pointer) ReadU32List(off int, def []uint32) []uint32 {
	if m := p.ReadPtr(off); m.typ != Struct {
		r := make([]uint32, m.Size())
		for i := range r {
			r[i] = m.Read32(i)
		}
		return r
	}
	return def
}
func (p Pointer) ReadU64List(off int, def []uint64) []uint64 {
	if m := p.ReadPtr(off); m.typ != Struct {
		r := make([]uint64, m.Size())
		for i := range r {
			r[i] = m.Read64(i)
		}
		return r
	}
	return def
}
func (p Pointer) ReadI8List(off int, def []int8) []int8 {
	if m := p.ReadPtr(off); m.typ != Struct {
		r := make([]int8, m.Size())
		for i := range r {
			r[i] = int8(m.Read8(i))
		}
		return r
	}
	return def
}
func (p Pointer) ReadI16List(off int, def []int16) []int16 {
	if m := p.ReadPtr(off); m.typ != Struct {
		r := make([]int16, m.Size())
		for i := range r {
			r[i] = int16(m.Read16(i))
		}
		return r
	}
	return def
}
func (p Pointer) ReadI32List(off int, def []int32) []int32 {
	if m := p.ReadPtr(off); m.typ != Struct {
		r := make([]int32, m.Size())
		for i := range r {
			r[i] = int32(m.Read32(i))
		}
		return r
	}
	return def
}
func (p Pointer) ReadI64List(off int, def []int64) []int64 {
	if m := p.ReadPtr(off); m.typ != Struct {
		r := make([]int64, m.Size())
		for i := range r {
			r[i] = int64(m.Read64(i))
		}
		return r
	}
	return def
}
func (p Pointer) ReadF32List(off int, def []float32) []float32 {
	if m := p.ReadPtr(off); m.typ != Struct {
		r := make([]float32, m.Size())
		for i := range r {
			r[i] = math.Float32frombits(m.Read32(i))
		}
		return r
	}
	return def
}
func (p Pointer) ReadF64List(off int, def []float64) []float64 {
	if m := p.ReadPtr(off); m.typ != Struct {
		r := make([]float64, m.Size())
		for i := range r {
			r[i] = math.Float64frombits(m.Read64(i))
		}
		return r
	}
	return def
}
