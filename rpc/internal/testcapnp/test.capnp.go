// Code generated by capnpc-go. DO NOT EDIT.

package testcapnp

import (
	capnp "capnproto.org/go/capnp/v3"
	text "capnproto.org/go/capnp/v3/encoding/text"
	schemas "capnproto.org/go/capnp/v3/schemas"
	server "capnproto.org/go/capnp/v3/server"
	stream "capnproto.org/go/capnp/v3/std/capnp/stream"
	context "context"
)

type PingPong struct{ Client *capnp.Client }

// PingPong_TypeID is the unique identifier for the type PingPong.
const PingPong_TypeID = 0xf004c474c2f8ee7a

func (c PingPong) EchoNum(ctx context.Context, params func(PingPong_echoNum_Params) error) (PingPong_echoNum_Results_Future, capnp.ReleaseFunc) {
	s := capnp.Send{
		Method: capnp.Method{
			InterfaceID:   0xf004c474c2f8ee7a,
			MethodID:      0,
			InterfaceName: "test.capnp:PingPong",
			MethodName:    "echoNum",
		},
	}
	if params != nil {
		s.ArgsSize = capnp.ObjectSize{DataSize: 8, PointerCount: 0}
		s.PlaceArgs = func(s capnp.Struct) error { return params(PingPong_echoNum_Params{Struct: s}) }
	}
	ans, release := c.Client.SendCall(ctx, s)
	return PingPong_echoNum_Results_Future{Future: ans.Future()}, release
}

func (c PingPong) AddRef() PingPong {
	return PingPong{
		Client: c.Client.AddRef(),
	}
}

func (c PingPong) Release() {
	c.Client.Release()
}

// A PingPong_Server is a PingPong with a local implementation.
type PingPong_Server interface {
	EchoNum(context.Context, PingPong_echoNum) error
}

// PingPong_NewServer creates a new Server from an implementation of PingPong_Server.
func PingPong_NewServer(s PingPong_Server, policy *server.Policy) *server.Server {
	c, _ := s.(server.Shutdowner)
	return server.New(PingPong_Methods(nil, s), s, c, policy)
}

// PingPong_ServerToClient creates a new Client from an implementation of PingPong_Server.
// The caller is responsible for calling Release on the returned Client.
func PingPong_ServerToClient(s PingPong_Server, policy *server.Policy) PingPong {
	return PingPong{Client: capnp.NewClient(PingPong_NewServer(s, policy))}
}

// PingPong_Methods appends Methods to a slice that invoke the methods on s.
// This can be used to create a more complicated Server.
func PingPong_Methods(methods []server.Method, s PingPong_Server) []server.Method {
	if cap(methods) == 0 {
		methods = make([]server.Method, 0, 1)
	}

	methods = append(methods, server.Method{
		Method: capnp.Method{
			InterfaceID:   0xf004c474c2f8ee7a,
			MethodID:      0,
			InterfaceName: "test.capnp:PingPong",
			MethodName:    "echoNum",
		},
		Impl: func(ctx context.Context, call *server.Call) error {
			return s.EchoNum(ctx, PingPong_echoNum{call})
		},
	})

	return methods
}

// PingPong_echoNum holds the state for a server call to PingPong.echoNum.
// See server.Call for documentation.
type PingPong_echoNum struct {
	*server.Call
}

// Args returns the call's arguments.
func (c PingPong_echoNum) Args() PingPong_echoNum_Params {
	return PingPong_echoNum_Params{Struct: c.Call.Args()}
}

