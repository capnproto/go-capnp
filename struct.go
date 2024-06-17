package capnp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"unsafe"

	"capnproto.org/go/capnp/v3/exc"
	"capnproto.org/go/capnp/v3/internal/str"
)

// Struct is a pointer to a struct.
type Struct StructKind

// The underlying type of Struct. We expose this so that
// we can use ~StructKind as a constraint in generics to
// capture any struct type.
type StructKind = struct {
	seg        *Segment
	off        address
	size       ObjectSize
	depthLimit uint
	flags      structFlags
}

// AllocateRootStruct allocates the root struct on the message as a struct of
// the passed size.
func AllocateRootStruct(msg *Message, sz ObjectSize) (Struct, error) {
	if !sz.isValid() {
		return Struct{}, errors.New("new struct: invalid size")
	}
	sz.DataSize = sz.DataSize.padToWord()
	s, addr, err := msg.AllocateAsRoot(sz)
	if err != nil {
		return Struct{}, exc.WrapError("new struct", err)
	}
	return Struct{
		seg:        s,
		off:        addr,
		size:       sz,
		depthLimit: maxDepth,
	}, nil
}

// NewStruct creates a new struct, preferring placement in s.
func NewStruct(s *Segment, sz ObjectSize) (Struct, error) {
	if !sz.isValid() {
		return Struct{}, errors.New("new struct: invalid size")
	}
	sz.DataSize = sz.DataSize.padToWord()
	seg, addr, err := alloc(s, sz.totalSize())
	if err != nil {
		return Struct{}, exc.WrapError("new struct", err)
	}
	return Struct{
		seg:        seg,
		off:        addr,
		size:       sz,
		depthLimit: maxDepth,
	}, nil
}

// NewRootStruct creates a new struct, preferring placement in s, then sets the
// message's root to the new struct.
func NewRootStruct(s *Segment, sz ObjectSize) (Struct, error) {
	st, err := NewStruct(s, sz)
	if err != nil {
		return st, err
	}
	if err := s.msg.SetRoot(st.ToPtr()); err != nil {
		return st, err
	}
	return st, nil
}

// ToPtr converts the struct to a generic pointer.
func (p Struct) ToPtr() Ptr {
	return Ptr{
		seg:        p.seg,
		off:        p.off,
		size:       p.size,
		depthLimit: p.depthLimit,
		flags:      structPtrFlag(p.flags),
	}
}

// Segment returns the segment the referenced struct is stored in or nil
// if the pointer is invalid.
func (p Struct) Segment() *Segment {
	return p.seg
}

// Message returns the message the referenced struct is stored in or nil
// if the pointer is invalid.
func (p Struct) Message() *Message {
	if p.seg == nil {
		return nil
	}
	return p.seg.msg
}

// IsValid returns whether the struct is valid.
func (p Struct) IsValid() bool {
	return p.seg != nil
}

// Size returns the size of the struct.
func (p Struct) Size() ObjectSize {
	return p.size
}

// CopyFrom copies content from another struct.  If the other struct's
// sections are larger than this struct's, the extra data is not copied,
// meaning there is a risk of data loss when copying from messages built
// with future versions of the protocol.
func (p Struct) CopyFrom(other Struct) error {
	if err := copyStruct(p, other); err != nil {
		return exc.WrapError("copy struct", err)
	}
	return nil
}

// readSize returns the struct's size for the purposes of read limit
// accounting.
func (p Struct) readSize() Size {
	if p.seg == nil {
		return 0
	}
	return p.size.totalSize()
}

// Ptr returns the i'th pointer in the struct.
func (p Struct) Ptr(i uint16) (Ptr, error) {
	if p.seg == nil || i >= p.size.PointerCount {
		return Ptr{}, nil
	}
	return p.seg.readPtr(p.pointerAddress(i), p.depthLimit)
}

// HasPtr reports whether the i'th pointer in the struct is non-null.
// It does not affect the read limit.
func (p Struct) HasPtr(i uint16) bool {
	if p.seg == nil || i >= p.size.PointerCount {
		return false
	}
	return p.seg.readRawPointer(p.pointerAddress(i)) != 0
}

