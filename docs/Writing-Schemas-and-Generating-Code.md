# Code Generation

Cap'n Proto works by generating code.  First, you create a schema using [Cap'n Proto's interface-definition language (IDL)](https://capnproto.org/language.html).  Then, you feed this schema into the Cap'n Proto compiler, which generates the corresponding code in your native language.

There is support for [other languages](https://capnproto.org/otherlang.html), too.

## Before you Begin

Ensure that you have installed the `capnp` compiler and the Go language plugin, and that your shell's `$PATH` variable includes `$GOPATH/bin`.

Refer to the [previous section](Getting-Started.md) for instructions.

## Example Schema File

Consider the following schema, stored in `foo/books.capnp`:

```capnp
using Go = import "/go.capnp";
@0x85d3acc39d94e0f8;
$Go.package("books");
$Go.import("foo/books");

struct Book {
    title @0 :Text;
    # Title of the book.

    pageCount @1 :Int32;
    # Number of pages in the book.
}
```

capnpc-go requires that two [annotations](https://capnproto.org/language.html#annotations) be present in your schema:

1. `$Go.package("books")`:  tells the compiler to place `package books` at the top of the generated Go files.
2. `$Go.import("foo/books")`:  declares the full import path within your project.  The compiler uses this to generate the import statement in the auto-generated code, when one of your schemas imports a type from another.

Compilation will fail unless these annotations are present.

## Compiling the Schema

To compile this schema into Go code, run the following command.   Note that the source path `/foo/books.capnp` must correspond to the import path declared in your annotations.

```bash
capnp compile -I /path/to/go-capnp/std -ogo foo/books.capnp
```

> **Tip** 👉 For more compilation options, see `capnp compile --help`.

This will output the `foo/books.capnp.go` file, containing Go structs that can be imported into your programs.  These are ordinary Go types that represent the schema you declared in `books.capnp`.  Each has accessor methods corresponding to the fields declared in the schema.  For example, the `Book` struct will have the methods `Title() (string, error)` and `SetTitle(string) error`.

In the next section, we will show how you can write these structs to a file or transmit them over the network.

# Next

Now that you have generated code from your schema, you should [learn how to work with Cap'n Proto types](Working-with-Capn-Proto-Types.md).
