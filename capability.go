package capn

import (
	"errors"
	"strconv"

	"golang.org/x/net/context"
)

var ErrNullClient = errors.New("capn: call on null client")

// A Client makes calls for an interface type.
type Client interface {
	// Call calls a method with an existing parameters struct.  This is
	// used when the RPC system receives a call for an exported object.
	Call(
		ctx context.Context,
		method *Method,
		params Struct) Answer

	// NewCall calls a method, allowing the client to allocate the
	// parameters struct.  This is used when application code is using
	// a client.
	NewCall(
		ctx context.Context,
		method *Method,
		paramsSize ObjectSize,
		params func(Struct)) Answer

	// Close releases any resources associated with this client.
	// No further calls to the client should be made after calling Close.
	Close() error
}

// An Answer is the deferred result of a client call, which is usually wrapped by a Promise.
type Answer interface {
	// Struct waits until the call is finished and returns the result.
	Struct() (Struct, error)

	// The following methods are the same as in Client except with an added transform parameter.
	// This parameter describes the path to the interface to use for the call.

	Call(
		ctx context.Context,
		transform []PromiseOp,
		method *Method,
		params Struct) Answer
	NewCall(
		ctx context.Context,
		transform []PromiseOp,
		method *Method,
		paramsSize ObjectSize,
		params func(Struct)) Answer
	CloseClient(transform []PromiseOp) error
}

// A Promise is a generic promise for an object.
type Promise struct {
	answer Answer
	parent *Promise
	op     PromiseOp
}

// NewPromise returns a new promise based on an answer.
func NewPromise(ans Answer) *Promise {
	return &Promise{answer: ans}
}

// Answer returns the answer the promise is derived from.
func (p *Promise) Answer() Answer {
	return p.answer
}

// Transform returns the operations needed to transform the root answer
// into the value p represents.
func (p *Promise) Transform() []PromiseOp {
	n := 0
	for q := p; q.parent != nil; q = q.parent {
		n++
	}
	xform := make([]PromiseOp, n)
	for i, q := n-1, p; q.parent != nil; i, q = i-1, q.parent {
		xform[i] = q.op
	}
	return xform
}

// Struct waits until the answer is resolved and returns the struct
// this promise represents.
func (p *Promise) Struct() (Struct, error) {
	s, err := p.answer.Struct()
	if err != nil {
		return Struct{}, err
	}
	return walk(Object(s), p.Transform()).ToStruct(), err
}

// Client returns the client version of p.
func (p *Promise) Client() *PromiseClient {
	return (*PromiseClient)(p)
}

// GetPromise returns a derived promise which yields the pointer field given.
func (p *Promise) GetPromise(off int) *Promise {
	return p.GetPromiseDefault(off, nil, 0)
}

// GetPromiseDefault returns a derived promise which yields the pointer field given,
// defaulting to the value given.
func (p *Promise) GetPromiseDefault(off int, dseg *Segment, doff int) *Promise {
	return &Promise{
		answer: p.answer,
		parent: p,
		op: PromiseOp{
			Field:          off,
			DefaultSegment: dseg,
			DefaultOffset:  doff,
		},
	}
}

// PromiseClient implements Client by calling to the promise's answer.
type PromiseClient Promise

func (pc *PromiseClient) transform() []PromiseOp {
	return (*Promise)(pc).Transform()
}

func (pc *PromiseClient) NewCall(ctx context.Context, method *Method, paramsSize ObjectSize, params func(Struct)) Answer {
	return pc.answer.NewCall(ctx, pc.transform(), method, paramsSize, params)
}

func (pc *PromiseClient) Call(ctx context.Context, method *Method, params Struct) Answer {
	return pc.answer.Call(ctx, pc.transform(), method, params)
}

func (pc *PromiseClient) Close() error {
	return pc.answer.CloseClient(pc.transform())
}

// A PromiseOp describes a step in transforming a promise.
// It maps closely with the PromisedAnswer.Op struct in rpc.capnp.
type PromiseOp struct {
	Field          int
	DefaultSegment *Segment
	DefaultOffset  int
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
	buf = append(buf, '/')
	if m.MethodName == "" {
		buf = append(buf, '@')
		buf = strconv.AppendUint(buf, uint64(m.MethodID), 10)
	} else {
		buf = append(buf, m.MethodName...)
	}
	return string(buf)
}

func walk(p Object, transform []PromiseOp) Object {
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

type immediateAnswer Object

// ImmediateAnswer returns an Answer that accesses s.
func ImmediateAnswer(s Object) Answer {
	return immediateAnswer(s)
}

func (ans immediateAnswer) Struct() (Struct, error) {
	return Struct(ans), nil
}

func (ans immediateAnswer) Call(ctx context.Context, transform []PromiseOp, method *Method, params Struct) Answer {
	c := walk(Object(ans), transform).ToInterface().Client()
	if c == nil {
		return ErrorAnswer(ErrNullClient)
	}
	return c.Call(ctx, method, params)
}

func (ans immediateAnswer) NewCall(ctx context.Context, transform []PromiseOp, method *Method, paramsSize ObjectSize, params func(Struct)) Answer {
	c := walk(Object(ans), transform).ToInterface().Client()
	if c == nil {
		return ErrorAnswer(ErrNullClient)
	}
	return c.NewCall(ctx, method, paramsSize, params)
}

func (ans immediateAnswer) CloseClient(transform []PromiseOp) error {
	c := walk(Object(ans), transform).ToInterface().Client()
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

func (ans errorAnswer) Call(ctx context.Context, transform []PromiseOp, method *Method, params Struct) Answer {
	return ans
}

func (ans errorAnswer) NewCall(ctx context.Context, transform []PromiseOp, method *Method, paramsSize ObjectSize, params func(Struct)) Answer {
	return ans
}

func (ans errorAnswer) CloseClient(transform []PromiseOp) error {
	return ans.e
}

type errorClient struct {
	e error
}

// ErrorClient returns a Client that always returns error e.
func ErrorClient(e error) Client {
	return errorClient{e}
}

func (ec errorClient) NewCall(ctx context.Context, method *Method, paramsSize ObjectSize, params func(Struct)) Answer {
	return ErrorAnswer(ec.e)
}

func (ec errorClient) Call(ctx context.Context, method *Method, params Struct) Answer {
	return ErrorAnswer(ec.e)
}

func (ec errorClient) Close() error {
	return nil
}
