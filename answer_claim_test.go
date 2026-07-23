package capnp

import (
	"testing"
)

type testAdmissionPipelineCaller struct {
	dummyPipelineCaller
}

func (*testAdmissionPipelineCaller) pipelineAdmissionController() {}

var _ pipelineAdmissionController = (*testAdmissionPipelineCaller)(nil)

func TestPromisePipelineCallClaimDrainsOnce(t *testing.T) {
	caller := new(testAdmissionPipelineCaller)
	resolverCalled := make(chan struct{})
	p := NewPromise(Method{}, caller, signalingPtrResolver{called: resolverCalled})

	claim, state := p.claimPipelineCall(caller)
	if state != pipelineUnresolved || claim == nil {
		t.Fatalf("claimPipelineCall = (%v, %v); want non-nil, unresolved", claim, state)
	}

	resolved := make(chan struct{})
	go func() {
		p.Resolve(Ptr{}, nil)
		close(resolved)
	}()
	<-resolverCalled
	select {
	case <-resolved:
		t.Fatal("Resolve returned while claim was outstanding")
	default:
	}

	claim.Done()
	claim.Done()
	<-resolved
}

func TestPromisePipelineCallClaimReportsResolutionState(t *testing.T) {
	t.Run("pending", func(t *testing.T) {
		caller := new(testAdmissionPipelineCaller)
		resolverCalled := make(chan struct{})
		p := NewPromise(Method{}, caller, signalingPtrResolver{called: resolverCalled})
		claim, state := p.claimPipelineCall(caller)
		if state != pipelineUnresolved || claim == nil {
			t.Fatalf("initial claim = (%v, %v); want non-nil, unresolved", claim, state)
		}

		resolved := make(chan struct{})
		go func() {
			p.Resolve(Ptr{}, nil)
			close(resolved)
		}()
		<-resolverCalled

		lateClaim, state := p.claimPipelineCall(caller)
		if state != pipelinePendingResolution || lateClaim != nil {
			t.Fatalf("claim during resolution = (%v, %v); want nil, pending", lateClaim, state)
		}
		claim.Done()
		<-resolved
	})

	t.Run("resolved", func(t *testing.T) {
		caller := new(testAdmissionPipelineCaller)
		p := NewPromise(Method{}, caller, nil)
		p.Resolve(Ptr{}, nil)

		claim, state := p.claimPipelineCall(caller)
		if state != pipelineResolved || claim != nil {
			t.Fatalf("claim after resolution = (%v, %v); want nil, resolved", claim, state)
		}
	})
}
