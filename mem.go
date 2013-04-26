package capn

import (
	"errors"
	"fmt"
)

func min(a, b int) int {
	if b < a {
		return b
	}
	return a
}

var errOutOfBounds = errors.New("capn: supplied offset is out of bounds")

func copyData(to, from Pointer, sz int) error {
	buf := make([]byte, sz)
	if err := from.Read(0, buf); err != nil {
		return err
	}

	if err := to.Write(0, buf); err != nil {
		return err
	}

	return nil
}

func copyPointers(to, from Pointer, sz int) error {
	buf := make([]Pointer, sz)
	if err := from.ReadPtrs(0, buf); err != nil {
		return err
	}

	if err := to.WritePtrs(0, buf); err != nil {
		return err
	}

	return nil
}

func Copy(to, from Pointer) error {
	ft := from.Type()
	tt := to.Type()

	if tt.Type() == CompositeList {
		switch ft.Type() {
		case PointerList, CompositeList:
			for i := 0; i < min(ft.Elements(), tt.Elements()); i++ {
				fp := [1]Pointer{}
				tp := [1]Pointer{}

				if err := from.ReadPtrs(i, fp[:]); err != nil {
					return err
				}

				if err := to.ReadPtrs(i, tp[:]); err != nil {
					return err
				}

				if err := Copy(tp[0], fp[0]); err != nil {
					return err
				}
			}
			return nil
		}
	}

	switch ft.Type() {
	case Struct:
		if tt.Type() == Struct {
			if err := copyData(to, from, min(ft.DataSize(), tt.DataSize())); err != nil {
				return err
			}
			return copyPointers(to, from, min(ft.PointerNum(), tt.PointerNum()))
		}

	case Bit1List, Byte1List, Byte2List, Byte4List, Byte8List:
		if tt.Type() == ft.Type() {
			return copyData(to, from, min(ft.DataSize(), tt.DataSize()))
		}

	case PointerList, CompositeList:
		if tt.Type() == PointerList || tt.Type() == CompositeList {
			return copyPointers(to, from, min(ft.PointerNum(), tt.PointerNum()))
		}
	}

	return errors.New("capn: incompatible copy src/target")
}

func Must(p Pointer, err error) Pointer {
	if err != nil {
		panic(err)
	}
	return p
}

func NewMemory(p PointerType) (Pointer, error) {
	return (*memory)(nil).New(p)
}

type memory struct {
	data []uint8
	ptrs []Pointer
	typ  PointerType
}

func (m *memory) New(p PointerType) (Pointer, error) {
	switch p.Type() {
	case CompositeList:
		ptrs := make([]Pointer, p.Elements())

		for i := range ptrs {
			var err error
			ptrs[i], err = m.New(p.CompositeType())
			if err != nil {
				return nil, err
			}
		}

		return &memory{typ: p, ptrs: ptrs}, nil

	default:
		return &memory{
			typ:  p,
			data: make([]byte, p.DataSize()),
			ptrs: make([]Pointer, p.PointerNum()),
		}, nil
	}
}

func (m *memory) Call(method int, args Pointer) (Pointer, error) {
	return nil, errors.New("capn: memory pointers do not support interfaces")
}

func (m *memory) Read(off int, v []uint8) error {
	copy(v, m.data[off:])
	return nil
}

func (m *memory) Write(off int, v []uint8) error {
	copy(m.data[off:], v)
	return nil
}

func (m *memory) Type() PointerType {
	return m.typ
}

func (m *memory) ReadPtrs(off int, v []Pointer) error {
	copy(v, m.ptrs[off:])
	return nil
}

func (m *memory) WritePtrs(off int, v []Pointer) error {
	copy(m.ptrs[off:], v)
	return nil
}

type Buffer struct {
	buf    []byte
	ptrs   map[Pointer]int
	Caller CallFunc
	Id     uint32
}

// NewBuffer creates a new segment buffer. Data is read/written in wire
// format.
func NewBuffer(buf []byte) *Buffer {
	return &Buffer{
		buf:  buf,
		ptrs: make(map[Pointer]int),
	}
}

func (b *Buffer) Bytes() []byte {
	return b.buf
}

func align8(x int) int {
	return (x + 7) &^ 7
}

// New allocates room for a pointer of type p at the end of the buffer and
// returns a Pointer to manipulate the new data.
func (b *Buffer) New(p PointerType) (Pointer, error) {
	sz := align8(p.DataSize() + p.PointerNum()*8)
	off := len(b.buf)

	if p.Type() == CompositeList {
		b.buf = append(b.buf, make([]byte, 8)...)
		putLittle64(b.buf[off:], p.Composite)

		off += 8
		sz *= p.Elements()
	}

	b.buf = append(b.buf, make([]byte, sz)...)
	println("new", p.Type(), sz, p.Elements(), p.DataSize(), p.PointerNum())
	return bufferPointer{b, off, p}, nil
}

// NewRoot creates a new root value. This writes a pointer tag that points to
// the data immediately after the tag. The returned pointer points to the tag
// and can be written to another segment to create a far pointer.
func (b *Buffer) NewRoot(v PointerType) (Pointer, error) {
	off := len(b.buf)

	v.SetOffset(0)
	b.buf = append(b.buf, make([]byte, 8)...)
	putLittle64(b.buf[off:], v.Value)

	if _, err := b.New(v); err != nil {
		return nil, err
	}

	return bufferPointer{b, off, MakeFarPointer(b.Id, off+8)}, nil
}

