package rpc_test

import (
	"sync"
	"testing"

	"golang.org/x/net/context"
	"zombiezen.com/go/capnproto"
	"zombiezen.com/go/capnproto/rpc"
	"zombiezen.com/go/capnproto/rpc/internal/testcapnp"
)

func TestRelease(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	p, q := newPipe()
	c := rpc.NewConn(p)
	defer c.Close()
	hf := new(HandleFactory)
	d := rpc.NewConn(q, rpc.MainInterface(testcapnp.HandleFactory_ServerToClient(hf).GenericClient()))
	defer d.Close()
	client := testcapnp.NewHandleFactory(c.Bootstrap(ctx))
	r1, err := client.NewHandle(ctx, func(r testcapnp.HandleFactory_newHandle_Params) {}).Get()
	if err != nil {
		t.Fatal("NewHandle #1:", err)
	}
	handle1 := r1.Handle()
	if n := hf.numHandles(); n != 1 {
		t.Fatalf("numHandles = %d; want 1", n)
	}

	if err := handle1.GenericClient().Close(); err != nil {
		t.Error("handle1.Close():", err)
	}

	if n := hf.numHandles(); n != 0 {
		t.Errorf("numHandles = %d; want 0", n)
	}
}

type Handle struct {
	f *HandleFactory
}

func (h Handle) Close() error {
	h.f.mu.Lock()
	h.f.n--
	h.f.mu.Unlock()
	return nil
}

type HandleFactory struct {
	n  int
	mu sync.Mutex
}

func (hf *HandleFactory) NewHandle(
	c context.Context,
	opts capnp.CallOptions,
	p testcapnp.HandleFactory_newHandle_Params,
	r testcapnp.HandleFactory_newHandle_Results) error {
	hf.mu.Lock()
	hf.n++
	hf.mu.Unlock()
	r.SetHandle(testcapnp.Handle_ServerToClient(Handle{f: hf}))
	return nil
}

func (hf *HandleFactory) numHandles() int {
	hf.mu.Lock()
	n := hf.n
	hf.mu.Unlock()
	return n
}
