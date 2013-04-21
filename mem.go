package capnproto

var errOutOfBounds = errors.New("capnproto: supplied offset is out of bounds")

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
			buf := make([]byte, min(ft.DataSize(), tt.DataSize()))

			for i := 0; i < min(ft.Elements(), tt.Elements()); i++ {
				fp := [1]Pointer{}
				tp := [1]Pointer{}

				if err := ft.ReadPtrs(i, fp[:]); err != nil {
					return err
				}

				if err := tt.ReadPtrs(i, tp[:]); err != nil {
					return err
				}

				if err := Copy(tp[:], fp[:]); err != nil {
					return err
				}
			}
			return nil
		}

	switch ft.Type() {
	case Struct:
		if tt.Type() == Struct {
			if err := copyData(to, from, min(ft.DataSize(), tt.DataSize())); err != nil {
				return err
			}
			return copyPointers(to, from, min(ft.PointerSize(), tt.PointerSize()))
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

	return errors.New("capnproto: incompatible copy src/target")
}

func NewMemory(p PointerType) (Pointer, error) {
	return (*memory)(nil).New(p), nil
}

type memory struct {
	data []uint8
	ptrs []Pointer
	typ  PointerType
}

func (*memory) New(p PointerType) (Pointer, error) {
	switch p.Type() {
	case CompositeList:
		ptrs := make([]Pointer, p.Elements())

		for i := range m.ptrs {
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
	return nil, errors.New("capnproto: memory pointers do not support interfaces")
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

func (m *memory) WritePtrs(off int, ptr Pointer) error {
	copy(m.ptrs[off:], v)
	return nil
}

func (m *memory) Free() {
	// nothing todo, the GC will clean up everything for us
}

type Buffer struct {
	buf  []byte
	ptrs map[Pointer]int
	caller Caller
}

// NewBuffer creates a new segment buffer. Data is read/written in wire
// format.
func NewBuffer(buf []byte) *Buffer {
	return &buffer{buf, make(map[Pointer]int)}
}

func (b *Buffer) SetCaller(c Caller) {
	b.caller = c
}

func (b *Buffer) Bytes() []byte {
	return b.buf
}

func align8(x int) int {
	return (x + 7) & ~7
}

// New allocates room for a pointer of type p at the end of the buffer and
// returns a Pointer to manipulate the new data.
func (b *Buffer) New(p PointerType) Pointer {
	sz := align8(p.DataSize() + p.PointerNum()*8)
	off := len(b.buf)

	if p.Type() == CompositeList {
		b.buf = append(b.buf, [8]byte{}[:]...)
		putLittle64(b.buf[off:], p.Composite)

		off += 8
		sz *= p.Elements()
	}

	b.buf = append(b.buf, make([]byte, sz)...)
	return bufferPointer{b, off, p}, nil
}

// NewRoot creates a new root value. This writes a pointer tag followed by
// room for the new data, so that the tag can be pointed to via a far pointer
// from another segment. This returns a pointer to the new data as well as the
// offset into the segment for use in creating a far pointer. The offset is
// relative to the beginning of the buffer provided in NewBuffer.
func (b *Buffer) NewRoot(v PointerType) (Pointer, int) {
	off := len(b.buf)
	v.SetOffset(0)
	b.buf = append(b.buf, [8]byte{}[:]...)
	putLittle64(b.buf[off:], v.Value)
	return b.New(v), off
}

func (b *Buffer) ReadRoot(off int) (Pointer, error) {
	p := bufferPointer{
		typ: MakeList(PointerList, len(b.buf)/8),
		off: 0,
		seg: b,
	}

	return p.ReadPtrs(off)
}

type bufferPointer struct {
	seg *stream
	off int
	typ PointerType
}

func (p bufferPointer) New(typ PointerType) (Pointer, error) {
	return p.seg.New(typ), nil
}

func (p bufferPointer) Call(method int, args Pointer) (Pointer, error) {
	if p.seg.caller == nil {
		return nil, errors.New("capnproto: no bound caller")
	}
	return p.seg.caller.Call(p, method, args)
}

func (p bufferPointer) Read(off int, v []uint8) (int, error) {
	off = p.off + off*8

	if off+len(v) > len(p.seg.buf) {
		return 0, errOutOfBounds
	}

	return copy(v, p.seg.buf[off:], nil)
}

func (p bufferPointer) Write(off int, v []uint8) (int, error) {
	off = p.off + off*8

	if off+len(v) > len(p.seg.buf) {
		return 0, errOutOfBounds
	}

	return copy(p.seg.buf[off:], v), nil
}

func (p bufferPointer) ReadInterface() (interface{}, error) {
	if p.seg.iface == nil {
		return nil, errNoConverter
	}

	return p.seg.iface.Deserialize(p)
}

func (p bufferPointer) WriteInterface(iface interface{}) error {
	if p.seg.iface == nil {
		return errNoConverter
	}

	return p.seg.iface.Serialize(p, off, iface)
}

func (p bufferPointer) ReadPtrs(off int, v []Pointer) error {
	if p.typ.Type() == CompositeList {
		if off+len(v) > p.typ.Elements() {
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

	off = p.off + 8*off
	if off + len(v)*8 > len(p.seg.buf) {
		return errOutOfBounds
	}

	for i := range v {
		typ := PointerType{Value: little64(p.seg.buf[off:])}
		off += 8

		v[i] = bufferPointer{seg: p.seg, off: off, typ: typ}

		switch typ.Type() {
		case FarPointer:
			v[i].off = 0

		case CompositeList:
			v[i].off += 8*typ.Offset()
			if v[i].off+8 > len(p.seg.buf) {
				return errOutOfBounds
			}

			v[i].typ.Composite = little64(p.seg.buf[v[i].off:])
			v[i].off += 8

			if v[i].off+typ.Elements()*(typ.DataSize()+typ.PointerNum()*8) > len(p.seg.buf) {
				return errOutOfBounds
			}

		default:
			v[i].off += 8 * typ.Offset()
			if v[i].off+typ.DataSize()+typ.PointerNum()*8 > len(p.seg.buf) {
				return errOutOfBounds
			}
		}
	}

	return nil
}

func (p bufferPointer) WritePtrs(off int, v []Pointer) error {
	switch p.typ.Type() {
	case CompositeList:
		for _, src := range v {
			tgt := [1]Pointer{}

			if err := p.ReadPtrs(off, tgt[:]); err != nil {
				return err
			}

			if err := Copy(tgt[:], src); err != nil {
				return err
			}
		}

		return nil
	}

	if off + len(v) > p.typ.PointerNum() {
		return errOutOfBounds
	}

	off = p.off + p.typ.DataSize() + off*8

	for _, src := range v {
		tgt := 0
		typ := ptr.Type()

		if bsrc, ok := src.(bufferPointer); ok && bsrc.seg == p.seg {
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

		if typ.DataType() != FarPointer {
			// Offsets are relative to the end of the pointer
			typ.SetOffset((tgt-off)/8 - 1)
		}

		putLittle64(p.seg[off:], typ.Value)
		off += 8
	}

	return nil
}

func (p bufferPointer) Type() PointerType {
	return p.typ
}
