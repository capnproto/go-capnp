package rpc

import (
	"context"
	"errors"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/exc"
	"capnproto.org/go/capnp/v3/exp/spsc"
	rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"
	"capnproto.org/go/capnp/v3/util/deferred"
)

func TestUnsupportedLevel3MessagesEchoUnimplemented(t *testing.T) {
	tests := []struct {
		name  string
		build func(rpccp.Message) error
		want  rpccp.Message_Which
	}{
		{
			name: "provide",
			build: func(message rpccp.Message) error {
				_, err := message.NewProvide()
				return err
			},
			want: rpccp.Message_Which_provide,
		},
		{
			name: "accept",
			build: func(message rpccp.Message) error {
				_, err := message.NewAccept()
				return err
			},
			want: rpccp.Message_Which_accept,
		},
		{
			name: "thirdPartyAnswer",
			build: func(message rpccp.Message) error {
				_, err := message.NewThirdPartyAnswer()
				return err
			},
			want: rpccp.Message_Which_thirdPartyAnswer,
		},
		{
			name: "call.sendResultsTo.thirdParty",
			build: func(message rpccp.Message) error {
				call, err := message.NewCall()
				if err != nil {
					return err
				}
				call.SetQuestionId(41)
				return call.SendResultsTo().SetThirdParty(capnp.Ptr{})
			},
			want: rpccp.Message_Which_call,
		},
		{
			name: "disembargo.context.accept",
			build: func(message rpccp.Message) error {
				disembargo, err := message.NewDisembargo()
				if err != nil {
					return err
				}
				target, err := disembargo.NewTarget()
				if err != nil {
					return err
				}
				target.SetImportedCap(0)
				return disembargo.Context().SetAccept([]byte("unsupported"))
			},
			want: rpccp.Message_Which_disembargo,
		},
	}

	for _, networked := range []bool{false, true} {
		for _, test := range tests {
			t.Run(test.name+"/network="+boolString(networked), func(t *testing.T) {
				var opts *Options
				if networked {
					opts = &Options{Network: resolveTestNetwork{}}
				}
				conn, peer := newResolveLifecycleConnWithOptions(t, opts)

				out, err := peer.NewMessage()
				if err != nil {
					t.Fatal(err)
				}
				if err := test.build(out.Message()); err != nil {
					out.Release()
					t.Fatal(err)
				}
				if err := out.Send(); err != nil {
					out.Release()
					t.Fatal(err)
				}
				out.Release()

				in := recvPeerMessage(t, peer)
				defer in.Release()
				if got := in.Message().Which(); got != rpccp.Message_Which_unimplemented {
					t.Fatalf("response type = %v; want unimplemented", got)
				}
				echo, err := in.Message().Unimplemented()
				if err != nil {
					t.Fatal(err)
				}
				if got := echo.Which(); got != test.want {
					t.Fatalf("echoed message type = %v; want %v", got, test.want)
				}
				select {
				case <-conn.Done():
					t.Fatal("unsupported Level 3 message closed connection")
				default:
				}
				requireEmptyLevel3Tables(t, conn)
			})
		}
	}
}

func TestUnsupportedLevel3MessageReleaseOnNewMessageFailure(t *testing.T) {
	tests := []struct {
		name   string
		build  func(rpccp.Message) error
		handle func(*Conn, *countingIncomingMessage)
	}{
		{
			name: "provide",
			build: func(message rpccp.Message) error {
				_, err := message.NewProvide()
				return err
			},
		},
		{
			name: "accept",
			build: func(message rpccp.Message) error {
				_, err := message.NewAccept()
				return err
			},
		},
		{
			name: "thirdPartyAnswer",
			build: func(message rpccp.Message) error {
				_, err := message.NewThirdPartyAnswer()
				return err
			},
		},
		{
			name: "call.sendResultsTo.thirdParty",
			build: func(message rpccp.Message) error {
				call, err := message.NewCall()
				if err != nil {
					return err
				}
				return call.SendResultsTo().SetThirdParty(capnp.Ptr{})
			},
			handle: func(conn *Conn, in *countingIncomingMessage) {
				_ = conn.handleCall(context.Background(), in)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			queue := spsc.New[asyncSend]()
			conn := &Conn{
				transport: newMessageErrorTransport{err: errors.New("allocate")},
			}
			conn.lk.sendTx = &queue.Tx

			_, seg := capnp.NewSingleSegmentMessage(nil)
			message, err := rpccp.NewRootMessage(seg)
			if err != nil {
				t.Fatal(err)
			}
			if err := test.build(message); err != nil {
				t.Fatal(err)
			}
			in := &countingIncomingMessage{message: message}

			if test.handle == nil {
				conn.handleUnknownMessageType(context.Background(), in)
			} else {
				test.handle(conn, in)
			}
			pending, ok := queue.Rx.TryRecv()
			if !ok {
				t.Fatal("unimplemented response was not queued")
			}
			if err := pending.Send(); err != nil {
				t.Fatalf("send outcome = %v; want non-fatal allocation failure", err)
			}
			if got := atomic.LoadInt32(&in.releases); got != 1 {
				t.Fatalf("incoming message releases = %d; want 1", got)
			}
		})
	}
}

