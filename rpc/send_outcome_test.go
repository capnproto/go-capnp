package rpc

import (
	"context"
	"errors"
	"strings"
	"sync"
	"testing"
	"time"

	capnp "capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/exp/spsc"
	"capnproto.org/go/capnp/v3/rpc/transport"
	rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"
)

type newMessageErrorTransport struct{ err error }

func (t newMessageErrorTransport) NewMessage() (transport.OutgoingMessage, error) {
	return nil, t.err
}

func (newMessageErrorTransport) RecvMessage() (transport.IncomingMessage, error) {
	return nil, errors.New("unused")
}

func (newMessageErrorTransport) Close() error { return nil }

func TestSendMessageNewMessageError(t *testing.T) {
	wantErr := errors.New("allocate")
	queue := spsc.New[asyncSend]()
	c := &Conn{transport: newMessageErrorTransport{err: wantErr}}
	c.lk.sendTx = &queue.Tx

	var got sendOutcome
	(*lockedConn)(c).sendMessageOutcome(
		context.Background(),
		func(rpccp.Message) error {
			t.Fatal("build called after NewMessage failed")
			return nil
		},
		false,
		func(outcome sendOutcome) { got = outcome },
	)

	pending, ok := queue.Rx.TryRecv()
	if !ok {
		t.Fatal("message was not queued")
	}
	if err := pending.Send(); err != nil {
		t.Fatalf("Send() error = %v; want non-fatal allocation failure", err)
	}
	if got.disposition != sendDefinitelyUnsent || got.fatal {
		t.Fatalf("outcome = %+v; want definitely-unsent non-fatal", got)
	}
	if !errors.Is(got.err, wantErr) {
		t.Fatalf("outcome error = %v; want wrapped %v", got.err, wantErr)
	}
}

func TestAsyncSendOutcomes(t *testing.T) {
	canceled, cancel := context.WithCancel(context.Background())
	cancel()

	tests := []struct {
		name      string
		as        asyncSend
		want      sendDisposition
		wantFatal bool
		abort     bool
	}{
		{
			name: "success",
			as:   asyncSend{ctx: context.Background(), send: func() error { return nil }},
			want: sendSucceeded,
		},
		{
			name: "definitely unsent",
			as:   asyncSend{ctx: context.Background(), preSendErr: errors.New("build")},
			want: sendDefinitelyUnsent,
		},
		{
			name:      "required message definitely unsent",
			as:        asyncSend{ctx: context.Background(), preSendErr: errors.New("build"), fatalIfUnsent: true},
			want:      sendDefinitelyUnsent,
			wantFatal: true,
		},
		{
			name: "canceled before send",
			as:   asyncSend{ctx: canceled, send: func() error { t.Fatal("send called with canceled context"); return nil }},
			want: sendDefinitelyUnsent,
		},
		{
			name:      "required message canceled before send",
			as:        asyncSend{ctx: canceled, fatalIfUnsent: true, send: func() error { t.Fatal("send called with canceled context"); return nil }},
			want:      sendDefinitelyUnsent,
			wantFatal: true,
		},
		{
			name:      "delivery ambiguous",
			as:        asyncSend{ctx: context.Background(), send: func() error { return errors.New("write") }},
			want:      sendDeliveryAmbiguous,
			wantFatal: true,
		},
		{
			name:  "shutdown abort",
			as:    asyncSend{ctx: context.Background()},
			want:  sendAbortedByShutdown,
			abort: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			releases := 0
			test.as.release = func() { releases++ }
			var got sendOutcome
			test.as.onSent = func(outcome sendOutcome) { got = outcome }

			var err error
			if test.abort {
				test.as.Abort(errors.New("closed"))
			} else {
				err = test.as.Send()
			}
			if got.disposition != test.want || got.fatal != test.wantFatal {
				t.Fatalf("outcome = %+v; want disposition %v, fatal %t", got, test.want, test.wantFatal)
			}
			if (err != nil) != test.wantFatal {
				t.Fatalf("Send() error = %v; want fatal %t", err, test.wantFatal)
			}
			if releases != 1 {
				t.Fatalf("release count = %d; want 1", releases)
			}
		})
	}
}

type failingSendTransport struct {
	mu           sync.Mutex
	firstSend    chan struct{}
	closed       chan struct{}
	closeOnce    sync.Once
	newMessages  int
	closeCount   int
	writeErr     error
	firstRelease chan struct{}
}

