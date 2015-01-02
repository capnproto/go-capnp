package capn

import (
	"errors"
	"strconv"
	"sync"

	"golang.org/x/net/context"
)

var ErrNullClient = errors.New("capn: call on null client")

// A Client makes calls for an interface type.
type Client interface {
	// Call calls a method with an existing parameters struct.  This is
	// used when the RPC system receives a call for an exported object.
	Call(
		ctx context.Context,
		method *Method,
		params Struct) Promise

	// NewCall calls a method, allowing the client to allocate the
	// parameters struct.  This is used when application code is using
	// a client.
	NewCall(
		ctx context.Context,
		method *Method,
		paramsSize ObjectSize,
		params func(Struct)) Promise

	// Close releases any resources associated with this client.
	// No further calls to the client should be made after calling Close.
	Close() error
}

// A Promise is a generic promise for a Struct.
type Promise interface {
	Get() (Struct, error)
	GetClient(off int) Client
	GetPromise(off int) Promise
	GetPromiseDefault(off int, s *Segment, tagoff int) Promise
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
	buf = append(buf, '/')
	if m.MethodName == "" {
		buf = append(buf, '@')
		buf = strconv.AppendUint(buf, uint64(m.MethodID), 10)
	} else {
		buf = append(buf, m.MethodName...)
	}
	return string(buf)
}

// A Fulfiller is a placeholder for a promised value or an error.
// The zero value is an empty placeholder.  A Fulfiller is safe to use
// from multiple goroutines.
type Fulfiller struct {
	once sync.Once
	done chan struct{}

	mu       sync.RWMutex
	resolved bool
	obj      Object
	err      error
}

// init initializes the fulfiller.  It is idempotent.
func (f *Fulfiller) init() {
	f.once.Do(func() {
		f.done = make(chan struct{})
	})
}

// get blocks until Fulfill or Reject is called and returns the result.
func (f *Fulfiller) get() (obj Object, err error) {
	f.init()
	<-f.done
	obj, err, _ = f.nowaitGet()
	return
}

// nowaitGet returns the current state of the fulfiller.
func (f *Fulfiller) nowaitGet() (obj Object, err error, resolved bool) {
	f.mu.RLock()
	obj, err, resolved = f.obj, f.err, f.resolved
	f.mu.RUnlock()
	return
}

// Fulfill resolves f with a successful value.
// This method is a no-op if f has already been resolved.
func (f *Fulfiller) Fulfill(obj Object) {
	f.init()
	f.mu.Lock()
	if !f.resolved {
		f.obj = obj
		close(f.done)
	}
	f.mu.Unlock()
}

// Reject resolves f with an error.
// This method is a no-op if f has already been resolved.
func (f *Fulfiller) Reject(err error) {
	f.init()
	f.mu.Lock()
	if !f.resolved {
		f.err = err
		close(f.done)
	}
	f.mu.Unlock()
}

// Promise returns a promise that resolves when f is resolved.
func (f *Fulfiller) Promise() Promise {
	obj, err, resolved := f.nowaitGet()
	if !resolved {
		return &fulfillerPromise{f: f, path: fulfillerPath{noop: true}}
	}
	if err != nil {
		return ErrorPromise(err)
	}
	return ImmediatePromise(obj.ToStruct())
}

// Client returns a client that makes calls to the fulfilled interface.
// If f is already resolved, then the fulfilled client is returned.
func (f *Fulfiller) Client() Client {
	obj, err, resolved := f.nowaitGet()
	if !resolved {
		return &fulfillerClient{f: f, path: fulfillerPath{noop: true}}
	}
	if err != nil {
		return ErrorClient(err)
	}
	return obj.ToInterface().Client()
}

type fulfillerPromise struct {
	path fulfillerPath
	f    *Fulfiller
}

func (fp *fulfillerPromise) Get() (Struct, error) {
	p, err := fp.f.get()
	if err != nil {
		return Struct{}, err
	}
	p = fp.path.traverse(p.ToStruct())
	return p.ToStruct(), nil
}

func (fp *fulfillerPromise) GetClient(off int) Client {
	return &fulfillerClient{
		f: fp.f,
		path: fulfillerPath{
			parent: &fp.path,
			field:  off,
		},
	}
}

func (fp *fulfillerPromise) GetPromise(off int) Promise {
	return fp.GetPromiseDefault(off, nil, 0)
}

func (fp *fulfillerPromise) GetPromiseDefault(off int, s *Segment, tagoff int) Promise {
	return &fulfillerPromise{
		f: fp.f,
		path: fulfillerPath{
			parent: &fp.path,
			field:  off,
			defseg: s,
			defoff: tagoff,
		},
	}
}

type fulfillerClient struct {
	path fulfillerPath
	f    *Fulfiller
}

