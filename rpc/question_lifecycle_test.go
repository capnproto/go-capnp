package rpc

import (
	"context"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/exp/spsc"
	transportpkg "capnproto.org/go/capnp/v3/rpc/transport"
	rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"
)

type gatedQuestionCaller struct {
	inner    *question
	admitted chan struct{}
	release  <-chan struct{}
}

func (g *gatedQuestionCaller) PipelineSend(ctx context.Context, transform []capnp.PipelineOp, s capnp.Send) (*capnp.Answer, capnp.ReleaseFunc) {
	g.admitted <- struct{}{}
	<-g.release
	return g.inner.PipelineSend(ctx, transform, s)
}

func (g *gatedQuestionCaller) PipelineRecv(ctx context.Context, transform []capnp.PipelineOp, r capnp.Recv) capnp.PipelineCaller {
	g.admitted <- struct{}{}
	<-g.release
	return g.inner.PipelineRecv(ctx, transform, r)
}

type lifecycleResolver struct{ result chan<- error }

func (r lifecycleResolver) Fulfill(capnp.Ptr) { r.result <- nil }
func (r lifecycleResolver) Reject(err error)  { r.result <- err }

type questionLifecycleFixture struct {
	conn           *Conn
	peer           Transport
	q              *question
	ans            *capnp.Answer
	admitted       chan struct{}
	gate           chan struct{}
	gateOnce       sync.Once
	resolverResult chan error
}

type countingIncomingMessage struct {
	message  rpccp.Message
	releases int32
	once     sync.Once
}

func (m *countingIncomingMessage) Message() rpccp.Message { return m.message }
func (m *countingIncomingMessage) Release() {
	atomic.AddInt32(&m.releases, 1)
	m.once.Do(func() { m.message.Message().Release() })
}

func newQuestionReturnMessage(t *testing.T, noFinishNeeded, withCapability bool) *countingIncomingMessage {
	t.Helper()
	_, seg := capnp.NewSingleSegmentMessage(nil)
	message, err := rpccp.NewRootMessage(seg)
	if err != nil {
		t.Fatal(err)
	}
	ret, err := message.NewReturn()
	if err != nil {
		t.Fatal(err)
	}
	ret.SetAnswerId(0)
	ret.SetNoFinishNeeded(noFinishNeeded)
	results, err := ret.NewResults()
	if err != nil {
		t.Fatal(err)
	}
	if withCapability {
		caps, err := results.NewCapTable(1)
		if err != nil {
			t.Fatal(err)
		}
		caps.At(0).SetSenderHosted(7)
		if err := results.SetContent(capnp.NewInterface(results.Segment(), 0).ToPtr()); err != nil {
			t.Fatal(err)
		}
	}
	return &countingIncomingMessage{message: message}
}

func newQuestionLifecycleFixture(t *testing.T) *questionLifecycleFixture {
	t.Helper()
	left, right := transportpkg.NewPipe(16)
	f := &questionLifecycleFixture{
		conn:           NewConn(NewTransport(left), nil),
		peer:           NewTransport(right),
		admitted:       make(chan struct{}, 1),
		gate:           make(chan struct{}),
		resolverResult: make(chan error, 1),
	}
	f.conn.withLocked(func(c *lockedConn) {
		callSendDone := make(chan struct{})
		close(callSendDone)
		f.q = &question{
			c:            f.conn,
			release:      func() {},
			callSend:     questionCallSendSucceeded,
			callSendDone: callSendDone,
			drainDone:    make(chan struct{}),
		}
		f.q.id = c.lk.questions.Add(f.q)
		gate := &gatedQuestionCaller{inner: f.q, admitted: f.admitted, release: f.gate}
		f.q.p = capnp.NewPromise(capnp.Method{}, gate, lifecycleResolver{result: f.resolverResult})
		c.setAnswerQuestion(f.q.p.Answer(), f.q)
		f.ans = f.q.p.Answer()
	})
	t.Cleanup(func() {
		f.releaseGate()
		_ = f.conn.Close()
		f.q.p.ReleaseClients()
		f.q.release()
		_ = f.peer.Close()
	})
	return f
}

