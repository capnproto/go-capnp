package test

import (
	"encoding/binary"
	C "github.com/jmckaskill/go-capnproto"
	M "math"
)

var (
	putLittle16 = binary.LittleEndian.PutUint16
	_           = M.Float32bits

	boolConst                = true
	int8Const        int8    = -123
	int16Const       int16   = -12345
	int32Const       int32   = -12345678
	int64Const       int64   = -123456789012345
	uInt8Const       uint8   = 234
	uInt16Const      uint16  = 45678
	uInt32Const      uint32  = 3456789012
	uInt64Const      uint64  = 12345678901234567890
	float32Const     float32 = 1234.5
	float64Const             = -1.23e+47
	textConst                = "foo"
	dataConst                = []byte("bar")
	structConst              = TestAllTypes{test_capnp.ReadPtr(0)}
	enumConst                = TestEnum(1)
	voidListConst            = make([]struct{}, 6)
	boolListConst            = test_capnp.ReadBitset(237, C.Bitset{})
	int8ListConst            = []int8{111, -111}
	int16ListConst           = []int16{11111, -11111}
	int32ListConst           = []int32{111111111, -111111111}
	int64ListConst           = []int64{1111111111111111111, -1111111111111111111}
	uInt8ListConst           = []uint8{111, 222}
	uInt16ListConst          = []uint16{33333, 44444}
	uInt32ListConst          = []uint32{3333333333}
	uInt64ListConst          = []uint64{11111111111111111111}
	float32ListConst         = []float32{5555.5, 2222.25}
	float64ListConst         = []float64{7777.75, 1111.125}
	textListConst            = []C.Pointer{test_capnp.ReadPtr(239), test_capnp.ReadPtr(241), test_capnp.ReadPtr(243)}
	dataListConst            = []C.Pointer{test_capnp.ReadPtr(245), test_capnp.ReadPtr(247), test_capnp.ReadPtr(250)}
	structListConst          = []TestAllTypes{test_capnp_0, test_capnp_1, test_capnp_2}
	enumListConst            = []TestEnum{TestEnum(1), TestEnum(2)}
)

type TestInterface interface {
	C.Marshaller
	testMethod1(v bool, arg1 string, arg2 uint16) TestAllTypes
	testMethod0(arg0 TestInterface) int8
	testMethod2(arg0 TestInterface)
	testMultiRet(arg0 bool, arg1 string) (v uint16, ret1 string)
}

type RemoteTestInterface struct{ P C.Pointer }

func (p RemoteTestInterface) MarshalCaptain(r C.Pointer, i int) error {
	return r.WritePtr(i, p.P)
}
func (p RemoteTestInterface) testMethod1(a0 bool, a1 string, a2 uint16) (r0 TestAllTypes) {
	c, err := p.P.Segment.Session.NewCall()
	if err != nil {
		return
	}
	c.Message.SetObject(p.P)
	c.Message.SetMethod(1)
	args, _ := c.Message.P.Segment.NewStruct(32, 1)
	args.WriteStruct1(0, a0)
	args.WriteString(0, a1, "")
	args.WriteStruct16(16, a2)
	c.Message.SetArguments(args)
	if c.Send(c) != nil {
		return
	}
	return TestAllTypes{c.Reply.ReadPtr(0)}
}
func (p RemoteTestInterface) testMethod0(a0 TestInterface) (r0 int8) {
	c, err := p.P.Segment.Session.NewCall()
	if err != nil {
		return
	}
	c.Message.SetObject(p.P)
	c.Message.SetMethod(0)
	args, _ := c.Message.P.Segment.NewStruct(0, 1)
	a0.MarshalCaptain(args, 0)
	c.Message.SetArguments(args)
	if c.Send(c) != nil {
		return
	}
	return int8(c.Reply.ReadStruct8(0))
}
func (p RemoteTestInterface) testMethod2(a0 TestInterface) {
	c, err := p.P.Segment.Session.NewCall()
	if err != nil {
		return
	}
	c.Message.SetObject(p.P)
	c.Message.SetMethod(2)
	args, _ := c.Message.P.Segment.NewStruct(0, 1)
	a0.MarshalCaptain(args, 0)
	c.Message.SetArguments(args)
	c.Send(c)
}
func (p RemoteTestInterface) testMultiRet(a0 bool, a1 string) (r0 uint16, r1 string) {
	r1 = "abc"
	c, err := p.P.Segment.Session.NewCall()
	if err != nil {
		return
	}
	c.Message.SetObject(p.P)
	c.Message.SetMethod(3)
	args, _ := c.Message.P.Segment.NewStruct(1, 1)
	args.WriteStruct1(0, a0)
	args.WriteString(0, a1, "")
	c.Message.SetArguments(args)
	if c.Send(c) != nil {
		return
	}
	return c.Reply.ReadStruct16(0), c.Reply.ReadString(0, "abc")
}

func DispatchTestInterface(p TestInterface, in, out C.Message) error {
	switch in.Method() {
	case 1:
		a := in.Arguments()
		r0 := p.testMethod1(a.ReadStruct1(0), a.ReadString(0, ""), a.ReadStruct16(16))
		ret, err := out.P.Segment.NewStruct(0, 1)
		if err != nil {
			return err
		}
		if err := ret.WritePtr(0, r0.P); err != nil {
			return err
		}
		return out.SetArguments(ret)
	case 0:
		a := in.Arguments()
		r0 := p.testMethod0(RemoteTestInterface{a.ReadPtr(0)})
		ret, err := out.P.Segment.NewStruct(8, 0)
		if err != nil {
			return err
		}
		if err := ret.WriteStruct8(0, uint8(r0)); err != nil {
			return err
		}
		return out.SetArguments(ret)
	case 2:
		a := in.Arguments()
		p.testMethod2(RemoteTestInterface{a.ReadPtr(0)})
		return nil
	case 3:
		a := in.Arguments()
		r0, r1 := p.testMultiRet(a.ReadStruct1(0), a.ReadString(0, ""))
		ret, err := out.P.Segment.NewStruct(16, 1)
		if err != nil {
			return err
		}
		if err := ret.WriteStruct16(0, r0); err != nil {
			return err
		}
		if err := ret.WriteString(0, r1, "abc"); err != nil {
			return err
		}
		return out.SetArguments(ret)
	}
	return C.ErrInvalidInterface
}

func ReadTestInterfaceList(p C.Pointer, i int, def []TestInterface) []TestInterface {
	if m := p.ReadPtr(i); m.Type() == C.List {
		r := make([]TestInterface, m.Size())
		for i := range r {
			r[i] = RemoteTestInterface{m.ReadPtr(i)}
		}
		return r
	}
	return def
}
func WriteTestInterfaceList(p C.Pointer, i int, v, def []TestInterface) (err error) {
	var m C.Pointer
	if !C.SliceEqual(v, def) {
		if m, err = p.Segment.NewPointerList(len(v)); err != nil {
			return err
		}
		for i, u := range v {
			if err := u.MarshalCaptain(m, i); err != nil {
				return err
			}
		}
	}
	return p.WritePtr(i, m)
}

type TestEnum uint16

