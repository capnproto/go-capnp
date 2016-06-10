package persistent

// AUTO GENERATED - DO NOT EDIT

import (
	context "golang.org/x/net/context"
	capnp "zombiezen.com/go/capnproto2"
	server "zombiezen.com/go/capnproto2/server"
)

const PersistentAnnotation = uint64(0xf622595091cafb67)

type Persistent struct{ Client capnp.Client }

func (c Persistent) Save(ctx context.Context, params func(Persistent_SaveParams) error, opts ...capnp.CallOption) Persistent_SaveResults_Promise {
	if c.Client == nil {
		return Persistent_SaveResults_Promise{Pipeline: capnp.NewPipeline(capnp.ErrorAnswer(capnp.ErrNullClient))}
	}
	call := &capnp.Call{
		Ctx: ctx,
		Method: capnp.Method{
			InterfaceID:   0xc8cb212fcd9f5691,
			MethodID:      0,
			InterfaceName: "persistent.capnp:Persistent",
			MethodName:    "save",
		},
		Options: capnp.NewCallOptions(opts),
	}
	if params != nil {
		call.ParamsSize = capnp.ObjectSize{DataSize: 0, PointerCount: 1}
		call.ParamsFunc = func(s capnp.Struct) error { return params(Persistent_SaveParams{Struct: s}) }
	}
	return Persistent_SaveResults_Promise{Pipeline: capnp.NewPipeline(c.Client.Call(call))}
}

type Persistent_Server interface {
	Save(Persistent_save) error
}

func Persistent_ServerToClient(s Persistent_Server) Persistent {
	c, _ := s.(server.Closer)
	return Persistent{Client: server.New(Persistent_Methods(nil, s), c)}
}

func Persistent_Methods(methods []server.Method, s Persistent_Server) []server.Method {
	if cap(methods) == 0 {
		methods = make([]server.Method, 0, 1)
	}

	methods = append(methods, server.Method{
		Method: capnp.Method{
			InterfaceID:   0xc8cb212fcd9f5691,
			MethodID:      0,
			InterfaceName: "persistent.capnp:Persistent",
			MethodName:    "save",
		},
		Impl: func(c context.Context, opts capnp.CallOptions, p, r capnp.Struct) error {
			call := Persistent_save{c, opts, Persistent_SaveParams{Struct: p}, Persistent_SaveResults{Struct: r}}
			return s.Save(call)
		},
		ResultsSize: capnp.ObjectSize{DataSize: 0, PointerCount: 1},
	})

	return methods
}

// Persistent_save holds the arguments for a server call to Persistent.save.
type Persistent_save struct {
	Ctx     context.Context
	Options capnp.CallOptions
	Params  Persistent_SaveParams
	Results Persistent_SaveResults
}

type Persistent_SaveParams struct{ capnp.Struct }

func NewPersistent_SaveParams(s *capnp.Segment) (Persistent_SaveParams, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Persistent_SaveParams{st}, err
}

func NewRootPersistent_SaveParams(s *capnp.Segment) (Persistent_SaveParams, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Persistent_SaveParams{st}, err
}

func ReadRootPersistent_SaveParams(msg *capnp.Message) (Persistent_SaveParams, error) {
	root, err := msg.RootPtr()
	return Persistent_SaveParams{root.Struct()}, err
}
func (s Persistent_SaveParams) SealFor() (capnp.Pointer, error) {
	return s.Struct.Pointer(0)
}

func (s Persistent_SaveParams) HasSealFor() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Persistent_SaveParams) SealForPtr() (capnp.Ptr, error) {
	return s.Struct.Ptr(0)
}

func (s Persistent_SaveParams) SetSealFor(v capnp.Pointer) error {
	return s.Struct.SetPointer(0, v)
}

func (s Persistent_SaveParams) SetSealForPtr(v capnp.Ptr) error {
	return s.Struct.SetPtr(0, v)
}

// Persistent_SaveParams_List is a list of Persistent_SaveParams.
type Persistent_SaveParams_List struct{ capnp.List }

// NewPersistent_SaveParams creates a new list of Persistent_SaveParams.
func NewPersistent_SaveParams_List(s *capnp.Segment, sz int32) (Persistent_SaveParams_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return Persistent_SaveParams_List{l}, err
}

func (s Persistent_SaveParams_List) At(i int) Persistent_SaveParams {
	return Persistent_SaveParams{s.List.Struct(i)}
}
func (s Persistent_SaveParams_List) Set(i int, v Persistent_SaveParams) error {
	return s.List.SetStruct(i, v.Struct)
}

// Persistent_SaveParams_Promise is a wrapper for a Persistent_SaveParams promised by a client call.
type Persistent_SaveParams_Promise struct{ *capnp.Pipeline }