// SetPtr sets the i'th pointer in the struct to src.
func (p Struct) SetPtr(i uint16, src Ptr) error {
	if p.seg == nil || i >= p.size.PointerCount {
		panic("capnp: set field outside struct boundaries")
	}
	return p.seg.writePtr(p.pointerAddress(i), src, false)
}

// SetText sets the i'th pointer to a newly allocated text or null if v is empty.
func (p Struct) SetText(i uint16, v string) error {
	if v == "" {
		return p.SetPtr(i, Ptr{})
	}
	return p.SetNewText(i, v)
}

func (p *Struct) GetText(i uint16) (string, error) {
	// p.Ptr(i)
	if p.seg == nil || i >= p.size.PointerCount {
		return "", nil
	}

	// p.pointerAddress(i)
	ptrStart, _ := p.off.addSize(p.size.DataSize)
	addr, _ := ptrStart.element(int32(i), wordSize)

	ptr, err := p.seg.readPtr(addr, p.depthLimit)
	if err != nil {
		return "", err
	}

	// return ptr.Text(), nil

	tb, ok := ptr.textp()
	if !ok {
		return "", nil
	}

	return string(tb), nil
}

func (p *Struct) GetTextUnsafe(i uint16) (string, error) {
	// p.Ptr(i)
	if p.seg == nil || i >= p.size.PointerCount {
		return "", nil
	}

	// p.pointerAddress(i)
	ptrStart, _ := p.off.addSize(p.size.DataSize)   // Start of pointers in struct (could be cached)
	addr, _ := ptrStart.element(int32(i), wordSize) // Address of ith pointer

	// ptr, err := p.seg.readPtr(addr, p.depthLimit)
	s, base, val, err := p.seg.resolveFarPointer(addr) // Read pointer (panics if out of bounds, but validated by initial check)

	if err != nil {
		return "", exc.WrapError("read pointer", err)
	}
	if val == 0 {
		return "", nil
	}
	if p.depthLimit == 0 {
		return "", errors.New("read pointer: depth limit reached")
	}
	if val.pointerType() != listPointer {
		return "", fmt.Errorf("not a list pointer")
	}
	/*
		lp, err := s.readListPtr(base, val)
		if err != nil {
			return "", exc.WrapError("read pointer", err)
		}
	*/
	if val.listType() != byte1List {
		return "", fmt.Errorf("not a byte1list")
	}
	laddr, ok := val.offset().resolve(base)
	if !ok {
		return "", errors.New("list pointer: invalid address")
	}
	/*
		if !s.msg.canRead(lp.readSize()) {
			return "", errors.New("read pointer: read traversal limit reached")
		}
		lp.depthLimit = p.depthLimit - 1
	*/
	// return ptr.Text(), nil

	// tb, ok := ptr.textp()
	tb := s.slice(laddr, Size(val.numListElements()))
	end := len(tb) - 1
	trimmed := false
	for ; end >= 0 && tb[end] == 0; end-- {
		trimmed = true
	}
	if !trimmed {
		// Not null terminated.
		return "", nil
	}
	tb = tb[:end+1]

	return *(*string)(unsafe.Pointer(&tb)), nil
}

func (p *Struct) GetTextSuperUnsafe(i uint16) (string, error) {
	// p.Ptr(i)
	// Bounds check pointer address within range.
	if p.seg == nil || i >= p.size.PointerCount {
		return "", nil
	}

	// p.pointerAddress(i)
	ptrStart := p.off + address(p.size.DataSize)
	addr := ptrStart + address(i)*address(wordSize)

	// ptr, err := p.seg.readPtr(addr, p.depthLimit)
	s, base, val, _ := p.seg.resolveFarPointer(addr) // Read pointer (panics if out of bounds, but validated by initial check)
	/*
		if err != nil {
			return "", exc.WrapError("read pointer", err)
		}
	*/
	if val == 0 {
		return "", nil
	}
	/*
		if p.depthLimit == 0 {
			return "", errors.New("read pointer: depth limit reached")
		}
		if val.pointerType() != listPointer {
			return "", fmt.Errorf("not a list pointer")
		}
		/*
			lp, err := s.readListPtr(base, val)
			if err != nil {
				return "", exc.WrapError("read pointer", err)
			}
	*/
	/*
		if val.listType() != byte1List {
			return "", fmt.Errorf("not a byte1list")
		}
	*/
	laddr := val.offset().resolveUnsafe(base)
	/*
		laddr, ok := val.offset().resolve(base)
		if !ok {
			return "", errors.New("list pointer: invalid address")
		}
	*/
	/*
		if !s.msg.canRead(lp.readSize()) {
			return "", errors.New("read pointer: read traversal limit reached")
		}
		lp.depthLimit = p.depthLimit - 1
	*/
	// return ptr.Text(), nil

	// tb, ok := ptr.textp()
	tb := s.slice(laddr, Size(val.numListElements()))
	end := len(tb) - 1
	trimmed := false
	for ; end >= 0 && tb[end] == 0; end-- {
		trimmed = true
	}
	if !trimmed {
		// Not null terminated.
		return "", nil
	}
	tb = tb[:end+1]

	return *(*string)(unsafe.Pointer(&tb)), nil
}

