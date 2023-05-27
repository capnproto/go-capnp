# Overview

So far, we've used Cap'n Proto as a faster version of Protocol Buffers.  This is already a huge improvement, but the real power of Cap'n Proto lies in its high-performance [RPC protocol][rpc].  Cap'n Proto RPC builds on the data serialization in much the same way as gRPC is built on Protocol Buffers.

But unlike gRPC, Cap'n Proto offers a rich [object capability model](https://en.wikipedia.org/wiki/Object-capability_model) and protocol-level optimizations, including network path-shortening and promise pipelining (shown below).

<p align="center">
  <img src="https://capnproto.org/images/time-travel.png" alt="Cap'n Proto Promise Pipelining"/>
</p>


## Object Capabilities

At its core, Cap'n Proto RPC is a **distributed object protocol**.  It allows you to call methods on objects residing in remote hosts.  To do this, you first obtain a reference to the remote object, called a **capability**.  Method calls on the capability are translated into RPC calls against the object it points to.

Capabilities are first-class objects, and this means you can:

1. embed them in a Cap'n Proto `struct`,
2. store them in a `List` type,
3. pass them as arguments to RPC methods, and
4. return them from RPC calls.

This is a huge improvement over typical RPC protocols.  JSON-RPC, Go's `net/rpc`, gRPC and Thrift only allow you to address global URLs or singleton objects that are registered with the server.  In contrast, Cap'n Proto RPC allows you to dynamically create new objects at runtime, and share them over the network.  In other words, you can do object-oriented programming (OOP) over the network!

