package capnproto

import (
	"encoding/binary"
	"errors"
	"math"
)

var (
	little16    = binary.LittleEndian.Uint16
	little32    = binary.LittleEndian.Uint32
	little64    = binary.LittleEndian.Uint64
	putLittle16 = binary.LittleEndian.PutUint16
	putLittle32 = binary.LittleEndian.PutUint32
	putLittle64 = binary.LittleEndian.PutUint64

	ErrInvalidInterface = errors.New("capnproto: invalid interface")
)

type NewFunc func(p PointerType) (Pointer, error)
type CallFunc func(p Pointer, method int, args Pointer) (Pointer, error)

type Pointer interface {
	// New allocates a compatible pointer (same or referencible segment)
	New(p PointerType) (Pointer, error)
	Call(method int, args Pointer) (Pointer, error)

	Type() PointerType

	// off specifies offset in bytes
	Read(off int, v []uint8) error
	Write(off int, v []uint8) error

	// off specifies offset in words/index
	ReadPtrs(off int, p []Pointer) error
	WritePtrs(off int, p []Pointer) error
}

type Marshaller interface {
	MarshalCaptain(new NewFunc) (Pointer, error)
}

/*
Struct pointer layout:

lsb                      struct pointer                       msb
+-+-----------------------------+---------------+---------------+
|A|             B               |       C       |       D       |
+-+-----------------------------+---------------+---------------+

A (2 bits) = 0, to indicate that this is a struct pointer.
B (30 bits) = Offset, in words, from the start of the pointer to the
    start of the struct's data section.  Signed.
C (16 bits) = Size of the struct's data section, in words.
D (16 bits) = Size of the struct's pointer section, in words.

List pointer layout:

lsb                       list pointer                        msb
+-+-----------------------------+--+----------------------------+
|A|             B               |C |             D              |
+-+-----------------------------+--+----------------------------+

A (2 bits) = 1, to indicate that this is a list pointer.
B (30 bits) = Offset, in words, from the start of the pointer to the
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
D (29 bits) = Number of elements in the list (C!=7). Size of list data in
    words minus the header (C==7).


Far pointer layout:

lsb                        far pointer                        msb
+-+-----------------------------+-------------------------------+
|A|             B               |               C               |
+-+-----------------------------+-------------------------------+

A (2 bits) = 2, to indicate that this is a far pointer.
B (30 bits) = Offset, in words, from the start of the target segment
    to the location of the far-pointer landing-pad within that
    segment.
C (32 bits) = ID of the target segment.  (Segments are numbered
    sequentially starting from zero.)

Composite header layout

lsb                      struct pointer                       msb
+-+-----------------------------+---------------+---------------+
|A|             B               |       C       |       D       |
+-+-----------------------------+---------------+---------------+

A (2 bits) = 0, to indicate that this is a struct pointer.
B (30 bits) = number of elements in the list
C (16 bits) = Size of the struct's data section, in words.
D (16 bits) = Size of the struct's pointer section, in words.
*/

type PointerType struct {
	Value     uint64
	Composite uint64
}

type DataType int

const (
	Struct        DataType = 0
	ListPointer            = 1
	FarPointer             = 2
	VoidList               = 8
	Bit1List               = 9
	Byte1List              = 10
	Byte2List              = 11
	Byte4List              = 12
	Byte8List              = 13
	PointerList            = 14
	CompositeList          = 15
)

func NewStruct(new NewFunc, data, ptrs int) (Pointer, error) {
	return new(PointerType{uint64(Struct) | uint64((data+7)/8)<<32 | uint64(ptrs)<<48, 0})
}

func NewList(new NewFunc, typ DataType, elems int) (Pointer, error) {
	return new(PointerType{uint64(ListPointer) | uint64(typ)<<31 | uint64(elems)<<35, 0})
}

func NewCompositeList(new NewFunc, elems, data, ptrs int) (Pointer, error) {
	data = (data + 7) / 8
	return new(PointerType{
		uint64(ListPointer) | uint64(CompositeList)<<31 | uint64(elems*(data+ptrs))<<35,
		uint64(Struct) | uint64(elems)<<2 | uint64(data)<<32 | uint64(ptrs)<<48,
	})
}

func MakeFarPointer(seg uint32, off int) uint64 {
	return uint64(FarPointer) | uint64(off)<<2 | uint64(seg)<<32
}

func (p PointerType) CompositeType() PointerType {
	return PointerType{p.Composite, 0}
}

func (p PointerType) Type() DataType {
	if p.Value&7 == ListPointer {
		return DataType((p.Value>>32)&7) + 8
	} else {
		return DataType(p.Value & 7)
	}
}

func (p PointerType) SegmentId() uint32 {
	return uint32(p.Value >> 32)
}

