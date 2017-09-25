package capnp

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"sync"
)

// An Interface is a reference to a client in a message's capability table.
type Interface struct {
	seg *Segment
	cap CapabilityID
}

// NewInterface creates a new interface pointer.
//
// No allocation is performed in the given segment: it is used purely
// to associate the interface pointer with a message.
func NewInterface(s *Segment, cap CapabilityID) Interface {
	return Interface{s, cap}
}

// ToPtr converts the interface to a generic pointer.
func (p Interface) ToPtr() Ptr {
	return Ptr{
		seg:      p.seg,
		lenOrCap: uint32(p.cap),
		flags:    interfacePtrFlag,
	}
}

// Message returns the message whose capability table the interface
// references or nil if the pointer is invalid.
func (i Interface) Message() *Message {
	if i.seg == nil {
		return nil
	}
	return i.seg.msg
}

// IsValid returns whether the interface is valid.
func (i Interface) IsValid() bool {
	return i.seg != nil
}

// Capability returns the capability ID of the interface.
func (i Interface) Capability() CapabilityID {
	return i.cap
}

// value returns a raw interface pointer with the capability ID.
func (i Interface) value(paddr Address) rawPointer {
	if i.seg == nil {
		return 0
	}
	return rawInterfacePointer(i.cap)
}

// Client returns the client stored in the message's capability table
// or nil if the pointer is invalid.
func (i Interface) Client() *Client {
	msg := i.Message()
	if msg == nil {
		return nil
	}
	tab := msg.CapTable
	if int64(i.cap) >= int64(len(tab)) {
		return nil
	}
	return tab[i.cap]
}

// A CapabilityID is an index into a message's capability table.
type CapabilityID uint32

// A Client is a reference to a Cap'n Proto capability.
// The zero value is a null capability reference.
// It is safe to use from multiple goroutines.
type Client struct {
	mu       sync.Mutex  // protects the struct
	h        *clientHook // nil if resolved to nil or released
	released bool
}

// clientHook is a reference-counted wrapper for a ClientHook.
// It is assumed that a clientHook's address uniquely identifies a hook,
// since they are only created in NewClient and NewPromisedClient.
type clientHook struct {
	// ClientHook will never be nil and will not change for the lifetime of
	// a clientHook.
	ClientHook

	// done is closed when refs == 0 and calls == 0.
	done chan struct{}

	// resolved is closed after resolvedHook is set
	resolved chan struct{}

	mu           sync.Mutex
	refs         int         // how many open Clients reference this clientHook
	calls        int         // number of outstanding ClientHook accesses
	resolvedHook *clientHook // valid only if resolved is closed
}

// NewClient creates the first reference to a capability.
// If hook is nil, then NewClient returns nil.
//
// Typically the RPC system will create a client for the application.
// Most applications will not need to use this directly.
func NewClient(hook ClientHook) *Client {
	if hook == nil {
		return nil
	}
	h := &clientHook{
		ClientHook: hook,
		done:       make(chan struct{}),
		refs:       1,
		resolved:   newClosedSignal(),
	}
	h.resolvedHook = h
	return &Client{h: h}
}

// NewPromisedClient creates the first reference to a capability that
// can resolve to a different capability.  The hook will be shut down
// when the promise is resolved or the client has no more references,
// whichever comes first.
//
// Typically the RPC system will create a client for the application.
// Most applications will not need to use this directly.
func NewPromisedClient(hook ClientHook) (*Client, *ClientPromise) {
	if hook == nil {
		panic("NewPromisedClient(nil)")
	}
	h := &clientHook{
		ClientHook: hook,
		done:       make(chan struct{}),
		refs:       1,
		resolved:   make(chan struct{}),
	}
	return &Client{h: h}, &ClientPromise{h: h}
}

// startCall holds onto a hook to prevent it from shutting down until
// finish is called.  It resolves the client's hook as much as possible
// first.  The caller must not be holding onto c.mu.
func (c *Client) startCall() (hook ClientHook, released bool, finish func()) {
	if c == nil {
		return nil, false, func() {}
	}
	defer c.mu.Unlock()
	c.mu.Lock()
	if c.h == nil {
		return nil, c.released, func() {}
	}
	c.h.mu.Lock()
	c.h = resolveHook(c.h)
	if c.h == nil {
		return nil, false, func() {}
	}
	c.h.calls++
	c.h.mu.Unlock()
	savedHook := c.h
	return savedHook.ClientHook, false, func() {
		savedHook.mu.Lock()
		savedHook.calls--
		if savedHook.refs == 0 && savedHook.calls == 0 {
			close(savedHook.done)
		}
		savedHook.mu.Unlock()
	}
}

