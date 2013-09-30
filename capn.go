package capn

import (
	"encoding/binary"
	"errors"
	"github.com/yasushi-saito/rbtree"
	"math"
)

var (
	little16    = binary.LittleEndian.Uint16
	little32    = binary.LittleEndian.Uint32
	little64    = binary.LittleEndian.Uint64
	putLittle16 = binary.LittleEndian.PutUint16
	putLittle32 = binary.LittleEndian.PutUint32
	putLittle64 = binary.LittleEndian.PutUint64

	ErrOverlarge   = errors.New("capn: overlarge struct/list")
	ErrOutOfBounds = errors.New("capn: write out of bounds")
	ErrCopyDepth   = errors.New("capn: copy depth too large")
	ErrOverlap     = errors.New("capn: overlapping data on copy")
	errListSize    = errors.New("capn: invalid list size")
	errObjectType  = errors.New("capn: invalid object type")
)

type ObjectType uint8

const (
	TypeNull ObjectType = iota
	TypeStruct
	TypeList
	TypePointerList
	TypeBitList
)

type Session interface {
	NewSegment(minsz int) (*Segment, error)
	Lookup(segid uint32) (*Segment, error)
}

type Segment struct {
	Session Session
	Data    []uint8
	Id      uint32
}

type Object struct {
	Segment *Segment
	off     int // in bytes
	length  int
	datasz  int // in bytes
	ptrs    int
	typ     ObjectType
	flags   uint
}

type Void struct{}
type Struct Object
type VoidList Object
type List1 Object
type ListI8 Object
type ListU8 Object
type ListI16 Object
type ListU16 Object
type ListI32 Object
type ListU32 Object
type ListF32 Object
type ListI64 Object
type ListU64 Object
type ListF64 Object
type PointerList Object

func (p VoidList) Len() int    {return p.length}
func (p List1) Len() int       { return p.length }
func (p ListI8) Len() int       { return p.length }
func (p ListU8) Len() int       { return p.length }
func (p ListI16) Len() int      { return p.length }
func (p ListU16) Len() int      { return p.length }
func (p ListI32) Len() int      { return p.length }
func (p ListU32) Len() int      { return p.length }
func (p ListF32) Len() int      { return p.length }
func (p ListI64) Len() int      { return p.length }
func (p ListU64) Len() int      { return p.length }
func (p ListF64) Len() int      { return p.length }
func (p PointerList) Len() int { return p.length }

func (p Object) HasData() bool {
	switch p.typ {
	case TypeList:
		return p.length > 0 && (p.datasz != 0 || p.ptrs != 0)
	case TypePointerList:
		return p.length > 0
	case TypeBitList:
		return p.length > 0
	case TypeStruct:
		return p.datasz != 0 || p.ptrs != 0
	default:
		return false
	}
}

const (
	maxDataSize = 0xFFFF * 8
	maxPtrs     = 0xFFFF

	// flags
	bitOffsetMask   = 7
	isBitListMember = 8
	isListMember    = 16
	isCompositeList = 32
	isRoot          = 64
	hasPointerTag   = 128
)

func (s *Segment) Root(off int) Object {
	if off + 8 > len(s.Data) {
		return Object{}
	}
	return s.readPtr(off)
}

func (s *Segment) NewRoot() (PointerList, int, error) {
	n, err := s.create(8, Object{typ: TypePointerList, length: 1, ptrs: 1, flags: isRoot})
	return PointerList(n), n.off / 8, err
}

func (s *Segment) NewPointerList(length int) (PointerList, error) {
	n, err := s.create(length*8, Object{typ: TypePointerList, length: length, ptrs: 1})
	return PointerList(n), err
}

func (s *Segment) NewVoidList(length int) (VoidList, error) {
	return VoidList{typ: TypeList, length: length, datasz: 0}, nil
}

func (s *Segment) NewList1(length int) (List1, error) {
	n, err := s.create((length+63)/8, Object{typ: TypeBitList, length: length})
	return List1(n), err
}

