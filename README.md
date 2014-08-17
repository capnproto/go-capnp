License
-------

MIT - see LICENSE file

Documentation
-------------
In godoc see http://godoc.org/github.com/glycerine/go-capnproto


News
----

5 April 2014: James McKaskill, the author of go-capnproto (https://github.com/jmckaskill/go-capnproto), 
has been super busy of late, so I agreed to take over as maintainer. This branch 
(https://github.com/glycerine/go-capnproto) includes my recent work to fix bugs in the
creation (originating) of structs for Go, and an implementation of the packing/unpacking capnp specification.
Thanks to Albert Strasheim (https://github.com/alberts/go-capnproto) of CloudFlare for a great set of packing tests. - Jason

Getting started
---------------

~~~
# first: be sure you have your GOPATH env variable setup.
$ go get -u -t github.com/glycerine/go-capnproto
$ cd $GOPATH/src/github.com/glycerine/go-capnproto
$ make # will install capnpc-go and compile the test schema aircraftlib/aircraft.capnp, which is used in the tests.
$ diff ./capnpc-go/capnpc-go `which capnpc-go` # you should verify that you are using the capnpc-go binary you just built. There should be no diff. Adjust your PATH if necessary to include the binary capnpc-go that you just built/installed from ./capnpc-go/capnpc-go.
$ go test -v  # confirm all tests are green
~~~

What is Cap'n Proto?
--------------------

The best cerealization...

http://kentonv.github.io/capnproto/

Note, go-capnproto doesn't support the RPC layer of capnp, which is a more recent work-in-progress than the serialization.  Personally I use capnp (for schema based serialization) with nanomsg (for network transport). Here is a toy example of using them together: https://github.com/glycerine/gozbus

