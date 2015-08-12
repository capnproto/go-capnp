package rpc

import (
	"log"

	"zombiezen.com/go/capnproto"
	"zombiezen.com/go/capnproto/internal/fulfiller"
	"zombiezen.com/go/capnproto/rpc/rpccapnp"
)

// nestedCall is called from the coordinate goroutine to make a client call.
// Since the client may point
func (c *Conn) nestedCall(client capnp.Client, cl *capnp.Call) capnp.Answer {
	client = extractRPCClient(client)
	ac := appCallFromClientCall(c, client, cl)
	if ac != nil {
		ans, err := c.handleCall(ac)
		if err != nil {
			log.Println("rpc: failed to handle call:", err)
			return capnp.ErrorAnswer(err)
		}
		return ans
	}
	// TODO(light): Add a CallOption that signals to bypass sync.
	// The above hack works in *most* cases.
	//
	// If your code is deadlocking here, you've hit the edge of the
	// compromise between these three goals:
	// 1) Package capnp is loosely coupled with package rpc
	// 2) Arbitrary implementations of Client may exist
	// 3) Local E-order must be preserved
	//
	// #3 is the one that creates a goroutine send cycle, since
	// application code must synchronize with the coordinate goroutine
	// to preserve order of delivery.  You can't really overcome this
	// without breaking one of the first two constraints.
	//
	// To avoid #2 as much as possible, implementing Client is discouraged
	// by several docs.
	return client.Call(cl)
}

func (c *Conn) descriptorForClient(desc rpccapnp.CapDescriptor, client capnp.Client) {
	client = extractRPCClient(client)
	if ic, ok := client.(*importClient); ok && isImportFromConn(ic, c) {
		desc.SetReceiverHosted(uint32(ic.id))
		return
	}
	if pc, ok := client.(*capnp.PipelineClient); ok {
		p := (*capnp.Pipeline)(pc)
		if q, ok := p.Answer().(*question); ok && isQuestionFromConn(q, c) {
			a := rpccapnp.NewPromisedAnswer(desc.Segment())
			a.SetQuestionId(uint32(q.id))
			transformToPromisedAnswer(desc.Segment(), a, p.Transform())
			desc.SetReceiverAnswer(a)
			return
		}
	}
	id := c.exports.add(client)
	desc.SetSenderHosted(uint32(id))
}

func appCallFromClientCall(c *Conn, client capnp.Client, cl *capnp.Call) *appCall {
	if ic, ok := client.(*importClient); ok && isImportFromConn(ic, c) {
		ac, _ := newAppImportCall(ic.id, cl)
		return ac
	}
	if pc, ok := client.(*capnp.PipelineClient); ok {
		p := (*capnp.Pipeline)(pc)
		if q, ok := p.Answer().(*question); ok && isQuestionFromConn(q, c) {
			ac, _ := newAppPipelineCall(q, p.Transform(), cl)
			return ac
		}
	}
	return nil
}

// extractRPCClient attempts to extract the client that is the most
// meaningful for further processing of RPCs.  For example, instead of a
// PipelineClient on a resolved answer, the client's capability.
func extractRPCClient(client capnp.Client) capnp.Client {
	for {
		switch c := client.(type) {
		case *importClient:
			return c
		case *capnp.PipelineClient:
			p := (*capnp.Pipeline)(c)
			next := extractRPCClientFromPipeline(p.Answer(), p.Transform())
			if next == nil {
				return client
			}
			client = next
		case clientWrapper:
			wc := c.WrappedClient()
			if wc == nil {
				return client
			}
			client = wc
		default:
			return client
		}
	}
}

func extractRPCClientFromPipeline(ans capnp.Answer, transform []capnp.PipelineOp) capnp.Client {
	if capnp.IsFixedAnswer(ans) {
		s, err := ans.Struct()
		return clientFromResolution(transform, capnp.Pointer(s), err)
	}
	switch a := ans.(type) {
	case *fulfiller.Fulfiller:
		ap := a.Peek()
		if ap == nil {
			// This can race, see TODO in nestedCall.
			return nil
		}
		s, err := ap.Struct()
		return clientFromResolution(transform, capnp.Pointer(s), err)
	case *question:
		_, obj, err, ok := a.peek()
		if !ok {
			// This can race, see TODO in nestedCall.
			return nil
		}
		return clientFromResolution(transform, obj, err)
	default:
		return nil
	}
}

// clientWrapper is an interface for types that wrap clients.
// If WrappedClient returns a non-nil value, that means that a Call to
// the wrapper passes through to the returned client.
// TODO(light): this should probably be exported at some point.
type clientWrapper interface {
	WrappedClient() capnp.Client
}

func isQuestionFromConn(q *question, c *Conn) bool {
	// TODO(light): ideally there would be better ways to check.
	return q.manager == &c.manager
}

func isImportFromConn(ic *importClient, c *Conn) bool {
	// TODO(light): ideally there would be better ways to check.
	return ic.manager == &c.manager
}
