// Helpers for formatting strings. We avoid the use of the fmt package for the benefit
// of environments that care about minimizing executable size.
package str

import "strconv"

// Utod formats unsigned integers as decimals.
func Utod[T Uint](n T) string {
	return strconv.FormatUint(uint64(n), 10)
}

// Itod formats signed integers as decimals.
func Itod[T Int](n T) string {
	return strconv.FormatInt(int64(n), 10)
}

// UToHex returns n formatted in hexidecimal.
func UToHex[T Uint](n T) string {
	return strconv.FormatUint(uint64(n), 16)
}

// ZeroPad pads value to the left with zeros, making the resulting string
// count bytes long.
func ZeroPad(count int, value string) string {
	pad := count - len(value)
	if pad < 0 {
		panic("ZeroPad: count is less than len(value)")
	}
	buf := make([]byte, count)
	for i := 0; i < pad; i++ {
		buf[i] = '0'
	}
	copy(buf[:pad], value[:])
	return string(buf)
}

type Uint interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint
}

type Int interface {
	~int8 | ~int16 | ~int32 | ~int64 | ~int
}
