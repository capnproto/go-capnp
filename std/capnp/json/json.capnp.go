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

const schema_8ef99297a43a5e34 = "x\xda\x8c\x94_h\x1cU\x14\xc6\xcf\x99{'\xdb\x98" +
	"\xd8\xdd\xe9$\x94B\xc2\xd6\x98\xd8\xbaI\xb3\x7f,\xa8" +
	"\x0b%M\xb5\"\xc1?\x9dT}\xc8C\xeddv\\" +
	"6\xce\xce,\xbb\xb3m\xb6Rb\x8bJ-T\xadE" +
	"\x84\x0a\x8a\xa0H\x0a\x82\x0f\xad4\xd8\xd8*V\x8b\x88" +
	"`)\xa2E!>\x08\x0a\x8aV|0Us=w" +
	"'\xbb;\xc4\x08>,\xdc\xfb\x9b;\xe7;\xf7\xfb\xce" +
	"lj\x0e\xb7+i5\xce\x00\x8c~\xb5M\x94\xaf\x0e" +
	"=}J\x0c\x1c\x01\xa3\x03\x15\xb1uO\xf6\x8d\x97_" +
	"\\|\x0evb$\x02\xa0W\xf1\x94~\x107\x01\xdc" +
	"v\x1c\x9fG\xc0\xd6\x01TA\xd1\xefd\x13\xfa6\xb6" +
	"^\xbf\x9f=\x08\x8a\xd8\xfb\xde\x95m\xddO\x9d\x7f\x05" +
	"\xb4nl\xd5U\x15Y\xe8uvY\x7f\x9b\xc9\xd5," +
	"\xdbOu\x86\x1e}\xf2\xfc\xab\x8b\x1f\xcf\x81\x11SQ" +
	"\x1cz u\xf6\x89\x83\x7f\xcc\x03\xa0\xbe\x81_\xd6\x07" +
	"8\x9d\xdc\xbd\x913)\xd9|\x887\x90\xa4\xc6\x0f\xeb" +
	"\xdd|\x93\x9e\xe6Rra\xb0vslf\xf6}0" +
	"\xda\xa9\xcc\x11\xe7\xdb%\xa3'\xf1\x85,s\x8c\x1f\xd5" +
	"_\xaa\x97y!(\xf3\xd9\xed\xb3\xbb\xbfzx\xfa\xc3" +
	"\xd5\xba;\xc4\xaf\xd2\x0br\xf5,\x97\xdd5\x0b!#" +
	"\xc9\xef\xf9\x98\xfe#\xdf\xaf\x8f\xaaq\x92<s\xf7\xda" +
	"[\xf0\xdd\xd4w\xff\x96\x1cU\x0f\xeb;U)\xb9]" +
	"\xadKNU<w\xd82K\xe8\x96\xb2c\xb4~$" +
	"j:U\xdbX\x83\xe1n\xda3!\xe3\xd4D\xfc\x9e" +
	"\x82\xed\xe4\xa2w\x99\x8ec\xf40\xde)\x04G\x00\xed" +
	"L\x82\xf2z\x87\xa1qN\xc1^\\\x12\xb1.\x94x" +
	"n\x07\xe1\xd3\x84/\x10V\xfe\x16\xd8\x85\x0a\xe1\xf9," +
	"\xe1\xb3\x84?R\xf0F\xf6\x97\xe8\xa2{\x80\xf6\x81\xa4" +
	"\xe7\x88^\"\xca\xff$\xca\x89^\xcc\x10\xbd@\xf4\x1b" +
	"\xa2\xeau\xa2*\xd1\xaf\xe5\xd9+D\x7f%\xda\xb6H" +
	"\xb4\x8d\xe8\xcf\xb2\x8b\x1f\x88\xfe\xae`\xd4\xad:\x0e\xb4" +
	"\xcdLz\x9ec\x9b.\xb5\xa3\xd0\x0fG\xdcjq\xd2" +
	".c\x07m;h[\xf1\xcb\x057\x8f\x9d\xb4\xed\x04" +
	"\x8c\x9b\xe5\xb2Y\xc3\xb5\x80\xbb\x18b\xac\x95\x03\xa0\x84" +
	"#\xde\xe4\x94m\xf9\xad\xe7M\x9b\x82\xe7Q\x8bl!" +
	"\xdc4\x8cp,d40\xb7dp\x0c\xe5\xab\xe1\xb8" +
	"\xa8[O\xce\x03\xdaF'\x0dyk\xf6z\x14\x1c\xdd" +
	"\x8c\xa1)\xea'\x90\xc2P\xc6C\x04\xeeCj\x9d*" +
	"e\xe2$\xe1\x96H\xb4\xbeMD\xa5l\xb0K\xe3\x04" +
	"\x8a\x03^q\xb2`\x1f\xb09\xf5\xe2\x15\x93y/Y" +
	"?_\xf6|/\x93\xac\xf8\xb9`\x9b\x9c\xaa0\xf9Z" +
	"\xb3i\xa51\x1dr8\x86)xtv!\x1ak\x18" +
	"\xa5S\xcf\xfe\xd61r}3\xb9\xbeUA\x0d1H" +
	">-\x03\x1a\"x\xaf\x82\xe2\xb1\xaak\xf9\x05Y\x15" +
	"\x1aN\x8f\x94\xcc\xb2Y\xac\xfc\xa7\xd5\"h\xc6\x1aT" +
	"\x06\x87\xeb\xcb\xack\x16\xedJ\xc9\xb4\xd0&yY\xa6" +
	"y\x04\x97\x8f@}t\x9b\xf6i\xed\xe3\xe2\xf3\xe3\xd7" +
	"\x96jo\xe5~\xa3MB4*\x90\xd1Q\xb9^\xc5" +
	"\xedT\xd8mi\xee\x1da\xb7\xe9\x82\xa3\x0f-\xbb}" +
	"2\x90\xcffMt]\xcf7\xe9~\xcc\xad4\xcc\xef" +
	"\x8bX\xd3\xd3\x0d\xef\xc7\xff\xb7\xf7\xd64\xd2[\"\xef" +
	"\x05\x17\xc2,u\xfb\xb8\x99\xb7\x01\x1aw^=\x15\xfa" +
	"(\x99\x93[\x11\x8b\xfc\x18\xfa)\x81T(\x96-\x99" +
	"VVu\x0f\x9a\x93\xbfO\x16Z\x91\x83\x9c\xddF/" +
	"\x94\x9d\xd1\x83!{\xb4\xf4\x8e\x965\xda\x96\xac\xd8s" +
	"\xe25c\xfe\xcb\xa3\x17I\xbaO|\xf2\xcb\xa7\xfd\x1b" +
	"N\xfbo\x826\xd0'\xd6-\x8c\xffT{f\xdf%" +
	"\xd0n\xca\x88\x13\x1b\x93\x0b'\xed\xd8u\xd0z'\xc4" +
	"\xb5c\xc9\xf5\xeb\xf6\xce\xd1\xffLobf\xf9\xb2#" +
	"\x85b\xc9+\xfb\x91\x9cgE|3\x1f\x97\xee\xe6\x85" +
	"U\xad\xf8^\xd1\xaf\x01+-\xa7\xc71\xfc\x07\xcb)" +
	"\x9b\xd8\xca/!\xe4dP\xb4\xe1\xe3?\x01\x00\x00\xff" +
	"\xff\xbcH\xda\xe6"