func (s *Segment) NewListI8(length int) (ListI8, error) {
	n, err := s.create(length, Object{typ: TypeList, length: length, datasz: 1})
	return ListI8(n), err
}

func (s *Segment) NewListU8(length int) (ListU8, error) {
	n, err := s.create(length, Object{typ: TypeList, length: length, datasz: 1})
	return ListU8(n), err
}

func (s *Segment) NewListI16(length int) (ListI16, error) {
	n, err := s.create(2*length, Object{typ: TypeList, length: length, datasz: 2})
	return ListI16(n), err
}

func (s *Segment) NewListU16(length int) (ListU16, error) {
	n, err := s.create(2*length, Object{typ: TypeList, length: length, datasz: 2})
	return ListU16(n), err
}

func (s *Segment) NewListI32(length int) (ListI32, error) {
	n, err := s.create(4*length, Object{typ: TypeList, length: length, datasz: 4})
	return ListI32(n), err
}

func (s *Segment) NewListU32(length int) (ListU32, error) {
	n, err := s.create(4*length, Object{typ: TypeList, length: length, datasz: 4})
	return ListU32(n), err
}

func (s *Segment) NewListF32(length int) (ListF32, error) {
	n, err := s.create(4*length, Object{typ: TypeList, length: length, datasz: 4})
	return ListF32(n), err
}

func (s *Segment) NewListI64(length int) (ListI64, error) {
	n, err := s.create(8*length, Object{typ: TypeList, length: length, datasz: 8})
	return ListI64(n), err
}

func (s *Segment) NewListU64(length int) (ListU64, error) {
	n, err := s.create(8*length, Object{typ: TypeList, length: length, datasz: 8})
	return ListU64(n), err
}


func (s *Segment) NewListF64(length int) (ListF64, error) {
	n, err := s.create(8*length, Object{typ: TypeList, length: length, datasz: 8})
	return ListF64(n), err
}

func (s *Segment) NewCompositeList(datasz, ptrs, length int) (PointerList, error) {
	flags := uint(0)
	if datasz < 0 || datasz > maxDataSize || ptrs < 0 || ptrs > maxPtrs {
		return PointerList{}, ErrOverlarge
	} else if ptrs > 0 || datasz > 8 {
		datasz = (datasz + 7) &^ 7
		flags |= isCompositeList
	} else if datasz > 4 {
		datasz = (datasz + 7) &^ 7
	} else if datasz > 2 {
		datasz = (datasz + 3) &^ 3
	}
	n, err := s.create(length*(datasz+8*ptrs), Object{typ: TypeList, length: length, datasz: datasz, ptrs: ptrs, flags: flags})
	return PointerList(n), err
}

func (s *Segment) NewStruct(datasz, ptrs int) (Struct, error) {
	if datasz < 0 || datasz > maxDataSize || ptrs < 0 || ptrs > maxPtrs {
		return Struct{}, ErrOverlarge
	}
	datasz = (datasz + 7) &^ 7
	n, err := s.create(datasz, Object{typ: TypeStruct, datasz: datasz, ptrs: ptrs})
	return Struct(n), err
}

func (s *Segment) create(sz int, n Object) (Object, error) {
	sz = (sz + 7) &^ 7

	if uint64(sz) > uint64(math.MaxUint32)-8 {
		return Object{}, ErrOverlarge
	}

	tag := false
	if len(s.Data)+sz > cap(s.Data) {
		// If we can't fit the data in the current segment, we always
		// return a far pointer to a tag in the new segment.
		if (n.flags & isRoot) != 0 {
			tag = true
			sz += 8
		}
		news, err := s.Session.NewSegment(sz)
		if err != nil {
			return Object{}, err
		}

		// NewSegment is allowed to grow the segment and return it. In
		// which case we don't want to create a tag.
		if tag && news == s {
			sz -= 8
			tag = false
		}

		s = news
	}

	n.Segment = s
	n.off = len(s.Data)
	s.Data = s.Data[:len(s.Data)+sz]

	if tag {
		n.off += 8
		putLittle64(s.Data[n.off-8:], n.value(n.off-8))
		n.flags |= hasPointerTag
	}

	return n, nil
}