func (p Persistent_SaveParams_Promise) Struct() (Persistent_SaveParams, error) {
	s, err := p.Pipeline.Struct()
	return Persistent_SaveParams{s}, err
}

func (p Persistent_SaveParams_Promise) SealFor() *capnp.Pipeline {
	return p.Pipeline.GetPipeline(0)
}

type Persistent_SaveResults struct{ capnp.Struct }

func NewPersistent_SaveResults(s *capnp.Segment) (Persistent_SaveResults, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Persistent_SaveResults{st}, err
}

func NewRootPersistent_SaveResults(s *capnp.Segment) (Persistent_SaveResults, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Persistent_SaveResults{st}, err
}

func ReadRootPersistent_SaveResults(msg *capnp.Message) (Persistent_SaveResults, error) {
	root, err := msg.RootPtr()
	return Persistent_SaveResults{root.Struct()}, err
}
func (s Persistent_SaveResults) SturdyRef() (capnp.Pointer, error) {
	return s.Struct.Pointer(0)
}

func (s Persistent_SaveResults) HasSturdyRef() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Persistent_SaveResults) SturdyRefPtr() (capnp.Ptr, error) {
	return s.Struct.Ptr(0)
}

func (s Persistent_SaveResults) SetSturdyRef(v capnp.Pointer) error {
	return s.Struct.SetPointer(0, v)
}

func (s Persistent_SaveResults) SetSturdyRefPtr(v capnp.Ptr) error {
	return s.Struct.SetPtr(0, v)
}

// Persistent_SaveResults_List is a list of Persistent_SaveResults.
type Persistent_SaveResults_List struct{ capnp.List }

// NewPersistent_SaveResults creates a new list of Persistent_SaveResults.
func NewPersistent_SaveResults_List(s *capnp.Segment, sz int32) (Persistent_SaveResults_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return Persistent_SaveResults_List{l}, err
}

func (s Persistent_SaveResults_List) At(i int) Persistent_SaveResults {
	return Persistent_SaveResults{s.List.Struct(i)}
}
func (s Persistent_SaveResults_List) Set(i int, v Persistent_SaveResults) error {
	return s.List.SetStruct(i, v.Struct)
}

// Persistent_SaveResults_Promise is a wrapper for a Persistent_SaveResults promised by a client call.
type Persistent_SaveResults_Promise struct{ *capnp.Pipeline }

func (p Persistent_SaveResults_Promise) Struct() (Persistent_SaveResults, error) {
	s, err := p.Pipeline.Struct()
	return Persistent_SaveResults{s}, err
}

func (p Persistent_SaveResults_Promise) SturdyRef() *capnp.Pipeline {
	return p.Pipeline.GetPipeline(0)
}

type RealmGateway struct{ Client capnp.Client }

func (c RealmGateway) Import(ctx context.Context, params func(RealmGateway_import_Params) error, opts ...capnp.CallOption) Persistent_SaveResults_Promise {
	if c.Client == nil {
		return Persistent_SaveResults_Promise{Pipeline: capnp.NewPipeline(capnp.ErrorAnswer(capnp.ErrNullClient))}
	}
	call := &capnp.Call{
		Ctx: ctx,
		Method: capnp.Method{
			InterfaceID:   0x84ff286cd00a3ed4,
			MethodID:      0,
			InterfaceName: "persistent.capnp:RealmGateway",
			MethodName:    "import",
		},
		Options: capnp.NewCallOptions(opts),
	}
	if params != nil {
		call.ParamsSize = capnp.ObjectSize{DataSize: 0, PointerCount: 2}
		call.ParamsFunc = func(s capnp.Struct) error { return params(RealmGateway_import_Params{Struct: s}) }
	}
	return Persistent_SaveResults_Promise{Pipeline: capnp.NewPipeline(c.Client.Call(call))}
}
func (c RealmGateway) Export(ctx context.Context, params func(RealmGateway_export_Params) error, opts ...capnp.CallOption) Persistent_SaveResults_Promise {
	if c.Client == nil {
		return Persistent_SaveResults_Promise{Pipeline: capnp.NewPipeline(capnp.ErrorAnswer(capnp.ErrNullClient))}
	}
	call := &capnp.Call{
		Ctx: ctx,
		Method: capnp.Method{
			InterfaceID:   0x84ff286cd00a3ed4,
			MethodID:      1,
			InterfaceName: "persistent.capnp:RealmGateway",
			MethodName:    "export",
		},
		Options: capnp.NewCallOptions(opts),
	}
	if params != nil {
		call.ParamsSize = capnp.ObjectSize{DataSize: 0, PointerCount: 2}
		call.ParamsFunc = func(s capnp.Struct) error { return params(RealmGateway_export_Params{Struct: s}) }
	}
	return Persistent_SaveResults_Promise{Pipeline: capnp.NewPipeline(c.Client.Call(call))}
}

