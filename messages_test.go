package capn_test

// AUTO GENERATED - DO NOT EDIT

import (
	"bufio"
	"bytes"
	"encoding/json"
	C "github.com/jmckaskill/go-capnproto"
	"io"
	"unsafe"
)

type Zdate C.Struct

func NewZdate(s *C.Segment) Zdate      { return Zdate(s.NewStruct(8, 0)) }
func NewRootZdate(s *C.Segment) Zdate  { return Zdate(s.NewRootStruct(8, 0)) }
func ReadRootZdate(s *C.Segment) Zdate { return Zdate(s.Root(0).ToStruct()) }
func (s Zdate) Year() int16            { return int16(C.Struct(s).Get16(0)) }
func (s Zdate) SetYear(v int16)        { C.Struct(s).Set16(0, uint16(v)) }
func (s Zdate) Month() uint8           { return C.Struct(s).Get8(2) }
func (s Zdate) SetMonth(v uint8)       { C.Struct(s).Set8(2, v) }
func (s Zdate) Day() uint8             { return C.Struct(s).Get8(3) }
func (s Zdate) SetDay(v uint8)         { C.Struct(s).Set8(3, v) }
func (s Zdate) WriteJSON(w io.Writer) error {
	b := bufio.NewWriter(w)
	var err error
	var buf []byte
	_ = buf
	err = b.WriteByte('{')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"year\":")
	if err != nil {
		return err
	}
	{
		s := s.Year()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(',')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"month\":")
	if err != nil {
		return err
	}
	{
		s := s.Month()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(',')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"day\":")
	if err != nil {
		return err
	}
	{
		s := s.Day()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte('}')
	if err != nil {
		return err
	}
	err = b.Flush()
	return err
}
func (s Zdate) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteJSON(&b)
	return b.Bytes(), err
}

type Zdate_List C.PointerList

func NewZdateList(s *C.Segment, sz int) Zdate_List { return Zdate_List(s.NewUInt32List(sz)) }
func (s Zdate_List) Len() int                      { return C.PointerList(s).Len() }
func (s Zdate_List) At(i int) Zdate                { return Zdate(C.PointerList(s).At(i).ToStruct()) }
func (s Zdate_List) ToArray() []Zdate              { return *(*[]Zdate)(unsafe.Pointer(C.PointerList(s).ToArray())) }

type Zdata C.Struct

func NewZdata(s *C.Segment) Zdata      { return Zdata(s.NewStruct(0, 1)) }
func NewRootZdata(s *C.Segment) Zdata  { return Zdata(s.NewRootStruct(0, 1)) }
func ReadRootZdata(s *C.Segment) Zdata { return Zdata(s.Root(0).ToStruct()) }
func (s Zdata) Data() []byte           { return C.Struct(s).GetObject(0).ToData() }
func (s Zdata) SetData(v []byte)       { C.Struct(s).SetObject(0, s.Segment.NewData(v)) }
func (s Zdata) WriteJSON(w io.Writer) error {
	b := bufio.NewWriter(w)
	var err error
	var buf []byte
	_ = buf
	err = b.WriteByte('{')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"data\":")
	if err != nil {
		return err
	}
	{
		s := s.Data()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte('}')
	if err != nil {
		return err
	}
	err = b.Flush()
	return err
}
func (s Zdata) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteJSON(&b)
	return b.Bytes(), err
}

type Zdata_List C.PointerList

func NewZdataList(s *C.Segment, sz int) Zdata_List { return Zdata_List(s.NewCompositeList(0, 1, sz)) }
func (s Zdata_List) Len() int                      { return C.PointerList(s).Len() }
func (s Zdata_List) At(i int) Zdata                { return Zdata(C.PointerList(s).At(i).ToStruct()) }
func (s Zdata_List) ToArray() []Zdata              { return *(*[]Zdata)(unsafe.Pointer(C.PointerList(s).ToArray())) }
