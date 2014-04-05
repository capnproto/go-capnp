package capn_test

import (
	"bytes"
	"encoding/hex"
	"fmt"

	capn "github.com/glycerine/go-capnproto"
	air "github.com/glycerine/go-capnproto/aircraftlib"
)

func ExampleReadFromStream() {
	s := capn.NewBuffer(nil)
	d := air.NewRootZdate(s)
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
	d = air.ReadRootZdate(s)
	fmt.Printf("year %d, month %d, day %d\n", d.Year(), d.Month(), d.Day())
}
