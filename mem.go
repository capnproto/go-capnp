package capnp

import (
	"encoding/binary"
	"errors"
	"io"
	"math"

	"zombiezen.com/go/capnproto/internal/packed"
)

// A Message is a tree of Cap'n Proto objects, split into one or more
// segments of contiguous memory.  The only required field is Arena.
type Message struct {
	Arena Arena

	// CapTable is the indexed list of the clients referenced in the
	// message.  Capability pointers inside the message will use this table
	// to map pointers to Clients.  The table is usually populated by the
	// RPC system.
	//
	// See https://capnproto.org/encoding.html#capabilities-interfaces for
	// more details on the capability table.
	CapTable []Client

	segs []*Segment
}

// NewMessage creates a message with a new root and returns the first
// segment.  It is an error to call NewMessage on an arena with data in it.
func NewMessage(arena Arena) (msg *Message, first *Segment, err error) {
	msg = &Message{Arena: arena}
	switch arena.NumSegments() {
	case 0:
		first, err = msg.segmentForAlloc(defaultBufferSize)
		if err != nil {
			return nil, nil, err
		}
	case 1:
		first, err = msg.Segment(0)
		if err != nil {
			return nil, nil, err
		}
		if len(first.data) > 0 {
			return nil, nil, errHasData
		}
		if !hasCapacity(first.data, wordSize) {
			return nil, nil, errSegmentTooSmall
		}
	default:
		return nil, nil, errHasData
	}
	alloc(first, wordSize) // allocate root
	return msg, first, nil
}

// Root returns the pointer to the message's root object.
func (m *Message) Root() (Pointer, error) {
	s, err := m.Segment(0)
	if err != nil {
		return nil, err
	}
	return s.root().At(0)
}

// SetRoot sets the message's root object to p.
func (m *Message) SetRoot(p Pointer) error {
	s, err := m.Segment(0)
	if err != nil {
		return err
	}
	return s.root().Set(0, p)
}

// AddCap appends a capability to the message's capability table and
// returns its ID.
func (m *Message) AddCap(c Client) CapabilityID {
	n := CapabilityID(len(m.CapTable))
	m.CapTable = append(m.CapTable, c)
	return n
}

// NumSegments returns the number of segments in the message.
func (m *Message) NumSegments() int {
	return m.Arena.NumSegments()
}

// Segment returns the segment with the given ID.
func (m *Message) Segment(id SegmentID) (*Segment, error) {
	if isInt32Bit() && id > maxInt32 {
		return nil, errSegment32Bit
	}
	if seg := m.segment(id); seg != nil {
		return seg, nil
	}
	if int(id) >= m.Arena.NumSegments() {
		return nil, errSegmentOutOfBounds
	}
	data, err := m.Arena.Data(id)
	if err != nil {
		return nil, err
	}
	return m.setSegment(id, data), nil
}

func (m *Message) segment(id SegmentID) *Segment {
	if uint32(id) >= uint32(len(m.segs)) {
		return nil
	}
	return m.segs[id]
}

func (m *Message) setSegment(id SegmentID, data []byte) *Segment {
	if seg := m.segment(id); seg != nil {
		seg.data = data
		return seg
	}
	for uint32(len(m.segs)) <= uint32(id) {
		m.segs = append(m.segs, nil)
	}
	seg := &Segment{
		id:   id,
		msg:  m,
		data: data,
	}
	m.segs[id] = seg
	return seg
}

// segmentForAlloc creates or resizes an existing segment such that
// cap(seg.Data) - len(seg.Data) >= sz.
func (m *Message) segmentForAlloc(sz Size) (*Segment, error) {
	for _, seg := range m.segs {
		if seg == nil {
			continue
		}
		if hasCapacity(seg.data, sz) {
			return seg, nil
		}
	}
	id, data, err := m.Arena.Allocate(sz, m.segs)
	if err != nil {
		return nil, err
	}
	if isInt32Bit() && id > maxInt32 {
		return nil, errSegment32Bit
	}
	return m.setSegment(id, data), nil
}

