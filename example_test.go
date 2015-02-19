package capnp_test

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"zombiezen.com/go/capnproto"
	air "zombiezen.com/go/capnproto/aircraftlib"
)

func ExampleReadFromStream() {
	s := capnp.NewBuffer(nil)
	d := air.NewRootZdate(s)
	d.SetYear(2004)
	d.SetMonth(12)
	d.SetDay(7)
	buf := bytes.Buffer{}
	s.WriteTo(&buf)

	fmt.Println(hex.EncodeToString(buf.Bytes()))

	// Read
	s, err := capnp.ReadFromStream(&buf, nil)
	if err != nil {
		fmt.Printf("read error %v\n", err)
		return
	}
	d = air.ReadRootZdate(s)
	fmt.Printf("year %d, month %d, day %d\n", d.Year(), d.Month(), d.Day())
}
