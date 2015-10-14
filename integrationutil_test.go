package capnp_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"zombiezen.com/go/capnproto2"
	air "zombiezen.com/go/capnproto2/internal/aircraftlib"
)

const schemaPath = "internal/aircraftlib/aircraft.capnp"

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

// some generally useful capnp/segment utilities

// shell out to display capnp bytes as human-readable text. Data flow:
//    in-memory capn segment -> stdin to capnp decode -> stdout human-readble string form
func CapnpDecodeSegment(seg *capnp.Segment, capnpExePath string, capnpSchemaFilePath string, typeName string) string {

	// set defaults
	if capnpExePath == "" {
		capnpExePath = CheckAndGetCapnpPath()
	}

	if capnpSchemaFilePath == "" {
		capnpSchemaFilePath = schemaPath
	}

	if typeName == "" {
		typeName = "Z"
	}

	cs := []string{"decode", "--short", capnpSchemaFilePath, typeName}
	cmd := exec.Command(capnpExePath, cs...)
	cmdline := capnpExePath + " " + strings.Join(cs, " ")

	buf, err := seg.Message().Marshal()
	if err != nil {
		panic(err)
	}

	cmd.Stdin = bytes.NewReader(buf)

	var errout bytes.Buffer
	cmd.Stderr = &errout

	bs, err := cmd.Output()
	if err != nil {
		if err.Error() == "exit status 1" {
			cwd, _ := os.Getwd()
			fmt.Fprintf(os.Stderr, "\nCall to capnp in CapnpDecodeSegment(): '%s' in dir '%s' failed with status 1\n", cmdline, cwd)
			fmt.Printf("stderr: '%s'\n", string(errout.Bytes()))
			fmt.Printf("stdout: '%s'\n", string(bs))
		}
		panic(err)
	}
	return strings.TrimSpace(string(bs))
}

// return path to capnp if 'which' can find it. Feel free to replace this with
//   a more general configuration mechanism.
func CheckAndGetCapnpPath() string {

	path, err := exec.LookPath("capnp")
	if err != nil {
		panic(fmt.Sprintf("could not locate the capnp executable: put the capnp executable in your path: %s", err))
	}

	cmd := exec.Command(path, "id")
	bs, err := cmd.Output()
	if err != nil || string(bs[:3]) != `@0x` {
		panic(fmt.Sprintf("%s id did not function: put a working capnp executable in your path. Err: %s", path, err))
	}

	return path
}

// take an already (packed or unpacked, depending on the packed flag) buffer of a serialized segment, and display it
func CapnpDecodeBuf(buf []byte, capnpExePath string, capnpSchemaFilePath string, typeName string, packed bool) string {

	// set defaults
	if capnpExePath == "" {
		capnpExePath = CheckAndGetCapnpPath()
	}

	if capnpSchemaFilePath == "" {
		capnpSchemaFilePath = schemaPath
	}

	if typeName == "" {
		typeName = "Z"
	}

	cs := []string{"decode", "--short", capnpSchemaFilePath, typeName}
	if packed {
		cs = []string{"decode", "--short", "--packed", capnpSchemaFilePath, typeName}
	}
	cmd := exec.Command(capnpExePath, cs...)
	cmdline := capnpExePath + " " + strings.Join(cs, " ")

	cmd.Stdin = bytes.NewReader(buf)

	var errout bytes.Buffer
	cmd.Stderr = &errout

	bs, err := cmd.Output()
	if err != nil {
		if err.Error() == "exit status 1" {
			cwd, _ := os.Getwd()
			fmt.Fprintf(os.Stderr, "\nCall to capnp in CapnpDecodeBuf(): '%s' in dir '%s' failed with status 1\n", cmdline, cwd)
			fmt.Printf("stderr: '%s'\n", string(errout.Bytes()))
			fmt.Printf("stdout: '%s'\n", string(bs))
		}
		panic(err)
	}
	return strings.TrimSpace(string(bs))
}

func CapnpEncode(msg string, typ string) []byte {
	capnpPath, err := exec.LookPath("capnp")
	//capnpPath, err := exec.LookPath("tee")
	if err != nil {
		panic(err)
	}
	if !FileExists(capnpPath) {
		panic(fmt.Sprintf("could not locate capnp tool in PATH"))
	}

	schfn := schemaPath
	args := []string{"encode", schfn, typ}
	cmdline := fmt.Sprintf("%s %s %s %s", capnpPath, "encode", schfn, typ)
	//fmt.Printf("cmdline = %s\n", cmdline)
	c := exec.Command(capnpPath, args...)

	var o bytes.Buffer
	c.Stdout = &o

	var in bytes.Buffer
	in.Write([]byte(msg))
	c.Stdin = &in

	err = c.Run()
	if err != nil {
		panic(fmt.Errorf("tried to run %s, got err:%s", cmdline, err))
	}
	return o.Bytes()
}

func CapnpDecode(input []byte, typ string) []byte {
	capnpPath, err := exec.LookPath("capnp")
	//capnpPath, err := exec.LookPath("tee")
	if err != nil {
		panic(err)
	}
	if !FileExists(capnpPath) {
		panic(fmt.Sprintf("could not locate capnp tool in PATH"))
	}

	schfn := schemaPath
	args := []string{"decode", "--short", schfn, typ}
	cmdline := fmt.Sprintf("%s %s %s %s %s", capnpPath, "decode", "--short", schfn, typ)
	fmt.Printf("cmdline = %s\n", cmdline)
	c := exec.Command(capnpPath, args...)

	var o bytes.Buffer
	c.Stdout = &o

	var e bytes.Buffer
	c.Stderr = &e

	var in bytes.Buffer
	in.Write(input)
	c.Stdin = &in

	err = c.Run()
	if err != nil {
		fmt.Printf("tried to run %s, got err:%s and stderr: '%s'", cmdline, err, e.Bytes())
		panic(err)
	}
	return o.Bytes()
}

func FileExists(name string) bool {
	fi, err := os.Stat(name)
	if err != nil {
		return false
	}
	if fi.IsDir() {
		return false
	}
	return true
}

// encodeTestMessage encodes the textual Cap'n Proto message to unpacked
// binary using the capnp tool, or returns the fallback if the tool fails.
func encodeTestMessage(typ string, text string, fallback []byte) ([]byte, error) {
	// TODO(light): use fallback if tool not present.
	b := CapnpEncode(text, typ)
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
