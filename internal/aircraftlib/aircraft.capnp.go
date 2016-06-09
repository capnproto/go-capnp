package aircraftlib

// AUTO GENERATED - DO NOT EDIT

import (
	context "golang.org/x/net/context"
	math "math"
	strconv "strconv"
	capnp "zombiezen.com/go/capnproto2"
	server "zombiezen.com/go/capnproto2/server"
)

// Constants defined in aircraft.capnp.
const (
	ConstEnum = Airport_jfk
)

// Constants defined in aircraft.capnp.
var (
	ConstDate = Zdate{Struct: capnp.MustUnmarshalRootPtr(x_832bcc6686a26d56[0:24]).Struct()}
	ConstList = Zdate_List{List: capnp.MustUnmarshalRootPtr(x_832bcc6686a26d56[24:64]).List()}
)

type Zdate struct{ capnp.Struct }

func NewZdate(s *capnp.Segment) (Zdate, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	if err != nil {
		return Zdate{}, err
	}
	return Zdate{st}, nil
}

func NewRootZdate(s *capnp.Segment) (Zdate, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	if err != nil {
		return Zdate{}, err
	}
	return Zdate{st}, nil
}

func ReadRootZdate(msg *capnp.Message) (Zdate, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Zdate{}, err
	}
	return Zdate{root.Struct()}, nil
}
func (s Zdate) Year() int16 {
	return int16(s.Struct.Uint16(0))
}

func (s Zdate) SetYear(v int16) {
	s.Struct.SetUint16(0, uint16(v))
}

func (s Zdate) Month() uint8 {
	return s.Struct.Uint8(2)
}

func (s Zdate) SetMonth(v uint8) {
	s.Struct.SetUint8(2, v)
}

func (s Zdate) Day() uint8 {
	return s.Struct.Uint8(3)
}

func (s Zdate) SetDay(v uint8) {
	s.Struct.SetUint8(3, v)
}

// Zdate_List is a list of Zdate.
type Zdate_List struct{ capnp.List }

// NewZdate creates a new list of Zdate.
func NewZdate_List(s *capnp.Segment, sz int32) (Zdate_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0}, sz)
	if err != nil {
		return Zdate_List{}, err
	}
	return Zdate_List{l}, nil
}

func (s Zdate_List) At(i int) Zdate           { return Zdate{s.List.Struct(i)} }
func (s Zdate_List) Set(i int, v Zdate) error { return s.List.SetStruct(i, v.Struct) }

// Zdate_Promise is a wrapper for a Zdate promised by a client call.
type Zdate_Promise struct{ *capnp.Pipeline }

func (p Zdate_Promise) Struct() (Zdate, error) {
	s, err := p.Pipeline.Struct()
	return Zdate{s}, err
}

type Zdata struct{ capnp.Struct }

func NewZdata(s *capnp.Segment) (Zdata, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Zdata{}, err
	}
	return Zdata{st}, nil
}

func NewRootZdata(s *capnp.Segment) (Zdata, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Zdata{}, err
	}
	return Zdata{st}, nil
}

func ReadRootZdata(msg *capnp.Message) (Zdata, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Zdata{}, err
	}
	return Zdata{root.Struct()}, nil
}
func (s Zdata) Data() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	return []byte(p.Data()), nil
}

func (s Zdata) HasData() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Zdata) SetData(v []byte) error {
	d, err := capnp.NewData(s.Struct.Segment(), []byte(v))
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, d.List.ToPtr())
}

// Zdata_List is a list of Zdata.
type Zdata_List struct{ capnp.List }

// NewZdata creates a new list of Zdata.
func NewZdata_List(s *capnp.Segment, sz int32) (Zdata_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return Zdata_List{}, err
	}
	return Zdata_List{l}, nil
}

func (s Zdata_List) At(i int) Zdata           { return Zdata{s.List.Struct(i)} }
func (s Zdata_List) Set(i int, v Zdata) error { return s.List.SetStruct(i, v.Struct) }

// Zdata_Promise is a wrapper for a Zdata promised by a client call.
type Zdata_Promise struct{ *capnp.Pipeline }

func (p Zdata_Promise) Struct() (Zdata, error) {
	s, err := p.Pipeline.Struct()
	return Zdata{s}, err
}

type Airport uint16

// Values of Airport.
const (
	Airport_none Airport = 0
	Airport_jfk  Airport = 1
	Airport_lax  Airport = 2
	Airport_sfo  Airport = 3
	Airport_luv  Airport = 4
	Airport_dfw  Airport = 5
	Airport_test Airport = 6
)

// String returns the enum's constant name.
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

// AirportFromString returns the enum value with a name,
// or the zero value if there's no such value.
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

type Airport_List struct{ capnp.List }

func NewAirport_List(s *capnp.Segment, sz int32) (Airport_List, error) {
	l, err := capnp.NewUInt16List(s, sz)
	if err != nil {
		return Airport_List{}, err
	}
	return Airport_List{l.List}, nil
}

func (l Airport_List) At(i int) Airport {
	ul := capnp.UInt16List{List: l.List}
	return Airport(ul.At(i))
}

func (l Airport_List) Set(i int, v Airport) {
	ul := capnp.UInt16List{List: l.List}
	ul.Set(i, uint16(v))
}

type PlaneBase struct{ capnp.Struct }

func NewPlaneBase(s *capnp.Segment) (PlaneBase, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 32, PointerCount: 2})
	if err != nil {
		return PlaneBase{}, err
	}
	return PlaneBase{st}, nil
}

func NewRootPlaneBase(s *capnp.Segment) (PlaneBase, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 32, PointerCount: 2})
	if err != nil {
		return PlaneBase{}, err
	}
	return PlaneBase{st}, nil
}

func ReadRootPlaneBase(msg *capnp.Message) (PlaneBase, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return PlaneBase{}, err
	}
	return PlaneBase{root.Struct()}, nil
}
func (s PlaneBase) Name() (string, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return "", err
	}
	return p.Text(), nil
}

func (s PlaneBase) HasName() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s PlaneBase) NameBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	d := p.Data()
	if len(d) == 0 {
		return d, nil
	}
	return d[:len(d)-1], nil
}

func (s PlaneBase) SetName(v string) error {
	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, t.List.ToPtr())
}

func (s PlaneBase) Homes() (Airport_List, error) {
	p, err := s.Struct.Ptr(1)
	if err != nil {
		return Airport_List{}, err
	}
	return Airport_List{List: p.List()}, nil
}

func (s PlaneBase) HasHomes() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s PlaneBase) SetHomes(v Airport_List) error {
	return s.Struct.SetPtr(1, v.List.ToPtr())
}

// NewHomes sets the homes field to a newly
// allocated Airport_List, preferring placement in s's segment.
func (s PlaneBase) NewHomes(n int32) (Airport_List, error) {
	l, err := NewAirport_List(s.Struct.Segment(), n)
	if err != nil {
		return Airport_List{}, err
	}
	err = s.Struct.SetPtr(1, l.List.ToPtr())
	return l, err
}

func (s PlaneBase) Rating() int64 {
	return int64(s.Struct.Uint64(0))
}

func (s PlaneBase) SetRating(v int64) {
	s.Struct.SetUint64(0, uint64(v))
}

func (s PlaneBase) CanFly() bool {
	return s.Struct.Bit(64)
}

func (s PlaneBase) SetCanFly(v bool) {
	s.Struct.SetBit(64, v)
}

func (s PlaneBase) Capacity() int64 {
	return int64(s.Struct.Uint64(16))
}

func (s PlaneBase) SetCapacity(v int64) {
	s.Struct.SetUint64(16, uint64(v))
}

func (s PlaneBase) MaxSpeed() float64 {
	return math.Float64frombits(s.Struct.Uint64(24))
}

func (s PlaneBase) SetMaxSpeed(v float64) {
	s.Struct.SetUint64(24, math.Float64bits(v))
}

// PlaneBase_List is a list of PlaneBase.
type PlaneBase_List struct{ capnp.List }

// NewPlaneBase creates a new list of PlaneBase.
func NewPlaneBase_List(s *capnp.Segment, sz int32) (PlaneBase_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 32, PointerCount: 2}, sz)
	if err != nil {
		return PlaneBase_List{}, err
	}
	return PlaneBase_List{l}, nil
}

func (s PlaneBase_List) At(i int) PlaneBase           { return PlaneBase{s.List.Struct(i)} }
func (s PlaneBase_List) Set(i int, v PlaneBase) error { return s.List.SetStruct(i, v.Struct) }

// PlaneBase_Promise is a wrapper for a PlaneBase promised by a client call.
type PlaneBase_Promise struct{ *capnp.Pipeline }

func (p PlaneBase_Promise) Struct() (PlaneBase, error) {
	s, err := p.Pipeline.Struct()
	return PlaneBase{s}, err
}

type B737 struct{ capnp.Struct }

func NewB737(s *capnp.Segment) (B737, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return B737{}, err
	}
	return B737{st}, nil
}

func NewRootB737(s *capnp.Segment) (B737, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return B737{}, err
	}
	return B737{st}, nil
}

func ReadRootB737(msg *capnp.Message) (B737, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return B737{}, err
	}
	return B737{root.Struct()}, nil
}
func (s B737) Base() (PlaneBase, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return PlaneBase{}, err
	}
	return PlaneBase{Struct: p.Struct()}, nil
}

func (s B737) HasBase() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s B737) SetBase(v PlaneBase) error {
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewBase sets the base field to a newly
// allocated PlaneBase struct, preferring placement in s's segment.
func (s B737) NewBase() (PlaneBase, error) {
	ss, err := NewPlaneBase(s.Struct.Segment())
	if err != nil {
		return PlaneBase{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

// B737_List is a list of B737.
type B737_List struct{ capnp.List }

// NewB737 creates a new list of B737.
func NewB737_List(s *capnp.Segment, sz int32) (B737_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return B737_List{}, err
	}
	return B737_List{l}, nil
}

func (s B737_List) At(i int) B737           { return B737{s.List.Struct(i)} }
func (s B737_List) Set(i int, v B737) error { return s.List.SetStruct(i, v.Struct) }

// B737_Promise is a wrapper for a B737 promised by a client call.
type B737_Promise struct{ *capnp.Pipeline }

func (p B737_Promise) Struct() (B737, error) {
	s, err := p.Pipeline.Struct()
	return B737{s}, err
}

func (p B737_Promise) Base() PlaneBase_Promise {
	return PlaneBase_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

type A320 struct{ capnp.Struct }

func NewA320(s *capnp.Segment) (A320, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return A320{}, err
	}
	return A320{st}, nil
}

func NewRootA320(s *capnp.Segment) (A320, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return A320{}, err
	}
	return A320{st}, nil
}

func ReadRootA320(msg *capnp.Message) (A320, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return A320{}, err
	}
	return A320{root.Struct()}, nil
}
func (s A320) Base() (PlaneBase, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return PlaneBase{}, err
	}
	return PlaneBase{Struct: p.Struct()}, nil
}

func (s A320) HasBase() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s A320) SetBase(v PlaneBase) error {
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewBase sets the base field to a newly
// allocated PlaneBase struct, preferring placement in s's segment.
func (s A320) NewBase() (PlaneBase, error) {
	ss, err := NewPlaneBase(s.Struct.Segment())
	if err != nil {
		return PlaneBase{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

// A320_List is a list of A320.
type A320_List struct{ capnp.List }

// NewA320 creates a new list of A320.
func NewA320_List(s *capnp.Segment, sz int32) (A320_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return A320_List{}, err
	}
	return A320_List{l}, nil
}

func (s A320_List) At(i int) A320           { return A320{s.List.Struct(i)} }
func (s A320_List) Set(i int, v A320) error { return s.List.SetStruct(i, v.Struct) }

// A320_Promise is a wrapper for a A320 promised by a client call.
type A320_Promise struct{ *capnp.Pipeline }

func (p A320_Promise) Struct() (A320, error) {
	s, err := p.Pipeline.Struct()
	return A320{s}, err
}

func (p A320_Promise) Base() PlaneBase_Promise {
	return PlaneBase_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

type F16 struct{ capnp.Struct }

func NewF16(s *capnp.Segment) (F16, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return F16{}, err
	}
	return F16{st}, nil
}

func NewRootF16(s *capnp.Segment) (F16, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return F16{}, err
	}
	return F16{st}, nil
}

func ReadRootF16(msg *capnp.Message) (F16, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return F16{}, err
	}
	return F16{root.Struct()}, nil
}
func (s F16) Base() (PlaneBase, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return PlaneBase{}, err
	}
	return PlaneBase{Struct: p.Struct()}, nil
}

func (s F16) HasBase() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s F16) SetBase(v PlaneBase) error {
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewBase sets the base field to a newly
// allocated PlaneBase struct, preferring placement in s's segment.
func (s F16) NewBase() (PlaneBase, error) {
	ss, err := NewPlaneBase(s.Struct.Segment())
	if err != nil {
		return PlaneBase{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

// F16_List is a list of F16.
type F16_List struct{ capnp.List }

// NewF16 creates a new list of F16.
func NewF16_List(s *capnp.Segment, sz int32) (F16_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return F16_List{}, err
	}
	return F16_List{l}, nil
}

func (s F16_List) At(i int) F16           { return F16{s.List.Struct(i)} }
func (s F16_List) Set(i int, v F16) error { return s.List.SetStruct(i, v.Struct) }

// F16_Promise is a wrapper for a F16 promised by a client call.
type F16_Promise struct{ *capnp.Pipeline }

func (p F16_Promise) Struct() (F16, error) {
	s, err := p.Pipeline.Struct()
	return F16{s}, err
}

func (p F16_Promise) Base() PlaneBase_Promise {
	return PlaneBase_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

type Regression struct{ capnp.Struct }

func NewRegression(s *capnp.Segment) (Regression, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 24, PointerCount: 3})
	if err != nil {
		return Regression{}, err
	}
	return Regression{st}, nil
}

func NewRootRegression(s *capnp.Segment) (Regression, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 24, PointerCount: 3})
	if err != nil {
		return Regression{}, err
	}
	return Regression{st}, nil
}

func ReadRootRegression(msg *capnp.Message) (Regression, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Regression{}, err
	}
	return Regression{root.Struct()}, nil
}
func (s Regression) Base() (PlaneBase, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return PlaneBase{}, err
	}
	return PlaneBase{Struct: p.Struct()}, nil
}

func (s Regression) HasBase() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Regression) SetBase(v PlaneBase) error {
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewBase sets the base field to a newly
// allocated PlaneBase struct, preferring placement in s's segment.
func (s Regression) NewBase() (PlaneBase, error) {
	ss, err := NewPlaneBase(s.Struct.Segment())
	if err != nil {
		return PlaneBase{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Regression) B0() float64 {
	return math.Float64frombits(s.Struct.Uint64(0))
}

func (s Regression) SetB0(v float64) {
	s.Struct.SetUint64(0, math.Float64bits(v))
}

func (s Regression) Beta() (capnp.Float64List, error) {
	p, err := s.Struct.Ptr(1)
	if err != nil {
		return capnp.Float64List{}, err
	}
	return capnp.Float64List{List: p.List()}, nil
}

func (s Regression) HasBeta() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s Regression) SetBeta(v capnp.Float64List) error {
	return s.Struct.SetPtr(1, v.List.ToPtr())
}

// NewBeta sets the beta field to a newly
// allocated capnp.Float64List, preferring placement in s's segment.
func (s Regression) NewBeta(n int32) (capnp.Float64List, error) {
	l, err := capnp.NewFloat64List(s.Struct.Segment(), n)
	if err != nil {
		return capnp.Float64List{}, err
	}
	err = s.Struct.SetPtr(1, l.List.ToPtr())
	return l, err
}

func (s Regression) Planes() (Aircraft_List, error) {
	p, err := s.Struct.Ptr(2)
	if err != nil {
		return Aircraft_List{}, err
	}
	return Aircraft_List{List: p.List()}, nil
}

func (s Regression) HasPlanes() bool {
	p, err := s.Struct.Ptr(2)
	return p.IsValid() || err != nil
}

func (s Regression) SetPlanes(v Aircraft_List) error {
	return s.Struct.SetPtr(2, v.List.ToPtr())
}

// NewPlanes sets the planes field to a newly
// allocated Aircraft_List, preferring placement in s's segment.
func (s Regression) NewPlanes(n int32) (Aircraft_List, error) {
	l, err := NewAircraft_List(s.Struct.Segment(), n)
	if err != nil {
		return Aircraft_List{}, err
	}
	err = s.Struct.SetPtr(2, l.List.ToPtr())
	return l, err
}

func (s Regression) Ymu() float64 {
	return math.Float64frombits(s.Struct.Uint64(8))
}

func (s Regression) SetYmu(v float64) {
	s.Struct.SetUint64(8, math.Float64bits(v))
}

func (s Regression) Ysd() float64 {
	return math.Float64frombits(s.Struct.Uint64(16))
}

func (s Regression) SetYsd(v float64) {
	s.Struct.SetUint64(16, math.Float64bits(v))
}

// Regression_List is a list of Regression.
type Regression_List struct{ capnp.List }

// NewRegression creates a new list of Regression.
func NewRegression_List(s *capnp.Segment, sz int32) (Regression_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 24, PointerCount: 3}, sz)
	if err != nil {
		return Regression_List{}, err
	}
	return Regression_List{l}, nil
}

func (s Regression_List) At(i int) Regression           { return Regression{s.List.Struct(i)} }
func (s Regression_List) Set(i int, v Regression) error { return s.List.SetStruct(i, v.Struct) }

// Regression_Promise is a wrapper for a Regression promised by a client call.
type Regression_Promise struct{ *capnp.Pipeline }

func (p Regression_Promise) Struct() (Regression, error) {
	s, err := p.Pipeline.Struct()
	return Regression{s}, err
}

func (p Regression_Promise) Base() PlaneBase_Promise {
	return PlaneBase_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

type Aircraft struct{ capnp.Struct }
type Aircraft_Which uint16

const (
	Aircraft_Which_void Aircraft_Which = 0
	Aircraft_Which_b737 Aircraft_Which = 1
	Aircraft_Which_a320 Aircraft_Which = 2
	Aircraft_Which_f16  Aircraft_Which = 3
)

func (w Aircraft_Which) String() string {
	const s = "voidb737a320f16"
	switch w {
	case Aircraft_Which_void:
		return s[0:4]
	case Aircraft_Which_b737:
		return s[4:8]
	case Aircraft_Which_a320:
		return s[8:12]
	case Aircraft_Which_f16:
		return s[12:15]

	}
	return "Aircraft_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
}

func NewAircraft(s *capnp.Segment) (Aircraft, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		return Aircraft{}, err
	}
	return Aircraft{st}, nil
}

func NewRootAircraft(s *capnp.Segment) (Aircraft, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		return Aircraft{}, err
	}
	return Aircraft{st}, nil
}

func ReadRootAircraft(msg *capnp.Message) (Aircraft, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Aircraft{}, err
	}
	return Aircraft{root.Struct()}, nil
}

func (s Aircraft) Which() Aircraft_Which {
	return Aircraft_Which(s.Struct.Uint16(0))
}
func (s Aircraft) SetVoid() {
	s.Struct.SetUint16(0, 0)

}

func (s Aircraft) B737() (B737, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return B737{}, err
	}
	return B737{Struct: p.Struct()}, nil
}

func (s Aircraft) HasB737() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Aircraft) SetB737(v B737) error {
	s.Struct.SetUint16(0, 1)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewB737 sets the b737 field to a newly
// allocated B737 struct, preferring placement in s's segment.
func (s Aircraft) NewB737() (B737, error) {
	s.Struct.SetUint16(0, 1)
	ss, err := NewB737(s.Struct.Segment())
	if err != nil {
		return B737{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Aircraft) A320() (A320, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return A320{}, err
	}
	return A320{Struct: p.Struct()}, nil
}

func (s Aircraft) HasA320() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Aircraft) SetA320(v A320) error {
	s.Struct.SetUint16(0, 2)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewA320 sets the a320 field to a newly
// allocated A320 struct, preferring placement in s's segment.
func (s Aircraft) NewA320() (A320, error) {
	s.Struct.SetUint16(0, 2)
	ss, err := NewA320(s.Struct.Segment())
	if err != nil {
		return A320{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Aircraft) F16() (F16, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return F16{}, err
	}
	return F16{Struct: p.Struct()}, nil
}

func (s Aircraft) HasF16() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Aircraft) SetF16(v F16) error {
	s.Struct.SetUint16(0, 3)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewF16 sets the f16 field to a newly
// allocated F16 struct, preferring placement in s's segment.
func (s Aircraft) NewF16() (F16, error) {
	s.Struct.SetUint16(0, 3)
	ss, err := NewF16(s.Struct.Segment())
	if err != nil {
		return F16{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

// Aircraft_List is a list of Aircraft.
type Aircraft_List struct{ capnp.List }

// NewAircraft creates a new list of Aircraft.
func NewAircraft_List(s *capnp.Segment, sz int32) (Aircraft_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1}, sz)
	if err != nil {
		return Aircraft_List{}, err
	}
	return Aircraft_List{l}, nil
}

func (s Aircraft_List) At(i int) Aircraft           { return Aircraft{s.List.Struct(i)} }
func (s Aircraft_List) Set(i int, v Aircraft) error { return s.List.SetStruct(i, v.Struct) }

// Aircraft_Promise is a wrapper for a Aircraft promised by a client call.
type Aircraft_Promise struct{ *capnp.Pipeline }

func (p Aircraft_Promise) Struct() (Aircraft, error) {
	s, err := p.Pipeline.Struct()
	return Aircraft{s}, err
}

func (p Aircraft_Promise) B737() B737_Promise {
	return B737_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Aircraft_Promise) A320() A320_Promise {
	return A320_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Aircraft_Promise) F16() F16_Promise {
	return F16_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

type Z struct{ capnp.Struct }
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

func (w Z_Which) String() string {
	const s = "voidzzf64f32i64i32i16i8u64u32u16u8booltextblobf64vecf32veci64veci32veci16veci8vecu64vecu32vecu16vecu8veczveczvecveczdatezdataaircraftvecaircraftregressionplanebaseairportb737a320f16zdateveczdatavecboolvec"
	switch w {
	case Z_Which_void:
		return s[0:4]
	case Z_Which_zz:
		return s[4:6]
	case Z_Which_f64:
		return s[6:9]
	case Z_Which_f32:
		return s[9:12]
	case Z_Which_i64:
		return s[12:15]
	case Z_Which_i32:
		return s[15:18]
	case Z_Which_i16:
		return s[18:21]
	case Z_Which_i8:
		return s[21:23]
	case Z_Which_u64:
		return s[23:26]
	case Z_Which_u32:
		return s[26:29]
	case Z_Which_u16:
		return s[29:32]
	case Z_Which_u8:
		return s[32:34]
	case Z_Which_bool:
		return s[34:38]
	case Z_Which_text:
		return s[38:42]
	case Z_Which_blob:
		return s[42:46]
	case Z_Which_f64vec:
		return s[46:52]
	case Z_Which_f32vec:
		return s[52:58]
	case Z_Which_i64vec:
		return s[58:64]
	case Z_Which_i32vec:
		return s[64:70]
	case Z_Which_i16vec:
		return s[70:76]
	case Z_Which_i8vec:
		return s[76:81]
	case Z_Which_u64vec:
		return s[81:87]
	case Z_Which_u32vec:
		return s[87:93]
	case Z_Which_u16vec:
		return s[93:99]
	case Z_Which_u8vec:
		return s[99:104]
	case Z_Which_zvec:
		return s[104:108]
	case Z_Which_zvecvec:
		return s[108:115]
	case Z_Which_zdate:
		return s[115:120]
	case Z_Which_zdata:
		return s[120:125]
	case Z_Which_aircraftvec:
		return s[125:136]
	case Z_Which_aircraft:
		return s[136:144]
	case Z_Which_regression:
		return s[144:154]
	case Z_Which_planebase:
		return s[154:163]
	case Z_Which_airport:
		return s[163:170]
	case Z_Which_b737:
		return s[170:174]
	case Z_Which_a320:
		return s[174:178]
	case Z_Which_f16:
		return s[178:181]
	case Z_Which_zdatevec:
		return s[181:189]
	case Z_Which_zdatavec:
		return s[189:197]
	case Z_Which_boolvec:
		return s[197:204]

	}
	return "Z_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
}

func NewZ(s *capnp.Segment) (Z, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 1})
	if err != nil {
		return Z{}, err
	}
	return Z{st}, nil
}

func NewRootZ(s *capnp.Segment) (Z, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 1})
	if err != nil {
		return Z{}, err
	}
	return Z{st}, nil
}

func ReadRootZ(msg *capnp.Message) (Z, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Z{}, err
	}
	return Z{root.Struct()}, nil
}

