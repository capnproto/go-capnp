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
	if err != nil {
		return Persistent_SaveParams{}, err
	}
	return Persistent_SaveParams{st}, nil
}

func NewRootPersistent_SaveParams(s *capnp.Segment) (Persistent_SaveParams, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Persistent_SaveParams{}, err
	}
	return Persistent_SaveParams{st}, nil
}

func ReadRootPersistent_SaveParams(msg *capnp.Message) (Persistent_SaveParams, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Persistent_SaveParams{}, err
	}
	return Persistent_SaveParams{root.Struct()}, nil
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
	if err != nil {
		return Persistent_SaveParams_List{}, err
	}
	return Persistent_SaveParams_List{l}, nil
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
	if err != nil {
		return Persistent_SaveResults{}, err
	}
	return Persistent_SaveResults{st}, nil
}

func NewRootPersistent_SaveResults(s *capnp.Segment) (Persistent_SaveResults, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Persistent_SaveResults{}, err
	}
	return Persistent_SaveResults{st}, nil
}

func ReadRootPersistent_SaveResults(msg *capnp.Message) (Persistent_SaveResults, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Persistent_SaveResults{}, err
	}
	return Persistent_SaveResults{root.Struct()}, nil
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
	if err != nil {
		return Persistent_SaveResults_List{}, err
	}
	return Persistent_SaveResults_List{l}, nil
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
	if err != nil {
		return RealmGateway_import_Params{}, err
	}
	return RealmGateway_import_Params{st}, nil
}

func NewRootRealmGateway_import_Params(s *capnp.Segment) (RealmGateway_import_Params, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	if err != nil {
		return RealmGateway_import_Params{}, err
	}
	return RealmGateway_import_Params{st}, nil
}

func ReadRootRealmGateway_import_Params(msg *capnp.Message) (RealmGateway_import_Params, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return RealmGateway_import_Params{}, err
	}
	return RealmGateway_import_Params{root.Struct()}, nil
}
func (s RealmGateway_import_Params) Cap() Persistent {
	p, err := s.Struct.Ptr(0)
	if err != nil {

		return Persistent{}
	}
	return Persistent{Client: p.Interface().Client()}
}

func (s RealmGateway_import_Params) HasCap() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s RealmGateway_import_Params) SetCap(v Persistent) error {
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

func (s RealmGateway_import_Params) Params() (Persistent_SaveParams, error) {
	p, err := s.Struct.Ptr(1)
	if err != nil {
		return Persistent_SaveParams{}, err
	}
	return Persistent_SaveParams{Struct: p.Struct()}, nil
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
	if err != nil {
		return RealmGateway_import_Params_List{}, err
	}
	return RealmGateway_import_Params_List{l}, nil
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
	if err != nil {
		return RealmGateway_export_Params{}, err
	}
	return RealmGateway_export_Params{st}, nil
}

func NewRootRealmGateway_export_Params(s *capnp.Segment) (RealmGateway_export_Params, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	if err != nil {
		return RealmGateway_export_Params{}, err
	}
	return RealmGateway_export_Params{st}, nil
}

func ReadRootRealmGateway_export_Params(msg *capnp.Message) (RealmGateway_export_Params, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return RealmGateway_export_Params{}, err
	}
	return RealmGateway_export_Params{root.Struct()}, nil
}
func (s RealmGateway_export_Params) Cap() Persistent {
	p, err := s.Struct.Ptr(0)
	if err != nil {

		return Persistent{}
	}
	return Persistent{Client: p.Interface().Client()}
}

