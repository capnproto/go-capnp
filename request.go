package capnp

import (
	"context"
)

type Request struct {
	method          Method
	args            Struct
	client          Client
	releaseResponse ReleaseFunc
	future          *Future
}

func NewRequest(client Client, method Method, argsSize ObjectSize) (*Request, error) {
	_, seg, err := NewMessage(MultiSegment(nil))
	if err != nil {
		return nil, err
	}
	args, err := NewStruct(seg, argsSize)
	if err != nil {
		return nil, err
	}
	return &Request{
		method: method,
		args:   args,
		client: client,
	}, nil
}

func (r *Request) Args() Struct {
	return r.args
}

func (r *Request) getSend() Send {
	return Send{
		Method: r.method,
		PlaceArgs: func(args Struct) error {
			err := args.CopyFrom(r.args)
			r.releaseArgs()
			return err
		},
		ArgsSize: r.args.Size(),
	}
}

func (r *Request) Send(ctx context.Context) *Future {
	ans, rel := r.client.SendCall(ctx, r.getSend())
	r.releaseResponse = rel
	r.future = ans.Future()
	return r.future
}

func (r *Request) SendStream(ctx context.Context) error {
	return r.client.SendStreamCall(ctx, r.getSend())
}

func (r *Request) Future() *Future {
	return r.future
}

func (r *Request) Release() {
	r.releaseArgs()
	rel := r.releaseResponse
	if rel != nil {
		r.releaseResponse = nil
		r.future = nil
		rel()
	}
}

func (r *Request) releaseArgs() {
	if r.args.IsValid() {
		return
	}
	msg := r.args.Message()
	r.args = Struct{}
	arena := msg.Arena
	msg.Reset(nil)
	arena.Release()
}
