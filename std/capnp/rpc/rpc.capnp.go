package rpc

// AUTO GENERATED - DO NOT EDIT

import (
	strconv "strconv"
	capnp "zombiezen.com/go/capnproto2"
	text "zombiezen.com/go/capnproto2/encoding/text"
	schemas "zombiezen.com/go/capnproto2/schemas"
)

type Message struct{ capnp.Struct }
type Message_Which uint16

const (
	Message_Which_unimplemented  Message_Which = 0
	Message_Which_abort          Message_Which = 1
	Message_Which_bootstrap      Message_Which = 8
	Message_Which_call           Message_Which = 2
	Message_Which_return         Message_Which = 3
	Message_Which_finish         Message_Which = 4
	Message_Which_resolve        Message_Which = 5
	Message_Which_release        Message_Which = 6
	Message_Which_disembargo     Message_Which = 13
	Message_Which_obsoleteSave   Message_Which = 7
	Message_Which_obsoleteDelete Message_Which = 9
	Message_Which_provide        Message_Which = 10
	Message_Which_accept         Message_Which = 11
	Message_Which_join           Message_Which = 12
)

func (w Message_Which) String() string {
	const s = "unimplementedabortbootstrapcallreturnfinishresolvereleasedisembargoobsoleteSaveobsoleteDeleteprovideacceptjoin"
	switch w {
	case Message_Which_unimplemented:
		return s[0:13]
	case Message_Which_abort:
		return s[13:18]
	case Message_Which_bootstrap:
		return s[18:27]
	case Message_Which_call:
		return s[27:31]
	case Message_Which_return:
		return s[31:37]
	case Message_Which_finish:
		return s[37:43]
	case Message_Which_resolve:
		return s[43:50]
	case Message_Which_release:
		return s[50:57]
	case Message_Which_disembargo:
		return s[57:67]
	case Message_Which_obsoleteSave:
		return s[67:79]
	case Message_Which_obsoleteDelete:
		return s[79:93]
	case Message_Which_provide:
		return s[93:100]
	case Message_Which_accept:
		return s[100:106]
	case Message_Which_join:
		return s[106:110]

	}
	return "Message_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
}

func NewMessage(s *capnp.Segment) (Message, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return Message{st}, err
}

func NewRootMessage(s *capnp.Segment) (Message, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return Message{st}, err
}

func ReadRootMessage(msg *capnp.Message) (Message, error) {
	root, err := msg.RootPtr()
	return Message{root.Struct()}, err
}

func (s Message) String() string {
	str, _ := text.Marshal(0x91b79f1f808db032, s.Struct)
	return str
}

func (s Message) Which() Message_Which {
	return Message_Which(s.Struct.Uint16(0))
}
func (s Message) Unimplemented() (Message, error) {
	p, err := s.Struct.Ptr(0)
	return Message{Struct: p.Struct()}, err
}

