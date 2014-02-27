.PHONY: prepare

prepare: test.capnp
	go install ./capnpc-go
	capnp compile -ogo test.capnp
	mv test.capnp.go messages_test.go
