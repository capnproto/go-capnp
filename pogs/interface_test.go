package pogs

import (
	"golang.org/x/net/context"
	"testing"
	"zombiezen.com/go/capnproto2"
	air "zombiezen.com/go/capnproto2/internal/aircraftlib"
)

type simpleEcho struct{}

func checkFatal(t *testing.T, name string, err error) {
	if err != nil {
		t.Fatalf("%s for TestInsertIFace: %v", name, err)
	}
}

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
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	checkFatal(t, "NewMessage", err)
	h, err := air.NewRootHoth(seg)
	checkFatal(t, "NewRootHoth", err)
	err = Insert(air.Hoth_TypeID, h.Struct, Hoth{
		Base: EchoBase{Echo: air.Echo_ServerToClient(simpleEcho{})},
	})
	checkFatal(t, "Insert", err)
	base, err := h.Base()
	checkFatal(t, "h.Base", err)
	echo := base.Echo()

	testEcho(t, echo)
}

func TestExtractIFace(t *testing.T) {
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	checkFatal(t, "NewMessage", err)
	h, err := air.NewRootHoth(seg)
	checkFatal(t, "NewRootHoth", err)
	base, err := air.NewEchoBase(seg)
	checkFatal(t, "NewEchoBase", err)
	h.SetBase(base)
	base.SetEcho(air.Echo_ServerToClient(simpleEcho{}))

	extractedHoth := Hoth{}
	err = Extract(&extractedHoth, air.Hoth_TypeID, h.Struct)
	checkFatal(t, "Extract", err)

	testEcho(t, extractedHoth.Base.Echo)
}

func testEcho(t *testing.T, echo air.Echo) {
	expected := "Hello!"
	result, err := echo.Echo(context.TODO(), func(p air.Echo_echo_Params) error {
		p.SetIn(expected)
		return nil
	}).Struct()
	checkFatal(t, "Echo", err)
	actual, err := result.Out()
	checkFatal(t, "result.Out", err)
	if actual != expected {
		t.Fatal("Echo result did not match input; "+
			"wanted %q but got %q.", expected, actual)
	}
}