// alloc allocates sz zero-filled bytes.  It prefers using s, but may
// use a different segment in the same message if there's not sufficient
// capacity.
func alloc(s *Segment, sz Size) (*Segment, Address, error) {
	sz = sz.padToWord()
	if sz > Size(math.MaxUint32)-wordSize {
		return nil, 0, errOverlarge
	}

	if !hasCapacity(s.data, sz) {
		// If we can't fit the data in the current segment, we always
		// return a far pointer to a tag in the new segment.
		var err error
		s, err = s.msg.segmentForAlloc(sz)
		if err != nil {
			return nil, 0, err
		}
	}

	addr := Address(len(s.data))
	end := addr.addSize(sz)
	s.data = s.data[:end]
	for i := addr; i < end; i++ {
		s.data[i] = 0
	}
	return s, addr, nil
}

// An Arena loads and allocates segments for a Message.  Segment IDs
// must be tightly packed in the range [0, NumSegments()).
type Arena interface {
	// NumSegments returns the number of segments in the arena.
	NumSegments() int

	// Data loads the data for the segment with the given ID.
	Data(id SegmentID) ([]byte, error)

	// Allocate allocates a byte slice such that cap(data) - len(data) >= minsz.
	// segs is a sparse list of already loaded segments, indexed by ID.
	// The arena may return an existing segment's ID, in which case the
	// arena is responsible for copying the data to the new byte slice.
	// Allocate must not modify the segments passed into it.
	Allocate(minsz Size, segs []*Segment) (SegmentID, []byte, error)
}

// Arena parameters.
const (
	defaultBufferSize      = 4096
	minSingleSegmentGrowth = 4096
)

type singleSegmentArena []byte

// SingleSegment returns a new arena with an expanding single-segment
// buffer.  b can be used to populate the buffer for reading or to
// reserve memory of a specific size.
func SingleSegment(b []byte) Arena {
	if cap(b) == 0 {
		b = make([]byte, 0, defaultBufferSize)
	}
	ssa := new(singleSegmentArena)
	*ssa = b
	return ssa
}

func (ssa *singleSegmentArena) NumSegments() int {
	return 1
}

func (ssa *singleSegmentArena) Data(id SegmentID) ([]byte, error) {
	if id != 0 {
		return nil, errSegmentOutOfBounds
	}
	return *ssa, nil
}

func (ssa *singleSegmentArena) Allocate(sz Size, segs []*Segment) (SegmentID, []byte, error) {
	if (len(segs) == 0 || segs[0] == nil) && hasCapacity(*ssa, sz) {
		return 0, *ssa, nil
	}
	// TODO(light): ensure len(data)+sz is word-aligned
	if sz < minSingleSegmentGrowth {
		sz = minSingleSegmentGrowth
	} else {
		sz = sz.padToWord()
	}
	buf := make([]byte, len(segs[0].data), cap(*ssa)+int(sz))
	copy(buf, *ssa)
	*ssa = buf
	return 0, *ssa, nil
}

type multiSegmentArena [][]byte

// MultiSegment returns a new arena that allocates new segments when
// they are full.  b can be used to populate the buffer for reading or
// to reserve memory of a specific size.
func MultiSegment(b [][]byte) Arena {
	msa := new(multiSegmentArena)
	*msa = b
	return msa
}

// demuxArena slices b into a multi-segment arena.
func demuxArena(sizes []Size, data []byte) Arena {
	segs := make([][]byte, len(sizes))
	for i, sz := range sizes {
		segs[i], data = data[:sz:sz], data[sz:]
	}
	return MultiSegment(segs)
}

func (msa *multiSegmentArena) NumSegments() int {
	return len(*msa)
}

func (msa *multiSegmentArena) Data(id SegmentID) ([]byte, error) {
	if int64(id) >= int64(len(*msa)) {
		return nil, errSegmentOutOfBounds
	}
	return (*msa)[id], nil
}

