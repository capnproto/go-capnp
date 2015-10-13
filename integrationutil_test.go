package capnp_test

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"unsafe"

	"zombiezen.com/go/capnproto2"
	air "zombiezen.com/go/capnproto2/internal/aircraftlib"
)

const schemaPath = "internal/aircraftlib/aircraft.capnp"

func zdateFilledMessage(t testing.TB, n int32) *capnp.Message {
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatal(err)
	}
	z, err := air.NewRootZ(seg)
	if err != nil {
		t.Fatal(err)
	}
	list, err := air.NewZdate_List(seg, n)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < int(n); i++ {
		d, err := air.NewZdate(seg)
		if err != nil {
			t.Fatal(err)
		}
		d.SetMonth(12)
		d.SetDay(7)
		d.SetYear(int16(2004 + i))
		list.Set(i, d)
	}
	z.SetZdatevec(list)

	return msg
}

func zdataFilledMessage(t testing.TB, n int) *capnp.Message {
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatal(err)
	}
	z, err := air.NewRootZ(seg)
	if err != nil {
		t.Fatal(err)
	}
	d, err := air.NewZdata(seg)
	if err != nil {
		t.Fatal(err)
	}
	b := make([]byte, n)
	for i := 0; i < len(b); i++ {
		b[i] = byte(i)
	}
	d.SetData(b)
	z.SetZdata(d)
	return msg
}

// some generally useful capnp/segment utilities

// shell out to display capnp bytes as human-readable text. Data flow:
//    in-memory capn segment -> stdin to capnp decode -> stdout human-readble string form
func CapnpDecodeSegment(seg *capnp.Segment, capnpExePath string, capnpSchemaFilePath string, typeName string) string {

	// set defaults
	if capnpExePath == "" {
		capnpExePath = CheckAndGetCapnpPath()
	}

	if capnpSchemaFilePath == "" {
		capnpSchemaFilePath = schemaPath
	}

	if typeName == "" {
		typeName = "Z"
	}

	cs := []string{"decode", "--short", capnpSchemaFilePath, typeName}
	cmd := exec.Command(capnpExePath, cs...)
	cmdline := capnpExePath + " " + strings.Join(cs, " ")

	buf, err := seg.Message().Marshal()
	if err != nil {
		panic(err)
	}

	cmd.Stdin = bytes.NewReader(buf)

	var errout bytes.Buffer
	cmd.Stderr = &errout

	bs, err := cmd.Output()
	if err != nil {
		if err.Error() == "exit status 1" {
			cwd, _ := os.Getwd()
			fmt.Fprintf(os.Stderr, "\nCall to capnp in CapnpDecodeSegment(): '%s' in dir '%s' failed with status 1\n", cmdline, cwd)
			fmt.Printf("stderr: '%s'\n", string(errout.Bytes()))
			fmt.Printf("stdout: '%s'\n", string(bs))
		}
		panic(err)
	}
	return strings.TrimSpace(string(bs))
}

// disk file of a capn segment -> in-memory capn segment -> stdin to capnp decode -> stdout human-readble string form
func CapnFileToText(serializedCapnpFilePathToDisplay string, capnpSchemaFilePath string, capnpExePath string) (string, error) {

	// a) read file into Segment

	byteslice, err := ioutil.ReadFile(serializedCapnpFilePathToDisplay)
	if err != nil {
		return "", err
	}

	msg, err := capnp.Unmarshal(byteslice)

	if err != nil {
		return "", err
	}
	seg, err := msg.Segment(0)
	if err != nil {
		return "", err
	}

	// b) tell CapnpDecodeSegment() to show the human-readable-text form of the message
	// warning: CapnpDecodeSegment() may panic on you. It is a testing utility so that
	//  is desirable. For production, do something else.
	return CapnpDecodeSegment(seg, capnpExePath, capnpSchemaFilePath, "Z"), nil
}

