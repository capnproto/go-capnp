package capnp

import (
	"encoding/binary"
	"errors"
	"math"

	"github.com/glycerine/rbtree"
)

var (
	errOverlarge   = errors.New("capn: overlarge struct/list")
	errOutOfBounds = errors.New("capn: write out of bounds")
	errCopyDepth   = errors.New("capn: copy depth too large")
	errOverlap     = errors.New("capn: overlapping data on copy")
	errListSize    = errors.New("capn: invalid list size")
	errObjectType  = errors.New("capn: invalid object type")
)

type Message interface {
	NewSegment(minsz Size) (*Segment, error)
	Lookup(segid SegmentID) (*Segment, error)

	CapTable() []Client
	AddCap(c Client) CapabilityID
}

// A SegmentID is a numeric identifier for a Segment.
type SegmentID uint32

// A Segment is an allocation arena for Cap'n Proto objects.
// It is part of a Message, which can contain other segments that
// reference each other.
type Segment struct {
	Message  Message
	Data     []byte
	Id       SegmentID
	RootDone bool
}

func (s *Segment) inBounds(addr Address) bool {
	return addr < Address(len(s.Data))
}

func (s *Segment) regionInBounds(base Address, sz Size) bool {
	return base.addSize(sz) <= Address(len(s.Data))
}

func (s *Segment) Root(addr Address) Pointer {
	if !s.regionInBounds(addr, wordSize) {
		return Pointer{}
	}
	return s.readPtr(addr)
}

func (s *Segment) NewRoot() (PointerList, error) {
	n, err := s.create(wordSize, Pointer{
		typ:    TypePointerList,
		length: 1,
		size:   ObjectSize{PointerCount: 1},
		flags:  isRoot,
	})
	return PointerList(n), err
}

const maxDataSize = 1<<32 - 1

func (s *Segment) NewText(v string) Pointer {
	if int64(len(v))+1 > maxDataSize {
		panic(errOverlarge)
	}
	n := s.NewUInt8List(int32(len(v) + 1))
	copy(n.seg.Data[n.off:], v)
	return Pointer(n)
}
func (s *Segment) NewData(v []byte) Pointer {
	if int64(len(v)) > maxDataSize {
		panic(errOverlarge)
	}
	n := s.NewUInt8List(int32(len(v)))
	copy(n.seg.Data[n.off:], v)
	return Pointer(n)
}

func (s *Segment) NewRootStruct(sz ObjectSize) Struct {
	r, err := s.NewRoot()
	if err != nil {
		return Struct{}
	}
	v := s.NewStruct(sz)
	r.Set(0, Pointer(v))
	return v
}

func (s *Segment) NewStruct(sz ObjectSize) Struct {
	if !sz.isValid() {
		return Struct{}
	}
	sz.DataSize = sz.DataSize.padToWord()
	n, _ := s.create(sz.totalSize(), Pointer{typ: TypeStruct, size: sz})
	return Struct(n)
}

// NewStructAR (AutoRoot): experimental Root setting: assumes the
// struct is the root iff it is the first allocation in a segment.
func (s *Segment) NewStructAR(sz ObjectSize) Struct {
	if s.RootDone {
		return s.NewStruct(sz)
	} else {
		s.RootDone = true
		return s.NewRootStruct(sz)
	}
}

// create allocates sz bytes in the segment, using n as a template.
func (s *Segment) create(sz Size, n Pointer) (Pointer, error) {
	sz = sz.padToWord()

	// TODO(light): this can overflow easily
	if uint64(sz) > uint64(math.MaxUint32)-8 {
		return Pointer{}, errOverlarge
	}

	if s == nil {
		s = NewBuffer(nil)
	}

	tag := false
	newSize := Size(len(s.Data)) + sz
	addr := Address(len(s.Data))
	end := addr.addSize(sz)
	if newSize > Size(cap(s.Data)) {
		// If we can't fit the data in the current segment, we always
		// return a far pointer to a tag in the new segment.
		if (n.flags & isRoot) != 0 {
			tag = true
			sz += wordSize
		}
		news, err := s.Message.NewSegment(sz)
		if err != nil {
			return Pointer{}, err
		}

		// NewSegment is allowed to grow the segment and return it. In
		// which case we don't want to create a tag.
		if tag && news == s {
			sz -= wordSize
			tag = false
		}

		s = news
	}

	n.seg = s
	n.off = addr
	s.Data = s.Data[:newSize] // NewSegment() makes this promise

	if tag {
		putRawPointer(s.Data[addr:], n.value(addr))
		n.off = n.off.addSize(wordSize)
		n.flags |= hasPointerTag
	}

	for i := n.off; i < end; i++ {
		s.Data[i] = 0
	}

	return n, nil
}

