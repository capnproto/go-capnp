package capnp_test

import (
	"bytes"

	"zombiezen.com/go/capnproto"
	air "zombiezen.com/go/capnproto/internal/aircraftlib"
)

const schemaPath = "internal/aircraftlib/aircraft.capnp"

func zdateFilledSegment(n int32, packed bool) (*capnp.Segment, []byte) {
	seg := capnp.NewBuffer(nil)
	z := air.NewRootZ(seg)
	list := air.NewZdate_List(seg, n)
	// hand added a Set() method to messages_test.go, so plist not needed
	plist := capnp.PointerList(list)

	for i := 0; i < int(n); i++ {
		d := air.NewZdate(seg)
		d.SetMonth(12)
		d.SetDay(7)
		d.SetYear(int16(2004 + i))
		plist.Set(i, capnp.Pointer(d))
		//list.Set(i, d)
	}
	z.SetZdatevec(list)

	buf := bytes.Buffer{}
	if packed {
		seg.WriteToPacked(&buf)
	} else {
		seg.WriteTo(&buf)
	}
	return seg, buf.Bytes()
}

func zdateReader(n int32, packed bool) *bytes.Reader {
	_, byteSlice := zdateFilledSegment(n, packed)
	return bytes.NewReader(byteSlice)
}

// actually return n segments back-to-back.
// WriteTo will automatically add the stream header word with message length.
//
func zdateReaderNBackToBack(n int, packed bool) *bytes.Reader {

	buf := bytes.Buffer{}

	for i := 0; i < n; i++ {
		seg := capnp.NewBuffer(nil)
		d := air.NewRootZdate(seg)
		d.SetMonth(12)
		d.SetDay(7)
		d.SetYear(int16(2004 + i))

		if packed {
			seg.WriteToPacked(&buf)
		} else {
			seg.WriteTo(&buf)
		}
	}

	return bytes.NewReader(buf.Bytes())
}

func zdataFilledSegment(n int) (*capnp.Segment, []byte) {
	seg := capnp.NewBuffer(nil)
	z := air.NewRootZ(seg)
	d := air.NewZdata(seg)

	b := make([]byte, n)
	for i := 0; i < len(b); i++ {
		b[i] = byte(i)
	}
	d.SetData(b)
	z.SetZdata(d)

	buf := bytes.Buffer{}
	seg.WriteTo(&buf)
	return seg, buf.Bytes()
}

func zdataReader(n int) *bytes.Reader {
	_, byteSlice := zdataFilledSegment(n)
	return bytes.NewReader(byteSlice)
}
