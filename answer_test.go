package capnp

import (
	"context"
	"errors"
	"strings"
	"testing"
)

func TestPromiseJoin(t *testing.T) {
	t.Run("BeforeReject", func(t *testing.T) {
		pa := NewPromise(dummyPipelineCaller{})
		pb := NewPromise(dummyPipelineCaller{})
		ansB := pb.Answer()
		doneB := ansB.Done()
		pb.Join(pa.Answer())
		pa.Reject(errors.New("omg bbq"))
		select {
		case <-doneB:
		default:
			t.Fatal("joined answer not resolved")
		}
		if _, err := ansB.Struct(); err == nil || !strings.Contains(err.Error(), "omg bbq") {
			t.Errorf("joined answer error = %v; want \"omg bbq\"", err)
		}
	})
	t.Run("AfterReject", func(t *testing.T) {
		pa := NewPromise(dummyPipelineCaller{})
		pb := NewPromise(dummyPipelineCaller{})
		ansB := pb.Answer()
		doneB := ansB.Done()
		pa.Reject(errors.New("omg bbq"))
		pb.Join(pa.Answer())
		select {
		case <-doneB:
		default:
			t.Fatal("joined answer not resolved")
		}
		if _, err := ansB.Struct(); err == nil || !strings.Contains(err.Error(), "omg bbq") {
			t.Errorf("joined answer error = %v; want \"omg bbq\"", err)
		}
	})
	t.Run("MultipleJoinReject", func(t *testing.T) {
		pa := NewPromise(dummyPipelineCaller{})
		pb := NewPromise(dummyPipelineCaller{})
		pc := NewPromise(dummyPipelineCaller{})
		pc.Join(pb.Answer())
		pb.Join(pa.Answer())
		pa.Reject(errors.New("omg bbq"))
		select {
		case <-pb.Answer().Done():
			if _, err := pb.Answer().Struct(); err == nil || !strings.Contains(err.Error(), "omg bbq") {
				t.Errorf("directly joined answer error = %v; want \"omg bbq\"", err)
			}
		default:
			t.Error("directly joined answer not resolved")
		}
		select {
		case <-pc.Answer().Done():
			if _, err := pc.Answer().Struct(); err == nil || !strings.Contains(err.Error(), "omg bbq") {
				t.Errorf("transitively joined answer error = %v; want \"omg bbq\"", err)
			}
		default:
			t.Error("transitively joined answer not resolved")
		}
	})
}

type dummyPipelineCaller struct{}

func (dummyPipelineCaller) PipelineRecv(ctx context.Context, transform []PipelineOp, r Recv) (*Answer, ReleaseFunc) {
	r.ReleaseArgs()
	return ErrorAnswer(errors.New("dummy call")), func() {}
}

func (dummyPipelineCaller) PipelineSend(ctx context.Context, transform []PipelineOp, s Send) (*Answer, ReleaseFunc) {
	return ErrorAnswer(errors.New("dummy call")), func() {}
}
