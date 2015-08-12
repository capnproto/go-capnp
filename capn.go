package capnp

import (
	"encoding/binary"
	"errors"
	"math"
	"strconv"

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

// PointerType is an enumeration of pointer types.
type PointerType uint8

// Pointer types.
const (
	TypeNull PointerType = iota
	TypeStruct
	TypeList
	TypePointerList
	TypeBitList
	TypeInterface
)

// String returns the constant name of the pointer type.
func (t PointerType) String() string {
	switch t {
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
	case TypeInterface:
		return "TypeInterface"
	default:
		return "PointerType(" + strconv.FormatUint(uint64(t), 10) + ")"
	}
}

type Message interface {
	NewSegment(minsz Size) (*Segment, error)
	Lookup(segid SegmentID) (*Segment, error)

	CapTable() []Client
	AddCap(c Client) CapabilityID
}

// A Segment is an allocation arena for Cap'n Proto objects.
// It is part of a Message, which can contain other segments that
// reference each other.
type Segment struct {
	Message  Message
	Data     []byte
	Id       SegmentID
	RootDone bool
}

func (s *Segment) inBounds(addr Address) bool {
	return addr < Address(len(s.Data))
}

func (s *Segment) regionInBounds(base Address, sz Size) bool {
	return base.addSize(sz) <= Address(len(s.Data))
}

// Address returns the address the pointer references.
func (p Pointer) Address() Address {
	return p.off
}

// A SegmentID is a numeric identifier for a Segment.
type SegmentID uint32

// ObjectSize records section sizes for a struct or list.
type ObjectSize struct {
	DataSize     Size
	PointerCount uint16
}

// isZero reports whether sz is the zero size.
func (sz ObjectSize) isZero() bool {
	return sz.DataSize == 0 && sz.PointerCount == 0
}

// isOneByte reports whether the object size is one byte (for Text/Data element sizes).
func (sz ObjectSize) isOneByte() bool {
	return sz.DataSize == 1 && sz.PointerCount == 0
}

// isValid reports whether sz's fields are in range.
func (sz ObjectSize) isValid() bool {
	return sz.DataSize <= 0xffff*wordSize
}

// pointerSize returns the number of bytes the pointer section occupies.
func (sz ObjectSize) pointerSize() Size {
	return Size(sz.PointerCount) * wordSize
}

// totalSize returns the number of bytes that the object occupies.
func (sz ObjectSize) totalSize() Size {
	return sz.DataSize + sz.pointerSize()
}

// dataWordCount returns the number of words in the data section.
func (sz ObjectSize) dataWordCount() int32 {
	if sz.DataSize%wordSize != 0 {
		panic("data size not aligned by word")
	}
	return int32(sz.DataSize / wordSize)
}

// totalWordCount returns the number of words that the object occupies.
func (sz ObjectSize) totalWordCount() int32 {
	return sz.dataWordCount() + int32(sz.PointerCount)
}

// A Pointer is a reference to a Cap'n Proto struct, list, or interface.
// The zero value is a null pointer.
type Pointer struct {
	seg    *Segment
	off    Address
	length int32
	size   ObjectSize
	typ    PointerType
	flags  pointerFlags
	cap    CapabilityID
}

// IsNull reports whether this is a null pointer.
func (p Pointer) IsNull() bool {
	return p.typ == TypeNull
}

// Segment returns the segment this pointer came from.
func (p Pointer) Segment() *Segment {
	return p.seg
}

func (p Pointer) DupWithOff(off Address) Pointer {
	return Pointer{
		seg:    p.seg,
		off:    off,
		length: p.length,
		size:   p.size,
		typ:    p.typ,
		flags:  p.flags,
		cap:    p.cap,
	}
}

type Void struct{}
type Struct Pointer
type VoidList Pointer
type BitList Pointer
type Int8List Pointer
type UInt8List Pointer
type Int16List Pointer
type UInt16List Pointer
type Int32List Pointer
type UInt32List Pointer
type Float32List Pointer
type Int64List Pointer
type UInt64List Pointer
type Float64List Pointer
type PointerList Pointer
type TextList Pointer
type DataList Pointer

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

// Segment returns the segment this pointer came from.
func (s Struct) Segment() *Segment {
	return s.seg
}

// HasData reports whether the object referenced by p has non-zero size.
func (p Pointer) HasData() bool {
	switch p.typ {
	case TypeList:
		return p.length > 0 && !p.size.isZero()
	case TypePointerList:
		return p.length > 0
	case TypeBitList:
		return p.length > 0
	case TypeStruct:
		return !p.size.isZero()
	case TypeInterface:
		return true
	default:
		return false
	}
}

// pointerFlags is a bitmask of flags for a pointer.
type pointerFlags uint8

// Pointer flags.
const (
	isBitListMember pointerFlags = 8 << iota
	isListMember
	isCompositeList
	isRoot
	hasPointerTag

	bitOffsetMask pointerFlags = 7
)

func (flags pointerFlags) bitOffset() bitOffset {
	return bitOffset(flags & bitOffsetMask)
}

func (s *Segment) Root(addr Address) Pointer {
	if !s.regionInBounds(addr, wordSize) {
		return Pointer{}
	}
	return s.readPtr(addr)
}

func (s *Segment) NewRoot() (PointerList, error) {
	n, err := s.create(wordSize, Pointer{
		typ:    TypePointerList,
		length: 1,
		size:   ObjectSize{PointerCount: 1},
		flags:  isRoot,
	})
	return PointerList(n), err
}

const maxDataSize = 1<<32 - 1

func (s *Segment) NewText(v string) Pointer {
	if int64(len(v))+1 > maxDataSize {
		panic(ErrOverlarge)
	}
	n := s.NewUInt8List(int32(len(v) + 1))
	copy(n.seg.Data[n.off:], v)
	return Pointer(n)
}
func (s *Segment) NewData(v []byte) Pointer {
	if int64(len(v)) > maxDataSize {
		panic(ErrOverlarge)
	}
	n := s.NewUInt8List(int32(len(v)))
	copy(n.seg.Data[n.off:], v)
	return Pointer(n)
}

func (s *Segment) NewBitList(n int32) BitList {
	p, _ := s.create(Size((n+63)/8), Pointer{typ: TypeBitList, length: n})
	return BitList(p)
}

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

func (s *Segment) NewRootStruct(sz ObjectSize) Struct {
	r, err := s.NewRoot()
	if err != nil {
		return Struct{}
	}
	v := s.NewStruct(sz)
	r.Set(0, Pointer(v))
	return v
}

func (s *Segment) NewStruct(sz ObjectSize) Struct {
	if !sz.isValid() {
		return Struct{}
	}
	sz.DataSize = sz.DataSize.padToWord()
	n, _ := s.create(sz.totalSize(), Pointer{typ: TypeStruct, size: sz})
	return Struct(n)
}

// NewStructAR (AutoRoot): experimental Root setting: assumes the
// struct is the root iff it is the first allocation in a segment.
func (s *Segment) NewStructAR(sz ObjectSize) Struct {
	if s.RootDone {
		return s.NewStruct(sz)
	} else {
		s.RootDone = true
		return s.NewRootStruct(sz)
	}
}

// create allocates sz bytes in the segment, using n as a template.
func (s *Segment) create(sz Size, n Pointer) (Pointer, error) {
	sz = sz.padToWord()

	// TODO(light): this can overflow easily
	if uint64(sz) > uint64(math.MaxUint32)-8 {
		return Pointer{}, ErrOverlarge
	}

	if s == nil {
		s = NewBuffer(nil)
	}

	tag := false
	newSize := Size(len(s.Data)) + sz
	addr := Address(len(s.Data))
	end := addr.addSize(sz)
	if newSize > Size(cap(s.Data)) {
		// If we can't fit the data in the current segment, we always
		// return a far pointer to a tag in the new segment.
		if (n.flags & isRoot) != 0 {
			tag = true
			sz += wordSize
		}
		news, err := s.Message.NewSegment(sz)
		if err != nil {
			return Pointer{}, err
		}

		// NewSegment is allowed to grow the segment and return it. In
		// which case we don't want to create a tag.
		if tag && news == s {
			sz -= wordSize
			tag = false
		}

		s = news
	}

	n.seg = s
	n.off = addr
	s.Data = s.Data[:newSize] // NewSegment() makes this promise

	if tag {
		putRawPointer(s.Data[addr:], n.value(addr))
		n.off = n.off.addSize(wordSize)
		n.flags |= hasPointerTag
	}

	for i := n.off; i < end; i++ {
		s.Data[i] = 0
	}

	return n, nil
}

// NewInterface returns a new interface pointer.
func (s *Segment) NewInterface(cap CapabilityID) Interface {
	return Interface(Pointer{
		seg: s,
		typ: TypeInterface,
		cap: cap,
	})
}

// Type returns the type of object referenced by p.
func (p Pointer) Type() PointerType { return p.typ }

func (p Pointer) ToStruct() Struct {
	if p.typ == TypeStruct {
		return Struct(p)
	} else {
		return Struct{}
	}
}

func (p Pointer) ToStructDefault(s *Segment, tagAddr Address) Struct {
	if p.typ == TypeStruct {
		return Struct(p)
	} else {
		return s.Root(tagAddr).ToStruct()
	}
}

func (p Pointer) ToInterface() Interface {
	if p.typ == TypeInterface {
		return Interface(p)
	} else {
		return Interface{}
	}
}

func (p Pointer) ToText() string { return p.ToTextDefault("") }
func (p Pointer) ToTextDefault(def string) string {
	lastAddr := p.off.element(p.length-1, 1)
	if p.typ != TypeList || !p.size.isOneByte() || p.length == 0 || p.seg.Data[lastAddr] != 0 {
		return def
	}

	return string(p.seg.Data[p.off:lastAddr])
}

func (p Pointer) ToData() []byte { return p.ToDataDefault(nil) }
func (p Pointer) ToDataDefault(def []byte) []byte {
	if p.typ != TypeList || !p.size.isOneByte() {
		return def
	}

	return p.seg.Data[p.off:p.off.element(p.length, 1)]
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

func (p Pointer) ToObjectDefault(s *Segment, tagAddr Address) Pointer {
	if p.typ == TypeNull {
		return s.Root(tagAddr)
	} else {
		return p
	}
}

// Pointer returns the i'th pointer in the struct.
func (p Struct) Pointer(i uint16) Pointer {
	if i < p.size.PointerCount {
		return p.seg.readPtr(p.pointerAddress(i))
	} else {
		return Pointer{}
	}
}

// SetPointer sets the i'th pointer in the struct to src.
func (p Struct) SetPointer(i uint16, src Pointer) {
	if i < p.size.PointerCount {
		//replaceMe := p.seg.readPtr(p.off + p.datasz + i*8)
		//copyStructHandlingVersionSkew(replaceMe, src, nil, 0, 0, 0)
		p.seg.writePtr(p.pointerAddress(i), src, nil, 0)
	}
}

func (p Struct) pointerAddress(i uint16) Address {
	ptrStart := p.off.addSize(p.size.DataSize)
	return ptrStart.element(int32(i), wordSize)
}

func (p Struct) bitOffset(bitoff uint32) bitOffset {
	if bitoff == 0 && (p.flags&isBitListMember) != 0 {
		return p.flags.bitOffset()
	}
	return bitOffset(bitoff)
}

// Bit returns the bit that is bitoff bits from the start of the struct.
func (p Struct) Bit(bitoff uint32) bool {
	o := p.bitOffset(bitoff)
	if o >= bitOffset(p.size.DataSize*8) {
		return false
	}
	addr := p.off.addOffset(o.offset())
	mask := o.mask()
	return p.seg.Data[addr]&mask != 0
}

// SetBit sets the bit that is bitoff bits from the start of the struct to v.
func (p Struct) SetBit(bitoff uint32, v bool) {
	o := p.bitOffset(bitoff)
	if o >= bitOffset(p.size.DataSize*8) {
		return
	}
	addr := p.off.addOffset(o.offset())
	mask := o.mask()
	if v {
		p.seg.Data[addr] |= mask
	} else {
		p.seg.Data[addr] &^= mask
	}
}

func (p Struct) dataAddress(off DataOffset, sz Size) (addr Address, ok bool) {
	if Size(off)+sz > p.size.DataSize {
		return 0, false
	}
	return p.off.addOffset(off), true
}

// Uint8 returns an 8-bit integer from the struct's data section.
func (p Struct) Uint8(off DataOffset) uint8 {
	addr, ok := p.dataAddress(off, 1)
	if !ok {
		return 0
	}
	return p.seg.Data[addr]
}

// Uint16 returns a 16-bit integer from the struct's data section.
func (p Struct) Uint16(off DataOffset) uint16 {
	addr, ok := p.dataAddress(off, 2)
	if !ok {
		return 0
	}
	return binary.LittleEndian.Uint16(p.seg.Data[addr:])
}

// Uint32 returns a 32-bit integer from the struct's data section.
func (p Struct) Uint32(off DataOffset) uint32 {
	addr, ok := p.dataAddress(off, 4)
	if !ok {
		return 0
	}
	return binary.LittleEndian.Uint32(p.seg.Data[addr:])
}

// Uint64 returns a 64-bit integer from the struct's data section.
func (p Struct) Uint64(off DataOffset) uint64 {
	addr, ok := p.dataAddress(off, 8)
	if !ok {
		return 0
	}
	return binary.LittleEndian.Uint64(p.seg.Data[addr:])
}

// SetUint8 sets the 8-bit integer that is off bytes from the start of the struct to v.
func (p Struct) SetUint8(off DataOffset, v uint8) {
	addr, ok := p.dataAddress(off, 1)
	if !ok {
		panic(ErrOutOfBounds)
	}
	p.seg.Data[addr] = v
}

// SetUint16 sets the 16-bit integer that is off bytes from the start of the struct to v.
func (p Struct) SetUint16(off DataOffset, v uint16) {
	addr, ok := p.dataAddress(off, 2)
	if !ok {
		panic(ErrOutOfBounds)
	}
	binary.LittleEndian.PutUint16(p.seg.Data[addr:], v)
}

// SetUint32 sets the 32-bit integer that is off bytes from the start of the struct to v.
func (p Struct) SetUint32(off DataOffset, v uint32) {
	addr, ok := p.dataAddress(off, 4)
	if !ok {
		panic(ErrOutOfBounds)
	}
	binary.LittleEndian.PutUint32(p.seg.Data[addr:], v)
}

// SetUint64 sets the 64-bit integer that is off bytes from the start of the struct to v.
func (p Struct) SetUint64(off DataOffset, v uint64) {
	addr, ok := p.dataAddress(off, 8)
	if !ok {
		panic(ErrOutOfBounds)
	}
	binary.LittleEndian.PutUint64(p.seg.Data[addr:], v)
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

// listpos allows us to use this routine for copying elements between lists
func copyStructHandlingVersionSkew(dest, src Pointer, copies *rbtree.Tree, depth int, destListPos, srcListPos int32) error {
	// handle VoidList destinations
	if dest.seg == nil {
		return nil
	}

	destBase := dest.off.element(destListPos, dest.size.totalSize())
	destDataEnd := destBase.addSize(dest.size.DataSize)
	srcBase := src.off.element(srcListPos, src.size.totalSize())
	srcDataEnd := srcBase.addSize(src.size.DataSize)

	// Q: how does version handling happen here, when the
	//    desination toData[] slice can be bigger or smaller
	//    than the source data slice, which is in
	//    src.seg.Data[src.off:src.off+src.size.DataSize] ?
	//
	// A: Newer fields only come *after* old fields. Note that
	//    copy only copies min(len(src), len(dst)) size,
	//    and then we manually zero the rest in the for loop
	//    that writes toData[j] = 0.
	//

	// data section:
	toData := dest.seg.Data[destBase:destDataEnd]
	from := src.seg.Data[srcBase:srcDataEnd]
	copyCount := copy(toData, from)
	toData = toData[copyCount:]
	for j := range toData {
		toData[j] = 0
	}

	// ptrs section:

	// version handling: we ignore any extra-newer-pointers in src,
	// i.e. the case when srcPtrSize > dstPtrSize, by only
	// running j over the size of dstPtrSize, the destination size.
	numSrcPtrs := src.size.PointerCount
	numDstPtrs := dest.size.PointerCount
	for j := uint16(0); j < numSrcPtrs && j < numDstPtrs; j++ {
		m := src.seg.readPtr(srcDataEnd.element(int32(j), wordSize))
		err := dest.seg.writePtr(destDataEnd.element(int32(j), wordSize), m, copies, depth+1)
		if err != nil {
			return err
		}
	}
	for j := numSrcPtrs; j < numDstPtrs; j++ {
		// destination p is a newer version than source so these extra new pointer fields in p must be zeroed.
		addr := destDataEnd.element(int32(j), wordSize)
		putRawPointer(dest.seg.Data[addr:], 0)
	}
	// Nothing more here: so any other pointers in srcPtrSize beyond
	// those in dstPtrSize are ignored and discarded.

	return nil

} // end copyStructHandlingVersionSkew()

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
		boff := bitOffset(i)
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

func (s *Segment) lookupSegment(id SegmentID) (*Segment, error) {
	if s.Id != id {
		return s.Message.Lookup(id)
	} else {
		return s, nil
	}
}

func (s *Segment) readPtr(off Address) Pointer {
	var err error
	val := rawPointer(binary.LittleEndian.Uint64(s.Data[off:]))

	switch val.pointerType() {
	case doubleFarPointer:
		// A double far pointer points to a double pointer, where the
		// first points to the actual data, and the second is the tag
		// that would normally be placed right before the data (offset
		// == 0).

		faroff, segid := val.farAddress(), val.farSegment()
		if s, err = s.lookupSegment(segid); err != nil || !s.regionInBounds(faroff, wordSize.times(2)) {
			return Pointer{}
		}
		far := rawPointer(binary.LittleEndian.Uint64(s.Data[faroff:]))
		tag := rawPointer(binary.LittleEndian.Uint64(s.Data[faroff.addSize(wordSize):]))
		if far.pointerType() != farPointer || tag.offset() != 0 {
			return Pointer{}
		}
		segid = far.farSegment()
		if s, err = s.lookupSegment(segid); err != nil {
			return Pointer{}
		}

		// -8 because far pointers reference from the start of the
		// segment, but offsets reference the end of the pointer data.
		val = tag | rawPointer(far.farAddress()-Address(wordSize))<<2

	case farPointer:
		faroff, segid := val.farAddress(), val.farSegment()

		if s, err = s.lookupSegment(segid); err != nil || !s.regionInBounds(faroff, wordSize) {
			return Pointer{}
		}

		off = faroff
		val = rawPointer(binary.LittleEndian.Uint64(s.Data[faroff:]))
	}

	if val == 0 {
		// This is a null pointer.
		return Pointer{}
	}

	// Be wary of overflow. Offset is 30 bits signed. List size is 29 bits
	// unsigned. For both of these we need to check in terms of words if
	// using 32 bit maths as bits or bytes will overflow.
	switch val.pointerType() {
	case structPointer:
		addr, ok := val.offset().resolve(off)
		if !ok {
			return Pointer{}
		}
		sz := val.structSize()
		if !s.regionInBounds(addr, sz.totalSize()) {
			return Pointer{}
		}
		return Pointer{
			seg:  s,
			typ:  TypeStruct,
			off:  addr,
			size: sz,
		}
	case listPointer:
		addr, ok := val.offset().resolve(off)
		if !ok {
			return Pointer{}
		}
		lt, lsize := val.listType(), val.totalListSize()
		if !s.regionInBounds(addr, lsize) {
			return Pointer{}
		}
		if lt == compositeList {
			hdr := rawPointer(binary.LittleEndian.Uint64(s.Data[addr:]))
			addr = addr.addSize(wordSize)
			if hdr.pointerType() != structPointer {
				return Pointer{}
			}
			sz := hdr.structSize()
			n := int32(hdr.offset())
			if !s.regionInBounds(addr, sz.totalSize().times(n)) {
				return Pointer{}
			}
			return Pointer{
				seg:    s,
				typ:    TypeList,
				size:   sz,
				off:    addr,
				length: n,
				flags:  isCompositeList,
			}
		}
		return Pointer{
			seg:    s,
			typ:    listPointerType(lt),
			size:   val.elementSize(),
			off:    addr,
			length: val.numListElements(),
		}
	case otherPointer:
		if val.otherPointerType() != 0 {
			return Pointer{}
		}
		return Pointer{
			seg: s,
			typ: TypeInterface,
			cap: val.capabilityIndex(),
		}

	default:
		return Pointer{}
	}
}

func listPointerType(listType int) PointerType {
	switch listType {
	case bit1List:
		return TypeBitList
	case pointerList:
		return TypePointerList
	default:
		return TypeList
	}
}

// value converts the pointer into a raw near pointer.
// paddr is where the pointer will be located in the segment.
func (p Pointer) value(paddr Address) rawPointer {
	off := makePointerOffset(paddr, p.off)

	switch p.typ {
	case TypeStruct:
		return rawStructPointer(off, p.size)
	case TypePointerList:
		return rawListPointer(off, pointerList, p.length)
	case TypeList:
		if (p.flags & isCompositeList) != 0 {
			// p.off points to the data not the header
			return rawListPointer(off-1, compositeList, p.length*p.size.totalWordCount())
		}

		switch p.size.DataSize {
		case 0:
			return rawListPointer(off, voidList, p.length)
		case 1:
			return rawListPointer(off, byte1List, p.length)
		case 2:
			return rawListPointer(off, byte2List, p.length)
		case 4:
			return rawListPointer(off, byte4List, p.length)
		case 8:
			return rawListPointer(off, byte8List, p.length)
		default:
			panic(errListSize)
		}

	case TypeBitList:
		return rawListPointer(off, bit1List, p.length)
	case TypeInterface:
		return rawInterfacePointer(p.cap)
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

type offset struct {
	id         SegmentID
	boff, bend int64 // in bits
	newval     Pointer
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

// dataEnd returns the first address greater than p.off that is not part of p's bounds.
func (p Pointer) dataEnd() Address {
	switch p.typ {
	case TypeList:
		return p.off.addSize(p.size.totalSize().times(p.length))
	case TypePointerList:
		return p.off.addSize(wordSize.times(p.length))
	case TypeStruct:
		return p.off.addSize(p.size.totalSize())
	case TypeBitList:
		return p.off.addSize(Size((p.length + 7) / 8))
	default:
		return p.off
	}
}

func isEmptyStruct(src Pointer) bool {
	return src.typ == TypeStruct && src.length == 0 && src.size.isZero() && src.flags == 0
}

func (destSeg *Segment) writePtr(off Address, src Pointer, copies *rbtree.Tree, depth int) error {
	// handle size-zero Objects/empty structs
	if src.seg == nil {
		return nil
	}
	srcSeg := src.seg

	if src.typ == TypeNull || isEmptyStruct(src) {
		binary.LittleEndian.PutUint64(destSeg.Data[off:], 0)
		return nil

	} else if destSeg == srcSeg || src.typ == TypeInterface && destSeg.Message == srcSeg.Message {
		// Same segment
		binary.LittleEndian.PutUint64(destSeg.Data[off:], uint64(src.value(off)))
		return nil

	} else if src.typ == TypeInterface {
		// Different messages.  Need to copy table entry.
		c := destSeg.Message.AddCap(Interface(src).Client())
		p := Pointer(destSeg.NewInterface(c))
		binary.LittleEndian.PutUint64(destSeg.Data[off:], uint64(p.value(off)))
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
			newval: Pointer{
				typ:    src.typ,
				length: src.length,
				size:   src.size,
				flags:  (src.flags & isCompositeList),
			},
		}

		if (src.flags & isBitListMember) != 0 {
			key.boff += int64(src.flags & bitOffsetMask)
			key.bend = key.boff + 1
			key.newval.size.DataSize = 8
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
		n, err := destSeg.create(Size((key.bend-key.boff)/8), key.newval)
		if err != nil {
			return err
		}

		// n is possibly in a new segment, if destSeg was full.
		newSeg := n.seg

		if (n.flags & isCompositeList) != 0 {
			copy(newSeg.Data[n.off:], srcSeg.Data[src.off-8:src.off])
			n.off += 8
		}

		key.newval = n
		copies.Insert(key)

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
				dest := Pointer{
					seg:  newSeg,
					off:  n.off,
					size: n.size,
				}

				if err := copyStructHandlingVersionSkew(dest, src, copies, depth, 0, 0); err != nil {
					return err
				}
			}

		case TypeList:
			// recognize Data and Text, both List(Byte), as special cases for speed.
			if n.size.isOneByte() && src.size.isOneByte() {
				copy(newSeg.Data[n.off:], srcSeg.Data[src.off:src.off.element(n.length, 1)])
				break
			}

			dest := Pointer{
				seg:  newSeg,
				off:  n.off,
				size: n.size,
			}
			for i := int32(0); i < n.length; i++ {
				if err := copyStructHandlingVersionSkew(dest, src, copies, depth, i, i); err != nil {
					return err
				}
			}

		case TypePointerList:
			for i := int32(0); i < n.length; i++ {
				c := srcSeg.readPtr(src.off.element(i, wordSize))
				if err := newSeg.writePtr(n.off.element(i, wordSize), c, copies, depth+1); err != nil {
					return err
				}
			}

		case TypeBitList:
			// TODO(light): DataSize is not populated for BitList.
			nbytes := int(src.length + 7/8)
			listSize := Size(nbytes).padToWord()
			copy(newSeg.Data[n.off:], srcSeg.Data[src.off:src.off.addSize(listSize)])
		}
		return destSeg.writePtr(off, key.newval, nil, depth+1)

	} else if (src.flags & hasPointerTag) != 0 {
		// By lucky chance, the data has a tag in front of it. This
		// happens when create had to move the data to a new segment.
		binary.LittleEndian.PutUint64(destSeg.Data[off:], uint64(rawFarPointer(srcSeg.Id, src.off-8)))
		return nil

	} else if len(srcSeg.Data)+int(wordSize) <= cap(srcSeg.Data) {
		// Have room in the target for a tag
		srcAddr := Address(len(srcSeg.Data))
		srcSeg.Data = srcSeg.Data[:srcAddr.addSize(wordSize)]
		binary.LittleEndian.PutUint64(srcSeg.Data[srcAddr:], uint64(src.value(srcAddr)))
		binary.LittleEndian.PutUint64(destSeg.Data[off:], uint64(rawFarPointer(srcSeg.Id, srcAddr)))
		return nil

	} else {
		// Need to create a double far pointer. Try and create it in
		// the originating segment if we can.
		t := destSeg
		const landingSize = wordSize * 2
		if len(t.Data)+int(landingSize) > cap(t.Data) {
			var err error
			if t, err = t.Message.NewSegment(landingSize); err != nil {
				return err
			}
		}

		dstAddr := Address(len(t.Data))
		binary.LittleEndian.PutUint64(t.Data[dstAddr:], uint64(rawFarPointer(srcSeg.Id, src.off)))
		binary.LittleEndian.PutUint64(t.Data[dstAddr.addSize(wordSize):], uint64(src.value(src.off-8)))
		binary.LittleEndian.PutUint64(destSeg.Data[off:], uint64(rawDoubleFarPointer(t.Id, dstAddr)))
		t.Data = t.Data[:len(t.Data)+16]
		return nil
	}
}

// An Interface is a wrapper for a Pointer that provides methods to
// access interface information.
type Interface Pointer

// Segment returns the segment this pointer came from.
func (i Interface) Segment() *Segment {
	return i.seg
}

// Capability returns the capability ID of the interface.
func (i Interface) Capability() CapabilityID {
	return i.cap
}

// Client returns the client stored in the message's capability table
// or nil if the pointer is invalid.
func (i Interface) Client() Client {
	if i.seg == nil || i.typ == TypeNull {
		return nil
	}
	tab := i.seg.Message.CapTable()
	if uint64(i.cap) >= uint64(len(tab)) {
		return nil
	}
	return tab[i.cap]
}