func (s Message) HasUnimplemented() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Message) SetUnimplemented(v Message) error {
	s.Struct.SetUint16(0, 0)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewUnimplemented sets the unimplemented field to a newly
// allocated Message struct, preferring placement in s's segment.
func (s Message) NewUnimplemented() (Message, error) {
	s.Struct.SetUint16(0, 0)
	ss, err := NewMessage(s.Struct.Segment())
	if err != nil {
		return Message{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Message) Abort() (Exception, error) {
	p, err := s.Struct.Ptr(0)
	return Exception{Struct: p.Struct()}, err
}

func (s Message) HasAbort() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Message) SetAbort(v Exception) error {
	s.Struct.SetUint16(0, 1)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewAbort sets the abort field to a newly
// allocated Exception struct, preferring placement in s's segment.
func (s Message) NewAbort() (Exception, error) {
	s.Struct.SetUint16(0, 1)
	ss, err := NewException(s.Struct.Segment())
	if err != nil {
		return Exception{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Message) Bootstrap() (Bootstrap, error) {
	p, err := s.Struct.Ptr(0)
	return Bootstrap{Struct: p.Struct()}, err
}

func (s Message) HasBootstrap() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Message) SetBootstrap(v Bootstrap) error {
	s.Struct.SetUint16(0, 8)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewBootstrap sets the bootstrap field to a newly
// allocated Bootstrap struct, preferring placement in s's segment.
func (s Message) NewBootstrap() (Bootstrap, error) {
	s.Struct.SetUint16(0, 8)
	ss, err := NewBootstrap(s.Struct.Segment())
	if err != nil {
		return Bootstrap{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Message) Call() (Call, error) {
	p, err := s.Struct.Ptr(0)
	return Call{Struct: p.Struct()}, err
}

func (s Message) HasCall() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Message) SetCall(v Call) error {
	s.Struct.SetUint16(0, 2)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewCall sets the call field to a newly
// allocated Call struct, preferring placement in s's segment.
func (s Message) NewCall() (Call, error) {
	s.Struct.SetUint16(0, 2)
	ss, err := NewCall(s.Struct.Segment())
	if err != nil {
		return Call{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Message) Return() (Return, error) {
	p, err := s.Struct.Ptr(0)
	return Return{Struct: p.Struct()}, err
}

func (s Message) HasReturn() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Message) SetReturn(v Return) error {
	s.Struct.SetUint16(0, 3)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewReturn sets the return field to a newly
// allocated Return struct, preferring placement in s's segment.
func (s Message) NewReturn() (Return, error) {
	s.Struct.SetUint16(0, 3)
	ss, err := NewReturn(s.Struct.Segment())
	if err != nil {
		return Return{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Message) Finish() (Finish, error) {
	p, err := s.Struct.Ptr(0)
	return Finish{Struct: p.Struct()}, err
}

func (s Message) HasFinish() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Message) SetFinish(v Finish) error {
	s.Struct.SetUint16(0, 4)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewFinish sets the finish field to a newly
// allocated Finish struct, preferring placement in s's segment.
func (s Message) NewFinish() (Finish, error) {
	s.Struct.SetUint16(0, 4)
	ss, err := NewFinish(s.Struct.Segment())
	if err != nil {
		return Finish{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Message) Resolve() (Resolve, error) {
	p, err := s.Struct.Ptr(0)
	return Resolve{Struct: p.Struct()}, err
}

func (s Message) HasResolve() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Message) SetResolve(v Resolve) error {
	s.Struct.SetUint16(0, 5)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewResolve sets the resolve field to a newly
// allocated Resolve struct, preferring placement in s's segment.
func (s Message) NewResolve() (Resolve, error) {
	s.Struct.SetUint16(0, 5)
	ss, err := NewResolve(s.Struct.Segment())
	if err != nil {
		return Resolve{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Message) Release() (Release, error) {
	p, err := s.Struct.Ptr(0)
	return Release{Struct: p.Struct()}, err
}

func (s Message) HasRelease() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Message) SetRelease(v Release) error {
	s.Struct.SetUint16(0, 6)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewRelease sets the release field to a newly
// allocated Release struct, preferring placement in s's segment.
func (s Message) NewRelease() (Release, error) {
	s.Struct.SetUint16(0, 6)
	ss, err := NewRelease(s.Struct.Segment())
	if err != nil {
		return Release{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Message) Disembargo() (Disembargo, error) {
	p, err := s.Struct.Ptr(0)
	return Disembargo{Struct: p.Struct()}, err
}

func (s Message) HasDisembargo() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Message) SetDisembargo(v Disembargo) error {
	s.Struct.SetUint16(0, 13)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewDisembargo sets the disembargo field to a newly
// allocated Disembargo struct, preferring placement in s's segment.
func (s Message) NewDisembargo() (Disembargo, error) {
	s.Struct.SetUint16(0, 13)
	ss, err := NewDisembargo(s.Struct.Segment())
	if err != nil {
		return Disembargo{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Message) ObsoleteSave() (capnp.Pointer, error) {
	return s.Struct.Pointer(0)
}

func (s Message) HasObsoleteSave() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Message) ObsoleteSavePtr() (capnp.Ptr, error) {
	return s.Struct.Ptr(0)
}

func (s Message) SetObsoleteSave(v capnp.Pointer) error {
	s.Struct.SetUint16(0, 7)
	return s.Struct.SetPointer(0, v)
}

func (s Message) SetObsoleteSavePtr(v capnp.Ptr) error {
	s.Struct.SetUint16(0, 7)
	return s.Struct.SetPtr(0, v)
}

func (s Message) ObsoleteDelete() (capnp.Pointer, error) {
	return s.Struct.Pointer(0)
}

func (s Message) HasObsoleteDelete() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Message) ObsoleteDeletePtr() (capnp.Ptr, error) {
	return s.Struct.Ptr(0)
}

func (s Message) SetObsoleteDelete(v capnp.Pointer) error {
	s.Struct.SetUint16(0, 9)
	return s.Struct.SetPointer(0, v)
}

func (s Message) SetObsoleteDeletePtr(v capnp.Ptr) error {
	s.Struct.SetUint16(0, 9)
	return s.Struct.SetPtr(0, v)
}

func (s Message) Provide() (Provide, error) {
	p, err := s.Struct.Ptr(0)
	return Provide{Struct: p.Struct()}, err
}

func (s Message) HasProvide() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Message) SetProvide(v Provide) error {
	s.Struct.SetUint16(0, 10)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewProvide sets the provide field to a newly
// allocated Provide struct, preferring placement in s's segment.
func (s Message) NewProvide() (Provide, error) {
	s.Struct.SetUint16(0, 10)
	ss, err := NewProvide(s.Struct.Segment())
	if err != nil {
		return Provide{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Message) Accept() (Accept, error) {
	p, err := s.Struct.Ptr(0)
	return Accept{Struct: p.Struct()}, err
}

func (s Message) HasAccept() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Message) SetAccept(v Accept) error {
	s.Struct.SetUint16(0, 11)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewAccept sets the accept field to a newly
// allocated Accept struct, preferring placement in s's segment.
func (s Message) NewAccept() (Accept, error) {
	s.Struct.SetUint16(0, 11)
	ss, err := NewAccept(s.Struct.Segment())
	if err != nil {
		return Accept{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Message) Join() (Join, error) {
	p, err := s.Struct.Ptr(0)
	return Join{Struct: p.Struct()}, err
}

func (s Message) HasJoin() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Message) SetJoin(v Join) error {
	s.Struct.SetUint16(0, 12)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewJoin sets the join field to a newly
// allocated Join struct, preferring placement in s's segment.
func (s Message) NewJoin() (Join, error) {
	s.Struct.SetUint16(0, 12)
	ss, err := NewJoin(s.Struct.Segment())
	if err != nil {
		return Join{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

// Message_List is a list of Message.
type Message_List struct{ capnp.List }

// NewMessage creates a new list of Message.
func NewMessage_List(s *capnp.Segment, sz int32) (Message_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1}, sz)
	return Message_List{l}, err
}

func (s Message_List) At(i int) Message { return Message{s.List.Struct(i)} }

func (s Message_List) Set(i int, v Message) error { return s.List.SetStruct(i, v.Struct) }

// Message_Promise is a wrapper for a Message promised by a client call.
type Message_Promise struct{ *capnp.Pipeline }

func (p Message_Promise) Struct() (Message, error) {
	s, err := p.Pipeline.Struct()
	return Message{s}, err
}

func (p Message_Promise) Unimplemented() Message_Promise {
	return Message_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Message_Promise) Abort() Exception_Promise {
	return Exception_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Message_Promise) Bootstrap() Bootstrap_Promise {
	return Bootstrap_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Message_Promise) Call() Call_Promise {
	return Call_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Message_Promise) Return() Return_Promise {
	return Return_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Message_Promise) Finish() Finish_Promise {
	return Finish_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Message_Promise) Resolve() Resolve_Promise {
	return Resolve_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Message_Promise) Release() Release_Promise {
	return Release_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Message_Promise) Disembargo() Disembargo_Promise {
	return Disembargo_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Message_Promise) ObsoleteSave() *capnp.Pipeline {
	return p.Pipeline.GetPipeline(0)
}

func (p Message_Promise) ObsoleteDelete() *capnp.Pipeline {
	return p.Pipeline.GetPipeline(0)
}

func (p Message_Promise) Provide() Provide_Promise {
	return Provide_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Message_Promise) Accept() Accept_Promise {
	return Accept_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Message_Promise) Join() Join_Promise {
	return Join_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

type Bootstrap struct{ capnp.Struct }

func NewBootstrap(s *capnp.Segment) (Bootstrap, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return Bootstrap{st}, err
}

func NewRootBootstrap(s *capnp.Segment) (Bootstrap, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return Bootstrap{st}, err
}

func ReadRootBootstrap(msg *capnp.Message) (Bootstrap, error) {
	root, err := msg.RootPtr()
	return Bootstrap{root.Struct()}, err
}

func (s Bootstrap) String() string {
	str, _ := text.Marshal(0xe94ccf8031176ec4, s.Struct)
	return str
}

func (s Bootstrap) QuestionId() uint32 {
	return s.Struct.Uint32(0)
}

func (s Bootstrap) SetQuestionId(v uint32) {
	s.Struct.SetUint32(0, v)
}

func (s Bootstrap) DeprecatedObjectId() (capnp.Pointer, error) {
	return s.Struct.Pointer(0)
}

func (s Bootstrap) HasDeprecatedObjectId() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Bootstrap) DeprecatedObjectIdPtr() (capnp.Ptr, error) {
	return s.Struct.Ptr(0)
}

func (s Bootstrap) SetDeprecatedObjectId(v capnp.Pointer) error {
	return s.Struct.SetPointer(0, v)
}

func (s Bootstrap) SetDeprecatedObjectIdPtr(v capnp.Ptr) error {
	return s.Struct.SetPtr(0, v)
}

// Bootstrap_List is a list of Bootstrap.
type Bootstrap_List struct{ capnp.List }

// NewBootstrap creates a new list of Bootstrap.
func NewBootstrap_List(s *capnp.Segment, sz int32) (Bootstrap_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1}, sz)
	return Bootstrap_List{l}, err
}

func (s Bootstrap_List) At(i int) Bootstrap { return Bootstrap{s.List.Struct(i)} }

func (s Bootstrap_List) Set(i int, v Bootstrap) error { return s.List.SetStruct(i, v.Struct) }

// Bootstrap_Promise is a wrapper for a Bootstrap promised by a client call.
type Bootstrap_Promise struct{ *capnp.Pipeline }

func (p Bootstrap_Promise) Struct() (Bootstrap, error) {
	s, err := p.Pipeline.Struct()
	return Bootstrap{s}, err
}

func (p Bootstrap_Promise) DeprecatedObjectId() *capnp.Pipeline {
	return p.Pipeline.GetPipeline(0)
}

type Call struct{ capnp.Struct }
type Call_sendResultsTo Call
type Call_sendResultsTo_Which uint16

const (
	Call_sendResultsTo_Which_caller     Call_sendResultsTo_Which = 0
	Call_sendResultsTo_Which_yourself   Call_sendResultsTo_Which = 1
	Call_sendResultsTo_Which_thirdParty Call_sendResultsTo_Which = 2
)

func (w Call_sendResultsTo_Which) String() string {
	const s = "calleryourselfthirdParty"
	switch w {
	case Call_sendResultsTo_Which_caller:
		return s[0:6]
	case Call_sendResultsTo_Which_yourself:
		return s[6:14]
	case Call_sendResultsTo_Which_thirdParty:
		return s[14:24]

	}
	return "Call_sendResultsTo_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
}

func NewCall(s *capnp.Segment) (Call, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 24, PointerCount: 3})
	return Call{st}, err
}

func NewRootCall(s *capnp.Segment) (Call, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 24, PointerCount: 3})
	return Call{st}, err
}

func ReadRootCall(msg *capnp.Message) (Call, error) {
	root, err := msg.RootPtr()
	return Call{root.Struct()}, err
}

func (s Call) String() string {
	str, _ := text.Marshal(0x836a53ce789d4cd4, s.Struct)
	return str
}

func (s Call) QuestionId() uint32 {
	return s.Struct.Uint32(0)
}

func (s Call) SetQuestionId(v uint32) {
	s.Struct.SetUint32(0, v)
}

func (s Call) Target() (MessageTarget, error) {
	p, err := s.Struct.Ptr(0)
	return MessageTarget{Struct: p.Struct()}, err
}

func (s Call) HasTarget() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Call) SetTarget(v MessageTarget) error {
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewTarget sets the target field to a newly
// allocated MessageTarget struct, preferring placement in s's segment.
func (s Call) NewTarget() (MessageTarget, error) {
	ss, err := NewMessageTarget(s.Struct.Segment())
	if err != nil {
		return MessageTarget{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Call) InterfaceId() uint64 {
	return s.Struct.Uint64(8)
}

func (s Call) SetInterfaceId(v uint64) {
	s.Struct.SetUint64(8, v)
}

func (s Call) MethodId() uint16 {
	return s.Struct.Uint16(4)
}

func (s Call) SetMethodId(v uint16) {
	s.Struct.SetUint16(4, v)
}

func (s Call) AllowThirdPartyTailCall() bool {
	return s.Struct.Bit(128)
}

func (s Call) SetAllowThirdPartyTailCall(v bool) {
	s.Struct.SetBit(128, v)
}

func (s Call) Params() (Payload, error) {
	p, err := s.Struct.Ptr(1)
	return Payload{Struct: p.Struct()}, err
}

func (s Call) HasParams() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s Call) SetParams(v Payload) error {
	return s.Struct.SetPtr(1, v.Struct.ToPtr())
}

// NewParams sets the params field to a newly
// allocated Payload struct, preferring placement in s's segment.
func (s Call) NewParams() (Payload, error) {
	ss, err := NewPayload(s.Struct.Segment())
	if err != nil {
		return Payload{}, err
	}
	err = s.Struct.SetPtr(1, ss.Struct.ToPtr())
	return ss, err
}

func (s Call) SendResultsTo() Call_sendResultsTo { return Call_sendResultsTo(s) }

func (s Call_sendResultsTo) Which() Call_sendResultsTo_Which {
	return Call_sendResultsTo_Which(s.Struct.Uint16(6))
}
func (s Call_sendResultsTo) SetCaller() {
	s.Struct.SetUint16(6, 0)

}

func (s Call_sendResultsTo) SetYourself() {
	s.Struct.SetUint16(6, 1)

}

func (s Call_sendResultsTo) ThirdParty() (capnp.Pointer, error) {
	return s.Struct.Pointer(2)
}

func (s Call_sendResultsTo) HasThirdParty() bool {
	p, err := s.Struct.Ptr(2)
	return p.IsValid() || err != nil
}

func (s Call_sendResultsTo) ThirdPartyPtr() (capnp.Ptr, error) {
	return s.Struct.Ptr(2)
}

func (s Call_sendResultsTo) SetThirdParty(v capnp.Pointer) error {
	s.Struct.SetUint16(6, 2)
	return s.Struct.SetPointer(2, v)
}

func (s Call_sendResultsTo) SetThirdPartyPtr(v capnp.Ptr) error {
	s.Struct.SetUint16(6, 2)
	return s.Struct.SetPtr(2, v)
}

// Call_List is a list of Call.
type Call_List struct{ capnp.List }

// NewCall creates a new list of Call.
func NewCall_List(s *capnp.Segment, sz int32) (Call_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 24, PointerCount: 3}, sz)
	return Call_List{l}, err
}

func (s Call_List) At(i int) Call { return Call{s.List.Struct(i)} }

func (s Call_List) Set(i int, v Call) error { return s.List.SetStruct(i, v.Struct) }

// Call_Promise is a wrapper for a Call promised by a client call.
type Call_Promise struct{ *capnp.Pipeline }

func (p Call_Promise) Struct() (Call, error) {
	s, err := p.Pipeline.Struct()
	return Call{s}, err
}

func (p Call_Promise) Target() MessageTarget_Promise {
	return MessageTarget_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Call_Promise) Params() Payload_Promise {
	return Payload_Promise{Pipeline: p.Pipeline.GetPipeline(1)}
}

func (p Call_Promise) SendResultsTo() Call_sendResultsTo_Promise {
	return Call_sendResultsTo_Promise{p.Pipeline}
}

// Call_sendResultsTo_Promise is a wrapper for a Call_sendResultsTo promised by a client call.
type Call_sendResultsTo_Promise struct{ *capnp.Pipeline }

func (p Call_sendResultsTo_Promise) Struct() (Call_sendResultsTo, error) {
	s, err := p.Pipeline.Struct()
	return Call_sendResultsTo{s}, err
}

func (p Call_sendResultsTo_Promise) ThirdParty() *capnp.Pipeline {
	return p.Pipeline.GetPipeline(2)
}

type Return struct{ capnp.Struct }
type Return_Which uint16

const (
	Return_Which_results               Return_Which = 0
	Return_Which_exception             Return_Which = 1
	Return_Which_canceled              Return_Which = 2
	Return_Which_resultsSentElsewhere  Return_Which = 3
	Return_Which_takeFromOtherQuestion Return_Which = 4
	Return_Which_acceptFromThirdParty  Return_Which = 5
)

func (w Return_Which) String() string {
	const s = "resultsexceptioncanceledresultsSentElsewheretakeFromOtherQuestionacceptFromThirdParty"
	switch w {
	case Return_Which_results:
		return s[0:7]
	case Return_Which_exception:
		return s[7:16]
	case Return_Which_canceled:
		return s[16:24]
	case Return_Which_resultsSentElsewhere:
		return s[24:44]
	case Return_Which_takeFromOtherQuestion:
		return s[44:65]
	case Return_Which_acceptFromThirdParty:
		return s[65:85]

	}
	return "Return_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
}

func NewReturn(s *capnp.Segment) (Return, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 1})
	return Return{st}, err
}

func NewRootReturn(s *capnp.Segment) (Return, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 1})
	return Return{st}, err
}

func ReadRootReturn(msg *capnp.Message) (Return, error) {
	root, err := msg.RootPtr()
	return Return{root.Struct()}, err
}

func (s Return) String() string {
	str, _ := text.Marshal(0x9e19b28d3db3573a, s.Struct)
	return str
}

func (s Return) Which() Return_Which {
	return Return_Which(s.Struct.Uint16(6))
}
func (s Return) AnswerId() uint32 {
	return s.Struct.Uint32(0)
}

func (s Return) SetAnswerId(v uint32) {
	s.Struct.SetUint32(0, v)
}

func (s Return) ReleaseParamCaps() bool {
	return !s.Struct.Bit(32)
}

func (s Return) SetReleaseParamCaps(v bool) {
	s.Struct.SetBit(32, !v)
}

func (s Return) Results() (Payload, error) {
	p, err := s.Struct.Ptr(0)
	return Payload{Struct: p.Struct()}, err
}

func (s Return) HasResults() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Return) SetResults(v Payload) error {
	s.Struct.SetUint16(6, 0)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewResults sets the results field to a newly
// allocated Payload struct, preferring placement in s's segment.
func (s Return) NewResults() (Payload, error) {
	s.Struct.SetUint16(6, 0)
	ss, err := NewPayload(s.Struct.Segment())
	if err != nil {
		return Payload{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Return) Exception() (Exception, error) {
	p, err := s.Struct.Ptr(0)
	return Exception{Struct: p.Struct()}, err
}

func (s Return) HasException() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Return) SetException(v Exception) error {
	s.Struct.SetUint16(6, 1)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewException sets the exception field to a newly
// allocated Exception struct, preferring placement in s's segment.
func (s Return) NewException() (Exception, error) {
	s.Struct.SetUint16(6, 1)
	ss, err := NewException(s.Struct.Segment())
	if err != nil {
		return Exception{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Return) SetCanceled() {
	s.Struct.SetUint16(6, 2)

}

func (s Return) SetResultsSentElsewhere() {
	s.Struct.SetUint16(6, 3)

}

func (s Return) TakeFromOtherQuestion() uint32 {
	return s.Struct.Uint32(8)
}

func (s Return) SetTakeFromOtherQuestion(v uint32) {
	s.Struct.SetUint16(6, 4)
	s.Struct.SetUint32(8, v)
}

func (s Return) AcceptFromThirdParty() (capnp.Pointer, error) {
	return s.Struct.Pointer(0)
}

func (s Return) HasAcceptFromThirdParty() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Return) AcceptFromThirdPartyPtr() (capnp.Ptr, error) {
	return s.Struct.Ptr(0)
}

func (s Return) SetAcceptFromThirdParty(v capnp.Pointer) error {
	s.Struct.SetUint16(6, 5)
	return s.Struct.SetPointer(0, v)
}

func (s Return) SetAcceptFromThirdPartyPtr(v capnp.Ptr) error {
	s.Struct.SetUint16(6, 5)
	return s.Struct.SetPtr(0, v)
}

// Return_List is a list of Return.
type Return_List struct{ capnp.List }

// NewReturn creates a new list of Return.
func NewReturn_List(s *capnp.Segment, sz int32) (Return_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 16, PointerCount: 1}, sz)
	return Return_List{l}, err
}

func (s Return_List) At(i int) Return { return Return{s.List.Struct(i)} }

func (s Return_List) Set(i int, v Return) error { return s.List.SetStruct(i, v.Struct) }

// Return_Promise is a wrapper for a Return promised by a client call.
type Return_Promise struct{ *capnp.Pipeline }

func (p Return_Promise) Struct() (Return, error) {
	s, err := p.Pipeline.Struct()
	return Return{s}, err
}

func (p Return_Promise) Results() Payload_Promise {
	return Payload_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Return_Promise) Exception() Exception_Promise {
	return Exception_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Return_Promise) AcceptFromThirdParty() *capnp.Pipeline {
	return p.Pipeline.GetPipeline(0)
}

type Finish struct{ capnp.Struct }

func NewFinish(s *capnp.Segment) (Finish, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	return Finish{st}, err
}

func NewRootFinish(s *capnp.Segment) (Finish, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	return Finish{st}, err
}

func ReadRootFinish(msg *capnp.Message) (Finish, error) {
	root, err := msg.RootPtr()
	return Finish{root.Struct()}, err
}

func (s Finish) String() string {
	str, _ := text.Marshal(0xd37d2eb2c2f80e63, s.Struct)
	return str
}

func (s Finish) QuestionId() uint32 {
	return s.Struct.Uint32(0)
}

func (s Finish) SetQuestionId(v uint32) {
	s.Struct.SetUint32(0, v)
}

func (s Finish) ReleaseResultCaps() bool {
	return !s.Struct.Bit(32)
}

func (s Finish) SetReleaseResultCaps(v bool) {
	s.Struct.SetBit(32, !v)
}

// Finish_List is a list of Finish.
type Finish_List struct{ capnp.List }

// NewFinish creates a new list of Finish.
func NewFinish_List(s *capnp.Segment, sz int32) (Finish_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0}, sz)
	return Finish_List{l}, err
}

func (s Finish_List) At(i int) Finish { return Finish{s.List.Struct(i)} }

func (s Finish_List) Set(i int, v Finish) error { return s.List.SetStruct(i, v.Struct) }

// Finish_Promise is a wrapper for a Finish promised by a client call.
type Finish_Promise struct{ *capnp.Pipeline }

func (p Finish_Promise) Struct() (Finish, error) {
	s, err := p.Pipeline.Struct()
	return Finish{s}, err
}

type Resolve struct{ capnp.Struct }
type Resolve_Which uint16

const (
	Resolve_Which_cap       Resolve_Which = 0
	Resolve_Which_exception Resolve_Which = 1
)

func (w Resolve_Which) String() string {
	const s = "capexception"
	switch w {
	case Resolve_Which_cap:
		return s[0:3]
	case Resolve_Which_exception:
		return s[3:12]

	}
	return "Resolve_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
}

func NewResolve(s *capnp.Segment) (Resolve, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return Resolve{st}, err
}

func NewRootResolve(s *capnp.Segment) (Resolve, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return Resolve{st}, err
}

func ReadRootResolve(msg *capnp.Message) (Resolve, error) {
	root, err := msg.RootPtr()
	return Resolve{root.Struct()}, err
}

func (s Resolve) String() string {
	str, _ := text.Marshal(0xbbc29655fa89086e, s.Struct)
	return str
}

func (s Resolve) Which() Resolve_Which {
	return Resolve_Which(s.Struct.Uint16(4))
}
func (s Resolve) PromiseId() uint32 {
	return s.Struct.Uint32(0)
}

func (s Resolve) SetPromiseId(v uint32) {
	s.Struct.SetUint32(0, v)
}

func (s Resolve) Cap() (CapDescriptor, error) {
	p, err := s.Struct.Ptr(0)
	return CapDescriptor{Struct: p.Struct()}, err
}

func (s Resolve) HasCap() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Resolve) SetCap(v CapDescriptor) error {
	s.Struct.SetUint16(4, 0)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewCap sets the cap field to a newly
// allocated CapDescriptor struct, preferring placement in s's segment.
func (s Resolve) NewCap() (CapDescriptor, error) {
	s.Struct.SetUint16(4, 0)
	ss, err := NewCapDescriptor(s.Struct.Segment())
	if err != nil {
		return CapDescriptor{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Resolve) Exception() (Exception, error) {
	p, err := s.Struct.Ptr(0)
	return Exception{Struct: p.Struct()}, err
}

func (s Resolve) HasException() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Resolve) SetException(v Exception) error {
	s.Struct.SetUint16(4, 1)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewException sets the exception field to a newly
// allocated Exception struct, preferring placement in s's segment.
func (s Resolve) NewException() (Exception, error) {
	s.Struct.SetUint16(4, 1)
	ss, err := NewException(s.Struct.Segment())
	if err != nil {
		return Exception{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

// Resolve_List is a list of Resolve.
type Resolve_List struct{ capnp.List }

// NewResolve creates a new list of Resolve.
func NewResolve_List(s *capnp.Segment, sz int32) (Resolve_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1}, sz)
	return Resolve_List{l}, err
}

func (s Resolve_List) At(i int) Resolve { return Resolve{s.List.Struct(i)} }

func (s Resolve_List) Set(i int, v Resolve) error { return s.List.SetStruct(i, v.Struct) }

// Resolve_Promise is a wrapper for a Resolve promised by a client call.
type Resolve_Promise struct{ *capnp.Pipeline }

func (p Resolve_Promise) Struct() (Resolve, error) {
	s, err := p.Pipeline.Struct()
	return Resolve{s}, err
}

func (p Resolve_Promise) Cap() CapDescriptor_Promise {
	return CapDescriptor_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Resolve_Promise) Exception() Exception_Promise {
	return Exception_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

type Release struct{ capnp.Struct }

func NewRelease(s *capnp.Segment) (Release, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	return Release{st}, err
}

func NewRootRelease(s *capnp.Segment) (Release, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	return Release{st}, err
}

func ReadRootRelease(msg *capnp.Message) (Release, error) {
	root, err := msg.RootPtr()
	return Release{root.Struct()}, err
}

func (s Release) String() string {
	str, _ := text.Marshal(0xad1a6c0d7dd07497, s.Struct)
	return str
}

func (s Release) Id() uint32 {
	return s.Struct.Uint32(0)
}

func (s Release) SetId(v uint32) {
	s.Struct.SetUint32(0, v)
}

func (s Release) ReferenceCount() uint32 {
	return s.Struct.Uint32(4)
}

func (s Release) SetReferenceCount(v uint32) {
	s.Struct.SetUint32(4, v)
}

// Release_List is a list of Release.
type Release_List struct{ capnp.List }

// NewRelease creates a new list of Release.
func NewRelease_List(s *capnp.Segment, sz int32) (Release_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0}, sz)
	return Release_List{l}, err
}

func (s Release_List) At(i int) Release { return Release{s.List.Struct(i)} }

func (s Release_List) Set(i int, v Release) error { return s.List.SetStruct(i, v.Struct) }

// Release_Promise is a wrapper for a Release promised by a client call.
type Release_Promise struct{ *capnp.Pipeline }

func (p Release_Promise) Struct() (Release, error) {
	s, err := p.Pipeline.Struct()
	return Release{s}, err
}

type Disembargo struct{ capnp.Struct }
type Disembargo_context Disembargo
type Disembargo_context_Which uint16

const (
	Disembargo_context_Which_senderLoopback   Disembargo_context_Which = 0
	Disembargo_context_Which_receiverLoopback Disembargo_context_Which = 1
	Disembargo_context_Which_accept           Disembargo_context_Which = 2
	Disembargo_context_Which_provide          Disembargo_context_Which = 3
)

func (w Disembargo_context_Which) String() string {
	const s = "senderLoopbackreceiverLoopbackacceptprovide"
	switch w {
	case Disembargo_context_Which_senderLoopback:
		return s[0:14]
	case Disembargo_context_Which_receiverLoopback:
		return s[14:30]
	case Disembargo_context_Which_accept:
		return s[30:36]
	case Disembargo_context_Which_provide:
		return s[36:43]

	}
	return "Disembargo_context_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
}

func NewDisembargo(s *capnp.Segment) (Disembargo, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return Disembargo{st}, err
}

func NewRootDisembargo(s *capnp.Segment) (Disembargo, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return Disembargo{st}, err
}

func ReadRootDisembargo(msg *capnp.Message) (Disembargo, error) {
	root, err := msg.RootPtr()
	return Disembargo{root.Struct()}, err
}

func (s Disembargo) String() string {
	str, _ := text.Marshal(0xf964368b0fbd3711, s.Struct)
	return str
}

func (s Disembargo) Target() (MessageTarget, error) {
	p, err := s.Struct.Ptr(0)
	return MessageTarget{Struct: p.Struct()}, err
}

func (s Disembargo) HasTarget() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Disembargo) SetTarget(v MessageTarget) error {
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewTarget sets the target field to a newly
// allocated MessageTarget struct, preferring placement in s's segment.
func (s Disembargo) NewTarget() (MessageTarget, error) {
	ss, err := NewMessageTarget(s.Struct.Segment())
	if err != nil {
		return MessageTarget{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Disembargo) Context() Disembargo_context { return Disembargo_context(s) }

func (s Disembargo_context) Which() Disembargo_context_Which {
	return Disembargo_context_Which(s.Struct.Uint16(4))
}
func (s Disembargo_context) SenderLoopback() uint32 {
	return s.Struct.Uint32(0)
}

func (s Disembargo_context) SetSenderLoopback(v uint32) {
	s.Struct.SetUint16(4, 0)
	s.Struct.SetUint32(0, v)
}

func (s Disembargo_context) ReceiverLoopback() uint32 {
	return s.Struct.Uint32(0)
}

func (s Disembargo_context) SetReceiverLoopback(v uint32) {
	s.Struct.SetUint16(4, 1)
	s.Struct.SetUint32(0, v)
}

func (s Disembargo_context) SetAccept() {
	s.Struct.SetUint16(4, 2)

}

func (s Disembargo_context) Provide() uint32 {
	return s.Struct.Uint32(0)
}

func (s Disembargo_context) SetProvide(v uint32) {
	s.Struct.SetUint16(4, 3)
	s.Struct.SetUint32(0, v)
}

// Disembargo_List is a list of Disembargo.
type Disembargo_List struct{ capnp.List }

// NewDisembargo creates a new list of Disembargo.
func NewDisembargo_List(s *capnp.Segment, sz int32) (Disembargo_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1}, sz)
	return Disembargo_List{l}, err
}

func (s Disembargo_List) At(i int) Disembargo { return Disembargo{s.List.Struct(i)} }

func (s Disembargo_List) Set(i int, v Disembargo) error { return s.List.SetStruct(i, v.Struct) }

// Disembargo_Promise is a wrapper for a Disembargo promised by a client call.
type Disembargo_Promise struct{ *capnp.Pipeline }

func (p Disembargo_Promise) Struct() (Disembargo, error) {
	s, err := p.Pipeline.Struct()
	return Disembargo{s}, err
}

func (p Disembargo_Promise) Target() MessageTarget_Promise {
	return MessageTarget_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Disembargo_Promise) Context() Disembargo_context_Promise {
	return Disembargo_context_Promise{p.Pipeline}
}

// Disembargo_context_Promise is a wrapper for a Disembargo_context promised by a client call.
type Disembargo_context_Promise struct{ *capnp.Pipeline }

func (p Disembargo_context_Promise) Struct() (Disembargo_context, error) {
	s, err := p.Pipeline.Struct()
	return Disembargo_context{s}, err
}

type Provide struct{ capnp.Struct }

func NewProvide(s *capnp.Segment) (Provide, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2})
	return Provide{st}, err
}

func NewRootProvide(s *capnp.Segment) (Provide, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2})
	return Provide{st}, err
}

func ReadRootProvide(msg *capnp.Message) (Provide, error) {
	root, err := msg.RootPtr()
	return Provide{root.Struct()}, err
}

func (s Provide) String() string {
	str, _ := text.Marshal(0x9c6a046bfbc1ac5a, s.Struct)
	return str
}

func (s Provide) QuestionId() uint32 {
	return s.Struct.Uint32(0)
}

func (s Provide) SetQuestionId(v uint32) {
	s.Struct.SetUint32(0, v)
}

func (s Provide) Target() (MessageTarget, error) {
	p, err := s.Struct.Ptr(0)
	return MessageTarget{Struct: p.Struct()}, err
}

func (s Provide) HasTarget() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Provide) SetTarget(v MessageTarget) error {
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewTarget sets the target field to a newly
// allocated MessageTarget struct, preferring placement in s's segment.
func (s Provide) NewTarget() (MessageTarget, error) {
	ss, err := NewMessageTarget(s.Struct.Segment())
	if err != nil {
		return MessageTarget{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Provide) Recipient() (capnp.Pointer, error) {
	return s.Struct.Pointer(1)
}

func (s Provide) HasRecipient() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s Provide) RecipientPtr() (capnp.Ptr, error) {
	return s.Struct.Ptr(1)
}

func (s Provide) SetRecipient(v capnp.Pointer) error {
	return s.Struct.SetPointer(1, v)
}

func (s Provide) SetRecipientPtr(v capnp.Ptr) error {
	return s.Struct.SetPtr(1, v)
}

// Provide_List is a list of Provide.
type Provide_List struct{ capnp.List }

// NewProvide creates a new list of Provide.
func NewProvide_List(s *capnp.Segment, sz int32) (Provide_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2}, sz)
	return Provide_List{l}, err
}

func (s Provide_List) At(i int) Provide { return Provide{s.List.Struct(i)} }

func (s Provide_List) Set(i int, v Provide) error { return s.List.SetStruct(i, v.Struct) }

// Provide_Promise is a wrapper for a Provide promised by a client call.
type Provide_Promise struct{ *capnp.Pipeline }

func (p Provide_Promise) Struct() (Provide, error) {
	s, err := p.Pipeline.Struct()
	return Provide{s}, err
}

func (p Provide_Promise) Target() MessageTarget_Promise {
	return MessageTarget_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Provide_Promise) Recipient() *capnp.Pipeline {
	return p.Pipeline.GetPipeline(1)
}

type Accept struct{ capnp.Struct }

func NewAccept(s *capnp.Segment) (Accept, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return Accept{st}, err
}

func NewRootAccept(s *capnp.Segment) (Accept, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return Accept{st}, err
}

func ReadRootAccept(msg *capnp.Message) (Accept, error) {
	root, err := msg.RootPtr()
	return Accept{root.Struct()}, err
}

func (s Accept) String() string {
	str, _ := text.Marshal(0xd4c9b56290554016, s.Struct)
	return str
}

func (s Accept) QuestionId() uint32 {
	return s.Struct.Uint32(0)
}

func (s Accept) SetQuestionId(v uint32) {
	s.Struct.SetUint32(0, v)
}

func (s Accept) Provision() (capnp.Pointer, error) {
	return s.Struct.Pointer(0)
}

func (s Accept) HasProvision() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Accept) ProvisionPtr() (capnp.Ptr, error) {
	return s.Struct.Ptr(0)
}

func (s Accept) SetProvision(v capnp.Pointer) error {
	return s.Struct.SetPointer(0, v)
}

func (s Accept) SetProvisionPtr(v capnp.Ptr) error {
	return s.Struct.SetPtr(0, v)
}

func (s Accept) Embargo() bool {
	return s.Struct.Bit(32)
}

func (s Accept) SetEmbargo(v bool) {
	s.Struct.SetBit(32, v)
}

// Accept_List is a list of Accept.
type Accept_List struct{ capnp.List }

// NewAccept creates a new list of Accept.
func NewAccept_List(s *capnp.Segment, sz int32) (Accept_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1}, sz)
	return Accept_List{l}, err
}

func (s Accept_List) At(i int) Accept { return Accept{s.List.Struct(i)} }

func (s Accept_List) Set(i int, v Accept) error { return s.List.SetStruct(i, v.Struct) }

// Accept_Promise is a wrapper for a Accept promised by a client call.
type Accept_Promise struct{ *capnp.Pipeline }

func (p Accept_Promise) Struct() (Accept, error) {
	s, err := p.Pipeline.Struct()
	return Accept{s}, err
}

func (p Accept_Promise) Provision() *capnp.Pipeline {
	return p.Pipeline.GetPipeline(0)
}

type Join struct{ capnp.Struct }

func NewJoin(s *capnp.Segment) (Join, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2})
	return Join{st}, err
}

func NewRootJoin(s *capnp.Segment) (Join, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2})
	return Join{st}, err
}

