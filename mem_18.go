// +build go1.8

package capnp

import "net"

func (e *Encoder) write(bufs net.Buffers) error {
	_, err := bufs.WriteTo(e.w)
	return err
}
