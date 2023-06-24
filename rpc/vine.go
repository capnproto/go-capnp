package rpc

import (
	"context"

	"capnproto.org/go/capnp/v3"
)

// A vine is an implementation of capnp.ClientHook intended to be used to
// implement the logic discussed in rpc.capnp under
// ThirdPartyCapDescriptor.vineId. It forwards calls to an underlying
// capnp.ClientSnapshot
type vine struct {
	used         bool
	providerConn *Conn
	provideID    questionID
	snapshot     capnp.ClientSnapshot
}

func newVine(c *Conn, qid questionID, snapshot capnp.ClientSnapshot) *vine {
	return &vine{
		providerConn: c,
		provideID:    qid,
		snapshot:     snapshot,
	}
}

func (v *vine) Send(ctx context.Context, s capnp.Send) (*capnp.Answer, capnp.ReleaseFunc) {
	v.markUsed()
	return v.snapshot.Send(ctx, s)
}

func (v *vine) Recv(ctx context.Context, r capnp.Recv) capnp.PipelineCaller {
	v.markUsed()
	return v.snapshot.Recv(ctx, r)
}

func (v *vine) Brand() capnp.Brand {
	return v.snapshot.Brand()
}

func (v *vine) Shutdown() {
	v.snapshot.Release()
}

func (v *vine) String() string {
	// TODO: include other fields?
	return "&vine{snapshot: " + v.snapshot.String() + ", ...}"
}

func (v *vine) markUsed() {
	if v.used {
		return
	}
	v.used = true
	go func() {
		v.providerConn.withLocked(func(c *lockedConn) {
			panic("TODO: send finish, manipulate tables...")
		})
	}()
}