func ReadRootJoin(msg *capnp.Message) (Join, error) {
	root, err := msg.RootPtr()
	return Join{root.Struct()}, err
}

func (s Join) String() string {
	str, _ := text.Marshal(0xfbe1980490e001af, s.Struct)
	return str
}

func (s Join) QuestionId() uint32 {
	return s.Struct.Uint32(0)
}

func (s Join) SetQuestionId(v uint32) {
	s.Struct.SetUint32(0, v)
}

func (s Join) Target() (MessageTarget, error) {
	p, err := s.Struct.Ptr(0)
	return MessageTarget{Struct: p.Struct()}, err
}

func (s Join) HasTarget() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Join) SetTarget(v MessageTarget) error {
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewTarget sets the target field to a newly
// allocated MessageTarget struct, preferring placement in s's segment.
func (s Join) NewTarget() (MessageTarget, error) {
	ss, err := NewMessageTarget(s.Struct.Segment())
	if err != nil {
		return MessageTarget{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Join) KeyPart() (capnp.Pointer, error) {
	return s.Struct.Pointer(1)
}

func (s Join) HasKeyPart() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s Join) KeyPartPtr() (capnp.Ptr, error) {
	return s.Struct.Ptr(1)
}

func (s Join) SetKeyPart(v capnp.Pointer) error {
	return s.Struct.SetPointer(1, v)
}

func (s Join) SetKeyPartPtr(v capnp.Ptr) error {
	return s.Struct.SetPtr(1, v)
}

// Join_List is a list of Join.
type Join_List struct{ capnp.List }

// NewJoin creates a new list of Join.
func NewJoin_List(s *capnp.Segment, sz int32) (Join_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2}, sz)
	return Join_List{l}, err
}

