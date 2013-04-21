package capnproto

/*
struct Message {
	Cookie @1 :Pointer;
	Method @2 :UInt16;
	ReturnCookie @3 :Pointer;
	ReturnMethod @4 :UInt16;
	Arguments @5 :Pointer
}
*/

type Message struct {
	Ptr Pointer
}

func (m *Message) Cookie() Pointer {
	return m.Ptr.ReadPtr(0)
}

func (m *Message) SetCookie(v Pointer) error {
	return m.Ptr.WritePtr(0, v)
}

func (m *Message) Method() uint16 {
	return m.Ptr.ReadUInt16(0)
}

func (m *Message) SetMethod(v uint16) error {
	return m.Ptr.WriteUint16(0, v)
}

func (m *Message) ReturnCookie() Pointer {
	return m.Ptr.ReadPtr(1)
}

func (m *Message) SetReturnCookie(v Pointer) error {
	return m.Ptr.WritePtr(1, v)
}

func (m *Message) ReturnMethod() uint16 {
	return m.Ptr.ReadUInt16(2)
}

func (m *Message) SetReturnMethod(v uint16) error {
	return m.Ptr.WriteUint16(2, v)
}

type Connection interface {
	Call(iface Pointer, method int, args Pointer) (Pointer, error)
}

type MessageReader struct {
	r *StreamDecompressor
	connection Connection
}

func NewMessageReader(r io.Reader, c Connection) {
	return &MessageReader{
		r: NewStreamDecompressor(r),
		connection: c,
	}
}

func (s *MessageReader) ReadMessage() (Message, error) {
	buf := make([]byte, 8)
	if _, err := io.ReadFull(s.r, buf); err != nil {
		return Message{}, err
	}

	sz := little64(hdr[:])
	if sz > 1<<30 {
		return Message{}, errOverlargeSegment
	}

	typ := PointerType{Value: little64(hdr[:])}
	method := typ.Offset()

	sz := typ.DataSize()
	if typ.Type() == CompositeList {
		sz = typ.CompositeDataSize()
	}

	// Put the pointer back into the buffer after resetting the
	// offset to 0 as this contains the method index
	typ.SetOffset(0)
	putLittle64(hdr[:], typ.Value)

	buf = append(buf, hdr[:])
	buf = append(buf, make([]byte, sz)...)
	if _, err := io.ReadFull(s.r, buf[8:]); err != nil {
		return Message{}, err
	}

	seg := NewBuffer(buf)
	ptr, err := seg.ReadRoot(0)
	if err != nil {
		return Message{}, err
	}

	return Message{Ptr: ptr}, nil
}

type MessageWriter struct {
	w *StreamCompressor
}

func NewMessageWriter(w io.Writer) {
	return &MessageWriter{
		w: NewStreamCompressor(w),
	}
}

func (s *MessageWriter) WriteMessage(m Message) error {
	buf := NewBuffer(make([]byte, 8))
	ptr, _ := buf.NewRoot(m.Ptr.Type())
	if err := Copy(ptr, m.Ptr); err != nil {
		return err
	}

	data := buf.Bytes()
	putLittle64(data, len(data)-8)

	if _, err := s.w.Write(data); err != nil {
		return err
	}

	return nil
}
