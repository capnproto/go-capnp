// Package packed provides functions to read and write the "packed"
// compression scheme described at https://capnproto.org/encoding.html#packing.
package packed

import "io"

const wordSize = 8

// Special case tags.
const (
	zeroTag     byte = 0x00
	unpackedTag byte = 0xff
)

// Pack appends the packed version of src to dst and returns the
// resulting slice.  len(src) must be a multiple of 8 or Pack panics.
func Pack(dst, src []byte) []byte {
	if len(src)%wordSize != 0 {
		panic("packed.Pack len(src) must be a multiple of 8")
	}
	var buf [wordSize]byte
	for len(src) > 0 {
		var hdr byte
		n := 0
		for i := uint(0); i < wordSize; i++ {
			if src[i] != 0 {
				hdr |= 1 << i
				buf[n] = src[i]
				n++
			}
		}
		dst = append(dst, hdr)
		dst = append(dst, buf[:n]...)
		src = src[wordSize:]

		switch hdr {
		case zeroTag:
			z := min(numZeroWords(src), 0xff)
			dst = append(dst, byte(z))
			src = src[z*wordSize:]
		case unpackedTag:
			i := 0
			end := min(len(src), 0xff*wordSize)
			for i < end {
				zeros := 0
				for _, b := range src[i : i+wordSize] {
					if b == 0 {
						zeros++
					}
				}

				if zeros > 1 {
					break
				}
				i += wordSize
			}

			rawWords := byte(i / wordSize)
			dst = append(dst, rawWords)
			dst = append(dst, src[:i]...)
			src = src[i:]
		}
	}
	return dst
}

// numZeroWords returns the number of leading zero words in b.
func numZeroWords(b []byte) int {
	for i, bb := range b {
		if bb != 0 {
			return i / wordSize
		}
	}
	return len(b) / wordSize
}

type decompressor struct {
	r     io.Reader
	buf   [wordSize]byte
	bufsz int

	// track the bytes after a 0xff raw tag
	ffBuf          [wordSize]byte
	ffBufLoadCount int // count of bytes loaded from r into ffBuf (max wordSize)
	ffBufUsedCount int // count of bytes supplied to v during Read().

	zeros int
	raw   int // number of raw bytes left to copy through
	state decompressorState
}

// NewReader returns a reader that decompresses a packed stream from r.
func NewReader(r io.Reader) io.Reader {
	return &decompressor{r: r}
}

func min(a, b int) int {
	if b < a {
		return b
	}
	return a
}

func (c *decompressor) Read(v []byte) (n int, err error) {

	var b [1]byte
	var bytesRead int

	for {
		if len(v) == 0 {
			return
		}

		switch c.state {

		case rawState:
			if c.raw > 0 {
				bytesRead, err = c.r.Read(v[:min(len(v), c.raw)])
				c.raw -= bytesRead
				v = v[bytesRead:]
				n += bytesRead

				if err != nil {
					return
				}
			} else {
				c.state = normalState
			}

		case postFFState:
			if c.ffBufUsedCount >= wordSize {
				c.state = readnState
				continue
			}
			// invar: c.ffBufUsedCount < wordSize

			// before reading more from r, first empty any residual in buffer. Such
			// bytes were already read from r, are now
			// waiting in c.ffBuf, and have not yet been given to v: so
			// these bytes are first in line to go.
			if c.ffBufUsedCount < c.ffBufLoadCount {
				br := copy(v, c.ffBuf[c.ffBufUsedCount:c.ffBufLoadCount])
				c.ffBufUsedCount += br
				v = v[br:]
				n += br
			}
			if c.ffBufUsedCount >= wordSize {
				c.state = readnState
				continue
			}
			// invar: c.ffBufUsedCount < wordSize

			// io.ReadFull, try to read exactly (wordSize - cc.ffBufLoadCount) bytes
			// io.ReadFull returns EOF only if no bytes were read
			if c.ffBufLoadCount < wordSize {
				bytesRead, err = io.ReadFull(c.r, c.ffBuf[c.ffBufLoadCount:]) // read up to wordSize bytes into c.buf
				if bytesRead > 0 {
					c.ffBufLoadCount += bytesRead
				} else {
					return
				}
				if err != nil {
					return
				}
			}
			// stay in postFFState

		case readnState:
			if bytesRead, err = c.r.Read(b[:]); err != nil {
				return
			}
			if bytesRead == 0 {
				return
			}
			c.raw = int(b[0]) * wordSize
			c.state = rawState

		case normalState:

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
				nc := copy(v, c.buf[wordSize-c.bufsz:])
				c.bufsz -= nc
				n += nc
				v = v[nc:]
				if c.bufsz > 0 {
					return n, nil
				}
			}
			// INVAR: c.bufz == 0

			for c.state == normalState && len(v) > 0 {

				if _, err = c.r.Read(b[:]); err != nil {
					return
				}

				switch b[0] {
				case unpackedTag:
					c.ffBufLoadCount = 0
					c.ffBufUsedCount = 0
					c.state = postFFState

					break

				case zeroTag:
					if _, err = c.r.Read(b[:]); err != nil {
						return
					}

					requestedZeroBytes := (int(b[0]) + 1) * wordSize
					zeros := min(requestedZeroBytes, len(v))

					for i := 0; i < zeros; i++ {
						v[i] = 0
					}
					v = v[zeros:]
					n += zeros
					// remember the leftover zeros to write
					c.zeros = requestedZeroBytes - zeros

				default:
					ones := 0
					var buf [wordSize]byte
					for i := 0; i < wordSize; i++ {
						if (b[0] & (1 << uint(i))) != 0 {
							ones++
						}
					}

					_, err = io.ReadFull(c.r, buf[:ones])
					if err != nil {
						return
					}

					for i, j := 0, 0; i < wordSize; i++ {
						if (b[0] & (1 << uint(i))) != 0 {
							c.buf[i] = buf[j]
							j++
						} else {
							c.buf[i] = 0
						}
					}

					use := copy(v, c.buf[:])
					v = v[use:]
					n += use
					c.bufsz = wordSize - use
				}
			}
		}

	}
	return
}

// decompressorState is the state of a decompressor.
type decompressorState uint8

// Decompressor states
const (
	normalState decompressorState = iota

	// These states are for dealing with the 0xFF tag and the raw bytes that follow.
	// They tell us where to pick up if we are interrupted in the middle of anything
	// after the 0xFF tag, until we are done with the raw read.
	postFFState
	readnState
	rawState
)
