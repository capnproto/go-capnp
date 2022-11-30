package rpc_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"capnproto.org/go/capnp/v3/rpc"
)

//clientHook := {}

func TestServe(t *testing.T) {
	lis, err := net.Listen("tcp", ":8888")
	assert.NoError(t, err)

	errChannel := make(chan error)
	go func() {
		err2 := rpc.Serve(lis, &rpc.Options{})
		errChannel <- err2
	}()
	time.Sleep(time.Second)
	t.Log("Closing listener")
	err = lis.Close()
	assert.NoError(t, err)
	err = <-errChannel
	assert.ErrorIs(t, err, net.ErrClosed)
}

func TestListenAndServe(t *testing.T) {
	var err error
	ctx, cancelFunc := context.WithCancel(context.Background())
	errChannel := make(chan error)
	go func() {
		err2 := rpc.ListenAndServe(ctx, ":8888", &rpc.Options{})
		errChannel <- err2
		t.Log("Serve has ended")
	}()
	time.Sleep(time.Second)
	t.Log("Closing listener")
	cancelFunc()
	err = <-errChannel
	assert.ErrorIs(t, err, net.ErrClosed)
}