func (s Z) Which() Z_Which {
	return Z_Which(s.Struct.Uint16(0))
}
func (s Z) SetVoid() {
	s.Struct.SetUint16(0, 0)

}

func (s Z) Zz() (Z, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return Z{}, err
	}
	return Z{Struct: p.Struct()}, nil
}

func (s Z) HasZz() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Z) SetZz(v Z) error {
	s.Struct.SetUint16(0, 1)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewZz sets the zz field to a newly
// allocated Z struct, preferring placement in s's segment.
func (s Z) NewZz() (Z, error) {
	s.Struct.SetUint16(0, 1)
	ss, err := NewZ(s.Struct.Segment())
	if err != nil {
		return Z{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Z) F64() float64 {
	return math.Float64frombits(s.Struct.Uint64(8))
}

func (s Z) SetF64(v float64) {
	s.Struct.SetUint16(0, 2)
	s.Struct.SetUint64(8, math.Float64bits(v))
}

func (s Z) F32() float32 {
	return math.Float32frombits(s.Struct.Uint32(8))
}

func (s Z) SetF32(v float32) {
	s.Struct.SetUint16(0, 3)
	s.Struct.SetUint32(8, math.Float32bits(v))
}

func (s Z) I64() int64 {
	return int64(s.Struct.Uint64(8))
}

func (s Z) SetI64(v int64) {
	s.Struct.SetUint16(0, 4)
	s.Struct.SetUint64(8, uint64(v))
}

func (s Z) I32() int32 {
	return int32(s.Struct.Uint32(8))
}

func (s Z) SetI32(v int32) {
	s.Struct.SetUint16(0, 5)
	s.Struct.SetUint32(8, uint32(v))
}

func (s Z) I16() int16 {
	return int16(s.Struct.Uint16(8))
}

func (s Z) SetI16(v int16) {
	s.Struct.SetUint16(0, 6)
	s.Struct.SetUint16(8, uint16(v))
}

func (s Z) I8() int8 {
	return int8(s.Struct.Uint8(8))
}

func (s Z) SetI8(v int8) {
	s.Struct.SetUint16(0, 7)
	s.Struct.SetUint8(8, uint8(v))
}

func (s Z) U64() uint64 {
	return s.Struct.Uint64(8)
}

func (s Z) SetU64(v uint64) {
	s.Struct.SetUint16(0, 8)
	s.Struct.SetUint64(8, v)
}

func (s Z) U32() uint32 {
	return s.Struct.Uint32(8)
}

func (s Z) SetU32(v uint32) {
	s.Struct.SetUint16(0, 9)
	s.Struct.SetUint32(8, v)
}

func (s Z) U16() uint16 {
	return s.Struct.Uint16(8)
}

func (s Z) SetU16(v uint16) {
	s.Struct.SetUint16(0, 10)
	s.Struct.SetUint16(8, v)
}

func (s Z) U8() uint8 {
	return s.Struct.Uint8(8)
}

func (s Z) SetU8(v uint8) {
	s.Struct.SetUint16(0, 11)
	s.Struct.SetUint8(8, v)
}

func (s Z) Bool() bool {
	return s.Struct.Bit(64)
}

func (s Z) SetBool(v bool) {
	s.Struct.SetUint16(0, 12)
	s.Struct.SetBit(64, v)
}

func (s Z) Text() (string, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return "", err
	}
	return p.Text(), nil
}

func (s Z) HasText() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Z) TextBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	d := p.Data()
	if len(d) == 0 {
		return d, nil
	}
	return d[:len(d)-1], nil
}

func (s Z) SetText(v string) error {
	s.Struct.SetUint16(0, 13)
	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, t.List.ToPtr())
}

func (s Z) Blob() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	return []byte(p.Data()), nil
}

func (s Z) HasBlob() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Z) SetBlob(v []byte) error {
	s.Struct.SetUint16(0, 14)
	d, err := capnp.NewData(s.Struct.Segment(), []byte(v))
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, d.List.ToPtr())
}

func (s Z) F64vec() (capnp.Float64List, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return capnp.Float64List{}, err
	}
	return capnp.Float64List{List: p.List()}, nil
}

func (s Z) HasF64vec() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Z) SetF64vec(v capnp.Float64List) error {
	s.Struct.SetUint16(0, 15)
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewF64vec sets the f64vec field to a newly
// allocated capnp.Float64List, preferring placement in s's segment.
func (s Z) NewF64vec(n int32) (capnp.Float64List, error) {
	s.Struct.SetUint16(0, 15)
	l, err := capnp.NewFloat64List(s.Struct.Segment(), n)
	if err != nil {
		return capnp.Float64List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

func (s Z) F32vec() (capnp.Float32List, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return capnp.Float32List{}, err
	}
	return capnp.Float32List{List: p.List()}, nil
}

func (s Z) HasF32vec() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Z) SetF32vec(v capnp.Float32List) error {
	s.Struct.SetUint16(0, 16)
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewF32vec sets the f32vec field to a newly
// allocated capnp.Float32List, preferring placement in s's segment.
func (s Z) NewF32vec(n int32) (capnp.Float32List, error) {
	s.Struct.SetUint16(0, 16)
	l, err := capnp.NewFloat32List(s.Struct.Segment(), n)
	if err != nil {
		return capnp.Float32List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

func (s Z) I64vec() (capnp.Int64List, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return capnp.Int64List{}, err
	}
	return capnp.Int64List{List: p.List()}, nil
}

func (s Z) HasI64vec() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Z) SetI64vec(v capnp.Int64List) error {
	s.Struct.SetUint16(0, 17)
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewI64vec sets the i64vec field to a newly
// allocated capnp.Int64List, preferring placement in s's segment.
func (s Z) NewI64vec(n int32) (capnp.Int64List, error) {
	s.Struct.SetUint16(0, 17)
	l, err := capnp.NewInt64List(s.Struct.Segment(), n)
	if err != nil {
		return capnp.Int64List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

func (s Z) I32vec() (capnp.Int32List, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return capnp.Int32List{}, err
	}
	return capnp.Int32List{List: p.List()}, nil
}

func (s Z) HasI32vec() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Z) SetI32vec(v capnp.Int32List) error {
	s.Struct.SetUint16(0, 18)
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewI32vec sets the i32vec field to a newly
// allocated capnp.Int32List, preferring placement in s's segment.
func (s Z) NewI32vec(n int32) (capnp.Int32List, error) {
	s.Struct.SetUint16(0, 18)
	l, err := capnp.NewInt32List(s.Struct.Segment(), n)
	if err != nil {
		return capnp.Int32List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

func (s Z) I16vec() (capnp.Int16List, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return capnp.Int16List{}, err
	}
	return capnp.Int16List{List: p.List()}, nil
}

func (s Z) HasI16vec() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Z) SetI16vec(v capnp.Int16List) error {
	s.Struct.SetUint16(0, 19)
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewI16vec sets the i16vec field to a newly
// allocated capnp.Int16List, preferring placement in s's segment.
func (s Z) NewI16vec(n int32) (capnp.Int16List, error) {
	s.Struct.SetUint16(0, 19)
	l, err := capnp.NewInt16List(s.Struct.Segment(), n)
	if err != nil {
		return capnp.Int16List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

func (s Z) I8vec() (capnp.Int8List, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return capnp.Int8List{}, err
	}
	return capnp.Int8List{List: p.List()}, nil
}

func (s Z) HasI8vec() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Z) SetI8vec(v capnp.Int8List) error {
	s.Struct.SetUint16(0, 20)
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewI8vec sets the i8vec field to a newly
// allocated capnp.Int8List, preferring placement in s's segment.
func (s Z) NewI8vec(n int32) (capnp.Int8List, error) {
	s.Struct.SetUint16(0, 20)
	l, err := capnp.NewInt8List(s.Struct.Segment(), n)
	if err != nil {
		return capnp.Int8List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

func (s Z) U64vec() (capnp.UInt64List, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return capnp.UInt64List{}, err
	}
	return capnp.UInt64List{List: p.List()}, nil
}

func (s Z) HasU64vec() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Z) SetU64vec(v capnp.UInt64List) error {
	s.Struct.SetUint16(0, 21)
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewU64vec sets the u64vec field to a newly
// allocated capnp.UInt64List, preferring placement in s's segment.
func (s Z) NewU64vec(n int32) (capnp.UInt64List, error) {
	s.Struct.SetUint16(0, 21)
	l, err := capnp.NewUInt64List(s.Struct.Segment(), n)
	if err != nil {
		return capnp.UInt64List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

func (s Z) U32vec() (capnp.UInt32List, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return capnp.UInt32List{}, err
	}
	return capnp.UInt32List{List: p.List()}, nil
}

func (s Z) HasU32vec() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Z) SetU32vec(v capnp.UInt32List) error {
	s.Struct.SetUint16(0, 22)
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewU32vec sets the u32vec field to a newly
// allocated capnp.UInt32List, preferring placement in s's segment.
func (s Z) NewU32vec(n int32) (capnp.UInt32List, error) {
	s.Struct.SetUint16(0, 22)
	l, err := capnp.NewUInt32List(s.Struct.Segment(), n)
	if err != nil {
		return capnp.UInt32List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

func (s Z) U16vec() (capnp.UInt16List, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return capnp.UInt16List{}, err
	}
	return capnp.UInt16List{List: p.List()}, nil
}

func (s Z) HasU16vec() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Z) SetU16vec(v capnp.UInt16List) error {
	s.Struct.SetUint16(0, 23)
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewU16vec sets the u16vec field to a newly
// allocated capnp.UInt16List, preferring placement in s's segment.
func (s Z) NewU16vec(n int32) (capnp.UInt16List, error) {
	s.Struct.SetUint16(0, 23)
	l, err := capnp.NewUInt16List(s.Struct.Segment(), n)
	if err != nil {
		return capnp.UInt16List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

func (s Z) U8vec() (capnp.UInt8List, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return capnp.UInt8List{}, err
	}
	return capnp.UInt8List{List: p.List()}, nil
}

func (s Z) HasU8vec() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Z) SetU8vec(v capnp.UInt8List) error {
	s.Struct.SetUint16(0, 24)
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewU8vec sets the u8vec field to a newly
// allocated capnp.UInt8List, preferring placement in s's segment.
func (s Z) NewU8vec(n int32) (capnp.UInt8List, error) {
	s.Struct.SetUint16(0, 24)
	l, err := capnp.NewUInt8List(s.Struct.Segment(), n)
	if err != nil {
		return capnp.UInt8List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

func (s Z) Zvec() (Z_List, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return Z_List{}, err
	}
	return Z_List{List: p.List()}, nil
}

func (s Z) HasZvec() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Z) SetZvec(v Z_List) error {
	s.Struct.SetUint16(0, 25)
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewZvec sets the zvec field to a newly
// allocated Z_List, preferring placement in s's segment.
func (s Z) NewZvec(n int32) (Z_List, error) {
	s.Struct.SetUint16(0, 25)
	l, err := NewZ_List(s.Struct.Segment(), n)
	if err != nil {
		return Z_List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

func (s Z) Zvecvec() (capnp.PointerList, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return capnp.PointerList{}, err
	}
	return capnp.PointerList{List: p.List()}, nil
}

func (s Z) HasZvecvec() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Z) SetZvecvec(v capnp.PointerList) error {
	s.Struct.SetUint16(0, 26)
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewZvecvec sets the zvecvec field to a newly
// allocated capnp.PointerList, preferring placement in s's segment.
func (s Z) NewZvecvec(n int32) (capnp.PointerList, error) {
	s.Struct.SetUint16(0, 26)
	l, err := capnp.NewPointerList(s.Struct.Segment(), n)
	if err != nil {
		return capnp.PointerList{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

func (s Z) Zdate() (Zdate, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return Zdate{}, err
	}
	return Zdate{Struct: p.Struct()}, nil
}

func (s Z) HasZdate() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Z) SetZdate(v Zdate) error {
	s.Struct.SetUint16(0, 27)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewZdate sets the zdate field to a newly
// allocated Zdate struct, preferring placement in s's segment.
func (s Z) NewZdate() (Zdate, error) {
	s.Struct.SetUint16(0, 27)
	ss, err := NewZdate(s.Struct.Segment())
	if err != nil {
		return Zdate{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Z) Zdata() (Zdata, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return Zdata{}, err
	}
	return Zdata{Struct: p.Struct()}, nil
}

func (s Z) HasZdata() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Z) SetZdata(v Zdata) error {
	s.Struct.SetUint16(0, 28)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewZdata sets the zdata field to a newly
// allocated Zdata struct, preferring placement in s's segment.
func (s Z) NewZdata() (Zdata, error) {
	s.Struct.SetUint16(0, 28)
	ss, err := NewZdata(s.Struct.Segment())
	if err != nil {
		return Zdata{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Z) Aircraftvec() (Aircraft_List, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return Aircraft_List{}, err
	}
	return Aircraft_List{List: p.List()}, nil
}

func (s Z) HasAircraftvec() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Z) SetAircraftvec(v Aircraft_List) error {
	s.Struct.SetUint16(0, 29)
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewAircraftvec sets the aircraftvec field to a newly
// allocated Aircraft_List, preferring placement in s's segment.
func (s Z) NewAircraftvec(n int32) (Aircraft_List, error) {
	s.Struct.SetUint16(0, 29)
	l, err := NewAircraft_List(s.Struct.Segment(), n)
	if err != nil {
		return Aircraft_List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

func (s Z) Aircraft() (Aircraft, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return Aircraft{}, err
	}
	return Aircraft{Struct: p.Struct()}, nil
}

func (s Z) HasAircraft() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Z) SetAircraft(v Aircraft) error {
	s.Struct.SetUint16(0, 30)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewAircraft sets the aircraft field to a newly
// allocated Aircraft struct, preferring placement in s's segment.
func (s Z) NewAircraft() (Aircraft, error) {
	s.Struct.SetUint16(0, 30)
	ss, err := NewAircraft(s.Struct.Segment())
	if err != nil {
		return Aircraft{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Z) Regression() (Regression, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return Regression{}, err
	}
	return Regression{Struct: p.Struct()}, nil
}

func (s Z) HasRegression() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Z) SetRegression(v Regression) error {
	s.Struct.SetUint16(0, 31)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewRegression sets the regression field to a newly
// allocated Regression struct, preferring placement in s's segment.
func (s Z) NewRegression() (Regression, error) {
	s.Struct.SetUint16(0, 31)
	ss, err := NewRegression(s.Struct.Segment())
	if err != nil {
		return Regression{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Z) Planebase() (PlaneBase, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return PlaneBase{}, err
	}
	return PlaneBase{Struct: p.Struct()}, nil
}

func (s Z) HasPlanebase() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Z) SetPlanebase(v PlaneBase) error {
	s.Struct.SetUint16(0, 32)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewPlanebase sets the planebase field to a newly
// allocated PlaneBase struct, preferring placement in s's segment.
func (s Z) NewPlanebase() (PlaneBase, error) {
	s.Struct.SetUint16(0, 32)
	ss, err := NewPlaneBase(s.Struct.Segment())
	if err != nil {
		return PlaneBase{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Z) Airport() Airport {
	return Airport(s.Struct.Uint16(8))
}

func (s Z) SetAirport(v Airport) {
	s.Struct.SetUint16(0, 33)
	s.Struct.SetUint16(8, uint16(v))
}

func (s Z) B737() (B737, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return B737{}, err
	}
	return B737{Struct: p.Struct()}, nil
}

func (s Z) HasB737() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Z) SetB737(v B737) error {
	s.Struct.SetUint16(0, 34)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewB737 sets the b737 field to a newly
// allocated B737 struct, preferring placement in s's segment.
func (s Z) NewB737() (B737, error) {
	s.Struct.SetUint16(0, 34)
	ss, err := NewB737(s.Struct.Segment())
	if err != nil {
		return B737{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Z) A320() (A320, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return A320{}, err
	}
	return A320{Struct: p.Struct()}, nil
}

func (s Z) HasA320() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Z) SetA320(v A320) error {
	s.Struct.SetUint16(0, 35)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewA320 sets the a320 field to a newly
// allocated A320 struct, preferring placement in s's segment.
func (s Z) NewA320() (A320, error) {
	s.Struct.SetUint16(0, 35)
	ss, err := NewA320(s.Struct.Segment())
	if err != nil {
		return A320{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Z) F16() (F16, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return F16{}, err
	}
	return F16{Struct: p.Struct()}, nil
}

func (s Z) HasF16() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Z) SetF16(v F16) error {
	s.Struct.SetUint16(0, 36)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewF16 sets the f16 field to a newly
// allocated F16 struct, preferring placement in s's segment.
func (s Z) NewF16() (F16, error) {
	s.Struct.SetUint16(0, 36)
	ss, err := NewF16(s.Struct.Segment())
	if err != nil {
		return F16{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Z) Zdatevec() (Zdate_List, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return Zdate_List{}, err
	}
	return Zdate_List{List: p.List()}, nil
}

func (s Z) HasZdatevec() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Z) SetZdatevec(v Zdate_List) error {
	s.Struct.SetUint16(0, 37)
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewZdatevec sets the zdatevec field to a newly
// allocated Zdate_List, preferring placement in s's segment.
func (s Z) NewZdatevec(n int32) (Zdate_List, error) {
	s.Struct.SetUint16(0, 37)
	l, err := NewZdate_List(s.Struct.Segment(), n)
	if err != nil {
		return Zdate_List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

func (s Z) Zdatavec() (Zdata_List, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return Zdata_List{}, err
	}
	return Zdata_List{List: p.List()}, nil
}

func (s Z) HasZdatavec() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Z) SetZdatavec(v Zdata_List) error {
	s.Struct.SetUint16(0, 38)
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewZdatavec sets the zdatavec field to a newly
// allocated Zdata_List, preferring placement in s's segment.
func (s Z) NewZdatavec(n int32) (Zdata_List, error) {
	s.Struct.SetUint16(0, 38)
	l, err := NewZdata_List(s.Struct.Segment(), n)
	if err != nil {
		return Zdata_List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

func (s Z) Boolvec() (capnp.BitList, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return capnp.BitList{}, err
	}
	return capnp.BitList{List: p.List()}, nil
}

func (s Z) HasBoolvec() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Z) SetBoolvec(v capnp.BitList) error {
	s.Struct.SetUint16(0, 39)
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewBoolvec sets the boolvec field to a newly
// allocated capnp.BitList, preferring placement in s's segment.
func (s Z) NewBoolvec(n int32) (capnp.BitList, error) {
	s.Struct.SetUint16(0, 39)
	l, err := capnp.NewBitList(s.Struct.Segment(), n)
	if err != nil {
		return capnp.BitList{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

// Z_List is a list of Z.
type Z_List struct{ capnp.List }

// NewZ creates a new list of Z.
func NewZ_List(s *capnp.Segment, sz int32) (Z_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 16, PointerCount: 1}, sz)
	if err != nil {
		return Z_List{}, err
	}
	return Z_List{l}, nil
}

func (s Z_List) At(i int) Z           { return Z{s.List.Struct(i)} }
func (s Z_List) Set(i int, v Z) error { return s.List.SetStruct(i, v.Struct) }

// Z_Promise is a wrapper for a Z promised by a client call.
type Z_Promise struct{ *capnp.Pipeline }

func (p Z_Promise) Struct() (Z, error) {
	s, err := p.Pipeline.Struct()
	return Z{s}, err
}

func (p Z_Promise) Zz() Z_Promise {
	return Z_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Z_Promise) Zdate() Zdate_Promise {
	return Zdate_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Z_Promise) Zdata() Zdata_Promise {
	return Zdata_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Z_Promise) Aircraft() Aircraft_Promise {
	return Aircraft_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Z_Promise) Regression() Regression_Promise {
	return Regression_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Z_Promise) Planebase() PlaneBase_Promise {
	return PlaneBase_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Z_Promise) B737() B737_Promise {
	return B737_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Z_Promise) A320() A320_Promise {
	return A320_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Z_Promise) F16() F16_Promise {
	return F16_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

type Counter struct{ capnp.Struct }

func NewCounter(s *capnp.Segment) (Counter, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2})
	if err != nil {
		return Counter{}, err
	}
	return Counter{st}, nil
}

func NewRootCounter(s *capnp.Segment) (Counter, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2})
	if err != nil {
		return Counter{}, err
	}
	return Counter{st}, nil
}

func ReadRootCounter(msg *capnp.Message) (Counter, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Counter{}, err
	}
	return Counter{root.Struct()}, nil
}
func (s Counter) Size() int64 {
	return int64(s.Struct.Uint64(0))
}

func (s Counter) SetSize(v int64) {
	s.Struct.SetUint64(0, uint64(v))
}

func (s Counter) Words() (string, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return "", err
	}
	return p.Text(), nil
}

func (s Counter) HasWords() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Counter) WordsBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	d := p.Data()
	if len(d) == 0 {
		return d, nil
	}
	return d[:len(d)-1], nil
}

func (s Counter) SetWords(v string) error {
	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, t.List.ToPtr())
}

func (s Counter) Wordlist() (capnp.TextList, error) {
	p, err := s.Struct.Ptr(1)
	if err != nil {
		return capnp.TextList{}, err
	}
	return capnp.TextList{List: p.List()}, nil
}

func (s Counter) HasWordlist() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s Counter) SetWordlist(v capnp.TextList) error {
	return s.Struct.SetPtr(1, v.List.ToPtr())
}

// NewWordlist sets the wordlist field to a newly
// allocated capnp.TextList, preferring placement in s's segment.
func (s Counter) NewWordlist(n int32) (capnp.TextList, error) {
	l, err := capnp.NewTextList(s.Struct.Segment(), n)
	if err != nil {
		return capnp.TextList{}, err
	}
	err = s.Struct.SetPtr(1, l.List.ToPtr())
	return l, err
}

// Counter_List is a list of Counter.
type Counter_List struct{ capnp.List }

// NewCounter creates a new list of Counter.
func NewCounter_List(s *capnp.Segment, sz int32) (Counter_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2}, sz)
	if err != nil {
		return Counter_List{}, err
	}
	return Counter_List{l}, nil
}

func (s Counter_List) At(i int) Counter           { return Counter{s.List.Struct(i)} }
func (s Counter_List) Set(i int, v Counter) error { return s.List.SetStruct(i, v.Struct) }

// Counter_Promise is a wrapper for a Counter promised by a client call.
type Counter_Promise struct{ *capnp.Pipeline }

func (p Counter_Promise) Struct() (Counter, error) {
	s, err := p.Pipeline.Struct()
	return Counter{s}, err
}

type Bag struct{ capnp.Struct }

func NewBag(s *capnp.Segment) (Bag, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Bag{}, err
	}
	return Bag{st}, nil
}

func NewRootBag(s *capnp.Segment) (Bag, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Bag{}, err
	}
	return Bag{st}, nil
}

func ReadRootBag(msg *capnp.Message) (Bag, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Bag{}, err
	}
	return Bag{root.Struct()}, nil
}
func (s Bag) Counter() (Counter, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return Counter{}, err
	}
	return Counter{Struct: p.Struct()}, nil
}

func (s Bag) HasCounter() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Bag) SetCounter(v Counter) error {
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewCounter sets the counter field to a newly
// allocated Counter struct, preferring placement in s's segment.
func (s Bag) NewCounter() (Counter, error) {
	ss, err := NewCounter(s.Struct.Segment())
	if err != nil {
		return Counter{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

// Bag_List is a list of Bag.
type Bag_List struct{ capnp.List }

// NewBag creates a new list of Bag.
func NewBag_List(s *capnp.Segment, sz int32) (Bag_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return Bag_List{}, err
	}
	return Bag_List{l}, nil
}

func (s Bag_List) At(i int) Bag           { return Bag{s.List.Struct(i)} }
func (s Bag_List) Set(i int, v Bag) error { return s.List.SetStruct(i, v.Struct) }

// Bag_Promise is a wrapper for a Bag promised by a client call.
type Bag_Promise struct{ *capnp.Pipeline }

func (p Bag_Promise) Struct() (Bag, error) {
	s, err := p.Pipeline.Struct()
	return Bag{s}, err
}

func (p Bag_Promise) Counter() Counter_Promise {
	return Counter_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

type Zserver struct{ capnp.Struct }

func NewZserver(s *capnp.Segment) (Zserver, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Zserver{}, err
	}
	return Zserver{st}, nil
}

func NewRootZserver(s *capnp.Segment) (Zserver, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Zserver{}, err
	}
	return Zserver{st}, nil
}

func ReadRootZserver(msg *capnp.Message) (Zserver, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Zserver{}, err
	}
	return Zserver{root.Struct()}, nil
}
func (s Zserver) Waitingjobs() (Zjob_List, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return Zjob_List{}, err
	}
	return Zjob_List{List: p.List()}, nil
}

func (s Zserver) HasWaitingjobs() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Zserver) SetWaitingjobs(v Zjob_List) error {
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewWaitingjobs sets the waitingjobs field to a newly
// allocated Zjob_List, preferring placement in s's segment.
func (s Zserver) NewWaitingjobs(n int32) (Zjob_List, error) {
	l, err := NewZjob_List(s.Struct.Segment(), n)
	if err != nil {
		return Zjob_List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

// Zserver_List is a list of Zserver.
type Zserver_List struct{ capnp.List }

// NewZserver creates a new list of Zserver.
func NewZserver_List(s *capnp.Segment, sz int32) (Zserver_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return Zserver_List{}, err
	}
	return Zserver_List{l}, nil
}

func (s Zserver_List) At(i int) Zserver           { return Zserver{s.List.Struct(i)} }
func (s Zserver_List) Set(i int, v Zserver) error { return s.List.SetStruct(i, v.Struct) }

// Zserver_Promise is a wrapper for a Zserver promised by a client call.
type Zserver_Promise struct{ *capnp.Pipeline }

func (p Zserver_Promise) Struct() (Zserver, error) {
	s, err := p.Pipeline.Struct()
	return Zserver{s}, err
}

type Zjob struct{ capnp.Struct }

func NewZjob(s *capnp.Segment) (Zjob, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	if err != nil {
		return Zjob{}, err
	}
	return Zjob{st}, nil
}

func NewRootZjob(s *capnp.Segment) (Zjob, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	if err != nil {
		return Zjob{}, err
	}
	return Zjob{st}, nil
}

func ReadRootZjob(msg *capnp.Message) (Zjob, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Zjob{}, err
	}
	return Zjob{root.Struct()}, nil
}
func (s Zjob) Cmd() (string, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return "", err
	}
	return p.Text(), nil
}

func (s Zjob) HasCmd() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Zjob) CmdBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	d := p.Data()
	if len(d) == 0 {
		return d, nil
	}
	return d[:len(d)-1], nil
}

func (s Zjob) SetCmd(v string) error {
	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, t.List.ToPtr())
}

func (s Zjob) Args() (capnp.TextList, error) {
	p, err := s.Struct.Ptr(1)
	if err != nil {
		return capnp.TextList{}, err
	}
	return capnp.TextList{List: p.List()}, nil
}

func (s Zjob) HasArgs() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s Zjob) SetArgs(v capnp.TextList) error {
	return s.Struct.SetPtr(1, v.List.ToPtr())
}

// NewArgs sets the args field to a newly
// allocated capnp.TextList, preferring placement in s's segment.
func (s Zjob) NewArgs(n int32) (capnp.TextList, error) {
	l, err := capnp.NewTextList(s.Struct.Segment(), n)
	if err != nil {
		return capnp.TextList{}, err
	}
	err = s.Struct.SetPtr(1, l.List.ToPtr())
	return l, err
}

// Zjob_List is a list of Zjob.
type Zjob_List struct{ capnp.List }

// NewZjob creates a new list of Zjob.
func NewZjob_List(s *capnp.Segment, sz int32) (Zjob_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2}, sz)
	if err != nil {
		return Zjob_List{}, err
	}
	return Zjob_List{l}, nil
}

func (s Zjob_List) At(i int) Zjob           { return Zjob{s.List.Struct(i)} }
func (s Zjob_List) Set(i int, v Zjob) error { return s.List.SetStruct(i, v.Struct) }

// Zjob_Promise is a wrapper for a Zjob promised by a client call.
type Zjob_Promise struct{ *capnp.Pipeline }

func (p Zjob_Promise) Struct() (Zjob, error) {
	s, err := p.Pipeline.Struct()
	return Zjob{s}, err
}

type VerEmpty struct{ capnp.Struct }

func NewVerEmpty(s *capnp.Segment) (VerEmpty, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0})
	if err != nil {
		return VerEmpty{}, err
	}
	return VerEmpty{st}, nil
}

func NewRootVerEmpty(s *capnp.Segment) (VerEmpty, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0})
	if err != nil {
		return VerEmpty{}, err
	}
	return VerEmpty{st}, nil
}

