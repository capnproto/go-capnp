package text

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"zombiezen.com/go/capnproto2"
	"zombiezen.com/go/capnproto2/schemas"
	"zombiezen.com/go/capnproto2/std/capnp/schema"
)

func readTestFile(name string) ([]byte, error) {
	path := filepath.Join("testdata", name)
	return ioutil.ReadFile(path)
}

func TestEncode(t *testing.T) {
	tests := []struct {
		constID uint64
		text    string
	}{
		{0xc0b634e19e5a9a4e, `(key = "42", value = (int32 = -123))`},
		{0x967c8fe21790b0fb, `(key = "float", value = (float64 = 3.14))`},
		{0xdf35cb2e1f5ea087, `(key = "bool", value = (bool = false))`},
	}

	data, err := readTestFile("txt.capnp.out")
	if err != nil {
		t.Fatal(err)
	}
	reg := new(schemas.Registry)
	err = reg.Register(&schemas.Schema{
		Bytes: data,
		Nodes: []uint64{
			0x8df8bc5abdc060a6,
			0xd3602730c572a43b,
		},
	})
	if err != nil {
		t.Fatal("Adding to registry: %v", err)
	}
	msg, err := capnp.Unmarshal(data)
	if err != nil {
		t.Fatal("Unmarshaling txt.capnp.out:", err)
	}
	req, err := schema.ReadRootCodeGeneratorRequest(msg)
	if err != nil {
		t.Fatal("Reading code generator request txt.capnp.out:", err)
	}
	nodes, err := req.Nodes()
	if err != nil {
		t.Fatal(err)
	}
	nodeMap := make(map[uint64]schema.Node, nodes.Len())
	for i := 0; i < nodes.Len(); i++ {
		n := nodes.At(i)
		nodeMap[n.Id()] = n
	}

	for _, test := range tests {
		c := nodeMap[test.constID]
		if !c.IsValid() {
			t.Errorf("Can't find node %#x; skipping", test.constID)
			continue
		}
		dn, _ := c.DisplayName()
		if c.Which() != schema.Node_Which_const {
			t.Errorf("%s @%#x is a %v, not const; skipping", dn, test.constID, c.Which())
			continue
		}

		typ, err := c.Const().Type()
		if err != nil {
			t.Errorf("(%s @%#x - %s).const.value: %v", dn, test.constID, err)
			continue
		}
		if typ.Which() != schema.Type_Which_structType {
			t.Errorf("(%s @%#x).const.type is a %v; want struct", dn, test.constID, typ.Which())
			continue
		}
		tid := typ.StructType().TypeId()

		v, err := c.Const().Value()
		if err != nil {
			t.Errorf("(%s @%#x).const.value: %v", dn, test.constID, err)
			continue
		}
		if v.Which() != schema.Value_Which_structValue {
			t.Errorf("(%s @%#x).const.value is a %v; want struct", dn, test.constID, v.Which())
			continue
		}
		sv, err := v.StructValuePtr()
		if err != nil {
			t.Errorf("(%s @%#x).const.value.struct: %v", dn, test.constID, err)
			continue
		}

		buf := new(bytes.Buffer)
		enc := NewEncoder(buf)
		enc.UseRegistry(reg)
		if err := enc.Encode(tid, sv.Struct()); err != nil {
			t.Errorf("Encode(%#x, (%s @%#x).const.value.struct): %v", tid, dn, test.constID, err)
			continue
		}
		if text := buf.String(); text != test.text {
			t.Errorf("Encode(%#x, (%s @%#x).const.value.struct) = %q; want %q", tid, dn, test.constID, text, test.text)
			continue
		}
	}
}
