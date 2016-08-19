package rpc

import (
	"zombiezen.com/go/capnproto2"
	"zombiezen.com/go/capnproto2/internal/fulfiller"
	rpccapnp "zombiezen.com/go/capnproto2/std/capnp/rpc"
)

// lockedCall is used to make a call to an arbitrary client while
// holding onto c.mu.  Since the client could point back to c, naively
// calling c.Call could deadlock.
func (c *Conn) lockedCall(client capnp.Client, cl *capnp.Call) capnp.Answer {
	client = extractRPCClient(client)
	switch client := client.(type) {
	case *importClient:
		if client.conn != c {
			return client.Call(cl)
		}
		return client.lockedCall(cl)
	case *capnp.PipelineClient:
		p := (*capnp.Pipeline)(client)
		if q, ok := p.Answer().(*question); ok && q.conn == c {
			return q.lockedPipelineCall(p.Transform(), cl)
		}
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
	// #3 is the one that creates a deadlock, since application code must
	// acquire the connection mutex to preserve order of delivery.  You
	// can't really overcome this without breaking one of the first two
	// constraints.
	//
	// To avoid #2 as much as possible, implementing Client is discouraged
	// by several docs.
	return client.Call(cl)
}

func (c *Conn) descriptorForClient(desc rpccapnp.CapDescriptor, client capnp.Client) error {
	client = extractRPCClient(client)
	if ic, ok := client.(*importClient); ok && ic.conn == c {
		desc.SetReceiverHosted(uint32(ic.id))
		return nil
	}
	if pc, ok := client.(*capnp.PipelineClient); ok {
		p := (*capnp.Pipeline)(pc)
		if q, ok := p.Answer().(*question); ok && q.conn == c {
			a, err := desc.NewReceiverAnswer()
			if err != nil {
				return err
			}
			a.SetQuestionId(uint32(q.id))
			err = transformToPromisedAnswer(desc.Segment(), a, p.Transform())
			if err != nil {
				return err
			}
			return nil
		}
	}
	id := c.addExport(client)
	desc.SetSenderHosted(uint32(id))
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
		return clientFromResolution(transform, s.ToPtr(), err)
	}
	switch a := ans.(type) {
	case *fulfiller.Fulfiller:
		ap := a.Peek()
		if ap == nil {
			// This can race, see TODO in lockedCall.
			return nil
		}
		s, err := ap.Struct()
		return clientFromResolution(transform, s.ToPtr(), err)
	case *question:
		a.mu.RLock()
		obj, err, state := a.obj, a.err, a.state
		a.mu.RUnlock()
		if state == questionInProgress {
			// This can race, see TODO in lockedCall.
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