// NewInterface returns a new interface pointer.
func (s *Segment) NewInterface(cap CapabilityID) Interface {
	return Interface(Pointer{
		seg: s,
		typ: TypeInterface,
		cap: cap,
	})
}

// listpos allows us to use this routine for copying elements between lists
func copyStructHandlingVersionSkew(dest, src Pointer, copies *rbtree.Tree, depth int, destListPos, srcListPos int32) error {
	// handle VoidList destinations
	if dest.seg == nil {
		return nil
	}

	destBase := dest.off.element(destListPos, dest.size.totalSize())
	destDataEnd := destBase.addSize(dest.size.DataSize)
	srcBase := src.off.element(srcListPos, src.size.totalSize())
	srcDataEnd := srcBase.addSize(src.size.DataSize)

	// Q: how does version handling happen here, when the
	//    desination toData[] slice can be bigger or smaller
	//    than the source data slice, which is in
	//    src.seg.Data[src.off:src.off+src.size.DataSize] ?
	//
	// A: Newer fields only come *after* old fields. Note that
	//    copy only copies min(len(src), len(dst)) size,
	//    and then we manually zero the rest in the for loop
	//    that writes toData[j] = 0.
	//

	// data section:
	toData := dest.seg.Data[destBase:destDataEnd]
	from := src.seg.Data[srcBase:srcDataEnd]
	copyCount := copy(toData, from)
	toData = toData[copyCount:]
	for j := range toData {
		toData[j] = 0
	}

	// ptrs section:

	// version handling: we ignore any extra-newer-pointers in src,
	// i.e. the case when srcPtrSize > dstPtrSize, by only
	// running j over the size of dstPtrSize, the destination size.
	numSrcPtrs := src.size.PointerCount
	numDstPtrs := dest.size.PointerCount
	for j := uint16(0); j < numSrcPtrs && j < numDstPtrs; j++ {
		m := src.seg.readPtr(srcDataEnd.element(int32(j), wordSize))
		err := dest.seg.writePtr(destDataEnd.element(int32(j), wordSize), m, copies, depth+1)
		if err != nil {
			return err
		}
	}
	for j := numSrcPtrs; j < numDstPtrs; j++ {
		// destination p is a newer version than source so these extra new pointer fields in p must be zeroed.
		addr := destDataEnd.element(int32(j), wordSize)
		putRawPointer(dest.seg.Data[addr:], 0)
	}
	// Nothing more here: so any other pointers in srcPtrSize beyond
	// those in dstPtrSize are ignored and discarded.

	return nil

} // end copyStructHandlingVersionSkew()

func (s *Segment) lookupSegment(id SegmentID) (*Segment, error) {
	if s.Id != id {
		return s.Message.Lookup(id)
	} else {
		return s, nil
	}
}