func (s RealmGateway_export_Params) HasCap() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s RealmGateway_export_Params) SetCap(v Persistent) error {
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

func (s RealmGateway_export_Params) Params() (Persistent_SaveParams, error) {
	p, err := s.Struct.Ptr(1)
	if err != nil {
		return Persistent_SaveParams{}, err
	}
	return Persistent_SaveParams{Struct: p.Struct()}, nil
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
	if err != nil {
		return RealmGateway_export_Params_List{}, err
	}
	return RealmGateway_export_Params_List{l}, nil
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

var schema_b8630836983feed7 = []byte{
	120, 218, 188, 86, 107, 108, 20, 213,
	23, 191, 231, 206, 78, 167, 203, 159,
	254, 187, 183, 3, 1, 146, 226, 66,
	45, 22, 202, 99, 219, 138, 138, 27,
	113, 187, 4, 4, 76, 148, 206, 214,
	104, 104, 226, 99, 186, 29, 106, 181,
	59, 179, 238, 76, 233, 46, 134, 240,
	136, 175, 20, 253, 64, 248, 162, 36,
	74, 140, 33, 70, 34, 49, 38, 62,
	192, 4, 137, 68, 124, 32, 53, 214,
	24, 249, 100, 130, 124, 51, 33, 34,
	36, 248, 76, 184, 158, 59, 143, 157,
	105, 183, 161, 49, 24, 63, 236, 38,
	115, 238, 220, 115, 206, 239, 119, 126,
	231, 156, 233, 216, 71, 187, 105, 167,
	124, 46, 70, 136, 182, 70, 174, 227,
	135, 237, 142, 249, 91, 63, 178, 126,
	35, 172, 25, 248, 254, 7, 15, 141,
	167, 22, 127, 245, 5, 145, 65, 73,
	192, 173, 26, 77, 131, 170, 83, 133,
	16, 245, 97, 154, 33, 192, 207, 47,
	175, 220, 156, 216, 245, 214, 199, 68,
	139, 203, 192, 95, 24, 254, 225, 154,
	214, 220, 62, 65, 8, 168, 59, 233,
	152, 250, 172, 120, 179, 119, 55, 149,
	128, 68, 14, 65, 34, 84, 125, 138,
	222, 171, 142, 208, 81, 117, 156, 38,
	9, 229, 151, 95, 74, 205, 107, 122,
	236, 248, 41, 50, 17, 151, 175, 53,
	78, 242, 51, 78, 75, 234, 132, 235,
	231, 172, 231, 231, 187, 187, 103, 125,
	51, 188, 148, 63, 67, 24, 147, 248,
	185, 159, 51, 47, 223, 94, 159, 63,
	70, 8, 73, 128, 122, 130, 94, 81,
	191, 116, 179, 59, 77, 159, 87, 199,
	37, 5, 127, 109, 124, 240, 175, 51,
	251, 123, 182, 182, 252, 74, 38, 152,
	188, 8, 34, 87, 64, 253, 73, 186,
	160, 94, 197, 215, 174, 74, 201, 222,
	230, 152, 235, 254, 200, 27, 43, 55,
	237, 123, 253, 157, 139, 132, 221, 4,
	132, 200, 84, 192, 110, 136, 245, 227,
	145, 186, 32, 54, 74, 34, 247, 225,
	127, 136, 67, 150, 199, 212, 6, 249,
	14, 117, 173, 188, 5, 113, 236, 185,
	191, 227, 216, 211, 59, 127, 63, 65,
	96, 22, 158, 189, 34, 239, 85, 95,
	149, 219, 212, 227, 238, 217, 201, 75,
	221, 47, 126, 178, 233, 241, 15, 167,
	227, 117, 110, 221, 58, 80, 151, 212,
	137, 204, 23, 215, 9, 94, 239, 123,
	40, 126, 231, 194, 179, 167, 126, 137,
	102, 177, 161, 206, 205, 66, 171, 19,
	89, 172, 120, 116, 247, 201, 215, 254,
	248, 236, 56, 209, 18, 72, 124, 53,
	46, 158, 175, 85, 190, 85, 55, 43,
	130, 176, 245, 138, 139, 168, 26, 171,
	150, 176, 219, 148, 11, 106, 86, 105,
	19, 229, 84, 54, 170, 71, 21, 5,
	127, 243, 248, 251, 235, 255, 127, 11,
	124, 208, 241, 99, 109, 81, 143, 42,
	123, 213, 119, 93, 223, 111, 123, 190,
	139, 70, 201, 30, 178, 29, 67, 50,
	157, 85, 121, 189, 104, 22, 211, 61,
	190, 5, 13, 189, 250, 118, 35, 211,
	163, 151, 244, 130, 221, 3, 160, 197,
	36, 148, 88, 12, 225, 176, 134, 117,
	168, 181, 122, 9, 180, 57, 20, 118,
	217, 134, 62, 124, 143, 85, 130, 38,
	25, 34, 185, 98, 180, 38, 244, 63,
	104, 121, 110, 33, 93, 212, 243, 79,
	234, 131, 6, 33, 232, 10, 102, 35,
	159, 193, 17, 190, 170, 53, 67, 68,
	138, 172, 115, 93, 8, 129, 173, 76,
	243, 71, 14, 28, 210, 78, 124, 63,
	118, 154, 176, 101, 45, 252, 243, 75,
	103, 90, 23, 188, 231, 28, 38, 108,
	73, 11, 111, 58, 159, 187, 88, 121,
	110, 59, 82, 179, 184, 139, 31, 88,
	148, 58, 127, 208, 72, 252, 73, 216,
	194, 190, 80, 145, 108, 97, 251, 46,
	63, 118, 102, 168, 80, 180, 74, 142,
	50, 96, 229, 21, 71, 31, 76, 154,
	22, 254, 243, 252, 136, 237, 88, 5,
	167, 66, 164, 162, 209, 104, 234, 5,
	67, 139, 1, 141, 52, 70, 140, 66,
	54, 33, 82, 38, 12, 186, 146, 126,
	202, 97, 246, 201, 180, 184, 19, 128,
	10, 248, 164, 85, 62, 115, 72, 79,
	97, 163, 158, 116, 140, 81, 189, 34,
	104, 172, 151, 100, 188, 31, 8, 4,
	2, 89, 177, 206, 52, 33, 217, 21,
	144, 189, 11, 216, 30, 5, 66, 33,
	135, 111, 140, 136, 55, 138, 144, 221,
	13, 236, 180, 226, 131, 97, 144, 20,
	25, 134, 93, 5, 162, 64, 190, 49,
	82, 13, 6, 109, 90, 61, 38, 41,
	226, 139, 255, 4, 254, 55, 65, 228,
	158, 40, 152, 140, 192, 35, 134, 110,
	200, 24, 229, 27, 14, 34, 215, 70,
	145, 38, 71, 65, 86, 144, 97, 214,
	208, 207, 88, 63, 155, 91, 98, 11,
	74, 124, 179, 233, 24, 37, 83, 31,
	38, 74, 206, 216, 198, 55, 148, 163,
	79, 213, 179, 228, 150, 81, 211, 40,
	133, 167, 254, 115, 109, 13, 124, 139,
	130, 162, 118, 133, 12, 145, 137, 229,
	214, 119, 190, 95, 223, 35, 60, 144,
	63, 152, 78, 214, 20, 2, 105, 116,
	134, 44, 147, 196, 170, 94, 99, 83,
	43, 235, 22, 118, 149, 71, 84, 43,
	118, 139, 130, 237, 130, 248, 131, 94,
	89, 214, 130, 250, 110, 197, 94, 41,
	83, 96, 0, 115, 92, 234, 68, 37,
	181, 34, 26, 63, 165, 160, 160, 55,
	96, 81, 26, 187, 225, 95, 42, 32,
	48, 2, 153, 162, 219, 192, 144, 8,
	23, 195, 63, 137, 48, 83, 245, 32,
	17, 153, 34, 16, 112, 131, 232, 102,
	195, 164, 201, 213, 23, 157, 251, 79,
	132, 51, 93, 156, 4, 156, 19, 156,
	66, 60, 96, 149, 52, 10, 94, 171,
	174, 197, 25, 250, 164, 145, 193, 217,
	140, 133, 91, 26, 93, 97, 173, 104,
	88, 13, 145, 241, 215, 129, 134, 7,
	166, 182, 174, 247, 216, 55, 201, 179,
	103, 237, 132, 49, 224, 59, 172, 66,
	255, 144, 177, 195, 144, 205, 85, 121,
	171, 144, 26, 180, 82, 238, 189, 146,
	229, 88, 93, 41, 219, 25, 240, 30,
	83, 197, 234, 152, 68, 252, 158, 41,
	191, 28, 150, 251, 240, 5, 143, 97,
	166, 44, 158, 227, 95, 239, 191, 124,
	173, 242, 230, 192, 21, 124, 104, 231,
	98, 98, 216, 56, 148, 8, 248, 19,
	167, 22, 88, 71, 20, 216, 10, 52,
	172, 137, 2, 91, 29, 1, 118, 208,
	11, 159, 78, 235, 224, 10, 22, 245,
	42, 153, 118, 128, 179, 69, 201, 151,
	203, 1, 188, 92, 8, 47, 118, 125,
	120, 249, 50, 224, 173, 105, 68, 63,
	101, 61, 228, 12, 123, 100, 216, 177,
	201, 148, 253, 144, 19, 2, 64, 121,
	207, 71, 133, 217, 206, 72, 105, 160,
	146, 51, 8, 108, 115, 5, 59, 117,
	69, 204, 208, 88, 222, 152, 251, 175,
	27, 107, 70, 217, 223, 120, 103, 205,
	212, 187, 137, 136, 176, 168, 47, 172,
	116, 160, 28, 184, 206, 198, 233, 9,
	167, 157, 171, 195, 240, 131, 48, 222,
	23, 249, 138, 137, 247, 115, 81, 65,
	177, 223, 137, 84, 176, 121, 80, 78,
	162, 96, 65, 177, 152, 98, 75, 5,
	87, 195, 29, 196, 218, 113, 7, 205,
	134, 108, 51, 176, 149, 74, 163, 141,
	119, 106, 240, 78, 187, 28, 132, 209,
	155, 244, 245, 192, 228, 28, 139, 119,
	241, 222, 80, 24, 254, 228, 14, 191,
	22, 188, 170, 7, 223, 10, 127, 7,
	0, 0, 255, 255, 110, 87, 69, 87,
}
