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
[gocompat]: https://golang.org/doc/go1compat
[godoc]: https://godoc.org/zombiezen.com/go/capnproto2
[license]: https://github.com/zombiezen/go-capnproto/blob/master/LICENSE
