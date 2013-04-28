package capn

import (
	"encoding/binary"
	"errors"
	"math"
)

var (
	little16    = binary.LittleEndian.Uint16
	little32    = binary.LittleEndian.Uint32
	little64    = binary.LittleEndian.Uint64
	putLittle16 = binary.LittleEndian.PutUint16
	putLittle32 = binary.LittleEndian.PutUint32
	putLittle64 = binary.LittleEndian.PutUint64

	ErrInvalidInterface = errors.New("capn: invalid interface")
	ErrInvalidPointer   = errors.New("capn: invalid pointer")
	ErrAlignment        = errors.New("capn: segment is not aligned to 8 byte boundary")
	ErrOverlarge        = errors.New("capn: segment/struct/list is too large")
	ErrOutOfBounds      = errors.New("capn: write out of bounds")
)

type PointerType uint8

const (
	Null PointerType = iota
	Struct
	List
	PointerList
	FarPointer
	DoubleFarPointer
	Root
)

type Session interface {
	NewSegment(minsz int) (*Segment, error)
	Lookup(segid uint32) (*Segment, error)
	NewCall() (*Call, error)
}

type Marshaller interface {
	MarshalCaptain(p Pointer, off int) error
}

type Segment struct {
	Session Session
	Data    []uint8
	Id      uint32
}

type Pointer struct {
	Segment  *Segment
	off      uint32 // in bits
	size     uint32 // also used for the segment id
	dataBits uint32
	ptrs     uint16
	typ      PointerType
}

const (
	maxDataBits = 0xFFFF * 64
	maxPtrs     = 0xFFFF
)

func (s *Segment) NewRoot() (Pointer, int, error) {
	n := Pointer{typ: Root, size: 1, ptrs: 1}
	n, err := s.create(n)
	return n, int(n.off / 64), err
}

func (s *Segment) Root() Pointer {
	return Pointer{Segment: s, typ: Root, size: uint32(len(s.Data) / 8), ptrs: 1}
}

func (s *Segment) NewPointerList(size int) (Pointer, error) {
	n := Pointer{typ: PointerList, size: uint32(size), ptrs: 1}
	return s.create(n)
}

func (s *Segment) NewList(dataBits, ptrs, size int) (Pointer, error) {
	if dataBits < 0 || dataBits > maxDataBits || ptrs < 0 || ptrs > maxPtrs {
		return Pointer{}, ErrOverlarge
	} else if ptrs > 0 || dataBits > 32 {
		dataBits = (dataBits + 63) &^ 63
	} else if dataBits > 16 {
		dataBits = (dataBits + 31) &^ 31
	} else if dataBits > 8 {
		dataBits = (dataBits + 15) &^ 15
	} else if dataBits > 1 {
		dataBits = (dataBits + 7) &^ 7
	}
	n := Pointer{typ: List, size: uint32(size), dataBits: uint32(dataBits), ptrs: uint16(ptrs)}
	return s.create(n)
}

func (s *Segment) NewStruct(dataBits, ptrs int) (Pointer, error) {
	if dataBits < 0 || dataBits > maxDataBits || ptrs < 0 || ptrs > maxPtrs {
		return Pointer{}, ErrOverlarge
	}
	dataBits = (dataBits + 63) &^ 63
	n := Pointer{typ: Struct, size: 1, dataBits: uint32(dataBits), ptrs: uint16(ptrs)}
	return s.create(n)
}

func (p Pointer) Type() PointerType { return p.typ }
func (p Pointer) Size() int         { return int(p.size) }
func (p Pointer) end() uint32       { return (p.off + p.size*(p.dataBits+uint32(p.ptrs)*64) + 7) / 8 }

func (p Pointer) Data() []byte {
	if q, err := p.Deref(); err == nil {
		return q.Segment.Data[q.off/8 : q.end()]
	}
	return nil
}

func (p Pointer) ReadStruct1(bitoff int) (ret bool) {
	if q, err := p.Deref(); err == nil && uint(bitoff) < uint(q.dataBits) {
		off := uint(q.off) + uint(bitoff)
		ret = q.Segment.Data[off/8]&(1<<uint(off%8)) != 0
	}
	return
}

