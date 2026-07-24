package capnp

import (
	"context"
	"errors"
	"sync"
	"testing"

	"capnproto.org/go/capnp/v3/flowcontrol"
)

type gateCompletion struct {
	kind flowcontrol.MessageOutcomeKind
	err  error
}

type testGateController struct {
	waitCalled  chan struct{}
	completions chan gateCompletion
	wait        func(context.Context) error
}

func (g *testGateController) CommitMessage(uint64) (func(context.Context) error, func(flowcontrol.MessageOutcomeKind, error)) {
	return func(ctx context.Context) error {
			g.waitCalled <- struct{}{}
			if g.wait != nil {
				return g.wait(ctx)
			}
			return nil
		}, func(kind flowcontrol.MessageOutcomeKind, err error) {
			g.completions <- gateCompletion{kind: kind, err: err}
		}
}

func (*testGateController) Poison(error) {}

type testGateLimiter struct {
	*ticketTestLimiter
	gate *testGateController
}

func (l *testGateLimiter) GateNext() flowcontrol.GateNextController { return l.gate }

var _ flowcontrol.GateNextFlowLimiter = (*testGateLimiter)(nil)

type testPreparedSend struct {
	size   uint64
	commit func(func(flowcontrol.MessageOutcomeKind, error)) (*Answer, ReleaseFunc, error)

	mu     sync.Mutex
	aborts int
}

func (p *testPreparedSend) Size() uint64 { return p.size }

func (p *testPreparedSend) Commit(terminal func(flowcontrol.MessageOutcomeKind, error)) (*Answer, ReleaseFunc, error) {
	return p.commit(terminal)
}

func (p *testPreparedSend) Abort() {
	p.mu.Lock()
	p.aborts++
	p.mu.Unlock()
}

func (p *testPreparedSend) abortCount() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.aborts
}

var _ PreparedSend = (*testPreparedSend)(nil)

type testPipelinePreparer struct {
	dummyPipelineCaller
	prepared PreparedSend
	called   chan struct{}
}

func (p *testPipelinePreparer) PreparePipelineSend(_ context.Context, _ []PipelineOp, s Send) (PreparedSend, error) {
	close(p.called)
	if s.PlaceArgs != nil {
		if err := s.PlaceArgs(Struct{}); err != nil {
			return nil, err
		}
	}
	return p.prepared, nil
}

var _ PipelineCallerPreparer = (*testPipelinePreparer)(nil)

func newTestGateLimiter() *testGateLimiter {
	return &testGateLimiter{
		ticketTestLimiter: new(ticketTestLimiter),
		gate: &testGateController{
			waitCalled:  make(chan struct{}, 1),
			completions: make(chan gateCompletion, 2),
		},
	}
}

func TestAdmitPreparedSendPublishesAfterCommit(t *testing.T) {
	lim := newTestGateLimiter()
	flow := newClientFlow(lim)
	first := flow.ticketCurrent()
	second := flow.ticketCurrent()
	commitStarted := make(chan struct{})
	finishCommit := make(chan struct{})
	var terminal func(flowcontrol.MessageOutcomeKind, error)
	prepared := &testPreparedSend{
		size: 1,
		commit: func(t func(flowcontrol.MessageOutcomeKind, error)) (*Answer, ReleaseFunc, error) {
			terminal = t
			close(commitStarted)
			<-finishCommit
			return nil, nil, nil
		},
	}

	admitted := make(chan error, 1)
	go func() {
		_, _, err := admitPreparedSend(context.Background(), first, prepared)
		admitted <- err
	}()
	<-commitStarted
	select {
	case <-first.ready:
		t.Fatal("published successor before Commit returned")
	default:
	}
	close(finishCommit)
	if err := <-admitted; err != nil {
		t.Fatalf("admitPreparedSend() = %v; want nil", err)
	}
	select {
	case <-first.ready:
	default:
		t.Fatal("did not publish successor after Commit returned")
	}
	if err := second.await(context.Background()); err != nil {
		t.Fatalf("successor await() = %v; want nil", err)
	}
	select {
	case <-lim.gate.waitCalled:
	default:
		t.Fatal("successor did not use GateNext permission")
	}

	terminal(flowcontrol.MessageOutcomeSucceeded, nil)
	terminal(flowcontrol.MessageOutcomeFatal, errors.New("duplicate terminal"))
	completion := <-lim.gate.completions
	if completion.kind != flowcontrol.MessageOutcomeSucceeded || completion.err != nil {
		t.Fatalf("completion = %+v; want success", completion)
	}
	select {
	case completion := <-lim.gate.completions:
		t.Fatalf("duplicate completion = %+v", completion)
	default:
	}
	select {
	case <-lim.gate.waitCalled:
		t.Fatal("successor invoked GateNext permission more than once")
	default:
	}
	second.finish()
	if release := flow.close(); release != lim {
		t.Fatalf("close release = %v; want limiter", release)
	}
	lim.Release()
}

