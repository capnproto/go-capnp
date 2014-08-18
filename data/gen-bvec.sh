echo "(boolvec = [true, false, true])" | capnp encode ../aircraftlib/aircraft.capnp Z > bvec5.bin
echo "(boolvec = [true, true])" | capnp encode ../aircraftlib/aircraft.capnp Z > bvec3.bin
echo "(boolvec = [true])" | capnp encode ../aircraftlib/aircraft.capnp Z > bvec1.bin
echo "(boolvec = [false])" | capnp encode ../aircraftlib/aircraft.capnp Z > bvec0.bin
