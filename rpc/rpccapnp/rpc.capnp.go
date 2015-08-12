package rpccapnp

// AUTO GENERATED - DO NOT EDIT

import (
	strconv "strconv"
	C "zombiezen.com/go/capnproto"
)

type Message C.Struct
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

func NewMessage(s *C.Segment) Message {
	return Message(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func NewRootMessage(s *C.Segment) Message {
	return Message(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func AutoNewMessage(s *C.Segment) Message {
	return Message(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func ReadRootMessage(s *C.Segment) Message { return Message(s.Root(0).ToStruct()) }

func (s Message) Segment() *C.Segment {
	return C.Struct(s).Segment()
}

func (s Message) Which() Message_Which {
	return Message_Which(C.Struct(s).Uint16(0))
}

func (s Message) Unimplemented() Message { return Message(C.Struct(s).Pointer(0).ToStruct()) }
func (s Message) SetUnimplemented(v Message) {
	C.Struct(s).SetUint16(0, 0)
	C.Struct(s).SetPointer(0, C.Pointer(v))
}
func (s Message) Abort() Exception { return Exception(C.Struct(s).Pointer(0).ToStruct()) }
func (s Message) SetAbort(v Exception) {
	C.Struct(s).SetUint16(0, 1)
	C.Struct(s).SetPointer(0, C.Pointer(v))
}
func (s Message) Bootstrap() Bootstrap { return Bootstrap(C.Struct(s).Pointer(0).ToStruct()) }
func (s Message) SetBootstrap(v Bootstrap) {
	C.Struct(s).SetUint16(0, 8)
	C.Struct(s).SetPointer(0, C.Pointer(v))
}
func (s Message) Call() Call     { return Call(C.Struct(s).Pointer(0).ToStruct()) }
func (s Message) SetCall(v Call) { C.Struct(s).SetUint16(0, 2); C.Struct(s).SetPointer(0, C.Pointer(v)) }
func (s Message) Return() Return { return Return(C.Struct(s).Pointer(0).ToStruct()) }
func (s Message) SetReturn(v Return) {
	C.Struct(s).SetUint16(0, 3)
	C.Struct(s).SetPointer(0, C.Pointer(v))
}
func (s Message) Finish() Finish { return Finish(C.Struct(s).Pointer(0).ToStruct()) }
func (s Message) SetFinish(v Finish) {
	C.Struct(s).SetUint16(0, 4)
	C.Struct(s).SetPointer(0, C.Pointer(v))
}
func (s Message) Resolve() Resolve { return Resolve(C.Struct(s).Pointer(0).ToStruct()) }
func (s Message) SetResolve(v Resolve) {
	C.Struct(s).SetUint16(0, 5)
	C.Struct(s).SetPointer(0, C.Pointer(v))
}
func (s Message) Release() Release { return Release(C.Struct(s).Pointer(0).ToStruct()) }
func (s Message) SetRelease(v Release) {
	C.Struct(s).SetUint16(0, 6)
	C.Struct(s).SetPointer(0, C.Pointer(v))
}
func (s Message) Disembargo() Disembargo { return Disembargo(C.Struct(s).Pointer(0).ToStruct()) }
func (s Message) SetDisembargo(v Disembargo) {
	C.Struct(s).SetUint16(0, 13)
	C.Struct(s).SetPointer(0, C.Pointer(v))
}
func (s Message) ObsoleteSave() C.Pointer { return C.Struct(s).Pointer(0) }
func (s Message) SetObsoleteSave(v C.Pointer) {
	C.Struct(s).SetUint16(0, 7)
	C.Struct(s).SetPointer(0, v)
}
func (s Message) ObsoleteDelete() C.Pointer { return C.Struct(s).Pointer(0) }
func (s Message) SetObsoleteDelete(v C.Pointer) {
	C.Struct(s).SetUint16(0, 9)
	C.Struct(s).SetPointer(0, v)
}
func (s Message) Provide() Provide { return Provide(C.Struct(s).Pointer(0).ToStruct()) }
func (s Message) SetProvide(v Provide) {
	C.Struct(s).SetUint16(0, 10)
	C.Struct(s).SetPointer(0, C.Pointer(v))
}
func (s Message) Accept() Accept { return Accept(C.Struct(s).Pointer(0).ToStruct()) }
func (s Message) SetAccept(v Accept) {
	C.Struct(s).SetUint16(0, 11)
	C.Struct(s).SetPointer(0, C.Pointer(v))
}
func (s Message) Join() Join { return Join(C.Struct(s).Pointer(0).ToStruct()) }
func (s Message) SetJoin(v Join) {
	C.Struct(s).SetUint16(0, 12)
	C.Struct(s).SetPointer(0, C.Pointer(v))
}

type Message_List C.PointerList

func NewMessage_List(s *C.Segment, sz int32) Message_List {
	return Message_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 1}, sz))
}
func (s Message_List) Len() int                { return C.PointerList(s).Len() }
func (s Message_List) At(i int) Message        { return Message(C.PointerList(s).At(i).ToStruct()) }
func (s Message_List) Set(i int, item Message) { C.PointerList(s).Set(i, C.Pointer(item)) }

type Message_Promise C.Pipeline

func (p *Message_Promise) Get() (Message, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Message(s), err
}

func (p *Message_Promise) Unimplemented() *Message_Promise {
	return (*Message_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *Message_Promise) Abort() *Exception_Promise {
	return (*Exception_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *Message_Promise) Bootstrap() *Bootstrap_Promise {
	return (*Bootstrap_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *Message_Promise) Call() *Call_Promise {
	return (*Call_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *Message_Promise) Return() *Return_Promise {
	return (*Return_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *Message_Promise) Finish() *Finish_Promise {
	return (*Finish_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *Message_Promise) Resolve() *Resolve_Promise {
	return (*Resolve_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *Message_Promise) Release() *Release_Promise {
	return (*Release_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *Message_Promise) Disembargo() *Disembargo_Promise {
	return (*Disembargo_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *Message_Promise) ObsoleteSave() *C.Pipeline {
	return (*C.Pipeline)(p).GetPipeline(0)
}

func (p *Message_Promise) ObsoleteDelete() *C.Pipeline {
	return (*C.Pipeline)(p).GetPipeline(0)
}

func (p *Message_Promise) Provide() *Provide_Promise {
	return (*Provide_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *Message_Promise) Accept() *Accept_Promise {
	return (*Accept_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *Message_Promise) Join() *Join_Promise {
	return (*Join_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

type Bootstrap C.Struct

func NewBootstrap(s *C.Segment) Bootstrap {
	return Bootstrap(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func NewRootBootstrap(s *C.Segment) Bootstrap {
	return Bootstrap(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func AutoNewBootstrap(s *C.Segment) Bootstrap {
	return Bootstrap(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func ReadRootBootstrap(s *C.Segment) Bootstrap { return Bootstrap(s.Root(0).ToStruct()) }

func (s Bootstrap) Segment() *C.Segment {
	return C.Struct(s).Segment()
}

func (s Bootstrap) QuestionId() uint32                { return C.Struct(s).Uint32(0) }
func (s Bootstrap) SetQuestionId(v uint32)            { C.Struct(s).SetUint32(0, v) }
func (s Bootstrap) DeprecatedObjectId() C.Pointer     { return C.Struct(s).Pointer(0) }
func (s Bootstrap) SetDeprecatedObjectId(v C.Pointer) { C.Struct(s).SetPointer(0, v) }

type Bootstrap_List C.PointerList

func NewBootstrap_List(s *C.Segment, sz int32) Bootstrap_List {
	return Bootstrap_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 1}, sz))
}
func (s Bootstrap_List) Len() int                  { return C.PointerList(s).Len() }
func (s Bootstrap_List) At(i int) Bootstrap        { return Bootstrap(C.PointerList(s).At(i).ToStruct()) }
func (s Bootstrap_List) Set(i int, item Bootstrap) { C.PointerList(s).Set(i, C.Pointer(item)) }

type Bootstrap_Promise C.Pipeline

func (p *Bootstrap_Promise) Get() (Bootstrap, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Bootstrap(s), err
}

func (p *Bootstrap_Promise) DeprecatedObjectId() *C.Pipeline {
	return (*C.Pipeline)(p).GetPipeline(0)
}

type Call C.Struct
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

func NewCall(s *C.Segment) Call { return Call(s.NewStruct(C.ObjectSize{DataSize: 24, PointerCount: 3})) }
func NewRootCall(s *C.Segment) Call {
	return Call(s.NewRootStruct(C.ObjectSize{DataSize: 24, PointerCount: 3}))
}
func AutoNewCall(s *C.Segment) Call {
	return Call(s.NewStructAR(C.ObjectSize{DataSize: 24, PointerCount: 3}))
}
func ReadRootCall(s *C.Segment) Call { return Call(s.Root(0).ToStruct()) }

func (s Call) Segment() *C.Segment {
	return C.Struct(s).Segment()
}

func (s Call) QuestionId() uint32                { return C.Struct(s).Uint32(0) }
func (s Call) SetQuestionId(v uint32)            { C.Struct(s).SetUint32(0, v) }
func (s Call) Target() MessageTarget             { return MessageTarget(C.Struct(s).Pointer(0).ToStruct()) }
func (s Call) SetTarget(v MessageTarget)         { C.Struct(s).SetPointer(0, C.Pointer(v)) }
func (s Call) InterfaceId() uint64               { return C.Struct(s).Uint64(8) }
func (s Call) SetInterfaceId(v uint64)           { C.Struct(s).SetUint64(8, v) }
func (s Call) MethodId() uint16                  { return C.Struct(s).Uint16(4) }
func (s Call) SetMethodId(v uint16)              { C.Struct(s).SetUint16(4, v) }
func (s Call) AllowThirdPartyTailCall() bool     { return C.Struct(s).Bit(128) }
func (s Call) SetAllowThirdPartyTailCall(v bool) { C.Struct(s).SetBit(128, v) }
func (s Call) Params() Payload                   { return Payload(C.Struct(s).Pointer(1).ToStruct()) }
func (s Call) SetParams(v Payload)               { C.Struct(s).SetPointer(1, C.Pointer(v)) }
func (s Call) SendResultsTo() Call_sendResultsTo { return Call_sendResultsTo(s) }

func (s Call_sendResultsTo) Segment() *C.Segment {
	return C.Struct(s).Segment()
}

func (s Call_sendResultsTo) Which() Call_sendResultsTo_Which {
	return Call_sendResultsTo_Which(C.Struct(s).Uint16(6))
}

func (s Call_sendResultsTo) SetCaller()            { C.Struct(s).SetUint16(6, 0) }
func (s Call_sendResultsTo) SetYourself()          { C.Struct(s).SetUint16(6, 1) }
func (s Call_sendResultsTo) ThirdParty() C.Pointer { return C.Struct(s).Pointer(2) }
func (s Call_sendResultsTo) SetThirdParty(v C.Pointer) {
	C.Struct(s).SetUint16(6, 2)
	C.Struct(s).SetPointer(2, v)
}

type Call_List C.PointerList

func NewCall_List(s *C.Segment, sz int32) Call_List {
	return Call_List(s.NewCompositeList(C.ObjectSize{DataSize: 24, PointerCount: 3}, sz))
}
func (s Call_List) Len() int             { return C.PointerList(s).Len() }
func (s Call_List) At(i int) Call        { return Call(C.PointerList(s).At(i).ToStruct()) }
func (s Call_List) Set(i int, item Call) { C.PointerList(s).Set(i, C.Pointer(item)) }

type Call_Promise C.Pipeline

func (p *Call_Promise) Get() (Call, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Call(s), err
}

func (p *Call_Promise) Target() *MessageTarget_Promise {
	return (*MessageTarget_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *Call_Promise) Params() *Payload_Promise {
	return (*Payload_Promise)((*C.Pipeline)(p).GetPipeline(1))
}
func (p *Call_Promise) SendResultsTo() *Call_sendResultsTo_Promise {
	return (*Call_sendResultsTo_Promise)(p)
}

type Call_sendResultsTo_Promise C.Pipeline

func (p *Call_sendResultsTo_Promise) Get() (Call_sendResultsTo, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Call_sendResultsTo(s), err
}

func (p *Call_sendResultsTo_Promise) ThirdParty() *C.Pipeline {
	return (*C.Pipeline)(p).GetPipeline(2)
}

type Return C.Struct
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

func NewReturn(s *C.Segment) Return {
	return Return(s.NewStruct(C.ObjectSize{DataSize: 16, PointerCount: 1}))
}
func NewRootReturn(s *C.Segment) Return {
	return Return(s.NewRootStruct(C.ObjectSize{DataSize: 16, PointerCount: 1}))
}
func AutoNewReturn(s *C.Segment) Return {
	return Return(s.NewStructAR(C.ObjectSize{DataSize: 16, PointerCount: 1}))
}
func ReadRootReturn(s *C.Segment) Return { return Return(s.Root(0).ToStruct()) }

func (s Return) Segment() *C.Segment {
	return C.Struct(s).Segment()
}

func (s Return) Which() Return_Which {
	return Return_Which(C.Struct(s).Uint16(6))
}

func (s Return) AnswerId() uint32           { return C.Struct(s).Uint32(0) }
func (s Return) SetAnswerId(v uint32)       { C.Struct(s).SetUint32(0, v) }
func (s Return) ReleaseParamCaps() bool     { return !C.Struct(s).Bit(32) }
func (s Return) SetReleaseParamCaps(v bool) { C.Struct(s).SetBit(32, !v) }
func (s Return) Results() Payload           { return Payload(C.Struct(s).Pointer(0).ToStruct()) }
func (s Return) SetResults(v Payload) {
	C.Struct(s).SetUint16(6, 0)
	C.Struct(s).SetPointer(0, C.Pointer(v))
}
func (s Return) Exception() Exception { return Exception(C.Struct(s).Pointer(0).ToStruct()) }
func (s Return) SetException(v Exception) {
	C.Struct(s).SetUint16(6, 1)
	C.Struct(s).SetPointer(0, C.Pointer(v))
}
func (s Return) SetCanceled()                  { C.Struct(s).SetUint16(6, 2) }
func (s Return) SetResultsSentElsewhere()      { C.Struct(s).SetUint16(6, 3) }
func (s Return) TakeFromOtherQuestion() uint32 { return C.Struct(s).Uint32(8) }
func (s Return) SetTakeFromOtherQuestion(v uint32) {
	C.Struct(s).SetUint16(6, 4)
	C.Struct(s).SetUint32(8, v)
}
func (s Return) AcceptFromThirdParty() C.Pointer { return C.Struct(s).Pointer(0) }
func (s Return) SetAcceptFromThirdParty(v C.Pointer) {
	C.Struct(s).SetUint16(6, 5)
	C.Struct(s).SetPointer(0, v)
}

type Return_List C.PointerList

func NewReturn_List(s *C.Segment, sz int32) Return_List {
	return Return_List(s.NewCompositeList(C.ObjectSize{DataSize: 16, PointerCount: 1}, sz))
}
func (s Return_List) Len() int               { return C.PointerList(s).Len() }
func (s Return_List) At(i int) Return        { return Return(C.PointerList(s).At(i).ToStruct()) }
func (s Return_List) Set(i int, item Return) { C.PointerList(s).Set(i, C.Pointer(item)) }

type Return_Promise C.Pipeline

func (p *Return_Promise) Get() (Return, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Return(s), err
}

func (p *Return_Promise) Results() *Payload_Promise {
	return (*Payload_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *Return_Promise) Exception() *Exception_Promise {
	return (*Exception_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *Return_Promise) AcceptFromThirdParty() *C.Pipeline {
	return (*C.Pipeline)(p).GetPipeline(0)
}

type Finish C.Struct

func NewFinish(s *C.Segment) Finish {
	return Finish(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 0}))
}
func NewRootFinish(s *C.Segment) Finish {
	return Finish(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 0}))
}
func AutoNewFinish(s *C.Segment) Finish {
	return Finish(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 0}))
}
func ReadRootFinish(s *C.Segment) Finish { return Finish(s.Root(0).ToStruct()) }

func (s Finish) Segment() *C.Segment {
	return C.Struct(s).Segment()
}

func (s Finish) QuestionId() uint32          { return C.Struct(s).Uint32(0) }
func (s Finish) SetQuestionId(v uint32)      { C.Struct(s).SetUint32(0, v) }
func (s Finish) ReleaseResultCaps() bool     { return !C.Struct(s).Bit(32) }
func (s Finish) SetReleaseResultCaps(v bool) { C.Struct(s).SetBit(32, !v) }

type Finish_List C.PointerList

func NewFinish_List(s *C.Segment, sz int32) Finish_List {
	return Finish_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 0}, sz))
}
func (s Finish_List) Len() int               { return C.PointerList(s).Len() }
func (s Finish_List) At(i int) Finish        { return Finish(C.PointerList(s).At(i).ToStruct()) }
func (s Finish_List) Set(i int, item Finish) { C.PointerList(s).Set(i, C.Pointer(item)) }

type Finish_Promise C.Pipeline

func (p *Finish_Promise) Get() (Finish, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Finish(s), err
}

type Resolve C.Struct
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

func NewResolve(s *C.Segment) Resolve {
	return Resolve(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func NewRootResolve(s *C.Segment) Resolve {
	return Resolve(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func AutoNewResolve(s *C.Segment) Resolve {
	return Resolve(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func ReadRootResolve(s *C.Segment) Resolve { return Resolve(s.Root(0).ToStruct()) }

func (s Resolve) Segment() *C.Segment {
	return C.Struct(s).Segment()
}

func (s Resolve) Which() Resolve_Which {
	return Resolve_Which(C.Struct(s).Uint16(4))
}

func (s Resolve) PromiseId() uint32     { return C.Struct(s).Uint32(0) }
func (s Resolve) SetPromiseId(v uint32) { C.Struct(s).SetUint32(0, v) }
func (s Resolve) Cap() CapDescriptor    { return CapDescriptor(C.Struct(s).Pointer(0).ToStruct()) }
func (s Resolve) SetCap(v CapDescriptor) {
	C.Struct(s).SetUint16(4, 0)
	C.Struct(s).SetPointer(0, C.Pointer(v))
}
func (s Resolve) Exception() Exception { return Exception(C.Struct(s).Pointer(0).ToStruct()) }
func (s Resolve) SetException(v Exception) {
	C.Struct(s).SetUint16(4, 1)
	C.Struct(s).SetPointer(0, C.Pointer(v))
}

type Resolve_List C.PointerList

func NewResolve_List(s *C.Segment, sz int32) Resolve_List {
	return Resolve_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 1}, sz))
}
func (s Resolve_List) Len() int                { return C.PointerList(s).Len() }
func (s Resolve_List) At(i int) Resolve        { return Resolve(C.PointerList(s).At(i).ToStruct()) }
func (s Resolve_List) Set(i int, item Resolve) { C.PointerList(s).Set(i, C.Pointer(item)) }

type Resolve_Promise C.Pipeline

func (p *Resolve_Promise) Get() (Resolve, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Resolve(s), err
}

func (p *Resolve_Promise) Cap() *CapDescriptor_Promise {
	return (*CapDescriptor_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *Resolve_Promise) Exception() *Exception_Promise {
	return (*Exception_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

type Release C.Struct

func NewRelease(s *C.Segment) Release {
	return Release(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 0}))
}
func NewRootRelease(s *C.Segment) Release {
	return Release(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 0}))
}
func AutoNewRelease(s *C.Segment) Release {
	return Release(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 0}))
}
func ReadRootRelease(s *C.Segment) Release { return Release(s.Root(0).ToStruct()) }

func (s Release) Segment() *C.Segment {
	return C.Struct(s).Segment()
}

func (s Release) Id() uint32                 { return C.Struct(s).Uint32(0) }
func (s Release) SetId(v uint32)             { C.Struct(s).SetUint32(0, v) }
func (s Release) ReferenceCount() uint32     { return C.Struct(s).Uint32(4) }
func (s Release) SetReferenceCount(v uint32) { C.Struct(s).SetUint32(4, v) }

type Release_List C.PointerList

func NewRelease_List(s *C.Segment, sz int32) Release_List {
	return Release_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 0}, sz))
}
func (s Release_List) Len() int                { return C.PointerList(s).Len() }
func (s Release_List) At(i int) Release        { return Release(C.PointerList(s).At(i).ToStruct()) }
func (s Release_List) Set(i int, item Release) { C.PointerList(s).Set(i, C.Pointer(item)) }

type Release_Promise C.Pipeline

func (p *Release_Promise) Get() (Release, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Release(s), err
}

type Disembargo C.Struct
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

func NewDisembargo(s *C.Segment) Disembargo {
	return Disembargo(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func NewRootDisembargo(s *C.Segment) Disembargo {
	return Disembargo(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func AutoNewDisembargo(s *C.Segment) Disembargo {
	return Disembargo(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func ReadRootDisembargo(s *C.Segment) Disembargo { return Disembargo(s.Root(0).ToStruct()) }

func (s Disembargo) Segment() *C.Segment {
	return C.Struct(s).Segment()
}

func (s Disembargo) Target() MessageTarget       { return MessageTarget(C.Struct(s).Pointer(0).ToStruct()) }
func (s Disembargo) SetTarget(v MessageTarget)   { C.Struct(s).SetPointer(0, C.Pointer(v)) }
func (s Disembargo) Context() Disembargo_context { return Disembargo_context(s) }

func (s Disembargo_context) Segment() *C.Segment {
	return C.Struct(s).Segment()
}

func (s Disembargo_context) Which() Disembargo_context_Which {
	return Disembargo_context_Which(C.Struct(s).Uint16(4))
}

func (s Disembargo_context) SenderLoopback() uint32 { return C.Struct(s).Uint32(0) }
func (s Disembargo_context) SetSenderLoopback(v uint32) {
	C.Struct(s).SetUint16(4, 0)
	C.Struct(s).SetUint32(0, v)
}
func (s Disembargo_context) ReceiverLoopback() uint32 { return C.Struct(s).Uint32(0) }
func (s Disembargo_context) SetReceiverLoopback(v uint32) {
	C.Struct(s).SetUint16(4, 1)
	C.Struct(s).SetUint32(0, v)
}
func (s Disembargo_context) SetAccept()      { C.Struct(s).SetUint16(4, 2) }
func (s Disembargo_context) Provide() uint32 { return C.Struct(s).Uint32(0) }
func (s Disembargo_context) SetProvide(v uint32) {
	C.Struct(s).SetUint16(4, 3)
	C.Struct(s).SetUint32(0, v)
}

type Disembargo_List C.PointerList

func NewDisembargo_List(s *C.Segment, sz int32) Disembargo_List {
	return Disembargo_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 1}, sz))
}
func (s Disembargo_List) Len() int                   { return C.PointerList(s).Len() }
func (s Disembargo_List) At(i int) Disembargo        { return Disembargo(C.PointerList(s).At(i).ToStruct()) }
func (s Disembargo_List) Set(i int, item Disembargo) { C.PointerList(s).Set(i, C.Pointer(item)) }

type Disembargo_Promise C.Pipeline

func (p *Disembargo_Promise) Get() (Disembargo, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Disembargo(s), err
}

func (p *Disembargo_Promise) Target() *MessageTarget_Promise {
	return (*MessageTarget_Promise)((*C.Pipeline)(p).GetPipeline(0))
}
func (p *Disembargo_Promise) Context() *Disembargo_context_Promise {
	return (*Disembargo_context_Promise)(p)
}

type Disembargo_context_Promise C.Pipeline

func (p *Disembargo_context_Promise) Get() (Disembargo_context, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Disembargo_context(s), err
}

type Provide C.Struct

func NewProvide(s *C.Segment) Provide {
	return Provide(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 2}))
}
func NewRootProvide(s *C.Segment) Provide {
	return Provide(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 2}))
}
func AutoNewProvide(s *C.Segment) Provide {
	return Provide(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 2}))
}
func ReadRootProvide(s *C.Segment) Provide { return Provide(s.Root(0).ToStruct()) }

func (s Provide) Segment() *C.Segment {
	return C.Struct(s).Segment()
}

func (s Provide) QuestionId() uint32        { return C.Struct(s).Uint32(0) }
func (s Provide) SetQuestionId(v uint32)    { C.Struct(s).SetUint32(0, v) }
func (s Provide) Target() MessageTarget     { return MessageTarget(C.Struct(s).Pointer(0).ToStruct()) }
func (s Provide) SetTarget(v MessageTarget) { C.Struct(s).SetPointer(0, C.Pointer(v)) }
func (s Provide) Recipient() C.Pointer      { return C.Struct(s).Pointer(1) }
func (s Provide) SetRecipient(v C.Pointer)  { C.Struct(s).SetPointer(1, v) }

type Provide_List C.PointerList

func NewProvide_List(s *C.Segment, sz int32) Provide_List {
	return Provide_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 2}, sz))
}
func (s Provide_List) Len() int                { return C.PointerList(s).Len() }
func (s Provide_List) At(i int) Provide        { return Provide(C.PointerList(s).At(i).ToStruct()) }
func (s Provide_List) Set(i int, item Provide) { C.PointerList(s).Set(i, C.Pointer(item)) }

type Provide_Promise C.Pipeline

func (p *Provide_Promise) Get() (Provide, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Provide(s), err
}

func (p *Provide_Promise) Target() *MessageTarget_Promise {
	return (*MessageTarget_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *Provide_Promise) Recipient() *C.Pipeline {
	return (*C.Pipeline)(p).GetPipeline(1)
}

type Accept C.Struct

func NewAccept(s *C.Segment) Accept {
	return Accept(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func NewRootAccept(s *C.Segment) Accept {
	return Accept(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func AutoNewAccept(s *C.Segment) Accept {
	return Accept(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func ReadRootAccept(s *C.Segment) Accept { return Accept(s.Root(0).ToStruct()) }

func (s Accept) Segment() *C.Segment {
	return C.Struct(s).Segment()
}

func (s Accept) QuestionId() uint32       { return C.Struct(s).Uint32(0) }
func (s Accept) SetQuestionId(v uint32)   { C.Struct(s).SetUint32(0, v) }
func (s Accept) Provision() C.Pointer     { return C.Struct(s).Pointer(0) }
func (s Accept) SetProvision(v C.Pointer) { C.Struct(s).SetPointer(0, v) }
func (s Accept) Embargo() bool            { return C.Struct(s).Bit(32) }
func (s Accept) SetEmbargo(v bool)        { C.Struct(s).SetBit(32, v) }

type Accept_List C.PointerList

func NewAccept_List(s *C.Segment, sz int32) Accept_List {
	return Accept_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 1}, sz))
}
func (s Accept_List) Len() int               { return C.PointerList(s).Len() }
func (s Accept_List) At(i int) Accept        { return Accept(C.PointerList(s).At(i).ToStruct()) }
func (s Accept_List) Set(i int, item Accept) { C.PointerList(s).Set(i, C.Pointer(item)) }

type Accept_Promise C.Pipeline

func (p *Accept_Promise) Get() (Accept, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Accept(s), err
}

func (p *Accept_Promise) Provision() *C.Pipeline {
	return (*C.Pipeline)(p).GetPipeline(0)
}

type Join C.Struct

func NewJoin(s *C.Segment) Join { return Join(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 2})) }
func NewRootJoin(s *C.Segment) Join {
	return Join(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 2}))
}
func AutoNewJoin(s *C.Segment) Join {
	return Join(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 2}))
}
func ReadRootJoin(s *C.Segment) Join { return Join(s.Root(0).ToStruct()) }

func (s Join) Segment() *C.Segment {
	return C.Struct(s).Segment()
}

func (s Join) QuestionId() uint32        { return C.Struct(s).Uint32(0) }
func (s Join) SetQuestionId(v uint32)    { C.Struct(s).SetUint32(0, v) }
func (s Join) Target() MessageTarget     { return MessageTarget(C.Struct(s).Pointer(0).ToStruct()) }
func (s Join) SetTarget(v MessageTarget) { C.Struct(s).SetPointer(0, C.Pointer(v)) }
func (s Join) KeyPart() C.Pointer        { return C.Struct(s).Pointer(1) }
func (s Join) SetKeyPart(v C.Pointer)    { C.Struct(s).SetPointer(1, v) }

type Join_List C.PointerList

func NewJoin_List(s *C.Segment, sz int32) Join_List {
	return Join_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 2}, sz))
}
func (s Join_List) Len() int             { return C.PointerList(s).Len() }
func (s Join_List) At(i int) Join        { return Join(C.PointerList(s).At(i).ToStruct()) }
func (s Join_List) Set(i int, item Join) { C.PointerList(s).Set(i, C.Pointer(item)) }

