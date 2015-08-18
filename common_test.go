package capnp_test

import (
	"bytes"

	"zombiezen.com/go/capnproto"
	air "zombiezen.com/go/capnproto/internal/aircraftlib"
)

const schemaPath = "internal/aircraftlib/aircraft.capnp"

func zdateFilledSegment(n int32, packed bool) (*capnp.Segment, []byte) {
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		panic(err)
	}
	z, err := air.NewRootZ(seg)
	if err != nil {
		panic(err)
	}
	list, err := air.NewZdate_List(seg, n)
	if err != nil {
		panic(err)
	}

	for i := 0; i < int(n); i++ {
		d, err := air.NewZdate(seg)
		if err != nil {
			panic(err)
		}
		d.SetMonth(12)
		d.SetDay(7)
		d.SetYear(int16(2004 + i))
		list.Set(i, d)
	}
	z.SetZdatevec(list)

	if packed {
		b, err := msg.MarshalPacked()
		if err != nil {
			panic(err)
		}
		return seg, b
	}
	b, err := msg.Marshal()
	if err != nil {
		panic(err)
	}
	return seg, b
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
		msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
		if err != nil {
			panic(err)
		}
		d, err := air.NewRootZdate(seg)
		if err != nil {
			panic(err)
		}
		d.SetMonth(12)
		d.SetDay(7)
		d.SetYear(int16(2004 + i))

		if packed {
			b, err := msg.MarshalPacked()
			if err != nil {
				panic(err)
			}
			buf.Write(b)
		} else {
			b, err := msg.Marshal()
			if err != nil {
				panic(err)
			}
			buf.Write(b)
		}
	}

	return bytes.NewReader(buf.Bytes())
}

func zdataFilledSegment(n int) (*capnp.Segment, []byte) {
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		panic(err)
	}
	z, err := air.NewRootZ(seg)
	if err != nil {
		panic(err)
	}
	d, err := air.NewZdata(seg)
	if err != nil {
		panic(err)
	}

	b := make([]byte, n)
	for i := 0; i < len(b); i++ {
		b[i] = byte(i)
	}
	d.SetData(b)
	z.SetZdata(d)

	buf, err := msg.Marshal()
	if err != nil {
		panic(err)
	}
	return seg, buf
}

func zdataReader(n int) *bytes.Reader {
	_, byteSlice := zdataFilledSegment(n)
	return bytes.NewReader(byteSlice)
}