const (
	TestEnum_FOO TestEnum = 1
	TestEnum_BAR TestEnum = 2
)

func ReadTestEnumList(p C.Pointer, i int, def []TestEnum) []TestEnum {
	if m := p.ReadPtr(i); m.Type() == C.List {
		r := make([]TestEnum, m.Size())
		for i := range r {
			r[i] = TestEnum(m.Read16(i))
		}
		return r
	}
	return def
}
func WriteTestEnumList(p C.Pointer, i int, v, def []TestEnum) (err error) {
	var m C.Pointer
	if !C.SliceEqual(v, def) {
		if m, err = p.Segment.NewList(16, 0, len(v)); err != nil {
			return err
		}
		d := p.Data()
		for i, u := range v {
			putLittle16(d[2*i:], uint16(u))
		}
	}
	return p.WritePtr(i, m)
}

type TestAllTypes_unionField uint16

const (
	TestAllTypes_voidUnion      TestAllTypes_unionField = 35
	TestAllTypes_boolUnion      TestAllTypes_unionField = 36
	TestAllTypes_int8Union      TestAllTypes_unionField = 37
	TestAllTypes_uint8Union     TestAllTypes_unionField = 38
	TestAllTypes_int16Union     TestAllTypes_unionField = 39
	TestAllTypes_uint16Union    TestAllTypes_unionField = 40
	TestAllTypes_int32Union     TestAllTypes_unionField = 41
	TestAllTypes_uint32Union    TestAllTypes_unionField = 42
	TestAllTypes_int64Union     TestAllTypes_unionField = 43
	TestAllTypes_uint64Union    TestAllTypes_unionField = 44
	TestAllTypes_float32Union   TestAllTypes_unionField = 45
	TestAllTypes_float64Union   TestAllTypes_unionField = 46
	TestAllTypes_textUnion      TestAllTypes_unionField = 47
	TestAllTypes_dataUnion      TestAllTypes_unionField = 48
	TestAllTypes_structUnion    TestAllTypes_unionField = 49
	TestAllTypes_enumUnion      TestAllTypes_unionField = 50
	TestAllTypes_interfaceUnion TestAllTypes_unionField = 51
)

type TestAllTypes struct{ P C.Pointer }

func NewTestAllTypes(seg *C.Segment) (TestAllTypes, error) {
	p, err := seg.NewStruct(576, 22)
	return TestAllTypes{p}, err
}
func NewTestAllTypesList(seg *C.Segment, sz int) (C.Pointer, error) {
	return seg.NewList(576, 22, sz)
}
func (p TestAllTypes) MarshalCaptain(r C.Pointer, i int) error {
	return r.WritePtr(i, p.P)
}

func (p TestAllTypes) boolField() bool                { return p.P.ReadStruct1(0) }
func (p TestAllTypes) int8Field() int8                { return int8(p.P.ReadStruct8(8)) }
func (p TestAllTypes) int16Field() int16              { return int16(p.P.ReadStruct16(16)) }
func (p TestAllTypes) int32Field() int32              { return int32(p.P.ReadStruct32(32)) }
func (p TestAllTypes) int64Field() int64              { return int64(p.P.ReadStruct64(64)) }
func (p TestAllTypes) uInt8Field() uint8              { return p.P.ReadStruct8(128) }
func (p TestAllTypes) uInt16Field() uint16            { return p.P.ReadStruct16(144) }
func (p TestAllTypes) uInt32Field() uint32            { return p.P.ReadStruct32(160) }
func (p TestAllTypes) uInt64Field() uint64            { return p.P.ReadStruct64(192) }
func (p TestAllTypes) float32Field() float32          { return M.Float32frombits(p.P.ReadStruct32(256)) }
func (p TestAllTypes) float64Field() float64          { return M.Float64frombits(p.P.ReadStruct64(320)) }
func (p TestAllTypes) textField() string              { return p.P.ReadString(0, "") }
func (p TestAllTypes) dataField() []byte              { return p.P.ReadData(1, nil) }
func (p TestAllTypes) structField() TestAllTypes      { return TestAllTypes{p.P.ReadPtr(2)} }
func (p TestAllTypes) enumField() TestEnum            { return TestEnum(p.P.ReadStruct16(384)) }
func (p TestAllTypes) interfaceField() TestInterface  { return RemoteTestInterface{p.P.ReadPtr(3)} }
func (p TestAllTypes) voidList() []struct{}           { return p.P.ReadVoidList(4, nil) }
func (p TestAllTypes) boolList() C.Bitset             { return p.P.ReadBitset(5, C.Bitset{}) }
func (p TestAllTypes) int8List() []int8               { return p.P.ReadI8List(6, nil) }
func (p TestAllTypes) int16List() []int16             { return p.P.ReadI16List(7, nil) }
func (p TestAllTypes) int32List() []int32             { return p.P.ReadI32List(8, nil) }
func (p TestAllTypes) int64List() []int64             { return p.P.ReadI64List(9, nil) }
func (p TestAllTypes) uInt8List() []uint8             { return p.P.ReadU8List(10, nil) }
func (p TestAllTypes) uInt16List() []uint16           { return p.P.ReadU16List(11, nil) }
func (p TestAllTypes) uInt32List() []uint32           { return p.P.ReadU32List(12, nil) }
func (p TestAllTypes) uInt64List() []uint64           { return p.P.ReadU64List(13, nil) }
func (p TestAllTypes) float32List() []float32         { return p.P.ReadF32List(14, nil) }
func (p TestAllTypes) float64List() []float64         { return p.P.ReadF64List(15, nil) }
func (p TestAllTypes) textList() []C.Pointer          { return p.P.ReadPointerList(16, nil) }
func (p TestAllTypes) dataList() []C.Pointer          { return p.P.ReadPointerList(17, nil) }
func (p TestAllTypes) structList() []TestAllTypes     { return ReadTestAllTypesList(p.P, 18, nil) }
func (p TestAllTypes) enumList() []TestEnum           { return ReadTestEnumList(p.P, 19, nil) }
func (p TestAllTypes) interfaceList() []TestInterface { return ReadTestInterfaceList(p.P, 20, nil) }
func (p TestAllTypes) unionField() TestAllTypes_unionField {
	return TestAllTypes_unionField(p.P.ReadStruct16(400))
}
func (p TestAllTypes) boolUnion() (ret bool) {
	if p.P.ReadStruct16(400) != 36 {
		return
	}
	return p.P.ReadStruct1(416)
}
func (p TestAllTypes) int8Union() (ret int8) {
	if p.P.ReadStruct16(400) != 37 {
		return
	}
	return int8(p.P.ReadStruct8(424))
}
func (p TestAllTypes) uint8Union() (ret uint8) {
	if p.P.ReadStruct16(400) != 38 {
		return
	}
	return p.P.ReadStruct8(424)
}
func (p TestAllTypes) int16Union() (ret int16) {
	if p.P.ReadStruct16(400) != 39 {
		return
	}
	return int16(p.P.ReadStruct16(432))
}
func (p TestAllTypes) uint16Union() (ret uint16) {
	if p.P.ReadStruct16(400) != 40 {
		return
	}
	return p.P.ReadStruct16(432)
}
func (p TestAllTypes) int32Union() (ret int32) {
	if p.P.ReadStruct16(400) != 41 {
		return
	}
	return int32(p.P.ReadStruct32(448))
}
func (p TestAllTypes) uint32Union() (ret uint32) {
	if p.P.ReadStruct16(400) != 42 {
		return
	}
	return p.P.ReadStruct32(448)
}
func (p TestAllTypes) int64Union() (ret int64) {
	if p.P.ReadStruct16(400) != 43 {
		return
	}
	return int64(p.P.ReadStruct64(512))
}
func (p TestAllTypes) uint64Union() (ret uint64) {
	if p.P.ReadStruct16(400) != 44 {
		return
	}
	return p.P.ReadStruct64(512)
}
func (p TestAllTypes) float32Union() (ret float32) {
	if p.P.ReadStruct16(400) != 45 {
		return
	}
	return M.Float32frombits(p.P.ReadStruct32(512))
}
func (p TestAllTypes) float64Union() (ret float64) {
	if p.P.ReadStruct16(400) != 46 {
		return
	}
	return M.Float64frombits(p.P.ReadStruct64(512))
}
func (p TestAllTypes) textUnion() (ret string) {
	if p.P.ReadStruct16(400) != 47 {
		return
	}
	return p.P.ReadString(21, "")
}
func (p TestAllTypes) dataUnion() (ret []byte) {
	if p.P.ReadStruct16(400) != 48 {
		return
	}
	return p.P.ReadData(21, nil)
}
func (p TestAllTypes) structUnion() (ret TestAllTypes) {
	if p.P.ReadStruct16(400) != 49 {
		return
	}
	return TestAllTypes{p.P.ReadPtr(21)}
}
func (p TestAllTypes) enumUnion() (ret TestEnum) {
	if p.P.ReadStruct16(400) != 50 {
		return
	}
	return TestEnum(p.P.ReadStruct16(512))
}
func (p TestAllTypes) interfaceUnion() (ret TestInterface) {
	if p.P.ReadStruct16(400) != 51 {
		return
	}
	return RemoteTestInterface{p.P.ReadPtr(21)}
}

