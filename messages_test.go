package capn_test

// AUTO GENERATED - DO NOT EDIT

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"math"
	"unsafe"

	C "github.com/glycerine/go-capnproto"
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
func (s Zdate) WriteJSON(w io.Writer) error {
	b := bufio.NewWriter(w)
	var err error
	var buf []byte
	_ = buf
	err = b.WriteByte('{')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"year\":")
	if err != nil {
		return err
	}
	{
		s := s.Year()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(',')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"month\":")
	if err != nil {
		return err
	}
	{
		s := s.Month()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte(',')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"day\":")
	if err != nil {
		return err
	}
	{
		s := s.Day()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte('}')
	if err != nil {
		return err
	}
	err = b.Flush()
	return err
}
func (s Zdate) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteJSON(&b)
	return b.Bytes(), err
}

type Zdate_List C.PointerList

func NewZdateList(s *C.Segment, sz int) Zdate_List { return Zdate_List(s.NewUInt32List(sz)) }
func (s Zdate_List) Len() int                      { return C.PointerList(s).Len() }
func (s Zdate_List) At(i int) Zdate                { return Zdate(C.PointerList(s).At(i).ToStruct()) }
func (s Zdate_List) Set(i int, src Zdate)          { C.PointerList(s).Set(i, C.Object(src)) }
func (s Zdate_List) ToArray() []Zdate              { return *(*[]Zdate)(unsafe.Pointer(C.PointerList(s).ToArray())) }

type Zdata C.Struct

func NewZdata(s *C.Segment) Zdata      { return Zdata(s.NewStruct(0, 1)) }
func NewRootZdata(s *C.Segment) Zdata  { return Zdata(s.NewRootStruct(0, 1)) }
func ReadRootZdata(s *C.Segment) Zdata { return Zdata(s.Root(0).ToStruct()) }
func (s Zdata) Data() []byte           { return C.Struct(s).GetObject(0).ToData() }
func (s Zdata) SetData(v []byte)       { C.Struct(s).SetObject(0, s.Segment.NewData(v)) }
func (s Zdata) WriteJSON(w io.Writer) error {
	b := bufio.NewWriter(w)
	var err error
	var buf []byte
	_ = buf
	err = b.WriteByte('{')
	if err != nil {
		return err
	}
	_, err = b.WriteString("\"data\":")
	if err != nil {
		return err
	}
	{
		s := s.Data()
		buf, err = json.Marshal(s)
		if err != nil {
			return err
		}
		_, err = b.Write(buf)
		if err != nil {
			return err
		}
	}
	err = b.WriteByte('}')
	if err != nil {
		return err
	}
	err = b.Flush()
	return err
}
func (s Zdata) MarshalJSON() ([]byte, error) {
	b := bytes.Buffer{}
	err := s.WriteJSON(&b)
	return b.Bytes(), err
}

type Zdata_List C.PointerList

func NewZdataList(s *C.Segment, sz int) Zdata_List { return Zdata_List(s.NewCompositeList(0, 1, sz)) }
func (s Zdata_List) Len() int                      { return C.PointerList(s).Len() }
func (s Zdata_List) At(i int) Zdata                { return Zdata(C.PointerList(s).At(i).ToStruct()) }
func (s Zdata_List) ToArray() []Zdata              { return *(*[]Zdata)(unsafe.Pointer(C.PointerList(s).ToArray())) }

type Airport uint16

const (
	AIRPORT_NONE Airport = 0
	AIRPORT_JFK          = 1
	AIRPORT_LAX          = 2
	AIRPORT_SFO          = 3
	AIRPORT_LUV          = 4
	AIRPORT_DFW          = 5
	AIRPORT_TEST         = 6
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
	AIRCRAFT_B737                = 1
	AIRCRAFT_A320                = 2
	AIRCRAFT_F16                 = 3
)

func NewAircraft(s *C.Segment) Aircraft      { return Aircraft(s.NewStruct(8, 1)) }
func NewRootAircraft(s *C.Segment) Aircraft  { return Aircraft(s.NewRootStruct(8, 1)) }
func ReadRootAircraft(s *C.Segment) Aircraft { return Aircraft(s.Root(0).ToStruct()) }
func (s Aircraft) Which() Aircraft_Which     { return Aircraft_Which(C.Struct(s).Get16(0)) }
func (s Aircraft) B737() B737                { return B737(C.Struct(s).GetObject(0).ToStruct()) }
func (s Aircraft) SetB737(v B737)            { C.Struct(s).Set16(0, 1); C.Struct(s).SetObject(0, C.Object(v)) }
func (s Aircraft) A320() A320                { return A320(C.Struct(s).GetObject(0).ToStruct()) }
func (s Aircraft) SetA320(v A320)            { C.Struct(s).Set16(0, 2); C.Struct(s).SetObject(0, C.Object(v)) }
func (s Aircraft) F16() F16                  { return F16(C.Struct(s).GetObject(0).ToStruct()) }
func (s Aircraft) SetF16(v F16)              { C.Struct(s).Set16(0, 3); C.Struct(s).SetObject(0, C.Object(v)) }

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
	Z_ZZ                  = 1
	Z_F64                 = 2
	Z_F32                 = 3
	Z_I64                 = 4
	Z_I32                 = 5
	Z_I16                 = 6
	Z_I8                  = 7
	Z_U64                 = 8
	Z_U32                 = 9
	Z_U16                 = 10
	Z_U8                  = 11
	Z_BOOL                = 12
	Z_TEXT                = 13
	Z_BLOB                = 14
	Z_F64VEC              = 15
	Z_F32VEC              = 16
	Z_I64VEC              = 17
	Z_I32VEC              = 18
	Z_I16VEC              = 19
	Z_I8VEC               = 20
	Z_U64VEC              = 21
	Z_U32VEC              = 22
	Z_U16VEC              = 23
	Z_U8VEC               = 24
	Z_ZVEC                = 25
	Z_ZVECVEC             = 26
	Z_ZDATE               = 27
	Z_ZDATA               = 28
	Z_AIRCRAFTVEC         = 29
	Z_AIRCRAFT            = 30
	Z_REGRESSION          = 31
	Z_PLANEBASE           = 32
	Z_AIRPORT             = 33
	Z_B737                = 34
	Z_A320                = 35
	Z_F16                 = 36
	Z_ZDATEVEC            = 37
)

func NewZ(s *C.Segment) Z             { return Z(s.NewStruct(16, 1)) }
func NewRootZ(s *C.Segment) Z         { return Z(s.NewRootStruct(16, 1)) }
func ReadRootZ(s *C.Segment) Z        { return Z(s.Root(0).ToStruct()) }
func (s Z) Which() Z_Which            { return Z_Which(C.Struct(s).Get16(0)) }
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

type Z_List C.PointerList

func NewZList(s *C.Segment, sz int) Z_List { return Z_List(s.NewCompositeList(16, 1, sz)) }
func (s Z_List) Len() int                  { return C.PointerList(s).Len() }
func (s Z_List) At(i int) Z                { return Z(C.PointerList(s).At(i).ToStruct()) }
func (s Z_List) ToArray() []Z              { return *(*[]Z)(unsafe.Pointer(C.PointerList(s).ToArray())) }
