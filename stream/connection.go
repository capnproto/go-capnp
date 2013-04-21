package capnproto

import (
	"encoding/binary"
	"errors"
	"github.com/jmckaskill/go-capnproto"
	"github.com/jmckaskill/go-capnproto/msgs"
	"io"
)

var (
	little64    = binary.LittleEndian.Uint64
	putLittle64 = binary.LittleEndian.PutUint64
)

type MessageReader struct {
	r *Decompressor
}

func NewMessageReader(r io.Reader) *MessageReader {
	return &MessageReader{
		r: NewDecompressor(r),
	}
}

func (s *MessageReader) ReadMessage() (ret msgs.Message, err error) {
	hdr := [8]byte{}
	if _, err = io.ReadFull(s.r, hdr[:]); err != nil {
		return
	}

	sz := little64(hdr[:])
	if sz > 1<<27 {
		err = errors.New("capnproto: overlarge message")
		return
	}

	buf := capnproto.NewBuffer(make([]byte, int(sz)))
	if _, err = io.ReadFull(s.r, buf.Bytes()); err != nil {
		return
	}

	ret.Ptr, err = buf.ReadRoot(0)
	return
}

type MessageWriter struct {
	w *Compressor
}

func NewMessageWriter(w io.Writer) *MessageWriter {
	return &MessageWriter{
		w: NewCompressor(w),
	}
}

func (s *MessageWriter) WriteMessage(m msgs.Message) error {
	buf := capnproto.NewBuffer(make([]byte, 8))
	buf.NewRoot(m.Ptr.Type())
	ptr, err := buf.New(m.Ptr.Type())
	if err != nil {
		return err
	}
	if err := capnproto.Copy(ptr, m.Ptr); err != nil {
		return err
	}

	data := buf.Bytes()
	putLittle64(data, uint64(len(data)-8))

	if _, err := s.w.Write(data); err != nil {
		return err
	}

	return nil
}
