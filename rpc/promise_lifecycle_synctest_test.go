package rpc_test

import (
	"context"
	"fmt"
	"testing"
	"testing/synctest"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc"
	"capnproto.org/go/capnp/v3/server"
	rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"
)

func sendLifecycleCapabilityReturn(
	t *testing.T,
	peer rpc.Transport,
	answerID uint32,
	descriptor func(rpccp.CapDescriptor),
) {
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
	ret.SetAnswerId(answerID)
	results, err := ret.NewResults()
	if err != nil {
		t.Fatal(err)
	}
	caps, err := results.NewCapTable(1)
	if err != nil {
		t.Fatal(err)
	}
	descriptor(caps.At(0))
	if err := results.SetContent(capnp.NewInterface(results.Segment(), 0).ToPtr()); err != nil {
		t.Fatal(err)
	}
	if err := out.Send(); err != nil {
		t.Fatal(err)
	}
}

func sendLifecycleCapabilityFieldReturn(
	t *testing.T,
	peer rpc.Transport,
	answerID uint32,
	releaseParamCaps bool,
	descriptor func(rpccp.CapDescriptor),
) {
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
	ret.SetAnswerId(answerID)
	ret.SetReleaseParamCaps(releaseParamCaps)
	results, err := ret.NewResults()
	if err != nil {
		t.Fatal(err)
	}
	caps, err := results.NewCapTable(1)
	if err != nil {
		t.Fatal(err)
	}
	descriptor(caps.At(0))
	content, err := capnp.NewStruct(results.Segment(), capnp.ObjectSize{PointerCount: 1})
	if err != nil {
		t.Fatal(err)
	}
	if err := content.SetPtr(0, capnp.NewInterface(results.Segment(), 0).ToPtr()); err != nil {
		t.Fatal(err)
	}
	if err := results.SetContent(content.ToPtr()); err != nil {
		t.Fatal(err)
	}
	if err := out.Send(); err != nil {
		t.Fatal(err)
	}
}

func sendLifecycleCallWithPromiseParam(
	t *testing.T,
	peer rpc.Transport,
	questionID, exportID, promiseID uint32,
) {
	t.Helper()

	out, err := peer.NewMessage()
	if err != nil {
		t.Fatal(err)
	}
	defer out.Release()
	call, err := out.Message().NewCall()
	if err != nil {
		t.Fatal(err)
	}
	call.SetQuestionId(questionID)
	call.SetInterfaceId(interfaceID)
	call.SetMethodId(methodID)
	target, err := call.NewTarget()
	if err != nil {
		t.Fatal(err)
	}
	target.SetImportedCap(exportID)
	params, err := call.NewParams()
	if err != nil {
		t.Fatal(err)
	}
	caps, err := params.NewCapTable(1)
	if err != nil {
		t.Fatal(err)
	}
	caps.At(0).SetSenderPromise(promiseID)
	content, err := capnp.NewStruct(params.Segment(), capnp.ObjectSize{PointerCount: 1})
	if err != nil {
		t.Fatal(err)
	}
	if err := content.SetPtr(0, capnp.NewInterface(params.Segment(), 0).ToPtr()); err != nil {
		t.Fatal(err)
	}
	if err := params.SetContent(content.ToPtr()); err != nil {
		t.Fatal(err)
	}
	if err := out.Send(); err != nil {
		t.Fatal(err)
	}
}

func recvLifecycleMessage(t *testing.T, peer rpc.Transport) (*rpcMessage, capnp.ReleaseFunc) {
	t.Helper()
	msg, release, err := recvMessage(context.Background(), peer)
	if err != nil {
		t.Fatal(err)
	}
	return msg, release
}

func requireLifecycleMessage(t *testing.T, peer rpc.Transport, want rpccp.Message_Which) *rpcMessage {
	t.Helper()
	msg, release := recvLifecycleMessage(t, peer)
	t.Cleanup(release)
	if msg.Which != want {
		t.Fatalf("message = %v; want %v", msg.Which, want)
	}
	return msg
}

