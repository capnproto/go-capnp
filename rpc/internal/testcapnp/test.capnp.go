package testcapnp

// AUTO GENERATED - DO NOT EDIT

import (
	context "golang.org/x/net/context"
	capnp "zombiezen.com/go/capnproto2"
	server "zombiezen.com/go/capnproto2/server"
)

type Handle struct{ Client capnp.Client }

type Handle_Server interface {
}

func Handle_ServerToClient(s Handle_Server) Handle {
	c, _ := s.(server.Closer)
	return Handle{Client: server.New(Handle_Methods(nil, s), c)}
}

func Handle_Methods(methods []server.Method, s Handle_Server) []server.Method {
	if cap(methods) == 0 {
		methods = make([]server.Method, 0, 0)
	}

	return methods
}

type HandleFactory struct{ Client capnp.Client }

func (c HandleFactory) NewHandle(ctx context.Context, params func(HandleFactory_newHandle_Params) error, opts ...capnp.CallOption) HandleFactory_newHandle_Results_Promise {
	if c.Client == nil {
		return HandleFactory_newHandle_Results_Promise{Pipeline: capnp.NewPipeline(capnp.ErrorAnswer(capnp.ErrNullClient))}
	}
	call := &capnp.Call{
		Ctx: ctx,
		Method: capnp.Method{
			InterfaceID:   0x8491a7fe75fe0bce,
			MethodID:      0,
			InterfaceName: "test.capnp:HandleFactory",
			MethodName:    "newHandle",
		},
		Options: capnp.NewCallOptions(opts),
	}
	if params != nil {
		call.ParamsSize = capnp.ObjectSize{DataSize: 0, PointerCount: 0}
		call.ParamsFunc = func(s capnp.Struct) error { return params(HandleFactory_newHandle_Params{Struct: s}) }
	}
	return HandleFactory_newHandle_Results_Promise{Pipeline: capnp.NewPipeline(c.Client.Call(call))}
}

type HandleFactory_Server interface {
	NewHandle(HandleFactory_newHandle) error
}

func HandleFactory_ServerToClient(s HandleFactory_Server) HandleFactory {
	c, _ := s.(server.Closer)
	return HandleFactory{Client: server.New(HandleFactory_Methods(nil, s), c)}
}

func HandleFactory_Methods(methods []server.Method, s HandleFactory_Server) []server.Method {
	if cap(methods) == 0 {
		methods = make([]server.Method, 0, 1)
	}

	methods = append(methods, server.Method{
		Method: capnp.Method{
			InterfaceID:   0x8491a7fe75fe0bce,
			MethodID:      0,
			InterfaceName: "test.capnp:HandleFactory",
			MethodName:    "newHandle",
		},
		Impl: func(c context.Context, opts capnp.CallOptions, p, r capnp.Struct) error {
			call := HandleFactory_newHandle{c, opts, HandleFactory_newHandle_Params{Struct: p}, HandleFactory_newHandle_Results{Struct: r}}
			return s.NewHandle(call)
		},
		ResultsSize: capnp.ObjectSize{DataSize: 0, PointerCount: 1},
	})

	return methods
}

// HandleFactory_newHandle holds the arguments for a server call to HandleFactory.newHandle.
type HandleFactory_newHandle struct {
	Ctx     context.Context
	Options capnp.CallOptions
	Params  HandleFactory_newHandle_Params
	Results HandleFactory_newHandle_Results
}

type HandleFactory_newHandle_Params struct{ capnp.Struct }

func NewHandleFactory_newHandle_Params(s *capnp.Segment) (HandleFactory_newHandle_Params, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0})
	if err != nil {
		return HandleFactory_newHandle_Params{}, err
	}
	return HandleFactory_newHandle_Params{st}, nil
}

func NewRootHandleFactory_newHandle_Params(s *capnp.Segment) (HandleFactory_newHandle_Params, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0})
	if err != nil {
		return HandleFactory_newHandle_Params{}, err
	}
	return HandleFactory_newHandle_Params{st}, nil
}

func ReadRootHandleFactory_newHandle_Params(msg *capnp.Message) (HandleFactory_newHandle_Params, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return HandleFactory_newHandle_Params{}, err
	}
	return HandleFactory_newHandle_Params{root.Struct()}, nil
}

// HandleFactory_newHandle_Params_List is a list of HandleFactory_newHandle_Params.
type HandleFactory_newHandle_Params_List struct{ capnp.List }

// NewHandleFactory_newHandle_Params creates a new list of HandleFactory_newHandle_Params.
func NewHandleFactory_newHandle_Params_List(s *capnp.Segment, sz int32) (HandleFactory_newHandle_Params_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0}, sz)
	if err != nil {
		return HandleFactory_newHandle_Params_List{}, err
	}
	return HandleFactory_newHandle_Params_List{l}, nil
}

func (s HandleFactory_newHandle_Params_List) At(i int) HandleFactory_newHandle_Params {
	return HandleFactory_newHandle_Params{s.List.Struct(i)}
}
func (s HandleFactory_newHandle_Params_List) Set(i int, v HandleFactory_newHandle_Params) error {
	return s.List.SetStruct(i, v.Struct)
}

// HandleFactory_newHandle_Params_Promise is a wrapper for a HandleFactory_newHandle_Params promised by a client call.
type HandleFactory_newHandle_Params_Promise struct{ *capnp.Pipeline }

func (p HandleFactory_newHandle_Params_Promise) Struct() (HandleFactory_newHandle_Params, error) {
	s, err := p.Pipeline.Struct()
	return HandleFactory_newHandle_Params{s}, err
}

type HandleFactory_newHandle_Results struct{ capnp.Struct }

func NewHandleFactory_newHandle_Results(s *capnp.Segment) (HandleFactory_newHandle_Results, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return HandleFactory_newHandle_Results{}, err
	}
	return HandleFactory_newHandle_Results{st}, nil
}

