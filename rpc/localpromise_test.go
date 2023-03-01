package rpc_test

import (
	"context"
	"errors"
	"testing"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc/internal/testcapnp"
	"github.com/stretchr/testify/assert"
)

func TestLocalPromiseFulfill(t *testing.T) {
	ctx := context.Background()
	p, r := capnp.NewLocalPromise[testcapnp.PingPong]()
	defer p.Release()

	fut1, rel1 := p.EchoNum(ctx, func(p testcapnp.PingPong_echoNum_Params) error {
		p.SetN(1)
		return nil
	})
	defer rel1()

	fut2, rel2 := p.EchoNum(ctx, func(p testcapnp.PingPong_echoNum_Params) error {
		p.SetN(2)
		return nil
	})
	defer rel2()

	pp := testcapnp.PingPong_ServerToClient(&pingPonger{})
	defer pp.Release()
	r.Fulfill(pp)

	fut3, rel3 := p.EchoNum(ctx, func(p testcapnp.PingPong_echoNum_Params) error {
		p.SetN(3)
		return nil
	})
	defer rel3()

	res1, err := fut1.Struct()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), res1.N())

	res2, err := fut2.Struct()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), res2.N())

	res3, err := fut3.Struct()
	assert.NoError(t, err)
	assert.Equal(t, int64(3), res3.N())
}

func TestLocalPromiseReject(t *testing.T) {
	ctx := context.Background()
	p, r := capnp.NewLocalPromise[testcapnp.PingPong]()
	defer p.Release()

	fut1, rel1 := p.EchoNum(ctx, func(p testcapnp.PingPong_echoNum_Params) error {
		p.SetN(1)
		return nil
	})
	defer rel1()

	fut2, rel2 := p.EchoNum(ctx, func(p testcapnp.PingPong_echoNum_Params) error {
		p.SetN(2)
		return nil
	})
	defer rel2()

	r.Reject(errors.New("Promise rejected"))

	fut3, rel3 := p.EchoNum(ctx, func(p testcapnp.PingPong_echoNum_Params) error {
		p.SetN(3)
		return nil
	})
	defer rel3()

	_, err := fut1.Struct()
	assert.NotNil(t, err)

	_, err = fut2.Struct()
	assert.NotNil(t, err)

	_, err = fut3.Struct()
	assert.NotNil(t, err)
}
