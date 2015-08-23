// +build !nocapnpstrings

package capnp

import (
	"fmt"
)

// String returns the address in hex format.
func (addr Address) String() string {
	return fmt.Sprintf("%#016x", uint64(addr))
}

// GoString returns the address in hex format.
func (addr Address) GoString() string {
	return fmt.Sprintf("capnp.Address(%#016x)", uint64(addr))
}

// String returns the size in the format "X bytes".
func (sz Size) String() string {
	if sz == 1 {
		return "1 byte"
	}
	return fmt.Sprintf("%d bytes", sz)
}

// GoString returns the size as a Go expression.
func (sz Size) GoString() string {
	return fmt.Sprintf("capnp.Size(%d)", sz)
}

// String returns the offset in the format "+X bytes".
func (off DataOffset) String() string {
	if off == 1 {
		return "+1 byte"
	}
	return fmt.Sprintf("+%d bytes", off)
}

// GoString returns the offset as a Go expression.
func (off DataOffset) GoString() string {
	return fmt.Sprintf("capnp.DataOffset(%d)", off)
}

// String returns a short, human readable representation of the object
// size.
func (sz ObjectSize) String() string {
	return fmt.Sprintf("{datasz=%d ptrs=%d}", sz.DataSize, sz.PointerCount)
}

// GoString formats the ObjectSize as a keyed struct literal.
func (sz ObjectSize) GoString() string {
	return fmt.Sprintf("capnp.ObjectSize{DataSize: %d, PointerCount: %d}", sz.DataSize, sz.PointerCount)
}

// String returns the offset in the format "bit X".
func (bit BitOffset) String() string {
	return fmt.Sprintf("bit %d", bit)
}

// GoString returns the offset as a Go expression.
func (bit BitOffset) GoString() string {
	return fmt.Sprintf("capnp.BitOffset(%d)", bit)
}

// GoString formats the pointer as a call to one of the rawPointer
// construction functions.
func (p rawPointer) GoString() string {
	if p == 0 {
		return "rawPointer(0)"
	}
	switch p.pointerType() {
	case structPointer:
		return fmt.Sprintf("rawStructPointer(%d, %#v)", p.offset(), p.structSize())
	case listPointer:
		var lt string
		switch p.listType() {
		case voidList:
			lt = "voidList"
		case bit1List:
			lt = "bit1List"
		case byte1List:
			lt = "byte1List"
		case byte2List:
			lt = "byte2List"
		case byte4List:
			lt = "byte4List"
		case byte8List:
			lt = "byte8List"
		case pointerList:
			lt = "pointerList"
		case compositeList:
			lt = "compositeList"
		}
		return fmt.Sprintf("rawListPointer(%d, %s, %d)", p.offset(), lt, p.numListElements())
	case farPointer:
		return fmt.Sprintf("rawFarPointer(%d, %v)", p.farSegment(), p.farAddress())
	case doubleFarPointer:
		return fmt.Sprintf("rawDoubleFarPointer(%d, %v)", p.farSegment(), p.farAddress())
	default:
		// other pointer
		if p.otherPointerType() != 0 {
			return fmt.Sprintf("rawPointer(%#016x)", uint64(p))
		}
		return fmt.Sprintf("rawInterfacePointer(%d)", p.capabilityIndex())
	}
}
