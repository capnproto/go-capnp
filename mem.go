package capn

import (
	"errors"
)

var (
	errBufferCall     = errors.New("capn: can't call on a memory buffer")
	ErrInvalidSegment = errors.New("capn: invalid segment id")
)

type buffer Segment

func (b *buffer) NewCall() (*Call, error) { return nil, errBufferCall }

func NewBuffer(data []byte) *Segment {
	b := &buffer{}
	b.Session = b
	b.Data = data
	return (*Segment)(b)
}

func (b *buffer) NewSegment(minsz int) (*Segment, error) {
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

func (b *multiBuffer) NewCall() (*Call, error) { return nil, errBufferCall }

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
