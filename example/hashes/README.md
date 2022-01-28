cd to wherever you installed the source, for example:

`cd $GOPATH/src/go-capnproto2/example/hashes`

Compile the capnp file:

`capnp compile -I$GOPATH/src/capnproto.org/go/capnp/std/ -ogo hashes/hashes.capnp`

Build your code, and run it:

```
go build .
./hashtest
```