type RealmGateway_Server interface {
	Import(RealmGateway_import) error

	Export(RealmGateway_export) error
}

func RealmGateway_ServerToClient(s RealmGateway_Server) RealmGateway {
	c, _ := s.(server.Closer)
	return RealmGateway{Client: server.New(RealmGateway_Methods(nil, s), c)}
}

func RealmGateway_Methods(methods []server.Method, s RealmGateway_Server) []server.Method {
	if cap(methods) == 0 {
		methods = make([]server.Method, 0, 2)
	}

	methods = append(methods, server.Method{
		Method: capnp.Method{
			InterfaceID:   0x84ff286cd00a3ed4,
			MethodID:      0,
			InterfaceName: "persistent.capnp:RealmGateway",
			MethodName:    "import",
		},
		Impl: func(c context.Context, opts capnp.CallOptions, p, r capnp.Struct) error {
			call := RealmGateway_import{c, opts, RealmGateway_import_Params{Struct: p}, Persistent_SaveResults{Struct: r}}
			return s.Import(call)
		},
		ResultsSize: capnp.ObjectSize{DataSize: 0, PointerCount: 1},
	})

	methods = append(methods, server.Method{
		Method: capnp.Method{
			InterfaceID:   0x84ff286cd00a3ed4,
			MethodID:      1,
			InterfaceName: "persistent.capnp:RealmGateway",
			MethodName:    "export",
		},
		Impl: func(c context.Context, opts capnp.CallOptions, p, r capnp.Struct) error {
			call := RealmGateway_export{c, opts, RealmGateway_export_Params{Struct: p}, Persistent_SaveResults{Struct: r}}
			return s.Export(call)
		},
		ResultsSize: capnp.ObjectSize{DataSize: 0, PointerCount: 1},
	})

	return methods
}

// RealmGateway_import holds the arguments for a server call to RealmGateway.import.
type RealmGateway_import struct {
	Ctx     context.Context
	Options capnp.CallOptions
	Params  RealmGateway_import_Params
	Results Persistent_SaveResults
}

// RealmGateway_export holds the arguments for a server call to RealmGateway.export.
type RealmGateway_export struct {
	Ctx     context.Context
	Options capnp.CallOptions
	Params  RealmGateway_export_Params
	Results Persistent_SaveResults
}

type RealmGateway_import_Params struct{ capnp.Struct }

func NewRealmGateway_import_Params(s *capnp.Segment) (RealmGateway_import_Params, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	return RealmGateway_import_Params{st}, err
}

func NewRootRealmGateway_import_Params(s *capnp.Segment) (RealmGateway_import_Params, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	return RealmGateway_import_Params{st}, err
}

func ReadRootRealmGateway_import_Params(msg *capnp.Message) (RealmGateway_import_Params, error) {
	root, err := msg.RootPtr()
	return RealmGateway_import_Params{root.Struct()}, err
}
func (s RealmGateway_import_Params) Cap() Persistent {
	p, _ := s.Struct.Ptr(0)
	return Persistent{Client: p.Interface().Client()}
}

func (s RealmGateway_import_Params) HasCap() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s RealmGateway_import_Params) SetCap(v Persistent) error {
	if v.Client == nil {
		return s.Struct.SetPtr(0, capnp.Ptr{})
	}
	seg := s.Segment()
	in := capnp.NewInterface(seg, seg.Message().AddCap(v.Client))
	return s.Struct.SetPtr(0, in.ToPtr())
}

func (s RealmGateway_import_Params) Params() (Persistent_SaveParams, error) {
	p, err := s.Struct.Ptr(1)
	return Persistent_SaveParams{Struct: p.Struct()}, err
}

func (s RealmGateway_import_Params) HasParams() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s RealmGateway_import_Params) SetParams(v Persistent_SaveParams) error {
	return s.Struct.SetPtr(1, v.Struct.ToPtr())
}

// NewParams sets the params field to a newly
// allocated Persistent_SaveParams struct, preferring placement in s's segment.
func (s RealmGateway_import_Params) NewParams() (Persistent_SaveParams, error) {
	ss, err := NewPersistent_SaveParams(s.Struct.Segment())
	if err != nil {
		return Persistent_SaveParams{}, err
	}
	err = s.Struct.SetPtr(1, ss.Struct.ToPtr())
	return ss, err
}

// RealmGateway_import_Params_List is a list of RealmGateway_import_Params.
type RealmGateway_import_Params_List struct{ capnp.List }

