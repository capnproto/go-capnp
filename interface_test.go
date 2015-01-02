package capn_test

import (
	"errors"
	"testing"

	"github.com/glycerine/go-capnproto"
	air "github.com/glycerine/go-capnproto/aircraftlib"
)

func TestInterfaceSet(t *testing.T) {
	cl := air.NewEcho(capn.ErrorClient(errors.New("foo")))
	s := capn.NewBuffer(nil)
	base := air.NewRootEchoBase(s)

	base.SetEcho(cl)

	if base.Echo() != cl {
		t.Errorf("base.Echo() = %#v; want %#v", base.Echo(), cl)
	}
}

func TestInterfaceCopyToOtherMessage(t *testing.T) {
	cl := air.NewEcho(capn.ErrorClient(errors.New("foo")))
	s1 := capn.NewBuffer(nil)
	base1 := air.NewRootEchoBase(s1)
	base1.SetEcho(cl)

	s2 := capn.NewBuffer(nil)
	hoth2 := air.NewRootHoth(s2)
	hoth2.SetBase(base1)

	if hoth2.Base().Echo() != cl {
		t.Errorf("hoth2.Base().Echo() = %#v; want %#v", hoth2.Base().Echo(), cl)
	}
	tab2 := s2.Message.CapTable()
	if len(tab2) == 1 {
		if tab2[0] != cl.GenericClient() {
			t.Error("s2.Message.CapTable()[0] = %#v; want %#v", tab2[0], cl.GenericClient())
		}
	} else {
		t.Errorf("len(s2.Message.CapTable()) = %d; want 1", len(tab2))
	}
}

func TestInterfaceCopyToOtherMessageWithCaps(t *testing.T) {
	cl := air.NewEcho(capn.ErrorClient(errors.New("foo")))
	s1 := capn.NewBuffer(nil)
	base1 := air.NewRootEchoBase(s1)
	base1.SetEcho(cl)

	s2 := capn.NewBuffer(nil)
	s2.Message.AddCap(nil)
	hoth2 := air.NewRootHoth(s2)
	hoth2.SetBase(base1)

	if hoth2.Base().Echo() != cl {
		t.Errorf("hoth2.Base().Echo() = %#v; want %#v", hoth2.Base().Echo(), cl)
	}
	tab2 := s2.Message.CapTable()
	if len(tab2) != 2 {
		t.Errorf("len(s2.Message.CapTable()) = %d; want 2", len(tab2))
	}
}
