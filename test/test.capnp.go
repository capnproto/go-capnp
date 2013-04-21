package test

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

type _TestInterface_testMethod1_args struct {
	Ptr C.Pointer
}

func new_TestInterface_testMethod1_args(new C.NewFunc) (_TestInterface_testMethod1_args, error) {
	ptr, err := new(C.MakeStruct(1, 1))
	return _TestInterface_testMethod1_args{Ptr: ptr}, err
}
func (p _TestInterface_testMethod1_args) v() bool {
	return (C.ReadUInt8(p.Ptr, 0) & 0) != 0
}
func (p _TestInterface_testMethod1_args) setV(v bool) error {
	return C.WriteBool(p.Ptr, 0, v)
}
func (p _TestInterface_testMethod1_args) arg1() string {
	ptr := C.ReadPtr(p.Ptr, 0)
	return C.ToString(ptr, "")
}
func (p _TestInterface_testMethod1_args) setArg1(v string) error {
	data, err := C.NewString(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(0, []C.Pointer{data})
}
func (p _TestInterface_testMethod1_args) arg2() uint16 {
	return C.ReadUInt16(p.Ptr, 2)
}
func (p _TestInterface_testMethod1_args) setArg2(v uint16) error {
	return C.WriteUInt16(p.Ptr, 2, v)
}

type _TestInterface_testMethod0_args struct {
	Ptr C.Pointer
}

func new_TestInterface_testMethod0_args(new C.NewFunc) (_TestInterface_testMethod0_args, error) {
	ptr, err := new(C.MakeStruct(0, 1))
	return _TestInterface_testMethod0_args{Ptr: ptr}, err
}
func (p _TestInterface_testMethod0_args) arg0() TestInterface {
	ptr := C.ReadPtr(p.Ptr, 0)
	return _TestInterface_remote{Ptr: ptr}
}
func (p _TestInterface_testMethod0_args) setArg0(v TestInterface) error {
	data, err := v.MarshalCaptain(p.Ptr.New)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(0, []C.Pointer{data})
}

type _TestInterface_testMethod2_args struct {
	Ptr C.Pointer
}

func new_TestInterface_testMethod2_args(new C.NewFunc) (_TestInterface_testMethod2_args, error) {
	ptr, err := new(C.MakeStruct(0, 1))
	return _TestInterface_testMethod2_args{Ptr: ptr}, err
}
func (p _TestInterface_testMethod2_args) arg0() TestInterface {
	ptr := C.ReadPtr(p.Ptr, 0)
	return _TestInterface_remote{Ptr: ptr}
}
func (p _TestInterface_testMethod2_args) setArg0(v TestInterface) error {
	data, err := v.MarshalCaptain(p.Ptr.New)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(0, []C.Pointer{data})
}

type TestInterface interface {
	C.Marshaller
	testMethod1(v bool, arg1 string, arg2 uint16) (TestAllTypes, error)
	testMethod0(arg0 TestInterface) error
	testMethod2(arg0 TestInterface) error
}
type _TestInterface_remote struct {
	Ptr C.Pointer
}

func (p _TestInterface_remote) MarshalCaptain(new C.NewFunc) (C.Pointer, error) {
	return p.Ptr, nil
}
func getret_TestInterface_testMethod1(p C.Pointer) TestAllTypes {
	return TestAllTypes{Ptr: p}
}
func setret_TestInterface_testMethod1(new C.NewFunc, v TestAllTypes) (C.Pointer, error) {
	return v.Ptr, nil
}
func (p _TestInterface_remote) testMethod1(a0 bool, a1 string, a2 uint16) (ret TestAllTypes, err error) {
	var args _TestInterface_testMethod1_args
	args, err = new_TestInterface_testMethod1_args(p.Ptr.New)
	if err != nil {
		return
	}
	args.setV(a0)
	args.setArg1(a1)
	args.setArg2(a2)
	var rets C.Pointer
	rets, err = p.Ptr.Call(1, args.Ptr)
	if err != nil {
		return
	}
	ret = getret_TestInterface_testMethod1(rets)
	return
}
func (p _TestInterface_remote) testMethod0(a0 TestInterface) (err error) {
	var args _TestInterface_testMethod0_args
	args, err = new_TestInterface_testMethod0_args(p.Ptr.New)
	if err != nil {
		return
	}
	args.setArg0(a0)
	_, err = p.Ptr.Call(0, args.Ptr)
	return
}
func (p _TestInterface_remote) testMethod2(a0 TestInterface) (err error) {
	var args _TestInterface_testMethod2_args
	args, err = new_TestInterface_testMethod2_args(p.Ptr.New)
	if err != nil {
		return
	}
	args.setArg0(a0)
	_, err = p.Ptr.Call(2, args.Ptr)
	return
}
func DispatchTestInterface(iface interface{}, method int, args C.Pointer, retnew C.NewFunc) (C.Pointer, error) {
	p, ok := iface.(TestInterface)
	if !ok {
		return nil, C.ErrInvalidInterface
	}
	switch method {
	case 1:
		a := _TestInterface_testMethod1_args{Ptr: args}
		r, err := p.testMethod1(a.v(), a.arg1(), a.arg2())
		if err != nil {
			return nil, err
		}
		return setret_TestInterface_testMethod1(retnew, r)
	case 0:
		a := _TestInterface_testMethod0_args{Ptr: args}
		err := p.testMethod0(a.arg0())
		if err != nil {
			return nil, err
		}
		return nil, nil
	case 2:
		a := _TestInterface_testMethod2_args{Ptr: args}
		err := p.testMethod2(a.arg0())
		if err != nil {
			return nil, err
		}
		return nil, nil
	default:
		return nil, C.ErrInvalidInterface
	}
}

type TestEnum uint16

const (
	FOO TestEnum = 1
	BAR TestEnum = 2
)

type TestAllTypes struct {
	Ptr C.Pointer
}

func NewTestAllTypes(new C.NewFunc) (TestAllTypes, error) {
	ptr, err := new(C.MakeStruct(7, 21))
	return TestAllTypes{Ptr: ptr}, err
}
func (p TestAllTypes) boolField() bool {
	return (C.ReadUInt8(p.Ptr, 0) & 0) != 0
}
func (p TestAllTypes) setBoolField(v bool) error {
	return C.WriteBool(p.Ptr, 0, v)
}
func (p TestAllTypes) int8Field() int8 {
	return int8(C.ReadUInt8(p.Ptr, 1))
}
func (p TestAllTypes) setInt8Field(v int8) error {
	return C.WriteUInt8(p.Ptr, 1, uint8(v))
}
func (p TestAllTypes) int16Field() int16 {
	return int16(C.ReadUInt16(p.Ptr, 2))
}
func (p TestAllTypes) setInt16Field(v int16) error {
	return C.WriteUInt16(p.Ptr, 2, uint16(v))
}
func (p TestAllTypes) int32Field() int32 {
	return int32(C.ReadUInt32(p.Ptr, 4))
}
func (p TestAllTypes) setInt32Field(v int32) error {
	return C.WriteUInt32(p.Ptr, 4, uint32(v))
}
func (p TestAllTypes) int64Field() int64 {
	return int64(C.ReadUInt64(p.Ptr, 8))
}
func (p TestAllTypes) setInt64Field(v int64) error {
	return C.WriteUInt64(p.Ptr, 8, uint64(v))
}
func (p TestAllTypes) uInt8Field() uint8 {
	return C.ReadUInt8(p.Ptr, 16)
}
func (p TestAllTypes) setUInt8Field(v uint8) error {
	return C.WriteUInt8(p.Ptr, 16, v)
}
func (p TestAllTypes) uInt16Field() uint16 {
	return C.ReadUInt16(p.Ptr, 18)
}
func (p TestAllTypes) setUInt16Field(v uint16) error {
	return C.WriteUInt16(p.Ptr, 18, v)
}
func (p TestAllTypes) uInt32Field() uint32 {
	return C.ReadUInt32(p.Ptr, 20)
}
func (p TestAllTypes) setUInt32Field(v uint32) error {
	return C.WriteUInt32(p.Ptr, 20, v)
}
func (p TestAllTypes) uInt64Field() uint64 {
	return C.ReadUInt64(p.Ptr, 24)
}
func (p TestAllTypes) setUInt64Field(v uint64) error {
	return C.WriteUInt64(p.Ptr, 24, v)
}
func (p TestAllTypes) float32Field() float32 {
	return M.Float32frombits(C.ReadUInt32(p.Ptr, 32))
}
func (p TestAllTypes) setFloat32Field(v float32) error {
	return C.WriteUInt32(p.Ptr, 32, M.Float32bits(v))
}
func (p TestAllTypes) float64Field() float64 {
	return M.Float64frombits(C.ReadUInt64(p.Ptr, 40))
}
func (p TestAllTypes) setFloat64Field(v float64) error {
	return C.WriteUInt64(p.Ptr, 40, uint64(M.Float64bits(v)))
}
func (p TestAllTypes) textField() string {
	ptr := C.ReadPtr(p.Ptr, 0)
	return C.ToString(ptr, "")
}
func (p TestAllTypes) setTextField(v string) error {
	data, err := C.NewString(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(0, []C.Pointer{data})
}
func (p TestAllTypes) dataField() []uint8 {
	ptr := C.ReadPtr(p.Ptr, 1)
	return C.ToUInt8List(ptr)
}
func (p TestAllTypes) setDataField(v []uint8) error {
	data, err := C.NewUInt8List(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(1, []C.Pointer{data})
}
func (p TestAllTypes) structField() TestAllTypes {
	ptr := C.ReadPtr(p.Ptr, 2)
	return TestAllTypes{Ptr: ptr}
}
func (p TestAllTypes) setStructField(v TestAllTypes) error {
	return p.Ptr.WritePtrs(2, []C.Pointer{v.Ptr})
}
func (p TestAllTypes) enumField() TestEnum {
	return TestEnum(C.ReadUInt16(p.Ptr, 48))
}
func (p TestAllTypes) setEnumField(v TestEnum) error {
	return C.WriteUInt16(p.Ptr, 48, uint16(v))
}
func (p TestAllTypes) interfaceField() TestInterface {
	ptr := C.ReadPtr(p.Ptr, 3)
	return _TestInterface_remote{Ptr: ptr}
}
func (p TestAllTypes) setInterfaceField(v TestInterface) error {
	data, err := v.MarshalCaptain(p.Ptr.New)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(3, []C.Pointer{data})
}
func (p TestAllTypes) voidList() []struct{} {
	ptr := C.ReadPtr(p.Ptr, 4)
	return C.ToVoidList(ptr)
}
func (p TestAllTypes) setVoidList(v []struct{}) error {
	data, err := C.NewVoidList(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(4, []C.Pointer{data})
}
func (p TestAllTypes) boolList() C.Bitset {
	ptr := C.ReadPtr(p.Ptr, 5)
	return C.ToBitset(ptr)
}
func (p TestAllTypes) setBoolList(v C.Bitset) error {
	data, err := C.NewBitset(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(5, []C.Pointer{data})
}
func (p TestAllTypes) int8List() []int8 {
	ptr := C.ReadPtr(p.Ptr, 6)
	return C.ToInt8List(ptr)
}
func (p TestAllTypes) setInt8List(v []int8) error {
	data, err := C.NewInt8List(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(6, []C.Pointer{data})
}
func (p TestAllTypes) int16List() []int16 {
	ptr := C.ReadPtr(p.Ptr, 7)
	return C.ToInt16List(ptr)
}
func (p TestAllTypes) setInt16List(v []int16) error {
	data, err := C.NewInt16List(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(7, []C.Pointer{data})
}
func (p TestAllTypes) int32List() []int32 {
	ptr := C.ReadPtr(p.Ptr, 8)
	return C.ToInt32List(ptr)
}
func (p TestAllTypes) setInt32List(v []int32) error {
	data, err := C.NewInt32List(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(8, []C.Pointer{data})
}
func (p TestAllTypes) int64List() []int64 {
	ptr := C.ReadPtr(p.Ptr, 9)
	return C.ToInt64List(ptr)
}
func (p TestAllTypes) setInt64List(v []int64) error {
	data, err := C.NewInt64List(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(9, []C.Pointer{data})
}
func (p TestAllTypes) uInt8List() []uint8 {
	ptr := C.ReadPtr(p.Ptr, 10)
	return C.ToUInt8List(ptr)
}
func (p TestAllTypes) setUInt8List(v []uint8) error {
	data, err := C.NewUInt8List(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(10, []C.Pointer{data})
}
func (p TestAllTypes) uInt16List() []uint16 {
	ptr := C.ReadPtr(p.Ptr, 11)
	return C.ToUInt16List(ptr)
}
func (p TestAllTypes) setUInt16List(v []uint16) error {
	data, err := C.NewUInt16List(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(11, []C.Pointer{data})
}
func (p TestAllTypes) uInt32List() []uint32 {
	ptr := C.ReadPtr(p.Ptr, 12)
	return C.ToUInt32List(ptr)
}
func (p TestAllTypes) setUInt32List(v []uint32) error {
	data, err := C.NewUInt32List(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(12, []C.Pointer{data})
}
func (p TestAllTypes) uInt64List() []uint64 {
	ptr := C.ReadPtr(p.Ptr, 13)
	return C.ToUInt64List(ptr)
}
func (p TestAllTypes) setUInt64List(v []uint64) error {
	data, err := C.NewUInt64List(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(13, []C.Pointer{data})
}
func (p TestAllTypes) float32List() []float32 {
	ptr := C.ReadPtr(p.Ptr, 14)
	return C.ToFloat32List(ptr)
}
func (p TestAllTypes) setFloat32List(v []float32) error {
	data, err := C.NewFloat32List(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(14, []C.Pointer{data})
}
func (p TestAllTypes) float64List() []float64 {
	ptr := C.ReadPtr(p.Ptr, 15)
	return C.ToFloat64List(ptr)
}
func (p TestAllTypes) setFloat64List(v []float64) error {
	data, err := C.NewFloat64List(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(15, []C.Pointer{data})
}
func (p TestAllTypes) textList() []string {
	ptr := C.ReadPtr(p.Ptr, 16)
	return C.ToStringList(ptr)
}
func (p TestAllTypes) setTextList(v []string) error {
	data, err := C.NewStringList(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(16, []C.Pointer{data})
}
func (p TestAllTypes) dataList() []C.Pointer {
	ptr := C.ReadPtr(p.Ptr, 17)
	return C.ToPointerList(ptr)
}
func (p TestAllTypes) setDataList(v []C.Pointer) error {
	data, err := C.NewPointerList(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(17, []C.Pointer{data})
}
func (p TestAllTypes) structList() []TestAllTypes {
	ptr := C.ReadPtr(p.Ptr, 18)
	list := C.ToPointerList(ptr)
	ret := []TestAllTypes(nil)
	pret := (*reflect.SliceHeader)(unsafe.Pointer(&ret))
	*pret = *(*reflect.SliceHeader)(unsafe.Pointer(&list))
	return ret
}
func (p TestAllTypes) setStructList(v []TestAllTypes) error {
	ptrs := []C.Pointer(nil)
	pptrs := (*reflect.SliceHeader)(unsafe.Pointer(&ptrs))
	*pptrs = *(*reflect.SliceHeader)(unsafe.Pointer(&v))
	data, err := C.NewPointerList(p.Ptr.New, ptrs)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(18, []C.Pointer{data})
}
func (p TestAllTypes) enumList() []TestEnum {
	ptr := C.ReadPtr(p.Ptr, 19)
	u16 := C.ToUInt16List(ptr)
	ret := []TestEnum(nil)
	pret := (*reflect.SliceHeader)(unsafe.Pointer(&ret))
	*pret = *(*reflect.SliceHeader)(unsafe.Pointer(&u16))
	return ret
}
func (p TestAllTypes) setEnumList(v []TestEnum) error {
	u16 := []uint16(nil)
	pu16 := (*reflect.SliceHeader)(unsafe.Pointer(&u16))
	*pu16 = *(*reflect.SliceHeader)(unsafe.Pointer(&v))
	data, err := C.NewUInt16List(p.Ptr.New, u16)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(19, []C.Pointer{data})
}
func (p TestAllTypes) interfaceList() []TestInterface {
	ptr := C.ReadPtr(p.Ptr, 20)
	ptrs := C.ToPointerList(ptr)
	ret := make([]TestInterface, len(ptrs))
	for i := range ptrs {
		ret[i] = _TestInterface_remote{Ptr: ptr}
	}
	return ret
}
func (p TestAllTypes) setInterfaceList(v []TestInterface) error {
	cookies := make([]C.Pointer, len(v))
	for i, iface := range v {
		var err error
		cookies[i], err = iface.MarshalCaptain(p.Ptr.New)
		if err != nil {
			return err
		}
	}
	data, err := C.NewPointerList(p.Ptr.New, cookies)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(20, []C.Pointer{data})
}

type TestDefaults struct {
	Ptr C.Pointer
}

func NewTestDefaults(new C.NewFunc) (TestDefaults, error) {
	ptr, err := new(C.MakeStruct(7, 21))
	return TestDefaults{Ptr: ptr}, err
}
func (p TestDefaults) boolField() bool {
	return (C.ReadUInt8(p.Ptr, 0) & 0) != 1
}
func (p TestDefaults) setBoolField(v bool) error {
	return C.WriteBool(p.Ptr, 0, !v)
}
func (p TestDefaults) int8Field() int8 {
	return int8(C.ReadUInt8(p.Ptr, 1)) ^ -123
}
func (p TestDefaults) setInt8Field(v int8) error {
	return C.WriteUInt8(p.Ptr, 1, uint8(v^-123))
}
func (p TestDefaults) int16Field() int16 {
	return int16(C.ReadUInt16(p.Ptr, 2)) ^ -12345
}
func (p TestDefaults) setInt16Field(v int16) error {
	return C.WriteUInt16(p.Ptr, 2, uint16(v^-12345))
}
func (p TestDefaults) int32Field() int32 {
	return int32(C.ReadUInt32(p.Ptr, 4)) ^ -12345678
}
func (p TestDefaults) setInt32Field(v int32) error {
	return C.WriteUInt32(p.Ptr, 4, uint32(v^-12345678))
}
func (p TestDefaults) int64Field() int64 {
	return int64(C.ReadUInt64(p.Ptr, 8)) ^ -123456789012345
}
func (p TestDefaults) setInt64Field(v int64) error {
	return C.WriteUInt64(p.Ptr, 8, uint64(v^-123456789012345))
}
func (p TestDefaults) uInt8Field() uint8 {
	return C.ReadUInt8(p.Ptr, 16) ^ 234
}
func (p TestDefaults) setUInt8Field(v uint8) error {
	return C.WriteUInt8(p.Ptr, 16, v^234)
}
func (p TestDefaults) uInt16Field() uint16 {
	return C.ReadUInt16(p.Ptr, 18) ^ 45678
}
func (p TestDefaults) setUInt16Field(v uint16) error {
	return C.WriteUInt16(p.Ptr, 18, v^45678)
}
func (p TestDefaults) uInt32Field() uint32 {
	return C.ReadUInt32(p.Ptr, 20) ^ 3456789012
}
func (p TestDefaults) setUInt32Field(v uint32) error {
	return C.WriteUInt32(p.Ptr, 20, v^3456789012)
}
func (p TestDefaults) uInt64Field() uint64 {
	return C.ReadUInt64(p.Ptr, 24) ^ 12345678901234567890
}
func (p TestDefaults) setUInt64Field(v uint64) error {
	return C.WriteUInt64(p.Ptr, 24, v^12345678901234567890)
}
func (p TestDefaults) float32Field() float32 {
	u := C.ReadUInt32(p.Ptr, 32)
	u ^= M.Float32bits(1234.5)
	return M.Float32frombits(u)
}
func (p TestDefaults) setFloat32Field(v float32) error {
	return C.WriteUInt32(p.Ptr, 32, M.Float32bits(v)^M.Float32bits(1234.5))
}
func (p TestDefaults) float64Field() float64 {
	u := C.ReadUInt64(p.Ptr, 40)
	u ^= M.Float64bits(-123e45)
	return M.Float64frombits(u)
}
func (p TestDefaults) setFloat64Field(v float64) error {
	return C.WriteUInt64(p.Ptr, 40, uint64(M.Float64bits(v)^M.Float64bits(-123e45)))
}
func (p TestDefaults) textField() string {
	ptr := C.ReadPtr(p.Ptr, 0)
	return C.ToString(ptr, "foo")
}
func (p TestDefaults) setTextField(v string) error {
	data, err := C.NewString(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(0, []C.Pointer{data})
}
func (p TestDefaults) dataField() []uint8 {
	ptr := C.ReadPtr(p.Ptr, 1)
	ret := C.ToUInt8List(ptr)
	if ret == nil {
		ret = []byte("bar")
	}
	return ret
}
func (p TestDefaults) setDataField(v []uint8) error {
	data, err := C.NewUInt8List(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(1, []C.Pointer{data})
}
func (p TestDefaults) structField() TestAllTypes {
	ptr := C.ReadPtr(p.Ptr, 2)
	ret := TestAllTypes{Ptr: ptr}
	if ret.Ptr == nil {
		ret = func() TestAllTypes {
			p, _ := NewTestAllTypes(C.NewMemory)
			p.setBoolField(true)
			p.setInt8Field(-12)
			p.setInt16Field(3456)
			p.setInt32Field(-78901234)
			p.setInt64Field(56789012345678)
			p.setUInt8Field(90)
			p.setUInt16Field(1234)
			p.setUInt32Field(56789012)
			p.setUInt64Field(345678901234567890)
			p.setFloat32Field(-1.25e-10)
			p.setFloat64Field(345)
			p.setTextField("baz")
			p.setDataField([]byte("qux"))
			p.setStructField(func() TestAllTypes {
				p, _ := NewTestAllTypes(C.NewMemory)
				p.setTextField("nested")
				p.setStructField(func() TestAllTypes {
					p, _ := NewTestAllTypes(C.NewMemory)
					p.setTextField("really nested")
					return p
				}())
				return p
			}())
			p.setEnumField(FOO)
			p.setVoidList(make([]struct{}, 3))
			p.setBoolList(C.Bitset{0x1a})
			p.setInt8List([]int8{12, -34, -0x80, 0x7f})
			p.setInt16List([]int16{1234, -5678, -0x8000, 0x7fff})
			p.setInt32List([]int32{12345678, -90123456, -0x8000000, 0x7ffffff})
			p.setInt64List([]int64{123456789012345, -678901234567890, -0x8000000000000000, 0x7fffffffffffffff})
			p.setUInt8List([]uint8{12, 34, 0, 0xff})
			p.setUInt16List([]uint16{1234, 5678, 0, 0xffff})
			p.setUInt32List([]uint32{12345678, 90123456, 0, 0xffffffff})
			p.setUInt64List([]uint64{123456789012345, 678901234567890, 0, 0xffffffffffffffff})
			p.setFloat32List([]float32{0, 1234567, 1e37, -1e37, 1e-37, -1e-37})
			p.setFloat64List([]float64{0, 123456789012345, 1e306, -1e306, 1e-306, -1e-306})
			p.setTextList([]string{"quux", "corge", "grault"})
			p.setDataList([]C.Pointer{C.Must(C.NewUInt8List(C.NewMemory, []byte("garply"))), C.Must(C.NewUInt8List(C.NewMemory, []byte("waldo"))), C.Must(C.NewUInt8List(C.NewMemory, []byte("fred")))})
			p.setStructList([]TestAllTypes{func() TestAllTypes {
				p, _ := NewTestAllTypes(C.NewMemory)
				p.setTextField("x structlist 1")
				return p
			}(), func() TestAllTypes {
				p, _ := NewTestAllTypes(C.NewMemory)
				p.setTextField("x structlist 2")
				return p
			}(), func() TestAllTypes {
				p, _ := NewTestAllTypes(C.NewMemory)
				p.setTextField("x structlist 3")
				return p
			}()})
			p.setEnumList([]TestEnum{BAR, FOO})
			return p
		}()
	}
	return ret
}
func (p TestDefaults) setStructField(v TestAllTypes) error {
	return p.Ptr.WritePtrs(2, []C.Pointer{v.Ptr})
}
func (p TestDefaults) enumField() TestEnum {
	return TestEnum(C.ReadUInt16(p.Ptr, 48)) ^ FOO
}
func (p TestDefaults) setEnumField(v TestEnum) error {
	return C.WriteUInt16(p.Ptr, 48, uint16(v^FOO))
}

/* interface can't have a default
 */
func (p TestDefaults) interfaceField() TestInterface {
	ptr := C.ReadPtr(p.Ptr, 3)
	return _TestInterface_remote{Ptr: ptr}
}
func (p TestDefaults) setInterfaceField(v TestInterface) error {
	data, err := v.MarshalCaptain(p.Ptr.New)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(3, []C.Pointer{data})
}
func (p TestDefaults) voidList() []struct{} {
	ptr := C.ReadPtr(p.Ptr, 4)
	ret := C.ToVoidList(ptr)
	if ret == nil {
		ret = make([]struct{}, 6)
	}
	return ret
}
func (p TestDefaults) setVoidList(v []struct{}) error {
	data, err := C.NewVoidList(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(4, []C.Pointer{data})
}
func (p TestDefaults) boolList() C.Bitset {
	ptr := C.ReadPtr(p.Ptr, 5)
	ret := C.ToBitset(ptr)
	if ret == nil {
		ret = C.Bitset{0x09}
	}
	return ret
}
func (p TestDefaults) setBoolList(v C.Bitset) error {
	data, err := C.NewBitset(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(5, []C.Pointer{data})
}
func (p TestDefaults) int8List() []int8 {
	ptr := C.ReadPtr(p.Ptr, 6)
	ret := C.ToInt8List(ptr)
	if ret == nil {
		ret = []int8{111, -111}
	}
	return ret
}
func (p TestDefaults) setInt8List(v []int8) error {
	data, err := C.NewInt8List(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(6, []C.Pointer{data})
}
func (p TestDefaults) int16List() []int16 {
	ptr := C.ReadPtr(p.Ptr, 7)
	ret := C.ToInt16List(ptr)
	if ret == nil {
		ret = []int16{11111, -11111}
	}
	return ret
}
func (p TestDefaults) setInt16List(v []int16) error {
	data, err := C.NewInt16List(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(7, []C.Pointer{data})
}
func (p TestDefaults) int32List() []int32 {
	ptr := C.ReadPtr(p.Ptr, 8)
	ret := C.ToInt32List(ptr)
	if ret == nil {
		ret = []int32{111111111, -111111111}
	}
	return ret
}
func (p TestDefaults) setInt32List(v []int32) error {
	data, err := C.NewInt32List(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(8, []C.Pointer{data})
}
func (p TestDefaults) int64List() []int64 {
	ptr := C.ReadPtr(p.Ptr, 9)
	ret := C.ToInt64List(ptr)
	if ret == nil {
		ret = []int64{1111111111111111111, -1111111111111111111}
	}
	return ret
}
func (p TestDefaults) setInt64List(v []int64) error {
	data, err := C.NewInt64List(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(9, []C.Pointer{data})
}
func (p TestDefaults) uInt8List() []uint8 {
	ptr := C.ReadPtr(p.Ptr, 10)
	ret := C.ToUInt8List(ptr)
	if ret == nil {
		ret = []uint8{111, 222}
	}
	return ret
}
func (p TestDefaults) setUInt8List(v []uint8) error {
	data, err := C.NewUInt8List(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(10, []C.Pointer{data})
}
func (p TestDefaults) uInt16List() []uint16 {
	ptr := C.ReadPtr(p.Ptr, 11)
	ret := C.ToUInt16List(ptr)
	if ret == nil {
		ret = []uint16{33333, 44444}
	}
	return ret
}
func (p TestDefaults) setUInt16List(v []uint16) error {
	data, err := C.NewUInt16List(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(11, []C.Pointer{data})
}
func (p TestDefaults) uInt32List() []uint32 {
	ptr := C.ReadPtr(p.Ptr, 12)
	ret := C.ToUInt32List(ptr)
	if ret == nil {
		ret = []uint32{3333333333}
	}
	return ret
}
func (p TestDefaults) setUInt32List(v []uint32) error {
	data, err := C.NewUInt32List(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(12, []C.Pointer{data})
}
func (p TestDefaults) uInt64List() []uint64 {
	ptr := C.ReadPtr(p.Ptr, 13)
	ret := C.ToUInt64List(ptr)
	if ret == nil {
		ret = []uint64{11111111111111111111}
	}
	return ret
}
func (p TestDefaults) setUInt64List(v []uint64) error {
	data, err := C.NewUInt64List(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(13, []C.Pointer{data})
}
func (p TestDefaults) float32List() []float32 {
	ptr := C.ReadPtr(p.Ptr, 14)
	ret := C.ToFloat32List(ptr)
	if ret == nil {
		ret = []float32{5555.5, 2222.25}
	}
	return ret
}
func (p TestDefaults) setFloat32List(v []float32) error {
	data, err := C.NewFloat32List(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(14, []C.Pointer{data})
}
func (p TestDefaults) float64List() []float64 {
	ptr := C.ReadPtr(p.Ptr, 15)
	ret := C.ToFloat64List(ptr)
	if ret == nil {
		ret = []float64{7777.75, 1111.125}
	}
	return ret
}
func (p TestDefaults) setFloat64List(v []float64) error {
	data, err := C.NewFloat64List(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(15, []C.Pointer{data})
}
func (p TestDefaults) textList() []string {
	ptr := C.ReadPtr(p.Ptr, 16)
	ret := C.ToStringList(ptr)
	if ret == nil {
		ret = []string{"plugh", "xyzzy", "thud"}
	}
	return ret
}
func (p TestDefaults) setTextList(v []string) error {
	data, err := C.NewStringList(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(16, []C.Pointer{data})
}
func (p TestDefaults) dataList() []C.Pointer {
	ptr := C.ReadPtr(p.Ptr, 17)
	ret := C.ToPointerList(ptr)
	if ret == nil {
		ret = []C.Pointer{C.Must(C.NewUInt8List(C.NewMemory, []byte("oops"))), C.Must(C.NewUInt8List(C.NewMemory, []byte("exhausted"))), C.Must(C.NewUInt8List(C.NewMemory, []byte("rfc3092")))}
	}
	return ret
}
func (p TestDefaults) setDataList(v []C.Pointer) error {
	data, err := C.NewPointerList(p.Ptr.New, v)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(17, []C.Pointer{data})
}
func (p TestDefaults) structList() []TestAllTypes {
	ptr := C.ReadPtr(p.Ptr, 18)
	list := C.ToPointerList(ptr)
	ret := []TestAllTypes(nil)
	pret := (*reflect.SliceHeader)(unsafe.Pointer(&ret))
	*pret = *(*reflect.SliceHeader)(unsafe.Pointer(&list))
	if ret == nil {
		ret = []TestAllTypes{func() TestAllTypes {
			p, _ := NewTestAllTypes(C.NewMemory)
			p.setTextField("structlist 1")
			return p
		}(), func() TestAllTypes {
			p, _ := NewTestAllTypes(C.NewMemory)
			p.setTextField("structlist 2")
			return p
		}(), func() TestAllTypes {
			p, _ := NewTestAllTypes(C.NewMemory)
			p.setTextField("structlist 3")
			return p
		}()}
	}
	return ret
}
func (p TestDefaults) setStructList(v []TestAllTypes) error {
	ptrs := []C.Pointer(nil)
	pptrs := (*reflect.SliceHeader)(unsafe.Pointer(&ptrs))
	*pptrs = *(*reflect.SliceHeader)(unsafe.Pointer(&v))
	data, err := C.NewPointerList(p.Ptr.New, ptrs)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(18, []C.Pointer{data})
}
func (p TestDefaults) enumList() []TestEnum {
	ptr := C.ReadPtr(p.Ptr, 19)
	u16 := C.ToUInt16List(ptr)
	ret := []TestEnum(nil)
	pret := (*reflect.SliceHeader)(unsafe.Pointer(&ret))
	*pret = *(*reflect.SliceHeader)(unsafe.Pointer(&u16))
	if ret == nil {
		ret = []TestEnum{FOO, BAR}
	}
	return ret
}
func (p TestDefaults) setEnumList(v []TestEnum) error {
	u16 := []uint16(nil)
	pu16 := (*reflect.SliceHeader)(unsafe.Pointer(&u16))
	*pu16 = *(*reflect.SliceHeader)(unsafe.Pointer(&v))
	data, err := C.NewUInt16List(p.Ptr.New, u16)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(19, []C.Pointer{data})
}
func (p TestDefaults) interfaceList() []TestInterface {
	ptr := C.ReadPtr(p.Ptr, 20)
	ptrs := C.ToPointerList(ptr)
	ret := make([]TestInterface, len(ptrs))
	for i := range ptrs {
		ret[i] = _TestInterface_remote{Ptr: ptr}
	}
	return ret
}
func (p TestDefaults) setInterfaceList(v []TestInterface) error {
	cookies := make([]C.Pointer, len(v))
	for i, iface := range v {
		var err error
		cookies[i], err = iface.MarshalCaptain(p.Ptr.New)
		if err != nil {
			return err
		}
	}
	data, err := C.NewPointerList(p.Ptr.New, cookies)
	if err != nil {
		return err
	}
	return p.Ptr.WritePtrs(20, []C.Pointer{data})
}