func (p *PointerType) SetOffset(off int) {
	p.Value &= 0xFFFFFFFF00000003
	p.Value |= uint64(uint32(int32(off << 2)))
}

func (p PointerType) Offset() int {
	return int(int32(uint32(p.Value))) >> 2
}

func (p PointerType) Elements() int {
	switch p.Type() {
	case VoidList, Bit1List, Byte1List, Byte2List, Byte4List, Byte8List, PointerList:
		return int(p.Value >> 35)
	case CompositeList:
		return int(uint32(p.Composite) >> 2)
	default:
		return 0
	}
}

func (p PointerType) DataSize() int {
	switch p.Type() {
	case Struct:
		return int(uint16(p.Value>>32)) * 8
	case Bit1List:
		return (int(p.Value>>35) + 7) / 8
	case Byte1List:
		return int(p.Value >> 35)
	case Byte2List:
		return int(p.Value>>35) * 2
	case Byte4List:
		return int(p.Value>>35) * 4
	case Byte8List:
		return int(p.Value>>35) * 8
	case CompositeList:
		return int(uint16(p.Composite>>32)) * 8
	default:
		return 0
	}
}

func (p PointerType) PointerNum() int {
	switch p.Type() {
	case Struct:
		return int(uint16(p.Value >> 48))
	case PointerList:
		return int(p.Value >> 35)
	case CompositeList:
		return int(uint16(p.Composite >> 48))
	default:
		return 0
	}
}

func ReadPtr(ptr Pointer, off int) Pointer {
	ret := [1]Pointer{}
	ptr.ReadPtrs(off, ret[:])
	return ret[0]
}

func ReadUInt8(ptr Pointer, off int) uint8 {
	// ignore read error, buf is still zero
	buf := [1]byte{}
	ptr.Read(off, buf[:])
	return buf[0]
}

func ReadUInt16(ptr Pointer, off int) uint16 {
	buf := [2]byte{}
	if err := ptr.Read(off, buf[:]); err != nil {
		return 0
	}
	return little16(buf[:])
}

func ReadUInt32(ptr Pointer, off int) uint32 {
	buf := [4]byte{}
	if err := ptr.Read(off, buf[:]); err != nil {
		return 0
	}
	return little32(buf[:])
}

func ReadUInt64(ptr Pointer, off int) uint64 {
	buf := [8]byte{}
	if err := ptr.Read(off, buf[:]); err != nil {
		return 0
	}
	return little64(buf[:])
}

func ToPointerList(p Pointer) []Pointer {
	typ := p.Type()
	if typ.Type() != PointerList && typ.Type() != CompositeList {
		return nil
	}
	ret := make([]Pointer, typ.Elements())
	if err := p.ReadPtrs(0, ret); err != nil {
		return nil
	}
	return ret
}

func tolist(p Pointer, expect DataType, esz int) []uint8 {
	typ := p.Type()

	switch typ.Type() {
	case expect:
		u8 := make([]byte, typ.DataSize())
		if err := p.Read(0, u8); err != nil {
			return nil
		}
		return u8

	case PointerList, CompositeList:
		// List(T) has been converted into a List(struct). We handle
		// it as if the first entry of each struct is an item in the
		// list.
		u8 := make([]byte, typ.Elements()*esz)
		ptrs := make([]Pointer, typ.Elements())
		if err := p.ReadPtrs(0, ptrs); err != nil {
			return nil
		}
		for i, ptr := range ptrs {
			if err := ptr.Read(0, u8[i*esz:]); err != nil {
				return nil
			}
		}
		return u8

	default:
		return nil
	}
}

func ToVoidList(p Pointer) []struct{} {
	typ := p.Type()

	switch typ.Type() {
	case VoidList, PointerList, CompositeList:
		return make([]struct{}, typ.Elements())
	default:
		return nil
	}
}

// Litle endian bits
type Bitset []byte

func (b Bitset) Test(i int) bool {
	return (b[i/8] & (1 << uint(i%8))) != 0
}

func (b *Bitset) Set(i int) {
	(*b)[i/8] |= 1 << uint(i%8)
}

func (b *Bitset) Clear(i int) {
	(*b)[i/8] &^= 1 << uint(i%8)
}

func ToBitset(p Pointer) Bitset {
	typ := p.Type()

	switch typ.Type() {
	case Bit1List:
		u8 := make([]byte, typ.DataSize())
		if err := p.Read(0, u8); err != nil {
			return nil
		}
		return Bitset(u8)

	case PointerList, CompositeList:
		// List(T) has been converted into a List(struct). We handle
		// it as if the first entry of each struct is an item in the
		// list.
		bits := make(Bitset, (typ.Elements()+7)/8)
		ptrs := make([]Pointer, typ.Elements())
		if err := p.ReadPtrs(0, ptrs); err != nil {
			return nil
		}
		for i, ptr := range ptrs {
			u8 := [1]byte{}
			if err := ptr.Read(0, u8[:]); err != nil {
				return nil
			}
			if (u8[0] & 1) != 0 {
				bits.Set(i)
			}
		}
		return bits

	default:
		return nil
	}
}

