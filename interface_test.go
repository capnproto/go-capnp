package capn_test

import (
	"testing"

	"github.com/glycerine/go-capnproto"
	air "github.com/glycerine/go-capnproto/aircraftlib"
)

func TestInterfaceSet(t *testing.T) {
	s := capn.NewBuffer(nil)
	base := air.NewRootEchoBase(s)
	base.SetEcho(air.NewEcho(nil))
}

func TestInterfaceCopyToOtherMessage(t *testing.T) {
	s1 := capn.NewBuffer(nil)
	base1 := air.NewRootEchoBase(s1)
	base1.SetEcho(air.NewEcho(nil))

	s2 := capn.NewBuffer(nil)
	hoth2 := air.NewRootHoth(s1)
	hoth2.SetBase(base1)

	if hoth2.Base().Echo().IsNull() {
		t.Error("hoth2.Base().Echo() = nil")
	}
	tab2 := s2.Message.CapTable()
	if len(tab2) == 1 {
		if tab2[0] == nil {
			t.Error("s2.Message.CapTable()[0] = nil")
		}
	} else {
		t.Errorf("len(s2.Message.CapTable()) = %d; want 1", len(tab2))
	}
}

func TestInterfaceCopyToOtherMessageWithCaps(t *testing.T) {
	s1 := capn.NewBuffer(nil)
	base1 := air.NewRootEchoBase(s1)
	base1.SetEcho(air.NewEcho(nil))

	s2 := capn.NewBuffer(nil)
	s2.Message.AddCap(nil)
	hoth2 := air.NewRootHoth(s1)
	hoth2.SetBase(base1)

	if hoth2.Base().Echo().IsNull() {
		t.Error("hoth2.Base().Echo() = nil")
	}
	tab2 := s2.Message.CapTable()
	if len(tab2) != 2 {
		t.Errorf("len(s2.Message.CapTable()) = %d; want 2", len(tab2))
	}
}
