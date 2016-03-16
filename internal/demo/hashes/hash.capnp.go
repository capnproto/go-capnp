package hashes

// AUTO GENERATED - DO NOT EDIT

import (
	context "golang.org/x/net/context"
	capnp "zombiezen.com/go/capnproto2"
	server "zombiezen.com/go/capnproto2/server"
)

type HashFactory struct{ Client capnp.Client }

func (c HashFactory) NewSha1(ctx context.Context, params func(HashFactory_newSha1_Params) error, opts ...capnp.CallOption) HashFactory_newSha1_Results_Promise {
	if c.Client == nil {
		return HashFactory_newSha1_Results_Promise{Pipeline: capnp.NewPipeline(capnp.ErrorAnswer(capnp.ErrNullClient))}
	}
	call := &capnp.Call{
		Ctx: ctx,
		Method: capnp.Method{

			InterfaceID:   0xaead580f97fddabc,
			MethodID:      0,
			InterfaceName: "hash.capnp:HashFactory",
			MethodName:    "newSha1",
		},
		Options: capnp.NewCallOptions(opts),
	}
	if params != nil {
		call.ParamsSize = capnp.ObjectSize{DataSize: 0, PointerCount: 0}
		call.ParamsFunc = func(s capnp.Struct) error { return params(HashFactory_newSha1_Params{Struct: s}) }
	}
	return HashFactory_newSha1_Results_Promise{Pipeline: capnp.NewPipeline(c.Client.Call(call))}
}

type HashFactory_Server interface {
	NewSha1(HashFactory_newSha1) error
}

func HashFactory_ServerToClient(s HashFactory_Server) HashFactory {
	c, _ := s.(server.Closer)
	return HashFactory{Client: server.New(HashFactory_Methods(nil, s), c)}
}

func HashFactory_Methods(methods []server.Method, s HashFactory_Server) []server.Method {
	if cap(methods) == 0 {
		methods = make([]server.Method, 0, 1)
	}

	methods = append(methods, server.Method{
		Method: capnp.Method{

			InterfaceID:   0xaead580f97fddabc,
			MethodID:      0,
			InterfaceName: "hash.capnp:HashFactory",
			MethodName:    "newSha1",
		},
		Impl: func(c context.Context, opts capnp.CallOptions, p, r capnp.Struct) error {
			call := HashFactory_newSha1{c, opts, HashFactory_newSha1_Params{Struct: p}, HashFactory_newSha1_Results{Struct: r}}
			return s.NewSha1(call)
		},
		ResultsSize: capnp.ObjectSize{DataSize: 0, PointerCount: 1},
	})

	return methods
}

// HashFactory_newSha1 holds the arguments for a server call to HashFactory.newSha1.
type HashFactory_newSha1 struct {
	Ctx     context.Context
	Options capnp.CallOptions
	Params  HashFactory_newSha1_Params
	Results HashFactory_newSha1_Results
}

type HashFactory_newSha1_Params struct{ capnp.Struct }

func NewHashFactory_newSha1_Params(s *capnp.Segment) (HashFactory_newSha1_Params, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0})
	if err != nil {
		return HashFactory_newSha1_Params{}, err
	}
	return HashFactory_newSha1_Params{st}, nil
}

func NewRootHashFactory_newSha1_Params(s *capnp.Segment) (HashFactory_newSha1_Params, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0})
	if err != nil {
		return HashFactory_newSha1_Params{}, err
	}
	return HashFactory_newSha1_Params{st}, nil
}

func ReadRootHashFactory_newSha1_Params(msg *capnp.Message) (HashFactory_newSha1_Params, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return HashFactory_newSha1_Params{}, err
	}
	return HashFactory_newSha1_Params{root.Struct()}, nil
}

// HashFactory_newSha1_Params_List is a list of HashFactory_newSha1_Params.
type HashFactory_newSha1_Params_List struct{ capnp.List }

// NewHashFactory_newSha1_Params creates a new list of HashFactory_newSha1_Params.
func NewHashFactory_newSha1_Params_List(s *capnp.Segment, sz int32) (HashFactory_newSha1_Params_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0}, sz)
	if err != nil {
		return HashFactory_newSha1_Params_List{}, err
	}
	return HashFactory_newSha1_Params_List{l}, nil
}

func (s HashFactory_newSha1_Params_List) At(i int) HashFactory_newSha1_Params {
	return HashFactory_newSha1_Params{s.List.Struct(i)}
}
func (s HashFactory_newSha1_Params_List) Set(i int, v HashFactory_newSha1_Params) error {
	return s.List.SetStruct(i, v.Struct)
}

// HashFactory_newSha1_Params_Promise is a wrapper for a HashFactory_newSha1_Params promised by a client call.
type HashFactory_newSha1_Params_Promise struct{ *capnp.Pipeline }

