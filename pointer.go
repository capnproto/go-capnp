package capnp

import (
	"strconv"
)

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

// Address returns the address the pointer references.
func (p Pointer) Address() Address {
	return p.off
}

// IsNull reports whether this is a null pointer.
func (p Pointer) IsNull() bool {
	return p.typ == TypeNull
}

// Segment returns the segment this pointer came from.
func (p Pointer) Segment() *Segment {
	return p.seg
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

func (p Pointer) ToObjectDefault(s *Segment, tagAddr Address) Pointer {
	if p.typ == TypeNull {
		return s.Root(tagAddr)
	} else {
		return p
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

// objectEnd returns the first address greater than p.off that is not
// part of the object's bounds.
func (p Pointer) objectEnd() Address {
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
