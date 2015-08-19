package capnp_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	cv "github.com/smartystreets/goconvey/convey"
	"zombiezen.com/go/capnproto"
	air "zombiezen.com/go/capnproto/internal/aircraftlib"
)

func Example() {
	// Make a brand new empty message.
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))

	// If you want runtime-type identification, this is easily obtained. Just
	// wrap everything in a struct that contains a single anoymous union (e.g. struct Z).
	// Then always set a Z as the root object in you message/first segment.
	// The cost of the extra word of storage is usually worth it, as
	// then human readable output is easily obtained via a shell command such as
	//
	// $ cat binary.cpz | capnp decode aircraft.capnp Z
	//
	// If you need to conserve space, and know your content in advance, it
	// isn't necessary to use an anonymous union. Just supply the type name
	// in place of 'Z' in the decode command above.

	// There can only be one root.  Subsequent NewRoot* calls will set the root
	// pointer and orphan the previous root.
	z, err := air.NewRootZ(seg)
	if err != nil {
		panic(err)
	}

	// then non-root objects:
	aircraft, err := z.NewAircraft()
	if err != nil {
		panic(err)
	}
	b737, err := aircraft.NewB737()
	if err != nil {
		panic(err)
	}
	planebase, err := b737.NewBase()
	if err != nil {
		panic(err)
	}

	// Set primitive fields
	planebase.SetCanFly(true)
	planebase.SetName("Henrietta")
	planebase.SetRating(100)
	planebase.SetMaxSpeed(876) // km/hr
	// if we don't set capacity, it will get the default value, in this case 0.
	//planebase.SetCapacity(26020) // Liters fuel

	// Creating a list
	homes, err := air.NewAirport_List(seg, 2)
	if err != nil {
		panic(err)
	}
	homes.Set(0, air.Airport_jfk)
	homes.Set(1, air.Airport_lax)
	// Setting a list field
	planebase.SetHomes(homes)

	// Ready to write!

	// You can write to memory...
	buf, err := msg.Marshal()
	if err != nil {
		panic(err)
	}
	_ = buf

	// ... or write to an io.Writer.
	file, err := ioutil.TempFile("", "go-capnproto")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	defer os.Remove(file.Name())
	err = capnp.NewEncoder(file).Encode(msg)
	if err != nil {
		panic(err)
	}

	// Read back and view that file in human readable format. Defined in util_test.go
	text, err := CapnFileToText(file.Name(), schemaPath, "")
	if err != nil {
		panic(err)
	}
	fmt.Printf("here is our aircraft:\n")
	fmt.Printf("%s\n", text)

	// Output:
	// here is our aircraft:
	// (aircraft = (b737 = (base = (name = "Henrietta", homes = [jfk, lax], rating = 100, canFly = true, capacity = 0, maxSpeed = 876))))
}

func TestVoidUnionSetters(t *testing.T) {
	want := CapnpEncode(`(b = void)`, "VoidUnion")

	cv.Convey("Given a VoidUnion set to b", t, func() {
		cv.Convey("then the go-capnproto serialization should match the capnp c++ serialization", func() {
			msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
			cv.So(err, cv.ShouldEqual, nil)
			voidUnion, err := air.NewRootVoidUnion(seg)
			cv.So(err, cv.ShouldEqual, nil)
			voidUnion.SetB()

			act, err := msg.Marshal()
			cv.So(err, cv.ShouldEqual, nil)

			cv.So(act, cv.ShouldResemble, want)
		})
	})
}