func NewRootHandleFactory_newHandle_Results(s *capnp.Segment) (HandleFactory_newHandle_Results, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return HandleFactory_newHandle_Results{}, err
	}
	return HandleFactory_newHandle_Results{st}, nil
}

func ReadRootHandleFactory_newHandle_Results(msg *capnp.Message) (HandleFactory_newHandle_Results, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return HandleFactory_newHandle_Results{}, err
	}
	return HandleFactory_newHandle_Results{root.Struct()}, nil
}
func (s HandleFactory_newHandle_Results) Handle() Handle {
	p, err := s.Struct.Ptr(0)
	if err != nil {

		return Handle{}
	}
	return Handle{Client: p.Interface().Client()}
}

func (s HandleFactory_newHandle_Results) HasHandle() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s HandleFactory_newHandle_Results) SetHandle(v Handle) error {
	seg := s.Segment()
	if seg == nil {

		return nil
	}
	var in capnp.Interface
	if v.Client != nil {
		in = capnp.NewInterface(seg, seg.Message().AddCap(v.Client))
	}
	return s.Struct.SetPtr(0, in.ToPtr())
}

// HandleFactory_newHandle_Results_List is a list of HandleFactory_newHandle_Results.
type HandleFactory_newHandle_Results_List struct{ capnp.List }

// NewHandleFactory_newHandle_Results creates a new list of HandleFactory_newHandle_Results.
func NewHandleFactory_newHandle_Results_List(s *capnp.Segment, sz int32) (HandleFactory_newHandle_Results_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return HandleFactory_newHandle_Results_List{}, err
	}
	return HandleFactory_newHandle_Results_List{l}, nil
}

func (s HandleFactory_newHandle_Results_List) At(i int) HandleFactory_newHandle_Results {
	return HandleFactory_newHandle_Results{s.List.Struct(i)}
}
func (s HandleFactory_newHandle_Results_List) Set(i int, v HandleFactory_newHandle_Results) error {
	return s.List.SetStruct(i, v.Struct)
}

// HandleFactory_newHandle_Results_Promise is a wrapper for a HandleFactory_newHandle_Results promised by a client call.
type HandleFactory_newHandle_Results_Promise struct{ *capnp.Pipeline }

func (p HandleFactory_newHandle_Results_Promise) Struct() (HandleFactory_newHandle_Results, error) {
	s, err := p.Pipeline.Struct()
	return HandleFactory_newHandle_Results{s}, err
}

func (p HandleFactory_newHandle_Results_Promise) Handle() Handle {
	return Handle{Client: p.Pipeline.GetPipeline(0).Client()}
}

type Hanger struct{ Client capnp.Client }

func (c Hanger) Hang(ctx context.Context, params func(Hanger_hang_Params) error, opts ...capnp.CallOption) Hanger_hang_Results_Promise {
	if c.Client == nil {
		return Hanger_hang_Results_Promise{Pipeline: capnp.NewPipeline(capnp.ErrorAnswer(capnp.ErrNullClient))}
	}
	call := &capnp.Call{
		Ctx: ctx,
		Method: capnp.Method{
			InterfaceID:   0x8ae08044aae8a26e,
			MethodID:      0,
			InterfaceName: "test.capnp:Hanger",
			MethodName:    "hang",
		},
		Options: capnp.NewCallOptions(opts),
	}
	if params != nil {
		call.ParamsSize = capnp.ObjectSize{DataSize: 0, PointerCount: 0}
		call.ParamsFunc = func(s capnp.Struct) error { return params(Hanger_hang_Params{Struct: s}) }
	}
	return Hanger_hang_Results_Promise{Pipeline: capnp.NewPipeline(c.Client.Call(call))}
}

type Hanger_Server interface {
	Hang(Hanger_hang) error
}

func Hanger_ServerToClient(s Hanger_Server) Hanger {
	c, _ := s.(server.Closer)
	return Hanger{Client: server.New(Hanger_Methods(nil, s), c)}
}

func Hanger_Methods(methods []server.Method, s Hanger_Server) []server.Method {
	if cap(methods) == 0 {
		methods = make([]server.Method, 0, 1)
	}

	methods = append(methods, server.Method{
		Method: capnp.Method{
			InterfaceID:   0x8ae08044aae8a26e,
			MethodID:      0,
			InterfaceName: "test.capnp:Hanger",
			MethodName:    "hang",
		},
		Impl: func(c context.Context, opts capnp.CallOptions, p, r capnp.Struct) error {
			call := Hanger_hang{c, opts, Hanger_hang_Params{Struct: p}, Hanger_hang_Results{Struct: r}}
			return s.Hang(call)
		},
		ResultsSize: capnp.ObjectSize{DataSize: 0, PointerCount: 0},
	})

	return methods
}

// Hanger_hang holds the arguments for a server call to Hanger.hang.
type Hanger_hang struct {
	Ctx     context.Context
	Options capnp.CallOptions
	Params  Hanger_hang_Params
	Results Hanger_hang_Results
}

type Hanger_hang_Params struct{ capnp.Struct }

func NewHanger_hang_Params(s *capnp.Segment) (Hanger_hang_Params, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0})
	if err != nil {
		return Hanger_hang_Params{}, err
	}
	return Hanger_hang_Params{st}, nil
}

func NewRootHanger_hang_Params(s *capnp.Segment) (Hanger_hang_Params, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0})
	if err != nil {
		return Hanger_hang_Params{}, err
	}
	return Hanger_hang_Params{st}, nil
}

func ReadRootHanger_hang_Params(msg *capnp.Message) (Hanger_hang_Params, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Hanger_hang_Params{}, err
	}
	return Hanger_hang_Params{root.Struct()}, nil
}

// Hanger_hang_Params_List is a list of Hanger_hang_Params.
type Hanger_hang_Params_List struct{ capnp.List }