func (p Object) Type() ObjectType { return p.typ }

func (p Object) ToStruct() Struct {
	if p.typ == TypeStruct {
		return Struct(p)
	} else {
		return Struct{}
	}
}

func (p Object) ToStructDefault(s *Segment, tagoff int) Struct {
	if p.typ == TypeStruct {
		return Struct(p)
	} else {
		return s.Root(tagoff).ToStruct()
	}
}

func (p Object) ToString() string { return p.ToStringDefault("") }
func (p Object) ToStringDefault(def string) string {
	if p.typ != TypeList || p.datasz != 1 || p.length == 0 || p.Segment.Data[p.off+p.length-1] != 0 {
		return def
	}

	return string(p.Segment.Data[p.off : p.off+p.length-1])
}

func (p Object) ToData() []byte { return p.ToDataDefault(nil) }
func (p Object) ToDataDefault(def []byte) []byte {
	if p.typ != TypeList || p.datasz != 1 {
		return def
	}

	return p.Segment.Data[p.off : p.off+p.length]
}

// There is no need to check the object type for lists as:
// 1. Its a list (TypeList, TypeBitList, TypePointerList)
// 2. Its TypeStruct, but then the length is 0
// 3. Its TypeNull, but then the length is 0

func (p Object) ToVoidList() VoidList { return VoidList(p) }
func (p Object) ToListU1() List1             { return List1(p) }
func (p Object) ToListI8() ListI8             { return ListI8(p) }
func (p Object) ToListU8() ListU8             { return ListU8(p) }
func (p Object) ToListI16() ListI16           { return ListI16(p) }
func (p Object) ToListU16() ListU16           { return ListU16(p) }
func (p Object) ToListI32() ListI32           { return ListI32(p) }
func (p Object) ToListU32() ListU32           { return ListU32(p) }
func (p Object) ToListF32() ListF32           { return ListF32(p) }
func (p Object) ToListI64() ListI64           { return ListI64(p) }
func (p Object) ToListU64() ListU64           { return ListU64(p) }
func (p Object) ToListF64() ListF64           { return ListF64(p) }
func (p Object) ToPointerList() PointerList { return PointerList(p) }

func (p Object) ToListDefault(s *Segment, tagoff int) Object {
	switch p.typ {
	case TypeList, TypeBitList, TypePointerList:
		return p
	default:
		return s.Root(tagoff)
	}
}

func (p Object) ToObjectDefault(s *Segment, tagoff int) Object {
	if p.typ == TypeNull {
		return s.Root(tagoff)
	} else {
		return p
	}
}

func (p Struct) GetObject(off int) Object {
	if uint(off) < uint(p.ptrs) {
		return p.Segment.readPtr(p.off + p.datasz + off*8)
	} else {
		return Object{}
	}
}

func (p Struct) SetObject(i int, tgt Object) {
	if uint(i) < uint(p.ptrs) {
		p.Segment.writePtr(p.off+p.datasz+i*8, tgt, nil, 0)
	}
}

func (p Struct) Get1(bitoff int) bool {
	off := uint(p.off*8 + bitoff)

	if bitoff == 0 && (p.flags&isBitListMember) != 0 {
		off += p.flags & bitOffsetMask
	} else if bitoff < 0 || bitoff >= p.datasz*8 {
		return false
	}

	return p.Segment.Data[off/8]&(1<<uint(off%8)) != 0
}

func (p Struct) Set1(bitoff int, v bool) {
	off := uint(p.off*8 + bitoff)

	if bitoff == 0 && (p.flags&isBitListMember) != 0 {
		off += p.flags & bitOffsetMask
	} else if bitoff < 0 || bitoff >= p.datasz*8 {
		return
	}

	if v {
		p.Segment.Data[off/8] |= 1 << (off % 8)
	} else {
		p.Segment.Data[off/8] &^= 1 << (off % 8)
	}
}

