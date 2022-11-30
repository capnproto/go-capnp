package rpc_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"capnproto.org/go/capnp/v3/rpc"
)

func TestServe(t *testing.T) {
	t.Parallel()
	t.Log("Opening listener")
	lis, err := net.Listen("tcp", ":0")
	defer lis.Close()
	assert.NoError(t, err)

	errChannel := make(chan error)
	go func() {
		err2 := rpc.Serve(lis, nil)
		t.Log("Serve has ended")
		errChannel <- err2
	}()

	time.Sleep(time.Second)

	t.Log("Closing listener")
	err = lis.Close()
	assert.NoError(t, err)
	select {
	case <-time.After(time.Second * 2):
		t.Fail()
	case err = <-errChannel:
		assert.ErrorIs(t, err, net.ErrClosed)
	}
}

func TestListenAndServe(t *testing.T) {
	var err error
	t.Parallel()
	ctx, cancelFunc := context.WithCancel(context.Background())
	errChannel := make(chan error)

	go func() {
		t.Log("Starting ListenAndServe")
		err2 := rpc.ListenAndServe(ctx, "tcp", ":0", nil)
		errChannel <- err2
		t.Log("ListenAndServe has ended")
	}()

	time.Sleep(time.Second)
	t.Log("Closing context")

	cancelFunc()
	select {
	case <-time.After(time.Second * 2):
		t.Error("Cancelling context didn't end the listener")
	case err = <-errChannel:
		assert.ErrorIs(t, err, net.ErrClosed)
	}
}
