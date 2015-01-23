package aircraftlib

// AUTO GENERATED - DO NOT EDIT

import (
	C "github.com/glycerine/go-capnproto"
	context "golang.org/x/net/context"
	math "math"
	net "net"
)

type Zdate C.Struct

func NewZdate(s *C.Segment) Zdate {
	return Zdate(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 0}))
}
func NewRootZdate(s *C.Segment) Zdate {
	return Zdate(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 0}))
}
func AutoNewZdate(s *C.Segment) Zdate {
	return Zdate(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 0}))
}
func ReadRootZdate(s *C.Segment) Zdate { return Zdate(s.Root(0).ToStruct()) }
func (s Zdate) Year() int16            { return int16(C.Struct(s).Get16(0)) }
func (s Zdate) SetYear(v int16)        { C.Struct(s).Set16(0, uint16(v)) }
func (s Zdate) Month() uint8           { return C.Struct(s).Get8(2) }
func (s Zdate) SetMonth(v uint8)       { C.Struct(s).Set8(2, v) }
func (s Zdate) Day() uint8             { return C.Struct(s).Get8(3) }
func (s Zdate) SetDay(v uint8)         { C.Struct(s).Set8(3, v) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Zdate) MarshalJSON() (bs []byte, err error) { return }

type Zdate_List C.PointerList

func NewZdate_List(s *C.Segment, sz int) Zdate_List {
	return Zdate_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 0}, sz))
}
func (s Zdate_List) Len() int              { return C.PointerList(s).Len() }
func (s Zdate_List) At(i int) Zdate        { return Zdate(C.PointerList(s).At(i).ToStruct()) }
func (s Zdate_List) Set(i int, item Zdate) { C.PointerList(s).Set(i, C.Object(item)) }

type Zdate_Promise C.Pipeline

func (p *Zdate_Promise) Get() (Zdate, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Zdate(s), err
}

type Zdata C.Struct