func (msa *multiSegmentArena) Allocate(sz Size, segs []*Segment) (SegmentID, []byte, error) {
	for i, data := range *msa {
		if i < len(segs) && segs[i] != nil {
			// The message would have already considered this segment.
			continue
		}
		if hasCapacity(data, sz) {
			return SegmentID(i), data, nil
		}
	}
	if sz < defaultBufferSize {
		sz = defaultBufferSize
	} else {
		sz = sz.padToWord()
	}
	buf := make([]byte, 0, int(sz))
	id := SegmentID(len(*msa))
	*msa = append(*msa, buf)
	return id, buf, nil
}

// A Decoder represents a framer that deserializes a particular Cap'n
// Proto input stream.
type Decoder struct {
	r io.Reader
}

// NewDecoder creates a new Cap'n Proto framer that reads from r.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{r: r}
}

// NewPackedDecoder creates a new Cap'n Proto framer that reads from a
// packed stream r.
func NewPackedDecoder(r io.Reader) *Decoder {
	return NewDecoder(packed.NewReader(r))
}

// Decode reads a message from the decoder stream.
func (d *Decoder) Decode() (*Message, error) {
	var maxSegBuf [msgHeaderSize]byte
	if _, err := io.ReadFull(d.r, maxSegBuf[:]); err != nil {
		return nil, err
	}
	maxSeg := binary.LittleEndian.Uint32(maxSegBuf[:])
	hdrSize := streamHeaderSize(maxSeg)
	hdr := make([]byte, hdrSize)
	copy(hdr, maxSegBuf[:])
	if _, err := io.ReadFull(d.r, hdr[msgHeaderSize:]); err != nil {
		return nil, err
	}
	sizes, _, err := unmarshalStreamHeader(hdr)
	if err != nil {
		return nil, err
	}
	var total int64
	for _, sz := range sizes {
		total += int64(sz)
	}
	// TODO(light): size check
	buf := make([]byte, int(total))
	if _, err := io.ReadFull(d.r, buf); err != nil {
		return nil, err
	}
	return &Message{Arena: demuxArena(sizes, buf)}, nil
}

// Unmarshal reads an unpacked serialized stream into a message.  No
// copying is performed, so the objects in the returned message read
// directly from data.
func Unmarshal(data []byte) (*Message, error) {
	sizes, data, err := unmarshalStreamHeader(data)
	if err != nil {
		return nil, err
	}
	// TODO(light): check len(data) to size
	return &Message{Arena: demuxArena(sizes, data)}, nil
}

// An Encoder represents a framer for serializing a particular Cap'n
// Proto stream.
type Encoder struct {
	w      io.Writer
	hdrbuf []byte

	packed  bool
	packbuf []byte
}

// NewEncoder creates a new Cap'n Proto framer that writes to w.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

// NewPackedEncoder creates a new Cap'n Proto framer that writes to a
// packed stream w.
func NewPackedEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w, packed: true}
}

// Encode writes a message to the encoder stream.
func (e *Encoder) Encode(m *Message) error {
	// TODO(light): Lazily load from arena, don't necessarily need to fit in memory.
	nsegs := m.NumSegments()
	sizes, err := m.segmentSizes()
	if err != nil {
		return err
	}
	maxSeg := uint32(nsegs - 1)
	hdrSize := streamHeaderSize(maxSeg)
	if cap(e.hdrbuf) < hdrSize {
		e.hdrbuf = make([]byte, hdrSize)
	} else {
		e.hdrbuf = e.hdrbuf[:hdrSize]
	}
	marshalStreamHeader(e.hdrbuf, sizes)
	if err := e.write(e.hdrbuf); err != nil {
		return err
	}
	for i := 0; i < nsegs; i++ {
		s, err := m.Segment(SegmentID(i))
		if err != nil {
			return err
		}
		if err := e.write(s.data); err != nil {
			return err
		}
	}
	return nil
}

func (e *Encoder) write(b []byte) error {
	if e.packed {
		e.packbuf = packed.Pack(e.packbuf[:0], b)
		b = e.packbuf
	}
	_, err := e.w.Write(b)
	return err
}

