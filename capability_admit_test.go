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
}

func (g *testGateController) CommitMessage(uint64) (func(context.Context) error, func(flowcontrol.MessageOutcomeKind, error)) {
	return func(context.Context) error {
			g.waitCalled <- struct{}{}
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
