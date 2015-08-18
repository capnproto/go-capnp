package capnp

import (
	"encoding/binary"
	"errors"

	"github.com/glycerine/rbtree"
)

// A SegmentID is a numeric identifier for a Segment.
type SegmentID uint32

// A Segment is an allocation arena for Cap'n Proto objects.
// It is part of a Message, which can contain other segments that
// reference each other.
type Segment struct {
	msg  *Message
	id   SegmentID
	data []byte
}

// Message returns the message that contains s.
func (s *Segment) Message() *Message {
	return s.msg
}

// ID returns the segment's ID.
func (s *Segment) ID() SegmentID {
	return s.id
}

// Data returns the raw byte slice for the segment.
func (s *Segment) Data() []byte {
	return s.data
}

func (s *Segment) inBounds(addr Address) bool {
	return addr < Address(len(s.data))
}

func (s *Segment) regionInBounds(base Address, sz Size) bool {
	return base.addSize(sz) <= Address(len(s.data))
}

// slice returns the segment of data from base to base+sz.
func (s *Segment) slice(base Address, sz Size) []byte {
	if !s.regionInBounds(base, sz) {
		return nil
	}
	return s.data[base:base.addSize(sz)]
}

func (s *Segment) readUint16(addr Address) uint16 {
	data := s.slice(addr, 2)
	if data == nil {
		panic(errOutOfBounds)
	}
	return binary.LittleEndian.Uint16(data)
}

func (s *Segment) readUint32(addr Address) uint32 {
	data := s.slice(addr, 4)
	if data == nil {
		panic(errOutOfBounds)
	}
	return binary.LittleEndian.Uint32(data)
}

func (s *Segment) readUint64(addr Address) uint64 {
	data := s.slice(addr, 8)
	if data == nil {
		panic(errOutOfBounds)
	}
	return binary.LittleEndian.Uint64(data)
}

func (s *Segment) readRawPointer(addr Address) rawPointer {
	return rawPointer(s.readUint64(addr))
}

func (s *Segment) writeUint16(addr Address, val uint16) {
	data := s.slice(addr, 2)
	if data == nil {
		panic(errOutOfBounds)
	}
	binary.LittleEndian.PutUint16(data, val)
}

func (s *Segment) writeUint32(addr Address, val uint32) {
	data := s.slice(addr, 4)
	if data == nil {
		panic(errOutOfBounds)
	}
	binary.LittleEndian.PutUint32(data, val)
}

func (s *Segment) writeUint64(addr Address, val uint64) {
	data := s.slice(addr, 8)
	if data == nil {
		panic(errOutOfBounds)
	}
	binary.LittleEndian.PutUint64(data, val)
}

func (s *Segment) writeRawPointer(addr Address, val rawPointer) {
	s.writeUint64(addr, uint64(val))
}

// root returns a 1-element pointer list that references the first word
// in the segment.  This only makes sense to call on the first segment
// in a message.
func (s *Segment) root() PointerList {
	sz := ObjectSize{PointerCount: 1}
	if !s.regionInBounds(0, sz.totalSize()) {
		return PointerList{}
	}
	return PointerList{List{
		seg:    s,
		length: 1,
		size:   sz,
	}}
}

func (s *Segment) lookupSegment(id SegmentID) (*Segment, error) {
	if s.id == id {
		return s, nil
	}
	return s.msg.Segment(id)
}

