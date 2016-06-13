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
	return Side_List{l.List}, err
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
	return VatId{st}, err
}

func NewRootVatId(s *capnp.Segment) (VatId, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	return VatId{st}, err
}

func ReadRootVatId(msg *capnp.Message) (VatId, error) {
	root, err := msg.RootPtr()
	return VatId{root.Struct()}, err
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
	return VatId_List{l}, err
}

func (s VatId_List) At(i int) VatId { return VatId{s.List.Struct(i)} }

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
	return ProvisionId{st}, err
}

func NewRootProvisionId(s *capnp.Segment) (ProvisionId, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	return ProvisionId{st}, err
}

func ReadRootProvisionId(msg *capnp.Message) (ProvisionId, error) {
	root, err := msg.RootPtr()
	return ProvisionId{root.Struct()}, err
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
	return ProvisionId_List{l}, err
}

func (s ProvisionId_List) At(i int) ProvisionId { return ProvisionId{s.List.Struct(i)} }

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
	return RecipientId{st}, err
}

func NewRootRecipientId(s *capnp.Segment) (RecipientId, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0})
	return RecipientId{st}, err
}

func ReadRootRecipientId(msg *capnp.Message) (RecipientId, error) {
	root, err := msg.RootPtr()
	return RecipientId{root.Struct()}, err
}

// RecipientId_List is a list of RecipientId.
type RecipientId_List struct{ capnp.List }

// NewRecipientId creates a new list of RecipientId.
func NewRecipientId_List(s *capnp.Segment, sz int32) (RecipientId_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0}, sz)
	return RecipientId_List{l}, err
}

func (s RecipientId_List) At(i int) RecipientId { return RecipientId{s.List.Struct(i)} }

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
	return ThirdPartyCapId{st}, err
}

func NewRootThirdPartyCapId(s *capnp.Segment) (ThirdPartyCapId, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0})
	return ThirdPartyCapId{st}, err
}

func ReadRootThirdPartyCapId(msg *capnp.Message) (ThirdPartyCapId, error) {
	root, err := msg.RootPtr()
	return ThirdPartyCapId{root.Struct()}, err
}

// ThirdPartyCapId_List is a list of ThirdPartyCapId.
type ThirdPartyCapId_List struct{ capnp.List }

// NewThirdPartyCapId creates a new list of ThirdPartyCapId.
func NewThirdPartyCapId_List(s *capnp.Segment, sz int32) (ThirdPartyCapId_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0}, sz)
	return ThirdPartyCapId_List{l}, err
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
	return JoinKeyPart{st}, err
}

func NewRootJoinKeyPart(s *capnp.Segment) (JoinKeyPart, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	return JoinKeyPart{st}, err
}

func ReadRootJoinKeyPart(msg *capnp.Message) (JoinKeyPart, error) {
	root, err := msg.RootPtr()
	return JoinKeyPart{root.Struct()}, err
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
	return JoinKeyPart_List{l}, err
}

func (s JoinKeyPart_List) At(i int) JoinKeyPart { return JoinKeyPart{s.List.Struct(i)} }

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
	return JoinResult{st}, err
}

func NewRootJoinResult(s *capnp.Segment) (JoinResult, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return JoinResult{st}, err
}

func ReadRootJoinResult(msg *capnp.Message) (JoinResult, error) {
	root, err := msg.RootPtr()
	return JoinResult{root.Struct()}, err
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
	return JoinResult_List{l}, err
}

func (s JoinResult_List) At(i int) JoinResult { return JoinResult{s.List.Struct(i)} }

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
