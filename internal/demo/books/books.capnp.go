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

var schema_85d3acc39d94e0f8 = []byte{
	120, 218, 100, 146, 63, 104, 19, 97,
	24, 198, 223, 231, 114, 105, 42, 20,
	155, 175, 87, 168, 133, 148, 104, 173,
	98, 163, 233, 53, 89, 44, 89, 212,
	234, 228, 228, 213, 205, 65, 188, 164,
	71, 8, 205, 221, 119, 164, 95, 172,
	41, 74, 177, 160, 148, 128, 83, 157,
	4, 113, 114, 16, 58, 170, 224, 160,
	5, 165, 248, 7, 29, 196, 162, 147,
	16, 119, 65, 220, 196, 193, 207, 247,
	210, 52, 70, 58, 220, 29, 223, 243,
	242, 254, 222, 231, 125, 238, 75, 110,
	156, 54, 114, 241, 77, 16, 57, 195,
	241, 190, 205, 235, 219, 219, 95, 214,
	222, 223, 116, 6, 0, 253, 171, 117,
	247, 254, 171, 141, 79, 183, 40, 142,
	4, 145, 216, 106, 138, 15, 209, 247,
	237, 18, 65, 63, 57, 183, 255, 40,
	158, 78, 127, 35, 103, 95, 28, 122,
	173, 250, 245, 143, 147, 202, 124, 36,
	130, 53, 137, 85, 43, 27, 53, 92,
	60, 134, 24, 168, 167, 136, 24, 25,
	214, 40, 206, 91, 99, 88, 178, 234,
	72, 147, 161, 91, 199, 27, 135, 147,
	43, 143, 94, 236, 197, 212, 209, 180,
	110, 180, 49, 215, 118, 48, 93, 55,
	232, 99, 76, 5, 69, 203, 199, 136,
	213, 192, 73, 198, 20, 165, 92, 88,
	156, 42, 185, 8, 131, 176, 48, 203,
	7, 162, 11, 128, 211, 31, 51, 137,
	76, 94, 77, 76, 230, 121, 191, 137,
	24, 156, 105, 3, 192, 48, 34, 45,
	59, 199, 218, 9, 214, 102, 12, 164,
	85, 69, 85, 61, 12, 144, 193, 15,
	116, 232, 150, 189, 179, 178, 30, 16,
	20, 76, 214, 76, 214, 202, 146, 7,
	48, 31, 133, 138, 31, 202, 154, 138,
	38, 68, 13, 221, 2, 49, 46, 197,
	177, 117, 119, 18, 185, 217, 127, 57,
	137, 108, 65, 95, 94, 127, 224, 60,
	255, 220, 220, 98, 63, 227, 250, 245,
	143, 119, 19, 163, 143, 213, 67, 18,
	71, 198, 245, 80, 107, 238, 123, 227,
	246, 213, 55, 36, 14, 229, 245, 250,
	65, 187, 117, 207, 75, 254, 38, 49,
	118, 73, 255, 188, 99, 143, 12, 93,
	121, 246, 146, 15, 153, 149, 208, 45,
	45, 176, 181, 83, 59, 14, 18, 243,
	178, 148, 80, 110, 57, 29, 72, 126,
	235, 82, 125, 81, 73, 95, 53, 40,
	22, 122, 131, 129, 235, 123, 142, 137,
	222, 132, 77, 3, 103, 146, 145, 101,
	18, 200, 167, 59, 150, 123, 214, 234,
	192, 169, 187, 216, 110, 170, 148, 224,
	58, 195, 176, 123, 59, 4, 50, 131,
	81, 202, 78, 255, 127, 3, 146, 60,
	32, 213, 123, 53, 14, 176, 48, 211,
	157, 216, 198, 241, 175, 229, 99, 14,
	171, 208, 203, 210, 47, 86, 188, 101,
	47, 30, 76, 149, 164, 111, 151, 165,
	221, 54, 82, 147, 74, 230, 237, 74,
	160, 188, 90, 224, 86, 237, 121, 207,
	151, 118, 167, 149, 254, 6, 0, 0,
	255, 255, 147, 176, 205, 234,
}
