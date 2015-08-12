package capnp

import (
	"encoding/binary"
)

type Struct Pointer

// Segment returns the segment this pointer came from.
func (s Struct) Segment() *Segment {
	return s.seg
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
