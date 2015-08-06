package capnp

import (
	"errors"
	"testing"
)

func TestFulfiller_NewShouldBeUnresolved(t *testing.T) {
	f := new(Fulfiller)

	if a := f.Peek(); a != nil {
		t.Errorf("f.Peek() = %v; want nil", a)
	}
	select {
	case <-f.Done():
		t.Error("Done closed early")
	default:
		// success
	}
}

func TestFulfiller_FulfillShouldResolve(t *testing.T) {
	f := new(Fulfiller)
	s := NewBuffer(nil)
	st := s.NewRootStruct(ObjectSize{0, 0})

	f.Fulfill(st)

	select {
	case <-f.Done():
	default:
		t.Error("Done still closed after Fulfill")
	}
	ret, err := f.Struct()
	if err != nil {
		t.Errorf("f.Struct() error: %v", err)
	}
	if ret != st {
		t.Errorf("f.Struct() = %v; want %v", ret, st)
	}
}

func TestFulfiller_RejectShouldResolve(t *testing.T) {
	f := new(Fulfiller)
	e := errors.New("failure and rejection")

	f.Reject(e)

	select {
	case <-f.Done():
	default:
		t.Error("Done still closed after Reject")
	}
	ret, err := f.Struct()
	if err != e {
		t.Errorf("f.Struct() error = %v; want %v", err, e)
	}
	if Object(ret).Type() != TypeNull {
		t.Errorf("f.Struct() = %v; want null", ret)
	}
}

func TestFulfiller_QueuedCallsDeliveredInOrder(t *testing.T) {
	f := new(Fulfiller)
	oc := new(orderClient)
	result := NewBuffer(nil).NewRootStruct(ObjectSize{PointerCount: 1})
	in := result.Segment.Message.AddCap(oc)
	result.SetObject(0, Object(result.Segment.NewInterface(in)))

	ans1 := f.PipelineCall([]PipelineOp{{Field: 0}}, new(Call))
	ans2 := f.PipelineCall([]PipelineOp{{Field: 0}}, new(Call))
	f.Fulfill(result)
	ans3 := f.PipelineCall([]PipelineOp{{Field: 0}}, new(Call))
	ans3.Struct()
	ans4 := f.PipelineCall([]PipelineOp{{Field: 0}}, new(Call))

	check := func(a Answer, n uint64) {
		r, err := a.Struct()
		if r.Get64(0) != n {
			t.Errorf("r%d = %d; want %d", n+1, r.Get64(0), n)
		}
		if err != nil {
			t.Errorf("err%d = %v", n+1, err)
		}
	}
	check(ans1, 0)
	check(ans2, 1)
	check(ans3, 2)
	check(ans4, 3)
}

type orderClient int

func (oc *orderClient) Call(cl *Call) Answer {
	s := NewBuffer(nil)
	st := s.NewRootStruct(ObjectSize{DataSize: 8})
	st.Set64(0, uint64(*oc))
	*oc++
	return ImmediateAnswer(st)
}

func (oc *orderClient) Close() error {
	return nil
}
