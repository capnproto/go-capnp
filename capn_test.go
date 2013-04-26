package capn

import (
	"testing"
)

type testPointer PointerType

func newTestPointer(p PointerType) (Pointer, error)                  { return testPointer(p), nil }
func (t testPointer) New(p PointerType) (Pointer, error)             { return newTestPointer(p) }
func (t testPointer) Call(method int, args Pointer) (Pointer, error) { panic("todo") }
func (t testPointer) Type() PointerType                              { return PointerType(t) }
func (t testPointer) Read(off int, v []uint8) error                  { return nil }
func (t testPointer) Write(off int, v []uint8) error                 { return nil }
func (t testPointer) ReadPtrs(off int, v []Pointer) error            { return nil }
func (t testPointer) WritePtrs(off int, v []Pointer) error           { return nil }

func TestNewStruct(t *testing.T) {
	p := Must(NewStruct(newTestPointer, 8, 4))
	typ := p.Type()
	if typ.Type() != Struct {
		t.Fatalf("expected struct got %v", typ.Type())
	}
	if typ.DataSize() != 8 {
		t.Fatalf("expected struct data size 8 got %d", typ.DataSize())
	}
	if typ.PointerNum() != 4 {
		t.Fatalf("expected struct pointer num 4 got %d", typ.PointerNum())
	}
	if typ.Offset() != 0 {
		t.Fatalf("expected initial offset of 0 got %d", typ.Offset())
	}
	typ.SetOffset(0x1FFFFFFF)
	if typ.Offset() != 0x1FFFFFFF {
		t.Fatalf("expected set offset of 0x1FFFFFFF got %d", typ.Offset())
	}
	typ.SetOffset(-0x20000000)
	if typ.Offset() != -0x20000000 {
		t.Fatalf("expected set offset of -0x20000000 got %d", typ.Offset())
	}
}