func (m *Message) segmentSizes() ([]Size, error) {
	nsegs := m.NumSegments()
	sizes := make([]Size, nsegs)
	for i := 0; i < nsegs; i++ {
		s, err := m.Segment(SegmentID(i))
		if err != nil {
			return sizes[:i], err
		}
		n := len(s.data)
		if int64(n) > int64(maxSize) {
			return sizes[:i], errOverlarge
		}
		sizes[i] = Size(n)
	}
	return sizes, nil
}

// Marshal concatenates the segments in the message into a single byte
// slice including framing.
func (m *Message) Marshal() ([]byte, error) {
	// Compute buffer size.
	// TODO(light): error out if too many segments
	nsegs := m.NumSegments()
	maxSeg := uint32(nsegs - 1)
	hdrSize := streamHeaderSize(maxSeg)
	sizes, err := m.segmentSizes()
	if err != nil {
		return nil, err
	}
	total := hdrSize
	for _, sz := range sizes {
		// TODO(light): error out if too large
		total += int(sz)
	}

	// Fill in buffer.
	buf := make([]byte, hdrSize, total)
	marshalStreamHeader(buf, sizes)
	for i := 0; i < nsegs; i++ {
		s, err := m.Segment(SegmentID(i))
		if err != nil {
			return nil, err
		}
		buf = append(buf, s.data...)
	}
	return buf, nil
}

// MarshalPacked marshals the message in packed form.
func (m *Message) MarshalPacked() ([]byte, error) {
	data, err := m.Marshal()
	if err != nil {
		return nil, err
	}
	buf := make([]byte, 0, len(data))
	buf = packed.Pack(buf, data)
	return buf, nil
}

// Stream header sizes.
const (
	msgHeaderSize = 4
	segHeaderSize = 4
)

// streamHeaderSize returns the size of the header, given the
// first 32-bit number.
func streamHeaderSize(n uint32) int {
	return (msgHeaderSize + segHeaderSize*(int(n)+1) + 7) &^ 7
}

// marshalStreamHeader marshals the sizes into the byte slice, which
// must be of size streamHeaderSize(len(sizes) - 1).
func marshalStreamHeader(b []byte, sizes []Size) {
	binary.LittleEndian.PutUint32(b, uint32(len(sizes)-1))
	for i, sz := range sizes {
		loc := msgHeaderSize + i*segHeaderSize
		binary.LittleEndian.PutUint32(b[loc:], uint32(sz/Size(wordSize)))
	}
}

// unmarshalStreamHeader parses the header of the stream framing format.
func unmarshalStreamHeader(data []byte) (sizes []Size, tail []byte, err error) {
	if len(data) < streamHeaderSize(0) {
		return nil, nil, io.ErrUnexpectedEOF
	}
	maxSeg := binary.LittleEndian.Uint32(data)
	// TODO(light): check int
	hdrSize := streamHeaderSize(maxSeg)
	if len(data) < hdrSize {
		return nil, nil, io.ErrUnexpectedEOF
	}
	n := int(maxSeg + 1)
	sizes = make([]Size, n)
	for i := 0; i < n; i++ {
		s := binary.LittleEndian.Uint32(data[msgHeaderSize+i*segHeaderSize:])
		sizes[i] = wordSize.times(int32(s))
	}
	return sizes, data[hdrSize:], nil
}

func hasCapacity(b []byte, sz Size) bool {
	return sz <= Size(cap(b)-len(b))
}

const maxInt32 = 0x7fffffff

// isInt32Bit reports whether the built-in int type is 32 bits.
func isInt32Bit() bool {
	const maxInt = int(^uint(0) >> 1)
	return maxInt == maxInt32
}

var (
	errBufferCall         = errors.New("capn: can't call on a memory buffer")
	errSegmentOutOfBounds = errors.New("capn: segment ID out of bounds")
	errSegment32Bit       = errors.New("capn: segment ID larger than 31 bits")
	errHasData            = errors.New("capnp: NewMessage called on arena with data")
	errTooMuchData        = errors.New("capn: too much data in stream")
	errSegmentTooSmall    = errors.New("capn: segment too small")
	errStreamHeader       = errors.New("capnp: invalid stream header")
)
