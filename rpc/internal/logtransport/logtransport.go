// Package logtransport provides a transport that logs all of its messages.
package logtransport

import (
	"bytes"
	"fmt"
	"io"
	"log"

	"golang.org/x/net/context"
	"zombiezen.com/go/capnproto2/rpc"
	"zombiezen.com/go/capnproto2/rpc/internal/logutil"
	rpccapnp "zombiezen.com/go/capnproto2/std/capnp/rpc"
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
		mabort, _ := m.Abort()
		reason, _ := mabort.Reason()
		fmt.Fprintf(w, "abort type=%v: %s", mabort.Type(), reason)
	case rpccapnp.Message_Which_bootstrap:
		mboot, _ := m.Bootstrap()
		fmt.Fprintf(w, "bootstrap id=%d", mboot.QuestionId())
	case rpccapnp.Message_Which_call:
		c, _ := m.Call()
		fmt.Fprintf(w, "call id=%d target=<", c.QuestionId())
		tgt, _ := c.Target()
		formatMessageTarget(w, tgt)
		fmt.Fprintf(w, "> @%#x/@%d", c.InterfaceId(), c.MethodId())
	case rpccapnp.Message_Which_return:
		r, _ := m.Return()
		fmt.Fprintf(w, "return id=%d", r.AnswerId())
		if r.ReleaseParamCaps() {
			fmt.Fprint(w, " releaseParamCaps")
		}
		switch r.Which() {
		case rpccapnp.Return_Which_results:
		case rpccapnp.Return_Which_exception:
			exc, _ := r.Exception()
			reason, _ := exc.Reason()
			fmt.Fprintf(w, ", exception type=%v: %s", exc.Type(), reason)
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
		fin, _ := m.Finish()
		fmt.Fprintf(w, "finish id=%d", fin.QuestionId())
		if fin.ReleaseResultCaps() {
			fmt.Fprint(w, " releaseResultCaps")
		}
	case rpccapnp.Message_Which_resolve:
		r, _ := m.Resolve()
		fmt.Fprintf(w, "resolve id=%d ", r.PromiseId())
		switch r.Which() {
		case rpccapnp.Resolve_Which_cap:
			fmt.Fprint(w, "capability=")
			c, _ := r.Cap()
			formatCapDescriptor(w, c)
		case rpccapnp.Resolve_Which_exception:
			exc, _ := r.Exception()
			reason, _ := exc.Reason()
			fmt.Fprintf(w, "exception type=%v: %s", exc.Type(), reason)
		default:
			fmt.Fprintf(w, "UNKNOWN RESOLUTION which=%v", r.Which())
		}
	case rpccapnp.Message_Which_release:
		rel, _ := m.Release()
		fmt.Fprintf(w, "release id=%d by %d", rel.Id(), rel.ReferenceCount())
	case rpccapnp.Message_Which_disembargo:
		de, _ := m.Disembargo()
		tgt, _ := de.Target()
		fmt.Fprint(w, "disembargo <")
		formatMessageTarget(w, tgt)
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
		prov, _ := m.Provide()
		tgt, _ := prov.Target()
		fmt.Fprintf(w, "provide id=%d <", prov.QuestionId())
		formatMessageTarget(w, tgt)
		fmt.Fprint(w, ">")
	case rpccapnp.Message_Which_accept:
		acc, _ := m.Accept()
		fmt.Fprintf(w, "accept id=%d", acc.QuestionId())
		if acc.Embargo() {
			fmt.Fprint(w, " with embargo")
		}
	case rpccapnp.Message_Which_join:
		join, _ := m.Join()
		tgt, _ := join.Target()
		fmt.Fprintf(w, "join id=%d <", join.QuestionId())
		formatMessageTarget(w, tgt)
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
		pa, _ := t.PromisedAnswer()
		formatPromisedAnswer(w, pa)
	default:
		fmt.Fprintf(w, "UNKNOWN TARGET which=%v", t.Which())
	}
}

func formatPromisedAnswer(w io.Writer, a rpccapnp.PromisedAnswer) {
	fmt.Fprintf(w, "(question %d)", a.QuestionId())
	trans, _ := a.Transform()
	for i := 0; i < trans.Len(); i++ {
		t := trans.At(i)
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
		ans, _ := c.ReceiverAnswer()
		fmt.Fprint(w, "receiver answer ")
		formatPromisedAnswer(w, ans)
	case rpccapnp.CapDescriptor_Which_thirdPartyHosted:
		fmt.Fprint(w, "third-party hosted")
	default:
		fmt.Fprintf(w, "UNKNOWN CAPABILITY which=%v", c.Which())
	}
}
