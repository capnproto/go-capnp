package server_test

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"

	"capnproto.org/go/capnp/v3"
	air "capnproto.org/go/capnp/v3/internal/aircraftlib"
	"capnproto.org/go/capnp/v3/server"

	"github.com/stretchr/testify/assert"
)

type echoImpl struct{}

func (echoImpl) Echo(ctx context.Context, call air.Echo_echo) error {
	in, err := call.Args().In()
	if err != nil {
		return err
	}
	r, err := call.AllocResults()
	if err != nil {
		return err
	}
	r.SetOut(in + in)
	return nil
}

type errorEchoImpl struct{}

func (errorEchoImpl) Echo(_ context.Context, call air.Echo_echo) error {
	call.Ack()
	return errors.New("reverb stopped")
}

func TestServerCall(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		echo := air.Echo_ServerToClient(echoImpl{}, nil)
		defer echo.Client.Release()

		ans, finish := echo.Echo(context.Background(), func(p air.Echo_echo_Params) error {
			err := p.SetIn("foo")
			return err
		})
		defer finish()
		result, err := ans.Struct()
		if err != nil {
			t.Errorf("echo.Echo() error: %v", err)
		}
		if out, err := result.Out(); err != nil {
			t.Errorf("echo.Echo() error: %v", err)
		} else if out != "foofoo" {
			t.Errorf("echo.Echo() = %q; want %q", out, "foofoo")
		}
	})
	t.Run("Error", func(t *testing.T) {
		echo := air.Echo_ServerToClient(errorEchoImpl{}, nil)
		defer echo.Client.Release()

		ans, finish := echo.Echo(context.Background(), func(p air.Echo_echo_Params) error {
			err := p.SetIn("foo")
			return err
		})
		defer finish()
		_, err := ans.Struct()
		if err == nil || !strings.Contains(err.Error(), "reverb stopped") || !strings.Contains(err.Error(), "echo") {
			t.Errorf("echo.Echo() error = %v; want error containing \"reverb stopped\" and \"echo\"", err)
		}
	})
	t.Run("Unimplemented", func(t *testing.T) {
		echo := air.Echo{Client: capnp.NewClient(server.New(nil, nil, nil, nil))}
		defer echo.Client.Release()

		ans, finish := echo.Echo(context.Background(), func(p air.Echo_echo_Params) error {
			err := p.SetIn("foo")
			return err
		})
		defer finish()
		_, err := ans.Struct()
		if err == nil {
			t.Error("echo.Echo() error = <nil>; want unimplemented")
		} else {
			if !capnp.IsUnimplemented(err) {
				t.Errorf("echo.Echo() error = %v; want unimplemented", err)
			}
			if !strings.Contains(err.Error(), "echo") {
				t.Errorf("echo.Echo() error = %v; want error containing \"echo\"", err)
			}
		}
	})
}

type callSeq uint32

func (seq *callSeq) GetNumber(ctx context.Context, call air.CallSequence_getNumber) error {
	r, err := call.AllocResults()
	if err != nil {
		return err
	}
	r.SetN(uint32(*seq))
	*seq++
	return nil
}

type lockCallSeq struct {
	n  uint32
	mu sync.Mutex
}

func (seq *lockCallSeq) GetNumber(ctx context.Context, call air.CallSequence_getNumber) error {
	seq.mu.Lock()
	defer seq.mu.Unlock()
	call.Ack()

	r, err := call.AllocResults()
	if err != nil {
		return err
	}
	r.SetN(seq.n)
	seq.n++
	return nil
}

func TestServerCallOrder(t *testing.T) {
	tests := []struct {
		name string
		seq  air.CallSequence
	}{
		{"NoAck", air.CallSequence_ServerToClient(new(callSeq), nil)},
		{"AckWithLocks", air.CallSequence_ServerToClient(new(callSeq), nil)},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			send := func() (air.CallSequence_getNumber_Results_Future, capnp.ReleaseFunc) {
				return test.seq.GetNumber(ctx, nil)
			}
			check := func(p air.CallSequence_getNumber_Results_Future, n uint32) {
				result, err := p.Struct()
				if err != nil {
					t.Errorf("seq.getNumber() error: %v; want %d", err, n)
				} else if result.N() != n {
					t.Errorf("seq.getNumber() = %d; want %d", result.N(), n)
				}
			}

			call0, finish := send()
			defer finish()
			call1, finish := send()
			defer finish()
			call2, finish := send()
			defer finish()
			call3, finish := send()
			defer finish()
			call4, finish := send()
			defer finish()

			check(call0, 0)
			check(call1, 1)
			check(call2, 2)
			check(call3, 3)
			check(call4, 4)
		})
		test.seq.Client.Release()
	}
}

func TestServerMaxConcurrentCalls(t *testing.T) {
	wait := make(chan struct{})
	echo := air.Echo_ServerToClient(blockingEchoImpl{wait}, &server.Policy{
		MaxConcurrentCalls: 2,
	})
	defer echo.Client.Release()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	call1, finish := echo.Echo(ctx, nil)
	defer finish()
	call2, finish := echo.Echo(ctx, nil)
	defer finish()
	go close(wait)
	call3, finish := echo.Echo(ctx, nil)
	defer finish()
	<-wait
	if _, err := call1.Struct(); err != nil {
		t.Error("Echo #1:", err)
	}
	if _, err := call2.Struct(); err != nil {
		t.Error("Echo #2:", err)
	}
	if _, err := call3.Struct(); err != nil {
		t.Error("Echo #3:", err)
	}
}

