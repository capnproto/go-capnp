package rpc

import (
	"errors"
	"testing"

	"capnproto.org/go/capnp/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPromisedAnswerTarget(t *testing.T) {
	t.Parallel()

	msg, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	require.NoError(t, err)
	defer msg.Release()

	t.Run("MissingCapability", func(t *testing.T) {
		target, err := promisedAnswerTarget(capnp.NewInterface(seg, 1).ToPtr(), nil)
		assert.Error(t, err)
		assert.False(t, target.IsValid())
	})

	t.Run("InvalidCapability", func(t *testing.T) {
		target, err := promisedAnswerTarget(capnp.NewInterface(seg, 0).ToPtr(), []capnp.ClientSnapshot{{}})
		assert.Error(t, err)
		assert.False(t, target.IsValid())
	})

	t.Run("NonCapability", func(t *testing.T) {
		st, err := capnp.NewStruct(seg, capnp.ObjectSize{})
		require.NoError(t, err)
		target, err := promisedAnswerTarget(st.ToPtr(), nil)
		require.NoError(t, err)
		defer target.Release()
		assert.True(t, target.IsValid())
	})

	t.Run("Capability", func(t *testing.T) {
		client := capnp.ErrorClient(errors.New("test"))
		defer client.Release()
		snapshot := client.Snapshot()
		defer snapshot.Release()
		target, err := promisedAnswerTarget(capnp.NewInterface(seg, 0).ToPtr(), []capnp.ClientSnapshot{snapshot})
		require.NoError(t, err)
		defer target.Release()
		assert.True(t, target.IsValid())
	})
}
