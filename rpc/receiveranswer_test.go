package rpc

import (
	"context"
	"fmt"
	"net"
	"testing"

	//"github.com/stretchr/testify/assert"
	//rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc/internal/testcapnp"
	"capnproto.org/go/capnp/v3/server"
	"zenhack.net/go/util"
)

type capArgsTest struct {
	Errs chan<- error
}

func (me *capArgsTest) Self(ctx context.Context, p testcapnp.CapArgsTest_self) error {
	res, err := p.AllocResults()
	if err != nil {
		return err
	}
	res.SetSelf(testcapnp.CapArgsTest_ServerToClient(me))
	return nil
}

func (me *capArgsTest) Call(ctx context.Context, p testcapnp.CapArgsTest_call) error {
	defer close(me.Errs)
	client := p.Args().Cap()
	util.Chkfatal(client.Resolve(ctx))
	snapshot := client.Snapshot()
	defer snapshot.Release()
	brand, ok := server.IsServer(snapshot.Brand())
	if !ok {
		err := fmt.Errorf("server.IsServer returned !ok")
		me.Errs <- err
		return err
	}
	other := brand.(*capArgsTest)
	if other != me {
		me.Errs <- fmt.Errorf(
			"Passed something other than ourselves: wanted %v but got %v",
			me, other)
	}
	return nil
}

func TestBootstrapReceiverAnswerRpc(t *testing.T) {
	t.Parallel()

	cClient, cServer := net.Pipe()
	defer cClient.Close()
	defer cServer.Close()

	errChan := make(chan error)
	srv := &capArgsTest{Errs: errChan}

	// start server:
	serverConn := NewConn(
		NewStreamTransport(cServer),
		&Options{
			BootstrapClient: capnp.Client(testcapnp.CapArgsTest_ServerToClient(srv)),
		},
	)
	defer serverConn.Close()

	clientConn := NewConn(NewStreamTransport(cClient), nil)
	defer clientConn.Close()

	ctx := context.Background()
	c := testcapnp.CapArgsTest(clientConn.Bootstrap(ctx))

	res, rel := c.Call(ctx, func(p testcapnp.CapArgsTest_call_Params) error {
		return p.SetCap(capnp.Client(c.AddRef()))
	})
	defer rel()
	c.Release()

	_, err := res.Struct()
	util.Chkfatal(err)

	for err := range errChan {
		t.Errorf("Error: %v", err)
	}
}

func TestCallReceiverAnswerRpc(t *testing.T) {
	t.Parallel()

	cClient, cServer := net.Pipe()
	defer cClient.Close()
	defer cServer.Close()

	errChan := make(chan error)
	srv := &capArgsTest{Errs: errChan}

	// start server:
	serverConn := NewConn(
		NewStreamTransport(cServer),
		&Options{
			BootstrapClient: capnp.Client(testcapnp.CapArgsTest_ServerToClient(srv)),
		},
	)
	defer serverConn.Close()

	clientConn := NewConn(NewStreamTransport(cClient), nil)
	defer clientConn.Close()

	ctx := context.Background()
	bs := testcapnp.CapArgsTest(clientConn.Bootstrap(ctx))
	defer bs.Release()

	selfRes, rel := bs.Self(ctx, nil)
	defer rel()
	self := selfRes.Self()
	callRes, rel := self.Call(ctx, func(p testcapnp.CapArgsTest_call_Params) error {
		return p.SetCap(capnp.Client(self.AddRef()))
	})
	defer rel()

	_, err := selfRes.Struct()
	util.Chkfatal(err)
	_, err = callRes.Struct()
	util.Chkfatal(err)

	for err = range errChan {
		t.Errorf("Error: %v", err)
	}
}