type Join_Promise C.Pipeline

func (p *Join_Promise) Get() (Join, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Join(s), err
}

func (p *Join_Promise) Target() *MessageTarget_Promise {
	return (*MessageTarget_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *Join_Promise) KeyPart() *C.Pipeline {
	return (*C.Pipeline)(p).GetPipeline(1)
}

type MessageTarget C.Struct
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

func NewMessageTarget(s *C.Segment) MessageTarget {
	return MessageTarget(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func NewRootMessageTarget(s *C.Segment) MessageTarget {
	return MessageTarget(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func AutoNewMessageTarget(s *C.Segment) MessageTarget {
	return MessageTarget(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func ReadRootMessageTarget(s *C.Segment) MessageTarget { return MessageTarget(s.Root(0).ToStruct()) }

func (s MessageTarget) Segment() *C.Segment {
	return C.Struct(s).Segment()
}

func (s MessageTarget) Which() MessageTarget_Which {
	return MessageTarget_Which(C.Struct(s).Uint16(4))
}

func (s MessageTarget) ImportedCap() uint32 { return C.Struct(s).Uint32(0) }
func (s MessageTarget) SetImportedCap(v uint32) {
	C.Struct(s).SetUint16(4, 0)
	C.Struct(s).SetUint32(0, v)
}
func (s MessageTarget) PromisedAnswer() PromisedAnswer {
	return PromisedAnswer(C.Struct(s).Pointer(0).ToStruct())
}
func (s MessageTarget) SetPromisedAnswer(v PromisedAnswer) {
	C.Struct(s).SetUint16(4, 1)
	C.Struct(s).SetPointer(0, C.Pointer(v))
}

type MessageTarget_List C.PointerList

func NewMessageTarget_List(s *C.Segment, sz int32) MessageTarget_List {
	return MessageTarget_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 1}, sz))
}
func (s MessageTarget_List) Len() int { return C.PointerList(s).Len() }
func (s MessageTarget_List) At(i int) MessageTarget {
	return MessageTarget(C.PointerList(s).At(i).ToStruct())
}
func (s MessageTarget_List) Set(i int, item MessageTarget) { C.PointerList(s).Set(i, C.Pointer(item)) }

type MessageTarget_Promise C.Pipeline

func (p *MessageTarget_Promise) Get() (MessageTarget, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return MessageTarget(s), err
}

func (p *MessageTarget_Promise) PromisedAnswer() *PromisedAnswer_Promise {
	return (*PromisedAnswer_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

type Payload C.Struct

func NewPayload(s *C.Segment) Payload {
	return Payload(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 2}))
}
func NewRootPayload(s *C.Segment) Payload {
	return Payload(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 2}))
}
func AutoNewPayload(s *C.Segment) Payload {
	return Payload(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 2}))
}
func ReadRootPayload(s *C.Segment) Payload { return Payload(s.Root(0).ToStruct()) }

