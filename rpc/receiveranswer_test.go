package rpc

import (
	"context"
	"fmt"
	"net"
	"os"
	"syscall"
	"testing"

	//"github.com/stretchr/testify/assert"
	//rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc/internal/testcapnp"
	"capnproto.org/go/capnp/v3/server"
)

// A variant of net.Pipe() that uses the socketpair() syscall, instead of
// using an in-proces transport. This is also buffered, which works around
// #189. TODO: once that issue is fixed, delete this and just use net.Pipe().
func netPipe() (net.Conn, net.Conn) {
	fds, err := syscall.Socketpair(syscall.AF_UNIX, syscall.SOCK_STREAM, 0)
	chkfatal(err)
	mkConn := func(fd int, name string) net.Conn {
		conn, err := net.FileConn(os.NewFile(uintptr(fd), name))
		chkfatal(err)
		return conn
	}
	return mkConn(fds[0], "pipe0"), mkConn(fds[1], "pipe1")
}

type capArgsTest struct {
	Errs chan<- error
}

func (me *capArgsTest) Self(ctx context.Context, p testcapnp.CapArgsTest_self) error {
	res, err := p.AllocResults()
	if err != nil {
		return err
	}
	res.SetSelf(testcapnp.CapArgsTest_ServerToClient(me, nil))
	return nil
}

