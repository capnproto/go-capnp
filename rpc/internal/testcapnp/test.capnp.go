package testcapnp

// AUTO GENERATED - DO NOT EDIT

import (
	context "golang.org/x/net/context"
	C "zombiezen.com/go/capnproto"
)

type Handle struct{ c C.Client }

func NewHandle(c C.Client) Handle { return Handle{c} }

func (c Handle) GenericClient() C.Client { return c.c }

func (c Handle) IsNull() bool { return c.c == nil }

type Handle_Server interface {
}

func Handle_ServerToClient(s Handle_Server) Handle {
	c, _ := s.(C.Closer)
	return NewHandle(C.NewServer(Handle_Methods(nil, s), c))
}

func Handle_Methods(methods []C.ServerMethod, server Handle_Server) []C.ServerMethod {
	if cap(methods) == 0 {
		methods = make([]C.ServerMethod, 0, 0)
	}

	return methods
}

type HandleFactory struct{ c C.Client }

func NewHandleFactory(c C.Client) HandleFactory { return HandleFactory{c} }

func (c HandleFactory) GenericClient() C.Client { return c.c }

func (c HandleFactory) IsNull() bool { return c.c == nil }

func (c HandleFactory) NewHandle(ctx context.Context, params func(HandleFactory_newHandle_Params), opts ...C.CallOption) *HandleFactory_newHandle_Results_Promise {
	if c.c == nil {
		return (*HandleFactory_newHandle_Results_Promise)(C.NewPipeline(C.ErrorAnswer(C.ErrNullClient)))
	}
	return (*HandleFactory_newHandle_Results_Promise)(C.NewPipeline(c.c.Call(&C.Call{
		Ctx: ctx,
		Method: C.Method{

			InterfaceID:   0x8491a7fe75fe0bce,
			MethodID:      0,
			InterfaceName: "test.capnp:HandleFactory",
			MethodName:    "newHandle",
		},
		ParamsSize: C.ObjectSize{DataSize: 0, PointerCount: 0},
		ParamsFunc: func(s C.Struct) { params(HandleFactory_newHandle_Params(s)) },
		Options:    C.NewCallOptions(opts),
	})))
}

type HandleFactory_Server interface {
	NewHandle(HandleFactory_newHandle) error
}

func HandleFactory_ServerToClient(s HandleFactory_Server) HandleFactory {
	c, _ := s.(C.Closer)
	return NewHandleFactory(C.NewServer(HandleFactory_Methods(nil, s), c))
}

func HandleFactory_Methods(methods []C.ServerMethod, server HandleFactory_Server) []C.ServerMethod {
	if cap(methods) == 0 {
		methods = make([]C.ServerMethod, 0, 1)
	}

	methods = append(methods, C.ServerMethod{
		Method: C.Method{

			InterfaceID:   0x8491a7fe75fe0bce,
			MethodID:      0,
			InterfaceName: "test.capnp:HandleFactory",
			MethodName:    "newHandle",
		},
		Impl: func(c context.Context, opts C.CallOptions, p, r C.Struct) error {
			call := HandleFactory_newHandle{c, opts, HandleFactory_newHandle_Params(p), HandleFactory_newHandle_Results(r)}
			return server.NewHandle(call)
		},
		ResultsSize: C.ObjectSize{DataSize: 0, PointerCount: 1},
	})

	return methods
}

// HandleFactory_newHandle holds the arguments for a server call to HandleFactory.newHandle.
type HandleFactory_newHandle struct {
	Ctx     context.Context
	Options C.CallOptions
	Params  HandleFactory_newHandle_Params
	Results HandleFactory_newHandle_Results
}

type HandleFactory_newHandle_Params C.Struct

func NewHandleFactory_newHandle_Params(s *C.Segment) HandleFactory_newHandle_Params {
	return HandleFactory_newHandle_Params(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 0}))
}
func NewRootHandleFactory_newHandle_Params(s *C.Segment) HandleFactory_newHandle_Params {
	return HandleFactory_newHandle_Params(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 0}))
}
func AutoNewHandleFactory_newHandle_Params(s *C.Segment) HandleFactory_newHandle_Params {
	return HandleFactory_newHandle_Params(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 0}))
}
func ReadRootHandleFactory_newHandle_Params(s *C.Segment) HandleFactory_newHandle_Params {
	return HandleFactory_newHandle_Params(s.Root(0).ToStruct())
}