func (p Pointer) ReadStruct8(bitoff int) (ret uint8) {
	if q, err := p.Deref(); err == nil && uint(bitoff) < uint(q.dataBits) {
		ret = q.Segment.Data[(uint(q.off)+uint(bitoff))/8]
	}
	return
}

func (p Pointer) ReadStruct16(bitoff int) (ret uint16) {
	if q, err := p.Deref(); err == nil && uint(bitoff) < uint(q.dataBits) {
		ret = little16(q.Segment.Data[(uint(q.off)+uint(bitoff))/8:])
	}
	return
}

func (p Pointer) ReadStruct32(bitoff int) (ret uint32) {
	if q, err := p.Deref(); err == nil && uint(bitoff) < uint(q.dataBits) {
		ret = little32(q.Segment.Data[(uint(q.off)+uint(bitoff))/8:])
	}
	return
}

func (p Pointer) ReadStruct64(bitoff int) (ret uint64) {
	if q, err := p.Deref(); err == nil && uint(bitoff) < uint(q.dataBits) {
		ret = little64(q.Segment.Data[(uint(q.off)+uint(bitoff))/8:])
	}
	return
}

func (p Pointer) WriteStruct1(bitoff int, v bool) error {
	q, err := p.Deref()
	if err != nil {
		return err
	} else if uint(bitoff)+1 > uint(q.dataBits) {
		return ErrOutOfBounds
	}
	off := uint(q.off) + uint(bitoff)
	if v {
		q.Segment.Data[off/8] |= 1 << uint(off%8)
	} else {
		q.Segment.Data[off/8] &^= 1 << uint(off%8)
	}
	return nil
}

func (p Pointer) WriteStruct8(bitoff int, v uint8) error {
	q, err := p.Deref()
	if err != nil {
		return err
	} else if uint(bitoff)+8 > uint(q.dataBits) {
		return ErrOutOfBounds
	}
	q.Segment.Data[(uint(q.off)+uint(bitoff))/8] = v
	return nil
}

func (p Pointer) WriteStruct16(bitoff int, v uint16) error {
	q, err := p.Deref()
	if err != nil {
		return err
	} else if uint(bitoff)+16 > uint(q.dataBits) {
		return ErrOutOfBounds
	}
	putLittle16(q.Segment.Data[(uint(q.off)+uint(bitoff))/8:], v)
	return nil
}

func (p Pointer) WriteStruct32(bitoff int, v uint32) error {
	q, err := p.Deref()
	if err != nil {
		return err
	} else if uint(bitoff)+32 > uint(q.dataBits) {
		return ErrOutOfBounds
	}
	putLittle32(q.Segment.Data[(uint(q.off)+uint(bitoff))/8:], v)
	return nil
}

func (p Pointer) WriteStruct64(bitoff int, v uint64) error {
	q, err := p.Deref()
	if err != nil {
		return err
	} else if uint(bitoff)+64 > uint(q.dataBits) {
		return ErrOutOfBounds
	}
	putLittle64(q.Segment.Data[(uint(q.off)+uint(bitoff))/8:], v)
	return nil
}

func (p Pointer) listData(i int, sz uint32) ([]byte, uint32, error) {
	q, err := p.Deref()
	if err != nil {
		return nil, 0, err
	} else if uint(i) >= uint(q.size) {
		return nil, 0, ErrOutOfBounds
	}

	switch q.typ {
	case PointerList:
		m := q.ReadPtr(i)
		if sz > m.dataBits {
			return nil, 0, ErrOutOfBounds
		}
		return m.Segment.Data[m.off/8:], m.off, nil
	case List:
		if sz > q.dataBits {
			return nil, 0, ErrOutOfBounds
		}
		off := q.off + uint32(i)*(q.dataBits+uint32(q.ptrs)*64)
		return q.Segment.Data[off/8:], off, nil
	case Struct:
		panic("capn: Read/Write list value on a struct")
	default:
		panic("unhandled")
	}
}

func (p Pointer) Read1(i int) (ret bool) {
	if data, off, err := p.listData(i, 1); err == nil {
		ret = data[0]&(1<<uint(off%8)) != 0
	}
	return
}

