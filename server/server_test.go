package server_test

import (
	"context"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

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
	call.Go()
	return errors.New("reverb stopped")
}

func TestServerCall(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		echo := air.Echo_ServerToClient(echoImpl{})
		defer echo.Release()

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
		echo := air.Echo_ServerToClient(errorEchoImpl{})
		defer echo.Release()

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
		echo := air.Echo(capnp.NewClient(server.New(nil, nil, nil)))
		defer echo.Release()

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
	t.Run("Unimplemented hook", func(t *testing.T) {
		t.Parallel()
		var echoText = "are you there?"
		var proxyReceived atomic.Value

		// start a proxy server with hook
		srv := server.New(nil, nil, nil)
		srv.HandleUnknownMethod = func(method capnp.Method) *server.Method {
			sm := server.Method{
				Method: method,
				Impl:   nil,
			}
			sm.Impl = func(ctx context.Context, call *server.Call) error {
				echoArgs := air.Echo_echo_Params(call.Args())
				inText, err := echoArgs.In()
				require.NoError(t, err)
				proxyReceived.Store(inText)
				// pretend we received an answer
				echo := air.Echo_echo{Call: call}
				resp, _ := echo.AllocResults()
				err = resp.SetOut(inText)
				return err
			}
			return &sm
		}
		blankBoot := capnp.NewClient(srv)
		echoClient := air.Echo(blankBoot)
		defer echoClient.Release()

		ans, finish := echoClient.Echo(context.Background(), func(p air.Echo_echo_Params) error {
			err := p.SetIn(echoText)
			return err
		})
		defer finish()
		resp, err := ans.Struct()
		answerOut, _ := resp.Out()
		rxValue := proxyReceived.Load()
		require.Equal(t, echoText, rxValue)
		assert.Equal(t, echoText, answerOut)
		assert.NoError(t, err, "echo.Echo() error != <nil>; want success")
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

func TestServerCallOrder(t *testing.T) {
	tests := []struct {
		name string
		seq  air.CallSequence
	}{
		{"NoGo", air.CallSequence_ServerToClient(new(callSeq))},
		{"GoWithLocks", air.CallSequence_ServerToClient(new(callSeq))},
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
		test.seq.Release()
	}
}

func TestServerShutdown(t *testing.T) {
	wait := make(chan struct{})
	echo := air.Echo_ServerToClient(blockingEchoImpl{wait})
	defer echo.Release()
	ctx, cancel := context.WithCancel(context.Background())
	call, finish := echo.Echo(ctx, nil)
	defer finish()
	echo.Release()

	// Even though we've dropped the client, existing calls should
	// still go through:
	select {
	case <-call.Done():
		t.Error("call finished before cancel()")
	case <-time.After(10 * time.Millisecond):
	}

	cancel()
	<-call.Done() // Will hang if cancel doesn't stop the call.

	if _, err := call.Struct(); err == nil {
		t.Error("call finished without error")
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
	call.Go()
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
	})

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
	assert.NoError(t, err)
	assert.Equal(t, uint32(0), result1.N())

	result2, err := ans2.Struct()
	assert.NoError(t, err)
	assert.Equal(t, uint32(1), result2.N())
}

func TestBrokenPipelineCall(t *testing.T) {
	wait := make(chan struct{})
	p := air.Pipeliner_ServerToClient(brokenPipeliner{wait})

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
	call.Go()
	q, err := p.factory(ctx)
	if err != nil {
		return err
	}
	r, err := call.AllocResults()
	if err != nil {
		return err
	}
	r.SetPipeliner(air.Pipeliner_ServerToClient(q))
	return nil
}

type brokenPipeliner struct {
	ready chan struct{}
}

func (p brokenPipeliner) GetNumber(ctx context.Context, call air.CallSequence_getNumber) error {
	call.Go()
	<-p.ready
	return errors.New("got no number")
}

func (p brokenPipeliner) NewPipeliner(ctx context.Context, call air.Pipeliner_newPipeliner) error {
	call.Go()
	<-p.ready
	return errors.New("got no pipe")
}

// Verify that if the first call calls .Go(), the second will proceed without
// waiting for it to return.
func TestGoDoesntBlock(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client := air.CallSequence_ServerToClient(&callSeqBlockN{
		blockingCalls: 1,
	})

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
		p.Go()
		<-ctx.Done()
		return ctx.Err()
	}
}
