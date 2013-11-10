package main

import (
	"github.com/kaos/capnp_test"
	"github.com/jmckaskill/go-capnproto"
	"fmt"
	"os"
	"strconv"
)

func main() {
	decode := os.Args[1] == "decode"
	test := os.Args[2]

	if decode {
		msg, err := capn.ReadFromStream(os.Stdin, nil)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read error %v\n", err)
			os.Exit(1)
		}

		switch test {
		case "simpleTest":
			s := capnp_test.ReadRootSimpleTestStruct(msg)
			fmt.Printf("(int = %d, msg = %s)\n",
				s.Int(), strconv.Quote(s.Msg()))
		case "textListTypeTest":
			s := capnp_test.ReadRootListTest(msg)
			fmt.Printf("(textList = [")
			for i, v := range s.TextList().ToArray() {
				if i > 0 {
					fmt.Printf(", ")
				}
				fmt.Printf("%s", strconv.Quote(v))
			}
			fmt.Printf("])\n")
		case "uInt8DefaultValueTest":
			s := capnp_test.ReadRootTestDefaults(msg)
			fmt.Printf("(uInt8Field = %d)\n", s.UInt8Field())
		case "constTest":
			s := capnp_test.ReadRootSimpleTestStruct(msg)
			fmt.Printf("(msg = %s)\n", strconv.Quote(s.Msg()))
		default:
			os.Exit(127)
		}

	} else {
		msg := capn.NewBuffer(nil)

		switch test {
		case "simpleTest":
			s := capnp_test.NewRootSimpleTestStruct(msg)
			s.SetInt(1234567890)
			s.SetMsg("a short message...")
		case "textListTypeTest":
			s := capnp_test.NewRootListTest(msg)
			l := s.Segment.NewTextList(3)
			l.Set(0, "foo")
			l.Set(1, "bar")
			l.Set(2, "baz")
			s.SetTextList(l)
		case "uInt8DefaultValueTest":
			s := capnp_test.NewRootTestDefaults(msg)
			s.SetUInt8Field(0)
		case "constTest":
			s := capnp_test.NewRootSimpleTestStruct(msg)
			s.SetMsg(capnp_test.ConstTestValue)
		default:
			os.Exit(127)
		}

		msg.WriteTo(os.Stdout)
	}
}
