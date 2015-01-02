package capn

import (
	"errors"
	"strconv"

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
