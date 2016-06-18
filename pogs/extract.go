package pogs

import (
	"zombiezen.com/go/capnproto2"
)

// Extract copies s into val, a pointer to a Go struct.
func Extract(val interface{}, typeID uint64, s capnp.Struct) error {
	return nil
}
