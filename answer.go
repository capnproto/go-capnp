package capnp

import (
	"context"
	"strconv"
	"sync"
)

// A Promise holds the result of an RPC call.  Only one of Fulfill,
// Reject, or Join can be called on a Promise.  Before the result is
// written, calls can be queued up using the Answer methods — this is
// promise pipelining.
//
// Promise is most useful for implementing ClientHook.
// Most applications will use Answer, since that what is returned by
// a Client.
type Promise struct {
	resolved chan struct{} // can only be closed while holding mu
	ans      Answer

	mu       sync.RWMutex
	caller   PipelineCaller // nil if resolved
	children []*Promise
	joined   bool
	result   Ptr
	err      error
}

// NewPromise creates a new unresolved promise.  The PipelineCaller will
// be used to make pipelined calls before the promise resolves.
func NewPromise(pc PipelineCaller) *Promise {
	p := &Promise{
		caller:   pc,
		resolved: make(chan struct{}),
	}
	p.ans.f.promise = p
	return p
}

// Fulfill resolves the promise with a successful result.
func (p *Promise) Fulfill(result Ptr) {
	defer p.mu.Unlock()
	p.mu.Lock()
	if p.joined {
		panic("Promise.Fulfill called after Promise.Join")
	}
	select {
	case <-p.resolved:
		panic("Promise.Fulfill or Promise.Reject called more than once")
	default:
		p.lockedResolve(result, nil)
	}
}

// Reject resolves the promise with a failure.
func (p *Promise) Reject(e error) {
	if e == nil {
		panic("Promise.Reject(nil)")
	}
	defer p.mu.Unlock()
	p.mu.Lock()
	if p.joined {
		panic("Promise.Reject called after Promise.Join")
	}
	select {
	case <-p.resolved:
		panic("Promise.Fulfill or Promise.Reject called more than once")
	default:
		p.lockedResolve(Ptr{}, e)
	}
}

// Join ties the outcome of a promise to another promise's outcome.
func (p *Promise) Join(from *Answer) {
	parent := from.f.promise
	parent.mu.Lock()
	select {
	case <-parent.resolved:
		result, err := parent.result, parent.err
		parent.mu.Unlock()

		p.mu.Lock()
		if p.joined {
			panic("double Promise.Join")
		}
		select {
		case <-p.resolved:
			panic("Promise.Join after resolved")
		default:
		}
		p.lockedResolve(result, err)
		p.mu.Unlock()
		return
	default:
	}
	p.mu.Lock()
	if p.joined {
		panic("double Promise.Join")
	}
	select {
	case <-p.resolved:
		panic("Promise.Join after resolved")
	default:
	}
	p.caller = parent.caller
	children := p.children
	p.children = nil
	p.joined = true
	p.mu.Unlock()

	parent.children = append(parent.children, p)
	parent.children = append(parent.children, children...)
	for _, c := range children {
		c.mu.Lock()
		c.caller = parent.caller
		c.mu.Unlock()
	}
	parent.mu.Unlock()
}

func (p *Promise) lockedResolve(r Ptr, e error) {
	p.result, p.err = r, e
	p.caller = nil
	for _, c := range p.children {
		c.mu.Lock()
		c.result, c.err = r, e
		c.caller = nil
		close(c.resolved)
		c.mu.Unlock()
	}
	p.children = nil
	close(p.resolved)
}

// Answer returns a read-only view of the promise.
func (p *Promise) Answer() *Answer {
	return &p.ans
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
	PipelineSend(ctx context.Context, transform []PipelineOp, s Send) (*Answer, ReleaseFunc)
	PipelineRecv(ctx context.Context, transform []PipelineOp, r Recv) (*Answer, ReleaseFunc)
}

// An Answer is a deferred result of a client call.  Conceptually, this is a
// future.  It is safe to use from multiple goroutines.
type Answer struct {
	f Future
}

// ErrorAnswer returns a Answer that always returns error e.
func ErrorAnswer(e error) *Answer {
	p := &Promise{
		resolved: newClosedSignal(),
		err:      e,
	}
	p.ans.f.promise = p
	return &p.ans
}

