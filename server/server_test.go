package server_test

import (
	"sync"
	"testing"

	"context"
	"zombiezen.com/go/capnproto2"
	air "zombiezen.com/go/capnproto2/internal/aircraftlib"
	"zombiezen.com/go/capnproto2/server"
)

type echoImpl struct{}

func (echoImpl) Echo(ctx context.Context, call air.Echo_echo) error {
	in, err := call.Args.In()
	if err != nil {
		return err
	}
	call.Results.SetOut(in + in)
	return nil
}

func TestServerCall(t *testing.T) {
	echo := air.Echo_ServerToClient(echoImpl{}, nil)
	defer func() {
		if err := echo.Client.Close(); err != nil {
			t.Error("Close:", err)
		}
	}()

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
}

type callSeq uint32

func (seq *callSeq) GetNumber(ctx context.Context, call air.CallSequence_getNumber) error {
	call.Results.SetN(uint32(*seq))
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

	call.Results.SetN(seq.n)
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
		if err := test.seq.Client.Close(); err != nil {
			t.Errorf("Close %s: %v", test.name, err)
		}
	}
}

func TestServerMaxConcurrentCalls(t *testing.T) {
	wait := make(chan struct{})
	echo := air.Echo_ServerToClient(blockingEchoImpl{wait}, &server.Policy{
		MaxConcurrentCalls: 2,
	})
	defer func() {
		if err := echo.Client.Close(); err != nil {
			t.Error("Close:", err)
		}
	}()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	call1, finish := echo.Echo(ctx, nil)
	defer finish()
	call2, finish := echo.Echo(ctx, nil)
	defer finish()
	call3, finish := echo.Echo(ctx, nil)
	defer finish()
	select {
	case <-call3.Done():
		if _, err := call3.Struct(); err == nil {
			t.Error("overload call returned early success")
		}
	default:
		t.Error("overload call not finished; want immediate overload error")
	}
	close(wait)
	call1.Struct()
	call2.Struct()
}

func TestServerClose(t *testing.T) {
	wait := make(chan struct{})
	echo := air.Echo_ServerToClient(blockingEchoImpl{wait}, nil)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	call, finish := echo.Echo(ctx, nil)
	defer finish()
	if err := echo.Client.Close(); err != nil {
		t.Error("Close:", err)
	}
	select {
	case <-call.Done():
		if _, err := call.Struct(); err == nil {
			t.Error("call finished without error")
		}
	default:
		t.Error("call not done after Close")
	}
}

type blockingEchoImpl struct {
	wait <-chan struct{}
}

func (echo blockingEchoImpl) Echo(ctx context.Context, call air.Echo_echo) error {
	in, err := call.Args.In()
	if err != nil {
		return err
	}
	call.Ack()
	select {
	case <-echo.wait:
	case <-ctx.Done():
		return ctx.Err()
	}
	call.Results.SetOut(in)
	return nil
}
