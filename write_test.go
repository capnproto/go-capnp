package capn_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	capn "github.com/glycerine/go-capnproto"
	air "github.com/glycerine/go-capnproto/aircraftlib"
	cv "github.com/smartystreets/goconvey/convey"
)

func ExampleAirplaneWrite() string {

	fname := "out.write_test.airplane.cpz"

	// make a brand new, empty segment (message)
	seg := capn.NewBuffer(nil)

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

	z := air.NewRootZ(seg) // root should be allocated first.
	// There must be only one root.

	// then non-root objects:
	aircraft := air.NewAircraft(seg)
	b737 := air.NewB737(seg)
	planebase := air.NewPlaneBase(seg)

	// how to create a list. Requires a cast at the moment.
	homes := air.NewAirportList(seg, 2)
	uint16list := capn.UInt16List(homes) // cast to the underlying type
	uint16list.Set(0, uint16(air.AIRPORT_JFK))
	uint16list.Set(1, uint16(air.AIRPORT_LAX))

	// set the primitive fields
	planebase.SetCanFly(true)
	planebase.SetName("Henrietta")
	planebase.SetRating(100)
	planebase.SetMaxSpeed(876) // km/hr
	// if we don't set capacity, it will get the default value, in this case 0.
	//planebase.SetCapacity(26020) // Liters fuel

	// set a list field
	planebase.SetHomes(homes)

	// wire up the pointers between objects
	b737.SetBase(planebase)
	aircraft.SetB737(b737)
	z.SetAircraft(aircraft)

	// ready to write

	// example of writing to memory
	buf := bytes.Buffer{}
	seg.WriteTo(&buf)

	// example of writing to file. Just use WriteTo().
	// We could have used SegToFile(seg, fname) from
	// util_test.go intead, but this makes it clear how easy it is.
	file, err := os.Create(fname)
	defer file.Close()
	if err != nil {
		panic(err)
	}
	seg.WriteTo(file)

	// readback and view that file in human readable format. Defined in util_test.go
	text, err := CapnFileToText(fname, "aircraftlib/aircraft.capnp", "")
	if err != nil {
		panic(err)
	}
	fmt.Printf("here is our aircraft:\n")
	fmt.Printf("%s\n", text)

	return text
}

func TestAircraftWrite(t *testing.T) {

	observedText := ExampleAirplaneWrite()
	expectedText := `(aircraft = (b737 = (base = (name = "Henrietta", homes = [jfk, lax], rating = 100, canFly = true, capacity = 0, maxSpeed = 876))))`

	cv.Convey("When we run the ExampleAirplaneWrite() function in write_test.go", t, func() {
		cv.Convey("Then we should see the human readable B737 example struct we expect", func() {
			cv.So(observedText, cv.ShouldEqual, expectedText)
		})
	})

}

func TestVoidUnionSetters(t *testing.T) {
	want := CapnpEncode(`(b = void)`, "VoidUnion")

	cv.Convey("Given a VoidUnion set to b", t, func() {
		cv.Convey("then the go-capnproto serialization should match the capnp c++ serialization", func() {
			seg := capn.NewBuffer(nil)
			voidUnion := air.NewRootVoidUnion(seg)
			voidUnion.SetB()

			var buf bytes.Buffer
			seg.WriteTo(&buf)
			act := buf.Bytes()

			cv.So(act, cv.ShouldResemble, want)
		})
	})
}
