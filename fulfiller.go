package capnp

import (
	"errors"
	"sync"
)

// Fulfiller is a placeholder for a promised answer.  The zero value is
// an empty placeholder.  Calls to the fulfiller will queue up until it
// is fulfilled.  A Fulfiller is safe to use from multiple goroutines.
type Fulfiller struct {
	// Initialized by init()
	once      sync.Once
	fulfilled chan struct{}

	// Protected by mu
	mu        sync.RWMutex
	answer    Answer
	dequeuing bool
	queue     chan pcall // TODO(light): use ring buffer
}

// init initializes the Fulfiller.
// Should be called for each method on Fulfiller.
func (f *Fulfiller) init() {
	f.once.Do(func() {
		f.fulfilled = make(chan struct{})
		f.queue = make(chan pcall, 64)
	})
}

// Fulfill sets the fulfiller's answer to a.  If there are outstanding
// pipeline calls, a goroutine is started to flush the queue and once
// all the calls have been sent, the fulfiller will be marked fulfilled.
// Fulfill will panic if a is nil or the fulfiller has already been
// fulfilled.
func (f *Fulfiller) Fulfill(a Answer) {
	if a == nil {
		panic("Fulfiller.Fulfill called with nil")
	}
	f.init()
	f.mu.Lock()
	defer f.mu.Unlock()

	// TODO(light): use embargoed clients
	if f.dequeuing || f.answer != nil {
		panic("Fulfiller.Fulfill called more than once")
	}
	if len(f.queue) == 0 {
		f.answer = a
		close(f.fulfilled)
		return
	}
	f.dequeuing = true
	go func() {
		for {
			f.mu.Lock()
			if len(f.queue) == 0 {
				f.answer = a
				close(f.fulfilled)
				f.mu.Unlock()
				return
			}
			pc := <-f.queue
			f.mu.Unlock()
			pc.f.Fulfill(a.PipelineCall(pc.transform, pc.call))
		}
	}()
}

// Done returns a channel that is closed once f is fulfilled.
func (f *Fulfiller) Done() <-chan struct{} {
	f.init()
	return f.fulfilled
}

// Peek returns f's fulfilled answer or nil if Fulfill has not been
// called yet.
func (f *Fulfiller) Peek() Answer {
	f.init()
	f.mu.RLock()
	a := f.answer
	f.mu.RUnlock()
	return a
}

// Struct waits until f is fulfilled and then calls Struct on the
// fulfilled answer.
func (f *Fulfiller) Struct() (Struct, error) {
	<-f.Done()
	return f.Peek().Struct()
}

// PipelineCall calls PipelineCall on the fulfilled answer or queues the
// call if f has not been fulfilled.
func (f *Fulfiller) PipelineCall(transform []PipelineOp, call *Call) Answer {
	f.init()

	// Fast path: pass-through after fulfilled.
	if a := f.Peek(); a != nil {
		return a.PipelineCall(transform, call)
	}

	f.mu.Lock()
	defer f.mu.Unlock()
	// Make sure that f wasn't fulfilled.
	if a := f.answer; a != nil {
		return a.PipelineCall(transform, call)
	}
	pc := pcall{
		transform: transform,
		call: &Call{
			Ctx:     call.Ctx,
			Method:  call.Method,
			Params:  call.PlaceParams(nil),
			Options: call.Options,
		},
		f: new(Fulfiller),
	}
	select {
	case f.queue <- pc:
		return pc.f
	default:
		return ErrorAnswer(errFulfillerQueueFull)
	}
}

// PipelineClose waits until f is fulfilled and then calls PipelineClose
// on the fulfilled answer.
func (f *Fulfiller) PipelineClose(transform []PipelineOp) error {
	<-f.Done()
	return f.Peek().PipelineClose(transform)
}

type pcall struct {
	transform []PipelineOp
	call      *Call
	f         *Fulfiller
}

var errFulfillerQueueFull = errors.New("capnp: fulfiller pipeline call: queue full")
