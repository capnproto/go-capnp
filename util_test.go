package capn_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	capn "github.com/glycerine/go-capnproto"
)

// some generally useful capnp/segment utilities

// shell out to display capnp bytes as human-readable text. Data flow:
//    in-memory capn segment -> stdin to capnp decode -> stdout human-readble string form
func CapnpDecodeSegment(seg *capn.Segment, capnpExePath string, capnpSchemaFilePath string, typeName string) string {

	// set defaults
	if capnpExePath == "" {
		capnpExePath = CheckAndGetCapnpPath()
	}

	if capnpSchemaFilePath == "" {
		capnpSchemaFilePath = "aircraftlib/aircraft.capnp"
	}

	if typeName == "" {
		typeName = "Z"
	}

	cs := []string{"decode", "--short", capnpSchemaFilePath, typeName}
	cmd := exec.Command(capnpExePath, cs...)
	cmdline := capnpExePath + " " + strings.Join(cs, " ")

	buf := new(bytes.Buffer)
	seg.WriteTo(buf)

	cmd.Stdin = buf

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

// reduce boilerplate, dump this segment to disk.
func SegToFile(seg *capn.Segment, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	seg.WriteTo(file)
	file.Close()
}

// disk file of a capn segment -> in-memory capn segment -> stdin to capnp decode -> stdout human-readble string form
func CapnFileToText(serializedCapnpFilePathToDisplay string, capnpSchemaFilePath string, capnpExePath string) (string, error) {

	// a) read file into Segment

	byteslice, err := ioutil.ReadFile(serializedCapnpFilePathToDisplay)
	if err != nil {
		return "", err
	}

	seg, nbytes, err := capn.ReadFromMemoryZeroCopy(byteslice)

	if err == io.EOF {
		return "", err
	}
	if err != nil {
		return "", err
	}
	if nbytes == 0 {
		return "", errors.New(fmt.Sprintf("did not expect 0 bytes back from capn.ReadFromMemoryZeroCopy() on reading file '%s'", serializedCapnpFilePathToDisplay))
	}

	// b) tell CapnpDecodeSegment() to show the human-readable-text form of the message
	// warning: CapnpDecodeSegment() may panic on you. It is a testing utility so that
	//  is desirable. For production, do something else.
	return CapnpDecodeSegment(seg, capnpExePath, capnpSchemaFilePath, "Z"), nil
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
		capnpSchemaFilePath = "aircraftlib/aircraft.capnp"
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
