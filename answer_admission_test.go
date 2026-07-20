package capnp

import (
	"context"
	"errors"
	"sync"
	"testing"
)

type blockingPipelineCaller struct {
	admitted      chan struct{}
	release       <-chan struct{}
	panicOnReturn bool
}

func (c *blockingPipelineCaller) wait() {
	c.admitted <- struct{}{}
	<-c.release
	if c.panicOnReturn {
		panic("pipeline caller panic")
	}
}

func (c *blockingPipelineCaller) PipelineSend(context.Context, []PipelineOp, Send) (*Answer, ReleaseFunc) {
	c.wait()
	return ErrorAnswer(Method{}, errors.New("test call")), func() {}
}

func (c *blockingPipelineCaller) PipelineRecv(context.Context, []PipelineOp, Recv) PipelineCaller {
	c.wait()
	return nil
}

type signalingPtrResolver struct {
	called  chan struct{}
	release <-chan struct{}
}

func (r signalingPtrResolver) signal() {
	close(r.called)
	if r.release != nil {
		<-r.release
	}
}

func (r signalingPtrResolver) Fulfill(Ptr)  { r.signal() }
func (r signalingPtrResolver) Reject(error) { r.signal() }

func TestPromiseResolutionDrainsAdmittedCalls(t *testing.T) {
	for _, recv := range []bool{false, true} {
		name := "send"
		if recv {
			name = "recv"
		}
		t.Run(name, func(t *testing.T) {
			const callCount = 64
			admitted := make(chan struct{}, callCount)
			release := make(chan struct{})
			resolverCalled := make(chan struct{})
			resolverRelease := make(chan struct{})
			p := NewPromise(Method{}, &blockingPipelineCaller{
				admitted: admitted,
				release:  release,
			}, signalingPtrResolver{called: resolverCalled, release: resolverRelease})

			// Materialize a promised client so the test also verifies that
			// promised clients resolve before admission drain blocks Resolve.
			promisedClient := p.Answer().Field(0, nil).Client()
			clientResolved := make(chan struct{})
			go func() {
				_ = promisedClient.Resolve(context.Background())
				close(clientResolved)
			}()

			var calls sync.WaitGroup
			calls.Add(callCount)
			for i := 0; i < callCount; i++ {
				go func() {
					defer calls.Done()
					if recv {
						p.Answer().PipelineRecv(context.Background(), nil, Recv{})
						return
					}
					_, release := p.Answer().PipelineSend(context.Background(), nil, Send{Method: Method{}})
					release()
				}()
			}
			for i := 0; i < callCount; i++ {
				<-admitted
			}

			msg, seg := NewSingleSegmentMessage(nil)
			defer msg.Release()
			result, err := NewStruct(seg, ObjectSize{PointerCount: 1})
			if err != nil {
				t.Fatal(err)
			}
			client := ErrorClient(errors.New("resolved client"))
			if err := result.SetPtr(0, NewInterface(seg, msg.CapTable().Add(client)).ToPtr()); err != nil {
				t.Fatal(err)
			}

			resolveReturned := make(chan struct{})
			go func() {
				p.Resolve(result.ToPtr(), nil)
				close(resolveReturned)
			}()

			<-resolverCalled
			select {
			case <-clientResolved:
				t.Fatal("promised client resolved before Resolver returned")
			default:
			}
			close(resolverRelease)
			<-clientResolved
			select {
			case <-p.Answer().Done():
				t.Fatal("result published before admitted calls yielded")
			default:
			}
			select {
			case <-resolveReturned:
				t.Fatal("Resolve returned before admitted calls yielded")
			default:
			}
			lateCtx, cancelLateCall := context.WithCancel(context.Background())
			lateCallReturned := make(chan error, 1)
			go func() {
				ans, release := p.Answer().PipelineSend(lateCtx, nil, Send{Method: Method{}})
				_, err := ans.Struct()
				release()
				lateCallReturned <- err
			}()
			cancelLateCall()
			if err := <-lateCallReturned; !errors.Is(err, context.Canceled) {
				t.Fatalf("call arriving during admission drain returned %v; want context.Canceled", err)
			}
			select {
			case <-admitted:
				t.Fatal("call admitted to old PipelineCaller after resolution began")
			default:
			}

			close(release)
			calls.Wait()
			<-resolveReturned
			<-p.Answer().Done()
			ans, releaseCall := p.Answer().PipelineSend(context.Background(), nil, Send{Method: Method{}})
			_, _ = ans.Struct()
			releaseCall()
			select {
			case <-admitted:
				t.Fatal("call admitted to old PipelineCaller after result publication")
			default:
			}
			p.ReleaseClients()
		})
	}
}

func TestPromiseAdmissionReleasedWhenPipelineCallerPanics(t *testing.T) {
	admitted := make(chan struct{}, 1)
	release := make(chan struct{})
	resolverCalled := make(chan struct{})
	p := NewPromise(Method{}, &blockingPipelineCaller{
		admitted:      admitted,
		release:       release,
		panicOnReturn: true,
	}, signalingPtrResolver{called: resolverCalled})

	panicked := make(chan struct{})
	go func() {
		defer close(panicked)
		defer func() { _ = recover() }()
		p.Answer().PipelineSend(context.Background(), nil, Send{Method: Method{}})
	}()
	<-admitted

	resolveReturned := make(chan struct{})
	go func() {
		p.Resolve(Ptr{}, nil)
		close(resolveReturned)
	}()
	<-resolverCalled
	close(release)
	<-panicked
	<-resolveReturned
	p.ReleaseClients()
}