func (f *questionLifecycleFixture) releaseGate() {
	f.gateOnce.Do(func() { close(f.gate) })
}

func (f *questionLifecycleFixture) startAdmittedCall() <-chan struct{} {
	returned := make(chan struct{})
	go func() {
		defer close(returned)
		_, _ = f.ans.PipelineSend(context.Background(), nil, capnp.Send{Method: capnp.Method{}})
	}()
	<-f.admitted
	return returned
}

func sendQuestionReturn(t *testing.T, peer Transport, noFinishNeeded, withCapability bool) {
	t.Helper()
	out, err := peer.NewMessage()
	if err != nil {
		t.Fatal(err)
	}
	defer out.Release()
	ret, err := out.Message().NewReturn()
	if err != nil {
		t.Fatal(err)
	}
	ret.SetAnswerId(0)
	ret.SetNoFinishNeeded(noFinishNeeded)
	results, err := ret.NewResults()
	if err != nil {
		t.Fatal(err)
	}
	if withCapability {
		caps, err := results.NewCapTable(1)
		if err != nil {
			t.Fatal(err)
		}
		caps.At(0).SetSenderHosted(7)
		if err := results.SetContent(capnp.NewInterface(results.Segment(), 0).ToPtr()); err != nil {
			t.Fatal(err)
		}
	}
	if err := out.Send(); err != nil {
		t.Fatal(err)
	}
}

func sendBootstrapMarker(t *testing.T, peer Transport, id uint32) {
	t.Helper()
	out, err := peer.NewMessage()
	if err != nil {
		t.Fatal(err)
	}
	defer out.Release()
	bootstrap, err := out.Message().NewBootstrap()
	if err != nil {
		t.Fatal(err)
	}
	bootstrap.SetQuestionId(id)
	if err := out.Send(); err != nil {
		t.Fatal(err)
	}
}

func recvPeerMessage(t *testing.T, peer Transport) transportpkg.IncomingMessage {
	t.Helper()
	type result struct {
		in  transportpkg.IncomingMessage
		err error
	}
	resultCh := make(chan result, 1)
	go func() {
		in, err := peer.RecvMessage()
		resultCh <- result{in: in, err: err}
	}()
	select {
	case result := <-resultCh:
		if result.err != nil {
			t.Fatal(result.err)
		}
		return result.in
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for peer message")
		return nil
	}
}

func assertQuestionIDReusable(t *testing.T, conn *Conn, want questionID) {
	t.Helper()
	conn.withLocked(func(c *lockedConn) {
		id := c.lk.questions.Add(nil)
		if id != want {
			t.Errorf("next question ID = %d; want %d", id, want)
		}
		c.lk.questions.Remove(id)
	})
}

func TestQuestionReturnDrainsAdmittedCallBeforeFinish(t *testing.T) {
	f := newQuestionLifecycleFixture(t)
	const callCount = 4
	callReturned := make([]<-chan struct{}, callCount)
	for i := range callReturned {
		callReturned[i] = f.startAdmittedCall()
	}
	incoming := newQuestionReturnMessage(t, false, false)
	if err := f.conn.handleReturn(f.conn.bgctx, incoming); err != nil {
		t.Fatal(err)
	}
	if err := <-f.resolverResult; err != nil {
		t.Fatalf("Promise resolver error = %v; want nil", err)
	}
	f.conn.withLocked(func(c *lockedConn) {
		if f.q.owner != questionOwnerReturn || f.q.drainComplete || f.q.finish != questionFinishNotQueued {
			t.Errorf("question advanced before admitted call yielded: %+v", *f.q)
		}
	})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	f.q.handleCancel(ctx)
	f.conn.withLocked(func(c *lockedConn) {
		if f.q.owner != questionOwnerReturn {
			t.Errorf("cancellation stole terminal ownership after Return: owner %d", f.q.owner)
		}
	})

	f.releaseGate()
	for _, returned := range callReturned {
		<-returned
	}
	for i := 0; i < callCount; i++ {
		call := recvPeerMessage(t, f.peer)
		if call.Message().Which() != rpccp.Message_Which_call {
			t.Fatalf("message %d = %v; want Call", i, call.Message().Which())
		}
		call.Release()
	}
	finishMessage := recvPeerMessage(t, f.peer)
	if finishMessage.Message().Which() != rpccp.Message_Which_finish {
		t.Fatalf("message after Calls = %v; want Finish", finishMessage.Message().Which())
	}
	finish, err := finishMessage.Message().Finish()
	if err != nil {
		t.Fatal(err)
	}
	if finish.ReleaseResultCaps() {
		t.Error("normal Return sent Finish(releaseResultCaps=true)")
	}
	finishMessage.Release()

	// A later outbound marker proves the Finish callback has run before ID
	// reuse is checked, because the sender loop invokes callbacks serially.
	sendBootstrapMarker(t, f.peer, 99)
	marker := recvPeerMessage(t, f.peer)
	if marker.Message().Which() != rpccp.Message_Which_return {
		t.Fatalf("marker response = %v; want Return", marker.Message().Which())
	}
	marker.Release()
	assertQuestionIDReusable(t, f.conn, 0)
	if got := atomic.LoadInt32(&incoming.releases); got != 0 {
		t.Fatalf("successful Return released incoming message before response release: %d", got)
	}
	f.q.release()
	f.q.release = func() {}
	if got := atomic.LoadInt32(&incoming.releases); got != 1 {
		t.Fatalf("incoming message release count = %d; want 1", got)
	}
}