// return path to capnp if 'which' can find it. Feel free to replace this with
//   a more general configuration mechanism.
func CheckAndGetCapnpPath() string {

	path, err := exec.LookPath("capnp")
	if err != nil {
		panic(fmt.Sprintf("could not locate the capnp executable: put the capnp executable in your path: %s", err))
	}

	cmd := exec.Command(path, "id")
	bs, err := cmd.Output()
	if err != nil || string(bs[:3]) != `@0x` {
		panic(fmt.Sprintf("%s id did not function: put a working capnp executable in your path. Err: %s", path, err))
	}

	return path
}

// take an already (packed or unpacked, depending on the packed flag) buffer of a serialized segment, and display it
func CapnpDecodeBuf(buf []byte, capnpExePath string, capnpSchemaFilePath string, typeName string, packed bool) string {

	// set defaults
	if capnpExePath == "" {
		capnpExePath = CheckAndGetCapnpPath()
	}

	if capnpSchemaFilePath == "" {
		capnpSchemaFilePath = schemaPath
	}

	if typeName == "" {
		typeName = "Z"
	}

	cs := []string{"decode", "--short", capnpSchemaFilePath, typeName}
	if packed {
		cs = []string{"decode", "--short", "--packed", capnpSchemaFilePath, typeName}
	}
	cmd := exec.Command(capnpExePath, cs...)
	cmdline := capnpExePath + " " + strings.Join(cs, " ")

	cmd.Stdin = bytes.NewReader(buf)

	var errout bytes.Buffer
	cmd.Stderr = &errout

	bs, err := cmd.Output()
	if err != nil {
		if err.Error() == "exit status 1" {
			cwd, _ := os.Getwd()
			fmt.Fprintf(os.Stderr, "\nCall to capnp in CapnpDecodeBuf(): '%s' in dir '%s' failed with status 1\n", cmdline, cwd)
			fmt.Printf("stderr: '%s'\n", string(errout.Bytes()))
			fmt.Printf("stdout: '%s'\n", string(bs))
		}
		panic(err)
	}
	return strings.TrimSpace(string(bs))
}

const zerohi32 uint64 = ^(^0 << 32)

func CapnpEncode(msg string, typ string) []byte {
	capnpPath, err := exec.LookPath("capnp")
	//capnpPath, err := exec.LookPath("tee")
	if err != nil {
		panic(err)
	}
	if !FileExists(capnpPath) {
		panic(fmt.Sprintf("could not locate capnp tool in PATH"))
	}

	schfn := schemaPath
	args := []string{"encode", schfn, typ}
	cmdline := fmt.Sprintf("%s %s %s %s", capnpPath, "encode", schfn, typ)
	//fmt.Printf("cmdline = %s\n", cmdline)
	c := exec.Command(capnpPath, args...)

	var o bytes.Buffer
	c.Stdout = &o

	var in bytes.Buffer
	in.Write([]byte(msg))
	c.Stdin = &in

	err = c.Run()
	if err != nil {
		panic(fmt.Errorf("tried to run %s, got err:%s", cmdline, err))
	}
	return o.Bytes()
}

func CapnpDecode(input []byte, typ string) []byte {
	capnpPath, err := exec.LookPath("capnp")
	//capnpPath, err := exec.LookPath("tee")
	if err != nil {
		panic(err)
	}
	if !FileExists(capnpPath) {
		panic(fmt.Sprintf("could not locate capnp tool in PATH"))
	}

	schfn := schemaPath
	args := []string{"decode", "--short", schfn, typ}
	cmdline := fmt.Sprintf("%s %s %s %s %s", capnpPath, "decode", "--short", schfn, typ)
	fmt.Printf("cmdline = %s\n", cmdline)
	c := exec.Command(capnpPath, args...)

	var o bytes.Buffer
	c.Stdout = &o

	var e bytes.Buffer
	c.Stderr = &e

	var in bytes.Buffer
	in.Write(input)
	c.Stdin = &in

	err = c.Run()
	if err != nil {
		fmt.Printf("tried to run %s, got err:%s and stderr: '%s'", cmdline, err, e.Bytes())
		panic(err)
	}
	return o.Bytes()
}