// NewHanger_hang_Params creates a new list of Hanger_hang_Params.
func NewHanger_hang_Params_List(s *capnp.Segment, sz int32) (Hanger_hang_Params_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0}, sz)
	if err != nil {
		return Hanger_hang_Params_List{}, err
	}
	return Hanger_hang_Params_List{l}, nil
}

func (s Hanger_hang_Params_List) At(i int) Hanger_hang_Params {
	return Hanger_hang_Params{s.List.Struct(i)}
}
func (s Hanger_hang_Params_List) Set(i int, v Hanger_hang_Params) error {
	return s.List.SetStruct(i, v.Struct)
}

// Hanger_hang_Params_Promise is a wrapper for a Hanger_hang_Params promised by a client call.
type Hanger_hang_Params_Promise struct{ *capnp.Pipeline }

func (p Hanger_hang_Params_Promise) Struct() (Hanger_hang_Params, error) {
	s, err := p.Pipeline.Struct()
	return Hanger_hang_Params{s}, err
}

type Hanger_hang_Results struct{ capnp.Struct }

func NewHanger_hang_Results(s *capnp.Segment) (Hanger_hang_Results, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0})
	if err != nil {
		return Hanger_hang_Results{}, err
	}
	return Hanger_hang_Results{st}, nil
}

func NewRootHanger_hang_Results(s *capnp.Segment) (Hanger_hang_Results, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0})
	if err != nil {
		return Hanger_hang_Results{}, err
	}
	return Hanger_hang_Results{st}, nil
}

func ReadRootHanger_hang_Results(msg *capnp.Message) (Hanger_hang_Results, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Hanger_hang_Results{}, err
	}
	return Hanger_hang_Results{root.Struct()}, nil
}

// Hanger_hang_Results_List is a list of Hanger_hang_Results.
type Hanger_hang_Results_List struct{ capnp.List }

// NewHanger_hang_Results creates a new list of Hanger_hang_Results.
func NewHanger_hang_Results_List(s *capnp.Segment, sz int32) (Hanger_hang_Results_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0}, sz)
	if err != nil {
		return Hanger_hang_Results_List{}, err
	}
	return Hanger_hang_Results_List{l}, nil
}

func (s Hanger_hang_Results_List) At(i int) Hanger_hang_Results {
	return Hanger_hang_Results{s.List.Struct(i)}
}
func (s Hanger_hang_Results_List) Set(i int, v Hanger_hang_Results) error {
	return s.List.SetStruct(i, v.Struct)
}

// Hanger_hang_Results_Promise is a wrapper for a Hanger_hang_Results promised by a client call.
type Hanger_hang_Results_Promise struct{ *capnp.Pipeline }

func (p Hanger_hang_Results_Promise) Struct() (Hanger_hang_Results, error) {
	s, err := p.Pipeline.Struct()
	return Hanger_hang_Results{s}, err
}

type CallOrder struct{ Client capnp.Client }

func (c CallOrder) GetCallSequence(ctx context.Context, params func(CallOrder_getCallSequence_Params) error, opts ...capnp.CallOption) CallOrder_getCallSequence_Results_Promise {
	if c.Client == nil {
		return CallOrder_getCallSequence_Results_Promise{Pipeline: capnp.NewPipeline(capnp.ErrorAnswer(capnp.ErrNullClient))}
	}
	call := &capnp.Call{
		Ctx: ctx,
		Method: capnp.Method{
			InterfaceID:   0x92c5ca8314cdd2a5,
			MethodID:      0,
			InterfaceName: "test.capnp:CallOrder",
			MethodName:    "getCallSequence",
		},
		Options: capnp.NewCallOptions(opts),
	}
	if params != nil {
		call.ParamsSize = capnp.ObjectSize{DataSize: 8, PointerCount: 0}
		call.ParamsFunc = func(s capnp.Struct) error { return params(CallOrder_getCallSequence_Params{Struct: s}) }
	}
	return CallOrder_getCallSequence_Results_Promise{Pipeline: capnp.NewPipeline(c.Client.Call(call))}
}

type CallOrder_Server interface {
	GetCallSequence(CallOrder_getCallSequence) error
}

func CallOrder_ServerToClient(s CallOrder_Server) CallOrder {
	c, _ := s.(server.Closer)
	return CallOrder{Client: server.New(CallOrder_Methods(nil, s), c)}
}

func CallOrder_Methods(methods []server.Method, s CallOrder_Server) []server.Method {
	if cap(methods) == 0 {
		methods = make([]server.Method, 0, 1)
	}

	methods = append(methods, server.Method{
		Method: capnp.Method{
			InterfaceID:   0x92c5ca8314cdd2a5,
			MethodID:      0,
			InterfaceName: "test.capnp:CallOrder",
			MethodName:    "getCallSequence",
		},
		Impl: func(c context.Context, opts capnp.CallOptions, p, r capnp.Struct) error {
			call := CallOrder_getCallSequence{c, opts, CallOrder_getCallSequence_Params{Struct: p}, CallOrder_getCallSequence_Results{Struct: r}}
			return s.GetCallSequence(call)
		},
		ResultsSize: capnp.ObjectSize{DataSize: 8, PointerCount: 0},
	})

	return methods
}

// CallOrder_getCallSequence holds the arguments for a server call to CallOrder.getCallSequence.
type CallOrder_getCallSequence struct {
	Ctx     context.Context
	Options capnp.CallOptions
	Params  CallOrder_getCallSequence_Params
	Results CallOrder_getCallSequence_Results
}

type CallOrder_getCallSequence_Params struct{ capnp.Struct }

func NewCallOrder_getCallSequence_Params(s *capnp.Segment) (CallOrder_getCallSequence_Params, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	if err != nil {
		return CallOrder_getCallSequence_Params{}, err
	}
	return CallOrder_getCallSequence_Params{st}, nil
}

