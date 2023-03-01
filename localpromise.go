package capnp

import (
	"context"
)

// ClientHook for a promise that will be resolved to some other capability
// at some point. Buffers calls in a queue until the promsie is fulfilled,
// then forwards them.
type localPromise struct {
	aq *AnswerQueue
}

func NewLocalPromise() (Client, Fulfiller[Client]) {
	lp := newLocalPromise()
	p, f := NewPromisedClient(lp)
	return p, localFulfiller{
		lp:              lp,
		clientFulfiller: f,
	}
}

func newLocalPromise() localPromise {
	return localPromise{NewAnswerQueue(Method{})}
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

func (lp localPromise) Fulfill(c Client) {
	msg, seg := NewSingleSegmentMessage(nil)
	capID := msg.AddCap(c)
	lp.aq.Fulfill(NewInterface(seg, capID).ToPtr())
}

func (lp localPromise) Reject(err error) {
	lp.aq.Reject(err)
}

type localFulfiller struct {
	lp              localPromise
	clientFulfiller Fulfiller[Client]
}

func (lf localFulfiller) Fulfill(c Client) {
	lf.lp.Fulfill(c)
	lf.clientFulfiller.Fulfill(c)
}

func (lf localFulfiller) Reject(err error) {
	lf.lp.Reject(err)
	lf.clientFulfiller.Reject(err)
}
