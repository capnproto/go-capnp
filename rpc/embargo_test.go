package rpc_test

import (
	"sync"
	"testing"

	"golang.org/x/net/context"
	"zombiezen.com/go/capnproto"
	"zombiezen.com/go/capnproto/rpc"
	"zombiezen.com/go/capnproto/rpc/internal/logtransport"
	"zombiezen.com/go/capnproto/rpc/internal/testcapnp"
)

func TestEmbargo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	p, q := newPipe()
	if *logMessages {
		p = logtransport.New(nil, p)
	}
	c := rpc.NewConn(p)
	echoSrv := testcapnp.Echoer_ServerToClient(new(Echoer))
	d := rpc.NewConn(q, rpc.MainInterface(echoSrv.GenericClient()))
	defer d.Wait()
	defer c.Close()
	client := testcapnp.NewEchoer(c.Bootstrap(ctx))
	localCap := testcapnp.CallOrder_ServerToClient(new(CallOrder))

	earlyCall := callseq(ctx, client.GenericClient(), 0)
	echo := client.Echo(ctx, func(p testcapnp.Echoer_echo_Params) {
		p.SetCap(localCap)
	})
	pipeline := echo.Cap()
	call0 := callseq(ctx, pipeline.GenericClient(), 0)
	call1 := callseq(ctx, pipeline.GenericClient(), 1)
	_, err := earlyCall.Get()
	if err != nil {
		t.Errorf("earlyCall error: %v", err)
	}
	call2 := callseq(ctx, pipeline.GenericClient(), 2)
	_, err = echo.Get()
	if err != nil {
		t.Errorf("echo.Get() error: %v", err)
	}
	call3 := callseq(ctx, pipeline.GenericClient(), 3)
	call4 := callseq(ctx, pipeline.GenericClient(), 4)
	call5 := callseq(ctx, pipeline.GenericClient(), 5)

	check := func(promise *testcapnp.CallOrder_getCallSequence_Results_Promise, n uint32) {
		r, err := promise.Get()
		if err != nil {
			t.Errorf("call%d error: %v", n, err)
		}
		if r.N() != n {
			t.Errorf("call%d = %d; want %d", n, r.N(), n)
		}
	}
	check(call0, 0)
	check(call1, 1)
	check(call2, 2)
	check(call3, 3)
	check(call4, 4)
	check(call5, 5)
}

func callseq(c context.Context, client capnp.Client, n uint32) *testcapnp.CallOrder_getCallSequence_Results_Promise {
	return testcapnp.NewCallOrder(client).GetCallSequence(c, func(p testcapnp.CallOrder_getCallSequence_Params) {
		p.SetExpected(n)
	})
}

type CallOrder struct {
	mu sync.Mutex
	n  uint32
}

func (co *CallOrder) GetCallSequence(
	c context.Context,
	opts capnp.CallOptions,
	p testcapnp.CallOrder_getCallSequence_Params,
	r testcapnp.CallOrder_getCallSequence_Results) error {
	r.SetN(co.n)
	co.n++
	return nil
}

type Echoer struct {
	CallOrder
}

func (*Echoer) Echo(c context.Context, opts capnp.CallOptions, p testcapnp.Echoer_echo_Params, r testcapnp.Echoer_echo_Results) error {
	r.SetCap(p.Cap())
	return nil
}