func NewRootCallOrder_getCallSequence_Params(s *capnp.Segment) (CallOrder_getCallSequence_Params, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	if err != nil {
		return CallOrder_getCallSequence_Params{}, err
	}
	return CallOrder_getCallSequence_Params{st}, nil
}

func ReadRootCallOrder_getCallSequence_Params(msg *capnp.Message) (CallOrder_getCallSequence_Params, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return CallOrder_getCallSequence_Params{}, err
	}
	return CallOrder_getCallSequence_Params{root.Struct()}, nil
}
func (s CallOrder_getCallSequence_Params) Expected() uint32 {
	return s.Struct.Uint32(0)
}

func (s CallOrder_getCallSequence_Params) SetExpected(v uint32) {
	s.Struct.SetUint32(0, v)
}

// CallOrder_getCallSequence_Params_List is a list of CallOrder_getCallSequence_Params.
type CallOrder_getCallSequence_Params_List struct{ capnp.List }

// NewCallOrder_getCallSequence_Params creates a new list of CallOrder_getCallSequence_Params.
func NewCallOrder_getCallSequence_Params_List(s *capnp.Segment, sz int32) (CallOrder_getCallSequence_Params_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0}, sz)
	if err != nil {
		return CallOrder_getCallSequence_Params_List{}, err
	}
	return CallOrder_getCallSequence_Params_List{l}, nil
}

func (s CallOrder_getCallSequence_Params_List) At(i int) CallOrder_getCallSequence_Params {
	return CallOrder_getCallSequence_Params{s.List.Struct(i)}
}
func (s CallOrder_getCallSequence_Params_List) Set(i int, v CallOrder_getCallSequence_Params) error {
	return s.List.SetStruct(i, v.Struct)
}

// CallOrder_getCallSequence_Params_Promise is a wrapper for a CallOrder_getCallSequence_Params promised by a client call.
type CallOrder_getCallSequence_Params_Promise struct{ *capnp.Pipeline }

func (p CallOrder_getCallSequence_Params_Promise) Struct() (CallOrder_getCallSequence_Params, error) {
	s, err := p.Pipeline.Struct()
	return CallOrder_getCallSequence_Params{s}, err
}

type CallOrder_getCallSequence_Results struct{ capnp.Struct }

func NewCallOrder_getCallSequence_Results(s *capnp.Segment) (CallOrder_getCallSequence_Results, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	if err != nil {
		return CallOrder_getCallSequence_Results{}, err
	}
	return CallOrder_getCallSequence_Results{st}, nil
}

func NewRootCallOrder_getCallSequence_Results(s *capnp.Segment) (CallOrder_getCallSequence_Results, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	if err != nil {
		return CallOrder_getCallSequence_Results{}, err
	}
	return CallOrder_getCallSequence_Results{st}, nil
}

func ReadRootCallOrder_getCallSequence_Results(msg *capnp.Message) (CallOrder_getCallSequence_Results, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return CallOrder_getCallSequence_Results{}, err
	}
	return CallOrder_getCallSequence_Results{root.Struct()}, nil
}
func (s CallOrder_getCallSequence_Results) N() uint32 {
	return s.Struct.Uint32(0)
}

func (s CallOrder_getCallSequence_Results) SetN(v uint32) {
	s.Struct.SetUint32(0, v)
}

// CallOrder_getCallSequence_Results_List is a list of CallOrder_getCallSequence_Results.
type CallOrder_getCallSequence_Results_List struct{ capnp.List }

// NewCallOrder_getCallSequence_Results creates a new list of CallOrder_getCallSequence_Results.
func NewCallOrder_getCallSequence_Results_List(s *capnp.Segment, sz int32) (CallOrder_getCallSequence_Results_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0}, sz)
	if err != nil {
		return CallOrder_getCallSequence_Results_List{}, err
	}
	return CallOrder_getCallSequence_Results_List{l}, nil
}

func (s CallOrder_getCallSequence_Results_List) At(i int) CallOrder_getCallSequence_Results {
	return CallOrder_getCallSequence_Results{s.List.Struct(i)}
}
func (s CallOrder_getCallSequence_Results_List) Set(i int, v CallOrder_getCallSequence_Results) error {
	return s.List.SetStruct(i, v.Struct)
}

// CallOrder_getCallSequence_Results_Promise is a wrapper for a CallOrder_getCallSequence_Results promised by a client call.
type CallOrder_getCallSequence_Results_Promise struct{ *capnp.Pipeline }

func (p CallOrder_getCallSequence_Results_Promise) Struct() (CallOrder_getCallSequence_Results, error) {
	s, err := p.Pipeline.Struct()
	return CallOrder_getCallSequence_Results{s}, err
}

type Echoer struct{ Client capnp.Client }

func (c Echoer) Echo(ctx context.Context, params func(Echoer_echo_Params) error, opts ...capnp.CallOption) Echoer_echo_Results_Promise {
	if c.Client == nil {
		return Echoer_echo_Results_Promise{Pipeline: capnp.NewPipeline(capnp.ErrorAnswer(capnp.ErrNullClient))}
	}
	call := &capnp.Call{
		Ctx: ctx,
		Method: capnp.Method{
			InterfaceID:   0x841756c6a41b2a45,
			MethodID:      0,
			InterfaceName: "test.capnp:Echoer",
			MethodName:    "echo",
		},
		Options: capnp.NewCallOptions(opts),
	}
	if params != nil {
		call.ParamsSize = capnp.ObjectSize{DataSize: 0, PointerCount: 1}
		call.ParamsFunc = func(s capnp.Struct) error { return params(Echoer_echo_Params{Struct: s}) }
	}
	return Echoer_echo_Results_Promise{Pipeline: capnp.NewPipeline(c.Client.Call(call))}
}
func (c Echoer) GetCallSequence(ctx context.Context, params func(CallOrder_getCallSequence_Params) error, opts ...capnp.CallOption) CallOrder_getCallSequence_Results_Promise {
	if c.Client == nil {
		return CallOrder_getCallSequence_Results_Promise{Pipeline: capnp.NewPipeline(capnp.ErrorAnswer(capnp.ErrNullClient))}
	}
	call := &capnp.Call{
		Ctx: ctx,
		Method: capnp.Method{
			InterfaceID:   0x92c5ca8314cdd2a5,
			MethodID:      0,
			InterfaceName: "test.capnp:CallOrder",
			MethodName:    "getCallSequence",
		},
		Options: capnp.NewCallOptions(opts),
	}
	if params != nil {
		call.ParamsSize = capnp.ObjectSize{DataSize: 8, PointerCount: 0}
		call.ParamsFunc = func(s capnp.Struct) error { return params(CallOrder_getCallSequence_Params{Struct: s}) }
	}
	return CallOrder_getCallSequence_Results_Promise{Pipeline: capnp.NewPipeline(c.Client.Call(call))}
}

