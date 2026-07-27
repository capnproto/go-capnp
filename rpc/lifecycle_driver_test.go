package rpc_test

import (
	"sync"
	"testing"

	"capnproto.org/go/capnp/v3/rpc"
	"capnproto.org/go/capnp/v3/rpc/transport"
	rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"
)

// lifecycleWireFlow names a message direction relative to the Go connection
// under test.
type lifecycleWireFlow uint8

const (
	lifecyclePeerToGo lifecycleWireFlow = iota
	lifecycleGoToPeer
)

func (f lifecycleWireFlow) String() string {
	switch f {
	case lifecyclePeerToGo:
		return "peer->go"
	case lifecycleGoToPeer:
		return "go->peer"
	default:
		return "unknown"
	}
}

// lifecycleWireObservation is the scalar protocol state needed by lifecycle
// tests.  It intentionally does not retain pointers into the message arena.
type lifecycleWireObservation struct {
	Sequence   uint64
	Flow       lifecycleWireFlow
	Which      rpccp.Message_Which
	SummaryErr string

	QuestionID uint32
	AnswerID   uint32
	PromiseID  uint32
	ImportID   uint32

	NoFinishNeeded    bool
	ReleaseResultCaps bool
	ReferenceCount    uint32
	ResultCapCount    int

	ResolveWhich      rpccp.Resolve_Which
	ResolveCapWhich   rpccp.CapDescriptor_Which
	ResolveCapID      uint32
	TargetWhich       rpccp.MessageTarget_Which
	PromisedQuestion  uint32
	DisembargoContext rpccp.Disembargo_context_Which
	DisembargoID      uint32

	incomingReleaseCalls int
}

func (o lifecycleWireObservation) IncomingReleaseCalls() int {
	return o.incomingReleaseCalls
}

// lifecycleObservedTransport is a protocol-only decorator.  A connection can
// use it anywhere it accepts rpc.Transport; tests retain the decorator to
// inspect a stable copy of messages and exact IncomingMessage.Release calls.
//
// Sequence is the order in which the decorator observed transport operations.
// It is useful for diagnostics, but tests should use protocol acknowledgements
// (and synctest.Wait where appropriate) to establish causal ordering.
type lifecycleObservedTransport struct {
	rpc.Transport

	mu           sync.Mutex
	nextSequence uint64
	observations []lifecycleWireObservation
}

func observeLifecycleTransport(inner rpc.Transport) *lifecycleObservedTransport {
	return &lifecycleObservedTransport{Transport: inner}
}

func (t *lifecycleObservedTransport) NewMessage() (transport.OutgoingMessage, error) {
	out, err := t.Transport.NewMessage()
	if err != nil {
		return nil, err
	}
	return &lifecycleObservedOutgoing{OutgoingMessage: out, observer: t}, nil
}

func (t *lifecycleObservedTransport) RecvMessage() (transport.IncomingMessage, error) {
	in, err := t.Transport.RecvMessage()
	if err != nil {
		return nil, err
	}
	index := t.record(lifecyclePeerToGo, in.Message())
	return &lifecycleObservedIncoming{
		IncomingMessage: in,
		observer:        t,
		observation:     index,
	}, nil
}

func (t *lifecycleObservedTransport) record(flow lifecycleWireFlow, msg rpccp.Message) int {
	return t.recordObservation(summarizeLifecycleWireMessage(flow, msg))
}

func (t *lifecycleObservedTransport) recordObservation(observation lifecycleWireObservation) int {
	t.mu.Lock()
	observation.Sequence = t.nextSequence
	t.nextSequence++
	t.observations = append(t.observations, observation)
	index := len(t.observations) - 1
	t.mu.Unlock()
	return index
}

func (t *lifecycleObservedTransport) recordIncomingRelease(index int) {
	t.mu.Lock()
	t.observations[index].incomingReleaseCalls++
	t.mu.Unlock()
}