func (p TestAllTypes) setBoolField(v bool) error     { return p.P.WriteStruct1(0, v) }
func (p TestAllTypes) setInt8Field(v int8) error     { return p.P.WriteStruct8(8, uint8(v)) }
func (p TestAllTypes) setInt16Field(v int16) error   { return p.P.WriteStruct16(16, uint16(v)) }
func (p TestAllTypes) setInt32Field(v int32) error   { return p.P.WriteStruct32(32, uint32(v)) }
func (p TestAllTypes) setInt64Field(v int64) error   { return p.P.WriteStruct64(64, uint64(v)) }
func (p TestAllTypes) setUInt8Field(v uint8) error   { return p.P.WriteStruct8(128, v) }
func (p TestAllTypes) setUInt16Field(v uint16) error { return p.P.WriteStruct16(144, v) }
func (p TestAllTypes) setUInt32Field(v uint32) error { return p.P.WriteStruct32(160, v) }
func (p TestAllTypes) setUInt64Field(v uint64) error { return p.P.WriteStruct64(192, v) }
func (p TestAllTypes) setFloat32Field(v float32) error {
	return p.P.WriteStruct32(256, M.Float32bits(v))
}
func (p TestAllTypes) setFloat64Field(v float64) error {
	return p.P.WriteStruct64(320, M.Float64bits(v))
}
func (p TestAllTypes) setTextField(v string) error             { return p.P.WriteString(0, v, "") }
func (p TestAllTypes) setDataField(v []byte) error             { return p.P.WriteU8List(1, v, nil) }
func (p TestAllTypes) setStructField(v TestAllTypes) error     { return p.P.WritePtr(2, v.P) }
func (p TestAllTypes) setEnumField(v TestEnum) error           { return p.P.WriteStruct16(384, uint16(v)) }
func (p TestAllTypes) setInterfaceField(v TestInterface) error { return v.MarshalCaptain(p.P, 3) }
func (p TestAllTypes) setVoidList(v []struct{}) error          { return p.P.WriteVoidList(4, v, nil) }
func (p TestAllTypes) setBoolList(v C.Bitset) error            { return p.P.WriteBitset(5, v, C.Bitset{}) }
func (p TestAllTypes) setInt8List(v []int8) error              { return p.P.WriteI8List(6, v, nil) }
func (p TestAllTypes) setInt16List(v []int16) error            { return p.P.WriteI16List(7, v, nil) }
func (p TestAllTypes) setInt32List(v []int32) error            { return p.P.WriteI32List(8, v, nil) }
func (p TestAllTypes) setInt64List(v []int64) error            { return p.P.WriteI64List(9, v, nil) }
func (p TestAllTypes) setUInt8List(v []uint8) error            { return p.P.WriteU8List(10, v, nil) }
func (p TestAllTypes) setUInt16List(v []uint16) error          { return p.P.WriteU16List(11, v, nil) }
func (p TestAllTypes) setUInt32List(v []uint32) error          { return p.P.WriteU32List(12, v, nil) }
func (p TestAllTypes) setUInt64List(v []uint64) error          { return p.P.WriteU64List(13, v, nil) }
func (p TestAllTypes) setFloat32List(v []float32) error        { return p.P.WriteF32List(14, v, nil) }
func (p TestAllTypes) setFloat64List(v []float64) error        { return p.P.WriteF64List(15, v, nil) }
func (p TestAllTypes) setTextList(v []C.Pointer) error         { return p.P.WritePointerList(16, v, nil) }
func (p TestAllTypes) setDataList(v []C.Pointer) error         { return p.P.WritePointerList(17, v, nil) }
func (p TestAllTypes) setStructList(v []TestAllTypes) error {
	return WriteTestAllTypesList(p.P, 18, v, nil)
}
func (p TestAllTypes) setEnumList(v []TestEnum) error { return WriteTestEnumList(p.P, 19, v, nil) }
func (p TestAllTypes) setInterfaceList(v []TestInterface) error {
	return WriteTestInterfaceList(p.P, 20, v, nil)
}
func (p TestAllTypes) setBoolUnion(v bool) error {
	if err := p.P.WriteStruct16(400, 36); err != nil {
		return err
	}
	return p.P.WriteStruct1(416, v)
}
func (p TestAllTypes) setInt8Union(v int8) error {
	if err := p.P.WriteStruct16(400, 37); err != nil {
		return err
	}
	return p.P.WriteStruct8(424, uint8(v))
}
func (p TestAllTypes) setUint8Union(v uint8) error {
	if err := p.P.WriteStruct16(400, 38); err != nil {
		return err
	}
	return p.P.WriteStruct8(424, v)
}
func (p TestAllTypes) setInt16Union(v int16) error {
	if err := p.P.WriteStruct16(400, 39); err != nil {
		return err
	}
	return p.P.WriteStruct16(432, uint16(v))
}
func (p TestAllTypes) setUint16Union(v uint16) error {
	if err := p.P.WriteStruct16(400, 40); err != nil {
		return err
	}
	return p.P.WriteStruct16(432, v)
}
func (p TestAllTypes) setInt32Union(v int32) error {
	if err := p.P.WriteStruct16(400, 41); err != nil {
		return err
	}
	return p.P.WriteStruct32(448, uint32(v))
}
func (p TestAllTypes) setUint32Union(v uint32) error {
	if err := p.P.WriteStruct16(400, 42); err != nil {
		return err
	}
	return p.P.WriteStruct32(448, v)
}
func (p TestAllTypes) setInt64Union(v int64) error {
	if err := p.P.WriteStruct16(400, 43); err != nil {
		return err
	}
	return p.P.WriteStruct64(512, uint64(v))
}
func (p TestAllTypes) setUint64Union(v uint64) error {
	if err := p.P.WriteStruct16(400, 44); err != nil {
		return err
	}
	return p.P.WriteStruct64(512, v)
}
func (p TestAllTypes) setFloat32Union(v float32) error {
	if err := p.P.WriteStruct16(400, 45); err != nil {
		return err
	}
	return p.P.WriteStruct32(512, M.Float32bits(v))
}
func (p TestAllTypes) setFloat64Union(v float64) error {
	if err := p.P.WriteStruct16(400, 46); err != nil {
		return err
	}
	return p.P.WriteStruct64(512, M.Float64bits(v))
}
func (p TestAllTypes) setTextUnion(v string) error {
	if err := p.P.WriteStruct16(400, 47); err != nil {
		return err
	}
	return p.P.WriteString(21, v, "")
}
func (p TestAllTypes) setDataUnion(v []byte) error {
	if err := p.P.WriteStruct16(400, 48); err != nil {
		return err
	}
	return p.P.WriteU8List(21, v, nil)
}
func (p TestAllTypes) setStructUnion(v TestAllTypes) error {
	if err := p.P.WriteStruct16(400, 49); err != nil {
		return err
	}
	return p.P.WritePtr(21, v.P)
}
func (p TestAllTypes) setEnumUnion(v TestEnum) error {
	if err := p.P.WriteStruct16(400, 50); err != nil {
		return err
	}
	return p.P.WriteStruct16(512, uint16(v))
}
func (p TestAllTypes) setInterfaceUnion(v TestInterface) error {
	if err := p.P.WriteStruct16(400, 51); err != nil {
		return err
	}
	return v.MarshalCaptain(p.P, 21)
}