func TestBootstrapReceiverAnswer(t *testing.T) {
	t.Parallel()

	cClient, cServer := net.Pipe()
	defer cClient.Close()
	defer cServer.Close()

	errChan := make(chan error)
	srv := &capArgsTest{Errs: errChan}

	conn := NewConn(
		NewStreamTransport(cServer),
		&Options{
			BootstrapClient: capnp.Client(testcapnp.CapArgsTest_ServerToClient(srv)),
		},
	)
	defer conn.Close()
	trans := NewStreamTransport(cClient)

	outMsg, err := trans.NewMessage()
	util.Chkfatal(err)

	bs, err := outMsg.Message().NewBootstrap()
	util.Chkfatal(err)
	bs.SetQuestionId(0)
	outMsg.Send()
	outMsg.Release()

	outMsg, err = trans.NewMessage()
	util.Chkfatal(err)

	// bootstrap.call(cap = bootstrap)
	call, err := outMsg.Message().NewCall()
	util.Chkfatal(err)
	call.SetQuestionId(1)
	tgt, err := call.NewTarget()
	util.Chkfatal(err)
	pa, err := tgt.NewPromisedAnswer()
	util.Chkfatal(err)
	pa.SetQuestionId(0)
	// Can leave off transform, since the root of the response is the
	// bootstrap capability.
	call.SetInterfaceId(testcapnp.CapArgsTest_TypeID)
	call.SetMethodId(0)
	params, err := call.NewParams()
	util.Chkfatal(err)
	capTable, err := params.NewCapTable(1)
	util.Chkfatal(err)
	capDesc := capTable.At(0)
	ra, err := capDesc.NewReceiverAnswer()
	util.Chkfatal(err)
	ra.SetQuestionId(0)
	seg := params.Segment()
	argStruct, err := capnp.NewStruct(seg, capnp.ObjectSize{PointerCount: 1})
	util.Chkfatal(err)
	argStruct.SetPtr(0, capnp.NewInterface(seg, 0).ToPtr())
	params.SetContent(argStruct.ToPtr())
	outMsg.Send()
	outMsg.Release()

	for err = range errChan {
		t.Errorf("Error: %v", err)
	}
}

func TestCallReceiverAnswer(t *testing.T) {
	t.Parallel()

	cClient, cServer := net.Pipe()
	defer cClient.Close()
	defer cServer.Close()

	errChan := make(chan error)
	srv := &capArgsTest{Errs: errChan}

	conn := NewConn(
		NewStreamTransport(cServer),
		&Options{
			BootstrapClient: capnp.Client(testcapnp.CapArgsTest_ServerToClient(srv)),
		},
	)
	defer conn.Close()
	trans := NewStreamTransport(cClient)

	outMsg, err := trans.NewMessage()
	util.Chkfatal(err)

	bs, err := outMsg.Message().NewBootstrap()
	util.Chkfatal(err)
	bs.SetQuestionId(0)
	outMsg.Send()
	outMsg.Release()

	outMsg, err = trans.NewMessage()
	util.Chkfatal(err)

	// qid1 = bootstrap.self()
	call, err := outMsg.Message().NewCall()
	util.Chkfatal(err)
	call.SetQuestionId(1)
	tgt, err := call.NewTarget()
	util.Chkfatal(err)
	pa, err := tgt.NewPromisedAnswer()
	util.Chkfatal(err)
	pa.SetQuestionId(0)
	call.SetInterfaceId(testcapnp.CapArgsTest_TypeID)
	call.SetMethodId(1)
	outMsg.Send()
	outMsg.Release()

	outMsg, err = trans.NewMessage()
	util.Chkfatal(err)

	// qid1.self.call(cap = qid1.self)
	call, err = outMsg.Message().NewCall()
	util.Chkfatal(err)
	call.SetQuestionId(2)
	tgt, err = call.NewTarget()
	util.Chkfatal(err)
	pa, err = tgt.NewPromisedAnswer()
	util.Chkfatal(err)
	pa.SetQuestionId(1)
	transform, err := pa.NewTransform(1)
	util.Chkfatal(err)
	transform.At(0).SetGetPointerField(0)
	call.SetInterfaceId(testcapnp.CapArgsTest_TypeID)
	call.SetMethodId(0)
	params, err := call.NewParams()
	util.Chkfatal(err)
	capTable, err := params.NewCapTable(1)
	util.Chkfatal(err)
	capDesc := capTable.At(0)
	ra, err := capDesc.NewReceiverAnswer()
	util.Chkfatal(err)
	transform.At(0).SetGetPointerField(0)
	ra.SetQuestionId(1)
	transform, err = ra.NewTransform(1)
	util.Chkfatal(err)
	transform.At(0).SetGetPointerField(0)
	seg := params.Segment()
	argStruct, err := capnp.NewStruct(seg, capnp.ObjectSize{PointerCount: 1})
	util.Chkfatal(err)
	argStruct.SetPtr(0, capnp.NewInterface(seg, 0).ToPtr())
	params.SetContent(argStruct.ToPtr())
	outMsg.Send()
	outMsg.Release()

	for err = range errChan {
		t.Errorf("Error: %v", err)
	}
}
