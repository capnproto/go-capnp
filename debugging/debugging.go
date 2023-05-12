// Package debugging exports various helpers useful when debugging. In general,
// code from this package should not be included in production builds.
package debugging

import "runtime"

// Wrapper around runtime.Stack, which allocates a buffer and returns the stack
// trace as a string. Stack traces over 1MB will be truncated.
func StackString(all bool) string {
	var buf [1e6]byte
	n := runtime.Stack(buf[:], all)
	return string(buf[:n])
}
