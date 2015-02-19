package capn_test

import (
	"testing"

	"golang.org/x/net/context"
	"zombiezen.com/go/capnproto"
	air "zombiezen.com/go/capnproto/aircraftlib"
)

type echoImpl struct{}

func (echoImpl) Echo(ctx context.Context, opts capn.CallOptions, p air.Echo_echo_Params, r air.Echo_echo_Results) error {
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
