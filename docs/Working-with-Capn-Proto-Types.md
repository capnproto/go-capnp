# Data Model

The Go types you generated in the previous section are actually wrappers around a `[]byte` buffer, with getters and setters that operate by indexing the buffer at specific offsets.  So technically speaking, **there there is no (un)marshal step** in Cap'n Proto.  Converting a generated type to and from `[]byte` is a simple matter of wrapping and unwrapping the buffer:  a constant-time operation on the order of nanoseconds.

In this section, we will learn how to:

1. create and interact with Cap'n Proto types,
2. convert Cap'n Proto types to/from their underlying `[]byte` buffers; and,
3. stream Cap'n Proto types to/from byte streams (e.g. network connections).

## Using Generated Types

Instantiating a type that was generated from your a schema is a three-step process:

1. Instantiate a `capnp.Arena`, which exposes a low-level API to a `[]byte` buffer.
2. Instantiate a new `*capnp.Message`, which allocates capnp structs in the above arena.
3. Instantiate the schema-generated type, which wraps the above `Message` and provides a high-level getter/setter API.

The code below creates a new `books.Book` from the schema you generated in the [previous section](Writing-Schemas-and-Generating-Code.md), and populates its fields.  Steps 1 through 3 are repeated in the comments, for clarity.

```go
package main

import (
    "foo/books"
    "capnproto.org/go/capnp/v3"
)

func main() {
    // Create a new Arena for a books.Book type.  The Arena wraps the underlying
    // buffer, providing a low-level access API.  You probably won't ever need to
    // interact with it directly.  We will ignore the meaning of "single segment"
    // for now.
    arena := capnp.SingleSegment(nil)

    // Make a brand new empty message.  A Message allocates Cap'n Proto structs within
    // its arena.  For convenience, NewMessage also returns the root "segment" of the
    // message, which is needed to instantiate the Book struct.  You don't need to
    // understand segments and roots yet (or maybe ever), but if you're curious, messages
    // and segments are documented here:  https://capnproto.org/encoding.html
    msg, seg, err := capnp.NewMessage(arena)
    if err != nil {
        panic(err)
    }

    // Create a new Book struct.  Every message must have a root struct.  Again, it is
    // not important to understand "root structs" at this point.  For now, just understand
    // that every type you instantiate needs to be a "root", unless you plan on assigning
    // it to another object.  When in doubt, use NewRootXXX.
    //
    // If you're insatiably curious, see:  https://capnproto.org/encoding.html#messages
    book, err := books.NewRootBook(seg)
    if err != nil {
        panic(err)
    }

    // Great, we have our book!  Now let's set some fields.  Each field you declared in
    // your schema will produce two methods on the generated type.  The "getter" method
    // has the name of the field, for example:  Book.Title().  The corresponding "setter"
    // method is prefixed with "Set", for example:  Book.SetTitle().
    //
    // Some getters and setters return errors, which we are ignoring in this example for
    // the sake of clarity.  Your code SHOULD check these errors and handle them.
    //
    // To begin, we set the book's title to "War and Peace".
    _ = book.SetTitle("War and Peace")

    // Then, we set the page count.
    book.SetPageCount(1440)

    // Finally, we "get" these fields and print them.
    title, _ := book.Title()
    fmt.Printf("%s (%d pages)", title, book.Pages())
}
```

So far, this looks a lot like Protocol Buffers.  In the next few sections, we'll show you where Cap'n Proto really comes into its own:  data serialization.  This will also show you where the `*capnp.Message` type is used.

## Marshalling and Unmarshalling

In a narrow technical sense, there is no "marshalling" or "unmarshalling" in Cap'n Proto.  This is because the get and set operations you saw above act **directly on the object's underlying `[]byte` buffer**.  This is what makes reading and writing Cap'n Proto data blisteringly fast!

Despite this slight inaccuracy, we still refer to the process of converting objects to and from `[]byte`s as *marshalling* and *unmarshalling*, and provide the corresponding methods [`Message.Marshal`](https://pkg.go.dev/capnproto.org/go/capnp/v3#Message.Marshal) and [`Message.Unmarshal`](https://pkg.go.dev/capnproto.org/go/capnp/v3#Message.Unmarshal).   We use this terminology for two reasons:

1. the terms are familiar to most Go developers, and
2. they correctly convey the essence of the operations:  `type -> []byte` and `[]byte -> type`, respectively for marshal and unmarshal.

### Marshalling

To marshal the `books.Book` instance from the previous example, we need only call the `Marshal` method on the corresponding `*capnp.Message`:

```go
b, err := msg.Marshal()
if err != nil {
    panic(err)
}

// send b over the network, or write it to a file, or whatever...
```

### Unmarshalling

In the above example, `Marshal` returns the book's underlying buffer unmodified.  To unmarshal it into a new `books.Book` object, you need to:

1. call `capnp.Unmarshal` to obtain a new `*capnp.Message`; then,
2. call `books.ReadRootBook`, passing in the newly-obtained message.

Here is the code to do so:

```go
msg, err := capnp.Unmarshal(b)
if err != nil {
    panic(err)
}

// Again, don't worry about the meaning of "root" for now.
// When in doubt, use the "root" version of functions.
book, err := books.ReadRootBook(msg)
if err != nil {
    panic(err)
}

title, _ := book.Title()
fmt.Printf("%s (%d pages)", title, book.Pages())
```

### Using the Packed Encoding

Cap'n Proto supports a [packed encoding](https://capnproto.org/encoding.html#packing), that provides ultra-fast compression.  To use the packed encoding, substitute `Message.Marshal` with `Message.MarshalPacked` and `capnp.Unmarshal` with `capnp.UnmarshalPacked`.

## Streaming

The `Marshal` and `MarshalPacked` methods are suitable for such things as writing capnp types to a file, or sending a single object over the network, e.g. in an HTTP request.  But if you want to stream multiple objects over, say, a network connection, you'll need a way of "framing" the stream, _i.e._ of separating the byte stream into different objects.  For this, we use the `Encoder` and `Decoder` types.

The `Encoder` and `Decoder` types allow you to stream generated types to and from any `io.Writer` or `io.Reader`, respectively.

### Writing to a Byte Stream

```go
// Create a new encoder that streams messages to stdout.
// You can also use NewPackedEncoder if you want to compress
// the data.
encoder := capnp.NewEncoder(os.Stdout)

// Send the book's underlying *capnp.Message.  Note that we
// could have also passed the 'msg' variable we obtained from
// our previous call to capnp.Unmarshal or capnp.NewMessage.
// In most cases, however, it is more convenient to use the
// generated type's Message() method.
err = encoder.Encode(book.Message())
if err != nil {
    panic(err)
}
```

### Reading from a Byte Stream

```go
// Create a new decoder that reads from stdin.
// Use capnp.NewPackedDecoder if you are expecting a
// packed byte-stream.
decoder := capnp.NewDecoder(os.Stdin)

// Read the message from stdin.
msg, err := decoder.Decode()
if err != nil {
    panic(err)
}

// Extract the root struct from the message.
book, err := books.ReadRootBook(msg)
if err != nil {
    panic(err)
}

// Access fields from the struct.  Again, we're
// ignoring errors, but you definitely shouldn't.
title, _ := book.Title()
pageCount := book.PageCount()
fmt.Printf("%q has %d pages\n", title, pageCount)
```

# Next

Now that you understand how marshalling works, you're ready to [write your first RPC service](Remote-Procedure-Calls-using-Interfaces.md).
