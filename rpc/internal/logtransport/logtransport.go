// Package logtransport provides a transport that logs all of its messages.
package logtransport

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"golang.org/x/net/context"
	"zombiezen.com/go/capnproto/rpc"
	"zombiezen.com/go/capnproto/rpc/internal/logutil"
	"zombiezen.com/go/capnproto/rpc/rpccapnp"
)

type transport struct {
	rpc.Transport
	l       *log.Logger
	sendBuf bytes.Buffer
	recvBuf bytes.Buffer
}

// New creates a new logger that proxies messages to and from t and
// logs them to l.  If l is nil, then the log package's default
// logger is used.
func New(l *log.Logger, t rpc.Transport) rpc.Transport {
	return &transport{Transport: t, l: l}
}

func (t *transport) SendMessage(ctx context.Context, msg rpccapnp.Message) error {
	t.sendBuf.Reset()
	t.sendBuf.WriteString("<- ")
	formatMsg(&t.sendBuf, msg)
	logutil.Print(t.l, t.sendBuf.String())
	return t.Transport.SendMessage(ctx, msg)
}

func (t *transport) RecvMessage(ctx context.Context) (rpccapnp.Message, error) {
	msg, err := t.Transport.RecvMessage(ctx)
	if err != nil {
		return msg, err
	}
	t.recvBuf.Reset()
	t.recvBuf.WriteString("-> ")
	formatMsg(&t.recvBuf, msg)
	logutil.Print(t.l, t.recvBuf.String())
	return msg, nil
}

func formatMsg(w io.Writer, m rpccapnp.Message) {
	switch m.Which() {
	case rpccapnp.Message_Which_unimplemented:
		fmt.Fprint(w, "unimplemented")
	case rpccapnp.Message_Which_abort:
		fmt.Fprintf(w, "abort type=%v: %s", m.Abort().Type(), m.Abort().Reason())
	case rpccapnp.Message_Which_bootstrap:
		fmt.Fprintf(w, "bootstrap id=%d", m.Bootstrap().QuestionId())
	case rpccapnp.Message_Which_call:
		c := m.Call()
		fmt.Fprintf(w, "call id=%d target=<", c.QuestionId())
		formatMessageTarget(w, c.Target())
		fmt.Fprintf(w, "> @%#x/@%d", c.InterfaceId(), c.MethodId())
	case rpccapnp.Message_Which_return:
		r := m.Return()
		fmt.Fprintf(w, "return id=%d", r.AnswerId())
		if r.ReleaseParamCaps() {
			fmt.Fprint(w, " releaseParamCaps")
		}
		switch r.Which() {
		case rpccapnp.Return_Which_results:
		case rpccapnp.Return_Which_exception:
			fmt.Fprintf(w, ", exception type=%v: %s", r.Exception().Type(), r.Exception().Reason())
		case rpccapnp.Return_Which_canceled:
			fmt.Fprint(w, ", canceled")
		case rpccapnp.Return_Which_resultsSentElsewhere:
			fmt.Fprint(w, ", results sent elsewhere")
		case rpccapnp.Return_Which_takeFromOtherQuestion:
			fmt.Fprint(w, ", results sent elsewhere")
		case rpccapnp.Return_Which_acceptFromThirdParty:
			fmt.Fprint(w, ", accept from third party")
		default:
			fmt.Fprintf(w, ", UNKNOWN RESULT which=%v", r.Which())
		}
	case rpccapnp.Message_Which_finish:
		fmt.Fprintf(w, "finish id=%d", m.Finish().QuestionId())
		if m.Finish().ReleaseResultCaps() {
			fmt.Fprint(w, " releaseResultCaps")
		}
	case rpccapnp.Message_Which_resolve:
		r := m.Resolve()
		fmt.Fprintf(w, "resolve id=%d ", r.PromiseId())
		switch r.Which() {
		case rpccapnp.Resolve_Which_cap:
			fmt.Fprint(w, "capability=")
			formatCapDescriptor(w, r.Cap())
		case rpccapnp.Resolve_Which_exception:
			fmt.Fprintf(w, "exception type=%v: %s", r.Exception().Type(), r.Exception().Reason())
		default:
			fmt.Fprintf(w, "UNKNOWN RESOLUTION which=%v", r.Which())
		}
	case rpccapnp.Message_Which_release:
		fmt.Fprintf(w, "release id=%d by %d", m.Release().Id(), m.Release().ReferenceCount())
	case rpccapnp.Message_Which_disembargo:
		de := m.Disembargo()
		fmt.Fprint(w, "disembargo <")
		formatMessageTarget(w, de.Target())
		fmt.Fprint(w, "> ")
		dc := de.Context()
		switch dc.Which() {
		case rpccapnp.Disembargo_context_Which_senderLoopback:
			fmt.Fprintf(w, "sender loopback id=%d", dc.SenderLoopback())
		case rpccapnp.Disembargo_context_Which_receiverLoopback:
			fmt.Fprintf(w, "receiver loopback id=%d", dc.ReceiverLoopback())
		case rpccapnp.Disembargo_context_Which_accept:
			fmt.Fprint(w, "accept")
		case rpccapnp.Disembargo_context_Which_provide:
			fmt.Fprintf(w, "provide id=%d", dc.Provide())
		default:
			fmt.Fprintf(w, "UNKNOWN CONTEXT which=%v", dc.Which())
		}
	case rpccapnp.Message_Which_obsoleteSave:
		fmt.Fprint(w, "save")
	case rpccapnp.Message_Which_obsoleteDelete:
		fmt.Fprint(w, "delete")
	case rpccapnp.Message_Which_provide:
		fmt.Fprintf(w, "provide id=%d <", m.Provide().QuestionId())
		formatMessageTarget(w, m.Provide().Target())
		fmt.Fprint(w, ">")
	case rpccapnp.Message_Which_accept:
		fmt.Fprintf(w, "accept id=%d", m.Accept().QuestionId())
		if m.Accept().Embargo() {
			fmt.Fprint(w, " with embargo")
		}
	case rpccapnp.Message_Which_join:
		fmt.Fprintf(w, "join id=%d <", m.Join().QuestionId())
		formatMessageTarget(w, m.Join().Target())
		fmt.Fprint(w, ">")
	default:
		fmt.Fprintf(w, "UNKNOWN MESSAGE which=%v", m.Which())
	}
}

