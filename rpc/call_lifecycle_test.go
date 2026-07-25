package rpc

import (
	"context"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/rpc/transport"
	rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"
)

type callReleaseTransport struct {
	outgoing *callReleaseOutgoingMessage
}

func (t *callReleaseTransport) NewMessage() (transport.OutgoingMessage, error) {
	_, seg := capnp.NewSingleSegmentMessage(nil)
	message, err := rpccp.NewRootMessage(seg)
	if err != nil {
		return nil, err
	}
	t.outgoing = &callReleaseOutgoingMessage{message: message}
	return t.outgoing, nil
}

func (*callReleaseTransport) RecvMessage() (transport.IncomingMessage, error) {
	return nil, errors.New("unused")
}

func (*callReleaseTransport) Close() error { return nil }

type callReleaseOutgoingMessage struct {
	message  rpccp.Message
	releases atomic.Int32
	once     sync.Once
}

func (m *callReleaseOutgoingMessage) Message() rpccp.Message { return m.message }
func (*callReleaseOutgoingMessage) Send() error {
	panic("rejected Call sent its Return")
}
func (m *callReleaseOutgoingMessage) Release() {
	m.releases.Add(1)
	m.once.Do(m.message.Message().Release)
}

func newRejectedTargetCall(
	t *testing.T,
	id answerID,
	target func(rpccp.MessageTarget) error,
) *countingIncomingMessage {
	t.Helper()

	_, seg := capnp.NewSingleSegmentMessage(nil)
	message, err := rpccp.NewRootMessage(seg)
	if err != nil {
		t.Fatal(err)
	}
	call, err := message.NewCall()
	if err != nil {
		t.Fatal(err)
	}
	call.SetQuestionId(uint32(id))
	callTarget, err := call.NewTarget()
	if err != nil {
		t.Fatal(err)
	}
	if err := target(callTarget); err != nil {
		t.Fatal(err)
	}
	return &countingIncomingMessage{message: message}
}

func TestRejectedCallTargetReleasesUnsentReturn(t *testing.T) {
	const (
		callID   answerID = 17
		targetID answerID = 23
	)

	tests := []struct {
		name    string
		prepare func(*Conn)
		target  func(rpccp.MessageTarget) error
		wantErr string
	}{
		{
			name: "unknown imported capability",
			target: func(target rpccp.MessageTarget) error {
				target.SetImportedCap(31)
				return nil
			},
			wantErr: "unknown export ID",
		},
		{
			name: "unknown promised answer",
			target: func(target rpccp.MessageTarget) error {
				promised, err := target.NewPromisedAnswer()
				if err == nil {
					promised.SetQuestionId(uint32(targetID))
				}
				return err
			},
			wantErr: "use of unknown or finished answer ID",
		},
		{
			name: "finished promised answer",
			prepare: func(conn *Conn) {
				conn.lk.answers.Create(targetID, &ansent{flags: finishReceived})
			},
			target: func(target rpccp.MessageTarget) error {
				promised, err := target.NewPromisedAnswer()
				if err == nil {
					promised.SetQuestionId(uint32(targetID))
				}
				return err
			},
			wantErr: "use of unknown or finished answer ID",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			trans := new(callReleaseTransport)
			conn := &Conn{
				transport: trans,
				bgctx:     context.Background(),
			}
			if test.prepare != nil {
				test.prepare(conn)
			}
			in := newRejectedTargetCall(t, callID, test.target)

			err := conn.handleCall(context.Background(), in)
			if err == nil || !strings.Contains(err.Error(), test.wantErr) {
				t.Fatalf("handleCall() error = %v; want error containing %q", err, test.wantErr)
			}
			if got := atomic.LoadInt32(&in.releases); got != 1 {
				t.Errorf("incoming Call releases = %d; want 1", got)
			}
			if trans.outgoing == nil {
				t.Fatal("Return message was not allocated")
			}
			if got := trans.outgoing.releases.Load(); got != 1 {
				t.Errorf("outgoing Return releases = %d; want 1", got)
			}
		})
	}
}
