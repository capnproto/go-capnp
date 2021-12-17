package main

import (
	"crypto/sha1"
	"hash"
	"io"
	"fmt"
	"net"
	
	"hashes"
	"golang.org/x/net/context"
	"capnproto.org/go/capnp/v3/rpc"
)

// hashFactory is a local implementation of HashFactory.
type hashFactory struct{}

func (hf hashFactory) NewSha1(_ context.Context, call hashes.HashFactory_newSha1) error {
	// Create a new locally implemented Hash capability.
	hs := hashes.Hash_ServerToClient(hashServer{sha1.New()},nil)
	// Notice that methods can return other interfaces.
	res, err := call.AllocResults()
	if err != nil {
		return err
	}
	
	return res.SetHash(hs)
}

// hashServer is a local implementation of Hash.
type hashServer struct {
    h hash.Hash
}

func (hs hashServer) Write(_ context.Context, call hashes.Hash_write) error {
    data, err := call.Args().Data()
    if err != nil {
        return err
    }

    _, err = hs.h.Write(data)
    return err
}

func (hs hashServer) Sum(_ context.Context, call hashes.Hash_sum) error {
    res, err := call.AllocResults()
    if err != nil {
        return err
    }

    b := hs.h.Sum(nil)
    return res.SetHash(b)
}

func serveHash(ctx context.Context, rwc io.ReadWriteCloser) error {
    // Create a new locally implemented HashFactory.
    main := hashes.HashFactory_ServerToClient(hashFactory{}, nil)

    // Listen for calls, using the HashFactory as the bootstrap interface.
    conn := rpc.NewConn(rpc.NewStreamTransport(rwc), &rpc.Options{
        BootstrapClient:  main.Client,
    })
    defer conn.Close()


    // Wait for connection to abort.
    select {
    case <-conn.Done():
        return nil
    case <-ctx.Done():
        return conn.Close()
    }
}

func client(ctx context.Context, rwc io.ReadWriteCloser) error {
	// Create a connection that we can use to get the HashFactory.
	conn := rpc.NewConn(rpc.NewStreamTransport(rwc), nil) // nil sets default options
	defer conn.Close()
	
	// Get the "bootstrap" interface.  This is the capability set with
	// rpc.MainInterface on the remote side.
	hf := hashes.HashFactory{Client: conn.Bootstrap(ctx)}
	
	// Now we can call methods on hf, and they will be sent over c.
	// The NewSha1 method does not have any parameters we can set, so we
	// pass a nil function.
	f, free := hf.NewSha1(ctx, nil)

	defer free()
	
	// 'NewSha1' returns a future, which allows us to pipeline calls to
	// returned values before they are actually delivered.  Here, we issue
	// calls to an as-of-yet-unresolved Sha1 instance.
	s := f.Hash()

	// s refers to a remote Hash.  Method calls are delivered in order.
	f, free = s.Write(ctx, func(p hashes.Hash_write_Params) error {
		return p.SetData([]byte("Hello, "))
	})
	defer free()
	
	f, free = s.Write(ctx, func(p hashes.Hash_write_Params) error {
		return p.SetData([]byte("World!"))
	})
	defer free()
	
	// Get the sum, waiting for the result.
	f, free = s.Sum(ctx, nil)

	defer free()

	result, err := f.Struct()
	
	if err != nil {
		return err
	}
	
	// Display the result.
	sha1Val, err := result.Hash()
	if err != nil {
		return err
	}
	fmt.Printf("sha1: %x\n", sha1Val)
	return nil
}

func main() {
    ctx := context.Background()
    c1, c2 := net.Pipe()
    go serveHash(ctx, c1)
    client(ctx, c2)
}