func (p Pointer) Read8(i int) (ret uint8) {
	if data, _, err := p.listData(i, 8); err == nil {
		ret = data[0]
	}
	return
}

func (p Pointer) Read16(i int) (ret uint16) {
	if data, _, err := p.listData(i, 16); err == nil {
		ret = little16(data)
	}
	return
}

func (p Pointer) Read32(i int) (ret uint32) {
	if data, _, err := p.listData(i, 32); err == nil {
		ret = little32(data)
	}
	return
}

func (p Pointer) Read64(i int) (ret uint64) {
	if data, _, err := p.listData(i, 64); err == nil {
		ret = little64(data)
	}
	return
}

func (p Pointer) Write1(i int, v bool) error {
	if data, off, err := p.listData(i, 1); err != nil {
		return err
	} else if v {
		data[0] |= 1 << uint(off%8)
	} else {
		data[0] &^= 1 << uint(off%8)
	}
	return nil
}

func (p Pointer) Write8(i int, v uint8) error {
	if data, _, err := p.listData(i, 16); err != nil {
		return err
	} else {
		data[0] = v
		return nil
	}
}

func (p Pointer) Write16(i int, v uint16) error {
	if data, _, err := p.listData(i, 16); err != nil {
		return err
	} else {
		putLittle16(data, v)
		return nil
	}
}

func (p Pointer) Write32(i int, v uint32) error {
	if data, _, err := p.listData(i, 32); err != nil {
		return err
	} else {
		putLittle32(data, v)
		return nil
	}
}

func (p Pointer) Write64(i int, v uint64) error {
	if data, _, err := p.listData(i, 64); err != nil {
		return err
	} else {
		putLittle64(data, v)
		return nil
	}
}

func (s *Segment) check() error {
	if (len(s.Data) % 8) != 0 {
		return ErrAlignment
	}
	if len(s.Data) > int(math.MaxUint32/64) {
		return ErrOverlarge
	}
	return nil
}

func (s *Segment) create(n Pointer) (Pointer, error) {
	if err := s.check(); err != nil {
		return Pointer{}, err
	}

	// switch to 64 bit multiply to overvoid overflow
	sz := (uint64(n.size)*uint64(n.dataBits+uint32(n.ptrs)*64) + 63) &^ 63
	if sz > uint64(math.MaxUint32)-64 {
		return Pointer{}, ErrOverlarge
	}

	tag := false
	if len(s.Data)+int(sz/8) > cap(s.Data) {
		// If we can't fit the data in the current segment, we always
		// return a far pointer to a tag in the new segment.
		if n.typ != Root {
			sz += 64
		}
		news, err := s.Session.NewSegment(int(sz / 8))
		if err != nil {
			return Pointer{}, err
		}

		if err := news.check(); err != nil {
			return Pointer{}, err
		}

		// NewSegment is allowed to grow the segment and return it. In
		// which case we don't want to create a far pointer.
		if n.typ != Root {
			if news == s {
				sz -= 64
			} else {
				tag = true
			}
		}

		s = news
	}

	n.Segment = s
	n.off = uint32(len(s.Data) * 8)
	s.Data = s.Data[:len(s.Data)+int(sz/8)]

	if tag {
		t := s.makeFarPointer(s.Id, n.off)
		n.off += 64
		putLittle64(s.Data[t.off/8:], n.value(t.off))
		return t, nil
	} else {
		return n, nil
	}
}

const (
	structPointer = 0
	listPointer   = 1
	farPointer    = 2

	voidList      = 0
	bit1List      = 1
	byte1List     = 2
	byte2List     = 3
	byte4List     = 4
	byte8List     = 5
	pointerList   = 6
	compositeList = 7
)

