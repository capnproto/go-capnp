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
	mu     sync.Mutex  // protects the struct
	h      *clientHook // nil if resolved to nil or closed
	closed bool
}

// clientHook is a reference-counted wrapper for a ClientHook.
// It is assumed that a clientHook's address uniquely identifies a hook,
// since they are only created in NewClient and NewPromisedClient.
type clientHook struct {
	// ClientHook will never be nil and will not change for the lifetime of
	// a clientHook.
	ClientHook

	// closer will be non-nil if created by NewClient and will not change
	// for the lifetime of a clientHook.
	closer

	// done is closed when refs == 0 and calls == 0.
	// Only non-nil for clients created by NewClient.
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
func NewClient(hook ClientHookCloser) *Client {
	if hook == nil {
		return nil
	}
	h := &clientHook{
		ClientHook: hook,
		closer:     hook,
		done:       make(chan struct{}),
		refs:       1,
		resolved:   newClosedSignal(),
	}
	h.resolvedHook = h
	return &Client{h: h}
}

// NewPromisedClient creates the first reference to a capability that
// can resolve to a different capability.
//
// Typically the RPC system will create a client for the application.
// Most applications will not need to use this directly.
func NewPromisedClient(hook ClientHook) (*Client, *ClientPromise) {
	if hook == nil {
		panic("NewPromisedClient(nil)")
	}
	h := &clientHook{
		ClientHook: hook,
		refs:       2, // to keep open until resolved
		resolved:   make(chan struct{}),
	}
	return &Client{h: h}, &ClientPromise{h: h}
}

// startCall holds onto a hook to prevent it from closing until finish is called.
// It resolves the client's hook as much as possible first.
// The caller must not be holding onto c.mu.
func (c *Client) startCall() (hook ClientHook, closed bool, finish func()) {
	if c == nil {
		return nil, false, func() {}
	}
	defer c.mu.Unlock()
	c.mu.Lock()
	if c.h == nil {
		return nil, c.closed, func() {}
	}
	c.h.mu.Lock()
	c.lockedResolve()
	if c.h == nil {
		return nil, false, func() {}
	}
	c.h.calls++
	c.h.mu.Unlock()
	savedHook := c.h
	return savedHook.ClientHook, false, func() {
		savedHook.mu.Lock()
		savedHook.calls--
		if savedHook.refs == 0 && savedHook.calls == 0 && savedHook.done != nil {
			close(savedHook.done)
		}
		savedHook.mu.Unlock()
	}
}

func (c *Client) peek() (hook *clientHook, closed bool, resolved bool) {
	if c == nil {
		return nil, false, true
	}
	defer c.mu.Unlock()
	c.mu.Lock()
	if c.h == nil {
		return nil, c.closed, true
	}
	c.h.mu.Lock()
	c.lockedResolve()
	if c.h == nil {
		return nil, false, true
	}
	select {
	case <-c.h.resolved:
		resolved = true
	default:
	}
	c.h.mu.Unlock()
	return c.h, false, resolved
}

// lockedResolve resolves c.h as much as possible.
// The caller must be holding onto c.mu and c.h.mu.
func (c *Client) lockedResolve() {
	for {
		select {
		case <-c.h.resolved:
		default:
			return
		}
		rh := c.h.resolvedHook
		if rh == c.h {
			return
		}
		c.h.refs--
		c.h.mu.Unlock()
		c.h = rh
		if rh == nil {
			return
		}
		rh.mu.Lock()
	}
}

// SendCall allocates space for parameters, calls args.Place to fill out
// the parameters, then starts executing a method, returning an answer
// that will hold the result.  The caller must call the returned release
// function when it no longer needs the answer's data.
func (c *Client) SendCall(ctx context.Context, s Send) (*Answer, ReleaseFunc) {
	h, closed, finish := c.startCall()
	defer finish()
	if closed {
		return ErrorAnswer(errors.New("capnp: call on closed client")), func() {}
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
	h, closed, finish := c.startCall()
	defer finish()
	if closed {
		return ErrorAnswer(errors.New("capnp: call on closed client")), func() {}
	}
	if h == nil {
		return ErrorAnswer(errors.New("capnp: call on null client")), func() {}
	}
	return h.Recv(ctx, r)
}

// IsValid reports whether c is a valid reference to a capability.
// A reference is invalid if it is nil, has resolved to null, or has
// been closed.
func (c *Client) IsValid() bool {
	h, closed, _ := c.peek()
	return !closed && h != nil
}

// IsSame reports whether c and c2 refer to a capability created by the
// same call to NewClient.  This can return false negatives if c or c2
// are not fully resolved: use Resolve if this is an issue.  If either
// c or c2 are closed, then IsSame panics.
func (c *Client) IsSame(c2 *Client) bool {
	h1, closed, _ := c.peek()
	if closed {
		panic("IsSame on closed client")
	}
	h2, closed, _ := c2.peek()
	if closed {
		panic("IsSame on closed client")
	}
	return h1 == h2
}

// Resolve blocks until the capability is fully resolved or the Context is Done.
func (c *Client) Resolve(ctx context.Context) error {
	for {
		h, closed, resolved := c.peek()
		if closed {
			return errors.New("capnp: cannot resolve closed client")
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
	if c.closed {
		panic("AddRef on closed client")
	}
	if c.h == nil {
		return nil
	}
	c.h.mu.Lock()
	c.lockedResolve()
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
	h, closed, _ := c.peek()
	if closed {
		panic("WeakRef on closed client")
	}
	if h == nil {
		return nil
	}
	return &WeakClient{h: h}
}

// Brand returns the current underlying hook's Brand method or nil if
// c is nil, has resolved to null, or has been closed.
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
	h, closed, resolved := c.peek()
	if closed {
		return "<closed client>"
	}
	if h == nil {
		return "<nil>"
	}
	if !resolved {
		return fmt.Sprintf("<unresolved client %p>", h)
	}
	return fmt.Sprintf("<client %p>", h)
}

// Close releases a capability reference.  If this is the last reference
// to the capability, then the underlying resources associated with the
// capability will be released and any error will be returned.
//
// Close also returns an error if c has already been closed, but not if
// c is nil or resolved to null.
func (c *Client) Close() error {
	if c == nil {
		return nil
	}
	c.mu.Lock()
	if c.closed {
		c.mu.Unlock()
		return errors.New("capnp: double close on Client")
	}
	if c.h == nil {
		c.mu.Unlock()
		return nil
	}
	c.closed = true
	c.h.mu.Lock()
	c.lockedResolve()
	if c.h == nil {
		c.mu.Unlock()
		return nil
	}
	h := c.h
	c.h = nil
	h.refs--
	if h.refs > 0 {
		h.mu.Unlock()
		c.mu.Unlock()
		return nil
	}
	done := h.done
	if done == nil {
		h.mu.Unlock()
		c.mu.Unlock()
		return nil
	}
	if h.calls == 0 {
		close(done)
	}
	h.mu.Unlock()
	c.mu.Unlock()
	<-done
	return h.closer.Close()
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
// be sent to c.
func (cp *ClientPromise) Fulfill(c *Client) {
	// Obtain next client hook.
	var rh *clientHook
	if c != nil {
		c.mu.Lock()
		if c.closed {
			c.mu.Unlock()
			panic("ClientPromise.Resolve with a closed client")
		}
		// TODO(maybe): c.lockedResolve()?
		rh = c.h
		c.mu.Unlock()
	}

	// Mark hook as resolved.
	cp.h.mu.Lock()
	select {
	case <-cp.h.resolved:
		cp.h.mu.Unlock()
		panic("ClientPromise.Resolve called more than once")
	default:
	}
	cp.h.resolvedHook = rh
	close(cp.h.resolved)
	refs := cp.h.refs - 1
	cp.h.refs = refs
	if rh == nil || rh == cp.h || refs == 0 {
		cp.h.mu.Unlock()
		return
	}

	// Add refs down resolution chain.
	rh.mu.Lock()
	cp.h.mu.Unlock()
	for rh.isResolved() && rh.resolvedHook != rh {
		rh.refs += refs
		rh2 := rh.resolvedHook
		if rh2 == nil {
			rh.mu.Unlock()
			return
		}
		rh2.mu.Lock()
		rh.mu.Unlock()
		rh = rh2
	}
	rh.refs += refs
	rh.mu.Unlock()
}

// A WeakClient is a weak reference to a capability: it refers to a
// capability without preventing it from being closed.  The zero value
// is a null reference.
type WeakClient struct {
	h *clientHook
}

// AddRef creates a new Client that refers to the same capability as c
// as long as the capability hasn't already been closed.
func (wc *WeakClient) AddRef() (c *Client, ok bool) {
	if wc == nil {
		return nil, true
	}
	for {
		if wc.h == nil {
			return nil, true
		}
		wc.h.mu.Lock()
		if wc.h.isResolved() {
			if r := wc.h.resolvedHook; r != wc.h {
				wc.h.mu.Unlock()
				wc.h = r
				continue
			}
		}
		if wc.h.refs == 0 {
			wc.h.mu.Unlock()
			return nil, false
		}
		wc.h.refs++
		wc.h.mu.Unlock()
		return &Client{h: wc.h}, true
	}
}

// A ClientHook represents a Cap'n Proto capability.  Application code
// should not pass around ClientHooks; applications should pass around
// Clients.  A ClientHook must be safe to use from multiple goroutines.
//
// Calls must be delivered to the capability in the order they are made.
// This guarantee is based on the concept of a capability
// acknowledging delivery of a call: this is specific to an
// implementation of ClientHook.  A type that implements ClientHook
// must guarantee that if foo() then bar() is called on a client, that
// acknowledging foo() happens before acknowledging bar().
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

	// Options is the set of call options.
	Options CallOptions
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

	// Options is the set of call options.
	Options CallOptions
}

// ClientHookCloser is the interface that groups the ClientHook and
// Close methods.
type ClientHookCloser interface {
	ClientHook

	// Close releases any resources associated with this capability.
	// The behavior of calling any methods on the receiver after calling
	// Close is undefined.  It is expected for the ClientHook to reject any
	// outstanding call futures.
	Close() error
}

type closer interface {
	Close() error
}

// A ReleaseFunc tells the RPC system that a parameter or result struct
// is no longer in use and may be reclaimed.  After the first call,
// subsequent calls to a ReleaseFunc do nothing.
type ReleaseFunc func()

// CallOptions holds RPC-specific options for an interface call.
// The zero value is an empty set of options.  CallOptions is safe to
// use from multiple goroutines.
//
// Its usage is similar to the values in context.Context, but is only
// used for a single call: its values are not intended to propagate to
// other callees.  An example of an option would be the
// Call.sendResultsTo field in rpc.capnp.
type CallOptions struct {
	m map[interface{}]interface{}
}

// NewCallOptions builds a CallOptions value from a list of individual options.
func NewCallOptions(opts []CallOption) CallOptions {
	co := CallOptions{make(map[interface{}]interface{})}
	for _, o := range opts {
		o.f(co)
	}
	return co
}

// Value retrieves the value associated with the options for this key,
// or nil if no value is associated with this key.
func (co CallOptions) Value(key interface{}) interface{} {
	return co.m[key]
}

// With creates a copy of the CallOptions value with other options applied.
func (co CallOptions) With(opts []CallOption) CallOptions {
	newopts := CallOptions{make(map[interface{}]interface{})}
	for k, v := range co.m {
		newopts.m[k] = v
	}
	for _, o := range opts {
		o.f(newopts)
	}
	return newopts
}

// A CallOption is a function that modifies options on an interface call.
type CallOption struct {
	f func(CallOptions)
}

// SetOptionValue returns a call option that associates a value to an
// option key.  This can be retrieved later with CallOptions.Value.
func SetOptionValue(key, value interface{}) CallOption {
	return CallOption{func(co CallOptions) {
		co.m[key] = value
	}}
}

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
// An ErrorClient does not need to be closed: it is a sentinel like a
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

func (ec errorClient) Close() error {
	return nil
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
