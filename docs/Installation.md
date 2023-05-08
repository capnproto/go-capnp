# Toolchain & Configuration

First, [install the Cap'n Proto tools](https://capnproto.org/install.html).

Then, run the following command to install the compiler plugin for the
Go language:

```bash
$ go install capnproto.org/go/capnp/v3/capnpc-go@latest  # install go compiler plugin
```

This will install a `capnpc-go` executable under `$(go env GOPATH)/bin`,
which you should make sure has been added to your shell's `$PATH` variable.

You will also need a checkout of of the `go-capnp` repository, for the
included schema:

```
$ git clone https://github.com/capnproto/go-capnp
```

If you get stuck at any point, please [ask us for help](https://matrix.to/#/#go-capnp:matrix.org)!

# Next

Once you have installed `capnp` and the Go plugin, you should [write and
compile your first schema](Writing-Schemas-and-Generating-Code.md).