// capnp.JSON_enabled == false so we stub MarshalJSON().
func (s HandleFactory_newHandle_Params) MarshalJSON() (bs []byte, err error) { return }

type HandleFactory_newHandle_Params_List C.PointerList

func NewHandleFactory_newHandle_Params_List(s *C.Segment, sz int) HandleFactory_newHandle_Params_List {
	return HandleFactory_newHandle_Params_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 0}, sz))
}
func (s HandleFactory_newHandle_Params_List) Len() int { return C.PointerList(s).Len() }
func (s HandleFactory_newHandle_Params_List) At(i int) HandleFactory_newHandle_Params {
	return HandleFactory_newHandle_Params(C.PointerList(s).At(i).ToStruct())
}
func (s HandleFactory_newHandle_Params_List) Set(i int, item HandleFactory_newHandle_Params) {
	C.PointerList(s).Set(i, C.Object(item))
}

type HandleFactory_newHandle_Params_Promise C.Pipeline

func (p *HandleFactory_newHandle_Params_Promise) Get() (HandleFactory_newHandle_Params, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return HandleFactory_newHandle_Params(s), err
}

type HandleFactory_newHandle_Results C.Struct

func NewHandleFactory_newHandle_Results(s *C.Segment) HandleFactory_newHandle_Results {
	return HandleFactory_newHandle_Results(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func NewRootHandleFactory_newHandle_Results(s *C.Segment) HandleFactory_newHandle_Results {
	return HandleFactory_newHandle_Results(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewHandleFactory_newHandle_Results(s *C.Segment) HandleFactory_newHandle_Results {
	return HandleFactory_newHandle_Results(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootHandleFactory_newHandle_Results(s *C.Segment) HandleFactory_newHandle_Results {
	return HandleFactory_newHandle_Results(s.Root(0).ToStruct())
}
func (s HandleFactory_newHandle_Results) Handle() Handle {
	return NewHandle(C.Struct(s).GetObject(0).ToInterface().Client())
}
func (s HandleFactory_newHandle_Results) SetHandle(v Handle) {
	if s.Segment == nil {
		return
	}
	ci := s.Segment.Message.AddCap(v.GenericClient())
	C.Struct(s).SetObject(0, C.Object(s.Segment.NewInterface(ci)))
}

// capnp.JSON_enabled == false so we stub MarshalJSON().
func (s HandleFactory_newHandle_Results) MarshalJSON() (bs []byte, err error) { return }

type HandleFactory_newHandle_Results_List C.PointerList

func NewHandleFactory_newHandle_Results_List(s *C.Segment, sz int) HandleFactory_newHandle_Results_List {
	return HandleFactory_newHandle_Results_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s HandleFactory_newHandle_Results_List) Len() int { return C.PointerList(s).Len() }
func (s HandleFactory_newHandle_Results_List) At(i int) HandleFactory_newHandle_Results {
	return HandleFactory_newHandle_Results(C.PointerList(s).At(i).ToStruct())
}
func (s HandleFactory_newHandle_Results_List) Set(i int, item HandleFactory_newHandle_Results) {
	C.PointerList(s).Set(i, C.Object(item))
}

type HandleFactory_newHandle_Results_Promise C.Pipeline

func (p *HandleFactory_newHandle_Results_Promise) Get() (HandleFactory_newHandle_Results, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return HandleFactory_newHandle_Results(s), err
}

func (p *HandleFactory_newHandle_Results_Promise) Handle() Handle {
	return NewHandle((*C.Pipeline)(p).GetPipeline(0).Client())
}

type Hanger struct{ c C.Client }

func NewHanger(c C.Client) Hanger { return Hanger{c} }

func (c Hanger) GenericClient() C.Client { return c.c }

func (c Hanger) IsNull() bool { return c.c == nil }

func (c Hanger) Hang(ctx context.Context, params func(Hanger_hang_Params), opts ...C.CallOption) *Hanger_hang_Results_Promise {
	if c.c == nil {
		return (*Hanger_hang_Results_Promise)(C.NewPipeline(C.ErrorAnswer(C.ErrNullClient)))
	}
	return (*Hanger_hang_Results_Promise)(C.NewPipeline(c.c.Call(&C.Call{
		Ctx: ctx,
		Method: C.Method{

			InterfaceID:   0x8ae08044aae8a26e,
			MethodID:      0,
			InterfaceName: "test.capnp:Hanger",
			MethodName:    "hang",
		},
		ParamsSize: C.ObjectSize{DataSize: 0, PointerCount: 0},
		ParamsFunc: func(s C.Struct) { params(Hanger_hang_Params(s)) },
		Options:    C.NewCallOptions(opts),
	})))
}

type Hanger_Server interface {
	Hang(Hanger_hang) error
}

func Hanger_ServerToClient(s Hanger_Server) Hanger {
	c, _ := s.(C.Closer)
	return NewHanger(C.NewServer(Hanger_Methods(nil, s), c))
}

func Hanger_Methods(methods []C.ServerMethod, server Hanger_Server) []C.ServerMethod {
	if cap(methods) == 0 {
		methods = make([]C.ServerMethod, 0, 1)
	}

	methods = append(methods, C.ServerMethod{
		Method: C.Method{

			InterfaceID:   0x8ae08044aae8a26e,
			MethodID:      0,
			InterfaceName: "test.capnp:Hanger",
			MethodName:    "hang",
		},
		Impl: func(c context.Context, opts C.CallOptions, p, r C.Struct) error {
			call := Hanger_hang{c, opts, Hanger_hang_Params(p), Hanger_hang_Results(r)}
			return server.Hang(call)
		},
		ResultsSize: C.ObjectSize{DataSize: 0, PointerCount: 0},
	})

	return methods
}

// Hanger_hang holds the arguments for a server call to Hanger.hang.
type Hanger_hang struct {
	Ctx     context.Context
	Options C.CallOptions
	Params  Hanger_hang_Params
	Results Hanger_hang_Results
}

type Hanger_hang_Params C.Struct

func NewHanger_hang_Params(s *C.Segment) Hanger_hang_Params {
	return Hanger_hang_Params(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 0}))
}
func NewRootHanger_hang_Params(s *C.Segment) Hanger_hang_Params {
	return Hanger_hang_Params(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 0}))
}
func AutoNewHanger_hang_Params(s *C.Segment) Hanger_hang_Params {
	return Hanger_hang_Params(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 0}))
}
func ReadRootHanger_hang_Params(s *C.Segment) Hanger_hang_Params {
	return Hanger_hang_Params(s.Root(0).ToStruct())
}