func (p Struct) FlatSetNewText(i uint16, v string) error {
	// NewText().newPrimitiveList
	sz := Size(1)
	n := int32(len(v) + 1)
	total := sz.timesUnchecked(n)
	s, addr, err := alloc(p.seg, total)
	if err != nil {
		return err
	}

	// NewText()
	copy(s.slice(addr, Size(len(v))), v)

	// SetPtr().pointerAddress()
	// offInsideP := p.pointerAddress(i)
	ptrStart, _ := p.off.addSize(p.size.DataSize)
	offInsideP, _ := ptrStart.element(int32(i), wordSize)

	// SetPtr().writePtr()
	srcAddr := addr                                       // srcAddr = l.off
	srcRaw := rawListPointer(0, byte1List, int32(len(v))) // srcRaw = l.raw()
	s.writeRawPointer(offInsideP, srcRaw.withOffset(nearPointerOffset(offInsideP, srcAddr)))
	return nil
}

// EXPERIMENTAL unrolled version of updating a text field in-place when it
// already has enough space to hold the value (as opposed to allocating a new
// object in the message).
func (p Struct) UpdateText(i uint16, v string) error {
	// Determine pointer offset.
	ptrStart, _ := p.off.addSize(p.size.DataSize)
	offInsideP, _ := ptrStart.element(int32(i), wordSize)

	// ptr, err := p.seg.readPtr(offInsideP, p.depthLimit)
	s, base, val, err := p.seg.resolveFarPointer(offInsideP)
	if err != nil {
		return exc.WrapError("read pointer", err)
	}

	// TODO: depth limit read check?

	if val == 0 {
		// Existing pointer is empty/void.
		return p.FlatSetNewText(i, v)
	}

	// lp, err := s.readListPtr(base, val)
	addr, ok := val.offset().resolve(base)
	if !ok {
		return errors.New("list pointer: invalid address")
	}
	// TODO: list checks from readListPtr()?

	if err != nil {
		return exc.WrapError("read pointer", err)
	}

	length := int(val.numListElements())

	if length < len(v)+1 {
		// Existing buffer does not have enough space for new text.
		return p.FlatSetNewText(i, v)
	}

	// Existing buffer location has space for new text. Copy text over it.
	dst := s.slice(addr, Size(length))
	n := copy(dst, []byte(v))

	// Pad with zeros (clear leftover). Last byte is already zero.
	//
	// TODO: replace with clear(dst[n:length-1]) after go1.21.
	for i := n; i < int(length-1); i++ {
		dst[i] = 0
	}

	return nil
}

type TextField struct {
	// Pointer location
	pSeg  *Segment
	pAddr address

	// Current value location
	vSeg  *Segment
	vAddr address
	vLen  int
}

// EXPERIMENTAL: return ith pointer as a text field.
func (p Struct) TextField(i uint16) (TextField, error) {
	ptrStart, _ := p.off.addSize(p.size.DataSize)
	offInsideP, _ := ptrStart.element(int32(i), wordSize)

	// ptr, err := p.seg.readPtr(offInsideP, p.depthLimit)
	s, base, val, err := p.seg.resolveFarPointer(offInsideP)
	if err != nil {
		return TextField{}, exc.WrapError("read pointer", err)
	}

	tf := TextField{pSeg: p.seg, pAddr: offInsideP}

	if val == 0 {
		return tf, nil
	}

	addr, ok := val.offset().resolve(base)
	if !ok {
		return TextField{}, errors.New("list pointer: invalid address")
	}

	tf.vSeg = s
	tf.vLen = int(val.numListElements())
	tf.vAddr = addr

	return tf, nil
}

