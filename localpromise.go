package capnp

// ClientHook for a promise that will be resolved to some other capability
// at some point. Buffers calls in a queue until the promise is fulfilled,
// then forwards them.
type localPromise struct {
	aq *AnswerQueue
}

// NewLocalPromise returns a client that will eventually resolve to a capability,
// supplied via resolver. resolver.Fulfill steals the reference to its argument.
func NewLocalPromise[C ~ClientKind]() (promise C, resolver Resolver[C]) {
	aq := NewAnswerQueue(Method{})
	f := NewPromise(Method{}, aq, aq)
	p := f.Answer().Client().AddRef()
	return C(p), localResolver[C]{
		p: f,
	}
}

type localResolver[C ~ClientKind] struct {
	p *Promise
}

func (lf localResolver[C]) Fulfill(c C) {
	msg, seg := NewSingleSegmentMessage(nil)
	capID := msg.CapTable().Add(Client(c))
	iface := NewInterface(seg, capID)
	lf.p.Fulfill(iface.ToPtr())
	lf.p.ReleaseClients()
	iface.Client().Release()
}

func (lf localResolver[C]) Reject(err error) {
	lf.p.Reject(err)
	lf.p.ReleaseClients()
}
