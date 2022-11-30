package transport_test

import (
	"io"
	"testing"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc/transport"
	"github.com/stretchr/testify/require"
)

func TestPipe(t *testing.T) {
	t.Parallel()

	m, _ := capnp.NewSingleSegmentMessage(nil)

	p1, p2 := transport.NewPipe(1)

	err := p1.Encode(m)
	require.NoError(t, err)

	m2, err := p2.Decode()
	require.NoError(t, err)
	require.NotEqual(t, m, m2, "message should have been copied")

	err = p1.Close()
	require.NoError(t, err)

	err = p1.Encode(m)
	require.ErrorIs(t, err, io.ErrClosedPipe)

	m, err = p2.Decode()
	require.Nil(t, m)
	require.ErrorIs(t, err, io.ErrClosedPipe)
}