func TestAdmittedPipelineSendWaitsBeforePreparation(t *testing.T) {
	grant := make(chan struct{})
	lim := newTestGateLimiter()
	lim.gate.wait = func(ctx context.Context) error {
		select {
		case <-grant:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	flow := newClientFlow(lim)
	parentTicket := flow.ticketCurrent()

	childDone := make(chan func(flowcontrol.MessageOutcomeKind, error), 1)
	childPrepared := &testPreparedSend{
		size: 1,
		commit: func(terminal func(flowcontrol.MessageOutcomeKind, error)) (*Answer, ReleaseFunc, error) {
			childDone <- terminal
			return ImmediateAnswer(Method{}, Ptr{}), func() {}, nil
		},
	}
	preparer := &testPipelinePreparer{
		prepared: childPrepared,
		called:   make(chan struct{}),
	}
	parentPromise := NewPromise(Method{}, preparer, nil)
	parentDone := make(chan func(flowcontrol.MessageOutcomeKind, error), 1)
	parentPrepared := &testPreparedSend{
		size: 1,
		commit: func(terminal func(flowcontrol.MessageOutcomeKind, error)) (*Answer, ReleaseFunc, error) {
			parentDone <- terminal
			return parentPromise.Answer(), func() {}, nil
		},
	}
	if _, _, err := admitPreparedSend(context.Background(), parentTicket, parentPrepared); err != nil {
		t.Fatalf("admit parent = %v", err)
	}

	l := parentPromise.state.Lock()
	current := l.Value().caller
	l.Unlock()
	if _, ok := current.(*admittedPipelineCaller); !ok {
		t.Fatalf("parent caller = %T; want *admittedPipelineCaller", current)
	}
	if _, ok := current.(PipelineCallerPreparer); ok {
		t.Fatal("admitted pipeline caller implements PipelineCallerPreparer")
	}

	placeArgs := make(chan struct{})
	callDone := make(chan *Answer, 1)
	go func() {
		ans, _ := parentPromise.Answer().PipelineSend(context.Background(), nil, Send{
			PlaceArgs: func(Struct) error {
				close(placeArgs)
				return nil
			},
		})
		callDone <- ans
	}()
	<-lim.gate.waitCalled
	select {
	case <-preparer.called:
		t.Fatal("prepared pipeline call behind a closed gate")
	default:
	}
	select {
	case <-placeArgs:
		t.Fatal("called PlaceArgs behind a closed gate")
	default:
	}

	close(grant)
	<-preparer.called
	<-placeArgs
	if ans := <-callDone; ans == nil {
		t.Fatal("PipelineSend returned a nil Answer")
	}

	(<-childDone)(flowcontrol.MessageOutcomeSucceeded, nil)
	(<-parentDone)(flowcontrol.MessageOutcomeSucceeded, nil)
	parentPromise.Resolve(Ptr{}, nil)
	if release := flow.close(); release != lim {
		t.Fatalf("close release = %v; want limiter", release)
	}
	lim.Release()
}

func TestAdmittedPipelineWaitDoesNotBlockResolution(t *testing.T) {
	grant := make(chan struct{})
	lim := newTestGateLimiter()
	lim.gate.wait = func(ctx context.Context) error {
		select {
		case <-grant:
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	flow := newClientFlow(lim)
	parentTicket := flow.ticketCurrent()
	waitNext, complete := lim.gate.CommitMessage(1)
	parentTicket.publish(waitNext)

	preparer := &testPipelinePreparer{
		prepared: &testPreparedSend{
			size: 1,
			commit: func(func(flowcontrol.MessageOutcomeKind, error)) (*Answer, ReleaseFunc, error) {
				t.Fatal("prepared unresolved route after Promise rejection")
				return nil, nil, nil
			},
		},
		called: make(chan struct{}),
	}
	p := NewPromise(Method{}, preparer, nil)
	installPipelineAdmission(p.Answer(), flow)

	callDone := make(chan *Answer, 1)
	go func() {
		ans, _ := p.Answer().PipelineSend(context.Background(), nil, Send{})
		callDone <- ans
	}()
	<-lim.gate.waitCalled

	resolvedErr := errors.New("resolved while gated")
	p.Reject(resolvedErr)
	close(grant)
	ans := <-callDone
	if _, err := ans.Struct(); !errors.Is(err, resolvedErr) {
		t.Fatalf("pipelined call after resolution = %v; want %v", err, resolvedErr)
	}
	select {
	case <-preparer.called:
		t.Fatal("used unresolved preparer after Promise resolution won")
	default:
	}

	complete(flowcontrol.MessageOutcomeSucceeded, nil)
	parentTicket.finish()
	if release := flow.close(); release != lim {
		t.Fatalf("close release = %v; want limiter", release)
	}
	lim.Release()
}

func TestAdmitPreparedSendFailureRetiresBeforePublishing(t *testing.T) {
	lim := newTestGateLimiter()
	flow := newClientFlow(lim)
	first := flow.ticketCurrent()
	second := flow.ticketCurrent()
	want := errors.New("commit failed before enqueue")
	prepared := &testPreparedSend{
		size: 1,
		commit: func(func(flowcontrol.MessageOutcomeKind, error)) (*Answer, ReleaseFunc, error) {
			return nil, nil, want
		},
	}

	_, _, err := admitPreparedSend(context.Background(), first, prepared)
	if !errors.Is(err, want) {
		t.Fatalf("admitPreparedSend() = %v; want %v", err, want)
	}
	if got := prepared.abortCount(); got != 1 {
		t.Fatalf("Abort called %d times; want 1", got)
	}
	completion := <-lim.gate.completions
	if completion.kind != flowcontrol.MessageOutcomeAbortedBeforeEnqueue || !errors.Is(completion.err, want) {
		t.Fatalf("completion = %+v; want aborted with %v", completion, want)
	}
	select {
	case <-first.ready:
	default:
		t.Fatal("did not publish successor after local cleanup")
	}
	if err := second.await(context.Background()); err != nil {
		t.Fatalf("successor await() = %v; want nil", err)
	}
	<-lim.gate.waitCalled
	select {
	case completion := <-lim.gate.completions:
		t.Fatalf("duplicate completion = %+v", completion)
	default:
	}
	select {
	case <-lim.gate.waitCalled:
		t.Fatal("successor invoked GateNext permission more than once")
	default:
	}

	second.finish()
	if release := flow.close(); release != lim {
		t.Fatalf("close release = %v; want limiter", release)
	}
	lim.Release()
}

func TestAdmitPreparedSendPanicPoisonsAndPublishes(t *testing.T) {
	lim := newTestGateLimiter()
	flow := newClientFlow(lim)
	first := flow.ticketCurrent()
	second := flow.ticketCurrent()
	prepared := &testPreparedSend{
		size: 1,
		commit: func(func(flowcontrol.MessageOutcomeKind, error)) (*Answer, ReleaseFunc, error) {
			panic("commit panic")
		},
	}

	func() {
		defer func() {
			if got := recover(); got != "commit panic" {
				t.Fatalf("panic = %v; want commit panic", got)
			}
		}()
		_, _, _ = admitPreparedSend(context.Background(), first, prepared)
	}()
	if got := prepared.abortCount(); got != 0 {
		t.Fatalf("Abort called %d times after ambiguous panic; want 0", got)
	}
	completion := <-lim.gate.completions
	if completion.kind != flowcontrol.MessageOutcomeFatal {
		t.Fatalf("completion kind = %v; want fatal", completion.kind)
	}
	select {
	case <-first.ready:
	default:
		t.Fatal("did not publish successor after panic")
	}
	if err := second.await(context.Background()); err != nil {
		t.Fatalf("successor await() = %v; want nil", err)
	}
	<-lim.gate.waitCalled
	select {
	case completion := <-lim.gate.completions:
		t.Fatalf("duplicate completion = %+v", completion)
	default:
	}
	select {
	case <-lim.gate.waitCalled:
		t.Fatal("successor invoked GateNext permission more than once")
	default:
	}

	second.finish()
	if release := flow.close(); release != lim {
		t.Fatalf("close release = %v; want limiter", release)
	}
	lim.Release()
}
