// +build go1.8

package capnp

import (
	"net"

	"zombiezen.com/go/capnproto2/internal/packed"
)

func (e *Encoder) write(bufs net.Buffers) error {
	if e.packed {
		for _, b := range bufs {
			e.packbuf = packed.Pack(e.packbuf[:0], b)
			if _, err := e.w.Write(e.packbuf); err != nil {
				return err
			}
		}
		return nil
	}
	_, err := bufs.WriteTo(e.w)
	return err
}
