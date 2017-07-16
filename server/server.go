// Package server provides runtime support for implementing Cap'n Proto
// interfaces locally.
package server // import "zombiezen.com/go/capnproto2/server"

import (
	"context"
	"errors"
	"sort"
	"sync"

	"zombiezen.com/go/capnproto2"
)

// A Method describes a single method on a server object.
type Method struct {
	capnp.Method
	Impl        Func
	ResultsSize capnp.ObjectSize
}

// A Func is a function that implements a single method.
type Func func(ctx context.Context, params, results capnp.Struct, opts capnp.CallOptions) error

// Closer is the interface that wraps the Close method.
type Closer interface {
	Close() error
}

// A Server is a locally implemented interface.  It implements the
// capnp.ClientHookCloser interface.
type Server struct {
	methods sortedMethods
	brand   interface{}
	closer  Closer

	calls sync.WaitGroup
	mu    sync.Mutex // serializes calls
}

// Policy is a set of behavioral parameters for a Server.
// They're not specific to a particular server and are generally set at
// an application level.  Library functions are encouraged to accept a
// Policy from a caller instead of creating their own.
type Policy struct {
	// TODO(soon): Add queue size etc.
}

// New returns a client that makes calls to a set of methods.
// If closer is nil then the client's Close is a no-op.  The server
// guarantees message delivery order by blocking each call on the
// return or acknowledgment of the previous call.  See the Ack function
// for more details.
func New(methods []Method, brand interface{}, closer Closer, policy *Policy) *Server {
	s := &Server{
		methods: make(sortedMethods, len(methods)),
		brand:   brand,
		closer:  closer,
	}
	copy(s.methods, methods)
	sort.Sort(s.methods)
	return s
}

// Send starts a method call.
func (s *Server) Send(ctx context.Context, m capnp.Method, a capnp.SendArgs, opts capnp.CallOptions) (*capnp.Answer, capnp.ReleaseFunc) {
	mm := s.methods.find(m)
	if mm == nil {
		// TODO(soon): signal unimplemented.
		return capnp.ErrorAnswer(errors.New("unimplemented")), func() {}
	}
	args, err := argsToStruct(a)
	if err != nil {
		return capnp.ErrorAnswer(err), func() {}
	}
	rargs := capnp.RecvArgs{
		Args: args,
		Release: func() {
			// TODO(someday): log error from ClearCaps
			if seg := args.Segment(); seg != nil {
				seg.Message().Reset(nil)
			}
		},
	}
	return s.start(ctx, mm, rargs, opts)
}

// Recv starts a method call.
func (s *Server) Recv(ctx context.Context, m capnp.Method, a capnp.RecvArgs, opts capnp.CallOptions) (*capnp.Answer, capnp.ReleaseFunc) {
	mm := s.methods.find(m)
	if mm == nil {
		// TODO(soon): signal unimplemented.
		return capnp.ErrorAnswer(errors.New("unimplemented")), func() {}
	}
	return s.start(ctx, mm, a, opts)
}

func (s *Server) start(ctx context.Context, m *Method, a capnp.RecvArgs, opts capnp.CallOptions) (*capnp.Answer, capnp.ReleaseFunc) {
	// TODO(someday): Throttle number of concurrent calls.
	defer s.mu.Unlock()
	s.calls.Add(1)
	s.mu.Lock()
	results, err := newBlankStruct(m.ResultsSize)
	if err != nil {
		a.Release()
		return capnp.ErrorAnswer(err), func() {}
	}
	acksig := newAckSignal()
	opts = opts.With([]capnp.CallOption{capnp.SetOptionValue(ackSignalKey{}, acksig)})
	p := capnp.NewPromise(new(pipelineQueue))
	go func() {
		defer s.calls.Done()
		err := m.Impl(ctx, a.Args, results, opts)
		a.Release()
		if err != nil {
			p.Reject(err)
			// TODO(someday): log error from ClearCaps
			results.Segment().Message().Reset(nil)
			return
		}
		p.Fulfill(results.ToPtr())
	}()
	ans := p.Answer()
	select {
	case <-acksig.c:
	case <-ans.Done():
		// Implementation functions may not call Ack, which is fine for
		// smaller functions.
	case <-ctx.Done():
		// Ideally, this would reject the answer immediately, but that
		// would race with the goroutine.
	}
	return ans, func() {
		<-ans.Done()
		// TODO(someday): log error from ClearCaps
		results.Segment().Message().Reset(nil)
	}
}

