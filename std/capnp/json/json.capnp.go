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

var schema_8ef99297a43a5e34 = []byte{
	120, 218, 140, 148, 91, 104, 28, 85,
	28, 198, 207, 127, 46, 217, 198, 164,
	221, 157, 78, 74, 17, 18, 86, 99,
	106, 235, 38, 205, 110, 214, 130, 178,
	32, 209, 122, 65, 130, 104, 39, 85,
	31, 242, 80, 59, 153, 29, 151, 141,
	179, 51, 203, 238, 108, 155, 173, 72,
	172, 24, 41, 5, 21, 237, 131, 88,
	81, 95, 20, 105, 65, 240, 161, 149,
	6, 27, 91, 197, 106, 17, 17, 148,
	86, 108, 49, 144, 130, 15, 22, 20,
	141, 248, 96, 170, 230, 248, 157, 153,
	189, 12, 177, 1, 31, 6, 102, 126,
	231, 204, 255, 246, 125, 231, 100, 230,
	232, 110, 105, 68, 77, 202, 140, 25,
	3, 106, 7, 63, 121, 223, 134, 91,
	233, 195, 204, 21, 102, 116, 170, 196,
	15, 57, 11, 43, 70, 111, 234, 27,
	198, 72, 175, 209, 115, 122, 157, 98,
	140, 237, 246, 73, 38, 70, 252, 224,
	195, 153, 83, 79, 63, 243, 231, 60,
	163, 27, 152, 164, 219, 88, 46, 210,
	86, 125, 150, 30, 97, 18, 95, 28,
	172, 223, 146, 152, 57, 246, 241, 127,
	195, 92, 161, 195, 250, 213, 32, 204,
	143, 97, 152, 202, 229, 161, 217, 227,
	124, 203, 33, 102, 116, 145, 196, 119,
	236, 201, 189, 243, 218, 171, 203, 47,
	177, 251, 41, 134, 77, 250, 69, 58,
	174, 47, 208, 86, 198, 110, 95, 162,
	151, 197, 246, 214, 6, 82, 145, 245,
	117, 121, 66, 127, 83, 222, 172, 191,
	47, 139, 172, 123, 63, 186, 112, 215,
	166, 231, 207, 188, 193, 180, 77, 145,
	184, 170, 36, 2, 145, 242, 173, 190,
	94, 17, 111, 157, 202, 126, 196, 249,
	234, 142, 99, 187, 191, 127, 108, 250,
	211, 235, 237, 173, 43, 151, 245, 217,
	96, 239, 193, 96, 239, 208, 19, 207,
	158, 121, 107, 249, 243, 57, 102, 36,
	212, 72, 223, 232, 102, 1, 81, 175,
	42, 65, 55, 74, 208, 77, 171, 85,
	146, 81, 222, 69, 101, 76, 191, 164,
	236, 215, 49, 96, 148, 87, 240, 134,
	45, 179, 236, 150, 41, 87, 44, 149,
	189, 138, 207, 118, 17, 81, 55, 22,
	2, 154, 182, 6, 105, 48, 220, 192,
	140, 117, 20, 201, 170, 117, 142, 243,
	175, 95, 89, 90, 169, 191, 151, 255,
	29, 31, 41, 238, 154, 37, 187, 90,
	54, 45, 70, 118, 92, 188, 27, 221,
	152, 92, 187, 200, 94, 137, 238, 201,
	80, 68, 131, 33, 128, 59, 41, 162,
	237, 14, 128, 71, 69, 110, 166, 209,
	209, 48, 125, 46, 103, 146, 235, 122,
	190, 233, 23, 61, 217, 173, 162, 185,
	96, 181, 63, 102, 77, 79, 135, 31,
	35, 52, 78, 252, 128, 87, 154, 44,
	218, 7, 108, 197, 29, 182, 188, 82,
	186, 224, 165, 131, 191, 43, 158, 239,
	101, 211, 85, 63, 159, 110, 244, 50,
	77, 248, 43, 210, 49, 170, 125, 202,
	44, 216, 172, 213, 243, 84, 213, 115,
	197, 34, 33, 245, 24, 222, 31, 143,
	155, 78, 205, 14, 26, 111, 75, 211,
	153, 141, 104, 170, 166, 146, 15, 20,
	109, 39, 31, 191, 215, 116, 28, 163,
	87, 86, 186, 57, 87, 8, 69, 158,
	76, 193, 188, 31, 200, 100, 156, 150,
	168, 143, 86, 120, 162, 135, 4, 158,
	219, 9, 124, 2, 248, 44, 176, 244,
	15, 167, 30, 146, 128, 231, 115, 192,
	167, 128, 63, 147, 104, 189, 252, 55,
	239, 129, 82, 76, 251, 68, 208, 211,
	160, 231, 65, 149, 191, 64, 21, 208,
	115, 89, 208, 179, 160, 63, 128, 170,
	215, 64, 85, 208, 75, 98, 239, 5,
	208, 223, 64, 59, 150, 65, 59, 64,
	127, 17, 85, 252, 4, 250, 135, 68,
	113, 183, 230, 56, 172, 99, 102, 210,
	243, 28, 219, 116, 81, 142, 132, 135,
	70, 221, 90, 105, 210, 174, 80, 23,
	62, 187, 240, 89, 245, 43, 69, 183,
	32, 198, 129, 135, 146, 102, 165, 98,
	214, 105, 3, 163, 93, 50, 81, 162,
	109, 74, 70, 2, 142, 122, 147, 83,
	182, 229, 183, 215, 91, 99, 10, 215,
	227, 22, 198, 2, 220, 26, 24, 112,
	2, 142, 108, 14, 154, 201, 110, 217,
	80, 40, 98, 118, 141, 198, 121, 48,
	122, 76, 30, 86, 186, 142, 139, 182,
	69, 93, 52, 16, 218, 170, 237, 34,
	97, 171, 135, 26, 46, 202, 38, 67,
	231, 54, 108, 147, 138, 139, 180, 77,
	223, 76, 252, 111, 223, 76, 85, 101,
	241, 91, 171, 104, 169, 233, 14, 97,
	142, 97, 8, 79, 14, 12, 100, 172,
	147, 161, 78, 160, 253, 109, 99, 152,
	250, 54, 76, 29, 150, 214, 136, 66,
	229, 71, 132, 64, 67, 128, 15, 74,
	196, 159, 172, 185, 22, 60, 141, 168,
	172, 57, 233, 209, 178, 89, 49, 75,
	213, 53, 71, 189, 70, 122, 184, 79,
	118, 242, 171, 242, 11, 213, 7, 144,
	42, 19, 201, 191, 61, 219, 46, 42,
	56, 160, 45, 137, 247, 137, 64, 171,
	18, 10, 145, 154, 55, 128, 212, 184,
	1, 114, 205, 35, 78, 118, 243, 192,
	52, 207, 18, 250, 48, 122, 41, 34,
	140, 54, 178, 179, 45, 138, 182, 61,
	199, 247, 28, 121, 219, 152, 255, 238,
	240, 57, 84, 215, 207, 191, 248, 245,
	203, 129, 27, 79, 248, 239, 50, 109,
	75, 63, 223, 184, 56, 254, 115, 253,
	133, 125, 231, 153, 118, 115, 150, 31,
	185, 41, 189, 120, 212, 78, 92, 99,
	90, 223, 4, 95, 122, 49, 189, 121,
	227, 222, 57, 156, 185, 190, 212, 76,
	227, 176, 142, 134, 215, 84, 44, 239,
	89, 49, 223, 44, 36, 197, 237, 80,
	224, 86, 173, 234, 123, 37, 191, 206,
	228, 114, 227, 246, 81, 40, 122, 225,
	43, 112, 69, 98, 181, 43, 254, 13,
	0, 0, 255, 255, 95, 36, 223, 195,
}