func ShowBytes(b []byte, indent int) {
	c := NewCap()
	k := 0
	ind := strings.Repeat(" ", indent)
	fmt.Printf("\n%s", ind)
	line := 0
	for i := 0; i < len(b)/8; i++ {
		for j := 0; j < 8; j++ {
			fmt.Printf("%02x ", b[k])
			k++
			if k == len(b) {
				break
			}
		}
		fmt.Printf("   ==(line %02d)> %s\n%s", line, c.Interp(line, binary.LittleEndian.Uint64(b[k-8:k]), b), ind)
		line++
	}
}

type Cap struct {
	nextTag  bool
	expected map[int]string
}

func NewCap() *Cap {
	return &Cap{
		expected: make(map[int]string),
	}
}

func (c *Cap) Interp(line int, val uint64, b []byte) string {
	r := ""

	// allowing store of state and re-discovery
	if k, ok := c.expected[line]; ok {
		return k
	}

	if line == 0 {
		numSeg := val&zerohi32 + 1
		words := val >> 32
		return fmt.Sprintf("stream header: %d segment(s), this segment has %d words", numSeg, words)
	} else {
		// assume single segment for now
		switch A(val) {
		case structPointer:
			return c.StructPointer(val, line)
		case listPointer:
			//fmt.Printf("\ndetected List with element count = %d (unless this is a composite). ListB = %d, ListC = %d\n", ListD(val), B(val), ListC(val))

			if ListC(val) == bit1List {
				listSize := ListD(val)
				bytesRequired := (listSize + 7) / 8
				szBytesWordBoundary := (bytesRequired + 7) &^ 7
				eline := line + 1 + B(val)
				listContent := fmt.Sprintf("% 02x", b[eline*8:(eline*8+szBytesWordBoundary)])
				c.expected[eline] = fmt.Sprintf("bit-list contents: %s", listContent)
				return fmt.Sprintf("list of %d bits (pointer to: '%s' at line %d)", listSize, listContent, eline)
			}

			if ListC(val) == byte1List {
				// assume it will be text
				eline := line + 1 + B(val)
				c.expected[eline] = fmt.Sprintf("text contents: %s", string(b[eline*8:(eline*8+ListD(val))]))
				return fmt.Sprintf("list of bytes/Text (pointer to: '%s' at line %d)", string(b[eline*8:(eline*8+ListD(val)-1)]), eline)
			}

			if ListC(val) == compositeList {
				c.nextTag = true
				tagline := line + 1 + B(val)
				tag := binary.LittleEndian.Uint64(b[(tagline)*8 : (tagline+1)*8])
				r = fmt.Sprintf("list-of-composite, count: %d. (from tag at line %d). total-words-not-counting-tag-word: %d", B(tag), line+1+B(val), ListD(val))
				c.expected[tagline] = CompositeTag(tag)
				return r
			}
			eline := line + 1 + B(val)
			return fmt.Sprintf("list, first element starts %d words from here (at line %d). Size: %s, num-elem: %d", B(val), eline, ListCString(val), ListD(val))

		default:
			r += "other"
		}
	}
	return r
}

// lsb                      struct pointer                       msb
// +-+-----------------------------+---------------+---------------+
// |A|             B               |       C       |       D       |
// +-+-----------------------------+---------------+---------------+
//
// A (2 bits) = 0, to indicate that this is a struct pointer.
// B (30 bits) = Offset, in words, from the end of the pointer to the
//     start of the struct's data section.  Signed.
// C (16 bits) = Size of the struct's data section, in words.
// D (16 bits) = Size of the struct's pointer section, in words.
//
// (B is the same for list pointers, but C and D have different size
// and meaning)
//
// B(): extract the count from the B section of a struct pointer
// a.k.a. signedOffsetFromStructPointer()
func B(val uint64) int {
	u64 := uint64(val) & zerohi32
	u32 := uint32(u64)
	s32 := int32(u32) >> 2
	return int(s32)
}

