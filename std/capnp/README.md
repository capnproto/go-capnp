This directory contains go packages for all of the capnproto
schemas that ship with capnproto itself. They are generated with
the help of `./gen.sh`. Though some manual modifications have been made
to the schema to correct name collisions (See below).

# `./gen.sh` Usage

Executing:

    ./gen.sh import /path/to/schemas

Will copy all the `*.capnp` files from the given directory into this
directory, and add annotations to them specifying package names and
import paths for `capnpc-go`.

    ./gen.sh compile

Will generate go packages for each of the schemas in the current
directory.

    ./gen.sh clean-go

Will remove all generated go source files from this directory. Finally,

    ./gen.sh clean-all

Will remove both the go source files and the imported schemas.

`gen.sh` does some name mangling to ensure that the generated packages
are actually legal go. However, this is not meant to be a
general-purpose solution; it is only intended to work for the base
schemas.

# Footnote: annotations get an underscore in certain cases

Under certain circumstances, `capnpc-go` will rename the identifier when it
generates go code. As an example, if two declarations `struct Foo` and
`annotation foo` exist in the schema, `capnpc-go` has to capitalize
`annotation foo` in the generated code. But that would not compile since
there is already a `Foo` from `struct Foo`.

To address this kind of issue, `capnpc-go` will change the annotation to
`Foo_` (by adding the trailing "_"). This can be overridden on a
case-by-case basis if you need some other name. Just add `$Go.name`:

```
annotation foo(struct, field) :Void $Go.name("fooAnnotation");
```

# Directory Structure

The directory structure of this repository is designed such that when
compiling other schema, it should be sufficient to execute:

    capnp compile -I ${path_to_this_repository}/std -ogo ${schama_name}.capnp

And have the `$import` statements in existing schema "just work."

To achieve this, the base schemas themselves are stored as
`/std/capnp/${schema_name}.capnp`. The generated go source files are
stored in  a subdirectory, to make them their own package:
`/std/capnp/${mangled_schema_name}/${mangled_schema_name}.capnp.go`.

In addition to the upstream base schemas, we also ship a schema
`/std/go.capnp`, which contains annotations used by `go-capnpc`. Its
usage is described in the top-level README. The generated source is
placed in the root of the repository, making it part of the go package
`capnproto.org/go/capnp/v3`.
