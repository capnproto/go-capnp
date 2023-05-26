package server_test

import (
	"fmt"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/server"
)

func ExampleIsServer() {
	x := int(42)
	c := capnp.NewClient(server.New([]server.Method{}, x, nil))
	snapshot := c.Snapshot()
	brand := snapshot.Brand()
	snapshot.Release()
	if brand, ok := server.IsServer(brand); ok {
		fmt.Println("Client is a server, got brand:", brand)
	}
	// Output:
	// Client is a server, got brand: 42
}