// AllocResults allocates the results struct.
func (c PingPong_echoNum) AllocResults() (PingPong_echoNum_Results, error) {
	r, err := c.Call.AllocResults(capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	return PingPong_echoNum_Results{Struct: r}, err
}

type PingPong_echoNum_Params struct{ capnp.Struct }

// PingPong_echoNum_Params_TypeID is the unique identifier for the type PingPong_echoNum_Params.
const PingPong_echoNum_Params_TypeID = 0xd797e0a99edf0921

func NewPingPong_echoNum_Params(s *capnp.Segment) (PingPong_echoNum_Params, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	return PingPong_echoNum_Params{st}, err
}

func NewRootPingPong_echoNum_Params(s *capnp.Segment) (PingPong_echoNum_Params, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	return PingPong_echoNum_Params{st}, err
}

func ReadRootPingPong_echoNum_Params(msg *capnp.Message) (PingPong_echoNum_Params, error) {
	root, err := msg.Root()
	return PingPong_echoNum_Params{root.Struct()}, err
}

func (s PingPong_echoNum_Params) String() string {
	str, _ := text.Marshal(0xd797e0a99edf0921, s.Struct)
	return str
}

func (s PingPong_echoNum_Params) N() int64 {
	return int64(s.Struct.Uint64(0))
}

func (s PingPong_echoNum_Params) SetN(v int64) {
	s.Struct.SetUint64(0, uint64(v))
}

// PingPong_echoNum_Params_List is a list of PingPong_echoNum_Params.
type PingPong_echoNum_Params_List struct{ capnp.List }

// NewPingPong_echoNum_Params creates a new list of PingPong_echoNum_Params.
func NewPingPong_echoNum_Params_List(s *capnp.Segment, sz int32) (PingPong_echoNum_Params_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0}, sz)
	return PingPong_echoNum_Params_List{l}, err
}

func (s PingPong_echoNum_Params_List) At(i int) PingPong_echoNum_Params {
	return PingPong_echoNum_Params{s.List.Struct(i)}
}

func (s PingPong_echoNum_Params_List) Set(i int, v PingPong_echoNum_Params) error {
	return s.List.SetStruct(i, v.Struct)
}

func (s PingPong_echoNum_Params_List) String() string {
	str, _ := text.MarshalList(0xd797e0a99edf0921, s.List)
	return str
}

// PingPong_echoNum_Params_Future is a wrapper for a PingPong_echoNum_Params promised by a client call.
type PingPong_echoNum_Params_Future struct{ *capnp.Future }

func (p PingPong_echoNum_Params_Future) Struct() (PingPong_echoNum_Params, error) {
	s, err := p.Future.Struct()
	return PingPong_echoNum_Params{s}, err
}

type PingPong_echoNum_Results struct{ capnp.Struct }

// PingPong_echoNum_Results_TypeID is the unique identifier for the type PingPong_echoNum_Results.
const PingPong_echoNum_Results_TypeID = 0x85ddfd96db252600

func NewPingPong_echoNum_Results(s *capnp.Segment) (PingPong_echoNum_Results, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	return PingPong_echoNum_Results{st}, err
}

func NewRootPingPong_echoNum_Results(s *capnp.Segment) (PingPong_echoNum_Results, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	return PingPong_echoNum_Results{st}, err
}

func ReadRootPingPong_echoNum_Results(msg *capnp.Message) (PingPong_echoNum_Results, error) {
	root, err := msg.Root()
	return PingPong_echoNum_Results{root.Struct()}, err
}

func (s PingPong_echoNum_Results) String() string {
	str, _ := text.Marshal(0x85ddfd96db252600, s.Struct)
	return str
}

func (s PingPong_echoNum_Results) N() int64 {
	return int64(s.Struct.Uint64(0))
}

func (s PingPong_echoNum_Results) SetN(v int64) {
	s.Struct.SetUint64(0, uint64(v))
}

// PingPong_echoNum_Results_List is a list of PingPong_echoNum_Results.
type PingPong_echoNum_Results_List struct{ capnp.List }

// NewPingPong_echoNum_Results creates a new list of PingPong_echoNum_Results.
func NewPingPong_echoNum_Results_List(s *capnp.Segment, sz int32) (PingPong_echoNum_Results_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0}, sz)
	return PingPong_echoNum_Results_List{l}, err
}

func (s PingPong_echoNum_Results_List) At(i int) PingPong_echoNum_Results {
	return PingPong_echoNum_Results{s.List.Struct(i)}
}

func (s PingPong_echoNum_Results_List) Set(i int, v PingPong_echoNum_Results) error {
	return s.List.SetStruct(i, v.Struct)
}

func (s PingPong_echoNum_Results_List) String() string {
	str, _ := text.MarshalList(0x85ddfd96db252600, s.List)
	return str
}

// PingPong_echoNum_Results_Future is a wrapper for a PingPong_echoNum_Results promised by a client call.
type PingPong_echoNum_Results_Future struct{ *capnp.Future }

func (p PingPong_echoNum_Results_Future) Struct() (PingPong_echoNum_Results, error) {
	s, err := p.Future.Struct()
	return PingPong_echoNum_Results{s}, err
}

type StreamTest struct{ Client *capnp.Client }

