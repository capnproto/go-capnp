package rpc

// AUTO GENERATED - DO NOT EDIT

import (
	strconv "strconv"
	capnp "zombiezen.com/go/capnproto2"
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
	if err != nil {
		return Message{}, err
	}
	return Message{st}, nil
}

func NewRootMessage(s *capnp.Segment) (Message, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		return Message{}, err
	}
	return Message{st}, nil
}

func ReadRootMessage(msg *capnp.Message) (Message, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Message{}, err
	}
	return Message{root.Struct()}, nil
}

func (s Message) Which() Message_Which {
	return Message_Which(s.Struct.Uint16(0))
}
func (s Message) Unimplemented() (Message, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return Message{}, err
	}
	return Message{Struct: p.Struct()}, nil
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
	if err != nil {
		return Exception{}, err
	}
	return Exception{Struct: p.Struct()}, nil
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
	if err != nil {
		return Bootstrap{}, err
	}
	return Bootstrap{Struct: p.Struct()}, nil
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
	if err != nil {
		return Call{}, err
	}
	return Call{Struct: p.Struct()}, nil
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
	if err != nil {
		return Return{}, err
	}
	return Return{Struct: p.Struct()}, nil
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
	if err != nil {
		return Finish{}, err
	}
	return Finish{Struct: p.Struct()}, nil
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
	if err != nil {
		return Resolve{}, err
	}
	return Resolve{Struct: p.Struct()}, nil
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
	if err != nil {
		return Release{}, err
	}
	return Release{Struct: p.Struct()}, nil
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
	if err != nil {
		return Disembargo{}, err
	}
	return Disembargo{Struct: p.Struct()}, nil
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
	if err != nil {
		return Provide{}, err
	}
	return Provide{Struct: p.Struct()}, nil
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
	if err != nil {
		return Accept{}, err
	}
	return Accept{Struct: p.Struct()}, nil
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
	if err != nil {
		return Join{}, err
	}
	return Join{Struct: p.Struct()}, nil
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
	if err != nil {
		return Message_List{}, err
	}
	return Message_List{l}, nil
}

func (s Message_List) At(i int) Message           { return Message{s.List.Struct(i)} }
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
	if err != nil {
		return Bootstrap{}, err
	}
	return Bootstrap{st}, nil
}

func NewRootBootstrap(s *capnp.Segment) (Bootstrap, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		return Bootstrap{}, err
	}
	return Bootstrap{st}, nil
}

func ReadRootBootstrap(msg *capnp.Message) (Bootstrap, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Bootstrap{}, err
	}
	return Bootstrap{root.Struct()}, nil
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
	if err != nil {
		return Bootstrap_List{}, err
	}
	return Bootstrap_List{l}, nil
}

func (s Bootstrap_List) At(i int) Bootstrap           { return Bootstrap{s.List.Struct(i)} }
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
	if err != nil {
		return Call{}, err
	}
	return Call{st}, nil
}

func NewRootCall(s *capnp.Segment) (Call, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 24, PointerCount: 3})
	if err != nil {
		return Call{}, err
	}
	return Call{st}, nil
}

func ReadRootCall(msg *capnp.Message) (Call, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Call{}, err
	}
	return Call{root.Struct()}, nil
}
func (s Call) QuestionId() uint32 {
	return s.Struct.Uint32(0)
}

func (s Call) SetQuestionId(v uint32) {
	s.Struct.SetUint32(0, v)
}

func (s Call) Target() (MessageTarget, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return MessageTarget{}, err
	}
	return MessageTarget{Struct: p.Struct()}, nil
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
	if err != nil {
		return Payload{}, err
	}
	return Payload{Struct: p.Struct()}, nil
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
	if err != nil {
		return Call_List{}, err
	}
	return Call_List{l}, nil
}

func (s Call_List) At(i int) Call           { return Call{s.List.Struct(i)} }
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
	if err != nil {
		return Return{}, err
	}
	return Return{st}, nil
}

func NewRootReturn(s *capnp.Segment) (Return, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 1})
	if err != nil {
		return Return{}, err
	}
	return Return{st}, nil
}

func ReadRootReturn(msg *capnp.Message) (Return, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Return{}, err
	}
	return Return{root.Struct()}, nil
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
	if err != nil {
		return Payload{}, err
	}
	return Payload{Struct: p.Struct()}, nil
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
	if err != nil {
		return Exception{}, err
	}
	return Exception{Struct: p.Struct()}, nil
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
	if err != nil {
		return Return_List{}, err
	}
	return Return_List{l}, nil
}

func (s Return_List) At(i int) Return           { return Return{s.List.Struct(i)} }
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
	if err != nil {
		return Finish{}, err
	}
	return Finish{st}, nil
}

func NewRootFinish(s *capnp.Segment) (Finish, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	if err != nil {
		return Finish{}, err
	}
	return Finish{st}, nil
}

func ReadRootFinish(msg *capnp.Message) (Finish, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Finish{}, err
	}
	return Finish{root.Struct()}, nil
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
	if err != nil {
		return Finish_List{}, err
	}
	return Finish_List{l}, nil
}

