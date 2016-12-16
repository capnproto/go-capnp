package hashes

// AUTO GENERATED - DO NOT EDIT

import (
	context "golang.org/x/net/context"
	capnp "zombiezen.com/go/capnproto2"
	text "zombiezen.com/go/capnproto2/encoding/text"
	schemas "zombiezen.com/go/capnproto2/schemas"
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

// HashFactory_newSha1_Params_TypeID is the unique identifier for the type HashFactory_newSha1_Params.
const HashFactory_newSha1_Params_TypeID = 0x92b20ad1a58ca0ca

func NewHashFactory_newSha1_Params(s *capnp.Segment) (HashFactory_newSha1_Params, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0})
	return HashFactory_newSha1_Params{st}, err
}

func NewRootHashFactory_newSha1_Params(s *capnp.Segment) (HashFactory_newSha1_Params, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0})
	return HashFactory_newSha1_Params{st}, err
}

func ReadRootHashFactory_newSha1_Params(msg *capnp.Message) (HashFactory_newSha1_Params, error) {
	root, err := msg.RootPtr()
	return HashFactory_newSha1_Params{root.Struct()}, err
}

func (s HashFactory_newSha1_Params) String() string {
	str, _ := text.Marshal(0x92b20ad1a58ca0ca, s.Struct)
	return str
}

// HashFactory_newSha1_Params_List is a list of HashFactory_newSha1_Params.
type HashFactory_newSha1_Params_List struct{ capnp.List }

// NewHashFactory_newSha1_Params creates a new list of HashFactory_newSha1_Params.
func NewHashFactory_newSha1_Params_List(s *capnp.Segment, sz int32) (HashFactory_newSha1_Params_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0}, sz)
	return HashFactory_newSha1_Params_List{l}, err
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

// HashFactory_newSha1_Results_TypeID is the unique identifier for the type HashFactory_newSha1_Results.
const HashFactory_newSha1_Results_TypeID = 0xea3e50f7663f7bdf

func NewHashFactory_newSha1_Results(s *capnp.Segment) (HashFactory_newSha1_Results, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return HashFactory_newSha1_Results{st}, err
}

func NewRootHashFactory_newSha1_Results(s *capnp.Segment) (HashFactory_newSha1_Results, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return HashFactory_newSha1_Results{st}, err
}

func ReadRootHashFactory_newSha1_Results(msg *capnp.Message) (HashFactory_newSha1_Results, error) {
	root, err := msg.RootPtr()
	return HashFactory_newSha1_Results{root.Struct()}, err
}

func (s HashFactory_newSha1_Results) String() string {
	str, _ := text.Marshal(0xea3e50f7663f7bdf, s.Struct)
	return str
}

func (s HashFactory_newSha1_Results) Hash() Hash {
	p, _ := s.Struct.Ptr(0)
	return Hash{Client: p.Interface().Client()}
}

func (s HashFactory_newSha1_Results) HasHash() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s HashFactory_newSha1_Results) SetHash(v Hash) error {
	if v.Client == nil {
		return s.Struct.SetPtr(0, capnp.Ptr{})
	}
	seg := s.Segment()
	in := capnp.NewInterface(seg, seg.Message().AddCap(v.Client))
	return s.Struct.SetPtr(0, in.ToPtr())
}

// HashFactory_newSha1_Results_List is a list of HashFactory_newSha1_Results.
type HashFactory_newSha1_Results_List struct{ capnp.List }

// NewHashFactory_newSha1_Results creates a new list of HashFactory_newSha1_Results.
func NewHashFactory_newSha1_Results_List(s *capnp.Segment, sz int32) (HashFactory_newSha1_Results_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return HashFactory_newSha1_Results_List{l}, err
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

// Hash_write_Params_TypeID is the unique identifier for the type Hash_write_Params.
const Hash_write_Params_TypeID = 0xdffe94ae546cdee3

func NewHash_write_Params(s *capnp.Segment) (Hash_write_Params, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Hash_write_Params{st}, err
}

func NewRootHash_write_Params(s *capnp.Segment) (Hash_write_Params, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Hash_write_Params{st}, err
}

func ReadRootHash_write_Params(msg *capnp.Message) (Hash_write_Params, error) {
	root, err := msg.RootPtr()
	return Hash_write_Params{root.Struct()}, err
}

func (s Hash_write_Params) String() string {
	str, _ := text.Marshal(0xdffe94ae546cdee3, s.Struct)
	return str
}

func (s Hash_write_Params) Data() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	return []byte(p.Data()), err
}

