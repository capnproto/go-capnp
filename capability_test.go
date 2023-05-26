package capnp

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	ctx := context.Background()
	h := &dummyHook{brand: Brand{Value: int(42)}}
	c := NewClient(h)
	defer c.Release()

	if !c.IsSame(c) {
		t.Error("!c.IsSame(c)")
	}
	if !c.IsValid() {
		t.Error("new client is not valid")
	}
	state := c.Snapshot()
	if state.IsPromise() {
		t.Error("c.State().IsPromise = true; want false")
	}
	if state.Brand().Value != int(42) {
		t.Errorf("c.State().Brand().Value = %#v; want 42", state.Brand().Value)
	}
	state.Release()
	ans, finish := c.SendCall(ctx, Send{})
	if _, err := ans.Struct(); err != nil {
		t.Error("SendCall:", err)
	}
	finish()
	if h.calls != 1 {
		t.Errorf("after SendCall, h.calls = %d; want 1", h.calls)
	}
	ret := new(dummyReturner)
	pcall := c.RecvCall(ctx, Recv{
		ReleaseArgs: func() {},
		Returner:    ret,
	})
	if !ret.returned {
		t.Error("RecvCall did not return")
	} else {
		if ret.err != nil {
			t.Error("RecvCall returned error:", ret.err)
		}
		if pcall != nil {
			t.Error("RecvCall returned a PipelineCaller")
		}
	}
	if h.calls != 2 {
		t.Errorf("after RecvCall, h.calls = %d; want 2", h.calls)
	}
	rctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	if err := c.Resolve(rctx); err != nil {
		t.Error("Resolve failed:", err)
	}
	cancel()
	c.Release()
	if h.shutdowns == 0 {
		t.Error("Release did not call ClientHook.Shutdown")
	} else if h.shutdowns > 1 {
		t.Error("Release called ClientHook.Shutdown multiple times")
	}
}

func TestReleasedClient(t *testing.T) {
	ctx := context.Background()
	h := &dummyHook{brand: Brand{Value: int(42)}}
	c := NewClient(h)
	c.Release()

	if c.IsValid() {
		t.Error("released client is valid")
	}
	state := c.Snapshot()
	if state.Brand().Value != nil {
		t.Errorf("c.Snapshot().Brand().Value = %#v; want <nil>", state.Brand().Value)
	}
	if state.IsPromise() {
		t.Error("c.Snapshot().IsPromise = true; want false")
	}
	state.Release()
	ans, finish := c.SendCall(ctx, Send{})
	if _, err := ans.Struct(); err == nil {
		t.Error("SendCall did not return error")
	}
	finish()
	if h.calls != 0 {
		t.Errorf("after SendCall, h.calls = %d; want 0", h.calls)
	}
	ret := new(dummyReturner)
	c.RecvCall(ctx, Recv{
		ReleaseArgs: func() {},
		Returner:    ret,
	})
	if !ret.returned {
		t.Error("RecvCall did not return")
	} else if ret.err == nil {
		t.Error("RecvCall did not return error")
	}
	if err := c.Resolve(ctx); err == nil {
		t.Error("Resolve did not return error")
	}
	if h.calls != 0 {
		t.Errorf("after RecvCall, h.calls = %d; want 0", h.calls)
	}

	// Double release
	c.Release()
	if h.shutdowns > 1 {
		t.Error("second Release made more calls to ClientHook.Shutdown")
	}
}
func TestResolve(t *testing.T) {
	test := func(t *testing.T, name string, f func(t *testing.T, p1, p2 Client, r1, r2 Resolver[Client])) {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			p1, r1 := NewLocalPromise[Client]()
			p2, r2 := NewLocalPromise[Client]()
			defer p1.Release()
			defer p2.Release()
			f(t, p1, p2, r1, r2)
		})
	}
	t.Run("Clients", func(t *testing.T) {
		test(t, "Waits for the full chain", func(t *testing.T, p1, p2 Client, r1, r2 Resolver[Client]) {
			r1.Fulfill(p2)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second/10)
			defer cancel()
			require.NotNil(t, p1.Resolve(ctx), "blocks on second promise")
			r2.Fulfill(Client{})
			require.NoError(t, p1.Resolve(context.Background()), "resolves after second resolution")
			assert.True(t, p1.IsSame(Client{}), "p1 resolves to null")
			assert.True(t, p2.IsSame(Client{}), "p2 resolves to null")
			assert.True(t, p1.IsSame(p2), "p1 & p2 are the same")
		})
	})
	t.Run("Snapshots", func(t *testing.T) {
		test(t, "Resolve1 only waits for one link", func(t *testing.T, p1, p2 Client, r1, r2 Resolver[Client]) {
			s1 := p1.Snapshot()
			defer s1.Release()
			r1.Fulfill(p2)
			require.NoError(t, s1.Resolve1(context.Background()), "Resolve1 returns after first resolution")
		})
		test(t, "Resolve waits for the full chain", func(t *testing.T, p1, p2 Client, r1, r2 Resolver[Client]) {
			s1 := p1.Snapshot()
			defer s1.Release()
			r1.Fulfill(p2)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second/10)
			defer cancel()
			require.NotNil(t, s1.Resolve(ctx), "blocks on second promise")
			r2.Fulfill(Client{})
			require.NoError(t, s1.Resolve(context.Background()), "resolves after second resolution")
		})
	})
}