func TestQuestionCancellationDrainsBeforeSingleFinish(t *testing.T) {
	f := newQuestionLifecycleFixture(t)
	callReturned := f.startAdmittedCall()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	cancelReturned := make(chan struct{})
	go func() {
		f.q.handleCancel(ctx)
		close(cancelReturned)
	}()
	if err := <-f.resolverResult; !errors.Is(err, context.Canceled) {
		t.Fatalf("Promise resolver error = %v; want context.Canceled", err)
	}
	f.conn.withLocked(func(c *lockedConn) {
		if f.q.owner != questionOwnerCancel || f.q.drainComplete || f.q.finish != questionFinishNotQueued {
			t.Errorf("cancellation advanced before admitted call yielded: %+v", *f.q)
		}
	})

	f.releaseGate()
	<-callReturned
	<-cancelReturned
	call := recvPeerMessage(t, f.peer)
	if call.Message().Which() != rpccp.Message_Which_call {
		t.Fatalf("first message = %v; want Call", call.Message().Which())
	}
	call.Release()
	finishMessage := recvPeerMessage(t, f.peer)
	finish, err := finishMessage.Message().Finish()
	if err != nil {
		t.Fatal(err)
	}
	if !finish.ReleaseResultCaps() {
		t.Error("cancellation sent Finish(releaseResultCaps=false)")
	}
	finishMessage.Release()

	// A late contradictory hint cannot suppress or duplicate cancellation's
	// already-owned Finish(true).
	lateReturn := newQuestionReturnMessage(t, true, false)
	if err := f.conn.handleReturn(f.conn.bgctx, lateReturn); err != nil {
		t.Fatal(err)
	}
	if got := atomic.LoadInt32(&lateReturn.releases); got != 1 {
		t.Fatalf("late canceled Return release count = %d; want 1", got)
	}
	sendBootstrapMarker(t, f.peer, 99)
	marker := recvPeerMessage(t, f.peer)
	if marker.Message().Which() != rpccp.Message_Which_return {
		t.Fatalf("message after late Return = %v; want marker Return (no second Finish)", marker.Message().Which())
	}
	marker.Release()
	f.conn.withLocked(func(c *lockedConn) {
		if !f.q.returnReceived || f.q.finish != questionFinishSent {
			t.Errorf("late Return state = returnReceived %t, finish %d", f.q.returnReceived, f.q.finish)
		}
	})
	assertQuestionIDReusable(t, f.conn, 0)
}

