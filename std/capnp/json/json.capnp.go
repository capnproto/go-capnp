package json

// AUTO GENERATED - DO NOT EDIT

import (
	math "math"
	strconv "strconv"
	capnp "zombiezen.com/go/capnproto2"
)

type JsonValue struct{ capnp.Struct }
type JsonValue_Which uint16

const (
	JsonValue_Which_null    JsonValue_Which = 0
	JsonValue_Which_boolean JsonValue_Which = 1
	JsonValue_Which_number  JsonValue_Which = 2
	JsonValue_Which_string  JsonValue_Which = 3
	JsonValue_Which_array   JsonValue_Which = 4
	JsonValue_Which_object  JsonValue_Which = 5
	JsonValue_Which_call    JsonValue_Which = 6
)

func (w JsonValue_Which) String() string {
	const s = "nullbooleannumberstringarrayobjectcall"
	switch w {
	case JsonValue_Which_null:
		return s[0:4]
	case JsonValue_Which_boolean:
		return s[4:11]
	case JsonValue_Which_number:
		return s[11:17]
	case JsonValue_Which_string:
		return s[17:23]
	case JsonValue_Which_array:
		return s[23:28]
	case JsonValue_Which_object:
		return s[28:34]
	case JsonValue_Which_call:
		return s[34:38]

	}
	return "JsonValue_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
}

func NewJsonValue(s *capnp.Segment) (JsonValue, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 1})
	return JsonValue{st}, err
}

func NewRootJsonValue(s *capnp.Segment) (JsonValue, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 1})
	return JsonValue{st}, err
}

func ReadRootJsonValue(msg *capnp.Message) (JsonValue, error) {
	root, err := msg.RootPtr()
	return JsonValue{root.Struct()}, err
}

func (s JsonValue) Which() JsonValue_Which {
	return JsonValue_Which(s.Struct.Uint16(0))
}
func (s JsonValue) SetNull() {
	s.Struct.SetUint16(0, 0)

}

func (s JsonValue) Boolean() bool {
	return s.Struct.Bit(16)
}

func (s JsonValue) SetBoolean(v bool) {
	s.Struct.SetUint16(0, 1)
	s.Struct.SetBit(16, v)
}

func (s JsonValue) Number() float64 {
	return math.Float64frombits(s.Struct.Uint64(8))
}

func (s JsonValue) SetNumber(v float64) {
	s.Struct.SetUint16(0, 2)
	s.Struct.SetUint64(8, math.Float64bits(v))
}

func (s JsonValue) String() (string, error) {
	p, err := s.Struct.Ptr(0)
	return p.Text(), err
}

func (s JsonValue) HasString() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s JsonValue) StringBytes() ([]byte, error) {
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

func (s JsonValue) SetString(v string) error {
	s.Struct.SetUint16(0, 3)
	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, t.List.ToPtr())
}

func (s JsonValue) Array() (JsonValue_List, error) {
	p, err := s.Struct.Ptr(0)
	return JsonValue_List{List: p.List()}, err
}

