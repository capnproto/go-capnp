package rpc_test

import (
	"testing"

	"golang.org/x/net/context"
	"zombiezen.com/go/capnproto/rpc"
	"zombiezen.com/go/capnproto/rpc/internal/logtransport"
	"zombiezen.com/go/capnproto/rpc/internal/pipetransport"
	"zombiezen.com/go/capnproto/rpc/internal/testcapnp"
	"zombiezen.com/go/capnproto/server"
)

func TestPromisedCapability(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	p, q := pipetransport.New()
	if *logMessages {
		p = logtransport.New(nil, p)
	}
	c := rpc.NewConn(p)
	delay := make(chan struct{})
	echoSrv := testcapnp.Echoer_ServerToClient(&DelayEchoer{delay: delay})
	d := rpc.NewConn(q, rpc.MainInterface(echoSrv.GenericClient()))
	defer d.Wait()
	defer c.Close()
	client := testcapnp.NewEchoer(c.Bootstrap(ctx))

	echo := client.Echo(ctx, func(p testcapnp.Echoer_echo_Params) {
		p.SetCap(testcapnp.NewCallOrder(client.GenericClient()))
	})
	pipeline := echo.Cap()
	call0 := callseq(ctx, pipeline.GenericClient(), 0)
	call1 := callseq(ctx, pipeline.GenericClient(), 1)
	close(delay)

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
}

type DelayEchoer struct {
	Echoer
	delay chan struct{}
}

func (de *DelayEchoer) Echo(call testcapnp.Echoer_echo) error {
	server.Ack(call.Options)
	<-de.delay
	return de.Echoer.Echo(call)
}