func TestQuestionNoFinishNeededSuppressesAndReusesID(t *testing.T) {
	f := newQuestionLifecycleFixture(t)
	incoming := newQuestionReturnMessage(t, true, false)
	if err := f.conn.handleReturn(f.conn.bgctx, incoming); err != nil {
		t.Fatal(err)
	}
	if err := <-f.resolverResult; err != nil {
		t.Fatalf("Promise resolver error = %v; want nil", err)
	}
	<-f.q.drainDone
	f.conn.withLocked(func(c *lockedConn) {
		if !f.q.noFinishNeeded || !f.q.drainComplete || f.q.finish != questionFinishSuppressed {
			t.Errorf("noFinishNeeded state = hint %t, drain %t, finish %d",
				f.q.noFinishNeeded, f.q.drainComplete, f.q.finish)
		}
	})

	// Once the drain is complete, a marker Return must be the next outbound
	// message: a Finish would contradict the valid hint.
	sendBootstrapMarker(t, f.peer, 99)
	marker := recvPeerMessage(t, f.peer)
	if marker.Message().Which() != rpccp.Message_Which_return {
		t.Fatalf("message after noFinishNeeded drain = %v; want marker Return", marker.Message().Which())
	}
	marker.Release()
	assertQuestionIDReusable(t, f.conn, 0)
	if got := atomic.LoadInt32(&incoming.releases); got != 0 {
		t.Fatalf("noFinishNeeded Return released response before ReleaseFunc: %d", got)
	}
	f.q.release()
	f.q.release = func() {}
	if got := atomic.LoadInt32(&incoming.releases); got != 1 {
		t.Fatalf("noFinishNeeded incoming release count = %d; want 1", got)
	}
}

func TestQuestionContradictoryNoFinishNeededSendsFinish(t *testing.T) {
	f := newQuestionLifecycleFixture(t)
	logger := new(recordingLogger)
	f.conn.er = errReporter{Logger: logger}
	sendQuestionReturn(t, f.peer, true, true)
	if err := <-f.resolverResult; err != nil {
		t.Fatalf("Promise resolver error = %v; want nil", err)
	}
	finishMessage := recvPeerMessage(t, f.peer)
	if finishMessage.Message().Which() != rpccp.Message_Which_finish {
		t.Fatalf("message = %v; want Finish", finishMessage.Message().Which())
	}
	finish, err := finishMessage.Message().Finish()
	if err != nil {
		t.Fatal(err)
	}
	if finish.ReleaseResultCaps() {
		t.Error("capability-bearing Return sent Finish(releaseResultCaps=true)")
	}
	finishMessage.Release()

	logger.mu.Lock()
	loggedErrors := append([]string(nil), logger.errors...)
	logger.mu.Unlock()
	if len(loggedErrors) != 1 || !strings.Contains(loggedErrors[0], "noFinishNeeded") {
		t.Fatalf("logged errors = %q; want one noFinishNeeded diagnostic", loggedErrors)
	}
}

func TestQuestionShutdownWaitsForReturnDrain(t *testing.T) {
	f := newQuestionLifecycleFixture(t)
	callReturned := f.startAdmittedCall()
	incoming := newQuestionReturnMessage(t, false, false)
	if err := f.conn.handleReturn(f.conn.bgctx, incoming); err != nil {
		t.Fatal(err)
	}
	if err := <-f.resolverResult; err != nil {
		t.Fatalf("Promise resolver error = %v; want nil", err)
	}

	closeReturned := make(chan error, 1)
	go func() { closeReturned <- f.conn.Close() }()
	<-f.conn.bgctx.Done()
	select {
	case err := <-closeReturned:
		t.Fatalf("Close returned before admitted call yielded: %v", err)
	default:
	}
	f.conn.withLocked(func(c *lockedConn) {
		if f.q.finish != questionFinishNotQueued {
			t.Errorf("Finish state during blocked drain = %d; want not queued", f.q.finish)
		}
	})

	f.releaseGate()
	<-callReturned
	if err := <-closeReturned; err != nil {
		t.Fatal(err)
	}
	if f.q.finish != questionFinishFailed {
		t.Errorf("Finish state after shutdown = %d; want failed", f.q.finish)
	}
	if got := atomic.LoadInt32(&incoming.releases); got != 0 {
		t.Fatalf("shutdown released successful response still owned by caller: %d", got)
	}
	f.q.release()
	f.q.release = func() {}
	if got := atomic.LoadInt32(&incoming.releases); got != 1 {
		t.Fatalf("shutdown response release count = %d; want 1", got)
	}
}

