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
--------------

~~~
$ make # will install capnpc-go and compile the test schema aircraftlib/aircraft.capnp, which is used in the tests.
$ diff ./capnpc-go/capnpc-go `which capnpc-go` # you should verify that you are using the capnpc-go binary you just built. There should be no diff. Adjust your PATH if necessary to include the binary capnpc-go that you just built/installed from ./capnpc-go/capnpc-go.
$ go test -v  # confirm all tests are green
~~~