// StreamTest_TypeID is the unique identifier for the type StreamTest.
const StreamTest_TypeID = 0xbb3ca85b01eea465

func (c StreamTest) Push(ctx context.Context, params func(StreamTest_push_Params) error) (stream.StreamResult_Future, capnp.ReleaseFunc) {
	s := capnp.Send{
		Method: capnp.Method{
			InterfaceID:   0xbb3ca85b01eea465,
			MethodID:      0,
			InterfaceName: "test.capnp:StreamTest",
			MethodName:    "push",
		},
	}
	if params != nil {
		s.ArgsSize = capnp.ObjectSize{DataSize: 0, PointerCount: 1}
		s.PlaceArgs = func(s capnp.Struct) error { return params(StreamTest_push_Params{Struct: s}) }
	}
	ans, release := c.Client.SendCall(ctx, s)
	return stream.StreamResult_Future{Future: ans.Future()}, release
}

func (c StreamTest) AddRef() StreamTest {
	return StreamTest{
		Client: c.Client.AddRef(),
	}
}

func (c StreamTest) Release() {
	c.Client.Release()
}

// A StreamTest_Server is a StreamTest with a local implementation.
type StreamTest_Server interface {
	Push(context.Context, StreamTest_push) error
}

// StreamTest_NewServer creates a new Server from an implementation of StreamTest_Server.
func StreamTest_NewServer(s StreamTest_Server, policy *server.Policy) *server.Server {
	c, _ := s.(server.Shutdowner)
	return server.New(StreamTest_Methods(nil, s), s, c, policy)
}

// StreamTest_ServerToClient creates a new Client from an implementation of StreamTest_Server.
// The caller is responsible for calling Release on the returned Client.
func StreamTest_ServerToClient(s StreamTest_Server, policy *server.Policy) StreamTest {
	return StreamTest{Client: capnp.NewClient(StreamTest_NewServer(s, policy))}
}

// StreamTest_Methods appends Methods to a slice that invoke the methods on s.
// This can be used to create a more complicated Server.
func StreamTest_Methods(methods []server.Method, s StreamTest_Server) []server.Method {
	if cap(methods) == 0 {
		methods = make([]server.Method, 0, 1)
	}

	methods = append(methods, server.Method{
		Method: capnp.Method{
			InterfaceID:   0xbb3ca85b01eea465,
			MethodID:      0,
			InterfaceName: "test.capnp:StreamTest",
			MethodName:    "push",
		},
		Impl: func(ctx context.Context, call *server.Call) error {
			return s.Push(ctx, StreamTest_push{call})
		},
	})

	return methods
}

// StreamTest_push holds the state for a server call to StreamTest.push.
// See server.Call for documentation.
type StreamTest_push struct {
	*server.Call
}

// Args returns the call's arguments.
func (c StreamTest_push) Args() StreamTest_push_Params {
	return StreamTest_push_Params{Struct: c.Call.Args()}
}

// AllocResults allocates the results struct.
func (c StreamTest_push) AllocResults() (stream.StreamResult, error) {
	r, err := c.Call.AllocResults(capnp.ObjectSize{DataSize: 0, PointerCount: 0})
	return stream.StreamResult{Struct: r}, err
}

type StreamTest_push_Params struct{ capnp.Struct }

// StreamTest_push_Params_TypeID is the unique identifier for the type StreamTest_push_Params.
const StreamTest_push_Params_TypeID = 0xf838dca6c8721bdb

func NewStreamTest_push_Params(s *capnp.Segment) (StreamTest_push_Params, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return StreamTest_push_Params{st}, err
}

func NewRootStreamTest_push_Params(s *capnp.Segment) (StreamTest_push_Params, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return StreamTest_push_Params{st}, err
}

func ReadRootStreamTest_push_Params(msg *capnp.Message) (StreamTest_push_Params, error) {
	root, err := msg.Root()
	return StreamTest_push_Params{root.Struct()}, err
}

func (s StreamTest_push_Params) String() string {
	str, _ := text.Marshal(0xf838dca6c8721bdb, s.Struct)
	return str
}

func (s StreamTest_push_Params) Data() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	return []byte(p.Data()), err
}

func (s StreamTest_push_Params) HasData() bool {
	return s.Struct.HasPtr(0)
}

func (s StreamTest_push_Params) SetData(v []byte) error {
	return s.Struct.SetData(0, v)
}