func newFailingSendTransport(writeErr error) *failingSendTransport {
	return &failingSendTransport{
		firstSend:    make(chan struct{}),
		closed:       make(chan struct{}),
		writeErr:     writeErr,
		firstRelease: make(chan struct{}),
	}
}

func (t *failingSendTransport) NewMessage() (transport.OutgoingMessage, error) {
	_, seg := capnp.NewMultiSegmentMessage(nil)
	m, err := rpccp.NewRootMessage(seg)
	if err != nil {
		return nil, err
	}

	t.mu.Lock()
	t.newMessages++
	messageNumber := t.newMessages
	t.mu.Unlock()

	return &failingOutgoingMessage{
		message: m,
		send: func() error {
			if messageNumber == 1 {
				<-t.firstSend
				return t.writeErr
			}
			return nil
		},
		onRelease: func() {
			if messageNumber == 1 {
				close(t.firstRelease)
			}
		},
	}, nil
}

func (t *failingSendTransport) RecvMessage() (transport.IncomingMessage, error) {
	<-t.closed
	return nil, errors.New("transport closed")
}

func (t *failingSendTransport) Close() error {
	t.closeOnce.Do(func() {
		t.mu.Lock()
		t.closeCount++
		t.mu.Unlock()
		close(t.closed)
	})
	return nil
}

type failingOutgoingMessage struct {
	message   rpccp.Message
	send      func() error
	onRelease func()
	once      sync.Once
}

func (m *failingOutgoingMessage) Message() rpccp.Message { return m.message }
func (m *failingOutgoingMessage) Send() error            { return m.send() }
func (m *failingOutgoingMessage) Release() {
	m.once.Do(func() {
		m.message.Message().Release()
		if m.onRelease != nil {
			m.onRelease()
		}
	})
}

type recordingLogger struct {
	mu     sync.Mutex
	errors []string
}

func (*recordingLogger) Debug(string, ...any) {}
func (*recordingLogger) Info(string, ...any)  {}
func (*recordingLogger) Warn(string, ...any)  {}
func (l *recordingLogger) Error(message string, _ ...any) {
	l.mu.Lock()
	l.errors = append(l.errors, message)
	l.mu.Unlock()
}

func TestSendFailureShutsDownAndDrainsQueue(t *testing.T) {
	writeErr := errors.New("write failed")
	tr := newFailingSendTransport(writeErr)
	logger := new(recordingLogger)
	c := NewConn(tr, &Options{Logger: logger})

	var mu sync.Mutex
	var outcomes [2][]sendOutcome
	c.withLocked(func(c *lockedConn) {
		for i := range outcomes {
			i := i
			c.sendMessageOutcome(context.Background(), func(m rpccp.Message) error {
				_, err := m.NewUnimplemented()
				return err
			}, false, func(outcome sendOutcome) {
				mu.Lock()
				outcomes[i] = append(outcomes[i], outcome)
				mu.Unlock()
			})
		}
	})
	close(tr.firstSend)

	select {
	case <-c.Done():
	case <-time.After(2 * time.Second):
		t.Fatal("connection did not shut down after fatal send failure")
	}
	select {
	case <-tr.firstRelease:
	default:
		t.Fatal("failed outgoing message was not released")
	}

	mu.Lock()
	gotOutcomes := outcomes
	mu.Unlock()
	if len(gotOutcomes[0]) != 1 || gotOutcomes[0][0].disposition != sendDeliveryAmbiguous || !gotOutcomes[0][0].fatal {
		t.Fatalf("first outcomes = %+v; want one fatal delivery-ambiguous outcome", gotOutcomes[0])
	}
	if len(gotOutcomes[1]) != 1 || gotOutcomes[1][0].disposition != sendAbortedByShutdown || gotOutcomes[1][0].fatal {
		t.Fatalf("queued outcomes = %+v; want one non-fatal shutdown-aborted outcome", gotOutcomes[1])
	}

	tr.mu.Lock()
	closeCount := tr.closeCount
	tr.mu.Unlock()
	if closeCount != 1 {
		t.Fatalf("transport close count = %d; want 1", closeCount)
	}

	logger.mu.Lock()
	loggedErrors := append([]string(nil), logger.errors...)
	logger.mu.Unlock()
	if len(loggedErrors) != 1 || !strings.Contains(loggedErrors[0], writeErr.Error()) {
		t.Fatalf("logged errors = %q; want one error containing %q", loggedErrors, writeErr)
	}
}