func TestNullClient(t *testing.T) {
	ctx := context.Background()
	c, p := NewPromisedClient(new(dummyHook))
	p.Fulfill(Client{})
	tests := []struct {
		name string
		c    Client
	}{
		{"nil", Client{}},
		{"promised nil", c},
	}

	for _, test := range tests {
		c := test.c
		t.Run(test.name, func(t *testing.T) {
			if (NewClient(nil) != Client{}) {
				t.Error("NewClient(nil) != nil")
			}
			if !c.IsSame(c) {
				t.Error("!<nil>.IsSame(<nil>)")
			}
			if c.IsValid() {
				t.Error("null client is valid")
			}
			state := c.Snapshot()
			if state.Brand().Value != nil {
				t.Errorf("c.Snapshot().Brand() = %#v; want <nil>", state.Brand())
			}
			if state.IsPromise() {
				t.Error("c.Snapshot().IsPromise = true; want false")
			}
			state.Release()
			ans, finish := c.SendCall(ctx, Send{})
			if _, err := ans.Struct(); err == nil {
				t.Error("SendCall did not return error")
			}
			finish()
			ret := new(dummyReturner)
			c.RecvCall(ctx, Recv{
				ReleaseArgs: func() {},
				Returner:    ret,
			})
			if !ret.returned {
				t.Error("RecvCall did not return")
			} else if ret.err == nil {
				t.Error("RecvCall did not return error")
			}
			rctx, cancel := context.WithTimeout(ctx, 1*time.Second)
			if err := c.Resolve(rctx); err != nil {
				t.Error("Resolve failed:", err)
			}
			cancel()
			c.Release()
			c.Release() // should not panic
		})
	}
}

func TestPromisedClient(t *testing.T) {
	a := &dummyHook{brand: Brand{Value: int(111)}}
	b := &dummyHook{brand: Brand{Value: int(222)}}
	ca, pa := NewPromisedClient(a)
	defer ca.Release()
	cb := NewClient(b)
	defer cb.Release()
	ctx := context.Background()

	if ca.IsSame(cb) {
		t.Error("before resolution, ca == cb")
	}
	state := ca.Snapshot()
	if state.Brand().Value != int(111) {
		t.Errorf("before resolution, ca.Snapshot().Brand().Value = %#v; want 111", state.Brand().Value)
	}
	if !state.IsPromise() {
		t.Error("before resolution, ca.Snapshot().IsPromise = false; want true")
	}
	state.Release()
	_, finish := ca.SendCall(ctx, Send{})
	finish()
	pa.Fulfill(cb)
	if a.shutdowns == 0 {
		t.Error("a not shut down after fulfilling ClientPromise")
	} else if a.shutdowns > 1 {
		t.Error("a shut down multiple times")
	}
	_, finish = ca.SendCall(ctx, Send{})
	finish()

	if !ca.IsSame(cb) {
		t.Errorf("after resolution, ca != cb (%v vs. %v)", ca, cb)
	}
	state = ca.Snapshot()
	if state.Brand().Value != int(222) {
		t.Errorf("after resolution, ca.Snapshot().Brand().Value = %#v; want 222", state.Brand().Value)
	}
	if state.IsPromise() {
		t.Error("after resolution, ca.Snapshot().IsPromise = true; want false")
	}
	state.Release()

	if b.shutdowns > 0 {
		t.Error("b shut down before clients released")
	}
	ca.Release()
	if b.shutdowns > 0 {
		t.Error("b shut down after ca.Release but before cb.Release")
	}
	cb.Release()
	if b.shutdowns == 0 {
		t.Error("b not shut down after calling ca.Release and cb.Release")
	} else if b.shutdowns > 1 {
		t.Error("b shut down multiple times")
	}
}

