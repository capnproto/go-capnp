.PHONY: prepare

prepare: test.capnp
	go install ./capnpc-go
	capnp compile -ogo test.capnp
	mv test.capnp.go messages_test.go


check:
	cat check.zdate.cpz | capnp decode test.capnp  Zdate 

testbuild:
	go test -c -gcflags "-N -l" -v
