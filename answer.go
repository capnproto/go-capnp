package capnp

import (
	"context"
	"strconv"
	"sync"
)

// A Promise holds the result of an RPC call.  Only one result can be
// written to a promise.  Before the result is written, calls can be
// queued up using the Answer methods — this is promise pipelining.
//
// Promise is most useful for implementing ClientHook.
// Most applications will use Answer, since that what is returned by
// a Client.
type Promise struct {
	resolved chan struct{} // can only be closed while holding mu

	mu     sync.RWMutex
	caller PipelineCaller // nil if resolved
	result Ptr
	err    error
}

// NewPromise creates a new unresolved promise.  The PipelineCaller will
// be used to make pipelined calls before the promise resolves.
func NewPromise(pc PipelineCaller) *Promise {
	return &Promise{
		caller:   pc,
		resolved: make(chan struct{}),
	}
}

// Fulfill resolves the promise with a successful result.
func (p *Promise) Fulfill(result Ptr) {
	defer p.mu.Unlock()
	p.mu.Lock()
	select {
	case <-p.resolved:
		panic("Promise.Fulfill or Promise.Reject called more than once")
	default:
		p.result = result
		p.caller = nil
		close(p.resolved)
	}
}

// Reject resolves the promise with a failure.
func (p *Promise) Reject(e error) {
	if e == nil {
		panic("Promise.Reject(nil)")
	}
	defer p.mu.Unlock()
	p.mu.Lock()
	select {
	case <-p.resolved:
		panic("Promise.Fulfill or Promise.Reject called more than once")
	default:
		p.err = e
		p.caller = nil
		close(p.resolved)
	}
}

// Answer returns a read-only view of the promise.
func (p *Promise) Answer() *Answer {
	return &Answer{promise: p}
}

func (p *Promise) client(t []PipelineOp) *Client {
	p.mu.RLock()
	if p.caller == nil {
		c := clientFromResolution(p.result, p.err, t)
		p.mu.RUnlock()
		return c
	}
	p.mu.RUnlock()
	p.mu.Lock()
	if p.caller == nil {
		c := clientFromResolution(p.result, p.err, t)
		p.mu.Unlock()
		return c
	}
	// TODO(soon): cache clients in a table.
	c, _ := NewPromisedClient(pipelineClient{
		p:         p,
		transform: t,
	})
	p.mu.Unlock()
	return c
}

// A PipelineCaller implements promise pipelining.
//
// See the counterpart methods in ClientHook for a description.
type PipelineCaller interface {
	PipelineSend(ctx context.Context, transform []PipelineOp, m Method, a SendArgs, opts CallOptions) (*Answer, ReleaseFunc)
	PipelineRecv(ctx context.Context, transform []PipelineOp, m Method, a RecvArgs, opts CallOptions) (*Answer, ReleaseFunc)
}

// An Answer is a deferred result of a client call.  Conceptually, this is a
// future.  It is safe to use from multiple goroutines.
type Answer struct {
	promise *Promise
	parent  *Answer // nil if root answer
	op      PipelineOp
}

// ErrorAnswer returns a Answer that always returns error e.
func ErrorAnswer(e error) *Answer {
	return &Answer{
		promise: &Promise{
			resolved: newClosedSignal(),
			err:      e,
		},
	}
}

// ImmediateAnswer returns an Answer that accesses s.
func ImmediateAnswer(s Struct) *Answer {
	return &Answer{
		promise: &Promise{
			resolved: newClosedSignal(),
			result:   s.ToPtr(),
		},
	}
}

// transform returns the operations needed to transform the root answer
// into the value p represents.
func (ans *Answer) transform() []PipelineOp {
	n := 0
	for a := ans; a.parent != nil; a = a.parent {
		n++
	}
	xform := make([]PipelineOp, n)
	for i, a := n-1, ans; a.parent != nil; i, a = i-1, a.parent {
		xform[i] = a.op
	}
	return xform
}

// Done returns a channel that is closed when the answer's call is finished.
func (ans *Answer) Done() <-chan struct{} {
	return ans.promise.resolved
}

