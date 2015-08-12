package capnp

import (
	"encoding/binary"
	"math"
)

// Primitive list types.
type (
	VoidList    Pointer
	BitList     Pointer
	Int8List    Pointer
	UInt8List   Pointer
	Int16List   Pointer
	UInt16List  Pointer
	Int32List   Pointer
	UInt32List  Pointer
	Float32List Pointer
	Int64List   Pointer
	UInt64List  Pointer
	Float64List Pointer
	PointerList Pointer
	TextList    Pointer
	DataList    Pointer
)

func (p VoidList) Len() int    { return int(p.length) }
func (p BitList) Len() int     { return int(p.length) }
func (p Int8List) Len() int    { return int(p.length) }
func (p UInt8List) Len() int   { return int(p.length) }
func (p Int16List) Len() int   { return int(p.length) }
func (p UInt16List) Len() int  { return int(p.length) }
func (p Int32List) Len() int   { return int(p.length) }
func (p UInt32List) Len() int  { return int(p.length) }
func (p Float32List) Len() int { return int(p.length) }
func (p Int64List) Len() int   { return int(p.length) }
func (p UInt64List) Len() int  { return int(p.length) }
func (p Float64List) Len() int { return int(p.length) }
func (p PointerList) Len() int { return int(p.length) }
func (p TextList) Len() int    { return int(p.length) }
func (p DataList) Len() int    { return int(p.length) }

func (s *Segment) NewVoidList(n int32) VoidList       { return VoidList(s.newList(0, n)) }
func (s *Segment) NewInt8List(n int32) Int8List       { return Int8List(s.newList(1, n)) }
func (s *Segment) NewUInt8List(n int32) UInt8List     { return UInt8List(s.newList(1, n)) }
func (s *Segment) NewInt16List(n int32) Int16List     { return Int16List(s.newList(2, n)) }
func (s *Segment) NewUInt16List(n int32) UInt16List   { return UInt16List(s.newList(2, n)) }
func (s *Segment) NewFloat32List(n int32) Float32List { return Float32List(s.newList(4, n)) }
func (s *Segment) NewInt32List(n int32) Int32List     { return Int32List(s.newList(4, n)) }
func (s *Segment) NewUInt32List(n int32) UInt32List   { return UInt32List(s.newList(4, n)) }
func (s *Segment) NewFloat64List(n int32) Float64List { return Float64List(s.newList(8, n)) }
func (s *Segment) NewInt64List(n int32) Int64List     { return Int64List(s.newList(8, n)) }
func (s *Segment) NewUInt64List(n int32) UInt64List   { return UInt64List(s.newList(8, n)) }
func (s *Segment) newList(datasz Size, length int32) Pointer {
	n, _ := s.create(datasz*Size(length), Pointer{
		typ:    TypeList,
		length: length,
		size:   ObjectSize{DataSize: datasz},
	})
	return n
}

func (s *Segment) NewTextList(n int32) TextList { return TextList(s.NewPointerList(n)) }
func (s *Segment) NewDataList(n int32) DataList { return DataList(s.NewPointerList(n)) }
func (s *Segment) NewPointerList(length int32) PointerList {
	n, _ := s.create(wordSize.times(length), Pointer{
		typ:    TypePointerList,
		length: length,
		size:   ObjectSize{PointerCount: 1},
	})
	return PointerList(n)
}

func (s *Segment) NewCompositeList(sz ObjectSize, length int32) PointerList {
	if !sz.isValid() {
		return PointerList{}
	}
	sz.DataSize = sz.DataSize.padToWord()
	n, _ := s.create(wordSize+sz.totalSize().times(length), Pointer{
		typ:    TypeList,
		length: length,
		size:   sz,
		flags:  isCompositeList,
	})
	// Add tag word
	putRawPointer(s.Data[n.off:], rawStructPointer(pointerOffset(length), sz))
	n.off = n.off.addSize(wordSize)
	return PointerList(n)
}

func (s *Segment) NewBitList(n int32) BitList {
	p, _ := s.create(Size((n+63)/8), Pointer{typ: TypeBitList, length: n})
	return BitList(p)
}

// There is no need to check the object type for lists as:
// 1. Its a list (TypeList, TypeBitList, TypePointerList)
// 2. Its TypeStruct, but then the length is 0
// 3. Its TypeNull, but then the length is 0

func (p Pointer) ToVoidList() VoidList       { return VoidList(p) }
func (p Pointer) ToBitList() BitList         { return BitList(p) }
func (p Pointer) ToInt8List() Int8List       { return Int8List(p) }
func (p Pointer) ToUInt8List() UInt8List     { return UInt8List(p) }
func (p Pointer) ToInt16List() Int16List     { return Int16List(p) }
func (p Pointer) ToUInt16List() UInt16List   { return UInt16List(p) }
func (p Pointer) ToInt32List() Int32List     { return Int32List(p) }
func (p Pointer) ToUInt32List() UInt32List   { return UInt32List(p) }
func (p Pointer) ToFloat32List() Float32List { return Float32List(p) }
func (p Pointer) ToInt64List() Int64List     { return Int64List(p) }
func (p Pointer) ToUInt64List() UInt64List   { return UInt64List(p) }
func (p Pointer) ToFloat64List() Float64List { return Float64List(p) }
func (p Pointer) ToPointerList() PointerList { return PointerList(p) }
func (p Pointer) ToTextList() TextList       { return TextList(p) }
func (p Pointer) ToDataList() DataList       { return DataList(p) }

