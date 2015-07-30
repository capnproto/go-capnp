package capnp

import (
	"errors"
	"sync"
)

// callQueueSize is the maximum number of calls
const callQueueSize = 64

// Fulfiller is a promise for a Struct.  The zero value is an unresolved
// answer.  A Fulfiller is considered to be resolved once Fulfill or
// Reject is called.  Calls to the Fulfiller will queue up until it is
// resolved.  A Fulfiller is safe to use from multiple goroutines.
type Fulfiller struct {
	once     sync.Once
	resolved chan struct{} // initialized by init()

	// Protected by mu
	mu     sync.RWMutex
	answer Answer
	queue  []pcall // initialized by init()
}

// init initializes the Fulfiller.  It is idempotent.
// Should be called for each method on Fulfiller.
func (f *Fulfiller) init() {
	f.once.Do(func() {
		f.resolved = make(chan struct{})
		f.queue = make([]pcall, 0, callQueueSize)
	})
}

// Fulfill sets the fulfiller's answer to s.  If there are queued
// pipeline calls, the capabilities on the struct will be embargoed
// until the queued calls finish.  Fulfill will panic if the fulfiller
// has already been resolved.
func (f *Fulfiller) Fulfill(s Struct) {
	f.init()
	f.mu.Lock()
	if f.answer != nil {
		f.mu.Unlock()
		panic("Fulfiller.Fulfill called more than once")
	}
	f.answer = ImmediateAnswer(s)
	queues := f.emptyQueue(s)
	ctab := s.Segment.Message.CapTable()
	for capIdx, q := range queues {
		ctab[capIdx] = newEmbargoClient(ctab[capIdx], q)
	}
	close(f.resolved)
	f.mu.Unlock()
}

// emptyQueue splits the queue by which capability it targets and
// drops any invalid calls.  Once this function returns, f.queue will
// be nil.
func (f *Fulfiller) emptyQueue(s Struct) map[uint32][]ecall {
	qs := make(map[uint32][]ecall, len(f.queue))
	for i, pc := range f.queue {
		c := TransformObject(Object(s), pc.transform)
		if c.Type() != TypeInterface {
			pc.f.Reject(ErrNullClient)
			continue
		}
		cn := c.ToInterface().Capability()
		if qs[cn] == nil {
			qs[cn] = make([]ecall, 0, len(f.queue)-i)
		}
		qs[cn] = append(qs[cn], ecall{
			call: pc.call,
			f:    pc.f,
		})
	}
	f.queue = nil
	return qs
}

// Reject sets the fulfiller's answer to err.  If there are queued
// pipeline calls, they will all return errors.  Reject will panic if
// the error is nil or the fulfiller has already been resolved.
func (f *Fulfiller) Reject(err error) {
	if err == nil {
		panic("Fulfiller.Reject called with nil")
	}
	f.init()
	f.mu.Lock()
	if f.answer != nil {
		f.mu.Unlock()
		panic("Fulfiller.Reject called more than once")
	}
	f.answer = ErrorAnswer(err)
	for i := range f.queue {
		f.queue[i].f.Reject(err)
		f.queue[i] = pcall{}
	}
	close(f.resolved)
	f.mu.Unlock()
}

// Done returns a channel that is closed once f is resolved.
func (f *Fulfiller) Done() <-chan struct{} {
	f.init()
	return f.resolved
}

// Peek returns f's resolved answer or nil if f has not been resolved.
func (f *Fulfiller) Peek() Answer {
	f.init()
	f.mu.RLock()
	a := f.answer
	f.mu.RUnlock()
	return a
}

// Struct waits until f is resolved and returns its struct if fulfilled
// or an error if rejected.
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
	// Make sure that f wasn't fulfilled.
	if a := f.answer; a != nil {
		f.mu.Unlock()
		return a.PipelineCall(transform, call)
	}
	if len(f.queue) == cap(f.queue) {
		f.mu.Unlock()
		return ErrorAnswer(errCallQueueFull)
	}
	g := new(Fulfiller)
	f.queue = append(f.queue, pcall{
		transform: transform,
		call: &Call{
			Ctx:     call.Ctx,
			Method:  call.Method,
			Params:  call.PlaceParams(nil),
			Options: call.Options,
		},
		f: new(Fulfiller),
	})
	f.mu.Unlock()
	return g
}

// PipelineClose waits until f is resolved and then calls PipelineClose
// on the fulfilled answer.
func (f *Fulfiller) PipelineClose(transform []PipelineOp) error {
	<-f.Done()
	return f.Peek().PipelineClose(transform)
}

// pcall is a queued pipeline call.
type pcall struct {
	transform []PipelineOp
	call      *Call
	f         *Fulfiller
}

// embargoClient is a client that flushes a queue of calls.
type embargoClient struct {
	client Client

	mu       sync.RWMutex
	queue    []ecall
	start, n int
}

func newEmbargoClient(client Client, queue []ecall) Client {
	ec := &embargoClient{client: client, queue: make([]ecall, callQueueSize)}
	ec.n = copy(ec.queue, queue)
	go ec.flushQueue()
	return ec
}

func (ec *embargoClient) push(cl *Call) Answer {
	if ec.n == len(ec.queue) {
		return ErrorAnswer(errCallQueueFull)
	}
	f := new(Fulfiller)
	i := (ec.start + ec.n) % len(ec.queue)
	ec.queue[i] = ecall{cl, f}
	return f
}

func (ec *embargoClient) pop() ecall {
	if ec.n == 0 {
		return ecall{}
	}
	c := ec.queue[ec.start]
	ec.queue[ec.start] = ecall{}
	ec.start = (ec.start + 1) % len(ec.queue)
	ec.n--
	ec.mu.Unlock()
	return c
}

// flushQueue is run in its own goroutine.
func (ec *embargoClient) flushQueue() {
	for {
		ec.mu.Lock()
		c := ec.pop()
		ec.mu.Unlock()
		if c.call == nil {
			return
		}
		ans := ec.client.Call(c.call)
		go func(f *Fulfiller, ans Answer) {
			s, err := ans.Struct()
			if err == nil {
				f.Fulfill(s)
			} else {
				f.Reject(err)
			}
		}(c.f, ans)
	}
}

func (ec *embargoClient) Call(cl *Call) Answer {
	// Fast path: queue is flushed.
	ec.mu.RLock()
	n := ec.n
	ec.mu.RUnlock()
	if n == 0 {
		return ec.client.Call(cl)
	}

	// Add to queue.
	ec.mu.Lock()
	// Since we released the lock, check that the queue hasn't been flushed.
	if ec.n == 0 {
		ec.mu.Unlock()
		return ec.client.Call(cl)
	}
	ans := ec.push(cl)
	ec.mu.Unlock()
	return ans
}

func (ec *embargoClient) Close() error {
	ec.mu.Lock()
	// reject all queued calls
	for {
		c := ec.pop()
		if c.call == nil {
			break
		}
		c.f.Reject(errQueueCallCancel)
	}
	ec.mu.Unlock()
	return ec.client.Close()
}

// ecall is an queued embargoed call.
type ecall struct {
	call *Call
	f    *Fulfiller
}

var (
	errCallQueueFull   = errors.New("capnp: promised answer call queue full")
	errQueueCallCancel = errors.New("capnp: queued call canceled")
)
