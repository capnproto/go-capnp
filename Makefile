.PHONY: prepare

prepare: test.capnp
	go install ./capnpc-go
	capnp compile -ogo test.capnp
	mv test.capnp.go messages_test.go

check:
	cat data/check.zdate.cpz | capnp decode aircraftlib/aircraft.capnp  Zdate 

checkp:
	cat data/zdate2.packed.dat | bin/decp

testbuild:
	go test -c -gcflags "-N -l" -v

clean:
	rm -f go-capnproto.test *~
