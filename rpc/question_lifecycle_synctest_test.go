package rpc_test

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc"
	rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"
)

const (
	lifecycleBootstrapExportID  = 80
	lifecycleLateResultExportID = 81
)

type questionTraceFixture struct {
	observed *lifecycleObservedTransport
	peer     rpc.Transport
	conn     *rpc.Conn
	target   capnp.Client
}

func newQuestionTraceFixture(t *testing.T) *questionTraceFixture {
	t.Helper()
	observed, peer := newLifecycleDriver(t, 16)
	conn := rpc.NewConn(observed, &rpc.Options{
		Logger: testErrorReporter{tb: t, fail: true},
	})
	fixture := &questionTraceFixture{
		observed: observed,
		peer:     peer,
		conn:     conn,
	}
	fixture.target = fixture.bootstrapTarget(t)
	t.Cleanup(func() {
		fixture.target.Release()
		// Let the sender acknowledge the final import Release before Close
		// cancels the connection context.
		synctest.Wait()
		if err := fixture.conn.Close(); err != nil {
			t.Error("close lifecycle connection:", err)
		}
	})
	return fixture
}

func (f *questionTraceFixture) bootstrapTarget(t *testing.T) capnp.Client {
	t.Helper()
	client := f.conn.Bootstrap(context.Background())
	message, release, err := recvMessage(context.Background(), f.peer)
	if err != nil {
		t.Fatal("receive bootstrap:", err)
	}
	if message.Which != rpccp.Message_Which_bootstrap {
		release()
		t.Fatalf("bootstrap setup message = %v; want Bootstrap", message.Which)
	}
	questionID := message.Bootstrap.QuestionID
	release()

	sendLifecycleReturn(t, f.peer, questionID, false, lifecycleBootstrapExportID)
	if err := client.Resolve(context.Background()); err != nil {
		client.Release()
		t.Fatal("resolve bootstrap target:", err)
	}
	finish, release, err := recvMessage(context.Background(), f.peer)
	if err != nil {
		client.Release()
		t.Fatal("receive bootstrap Finish:", err)
	}
	if finish.Which != rpccp.Message_Which_finish ||
		finish.Finish.QuestionID != questionID ||
		finish.Finish.ReleaseResultCaps {
		release()
		client.Release()
		t.Fatalf("bootstrap Finish = %+v; want Finish(%d, false)", finish, questionID)
	}
	release()
	synctest.Wait()
	return client
}

func (f *questionTraceFixture) startCall(
	t *testing.T,
	ctx context.Context,
) (*capnp.Answer, capnp.ReleaseFunc, uint32) {
	t.Helper()
	answer, releaseAnswer := f.target.SendCall(ctx, capnp.Send{
		Method: capnp.Method{InterfaceID: interfaceID, MethodID: methodID},
	})
	message, releaseMessage, err := recvMessage(context.Background(), f.peer)
	if err != nil {
		t.Fatal("receive Call:", err)
	}
	if message.Which != rpccp.Message_Which_call {
		releaseMessage()
		t.Fatalf("outbound message = %v; want Call", message.Which)
	}
	questionID := message.Call.QuestionID
	releaseMessage()
	return answer, releaseAnswer, questionID
}

// sendLifecycleReturn sends a results Return.  exportID == 0 means that the
// result is capability-free; otherwise content points at one senderHosted
// capability with the supplied export ID.
func sendLifecycleReturn(
	t *testing.T,
	peer rpc.Transport,
	questionID uint32,
	noFinishNeeded bool,
	exportID uint32,
) {
	t.Helper()
	out, err := peer.NewMessage()
	if err != nil {
		t.Fatal("allocate Return:", err)
	}
	defer out.Release()
	ret, err := out.Message().NewReturn()
	if err != nil {
		t.Fatal("build Return:", err)
	}
	ret.SetAnswerId(questionID)
	ret.SetNoFinishNeeded(noFinishNeeded)
	results, err := ret.NewResults()
	if err != nil {
		t.Fatal("build Return results:", err)
	}
	if exportID != 0 {
		caps, err := results.NewCapTable(1)
		if err != nil {
			t.Fatal("build Return cap table:", err)
		}
		caps.At(0).SetSenderHosted(exportID)
		if err := results.SetContent(capnp.NewInterface(results.Segment(), 0).ToPtr()); err != nil {
			t.Fatal("set Return capability content:", err)
		}
	}
	if err := out.Send(); err != nil {
		t.Fatal("send Return:", err)
	}
}