func ReadRootVerEmpty(msg *capnp.Message) (VerEmpty, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return VerEmpty{}, err
	}
	return VerEmpty{root.Struct()}, nil
}

// VerEmpty_List is a list of VerEmpty.
type VerEmpty_List struct{ capnp.List }

// NewVerEmpty creates a new list of VerEmpty.
func NewVerEmpty_List(s *capnp.Segment, sz int32) (VerEmpty_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0}, sz)
	if err != nil {
		return VerEmpty_List{}, err
	}
	return VerEmpty_List{l}, nil
}

func (s VerEmpty_List) At(i int) VerEmpty           { return VerEmpty{s.List.Struct(i)} }
func (s VerEmpty_List) Set(i int, v VerEmpty) error { return s.List.SetStruct(i, v.Struct) }

// VerEmpty_Promise is a wrapper for a VerEmpty promised by a client call.
type VerEmpty_Promise struct{ *capnp.Pipeline }

func (p VerEmpty_Promise) Struct() (VerEmpty, error) {
	s, err := p.Pipeline.Struct()
	return VerEmpty{s}, err
}

type VerOneData struct{ capnp.Struct }

func NewVerOneData(s *capnp.Segment) (VerOneData, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	if err != nil {
		return VerOneData{}, err
	}
	return VerOneData{st}, nil
}

func NewRootVerOneData(s *capnp.Segment) (VerOneData, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	if err != nil {
		return VerOneData{}, err
	}
	return VerOneData{st}, nil
}

func ReadRootVerOneData(msg *capnp.Message) (VerOneData, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return VerOneData{}, err
	}
	return VerOneData{root.Struct()}, nil
}
func (s VerOneData) Val() int16 {
	return int16(s.Struct.Uint16(0))
}

func (s VerOneData) SetVal(v int16) {
	s.Struct.SetUint16(0, uint16(v))
}

// VerOneData_List is a list of VerOneData.
type VerOneData_List struct{ capnp.List }

// NewVerOneData creates a new list of VerOneData.
func NewVerOneData_List(s *capnp.Segment, sz int32) (VerOneData_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0}, sz)
	if err != nil {
		return VerOneData_List{}, err
	}
	return VerOneData_List{l}, nil
}

func (s VerOneData_List) At(i int) VerOneData           { return VerOneData{s.List.Struct(i)} }
func (s VerOneData_List) Set(i int, v VerOneData) error { return s.List.SetStruct(i, v.Struct) }

// VerOneData_Promise is a wrapper for a VerOneData promised by a client call.
type VerOneData_Promise struct{ *capnp.Pipeline }

func (p VerOneData_Promise) Struct() (VerOneData, error) {
	s, err := p.Pipeline.Struct()
	return VerOneData{s}, err
}

type VerTwoData struct{ capnp.Struct }

func NewVerTwoData(s *capnp.Segment) (VerTwoData, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 0})
	if err != nil {
		return VerTwoData{}, err
	}
	return VerTwoData{st}, nil
}

func NewRootVerTwoData(s *capnp.Segment) (VerTwoData, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 0})
	if err != nil {
		return VerTwoData{}, err
	}
	return VerTwoData{st}, nil
}

func ReadRootVerTwoData(msg *capnp.Message) (VerTwoData, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return VerTwoData{}, err
	}
	return VerTwoData{root.Struct()}, nil
}
func (s VerTwoData) Val() int16 {
	return int16(s.Struct.Uint16(0))
}

func (s VerTwoData) SetVal(v int16) {
	s.Struct.SetUint16(0, uint16(v))
}

func (s VerTwoData) Duo() int64 {
	return int64(s.Struct.Uint64(8))
}

func (s VerTwoData) SetDuo(v int64) {
	s.Struct.SetUint64(8, uint64(v))
}

// VerTwoData_List is a list of VerTwoData.
type VerTwoData_List struct{ capnp.List }

// NewVerTwoData creates a new list of VerTwoData.
func NewVerTwoData_List(s *capnp.Segment, sz int32) (VerTwoData_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 16, PointerCount: 0}, sz)
	if err != nil {
		return VerTwoData_List{}, err
	}
	return VerTwoData_List{l}, nil
}

func (s VerTwoData_List) At(i int) VerTwoData           { return VerTwoData{s.List.Struct(i)} }
func (s VerTwoData_List) Set(i int, v VerTwoData) error { return s.List.SetStruct(i, v.Struct) }

// VerTwoData_Promise is a wrapper for a VerTwoData promised by a client call.
type VerTwoData_Promise struct{ *capnp.Pipeline }

func (p VerTwoData_Promise) Struct() (VerTwoData, error) {
	s, err := p.Pipeline.Struct()
	return VerTwoData{s}, err
}

type VerOnePtr struct{ capnp.Struct }

func NewVerOnePtr(s *capnp.Segment) (VerOnePtr, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return VerOnePtr{}, err
	}
	return VerOnePtr{st}, nil
}

func NewRootVerOnePtr(s *capnp.Segment) (VerOnePtr, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return VerOnePtr{}, err
	}
	return VerOnePtr{st}, nil
}

func ReadRootVerOnePtr(msg *capnp.Message) (VerOnePtr, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return VerOnePtr{}, err
	}
	return VerOnePtr{root.Struct()}, nil
}
func (s VerOnePtr) Ptr() (VerOneData, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return VerOneData{}, err
	}
	return VerOneData{Struct: p.Struct()}, nil
}

func (s VerOnePtr) HasPtr() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s VerOnePtr) SetPtr(v VerOneData) error {
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewPtr sets the ptr field to a newly
// allocated VerOneData struct, preferring placement in s's segment.
func (s VerOnePtr) NewPtr() (VerOneData, error) {
	ss, err := NewVerOneData(s.Struct.Segment())
	if err != nil {
		return VerOneData{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

// VerOnePtr_List is a list of VerOnePtr.
type VerOnePtr_List struct{ capnp.List }

// NewVerOnePtr creates a new list of VerOnePtr.
func NewVerOnePtr_List(s *capnp.Segment, sz int32) (VerOnePtr_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return VerOnePtr_List{}, err
	}
	return VerOnePtr_List{l}, nil
}

func (s VerOnePtr_List) At(i int) VerOnePtr           { return VerOnePtr{s.List.Struct(i)} }
func (s VerOnePtr_List) Set(i int, v VerOnePtr) error { return s.List.SetStruct(i, v.Struct) }

// VerOnePtr_Promise is a wrapper for a VerOnePtr promised by a client call.
type VerOnePtr_Promise struct{ *capnp.Pipeline }

func (p VerOnePtr_Promise) Struct() (VerOnePtr, error) {
	s, err := p.Pipeline.Struct()
	return VerOnePtr{s}, err
}

func (p VerOnePtr_Promise) Ptr() VerOneData_Promise {
	return VerOneData_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

type VerTwoPtr struct{ capnp.Struct }

func NewVerTwoPtr(s *capnp.Segment) (VerTwoPtr, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	if err != nil {
		return VerTwoPtr{}, err
	}
	return VerTwoPtr{st}, nil
}

func NewRootVerTwoPtr(s *capnp.Segment) (VerTwoPtr, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	if err != nil {
		return VerTwoPtr{}, err
	}
	return VerTwoPtr{st}, nil
}

func ReadRootVerTwoPtr(msg *capnp.Message) (VerTwoPtr, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return VerTwoPtr{}, err
	}
	return VerTwoPtr{root.Struct()}, nil
}
func (s VerTwoPtr) Ptr1() (VerOneData, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return VerOneData{}, err
	}
	return VerOneData{Struct: p.Struct()}, nil
}

func (s VerTwoPtr) HasPtr1() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s VerTwoPtr) SetPtr1(v VerOneData) error {
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewPtr1 sets the ptr1 field to a newly
// allocated VerOneData struct, preferring placement in s's segment.
func (s VerTwoPtr) NewPtr1() (VerOneData, error) {
	ss, err := NewVerOneData(s.Struct.Segment())
	if err != nil {
		return VerOneData{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s VerTwoPtr) Ptr2() (VerOneData, error) {
	p, err := s.Struct.Ptr(1)
	if err != nil {
		return VerOneData{}, err
	}
	return VerOneData{Struct: p.Struct()}, nil
}

func (s VerTwoPtr) HasPtr2() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s VerTwoPtr) SetPtr2(v VerOneData) error {
	return s.Struct.SetPtr(1, v.Struct.ToPtr())
}

// NewPtr2 sets the ptr2 field to a newly
// allocated VerOneData struct, preferring placement in s's segment.
func (s VerTwoPtr) NewPtr2() (VerOneData, error) {
	ss, err := NewVerOneData(s.Struct.Segment())
	if err != nil {
		return VerOneData{}, err
	}
	err = s.Struct.SetPtr(1, ss.Struct.ToPtr())
	return ss, err
}

// VerTwoPtr_List is a list of VerTwoPtr.
type VerTwoPtr_List struct{ capnp.List }

// NewVerTwoPtr creates a new list of VerTwoPtr.
func NewVerTwoPtr_List(s *capnp.Segment, sz int32) (VerTwoPtr_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2}, sz)
	if err != nil {
		return VerTwoPtr_List{}, err
	}
	return VerTwoPtr_List{l}, nil
}

func (s VerTwoPtr_List) At(i int) VerTwoPtr           { return VerTwoPtr{s.List.Struct(i)} }
func (s VerTwoPtr_List) Set(i int, v VerTwoPtr) error { return s.List.SetStruct(i, v.Struct) }

// VerTwoPtr_Promise is a wrapper for a VerTwoPtr promised by a client call.
type VerTwoPtr_Promise struct{ *capnp.Pipeline }

func (p VerTwoPtr_Promise) Struct() (VerTwoPtr, error) {
	s, err := p.Pipeline.Struct()
	return VerTwoPtr{s}, err
}

func (p VerTwoPtr_Promise) Ptr1() VerOneData_Promise {
	return VerOneData_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p VerTwoPtr_Promise) Ptr2() VerOneData_Promise {
	return VerOneData_Promise{Pipeline: p.Pipeline.GetPipeline(1)}
}

type VerTwoDataTwoPtr struct{ capnp.Struct }

func NewVerTwoDataTwoPtr(s *capnp.Segment) (VerTwoDataTwoPtr, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 2})
	if err != nil {
		return VerTwoDataTwoPtr{}, err
	}
	return VerTwoDataTwoPtr{st}, nil
}

func NewRootVerTwoDataTwoPtr(s *capnp.Segment) (VerTwoDataTwoPtr, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 2})
	if err != nil {
		return VerTwoDataTwoPtr{}, err
	}
	return VerTwoDataTwoPtr{st}, nil
}

func ReadRootVerTwoDataTwoPtr(msg *capnp.Message) (VerTwoDataTwoPtr, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return VerTwoDataTwoPtr{}, err
	}
	return VerTwoDataTwoPtr{root.Struct()}, nil
}
func (s VerTwoDataTwoPtr) Val() int16 {
	return int16(s.Struct.Uint16(0))
}

func (s VerTwoDataTwoPtr) SetVal(v int16) {
	s.Struct.SetUint16(0, uint16(v))
}

func (s VerTwoDataTwoPtr) Duo() int64 {
	return int64(s.Struct.Uint64(8))
}

func (s VerTwoDataTwoPtr) SetDuo(v int64) {
	s.Struct.SetUint64(8, uint64(v))
}

func (s VerTwoDataTwoPtr) Ptr1() (VerOneData, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return VerOneData{}, err
	}
	return VerOneData{Struct: p.Struct()}, nil
}

func (s VerTwoDataTwoPtr) HasPtr1() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s VerTwoDataTwoPtr) SetPtr1(v VerOneData) error {
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewPtr1 sets the ptr1 field to a newly
// allocated VerOneData struct, preferring placement in s's segment.
func (s VerTwoDataTwoPtr) NewPtr1() (VerOneData, error) {
	ss, err := NewVerOneData(s.Struct.Segment())
	if err != nil {
		return VerOneData{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s VerTwoDataTwoPtr) Ptr2() (VerOneData, error) {
	p, err := s.Struct.Ptr(1)
	if err != nil {
		return VerOneData{}, err
	}
	return VerOneData{Struct: p.Struct()}, nil
}

func (s VerTwoDataTwoPtr) HasPtr2() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s VerTwoDataTwoPtr) SetPtr2(v VerOneData) error {
	return s.Struct.SetPtr(1, v.Struct.ToPtr())
}

// NewPtr2 sets the ptr2 field to a newly
// allocated VerOneData struct, preferring placement in s's segment.
func (s VerTwoDataTwoPtr) NewPtr2() (VerOneData, error) {
	ss, err := NewVerOneData(s.Struct.Segment())
	if err != nil {
		return VerOneData{}, err
	}
	err = s.Struct.SetPtr(1, ss.Struct.ToPtr())
	return ss, err
}

// VerTwoDataTwoPtr_List is a list of VerTwoDataTwoPtr.
type VerTwoDataTwoPtr_List struct{ capnp.List }

// NewVerTwoDataTwoPtr creates a new list of VerTwoDataTwoPtr.
func NewVerTwoDataTwoPtr_List(s *capnp.Segment, sz int32) (VerTwoDataTwoPtr_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 16, PointerCount: 2}, sz)
	if err != nil {
		return VerTwoDataTwoPtr_List{}, err
	}
	return VerTwoDataTwoPtr_List{l}, nil
}