// Brand returns a value that will match IsServer.
func (s *Server) Brand() interface{} {
	return serverBrand{s.brand}
}

// Close waits for ongoing calls to finish and calls Close to the Closer.
func (s *Server) Close() error {
	// TODO(someday): cancel all outstanding calls.
	s.calls.Wait()
	if s.closer == nil {
		return nil
	}
	return s.closer.Close()
}

// IsServer reports whether a brand returned by capnp.Client.Brand
// originated from Server.Brand, and returns the brand argument passed
// to New.
func IsServer(brand interface{}) (_ interface{}, ok bool) {
	sb, ok := brand.(serverBrand)
	return sb.x, ok
}

type serverBrand struct {
	x interface{}
}

// Ack acknowledges delivery of a server call, allowing other methods
// to be called on the server.  It is intended to be used inside the
// implementation of a server function.  Calling Ack on options that
// aren't from a server method implementation is a no-op.
//
// Example:
//
//	func (my *myServer) MyMethod(call schema.MyServer_myMethod) error {
//		server.Ack(call.Options)
//		// ... do long-running operation ...
//		return nil
//	}
//
// Ack need not be the first call in a function nor is it required.
// Since the function's return is also an acknowledgment of delivery,
// short functions can return without calling Ack.  However, since
// clients will not return an Answer until the delivery is acknowledged,
// it is advisable to ack early.
func Ack(opts capnp.CallOptions) {
	if ack, _ := opts.Value(ackSignalKey{}).(*ackSignal); ack != nil {
		ack.signal()
	}
}

type pipelineQueue struct {
}

func (pq *pipelineQueue) PipelineRecv(ctx context.Context, transform []capnp.PipelineOp, m capnp.Method, a capnp.RecvArgs, opts capnp.CallOptions) (*capnp.Answer, capnp.ReleaseFunc) {
	a.Release()
	return capnp.ErrorAnswer(errors.New("TODO(soon)")), func() {}
}

func (pq *pipelineQueue) PipelineSend(ctx context.Context, transform []capnp.PipelineOp, m capnp.Method, a capnp.SendArgs, opts capnp.CallOptions) (*capnp.Answer, capnp.ReleaseFunc) {
	return capnp.ErrorAnswer(errors.New("TODO(soon)")), func() {}
}

func argsToStruct(a capnp.SendArgs) (capnp.Struct, error) {
	if a.Place == nil {
		return capnp.Struct{}, nil
	}
	s, err := newBlankStruct(a.Size)
	if err != nil {
		return capnp.Struct{}, err
	}
	if err := a.Place(s); err != nil {
		s.Segment().Message().Reset(nil)
		return capnp.Struct{}, err
	}
	return s, nil
}

func newBlankStruct(sz capnp.ObjectSize) (capnp.Struct, error) {
	_, seg, err := capnp.NewMessage(capnp.MultiSegment(nil))
	if err != nil {
		return capnp.Struct{}, err
	}
	s, err := capnp.NewRootStruct(seg, sz)
	if err != nil {
		return capnp.Struct{}, err
	}
	return s, nil
}

type sortedMethods []Method

// find returns the method with the given ID or nil.
func (sm sortedMethods) find(id capnp.Method) *Method {
	i := sort.Search(len(sm), func(i int) bool {
		m := &sm[i]
		if m.InterfaceID != id.InterfaceID {
			return m.InterfaceID >= id.InterfaceID
		}
		return m.MethodID >= id.MethodID
	})
	if i == len(sm) {
		return nil
	}
	m := &sm[i]
	if m.InterfaceID != id.InterfaceID || m.MethodID != id.MethodID {
		return nil
	}
	return m
}

func (sm sortedMethods) Len() int {
	return len(sm)
}

func (sm sortedMethods) Less(i, j int) bool {
	if id1, id2 := sm[i].InterfaceID, sm[j].InterfaceID; id1 != id2 {
		return id1 < id2
	}
	return sm[i].MethodID < sm[j].MethodID
}

func (sm sortedMethods) Swap(i, j int) {
	sm[i], sm[j] = sm[j], sm[i]
}

type ackSignal struct {
	c    chan struct{}
	once sync.Once
}

func newAckSignal() *ackSignal {
	return &ackSignal{c: make(chan struct{})}
}

func (ack *ackSignal) signal() {
	ack.once.Do(func() {
		close(ack.c)
	})
}

type ackSignalKey struct{}
