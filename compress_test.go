package capn_test

import (
	"bytes"
	"github.com/jmckaskill/go-capnproto"
	"io/ioutil"
	"testing"
)

func TestPackingRoundtrip(t *testing.T) {
	var buf bytes.Buffer
	c := capn.NewCompressor(&buf)
	b := []byte{0x8, 0, 0, 0, 0x3, 0, 0x2, 0, 0x19, 0, 0, 0, 0xaa, 0x1, 0, 0}
	n, err := c.Write(b)
	if err != nil {
		panic(err)
	}
	var expected = len(b)
	if expected != n {
		t.Fatalf("expected %v bytes, got %v", expected, n)
	}
	exp := []byte{0x51, 0x08, 0x03, 0x02, 0x31, 0x19, 0xaa, 0x01}
	if !bytes.Equal(buf.Bytes(), exp) {
		t.Fatalf("expected %x bytes, got %x", exp, buf.Bytes())
	}
	dc := capn.NewDecompressor(&buf)
	readBuf, err := ioutil.ReadAll(dc)
	if err != nil {
		panic(err)
	}
	readBuf = readBuf[:n]
	if !bytes.Equal(b, readBuf) {
		t.Fatalf("expected %x bytes, got %x", b, readBuf)
	}
}

func TestPackingTag1(t *testing.T) {
	var buf bytes.Buffer
	c := capn.NewCompressor(&buf)
	_, err := c.Write(make([]byte, 32))
	if err != nil {
		panic(err)
	}
	if bytes.Compare(buf.Bytes(), []byte{0x0, 0x3}) != 0 {
		t.Fatalf("invalid packing")
	}
}

func TestPackingTag2(t *testing.T) {
	var buf bytes.Buffer
	c := capn.NewCompressor(&buf)
	b := make([]byte, 32)
	for i := 0; i < len(b); i++ {
		b[i] = 0x8a
	}
	_, err := c.Write(b)
	if err != nil {
		panic(err)
	}
	b = make([]byte, 34)
	for i := 0; i < len(b); i++ {
		b[i] = 0x8a
	}
	b[0] = 0xff
	b[9] = 0x3
	if bytes.Compare(buf.Bytes(), b) != 0 {
		t.Fatalf("invalid packing")
	}
}
