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
		p := NewPromise(dummyMethod, dummyPipelineCaller{})
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
		p := NewPromise(dummyMethod, dummyPipelineCaller{})
		defer p.ReleaseClients()
		ans := p.Answer()
		p.Reject(errors.New("omg bbq"))
		if _, err := ans.Struct(); err == nil || !strings.Contains(err.Error(), "omg bbq") || !strings.Contains(err.Error(), "Foo.bar") {
			t.Errorf("answer error = %v; want message containing \"omg bbq\" and \"Foo.bar\"", err)
		}
	})
	t.Run("Client", func(t *testing.T) {
		p := NewPromise(dummyMethod, dummyPipelineCaller{})
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

func TestPromiseFulfill(t *testing.T) {
	t.Run("Done", func(t *testing.T) {
		p := NewPromise(dummyMethod, dummyPipelineCaller{})
		done := p.Answer().Done()
		msg, seg, _ := NewMessage(SingleSegment(nil))
		defer msg.Reset(nil)
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
		p := NewPromise(dummyMethod, dummyPipelineCaller{})
		defer p.ReleaseClients()
		ans := p.Answer()
		msg, seg, _ := NewMessage(SingleSegment(nil))
		defer msg.Reset(nil)
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
		p := NewPromise(dummyMethod, dummyPipelineCaller{})
		defer p.ReleaseClients()
		pc := p.Answer().Field(1, nil).Client()

		h := new(dummyHook)
		c := NewClient(h)
		defer c.Release()
		msg, seg, _ := NewMessage(SingleSegment(nil))
		defer msg.Reset(nil)
		res, _ := NewStruct(seg, ObjectSize{PointerCount: 3})
		res.SetPtr(1, NewInterface(seg, msg.AddCap(c.AddRef())).ToPtr())

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

func TestPromiseJoin(t *testing.T) {
	t.Run("BeforeReject", func(t *testing.T) {
		pa := NewPromise(dummyMethod, dummyPipelineCaller{})
		pb := NewPromise(dummyMethod, dummyPipelineCaller{})
		ansB := pb.Answer()
		doneB := ansB.Done()
		pb.Join(pa.Answer())
		pa.Reject(errors.New("omg bbq"))
		select {
		case <-doneB:
			if _, err := ansB.Struct(); err == nil || !strings.Contains(err.Error(), "omg bbq") {
				t.Errorf("joined answer error = %v; want \"omg bbq\"", err)
			}
			pa.ReleaseClients()
			pb.ReleaseClients()
		default:
			t.Error("joined answer not resolved")
		}
	})
	t.Run("AfterReject", func(t *testing.T) {
		pa := NewPromise(dummyMethod, dummyPipelineCaller{})
		pb := NewPromise(dummyMethod, dummyPipelineCaller{})
		ansB := pb.Answer()
		doneB := ansB.Done()
		pa.Reject(errors.New("omg bbq"))
		pb.Join(pa.Answer())
		select {
		case <-doneB:
			if _, err := ansB.Struct(); err == nil || !strings.Contains(err.Error(), "omg bbq") {
				t.Errorf("joined answer error = %v; want \"omg bbq\"", err)
			}
			pa.ReleaseClients()
			pb.ReleaseClients()
		default:
			t.Error("joined answer not resolved")
		}
	})
	t.Run("MultipleJoinReject", func(t *testing.T) {
		pa := NewPromise(dummyMethod, dummyPipelineCaller{})
		pb := NewPromise(dummyMethod, dummyPipelineCaller{})
		pc := NewPromise(dummyMethod, dummyPipelineCaller{})
		pc.Join(pb.Answer())
		pb.Join(pa.Answer())
		pa.Reject(errors.New("omg bbq"))
		select {
		case <-pb.Answer().Done():
			if _, err := pb.Answer().Struct(); err == nil || !strings.Contains(err.Error(), "omg bbq") {
				t.Errorf("directly joined answer error = %v; want \"omg bbq\"", err)
			}
			pb.ReleaseClients()
		default:
			t.Error("directly joined answer not resolved")
		}
		select {
		case <-pc.Answer().Done():
			if _, err := pc.Answer().Struct(); err == nil || !strings.Contains(err.Error(), "omg bbq") {
				t.Errorf("transitively joined answer error = %v; want \"omg bbq\"", err)
			}
			pc.ReleaseClients()
		default:
			t.Error("transitively joined answer not resolved")
		}
		pa.ReleaseClients()
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
