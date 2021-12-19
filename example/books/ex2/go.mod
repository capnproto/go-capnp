module bookstest2

go 1.17

require (
	books v1.0.0
	capnproto.org/go/capnp/v3 v3.0.0-alpha.1
)

replace books v1.0.0 => ../books
