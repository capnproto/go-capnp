package capnp_test

import (
	"testing"

	"golang.org/x/net/context"
	"zombiezen.com/go/capnproto"
	air "zombiezen.com/go/capnproto/aircraftlib"
)

type echoImpl struct{}

func (echoImpl) Echo(ctx context.Context, opts capnp.CallOptions, p air.Echo_echo_Params, r air.Echo_echo_Results) error {
	in := p.In()
	r.SetOut(in + in)
	return nil
}

func TestServerCall(t *testing.T) {
	echo := air.Echo_ServerToClient(echoImpl{})

	result, err := echo.Echo(context.Background(), func(p air.Echo_echo_Params) {
		p.SetIn("foo")
	}).Get()

	if err != nil {
		t.Errorf("echo.Echo() error: %v", err)
	}
	if out := result.Out(); out != "foofoo" {
		t.Errorf("echo.Echo() = %q; want %q", out, "foofoo")
	}
}

type callSeq uint32

func (seq *callSeq) GetNumber(
	ctx context.Context,
	opts capnp.CallOptions,
	p air.CallSequence_getNumber_Params,
	r air.CallSequence_getNumber_Results) error {
	r.SetN(uint32(*seq))
	*seq++
	capnp.Ack(opts)
	return nil
}

func TestServerCallOrder(t *testing.T) {
	seq := air.CallSequence_ServerToClient(new(callSeq))
	ctx := context.Background()
	send := func() *air.CallSequence_getNumber_Results_Promise {
		return seq.GetNumber(ctx, func(air.CallSequence_getNumber_Params) {})
	}
	check := func(p *air.CallSequence_getNumber_Results_Promise, n uint32) {
		result, err := p.Get()
		if err != nil {
			t.Errorf("seq.getNumber() error: %v; want %d", err, n)
		} else if result.N() != n {
			t.Errorf("seq.getNumber() = %d; want %d", result.N(), n)
		}
	}

	call0 := send()
	call1 := send()
	call2 := send()
	call3 := send()
	call4 := send()

	check(call0, 0)
	check(call1, 1)
	check(call2, 2)
	check(call3, 3)
	check(call4, 4)
}