func ReadTestAllTypesList(p C.Pointer, i int, def []TestAllTypes) []TestAllTypes {
	if m := p.ReadPtr(i); m.Type() == C.List {
		r := make([]TestAllTypes, m.Size())
		for i := range r {
			r[i] = TestAllTypes{m.ReadPtr(i)}
		}
		return r
	}
	return def
}
func WriteTestAllTypesList(p C.Pointer, i int, v, def []TestAllTypes) (err error) {
	var m C.Pointer
	if !C.SliceEqual(v, def) {
		if m, err = p.Segment.NewPointerList(len(v)); err != nil {
			return err
		}
		for i, u := range v {
			if err := u.MarshalCaptain(m, i); err != nil {
				return err
			}
		}
	}
	return p.WritePtr(i, m)
}

type TestDefaults_unionField uint16

const (
	TestDefaults_voidUnion      TestDefaults_unionField = 35
	TestDefaults_boolUnion      TestDefaults_unionField = 36
	TestDefaults_int8Union      TestDefaults_unionField = 37
	TestDefaults_uint8Union     TestDefaults_unionField = 38
	TestDefaults_int16Union     TestDefaults_unionField = 39
	TestDefaults_uint16Union    TestDefaults_unionField = 40
	TestDefaults_int32Union     TestDefaults_unionField = 41
	TestDefaults_uint32Union    TestDefaults_unionField = 42
	TestDefaults_int64Union     TestDefaults_unionField = 43
	TestDefaults_uint64Union    TestDefaults_unionField = 44
	TestDefaults_float32Union   TestDefaults_unionField = 45
	TestDefaults_float64Union   TestDefaults_unionField = 46
	TestDefaults_textUnion      TestDefaults_unionField = 47
	TestDefaults_dataUnion      TestDefaults_unionField = 48
	TestDefaults_structUnion    TestDefaults_unionField = 49
	TestDefaults_enumUnion      TestDefaults_unionField = 50
	TestDefaults_interfaceUnion TestDefaults_unionField = 51
)

type TestDefaults struct{ P C.Pointer }

func NewTestDefaults(seg *C.Segment) (TestDefaults, error) {
	p, err := seg.NewStruct(576, 22)
	return TestDefaults{p}, err
}
func NewTestDefaultsList(seg *C.Segment, sz int) (C.Pointer, error) {
	return seg.NewList(576, 22, sz)
}
func (p TestDefaults) MarshalCaptain(r C.Pointer, i int) error {
	return r.WritePtr(i, p.P)
}

func (p TestDefaults) boolField() bool       { return !p.P.ReadStruct1(0) }
func (p TestDefaults) int8Field() int8       { return int8(p.P.ReadStruct8(8)) ^ -123 }
func (p TestDefaults) int16Field() int16     { return int16(p.P.ReadStruct16(16)) ^ -12345 }
func (p TestDefaults) int32Field() int32     { return int32(p.P.ReadStruct32(32)) ^ -12345678 }
func (p TestDefaults) int64Field() int64     { return int64(p.P.ReadStruct64(64)) ^ -123456789012345 }
func (p TestDefaults) uInt8Field() uint8     { return p.P.ReadStruct8(128) ^ 234 }
func (p TestDefaults) uInt16Field() uint16   { return p.P.ReadStruct16(144) ^ 45678 }
func (p TestDefaults) uInt32Field() uint32   { return p.P.ReadStruct32(160) ^ 3456789012 }
func (p TestDefaults) uInt64Field() uint64   { return p.P.ReadStruct64(192) ^ 12345678901234567890 }
func (p TestDefaults) float32Field() float32 { return p.P.ReadStructF32(256, 1234.5) }
func (p TestDefaults) float64Field() float64 { return p.P.ReadStructF64(320, -1.23e+47) }
func (p TestDefaults) textField() string     { return p.P.ReadString(0, "foo") }
func (p TestDefaults) dataField() []byte     { return p.P.ReadData(1, test_capnp_3) }
func (p TestDefaults) structField() TestAllTypes {
	return TestAllTypes{p.P.ReadStruct(2, test_capnp_4.P)}
}
func (p TestDefaults) enumField() TestEnum { return TestEnum(p.P.ReadStruct16(384) ^ 1) }