func (s Payload) Segment() *C.Segment {
	return C.Struct(s).Segment()
}

func (s Payload) Content() C.Pointer               { return C.Struct(s).Pointer(0) }
func (s Payload) SetContent(v C.Pointer)           { C.Struct(s).SetPointer(0, v) }
func (s Payload) CapTable() CapDescriptor_List     { return CapDescriptor_List(C.Struct(s).Pointer(1)) }
func (s Payload) SetCapTable(v CapDescriptor_List) { C.Struct(s).SetPointer(1, C.Pointer(v)) }

type Payload_List C.PointerList

func NewPayload_List(s *C.Segment, sz int32) Payload_List {
	return Payload_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 2}, sz))
}
func (s Payload_List) Len() int                { return C.PointerList(s).Len() }
func (s Payload_List) At(i int) Payload        { return Payload(C.PointerList(s).At(i).ToStruct()) }
func (s Payload_List) Set(i int, item Payload) { C.PointerList(s).Set(i, C.Pointer(item)) }

type Payload_Promise C.Pipeline

func (p *Payload_Promise) Get() (Payload, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Payload(s), err
}

func (p *Payload_Promise) Content() *C.Pipeline {
	return (*C.Pipeline)(p).GetPipeline(0)
}

type CapDescriptor C.Struct
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

