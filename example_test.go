package capnp_test

import (
	"encoding/hex"
	"fmt"

	"zombiezen.com/go/capnproto"
	air "zombiezen.com/go/capnproto/internal/aircraftlib"
)

func ExampleUnmarshal() {
	msg, s, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		fmt.Printf("allocation error %v\n", err)
		return
	}
	d, err := air.NewRootZdate(s)
	if err != nil {
		fmt.Printf("root error %v\n", err)
		return
	}
	d.SetYear(2004)
	d.SetMonth(12)
	d.SetDay(7)
	data, err := msg.Marshal()
	if err != nil {
		fmt.Printf("marshal error %v\n", err)
		return
	}

	fmt.Println(hex.EncodeToString(data))

	// Read
	msg, err = capnp.Unmarshal(data)
	if err != nil {
		fmt.Printf("unmarshal error %v\n", err)
		return
	}
	d, err = air.ReadRootZdate(msg)
	if err != nil {
		fmt.Printf("read root error %v\n", err)
		return
	}
	fmt.Printf("year %d, month %d, day %d\n", d.Year(), d.Month(), d.Day())
	// Output:
	// year 2004, month 12, day 7
}