func (p Struct) Get8(off int) uint8 {
	if off < p.datasz {
		return p.Segment.Data[uint(p.off)+uint(off)]
	} else {
		return 0
	}
}

func (p Struct) Get16(off int) uint16 {
	if off < p.datasz {
		return little16(p.Segment.Data[uint(p.off)+uint(off):])
	} else {
		return 0
	}
}

func (p Struct) Get32(off int) uint32 {
	if off < p.datasz {
		return little32(p.Segment.Data[uint(p.off)+uint(off):])
	} else {
		return 0
	}
}

func (p Struct) Get64(off int) uint64 {
	if off < p.datasz {
		return little64(p.Segment.Data[uint(p.off)+uint(off):])
	} else {
		return 0
	}
}

func (p Struct) Set8(off int, v uint8) {
	if uint(off) < uint(p.datasz) {
		p.Segment.Data[uint(p.off)+uint(off)] = v
	}
}

func (p Struct) Set16(off int, v uint16) {
	if uint(off) < uint(p.datasz) {
		putLittle16(p.Segment.Data[uint(p.off)+uint(off):], v)
	}
}

func (p Struct) Set32(off int, v uint32) {
	if uint(off) < uint(p.datasz) {
		putLittle32(p.Segment.Data[uint(p.off)+uint(off):], v)
	}
}

func (p Struct) Set64(off int, v uint64) {
	if uint(off) < uint(p.datasz) {
		putLittle64(p.Segment.Data[uint(p.off)+uint(off):], v)
	}
}

func (p List1) At(i int) bool {
	if i < 0 || i >= p.length {
		return false
	}

	switch p.typ {
	case TypePointerList:
		m := p.Segment.readPtr(p.off + i*8)
		return m.typ == TypeStruct && m.datasz > 0 && (m.Segment.Data[0]&1) != 0
	case TypeList:
		off := p.off + i*(p.datasz+p.ptrs*8)
		return (p.Segment.Data[off] & 1) != 0
	case TypeBitList:
		return (p.Segment.Data[i/8] & (1 << uint(i%8))) != 0
	default:
		return false
	}
}

func (p List1) Set(i int, v bool) {
	if i < 0 || i >= p.length {
		return
	}

	switch p.typ {
	case TypePointerList:
		m := p.Segment.readPtr(p.off + i*8)
		if m.typ == TypeStruct && m.datasz > 0 {
			if v {
				m.Segment.Data[0] |= 1
			} else {
				m.Segment.Data[0] &^= 1
			}
		}
	case TypeList:
		off := p.off + i*(p.datasz+p.ptrs*8)
		if v {
			p.Segment.Data[off] |= 1
		} else {
			p.Segment.Data[off] &^= 1
		}
	case TypeBitList:
		if v {
			p.Segment.Data[i/8] |= 1 << uint(i%8)
		} else {
			p.Segment.Data[i/8] &^= 1 << uint(i%8)
		}
	}
}

func (p Object) listData(i int, sz int) []byte {
	if i < 0 || i >= p.length {
		return nil
	}

	switch p.typ {
	case TypePointerList:
		m := p.Segment.readPtr(p.off + i*8)
		if m.typ != TypeStruct || sz > m.datasz*8 {
			return nil
		}
		return m.Segment.Data[m.off:]

	case TypeList:
		if sz > p.datasz*8 {
			return nil
		}
		off := p.off + i*(p.datasz+p.ptrs*8)
		return p.Segment.Data[off:]

	default: // including TypeBitList as this is only used for 8 bit and larger
		return nil
	}
}

func (p ListI8) At(i int) int8 { return int8(ListU8(p).At(i)) }
func (p ListU8) At(i int) uint8 {
	if data := Object(p).listData(i, 8); data != nil {
		return data[0]
	} else {
		return 0
	}
}

func (p ListI16) At(i int) int16 { return int16(ListU16(p).At(i)) }
func (p ListU16) At(i int) uint16 {
	if data := Object(p).listData(i, 16); data != nil {
		return little16(data)
	} else {
		return 0
	}
}

