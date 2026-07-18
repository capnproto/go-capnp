package capnp

import (
	"context"
	"errors"
	"strings"
	"testing"
)

var dummyMethod = Method{
	InterfaceID:   0xa7317bd7216570aa,
	InterfaceName: "Foo",
	MethodID:      9,
	MethodName:    "bar",
}

func TestPromiseReject(t *testing.T) {
	t.Run("Done", func(t *testing.T) {
		p := NewPromise(dummyMethod, dummyPipelineCaller{}, nil)
		done := p.Answer().Done()
		p.Reject(errors.New("omg bbq"))
		select {
		case <-done:
			p.ReleaseClients()
		default:
			t.Error("answer not resolved")
		}
	})
	t.Run("Struct", func(t *testing.T) {
		p := NewPromise(dummyMethod, dummyPipelineCaller{}, nil)
		defer p.ReleaseClients()
		ans := p.Answer()
		p.Reject(errors.New("omg bbq"))
		if _, err := ans.Struct(); err == nil || !strings.Contains(err.Error(), "omg bbq") || !strings.Contains(err.Error(), "Foo.bar") {
			t.Errorf("answer error = %v; want message containing \"omg bbq\" and \"Foo.bar\"", err)
		}
	})
	t.Run("Client", func(t *testing.T) {
		p := NewPromise(dummyMethod, dummyPipelineCaller{}, nil)
		defer p.ReleaseClients()
		pc := p.Answer().Field(1, nil).Client()
		p.Reject(errors.New("omg bbq"))
		ctx := context.Background()
		if err := pc.Resolve(ctx); err != nil {
			t.Error("pc.Resolve:", err)
		}
		ans, release := pc.SendCall(ctx, Send{})
		_, err := ans.Struct()
		release()
		if err == nil || !strings.Contains(err.Error(), "omg bbq") || !strings.Contains(err.Error(), "Foo.bar") {
			t.Errorf("pc.SendCall error = %v; want message containing \"omg bbq\"", err)
		}
	})
}

func TestAnswerQueueFulfillLeafWithoutPipelineCaller(t *testing.T) {
	aq := NewAnswerQueue(dummyMethod)
	ret1 := new(dummyReturner)
	pcall := aq.PipelineRecv(context.Background(), nil, Recv{
		Method:      dummyMethod,
		Returner:    ret1,
		ReleaseArgs: func() {},
	})
	ret2 := new(dummyReturner)
	pcall.PipelineRecv(context.Background(), nil, Recv{
		Method:      dummyMethod,
		Returner:    ret2,
		ReleaseArgs: func() {},
	})

	msg, seg := NewSingleSegmentMessage(nil)
	defer msg.Release()
	capID := msg.CapTable().Add(ErrorClient(errors.New("test error")))
	aq.Fulfill(NewInterface(seg, capID).ToPtr())

	if !ret1.returned {
		t.Fatal("first queued call was not returned")
	}
	if ret1.err == nil {
		t.Fatal("first queued call was not rejected")
	}
	if !ret2.returned {
		t.Fatal("dependent queued call was not returned")
	}
	if ret2.err == nil {
		t.Fatal("dependent queued call was not rejected")
	}
}

func TestPromiseFulfill(t *testing.T) {
	t.Parallel()

	t.Run("Done", func(t *testing.T) {
		p := NewPromise(dummyMethod, dummyPipelineCaller{}, nil)
		done := p.Answer().Done()
		msg, seg := NewSingleSegmentMessage(nil)
		defer msg.Release()

		res, _ := NewStruct(seg, ObjectSize{DataSize: 8})
		p.Fulfill(res.ToPtr())
		select {
		case <-done:
			p.ReleaseClients()
		default:
			t.Error("answer not resolved")
		}
	})
	t.Run("Struct", func(t *testing.T) {
		p := NewPromise(dummyMethod, dummyPipelineCaller{}, nil)
		defer p.ReleaseClients()
		ans := p.Answer()
		msg, seg := NewSingleSegmentMessage(nil)
		defer msg.Release()

		res, _ := NewStruct(seg, ObjectSize{DataSize: 8})
		res.SetUint32(0, 0xdeadbeef)
		p.Fulfill(res.ToPtr())
		s, err := ans.Struct()
		if eq, err := Equal(res.ToPtr(), s.ToPtr()); err != nil {
			t.Error("Equal(p.Answer.Struct(), res):", err)
		} else if !eq {
			t.Error("p.Answer().Struct() != res")
		}
		if err != nil {
			t.Error("p.Answer().Struct():", err)
		}
	})
	t.Run("Client", func(t *testing.T) {
		p := NewPromise(dummyMethod, dummyPipelineCaller{}, nil)
		defer p.ReleaseClients()
		pc := p.Answer().Field(1, nil).Client()

		h := new(dummyHook)
		c := NewClient(h)
		defer c.Release()
		msg, seg := NewSingleSegmentMessage(nil)
		defer msg.Release()

		res, _ := NewStruct(seg, ObjectSize{PointerCount: 3})
		res.SetPtr(1, NewInterface(seg, msg.CapTable().Add(c.AddRef())).ToPtr())

		p.Fulfill(res.ToPtr())

		ctx := context.Background()
		if err := pc.Resolve(ctx); err != nil {
			t.Error("pc.Resolve:", err)
		}
		if !pc.IsSame(c) {
			t.Errorf("pc != c; pc = %v, c = %v", pc, c)
		}
		c.Release()
		ans, release := pc.SendCall(ctx, Send{})
		_, err := ans.Struct()
		release()
		if err != nil {
			t.Error("pc.SendCall:", err)
		}
		if h.calls == 0 {
			t.Error("hook never called")
		}
	})
}

type dummyPipelineCaller struct{}

func (dummyPipelineCaller) PipelineRecv(ctx context.Context, transform []PipelineOp, r Recv) PipelineCaller {
	r.Reject(errors.New("dummy call"))
	return nil
}

func (dummyPipelineCaller) PipelineSend(ctx context.Context, transform []PipelineOp, s Send) (*Answer, ReleaseFunc) {
	return ErrorAnswer(s.Method, errors.New("dummy call")), func() {}
}
