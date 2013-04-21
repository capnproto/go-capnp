package msgs

import (
	C "github.com/jmckaskill/go-capnproto"
	M "math"
	"reflect"
	"unsafe"
)

var (
	_ = M.Float32bits
	_ = reflect.SliceHeader{}
	_ = unsafe.Pointer(nil)
)

type Pointer struct {
	Ptr C.Pointer
}

func NewPointer(new C.NewFunc) (Pointer, error) {
	ptr, err := C.NewStruct(new, 0, 0)
	return Pointer{Ptr: ptr}, err
}

type Message struct {
	Ptr C.Pointer
}

func NewMessage(new C.NewFunc) (Message, error) {
	ptr, err := C.NewStruct(new, 1, 3)
	return Message{Ptr: ptr}, err
}
func (p Message) Cookie() Pointer {
	ptr := C.ReadPtr(p.Ptr, 0)
	return Pointer{Ptr: ptr}
}
func (p Message) SetCookie(v Pointer) error {
	return p.Ptr.WritePtrs(0, []C.Pointer{v.Ptr})
}
func (p Message) Method() uint16 {
	return C.ReadUInt16(p.Ptr, 0)
}
func (p Message) SetMethod(v uint16) error {
	return C.WriteUInt16(p.Ptr, 0, v)
}
func (p Message) ReturnCookie() Pointer {
	ptr := C.ReadPtr(p.Ptr, 1)
	return Pointer{Ptr: ptr}
}
func (p Message) SetReturnCookie(v Pointer) error {
	return p.Ptr.WritePtrs(1, []C.Pointer{v.Ptr})
}
func (p Message) ReturnMethod() uint16 {
	return C.ReadUInt16(p.Ptr, 2)
}
func (p Message) SetReturnMethod(v uint16) error {
	return C.WriteUInt16(p.Ptr, 2, v)
}
func (p Message) Arguments() Pointer {
	ptr := C.ReadPtr(p.Ptr, 2)
	return Pointer{Ptr: ptr}
}
func (p Message) SetArguments(v Pointer) error {
	return p.Ptr.WritePtrs(2, []C.Pointer{v.Ptr})
}
