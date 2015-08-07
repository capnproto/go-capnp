// Package server provides runtime support for implementing Cap'n Proto
// interfaces locally.
package server

import (
	"sort"
	"sync"

	"golang.org/x/net/context"
	"zombiezen.com/go/capnproto"
	"zombiezen.com/go/capnproto/internal/fulfiller"
)

// A Method describes a single method on a server object.
type Method struct {
	capnp.Method
	Impl        Func
	ResultsSize capnp.ObjectSize
}

// A Func is a function that implements a single method.
type Func func(ctx context.Context, options capnp.CallOptions, params, results capnp.Struct) error

// Closer is the interface that wraps the Close method.
type Closer interface {
	Close() error
}

// A server is a locally implemented interface.
type server struct {
	methods sortedMethods
	closer  Closer
	lock    chan struct{}
}

// New returns a client that makes calls to a set of methods.
// If closer is nil then the client's Close is a no-op.  The server
// guarantees message delivery order by blocking each call on the
// return or acknowledgement of the previous call.  See the Ack function
// for more details.
func New(methods []Method, closer Closer) capnp.Client {
	s := &server{
		methods: make(sortedMethods, len(methods)),
		closer:  closer,
		lock:    make(chan struct{}, 1),
	}
	s.lock <- struct{}{}
	copy(s.methods, methods)
	sort.Sort(s.methods)
	return s
}

func (s *server) Call(call *capnp.Call) capnp.Answer {
	// Find method
	sm := s.methods.find(&call.Method)
	if sm == nil {
		return capnp.ErrorAnswer(&capnp.MethodError{
			Method: &call.Method,
			Err:    capnp.ErrUnimplemented,
		})
	}
	// Acquire lock for ordering
	select {
	case <-s.lock:
		defer func() {
			s.lock <- struct{}{}
		}()
	case <-call.Ctx.Done():
		return capnp.ErrorAnswer(call.Ctx.Err())
	}
	// Call implementation function
	ans := new(fulfiller.Fulfiller)
	params := call.PlaceParams(nil)
	out := capnp.NewBuffer(nil)
	results := out.NewRootStruct(sm.ResultsSize)
	acksig := &ackSignal{c: make(chan struct{})}
	opts := call.Options.With([]capnp.CallOption{capnp.SetOptionValue(ackSignalKey, acksig)})
	go func() {
		err := sm.Impl(call.Ctx, opts, params, results)
		if err == nil {
			ans.Fulfill(results)
		} else {
			ans.Reject(err)
		}
	}()
	// Wait for resolution
	select {
	case <-acksig.c:
		return ans
	case <-ans.Done():
		// Implementation functions may not call Ack, which is fine for smaller functions.
		return ans
	case <-call.Ctx.Done():
		return capnp.ErrorAnswer(call.Ctx.Err())
	}
}

func (s *server) Close() error {
	if s.closer == nil {
		return nil
	}
	return s.closer.Close()
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
// Since the function's return is also an acknowledgement of delivery,
// short functions can return without calling Ack.  However, since
// clients will not return an Answer until the delivery is acknowledged,
// it is advisable to ack early.
func Ack(opts capnp.CallOptions) {
	if ack, _ := opts.Value(ackSignalKey).(*ackSignal); ack != nil {
		ack.signal()
	}
}

type sortedMethods []Method

// find returns the method with the given ID or nil.
func (sm sortedMethods) find(id *capnp.Method) *Method {
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

func (ack *ackSignal) signal() {
	ack.once.Do(func() {
		close(ack.c)
	})
}

// callOptionKey is the unexported key type for predefined options.
type callOptionKey int

// Predefined call options
const (
	invalidOptionKey callOptionKey = iota
	ackSignalKey
)
