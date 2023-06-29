package rpc_test

import (
	"context"
	"errors"
	"testing"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc/internal/testcapnp"
	"github.com/stretchr/testify/assert"
)

type echoNumOrderChecker struct {
	t       *testing.T
	nextNum int64
}

func (e *echoNumOrderChecker) EchoNum(ctx context.Context, p testcapnp.PingPong_echoNum) error {
	assert.Equal(e.t, e.nextNum, p.Args().N())
	e.nextNum++
	results, err := p.AllocResults()
	if err != nil {
		return err
	}
	results.SetN(p.Args().N())
	return nil
}

func TestLocalPromiseFulfill(t *testing.T) {
	t.Parallel()

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

	pp := testcapnp.PingPong_ServerToClient(&echoNumOrderChecker{
		t:       t,
		nextNum: 1,
	})
	defer pp.Release()
	r.Fulfill(pp)

	fut3, rel3 := p.EchoNum(ctx, func(p testcapnp.PingPong_echoNum_Params) error {
		p.SetN(3)
		return nil
	})
	defer rel3()

	res1, err1 := fut1.Struct()
	res2, err2 := fut2.Struct()
	res3, err3 := fut3.Struct()

	assert.NoError(t, err1)
	assert.Equal(t, int64(1), res1.N())
	assert.NoError(t, err2)
	assert.Equal(t, int64(2), res2.N())
	assert.NoError(t, err3)
	assert.Equal(t, int64(3), res3.N())
}

func echoNum(ctx context.Context, pp testcapnp.PingPong, n int64) (testcapnp.PingPong_echoNum_Results_Future, capnp.ReleaseFunc) {
	return pp.EchoNum(ctx, func(p testcapnp.PingPong_echoNum_Params) error {
		p.SetN(n)
		return nil
	})
}

func TestLocalPromiseReject(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	p, r := capnp.NewLocalPromise[testcapnp.PingPong]()
	defer p.Release()

	fut1, rel1 := echoNum(ctx, p, 1)
	defer rel1()

	fut2, rel2 := echoNum(ctx, p, 2)
	defer rel2()

	r.Reject(errors.New("Promise rejected"))

	fut3, rel3 := echoNum(ctx, p, 3)
	defer rel3()

	_, err := fut1.Struct()
	assert.NotNil(t, err)

	_, err = fut2.Struct()
	assert.NotNil(t, err)

	_, err = fut3.Struct()
	assert.NotNil(t, err)
}
