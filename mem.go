package capn

import (
	"errors"
	"io"
	"math"
)

var (
	errBufferCall     = errors.New("capn: can't call on a memory buffer")
	ErrInvalidSegment = errors.New("capn: invalid segment id")
	ErrTooMuchData    = errors.New("capn: too much data in stream")
)

type buffer Segment

func NewBuffer(data []byte) *Segment {
	if uint64(len(data)) > uint64(math.MaxUint32) {
		return nil
	}

	b := &buffer{}
	b.Session = b
	b.Data = data
	return (*Segment)(b)
}

func (b *buffer) NewSegment(minsz int) (*Segment, error) {
	if uint64(len(b.Data)) > uint64(math.MaxUint32-minsz) {
		return nil, ErrOverlarge
	}
	b.Data = append(b.Data, make([]byte, minsz)...)
	b.Data = b.Data[:len(b.Data)-minsz]
	return (*Segment)(b), nil
}

func (b *buffer) Lookup(segid uint32) (*Segment, error) {
	if segid == 0 {
		return (*Segment)(b), nil
	} else {
		return nil, ErrInvalidSegment
	}
}

type multiBuffer struct {
	segments []*Segment
}

func NewMultiBuffer(data [][]byte) *Segment {
	m := &multiBuffer{make([]*Segment, len(data))}
	for i, d := range data {
		m.segments[i] = &Segment{m, d, uint32(i)}
	}
	if len(data) > 0 {
		return m.segments[0]
	}
	return &Segment{m, nil, 0xFFFFFFFF}
}

var (
	MaxSegmentNumber = 1024
	MaxTotalSize     = 1024 * 1024 * 1024
)

func ReadFromStream(r io.Reader) (*Segment, error) {
	var segnumv [4]byte
	if _, err := io.ReadFull(r, segnumv[:]); err != nil {
		return nil, err
	}

	if little32(segnumv[:]) >= uint32(MaxSegmentNumber) {
		return nil, ErrTooMuchData
	}

	segnum := int(little32(segnumv[:]) + 1)
	hdrsv := make([]byte, 8*(segnum/2)+4) // also include the padding

	if _, err := io.ReadFull(r, hdrsv[:]); err != nil {
		return nil, err
	}

	total := 0
	for i := 0; i < segnum; i++ {
		sz := little32(hdrsv[4*i:])
		if uint64(total)+uint64(sz)*8 > uint64(MaxTotalSize) {
			return nil, ErrTooMuchData
		}
		total += int(sz) * 8
	}

	datav := make([]byte, total)
	if _, err := io.ReadFull(r, datav[:]); err != nil {
		return nil, err
	}

	m := &multiBuffer{make([]*Segment, segnum)}
	for i := 0; i < segnum; i++ {
		sz := int(little32(hdrsv[4*i:])) * 8
		m.segments[i] = &Segment{m, datav[:sz], uint32(i)}
		datav = datav[sz:]
	}

	return m.segments[0], nil
}

func (m *multiBuffer) NewSegment(minsz int) (*Segment, error) {
	for _, s := range m.segments {
		if len(s.Data)+minsz <= cap(s.Data) {
			return s, nil
		}
	}

	if minsz < 4096 {
		minsz = 4096
	}
	s := &Segment{m, make([]byte, 0, minsz), uint32(len(m.segments))}
	m.segments = append(m.segments, s)
	return s, nil
}

func (m *multiBuffer) Lookup(segid uint32) (*Segment, error) {
	if uint(segid) < uint(len(m.segments)) {
		return m.segments[segid], nil
	} else {
		return nil, ErrInvalidSegment
	}
}
