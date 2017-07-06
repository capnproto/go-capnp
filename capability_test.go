package capnp

import (
	"bytes"
	"context"
	"errors"
	"testing"
)

func TestResolveClient(t *testing.T) {
	a := new(dummyHook)
	b := new(dummyHook)
	ca, pa := NewPromisedClient(a)
	cb := NewClient(b)
	ctx := context.Background()

	if ca.IsSame(cb) {
		t.Error("before resolution, ca == cb")
	}
	args := SendArgs{
		Place: func(Struct) error {
			return nil
		},
	}
	_, finish := ca.SendCall(ctx, Method{}, args, CallOptions{})
	finish()
	pa.Resolve(cb)
	_, finish = ca.SendCall(ctx, Method{}, args, CallOptions{})
	finish()

	if !ca.IsSame(cb) {
		t.Error("after resolution, ca != cb")
	}
	if b.closed {
		t.Error("b.closed = true before closing")
	}
	if err := ca.Close(); err != nil {
		t.Error("ca.Close() =", err)
	}
	if b.closed {
		t.Error("b.closed = true after closing ca, but not cb")
	}
	if err := cb.Close(); err != nil {
		t.Error("cb.Close() =", err)
	}
	if !b.closed {
		t.Error("b.closed = false after closing")
	}
}

type dummyHook struct {
	closed bool
}

func (dh *dummyHook) Send(context.Context, Method, SendArgs, CallOptions) (*Answer, ReleaseFunc) {
	return ErrorAnswer(errors.New("dummy hook answer")), func() {}
}

func (dh *dummyHook) Recv(_ context.Context, _ Method, a RecvArgs, _ CallOptions) (*Answer, ReleaseFunc) {
	a.Release()
	return ErrorAnswer(errors.New("dummy hook answer")), func() {}
}

func (dh *dummyHook) Brand() interface{} {
	return nil
}

func (dh *dummyHook) Close() error {
	dh.closed = true
	return nil
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
	if h.closed {
		t.Fatal("Closing second reference closed capability")
	}
	if err := c1.Close(); err != nil {
		t.Errorf("w.AddRef().Close(): %v", err)
	}
	if !h.closed {
		t.Errorf("Closing strong reference with weak reference kept capability open")
	}
	if _, ok := w.AddRef(); ok {
		t.Error("w.AddRef() after close did not fail")
	}
}

func TestResolveWeakClient(t *testing.T) {
	a := new(dummyHook)
	b := new(dummyHook)
	ca, pa := NewPromisedClient(a)
	cb := NewClient(b)
	wa := ca.WeakRef()
	ctx := context.Background()

	args := SendArgs{
		Place: func(Struct) error {
			return nil
		},
	}
	_, finish := ca.SendCall(ctx, Method{}, args, CallOptions{})
	finish()
	pa.Resolve(cb)
	_, finish = ca.SendCall(ctx, Method{}, args, CallOptions{})
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
	if b.closed {
		t.Error("b.closed = true before closing cb2")
	}
	if err := cb2.Close(); err != nil {
		t.Error("cb2.Close() =", err)
	}
	if !b.closed {
		t.Error("b.closed = false after closing cb2")
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
