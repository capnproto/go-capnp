//go:build go1.9
// +build go1.9

package transport

import (
	"errors"
	"io"
	"os"
	"testing"
	"time"
)

// It is only safe to call Read and Close concurrently on an *os.File in
// Go 1.9 and above.  See https://golang.org/issue/7970

func TestFileStreamTransport(t *testing.T) {
	testTransport(t, func() (t1, t2 Transport, err error) {
		r1, w1, err := os.Pipe()
		if err != nil {
			return nil, nil, err
		}
		r2, w2, err := os.Pipe()
		if err != nil {
			r1.Close()
			w1.Close()
			return nil, nil, err
		}
		t1 = NewStream(readWriteCloser{r1, w2})
		t2 = NewStream(readWriteCloser{r2, w1})
		return t1, t2, nil
	})
}

type readWriteCloser struct {
	r io.ReadCloser
	w io.WriteCloser
}

func (rwc readWriteCloser) Read(p []byte) (int, error) {
	return rwc.r.Read(p)
}

func (rwc readWriteCloser) SetReadDeadline(t time.Time) error {
	d, ok := rwc.r.(interface {
		SetReadDeadline(time.Time) error
	})
	if !ok {
		return errors.New("read deadline not implemented")
	}
	return d.SetReadDeadline(t)
}

func (rwc readWriteCloser) Write(p []byte) (int, error) {
	return rwc.w.Write(p)
}

func (rwc readWriteCloser) SetWriteDeadline(t time.Time) error {
	d, ok := rwc.w.(interface {
		SetWriteDeadline(time.Time) error
	})
	if !ok {
		return errors.New("write deadline not implemented")
	}
	return d.SetWriteDeadline(t)
}

func (rwc readWriteCloser) Close() error {
	werr := rwc.w.Close()
	rerr := rwc.r.Close()
	if werr != nil {
		return werr
	}
	if rerr != nil {
		return rerr
	}
	return nil
}
