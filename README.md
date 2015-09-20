# Cap'n Proto bindings for Go

[![Build Status](https://travis-ci.org/zombiezen/go-capnproto2.svg?branch=master)](https://travis-ci.org/zombiezen/go-capnproto2)
[![GoDoc](https://godoc.org/zombiezen.com/go/capnproto2?status.svg)][godoc]

go-capnproto consists of:
- a Go code generator for [Cap'n Proto][capnproto]
- a Go package that provides runtime support
- a Go package that implements the RPC protocol

## Getting started

You will need the `capnp` tool to compile schemas into Go.  This package has
been tested with Cap'n Proto 0.5.0.

```
# first: be sure you have your GOPATH env variable setup.
$ go get -u -t zombiezen.com/go/capnproto2/...
$ go test -v zombiezen.com/go/capnproto2/...
```

Then read [the Getting Started guide][gettingstarted].

## News

22 August 2015: API breakage time!  Grep through the commit history for "API
change" to see precise changes.  The main impact on application code is that
more functions that could fail previously and silently swallow errors now return
them.  Most of the methods on Struct or Pointer are now package functions to
improve consistency with generated code.

On the flip side, I am now marking the API as stable, barring major changes to
the Cap'n Proto specification.  See the API stability section below.

16 August 2015: I'm cleaning up the API to make it more Go-like.  This change
will mostly affect those that were using the runtime library directly, but the
generated code will now expose errors in places that it hasn't before.  Watch
the `cleanup` branch for changes, and expect the branch to be merged in the next
few weeks.

*Why the change?* Since most users of the go-capnproto package are depending on
Jason's (@glycerine) fork, I want to take the opportunity to clean up
non-idiomatic parts of the API.  In particular, the current design makes it
difficult to implement new allocation algorithms or make changes to internals
without breaking callers.  The main goals are:

- Surface errors from `Message` that were being silenced before.
- Make all integer parameters use types (e.g. addresses can't be mixed with
  sizes, etc.).
- Make `Pointer` into an interface instead of a struct.  `Pointer` is already
  essentially a generic type, but its fields are not well documented and
  confusing.  By making the generated code embed exactly the pointer type they
  need, this should reduce memory usage and provide more type checking.

6 August 2015: **Level 1 RPC support** with some [known issues][issues].  I've
added a section about compatibility guarantees below. -Ross

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

Consider this package's API as beta software, since the Cap'n Proto spec is not
final.  In the spirit of [the Go 1 compatibility guarantee][gocompat], I will
make every effort to avoid making breaking API changes.  The major cases where I
reserve the right to make breaking changes are:

- Security.
- Changes in the Cap'n Proto specification.
- Bugs.

## Documentation

See the docs on [godoc.org][godoc].

## What is Cap'n Proto?

The best cerealization...

https://capnproto.org/

## License

MIT - see [LICENSE][license] file

[capnproto]: https://capnproto.org/
[gettingstarted]: https://github.com/zombiezen/go-capnproto/wiki/Getting-Started
[glycerine]: https://github.com/glycerine/go-capnproto
[gocompat]: https://golang.org/doc/go1compat
[godoc]: https://godoc.org/zombiezen.com/go/capnproto2
[issue1]: https://github.com/zombiezen/go-capnproto/issues/1
[issues]: https://github.com/zombiezen/go-capnproto/issues
[license]: https://github.com/zombiezen/go-capnproto/blob/master/LICENSE
