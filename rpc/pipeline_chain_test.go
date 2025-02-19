package rpc_test

import (
	"context"
	"testing"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc/internal/testcapnp"
	"github.com/stretchr/testify/assert"
)

func TestPipelineChaining(t *testing.T) {
	t.Parallel()

	bg := context.Background()

	np := func(cat testcapnp.CapArgsTest_self_Params) error {
		return nil
	}

	ctx, cancel := context.WithCancel(bg)
	cli := testcapnp.CapArgsTest_ServerToClient(&delayingCapReturner{ctx})
	defer cli.Release()

	// Originally, chaining pipelines like this didn't work correctly.

	fut1, rel1 := cli.Self(bg, np)
	defer rel1()
	// This one was fine
	fut2, rel2 := fut1.Self().Self(bg, np)
	defer rel2()
	// This one would be delivered to the wrong place
	// (the same place as the previous one)
	fut3, rel3 := fut2.Self().Self(bg, np)
	defer rel3()
	// This one would segfault
	fut4, rel4 := fut3.Self().Self(bg, np)
	defer rel4()

	cancel()
	_, err := fut4.Struct()

	assert.Nil(t, err)
}

type delayingCapReturner struct {
	context.Context
}

func (c delayingCapReturner) Call(ctx context.Context, call testcapnp.CapArgsTest_call) error {
	return capnp.Unimplemented("yes")
}

func (c delayingCapReturner) Self(ctx context.Context, call testcapnp.CapArgsTest_self) error {
	call.Go()
	select {
	case <-ctx.Done():
	case <-c.Done():
	}

	results, err := call.AllocResults()
	if err != nil {
		return err
	}

	results.SetSelf(testcapnp.CapArgsTest_ServerToClient(selfCapReturner{}))
	return nil
}

type selfCapReturner struct{}

func (c selfCapReturner) Call(ctx context.Context, call testcapnp.CapArgsTest_call) error {
	return capnp.Unimplemented("yes")
}

func (c selfCapReturner) Self(ctx context.Context, call testcapnp.CapArgsTest_self) error {
	results, err := call.AllocResults()
	if err != nil {
		return err
	}

	results.SetSelf(testcapnp.CapArgsTest_ServerToClient(c))
	return nil
}