func (s VerTwoDataTwoPtr_List) At(i int) VerTwoDataTwoPtr { return VerTwoDataTwoPtr{s.List.Struct(i)} }
func (s VerTwoDataTwoPtr_List) Set(i int, v VerTwoDataTwoPtr) error {
	return s.List.SetStruct(i, v.Struct)
}

// VerTwoDataTwoPtr_Promise is a wrapper for a VerTwoDataTwoPtr promised by a client call.
type VerTwoDataTwoPtr_Promise struct{ *capnp.Pipeline }

func (p VerTwoDataTwoPtr_Promise) Struct() (VerTwoDataTwoPtr, error) {
	s, err := p.Pipeline.Struct()
	return VerTwoDataTwoPtr{s}, err
}

func (p VerTwoDataTwoPtr_Promise) Ptr1() VerOneData_Promise {
	return VerOneData_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p VerTwoDataTwoPtr_Promise) Ptr2() VerOneData_Promise {
	return VerOneData_Promise{Pipeline: p.Pipeline.GetPipeline(1)}
}

type HoldsVerEmptyList struct{ capnp.Struct }

func NewHoldsVerEmptyList(s *capnp.Segment) (HoldsVerEmptyList, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return HoldsVerEmptyList{}, err
	}
	return HoldsVerEmptyList{st}, nil
}

func NewRootHoldsVerEmptyList(s *capnp.Segment) (HoldsVerEmptyList, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return HoldsVerEmptyList{}, err
	}
	return HoldsVerEmptyList{st}, nil
}

func ReadRootHoldsVerEmptyList(msg *capnp.Message) (HoldsVerEmptyList, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return HoldsVerEmptyList{}, err
	}
	return HoldsVerEmptyList{root.Struct()}, nil
}
func (s HoldsVerEmptyList) Mylist() (VerEmpty_List, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return VerEmpty_List{}, err
	}
	return VerEmpty_List{List: p.List()}, nil
}

func (s HoldsVerEmptyList) HasMylist() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s HoldsVerEmptyList) SetMylist(v VerEmpty_List) error {
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewMylist sets the mylist field to a newly
// allocated VerEmpty_List, preferring placement in s's segment.
func (s HoldsVerEmptyList) NewMylist(n int32) (VerEmpty_List, error) {
	l, err := NewVerEmpty_List(s.Struct.Segment(), n)
	if err != nil {
		return VerEmpty_List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

// HoldsVerEmptyList_List is a list of HoldsVerEmptyList.
type HoldsVerEmptyList_List struct{ capnp.List }

// NewHoldsVerEmptyList creates a new list of HoldsVerEmptyList.
func NewHoldsVerEmptyList_List(s *capnp.Segment, sz int32) (HoldsVerEmptyList_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return HoldsVerEmptyList_List{}, err
	}
	return HoldsVerEmptyList_List{l}, nil
}

func (s HoldsVerEmptyList_List) At(i int) HoldsVerEmptyList {
	return HoldsVerEmptyList{s.List.Struct(i)}
}
func (s HoldsVerEmptyList_List) Set(i int, v HoldsVerEmptyList) error {
	return s.List.SetStruct(i, v.Struct)
}

// HoldsVerEmptyList_Promise is a wrapper for a HoldsVerEmptyList promised by a client call.
type HoldsVerEmptyList_Promise struct{ *capnp.Pipeline }

func (p HoldsVerEmptyList_Promise) Struct() (HoldsVerEmptyList, error) {
	s, err := p.Pipeline.Struct()
	return HoldsVerEmptyList{s}, err
}

type HoldsVerOneDataList struct{ capnp.Struct }

func NewHoldsVerOneDataList(s *capnp.Segment) (HoldsVerOneDataList, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return HoldsVerOneDataList{}, err
	}
	return HoldsVerOneDataList{st}, nil
}

func NewRootHoldsVerOneDataList(s *capnp.Segment) (HoldsVerOneDataList, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return HoldsVerOneDataList{}, err
	}
	return HoldsVerOneDataList{st}, nil
}

func ReadRootHoldsVerOneDataList(msg *capnp.Message) (HoldsVerOneDataList, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return HoldsVerOneDataList{}, err
	}
	return HoldsVerOneDataList{root.Struct()}, nil
}
func (s HoldsVerOneDataList) Mylist() (VerOneData_List, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return VerOneData_List{}, err
	}
	return VerOneData_List{List: p.List()}, nil
}

func (s HoldsVerOneDataList) HasMylist() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s HoldsVerOneDataList) SetMylist(v VerOneData_List) error {
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewMylist sets the mylist field to a newly
// allocated VerOneData_List, preferring placement in s's segment.
func (s HoldsVerOneDataList) NewMylist(n int32) (VerOneData_List, error) {
	l, err := NewVerOneData_List(s.Struct.Segment(), n)
	if err != nil {
		return VerOneData_List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

// HoldsVerOneDataList_List is a list of HoldsVerOneDataList.
type HoldsVerOneDataList_List struct{ capnp.List }

// NewHoldsVerOneDataList creates a new list of HoldsVerOneDataList.
func NewHoldsVerOneDataList_List(s *capnp.Segment, sz int32) (HoldsVerOneDataList_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return HoldsVerOneDataList_List{}, err
	}
	return HoldsVerOneDataList_List{l}, nil
}

func (s HoldsVerOneDataList_List) At(i int) HoldsVerOneDataList {
	return HoldsVerOneDataList{s.List.Struct(i)}
}
func (s HoldsVerOneDataList_List) Set(i int, v HoldsVerOneDataList) error {
	return s.List.SetStruct(i, v.Struct)
}

// HoldsVerOneDataList_Promise is a wrapper for a HoldsVerOneDataList promised by a client call.
type HoldsVerOneDataList_Promise struct{ *capnp.Pipeline }

func (p HoldsVerOneDataList_Promise) Struct() (HoldsVerOneDataList, error) {
	s, err := p.Pipeline.Struct()
	return HoldsVerOneDataList{s}, err
}

type HoldsVerTwoDataList struct{ capnp.Struct }

func NewHoldsVerTwoDataList(s *capnp.Segment) (HoldsVerTwoDataList, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return HoldsVerTwoDataList{}, err
	}
	return HoldsVerTwoDataList{st}, nil
}

func NewRootHoldsVerTwoDataList(s *capnp.Segment) (HoldsVerTwoDataList, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return HoldsVerTwoDataList{}, err
	}
	return HoldsVerTwoDataList{st}, nil
}

func ReadRootHoldsVerTwoDataList(msg *capnp.Message) (HoldsVerTwoDataList, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return HoldsVerTwoDataList{}, err
	}
	return HoldsVerTwoDataList{root.Struct()}, nil
}
func (s HoldsVerTwoDataList) Mylist() (VerTwoData_List, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return VerTwoData_List{}, err
	}
	return VerTwoData_List{List: p.List()}, nil
}

func (s HoldsVerTwoDataList) HasMylist() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s HoldsVerTwoDataList) SetMylist(v VerTwoData_List) error {
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewMylist sets the mylist field to a newly
// allocated VerTwoData_List, preferring placement in s's segment.
func (s HoldsVerTwoDataList) NewMylist(n int32) (VerTwoData_List, error) {
	l, err := NewVerTwoData_List(s.Struct.Segment(), n)
	if err != nil {
		return VerTwoData_List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

// HoldsVerTwoDataList_List is a list of HoldsVerTwoDataList.
type HoldsVerTwoDataList_List struct{ capnp.List }

// NewHoldsVerTwoDataList creates a new list of HoldsVerTwoDataList.
func NewHoldsVerTwoDataList_List(s *capnp.Segment, sz int32) (HoldsVerTwoDataList_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return HoldsVerTwoDataList_List{}, err
	}
	return HoldsVerTwoDataList_List{l}, nil
}

func (s HoldsVerTwoDataList_List) At(i int) HoldsVerTwoDataList {
	return HoldsVerTwoDataList{s.List.Struct(i)}
}
func (s HoldsVerTwoDataList_List) Set(i int, v HoldsVerTwoDataList) error {
	return s.List.SetStruct(i, v.Struct)
}

// HoldsVerTwoDataList_Promise is a wrapper for a HoldsVerTwoDataList promised by a client call.
type HoldsVerTwoDataList_Promise struct{ *capnp.Pipeline }

func (p HoldsVerTwoDataList_Promise) Struct() (HoldsVerTwoDataList, error) {
	s, err := p.Pipeline.Struct()
	return HoldsVerTwoDataList{s}, err
}

type HoldsVerOnePtrList struct{ capnp.Struct }

func NewHoldsVerOnePtrList(s *capnp.Segment) (HoldsVerOnePtrList, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return HoldsVerOnePtrList{}, err
	}
	return HoldsVerOnePtrList{st}, nil
}

func NewRootHoldsVerOnePtrList(s *capnp.Segment) (HoldsVerOnePtrList, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return HoldsVerOnePtrList{}, err
	}
	return HoldsVerOnePtrList{st}, nil
}

func ReadRootHoldsVerOnePtrList(msg *capnp.Message) (HoldsVerOnePtrList, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return HoldsVerOnePtrList{}, err
	}
	return HoldsVerOnePtrList{root.Struct()}, nil
}
func (s HoldsVerOnePtrList) Mylist() (VerOnePtr_List, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return VerOnePtr_List{}, err
	}
	return VerOnePtr_List{List: p.List()}, nil
}

func (s HoldsVerOnePtrList) HasMylist() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s HoldsVerOnePtrList) SetMylist(v VerOnePtr_List) error {
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewMylist sets the mylist field to a newly
// allocated VerOnePtr_List, preferring placement in s's segment.
func (s HoldsVerOnePtrList) NewMylist(n int32) (VerOnePtr_List, error) {
	l, err := NewVerOnePtr_List(s.Struct.Segment(), n)
	if err != nil {
		return VerOnePtr_List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

// HoldsVerOnePtrList_List is a list of HoldsVerOnePtrList.
type HoldsVerOnePtrList_List struct{ capnp.List }

// NewHoldsVerOnePtrList creates a new list of HoldsVerOnePtrList.
func NewHoldsVerOnePtrList_List(s *capnp.Segment, sz int32) (HoldsVerOnePtrList_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return HoldsVerOnePtrList_List{}, err
	}
	return HoldsVerOnePtrList_List{l}, nil
}

func (s HoldsVerOnePtrList_List) At(i int) HoldsVerOnePtrList {
	return HoldsVerOnePtrList{s.List.Struct(i)}
}
func (s HoldsVerOnePtrList_List) Set(i int, v HoldsVerOnePtrList) error {
	return s.List.SetStruct(i, v.Struct)
}

// HoldsVerOnePtrList_Promise is a wrapper for a HoldsVerOnePtrList promised by a client call.
type HoldsVerOnePtrList_Promise struct{ *capnp.Pipeline }

func (p HoldsVerOnePtrList_Promise) Struct() (HoldsVerOnePtrList, error) {
	s, err := p.Pipeline.Struct()
	return HoldsVerOnePtrList{s}, err
}

type HoldsVerTwoPtrList struct{ capnp.Struct }

func NewHoldsVerTwoPtrList(s *capnp.Segment) (HoldsVerTwoPtrList, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return HoldsVerTwoPtrList{}, err
	}
	return HoldsVerTwoPtrList{st}, nil
}

func NewRootHoldsVerTwoPtrList(s *capnp.Segment) (HoldsVerTwoPtrList, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return HoldsVerTwoPtrList{}, err
	}
	return HoldsVerTwoPtrList{st}, nil
}

func ReadRootHoldsVerTwoPtrList(msg *capnp.Message) (HoldsVerTwoPtrList, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return HoldsVerTwoPtrList{}, err
	}
	return HoldsVerTwoPtrList{root.Struct()}, nil
}
func (s HoldsVerTwoPtrList) Mylist() (VerTwoPtr_List, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return VerTwoPtr_List{}, err
	}
	return VerTwoPtr_List{List: p.List()}, nil
}

func (s HoldsVerTwoPtrList) HasMylist() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s HoldsVerTwoPtrList) SetMylist(v VerTwoPtr_List) error {
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewMylist sets the mylist field to a newly
// allocated VerTwoPtr_List, preferring placement in s's segment.
func (s HoldsVerTwoPtrList) NewMylist(n int32) (VerTwoPtr_List, error) {
	l, err := NewVerTwoPtr_List(s.Struct.Segment(), n)
	if err != nil {
		return VerTwoPtr_List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

// HoldsVerTwoPtrList_List is a list of HoldsVerTwoPtrList.
type HoldsVerTwoPtrList_List struct{ capnp.List }

// NewHoldsVerTwoPtrList creates a new list of HoldsVerTwoPtrList.
func NewHoldsVerTwoPtrList_List(s *capnp.Segment, sz int32) (HoldsVerTwoPtrList_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return HoldsVerTwoPtrList_List{}, err
	}
	return HoldsVerTwoPtrList_List{l}, nil
}

func (s HoldsVerTwoPtrList_List) At(i int) HoldsVerTwoPtrList {
	return HoldsVerTwoPtrList{s.List.Struct(i)}
}
func (s HoldsVerTwoPtrList_List) Set(i int, v HoldsVerTwoPtrList) error {
	return s.List.SetStruct(i, v.Struct)
}

// HoldsVerTwoPtrList_Promise is a wrapper for a HoldsVerTwoPtrList promised by a client call.
type HoldsVerTwoPtrList_Promise struct{ *capnp.Pipeline }

func (p HoldsVerTwoPtrList_Promise) Struct() (HoldsVerTwoPtrList, error) {
	s, err := p.Pipeline.Struct()
	return HoldsVerTwoPtrList{s}, err
}

type HoldsVerTwoTwoList struct{ capnp.Struct }

func NewHoldsVerTwoTwoList(s *capnp.Segment) (HoldsVerTwoTwoList, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return HoldsVerTwoTwoList{}, err
	}
	return HoldsVerTwoTwoList{st}, nil
}

func NewRootHoldsVerTwoTwoList(s *capnp.Segment) (HoldsVerTwoTwoList, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return HoldsVerTwoTwoList{}, err
	}
	return HoldsVerTwoTwoList{st}, nil
}

func ReadRootHoldsVerTwoTwoList(msg *capnp.Message) (HoldsVerTwoTwoList, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return HoldsVerTwoTwoList{}, err
	}
	return HoldsVerTwoTwoList{root.Struct()}, nil
}
func (s HoldsVerTwoTwoList) Mylist() (VerTwoDataTwoPtr_List, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return VerTwoDataTwoPtr_List{}, err
	}
	return VerTwoDataTwoPtr_List{List: p.List()}, nil
}

func (s HoldsVerTwoTwoList) HasMylist() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s HoldsVerTwoTwoList) SetMylist(v VerTwoDataTwoPtr_List) error {
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewMylist sets the mylist field to a newly
// allocated VerTwoDataTwoPtr_List, preferring placement in s's segment.
func (s HoldsVerTwoTwoList) NewMylist(n int32) (VerTwoDataTwoPtr_List, error) {
	l, err := NewVerTwoDataTwoPtr_List(s.Struct.Segment(), n)
	if err != nil {
		return VerTwoDataTwoPtr_List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

// HoldsVerTwoTwoList_List is a list of HoldsVerTwoTwoList.
type HoldsVerTwoTwoList_List struct{ capnp.List }

// NewHoldsVerTwoTwoList creates a new list of HoldsVerTwoTwoList.
func NewHoldsVerTwoTwoList_List(s *capnp.Segment, sz int32) (HoldsVerTwoTwoList_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return HoldsVerTwoTwoList_List{}, err
	}
	return HoldsVerTwoTwoList_List{l}, nil
}

func (s HoldsVerTwoTwoList_List) At(i int) HoldsVerTwoTwoList {
	return HoldsVerTwoTwoList{s.List.Struct(i)}
}
func (s HoldsVerTwoTwoList_List) Set(i int, v HoldsVerTwoTwoList) error {
	return s.List.SetStruct(i, v.Struct)
}

// HoldsVerTwoTwoList_Promise is a wrapper for a HoldsVerTwoTwoList promised by a client call.
type HoldsVerTwoTwoList_Promise struct{ *capnp.Pipeline }

func (p HoldsVerTwoTwoList_Promise) Struct() (HoldsVerTwoTwoList, error) {
	s, err := p.Pipeline.Struct()
	return HoldsVerTwoTwoList{s}, err
}

type HoldsVerTwoTwoPlus struct{ capnp.Struct }

func NewHoldsVerTwoTwoPlus(s *capnp.Segment) (HoldsVerTwoTwoPlus, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return HoldsVerTwoTwoPlus{}, err
	}
	return HoldsVerTwoTwoPlus{st}, nil
}

func NewRootHoldsVerTwoTwoPlus(s *capnp.Segment) (HoldsVerTwoTwoPlus, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return HoldsVerTwoTwoPlus{}, err
	}
	return HoldsVerTwoTwoPlus{st}, nil
}

func ReadRootHoldsVerTwoTwoPlus(msg *capnp.Message) (HoldsVerTwoTwoPlus, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return HoldsVerTwoTwoPlus{}, err
	}
	return HoldsVerTwoTwoPlus{root.Struct()}, nil
}
func (s HoldsVerTwoTwoPlus) Mylist() (VerTwoTwoPlus_List, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return VerTwoTwoPlus_List{}, err
	}
	return VerTwoTwoPlus_List{List: p.List()}, nil
}

func (s HoldsVerTwoTwoPlus) HasMylist() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s HoldsVerTwoTwoPlus) SetMylist(v VerTwoTwoPlus_List) error {
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewMylist sets the mylist field to a newly
// allocated VerTwoTwoPlus_List, preferring placement in s's segment.
func (s HoldsVerTwoTwoPlus) NewMylist(n int32) (VerTwoTwoPlus_List, error) {
	l, err := NewVerTwoTwoPlus_List(s.Struct.Segment(), n)
	if err != nil {
		return VerTwoTwoPlus_List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

// HoldsVerTwoTwoPlus_List is a list of HoldsVerTwoTwoPlus.
type HoldsVerTwoTwoPlus_List struct{ capnp.List }

// NewHoldsVerTwoTwoPlus creates a new list of HoldsVerTwoTwoPlus.
func NewHoldsVerTwoTwoPlus_List(s *capnp.Segment, sz int32) (HoldsVerTwoTwoPlus_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return HoldsVerTwoTwoPlus_List{}, err
	}
	return HoldsVerTwoTwoPlus_List{l}, nil
}

func (s HoldsVerTwoTwoPlus_List) At(i int) HoldsVerTwoTwoPlus {
	return HoldsVerTwoTwoPlus{s.List.Struct(i)}
}
func (s HoldsVerTwoTwoPlus_List) Set(i int, v HoldsVerTwoTwoPlus) error {
	return s.List.SetStruct(i, v.Struct)
}

// HoldsVerTwoTwoPlus_Promise is a wrapper for a HoldsVerTwoTwoPlus promised by a client call.
type HoldsVerTwoTwoPlus_Promise struct{ *capnp.Pipeline }

func (p HoldsVerTwoTwoPlus_Promise) Struct() (HoldsVerTwoTwoPlus, error) {
	s, err := p.Pipeline.Struct()
	return HoldsVerTwoTwoPlus{s}, err
}

type VerTwoTwoPlus struct{ capnp.Struct }

func NewVerTwoTwoPlus(s *capnp.Segment) (VerTwoTwoPlus, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 24, PointerCount: 3})
	if err != nil {
		return VerTwoTwoPlus{}, err
	}
	return VerTwoTwoPlus{st}, nil
}

func NewRootVerTwoTwoPlus(s *capnp.Segment) (VerTwoTwoPlus, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 24, PointerCount: 3})
	if err != nil {
		return VerTwoTwoPlus{}, err
	}
	return VerTwoTwoPlus{st}, nil
}

func ReadRootVerTwoTwoPlus(msg *capnp.Message) (VerTwoTwoPlus, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return VerTwoTwoPlus{}, err
	}
	return VerTwoTwoPlus{root.Struct()}, nil
}
func (s VerTwoTwoPlus) Val() int16 {
	return int16(s.Struct.Uint16(0))
}

func (s VerTwoTwoPlus) SetVal(v int16) {
	s.Struct.SetUint16(0, uint16(v))
}

func (s VerTwoTwoPlus) Duo() int64 {
	return int64(s.Struct.Uint64(8))
}

func (s VerTwoTwoPlus) SetDuo(v int64) {
	s.Struct.SetUint64(8, uint64(v))
}

func (s VerTwoTwoPlus) Ptr1() (VerTwoDataTwoPtr, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return VerTwoDataTwoPtr{}, err
	}
	return VerTwoDataTwoPtr{Struct: p.Struct()}, nil
}

