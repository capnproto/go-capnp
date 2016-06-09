package rpctwoparty

// AUTO GENERATED - DO NOT EDIT

import (
	capnp "zombiezen.com/go/capnproto2"
)

type Side uint16

// Values of Side.
const (
	Side_server Side = 0
	Side_client Side = 1
)

// String returns the enum's constant name.
func (c Side) String() string {
	switch c {
	case Side_server:
		return "server"
	case Side_client:
		return "client"

	default:
		return ""
	}
}

// SideFromString returns the enum value with a name,
// or the zero value if there's no such value.
func SideFromString(c string) Side {
	switch c {
	case "server":
		return Side_server
	case "client":
		return Side_client

	default:
		return 0
	}
}

type Side_List struct{ capnp.List }

func NewSide_List(s *capnp.Segment, sz int32) (Side_List, error) {
	l, err := capnp.NewUInt16List(s, sz)
	if err != nil {
		return Side_List{}, err
	}
	return Side_List{l.List}, nil
}

func (l Side_List) At(i int) Side {
	ul := capnp.UInt16List{List: l.List}
	return Side(ul.At(i))
}

func (l Side_List) Set(i int, v Side) {
	ul := capnp.UInt16List{List: l.List}
	ul.Set(i, uint16(v))
}

type VatId struct{ capnp.Struct }

func NewVatId(s *capnp.Segment) (VatId, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	if err != nil {
		return VatId{}, err
	}
	return VatId{st}, nil
}

func NewRootVatId(s *capnp.Segment) (VatId, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	if err != nil {
		return VatId{}, err
	}
	return VatId{st}, nil
}

func ReadRootVatId(msg *capnp.Message) (VatId, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return VatId{}, err
	}
	return VatId{root.Struct()}, nil
}
func (s VatId) Side() Side {
	return Side(s.Struct.Uint16(0))
}

func (s VatId) SetSide(v Side) {
	s.Struct.SetUint16(0, uint16(v))
}

// VatId_List is a list of VatId.
type VatId_List struct{ capnp.List }

// NewVatId creates a new list of VatId.
func NewVatId_List(s *capnp.Segment, sz int32) (VatId_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0}, sz)
	if err != nil {
		return VatId_List{}, err
	}
	return VatId_List{l}, nil
}

func (s VatId_List) At(i int) VatId           { return VatId{s.List.Struct(i)} }
func (s VatId_List) Set(i int, v VatId) error { return s.List.SetStruct(i, v.Struct) }

// VatId_Promise is a wrapper for a VatId promised by a client call.
type VatId_Promise struct{ *capnp.Pipeline }

func (p VatId_Promise) Struct() (VatId, error) {
	s, err := p.Pipeline.Struct()
	return VatId{s}, err
}

type ProvisionId struct{ capnp.Struct }

func NewProvisionId(s *capnp.Segment) (ProvisionId, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	if err != nil {
		return ProvisionId{}, err
	}
	return ProvisionId{st}, nil
}

func NewRootProvisionId(s *capnp.Segment) (ProvisionId, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	if err != nil {
		return ProvisionId{}, err
	}
	return ProvisionId{st}, nil
}

func ReadRootProvisionId(msg *capnp.Message) (ProvisionId, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return ProvisionId{}, err
	}
	return ProvisionId{root.Struct()}, nil
}
func (s ProvisionId) JoinId() uint32 {
	return s.Struct.Uint32(0)
}

func (s ProvisionId) SetJoinId(v uint32) {
	s.Struct.SetUint32(0, v)
}

// ProvisionId_List is a list of ProvisionId.
type ProvisionId_List struct{ capnp.List }

// NewProvisionId creates a new list of ProvisionId.
func NewProvisionId_List(s *capnp.Segment, sz int32) (ProvisionId_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0}, sz)
	if err != nil {
		return ProvisionId_List{}, err
	}
	return ProvisionId_List{l}, nil
}

func (s ProvisionId_List) At(i int) ProvisionId           { return ProvisionId{s.List.Struct(i)} }
func (s ProvisionId_List) Set(i int, v ProvisionId) error { return s.List.SetStruct(i, v.Struct) }