func (p ListI32) At(i int) int32 { return int32(ListU32(p).At(i)) }
func (p ListF32) At(i int) float32 { return math.Float32frombits(ListU32(p).At(i)) }
func (p ListU32) At(i int) uint32 {
	if data := Object(p).listData(i, 32); data != nil {
		return little32(data)
	} else {
		return 0
	}
}

func (p ListI64) At(i int) int64 { return int64(ListU64(p).At(i)) }
func (p ListF64) At(i int) float64 { return math.Float64frombits(ListU64(p).At(i)) }
func (p ListU64) At(i int) uint64 {
	if data := Object(p).listData(i, 64); data != nil {
		return little64(data)
	} else {
		return 0
	}
}

func (p ListI8) Set(i int, v int8) { ListU8(p).Set(i, uint8(v)) }
func (p ListU8) Set(i int, v uint8) {
	if data := Object(p).listData(i, 8); data != nil {
		data[0] = v
	}
}

func (p ListI16) Set(i int, v int16) { ListU16(p).Set(i, uint16(v)) }
func (p ListU16) Set(i int, v uint16) {
	if data := Object(p).listData(i, 16); data != nil {
		putLittle16(data, v)
	}
}

func (p ListI32) Set(i int, v int32) { ListU32(p).Set(i, uint32(v)) }
func (p ListF32) Set(i int, v float32) { ListU32(p).Set(i, math.Float32bits(v)) }
func (p ListU32) Set(i int, v uint32) {
	if data := Object(p).listData(i, 32); data != nil {
		putLittle32(data, v)
	}
}

func (p ListI64) Set(i int, v int64) { ListU64(p).Set(i, uint64(v)) }
func (p ListF64) Set(i int, v float64) { ListU64(p).Set(i, math.Float64bits(v)) }
func (p ListU64) Set(i int, v uint64) {
	if data := Object(p).listData(i, 64); data != nil {
		putLittle64(data, v)
	}
}

func (p PointerList) ToArray() *[]Object {
	v := make([]Object, p.Len())
	for i := range v {
		v[i] = p.At(i)
	}
	return &v
}

func (p PointerList) At(i int) Object {
	if i < 0 || i >= p.length {
		return Object{}
	}

	switch p.typ {
	case TypeList:
		return Object{
			Segment: p.Segment,
			typ:     TypeStruct,
			off:     p.off + i*(p.datasz+p.ptrs*8),
			datasz:  p.datasz,
			ptrs:    p.ptrs,
			flags:   isListMember,
		}

	case TypePointerList:
		return p.Segment.readPtr(p.off + i*8)

	case TypeBitList:
		return Object{
			Segment: p.Segment,
			typ:     TypeStruct,
			off:     p.off + i/8,
			flags:   uint(i%8) | isBitListMember,
			datasz:  0,
			ptrs:    0,
		}

	default:
		return Object{}
	}
}

func (p PointerList) Set(i int, tgt Object) error {
	if i < 0 || i >= p.length {
		return nil
	}

	switch p.typ {
	case TypeList:
		if tgt.typ != TypeStruct {
			tgt = Object{}
		}

		off := p.off + i*(p.datasz+p.ptrs*8)
		data := p.Segment.Data[off:p.datasz]
		data = data[copy(data, tgt.Segment.Data[tgt.off:tgt.off+tgt.datasz]):]
		for j := range data {
			data[j] = 0
		}

		j := int(p.ptrs * 8)
		for ; j > tgt.ptrs*8 && j >= 0; j -= 8 {
			putLittle64(p.Segment.Data[off+p.datasz+j:], 0)
		}

		for ; j >= 0; j -= 8 {
			m := tgt.Segment.readPtr(tgt.off+tgt.datasz+j)
			if err := p.Segment.writePtr(off+p.datasz+j, m, nil, 0); err != nil {
				return err
			}
		}
		return nil

	case TypePointerList:
		return p.Segment.writePtr(p.off+i*8, tgt, nil, 0)

	case TypeBitList:
		if tgt.ToStruct().Get1(0) {
			p.Segment.Data[p.off+i/8] |= 1 << uint(i%8)
		} else {
			p.Segment.Data[p.off+i/8] &^= 1 << uint(i%8)
		}
		return nil

	default:
		return nil
	}
}