// capnp.JSON_enabled == false so we stub MarshalJSON().
func (s Hanger_hang_Params) MarshalJSON() (bs []byte, err error) { return }

type Hanger_hang_Params_List C.PointerList

func NewHanger_hang_Params_List(s *C.Segment, sz int) Hanger_hang_Params_List {
	return Hanger_hang_Params_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 0}, sz))
}
func (s Hanger_hang_Params_List) Len() int { return C.PointerList(s).Len() }
func (s Hanger_hang_Params_List) At(i int) Hanger_hang_Params {
	return Hanger_hang_Params(C.PointerList(s).At(i).ToStruct())
}
func (s Hanger_hang_Params_List) Set(i int, item Hanger_hang_Params) {
	C.PointerList(s).Set(i, C.Object(item))
}

type Hanger_hang_Params_Promise C.Pipeline

func (p *Hanger_hang_Params_Promise) Get() (Hanger_hang_Params, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Hanger_hang_Params(s), err
}

type Hanger_hang_Results C.Struct

func NewHanger_hang_Results(s *C.Segment) Hanger_hang_Results {
	return Hanger_hang_Results(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 0}))
}
func NewRootHanger_hang_Results(s *C.Segment) Hanger_hang_Results {
	return Hanger_hang_Results(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 0}))
}
func AutoNewHanger_hang_Results(s *C.Segment) Hanger_hang_Results {
	return Hanger_hang_Results(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 0}))
}
func ReadRootHanger_hang_Results(s *C.Segment) Hanger_hang_Results {
	return Hanger_hang_Results(s.Root(0).ToStruct())
}

