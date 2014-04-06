all:
	capnp compile -ogo aircraft.capnp

aircraft.capnp.go: aircraft.capnp
	capnp compile -ogo aircraft.capnp

clean:
	rm -f aircraft.capnp.go *~