func (s *Segment) rawReadPtr(off uint32, val uint64) Pointer {
	// Be wary of overflow. Offset is 30 bits signed. List size is 29 bits
	// unsigned. For both of these we need to checks in terms of words if
	// using 32 bit maths as bits or bytes will overflow.
	switch val & 3 {
	case structPointer:
		offw := int32(off/64) + 1 + int32(uint32(val))>>2
		if offw < 0 || int(offw) >= len(s.Data)/8 {
			return Pointer{}
		}

		p := Pointer{
			Segment:  s,
			typ:      Struct,
			off:      uint32(offw) * 64,
			size:     1,
			dataBits: uint32(uint16(val>>32)) * 64,
			ptrs:     uint16(val >> 48),
		}

		if p.end() > uint32(len(s.Data)) {
			return Pointer{}
		}

		return p

	case listPointer:
		offw := int32(off/64) + 1 + int32(uint32(val))>>2
		if offw < 0 || int(offw) >= len(s.Data)/8 {
			return Pointer{}
		}

		p := Pointer{
			Segment: s,
			typ:     List,
			off:     uint32(offw) * 64,
			size:    uint32(val >> 35),
		}

		words := p.size

		switch (val >> 35) & 7 {
		case bit1List:
			p.dataBits = 1
			words = (p.size + 63) / 64
		case byte1List:
			p.dataBits = 8
			words = (p.size + 7) / 8
		case byte2List:
			p.dataBits = 16
			words = (p.size + 3) / 4
		case byte4List:
			p.dataBits = 32
			words = (p.size + 1) / 2
		case byte8List:
			p.dataBits = 64
			words = p.size
		case pointerList:
			p.ptrs = 1
			p.typ = PointerList
			words = p.size
		case compositeList:
			hdr := little64(p.Segment.Data[p.off/8:])
			p.off += 64
			if hdr&2 != structPointer {
				return Pointer{}
			}

			p.size = uint32(uint32(hdr >> 2))
			p.dataBits = uint32(uint16(hdr>>32)) * 64
			p.ptrs = uint16(hdr >> 48)

			if p.size*(p.dataBits/64+uint32(p.ptrs)) != words {
				return Pointer{}
			}
		}

		// Largest possible message is 30 bits * 1 word, with either a
		// composite, ptr, or 8 byte list. If we do a size check using
		// bits or bytes, we overflow.
		if words > uint32(len(s.Data))/8-uint32(offw) {
			return Pointer{}
		}

		return p

	case farPointer:
		if (val & 4) != 0 {
			return s.makeDoublePointer(uint32(val>>32), uint32(val>>3)*64)
		} else {
			return s.makeFarPointer(uint32(val>>32), uint32(val>>3)*64)
		}

	default:
		return Pointer{}
	}
}

func (p Pointer) ReadPtr(i int) Pointer {
	switch p.typ {
	case Root:
		if uint(i) >= uint(p.size) {
			return Pointer{}
		}

		off := p.off + uint32(i)*64
		n := p.Segment.rawReadPtr(off, little64(p.Segment.Data[off/8:]))
		n, _ = n.Deref()
		return n
	}

	q, err := p.Deref()
	if err != nil {
		return Pointer{}
	}

	switch q.typ {
	case Struct:
		if uint(i) < uint(q.ptrs) {
			off := q.off + q.dataBits + uint32(i)*64
			n := q.Segment.rawReadPtr(off, little64(q.Segment.Data[off/8:]))
			n, _ = n.Deref()
			return n
		}

	case List:
		if uint(i) < uint(q.size) {
			return Pointer{
				Segment:  q.Segment,
				typ:      Struct,
				off:      q.off + uint32(i)*(q.dataBits+uint32(q.ptrs)*64),
				size:     1,
				dataBits: q.dataBits,
				ptrs:     q.ptrs,
			}
		}
	case PointerList:
		if uint(i) < uint(q.size) {
			off := q.off + uint32(i)*64
			n := q.Segment.rawReadPtr(off, little64(q.Segment.Data[off/8:]))
			n, _ = n.Deref()
			return n
		}
	}

	return Pointer{}
}

func (p Pointer) farSegmentId() uint32 { return p.size }

func (s *Segment) makeFarPointer(id, off uint32) Pointer {
	return Pointer{Segment: s, typ: FarPointer, dataBits: 8, off: off, size: id}
}

func (s *Segment) makeDoublePointer(id, off uint32) Pointer {
	return Pointer{Segment: s, typ: DoubleFarPointer, dataBits: 16, off: off, size: id}
}