func (s *Segment) readPtr(off Address) Pointer {
	var err error
	val := rawPointer(binary.LittleEndian.Uint64(s.Data[off:]))

	switch val.pointerType() {
	case doubleFarPointer:
		// A double far pointer points to a double pointer, where the
		// first points to the actual data, and the second is the tag
		// that would normally be placed right before the data (offset
		// == 0).

		faroff, segid := val.farAddress(), val.farSegment()
		if s, err = s.lookupSegment(segid); err != nil || !s.regionInBounds(faroff, wordSize.times(2)) {
			return Pointer{}
		}
		far := rawPointer(binary.LittleEndian.Uint64(s.Data[faroff:]))
		tag := rawPointer(binary.LittleEndian.Uint64(s.Data[faroff.addSize(wordSize):]))
		if far.pointerType() != farPointer || tag.offset() != 0 {
			return Pointer{}
		}
		segid = far.farSegment()
		if s, err = s.lookupSegment(segid); err != nil {
			return Pointer{}
		}

		// -8 because far pointers reference from the start of the
		// segment, but offsets reference the end of the pointer data.
		val = tag | rawPointer(far.farAddress()-Address(wordSize))<<2

	case farPointer:
		faroff, segid := val.farAddress(), val.farSegment()

		if s, err = s.lookupSegment(segid); err != nil || !s.regionInBounds(faroff, wordSize) {
			return Pointer{}
		}

		off = faroff
		val = rawPointer(binary.LittleEndian.Uint64(s.Data[faroff:]))
	}

	if val == 0 {
		// This is a null pointer.
		return Pointer{}
	}

	// Be wary of overflow. Offset is 30 bits signed. List size is 29 bits
	// unsigned. For both of these we need to check in terms of words if
	// using 32 bit maths as bits or bytes will overflow.
	switch val.pointerType() {
	case structPointer:
		addr, ok := val.offset().resolve(off)
		if !ok {
			return Pointer{}
		}
		sz := val.structSize()
		if !s.regionInBounds(addr, sz.totalSize()) {
			return Pointer{}
		}
		return Pointer{
			seg:  s,
			typ:  TypeStruct,
			off:  addr,
			size: sz,
		}
	case listPointer:
		addr, ok := val.offset().resolve(off)
		if !ok {
			return Pointer{}
		}
		lt, lsize := val.listType(), val.totalListSize()
		if !s.regionInBounds(addr, lsize) {
			return Pointer{}
		}
		if lt == compositeList {
			hdr := rawPointer(binary.LittleEndian.Uint64(s.Data[addr:]))
			addr = addr.addSize(wordSize)
			if hdr.pointerType() != structPointer {
				return Pointer{}
			}
			sz := hdr.structSize()
			n := int32(hdr.offset())
			if !s.regionInBounds(addr, sz.totalSize().times(n)) {
				return Pointer{}
			}
			return Pointer{
				seg:    s,
				typ:    TypeList,
				size:   sz,
				off:    addr,
				length: n,
				flags:  isCompositeList,
			}
		}
		return Pointer{
			seg:    s,
			typ:    listPointerType(lt),
			size:   val.elementSize(),
			off:    addr,
			length: val.numListElements(),
		}
	case otherPointer:
		if val.otherPointerType() != 0 {
			return Pointer{}
		}
		return Pointer{
			seg: s,
			typ: TypeInterface,
			cap: val.capabilityIndex(),
		}

	default:
		return Pointer{}
	}
}

func listPointerType(listType int) PointerType {
	switch listType {
	case bit1List:
		return TypeBitList
	case pointerList:
		return TypePointerList
	default:
		return TypeList
	}
}

/*
lsb                       list pointer                        msb
+-+-----------------------------+--+----------------------------+
|A|             B               |C |             D              |
+-+-----------------------------+--+----------------------------+

A (2 bits) = 1, to indicate that this is a list pointer.
B (30 bits) = Offset, in words, from the end of the pointer to the
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
D (29 bits) = Number of elements in the list, except when C is 7
    (see below).

The pointed-to values are tightly-packed. In particular, Bools are packed bit-by-bit in little-endian order (the first bit is the least-significant bit of the first byte).

Lists of structs use the smallest element size in which the struct can fit. So, a list of structs that each contain two UInt8 fields and nothing else could be encoded with C = 3 (2-byte elements). A list of structs that each contain a single Text field would be encoded as C = 6 (pointer elements). A list of structs that each contain a single Bool field would be encoded using C = 1 (1-bit elements). A list structs which are each more than one word in size must be be encoded using C = 7 (composite).

When C = 7, the elements of the list are fixed-width composite values – usually, structs. In this case, the list content is prefixed by a "tag" word that describes each individual element. The tag has the same layout as a struct pointer, except that the pointer offset (B) instead indicates the number of elements in the list. Meanwhile, section (D) of the list pointer – which normally would store this element count – instead stores the total number of words in the list (not counting the tag word). The reason we store a word count in the pointer rather than an element count is to ensure that the extents of the list’s location can always be determined by inspecting the pointer alone, without having to look at the tag; this may allow more-efficient prefetching in some use cases. The reason we don’t store struct lists as a list of pointers is because doing so would take significantly more space (an extra pointer per element) and may be less cache-friendly.
*/

