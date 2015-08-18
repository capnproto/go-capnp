package capnp_test

import (
	"errors"
	"testing"

	"zombiezen.com/go/capnproto"
	air "zombiezen.com/go/capnproto/internal/aircraftlib"
)

func TestInterfaceSet(t *testing.T) {
	cl := air.Echo{Client: capnp.ErrorClient(errors.New("foo"))}
	_, s, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatal(err)
	}
	base, err := air.NewRootEchoBase(s)
	if err != nil {
		t.Fatal(err)
	}

	base.SetEcho(cl)

	if base.Echo() != cl {
		t.Errorf("base.Echo() = %#v; want %#v", base.Echo(), cl)
	}
}

func TestInterfaceCopyToOtherMessage(t *testing.T) {
	cl := air.Echo{Client: capnp.ErrorClient(errors.New("foo"))}
	_, s1, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatal(err)
	}
	base1, err := air.NewRootEchoBase(s1)
	if err != nil {
		t.Fatal(err)
	}
	if err := base1.SetEcho(cl); err != nil {
		t.Fatal(err)
	}

	_, s2, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatal(err)
	}
	hoth2, err := air.NewRootHoth(s2)
	if err != nil {
		t.Fatal(err)
	}
	if err := hoth2.SetBase(base1); err != nil {
		t.Fatal(err)
	}

	if base, err := hoth2.Base(); err != nil {
		t.Errorf("hoth2.Base() error: %v", err)
	} else if base.Echo() != cl {
		t.Errorf("hoth2.Base().Echo() = %#v; want %#v", base.Echo(), cl)
	}
	tab2 := s2.Message().CapTable
	if len(tab2) == 1 {
		if tab2[0] != cl.Client {
			t.Errorf("s2.Message().CapTable[0] = %#v; want %#v", tab2[0], cl.Client)
		}
	} else {
		t.Errorf("len(s2.Message().CapTable) = %d; want 1", len(tab2))
	}
}

func TestInterfaceCopyToOtherMessageWithCaps(t *testing.T) {
	cl := air.Echo{Client: capnp.ErrorClient(errors.New("foo"))}
	_, s1, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatal(err)
	}
	base1, err := air.NewRootEchoBase(s1)
	if err != nil {
		t.Fatal(err)
	}
	if err := base1.SetEcho(cl); err != nil {
		t.Fatal(err)
	}

	_, s2, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatal(err)
	}
	s2.Message().AddCap(nil)
	hoth2, err := air.NewRootHoth(s2)
	if err != nil {
		t.Fatal(err)
	}
	if err := hoth2.SetBase(base1); err != nil {
		t.Fatal(err)
	}

	if base, err := hoth2.Base(); err != nil {
		t.Errorf("hoth2.Base() error: %v", err)
	} else if base.Echo() != cl {
		t.Errorf("hoth2.Base().Echo() = %#v; want %#v", base.Echo(), cl)
	}
	tab2 := s2.Message().CapTable
	if len(tab2) != 2 {
		t.Errorf("len(s2.Message().CapTable) = %d; want 2", len(tab2))
	}
}

func TestMethodString(t *testing.T) {
	tests := []struct {
		m *capnp.Method
		s string
	}{
		{
			&capnp.Method{
				InterfaceID: 0x8e5322c1e9282534,
				MethodID:    1,
			},
			"@0x8e5322c1e9282534.@1",
		},
		{
			&capnp.Method{
				InterfaceID:   0x8e5322c1e9282534,
				MethodID:      1,
				InterfaceName: "aircraftlib:Echo",
				MethodName:    "foo",
			},
			"aircraftlib:Echo.foo",
		},
	}
	for _, test := range tests {
		if s := test.m.String(); s != test.s {
			t.Errorf("%#v.String() = %q; want %q", test.m, s, test.s)
		}
	}
}
