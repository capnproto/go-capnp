package aircraftlib

import capnp "capnproto.org/go/capnp/v3"

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