func (p Pointer) lookupSegment() (Pointer, error) {
	if p.Segment.Id != p.farSegmentId() {
		seg, err := p.Segment.Session.Lookup(p.farSegmentId())
		if err != nil {
			return Pointer{}, err
		}
		seg.check()
		p.Segment = seg
	}

	if p.off+p.dataBits > uint32(len(p.Segment.Data)) {
		return Pointer{}, ErrInvalidPointer
	}

	return p, nil
}

func (p Pointer) Deref() (Pointer, error) {
	var err error

	switch p.typ {
	case FarPointer:
		if p, err = p.lookupSegment(); err != nil {
			return Pointer{}, err
		}

		val := little64(p.Segment.Data[p.off/8:])
		if val&3 == farPointer {
			return Pointer{}, ErrInvalidPointer
		}

		if p = p.Segment.rawReadPtr(p.off, val); p.typ == Null {
			return Pointer{}, ErrInvalidPointer
		}

		return p, nil

	case DoubleFarPointer:
		// A double far pointer points to a double pointer, where the
		// first points to the actual data, and the second is the tag
		// that would normally be placed right before the data (offset
		// == 0).

		if p, err = p.lookupSegment(); err != nil {
			return Pointer{}, err
		}

		tag := little64(p.Segment.Data[(p.off+64)/8:])
		p = p.Segment.rawReadPtr(p.off, little64(p.Segment.Data[p.off/8:]))

		// The low bits are the type (0 struct, 1 list) followed by
		// the offset. So we can quickly check for a struct/list with
		// 0 offset.
		if p.typ != FarPointer || uint32(tag) > listPointer {
			return Pointer{}, ErrInvalidPointer
		}

		if p, err = p.lookupSegment(); err != nil {
			return Pointer{}, err
		}

		if p = p.Segment.rawReadPtr(p.off-64, tag); p.typ == Null {
			return Pointer{}, ErrInvalidPointer
		}

		return p, nil

	case Struct, PointerList, List, Root:
		return p, nil
	case Null:
		return p, ErrInvalidPointer
	default:
		panic("unhandled")
	}
}

func (p Pointer) value(off uint32) uint64 {
	d := uint64(uint32(int32(p.off)-int32(off)-64) / 64 << 2)

	switch p.typ {
	case FarPointer:
		return farPointer | uint64(p.off/64)<<3 | uint64(p.farSegmentId())<<32
	case DoubleFarPointer:
		return farPointer | 4 | uint64(p.off/64)<<3 | uint64(p.farSegmentId())<<32
	case Struct:
		return structPointer | d | uint64(p.dataBits/64)<<32 | uint64(p.ptrs)<<48
	case PointerList:
		return listPointer | d | pointerList<<32 | uint64(p.size)<<35
	case List:
		if p.ptrs > 0 || p.dataBits > 64 {
			return listPointer | d | compositeList<<32 | uint64(p.size)<<35
		}

		switch p.dataBits {
		case 0:
			return listPointer | d | voidList<<32 | uint64(p.size)<<35
		case 1:
			return listPointer | d | bit1List<<32 | uint64(p.size)<<35
		case 8:
			return listPointer | d | byte1List<<32 | uint64(p.size)<<35
		case 16:
			return listPointer | d | byte2List<<32 | uint64(p.size)<<35
		case 32:
			return listPointer | d | byte4List<<32 | uint64(p.size)<<35
		case 64:
			return listPointer | d | byte8List<<32 | uint64(p.size)<<35
		}
	case Null:
		return 0
	}
	println(p.typ, p.dataBits, p.ptrs)
	panic("unhandled")
}