func ToUInt8List(p Pointer) []uint8 {
	return tolist(p, Byte1List, 1)
}

func ToString(p Pointer, def string) string {
	u8 := ToUInt8List(p)
	if u8 == nil || len(u8) < 1 || u8[len(u8)-1] != 0 {
		return def
	}
	return string(u8[:len(u8)-1])
}

func ToInt8List(p Pointer) []int8 {
	u8 := tolist(p, Byte1List, 1)
	if u8 == nil {
		return nil
	}
	ret := make([]int8, len(u8))
	for i, u := range u8 {
		ret[i] = int8(u)
	}
	return ret
}

func ToUInt16List(p Pointer) []uint16 {
	u8 := tolist(p, Byte2List, 2)
	if u8 == nil {
		return nil
	}
	ret := make([]uint16, len(u8)/2)
	for i := range ret {
		ret[i] = little16(u8[2*i:])
	}
	return ret
}

func ToInt16List(p Pointer) []int16 {
	u8 := tolist(p, Byte2List, 2)
	if u8 == nil {
		return nil
	}
	ret := make([]int16, len(u8)/2)
	for i := range ret {
		ret[i] = int16(little16(u8[2*i:]))
	}
	return ret
}

func ToUInt32List(p Pointer) []uint32 {
	u8 := tolist(p, Byte4List, 4)
	if u8 == nil {
		return nil
	}
	ret := make([]uint32, len(u8)/4)
	for i := range ret {
		ret[i] = little32(u8[4*i:])
	}
	return ret
}

func ToInt32List(p Pointer) []int32 {
	u8 := tolist(p, Byte4List, 4)
	if u8 == nil {
		return nil
	}
	ret := make([]int32, len(u8)/4)
	for i := range ret {
		ret[i] = int32(little32(u8[4*i:]))
	}
	return ret
}

func ToUInt64List(p Pointer) []uint64 {
	u8 := tolist(p, Byte8List, 8)
	if u8 == nil {
		return nil
	}
	ret := make([]uint64, len(u8)/8)
	for i := range ret {
		ret[i] = little64(u8[8*i:])
	}
	return ret
}

func ToInt64List(p Pointer) []int64 {
	u8 := tolist(p, Byte8List, 8)
	if u8 == nil {
		return nil
	}
	ret := make([]int64, len(u8)/8)
	for i := range ret {
		ret[i] = int64(little16(u8[8*i:]))
	}
	return ret
}

func ToFloat32List(p Pointer) []float32 {
	u8 := tolist(p, Byte4List, 4)
	if u8 == nil {
		return nil
	}
	ret := make([]float32, len(u8)/4)
	for i := range ret {
		ret[i] = math.Float32frombits(little32(u8[4*i:]))
	}
	return ret
}

func ToFloat64List(p Pointer) []float64 {
	u8 := tolist(p, Byte8List, 8)
	if u8 == nil {
		return nil
	}
	ret := make([]float64, len(u8)/8)
	for i := range ret {
		ret[i] = math.Float64frombits(little64(u8[8*i:]))
	}
	return ret
}

func ToStringList(p Pointer) []string {
	typ := p.Type()
	if typ.Type() != PointerList && typ.Type() != CompositeList {
		return nil
	}
	ret := make([]string, typ.Elements())
	for i := range ret {
		ret[i] = ToString(ReadPtr(p, i), "")
	}
	return ret
}

func ToBitsetList(p Pointer) []Bitset {
	typ := p.Type()
	if typ.Type() != PointerList && typ.Type() != CompositeList {
		return nil
	}
	ret := make([]Bitset, typ.Elements())
	for i := range ret {
		ret[i] = ToBitset(ReadPtr(p, i))
	}
	return ret
}

func WriteBool(ptr Pointer, off int, v bool) error {
	u := ReadUInt8(ptr, off/8)
	if v {
		u &= 1 << uint(off%8)
	} else {
		u &^= 1 << uint(off%8)
	}
	return WriteUInt8(ptr, off/8, u)
}

func WriteUInt8(ptr Pointer, off int, v uint8) error {
	return ptr.Write(off, []uint8{v})
}

func WriteUInt16(ptr Pointer, off int, v uint16) error {
	u8 := [2]uint8{}
	putLittle16(u8[:], v)
	return ptr.Write(off, u8[:])
}