type offset struct {
	id         SegmentID
	boff, bend int64 // in bits
	newval     Pointer
}

func compare(a, b rbtree.Item) int {
	ao := a.(offset)
	bo := b.(offset)
	if ao.id != bo.id {
		return int(ao.id - bo.id)
	} else if ao.boff > bo.boff {
		return 1
	} else if ao.boff < bo.boff {
		return -1
	} else {
		return 0
	}
}

func (destSeg *Segment) writePtr(off Address, src Pointer, copies *rbtree.Tree, depth int) error {
	// handle size-zero Objects/empty structs
	if src.seg == nil {
		return nil
	}
	srcSeg := src.seg

	if src.typ == TypeNull || isEmptyStruct(src) {
		binary.LittleEndian.PutUint64(destSeg.Data[off:], 0)
		return nil

	} else if destSeg == srcSeg || src.typ == TypeInterface && destSeg.Message == srcSeg.Message {
		// Same segment
		binary.LittleEndian.PutUint64(destSeg.Data[off:], uint64(src.value(off)))
		return nil

	} else if src.typ == TypeInterface {
		// Different messages.  Need to copy table entry.
		c := destSeg.Message.AddCap(Interface(src).Client())
		p := Pointer(destSeg.NewInterface(c))
		binary.LittleEndian.PutUint64(destSeg.Data[off:], uint64(p.value(off)))
		return nil

	} else if destSeg.Message != srcSeg.Message || (src.flags&isListMember) != 0 || (src.flags&isBitListMember) != 0 {
		// We need to clone the target.

		if depth >= 32 {
			return errCopyDepth
		}

		// First see if the ptr has already been copied
		if copies == nil {
			copies = rbtree.NewTree(compare)
		}

		key := offset{
			id:   srcSeg.Id,
			boff: int64(src.off) * 8,
			bend: int64(src.objectEnd()) * 8,
			newval: Pointer{
				typ:    src.typ,
				length: src.length,
				size:   src.size,
				flags:  (src.flags & isCompositeList),
			},
		}

		if (src.flags & isBitListMember) != 0 {
			key.boff += int64(src.flags & bitOffsetMask)
			key.bend = key.boff + 1
			key.newval.size.DataSize = 8
		}

		if (src.flags & isCompositeList) != 0 {
			key.boff -= 64 //  Q: what the heck does this do? why is it here? A: Accounts for the Tag word, perhaps because objectEnd() does not.
		}

		iter := copies.FindLE(key)

		if key.bend > key.boff {
			if !iter.NegativeLimit() {
				other := iter.Item().(offset)
				if key.id == other.id {
					if key.boff == other.boff && key.bend == other.bend {
						return destSeg.writePtr(off, other.newval, nil, depth+1)
					} else if other.bend >= key.bend {
						return errOverlap
					}
				}
			}

			iter = iter.Next()

			if !iter.Limit() {
				other := iter.Item().(offset)
				if key.id == other.id && other.boff < key.bend {
					return errOverlap
				}
			}
		}

		// No copy nor overlap found, so we need to clone the target
		n, err := destSeg.create(Size((key.bend-key.boff)/8), key.newval)
		if err != nil {
			return err
		}

		// n is possibly in a new segment, if destSeg was full.
		newSeg := n.seg

		if (n.flags & isCompositeList) != 0 {
			copy(newSeg.Data[n.off:], srcSeg.Data[src.off-8:src.off])
			n.off += 8
		}

		key.newval = n
		copies.Insert(key)

		switch src.typ {
		case TypeStruct:
			if (src.flags & isBitListMember) != 0 {
				if (srcSeg.Data[src.off] & (1 << (src.flags & bitOffsetMask))) != 0 {
					newSeg.Data[n.off] = 1
				} else {
					newSeg.Data[n.off] = 0
				}

				for i := range newSeg.Data[n.off+1 : n.off+8] {
					newSeg.Data[i] = 0
				}
			} else {
				dest := Pointer{
					seg:  newSeg,
					off:  n.off,
					size: n.size,
				}

				if err := copyStructHandlingVersionSkew(dest, src, copies, depth, 0, 0); err != nil {
					return err
				}
			}

		case TypeList:
			// recognize Data and Text, both List(Byte), as special cases for speed.
			if n.size.isOneByte() && src.size.isOneByte() {
				copy(newSeg.Data[n.off:], srcSeg.Data[src.off:src.off.element(n.length, 1)])
				break
			}

			dest := Pointer{
				seg:  newSeg,
				off:  n.off,
				size: n.size,
			}
			for i := int32(0); i < n.length; i++ {
				if err := copyStructHandlingVersionSkew(dest, src, copies, depth, i, i); err != nil {
					return err
				}
			}

		case TypePointerList:
			for i := int32(0); i < n.length; i++ {
				c := srcSeg.readPtr(src.off.element(i, wordSize))
				if err := newSeg.writePtr(n.off.element(i, wordSize), c, copies, depth+1); err != nil {
					return err
				}
			}

		case TypeBitList:
			// TODO(light): DataSize is not populated for BitList.
			nbytes := int(src.length + 7/8)
			listSize := Size(nbytes).padToWord()
			copy(newSeg.Data[n.off:], srcSeg.Data[src.off:src.off.addSize(listSize)])
		}
		return destSeg.writePtr(off, key.newval, nil, depth+1)

	} else if (src.flags & hasPointerTag) != 0 {
		// By lucky chance, the data has a tag in front of it. This
		// happens when create had to move the data to a new segment.
		binary.LittleEndian.PutUint64(destSeg.Data[off:], uint64(rawFarPointer(srcSeg.Id, src.off-8)))
		return nil

	} else if len(srcSeg.Data)+int(wordSize) <= cap(srcSeg.Data) {
		// Have room in the target for a tag
		srcAddr := Address(len(srcSeg.Data))
		srcSeg.Data = srcSeg.Data[:srcAddr.addSize(wordSize)]
		binary.LittleEndian.PutUint64(srcSeg.Data[srcAddr:], uint64(src.value(srcAddr)))
		binary.LittleEndian.PutUint64(destSeg.Data[off:], uint64(rawFarPointer(srcSeg.Id, srcAddr)))
		return nil

	} else {
		// Need to create a double far pointer. Try and create it in
		// the originating segment if we can.
		t := destSeg
		const landingSize = wordSize * 2
		if len(t.Data)+int(landingSize) > cap(t.Data) {
			var err error
			if t, err = t.Message.NewSegment(landingSize); err != nil {
				return err
			}
		}

		dstAddr := Address(len(t.Data))
		binary.LittleEndian.PutUint64(t.Data[dstAddr:], uint64(rawFarPointer(srcSeg.Id, src.off)))
		binary.LittleEndian.PutUint64(t.Data[dstAddr.addSize(wordSize):], uint64(src.value(src.off-8)))
		binary.LittleEndian.PutUint64(destSeg.Data[off:], uint64(rawDoubleFarPointer(t.Id, dstAddr)))
		t.Data = t.Data[:len(t.Data)+16]
		return nil
	}
}

// An Interface is a wrapper for a Pointer that provides methods to
// access interface information.
type Interface Pointer

// Segment returns the segment this pointer came from.
func (i Interface) Segment() *Segment {
	return i.seg
}

// Capability returns the capability ID of the interface.
func (i Interface) Capability() CapabilityID {
	return i.cap
}

// Client returns the client stored in the message's capability table
// or nil if the pointer is invalid.
func (i Interface) Client() Client {
	if i.seg == nil || i.typ == TypeNull {
		return nil
	}
	tab := i.seg.Message.CapTable()
	if uint64(i.cap) >= uint64(len(tab)) {
		return nil
	}
	return tab[i.cap]
}