// ProvisionId_Promise is a wrapper for a ProvisionId promised by a client call.
type ProvisionId_Promise struct{ *capnp.Pipeline }

func (p ProvisionId_Promise) Struct() (ProvisionId, error) {
	s, err := p.Pipeline.Struct()
	return ProvisionId{s}, err
}

type RecipientId struct{ capnp.Struct }

func NewRecipientId(s *capnp.Segment) (RecipientId, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0})
	if err != nil {
		return RecipientId{}, err
	}
	return RecipientId{st}, nil
}

func NewRootRecipientId(s *capnp.Segment) (RecipientId, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0})
	if err != nil {
		return RecipientId{}, err
	}
	return RecipientId{st}, nil
}

func ReadRootRecipientId(msg *capnp.Message) (RecipientId, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return RecipientId{}, err
	}
	return RecipientId{root.Struct()}, nil
}

// RecipientId_List is a list of RecipientId.
type RecipientId_List struct{ capnp.List }

// NewRecipientId creates a new list of RecipientId.
func NewRecipientId_List(s *capnp.Segment, sz int32) (RecipientId_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0}, sz)
	if err != nil {
		return RecipientId_List{}, err
	}
	return RecipientId_List{l}, nil
}

func (s RecipientId_List) At(i int) RecipientId           { return RecipientId{s.List.Struct(i)} }
func (s RecipientId_List) Set(i int, v RecipientId) error { return s.List.SetStruct(i, v.Struct) }

// RecipientId_Promise is a wrapper for a RecipientId promised by a client call.
type RecipientId_Promise struct{ *capnp.Pipeline }

func (p RecipientId_Promise) Struct() (RecipientId, error) {
	s, err := p.Pipeline.Struct()
	return RecipientId{s}, err
}

type ThirdPartyCapId struct{ capnp.Struct }

func NewThirdPartyCapId(s *capnp.Segment) (ThirdPartyCapId, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0})
	if err != nil {
		return ThirdPartyCapId{}, err
	}
	return ThirdPartyCapId{st}, nil
}

func NewRootThirdPartyCapId(s *capnp.Segment) (ThirdPartyCapId, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0})
	if err != nil {
		return ThirdPartyCapId{}, err
	}
	return ThirdPartyCapId{st}, nil
}

func ReadRootThirdPartyCapId(msg *capnp.Message) (ThirdPartyCapId, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return ThirdPartyCapId{}, err
	}
	return ThirdPartyCapId{root.Struct()}, nil
}

// ThirdPartyCapId_List is a list of ThirdPartyCapId.
type ThirdPartyCapId_List struct{ capnp.List }

// NewThirdPartyCapId creates a new list of ThirdPartyCapId.
func NewThirdPartyCapId_List(s *capnp.Segment, sz int32) (ThirdPartyCapId_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0}, sz)
	if err != nil {
		return ThirdPartyCapId_List{}, err
	}
	return ThirdPartyCapId_List{l}, nil
}

func (s ThirdPartyCapId_List) At(i int) ThirdPartyCapId { return ThirdPartyCapId{s.List.Struct(i)} }
func (s ThirdPartyCapId_List) Set(i int, v ThirdPartyCapId) error {
	return s.List.SetStruct(i, v.Struct)
}

// ThirdPartyCapId_Promise is a wrapper for a ThirdPartyCapId promised by a client call.
type ThirdPartyCapId_Promise struct{ *capnp.Pipeline }

func (p ThirdPartyCapId_Promise) Struct() (ThirdPartyCapId, error) {
	s, err := p.Pipeline.Struct()
	return ThirdPartyCapId{s}, err
}

type JoinKeyPart struct{ capnp.Struct }

func NewJoinKeyPart(s *capnp.Segment) (JoinKeyPart, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	if err != nil {
		return JoinKeyPart{}, err
	}
	return JoinKeyPart{st}, nil
}

func NewRootJoinKeyPart(s *capnp.Segment) (JoinKeyPart, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	if err != nil {
		return JoinKeyPart{}, err
	}
	return JoinKeyPart{st}, nil
}

