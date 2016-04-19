package main

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"zombiezen.com/go/capnproto2"
	"zombiezen.com/go/capnproto2/internal/capnptool"
	"zombiezen.com/go/capnproto2/internal/schema"
)

type schemaFile struct {
	name string
	src  string
	skip bool // if false, use as a target
}

func mustCompileSchema(t *testing.T, fallback []byte, files ...schemaFile) schema.CodeGeneratorRequest {
	req, err := compileSchema(fallback, files...)
	if err != nil {
		t.Fatal(err)
	}
	return req
}

// compileSchema either invokes the capnp tool to compile the given
// files or unmarshals the base64-encoded, gzip-compressed fallback.
func compileSchema(fallback []byte, files ...schemaFile) (schema.CodeGeneratorRequest, error) {
	tool, err := capnptool.Find()
	if err != nil {
		// TODO(light): log tool missing
		msg, err := capnp.Unmarshal(fallback)
		if err != nil {
			return schema.CodeGeneratorRequest{}, err
		}
		return schema.ReadRootCodeGeneratorRequest(msg)
	}
	d, err := ioutil.TempDir("", "capnpc-go-tests")
	if err != nil {
		return schema.CodeGeneratorRequest{}, err
	}
	defer os.RemoveAll(d)
	for _, f := range files {
		ff, err := os.Create(filepath.Join(d, f.name))
		if err != nil {
			return schema.CodeGeneratorRequest{}, err
		}
		_, err = ff.WriteString(f.src)
		cerr := ff.Close()
		if err != nil {
			return schema.CodeGeneratorRequest{}, err
		}
		if cerr != nil {
			return schema.CodeGeneratorRequest{}, cerr
		}
	}

	args := make([]string, 0, len(files)+3)
	args = append(args, "compile", "-o-", "--src-prefix="+d)
	for _, f := range files {
		if f.skip {
			continue
		}
		args = append(args, filepath.Join(d, f.name))
	}
	out, err := tool.Run(nil, args...)
	if err != nil {
		return schema.CodeGeneratorRequest{}, err
	}
	if !bytes.Equal(out, fallback) {
		// TODO(light): print schema names
		return schema.CodeGeneratorRequest{}, fmt.Errorf("schema =\n%s\nFallback out of date?", encodeByteLiteral(out))
	}
	msg, err := capnp.Unmarshal(out)
	if err != nil {
		return schema.CodeGeneratorRequest{}, err
	}
	return schema.ReadRootCodeGeneratorRequest(msg)
}

type byteLiteral string

func encodeByteLiteral(in []byte) byteLiteral {
	buf := new(bytes.Buffer)
	bw := base64.NewEncoder(base64.StdEncoding, buf)
	zw := gzip.NewWriter(bw)
	zw.Write(in)
	zw.Close()
	bw.Close()
	return byteLiteral(buf.String())
}

func (lit byteLiteral) String() string {
	return string(lit)
}

func (lit byteLiteral) decode() []byte {
	r, err := gzip.NewReader(base64.NewDecoder(base64.StdEncoding, strings.NewReader(string(lit))))
	if err != nil {
		return nil
	}
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil
	}
	return b
}

func TestGoCapnpNodeMap(t *testing.T) {
	req := mustCompileSchema(t,
		byteLiteral("H4sIAAAJbogA/2IAgh0MMMDEwAokjzMxMOwEYnkgWwGIWRnYGD726osLJew8xAnm/+PvyLn7L1BG6yIDGvBkZGAoAtK+QJodSdwDyGdmYERXzjBZQf/BrFSBXxBzFRhwmWsJ1LoISDuimWuDw1yhB0FvKtvKTkDM5cBpri5QaxWQNkQzVweHucffn1KR3FKylJC5skCtWUBaEc1cGRzmbnXhU2PcZvAQYi4jTnN5gVqbgLQgmrk8OMx9oF2pLFC/ch8hc/8CcRdIDZq5DDjMjZs8P3Dvta6jEHNNcJr7kQESDl8ZUM39wIDdXJg5zBgyqOAhEHsB8VMgLgdiW0ZIesUF0vP1khML8gqs8hJzU+FeY2TgIWAPXF9yaXFJfm5JZUEqXI4k/Xn5JYnpCH2EAFwfXBeJ9mXmFuQXlZCuryAxOTsxPRVZjiT9KfnJ5OiD8WWg+mDp1hDId2JA5A9dINuKAZH+NIFsJQZEvlSF8mH5XxHINmJAlDOyQHYUEMPKMxBfC4hh3oYFG8wbsOCHRR96MoAlJxawu5ng7gbxmYAiAtCQQA8HRqi7UDyPZA4sH4DKY1A6B+UzdizhBQtnQAAAAP//BZSLzMgFAAA=").decode(),
		schemaFile{
			name: "go.capnp",
			src: `@0xd12a1c51fedd6c88;

annotation package(file) :Text;
annotation import(file) :Text;
annotation doc(struct, field, enum) :Text;
annotation tag(enumerant) :Text;
annotation notag(enumerant) :Void;
annotation customtype(field) :Text;
annotation name(struct, field, union, enum, enumerant, interface, method, param, annotation, const, group) :Text;
$package("capnp");`,
		})
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
