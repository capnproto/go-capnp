package rpc

import (
	"context"
	"strings"
	"sync/atomic"
	"testing"
	"testing/synctest"

	"capnproto.org/go/capnp/v3"
)

type cyclePipelineCaller struct{}

func (cyclePipelineCaller) PipelineSend(
	context.Context,
	[]capnp.PipelineOp,
	capnp.Send,
) (*capnp.Answer, capnp.ReleaseFunc) {
	panic("unexpected pipeline send")
}

func (cyclePipelineCaller) PipelineRecv(
	context.Context,
	[]capnp.PipelineOp,
	capnp.Recv,
) capnp.PipelineCaller {
	panic("unexpected pipeline recv")
}

func TestHandleResolveRejectsDelayedReceiverAnswerCycle(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		conn, _ := newResolveLifecycleConn(t)
		promise := addResolvePromise(t, conn)
		promiseSnapshot := promise.Snapshot()
		defer promiseSnapshot.Release()

		// Preserve a cursor rooted at the original promise hook.  Unlike
		// Client.AddRef, ClientSnapshot.Client creates an independent cursor.
		resultMessage, resultSegment := capnp.NewSingleSegmentMessage(nil)
		resultClient := promiseSnapshot.Client()
		resultCap := resultMessage.CapTable().Add(resultClient)
		result := capnp.NewInterface(resultSegment, resultCap).ToPtr()
		defer resultMessage.Release()

		answerPromise := capnp.NewPromise(capnp.Method{}, cyclePipelineCaller{}, nil)
		const delayedAnswerID = answerID(45)
		conn.withLocked(func(c *lockedConn) {
			// Mark this pipeline as belonging to the same connection.  This
			// ensures isLocalClient does not hide the promise graph behind an
			// embargo, and the assertion below pins that prerequisite.
			c.setAnswerQuestion(answerPromise.Answer(), &question{c: conn})
			if !c.lk.answers.Create(delayedAnswerID, &ansent{promise: answerPromise}) {
				t.Fatal("create delayed receiver answer")
			}
		})
		defer func() {
			conn.withLocked(func(c *lockedConn) {
				c.lk.answers.Remove(delayedAnswerID)
			})
		}()

		in := resolveToReceiverAnswer(t, resolvePromiseID, delayedAnswerID)
		if err := conn.handleResolve(context.Background(), in); err != nil {
			t.Fatal("handle Resolve:", err)
		}
		if got := atomic.LoadInt32(&in.releases); got != 1 {
			t.Fatalf("incoming message releases = %d; want 1", got)
		}
		requirePromiseResolverConsumed(t, conn, resolvePromiseID)
		conn.withLocked(func(c *lockedConn) {
			embargoes := 0
			c.lk.embargoes.Range(func(_ embargoID, _ *embargo) bool {
				embargoes++
				return true
			})
			if embargoes != 0 {
				t.Fatalf("delayed receiver answer created %d embargoes; want 0", embargoes)
			}
		})

		fulfilled := make(chan struct{})
		released := make(chan struct{})
		go func() {
			answerPromise.Fulfill(result)
			close(fulfilled)
		}()
		go func() {
			promise.Release()
			close(released)
		}()
		synctest.Wait()
		<-fulfilled
		<-released
		answerPromise.ReleaseClients()

		if err := promiseSnapshot.Resolve(context.Background()); err != nil {
			t.Fatal("resolve imported promise snapshot:", err)
		}
		resolutionErr, ok := promiseSnapshot.Brand().Value.(error)
		if !ok || !strings.Contains(resolutionErr.Error(), "client promise resolution cycle") {
			t.Fatalf("resolved promise brand = %v; want resolution-cycle error", resolutionErr)
		}

		lockObserved := false
		conn.withLocked(func(*lockedConn) {
			lockObserved = true
		})
		if !lockObserved {
			t.Fatal("connection lock remained unavailable after cycle rejection")
		}
	})
}
