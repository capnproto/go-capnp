package capnp

import (
	"errors"
	"strconv"

	"golang.org/x/net/context"
)

var ErrNullClient = errors.New("capn: call on null client")

// A Client represents an Cap'n Proto interface type.
//
// Generally, only RPC protocol implementers should provide types that
// implement Client: call ordering guarantees, promises, and
// synchronization are tricky to get right.  Prefer creating a server
// that wraps another interface than trying to implement Client.
type Client interface {
	// Call starts executing a method and returns an answer that will hold
	// the resulting struct.  The call's parameters must be placed before
	// Call() returns.
	//
	// Calls are delivered to the capability in the order they are made.
	// This guarantee is based on the concept of a capability
	// acknowledging delivery of a call: this is specific to an
	// implementation of Client.  A type that implements Client must
	// guarantee that if foo() then bar() is called on a client, that
	// acknowledging foo() happens before acknowledging bar().
	Call(call *Call) Answer

	// Close releases any resources associated with this client.
	// No further calls to the client should be made after calling Close.
	Close() error
}

// The Call type holds the record for an outgoing interface call.
type Call struct {
	// Ctx is the context of the call.
	Ctx context.Context

	// Method is the interface ID and method ID, along with the optional name,
	// of the method to call.
	Method Method

	// Params is a struct containing parameters for the call.
	// This should be set when the RPC system receives a call for an
	// exported interface.  It is mutually exclusive with ParamsFunc
	// and ParamsSize.
	Params Struct
	// ParamsFunc is a function that populates an allocated struct with
	// the parameters for the call.  ParamsSize determines the size of the
	// struct to allocate.  This is used when application code is using a
	// client.  These settings should be set together; they are mutually
	// exclusive with Params.
	ParamsFunc func(Struct)
	ParamsSize ObjectSize

	// Options passes RPC-specific options for the call.
	Options CallOptions
}

// Copy clones a call, ensuring that its Params are placed.
// If Call.ParamsFunc is nil, then the same pointer will be returned.
func (call *Call) Copy(s *Segment) *Call {
	if call.ParamsFunc == nil {
		return call
	}
	return &Call{
		Ctx:     call.Ctx,
		Method:  call.Method,
		Params:  call.PlaceParams(s),
		Options: call.Options,
	}
}

// PlaceParams returns the parameters struct, allocating it inside
// segment s as necessary.  If s is nil, a new segment is allocated.
func (call *Call) PlaceParams(s *Segment) Struct {
	if call.ParamsFunc == nil {
		return call.Params
	}
	if s == nil {
		s = NewBuffer(nil)
	}
	p := s.NewStruct(call.ParamsSize)
	call.ParamsFunc(p)
	return p
}

// CallOptions holds RPC-specific options for an interface call.
// Its usage is similar to the values in context.Context, but is only
// used for a single call: its values are not intended to propagate to
// other callees.  An example of an option would be the
// Call.sendResultsTo field in rpc.capnp.
type CallOptions struct {
	m map[interface{}]interface{}
}

// NewCallOptions builds a CallOptions value from a list of individual options.
func NewCallOptions(opts []CallOption) CallOptions {
	co := CallOptions{make(map[interface{}]interface{})}
	for _, o := range opts {
		o.f(co)
	}
	return co
}

// Value retrieves the value associated with the options for this key,
// or nil if no value is associated with this key.
func (co CallOptions) Value(key interface{}) interface{} {
	return co.m[key]
}

// With creates a copy of the CallOptions value with other options applied.
func (co CallOptions) With(opts []CallOption) CallOptions {
	newopts := CallOptions{make(map[interface{}]interface{})}
	for k, v := range co.m {
		newopts.m[k] = v
	}
	for _, o := range opts {
		o.f(newopts)
	}
	return newopts
}

// A CallOption is a function that modifies options on an interface call.
type CallOption struct {
	f func(CallOptions)
}

// SetOptionValue returns a call option that associates a value to an
// option key.  This can be retrieved later with CallOptions.Value.
func SetOptionValue(key, value interface{}) CallOption {
	return CallOption{func(co CallOptions) {
		co.m[key] = value
	}}
}

// callOptionKey is the unexported key type for predefined options.
type callOptionKey int

// Predefined call options
const (
	invalidOptionKey callOptionKey = iota
	ackSignalKey
)

// An Answer is the deferred result of a client call, which is usually wrapped by a Pipeline.
type Answer interface {
	// Struct waits until the call is finished and returns the result.
	Struct() (Struct, error)

	// The following methods are the same as in Client except with
	// an added transform parameter -- a path to the interface to use.

	PipelineCall(transform []PipelineOp, call *Call) Answer
	PipelineClose(transform []PipelineOp) error
}

// A Pipeline is a generic wrapper for an answer.
type Pipeline struct {
	answer Answer
	parent *Pipeline
	op     PipelineOp
}

// NewPipeline returns a new pipeline based on an answer.
func NewPipeline(ans Answer) *Pipeline {
	return &Pipeline{answer: ans}
}

// Answer returns the answer the pipeline is derived from.
func (p *Pipeline) Answer() Answer {
	return p.answer
}

// Transform returns the operations needed to transform the root answer
// into the value p represents.
func (p *Pipeline) Transform() []PipelineOp {
	n := 0
	for q := p; q.parent != nil; q = q.parent {
		n++
	}
	xform := make([]PipelineOp, n)
	for i, q := n-1, p; q.parent != nil; i, q = i-1, q.parent {
		xform[i] = q.op
	}
	return xform
}

