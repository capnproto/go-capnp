package queue

import (
	"testing"
)

func TestNew(t *testing.T) {
	qi := make(ints, 5)

	q := New(qi, 0)

	if n := q.Len(); n != 0 {
		t.Errorf("New(qi, 0).Len() = %d; want 0", n)
	}
}

func TestPrepush(t *testing.T) {
	qi := make(ints, 5)
	qi[0] = 42

	q := New(qi, 1)

	if n := q.Len(); n != 1 {
		t.Fatalf("New(qi, 1).Len() = %d; want 1", n)
	}
	if x := q.Pop().(int); x != 42 {
		t.Errorf("Pop() = %d; want 42", x)
	}
}

func TestPush(t *testing.T) {
	qi := make(ints, 5)
	q := New(qi, 0)

	ok := q.Push(42)

	if !ok {
		t.Error("q.Push(42) returned false")
	}
	if n := q.Len(); n != 1 {
		t.Errorf("q.Len() after push = %d; want 1", n)
	}
}

func TestPushFull(t *testing.T) {
	qi := make(ints, 5)
	q := New(qi, 0)
	var ok [6]bool

	ok[0] = q.Push(10)
	ok[1] = q.Push(11)
	ok[2] = q.Push(12)
	ok[3] = q.Push(13)
	ok[4] = q.Push(14)
	ok[5] = q.Push(15)

	for i := 0; i < 5; i++ {
		if !ok[i] {
			t.Errorf("q.Push(%d) returned false", 10+i)
		}
	}
	if ok[5] {
		t.Error("q.Push(15) returned true")
	}
	if n := q.Len(); n != 5 {
		t.Errorf("q.Len() after full = %d; want 5", n)
	}
}

func TestPop(t *testing.T) {
	qi := make(ints, 5)
	q := New(qi, 0)
	q.Push(1)
	q.Push(2)
	q.Push(3)

	outs := make([]int, 0, len(qi))
	outs = append(outs, q.Pop().(int))
	outs = append(outs, q.Pop().(int))
	outs = append(outs, q.Pop().(int))

	if n := q.Len(); n != 0 {
		t.Errorf("q.Len() after pops = %d; want 0", n)
	}
	if outs[0] != 1 {
		t.Errorf("first pop = %d; want 1", outs[0])
	}
	if outs[1] != 2 {
		t.Errorf("first pop = %d; want 2", outs[1])
	}
	if outs[2] != 3 {
		t.Errorf("first pop = %d; want 3", outs[2])
	}
	for i := range qi {
		if qi[i] != 0 {
			t.Errorf("qi[%d] = %d; want 0 (not cleared)", i, qi[i])
		}
	}
}

func TestWrap(t *testing.T) {
	qi := make(ints, 5)
	q := New(qi, 0)
	var ok [7]bool

	ok[0] = q.Push(10)
	ok[1] = q.Push(11)
	ok[2] = q.Push(12)
	q.Pop()
	q.Pop()
	ok[3] = q.Push(13)
	ok[4] = q.Push(14)
	ok[5] = q.Push(15)
	ok[6] = q.Push(16)

	for i := 0; i < 6; i++ {
		if !ok[i] {
			t.Errorf("q.Push(%d) returned false", 10+i)
		}
	}
	if n := q.Len(); n != 5 {
		t.Errorf("q.Len() = %d; want 5", n)
	}
	for i := 12; q.Len() > 0; i++ {
		if x := q.Pop().(int); x != i {
			t.Errorf("q.Pop() = %d; want %d", x, i)
		}
	}
}

type ints []int

func (is ints) Len() int {
	return len(is)
}

func (is ints) At(i int) interface{} {
	return is[i]
}

func (is ints) Set(i int, x interface{}) {
	if x == nil {
		is[i] = 0
	} else {
		is[i] = x.(int)
	}
}