func (p Pointer) ToListDefault(s *Segment, tagAddr Address) Pointer {
	switch p.typ {
	case TypeList, TypeBitList, TypePointerList:
		return p
	default:
		return s.Root(tagAddr)
	}
}

func (p BitList) At(i int) bool {
	if i < 0 || i >= int(p.length) {
		return false
	}

	switch p.typ {
	case TypePointerList:
		addr := p.off.element(int32(i), wordSize)
		m := p.seg.readPtr(addr)
		// TODO(light): the m.seg.Data[0] seems wrong
		return m.typ == TypeStruct && m.size.DataSize > 0 && (m.seg.Data[0]&1) != 0
	case TypeList:
		addr := p.off.element(int32(i), p.size.totalSize())
		return (p.seg.Data[addr] & 1) != 0
	case TypeBitList:
		addr := p.off.element(int32(i/8), 1)
		mask := byte(1 << uint(i%8))
		return p.seg.Data[addr]&mask != 0
	default:
		return false
	}
}

func (p BitList) Set(i int, v bool) {
	if i < 0 || i >= int(p.length) {
		return
	}

	switch p.typ {
	case TypePointerList:
		addr := p.off.element(int32(i), wordSize)
		m := p.seg.readPtr(addr)
		if m.typ == TypeStruct && m.size.DataSize > 0 {
			if v {
				// TODO(light): the m.seg.Data[0] seems wrong
				m.seg.Data[0] |= 1
			} else {
				m.seg.Data[0] &^= 1
			}
		}
	case TypeList:
		addr := p.off.element(int32(i), p.size.totalSize())
		if v {
			p.seg.Data[addr] |= 1
		} else {
			p.seg.Data[addr] &^= 1
		}
	case TypeBitList:
		addr := p.off.element(int32(i/8), 1)
		mask := byte(1 << uint(i%8))
		if v {
			p.seg.Data[addr] |= mask
		} else {
			p.seg.Data[addr] &^= mask
		}
	}
}

func (p Pointer) listData(i int, sz Size) []byte {
	if i < 0 || i >= int(p.length) {
		return nil
	}

	switch p.typ {
	case TypePointerList:
		m := p.seg.readPtr(p.off.element(int32(i), wordSize))
		if m.typ != TypeStruct || sz > m.size.DataSize {
			return nil
		}
		return m.seg.Data[m.off:]

	case TypeList:
		if sz > p.size.DataSize {
			return nil
		}
		addr := p.off.element(int32(i), p.size.totalSize())
		return p.seg.Data[addr:]

	default: // including TypeBitList as this is only used for 8 bit and larger
		return nil
	}
}

func (p Int8List) At(i int) int8 { return int8(UInt8List(p).At(i)) }
func (p UInt8List) At(i int) uint8 {
	if data := Pointer(p).listData(i, 1); data != nil {
		return data[0]
	} else {
		return 0
	}
}

func (p Int16List) At(i int) int16 { return int16(UInt16List(p).At(i)) }
func (p UInt16List) At(i int) uint16 {
	if data := Pointer(p).listData(i, 2); data != nil {
		return binary.LittleEndian.Uint16(data)
	} else {
		return 0
	}
}

func (p Int32List) At(i int) int32     { return int32(UInt32List(p).At(i)) }
func (p Float32List) At(i int) float32 { return math.Float32frombits(UInt32List(p).At(i)) }
func (p UInt32List) At(i int) uint32 {
	if data := Pointer(p).listData(i, 4); data != nil {
		return binary.LittleEndian.Uint32(data)
	} else {
		return 0
	}
}

func (p Int64List) At(i int) int64     { return int64(UInt64List(p).At(i)) }
func (p Float64List) At(i int) float64 { return math.Float64frombits(UInt64List(p).At(i)) }
func (p UInt64List) At(i int) uint64 {
	if data := Pointer(p).listData(i, 8); data != nil {
		return binary.LittleEndian.Uint64(data)
	} else {
		return 0
	}
}

func (p TextList) At(i int) string { return PointerList(p).At(i).ToText() }
func (p DataList) At(i int) []byte { return PointerList(p).At(i).ToData() }
func (p PointerList) At(i int) Pointer {
	if i < 0 || i >= int(p.length) {
		return Pointer{}
	}

	switch p.typ {
	case TypeList:
		return Pointer{
			seg:   p.seg,
			typ:   TypeStruct,
			off:   p.off.element(int32(i), p.size.totalSize()),
			size:  p.size,
			flags: isListMember,
		}

	case TypePointerList:
		return p.seg.readPtr(p.off.element(int32(i), wordSize))

	case TypeBitList:
		return Pointer{
			seg:   p.seg,
			typ:   TypeStruct,
			off:   p.off.element(int32(i/8), 1),
			flags: pointerFlags(i%8) | isBitListMember,
		}

	default:
		return Pointer{}
	}
}