func (s VerTwoTwoPlus) HasPtr1() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s VerTwoTwoPlus) SetPtr1(v VerTwoDataTwoPtr) error {
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewPtr1 sets the ptr1 field to a newly
// allocated VerTwoDataTwoPtr struct, preferring placement in s's segment.
func (s VerTwoTwoPlus) NewPtr1() (VerTwoDataTwoPtr, error) {
	ss, err := NewVerTwoDataTwoPtr(s.Struct.Segment())
	if err != nil {
		return VerTwoDataTwoPtr{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s VerTwoTwoPlus) Ptr2() (VerTwoDataTwoPtr, error) {
	p, err := s.Struct.Ptr(1)
	if err != nil {
		return VerTwoDataTwoPtr{}, err
	}
	return VerTwoDataTwoPtr{Struct: p.Struct()}, nil
}

func (s VerTwoTwoPlus) HasPtr2() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s VerTwoTwoPlus) SetPtr2(v VerTwoDataTwoPtr) error {
	return s.Struct.SetPtr(1, v.Struct.ToPtr())
}

// NewPtr2 sets the ptr2 field to a newly
// allocated VerTwoDataTwoPtr struct, preferring placement in s's segment.
func (s VerTwoTwoPlus) NewPtr2() (VerTwoDataTwoPtr, error) {
	ss, err := NewVerTwoDataTwoPtr(s.Struct.Segment())
	if err != nil {
		return VerTwoDataTwoPtr{}, err
	}
	err = s.Struct.SetPtr(1, ss.Struct.ToPtr())
	return ss, err
}

func (s VerTwoTwoPlus) Tre() int64 {
	return int64(s.Struct.Uint64(16))
}

func (s VerTwoTwoPlus) SetTre(v int64) {
	s.Struct.SetUint64(16, uint64(v))
}

func (s VerTwoTwoPlus) Lst3() (capnp.Int64List, error) {
	p, err := s.Struct.Ptr(2)
	if err != nil {
		return capnp.Int64List{}, err
	}
	return capnp.Int64List{List: p.List()}, nil
}

func (s VerTwoTwoPlus) HasLst3() bool {
	p, err := s.Struct.Ptr(2)
	return p.IsValid() || err != nil
}

func (s VerTwoTwoPlus) SetLst3(v capnp.Int64List) error {
	return s.Struct.SetPtr(2, v.List.ToPtr())
}

// NewLst3 sets the lst3 field to a newly
// allocated capnp.Int64List, preferring placement in s's segment.
func (s VerTwoTwoPlus) NewLst3(n int32) (capnp.Int64List, error) {
	l, err := capnp.NewInt64List(s.Struct.Segment(), n)
	if err != nil {
		return capnp.Int64List{}, err
	}
	err = s.Struct.SetPtr(2, l.List.ToPtr())
	return l, err
}

// VerTwoTwoPlus_List is a list of VerTwoTwoPlus.
type VerTwoTwoPlus_List struct{ capnp.List }

// NewVerTwoTwoPlus creates a new list of VerTwoTwoPlus.
func NewVerTwoTwoPlus_List(s *capnp.Segment, sz int32) (VerTwoTwoPlus_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 24, PointerCount: 3}, sz)
	if err != nil {
		return VerTwoTwoPlus_List{}, err
	}
	return VerTwoTwoPlus_List{l}, nil
}

func (s VerTwoTwoPlus_List) At(i int) VerTwoTwoPlus           { return VerTwoTwoPlus{s.List.Struct(i)} }
func (s VerTwoTwoPlus_List) Set(i int, v VerTwoTwoPlus) error { return s.List.SetStruct(i, v.Struct) }

// VerTwoTwoPlus_Promise is a wrapper for a VerTwoTwoPlus promised by a client call.
type VerTwoTwoPlus_Promise struct{ *capnp.Pipeline }

func (p VerTwoTwoPlus_Promise) Struct() (VerTwoTwoPlus, error) {
	s, err := p.Pipeline.Struct()
	return VerTwoTwoPlus{s}, err
}

func (p VerTwoTwoPlus_Promise) Ptr1() VerTwoDataTwoPtr_Promise {
	return VerTwoDataTwoPtr_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p VerTwoTwoPlus_Promise) Ptr2() VerTwoDataTwoPtr_Promise {
	return VerTwoDataTwoPtr_Promise{Pipeline: p.Pipeline.GetPipeline(1)}
}

type HoldsText struct{ capnp.Struct }

func NewHoldsText(s *capnp.Segment) (HoldsText, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 3})
	if err != nil {
		return HoldsText{}, err
	}
	return HoldsText{st}, nil
}

func NewRootHoldsText(s *capnp.Segment) (HoldsText, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 3})
	if err != nil {
		return HoldsText{}, err
	}
	return HoldsText{st}, nil
}

func ReadRootHoldsText(msg *capnp.Message) (HoldsText, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return HoldsText{}, err
	}
	return HoldsText{root.Struct()}, nil
}
func (s HoldsText) Txt() (string, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return "", err
	}
	return p.Text(), nil
}

func (s HoldsText) HasTxt() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s HoldsText) TxtBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	d := p.Data()
	if len(d) == 0 {
		return d, nil
	}
	return d[:len(d)-1], nil
}

func (s HoldsText) SetTxt(v string) error {
	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, t.List.ToPtr())
}

func (s HoldsText) Lst() (capnp.TextList, error) {
	p, err := s.Struct.Ptr(1)
	if err != nil {
		return capnp.TextList{}, err
	}
	return capnp.TextList{List: p.List()}, nil
}

func (s HoldsText) HasLst() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s HoldsText) SetLst(v capnp.TextList) error {
	return s.Struct.SetPtr(1, v.List.ToPtr())
}

// NewLst sets the lst field to a newly
// allocated capnp.TextList, preferring placement in s's segment.
func (s HoldsText) NewLst(n int32) (capnp.TextList, error) {
	l, err := capnp.NewTextList(s.Struct.Segment(), n)
	if err != nil {
		return capnp.TextList{}, err
	}
	err = s.Struct.SetPtr(1, l.List.ToPtr())
	return l, err
}

func (s HoldsText) Lstlst() (capnp.PointerList, error) {
	p, err := s.Struct.Ptr(2)
	if err != nil {
		return capnp.PointerList{}, err
	}
	return capnp.PointerList{List: p.List()}, nil
}

func (s HoldsText) HasLstlst() bool {
	p, err := s.Struct.Ptr(2)
	return p.IsValid() || err != nil
}

func (s HoldsText) SetLstlst(v capnp.PointerList) error {
	return s.Struct.SetPtr(2, v.List.ToPtr())
}

// NewLstlst sets the lstlst field to a newly
// allocated capnp.PointerList, preferring placement in s's segment.
func (s HoldsText) NewLstlst(n int32) (capnp.PointerList, error) {
	l, err := capnp.NewPointerList(s.Struct.Segment(), n)
	if err != nil {
		return capnp.PointerList{}, err
	}
	err = s.Struct.SetPtr(2, l.List.ToPtr())
	return l, err
}

// HoldsText_List is a list of HoldsText.
type HoldsText_List struct{ capnp.List }

// NewHoldsText creates a new list of HoldsText.
func NewHoldsText_List(s *capnp.Segment, sz int32) (HoldsText_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 3}, sz)
	if err != nil {
		return HoldsText_List{}, err
	}
	return HoldsText_List{l}, nil
}

func (s HoldsText_List) At(i int) HoldsText           { return HoldsText{s.List.Struct(i)} }
func (s HoldsText_List) Set(i int, v HoldsText) error { return s.List.SetStruct(i, v.Struct) }

// HoldsText_Promise is a wrapper for a HoldsText promised by a client call.
type HoldsText_Promise struct{ *capnp.Pipeline }

func (p HoldsText_Promise) Struct() (HoldsText, error) {
	s, err := p.Pipeline.Struct()
	return HoldsText{s}, err
}

type WrapEmpty struct{ capnp.Struct }

func NewWrapEmpty(s *capnp.Segment) (WrapEmpty, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return WrapEmpty{}, err
	}
	return WrapEmpty{st}, nil
}

func NewRootWrapEmpty(s *capnp.Segment) (WrapEmpty, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return WrapEmpty{}, err
	}
	return WrapEmpty{st}, nil
}

func ReadRootWrapEmpty(msg *capnp.Message) (WrapEmpty, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return WrapEmpty{}, err
	}
	return WrapEmpty{root.Struct()}, nil
}
func (s WrapEmpty) MightNotBeReallyEmpty() (VerEmpty, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return VerEmpty{}, err
	}
	return VerEmpty{Struct: p.Struct()}, nil
}

func (s WrapEmpty) HasMightNotBeReallyEmpty() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s WrapEmpty) SetMightNotBeReallyEmpty(v VerEmpty) error {
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewMightNotBeReallyEmpty sets the mightNotBeReallyEmpty field to a newly
// allocated VerEmpty struct, preferring placement in s's segment.
func (s WrapEmpty) NewMightNotBeReallyEmpty() (VerEmpty, error) {
	ss, err := NewVerEmpty(s.Struct.Segment())
	if err != nil {
		return VerEmpty{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

// WrapEmpty_List is a list of WrapEmpty.
type WrapEmpty_List struct{ capnp.List }

// NewWrapEmpty creates a new list of WrapEmpty.
func NewWrapEmpty_List(s *capnp.Segment, sz int32) (WrapEmpty_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return WrapEmpty_List{}, err
	}
	return WrapEmpty_List{l}, nil
}

func (s WrapEmpty_List) At(i int) WrapEmpty           { return WrapEmpty{s.List.Struct(i)} }
func (s WrapEmpty_List) Set(i int, v WrapEmpty) error { return s.List.SetStruct(i, v.Struct) }

// WrapEmpty_Promise is a wrapper for a WrapEmpty promised by a client call.
type WrapEmpty_Promise struct{ *capnp.Pipeline }

func (p WrapEmpty_Promise) Struct() (WrapEmpty, error) {
	s, err := p.Pipeline.Struct()
	return WrapEmpty{s}, err
}

func (p WrapEmpty_Promise) MightNotBeReallyEmpty() VerEmpty_Promise {
	return VerEmpty_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

type Wrap2x2 struct{ capnp.Struct }

func NewWrap2x2(s *capnp.Segment) (Wrap2x2, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Wrap2x2{}, err
	}
	return Wrap2x2{st}, nil
}

func NewRootWrap2x2(s *capnp.Segment) (Wrap2x2, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Wrap2x2{}, err
	}
	return Wrap2x2{st}, nil
}

func ReadRootWrap2x2(msg *capnp.Message) (Wrap2x2, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Wrap2x2{}, err
	}
	return Wrap2x2{root.Struct()}, nil
}
func (s Wrap2x2) MightNotBeReallyEmpty() (VerTwoDataTwoPtr, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return VerTwoDataTwoPtr{}, err
	}
	return VerTwoDataTwoPtr{Struct: p.Struct()}, nil
}

func (s Wrap2x2) HasMightNotBeReallyEmpty() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Wrap2x2) SetMightNotBeReallyEmpty(v VerTwoDataTwoPtr) error {
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewMightNotBeReallyEmpty sets the mightNotBeReallyEmpty field to a newly
// allocated VerTwoDataTwoPtr struct, preferring placement in s's segment.
func (s Wrap2x2) NewMightNotBeReallyEmpty() (VerTwoDataTwoPtr, error) {
	ss, err := NewVerTwoDataTwoPtr(s.Struct.Segment())
	if err != nil {
		return VerTwoDataTwoPtr{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

// Wrap2x2_List is a list of Wrap2x2.
type Wrap2x2_List struct{ capnp.List }

// NewWrap2x2 creates a new list of Wrap2x2.
func NewWrap2x2_List(s *capnp.Segment, sz int32) (Wrap2x2_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return Wrap2x2_List{}, err
	}
	return Wrap2x2_List{l}, nil
}

func (s Wrap2x2_List) At(i int) Wrap2x2           { return Wrap2x2{s.List.Struct(i)} }
func (s Wrap2x2_List) Set(i int, v Wrap2x2) error { return s.List.SetStruct(i, v.Struct) }

// Wrap2x2_Promise is a wrapper for a Wrap2x2 promised by a client call.
type Wrap2x2_Promise struct{ *capnp.Pipeline }

func (p Wrap2x2_Promise) Struct() (Wrap2x2, error) {
	s, err := p.Pipeline.Struct()
	return Wrap2x2{s}, err
}

func (p Wrap2x2_Promise) MightNotBeReallyEmpty() VerTwoDataTwoPtr_Promise {
	return VerTwoDataTwoPtr_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

type Wrap2x2plus struct{ capnp.Struct }

func NewWrap2x2plus(s *capnp.Segment) (Wrap2x2plus, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Wrap2x2plus{}, err
	}
	return Wrap2x2plus{st}, nil
}

func NewRootWrap2x2plus(s *capnp.Segment) (Wrap2x2plus, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Wrap2x2plus{}, err
	}
	return Wrap2x2plus{st}, nil
}

func ReadRootWrap2x2plus(msg *capnp.Message) (Wrap2x2plus, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Wrap2x2plus{}, err
	}
	return Wrap2x2plus{root.Struct()}, nil
}
func (s Wrap2x2plus) MightNotBeReallyEmpty() (VerTwoTwoPlus, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return VerTwoTwoPlus{}, err
	}
	return VerTwoTwoPlus{Struct: p.Struct()}, nil
}

func (s Wrap2x2plus) HasMightNotBeReallyEmpty() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Wrap2x2plus) SetMightNotBeReallyEmpty(v VerTwoTwoPlus) error {
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewMightNotBeReallyEmpty sets the mightNotBeReallyEmpty field to a newly
// allocated VerTwoTwoPlus struct, preferring placement in s's segment.
func (s Wrap2x2plus) NewMightNotBeReallyEmpty() (VerTwoTwoPlus, error) {
	ss, err := NewVerTwoTwoPlus(s.Struct.Segment())
	if err != nil {
		return VerTwoTwoPlus{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

// Wrap2x2plus_List is a list of Wrap2x2plus.
type Wrap2x2plus_List struct{ capnp.List }

// NewWrap2x2plus creates a new list of Wrap2x2plus.
func NewWrap2x2plus_List(s *capnp.Segment, sz int32) (Wrap2x2plus_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return Wrap2x2plus_List{}, err
	}
	return Wrap2x2plus_List{l}, nil
}

func (s Wrap2x2plus_List) At(i int) Wrap2x2plus           { return Wrap2x2plus{s.List.Struct(i)} }
func (s Wrap2x2plus_List) Set(i int, v Wrap2x2plus) error { return s.List.SetStruct(i, v.Struct) }

// Wrap2x2plus_Promise is a wrapper for a Wrap2x2plus promised by a client call.
type Wrap2x2plus_Promise struct{ *capnp.Pipeline }

func (p Wrap2x2plus_Promise) Struct() (Wrap2x2plus, error) {
	s, err := p.Pipeline.Struct()
	return Wrap2x2plus{s}, err
}

func (p Wrap2x2plus_Promise) MightNotBeReallyEmpty() VerTwoTwoPlus_Promise {
	return VerTwoTwoPlus_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

type VoidUnion struct{ capnp.Struct }
type VoidUnion_Which uint16

const (
	VoidUnion_Which_a VoidUnion_Which = 0
	VoidUnion_Which_b VoidUnion_Which = 1
)

func (w VoidUnion_Which) String() string {
	const s = "ab"
	switch w {
	case VoidUnion_Which_a:
		return s[0:1]
	case VoidUnion_Which_b:
		return s[1:2]

	}
	return "VoidUnion_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
}

func NewVoidUnion(s *capnp.Segment) (VoidUnion, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	if err != nil {
		return VoidUnion{}, err
	}
	return VoidUnion{st}, nil
}

func NewRootVoidUnion(s *capnp.Segment) (VoidUnion, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	if err != nil {
		return VoidUnion{}, err
	}
	return VoidUnion{st}, nil
}

func ReadRootVoidUnion(msg *capnp.Message) (VoidUnion, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return VoidUnion{}, err
	}
	return VoidUnion{root.Struct()}, nil
}

func (s VoidUnion) Which() VoidUnion_Which {
	return VoidUnion_Which(s.Struct.Uint16(0))
}
func (s VoidUnion) SetA() {
	s.Struct.SetUint16(0, 0)

}

func (s VoidUnion) SetB() {
	s.Struct.SetUint16(0, 1)

}

// VoidUnion_List is a list of VoidUnion.
type VoidUnion_List struct{ capnp.List }

// NewVoidUnion creates a new list of VoidUnion.
func NewVoidUnion_List(s *capnp.Segment, sz int32) (VoidUnion_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0}, sz)
	if err != nil {
		return VoidUnion_List{}, err
	}
	return VoidUnion_List{l}, nil
}

func (s VoidUnion_List) At(i int) VoidUnion           { return VoidUnion{s.List.Struct(i)} }
func (s VoidUnion_List) Set(i int, v VoidUnion) error { return s.List.SetStruct(i, v.Struct) }

// VoidUnion_Promise is a wrapper for a VoidUnion promised by a client call.
type VoidUnion_Promise struct{ *capnp.Pipeline }

func (p VoidUnion_Promise) Struct() (VoidUnion, error) {
	s, err := p.Pipeline.Struct()
	return VoidUnion{s}, err
}

type Nester1Capn struct{ capnp.Struct }

func NewNester1Capn(s *capnp.Segment) (Nester1Capn, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Nester1Capn{}, err
	}
	return Nester1Capn{st}, nil
}

func NewRootNester1Capn(s *capnp.Segment) (Nester1Capn, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Nester1Capn{}, err
	}
	return Nester1Capn{st}, nil
}

func ReadRootNester1Capn(msg *capnp.Message) (Nester1Capn, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Nester1Capn{}, err
	}
	return Nester1Capn{root.Struct()}, nil
}
func (s Nester1Capn) Strs() (capnp.TextList, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return capnp.TextList{}, err
	}
	return capnp.TextList{List: p.List()}, nil
}

func (s Nester1Capn) HasStrs() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Nester1Capn) SetStrs(v capnp.TextList) error {
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewStrs sets the strs field to a newly
// allocated capnp.TextList, preferring placement in s's segment.
func (s Nester1Capn) NewStrs(n int32) (capnp.TextList, error) {
	l, err := capnp.NewTextList(s.Struct.Segment(), n)
	if err != nil {
		return capnp.TextList{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

// Nester1Capn_List is a list of Nester1Capn.
type Nester1Capn_List struct{ capnp.List }

// NewNester1Capn creates a new list of Nester1Capn.
func NewNester1Capn_List(s *capnp.Segment, sz int32) (Nester1Capn_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return Nester1Capn_List{}, err
	}
	return Nester1Capn_List{l}, nil
}

func (s Nester1Capn_List) At(i int) Nester1Capn           { return Nester1Capn{s.List.Struct(i)} }
func (s Nester1Capn_List) Set(i int, v Nester1Capn) error { return s.List.SetStruct(i, v.Struct) }

// Nester1Capn_Promise is a wrapper for a Nester1Capn promised by a client call.
type Nester1Capn_Promise struct{ *capnp.Pipeline }

func (p Nester1Capn_Promise) Struct() (Nester1Capn, error) {
	s, err := p.Pipeline.Struct()
	return Nester1Capn{s}, err
}

type RWTestCapn struct{ capnp.Struct }

func NewRWTestCapn(s *capnp.Segment) (RWTestCapn, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return RWTestCapn{}, err
	}
	return RWTestCapn{st}, nil
}

func NewRootRWTestCapn(s *capnp.Segment) (RWTestCapn, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return RWTestCapn{}, err
	}
	return RWTestCapn{st}, nil
}

func ReadRootRWTestCapn(msg *capnp.Message) (RWTestCapn, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return RWTestCapn{}, err
	}
	return RWTestCapn{root.Struct()}, nil
}
func (s RWTestCapn) NestMatrix() (capnp.PointerList, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return capnp.PointerList{}, err
	}
	return capnp.PointerList{List: p.List()}, nil
}

func (s RWTestCapn) HasNestMatrix() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s RWTestCapn) SetNestMatrix(v capnp.PointerList) error {
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewNestMatrix sets the nestMatrix field to a newly
// allocated capnp.PointerList, preferring placement in s's segment.
func (s RWTestCapn) NewNestMatrix(n int32) (capnp.PointerList, error) {
	l, err := capnp.NewPointerList(s.Struct.Segment(), n)
	if err != nil {
		return capnp.PointerList{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

// RWTestCapn_List is a list of RWTestCapn.
type RWTestCapn_List struct{ capnp.List }

// NewRWTestCapn creates a new list of RWTestCapn.
func NewRWTestCapn_List(s *capnp.Segment, sz int32) (RWTestCapn_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return RWTestCapn_List{}, err
	}
	return RWTestCapn_List{l}, nil
}

func (s RWTestCapn_List) At(i int) RWTestCapn           { return RWTestCapn{s.List.Struct(i)} }
func (s RWTestCapn_List) Set(i int, v RWTestCapn) error { return s.List.SetStruct(i, v.Struct) }

// RWTestCapn_Promise is a wrapper for a RWTestCapn promised by a client call.
type RWTestCapn_Promise struct{ *capnp.Pipeline }

func (p RWTestCapn_Promise) Struct() (RWTestCapn, error) {
	s, err := p.Pipeline.Struct()
	return RWTestCapn{s}, err
}

type ListStructCapn struct{ capnp.Struct }

func NewListStructCapn(s *capnp.Segment) (ListStructCapn, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return ListStructCapn{}, err
	}
	return ListStructCapn{st}, nil
}

func NewRootListStructCapn(s *capnp.Segment) (ListStructCapn, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return ListStructCapn{}, err
	}
	return ListStructCapn{st}, nil
}

func ReadRootListStructCapn(msg *capnp.Message) (ListStructCapn, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return ListStructCapn{}, err
	}
	return ListStructCapn{root.Struct()}, nil
}
func (s ListStructCapn) Vec() (Nester1Capn_List, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return Nester1Capn_List{}, err
	}
	return Nester1Capn_List{List: p.List()}, nil
}

