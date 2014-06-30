package capn

import (
	"bufio"
	"bytes"
	"encoding/binary"
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

	// The 1-3 states are for dealing with the 0xFF tag and the raw bytes that follow.
	// They tell us where to pick up if we are interrupted in the middle of anything
	// after the 0xFF tag, until we are done with the raw read.
	S_POSTFF = 1
	S_READN  = 2
	S_RAW    = 3
)

type Decompressor struct {
	r     io.Reader
	buf   [8]byte
	bufsz int

	// track the bytes after a 0xff raw tag
	ffBuf          [8]byte
	ffBufLoadCount int // count of bytes loaded from r into ffBuf (max 8)
	ffBufUsedCount int // count of bytes supplied to v during Read().

	zeros int
	raw   int // number of raw bytes left to copy through
	state DecompParseState
}

// externally available flag for compiling with debug info on/off
const VerboseDecomp = false
const VerboseCompress = false

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
	var bytesRead int

	for {
		if len(v) == 0 {
			return
		}

		switch c.state {

		case S_RAW:
			if c.raw > 0 {
				bytesRead, err = c.r.Read(v[:min(len(v), c.raw)])
				if VerboseDecomp {
					fmt.Printf("decompression copied in %d raw bytes: %v\n", bytesRead, v[:bytesRead])
				}
				c.raw -= bytesRead
				v = v[bytesRead:]
				n += bytesRead

				if err != nil {
					return
				}
			} else {
				c.state = S_NORMAL
			}

		case S_POSTFF:
			if c.ffBufUsedCount >= 8 {
				c.state = S_READN
				continue
			}
			// invar: c.ffBufUsedCount < 8

			// before reading more from r, first empty any residual in buffer. Such
			// bytes were already read from r, are now
			// waiting in c.ffBuf, and have not yet been given to v: so
			// these bytes are first in line to go.
			if c.ffBufUsedCount < c.ffBufLoadCount {
				br := copy(v, c.ffBuf[c.ffBufUsedCount:c.ffBufLoadCount])
				if VerboseDecomp {
					fmt.Printf("decompression copied in %d bytes: %v\n", br, v[:br])
				}
				c.ffBufUsedCount += br
				v = v[br:]
				n += br
			}
			if c.ffBufUsedCount >= 8 {
				c.state = S_READN
				continue
			}
			// invar: c.ffBufUsedCount < 8

			// io.ReadFull, try to read exactly (8 - cc.ffBufLoadCount) bytes
			// io.ReadFull returns EOF only if no bytes were read
			if c.ffBufLoadCount < 8 {
				bytesRead, err = io.ReadFull(c.r, c.ffBuf[c.ffBufLoadCount:]) // read up to 8 bytes into c.buf
				if bytesRead > 0 {
					c.ffBufLoadCount += bytesRead
				} else {
					return
				}
				if err != nil {
					return
				}
			}
			// stay in S_POSTFF

		case S_READN:
			if bytesRead, err = c.r.Read(b[:]); err != nil {
				return
			}
			if bytesRead == 0 {
				return
			}
			c.raw = int(b[0]) * 8
			c.state = S_RAW

		case S_NORMAL:

			if c.zeros > 0 {
				num0 := min(len(v), c.zeros)
				x := v[:num0]
				for i := range x {
					x[i] = 0
				}
				c.zeros -= num0
				n += num0
				if c.zeros > 0 {
					return n, nil
				}
				v = v[num0:]
				if len(v) == 0 {
					return n, nil
				}
			}
			// INVAR: c.zeros == 0

			if c.bufsz > 0 {
				nc := copy(v, c.buf[8-c.bufsz:])
				c.bufsz -= nc
				n += nc
				v = v[nc:]
				if c.bufsz > 0 {
					return n, nil
				}
			}
			// INVAR: c.bufz == 0

			for c.state == S_NORMAL && len(v) > 0 {

				if _, err = c.r.Read(b[:]); err != nil {
					return
				}
				if VerboseDecomp {
					fmt.Printf("decompression read TAG byte b: %#v\n", b)
				}

				switch b[0] {
				case 0xFF:
					c.ffBufLoadCount = 0
					c.ffBufUsedCount = 0
					c.state = S_POSTFF

					break

				case 0x00:
					if _, err = c.r.Read(b[:]); err != nil {
						return
					}
					if VerboseDecomp {
						fmt.Printf("decompression read byte Zero-word -1 count byte b: %#v\n", b)
					}

					requestedZeroBytes := (int(b[0]) + 1) * 8
					zeros := min(requestedZeroBytes, len(v))

					if VerboseDecomp {
						fmt.Printf("decompression writing zeros to n=%d to n+zeros=%d  &v[0]=%p\n", n, n+zeros, &v[0])
					} // this next is obliterating out v[4] wierdly
					for i := 0; i < zeros; i++ {
						v[i] = 0
					}
					v = v[zeros:]
					n += zeros
					// remember the leftover zeros to write
					c.zeros = requestedZeroBytes - zeros

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

					use := copy(v, c.buf[:])
					if VerboseDecomp {
						fmt.Printf("decompression copied in %d bytes: %v\n", use, c.buf[:])
					}
					v = v[use:]
					n += use
					c.bufsz = 8 - use
				}
			}
		}

	}
	return
}

func (c *Compressor) Write(v []byte) (n int, err error) {
	origVlen := len(v)
	if (origVlen % 8) != 0 {
		return 0, errors.New("capnproto: compressor relies on word aligned data")
	}
	buf := make([]byte, 0, 8)
	for len(v) > 0 {
		var hdr byte
		buf = buf[:0]
		for i, b := range v[:8] {
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
		v = v[8:]

		switch hdr {
		case 0x00:
			i := 0
			for len(v) > 0 && binary.LittleEndian.Uint64(v) == 0 && i < 0xFF {
				i++
				n += 8
				v = v[8:]
			}
			err = c.w.WriteByte(byte(i))
			if err != nil {
				return n, err
			}
		case 0xFF:
			i := 0
			end := min(len(v), 0xFF*8)
			for i < end {
				zeros := 0
				for _, b := range v[i : i+8] {
					if b == 0 {
						zeros++
					}
				}

				if zeros > 1 {
					break
				}
				i += 8
			}

			rawWords := byte(i / 8)
			err := c.w.WriteByte(rawWords)
			if err != nil {
				return n, err
			}
			_, err = c.w.Write(v[:i])
			if err != nil {
				return n, err
			}
			n += i
			v = v[i:]
		}
	}
	err = c.w.Flush()
	return n, err
}