func TestPromisedClient_EarlyClose(t *testing.T) {
	a := new(dummyHook)
	b := new(dummyHook)
	ca, p := NewPromisedClient(a)
	defer ca.Release()
	cb := NewClient(b)
	defer cb.Release()
	ctx := context.Background()

	ca.Release()
	if a.shutdowns == 0 {
		t.Error("a not shut down after releasing only reference")
	} else if a.shutdowns > 1 {
		t.Error("a shut down multiple times")
	}
	p.Fulfill(cb)
	_, finish := ca.SendCall(ctx, Send{})
	finish()
	if a.calls > 0 {
		t.Error("a called after shut down")
	}
	if b.calls > 0 {
		t.Error("b called after shut down")
	}
	if b.shutdowns > 0 {
		t.Error("b shut down after Fulfill")
	}
	cb.Release()
	if b.shutdowns == 0 {
		t.Error("b not shut down after releasing only reference")
	} else if b.shutdowns > 1 {
		t.Error("b shut down multiple times")
	}
}

type dummyHook struct {
	calls     int
	brand     Brand
	shutdowns int
}

func (dh *dummyHook) String() string {
	return fmt.Sprintf(
		"&dummyHook{calls: %v, brand: %v, shutdowns: %v}",
		dh.calls, dh.brand, dh.shutdowns,
	)
}

func (dh *dummyHook) Send(_ context.Context, s Send) (*Answer, ReleaseFunc) {
	dh.calls++
	return ImmediateAnswer(s.Method, newEmptyStruct().ToPtr()), func() {}
}

func (dh *dummyHook) Recv(_ context.Context, r Recv) PipelineCaller {
	dh.calls++
	r.AllocResults(ObjectSize{})
	r.Return()
	return nil
}

func (dh *dummyHook) Brand() Brand {
	return dh.brand
}

func (dh *dummyHook) Shutdown() {
	dh.shutdowns++
}

type dummyReturner struct {
	s        Struct
	returned bool
	err      error
}

func (dr *dummyReturner) AllocResults(sz ObjectSize) (Struct, error) {
	if dr.s.IsValid() {
		return Struct{}, errors.New("AllocResults called multiple times")
	}
	_, seg, err := NewMessage(SingleSegment(nil))
	if err != nil {
		return Struct{}, err
	}
	dr.s, err = NewRootStruct(seg, sz)
	return dr.s, err
}

func (dr *dummyReturner) PrepareReturn(e error) {
	dr.err = e
}

func (dr *dummyReturner) Return() {
	dr.returned = true
}

func (dr *dummyReturner) ReleaseResults() {
}

func TestToInterface(t *testing.T) {
	_, seg, err := NewMessage(SingleSegment(nil))
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		ptr Ptr
		in  Interface
	}{
		{Ptr{}, Interface{}},
		{Struct{}.ToPtr(), Interface{}},
		{Struct{seg: seg, off: 0, depthLimit: maxDepth}.ToPtr(), Interface{}},
		{Interface{}.ToPtr(), Interface{}},
		{Interface{seg, 42}.ToPtr(), Interface{seg, 42}},
	}
	for _, test := range tests {
		if in := test.ptr.Interface(); in != test.in {
			t.Errorf("ToInterface(%#v) = %#v; want %#v", test.ptr, in, test.in)
		}
	}
}

