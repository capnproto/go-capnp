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

func (s HashFactory_newSha1_Results) HasHash() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
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

func (s Hash_write_Params) HasData() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
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

func (s Hash_sum_Results) HasHash() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
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

var schema_db8274f9144abc7e = []byte{
	120, 218, 132, 83, 95, 72, 44, 101,
	20, 255, 206, 55, 51, 205, 181, 220,
	118, 191, 59, 23, 110, 11, 91, 155,
	166, 93, 238, 198, 221, 117, 247, 62,
	4, 91, 180, 118, 41, 11, 35, 216,
	89, 11, 194, 135, 104, 92, 167, 117,
	105, 103, 102, 217, 29, 147, 53, 200,
	16, 236, 143, 36, 81, 90, 66, 20,
	38, 228, 67, 80, 62, 72, 6, 133,
	5, 5, 98, 248, 164, 47, 61, 24,
	166, 214, 67, 32, 65, 228, 75, 89,
	228, 215, 153, 217, 29, 119, 204, 43,
	62, 40, 243, 237, 239, 156, 223, 57,
	231, 247, 59, 167, 235, 73, 232, 166,
	73, 233, 138, 72, 136, 218, 37, 221,
	194, 119, 94, 202, 60, 255, 103, 246,
	161, 125, 194, 46, 3, 33, 18, 200,
	132, 92, 63, 132, 52, 16, 80, 128,
	102, 8, 240, 95, 126, 42, 61, 181,
	56, 115, 180, 67, 88, 200, 11, 80,
	218, 232, 62, 226, 157, 46, 254, 193,
	95, 173, 225, 141, 165, 39, 126, 173,
	227, 162, 3, 63, 74, 183, 136, 200,
	95, 47, 109, 31, 169, 145, 216, 38,
	1, 129, 80, 229, 42, 237, 85, 174,
	209, 17, 101, 130, 70, 9, 229, 235,
	31, 189, 185, 176, 121, 235, 210, 59,
	245, 170, 78, 210, 245, 9, 154, 2,
	204, 210, 223, 250, 103, 237, 46, 251,
	211, 87, 124, 116, 6, 61, 64, 32,
	57, 211, 95, 216, 158, 157, 59, 32,
	236, 54, 129, 191, 188, 210, 123, 233,
	208, 30, 255, 145, 96, 27, 42, 29,
	87, 158, 166, 78, 160, 74, 95, 83,
	230, 157, 47, 190, 178, 245, 239, 108,
	240, 153, 207, 22, 79, 5, 191, 65,
	191, 82, 222, 118, 131, 167, 232, 99,
	202, 178, 27, 44, 239, 125, 252, 238,
	3, 239, 77, 111, 248, 39, 252, 144,
	254, 140, 209, 243, 238, 132, 199, 249,
	32, 225, 32, 223, 210, 126, 101, 149,
	94, 81, 246, 232, 253, 56, 200, 238,
	125, 181, 123, 66, 99, 159, 124, 67,
	212, 22, 9, 154, 35, 99, 110, 88,
	152, 84, 218, 4, 228, 234, 139, 8,
	2, 202, 201, 151, 31, 185, 253, 94,
	248, 162, 107, 239, 116, 104, 64, 24,
	87, 152, 27, 218, 90, 15, 29, 210,
	170, 67, 241, 188, 86, 22, 204, 114,
	250, 113, 252, 238, 209, 242, 182, 85,
	169, 197, 77, 125, 164, 111, 72, 75,
	118, 228, 50, 122, 117, 184, 100, 87,
	85, 81, 64, 31, 69, 236, 154, 5,
	98, 104, 232, 5, 1, 212, 75, 20,
	130, 78, 62, 176, 166, 100, 4, 128,
	249, 104, 105, 131, 54, 62, 82, 41,
	218, 122, 71, 86, 11, 86, 52, 227,
	76, 178, 65, 205, 214, 32, 64, 40,
	254, 221, 132, 162, 58, 108, 32, 65,
	69, 19, 140, 42, 47, 88, 14, 100,
	150, 9, 102, 71, 0, 154, 226, 176,
	228, 141, 230, 248, 236, 90, 154, 63,
	59, 61, 167, 126, 253, 195, 228, 42,
	97, 87, 219, 249, 218, 239, 235, 29,
	225, 207, 237, 5, 194, 58, 219, 249,
	197, 221, 220, 111, 181, 87, 95, 252,
	158, 176, 182, 20, 159, 190, 59, 177,
	251, 190, 30, 250, 155, 176, 59, 251,
	249, 31, 83, 137, 203, 23, 159, 251,
	242, 59, 124, 196, 198, 202, 90, 254,
	5, 173, 160, 103, 138, 70, 217, 170,
	216, 242, 160, 149, 151, 109, 173, 16,
	53, 45, 252, 207, 243, 195, 85, 219,
	50, 236, 26, 17, 202, 122, 208, 212,
	12, 93, 21, 193, 111, 149, 72, 225,
	225, 16, 64, 43, 78, 10, 169, 104,
	163, 229, 115, 69, 207, 70, 53, 71,
	166, 179, 84, 204, 233, 81, 215, 147,
	99, 28, 26, 56, 201, 2, 160, 152,
	18, 150, 240, 110, 9, 188, 45, 103,
	201, 20, 161, 172, 83, 134, 230, 29,
	129, 183, 142, 44, 220, 142, 88, 64,
	142, 186, 252, 221, 32, 163, 212, 221,
	128, 100, 167, 42, 244, 100, 234, 173,
	58, 133, 68, 183, 144, 119, 95, 224,
	157, 55, 99, 55, 144, 76, 146, 199,
	26, 211, 156, 36, 58, 225, 102, 14,
	119, 75, 62, 111, 185, 254, 191, 15,
	4, 85, 83, 47, 32, 101, 243, 242,
	90, 6, 124, 55, 219, 18, 227, 158,
	166, 68, 198, 86, 131, 206, 11, 19,
	252, 182, 132, 208, 150, 136, 255, 78,
	238, 192, 31, 30, 108, 248, 148, 206,
	56, 165, 244, 42, 238, 50, 190, 147,
	48, 9, 124, 212, 50, 6, 138, 250,
	168, 46, 153, 241, 188, 101, 36, 10,
	86, 194, 181, 178, 98, 217, 86, 42,
	81, 52, 109, 189, 98, 106, 165, 196,
	160, 110, 88, 137, 70, 46, 28, 175,
	40, 164, 27, 27, 68, 28, 123, 144,
	146, 250, 160, 250, 82, 121, 192, 127,
	1, 0, 0, 255, 255, 76, 34, 149,
	126,
}
