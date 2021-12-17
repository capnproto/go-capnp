module test

go 1.17

require (
	hashes v1.0.0
	capnproto.org/go/capnp/v3 v3.0.0-alpha.1
	golang.org/x/net v0.0.0-20211209124913-491a49abca63
)

replace (
	hashes v1.0.0 => ./hashes
)