func (s Join_List) At(i int) Join { return Join{s.List.Struct(i)} }

func (s Join_List) Set(i int, v Join) error { return s.List.SetStruct(i, v.Struct) }

// Join_Promise is a wrapper for a Join promised by a client call.
type Join_Promise struct{ *capnp.Pipeline }

func (p Join_Promise) Struct() (Join, error) {
	s, err := p.Pipeline.Struct()
	return Join{s}, err
}

func (p Join_Promise) Target() MessageTarget_Promise {
	return MessageTarget_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Join_Promise) KeyPart() *capnp.Pipeline {
	return p.Pipeline.GetPipeline(1)
}

type MessageTarget struct{ capnp.Struct }
type MessageTarget_Which uint16

const (
	MessageTarget_Which_importedCap    MessageTarget_Which = 0
	MessageTarget_Which_promisedAnswer MessageTarget_Which = 1
)

func (w MessageTarget_Which) String() string {
	const s = "importedCappromisedAnswer"
	switch w {
	case MessageTarget_Which_importedCap:
		return s[0:11]
	case MessageTarget_Which_promisedAnswer:
		return s[11:25]

	}
	return "MessageTarget_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
}

func NewMessageTarget(s *capnp.Segment) (MessageTarget, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return MessageTarget{st}, err
}