// Struct waits until the answer is resolved and returns the struct
// this answer represents.
func (ans *Answer) Struct() (Struct, error) {
	<-ans.promise.resolved
	ans.promise.mu.RLock()
	ptr, err := ans.promise.result, ans.promise.err
	ans.promise.mu.RUnlock()
	if err != nil {
		return Struct{}, err
	}
	ptr, err = Transform(ptr, ans.transform())
	if err != nil {
		return Struct{}, err
	}
	return ptr.Struct(), nil
}

// Client returns the answer as a client.  If the answer's originating
// call has not completed, then calls will be queued until the original
// call's completion.
func (ans *Answer) Client() *Client {
	return ans.promise.client(ans.transform())
}

// Field returns a derived pipeline which yields the pointer field given,
// defaulting to the value given.
func (ans *Answer) Field(off uint16, def []byte) *Answer {
	return &Answer{
		promise: ans.promise,
		parent:  ans,
		op: PipelineOp{
			Field:        off,
			DefaultValue: def,
		},
	}
}

// pipelineClient implements ClientHook by calling to the pipeline's answer.
type pipelineClient struct {
	p         *Promise
	transform []PipelineOp
}

func (pc pipelineClient) Send(ctx context.Context, m Method, a SendArgs, opts CallOptions) (*Answer, ReleaseFunc) {
	defer pc.p.mu.RUnlock()
	pc.p.mu.RLock()
	if pc.p.caller != nil {
		return pc.p.caller.PipelineSend(ctx, pc.transform, m, a, opts)
	}
	c := clientFromResolution(pc.p.result, pc.p.err, pc.transform)
	return c.SendCall(ctx, m, a, opts)
}

func (pc pipelineClient) Recv(ctx context.Context, m Method, a RecvArgs, opts CallOptions) (*Answer, ReleaseFunc) {
	defer pc.p.mu.RUnlock()
	pc.p.mu.RLock()
	if pc.p.caller != nil {
		return pc.p.caller.PipelineRecv(ctx, pc.transform, m, a, opts)
	}
	c := clientFromResolution(pc.p.result, pc.p.err, pc.transform)
	return c.RecvCall(ctx, m, a, opts)
}

func (pc pipelineClient) Brand() interface{} {
	defer pc.p.mu.RUnlock()
	pc.p.mu.RLock()
	if pc.p.caller != nil {
		// TODO(someday): allow people to obtain the underlying answer.
		return nil
	}
	c := clientFromResolution(pc.p.result, pc.p.err, pc.transform)
	return c.Brand()
}

// A PipelineOp describes a step in transforming a pipeline.
// It maps closely with the PromisedAnswer.Op struct in rpc.capnp.
type PipelineOp struct {
	Field        uint16
	DefaultValue []byte
}

// String returns a human-readable description of op.
func (op PipelineOp) String() string {
	s := make([]byte, 0, 32)
	s = append(s, "get field "...)
	s = strconv.AppendInt(s, int64(op.Field), 10)
	if op.DefaultValue == nil {
		return string(s)
	}
	s = append(s, " with default"...)
	return string(s)
}

// Transform applies a sequence of pipeline operations to a pointer
// and returns the result.
func Transform(p Ptr, transform []PipelineOp) (Ptr, error) {
	n := len(transform)
	if n == 0 {
		return p, nil
	}
	s := p.Struct()
	for _, op := range transform[:n-1] {
		field, err := s.Ptr(op.Field)
		if err != nil {
			return Ptr{}, err
		}
		s, err = field.StructDefault(op.DefaultValue)
		if err != nil {
			return Ptr{}, err
		}
	}
	op := transform[n-1]
	p, err := s.Ptr(op.Field)
	if err != nil {
		return Ptr{}, err
	}
	if op.DefaultValue != nil {
		p, err = p.Default(op.DefaultValue)
	}
	return p, err
}

// clientFromResolution retrieves a client from a resolved answer by
// applying a transform.
func clientFromResolution(obj Ptr, err error, transform []PipelineOp) *Client {
	if err != nil {
		return ErrorClient(err)
	}
	out, err := Transform(obj, transform)
	if err != nil {
		return ErrorClient(err)
	}
	return out.Interface().Client()
}