func (p Int8List) Set(i int, v int8) { UInt8List(p).Set(i, uint8(v)) }
func (p UInt8List) Set(i int, v uint8) {
	if data := Pointer(p).listData(i, 1); data != nil {
		data[0] = v
	}
}

func (p Int16List) Set(i int, v int16) { UInt16List(p).Set(i, uint16(v)) }
func (p UInt16List) Set(i int, v uint16) {
	if data := Pointer(p).listData(i, 2); data != nil {
		binary.LittleEndian.PutUint16(data, v)
	}
}

func (p Int32List) Set(i int, v int32)     { UInt32List(p).Set(i, uint32(v)) }
func (p Float32List) Set(i int, v float32) { UInt32List(p).Set(i, math.Float32bits(v)) }
func (p UInt32List) Set(i int, v uint32) {
	if data := Pointer(p).listData(i, 4); data != nil {
		binary.LittleEndian.PutUint32(data, v)
	}
}

func (p Int64List) Set(i int, v int64)     { UInt64List(p).Set(i, uint64(v)) }
func (p Float64List) Set(i int, v float64) { UInt64List(p).Set(i, math.Float64bits(v)) }
func (p UInt64List) Set(i int, v uint64) {
	if data := Pointer(p).listData(i, 8); data != nil {
		binary.LittleEndian.PutUint64(data, v)
	}
}

func (p TextList) Set(i int, v string) { PointerList(p).Set(i, p.seg.NewText(v)) }
func (p DataList) Set(i int, v []byte) { PointerList(p).Set(i, p.seg.NewData(v)) }
func (p PointerList) Set(i int, src Pointer) error {
	if i < 0 || i >= int(p.length) {
		return nil
	}

	switch p.typ {
	case TypeList:
		if src.typ != TypeStruct {
			src = Pointer{}
		}

		err := copyStructHandlingVersionSkew(Pointer(p), src, nil, 0, int32(i), 0)
		if err != nil {
			return err
		}
		return nil

	case TypePointerList:
		return p.seg.writePtr(p.off.element(int32(i), wordSize), src, nil, 0)

	case TypeBitList:
		boff := BitOffset(i)
		addr := p.off + Address(boff.offset())
		if src.ToStruct().Bit(0) {
			p.seg.Data[addr] |= boff.mask()
		} else {
			p.seg.Data[addr] &^= boff.mask()
		}
		return nil

	default:
		return nil
	}
}

func (p BitList) ToArray() []bool {
	v := make([]bool, p.Len())
	for i := range v {
		v[i] = p.At(i)
	}
	return v
}

func (p UInt8List) ToArray() []uint8 {
	if p.typ == TypeList && p.size.isOneByte() {
		return p.seg.Data[p.off:p.off.element(p.length, 1)]
	}

	v := make([]uint8, p.Len())
	for i := range v {
		v[i] = p.At(i)
	}
	return v
}

func (p Int8List) ToArray() []int8 {
	v := make([]int8, p.Len())
	for i := range v {
		v[i] = p.At(i)
	}
	return v
}

func (p UInt16List) ToArray() []uint16 {
	v := make([]uint16, p.Len())
	for i := range v {
		v[i] = p.At(i)
	}
	return v
}

func (p Int16List) ToArray() []int16 {
	v := make([]int16, p.Len())
	for i := range v {
		v[i] = p.At(i)
	}
	return v
}

func (p UInt32List) ToArray() []uint32 {
	v := make([]uint32, p.Len())
	for i := range v {
		v[i] = p.At(i)
	}
	return v
}

func (p Float32List) ToArray() []float32 {
	v := make([]float32, p.Len())
	for i := range v {
		v[i] = p.At(i)
	}
	return v
}

func (p Int32List) ToArray() []int32 {
	v := make([]int32, p.Len())
	for i := range v {
		v[i] = p.At(i)
	}
	return v
}

func (p Int64List) ToArray() []int64 {
	v := make([]int64, p.Len())
	for i := range v {
		v[i] = p.At(i)
	}
	return v
}

func (p Float64List) ToArray() []float64 {
	v := make([]float64, p.Len())
	for i := range v {
		v[i] = p.At(i)
	}
	return v
}

func (p UInt64List) ToArray() []uint64 {
	v := make([]uint64, p.Len())
	for i := range v {
		v[i] = p.At(i)
	}
	return v
}

func (p TextList) ToArray() []string {
	v := make([]string, p.Len())
	for i := range v {
		v[i] = p.At(i)
	}
	return v
}

func (p DataList) ToArray() [][]byte {
	v := make([][]byte, p.Len())
	for i := range v {
		v[i] = p.At(i)
	}
	return v
}