func TestLifecycleReleasedSenderPromiseBalancesLateResolve(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		const (
			promiseID     = 41
			replacementID = 42
			markerID      = 91
		)

		observed, peer := newLifecycleDriver(t, 16)
		conn := rpc.NewConn(observed, &rpc.Options{Logger: testErrorReporter{tb: t}})
		t.Cleanup(func() { _ = conn.Close() })

		promise := conn.Bootstrap(context.Background())
		bootstrap := requireLifecycleMessage(t, peer, rpccp.Message_Which_bootstrap)
		sendLifecycleCapabilityReturn(t, peer, bootstrap.Bootstrap.QuestionID, func(cap rpccp.CapDescriptor) {
			cap.SetSenderPromise(promiseID)
		})

		finish := requireLifecycleMessage(t, peer, rpccp.Message_Which_finish)
		if finish.Finish.QuestionID != bootstrap.Bootstrap.QuestionID {
			t.Fatalf("Finish question = %d; want %d", finish.Finish.QuestionID, bootstrap.Bootstrap.QuestionID)
		}

		promise.Release()
		release := requireLifecycleMessage(t, peer, rpccp.Message_Which_release)
		if release.Release.ID != promiseID || release.Release.ReferenceCount != 1 {
			t.Fatalf("promise Release = %+v; want id=%d count=1", release.Release, promiseID)
		}

		if err := sendMessage(context.Background(), peer, &rpcMessage{
			Which: rpccp.Message_Which_resolve,
			Resolve: &rpcResolve{
				PromiseID: promiseID,
				Which:     rpccp.Resolve_Which_cap,
				Cap: &rpcCapDescriptor{
					Which:        rpccp.CapDescriptor_Which_senderHosted,
					SenderHosted: replacementID,
				},
			},
		}); err != nil {
			t.Fatal(err)
		}
		if err := sendMessage(context.Background(), peer, &rpcMessage{
			Which:     rpccp.Message_Which_bootstrap,
			Bootstrap: &rpcBootstrap{QuestionID: markerID},
		}); err != nil {
			t.Fatal(err)
		}

		replacementRelease := requireLifecycleMessage(t, peer, rpccp.Message_Which_release)
		if replacementRelease.Release.ID != replacementID ||
			replacementRelease.Release.ReferenceCount != 1 {
			t.Fatalf(
				"replacement Release = %+v; want id=%d count=1",
				replacementRelease.Release,
				replacementID,
			)
		}
		marker := requireLifecycleMessage(t, peer, rpccp.Message_Which_return)
		if marker.Return.AnswerID != markerID {
			t.Fatalf("marker Return answer = %d; want %d", marker.Return.AnswerID, markerID)
		}
		select {
		case <-conn.Done():
			t.Fatal("late Resolve closed the connection")
		default:
		}

		oracle := newLifecycleOracle()
		imported := oracle.observeImportedPromise(promiseID)
		if err := oracle.localRelease(imported); err != nil {
			t.Fatal(err)
		}
		if err := oracle.peerResolve(imported); err != nil {
			t.Fatal(err)
		}
		replacement := oracle.ids.capability(sidePeer, replacementID, false)
		ledger := lifecycleLedger{
			{Kind: actionObserveImportPromise, Capability: imported, ReferenceDelta: 1},
			{Kind: actionWireRelease, Capability: imported, ReferenceDelta: -1},
			{Kind: actionPeerResolve, Capability: imported},
			{Kind: actionObserveImportPromise, Capability: replacement, ReferenceDelta: 1},
			{Kind: actionWireRelease, Capability: replacement, ReferenceDelta: -1},
		}
		constraints, err := buildLateResolveConstraints(ruleInput{
			Action:      actionPeerResolve,
			Capability:  imported,
			Replacement: replacement,
		})
		if err != nil {
			t.Fatal(err)
		}
		if err := checkConstraintAlternatives(constraints, ledger); err != nil {
			t.Fatal(err)
		}
		for _, observation := range observed.Observations() {
			if observation.Which == rpccp.Message_Which_disembargo {
				t.Fatalf("late Resolve produced Disembargo: %+v", observation)
			}
		}
	})
}

