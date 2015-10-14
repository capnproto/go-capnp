package capnp_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"

	"zombiezen.com/go/capnproto2"
	air "zombiezen.com/go/capnproto2/internal/aircraftlib"
)

const schemaPath = "internal/aircraftlib/aircraft.capnp"

// capnpTool is the path to the capnp command-line tool.
type capnpTool string

var toolCache struct {
	init sync.Once
	tool capnpTool
	err  error
}

// findCapnpTool searches the PATH for the capnp tool.
func findCapnpTool() (capnpTool, error) {
	toolCache.init.Do(func() {
		path, err := exec.LookPath("capnp")
		if err != nil {
			toolCache.err = err
			return
		}
		toolCache.tool = capnpTool(path)
	})
	return toolCache.tool, toolCache.err
}

// run executes the tool with the given stdin and arguments returns the stdout.
func (tool capnpTool) run(stdin io.Reader, args ...string) ([]byte, error) {
	c := exec.Command(string(tool), args...)
	c.Stdin = stdin
	stderr := new(bytes.Buffer)
	c.Stderr = stderr
	out, err := c.Output()
	if err != nil {
		return nil, fmt.Errorf("run `%s`: %v; stderr:\n%s", strings.Join(c.Args, " "), err, stderr)
	}
	return out, nil
}

// encode encodes Cap'n Proto text into the binary representation.
func (tool capnpTool) encode(typ string, text string) ([]byte, error) {
	return tool.run(strings.NewReader(text), "encode", schemaPath, typ)
}

// decode decodes a Cap'n Proto message into text.
func (tool capnpTool) decode(typ string, r io.Reader) (string, error) {
	out, err := tool.run(r, "decode", "--short", schemaPath, typ)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// decodePacked decodes a packed Cap'n Proto message into text.
func (tool capnpTool) decodePacked(typ string, r io.Reader) (string, error) {
	out, err := tool.run(r, "decode", "--short", "--packed", schemaPath, typ)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func initNester(t *testing.T, n air.Nester1Capn, strs ...string) {
	tl, err := capnp.NewTextList(n.Segment(), int32(len(strs)))
	if err != nil {
		t.Fatalf("initNester(..., %q): NewTextList: %v", strs, err)
	}
	if err := n.SetStrs(tl); err != nil {
		t.Fatalf("initNester(..., %q): SetStrs: %v", strs, err)
	}
	for i, s := range strs {
		if err := tl.Set(i, s); err != nil {
			t.Fatalf("initNester(..., %q): set strs[%d]: %v", strs, i, err)
		}
	}
}

func zdateFilledMessage(t testing.TB, n int32) *capnp.Message {
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatal(err)
	}
	z, err := air.NewRootZ(seg)
	if err != nil {
		t.Fatal(err)
	}
	list, err := air.NewZdate_List(seg, n)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < int(n); i++ {
		d, err := air.NewZdate(seg)
		if err != nil {
			t.Fatal(err)
		}
		d.SetMonth(12)
		d.SetDay(7)
		d.SetYear(int16(2004 + i))
		list.Set(i, d)
	}
	z.SetZdatevec(list)

	return msg
}

func zdataFilledMessage(t testing.TB, n int) *capnp.Message {
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatal(err)
	}
	z, err := air.NewRootZ(seg)
	if err != nil {
		t.Fatal(err)
	}
	d, err := air.NewZdata(seg)
	if err != nil {
		t.Fatal(err)
	}
	b := make([]byte, n)
	for i := 0; i < len(b); i++ {
		b[i] = byte(i)
	}
	d.SetData(b)
	z.SetZdata(d)
	return msg
}

// encodeTestMessage encodes the textual Cap'n Proto message to unpacked
// binary using the capnp tool, or returns the fallback if the tool fails.
func encodeTestMessage(typ string, text string, fallback []byte) ([]byte, error) {
	tool, err := findCapnpTool()
	if err != nil {
		// TODO(light): log tool missing
		return fallback, nil
	}
	b, err := tool.encode(typ, text)
	if err != nil {
		return nil, fmt.Errorf("%s value %q encode failed: %v", typ, text, err)
	}
	if !bytes.Equal(b, fallback) {
		return nil, fmt.Errorf("%s value %q =\n%s; fallback is\n%s\nFallback out of date?", typ, text, hex.Dump(b), hex.Dump(fallback))
	}
	return b, nil
}

// mustEncodeTestMessage encodes the textual Cap'n Proto message to unpacked
// binary using the capnp tool, or returns the fallback if the tool fails.
func mustEncodeTestMessage(t testing.TB, typ string, text string, fallback []byte) []byte {
	b, err := encodeTestMessage(typ, text, fallback)
	if err != nil {
		if _, fname, line, ok := runtime.Caller(1); ok {
			t.Fatalf("%s:%d: %v", filepath.Base(fname), line, err)
		} else {
			t.Fatal(err)
		}
	}
	return b
}
