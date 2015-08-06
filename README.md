# Cap'n Proto bindings for Go

[![Build Status](https://travis-ci.org/zombiezen/go-capnproto.svg?branch=master)](https://travis-ci.org/zombiezen/go-capnproto)
[![GoDoc](https://godoc.org/zombiezen.com/go/capnproto?status.svg)][godoc]

go-capnproto consists of:
- a Go code generator for Cap'n Proto
- a Go package that provides runtime support
- a Go package that implements the RPC protocol

## News

6 August 2015: **Level 1 RPC support** with some [known issues][issues].  I've
added a section about compatibility guarantees below.

23 July 2015: **Level 0 RPC support (and parts of Level 1)!**

The rest of Level 1 will be coming soon. -Ross

19 February 2015: This is a fork of Jason Aten's [go-capnproto branch][glycerine]
that supports Cap'n Proto interfaces.  It's not API compatible, as I chose to
clean up some of the naming rules to bring it more in line with the Protobuf
code generator.  Jason has agreed that this should live as its own fork for a
while, but I will try to upstream as much as possible.  As a result, some
branches in this repo may be used to push smaller features back upstream. -Ross

5 April 2014: James McKaskill, the author of go-capnproto (https://github.com/jmckaskill/go-capnproto), 
has been super busy of late, so I agreed to take over as maintainer. This branch 
(https://github.com/glycerine/go-capnproto) includes my recent work to fix bugs in the
creation (originating) of structs for Go, and an implementation of the packing/unpacking capnp specification.
Thanks to Albert Strasheim (https://github.com/alberts/go-capnproto) of CloudFlare for a great set of packing tests. - Jason

## API Compatibility

Consider this package's API as beta software.  In the spirit of
[the Go 1 compatibility guarantee][gocompat], I will make every effort to avoid
making breaking API changes.  The major cases where I reserve the right to make
breaking changes are:

- Security.
- Changes in the Cap'n Proto specification
- Bugs
- And this code cleanup: #1 (but this will go away soon)


## Getting started

You will need the `capnp` tool to compile schemas into Go.  This package has
been tested with Cap'n Proto 0.5.0.

```
# first: be sure you have your GOPATH env variable setup.
$ go get -u -t zombiezen.com/go/capnproto
$ cd $GOPATH/src/zombiezen.com/go/capnproto
$ make # will install capnpc-go and compile the test schema aircraftlib/aircraft.capnp, which is used in the tests.
$ diff ./capnpc-go/capnpc-go `which capnpc-go` # you should verify that you are using the capnpc-go binary you just built. There should be no diff. Adjust your PATH if necessary to include the binary capnpc-go that you just built/installed from ./capnpc-go/capnpc-go.
$ go test -v  # confirm all tests are green
```

## Documentation

See the docs on [godoc.org][godoc].

## What is Cap'n Proto?

The best cerealization...

https://capnproto.org/

## License

MIT - see [LICENSE][license] file

[gocompat]: https://golang.org/doc/go1compat
[godoc]: https://godoc.org/zombiezen.com/go/capnproto
[issues]: https://github.com/zombiezen/go-capnproto/issues
[license]: https://github.com/zombiezen/go-capnproto/blob/master/LICENSE
[glycerine]: https://github.com/glycerine/go-capnproto
