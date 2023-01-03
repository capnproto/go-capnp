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

type Uint interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint
}

type Int interface {
	~int8 | ~int16 | ~int32 | ~int64 | ~int
}
