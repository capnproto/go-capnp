package rpc_test

import (
	"fmt"
	"golang.org/x/net/context"
	"net"

	"zombiezen.com/go/capnproto"
	"zombiezen.com/go/capnproto/rpc"
	"zombiezen.com/go/capnproto/rpc/internal/testcapnp"
)

func Example() {
	// Create an in-memory transport.  In a real application, you would probably
	// use a net.TCPConn (for RPC) or an os.Pipe (for IPC).
	p1, p2 := net.Pipe()
	t1, t2 := rpc.StreamTransport(p1), rpc.StreamTransport(p2)

	// Server-side
	srv := testcapnp.Adder_ServerToClient(AdderServer{})
	serverConn := rpc.NewConn(t1, rpc.MainInterface(srv.GenericClient()))
	defer serverConn.Wait()

	// Client-side
	ctx := context.Background()
	clientConn := rpc.NewConn(t2)
	defer clientConn.Close()
	adderClient := testcapnp.NewAdder(clientConn.Bootstrap(ctx))
	// Every client call returns a promise.  You can make multiple calls
	// concurrently.
	call1 := adderClient.Add(ctx, func(p testcapnp.Adder_add_Params) {
		p.SetA(5)
		p.SetB(2)
	})
	call2 := adderClient.Add(ctx, func(p testcapnp.Adder_add_Params) {
		p.SetA(10)
		p.SetB(20)
	})
	// Calling Get() on a promise waits until it returns.
	result1, err := call1.Get()
	if err != nil {
		fmt.Println("Add #1 failed:", err)
		return
	}
	result2, err := call2.Get()
	if err != nil {
		fmt.Println("Add #2 failed:", err)
		return
	}

	fmt.Println("Results:", result1.Result(), result2.Result())
	// Output:
	// Results: 7 30
}

// An AdderServer is a local implementation of the Adder interface.
type AdderServer struct{}

// Add implements a method
func (AdderServer) Add(call testcapnp.Adder_add) error {
	// Acknowledging the call allows other calls to be made (it returns the Answer
	// to the caller).
	capnp.Ack(call.Options)

	// Parameters are accessed with call.Params.
	a := call.Params.A()
	b := call.Params.B()

	// A result struct is allocated for you at call.Results.
	call.Results.SetResult(a + b)

	return nil
}
