package main

import (
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
