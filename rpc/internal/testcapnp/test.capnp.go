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
	NewHandle(ctx context.Context, opts C.CallOptions, params HandleFactory_newHandle_Params, results HandleFactory_newHandle_Results) error
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
			return server.NewHandle(c, opts, HandleFactory_newHandle_Params(p), HandleFactory_newHandle_Results(r))
		},
		ResultsSize: C.ObjectSize{DataSize: 0, PointerCount: 1},
	})

	return methods
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
