package capnp

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"zombiezen.com/go/capnproto2/internal/errors"
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
	end, ok := base.addSize(sz)
	if !ok {
		return false
	}
	return end <= Address(len(s.data))
}

// slice returns the segment of data from base to base+sz.
func (s *Segment) slice(base Address, sz Size) []byte {
	// Bounds check should have happened before calling slice.
	return s.data[base : base+Address(sz)]
}

func (s *Segment) readUint8(addr Address) uint8 {
	return s.slice(addr, 1)[0]
}

func (s *Segment) readUint16(addr Address) uint16 {
	return binary.LittleEndian.Uint16(s.slice(addr, 2))
}

func (s *Segment) readUint32(addr Address) uint32 {
	return binary.LittleEndian.Uint32(s.slice(addr, 4))
}

func (s *Segment) readUint64(addr Address) uint64 {
	return binary.LittleEndian.Uint64(s.slice(addr, 8))
}

func (s *Segment) readRawPointer(addr Address) rawPointer {
	return rawPointer(s.readUint64(addr))
}

func (s *Segment) writeUint8(addr Address, val uint8) {
	s.slice(addr, 1)[0] = val
}

func (s *Segment) writeUint16(addr Address, val uint16) {
	binary.LittleEndian.PutUint16(s.slice(addr, 2), val)
}

func (s *Segment) writeUint32(addr Address, val uint32) {
	binary.LittleEndian.PutUint32(s.slice(addr, 4), val)
}

func (s *Segment) writeUint64(addr Address, val uint64) {
	binary.LittleEndian.PutUint64(s.slice(addr, 8), val)
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
		seg:        s,
		length:     1,
		size:       sz,
		depthLimit: s.msg.depthLimit(),
	}}
}

func (s *Segment) lookupSegment(id SegmentID) (*Segment, error) {
	if s.id == id {
		return s, nil
	}
	return s.msg.Segment(id)
}

func (s *Segment) readPtr(paddr Address, depthLimit uint) (ptr Ptr, err error) {
	s, base, val, err := s.resolveFarPointer(paddr)
	if err != nil {
		return Ptr{}, annotate(err).errorf("read pointer")
	}
	if val == 0 {
		return Ptr{}, nil
	}
	if depthLimit == 0 {
		return Ptr{}, newError("read pointer: depth limit reached")
	}
	switch val.pointerType() {
	case structPointer:
		sp, err := s.readStructPtr(base, val)
		if err != nil {
			return Ptr{}, annotate(err).errorf("read pointer")
		}
		if !s.msg.canRead(sp.readSize()) {
			return Ptr{}, newError("read pointer: read traversal limit reached")
		}
		sp.depthLimit = depthLimit - 1
		return sp.ToPtr(), nil
	case listPointer:
		lp, err := s.readListPtr(base, val)
		if err != nil {
			return Ptr{}, annotate(err).errorf("read pointer")
		}
		if !s.msg.canRead(lp.readSize()) {
			return Ptr{}, newError("read pointer: read traversal limit reached")
		}
		lp.depthLimit = depthLimit - 1
		return lp.ToPtr(), nil
	case otherPointer:
		if val.otherPointerType() != 0 {
			return Ptr{}, newError("read pointer: unknown pointer type")
		}
		return Interface{
			seg: s,
			cap: val.capabilityIndex(),
		}.ToPtr(), nil
	default:
		// Only other types are far pointers.
		return Ptr{}, newError("read pointer: far pointer landing pad is a far pointer")
	}
}

func (s *Segment) readStructPtr(base Address, val rawPointer) (Struct, error) {
	addr, ok := val.offset().resolve(base)
	if !ok {
		return Struct{}, newError("struct pointer: invalid address")
	}
	sz := val.structSize()
	if !s.regionInBounds(addr, sz.totalSize()) {
		return Struct{}, newError("struct pointer: invalid address")
	}
	return Struct{
		seg:  s,
		off:  addr,
		size: sz,
	}, nil
}

