package main

import (
	"os"

	"books"
	"capnproto.org/go/capnp/v3"
)

func main() {
	// Make a brand new empty message.  A Message allocates Cap'n Proto structs.
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		panic(err)
	}

	// Create a new Book struct.  Every message must have a root struct.
	book, err := books.NewRootBook(seg)
	if err != nil {
		panic(err)
	}
	_ = book.SetTitle("War and Peace")
	book.SetPageCount(1440)

	// Write the message to stdout.
	err = capnp.NewEncoder(os.Stdout).Encode(msg)
	if err != nil {
		panic(err)
	}
}