/* interface can't have a default */
func (p TestDefaults) interfaceField() TestInterface  { return RemoteTestInterface{p.P.ReadPtr(3)} }
func (p TestDefaults) voidList() []struct{}           { return p.P.ReadVoidList(4, test_capnp_5) }
func (p TestDefaults) boolList() C.Bitset             { return p.P.ReadBitset(5, test_capnp_6) }
func (p TestDefaults) int8List() []int8               { return p.P.ReadI8List(6, test_capnp_7) }
func (p TestDefaults) int16List() []int16             { return p.P.ReadI16List(7, test_capnp_8) }
func (p TestDefaults) int32List() []int32             { return p.P.ReadI32List(8, test_capnp_9) }
func (p TestDefaults) int64List() []int64             { return p.P.ReadI64List(9, test_capnp_10) }
func (p TestDefaults) uInt8List() []uint8             { return p.P.ReadU8List(10, test_capnp_11) }
func (p TestDefaults) uInt16List() []uint16           { return p.P.ReadU16List(11, test_capnp_12) }
func (p TestDefaults) uInt32List() []uint32           { return p.P.ReadU32List(12, test_capnp_13) }
func (p TestDefaults) uInt64List() []uint64           { return p.P.ReadU64List(13, test_capnp_14) }
func (p TestDefaults) float32List() []float32         { return p.P.ReadF32List(14, test_capnp_15) }
func (p TestDefaults) float64List() []float64         { return p.P.ReadF64List(15, test_capnp_16) }
func (p TestDefaults) textList() []C.Pointer          { return p.P.ReadPointerList(16, test_capnp_17) }
func (p TestDefaults) dataList() []C.Pointer          { return p.P.ReadPointerList(17, test_capnp_18) }
func (p TestDefaults) structList() []TestAllTypes     { return ReadTestAllTypesList(p.P, 18, test_capnp_22) }
func (p TestDefaults) enumList() []TestEnum           { return ReadTestEnumList(p.P, 19, test_capnp_23) }
func (p TestDefaults) interfaceList() []TestInterface { return ReadTestInterfaceList(p.P, 20, nil) }
func (p TestDefaults) unionField() TestDefaults_unionField {
	return TestDefaults_unionField(p.P.ReadStruct16(400))
}
func (p TestDefaults) boolUnion() (ret bool) {
	if p.P.ReadStruct16(400) != 36 {
		return true
	}
	return !p.P.ReadStruct1(416)
}
func (p TestDefaults) int8Union() (ret int8) {
	if p.P.ReadStruct16(400) != 37 {
		return -123
	}
	return int8(p.P.ReadStruct8(424)) ^ -123
}
func (p TestDefaults) uint8Union() (ret uint8) {
	if p.P.ReadStruct16(400) != 38 {
		return 124
	}
	return p.P.ReadStruct8(424) ^ 124
}
func (p TestDefaults) int16Union() (ret int16) {
	if p.P.ReadStruct16(400) != 39 {
		return -12345
	}
	return int16(p.P.ReadStruct16(432)) ^ -12345
}
func (p TestDefaults) uint16Union() (ret uint16) {
	if p.P.ReadStruct16(400) != 40 {
		return 12456
	}
	return p.P.ReadStruct16(432) ^ 12456
}
func (p TestDefaults) int32Union() (ret int32) {
	if p.P.ReadStruct16(400) != 41 {
		return -125678
	}
	return int32(p.P.ReadStruct32(448)) ^ -125678
}
func (p TestDefaults) uint32Union() (ret uint32) {
	if p.P.ReadStruct16(400) != 42 {
		return 345786
	}
	return p.P.ReadStruct32(448) ^ 345786
}
func (p TestDefaults) int64Union() (ret int64) {
	if p.P.ReadStruct16(400) != 43 {
		return -123567379234
	}
	return int64(p.P.ReadStruct64(512)) ^ -123567379234
}
func (p TestDefaults) uint64Union() (ret uint64) {
	if p.P.ReadStruct16(400) != 44 {
		return 1235768497284
	}
	return p.P.ReadStruct64(512) ^ 1235768497284
}
func (p TestDefaults) float32Union() (ret float32) {
	if p.P.ReadStruct16(400) != 45 {
		return 33.3
	}
	return p.P.ReadStructF32(512, 33.3)
}
func (p TestDefaults) float64Union() (ret float64) {
	if p.P.ReadStruct16(400) != 46 {
		return 340000
	}
	return p.P.ReadStructF64(512, 340000)
}
func (p TestDefaults) textUnion() (ret string) {
	if p.P.ReadStruct16(400) != 47 {
		return "foo"
	}
	return p.P.ReadString(21, "foo")
}
func (p TestDefaults) dataUnion() (ret []byte) {
	if p.P.ReadStruct16(400) != 48 {
		return test_capnp_24
	}
	return p.P.ReadData(21, test_capnp_24)
}
func (p TestDefaults) structUnion() (ret TestAllTypes) {
	if p.P.ReadStruct16(400) != 49 {
		return test_capnp_25
	}
	return TestAllTypes{p.P.ReadStruct(21, test_capnp_25.P)}
}
func (p TestDefaults) enumUnion() (ret TestEnum) {
	if p.P.ReadStruct16(400) != 50 {
		return TestEnum(1)
	}
	return TestEnum(p.P.ReadStruct16(512) ^ 1)
}
func (p TestDefaults) interfaceUnion() (ret TestInterface) {
	if p.P.ReadStruct16(400) != 51 {
		return
	}
	return RemoteTestInterface{p.P.ReadPtr(21)}
}

func (p TestDefaults) setBoolField(v bool) error   { return p.P.WriteStruct1(0, !v) }
func (p TestDefaults) setInt8Field(v int8) error   { return p.P.WriteStruct8(8, uint8(v^-123)) }
func (p TestDefaults) setInt16Field(v int16) error { return p.P.WriteStruct16(16, uint16(v^-12345)) }
func (p TestDefaults) setInt32Field(v int32) error { return p.P.WriteStruct32(32, uint32(v^-12345678)) }
func (p TestDefaults) setInt64Field(v int64) error {
	return p.P.WriteStruct64(64, uint64(v^-123456789012345))
}
func (p TestDefaults) setUInt8Field(v uint8) error   { return p.P.WriteStruct8(128, v^234) }
func (p TestDefaults) setUInt16Field(v uint16) error { return p.P.WriteStruct16(144, v^45678) }
func (p TestDefaults) setUInt32Field(v uint32) error { return p.P.WriteStruct32(160, v^3456789012) }
func (p TestDefaults) setUInt64Field(v uint64) error {
	return p.P.WriteStruct64(192, v^12345678901234567890)
}
func (p TestDefaults) setFloat32Field(v float32) error { return p.P.WriteStructF32(256, v, 1234.5) }
func (p TestDefaults) setFloat64Field(v float64) error { return p.P.WriteStructF64(320, v, -1.23e+47) }
func (p TestDefaults) setTextField(v string) error     { return p.P.WriteString(0, v, "foo") }
func (p TestDefaults) setDataField(v []byte) error     { return p.P.WriteU8List(1, v, test_capnp_3) }
func (p TestDefaults) setStructField(v TestAllTypes) error {
	return p.P.WriteStruct(2, v.P, test_capnp_4.P)
}
func (p TestDefaults) setEnumField(v TestEnum) error { return p.P.WriteStruct16(384, uint16(v)^1) }