func NewCapDescriptor(s *C.Segment) CapDescriptor {
	return CapDescriptor(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func NewRootCapDescriptor(s *C.Segment) CapDescriptor {
	return CapDescriptor(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func AutoNewCapDescriptor(s *C.Segment) CapDescriptor {
	return CapDescriptor(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func ReadRootCapDescriptor(s *C.Segment) CapDescriptor { return CapDescriptor(s.Root(0).ToStruct()) }

func (s CapDescriptor) Segment() *C.Segment {
	return C.Struct(s).Segment()
}

func (s CapDescriptor) Which() CapDescriptor_Which {
	return CapDescriptor_Which(C.Struct(s).Uint16(0))
}

func (s CapDescriptor) SetNone()             { C.Struct(s).SetUint16(0, 0) }
func (s CapDescriptor) SenderHosted() uint32 { return C.Struct(s).Uint32(4) }
func (s CapDescriptor) SetSenderHosted(v uint32) {
	C.Struct(s).SetUint16(0, 1)
	C.Struct(s).SetUint32(4, v)
}
func (s CapDescriptor) SenderPromise() uint32 { return C.Struct(s).Uint32(4) }
func (s CapDescriptor) SetSenderPromise(v uint32) {
	C.Struct(s).SetUint16(0, 2)
	C.Struct(s).SetUint32(4, v)
}
func (s CapDescriptor) ReceiverHosted() uint32 { return C.Struct(s).Uint32(4) }
func (s CapDescriptor) SetReceiverHosted(v uint32) {
	C.Struct(s).SetUint16(0, 3)
	C.Struct(s).SetUint32(4, v)
}
func (s CapDescriptor) ReceiverAnswer() PromisedAnswer {
	return PromisedAnswer(C.Struct(s).Pointer(0).ToStruct())
}
func (s CapDescriptor) SetReceiverAnswer(v PromisedAnswer) {
	C.Struct(s).SetUint16(0, 4)
	C.Struct(s).SetPointer(0, C.Pointer(v))
}
func (s CapDescriptor) ThirdPartyHosted() ThirdPartyCapDescriptor {
	return ThirdPartyCapDescriptor(C.Struct(s).Pointer(0).ToStruct())
}
func (s CapDescriptor) SetThirdPartyHosted(v ThirdPartyCapDescriptor) {
	C.Struct(s).SetUint16(0, 5)
	C.Struct(s).SetPointer(0, C.Pointer(v))
}

type CapDescriptor_List C.PointerList

func NewCapDescriptor_List(s *C.Segment, sz int32) CapDescriptor_List {
	return CapDescriptor_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 1}, sz))
}
func (s CapDescriptor_List) Len() int { return C.PointerList(s).Len() }
func (s CapDescriptor_List) At(i int) CapDescriptor {
	return CapDescriptor(C.PointerList(s).At(i).ToStruct())
}
func (s CapDescriptor_List) Set(i int, item CapDescriptor) { C.PointerList(s).Set(i, C.Pointer(item)) }

type CapDescriptor_Promise C.Pipeline

func (p *CapDescriptor_Promise) Get() (CapDescriptor, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return CapDescriptor(s), err
}

func (p *CapDescriptor_Promise) ReceiverAnswer() *PromisedAnswer_Promise {
	return (*PromisedAnswer_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *CapDescriptor_Promise) ThirdPartyHosted() *ThirdPartyCapDescriptor_Promise {
	return (*ThirdPartyCapDescriptor_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

type PromisedAnswer C.Struct

func NewPromisedAnswer(s *C.Segment) PromisedAnswer {
	return PromisedAnswer(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func NewRootPromisedAnswer(s *C.Segment) PromisedAnswer {
	return PromisedAnswer(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func AutoNewPromisedAnswer(s *C.Segment) PromisedAnswer {
	return PromisedAnswer(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func ReadRootPromisedAnswer(s *C.Segment) PromisedAnswer { return PromisedAnswer(s.Root(0).ToStruct()) }

func (s PromisedAnswer) Segment() *C.Segment {
	return C.Struct(s).Segment()
}

func (s PromisedAnswer) QuestionId() uint32     { return C.Struct(s).Uint32(0) }
func (s PromisedAnswer) SetQuestionId(v uint32) { C.Struct(s).SetUint32(0, v) }
func (s PromisedAnswer) Transform() PromisedAnswer_Op_List {
	return PromisedAnswer_Op_List(C.Struct(s).Pointer(0))
}
func (s PromisedAnswer) SetTransform(v PromisedAnswer_Op_List) {
	C.Struct(s).SetPointer(0, C.Pointer(v))
}

type PromisedAnswer_List C.PointerList

func NewPromisedAnswer_List(s *C.Segment, sz int32) PromisedAnswer_List {
	return PromisedAnswer_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 1}, sz))
}
func (s PromisedAnswer_List) Len() int { return C.PointerList(s).Len() }
func (s PromisedAnswer_List) At(i int) PromisedAnswer {
	return PromisedAnswer(C.PointerList(s).At(i).ToStruct())
}
func (s PromisedAnswer_List) Set(i int, item PromisedAnswer) { C.PointerList(s).Set(i, C.Pointer(item)) }

type PromisedAnswer_Promise C.Pipeline

func (p *PromisedAnswer_Promise) Get() (PromisedAnswer, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return PromisedAnswer(s), err
}

type PromisedAnswer_Op C.Struct
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

func NewPromisedAnswer_Op(s *C.Segment) PromisedAnswer_Op {
	return PromisedAnswer_Op(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 0}))
}
func NewRootPromisedAnswer_Op(s *C.Segment) PromisedAnswer_Op {
	return PromisedAnswer_Op(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 0}))
}
func AutoNewPromisedAnswer_Op(s *C.Segment) PromisedAnswer_Op {
	return PromisedAnswer_Op(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 0}))
}
func ReadRootPromisedAnswer_Op(s *C.Segment) PromisedAnswer_Op {
	return PromisedAnswer_Op(s.Root(0).ToStruct())
}

func (s PromisedAnswer_Op) Segment() *C.Segment {
	return C.Struct(s).Segment()
}

func (s PromisedAnswer_Op) Which() PromisedAnswer_Op_Which {
	return PromisedAnswer_Op_Which(C.Struct(s).Uint16(0))
}

func (s PromisedAnswer_Op) SetNoop()                { C.Struct(s).SetUint16(0, 0) }
func (s PromisedAnswer_Op) GetPointerField() uint16 { return C.Struct(s).Uint16(2) }
func (s PromisedAnswer_Op) SetGetPointerField(v uint16) {
	C.Struct(s).SetUint16(0, 1)
	C.Struct(s).SetUint16(2, v)
}

type PromisedAnswer_Op_List C.PointerList

func NewPromisedAnswer_Op_List(s *C.Segment, sz int32) PromisedAnswer_Op_List {
	return PromisedAnswer_Op_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 0}, sz))
}
func (s PromisedAnswer_Op_List) Len() int { return C.PointerList(s).Len() }
func (s PromisedAnswer_Op_List) At(i int) PromisedAnswer_Op {
	return PromisedAnswer_Op(C.PointerList(s).At(i).ToStruct())
}
func (s PromisedAnswer_Op_List) Set(i int, item PromisedAnswer_Op) {
	C.PointerList(s).Set(i, C.Pointer(item))
}