func receiveLifecycleFinish(
	t *testing.T,
	peer rpc.Transport,
	questionID uint32,
	releaseResultCaps bool,
) {
	t.Helper()
	message, release, err := recvMessage(context.Background(), peer)
	if err != nil {
		t.Fatal("receive Finish:", err)
	}
	defer release()
	if message.Which != rpccp.Message_Which_finish ||
		message.Finish.QuestionID != questionID ||
		message.Finish.ReleaseResultCaps != releaseResultCaps {
		t.Fatalf(
			"terminal message = %+v; want Finish(%d, %t)",
			message,
			questionID,
			releaseResultCaps,
		)
	}
}

func TestLifecycleQuestionReturnFinishThenReuse(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		fixture := newQuestionTraceFixture(t)
		oracle := newLifecycleOracle()

		firstAnswer, releaseFirst, firstWireID := fixture.startCall(t, context.Background())
		if err := oracle.localCall(); err != nil {
			t.Fatal(err)
		}
		first, err := oracle.wireCall(firstWireID)
		if err != nil {
			t.Fatal(err)
		}

		sendLifecycleReturn(t, fixture.peer, firstWireID, false, 0)
		if _, err := firstAnswer.Struct(); err != nil {
			t.Fatal("read first result:", err)
		}
		releaseFirst()
		if err := oracle.peerReturn(first, false, 0); err != nil {
			t.Fatal(err)
		}

		receiveLifecycleFinish(t, fixture.peer, firstWireID, false)
		if err := oracle.wireFinish(first); err != nil {
			t.Fatal(err)
		}
		synctest.Wait()

		secondAnswer, releaseSecond, secondWireID := fixture.startCall(t, context.Background())
		if err := oracle.localCall(); err != nil {
			t.Fatal(err)
		}
		second, err := oracle.wireCall(secondWireID)
		if err != nil {
			t.Fatal(err)
		}
		if secondWireID != firstWireID {
			t.Fatalf("later Call question ID = %d; want reused ID %d", secondWireID, firstWireID)
		}
		if second.ID.Generation != first.ID.Generation+1 {
			t.Fatalf(
				"later Call generation = %d; want %d",
				second.ID.Generation,
				first.ID.Generation+1,
			)
		}

		sendLifecycleReturn(t, fixture.peer, secondWireID, true, 0)
		if _, err := secondAnswer.Struct(); err != nil {
			t.Fatal("read cleanup result:", err)
		}
		releaseSecond()
	})
}