func WriteUInt32(ptr Pointer, off int, v uint32) error {
	u8 := [4]uint8{}
	putLittle32(u8[:], v)
	return ptr.Write(off, u8[:])
}

func WriteUInt64(ptr Pointer, off int, v uint64) error {
	u8 := [8]uint8{}
	putLittle64(u8[:], v)
	return ptr.Write(off, u8[:])
}

func newList(new NewFunc, typ DataType, sz int, data []uint8) (Pointer, error) {
	to, err := NewList(new, typ, sz)
	if err != nil {
		return nil, err
	}
	if err := to.Write(0, data); err != nil {
		return nil, err
	}
	return to, nil
}

func NewString(new NewFunc, v string) (Pointer, error) {
	return newList(new, Byte1List, len(v)+1, []byte(v))
}

func NewBitset(new NewFunc, v Bitset) (Pointer, error) {
	return newList(new, Bit1List, len(v)*8, []byte(v))
}

func NewUInt8List(new NewFunc, v []uint8) (Pointer, error) {
	return newList(new, Byte1List, len(v), v)
}

func NewInt8List(new NewFunc, v []int8) (Pointer, error) {
	u8 := make([]uint8, len(v))
	for i, u := range v {
		u8[i] = uint8(u)
	}
	return newList(new, Byte2List, len(v), u8)
}

func NewUInt16List(new NewFunc, v []uint16) (Pointer, error) {
	u8 := make([]uint8, len(v)*2)
	for i, u := range v {
		putLittle16(u8[i*2:], u)
	}
	return newList(new, Byte2List, len(v), u8)
}

func NewInt16List(new NewFunc, v []int16) (Pointer, error) {
	u8 := make([]uint8, len(v)*2)
	for i, u := range v {
		putLittle16(u8[i*2:], uint16(u))
	}
	return newList(new, Byte2List, len(v), u8)
}

func NewUInt32List(new NewFunc, v []uint32) (Pointer, error) {
	u8 := make([]uint8, len(v)*4)
	for i, u := range v {
		putLittle32(u8[i*4:], u)
	}
	return newList(new, Byte4List, len(v), u8)
}

func NewInt32List(new NewFunc, v []int32) (Pointer, error) {
	u8 := make([]uint8, len(v)*4)
	for i, u := range v {
		putLittle32(u8[i*4:], uint32(u))
	}
	return newList(new, Byte4List, len(v), u8)
}

func NewUInt64List(new NewFunc, v []uint64) (Pointer, error) {
	u8 := make([]uint8, len(v)*8)
	for i, u := range v {
		putLittle64(u8[i*8:], u)
	}
	return newList(new, Byte8List, len(v), u8)
}

func NewInt64List(new NewFunc, v []int64) (Pointer, error) {
	u8 := make([]uint8, len(v)*8)
	for i, u := range v {
		putLittle64(u8[i*8:], uint64(u))
	}
	return newList(new, Byte8List, len(v), u8)
}

func NewFloat32List(new NewFunc, v []float32) (Pointer, error) {
	u8 := make([]uint8, len(v)*4)
	for i, u := range v {
		putLittle32(u8[i*4:], math.Float32bits(u))
	}
	return newList(new, Byte4List, len(v), u8)
}

func NewFloat64List(new NewFunc, v []float64) (Pointer, error) {
	u8 := make([]uint8, len(v)*8)
	for i, u := range v {
		putLittle64(u8[i*8:], math.Float64bits(u))
	}
	return newList(new, Byte8List, len(v), u8)
}

func NewVoidList(new NewFunc, v []struct{}) (Pointer, error) {
	if v == nil {
		return nil, nil
	}
	return NewList(new, VoidList, len(v))
}

func NewPointerList(new NewFunc, v []Pointer) (Pointer, error) {
	if v == nil {
		return nil, nil
	}
	to, err := NewList(new, PointerList, len(v))
	if err != nil {
		return nil, err
	}
	if err := to.WritePtrs(0, v); err != nil {
		return nil, err
	}
	return to, nil
}

func NewStringList(new NewFunc, v []string) (Pointer, error) {
	if v == nil {
		return nil, nil
	}
	ptrs := make([]Pointer, len(v))
	for i, s := range v {
		ps, err := NewString(new, s)
		if err != nil {
			return nil, err
		}
		ptrs[i] = ps
	}

	return NewPointerList(new, ptrs)
}

func NewBitsetList(new NewFunc, v []Bitset) (Pointer, error) {
	if v == nil {
		return nil, nil
	}
	ptrs := make([]Pointer, len(v))
	for i, b := range v {
		pb, err := NewBitset(new, b)
		if err != nil {
			return nil, err
		}
		ptrs[i] = pb
	}

	return NewPointerList(new, ptrs)
}
