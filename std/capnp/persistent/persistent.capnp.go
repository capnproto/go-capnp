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

const schema_b8630836983feed7 = "x\xda\xbcU[l\x14e\x14\xfe\xcf\\\x98.\x82\xdd" +
	"\xbf\x03\x01\x92\xe2B\xad\x16\xcae\xb7\x95(l\xd4\xed" +
	"\x12\x154A:[/i\x13\xd4\xe9vX\xab\xbb3" +
	"\xcb\xce\x94\xeeb\x08\x97\x18/E\x1e\x1a^\x94\x07\x89" +
	"1\x8dBB|\xb1ZL*\x81\x88\x1a\xbc\xc4\x1a#" +
	"O&\x95\xc47\"\xc2\x83\xd7\xa4\xe3\xf9\xe7\xde\x0b6" +
	"\x06\xe3\xc3l\xf6?\xff\xf9\xcf\xe5;\xe7|'u\x84" +
	"\xeb\xe0\xda\xc4K\x02!\xcafq\x81\xfd\xdd\xfd\x0b\xbf" +
	")\xae\xb1_$\x94\xf2\xf6\xa5\x9f3\xaf\xdf]\x97\x1f" +
	"#\x84\xc4AV\xb8\xeb\xf2.N\"D\xee\xe6^\x96" +
	"U^\xc2\xaf\xc5>{\xb5\xe3\xb5s\xdb\x9f\xfd\x90\xd0" +
	"F\xb0\x87\x9f8\xf1Ur\xf5\x17\x9f\x13\x11\xa48\xdc" +
	"\xb5\x9f\xdf\x0a\xf2Q\x9e=y\x95\xcf\x10\x08\x0d\xc2-" +
	"\x84\x93O\xf3C\xf2(\x7f\x8f\xfc\x13\xbf\x93p\xf6\xfa" +
	"\xa7\x0f\x9e}\xf3\x8fO\xcf\x10%.\x82}\xe8\xd1\xd4" +
	"\xd8\x0b\xfb\x7f\x1f'\x04\xe4\xac\xf0\xad\xbcC@+]" +
	"\xdb\x05\x1eH\xe4\x12\x16\xa2\x99-\xc2a\xf9>\xa1E" +
	"\xee\x16\x98\x99\xc9u\xb5\xdb\xe3\x07N~L\x94\x18\x9a" +
	"y\xa5\xf8\xc3\x94\xd2\xd8:\xc1\xcc\x8c\x0aC\xf2\xb8c" +
	"f\xcc5s\xedhrY\xc33g\xce\x93\x89\x988" +
	"U?M\xf7\xa4P\x91O;\xba\xef\xba\xbaAf\xb3" +
	"qyC\xb8,\x8f\x08-\x98\xe49a\x9bLE\x09" +
	"\xbfe\xa15\xe01F*>\"/\x15\x07\xe5\x92\x98" +
	"\xc0\x18G\x1f\xb8\xf5N\xf8 \xf5\xe3\xec\x18K\xe2a" +
	"y\x8f\xc8\xfc\x16E\xc7\xef\xa9\xb77l?\xf2\xd6{" +
	"W\x08\xbd\x0d\x08\x119\x06\xeb.\xb1\x17\xaf\xe4~q" +
	"\x10\x15v<\x19\xdb\xb2\xf2\xcb\xf3\xbfD\x15\x1e_\xe0" +
	"(\xa8\x0b\x98B\xe1\xaf\x8b\xc3\x9d\xddM\xbf\x92\x09*" +
	"\xae\x8a\x94\x00\x15vH\x97\xe5nI\xc2/\xd1uP" +
	"r\xfc\x8d\x98\xa9\xe5\xdd\x1f\x19\xbf\xcdU\xce=R\x1a" +
	"\xe4C\x12+\xe7~\x89\x95\xb3\xacU\xcc~\xd3\xd28" +
	"\xdd\xda\x98W\xcbz9\x9d\xd3\xd4bi\x9b\x9a\xb0\xb4" +
	"A\xb5\xd6\x09\xa0\xd4\xf1\"!A\x90\xe0w\x0bmK" +
	"\x13\x92]\x0f\xd9{\x81\x1e\x92 \xcc3\xd4\x18`\x1a" +
	"e\xc8\x1e\x04zA\xca\xf4\x97\xcaF\xc5\xa2\x90P\x04" +
	"\x0e\xc2.\x05\xcc\xd9\x17\x06\xd12Q\x8bR\x07\xe0\xf8" +
	"g\xbfq\xfcm\x80\xc8;\xcc\xbeA\x04.*\xe8\x80" +
	"\x8cV\xbdi'\xe2l/\xfct/\x88J6\x0et" +
	"q/\xa5\xbdti\x85\xae\xa8\xd8\x0f\xeb\x96V\xd1\xd5" +
	"\"\x91r\xdan\xfb\xc1j\xf4\x14\xdc%v\x0e\xeaZ" +
	"%\xbc\xf5\xce~\x0d\x84\xa0\x06\x9d\x9e\x04\x05]\xea^" +
	"-\xa7\x99\x03E\xcb$\xac\x1a\x02\x8fs.\xb0t\x16" +
	"\xe7p\xe0\x17\xf1\xa0,\xc7\xa4Lk\xa0\xd2W\xcbi" +
	"\x04v;0E\x92\x84\x86H\x9d\xc1\xf7\xc1\xde\xc2\xb4" +
	"\xa1\xe8\x892\xc7sa\xd3\xb1\x1b? \xc2\xeb\x96\xed" +
	"u\x88E\xeaY\x8b\x04\xa6\xd9\x1d\xda\x8c\xd2@#\x07" +
	"\xd95\x10\x19\xe8f\x14l\x82\xc8\xf4\xa4P\xf0\x18\xc0" +
	"\"V\x8b\xf6\x84\x17\x99{\xec\x99f\xd9\x95\xb6\xc1\x10" +
	"\xd8\xfb\x8cRo\xbf\xb6O\x13\xf5\x8dy\xa3\x94,\x18" +
	"I\xe7]\xc5\xb0\x8c\xf6\xa4i\xf5\xb9\xc7d9\x00\x11" +
	"\xf3wE\xf9u\xdc:\x0fb]-ifY\xcd\x83" +
	"\x86\xa8\xa2i.P\x01O\x85\xb0\xd6\x08\x93\xa1\xb1\x9c" +
	"\xfd\xf5\xf0\xb5\xa9\xda;}\xd7\xf1\xd0j\xfb\x16\x08h" +
	"\xf5\xec\xff\x1c\xb9\xa7\xa2\xb9\xafG\xc1\xe6h\xee\x9b\"" +
	"\xb9\x1fw\xdd\xa7\xd3*\xe8\xbaa\xa9V\xbf\xc1\xeb\xa6" +
	"\x0fE\x93\x94\xafV}\x04r!\x02\xc2?#\x90\xaf" +
	"\x02\xbe\xb2\x0b\x86\x9b\x10\xa41\xda\xe7\xd5\x82F\x88\x9f" +
	"\xb3\x7fE\x12\x0e \xbex6+x\x0d !\x9a\x0e" +
	",!\xc9\xc4z\"\x0b$\xd6k\xb3\x86\xedT+*" +
	"\xe1K\xa6\xedw/\x91\xb0\x7f\xb1w\x19\x93\xf8OC" +
	"\x9e\xa0\xad\xc8\x13\x8b \xdb\x08t\x83To\xe2\x9bY" +
	"\xb3:\xe7\x003\xa1;\x8du@\xc5\x1c\x8d\xb5\xdb]" +
	"\xe1\x1cx\xd3\x15\xa4\x88\x1d\xdf\x08\x91\x82\xd0\xb6\xada" +
	"1\xe8\x86\xb4\xfd\xd4\xb1\x13\xca\xf8\xf7C\x17\x08]\xdb" +
	"d\x7fv\xf5b\xf3\x8a\xf7\xad\x11B\xefh\xb2\x1b&" +
	"sWj/\xed\xc5)Y\xddn\x1f[\x95\x9c<\xae" +
	"\xc5\xff$teO\xb8\x80\xe8\xca\xd6\x03\x1e\xbc\x1e\xd3" +
	"I}F^\xb2\xd4B\x82\xd5\xb3`\xe7\x07L\xcb(" +
	"Y5\xc2\x97\xbd~\x11 \xba\xeb03$\x95\x19\x93" +
	"\x10\xa9\x9dk\x94\xcc,\x910\x93\xb8\x1d\xde\xde\xe8\xf2" +
	"`3\x16BRK&\xd2\x9b\xcf\x1ak\x9b\x10\x87f" +
	"d\x8d*\x07\x14`\x89\x03,#j\xa5\x8c\xc2O8" +
	"\x90\xd0\x1a\xd0(\x81t\xc0\x7f\xc4\xcf@\x09d\xca\xd8" +
	"\x1b%\x13\xe2a\x0b\xfd\x1b\x0f\xf3\x913\xc4#dw" +
	"\x03l\\$\xffol\xe6\x8d\xfc\xe6\xc1\x99\x0f\xfe\xf8" +
	"\x9c\x0b\xbf\x1c\x8e\xb6\xb3`\xf0Q\xd0\xd5NW.\xf7" +
	"\xba\xf2T\xb0\x06p\x8bd\x1d\x9a\xaaG\x9e\xd2\x89\x10" +
	"X\xe5o\xb4\xc22\x9dNf36\xd8V\x84\x16\xe3" +
	"W\x96pp\xc0\xc4\x12=dT\\\xa0f\xac\xaf\xbf" +
	"\x03\x00\x00\xff\xff\xb4QC^"