func A(val uint64) int {
	return int(val & 3)
}

func StructC(val uint64) int {
	return int(uint16(val >> 32))
}

func StructD(val uint64) int {
	return int(uint16(val >> 48))
}

func ListC(val uint64) int {
	return int((val >> 32) & 7)
}

func ListCString(val uint64) string {
	switch ListC(val) {
	case voidList:
		return "void"
	case bit1List:
		return "1bit"
	case byte1List:
		return "1byte"
	case byte2List:
		return "2bytes"
	case byte4List:
		return "4bytes"
	case byte8List:
		return "8bytes"
	case pointerList:
		return "pointer"
	case compositeList:
		return "composite"
	default:
		panic("unknown list element size")
	}
}

func ListD(val uint64) int {
	return int(uint32(val >> 35))
}

const (
	structPointer    = 0
	listPointer      = 1
	farPointer       = 2
	doubleFarPointer = 6

	voidList      = 0
	bit1List      = 1
	byte1List     = 2
	byte2List     = 3
	byte4List     = 4
	byte8List     = 5
	pointerList   = 6
	compositeList = 7
)

/*
lsb                       list pointer                        msb
+-+-----------------------------+--+----------------------------+
|A|             B               |C |             D              |
+-+-----------------------------+--+----------------------------+

A (2 bits) = 1, to indicate that this is a list pointer.
B (30 bits) = Offset, in words, from the end of the pointer to the
    start of the first element of the list.  Signed.
C (3 bits) = Size of each element:
    0 = 0 (e.g. List(Void))
    1 = 1 bit
    2 = 1 byte
    3 = 2 bytes
    4 = 4 bytes
    5 = 8 bytes (non-pointer)
    6 = 8 bytes (pointer)
    7 = composite (see below)
D (29 bits) = Number of elements in the list, except when C is 7
    (see below).

The pointed-to values are tightly-packed. In particular, Bools are packed bit-by-bit in little-endian order (the first bit is the least-significant bit of the first byte).

Lists of structs use the smallest element size in which the struct can fit. So, a list of structs that each contain two UInt8 fields and nothing else could be encoded with C = 3 (2-byte elements). A list of structs that each contain a single Text field would be encoded as C = 6 (pointer elements). A list of structs that each contain a single Bool field would be encoded using C = 1 (1-bit elements). A list structs which are each more than one word in size must be be encoded using C = 7 (composite).

When C = 7, the elements of the list are fixed-width composite values – usually, structs. In this case, the list content is prefixed by a "tag" word that describes each individual element. The tag has the same layout as a struct pointer, except that the pointer offset (B) instead indicates the number of elements in the list. Meanwhile, section (D) of the list pointer – which normally would store this element count – instead stores the total number of words in the list (not counting the tag word). The reason we store a word count in the pointer rather than an element count is to ensure that the extents of the list’s location can always be determined by inspecting the pointer alone, without having to look at the tag; this may allow more-efficient prefetching in some use cases. The reason we don’t store struct lists as a list of pointers is because doing so would take significantly more space (an extra pointer per element) and may be less cache-friendly.
*/

func CompositeTag(val uint64) string {
	//return fmt.Sprintf("composite-tag, num elements in list: %d. Each elem: {prim: %d words. pointers: %d words}.", B(val), StructC(val), StructD(val))
	return fmt.Sprintf("composite-tag {prim: %d, pointers: %d words}.", StructC(val), StructD(val))
}

func (c *Cap) StructPointer(val uint64, line int) string {
	if val == 0 {
		return "empty struct, zero valued."
	}
	eline := line + 1 + B(val)
	numprim := StructC(val)
	if numprim > 0 {
		for i := 0; i < numprim; i++ {
			c.expected[eline+i] = fmt.Sprintf("primitive data for struct on line %d", line)
		}
	}
	return fmt.Sprintf("struct-pointer, data starts at +%d words (line %d). {prim: %d, pointers: %d words}.", B(val), eline, StructC(val), StructD(val))
}

