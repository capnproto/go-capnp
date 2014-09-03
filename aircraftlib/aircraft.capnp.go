package aircraftlib

// AUTO GENERATED - DO NOT EDIT

import (
	C "github.com/glycerine/go-capnproto"
	"math"
	"net"
	"unsafe"
)

type Zdate C.Struct

func NewZdate(s *C.Segment) Zdate      { return Zdate(s.NewStruct(8, 0)) }
func NewRootZdate(s *C.Segment) Zdate  { return Zdate(s.NewRootStruct(8, 0)) }
func ReadRootZdate(s *C.Segment) Zdate { return Zdate(s.Root(0).ToStruct()) }
func (s Zdate) Year() int16            { return int16(C.Struct(s).Get16(0)) }
func (s Zdate) SetYear(v int16)        { C.Struct(s).Set16(0, uint16(v)) }
func (s Zdate) Month() uint8           { return C.Struct(s).Get8(2) }
func (s Zdate) SetMonth(v uint8)       { C.Struct(s).Set8(2, v) }
func (s Zdate) Day() uint8             { return C.Struct(s).Get8(3) }
func (s Zdate) SetDay(v uint8)         { C.Struct(s).Set8(3, v) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s Zdate) MarshalJSON() (bs []byte, err error) { return }

type Zdate_List C.PointerList

func NewZdateList(s *C.Segment, sz int) Zdate_List { return Zdate_List(s.NewUInt32List(sz)) }
func (s Zdate_List) Len() int                      { return C.PointerList(s).Len() }
func (s Zdate_List) At(i int) Zdate                { return Zdate(C.PointerList(s).At(i).ToStruct()) }
func (s Zdate_List) ToArray() []Zdate              { return *(*[]Zdate)(unsafe.Pointer(C.PointerList(s).ToArray())) }

type Zdata C.Struct

