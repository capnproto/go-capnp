package json

// AUTO GENERATED - DO NOT EDIT

import (
	math "math"
	strconv "strconv"
	capnp "zombiezen.com/go/capnproto2"
	schemas "zombiezen.com/go/capnproto2/schemas"
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
	if err != nil {
		return JsonValue{}, err
	}
	return JsonValue{st}, nil
}

func NewRootJsonValue(s *capnp.Segment) (JsonValue, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 1})
	if err != nil {
		return JsonValue{}, err
	}
	return JsonValue{st}, nil
}

func ReadRootJsonValue(msg *capnp.Message) (JsonValue, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return JsonValue{}, err
	}
	return JsonValue{root.Struct()}, nil
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
	if err != nil {
		return "", err
	}
	return p.Text(), nil
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
	if err != nil {
		return JsonValue_List{}, err
	}
	return JsonValue_List{List: p.List()}, nil
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
	if err != nil {
		return JsonValue_Field_List{}, err
	}
	return JsonValue_Field_List{List: p.List()}, nil
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
	if err != nil {
		return JsonValue_Call{}, err
	}
	return JsonValue_Call{Struct: p.Struct()}, nil
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
	if err != nil {
		return JsonValue_List{}, err
	}
	return JsonValue_List{l}, nil
}

func (s JsonValue_List) At(i int) JsonValue           { return JsonValue{s.List.Struct(i)} }
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
	if err != nil {
		return JsonValue_Field{}, err
	}
	return JsonValue_Field{st}, nil
}

func NewRootJsonValue_Field(s *capnp.Segment) (JsonValue_Field, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	if err != nil {
		return JsonValue_Field{}, err
	}
	return JsonValue_Field{st}, nil
}

func ReadRootJsonValue_Field(msg *capnp.Message) (JsonValue_Field, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return JsonValue_Field{}, err
	}
	return JsonValue_Field{root.Struct()}, nil
}
func (s JsonValue_Field) Name() (string, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return "", err
	}
	return p.Text(), nil
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
	if err != nil {
		return JsonValue{}, err
	}
	return JsonValue{Struct: p.Struct()}, nil
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
	if err != nil {
		return JsonValue_Field_List{}, err
	}
	return JsonValue_Field_List{l}, nil
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
	if err != nil {
		return JsonValue_Call{}, err
	}
	return JsonValue_Call{st}, nil
}

func NewRootJsonValue_Call(s *capnp.Segment) (JsonValue_Call, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	if err != nil {
		return JsonValue_Call{}, err
	}
	return JsonValue_Call{st}, nil
}

func ReadRootJsonValue_Call(msg *capnp.Message) (JsonValue_Call, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return JsonValue_Call{}, err
	}
	return JsonValue_Call{root.Struct()}, nil
}
func (s JsonValue_Call) Function() (string, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return "", err
	}
	return p.Text(), nil
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
	if err != nil {
		return JsonValue_List{}, err
	}
	return JsonValue_List{List: p.List()}, nil
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
	if err != nil {
		return JsonValue_Call_List{}, err
	}
	return JsonValue_Call_List{l}, nil
}

func (s JsonValue_Call_List) At(i int) JsonValue_Call           { return JsonValue_Call{s.List.Struct(i)} }
func (s JsonValue_Call_List) Set(i int, v JsonValue_Call) error { return s.List.SetStruct(i, v.Struct) }

// JsonValue_Call_Promise is a wrapper for a JsonValue_Call promised by a client call.
type JsonValue_Call_Promise struct{ *capnp.Pipeline }

func (p JsonValue_Call_Promise) Struct() (JsonValue_Call, error) {
	s, err := p.Pipeline.Struct()
	return JsonValue_Call{s}, err
}

