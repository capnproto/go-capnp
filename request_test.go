package capnp

import (
	"context"
	"encoding/binary"
	"strings"
	"testing"
)

type requestReleaseArena struct {
	Arena
	releases int
}

func (a *requestReleaseArena) Release() {
	a.releases++
	a.Arena.Release()
}

type requestLifecycleHook struct {
	calls        int
	placeArgs    bool
	placeArgsErr error
}

func (h *requestLifecycleHook) Send(_ context.Context, s Send) (*Answer, ReleaseFunc) {
	h.calls++
	if h.placeArgs {
		_, seg := NewSingleSegmentMessage(nil)
		args, err := NewStruct(seg, s.ArgsSize)
		if err != nil {
			return ErrorAnswer(s.Method, err), func() {}
		}
		if err := s.PlaceArgs(args); err != nil {
			h.placeArgsErr = err
			return ErrorAnswer(s.Method, err), func() {}
		}
	}
	return ImmediateAnswer(s.Method, newEmptyStruct().ToPtr()), func() {}
}

func (*requestLifecycleHook) Recv(context.Context, Recv) PipelineCaller { return nil }
func (*requestLifecycleHook) Brand() Brand                              { return Brand{} }
func (*requestLifecycleHook) Shutdown()                                 {}
func (*requestLifecycleHook) String() string                            { return "requestLifecycleHook" }

func newLifecycleRequest(t *testing.T, hook *requestLifecycleHook, size ObjectSize) (*Request, *requestReleaseArena, Client) {
	t.Helper()
	arena := &requestReleaseArena{Arena: MultiSegment(nil)}
	_, seg, err := NewMessage(arena)
	if err != nil {
		t.Fatal(err)
	}
	args, err := NewStruct(seg, size)
	if err != nil {
		t.Fatal(err)
	}
	client := NewClient(hook)
	return &Request{method: dummyMethod, args: args, client: client}, arena, client
}

func TestRequestReleaseArgsExactlyOnce(t *testing.T) {
	hook := &requestLifecycleHook{}
	req, arena, client := newLifecycleRequest(t, hook, ObjectSize{})
	defer client.Release()

	req.Release()
	req.Release()

	if arena.releases != 1 {
		t.Fatalf("argument arena release count = %d; want 1", arena.releases)
	}
	if req.Args().IsValid() {
		t.Fatal("request arguments remain valid after Release")
	}

	// A valid Struct can still have no message when constructed from an
	// unbound segment. Releasing it must be a no-op rather than a panic.
	noMessage := &Request{args: Struct{seg: new(Segment)}}
	noMessage.releaseArgs()
	if noMessage.Args().IsValid() {
		t.Fatal("request arguments with no message remain valid after release")
	}
}

func TestRequestSendReleasesArgsWhenNotPlaced(t *testing.T) {
	hook := &requestLifecycleHook{}
	req, arena, client := newLifecycleRequest(t, hook, ObjectSize{})
	defer client.Release()

	if future := req.Send(context.Background()); future == nil {
		t.Fatal("Send returned a nil future")
	}
	req.Release()

	if arena.releases != 1 {
		t.Fatalf("argument arena release count = %d; want 1", arena.releases)
	}
}

func TestRequestSendReleasesArgsAfterCopyFailure(t *testing.T) {
	hook := &requestLifecycleHook{placeArgs: true}
	req, arena, client := newLifecycleRequest(t, hook, ObjectSize{PointerCount: 1})
	defer client.Release()

	// A malformed pointer makes Struct.CopyFrom fail while placing arguments.
	binary.LittleEndian.PutUint64(req.args.seg.data[req.args.off:], ^uint64(0))
	future := req.Send(context.Background())
	if _, err := future.Struct(); err == nil {
		t.Fatal("Send error = nil; want argument-copy failure")
	}
	if hook.placeArgsErr == nil {
		t.Fatal("PlaceArgs did not observe the argument-copy failure")
	}
	req.Release()

	if arena.releases != 1 {
		t.Fatalf("argument arena release count = %d; want 1", arena.releases)
	}
}

func TestRequestSendAndSendStreamAreSingleUse(t *testing.T) {
	for _, test := range []struct {
		name   string
		first  func(*Request) error
		second func(*Request) error
	}{
		{
			name:  "SendThenSend",
			first: func(r *Request) error { r.Send(context.Background()); return nil },
			second: func(r *Request) error {
				_, err := r.Send(context.Background()).Struct()
				return err
			},
		},
		{
			name:   "SendStreamThenSendStream",
			first:  func(r *Request) error { return r.SendStream(context.Background()) },
			second: func(r *Request) error { return r.SendStream(context.Background()) },
		},
		{
			name:   "SendThenSendStream",
			first:  func(r *Request) error { r.Send(context.Background()); return nil },
			second: func(r *Request) error { return r.SendStream(context.Background()) },
		},
		{
			name:  "SendStreamThenSend",
			first: func(r *Request) error { return r.SendStream(context.Background()) },
			second: func(r *Request) error {
				_, err := r.Send(context.Background()).Struct()
				return err
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			hook := &requestLifecycleHook{placeArgs: true}
			req, arena, client := newLifecycleRequest(t, hook, ObjectSize{})
			defer client.Release()

			if err := test.first(req); err != nil {
				t.Fatal(err)
			}
			if err := test.second(req); err == nil || !strings.Contains(err.Error(), "sent the same request twice") {
				t.Fatalf("second send error = %v; want sent the same request twice", err)
			}
			req.Release()

			if hook.calls != 1 {
				t.Errorf("hook calls = %d; want 1", hook.calls)
			}
			if arena.releases != 1 {
				t.Errorf("argument arena release count = %d; want 1", arena.releases)
			}
		})
	}
}