func (s *Segment) readPtr(off Address) (Pointer, error) {
	var err error
	val := s.readRawPointer(off)
	s, off, val, err = s.resolveFarPointer(off, val)
	if err != nil {
		return nil, err
	}
	if val == 0 {
		return nil, nil
	}
	// Be wary of overflow. Offset is 30 bits signed. List size is 29 bits
	// unsigned. For both of these we need to check in terms of words if
	// using 32 bit maths as bits or bytes will overflow.
	switch val.pointerType() {
	case structPointer:
		addr, ok := val.offset().resolve(off)
		if !ok {
			return nil, errPointerAddress
		}
		sz := val.structSize()
		if !s.regionInBounds(addr, sz.totalSize()) {
			return nil, errPointerAddress
		}
		return Struct{
			seg:  s,
			off:  addr,
			size: sz,
		}, nil
	case listPointer:
		addr, ok := val.offset().resolve(off)
		if !ok {
			return nil, errPointerAddress
		}
		lt, lsize := val.listType(), val.totalListSize()
		if !s.regionInBounds(addr, lsize) {
			return nil, errPointerAddress
		}
		if lt == compositeList {
			hdr := s.readRawPointer(addr)
			addr = addr.addSize(wordSize)
			if hdr.pointerType() != structPointer {
				return nil, errBadTag
			}
			sz := hdr.structSize()
			n := int32(hdr.offset())
			// TODO(light): check that this has the same end address
			if !s.regionInBounds(addr, sz.totalSize().times(n)) {
				return nil, errPointerAddress
			}
			return List{
				seg:    s,
				size:   sz,
				off:    addr,
				length: n,
				flags:  isCompositeList,
			}, nil
		}
		if lt == bit1List {
			return List{
				seg:    s,
				off:    addr,
				length: val.numListElements(),
				flags:  isBitList,
			}, nil
		}
		return List{
			seg:    s,
			size:   val.elementSize(),
			off:    addr,
			length: val.numListElements(),
		}, nil
	case otherPointer:
		if val.otherPointerType() != 0 {
			return nil, errOtherPointer
		}
		return Interface{
			seg: s,
			cap: val.capabilityIndex(),
		}, nil
	default:
		// Only other types are far pointers.
		return nil, errBadLandingPad
	}
}

