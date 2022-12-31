package capnp

// Helpers for formatting strings. We avoid the use of the fmt package for the benefit
// of environments that care about minimizing executable size.

import "strconv"

// An address is an index inside a segment's data (in bytes).
// It is bounded to [0, maxSegmentSize).
type address uint32

// fmtUdecimal is a helper for formatting unsigned integers as decimals.
func fmtUdecimal[T isUint](n T) string {
	return strconv.FormatUint(uint64(n), 10)
}

// Like fmtIdecimal, but for signed integers.
func fmtIdecimal[T isInt](n T) string {
	return strconv.FormatInt(int64(n), 10)
}

type isUint interface {
	~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uint
}

type isInt interface {
	~int8 | ~int16 | ~int32 | ~int64 | ~int
}
