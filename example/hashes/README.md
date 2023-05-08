cd to wherever you installed the source:

`cd /path/to/go-capnp/example/hashes`

Compile the capnp file:

`capnp compile -I ../../std/ -ogo hashes.capnp`

Build your code, and run it:

```
go build .
./hashtest
```
