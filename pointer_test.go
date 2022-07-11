package capnp

import (
	"errors"
	"testing"
)

func TestEqual(t *testing.T) {
	msg, seg, _ := NewMessage(SingleSegment(nil))
	emptyStruct1, _ := NewStruct(seg, ObjectSize{})
	emptyStruct2, _ := NewStruct(seg, ObjectSize{})
	zeroStruct1, _ := NewStruct(seg, ObjectSize{DataSize: 8, PointerCount: 1})
	zeroStruct2, _ := NewStruct(seg, ObjectSize{DataSize: 8, PointerCount: 1})
	structA1, _ := NewStruct(seg, ObjectSize{DataSize: 16, PointerCount: 1})
	structA1.SetUint32(0, 0xdeadbeef)
	subA1, _ := NewStruct(seg, ObjectSize{DataSize: 8})
	subA1.SetUint32(0, 0x0cafefe0)
	structA1.SetPtr(0, subA1.ToPtr())
	structA2, _ := NewStruct(seg, ObjectSize{DataSize: 16, PointerCount: 1})
	structA2.SetUint32(0, 0xdeadbeef)
	subA2, _ := NewStruct(seg, ObjectSize{DataSize: 8})
	subA2.SetUint32(0, 0x0cafefe0)
	structA2.SetPtr(0, subA2.ToPtr())
	structB, _ := NewStruct(seg, ObjectSize{DataSize: 8, PointerCount: 2})
	structB.SetUint32(0, 0xdeadbeef)
	subB, _ := NewStruct(seg, ObjectSize{DataSize: 8})
	subB.SetUint32(0, 0x0cafefe0)
	structB.SetPtr(0, subB.ToPtr())
	structB.SetPtr(1, emptyStruct1.ToPtr())
	structC, _ := NewStruct(seg, ObjectSize{DataSize: 16, PointerCount: 1})
	structC.SetUint32(0, 0xfeed1234)
	subC, _ := NewStruct(seg, ObjectSize{DataSize: 8})
	subC.SetUint32(0, 0x0cafefe0)
	structC.SetPtr(0, subA2.ToPtr())
	structD, _ := NewStruct(seg, ObjectSize{DataSize: 16, PointerCount: 1})
	structD.SetUint32(0, 0xdeadbeef)
	subD, _ := NewStruct(seg, ObjectSize{DataSize: 8})
	subD.SetUint32(0, 0x12345678)
	structD.SetPtr(0, subD.ToPtr())
	emptyStructList1, _ := NewCompositeList(seg, ObjectSize{DataSize: 8, PointerCount: 1}, 0)
	emptyStructList2, _ := NewCompositeList(seg, ObjectSize{DataSize: 8, PointerCount: 1}, 0)
	emptyInt32List, _ := NewInt32List(seg, 0)
	emptyFloat32List, _ := NewFloat32List(seg, 0)
	emptyFloat64List, _ := NewFloat64List(seg, 0)
	list123Int, _ := NewInt32List(seg, 3)
	list123Int.Set(0, 1)
	list123Int.Set(1, 2)
	list123Int.Set(2, 3)
	list12Int, _ := NewInt32List(seg, 2)
	list12Int.Set(0, 1)
	list12Int.Set(1, 2)
	list456Int, _ := NewInt32List(seg, 3)
	list456Int.Set(0, 4)
	list456Int.Set(1, 5)
	list456Int.Set(2, 6)
	list123Struct, _ := NewCompositeList(seg, ObjectSize{DataSize: 8}, 3)
	list123Struct.Struct(0).SetUint32(0, 1)
	list123Struct.Struct(1).SetUint32(0, 2)
	list123Struct.Struct(2).SetUint32(0, 3)
	list12Struct, _ := NewCompositeList(seg, ObjectSize{DataSize: 8}, 2)
	list12Struct.Struct(0).SetUint32(0, 1)
	list12Struct.Struct(1).SetUint32(0, 2)
	list456Struct, _ := NewCompositeList(seg, ObjectSize{DataSize: 8}, 3)
	list456Struct.Struct(0).SetUint32(0, 4)
	list456Struct.Struct(1).SetUint32(0, 5)
	list456Struct.Struct(2).SetUint32(0, 6)
	plistA1, _ := NewPointerList(seg, 1)
	plistA1.Set(0, structA1.ToPtr())
	plistA2, _ := NewPointerList(seg, 1)
	plistA2.Set(0, structA2.ToPtr())
	plistB, _ := NewPointerList(seg, 1)
	plistB.Set(0, structB.ToPtr())
	ec := ErrorClient(errors.New("boo"))
	msg.CapTable = []Client{
		0: ec,
		1: ec,
		2: ErrorClient(errors.New("another boo")),
		3: Client{},
		4: Client{},
	}
	iface1 := NewInterface(seg, 0)
	iface2 := NewInterface(seg, 1)
	ifaceAlt := NewInterface(seg, 2)
	ifaceMissing1 := NewInterface(seg, 3)
	ifaceMissing2 := NewInterface(seg, 4)
	ifaceOOB1 := NewInterface(seg, 5)
	ifaceOOB2 := NewInterface(seg, 6)

	tests := []struct {
		name   string
		p1, p2 Ptr
		equal  bool
	}{

		// Structs
		{"EmptyStruct_EmptyStruct", emptyStruct1.ToPtr(), emptyStruct2.ToPtr(), true},
		{"EmptyStruct_ZeroStruct", emptyStruct1.ToPtr(), zeroStruct2.ToPtr(), true},
		{"ZeroStruct_ZeroStruct", zeroStruct1.ToPtr(), zeroStruct2.ToPtr(), true},
		{"EmptyStruct_StructA", emptyStruct1.ToPtr(), structA1.ToPtr(), false},
		{"StructA_EmptyStruct", structA1.ToPtr(), emptyStruct1.ToPtr(), false},
		{"StructA1_StructA1", structA1.ToPtr(), structA1.ToPtr(), true},
		{"StructA1_StructA2", structA1.ToPtr(), structA2.ToPtr(), true},
		{"StructA2_StructA1", structA2.ToPtr(), structA1.ToPtr(), true},
		{"StructA_StructB", structA1.ToPtr(), structB.ToPtr(), false},
		{"StructB_StructA", structB.ToPtr(), structA1.ToPtr(), false},
		{"StructA_StructC", structA1.ToPtr(), structC.ToPtr(), false},
		{"StructC_StructA", structC.ToPtr(), structA1.ToPtr(), false},
		{"StructA_StructD", structA1.ToPtr(), structD.ToPtr(), false},
		{"StructD_StructA", structD.ToPtr(), structA1.ToPtr(), false},

		// Lists
		{"EmptyStructList_EmptyStructList", emptyStructList1.ToPtr(), emptyStructList2.ToPtr(), true},
		{"EmptyInt32List_EmptyFloat32List", emptyInt32List.ToPtr(), emptyFloat32List.ToPtr(), true}, // identical on wire
		{"EmptyInt32List_EmptyFloat64List", emptyInt32List.ToPtr(), emptyFloat64List.ToPtr(), false},
		{"List123Int_List456Int", list123Int.ToPtr(), list456Int.ToPtr(), false},
		{"List123Struct_List456Struct", list123Struct.ToPtr(), list456Struct.ToPtr(), false},
		{"List123Int_List123Struct", list123Int.ToPtr(), list123Struct.ToPtr(), true},
		{"List123Struct_List123Int", list123Struct.ToPtr(), list123Int.ToPtr(), true},
		{"List123Int_List12Int", list123Int.ToPtr(), list12Int.ToPtr(), false},
		{"List123Struct_List12Struct", list123Struct.ToPtr(), list12Struct.ToPtr(), false},
		{"PointerListA1_PointerListA2", plistA1.ToPtr(), plistA2.ToPtr(), true},
		{"PointerListA2_PointerListA1", plistA2.ToPtr(), plistA1.ToPtr(), true},
		{"PointerListA_PointerListB", plistA1.ToPtr(), plistB.ToPtr(), false},
		{"PointerListB_PointerListA", plistB.ToPtr(), plistA1.ToPtr(), false},

		// Interfaces
		{"InterfaceA1_InterfaceA1", iface1.ToPtr(), iface1.ToPtr(), true},
		{"InterfaceA1_InterfaceA2", iface1.ToPtr(), iface2.ToPtr(), true},
		{"InterfaceA_InterfaceB", iface1.ToPtr(), ifaceAlt.ToPtr(), false},
		{"InterfaceMissingCap_Null", ifaceMissing1.ToPtr(), Ptr{}, false},
		{"InterfaceMissingCap_InterfaceMissingCap", ifaceMissing1.ToPtr(), ifaceMissing2.ToPtr(), true},
		{"InterfaceOOB1_InterfaceOOB1", ifaceOOB1.ToPtr(), ifaceOOB1.ToPtr(), true},
		{"InterfaceOOB1_InterfaceOOB2", ifaceOOB1.ToPtr(), ifaceOOB2.ToPtr(), false},
		{"InterfaceOOB_InterfaceMissingCap", ifaceOOB1.ToPtr(), ifaceMissing1.ToPtr(), false},

		// Null
		{"Null_Null", Ptr{}, Ptr{}, true},
		{"EmptyStruct_Null", emptyStruct1.ToPtr(), Ptr{}, false},
		{"Null_EmptyStruct", Ptr{}, emptyStruct1.ToPtr(), false},
		{"Null_EmptyStructList", Ptr{}, emptyStructList1.ToPtr(), false},
		{"EmptyStructList_Null", emptyStructList1.ToPtr(), Ptr{}, false},
		{"Interface_Null", iface1.ToPtr(), Ptr{}, false},
		{"Null_Interface", Ptr{}, iface1.ToPtr(), false},

		// Misc combinations that shouldn't be equal.
		{"EmptyStruct_EmptyList", emptyStruct1.ToPtr(), emptyStructList1.ToPtr(), false},
		{"EmptyStruct_InterfaceA", emptyStruct1.ToPtr(), iface1.ToPtr(), false},
		{"EmptyList_InterfaceA", emptyStructList1.ToPtr(), iface1.ToPtr(), false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Equal(test.p1, test.p2)
			if err != nil {
				t.Fatal(err)
			}
			if got != test.equal {
				if got {
					t.Error("p1 equals p2; want not equal")
				} else {
					t.Error("p1 does not equal p2; want equal")
				}
			}
		})
	}
}
