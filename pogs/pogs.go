// Package pogs provides functions to convert Cap'n Proto messages to and from Go structs.
package pogs // import "zombiezen.com/go/capnproto2/pogs"

import (
	"zombiezen.com/go/capnproto2"
)

// Extract copies s into val, a pointer to a Go struct.
func Extract(val interface{}, typeID uint64, s capnp.Struct) error {
	return nil
}

// Insert copies val, a pointer to a Go struct, into s.
func Insert(typeID uint64, s capnp.Struct, val interface{}) error {
	return nil
}