// UpdateText updates the value of the text field.
func (tf *TextField) Set(v string) error {
	if tf.vLen < len(v)+1 || tf.vSeg == nil {
		// TODO: handle this case. Needs to alloc and set pointer.
		// Needs to set tf.vSeg, tf.vLen and tf.vAddr.
		panic("we can work it out")
	}

	// Existing buffer location has space for new text. Copy text over it.
	dst := tf.vSeg.slice(tf.vAddr, Size(tf.vLen))
	n := copy(dst, []byte(v))

	// Pad with zeros (clear leftover). Last byte is already zero.
	//
	// TODO: replace with clear(dst[n:length-1]) after go1.21.
	for i := n; i < int(tf.vLen-1); i++ {
		dst[i] = 0
	}

	return nil
}

func trimZero(r rune) bool {
	return r == 0
}

func (tf *TextField) Get() string {
	if tf.vSeg == nil {
		panic("not allocated")
	}

	if tf.vLen == 0 {
		return ""
	}

	b := tf.vSeg.slice(tf.vAddr, Size(tf.vLen))
	return string(bytes.TrimRightFunc(b, trimZero))
}

func (tf *TextField) GetUnsafe() string {
	b := tf.vSeg.slice(tf.vAddr, Size(tf.vLen))
	b = bytes.TrimRightFunc(b, trimZero)
	return *(*string)(unsafe.Pointer(&b))
}

// SetNewText sets the i'th pointer to a newly allocated text.
func (p Struct) SetNewText(i uint16, v string) error {
	t, err := NewText(p.seg, v)
	if err != nil {
		return err
	}
	return p.SetPtr(i, t.ToPtr())
}

// SetTextFromBytes sets the i'th pointer to a newly allocated text or null if v is nil.
func (p Struct) SetTextFromBytes(i uint16, v []byte) error {
	if v == nil {
		return p.SetPtr(i, Ptr{})
	}
	t, err := NewTextFromBytes(p.seg, v)
	if err != nil {
		return err
	}
	return p.SetPtr(i, t.ToPtr())
}

// SetData sets the i'th pointer to a newly allocated data or null if v is nil.
func (p Struct) SetData(i uint16, v []byte) error {
	if v == nil {
		return p.SetPtr(i, Ptr{})
	}
	d, err := NewData(p.seg, v)
	if err != nil {
		return err
	}
	return p.SetPtr(i, d.ToPtr())
}

func (p Struct) pointerAddress(i uint16) address {
	// Struct already had bounds check
	ptrStart, _ := p.off.addSize(p.size.DataSize)
	a, _ := ptrStart.element(int32(i), wordSize)
	return a
}

// bitInData reports whether bit is inside p's data section.
func (p *Struct) bitInData(bit BitOffset) bool {
	return p.seg != nil && bit < BitOffset(p.size.DataSize*8)
}

// Bit returns the bit that is n bits from the start of the struct.
func (p Struct) Bit(n BitOffset) bool {
	if !p.bitInData(n) {
		return false
	}
	addr := p.off.addOffset(n.offset())
	return p.seg.readUint8(addr)&n.mask() != 0
}

func (p *Struct) Bitp(n BitOffset) bool {
	addr := p.dataAddressUnchecked(n.offset())
	return p.seg.readUint8(addr)&n.mask() != 0
}

// SetBit sets the bit that is n bits from the start of the struct to v.
func (p Struct) SetBit(n BitOffset, v bool) {
	if !p.bitInData(n) {
		panic("capnp: set field outside struct boundaries")
	}
	addr := p.off.addOffset(n.offset())
	b := p.seg.readUint8(addr)
	if v {
		b |= n.mask()
	} else {
		b &^= n.mask()
	}
	p.seg.writeUint8(addr, b)
}

