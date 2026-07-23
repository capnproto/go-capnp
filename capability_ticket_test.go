package capnp

import (
	"context"
	"sync"
	"testing"

	"capnproto.org/go/capnp/v3/flowcontrol"
)

type ticketTestLimiter struct {
	releases int
}

func (*ticketTestLimiter) StartMessage(context.Context, uint64) (func(), error) {
	return func() {}, nil
}

func (l *ticketTestLimiter) Release() {
	l.releases++
}

var _ flowcontrol.FlowLimiter = (*ticketTestLimiter)(nil)

type observedAwaitContext struct {
	context.Context
	entered chan struct{}
	once    sync.Once
}

func (c *observedAwaitContext) Done() <-chan struct{} {
	c.once.Do(func() { close(c.entered) })
	return c.Context.Done()
}

type steppedAwaitContext struct {
	context.Context
	steps chan struct{}
}

func (c *steppedAwaitContext) Done() <-chan struct{} {
	c.steps <- struct{}{}
	return c.Context.Done()
}

func TestFlowTicketAwaitSameGenerationUsesPredecessorGate(t *testing.T) {
	lim := new(ticketTestLimiter)
	flow := newClientFlow(lim)
	first := flow.ticketCurrent()
	second := flow.ticketCurrent()
	called := make(chan struct{}, 1)
	first.publish(func(context.Context) error {
		called <- struct{}{}
		return nil
	})

	if err := second.await(context.Background()); err != nil {
		t.Fatalf("await() = %v; want nil", err)
	}
	select {
	case <-called:
	default:
		t.Fatal("await did not invoke the predecessor gate")
	}

	first.finish()
	second.finish()
	if release := flow.close(); release != lim {
		t.Fatalf("close release = %v; want current limiter", release)
	}
	if lim.releases != 0 {
		t.Fatalf("Release called %d times before caller released returned limiter", lim.releases)
	}
	lim.Release()
}

func TestFlowTicketAwaitCrossGenerationOnlyWaitsForPublication(t *testing.T) {
	old := new(ticketTestLimiter)
	replacement := new(ticketTestLimiter)
	flow := newClientFlow(old)
	first := flow.ticketCurrent()
	if release := flow.replace(replacement); release != nil {
		t.Fatal("replacement released old limiter before its ticket drained")
	}
	second := flow.ticketCurrent()

	ctx := &observedAwaitContext{Context: context.Background(), entered: make(chan struct{})}
	awaited := make(chan error, 1)
	go func() { awaited <- second.await(ctx) }()
	<-ctx.entered
	select {
	case err := <-awaited:
		t.Fatalf("await returned before predecessor publication: %v", err)
	default:
	}

	oldGateCalled := make(chan struct{}, 1)
	first.publish(func(context.Context) error {
		oldGateCalled <- struct{}{}
		return nil
	})
	if err := <-awaited; err != nil {
		t.Fatalf("await() = %v; want nil", err)
	}
	select {
	case <-oldGateCalled:
		t.Fatal("cross-generation await invoked the retired limiter gate")
	default:
	}

	first.finish()
	if old.releases != 1 {
		t.Fatalf("old Release called %d times after old ticket drained; want 1", old.releases)
	}
	second.finish()
	if release := flow.close(); release != replacement {
		t.Fatalf("close release = %v; want current limiter", release)
	}
	replacement.Release()
}

func TestFlowTicketAwaitCrossGenerationWalksAbandonedPredecessor(t *testing.T) {
	old := new(ticketTestLimiter)
	replacement := new(ticketTestLimiter)
	flow := newClientFlow(old)
	first := flow.ticketCurrent()
	abandoned := flow.ticketCurrent()
	abandoned.abandon()
	abandoned.finish()
	if release := flow.replace(replacement); release != nil {
		t.Fatal("replacement released old limiter before its first ticket drained")
	}
	successor := flow.ticketCurrent()

	ctx := &steppedAwaitContext{
		Context: context.Background(),
		steps:   make(chan struct{}),
	}
	awaited := make(chan error, 1)
	go func() { awaited <- successor.await(ctx) }()

	// The first step observes the already-published abandoned ticket. The
	// second proves await walked through it and reached first.ready.
	<-ctx.steps
	select {
	case <-ctx.steps:
	case err := <-awaited:
		t.Fatalf("await returned before transitive predecessor publication: %v", err)
	}

	oldGateCalled := make(chan struct{}, 1)
	first.publish(func(context.Context) error {
		oldGateCalled <- struct{}{}
		return nil
	})
	if err := <-awaited; err != nil {
		t.Fatalf("await() = %v; want nil", err)
	}
	select {
	case <-oldGateCalled:
		t.Fatal("cross-generation await invoked the retired limiter gate")
	default:
	}

	first.finish()
	if old.releases != 1 {
		t.Fatalf("old Release called %d times after old tickets drained; want 1", old.releases)
	}
	successor.finish()
	if release := flow.close(); release != replacement {
		t.Fatalf("close release = %v; want current limiter", release)
	}
	replacement.Release()
}

func TestFlowTicketAwaitSameGenerationWalksAbandonedPredecessor(t *testing.T) {
	lim := new(ticketTestLimiter)
	flow := newClientFlow(lim)
	first := flow.ticketCurrent()
	abandoned := flow.ticketCurrent()
	abandoned.abandon()
	abandoned.finish()
	successor := flow.ticketCurrent()

	gateCalled := make(chan struct{}, 1)
	first.publish(func(context.Context) error {
		gateCalled <- struct{}{}
		return nil
	})
	if err := successor.await(context.Background()); err != nil {
		t.Fatalf("await() = %v; want nil", err)
	}
	select {
	case <-gateCalled:
	default:
		t.Fatal("same-generation successor skipped transitive predecessor gate")
	}

	first.finish()
	successor.finish()
	if release := flow.close(); release != lim {
		t.Fatalf("close release = %v; want limiter", release)
	}
	lim.Release()
}