func (s *Segment) resolveFarPointer(off Address, val rawPointer) (*Segment, Address, rawPointer, error) {
	switch val.pointerType() {
	case doubleFarPointer:
		// A double far pointer points to a double pointer, where the
		// first points to the actual data, and the second is the tag
		// that would normally be placed right before the data (offset
		// == 0).

		faroff, segid := val.farAddress(), val.farSegment()
		s, err := s.lookupSegment(segid)
		if err != nil {
			return nil, 0, 0, err
		}
		if !s.regionInBounds(faroff, wordSize.times(2)) {
			return nil, 0, 0, errPointerAddress
		}
		far := s.readRawPointer(faroff)
		tag := s.readRawPointer(faroff.addSize(wordSize))
		if far.pointerType() != farPointer || tag.offset() != 0 {
			return nil, 0, 0, errPointerAddress
		}
		segid = far.farSegment()
		if s, err = s.lookupSegment(segid); err != nil {
			return nil, 0, 0, errBadLandingPad
		}
		return s, 0, landingPadNearPointer(far, tag), nil
	case farPointer:
		faroff, segid := val.farAddress(), val.farSegment()
		s, err := s.lookupSegment(segid)
		if err != nil {
			return nil, 0, 0, err
		}
		if !s.regionInBounds(faroff, wordSize) {
			return nil, 0, 0, errPointerAddress
		}
		val = s.readRawPointer(faroff)
		return s, faroff, val, nil
	default:
		return s, off, val, nil
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

func makeOffsetKey(p Pointer) offset {
	switch p := p.underlying().(type) {
	case Struct:
		return offset{
			id:   p.seg.id,
			boff: int64(p.off) * 8,
			bend: int64(p.off.addSize(p.size.totalSize())) * 8,
		}
	case List:
		key := offset{
			id:   p.seg.id,
			boff: int64(p.off) * 8,
			bend: int64(p.off.addSize(p.size.totalSize())) * 8,
		}
		if p.flags&isBitList != 0 {
			key.bend = int64(p.off)*8 + int64(p.length)
		} else {
			key.bend = int64(p.off.addSize(p.size.totalSize())) * 8
		}
		if p.flags&isCompositeList != 0 {
			// Composite lists' offsets are after the tag word.
			key.boff -= int64(wordSize) * 8
		}
		return key
	default:
		panic("unreachable")
	}
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

func needsCopy(dest *Segment, src Pointer) bool {
	if src.Segment().msg != dest.msg {
		return true
	}
	if s := ToStruct(src); IsValid(s) {
		// Structs can only be referenced if they're not list members.
		return s.flags&isListMember != 0
	}
	return false
}

func (destSeg *Segment) writePtr(cc copyContext, off Address, src Pointer) error {
	// handle nulls
	if !IsValid(src) {
		destSeg.writeRawPointer(off, 0)
		return nil
	}
	srcSeg := src.Segment()

	if i := ToInterface(src); IsValid(i) {
		if destSeg.msg != srcSeg.msg {
			c := destSeg.msg.AddCap(i.Client())
			src = Pointer(NewInterface(destSeg, c))
		}
		destSeg.writeRawPointer(off, src.value(off))
		return nil
	}
	if destSeg != srcSeg {
		// Different segments
		if needsCopy(destSeg, src) {
			return copyPointer(cc, destSeg, off, src)
		}
		if !hasCapacity(srcSeg.data, wordSize) {
			// Double far pointer needed.
			const landingSize = wordSize * 2
			t, dstAddr, err := alloc(destSeg, landingSize)
			if err != nil {
				return err
			}

			srcAddr := pointerAddress(src)
			t.writeRawPointer(dstAddr, rawFarPointer(srcSeg.id, srcAddr))
			t.writeRawPointer(dstAddr.addSize(wordSize), src.value(srcAddr-Address(wordSize)))
			destSeg.writeRawPointer(off, rawDoubleFarPointer(t.id, dstAddr))
			return nil
		}
		// Have room in the target for a tag
		_, srcAddr, _ := alloc(srcSeg, wordSize)
		srcSeg.writeRawPointer(srcAddr, src.value(srcAddr))
		destSeg.writeRawPointer(off, rawFarPointer(srcSeg.id, srcAddr))
		return nil
	}
	destSeg.writeRawPointer(off, src.value(off))
	return nil
}

func copyPointer(cc copyContext, dstSeg *Segment, dstAddr Address, src Pointer) error {
	if cc.depth >= 32 {
		return errCopyDepth
	}
	cc = cc.init()
	// First, see if the ptr has already been copied.
	key := makeOffsetKey(src)
	iter := cc.copies.FindLE(key)
	if key.bend > key.boff {
		if !iter.NegativeLimit() {
			other := iter.Item().(offset)
			if key.id == other.id {
				if key.boff == other.boff && key.bend == other.bend {
					return dstSeg.writePtr(cc.incDepth(), dstAddr, other.newval)
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
	newSeg, newAddr, err := alloc(dstSeg, Size((key.bend-key.boff)/8))
	if err != nil {
		return err
	}
	switch src := src.underlying().(type) {
	case Struct:
		dst := Struct{
			seg:  newSeg,
			off:  newAddr,
			size: src.size,
			// clear flags
		}
		key.newval = dst
		cc.copies.Insert(key)
		if err := copyStruct(cc, dst, src); err != nil {
			return err
		}
	case List:
		dst := List{
			seg:    newSeg,
			off:    newAddr,
			length: src.length,
			size:   src.size,
			flags:  src.flags,
		}
		if dst.flags&isCompositeList != 0 {
			// Copy tag word
			copy(newSeg.data[newAddr:], src.seg.data[src.off-Address(wordSize):src.off])
			dst.off = dst.off.addSize(wordSize)
		}
		key.newval = dst
		cc.copies.Insert(key)
		// TODO(light): fast path for copying text/data
		if dst.flags&isBitList != 0 {
			copy(newSeg.data[newAddr:], src.seg.data[src.off:src.length+7/8])
		} else {
			for i := 0; i < src.Len(); i++ {
				err := copyStruct(cc, dst.Struct(i), src.Struct(i))
				if err != nil {
					return err
				}
			}
		}
	default:
		panic("unreachable")
	}
	return dstSeg.writePtr(cc.incDepth(), dstAddr, key.newval)
}

type copyContext struct {
	copies *rbtree.Tree
	depth  int
}

func (cc copyContext) init() copyContext {
	if cc.copies == nil {
		return copyContext{
			copies: rbtree.NewTree(compare),
		}
	}
	return cc
}

func (cc copyContext) incDepth() copyContext {
	return copyContext{
		copies: cc.copies,
		depth:  cc.depth + 1,
	}
}

var (
	errPointerAddress = errors.New("capnp: invalid pointer address")
	errBadLandingPad  = errors.New("capnp: invalid far pointer landing pad")
	errBadTag         = errors.New("capnp: invalid tag word")
	errOtherPointer   = errors.New("capnp: unknown pointer type")
	errObjectSize     = errors.New("capnp: invalid object size")
)

var (
	errOverlarge   = errors.New("capn: overlarge struct/list")
	errOutOfBounds = errors.New("capn: write out of bounds")
	errCopyDepth   = errors.New("capn: copy depth too large")
	errOverlap     = errors.New("capn: overlapping data on copy")
	errListSize    = errors.New("capn: invalid list size")
	errObjectType  = errors.New("capn: invalid object type")
)
