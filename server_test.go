package capnp_test

import (
	"testing"

	"golang.org/x/net/context"
	air "zombiezen.com/go/capnproto/aircraftlib"
)

type echoImpl struct{}

func (echoImpl) Echo(call air.Echo_echo) error {
	in := call.Params.In()
	call.Results.SetOut(in + in)
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