func (p *Struct) SetBitp(n BitOffset, v bool) {
	addr := p.dataAddressUnchecked(n.offset())
	if v {
		p.seg.setBit(addr, uint8(n%8))
	} else {
		p.seg.clearBit(addr, uint8(n%8))
	}
	/*
		if !p.bitInData(n) {
			panic("capnp: set field outside struct boundaries")
		}
		addr := p.off.addOffset(n.offset())
		b := p.seg.readUint8(addr)
		if v {
			b |= n.mask()
		} else {
			b &^= n.mask()
		}
		p.seg.writeUint8(addr, b)
	*/
}

func (p *Struct) dataAddress(off DataOffset, sz Size) (addr address, ok bool) {
	if p.seg == nil || Size(off)+sz > p.size.DataSize {
		return 0, false
	}
	return p.off.addOffset(off), true
	// return p.off + address(off), true
}

func (p *Struct) dataAddressUnchecked(off DataOffset) (addr address) {
	return p.off + address(off)
}

// Uint8 returns an 8-bit integer from the struct's data section.
func (p Struct) Uint8(off DataOffset) uint8 {
	addr, ok := p.dataAddress(off, 1)
	if !ok {
		return 0
	}
	return p.seg.readUint8(addr)
}

// Uint16 returns a 16-bit integer from the struct's data section.
func (p Struct) Uint16(off DataOffset) uint16 {
	addr, ok := p.dataAddress(off, 2)
	if !ok {
		return 0
	}
	return p.seg.readUint16(addr)
}

// Uint32 returns a 32-bit integer from the struct's data section.
func (p Struct) Uint32(off DataOffset) uint32 {
	addr, ok := p.dataAddress(off, 4)
	if !ok {
		return 0
	}
	return p.seg.readUint32(addr)
}

// Uint32p returns a 32-bit integer from the struct's data section.
func (p *Struct) Uint32p(off DataOffset) uint32 {
	addr, ok := p.dataAddress(off, 4)
	if !ok {
		return 0
	}
	return p.seg.readUint32(addr)
}

// Uint64 returns a 64-bit integer from the struct's data section.
func (p Struct) Uint64(off DataOffset) uint64 {
	addr, ok := p.dataAddress(off, 8)
	if !ok {
		return 0
	}
	return p.seg.readUint64(addr)
}

func (p *Struct) Uint64p(off DataOffset) uint64 {

	/*
		addr, ok := p.dataAddress(off, 8)
		if !ok {
			return 0
		}
	*/
	addr := p.dataAddressUnchecked(off)
	return p.seg.readUint64(addr)
}

// SetUint8 sets the 8-bit integer that is off bytes from the start of the struct to v.
func (p Struct) SetUint8(off DataOffset, v uint8) {
	addr, ok := p.dataAddress(off, 1)
	if !ok {
		panic("capnp: set field outside struct boundaries")
	}
	p.seg.writeUint8(addr, v)
}

// SetUint16 sets the 16-bit integer that is off bytes from the start of the struct to v.
func (p Struct) SetUint16(off DataOffset, v uint16) {
	addr, ok := p.dataAddress(off, 2)
	if !ok {
		panic("capnp: set field outside struct boundaries")
	}
	p.seg.writeUint16(addr, v)
}

// SetUint32 sets the 32-bit integer that is off bytes from the start of the struct to v.
func (p Struct) SetUint32(off DataOffset, v uint32) {
	addr, ok := p.dataAddress(off, 4)
	if !ok {
		panic("capnp: set field outside struct boundaries")
	}
	p.seg.writeUint32(addr, v)
}

func (p *Struct) SetUint32p(off DataOffset, v uint32) {
	addr, ok := p.dataAddress(off, 4)
	if !ok {
		panic("capnp: set field outside struct boundaries")
	}
	p.seg.writeUint32(addr, v)
}

// SetUint64 sets the 64-bit integer that is off bytes from the start of the struct to v.
func (p Struct) SetUint64(off DataOffset, v uint64) {
	addr, ok := p.dataAddress(off, 8)
	if !ok {
		panic("capnp: set field outside struct boundaries")
	}
	p.seg.writeUint64(addr, v)
}

