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

const schema_db8274f9144abc7e = "x\xda\x84T_HdU\x18?\xdf\xb9s\xbbZ\xda" +
	"\xcc\xf1\x1a&XS\xa6I\x139\xce\xf8\x10L\xd1\x98" +
	"D\x85!\xcc\xd5\x1e\xc2\x87\xe8:\xde\xc6\xa1\xb9\xf7\x0e" +
	"3w\x921\xc8\xf0\xa1B\x8a -!\x0a{\xc8\x07" +
	"\xa1|\x90\x0a\x0c\xdb\x87\x05\xf1\x0f\x08\xfa\xb2\xb0\xee\xba" +
	"\xea\xee\xc3\x82,,\x0a\xcb\xae\xbb\xacw\xbf33\xd7" +
	"\xb90\xe8>\xccp\xce\xfd\xfe\xfc~\xe7\xf7\xfd\xce\xe9" +
	"\xe8\x85.\x1a\x12\xdb<\x84(\x1d\xe2\x13\xb6\xf6\xc3\x83" +
	"\x95\xe7\xad?\xbf\"\xcc\x07\x84x$B\xe4c8\"" +
	"\x1e{\xfd\xf7\xeff\xb7\x9e\\\xf8\x91\xb0\x86R\xa0s" +
	"\x07\xc2\x80\x91\xa5\xed\x87\xd3\xde\x8f\xfe\x9a'\xec)\xc1" +
	"\xfer\xa9\xa7\xfe\xd8\x1a\xbfB\x08\xc8\xcb\xf0\x9f\xbc\x01" +
	"\xbc\xc5\x1a\xbc'\x1f\xf2\x95\xbd\xf7j\xfe%\xdf\xd8\xdc" +
	"\x05\xa2T\x8b`\x7f\x9b\xda9Q\x9a\x02[<\xfb2" +
	"L\xc8\xfb<\xa7\xff*\x08@\xc0\x96\xf6\xff\xf8\xe9\x8d" +
	"\x9f'7\x8bT\xc4B\xa3\x0d\xb8\x8e\xa9[\x10%\xae" +
	"b\x10\x08\xc5\xf6=\xf2\x1d\x18\x91{\xa9\x9f\xd02\x0d" +
	"\x101\xd6K\x07d\x85\xb6\xc9:}\x1dc7\xae\xa5" +
	">\x9c\x9f:\xd9u\xf7\x9d\xa5\x07\xd8w\x8e\xf2\xbe\xff" +
	"\xbc\xf3\xf4\xcb\xf0o\xc7~%\xc55:.oPN" +
	"q\x95\x16(\xfez\xaf\xa6qs\xe1\x83\x9b.\xb5\x16" +
	"\xe96j\xb2\xfbE\xf4\xd3\xbb\xb1\xb7\x0e\x8aj\x150" +
	":\x7f\xa3\x11\xacA(\x0e\x12\x9a\x1aH\xecL\xcf\x1c" +
	"U\x8a\x86\x18k\x1c\x03W\xdf\xc8\xcf\x08\\\xb4a5" +
	";\xdc\x1eW\xd3\xd4HG\xde\xe7\xeb\x91L\xd2\xd2Z" +
	"\xfa4\x7f6\x97\xb2\xb2\xa7q\xa1\x14\x7fW\x8d[f" +
	"&\xdfnh#\xfd\xc3j\xa8%\xe6W3\xaa^\xce" +
	"\x03'/ZL\x8c\x01(\x1eAD$g\xca\xe0\x1c" +
	"\x80\xb1nB\x99(\x8d\x95zu\x01&\xdb\x09\x93\xb7" +
	"1\xd2\x10I\xab\xf1\xcf\xd4\x84F\x08~\x86\x1aT\xb7" +
	"\x82k6\xa7#\xd3lNB\xa6\x88\x82F\xf3\xa0$" +
	"\xac6\x80\x8e\xab\x12@\xa9\xa7\xe0\xe5EPK(\xfe" +
	"\xca\xcd\x09&4!\xd8\xa9eX\xa8\xbb<\x1c\xf6Z" +
	"\xc4\xfexrF\xf9\xff\xd2\xc42a\xaf4\xdb+\xb7" +
	"\xd7[\x1a\xff\xb6f\x09km\xb6\xeb\xf6\xfan\xe5\xbf" +
	"\xfe|\x95\xb0\x17\xc3\xf6\xe4\x0b\xc1\xbd_4\xdf}\xc2" +
	"\x9e\x1b\xb0\x0f\xbf\x0f6\xd4}\xb2x\x117\x81\xb1\x12" +
	"\xfbhRO\x9b\x19K\x1a2\xe3\x92\xa5&\xfc\x86\x89" +
	"\xffv<\x97\xb5L\xdd\xca\x13!\xady\x0dU\xd7\x14" +
	"\x0fP\x97\x81=\x14\xde\xf6\xf1C\x13\x06a\x7f\x89\xf2" +
	"\xe9\xf9\x09\x0eC\xa9B\xfa\xe5\xebQ=\xe8\x1a{u" +
	"\xc0vFE$\x9c\x81\x97\xef\xb0\xc0\x8d\xe0C\x84&" +
	"\xb7!\x9f\xc5\x0fo\x96 #Q\x0e\xa5e\xd14\xb8" +
	"\x0f\xc1\x04\xd8\xa3\xa6>\x98\xd4F5\xd1h\x8f\x9bz" +
	"0a\x06\x0b\xac2\xa6e\x86\x83I\xc3\xd22\x86\x9a" +
	"\x0a\x0ei\xba\x19,\xd5\xc2Y\xde\x8a\xa9^n\x99\xb3" +
	"\x066\xa4Zj\xc5\xc0 R\xd4\xf1|/\xc4\xd0\x8b" +
	"\x82\xfex\xd3\xf6E\xb5\x82\xbb\xcf\xf5\x0c++J\x00" +
	"\x98\xeb<\x8e\xc79\x19\xac\xe1\xe6v.>8\x8f\x1c" +
	"\x0b\x85\xd1\xdc\xad\x12\x94o28\xaf\x0ekl\xc6X" +
	"\xad\xe4/\xe8\xd1\x05\x12r/\x98\xffQ\x00\x00\x00\xff" +
	"\xff\x1bz\x9e\xa5"