func (s Finish_List) At(i int) Finish           { return Finish{s.List.Struct(i)} }
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
	if err != nil {
		return Resolve{}, err
	}
	return Resolve{st}, nil
}

func NewRootResolve(s *capnp.Segment) (Resolve, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		return Resolve{}, err
	}
	return Resolve{st}, nil
}

func ReadRootResolve(msg *capnp.Message) (Resolve, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Resolve{}, err
	}
	return Resolve{root.Struct()}, nil
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
	if err != nil {
		return CapDescriptor{}, err
	}
	return CapDescriptor{Struct: p.Struct()}, nil
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
	if err != nil {
		return Exception{}, err
	}
	return Exception{Struct: p.Struct()}, nil
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
	if err != nil {
		return Resolve_List{}, err
	}
	return Resolve_List{l}, nil
}

func (s Resolve_List) At(i int) Resolve           { return Resolve{s.List.Struct(i)} }
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
	if err != nil {
		return Release{}, err
	}
	return Release{st}, nil
}

func NewRootRelease(s *capnp.Segment) (Release, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	if err != nil {
		return Release{}, err
	}
	return Release{st}, nil
}

func ReadRootRelease(msg *capnp.Message) (Release, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Release{}, err
	}
	return Release{root.Struct()}, nil
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
	if err != nil {
		return Release_List{}, err
	}
	return Release_List{l}, nil
}

func (s Release_List) At(i int) Release           { return Release{s.List.Struct(i)} }
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
	if err != nil {
		return Disembargo{}, err
	}
	return Disembargo{st}, nil
}

func NewRootDisembargo(s *capnp.Segment) (Disembargo, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		return Disembargo{}, err
	}
	return Disembargo{st}, nil
}

func ReadRootDisembargo(msg *capnp.Message) (Disembargo, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Disembargo{}, err
	}
	return Disembargo{root.Struct()}, nil
}
func (s Disembargo) Target() (MessageTarget, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return MessageTarget{}, err
	}
	return MessageTarget{Struct: p.Struct()}, nil
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
	if err != nil {
		return Disembargo_List{}, err
	}
	return Disembargo_List{l}, nil
}

func (s Disembargo_List) At(i int) Disembargo           { return Disembargo{s.List.Struct(i)} }
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
	if err != nil {
		return Provide{}, err
	}
	return Provide{st}, nil
}

func NewRootProvide(s *capnp.Segment) (Provide, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2})
	if err != nil {
		return Provide{}, err
	}
	return Provide{st}, nil
}

func ReadRootProvide(msg *capnp.Message) (Provide, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Provide{}, err
	}
	return Provide{root.Struct()}, nil
}
func (s Provide) QuestionId() uint32 {
	return s.Struct.Uint32(0)
}

func (s Provide) SetQuestionId(v uint32) {
	s.Struct.SetUint32(0, v)
}

func (s Provide) Target() (MessageTarget, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return MessageTarget{}, err
	}
	return MessageTarget{Struct: p.Struct()}, nil
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
	if err != nil {
		return Provide_List{}, err
	}
	return Provide_List{l}, nil
}

func (s Provide_List) At(i int) Provide           { return Provide{s.List.Struct(i)} }
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
	if err != nil {
		return Accept{}, err
	}
	return Accept{st}, nil
}

func NewRootAccept(s *capnp.Segment) (Accept, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		return Accept{}, err
	}
	return Accept{st}, nil
}

func ReadRootAccept(msg *capnp.Message) (Accept, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Accept{}, err
	}
	return Accept{root.Struct()}, nil
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
	if err != nil {
		return Accept_List{}, err
	}
	return Accept_List{l}, nil
}

func (s Accept_List) At(i int) Accept           { return Accept{s.List.Struct(i)} }
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
	if err != nil {
		return Join{}, err
	}
	return Join{st}, nil
}

func NewRootJoin(s *capnp.Segment) (Join, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2})
	if err != nil {
		return Join{}, err
	}
	return Join{st}, nil
}

func ReadRootJoin(msg *capnp.Message) (Join, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Join{}, err
	}
	return Join{root.Struct()}, nil
}
func (s Join) QuestionId() uint32 {
	return s.Struct.Uint32(0)
}

func (s Join) SetQuestionId(v uint32) {
	s.Struct.SetUint32(0, v)
}

func (s Join) Target() (MessageTarget, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return MessageTarget{}, err
	}
	return MessageTarget{Struct: p.Struct()}, nil
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
	if err != nil {
		return Join_List{}, err
	}
	return Join_List{l}, nil
}

func (s Join_List) At(i int) Join           { return Join{s.List.Struct(i)} }
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
	if err != nil {
		return MessageTarget{}, err
	}
	return MessageTarget{st}, nil
}

func NewRootMessageTarget(s *capnp.Segment) (MessageTarget, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		return MessageTarget{}, err
	}
	return MessageTarget{st}, nil
}

