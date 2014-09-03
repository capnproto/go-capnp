/*
Package capn is a capnproto library for go

see http://kentonv.github.io/capnproto/

capnpc-go provides the compiler backend for capnp
after installing to $PATH capnp files can be compiled with

	capnp compile -ogo *.capnp

capnpc-go requires two annotations for all types. This is the package and
import found in go.capnp. Package is needed to know what package to place at
the head of the generated file and what go name to use when referring to the
type from another package. Import should be the fully qualified import path
and is used to generate import statement from other packages and to detect
when two types are in the same package. Typically these are added as file
annotations. For example:

	using Go = import "github.com/glycerine/go-capnproto/go.capnp";
	$Go.package("main");
	$Go.import("github.com/glycerine/go-capnproto/example");

In capnproto, the unit of communication is a message. A message
consists of one or more of segments to allow easier allocation, but
ideally and typically you just make one segment per message.

Logically, a message organized in a tree of objects, with the root
always being a struct (as opposed to a list or primitive).

Here is an example of writing a new message. We use the demo schema
aircraft.capnp from the aircraftlib directory. You may wish to read
the schema before reading this example.

<< Example moved to its own file: See the file, write_test.go >>

In summary, when you make a new message, you should first make new segment,
and then create the root struct in that segment. Then add your non-child
(contained) objects. This is because, as the spec says:

   The first word of the first segment of the message
   is always a pointer pointing to the message's root
   struct.


All objects are values with pointer semantics that point into the data
in a message or segment. Messages can be read/written from a stream
uncompressed or using the capnproto compression.

In this library a *Segment is taken to refer to both a specific segment as
well as the containing message. This is to reduce the number of types generic
code needs to deal with and allows objects to be created in the same segment
as their outer object (thus reducing the number of far pointers).

Most getters/setters in the library don't return an error. Instead a get that
fails due to an invalid pointer, out of bounds, etc will return the default
value. A invalid set will be noop'ed. If you really need to know whether a set
succeeds then errors are provided by the lower level Object methods.

Since go doesn't have any templating, lists are created for the basic types
and one level of named types. The list of basic types (e.g. List(UInt8),
List(Text), etc) are all provided in this library. Lists of user named types
are created with the user types (e.g. user struct Foo will create a Foo_List
type). capnp schemas that use deeper lists (e.g. List(List(UInt8))) will use
PointerList and the user will have to use the Object.ToList* functions to cast
to the correct type.

For adding documentation comments to the generated code, there's the doc
annotation. This annotation adds the comment to a struct, enum or field so
that godoc will pick it up. For Example:

	struct Zdate $Go.doc("Zdate represents an instance in time") {
	  year  @0   :Int16;
	  month @1   :UInt8;
	  day   @2   :UInt8 ;
	}

Structs

capnpc-go will generate the following for structs:

	// Foo is a value with pointer semantics referencing the data in a
	// segment. Member functions are provided to get/set members in the
	// struct. Getters/setters of an outer struct will use values of type
	// Foo to set/get pointers.
	type Foo capn.Struct

	// NewFoo creates a new orphaned Foo struct. This can then be added to
	// a message by using a Set function which takes a Foo argument.
	func NewFoo(s *capn.Segment) Foo

	// NewRootFoo creates a new root of type Foo in the next unused space in the
	// provided segment. This is distinct from NewFoo as this always
	// creates a root tag. Typically the provided segment should be empty.
	// Remember that a message is a tree of objects with a single root, and
	// you usually have to create the root before any other object in a
	// segment. The only exception would be for a multi-segment message.
	func NewRootFoo(s *capn.Segment) Foo

	// ReadRootFoo reads the root tag at the beginning of the provided
	// segment and returns it as a Foo struct.
	func ReadRootFoo(s *capn.Segment) Foo

	// Foo_List is a value with pointer semantics. It is created for all
	// structs, and is used for List(Foo) in the capnp file.
	type Foo_List capn.List

	// NewFooList creates a new orphaned List(Foo). This can then be added
	// to a message by using a Set function which takes a Foo_List. sz
	// specifies the list size. Due to the list using memory directly in
	// the outgoing buffer (i.e. arena style memory management), the size
	// can not be changed after creation.
	func NewFooList(s *capn.Segment, sz int) Foo_List

	// Len returns the list length. For composite lists this is the number
	// of list elements.
	func (s Foo_List) Len() int

	// At returns a pointer to the i'th element. If i is an invalid index,
	// this will return a null Foo (all getters will return default
	// values, setters will fail). For a composite list the returned value
	// will be a list member. Setting another value to point to list
	// members forces a copy of the data. For pointer lists, the pointer
	// value will be auto-derefenced.
	func (s Foo_List) At(i int) Foo

	// ToArray converts the capnproto list into a go list. For large lists
	// this is inefficient as it has to read all elements. This can be
	// quite convenient especially for iterating as it lets you use a for
	// range clause:
	//	for i, f := range mylist.ToArray() {}
	func (s Foo_List) ToArray() []Foo



Groups

For each group a typedef is created with a different method set for just the
groups fields:

	struct Foo {
		group :Group {
			field @0 :Bool;
		}
	}

	type Foo capn.Struct
	type FooGroup Foo

	func (s Foo) Group() FooGroup
	func (s FooGroup) Field() bool

That way the following may be used to access a field in a group:

	var f Foo
	value := f.Group().Field()

Note that Group accessors just cast the type and so have no overhead

	func (s Foo) Group() FooGroup {return FooGroup(s)}



Unions

Named unions are treated as a group with an inner unnamed union. Unnamed
unions generate an enum Type_Which and a corresponding Which() function:

	struct Foo {
		union {
			a @0 :Bool;
			b @1 :Bool;
		}
	}

	type Foo_Which uint16

	const (
		FOO_A Foo_Which = 0
		FOO_B           = 1
	)

	func (s Foo) A() bool
	func (s Foo) B() bool
	func (s Foo) SetA(v bool)
	func (s Foo) SetB(v bool)
	func (s Foo) Which() Foo_Which

Which() should be checked before using the getters, and the default case must
always be handled.

Setters for single values will set the union discriminator as well as set the
value.

For voids in unions, there is a void setter that just sets the discriminator.
For example:

	struct Foo {
		union {
			a @0 :Void;
			b @1 :Void;
		}
	}

	f.SetA() // Set that we are using A
	f.SetB() // Set that we are using B

For groups in unions, there is a group setter that just sets the
discriminator. This must be called before the group getter can be used to set
values. For example:

	struct Foo {
		union {
			a :group {
				v :Bool
			}
			b :group {
				v :Bool
			}
		}
	}

	f.SetA()         // Set that we are using group A
	f.A().SetV(true) // then we can use the group A getter to set the inner values



Enums

capnpc-go generates enum values in all caps. For example in the capnp file:

	enum ElementSize {
	  empty @0;
	  bit @1;
	  byte @2;
	  twoBytes @3;
	  fourBytes @4;
	  eightBytes @5;
	  pointer @6;
	  inlineComposite @7;
	}

In the generated capnp.go file:

	type ElementSize uint16

	const (
		ELEMENTSIZE_EMPTY           ElementSize = 0
		ELEMENTSIZE_BIT                         = 1
		ELEMENTSIZE_BYTE                        = 2
		ELEMENTSIZE_TWOBYTES                    = 3
		ELEMENTSIZE_FOURBYTES                   = 4
		ELEMENTSIZE_EIGHTBYTES                  = 5
		ELEMENTSIZE_POINTER                     = 6
		ELEMENTSIZE_INLINECOMPOSITE             = 7
	)

In addition an enum.String() function is generated that will convert the constants to a string
for debugging or logging purposes. By default, the enum name is used as the tag value,
but the tags can be customized with a $Go.tag or $Go.notag annotation.

For example:

	enum ElementSize {
		empty @0           $Go.tag("void");
		bit @1             $Go.tag("1 bit");
		byte @2            $Go.tag("8 bits");
		inlineComposite @7 $Go.notag;
	}

In the generated go file:

	func (c ElementSize) String() string {
		switch c {
		case ELEMENTSIZE_EMPTY:
			return "void"
		case ELEMENTSIZE_BIT:
			return "1 bit"
		case ELEMENTSIZE_BYTE:
			return "8 bits"
		default:
			return ""
		}
	}
*/
package capn