func InspectSlice(slice []byte) {
	// Capture the address to the slice structure
	address := unsafe.Pointer(&slice)

	// Create a pointer to the underlying array
	addPtr := (*[8]byte)(unsafe.Pointer(*(*uintptr)(address)))

	fmt.Printf("underlying array Addr[%p]\n", addPtr)
	fmt.Printf("\n\n")
}

func FileExists(name string) bool {
	fi, err := os.Stat(name)
	if err != nil {
		return false
	}
	if fi.IsDir() {
		return false
	}
	return true
}

func ValAtBit(value int64, bitPosition uint) bool {
	return (int64(1)<<bitPosition)&value != 0
}

func TestValAtBit(t *testing.T) {
	const two_to_62 = int64(2) << 61
	tests := []struct {
		value       int64
		bitPosition uint
		bit         bool
	}{
		{0, 0, false},

		{1, 0, true},

		{2, 1, true},
		{2, 0, false},

		{3, 2, false},
		{3, 1, true},
		{3, 0, true},

		{4, 3, false},
		{4, 2, true},
		{4, 1, false},
		{4, 0, false},

		{5, 3, false},
		{5, 2, true},
		{5, 1, false},
		{5, 0, true},

		{6, 3, false},
		{6, 2, true},
		{6, 1, true},
		{6, 0, false},

		{7, 3, false},
		{7, 2, true},
		{7, 1, true},
		{7, 0, true},

		{8, 3, true},
		{8, 2, false},
		{8, 1, false},
		{8, 0, false},

		{two_to_62, 62, true},
		{two_to_62, 2, false},
		{two_to_62, 1, false},
		{two_to_62, 0, false},

		{9, 3, true},
		{9, 2, false},
		{9, 1, false},
		{9, 0, true},
	}
	for _, test := range tests {
		if bit := ValAtBit(test.value, test.bitPosition); bit != test.bit {
			t.Errorf("ValAtBit(%#x, %d) = %t; want %t", test.value, test.bitPosition, bit, test.bit)
		}
	}
}

func zboolvec_value_FilledSegment(value int64, elementCount uint) (*capnp.Segment, []byte) {
	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		panic(err)
	}
	z, err := air.NewRootZ(seg)
	if err != nil {
		panic(err)
	}
	list, err := capnp.NewBitList(seg, int32(elementCount))
	if err != nil {
		panic(err)
	}
	if value > 0 {
		for i := uint(0); i < elementCount; i++ {
			list.Set(int(i), ValAtBit(value, i))
		}
	}
	z.SetBoolvec(list)

	b, err := msg.Marshal()
	if err != nil {
		panic(err)
	}
	return seg, b
}

// encodeTestMessage encodes the textual Cap'n Proto message to unpacked
// binary using the capnp tool, or returns the fallback if the tool fails.
func encodeTestMessage(typ string, text string, fallback []byte) ([]byte, error) {
	// TODO(light): use fallback if tool not present.
	b := CapnpEncode(text, typ)
	if !bytes.Equal(b, fallback) {
		return nil, fmt.Errorf("%s value %q =\n%s; fallback is\n%s\nFallback out of date?", typ, text, hex.Dump(b), hex.Dump(fallback))
	}
	return b, nil
}

// mustEncodeTestMessage encodes the textual Cap'n Proto message to unpacked
// binary using the capnp tool, or returns the fallback if the tool fails.
func mustEncodeTestMessage(t testing.TB, typ string, text string, fallback []byte) []byte {
	b, err := encodeTestMessage(typ, text, fallback)
	if err != nil {
		if _, fname, line, ok := runtime.Caller(1); ok {
			t.Fatalf("%s:%d: %v", filepath.Base(fname), line, err)
		} else {
			t.Fatal(err)
		}
	}
	return b
}
