package pogs

import (
	"testing"

	"zombiezen.com/go/capnproto2"
	air "zombiezen.com/go/capnproto2/internal/aircraftlib"
)

type VerVal struct {
	Val int16
}

type VerOneData struct {
	VerVal
}

type VerTwoData struct {
	*VerVal
	Duo int64
}

type VerTwoDataOmit struct {
	VerVal `capnp:"-"`
	Duo    int64
}

type F16 struct {
	PlaneBase `capnp:"base"`
}

func TestExtract_Embed(t *testing.T) {
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatalf("NewMessage: %v", err)
	}
	v1, err := air.NewRootVerOneData(seg)
	if err != nil {
		t.Fatalf("NewRootVerOneData: %v", err)
	}
	v1.SetVal(123)
	out := new(VerOneData)
	if err := Extract(out, air.VerOneData_TypeID, v1.Struct); err != nil {
		t.Errorf("Extract error: %v", err)
	}
	if out.Val != 123 {
		t.Errorf("Extract produced %s; want %s", zpretty.Sprint(out), zpretty.Sprint(&VerOneData{VerVal{123}}))
	}
}

func TestExtract_EmbedPtr(t *testing.T) {
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatalf("NewMessage: %v", err)
	}
	v2, err := air.NewRootVerTwoData(seg)
	if err != nil {
		t.Fatalf("NewRootVerTwoData: %v", err)
	}
	v2.SetVal(123)
	v2.SetDuo(456)
	out := new(VerTwoData)
	if err := Extract(out, air.VerTwoData_TypeID, v2.Struct); err != nil {
		t.Errorf("Extract error: %v", err)
	}
	if out.VerVal == nil || out.Val != 123 || out.Duo != 456 {
		t.Errorf("Extract produced %s; want %s", zpretty.Sprint(out), zpretty.Sprint(&VerTwoData{&VerVal{123}, 456}))
	}
}

func TestExtract_EmbedOmit(t *testing.T) {
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatalf("NewMessage: %v", err)
	}
	v2, err := air.NewRootVerTwoData(seg)
	if err != nil {
		t.Fatalf("NewRootVerTwoData: %v", err)
	}
	v2.SetVal(123)
	v2.SetDuo(456)
	out := new(VerTwoDataOmit)
	if err := Extract(out, air.VerTwoData_TypeID, v2.Struct); err != nil {
		t.Errorf("Extract error: %v", err)
	}
	if out.Val != 0 || out.Duo != 456 {
		t.Errorf("Extract produced %s; want %s", zpretty.Sprint(out), zpretty.Sprint(&VerTwoDataOmit{VerVal{}, 456}))
	}
}

func TestExtract_EmbedName(t *testing.T) {
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatalf("NewMessage: %v", err)
	}
	f16, err := air.NewRootF16(seg)
	if err != nil {
		t.Fatalf("NewRootF16: %v", err)
	}
	base, err := f16.NewBase()
	if err != nil {
		t.Fatalf("F16.NewBase: %v", err)
	}
	if err := base.SetName("ALL YOUR BASE"); err != nil {
		t.Fatalf("Planebase.SetName: %v", err)
	}
	base.SetRating(5)
	base.SetCanFly(true)

	out := new(F16)
	if err := Extract(out, air.F16_TypeID, f16.Struct); err != nil {
		t.Errorf("Extract error: %v", err)
	}
	if out.Name != "ALL YOUR BASE" || out.Rating != 5 || !out.CanFly {
		t.Errorf("Extract produced %s; want %s", zpretty.Sprint(out), zpretty.Sprint(&F16{PlaneBase{Name: "ALL YOUR BASE", Rating: 5, CanFly: true}}))
	}
}

func TestInsert_Embed(t *testing.T) {
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatalf("NewMessage: %v", err)
	}
	v1, err := air.NewRootVerOneData(seg)
	if err != nil {
		t.Fatalf("NewRootVerOneData: %v", err)
	}
	gv1 := &VerOneData{VerVal{123}}
	err = Insert(air.VerOneData_TypeID, v1.Struct, gv1)
	if err != nil {
		t.Errorf("Insert(%s) error: %v", zpretty.Sprint(gv1), err)
	}
	if v1.Val() != 123 {
		t.Errorf("Insert(%s) produced %v", zpretty.Sprint(gv1), v1)
	}
}

func TestInsert_EmbedPtr(t *testing.T) {
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatalf("NewMessage: %v", err)
	}
	v2, err := air.NewRootVerTwoData(seg)
	if err != nil {
		t.Fatalf("NewRootVerTwoData: %v", err)
	}
	gv2 := &VerTwoData{&VerVal{123}, 456}
	err = Insert(air.VerTwoData_TypeID, v2.Struct, gv2)
	if err != nil {
		t.Errorf("Insert(%s) error: %v", zpretty.Sprint(gv2), err)
	}
	if v2.Val() != 123 || v2.Duo() != 456 {
		t.Errorf("Insert(%s) produced %v", zpretty.Sprint(gv2), v2)
	}
}

func TestInsert_EmbedNilPtr(t *testing.T) {
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatalf("NewMessage: %v", err)
	}
	v2, err := air.NewRootVerTwoData(seg)
	if err != nil {
		t.Fatalf("NewRootVerTwoData: %v", err)
	}
	gv2 := &VerTwoData{nil, 456}
	err = Insert(air.VerTwoData_TypeID, v2.Struct, gv2)
	if err != nil {
		t.Errorf("Insert(%s) error: %v", zpretty.Sprint(gv2), err)
	}
	if v2.Val() != 0 || v2.Duo() != 456 {
		t.Errorf("Insert(%s) produced %v", zpretty.Sprint(gv2), v2)
	}
}

func TestInsert_EmbedOmit(t *testing.T) {
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatalf("NewMessage: %v", err)
	}
	v2, err := air.NewRootVerTwoData(seg)
	if err != nil {
		t.Fatalf("NewRootVerTwoData: %v", err)
	}
	in := &VerTwoDataOmit{VerVal{123}, 456}
	err = Insert(air.VerTwoData_TypeID, v2.Struct, in)
	if err != nil {
		t.Errorf("Insert(%s) error: %v", zpretty.Sprint(in), err)
	}
	if v2.Val() != 0 || v2.Duo() != 456 {
		t.Errorf("Insert(%s) produced %v", zpretty.Sprint(in), v2)
	}
}

func TestInsert_EmbedNamed(t *testing.T) {
	_, seg, err := capnp.NewMessage(capnp.SingleSegment(nil))
	if err != nil {
		t.Fatalf("NewMessage: %v", err)
	}
	f16, err := air.NewRootF16(seg)
	if err != nil {
		t.Fatalf("NewRootF16: %v", err)
	}
	in := &F16{PlaneBase{Name: "ALL YOUR BASE", Rating: 5, CanFly: true}}
	err = Insert(air.F16_TypeID, f16.Struct, in)
	if err != nil {
		t.Errorf("Insert(%s) error: %v", zpretty.Sprint(in), err)
	}
	base, err := f16.Base()
	if err != nil {
		t.Errorf("f16.base: %v", err)
	}
	name, err := base.Name()
	if err != nil {
		t.Errorf("f16.base.name: %v", err)
	}
	if name != "ALL YOUR BASE" || base.Rating() != 5 || !base.CanFly() {
		t.Errorf("Insert(%s) produced %v", zpretty.Sprint(in), f16)
	}
}
