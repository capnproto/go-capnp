package rpc

import (
	"io"

	capnp "zombiezen.com/go/capnproto2"
)

// TODO(performance):  maintain (global?) sync.Pools for (un)packed encoders & decoders.

// Codec is a factory for encoders and decoders.
type Codec interface {
	Encode(*capnp.Message) ([]byte, error)
	Decode(io.Reader) (*capnp.Message, error)
}

// UnpackedCodec is a standard (unpacked) encoder/decoder pair.
type UnpackedCodec struct{}

// Encode unpacked.
func (c UnpackedCodec) Encode(msg *capnp.Message) ([]byte, error) {
	return msg.Marshal()
}

// Decode unpacked.
func (c UnpackedCodec) Decode(r io.Reader) (*capnp.Message, error) {
	return capnp.NewDecoder(r).Decode()
}

// PackedCodec is a packed encoder/decoder pair.
type PackedCodec struct{}

// Encode packed.
func (c PackedCodec) Encode(msg *capnp.Message) ([]byte, error) {
	return msg.MarshalPacked()
}

// Decode unpacked.
func (c PackedCodec) Decode(r io.Reader) (*capnp.Message, error) {
	return capnp.NewPackedDecoder(r).Decode()
}