func (t *lifecycleObservedTransport) Observations() []lifecycleWireObservation {
	t.mu.Lock()
	defer t.mu.Unlock()
	return append([]lifecycleWireObservation(nil), t.observations...)
}

type lifecycleObservedIncoming struct {
	transport.IncomingMessage
	observer    *lifecycleObservedTransport
	observation int
}

func (in *lifecycleObservedIncoming) Release() {
	in.observer.recordIncomingRelease(in.observation)
	in.IncomingMessage.Release()
}

type lifecycleObservedOutgoing struct {
	transport.OutgoingMessage
	observer *lifecycleObservedTransport
}

func (out *lifecycleObservedOutgoing) Send() error {
	observation := summarizeLifecycleWireMessage(lifecycleGoToPeer, out.Message())
	if err := out.OutgoingMessage.Send(); err != nil {
		return err
	}
	out.observer.recordObservation(observation)
	return nil
}

func summarizeLifecycleWireMessage(flow lifecycleWireFlow, msg rpccp.Message) lifecycleWireObservation {
	o := lifecycleWireObservation{
		Flow:  flow,
		Which: msg.Which(),
	}
	switch o.Which {
	case rpccp.Message_Which_bootstrap:
		m, err := msg.Bootstrap()
		if err != nil {
			o.SummaryErr = err.Error()
			break
		}
		o.QuestionID = m.QuestionId()
	case rpccp.Message_Which_call:
		m, err := msg.Call()
		if err != nil {
			o.SummaryErr = err.Error()
			break
		}
		o.QuestionID = m.QuestionId()
		target, err := m.Target()
		if err != nil {
			o.SummaryErr = err.Error()
			break
		}
		summarizeLifecycleTarget(&o, target)
	case rpccp.Message_Which_return:
		m, err := msg.Return()
		if err != nil {
			o.SummaryErr = err.Error()
			break
		}
		o.AnswerID = m.AnswerId()
		o.NoFinishNeeded = m.NoFinishNeeded()
		if m.Which() == rpccp.Return_Which_results {
			results, err := m.Results()
			if err != nil {
				o.SummaryErr = err.Error()
				break
			}
			caps, err := results.CapTable()
			if err != nil {
				o.SummaryErr = err.Error()
				break
			}
			o.ResultCapCount = caps.Len()
		}
	case rpccp.Message_Which_finish:
		m, err := msg.Finish()
		if err != nil {
			o.SummaryErr = err.Error()
			break
		}
		o.QuestionID = m.QuestionId()
		o.ReleaseResultCaps = m.ReleaseResultCaps()
	case rpccp.Message_Which_resolve:
		m, err := msg.Resolve()
		if err != nil {
			o.SummaryErr = err.Error()
			break
		}
		o.PromiseID = m.PromiseId()
		o.ResolveWhich = m.Which()
		if o.ResolveWhich == rpccp.Resolve_Which_cap {
			cap, err := m.Cap()
			if err != nil {
				o.SummaryErr = err.Error()
				break
			}
			o.ResolveCapWhich = cap.Which()
			o.ResolveCapID = lifecycleCapDescriptorID(cap)
		}
	case rpccp.Message_Which_release:
		m, err := msg.Release()
		if err != nil {
			o.SummaryErr = err.Error()
			break
		}
		o.ImportID = m.Id()
		o.ReferenceCount = m.ReferenceCount()
	case rpccp.Message_Which_disembargo:
		m, err := msg.Disembargo()
		if err != nil {
			o.SummaryErr = err.Error()
			break
		}
		target, err := m.Target()
		if err != nil {
			o.SummaryErr = err.Error()
			break
		}
		summarizeLifecycleTarget(&o, target)
		context := m.Context()
		o.DisembargoContext = context.Which()
		switch o.DisembargoContext {
		case rpccp.Disembargo_context_Which_senderLoopback:
			o.DisembargoID = context.SenderLoopback()
		case rpccp.Disembargo_context_Which_receiverLoopback:
			o.DisembargoID = context.ReceiverLoopback()
		}
	}
	return o
}