func formatMessageTarget(w io.Writer, t rpccapnp.MessageTarget) {
	switch t.Which() {
	case rpccapnp.MessageTarget_Which_importedCap:
		fmt.Fprintf(w, "import %d", t.ImportedCap())
	case rpccapnp.MessageTarget_Which_promisedAnswer:
		fmt.Fprint(w, "promise ")
		formatPromisedAnswer(w, t.PromisedAnswer())
	default:
		fmt.Fprintf(w, "UNKNOWN TARGET which=%v", t.Which())
	}
}

func formatPromisedAnswer(w io.Writer, a rpccapnp.PromisedAnswer) {
	fmt.Fprintf(w, "(question %d)", a.QuestionId())
	for i := 0; i < a.Transform().Len(); i++ {
		t := a.Transform().At(i)
		switch t.Which() {
		case rpccapnp.PromisedAnswer_Op_Which_noop:
		case rpccapnp.PromisedAnswer_Op_Which_getPointerField:
			fmt.Fprintf(w, ".getPointerField(%d)", t.GetPointerField())
		default:
			fmt.Fprintf(w, ".UNKNOWN(%v)", t.Which())
		}
	}
}

func formatCapDescriptor(w io.Writer, c rpccapnp.CapDescriptor) {
	switch c.Which() {
	case rpccapnp.CapDescriptor_Which_none:
		fmt.Fprint(w, "none")
	case rpccapnp.CapDescriptor_Which_senderHosted:
		fmt.Fprintf(w, "sender-hosted %d", c.SenderHosted())
	case rpccapnp.CapDescriptor_Which_senderPromise:
		fmt.Fprintf(w, "sender promise %d", c.SenderPromise())
	case rpccapnp.CapDescriptor_Which_receiverHosted:
		fmt.Fprintf(w, "receiver-hosted %d", c.ReceiverHosted())
	case rpccapnp.CapDescriptor_Which_receiverAnswer:
		fmt.Fprint(w, "receiver answer ")
		formatPromisedAnswer(w, c.ReceiverAnswer())
	case rpccapnp.CapDescriptor_Which_thirdPartyHosted:
		fmt.Fprint(w, "third-party hosted")
	default:
		fmt.Fprintf(w, "UNKNOWN CAPABILITY which=%v", c.Which())
	}
}
