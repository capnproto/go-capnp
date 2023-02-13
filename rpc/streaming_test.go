package rpc_test

import (
	"context"
	"fmt"
	"testing"

	"capnproto.org/go/capnp/v3/rpc/internal/testcapnp"
	"github.com/stretchr/testify/assert"
)

// TestStreamingWaitOk verifies that if no errors occur in streaming calls,
// WaitStreaming retnurs nil.
func TestStreamingWaitOk(t *testing.T) {
	ctx := context.Background()
	client := testcapnp.StreamTest_ServerToClient(&maxPushStream{limit: 1})
	defer client.Release()
	assert.NoError(t, client.Push(ctx, nil))
	assert.NoError(t, client.WaitStreaming())
}

// TestStreamingWaitErr verifies that if an error occurs in a streaming call,
// it shows up in a subsequent call to WaitStreaming().
func TestStreamingWaitErr(t *testing.T) {
	ctx := context.Background()
	client := testcapnp.StreamTest_ServerToClient(&maxPushStream{limit: 0})
	defer client.Release()
	assert.NoError(t, client.Push(ctx, nil))
	assert.NotNil(t, client.WaitStreaming())
}

// A maxPushStream is an implementation of StreamTest that
// starts returning errors after a specified number of calls.
type maxPushStream struct {
	count int // How many calls have we seen?
	limit int // How many calls are permitted?
}

func (m *maxPushStream) Push(context.Context, testcapnp.StreamTest_push) error {
	m.count++
	if m.count > m.limit {
		return fmt.Errorf("Exceeded limit of %v calls", m.limit)
	}
	return nil
}
