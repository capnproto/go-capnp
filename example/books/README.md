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