// StreamTest_push_Params_List is a list of StreamTest_push_Params.
type StreamTest_push_Params_List struct{ capnp.List }

// NewStreamTest_push_Params creates a new list of StreamTest_push_Params.
func NewStreamTest_push_Params_List(s *capnp.Segment, sz int32) (StreamTest_push_Params_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return StreamTest_push_Params_List{l}, err
}

func (s StreamTest_push_Params_List) At(i int) StreamTest_push_Params {
	return StreamTest_push_Params{s.List.Struct(i)}
}

func (s StreamTest_push_Params_List) Set(i int, v StreamTest_push_Params) error {
	return s.List.SetStruct(i, v.Struct)
}

func (s StreamTest_push_Params_List) String() string {
	str, _ := text.MarshalList(0xf838dca6c8721bdb, s.List)
	return str
}

// StreamTest_push_Params_Future is a wrapper for a StreamTest_push_Params promised by a client call.
type StreamTest_push_Params_Future struct{ *capnp.Future }

func (p StreamTest_push_Params_Future) Struct() (StreamTest_push_Params, error) {
	s, err := p.Future.Struct()
	return StreamTest_push_Params{s}, err
}

const schema_ef12a34b9807e19c = "x\xda\x12\xf8\xe4\xc0d\xc8Z\xcf\xc2\xc0\x10h\xc2\xca" +
	"\xf6OM\xf5\xf6\xb4\xbfw[\x03E\x18\x19\x19\x18X" +
	"\xd8\x19\x18\x8cU\x99\x94\x18\x19\x18\x85u\x99\xec\x19\x18" +
	"\xff\xa7.y\xc7\x18\xbd\xc2f7\x83 7\xf3\xff9" +
	"\x0f\xd9gx/\x16z\xcf\xc0\xc0(\xec\xcb\xb4I8" +
	"\x94\x89\x9d\x81A8\x90\xc9]\xb8\x12\xc4\xfa\xaf\xc8y" +
	"\x7f\xde\xca\x07\xd3\xaf3 \x99\x96\xc8$\x052-\x13" +
	"lZ\xd5\xbb\x1f\x87J\x8e\xb0|\xc00\xad\x93i\x91" +
	"\xf0D\xb0i\xbdL\xee\xc2[\xc1\xa6\xdd\x96.:\xb1" +
	"\xec\x8e\xc5\x0f\x06A1F\x06\x06VF\x90is\x99" +
	"\x84@\xa6-e\xb2g\x88\xfc_\x92Z\\\xa2\x97\x9c" +
	"X\xc0\x9cW`\x15\x90\x99\x97\x1e\x90\x9f\x97\xae\x97\x9a" +
	"\x9c\x91\xefW\x9a\xab\x12\x94Z\\\xca\x9eSR\x1c\xc8" +
	"\xc2\xcc\xc2\xc0\xc0\xc2\xc8\xc0 \xc8+\xc4\xc0\x10\xc8\xc1" +
	"\xcc\x18(\xc2\xc4\xc8\x98\xc7\xc8\xca\xc0\xc4\xc8\xca\xc0\x08" +
	"7\x861\xaf\xc0*\xb8\xa4(5Q>7$\xb5\xb8" +
	"$\x80\x911\x90\x85\x99\x15\xc9!\x8cy\x1b\x0f\x94\x1b" +
	"\xcf\x8a\x9f)(\xa8\xc5\xc0$\xc8\xca\xce_PZ\x9c" +
	"\xe1\xc0\x18\xc0\xc8\x88\xdf-\x01\x89E\x89\xcc\xb9$:" +
	"\x05l\x0a{~^:\xc2!\xb0\xf0ed\x80\xc6\x9a" +
	"\xa0\xa0\x13\xd8!\xf5P\x9b0\xdd\x02\xf6\x10\xd8?z" +
	" \xc7\x82\x9d\x92\xcb\x88\xe2\x14-\x84S\xf8S\x12K" +
	"\x12\x19y\x19\x98\x18y\x19\x18\x01\x01\x00\x00\xff\xff)" +
	"\x11\x93\xee"

func init() {
	schemas.Register(schema_ef12a34b9807e19c,
		0x85ddfd96db252600,
		0xbb3ca85b01eea465,
		0xd797e0a99edf0921,
		0xf004c474c2f8ee7a,
		0xf838dca6c8721bdb)
}
