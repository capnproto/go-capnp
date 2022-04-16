//go:build mktemplates
// +build mktemplates

// Build tag so that users who run `go get capnproto.org/go/capnp/v3/...` don't install this command.
// cd internal/cmd/mktemplates && go build -tags=mktemplates

// mktemplates is a command to regenerate capnpc-go/templates.go.
package main

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"os"
	"sort"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Fprintln(os.Stderr, "usage: mktemplates OUT DIR")
		os.Exit(64)
	}
	dir := os.Args[2]
	names, err := listdir(dir)
	if err != nil {
		fatalln(err)
	}
	genbuf := new(bytes.Buffer)
	err = generateGo(genbuf, os.Args, names)
	if err != nil {
		fatalln("generating code:", err)
	}
	code, err := format.Source(genbuf.Bytes())
	if err != nil {
		fatalln("formatting code:", err)
	}
	outname := os.Args[1]
	out, err := os.Create(outname)
	if err != nil {
		fatalf("opening destination %s: %v", outname, err)
	}
	_, err = out.Write(code)
	cerr := out.Close()
	if err != nil {
		fatalf("write to %s: %v", outname, err)
	}
	if cerr != nil {
		fatalln(err)
	}
}

func generateGo(w io.Writer, args []string, names []string) error {
	// TODO(light): collect errors
	fmt.Fprintln(w, "// Code generated from templates directory. DO NOT EDIT.")
	fmt.Fprintln(w)
	fmt.Fprintln(w, "//go:generate", strings.Join(args, " "))
	fmt.Fprintln(w)
	fmt.Fprintln(w, "package main")
	for _, name := range names {
		if strings.HasPrefix(name, "_") {
			continue
		}
		fmt.Fprintf(w, "func render%s(r renderer, p %sParams) error {\n\treturn r.Render(%[2]q, p)\n}\n", strings.Title(name), name)
	}
	return nil
}

type template struct {
	name    string
	content string
}

func listdir(name string) ([]string, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	names, err := f.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	n := 0
	for _, name := range names {
		if !strings.HasPrefix(name, ".") {
			names[n] = name
			n++
		}
	}
	names = names[:n]
	sort.Strings(names)
	return names, nil
}

func fatalln(args ...interface{}) {
	var buf bytes.Buffer
	buf.WriteString("mktemplates: ")
	fmt.Fprintln(&buf, args...)
	os.Stderr.Write(buf.Bytes())
	os.Exit(1)
}

func fatalf(format string, args ...interface{}) {
	var buf bytes.Buffer
	buf.WriteString("mktemplates: ")
	fmt.Fprintf(&buf, format, args...)
	if !bytes.HasSuffix(buf.Bytes(), []byte{'\n'}) {
		buf.Write([]byte{'\n'})
	}
	os.Stderr.Write(buf.Bytes())
	os.Exit(1)
}