func (s Hash_write_Params) HasData() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Hash_write_Params) SetData(v []byte) error {
	return s.Struct.SetData(0, v)
}

// Hash_write_Params_List is a list of Hash_write_Params.
type Hash_write_Params_List struct{ capnp.List }

// NewHash_write_Params creates a new list of Hash_write_Params.
func NewHash_write_Params_List(s *capnp.Segment, sz int32) (Hash_write_Params_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return Hash_write_Params_List{l}, err
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

// Hash_write_Results_TypeID is the unique identifier for the type Hash_write_Results.
const Hash_write_Results_TypeID = 0x80ac741ec7fb8f65

func NewHash_write_Results(s *capnp.Segment) (Hash_write_Results, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0})
	return Hash_write_Results{st}, err
}

func NewRootHash_write_Results(s *capnp.Segment) (Hash_write_Results, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0})
	return Hash_write_Results{st}, err
}

func ReadRootHash_write_Results(msg *capnp.Message) (Hash_write_Results, error) {
	root, err := msg.RootPtr()
	return Hash_write_Results{root.Struct()}, err
}

func (s Hash_write_Results) String() string {
	str, _ := text.Marshal(0x80ac741ec7fb8f65, s.Struct)
	return str
}

// Hash_write_Results_List is a list of Hash_write_Results.
type Hash_write_Results_List struct{ capnp.List }

// NewHash_write_Results creates a new list of Hash_write_Results.
func NewHash_write_Results_List(s *capnp.Segment, sz int32) (Hash_write_Results_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0}, sz)
	return Hash_write_Results_List{l}, err
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

// Hash_sum_Params_TypeID is the unique identifier for the type Hash_sum_Params.
const Hash_sum_Params_TypeID = 0xe74bb2d0190cf89c

func NewHash_sum_Params(s *capnp.Segment) (Hash_sum_Params, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0})
	return Hash_sum_Params{st}, err
}

func NewRootHash_sum_Params(s *capnp.Segment) (Hash_sum_Params, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0})
	return Hash_sum_Params{st}, err
}

func ReadRootHash_sum_Params(msg *capnp.Message) (Hash_sum_Params, error) {
	root, err := msg.RootPtr()
	return Hash_sum_Params{root.Struct()}, err
}

func (s Hash_sum_Params) String() string {
	str, _ := text.Marshal(0xe74bb2d0190cf89c, s.Struct)
	return str
}

// Hash_sum_Params_List is a list of Hash_sum_Params.
type Hash_sum_Params_List struct{ capnp.List }

// NewHash_sum_Params creates a new list of Hash_sum_Params.
func NewHash_sum_Params_List(s *capnp.Segment, sz int32) (Hash_sum_Params_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0}, sz)
	return Hash_sum_Params_List{l}, err
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

// Hash_sum_Results_TypeID is the unique identifier for the type Hash_sum_Results.
const Hash_sum_Results_TypeID = 0xd093963b95a4e107

func NewHash_sum_Results(s *capnp.Segment) (Hash_sum_Results, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Hash_sum_Results{st}, err
}

func NewRootHash_sum_Results(s *capnp.Segment) (Hash_sum_Results, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return Hash_sum_Results{st}, err
}

func ReadRootHash_sum_Results(msg *capnp.Message) (Hash_sum_Results, error) {
	root, err := msg.RootPtr()
	return Hash_sum_Results{root.Struct()}, err
}

func (s Hash_sum_Results) String() string {
	str, _ := text.Marshal(0xd093963b95a4e107, s.Struct)
	return str
}

func (s Hash_sum_Results) Hash() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	return []byte(p.Data()), err
}

func (s Hash_sum_Results) HasHash() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Hash_sum_Results) SetHash(v []byte) error {
	return s.Struct.SetData(0, v)
}