func (b *Buffer) ReadRoot(off int) (Pointer, error) {
	return b.readptr(off * 8)
}

type bufferPointer struct {
	seg *Buffer
	off int
	typ PointerType
}

func (p bufferPointer) New(typ PointerType) (Pointer, error) {
	return p.seg.New(typ)
}

func (p bufferPointer) Call(method int, args Pointer) (Pointer, error) {
	if p.seg.Caller == nil {
		return nil, errors.New("capn: no bound caller")
	}
	return p.seg.Caller(p, method, args)
}

func (p bufferPointer) Read(off int, v []uint8) error {
	if err := p.deref(); err != nil {
		return err
	}

	off = p.off + off*8

	if off+len(v) > len(p.seg.buf) {
		panic("out of bounds")
		return errOutOfBounds
	}

	copy(v, p.seg.buf[off:])
	return nil
}

func (p bufferPointer) Write(off int, v []uint8) error {
	if err := p.deref(); err != nil {
		return err
	}

	off = p.off + off*8

	if off+len(v) > len(p.seg.buf) {
		println(off, len(v), len(p.seg.buf))
		panic("out of bounds")
		return errOutOfBounds
	}

	copy(p.seg.buf[off:], v)
	return nil
}

func (p *bufferPointer) deref() error {
	if p.typ.Type() == FarPointer && p.typ.SegmentId() == p.seg.Id {
		p2, err := p.seg.readptr(p.typ.SegmentOffset())
		if err != nil {
			return err
		}
		*p = p2.(bufferPointer)
	}
	return nil
}

func (b *Buffer) readptr(off int) (Pointer, error) {
	if off+8 > len(b.buf) {
		panic("out of bounds")
		return nil, errOutOfBounds
	}

	typ := PointerType{Value: little64(b.buf[off:])}
	ptr := bufferPointer{seg: b, off: off + 8, typ: typ}

	switch typ.Type() {
	case FarPointer:
		ptr.off = 0

	case CompositeList:
		ptr.off += 8 * typ.Offset()
		if ptr.off+8 > len(b.buf) {
			panic("out of bounds")
			return nil, errOutOfBounds
		}

		ptr.typ.Composite = little64(b.buf[ptr.off:])
		ptr.off += 8

		if ptr.off+typ.Elements()*(typ.DataSize()+typ.PointerNum()*8) > len(b.buf) {
			panic("out of bounds")
			return nil, errOutOfBounds
		}

	default:
		println("readptr", typ.Offset(), typ.DataSize(), typ.PointerNum())
		ptr.off += 8 * typ.Offset()
		if ptr.off+typ.DataSize()+typ.PointerNum()*8 > len(b.buf) {
			panic("out of bounds")
			return nil, errOutOfBounds
		}
	}

	return ptr, nil
}

func (p bufferPointer) ReadPtrs(off int, v []Pointer) error {
	if err := p.deref(); err != nil {
		return err
	}

	if p.typ.Type() == CompositeList {
		if off+len(v) > p.typ.Elements() {
			panic("out of bounds")
			return errOutOfBounds
		}

		for i := range v {
			v[i] = bufferPointer{
				seg: p.seg,
				off: p.off + off*(p.typ.DataSize()+8*p.typ.PointerNum()),
				typ: p.typ.CompositeType(),
			}
		}

		return nil
	}

	for i := range v {
		var err error
		v[i], err = p.seg.readptr(p.off + 8*(off+i))
		if err != nil {
			return err
		}
	}

	return nil
}

func (p bufferPointer) WritePtrs(off int, v []Pointer) error {
	if err := p.deref(); err != nil {
		return err
	}

	switch p.typ.Type() {
	case CompositeList:
		for _, src := range v {
			tgt := [1]Pointer{}

			if err := p.ReadPtrs(off, tgt[:]); err != nil {
				return err
			}

			if err := Copy(tgt[0], src); err != nil {
				return err
			}
		}

		return nil
	}

	if off+len(v) > p.typ.PointerNum() {
		fmt.Println(p, off, len(v), p.typ.PointerNum())
		panic("out of bounds")
		return errOutOfBounds
	}

	off = p.off + p.typ.DataSize() + off*8

	for _, src := range v {
		tgt := 0
		typ := PointerType{}

		if src != nil {
			typ = src.Type()
		}

		if src == nil {
			// nothing to do

		} else if bsrc, ok := src.(bufferPointer); ok && bsrc.seg == p.seg {
			tgt = bsrc.off

		} else if copied, ok := p.seg.ptrs[src]; ok {
			tgt = copied

		} else {
			to, err := p.seg.New(typ)
			if err != nil {
				return err
			}

			tgt = to.(bufferPointer).off
			p.seg.ptrs[src] = tgt

			if err := Copy(to, src); err != nil {
				return err
			}
		}

		if src != nil && typ.Type() != FarPointer {
			// Offsets are relative to the end of the pointer
			typ.SetOffset((tgt-off)/8 - 1)
		}

		putLittle64(p.seg.buf[off:], typ.Value)
		off += 8
	}

	return nil
}

func (p bufferPointer) Type() PointerType {
	return p.typ
}
