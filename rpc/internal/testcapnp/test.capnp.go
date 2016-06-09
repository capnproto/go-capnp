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

var schema_ef12a34b9807e19c = []byte{
	120, 218, 156, 86, 91, 76, 28, 101,
	20, 254, 207, 92, 186, 52, 160, 48,
	59, 40, 109, 138, 217, 22, 169, 22,
	12, 44, 160, 137, 233, 62, 116, 169,
	178, 98, 138, 23, 102, 73, 141, 69,
	109, 28, 118, 39, 91, 234, 238, 204,
	50, 204, 138, 219, 134, 96, 177, 104,
	69, 137, 17, 98, 34, 77, 169, 169,
	74, 154, 54, 62, 104, 66, 77, 120,
	168, 137, 38, 104, 75, 162, 17, 13,
	68, 30, 154, 80, 163, 169, 53, 54,
	84, 244, 65, 73, 202, 120, 254, 217,
	253, 217, 129, 213, 166, 241, 97, 129,
	157, 115, 206, 119, 190, 243, 157, 203,
	80, 215, 194, 53, 114, 245, 226, 149,
	13, 132, 40, 173, 226, 6, 123, 225,
	190, 244, 221, 37, 125, 103, 62, 35,
	202, 70, 17, 236, 99, 241, 75, 43,
	74, 121, 245, 12, 33, 32, 239, 227,
	7, 101, 149, 247, 16, 210, 246, 44,
	207, 3, 1, 123, 127, 225, 248, 245,
	145, 178, 254, 81, 34, 109, 6, 66,
	4, 180, 220, 255, 56, 31, 6, 34,
	216, 198, 157, 143, 141, 207, 169, 187,
	70, 137, 114, 7, 48, 211, 78, 190,
	3, 131, 228, 16, 31, 196, 208, 55,
	133, 177, 41, 73, 57, 60, 155, 9,
	21, 129, 58, 116, 242, 237, 212, 33,
	229, 56, 180, 45, 221, 171, 165, 119,
	52, 157, 38, 138, 196, 16, 228, 119,
	248, 95, 209, 62, 234, 216, 199, 191,
	251, 186, 244, 149, 233, 169, 97, 34,
	21, 242, 246, 137, 203, 158, 119, 91,
	62, 240, 46, 82, 150, 147, 252, 89,
	249, 115, 202, 82, 62, 207, 55, 203,
	87, 233, 95, 182, 254, 254, 47, 103,
	155, 94, 94, 24, 204, 115, 158, 225,
	135, 229, 121, 199, 121, 22, 157, 111,
	56, 206, 231, 154, 110, 191, 7, 62,
	173, 187, 156, 47, 192, 207, 124, 127,
	6, 176, 237, 167, 140, 0, 223, 159,
	154, 153, 157, 14, 29, 188, 74, 36,
	47, 171, 2, 129, 150, 208, 117, 222,
	33, 185, 56, 247, 4, 119, 238, 228,
	252, 73, 119, 17, 127, 243, 63, 162,
	253, 6, 223, 131, 246, 11, 59, 127,
	56, 124, 116, 34, 244, 134, 59, 126,
	159, 176, 140, 246, 231, 4, 26, 31,
	170, 222, 242, 225, 151, 79, 149, 29,
	205, 227, 221, 43, 12, 203, 3, 14,
	220, 17, 161, 89, 62, 35, 148, 33,
	239, 85, 59, 136, 132, 147, 71, 133,
	118, 121, 76, 136, 201, 215, 132, 7,
	9, 151, 171, 2, 120, 180, 85, 137,
	123, 228, 26, 177, 71, 30, 16, 125,
	104, 155, 249, 182, 229, 202, 31, 151,
	212, 35, 121, 73, 6, 196, 97, 121,
	72, 164, 73, 94, 23, 61, 244, 131,
	73, 190, 41, 92, 73, 173, 156, 126,
	59, 159, 81, 151, 56, 45, 247, 58,
	206, 105, 177, 89, 30, 119, 156, 141,
	235, 3, 69, 229, 53, 202, 68, 166,
	60, 135, 238, 144, 184, 132, 227, 241,
	244, 174, 215, 250, 22, 55, 254, 117,
	204, 61, 30, 93, 226, 65, 218, 253,
	94, 145, 22, 238, 251, 115, 75, 225,
	222, 143, 78, 188, 149, 151, 102, 76,
	28, 204, 128, 203, 167, 48, 205, 69,
	39, 205, 199, 205, 215, 246, 62, 19,
	186, 56, 233, 74, 243, 137, 184, 140,
	105, 98, 70, 109, 68, 77, 234, 73,
	8, 36, 213, 200, 11, 106, 76, 35,
	164, 21, 0, 138, 176, 104, 75, 235,
	182, 168, 81, 208, 147, 129, 71, 85,
	61, 26, 215, 30, 81, 35, 150, 97,
	166, 107, 117, 173, 39, 243, 160, 178,
	85, 53, 213, 4, 116, 175, 241, 125,
	88, 141, 199, 159, 52, 163, 154, 89,
	27, 211, 44, 250, 165, 77, 235, 74,
	105, 122, 196, 241, 246, 168, 137, 110,
	69, 224, 5, 36, 129, 68, 164, 219,
	246, 224, 74, 21, 241, 160, 108, 226,
	192, 214, 94, 74, 106, 17, 75, 139,
	18, 172, 162, 128, 112, 248, 129, 91,
	34, 17, 214, 186, 83, 113, 222, 90,
	131, 27, 64, 220, 2, 196, 45, 229,
	32, 120, 192, 113, 3, 41, 215, 70,
	2, 32, 185, 192, 57, 4, 223, 29,
	165, 140, 213, 104, 148, 194, 21, 167,
	226, 255, 13, 103, 210, 116, 22, 8,
	200, 80, 112, 129, 0, 43, 189, 152,
	214, 142, 34, 98, 188, 72, 59, 156,
	221, 115, 96, 29, 149, 164, 126, 194,
	73, 27, 61, 54, 147, 7, 178, 250,
	144, 70, 192, 176, 53, 128, 88, 97,
	76, 3, 55, 90, 118, 94, 128, 117,
	84, 146, 170, 17, 77, 244, 20, 99,
	149, 177, 12, 64, 174, 167, 157, 137,
	164, 97, 90, 121, 29, 165, 245, 134,
	34, 7, 12, 44, 88, 195, 95, 149,
	173, 62, 218, 198, 53, 5, 87, 228,
	10, 246, 80, 46, 82, 238, 154, 220,
	84, 60, 214, 226, 130, 85, 168, 42,
	47, 66, 85, 34, 84, 29, 7, 184,
	226, 165, 116, 154, 165, 26, 250, 112,
	7, 62, 124, 128, 3, 80, 153, 150,
	208, 145, 167, 234, 122, 170, 225, 160,
	163, 254, 255, 228, 10, 12, 205, 45,
	41, 187, 80, 192, 78, 205, 170, 164,
	52, 99, 35, 186, 129, 235, 150, 98,
	0, 67, 35, 188, 158, 84, 54, 161,
	49, 119, 30, 170, 2, 174, 245, 223,
	110, 186, 174, 234, 246, 128, 235, 30,
	111, 11, 187, 238, 214, 182, 128, 107,
	151, 239, 106, 8, 102, 166, 218, 102,
	211, 78, 124, 206, 188, 7, 157, 81,
	48, 109, 182, 93, 4, 204, 96, 70,
	24, 159, 163, 190, 82, 0, 156, 235,
	197, 84, 194, 193, 238, 173, 224, 58,
	212, 229, 248, 160, 145, 142, 1, 145,
	32, 236, 148, 64, 103, 132, 64, 18,
	53, 194, 135, 245, 112, 28, 236, 67,
	70, 162, 163, 83, 59, 164, 137, 122,
	109, 196, 72, 248, 99, 134, 223, 153,
	35, 211, 176, 140, 6, 191, 153, 140,
	248, 59, 117, 75, 51, 117, 53, 238,
	207, 198, 163, 0, 168, 7, 155, 55,
	212, 70, 41, 71, 57, 86, 89, 72,
	245, 15, 229, 24, 72, 53, 1, 123,
	255, 200, 123, 202, 249, 185, 193, 41,
	20, 170, 194, 254, 106, 113, 186, 114,
	243, 132, 53, 142, 218, 84, 216, 222,
	133, 240, 111, 233, 87, 95, 188, 128,
	114, 52, 216, 35, 91, 253, 11, 199,
	181, 146, 101, 148, 163, 221, 254, 125,
	200, 95, 230, 125, 126, 242, 11, 252,
	82, 221, 151, 61, 82, 193, 204, 96,
	123, 162, 70, 196, 99, 169, 49, 159,
	110, 224, 79, 59, 146, 234, 182, 140,
	132, 149, 38, 124, 82, 43, 214, 213,
	132, 134, 157, 115, 107, 34, 160, 4,
	37, 89, 9, 26, 124, 89, 202, 235,
	215, 45, 26, 7, 13, 103, 163, 149,
	23, 221, 187, 200, 173, 63, 63, 206,
	153, 204, 14, 16, 123, 199, 3, 123,
	99, 75, 82, 56, 179, 225, 236, 70,
	17, 208, 214, 238, 54, 199, 118, 219,
	172, 165, 123, 203, 54, 240, 22, 239,
	232, 191, 93, 40, 111, 110, 9, 64,
	207, 59, 159, 192, 150, 212, 205, 155,
	189, 122, 129, 253, 35, 33, 73, 21,
	206, 224, 123, 112, 147, 111, 206, 151,
	173, 225, 63, 1, 0, 0, 255, 255,
	247, 222, 187, 47,
}
