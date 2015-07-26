package rpc_test

import (
	"testing"

	"golang.org/x/net/context"
	"zombiezen.com/go/capnproto"
	"zombiezen.com/go/capnproto/rpc"
	"zombiezen.com/go/capnproto/rpc/internal/testcapnp"
)

func TestCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	p, q := newPipe()
	c := rpc.NewConn(p)
	notify := make(chan struct{})
	hanger := testcapnp.Hanger_ServerToClient(Hanger{notify: notify})
	d := rpc.NewConn(q, rpc.MainInterface(hanger.GenericClient()))
	defer d.Wait()
	defer c.Close()
	client := testcapnp.NewHanger(c.Bootstrap(ctx))

	subctx, subcancel := context.WithCancel(ctx)
	promise := client.Hang(subctx, func(r testcapnp.Hanger_hang_Params) {})
	<-notify
	subcancel()
	_, err := promise.Get()
	<-notify // test will deadlock if cancel not delivered

	if err != context.Canceled {
		t.Error("promise.Get() error: %v; want %v", err, context.Canceled)
	}
}

type Hanger struct {
	notify chan struct{}
}

func (h Hanger) Hang(c context.Context, opts capnp.CallOptions, p testcapnp.Hanger_hang_Params, r testcapnp.Hanger_hang_Results) error {
	h.notify <- struct{}{}
	<-c.Done()
	close(h.notify)
	return nil
}