// ImmediateAnswer returns an Answer that accesses s.
func ImmediateAnswer(s Struct) *Answer {
	p := &Promise{
		resolved: newClosedSignal(),
		result:   s.ToPtr(),
	}
	p.ans.f.promise = p
	return &p.ans
}

// Future returns a future that is equivalent to ans.
func (ans *Answer) Future() *Future {
	return &ans.f
}

// Done returns a channel that is closed when the answer's call is finished.
func (ans *Answer) Done() <-chan struct{} {
	return ans.f.Done()
}

// Struct waits until the answer is resolved and returns the struct
// this answer represents.
func (ans *Answer) Struct() (Struct, error) {
	return ans.f.Struct()
}

// Client returns the answer as a client.  If the answer's originating
// call has not completed, then calls will be queued until the original
// call's completion.
func (ans *Answer) Client() *Client {
	return ans.f.Client()
}

// Field returns a derived future which yields the pointer field given,
// defaulting to the value given.
func (ans *Answer) Field(off uint16, def []byte) *Future {
	return ans.f.Field(off, def)
}

// SendCall starts a pipelined call.
func (ans *Answer) SendCall(ctx context.Context, transform []PipelineOp, s Send) (*Answer, ReleaseFunc) {
	return ans.f.promise.client(transform).SendCall(ctx, s)
}

// RecvCall starts a pipelined call.
func (ans *Answer) RecvCall(ctx context.Context, transform []PipelineOp, r Recv) (*Answer, ReleaseFunc) {
	return ans.f.promise.client(transform).RecvCall(ctx, r)
}

// A Future accesses a portion of an Answer.  It is safe to use from
// multiple goroutines.
type Future struct {
	promise *Promise
	parent  *Future // nil if root future
	op      PipelineOp
}

// transform returns the operations needed to transform the root answer
// into the value f represents.
func (f *Future) transform() []PipelineOp {
	if f.parent == nil {
		return nil
	}
	n := 0
	for g := f; g.parent != nil; g = g.parent {
		n++
	}
	xform := make([]PipelineOp, n)
	for i, g := n-1, f; g.parent != nil; i, g = i-1, g.parent {
		xform[i] = g.op
	}
	return xform
}

// Done returns a channel that is closed when the answer's call is finished.
func (f *Future) Done() <-chan struct{} {
	return f.promise.resolved
}

// Struct waits until the answer is resolved and returns the struct
// this future represents.
func (f *Future) Struct() (Struct, error) {
	<-f.promise.resolved
	f.promise.mu.RLock()
	ptr, err := f.promise.result, f.promise.err
	f.promise.mu.RUnlock()
	if err != nil {
		return Struct{}, err
	}
	ptr, err = Transform(ptr, f.transform())
	if err != nil {
		return Struct{}, err
	}
	return ptr.Struct(), nil
}

// Client returns the future as a client.  If the answer's originating
// call has not completed, then calls will be queued until the original
// call's completion.
func (f *Future) Client() *Client {
	return f.promise.client(f.transform())
}

// Field returns a derived future which yields the pointer field given,
// defaulting to the value given.
func (f *Future) Field(off uint16, def []byte) *Future {
	return &Future{
		promise: f.promise,
		parent:  f,
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

func (pc pipelineClient) Send(ctx context.Context, s Send) (*Answer, ReleaseFunc) {
	defer pc.p.mu.RUnlock()
	pc.p.mu.RLock()
	if pc.p.caller != nil {
		return pc.p.caller.PipelineSend(ctx, pc.transform, s)
	}
	c := clientFromResolution(pc.p.result, pc.p.err, pc.transform)
	return c.SendCall(ctx, s)
}

func (pc pipelineClient) Recv(ctx context.Context, r Recv) (*Answer, ReleaseFunc) {
	defer pc.p.mu.RUnlock()
	pc.p.mu.RLock()
	if pc.p.caller != nil {
		return pc.p.caller.PipelineRecv(ctx, pc.transform, r)
	}
	c := clientFromResolution(pc.p.result, pc.p.err, pc.transform)
	return c.RecvCall(ctx, r)
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
