package thunk

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLazy(t *testing.T) {
	called := 0

	th := Lazy(func() int {
		called++
		return 7
	})

	assert.Equal(t, 0, called, "Lazy does not call the function immediately.")
	assert.Equal(t, 7, th.Force(), "Force returns the expected result")
	assert.Equal(t, 1, called, "Force calls the function once")
	assert.Equal(t, 7, th.Force(), "Force returns the same result a second time")
	assert.Equal(t, 1, called, "Force does not call the function twice")
}

func TestGo(t *testing.T) {
	done := make(chan struct{})

	th := Go(func() int {
		close(done)
		return 1
	})
	<-done // Should call the function immediately; otherwise this will hang.

	assert.Equal(t, 1, th.Force(), "Force returns the expected result")
}

func TestPromise(t *testing.T) {
	th, fulfill := Promise[int]()

	resultChan := make(chan int)

	go func() {
		v := th.Force()
		resultChan <- v
	}()

	select {
	case <-time.NewTimer(10 * time.Millisecond).C:
	case <-resultChan:
		t.Fatal("Result should not be ready before fulfill() is called.")
	}
	fulfill(1)
	assert.Equal(t, 1, <-resultChan, "Result should be ready after fulfill.")
	assert.Equal(t, 1, th.Force(), "Force should return the result a second time.")
	assert.Panics(t, func() {
		fulfill(1)
	}, "Calling fulfill a second time panics.")
}