func (s *Segment) readListPtr(base Address, val rawPointer) (List, error) {
	addr, ok := val.offset().resolve(base)
	if !ok {
		return List{}, newError("list pointer: invalid address")
	}
	lsize, ok := val.totalListSize()
	if !ok {
		return List{}, newError("list pointer: size overflow")
	}
	if !s.regionInBounds(addr, lsize) {
		return List{}, newError("list pointer: address out of bounds")
	}
	lt := val.listType()
	if lt == compositeList {
		hdr := s.readRawPointer(addr)
		var ok bool
		addr, ok = addr.addSize(wordSize)
		if !ok {
			return List{}, newError("composite list pointer: content address overflow")
		}
		if hdr.pointerType() != structPointer {
			return List{}, newError("composite list pointer: tag word is not a struct")
		}
		sz := hdr.structSize()
		n := int32(hdr.offset())
		// TODO(light): check that this has the same end address
		if tsize, ok := sz.totalSize().times(n); !ok {
			return List{}, newError("composite list pointer: size overflow")
		} else if !s.regionInBounds(addr, tsize) {
			return List{}, newError("composite list pointer: address out of bounds")
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
}

func (s *Segment) resolveFarPointer(paddr Address) (dst *Segment, base Address, resolved rawPointer, err error) {
	// Encoding details at https://capnproto.org/encoding.html#inter-segment-pointers

	val := s.readRawPointer(paddr)
	switch val.pointerType() {
	case doubleFarPointer:
		padSeg, err := s.lookupSegment(val.farSegment())
		if err != nil {
			return nil, 0, 0, annotate(err).errorf("double-far pointer")
		}
		padAddr := val.farAddress()
		if !padSeg.regionInBounds(padAddr, wordSize*2) {
			return nil, 0, 0, newError("double-far pointer: address out of bounds")
		}
		far := padSeg.readRawPointer(padAddr)
		if far.pointerType() != farPointer {
			return nil, 0, 0, newError("double-far pointer: first word in landing pad is not a far pointer")
		}
		tagAddr, ok := padAddr.addSize(wordSize)
		if !ok {
			return nil, 0, 0, newError("double-far pointer: landing pad address overflow")
		}
		tag := padSeg.readRawPointer(tagAddr)
		if pt := tag.pointerType(); (pt != structPointer && pt != listPointer) || tag.offset() != 0 {
			return nil, 0, 0, newError("double-far pointer: second word is not a struct or list with zero offset")
		}
		if dst, err = s.lookupSegment(far.farSegment()); err != nil {
			return nil, 0, 0, annotate(err).errorf("double-far pointer")
		}
		return dst, 0, landingPadNearPointer(far, tag), nil
	case farPointer:
		var err error
		dst, err = s.lookupSegment(val.farSegment())
		if err != nil {
			return nil, 0, 0, annotate(err).errorf("far pointer")
		}
		padAddr := val.farAddress()
		if !dst.regionInBounds(padAddr, wordSize) {
			return nil, 0, 0, newError("far pointer: address out of bounds")
		}
		var ok bool
		base, ok = padAddr.addSize(wordSize)
		if !ok {
			return nil, 0, 0, newError("far pointer: landing pad address overflow")
		}
		return dst, base, dst.readRawPointer(padAddr), nil
	default:
		var ok bool
		base, ok = paddr.addSize(wordSize)
		if !ok {
			return nil, 0, 0, newError("pointer base address overflow")
		}
		return s, base, val, nil
	}
}

func (s *Segment) writePtr(off Address, src Ptr, forceCopy bool) error {
	if !src.IsValid() {
		s.writeRawPointer(off, 0)
		return nil
	}

	// Copy src, if needed, and process pointers where placement is
	// irrelevant (capabilities and zero-sized structs).
	var srcAddr Address
	var srcRaw rawPointer
	switch src.flags.ptrType() {
	case structPtrType:
		st := src.Struct()
		if st.size.isZero() {
			// Zero-sized structs should always be encoded with offset -1 in
			// order to avoid conflating with null.  No allocation needed.
			s.writeRawPointer(off, rawStructPointer(-1, ObjectSize{}))
			return nil
		}
		if forceCopy || src.seg.msg != s.msg || st.flags&isListMember != 0 {
			newSeg, newAddr, err := alloc(s, st.size.totalSize())
			if err != nil {
				return annotate(err).errorf("write pointer: copy")
			}
			dst := Struct{
				seg:        newSeg,
				off:        newAddr,
				size:       st.size,
				depthLimit: maxDepth,
				// clear flags
			}
			if err := copyStruct(dst, st); err != nil {
				return annotate(err).errorf("write pointer")
			}
			st = dst
			src = dst.ToPtr()
		}
		srcAddr = st.off
		srcRaw = rawStructPointer(0, st.size)
	case listPtrType:
		l := src.List()
		if forceCopy || src.seg.msg != s.msg {
			sz := l.allocSize()
			newSeg, newAddr, err := alloc(s, sz)
			if err != nil {
				return annotate(err).errorf("write pointer: copy")
			}
			dst := List{
				seg:        newSeg,
				off:        newAddr,
				length:     l.length,
				size:       l.size,
				flags:      l.flags,
				depthLimit: maxDepth,
			}
			if dst.flags&isCompositeList != 0 {
				// Copy tag word
				newSeg.writeRawPointer(newAddr, l.seg.readRawPointer(l.off-Address(wordSize)))
				var ok bool
				dst.off, ok = dst.off.addSize(wordSize)
				if !ok {
					return newError("write pointer: copy composite list: content address overflow")
				}
				sz -= wordSize
			}
			if dst.flags&isBitList != 0 || dst.size.PointerCount == 0 {
				end, _ := l.off.addSize(sz) // list was already validated
				copy(newSeg.data[dst.off:], l.seg.data[l.off:end])
			} else {
				for i := 0; i < l.Len(); i++ {
					err := copyStruct(dst.Struct(i), l.Struct(i))
					if err != nil {
						return annotate(err).errorf("write pointer: copy list element %d", i)
					}
				}
			}
			l = dst
			src = dst.ToPtr()
		}
		srcAddr = l.off
		if l.flags&isCompositeList != 0 {
			srcAddr -= Address(wordSize)
		}
		srcRaw = l.raw()
	case interfacePtrType:
		i := src.Interface()
		if src.seg.msg != s.msg {
			c := s.msg.AddCap(i.Client())
			i = NewInterface(s, c)
		}
		s.writeRawPointer(off, i.value(off))
		return nil
	default:
		panic("unreachable")
	}

	switch {
	case src.seg == s:
		// Common case: src is in same segment as pointer.
		// Use a near pointer.
		s.writeRawPointer(off, srcRaw.withOffset(nearPointerOffset(off, srcAddr)))
		return nil
	case hasCapacity(src.seg.data, wordSize):
		// Enough room adjacent to src to write a far pointer landing pad.
		_, padAddr, _ := alloc(src.seg, wordSize)
		src.seg.writeRawPointer(padAddr, srcRaw.withOffset(nearPointerOffset(padAddr, srcAddr)))
		s.writeRawPointer(off, rawFarPointer(src.seg.id, padAddr))
		return nil
	default:
		// Not enough room for a landing pad, need to use a double-far pointer.
		padSeg, padAddr, err := alloc(s, wordSize*2)
		if err != nil {
			return annotate(err).errorf("write pointer: make landing pad")
		}
		padSeg.writeRawPointer(padAddr, rawFarPointer(src.seg.id, srcAddr))
		padSeg.writeRawPointer(padAddr+Address(wordSize), srcRaw)
		s.writeRawPointer(off, rawDoubleFarPointer(padSeg.id, padAddr))
		return nil
	}
}

// Equal returns true iff p1 and p2 are equal.
//
// Equality is defined to be:
//
//	- Two structs are equal iff all of their fields are equal.  If one
//	  struct has more fields than the other, the extra fields must all be
//		zero.
//	- Two lists are equal iff they have the same length and their
//	  corresponding elements are equal.  If one list is a list of
//	  primitives and the other is a list of structs, then the list of
//	  primitives is treated as if it was a list of structs with the
//	  element value as the sole field.
//	- Two interfaces are equal iff they point to a capability created by
//	  the same call to NewClient or they are referring to the same
//	  capability table index in the same message.  The latter is
//	  significant when the message's capability table has not been
//	  populated.
//	- Two null pointers are equal.
//	- All other combinations of things are not equal.
func Equal(p1, p2 Ptr) (bool, error) {
	if !p1.IsValid() && !p2.IsValid() {
		return true, nil
	}
	if !p1.IsValid() || !p2.IsValid() {
		return false, nil
	}
	pt := p1.flags.ptrType()
	if pt != p2.flags.ptrType() {
		return false, nil
	}
	switch pt {
	case structPtrType:
		s1, s2 := p1.Struct(), p2.Struct()
		data1 := s1.seg.slice(s1.off, s1.size.DataSize)
		data2 := s2.seg.slice(s2.off, s2.size.DataSize)
		switch {
		case len(data1) < len(data2):
			if !bytes.Equal(data1, data2[:len(data1)]) {
				return false, nil
			}
			if !isZeroFilled(data2[len(data1):]) {
				return false, nil
			}
		case len(data1) > len(data2):
			if !bytes.Equal(data1[:len(data2)], data2) {
				return false, nil
			}
			if !isZeroFilled(data1[len(data2):]) {
				return false, nil
			}
		default:
			if !bytes.Equal(data1, data2) {
				return false, nil
			}
		}
		n := int(s1.size.PointerCount)
		if n2 := int(s2.size.PointerCount); n2 < n {
			n = n2
		}
		for i := 0; i < n; i++ {
			sp1, err := s1.Ptr(uint16(i))
			if err != nil {
				return false, annotate(err).errorf("equal")
			}
			sp2, err := s2.Ptr(uint16(i))
			if err != nil {
				return false, annotate(err).errorf("equal")
			}
			if ok, err := Equal(sp1, sp2); !ok || err != nil {
				return false, err
			}
		}
		for i := n; i < int(s1.size.PointerCount); i++ {
			if s1.HasPtr(uint16(i)) {
				return false, nil
			}
		}
		for i := n; i < int(s2.size.PointerCount); i++ {
			if s2.HasPtr(uint16(i)) {
				return false, nil
			}
		}
		return true, nil
	case listPtrType:
		l1, l2 := p1.List(), p2.List()
		if l1.Len() != l2.Len() {
			return false, nil
		}
		if l1.flags&isCompositeList == 0 && l2.flags&isCompositeList == 0 && l1.size != l2.size {
			return false, nil
		}
		if l1.size.PointerCount == 0 && l2.size.PointerCount == 0 && l1.size.DataSize == l2.size.DataSize {
			// Optimization: pure data lists can be compared bytewise.
			sz, _ := l1.size.totalSize().times(l1.length) // both list bounds have been validated
			return bytes.Equal(l1.seg.slice(l1.off, sz), l2.seg.slice(l2.off, sz)), nil
		}
		for i := 0; i < l1.Len(); i++ {
			e1, e2 := l1.Struct(i), l2.Struct(i)
			if ok, err := Equal(e1.ToPtr(), e2.ToPtr()); err != nil {
				return false, annotate(err).errorf("equal: list element %d", i)
			} else if !ok {
				return false, nil
			}
		}
		return true, nil
	case interfacePtrType:
		i1, i2 := p1.Interface(), p2.Interface()
		if i1.Message() == i2.Message() {
			if i1.Capability() == i2.Capability() {
				return true, nil
			}
			ntab := len(i1.Message().CapTable)
			if int64(i1.Capability()) >= int64(ntab) || int64(i2.Capability()) >= int64(ntab) {
				return false, nil
			}
		}
		return i1.Client().IsSame(i2.Client()), nil
	default:
		panic("unreachable")
	}
}

func isZeroFilled(b []byte) bool {
	for _, bb := range b {
		if bb != 0 {
			return false
		}
	}
	return true
}

func newError(msg string) error {
	return errors.New(errors.Failed, "capnp", msg)
}

func errorf(format string, args ...interface{}) error {
	return newError(fmt.Sprintf(format, args...))
}

type annotater struct {
	err error
}

func annotate(err error) annotater {
	return annotater{err}
}

func (a annotater) errorf(format string, args ...interface{}) error {
	return errors.Annotate("capnp", fmt.Sprintf(format, args...), a.err)
}
