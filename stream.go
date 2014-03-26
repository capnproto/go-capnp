package capn

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
)

type Compressor struct {
	w *bufio.Writer
}

type DecompParseState uint8

const (
	S_NORMAL DecompParseState = 0
	S_POSTFF                  = 1
	S_READN                   = 2
	S_RAW                     = 3
)

type Decompressor struct {
	r     io.Reader
	buf   [8]byte
	bufsz int
	zeros int
	raw   int // number of raw bytes left to copy through
	preN  int // number of bytes that we have left to read, after the 0xFF and before the N. (preN has max of 8)
	state DecompParseState
}

func NewCompressor(w io.Writer) *Compressor {
	return &Compressor{
		w: bufio.NewWriter(w),
	}
}

// WriteToPacked writes the message that the segment is part of to the
// provided stream in packed form.
func (s *Segment) WriteToPacked(w io.Writer) (int64, error) {
	c := NewCompressor(w)
	return s.WriteTo(c)
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

	var b [1]byte

	for {
		switch c.state {

		case S_RAW:
			if c.raw > 0 {
				n, err = c.r.Read(v[:min(len(v), c.raw)])
				c.raw -= n
				if c.raw == 0 {
					c.state = S_NORMAL
				}
				return
			}

		case S_POSTFF:
			bytesRead, err = io.ReadFull(c.r, c.buf[:c.PreN]) // read 8 bytes
			if bytesRead > 0 {
				br := copy(p, c.buf[:c.PreN])
				n += br
			}
			c.preN -= bytesRead
			if c.preN == 0 {
				c.state = S_READN
			}

			if err != nil {
				return
				// panic("Decompression error: truncated stream after 0xFF tag: we could not read 8 bytes after 0xFF")
			}

		case S_READN:
			if _, err = c.r.Read(b[:]); err != nil {
				return
			}
			c.raw = int(b) * 8
			c.state = S_RAW

		case S_NORMAL:

			if c.zeros > 0 {
				n = min(len(v), c.zeros)
				x := v[:n]
				for i := range x {
					x[i] = 0
				}
				c.zeros -= n
				return
			}

			//	if len(v) > 5 {
			//		fmt.Printf("address of v[4] is: %p\n", &v[4])
			//	}

			if c.bufsz > 0 {
				n = copy(v, c.buf[8-c.bufsz:])
				c.bufsz -= n
			}

			for c.state == S_NORMAL && n < len(v) {

				if _, err = c.r.Read(b[:]); err != nil {
					return
				}
				fmt.Printf("decompression read byte TAG byte b: %#v\n", b)

				switch b[0] {
				case 0xFF:
					c.preN = 8
					c.state = S_POSTFF
					break

				case 0x00:
					if _, err = c.r.Read(b[:]); err != nil {
						return
					}
					fmt.Printf("decompression read byte Zero-word -1 count byte b: %#v\n", b)

					requestedZeroBytes := (int(b[0]) + 1) * 8
					zeros := min(requestedZeroBytes, len(v)-n)

					fmt.Printf("decompression writing zeros to n=%d to n+zeros=%d  &v[0]=%p\n", n, n+zeros, &v[0]) // this next is obliterating out v[4] wierdly
					//for i := range v[n : n+zeros] { // this is a bug: the range is from 0...zeros, obliterating the already written stuff. not what we want.
					for i := n; i < n+zeros; i++ {
						v[i] = 0
					}
					c.zeros = requestedZeroBytes - zeros
					n += zeros
				default:
					ones := 0
					var buf [8]byte
					for i := 0; i < 8; i++ {
						if (b[0] & (1 << uint(i))) != 0 {
							ones++
						}
					}

					_, err = io.ReadFull(c.r, buf[:ones])
					if err != nil {
						return
					}

					for i, j := 0, 0; i < 8; i++ {
						if (b[0] & (1 << uint(i))) != 0 {
							c.buf[i] = buf[j]
							j++
						} else {
							c.buf[i] = 0
						}
					}

					use := copy(v[n:], c.buf[:])
					fmt.Printf("decompression copied in %d bytes: %v\n", use, c.buf[:])
					n += use
					c.bufsz = 8 - use
				}
			}
		}
	}
	return
}

func (c *Compressor) Write(v []byte) (n int, err error) {
	if (len(v) % 8) != 0 {
		return 0, errors.New("capnproto: compressor relies on word aligned data")
	}
	buf := make([]byte, 0, 8)
	for n < len(v) {
		var hdr byte
		buf = buf[:0]
		for i, b := range v[n : n+8] {
			if b != 0 {
				hdr |= 1 << uint(i)
				buf = append(buf, b)
			}
		}
		err = c.w.WriteByte(hdr)
		if err != nil {
			return n, err
		}
		_, err = c.w.Write(buf)
		if err != nil {
			return n, err
		}
		n += 8

		switch hdr {
		case 0x00:
			i := 0
			for n < len(v) && little64(v[n:]) == 0 && i < 0xFF {
				i++
				n += 8
			}
			err = c.w.WriteByte(byte(i))
			if err != nil {
				return n, err
			}
		case 0xFF:
			i := n
			end := min(len(v), n+0xFF*8)
			for i < end {
				zeros := 0
				for _, b := range v[i : i+8] {
					if b == 0 {
						zeros++
					}
				}

				if zeros > 7 {
					break
				}
				i += 8
			}

			rawWords := byte((i - n) / 8)
			err := c.w.WriteByte(rawWords)
			if err != nil {
				return n, err
			}
			_, err = c.w.Write(v[n:i])
			if err != nil {
				return n, err
			}
			n = i
		}
	}
	err = c.w.Flush()
	return n, err
}
