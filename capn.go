package capn

import (
	"encoding/binary"
	"errors"
	"math"

	"github.com/glycerine/rbtree"
)

var (
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

func (o ObjectType) String() string {
	switch o {
	case TypeNull:
		return "TypeNull"
	case TypeStruct:
		return "TypeStruct"
	case TypeList:
		return "TypeList"
	case TypePointerList:
		return "TypePointerList"
	case TypeBitList:
		return "TypeBitList"
	default:
		return "Unknown ObjectType"
	}
}

type Message interface {
	NewSegment(minsz int) (*Segment, error)
	Lookup(segid uint32) (*Segment, error)
}

type Segment struct {
	Message  Message
	Data     []uint8
	Id       uint32
	RootDone bool
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

func (o Object) DupWithOff(off int) Object {
	return Object{
		Segment: o.Segment,
		off:     off,
		length:  o.length,
		datasz:  o.datasz,
		ptrs:    o.ptrs,
		typ:     o.typ,
		flags:   o.flags,
	}
}

type Void struct{}
type Struct Object
type VoidList Object
type BitList Object
type Int8List Object
type UInt8List Object
type Int16List Object
type UInt16List Object
type Int32List Object
type UInt32List Object
type Float32List Object
type Int64List Object
type UInt64List Object
type Float64List Object
type PointerList Object
type TextList Object
type DataList Object

func (p VoidList) Len() int    { return p.length }
func (p BitList) Len() int     { return p.length }
func (p Int8List) Len() int    { return p.length }
func (p UInt8List) Len() int   { return p.length }
func (p Int16List) Len() int   { return p.length }
func (p UInt16List) Len() int  { return p.length }
func (p Int32List) Len() int   { return p.length }
func (p UInt32List) Len() int  { return p.length }
func (p Float32List) Len() int { return p.length }
func (p Int64List) Len() int   { return p.length }
func (p UInt64List) Len() int  { return p.length }
func (p Float64List) Len() int { return p.length }
func (p PointerList) Len() int { return p.length }
func (p TextList) Len() int    { return p.length }
func (p DataList) Len() int    { return p.length }

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
	if off+8 > len(s.Data) {
		return Object{}
	}
	return s.readPtr(off)
}

func (s *Segment) NewRoot() (PointerList, int, error) {
	n, err := s.create(8, Object{typ: TypePointerList, length: 1, ptrs: 1, flags: isRoot})
	return PointerList(n), n.off / 8, err
}

func (s *Segment) NewText(v string) Object {
	n := s.NewUInt8List(len(v) + 1)
	copy(n.Segment.Data[n.off:], v)
	return Object(n)
}
func (s *Segment) NewData(v []byte) Object {
	n := s.NewUInt8List(len(v))
	copy(n.Segment.Data[n.off:], v)
	return Object(n)
}

func (s *Segment) NewBitList(sz int) BitList {
	n, _ := s.create((sz+63)/8, Object{typ: TypeBitList, length: sz})
	return BitList(n)
}

//func (s *Segment) NewVoidList(sz int) VoidList       { return VoidList{typ: TypeList, length: sz, datasz: 0} }
func (s *Segment) NewVoidList(sz int) VoidList       { return VoidList(s.newList(0, sz)) }
func (s *Segment) NewInt8List(sz int) Int8List       { return Int8List(s.newList(1, sz)) }
func (s *Segment) NewUInt8List(sz int) UInt8List     { return UInt8List(s.newList(1, sz)) }
func (s *Segment) NewInt16List(sz int) Int16List     { return Int16List(s.newList(2, sz)) }
func (s *Segment) NewUInt16List(sz int) UInt16List   { return UInt16List(s.newList(2, sz)) }
func (s *Segment) NewFloat32List(sz int) Float32List { return Float32List(s.newList(4, sz)) }
func (s *Segment) NewInt32List(sz int) Int32List     { return Int32List(s.newList(4, sz)) }
func (s *Segment) NewUInt32List(sz int) UInt32List   { return UInt32List(s.newList(4, sz)) }
func (s *Segment) NewFloat64List(sz int) Float64List { return Float64List(s.newList(8, sz)) }
func (s *Segment) NewInt64List(sz int) Int64List     { return Int64List(s.newList(8, sz)) }
func (s *Segment) NewUInt64List(sz int) UInt64List   { return UInt64List(s.newList(8, sz)) }
func (s *Segment) newList(datasz, length int) Object {
	n, _ := s.create(datasz*length, Object{typ: TypeList, length: length, datasz: datasz})
	return n
}

func (s *Segment) NewTextList(sz int) TextList { return TextList(s.NewPointerList(sz)) }
func (s *Segment) NewDataList(sz int) DataList { return DataList(s.NewPointerList(sz)) }
func (s *Segment) NewPointerList(sz int) PointerList {
	n, _ := s.create(sz*8, Object{typ: TypePointerList, length: sz, ptrs: 1})
	return PointerList(n)
}

/*
Canonicalization discussion on the Capnproto mailing list:

https://groups.google.com/d/msg/capnproto/V_NysVkvNgs/P5RfeQyMvpkJ

> Q:
> 1: lists of non-trivial but honogeneous structs;
> These will often be 7 (inline-composite). However, is it legal (and: common?) for such to be stored as 6 (pointer)? At the moment I can support both - I just want to make sure I'm not emitting something unusual.

Actually, through Cap'n Proto 0.4.x, it is not only allowed, it is expected. In fact, a struct list can be encoded with any element size, as long as the struct fits. For example, if a struct type Foo has two int16 fields, then List(Foo) can be encoded as a list of 32-bit values.

This rule was intended partly for optimization and partly to solve a common problem, for which I'll give a real example: In schema.capnp, to represent an interface's list of superclasses, I used a List(UInt64) containing the type IDs. Recently, when generics were implemented, I found I now additionally had to specify a type parameterization for each superclass, in addition to its type ID, since they could be generic. Because of the rule that struct lists may be encoded as primitive lists, I was able to simply replace my List(UInt64) with List(Superclass), where Superclass is a struct containing a UInt64 type ID as its @0 field. This change is backwards-compatible. Without this ability, I would have had to use parallel arrays, which would be terrible.

With all that said, in Cap'n Proto 0.5.x, we are making a change: When reading a struct list, the above still applies (a list of non-structs must be accepted), but when writing a struct list, the implementation should always prefer to encode the list using the inline-composite element type.

The reason for this change is so that we can define a canonicalization algorithm for Cap'n Proto messages that does not require knowledge of the schema yet produces stable canonicalizations even as new fields are added to struct types. The canonicalization algorithm needs to know when a value is a struct so that it can truncate trailing zeroes.

List(primitive) -> List(struct) updates are still allowed (because they have proven useful), but with the understanding that this kind of upgrade breaks canonicalization. This seems like an acceptable trade-off.

We also made a second change: Previously, consistent with the above rule, a list of structs where each struct contains a single boolean field was allowed to be encoded as a bit list. We now make a special exception to say that this is *not* allowed, because the implementation burden was too high and we have doubts about whether it would ever be used in practice. In other words, we used to allow List(Bool) -> List(struct) updates where the struct's @0 field was of type Bool, but we no longer allow this. For all types other than Bool, it is still allowed.

*/

const CanonicalizableOn = true

func (s *Segment) NewCompositeList(datasz, ptrs, length int) PointerList {
	if datasz < 0 || datasz > maxDataSize || ptrs < 0 || ptrs > maxPtrs {
		return PointerList{}
	} else if ptrs > 0 || datasz > 8 || CanonicalizableOn {
		datasz = (datasz + 7) &^ 7
		n, _ := s.create(8+length*(datasz+8*ptrs), Object{typ: TypeList, length: length, datasz: datasz, ptrs: ptrs, flags: isCompositeList})
		n.off += 8
		hdr := structPointer | uint64(length)<<2 | uint64(datasz/8)<<32 | uint64(ptrs)<<48
		binary.LittleEndian.PutUint64(s.Data[n.off-8:], hdr)
		return PointerList(n)
	} else if datasz > 4 {
		datasz = (datasz + 7) &^ 7
	} else if datasz > 2 {
		datasz = (datasz + 3) &^ 3
	}

	n, _ := s.create(length*(datasz+8*ptrs), Object{typ: TypeList, length: length, datasz: datasz, ptrs: ptrs})
	return PointerList(n)
}

func (s *Segment) NewRootStruct(datasz, ptrs int) Struct {
	r, _, err := s.NewRoot()
	if err != nil {
		return Struct{}
	}
	v := s.NewStruct(datasz, ptrs)
	r.Set(0, Object(v))
	return v
}

func (s *Segment) NewStruct(datasz, ptrs int) Struct {
	if datasz < 0 || datasz > maxDataSize || ptrs < 0 || ptrs > maxPtrs {
		return Struct{}
	}
	datasz = (datasz + 7) &^ 7
	n, _ := s.create(datasz+ptrs*8, Object{typ: TypeStruct, datasz: datasz, ptrs: ptrs})
	return Struct(n)
}

// NewStructAR (AutoRoot): experimental Root setting: assumes the
// struct is the root iff it is the first allocation in a segment.
func (s *Segment) NewStructAR(datasz, ptrs int) Struct {
	if s.RootDone {
		return s.NewStruct(datasz, ptrs)
	} else {
		s.RootDone = true
		return s.NewRootStruct(datasz, ptrs)
	}
}

// sz is in bytes
func (s *Segment) create(sz int, n Object) (Object, error) {
	// rounded up to word-boundary number of bytes:
	sz = (sz + 7) &^ 7

	if uint64(sz) > uint64(math.MaxUint32)-8 {
		return Object{}, ErrOverlarge
	}

	if s == nil {
		s = NewBuffer(nil)
	}

	tag := false
	if len(s.Data)+sz > cap(s.Data) {
		// If we can't fit the data in the current segment, we always
		// return a far pointer to a tag in the new segment.
		if (n.flags & isRoot) != 0 {
			tag = true
			sz += 8
		}
		news, err := s.Message.NewSegment(sz)
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
	s.Data = s.Data[:len(s.Data)+sz] // NewSegment() makes this promise

	if tag {
		n.off += 8
		binary.LittleEndian.PutUint64(s.Data[n.off-8:], uint64(n.value(n.off-8)))
		n.flags |= hasPointerTag
	}

	for i := n.off; i < len(s.Data); i++ {
		s.Data[i] = 0
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

func (p Object) ToText() string { return p.ToTextDefault("") }
func (p Object) ToTextDefault(def string) string {
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

func (p Object) ToVoidList() VoidList       { return VoidList(p) }
func (p Object) ToBitList() BitList         { return BitList(p) }
func (p Object) ToInt8List() Int8List       { return Int8List(p) }
func (p Object) ToUInt8List() UInt8List     { return UInt8List(p) }
func (p Object) ToInt16List() Int16List     { return Int16List(p) }
func (p Object) ToUInt16List() UInt16List   { return UInt16List(p) }
func (p Object) ToInt32List() Int32List     { return Int32List(p) }
func (p Object) ToUInt32List() UInt32List   { return UInt32List(p) }
func (p Object) ToFloat32List() Float32List { return Float32List(p) }
func (p Object) ToInt64List() Int64List     { return Int64List(p) }
func (p Object) ToUInt64List() UInt64List   { return UInt64List(p) }
func (p Object) ToFloat64List() Float64List { return Float64List(p) }
func (p Object) ToPointerList() PointerList { return PointerList(p) }
func (p Object) ToTextList() TextList       { return TextList(p) }
func (p Object) ToDataList() DataList       { return DataList(p) }

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

func (p Struct) SetObject(i int, src Object) {
	if uint(i) < uint(p.ptrs) {
		//replaceMe := p.Segment.readPtr(p.off + p.datasz + i*8)
		//copyStructHandlingVersionSkew(replaceMe, src, nil, 0, 0, 0)
		p.Segment.writePtr(p.off+p.datasz+i*8, src, nil, 0)
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
		return binary.LittleEndian.Uint16(p.Segment.Data[uint(p.off)+uint(off):])
	} else {
		return 0
	}
}

func (p Struct) Get32(off int) uint32 {
	if off < p.datasz {
		return binary.LittleEndian.Uint32(p.Segment.Data[uint(p.off)+uint(off):])
	} else {
		return 0
	}
}

func (p Struct) Get64(off int) uint64 {
	if off < p.datasz {
		return binary.LittleEndian.Uint64(p.Segment.Data[uint(p.off)+uint(off):])
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
		binary.LittleEndian.PutUint16(p.Segment.Data[uint(p.off)+uint(off):], v)
	}
}

func (p Struct) Set32(off int, v uint32) {
	if uint(off) < uint(p.datasz) {
		binary.LittleEndian.PutUint32(p.Segment.Data[uint(p.off)+uint(off):], v)
	}
}

func (p Struct) Set64(off int, v uint64) {
	if uint(off) < uint(p.datasz) {
		binary.LittleEndian.PutUint64(p.Segment.Data[uint(p.off)+uint(off):], v)
	}
}

func (p BitList) At(i int) bool {
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
		return (p.Segment.Data[p.off+i/8] & (1 << uint(i%8))) != 0
	default:
		return false
	}
}

func (p BitList) Set(i int, v bool) {
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
			p.Segment.Data[p.off+i/8] |= 1 << uint(i%8)
		} else {
			p.Segment.Data[p.off+i/8] &^= 1 << uint(i%8)
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

func (p Int8List) At(i int) int8 { return int8(UInt8List(p).At(i)) }
func (p UInt8List) At(i int) uint8 {
	if data := Object(p).listData(i, 8); data != nil {
		return data[0]
	} else {
		return 0
	}
}

func (p Int16List) At(i int) int16 { return int16(UInt16List(p).At(i)) }
func (p UInt16List) At(i int) uint16 {
	if data := Object(p).listData(i, 16); data != nil {
		return binary.LittleEndian.Uint16(data)
	} else {
		return 0
	}
}

func (p Int32List) At(i int) int32     { return int32(UInt32List(p).At(i)) }
func (p Float32List) At(i int) float32 { return math.Float32frombits(UInt32List(p).At(i)) }
func (p UInt32List) At(i int) uint32 {
	if data := Object(p).listData(i, 32); data != nil {
		return binary.LittleEndian.Uint32(data)
	} else {
		return 0
	}
}

func (p Int64List) At(i int) int64     { return int64(UInt64List(p).At(i)) }
func (p Float64List) At(i int) float64 { return math.Float64frombits(UInt64List(p).At(i)) }
func (p UInt64List) At(i int) uint64 {
	if data := Object(p).listData(i, 64); data != nil {
		return binary.LittleEndian.Uint64(data)
	} else {
		return 0
	}
}

func (p Int8List) Set(i int, v int8) { UInt8List(p).Set(i, uint8(v)) }
func (p UInt8List) Set(i int, v uint8) {
	if data := Object(p).listData(i, 8); data != nil {
		data[0] = v
	}
}

func (p Int16List) Set(i int, v int16) { UInt16List(p).Set(i, uint16(v)) }
func (p UInt16List) Set(i int, v uint16) {
	if data := Object(p).listData(i, 16); data != nil {
		binary.LittleEndian.PutUint16(data, v)
	}
}

func (p Int32List) Set(i int, v int32)     { UInt32List(p).Set(i, uint32(v)) }
func (p Float32List) Set(i int, v float32) { UInt32List(p).Set(i, math.Float32bits(v)) }
func (p UInt32List) Set(i int, v uint32) {
	if data := Object(p).listData(i, 32); data != nil {
		binary.LittleEndian.PutUint32(data, v)
	}
}

func (p Int64List) Set(i int, v int64)     { UInt64List(p).Set(i, uint64(v)) }
func (p Float64List) Set(i int, v float64) { UInt64List(p).Set(i, math.Float64bits(v)) }
func (p UInt64List) Set(i int, v uint64) {
	if data := Object(p).listData(i, 64); data != nil {
		binary.LittleEndian.PutUint64(data, v)
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
	if p.typ == TypeList && p.datasz == 1 && p.ptrs == 0 {
		return p.Segment.Data[p.off : p.off+p.length]
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

func (p UInt16List) ToEnumArray() *[]uint16 {
	v := make([]uint16, p.Len())
	for i := range v {
		v[i] = p.At(i)
	}
	return &v
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

func (p Int64List) ToIntArray() []int {
	v := make([]int, p.Len())
	for i := range v {
		v[i] = int(p.At(i))
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

func (p PointerList) ToArray() *[]Object {
	v := make([]Object, p.Len())
	for i := range v {
		v[i] = p.At(i)
	}
	return &v
}

func (p TextList) At(i int) string { return PointerList(p).At(i).ToText() }
func (p DataList) At(i int) []byte { return PointerList(p).At(i).ToData() }
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

// listpos allows us to use this routine for copying elements between lists
func copyStructHandlingVersionSkew(dest Object, src Object, copies *rbtree.Tree, depth int, destListPos int, srcListPos int) error {

	// handle VoidList destinations
	if dest.Segment == nil {
		return nil
	}

	destElemSz := dest.datasz + dest.ptrs*8
	srcElemSz := src.datasz + src.ptrs*8

	// handle assignment into a size-zero object/empty struct/old version
	//if destElemSz == 0 {
	// return nil
	// }

	destListInc := destElemSz * destListPos
	srcListInc := srcElemSz * srcListPos

	toData := dest.Segment.Data[dest.off+destListInc : dest.off+destListInc+dest.datasz]

	//fmt.Printf("\n\n debug: destElemSz = %d, srcElemSz = %d, destListInc = %d, srcListInc = %d, toData = %#v, len(toData)=%d\n", destElemSz, srcElemSz, destListInc, srcListInc, toData, len(toData))

	// Q: how does version handling happen here, when the
	//    desination toData[] slice can be bigger or smaller
	//    than the source data slice, which is in
	//    src.Segment.Data[src.off:src.off+src.datasz] ?
	//
	// A: Newer fields only come *after* old fields. Note that
	//    copy only copies min(len(src), len(dst)) size,
	//    and then we manually zero the rest in the for loop
	//    that writes toData[j] = 0.
	//

	// data section:
	//fmt.Printf("\n\n  debug: len(src.Segment.Data) = %d, src.off(%d)+srcListInc(%d) = %d\n", len(src.Segment.Data), src.off, srcListInc, src.off+srcListInc)
	from := src.Segment.Data[src.off+srcListInc : src.off+srcListInc+src.datasz]
	//fmt.Printf("\n\n  debug: len(src.Segment.Data) = %d, src.off(%d)+srcListInc(%d) = %d,  len(from)=%d\n", len(src.Segment.Data), src.off, srcListInc, src.off+srcListInc, len(from))

	copyCount := copy(toData, from)
	//fmt.Printf("\n\n  debug: copyCount = %d\n", copyCount)
	toData = toData[copyCount:]
	for j := range toData {
		toData[j] = 0
	}

	// ptrs section:

	// version handling: we ignore any extra-newer-pointers in src,
	// i.e. the case when srcPtrSize > dstPtrSize, by only
	// running j over the size of dstPtrSize, the destination size.
	srcPtrSize := src.ptrs * 8
	dstPtrSize := int(dest.ptrs * 8)
	for j := 0; j < dstPtrSize; j += 8 {
		if j < srcPtrSize {
			m := src.Segment.readPtr(src.off + srcListInc + src.datasz + j)
			//fmt.Printf("debug: PointerList.Set(i=%d, src=%#v). source ptr is m = %#v\n", i, src, m)
			if err := dest.Segment.writePtr(dest.off+destListInc+dest.datasz+j, m, copies, depth+1); err != nil {
				return err
			}
		} else {
			// destination p is a newer version than source
			//  so these extra new pointer fields in p must be zeroed.
			binary.LittleEndian.PutUint64(dest.Segment.Data[dest.off+destListInc+dest.datasz+j:], 0)
		}
	}
	// Nothing more here: so any other pointers in srcPtrSize beyond
	// those in dstPtrSize are ignored and discarded.

	return nil

} // end copyStructHandlingVersionSkew()

func (p TextList) Set(i int, v string) { PointerList(p).Set(i, p.Segment.NewText(v)) }
func (p DataList) Set(i int, v []byte) { PointerList(p).Set(i, p.Segment.NewData(v)) }
func (p PointerList) Set(i int, src Object) error {
	if i < 0 || i >= p.length {
		return nil
	}

	switch p.typ {
	case TypeList:
		if src.typ != TypeStruct {
			src = Object{}
		}

		err := copyStructHandlingVersionSkew(Object(p), src, nil, 0, i, 0)
		if err != nil {
			return err
		}
		return nil

	case TypePointerList:
		return p.Segment.writePtr(p.off+i*8, src, nil, 0)

	case TypeBitList:
		if src.ToStruct().Get1(0) {
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
		return s.Message.Lookup(id)
	} else {
		return s, nil
	}
}

func (s *Segment) readPtr(off int) Object {
	var err error
	val := pointer(binary.LittleEndian.Uint64(s.Data[off:]))

	//fmt.Printf("readPtr see val= %x\n", val)

	switch val.pointerType() {
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

		far := binary.LittleEndian.Uint64(s.Data[faroff:])
		tag := binary.LittleEndian.Uint64(s.Data[faroff+8:])

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
		val = pointer(uint64(uint32(far)>>3<<2) | tag)

	case farPointer:
		segid := uint32(val >> 32)
		faroff := int((uint32(val) >> 3) * 8)

		if s, err = s.lookupSegment(segid); err != nil || faroff+8 > len(s.Data) {
			return Object{}
		}

		off = faroff
		val = pointer(binary.LittleEndian.Uint64(s.Data[faroff:]))
	}

	// Be wary of overflow. Offset is 30 bits signed. List size is 29 bits
	// unsigned. For both of these we need to check in terms of words if
	// using 32 bit maths as bits or bytes will overflow.
	switch val.pointerType() {
	case structPointer:
		offw := off/8 + 1 + val.structSignedOffset()

		if offw < 0 || offw >= len(s.Data)/8 {
			return Object{}
		}

		p := Object{
			Segment: s,
			typ:     TypeStruct,
			off:     offw * 8,
			datasz:  pointer(val).structC() * 8,
			ptrs:    pointer(val).structD(),
		}

		if p.off+p.datasz+p.ptrs*8 > len(s.Data) {
			return Object{}
		}

		return p

	case listPointer:
		offw := off/8 + 1 + val.structSignedOffset()
		//fmt.Printf("offw = %d,  len(s.Data)=%d", offw, len(s.Data))
		if val.listC() != voidList {
			if offw < 0 || offw >= len(s.Data)/8 {
				return Object{}
			}
		}

		p := Object{
			Segment: s,
			typ:     TypeList,
			off:     offw * 8,
			length:  int(uint32(val >> 35)),
		}

		words := p.length

		switch val.listC() {
		case voidList:
			//fmt.Printf("we see a voidList with len: %d\n", p.length)
			return p
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
			hdr := binary.LittleEndian.Uint64(p.Segment.Data[p.off:])
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

// pointer is an encoded pointer.
type pointer uint64

// Pointer types
const (
	structPointer    = 0
	listPointer      = 1
	farPointer       = 2
	doubleFarPointer = 6
)

// List pointer types
const (
	voidList      = 0
	bit1List      = 1
	byte1List     = 2
	byte2List     = 3
	byte4List     = 4
	byte8List     = 5
	pointerList   = 6
	compositeList = 7
)

func (p pointer) pointerType() int {
	t := p & 3
	if t == farPointer {
		return int(p & 7)
	}
	return int(t)
}

func (p pointer) structC() int {
	return int(uint16(p >> 32))
}

func (p pointer) structD() int {
	return int(uint16(p >> 48))
}

func (p pointer) listC() int {
	return int((p >> 32) & 7)
}

// used in orable30BitOffsetPart() and pointer.structSignedOffset()
const zerohi32 pointer = ^(^0 << 32)

// orable30BitOffsetPart(): get an or-able value that handles sign
// conversion. Creates part B in a struct (or list) pointer, leaving
// parts A, C, and D completely zeroed in the returned uint64.
//
// From the spec:
//
// lsb                      struct pointer                       msb
// +-+-----------------------------+---------------+---------------+
// |A|             B               |       C       |       D       |
// +-+-----------------------------+---------------+---------------+
//
// A (2 bits) = 0, to indicate that this is a struct pointer.
// B (30 bits) = Offset, in words, from the end of the pointer to the
//     start of the struct's data section.  Signed.
// C (16 bits) = Size of the struct's data section, in words.
// D (16 bits) = Size of the struct's pointer section, in words.
//
// (B is the same for list pointers, but C and D have different size
// and meaning)
//
func orable30BitOffsetPart(signedOff int) pointer {
	d32 := int32(signedOff) << 2
	return pointer(d32) & zerohi32
}

// and convert in the other direction, extracting the count from
// the B section into an int
func (p pointer) structSignedOffset() int {
	u64 := p & zerohi32
	u32 := uint32(u64)
	s32 := int32(u32) >> 2
	return int(s32)
}

func (p Object) value(off int) pointer {
	d := orable30BitOffsetPart(p.off/8 - off/8 - 1)
	//fmt.Printf(" debug4 in value(off=%d): p.off/8 = %d, off/8 = %d,  and p.off/8 - off/8 -1 = %d    d=%v\n", off, p.off/8, off/8, p.off/8-off/8-1, d)

	switch p.typ {
	case TypeStruct:
		return structPointer | d | pointer(p.datasz/8)<<32 | pointer(p.ptrs)<<48
	case TypePointerList:
		return listPointer | d | pointerList<<32 | pointer(p.length)<<35
	case TypeList:
		if (p.flags & isCompositeList) != 0 {
			d -= 1 << 2 // p.off points to the data not the header
			return listPointer | d | compositeList<<32 | pointer(p.length*(p.datasz/8+p.ptrs))<<35
		}

		switch p.datasz {
		case 0:
			return listPointer | d | voidList<<32 | pointer(p.length)<<35
		case 1:
			return listPointer | d | byte1List<<32 | pointer(p.length)<<35
		case 2:
			return listPointer | d | byte2List<<32 | pointer(p.length)<<35
		case 4:
			return listPointer | d | byte4List<<32 | pointer(p.length)<<35
		case 8:
			return listPointer | d | byte8List<<32 | pointer(p.length)<<35
		default:
			panic(errListSize)
		}

	case TypeBitList:
		return listPointer | d | bit1List<<32 | pointer(p.length)<<35
	case TypeNull:
		return 0
	default:
		panic(errObjectType)
	}
}

/*
lsb                       list pointer                        msb
+-+-----------------------------+--+----------------------------+
|A|             B               |C |             D              |
+-+-----------------------------+--+----------------------------+

A (2 bits) = 1, to indicate that this is a list pointer.
B (30 bits) = Offset, in words, from the end of the pointer to the
    start of the first element of the list.  Signed.
C (3 bits) = Size of each element:
    0 = 0 (e.g. List(Void))
    1 = 1 bit
    2 = 1 byte
    3 = 2 bytes
    4 = 4 bytes
    5 = 8 bytes (non-pointer)
    6 = 8 bytes (pointer)
    7 = composite (see below)
D (29 bits) = Number of elements in the list, except when C is 7
    (see below).

The pointed-to values are tightly-packed. In particular, Bools are packed bit-by-bit in little-endian order (the first bit is the least-significant bit of the first byte).

Lists of structs use the smallest element size in which the struct can fit. So, a list of structs that each contain two UInt8 fields and nothing else could be encoded with C = 3 (2-byte elements). A list of structs that each contain a single Text field would be encoded as C = 6 (pointer elements). A list of structs that each contain a single Bool field would be encoded using C = 1 (1-bit elements). A list structs which are each more than one word in size must be be encoded using C = 7 (composite).

When C = 7, the elements of the list are fixed-width composite values – usually, structs. In this case, the list content is prefixed by a "tag" word that describes each individual element. The tag has the same layout as a struct pointer, except that the pointer offset (B) instead indicates the number of elements in the list. Meanwhile, section (D) of the list pointer – which normally would store this element count – instead stores the total number of words in the list (not counting the tag word). The reason we store a word count in the pointer rather than an element count is to ensure that the extents of the list’s location can always be determined by inspecting the pointer alone, without having to look at the tag; this may allow more-efficient prefetching in some use cases. The reason we don’t store struct lists as a list of pointers is because doing so would take significantly more space (an extra pointer per element) and may be less cache-friendly.
*/

func (s *Segment) farPtrValue(farType int, off int) pointer {
	return pointer(farType) | pointer(off) | (pointer(s.Id) << 32)
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

func isEmptyStruct(src Object) bool {
	if src.typ == TypeStruct && src.length == 0 && src.datasz == 0 && src.ptrs == 0 && src.flags == 0 {
		return true
	}
	return false
}

func (destSeg *Segment) writePtr(off int, src Object, copies *rbtree.Tree, depth int) error {

	// handle size-zero Objects/empty structs
	if src.Segment == nil {
		return nil
	}
	srcSeg := src.Segment

	if src.typ == TypeNull || isEmptyStruct(src) {
		binary.LittleEndian.PutUint64(destSeg.Data[off:], 0)
		return nil

	} else if destSeg == srcSeg {
		// Same segment
		binary.LittleEndian.PutUint64(destSeg.Data[off:], uint64(src.value(off)))
		return nil

	} else if destSeg.Message != srcSeg.Message || (src.flags&isListMember) != 0 || (src.flags&isBitListMember) != 0 {
		// We need to clone the target.

		if depth >= 32 {
			return ErrCopyDepth
		}

		// First see if the ptr has already been copied
		if copies == nil {
			copies = rbtree.NewTree(compare)
		}

		key := offset{
			id:   srcSeg.Id,
			boff: int64(src.off) * 8,
			bend: int64(src.dataEnd()) * 8,
			newval: Object{
				typ:    src.typ,
				length: src.length,
				datasz: src.datasz,
				ptrs:   src.ptrs,
				flags:  (src.flags & isCompositeList),
			},
		}

		if (src.flags & isBitListMember) != 0 {
			key.boff += int64(src.flags & bitOffsetMask)
			key.bend = key.boff + 1
			key.newval.datasz = 8
		}

		if (src.flags & isCompositeList) != 0 {
			key.boff -= 64 //  Q: what the heck does this do? why is it here? A: Accounts for the Tag word, perhaps because dataEnd() does not.
		}

		iter := copies.FindLE(key)

		if key.bend > key.boff {
			if !iter.NegativeLimit() {
				other := iter.Item().(offset)
				if key.id == other.id {
					if key.boff == other.boff && key.bend == other.bend {
						return destSeg.writePtr(off, other.newval, nil, depth+1)
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
		n, err := destSeg.create(int((key.bend-key.boff)/8), key.newval)
		if err != nil {
			return err
		}

		// n is possibly in a new segment, if destSeg was full.
		newSeg := n.Segment

		if (n.flags & isCompositeList) != 0 {
			copy(newSeg.Data[n.off:], srcSeg.Data[src.off-8:src.off])
			n.off += 8
		}

		key.newval = n
		copies.Insert(key)

		//fmt.Printf(" .....  need to clone target: key.newval: %#v of type '%s'\n", key.newval, key.newval.typ.String())

		switch src.typ {
		case TypeStruct:
			if (src.flags & isBitListMember) != 0 {
				if (srcSeg.Data[src.off] & (1 << (src.flags & bitOffsetMask))) != 0 {
					newSeg.Data[n.off] = 1
				} else {
					newSeg.Data[n.off] = 0
				}

				for i := range newSeg.Data[n.off+1 : n.off+8] {
					newSeg.Data[i] = 0
				}
			} else {
				dest := Object{
					Segment: newSeg,
					off:     n.off,
					datasz:  n.datasz,
					ptrs:    n.ptrs,
				}

				if err := copyStructHandlingVersionSkew(dest, src, copies, depth, 0, 0); err != nil {
					return err
				}
			}

		case TypeList:
			// recognize Data and Text, both List(Byte), as special cases for speed.
			if n.datasz == 1 && n.ptrs == 0 && src.datasz == 1 && src.ptrs == 0 {
				//fmt.Printf("\n\n    *** special case for Text and Data kicking in *** \n\n")
				copy(newSeg.Data[n.off:], srcSeg.Data[src.off:src.off+n.length+1])
				break
			}

			dest := Object{
				Segment: newSeg,
				off:     n.off,
				datasz:  n.datasz,
				ptrs:    n.ptrs,
			}
			for i := 0; i < n.length; i++ {
				if err := copyStructHandlingVersionSkew(dest, src, copies, depth, i, i); err != nil {
					return err
				}
			}

		case TypePointerList:
			for i := 0; i < n.length; i++ {
				c := srcSeg.readPtr(src.off + i*8)
				if err := newSeg.writePtr(n.off+i*8, c, copies, depth+1); err != nil {
					return err
				}
			}

		case TypeBitList:
			copy(newSeg.Data[n.off:], srcSeg.Data[src.off:src.off+src.datasz])
		}
		return destSeg.writePtr(off, key.newval, nil, depth+1)

	} else if (src.flags & hasPointerTag) != 0 {
		// By lucky chance, the data has a tag in front of it. This
		// happens when create had to move the data to a new segment.
		binary.LittleEndian.PutUint64(destSeg.Data[off:], uint64(srcSeg.farPtrValue(farPointer, src.off-8)))
		return nil

	} else if len(srcSeg.Data)+8 <= cap(srcSeg.Data) {
		// Have room in the target for a tag
		binary.LittleEndian.PutUint64(srcSeg.Data[len(srcSeg.Data):], uint64(src.value(len(srcSeg.Data))))
		binary.LittleEndian.PutUint64(destSeg.Data[off:], uint64(srcSeg.farPtrValue(farPointer, len(srcSeg.Data))))
		srcSeg.Data = srcSeg.Data[:len(srcSeg.Data)+8]
		return nil

	} else {
		// Need to create a double far pointer. Try and create it in
		// the originating segment if we can.
		t := destSeg
		if len(t.Data)+16 > cap(t.Data) {
			var err error
			if t, err = t.Message.NewSegment(16); err != nil {
				return err
			}
		}

		binary.LittleEndian.PutUint64(t.Data[len(t.Data):], uint64(srcSeg.farPtrValue(farPointer, src.off)))
		binary.LittleEndian.PutUint64(t.Data[len(t.Data)+8:], uint64(src.value(src.off-8)))
		binary.LittleEndian.PutUint64(destSeg.Data[off:], uint64(t.farPtrValue(doubleFarPointer, len(t.Data))))
		t.Data = t.Data[:len(t.Data)+16]
		return nil
	}
}