func (s *Segment) rawWritePtr(off uint32, tgt Pointer) error {
	if tgt.typ == Null {
		putLittle64(s.Data[off/8:], 0)
		return nil
	}

	q, err := tgt.Deref()
	if err != nil {
		return err
	}

	if s == q.Segment {
		// Same segment
		putLittle64(s.Data[off/8:], q.value(off))
		return nil

	} else if s.Session != q.Segment.Session {
		// Different session - deep copy
		var err error
		var n Pointer

		switch q.typ {
		case Struct:
			n, err = s.NewStruct(int(q.dataBits), int(q.ptrs))
		case List:
			n, err = s.NewList(int(q.dataBits), int(q.ptrs), int(q.size))
		case PointerList:
			n, err = s.NewPointerList(int(q.size))
		default:
			panic("unhandled")
		}

		if err != nil {
			return err
		}

		putLittle64(s.Data[off/8:], n.value(off))
		return Copy(n, q)

	} else if tgt.typ == FarPointer || tgt.typ == DoubleFarPointer {
		// Already have a far pointer to point to
		putLittle64(s.Data[off/8:], tgt.value(0))
		return nil

	} else if qs := q.Segment; len(qs.Data)+8 <= cap(qs.Data) {
		// Have room in the target for a tag
		far := qs.makeFarPointer(qs.Id, uint32(len(qs.Data))*8)
		qs.Data = qs.Data[:len(qs.Data)+8]

		putLittle64(qs.Data[far.off/8:], q.value(far.off))
		putLittle64(s.Data[off/8:], far.value(0))
		return nil

	} else {
		// Need to create a double far pointer
		t := s
		if len(t.Data)+16 > cap(t.Data) {
			news, err := t.Session.NewSegment(16)
			if err != nil {
				return err
			}

			if err := news.check(); err != nil {
				return err
			}

			t = news
		}

		far1 := t.makeFarPointer(q.Segment.Id, q.off)
		far2 := t.makeDoublePointer(t.Id, uint32(len(t.Data))*8)

		s.Data = s.Data[:len(s.Data)+16]

		putLittle64(t.Data[far2.off/8:], far1.value(0))
		putLittle64(t.Data[far2.off/8+8:], q.value(q.off+64))
		putLittle64(s.Data[off/8:], far2.value(0))
		return nil
	}
}

func (p Pointer) MarshalCaptain(in Pointer, off int) error {
	return in.WritePtr(off, p)
}

func (p Pointer) WritePtr(i int, tgt Pointer) error {
	switch p.typ {
	case Root:
		if uint(i) >= uint(p.size) {
			return ErrOutOfBounds
		}
		off := p.off + uint32(i)*64
		return p.Segment.rawWritePtr(off, tgt)
	}

	q, err := p.Deref()
	if err != nil {
		return err
	}

	switch q.typ {
	case Root:
	case Struct:
		if uint(i) < uint(q.ptrs) {
			return q.Segment.rawWritePtr(q.off+q.dataBits+64*uint32(i), tgt)
		}
	case List:
		if uint(i) < uint(q.size) {
			return Copy(q.ReadPtr(i), tgt)
		}
	case PointerList:
		if uint(i) < uint(q.size) {
			return q.Segment.rawWritePtr(q.off+64*uint32(i), tgt)
		}
	default:
		panic("unhandled")
	}
	return ErrOutOfBounds
}

func copyData(to, from Pointer) {
	td := to.Data()
	fd := from.Data()
	for i := range td[copy(td, fd):] {
		td[i] = 0
	}
}

func copyPtrs(to, from Pointer) error {
	ptrs := int(from.ptrs)
	if ptrs > int(to.ptrs) {
		ptrs = int(to.ptrs)
	}

	for i := 0; i < ptrs; i++ {
		if err := to.WritePtr(i, from.ReadPtr(i)); err != nil {
			return err
		}
	}

	for i := ptrs; i < int(to.ptrs); i++ {
		to.WritePtr(i, Pointer{})
	}

	return nil
}

func Copy(to, from Pointer) error {
	var err error

	if to, err = to.Deref(); err != nil {
		return err
	}
	if from, err = from.Deref(); err != nil {
		return err
	}

	switch to.typ {
	case Struct:
		copyData(to, from)
		return copyPtrs(to, from)
	case List:
		if from.ptrs == 0 && to.ptrs == 0 && to.dataBits == from.dataBits {
			copyData(to, from)
			return nil
		} else {
			return copyPtrs(to, from)
		}
	case PointerList:
		return copyPtrs(to, from)
	case Null:
		return nil
	default:
		panic("unhandled")
	}
}