// capnp.JSON_enabled == false so we stub MarshalJSON().
func (s Hanger_hang_Results) MarshalJSON() (bs []byte, err error) { return }

type Hanger_hang_Results_List C.PointerList

func NewHanger_hang_Results_List(s *C.Segment, sz int) Hanger_hang_Results_List {
	return Hanger_hang_Results_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 0}, sz))
}
func (s Hanger_hang_Results_List) Len() int { return C.PointerList(s).Len() }
func (s Hanger_hang_Results_List) At(i int) Hanger_hang_Results {
	return Hanger_hang_Results(C.PointerList(s).At(i).ToStruct())
}
func (s Hanger_hang_Results_List) Set(i int, item Hanger_hang_Results) {
	C.PointerList(s).Set(i, C.Object(item))
}

type Hanger_hang_Results_Promise C.Pipeline

func (p *Hanger_hang_Results_Promise) Get() (Hanger_hang_Results, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Hanger_hang_Results(s), err
}

type CallOrder struct{ c C.Client }

func NewCallOrder(c C.Client) CallOrder { return CallOrder{c} }

func (c CallOrder) GenericClient() C.Client { return c.c }

func (c CallOrder) IsNull() bool { return c.c == nil }

func (c CallOrder) GetCallSequence(ctx context.Context, params func(CallOrder_getCallSequence_Params), opts ...C.CallOption) *CallOrder_getCallSequence_Results_Promise {
	if c.c == nil {
		return (*CallOrder_getCallSequence_Results_Promise)(C.NewPipeline(C.ErrorAnswer(C.ErrNullClient)))
	}
	return (*CallOrder_getCallSequence_Results_Promise)(C.NewPipeline(c.c.Call(&C.Call{
		Ctx: ctx,
		Method: C.Method{

			InterfaceID:   0x92c5ca8314cdd2a5,
			MethodID:      0,
			InterfaceName: "test.capnp:CallOrder",
			MethodName:    "getCallSequence",
		},
		ParamsSize: C.ObjectSize{DataSize: 8, PointerCount: 0},
		ParamsFunc: func(s C.Struct) { params(CallOrder_getCallSequence_Params(s)) },
		Options:    C.NewCallOptions(opts),
	})))
}

type CallOrder_Server interface {
	GetCallSequence(CallOrder_getCallSequence) error
}

func CallOrder_ServerToClient(s CallOrder_Server) CallOrder {
	c, _ := s.(C.Closer)
	return NewCallOrder(C.NewServer(CallOrder_Methods(nil, s), c))
}

func CallOrder_Methods(methods []C.ServerMethod, server CallOrder_Server) []C.ServerMethod {
	if cap(methods) == 0 {
		methods = make([]C.ServerMethod, 0, 1)
	}

	methods = append(methods, C.ServerMethod{
		Method: C.Method{

			InterfaceID:   0x92c5ca8314cdd2a5,
			MethodID:      0,
			InterfaceName: "test.capnp:CallOrder",
			MethodName:    "getCallSequence",
		},
		Impl: func(c context.Context, opts C.CallOptions, p, r C.Struct) error {
			call := CallOrder_getCallSequence{c, opts, CallOrder_getCallSequence_Params(p), CallOrder_getCallSequence_Results(r)}
			return server.GetCallSequence(call)
		},
		ResultsSize: C.ObjectSize{DataSize: 8, PointerCount: 0},
	})

	return methods
}

// CallOrder_getCallSequence holds the arguments for a server call to CallOrder.getCallSequence.
type CallOrder_getCallSequence struct {
	Ctx     context.Context
	Options C.CallOptions
	Params  CallOrder_getCallSequence_Params
	Results CallOrder_getCallSequence_Results
}

type CallOrder_getCallSequence_Params C.Struct