type Echoer_Server interface {
	Echo(Echoer_echo) error

	GetCallSequence(CallOrder_getCallSequence) error
}

func Echoer_ServerToClient(s Echoer_Server) Echoer {
	c, _ := s.(server.Closer)
	return Echoer{Client: server.New(Echoer_Methods(nil, s), c)}
}

func Echoer_Methods(methods []server.Method, s Echoer_Server) []server.Method {
	if cap(methods) == 0 {
		methods = make([]server.Method, 0, 2)
	}

	methods = append(methods, server.Method{
		Method: capnp.Method{
			InterfaceID:   0x841756c6a41b2a45,
			MethodID:      0,
			InterfaceName: "test.capnp:Echoer",
			MethodName:    "echo",
		},
		Impl: func(c context.Context, opts capnp.CallOptions, p, r capnp.Struct) error {
			call := Echoer_echo{c, opts, Echoer_echo_Params{Struct: p}, Echoer_echo_Results{Struct: r}}
			return s.Echo(call)
		},
		ResultsSize: capnp.ObjectSize{DataSize: 0, PointerCount: 1},
	})

	methods = append(methods, server.Method{
		Method: capnp.Method{
			InterfaceID:   0x92c5ca8314cdd2a5,
			MethodID:      0,
			InterfaceName: "test.capnp:CallOrder",
			MethodName:    "getCallSequence",
		},
		Impl: func(c context.Context, opts capnp.CallOptions, p, r capnp.Struct) error {
			call := CallOrder_getCallSequence{c, opts, CallOrder_getCallSequence_Params{Struct: p}, CallOrder_getCallSequence_Results{Struct: r}}
			return s.GetCallSequence(call)
		},
		ResultsSize: capnp.ObjectSize{DataSize: 8, PointerCount: 0},
	})

	return methods
}

// Echoer_echo holds the arguments for a server call to Echoer.echo.
type Echoer_echo struct {
	Ctx     context.Context
	Options capnp.CallOptions
	Params  Echoer_echo_Params
	Results Echoer_echo_Results
}

type Echoer_echo_Params struct{ capnp.Struct }

func NewEchoer_echo_Params(s *capnp.Segment) (Echoer_echo_Params, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Echoer_echo_Params{}, err
	}
	return Echoer_echo_Params{st}, nil
}

func NewRootEchoer_echo_Params(s *capnp.Segment) (Echoer_echo_Params, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Echoer_echo_Params{}, err
	}
	return Echoer_echo_Params{st}, nil
}

func ReadRootEchoer_echo_Params(msg *capnp.Message) (Echoer_echo_Params, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Echoer_echo_Params{}, err
	}
	return Echoer_echo_Params{root.Struct()}, nil
}
func (s Echoer_echo_Params) Cap() CallOrder {
	p, err := s.Struct.Ptr(0)
	if err != nil {

		return CallOrder{}
	}
	return CallOrder{Client: p.Interface().Client()}
}

func (s Echoer_echo_Params) HasCap() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Echoer_echo_Params) SetCap(v CallOrder) error {
	seg := s.Segment()
	if seg == nil {

		return nil
	}
	var in capnp.Interface
	if v.Client != nil {
		in = capnp.NewInterface(seg, seg.Message().AddCap(v.Client))
	}
	return s.Struct.SetPtr(0, in.ToPtr())
}

// Echoer_echo_Params_List is a list of Echoer_echo_Params.
type Echoer_echo_Params_List struct{ capnp.List }

// NewEchoer_echo_Params creates a new list of Echoer_echo_Params.
func NewEchoer_echo_Params_List(s *capnp.Segment, sz int32) (Echoer_echo_Params_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return Echoer_echo_Params_List{}, err
	}
	return Echoer_echo_Params_List{l}, nil
}

func (s Echoer_echo_Params_List) At(i int) Echoer_echo_Params {
	return Echoer_echo_Params{s.List.Struct(i)}
}
func (s Echoer_echo_Params_List) Set(i int, v Echoer_echo_Params) error {
	return s.List.SetStruct(i, v.Struct)
}

// Echoer_echo_Params_Promise is a wrapper for a Echoer_echo_Params promised by a client call.
type Echoer_echo_Params_Promise struct{ *capnp.Pipeline }

func (p Echoer_echo_Params_Promise) Struct() (Echoer_echo_Params, error) {
	s, err := p.Pipeline.Struct()
	return Echoer_echo_Params{s}, err
}

func (p Echoer_echo_Params_Promise) Cap() CallOrder {
	return CallOrder{Client: p.Pipeline.GetPipeline(0).Client()}
}

type Echoer_echo_Results struct{ capnp.Struct }

func NewEchoer_echo_Results(s *capnp.Segment) (Echoer_echo_Results, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Echoer_echo_Results{}, err
	}
	return Echoer_echo_Results{st}, nil
}