func (s ListStructCapn) HasVec() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s ListStructCapn) SetVec(v Nester1Capn_List) error {
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewVec sets the vec field to a newly
// allocated Nester1Capn_List, preferring placement in s's segment.
func (s ListStructCapn) NewVec(n int32) (Nester1Capn_List, error) {
	l, err := NewNester1Capn_List(s.Struct.Segment(), n)
	if err != nil {
		return Nester1Capn_List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

// ListStructCapn_List is a list of ListStructCapn.
type ListStructCapn_List struct{ capnp.List }

// NewListStructCapn creates a new list of ListStructCapn.
func NewListStructCapn_List(s *capnp.Segment, sz int32) (ListStructCapn_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return ListStructCapn_List{}, err
	}
	return ListStructCapn_List{l}, nil
}

func (s ListStructCapn_List) At(i int) ListStructCapn           { return ListStructCapn{s.List.Struct(i)} }
func (s ListStructCapn_List) Set(i int, v ListStructCapn) error { return s.List.SetStruct(i, v.Struct) }

// ListStructCapn_Promise is a wrapper for a ListStructCapn promised by a client call.
type ListStructCapn_Promise struct{ *capnp.Pipeline }

func (p ListStructCapn_Promise) Struct() (ListStructCapn, error) {
	s, err := p.Pipeline.Struct()
	return ListStructCapn{s}, err
}

type Echo struct{ Client capnp.Client }

func (c Echo) Echo(ctx context.Context, params func(Echo_echo_Params) error, opts ...capnp.CallOption) Echo_echo_Results_Promise {
	if c.Client == nil {
		return Echo_echo_Results_Promise{Pipeline: capnp.NewPipeline(capnp.ErrorAnswer(capnp.ErrNullClient))}
	}
	call := &capnp.Call{
		Ctx: ctx,
		Method: capnp.Method{
			InterfaceID:   0x8e5322c1e9282534,
			MethodID:      0,
			InterfaceName: "aircraft.capnp:Echo",
			MethodName:    "echo",
		},
		Options: capnp.NewCallOptions(opts),
	}
	if params != nil {
		call.ParamsSize = capnp.ObjectSize{DataSize: 0, PointerCount: 1}
		call.ParamsFunc = func(s capnp.Struct) error { return params(Echo_echo_Params{Struct: s}) }
	}
	return Echo_echo_Results_Promise{Pipeline: capnp.NewPipeline(c.Client.Call(call))}
}

type Echo_Server interface {
	Echo(Echo_echo) error
}

func Echo_ServerToClient(s Echo_Server) Echo {
	c, _ := s.(server.Closer)
	return Echo{Client: server.New(Echo_Methods(nil, s), c)}
}

func Echo_Methods(methods []server.Method, s Echo_Server) []server.Method {
	if cap(methods) == 0 {
		methods = make([]server.Method, 0, 1)
	}

	methods = append(methods, server.Method{
		Method: capnp.Method{
			InterfaceID:   0x8e5322c1e9282534,
			MethodID:      0,
			InterfaceName: "aircraft.capnp:Echo",
			MethodName:    "echo",
		},
		Impl: func(c context.Context, opts capnp.CallOptions, p, r capnp.Struct) error {
			call := Echo_echo{c, opts, Echo_echo_Params{Struct: p}, Echo_echo_Results{Struct: r}}
			return s.Echo(call)
		},
		ResultsSize: capnp.ObjectSize{DataSize: 0, PointerCount: 1},
	})

	return methods
}

// Echo_echo holds the arguments for a server call to Echo.echo.
type Echo_echo struct {
	Ctx     context.Context
	Options capnp.CallOptions
	Params  Echo_echo_Params
	Results Echo_echo_Results
}

type Echo_echo_Params struct{ capnp.Struct }

func NewEcho_echo_Params(s *capnp.Segment) (Echo_echo_Params, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Echo_echo_Params{}, err
	}
	return Echo_echo_Params{st}, nil
}

func NewRootEcho_echo_Params(s *capnp.Segment) (Echo_echo_Params, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Echo_echo_Params{}, err
	}
	return Echo_echo_Params{st}, nil
}

func ReadRootEcho_echo_Params(msg *capnp.Message) (Echo_echo_Params, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Echo_echo_Params{}, err
	}
	return Echo_echo_Params{root.Struct()}, nil
}
func (s Echo_echo_Params) In() (string, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return "", err
	}
	return p.Text(), nil
}

func (s Echo_echo_Params) HasIn() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Echo_echo_Params) InBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	d := p.Data()
	if len(d) == 0 {
		return d, nil
	}
	return d[:len(d)-1], nil
}

func (s Echo_echo_Params) SetIn(v string) error {
	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, t.List.ToPtr())
}

// Echo_echo_Params_List is a list of Echo_echo_Params.
type Echo_echo_Params_List struct{ capnp.List }

// NewEcho_echo_Params creates a new list of Echo_echo_Params.
func NewEcho_echo_Params_List(s *capnp.Segment, sz int32) (Echo_echo_Params_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return Echo_echo_Params_List{}, err
	}
	return Echo_echo_Params_List{l}, nil
}

func (s Echo_echo_Params_List) At(i int) Echo_echo_Params { return Echo_echo_Params{s.List.Struct(i)} }
func (s Echo_echo_Params_List) Set(i int, v Echo_echo_Params) error {
	return s.List.SetStruct(i, v.Struct)
}

// Echo_echo_Params_Promise is a wrapper for a Echo_echo_Params promised by a client call.
type Echo_echo_Params_Promise struct{ *capnp.Pipeline }

func (p Echo_echo_Params_Promise) Struct() (Echo_echo_Params, error) {
	s, err := p.Pipeline.Struct()
	return Echo_echo_Params{s}, err
}

type Echo_echo_Results struct{ capnp.Struct }

func NewEcho_echo_Results(s *capnp.Segment) (Echo_echo_Results, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Echo_echo_Results{}, err
	}
	return Echo_echo_Results{st}, nil
}

func NewRootEcho_echo_Results(s *capnp.Segment) (Echo_echo_Results, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Echo_echo_Results{}, err
	}
	return Echo_echo_Results{st}, nil
}

func ReadRootEcho_echo_Results(msg *capnp.Message) (Echo_echo_Results, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Echo_echo_Results{}, err
	}
	return Echo_echo_Results{root.Struct()}, nil
}
func (s Echo_echo_Results) Out() (string, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return "", err
	}
	return p.Text(), nil
}

func (s Echo_echo_Results) HasOut() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Echo_echo_Results) OutBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	d := p.Data()
	if len(d) == 0 {
		return d, nil
	}
	return d[:len(d)-1], nil
}

func (s Echo_echo_Results) SetOut(v string) error {
	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, t.List.ToPtr())
}

// Echo_echo_Results_List is a list of Echo_echo_Results.
type Echo_echo_Results_List struct{ capnp.List }

// NewEcho_echo_Results creates a new list of Echo_echo_Results.
func NewEcho_echo_Results_List(s *capnp.Segment, sz int32) (Echo_echo_Results_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return Echo_echo_Results_List{}, err
	}
	return Echo_echo_Results_List{l}, nil
}

func (s Echo_echo_Results_List) At(i int) Echo_echo_Results {
	return Echo_echo_Results{s.List.Struct(i)}
}
func (s Echo_echo_Results_List) Set(i int, v Echo_echo_Results) error {
	return s.List.SetStruct(i, v.Struct)
}

// Echo_echo_Results_Promise is a wrapper for a Echo_echo_Results promised by a client call.
type Echo_echo_Results_Promise struct{ *capnp.Pipeline }

func (p Echo_echo_Results_Promise) Struct() (Echo_echo_Results, error) {
	s, err := p.Pipeline.Struct()
	return Echo_echo_Results{s}, err
}

type Hoth struct{ capnp.Struct }

func NewHoth(s *capnp.Segment) (Hoth, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Hoth{}, err
	}
	return Hoth{st}, nil
}

func NewRootHoth(s *capnp.Segment) (Hoth, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Hoth{}, err
	}
	return Hoth{st}, nil
}

func ReadRootHoth(msg *capnp.Message) (Hoth, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Hoth{}, err
	}
	return Hoth{root.Struct()}, nil
}
func (s Hoth) Base() (EchoBase, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return EchoBase{}, err
	}
	return EchoBase{Struct: p.Struct()}, nil
}

func (s Hoth) HasBase() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Hoth) SetBase(v EchoBase) error {
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewBase sets the base field to a newly
// allocated EchoBase struct, preferring placement in s's segment.
func (s Hoth) NewBase() (EchoBase, error) {
	ss, err := NewEchoBase(s.Struct.Segment())
	if err != nil {
		return EchoBase{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

// Hoth_List is a list of Hoth.
type Hoth_List struct{ capnp.List }

// NewHoth creates a new list of Hoth.
func NewHoth_List(s *capnp.Segment, sz int32) (Hoth_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return Hoth_List{}, err
	}
	return Hoth_List{l}, nil
}

func (s Hoth_List) At(i int) Hoth           { return Hoth{s.List.Struct(i)} }
func (s Hoth_List) Set(i int, v Hoth) error { return s.List.SetStruct(i, v.Struct) }

// Hoth_Promise is a wrapper for a Hoth promised by a client call.
type Hoth_Promise struct{ *capnp.Pipeline }

func (p Hoth_Promise) Struct() (Hoth, error) {
	s, err := p.Pipeline.Struct()
	return Hoth{s}, err
}

func (p Hoth_Promise) Base() EchoBase_Promise {
	return EchoBase_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

type EchoBase struct{ capnp.Struct }

func NewEchoBase(s *capnp.Segment) (EchoBase, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return EchoBase{}, err
	}
	return EchoBase{st}, nil
}

func NewRootEchoBase(s *capnp.Segment) (EchoBase, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return EchoBase{}, err
	}
	return EchoBase{st}, nil
}

func ReadRootEchoBase(msg *capnp.Message) (EchoBase, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return EchoBase{}, err
	}
	return EchoBase{root.Struct()}, nil
}
func (s EchoBase) Echo() Echo {
	p, err := s.Struct.Ptr(0)
	if err != nil {

		return Echo{}
	}
	return Echo{Client: p.Interface().Client()}
}

func (s EchoBase) HasEcho() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s EchoBase) SetEcho(v Echo) error {
	seg := s.Segment()
	if seg == nil {

		return nil
	}
	var in capnp.Interface
	if v.Client != nil {
		in = capnp.NewInterface(seg, seg.Message().AddCap(v.Client))
	}
	return s.Struct.SetPtr(0, in.ToPtr())
}

// EchoBase_List is a list of EchoBase.
type EchoBase_List struct{ capnp.List }

// NewEchoBase creates a new list of EchoBase.
func NewEchoBase_List(s *capnp.Segment, sz int32) (EchoBase_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return EchoBase_List{}, err
	}
	return EchoBase_List{l}, nil
}

func (s EchoBase_List) At(i int) EchoBase           { return EchoBase{s.List.Struct(i)} }
func (s EchoBase_List) Set(i int, v EchoBase) error { return s.List.SetStruct(i, v.Struct) }

// EchoBase_Promise is a wrapper for a EchoBase promised by a client call.
type EchoBase_Promise struct{ *capnp.Pipeline }

func (p EchoBase_Promise) Struct() (EchoBase, error) {
	s, err := p.Pipeline.Struct()
	return EchoBase{s}, err
}

func (p EchoBase_Promise) Echo() Echo {
	return Echo{Client: p.Pipeline.GetPipeline(0).Client()}
}

type StackingRoot struct{ capnp.Struct }

func NewStackingRoot(s *capnp.Segment) (StackingRoot, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	if err != nil {
		return StackingRoot{}, err
	}
	return StackingRoot{st}, nil
}

func NewRootStackingRoot(s *capnp.Segment) (StackingRoot, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	if err != nil {
		return StackingRoot{}, err
	}
	return StackingRoot{st}, nil
}

func ReadRootStackingRoot(msg *capnp.Message) (StackingRoot, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return StackingRoot{}, err
	}
	return StackingRoot{root.Struct()}, nil
}
func (s StackingRoot) A() (StackingA, error) {
	p, err := s.Struct.Ptr(1)
	if err != nil {
		return StackingA{}, err
	}
	return StackingA{Struct: p.Struct()}, nil
}

func (s StackingRoot) HasA() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s StackingRoot) SetA(v StackingA) error {
	return s.Struct.SetPtr(1, v.Struct.ToPtr())
}

// NewA sets the a field to a newly
// allocated StackingA struct, preferring placement in s's segment.
func (s StackingRoot) NewA() (StackingA, error) {
	ss, err := NewStackingA(s.Struct.Segment())
	if err != nil {
		return StackingA{}, err
	}
	err = s.Struct.SetPtr(1, ss.Struct.ToPtr())
	return ss, err
}

func (s StackingRoot) AWithDefault() (StackingA, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return StackingA{}, err
	}
	ss, err := p.StructDefault(x_832bcc6686a26d56[64:96])
	if err != nil {
		return StackingA{}, err
	}
	return StackingA{Struct: ss}, nil
}

func (s StackingRoot) HasAWithDefault() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s StackingRoot) SetAWithDefault(v StackingA) error {
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewAWithDefault sets the aWithDefault field to a newly
// allocated StackingA struct, preferring placement in s's segment.
func (s StackingRoot) NewAWithDefault() (StackingA, error) {
	ss, err := NewStackingA(s.Struct.Segment())
	if err != nil {
		return StackingA{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

// StackingRoot_List is a list of StackingRoot.
type StackingRoot_List struct{ capnp.List }

// NewStackingRoot creates a new list of StackingRoot.
func NewStackingRoot_List(s *capnp.Segment, sz int32) (StackingRoot_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2}, sz)
	if err != nil {
		return StackingRoot_List{}, err
	}
	return StackingRoot_List{l}, nil
}

func (s StackingRoot_List) At(i int) StackingRoot           { return StackingRoot{s.List.Struct(i)} }
func (s StackingRoot_List) Set(i int, v StackingRoot) error { return s.List.SetStruct(i, v.Struct) }

// StackingRoot_Promise is a wrapper for a StackingRoot promised by a client call.
type StackingRoot_Promise struct{ *capnp.Pipeline }

func (p StackingRoot_Promise) Struct() (StackingRoot, error) {
	s, err := p.Pipeline.Struct()
	return StackingRoot{s}, err
}

func (p StackingRoot_Promise) A() StackingA_Promise {
	return StackingA_Promise{Pipeline: p.Pipeline.GetPipeline(1)}
}

func (p StackingRoot_Promise) AWithDefault() StackingA_Promise {
	return StackingA_Promise{Pipeline: p.Pipeline.GetPipelineDefault(0, x_832bcc6686a26d56[96:128])}
}

type StackingA struct{ capnp.Struct }

func NewStackingA(s *capnp.Segment) (StackingA, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		return StackingA{}, err
	}
	return StackingA{st}, nil
}

func NewRootStackingA(s *capnp.Segment) (StackingA, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		return StackingA{}, err
	}
	return StackingA{st}, nil
}

func ReadRootStackingA(msg *capnp.Message) (StackingA, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return StackingA{}, err
	}
	return StackingA{root.Struct()}, nil
}
func (s StackingA) Num() int32 {
	return int32(s.Struct.Uint32(0))
}

func (s StackingA) SetNum(v int32) {
	s.Struct.SetUint32(0, uint32(v))
}

func (s StackingA) B() (StackingB, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return StackingB{}, err
	}
	return StackingB{Struct: p.Struct()}, nil
}

func (s StackingA) HasB() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s StackingA) SetB(v StackingB) error {
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewB sets the b field to a newly
// allocated StackingB struct, preferring placement in s's segment.
func (s StackingA) NewB() (StackingB, error) {
	ss, err := NewStackingB(s.Struct.Segment())
	if err != nil {
		return StackingB{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

// StackingA_List is a list of StackingA.
type StackingA_List struct{ capnp.List }

// NewStackingA creates a new list of StackingA.
func NewStackingA_List(s *capnp.Segment, sz int32) (StackingA_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1}, sz)
	if err != nil {
		return StackingA_List{}, err
	}
	return StackingA_List{l}, nil
}

func (s StackingA_List) At(i int) StackingA           { return StackingA{s.List.Struct(i)} }
func (s StackingA_List) Set(i int, v StackingA) error { return s.List.SetStruct(i, v.Struct) }

// StackingA_Promise is a wrapper for a StackingA promised by a client call.
type StackingA_Promise struct{ *capnp.Pipeline }

func (p StackingA_Promise) Struct() (StackingA, error) {
	s, err := p.Pipeline.Struct()
	return StackingA{s}, err
}

func (p StackingA_Promise) B() StackingB_Promise {
	return StackingB_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

type StackingB struct{ capnp.Struct }

func NewStackingB(s *capnp.Segment) (StackingB, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	if err != nil {
		return StackingB{}, err
	}
	return StackingB{st}, nil
}

func NewRootStackingB(s *capnp.Segment) (StackingB, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	if err != nil {
		return StackingB{}, err
	}
	return StackingB{st}, nil
}

func ReadRootStackingB(msg *capnp.Message) (StackingB, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return StackingB{}, err
	}
	return StackingB{root.Struct()}, nil
}
func (s StackingB) Num() int32 {
	return int32(s.Struct.Uint32(0))
}

func (s StackingB) SetNum(v int32) {
	s.Struct.SetUint32(0, uint32(v))
}

// StackingB_List is a list of StackingB.
type StackingB_List struct{ capnp.List }

// NewStackingB creates a new list of StackingB.
func NewStackingB_List(s *capnp.Segment, sz int32) (StackingB_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0}, sz)
	if err != nil {
		return StackingB_List{}, err
	}
	return StackingB_List{l}, nil
}

func (s StackingB_List) At(i int) StackingB           { return StackingB{s.List.Struct(i)} }
func (s StackingB_List) Set(i int, v StackingB) error { return s.List.SetStruct(i, v.Struct) }

// StackingB_Promise is a wrapper for a StackingB promised by a client call.
type StackingB_Promise struct{ *capnp.Pipeline }

func (p StackingB_Promise) Struct() (StackingB, error) {
	s, err := p.Pipeline.Struct()
	return StackingB{s}, err
}

type CallSequence struct{ Client capnp.Client }

func (c CallSequence) GetNumber(ctx context.Context, params func(CallSequence_getNumber_Params) error, opts ...capnp.CallOption) CallSequence_getNumber_Results_Promise {
	if c.Client == nil {
		return CallSequence_getNumber_Results_Promise{Pipeline: capnp.NewPipeline(capnp.ErrorAnswer(capnp.ErrNullClient))}
	}
	call := &capnp.Call{
		Ctx: ctx,
		Method: capnp.Method{
			InterfaceID:   0xabaedf5f7817c820,
			MethodID:      0,
			InterfaceName: "aircraft.capnp:CallSequence",
			MethodName:    "getNumber",
		},
		Options: capnp.NewCallOptions(opts),
	}
	if params != nil {
		call.ParamsSize = capnp.ObjectSize{DataSize: 0, PointerCount: 0}
		call.ParamsFunc = func(s capnp.Struct) error { return params(CallSequence_getNumber_Params{Struct: s}) }
	}
	return CallSequence_getNumber_Results_Promise{Pipeline: capnp.NewPipeline(c.Client.Call(call))}
}

type CallSequence_Server interface {
	GetNumber(CallSequence_getNumber) error
}

func CallSequence_ServerToClient(s CallSequence_Server) CallSequence {
	c, _ := s.(server.Closer)
	return CallSequence{Client: server.New(CallSequence_Methods(nil, s), c)}
}

func CallSequence_Methods(methods []server.Method, s CallSequence_Server) []server.Method {
	if cap(methods) == 0 {
		methods = make([]server.Method, 0, 1)
	}

	methods = append(methods, server.Method{
		Method: capnp.Method{
			InterfaceID:   0xabaedf5f7817c820,
			MethodID:      0,
			InterfaceName: "aircraft.capnp:CallSequence",
			MethodName:    "getNumber",
		},
		Impl: func(c context.Context, opts capnp.CallOptions, p, r capnp.Struct) error {
			call := CallSequence_getNumber{c, opts, CallSequence_getNumber_Params{Struct: p}, CallSequence_getNumber_Results{Struct: r}}
			return s.GetNumber(call)
		},
		ResultsSize: capnp.ObjectSize{DataSize: 8, PointerCount: 0},
	})

	return methods
}

// CallSequence_getNumber holds the arguments for a server call to CallSequence.getNumber.
type CallSequence_getNumber struct {
	Ctx     context.Context
	Options capnp.CallOptions
	Params  CallSequence_getNumber_Params
	Results CallSequence_getNumber_Results
}

type CallSequence_getNumber_Params struct{ capnp.Struct }

func NewCallSequence_getNumber_Params(s *capnp.Segment) (CallSequence_getNumber_Params, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0})
	if err != nil {
		return CallSequence_getNumber_Params{}, err
	}
	return CallSequence_getNumber_Params{st}, nil
}

func NewRootCallSequence_getNumber_Params(s *capnp.Segment) (CallSequence_getNumber_Params, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0})
	if err != nil {
		return CallSequence_getNumber_Params{}, err
	}
	return CallSequence_getNumber_Params{st}, nil
}

func ReadRootCallSequence_getNumber_Params(msg *capnp.Message) (CallSequence_getNumber_Params, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return CallSequence_getNumber_Params{}, err
	}
	return CallSequence_getNumber_Params{root.Struct()}, nil
}

// CallSequence_getNumber_Params_List is a list of CallSequence_getNumber_Params.
type CallSequence_getNumber_Params_List struct{ capnp.List }

// NewCallSequence_getNumber_Params creates a new list of CallSequence_getNumber_Params.
func NewCallSequence_getNumber_Params_List(s *capnp.Segment, sz int32) (CallSequence_getNumber_Params_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 0}, sz)
	if err != nil {
		return CallSequence_getNumber_Params_List{}, err
	}
	return CallSequence_getNumber_Params_List{l}, nil
}

func (s CallSequence_getNumber_Params_List) At(i int) CallSequence_getNumber_Params {
	return CallSequence_getNumber_Params{s.List.Struct(i)}
}
func (s CallSequence_getNumber_Params_List) Set(i int, v CallSequence_getNumber_Params) error {
	return s.List.SetStruct(i, v.Struct)
}

// CallSequence_getNumber_Params_Promise is a wrapper for a CallSequence_getNumber_Params promised by a client call.
type CallSequence_getNumber_Params_Promise struct{ *capnp.Pipeline }

func (p CallSequence_getNumber_Params_Promise) Struct() (CallSequence_getNumber_Params, error) {
	s, err := p.Pipeline.Struct()
	return CallSequence_getNumber_Params{s}, err
}

type CallSequence_getNumber_Results struct{ capnp.Struct }

func NewCallSequence_getNumber_Results(s *capnp.Segment) (CallSequence_getNumber_Results, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	if err != nil {
		return CallSequence_getNumber_Results{}, err
	}
	return CallSequence_getNumber_Results{st}, nil
}

func NewRootCallSequence_getNumber_Results(s *capnp.Segment) (CallSequence_getNumber_Results, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0})
	if err != nil {
		return CallSequence_getNumber_Results{}, err
	}
	return CallSequence_getNumber_Results{st}, nil
}