func NewCallOrder_getCallSequence_Params(s *C.Segment) CallOrder_getCallSequence_Params {
	return CallOrder_getCallSequence_Params(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 0}))
}
func NewRootCallOrder_getCallSequence_Params(s *C.Segment) CallOrder_getCallSequence_Params {
	return CallOrder_getCallSequence_Params(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 0}))
}
func AutoNewCallOrder_getCallSequence_Params(s *C.Segment) CallOrder_getCallSequence_Params {
	return CallOrder_getCallSequence_Params(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 0}))
}
func ReadRootCallOrder_getCallSequence_Params(s *C.Segment) CallOrder_getCallSequence_Params {
	return CallOrder_getCallSequence_Params(s.Root(0).ToStruct())
}
func (s CallOrder_getCallSequence_Params) Expected() uint32     { return C.Struct(s).Get32(0) }
func (s CallOrder_getCallSequence_Params) SetExpected(v uint32) { C.Struct(s).Set32(0, v) }

// capnp.JSON_enabled == false so we stub MarshalJSON().
func (s CallOrder_getCallSequence_Params) MarshalJSON() (bs []byte, err error) { return }

type CallOrder_getCallSequence_Params_List C.PointerList

func NewCallOrder_getCallSequence_Params_List(s *C.Segment, sz int) CallOrder_getCallSequence_Params_List {
	return CallOrder_getCallSequence_Params_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 0}, sz))
}
func (s CallOrder_getCallSequence_Params_List) Len() int { return C.PointerList(s).Len() }
func (s CallOrder_getCallSequence_Params_List) At(i int) CallOrder_getCallSequence_Params {
	return CallOrder_getCallSequence_Params(C.PointerList(s).At(i).ToStruct())
}
func (s CallOrder_getCallSequence_Params_List) Set(i int, item CallOrder_getCallSequence_Params) {
	C.PointerList(s).Set(i, C.Object(item))
}

type CallOrder_getCallSequence_Params_Promise C.Pipeline

func (p *CallOrder_getCallSequence_Params_Promise) Get() (CallOrder_getCallSequence_Params, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return CallOrder_getCallSequence_Params(s), err
}

type CallOrder_getCallSequence_Results C.Struct

func NewCallOrder_getCallSequence_Results(s *C.Segment) CallOrder_getCallSequence_Results {
	return CallOrder_getCallSequence_Results(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 0}))
}
func NewRootCallOrder_getCallSequence_Results(s *C.Segment) CallOrder_getCallSequence_Results {
	return CallOrder_getCallSequence_Results(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 0}))
}
func AutoNewCallOrder_getCallSequence_Results(s *C.Segment) CallOrder_getCallSequence_Results {
	return CallOrder_getCallSequence_Results(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 0}))
}
func ReadRootCallOrder_getCallSequence_Results(s *C.Segment) CallOrder_getCallSequence_Results {
	return CallOrder_getCallSequence_Results(s.Root(0).ToStruct())
}
func (s CallOrder_getCallSequence_Results) N() uint32     { return C.Struct(s).Get32(0) }
func (s CallOrder_getCallSequence_Results) SetN(v uint32) { C.Struct(s).Set32(0, v) }

// capnp.JSON_enabled == false so we stub MarshalJSON().
func (s CallOrder_getCallSequence_Results) MarshalJSON() (bs []byte, err error) { return }

type CallOrder_getCallSequence_Results_List C.PointerList

func NewCallOrder_getCallSequence_Results_List(s *C.Segment, sz int) CallOrder_getCallSequence_Results_List {
	return CallOrder_getCallSequence_Results_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 0}, sz))
}
func (s CallOrder_getCallSequence_Results_List) Len() int { return C.PointerList(s).Len() }
func (s CallOrder_getCallSequence_Results_List) At(i int) CallOrder_getCallSequence_Results {
	return CallOrder_getCallSequence_Results(C.PointerList(s).At(i).ToStruct())
}
func (s CallOrder_getCallSequence_Results_List) Set(i int, item CallOrder_getCallSequence_Results) {
	C.PointerList(s).Set(i, C.Object(item))
}

type CallOrder_getCallSequence_Results_Promise C.Pipeline

func (p *CallOrder_getCallSequence_Results_Promise) Get() (CallOrder_getCallSequence_Results, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return CallOrder_getCallSequence_Results(s), err
}

type Echoer struct{ c C.Client }

