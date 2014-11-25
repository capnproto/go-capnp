package capn_test

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"unsafe"
)

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

	schfn := "aircraftlib/aircraft.capnp"
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

	schfn := "aircraftlib/aircraft.capnp"
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

func MakeAndMoveToTempDir() (origdir string, tmpdir string) {

	var err error
	origdir, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	tmpdir, err = ioutil.TempDir(origdir, "tempgocapnpdir")
	if err != nil {
		panic(err)
	}
	err = os.Chdir(tmpdir)
	if err != nil {
		panic(err)
	}

	return origdir, tmpdir
}

func TempDirCleanup(origdir string, tmpdir string) {
	// cleanup
	os.Chdir(origdir)
	err := os.RemoveAll(tmpdir)
	if err != nil {
		panic(err)
	}
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
				listContent := BytesToWordString(b[eline*8 : (eline*8 + szBytesWordBoundary)])
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
	return ""
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

func save(b []byte, fn string) {
	file, err := os.Create(fn)
	if err != nil {
		panic(err)
	}
	file.Write(b)
	file.Close()

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

func DirExists(name string) bool {
	fi, err := os.Stat(name)
	if err != nil {
		return false
	}
	if fi.IsDir() {
		return true
	}
	return false
}

func BytesToWordString(b []byte) string {
	var s string
	k := 0
	for i := 0; i < len(b)/8; i++ {
		for j := 0; j < 8; j++ {
			s += fmt.Sprintf("%02x ", b[k])
			k++
			if k == len(b) {
				break
			}
		}
	}
	return s
}