func ReadRootCallSequence_getNumber_Results(msg *capnp.Message) (CallSequence_getNumber_Results, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return CallSequence_getNumber_Results{}, err
	}
	return CallSequence_getNumber_Results{root.Struct()}, nil
}
func (s CallSequence_getNumber_Results) N() uint32 {
	return s.Struct.Uint32(0)
}

func (s CallSequence_getNumber_Results) SetN(v uint32) {
	s.Struct.SetUint32(0, v)
}

// CallSequence_getNumber_Results_List is a list of CallSequence_getNumber_Results.
type CallSequence_getNumber_Results_List struct{ capnp.List }

// NewCallSequence_getNumber_Results creates a new list of CallSequence_getNumber_Results.
func NewCallSequence_getNumber_Results_List(s *capnp.Segment, sz int32) (CallSequence_getNumber_Results_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 0}, sz)
	if err != nil {
		return CallSequence_getNumber_Results_List{}, err
	}
	return CallSequence_getNumber_Results_List{l}, nil
}

func (s CallSequence_getNumber_Results_List) At(i int) CallSequence_getNumber_Results {
	return CallSequence_getNumber_Results{s.List.Struct(i)}
}
func (s CallSequence_getNumber_Results_List) Set(i int, v CallSequence_getNumber_Results) error {
	return s.List.SetStruct(i, v.Struct)
}

// CallSequence_getNumber_Results_Promise is a wrapper for a CallSequence_getNumber_Results promised by a client call.
type CallSequence_getNumber_Results_Promise struct{ *capnp.Pipeline }

func (p CallSequence_getNumber_Results_Promise) Struct() (CallSequence_getNumber_Results, error) {
	s, err := p.Pipeline.Struct()
	return CallSequence_getNumber_Results{s}, err
}

type Defaults struct{ capnp.Struct }

func NewDefaults(s *capnp.Segment) (Defaults, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 2})
	if err != nil {
		return Defaults{}, err
	}
	return Defaults{st}, nil
}

func NewRootDefaults(s *capnp.Segment) (Defaults, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 2})
	if err != nil {
		return Defaults{}, err
	}
	return Defaults{st}, nil
}

func ReadRootDefaults(msg *capnp.Message) (Defaults, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Defaults{}, err
	}
	return Defaults{root.Struct()}, nil
}
func (s Defaults) Text() (string, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return "", err
	}
	return p.TextDefault("foo"), nil
}

func (s Defaults) HasText() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Defaults) TextBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	d := p.DataDefault([]byte("foo" + "\x00"))
	if len(d) == 0 {
		return d, nil
	}
	return d[:len(d)-1], nil
}

func (s Defaults) SetText(v string) error {
	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, t.List.ToPtr())
}

func (s Defaults) Data() ([]byte, error) {
	p, err := s.Struct.Ptr(1)
	if err != nil {
		return nil, err
	}
	return []byte(p.DataDefault([]byte{0x62, 0x61, 0x72})), nil
}

func (s Defaults) HasData() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s Defaults) SetData(v []byte) error {
	d, err := capnp.NewData(s.Struct.Segment(), []byte(v))
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(1, d.List.ToPtr())
}

func (s Defaults) Float() float32 {
	return math.Float32frombits(s.Struct.Uint32(0) ^ 0x4048f5c3)
}

func (s Defaults) SetFloat(v float32) {
	s.Struct.SetUint32(0, math.Float32bits(v)^0x4048f5c3)
}

func (s Defaults) Int() int32 {
	return int32(s.Struct.Uint32(4) ^ 4294967173)
}

func (s Defaults) SetInt(v int32) {
	s.Struct.SetUint32(4, uint32(v)^4294967173)
}

func (s Defaults) Uint() uint32 {
	return s.Struct.Uint32(8) ^ 42
}

func (s Defaults) SetUint(v uint32) {
	s.Struct.SetUint32(8, v^42)
}

// Defaults_List is a list of Defaults.
type Defaults_List struct{ capnp.List }

// NewDefaults creates a new list of Defaults.
func NewDefaults_List(s *capnp.Segment, sz int32) (Defaults_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 16, PointerCount: 2}, sz)
	if err != nil {
		return Defaults_List{}, err
	}
	return Defaults_List{l}, nil
}

func (s Defaults_List) At(i int) Defaults           { return Defaults{s.List.Struct(i)} }
func (s Defaults_List) Set(i int, v Defaults) error { return s.List.SetStruct(i, v.Struct) }

// Defaults_Promise is a wrapper for a Defaults promised by a client call.
type Defaults_Promise struct{ *capnp.Pipeline }

func (p Defaults_Promise) Struct() (Defaults, error) {
	s, err := p.Pipeline.Struct()
	return Defaults{s}, err
}

type BenchmarkA struct{ capnp.Struct }

func NewBenchmarkA(s *capnp.Segment) (BenchmarkA, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 24, PointerCount: 2})
	if err != nil {
		return BenchmarkA{}, err
	}
	return BenchmarkA{st}, nil
}

func NewRootBenchmarkA(s *capnp.Segment) (BenchmarkA, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 24, PointerCount: 2})
	if err != nil {
		return BenchmarkA{}, err
	}
	return BenchmarkA{st}, nil
}

func ReadRootBenchmarkA(msg *capnp.Message) (BenchmarkA, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return BenchmarkA{}, err
	}
	return BenchmarkA{root.Struct()}, nil
}
func (s BenchmarkA) Name() (string, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return "", err
	}
	return p.Text(), nil
}

func (s BenchmarkA) HasName() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s BenchmarkA) NameBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	d := p.Data()
	if len(d) == 0 {
		return d, nil
	}
	return d[:len(d)-1], nil
}

func (s BenchmarkA) SetName(v string) error {
	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, t.List.ToPtr())
}

func (s BenchmarkA) BirthDay() int64 {
	return int64(s.Struct.Uint64(0))
}

func (s BenchmarkA) SetBirthDay(v int64) {
	s.Struct.SetUint64(0, uint64(v))
}

func (s BenchmarkA) Phone() (string, error) {
	p, err := s.Struct.Ptr(1)
	if err != nil {
		return "", err
	}
	return p.Text(), nil
}

func (s BenchmarkA) HasPhone() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s BenchmarkA) PhoneBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(1)
	if err != nil {
		return nil, err
	}
	d := p.Data()
	if len(d) == 0 {
		return d, nil
	}
	return d[:len(d)-1], nil
}

func (s BenchmarkA) SetPhone(v string) error {
	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(1, t.List.ToPtr())
}

func (s BenchmarkA) Siblings() int32 {
	return int32(s.Struct.Uint32(8))
}

func (s BenchmarkA) SetSiblings(v int32) {
	s.Struct.SetUint32(8, uint32(v))
}

func (s BenchmarkA) Spouse() bool {
	return s.Struct.Bit(96)
}

func (s BenchmarkA) SetSpouse(v bool) {
	s.Struct.SetBit(96, v)
}

func (s BenchmarkA) Money() float64 {
	return math.Float64frombits(s.Struct.Uint64(16))
}

func (s BenchmarkA) SetMoney(v float64) {
	s.Struct.SetUint64(16, math.Float64bits(v))
}

// BenchmarkA_List is a list of BenchmarkA.
type BenchmarkA_List struct{ capnp.List }

// NewBenchmarkA creates a new list of BenchmarkA.
func NewBenchmarkA_List(s *capnp.Segment, sz int32) (BenchmarkA_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 24, PointerCount: 2}, sz)
	if err != nil {
		return BenchmarkA_List{}, err
	}
	return BenchmarkA_List{l}, nil
}

func (s BenchmarkA_List) At(i int) BenchmarkA           { return BenchmarkA{s.List.Struct(i)} }
func (s BenchmarkA_List) Set(i int, v BenchmarkA) error { return s.List.SetStruct(i, v.Struct) }

// BenchmarkA_Promise is a wrapper for a BenchmarkA promised by a client call.
type BenchmarkA_Promise struct{ *capnp.Pipeline }

func (p BenchmarkA_Promise) Struct() (BenchmarkA, error) {
	s, err := p.Pipeline.Struct()
	return BenchmarkA{s}, err
}