func ReadRootJoinKeyPart(msg *capnp.Message) (JoinKeyPart, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return JoinKeyPart{}, err
	}
	return JoinKeyPart{root.Struct()}, nil
}
func (s JoinKeyPart) JoinId() uint32 {
	return s.Struct.Uint32(0)
}

func (s JoinKeyPart) SetJoinId(v uint32) {
	s.Struct.SetUint32(0, v)
}

func (s JoinKeyPart) PartCount() uint16 {
	return s.Struct.Uint16(4)
}

func (s JoinKeyPart) SetPartCount(v uint16) {
	s.Struct.SetUint16(4, v)
}

func (s JoinKeyPart) PartNum() uint16 {
	return s.Struct.Uint16(6)
}

func (s JoinKeyPart) SetPartNum(v uint16) {
	s.Struct.SetUint16(6, v)
}

// JoinKeyPart_List is a list of JoinKeyPart.
type JoinKeyPart_List struct{ capnp.List }

// NewJoinKeyPart creates a new list of JoinKeyPart.
func NewJoinKeyPart_List(s *capnp.Segment, sz int32) (JoinKeyPart_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0}, sz)
	if err != nil {
		return JoinKeyPart_List{}, err
	}
	return JoinKeyPart_List{l}, nil
}

func (s JoinKeyPart_List) At(i int) JoinKeyPart           { return JoinKeyPart{s.List.Struct(i)} }
func (s JoinKeyPart_List) Set(i int, v JoinKeyPart) error { return s.List.SetStruct(i, v.Struct) }

// JoinKeyPart_Promise is a wrapper for a JoinKeyPart promised by a client call.
type JoinKeyPart_Promise struct{ *capnp.Pipeline }

func (p JoinKeyPart_Promise) Struct() (JoinKeyPart, error) {
	s, err := p.Pipeline.Struct()
	return JoinKeyPart{s}, err
}

type JoinResult struct{ capnp.Struct }

func NewJoinResult(s *capnp.Segment) (JoinResult, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		return JoinResult{}, err
	}
	return JoinResult{st}, nil
}

func NewRootJoinResult(s *capnp.Segment) (JoinResult, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		return JoinResult{}, err
	}
	return JoinResult{st}, nil
}

func ReadRootJoinResult(msg *capnp.Message) (JoinResult, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return JoinResult{}, err
	}
	return JoinResult{root.Struct()}, nil
}
func (s JoinResult) JoinId() uint32 {
	return s.Struct.Uint32(0)
}

func (s JoinResult) SetJoinId(v uint32) {
	s.Struct.SetUint32(0, v)
}

func (s JoinResult) Succeeded() bool {
	return s.Struct.Bit(32)
}

func (s JoinResult) SetSucceeded(v bool) {
	s.Struct.SetBit(32, v)
}

func (s JoinResult) Cap() (capnp.Pointer, error) {
	return s.Struct.Pointer(0)
}

func (s JoinResult) HasCap() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s JoinResult) CapPtr() (capnp.Ptr, error) {
	return s.Struct.Ptr(0)
}

func (s JoinResult) SetCap(v capnp.Pointer) error {
	return s.Struct.SetPointer(0, v)
}

func (s JoinResult) SetCapPtr(v capnp.Ptr) error {
	return s.Struct.SetPtr(0, v)
}

// JoinResult_List is a list of JoinResult.
type JoinResult_List struct{ capnp.List }

// NewJoinResult creates a new list of JoinResult.
func NewJoinResult_List(s *capnp.Segment, sz int32) (JoinResult_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1}, sz)
	if err != nil {
		return JoinResult_List{}, err
	}
	return JoinResult_List{l}, nil
}

func (s JoinResult_List) At(i int) JoinResult           { return JoinResult{s.List.Struct(i)} }
func (s JoinResult_List) Set(i int, v JoinResult) error { return s.List.SetStruct(i, v.Struct) }

// JoinResult_Promise is a wrapper for a JoinResult promised by a client call.
type JoinResult_Promise struct{ *capnp.Pipeline }

func (p JoinResult_Promise) Struct() (JoinResult, error) {
	s, err := p.Pipeline.Struct()
	return JoinResult{s}, err
}

