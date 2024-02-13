Change directory into the example directory, and compile the example
capnproto schema:

```
cd $GOPATH/src/go-capnproto2/example/books
capnp compile -I ../../std/ -ogo books/books.capnp
```

Then build and run each example:
```
cd ex1 && go build . && ./bookstest1

cd ex2 && go build . && capnp encode ../books/books.capnp Book < ./book.txt | ./bookstest2 && cd ..
```

If this results in:
```
../books/books.capnp.go:37:25: capnp.Struct(s).EncodeAsPtr undefined (type capnp.Struct has no field or method EncodeAsPtr)
../books/books.capnp.go:41:29: capnp.Struct{}.DecodeFromPtr undefined (type capnp.Struct has no field or method DecodeFromPtr)
```

This has been seen with go verion 1.21.6 and 1.22.  To fix you need a version newer than v2.18.0.  v3.0 should fix this, if not released yet, try the newest tagged version newer than v2.18.0.  For example:
```
go install capnproto.org/go/capnp/v3/capnpc-go@v3.0.0-alpha.1
```



