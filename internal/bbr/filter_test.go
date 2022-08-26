package bbr

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBtlBwFilter(t *testing.T) {
	f := newBtlBwFilter()

	assert.Equal(t, bytesPerNs(1e-10), f.Estimate, "Initial bandwidth estimate is 1 byte per 10s.")
	f.AddSample(4)
	assert.Equal(t, bytesPerNs(4), f.Estimate, "After adding one sample, estimate is that sample.")
	f.AddSample(2)
	assert.Equal(t, bytesPerNs(4), f.Estimate, "After adding a smaller sample, estimate is unchanged.")
	f.AddSample(6)
	assert.Equal(t, bytesPerNs(6), f.Estimate, "After adding a larger sample, estimate equals new sample.")

	for i := 0; i < btlBwFilterSize; i++ {
		f.AddSample(1)
	}
	assert.Equal(t, bytesPerNs(1), f.Estimate, "Older samples slide out of the window correctly.")
}

func TestRtPropFilter(t *testing.T) {
	f := newRtPropFilter()

	now := sampleStartTime

	addSample := func(rtt time.Duration) {
		f.AddSample(rtPropSample{
			now: now,
			rtt: rtt,
		})
	}

	addSample(4)
	assert.Equal(t, time.Duration(4), f.Estimate, "After adding one sample, estimate is that sample.")
	assert.Equal(t, 0, f.q.Len(), "Queue is empty after the first sample.")
	now = now.Add(time.Second / 10)
	addSample(6)
	assert.Equal(t, time.Duration(4), f.Estimate, "After adding a larger sample, estimate is unchanged.")
	assert.Equal(t, 0, f.q.Len(), "Queue is empty before 1 second has passed.")

	now = now.Add(1 * time.Second)
	addSample(2)
	assert.Equal(t, time.Duration(2), f.Estimate, "After adding a smaller sample, estimate equals new sample.")
	assert.Equal(t, 1, f.q.Len(), "After one second, old sample was added to the queue.")

	now = now.Add(15 * time.Second)
	addSample(8)
	assert.Equal(t, time.Duration(2), f.Estimate)
	now = now.Add(16 * time.Second)
	addSample(12)
	assert.Equal(t, time.Duration(8), f.Estimate, "After 30 seconds, old entries are removed from the queue.")

}
