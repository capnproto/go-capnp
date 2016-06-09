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

const schema_85d3acc39d94e0f8 = "x\xdad\x91Oh\x13]\x14\xc5\xef\x99L\x9a~P" +
	"\xbef:\x85*D\xa2\xb5\x0a\x8d&\xd3dc\xcdF" +
	"\xad\xae\\9u\xe7B\x9cL\x87\x10\x9a\x997\xa4/" +
	"\x96\x14\xa1XPJ\xc0U]\x09\xe2\xca\x85\xd0\xa5\x0a" +
	"\x82\xb1\xa8\x14\xff\xa0\x0b\xb1\xe8J\x88{A\xdc\x89\x0b" +
	"\x9fw\xd24\x0e\xb8\x98\x19\xee\xbdo\xce\xb9\xbf\xf3\xd2" +
	"\x9b\xa7\xb5br\x0bD\xf6xrh\xeb\xda\xce\xce\xe7" +
	"\xf5w\xd7\xed\x11@\xfd\xec\xde\xbe\xfbr\xf3\xe3\x0dJ" +
	"\"Edl\xb7\x8d\xf7\xd1\xf7\xcd2\xc5f\x18\"\xcd" +
	"\x9cF\xc5\xccc\xc2<\x89\x13\xa4\xa9\xee\xb1\xd6\xe1\xf4" +
	"\xea\x83gd\xff\x97\x84Z\xaf\x7f\xf9mgr\x1f\x88" +
	"`\xfah\x9b\xcdH\xecb\x88\x04(6D\x82e\x1c" +
	"\x9c7=,\x9b\x1ddY\xe6\xd1\xb9\xff\x8f\xe2\xf1\xcc" +
	"\xd7\x7fe:X3\x9f\xf7d\x9e\xee\xcaT\x84X\\" +
	"*\xb8\x0e\xc2 ,\xcfqAt\x01\xb0\x87\x13:\x91" +
	"\xceh\xc6t\x89\xf9\xa6\x12\xb0g4\x00\xe3\x88z\xf9" +
	"y\xee\x1d\xe7\xde\xac\x86\xac\xac\xc9\xba\x87\x11\xd2\xf8\x81" +
	"\x0a\x9d\xaawV4\x03\x82\x84\xce==\xe6A)6" +
	"\xb1u`/+\x03\xb9\xd1\xc8\xd3\x1eF\x9c=\xad\xe1" +
	"L\x061\x8a}\xdc\x98\x05\xeb\x93\x81R\xb6'\xc70" +
	"\\\x16\xb1\x06\xb5\"\xfcJ\xcd[\xf1\x92A\xc1\x15\xbe" +
	"U\x15\x96\xeb\xb0QCHQ\xb2j\x81\xf4\x1a\x81S" +
	"\xb7\x16<_X\xfd_IUE\xa1w\x08\xe5\xd0q" +
	"\x17y\xe7\x1evD1\x18\xf11;\xc3w9\xd8\xcb" +
	"(\xce\xfd\xdd\xc9\xc8\x97\xd5\xe5\x8d{v\xe7S{\x9b" +
	"C\x9aT\xaf\xbe\xbf\x9d\xda\xffP\xde'\xe3\xc8\xa4\x1a" +
	"\xeb\xce\x7fk\xdd\xbc\xfa\x9a\x8cC%\xb5q\xd0\xea\xde" +
	"\xf1\xd2\xbf\xc88pI\xfd\xb8eM\x8c]y\xf2\x82" +
	"\x8b\xdcj\xdf\xfbT\xcd\x0fEC\xa6\x16\x84\x9b\x92N" +
	"5\x1b\x08~+\xb7\xb9$\x85/[\x94\x08\xbd\xd1\xc0" +
	"\xf1=N.\x9e\x92\xce\xa1\xa4\x07\xa1\xf4W\x8e\x81\xed" +
	"\x8a\xeea\xfd\x09\x00\x00\xff\xffj\x83\xcf\xa6"