/* interface can't have a default */
func (p TestDefaults) setInterfaceField(v TestInterface) error { return v.MarshalCaptain(p.P, 3) }
func (p TestDefaults) setVoidList(v []struct{}) error          { return p.P.WriteVoidList(4, v, test_capnp_5) }
func (p TestDefaults) setBoolList(v C.Bitset) error            { return p.P.WriteBitset(5, v, test_capnp_6) }
func (p TestDefaults) setInt8List(v []int8) error              { return p.P.WriteI8List(6, v, test_capnp_7) }
func (p TestDefaults) setInt16List(v []int16) error            { return p.P.WriteI16List(7, v, test_capnp_8) }
func (p TestDefaults) setInt32List(v []int32) error            { return p.P.WriteI32List(8, v, test_capnp_9) }
func (p TestDefaults) setInt64List(v []int64) error            { return p.P.WriteI64List(9, v, test_capnp_10) }
func (p TestDefaults) setUInt8List(v []uint8) error            { return p.P.WriteU8List(10, v, test_capnp_11) }
func (p TestDefaults) setUInt16List(v []uint16) error          { return p.P.WriteU16List(11, v, test_capnp_12) }
func (p TestDefaults) setUInt32List(v []uint32) error          { return p.P.WriteU32List(12, v, test_capnp_13) }
func (p TestDefaults) setUInt64List(v []uint64) error          { return p.P.WriteU64List(13, v, test_capnp_14) }
func (p TestDefaults) setFloat32List(v []float32) error        { return p.P.WriteF32List(14, v, test_capnp_15) }
func (p TestDefaults) setFloat64List(v []float64) error        { return p.P.WriteF64List(15, v, test_capnp_16) }
func (p TestDefaults) setTextList(v []C.Pointer) error {
	return p.P.WritePointerList(16, v, test_capnp_17)
}
func (p TestDefaults) setDataList(v []C.Pointer) error {
	return p.P.WritePointerList(17, v, test_capnp_18)
}
func (p TestDefaults) setStructList(v []TestAllTypes) error {
	return WriteTestAllTypesList(p.P, 18, v, test_capnp_22)
}
func (p TestDefaults) setEnumList(v []TestEnum) error {
	return WriteTestEnumList(p.P, 19, v, test_capnp_23)
}
func (p TestDefaults) setInterfaceList(v []TestInterface) error {
	return WriteTestInterfaceList(p.P, 20, v, nil)
}
func (p TestDefaults) setBoolUnion(v bool) error {
	if err := p.P.WriteStruct16(400, 36); err != nil {
		return err
	}
	return p.P.WriteStruct1(416, !v)
}
func (p TestDefaults) setInt8Union(v int8) error {
	if err := p.P.WriteStruct16(400, 37); err != nil {
		return err
	}
	return p.P.WriteStruct8(424, uint8(v^-123))
}
func (p TestDefaults) setUint8Union(v uint8) error {
	if err := p.P.WriteStruct16(400, 38); err != nil {
		return err
	}
	return p.P.WriteStruct8(424, v^124)
}
func (p TestDefaults) setInt16Union(v int16) error {
	if err := p.P.WriteStruct16(400, 39); err != nil {
		return err
	}
	return p.P.WriteStruct16(432, uint16(v^-12345))
}
func (p TestDefaults) setUint16Union(v uint16) error {
	if err := p.P.WriteStruct16(400, 40); err != nil {
		return err
	}
	return p.P.WriteStruct16(432, v^12456)
}
func (p TestDefaults) setInt32Union(v int32) error {
	if err := p.P.WriteStruct16(400, 41); err != nil {
		return err
	}
	return p.P.WriteStruct32(448, uint32(v^-125678))
}
func (p TestDefaults) setUint32Union(v uint32) error {
	if err := p.P.WriteStruct16(400, 42); err != nil {
		return err
	}
	return p.P.WriteStruct32(448, v^345786)
}
func (p TestDefaults) setInt64Union(v int64) error {
	if err := p.P.WriteStruct16(400, 43); err != nil {
		return err
	}
	return p.P.WriteStruct64(512, uint64(v^-123567379234))
}
func (p TestDefaults) setUint64Union(v uint64) error {
	if err := p.P.WriteStruct16(400, 44); err != nil {
		return err
	}
	return p.P.WriteStruct64(512, v^1235768497284)
}
func (p TestDefaults) setFloat32Union(v float32) error {
	if err := p.P.WriteStruct16(400, 45); err != nil {
		return err
	}
	return p.P.WriteStructF32(512, v, 33.3)
}
func (p TestDefaults) setFloat64Union(v float64) error {
	if err := p.P.WriteStruct16(400, 46); err != nil {
		return err
	}
	return p.P.WriteStructF64(512, v, 340000)
}
func (p TestDefaults) setTextUnion(v string) error {
	if err := p.P.WriteStruct16(400, 47); err != nil {
		return err
	}
	return p.P.WriteString(21, v, "foo")
}
func (p TestDefaults) setDataUnion(v []byte) error {
	if err := p.P.WriteStruct16(400, 48); err != nil {
		return err
	}
	return p.P.WriteU8List(21, v, test_capnp_24)
}
func (p TestDefaults) setStructUnion(v TestAllTypes) error {
	if err := p.P.WriteStruct16(400, 49); err != nil {
		return err
	}
	return p.P.WriteStruct(21, v.P, test_capnp_25.P)
}
func (p TestDefaults) setEnumUnion(v TestEnum) error {
	if err := p.P.WriteStruct16(400, 50); err != nil {
		return err
	}
	return p.P.WriteStruct16(512, uint16(v)^1)
}
func (p TestDefaults) setInterfaceUnion(v TestInterface) error {
	if err := p.P.WriteStruct16(400, 51); err != nil {
		return err
	}
	return v.MarshalCaptain(p.P, 21)
}

func ReadTestDefaultsList(p C.Pointer, i int, def []TestDefaults) []TestDefaults {
	if m := p.ReadPtr(i); m.Type() == C.List {
		r := make([]TestDefaults, m.Size())
		for i := range r {
			r[i] = TestDefaults{m.ReadPtr(i)}
		}
		return r
	}
	return def
}
func WriteTestDefaultsList(p C.Pointer, i int, v, def []TestDefaults) (err error) {
	var m C.Pointer
	if !C.SliceEqual(v, def) {
		if m, err = p.Segment.NewPointerList(len(v)); err != nil {
			return err
		}
		for i, u := range v {
			if err := u.MarshalCaptain(m, i); err != nil {
				return err
			}
		}
	}
	return p.WritePtr(i, m)
}

