package capn

import (
	"fmt"
	"testing"
)

func TestList(t *testing.T) {
	s := NewBuffer(nil)
	p := s.NewList(33, 0, 20, true)
	fmt.Println(p)
	p.Write32(0, 86)
	if p.Read32(0) != 86 {
		t.Fatal("Read32")
	}
	p.Write32(-1, 10)
	p.Write32(19, 17)
	p.Write32(20, 15)
	if p.Read32(-1) != 0 || p.Read32(19) != 17 || p.Read32(20) != 0 {
		t.Fatal("Read32 out of bounds")
	}
	p.Write32(0, 2147483648)
	if p.Read32(0) != 2147483648 {
		t.Fatal("Read32 large")
	}
	fmt.Printf("%d\n", s.Data)
}

func TestStruct(t *testing.T) {
	s := NewBuffer(nil)
	p := s.NewStruct(3, 1, true)
	p.WriteStruct1(2, true)
	p.WriteStruct1(3, true)
	if p.ReadStruct1(1) || !p.ReadStruct1(2) || !p.ReadStruct1(3) || p.ReadStruct1(4) {
		t.Fatal("ReadStruct1")
	}
	fmt.Printf("%d\n", s.Data)
}
