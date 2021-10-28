//go:build go1.5
// +build go1.5

package packed

import (
	"bufio"
)

func discard(r *bufio.Reader, n int) {
	r.Discard(n)
}
