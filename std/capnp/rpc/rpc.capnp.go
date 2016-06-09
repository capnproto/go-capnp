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

const schema_b312981b2552a250 = "x\xda\x9cX}p\x14\xf5\xf9\xff~\xf7.\xb9\x10\x92" +
	"\xdcm6\x81\x1f\xfc`\x02(\xadD\xde\x82\xda\xdak" +
	"\x9d\x83\x18\x18`\xa0\xe4\x92P\x85Jus\xf7%\x1c" +
	"\\v\xcf\xbd\x0d\xe4R3\x01\xab\x8eR\x99\x02\xbeT" +
	"\x19\xf1mt\xaa\x16GD\xa8\xd8\x9a*\x0c\xd6\x97\xf1" +
	"\x8d\xfa2\xda\xd1)8\xe3T::\x15\xdf\xca\xfb\xf6" +
	"\xf3\xec\xee\xed.\xc91\xd6\xfe\x91\xc9\xde\xf3<\xfb}" +
	"\xde?\xcf\xf3\xdd\x99c#\xb3\xa5\xa6\xb2%#\x18K" +
	"f\xcb\xca\xad\xb7\x17m\xef}\xbd}\xf5\xafX\xb2\x92" +
	"\x87\xac\xd6\x07\xdb&\xff\xff]\xb5O\xb1\xb2P\x841" +
	"eS\xb8O\xd9\x12\xc6\xd3E\x9b\xc2\xbf\xe1\x8c[;" +
	"\xf7\xde4\xf2\xc0\x87\xe7\xddH\xd2\xdc\x97\x9e\xcb#\xe5" +
	"\x10\x97\xcb\xf7+c\xcaI\xbc\xbe\xfc\x0a\x12\x9f\xb5s" +
	"\xd3\xfa\x86\xfb\x9e\xde2\\\xbc\x06\xe2\x85\xc8VeC" +
	"\x84\xc4\xfb#\xa3C\x10\xdfwR\xb9\xb2\xa3\xee\xd9;" +
	"\x86\x8bK\\R\x06+\xf7+/T\x92Y\xfb*\xd7" +
	"A\xfa\xc7\xe6\x9d\x97MRk\xb61\xb92 \\&" +
	"\x91\xc4\xe4\x91[\x95i#\xe9i\xcaH\x92]\xbec" +
	"\xdf\xc95\xe1\xd5\xf7\x0c9\xd9\x11\xbe\x1b\xc2\x0f\xd8\xc2" +
	"\xdbG>\x01\xe1\xf8\x15O]\xb6i\xd7\x98{IX" +
	":\xdbI\x1eR.\xab\xda\xa8\xcc\xad\"\xab\xe7T\xbd" +
	"HN\xfe\xd6|\xb3\xbf:;\xf6\xf1!gS\xd8\x94" +
	"\xa6\x9a\xad\xca\x8fj\xe8\xe9\x92\x1a\xb2\xe3\xca\xc1E\x89" +
	"\x8f\xee\xbcu\x17\x93\xeb$kl\xe6\x8dx\xf9\xd3\x93" +
	"\xdfe\x8c+w\xd4\xbc\xa2<`\x0bn\xaf\xe9b\x81" +
	"cx\x98I\xcak5m\x17\x1d\xac\x19\xcd\x95\xc1\xe8" +
	"\x12&YS\xaf^\xff\xdc\xbd\xc7\xff\xf2\x0cK\xc6\xca" +
	"\xb8\xb5\xe1\xa73\xf7\xfe\xb2\xff\xd8 \x1dS\x1f\xfb\xab" +
	"21\x86c\xda\xc7\xc5Bd\x9cVq\xcb\x89\xa5w" +
	"\xee\xffS\xe9\x90V\xc7\xb6\xe2\x15\xd2*\xc7\xc8s\xef" +
	"(^\x09\xadO\xc6\xaeW\xf6\xc4\xbe\xaf\x1c\x8c\x91\xd2" +
	"C\x17\x16\xce\x8b\x0d<\xfag\x96\x1c\x01\xa57g?" +
	"<\x93\x1c\xd7x\x90\x94N\x937*\x97\xc8\xa4t\xa6" +
	"l+\xf5\x98<\x84c&\xca\x0b\x95\xc9\xf2:\xa5_" +
	"n\xc01\xfd\xd2\xe7\x87OGro\x0d\xcd\x04\xb7\xeb" +
	"@\xae\xe5\xca&:J\xb9E\xa6x\xa5j\x8e\xed\xdf" +
	"5\xbd\xff\xadR\xb1\xfd\x18j?\xb3e\x8f\xd8\xb2\xa3" +
	"f/\xdd\xdc\xb9\xe7\xe5\xb7K\x9d\xac$k7*\xcb" +
	"j\xe9ii-y\xba\xf8\xc3\x9f\x8b\xbf\xef\xee|\x87" +
	"%\xeb!,\xffp0\xfa\xeb\x1f\xa4\x8f\xb3\xa5<\xc2" +
	"\xc3\x08\xcc\xd7\xb5\xff\x84c\xc7k\xff\x01Q/M\xa5" +
	"\xce}_yP9\xac\x8c&#\x14\x08?w\xff8" +
	"\xfd\xb5w\x9f|\xaf\x94\xe8ku\xaf(\xef\xd7\x91\xe8" +
	"\xe1:\xb2\xf7\xee\xab\x7f?\xf6\x9b\x9d\x9f\xfc\x8d%\xa3" +
	"\xe8<\xaf\x0f\x97\x86\"<\x84:[ZO&,\xab" +
	"'k\xf7\xb4\xd4|\x8f\xffa\xe6\xe1\xe1\xb1?^\x7f" +
	"\xbdr\xba\x9eb\x7f\xac\xde\x8e\xfd\x01mt\xd3\xfa7" +
	"\x16\x1d)i\xed\x91\xfa\x07\x95\xa3$\xad|VO&" +
	"l\xd8\xfc\xb3\xfa\x96\xdbG}\xc9\x92c\xb8g{K" +
	"D\xa20\x8d\xfaHQG\x91\xe8\x8aQ$\xea\x85\xa8" +
	"\xd4\xb9\x83\xa3\x1eS^\xb0\x85\xf7\xd9\xc2O\xf0C\x9b" +
	"\xc3w\x1d>Y\xb2\xdd\xf8\xe8>\xa5l\xb4\xf3D\xce" +
	"\x19\xb9\xd4\xf4\x94\x9a\xd3X\"\x17\xbf\\\xcdf[9" +
	"O\x8e\x0b\x85\x19\x0bs\xc6\xe4=\xcb\x01T\xbbC<" +
	"\xf9\xbc\xc49\xaf\xe3D\x1b\x8c\x83\xb6\x17\xb4\x03\x12\x97" +
	"%\x10a\xb0\xbc\xaf\x13\xc4\xe7A|\x15\xc4\x90T\x87" +
	"\xdac\xf2\xcb\x0bA|\x09\xc4\xb7A,\x83$\x8e\x95" +
	"\x0f\xd2\xeb\xaf\x82\xf8\x1e\x8e,\xe7\x81L\xc8\xef\x18L" +
	"\x92\xc3\xeb\xebx\x05\xe7\xf2\xbe\xfd\x90;\x00\xb97%" +
	"n]\xdb#\xf2fF\xd7XhA\x9aW0\x09\x7f" +
	"<a\xaaF\x970y\xccG.\xc6y\x0c>e4" +
	"S\x18+\xd5\x14\x8b\x08\x88\x8f\x80\xf8\x08P\xbb\x85\xb9" +
	"JO/H3\x88E@\x8b\xe0\x88\x9cj\xa8\xddy" +
	"\x1c\xe1\xc1\x99{D^h\xe96\x91\xefa\x0dY3" +
	"\xdf\xa1[\x88\x8c\xbe\xaecUF2\xd2\xad\xaaa\x16" +
	":\xd4L\x96\xc2\x05q\x1c\xc5\x03\x81\x94(\x8e\xb9\x16" +
	"\x91O\x19\x99\x9c\xa9\x1b\x8c\"\xfa\x7f\xa1p\x95e\xd9" +
	"!\xbd\xbb\x11~\xdd\x0e\xbf\xee\x97\xf8x~\xc6r\xa3" +
	"\xba}5\xc8\xf7\x80\xfc\x08\xc8\xd2i\xcb\x8d\xeb\xc3\x06" +
	"\xc8\x0f\x81\xbc\x13\xe4\xd0)\"Sd\x1f\xef\x03y\x07" +
	"\xc8{%^\x1d>i9\xa1\xdd\xd3\xe7g\xab\xba\xec" +
	"\x04\xa8e\x94\xaf\x8d~j\xa2\x9a\xae\x09Vn\xbb'" +
	"\x8c\xf9:\x8b\xe6M\xe1E\xd4%\xb7\x1a\xacA\xef\xce" +
	"\xe4\x85G7DJd\xd6\x0a\x83%\xe6\xebg\xbd\xe0" +
	"3\xe6h\xf9u\xc2\xe0\xb1b\x1d\xbbq4We\xec" +
	"\x88q\xb3\xe0\xbc\x8a\x00\xc7|\x18r\xa5\x8a\xb1\xe3\xb9" +
	"\xf8b\x91\xcf\xab]\\P\xd4.\xf5\xa2\xa6\x148\x02" +
	"\xd1\xde\x8b\xdel\xbf\x81\xc3;\x04\xce\x8e\x9b\xb2\x81\xcf" +
	"\x02\xe3:b\xdcL\x8c\xd0i\xcb\x8e\x9cr#G\xa0" +
	"\xdb\xd7\x13\xe3Vb\x84OYv\xec\x94[8*\x10" +
	"\xa7\x80\xb1\x99\x18en\xf8\x94M6\xe3fb\xdcN" +
	"\x8cr7\x82\xca\x16\xde\x0c\xc6\xad\xc4\xb8\x8b\x18\x91\xe3" +
	"`\xd0\xc8\xbd\xc3fl&\xc6=\xc4\x18q\x0c\x0c{" +
	"\xa8q\xa4\x13\xc2`<D\x0c\xe9\xdf`T\x80\xf1\x00" +
	"o\x03\xe3~b\xec F\xe57``\x1bP\x1e\xe5" +
	"H^\xfb#\xc4\xd8M\x8c\x91_\x83Q\x09\xc6\x93\xb6" +
	"\x8e\x1d\xc4\xd8K\x8c\xaa\xaf\xc0\x18\x09\xc6\x1e\xdb\xdc\x9d" +
	"\xc4x\x96\x18\xd5_\x82Q\x05\xc63\xb6\xe7\xbb\x89\xf1" +
	"<1*\xbe\x00\xa3\x9a\x80\x82\xa3\x9d!\x0c\xc6K`" +
	"X=Z\xa6;\x97\x15\xdd\xacAh\x94\xd5\x98\xbf2" +
	"8\x89iP;u\x83:,0,\x89\x1eM\xa1\xf4" +
	"A\xf6`\xd3!'\x0ca\xf6\x18\x1a\x18\xde\x10w\x19" +
	"+3Z&\xbf\x0a\x0co\xa48\x8c\x01C\xe4\xf5\xec" +
	"Z\x01\x8e7+=NV\xa8y\xe2x#\xde\xad\x16" +
	"\xbd\x13\xef\x08S\xb0h\xbb\x8aWkQ\x8b\xb5 w" +
	"\xea\xba\x997\x0d\x95\xf1\x1c^\xf2\x90x\xe8K\x89\x16" +
	"A\xff\x8b\xaf\x0d\xe4\x0c}m&Mz\xbc5\xc55" +
	"ZM\xa5D\x8e\xbc\xf7f\x9b\xeb\xfdj=CNz" +
	"8\xeb\xaaH\xa3e\xba;U\x83\x85\xbat\xb0=\xcc" +
	"\x1eR\xe4R\xb1\xc8E\x87\x0d`6@T\xf8\x001" +
	"\x85\xa0\xf4\x02\xf4\xeb\xc5\x81:\x97\x9b\xa8\xb7g\x82\xfa" +
	"\x13$\x0eiC^\xd0L\x11@\x8d\xd7\x8c\xf0\x84\xba" +
	"6}\xcef\x0c\xb4Y\xabZ\xc8\xea*O\xbb\xba]" +
	"\xb8\x9f\x82RK\x9e\x0f%3\x01\xd8E\xbc\x9fF(" +
	">\x15\xc4\xf9\x12\x1fH\xe9\xa8\x14\xcd\xf4\x82\x8e\xe3:" +
	"\xd4\xce\xac P\xada\xbc5\x04M\xfe\x9e\x0a\xbd5" +
	"C\xf4\xda\xd1v\xda\xbb\xca\xd3;\x97\xc6L\x0bT\xb4" +
	"\xfacf1\xcd\x89\xf9\xa0u\x04\xc6L\x12\xdd\x93\x84" +
	"\x92\xe4U\xdfy(\x00\xaa2\xb9\x8c\xd0\x18\xf7\xad\x0f" +
	"\x18\xd6f\x97.\xb3\x931\xc13\xec \xf9\xfe&\xf4" +
	"}@\x01\x99\x00\xcb0\x99\xde'@\xfd\x00\xc4O\xa8" +
	"\xb3-\x07o\xe4\x8f)v\x87@\xfd\x94P\xe8\x8c\x03" +
	"6\xf2\x112\xf8\x13P\xbf\"\x08:\xed\x02\xf5Q:" +
	"\xf6sPO\x11\xfe\x9cr\x81\xfa\xf8c\xa0\x9eBs" +
	"V\xa09\xc7\x97\x9f\xb4$\x07e\xca\xf8.\xb4m\x05" +
	"\xb5m\x9d\x0d?'\\\x94\x919\xde\x00\x0d\x8c\x09\xd4" +
	"\xcf\xaa\x9dvg\xc2\xf9\x08m\xb7Q+\xa7I\x87j" +
	"\xc93{d\x95\xd1\xcc\xa2\xee\xeb\xc1|+1\xffD" +
	"/\xd5~Fg\\\x1b\xde\xfe\xc8\xba\x96\xc2\xb9P\x14" +
	"\xb1\xdc3\xda9\xcabn6/\xd6EW\x09\x83f" +
	"\x8c\xa9\xae\x11\xf3P\x92|\x89\x09J\xb2G4\xd8\xd9" +
	"\xf2,s\xdak\x9e\xc1\xf5\xee\x0e{JDi\xb0\x96" +
	"\xce\x0d\xf9\xe0\x14M\xa0X\xc7\x96*\xd6>\xb7X/" +
	"\x95x(\x13\x1cT+a\x95\x96b\x09q\xb9\xde\x83" +
	"\x02\xf6\x18~W\xceu}\xd6\xa6w\x14r\xc2)\x85" +
	"\x98\x9d\xdb)q\xf2\\\x9e\x882\xe5\x92<\x1e\xd8\xce" +
	"C\xf2\x18L\xa4\xc4J\xec\x01\"m\xe9\x98\x81\xe8\xa7" +
	"4\x0b\xe1\x07p\x00m\x829\x1bM\xa1G\x87\xa2\xac" +
	"\xbfq\xf1\\\xb2\x85\x07\xaei\xf2\x86f\x1f\xb6\xe4\xfe" +
	"6\x1f^\xe5\xfeF\x1fR\xe5B\xdc\x87Q\xb9'\xee" +
	"#\xa7|m\xb3\x0f\x96rw\xb3\x0fArf\xb9\x0f" +
	"or\xa6\xd9\x874Y\xc4} \x93\xd5F\xbfo\xe4" +
	"\x15F\xe0\x96\xb7\xa2\xd9oky\x99QD\x16yY" +
	"\x9f?\xcd\xe5e\xfb\xfdR\x91W\xb4\x0d\xb88g5" +
	"\xfb\xd0\x1c\xa5\xad)\xe1t[b\x9e=\x16\x06\xda\x9c" +
	"!0\xe0\xe4YX-\x01 \x1dp\x10C$\xe6\xd8" +
	"\xe5\x12]\x08\xec\xb5\x8a\xf8\xc9\x1a\xecn\x1fp\xd0," +
	"m\x15\x17/\xd6`\xaf^V\xeb\x10D\xb4:\xdcu" +
	"$d\x16\xce^\xd2\x8a\xb9\xe7Z\xb2\x8a\x07\xefx\xe3" +
	"$>g*\x0f\xdc\xbf.\x00\xe1\xe2\xe0\xa5\x00\xe57" +
	"\x07\xa5\x82\xc9+\xf3\xe5\x04\x88Z.\x1e\x87\xed\xb9\x14" +
	"\xca\xc4\xa6N\x8a\x18\xc5\x1fM\xbc\x8d[}zwg" +
	"F\xf4\x89\xb06=\xa5w\xcf\xe8\xd2g\xd8o\x19\xba" +
	"\xa9\xcf\x9a\x917\xd3\xce\xcf\x19F\x8e\xe3-\xe7\xc4\x19" +
	"\xa9\x0b\xa5\x0b\xa7;\x87kj\xb7\xc8\xe7\xd4\x94\xdd\x10" +
	"8T:\xbbW(\x96\xc3\x00\xb6\xcd\x07\xd8jn\xb9" +
	"3e\xf1$\x1fb\xab\xa53V\x09\x8cug\xca\x02" +
	"\xc6\xbd^\x8a@\xd3\x10\x90\xffo\x10\xc3q\x82\xbbN" +
	"\xb0$v|?\xcc\xf2\x886\xeb\xf5-G\xcf\x14~" +
	"\x97\xfe\x02?\x1a\xad\xa2\x8f\x8c\x8b(=\x97\xc8\xca\xcc" +
	"`V\xa6\x82pi0+\x18\x9es:\xdc\xacl+" +
	"fE\xe5\x9a\xa6\x9b*\xac\x0ciy?;\xa9\xde\xde" +
	"\xef\x9e\x9dT/\xc7[V\x97\xee8\xc4\xe3\xb0v\x0d" +
	"\x8a\xd2\xc6\x0d;+E\x16\xe089\x8e\x07\xac\x95\x9b" +
	"\x9a}K\xe5iq\xeb\x17\xb7\xdd\x97\x1c|w\xe3\x0b" +
	"@\x9aI\xd6\x8b\xffz\xe5\xfc1\xbb\xcd\x87\x99<y" +
	"\x92U{\xa8\xed\xd3\xc2Mk_b\xf2\xc4Y\xd6m" +
	"\x13f\x1c\xda&b'\x98<~\xb9ut\xd3\x8c\xd1" +
	"\xb5\xd7<\xb3\x1f?\x1a\x07\\\xdd\x09g1\x88\xa4\xf5" +
	"T\xc4T\xbb\x1a\xc8\xd9.+\xd5\x937\xf5n\xb3\x80" +
	"\x92t\x83\x89Kv va\x84*\xe6\x86jV\x83" +
	"k\xb2WS\xa1\\\xbc\xd85g7\x0d7\xce\x8d\xc7" +
	"\x1e\x1c\xc7\xfd]\x86\xe0\xd8\x05\xf7\xc4\xda\x8c&\x16\xa4" +
	"\x87\x810\xea\xd7\x81\x04\xc6\x86\x9c\xbd\xdc?\xc7\x9b\xc3" +
	"M[A\xbc\x18\xc4\xd9\xe7X\x06\x8a\xc3\xaf\x8d\xdb3" +
	"\x8av\xa5\xbc7\xfc\x82J\x1dlq\x94~\xcbVB" +
	"\xcd\xb1\x08\xb4+i+\x99\xe0t\xcc\xd2\xe6o\xd9J" +
	",{\xc9\xcc;\xcdQ\\<m\x88\xc3\xa6X\xe2\x02" +
	"Y\x04@\xaa ]\x8b\x9a\xa2\xd7L\xc6\xec\x0d\xd1\xb1" +
	"B\xa5)w\x0d\x14f\x8b+\"\x99\x91\xa1\xbd$\x0b" +
	"j/\xb5\xf3iw\x03\xe9\xa1\x14\xe4@\xbd\x8e\xf6\x92" +
	"S\xee\x06R \x93MP\xd7K\xc5k\xdf\"\x9d%" +
	"\xf4\\'*i\xd8\xf5\x8e/\xd2\x1d\x8e\xbfX\xb8\xdb" +
	"1+\xf7\x16\xe8\x12\xc9tP5\x82h\xa0\xe4\x82\xdf" +
	"\xdfxc\x94f,9\xe5\x06[%3\xaf\x82A\xab" +
	"\x10\x0d\xc9qS\xfc\x11\xb4U\xa0\x99\xf4Q\xc1]\x01" +
	"\xaf\xdd\xe6[.s\xf7KC?]\xaa{A\xbcA" +
	"\xa2[\x88\x9a\xc7\x9e\x81\x16\xc4\x9f\xbf\xf9\xf3\x05y\x9a" +
	"=\xc2H\xe4\xe7\xa9(\x07/\xf0\x9e@K\x8f\xa1v" +
	"f\xb2\x19L\x87\xe2\x17\x82\xa8\x093y\xd47\x1dx" +
	"\x16=;Y\xc5!\xe3\x8c\x18\xd8A\xaez\xdfvd" +
	">6\xb4$w\x8eR.VUS\x9b\xbb\xdc/:" +
	"W\x01axj\xf9\x95\x18T\xbc\xdb\xdf\xb3=%C" +
	"\xf6l\xc9\xf9\x943\xbd\xf8\x11#\x1b\xa5o\x18T\xd9" +
	"v\x05\xd1\xae9\x97\xc2=\xdb\xd1\xe8T\x106My" +
	"\xc1B\x7f \xd0G\x08\xc9\xde3\xe5\xe4r\xbf\xbe\x13" +
	");\x86\xd8\xef\x0az\x8f\x91\x17\xd9\x95\xb4\x04\x16\xaf" +
	"\xf9,\x14\xd8\xe0|ht`\x89\x95\x18W\xcef\x10" +
	"1\xd4\xdc\xb9\x1b\xde\x8b\xd2\xb6o\xeb\xf7\xb4\xc8\xa1d" +
	"U\x93\x8b\xf4\x92\xce\xd5\"e\x12s\xe8B9,e" +
	"\x91\xe9KrC\xef`\x8d>\x96\x05>\xd2L\xbb\xde" +
	"\xdf.\xa3\x1az\x02a\xc0\xfe\xd1\x8a\x9d\x04:\x8dy" +
	"\x19\x91M{\x1f\x97\x82n:\x0d\x1d\xa5\x8e\x1e\xe2g" +
	"<\x08\x9a<\xf0iT\x9e\xd6\xcc\xa4s^g\x9c\x8b" +
	"X\xafy\xd6\xf7;\xda\x8d\xfe\xd7\x8bU\xb3\x8fk\xdf" +
	"\xedb5\xb0F\x14h6\x14\xe3\xfc\x9f\x00\x00\x00\xff" +
	"\xff\x97q\xc8^"
