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

func TestCallBootstrapReceiverAnswer(t *testing.T) {
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
