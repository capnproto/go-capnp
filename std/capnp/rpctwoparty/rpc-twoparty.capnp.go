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

var schema_a184c7885cdaf2a1 = []byte{
	120, 218, 140, 84, 93, 104, 28, 85,
	20, 190, 231, 206, 140, 155, 104, 235,
	102, 50, 129, 82, 73, 73, 140, 166,
	146, 191, 110, 54, 22, 173, 131, 176,
	177, 53, 96, 106, 169, 153, 77, 16,
	45, 69, 157, 204, 14, 235, 106, 118,
	102, 156, 157, 252, 108, 69, 214, 196,
	191, 180, 16, 105, 45, 1, 27, 76,
	99, 30, 130, 212, 23, 65, 218, 218,
	128, 81, 82, 168, 213, 60, 8, 246,
	161, 15, 130, 144, 62, 43, 149, 164,
	47, 54, 208, 94, 207, 157, 236, 238,
	12, 73, 12, 62, 12, 115, 239, 185,
	231, 156, 251, 157, 239, 126, 231, 180,
	223, 130, 78, 26, 151, 10, 34, 33,
	218, 126, 233, 1, 54, 122, 180, 253,
	202, 187, 239, 253, 179, 64, 224, 65,
	66, 149, 70, 58, 166, 52, 209, 39,
	148, 46, 250, 18, 161, 236, 210, 243,
	15, 239, 133, 203, 237, 183, 136, 86,
	41, 1, 27, 31, 248, 227, 190, 86,
	219, 252, 27, 33, 160, 204, 162, 223,
	28, 141, 16, 210, 251, 37, 21, 128,
	0, 251, 84, 205, 221, 158, 57, 253,
	208, 13, 162, 41, 0, 108, 118, 245,
	247, 227, 227, 215, 63, 156, 37, 34,
	186, 40, 103, 232, 146, 50, 205, 157,
	149, 115, 52, 129, 190, 95, 188, 51,
	255, 201, 247, 231, 111, 206, 16, 89,
	161, 129, 43, 102, 157, 167, 87, 149,
	69, 223, 113, 129, 62, 77, 66, 105,
	96, 39, 98, 91, 164, 83, 202, 47,
	116, 88, 137, 11, 28, 91, 25, 13,
	8, 120, 54, 41, 28, 86, 206, 9,
	195, 202, 138, 80, 135, 103, 31, 44,
	21, 42, 95, 61, 253, 237, 228, 86,
	96, 86, 132, 53, 229, 158, 192, 87,
	119, 133, 111, 240, 142, 233, 75, 173,
	233, 124, 119, 225, 34, 130, 217, 224,
	250, 228, 164, 248, 8, 40, 115, 184,
	20, 88, 235, 235, 239, 255, 120, 254,
	238, 79, 243, 68, 171, 66, 38, 202,
	156, 33, 230, 81, 241, 134, 50, 193,
	221, 123, 199, 69, 159, 137, 219, 127,
	126, 215, 254, 148, 186, 119, 122, 195,
	229, 18, 240, 59, 243, 226, 170, 242,
	145, 143, 99, 84, 228, 183, 47, 183,
	228, 31, 171, 42, 92, 248, 97, 51,
	195, 123, 164, 83, 74, 163, 196, 243,
	214, 75, 126, 222, 93, 159, 223, 185,
	118, 161, 114, 226, 202, 86, 69, 201,
	210, 26, 6, 240, 213, 110, 137, 51,
	188, 56, 214, 121, 239, 242, 201, 59,
	39, 55, 23, 165, 60, 131, 174, 93,
	18, 175, 201, 208, 29, 203, 137, 25,
	45, 208, 178, 207, 95, 18, 173, 2,
	19, 151, 43, 149, 43, 147, 236, 215,
	51, 43, 247, 243, 95, 165, 86, 113,
	211, 204, 44, 61, 107, 230, 28, 221,
	32, 96, 70, 249, 90, 219, 1, 52,
	68, 76, 45, 133, 231, 218, 195, 21,
	181, 162, 225, 0, 132, 68, 180, 31,
	13, 125, 0, 59, 8, 145, 97, 106,
	253, 122, 85, 213, 193, 178, 108, 79,
	247, 50, 182, 96, 229, 176, 112, 255,
	180, 33, 98, 140, 140, 172, 111, 226,
	144, 4, 118, 194, 206, 246, 103, 204,
	19, 166, 104, 237, 51, 236, 108, 44,
	109, 199, 252, 104, 215, 246, 236, 142,
	88, 206, 75, 197, 138, 181, 140, 0,
	70, 177, 180, 189, 94, 16, 168, 153,
	172, 99, 187, 30, 233, 1, 126, 43,
	101, 174, 99, 180, 121, 195, 182, 67,
	117, 215, 203, 175, 251, 168, 47, 235,
	94, 119, 138, 112, 23, 77, 20, 176,
	41, 68, 64, 0, 59, 155, 177, 59,
	42, 4, 208, 106, 40, 68, 115, 153,
	148, 9, 209, 64, 181, 4, 32, 138,
	28, 111, 149, 172, 23, 93, 253, 84,
	21, 64, 49, 141, 172, 114, 103, 185,
	18, 127, 137, 156, 233, 14, 153, 110,
	194, 24, 200, 152, 150, 87, 14, 134,
	82, 176, 96, 57, 90, 45, 132, 123,
	35, 222, 28, 52, 149, 220, 214, 17,
	188, 191, 220, 212, 31, 122, 96, 220,
	4, 18, 110, 26, 11, 180, 207, 79,
	202, 90, 148, 155, 142, 69, 57, 184,
	58, 191, 92, 214, 227, 218, 67, 153,
	92, 198, 38, 17, 11, 119, 73, 211,
	200, 56, 8, 139, 68, 248, 89, 223,
	155, 25, 55, 213, 163, 187, 212, 203,
	31, 210, 29, 228, 230, 176, 157, 177,
	94, 52, 243, 61, 17, 196, 202, 248,
	38, 105, 230, 6, 137, 48, 224, 253,
	31, 1, 60, 27, 22, 192, 1, 52,
	28, 47, 10, 224, 235, 146, 0, 92,
	112, 12, 85, 229, 100, 68, 57, 25,
	37, 9, 244, 115, 142, 184, 85, 39,
	145, 178, 57, 14, 159, 5, 98, 144,
	182, 23, 67, 41, 220, 245, 32, 31,
	18, 5, 193, 151, 229, 68, 151, 97,
	202, 241, 131, 1, 68, 185, 77, 101,
	175, 157, 157, 209, 22, 110, 158, 186,
	134, 164, 53, 176, 235, 127, 47, 61,
	190, 251, 162, 55, 71, 228, 198, 6,
	86, 189, 156, 252, 43, 255, 241, 208,
	207, 68, 126, 180, 131, 157, 173, 143,
	45, 79, 153, 85, 107, 68, 222, 115,
	140, 173, 76, 196, 118, 85, 191, 49,
	127, 21, 55, 205, 5, 236, 145, 183,
	245, 180, 153, 88, 87, 95, 36, 101,
	27, 17, 79, 79, 215, 113, 153, 167,
	153, 49, 152, 243, 236, 44, 22, 36,
	56, 197, 54, 18, 145, 197, 128, 52,
	17, 57, 170, 42, 114, 212, 81, 87,
	132, 188, 165, 216, 252, 119, 73, 224,
	195, 160, 141, 107, 110, 71, 89, 190,
	93, 40, 56, 173, 19, 229, 123, 132,
	130, 12, 180, 6, 184, 177, 59, 137,
	198, 23, 208, 216, 135, 70, 42, 212,
	248, 10, 213, 14, 162, 241, 8, 26,
	95, 161, 144, 120, 11, 19, 118, 167,
	160, 130, 80, 252, 128, 113, 234, 14,
	217, 131, 22, 1, 15, 103, 23, 197,
	15, 10, 220, 118, 116, 48, 91, 218,
	151, 113, 9, 33, 92, 37, 1, 249,
	250, 17, 186, 83, 136, 173, 60, 107,
	104, 113, 214, 168, 165, 97, 2, 230,
	118, 205, 233, 171, 173, 14, 229, 54,
	176, 109, 133, 245, 155, 43, 228, 85,
	251, 5, 54, 252, 119, 129, 185, 65,
	195, 48, 77, 236, 88, 72, 97, 60,
	198, 16, 136, 224, 189, 80, 141, 235,
	106, 60, 15, 230, 72, 241, 69, 201,
	182, 147, 196, 111, 170, 4, 118, 149,
	229, 215, 28, 30, 39, 106, 48, 78,
	54, 130, 216, 42, 147, 223, 144, 9,
	62, 40, 252, 76, 255, 6, 0, 0,
	255, 255, 229, 6, 101, 60,
}