func (p HashFactory_newSha1_Params_Promise) Struct() (HashFactory_newSha1_Params, error) {
	s, err := p.Pipeline.Struct()
	return HashFactory_newSha1_Params{s}, err
}

type HashFactory_newSha1_Results struct{ capnp.Struct }

func NewHashFactory_newSha1_Results(s *capnp.Segment) (HashFactory_newSha1_Results, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return HashFactory_newSha1_Results{}, err
	}
	return HashFactory_newSha1_Results{st}, nil
}

func NewRootHashFactory_newSha1_Results(s *capnp.Segment) (HashFactory_newSha1_Results, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return HashFactory_newSha1_Results{}, err
	}
	return HashFactory_newSha1_Results{st}, nil
}

func ReadRootHashFactory_newSha1_Results(msg *capnp.Message) (HashFactory_newSha1_Results, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return HashFactory_newSha1_Results{}, err
	}
	return HashFactory_newSha1_Results{root.Struct()}, nil
}

func (s HashFactory_newSha1_Results) Hash() Hash {
	p, err := s.Struct.Ptr(0)
	if err != nil {

		return Hash{}
	}
	return Hash{Client: p.Interface().Client()}
}

func (s HashFactory_newSha1_Results) SetHash(v Hash) error {

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

// HashFactory_newSha1_Results_List is a list of HashFactory_newSha1_Results.
type HashFactory_newSha1_Results_List struct{ capnp.List }

// NewHashFactory_newSha1_Results creates a new list of HashFactory_newSha1_Results.
func NewHashFactory_newSha1_Results_List(s *capnp.Segment, sz int32) (HashFactory_newSha1_Results_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return HashFactory_newSha1_Results_List{}, err
	}
	return HashFactory_newSha1_Results_List{l}, nil
}

func (s HashFactory_newSha1_Results_List) At(i int) HashFactory_newSha1_Results {
	return HashFactory_newSha1_Results{s.List.Struct(i)}
}
func (s HashFactory_newSha1_Results_List) Set(i int, v HashFactory_newSha1_Results) error {
	return s.List.SetStruct(i, v.Struct)
}

// HashFactory_newSha1_Results_Promise is a wrapper for a HashFactory_newSha1_Results promised by a client call.
type HashFactory_newSha1_Results_Promise struct{ *capnp.Pipeline }

func (p HashFactory_newSha1_Results_Promise) Struct() (HashFactory_newSha1_Results, error) {
	s, err := p.Pipeline.Struct()
	return HashFactory_newSha1_Results{s}, err
}

func (p HashFactory_newSha1_Results_Promise) Hash() Hash {
	return Hash{Client: p.Pipeline.GetPipeline(0).Client()}
}

type Hash struct{ Client capnp.Client }

func (c Hash) Write(ctx context.Context, params func(Hash_write_Params) error, opts ...capnp.CallOption) Hash_write_Results_Promise {
	if c.Client == nil {
		return Hash_write_Results_Promise{Pipeline: capnp.NewPipeline(capnp.ErrorAnswer(capnp.ErrNullClient))}
	}
	call := &capnp.Call{
		Ctx: ctx,
		Method: capnp.Method{

			InterfaceID:   0xf29f97dd675a9431,
			MethodID:      0,
			InterfaceName: "hash.capnp:Hash",
			MethodName:    "write",
		},
		Options: capnp.NewCallOptions(opts),
	}
	if params != nil {
		call.ParamsSize = capnp.ObjectSize{DataSize: 0, PointerCount: 1}
		call.ParamsFunc = func(s capnp.Struct) error { return params(Hash_write_Params{Struct: s}) }
	}
	return Hash_write_Results_Promise{Pipeline: capnp.NewPipeline(c.Client.Call(call))}
}

func (c Hash) Sum(ctx context.Context, params func(Hash_sum_Params) error, opts ...capnp.CallOption) Hash_sum_Results_Promise {
	if c.Client == nil {
		return Hash_sum_Results_Promise{Pipeline: capnp.NewPipeline(capnp.ErrorAnswer(capnp.ErrNullClient))}
	}
	call := &capnp.Call{
		Ctx: ctx,
		Method: capnp.Method{

			InterfaceID:   0xf29f97dd675a9431,
			MethodID:      1,
			InterfaceName: "hash.capnp:Hash",
			MethodName:    "sum",
		},
		Options: capnp.NewCallOptions(opts),
	}
	if params != nil {
		call.ParamsSize = capnp.ObjectSize{DataSize: 0, PointerCount: 0}
		call.ParamsFunc = func(s capnp.Struct) error { return params(Hash_sum_Params{Struct: s}) }
	}
	return Hash_sum_Results_Promise{Pipeline: capnp.NewPipeline(c.Client.Call(call))}
}

