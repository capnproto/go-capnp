package capnp

import (
	"context"
)

// ClientHook for a promise that will be resolved to some other capability
// at some point. Buffers calls in a queue until the promise is fulfilled,
// then forwards them.
type localPromise struct {
	aq *AnswerQueue
}

// NewLocalPromise returns a client that will eventually resolve to a capability,
// supplied via resolver. resolver.Fulfill steals the reference to its argument.
func NewLocalPromise[C ~ClientKind]() (promise C, resolver Resolver[C]) {
	lp := newLocalPromise()
	p, f := NewPromisedClient(lp)
	return C(p), localResolver[C]{
		lp:             lp,
		clientResolver: f,
	}
}

func newLocalPromise() localPromise {
	return localPromise{aq: NewAnswerQueue(Method{})}
}

func (lp localPromise) Send(ctx context.Context, s Send) (*Answer, ReleaseFunc) {
	return lp.aq.PipelineSend(ctx, nil, s)
}

func (lp localPromise) Recv(ctx context.Context, r Recv) PipelineCaller {
	return lp.aq.PipelineRecv(ctx, nil, r)
}

func (lp localPromise) Brand() Brand {
	return Brand{}
}

func (lp localPromise) Shutdown() {}

func (lp localPromise) String() string {
	return "localPromise{...}"
}

type localResolver[C ~ClientKind] struct {
	lp             localPromise
	clientResolver Resolver[Client]
}

func (lf localResolver[C]) Fulfill(c C) {
	// This is convoluted, for a few reasons:
	//
	// 1. AnswerQueue wants a Ptr, not a Client, so we have to construct
	//    a message for this.
	// 2. Message.AddCap steals the reference, so we have to get the client
	//    back from the message instead of using the reference we already
	//    have.
	// 3. The semantics of NewPromisedClient differ from what we want, and
	//    are kindof odd: when it is resolved it does not steal the reference,
	//    nor does it borrow it -- instead it merges the refcounts, so that
	//    the two clients point to the same place. So we have to drop our
	//    reference to get the semantics we want.
	//
	//    TODO: We should probably push this part down into the implementation of
	//    clientPromise. That requires auditing its uses and adjusting call sites
	//    though.
	msg, seg := NewSingleSegmentMessage(nil)
	capID := msg.AddCap(Client(c))
	iface := NewInterface(seg, capID)
	client := iface.Client()
	defer client.Release()
	lf.lp.aq.Fulfill(iface.ToPtr())
	lf.clientResolver.Fulfill(client)
}

func (lf localResolver[C]) Reject(err error) {
	lf.lp.aq.Reject(err)
	lf.clientResolver.Reject(err)
}