// Struct waits until the answer is resolved and returns the struct
// this pipeline represents.
func (p *Pipeline) Struct() (Struct, error) {
	s, err := p.answer.Struct()
	if err != nil {
		return Struct{}, err
	}
	return TransformObject(Object(s), p.Transform()).ToStruct(), err
}

// Client returns the client version of p.
func (p *Pipeline) Client() *PipelineClient {
	return (*PipelineClient)(p)
}

// GetPipeline returns a derived pipeline which yields the pointer field given.
func (p *Pipeline) GetPipeline(off int) *Pipeline {
	return p.GetPipelineDefault(off, nil, 0)
}

// GetPipelineDefault returns a derived pipeline which yields the pointer field given,
// defaulting to the value given.
func (p *Pipeline) GetPipelineDefault(off int, dseg *Segment, doff int) *Pipeline {
	return &Pipeline{
		answer: p.answer,
		parent: p,
		op: PipelineOp{
			Field:          off,
			DefaultSegment: dseg,
			DefaultOffset:  doff,
		},
	}
}

// PipelineClient implements Client by calling to the pipeline's answer.
type PipelineClient Pipeline

func (pc *PipelineClient) transform() []PipelineOp {
	return (*Pipeline)(pc).Transform()
}

func (pc *PipelineClient) Call(call *Call) Answer {
	return pc.answer.PipelineCall(pc.transform(), call)
}

func (pc *PipelineClient) Close() error {
	return pc.answer.PipelineClose(pc.transform())
}

// A PipelineOp describes a step in transforming a pipeline.
// It maps closely with the PromisedAnswer.Op struct in rpc.capnp.
type PipelineOp struct {
	Field          int
	DefaultSegment *Segment
	DefaultOffset  int
}

// String returns a human-readable description of op.
func (op PipelineOp) String() string {
	s := make([]byte, 0, 32)
	s = append(s, "get field "...)
	s = strconv.AppendInt(s, int64(op.Field), 10)
	if op.DefaultSegment == nil {
		return string(s)
	}
	s = append(s, " with default"...)
	return string(s)
}

// A Method identifies a method along with an optional human-readable
// description of the method.
type Method struct {
	InterfaceID uint64
	MethodID    uint16

	// Canonical name of the interface.  May be empty.
	InterfaceName string
	// Method name as it appears in the schema.  May be empty.
	MethodName string
}

// String returns a formatted string containing the interface name or
// the method name if present, otherwise it uses the raw IDs.
// This is suitable for use in error messages and logs.
func (m *Method) String() string {
	buf := make([]byte, 0, 128)
	if m.InterfaceName == "" {
		buf = append(buf, '@', '0', 'x')
		buf = strconv.AppendUint(buf, m.InterfaceID, 16)
	} else {
		buf = append(buf, m.InterfaceName...)
	}
	buf = append(buf, '.')
	if m.MethodName == "" {
		buf = append(buf, '@')
		buf = strconv.AppendUint(buf, uint64(m.MethodID), 10)
	} else {
		buf = append(buf, m.MethodName...)
	}
	return string(buf)
}

// TransformObject applies a sequence of pipeline operations to an object
// and returns the result.
func TransformObject(p Object, transform []PipelineOp) Object {
	n := len(transform)
	if n == 0 {
		return p
	}
	s := p.ToStruct()
	for _, op := range transform[:n-1] {
		field := s.GetObject(op.Field)
		if op.DefaultSegment == nil {
			s = field.ToStruct()
		} else {
			s = field.ToStructDefault(op.DefaultSegment, op.DefaultOffset)
		}
	}
	op := transform[n-1]
	p = s.GetObject(op.Field)
	if op.DefaultSegment != nil {
		p = Object(p.ToStructDefault(op.DefaultSegment, op.DefaultOffset))
	}
	return p
}

type immediateAnswer Struct

// ImmediateAnswer returns an Answer that accesses s.
func ImmediateAnswer(s Struct) Answer {
	return immediateAnswer(s)
}

func (ans immediateAnswer) Struct() (Struct, error) {
	return Struct(ans), nil
}

func (ans immediateAnswer) PipelineCall(transform []PipelineOp, call *Call) Answer {
	c := TransformObject(Object(ans), transform).ToInterface().Client()
	if c == nil {
		return ErrorAnswer(ErrNullClient)
	}
	return c.Call(call)
}

func (ans immediateAnswer) PipelineClose(transform []PipelineOp) error {
	c := TransformObject(Object(ans), transform).ToInterface().Client()
	if c == nil {
		return ErrNullClient
	}
	return c.Close()
}

type errorAnswer struct {
	e error
}

// ErrorAnswer returns a Answer that always returns error e.
func ErrorAnswer(e error) Answer {
	return errorAnswer{e}
}

func (ans errorAnswer) Struct() (Struct, error) {
	return Struct{}, ans.e
}

func (ans errorAnswer) PipelineCall([]PipelineOp, *Call) Answer {
	return ans
}

func (ans errorAnswer) PipelineClose([]PipelineOp) error {
	return ans.e
}

// IsFixedAnswer reports whether an answer was created by
// ImmediateAnswer or ErrorAnswer.
func IsFixedAnswer(ans Answer) bool {
	switch ans.(type) {
	case immediateAnswer:
		return true
	case errorAnswer:
		return true
	default:
		return false
	}
}

type errorClient struct {
	e error
}

// ErrorClient returns a Client that always returns error e.
func ErrorClient(e error) Client {
	return errorClient{e}
}

func (ec errorClient) Call(*Call) Answer {
	return ErrorAnswer(ec.e)
}

func (ec errorClient) Close() error {
	return nil
}

// IsErrorClient reports whether c was created with ErrorClient.
func IsErrorClient(c Client) bool {
	_, ok := c.(errorClient)
	return ok
}
