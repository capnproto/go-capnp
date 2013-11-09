package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/jmckaskill/go-capnproto"
)

func main() {
	b := capn.NewBuffer(nil)
	d := NewZdate(b)
	d.SetYear(2004)
	d.SetMonth(12)
	d.SetDay(7)
	buf := bytes.Buffer{}
	b.WriteTo(&buf)
	fmt.Println(hex.EncodeToString(buf.Bytes()))
}
