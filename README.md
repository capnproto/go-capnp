[![Build Status](https://travis-ci.org/zombiezen/go-capnproto.svg?branch=master)](https://travis-ci.org/zombiezen/go-capnproto)
[![GoDoc](https://godoc.org/zombiezen.com/go/capnproto?status.svg)](https://godoc.org/zombiezen.com/go/capnproto)

go-capnproto consists of a Go code generator for Cap'n Proto and a Go
package that provides runtime support.  The RPC protocol is not yet
implemented, but there is support for generating interfaces.

News
----

19 February 2015: This is a fork of Jason Aten's [go-capnproto
branch](https://github.com/glycerine/go-capnproto) that supports Cap'n
Proto interfaces.  It's not API compatible, as I chose to clean up some of the
naming rules to bring it more in line with the Protobuf code generator.
Jason has agreed that this should live as its own fork for a while, but I
will try to upstream as much as possible.  As a result, some branches in
this repo may be used to push smaller features back upstream. -Ross

5 April 2014: James McKaskill, the author of go-capnproto (https://github.com/jmckaskill/go-capnproto), 
has been super busy of late, so I agreed to take over as maintainer. This branch 
(https://github.com/glycerine/go-capnproto) includes my recent work to fix bugs in the
creation (originating) of structs for Go, and an implementation of the packing/unpacking capnp specification.
Thanks to Albert Strasheim (https://github.com/alberts/go-capnproto) of CloudFlare for a great set of packing tests. - Jason

Getting started
---------------

You will need the `capnp` tool to compile schemas into Go.  This package has
been tested with Cap'n Proto 0.5.0.

~~~
# first: be sure you have your GOPATH env variable setup.
$ go get -u -t zombiezen.com/go/capnproto
$ cd $GOPATH/src/zombiezen.com/go/capnproto
$ make # will install capnpc-go and compile the test schema aircraftlib/aircraft.capnp, which is used in the tests.
$ diff ./capnpc-go/capnpc-go `which capnpc-go` # you should verify that you are using the capnpc-go binary you just built. There should be no diff. Adjust your PATH if necessary to include the binary capnpc-go that you just built/installed from ./capnpc-go/capnpc-go.
$ go test -v  # confirm all tests are green
~~~

Documentation
-------------

In godoc see http://godoc.org/zombiezen.com/go/capnproto

What is Cap'n Proto?
--------------------

The best cerealization...

https://capnproto.org/

License
-------

MIT - see LICENSE file
