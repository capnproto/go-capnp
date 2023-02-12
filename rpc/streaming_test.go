package rpc_test

import (
	"context"
	"fmt"
	"testing"

	"capnproto.org/go/capnp/v3/rpc/internal/testcapnp"
	"github.com/stretchr/testify/assert"
)

// TestStreamingDoneErr verifies that if an error occurs in a streaming call,
// it shows up in a subsequent call to done().
func TestStreamingDoneErr(t *testing.T) {
	ctx := context.Background()
	client := testcapnp.StreamTest_ServerToClient(&maxPushStream{limit: 0})
	defer client.Release()
	assert.NoError(t, client.Push(ctx, nil))
	fut, rel := client.Done(ctx, nil)
	defer rel()
	_, err := fut.Struct()
	assert.NotNil(t, err)
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

func (m *maxPushStream) Done(context.Context, testcapnp.StreamTest_done) error {
	// Note: important to always return nil here, so the tests can distinguish
	// between an error carried over from a call to push() and one returned
	// directly from done().
	return nil
}
