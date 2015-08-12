package capnp

// wordSize is the number of bytes in a Cap'n Proto word.
const wordSize Size = 8

// An Address is an index inside a segment's data (in bytes).
type Address uint32

// addSize returns the address a+sz.
func (a Address) addSize(sz Size) Address {
	return a.element(1, sz)
}

// element returns the address a+i*sz.
func (a Address) element(i int32, sz Size) Address {
	return a + Address(sz.times(i))
}

// addOffset returns the address a+o.
func (a Address) addOffset(o DataOffset) Address {
	return a + Address(o)
}

// A Size is a size (in bytes).
type Size uint32

// DataOffset is an offset in bytes from the beginning of a struct's data section.
type DataOffset uint32

// times returns the size sz*n.
func (sz Size) times(n int32) Size {
	const maxSize = 1<<32 - 1
	result := int64(sz) * int64(n)
	if result > maxSize {
		panic(ErrOverlarge)
	}
	return Size(result)
}

// padToWord adds padding to sz to make it divisible by wordSize.
func (sz Size) padToWord() Size {
	n := Size(wordSize - 1)
	return (sz + n) &^ n
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

type bitOffset uint32

func (boff bitOffset) offset() DataOffset {
	return DataOffset(boff / 8)
}

func (boff bitOffset) mask() byte {
	return byte(1 << (boff % 8))
}