func NewRootMessageTarget(s *capnp.Segment) (MessageTarget, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return MessageTarget{st}, err
}

func ReadRootMessageTarget(msg *capnp.Message) (MessageTarget, error) {
	root, err := msg.RootPtr()
	return MessageTarget{root.Struct()}, err
}

func (s MessageTarget) String() string {
	str, _ := text.Marshal(0x95bc14545813fbc1, s.Struct)
	return str
}

func (s MessageTarget) Which() MessageTarget_Which {
	return MessageTarget_Which(s.Struct.Uint16(4))
}
func (s MessageTarget) ImportedCap() uint32 {
	return s.Struct.Uint32(0)
}

func (s MessageTarget) SetImportedCap(v uint32) {
	s.Struct.SetUint16(4, 0)
	s.Struct.SetUint32(0, v)
}

func (s MessageTarget) PromisedAnswer() (PromisedAnswer, error) {
	p, err := s.Struct.Ptr(0)
	return PromisedAnswer{Struct: p.Struct()}, err
}

func (s MessageTarget) HasPromisedAnswer() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s MessageTarget) SetPromisedAnswer(v PromisedAnswer) error {
	s.Struct.SetUint16(4, 1)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewPromisedAnswer sets the promisedAnswer field to a newly
// allocated PromisedAnswer struct, preferring placement in s's segment.
func (s MessageTarget) NewPromisedAnswer() (PromisedAnswer, error) {
	s.Struct.SetUint16(4, 1)
	ss, err := NewPromisedAnswer(s.Struct.Segment())
	if err != nil {
		return PromisedAnswer{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

// MessageTarget_List is a list of MessageTarget.
type MessageTarget_List struct{ capnp.List }

// NewMessageTarget creates a new list of MessageTarget.
func NewMessageTarget_List(s *capnp.Segment, sz int32) (MessageTarget_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1}, sz)
	return MessageTarget_List{l}, err
}

func (s MessageTarget_List) At(i int) MessageTarget { return MessageTarget{s.List.Struct(i)} }

func (s MessageTarget_List) Set(i int, v MessageTarget) error { return s.List.SetStruct(i, v.Struct) }

// MessageTarget_Promise is a wrapper for a MessageTarget promised by a client call.
type MessageTarget_Promise struct{ *capnp.Pipeline }

func (p MessageTarget_Promise) Struct() (MessageTarget, error) {
	s, err := p.Pipeline.Struct()
	return MessageTarget{s}, err
}

func (p MessageTarget_Promise) PromisedAnswer() PromisedAnswer_Promise {
	return PromisedAnswer_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

type Payload struct{ capnp.Struct }

func NewPayload(s *capnp.Segment) (Payload, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	return Payload{st}, err
}

func NewRootPayload(s *capnp.Segment) (Payload, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	return Payload{st}, err
}

func ReadRootPayload(msg *capnp.Message) (Payload, error) {
	root, err := msg.RootPtr()
	return Payload{root.Struct()}, err
}

func (s Payload) String() string {
	str, _ := text.Marshal(0x9a0e61223d96743b, s.Struct)
	return str
}

func (s Payload) Content() (capnp.Pointer, error) {
	return s.Struct.Pointer(0)
}

func (s Payload) HasContent() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Payload) ContentPtr() (capnp.Ptr, error) {
	return s.Struct.Ptr(0)
}

func (s Payload) SetContent(v capnp.Pointer) error {
	return s.Struct.SetPointer(0, v)
}

func (s Payload) SetContentPtr(v capnp.Ptr) error {
	return s.Struct.SetPtr(0, v)
}

func (s Payload) CapTable() (CapDescriptor_List, error) {
	p, err := s.Struct.Ptr(1)
	return CapDescriptor_List{List: p.List()}, err
}

func (s Payload) HasCapTable() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s Payload) SetCapTable(v CapDescriptor_List) error {
	return s.Struct.SetPtr(1, v.List.ToPtr())
}

// NewCapTable sets the capTable field to a newly
// allocated CapDescriptor_List, preferring placement in s's segment.
func (s Payload) NewCapTable(n int32) (CapDescriptor_List, error) {
	l, err := NewCapDescriptor_List(s.Struct.Segment(), n)
	if err != nil {
		return CapDescriptor_List{}, err
	}
	err = s.Struct.SetPtr(1, l.List.ToPtr())
	return l, err
}

// Payload_List is a list of Payload.
type Payload_List struct{ capnp.List }

// NewPayload creates a new list of Payload.
func NewPayload_List(s *capnp.Segment, sz int32) (Payload_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2}, sz)
	return Payload_List{l}, err
}

func (s Payload_List) At(i int) Payload { return Payload{s.List.Struct(i)} }

func (s Payload_List) Set(i int, v Payload) error { return s.List.SetStruct(i, v.Struct) }

// Payload_Promise is a wrapper for a Payload promised by a client call.
type Payload_Promise struct{ *capnp.Pipeline }

func (p Payload_Promise) Struct() (Payload, error) {
	s, err := p.Pipeline.Struct()
	return Payload{s}, err
}

func (p Payload_Promise) Content() *capnp.Pipeline {
	return p.Pipeline.GetPipeline(0)
}

type CapDescriptor struct{ capnp.Struct }
type CapDescriptor_Which uint16

const (
	CapDescriptor_Which_none             CapDescriptor_Which = 0
	CapDescriptor_Which_senderHosted     CapDescriptor_Which = 1
	CapDescriptor_Which_senderPromise    CapDescriptor_Which = 2
	CapDescriptor_Which_receiverHosted   CapDescriptor_Which = 3
	CapDescriptor_Which_receiverAnswer   CapDescriptor_Which = 4
	CapDescriptor_Which_thirdPartyHosted CapDescriptor_Which = 5
)

func (w CapDescriptor_Which) String() string {
	const s = "nonesenderHostedsenderPromisereceiverHostedreceiverAnswerthirdPartyHosted"
	switch w {
	case CapDescriptor_Which_none:
		return s[0:4]
	case CapDescriptor_Which_senderHosted:
		return s[4:16]
	case CapDescriptor_Which_senderPromise:
		return s[16:29]
	case CapDescriptor_Which_receiverHosted:
		return s[29:43]
	case CapDescriptor_Which_receiverAnswer:
		return s[43:57]
	case CapDescriptor_Which_thirdPartyHosted:
		return s[57:73]

	}
	return "CapDescriptor_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
}

func NewCapDescriptor(s *capnp.Segment) (CapDescriptor, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return CapDescriptor{st}, err
}

func NewRootCapDescriptor(s *capnp.Segment) (CapDescriptor, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return CapDescriptor{st}, err
}

func ReadRootCapDescriptor(msg *capnp.Message) (CapDescriptor, error) {
	root, err := msg.RootPtr()
	return CapDescriptor{root.Struct()}, err
}

func (s CapDescriptor) String() string {
	str, _ := text.Marshal(0x8523ddc40b86b8b0, s.Struct)
	return str
}

func (s CapDescriptor) Which() CapDescriptor_Which {
	return CapDescriptor_Which(s.Struct.Uint16(0))
}
func (s CapDescriptor) SetNone() {
	s.Struct.SetUint16(0, 0)

}

func (s CapDescriptor) SenderHosted() uint32 {
	return s.Struct.Uint32(4)
}

func (s CapDescriptor) SetSenderHosted(v uint32) {
	s.Struct.SetUint16(0, 1)
	s.Struct.SetUint32(4, v)
}

func (s CapDescriptor) SenderPromise() uint32 {
	return s.Struct.Uint32(4)
}

func (s CapDescriptor) SetSenderPromise(v uint32) {
	s.Struct.SetUint16(0, 2)
	s.Struct.SetUint32(4, v)
}

func (s CapDescriptor) ReceiverHosted() uint32 {
	return s.Struct.Uint32(4)
}

func (s CapDescriptor) SetReceiverHosted(v uint32) {
	s.Struct.SetUint16(0, 3)
	s.Struct.SetUint32(4, v)
}

func (s CapDescriptor) ReceiverAnswer() (PromisedAnswer, error) {
	p, err := s.Struct.Ptr(0)
	return PromisedAnswer{Struct: p.Struct()}, err
}

