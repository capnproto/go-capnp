.PHONY: all check

all:
	( cd $(GOPATH)/src/github.com/kaos/capnp_test && capnp compile -ogo test.capnp )
	go install

check: all
	CAPNP_TEST_APP="$(GOPATH)/bin/go_capnp_test" make -C $(GOPATH)/src/github.com/kaos/capnp_test
