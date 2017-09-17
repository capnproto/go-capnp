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

// A Method describes a single RPC method on a server object.
type Method struct {
	capnp.Method
	Impl        func(context.Context, Call) error
	ResultsSize capnp.ObjectSize
}

// Call holds the state of an ongoing RPC method call.
type Call struct {
	// Args is a struct holding the call's arguments.
	Args capnp.Struct

	// Results is a struct that has enough space to hold the call's results.
	Results capnp.Struct

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
	Ack func()
}

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
	policy  Policy

	mu      sync.Mutex      // protects the following fields.  Should never be held while calling application code.
	inImpl  <-chan struct{} // non-nil if inside application code, closed when done
	drain   chan struct{}   // non-nil if draining, closed when drained
	ongoing []cstate
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

	// AnswerQueueSize is the maximum number of methods allowed to be
	// enqueued on a single returned Answer while it is unresolved.
	// Attempts to enqueue more calls than this limit will result in
	// immediate error answers.
	//
	// If this is zero, then a reasonably small default is used.
	AnswerQueueSize int
}

// New returns a client that makes calls to a set of methods.
// If closer is nil then the client's Close is a no-op.  The server
// guarantees message delivery order by blocking each call on the
// return or acknowledgment of the previous call.  See the Ack function
// for more details.
func New(methods []Method, brand interface{}, closer Closer, policy *Policy) *Server {
	srv := &Server{
		methods: make(sortedMethods, len(methods)),
		brand:   brand,
		closer:  closer,
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
	if srv.policy.AnswerQueueSize < 1 {
		srv.policy.AnswerQueueSize = 8
	}
	return srv
}

// Send starts a method call.
func (srv *Server) Send(ctx context.Context, s capnp.Send) (*capnp.Answer, capnp.ReleaseFunc) {
	mm := srv.methods.find(s.Method)
	if mm == nil {
		// TODO(soon): signal unimplemented.
		return capnp.ErrorAnswer(errors.New("unimplemented")), func() {}
	}
	args, err := sendArgsToStruct(s)
	if err != nil {
		return capnp.ErrorAnswer(err), func() {}
	}
	r := capnp.Recv{
		Args: args,
		ReleaseArgs: func() {
			// TODO(someday): log error from ClearCaps
			if seg := args.Segment(); seg != nil {
				seg.Message().Reset(nil)
			}
		},
	}
	return srv.start(ctx, mm, r)
}

// Recv starts a method call.
func (srv *Server) Recv(ctx context.Context, r capnp.Recv) (*capnp.Answer, capnp.ReleaseFunc) {
	mm := srv.methods.find(r.Method)
	if mm == nil {
		// TODO(soon): signal unimplemented.
		return capnp.ErrorAnswer(errors.New("unimplemented")), func() {}
	}
	return srv.start(ctx, mm, r)
}

func (srv *Server) start(ctx context.Context, m *Method, r capnp.Recv) (*capnp.Answer, capnp.ReleaseFunc) {
	results, err := newBlankStruct(m.ResultsSize)
	if err != nil {
		r.ReleaseArgs()
		return capnp.ErrorAnswer(err), func() {}
	}
	resultMsg := results.Message()
	ack := make(chan struct{})

	aq := newAnswerQueue(srv.policy.AnswerQueueSize)
	p := capnp.NewPromise(aq)

	// Acquire lock and block until all other concurrent calls have acknowledged delivery.
	srv.mu.Lock()
	for {
		if srv.drain != nil {
			srv.mu.Unlock()
			r.ReleaseArgs()
			return capnp.ErrorAnswer(errors.New("capnp server: call after Close")), func() {}
		}
		if srv.inImpl == nil {
			break
		}
		wait := srv.inImpl
		srv.mu.Unlock()
		select {
		case <-wait:
		case <-ctx.Done():
			r.ReleaseArgs()
			return capnp.ErrorAnswer(err), func() {}
		}
		srv.mu.Lock()
	}

	id := srv.nextID()
	if id == -1 {
		// TODO(soon): block (backpressure) instead of drop.
		srv.mu.Unlock()
		r.ReleaseArgs()
		return capnp.ErrorAnswer(errors.New("capnp server: too many concurrent calls")), func() {}
	}
	ctx, cancel := context.WithCancel(ctx)
	srv.ongoing[id] = cstate{cancel}
	done := make(chan struct{})
	srv.inImpl = done
	srv.mu.Unlock()
	defer func() {
		srv.mu.Lock()
		srv.inImpl = nil
		close(done)
		srv.mu.Unlock()
	}()

	go func() {
		once := new(sync.Once)
		err := m.Impl(ctx, Call{
			Args:    r.Args,
			Results: results,
			Ack: func() {
				once.Do(func() { close(ack) })
			},
		})
		r.ReleaseArgs()
		if err == nil {
			aq.drain(capnp.ImmediateAnswer(results), p)
		} else {
			aq.drain(capnp.ErrorAnswer(err), p)
			// TODO(someday): log error from ClearCaps
			resultMsg.Reset(nil)
		}
		srv.mu.Lock()
		srv.ongoing[id].cancel()
		srv.ongoing[id] = cstate{}
		if srv.drain != nil && !srv.hasOngoing() {
			close(srv.drain)
		}
		srv.mu.Unlock()
	}()
	ans := p.Answer()
	select {
	case <-ack:
	case <-ans.Done():
		// Implementation functions may not call Ack, which is fine for
		// smaller functions.
	}
	once := new(sync.Once)
	return ans, func() {
		once.Do(func() {
			<-ans.Done()
			p.ReleaseClients()
			// TODO(someday): log error from ClearCaps
			resultMsg.Reset(nil)
		})
	}
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
func (srv *Server) Brand() interface{} {
	return serverBrand{srv.brand}
}

// Close waits for ongoing calls to finish and calls Close to the Closer.
func (srv *Server) Close() error {
	srv.mu.Lock()
	if srv.drain != nil {
		srv.mu.Unlock()
		return errors.New("capnp server: Close called multiple times")
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
	if srv.closer == nil {
		return nil
	}
	return srv.closer.Close()
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
		return capnp.Struct{}, err
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
