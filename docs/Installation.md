# Toolchain & Configuration

First, [install the Cap'n Proto tools](https://capnproto.org/install.html).

Then, run the following commands to install the compiler plugin for the Go language:

```bash
$ go install capnproto.org/go/capnp/v3/capnpc-go@latest  # install go compiler plugin
$ GO111MODULE=off go get -u capnproto.org/go/capnp/v3/   # install go-capnproto to $GOPATH
```

Lastly, ensure `$GOPATH/bin` has been added to your shell's `$PATH` variable.

If you get stuck at any point, please [ask us for help](https://matrix.to/#/#go-capnp:matrix.org)!

# Next

Once you have installed `capnp` and the Go plugin, you should [write and
compile your first schema](Writing-Schemas-and-Generating-Code.md).