func (s *Segment) lookupSegment(id uint32) (*Segment, error) {
	if s.Id != id {
		return s.Session.Lookup(id)
	} else {
		return s, nil
	}
}

const (
	structPointer    = 0
	listPointer      = 1
	farPointer       = 2
	doubleFarPointer = 6

	voidList      = 0
	bit1List      = 1
	byte1List     = 2
	byte2List     = 3
	byte4List     = 4
	byte8List     = 5
	pointerList   = 6
	compositeList = 7
)

func (s *Segment) readPtr(off int) Object {
	var err error
	val := little64(s.Data[off:])

	switch val & 7 {
	case doubleFarPointer:
		// A double far pointer points to a double pointer, where the
		// first points to the actual data, and the second is the tag
		// that would normally be placed right before the data (offset
		// == 0).

		faroff := int((uint32(val) >> 3) * 8)
		segid := uint32(val >> 32)

		if s, err = s.lookupSegment(segid); err != nil || uint(faroff)+16 > uint(len(s.Data)) {
			return Object{}
		}

		far := little64(s.Data[faroff:])
		tag := little64(s.Data[faroff+8:])

		// The far tag should not be another double and the tag should
		// be struct/list with a 0 offset.
		if far&7 != farPointer || uint32(tag) > listPointer {
			return Object{}
		}

		segid = uint32(far >> 32)
		if s, err = s.lookupSegment(segid); err != nil {
			return Object{}
		}

		// -8 because far pointers reference from the start of the
		// segment, but offsets reference the end of the pointer data.
		off = -8
		val = uint64(uint32(far)>>3<<2) | tag

	case farPointer:
		segid := uint32(val >> 32)
		faroff := int((uint32(val) >> 3) * 8)

		if s, err = s.lookupSegment(segid); err != nil || faroff+8 > len(s.Data) {
			return Object{}
		}

		off = faroff
		val = little64(s.Data[faroff:])
	}

	// Be wary of overflow. Offset is 30 bits signed. List size is 29 bits
	// unsigned. For both of these we need to check in terms of words if
	// using 32 bit maths as bits or bytes will overflow.
	switch val & 3 {
	case structPointer:
		offw := off/8 + 1 + int(uint32(val)>>2)
		if offw < 0 || offw >= len(s.Data)/8 {
			return Object{}
		}

		p := Object{
			Segment: s,
			typ:     TypeStruct,
			off:     offw * 8,
			datasz:  int(uint16(val>>32)) * 8,
			ptrs:    int(uint16(val >> 48)),
		}

		if p.off+p.datasz+p.ptrs*8 > len(s.Data) {
			return Object{}
		}

		return p

	case listPointer:
		offw := off/8 + 1 + int(uint32(val))>>2
		if offw < 0 || offw >= len(s.Data)/8 {
			return Object{}
		}

		p := Object{
			Segment: s,
			typ:     TypeList,
			off:     offw * 8,
			length:  int(uint32(val >> 35)),
		}

		words := p.length

		switch (val >> 32) & 7 {
		case bit1List:
			p.typ = TypeBitList
			words = (p.length + 63) / 64
		case byte1List:
			p.datasz = 1
			words = (p.length + 7) / 8
		case byte2List:
			p.datasz = 2
			words = (p.length + 3) / 4
		case byte4List:
			p.datasz = 4
			words = (p.length + 1) / 2
		case byte8List:
			p.datasz = 8
		case pointerList:
			p.typ = TypePointerList
		case compositeList:
			hdr := little64(p.Segment.Data[p.off:])
			p.off += 8
			if hdr&2 != structPointer {
				return Object{}
			}

			p.flags |= isCompositeList
			p.length = int(uint32(hdr) >> 2)
			p.datasz = int(uint16(hdr>>32)) * 8
			p.ptrs = int(uint16(hdr >> 48))

			// Jump up to 64bit as length is 30 bits, datasz+ptrs is 17 bit
			if uint64(p.length)*uint64(p.datasz/8+p.ptrs) != uint64(words) {
				return Object{}
			}
		}

		// Largest possible message is 30 bits * 1 word, with either a
		// composite, ptr, or 8 byte list. If we do a size check using
		// bits or bytes, we overflow.
		if words > len(s.Data)/8-offw {
			return Object{}
		}

		return p

	default:
		return Object{}
	}
}

