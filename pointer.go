package capnp

// A Ptr is a reference to a Cap'n Proto struct, list, or interface.
type Ptr struct {
	Struct    Struct
	List      List
	Interface Interface
}

func toPtr(p Pointer) Ptr {
	if p == nil {
		return Ptr{}
	}
	switch p := p.underlying().(type) {
	case Struct:
		return Ptr{Struct: p}
	case List:
		return Ptr{List: p}
	case Interface:
		return Ptr{Interface: p}
	}
	return Ptr{}
}

func (p Ptr) toPointer() Pointer {
	switch {
	case p.Struct.seg != nil:
		return p.Struct
	case p.List.seg != nil:
		return p.List
	case p.Interface.seg != nil:
		return p.Interface
	}
	return nil
}

// IsValid reports whether p is valid.
func (p Ptr) IsValid() bool {
	return p.Struct.seg != nil || p.List.seg != nil || p.Interface.seg != nil
}

func (p Ptr) Segment() *Segment {
	switch {
	case p.Struct.seg != nil:
		return p.Struct.seg
	case p.List.seg != nil:
		return p.List.seg
	case p.Interface.seg != nil:
		return p.Interface.seg
	}
	return nil
}

// HasData returns true if the pointer is valid and has non-zero size.
func (p Ptr) HasData() bool {
	return p.Struct.HasData() || p.List.HasData() || p.Interface.HasData()
}

func (p Ptr) value(paddr Address) rawPointer {
	switch {
	case p.Struct.seg != nil:
		return p.Struct.value(paddr)
	case p.List.seg != nil:
		return p.List.value(paddr)
	case p.Interface.seg != nil:
		return p.Interface.value(paddr)
	}
	return 0
}

// address returns the pointer's address.  It panics if p is not a valid Struct or List.
func (p Ptr) address() Address {
	switch {
	case p.Struct.seg != nil:
		return p.Struct.Address()
	case p.List.seg != nil:
		return p.List.Address()
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
