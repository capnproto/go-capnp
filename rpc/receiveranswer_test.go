package rpc

import (
	"context"
	"net"
	"testing"

	//"github.com/stretchr/testify/assert"
	//rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc/internal/testcapnp"
	"capnproto.org/go/capnp/v3/server"
)

func TestCallBootstrapReceiverAnswer(t *testing.T) {
	cClient, cServer := net.Pipe()

	conn := NewConn(
		NewStreamTransport(cServer),
		&Options{
			BootstrapClient: testcapnp.PingPong_ServerToClient(nil, &server.Policy{
				MaxConcurrentCalls: 10,
			}).Client,
		})
	defer conn.Close()
	trans := NewStreamTransport(cClient)

	ctx := context.Background()

	chkfatal := func(err error) {
		if err != nil {
			panic(err)
		}
	}

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
	params.SetContent(capnp.NewInterface(params.Struct.Segment(), 0).ToPtr())
	send()
	release()
}
