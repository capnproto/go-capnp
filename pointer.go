package capnp

// A Ptr is a reference to a Cap'n Proto struct, list, or interface.
type Ptr struct {
	seg      *Segment
	off      Address
	lenOrCap uint32
	size     ObjectSize
	flags    ptrFlags
}

func toPtr(p Pointer) Ptr {
	if p == nil {
		return Ptr{}
	}
	switch p := p.underlying().(type) {
	case Struct:
		return p.ToPtr()
	case List:
		return p.ToPtr()
	case Interface:
		return p.ToPtr()
	}
	return Ptr{}
}

// Struct converts p to a Struct. If p does not hold a Struct pointer,
// the zero value is returned.
func (p Ptr) Struct() Struct {
	if p.flags.ptrType() != structPtrType {
		return Struct{}
	}
	return Struct{
		seg:   p.seg,
		off:   p.off,
		size:  p.size,
		flags: p.flags.structFlags(),
	}
}

func (p Ptr) List() List {
	if p.flags.ptrType() != listPtrType {
		return List{}
	}
	return List{
		seg:    p.seg,
		off:    p.off,
		length: int32(p.lenOrCap),
		size:   p.size,
		flags:  p.flags.listFlags(),
	}
}

func (p Ptr) Interface() Interface {
	if p.flags.ptrType() != interfacePtrType {
		return Interface{}
	}
	return Interface{
		seg: p.seg,
		cap: CapabilityID(p.lenOrCap),
	}
}

func (p Ptr) toPointer() Pointer {
	if p.seg == nil {
		return nil
	}
	switch p.flags.ptrType() {
	case structPtrType:
		return p.Struct()
	case listPtrType:
		return p.List()
	case interfacePtrType:
		return p.Interface()
	}
	return nil
}

// IsValid reports whether p is valid.
func (p Ptr) IsValid() bool {
	return p.seg != nil
}

func (p Ptr) Segment() *Segment {
	return p.seg
}

func (p Ptr) value(paddr Address) rawPointer {
	switch p.flags.ptrType() {
	case structPtrType:
		return p.Struct().value(paddr)
	case listPtrType:
		return p.List().value(paddr)
	case interfacePtrType:
		return p.Interface().value(paddr)
	}
	return 0
}

// address returns the pointer's address.  It panics if p is not a valid Struct or List.
func (p Ptr) address() Address {
	switch p.flags.ptrType() {
	case structPtrType:
		return p.Struct().Address()
	case listPtrType:
		return p.List().Address()
	}
	panic("ptr not a valid struct or list")
}

// Pointer is deprecated in favor of Ptr.
type Pointer interface {
	// Segment returns the segment this pointer points into.
	// If nil, then this is an invalid pointer.
	Segment() *Segment

	// HasData reports whether the object referenced by the pointer has
	// non-zero size.
	HasData() bool

	// value converts the pointer into a raw value.
	value(paddr Address) rawPointer

	// underlying returns a Pointer that is one of a Struct, a List, or an
	// Interface.
	underlying() Pointer
}

// IsValid is deprecated.
func IsValid(p Pointer) bool {
	return p != nil && p.Segment() != nil
}

// HasData is deprecated.
func HasData(p Pointer) bool {
	return IsValid(p) && p.HasData()
}

// PointerDefault is deprecated in favor of PtrDefault.
func PointerDefault(p Pointer, def []byte) (Pointer, error) {
	pp, err := PtrDefault(toPtr(p), def)
	return pp.toPointer(), err
}

// PtrDefault returns p if it is valid, otherwise it unmarshals def.
func PtrDefault(p Ptr, def []byte) (Ptr, error) {
	if !p.IsValid() {
		return unmarshalDefault(def)
	}
	return p, nil
}

func unmarshalDefault(def []byte) (Ptr, error) {
	msg, err := Unmarshal(def)
	if err != nil {
		return Ptr{}, err
	}
	p, err := msg.RootPtr()
	if err != nil {
		return Ptr{}, err
	}
	return p, nil
}

type ptrFlags uint8

const interfacePtrFlag ptrFlags = interfacePtrType << 6

func structPtrFlag(f structFlags) ptrFlags {
	return structPtrType<<6 | ptrFlags(f)&ptrLowerMask
}

func listPtrFlag(f listFlags) ptrFlags {
	return listPtrType<<6 | ptrFlags(f)&ptrLowerMask
}

const (
	structPtrType = iota
	listPtrType
	interfacePtrType
)

func (f ptrFlags) ptrType() int {
	return int(f >> 6)
}

const ptrLowerMask ptrFlags = 0x3f

func (f ptrFlags) listFlags() listFlags {
	return listFlags(f & ptrLowerMask)
}

func (f ptrFlags) structFlags() structFlags {
	return structFlags(f & ptrLowerMask)
}
