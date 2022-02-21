// Package server provides runtime support for implementing Cap'n Proto
// interfaces locally.
package server // import "capnproto.org/go/capnp/v3/server"

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/exc"
)

// A Method describes a single capability method on a server object.
type Method struct {
	capnp.Method
	Impl func(context.Context, *Call) error
}

// Call holds the state of an ongoing capability method call.
// A Call cannot be used after the server method returns.
type Call struct {
	args capnp.Struct

	alloced bool
	alloc   resultsAllocer
	results capnp.Struct

	ack   chan<- struct{}
	acked bool
}

func newCall(args capnp.Struct, ra resultsAllocer) (*Call, <-chan struct{}) {
	ack := make(chan struct{})
	return &Call{
		args:  args,
		alloc: ra,
		ack:   ack,
	}, ack
}

// Args returns the call's arguments.  Args is not safe to
// reference after a method implementation returns.  Args is safe to
// call and read from multiple goroutines.
func (c *Call) Args() capnp.Struct {
	return c.args
}

// AllocResults allocates the results struct.  It is an error to call
// AllocResults more than once.
func (c *Call) AllocResults(sz capnp.ObjectSize) (capnp.Struct, error) {
	if c.alloced {
		return capnp.Struct{}, newError("multiple calls to AllocResults")
	}
	var err error
	c.alloced = true
	c.results, err = c.alloc.AllocResults(sz)
	return c.results, err
}

// Ack is a function that is called to acknowledge the delivery of the
// RPC call, allowing other RPC methods to be called on the server.
// After the first call, subsequent calls to Ack do nothing.
//
// Ack need not be the first call in a function nor is it required.
// Since the function's return is also an acknowledgment of delivery,
// short functions can return without calling Ack.  However, since
// the server will not return an Answer until the delivery is
// acknowledged, failure to acknowledge a call before waiting on an
// RPC may cause deadlocks.
func (c *Call) Ack() {
	if c.acked {
		return
	}
	close(c.ack)
	c.acked = true
}

// Shutdowner is the interface that wraps the Shutdown method.
type Shutdowner interface {
	Shutdown()
}

// A Server is a locally implemented interface.  It implements the
// capnp.ClientHook interface.
type Server struct {
	methods  sortedMethods
	brand    interface{}
	shutdown Shutdowner
	policy   Policy

	// mu protects the following fields.
	// mu should never be held while calling application code.
	mu sync.Mutex

	// ongoing is a fixed-size list of ongoing calls.
	// It is used as a semaphore: when all elements are set, no new work
	// can be started until an element is cleared.
	ongoing []cstate

	// starting is non-nil if start() is waiting for acknowledgement of a
	// call.  It is closed when the acknowledgement is received.
	starting <-chan struct{}

	// full is non-nil if a start() is waiting for a space in ongoing to
	// free up.  It is closed and set to nil when the next call returns.
	full chan<- struct{}

	// drain is non-nil when Shutdown starts and is closed by the last
	// call to return.
	drain chan struct{}
}

type cstate struct {
	cancel context.CancelFunc // nil if slot free
}

// Policy is a set of behavioral parameters for a Server.
// They're not specific to a particular server and are generally set at
// an application level.  Library functions are encouraged to accept a
// Policy from a caller instead of creating their own.
type Policy struct {
	// MaxConcurrentCalls is the maximum number of methods allowed to be
	// executing on a single Server simultaneously.  Attempts to make more
	// calls than this limit will result in immediate error answers.
	//
	// If this is zero, then a reasonably small default is used.
	MaxConcurrentCalls int
}

// New returns a client hook that makes calls to a set of methods.
// If shutdown is nil then the server's shutdown is a no-op.  The server
// guarantees message delivery order by blocking each call on the
// return or acknowledgment of the previous call.  See Call.Ack for more
// details.
func New(methods []Method, brand interface{}, shutdown Shutdowner, policy *Policy) *Server {
	srv := &Server{
		methods:  make(sortedMethods, len(methods)),
		brand:    brand,
		shutdown: shutdown,
	}
	copy(srv.methods, methods)
	sort.Sort(srv.methods)
	if policy != nil {
		srv.policy = *policy
	}
	if srv.policy.MaxConcurrentCalls < 1 {
		srv.policy.MaxConcurrentCalls = 2
	}
	srv.ongoing = make([]cstate, srv.policy.MaxConcurrentCalls)
	return srv
}

// Send starts a method call.
func (srv *Server) Send(ctx context.Context, s capnp.Send) (*capnp.Answer, capnp.ReleaseFunc) {
	mm := srv.methods.find(s.Method)
	if mm == nil {
		return capnp.ErrorAnswer(s.Method, capnp.Unimplemented("unimplemented")), func() {}
	}
	args, err := sendArgsToStruct(s)
	if err != nil {
		return capnp.ErrorAnswer(mm.Method, err), func() {}
	}
	ret := new(structReturner)
	return ret.answer(mm.Method, srv.start(ctx, mm, capnp.Recv{
		Method: mm.Method, // pick up names from server method
		Args:   args,
		ReleaseArgs: func() {
			if msg := args.Message(); msg != nil {
				msg.Reset(nil)
				args = capnp.Struct{}
			}
		},
		Returner: ret,
	}))
}

// Recv starts a method call.
func (srv *Server) Recv(ctx context.Context, r capnp.Recv) capnp.PipelineCaller {
	mm := srv.methods.find(r.Method)
	if mm == nil {
		r.Reject(capnp.Unimplemented("unimplemented"))
		return nil
	}
	return srv.start(ctx, mm, r)
}