func TestQuestionReturnCancelRaceHasSingleOwner(t *testing.T) {
	f := newQuestionLifecycleFixture(t)
	incoming := newQuestionReturnMessage(t, false, false)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	start := make(chan struct{})
	returnResult := make(chan error, 1)
	cancelReturned := make(chan struct{})
	go func() {
		<-start
		returnResult <- f.conn.handleReturn(f.conn.bgctx, incoming)
	}()
	go func() {
		<-start
		f.q.handleCancel(ctx)
		close(cancelReturned)
	}()
	close(start)
	if err := <-returnResult; err != nil {
		t.Fatal(err)
	}
	<-cancelReturned
	resolvedErr := <-f.resolverResult
	<-f.q.drainDone

	f.conn.withLocked(func(c *lockedConn) {
		switch f.q.owner {
		case questionOwnerReturn:
			if resolvedErr != nil {
				t.Errorf("Return winner rejected Promise: %v", resolvedErr)
			}
		case questionOwnerCancel:
			if !errors.Is(resolvedErr, context.Canceled) {
				t.Errorf("cancellation winner resolved with %v; want context.Canceled", resolvedErr)
			}
		default:
			t.Errorf("terminal owner = %d; want Return or cancellation", f.q.owner)
		}
		if !f.q.returnReceived || (f.q.finish != questionFinishQueued && f.q.finish != questionFinishSent) {
			t.Errorf("race state = returnReceived %t, finish %d", f.q.returnReceived, f.q.finish)
		}
	})

	finishMessage := recvPeerMessage(t, f.peer)
	finish, err := finishMessage.Message().Finish()
	if err != nil {
		t.Fatal(err)
	}
	if got := finish.ReleaseResultCaps(); got != (f.q.owner == questionOwnerCancel) {
		t.Errorf("Finish.releaseResultCaps = %t; owner %d", got, f.q.owner)
	}
	finishMessage.Release()
	sendBootstrapMarker(t, f.peer, 99)
	marker := recvPeerMessage(t, f.peer)
	if marker.Message().Which() != rpccp.Message_Which_return {
		t.Fatalf("message after race Finish = %v; want marker Return", marker.Message().Which())
	}
	marker.Release()
	assertQuestionIDReusable(t, f.conn, 0)

	beforeRelease := atomic.LoadInt32(&incoming.releases)
	f.q.release()
	f.q.release = func() {}
	if f.q.owner == questionOwnerReturn {
		if beforeRelease != 0 || atomic.LoadInt32(&incoming.releases) != 1 {
			t.Errorf("Return winner incoming releases = %d before, %d after response release",
				beforeRelease, atomic.LoadInt32(&incoming.releases))
		}
	} else if beforeRelease != 1 || atomic.LoadInt32(&incoming.releases) != 1 {
		t.Errorf("cancellation winner incoming release count = %d before, %d after", beforeRelease,
			atomic.LoadInt32(&incoming.releases))
	}
}

type unusedPipelineCaller struct{}

func (unusedPipelineCaller) PipelineSend(context.Context, []capnp.PipelineOp, capnp.Send) (*capnp.Answer, capnp.ReleaseFunc) {
	panic("unexpected pipeline send")
}

func (unusedPipelineCaller) PipelineRecv(context.Context, []capnp.PipelineOp, capnp.Recv) capnp.PipelineCaller {
	panic("unexpected pipeline recv")
}

func newUnsentQuestion() (*Conn, *question) {
	c := &Conn{bgctx: context.Background()}
	q := &question{c: c, release: func() {}, callSendDone: make(chan struct{})}
	q.id = (*lockedConn)(c).lk.questions.Add(q)
	q.p = capnp.NewPromise(capnp.Method{}, unusedPipelineCaller{}, nil)
	return c, q
}

