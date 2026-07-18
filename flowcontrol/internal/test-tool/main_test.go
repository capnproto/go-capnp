package main

import (
	"testing"
	"time"
)

func TestTransferDelay(t *testing.T) {
	tests := []struct {
		name           string
		size           int
		bytesPerSecond uint64
		want           time.Duration
	}{
		{
			name:           "one MiB per second",
			size:           8192,
			bytesPerSecond: 1_000_000,
			want:           8_192 * time.Microsecond,
		},
		{
			name:           "one GiB per second",
			size:           8192,
			bytesPerSecond: 1_000_000_000,
			want:           8_192 * time.Nanosecond,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := transferDelay(test.size, test.bytesPerSecond); got != test.want {
				t.Errorf("transferDelay(%d, %d) = %v; want %v", test.size, test.bytesPerSecond, got, test.want)
			}
		})
	}
}