type Hash_Server interface {
	Write(Hash_write) error

	Sum(Hash_sum) error
}

func Hash_ServerToClient(s Hash_Server) Hash {
	c, _ := s.(server.Closer)
	return Hash{Client: server.New(Hash_Methods(nil, s), c)}
}

func Hash_Methods(methods []server.Method, s Hash_Server) []server.Method {
	if cap(methods) == 0 {
		methods = make([]server.Method, 0, 2)
	}

	methods = append(methods, server.Method{
		Method: capnp.Method{

			InterfaceID:   0xf29f97dd675a9431,
			MethodID:      0,
			InterfaceName: "hash.capnp:Hash",
			MethodName:    "write",
		},
		Impl: func(c context.Context, opts capnp.CallOptions, p, r capnp.Struct) error {
			call := Hash_write{c, opts, Hash_write_Params{Struct: p}, Hash_write_Results{Struct: r}}
			return s.Write(call)
		},
		ResultsSize: capnp.ObjectSize{DataSize: 0, PointerCount: 0},
	})

	methods = append(methods, server.Method{
		Method: capnp.Method{

			InterfaceID:   0xf29f97dd675a9431,
			MethodID:      1,
			InterfaceName: "hash.capnp:Hash",
			MethodName:    "sum",
		},
		Impl: func(c context.Context, opts capnp.CallOptions, p, r capnp.Struct) error {
			call := Hash_sum{c, opts, Hash_sum_Params{Struct: p}, Hash_sum_Results{Struct: r}}
			return s.Sum(call)
		},
		ResultsSize: capnp.ObjectSize{DataSize: 0, PointerCount: 1},
	})

	return methods
}

// Hash_write holds the arguments for a server call to Hash.write.
type Hash_write struct {
	Ctx     context.Context
	Options capnp.CallOptions
	Params  Hash_write_Params
	Results Hash_write_Results
}

// Hash_sum holds the arguments for a server call to Hash.sum.
type Hash_sum struct {
	Ctx     context.Context
	Options capnp.CallOptions
	Params  Hash_sum_Params
	Results Hash_sum_Results
}

type Hash_write_Params struct{ capnp.Struct }

func NewHash_write_Params(s *capnp.Segment) (Hash_write_Params, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Hash_write_Params{}, err
	}
	return Hash_write_Params{st}, nil
}

func NewRootHash_write_Params(s *capnp.Segment) (Hash_write_Params, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Hash_write_Params{}, err
	}
	return Hash_write_Params{st}, nil
}

func ReadRootHash_write_Params(msg *capnp.Message) (Hash_write_Params, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Hash_write_Params{}, err
	}
	return Hash_write_Params{root.Struct()}, nil
}

func (s Hash_write_Params) Data() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}

	return []byte(p.Data()), nil

}

func (s Hash_write_Params) SetData(v []byte) error {

	d, err := capnp.NewData(s.Struct.Segment(), []byte(v))
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, d.List.ToPtr())
}

// Hash_write_Params_List is a list of Hash_write_Params.
type Hash_write_Params_List struct{ capnp.List }

// NewHash_write_Params creates a new list of Hash_write_Params.
func NewHash_write_Params_List(s *capnp.Segment, sz int32) (Hash_write_Params_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return Hash_write_Params_List{}, err
	}
	return Hash_write_Params_List{l}, nil
}

func (s Hash_write_Params_List) At(i int) Hash_write_Params {
	return Hash_write_Params{s.List.Struct(i)}
}
func (s Hash_write_Params_List) Set(i int, v Hash_write_Params) error {
	return s.List.SetStruct(i, v.Struct)
}

// Hash_write_Params_Promise is a wrapper for a Hash_write_Params promised by a client call.
type Hash_write_Params_Promise struct{ *capnp.Pipeline }

func (p Hash_write_Params_Promise) Struct() (Hash_write_Params, error) {
	s, err := p.Pipeline.Struct()
	return Hash_write_Params{s}, err
}

type Hash_write_Results struct{ capnp.Struct }

func NewHash_write_Results(s *capnp.Segment) (Hash_write_Results, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0})
	if err != nil {
		return Hash_write_Results{}, err
	}
	return Hash_write_Results{st}, nil
}

func NewRootHash_write_Results(s *capnp.Segment) (Hash_write_Results, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0})
	if err != nil {
		return Hash_write_Results{}, err
	}
	return Hash_write_Results{st}, nil
}

func ReadRootHash_write_Results(msg *capnp.Message) (Hash_write_Results, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Hash_write_Results{}, err
	}
	return Hash_write_Results{root.Struct()}, nil
}

