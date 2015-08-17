package capnp

// A Pointer is a reference to a Cap'n Proto struct, list, or interface.
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

// IsValid reports whether p is non-nil and valid.
func IsValid(p Pointer) bool {
	return p != nil && p.Segment() != nil
}

// HasData returns true if the pointer is valid and has non-zero size.
func HasData(p Pointer) bool {
	return IsValid(p) && p.HasData()
}

/*
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
*/

// pointerAddress returns the pointer's address.
// It panics if p's underlying pointer is not a valid Struct or List.
func pointerAddress(p Pointer) Address {
	type addresser interface {
		Address() Address
	}
	a := p.underlying().(addresser)
	return a.Address()
}