func (p Object) value(off int) uint64 {
	d := uint64(p.off/8-off/8-1) << 2

	switch p.typ {
	case TypeStruct:
		return structPointer | d | uint64(p.datasz/8)<<32 | uint64(p.ptrs)<<48
	case TypePointerList:
		return listPointer | d | pointerList<<32 | uint64(p.length)<<35
	case TypeList:
		if (p.flags & isCompositeList) != 0 {
			return listPointer | d | compositeList<<32 | uint64(p.length*(p.datasz/8+p.ptrs))<<35
		}

		switch p.datasz {
		case 0:
			return listPointer | d | voidList<<32 | uint64(p.length)<<35
		case 1:
			return listPointer | d | byte1List<<32 | uint64(p.length)<<35
		case 2:
			return listPointer | d | byte2List<<32 | uint64(p.length)<<35
		case 4:
			return listPointer | d | byte4List<<32 | uint64(p.length)<<35
		case 8:
			return listPointer | d | byte8List<<32 | uint64(p.length)<<35
		default:
			panic(errListSize)
		}

	case TypeBitList:
		return listPointer | d | bit1List<<32 | uint64(p.length)<<35
	case TypeNull:
		return 0
	default:
		panic(errObjectType)
	}
}

func (s *Segment) farPtrValue(farType int, off int) uint64 {
	return uint64(farType) | uint64(off) | (uint64(s.Id) << 32)
}

type offset struct {
	id         uint32
	boff, bend int64 // in bits
	newval     Object
}

func compare(a, b rbtree.Item) int {
	ao := a.(offset)
	bo := b.(offset)
	if ao.id != bo.id {
		return int(ao.id - bo.id)
	} else if ao.boff > bo.boff {
		return 1
	} else if ao.boff < bo.boff {
		return -1
	} else {
		return 0
	}
}

func (p Object) dataEnd() int {
	switch p.typ {
	case TypeList:
		return p.off + p.length*(p.datasz+p.ptrs*8)
	case TypePointerList:
		return p.off + p.length*8
	case TypeStruct:
		return p.off + p.datasz + p.ptrs*8
	case TypeBitList:
		return p.off + (p.length+7)/8
	default:
		return p.off
	}
}

