package capn_test

import (
	"bytes"
	capn "github.com/jmckaskill/go-capnproto"
)

func zdateReader(n int, packed bool) *bytes.Reader {
	s := capn.NewBuffer(nil)
	d := NewRootZdate(s)
	d.SetMonth(12)
	d.SetDay(7)
	buf := bytes.Buffer{}
	for i := 0; i < n; i++ {
		d.SetYear(int16(2004 + i))
		if packed {
			s.WriteToPacked(&buf)
		} else {
			s.WriteTo(&buf)
		}
	}
	return bytes.NewReader(buf.Bytes())
}

func zdataReader(n int) *bytes.Reader {
	s := capn.NewBuffer(nil)
	d := NewRootZdata(s)
	b := make([]byte, n)
	for i := 0; i < len(b); i++ {
		b[i] = byte(i)
	}
	d.SetData(b)
	buf := bytes.Buffer{}
	s.WriteTo(&buf)
	return bytes.NewReader(buf.Bytes())
}