func TestThirdPartyHostedCapabilityFallback(t *testing.T) {
	for _, networked := range []bool{false, true} {
		t.Run("network="+boolString(networked), func(t *testing.T) {
			var opts *Options
			if networked {
				opts = &Options{Network: resolveTestNetwork{}}
			}
			conn, _ := newResolveLifecycleConnWithOptions(t, opts)

			_, seg := capnp.NewSingleSegmentMessage(nil)
			desc, err := rpccp.NewRootCapDescriptor(seg)
			if err != nil {
				t.Fatal(err)
			}
			thirdParty, err := desc.NewThirdPartyHosted()
			if err != nil {
				t.Fatal(err)
			}
			const vineID = importID(17)
			thirdParty.SetVineId(uint32(vineID))

			var client capnp.Client
			conn.withLocked(func(c *lockedConn) {
				client, err = c.recvCap(desc)
				_, imported := c.lk.imports.Find(vineID)
				if networked && imported {
					t.Error("unsupported third-party descriptor mutated imports")
				}
				if !networked && !imported {
					t.Error("vine fallback did not create import")
				}
			})
			if networked {
				if err == nil || !strings.Contains(err.Error(), "three-party handoff not implemented") {
					t.Fatalf("recvCap error = %v; want unsupported handoff error", err)
				}
				if got := exc.TypeOf(err); got != exc.Unimplemented {
					t.Fatalf("recvCap error type = %v; want unimplemented", got)
				}
				if client.IsValid() {
					client.Release()
					t.Fatal("recvCap returned a client after rejecting descriptor")
				}
			} else {
				if err != nil {
					client.Release()
					t.Fatal(err)
				}
				if !client.IsValid() {
					t.Fatal("vine fallback returned invalid client")
				}
				client.Release()
			}
		})
	}
}

func TestNetworkedThirdPartyPayloadReleasesPartialCapabilityTable(t *testing.T) {
	conn, _ := newResolveLifecycleConnWithOptions(t, &Options{Network: resolveTestNetwork{}})

	message, seg := capnp.NewSingleSegmentMessage(nil)
	defer message.Release()
	payload, err := rpccp.NewRootPayload(seg)
	if err != nil {
		t.Fatal(err)
	}
	descriptors, err := payload.NewCapTable(2)
	if err != nil {
		t.Fatal(err)
	}
	const (
		firstImport = importID(23)
		vineImport  = importID(24)
	)
	descriptors.At(0).SetSenderHosted(uint32(firstImport))
	thirdParty, err := descriptors.At(1).NewThirdPartyHosted()
	if err != nil {
		t.Fatal(err)
	}
	thirdParty.SetVineId(uint32(vineImport))

	dq := &deferred.Queue{}
	conn.withLocked(func(c *lockedConn) {
		_, _, err = c.recvPayload(dq, payload)
	})
	dq.Run()
	if err == nil || !strings.Contains(err.Error(), "three-party handoff not implemented") {
		t.Fatalf("recvPayload error = %v; want unsupported handoff error", err)
	}
	if got := message.CapTable().Len(); got != 0 {
		t.Errorf("partially decoded message capability table length = %d; want 0", got)
	}
	conn.withLocked(func(c *lockedConn) {
		if _, ok := c.lk.imports.Find(firstImport); ok {
			t.Error("partially decoded sender-hosted capability remains imported")
		}
		if _, ok := c.lk.imports.Find(vineImport); ok {
			t.Error("rejected third-party vine was imported")
		}
	})
}

