package capnp

import (
	"bytes"
	"context"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	ctx := context.Background()
	h := &dummyHook{brand: int(42)}
	c := NewClient(h)

	if !c.IsSame(c) {
		t.Error("!c.IsSame(c)")
	}
	if !c.IsValid() {
		t.Error("new client is not valid")
	}
	brand := c.Brand()
	if x, ok := brand.(int); !ok || x != 42 {
		t.Errorf("c.Brand() = %v; want 42", brand)
	}
	ans, finish := c.SendCall(ctx, Send{})
	if _, err := ans.Struct(); err != nil {
		t.Error("SendCall:", err)
	}
	finish()
	if h.calls != 1 {
		t.Errorf("after SendCall, h.calls = %d; want 1", h.calls)
	}
	ans, finish = c.RecvCall(ctx, Recv{ReleaseArgs: func() {}})
	if _, err := ans.Struct(); err != nil {
		t.Error("RecvCall:", err)
	}
	finish()
	if h.calls != 2 {
		t.Errorf("after RecvCall, h.calls = %d; want 2", h.calls)
	}
	rctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	if err := c.Resolve(rctx); err != nil {
		t.Error("Resolve failed:", err)
	}
	cancel()
	if err := c.Close(); err != nil {
		t.Error("Close:", err)
	}
	if h.closes == 0 {
		t.Error("Close did not call ClientHook.Close")
	} else if h.closes > 1 {
		t.Error("Close called ClientHook.Close multiple times")
	}
}

func TestClosedClient(t *testing.T) {
	ctx := context.Background()
	h := &dummyHook{brand: int(42)}
	c := NewClient(h)
	if err := c.Close(); err != nil {
		t.Error("Close:", err)
	}

	if c.IsValid() {
		t.Error("closed client is valid")
	}
	ans, finish := c.SendCall(ctx, Send{})
	if _, err := ans.Struct(); err == nil {
		t.Error("SendCall did not return error")
	}
	finish()
	if h.calls != 0 {
		t.Errorf("after SendCall, h.calls = %d; want 0", h.calls)
	}
	ans, finish = c.RecvCall(ctx, Recv{ReleaseArgs: func() {}})
	if _, err := ans.Struct(); err == nil {
		t.Error("RecvCall did not return error")
	}
	finish()
	if err := c.Resolve(ctx); err == nil {
		t.Error("Resolve did not return error")
	}
	if h.calls != 0 {
		t.Errorf("after RecvCall, h.calls = %d; want 0", h.calls)
	}
	if err := c.Close(); err == nil {
		t.Error("second Close did not return error")
	}
	if h.closes > 1 {
		t.Error("second Close made more calls to ClientHook.Close")
	}
}

func TestNullClient(t *testing.T) {
	ctx := context.Background()
	c, p := NewPromisedClient(new(dummyHook))
	p.Fulfill(nil)
	tests := []struct {
		name string
		c    *Client
	}{
		{"nil", nil},
		{"promised nil", c},
	}

	for _, test := range tests {
		c := test.c
		t.Run(test.name, func(t *testing.T) {
			if NewClient(nil) != nil {
				t.Error("NewClient(nil) != nil")
			}
			if !c.IsSame(c) {
				t.Error("!<nil>.IsSame(<nil>)")
			}
			if c.IsValid() {
				t.Error("null client is valid")
			}
			if b := c.Brand(); b != nil {
				t.Errorf("c.Brand() = %v; want <nil>", b)
			}
			ans, finish := c.SendCall(ctx, Send{})
			if _, err := ans.Struct(); err == nil {
				t.Error("SendCall did not return error")
			}
			finish()
			ans, finish = c.RecvCall(ctx, Recv{ReleaseArgs: func() {}})
			if _, err := ans.Struct(); err == nil {
				t.Error("RecvCall did not return error")
			}
			finish()
			rctx, cancel := context.WithTimeout(ctx, 1*time.Second)
			if err := c.Resolve(rctx); err != nil {
				t.Error("Resolve failed:", err)
			}
			cancel()
			if err := c.Close(); err != nil {
				t.Error("Close #1:", err)
			}
			if err := c.Close(); err != nil {
				t.Error("Close #2:", err)
			}
		})
	}
}

func TestPromisedClient(t *testing.T) {
	a := new(dummyHook)
	b := new(dummyHook)
	ca, pa := NewPromisedClient(a)
	cb := NewClient(b)
	ctx := context.Background()

	if ca.IsSame(cb) {
		t.Error("before resolution, ca == cb")
	}
	_, finish := ca.SendCall(ctx, Send{})
	finish()
	pa.Fulfill(cb)
	_, finish = ca.SendCall(ctx, Send{})
	finish()

	if !ca.IsSame(cb) {
		t.Error("after resolution, ca != cb")
	}
	if b.closes > 0 {
		t.Error("b closed before clients closed")
	}
	if err := ca.Close(); err != nil {
		t.Error("ca.Close() =", err)
	}
	if b.closes > 0 {
		t.Error("b closed after ca.Close but before cb.Close")
	}
	if err := cb.Close(); err != nil {
		t.Error("cb.Close() =", err)
	}
	if b.closes == 0 {
		t.Error("b not closed after calling ca.Close and cb.Close")
	} else if b.closes > 1 {
		t.Error("b closed multiple times")
	}
}

type dummyHook struct {
	calls    int
	brand    interface{}
	closes   int
	closeErr error
}

func (dh *dummyHook) Send(context.Context, Send) (*Answer, ReleaseFunc) {
	dh.calls++
	return ImmediateAnswer(newEmptyStruct()), func() {}
}

func (dh *dummyHook) Recv(_ context.Context, r Recv) (*Answer, ReleaseFunc) {
	r.ReleaseArgs()
	dh.calls++
	return ImmediateAnswer(newEmptyStruct()), func() {}
}

func (dh *dummyHook) Brand() interface{} {
	return dh.brand
}

func (dh *dummyHook) Close() error {
	dh.closes++
	return dh.closeErr
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
		for paddr := Address(0); paddr < 16; paddr++ {
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
	w := c1.WeakRef()
	c2, ok := w.AddRef()
	if !ok {
		t.Fatal("AddRef on open client failed")
	}
	if !c1.IsSame(c2) {
		t.Error("c1 != c2")
	}
	if err := c2.Close(); err != nil {
		t.Errorf("w.AddRef().Close(): %v", err)
	}
	if h.closes > 0 {
		t.Fatal("Closing second reference closed capability")
	}
	if err := c1.Close(); err != nil {
		t.Errorf("w.AddRef().Close(): %v", err)
	}
	if h.closes == 0 {
		t.Errorf("Closing strong reference with weak reference kept capability open")
	}
	if _, ok := w.AddRef(); ok {
		t.Error("w.AddRef() after close did not fail")
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

	if err := ca.Close(); err != nil {
		t.Error("ca.Close() =", err)
	}
	cb2, ok := wa.AddRef()
	if !ok {
		t.Error("wa.AddRef() failed after closing ca")
	}
	if !cb.IsSame(cb2) {
		t.Error("cb != cb2")
	}
	if err := cb.Close(); err != nil {
		t.Error("cb.Close() =", err)
	}
	if b.closes > 0 {
		t.Error("b closed before cb2.Close")
	}
	if err := cb2.Close(); err != nil {
		t.Error("cb2.Close() =", err)
	}
	if b.closes == 0 {
		t.Error("b not closed after cb2.Close")
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
