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

	return p.Data(), nil

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
