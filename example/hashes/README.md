# RPC Example: Hashses

## Running the example

Navigate to the example directory:

```bash
cd /path/to/go-capnp/example/hashes`
````

Compile the capnp file:

`capnp compile -I ../../std/ -ogo hashes.capnp`

Build your code, and run it:

```
go run cmd/hashesserver.go
```

The output of the example should be the following

```bash
sha1: 0a0a9f2a6772942557ab5355d76af442f8f65e01
```