type PromisedAnswer_Op_Promise C.Pipeline

func (p *PromisedAnswer_Op_Promise) Get() (PromisedAnswer_Op, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return PromisedAnswer_Op(s), err
}

type ThirdPartyCapDescriptor C.Struct

func NewThirdPartyCapDescriptor(s *C.Segment) ThirdPartyCapDescriptor {
	return ThirdPartyCapDescriptor(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func NewRootThirdPartyCapDescriptor(s *C.Segment) ThirdPartyCapDescriptor {
	return ThirdPartyCapDescriptor(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func AutoNewThirdPartyCapDescriptor(s *C.Segment) ThirdPartyCapDescriptor {
	return ThirdPartyCapDescriptor(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func ReadRootThirdPartyCapDescriptor(s *C.Segment) ThirdPartyCapDescriptor {
	return ThirdPartyCapDescriptor(s.Root(0).ToStruct())
}

func (s ThirdPartyCapDescriptor) Segment() *C.Segment {
	return C.Struct(s).Segment()
}

func (s ThirdPartyCapDescriptor) Id() C.Pointer      { return C.Struct(s).Pointer(0) }
func (s ThirdPartyCapDescriptor) SetId(v C.Pointer)  { C.Struct(s).SetPointer(0, v) }
func (s ThirdPartyCapDescriptor) VineId() uint32     { return C.Struct(s).Uint32(0) }
func (s ThirdPartyCapDescriptor) SetVineId(v uint32) { C.Struct(s).SetUint32(0, v) }

type ThirdPartyCapDescriptor_List C.PointerList

func NewThirdPartyCapDescriptor_List(s *C.Segment, sz int32) ThirdPartyCapDescriptor_List {
	return ThirdPartyCapDescriptor_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 1}, sz))
}
func (s ThirdPartyCapDescriptor_List) Len() int { return C.PointerList(s).Len() }
func (s ThirdPartyCapDescriptor_List) At(i int) ThirdPartyCapDescriptor {
	return ThirdPartyCapDescriptor(C.PointerList(s).At(i).ToStruct())
}
func (s ThirdPartyCapDescriptor_List) Set(i int, item ThirdPartyCapDescriptor) {
	C.PointerList(s).Set(i, C.Pointer(item))
}

type ThirdPartyCapDescriptor_Promise C.Pipeline

func (p *ThirdPartyCapDescriptor_Promise) Get() (ThirdPartyCapDescriptor, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return ThirdPartyCapDescriptor(s), err
}

func (p *ThirdPartyCapDescriptor_Promise) Id() *C.Pipeline {
	return (*C.Pipeline)(p).GetPipeline(0)
}

type Exception C.Struct

func NewException(s *C.Segment) Exception {
	return Exception(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func NewRootException(s *C.Segment) Exception {
	return Exception(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func AutoNewException(s *C.Segment) Exception {
	return Exception(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func ReadRootException(s *C.Segment) Exception { return Exception(s.Root(0).ToStruct()) }

func (s Exception) Segment() *C.Segment {
	return C.Struct(s).Segment()
}

func (s Exception) Reason() string                   { return C.Struct(s).Pointer(0).ToText() }
func (s Exception) SetReason(v string)               { C.Struct(s).SetPointer(0, s.Segment().NewText(v)) }
func (s Exception) Type() Exception_Type             { return Exception_Type(C.Struct(s).Uint16(4)) }
func (s Exception) SetType(v Exception_Type)         { C.Struct(s).SetUint16(4, uint16(v)) }
func (s Exception) ObsoleteIsCallersFault() bool     { return C.Struct(s).Bit(0) }
func (s Exception) SetObsoleteIsCallersFault(v bool) { C.Struct(s).SetBit(0, v) }
func (s Exception) ObsoleteDurability() uint16       { return C.Struct(s).Uint16(2) }
func (s Exception) SetObsoleteDurability(v uint16)   { C.Struct(s).SetUint16(2, v) }

type Exception_List C.PointerList

func NewException_List(s *C.Segment, sz int32) Exception_List {
	return Exception_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 1}, sz))
}
func (s Exception_List) Len() int                  { return C.PointerList(s).Len() }
func (s Exception_List) At(i int) Exception        { return Exception(C.PointerList(s).At(i).ToStruct()) }
func (s Exception_List) Set(i int, item Exception) { C.PointerList(s).Set(i, C.Pointer(item)) }

type Exception_Promise C.Pipeline

func (p *Exception_Promise) Get() (Exception, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Exception(s), err
}

type Exception_Type uint16

const (
	Exception_Type_failed        Exception_Type = 0
	Exception_Type_overloaded    Exception_Type = 1
	Exception_Type_disconnected  Exception_Type = 2
	Exception_Type_unimplemented Exception_Type = 3
)

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

type Exception_Type_List C.UInt16List

func NewException_Type_List(s *C.Segment, sz int32) Exception_Type_List {
	return Exception_Type_List(s.NewUInt16List(sz))
}
func (s Exception_Type_List) Len() int                { return C.UInt16List(s).Len() }
func (s Exception_Type_List) At(i int) Exception_Type { return Exception_Type(C.UInt16List(s).At(i)) }
