.PHONY: all clean

all: test.capnp.go

test.capnp.go: test.capnp
	capnp compile -ogo test.capnp

clean:
	rm -f test.capnp.go