func (c *Client) peek() (hook *clientHook, released bool, resolved bool) {
	if c == nil {
		return nil, false, true
	}
	defer c.mu.Unlock()
	c.mu.Lock()
	if c.h == nil {
		return nil, c.released, true
	}
	c.h.mu.Lock()
	c.h = resolveHook(c.h)
	if c.h == nil {
		return nil, false, true
	}
	resolved = c.h.isResolved()
	c.h.mu.Unlock()
	return c.h, false, resolved
}

// resolveHook resolves h as much as possible without blocking.
// The caller must be holding onto h.mu and when resolveHook returns, it
// will be holding onto the mutex of the returned hook if not nil.
func resolveHook(h *clientHook) *clientHook {
	for {
		if !h.isResolved() {
			return h
		}
		r := h.resolvedHook
		if r == h {
			return h
		}
		h.mu.Unlock()
		h = r
		if h == nil {
			return nil
		}
		h.mu.Lock()
	}
}

// SendCall allocates space for parameters, calls args.Place to fill out
// the parameters, then starts executing a method, returning an answer
// that will hold the result.  The caller must call the returned release
// function when it no longer needs the answer's data.
func (c *Client) SendCall(ctx context.Context, s Send) (*Answer, ReleaseFunc) {
	h, released, finish := c.startCall()
	defer finish()
	if released {
		return ErrorAnswer(errors.New("capnp: call on released client")), func() {}
	}
	if h == nil {
		return ErrorAnswer(errors.New("capnp: call on null client")), func() {}
	}
	return h.Send(ctx, s)
}

// RecvCall starts executing a method with the referenced arguments
// and returns an answer that will hold the result.  The hook will call
// a.Release when it no longer needs to reference the parameters.  The
// caller must call the returned release function when it no longer
// needs the answer's data.
func (c *Client) RecvCall(ctx context.Context, r Recv) (*Answer, ReleaseFunc) {
	h, released, finish := c.startCall()
	defer finish()
	if released {
		return ErrorAnswer(errors.New("capnp: call on released client")), func() {}
	}
	if h == nil {
		return ErrorAnswer(errors.New("capnp: call on null client")), func() {}
	}
	return h.Recv(ctx, r)
}

// IsValid reports whether c is a valid reference to a capability.
// A reference is invalid if it is nil, has resolved to null, or has
// been released.
func (c *Client) IsValid() bool {
	h, released, _ := c.peek()
	return !released && h != nil
}

// IsSame reports whether c and c2 refer to a capability created by the
// same call to NewClient.  This can return false negatives if c or c2
// are not fully resolved: use Resolve if this is an issue.  If either
// c or c2 are released, then IsSame panics.
func (c *Client) IsSame(c2 *Client) bool {
	h1, released, _ := c.peek()
	if released {
		panic("IsSame on released client")
	}
	h2, released, _ := c2.peek()
	if released {
		panic("IsSame on released client")
	}
	return h1 == h2
}