func TestInterface_value(t *testing.T) {
	_, seg, err := NewMessage(SingleSegment(nil))
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		in  Interface
		val rawPointer
	}{
		{Interface{}, 0},
		{NewInterface(seg, 0), 0x0000000000000003},
		{NewInterface(seg, 0xdeadbeef), 0xdeadbeef00000003},
	}
	for _, test := range tests {
		for paddr := address(0); paddr < 16; paddr++ {
			if val := test.in.value(paddr); val != test.val {
				t.Errorf("Interface{seg: %p, cap: %d}.value(%v) = %v; want %v", test.in.seg, test.in.cap, paddr, val, test.val)
			}
		}
	}
}

func TestTransform(t *testing.T) {
	_, s, err := NewMessage(SingleSegment(nil))
	if err != nil {
		t.Fatal(err)
	}
	root, err := NewStruct(s, ObjectSize{PointerCount: 2})
	if err != nil {
		t.Fatal(err)
	}
	a, err := NewStruct(s, ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		t.Fatal(err)
	}
	root.SetPtr(1, a.ToPtr())
	a.SetUint64(0, 1)
	b, err := NewStruct(s, ObjectSize{DataSize: 8})
	if err != nil {
		t.Fatal(err)
	}
	b.SetUint64(0, 2)
	a.SetPtr(0, b.ToPtr())

	dmsg, d, err := NewMessage(SingleSegment(nil))
	if err != nil {
		t.Fatal(err)
	}
	da, err := NewStruct(d, ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		t.Fatal(err)
	}
	if err := dmsg.SetRoot(da.ToPtr()); err != nil {
		t.Fatal(err)
	}
	da.SetUint64(0, 56)
	db, err := NewStruct(d, ObjectSize{DataSize: 8})
	if err != nil {
		t.Fatal(err)
	}
	db.SetUint64(0, 78)
	da.SetPtr(0, db.ToPtr())

	tests := []struct {
		p         Ptr
		transform []PipelineOp
		out       Ptr
	}{
		{
			root.ToPtr(),
			nil,
			root.ToPtr(),
		},
		{
			root.ToPtr(),
			[]PipelineOp{},
			root.ToPtr(),
		},
		{
			root.ToPtr(),
			[]PipelineOp{
				{Field: 0},
			},
			Ptr{},
		},
		{
			root.ToPtr(),
			[]PipelineOp{
				{Field: 0, DefaultValue: mustMarshal(t, dmsg)},
			},
			da.ToPtr(),
		},
		{
			root.ToPtr(),
			[]PipelineOp{
				{Field: 1},
			},
			a.ToPtr(),
		},
		{
			root.ToPtr(),
			[]PipelineOp{
				{Field: 1, DefaultValue: mustMarshal(t, dmsg)},
			},
			a.ToPtr(),
		},
		{
			root.ToPtr(),
			[]PipelineOp{
				{Field: 1},
				{Field: 0},
			},
			b.ToPtr(),
		},
		{
			root.ToPtr(),
			[]PipelineOp{
				{Field: 0},
				{Field: 0},
			},
			Ptr{},
		},
		{
			root.ToPtr(),
			[]PipelineOp{
				{Field: 0, DefaultValue: mustMarshal(t, dmsg)},
				{Field: 0},
			},
			db.ToPtr(),
		},
		{
			root.ToPtr(),
			[]PipelineOp{
				{Field: 0},
				{Field: 0, DefaultValue: mustMarshal(t, dmsg)},
			},
			da.ToPtr(),
		},
		{
			root.ToPtr(),
			[]PipelineOp{
				{Field: 0, DefaultValue: mustMarshal(t, dmsg)},
				{Field: 1, DefaultValue: mustMarshal(t, dmsg)},
			},
			da.ToPtr(),
		},
	}

	for _, test := range tests {
		out, err := Transform(test.p, test.transform)
		if !deepPointerEqual(out, test.out) {
			t.Errorf("Transform(%+v, %v) = %+v; want %+v", test.p, test.transform, out, test.out)
		}
		if err != nil {
			t.Errorf("Transform(%+v, %v) error: %v", test.p, test.transform, err)
		}
	}
}