func TestServerShutdown(t *testing.T) {
	wait := make(chan struct{})
	echo := air.Echo_ServerToClient(blockingEchoImpl{wait}, nil)
	defer echo.Client.Release()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	call, finish := echo.Echo(ctx, nil)
	defer finish()
	echo.Client.Release()
	select {
	case <-call.Done():
		if _, err := call.Struct(); err == nil {
			t.Error("call finished without error")
		}
	default:
		t.Error("call not done after Shutdown")
	}
}

type blockingEchoImpl struct {
	wait <-chan struct{}
}

func (echo blockingEchoImpl) Echo(ctx context.Context, call air.Echo_echo) error {
	in, err := call.Args().In()
	if err != nil {
		return err
	}
	call.Ack()
	select {
	case <-echo.wait:
	case <-ctx.Done():
		return ctx.Err()
	}
	r, err := call.AllocResults()
	if err != nil {
		return err
	}
	r.SetOut(in)
	return nil
}

func TestPipelineCall(t *testing.T) {
	wait := make(chan struct{})
	var once sync.Once
	p := air.Pipeliner_ServerToClient(&pipeliner{
		factory: func(ctx context.Context) (*pipeliner, error) {
			first := false
			once.Do(func() { first = true })
			if !first {
				return nil, errors.New("can only create one pipeliner")
			}
			select {
			case <-wait:
				return new(pipeliner), nil
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		},
	}, nil)

	ctx := context.Background()
	baseAns, finish := p.NewPipeliner(ctx, nil)
	defer finish()
	qq := baseAns.Pipeliner()
	ans1, finish := qq.GetNumber(ctx, nil)
	defer finish()
	close(wait)
	<-baseAns.Done()
	ans2, finish := qq.GetNumber(ctx, nil)
	defer finish()

	result1, err := ans1.Struct()
	if err != nil {
		t.Errorf("GetNumber() #1: %v", err)
	} else if result1.N() != 0 {
		t.Errorf("GetNumber() #1 = %d; want 0", result1.N())
	}
	result2, err := ans2.Struct()
	if err != nil {
		t.Errorf("GetNumber() #2: %v", err)
	} else if result2.N() != 1 {
		t.Errorf("GetNumber() #2 = %d; want 1", result2.N())
	}
}

func TestBrokenPipelineCall(t *testing.T) {
	wait := make(chan struct{})
	p := air.Pipeliner_ServerToClient(brokenPipeliner{wait}, nil)

	ctx := context.Background()
	baseAns, finish := p.NewPipeliner(ctx, nil)
	defer finish()
	qq := baseAns.Pipeliner()
	ans1, finish := qq.GetNumber(ctx, nil)
	defer finish()
	close(wait)
	<-baseAns.Done()
	ans2, finish := qq.GetNumber(ctx, nil)
	defer finish()

	if _, err := ans1.Struct(); err == nil || !strings.Contains(err.Error(), "got no pipe") {
		t.Errorf("GetNumber() #1 error = %v; want \"got no pipe\"", err)
	}
	if _, err := ans2.Struct(); err == nil || !strings.Contains(err.Error(), "got no pipe") {
		t.Errorf("GetNumber() #2 error = %v; want \"got no pipe\"", err)
	}
}

type pipeliner struct {
	callSeq
	factory func(context.Context) (*pipeliner, error)
}

func (p *pipeliner) NewPipeliner(ctx context.Context, call air.Pipeliner_newPipeliner) error {
	if p.factory == nil {
		return errors.New("no factory present")
	}
	call.Ack()
	q, err := p.factory(ctx)
	if err != nil {
		return err
	}
	r, err := call.AllocResults()
	if err != nil {
		return err
	}
	r.SetPipeliner(air.Pipeliner_ServerToClient(q, nil))
	return nil
}

type brokenPipeliner struct {
	ready chan struct{}
}

func (p brokenPipeliner) GetNumber(ctx context.Context, call air.CallSequence_getNumber) error {
	call.Ack()
	<-p.ready
	return errors.New("got no number")
}

func (p brokenPipeliner) NewPipeliner(ctx context.Context, call air.Pipeliner_newPipeliner) error {
	call.Ack()
	<-p.ready
	return errors.New("got no pipe")
}

// Verify that if the first call calls .Ack(), the second will proceed without
// waiting for it to return.
func TestAckDoesntBlock(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := air.CallSequence_ServerToClient(&callSeqBlockN{
		blockingCalls: 1,
	}, nil)

	fut1, rel := client.GetNumber(ctx, nil)
	defer rel()

	fut2, rel := client.GetNumber(ctx, nil)
	defer rel()

	res2, err := fut2.Struct()
	assert.Nil(t, err, "Second call returns successfully")
	assert.Equal(t, res2.N(), uint32(2), "Second call returns 2.")

	cancel()
	_, err = fut1.Struct()
	assert.NotNil(t, err, "First call returns an error after cancel()")
}

// An implementation of CallSequence, where the first n calls never
// actually return.
type callSeqBlockN struct {
	blockingCalls, currentCall uint32
}

func (c *callSeqBlockN) GetNumber(ctx context.Context, p air.CallSequence_getNumber) error {
	c.currentCall++
	if c.currentCall > c.blockingCalls {
		res, err := p.AllocResults()
		if err != nil {
			panic(err)
		}
		res.SetN(c.currentCall)
		return nil
	} else {
		p.Ack()
		<-ctx.Done()
		return ctx.Err()
	}
}