const schema_832bcc6686a26d56 = "x\xda\xacZ\x0dt\x14U\x96\xae\xd7?\xa9\x04\x92t" +
	"W\xaaC\xcco\x9b\x08J2\x07\xc9\x0f\x83\xc8\xcc\x9c" +
	"\x90\x98(z\xc0\xc9\x0f\x88d\x8d\xa6\x92T\x92\xc6\xee" +
	"\xae\xd0]\x0d\x09\x0e':\x03\x0a\x1c\xb3\xa3\xab\x8e\x8a" +
	"\xb2\xe3d\xf1\xac\x8c\xe2\x8a#s\x84U\x06\x18\x18!" +
	"\xe2\x8e\xcc\x01\x05V\x10p\xd5\x15a4\x8e\xac\xe0\x0f" +
	"\xb5\xf7\xbe\xfa\xed\xeej\x189sN~\xea\xdd\xef\xd6" +
	"}\xf7\xddw\xdf\xbd\xf7\xbdW\x95\x7f\xcd\x9c\xe5\xa8r" +
	"\xdf\xc43L\xf3qw\x9ar[h\xe4\xfe\x9e}?" +
	"\xf8\x05C2\x18\x07\xbf:cY\xcdp\x06\xeb\xe4c" +
	"\x99\xd71\x0ee\xf7\xb93\xefT\xde3i%\xd3\xec" +
	"!\xc4du\xb1\x0c\xc3o\xca\x1c\xe5_\xcf\xc4\xa7-" +
	"\x99\xb5\x0cQ\xda\xdf\xf4\xde\x99\xf1\xda\xec\x07\x12x\xdd" +
	"\x0ed9\x9c\xb9\x95?A\x99\x8ff\xfe\x070\x1f\xfa" +
	"\xdd\x85\xca\xabj\xff\xf8\x00\xc3y\xac\xbc\x048jV" +
	"g\xe5\x10\xfe\xc9,d~,\x0b%\xcf\x1c\x985\xfd" +
	"\xe5\xb7JW%Hn`\x1d\xc0\xb2;k\x94\xdfO" +
	"\x99\xdf\xcaZ\x0a\xcc#\x7f+|\xf7\x95\xbb&\xaca" +
	"8\x1fa4\x89\xe5\xd9\x0e\xc2\x10~J6J\x9b6" +
	"i\xf2';\xcaZ\xff\x19\xbav\x9a\xc2\x00\x9e\x9b=" +
	"\xc2\xcf\xcfFI\xcd\xd97\xf1\x83\xf8\xa4<\xb4\xc7\xff" +
	"\xa7\xba{^\xfce\x82\x9etLB\xf6I>D\xf9" +
	"\x03\xd9\xd8\xb3\x12\xd9w\xaa\xf9\xc9\xbd\x8f\xc4\xf3R[" +
	"\xed\xcf\xde\xc9\x1f\x06V\xa7\x92y\xe0W\xbbr\x9e\xaf" +
	"|\x14\x98\\q\xbdo\xce\x1e\xe5w\xa0\xb4\xd6\xd7\xb2" +
	"\x9d\xa4\xf5\x10\xaa\xcc(\xc3\xa5\xe9\xd3\xce-\xff\xc3\xa3" +
	"6v\xe2\x0f\xc2\x1b'h\xffG\xe9\xc0\xee\x0ew\xf8" +
	"j/l{\xcc\xce\xa6\xdfe\x83M\xb3<\xc8\x9c\xe1" +
	"A\xe6\xe1\x07\xfc\x7f\x9a\xbd\xfa\x83\xc7\xd1\xa6\x8e\xc4\x91" +
	"]\xef\xd9\xc9\xd7!s\xcdO<~\xb0\x9c\xb2\xf0\x05" +
	"\xc7S\x8f?\xb9y\xad\x9d\x1a\xf3\xbd\xa3\xbc\xe0\xc5\xa7" +
	"v/J^w\xcf\xd1-\xe5\xef^\xf7\x94u\x02V" +
	"z\xc7\xe1\x04\x0cS\x86=\xf3O\xba\xb7^\xf3\xcb\xa7" +
	"\x92L\xb0\x11$mAI\xad\xafx\xc1\x04\xdb\xbd\xd4" +
	"\x04\xb1\x1f\xb9\xeeW\xaa+\xd7%\xfa\x15\xed|\x13\xbc" +
	"\xf2:\xed|\x8b\x17\xe7\xe0\xf1\x92\xfdS\x9b\xce\x8b\xeb" +
	"\x99\xe6\"\x02/\xa3\xf1kr\xb9\x08v^\xcaa\xe7" +
	"=\xb3{\xbf\xb8\xc0\xff\xe19\xbb\x91\xd4q;\xf9\x9b" +
	"9|j\xa4\xbcW\xee\xc9\x1b\xb8\xeb\xfd\x17\x9fO\xf2" +
	"\x14\x91;\xc9/\xa6\x8c!\xee&\xfeI|R\xc6\x96" +
	"\xcf\xad\xa8\x9f\xff\xf6\xf3v\xd6\xbf\x8f+ \xfc\xc3\xf4" +
	"\x85a*y\xcb\xa7/I\x8dG\x1e\xd8h\xa7\xc5\xeb" +
	"\xdc\x08\xbf\x9b\xf2\xee\xa0\xbc\\\x7f\xf7\xbba\xf7\x0b\x9b" +
	"\xecxOp_\xf0g(\xef'\x94wh\xfa\x9d+" +
	"\xdag|\xbe\x09meQ\xd9\xedD\x96\xfc\x9c\xbf\xf0" +
	"\x93rP\x9f\xd2\x9c\x058\xab\xd1\xd1i\xca\xe9\x93%" +
	"\xbf\xb7\xf3\x81\x9aa\xde\x01\xab\x90\xa7\xab\x90\xff\x18\xb8" +
	"\x8f\xff`\xf0*\xef\xd0\x86mLs\x86\x9b(\xab\x82" +
	"G/4\x17U\xecGs\x04|k\xf8\xc5>\x9c\xb7" +
	"\xa0\xcf\x89\x82G\xc6\x7f1r^<\xf2\x86\x9d\xca\xed" +
	"\xbe\xdf\xf2\xa2\x8f\xae\x1f\x1f\xaa<e\xee\xf5\xaf\x7f\xf0" +
	"\xfc?\xed\xb5[a\xf7\xf9F\xf9a\xca\xbb\xda\x87\xb3" +
	"\xbb\xf5\xab\xf7\x0f\xde\xb5\xe4\xbd7\xedl\xfc\xa1\x0fl" +
	"|\x962\x8fQ\xc1O\xec\xda0\xfe#\xaea\x9f\x9d" +
	"\x12\xa5\xb9[\xf9\xf2\\|\x9a\x94\x8b\xbc\xc7\xaf\x9e\xd1" +
	"\xf1\xc1\xcb\xbf\xb3\xe5]\x98;\xc2\x0b\x94\xb7\x9d\xf2\xde" +
	"\xdct\xf2\xf0\xc9\x17\x1b\xfe\xcb\xd6\xc6\xcbsO\xf1\xab" +
	"s\xa9\xcb\xe7R\x1b\xef\xb9wG\xd1\xe8\xa9\xa7\xfel" +
	"\xa7\xf2\x89\x09\xb0(\xc7&\xe0{g&\xa0h\xc3\xa8" +
	"\xc4\x09\xb1\xb88\xef\x16\xbe4o)?\x98\xe7\x87P" +
	"\xbc\xed\xc8\xaa\xb1\xf5\xdfL\x7f\xc7N\xc5\xc1\xbc\xb5\xfc" +
	"}y\xb4\xff<\x94\xf3\xf0u\xff\xde\x11\xfe\xf3k\x87" +
	"PEW\xa2Q\xd7\xe5\x8d\xf2\x1b\x90\xb9\xe6\xd9<\xaa" +
	"\xe2\xf0\x9b\x87\x97\xae\xeax\xf0\xb0\x9d\xe4\x8c\xfc\x11\x9e" +
	"\xcb\xc7\xa7\xac|\x94\xcc/\xf9&\xd0S\xb7\xff\xa8\xdd" +
	"lM\x01\xde\x1fR\xde\xaa|\x9c\xad\xa2\x1d]\xe9\xbf" +
	"*\xa88\x96h(U\x8b\xfc\xbf\xf0\x1b\xf2\xa9\x16\xf9" +
	"T\x8bu\xf3\x16l\xfc\xcf\x17\x9b\x8e\xd9\xa5\x9a\xf3\x05" +
	"\xbf\xe5I!>}W\x80\xd9\xe3\x91\xca\xe7\xbe\xfe\xf1" +
	"\x81\x7f=fg\xd4\xc7\x0a\xc7\x11\xfeY\xca\xfc\x9bB" +
	"TysC\xf6\xd5\xe4\xf7\x95'\x92\xfdv\x7f\xe1\xcf" +
	"\xf9\x83\xc8\xd9\xfav!\xf5\xdb\x8d\xdbX\xee\xe0\xfe\x91" +
	"\x13v\x96\xd8Q\xb8\x95\xdfK\xc5\xee\xa6b\x85\xd6\x9a" +
	"\x9c\xdd\xa7\xf6\xda\xf2\x9e)\\\xcb\x9f\xa5\xbcc\x94\xf7" +
	"\xdc\xc2g~\xf1\xf4H\xfa\x87v\xfarE\xe0\x04\xa5" +
	"E\xc8\\\\\x84\xcc\x9b\xde\x98\x7f\xecE\xef\xad\x1f&" +
	"\x18\xa2\x91\xb0.\xe0\x99[\xb4\x93\x9fO\xb9\x9b\x8bp" +
	"U\x96N=W\xf8\xed\xcav\x14\xed\x88\x8bQ\x93\x8a" +
	"\xb7\xf2S\x8a\x91\xb1\xbc\x18m\xf6^\xda\xf9_\xaf\x18" +
	"\xba/Q\x07\xd5k\x8bG\xf9\xd5\x94w%\xe5]\xd8" +
	"W\xf0\xd9\xccOW|d7\xb6\xb9%G\xf8\x85%" +
	"4\xf0\x97\xd0\xacsp\xdb\xfa\x8d\x05\x8b?N\x8a\xe6" +
	"\xcbK@(2\xb6\xae(\x81h\xfeP\x09\x8d\xe6\x87" +
	"\x9fy\xe7G\x8f\x7fr\xf5\xa9\x84\x98\x03\xa3\x9b\x8c\xdd" +
	"\x97\xac\xe1\x87Kh\x11P\xf2F:\xfaPz{}" +
	"\xc6\xc6\x9f\x8e\xd9irv\xf2\x11\x9e\x94S\xaf\x98L" +
	"\x17|~\xce\x9a/\x7f\xfe\xc0Y\x86+\xd2C\xff\x94" +
	"\xf2E\x84q)\xdf\xf5\xdc\xb4\xb7\xf1=\xf7W\x09\x9d" +
	"R\xdf\xca/\x87\xa8H\xa5\x94\x96\xa3\xd7.\xba\"x" +
	"\x93\xafA\xf9\xca\xae\xc7\xe5\xc0\xbb\x9a\xf2\xae,\xc7\x1e" +
	"\x8f\xcd\xd9\xf6\xc8d\xf9\xdf\xbe\xb5\xf3\xd9\xbd\xc0{\x90" +
	"\xf2\xee\xa7\xbcB \xd2\x15\x11zd\xa6\xf6\xda.\xa1" +
	"?\xdc\xdf\xbc\x9dX\x8c\xc7\x0f\x92\x16\xb36\xa0-#" +
	"M\xd2\x96\xb1<\xa0Um\x06X>\x06-\xc3\x07\xf8" +
	"\xc5\xa4\xde\\\xfc|\x08\xde3\"\x1b\xb4*\xcc\xa5\xce" +
	"\x07\xa0e\xb80/\x9223w\xf0\x02i3\xbd\x10" +
	"Z\xb7\x98\xb3\x06\xad\x1c\xb3\xd2\xe3\xdb\xa1?#0\xf1" +
	"\x0bA\x8a\x11u\xf9\xf9\x80\x19\xe1\x82o\x86\xfe\x8c\x02" +
	"\x89\x9f\x0b2\x0d\xe3A\xab\xcd\x9c\"\xda2\xca\x1eh" +
	"\xb5\x98)\x82\xb6\x8c\xac\x05\xad5f(\x80\x1e\xfe\xc5" +
	"L\xc2\xd0\xfb\x88\x99.@\xb3\x11s\x11\x82\xd6k\xcd" +
	"\xb8\x0c#Zk\x96N`\x89\xb5fq\x0aVZk" +
	"\xc6{\xb0`\xc4\\F\xd4\xbaF]D[F\xf8\x80" +
	"V\xbd\xb9\x88@J\xa7Y\xc3B\xab\xc5tk\x8a\x19" +
	".\x07\xad63\xdfCk\x99Y\xad\xd2\x193\xea\x06" +
	"\xd0\xb3\xc2\xace\xe8\x1c\x19\xa5*\xb4\x16\x99U\x13\xb4" +
	"Z\xcc:\x9e\xb6\x8c\xba\x86r\x1a\x95 \x95b\x04m" +
	"\xea\x05]R8*7\x082CD\xf5yN \xca" +
	"\x10Y}n\x0c\xc7\x18\x12\xf2\xb7u\x0b\xb2H\xff\x0a" +
	"Cu\x81H\xbf\x14\x91\x95\xa6\xa0\x10\x16\xeb\x05\xe0\x15" +
	"=\xf5\xd7\xd5\\\xe7\xa9\xab\xa9\xaedo\xac\x9a\xae\xb4" +
	"\x88\xbd\x111\x1a\x0d0N)\xac\xd4\xe9k\x02\xe2E" +
	"\xdb\xd0\x0dR,,\x8b\x11\xb6^\xe8\x1dj\x8b\x8a\x91" +
	"%b\xc4\xd3\xb6H\xeaTn\x13#\x8d\xa1~y\x10" +
	"\xd8\xf0\xf9\xa7a\xb1A`\x9c\xb2\x80\x8dyK%\xb3" +
	"\x01H\x13\xe8\x1a\xd1\x00\xebs\x83@d\x81\xd2\" " +
	"e\xb6\x14\xec\x8e\x02@\xa8\\\x18\x14\x8c\xc9\xa0Q\xf9" +
	"\xb20\x87\x0dD-T*#\x89J{\x8c\xcc\x098" +
	"\x13Xm\x89\xf0cKl\x0a:cQM\xcfyK" +
	"\x19?\x10\xa0M\x99\xe6\x89\x03h\xef\x05\x11\xa1\x1fU" +
	"e\xc8\xe0\x10>W\x0fT+\xda\xff~\x86E\xee\xdb" +
	"\xa4@\xf7\xfcp@bHX\xb9U\x8c\x82\x1d\xabn" +
	"`X\x084J\xcb\x82y\xd0\xbe\x01\x8c\x04\x0d\x18\xaa" +
	"\xdc*GbLm\x17\x90\xfa\xc3\x9e\xc6\xae>\xc93" +
	"[\x92\xfb\x14|\x82)\x13\xd1\xca\xad\xb2\xd0uw " +
	"\xdc\xcbxZ$I6\x9b\xa4\xce\xf2\\\xaf\xdc \x04" +
	"\x83\xad\xe2\xe2\x18\xe3\x11\xc3]\xa2\xd2 \xf6\x08\xb1\xa0" +
	"\x1cE\x09\xf5@\xe9\x0b\x09\x11\xc6yw]s:D" +
	"]\xb3|\x84\xb2\xbe\xeeJk^.\x02\xc2,B2" +
	"\x19\x86\x83\xd5`\x04J6\x18\xe8\x04\xc7\x00r\x15Y" +
	"C\x94eR\xa83 .\x13\xdd\xe1k\xbb\xa4\xd0\xd4" +
	"^i*\x8d\xa3\x11I\x96\xaa\xa7\x06\xd0s\xc2Bp" +
	"\xaa\xfe6\xbek\x06]\x87\x1asg\xea\xca\xd73L" +
	"\x13!\xcd.'\xa4T\x17$\x0b.\xab\x0c\xb6\xc3\xe9" +
	"N\xd2\xecs\x106\x1c\x0b\x11\x17\xe3\x80_S\x02\xd1" +
	"$\xdcP\xabz)\xbe\x9ei\xbc\xdeX\x01\xaf\xcf\x82" +
	"\xd7\xe78\x08!>\xdczp7W\x03\xad\x01hM" +
	"\x0e\xc29\x80\x08\xbbUn\xee-@\x9c\x03\xc4>\x07" +
	"\xf1D\x03\xcbD\xe2\x86\x8e\xdc\x0c\xf1/\x95\"\xddQ" +
	"\x18\xac\x03~\x89\x82\xad \xcc\x15\xae\x8cl\x8649" +
	"\x09\x85\xb2-\x1a95\x8dtg\xd2})F\xa2\x09" +
	"c\x9b\xa9\x8dm\xa2\x83\xd4\x86\x06Q\xac.\xd3kF" +
	"7\x86\xc4I\xd7-\xa6\xfbUX\xb5X\xba\xd3\x95\xa9" +
	"(Tly\x0e\x88\x9d\x08b+\x1d$\x8b\\P\xd4" +
	"QOA\xead\xa0N\x03K\x08L\x1a\xe9d\xd2\x92" +
	"TFW\xbbV\x84?\x13\x9b\x84\x88\x10\x8a2Vm" +
	"\x0b\xcc\x99p\x06\xc2\x86E\x12'\xa2\x91\x85\xf7\xd5\x81" +
	"\xba\xc1\xdf\xf4m?\xd1\xb7\x9f\x1cW\xc1887\xeb" +
	"\xc1~f\x11\xe0L\xe9\x0d-,8\xb96:B\xe0" +
	"\x87p\xe5\x8b\xccax\xb5\x09\xbd\x1e\x876\x0dh\xb3" +
	"\x1c kA@\xee\x03\x87g<\xe8\xf2`H#\xe0" +
	"\x82!\xbd\xd4X\x04\x9c\x82\x08I\x90\x8d\x8d\xf5\xf8f" +
	"\xa7\xa4\x1eveU\xa61qFA@'\x0eVN" +
	"\x1e\x97N<\xef\xb3\xe9\x85\xf8\xa7\xc8\xae\x135X]" +
	"\xd4\xf5\xfb\xe5\x08\xc86Rs\x82\xc2\xf6.G\xc3\xe7" +
	"\xf7p9#\x81\xa7p9#\x8c\xa0L\x9f*\x13\xa6" +
	"d9\xae\xb1\x01\x90\xb9\x02\xd6\x13]d@\\\x89\xc4" +
	"{\x81\xf8 \xb8\x9b\x03\xd6\x18\xd0\x86q\xe1\xad\x02\xda" +
	"\xa3\xc0\xe8\x04F\xe8\x95{\x18G\xf9 \x10\x9f\x00\xa2" +
	"\x0b8A&\xf7\x18\xbe\xfd\x10\x10\x9f\x86\xd5(\x8b\x03" +
	"\xb2\xe6m`\xcd2\xb6G\x92<\x98\xd4H\x16\xd0\xb2" +
	"\x90V\xc0v\x0a\x11\x7fOP\x12d2\x8eq\x8c\x8d" +
	"\xfb\xe3\xd9\xd9\xb3\x18\xc2B\xf0\xc1x1\xe6Z\xa9(" +
	"\x0aC<1$\xa4\x83\xf7\xa5W\xd8\x0cO\x8f\xe5\x83" +
	"I\x13\xf12h\xe3\x05m \x1e*\xa1@o\x9f|" +
	"\xab$\x93z\xb1E\x840;\xe8\xa7\xef\x80\xf9\x8c*" +
	"*\xc5\xe4\x98\x8b\xabE\x8c\xaa\xe18\xd5lK19" +
	"i}\xc5\xf9\x1d\xa4\xf8\x90\xeaw\x1e\xb3\xcad\x88\xdb" +
	"C.\x12\\\xeb\x8cP\xa1\xf5Y^fF\x0a=:" +
	"Z\xe3\x845\xe0\x92N\x18\xa2Q\xae$\x0c\xd1\xa5\x07" +
	"a-\xef`\xd6\xb9\xb6W\x94o\x8d\x85:\xc5\x08\x8c" +
	"\xd7O\x07l\x1dn\x8e9\\\x12\xc6I!\xe96\xaa" +
	"\x1b\xb9/aF*\xcc\xb7i\x1c!\x9cY\x90\x81n" +
	"\x9c\x8d(C7\x16\x943\xe3\x93\xbe9!\xfa\x09\x15" +
	"\xc7\xb5\x80\x87d\xb0\x8a\xae?\xd4*\xf1a*q\xb5" +
	"\xe9u\x09\x96\x10\xdfc\xb9Y\x17s\xb6M \x9d\xcd" +
	"B\xf6\xbf\xc8\xb8;\xc1. \xc6(6S\x041\xbd" +
	"\xb2\xe8\xf2\xd3\xc2\"E\x8c\x01\x05\xd9%b\x97\xa9\x9d" +
	"Q\x09\xa7\x08\x06z\xbd(\x11*\xf2\x0aC\xe4\x93\xa8" +
	"\xe3\xa3 \xf2\x19\xd3\xa9\xd6a\xeex\x02h\xeb-)" +
	"\xf77\xc8\xf84\x10_\xc5p\x00+\xdf\x09\xc4\xcdh" +
	"\xb1\x97\x80\xb8\x0f\xc3\x01p\x82Xn/j\xb9\x0b\x88" +
	"o\x03\xd1\x0d\x9c0q\xdc[H\xdc\x03\xc4\x03\xa6-" +
	"\x8c-\x96j\x0bgg%\x19\x0f\x9e5\x1e\xd6~\xa7" +
	"\x08\x11C\x1b\xddx5c\xd7\xf6c5\x1c5\xc7l" +
	"\xec\xad\xd41\xb3\x83\xa1\x98\xfe>;\x18\xed\xd6\x9f\x93" +
	"\x1cA\xafd\x8dB\x16-\xe25,\"\xa0\xa6w\xa8" +
	"\xb5\x05\xa7\x9bDDb\x07\x10\x83Z\x80\x04Z\x00-" +
	"\xd2\x0d\xb4~=@\x021\x84\xc4> \xca8EB" +
	"\x10\x8f\xa0\xe0\x97\xb0\xdd1I/R<\x90\x1e\xaa\x92" +
	"\xf3\x03\x92\xabm\xd2F\xaf\xa4*Nf\xf6Cd\x10" +
	"zE5\x94@\xc4Ir\xc26\x1aj/\xe2\x85\x96" +
	"Hl\x9b@\xf5\xd2>!\xf2T\x98\x91\xc70\xc9\x94" +
	"\x0a3\xf4|\xcf\x11]$\x11^\xce\xd246\xbd)" +
	"\x96f[\xad\xba\x0dJ\x10\x09U1\xd4\x9e\xa4y2" +
	"\xe4\x89\xa5B@\x86\x98\xbb\x88a\xa5N\x8b\x8b\x19\x1b" +
	"\xef\x14\x92\xebY\xd8\x95]z\xd1'8z\x0a\xc3\xc3" +
	"\x8f\x07\xb7,\x09\x0b\xb4\xcc\\\xa0\x86\xed\xd7\x95\x99+" +
	"TwGc\x81>gq\xc7g\x91\xf8\x0c\x10_\xd0" +
	"\xf35\x107\xe0\xdb\xeb\x81\xf8\x92e\x81nD\xce\xe7" +
	"\x80\xb8\xeb\xd2\x8ek-=,\xd3\x9c@f\xe5\x88Q" +
	"\x97{\x82Q\xb9F\xb7\xab\xfb\xd2%8\xdd\xf8}\xbf" +
	"z\xc88\xde\xd0\xe6J_6P\xfb\xe3-\x86\xb9s" +
	"\xe2\xaa\xea\xcd]\x137e\xa6r\xe7#\xbfn~\xfd" +
	"\x9d5\xbb1\xc5*o|6:1\xff\x15\xf9Y\x86" +
	"\x9bT\xa6\xe4\x1co9=x\xff\x92=\x0cWZ\xad" +
	"<r\xe5\xd4\xe3kE\xef\xd7\x0cW\xdc\xa6\x8c\x0dO" +
	"\xcd\xcb\xe9\xd8\xb2\x13\x1a\x15C\xda\xba\xac\x0d\x84p\xdb" +
	"\xcevK]\xac,\xf4\xfa\xc3\x12\xfcU\xbabQY" +
	"\x0aA\x9d\xe2\xec\x17=a!$6\xbb\xe2\xf6q." +
	"\xd8\xb6y\xb5m[\xb5_S9\xd9\xd3\x9cBo\x82" +
	"1\xeaMG\x1b\xeaR\xb7S`\x06\xe3T)\x85\xaf" +
	"\xe9g\x0aj\x14I\x91\x09LG\xab\xd6\x1c\xed\x15\xd3" +
	"\xd16\xe1,\xbc\xa0g\x82Y\x09\x99\xe05\x8b\xa3m" +
	"\xc1m\xda\xab\xaaOqn\xa7\xeah;\x90\xb8]\xcd" +
	"\x19\xd4\x1ez\xed\xe4\xef\x93Bf|\x8f\xab\x94h\xfc" +
	"\x8f\x08\xb8Du\x8f\xaa\xed\x12\xc27\x06\x07AK\xd0" +
	"\x0c\xc6\x09\xa3\x13\xba\x02\xf4tCgQB\xc2@k" +
	"\xbf(v#-1\x1b\xe8\x86\xadck\xaa+/\x7f" +
	"\x09\x1b1\x86]$u\xa6\xae\xd7\x92\xa3f\x03,\xb3" +
	"\xaeP\xb7>x\x8f\x10\xe9\x8d\xa6\xda\xa3\xeas\xa7\x9f" +
	"\x06\xdcM\xea\xfe\x8e,~\x8b%\x1e\xe8Y\xfc\xd9j" +
	"K<\xd0\xb3\xf8\x86[\xb4\xa5\xff\x0a\xce]\x87:w" +
	"q\xb3\xac\x07\x89\xcd\xd5\xe6,\xc7\xcd\x9d\xd2\x19\x88\xc0" +
	"FN\xb0\x9a\xdf\xdf\xdf'\x85M\x8eh\xa03\x08\xf3" +
	"\x87\xc7\x1az\x95Z\x1b\xed\x97b`^m\x0e\xfd!" +
	"\xe0\x1fL9S4\xb9\x89\xa9O\x0c8t\xd0\xe4#" +
	"\x03\xa7vd\x80\xb31\x1b\x88\xf3@\xf5AQ\x88\xe8" +
	"\x01\x0e{\x95\xfbH\x1a\xb4\xd20\xdc\x09\x83\xfas\xca" +
	"\x00\xa5\x1fw\xc9I{\x8f\x8b\x05(\xeb\x8e#;." +
	"\xaf\xabq#eV_P\xab\x9eQ]\xd6>'!" +
	"('\x09\xbf\xd1Y5\xfd\xf2\xdd\xdf\xa6\xb8\xbe\x8c\xd0" +
	"m\x9cS\xa7\xa8^\x8d\x03N\xadR\xd3\xcfN\x84\x0a" +
	"\xb3T3\xcfN\xc4\x0a\xb3V\xcbr|\xa7$Wk" +
	"Y\xceo\x15\xad\\+3\xcb5\xcf\x12)\xd0\xcd\xa4" +
	"y:!\xa7\x83R\xc6\x91\xbf\x96\xe4\x04\x88\x13\xa8\xab" +
	"~\xf6\xaf%\xb9\x9e\xaa\xe9@5\xee\x00RX\xb9\xae" +
	"V=\xd2\xc5\x01\x14\xa9\xfeXA\xb7>\x8d\xa0\x00q" +
	"p?\xc1\x7fN\xee\x87\xf8\xcf\xc5M\xc1\x7fnH@" +
	"\xf0/\x0d\x12\x0c\xc3x\xc2\xb04\xd8E=w\xb3A" +
	"a\x80\x8d\xf6Hl0\xb6\x84\xed\xeeY\x0a\xfb\xed\xa8" +
	"\x9cd0\xfd\xd0SsP\xcbz)\xb3\xae\x17\xfd\x88" +
	"\xadL[/\x1d\xb8^\xb4(\xdf\x8e\x13v\xbbzN" +
	"\xc0\xca\xc6\x9e\x9e\xb0As\xe6\xb48U\x0b$\x0b\xf5" +
	"\x12AL?p\xc5\xd3\xda\xcb\xf2h\xebY\x9b7\xd5" +
	"\x96\x1b\xeaH\xadN\x8e?\xe3\xf12\x8cW=\xdfI" +
	"\x8e/\xaa\xb1\x9e1\xfc\x8bov\x81\xe9[\xe7\xb8\x9c" +
	"\xa4\xf5v\x97\xc5\xc5\xf8\xf9.\xd8\"\xb56!p\x07" +
	"\x00\xc5\xe0ej\x80\xe5\x17\xba\xc0\x98\xad\xf3\x10\xe9@" +
	"\x04<M\x8d\xb2|;EnG\xa4\x1b\x11\xd77\x8a" +
	"\xbaa\xe2\x05\x8a\xdc\x81H\x1f\"\xee\xaf\x155\xde\xf2" +
	"\"E:\x10\x09\"\x92v^q\xf9 41|\x80" +
	"\"\xdd\x88\xf4#\xc2\x9eS\xd2}\xf4v-Du\xeb" +
	"CDF$\xfd+\xec'\x1d\x90\xc5\xf4\x9d \"\x03" +
	"\x88d\xfc\x1f\xf6\x93\x01H\x8c\"\xfd\x88\xfc\x0c\x91q" +
	"g\xb1\x9fqx'N\x11\x19\x91{\x11\x19\xff%\xf6" +
	"3\x1eo\xf1h?\x03\x88\xac@$\xf3o\x0a\xd4\x02" +
	"\x99\xf8\xbd\x015\xdb\xcf\x10Y\x85f\xcb\xfa\x02\xcc\x96" +
	"\x85\xd7}\x14\xb8\x17\x81\x07\x11\xc8\x1e\x03 \x1b\xbfK" +
	"\xa0\xc0\x0a\x04\x1eB\xc0\xf39\x00\x1e\xfc\xce\xc3\x05n" +
	"\x08R\x00X\x8f\x80\xf73\x00\xbcx)M\x81\xa7\x11" +
	"x\x15\x01\xee\xaf\x00p\x00l\xa6\xc0K\x08\xecC " +
	"\xe7\x0c\x009x\x93H\x81]\x08\x1cG\x80?\x0d\x00" +
	"\x8f\x9f\xfdP\xe0\x10\x02_\"\xe0\xfb\x14\x00\x1f\xde:" +
	"\xbb \x91\xb4\x9eF \xdd\x0d@\xee)\x00r\x01p" +
	"\xbb\xe1\x8d\x167\xd0\x8b\x90>\xe1\x13\xa0O\xc0KP" +
	"\xa4\xb7\xfa\x10\xa8D \xef\x7f\x01\xc8\xc3;}\x0aL" +
	"F\xa0\x01\x81+>\x06\xe0\x0a\xfcV\xc6\x8d]\xfc\x18" +
	"\x81\xdb\x11\xc8\xff\x08\x80|\xf407Z\xa4\x09\x81 " +
	"\x02\x05\x1f\x02P\x80\x13\xef\xae\xc7\x89G\xe0!\x04\x0a" +
	"\xff\x07\x80B4\x15\x15\xb5\x0a\x81G\x11(\xfa\x00\x80" +
	"\"\x00\x1e\xa6\xc0\x83\x08<\x81@\xf1I\x00\x8a\xf1\x83" +
	"\x147l\x81\x80\x19\x80\x97\x10(9\x01@\x09\x00\x1b" +
	"\xddP\x11\xb4\xbe\x80\xc0\xab\x08\xf8\x8f\x03\xe0G\xe3\xba" +
	"\xdb\xf0\xb3\"\x04\xb6#p\xe5\xfb\x00\\\x89_\xdb\xb8" +
	"[\xf0\x93+\x04\xf6\x00P\\z\x0c\x1d\xa8\x14\xaf\xf9" +
	"\xa9\xbe\xdb\x11\xd9\x87\xaf\x94\x1d\x85W\xcap>\xe8\x08" +
	"w!\xf06\x02W\xbd\x07\xc0U\xf8\x01\x1a\x05\xf6 " +
	"p\x00\x81\x89\xff\x0d\xc0D\xbc\x01v\xa3/\xeeC\xe0" +
	"\x10\x02\x93\x8e\x000\x09?\xe1\xa2\xfa\x1e@\xe0s\x04" +
	"\xae>\x0c\xc0\xd5\xf8\x05\x01\x05N#\x90\x9e\x06\xc05" +
	"\x87\x00\xb8\x06g0\x0d\xb4jI\xc3\x19L3\"\xbf" +
	"s\xd92\x88\x16\xc6%\xad\x1e\xe0\xa7O3\x8e\x15z" +
	"j\xaa\xf1\xc8\x12~\x09\x1b\x00\xbaV\xe9\xb0\x01\xa0k" +
	"5\x0d\x1b\x80\x84\xa0\x95\x16\xce\xc0\x0c\x08\x0a\x0e\xf8%" +
	"l\x0c\xd83\xe09\x03\x9f\x81];Ccc\xc0\xce" +
	"\xc23\x0b\xec\xb1\x19z\xe9\xe1\xe9\x94\xa4\xa0^\x17Y" +
	"\xcfT\x01\x09J\x9d\xfa\x16\xbe\x16\x94\xb3\x9c\x03\xe9'" +
	"%\xa0\xa6\x85:N\xa3\x06\xe2x\xdd:5\x8e\xd7\xa5" +
	"S\xab\xa6[\xa8N\x95\xea\x0f\xcc\xb0\x10\x1d\x1ak," +
	"Nl\x86N\x8d\x13\x9b\xaeS\xe3\xc4\xb2\x9a\xd8\x98U" +
	"l\x9aJ\xf4,\x8b;\xdf\xb2N\x0a\x10\x87\x10\xb50" +
	"\xa4\xe2\xf3/\xc3\x921)\x05\xa8t\xbcS0\xbe\x00" +
	"H\xc8%L\xfc\xf1Z\xc2Q\x93\xc9\x06\x84x\x14\x85" +
	"D,\x97\xb3\x00\x1b_\x04h0=\xc7\xea\xa4\xb7\xba" +
	"\xc9\x95\xd5\x90\xa0\x16\x0a\x09\xfb\x1f\x0fN\xfc?\xa0." +
	"\xa1\xf6\x80\xa1Y.\xc4\x12\xef@T&!\x91\xc9j" +
	")\x9c\x01tP\x8b\x8dH\x8al\xaf_\xa2\x0a\xce\xa4" +
	"\xc3\xcc\x0a\xb3\"\xf4D\xe5H\xca\xdd\xcf%\x8e\xab\x9b" +
	"\x04\x0f^~\xa58T\x81:\x80\x08\x7f\xcf\xbe\xac\xcc" +
	"r\x90\x9e\xe2\xf8#\xf9TU\xbb\x13\xee'\x89ck" +
	"\xd3\xce\x96f@%\x13\x06\x9e\xb9\x82\x1ca\x9c\x81\x81" +
	"$\x87\xbd\xd4\xc1\xadq`M\x84\x8b\xdc7Y\x14\xfe" +
	"\xff\x00\x00\x00\xff\xff\xa5<\xb9\x0a"

var x_832bcc6686a26d56 = []byte{
	0, 0, 0, 0, 2, 0, 0, 0,
	0, 0, 0, 0, 1, 0, 0, 0,
	223, 7, 8, 27, 0, 0, 0, 0,
	0, 0, 0, 0, 4, 0, 0, 0,
	1, 0, 0, 0, 23, 0, 0, 0,
	8, 0, 0, 0, 1, 0, 0, 0,
	223, 7, 8, 27, 0, 0, 0, 0,
	223, 7, 8, 28, 0, 0, 0, 0,
	0, 0, 0, 0, 3, 0, 0, 0,
	0, 0, 0, 0, 1, 0, 1, 0,
	42, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 3, 0, 0, 0,
	0, 0, 0, 0, 1, 0, 1, 0,
	42, 0, 0, 0, 0, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0,
}