func summarizeLifecycleTarget(o *lifecycleWireObservation, target rpccp.MessageTarget) {
	o.TargetWhich = target.Which()
	switch o.TargetWhich {
	case rpccp.MessageTarget_Which_importedCap:
		o.ImportID = target.ImportedCap()
	case rpccp.MessageTarget_Which_promisedAnswer:
		answer, err := target.PromisedAnswer()
		if err != nil {
			o.SummaryErr = err.Error()
			return
		}
		o.PromisedQuestion = answer.QuestionId()
	}
}

func lifecycleCapDescriptorID(cap rpccp.CapDescriptor) uint32 {
	switch cap.Which() {
	case rpccp.CapDescriptor_Which_senderHosted:
		return cap.SenderHosted()
	case rpccp.CapDescriptor_Which_senderPromise:
		return cap.SenderPromise()
	case rpccp.CapDescriptor_Which_receiverHosted:
		return cap.ReceiverHosted()
	case rpccp.CapDescriptor_Which_thirdPartyHosted:
		thirdParty, err := cap.ThirdPartyHosted()
		if err == nil {
			return thirdParty.VineId()
		}
	}
	return 0
}

// newLifecycleDriver gives a Conn-side observed transport and a scripted peer.
// The peer uses the existing sendMessage/recvMessage wire helpers.
func newLifecycleDriver(t testing.TB, buffer int) (*lifecycleObservedTransport, rpc.Transport) {
	t.Helper()
	goSide, peerSide := transport.NewPipe(buffer)
	observed := observeLifecycleTransport(rpc.NewTransport(goSide))
	peer := rpc.NewTransport(peerSide)
	t.Cleanup(func() {
		_ = observed.Close()
		_ = peer.Close()
	})
	return observed, peer
}

func TestLifecycleObservedTransportRecordsProtocolAndIncomingReleases(t *testing.T) {
	observed, peer := newLifecycleDriver(t, 4)

	if err := sendMessage(t.Context(), peer, &rpcMessage{
		Which: rpccp.Message_Which_return,
		Return: &rpcReturn{
			AnswerID:       17,
			NoFinishNeeded: true,
			Which:          rpccp.Return_Which_results,
			Results:        &rpcPayload{},
		},
	}); err != nil {
		t.Fatal(err)
	}
	in, err := observed.RecvMessage()
	if err != nil {
		t.Fatal(err)
	}
	in.Release()
	in.Release()

	out, err := observed.NewMessage()
	if err != nil {
		t.Fatal(err)
	}
	finish, err := out.Message().NewFinish()
	if err != nil {
		t.Fatal(err)
	}
	finish.SetQuestionId(17)
	finish.SetReleaseResultCaps(true)
	if err := out.Send(); err != nil {
		t.Fatal(err)
	}
	out.Release()

	peerMessage, err := peer.RecvMessage()
	if err != nil {
		t.Fatal(err)
	}
	peerMessage.Release()

	got := observed.Observations()
	if len(got) != 2 {
		t.Fatalf("observations = %d; want 2", len(got))
	}
	if got[0].Flow != lifecyclePeerToGo ||
		got[0].Which != rpccp.Message_Which_return ||
		got[0].AnswerID != 17 ||
		!got[0].NoFinishNeeded ||
		got[0].IncomingReleaseCalls() != 2 ||
		got[0].SummaryErr != "" {
		t.Fatalf("incoming observation = %+v", got[0])
	}
	if got[1].Flow != lifecycleGoToPeer ||
		got[1].Which != rpccp.Message_Which_finish ||
		got[1].QuestionID != 17 ||
		!got[1].ReleaseResultCaps ||
		got[1].SummaryErr != "" {
		t.Fatalf("outgoing observation = %+v", got[1])
	}
}