var (
	test_capnp_0  = TestAllTypes{test_capnp.ReadPtr(252)}
	test_capnp_1  = TestAllTypes{test_capnp.ReadPtr(286)}
	test_capnp_2  = TestAllTypes{test_capnp.ReadPtr(320)}
	test_capnp_3  = []byte("bar")
	test_capnp_4  = TestAllTypes{test_capnp.ReadPtr(354)}
	test_capnp_5  = make([]struct{}, 6)
	test_capnp_6  = test_capnp.ReadBitset(591, C.Bitset{})
	test_capnp_7  = []int8{111, -111}
	test_capnp_8  = []int16{11111, -11111}
	test_capnp_9  = []int32{111111111, -111111111}
	test_capnp_10 = []int64{1111111111111111111, -1111111111111111111}
	test_capnp_11 = []uint8{111, 222}
	test_capnp_12 = []uint16{33333, 44444}
	test_capnp_13 = []uint32{3333333333}
	test_capnp_14 = []uint64{11111111111111111111}
	test_capnp_15 = []float32{5555.5, 2222.25}
	test_capnp_16 = []float64{7777.75, 1111.125}
	test_capnp_17 = []C.Pointer{test_capnp.ReadPtr(593), test_capnp.ReadPtr(595), test_capnp.ReadPtr(597)}
	test_capnp_18 = []C.Pointer{test_capnp.ReadPtr(599), test_capnp.ReadPtr(601), test_capnp.ReadPtr(604)}
	test_capnp_19 = TestAllTypes{test_capnp.ReadPtr(606)}
	test_capnp_20 = TestAllTypes{test_capnp.ReadPtr(640)}
	test_capnp_21 = TestAllTypes{test_capnp.ReadPtr(674)}
	test_capnp_22 = []TestAllTypes{test_capnp_19, test_capnp_20, test_capnp_21}
	test_capnp_23 = []TestEnum{TestEnum(1), TestEnum(2)}
	test_capnp_24 = []byte("bar")
	test_capnp_25 = TestAllTypes{test_capnp.ReadPtr(708)}
	test_capnp    = C.NewBuffer([]byte{
		0, 0, 0, 0, 9, 0, 22, 0,
		1, 244, 128, 13, 14, 16, 76, 251,
		78, 115, 232, 56, 166, 51, 0, 0,
		90, 0, 210, 4, 20, 136, 98, 3,
		210, 10, 111, 18, 33, 25, 204, 4,
		95, 112, 9, 175, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 144, 117, 64,
		1, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		85, 0, 0, 0, 34, 0, 0, 0,
		85, 0, 0, 0, 26, 0, 0, 0,
		84, 0, 0, 0, 9, 0, 22, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		81, 1, 0, 0, 24, 0, 0, 0,
		77, 1, 0, 0, 41, 0, 0, 0,
		77, 1, 0, 0, 34, 0, 0, 0,
		77, 1, 0, 0, 35, 0, 0, 0,
		77, 1, 0, 0, 36, 0, 0, 0,
		81, 1, 0, 0, 37, 0, 0, 0,
		93, 1, 0, 0, 34, 0, 0, 0,
		93, 1, 0, 0, 35, 0, 0, 0,
		93, 1, 0, 0, 36, 0, 0, 0,
		97, 1, 0, 0, 37, 0, 0, 0,
		109, 1, 0, 0, 52, 0, 0, 0,
		117, 1, 0, 0, 53, 0, 0, 0,
		137, 1, 0, 0, 30, 0, 0, 0,
		157, 1, 0, 0, 30, 0, 0, 0,
		177, 1, 0, 0, 31, 0, 0, 0,
		57, 3, 0, 0, 19, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		98, 97, 122, 0, 0, 0, 0, 0,
		113, 117, 120, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		85, 0, 0, 0, 58, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		80, 0, 0, 0, 9, 0, 22, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		110, 101, 115, 116, 101, 100, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		85, 0, 0, 0, 114, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		114, 101, 97, 108, 108, 121, 32, 110,
		101, 115, 116, 101, 100, 0, 0, 0,
		26, 0, 0, 0, 0, 0, 0, 0,
		12, 222, 128, 127, 0, 0, 0, 0,
		210, 4, 210, 233, 0, 128, 255, 127,
		78, 97, 188, 0, 64, 211, 160, 250,
		0, 0, 0, 248, 255, 255, 255, 7,
		121, 223, 13, 134, 72, 112, 0, 0,
		46, 117, 19, 253, 138, 150, 253, 255,
		0, 0, 0, 0, 0, 0, 0, 128,
		255, 255, 255, 255, 255, 255, 255, 127,
		12, 34, 0, 255, 0, 0, 0, 0,
		210, 4, 46, 22, 0, 0, 255, 255,
		78, 97, 188, 0, 192, 44, 95, 5,
		0, 0, 0, 0, 255, 255, 255, 255,
		121, 223, 13, 134, 72, 112, 0, 0,
		210, 138, 236, 2, 117, 105, 2, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		255, 255, 255, 255, 255, 255, 255, 255,
		0, 0, 0, 0, 56, 180, 150, 73,
		194, 189, 240, 124, 194, 189, 240, 252,
		234, 28, 8, 2, 234, 28, 8, 130,
		0, 0, 0, 0, 0, 0, 0, 0,
		64, 222, 119, 131, 33, 18, 220, 66,
		41, 144, 35, 202, 229, 200, 118, 127,
		41, 144, 35, 202, 229, 200, 118, 255,
		145, 247, 80, 55, 158, 120, 102, 0,
		145, 247, 80, 55, 158, 120, 102, 128,
		9, 0, 0, 0, 42, 0, 0, 0,
		9, 0, 0, 0, 50, 0, 0, 0,
		9, 0, 0, 0, 58, 0, 0, 0,
		113, 117, 117, 120, 0, 0, 0, 0,
		99, 111, 114, 103, 101, 0, 0, 0,
		103, 114, 97, 117, 108, 116, 0, 0,
		9, 0, 0, 0, 50, 0, 0, 0,
		9, 0, 0, 0, 42, 0, 0, 0,
		9, 0, 0, 0, 34, 0, 0, 0,
		103, 97, 114, 112, 108, 121, 0, 0,
		119, 97, 108, 100, 111, 0, 0, 0,
		102, 114, 101, 100, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		77, 1, 0, 0, 122, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		217, 0, 0, 0, 122, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		101, 0, 0, 0, 122, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		120, 32, 115, 116, 114, 117, 99, 116,
		108, 105, 115, 116, 32, 49, 0, 0,
		120, 32, 115, 116, 114, 117, 99, 116,
		108, 105, 115, 116, 32, 50, 0, 0,
		120, 32, 115, 116, 114, 117, 99, 116,
		108, 105, 115, 116, 32, 51, 0, 0,
		2, 0, 1, 0, 0, 0, 0, 0,
		1, 0, 0, 0, 33, 0, 0, 0,
		9, 0, 0, 0, 0, 0, 0, 0,
		1, 0, 0, 0, 50, 0, 0, 0,
		112, 108, 117, 103, 104, 0, 0, 0,
		1, 0, 0, 0, 50, 0, 0, 0,
		120, 121, 122, 122, 121, 0, 0, 0,
		1, 0, 0, 0, 42, 0, 0, 0,
		116, 104, 117, 100, 0, 0, 0, 0,
		1, 0, 0, 0, 34, 0, 0, 0,
		111, 111, 112, 115, 0, 0, 0, 0,
		1, 0, 0, 0, 74, 0, 0, 0,
		101, 120, 104, 97, 117, 115, 116, 101,
		100, 0, 0, 0, 0, 0, 0, 0,
		1, 0, 0, 0, 58, 0, 0, 0,
		114, 102, 99, 51, 48, 57, 50, 0,
		0, 0, 0, 0, 9, 0, 22, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		85, 0, 0, 0, 106, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		115, 116, 114, 117, 99, 116, 108, 105,
		115, 116, 32, 49, 0, 0, 0, 0,
		0, 0, 0, 0, 9, 0, 22, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		85, 0, 0, 0, 106, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		115, 116, 114, 117, 99, 116, 108, 105,
		115, 116, 32, 50, 0, 0, 0, 0,
		0, 0, 0, 0, 9, 0, 22, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		85, 0, 0, 0, 106, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		115, 116, 114, 117, 99, 116, 108, 105,
		115, 116, 32, 51, 0, 0, 0, 0,
		0, 0, 0, 0, 9, 0, 22, 0,
		1, 244, 128, 13, 14, 16, 76, 251,
		78, 115, 232, 56, 166, 51, 0, 0,
		90, 0, 210, 4, 20, 136, 98, 3,
		210, 10, 111, 18, 33, 25, 204, 4,
		95, 112, 9, 175, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 144, 117, 64,
		1, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		85, 0, 0, 0, 34, 0, 0, 0,
		85, 0, 0, 0, 26, 0, 0, 0,
		84, 0, 0, 0, 9, 0, 22, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		81, 1, 0, 0, 24, 0, 0, 0,
		77, 1, 0, 0, 41, 0, 0, 0,
		77, 1, 0, 0, 34, 0, 0, 0,
		77, 1, 0, 0, 35, 0, 0, 0,
		77, 1, 0, 0, 36, 0, 0, 0,
		81, 1, 0, 0, 37, 0, 0, 0,
		93, 1, 0, 0, 34, 0, 0, 0,
		93, 1, 0, 0, 35, 0, 0, 0,
		93, 1, 0, 0, 36, 0, 0, 0,
		97, 1, 0, 0, 37, 0, 0, 0,
		109, 1, 0, 0, 52, 0, 0, 0,
		117, 1, 0, 0, 53, 0, 0, 0,
		137, 1, 0, 0, 30, 0, 0, 0,
		157, 1, 0, 0, 30, 0, 0, 0,
		177, 1, 0, 0, 31, 0, 0, 0,
		57, 3, 0, 0, 19, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		98, 97, 122, 0, 0, 0, 0, 0,
		113, 117, 120, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		85, 0, 0, 0, 58, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		80, 0, 0, 0, 9, 0, 22, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		110, 101, 115, 116, 101, 100, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		85, 0, 0, 0, 114, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		114, 101, 97, 108, 108, 121, 32, 110,
		101, 115, 116, 101, 100, 0, 0, 0,
		26, 0, 0, 0, 0, 0, 0, 0,
		12, 222, 128, 127, 0, 0, 0, 0,
		210, 4, 210, 233, 0, 128, 255, 127,
		78, 97, 188, 0, 64, 211, 160, 250,
		0, 0, 0, 248, 255, 255, 255, 7,
		121, 223, 13, 134, 72, 112, 0, 0,
		46, 117, 19, 253, 138, 150, 253, 255,
		0, 0, 0, 0, 0, 0, 0, 128,
		255, 255, 255, 255, 255, 255, 255, 127,
		12, 34, 0, 255, 0, 0, 0, 0,
		210, 4, 46, 22, 0, 0, 255, 255,
		78, 97, 188, 0, 192, 44, 95, 5,
		0, 0, 0, 0, 255, 255, 255, 255,
		121, 223, 13, 134, 72, 112, 0, 0,
		210, 138, 236, 2, 117, 105, 2, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		255, 255, 255, 255, 255, 255, 255, 255,
		0, 0, 0, 0, 56, 180, 150, 73,
		194, 189, 240, 124, 194, 189, 240, 252,
		234, 28, 8, 2, 234, 28, 8, 130,
		0, 0, 0, 0, 0, 0, 0, 0,
		64, 222, 119, 131, 33, 18, 220, 66,
		41, 144, 35, 202, 229, 200, 118, 127,
		41, 144, 35, 202, 229, 200, 118, 255,
		145, 247, 80, 55, 158, 120, 102, 0,
		145, 247, 80, 55, 158, 120, 102, 128,
		9, 0, 0, 0, 42, 0, 0, 0,
		9, 0, 0, 0, 50, 0, 0, 0,
		9, 0, 0, 0, 58, 0, 0, 0,
		113, 117, 117, 120, 0, 0, 0, 0,
		99, 111, 114, 103, 101, 0, 0, 0,
		103, 114, 97, 117, 108, 116, 0, 0,
		9, 0, 0, 0, 50, 0, 0, 0,
		9, 0, 0, 0, 42, 0, 0, 0,
		9, 0, 0, 0, 34, 0, 0, 0,
		103, 97, 114, 112, 108, 121, 0, 0,
		119, 97, 108, 100, 111, 0, 0, 0,
		102, 114, 101, 100, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		77, 1, 0, 0, 122, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		217, 0, 0, 0, 122, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		101, 0, 0, 0, 122, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		120, 32, 115, 116, 114, 117, 99, 116,
		108, 105, 115, 116, 32, 49, 0, 0,
		120, 32, 115, 116, 114, 117, 99, 116,
		108, 105, 115, 116, 32, 50, 0, 0,
		120, 32, 115, 116, 114, 117, 99, 116,
		108, 105, 115, 116, 32, 51, 0, 0,
		2, 0, 1, 0, 0, 0, 0, 0,
		1, 0, 0, 0, 33, 0, 0, 0,
		9, 0, 0, 0, 0, 0, 0, 0,
		1, 0, 0, 0, 50, 0, 0, 0,
		112, 108, 117, 103, 104, 0, 0, 0,
		1, 0, 0, 0, 50, 0, 0, 0,
		120, 121, 122, 122, 121, 0, 0, 0,
		1, 0, 0, 0, 42, 0, 0, 0,
		116, 104, 117, 100, 0, 0, 0, 0,
		1, 0, 0, 0, 34, 0, 0, 0,
		111, 111, 112, 115, 0, 0, 0, 0,
		1, 0, 0, 0, 74, 0, 0, 0,
		101, 120, 104, 97, 117, 115, 116, 101,
		100, 0, 0, 0, 0, 0, 0, 0,
		1, 0, 0, 0, 58, 0, 0, 0,
		114, 102, 99, 51, 48, 57, 50, 0,
		0, 0, 0, 0, 9, 0, 22, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		85, 0, 0, 0, 106, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		115, 116, 114, 117, 99, 116, 108, 105,
		115, 116, 32, 49, 0, 0, 0, 0,
		0, 0, 0, 0, 9, 0, 22, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		85, 0, 0, 0, 106, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		115, 116, 114, 117, 99, 116, 108, 105,
		115, 116, 32, 50, 0, 0, 0, 0,
		0, 0, 0, 0, 9, 0, 22, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		85, 0, 0, 0, 106, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		115, 116, 114, 117, 99, 116, 108, 105,
		115, 116, 32, 51, 0, 0, 0, 0,
		0, 0, 0, 0, 9, 0, 22, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 133, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0,
	}).Root()
)
