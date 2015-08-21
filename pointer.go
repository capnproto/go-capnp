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

// PointerDefault returns p if it is valid, otherwise it unmarshals def.
func PointerDefault(p Pointer, def []byte) (Pointer, error) {
	if !IsValid(p) {
		return unmarshalDefault(def)
	}
	return p, nil
}

func unmarshalDefault(def []byte) (Pointer, error) {
	msg, err := Unmarshal(def)
	if err != nil {
		return nil, err
	}
	p, err := msg.Root()
	if err != nil {
		return nil, err
	}
	return p, nil
}

// pointerAddress returns the pointer's address.
// It panics if p's underlying pointer is not a valid Struct or List.
func pointerAddress(p Pointer) Address {
	type addresser interface {
		Address() Address
	}
	a := p.underlying().(addresser)
	return a.Address()
}
