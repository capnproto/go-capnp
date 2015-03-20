package rpc

import (
	"sync"

	"zombiezen.com/go/capnproto"
	"zombiezen.com/go/capnproto/rpc/internal/rpc"
)

// fulfiller is a placeholder for a promised value or an error.
// The zero value is an empty placeholder.  A fulfiller is safe
// to use from multiple goroutines.
type fulfiller struct {
	once     sync.Once
	resolved chan struct{} // closed after fulfill/reject

	// All of these fields are protected by mu.
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

func (f *fulfiller) resolve(a capnp.Answer) {
	f.init()
	f.mu.Lock()
	if f.answer == nil {
		f.answer = a
		close(f.resolved)
	}
	f.mu.Unlock()
}

func (f *fulfiller) Struct() (capnp.Struct, error) {
	f.init()
	<-f.resolved
	f.mu.RLock()
	a := f.answer
	f.mu.RUnlock()
	return a.Struct()
}

func (f *fulfiller) PipelineCall(transform []capnp.PipelineOp, call *capnp.Call) capnp.Answer {
	f.init()
	f.mu.RLock()
	a := f.answer
	f.mu.RUnlock()
	if a != nil {
		return a.PipelineCall(transform, call)
	}

	g := new(fulfiller)
	go func() {
		<-f.resolved
		f.mu.RLock()
		a := f.answer
		f.mu.RUnlock()
		g.resolve(a.PipelineCall(transform, call))
	}()
	return g
}

func (f *fulfiller) PipelineClose(transform []capnp.PipelineOp) error {
	<-f.resolved
	f.mu.RLock()
	a := f.answer
	f.mu.RUnlock()
	return a.PipelineClose(transform)
}

func transformToPromisedAnswer(s *capnp.Segment, answer rpc.PromisedAnswer, transform []capnp.PipelineOp) {
	opList := rpc.NewPromisedAnswer_Op_List(s, len(transform))
	for i, op := range transform {
		opList.At(i).SetGetPointerField(uint16(op.Field))
	}
	answer.SetTransform(opList)
}