func (fc *fulfillerClient) client(obj Object, err error) (Client, error) {
	if err != nil {
		return nil, err
	}
	client := fc.path.traverse(obj.ToStruct()).ToInterface().Client()
	if client == nil {
		return nil, ErrNullClient
	}
	return client, nil
}

func (fc *fulfillerClient) Call(ctx context.Context, method *Method, params Struct) Promise {
	obj, err, ok := fc.f.nowaitGet()
	if ok {
		client, err := fc.client(obj, err)
		if err != nil {
			return ErrorPromise(err)
		}
		return client.Call(ctx, method, params)
	}
	var results Fulfiller
	go func() {
		client, err := fc.client(fc.f.get())
		if err != nil {
			results.Reject(err)
			return
		}
		s, err := client.Call(ctx, method, params).Get()
		if err != nil {
			results.Reject(err)
			return
		}
		results.Fulfill(Object(s))
	}()
	return results.Promise()
}

func (fc *fulfillerClient) NewCall(ctx context.Context, method *Method, paramsSize ObjectSize, params func(Struct)) Promise {
	obj, err, ok := fc.f.nowaitGet()
	if ok {
		client, err := fc.client(obj, err)
		if err != nil {
			return ErrorPromise(err)
		}
		return client.NewCall(ctx, method, paramsSize, params)
	}
	var results Fulfiller
	go func() {
		client, err := fc.client(fc.f.get())
		if err != nil {
			results.Reject(err)
			return
		}
		s, err := client.NewCall(ctx, method, paramsSize, params).Get()
		if err != nil {
			results.Reject(err)
			return
		}
		results.Fulfill(Object(s))
	}()
	return results.Promise()
}

func (fc *fulfillerClient) Close() error {
	obj, err, ok := fc.f.nowaitGet()
	if !ok {
		return nil
	}
	if err != nil {
		return err
	}
	client, err := fc.client(obj, nil)
	if err != nil {
		return err
	}
	return client.Close()
}

type fulfillerPath struct {
	parent *fulfillerPath

	noop bool

	field  int
	defseg *Segment
	defoff int
}

func (p *fulfillerPath) len() int {
	var n int
	for curr := p; curr != nil; curr = curr.parent {
		if !curr.noop {
			n++
		}
	}
	return n
}

// traverse walks through a struct to get to a value.
func (p *fulfillerPath) traverse(s Struct) Object {
	n := p.len()
	ops := make([]*fulfillerPath, n)
	if n == 0 {
		return Object(s)
	}
	for curr, i := p, n-1; curr != nil; curr = curr.parent {
		if !curr.noop {
			ops[i] = curr
			i--
		}
	}
	for _, op := range ops[:n-1] {
		field := s.GetObject(int(op.field))
		if op.defseg == nil {
			s = field.ToStruct()
		} else {
			s = field.ToStructDefault(p.defseg, p.defoff)
		}
	}
	return s.GetObject(int(ops[n-1].field))
}

type immediatePromise Struct

// ImmediatePromise returns a Promise that accesses s.
func ImmediatePromise(s Struct) Promise {
	return immediatePromise(s)
}

func (ip immediatePromise) Get() (Struct, error) {
	return Struct(ip), nil
}

func (ip immediatePromise) GetClient(off int) Client {
	return Struct(ip).GetObject(off).ToInterface().Client()
}

func (ip immediatePromise) GetPromise(off int) Promise {
	return immediatePromise(Struct(ip).GetObject(off).ToStruct())
}

func (ip immediatePromise) GetPromiseDefault(off int, s *Segment, tagoff int) Promise {
	return immediatePromise(Struct(ip).GetObject(off).ToStructDefault(s, tagoff))
}

type errorPromise struct {
	e error
}

// ErrorPromise returns a Promise that always returns error e.
func ErrorPromise(e error) Promise {
	return errorPromise{e}
}

func (ep errorPromise) Get() (Struct, error) {
	return Struct{}, ep.e
}

func (ep errorPromise) GetClient(off int) Client {
	return ErrorClient(ep.e)
}

func (ep errorPromise) GetPromise(off int) Promise {
	return ep
}

func (ep errorPromise) GetPromiseDefault(off int, s *Segment, tagoff int) Promise {
	return ep
}

type errorClient struct {
	e error
}

// ErrorClient returns a Client that always returns error e.
func ErrorClient(e error) Client {
	return errorClient{e}
}

func (ec errorClient) NewCall(ctx context.Context, method *Method, paramsSize ObjectSize, params func(Struct)) Promise {
	return ErrorPromise(ec.e)
}

func (ec errorClient) Call(ctx context.Context, method *Method, params Struct) Promise {
	return ErrorPromise(ec.e)
}

func (ec errorClient) Close() error {
	return nil
}
