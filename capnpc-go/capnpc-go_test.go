package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"zombiezen.com/go/capnproto2"
	"zombiezen.com/go/capnproto2/internal/schema"
)

func mustReadTestFile(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", name)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func mustReadGeneratorRequest(t *testing.T, name string) schema.CodeGeneratorRequest {
	data := mustReadTestFile(t, name)
	msg, err := capnp.Unmarshal(data)
	if err != nil {
		t.Fatalf("Unmarshaling %s: %v", name, err)
	}
	req, err := schema.ReadRootCodeGeneratorRequest(msg)
	if err != nil {
		t.Fatalf("Reading code generator request %s: %v", name, err)
	}
	return req
}

func TestGoCapnpNodeMap(t *testing.T) {
	req := mustReadGeneratorRequest(t, "go.capnp.out")
	nodes, err := buildNodeMap(req)
	if err != nil {
		t.Error("buildNodeMap:", err)
	}
	want := []uint64{
		0xd12a1c51fedd6c88,
		0xbea97f1023792be0,
		0xe130b601260e44b5,
		0xc58ad6bd519f935e,
		0xa574b41924caefc7,
		0xc8768679ec52e012,
		0xfa10659ae02f2093,
		0xc2b96012172f8df1,
	}
	for _, k := range want {
		if nodes[k] == nil {
			t.Errorf("missing @%#x from node map", k)
		}
	}
}

func TestRemoteScope(t *testing.T) {
	type scopeTest struct {
		name        string
		varID       uint64
		initImports []importSpec

		remoteName string
		remoteNew  string
		imports    []importSpec
	}
	tests := []scopeTest{
		{
			name:       "same-file struct",
			varID:      0x84efedc75e99768d, // scopes.fooVar
			remoteName: "Foo",
			remoteNew:  "NewFoo",
		},
		{
			name:       "different file struct",
			varID:      0x836faf1834d91729, // scopes.otherFooVar
			remoteName: "otherscopes.Foo",
			remoteNew:  "otherscopes.NewFoo",
			imports: []importSpec{
				{name: "otherscopes", path: "zombiezen.com/go/capnproto2/capnpc-go/testdata/otherscopes"},
			},
		},
	}
	req := mustReadGeneratorRequest(t, "scopes.capnp.out")
	nodes, err := buildNodeMap(req)
	if err != nil {
		t.Fatal("buildNodeMap:", err)
	}
	collect := func(test scopeTest) (g *generator, n *node, from *node, ok bool) {
		g = newGenerator(0xd68755941d99d05e, nodes, genoptions{})
		v := nodes[test.varID]
		if v == nil {
			t.Errorf("Can't find const @%#x for %s test", test.varID, test.name)
			return nil, nil, nil, false
		}
		if v.Which() != schema.Node_Which_const {
			t.Errorf("Type of const @%#x in %s test is a %v node; want const. Check the test.", test.varID, test.name, v.Which())
			return nil, nil, nil, false
		}
		varType, _ := v.Const().Type()
		// TODO(light): just use the type
		varTypeNode := nodes[varType.StructType().TypeId()]
		for _, i := range test.initImports {
			g.imports.add(i)
		}
		return g, varTypeNode, v, true
	}
	for _, test := range tests {
		g, n, from, ok := collect(test)
		if !ok {
			continue
		}
		rn, err := g.RemoteName(n, from)
		if err != nil {
			t.Errorf("%s: g.RemoteName(nodes[%#x].Const().Type(), nodes[%#x]) error: %v", test.name, test.varID, test.varID, err)
			continue
		}
		if rn != test.remoteName {
			t.Errorf("%s: g.RemoteName(nodes[%#x].Const().Type(), nodes[%#x]) = %q; want %q", test.name, test.varID, test.varID, rn, test.remoteName)
			continue
		}
		if !hasExactImports(test.imports, g.imports) {
			t.Errorf("%s: g.RemoteName(nodes[%#x].Const().Type(), nodes[%#x]); g.imports = %s; want %s", test.name, test.varID, test.varID, formatImportSpecs(g.imports.usedImports()), formatImportSpecs(test.imports))
			continue
		}
	}
	// TODO(light): add RemoteNew tests
}

func hasExactImports(specs []importSpec, imp imports) bool {
	used := imp.usedImports()
	if len(used) != len(specs) {
		return false
	}
outer:
	for i := range specs {
		for j := range used {
			if specs[i] == used[j] {
				continue outer
			}
		}
		return false
	}
	return true
}

func formatImportSpecs(specs []importSpec) string {
	var buf bytes.Buffer
	for i, s := range specs {
		if i > 0 {
			buf.WriteString("; ")
		}
		buf.WriteString(s.String())
	}
	return buf.String()
}