func (s CapDescriptor) HasReceiverAnswer() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s CapDescriptor) SetReceiverAnswer(v PromisedAnswer) error {
	s.Struct.SetUint16(0, 4)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewReceiverAnswer sets the receiverAnswer field to a newly
// allocated PromisedAnswer struct, preferring placement in s's segment.
func (s CapDescriptor) NewReceiverAnswer() (PromisedAnswer, error) {
	s.Struct.SetUint16(0, 4)
	ss, err := NewPromisedAnswer(s.Struct.Segment())
	if err != nil {
		return PromisedAnswer{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s CapDescriptor) ThirdPartyHosted() (ThirdPartyCapDescriptor, error) {
	p, err := s.Struct.Ptr(0)
	return ThirdPartyCapDescriptor{Struct: p.Struct()}, err
}

func (s CapDescriptor) HasThirdPartyHosted() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s CapDescriptor) SetThirdPartyHosted(v ThirdPartyCapDescriptor) error {
	s.Struct.SetUint16(0, 5)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewThirdPartyHosted sets the thirdPartyHosted field to a newly
// allocated ThirdPartyCapDescriptor struct, preferring placement in s's segment.
func (s CapDescriptor) NewThirdPartyHosted() (ThirdPartyCapDescriptor, error) {
	s.Struct.SetUint16(0, 5)
	ss, err := NewThirdPartyCapDescriptor(s.Struct.Segment())
	if err != nil {
		return ThirdPartyCapDescriptor{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

// CapDescriptor_List is a list of CapDescriptor.
type CapDescriptor_List struct{ capnp.List }

// NewCapDescriptor creates a new list of CapDescriptor.
func NewCapDescriptor_List(s *capnp.Segment, sz int32) (CapDescriptor_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1}, sz)
	return CapDescriptor_List{l}, err
}

func (s CapDescriptor_List) At(i int) CapDescriptor { return CapDescriptor{s.List.Struct(i)} }

func (s CapDescriptor_List) Set(i int, v CapDescriptor) error { return s.List.SetStruct(i, v.Struct) }

// CapDescriptor_Promise is a wrapper for a CapDescriptor promised by a client call.
type CapDescriptor_Promise struct{ *capnp.Pipeline }

func (p CapDescriptor_Promise) Struct() (CapDescriptor, error) {
	s, err := p.Pipeline.Struct()
	return CapDescriptor{s}, err
}

func (p CapDescriptor_Promise) ReceiverAnswer() PromisedAnswer_Promise {
	return PromisedAnswer_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p CapDescriptor_Promise) ThirdPartyHosted() ThirdPartyCapDescriptor_Promise {
	return ThirdPartyCapDescriptor_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

type PromisedAnswer struct{ capnp.Struct }

func NewPromisedAnswer(s *capnp.Segment) (PromisedAnswer, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return PromisedAnswer{st}, err
}

func NewRootPromisedAnswer(s *capnp.Segment) (PromisedAnswer, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return PromisedAnswer{st}, err
}

func ReadRootPromisedAnswer(msg *capnp.Message) (PromisedAnswer, error) {
	root, err := msg.RootPtr()
	return PromisedAnswer{root.Struct()}, err
}

func (s PromisedAnswer) String() string {
	str, _ := text.Marshal(0xd800b1d6cd6f1ca0, s.Struct)
	return str
}

func (s PromisedAnswer) QuestionId() uint32 {
	return s.Struct.Uint32(0)
}

func (s PromisedAnswer) SetQuestionId(v uint32) {
	s.Struct.SetUint32(0, v)
}

func (s PromisedAnswer) Transform() (PromisedAnswer_Op_List, error) {
	p, err := s.Struct.Ptr(0)
	return PromisedAnswer_Op_List{List: p.List()}, err
}

func (s PromisedAnswer) HasTransform() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s PromisedAnswer) SetTransform(v PromisedAnswer_Op_List) error {
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewTransform sets the transform field to a newly
// allocated PromisedAnswer_Op_List, preferring placement in s's segment.
func (s PromisedAnswer) NewTransform(n int32) (PromisedAnswer_Op_List, error) {
	l, err := NewPromisedAnswer_Op_List(s.Struct.Segment(), n)
	if err != nil {
		return PromisedAnswer_Op_List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

// PromisedAnswer_List is a list of PromisedAnswer.
type PromisedAnswer_List struct{ capnp.List }

// NewPromisedAnswer creates a new list of PromisedAnswer.
func NewPromisedAnswer_List(s *capnp.Segment, sz int32) (PromisedAnswer_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1}, sz)
	return PromisedAnswer_List{l}, err
}

func (s PromisedAnswer_List) At(i int) PromisedAnswer { return PromisedAnswer{s.List.Struct(i)} }

func (s PromisedAnswer_List) Set(i int, v PromisedAnswer) error { return s.List.SetStruct(i, v.Struct) }

// PromisedAnswer_Promise is a wrapper for a PromisedAnswer promised by a client call.
type PromisedAnswer_Promise struct{ *capnp.Pipeline }

func (p PromisedAnswer_Promise) Struct() (PromisedAnswer, error) {
	s, err := p.Pipeline.Struct()
	return PromisedAnswer{s}, err
}

type PromisedAnswer_Op struct{ capnp.Struct }
type PromisedAnswer_Op_Which uint16

const (
	PromisedAnswer_Op_Which_noop            PromisedAnswer_Op_Which = 0
	PromisedAnswer_Op_Which_getPointerField PromisedAnswer_Op_Which = 1
)

func (w PromisedAnswer_Op_Which) String() string {
	const s = "noopgetPointerField"
	switch w {
	case PromisedAnswer_Op_Which_noop:
		return s[0:4]
	case PromisedAnswer_Op_Which_getPointerField:
		return s[4:19]

	}
	return "PromisedAnswer_Op_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
}

func NewPromisedAnswer_Op(s *capnp.Segment) (PromisedAnswer_Op, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	return PromisedAnswer_Op{st}, err
}

func NewRootPromisedAnswer_Op(s *capnp.Segment) (PromisedAnswer_Op, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	return PromisedAnswer_Op{st}, err
}

func ReadRootPromisedAnswer_Op(msg *capnp.Message) (PromisedAnswer_Op, error) {
	root, err := msg.RootPtr()
	return PromisedAnswer_Op{root.Struct()}, err
}

func (s PromisedAnswer_Op) String() string {
	str, _ := text.Marshal(0xf316944415569081, s.Struct)
	return str
}

func (s PromisedAnswer_Op) Which() PromisedAnswer_Op_Which {
	return PromisedAnswer_Op_Which(s.Struct.Uint16(0))
}
func (s PromisedAnswer_Op) SetNoop() {
	s.Struct.SetUint16(0, 0)

}

func (s PromisedAnswer_Op) GetPointerField() uint16 {
	return s.Struct.Uint16(2)
}

func (s PromisedAnswer_Op) SetGetPointerField(v uint16) {
	s.Struct.SetUint16(0, 1)
	s.Struct.SetUint16(2, v)
}

// PromisedAnswer_Op_List is a list of PromisedAnswer_Op.
type PromisedAnswer_Op_List struct{ capnp.List }

// NewPromisedAnswer_Op creates a new list of PromisedAnswer_Op.
func NewPromisedAnswer_Op_List(s *capnp.Segment, sz int32) (PromisedAnswer_Op_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0}, sz)
	return PromisedAnswer_Op_List{l}, err
}

func (s PromisedAnswer_Op_List) At(i int) PromisedAnswer_Op {
	return PromisedAnswer_Op{s.List.Struct(i)}
}

func (s PromisedAnswer_Op_List) Set(i int, v PromisedAnswer_Op) error {
	return s.List.SetStruct(i, v.Struct)
}

// PromisedAnswer_Op_Promise is a wrapper for a PromisedAnswer_Op promised by a client call.
type PromisedAnswer_Op_Promise struct{ *capnp.Pipeline }

func (p PromisedAnswer_Op_Promise) Struct() (PromisedAnswer_Op, error) {
	s, err := p.Pipeline.Struct()
	return PromisedAnswer_Op{s}, err
}

type ThirdPartyCapDescriptor struct{ capnp.Struct }

func NewThirdPartyCapDescriptor(s *capnp.Segment) (ThirdPartyCapDescriptor, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return ThirdPartyCapDescriptor{st}, err
}

func NewRootThirdPartyCapDescriptor(s *capnp.Segment) (ThirdPartyCapDescriptor, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return ThirdPartyCapDescriptor{st}, err
}

func ReadRootThirdPartyCapDescriptor(msg *capnp.Message) (ThirdPartyCapDescriptor, error) {
	root, err := msg.RootPtr()
	return ThirdPartyCapDescriptor{root.Struct()}, err
}

func (s ThirdPartyCapDescriptor) String() string {
	str, _ := text.Marshal(0xd37007fde1f0027d, s.Struct)
	return str
}

func (s ThirdPartyCapDescriptor) Id() (capnp.Pointer, error) {
	return s.Struct.Pointer(0)
}

func (s ThirdPartyCapDescriptor) HasId() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s ThirdPartyCapDescriptor) IdPtr() (capnp.Ptr, error) {
	return s.Struct.Ptr(0)
}

func (s ThirdPartyCapDescriptor) SetId(v capnp.Pointer) error {
	return s.Struct.SetPointer(0, v)
}

func (s ThirdPartyCapDescriptor) SetIdPtr(v capnp.Ptr) error {
	return s.Struct.SetPtr(0, v)
}

func (s ThirdPartyCapDescriptor) VineId() uint32 {
	return s.Struct.Uint32(0)
}

func (s ThirdPartyCapDescriptor) SetVineId(v uint32) {
	s.Struct.SetUint32(0, v)
}

// ThirdPartyCapDescriptor_List is a list of ThirdPartyCapDescriptor.
type ThirdPartyCapDescriptor_List struct{ capnp.List }

// NewThirdPartyCapDescriptor creates a new list of ThirdPartyCapDescriptor.
func NewThirdPartyCapDescriptor_List(s *capnp.Segment, sz int32) (ThirdPartyCapDescriptor_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1}, sz)
	return ThirdPartyCapDescriptor_List{l}, err
}

func (s ThirdPartyCapDescriptor_List) At(i int) ThirdPartyCapDescriptor {
	return ThirdPartyCapDescriptor{s.List.Struct(i)}
}

func (s ThirdPartyCapDescriptor_List) Set(i int, v ThirdPartyCapDescriptor) error {
	return s.List.SetStruct(i, v.Struct)
}

// ThirdPartyCapDescriptor_Promise is a wrapper for a ThirdPartyCapDescriptor promised by a client call.
type ThirdPartyCapDescriptor_Promise struct{ *capnp.Pipeline }

func (p ThirdPartyCapDescriptor_Promise) Struct() (ThirdPartyCapDescriptor, error) {
	s, err := p.Pipeline.Struct()
	return ThirdPartyCapDescriptor{s}, err
}

func (p ThirdPartyCapDescriptor_Promise) Id() *capnp.Pipeline {
	return p.Pipeline.GetPipeline(0)
}

type Exception struct{ capnp.Struct }

func NewException(s *capnp.Segment) (Exception, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return Exception{st}, err
}