func (s JsonValue) HasArray() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s JsonValue) SetArray(v JsonValue_List) error {
	s.Struct.SetUint16(0, 4)
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewArray sets the array field to a newly
// allocated JsonValue_List, preferring placement in s's segment.
func (s JsonValue) NewArray(n int32) (JsonValue_List, error) {
	s.Struct.SetUint16(0, 4)
	l, err := NewJsonValue_List(s.Struct.Segment(), n)
	if err != nil {
		return JsonValue_List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

func (s JsonValue) Object() (JsonValue_Field_List, error) {
	p, err := s.Struct.Ptr(0)
	return JsonValue_Field_List{List: p.List()}, err
}

func (s JsonValue) HasObject() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s JsonValue) SetObject(v JsonValue_Field_List) error {
	s.Struct.SetUint16(0, 5)
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewObject sets the object field to a newly
// allocated JsonValue_Field_List, preferring placement in s's segment.
func (s JsonValue) NewObject(n int32) (JsonValue_Field_List, error) {
	s.Struct.SetUint16(0, 5)
	l, err := NewJsonValue_Field_List(s.Struct.Segment(), n)
	if err != nil {
		return JsonValue_Field_List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

func (s JsonValue) Call() (JsonValue_Call, error) {
	p, err := s.Struct.Ptr(0)
	return JsonValue_Call{Struct: p.Struct()}, err
}

func (s JsonValue) HasCall() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s JsonValue) SetCall(v JsonValue_Call) error {
	s.Struct.SetUint16(0, 6)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewCall sets the call field to a newly
// allocated JsonValue_Call struct, preferring placement in s's segment.
func (s JsonValue) NewCall() (JsonValue_Call, error) {
	s.Struct.SetUint16(0, 6)
	ss, err := NewJsonValue_Call(s.Struct.Segment())
	if err != nil {
		return JsonValue_Call{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

// JsonValue_List is a list of JsonValue.
type JsonValue_List struct{ capnp.List }

// NewJsonValue creates a new list of JsonValue.
func NewJsonValue_List(s *capnp.Segment, sz int32) (JsonValue_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 16, PointerCount: 1}, sz)
	return JsonValue_List{l}, err
}

func (s JsonValue_List) At(i int) JsonValue { return JsonValue{s.List.Struct(i)} }

func (s JsonValue_List) Set(i int, v JsonValue) error { return s.List.SetStruct(i, v.Struct) }

// JsonValue_Promise is a wrapper for a JsonValue promised by a client call.
type JsonValue_Promise struct{ *capnp.Pipeline }

func (p JsonValue_Promise) Struct() (JsonValue, error) {
	s, err := p.Pipeline.Struct()
	return JsonValue{s}, err
}

func (p JsonValue_Promise) Call() JsonValue_Call_Promise {
	return JsonValue_Call_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

type JsonValue_Field struct{ capnp.Struct }

func NewJsonValue_Field(s *capnp.Segment) (JsonValue_Field, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	return JsonValue_Field{st}, err
}

func NewRootJsonValue_Field(s *capnp.Segment) (JsonValue_Field, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	return JsonValue_Field{st}, err
}

func ReadRootJsonValue_Field(msg *capnp.Message) (JsonValue_Field, error) {
	root, err := msg.RootPtr()
	return JsonValue_Field{root.Struct()}, err
}

func (s JsonValue_Field) Name() (string, error) {
	p, err := s.Struct.Ptr(0)
	return p.Text(), err
}

func (s JsonValue_Field) HasName() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s JsonValue_Field) NameBytes() ([]byte, error) {
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

func (s JsonValue_Field) SetName(v string) error {
	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, t.List.ToPtr())
}

func (s JsonValue_Field) Value() (JsonValue, error) {
	p, err := s.Struct.Ptr(1)
	return JsonValue{Struct: p.Struct()}, err
}

func (s JsonValue_Field) HasValue() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s JsonValue_Field) SetValue(v JsonValue) error {
	return s.Struct.SetPtr(1, v.Struct.ToPtr())
}

// NewValue sets the value field to a newly
// allocated JsonValue struct, preferring placement in s's segment.
func (s JsonValue_Field) NewValue() (JsonValue, error) {
	ss, err := NewJsonValue(s.Struct.Segment())
	if err != nil {
		return JsonValue{}, err
	}
	err = s.Struct.SetPtr(1, ss.Struct.ToPtr())
	return ss, err
}

// JsonValue_Field_List is a list of JsonValue_Field.
type JsonValue_Field_List struct{ capnp.List }

// NewJsonValue_Field creates a new list of JsonValue_Field.
func NewJsonValue_Field_List(s *capnp.Segment, sz int32) (JsonValue_Field_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2}, sz)
	return JsonValue_Field_List{l}, err
}

func (s JsonValue_Field_List) At(i int) JsonValue_Field { return JsonValue_Field{s.List.Struct(i)} }

func (s JsonValue_Field_List) Set(i int, v JsonValue_Field) error {
	return s.List.SetStruct(i, v.Struct)
}

// JsonValue_Field_Promise is a wrapper for a JsonValue_Field promised by a client call.
type JsonValue_Field_Promise struct{ *capnp.Pipeline }

func (p JsonValue_Field_Promise) Struct() (JsonValue_Field, error) {
	s, err := p.Pipeline.Struct()
	return JsonValue_Field{s}, err
}

func (p JsonValue_Field_Promise) Value() JsonValue_Promise {
	return JsonValue_Promise{Pipeline: p.Pipeline.GetPipeline(1)}
}

type JsonValue_Call struct{ capnp.Struct }

func NewJsonValue_Call(s *capnp.Segment) (JsonValue_Call, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	return JsonValue_Call{st}, err
}

func NewRootJsonValue_Call(s *capnp.Segment) (JsonValue_Call, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	return JsonValue_Call{st}, err
}

func ReadRootJsonValue_Call(msg *capnp.Message) (JsonValue_Call, error) {
	root, err := msg.RootPtr()
	return JsonValue_Call{root.Struct()}, err
}

func (s JsonValue_Call) Function() (string, error) {
	p, err := s.Struct.Ptr(0)
	return p.Text(), err
}

func (s JsonValue_Call) HasFunction() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s JsonValue_Call) FunctionBytes() ([]byte, error) {
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

func (s JsonValue_Call) SetFunction(v string) error {
	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, t.List.ToPtr())
}

func (s JsonValue_Call) Params() (JsonValue_List, error) {
	p, err := s.Struct.Ptr(1)
	return JsonValue_List{List: p.List()}, err
}

func (s JsonValue_Call) HasParams() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s JsonValue_Call) SetParams(v JsonValue_List) error {
	return s.Struct.SetPtr(1, v.List.ToPtr())
}

// NewParams sets the params field to a newly
// allocated JsonValue_List, preferring placement in s's segment.
func (s JsonValue_Call) NewParams(n int32) (JsonValue_List, error) {
	l, err := NewJsonValue_List(s.Struct.Segment(), n)
	if err != nil {
		return JsonValue_List{}, err
	}
	err = s.Struct.SetPtr(1, l.List.ToPtr())
	return l, err
}

// JsonValue_Call_List is a list of JsonValue_Call.
type JsonValue_Call_List struct{ capnp.List }

// NewJsonValue_Call creates a new list of JsonValue_Call.
func NewJsonValue_Call_List(s *capnp.Segment, sz int32) (JsonValue_Call_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2}, sz)
	return JsonValue_Call_List{l}, err
}

func (s JsonValue_Call_List) At(i int) JsonValue_Call { return JsonValue_Call{s.List.Struct(i)} }

func (s JsonValue_Call_List) Set(i int, v JsonValue_Call) error { return s.List.SetStruct(i, v.Struct) }

// JsonValue_Call_Promise is a wrapper for a JsonValue_Call promised by a client call.
type JsonValue_Call_Promise struct{ *capnp.Pipeline }

func (p JsonValue_Call_Promise) Struct() (JsonValue_Call, error) {
	s, err := p.Pipeline.Struct()
	return JsonValue_Call{s}, err
}
