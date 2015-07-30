package capnp

import (
	"errors"
	"sort"
	"sync"

	"golang.org/x/net/context"
)

// A ServerMethod describes a single method on a server object.
type ServerMethod struct {
	Method
	Impl        ServerFunc
	ResultsSize ObjectSize
}

// A ServerFunc is a function that implements a single method.
type ServerFunc func(ctx context.Context, options CallOptions, params, results Struct) error

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

// NewServer returns a client that makes calls to a set of methods.
// If closer is nil then the client's Close is a no-op.  The server
// guarantees message delivery order by blocking each call on the
// return or acknowledgement of the previous call.  See the Ack function
// for more details.
func NewServer(methods []ServerMethod, closer Closer) Client {
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

func (s *server) Call(call *Call) Answer {
	// Find method
	sm := s.methods.find(&call.Method)
	if sm == nil {
		return ErrorAnswer(&MethodError{
			Method: &call.Method,
			Err:    ErrUnimplemented,
		})
	}
	// Acquire lock for ordering
	select {
	case <-s.lock:
		defer func() {
			s.lock <- struct{}{}
		}()
	case <-call.Ctx.Done():
		return ErrorAnswer(call.Ctx.Err())
	}
	// Call implementation function
	ans := new(Fulfiller)
	params := call.PlaceParams(nil)
	out := NewBuffer(nil)
	results := out.NewRootStruct(sm.ResultsSize)
	acksig := &ackSignal{c: make(chan struct{})}
	opts := call.Options.With([]CallOption{SetOptionValue(ackSignalKey, acksig)})
	go func() {
		err := sm.Impl(call.Ctx, opts, params, results)
		if err == nil {
			ans.Fulfill(ImmediateAnswer(Object(results)))
		} else {
			ans.Fulfill(ErrorAnswer(err))
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
		return ErrorAnswer(call.Ctx.Err())
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
//	func (my *myServer) MyMethod(
//		ctx context.Context,
//		opts capnp.CallOptions,
//		p Interface_myMethod_Params,
//		r Interface_myMethod_Results) error {
//		capnp.Ack(opts)
//		// ... do long-running operation ...
//		return nil
//	}
//
// Ack need not be the first call in a function nor is it required.
// Since the function's return is also an acknowledgement of delivery,
// short functions can return without calling Ack.  However, since
// clients will not return an Answer until the delivery is acknowledged,
// it is advisable to ack early.
func Ack(opts CallOptions) {
	if ack, _ := opts.Value(ackSignalKey).(*ackSignal); ack != nil {
		ack.signal()
	}
}

// MethodError is an error on an associated method.
type MethodError struct {
	Method *Method
	Err    error
}

// Error returns the method name concatenated with the error string.
func (me *MethodError) Error() string {
	return me.Method.String() + ": " + me.Err.Error()
}

// ErrUnimplemented is the error returned when a method is called on
// a server that does not implement the method.
var ErrUnimplemented = errors.New("method not implemented")

// IsUnimplemented reports whether e indicates an unimplemented method error.
func IsUnimplemented(e error) bool {
	if me, ok := e.(*MethodError); ok {
		e = me
	}
	return e == ErrUnimplemented
}

type sortedMethods []ServerMethod

// find returns the method with the given ID or nil.
func (sm sortedMethods) find(id *Method) *ServerMethod {
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