// NewRealmGateway_import_Params creates a new list of RealmGateway_import_Params.
func NewRealmGateway_import_Params_List(s *capnp.Segment, sz int32) (RealmGateway_import_Params_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2}, sz)
	return RealmGateway_import_Params_List{l}, err
}

func (s RealmGateway_import_Params_List) At(i int) RealmGateway_import_Params {
	return RealmGateway_import_Params{s.List.Struct(i)}
}
func (s RealmGateway_import_Params_List) Set(i int, v RealmGateway_import_Params) error {
	return s.List.SetStruct(i, v.Struct)
}

// RealmGateway_import_Params_Promise is a wrapper for a RealmGateway_import_Params promised by a client call.
type RealmGateway_import_Params_Promise struct{ *capnp.Pipeline }

func (p RealmGateway_import_Params_Promise) Struct() (RealmGateway_import_Params, error) {
	s, err := p.Pipeline.Struct()
	return RealmGateway_import_Params{s}, err
}

func (p RealmGateway_import_Params_Promise) Cap() Persistent {
	return Persistent{Client: p.Pipeline.GetPipeline(0).Client()}
}

func (p RealmGateway_import_Params_Promise) Params() Persistent_SaveParams_Promise {
	return Persistent_SaveParams_Promise{Pipeline: p.Pipeline.GetPipeline(1)}
}

type RealmGateway_export_Params struct{ capnp.Struct }

func NewRealmGateway_export_Params(s *capnp.Segment) (RealmGateway_export_Params, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	return RealmGateway_export_Params{st}, err
}

func NewRootRealmGateway_export_Params(s *capnp.Segment) (RealmGateway_export_Params, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	return RealmGateway_export_Params{st}, err
}

func ReadRootRealmGateway_export_Params(msg *capnp.Message) (RealmGateway_export_Params, error) {
	root, err := msg.RootPtr()
	return RealmGateway_export_Params{root.Struct()}, err
}
func (s RealmGateway_export_Params) Cap() Persistent {
	p, _ := s.Struct.Ptr(0)
	return Persistent{Client: p.Interface().Client()}
}

func (s RealmGateway_export_Params) HasCap() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s RealmGateway_export_Params) SetCap(v Persistent) error {
	if v.Client == nil {
		return s.Struct.SetPtr(0, capnp.Ptr{})
	}
	seg := s.Segment()
	in := capnp.NewInterface(seg, seg.Message().AddCap(v.Client))
	return s.Struct.SetPtr(0, in.ToPtr())
}

func (s RealmGateway_export_Params) Params() (Persistent_SaveParams, error) {
	p, err := s.Struct.Ptr(1)
	return Persistent_SaveParams{Struct: p.Struct()}, err
}

func (s RealmGateway_export_Params) HasParams() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s RealmGateway_export_Params) SetParams(v Persistent_SaveParams) error {
	return s.Struct.SetPtr(1, v.Struct.ToPtr())
}

// NewParams sets the params field to a newly
// allocated Persistent_SaveParams struct, preferring placement in s's segment.
func (s RealmGateway_export_Params) NewParams() (Persistent_SaveParams, error) {
	ss, err := NewPersistent_SaveParams(s.Struct.Segment())
	if err != nil {
		return Persistent_SaveParams{}, err
	}
	err = s.Struct.SetPtr(1, ss.Struct.ToPtr())
	return ss, err
}

// RealmGateway_export_Params_List is a list of RealmGateway_export_Params.
type RealmGateway_export_Params_List struct{ capnp.List }

// NewRealmGateway_export_Params creates a new list of RealmGateway_export_Params.
func NewRealmGateway_export_Params_List(s *capnp.Segment, sz int32) (RealmGateway_export_Params_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2}, sz)
	return RealmGateway_export_Params_List{l}, err
}

func (s RealmGateway_export_Params_List) At(i int) RealmGateway_export_Params {
	return RealmGateway_export_Params{s.List.Struct(i)}
}
func (s RealmGateway_export_Params_List) Set(i int, v RealmGateway_export_Params) error {
	return s.List.SetStruct(i, v.Struct)
}

// RealmGateway_export_Params_Promise is a wrapper for a RealmGateway_export_Params promised by a client call.
type RealmGateway_export_Params_Promise struct{ *capnp.Pipeline }

func (p RealmGateway_export_Params_Promise) Struct() (RealmGateway_export_Params, error) {
	s, err := p.Pipeline.Struct()
	return RealmGateway_export_Params{s}, err
}

func (p RealmGateway_export_Params_Promise) Cap() Persistent {
	return Persistent{Client: p.Pipeline.GetPipeline(0).Client()}
}

func (p RealmGateway_export_Params_Promise) Params() Persistent_SaveParams_Promise {
	return Persistent_SaveParams_Promise{Pipeline: p.Pipeline.GetPipeline(1)}
}
