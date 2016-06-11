// Package schemas provides a program-wide registry for Cap'n Proto reflection data.
package schemas

import (
	"compress/zlib"
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
	"sync"

	"zombiezen.com/go/capnproto2/internal/packed"
)

var registry = make(map[uint64]*record)

type record struct {
	load sync.Once
	c    string
	d    []byte
}

// Register is called by generated code to associate a blob of zlib-
// compressed, packed Cap'n Proto data for a CodeGeneratorRequest with
// the IDs it contains.  It should only be called during init().
func Register(data string, ids ...uint64) {
	r := &record{c: data}
	for _, id := range ids {
		if _, dup := registry[id]; dup {
			panic(errors.New("schemas: registered ID @0x" + strconv.FormatUint(id, 16) + " twice"))
		}
		registry[id] = r
	}
}

// Find returns the CodeGeneratorRequest message for the given ID,
// suitable for capnp.Unmarshal, or nil if the ID was not found.
// It is safe to call Find from multiple goroutines, so the returned
// byte slice should not be modified.  However, it is not safe to
// call Find concurrently with Register.
func Find(id uint64) []byte {
	r := registry[id]
	if r == nil {
		return nil
	}
	r.load.Do(func() {
		z, err := zlib.NewReader(strings.NewReader(r.c))
		if err != nil {
			panic(decompressError(id, err))
		}
		defer z.Close()
		p := packed.NewReader(z)
		r.d, err = ioutil.ReadAll(p)
		if err != nil {
			panic(decompressError(id, err))
		}
	})
	return r.d
}

func decompressError(id uint64, e error) error {
	return errors.New("schemas: decompressing schema for @0x" + strconv.FormatUint(id, 16) + ": " + e.Error())
}
