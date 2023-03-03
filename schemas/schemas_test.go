package schemas_test

import (
	"testing"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/internal/schema"
	"capnproto.org/go/capnp/v3/schemas"
	gocp "capnproto.org/go/capnp/v3/std/go"
)

func TestDefaultFind(t *testing.T) {
	gocp.RegisterSchema(schemas.DefaultRegistry)
	if s := schemas.Find(0xdeadbeef); s != nil {
		t.Errorf("schemas.Find(0xdeadbeef) = %d-byte slice; want nil", len(s))
	}
	s := schemas.Find(gocp.Package)
	if s == nil {
		t.Fatalf("schemas.Find(%#x) = nil", gocp.Package)
	}
	msg, err := capnp.Unmarshal(s)
	if err != nil {
		t.Fatalf("capnp.Unmarshal(schemas.Find(%#x)) error: %v", gocp.Package, err)
	}
	req, err := schema.ReadRootCodeGeneratorRequest(msg)
	if err != nil {
		t.Fatalf("ReadRootCodeGeneratorRequest error: %v", err)
	}
	nodes, err := req.Nodes()
	if err != nil {
		t.Fatalf("req.Nodes() error: %v", err)
	}
	for i := 0; i < nodes.Len(); i++ {
		n := nodes.At(i)
		if n.Id() == gocp.Package {
			// Found
			if n.Which() != schema.Node_Which_annotation {
				t.Errorf("found node %#x which = %v; want annotation", gocp.Package, n.Which())
			}
			return
		}
	}
	t.Fatalf("could not find node %#x in registry", gocp.Package)
}

func TestNotFound(t *testing.T) {
	reg := new(schemas.Registry)
	_, err := reg.Find(0)
	if err == nil {
		t.Error("new(schemas.Registry).Find(0) = nil; want not found error")
	}
	if !schemas.IsNotFound(err) {
		t.Errorf("new(schemas.Registry).Find(0) = %v; want not found error", err)
	}
}