func TestMethodString(t *testing.T) {
	tests := []struct {
		m *Method
		s string
	}{
		{
			&Method{
				InterfaceID: 0x8e5322c1e9282534,
				MethodID:    1,
			},
			"@0x8e5322c1e9282534.@1",
		},
		{
			&Method{
				InterfaceID:   0x8e5322c1e9282534,
				MethodID:      1,
				InterfaceName: "aircraftlib:Echo",
				MethodName:    "foo",
			},
			"aircraftlib:Echo.foo",
		},
	}
	for _, test := range tests {
		if s := test.m.String(); s != test.s {
			t.Errorf("%#v.String() = %q; want %q", test.m, s, test.s)
		}
	}
}

func TestPipelineOpString(t *testing.T) {
	tests := []struct {
		op PipelineOp
		s  string
	}{
		{
			PipelineOp{Field: 4},
			"get field 4",
		},
		{
			PipelineOp{Field: 4, DefaultValue: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
			"get field 4 with default",
		},
	}
	for _, test := range tests {
		if s := test.op.String(); s != test.s {
			t.Errorf("%#v.String() = %q; want %q", test.op, s, test.s)
		}
	}
}

func TestWeakClient(t *testing.T) {
	h := new(dummyHook)
	c1 := NewClient(h)
	defer c1.Release()
	w := c1.WeakRef()
	c2, ok := w.AddRef()
	defer c2.Release()
	if !ok {
		t.Fatal("AddRef on open client failed")
	}
	if !c1.IsSame(c2) {
		t.Error("c1 != c2")
	}
	c2.Release()
	if h.shutdowns > 0 {
		t.Fatal("Releasing second reference shut down capability")
	}
	c1.Release()
	if h.shutdowns == 0 {
		t.Errorf("Releasing strong reference with weak reference kept capability open")
	}
	if _, ok := w.AddRef(); ok {
		t.Error("w.AddRef() after shut down did not fail")
	}
}

func TestWeakPromisedClient(t *testing.T) {
	a := new(dummyHook)
	b := new(dummyHook)
	ca, pa := NewPromisedClient(a)
	cb := NewClient(b)
	wa := ca.WeakRef()
	ctx := context.Background()

	_, finish := ca.SendCall(ctx, Send{})
	finish()
	pa.Fulfill(cb)
	_, finish = ca.SendCall(ctx, Send{})
	finish()

	ca.Release()
	defer ca.Release()
	cb2, ok := wa.AddRef()
	defer cb2.Release()
	assert.False(t, ok, "wa.AddRef() failed after releasing ca")
	assert.False(t, cb.IsSame(cb2), "cb != cb2")

	cb.Release()
	if b.shutdowns == 0 {
		t.Error("b not shut down after cb.Release")
	}
}

func mustMarshal(t *testing.T, msg *Message) []byte {
	data, err := msg.Marshal()
	if err != nil {
		t.Fatal("Marshal:", err)
	}
	return data
}

func deepPointerEqual(a, b Ptr) bool {
	if !a.IsValid() && !b.IsValid() {
		return true
	}
	if !a.IsValid() || !b.IsValid() {
		return false
	}
	msgA, _, _ := NewMessage(SingleSegment(nil))
	msgA.SetRoot(a)
	abytes, _ := msgA.Marshal()
	msgB, _, _ := NewMessage(SingleSegment(nil))
	msgB.SetRoot(b)
	bbytes, _ := msgB.Marshal()
	return bytes.Equal(abytes, bbytes)
}

func newEmptyStruct() Struct {
	_, seg, err := NewMessage(SingleSegment(nil))
	if err != nil {
		panic(err)
	}
	s, err := NewRootStruct(seg, ObjectSize{})
	if err != nil {
		panic(err)
	}
	return s
}
