package fulfiller

import (
	"errors"
	"testing"

	"zombiezen.com/go/capnproto"
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
	s := capnp.NewBuffer(nil)
	st := s.NewRootStruct(capnp.ObjectSize{})

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
	if capnp.Pointer(ret).Type() != capnp.TypeNull {
		t.Errorf("f.Struct() = %v; want null", ret)
	}
}

func TestFulfiller_QueuedCallsDeliveredInOrder(t *testing.T) {
	f := new(Fulfiller)
	oc := new(orderClient)
	result := capnp.NewBuffer(nil).NewRootStruct(capnp.ObjectSize{PointerCount: 1})
	in := result.Segment().Message.AddCap(oc)
	result.SetPointer(0, capnp.Pointer(result.Segment().NewInterface(in)))

	ans1 := f.PipelineCall([]capnp.PipelineOp{{Field: 0}}, new(capnp.Call))
	ans2 := f.PipelineCall([]capnp.PipelineOp{{Field: 0}}, new(capnp.Call))
	f.Fulfill(result)
	ans3 := f.PipelineCall([]capnp.PipelineOp{{Field: 0}}, new(capnp.Call))
	ans3.Struct()
	ans4 := f.PipelineCall([]capnp.PipelineOp{{Field: 0}}, new(capnp.Call))

	check := func(a capnp.Answer, n uint64) {
		r, err := a.Struct()
		if r.Uint64(0) != n {
			t.Errorf("r%d = %d; want %d", n+1, r.Uint64(0), n)
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

func (oc *orderClient) Call(cl *capnp.Call) capnp.Answer {
	s := capnp.NewBuffer(nil)
	st := s.NewRootStruct(capnp.ObjectSize{DataSize: 8})
	st.SetUint64(0, uint64(*oc))
	*oc++
	return capnp.ImmediateAnswer(st)
}

func (oc *orderClient) Close() error {
	return nil
}
