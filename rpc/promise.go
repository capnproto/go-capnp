package rpc

import (
	"sync"

	"zombiezen.com/go/capnproto"
	"zombiezen.com/go/capnproto/rpc/rpccapnp"
)

// fulfiller is a placeholder for a promised answer.
// The zero value is an empty placeholder.  A fulfiller is safe
// to use from multiple goroutines.
type fulfiller struct {
	once     sync.Once
	resolved chan struct{} // closed after resolve

	// Protected by mu
	mu     sync.RWMutex
	answer capnp.Answer
}

// init initializes the fulfiller.
// It is idempotent and is not necessary to call.
func (f *fulfiller) init() {
	f.once.Do(func() {
		f.resolved = make(chan struct{})
	})
}

// resolve sets the fulfiller's answer to a, which must not be nil.
// resolve should be called only once for a fulfiller.
func (f *fulfiller) resolve(a capnp.Answer) {
	f.init()
	f.mu.Lock()
	if f.answer == nil {
		f.answer = a
		close(f.resolved)
	}
	f.mu.Unlock()
}

// peek returns the answer held by the fulfiller or nil if resolve has
// not been called.  It does not block for resolution.
func (f *fulfiller) peek() capnp.Answer {
	f.mu.RLock()
	a := f.answer
	f.mu.RUnlock()
	return a
}

func (f *fulfiller) Struct() (capnp.Struct, error) {
	f.init()
	<-f.resolved
	return f.peek().Struct()
}

func (f *fulfiller) PipelineCall(transform []capnp.PipelineOp, call *capnp.Call) capnp.Answer {
	f.init()
	a := f.peek()
	if a != nil {
		return a.PipelineCall(transform, call)
	}

	g := new(fulfiller)
	go func() {
		<-f.resolved
		a := f.peek()
		g.resolve(a.PipelineCall(transform, call))
	}()
	return g
}

func (f *fulfiller) PipelineClose(transform []capnp.PipelineOp) error {
	<-f.resolved
	a := f.peek()
	return a.PipelineClose(transform)
}

func transformToPromisedAnswer(s *capnp.Segment, answer rpccapnp.PromisedAnswer, transform []capnp.PipelineOp) {
	opList := rpccapnp.NewPromisedAnswer_Op_List(s, len(transform))
	for i, op := range transform {
		opList.At(i).SetGetPointerField(uint16(op.Field))
	}
	answer.SetTransform(opList)
}
