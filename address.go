package capnp

import (
	"fmt"
)

// An address is an index inside a segment's data (in bytes).
type address uint32

// String returns the address in hex format.
func (addr address) String() string {
	return fmt.Sprintf("%#08x", uint64(addr))
}

// GoString returns the address in hex format.
func (addr address) GoString() string {
	return fmt.Sprintf("capnp.address(%#08x)", uint64(addr))
}

// addSize returns the address a+sz.
func (a address) addSize(sz Size) (b address, ok bool) {
	x := int64(a) + int64(sz)
	if x > int64(maxSize) {
		return 0, false
	}
	return address(x), true
}

// element returns the address a+i*sz.
func (a address) element(i int32, sz Size) (b address, ok bool) {
	x := int64(i) * int64(sz)
	if x > int64(maxSize) {
		return 0, false
	}
	x += int64(a)
	if x > int64(maxSize) {
		return 0, false
	}
	return address(x), true
}

// addOffset returns the address a+o.
func (a address) addOffset(o DataOffset) address {
	return a + address(o)
}

// A Size is a size (in bytes).
type Size uint32

// wordSize is the number of bytes in a Cap'n Proto word.
const wordSize Size = 8

// maxSize is the maximum representable size.
const maxSize Size = 1<<32 - 1

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

// times returns the size sz*n.
func (sz Size) times(n int32) (ns Size, ok bool) {
	x := int64(sz) * int64(n)
	if x > int64(maxSize) {
		return 0, false
	}
	return Size(x), true
}

// padToWord adds padding to sz to make it divisible by wordSize.
func (sz Size) padToWord() Size {
	n := Size(wordSize - 1)
	return (sz + n) &^ n
}

// DataOffset is an offset in bytes from the beginning of a struct's data section.
type DataOffset uint32

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
	// Guaranteed not to overflow
	return wordSize * Size(sz.PointerCount)
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

// String returns a short, human readable representation of the object
// size.
func (sz ObjectSize) String() string {
	return fmt.Sprintf("{datasz=%d ptrs=%d}", sz.DataSize, sz.PointerCount)
}

// GoString formats the ObjectSize as a keyed struct literal.
func (sz ObjectSize) GoString() string {
	return fmt.Sprintf("capnp.ObjectSize{DataSize: %d, PointerCount: %d}", sz.DataSize, sz.PointerCount)
}

// BitOffset is an offset in bits from the beginning of a struct's data section.
type BitOffset uint32

// offset returns the equivalent byte offset.
func (bit BitOffset) offset() DataOffset {
	return DataOffset(bit / 8)
}

// mask returns the bitmask for the bit.
func (bit BitOffset) mask() byte {
	return byte(1 << (bit % 8))
}

// String returns the offset in the format "bit X".
func (bit BitOffset) String() string {
	return fmt.Sprintf("bit %d", bit)
}

// GoString returns the offset as a Go expression.
func (bit BitOffset) GoString() string {
	return fmt.Sprintf("capnp.BitOffset(%d)", bit)
}
