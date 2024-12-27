package rpc

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc/internal/testcapnp"
)

func TestConnection_BaseContext(t *testing.T) {
	t.Run("background context", func(t *testing.T) {
		client, server := net.Pipe()
		doneCh := make(chan struct{}, 1)

		go func() {
			bootstrapClient := testcapnp.StreamTest_ServerToClient(slowStreamTestServer{})
			conn := NewConn(NewStreamTransport(server), &Options{
				BootstrapClient: capnp.Client(bootstrapClient),
			})
			defer conn.Close()

			<-conn.Done()
			close(doneCh)
		}()

		func() {
			conn := NewConn(NewStreamTransport(client), nil)
			defer conn.Close()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			client := testcapnp.StreamTest(conn.Bootstrap(ctx))
			defer client.Release()

			err := client.Push(ctx, func(st testcapnp.StreamTest_push_Params) error {
				return st.SetData(make([]byte, 1))
			})

			require.NoError(t, err)
		}()

		<-doneCh
	})

	t.Run("external context", func(t *testing.T) {
		_, server := net.Pipe()
		doneCh := make(chan struct{}, 1)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			bootstrapClient := testcapnp.StreamTest_ServerToClient(slowStreamTestServer{})
			conn := NewConn(NewStreamTransport(server), &Options{
				BootstrapClient: capnp.Client(bootstrapClient),
				BaseContext: func() context.Context {
					return ctx
				},
			})
			defer conn.Close()

			<-conn.Done()
			close(doneCh)
		}()

		cancel()
		<-doneCh
	})
}
