// +build !go1.8

package capnp

import "zombiezen.com/go/capnproto2/internal/packed"

func (e *Encoder) write(bufs [][]byte) error {
	for _, b := range bufs {
		if e.packed {
			e.packbuf = packed.Pack(e.packbuf[:0], b)
			b = e.packbuf
		}
		if _, err := e.w.Write(b); err != nil {
			return err
		}
	}
	return nil
}