func TestQuestionInitialCallFailureIDPolicy(t *testing.T) {
	t.Run("definitely unsent releases ID", func(t *testing.T) {
		c, q := newUnsentQuestion()
		queue := spsc.New[asyncSend]()
		c.lk.sendTx = &queue.Tx
		admitted := make(chan struct{}, 1)
		gateRelease := make(chan struct{})
		q.p = capnp.NewPromise(capnp.Method{}, &gatedQuestionCaller{
			inner: q, admitted: admitted, release: gateRelease,
		}, nil)
		pipelineReturned := make(chan struct{})
		go func() {
			_, _ = q.p.Answer().PipelineSend(context.Background(), nil, capnp.Send{Method: capnp.Method{}})
			close(pipelineReturned)
		}()
		<-admitted

		handlerReturned := make(chan struct{})
		go func() {
			q.handleCallSend(context.Background(), sendOutcome{
				disposition: sendDefinitelyUnsent,
				err:         errors.New("build failed"),
			}, "test call")
			close(handlerReturned)
		}()
		<-q.callSendDone
		select {
		case <-handlerReturned:
			t.Fatal("definitely-unsent handler did not drain admitted pipeline call")
		default:
		}
		close(gateRelease)
		<-pipelineReturned
		<-handlerReturned
		<-q.p.Answer().Done()
		if _, ok := queue.Rx.TryRecv(); ok {
			t.Fatal("dependent Call escaped after parent Call was definitely unsent")
		}
		assertQuestionIDReusable(t, c, q.id)
	})

	t.Run("delivery ambiguous retains ID", func(t *testing.T) {
		c, q := newUnsentQuestion()
		q.handleCallSend(context.Background(), sendOutcome{
			disposition: sendDeliveryAmbiguous,
			err:         errors.New("write failed"),
			fatal:       true,
		}, "test call")
		<-q.p.Answer().Done()
		if q.owner != questionOwnerSendFailure {
			t.Fatalf("terminal owner = %d; want send failure", q.owner)
		}
		c.withLocked(func(c *lockedConn) {
			if id := c.lk.questions.Add(nil); id == q.id {
				t.Errorf("delivery-ambiguous question ID %d was reused", q.id)
			}
		})
	})

	t.Run("shutdown abort leaves cleanup to shutdown", func(t *testing.T) {
		c, q := newUnsentQuestion()
		q.handleCallSend(context.Background(), sendOutcome{
			disposition: sendAbortedByShutdown,
			err:         ErrConnClosed,
		}, "test call")
		select {
		case <-q.p.Answer().Done():
			t.Fatal("queue abort independently rejected question")
		default:
		}
		c.withLocked(func(c *lockedConn) {
			if current, ok := c.lk.questions.Find(q.id); !ok || current != q {
				t.Fatal("queue abort removed question before shutdown cleanup")
			}
		})
	})
}

func TestQuestionFinishFailureRetainsID(t *testing.T) {
	tests := []struct {
		name      string
		transport Transport
		abort     bool
	}{
		{
			name:      "definitely unsent",
			transport: newMessageErrorTransport{err: errors.New("allocate failed")},
		},
		{
			name: "delivery ambiguous",
			transport: func() Transport {
				tr := newFailingSendTransport(errors.New("write failed"))
				close(tr.firstSend)
				return tr
			}(),
		},
		{
			name:      "shutdown abort",
			transport: newMessageErrorTransport{err: errors.New("unused")},
			abort:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			queue := spsc.New[asyncSend]()
			c := &Conn{transport: test.transport, bgctx: context.Background()}
			c.lk.sendTx = &queue.Tx
			q := &question{
				c:              c,
				release:        func() {},
				owner:          questionOwnerReturn,
				drainComplete:  true,
				returnReceived: true,
			}
			q.id = (*lockedConn)(c).lk.questions.Add(q)
			(*lockedConn)(c).lk.questions.take(q.id)
			q.queueFinish((*lockedConn)(c), false)

			pending, ok := queue.Rx.TryRecv()
			if !ok {
				t.Fatal("Finish was not queued")
			}
			if test.abort {
				pending.Abort(ErrConnClosed)
			} else if err := pending.Send(); err == nil {
				t.Fatal("Finish failure did not return fatal sender error")
			}
			if q.finish != questionFinishFailed {
				t.Fatalf("Finish state = %d; want failed", q.finish)
			}
			c.withLocked(func(c *lockedConn) {
				if id := c.lk.questions.Add(nil); id == q.id {
					t.Errorf("failed Finish released question ID %d", q.id)
				}
			})
		})
	}
}
