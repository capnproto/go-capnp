package books

// AUTO GENERATED - DO NOT EDIT

import (
	capnp "zombiezen.com/go/capnproto2"
)

type Book struct{ capnp.Struct }

func NewBook(s *capnp.Segment) (Book, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		return Book{}, err
	}
	return Book{st}, nil
}

func NewRootBook(s *capnp.Segment) (Book, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		return Book{}, err
	}
	return Book{st}, nil
}

func ReadRootBook(msg *capnp.Message) (Book, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Book{}, err
	}
	return Book{root.Struct()}, nil
}
func (s Book) Title() (string, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return "", err
	}
	return p.Text(), nil
}

func (s Book) HasTitle() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Book) TitleBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	d := p.Data()
	if len(d) == 0 {
		return d, nil
	}
	return d[:len(d)-1], nil
}

func (s Book) SetTitle(v string) error {
	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, t.List.ToPtr())
}

func (s Book) PageCount() int32 {
	return int32(s.Struct.Uint32(0))
}

func (s Book) SetPageCount(v int32) {
	s.Struct.SetUint32(0, uint32(v))
}

// Book_List is a list of Book.
type Book_List struct{ capnp.List }

// NewBook creates a new list of Book.
func NewBook_List(s *capnp.Segment, sz int32) (Book_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1}, sz)
	if err != nil {
		return Book_List{}, err
	}
	return Book_List{l}, nil
}

func (s Book_List) At(i int) Book           { return Book{s.List.Struct(i)} }
func (s Book_List) Set(i int, v Book) error { return s.List.SetStruct(i, v.Struct) }

// Book_Promise is a wrapper for a Book promised by a client call.
type Book_Promise struct{ *capnp.Pipeline }

func (p Book_Promise) Struct() (Book, error) {
	s, err := p.Pipeline.Struct()
	return Book{s}, err
}

const schema_85d3acc39d94e0f8 = "x\xda\x12\x88w`2d\xdd\xcf\xc8\xc0\x10(\xc2\xca" +
	"\xb6\xbf\xe6\xca\x95\xeb\x1dg\x1a\x03y\x18\x19\xff\xffx" +
	"0e\xee\xe15\x97[\x19X\x19\xd9\x19\x18\x04\x8fv" +
	"\x09\x9e\x05\xd1'\xcb\x19t\xff'\xe5\xe7g\x17\xeb%" +
	"'2\x16\xe4\x15X9\x019\x0c\x0c\x01\x8c\x8c\x81\x1c" +
	"\xcc,\x0c\x0c,@\xc3\x045\x8d\x80&\xaa03\x06" +
	"\x1a0122\x8a0\x82\xc4t\x83\x80b:@1" +
	"\x0b&F\xf9\x92\xcc\x92\x9cTF\x1e\x06& f\xfc" +
	"_\x90\x98\x9e\xea\x9c_\x9a\xc7\xc0X\xc2\xc8\x02\x14c" +
	"a`\x04\x04\x00\x00\xff\xffF\xa9$\xae"
