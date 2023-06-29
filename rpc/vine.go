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
	snapshot      capnp.ClientSnapshot
	cancelProvide context.CancelFunc
}

func newVine(snapshot capnp.ClientSnapshot, cancelProvide context.CancelFunc) *vine {
	return &vine{
		snapshot:      snapshot,
		cancelProvide: cancelProvide,
	}
}

func (v *vine) Send(ctx context.Context, s capnp.Send) (*capnp.Answer, capnp.ReleaseFunc) {
	v.cancelProvide()
	return v.snapshot.Send(ctx, s)
}

func (v *vine) Recv(ctx context.Context, r capnp.Recv) capnp.PipelineCaller {
	v.cancelProvide()
	return v.snapshot.Recv(ctx, r)
}

func (v *vine) Brand() capnp.Brand {
	return v.snapshot.Brand()
}

func (v *vine) Shutdown() {
	v.cancelProvide()
	v.snapshot.Release()
}

func (v *vine) String() string {
	return "&vine{snapshot: " + v.snapshot.String() + "}"
}