func NewRootEchoer_echo_Results(s *capnp.Segment) (Echoer_echo_Results, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Echoer_echo_Results{}, err
	}
	return Echoer_echo_Results{st}, nil
}

func ReadRootEchoer_echo_Results(msg *capnp.Message) (Echoer_echo_Results, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Echoer_echo_Results{}, err
	}
	return Echoer_echo_Results{root.Struct()}, nil
}
func (s Echoer_echo_Results) Cap() CallOrder {
	p, err := s.Struct.Ptr(0)
	if err != nil {

		return CallOrder{}
	}
	return CallOrder{Client: p.Interface().Client()}
}

func (s Echoer_echo_Results) HasCap() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Echoer_echo_Results) SetCap(v CallOrder) error {
	seg := s.Segment()
	if seg == nil {

		return nil
	}
	var in capnp.Interface
	if v.Client != nil {
		in = capnp.NewInterface(seg, seg.Message().AddCap(v.Client))
	}
	return s.Struct.SetPtr(0, in.ToPtr())
}

// Echoer_echo_Results_List is a list of Echoer_echo_Results.
type Echoer_echo_Results_List struct{ capnp.List }

// NewEchoer_echo_Results creates a new list of Echoer_echo_Results.
func NewEchoer_echo_Results_List(s *capnp.Segment, sz int32) (Echoer_echo_Results_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return Echoer_echo_Results_List{}, err
	}
	return Echoer_echo_Results_List{l}, nil
}

func (s Echoer_echo_Results_List) At(i int) Echoer_echo_Results {
	return Echoer_echo_Results{s.List.Struct(i)}
}
func (s Echoer_echo_Results_List) Set(i int, v Echoer_echo_Results) error {
	return s.List.SetStruct(i, v.Struct)
}

// Echoer_echo_Results_Promise is a wrapper for a Echoer_echo_Results promised by a client call.
type Echoer_echo_Results_Promise struct{ *capnp.Pipeline }

func (p Echoer_echo_Results_Promise) Struct() (Echoer_echo_Results, error) {
	s, err := p.Pipeline.Struct()
	return Echoer_echo_Results{s}, err
}

func (p Echoer_echo_Results_Promise) Cap() CallOrder {
	return CallOrder{Client: p.Pipeline.GetPipeline(0).Client()}
}

type Adder struct{ Client capnp.Client }

func (c Adder) Add(ctx context.Context, params func(Adder_add_Params) error, opts ...capnp.CallOption) Adder_add_Results_Promise {
	if c.Client == nil {
		return Adder_add_Results_Promise{Pipeline: capnp.NewPipeline(capnp.ErrorAnswer(capnp.ErrNullClient))}
	}
	call := &capnp.Call{
		Ctx: ctx,
		Method: capnp.Method{
			InterfaceID:   0x8f9cac550b1bf41f,
			MethodID:      0,
			InterfaceName: "test.capnp:Adder",
			MethodName:    "add",
		},
		Options: capnp.NewCallOptions(opts),
	}
	if params != nil {
		call.ParamsSize = capnp.ObjectSize{DataSize: 8, PointerCount: 0}
		call.ParamsFunc = func(s capnp.Struct) error { return params(Adder_add_Params{Struct: s}) }
	}
	return Adder_add_Results_Promise{Pipeline: capnp.NewPipeline(c.Client.Call(call))}
}

type Adder_Server interface {
	Add(Adder_add) error
}

func Adder_ServerToClient(s Adder_Server) Adder {
	c, _ := s.(server.Closer)
	return Adder{Client: server.New(Adder_Methods(nil, s), c)}
}

func Adder_Methods(methods []server.Method, s Adder_Server) []server.Method {
	if cap(methods) == 0 {
		methods = make([]server.Method, 0, 1)
	}

	methods = append(methods, server.Method{
		Method: capnp.Method{
			InterfaceID:   0x8f9cac550b1bf41f,
			MethodID:      0,
			InterfaceName: "test.capnp:Adder",
			MethodName:    "add",
		},
		Impl: func(c context.Context, opts capnp.CallOptions, p, r capnp.Struct) error {
			call := Adder_add{c, opts, Adder_add_Params{Struct: p}, Adder_add_Results{Struct: r}}
			return s.Add(call)
		},
		ResultsSize: capnp.ObjectSize{DataSize: 8, PointerCount: 0},
	})

	return methods
}

// Adder_add holds the arguments for a server call to Adder.add.
type Adder_add struct {
	Ctx     context.Context
	Options capnp.CallOptions
	Params  Adder_add_Params
	Results Adder_add_Results
}

type Adder_add_Params struct{ capnp.Struct }

func NewAdder_add_Params(s *capnp.Segment) (Adder_add_Params, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	if err != nil {
		return Adder_add_Params{}, err
	}
	return Adder_add_Params{st}, nil
}

func NewRootAdder_add_Params(s *capnp.Segment) (Adder_add_Params, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	if err != nil {
		return Adder_add_Params{}, err
	}
	return Adder_add_Params{st}, nil
}

func ReadRootAdder_add_Params(msg *capnp.Message) (Adder_add_Params, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Adder_add_Params{}, err
	}
	return Adder_add_Params{root.Struct()}, nil
}
func (s Adder_add_Params) A() int32 {
	return int32(s.Struct.Uint32(0))
}

func (s Adder_add_Params) SetA(v int32) {
	s.Struct.SetUint32(0, uint32(v))
}

func (s Adder_add_Params) B() int32 {
	return int32(s.Struct.Uint32(4))
}

func (s Adder_add_Params) SetB(v int32) {
	s.Struct.SetUint32(4, uint32(v))
}

// Adder_add_Params_List is a list of Adder_add_Params.
type Adder_add_Params_List struct{ capnp.List }

// NewAdder_add_Params creates a new list of Adder_add_Params.
func NewAdder_add_Params_List(s *capnp.Segment, sz int32) (Adder_add_Params_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0}, sz)
	if err != nil {
		return Adder_add_Params_List{}, err
	}
	return Adder_add_Params_List{l}, nil
}