func ReadRootMessageTarget(msg *capnp.Message) (MessageTarget, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return MessageTarget{}, err
	}
	return MessageTarget{root.Struct()}, nil
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
	if err != nil {
		return PromisedAnswer{}, err
	}
	return PromisedAnswer{Struct: p.Struct()}, nil
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
	if err != nil {
		return MessageTarget_List{}, err
	}
	return MessageTarget_List{l}, nil
}

func (s MessageTarget_List) At(i int) MessageTarget           { return MessageTarget{s.List.Struct(i)} }
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
	if err != nil {
		return Payload{}, err
	}
	return Payload{st}, nil
}

func NewRootPayload(s *capnp.Segment) (Payload, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	if err != nil {
		return Payload{}, err
	}
	return Payload{st}, nil
}

func ReadRootPayload(msg *capnp.Message) (Payload, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Payload{}, err
	}
	return Payload{root.Struct()}, nil
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
	if err != nil {
		return CapDescriptor_List{}, err
	}
	return CapDescriptor_List{List: p.List()}, nil
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
	if err != nil {
		return Payload_List{}, err
	}
	return Payload_List{l}, nil
}

func (s Payload_List) At(i int) Payload           { return Payload{s.List.Struct(i)} }
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
	if err != nil {
		return CapDescriptor{}, err
	}
	return CapDescriptor{st}, nil
}

func NewRootCapDescriptor(s *capnp.Segment) (CapDescriptor, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		return CapDescriptor{}, err
	}
	return CapDescriptor{st}, nil
}

func ReadRootCapDescriptor(msg *capnp.Message) (CapDescriptor, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return CapDescriptor{}, err
	}
	return CapDescriptor{root.Struct()}, nil
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
	if err != nil {
		return PromisedAnswer{}, err
	}
	return PromisedAnswer{Struct: p.Struct()}, nil
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
	if err != nil {
		return ThirdPartyCapDescriptor{}, err
	}
	return ThirdPartyCapDescriptor{Struct: p.Struct()}, nil
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
	if err != nil {
		return CapDescriptor_List{}, err
	}
	return CapDescriptor_List{l}, nil
}

func (s CapDescriptor_List) At(i int) CapDescriptor           { return CapDescriptor{s.List.Struct(i)} }
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
	if err != nil {
		return PromisedAnswer{}, err
	}
	return PromisedAnswer{st}, nil
}

func NewRootPromisedAnswer(s *capnp.Segment) (PromisedAnswer, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		return PromisedAnswer{}, err
	}
	return PromisedAnswer{st}, nil
}

func ReadRootPromisedAnswer(msg *capnp.Message) (PromisedAnswer, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return PromisedAnswer{}, err
	}
	return PromisedAnswer{root.Struct()}, nil
}
func (s PromisedAnswer) QuestionId() uint32 {
	return s.Struct.Uint32(0)
}

func (s PromisedAnswer) SetQuestionId(v uint32) {
	s.Struct.SetUint32(0, v)
}

func (s PromisedAnswer) Transform() (PromisedAnswer_Op_List, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return PromisedAnswer_Op_List{}, err
	}
	return PromisedAnswer_Op_List{List: p.List()}, nil
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
	if err != nil {
		return PromisedAnswer_List{}, err
	}
	return PromisedAnswer_List{l}, nil
}

func (s PromisedAnswer_List) At(i int) PromisedAnswer           { return PromisedAnswer{s.List.Struct(i)} }
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
	if err != nil {
		return PromisedAnswer_Op{}, err
	}
	return PromisedAnswer_Op{st}, nil
}

func NewRootPromisedAnswer_Op(s *capnp.Segment) (PromisedAnswer_Op, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	if err != nil {
		return PromisedAnswer_Op{}, err
	}
	return PromisedAnswer_Op{st}, nil
}

func ReadRootPromisedAnswer_Op(msg *capnp.Message) (PromisedAnswer_Op, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return PromisedAnswer_Op{}, err
	}
	return PromisedAnswer_Op{root.Struct()}, nil
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
	if err != nil {
		return PromisedAnswer_Op_List{}, err
	}
	return PromisedAnswer_Op_List{l}, nil
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
	if err != nil {
		return ThirdPartyCapDescriptor{}, err
	}
	return ThirdPartyCapDescriptor{st}, nil
}

func NewRootThirdPartyCapDescriptor(s *capnp.Segment) (ThirdPartyCapDescriptor, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		return ThirdPartyCapDescriptor{}, err
	}
	return ThirdPartyCapDescriptor{st}, nil
}

func ReadRootThirdPartyCapDescriptor(msg *capnp.Message) (ThirdPartyCapDescriptor, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return ThirdPartyCapDescriptor{}, err
	}
	return ThirdPartyCapDescriptor{root.Struct()}, nil
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
	if err != nil {
		return ThirdPartyCapDescriptor_List{}, err
	}
	return ThirdPartyCapDescriptor_List{l}, nil
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
	if err != nil {
		return Exception{}, err
	}
	return Exception{st}, nil
}

func NewRootException(s *capnp.Segment) (Exception, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		return Exception{}, err
	}
	return Exception{st}, nil
}

func ReadRootException(msg *capnp.Message) (Exception, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Exception{}, err
	}
	return Exception{root.Struct()}, nil
}
func (s Exception) Reason() (string, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return "", err
	}
	return p.Text(), nil
}