func (me *capArgsTest) Call(ctx context.Context, p testcapnp.CapArgsTest_call) error {
	defer close(me.Errs)
	cap, err := p.Args().Cap()
	chkfatal(err)
	client := cap.Interface().Client()
	chkfatal(client.Resolve(ctx))
	brand, ok := server.IsServer(client.State().Brand)
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

func chkfatal(err error) {
	if err != nil {
		panic(err)
	}
}

func TestBootstrapReceiverAnswerRpc(t *testing.T) {
	t.Parallel()

	cClient, cServer := netPipe()
	defer cClient.Close()
	defer cServer.Close()

	errChan := make(chan error)
	srv := &capArgsTest{Errs: errChan}

	// start server:
	serverConn := NewConn(
		NewStreamTransport(cServer),
		&Options{
			BootstrapClient: testcapnp.CapArgsTest_ServerToClient(srv, nil).Client,
		},
	)
	defer serverConn.Close()

	clientConn := NewConn(NewStreamTransport(cClient), nil)
	defer clientConn.Close()

	ctx := context.Background()
	c := testcapnp.CapArgsTest{Client: clientConn.Bootstrap(ctx)}

	res, rel := c.Call(ctx, func(p testcapnp.CapArgsTest_call_Params) error {
		capId := p.Message().AddCap(c.Client.AddRef())
		p.SetCap(capnp.NewInterface(p.Segment(), capId).ToPtr())
		return nil
	})
	defer rel()
	c.Release()

	_, err := res.Struct()
	chkfatal(err)

	for err := range errChan {
		t.Errorf("Error: %v", err)
	}
}

func TestCallReceiverAnswerRpc(t *testing.T) {
	t.Parallel()

	cClient, cServer := netPipe()
	defer cClient.Close()
	defer cServer.Close()

	errChan := make(chan error)
	srv := &capArgsTest{Errs: errChan}

	// start server:
	serverConn := NewConn(
		NewStreamTransport(cServer),
		&Options{
			BootstrapClient: testcapnp.CapArgsTest_ServerToClient(srv, nil).Client,
		},
	)
	defer serverConn.Close()

	clientConn := NewConn(NewStreamTransport(cClient), nil)
	defer clientConn.Close()

	ctx := context.Background()
	bs := testcapnp.CapArgsTest{Client: clientConn.Bootstrap(ctx)}
	defer bs.Release()

	selfRes, rel := bs.Self(ctx, nil)
	defer rel()
	self := selfRes.Self()
	callRes, rel := self.Call(ctx, func(p testcapnp.CapArgsTest_call_Params) error {
		capId := p.Message().AddCap(self.Client.AddRef())
		p.SetCap(capnp.NewInterface(p.Segment(), capId).ToPtr())
		return nil
	})
	self.Release()
	defer rel()

	_, err := selfRes.Struct()
	chkfatal(err)
	_, err = callRes.Struct()
	chkfatal(err)

	for err = range errChan {
		t.Errorf("Error: %v", err)
	}
}

func TestBootstrapReceiverAnswer(t *testing.T) {
	t.Parallel()

	cClient, cServer := netPipe()
	defer cClient.Close()
	defer cServer.Close()

	errChan := make(chan error)
	srv := &capArgsTest{Errs: errChan}

	conn := NewConn(
		NewStreamTransport(cServer),
		&Options{
			BootstrapClient: testcapnp.CapArgsTest_ServerToClient(srv, nil).Client,
		},
	)
	defer conn.Close()
	trans := NewStreamTransport(cClient)

	ctx := context.Background()

	msg, send, release, err := trans.NewMessage(ctx)
	chkfatal(err)

	bs, err := msg.NewBootstrap()
	chkfatal(err)
	bs.SetQuestionId(0)
	send()
	release()

	msg, send, release, err = trans.NewMessage(ctx)
	chkfatal(err)

	// bootstrap.call(cap = bootstrap)
	call, err := msg.NewCall()
	chkfatal(err)
	call.SetQuestionId(1)
	tgt, err := call.NewTarget()
	chkfatal(err)
	pa, err := tgt.NewPromisedAnswer()
	chkfatal(err)
	pa.SetQuestionId(0)
	// Can leave off transform, since the root of the response is the
	// bootstrap capability.
	call.SetInterfaceId(testcapnp.CapArgsTest_TypeID)
	call.SetMethodId(0)
	params, err := call.NewParams()
	chkfatal(err)
	capTable, err := params.NewCapTable(1)
	chkfatal(err)
	capDesc := capTable.At(0)
	ra, err := capDesc.NewReceiverAnswer()
	chkfatal(err)
	ra.SetQuestionId(0)
	seg := params.Struct.Segment()
	argStruct, err := capnp.NewStruct(seg, capnp.ObjectSize{PointerCount: 1})
	chkfatal(err)
	argStruct.SetPtr(0, capnp.NewInterface(seg, 0).ToPtr())
	params.SetContent(argStruct.ToPtr())
	send()
	release()

	for err = range errChan {
		t.Errorf("Error: %v", err)
	}
}

func TestCallReceiverAnswer(t *testing.T) {
	t.Parallel()

	cClient, cServer := netPipe()
	defer cClient.Close()
	defer cServer.Close()

	errChan := make(chan error)
	srv := &capArgsTest{Errs: errChan}

	conn := NewConn(
		NewStreamTransport(cServer),
		&Options{
			BootstrapClient: testcapnp.CapArgsTest_ServerToClient(srv, nil).Client,
		},
	)
	defer conn.Close()
	trans := NewStreamTransport(cClient)

	ctx := context.Background()

	msg, send, release, err := trans.NewMessage(ctx)
	chkfatal(err)

	bs, err := msg.NewBootstrap()
	chkfatal(err)
	bs.SetQuestionId(0)
	send()
	release()

	msg, send, release, err = trans.NewMessage(ctx)
	chkfatal(err)

	// qid1 = bootstrap.self()
	call, err := msg.NewCall()
	chkfatal(err)
	call.SetQuestionId(1)
	tgt, err := call.NewTarget()
	chkfatal(err)
	pa, err := tgt.NewPromisedAnswer()
	chkfatal(err)
	pa.SetQuestionId(0)
	call.SetInterfaceId(testcapnp.CapArgsTest_TypeID)
	call.SetMethodId(1)
	send()
	release()

	msg, send, release, err = trans.NewMessage(ctx)
	chkfatal(err)

	// qid1.self.call(cap = qid1.self)
	call, err = msg.NewCall()
	chkfatal(err)
	call.SetQuestionId(2)
	tgt, err = call.NewTarget()
	chkfatal(err)
	pa, err = tgt.NewPromisedAnswer()
	chkfatal(err)
	pa.SetQuestionId(1)
	transform, err := pa.NewTransform(1)
	chkfatal(err)
	transform.At(0).SetGetPointerField(0)
	call.SetInterfaceId(testcapnp.CapArgsTest_TypeID)
	call.SetMethodId(0)
	params, err := call.NewParams()
	chkfatal(err)
	capTable, err := params.NewCapTable(1)
	chkfatal(err)
	capDesc := capTable.At(0)
	ra, err := capDesc.NewReceiverAnswer()
	chkfatal(err)
	transform.At(0).SetGetPointerField(0)
	ra.SetQuestionId(1)
	transform, err = ra.NewTransform(1)
	chkfatal(err)
	transform.At(0).SetGetPointerField(0)
	seg := params.Struct.Segment()
	argStruct, err := capnp.NewStruct(seg, capnp.ObjectSize{PointerCount: 1})
	chkfatal(err)
	argStruct.SetPtr(0, capnp.NewInterface(seg, 0).ToPtr())
	params.SetContent(argStruct.ToPtr())
	send()
	release()

	for err = range errChan {
		t.Errorf("Error: %v", err)
	}
}
