# Cap'n Proto bindings for Go

[![GoDoc](https://godoc.org/capnproto.org/go/capnp/v3?status.svg)][godoc]
![License](https://img.shields.io/badge/license-MIT-brightgreen?style=flat-square)
![tests](https://github.com/capnproto/go-capnproto2/workflows/Go/badge.svg)

[Cap’n Proto](https://capnproto.org/) is an insanely fast data interchange format similar to [Protocol Buffers](https://github.com/protocolbuffers/protobuf), but much faster.

It also includes a sophisticated RPC system based on [Object Capabilities](https://en.wikipedia.org/wiki/Object-capability_model), ideal for secure, low-latency applications.

This package provides:
- Go code-generation for Cap'n Proto
- Runtime support for the Go language
- Level 1 support for the [Cap'n Proto RPC](https://capnproto.org/rpc.html) protocol

[godoc]: http://pkg.go.dev/capnproto.org/go/capnp/v3
## Installation

```
$ go get capnproto.org/go/capnp/v3
```

**NOTE:** You will need to install the [`capnp` tool](https://capnproto.org/capnp-tool.html) in order to compile your Cap'n Proto schemas into Go.  This package has been tested with version `0.8.0` of the `capnp` tool.

## Documentation

### Getting Started

Read the ["Getting Started" guide](https://github.com/capnproto/go-capnproto2/wiki/Getting-Started) for a high-level introduction to the package API and workflow.

Browse rest of the [Wiki](https://github.com/capnproto/go-capnproto2/wiki) for in depth explanations of concepts, migration guides, and tutorials.

### API Reference

Available on [GoDoc](http://pkg.go.dev/capnproto.org/go/capnp/v3).

## API Compatibility

Until the official Cap'n Proto spec is finalized, this repository should be considered <u>beta software</u>.

We use [semantic versioning](https://semver.org) to track compatibility and signal breaking changes.  In the spirit of the [Go 1 compatibility guarantee][gocompat], we will make every effort to avoid making breaking API changes within major version numbers, but nevertheless reserve the right to introduce breaking changes for reasons related to:

- Security.
- Changes in the Cap'n Proto specification.
- Bugs.

An exception to this rule is currently in place for the `pogs` package, which is relatively new and may change over time.  However, its functionality has been well-tested, and breaking changes are relatively unlikely.

Note also we may merge breaking changes to the `master` branch without notice.  Users are encouraged to pin their dependencies to a major version, e.g. using the semver-aware features of `go get`.

[gocompat]: https://golang.org/doc/go1compat
## License

MIT - see [LICENSE][] file

[LICENSE]: https://github.com/capnproto/go-capnproto2/blob/master/LICENSE
