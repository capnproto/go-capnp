module hashtest

go 1.17

require capnproto.org/go/capnp/v3 v3.0.0-alpha.4

require hashes v1.0.0 // indirect

replace hashes v1.0.0 => ./hashes