func (srv *Server) start(ctx context.Context, m *Method, r capnp.Recv) capnp.PipelineCaller {
	// Acquire "starting" condition variable.
	srv.mu.Lock()
	for {
		if srv.drain != nil {
			srv.mu.Unlock()
			r.Reject(exc.New(exc.Failed, "capnp server", "call after shutdown"))
			return nil
		}
		if srv.starting == nil {
			break
		}
		wait := srv.starting
		srv.mu.Unlock()
		select {
		case <-wait:
		case <-ctx.Done():
			r.Reject(ctx.Err())
			return nil
		}
		srv.mu.Lock()
	}
	starting := make(chan struct{})
	srv.starting = starting

	// Acquire an ID (semaphore).
	id := srv.nextID()
	if id == -1 {
		full := make(chan struct{})
		srv.full = full
		srv.mu.Unlock()
		select {
		case <-full:
		case <-ctx.Done():
			srv.mu.Lock()
			srv.starting = nil
			close(starting)
			srv.full = nil // full could be nil or non-nil, ensure it is nil.
			srv.mu.Unlock()
			r.Reject(ctx.Err())
			return nil
		}
		srv.mu.Lock()
		id = srv.nextID()
		if srv.drain != nil {
			srv.starting = nil
			close(starting)
			srv.mu.Unlock()
			r.Reject(exc.New(exc.Failed, "capnp server", "call after shutdown"))
			return nil
		}
	}

	// Bookkeeping: set starting to indicate we're waiting for an ack and
	// record the cancel function for draining.
	ctx, cancel := context.WithCancel(ctx)
	srv.ongoing[id] = cstate{cancel}
	srv.mu.Unlock()

	// Call implementation function.
	call, ack := newCall(r.Args, r.Returner)
	aq := newAnswerQueue(r.Method)
	done := make(chan struct{})
	go func() {
		err := m.Impl(ctx, call)
		r.ReleaseArgs()
		if err == nil {
			aq.fulfill(call.results)
			r.Returner.Return(nil)
		} else {
			aq.reject(err)
			r.Returner.Return(err)
		}
		srv.mu.Lock()
		srv.ongoing[id].cancel()
		srv.ongoing[id] = cstate{}
		if srv.drain != nil && !srv.hasOngoing() {
			close(srv.drain)
		}
		if srv.full != nil {
			close(srv.full)
			srv.full = nil
		}
		srv.mu.Unlock()
		close(done)
	}()
	var pcall capnp.PipelineCaller
	select {
	case <-ack:
		pcall = aq
	case <-done:
		// Implementation functions may not call Ack, which is fine for
		// smaller functions.
	}
	srv.mu.Lock()
	srv.starting = nil
	close(starting)
	srv.mu.Unlock()
	return pcall
}

// nextID returns the next available index in srv.ongoing or -1 if
// there are too many ongoing calls.  The caller must be holding onto
// srv.mu.
func (srv *Server) nextID() int {
	for i := range srv.ongoing {
		if srv.ongoing[i].cancel == nil {
			return i
		}
	}
	return -1
}

// hasOngoing reports whether there are any ongoing calls.
// The caller must be holding onto srv.mu.
func (srv *Server) hasOngoing() bool {
	for i := range srv.ongoing {
		if srv.ongoing[i].cancel != nil {
			return true
		}
	}
	return false
}

// Brand returns a value that will match IsServer.
func (srv *Server) Brand() capnp.Brand {
	return capnp.Brand{Value: serverBrand{srv.brand}}
}

// Shutdown waits for ongoing calls to finish and calls Shutdown on the
// Shutdowner passed into NewServer.  Shutdown must not be called more
// than once.
func (srv *Server) Shutdown() {
	srv.mu.Lock()
	if srv.drain != nil {
		srv.mu.Unlock()
		panic("capnp server: Shutdown called multiple times")
	}
	srv.drain = make(chan struct{})
	if srv.hasOngoing() {
		for _, cs := range srv.ongoing {
			if cs.cancel != nil {
				cs.cancel()
			}
		}
		srv.mu.Unlock()
		<-srv.drain
	} else {
		close(srv.drain)
		srv.mu.Unlock()
	}
	if srv.shutdown != nil {
		srv.shutdown.Shutdown()
	}
}

// IsServer reports whether a brand returned by capnp.Client.Brand
// originated from Server.Brand, and returns the brand argument passed
// to New.
func IsServer(brand capnp.Brand) (_ interface{}, ok bool) {
	sb, ok := brand.Value.(serverBrand)
	return sb.x, ok
}

type serverBrand struct {
	x interface{}
}

func sendArgsToStruct(s capnp.Send) (capnp.Struct, error) {
	if s.PlaceArgs == nil {
		return capnp.Struct{}, nil
	}
	st, err := newBlankStruct(s.ArgsSize)
	if err != nil {
		return capnp.Struct{}, err
	}
	if err := s.PlaceArgs(st); err != nil {
		st.Message().Reset(nil)
		// Using fmt.Errorf to ensure sendArgsToStruct returns a generic error.
		return capnp.Struct{}, fmt.Errorf("place args: %v", err)
	}
	return st, nil
}

func newBlankStruct(sz capnp.ObjectSize) (capnp.Struct, error) {
	_, seg, err := capnp.NewMessage(capnp.MultiSegment(nil))
	if err != nil {
		return capnp.Struct{}, err
	}
	st, err := capnp.NewRootStruct(seg, sz)
	if err != nil {
		return capnp.Struct{}, err
	}
	return st, nil
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

type resultsAllocer interface {
	AllocResults(capnp.ObjectSize) (capnp.Struct, error)
}

func newError(msg string) error {
	return exc.New(exc.Failed, "capnp server", msg)
}

func errorf(format string, args ...interface{}) error {
	return newError(fmt.Sprintf(format, args...))
}