func NewRootException(s *capnp.Segment) (Exception, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	return Exception{st}, err
}

func ReadRootException(msg *capnp.Message) (Exception, error) {
	root, err := msg.RootPtr()
	return Exception{root.Struct()}, err
}

func (s Exception) String() string {
	str, _ := text.Marshal(0xd625b7063acf691a, s.Struct)
	return str
}

func (s Exception) Reason() (string, error) {
	p, err := s.Struct.Ptr(0)
	return p.Text(), err
}

func (s Exception) HasReason() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Exception) ReasonBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	return p.TextBytes(), err
}

func (s Exception) SetReason(v string) error {
	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, t.List.ToPtr())
}

func (s Exception) Type() Exception_Type {
	return Exception_Type(s.Struct.Uint16(4))
}

func (s Exception) SetType(v Exception_Type) {
	s.Struct.SetUint16(4, uint16(v))
}

func (s Exception) ObsoleteIsCallersFault() bool {
	return s.Struct.Bit(0)
}

func (s Exception) SetObsoleteIsCallersFault(v bool) {
	s.Struct.SetBit(0, v)
}

func (s Exception) ObsoleteDurability() uint16 {
	return s.Struct.Uint16(2)
}

func (s Exception) SetObsoleteDurability(v uint16) {
	s.Struct.SetUint16(2, v)
}

// Exception_List is a list of Exception.
type Exception_List struct{ capnp.List }

// NewException creates a new list of Exception.
func NewException_List(s *capnp.Segment, sz int32) (Exception_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1}, sz)
	return Exception_List{l}, err
}

func (s Exception_List) At(i int) Exception { return Exception{s.List.Struct(i)} }

func (s Exception_List) Set(i int, v Exception) error { return s.List.SetStruct(i, v.Struct) }

// Exception_Promise is a wrapper for a Exception promised by a client call.
type Exception_Promise struct{ *capnp.Pipeline }

func (p Exception_Promise) Struct() (Exception, error) {
	s, err := p.Pipeline.Struct()
	return Exception{s}, err
}

type Exception_Type uint16

// Values of Exception_Type.
const (
	Exception_Type_failed        Exception_Type = 0
	Exception_Type_overloaded    Exception_Type = 1
	Exception_Type_disconnected  Exception_Type = 2
	Exception_Type_unimplemented Exception_Type = 3
)

// String returns the enum's constant name.
func (c Exception_Type) String() string {
	switch c {
	case Exception_Type_failed:
		return "failed"
	case Exception_Type_overloaded:
		return "overloaded"
	case Exception_Type_disconnected:
		return "disconnected"
	case Exception_Type_unimplemented:
		return "unimplemented"

	default:
		return ""
	}
}

// Exception_TypeFromString returns the enum value with a name,
// or the zero value if there's no such value.
func Exception_TypeFromString(c string) Exception_Type {
	switch c {
	case "failed":
		return Exception_Type_failed
	case "overloaded":
		return Exception_Type_overloaded
	case "disconnected":
		return Exception_Type_disconnected
	case "unimplemented":
		return Exception_Type_unimplemented

	default:
		return 0
	}
}

type Exception_Type_List struct{ capnp.List }

func NewException_Type_List(s *capnp.Segment, sz int32) (Exception_Type_List, error) {
	l, err := capnp.NewUInt16List(s, sz)
	return Exception_Type_List{l.List}, err
}

func (l Exception_Type_List) At(i int) Exception_Type {
	ul := capnp.UInt16List{List: l.List}
	return Exception_Type(ul.At(i))
}

func (l Exception_Type_List) Set(i int, v Exception_Type) {
	ul := capnp.UInt16List{List: l.List}
	ul.Set(i, uint16(v))
}

