package capn_test

import (
	"bytes"

	capn "github.com/glycerine/go-capnproto"
)

func zdateFilledSegment(n int, packed bool) (*capn.Segment, []byte) {
	seg := capn.NewBuffer(nil)
	z := NewRootZ(seg)
	list := NewZdateList(seg, n)
	plist := capn.PointerList(list)

	for i := 0; i < n; i++ {
		d := NewZdate(seg)
		d.SetMonth(12)
		d.SetDay(7)
		d.SetYear(int16(2004 + i))
		plist.Set(i, capn.Object(d))
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

func zdateReader(n int, packed bool) *bytes.Reader {
	_, byteSlice := zdateFilledSegment(n, packed)
	return bytes.NewReader(byteSlice)
}

func zdataFilledSegment(n int) (*capn.Segment, []byte) {
	seg := capn.NewBuffer(nil)
	d := NewRootZdata(seg)
	b := make([]byte, n)
	for i := 0; i < len(b); i++ {
		b[i] = byte(i)
	}
	d.SetData(b)
	buf := bytes.Buffer{}
	seg.WriteTo(&buf)
	return seg, buf.Bytes()
}

func zdataReader(n int) *bytes.Reader {
	_, byteSlice := zdataFilledSegment(n)
	return bytes.NewReader(byteSlice)
}