func NewZdata(s *C.Segment) Zdata {
	return Zdata(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func NewRootZdata(s *C.Segment) Zdata {
	return Zdata(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewZdata(s *C.Segment) Zdata {
	return Zdata(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootZdata(s *C.Segment) Zdata { return Zdata(s.Root(0).ToStruct()) }
func (s Zdata) Data() []byte           { return C.Struct(s).GetObject(0).ToData() }
func (s Zdata) SetData(v []byte)       { C.Struct(s).SetObject(0, s.Segment.NewData(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Zdata) MarshalJSON() (bs []byte, err error) { return }

type Zdata_List C.PointerList

func NewZdata_List(s *C.Segment, sz int) Zdata_List {
	return Zdata_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s Zdata_List) Len() int              { return C.PointerList(s).Len() }
func (s Zdata_List) At(i int) Zdata        { return Zdata(C.PointerList(s).At(i).ToStruct()) }
func (s Zdata_List) Set(i int, item Zdata) { C.PointerList(s).Set(i, C.Object(item)) }

type Zdata_Promise C.Pipeline

func (p *Zdata_Promise) Get() (Zdata, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Zdata(s), err
}

type Airport uint16

const (
	Airport_none Airport = 0
	Airport_jfk  Airport = 1
	Airport_lax  Airport = 2
	Airport_sfo  Airport = 3
	Airport_luv  Airport = 4
	Airport_dfw  Airport = 5
	Airport_test Airport = 6
)

func (c Airport) String() string {
	switch c {
	case Airport_none:
		return "none"
	case Airport_jfk:
		return "jfk"
	case Airport_lax:
		return "lax"
	case Airport_sfo:
		return "sfo"
	case Airport_luv:
		return "luv"
	case Airport_dfw:
		return "dfw"
	case Airport_test:
		return "test"
	default:
		return ""
	}
}

func AirportFromString(c string) Airport {
	switch c {
	case "none":
		return Airport_none
	case "jfk":
		return Airport_jfk
	case "lax":
		return Airport_lax
	case "sfo":
		return Airport_sfo
	case "luv":
		return Airport_luv
	case "dfw":
		return Airport_dfw
	case "test":
		return Airport_test
	default:
		return 0
	}
}

type Airport_List C.PointerList

func NewAirport_List(s *C.Segment, sz int) Airport_List { return Airport_List(s.NewUInt16List(sz)) }
func (s Airport_List) Len() int                         { return C.UInt16List(s).Len() }
func (s Airport_List) At(i int) Airport                 { return Airport(C.UInt16List(s).At(i)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Airport) MarshalJSON() (bs []byte, err error) { return }

type PlaneBase C.Struct

func NewPlaneBase(s *C.Segment) PlaneBase {
	return PlaneBase(s.NewStruct(C.ObjectSize{DataSize: 32, PointerCount: 2}))
}
func NewRootPlaneBase(s *C.Segment) PlaneBase {
	return PlaneBase(s.NewRootStruct(C.ObjectSize{DataSize: 32, PointerCount: 2}))
}
func AutoNewPlaneBase(s *C.Segment) PlaneBase {
	return PlaneBase(s.NewStructAR(C.ObjectSize{DataSize: 32, PointerCount: 2}))
}
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

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s PlaneBase) MarshalJSON() (bs []byte, err error) { return }

type PlaneBase_List C.PointerList

func NewPlaneBase_List(s *C.Segment, sz int) PlaneBase_List {
	return PlaneBase_List(s.NewCompositeList(C.ObjectSize{DataSize: 32, PointerCount: 2}, sz))
}
func (s PlaneBase_List) Len() int                  { return C.PointerList(s).Len() }
func (s PlaneBase_List) At(i int) PlaneBase        { return PlaneBase(C.PointerList(s).At(i).ToStruct()) }
func (s PlaneBase_List) Set(i int, item PlaneBase) { C.PointerList(s).Set(i, C.Object(item)) }

type PlaneBase_Promise C.Pipeline

func (p *PlaneBase_Promise) Get() (PlaneBase, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return PlaneBase(s), err
}

type B737 C.Struct

func NewB737(s *C.Segment) B737 { return B737(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1})) }
func NewRootB737(s *C.Segment) B737 {
	return B737(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewB737(s *C.Segment) B737 {
	return B737(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootB737(s *C.Segment) B737 { return B737(s.Root(0).ToStruct()) }
func (s B737) Base() PlaneBase       { return PlaneBase(C.Struct(s).GetObject(0).ToStruct()) }
func (s B737) SetBase(v PlaneBase)   { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s B737) MarshalJSON() (bs []byte, err error) { return }

type B737_List C.PointerList

func NewB737_List(s *C.Segment, sz int) B737_List {
	return B737_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s B737_List) Len() int             { return C.PointerList(s).Len() }
func (s B737_List) At(i int) B737        { return B737(C.PointerList(s).At(i).ToStruct()) }
func (s B737_List) Set(i int, item B737) { C.PointerList(s).Set(i, C.Object(item)) }

type B737_Promise C.Pipeline

func (p *B737_Promise) Get() (B737, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return B737(s), err
}

func (p *B737_Promise) Base() *PlaneBase_Promise {
	return (*PlaneBase_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

type A320 C.Struct

func NewA320(s *C.Segment) A320 { return A320(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1})) }
func NewRootA320(s *C.Segment) A320 {
	return A320(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewA320(s *C.Segment) A320 {
	return A320(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootA320(s *C.Segment) A320 { return A320(s.Root(0).ToStruct()) }
func (s A320) Base() PlaneBase       { return PlaneBase(C.Struct(s).GetObject(0).ToStruct()) }
func (s A320) SetBase(v PlaneBase)   { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s A320) MarshalJSON() (bs []byte, err error) { return }

type A320_List C.PointerList

func NewA320_List(s *C.Segment, sz int) A320_List {
	return A320_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s A320_List) Len() int             { return C.PointerList(s).Len() }
func (s A320_List) At(i int) A320        { return A320(C.PointerList(s).At(i).ToStruct()) }
func (s A320_List) Set(i int, item A320) { C.PointerList(s).Set(i, C.Object(item)) }

type A320_Promise C.Pipeline

func (p *A320_Promise) Get() (A320, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return A320(s), err
}

func (p *A320_Promise) Base() *PlaneBase_Promise {
	return (*PlaneBase_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

type F16 C.Struct

func NewF16(s *C.Segment) F16 { return F16(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1})) }
func NewRootF16(s *C.Segment) F16 {
	return F16(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewF16(s *C.Segment) F16 {
	return F16(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootF16(s *C.Segment) F16 { return F16(s.Root(0).ToStruct()) }
func (s F16) Base() PlaneBase      { return PlaneBase(C.Struct(s).GetObject(0).ToStruct()) }
func (s F16) SetBase(v PlaneBase)  { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s F16) MarshalJSON() (bs []byte, err error) { return }

type F16_List C.PointerList

func NewF16_List(s *C.Segment, sz int) F16_List {
	return F16_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s F16_List) Len() int            { return C.PointerList(s).Len() }
func (s F16_List) At(i int) F16        { return F16(C.PointerList(s).At(i).ToStruct()) }
func (s F16_List) Set(i int, item F16) { C.PointerList(s).Set(i, C.Object(item)) }

type F16_Promise C.Pipeline

func (p *F16_Promise) Get() (F16, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return F16(s), err
}

func (p *F16_Promise) Base() *PlaneBase_Promise {
	return (*PlaneBase_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

type Regression C.Struct

func NewRegression(s *C.Segment) Regression {
	return Regression(s.NewStruct(C.ObjectSize{DataSize: 24, PointerCount: 3}))
}
func NewRootRegression(s *C.Segment) Regression {
	return Regression(s.NewRootStruct(C.ObjectSize{DataSize: 24, PointerCount: 3}))
}
func AutoNewRegression(s *C.Segment) Regression {
	return Regression(s.NewStructAR(C.ObjectSize{DataSize: 24, PointerCount: 3}))
}
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

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Regression) MarshalJSON() (bs []byte, err error) { return }

type Regression_List C.PointerList

func NewRegression_List(s *C.Segment, sz int) Regression_List {
	return Regression_List(s.NewCompositeList(C.ObjectSize{DataSize: 24, PointerCount: 3}, sz))
}
func (s Regression_List) Len() int                   { return C.PointerList(s).Len() }
func (s Regression_List) At(i int) Regression        { return Regression(C.PointerList(s).At(i).ToStruct()) }
func (s Regression_List) Set(i int, item Regression) { C.PointerList(s).Set(i, C.Object(item)) }

type Regression_Promise C.Pipeline

func (p *Regression_Promise) Get() (Regression, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Regression(s), err
}

func (p *Regression_Promise) Base() *PlaneBase_Promise {
	return (*PlaneBase_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

type Aircraft C.Struct
type Aircraft_Which uint16

const (
	Aircraft_Which_void Aircraft_Which = 0
	Aircraft_Which_b737 Aircraft_Which = 1
	Aircraft_Which_a320 Aircraft_Which = 2
	Aircraft_Which_f16  Aircraft_Which = 3
)

func NewAircraft(s *C.Segment) Aircraft {
	return Aircraft(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func NewRootAircraft(s *C.Segment) Aircraft {
	return Aircraft(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func AutoNewAircraft(s *C.Segment) Aircraft {
	return Aircraft(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func ReadRootAircraft(s *C.Segment) Aircraft { return Aircraft(s.Root(0).ToStruct()) }
func (s Aircraft) Which() Aircraft_Which     { return Aircraft_Which(C.Struct(s).Get16(0)) }
func (s Aircraft) SetVoid()                  { C.Struct(s).Set16(0, 0) }
func (s Aircraft) B737() B737                { return B737(C.Struct(s).GetObject(0).ToStruct()) }
func (s Aircraft) SetB737(v B737)            { C.Struct(s).Set16(0, 1); C.Struct(s).SetObject(0, C.Object(v)) }
func (s Aircraft) A320() A320                { return A320(C.Struct(s).GetObject(0).ToStruct()) }
func (s Aircraft) SetA320(v A320)            { C.Struct(s).Set16(0, 2); C.Struct(s).SetObject(0, C.Object(v)) }
func (s Aircraft) F16() F16                  { return F16(C.Struct(s).GetObject(0).ToStruct()) }
func (s Aircraft) SetF16(v F16)              { C.Struct(s).Set16(0, 3); C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Aircraft) MarshalJSON() (bs []byte, err error) { return }

type Aircraft_List C.PointerList

func NewAircraft_List(s *C.Segment, sz int) Aircraft_List {
	return Aircraft_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 1}, sz))
}
func (s Aircraft_List) Len() int                 { return C.PointerList(s).Len() }
func (s Aircraft_List) At(i int) Aircraft        { return Aircraft(C.PointerList(s).At(i).ToStruct()) }
func (s Aircraft_List) Set(i int, item Aircraft) { C.PointerList(s).Set(i, C.Object(item)) }

type Aircraft_Promise C.Pipeline

func (p *Aircraft_Promise) Get() (Aircraft, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Aircraft(s), err
}

func (p *Aircraft_Promise) B737() *B737_Promise {
	return (*B737_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *Aircraft_Promise) A320() *A320_Promise {
	return (*A320_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *Aircraft_Promise) F16() *F16_Promise {
	return (*F16_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

type Z C.Struct
type Z_Which uint16

const (
	Z_Which_void        Z_Which = 0
	Z_Which_zz          Z_Which = 1
	Z_Which_f64         Z_Which = 2
	Z_Which_f32         Z_Which = 3
	Z_Which_i64         Z_Which = 4
	Z_Which_i32         Z_Which = 5
	Z_Which_i16         Z_Which = 6
	Z_Which_i8          Z_Which = 7
	Z_Which_u64         Z_Which = 8
	Z_Which_u32         Z_Which = 9
	Z_Which_u16         Z_Which = 10
	Z_Which_u8          Z_Which = 11
	Z_Which_bool        Z_Which = 12
	Z_Which_text        Z_Which = 13
	Z_Which_blob        Z_Which = 14
	Z_Which_f64vec      Z_Which = 15
	Z_Which_f32vec      Z_Which = 16
	Z_Which_i64vec      Z_Which = 17
	Z_Which_i32vec      Z_Which = 18
	Z_Which_i16vec      Z_Which = 19
	Z_Which_i8vec       Z_Which = 20
	Z_Which_u64vec      Z_Which = 21
	Z_Which_u32vec      Z_Which = 22
	Z_Which_u16vec      Z_Which = 23
	Z_Which_u8vec       Z_Which = 24
	Z_Which_zvec        Z_Which = 25
	Z_Which_zvecvec     Z_Which = 26
	Z_Which_zdate       Z_Which = 27
	Z_Which_zdata       Z_Which = 28
	Z_Which_aircraftvec Z_Which = 29
	Z_Which_aircraft    Z_Which = 30
	Z_Which_regression  Z_Which = 31
	Z_Which_planebase   Z_Which = 32
	Z_Which_airport     Z_Which = 33
	Z_Which_b737        Z_Which = 34
	Z_Which_a320        Z_Which = 35
	Z_Which_f16         Z_Which = 36
	Z_Which_zdatevec    Z_Which = 37
	Z_Which_zdatavec    Z_Which = 38
	Z_Which_boolvec     Z_Which = 39
)

func NewZ(s *C.Segment) Z             { return Z(s.NewStruct(C.ObjectSize{DataSize: 16, PointerCount: 1})) }
func NewRootZ(s *C.Segment) Z         { return Z(s.NewRootStruct(C.ObjectSize{DataSize: 16, PointerCount: 1})) }
func AutoNewZ(s *C.Segment) Z         { return Z(s.NewStructAR(C.ObjectSize{DataSize: 16, PointerCount: 1})) }
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

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Z) MarshalJSON() (bs []byte, err error) { return }

type Z_List C.PointerList

func NewZ_List(s *C.Segment, sz int) Z_List {
	return Z_List(s.NewCompositeList(C.ObjectSize{DataSize: 16, PointerCount: 1}, sz))
}
func (s Z_List) Len() int          { return C.PointerList(s).Len() }
func (s Z_List) At(i int) Z        { return Z(C.PointerList(s).At(i).ToStruct()) }
func (s Z_List) Set(i int, item Z) { C.PointerList(s).Set(i, C.Object(item)) }

type Z_Promise C.Pipeline

func (p *Z_Promise) Get() (Z, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Z(s), err
}

func (p *Z_Promise) Zz() *Z_Promise {
	return (*Z_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *Z_Promise) Zdate() *Zdate_Promise {
	return (*Zdate_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *Z_Promise) Zdata() *Zdata_Promise {
	return (*Zdata_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *Z_Promise) Aircraft() *Aircraft_Promise {
	return (*Aircraft_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *Z_Promise) Regression() *Regression_Promise {
	return (*Regression_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *Z_Promise) Planebase() *PlaneBase_Promise {
	return (*PlaneBase_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *Z_Promise) B737() *B737_Promise {
	return (*B737_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *Z_Promise) A320() *A320_Promise {
	return (*A320_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *Z_Promise) F16() *F16_Promise {
	return (*F16_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

type Counter C.Struct

func NewCounter(s *C.Segment) Counter {
	return Counter(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 2}))
}
func NewRootCounter(s *C.Segment) Counter {
	return Counter(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 2}))
}
func AutoNewCounter(s *C.Segment) Counter {
	return Counter(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 2}))
}
func ReadRootCounter(s *C.Segment) Counter { return Counter(s.Root(0).ToStruct()) }
func (s Counter) Size() int64              { return int64(C.Struct(s).Get64(0)) }
func (s Counter) SetSize(v int64)          { C.Struct(s).Set64(0, uint64(v)) }
func (s Counter) Words() string            { return C.Struct(s).GetObject(0).ToText() }
func (s Counter) SetWords(v string)        { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s Counter) Wordlist() C.TextList     { return C.TextList(C.Struct(s).GetObject(1)) }
func (s Counter) SetWordlist(v C.TextList) { C.Struct(s).SetObject(1, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Counter) MarshalJSON() (bs []byte, err error) { return }

type Counter_List C.PointerList

func NewCounter_List(s *C.Segment, sz int) Counter_List {
	return Counter_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 2}, sz))
}
func (s Counter_List) Len() int                { return C.PointerList(s).Len() }
func (s Counter_List) At(i int) Counter        { return Counter(C.PointerList(s).At(i).ToStruct()) }
func (s Counter_List) Set(i int, item Counter) { C.PointerList(s).Set(i, C.Object(item)) }

type Counter_Promise C.Pipeline

func (p *Counter_Promise) Get() (Counter, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Counter(s), err
}

type Bag C.Struct

func NewBag(s *C.Segment) Bag { return Bag(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1})) }
func NewRootBag(s *C.Segment) Bag {
	return Bag(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewBag(s *C.Segment) Bag {
	return Bag(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootBag(s *C.Segment) Bag { return Bag(s.Root(0).ToStruct()) }
func (s Bag) Counter() Counter     { return Counter(C.Struct(s).GetObject(0).ToStruct()) }
func (s Bag) SetCounter(v Counter) { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Bag) MarshalJSON() (bs []byte, err error) { return }

type Bag_List C.PointerList

func NewBag_List(s *C.Segment, sz int) Bag_List {
	return Bag_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s Bag_List) Len() int            { return C.PointerList(s).Len() }
func (s Bag_List) At(i int) Bag        { return Bag(C.PointerList(s).At(i).ToStruct()) }
func (s Bag_List) Set(i int, item Bag) { C.PointerList(s).Set(i, C.Object(item)) }

type Bag_Promise C.Pipeline

func (p *Bag_Promise) Get() (Bag, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Bag(s), err
}

func (p *Bag_Promise) Counter() *Counter_Promise {
	return (*Counter_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

type Zserver C.Struct

func NewZserver(s *C.Segment) Zserver {
	return Zserver(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func NewRootZserver(s *C.Segment) Zserver {
	return Zserver(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewZserver(s *C.Segment) Zserver {
	return Zserver(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootZserver(s *C.Segment) Zserver   { return Zserver(s.Root(0).ToStruct()) }
func (s Zserver) Waitingjobs() Zjob_List     { return Zjob_List(C.Struct(s).GetObject(0)) }
func (s Zserver) SetWaitingjobs(v Zjob_List) { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Zserver) MarshalJSON() (bs []byte, err error) { return }

type Zserver_List C.PointerList

func NewZserver_List(s *C.Segment, sz int) Zserver_List {
	return Zserver_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s Zserver_List) Len() int                { return C.PointerList(s).Len() }
func (s Zserver_List) At(i int) Zserver        { return Zserver(C.PointerList(s).At(i).ToStruct()) }
func (s Zserver_List) Set(i int, item Zserver) { C.PointerList(s).Set(i, C.Object(item)) }

type Zserver_Promise C.Pipeline

func (p *Zserver_Promise) Get() (Zserver, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Zserver(s), err
}

type Zjob C.Struct

func NewZjob(s *C.Segment) Zjob { return Zjob(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 2})) }
func NewRootZjob(s *C.Segment) Zjob {
	return Zjob(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 2}))
}
func AutoNewZjob(s *C.Segment) Zjob {
	return Zjob(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 2}))
}
func ReadRootZjob(s *C.Segment) Zjob { return Zjob(s.Root(0).ToStruct()) }
func (s Zjob) Cmd() string           { return C.Struct(s).GetObject(0).ToText() }
func (s Zjob) SetCmd(v string)       { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s Zjob) Args() C.TextList      { return C.TextList(C.Struct(s).GetObject(1)) }
func (s Zjob) SetArgs(v C.TextList)  { C.Struct(s).SetObject(1, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Zjob) MarshalJSON() (bs []byte, err error) { return }

type Zjob_List C.PointerList

func NewZjob_List(s *C.Segment, sz int) Zjob_List {
	return Zjob_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 2}, sz))
}
func (s Zjob_List) Len() int             { return C.PointerList(s).Len() }
func (s Zjob_List) At(i int) Zjob        { return Zjob(C.PointerList(s).At(i).ToStruct()) }
func (s Zjob_List) Set(i int, item Zjob) { C.PointerList(s).Set(i, C.Object(item)) }

type Zjob_Promise C.Pipeline

func (p *Zjob_Promise) Get() (Zjob, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Zjob(s), err
}

type VerEmpty C.Struct

func NewVerEmpty(s *C.Segment) VerEmpty {
	return VerEmpty(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 0}))
}
func NewRootVerEmpty(s *C.Segment) VerEmpty {
	return VerEmpty(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 0}))
}
func AutoNewVerEmpty(s *C.Segment) VerEmpty {
	return VerEmpty(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 0}))
}
func ReadRootVerEmpty(s *C.Segment) VerEmpty { return VerEmpty(s.Root(0).ToStruct()) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s VerEmpty) MarshalJSON() (bs []byte, err error) { return }

type VerEmpty_List C.PointerList

func NewVerEmpty_List(s *C.Segment, sz int) VerEmpty_List {
	return VerEmpty_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 0}, sz))
}
func (s VerEmpty_List) Len() int                 { return C.PointerList(s).Len() }
func (s VerEmpty_List) At(i int) VerEmpty        { return VerEmpty(C.PointerList(s).At(i).ToStruct()) }
func (s VerEmpty_List) Set(i int, item VerEmpty) { C.PointerList(s).Set(i, C.Object(item)) }

type VerEmpty_Promise C.Pipeline

func (p *VerEmpty_Promise) Get() (VerEmpty, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return VerEmpty(s), err
}

type VerOneData C.Struct

func NewVerOneData(s *C.Segment) VerOneData {
	return VerOneData(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 0}))
}
func NewRootVerOneData(s *C.Segment) VerOneData {
	return VerOneData(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 0}))
}
func AutoNewVerOneData(s *C.Segment) VerOneData {
	return VerOneData(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 0}))
}
func ReadRootVerOneData(s *C.Segment) VerOneData { return VerOneData(s.Root(0).ToStruct()) }
func (s VerOneData) Val() int16                  { return int16(C.Struct(s).Get16(0)) }
func (s VerOneData) SetVal(v int16)              { C.Struct(s).Set16(0, uint16(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s VerOneData) MarshalJSON() (bs []byte, err error) { return }

type VerOneData_List C.PointerList

func NewVerOneData_List(s *C.Segment, sz int) VerOneData_List {
	return VerOneData_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 0}, sz))
}
func (s VerOneData_List) Len() int                   { return C.PointerList(s).Len() }
func (s VerOneData_List) At(i int) VerOneData        { return VerOneData(C.PointerList(s).At(i).ToStruct()) }
func (s VerOneData_List) Set(i int, item VerOneData) { C.PointerList(s).Set(i, C.Object(item)) }

type VerOneData_Promise C.Pipeline

func (p *VerOneData_Promise) Get() (VerOneData, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return VerOneData(s), err
}

type VerTwoData C.Struct

func NewVerTwoData(s *C.Segment) VerTwoData {
	return VerTwoData(s.NewStruct(C.ObjectSize{DataSize: 16, PointerCount: 0}))
}
func NewRootVerTwoData(s *C.Segment) VerTwoData {
	return VerTwoData(s.NewRootStruct(C.ObjectSize{DataSize: 16, PointerCount: 0}))
}
func AutoNewVerTwoData(s *C.Segment) VerTwoData {
	return VerTwoData(s.NewStructAR(C.ObjectSize{DataSize: 16, PointerCount: 0}))
}
func ReadRootVerTwoData(s *C.Segment) VerTwoData { return VerTwoData(s.Root(0).ToStruct()) }
func (s VerTwoData) Val() int16                  { return int16(C.Struct(s).Get16(0)) }
func (s VerTwoData) SetVal(v int16)              { C.Struct(s).Set16(0, uint16(v)) }
func (s VerTwoData) Duo() int64                  { return int64(C.Struct(s).Get64(8)) }
func (s VerTwoData) SetDuo(v int64)              { C.Struct(s).Set64(8, uint64(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s VerTwoData) MarshalJSON() (bs []byte, err error) { return }

type VerTwoData_List C.PointerList

func NewVerTwoData_List(s *C.Segment, sz int) VerTwoData_List {
	return VerTwoData_List(s.NewCompositeList(C.ObjectSize{DataSize: 16, PointerCount: 0}, sz))
}
func (s VerTwoData_List) Len() int                   { return C.PointerList(s).Len() }
func (s VerTwoData_List) At(i int) VerTwoData        { return VerTwoData(C.PointerList(s).At(i).ToStruct()) }
func (s VerTwoData_List) Set(i int, item VerTwoData) { C.PointerList(s).Set(i, C.Object(item)) }

type VerTwoData_Promise C.Pipeline

func (p *VerTwoData_Promise) Get() (VerTwoData, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return VerTwoData(s), err
}

type VerOnePtr C.Struct

func NewVerOnePtr(s *C.Segment) VerOnePtr {
	return VerOnePtr(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func NewRootVerOnePtr(s *C.Segment) VerOnePtr {
	return VerOnePtr(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewVerOnePtr(s *C.Segment) VerOnePtr {
	return VerOnePtr(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootVerOnePtr(s *C.Segment) VerOnePtr { return VerOnePtr(s.Root(0).ToStruct()) }
func (s VerOnePtr) Ptr() VerOneData            { return VerOneData(C.Struct(s).GetObject(0).ToStruct()) }
func (s VerOnePtr) SetPtr(v VerOneData)        { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s VerOnePtr) MarshalJSON() (bs []byte, err error) { return }

type VerOnePtr_List C.PointerList

func NewVerOnePtr_List(s *C.Segment, sz int) VerOnePtr_List {
	return VerOnePtr_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s VerOnePtr_List) Len() int                  { return C.PointerList(s).Len() }
func (s VerOnePtr_List) At(i int) VerOnePtr        { return VerOnePtr(C.PointerList(s).At(i).ToStruct()) }
func (s VerOnePtr_List) Set(i int, item VerOnePtr) { C.PointerList(s).Set(i, C.Object(item)) }

type VerOnePtr_Promise C.Pipeline

func (p *VerOnePtr_Promise) Get() (VerOnePtr, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return VerOnePtr(s), err
}

func (p *VerOnePtr_Promise) Ptr() *VerOneData_Promise {
	return (*VerOneData_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

type VerTwoPtr C.Struct

func NewVerTwoPtr(s *C.Segment) VerTwoPtr {
	return VerTwoPtr(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 2}))
}
func NewRootVerTwoPtr(s *C.Segment) VerTwoPtr {
	return VerTwoPtr(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 2}))
}
func AutoNewVerTwoPtr(s *C.Segment) VerTwoPtr {
	return VerTwoPtr(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 2}))
}
func ReadRootVerTwoPtr(s *C.Segment) VerTwoPtr { return VerTwoPtr(s.Root(0).ToStruct()) }
func (s VerTwoPtr) Ptr1() VerOneData           { return VerOneData(C.Struct(s).GetObject(0).ToStruct()) }
func (s VerTwoPtr) SetPtr1(v VerOneData)       { C.Struct(s).SetObject(0, C.Object(v)) }
func (s VerTwoPtr) Ptr2() VerOneData           { return VerOneData(C.Struct(s).GetObject(1).ToStruct()) }
func (s VerTwoPtr) SetPtr2(v VerOneData)       { C.Struct(s).SetObject(1, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s VerTwoPtr) MarshalJSON() (bs []byte, err error) { return }

type VerTwoPtr_List C.PointerList

func NewVerTwoPtr_List(s *C.Segment, sz int) VerTwoPtr_List {
	return VerTwoPtr_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 2}, sz))
}
func (s VerTwoPtr_List) Len() int                  { return C.PointerList(s).Len() }
func (s VerTwoPtr_List) At(i int) VerTwoPtr        { return VerTwoPtr(C.PointerList(s).At(i).ToStruct()) }
func (s VerTwoPtr_List) Set(i int, item VerTwoPtr) { C.PointerList(s).Set(i, C.Object(item)) }

type VerTwoPtr_Promise C.Pipeline

func (p *VerTwoPtr_Promise) Get() (VerTwoPtr, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return VerTwoPtr(s), err
}

func (p *VerTwoPtr_Promise) Ptr1() *VerOneData_Promise {
	return (*VerOneData_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *VerTwoPtr_Promise) Ptr2() *VerOneData_Promise {
	return (*VerOneData_Promise)((*C.Pipeline)(p).GetPipeline(1))
}

type VerTwoDataTwoPtr C.Struct

func NewVerTwoDataTwoPtr(s *C.Segment) VerTwoDataTwoPtr {
	return VerTwoDataTwoPtr(s.NewStruct(C.ObjectSize{DataSize: 16, PointerCount: 2}))
}
func NewRootVerTwoDataTwoPtr(s *C.Segment) VerTwoDataTwoPtr {
	return VerTwoDataTwoPtr(s.NewRootStruct(C.ObjectSize{DataSize: 16, PointerCount: 2}))
}
func AutoNewVerTwoDataTwoPtr(s *C.Segment) VerTwoDataTwoPtr {
	return VerTwoDataTwoPtr(s.NewStructAR(C.ObjectSize{DataSize: 16, PointerCount: 2}))
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

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s VerTwoDataTwoPtr) MarshalJSON() (bs []byte, err error) { return }

type VerTwoDataTwoPtr_List C.PointerList

func NewVerTwoDataTwoPtr_List(s *C.Segment, sz int) VerTwoDataTwoPtr_List {
	return VerTwoDataTwoPtr_List(s.NewCompositeList(C.ObjectSize{DataSize: 16, PointerCount: 2}, sz))
}
func (s VerTwoDataTwoPtr_List) Len() int { return C.PointerList(s).Len() }
func (s VerTwoDataTwoPtr_List) At(i int) VerTwoDataTwoPtr {
	return VerTwoDataTwoPtr(C.PointerList(s).At(i).ToStruct())
}
func (s VerTwoDataTwoPtr_List) Set(i int, item VerTwoDataTwoPtr) {
	C.PointerList(s).Set(i, C.Object(item))
}

type VerTwoDataTwoPtr_Promise C.Pipeline

func (p *VerTwoDataTwoPtr_Promise) Get() (VerTwoDataTwoPtr, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return VerTwoDataTwoPtr(s), err
}

func (p *VerTwoDataTwoPtr_Promise) Ptr1() *VerOneData_Promise {
	return (*VerOneData_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *VerTwoDataTwoPtr_Promise) Ptr2() *VerOneData_Promise {
	return (*VerOneData_Promise)((*C.Pipeline)(p).GetPipeline(1))
}

type HoldsVerEmptyList C.Struct

func NewHoldsVerEmptyList(s *C.Segment) HoldsVerEmptyList {
	return HoldsVerEmptyList(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func NewRootHoldsVerEmptyList(s *C.Segment) HoldsVerEmptyList {
	return HoldsVerEmptyList(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewHoldsVerEmptyList(s *C.Segment) HoldsVerEmptyList {
	return HoldsVerEmptyList(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootHoldsVerEmptyList(s *C.Segment) HoldsVerEmptyList {
	return HoldsVerEmptyList(s.Root(0).ToStruct())
}
func (s HoldsVerEmptyList) Mylist() VerEmpty_List     { return VerEmpty_List(C.Struct(s).GetObject(0)) }
func (s HoldsVerEmptyList) SetMylist(v VerEmpty_List) { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s HoldsVerEmptyList) MarshalJSON() (bs []byte, err error) { return }

type HoldsVerEmptyList_List C.PointerList

func NewHoldsVerEmptyList_List(s *C.Segment, sz int) HoldsVerEmptyList_List {
	return HoldsVerEmptyList_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s HoldsVerEmptyList_List) Len() int { return C.PointerList(s).Len() }
func (s HoldsVerEmptyList_List) At(i int) HoldsVerEmptyList {
	return HoldsVerEmptyList(C.PointerList(s).At(i).ToStruct())
}
func (s HoldsVerEmptyList_List) Set(i int, item HoldsVerEmptyList) {
	C.PointerList(s).Set(i, C.Object(item))
}

type HoldsVerEmptyList_Promise C.Pipeline

func (p *HoldsVerEmptyList_Promise) Get() (HoldsVerEmptyList, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return HoldsVerEmptyList(s), err
}

type HoldsVerOneDataList C.Struct

func NewHoldsVerOneDataList(s *C.Segment) HoldsVerOneDataList {
	return HoldsVerOneDataList(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func NewRootHoldsVerOneDataList(s *C.Segment) HoldsVerOneDataList {
	return HoldsVerOneDataList(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewHoldsVerOneDataList(s *C.Segment) HoldsVerOneDataList {
	return HoldsVerOneDataList(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootHoldsVerOneDataList(s *C.Segment) HoldsVerOneDataList {
	return HoldsVerOneDataList(s.Root(0).ToStruct())
}
func (s HoldsVerOneDataList) Mylist() VerOneData_List {
	return VerOneData_List(C.Struct(s).GetObject(0))
}
func (s HoldsVerOneDataList) SetMylist(v VerOneData_List) { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s HoldsVerOneDataList) MarshalJSON() (bs []byte, err error) { return }

type HoldsVerOneDataList_List C.PointerList

func NewHoldsVerOneDataList_List(s *C.Segment, sz int) HoldsVerOneDataList_List {
	return HoldsVerOneDataList_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s HoldsVerOneDataList_List) Len() int { return C.PointerList(s).Len() }
func (s HoldsVerOneDataList_List) At(i int) HoldsVerOneDataList {
	return HoldsVerOneDataList(C.PointerList(s).At(i).ToStruct())
}
func (s HoldsVerOneDataList_List) Set(i int, item HoldsVerOneDataList) {
	C.PointerList(s).Set(i, C.Object(item))
}

type HoldsVerOneDataList_Promise C.Pipeline

func (p *HoldsVerOneDataList_Promise) Get() (HoldsVerOneDataList, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return HoldsVerOneDataList(s), err
}

type HoldsVerTwoDataList C.Struct

func NewHoldsVerTwoDataList(s *C.Segment) HoldsVerTwoDataList {
	return HoldsVerTwoDataList(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func NewRootHoldsVerTwoDataList(s *C.Segment) HoldsVerTwoDataList {
	return HoldsVerTwoDataList(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewHoldsVerTwoDataList(s *C.Segment) HoldsVerTwoDataList {
	return HoldsVerTwoDataList(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootHoldsVerTwoDataList(s *C.Segment) HoldsVerTwoDataList {
	return HoldsVerTwoDataList(s.Root(0).ToStruct())
}
func (s HoldsVerTwoDataList) Mylist() VerTwoData_List {
	return VerTwoData_List(C.Struct(s).GetObject(0))
}
func (s HoldsVerTwoDataList) SetMylist(v VerTwoData_List) { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s HoldsVerTwoDataList) MarshalJSON() (bs []byte, err error) { return }

type HoldsVerTwoDataList_List C.PointerList

func NewHoldsVerTwoDataList_List(s *C.Segment, sz int) HoldsVerTwoDataList_List {
	return HoldsVerTwoDataList_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s HoldsVerTwoDataList_List) Len() int { return C.PointerList(s).Len() }
func (s HoldsVerTwoDataList_List) At(i int) HoldsVerTwoDataList {
	return HoldsVerTwoDataList(C.PointerList(s).At(i).ToStruct())
}
func (s HoldsVerTwoDataList_List) Set(i int, item HoldsVerTwoDataList) {
	C.PointerList(s).Set(i, C.Object(item))
}

type HoldsVerTwoDataList_Promise C.Pipeline

func (p *HoldsVerTwoDataList_Promise) Get() (HoldsVerTwoDataList, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return HoldsVerTwoDataList(s), err
}

type HoldsVerOnePtrList C.Struct

func NewHoldsVerOnePtrList(s *C.Segment) HoldsVerOnePtrList {
	return HoldsVerOnePtrList(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func NewRootHoldsVerOnePtrList(s *C.Segment) HoldsVerOnePtrList {
	return HoldsVerOnePtrList(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewHoldsVerOnePtrList(s *C.Segment) HoldsVerOnePtrList {
	return HoldsVerOnePtrList(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootHoldsVerOnePtrList(s *C.Segment) HoldsVerOnePtrList {
	return HoldsVerOnePtrList(s.Root(0).ToStruct())
}
func (s HoldsVerOnePtrList) Mylist() VerOnePtr_List     { return VerOnePtr_List(C.Struct(s).GetObject(0)) }
func (s HoldsVerOnePtrList) SetMylist(v VerOnePtr_List) { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s HoldsVerOnePtrList) MarshalJSON() (bs []byte, err error) { return }

type HoldsVerOnePtrList_List C.PointerList

func NewHoldsVerOnePtrList_List(s *C.Segment, sz int) HoldsVerOnePtrList_List {
	return HoldsVerOnePtrList_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s HoldsVerOnePtrList_List) Len() int { return C.PointerList(s).Len() }
func (s HoldsVerOnePtrList_List) At(i int) HoldsVerOnePtrList {
	return HoldsVerOnePtrList(C.PointerList(s).At(i).ToStruct())
}
func (s HoldsVerOnePtrList_List) Set(i int, item HoldsVerOnePtrList) {
	C.PointerList(s).Set(i, C.Object(item))
}

type HoldsVerOnePtrList_Promise C.Pipeline

func (p *HoldsVerOnePtrList_Promise) Get() (HoldsVerOnePtrList, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return HoldsVerOnePtrList(s), err
}

type HoldsVerTwoPtrList C.Struct

func NewHoldsVerTwoPtrList(s *C.Segment) HoldsVerTwoPtrList {
	return HoldsVerTwoPtrList(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func NewRootHoldsVerTwoPtrList(s *C.Segment) HoldsVerTwoPtrList {
	return HoldsVerTwoPtrList(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewHoldsVerTwoPtrList(s *C.Segment) HoldsVerTwoPtrList {
	return HoldsVerTwoPtrList(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootHoldsVerTwoPtrList(s *C.Segment) HoldsVerTwoPtrList {
	return HoldsVerTwoPtrList(s.Root(0).ToStruct())
}
func (s HoldsVerTwoPtrList) Mylist() VerTwoPtr_List     { return VerTwoPtr_List(C.Struct(s).GetObject(0)) }
func (s HoldsVerTwoPtrList) SetMylist(v VerTwoPtr_List) { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s HoldsVerTwoPtrList) MarshalJSON() (bs []byte, err error) { return }

type HoldsVerTwoPtrList_List C.PointerList

func NewHoldsVerTwoPtrList_List(s *C.Segment, sz int) HoldsVerTwoPtrList_List {
	return HoldsVerTwoPtrList_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s HoldsVerTwoPtrList_List) Len() int { return C.PointerList(s).Len() }
func (s HoldsVerTwoPtrList_List) At(i int) HoldsVerTwoPtrList {
	return HoldsVerTwoPtrList(C.PointerList(s).At(i).ToStruct())
}
func (s HoldsVerTwoPtrList_List) Set(i int, item HoldsVerTwoPtrList) {
	C.PointerList(s).Set(i, C.Object(item))
}

type HoldsVerTwoPtrList_Promise C.Pipeline

func (p *HoldsVerTwoPtrList_Promise) Get() (HoldsVerTwoPtrList, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return HoldsVerTwoPtrList(s), err
}

type HoldsVerTwoTwoList C.Struct

func NewHoldsVerTwoTwoList(s *C.Segment) HoldsVerTwoTwoList {
	return HoldsVerTwoTwoList(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func NewRootHoldsVerTwoTwoList(s *C.Segment) HoldsVerTwoTwoList {
	return HoldsVerTwoTwoList(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewHoldsVerTwoTwoList(s *C.Segment) HoldsVerTwoTwoList {
	return HoldsVerTwoTwoList(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootHoldsVerTwoTwoList(s *C.Segment) HoldsVerTwoTwoList {
	return HoldsVerTwoTwoList(s.Root(0).ToStruct())
}
func (s HoldsVerTwoTwoList) Mylist() VerTwoDataTwoPtr_List {
	return VerTwoDataTwoPtr_List(C.Struct(s).GetObject(0))
}
func (s HoldsVerTwoTwoList) SetMylist(v VerTwoDataTwoPtr_List) { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s HoldsVerTwoTwoList) MarshalJSON() (bs []byte, err error) { return }

type HoldsVerTwoTwoList_List C.PointerList

func NewHoldsVerTwoTwoList_List(s *C.Segment, sz int) HoldsVerTwoTwoList_List {
	return HoldsVerTwoTwoList_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s HoldsVerTwoTwoList_List) Len() int { return C.PointerList(s).Len() }
func (s HoldsVerTwoTwoList_List) At(i int) HoldsVerTwoTwoList {
	return HoldsVerTwoTwoList(C.PointerList(s).At(i).ToStruct())
}
func (s HoldsVerTwoTwoList_List) Set(i int, item HoldsVerTwoTwoList) {
	C.PointerList(s).Set(i, C.Object(item))
}

type HoldsVerTwoTwoList_Promise C.Pipeline

func (p *HoldsVerTwoTwoList_Promise) Get() (HoldsVerTwoTwoList, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return HoldsVerTwoTwoList(s), err
}

type HoldsVerTwoTwoPlus C.Struct

func NewHoldsVerTwoTwoPlus(s *C.Segment) HoldsVerTwoTwoPlus {
	return HoldsVerTwoTwoPlus(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func NewRootHoldsVerTwoTwoPlus(s *C.Segment) HoldsVerTwoTwoPlus {
	return HoldsVerTwoTwoPlus(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewHoldsVerTwoTwoPlus(s *C.Segment) HoldsVerTwoTwoPlus {
	return HoldsVerTwoTwoPlus(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootHoldsVerTwoTwoPlus(s *C.Segment) HoldsVerTwoTwoPlus {
	return HoldsVerTwoTwoPlus(s.Root(0).ToStruct())
}
func (s HoldsVerTwoTwoPlus) Mylist() VerTwoTwoPlus_List {
	return VerTwoTwoPlus_List(C.Struct(s).GetObject(0))
}
func (s HoldsVerTwoTwoPlus) SetMylist(v VerTwoTwoPlus_List) { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s HoldsVerTwoTwoPlus) MarshalJSON() (bs []byte, err error) { return }

type HoldsVerTwoTwoPlus_List C.PointerList

func NewHoldsVerTwoTwoPlus_List(s *C.Segment, sz int) HoldsVerTwoTwoPlus_List {
	return HoldsVerTwoTwoPlus_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s HoldsVerTwoTwoPlus_List) Len() int { return C.PointerList(s).Len() }
func (s HoldsVerTwoTwoPlus_List) At(i int) HoldsVerTwoTwoPlus {
	return HoldsVerTwoTwoPlus(C.PointerList(s).At(i).ToStruct())
}
func (s HoldsVerTwoTwoPlus_List) Set(i int, item HoldsVerTwoTwoPlus) {
	C.PointerList(s).Set(i, C.Object(item))
}

type HoldsVerTwoTwoPlus_Promise C.Pipeline

func (p *HoldsVerTwoTwoPlus_Promise) Get() (HoldsVerTwoTwoPlus, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return HoldsVerTwoTwoPlus(s), err
}

type VerTwoTwoPlus C.Struct

func NewVerTwoTwoPlus(s *C.Segment) VerTwoTwoPlus {
	return VerTwoTwoPlus(s.NewStruct(C.ObjectSize{DataSize: 24, PointerCount: 3}))
}
func NewRootVerTwoTwoPlus(s *C.Segment) VerTwoTwoPlus {
	return VerTwoTwoPlus(s.NewRootStruct(C.ObjectSize{DataSize: 24, PointerCount: 3}))
}
func AutoNewVerTwoTwoPlus(s *C.Segment) VerTwoTwoPlus {
	return VerTwoTwoPlus(s.NewStructAR(C.ObjectSize{DataSize: 24, PointerCount: 3}))
}
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

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s VerTwoTwoPlus) MarshalJSON() (bs []byte, err error) { return }

type VerTwoTwoPlus_List C.PointerList

func NewVerTwoTwoPlus_List(s *C.Segment, sz int) VerTwoTwoPlus_List {
	return VerTwoTwoPlus_List(s.NewCompositeList(C.ObjectSize{DataSize: 24, PointerCount: 3}, sz))
}
func (s VerTwoTwoPlus_List) Len() int { return C.PointerList(s).Len() }
func (s VerTwoTwoPlus_List) At(i int) VerTwoTwoPlus {
	return VerTwoTwoPlus(C.PointerList(s).At(i).ToStruct())
}
func (s VerTwoTwoPlus_List) Set(i int, item VerTwoTwoPlus) { C.PointerList(s).Set(i, C.Object(item)) }

type VerTwoTwoPlus_Promise C.Pipeline

func (p *VerTwoTwoPlus_Promise) Get() (VerTwoTwoPlus, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return VerTwoTwoPlus(s), err
}

func (p *VerTwoTwoPlus_Promise) Ptr1() *VerTwoDataTwoPtr_Promise {
	return (*VerTwoDataTwoPtr_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

func (p *VerTwoTwoPlus_Promise) Ptr2() *VerTwoDataTwoPtr_Promise {
	return (*VerTwoDataTwoPtr_Promise)((*C.Pipeline)(p).GetPipeline(1))
}

type HoldsText C.Struct

func NewHoldsText(s *C.Segment) HoldsText {
	return HoldsText(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 3}))
}
func NewRootHoldsText(s *C.Segment) HoldsText {
	return HoldsText(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 3}))
}
func AutoNewHoldsText(s *C.Segment) HoldsText {
	return HoldsText(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 3}))
}
func ReadRootHoldsText(s *C.Segment) HoldsText { return HoldsText(s.Root(0).ToStruct()) }
func (s HoldsText) Txt() string                { return C.Struct(s).GetObject(0).ToText() }
func (s HoldsText) SetTxt(v string)            { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s HoldsText) Lst() C.TextList            { return C.TextList(C.Struct(s).GetObject(1)) }
func (s HoldsText) SetLst(v C.TextList)        { C.Struct(s).SetObject(1, C.Object(v)) }
func (s HoldsText) Lstlst() C.PointerList      { return C.PointerList(C.Struct(s).GetObject(2)) }
func (s HoldsText) SetLstlst(v C.PointerList)  { C.Struct(s).SetObject(2, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s HoldsText) MarshalJSON() (bs []byte, err error) { return }

type HoldsText_List C.PointerList

func NewHoldsText_List(s *C.Segment, sz int) HoldsText_List {
	return HoldsText_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 3}, sz))
}
func (s HoldsText_List) Len() int                  { return C.PointerList(s).Len() }
func (s HoldsText_List) At(i int) HoldsText        { return HoldsText(C.PointerList(s).At(i).ToStruct()) }
func (s HoldsText_List) Set(i int, item HoldsText) { C.PointerList(s).Set(i, C.Object(item)) }

type HoldsText_Promise C.Pipeline

func (p *HoldsText_Promise) Get() (HoldsText, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return HoldsText(s), err
}

type WrapEmpty C.Struct

func NewWrapEmpty(s *C.Segment) WrapEmpty {
	return WrapEmpty(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func NewRootWrapEmpty(s *C.Segment) WrapEmpty {
	return WrapEmpty(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewWrapEmpty(s *C.Segment) WrapEmpty {
	return WrapEmpty(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootWrapEmpty(s *C.Segment) WrapEmpty { return WrapEmpty(s.Root(0).ToStruct()) }
func (s WrapEmpty) MightNotBeReallyEmpty() VerEmpty {
	return VerEmpty(C.Struct(s).GetObject(0).ToStruct())
}
func (s WrapEmpty) SetMightNotBeReallyEmpty(v VerEmpty) { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s WrapEmpty) MarshalJSON() (bs []byte, err error) { return }

type WrapEmpty_List C.PointerList

func NewWrapEmpty_List(s *C.Segment, sz int) WrapEmpty_List {
	return WrapEmpty_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s WrapEmpty_List) Len() int                  { return C.PointerList(s).Len() }
func (s WrapEmpty_List) At(i int) WrapEmpty        { return WrapEmpty(C.PointerList(s).At(i).ToStruct()) }
func (s WrapEmpty_List) Set(i int, item WrapEmpty) { C.PointerList(s).Set(i, C.Object(item)) }

type WrapEmpty_Promise C.Pipeline

func (p *WrapEmpty_Promise) Get() (WrapEmpty, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return WrapEmpty(s), err
}

func (p *WrapEmpty_Promise) MightNotBeReallyEmpty() *VerEmpty_Promise {
	return (*VerEmpty_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

type Wrap2x2 C.Struct

func NewWrap2x2(s *C.Segment) Wrap2x2 {
	return Wrap2x2(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func NewRootWrap2x2(s *C.Segment) Wrap2x2 {
	return Wrap2x2(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewWrap2x2(s *C.Segment) Wrap2x2 {
	return Wrap2x2(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootWrap2x2(s *C.Segment) Wrap2x2 { return Wrap2x2(s.Root(0).ToStruct()) }
func (s Wrap2x2) MightNotBeReallyEmpty() VerTwoDataTwoPtr {
	return VerTwoDataTwoPtr(C.Struct(s).GetObject(0).ToStruct())
}
func (s Wrap2x2) SetMightNotBeReallyEmpty(v VerTwoDataTwoPtr) { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Wrap2x2) MarshalJSON() (bs []byte, err error) { return }

type Wrap2x2_List C.PointerList

func NewWrap2x2_List(s *C.Segment, sz int) Wrap2x2_List {
	return Wrap2x2_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s Wrap2x2_List) Len() int                { return C.PointerList(s).Len() }
func (s Wrap2x2_List) At(i int) Wrap2x2        { return Wrap2x2(C.PointerList(s).At(i).ToStruct()) }
func (s Wrap2x2_List) Set(i int, item Wrap2x2) { C.PointerList(s).Set(i, C.Object(item)) }

type Wrap2x2_Promise C.Pipeline

func (p *Wrap2x2_Promise) Get() (Wrap2x2, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Wrap2x2(s), err
}

func (p *Wrap2x2_Promise) MightNotBeReallyEmpty() *VerTwoDataTwoPtr_Promise {
	return (*VerTwoDataTwoPtr_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

type Wrap2x2plus C.Struct

func NewWrap2x2plus(s *C.Segment) Wrap2x2plus {
	return Wrap2x2plus(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func NewRootWrap2x2plus(s *C.Segment) Wrap2x2plus {
	return Wrap2x2plus(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewWrap2x2plus(s *C.Segment) Wrap2x2plus {
	return Wrap2x2plus(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootWrap2x2plus(s *C.Segment) Wrap2x2plus { return Wrap2x2plus(s.Root(0).ToStruct()) }
func (s Wrap2x2plus) MightNotBeReallyEmpty() VerTwoTwoPlus {
	return VerTwoTwoPlus(C.Struct(s).GetObject(0).ToStruct())
}
func (s Wrap2x2plus) SetMightNotBeReallyEmpty(v VerTwoTwoPlus) { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Wrap2x2plus) MarshalJSON() (bs []byte, err error) { return }

type Wrap2x2plus_List C.PointerList

func NewWrap2x2plus_List(s *C.Segment, sz int) Wrap2x2plus_List {
	return Wrap2x2plus_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s Wrap2x2plus_List) Len() int                    { return C.PointerList(s).Len() }
func (s Wrap2x2plus_List) At(i int) Wrap2x2plus        { return Wrap2x2plus(C.PointerList(s).At(i).ToStruct()) }
func (s Wrap2x2plus_List) Set(i int, item Wrap2x2plus) { C.PointerList(s).Set(i, C.Object(item)) }

type Wrap2x2plus_Promise C.Pipeline

func (p *Wrap2x2plus_Promise) Get() (Wrap2x2plus, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Wrap2x2plus(s), err
}

func (p *Wrap2x2plus_Promise) MightNotBeReallyEmpty() *VerTwoTwoPlus_Promise {
	return (*VerTwoTwoPlus_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

type Endpoint C.Struct

func NewEndpoint(s *C.Segment) Endpoint {
	return Endpoint(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 2}))
}
func NewRootEndpoint(s *C.Segment) Endpoint {
	return Endpoint(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 2}))
}
func AutoNewEndpoint(s *C.Segment) Endpoint {
	return Endpoint(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 2}))
}
func ReadRootEndpoint(s *C.Segment) Endpoint { return Endpoint(s.Root(0).ToStruct()) }
func (s Endpoint) Ip() net.IP                { return net.IP(C.Struct(s).GetObject(0).ToData()) }
func (s Endpoint) SetIp(v net.IP)            { C.Struct(s).SetObject(0, s.Segment.NewData([]byte(v))) }
func (s Endpoint) Port() int16               { return int16(C.Struct(s).Get16(0)) }
func (s Endpoint) SetPort(v int16)           { C.Struct(s).Set16(0, uint16(v)) }
func (s Endpoint) Hostname() string          { return C.Struct(s).GetObject(1).ToText() }
func (s Endpoint) SetHostname(v string)      { C.Struct(s).SetObject(1, s.Segment.NewText(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Endpoint) MarshalJSON() (bs []byte, err error) { return }

type Endpoint_List C.PointerList

func NewEndpoint_List(s *C.Segment, sz int) Endpoint_List {
	return Endpoint_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 2}, sz))
}
func (s Endpoint_List) Len() int                 { return C.PointerList(s).Len() }
func (s Endpoint_List) At(i int) Endpoint        { return Endpoint(C.PointerList(s).At(i).ToStruct()) }
func (s Endpoint_List) Set(i int, item Endpoint) { C.PointerList(s).Set(i, C.Object(item)) }

type Endpoint_Promise C.Pipeline

func (p *Endpoint_Promise) Get() (Endpoint, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Endpoint(s), err
}

type VoidUnion C.Struct
type VoidUnion_Which uint16

const (
	VoidUnion_Which_a VoidUnion_Which = 0
	VoidUnion_Which_b VoidUnion_Which = 1
)

func NewVoidUnion(s *C.Segment) VoidUnion {
	return VoidUnion(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 0}))
}
func NewRootVoidUnion(s *C.Segment) VoidUnion {
	return VoidUnion(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 0}))
}
func AutoNewVoidUnion(s *C.Segment) VoidUnion {
	return VoidUnion(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 0}))
}
func ReadRootVoidUnion(s *C.Segment) VoidUnion { return VoidUnion(s.Root(0).ToStruct()) }
func (s VoidUnion) Which() VoidUnion_Which     { return VoidUnion_Which(C.Struct(s).Get16(0)) }
func (s VoidUnion) SetA()                      { C.Struct(s).Set16(0, 0) }
func (s VoidUnion) SetB()                      { C.Struct(s).Set16(0, 1) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s VoidUnion) MarshalJSON() (bs []byte, err error) { return }

type VoidUnion_List C.PointerList

func NewVoidUnion_List(s *C.Segment, sz int) VoidUnion_List {
	return VoidUnion_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 0}, sz))
}
func (s VoidUnion_List) Len() int                  { return C.PointerList(s).Len() }
func (s VoidUnion_List) At(i int) VoidUnion        { return VoidUnion(C.PointerList(s).At(i).ToStruct()) }
func (s VoidUnion_List) Set(i int, item VoidUnion) { C.PointerList(s).Set(i, C.Object(item)) }

type VoidUnion_Promise C.Pipeline

func (p *VoidUnion_Promise) Get() (VoidUnion, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return VoidUnion(s), err
}

type Nester1Capn C.Struct

func NewNester1Capn(s *C.Segment) Nester1Capn {
	return Nester1Capn(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func NewRootNester1Capn(s *C.Segment) Nester1Capn {
	return Nester1Capn(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewNester1Capn(s *C.Segment) Nester1Capn {
	return Nester1Capn(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootNester1Capn(s *C.Segment) Nester1Capn { return Nester1Capn(s.Root(0).ToStruct()) }
func (s Nester1Capn) Strs() C.TextList             { return C.TextList(C.Struct(s).GetObject(0)) }
func (s Nester1Capn) SetStrs(v C.TextList)         { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Nester1Capn) MarshalJSON() (bs []byte, err error) { return }

type Nester1Capn_List C.PointerList

func NewNester1Capn_List(s *C.Segment, sz int) Nester1Capn_List {
	return Nester1Capn_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s Nester1Capn_List) Len() int                    { return C.PointerList(s).Len() }
func (s Nester1Capn_List) At(i int) Nester1Capn        { return Nester1Capn(C.PointerList(s).At(i).ToStruct()) }
func (s Nester1Capn_List) Set(i int, item Nester1Capn) { C.PointerList(s).Set(i, C.Object(item)) }

type Nester1Capn_Promise C.Pipeline

func (p *Nester1Capn_Promise) Get() (Nester1Capn, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Nester1Capn(s), err
}

type RWTestCapn C.Struct

func NewRWTestCapn(s *C.Segment) RWTestCapn {
	return RWTestCapn(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func NewRootRWTestCapn(s *C.Segment) RWTestCapn {
	return RWTestCapn(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewRWTestCapn(s *C.Segment) RWTestCapn {
	return RWTestCapn(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootRWTestCapn(s *C.Segment) RWTestCapn   { return RWTestCapn(s.Root(0).ToStruct()) }
func (s RWTestCapn) NestMatrix() C.PointerList     { return C.PointerList(C.Struct(s).GetObject(0)) }
func (s RWTestCapn) SetNestMatrix(v C.PointerList) { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s RWTestCapn) MarshalJSON() (bs []byte, err error) { return }

type RWTestCapn_List C.PointerList

func NewRWTestCapn_List(s *C.Segment, sz int) RWTestCapn_List {
	return RWTestCapn_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s RWTestCapn_List) Len() int                   { return C.PointerList(s).Len() }
func (s RWTestCapn_List) At(i int) RWTestCapn        { return RWTestCapn(C.PointerList(s).At(i).ToStruct()) }
func (s RWTestCapn_List) Set(i int, item RWTestCapn) { C.PointerList(s).Set(i, C.Object(item)) }

type RWTestCapn_Promise C.Pipeline

func (p *RWTestCapn_Promise) Get() (RWTestCapn, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return RWTestCapn(s), err
}

type ListStructCapn C.Struct

func NewListStructCapn(s *C.Segment) ListStructCapn {
	return ListStructCapn(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func NewRootListStructCapn(s *C.Segment) ListStructCapn {
	return ListStructCapn(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewListStructCapn(s *C.Segment) ListStructCapn {
	return ListStructCapn(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootListStructCapn(s *C.Segment) ListStructCapn { return ListStructCapn(s.Root(0).ToStruct()) }
func (s ListStructCapn) Vec() Nester1Capn_List           { return Nester1Capn_List(C.Struct(s).GetObject(0)) }
func (s ListStructCapn) SetVec(v Nester1Capn_List)       { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s ListStructCapn) MarshalJSON() (bs []byte, err error) { return }

type ListStructCapn_List C.PointerList

func NewListStructCapn_List(s *C.Segment, sz int) ListStructCapn_List {
	return ListStructCapn_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s ListStructCapn_List) Len() int { return C.PointerList(s).Len() }
func (s ListStructCapn_List) At(i int) ListStructCapn {
	return ListStructCapn(C.PointerList(s).At(i).ToStruct())
}
func (s ListStructCapn_List) Set(i int, item ListStructCapn) { C.PointerList(s).Set(i, C.Object(item)) }

type ListStructCapn_Promise C.Pipeline

func (p *ListStructCapn_Promise) Get() (ListStructCapn, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return ListStructCapn(s), err
}

type Echo struct{ c C.Client }

func NewEcho(c C.Client) Echo { return Echo{c} }

func (c Echo) GenericClient() C.Client { return c.c }

func (c Echo) IsNull() bool { return c.c == nil }

func (c Echo) Echo(ctx context.Context, params func(Echo_echo_Params), opts ...C.CallOption) *Echo_echo_Results_Promise {
	if c.c == nil {
		return (*Echo_echo_Results_Promise)(C.NewPipeline(C.ErrorAnswer(C.ErrNullClient)))
	}
	return (*Echo_echo_Results_Promise)(C.NewPipeline(c.c.Call(&C.Call{
		Ctx: ctx,
		Method: C.Method{

			InterfaceID:   0x8e5322c1e9282534,
			MethodID:      0,
			InterfaceName: "aircraft.capnp:Echo",
			MethodName:    "echo",
		},
		ParamsSize: C.ObjectSize{DataSize: 0, PointerCount: 1},
		ParamsFunc: func(s C.Struct) { params(Echo_echo_Params(s)) },
		Options:    C.NewCallOptions(opts),
	})))
}

type Echo_Server interface {
	Echo(ctx context.Context, opts C.CallOptions, params Echo_echo_Params, results Echo_echo_Results) error
}

func Echo_ServerToClient(s Echo_Server) Echo {
	return NewEcho(C.NewServer(Echo_Methods(nil, s)))
}

func Echo_Methods(methods []C.ServerMethod, server Echo_Server) []C.ServerMethod {
	if cap(methods) == 0 {
		methods = make([]C.ServerMethod, 0, 1)
	}

	methods = append(methods, C.ServerMethod{
		Method: C.Method{

			InterfaceID:   0x8e5322c1e9282534,
			MethodID:      0,
			InterfaceName: "aircraft.capnp:Echo",
			MethodName:    "echo",
		},
		Impl: func(c context.Context, opts C.CallOptions, p, r C.Struct) error {
			return server.Echo(c, opts, Echo_echo_Params(p), Echo_echo_Results(r))
		},
		ResultsSize: C.ObjectSize{DataSize: 0, PointerCount: 1},
	})

	return methods
}

type Echo_echo_Params C.Struct

func NewEcho_echo_Params(s *C.Segment) Echo_echo_Params {
	return Echo_echo_Params(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func NewRootEcho_echo_Params(s *C.Segment) Echo_echo_Params {
	return Echo_echo_Params(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewEcho_echo_Params(s *C.Segment) Echo_echo_Params {
	return Echo_echo_Params(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootEcho_echo_Params(s *C.Segment) Echo_echo_Params {
	return Echo_echo_Params(s.Root(0).ToStruct())
}
func (s Echo_echo_Params) In() string     { return C.Struct(s).GetObject(0).ToText() }
func (s Echo_echo_Params) SetIn(v string) { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Echo_echo_Params) MarshalJSON() (bs []byte, err error) { return }

type Echo_echo_Params_List C.PointerList

func NewEcho_echo_Params_List(s *C.Segment, sz int) Echo_echo_Params_List {
	return Echo_echo_Params_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s Echo_echo_Params_List) Len() int { return C.PointerList(s).Len() }
func (s Echo_echo_Params_List) At(i int) Echo_echo_Params {
	return Echo_echo_Params(C.PointerList(s).At(i).ToStruct())
}
func (s Echo_echo_Params_List) Set(i int, item Echo_echo_Params) {
	C.PointerList(s).Set(i, C.Object(item))
}

type Echo_echo_Params_Promise C.Pipeline

func (p *Echo_echo_Params_Promise) Get() (Echo_echo_Params, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Echo_echo_Params(s), err
}

type Echo_echo_Results C.Struct

func NewEcho_echo_Results(s *C.Segment) Echo_echo_Results {
	return Echo_echo_Results(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func NewRootEcho_echo_Results(s *C.Segment) Echo_echo_Results {
	return Echo_echo_Results(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewEcho_echo_Results(s *C.Segment) Echo_echo_Results {
	return Echo_echo_Results(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootEcho_echo_Results(s *C.Segment) Echo_echo_Results {
	return Echo_echo_Results(s.Root(0).ToStruct())
}
func (s Echo_echo_Results) Out() string     { return C.Struct(s).GetObject(0).ToText() }
func (s Echo_echo_Results) SetOut(v string) { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Echo_echo_Results) MarshalJSON() (bs []byte, err error) { return }

type Echo_echo_Results_List C.PointerList

func NewEcho_echo_Results_List(s *C.Segment, sz int) Echo_echo_Results_List {
	return Echo_echo_Results_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s Echo_echo_Results_List) Len() int { return C.PointerList(s).Len() }
func (s Echo_echo_Results_List) At(i int) Echo_echo_Results {
	return Echo_echo_Results(C.PointerList(s).At(i).ToStruct())
}
func (s Echo_echo_Results_List) Set(i int, item Echo_echo_Results) {
	C.PointerList(s).Set(i, C.Object(item))
}

type Echo_echo_Results_Promise C.Pipeline

func (p *Echo_echo_Results_Promise) Get() (Echo_echo_Results, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Echo_echo_Results(s), err
}

type Hoth C.Struct

func NewHoth(s *C.Segment) Hoth { return Hoth(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1})) }
func NewRootHoth(s *C.Segment) Hoth {
	return Hoth(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewHoth(s *C.Segment) Hoth {
	return Hoth(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootHoth(s *C.Segment) Hoth { return Hoth(s.Root(0).ToStruct()) }
func (s Hoth) Base() EchoBase        { return EchoBase(C.Struct(s).GetObject(0).ToStruct()) }
func (s Hoth) SetBase(v EchoBase)    { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Hoth) MarshalJSON() (bs []byte, err error) { return }

type Hoth_List C.PointerList

func NewHoth_List(s *C.Segment, sz int) Hoth_List {
	return Hoth_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s Hoth_List) Len() int             { return C.PointerList(s).Len() }
func (s Hoth_List) At(i int) Hoth        { return Hoth(C.PointerList(s).At(i).ToStruct()) }
func (s Hoth_List) Set(i int, item Hoth) { C.PointerList(s).Set(i, C.Object(item)) }

type Hoth_Promise C.Pipeline

func (p *Hoth_Promise) Get() (Hoth, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Hoth(s), err
}

func (p *Hoth_Promise) Base() *EchoBase_Promise {
	return (*EchoBase_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

type EchoBase C.Struct

func NewEchoBase(s *C.Segment) EchoBase {
	return EchoBase(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func NewRootEchoBase(s *C.Segment) EchoBase {
	return EchoBase(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewEchoBase(s *C.Segment) EchoBase {
	return EchoBase(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootEchoBase(s *C.Segment) EchoBase { return EchoBase(s.Root(0).ToStruct()) }
func (s EchoBase) Echo() Echo                { return NewEcho(C.Struct(s).GetObject(0).ToInterface().Client()) }
func (s EchoBase) SetEcho(v Echo) {
	if s.Segment == nil {
		return
	}
	ci := s.Segment.Message.AddCap(v.GenericClient())
	C.Struct(s).SetObject(0, C.Object(s.Segment.NewInterface(ci)))
}

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s EchoBase) MarshalJSON() (bs []byte, err error) { return }

type EchoBase_List C.PointerList

func NewEchoBase_List(s *C.Segment, sz int) EchoBase_List {
	return EchoBase_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s EchoBase_List) Len() int                 { return C.PointerList(s).Len() }
func (s EchoBase_List) At(i int) EchoBase        { return EchoBase(C.PointerList(s).At(i).ToStruct()) }
func (s EchoBase_List) Set(i int, item EchoBase) { C.PointerList(s).Set(i, C.Object(item)) }

type EchoBase_Promise C.Pipeline

func (p *EchoBase_Promise) Get() (EchoBase, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return EchoBase(s), err
}

func (p *EchoBase_Promise) Echo() Echo {
	return NewEcho((*C.Pipeline)(p).GetPipeline(0).Client())
}

type StackingRoot C.Struct

func NewStackingRoot(s *C.Segment) StackingRoot {
	return StackingRoot(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 2}))
}
func NewRootStackingRoot(s *C.Segment) StackingRoot {
	return StackingRoot(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 2}))
}
func AutoNewStackingRoot(s *C.Segment) StackingRoot {
	return StackingRoot(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 2}))
}
func ReadRootStackingRoot(s *C.Segment) StackingRoot { return StackingRoot(s.Root(0).ToStruct()) }
func (s StackingRoot) A() StackingA                  { return StackingA(C.Struct(s).GetObject(1).ToStruct()) }
func (s StackingRoot) SetA(v StackingA)              { C.Struct(s).SetObject(1, C.Object(v)) }
func (s StackingRoot) AWithDefault() StackingA {
	return StackingA(C.Struct(s).GetObject(0).ToStructDefault(x_832bcc6686a26d56, 0))
}
func (s StackingRoot) SetAWithDefault(v StackingA) { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s StackingRoot) MarshalJSON() (bs []byte, err error) { return }

type StackingRoot_List C.PointerList

func NewStackingRoot_List(s *C.Segment, sz int) StackingRoot_List {
	return StackingRoot_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 2}, sz))
}
func (s StackingRoot_List) Len() int { return C.PointerList(s).Len() }
func (s StackingRoot_List) At(i int) StackingRoot {
	return StackingRoot(C.PointerList(s).At(i).ToStruct())
}
func (s StackingRoot_List) Set(i int, item StackingRoot) { C.PointerList(s).Set(i, C.Object(item)) }

type StackingRoot_Promise C.Pipeline

func (p *StackingRoot_Promise) Get() (StackingRoot, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return StackingRoot(s), err
}

func (p *StackingRoot_Promise) A() *StackingA_Promise {
	return (*StackingA_Promise)((*C.Pipeline)(p).GetPipeline(1))
}

func (p *StackingRoot_Promise) AWithDefault() *StackingA_Promise {
	return (*StackingA_Promise)((*C.Pipeline)(p).GetPipelineDefault(0, x_832bcc6686a26d56, 3))
}

type StackingA C.Struct

func NewStackingA(s *C.Segment) StackingA {
	return StackingA(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func NewRootStackingA(s *C.Segment) StackingA {
	return StackingA(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func AutoNewStackingA(s *C.Segment) StackingA {
	return StackingA(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func ReadRootStackingA(s *C.Segment) StackingA { return StackingA(s.Root(0).ToStruct()) }
func (s StackingA) Num() int32                 { return int32(C.Struct(s).Get32(0)) }
func (s StackingA) SetNum(v int32)             { C.Struct(s).Set32(0, uint32(v)) }
func (s StackingA) B() StackingB               { return StackingB(C.Struct(s).GetObject(0).ToStruct()) }
func (s StackingA) SetB(v StackingB)           { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s StackingA) MarshalJSON() (bs []byte, err error) { return }

type StackingA_List C.PointerList

func NewStackingA_List(s *C.Segment, sz int) StackingA_List {
	return StackingA_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 1}, sz))
}
func (s StackingA_List) Len() int                  { return C.PointerList(s).Len() }
func (s StackingA_List) At(i int) StackingA        { return StackingA(C.PointerList(s).At(i).ToStruct()) }
func (s StackingA_List) Set(i int, item StackingA) { C.PointerList(s).Set(i, C.Object(item)) }

type StackingA_Promise C.Pipeline

func (p *StackingA_Promise) Get() (StackingA, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return StackingA(s), err
}

func (p *StackingA_Promise) B() *StackingB_Promise {
	return (*StackingB_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

type StackingB C.Struct

func NewStackingB(s *C.Segment) StackingB {
	return StackingB(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 0}))
}
func NewRootStackingB(s *C.Segment) StackingB {
	return StackingB(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 0}))
}
func AutoNewStackingB(s *C.Segment) StackingB {
	return StackingB(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 0}))
}
func ReadRootStackingB(s *C.Segment) StackingB { return StackingB(s.Root(0).ToStruct()) }
func (s StackingB) Num() int32                 { return int32(C.Struct(s).Get32(0)) }
func (s StackingB) SetNum(v int32)             { C.Struct(s).Set32(0, uint32(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s StackingB) MarshalJSON() (bs []byte, err error) { return }

type StackingB_List C.PointerList

func NewStackingB_List(s *C.Segment, sz int) StackingB_List {
	return StackingB_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 0}, sz))
}
func (s StackingB_List) Len() int                  { return C.PointerList(s).Len() }
func (s StackingB_List) At(i int) StackingB        { return StackingB(C.PointerList(s).At(i).ToStruct()) }
func (s StackingB_List) Set(i int, item StackingB) { C.PointerList(s).Set(i, C.Object(item)) }

type StackingB_Promise C.Pipeline

func (p *StackingB_Promise) Get() (StackingB, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return StackingB(s), err
}

var x_832bcc6686a26d56 = C.NewBuffer([]byte{
	0, 0, 0, 0, 1, 0, 1, 0,
	42, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 1, 0, 1, 0,
	42, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
})
