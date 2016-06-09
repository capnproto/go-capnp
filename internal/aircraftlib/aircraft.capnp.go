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
var schema_832bcc6686a26d56 = []byte{
	120, 218, 172, 90, 13, 116, 84, 69,
	150, 126, 213, 63, 121, 157, 64, 236,
	126, 169, 252, 64, 12, 182, 68, 242,
	59, 7, 201, 15, 131, 200, 204, 158,
	144, 152, 40, 122, 192, 201, 15, 136,
	100, 141, 230, 37, 121, 73, 26, 186,
	251, 133, 238, 215, 144, 224, 120, 226,
	204, 128, 34, 71, 118, 116, 213, 81,
	81, 119, 52, 139, 103, 101, 52, 174,
	56, 48, 71, 89, 101, 196, 129, 17,
	34, 236, 136, 7, 20, 88, 65, 192,
	69, 23, 196, 81, 28, 89, 65, 133,
	183, 247, 214, 235, 247, 147, 238, 215,
	176, 112, 230, 28, 66, 94, 221, 239,
	171, 91, 85, 183, 170, 238, 189, 85,
	149, 138, 191, 142, 157, 233, 168, 116,
	223, 68, 57, 174, 233, 176, 59, 77,
	157, 56, 229, 204, 149, 63, 172, 104,
	59, 198, 9, 94, 135, 122, 91, 104,
	232, 190, 238, 157, 63, 250, 21, 199,
	17, 250, 64, 250, 38, 250, 112, 58,
	207, 113, 116, 117, 250, 191, 115, 68,
	157, 90, 84, 122, 124, 75, 97, 203,
	63, 1, 209, 57, 138, 120, 36, 125,
	136, 30, 103, 196, 99, 233, 55, 209,
	204, 12, 248, 82, 31, 218, 238, 255,
	115, 237, 221, 47, 255, 26, 200, 196,
	36, 187, 29, 200, 58, 157, 126, 148,
	18, 100, 209, 115, 233, 75, 65, 241,
	246, 121, 71, 221, 155, 74, 126, 253,
	20, 112, 93, 163, 20, 75, 25, 35,
	116, 49, 18, 91, 130, 25, 78, 210,
	210, 159, 225, 32, 160, 122, 81, 184,
	61, 187, 230, 252, 230, 199, 18, 84,
	19, 224, 85, 7, 50, 178, 8, 29,
	96, 186, 99, 25, 53, 160, 123, 245,
	187, 251, 151, 174, 108, 127, 112, 191,
	13, 153, 62, 159, 49, 68, 135, 25,
	119, 29, 227, 26, 48, 73, 231, 28,
	116, 71, 198, 178, 234, 93, 25, 188,
	147, 110, 204, 188, 142, 115, 168, 5,
	91, 58, 61, 191, 201, 47, 63, 196,
	53, 121, 137, 51, 113, 68, 167, 50,
	223, 167, 231, 50, 177, 3, 103, 51,
	231, 19, 80, 53, 118, 207, 111, 182,
	102, 189, 88, 241, 104, 210, 144, 230,
	120, 71, 232, 2, 47, 14, 105, 174,
	23, 134, 212, 231, 101, 67, 90, 255,
	206, 188, 67, 47, 251, 110, 61, 134,
	202, 45, 253, 104, 32, 188, 11, 180,
	135, 188, 111, 211, 24, 214, 161, 139,
	189, 159, 129, 242, 199, 175, 218, 61,
	165, 241, 172, 180, 150, 107, 42, 32,
	80, 219, 133, 237, 54, 248, 34, 208,
	46, 109, 242, 225, 64, 22, 142, 11,
	222, 148, 93, 175, 126, 107, 55, 232,
	197, 190, 247, 233, 61, 62, 252, 26,
	96, 220, 205, 7, 86, 158, 90, 251,
	253, 180, 15, 236, 184, 111, 250, 214,
	208, 109, 140, 187, 133, 113, 15, 205,
	222, 252, 72, 169, 242, 175, 63, 36,
	116, 19, 59, 64, 143, 128, 222, 47,
	24, 247, 56, 227, 46, 120, 201, 241,
	212, 227, 79, 110, 92, 99, 167, 55,
	93, 24, 161, 57, 2, 126, 9, 2,
	155, 164, 137, 158, 169, 103, 238, 249,
	227, 163, 118, 220, 235, 129, 219, 192,
	184, 181, 140, 27, 251, 137, 235, 62,
	181, 170, 226, 153, 132, 62, 104, 100,
	17, 200, 33, 70, 14, 8, 184, 178,
	78, 221, 51, 167, 188, 110, 222, 123,
	47, 218, 45, 149, 109, 66, 62, 161,
	123, 25, 121, 55, 211, 124, 115, 227,
	209, 253, 71, 95, 174, 255, 207, 196,
	25, 118, 34, 133, 100, 157, 160, 153,
	89, 88, 47, 61, 139, 205, 112, 219,
	187, 190, 59, 211, 223, 152, 117, 127,
	98, 63, 216, 122, 144, 232, 38, 26,
	162, 172, 31, 20, 183, 206, 234, 251,
	253, 127, 158, 245, 192, 39, 143, 35,
	217, 145, 72, 206, 201, 126, 155, 78,
	200, 70, 213, 227, 179, 253, 168, 250,
	137, 241, 89, 171, 190, 249, 229, 253,
	167, 57, 161, 64, 159, 222, 202, 156,
	133, 132, 115, 169, 5, 158, 182, 186,
	244, 225, 159, 157, 178, 51, 212, 132,
	156, 3, 180, 44, 7, 191, 138, 114,
	112, 56, 31, 165, 157, 253, 237, 242,
	193, 95, 28, 75, 224, 178, 209, 204,
	203, 25, 161, 34, 227, 182, 229, 96,
	255, 102, 244, 207, 156, 246, 234, 174,
	137, 43, 19, 6, 83, 207, 59, 80,
	93, 238, 8, 173, 204, 69, 242, 228,
	92, 52, 234, 208, 223, 174, 252, 112,
	195, 93, 185, 171, 56, 33, 27, 186,
	167, 25, 243, 129, 92, 7, 46, 191,
	213, 185, 216, 242, 51, 119, 31, 124,
	189, 236, 195, 235, 158, 178, 18, 134,
	115, 51, 144, 176, 145, 17, 94, 255,
	252, 21, 185, 225, 192, 253, 195, 118,
	195, 216, 155, 59, 68, 15, 178, 214,
	246, 51, 238, 225, 226, 233, 237, 159,
	188, 250, 251, 157, 118, 220, 179, 192,
	37, 121, 204, 145, 48, 238, 166, 111,
	63, 222, 123, 215, 146, 143, 222, 181,
	155, 238, 9, 121, 48, 221, 147, 25,
	185, 44, 15, 201, 98, 75, 117, 214,
	182, 19, 59, 142, 216, 41, 94, 144,
	183, 134, 138, 140, 219, 198, 184, 67,
	99, 190, 30, 58, 43, 29, 120, 199,
	142, 123, 79, 222, 239, 232, 10, 198,
	253, 5, 227, 62, 82, 241, 194, 119,
	63, 221, 243, 47, 135, 236, 58, 241,
	92, 94, 6, 161, 235, 25, 121, 152,
	145, 159, 216, 186, 110, 204, 167, 66,
	189, 237, 232, 14, 230, 109, 162, 199,
	24, 247, 8, 227, 158, 89, 240, 236,
	175, 158, 30, 242, 28, 179, 83, 156,
	51, 14, 252, 94, 209, 56, 36, 79,
	28, 135, 228, 193, 105, 119, 46, 111,
	155, 254, 213, 122, 219, 197, 220, 52,
	238, 125, 218, 134, 228, 234, 5, 227,
	216, 98, 142, 142, 76, 85, 79, 30,
	189, 234, 15, 118, 235, 179, 122, 219,
	120, 7, 161, 187, 199, 99, 197, 93,
	227, 209, 255, 44, 218, 187, 121, 237,
	112, 254, 226, 207, 146, 156, 219, 115,
	249, 35, 116, 56, 31, 157, 219, 11,
	249, 224, 220, 54, 228, 51, 231, 246,
	240, 117, 255, 214, 30, 254, 203, 27,
	251, 80, 185, 43, 113, 241, 175, 131,
	42, 27, 177, 74, 245, 250, 124, 214,
	149, 171, 183, 231, 245, 223, 245, 241,
	203, 47, 38, 69, 153, 156, 130, 163,
	116, 98, 1, 91, 233, 5, 55, 209,
	155, 241, 75, 61, 252, 163, 129, 107,
	124, 131, 235, 54, 115, 77, 233, 110,
	162, 174, 12, 30, 60, 223, 84, 80,
	190, 27, 217, 149, 5, 171, 232, 245,
	200, 105, 153, 90, 224, 68, 189, 251,
	159, 253, 224, 39, 143, 31, 47, 62,
	145, 48, 68, 112, 177, 165, 184, 198,
	129, 62, 25, 233, 213, 101, 5, 239,
	120, 128, 174, 70, 118, 158, 104, 122,
	114, 199, 35, 163, 141, 205, 60, 221,
	99, 197, 111, 211, 103, 138, 121, 206,
	169, 238, 251, 253, 249, 138, 107, 106,
	254, 116, 191, 221, 140, 12, 20, 195,
	140, 60, 80, 140, 21, 86, 20, 227,
	140, 108, 191, 119, 75, 193, 200, 137,
	167, 254, 98, 71, 222, 136, 228, 109,
	140, 188, 133, 145, 135, 55, 243, 194,
	222, 221, 67, 182, 139, 243, 84, 241,
	38, 122, 150, 113, 79, 51, 110, 247,
	172, 158, 175, 207, 211, 63, 190, 96,
	235, 20, 74, 222, 166, 69, 37, 108,
	89, 148, 32, 119, 219, 153, 47, 62,
	168, 184, 187, 104, 133, 157, 7, 175,
	45, 25, 161, 115, 24, 247, 102, 198,
	221, 88, 127, 69, 49, 249, 67, 197,
	145, 100, 227, 6, 74, 126, 73, 67,
	200, 108, 233, 45, 97, 198, 165, 75,
	190, 15, 116, 215, 238, 62, 104, 23,
	237, 23, 148, 12, 81, 145, 169, 109,
	43, 65, 247, 97, 40, 34, 78, 140,
	178, 37, 183, 208, 93, 37, 75, 105,
	81, 169, 31, 130, 236, 51, 115, 231,
	15, 255, 199, 203, 141, 135, 236, 186,
	87, 84, 250, 59, 58, 185, 148, 237,
	223, 82, 244, 89, 11, 122, 243, 191,
	156, 241, 249, 242, 79, 237, 134, 189,
	165, 244, 0, 221, 197, 184, 59, 74,
	113, 40, 66, 95, 215, 135, 97, 247,
	75, 235, 109, 205, 89, 250, 53, 61,
	199, 184, 103, 25, 247, 92, 247, 77,
	59, 26, 62, 114, 127, 155, 176, 80,
	88, 31, 202, 202, 222, 167, 63, 46,
	195, 175, 202, 50, 28, 203, 228, 57,
	215, 191, 249, 201, 139, 255, 184, 195,
	110, 220, 171, 203, 70, 232, 147, 140,
	251, 24, 227, 138, 129, 72, 103, 68,
	236, 86, 200, 181, 157, 98, 95, 184,
	111, 70, 109, 77, 32, 210, 39, 71,
	148, 70, 66, 32, 142, 131, 163, 21,
	230, 148, 131, 113, 137, 208, 80, 8,
	191, 28, 194, 63, 224, 47, 167, 240,
	99, 252, 229, 18, 38, 227, 47, 183,
	80, 132, 191, 210, 132, 9, 192, 244,
	134, 229, 176, 196, 47, 236, 94, 196,
	7, 197, 126, 62, 218, 45, 243, 193,
	216, 18, 190, 171, 123, 169, 87, 145,
	162, 74, 82, 115, 13, 124, 103, 175,
	140, 109, 185, 156, 110, 216, 57, 186,
	19, 39, 186, 179, 22, 132, 114, 206,
	33, 184, 121, 175, 4, 188, 153, 4,
	152, 134, 10, 71, 92, 69, 139, 34,
	118, 46, 10, 132, 123, 154, 121, 89,
	102, 221, 246, 56, 93, 132, 192, 63,
	34, 148, 45, 132, 108, 178, 212, 73,
	154, 166, 58, 136, 143, 100, 99, 90,
	34, 92, 159, 5, 178, 169, 32, 155,
	233, 0, 93, 243, 3, 74, 111, 189,
	212, 205, 121, 197, 88, 80, 33, 62,
	51, 142, 195, 144, 125, 28, 7, 234,
	8, 14, 95, 76, 130, 146, 251, 209,
	41, 135, 163, 74, 67, 56, 22, 210,
	170, 121, 205, 20, 22, 76, 228, 37,
	150, 26, 206, 120, 141, 89, 114, 176,
	43, 122, 155, 20, 153, 187, 84, 134,
	127, 179, 3, 81, 162, 104, 166, 128,
	244, 202, 133, 125, 205, 156, 1, 125,
	245, 64, 95, 39, 57, 72, 77, 104,
	32, 24, 136, 42, 228, 10, 142, 52,
	58, 161, 7, 166, 119, 132, 254, 92,
	97, 55, 147, 124, 117, 85, 69, 130,
	190, 242, 184, 190, 108, 7, 241, 118,
	136, 81, 9, 212, 24, 126, 48, 97,
	88, 92, 141, 166, 167, 233, 45, 98,
	241, 173, 116, 128, 52, 155, 105, 36,
	43, 25, 121, 50, 43, 25, 123, 5,
	74, 85, 102, 132, 162, 49, 40, 25,
	6, 161, 139, 73, 157, 217, 48, 13,
	65, 61, 35, 160, 66, 169, 220, 204,
	143, 105, 0, 74, 70, 76, 164, 18,
	41, 52, 99, 8, 21, 73, 171, 153,
	163, 66, 233, 22, 211, 157, 66, 41,
	203, 76, 134, 104, 27, 180, 103, 164,
	148, 116, 1, 104, 49, 66, 28, 157,
	7, 152, 225, 41, 104, 19, 180, 103,
	120, 89, 58, 7, 116, 26, 233, 37,
	148, 90, 205, 125, 200, 74, 70, 138,
	8, 165, 102, 115, 223, 177, 146, 49,
	63, 80, 90, 101, 198, 95, 104, 225,
	159, 205, 12, 16, 90, 31, 50, 19,
	4, 232, 217, 144, 25, 80, 161, 215,
	107, 76, 255, 12, 35, 90, 99, 30,
	50, 192, 18, 107, 76, 71, 15, 86,
	90, 99, 38, 138, 96, 193, 136, 153,
	103, 49, 235, 26, 105, 47, 43, 25,
	110, 28, 74, 117, 166, 199, 2, 45,
	29, 102, 206, 5, 165, 102, 51, 179,
	99, 152, 145, 192, 67, 169, 213, 244,
	94, 80, 90, 102, 30, 195, 216, 140,
	25, 169, 20, 244, 179, 220, 12, 5,
	108, 142, 140, 51, 24, 148, 22, 154,
	187, 10, 74, 205, 102, 32, 96, 37,
	35, 232, 50, 166, 145, 173, 50, 45,
	198, 193, 135, 173, 2, 182, 247, 234,
	69, 133, 35, 146, 246, 13, 155, 137,
	35, 138, 170, 239, 73, 142, 132, 252,
	173, 93, 162, 34, 177, 255, 197, 193,
	90, 205, 201, 169, 141, 65, 49, 44,
	213, 137, 192, 149, 188, 117, 215, 85,
	95, 231, 173, 133, 61, 195, 223, 88,
	57, 77, 109, 150, 122, 34, 82, 52,
	26, 224, 156, 114, 88, 173, 213, 247,
	4, 108, 146, 214, 193, 27, 228, 88,
	88, 145, 34, 124, 157, 216, 51, 216,
	26, 149, 34, 75, 164, 136, 183, 117,
	161, 220, 161, 194, 86, 110, 8, 245,
	41, 3, 64, 195, 239, 159, 133, 165,
	122, 145, 115, 42, 162, 170, 237, 113,
	179, 0, 72, 35, 244, 53, 18, 7,
	172, 223, 245, 34, 81, 68, 38, 139,
	128, 22, 221, 67, 16, 166, 151, 121,
	8, 83, 198, 244, 43, 226, 108, 30,
	220, 130, 41, 101, 58, 146, 164, 172,
	197, 200, 236, 128, 51, 129, 106, 43,
	100, 206, 200, 70, 216, 24, 116, 198,
	162, 170, 238, 176, 56, 63, 8, 160,
	204, 72, 115, 165, 126, 180, 247, 252,
	136, 216, 135, 93, 229, 200, 192, 32,
	126, 87, 245, 87, 169, 241, 223, 125,
	28, 143, 236, 219, 228, 64, 215, 188,
	112, 64, 230, 72, 88, 189, 21, 34,
	130, 20, 169, 188, 129, 227, 193, 209,
	168, 205, 243, 231, 66, 249, 6, 48,
	18, 20, 96, 168, 74, 139, 18, 137,
	113, 53, 157, 32, 234, 11, 123, 27,
	192, 255, 123, 103, 201, 74, 175, 138,
	95, 48, 101, 18, 90, 89, 247, 253,
	156, 183, 25, 156, 191, 89, 36, 181,
	150, 239, 58, 245, 6, 49, 24, 108,
	145, 22, 199, 56, 175, 20, 238, 148,
	84, 112, 248, 232, 238, 163, 168, 161,
	14, 36, 189, 33, 49, 194, 57, 23,
	213, 54, 121, 32, 180, 154, 201, 155,
	207, 65, 106, 175, 182, 38, 28, 5,
	32, 152, 73, 200, 88, 240, 164, 176,
	27, 12, 71, 201, 7, 3, 29, 176,
	48, 64, 92, 73, 86, 17, 117, 153,
	28, 234, 8, 72, 203, 36, 119, 248,
	218, 78, 57, 52, 165, 71, 158, 194,
	252, 104, 68, 86, 228, 170, 41, 1,
	92, 57, 97, 49, 56, 69, 175, 141,
	117, 147, 99, 137, 222, 171, 69, 164,
	22, 61, 248, 56, 195, 131, 63, 137,
	30, 252, 81, 240, 224, 207, 58, 32,
	198, 104, 17, 237, 153, 91, 64, 246,
	52, 200, 94, 112, 16, 193, 1, 66,
	140, 218, 207, 87, 129, 240, 89, 16,
	190, 4, 66, 167, 35, 27, 210, 26,
	78, 88, 135, 204, 23, 64, 184, 1,
	132, 174, 246, 108, 2, 106, 133, 245,
	24, 100, 94, 2, 225, 107, 32, 116,
	3, 19, 2, 177, 176, 17, 171, 191,
	2, 194, 55, 32, 82, 132, 197, 144,
	4, 195, 115, 192, 15, 81, 59, 2,
	17, 8, 153, 34, 46, 114, 96, 58,
	224, 135, 248, 251, 122, 33, 232, 27,
	140, 104, 160, 35, 8, 134, 71, 243,
	66, 3, 14, 248, 33, 53, 209, 62,
	57, 6, 225, 134, 64, 17, 34, 161,
	63, 4, 252, 1, 50, 6, 74, 99,
	82, 69, 82, 92, 1, 90, 36, 53,
	34, 158, 17, 91, 88, 196, 131, 73,
	200, 19, 60, 196, 251, 49, 239, 185,
	18, 255, 43, 72, 210, 99, 108, 91,
	180, 161, 207, 233, 26, 171, 170, 204,
	136, 34, 26, 241, 14, 24, 92, 175,
	131, 100, 146, 243, 170, 102, 70, 9,
	165, 237, 32, 13, 130, 212, 113, 78,
	213, 236, 24, 64, 105, 23, 72, 251,
	64, 234, 252, 65, 213, 12, 25, 130,
	132, 167, 169, 23, 164, 10, 152, 103,
	9, 44, 105, 46, 205, 219, 1, 254,
	3, 58, 105, 4, 50, 45, 158, 122,
	69, 240, 40, 32, 54, 34, 154, 38,
	230, 187, 43, 167, 129, 212, 136, 108,
	9, 193, 215, 21, 31, 129, 190, 112,
	113, 217, 94, 219, 35, 41, 183, 198,
	66, 29, 82, 100, 82, 179, 228, 143,
	226, 2, 182, 134, 246, 44, 51, 180,
	147, 48, 241, 128, 105, 61, 54, 166,
	213, 55, 89, 31, 9, 39, 164, 6,
	173, 80, 127, 44, 212, 159, 14, 105,
	81, 24, 56, 115, 68, 240, 64, 206,
	64, 191, 110, 126, 99, 22, 140, 152,
	144, 34, 239, 168, 115, 138, 61, 9,
	186, 235, 204, 190, 13, 118, 106, 222,
	19, 20, 25, 225, 57, 69, 66, 165,
	187, 80, 133, 136, 9, 250, 10, 77,
	125, 252, 18, 49, 136, 41, 59, 252,
	36, 107, 208, 221, 209, 0, 91, 72,
	86, 13, 175, 130, 6, 31, 104, 128,
	45, 173, 134, 2, 61, 189, 202, 173,
	178, 66, 234, 164, 102, 9, 12, 62,
	224, 103, 117, 160, 135, 70, 34, 112,
	193, 30, 50, 87, 125, 129, 30, 246,
	41, 56, 90, 35, 141, 72, 161, 75,
	247, 87, 181, 154, 46, 143, 161, 171,
	12, 117, 77, 2, 93, 21, 230, 150,
	159, 156, 101, 38, 182, 60, 164, 156,
	250, 62, 35, 29, 208, 146, 17, 68,
	19, 90, 74, 76, 59, 245, 248, 129,
	174, 254, 18, 242, 78, 235, 64, 174,
	176, 55, 138, 22, 46, 188, 24, 30,
	18, 220, 87, 161, 233, 190, 4, 195,
	127, 161, 240, 9, 16, 174, 133, 1,
	58, 180, 109, 247, 92, 185, 197, 167,
	57, 137, 182, 235, 158, 47, 183, 248,
	52, 151, 67, 115, 95, 235, 176, 246,
	90, 16, 190, 98, 113, 95, 195, 229,
	113, 71, 183, 117, 244, 10, 225, 187,
	98, 178, 238, 181, 188, 48, 47, 149,
	9, 121, 180, 79, 19, 87, 37, 139,
	121, 37, 34, 25, 53, 131, 81, 165,
	90, 55, 8, 202, 236, 182, 193, 13,
	53, 218, 66, 71, 3, 140, 53, 12,
	208, 128, 29, 155, 9, 29, 155, 109,
	78, 230, 205, 232, 107, 235, 65, 214,
	104, 241, 223, 115, 208, 85, 207, 214,
	124, 148, 55, 26, 88, 102, 52, 238,
	95, 42, 71, 186, 162, 134, 179, 197,
	82, 144, 185, 74, 78, 239, 209, 216,
	132, 30, 233, 19, 99, 132, 61, 236,
	82, 182, 214, 37, 56, 36, 221, 131,
	93, 234, 135, 150, 150, 235, 115, 2,
	194, 21, 40, 188, 23, 132, 15, 198,
	231, 4, 100, 171, 177, 159, 43, 65,
	246, 168, 62, 39, 32, 124, 24, 205,
	255, 32, 8, 159, 208, 231, 4, 132,
	143, 97, 237, 135, 64, 248, 52, 116,
	94, 145, 250, 149, 120, 119, 193, 101,
	23, 242, 221, 178, 236, 197, 36, 140,
	100, 130, 44, 19, 101, 249, 124, 135,
	24, 241, 119, 7, 101, 216, 234, 25,
	156, 227, 84, 198, 159, 78, 207, 154,
	9, 38, 135, 96, 137, 11, 251, 148,
	107, 133, 170, 170, 96, 246, 24, 10,
	192, 179, 9, 158, 242, 255, 191, 175,
	108, 20, 189, 17, 49, 20, 77, 178,
	134, 158, 124, 136, 144, 110, 164, 56,
	37, 77, 66, 211, 43, 145, 232, 197,
	12, 171, 103, 64, 90, 188, 178, 206,
	118, 161, 57, 219, 198, 114, 191, 185,
	48, 62, 221, 237, 56, 221, 241, 245,
	222, 134, 59, 238, 118, 109, 18, 120,
	197, 48, 24, 225, 131, 230, 214, 139,
	55, 95, 3, 34, 139, 244, 34, 125,
	211, 243, 173, 176, 225, 86, 244, 40,
	88, 150, 101, 250, 21, 51, 10, 90,
	61, 11, 17, 185, 52, 210, 193, 165,
	37, 121, 16, 76, 193, 174, 197, 19,
	57, 88, 23, 141, 203, 89, 173, 151,
	111, 186, 62, 103, 32, 108, 44, 212,
	212, 58, 154, 165, 168, 150, 141, 165,
	114, 160, 114, 76, 73, 82, 163, 111,
	179, 89, 60, 228, 132, 23, 63, 229,
	26, 71, 144, 4, 159, 104, 4, 45,
	30, 162, 247, 229, 31, 150, 109, 78,
	244, 151, 227, 90, 141, 195, 101, 138,
	208, 122, 163, 179, 114, 218, 229, 119,
	82, 215, 210, 202, 182, 223, 5, 244,
	88, 118, 103, 234, 33, 234, 39, 18,
	37, 41, 246, 93, 104, 132, 214, 136,
	106, 55, 194, 214, 26, 237, 64, 149,
	160, 178, 35, 158, 156, 148, 66, 184,
	94, 42, 6, 20, 136, 147, 11, 57,
	94, 238, 136, 154, 154, 141, 35, 124,
	130, 102, 155, 168, 199, 14, 61, 151,
	118, 217, 98, 28, 244, 83, 4, 61,
	253, 184, 40, 107, 121, 213, 69, 18,
	246, 124, 51, 224, 25, 14, 223, 136,
	120, 175, 89, 18, 246, 141, 51, 226,
	105, 248, 78, 244, 174, 68, 139, 120,
	59, 112, 111, 108, 5, 225, 123, 150,
	136, 183, 11, 133, 219, 65, 184, 39,
	229, 66, 112, 118, 84, 232, 41, 183,
	183, 67, 130, 41, 142, 143, 110, 76,
	220, 175, 244, 225, 97, 216, 98, 81,
	227, 106, 69, 27, 51, 63, 16, 138,
	233, 245, 249, 129, 104, 87, 82, 250,
	238, 28, 149, 0, 224, 250, 143, 159,
	99, 227, 233, 119, 220, 34, 98, 161,
	153, 125, 27, 78, 81, 42, 52, 147,
	111, 61, 7, 176, 166, 222, 70, 14,
	16, 42, 55, 51, 239, 139, 69, 246,
	132, 148, 75, 143, 236, 23, 203, 196,
	244, 203, 4, 73, 59, 126, 140, 62,
	119, 248, 56, 206, 167, 157, 57, 146,
	234, 233, 151, 9, 90, 189, 20, 107,
	192, 204, 122, 170, 226, 139, 96, 131,
	57, 226, 81, 71, 49, 231, 204, 132,
	53, 240, 134, 37, 235, 121, 29, 211,
	131, 215, 180, 4, 71, 112, 59, 181,
	53, 176, 5, 133, 111, 105, 171, 101,
	212, 161, 205, 223, 43, 135, 204, 153,
	29, 117, 51, 201, 102, 62, 34, 226,
	142, 210, 205, 87, 211, 41, 134, 111,
	12, 14, 232, 135, 53, 21, 70, 39,
	118, 6, 20, 235, 137, 79, 13, 137,
	253, 45, 125, 146, 212, 133, 178, 84,
	199, 56, 35, 32, 243, 16, 145, 205,
	59, 94, 253, 29, 145, 232, 15, 198,
	130, 208, 12, 49, 61, 157, 87, 245,
	160, 205, 145, 136, 118, 213, 219, 35,
	107, 170, 200, 12, 232, 193, 34, 177,
	39, 62, 35, 48, 170, 100, 199, 161,
	25, 253, 89, 35, 190, 209, 38, 23,
	152, 189, 101, 182, 203, 73, 90, 110,
	119, 89, 66, 28, 157, 231, 130, 253,
	215, 210, 136, 192, 29, 0, 76, 128,
	179, 158, 182, 7, 233, 2, 87, 33,
	62, 129, 35, 210, 142, 8, 156, 247,
	180, 141, 72, 219, 24, 114, 59, 34,
	93, 136, 184, 190, 87, 181, 221, 72,
	69, 134, 220, 129, 72, 47, 34, 238,
	239, 84, 109, 75, 82, 137, 33, 237,
	136, 4, 17, 73, 59, 171, 186, 178,
	73, 26, 62, 191, 50, 164, 11, 145,
	62, 68, 248, 51, 170, 39, 155, 93,
	251, 135, 88, 223, 122, 17, 81, 16,
	241, 124, 139, 237, 120, 240, 133, 156,
	213, 9, 34, 210, 143, 72, 250, 255,
	98, 59, 233, 248, 39, 5, 12, 233,
	67, 228, 231, 136, 100, 156, 198, 118,
	50, 240, 45, 157, 33, 10, 34, 247,
	34, 50, 230, 27, 108, 103, 12, 62,
	15, 178, 118, 250, 17, 89, 142, 200,
	216, 191, 169, 176, 220, 198, 226, 115,
	33, 51, 219, 207, 17, 89, 137, 102,
	203, 252, 26, 204, 150, 137, 239, 69,
	12, 184, 23, 129, 7, 17, 184, 226,
	20, 0, 112, 12, 167, 15, 48, 96,
	57, 2, 15, 33, 224, 253, 10, 0,
	47, 190, 53, 184, 96, 233, 130, 22,
	0, 214, 34, 224, 251, 18, 0, 216,
	64, 244, 57, 6, 60, 141, 192, 107,
	8, 8, 127, 5, 64, 0, 96, 35,
	3, 94, 65, 96, 39, 2, 89, 95,
	0, 144, 133, 15, 38, 12, 216, 138,
	192, 97, 4, 232, 73, 0, 40, 62,
	71, 50, 96, 31, 2, 223, 32, 144,
	253, 57, 0, 217, 248, 128, 226, 130,
	13, 214, 114, 18, 1, 143, 27, 128,
	156, 19, 0, 228, 0, 224, 118, 67,
	141, 102, 55, 200, 11, 80, 158, 123,
	28, 228, 185, 32, 31, 143, 242, 150,
	108, 4, 42, 16, 200, 251, 31, 0,
	242, 240, 125, 153, 1, 165, 8, 212,
	35, 48, 238, 51, 0, 198, 225, 211,
	148, 27, 155, 248, 41, 2, 183, 35,
	48, 254, 83, 0, 198, 227, 10, 115,
	163, 69, 26, 17, 8, 34, 144, 127,
	12, 128, 124, 156, 120, 119, 29, 78,
	60, 2, 15, 33, 112, 229, 127, 3,
	112, 37, 154, 138, 169, 90, 137, 192,
	163, 8, 20, 124, 2, 64, 1, 0,
	15, 51, 224, 65, 4, 158, 64, 96,
	194, 81, 0, 38, 224, 243, 141, 27,
	130, 34, 144, 1, 120, 5, 129, 171,
	142, 0, 112, 21, 62, 225, 186, 193,
	13, 180, 188, 132, 192, 107, 8, 248,
	15, 3, 224, 71, 227, 186, 97, 147,
	180, 108, 64, 224, 45, 4, 174, 254,
	24, 128, 171, 241, 15, 42, 220, 205,
	0, 188, 129, 192, 118, 0, 38, 76,
	60, 132, 11, 104, 34, 32, 219, 88,
	127, 223, 66, 100, 39, 86, 41, 60,
	8, 85, 10, 113, 62, 216, 8, 183,
	34, 240, 30, 2, 215, 124, 4, 192,
	53, 248, 22, 203, 128, 237, 8, 236,
	65, 96, 210, 127, 1, 48, 9, 255,
	154, 193, 141, 107, 113, 39, 2, 251,
	16, 40, 58, 0, 64, 17, 62, 173,
	179, 254, 238, 65, 224, 43, 4, 138,
	247, 3, 80, 12, 192, 23, 12, 56,
	137, 128, 39, 13, 128, 146, 125, 0,
	148, 224, 12, 166, 65, 175, 154, 211,
	112, 6, 211, 140, 251, 23, 231, 178,
	101, 224, 170, 141, 7, 0, 253, 154,
	101, 218, 84, 35, 102, 117, 87, 87,
	225, 241, 2, 126, 224, 104, 1, 242,
	184, 47, 227, 3, 32, 143, 159, 159,
	249, 64, 229, 52, 61, 162, 56, 3,
	211, 193, 41, 56, 224, 135, 240, 49,
	160, 167, 195, 119, 58, 126, 3, 61,
	126, 189, 194, 199, 128, 206, 195, 55,
	15, 244, 216, 116, 216, 219, 14, 248,
	129, 232, 42, 203, 65, 221, 125, 90,
	207, 63, 128, 4, 229, 14, 61, 181,
	170, 129, 206, 45, 145, 58, 19, 195,
	48, 116, 211, 34, 205, 136, 75, 3,
	163, 184, 110, 93, 58, 138, 235, 210,
	165, 149, 211, 44, 82, 167, 38, 245,
	7, 166, 91, 132, 142, 56, 53, 54,
	74, 109, 186, 46, 29, 165, 214, 163,
	75, 71, 169, 229, 227, 106, 99, 86,
	181, 105, 154, 208, 187, 204, 34, 27,
	61, 41, 32, 28, 68, 212, 66, 72,
	197, 243, 47, 195, 27, 251, 164, 248,
	171, 201, 241, 69, 206, 120, 93, 74,
	124, 186, 226, 71, 53, 159, 144, 199,
	152, 52, 16, 140, 70, 81, 73, 196,
	114, 241, 15, 176, 241, 218, 20, 135,
	89, 146, 212, 193, 94, 12, 146, 51,
	172, 65, 81, 123, 87, 72, 8, 177,
	94, 156, 248, 191, 195, 237, 32, 179,
	7, 12, 205, 114, 214, 79, 188, 20,
	213, 72, 98, 34, 201, 106, 41, 156,
	1, 92, 160, 22, 27, 145, 84, 135,
	71, 253, 65, 195, 250, 238, 106, 255,
	122, 217, 24, 140, 145, 232, 37, 36,
	212, 198, 107, 213, 69, 210, 117, 253,
	141, 226, 210, 210, 117, 227, 37, 46,
	197, 49, 99, 126, 141, 246, 40, 113,
	89, 183, 130, 9, 55, 67, 73, 86,
	51, 30, 40, 46, 112, 196, 194, 131,
	47, 17, 204, 87, 51, 124, 97, 191,
	192, 165, 96, 221, 5, 47, 24, 45,
	23, 128, 150, 148, 41, 16, 194, 165,
	152, 58, 97, 226, 23, 202, 29, 169,
	111, 26, 141, 60, 117, 114, 121, 252,
	66, 160, 30, 90, 234, 12, 117, 25,
	174, 76, 140, 244, 36, 221, 140, 232,
	173, 195, 226, 195, 191, 2, 52, 31,
	83, 132, 202, 58, 243, 33, 69, 152,
	60, 67, 189, 243, 145, 223, 54, 189,
	249, 193, 170, 109, 216, 170, 250, 206,
	151, 35, 147, 198, 111, 80, 158, 231,
	132, 162, 66, 53, 235, 112, 243, 201,
	129, 251, 150, 108, 231, 132, 137, 85,
	234, 35, 87, 79, 57, 188, 70, 242,
	125, 199, 9, 19, 90, 213, 83, 171,
	167, 228, 101, 181, 191, 254, 54, 20,
	202, 7, 227, 25, 97, 141, 54, 76,
	190, 75, 238, 228, 21, 177, 199, 31,
	150, 225, 127, 181, 51, 22, 85, 228,
	16, 44, 92, 103, 159, 196, 50, 225,
	38, 215, 168, 167, 29, 200, 20, 106,
	125, 241, 151, 156, 42, 127, 188, 203,
	246, 7, 101, 41, 245, 5, 158, 128,
	121, 123, 242, 13, 158, 51, 126, 131,
	135, 214, 156, 5, 194, 185, 48, 223,
	3, 146, 24, 209, 3, 11, 190, 126,
	40, 189, 122, 176, 224, 187, 196, 1,
	253, 219, 246, 2, 27, 223, 206, 240,
	225, 237, 178, 214, 170, 117, 155, 217,
	173, 85, 253, 157, 173, 211, 207, 158,
	217, 82, 44, 178, 73, 120, 214, 178,
	186, 212, 84, 111, 0, 142, 196, 35,
	160, 118, 117, 127, 177, 37, 86, 104,
	185, 205, 78, 113, 166, 75, 209, 132,
	113, 247, 110, 105, 162, 252, 66, 171,
	24, 154, 184, 180, 227, 225, 255, 5,
	0, 0, 255, 255, 125, 251, 180, 183,
}