func TestLifecycleQuestionCancelBalancesLateCapabilityReturn(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		fixture := newQuestionTraceFixture(t)
		callCtx, cancelCall := context.WithCancel(context.Background())
		answer, releaseAnswer, questionID := fixture.startCall(t, callCtx)

		cancelCall()
		receiveLifecycleFinish(t, fixture.peer, questionID, true)
		if _, err := answer.Struct(); !errors.Is(err, context.Canceled) {
			t.Fatalf("canceled Call error = %v; want context.Canceled", err)
		}
		releaseAnswer()

		sendLifecycleReturn(
			t,
			fixture.peer,
			questionID,
			false,
			lifecycleLateResultExportID,
		)

		// The marker is ordered after the late Return on the peer->Go
		// transport.  Receiving its response acknowledges that the Return was
		// handled, without a timeout-based absence assertion.
		const markerQuestionID = 99
		if err := sendMessage(context.Background(), fixture.peer, &rpcMessage{
			Which: rpccp.Message_Which_bootstrap,
			Bootstrap: &rpcBootstrap{
				QuestionID: markerQuestionID,
			},
		}); err != nil {
			t.Fatal("send marker Bootstrap:", err)
		}
		marker, releaseMarker, err := recvMessage(context.Background(), fixture.peer)
		if err != nil {
			t.Fatal("receive marker response:", err)
		}
		if marker.Which != rpccp.Message_Which_return ||
			marker.Return.AnswerID != markerQuestionID {
			releaseMarker()
			t.Fatalf("marker response = %+v; want Return(%d)", marker, markerQuestionID)
		}
		releaseMarker()
		synctest.Wait()

		observations := fixture.observed.Observations()
		var cancelFinish, lateReturn *lifecycleWireObservation
		for i := range observations {
			observation := &observations[i]
			if observation.Flow == lifecycleGoToPeer &&
				observation.Which == rpccp.Message_Which_finish &&
				observation.QuestionID == questionID {
				cancelFinish = observation
			}
			if observation.Flow == lifecyclePeerToGo &&
				observation.Which == rpccp.Message_Which_return &&
				observation.AnswerID == questionID {
				lateReturn = observation
			}
			if observation.Flow == lifecycleGoToPeer &&
				observation.Which == rpccp.Message_Which_release &&
				observation.ImportID == lifecycleLateResultExportID {
				t.Fatalf("late result capability sent separate Release: %+v", *observation)
			}
		}
		if lateReturn == nil {
			t.Fatal("late capability Return was not observed")
		}
		if cancelFinish == nil || !cancelFinish.ReleaseResultCaps {
			t.Fatalf("cancel Finish observation = %+v; want releaseResultCaps=true", cancelFinish)
		}
		if lateReturn.ResultCapCount != 1 || lateReturn.IncomingReleaseCalls() != 1 {
			t.Fatalf("late Return observation = %+v; want one cap and one incoming release", *lateReturn)
		}

		oracle := newLifecycleOracle()
		capability := capRef{
			Export: oracle.ids.normalize(
				sidePeer,
				spaceExport,
				lifecycleLateResultExportID,
			),
		}
		input := ruleInput{Action: actionLocalCancel, Capability: capability}
		alternatives, err := buildCancelResultConstraints(input)
		if err != nil {
			t.Fatal(err)
		}
		ledger := lifecycleLedger{
			{Kind: actionLocalCancel},
			{Kind: actionObserveCompletion, Label: "canceled"},
			{
				Kind:           actionWireFinish,
				Label:          "releaseResultCaps=true",
				Capability:     capability,
				ReferenceDelta: -int64(lateReturn.ResultCapCount),
			},
			{
				Kind:           actionPeerReturn,
				Capability:     capability,
				ReferenceDelta: int64(lateReturn.ResultCapCount),
			},
		}
		if err := checkConstraintAlternatives(alternatives, ledger); err != nil {
			t.Fatal(err)
		}
	})
}

func TestLifecycleQuestionNoFinishNeededThenReuse(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		fixture := newQuestionTraceFixture(t)
		oracle := newLifecycleOracle()

		firstAnswer, releaseFirst, firstWireID := fixture.startCall(t, context.Background())
		if err := oracle.localCall(); err != nil {
			t.Fatal(err)
		}
		first, err := oracle.wireCall(firstWireID)
		if err != nil {
			t.Fatal(err)
		}

		sendLifecycleReturn(t, fixture.peer, firstWireID, true, 0)
		if _, err := firstAnswer.Struct(); err != nil {
			t.Fatal("read noFinishNeeded result:", err)
		}
		releaseFirst()
		if err := oracle.peerReturn(first, true, 0); err != nil {
			t.Fatal(err)
		}
		synctest.Wait()

		// The next outbound protocol message is the second Call itself.  This
		// proves both that no Finish was sent and that the capability-free
		// noFinishNeeded Return retired the original question.
		secondAnswer, releaseSecond, secondWireID := fixture.startCall(t, context.Background())
		if err := oracle.localCall(); err != nil {
			t.Fatal(err)
		}
		second, err := oracle.wireCall(secondWireID)
		if err != nil {
			t.Fatal(err)
		}
		if secondWireID != firstWireID {
			t.Fatalf("later Call question ID = %d; want reused ID %d", secondWireID, firstWireID)
		}
		if second.ID.Generation != first.ID.Generation+1 {
			t.Fatalf(
				"later Call generation = %d; want %d",
				second.ID.Generation,
				first.ID.Generation+1,
			)
		}

		input := ruleInput{
			Action:         actionPeerReturn,
			Question:       first,
			NoFinishNeeded: true,
		}
		alternatives, err := buildNoFinishNeededConstraints(input)
		if err != nil {
			t.Fatal(err)
		}
		if err := checkConstraintAlternatives(alternatives, lifecycleLedger{{
			Kind:     actionPeerReturn,
			Label:    "noFinishNeeded=true,resultCaps=0",
			Question: first,
		}}); err != nil {
			t.Fatal(err)
		}

		sendLifecycleReturn(t, fixture.peer, secondWireID, true, 0)
		if _, err := secondAnswer.Struct(); err != nil {
			t.Fatal("read cleanup result:", err)
		}
		releaseSecond()
	})
}