func (s *Segment) writePtr(off int, p Object, copies *rbtree.Tree, depth int) error {
	ps := p.Segment

	if p.typ == TypeNull {
		putLittle64(s.Data[off:], 0)
		return nil

	} else if s == p.Segment {
		// Same segment
		putLittle64(s.Data[off:], p.value(off))
		return nil

	} else if s.Session != ps.Session || (p.flags&isListMember) != 0 || (p.flags&isBitListMember) != 0 {
		// We need to clone the target.

		if depth >= 32 {
			return ErrCopyDepth
		}

		// First see if the ptr has already been copied
		if copies == nil {
			copies = rbtree.NewTree(compare)
		}

		key := offset{
			id:   ps.Id,
			boff: int64(p.off) * 8,
			bend: int64(p.dataEnd()) * 8,
			newval: Object{
				typ:    p.typ,
				length: p.length,
				datasz: p.datasz,
				ptrs:   p.ptrs,
				flags: (p.flags & isCompositeList),
			},
		}

		if (p.flags & isBitListMember) != 0 {
			key.boff += int64(p.flags & bitOffsetMask)
			key.bend = key.boff + 1
			key.newval.datasz = 8
		}

		if (p.flags & isCompositeList) != 0 {
			key.boff -= 64
		}

		iter := copies.FindLE(key)

		if key.bend > key.boff {
			if !iter.NegativeLimit() {
				other := iter.Item().(offset)
				if key.id == other.id {
					if key.boff == other.boff && key.bend == other.bend {
						return s.writePtr(off, other.newval, nil, depth+1)
					} else if other.bend >= key.bend {
						return ErrOverlap
					}
				}
			}

			iter = iter.Next()

			if !iter.Limit() {
				other := iter.Item().(offset)
				if key.id == other.id && other.boff < key.bend {
					return ErrOverlap
				}
			}
		}

		// No copy nor overlap found, so we need to clone the target
		n, err := s.create(int((key.bend-key.boff)/8), key.newval)
		if err != nil {
			return err
		}

		ns := n.Segment

		if (n.flags & isCompositeList) != 0 {
			copy(ns.Data[n.off:], ps.Data[p.off-8:p.off])
			n.off += 8
		}

		key.newval = n
		copies.Insert(key)

		switch p.typ {
		case TypeStruct:
			if (p.flags & isBitListMember) != 0 {
				if (ps.Data[p.off] & (1 << (p.flags & bitOffsetMask))) != 0 {
					ns.Data[n.off] = 1
				} else {
					ns.Data[n.off] = 0
				}

				for i := range ns.Data[n.off+1 : n.off+8] {
					ns.Data[i] = 0
				}
			} else {
				copy(ns.Data[n.off:], ps.Data[p.off:p.off+p.datasz])
				for i := 0; i < n.ptrs; i++ {
					c := ps.readPtr(p.off + p.datasz + i*8)
					if err := ns.writePtr(n.off+n.datasz+i*8, c, copies, depth+1); err != nil {
						return err
					}
				}
			}

		case TypeList:
			for i := 0; i < n.length; i++ {
				o := i*(n.datasz+n.ptrs*8)
				copy(ns.Data[n.off+o:], ps.Data[p.off+o:p.off+o+n.datasz])
				o += n.datasz

				for j := 0; j < n.ptrs; j++ {
					c := ps.readPtr(p.off + o)
					if err := ns.writePtr(n.off+o, c, copies, depth+1); err != nil {
						return err
					}
					o += 8
				}
			}

		case TypePointerList:
			for i := 0; i < n.ptrs; i++ {
				c := ps.readPtr(p.off + i*8)
				if err := ns.writePtr(n.off+i*8, c, copies, depth+1); err != nil {
					return err
				}
			}

		case TypeBitList:
			copy(ns.Data[n.off:], ps.Data[p.off:p.off+p.datasz])
		}

		return s.writePtr(off, key.newval, nil, depth+1)

	} else if (p.flags & hasPointerTag) != 0 {
		// By lucky chance, the data has a tag in front of it. This
		// happens when create had to move the data to a new segment.
		putLittle64(s.Data[off:], ps.farPtrValue(farPointer, p.off-8))
		return nil

	} else if len(ps.Data)+8 <= cap(ps.Data) {
		// Have room in the target for a tag
		putLittle64(ps.Data[len(ps.Data):], p.value(len(ps.Data)))
		putLittle64(s.Data[off:], ps.farPtrValue(farPointer, len(ps.Data)))
		ps.Data = ps.Data[:len(ps.Data)+8]
		return nil

	} else {
		// Need to create a double far pointer. Try and create it in
		// the originating segment if we can.
		t := s
		if len(t.Data)+16 > cap(t.Data) {
			var err error
			if t, err = t.Session.NewSegment(16); err != nil {
				return err
			}
		}

		putLittle64(t.Data[len(t.Data):], ps.farPtrValue(farPointer, p.off))
		putLittle64(t.Data[len(t.Data)+8:], p.value(p.off-8))
		putLittle64(s.Data[off:], t.farPtrValue(doubleFarPointer, len(t.Data)))
		t.Data = t.Data[:len(t.Data)+16]
		return nil
	}
}
