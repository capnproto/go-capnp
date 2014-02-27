package capn_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	capn "github.com/jmckaskill/go-capnproto"
)

func ExampleReadFromStream() {
	s := capn.NewBuffer(nil)
	d := NewRootZdate(s)
	d.SetYear(2004)
	d.SetMonth(12)
	d.SetDay(7)
	buf := bytes.Buffer{}
	s.WriteTo(&buf)

	fmt.Println(hex.EncodeToString(buf.Bytes()))

	// Read
	s, err := capn.ReadFromStream(&buf, nil)
	if err != nil {
		fmt.Printf("read error %v\n", err)
		return
	}
	d = ReadRootZdate(s)
	fmt.Printf("year %d, month %d, day %d\n", d.Year(), d.Month(), d.Day())
}