func TestLifecyclePromisePathShorteningPreservesDeliveryOrder(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		observed, peer := newLifecycleDriver(t, 16)
		conn := rpc.NewConn(observed, &rpc.Options{Logger: testErrorReporter{tb: t}})
		t.Cleanup(func() { _ = conn.Close() })

		bootstrapClient := conn.Bootstrap(context.Background())
		defer bootstrapClient.Release()
		bootstrap := requireLifecycleMessage(t, peer, rpccp.Message_Which_bootstrap)
		sendLifecycleCapabilityReturn(
			t,
			peer,
			bootstrap.Bootstrap.QuestionID,
			func(cap rpccp.CapDescriptor) { cap.SetSenderHosted(bootstrapExportID) },
		)
		requireLifecycleMessage(t, peer, rpccp.Message_Which_finish)

		deliveries := make(chan uint64, 4)
		local := newLifecycleServer(func(_ context.Context, call *server.Call) error {
			n := call.Args().Uint64(0)
			deliveries <- n
			results, err := call.AllocResults(capnp.ObjectSize{DataSize: 8})
			if err != nil {
				return err
			}
			results.SetUint64(0, n)
			return nil
		}, nil)
		defer local.Release()

		answerA, releaseA := bootstrapClient.SendCall(context.Background(), capnp.Send{
			Method:   capnp.Method{InterfaceID: interfaceID, MethodID: methodID},
			ArgsSize: capnp.ObjectSize{PointerCount: 1},
			PlaceArgs: func(args capnp.Struct) error {
				id := args.Message().CapTable().Add(local)
				return args.SetPtr(0, capnp.NewInterface(args.Segment(), id).ToPtr())
			},
		})
		defer releaseA()
		callA, releaseCallAMessage := recvLifecycleMessage(t, peer)
		if callA.Which != rpccp.Message_Which_call {
			releaseCallAMessage()
			t.Fatalf("message = %v; want call A", callA.Which)
		}
		questionA := callA.Call.QuestionID
		param, err := callA.Call.Params.Content.Struct().Ptr(0)
		if err != nil {
			releaseCallAMessage()
			t.Fatal(err)
		}
		capIndex := param.Interface().Capability()
		if int64(capIndex) >= int64(len(callA.Call.Params.CapTable)) {
			releaseCallAMessage()
			t.Fatalf("call A parameter capability = %d; table has %d entries", capIndex, len(callA.Call.Params.CapTable))
		}
		export := callA.Call.Params.CapTable[capIndex]
		if export.Which != rpccp.CapDescriptor_Which_senderHosted {
			releaseCallAMessage()
			t.Fatalf("call A parameter = %v; want senderHosted", export.Which)
		}
		releaseCallAMessage()

		sendPipelineCall := func(n uint64) (*capnp.Answer, capnp.ReleaseFunc) {
			return answerA.PipelineSend(
				context.Background(),
				[]capnp.PipelineOp{{Field: 0}},
				capnp.Send{
					Method:   capnp.Method{InterfaceID: interfaceID, MethodID: methodID},
					ArgsSize: capnp.ObjectSize{DataSize: 8},
					PlaceArgs: func(args capnp.Struct) error {
						args.SetUint64(0, n)
						return nil
					},
				},
			)
		}

		answerB, releaseB := sendPipelineCall(1)
		defer releaseB()
		callB, releaseCallBMessage := recvLifecycleMessage(t, peer)
		if callB.Which != rpccp.Message_Which_call {
			releaseCallBMessage()
			t.Fatalf("message = %v; want call B", callB.Which)
		}
		if callB.Call.Target.Which != rpccp.MessageTarget_Which_promisedAnswer ||
			callB.Call.Target.PromisedAnswer.QuestionID != questionA {
			releaseCallBMessage()
			t.Fatalf("call B target = %+v; want promised answer %d", callB.Call.Target, questionA)
		}
		questionB := callB.Call.QuestionID
		releaseCallBMessage()
		// Releasing the peer-side incoming message acknowledges the pipelined
		// send.  Let its callback drain before resolving call A; Disembargo
		// creation is deliberately ordered after admitted pipeline sends.
		synctest.Wait()

		sendLifecycleCapabilityFieldReturn(
			t,
			peer,
			questionA,
			false, // the peer reflects this parameter capability in a later Call
			func(cap rpccp.CapDescriptor) { cap.SetReceiverHosted(export.SenderHosted) },
		)
		var (
			disembargo *rpcMessage
			sawFinish  bool
		)
		for !sawFinish || disembargo == nil {
			message, release := recvLifecycleMessage(t, peer)
			defer release()
			switch message.Which {
			case rpccp.Message_Which_disembargo:
				disembargo = message
			case rpccp.Message_Which_finish:
				sawFinish = true
				if message.Finish.QuestionID != questionA {
					t.Fatalf("Finish question = %d; want %d", message.Finish.QuestionID, questionA)
				}
			default:
				t.Fatalf("message after call A Return = %v; want Finish or Disembargo", message.Which)
			}
		}
		if disembargo.Disembargo.Context.Which != rpccp.Disembargo_context_Which_senderLoopback {
			t.Fatalf(
				"outgoing Disembargo context = %v; want senderLoopback",
				disembargo.Disembargo.Context.Which,
			)
		}
		embargoID := disembargo.Disembargo.Context.SenderLoopback

		answerCReady := make(chan struct{})
		var (
			answerC  *capnp.Answer
			releaseC capnp.ReleaseFunc
		)
		go func() {
			answerC, releaseC = sendPipelineCall(2)
			close(answerCReady)
		}()
		defer func() {
			select {
			case <-answerCReady:
				releaseC()
			default:
			}
		}()

		const reflectedQuestion = 909
		_, argsSegment := capnp.NewSingleSegmentMessage(nil)
		reflectedArgs, err := capnp.NewStruct(argsSegment, capnp.ObjectSize{DataSize: 8})
		if err != nil {
			t.Fatal(err)
		}
		reflectedArgs.SetUint64(0, 1)
		if err := sendMessage(context.Background(), peer, &rpcMessage{
			Which: rpccp.Message_Which_call,
			Call: &rpcCall{
				QuestionID:  reflectedQuestion,
				InterfaceID: interfaceID,
				MethodID:    methodID,
				Target: rpcMessageTarget{
					Which:       rpccp.MessageTarget_Which_importedCap,
					ImportedCap: export.SenderHosted,
				},
				Params: rpcPayload{Content: reflectedArgs.ToPtr()},
				SendResultsTo: rpcCallSendResultsTo{
					Which: rpccp.Call_sendResultsTo_Which_caller,
				},
			},
		}); err != nil {
			t.Fatal(err)
		}
		synctest.Wait()
		select {
		case got := <-deliveries:
			if got != 1 {
				t.Fatalf("pre-resolution delivery = %d; want 1", got)
			}
		default:
			t.Fatal("pre-resolution call was not delivered")
		}
		select {
		case got := <-deliveries:
			t.Fatalf("post-resolution call %d overtook receiver-loopback Disembargo", got)
		default:
		}

		if err := sendMessage(context.Background(), peer, &rpcMessage{
			Which: rpccp.Message_Which_disembargo,
			Disembargo: &rpcDisembargo{
				Context: rpcDisembargoContext{
					Which:            rpccp.Disembargo_context_Which_receiverLoopback,
					ReceiverLoopback: embargoID,
				},
				Target: rpcMessageTarget{
					Which:       rpccp.MessageTarget_Which_importedCap,
					ImportedCap: export.SenderHosted,
				},
			},
		}); err != nil {
			t.Fatal(err)
		}
		synctest.Wait()
		select {
		case got := <-deliveries:
			if got != 2 {
				t.Fatalf("post-resolution delivery = %d; want 2", got)
			}
		default:
			t.Fatal("post-resolution call remained embargoed after receiver loopback")
		}

		reflectedReturn := requireLifecycleMessage(t, peer, rpccp.Message_Which_return)
		if reflectedReturn.Return.AnswerID != reflectedQuestion {
			t.Fatalf(
				"reflected Return answer = %d; want %d",
				reflectedReturn.Return.AnswerID,
				reflectedQuestion,
			)
		}
		if err := sendMessage(context.Background(), peer, &rpcMessage{
			Which: rpccp.Message_Which_return,
			Return: &rpcReturn{
				AnswerID: questionB,
				Which:    rpccp.Return_Which_results,
				Results:  reflectedReturn.Return.Results,
			},
		}); err != nil {
			t.Fatal(err)
		}
		if _, err := answerB.Struct(); err != nil {
			t.Fatal("call B:", err)
		}
		<-answerCReady
		if _, err := answerC.Struct(); err != nil {
			t.Fatal("call C:", err)
		}

		ledger := lifecycleLedger{
			{Kind: actionObserveDelivery, Label: "pre-resolution"},
			{Kind: actionWireDisembargo},
			{Kind: actionPeerDisembargo},
			{Kind: actionObserveDelivery, Label: "post-resolution"},
		}
		constraints, err := buildDisembargoConstraints(probeFor(actionObserveDelivery))
		if err != nil {
			t.Fatal(err)
		}
		if err := checkConstraintAlternatives(constraints[1:], ledger); err != nil {
			t.Fatal(err)
		}

		var senderLoopback, receiverLoopback int
		for _, observation := range observed.Observations() {
			if observation.Which != rpccp.Message_Which_disembargo {
				continue
			}
			switch observation.DisembargoContext {
			case rpccp.Disembargo_context_Which_senderLoopback:
				senderLoopback++
			case rpccp.Disembargo_context_Which_receiverLoopback:
				receiverLoopback++
			default:
				t.Fatal(fmt.Sprintf("unexpected Disembargo observation: %+v", observation))
			}
		}
		if senderLoopback != 1 || receiverLoopback != 1 {
			t.Fatalf(
				"Disembargo observations = senderLoopback:%d receiverLoopback:%d; want 1 each",
				senderLoopback,
				receiverLoopback,
			)
		}
	})
}

