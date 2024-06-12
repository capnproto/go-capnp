package aircraftlib

import (
	math "math"

	capnp "capnproto.org/go/capnp/v3"
)

// Experimental: This could replace NewRootBenchmarkA without having to modify
// the public API.
func AllocateNewRootBenchmark(msg *capnp.Message) (BenchmarkA, error) {
	st, err := capnp.AllocateRootStruct(msg, capnp.ObjectSize{DataSize: 24, PointerCount: 2})
	return BenchmarkA(st), err

}

// Exprimental: set the name using the flat/unrolled version of SetNewText.
//
// If the unrolled version is deemed good enough, it would just be replaced
// inside SetNewText, without having to alter the public API.
func (s BenchmarkA) FlatSetName(v string) error {
	return capnp.Struct(s).FlatSetNewText(0, v)
}

// Experimental: update the set in-place.
func (s BenchmarkA) UpdateName(v string) error {
	return capnp.Struct(s).UpdateText(0, v)
}

// Experimental: return the name as a field that can be mutated.
func (s BenchmarkA) NameField() (capnp.TextField, error) {
	return capnp.Struct(s).TextField(0)
}

func (s BenchmarkA) FlatSetPhone(v string) error {
	return capnp.Struct(s).FlatSetNewText(1, v)
}

func (s BenchmarkA) PhoneField() (capnp.TextField, error) {
	return capnp.Struct(s).TextField(1)
}

func (s *BenchmarkA) SetMoneyp(v float64) {
	bits := math.Float64bits(v)
	(*capnp.Struct)(s).SetUint64p(16, bits)
}

func (s *BenchmarkA) SetSpousep(v bool) {
	(*capnp.Struct)(s).SetBitp(96, v)
}

func (s *BenchmarkA) SetSiblingsp(v int32) {
	(*capnp.Struct)(s).SetUint32p(8, uint32(v))
}

func (s *BenchmarkA) SetBirthDayp(v int64) {
	(*capnp.Struct)(s).SetUint64p(0, uint64(v))
}
