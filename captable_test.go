package capnp_test

import (
	"errors"
	"testing"

	"capnproto.org/go/capnp/v3"
	"github.com/stretchr/testify/assert"
)

func TestCapTable(t *testing.T) {
	t.Parallel()

	var ct capnp.CapTable

	assert.Zero(t, ct.Len(),
		"zero-value CapTable should be empty")
	assert.Zero(t, ct.Add(capnp.Client{}),
		"first entry should have CapabilityID(0)")
	assert.Equal(t, 1, ct.Len(),
		"should increase length after adding capability")

	ct.Reset()
	assert.Zero(t, ct.Len(),
		"zero-value CapTable should be empty after Reset()")
	ct.Reset(capnp.Client{}, capnp.Client{})
	assert.Equal(t, 2, ct.Len(),
		"zero-value CapTable should be empty after Reset(c, c)")

	errTest := errors.New("test")
	ct.Set(capnp.CapabilityID(0), capnp.ErrorClient(errTest))
	snapshot := ct.At(0).Snapshot()
	defer snapshot.Release()
	err := snapshot.Brand().Value.(error)
	assert.ErrorIs(t, errTest, err, "should update client at index 0")
}