func NewEchoer(c C.Client) Echoer { return Echoer{c} }

func (c Echoer) GenericClient() C.Client { return c.c }

func (c Echoer) IsNull() bool { return c.c == nil }

func (c Echoer) Echo(ctx context.Context, params func(Echoer_echo_Params), opts ...C.CallOption) *Echoer_echo_Results_Promise {
	if c.c == nil {
		return (*Echoer_echo_Results_Promise)(C.NewPipeline(C.ErrorAnswer(C.ErrNullClient)))
	}
	return (*Echoer_echo_Results_Promise)(C.NewPipeline(c.c.Call(&C.Call{
		Ctx: ctx,
		Method: C.Method{

			InterfaceID:   0x841756c6a41b2a45,
			MethodID:      0,
			InterfaceName: "test.capnp:Echoer",
			MethodName:    "echo",
		},
		ParamsSize: C.ObjectSize{DataSize: 0, PointerCount: 1},
		ParamsFunc: func(s C.Struct) { params(Echoer_echo_Params(s)) },
		Options:    C.NewCallOptions(opts),
	})))
}

func (c Echoer) GetCallSequence(ctx context.Context, params func(CallOrder_getCallSequence_Params), opts ...C.CallOption) *CallOrder_getCallSequence_Results_Promise {
	if c.c == nil {
		return (*CallOrder_getCallSequence_Results_Promise)(C.NewPipeline(C.ErrorAnswer(C.ErrNullClient)))
	}
	return (*CallOrder_getCallSequence_Results_Promise)(C.NewPipeline(c.c.Call(&C.Call{
		Ctx: ctx,
		Method: C.Method{

			InterfaceID:   0x92c5ca8314cdd2a5,
			MethodID:      0,
			InterfaceName: "test.capnp:CallOrder",
			MethodName:    "getCallSequence",
		},
		ParamsSize: C.ObjectSize{DataSize: 8, PointerCount: 0},
		ParamsFunc: func(s C.Struct) { params(CallOrder_getCallSequence_Params(s)) },
		Options:    C.NewCallOptions(opts),
	})))
}

type Echoer_Server interface {
	Echo(Echoer_echo) error

	GetCallSequence(CallOrder_getCallSequence) error
}

func Echoer_ServerToClient(s Echoer_Server) Echoer {
	c, _ := s.(C.Closer)
	return NewEchoer(C.NewServer(Echoer_Methods(nil, s), c))
}

func Echoer_Methods(methods []C.ServerMethod, server Echoer_Server) []C.ServerMethod {
	if cap(methods) == 0 {
		methods = make([]C.ServerMethod, 0, 2)
	}

	methods = append(methods, C.ServerMethod{
		Method: C.Method{

			InterfaceID:   0x841756c6a41b2a45,
			MethodID:      0,
			InterfaceName: "test.capnp:Echoer",
			MethodName:    "echo",
		},
		Impl: func(c context.Context, opts C.CallOptions, p, r C.Struct) error {
			call := Echoer_echo{c, opts, Echoer_echo_Params(p), Echoer_echo_Results(r)}
			return server.Echo(call)
		},
		ResultsSize: C.ObjectSize{DataSize: 0, PointerCount: 1},
	})

	methods = append(methods, C.ServerMethod{
		Method: C.Method{

			InterfaceID:   0x92c5ca8314cdd2a5,
			MethodID:      0,
			InterfaceName: "test.capnp:CallOrder",
			MethodName:    "getCallSequence",
		},
		Impl: func(c context.Context, opts C.CallOptions, p, r C.Struct) error {
			call := CallOrder_getCallSequence{c, opts, CallOrder_getCallSequence_Params(p), CallOrder_getCallSequence_Results(r)}
			return server.GetCallSequence(call)
		},
		ResultsSize: C.ObjectSize{DataSize: 8, PointerCount: 0},
	})

	return methods
}

// Echoer_echo holds the arguments for a server call to Echoer.echo.
type Echoer_echo struct {
	Ctx     context.Context
	Options C.CallOptions
	Params  Echoer_echo_Params
	Results Echoer_echo_Results
}

type Echoer_echo_Params C.Struct