func (s Exception) HasReason() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Exception) ReasonBytes() ([]byte, error) {
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
	if err != nil {
		return Exception_List{}, err
	}
	return Exception_List{l}, nil
}

func (s Exception_List) At(i int) Exception           { return Exception{s.List.Struct(i)} }
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
	if err != nil {
		return Exception_Type_List{}, err
	}
	return Exception_Type_List{l.List}, nil
}

func (l Exception_Type_List) At(i int) Exception_Type {
	ul := capnp.UInt16List{List: l.List}
	return Exception_Type(ul.At(i))
}

func (l Exception_Type_List) Set(i int, v Exception_Type) {
	ul := capnp.UInt16List{List: l.List}
	ul.Set(i, uint16(v))
}

var schema_b312981b2552a250 = []byte{
	120, 218, 156, 88, 123, 112, 84, 229,
	21, 255, 190, 187, 155, 108, 18, 195,
	110, 110, 110, 66, 10, 133, 9, 162,
	180, 138, 188, 130, 210, 234, 182, 204,
	66, 12, 12, 48, 80, 114, 67, 168,
	66, 165, 122, 179, 123, 9, 11, 187,
	123, 47, 119, 111, 32, 75, 101, 2,
	86, 29, 165, 50, 5, 20, 11, 140,
	248, 26, 157, 170, 197, 17, 17, 42,
	182, 82, 133, 193, 250, 24, 241, 81,
	209, 209, 142, 78, 193, 25, 167, 218,
	177, 83, 223, 149, 231, 237, 239, 220,
	119, 146, 101, 24, 251, 71, 96, 247,
	252, 206, 189, 223, 249, 206, 227, 119,
	206, 217, 73, 195, 99, 211, 132, 150,
	138, 249, 213, 140, 201, 185, 138, 74,
	107, 95, 91, 252, 7, 252, 143, 147,
	142, 51, 185, 186, 130, 91, 183, 229,
	62, 56, 43, 143, 24, 251, 38, 99,
	92, 218, 24, 189, 73, 218, 28, 141,
	49, 182, 224, 142, 104, 132, 51, 110,
	77, 222, 189, 113, 93, 243, 125, 79,
	111, 102, 114, 13, 231, 86, 251, 131,
	29, 99, 190, 191, 173, 254, 41, 54,
	131, 199, 226, 140, 73, 107, 163, 91,
	164, 91, 72, 253, 242, 245, 209, 166,
	8, 227, 207, 221, 63, 66, 59, 242,
	206, 147, 239, 246, 87, 174, 224, 208,
	144, 14, 198, 94, 145, 142, 196, 154,
	240, 233, 237, 216, 106, 188, 121, 253,
	207, 38, 237, 255, 213, 218, 111, 15,
	48, 94, 195, 4, 105, 118, 213, 77,
	210, 188, 170, 31, 74, 217, 170, 249,
	76, 176, 14, 23, 154, 90, 214, 189,
	62, 247, 19, 86, 238, 69, 47, 84,
	61, 40, 29, 169, 162, 79, 47, 87,
	209, 139, 118, 239, 191, 245, 130, 195,
	31, 92, 116, 203, 96, 19, 43, 161,
	211, 82, 125, 72, 186, 170, 154, 76,
	156, 82, 125, 13, 221, 40, 121, 205,
	83, 83, 55, 238, 25, 118, 47, 169,
	11, 253, 213, 121, 68, 218, 94, 179,
	65, 122, 160, 134, 212, 119, 214, 188,
	72, 234, 71, 231, 238, 236, 125, 109,
	193, 242, 95, 147, 122, 36, 100, 74,
	132, 12, 216, 88, 187, 70, 218, 92,
	75, 218, 27, 107, 127, 75, 218, 190,
	47, 121, 4, 151, 18, 227, 115, 164,
	198, 248, 106, 41, 31, 111, 198, 165,
	142, 93, 86, 186, 168, 174, 239, 209,
	191, 12, 246, 122, 62, 190, 65, 234,
	137, 147, 215, 245, 184, 237, 245, 225,
	217, 215, 147, 149, 79, 143, 121, 167,
	236, 253, 149, 248, 131, 82, 54, 78,
	142, 92, 25, 255, 39, 148, 197, 31,
	31, 72, 252, 230, 71, 153, 19, 101,
	149, 23, 38, 30, 147, 150, 36, 232,
	211, 162, 4, 57, 107, 232, 180, 133,
	155, 186, 246, 189, 124, 180, 172, 242,
	246, 4, 46, 111, 43, 239, 76, 60,
	1, 229, 237, 215, 255, 97, 248, 55,
	187, 63, 254, 59, 147, 19, 184, 187,
	239, 137, 133, 145, 24, 143, 192, 87,
	83, 235, 254, 5, 227, 167, 215, 145,
	234, 188, 15, 126, 161, 254, 99, 111,
	215, 219, 76, 110, 228, 33, 147, 22,
	242, 24, 143, 114, 65, 58, 110, 171,
	126, 84, 71, 246, 254, 206, 124, 99,
	237, 144, 220, 240, 199, 7, 152, 16,
	181, 99, 43, 110, 145, 142, 136, 118,
	108, 69, 50, 119, 241, 174, 131, 167,
	86, 68, 151, 223, 51, 208, 92, 129,
	84, 46, 173, 223, 34, 181, 212, 211,
	167, 241, 245, 100, 195, 184, 235, 215,
	61, 119, 239, 137, 191, 62, 195, 228,
	186, 138, 80, 126, 225, 224, 35, 245,
	127, 147, 222, 35, 205, 5, 71, 235,
	109, 7, 167, 227, 223, 30, 218, 51,
	97, 237, 91, 101, 109, 168, 223, 128,
	7, 108, 27, 234, 201, 134, 159, 152,
	119, 79, 29, 173, 196, 119, 48, 177,
	102, 144, 9, 45, 210, 22, 233, 42,
	137, 62, 77, 145, 72, 215, 135, 121,
	20, 241, 127, 88, 234, 184, 252, 81,
	169, 137, 75, 91, 27, 40, 171, 175,
	61, 48, 55, 245, 225, 221, 119, 236,
	97, 98, 131, 16, 132, 24, 230, 253,
	187, 225, 21, 233, 68, 3, 189, 229,
	235, 134, 110, 188, 101, 173, 240, 217,
	241, 51, 49, 253, 173, 114, 65, 186,
	124, 76, 99, 61, 151, 166, 52, 218,
	167, 55, 218, 133, 180, 233, 231, 141,
	109, 119, 13, 253, 146, 201, 195, 184,
	95, 129, 109, 49, 129, 242, 179, 241,
	67, 105, 187, 173, 186, 213, 86, 45,
	84, 221, 126, 114, 225, 221, 135, 254,
	60, 184, 84, 4, 4, 233, 235, 198,
	45, 210, 25, 91, 251, 68, 35, 249,
	243, 9, 126, 108, 83, 116, 219, 241,
	83, 101, 157, 191, 125, 232, 26, 105,
	231, 80, 231, 19, 41, 31, 60, 37,
	93, 219, 217, 240, 236, 214, 242, 175,
	158, 210, 116, 72, 154, 218, 68, 218,
	87, 53, 145, 33, 221, 218, 132, 180,
	162, 23, 116, 158, 204, 230, 117, 205,
	48, 89, 59, 231, 188, 22, 62, 50,
	244, 180, 141, 112, 61, 57, 79, 45,
	22, 149, 110, 174, 2, 146, 175, 140,
	68, 107, 45, 43, 202, 241, 130, 18,
	55, 16, 199, 94, 100, 224, 130, 155,
	185, 192, 135, 240, 179, 86, 3, 39,
	96, 61, 159, 12, 224, 70, 2, 110,
	35, 32, 114, 6, 0, 121, 225, 22,
	62, 22, 192, 58, 2, 238, 32, 32,
	122, 26, 64, 4, 192, 237, 60, 9,
	224, 102, 2, 54, 17, 80, 113, 10,
	64, 148, 252, 102, 3, 183, 17, 112,
	23, 1, 149, 39, 1, 84, 0, 216,
	204, 91, 137, 27, 9, 216, 70, 64,
	236, 4, 0, 162, 153, 173, 54, 176,
	137, 128, 123, 8, 168, 254, 22, 128,
	237, 29, 190, 28, 192, 54, 2, 30,
	34, 64, 248, 47, 128, 42, 0, 15,
	240, 14, 0, 247, 19, 176, 139, 128,
	154, 111, 0, 128, 161, 165, 71, 249,
	26, 0, 143, 16, 176, 151, 128, 11,
	190, 6, 80, 3, 224, 73, 251, 140,
	93, 4, 236, 39, 160, 246, 43, 0,
	23, 0, 216, 103, 155, 187, 155, 128,
	103, 9, 24, 242, 37, 128, 90, 0,
	207, 216, 55, 223, 75, 192, 243, 4,
	84, 125, 1, 96, 8, 128, 3, 124,
	49, 128, 103, 9, 120, 9, 128, 213,
	83, 64, 28, 114, 106, 158, 53, 171,
	5, 83, 205, 240, 186, 128, 248, 25,
	231, 117, 140, 55, 43, 93, 8, 19,
	228, 161, 188, 37, 121, 34, 173, 228,
	114, 16, 251, 228, 224, 136, 83, 134,
	106, 246, 24, 5, 0, 62, 221, 186,
	192, 210, 108, 33, 91, 92, 6, 192,
	47, 65, 7, 232, 51, 212, 162, 150,
	91, 165, 2, 241, 179, 212, 71, 114,
	170, 82, 36, 196, 167, 14, 7, 177,
	180, 46, 60, 163, 154, 42, 75, 44,
	80, 240, 104, 61, 19, 240, 199, 173,
	46, 77, 51, 139, 166, 161, 48, 174,
	227, 33, 191, 153, 12, 124, 40, 213,
	166, 210, 255, 222, 99, 125, 186, 161,
	173, 202, 102, 232, 28, 159, 118, 92,
	163, 149, 116, 90, 213, 233, 246, 62,
	125, 186, 183, 95, 174, 101, 233, 146,
	126, 165, 184, 71, 100, 178, 69, 53,
	223, 165, 24, 44, 210, 173, 1, 246,
	233, 208, 133, 189, 36, 23, 244, 100,
	187, 161, 229, 161, 156, 153, 94, 40,
	174, 86, 145, 216, 114, 148, 135, 42,
	90, 228, 195, 35, 243, 117, 185, 42,
	130, 164, 164, 228, 23, 47, 69, 216,
	228, 75, 34, 92, 190, 66, 64, 193,
	216, 121, 47, 182, 32, 143, 228, 73,
	144, 205, 69, 32, 87, 246, 168, 69,
	51, 171, 21, 88, 100, 118, 6, 105,
	38, 224, 143, 91, 112, 69, 161, 184,
	84, 51, 24, 207, 243, 56, 227, 237,
	17, 152, 17, 28, 2, 163, 226, 68,
	137, 84, 143, 19, 211, 151, 241, 203,
	156, 210, 100, 114, 21, 15, 113, 170,
	88, 221, 97, 189, 182, 249, 243, 179,
	165, 223, 103, 190, 192, 151, 177, 86,
	65, 201, 171, 69, 93, 73, 51, 174,
	38, 232, 179, 92, 139, 118, 26, 80,
	240, 8, 129, 79, 159, 196, 67, 77,
	111, 28, 4, 87, 242, 208, 236, 129,
	43, 76, 239, 164, 170, 199, 53, 119,
	56, 199, 39, 147, 10, 47, 20, 52,
	83, 193, 21, 34, 133, 34, 76, 179,
	209, 209, 177, 116, 111, 175, 243, 165,
	133, 119, 112, 107, 141, 150, 239, 202,
	170, 107, 212, 104, 97, 66, 90, 203,
	79, 236, 214, 38, 218, 79, 27, 154,
	169, 77, 158, 88, 52, 51, 19, 221,
	187, 244, 114, 60, 21, 166, 148, 86,
	59, 49, 98, 134, 162, 19, 169, 156,
	199, 173, 59, 32, 187, 2, 178, 105,
	231, 114, 107, 70, 213, 13, 53, 173,
	152, 92, 205, 204, 239, 90, 174, 166,
	77, 2, 189, 28, 12, 197, 248, 106,
	69, 111, 83, 139, 105, 35, 171, 155,
	136, 1, 157, 252, 61, 159, 206, 196,
	237, 168, 80, 249, 46, 28, 115, 191,
	192, 71, 130, 203, 220, 211, 119, 130,
	53, 228, 123, 32, 126, 4, 98, 225,
	12, 137, 65, 101, 226, 195, 148, 34,
	15, 65, 188, 27, 226, 200, 105, 18,
	131, 200, 196, 199, 193, 24, 242, 46,
	136, 247, 19, 189, 185, 44, 38, 238,
	35, 233, 94, 72, 159, 39, 110, 115,
	41, 76, 60, 176, 1, 210, 231, 33,
	125, 85, 224, 137, 130, 86, 80, 89,
	165, 85, 84, 11, 25, 213, 152, 165,
	177, 68, 145, 170, 223, 187, 160, 35,
	110, 55, 88, 179, 157, 164, 190, 28,
	183, 86, 179, 171, 144, 174, 169, 89,
	90, 191, 7, 2, 192, 201, 103, 94,
	231, 245, 35, 55, 243, 205, 101, 89,
	35, 211, 174, 24, 220, 44, 57, 143,
	34, 170, 117, 65, 191, 27, 80, 31,
	136, 88, 135, 205, 35, 204, 246, 218,
	40, 63, 94, 111, 206, 193, 29, 222,
	192, 29, 222, 23, 184, 200, 71, 193,
	101, 156, 139, 239, 209, 197, 222, 135,
	240, 99, 162, 89, 203, 33, 127, 241,
	35, 112, 166, 124, 12, 210, 79, 169,
	37, 156, 117, 152, 95, 252, 132, 106,
	230, 99, 72, 191, 34, 135, 157, 113,
	29, 246, 57, 189, 246, 51, 72, 79,
	147, 195, 78, 187, 14, 59, 241, 24,
	164, 167, 193, 148, 85, 96, 202, 145,
	149, 167, 44, 193, 161, 252, 10, 190,
	7, 28, 90, 69, 28, 218, 96, 247,
	130, 147, 46, 229, 139, 28, 79, 64,
	6, 96, 20, 145, 171, 98, 251, 98,
	54, 238, 202, 66, 158, 178, 57, 173,
	157, 43, 134, 146, 71, 130, 80, 174,
	115, 38, 84, 224, 31, 155, 10, 123,
	114, 102, 17, 158, 241, 103, 15, 215,
	51, 106, 47, 17, 81, 86, 99, 188,
	48, 152, 139, 81, 67, 133, 52, 222,
	139, 131, 98, 150, 251, 142, 5, 28,
	132, 62, 35, 87, 84, 87, 39, 150,
	169, 6, 197, 218, 84, 86, 168, 51,
	193, 58, 124, 190, 9, 137, 220, 163,
	54, 219, 153, 237, 91, 230, 112, 221,
	76, 131, 107, 249, 78, 59, 90, 9,
	197, 48, 75, 131, 242, 154, 165, 40,
	177, 115, 57, 10, 204, 8, 63, 48,
	251, 22, 7, 41, 231, 21, 210, 1,
	52, 39, 121, 63, 100, 135, 17, 44,
	193, 77, 228, 131, 93, 65, 22, 138,
	17, 193, 137, 202, 203, 228, 255, 151,
	32, 60, 10, 97, 5, 119, 130, 242,
	38, 61, 254, 42, 132, 239, 226, 149,
	149, 60, 52, 149, 138, 111, 27, 76,
	16, 163, 235, 208, 79, 17, 254, 131,
	135, 160, 119, 24, 122, 111, 156, 163,
	94, 83, 166, 98, 116, 171, 196, 226,
	254, 172, 226, 250, 45, 139, 158, 103,
	44, 5, 149, 197, 84, 168, 87, 67,
	189, 26, 210, 188, 106, 46, 211, 50,
	78, 216, 98, 144, 197, 240, 10, 157,
	194, 85, 46, 48, 84, 42, 29, 112,
	57, 107, 134, 211, 59, 53, 11, 158,
	209, 86, 195, 127, 130, 157, 238, 102,
	169, 83, 201, 230, 200, 93, 118, 148,
	41, 200, 254, 8, 132, 215, 203, 35,
	120, 136, 41, 197, 150, 214, 128, 37,
	197, 241, 73, 235, 151, 119, 222, 39,
	31, 120, 103, 195, 11, 224, 169, 209,
	214, 139, 255, 121, 229, 226, 97, 123,
	205, 135, 153, 56, 102, 180, 85, 127,
	172, 227, 211, 210, 173, 171, 94, 98,
	226, 133, 147, 173, 59, 71, 77, 60,
	182, 67, 173, 59, 201, 196, 145, 139,
	173, 207, 55, 78, 108, 170, 191, 225,
	153, 67, 248, 50, 182, 15, 44, 189,
	66, 233, 86, 83, 206, 172, 21, 203,
	104, 233, 152, 169, 116, 55, 19, 209,
	118, 91, 233, 158, 162, 169, 229, 205,
	18, 139, 232, 46, 145, 99, 96, 15,
	241, 118, 20, 52, 93, 231, 210, 244,
	228, 102, 215, 228, 208, 0, 231, 190,
	156, 149, 27, 225, 102, 216, 233, 26,
	67, 32, 236, 190, 22, 12, 192, 124,
	108, 162, 179, 164, 171, 114, 157, 159,
	57, 10, 133, 249, 58, 132, 111, 25,
	220, 35, 56, 153, 163, 254, 9, 178,
	101, 144, 153, 148, 36, 110, 230, 172,
	36, 94, 54, 33, 92, 71, 181, 239,
	102, 206, 90, 98, 209, 94, 8, 111,
	22, 104, 246, 80, 138, 72, 104, 216,
	130, 191, 160, 223, 243, 217, 69, 138,
	128, 106, 164, 138, 51, 21, 212, 134,
	31, 9, 95, 161, 173, 199, 80, 186,
	178, 185, 108, 4, 249, 238, 70, 60,
	97, 194, 76, 158, 8, 76, 71, 0,
	19, 253, 41, 170, 205, 110, 246, 168,
	146, 110, 109, 64, 87, 161, 43, 93,
	12, 163, 38, 81, 49, 240, 208, 146,
	36, 142, 111, 101, 194, 57, 243, 177,
	47, 173, 33, 33, 123, 205, 240, 33,
	211, 237, 170, 116, 120, 176, 214, 63,
	97, 6, 149, 91, 27, 78, 104, 15,
	202, 109, 30, 81, 219, 92, 200, 174,
	165, 114, 27, 229, 56, 109, 33, 177,
	32, 250, 190, 124, 221, 185, 154, 153,
	61, 252, 20, 29, 102, 241, 6, 34,
	123, 134, 193, 4, 227, 57, 170, 95,
	79, 203, 229, 38, 120, 73, 159, 75,
	80, 206, 147, 93, 181, 150, 75, 154,
	51, 232, 238, 211, 156, 161, 196, 25,
	209, 65, 153, 226, 108, 42, 240, 89,
	144, 118, 58, 93, 77, 176, 9, 83,
	148, 23, 7, 214, 165, 210, 118, 140,
	64, 84, 37, 173, 199, 40, 170, 185,
	165, 196, 102, 94, 223, 96, 145, 50,
	84, 36, 120, 49, 32, 99, 49, 17,
	20, 18, 228, 61, 202, 45, 219, 28,
	59, 187, 168, 21, 222, 128, 19, 114,
	158, 57, 228, 149, 44, 117, 140, 28,
	164, 189, 212, 49, 206, 184, 189, 161,
	135, 76, 215, 33, 189, 145, 58, 198,
	105, 183, 55, 148, 90, 131, 180, 115,
	27, 227, 92, 141, 165, 52, 189, 11,
	249, 63, 168, 1, 242, 185, 154, 131,
	4, 148, 239, 14, 145, 172, 210, 159,
	51, 253, 135, 194, 237, 142, 218, 130,
	179, 243, 132, 18, 105, 120, 144, 72,
	162, 23, 231, 241, 116, 165, 113, 16,
	94, 41, 240, 72, 54, 220, 131, 151,
	130, 232, 11, 105, 150, 82, 175, 214,
	122, 10, 102, 185, 99, 218, 109, 11,
	156, 99, 206, 147, 77, 201, 32, 98,
	62, 121, 203, 29, 231, 201, 166, 115,
	83, 45, 28, 148, 213, 179, 106, 129,
	113, 211, 143, 163, 55, 121, 10, 238,
	228, 153, 244, 70, 75, 219, 192, 129,
	156, 50, 211, 94, 32, 156, 74, 40,
	63, 193, 249, 19, 65, 203, 150, 243,
	141, 112, 110, 27, 238, 224, 118, 183,
	68, 31, 230, 69, 191, 13, 247, 115,
	152, 82, 202, 105, 10, 207, 12, 56,
	180, 181, 92, 92, 230, 184, 113, 153,
	37, 184, 149, 92, 232, 119, 213, 78,
	165, 43, 167, 82, 79, 241, 39, 113,
	255, 7, 44, 119, 18, 247, 91, 44,
	215, 229, 54, 30, 250, 13, 78, 92,
	223, 26, 108, 51, 226, 218, 142, 96,
	235, 2, 1, 6, 155, 150, 88, 74,
	6, 219, 21, 210, 57, 88, 168, 196,
	149, 173, 193, 14, 37, 230, 91, 131,
	205, 68, 204, 46, 14, 182, 30, 49,
	219, 26, 108, 58, 162, 154, 12, 246,
	27, 81, 25, 27, 132, 84, 92, 98,
	132, 126, 28, 89, 210, 26, 92, 68,
	92, 100, 120, 131, 159, 184, 104, 77,
	48, 223, 137, 139, 14, 5, 67, 139,
	184, 164, 163, 207, 217, 241, 85, 171,
	53, 216, 216, 18, 68, 45, 41, 103,
	238, 75, 57, 193, 238, 235, 112, 118,
	195, 62, 167, 60, 84, 171, 45, 180,
	95, 245, 57, 201, 172, 166, 28, 138,
	76, 204, 193, 74, 102, 185, 239, 237,
	100, 205, 118, 34, 246, 57, 241, 203,
	88, 222, 40, 206, 154, 237, 97, 220,
	242, 214, 47, 111, 96, 181, 58, 93,
	162, 1, 207, 244, 31, 219, 103, 184,
	147, 23, 47, 148, 89, 116, 198, 133,
	23, 157, 75, 32, 184, 34, 188, 232,
	32, 59, 166, 183, 187, 29, 116, 177,
	183, 232, 192, 118, 61, 29, 44, 56,
	134, 247, 229, 187, 44, 56, 134, 206,
	211, 161, 5, 71, 240, 26, 46, 178,
	124, 2, 181, 87, 167, 72, 234, 236,
	154, 165, 102, 132, 146, 184, 16, 117,
	194, 5, 113, 36, 214, 11, 30, 17,
	135, 97, 157, 72, 45, 197, 104, 162,
	102, 44, 13, 180, 69, 62, 98, 17,
	124, 193, 2, 139, 212, 197, 110, 144,
	72, 99, 72, 31, 248, 243, 128, 127,
	98, 68, 79, 122, 14, 235, 239, 47,
	110, 156, 155, 193, 252, 66, 73, 6,
	37, 75, 4, 230, 214, 72, 106, 85,
	182, 160, 134, 75, 244, 92, 203, 114,
	108, 194, 124, 111, 141, 243, 150, 169,
	75, 199, 6, 199, 132, 150, 169, 241,
	55, 5, 84, 137, 149, 71, 211, 209,
	93, 144, 21, 237, 200, 20, 44, 111,
	198, 204, 172, 154, 203, 248, 51, 94,
	127, 58, 166, 188, 27, 196, 147, 29,
	1, 79, 14, 225, 94, 139, 153, 55,
	58, 96, 202, 33, 130, 215, 98, 250,
	81, 165, 238, 216, 63, 155, 113, 255,
	122, 49, 156, 52, 128, 2, 206, 63,
	231, 135, 103, 112, 74, 247, 255, 151,
	198, 91, 131, 73, 225, 187, 209, 120,
	223, 10, 181, 68, 49, 47, 215, 138,
	189, 210, 179, 31, 101, 3, 3, 212,
	21, 196, 220, 255, 229, 78, 108, 89,
	227, 254, 130, 241, 83, 24, 226, 12,
	169, 40, 201, 24, 209, 113, 104, 60,
	233, 87, 168, 3, 55, 203, 255, 5,
	0, 0, 255, 255, 55, 116, 202, 122,
}