const schema_8ef99297a43a5e34 = "x\xdat\x91\xcdk\x13A\x18\xc6\xe7\x99\xc9\xa6\x95D" +
	"\x93\x90\xf4\xa6\xe4RQK-i,\x88\x01\x89Z*" +
	"\xd2S\x17\xd1\xa3\xb8IWI\xd9\xcc\x86M\xd6\x8f\x83" +
	"xQ\xf0\xa4\xe8\xc5\x83W/\xf6\xa4\xa0`\xb1\xc5*" +
	"\x0a\x1e<I\x0fZ\xfc\x03\xbc\x08\x1e\xbc\xd4\xaf\xf1\x99" +
	"\xd6&\xa1\xd8\xc3\xc2\xcco\xde\x9d\xf7}~S*\xe1" +
	"\x98\x1cw\x8aJ\x08w\xd8I\x9ahu\xf4\xe6\xbc\xd9" +
	"{K\xb8)H3q\xae\xf2\xf0\xfe\xbd\xb5\xdbb\x0a" +
	"\x03\x03B\xe4c\xcc\xe7\xafa\x9f\x10\x87\xee\xe2\x0e\x04" +
	"\xcc\xf9\x17+G\x87n\xbc| rC\xe8\xfd\xebH" +
	"[|D}\xc8O)\xbb:\xae.\xb3\xf6\xfd\xe1G" +
	"\xa7?\x9e\xb9\xf2\xfa\x7f\xb5O\xd4j~i\xbdv\x81" +
	"\xb5\x93f\xae\x1d\xea\xb1\xba\xd7\x82nU\xa6\xb9>\x9b" +
	"\xf1\x82\xd8w\x07\xd1\x7f\xcd\x8er_\x7fg\xa4x\xb2" +
	"\xe1\x07\xb3\x99I/\x08\xdc\xdd*\x916&\x01!r" +
	"\xcfF\x18\xed\xb1\x82\xbb(\xb1\x07\x7fL\xb6\x00\x8b\x17" +
	"N\x10?%^&\x96\xbf\x0d\x0a\x90\xc4K\x15\xe2\xe7" +
	"\xc4o$v\xaa_\xa6\x00\x9a\xc9\xbd\xb2t\x91\xf4\x1d" +
	"i\xe2'i\x82\xf4m\x99t\x99\xf43\xa9\xf3\x83\xd4" +
	"!\xfddkWH\xbf\x91&\xd7H\x93\xa4_\xed\x14" +
	"_H\xbfKdt\x1c\x04\"y\xbd\x16\x86\x81\xefi" +
	"\x8e#\xf9\xa1\xaa\xe3f\xcd\x8f\x90\xe26\xc5m\xbb\x13" +
	"5\xf4E\xa4\xb9M\x0b\x14\xbd(\xf2\xaeb\x97\xc0\x8c" +
	"\x02\xb2=\x81\x02\x16V\xc3\xda\x9c_\xef\xf4\xce\xbb\x9a" +
	"6\xce3uj!\xee\x0a#\xce\xf2Q6E\xcbM" +
	"\xd1\xd6\xf3\x18\x1d\"\x98\x01\xdcA\xc5\xa0\xeb\x1a\x0fL" +
	"3\xc0~\x06\x98\x90\xc8\x01\x1b\x12\xc7m\xd6Q\xc2S" +
	"\x12\xe6B\xac\xeb\x9dF\xa8\x05\xef\xfe7t\xb5\xe5E" +
	"^\xb3\xbd\xed\xd4\xdb\xb4\xe7C\xaa`vK\x7f+p" +
	"\x98\xadJ}\xfd\x0f\x96{Ce\xb4\xd7\xf4\xbb\xb6." +
	"\xd9\x8b\xb64d\xde\xbf\x01\x00\x00\xff\xff\xaa\x8f\xc0\x9e"

func init() {
	schemas.Register(schema_8ef99297a43a5e34,
		0x8825ffaa852cda72,
		0x9bbf84153dd4bb60,
		0xc27855d853a937cc)
}