This pattern of [Object Capabilities](http://habitatchronicles.com/2017/05/what-are-capabilities/) provides a powerful framework for writing secure, performant protocols.  We'll explore this paradigm in more detail in a [later chapter](RPC-and-Capability-Oriented-Design.md).

For now, let's skip over the theory and proceed by example.

- **First**, you will learn how to declare a capability in your schema.  In this step, we'll also see how `capnp compile` uses this to generate a Go interface.
- **Next**, you'll learn how to implement the capability server, _i.e._ the implementation for the Go interface generated in the previous step.
- **Finally**, you'll learn how to obtain a _capability_ that points to your server, share it with a remote host, and use it to make RPC calls.

## Declare a Capability in your Schema

We will begin by reproducing the `Arith` RPC server from the Go [Standard Library RPC package documentation](https://pkg.go.dev/net/rpc).

We begin by defining a capability in the schema file.  Capabilities are declared via the `interface` keyword.  When the `capnp` compiler encounters an `interface` type in the schema, it will generate a capability type of the same name, along with a Go interface for the capability server.  We will examine both in more detail in a moment.

For now, open a new file called `arith.capnp`, and copy/paste the following schema:

```capnp
using Go = import "/go.capnp";

@0xf454c62f08bc504b;

$Go.package("arith");
$Go.import("arith");

# Declare the Arith capability, which provides multiplication and division.
interface Arith {
	multiply @0 (a :Int64, b :Int64) -> (product :Int64);
	divide   @1 (num :Int64, denom :Int64) -> (quo :Int64, rem :Int64);
}
```

Now, compile the schema as before:
```bash
capnp compile -I /path/to/go-capnp/std -ogo arith.capnp
```

You should take a moment to inspect the generated types in `arith.capnp.go`.  For interface declarations in your schema, the capnp compiler generates several types, summarized in the following tables.

### Server Types

| Name           | Go Type |   Descripton        |
| -------------- |-------------| -------------|
| `Arith_Server` | `interface` | Network-shareable object.  Methods are RPC endpoints.|
| `Arith_<method>` | `struct` | Method call parameters, _e.g._ `Arith_multiply`.<br />Received by `Arith_Server` method when handling RPC call.|

### Client Types

| Name           | Go Type |   Descripton        |
| -------------- |-------------| -------------|
| `Arith`        | `struct` | The "client" or "capability".  Instances point to a specific `Arith_Server`.<br />Method calls perform RPC against corresponding `Arith_Server`.      |
| `Arith_<method>_Params` | `struct` | Arguments to `<method>`,  _e.g._ `Arith_multiply_Params` |
| `Arith_<method>_Results` | `struct` | The results of an RPC call, _e.g._ `Arith_multiply_Results`. |
| `Arith_<method>_Results_Future` | `struct` | A [promise type](https://en.wikipedia.org/wiki/Futures_and_promises).  Represents in-flight RPC request, _e.g._ `Arith_multiply_Results_Future`.<br />Resolves to `Arith_<method>_Results`.|

## Implement the Server Interface

Now that we have compiled our schema and inspected the generated Go code, let's write an implementation for `Arith_Server`.  Most of the time, you will only write one implementation for your capability.  Note however that because `Arith_Server` is just an ordinary Go interface, you can have multiple implementations in your program.  Common uses for this are to create mock implementations for testing, and to implement restricted or revokable capabilities.  We will explore both patterns in a later chapter.  For now, let's keep it simple.

Let's define our `Arith` server.  In the same directory, create an `arith.go` file and paste the following code:

```go
package arith

import capnp "capnproto.org/go/capnp/v3"

// ArithServer satisfies the Arith_Server interface that was generated
// by the capnp compiler.
type ArithServer struct{}

// Multiply is the concrete implementation of the Multiply method that was
// defined in the schema. Notice that the method signature matches that of
// the Arith_Server interface.
//
// The Arith_multiply struct was generated by the capnp compiler.  You will
// find it in arith.capnp.go
func (Arith) Multiply(ctx context.Context, call Arith_multiply) error {
	res, err := call.AllocResults()  // allocate the results struct
	if err != nil {
		return err
	}

        // Set the result to be the product of the two arguments, A and B,
        // that we received. These are found in the Arith_multiply struct.
	res.SetProduct(call.Args().A() * call.Args().B())
	return nil
}

// Divide is analogous to Multiply.  All capability server methods follow the
// same pattern.
func (Arith) Divide(ctx context.Context, call Arith_divide) error {
	if call.Args().Denom() == 0 {
		return errors.New("divide by zero")
	}

	res, err := call.AllocResults()
	if err != nil {
		return err
	}

	res.SetQuo(call.Args().Num() / call.Args().Denom())
	res.SetRem(call.Args().Num() % call.Args().Denom())
	return nil
}
```

## Share the Capability and Perform RPC

We now have a working RPC server implementation for our schema interface.  Let's begin by starting a server and listening for incoming RPC calls.

The following snippet instantiates an `arith.Arith` server, and exports it over a bidirectional stream.

```go
// Instantiate a local ArithServer.
server := arith.ArithServer{}

// Derive a client capability that points to the server.  Note the
// return type of arith.ServerToClient.  It is of type arith.Arith,
// which is the client capability.  This capability is bound to the
// server instance above; calling client methods will result in RPC
// against the corresponding server method.
//
// The client can be shared over the network.
client := arith.Arith_ServerToClient(server)

// Expose the client over the network.  The 'rwc' parameter can be any
// io.ReadWriteCloser.  In practice, it is almost always a net.Conn.
//
// Note the BootstrapClient option.  This tells the RPC connection to
// immediately make the supplied client -- an arith.Arith, in our case
// -- to the remote endpoint.  The capability that an rpc.Conn exports
// by default is called the "bootstrap capability".
conn := rpc.NewConn(rpc.NewStreamTransport(rwc), &rpc.Options{
	// The BootstrapClient is the RPC interface that will be made available
	// to the remote endpoint by default.  In this case, Arith.
	BootstrapClient: capnp.Client(client),
})
defer conn.Close()

// Block until the connection terminates.
select {
case <-conn.Done():
	return nil
case <-ctx.Done():
	return conn.Close()
}
```

And here's the corresponding client setup and RPC call:

```go
// As before, rwc can be any io.ReadWriteCloser, and will typically be
// a net.Conn.  The rpc.Options can be nil, if you don't want to override
// the defaults.
//
// Here, we expect to receive an arith.Arith from the remote side.  The
// remote side is not expecting a capability in return, however, so we
// don't need to define a bootstrap interface.
//
// This last point bears emphasis:  capnp RPC is fully bidirectional!  Both
// sides of a connection MAY export a boostrap interface, and in such cases,
// the bootstrap interfaces need not be the same!
//
// Again, for the avoidance of doubt:  only the remote side is exporting a
// bootstrap interface in this example.
conn := rpc.NewConn(rpc.NewStreamTransport(rwc), nil)
defer conn.Close()

// Now we resolve the bootstrap interface from the remote ArithServer.
// Thanks to Cap'n Proto's promise pipelining, this function call does
// NOT block.  We can start making RPC calls with 'a' immediately, and
// these will transparently resolve when bootstrapping completes.
//
// The context can be used to time-out or otherwise abort the bootstrap
// call.   It is safe to cancel the context after the first method call
// on 'a' completes.
a := Arith(conn.Bootstrap(ctx))

// Okay! Let's make an RPC call!  Remember:  RPC is performed simply by
// calling a's methods.
//
// There are couple of interesting things to note here:
//  1. We pass a callback function to set parameters on the RPC call.  If the
//     call takes no arguments, you MAY pass nil.
//  2. We return a Future type, representing the in-flight RPC call.  As with
//     the earlier call to Bootstrap, a's methods do not block.  They instead
//     return a future that eventually resolves with the RPC results. We also
//     return a release function, which MUST be called when you're done with
//     the RPC call and its results.
f, release := a.Multiply(ctx, func(ps arith.Arith_multiply_Params) error {
	ps.SetA(2)
	ps.SetB(42)
	return nil
})
defer release()

// You can do other things while the RPC call is in-flight.  Everything
// is asynchronous. For simplicity, we're going to block until the call
// completes.
res, err := f.Struct()
if err != nil {
	return err
}

// Lastly, let's print the result.  Recall that 'product' is the name of
// the return value that we defined in the schema file.
log.Println(res.Product())  // prints 84
```

And that's it!  Let's reiterate the key points about calling RPC methods:

1. For the sake of simplicity, this example uses an in-memory pipe, but you can use TCP connections, Unix pipes, or any other type that implements `io.ReadWriteCloser`.
2. The return type for a client call is a promise, not an immediate value.
It isn't until the `Struct()` method is called on a method that the `client` function blocks on the remote side.

A few additional words on the Future type are in order.  If your RPC method returns another interface type, you can use the Future to immediately make calls against that as-of-yet-unreturned interface.  This relies on a feature of the Cap'n Proto RPC protocol called [promise pipelining][pipelining], the advantage of which is that Cap'n Proto can often optimize away the additional network round-trips when such method calls are chained. This is one of Cap'n Proto's key advantages, which we will use heavily in the next chapter.

## Streaming and Backpressure

Cap'n Proto supports streaming workflows. Unlike other RPC protocols
such as grpc, this can be done without any dedicated "streaming"
construct. Instead, you can define an interface such as:

```capnp
interface ByteStream {
  write @0 (data :Data);
  done @1 ();
}
```

The above is roughly analogous to the `io.WriteCloser` interface. If you
have a `ByteStream` interface, you can write your data into it in
chunks, and then call done() to signal that all data has been
written.

There are however two wrinkles.

The first is flow control. If you naively call write() in a loop, Cap'n
Proto will not by default provide any backpressure, resulting in excess
memory usage and high latency. But waiting for each call in turn results
in low throughput. To solve this, you need to attach a flow limiter,
from the `flowcontrol` package. You can do this with `SetFlowLimiter` on
any capability:

```go
import "capnproto.org/go/capnp/v3/flowcontrol"

// ...

// Limits in-flight data to 2^16 bytes = 64KiB:
client.SetFlowLimiter(flowcontrol.NewFixedLimiter(1 << 16))
```

If too much data is already in-flight, This will cause future rpc calls
to block until some existing ones have returned.

The second wrinkle is dealing with return values: even though
`ByteStream.write()` has no return value, each call will return a future
and ReleaseFunc which must be waited on at some point. You could
accumulate these in a slice and wait on each of them at the end, but
there is a better option: the schema language supports a special
`stream` return type:

```capnp
interface ByteStream {
  write @0 (data :Data) -> stream;
  done @1 ();
}
```

If the return type of a method is `stream`, instead of returning
a future and `ReleaseFunc`, the generated method will just return an
error:

```diff
- func (c ByteStream) Write(ctx context.Context, params func(ByteStream_write_Params) error) (stream.ByteStream_write_Results_Future, capnp.ReleaseFunc) {
+ func (c ByteStream) Write(ctx context.Context, params func(ByteStream_write_Params) error) error {
```

The implementation will take care of waiting on the Future. If the
error is non-nil, it means some prior streaming call failed; this
can be useful for short-circuiting long streaming workflows.

Additionally, each client type has a `WaitStreaming` method, which should
be called at the end of a streaming workflow. A full example might look
like:

```go
for chunk := range chunks {
    err := client.Write(ctx, func(p ByteStream_write_Params) error {
        return p.SetData(chunk)
    })
    if err != nil {
        return err
    }
}

future, release := client.Done(ctx, nil)
defer release()
_, err := future.Struct()
if err != nil {
    return err
}

if err := client.WaitStreaming(); err != nil {
    return err
}
```

[rpc]: https://capnproto.org/rpc.html
[pipelining]: https://capnproto.org/news/2013-12-13-promise-pipelining-capnproto-vs-ice.html

# Next

Now that you've learned the basics of Cap'n Proto RPC, you are ready to
[learn more about object capabilities and advanced RPC](RPC-and-Capability-Oriented-Design.md).
