package pogs

import (
	"golang.org/x/net/context"
	"testing"
	"zombiezen.com/go/capnproto2"
	air "zombiezen.com/go/capnproto2/internal/aircraftlib"
)

type simpleEcho struct{}

func (s simpleEcho) Echo(p air.Echo_echo) error {
	text, err := p.Params.In()
	if err != nil {
		return err
	}
	p.Results.SetOut(text)
	return nil
}

type EchoBase struct {
	Echo air.Echo
}

type Hoth struct {
	Base EchoBase
}

func TestInsertIFace(t *testing.T) {
	checkErr := func(name string, err error) {
		if err != nil {
			t.Fatalf("%s for TestInsertIFace: %v", name, err)
		}
	}
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	checkErr("NewMessage", err)
	h, err := air.NewRootHoth(seg)
	checkErr("NewRootHoth", err)
	err = Insert(air.Hoth_TypeID, h.Struct, Hoth{
		Base: EchoBase{Echo: air.Echo_ServerToClient(simpleEcho{})},
	})
	checkErr("Insert", err)
	base, err := h.Base()
	checkErr("h.Base", err)
	echo := base.Echo()
	expected := "Hello!"
	result, err := echo.Echo(context.TODO(), func(p air.Echo_echo_Params) error {
		p.SetIn(expected)
		return nil
	}).Struct()
	checkErr("Echo", err)
	actual, err := result.Out()
	checkErr("result.Out", err)
	if actual != expected {
		t.Fatal("TestInsertIFace: Echo result did not match input; "+
			"wanted %q but got %q.", expected, actual)
	}
}
