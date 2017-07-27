package chanmu

import (
	"context"
	"sync"
	"testing"
)

func TestMutex(t *testing.T) {
	var wg sync.WaitGroup
	mu := New()
	x := 0
	add := func() {
		mu.Lock()
		x++
		mu.Unlock()
		wg.Done()
	}
	wg.Add(2)
	go add()
	go add()
	wg.Wait()
	if x != 2 {
		t.Errorf("x = %d; want 2", x)
	}
}

func TestTryLock(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	mu := New()
	if err := mu.TryLock(ctx); err != nil {
		t.Fatal("TryLock:", err)
	}
	cancel()
	if err := mu.TryLock(ctx); err == nil {
		t.Fatal("TryLock on canceled did not fail")
	}
	mu.Unlock()
}