func TestLifecycleDisconnectDrainsQuestionAndPromiseBeforeDone(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		const (
			bootstrapQuestion = 100
			importQuestion    = 101
			promiseID         = 51
		)

		observed, peer := newLifecycleDriver(t, 16)
		imported := make(chan capnp.Client, 1)
		shutdownEntered := make(chan struct{})
		bootstrapLocal := newLifecycleServer(
			func(_ context.Context, call *server.Call) error {
				param, err := call.Args().Ptr(0)
				if err != nil {
					return err
				}
				imported <- param.Interface().Client().AddRef()
				_, err = call.AllocResults(capnp.ObjectSize{})
				return err
			},
			func() { close(shutdownEntered) },
		)
		conn := rpc.NewConn(observed, &rpc.Options{
			BootstrapClient: bootstrapLocal.AddRef(),
			Logger:          testErrorReporter{tb: t},
		})
		t.Cleanup(func() { _ = conn.Close() })

		if err := sendMessage(context.Background(), peer, &rpcMessage{
			Which:     rpccp.Message_Which_bootstrap,
			Bootstrap: &rpcBootstrap{QuestionID: bootstrapQuestion},
		}); err != nil {
			t.Fatal(err)
		}
		bootstrapReturn, releaseBootstrapReturn := recvLifecycleMessage(t, peer)
		if bootstrapReturn.Which != rpccp.Message_Which_return ||
			bootstrapReturn.Return.AnswerID != bootstrapQuestion ||
			bootstrapReturn.Return.Which != rpccp.Return_Which_results ||
			len(bootstrapReturn.Return.Results.CapTable) != 1 ||
			bootstrapReturn.Return.Results.CapTable[0].Which != rpccp.CapDescriptor_Which_senderHosted {
			releaseBootstrapReturn()
			t.Fatalf("bootstrap Return = %+v; want one senderHosted capability", bootstrapReturn)
		}
		exportID := bootstrapReturn.Return.Results.CapTable[0].SenderHosted
		releaseBootstrapReturn()
		if err := sendMessage(context.Background(), peer, &rpcMessage{
			Which: rpccp.Message_Which_finish,
			Finish: &rpcFinish{
				QuestionID: bootstrapQuestion,
			},
		}); err != nil {
			t.Fatal(err)
		}

		sendLifecycleCallWithPromiseParam(t, peer, importQuestion, exportID, promiseID)
		promise := <-imported
		defer promise.Release()
		importReturn, releaseImportReturn := recvLifecycleMessage(t, peer)
		if importReturn.Which != rpccp.Message_Which_return ||
			importReturn.Return.AnswerID != importQuestion {
			releaseImportReturn()
			t.Fatalf("promise import call Return = %+v; want answer %d", importReturn, importQuestion)
		}
		releaseImportReturn()
		if err := sendMessage(context.Background(), peer, &rpcMessage{
			Which: rpccp.Message_Which_finish,
			Finish: &rpcFinish{
				QuestionID: importQuestion,
			},
		}); err != nil {
			t.Fatal(err)
		}

		pending, releasePending := promise.SendCall(context.Background(), capnp.Send{
			Method: capnp.Method{InterfaceID: interfaceID, MethodID: methodID},
		})
		defer releasePending()

		call, releaseCall := recvLifecycleMessage(t, peer)
		if call.Which != rpccp.Message_Which_call {
			releaseCall()
			t.Fatalf("message = %v; want pending call", call.Which)
		}
		if call.Call.Target.Which != rpccp.MessageTarget_Which_importedCap ||
			call.Call.Target.ImportedCap != promiseID {
			releaseCall()
			t.Fatalf("pending call target = %+v; want imported capability %d", call.Call.Target, promiseID)
		}
		releaseCall()
		synctest.Wait()

		importSnapshot := promise.Snapshot()
		defer importSnapshot.Release()
		if !importSnapshot.IsPromise() || importSnapshot.IsResolved() {
			t.Fatal("senderPromise is not an unresolved promise before disconnect")
		}
		resolveSnapshot := promise.Snapshot()
		defer resolveSnapshot.Release()
		resolveResult := make(chan error, 1)
		go func() {
			resolveResult <- resolveSnapshot.Resolve(context.Background())
		}()
		synctest.Wait()

		if err := peer.Close(); err != nil {
			t.Fatal(err)
		}
		<-conn.Done()

		// Conn.Done is the observation boundary: both the pending question and
		// the unresolved imported capability must already be terminal here.
		select {
		case <-pending.Done():
		default:
			t.Fatal("pending question remained open when Conn.Done closed")
		}
		if _, err := pending.Struct(); !capnp.IsDisconnected(err) {
			t.Fatalf("pending question error = %v; want disconnected", err)
		}
		if importSnapshot.IsPromise() && !importSnapshot.IsResolved() {
			t.Fatal("unresolved imported promise remained open when Conn.Done closed")
		}
		synctest.Wait()
		if err := <-resolveResult; err != nil {
			t.Fatalf("imported promise Resolve = %v; want resolution to an error client", err)
		}
		resolutionErr, ok := resolveSnapshot.Brand().Value.(error)
		if !ok || !capnp.IsDisconnected(resolutionErr) {
			t.Fatalf("imported promise brand = %v; want disconnected error", resolutionErr)
		}

		var releasedImportCall bool
		for _, observation := range observed.Observations() {
			if observation.Flow == lifecyclePeerToGo &&
				observation.Which == rpccp.Message_Which_call &&
				observation.QuestionID == importQuestion {
				releasedImportCall = observation.IncomingReleaseCalls() == 1
			}
		}
		if !releasedImportCall {
			t.Fatal("promise-bearing Call was not released exactly once before disconnect completed")
		}
		bootstrapLocal.Release()
		synctest.Wait()
		select {
		case <-shutdownEntered:
		default:
			t.Fatal("exported capability was not eventually released during disconnect")
		}

		ledger := lifecycleLedger{
			{Kind: actionObserveCompletion, Label: "pending-question"},
			{Kind: actionObserveConnDone},
			{Kind: actionObserveCapabilityRelease, Label: "export"},
		}
		schemaConstraints, err := buildDisconnectConstraints(probeFor(actionPeerAbort))
		if err != nil {
			t.Fatal(err)
		}
		if err := checkConstraintAlternatives(schemaConstraints, ledger); err != nil {
			t.Fatal(err)
		}
		if err := checkConstraintAlternatives(goConnectionDrainConstraints(), ledger); err != nil {
			t.Fatal(err)
		}
	})
}

type lifecycleShutdownFunc func()

func (f lifecycleShutdownFunc) Shutdown() {
	if f != nil {
		f()
	}
}

func newLifecycleServer(
	impl func(context.Context, *server.Call) error,
	shutdown lifecycleShutdownFunc,
) capnp.Client {
	return capnp.NewClient(server.New([]server.Method{{
		Method: capnp.Method{InterfaceID: interfaceID, MethodID: methodID},
		Impl:   impl,
	}}, nil, shutdown))
}