func NewEchoer_echo_Params(s *C.Segment) Echoer_echo_Params {
	return Echoer_echo_Params(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func NewRootEchoer_echo_Params(s *C.Segment) Echoer_echo_Params {
	return Echoer_echo_Params(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewEchoer_echo_Params(s *C.Segment) Echoer_echo_Params {
	return Echoer_echo_Params(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootEchoer_echo_Params(s *C.Segment) Echoer_echo_Params {
	return Echoer_echo_Params(s.Root(0).ToStruct())
}
func (s Echoer_echo_Params) Cap() CallOrder {
	return NewCallOrder(C.Struct(s).GetObject(0).ToInterface().Client())
}
func (s Echoer_echo_Params) SetCap(v CallOrder) {
	if s.Segment == nil {
		return
	}
	ci := s.Segment.Message.AddCap(v.GenericClient())
	C.Struct(s).SetObject(0, C.Object(s.Segment.NewInterface(ci)))
}

// capnp.JSON_enabled == false so we stub MarshalJSON().
func (s Echoer_echo_Params) MarshalJSON() (bs []byte, err error) { return }

type Echoer_echo_Params_List C.PointerList

func NewEchoer_echo_Params_List(s *C.Segment, sz int) Echoer_echo_Params_List {
	return Echoer_echo_Params_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s Echoer_echo_Params_List) Len() int { return C.PointerList(s).Len() }
func (s Echoer_echo_Params_List) At(i int) Echoer_echo_Params {
	return Echoer_echo_Params(C.PointerList(s).At(i).ToStruct())
}
func (s Echoer_echo_Params_List) Set(i int, item Echoer_echo_Params) {
	C.PointerList(s).Set(i, C.Object(item))
}

type Echoer_echo_Params_Promise C.Pipeline

func (p *Echoer_echo_Params_Promise) Get() (Echoer_echo_Params, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Echoer_echo_Params(s), err
}

func (p *Echoer_echo_Params_Promise) Cap() CallOrder {
	return NewCallOrder((*C.Pipeline)(p).GetPipeline(0).Client())
}

type Echoer_echo_Results C.Struct

func NewEchoer_echo_Results(s *C.Segment) Echoer_echo_Results {
	return Echoer_echo_Results(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func NewRootEchoer_echo_Results(s *C.Segment) Echoer_echo_Results {
	return Echoer_echo_Results(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewEchoer_echo_Results(s *C.Segment) Echoer_echo_Results {
	return Echoer_echo_Results(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootEchoer_echo_Results(s *C.Segment) Echoer_echo_Results {
	return Echoer_echo_Results(s.Root(0).ToStruct())
}
func (s Echoer_echo_Results) Cap() CallOrder {
	return NewCallOrder(C.Struct(s).GetObject(0).ToInterface().Client())
}
func (s Echoer_echo_Results) SetCap(v CallOrder) {
	if s.Segment == nil {
		return
	}
	ci := s.Segment.Message.AddCap(v.GenericClient())
	C.Struct(s).SetObject(0, C.Object(s.Segment.NewInterface(ci)))
}

// capnp.JSON_enabled == false so we stub MarshalJSON().
func (s Echoer_echo_Results) MarshalJSON() (bs []byte, err error) { return }

type Echoer_echo_Results_List C.PointerList

func NewEchoer_echo_Results_List(s *C.Segment, sz int) Echoer_echo_Results_List {
	return Echoer_echo_Results_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s Echoer_echo_Results_List) Len() int { return C.PointerList(s).Len() }
func (s Echoer_echo_Results_List) At(i int) Echoer_echo_Results {
	return Echoer_echo_Results(C.PointerList(s).At(i).ToStruct())
}
func (s Echoer_echo_Results_List) Set(i int, item Echoer_echo_Results) {
	C.PointerList(s).Set(i, C.Object(item))
}

type Echoer_echo_Results_Promise C.Pipeline

func (p *Echoer_echo_Results_Promise) Get() (Echoer_echo_Results, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Echoer_echo_Results(s), err
}

func (p *Echoer_echo_Results_Promise) Cap() CallOrder {
	return NewCallOrder((*C.Pipeline)(p).GetPipeline(0).Client())
}
