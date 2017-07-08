package server_test

import (
	"fmt"

	"zombiezen.com/go/capnproto2"
	"zombiezen.com/go/capnproto2/server"
)

func ExampleIsServer() {
	x := int(42)
	c := capnp.NewClient(server.New([]server.Method{}, x, nil, nil))
	if brand, ok := server.IsServer(c.Brand()); ok {
		fmt.Println("Client is a server, got brand:", brand)
	}
	// Output:
	// Client is a server, got brand: 42
}