// Resolve blocks until the capability is fully resolved or the Context is Done.
func (c *Client) Resolve(ctx context.Context) error {
	for {
		h, released, resolved := c.peek()
		if released {
			return errors.New("capnp: cannot resolve released client")
		}
		if resolved {
			return nil
		}
		select {
		case <-h.resolved:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// AddRef creates a new Client that refers to the same capability as c.
// If c is nil or has resolved to null, then AddRef returns nil.
func (c *Client) AddRef() *Client {
	if c == nil {
		return nil
	}
	defer c.mu.Unlock()
	c.mu.Lock()
	if c.released {
		panic("AddRef on released client")
	}
	if c.h == nil {
		return nil
	}
	c.h.mu.Lock()
	c.h = resolveHook(c.h)
	if c.h == nil {
		return nil
	}
	c.h.refs++
	c.h.mu.Unlock()
	return &Client{h: c.h}
}

// WeakRef creates a new WeakClient that refers to the same capability
// as c.  If c is nil or has resolved to null, then WeakRef returns nil.
func (c *Client) WeakRef() *WeakClient {
	h, released, _ := c.peek()
	if released {
		panic("WeakRef on released client")
	}
	if h == nil {
		return nil
	}
	return &WeakClient{h: h}
}

// Brand returns the current underlying hook's Brand method or nil if
// c is nil, has resolved to null, or has been released.
func (c *Client) Brand() interface{} {
	h, _, finish := c.startCall()
	defer finish()
	if h == nil {
		return nil
	}
	return h.Brand()
}

// String returns a string that identifies this capability for debugging
// purposes.  Its format should not be depended on: in particular, it
// should not be used to compare clients.  Use IsSame to compare clients
// for equality.
func (c *Client) String() string {
	h, released, resolved := c.peek()
	if released {
		return "<released client>"
	}
	if h == nil {
		return "<nil>"
	}
	if !resolved {
		return fmt.Sprintf("<unresolved client %p>", h)
	}
	return fmt.Sprintf("<client %p>", h)
}

// Release releases a capability reference.  If this is the last
// reference to the capability, then the underlying resources associated
// with the capability will be released.
//
// Release will panic if c has already been released, but not if c is
// nil or resolved to null.
func (c *Client) Release() {
	if c == nil {
		return
	}
	c.mu.Lock()
	if c.released {
		c.mu.Unlock()
		panic("capnp: double Client.Release")
	}
	if c.h == nil {
		c.mu.Unlock()
		return
	}
	c.released = true
	c.h.mu.Lock()
	c.h = resolveHook(c.h)
	if c.h == nil {
		c.mu.Unlock()
		return
	}
	h := c.h
	c.h = nil
	h.refs--
	if h.refs > 0 {
		h.mu.Unlock()
		c.mu.Unlock()
		return
	}
	if h.calls == 0 {
		close(h.done)
	}
	h.mu.Unlock()
	c.mu.Unlock()
	<-h.done
	h.Shutdown()
}

// isResolve reports whether ch has been resolved.
// The caller must be holding onto ch.mu.
func (ch *clientHook) isResolved() bool {
	select {
	case <-ch.resolved:
		return true
	default:
		return false
	}
}

// A ClientPromise resolves the identity of a client created by NewPromisedClient.
type ClientPromise struct {
	h *clientHook
}

// Fulfill resolves the client promise to c.  After Fulfill returns,
// then all future calls to the client created by NewPromisedClient will
// be sent to c.  It is guaranteed that the hook passed to
// NewPromisedClient will be shut down after Fulfill returns, but the
// hook may have been shut down earlier if the client ran out of
// references.
func (cp *ClientPromise) Fulfill(c *Client) {
	// Obtain next client hook.
	var rh *clientHook
	if c != nil {
		c.mu.Lock()
		if c.released {
			c.mu.Unlock()
			panic("ClientPromise.Resolve with a released client")
		}
		// TODO(maybe): c.h = resolveHook(c.h)
		rh = c.h
		c.mu.Unlock()
	}

	// Mark hook as resolved.
	cp.h.mu.Lock()
	if cp.h.isResolved() {
		cp.h.mu.Unlock()
		panic("ClientPromise.Resolve called more than once")
	}
	cp.h.resolvedHook = rh
	close(cp.h.resolved)
	refs := cp.h.refs
	cp.h.refs = 0
	if refs == 0 {
		cp.h.mu.Unlock()
		return
	}

	// Client still had references, so we're responsible for shutting it down.
	if cp.h.calls == 0 {
		close(cp.h.done)
	}
	rh = resolveHook(cp.h) // swaps mutex on cp.h for mutex on rh
	if rh != nil {
		rh.refs += refs
		rh.mu.Unlock()
	}
	<-cp.h.done
	cp.h.Shutdown()
}

// A WeakClient is a weak reference to a capability: it refers to a
// capability without preventing it from being shut down.  The zero
// value is a null reference.
type WeakClient struct {
	h *clientHook
}

// AddRef creates a new Client that refers to the same capability as c
// as long as the capability hasn't already been shut down.
func (wc *WeakClient) AddRef() (c *Client, ok bool) {
	if wc == nil {
		return nil, true
	}
	if wc.h == nil {
		return nil, true
	}
	wc.h.mu.Lock()
	wc.h = resolveHook(wc.h)
	if wc.h == nil {
		return nil, true
	}
	if wc.h.refs == 0 {
		wc.h.mu.Unlock()
		return nil, false
	}
	wc.h.refs++
	wc.h.mu.Unlock()
	return &Client{h: wc.h}, true
}

// A ClientHook represents a Cap'n Proto capability.  Application code
// should not pass around ClientHooks; applications should pass around
// Clients.  A ClientHook must be safe to use from multiple goroutines.
//
// Calls must be delivered to the capability in the order they are made.
// This guarantee is based on the concept of a capability
// acknowledging delivery of a call: this is specific to an
// implementation of ClientHook.  A type that implements ClientHook
// must guarantee that if foo() then bar() is called on a client, then
// the capability acknowledging foo() happens before the capability
// observing bar().
type ClientHook interface {
	// Send allocates space for parameters, calls args.Place to fill out
	// the arguments, then starts executing a method, returning an answer
	// that will hold the result.  The caller must call the returned
	// release function when it no longer needs the answer's data.
	//
	// Send is typically used when application code is making a call.
	Send(ctx context.Context, s Send) (*Answer, ReleaseFunc)

	// Recv starts executing a method with the referenced arguments
	// and returns an answer that will hold the result.  The hook will call
	// args.Release when it no longer needs to reference the parameters.
	// The caller must call the returned release function when it no longer
	// needs the answer's data.
	//
	// Recv is typically used when the RPC system has received a call.
	Recv(ctx context.Context, r Recv) (*Answer, ReleaseFunc)

	// Brand returns an implementation-specific value.  This can be used
	// to introspect and identify kinds of clients.
	Brand() interface{}

	// Shutdown releases any resources associated with this capability.
	// The behavior of calling any methods on the receiver after calling
	// Shutdown is undefined.  It is expected for the ClientHook to reject
	// any outstanding call futures.
	Shutdown()
}

// Send is the input to ClientHook.Send.
type Send struct {
	// Method must have InterfaceID and MethodID filled in.
	Method Method

	// PlaceArgs is a function that will be called at most once before Send
	// returns to populate the arguments for the RPC.  PlaceArgs may be nil.
	PlaceArgs func(Struct) error

	// ArgsSize specifies the size of the struct to pass to PlaceArgs.
	ArgsSize ObjectSize
}

// Recv is the input to ClientHook.Recv.
type Recv struct {
	// Method must have InterfaceID and MethodID filled in.
	Method Method

	// Args is the set of arguments for the RPC.
	Args Struct

	// ReleaseArgs is called after Args is no longer referenced.
	// Must not be nil.
	ReleaseArgs ReleaseFunc
}

// A ReleaseFunc tells the RPC system that a parameter or result struct
// is no longer in use and may be reclaimed.  After the first call,
// subsequent calls to a ReleaseFunc do nothing.
type ReleaseFunc func()

// A Method identifies a method along with an optional human-readable
// description of the method.
type Method struct {
	InterfaceID uint64
	MethodID    uint16

	// Canonical name of the interface.  May be empty.
	InterfaceName string
	// Method name as it appears in the schema.  May be empty.
	MethodName string
}

// String returns a formatted string containing the interface name or
// the method name if present, otherwise it uses the raw IDs.
// This is suitable for use in error messages and logs.
func (m *Method) String() string {
	buf := make([]byte, 0, 128)
	if m.InterfaceName == "" {
		buf = append(buf, '@', '0', 'x')
		buf = strconv.AppendUint(buf, m.InterfaceID, 16)
	} else {
		buf = append(buf, m.InterfaceName...)
	}
	buf = append(buf, '.')
	if m.MethodName == "" {
		buf = append(buf, '@')
		buf = strconv.AppendUint(buf, uint64(m.MethodID), 10)
	} else {
		buf = append(buf, m.MethodName...)
	}
	return string(buf)
}

type errorClient struct {
	e error
}

// ErrorClient returns a Client that always returns error e.
// An ErrorClient does not need to be released: it is a sentinel like a
// nil Client.
func ErrorClient(e error) *Client {
	return NewClient(errorClient{e})
}

func (ec errorClient) Send(context.Context, Send) (*Answer, ReleaseFunc) {
	return ErrorAnswer(ec.e), func() {}
}

func (ec errorClient) Recv(_ context.Context, r Recv) (*Answer, ReleaseFunc) {
	r.ReleaseArgs()
	return ErrorAnswer(ec.e), func() {}
}

func (ec errorClient) Brand() interface{} {
	return nil
}

func (ec errorClient) Shutdown() {
}

// IsUnimplemented reports whether e indicates an unimplemented method error.
func IsUnimplemented(e error) bool {
	// TODO(soon)
	return false
}

var closedSignal <-chan struct{} = newClosedSignal()

func newClosedSignal() chan struct{} {
	c := make(chan struct{})
	close(c)
	return c
}