func (s Adder_add_Params_List) At(i int) Adder_add_Params { return Adder_add_Params{s.List.Struct(i)} }
func (s Adder_add_Params_List) Set(i int, v Adder_add_Params) error {
	return s.List.SetStruct(i, v.Struct)
}

// Adder_add_Params_Promise is a wrapper for a Adder_add_Params promised by a client call.
type Adder_add_Params_Promise struct{ *capnp.Pipeline }

func (p Adder_add_Params_Promise) Struct() (Adder_add_Params, error) {
	s, err := p.Pipeline.Struct()
	return Adder_add_Params{s}, err
}

type Adder_add_Results struct{ capnp.Struct }

func NewAdder_add_Results(s *capnp.Segment) (Adder_add_Results, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	if err != nil {
		return Adder_add_Results{}, err
	}
	return Adder_add_Results{st}, nil
}

func NewRootAdder_add_Results(s *capnp.Segment) (Adder_add_Results, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	if err != nil {
		return Adder_add_Results{}, err
	}
	return Adder_add_Results{st}, nil
}

func ReadRootAdder_add_Results(msg *capnp.Message) (Adder_add_Results, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Adder_add_Results{}, err
	}
	return Adder_add_Results{root.Struct()}, nil
}
func (s Adder_add_Results) Result() int32 {
	return int32(s.Struct.Uint32(0))
}

func (s Adder_add_Results) SetResult(v int32) {
	s.Struct.SetUint32(0, uint32(v))
}

// Adder_add_Results_List is a list of Adder_add_Results.
type Adder_add_Results_List struct{ capnp.List }

// NewAdder_add_Results creates a new list of Adder_add_Results.
func NewAdder_add_Results_List(s *capnp.Segment, sz int32) (Adder_add_Results_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0}, sz)
	if err != nil {
		return Adder_add_Results_List{}, err
	}
	return Adder_add_Results_List{l}, nil
}

func (s Adder_add_Results_List) At(i int) Adder_add_Results {
	return Adder_add_Results{s.List.Struct(i)}
}
func (s Adder_add_Results_List) Set(i int, v Adder_add_Results) error {
	return s.List.SetStruct(i, v.Struct)
}

// Adder_add_Results_Promise is a wrapper for a Adder_add_Results promised by a client call.
type Adder_add_Results_Promise struct{ *capnp.Pipeline }

func (p Adder_add_Results_Promise) Struct() (Adder_add_Results, error) {
	s, err := p.Pipeline.Struct()
	return Adder_add_Results{s}, err
}