func (p JoinResult_Promise) Cap() *capnp.Pipeline {
	return p.Pipeline.GetPipeline(0)
}

const schema_a184c7885cdaf2a1 = "x\xda\x8cT]h\x1cU\x14\xbe\xe7\xce\x8c\x9b\xd5\xc4" +
	"\xcddRJ%%5\xdaJ\xfe\xba\xbb\xb1jY\x94" +
	"M[\x0b\xa6\x86\x9a\xd9D\xb1\xa5U'\xb3\xc3\xba\x9a" +
	"\x9d\x19w&?[\x91\xd5(\xda\x16*\xd5\"H\xb0" +
	"\x0d}\x10\xa9 Rik\x17\x1a\xa5\x85Z\x09(\xd8" +
	"\x87>\x08B\xfa\xacT\x9a\xbeh\xc0^\xcf\x9d\xdd\xd9" +
	"\x19\x9a\xa4\xf4a\xd9{\xcf\xfd\xee\xb9\xdf\xf9\xcew&" +
	"q\x1d\xfaiR*\x8b\x84\xa8[\xa4\xfb\xd8\xc5\xe9\xfe" +
	"\xff\xce\x1d\xbau\x88\xc8\x0a\xb0\x93\x8b\xbf\xef;x\xe5" +
	"\x83\x93D\x8c\x10\xa2l\xa4KJ\x92F\x88\xc0\xde\x9f" +
	"/G\xf7\x1c\xfd\xee3\xa2*p'\xaa\x09Q\xeb(" +
	"_\xad\xa1\xdf\x12`7\xfe\xfc>\xf1dj\xd3\xf1;" +
	"\xb0\x12p\xc8i\xba\xa8\xccy\xe0\x8a\x07\xfe\xe2\xad\xca" +
	"G\x17N\\\x9b\xc5\xd7i\x80%\xa0\xbc(\\R\xf6" +
	"\x0b\x1c\xb8Gx\x8a\x84\xf2@\x13\xa1x2\xa3\x18\xc2" +
	"\xa4\xf2\x8b\xf0\x02\xa1\xec\xf8\xd9\x9e\\i\xa0|fy" +
	"\x09\x8f'\xc5\x87@\xd9&\xf2\x1a\xd6~~\xeb\xf2\xa9" +
	"\xe8\x91\xf3+\xd5\xb0^\\R:\xab5\x8bi|\xad" +
	"\xe7\xd5w\x7f<\xf1\xefO\x15\xa26K\xc0\xde\xdb\x9d" +
	"8\xff\xf6;\xff\xccq^;\xc5\xab\x8a\xca\x91\xc3\x83" +
	"\xa2\x00$t\x08\xf7#\xb1g\xc4i|\xee1e\xbf" +
	"\xc8\x89-t\x97\x1ei.\x9f\xfa\x81\xa8QLsp" +
	"\xec\x8f\xdbj[\xd7o<ME<\xac\\\xf4\xd2\\" +
	"\xa8\xa6\xa9\x1f\x82\x80i\xbe\x11w)\xa7\xc5I%*" +
	"\xb5c\x9a\x8fS\xce\x8d\xd9\xa3\x0f\\]\x89yT\x9a" +
	"W\xd6H|%K\x9c\xf9\xd9g\x1f\xdc\x04\xe7\x12\xd7" +
	"\x97?\x99\x94\xa6\x95'8r8!yO\x16m\xbd" +
	"\xd7\x9d\xb4l\xaa\x15\xdd\xd2f]\xb3M;\x951\xf4" +
	"\xbc\x9d\xce\x1b\xa6;\x90\x1d\x82\x951\xbb\xac\xbc\xf9|" +
	"\xda(\x0da\x0c1j\xa3\x80^\x12\x81\x10yg\x0a" +
	"M\xd5/\x80:HA\x06\xda\x0a<8\x90\xc1\xe0s" +
	"\x18\x1c\xc1 \x15Z\x81bP\xdd\x8e\xc1A\x0c\xbeL" +
	"!\xfd\x06&\x1c\xc8B\x03\xa1\xf8\x03fc\xde\x1d\xd6" +
	"\xb8I\xc0E\xcfP\xfcA\x99\xc7v\x8f\x17\xfc\xfd\xaa" +
	"\xbc2\xed\x863>vWZ\x1b\x96\xd3\xe2T=V" +
	"\x1d\xab\xb3r\xc6u\xdd0\xb2\x06\x81,\xde\xc7;\x04" +
	"\"\xf8.\xb4\xe0\xbae\x15F\xc3y\xc4s.\x0d^" +
	"z\x19y\x00\xc8Q\xfcK;Fq\xc2(\xa6\xf51" +
	"\xaeu\xfd2\xf8\x97\x05\xd3V\xdb <\x1f\xc9\xae\xc0" +
	"\x07ro_`g\xb9s44\xc4\xb8\x09\xc6\xa1s" +
	":\x98\\~R\x1fM\xb9so\x8c\x93k\x7fI\xc3" +
	"F\xb3\xa1\xa25\x91w\xf2\x16\x89`\xc5\xcc\xf3\x00\xd2" +
	"\"\x11~6\xf2z\xbe\x98\xc5VS\xb7\xb4C\xb3\x07" +
	"\xb2\xc4k?v?\x82\\\x99\xa79JN\x841W" +
	"m\x04\x1a\x9a\x9c6\x0a\xdb\x12\x10\x9a\x81\x1e\x0c<\x1d" +
	"v\xe8V\x0c\xec\x03hDi\xe0kV\x95,U\x04" +
	"[O\xa5\xb8\x181.\x06*\xe6\x9d\x8fr\x8dxT" +
	"#\x91z8\x09\x9f\x02;`\x15F\xf3\xc6\x01C2" +
	"7\xebV!\x9e\xb3\xe2^\xa6\xa2\xe5Z}q\xc7\xcd" +
	"V\xb7q\xffz\xd1\x05\xbc]W\\\x08\xb5\xcb/\xd5" +
	"\xabTX}\x00<\xb9\xd2\xa8\x97\xe9aT\xb1\xee\xb4" +
	"&\xee\xb4\x064P\xebr\x03Ui\xe8\xdd\xb4\xbb\x96" +
	"\xc6\xd4\x0a\x86ck:\x18\x98\x04\xcb\xa1u\x08\xd4 " +
	"\x98\x0bB\x1f#9\x9aa\xbf~r\xf3v\xe9\xab\xec" +
	"\"n\xba\x98\x9f\x81\x80\x11\xe3\xeb{\xe9\xc0\xd6p\x07" +
	"\xb6``\xa4\xd6\x81\x19\xbf\x03\x1a\x98\xa6\xe5jn\xde" +
	"\x12L\xc7\xd7\xbf#\xa2OM\xf9\xaag\x02\xd5\xc5\xbb" +
	"\xab\xaeO\x01\xdeb9\xabZ\x10\xa4\x90\xed\x9bZ\xce" +
	" \xc4\xaf\xd9?\"\xa8\x1cw|\x9d\xad\x9c\xdc\x1e0" +
	"\x95{S\xec\x95c\xb3\xea\xdc\xb5\xc3\x97\xd1\xbd\x1d\xec" +
	"\xca\xdf\xf3\x8f\xae;\xe3~I\xe4\x8d\x1d\xace!\xf3" +
	"W\xe9\xc3\x89\x9f\x89\xfcp\x1f;\xb6!\xbe0c4" +
	"/\x11y\xfd^v\xf3H|m\xcbk\x95K\xb8\xe9" +
	"*\xd7\xdeN\xe7\x0b\xb6Ut#YK\x8f\xb8Z\xae" +
	"\x9d\x17\x9bc\xfa\xb8\xe3Z\x05t\x96`\xd7\xc4\x14!" +
	"\xfc\x05\x17Q\xaa\xe6\x9aT}\xed5\xca+\xda\xc3\x9b" +
	")\xaf\xc0\xb03\xba\x02g\xc4\x1c\x1c=\x88\x05\xc3\x8d" +
	"_\x85\x18:$P\xa9\xca\xd0\xd7\xe8\xff\x00\x00\x00\xff" +
	"\xff\xf90^\xe8"