// Hash_sum_Results_List is a list of Hash_sum_Results.
type Hash_sum_Results_List struct{ capnp.List }

// NewHash_sum_Results creates a new list of Hash_sum_Results.
func NewHash_sum_Results_List(s *capnp.Segment, sz int32) (Hash_sum_Results_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return Hash_sum_Results_List{l}, err
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

const schema_db8274f9144abc7e = "x\xda\x84\x92?h\x13a\x18\xc6\x9f\xe7\xbb;\xaf\xa8" +
	"G\xfaq\x05\xe9bP\"\x82\xd8\xe2\xb5[\x05\x13\x1c" +
	"Tt\xb9\x8b\x0e\xe2\xf6QO#$\xb5\xe4.\x14\x11" +
	"\xff\xd0\xc5E\x10\xb4\xdaE\xd0A7\xed\xd0Q\xba\x8a" +
	"(\x08u\xb4\x12\x8b:\x08\xdd\xec\xa2\xa5\xd4\x93\xef\x92" +
	"k\x02!v\xfb\xe0y\xdf\xe7\x9e\xf7w\xcf\xe0fI" +
	"x\xd6a\x13\x08\x8eY\xbb\x92\xf0\xc1\xe6\xbb\xfd\xf1\xab" +
	"\xbb\x90\x83\x04L\x1bp7\xb8\x0e3\xf9\xf0\xfc\xfe\xcb" +
	"O\xbb\x17\x1fB\xeek\x0b\xe3M\x8e\x11f\xb2\xb4\xb2" +
	"5\x9f\xbb\xf8z\x01r\x8f\x91\xdc^:;\xb4\x11\xcf" +
	"~\x01\xe8\xbe\xe5\x1b\xf7#\xb5\xc5{\x9ev\x7f\xe9W" +
	"b\x7f{\xf1\xf8\xf8\x93G\xcb-\x7f+U?\xf3;" +
	"\xe86Y\x04\x93\x1f_\xab\x17\x16\xe6\xfe\xaev\xeb[" +
	"\\\x03]\x0a\xad?\xfd\xb3wxy\xf1\xdc\xcf\xae|" +
	"\x07\xc4\x0a\xccd\xf5f\xf1\xcao\xff\xc4Z+_\xba" +
	"8n\x89\x09\x82\xae\x93nzs\x97\xae6\xe7\x9f\xad" +
	"\xf7\xc4\x1c\x11\xb3\xae'\xb4\xd3\x88\xb8\xe7\xde\x126\x8e" +
	"&\x15\x15UF'\xd5\xb4\x98\x9a\x9e8\xa3\xdf3\xf5" +
	"kqX(\x87\xf9\xa8Q\x8d\xa3m\xddh\xeb\xa7\xd4" +
	"d|\xbd~ct*\x9c9_Q^\xc1\xcf\xab\xba" +
	"\xaau\xe6\x98\xcd\x15[\x83>\x19\x98\x86\x05lse" +
	"v\x80\x94'!\xa4e\xdfi{\x95\xe8\x93\xbd\x81\xa2" +
	"F\xadP\x0e\xa3\x86]\x8d\xa3\xc04L\xc0$ \x9d" +
	"#@0`0\x18\x12\xcc\xe9%:\x10t\xc0~'" +
	"\xf9*\xa7\x93\xf6\xb3\xb8\xacb\xd5\xdfB\x87\xf0U]" +
	"\x19\xb5\x9d\x91\x94\x8ba\xca\xee\xbfae\xe77\x81\x94" +
	"]\xdf\xcc\x08B\xa3\x1bH\xd1e]aVZ\xe9\x8d" +
	"A\xc8C6;=aV89|\x10B:v>" +
	"=\xbbD;j\xd4R\xb4\xff\x02\x00\x00\xff\xff<." +
	"\xe3\xa6"

func init() {
	schemas.Register(schema_db8274f9144abc7e,
		0x80ac741ec7fb8f65,
		0x92b20ad1a58ca0ca,
		0xaead580f97fddabc,
		0xd093963b95a4e107,
		0xdffe94ae546cdee3,
		0xe74bb2d0190cf89c,
		0xea3e50f7663f7bdf,
		0xf29f97dd675a9431)
}
