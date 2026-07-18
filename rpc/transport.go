package rpc

// Re-export things from the transport package

import (
	"io"

	"capnproto.org/go/capnp/v3/rpc/transport"
)

type Codec = transport.Codec
type Transport = transport.Transport
type NewTransportFunc func(io.ReadWriteCloser) Transport
type StreamTransportOptions = transport.StreamTransportOptions

// NewStreamTransport is an alias for as transport.NewStream
func NewStreamTransport(rwc io.ReadWriteCloser) Transport {
	return transport.NewStream(rwc)
}

// NewStreamTransportWithOptions creates a stream transport with the given
// options.
func NewStreamTransportWithOptions(rwc io.ReadWriteCloser, opts StreamTransportOptions) Transport {
	return transport.NewStreamWithOptions(rwc, opts)
}

// NewPackedStreamTransport is an alias for as transport.NewPackedStream
func NewPackedStreamTransport(rwc io.ReadWriteCloser) Transport {
	return transport.NewPackedStream(rwc)
}

// NewPackedStreamTransportWithOptions creates a packed stream transport with
// the given options.
func NewPackedStreamTransportWithOptions(rwc io.ReadWriteCloser, opts StreamTransportOptions) Transport {
	return transport.NewPackedStreamWithOptions(rwc, opts)
}

// NewTransport is an alias for as transport.New
func NewTransport(codec Codec) Transport {
	return transport.New(codec)
}