// Hash_write_Results_List is a list of Hash_write_Results.
type Hash_write_Results_List struct{ capnp.List }

// NewHash_write_Results creates a new list of Hash_write_Results.
func NewHash_write_Results_List(s *capnp.Segment, sz int32) (Hash_write_Results_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0}, sz)
	if err != nil {
		return Hash_write_Results_List{}, err
	}
	return Hash_write_Results_List{l}, nil
}

func (s Hash_write_Results_List) At(i int) Hash_write_Results {
	return Hash_write_Results{s.List.Struct(i)}
}
func (s Hash_write_Results_List) Set(i int, v Hash_write_Results) error {
	return s.List.SetStruct(i, v.Struct)
}

// Hash_write_Results_Promise is a wrapper for a Hash_write_Results promised by a client call.
type Hash_write_Results_Promise struct{ *capnp.Pipeline }

func (p Hash_write_Results_Promise) Struct() (Hash_write_Results, error) {
	s, err := p.Pipeline.Struct()
	return Hash_write_Results{s}, err
}

type Hash_sum_Params struct{ capnp.Struct }

func NewHash_sum_Params(s *capnp.Segment) (Hash_sum_Params, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0})
	if err != nil {
		return Hash_sum_Params{}, err
	}
	return Hash_sum_Params{st}, nil
}

func NewRootHash_sum_Params(s *capnp.Segment) (Hash_sum_Params, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0})
	if err != nil {
		return Hash_sum_Params{}, err
	}
	return Hash_sum_Params{st}, nil
}

func ReadRootHash_sum_Params(msg *capnp.Message) (Hash_sum_Params, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Hash_sum_Params{}, err
	}
	return Hash_sum_Params{root.Struct()}, nil
}

// Hash_sum_Params_List is a list of Hash_sum_Params.
type Hash_sum_Params_List struct{ capnp.List }

// NewHash_sum_Params creates a new list of Hash_sum_Params.
func NewHash_sum_Params_List(s *capnp.Segment, sz int32) (Hash_sum_Params_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0}, sz)
	if err != nil {
		return Hash_sum_Params_List{}, err
	}
	return Hash_sum_Params_List{l}, nil
}

func (s Hash_sum_Params_List) At(i int) Hash_sum_Params { return Hash_sum_Params{s.List.Struct(i)} }
func (s Hash_sum_Params_List) Set(i int, v Hash_sum_Params) error {
	return s.List.SetStruct(i, v.Struct)
}

// Hash_sum_Params_Promise is a wrapper for a Hash_sum_Params promised by a client call.
type Hash_sum_Params_Promise struct{ *capnp.Pipeline }

func (p Hash_sum_Params_Promise) Struct() (Hash_sum_Params, error) {
	s, err := p.Pipeline.Struct()
	return Hash_sum_Params{s}, err
}

type Hash_sum_Results struct{ capnp.Struct }

func NewHash_sum_Results(s *capnp.Segment) (Hash_sum_Results, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Hash_sum_Results{}, err
	}
	return Hash_sum_Results{st}, nil
}

func NewRootHash_sum_Results(s *capnp.Segment) (Hash_sum_Results, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Hash_sum_Results{}, err
	}
	return Hash_sum_Results{st}, nil
}

func ReadRootHash_sum_Results(msg *capnp.Message) (Hash_sum_Results, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Hash_sum_Results{}, err
	}
	return Hash_sum_Results{root.Struct()}, nil
}

func (s Hash_sum_Results) Hash() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}

	return []byte(p.Data()), nil

}

func (s Hash_sum_Results) SetHash(v []byte) error {

	d, err := capnp.NewData(s.Struct.Segment(), []byte(v))
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, d.List.ToPtr())
}

// Hash_sum_Results_List is a list of Hash_sum_Results.
type Hash_sum_Results_List struct{ capnp.List }

// NewHash_sum_Results creates a new list of Hash_sum_Results.
func NewHash_sum_Results_List(s *capnp.Segment, sz int32) (Hash_sum_Results_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return Hash_sum_Results_List{}, err
	}
	return Hash_sum_Results_List{l}, nil
}

func (s Hash_sum_Results_List) At(i int) Hash_sum_Results { return Hash_sum_Results{s.List.Struct(i)} }
func (s Hash_sum_Results_List) Set(i int, v Hash_sum_Results) error {
	return s.List.SetStruct(i, v.Struct)
}

// Hash_sum_Results_Promise is a wrapper for a Hash_sum_Results promised by a client call.
type Hash_sum_Results_Promise struct{ *capnp.Pipeline }

func (p Hash_sum_Results_Promise) Struct() (Hash_sum_Results, error) {
	s, err := p.Pipeline.Struct()
	return Hash_sum_Results{s}, err
}