func TestHandleCallRejectsNetworkedThirdPartyPayloadCleanly(t *testing.T) {
	conn, peer := newResolveLifecycleConnWithOptions(t, &Options{Network: resolveTestNetwork{}})

	_, seg := capnp.NewSingleSegmentMessage(nil)
	message, err := rpccp.NewRootMessage(seg)
	if err != nil {
		t.Fatal(err)
	}
	call, err := message.NewCall()
	if err != nil {
		t.Fatal(err)
	}
	const questionID = answerID(55)
	call.SetQuestionId(uint32(questionID))
	target, err := call.NewTarget()
	if err != nil {
		t.Fatal(err)
	}
	target.SetImportedCap(0)
	params, err := call.NewParams()
	if err != nil {
		t.Fatal(err)
	}
	descriptors, err := params.NewCapTable(2)
	if err != nil {
		t.Fatal(err)
	}
	const (
		firstImport = importID(25)
		vineImport  = importID(26)
	)
	descriptors.At(0).SetSenderHosted(uint32(firstImport))
	thirdParty, err := descriptors.At(1).NewThirdPartyHosted()
	if err != nil {
		t.Fatal(err)
	}
	thirdParty.SetVineId(uint32(vineImport))
	in := &countingIncomingMessage{message: message}

	if err := conn.handleCall(context.Background(), in); err != nil {
		t.Fatal(err)
	}
	if got := atomic.LoadInt32(&in.releases); got != 1 {
		t.Fatalf("incoming message releases = %d; want 1", got)
	}
	conn.withLocked(func(c *lockedConn) {
		if _, ok := c.lk.imports.Find(firstImport); ok {
			t.Error("partially decoded sender-hosted capability remains imported")
		}
		if _, ok := c.lk.imports.Find(vineImport); ok {
			t.Error("rejected third-party vine was imported")
		}
	})

	response := recvPeerMessage(t, peer)
	if got := response.Message().Which(); got != rpccp.Message_Which_return {
		response.Release()
		t.Fatalf("response type = %v; want return", got)
	}
	ret, err := response.Message().Return()
	if err != nil {
		response.Release()
		t.Fatal(err)
	}
	exception, err := ret.Exception()
	if err != nil {
		response.Release()
		t.Fatal(err)
	}
	if got := exception.Type(); got != rpccp.Exception_Type_unimplemented {
		response.Release()
		t.Fatalf("return exception type = %v; want unimplemented", got)
	}
	response.Release()

	finish := newFinishMessage(t, questionID)
	if err := conn.handleFinish(context.Background(), finish); err != nil {
		t.Fatal(err)
	}
	if got := atomic.LoadInt32(&finish.releases); got != 1 {
		t.Fatalf("finish message releases = %d; want 1", got)
	}
	conn.withLocked(func(c *lockedConn) {
		if _, ok := c.lk.answers.Find(questionID); ok {
			t.Error("error answer remains after Finish")
		}
	})
	requireEmptyLevel3Tables(t, conn)
}

