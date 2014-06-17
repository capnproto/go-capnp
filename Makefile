.PHONY: prepare

prepare:
	cd capnpc-go; go build
	go install ./capnpc-go
	cd aircraftlib; make
	which capnpc-go
	diff `which capnpc-go` ./capnpc-go/capnpc-go
	# if there was a diff above, adjust your PATH to use the most-recently build capnpc-go


check:
	cat data/check.zdate.cpz | capnp decode aircraftlib/aircraft.capnp  Zdate 

checkp:
	cat data/zdate2.packed.dat | bin/decp

testbuild:
	go test -c -gcflags "-N -l" -v

clean:
	rm -f go-capnproto.test *~
	cd aircraftlib; make clean

test:
	cd capnpc-go; go build; go install
	cd aircraftlib; make
	go test -v

