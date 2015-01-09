package capn_test

import (
	"testing"

	"github.com/glycerine/go-capnproto"
	air "github.com/glycerine/go-capnproto/aircraftlib"
	"golang.org/x/net/context"
)

type echoImpl struct{}

func (echoImpl) Echo(ctx context.Context, p air.Echo_echo_Params, r air.Echo_echo_Results) error {
	in := p.In()
	r.SetOut(in + in)
	return nil
}

func TestServerCall(t *testing.T) {
	echo := air.NewEcho(capn.NewServer(air.Echo_Methods(nil, echoImpl{})))

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