const schema_ef12a34b9807e19c = "x\xda\x9cV\x7fh\x1c\xc5\x17\x9f\xd9\x9d\xed%\xa4\xf9" +
	"&\x93\xcd\xd744\xe5\x9a\x98h\x13\xc9]\x92\x0a\xe2" +
	"\x81\xbdTsV\x1b\x7fd/Tl\xc4\xe2\xe6n\xb8" +
	"\xa4\xde\xed^\xf76\xa6i)\xd1X\xb5\x8d\x16i\x83" +
	"`K\xad\xf8#\x94\x16\x0b\x0a\xa9\xd0?*(Dk" +
	"@1J\x82\xfeQH\x84R+J\xd2jA\x03f" +
	"};w{\xb7\x97k@\xfd\xe3\x96\xdbyo>\xef" +
	"}>\xef\xbd\x99m\xe9\x14\xda\x85V\xe9\xea\x1a\x84\x94" +
	".i\x8d5\xfdM\xe7\xd5\xdf.\xab/ Z\"Z" +
	"'\xe7=ov\xbeW\xb1\x80\x10\x96w\x8a\xc7dU" +
	"\xf4 $?-z\xec\x1fBV\xa8i\xfd\xfb\x9f?" +
	"Qu\xb0\xc09\x04\xce\x8fr\xe7\x87\xc5m\xf2\x1e\xb1" +
	"\x0a\x9c\xbf.Y\x1eX>}\xb4\xd0Y\x15\xa7\xe4\x04" +
	"w\xee\x07\xe778\xf2\x93[^\x19^(\xfe\xe3\x10" +
	"R\xfe\x8f1B\x04\xd66\x1f\x10wcp?,\x06" +
	"\x11\xb6\xb4w\x7f:\xdb\xf1\xfc\xdch\x01\xda\x19\x08\xfd" +
	"\x11G;\x07h3\x1c\xed\xd2\xbd\xdf\xef?8\x11z" +
	"\x15\xd1\x0a\x00\x93\xb0m\xbd(.\x81\xf7\xa7\x1c\xcc\xfb" +
	"\xfb\xfa\x92\x1d\x1f\x9c|\xbd\x00l^\x1c\x95\xafq\xb0" +
	"+\x00Vj\xa7a\x8d\x7f\xfbU\xe5\x8bS\x93\xc7\x0a" +
	"\x9co\x8ag\xe5\xbf\xb8\xf3\x9f\xe0\xdc\xc0\x9d\xf5\xdb\x1e" +
	"\x19\x9fU\xb7\x1cw\xf3(%\xbd6\x8fjb\x87\xde" +
	"U2\xbe8V5r\x1c\xd1j\xc7~\x1f\x09cD" +
	"\xac\x85\xd9\xc7\x84\xf3\xa7~8\x85\x14\xeal\x95\x1b\xc9" +
	"\x8f\xb0\xb3\x99\x0c\xc2\xce\xee\x1bw\xb2\xa1M\x1d\xa7\xdd" +
	"\xf6\xc3\xe4g\xb0\x1f\xe1\xc8\xfa\xe2Kkk\x9a\x95\x89" +
	"4in>Gn\x00\xf0\x87\xdb~\xdd\xf1T\xe8\xcb" +
	"\x0b.\xc3Q\xb2\x04\x86\xb9\xbb\x86n/\x1f>\xf3\x09" +
	"R\x8a%l\x1d\x8a_^Vj\x9a\xa6mj\x03d" +
	"T>`{v\xef%\"$\x9f3b\x11\x09r?" +
	"\xd9.'\xc8\xa0<I\xbcH\xb0^#oMRe" +
	"\xffL\x9a\x12W{\xf3$\xe9\xb19O\xf3\xcc\xcew" +
	"\xfc\xef\x0e\xfcq\xcb|a\x9c\xebdD\xbe\xc9\xe3," +
	"\xa6\xe3|\xf7\xce\xf4\xccTh\xf75w\xe5\xe6\x81\x05" +
	"\x96\xafp\xa8\xac\xfeX\x82<$\xa9G.\x96b\xf2" +
	"N\xe9\x1e\xc8\xc3d)\xd3\x17Q\x93XK\x06\x1eR" +
	"\xb5h\x1c\xb3.\x8c\xbbD\x09\x9ey\xc6P\xa4Og" +
	"\xd8\x80e\x85\x88\x12T\xcd\x89\x8a\x9d\xc6\xa1\xb4\x09\x09" +
	"T\xf2\x941pm\x077@\xc8\xf6\x01lp\xd0\x04" +
	"'\x14{P\x8d\x98\xba1\x84P\x0e\xd5)5v\x04" +
	"\xa24\x0c\xa8\xc5\x1eKc\x83|\x17\xc2\xac\x1d\xbb\xb3" +
	"#\x80\xf7\x80\x1a\x8f?nD\x99\xe1\x8b1\xd3~\xe9" +
	"f{\x06\x98\x16a\xf5a\x96*\x1b\x88\x9b)\xc0'" +
	"PH\x10\x88\x96V\xc0\x1c\x17\x89X\xa9\x140\xd6p" +
	"\x11\x12\xe0\x87WJ\x11\xcbc\xeb4\x0av\x1a#\xcb" +
	"\xb6\x0f\\\xf3\x13\x12\x1c\xb9\x0c\x9f-E}8\xc8R" +
	"+S\xa8\xcb\xa5\xe0\xb1cR\x97T\x18\xd3\x15\xe9l" +
	"\x8d\x025\xb7LN\xdfc\xa7\xc1)\xad\xe3\xe9x\xd4" +
	"h4?\x1b\xec\xc8Sf\xeb\xe3b\x94\x99:\xec\x1c" +
	"#\x94\x8e\xa4\x95v$\xc4\x19\x0d\xd1\xbf\xd1\xbbK5" +
	"<j\"\x8f\xebv\xe0\xba\x16\xb8\xae\x13\xb0\xc5\xf6&" +
	"Y\xc4dQh\x88\x02\xe1\xc9\xca\xc6\xf095\xe7\xb0" +
	"j\x02\xa7\xf24\xe6\xaa\xf8\x80p6hQ6h\xa3" +
	"]\xe3z\x08\xda\"`\x18\xfdJ{\xfai\xb3\xbd\xb8" +
	"\x09\x16\xef\x86\xc2\xab\x98@|\x02j\xf7:\xffVA" +
	"\xbfU\x0b\x05r\xf5\x0b\x1a\xbc\xbc\xb7\x04\xe1}d\xf8" +
	"\xec\x1e\xa9\xef\xf2\xda\x1cR\xab\xda\x9d6\xb1b\xbam" +
	"\xd6\x928\x90T#\xcf\xaa1\xc6+\x8f\xd7\xc2\xb8:" +
	"&\x90O\xa9\x81\x9adO#\xdaz\x7f\xee\xc8\xa0\xcd" +
	"\x01k\xd7\xd8\xdb\xca\xc5\xd9\xd1I\x90\xa2\xce\xfaba" +
	"\xaa\xbez\xc2\x1cG\xb4\xa1\xce\xaa\x98\x0b\xff2\xf4\xf2" +
	"s\x97\x10\xadm\xb3\xc66\xfa\xe7N\xb0\xf2%D7" +
	"\xf4X\xd7\x8f\xf8\xab*\x9e\xb9\xf0\x19\xbc4\x0dgb" +
	"\x07\xfb\x13I\xdd0=Q=\xe21\xd5\x98W\xd3\xe1" +
	"iE\x06R\xa6\x9e0\x87\x90\x98de\x9a\x9a`0" +
	"\xed\x82\xebl$\x02\xdeZn\xa7\x8c(n\xf3fR" +
	"\xfeG\x85\x0e\xdb\"\x88\xab\x8b\xdd\xc7\xdd`^\xb2\x97" +
	"pf^r\xaa\xa53\xcej\xb6\xdaXf\xea\xf1\xdf" +
	"\xa6\x12\x89ZRY\x07\x15\xc8}\x0b4\x06\\\xd7w" +
	"\x83\xe1\xba}\x1b\x02\xae\x0b\xb16\xec\xfa$\xa8\x0d\xb8" +
	"\xee\xd5\x0dm\xc1\xb4\x08\x96#\x0e\xf2ry\x82\xe96" +
	"\xb1\x9c\xa9C\xd8\x08\xa6\x99xy\x97*Ey\xe2\x97" +
	"\x83\xf8\x1b\xdd7H\x0d,\xb4g\xaa\x11\xe6\x14l\xa5" +
	"\x10\x86\x9a\xd8\x8b\xad\xf8\x04\xb6\xf6\xe9\x89\xde~\xb6\x8f" +
	"I\x9a/\xa2'\xfc1\xdd\xcf\xd54tSo\xf3\x1b" +
	"\xc9\x88\xbf_3\x99\xa1\xa9q\x7ff\xbfh\x97\xf4\xef" +
	"\x00\x00\x00\xff\xff\x11\xd4\xb7\xba"
