package testnetwork

import (
	"testing"

	"github.com/stretchr/testify/require"

	rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"
)

func TestBasicConnect(t *testing.T) {
	j := NewJoiner()
	n1 := j.Join(nil)
	n2 := j.Join(nil)

	trans1, err := n1.DialTransport(n2.LocalID())
	require.NoError(t, err)
	trans2, err := n2.DialTransport(n1.LocalID())
	require.NoError(t, err)

	sendMsg, err := trans1.NewMessage()
	require.NoError(t, err)
	_, err = sendMsg.Message().NewCall()
	require.NoError(t, err)
	require.NoError(t, sendMsg.Send())

	recvMsg, err := trans2.RecvMessage()
	require.NoError(t, err)
	require.Equal(t, rpccp.Message_Which_call, recvMsg.Message().Which())
}
