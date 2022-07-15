package transport

import (
	"testing"
)

func TestPipeTransport(t *testing.T) {
	testTransport(t, func() (t1, t2 Transport, err error) {
		p1, p2 := NewPipe(1)
		return p1, p2, nil
	})
}