func TestSameNetworkCapabilitiesFallBackToProxying(t *testing.T) {
	tests := []struct {
		name      string
		newClient func(*testing.T, *Conn) capnp.Client
	}{
		{
			name: "import",
			newClient: func(t *testing.T, conn *Conn) (client capnp.Client) {
				conn.withLocked(func(c *lockedConn) {
					client = c.addImport(31, false)
				})
				return client
			},
		},
		{
			name: "pipeline",
			newClient: func(t *testing.T, conn *Conn) capnp.Client {
				return conn.Bootstrap(context.Background())
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			network := resolveTestNetwork{}
			dst, _ := newResolveLifecycleConnWithOptions(t, &Options{Network: network})
			src, _ := newResolveLifecycleConnWithOptions(t, &Options{Network: network})
			client := test.newClient(t, src)
			defer client.Release()

			_, seg := capnp.NewSingleSegmentMessage(nil)
			desc, err := rpccp.NewRootCapDescriptor(seg)
			if err != nil {
				t.Fatal(err)
			}

			var (
				exported exportID
				isExport bool
				isLocal  bool
			)
			dst.withLocked(func(c *lockedConn) {
				isLocal = c.isLocalClient(client)
				exported, isExport, err = c.sendCap(desc, client.Snapshot())
			})
			if err != nil {
				t.Fatal(err)
			}
			if !isLocal {
				t.Error("same-network proxied client not treated as local")
			}
			if !isExport {
				t.Fatal("same-network client was not proxied through an export")
			}
			switch got := desc.Which(); got {
			case rpccp.CapDescriptor_Which_senderHosted, rpccp.CapDescriptor_Which_senderPromise:
			default:
				t.Fatalf("descriptor type = %v; want sender-hosted or sender-promise proxy", got)
			}

			dq := &deferred.Queue{}
			dst.withLocked(func(c *lockedConn) {
				err = c.releaseExport(dq, exported, 1)
			})
			dq.Run()
			if err != nil {
				t.Fatal(err)
			}
			if err := src.Close(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestAwaitFromThirdPartyReturnFailsQuestionCleanly(t *testing.T) {
	for _, networked := range []bool{false, true} {
		t.Run("network="+boolString(networked), func(t *testing.T) {
			var opts *Options
			if networked {
				opts = &Options{Network: resolveTestNetwork{}}
			}
			conn, peer := newResolveLifecycleConnWithOptions(t, opts)
			client := conn.Bootstrap(context.Background())
			defer client.Release()

			bootstrap := recvPeerMessage(t, peer)
			boot, err := bootstrap.Message().Bootstrap()
			if err != nil {
				bootstrap.Release()
				t.Fatal(err)
			}
			questionID := boot.QuestionId()
			bootstrap.Release()

			out, err := peer.NewMessage()
			if err != nil {
				t.Fatal(err)
			}
			ret, err := out.Message().NewReturn()
			if err != nil {
				out.Release()
				t.Fatal(err)
			}
			ret.SetAnswerId(questionID)
			if err := ret.SetAwaitFromThirdParty(capnp.Ptr{}); err != nil {
				out.Release()
				t.Fatal(err)
			}
			if err := out.Send(); err != nil {
				out.Release()
				t.Fatal(err)
			}
			out.Release()

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			if err := client.Resolve(ctx); err != nil {
				t.Fatal(err)
			}
			resolved := client.Snapshot()
			resolutionErr, ok := resolved.Brand().Value.(error)
			resolved.Release()
			if !ok || !strings.Contains(resolutionErr.Error(), "awaitFromThirdParty") {
				t.Fatalf("bootstrap resolution = %v; want unsupported return error", resolutionErr)
			}
			if got := exc.TypeOf(resolutionErr); got != exc.Unimplemented {
				t.Fatalf("bootstrap resolution error type = %v; want unimplemented", got)
			}

			finish := recvPeerMessage(t, peer)
			defer finish.Release()
			if got := finish.Message().Which(); got != rpccp.Message_Which_finish {
				t.Fatalf("response type = %v; want finish", got)
			}
			fin, err := finish.Message().Finish()
			if err != nil {
				t.Fatal(err)
			}
			if got := fin.QuestionId(); got != questionID {
				t.Fatalf("finish question ID = %d; want %d", got, questionID)
			}
			select {
			case <-conn.Done():
				t.Fatal("unsupported Return closed connection")
			default:
			}
			requireEmptyLevel3Tables(t, conn)
		})
	}
}

func requireEmptyLevel3Tables(t *testing.T, conn *Conn) {
	t.Helper()

	var nonempty []string
	conn.withLocked(func(c *lockedConn) {
		c.lk.questions.Range(func(_ questionID, _ *question) bool {
			nonempty = append(nonempty, "questions")
			return false
		})
		c.lk.answers.Range(func(_ answerID, _ *ansent) bool {
			nonempty = append(nonempty, "answers")
			return false
		})
		c.lk.imports.Range(func(_ importID, _ *impent) bool {
			nonempty = append(nonempty, "imports")
			return false
		})
		c.lk.exports.Range(func(_ exportID, _ *expent) bool {
			nonempty = append(nonempty, "exports")
			return false
		})
		c.lk.embargoes.Range(func(_ embargoID, _ *embargo) bool {
			nonempty = append(nonempty, "embargoes")
			return false
		})
	})
	if len(nonempty) > 0 {
		t.Errorf("unexpected table entries after Level 3 rejection: %v", nonempty)
	}
}

func boolString(v bool) string {
	if v {
		return "on"
	}
	return "off"
}