func (p *Struct) SetUint64p(off DataOffset, v uint64) {
	addr, ok := p.dataAddress(off, 8)
	if !ok {
		panic("capnp: set field outside struct boundaries")
	}

	// p.seg.writeUint64(addr, v)
	b := p.seg.slice(addr, 8)
	binary.LittleEndian.PutUint64(b, v)

	/*
		b := p.seg.slice(addr, 8)
		b[0] = byte(v)
		b[1] = byte(v >> 8)
		b[2] = byte(v >> 16)
		b[3] = byte(v >> 24)
		b[4] = byte(v >> 32)
		b[5] = byte(v >> 40)
		b[6] = byte(v >> 48)
		b[7] = byte(v >> 56)
	*/
	/*
		b := p.seg.data
		_ = b[addr+7]
		b[addr] = byte(v)
		b[addr+1] = byte(v >> 8)
		b[addr+2] = byte(v >> 16)
		b[addr+3] = byte(v >> 24)
		b[addr+4] = byte(v >> 32)
		b[addr+5] = byte(v >> 40)
		b[addr+6] = byte(v >> 48)
		b[addr+7] = byte(v >> 56)
	*/

}

func (p *Struct) SetFloat64p(off DataOffset, v float64) {
	/*
		addr, ok := p.dataAddress(off, 8)
		if !ok {
			panic("capnp: set field outside struct boundaries")
		}
	*/
	addr := p.dataAddressUnchecked(off)

	p.seg.writeFloat64(addr, v)
}

// structFlags is a bitmask of flags for a pointer.
type structFlags uint8

// Pointer flags.
const (
	isListMember structFlags = 1 << iota
)

// copyStruct makes a deep copy of src into dst.
func copyStruct(dst, src Struct) error {
	if dst.seg == nil {
		panic("copy struct into invalid pointer")
	}
	if src.seg == nil {
		return nil
	}

	// Q: how does version handling happen here, when the
	//    destination toData[] slice can be bigger or smaller
	//    than the source data slice, which is in
	//    src.seg.Data[src.off:src.off+src.size.DataSize] ?
	//
	// A: Newer fields only come *after* old fields. Note that
	//    copy only copies min(len(src), len(dst)) size,
	//    and then we manually zero the rest in the for loop
	//    that writes toData[j] = 0.
	//

	// data section:
	srcData := src.seg.slice(src.off, src.size.DataSize)
	dstData := dst.seg.slice(dst.off, dst.size.DataSize)
	copyCount := copy(dstData, srcData)
	dstData = dstData[copyCount:]
	for j := range dstData {
		dstData[j] = 0
	}

	// ptrs section:

	// version handling: we ignore any extra-newer-pointers in src,
	// i.e. the case when srcPtrSize > dstPtrSize, by only
	// running j over the size of dstPtrSize, the destination size.
	srcPtrSect, _ := src.off.addSize(src.size.DataSize)
	dstPtrSect, _ := dst.off.addSize(dst.size.DataSize)
	numSrcPtrs := src.size.PointerCount
	numDstPtrs := dst.size.PointerCount
	for j := uint16(0); j < numSrcPtrs && j < numDstPtrs; j++ {
		srcAddr, _ := srcPtrSect.element(int32(j), wordSize)
		dstAddr, _ := dstPtrSect.element(int32(j), wordSize)
		m, err := src.seg.readPtr(srcAddr, src.depthLimit)
		if err != nil {
			return exc.WrapError("copy struct pointer "+str.Utod(j), err)
		}
		err = dst.seg.writePtr(dstAddr, m, true)
		if err != nil {
			return exc.WrapError("copy struct pointer "+str.Utod(j), err)
		}
	}
	for j := numSrcPtrs; j < numDstPtrs; j++ {
		// destination p is a newer version than source so these extra new pointer fields in p must be zeroed.
		addr, _ := dstPtrSect.element(int32(j), wordSize)
		dst.seg.writeRawPointer(addr, 0)
	}
	// Nothing more here: so any other pointers in srcPtrSize beyond
	// those in dstPtrSize are ignored and discarded.

	return nil
}

// s.EncodeAsPtr is equivalent to s.ToPtr(); for implementing TypeParam.
// The segment argument is ignored.
func (s Struct) EncodeAsPtr(*Segment) Ptr { return s.ToPtr() }

// DecodeFromPtr(p) is equivalent to p.Struct() (the receiver is ignored).
// for implementing TypeParam.
func (Struct) DecodeFromPtr(p Ptr) Struct { return p.Struct() }

var _ TypeParam[Struct] = Struct{}
