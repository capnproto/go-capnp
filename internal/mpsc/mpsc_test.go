package mpsc

import (
	"context"
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTryRecvEmpty(t *testing.T) {
	t.Parallel()
	q := New()
	v, ok := q.TryRecv()
	assert.False(t, ok, "TryRecv() on an empty queue succeeded; recevied: ", v)
}

func TestRecvEmpty(t *testing.T) {
	t.Parallel()
	q := New()

	// Recv() on an empty queue should block until the context is canceled.
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()

	v, err := q.Recv(ctx)
	assert.Equal(t, ctx.Err(), err, "Returned error is not ctx.Err()")
	assert.NotNil(t, err, "Returned error is nil.")
	assert.Nil(t, v, "Return value is not nil.")
}

func TestSendOne(t *testing.T) {
	t.Parallel()
	q := New()
	want := Value(1)
	q.Send(want)
	got, err := q.Recv(context.Background())
	assert.Nil(t, err, "Non-nil error from Recv()")
	assert.Equal(t, want, got, "Sent and received different values.")
}

func TestSendManySequential(t *testing.T) {
	t.Parallel()

	inputs := []int{}
	outputs := []int{}
	for i := 0; i < 100; i++ {
		inputs = append(inputs, i)
	}

	ctx := context.Background()

	q := New()

	for _, v := range inputs {
		q.Send(v)
	}

	for len(outputs) != len(inputs) {
		v, err := q.Recv(ctx)
		assert.Nil(t, err, "Non-nil error from Recv()")
		outputs = append(outputs, v.(int))
	}
	assert.Equal(t, inputs, outputs, "Received sequence was different from sent.")

	v, ok := q.TryRecv()
	assert.False(t, ok, "Recieved more values than expected: ", v)
}

func TestSendManyConcurrent(t *testing.T) {
	t.Parallel()

	q := New()

	for i := 0; i < 100; i += 10 {
		for j := 0; j < 10; j++ {
			value := i + j
			delay := time.Duration(rand.Float64() * 100 * float64(time.Millisecond))
			go func() {
				// Sleep a random amount of time before proceeding, rather
				// than relying on the scheduler to give us good test coverage.
				//
				// But, also, we spawn more than one goroutine with the same
				// delay, so that we'll get some goroutines scheduled close to
				// one another.
				time.Sleep(delay)
				q.Send(value)
			}()
		}
	}

	expected := []int{}
	actual := []int{}
	for i := 0; i < 100; i++ {
		expected = append(expected, i)
	}

	ctx := context.Background()
	for i := 0; i < 100; i++ {
		v, err := q.Recv(ctx)
		assert.Nil(t, err, "Failed to receive from queue: ", err)
		actual = append(actual, v.(int))
	}
	// Values come out in random order, so sort them:
	sort.Slice(actual, func(i, j int) bool {
		return actual[i] < actual[j]
	})

	assert.Equal(t, expected, actual, "Different values came out of the queue than went in.")

	v, ok := q.TryRecv()
	assert.False(t, ok, "Recieved more values than expected: ", v)
}
