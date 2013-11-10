package capn

import (
	"bytes"
	"errors"
	"io"
)

type Compressor struct {
	w io.Writer
}

type Decompressor struct {
	r     io.Reader
	buf   [8]byte
	bufsz int
	zeros int
	raw   int
}

func NewCompressor(w io.Writer) *Compressor {
	return &Compressor{w: w}
}

// WriteToPacked writes the message that the segment is part of to the
// provided stream in packed form.
func (s *Segment) WriteToPacked(w io.Writer) (int64, error) {
	c := Compressor{w: w}
	return s.WriteTo(&c)
}

func NewDecompressor(r io.Reader) *Decompressor {
	return &Decompressor{r: r}
}

// ReadFromPackedStream reads a single message from the stream r in packed
// form returning the first segment. buf can be specified in order to reuse
// the buffer (or it is allocated each call if nil).
func ReadFromPackedStream(r io.Reader, buf *bytes.Buffer) (*Segment, error) {
	c := Decompressor{r: r}
	return ReadFromStream(&c, buf)
}

func min(a, b int) int {
	if b < a {
		return b
	}
	return a
}

func (c *Decompressor) Read(v []byte) (n int, err error) {
	if c.raw > 0 {
		n, err = c.r.Read(v[:min(len(v), c.raw)])
		c.raw -= n
		return
	}

	if c.zeros > 0 {
		n = min(len(v), c.zeros)
		for i := range v[:n] {
			v[i] = 0
		}
		c.zeros -= n
		return
	}

	if c.bufsz > 0 {
		n = copy(v, c.buf[:c.bufsz])
		c.bufsz -= n
	}

	for n < len(v) {
		b := [1]byte{}
		if _, err = c.r.Read(b[:]); err != nil {
			return
		}

		switch b[0] {
		case 0xFF:
			io.ReadFull(c.r, c.buf[:])
		case 0x00:
			if _, err = c.r.Read(b[:]); err != nil {
				return
			}
			zeros := min(int(b[0]), len(v)-n)
			for i := range v[n : n+zeros] {
				v[i] = 0
			}
			c.zeros = int(b[0]) - zeros
			n += zeros
		default:
			ones := 0
			buf := [8]byte{}
			for i := 0; i < 8; i++ {
				if (b[0] & (1 << uint(i))) != 0 {
					ones++
				}
			}

			io.ReadFull(c.r, buf[:ones])

			for i, j := 0, 0; i < 8; i++ {
				if (b[0] & (1 << uint(i))) != 0 {
					c.buf[i] = buf[j]
					j++
				} else {
					c.buf[i] = 0
				}
			}

			use := copy(v[n:], c.buf[:])
			n += use
			c.bufsz = 8 - use
		}
	}

	return
}

func (c *Compressor) Write(v []byte) (n int, err error) {
	buf := []byte{}

	if (len(v) % 8) != 0 {
		return 0, errors.New("capnproto: compressor relies on word aligned data")
	}

	for n < len(v) {
		var hdr byte
		for i, b := range v[n : n+8] {
			if b != 0 {
				hdr |= 1 << uint(i)
			}
		}

		buf = append(buf, hdr)

		switch hdr {
		case 0x00:
			n += 8

			i := 0
			for n < len(v) && little64(v[n:]) == 0 && i < 0xFF {
				i++
				n += 8
			}

			buf = append(buf, byte(i))

		case 0xFF:
			buf = append(buf, v[n:n+8]...)
			n += 8

			i := n
			for i < min(len(v), n+0xFF*8) {
				zeros := 0
				for _, b := range v[i : i+8] {
					if b == 0 {
						zeros++
					}
				}

				if zeros < 7 {
					break
				}

				i += 8
			}

			buf = append(buf, byte((i-n)/8))
			buf = append(buf, v[n:i]...)
			n = i

		default:
			for _, b := range v[n : n+8] {
				if b != 0 {
					buf = append(buf, b)
				}
			}
			v = v[8:]
			n += 8
		}

	}

	if w, err := c.w.Write(buf); err != nil {
		return 0, err
	} else if w != len(buf) {
		return 0, io.ErrShortWrite
	}

	return n, nil
}