const schema_b312981b2552a250 = "x\xda\x9cW\x7f\x8c\x14W\x1d\x7fof\xf7\xf6\xb8\x1f" +
	"\xbb;7\xd7\"XrX%\x11\"\x84\xabU\xdbS" +
	"r@\xef\x08\xd7@\xb8\xbd;,b\x133\xbb\xfb\xe0" +
	"\x06\xf6f\x96\x999\xe0\x88\x17@[\xd3\"D \x80" +
	"\xd0\x80\x02\xa9I\xad4P\x0a)j\x89@0\xb6M" +
	"k!\x85\xa6\x98\x12K\x13#$\xfeQl\xad\xe5\xe7" +
	"\xf8\xf9\xce\xcc\xce\xec\xed\xed\x85\xd4?\x08s\x9f\xcf\xf7" +
	"\xbd\xf7}\xdf_\x9f\xb73\xaf\xc4gK\xad\xf1E\xe3" +
	"\x18\xcb\x14\xe25\xee\x85\x05\xfb\xd6\xfe\xb5w\xc5OY" +
	"\xa6\x8e\xcbn\xf7\xc1\x9e)_\xde\xdd\xf4\x0a\x8b\xcb\x09" +
	"\xc6\xd4-\xb1u\xea\xb6\x18\xbe\xbe\xb9%\xf6\x0b\xce\xb8" +
	"{\xe4\xc4\xcf\xea\xcf^\xfe\xea\xd3d\xcd#\xebN\x9e" +
	"\xa8\x81\xb9RsF\x9dPC\xe6\xf7\xd5<A\xe6\x0f" +
	"\x1d\xd9\xb2\xa1\xe5\xd7\xafn\x1bm\x9e\x84\xf9Pb\xbb" +
	"\xba1A\xe6\xc3\x89\xf12\xccO\xdfR\x97\xf45\xbf" +
	"\xb6s\xb4\xb9\xc4%\xf5d\xdd\x19\xf5\xcfu\xe4\xd6\xe9" +
	"\xba5\xb0\xfe\xae\xb3k\xd6\x83Z\xf29\xa6\xd4\x95\x19" +
	"\xc7%\xb2\x98R\xbf]\x9d^O_S\xeb\xc9v\xe9" +
	"\xa1\xd3\xb7V\xc6V\xec\xad\xd8\xd97\xde\x03\xe3\x03\x9e" +
	"\xf1\xbe\xfa\xc30n{\xe2\x95Y[\x8eN\xf8\x15\x19" +
	"K#/\xc9euV\xc3&\xb5\xb3\x81\xbc\x9e\xd3\xf0" +
	"\x17\xba\xe4/\x9ds\xc3\x8d\x85\x89/U\xecMaS" +
	"[\x93\xdb\xd5G\x93\xf4\xf5\xad$\xf9\xb1\xe4\xe4\x82\xf6" +
	"\x8fvm>\xca\x94f\xc9\x9d\xa8\xbf\xd3V\xf3\xea\x94" +
	"\xf7\x18\xe3\xea\xce\xe4\x9b\xea\x01\xcfp_r9\x0c\x8d" +
	"\xdago.\xdeu\xe6\x8f\xd5C\xf16\xb6\xbd\xe8Y" +
	"\x9fO\x92\xc7\xc3\xd2\xc7W\xee$\x8a\xefV^\x8f\x93" +
	"\x9b\xabRM\\\xdd\x98\"\xeb\xe1\x149\x91K~~" +
	"\xe6\xe8\x8c\xe1w\xab9|)\xb5I\xbd\xe2\xd9^\xf6" +
	"l\xef\x9f\xbdxk\xf6\xf8\x1b\x17\xaa\xed\xacv\xa67" +
	"\xa9\x0b\xd3\xf4\xd5\x95&7\xc2\x0bU3\xbe\x96>\xa8" +
	"^O\x8f\xc7\xd7\x8d\xf4?\x19\xff\xd3\xfe\x07\xcc\xb7\xdf" +
	"{\xf9\xfdj\xa6\x97\x957\xd5k\x0a\x99^W\xc8\x89" +
	"\xb3\xc6\xf8\xd6\x0d\xef,\xb8Vu_\xd1tP\x1dh" +
	"\xa2/\xbd\x89\x8c7n\xfd\xfe}\x1d;\xee\xff\x84e" +
	"&\xf0\xf0\x94\x8e\x84D\xc1j\xfaH\xbd\xec\x99^\xf2" +
	"L\x95\xef\x9cL\xfd\xfc\xdb\xf9\x1bU\xf7}T}Q" +
	"\x9d\xa3\xd2\xd7,\x95\x8c\x0f\xf3\x0f\xb7\xc6v_\xb9U" +
	"\xb5\x84\x86\xd5u\xeaF\xd5\xff:\xcc~\xe0Z\xc5\xdc" +
	"\x8c\x9cV4X{\xb1\xed1\xadP\xe8\xe6<\xf3\x80" +
	"\x1cc,\xc6\x19S\x8e/E\xf3\x1d\x93y\xe6\x94\xc4" +
	"9o\xe6\x84\x9dl\x03v\x02\xd8Y\x89+\x12@8" +
	"\xac\x9c\xce\x02<\x05\xf0-\x80\xb2\xd4\xcce\x80o<" +
	"\x0e\xf0u\x80\x17\x00\xc6a\x89m\x95\xf3\xb4\xfc-\x80" +
	"\xefc\xcb\x1a8\xb8\xe7G\xbf\x9b\xf8\xd9\x91\xab\x7fc" +
	"\xcaE\x8bIJlC3\xaf\xe5\\9}\x06vg" +
	"awN\xe2\xee\xaaAa;\xbai0\xb9+\xcfk" +
	"\x99\x84\x7f\xbc\xdd\xd1\xac\xe5\xc2\xe1\xe9\xa8\x1b\x19\xe7i" +
	"\x04@7\x1ca-\xd3r,!`>\x0e\xe6\xe3\x80" +
	"\x0e\x08\xa7\xdf\xccw\xe5\x19\xcc\x12\xc0\x12\xd8\xa2\xa8Y" +
	"\xda\x80\x8d-\xc2\x16\x0d\xb6\xb0\x85\x91\xef\x11\xf6 k" +
	")8v\x9f\xe9\"2\xe6\x9a\xbe~]\xb2\xf2\xdd\x9a" +
	"\xe5\x0c\xf5iz\x81\xc2\x05slE\xcdU\x0a\xa4D" +
	"q,v\x08;g\xe9E\xc7\xb4\x18E\xf4Kr\xac" +
	"\xc1u\xbd\x90\xee\x99\x86{\xed\xc0\xbd\xf6K|\x12\xbf" +
	"\xeb\x06Q\xdd\xb7\x02\xf0^\xc0/\x00\x96\xee\xb8A\\" +
	"\x7fc\x01~\x1e\xf0\x11\xc0\xf2m\x82)\xb2/\xad\x03" +
	"|\x08\xf0\x09\x897\xc6n\xb9~h\x8f\xaf\x8b\xb2\xd5" +
	"\x18\xbf\x094N\xf9\xda\x14\xa5&e\x98\x86`5\xde" +
	"\xf5\x845\xdfd)\xdb\x11aD\x03\xb8\xdbb-\xe6" +
	"\x80n\x8b\x10\xb7DN\xe8\xab\x85\xc5\xda\xe7\x9b#\x16" +
	"D\xc4\x1c\xc3^#,\x9e.\xd5q\x10G\xa7_\xf7" +
	"\"\xc6\x9d!\x7f)\x02\x9c\x8e\xa6@`U\x8a\x1d/" +
	"\xb6-\x14\xb6\xad-\xe7\x82\xa2\xf6H\x185u\x88#" +
	"\x10\xbdk\xb9\xcc{\x9f\xe2\xb8\x1d\x02\xe7\xc5M\xdd\xc8" +
	"\x1f\x02\xf1c\"\x9e!B\xbe\xe3z\x91S\x9f\xe6\x08" +
	"t\xef\x06\"6\x13\x11\xbb\xedz\xb1S\x9f\xe5\xa8@" +
	"\xec\x02b+\x11\xf1 |\xea\x16\x8fx\x86\x88\x1dD" +
	"\xd4\x04\x11T\xb7\xf1\xb9 6\x13\xb1\x9b\x88\xc4\x0d\x10" +
	"$#;=b+\x11{\x89\x18\xf79\x08oPs" +
	"\xa4\x13\xc6 \x9e'B\xfa/\x88Z\x10\x07x\x0f\x88" +
	"\xfdD\x1c\"\xa2\xee3\x10P8\xf5\xb7\x1c\xc9\xeb}" +
	"\x81\x88cD\xd4\xff\x07D\x1d\x88\x97\xbd3\x0e\x11q" +
	"\x82\x88\x86OA\xd4\x838\xee\xb9{\x84\x88\xd7\x88h" +
	"\xfc\x04D\x03\x88\xdf{7?F\xc4)\"j\xff\x0d" +
	"\xa2\x11\xc4I\x8ev\x861\x88\xd7A\xb8\x83\x86>P" +
	",\x88\x01\xd6\"\x0c\xcaj:\x92A?1-Z\xd6" +
	"\xb4\xa8\xc3\xca\x04\x80\xf0T\x0e\xa5\x0f8\x94d\x1fn" +
	"\xb7\x843h\x19 Ba\x0a\x88e\xba\xa1\xdb\xfd " +
	"\xc2\x89\xee\x13\xeb-a\x9b\x85\xd5\x02L\xa8#!S" +
	"\x10\x9aML([A\xb5\x98Y\xac\x11\x8e`\xa9^" +
	"\x0dK\x9bP\x8bM\x80\xb3\xa6\xe9\xd8\x8e\xa51^\xc4" +
	"\xa2p\x12W.j\xef\x10\xf4\x7fi\xd9\xfa\xa2e\xae" +
	"\xd6\xf3tN(\xbd\x81\xd3Z.'\x8at\xfbPZ" +
	"\x82\xdb\xaf0u\xbad8g\x83#\xf2h\x99\x81\xac" +
	"f1y\xb9\x09:\x9c\xd9\x15E.\x95\x8a\\\xf4y" +
	"\x03\xcc\x1b\x10\xb5\xd1\x80\x98J\xa3\xf4\xeb\xe8\xd7\x87\xcb" +
	"\xea\\i\xa5\xde\x9e\x09\xf4{H\x1c\xd2\x86\xbc\xa0\x99" +
	"\x12\x185a3\xe2&\xd4\xb5\xf91\x9b\xb1\xac\xcd\xba" +
	"\xb5\xa1\x82\xa9\xf1|pv0\xee\xa7\xa2\xd42_\xc3" +
	"!31\xb0K\xf3~:M\xf1o\x00\x9c/\xf1\xf5" +
	"9\x13\x95b8a\xd0\xb1]\x9f\x96-\x08\x1a\xaaI" +
	"\xc6\xbbe\x9c\x14\xbd\xbdpn\xb2\xe2\\/\xda~{" +
	"7\x84\xe7v\x92\xcct\xe0\x88\xeeHf\x16\x92N\xcc" +
	"\x07\xd6W&3\x19tO\x06\x87d\x9e\xfc\xc2\xa2\x80" +
	"Q\xa5\x17ua0\x1ey_\xe6X\x8fW\xba\xccK" +
	"\xc6\xe4\xd0\xb1\xf3t\xf7s8\xef\x03\x0a\xc8dx\x06" +
	"e\xbaD\x03\xf5\x03\x80W\xa9\xb3]\x7f\xde(\xff\xa0" +
	"\xd8}\x08\xf4_4\x85\xee\xfa\xc3F\xb9F\x0e_\x05" +
	"\xfa)\x8d\xa0;\xc1\xa0\xbeN\xdb~\x0c\xf46\xcd\x9f" +
	"\xdb\xc1\xa0\xbe\xf1\"\xd0\xdbh\xceZ4\xe7\xa4\x9a[" +
	"\xae\xe4O\x998?\x8a\xb6\xad\xa5\xb6m\xf6\xc6\xcf\xcd" +
	"`\xca(\x1c+\x80\x81\x98L\xfd\xacyi\xf7\x15." +
	"\x9a\xd0^\x1busR:T\x8b\xcd<\xc9\x8a\x93f" +
	"Q\xf7\x0dB\xdf\xaa\xe8\x9fXK\xb5\xaf\x9b\x8c\x1b\xa3" +
	"\xdb\x1fY7r\xd8\x17\x07%\xdc`\x8f^\x8e\xb2\xe8" +
	",\xd8bM\xaa_X\xa41\x8e\xb6R\xccCI\xf2" +
	"E\x0e\x90\xcc\xa0h\xf1\xb2\x15z\xe6\xb7\xd7<\x8b\x9b" +
	"\x03}\x9eJ\xa4HX\xab\xe7\x86\xee\xe0\x17MY\xb1" +
	"N\xacV\xac\xeb\x82b}D\xe2\xb2^.T\xcb\xe0" +
	"\x95\x91c\xed\xe21s\x10\x05\x1c\x12QWv\x06w" +
	"6f\xf4\x0d\x15\x85_\x0ai/\xb7S\xdb\xe8\xe6\xca" +
	"WP\xa6\\R&a\xb6sY\x99\x00Ej_\x86" +
	"w\x80\xc8\xbb&4\x10\xfd\x94g2\xfe\xc0\x1c@\x9b" +
	"@gS9\xf4h\xe5\x94\x1dy1\x9a~\xa3\xba\xa1" +
	"'\xea\x86F\xee\x06\x03`\xe1\x83Q?4Jw\xdd" +
	"*\x0d\x11\x0c\x80.\xc6\xc3\x8b'pREG\xde;" +
	"\xbd%\x0f\xe5b[_\xa0\xdf\xceP\xf9\xa3\x86[c" +
	"\xa7\"\xccD[4\xc6(\x13A^\xdbW\xeb\x86\xe8" +
	"\xca\x8f\x8a?\xa21\xcf\x13\x09\xc6*\xf6^\x1a\xed\x13" +
	"\xb6`\xebv\x80\x0f\x03\x9c=\xc6\x1c(\xd5}\x0f\xf7" +
	"\xca\x93\xc6\xa4\x1d\xd6}\xf9\xa1s\xbc*\xf4\x0f\xbd\xc7" +
	"@\xa2P/\x00\xb6\x84\x06\xd2d?\xfe\x8b\xe7\xdec" +
	" \xb9\x9e\xbe\xd8~\xa8K\x9a\xe3\xc9\x04D\xa2\xf2\xed" +
	"\xc8KE\x98\xc0.\x99\x18/\xff\x19\xc6\xa7\xa5\xa8," +
	"3\xe9\xd0I\x8d\"\xfc$\xce\xee\xc7.\x92\xef\xa4\xf8" +
	"\x03\xb0~`\x0e\xbd\xc3\x83\xa9\xb9\x0am\x9dq\x00n" +
	"\xa0\x10\x06\x8f\xf3az\x87\xae\x05\xf8\x94D\xc2\xad\xd9" +
	"h\xcd\x06\xf8\xd3P&\x96\xbc\xcb\xa6G\xae\xb0\xda\xed" +
	"y\x1a\xc2\x18:\x1c\x1at\x0cZZV/\xe82:" +
	"7xT\xa7\x1c\xb8\xc9S\x91\xeb\xa8\xaa\xd4\xc8N\xeb" +
	"\x0e\x94\xca\xd7)\xf8AW\x0d\x7f\x0e)|\xa2\xbc\xa8" +
	"8F\x09\x94\xb2\xd1\xda\x13\xe8\xe1\x82\xb1\x02\x8f\xa7\x80" +
	"a/\xc3\x03\x9c\x0fD\xd2\x14\x1e2Z\x9a\xe6z\xef" +
	"\x87\x84\xa5\x15\xc7.\xc0\xf0\xf4\xe7\xeeU\x7fyQ\x84" +
	"\xe2h\x0e\x17\xf9E\xd9\x15\"\xe7\x10Y9\xdbF\x85" +
	"\"1cQ\xb1\xf290-\xea\xad\xb2\xdf\x0b\xd3\x7f" +
	"\x12\x0d:\xbc\xea\xcd\"&.\x84\xaf\x1bO\x13\x9ci" +
	"\xcd\xd3E!\x1f\xfe\xce)\xbff\x87\xf7JIQ\xfd" +
	"U\xdc\xb3\xad\xbc\x89\x91\x90\x85\x97\x7f(\xfe~,{" +
	"\x11\x87\xcde\xd2\x98\xca\xea\xbf\x09\xd6:#~J>" +
	"\x0e?\xfe_\x8d\x9f\x1b\xf5\xd9\x17\xd3\xf8\xf5+\xc5\x10" +
	"\xcd\xaaR\x9c\xff\x17\x00\x00\xff\xff\x03\xe1\x99\x85"

func init() {
	schemas.Register(schema_b312981b2552a250,
		0x836a53ce789d4cd4,
		0x8523ddc40b86b8b0,
		0x91b79f1f808db032,
		0x95bc14545813fbc1,
		0x9a0e61223d96743b,
		0x9c6a046bfbc1ac5a,
		0x9e19b28d3db3573a,
		0xad1a6c0d7dd07497,
		0xb28c96e23f4cbd58,
		0xbbc29655fa89086e,
		0xd37007fde1f0027d,
		0xd37d2eb2c2f80e63,
		0xd4c9b56290554016,
		0xd625b7063acf691a,
		0xd800b1d6cd6f1ca0,
		0xe94ccf8031176ec4,
		0xf316944415569081,
		0xf964368b0fbd3711,
		0xfbe1980490e001af)
}