func NewZdata(s *C.Segment) Zdata      { return Zdata(s.NewStruct(0, 1)) }
func NewRootZdata(s *C.Segment) Zdata  { return Zdata(s.NewRootStruct(0, 1)) }
func ReadRootZdata(s *C.Segment) Zdata { return Zdata(s.Root(0).ToStruct()) }
func (s Zdata) Data() []byte           { return C.Struct(s).GetObject(0).ToData() }
func (s Zdata) SetData(v []byte)       { C.Struct(s).SetObject(0, s.Segment.NewData(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s Zdata) MarshalJSON() (bs []byte, err error) { return }

type Zdata_List C.PointerList

func NewZdataList(s *C.Segment, sz int) Zdata_List { return Zdata_List(s.NewCompositeList(0, 1, sz)) }
func (s Zdata_List) Len() int                      { return C.PointerList(s).Len() }
func (s Zdata_List) At(i int) Zdata                { return Zdata(C.PointerList(s).At(i).ToStruct()) }
func (s Zdata_List) ToArray() []Zdata              { return *(*[]Zdata)(unsafe.Pointer(C.PointerList(s).ToArray())) }

type Airport uint16

const (
	AIRPORT_NONE Airport = 0
	AIRPORT_JFK  Airport = 1
	AIRPORT_LAX  Airport = 2
	AIRPORT_SFO  Airport = 3
	AIRPORT_LUV  Airport = 4
	AIRPORT_DFW  Airport = 5
	AIRPORT_TEST Airport = 6
)

func (c Airport) String() string {
	switch c {
	case AIRPORT_NONE:
		return "none"
	case AIRPORT_JFK:
		return "jfk"
	case AIRPORT_LAX:
		return "lax"
	case AIRPORT_SFO:
		return "sfo"
	case AIRPORT_LUV:
		return "luv"
	case AIRPORT_DFW:
		return "dfw"
	case AIRPORT_TEST:
		return "test"
	default:
		return ""
	}
}

type Airport_List C.PointerList

func NewAirportList(s *C.Segment, sz int) Airport_List { return Airport_List(s.NewUInt16List(sz)) }
func (s Airport_List) Len() int                        { return C.UInt16List(s).Len() }
func (s Airport_List) At(i int) Airport                { return Airport(C.UInt16List(s).At(i)) }
func (s Airport_List) ToArray() []Airport {
	return *(*[]Airport)(unsafe.Pointer(C.UInt16List(s).ToEnumArray()))
}

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s Airport) MarshalJSON() (bs []byte, err error) { return }

type PlaneBase C.Struct

func NewPlaneBase(s *C.Segment) PlaneBase      { return PlaneBase(s.NewStruct(32, 2)) }
func NewRootPlaneBase(s *C.Segment) PlaneBase  { return PlaneBase(s.NewRootStruct(32, 2)) }
func ReadRootPlaneBase(s *C.Segment) PlaneBase { return PlaneBase(s.Root(0).ToStruct()) }
func (s PlaneBase) Name() string               { return C.Struct(s).GetObject(0).ToText() }
func (s PlaneBase) SetName(v string)           { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s PlaneBase) Homes() Airport_List        { return Airport_List(C.Struct(s).GetObject(1)) }
func (s PlaneBase) SetHomes(v Airport_List)    { C.Struct(s).SetObject(1, C.Object(v)) }
func (s PlaneBase) Rating() int64              { return int64(C.Struct(s).Get64(0)) }
func (s PlaneBase) SetRating(v int64)          { C.Struct(s).Set64(0, uint64(v)) }
func (s PlaneBase) CanFly() bool               { return C.Struct(s).Get1(64) }
func (s PlaneBase) SetCanFly(v bool)           { C.Struct(s).Set1(64, v) }
func (s PlaneBase) Capacity() int64            { return int64(C.Struct(s).Get64(16)) }
func (s PlaneBase) SetCapacity(v int64)        { C.Struct(s).Set64(16, uint64(v)) }
func (s PlaneBase) MaxSpeed() float64          { return math.Float64frombits(C.Struct(s).Get64(24)) }
func (s PlaneBase) SetMaxSpeed(v float64)      { C.Struct(s).Set64(24, math.Float64bits(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s PlaneBase) MarshalJSON() (bs []byte, err error) { return }

type PlaneBase_List C.PointerList

func NewPlaneBaseList(s *C.Segment, sz int) PlaneBase_List {
	return PlaneBase_List(s.NewCompositeList(32, 2, sz))
}
func (s PlaneBase_List) Len() int           { return C.PointerList(s).Len() }
func (s PlaneBase_List) At(i int) PlaneBase { return PlaneBase(C.PointerList(s).At(i).ToStruct()) }
func (s PlaneBase_List) ToArray() []PlaneBase {
	return *(*[]PlaneBase)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type B737 C.Struct

func NewB737(s *C.Segment) B737      { return B737(s.NewStruct(0, 1)) }
func NewRootB737(s *C.Segment) B737  { return B737(s.NewRootStruct(0, 1)) }
func ReadRootB737(s *C.Segment) B737 { return B737(s.Root(0).ToStruct()) }
func (s B737) Base() PlaneBase       { return PlaneBase(C.Struct(s).GetObject(0).ToStruct()) }
func (s B737) SetBase(v PlaneBase)   { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s B737) MarshalJSON() (bs []byte, err error) { return }

type B737_List C.PointerList

func NewB737List(s *C.Segment, sz int) B737_List { return B737_List(s.NewCompositeList(0, 1, sz)) }
func (s B737_List) Len() int                     { return C.PointerList(s).Len() }
func (s B737_List) At(i int) B737                { return B737(C.PointerList(s).At(i).ToStruct()) }
func (s B737_List) ToArray() []B737              { return *(*[]B737)(unsafe.Pointer(C.PointerList(s).ToArray())) }

type A320 C.Struct

func NewA320(s *C.Segment) A320      { return A320(s.NewStruct(0, 1)) }
func NewRootA320(s *C.Segment) A320  { return A320(s.NewRootStruct(0, 1)) }
func ReadRootA320(s *C.Segment) A320 { return A320(s.Root(0).ToStruct()) }
func (s A320) Base() PlaneBase       { return PlaneBase(C.Struct(s).GetObject(0).ToStruct()) }
func (s A320) SetBase(v PlaneBase)   { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s A320) MarshalJSON() (bs []byte, err error) { return }

type A320_List C.PointerList

func NewA320List(s *C.Segment, sz int) A320_List { return A320_List(s.NewCompositeList(0, 1, sz)) }
func (s A320_List) Len() int                     { return C.PointerList(s).Len() }
func (s A320_List) At(i int) A320                { return A320(C.PointerList(s).At(i).ToStruct()) }
func (s A320_List) ToArray() []A320              { return *(*[]A320)(unsafe.Pointer(C.PointerList(s).ToArray())) }

type F16 C.Struct

func NewF16(s *C.Segment) F16      { return F16(s.NewStruct(0, 1)) }
func NewRootF16(s *C.Segment) F16  { return F16(s.NewRootStruct(0, 1)) }
func ReadRootF16(s *C.Segment) F16 { return F16(s.Root(0).ToStruct()) }
func (s F16) Base() PlaneBase      { return PlaneBase(C.Struct(s).GetObject(0).ToStruct()) }
func (s F16) SetBase(v PlaneBase)  { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s F16) MarshalJSON() (bs []byte, err error) { return }

type F16_List C.PointerList

func NewF16List(s *C.Segment, sz int) F16_List { return F16_List(s.NewCompositeList(0, 1, sz)) }
func (s F16_List) Len() int                    { return C.PointerList(s).Len() }
func (s F16_List) At(i int) F16                { return F16(C.PointerList(s).At(i).ToStruct()) }
func (s F16_List) ToArray() []F16              { return *(*[]F16)(unsafe.Pointer(C.PointerList(s).ToArray())) }

type Regression C.Struct

func NewRegression(s *C.Segment) Regression      { return Regression(s.NewStruct(24, 3)) }
func NewRootRegression(s *C.Segment) Regression  { return Regression(s.NewRootStruct(24, 3)) }
func ReadRootRegression(s *C.Segment) Regression { return Regression(s.Root(0).ToStruct()) }
func (s Regression) Base() PlaneBase             { return PlaneBase(C.Struct(s).GetObject(0).ToStruct()) }
func (s Regression) SetBase(v PlaneBase)         { C.Struct(s).SetObject(0, C.Object(v)) }
func (s Regression) B0() float64                 { return math.Float64frombits(C.Struct(s).Get64(0)) }
func (s Regression) SetB0(v float64)             { C.Struct(s).Set64(0, math.Float64bits(v)) }
func (s Regression) Beta() C.Float64List         { return C.Float64List(C.Struct(s).GetObject(1)) }
func (s Regression) SetBeta(v C.Float64List)     { C.Struct(s).SetObject(1, C.Object(v)) }
func (s Regression) Planes() Aircraft_List       { return Aircraft_List(C.Struct(s).GetObject(2)) }
func (s Regression) SetPlanes(v Aircraft_List)   { C.Struct(s).SetObject(2, C.Object(v)) }
func (s Regression) Ymu() float64                { return math.Float64frombits(C.Struct(s).Get64(8)) }
func (s Regression) SetYmu(v float64)            { C.Struct(s).Set64(8, math.Float64bits(v)) }
func (s Regression) Ysd() float64                { return math.Float64frombits(C.Struct(s).Get64(16)) }
func (s Regression) SetYsd(v float64)            { C.Struct(s).Set64(16, math.Float64bits(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s Regression) MarshalJSON() (bs []byte, err error) { return }

type Regression_List C.PointerList

func NewRegressionList(s *C.Segment, sz int) Regression_List {
	return Regression_List(s.NewCompositeList(24, 3, sz))
}
func (s Regression_List) Len() int            { return C.PointerList(s).Len() }
func (s Regression_List) At(i int) Regression { return Regression(C.PointerList(s).At(i).ToStruct()) }
func (s Regression_List) ToArray() []Regression {
	return *(*[]Regression)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type Aircraft C.Struct
type Aircraft_Which uint16

const (
	AIRCRAFT_VOID Aircraft_Which = 0
	AIRCRAFT_B737 Aircraft_Which = 1
	AIRCRAFT_A320 Aircraft_Which = 2
	AIRCRAFT_F16  Aircraft_Which = 3
)

func NewAircraft(s *C.Segment) Aircraft      { return Aircraft(s.NewStruct(8, 1)) }
func NewRootAircraft(s *C.Segment) Aircraft  { return Aircraft(s.NewRootStruct(8, 1)) }
func ReadRootAircraft(s *C.Segment) Aircraft { return Aircraft(s.Root(0).ToStruct()) }
func (s Aircraft) Which() Aircraft_Which     { return Aircraft_Which(C.Struct(s).Get16(0)) }
func (s Aircraft) SetVoid()                  { C.Struct(s).Set16(0, 0) }
func (s Aircraft) B737() B737                { return B737(C.Struct(s).GetObject(0).ToStruct()) }
func (s Aircraft) SetB737(v B737)            { C.Struct(s).Set16(0, 1); C.Struct(s).SetObject(0, C.Object(v)) }
func (s Aircraft) A320() A320                { return A320(C.Struct(s).GetObject(0).ToStruct()) }
func (s Aircraft) SetA320(v A320)            { C.Struct(s).Set16(0, 2); C.Struct(s).SetObject(0, C.Object(v)) }
func (s Aircraft) F16() F16                  { return F16(C.Struct(s).GetObject(0).ToStruct()) }
func (s Aircraft) SetF16(v F16)              { C.Struct(s).Set16(0, 3); C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s Aircraft) MarshalJSON() (bs []byte, err error) { return }

type Aircraft_List C.PointerList

func NewAircraftList(s *C.Segment, sz int) Aircraft_List {
	return Aircraft_List(s.NewCompositeList(8, 1, sz))
}
func (s Aircraft_List) Len() int          { return C.PointerList(s).Len() }
func (s Aircraft_List) At(i int) Aircraft { return Aircraft(C.PointerList(s).At(i).ToStruct()) }
func (s Aircraft_List) ToArray() []Aircraft {
	return *(*[]Aircraft)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type Z C.Struct
type Z_Which uint16

const (
	Z_VOID        Z_Which = 0
	Z_ZZ          Z_Which = 1
	Z_F64         Z_Which = 2
	Z_F32         Z_Which = 3
	Z_I64         Z_Which = 4
	Z_I32         Z_Which = 5
	Z_I16         Z_Which = 6
	Z_I8          Z_Which = 7
	Z_U64         Z_Which = 8
	Z_U32         Z_Which = 9
	Z_U16         Z_Which = 10
	Z_U8          Z_Which = 11
	Z_BOOL        Z_Which = 12
	Z_TEXT        Z_Which = 13
	Z_BLOB        Z_Which = 14
	Z_F64VEC      Z_Which = 15
	Z_F32VEC      Z_Which = 16
	Z_I64VEC      Z_Which = 17
	Z_I32VEC      Z_Which = 18
	Z_I16VEC      Z_Which = 19
	Z_I8VEC       Z_Which = 20
	Z_U64VEC      Z_Which = 21
	Z_U32VEC      Z_Which = 22
	Z_U16VEC      Z_Which = 23
	Z_U8VEC       Z_Which = 24
	Z_ZVEC        Z_Which = 25
	Z_ZVECVEC     Z_Which = 26
	Z_ZDATE       Z_Which = 27
	Z_ZDATA       Z_Which = 28
	Z_AIRCRAFTVEC Z_Which = 29
	Z_AIRCRAFT    Z_Which = 30
	Z_REGRESSION  Z_Which = 31
	Z_PLANEBASE   Z_Which = 32
	Z_AIRPORT     Z_Which = 33
	Z_B737        Z_Which = 34
	Z_A320        Z_Which = 35
	Z_F16         Z_Which = 36
	Z_ZDATEVEC    Z_Which = 37
	Z_ZDATAVEC    Z_Which = 38
	Z_BOOLVEC     Z_Which = 39
)

func NewZ(s *C.Segment) Z             { return Z(s.NewStruct(16, 1)) }
func NewRootZ(s *C.Segment) Z         { return Z(s.NewRootStruct(16, 1)) }
func ReadRootZ(s *C.Segment) Z        { return Z(s.Root(0).ToStruct()) }
func (s Z) Which() Z_Which            { return Z_Which(C.Struct(s).Get16(0)) }
func (s Z) SetVoid()                  { C.Struct(s).Set16(0, 0) }
func (s Z) Zz() Z                     { return Z(C.Struct(s).GetObject(0).ToStruct()) }
func (s Z) SetZz(v Z)                 { C.Struct(s).Set16(0, 1); C.Struct(s).SetObject(0, C.Object(v)) }
func (s Z) F64() float64              { return math.Float64frombits(C.Struct(s).Get64(8)) }
func (s Z) SetF64(v float64)          { C.Struct(s).Set16(0, 2); C.Struct(s).Set64(8, math.Float64bits(v)) }
func (s Z) F32() float32              { return math.Float32frombits(C.Struct(s).Get32(8)) }
func (s Z) SetF32(v float32)          { C.Struct(s).Set16(0, 3); C.Struct(s).Set32(8, math.Float32bits(v)) }
func (s Z) I64() int64                { return int64(C.Struct(s).Get64(8)) }
func (s Z) SetI64(v int64)            { C.Struct(s).Set16(0, 4); C.Struct(s).Set64(8, uint64(v)) }
func (s Z) I32() int32                { return int32(C.Struct(s).Get32(8)) }
func (s Z) SetI32(v int32)            { C.Struct(s).Set16(0, 5); C.Struct(s).Set32(8, uint32(v)) }
func (s Z) I16() int16                { return int16(C.Struct(s).Get16(8)) }
func (s Z) SetI16(v int16)            { C.Struct(s).Set16(0, 6); C.Struct(s).Set16(8, uint16(v)) }
func (s Z) I8() int8                  { return int8(C.Struct(s).Get8(8)) }
func (s Z) SetI8(v int8)              { C.Struct(s).Set16(0, 7); C.Struct(s).Set8(8, uint8(v)) }
func (s Z) U64() uint64               { return C.Struct(s).Get64(8) }
func (s Z) SetU64(v uint64)           { C.Struct(s).Set16(0, 8); C.Struct(s).Set64(8, v) }
func (s Z) U32() uint32               { return C.Struct(s).Get32(8) }
func (s Z) SetU32(v uint32)           { C.Struct(s).Set16(0, 9); C.Struct(s).Set32(8, v) }
func (s Z) U16() uint16               { return C.Struct(s).Get16(8) }
func (s Z) SetU16(v uint16)           { C.Struct(s).Set16(0, 10); C.Struct(s).Set16(8, v) }
func (s Z) U8() uint8                 { return C.Struct(s).Get8(8) }
func (s Z) SetU8(v uint8)             { C.Struct(s).Set16(0, 11); C.Struct(s).Set8(8, v) }
func (s Z) Bool() bool                { return C.Struct(s).Get1(64) }
func (s Z) SetBool(v bool)            { C.Struct(s).Set16(0, 12); C.Struct(s).Set1(64, v) }
func (s Z) Text() string              { return C.Struct(s).GetObject(0).ToText() }
func (s Z) SetText(v string)          { C.Struct(s).Set16(0, 13); C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s Z) Blob() []byte              { return C.Struct(s).GetObject(0).ToData() }
func (s Z) SetBlob(v []byte)          { C.Struct(s).Set16(0, 14); C.Struct(s).SetObject(0, s.Segment.NewData(v)) }
func (s Z) F64vec() C.Float64List     { return C.Float64List(C.Struct(s).GetObject(0)) }
func (s Z) SetF64vec(v C.Float64List) { C.Struct(s).Set16(0, 15); C.Struct(s).SetObject(0, C.Object(v)) }
func (s Z) F32vec() C.Float32List     { return C.Float32List(C.Struct(s).GetObject(0)) }
func (s Z) SetF32vec(v C.Float32List) { C.Struct(s).Set16(0, 16); C.Struct(s).SetObject(0, C.Object(v)) }
func (s Z) I64vec() C.Int64List       { return C.Int64List(C.Struct(s).GetObject(0)) }
func (s Z) SetI64vec(v C.Int64List)   { C.Struct(s).Set16(0, 17); C.Struct(s).SetObject(0, C.Object(v)) }
func (s Z) I32vec() C.Int32List       { return C.Int32List(C.Struct(s).GetObject(0)) }
func (s Z) SetI32vec(v C.Int32List)   { C.Struct(s).Set16(0, 18); C.Struct(s).SetObject(0, C.Object(v)) }
func (s Z) I16vec() C.Int16List       { return C.Int16List(C.Struct(s).GetObject(0)) }
func (s Z) SetI16vec(v C.Int16List)   { C.Struct(s).Set16(0, 19); C.Struct(s).SetObject(0, C.Object(v)) }
func (s Z) I8vec() C.Int8List         { return C.Int8List(C.Struct(s).GetObject(0)) }
func (s Z) SetI8vec(v C.Int8List)     { C.Struct(s).Set16(0, 20); C.Struct(s).SetObject(0, C.Object(v)) }
func (s Z) U64vec() C.UInt64List      { return C.UInt64List(C.Struct(s).GetObject(0)) }
func (s Z) SetU64vec(v C.UInt64List)  { C.Struct(s).Set16(0, 21); C.Struct(s).SetObject(0, C.Object(v)) }
func (s Z) U32vec() C.UInt32List      { return C.UInt32List(C.Struct(s).GetObject(0)) }
func (s Z) SetU32vec(v C.UInt32List)  { C.Struct(s).Set16(0, 22); C.Struct(s).SetObject(0, C.Object(v)) }
func (s Z) U16vec() C.UInt16List      { return C.UInt16List(C.Struct(s).GetObject(0)) }
func (s Z) SetU16vec(v C.UInt16List)  { C.Struct(s).Set16(0, 23); C.Struct(s).SetObject(0, C.Object(v)) }
func (s Z) U8vec() C.UInt8List        { return C.UInt8List(C.Struct(s).GetObject(0)) }
func (s Z) SetU8vec(v C.UInt8List)    { C.Struct(s).Set16(0, 24); C.Struct(s).SetObject(0, C.Object(v)) }
func (s Z) Zvec() Z_List              { return Z_List(C.Struct(s).GetObject(0)) }
func (s Z) SetZvec(v Z_List)          { C.Struct(s).Set16(0, 25); C.Struct(s).SetObject(0, C.Object(v)) }
func (s Z) Zvecvec() C.PointerList    { return C.PointerList(C.Struct(s).GetObject(0)) }
func (s Z) SetZvecvec(v C.PointerList) {
	C.Struct(s).Set16(0, 26)
	C.Struct(s).SetObject(0, C.Object(v))
}
func (s Z) Zdate() Zdate               { return Zdate(C.Struct(s).GetObject(0).ToStruct()) }
func (s Z) SetZdate(v Zdate)           { C.Struct(s).Set16(0, 27); C.Struct(s).SetObject(0, C.Object(v)) }
func (s Z) Zdata() Zdata               { return Zdata(C.Struct(s).GetObject(0).ToStruct()) }
func (s Z) SetZdata(v Zdata)           { C.Struct(s).Set16(0, 28); C.Struct(s).SetObject(0, C.Object(v)) }
func (s Z) Aircraftvec() Aircraft_List { return Aircraft_List(C.Struct(s).GetObject(0)) }
func (s Z) SetAircraftvec(v Aircraft_List) {
	C.Struct(s).Set16(0, 29)
	C.Struct(s).SetObject(0, C.Object(v))
}
func (s Z) Aircraft() Aircraft     { return Aircraft(C.Struct(s).GetObject(0).ToStruct()) }
func (s Z) SetAircraft(v Aircraft) { C.Struct(s).Set16(0, 30); C.Struct(s).SetObject(0, C.Object(v)) }
func (s Z) Regression() Regression { return Regression(C.Struct(s).GetObject(0).ToStruct()) }
func (s Z) SetRegression(v Regression) {
	C.Struct(s).Set16(0, 31)
	C.Struct(s).SetObject(0, C.Object(v))
}
func (s Z) Planebase() PlaneBase     { return PlaneBase(C.Struct(s).GetObject(0).ToStruct()) }
func (s Z) SetPlanebase(v PlaneBase) { C.Struct(s).Set16(0, 32); C.Struct(s).SetObject(0, C.Object(v)) }
func (s Z) Airport() Airport         { return Airport(C.Struct(s).Get16(8)) }
func (s Z) SetAirport(v Airport)     { C.Struct(s).Set16(0, 33); C.Struct(s).Set16(8, uint16(v)) }
func (s Z) B737() B737               { return B737(C.Struct(s).GetObject(0).ToStruct()) }
func (s Z) SetB737(v B737)           { C.Struct(s).Set16(0, 34); C.Struct(s).SetObject(0, C.Object(v)) }
func (s Z) A320() A320               { return A320(C.Struct(s).GetObject(0).ToStruct()) }
func (s Z) SetA320(v A320)           { C.Struct(s).Set16(0, 35); C.Struct(s).SetObject(0, C.Object(v)) }
func (s Z) F16() F16                 { return F16(C.Struct(s).GetObject(0).ToStruct()) }
func (s Z) SetF16(v F16)             { C.Struct(s).Set16(0, 36); C.Struct(s).SetObject(0, C.Object(v)) }
func (s Z) Zdatevec() Zdate_List     { return Zdate_List(C.Struct(s).GetObject(0)) }
func (s Z) SetZdatevec(v Zdate_List) { C.Struct(s).Set16(0, 37); C.Struct(s).SetObject(0, C.Object(v)) }
func (s Z) Zdatavec() Zdata_List     { return Zdata_List(C.Struct(s).GetObject(0)) }
func (s Z) SetZdatavec(v Zdata_List) { C.Struct(s).Set16(0, 38); C.Struct(s).SetObject(0, C.Object(v)) }
func (s Z) Boolvec() C.BitList       { return C.BitList(C.Struct(s).GetObject(0)) }
func (s Z) SetBoolvec(v C.BitList)   { C.Struct(s).Set16(0, 39); C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s Z) MarshalJSON() (bs []byte, err error) { return }

type Z_List C.PointerList

func NewZList(s *C.Segment, sz int) Z_List { return Z_List(s.NewCompositeList(16, 1, sz)) }
func (s Z_List) Len() int                  { return C.PointerList(s).Len() }
func (s Z_List) At(i int) Z                { return Z(C.PointerList(s).At(i).ToStruct()) }
func (s Z_List) ToArray() []Z              { return *(*[]Z)(unsafe.Pointer(C.PointerList(s).ToArray())) }

type Counter C.Struct

func NewCounter(s *C.Segment) Counter      { return Counter(s.NewStruct(8, 2)) }
func NewRootCounter(s *C.Segment) Counter  { return Counter(s.NewRootStruct(8, 2)) }
func ReadRootCounter(s *C.Segment) Counter { return Counter(s.Root(0).ToStruct()) }
func (s Counter) Size() int64              { return int64(C.Struct(s).Get64(0)) }
func (s Counter) SetSize(v int64)          { C.Struct(s).Set64(0, uint64(v)) }
func (s Counter) Words() string            { return C.Struct(s).GetObject(0).ToText() }
func (s Counter) SetWords(v string)        { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s Counter) Wordlist() C.TextList     { return C.TextList(C.Struct(s).GetObject(1)) }
func (s Counter) SetWordlist(v C.TextList) { C.Struct(s).SetObject(1, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s Counter) MarshalJSON() (bs []byte, err error) { return }

type Counter_List C.PointerList

func NewCounterList(s *C.Segment, sz int) Counter_List {
	return Counter_List(s.NewCompositeList(8, 2, sz))
}
func (s Counter_List) Len() int         { return C.PointerList(s).Len() }
func (s Counter_List) At(i int) Counter { return Counter(C.PointerList(s).At(i).ToStruct()) }
func (s Counter_List) ToArray() []Counter {
	return *(*[]Counter)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type Bag C.Struct

func NewBag(s *C.Segment) Bag      { return Bag(s.NewStruct(0, 1)) }
func NewRootBag(s *C.Segment) Bag  { return Bag(s.NewRootStruct(0, 1)) }
func ReadRootBag(s *C.Segment) Bag { return Bag(s.Root(0).ToStruct()) }
func (s Bag) Counter() Counter     { return Counter(C.Struct(s).GetObject(0).ToStruct()) }
func (s Bag) SetCounter(v Counter) { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s Bag) MarshalJSON() (bs []byte, err error) { return }

type Bag_List C.PointerList

func NewBagList(s *C.Segment, sz int) Bag_List { return Bag_List(s.NewCompositeList(0, 1, sz)) }
func (s Bag_List) Len() int                    { return C.PointerList(s).Len() }
func (s Bag_List) At(i int) Bag                { return Bag(C.PointerList(s).At(i).ToStruct()) }
func (s Bag_List) ToArray() []Bag              { return *(*[]Bag)(unsafe.Pointer(C.PointerList(s).ToArray())) }

type Zserver C.Struct

func NewZserver(s *C.Segment) Zserver        { return Zserver(s.NewStruct(0, 1)) }
func NewRootZserver(s *C.Segment) Zserver    { return Zserver(s.NewRootStruct(0, 1)) }
func ReadRootZserver(s *C.Segment) Zserver   { return Zserver(s.Root(0).ToStruct()) }
func (s Zserver) Waitingjobs() Zjob_List     { return Zjob_List(C.Struct(s).GetObject(0)) }
func (s Zserver) SetWaitingjobs(v Zjob_List) { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s Zserver) MarshalJSON() (bs []byte, err error) { return }

type Zserver_List C.PointerList

func NewZserverList(s *C.Segment, sz int) Zserver_List {
	return Zserver_List(s.NewCompositeList(0, 1, sz))
}
func (s Zserver_List) Len() int         { return C.PointerList(s).Len() }
func (s Zserver_List) At(i int) Zserver { return Zserver(C.PointerList(s).At(i).ToStruct()) }
func (s Zserver_List) ToArray() []Zserver {
	return *(*[]Zserver)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type Zjob C.Struct

func NewZjob(s *C.Segment) Zjob      { return Zjob(s.NewStruct(0, 2)) }
func NewRootZjob(s *C.Segment) Zjob  { return Zjob(s.NewRootStruct(0, 2)) }
func ReadRootZjob(s *C.Segment) Zjob { return Zjob(s.Root(0).ToStruct()) }
func (s Zjob) Cmd() string           { return C.Struct(s).GetObject(0).ToText() }
func (s Zjob) SetCmd(v string)       { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s Zjob) Args() C.TextList      { return C.TextList(C.Struct(s).GetObject(1)) }
func (s Zjob) SetArgs(v C.TextList)  { C.Struct(s).SetObject(1, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s Zjob) MarshalJSON() (bs []byte, err error) { return }

type Zjob_List C.PointerList

func NewZjobList(s *C.Segment, sz int) Zjob_List { return Zjob_List(s.NewCompositeList(0, 2, sz)) }
func (s Zjob_List) Len() int                     { return C.PointerList(s).Len() }
func (s Zjob_List) At(i int) Zjob                { return Zjob(C.PointerList(s).At(i).ToStruct()) }
func (s Zjob_List) ToArray() []Zjob              { return *(*[]Zjob)(unsafe.Pointer(C.PointerList(s).ToArray())) }

type VerEmpty C.Struct

func NewVerEmpty(s *C.Segment) VerEmpty      { return VerEmpty(s.NewStruct(0, 0)) }
func NewRootVerEmpty(s *C.Segment) VerEmpty  { return VerEmpty(s.NewRootStruct(0, 0)) }
func ReadRootVerEmpty(s *C.Segment) VerEmpty { return VerEmpty(s.Root(0).ToStruct()) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s VerEmpty) MarshalJSON() (bs []byte, err error) { return }

type VerEmpty_List C.PointerList

func NewVerEmptyList(s *C.Segment, sz int) VerEmpty_List { return VerEmpty_List(s.NewVoidList(sz)) }
func (s VerEmpty_List) Len() int                         { return C.PointerList(s).Len() }
func (s VerEmpty_List) At(i int) VerEmpty                { return VerEmpty(C.PointerList(s).At(i).ToStruct()) }
func (s VerEmpty_List) ToArray() []VerEmpty {
	return *(*[]VerEmpty)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type VerOneData C.Struct

func NewVerOneData(s *C.Segment) VerOneData      { return VerOneData(s.NewStruct(8, 0)) }
func NewRootVerOneData(s *C.Segment) VerOneData  { return VerOneData(s.NewRootStruct(8, 0)) }
func ReadRootVerOneData(s *C.Segment) VerOneData { return VerOneData(s.Root(0).ToStruct()) }
func (s VerOneData) Val() int16                  { return int16(C.Struct(s).Get16(0)) }
func (s VerOneData) SetVal(v int16)              { C.Struct(s).Set16(0, uint16(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s VerOneData) MarshalJSON() (bs []byte, err error) { return }

type VerOneData_List C.PointerList

func NewVerOneDataList(s *C.Segment, sz int) VerOneData_List {
	return VerOneData_List(s.NewUInt16List(sz))
}
func (s VerOneData_List) Len() int            { return C.PointerList(s).Len() }
func (s VerOneData_List) At(i int) VerOneData { return VerOneData(C.PointerList(s).At(i).ToStruct()) }
func (s VerOneData_List) ToArray() []VerOneData {
	return *(*[]VerOneData)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type VerTwoData C.Struct

func NewVerTwoData(s *C.Segment) VerTwoData      { return VerTwoData(s.NewStruct(16, 0)) }
func NewRootVerTwoData(s *C.Segment) VerTwoData  { return VerTwoData(s.NewRootStruct(16, 0)) }
func ReadRootVerTwoData(s *C.Segment) VerTwoData { return VerTwoData(s.Root(0).ToStruct()) }
func (s VerTwoData) Val() int16                  { return int16(C.Struct(s).Get16(0)) }
func (s VerTwoData) SetVal(v int16)              { C.Struct(s).Set16(0, uint16(v)) }
func (s VerTwoData) Duo() int64                  { return int64(C.Struct(s).Get64(8)) }
func (s VerTwoData) SetDuo(v int64)              { C.Struct(s).Set64(8, uint64(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s VerTwoData) MarshalJSON() (bs []byte, err error) { return }

type VerTwoData_List C.PointerList

func NewVerTwoDataList(s *C.Segment, sz int) VerTwoData_List {
	return VerTwoData_List(s.NewCompositeList(16, 0, sz))
}
func (s VerTwoData_List) Len() int            { return C.PointerList(s).Len() }
func (s VerTwoData_List) At(i int) VerTwoData { return VerTwoData(C.PointerList(s).At(i).ToStruct()) }
func (s VerTwoData_List) ToArray() []VerTwoData {
	return *(*[]VerTwoData)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type VerOnePtr C.Struct

func NewVerOnePtr(s *C.Segment) VerOnePtr      { return VerOnePtr(s.NewStruct(0, 1)) }
func NewRootVerOnePtr(s *C.Segment) VerOnePtr  { return VerOnePtr(s.NewRootStruct(0, 1)) }
func ReadRootVerOnePtr(s *C.Segment) VerOnePtr { return VerOnePtr(s.Root(0).ToStruct()) }
func (s VerOnePtr) Ptr() VerOneData            { return VerOneData(C.Struct(s).GetObject(0).ToStruct()) }
func (s VerOnePtr) SetPtr(v VerOneData)        { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s VerOnePtr) MarshalJSON() (bs []byte, err error) { return }

type VerOnePtr_List C.PointerList

func NewVerOnePtrList(s *C.Segment, sz int) VerOnePtr_List {
	return VerOnePtr_List(s.NewCompositeList(0, 1, sz))
}
func (s VerOnePtr_List) Len() int           { return C.PointerList(s).Len() }
func (s VerOnePtr_List) At(i int) VerOnePtr { return VerOnePtr(C.PointerList(s).At(i).ToStruct()) }
func (s VerOnePtr_List) ToArray() []VerOnePtr {
	return *(*[]VerOnePtr)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type VerTwoPtr C.Struct

func NewVerTwoPtr(s *C.Segment) VerTwoPtr      { return VerTwoPtr(s.NewStruct(0, 2)) }
func NewRootVerTwoPtr(s *C.Segment) VerTwoPtr  { return VerTwoPtr(s.NewRootStruct(0, 2)) }
func ReadRootVerTwoPtr(s *C.Segment) VerTwoPtr { return VerTwoPtr(s.Root(0).ToStruct()) }
func (s VerTwoPtr) Ptr1() VerOneData           { return VerOneData(C.Struct(s).GetObject(0).ToStruct()) }
func (s VerTwoPtr) SetPtr1(v VerOneData)       { C.Struct(s).SetObject(0, C.Object(v)) }
func (s VerTwoPtr) Ptr2() VerOneData           { return VerOneData(C.Struct(s).GetObject(1).ToStruct()) }
func (s VerTwoPtr) SetPtr2(v VerOneData)       { C.Struct(s).SetObject(1, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s VerTwoPtr) MarshalJSON() (bs []byte, err error) { return }

type VerTwoPtr_List C.PointerList

func NewVerTwoPtrList(s *C.Segment, sz int) VerTwoPtr_List {
	return VerTwoPtr_List(s.NewCompositeList(0, 2, sz))
}
func (s VerTwoPtr_List) Len() int           { return C.PointerList(s).Len() }
func (s VerTwoPtr_List) At(i int) VerTwoPtr { return VerTwoPtr(C.PointerList(s).At(i).ToStruct()) }
func (s VerTwoPtr_List) ToArray() []VerTwoPtr {
	return *(*[]VerTwoPtr)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type VerTwoDataTwoPtr C.Struct

func NewVerTwoDataTwoPtr(s *C.Segment) VerTwoDataTwoPtr { return VerTwoDataTwoPtr(s.NewStruct(16, 2)) }
func NewRootVerTwoDataTwoPtr(s *C.Segment) VerTwoDataTwoPtr {
	return VerTwoDataTwoPtr(s.NewRootStruct(16, 2))
}
func ReadRootVerTwoDataTwoPtr(s *C.Segment) VerTwoDataTwoPtr {
	return VerTwoDataTwoPtr(s.Root(0).ToStruct())
}
func (s VerTwoDataTwoPtr) Val() int16           { return int16(C.Struct(s).Get16(0)) }
func (s VerTwoDataTwoPtr) SetVal(v int16)       { C.Struct(s).Set16(0, uint16(v)) }
func (s VerTwoDataTwoPtr) Duo() int64           { return int64(C.Struct(s).Get64(8)) }
func (s VerTwoDataTwoPtr) SetDuo(v int64)       { C.Struct(s).Set64(8, uint64(v)) }
func (s VerTwoDataTwoPtr) Ptr1() VerOneData     { return VerOneData(C.Struct(s).GetObject(0).ToStruct()) }
func (s VerTwoDataTwoPtr) SetPtr1(v VerOneData) { C.Struct(s).SetObject(0, C.Object(v)) }
func (s VerTwoDataTwoPtr) Ptr2() VerOneData     { return VerOneData(C.Struct(s).GetObject(1).ToStruct()) }
func (s VerTwoDataTwoPtr) SetPtr2(v VerOneData) { C.Struct(s).SetObject(1, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s VerTwoDataTwoPtr) MarshalJSON() (bs []byte, err error) { return }

type VerTwoDataTwoPtr_List C.PointerList

func NewVerTwoDataTwoPtrList(s *C.Segment, sz int) VerTwoDataTwoPtr_List {
	return VerTwoDataTwoPtr_List(s.NewCompositeList(16, 2, sz))
}
func (s VerTwoDataTwoPtr_List) Len() int { return C.PointerList(s).Len() }
func (s VerTwoDataTwoPtr_List) At(i int) VerTwoDataTwoPtr {
	return VerTwoDataTwoPtr(C.PointerList(s).At(i).ToStruct())
}
func (s VerTwoDataTwoPtr_List) ToArray() []VerTwoDataTwoPtr {
	return *(*[]VerTwoDataTwoPtr)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type HoldsVerEmptyList C.Struct

func NewHoldsVerEmptyList(s *C.Segment) HoldsVerEmptyList { return HoldsVerEmptyList(s.NewStruct(0, 1)) }
func NewRootHoldsVerEmptyList(s *C.Segment) HoldsVerEmptyList {
	return HoldsVerEmptyList(s.NewRootStruct(0, 1))
}
func ReadRootHoldsVerEmptyList(s *C.Segment) HoldsVerEmptyList {
	return HoldsVerEmptyList(s.Root(0).ToStruct())
}
func (s HoldsVerEmptyList) Mylist() VerEmpty_List     { return VerEmpty_List(C.Struct(s).GetObject(0)) }
func (s HoldsVerEmptyList) SetMylist(v VerEmpty_List) { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s HoldsVerEmptyList) MarshalJSON() (bs []byte, err error) { return }

type HoldsVerEmptyList_List C.PointerList

func NewHoldsVerEmptyListList(s *C.Segment, sz int) HoldsVerEmptyList_List {
	return HoldsVerEmptyList_List(s.NewCompositeList(0, 1, sz))
}
func (s HoldsVerEmptyList_List) Len() int { return C.PointerList(s).Len() }
func (s HoldsVerEmptyList_List) At(i int) HoldsVerEmptyList {
	return HoldsVerEmptyList(C.PointerList(s).At(i).ToStruct())
}
func (s HoldsVerEmptyList_List) ToArray() []HoldsVerEmptyList {
	return *(*[]HoldsVerEmptyList)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type HoldsVerOneDataList C.Struct

func NewHoldsVerOneDataList(s *C.Segment) HoldsVerOneDataList {
	return HoldsVerOneDataList(s.NewStruct(0, 1))
}
func NewRootHoldsVerOneDataList(s *C.Segment) HoldsVerOneDataList {
	return HoldsVerOneDataList(s.NewRootStruct(0, 1))
}
func ReadRootHoldsVerOneDataList(s *C.Segment) HoldsVerOneDataList {
	return HoldsVerOneDataList(s.Root(0).ToStruct())
}
func (s HoldsVerOneDataList) Mylist() VerOneData_List {
	return VerOneData_List(C.Struct(s).GetObject(0))
}
func (s HoldsVerOneDataList) SetMylist(v VerOneData_List) { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s HoldsVerOneDataList) MarshalJSON() (bs []byte, err error) { return }

type HoldsVerOneDataList_List C.PointerList

func NewHoldsVerOneDataListList(s *C.Segment, sz int) HoldsVerOneDataList_List {
	return HoldsVerOneDataList_List(s.NewCompositeList(0, 1, sz))
}
func (s HoldsVerOneDataList_List) Len() int { return C.PointerList(s).Len() }
func (s HoldsVerOneDataList_List) At(i int) HoldsVerOneDataList {
	return HoldsVerOneDataList(C.PointerList(s).At(i).ToStruct())
}
func (s HoldsVerOneDataList_List) ToArray() []HoldsVerOneDataList {
	return *(*[]HoldsVerOneDataList)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type HoldsVerTwoDataList C.Struct

func NewHoldsVerTwoDataList(s *C.Segment) HoldsVerTwoDataList {
	return HoldsVerTwoDataList(s.NewStruct(0, 1))
}
func NewRootHoldsVerTwoDataList(s *C.Segment) HoldsVerTwoDataList {
	return HoldsVerTwoDataList(s.NewRootStruct(0, 1))
}
func ReadRootHoldsVerTwoDataList(s *C.Segment) HoldsVerTwoDataList {
	return HoldsVerTwoDataList(s.Root(0).ToStruct())
}
func (s HoldsVerTwoDataList) Mylist() VerTwoData_List {
	return VerTwoData_List(C.Struct(s).GetObject(0))
}
func (s HoldsVerTwoDataList) SetMylist(v VerTwoData_List) { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s HoldsVerTwoDataList) MarshalJSON() (bs []byte, err error) { return }

type HoldsVerTwoDataList_List C.PointerList

func NewHoldsVerTwoDataListList(s *C.Segment, sz int) HoldsVerTwoDataList_List {
	return HoldsVerTwoDataList_List(s.NewCompositeList(0, 1, sz))
}
func (s HoldsVerTwoDataList_List) Len() int { return C.PointerList(s).Len() }
func (s HoldsVerTwoDataList_List) At(i int) HoldsVerTwoDataList {
	return HoldsVerTwoDataList(C.PointerList(s).At(i).ToStruct())
}
func (s HoldsVerTwoDataList_List) ToArray() []HoldsVerTwoDataList {
	return *(*[]HoldsVerTwoDataList)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type HoldsVerOnePtrList C.Struct

func NewHoldsVerOnePtrList(s *C.Segment) HoldsVerOnePtrList {
	return HoldsVerOnePtrList(s.NewStruct(0, 1))
}
func NewRootHoldsVerOnePtrList(s *C.Segment) HoldsVerOnePtrList {
	return HoldsVerOnePtrList(s.NewRootStruct(0, 1))
}
func ReadRootHoldsVerOnePtrList(s *C.Segment) HoldsVerOnePtrList {
	return HoldsVerOnePtrList(s.Root(0).ToStruct())
}
func (s HoldsVerOnePtrList) Mylist() VerOnePtr_List     { return VerOnePtr_List(C.Struct(s).GetObject(0)) }
func (s HoldsVerOnePtrList) SetMylist(v VerOnePtr_List) { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s HoldsVerOnePtrList) MarshalJSON() (bs []byte, err error) { return }

type HoldsVerOnePtrList_List C.PointerList

func NewHoldsVerOnePtrListList(s *C.Segment, sz int) HoldsVerOnePtrList_List {
	return HoldsVerOnePtrList_List(s.NewCompositeList(0, 1, sz))
}
func (s HoldsVerOnePtrList_List) Len() int { return C.PointerList(s).Len() }
func (s HoldsVerOnePtrList_List) At(i int) HoldsVerOnePtrList {
	return HoldsVerOnePtrList(C.PointerList(s).At(i).ToStruct())
}
func (s HoldsVerOnePtrList_List) ToArray() []HoldsVerOnePtrList {
	return *(*[]HoldsVerOnePtrList)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type HoldsVerTwoPtrList C.Struct

func NewHoldsVerTwoPtrList(s *C.Segment) HoldsVerTwoPtrList {
	return HoldsVerTwoPtrList(s.NewStruct(0, 1))
}
func NewRootHoldsVerTwoPtrList(s *C.Segment) HoldsVerTwoPtrList {
	return HoldsVerTwoPtrList(s.NewRootStruct(0, 1))
}
func ReadRootHoldsVerTwoPtrList(s *C.Segment) HoldsVerTwoPtrList {
	return HoldsVerTwoPtrList(s.Root(0).ToStruct())
}
func (s HoldsVerTwoPtrList) Mylist() VerTwoPtr_List     { return VerTwoPtr_List(C.Struct(s).GetObject(0)) }
func (s HoldsVerTwoPtrList) SetMylist(v VerTwoPtr_List) { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s HoldsVerTwoPtrList) MarshalJSON() (bs []byte, err error) { return }

type HoldsVerTwoPtrList_List C.PointerList

func NewHoldsVerTwoPtrListList(s *C.Segment, sz int) HoldsVerTwoPtrList_List {
	return HoldsVerTwoPtrList_List(s.NewCompositeList(0, 1, sz))
}
func (s HoldsVerTwoPtrList_List) Len() int { return C.PointerList(s).Len() }
func (s HoldsVerTwoPtrList_List) At(i int) HoldsVerTwoPtrList {
	return HoldsVerTwoPtrList(C.PointerList(s).At(i).ToStruct())
}
func (s HoldsVerTwoPtrList_List) ToArray() []HoldsVerTwoPtrList {
	return *(*[]HoldsVerTwoPtrList)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type HoldsVerTwoTwoList C.Struct

func NewHoldsVerTwoTwoList(s *C.Segment) HoldsVerTwoTwoList {
	return HoldsVerTwoTwoList(s.NewStruct(0, 1))
}
func NewRootHoldsVerTwoTwoList(s *C.Segment) HoldsVerTwoTwoList {
	return HoldsVerTwoTwoList(s.NewRootStruct(0, 1))
}
func ReadRootHoldsVerTwoTwoList(s *C.Segment) HoldsVerTwoTwoList {
	return HoldsVerTwoTwoList(s.Root(0).ToStruct())
}
func (s HoldsVerTwoTwoList) Mylist() VerTwoDataTwoPtr_List {
	return VerTwoDataTwoPtr_List(C.Struct(s).GetObject(0))
}
func (s HoldsVerTwoTwoList) SetMylist(v VerTwoDataTwoPtr_List) { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s HoldsVerTwoTwoList) MarshalJSON() (bs []byte, err error) { return }

type HoldsVerTwoTwoList_List C.PointerList

func NewHoldsVerTwoTwoListList(s *C.Segment, sz int) HoldsVerTwoTwoList_List {
	return HoldsVerTwoTwoList_List(s.NewCompositeList(0, 1, sz))
}
func (s HoldsVerTwoTwoList_List) Len() int { return C.PointerList(s).Len() }
func (s HoldsVerTwoTwoList_List) At(i int) HoldsVerTwoTwoList {
	return HoldsVerTwoTwoList(C.PointerList(s).At(i).ToStruct())
}
func (s HoldsVerTwoTwoList_List) ToArray() []HoldsVerTwoTwoList {
	return *(*[]HoldsVerTwoTwoList)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type HoldsVerTwoTwoPlus C.Struct

func NewHoldsVerTwoTwoPlus(s *C.Segment) HoldsVerTwoTwoPlus {
	return HoldsVerTwoTwoPlus(s.NewStruct(0, 1))
}
func NewRootHoldsVerTwoTwoPlus(s *C.Segment) HoldsVerTwoTwoPlus {
	return HoldsVerTwoTwoPlus(s.NewRootStruct(0, 1))
}
func ReadRootHoldsVerTwoTwoPlus(s *C.Segment) HoldsVerTwoTwoPlus {
	return HoldsVerTwoTwoPlus(s.Root(0).ToStruct())
}
func (s HoldsVerTwoTwoPlus) Mylist() VerTwoTwoPlus_List {
	return VerTwoTwoPlus_List(C.Struct(s).GetObject(0))
}
func (s HoldsVerTwoTwoPlus) SetMylist(v VerTwoTwoPlus_List) { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s HoldsVerTwoTwoPlus) MarshalJSON() (bs []byte, err error) { return }

type HoldsVerTwoTwoPlus_List C.PointerList

func NewHoldsVerTwoTwoPlusList(s *C.Segment, sz int) HoldsVerTwoTwoPlus_List {
	return HoldsVerTwoTwoPlus_List(s.NewCompositeList(0, 1, sz))
}
func (s HoldsVerTwoTwoPlus_List) Len() int { return C.PointerList(s).Len() }
func (s HoldsVerTwoTwoPlus_List) At(i int) HoldsVerTwoTwoPlus {
	return HoldsVerTwoTwoPlus(C.PointerList(s).At(i).ToStruct())
}
func (s HoldsVerTwoTwoPlus_List) ToArray() []HoldsVerTwoTwoPlus {
	return *(*[]HoldsVerTwoTwoPlus)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type VerTwoTwoPlus C.Struct

func NewVerTwoTwoPlus(s *C.Segment) VerTwoTwoPlus      { return VerTwoTwoPlus(s.NewStruct(24, 3)) }
func NewRootVerTwoTwoPlus(s *C.Segment) VerTwoTwoPlus  { return VerTwoTwoPlus(s.NewRootStruct(24, 3)) }
func ReadRootVerTwoTwoPlus(s *C.Segment) VerTwoTwoPlus { return VerTwoTwoPlus(s.Root(0).ToStruct()) }
func (s VerTwoTwoPlus) Val() int16                     { return int16(C.Struct(s).Get16(0)) }
func (s VerTwoTwoPlus) SetVal(v int16)                 { C.Struct(s).Set16(0, uint16(v)) }
func (s VerTwoTwoPlus) Duo() int64                     { return int64(C.Struct(s).Get64(8)) }
func (s VerTwoTwoPlus) SetDuo(v int64)                 { C.Struct(s).Set64(8, uint64(v)) }
func (s VerTwoTwoPlus) Ptr1() VerTwoDataTwoPtr {
	return VerTwoDataTwoPtr(C.Struct(s).GetObject(0).ToStruct())
}
func (s VerTwoTwoPlus) SetPtr1(v VerTwoDataTwoPtr) { C.Struct(s).SetObject(0, C.Object(v)) }
func (s VerTwoTwoPlus) Ptr2() VerTwoDataTwoPtr {
	return VerTwoDataTwoPtr(C.Struct(s).GetObject(1).ToStruct())
}
func (s VerTwoTwoPlus) SetPtr2(v VerTwoDataTwoPtr) { C.Struct(s).SetObject(1, C.Object(v)) }
func (s VerTwoTwoPlus) Tre() int64                 { return int64(C.Struct(s).Get64(16)) }
func (s VerTwoTwoPlus) SetTre(v int64)             { C.Struct(s).Set64(16, uint64(v)) }
func (s VerTwoTwoPlus) Lst3() C.Int64List          { return C.Int64List(C.Struct(s).GetObject(2)) }
func (s VerTwoTwoPlus) SetLst3(v C.Int64List)      { C.Struct(s).SetObject(2, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s VerTwoTwoPlus) MarshalJSON() (bs []byte, err error) { return }

type VerTwoTwoPlus_List C.PointerList

func NewVerTwoTwoPlusList(s *C.Segment, sz int) VerTwoTwoPlus_List {
	return VerTwoTwoPlus_List(s.NewCompositeList(24, 3, sz))
}
func (s VerTwoTwoPlus_List) Len() int { return C.PointerList(s).Len() }
func (s VerTwoTwoPlus_List) At(i int) VerTwoTwoPlus {
	return VerTwoTwoPlus(C.PointerList(s).At(i).ToStruct())
}
func (s VerTwoTwoPlus_List) ToArray() []VerTwoTwoPlus {
	return *(*[]VerTwoTwoPlus)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type HoldsText C.Struct

func NewHoldsText(s *C.Segment) HoldsText      { return HoldsText(s.NewStruct(0, 3)) }
func NewRootHoldsText(s *C.Segment) HoldsText  { return HoldsText(s.NewRootStruct(0, 3)) }
func ReadRootHoldsText(s *C.Segment) HoldsText { return HoldsText(s.Root(0).ToStruct()) }
func (s HoldsText) Txt() string                { return C.Struct(s).GetObject(0).ToText() }
func (s HoldsText) SetTxt(v string)            { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s HoldsText) Lst() C.TextList            { return C.TextList(C.Struct(s).GetObject(1)) }
func (s HoldsText) SetLst(v C.TextList)        { C.Struct(s).SetObject(1, C.Object(v)) }
func (s HoldsText) Lstlst() C.PointerList      { return C.PointerList(C.Struct(s).GetObject(2)) }
func (s HoldsText) SetLstlst(v C.PointerList)  { C.Struct(s).SetObject(2, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s HoldsText) MarshalJSON() (bs []byte, err error) { return }

type HoldsText_List C.PointerList

func NewHoldsTextList(s *C.Segment, sz int) HoldsText_List {
	return HoldsText_List(s.NewCompositeList(0, 3, sz))
}
func (s HoldsText_List) Len() int           { return C.PointerList(s).Len() }
func (s HoldsText_List) At(i int) HoldsText { return HoldsText(C.PointerList(s).At(i).ToStruct()) }
func (s HoldsText_List) ToArray() []HoldsText {
	return *(*[]HoldsText)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type WrapEmpty C.Struct

func NewWrapEmpty(s *C.Segment) WrapEmpty      { return WrapEmpty(s.NewStruct(0, 1)) }
func NewRootWrapEmpty(s *C.Segment) WrapEmpty  { return WrapEmpty(s.NewRootStruct(0, 1)) }
func ReadRootWrapEmpty(s *C.Segment) WrapEmpty { return WrapEmpty(s.Root(0).ToStruct()) }
func (s WrapEmpty) MightNotBeReallyEmpty() VerEmpty {
	return VerEmpty(C.Struct(s).GetObject(0).ToStruct())
}
func (s WrapEmpty) SetMightNotBeReallyEmpty(v VerEmpty) { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s WrapEmpty) MarshalJSON() (bs []byte, err error) { return }

type WrapEmpty_List C.PointerList

func NewWrapEmptyList(s *C.Segment, sz int) WrapEmpty_List {
	return WrapEmpty_List(s.NewCompositeList(0, 1, sz))
}
func (s WrapEmpty_List) Len() int           { return C.PointerList(s).Len() }
func (s WrapEmpty_List) At(i int) WrapEmpty { return WrapEmpty(C.PointerList(s).At(i).ToStruct()) }
func (s WrapEmpty_List) ToArray() []WrapEmpty {
	return *(*[]WrapEmpty)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type Wrap2x2 C.Struct

func NewWrap2x2(s *C.Segment) Wrap2x2      { return Wrap2x2(s.NewStruct(0, 1)) }
func NewRootWrap2x2(s *C.Segment) Wrap2x2  { return Wrap2x2(s.NewRootStruct(0, 1)) }
func ReadRootWrap2x2(s *C.Segment) Wrap2x2 { return Wrap2x2(s.Root(0).ToStruct()) }
func (s Wrap2x2) MightNotBeReallyEmpty() VerTwoDataTwoPtr {
	return VerTwoDataTwoPtr(C.Struct(s).GetObject(0).ToStruct())
}
func (s Wrap2x2) SetMightNotBeReallyEmpty(v VerTwoDataTwoPtr) { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s Wrap2x2) MarshalJSON() (bs []byte, err error) { return }

type Wrap2x2_List C.PointerList

func NewWrap2x2List(s *C.Segment, sz int) Wrap2x2_List {
	return Wrap2x2_List(s.NewCompositeList(0, 1, sz))
}
func (s Wrap2x2_List) Len() int         { return C.PointerList(s).Len() }
func (s Wrap2x2_List) At(i int) Wrap2x2 { return Wrap2x2(C.PointerList(s).At(i).ToStruct()) }
func (s Wrap2x2_List) ToArray() []Wrap2x2 {
	return *(*[]Wrap2x2)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type Wrap2x2plus C.Struct

func NewWrap2x2plus(s *C.Segment) Wrap2x2plus      { return Wrap2x2plus(s.NewStruct(0, 1)) }
func NewRootWrap2x2plus(s *C.Segment) Wrap2x2plus  { return Wrap2x2plus(s.NewRootStruct(0, 1)) }
func ReadRootWrap2x2plus(s *C.Segment) Wrap2x2plus { return Wrap2x2plus(s.Root(0).ToStruct()) }
func (s Wrap2x2plus) MightNotBeReallyEmpty() VerTwoTwoPlus {
	return VerTwoTwoPlus(C.Struct(s).GetObject(0).ToStruct())
}
func (s Wrap2x2plus) SetMightNotBeReallyEmpty(v VerTwoTwoPlus) { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s Wrap2x2plus) MarshalJSON() (bs []byte, err error) { return }

type Wrap2x2plus_List C.PointerList

func NewWrap2x2plusList(s *C.Segment, sz int) Wrap2x2plus_List {
	return Wrap2x2plus_List(s.NewCompositeList(0, 1, sz))
}
func (s Wrap2x2plus_List) Len() int             { return C.PointerList(s).Len() }
func (s Wrap2x2plus_List) At(i int) Wrap2x2plus { return Wrap2x2plus(C.PointerList(s).At(i).ToStruct()) }
func (s Wrap2x2plus_List) ToArray() []Wrap2x2plus {
	return *(*[]Wrap2x2plus)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type Endpoint C.Struct

func NewEndpoint(s *C.Segment) Endpoint      { return Endpoint(s.NewStruct(8, 2)) }
func NewRootEndpoint(s *C.Segment) Endpoint  { return Endpoint(s.NewRootStruct(8, 2)) }
func ReadRootEndpoint(s *C.Segment) Endpoint { return Endpoint(s.Root(0).ToStruct()) }
func (s Endpoint) Ip() net.IP                { return net.IP(C.Struct(s).GetObject(0).ToData()) }
func (s Endpoint) SetIp(v net.IP)            { C.Struct(s).SetObject(0, s.Segment.NewData([]byte(v))) }
func (s Endpoint) Port() int16               { return int16(C.Struct(s).Get16(0)) }
func (s Endpoint) SetPort(v int16)           { C.Struct(s).Set16(0, uint16(v)) }
func (s Endpoint) Hostname() string          { return C.Struct(s).GetObject(1).ToText() }
func (s Endpoint) SetHostname(v string)      { C.Struct(s).SetObject(1, s.Segment.NewText(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s Endpoint) MarshalJSON() (bs []byte, err error) { return }

type Endpoint_List C.PointerList

func NewEndpointList(s *C.Segment, sz int) Endpoint_List {
	return Endpoint_List(s.NewCompositeList(8, 2, sz))
}
func (s Endpoint_List) Len() int          { return C.PointerList(s).Len() }
func (s Endpoint_List) At(i int) Endpoint { return Endpoint(C.PointerList(s).At(i).ToStruct()) }
func (s Endpoint_List) ToArray() []Endpoint {
	return *(*[]Endpoint)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type VoidUnion C.Struct
type VoidUnion_Which uint16

const (
	VOIDUNION_A VoidUnion_Which = 0
	VOIDUNION_B VoidUnion_Which = 1
)

func NewVoidUnion(s *C.Segment) VoidUnion      { return VoidUnion(s.NewStruct(8, 0)) }
func NewRootVoidUnion(s *C.Segment) VoidUnion  { return VoidUnion(s.NewRootStruct(8, 0)) }
func ReadRootVoidUnion(s *C.Segment) VoidUnion { return VoidUnion(s.Root(0).ToStruct()) }
func (s VoidUnion) Which() VoidUnion_Which     { return VoidUnion_Which(C.Struct(s).Get16(0)) }
func (s VoidUnion) SetA()                      { C.Struct(s).Set16(0, 0) }
func (s VoidUnion) SetB()                      { C.Struct(s).Set16(0, 1) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s VoidUnion) MarshalJSON() (bs []byte, err error) { return }

type VoidUnion_List C.PointerList

func NewVoidUnionList(s *C.Segment, sz int) VoidUnion_List { return VoidUnion_List(s.NewUInt16List(sz)) }
func (s VoidUnion_List) Len() int                          { return C.PointerList(s).Len() }
func (s VoidUnion_List) At(i int) VoidUnion                { return VoidUnion(C.PointerList(s).At(i).ToStruct()) }
func (s VoidUnion_List) ToArray() []VoidUnion {
	return *(*[]VoidUnion)(unsafe.Pointer(C.PointerList(s).ToArray()))
}
