/*
Package capnp is a Cap'n Proto library for Go.

see https://capnproto.org/

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

	using Go = import "zombiezen.com/go/capnproto/go.capnp";
	$Go.package("main");
	$Go.import("zombiezen.com/go/capnproto/example");

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
	type Foo capnp.Struct

	// NewFoo creates a new orphaned Foo struct. This can then be added to
	// a message by using a Set function which takes a Foo argument.
	func NewFoo(s *capnp.Segment) Foo

	// NewRootFoo creates a new root of type Foo in the next unused space in the
	// provided segment. This is distinct from NewFoo as this always
	// creates a root tag. Typically the provided segment should be empty.
	// Remember that a message is a tree of objects with a single root, and
	// you usually have to create the root before any other object in a
	// segment. The only exception would be for a multi-segment message.
	func NewRootFoo(s *capnp.Segment) Foo

	// ReadRootFoo reads the root tag at the beginning of the provided
	// segment and returns it as a Foo struct.
	func ReadRootFoo(s *capnp.Segment) Foo

	// Segment returns the struct's segment.
	func (s Foo) Segment() *capnp.Segment

	// Foo_List is a value with pointer semantics. It is created for all
	// structs, and is used for List(Foo) in the capnp file.
	type Foo_List capnp.List

	// NewFoo_List creates a new orphaned List(Foo). This can then be added
	// to a message by using a Set function which takes a Foo_List. sz
	// specifies the list size. Due to the list using memory directly in
	// the outgoing buffer (i.e. arena style memory management), the size
	// can not be changed after creation.
	func NewFoo_List(s *capnp.Segment, sz int) Foo_List

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

	// Foo_Promise is a promise for a Foo.  Methods are provided to get
	// promises of struct and interface fields.
	type Foo_Promise capnp.Pipeline

	// Get waits until the promise is resolved and returns the result.
	func (p Foo_Promise) Get() (Foo, error)


Groups

For each group a typedef is created with a different method set for just the
groups fields:

	struct Foo {
		group :Group {
			field @0 :Bool;
		}
	}

	type Foo capnp.Struct
	type Foo_group Foo

	func (s Foo) Group() Foo_group
	func (s Foo_group) Field() bool

That way the following may be used to access a field in a group:

	var f Foo
	value := f.Group().Field()

Note that Group accessors just cast the type and so have no overhead

	func (s Foo) Group() Foo_group {return Foo_group(s)}



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
		Foo_Which_a Foo_Which = 0
		Foo_Which_b Foo_Which = 1
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

capnpc-go generates enum values as constants. For example in the capnp file:

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
		ElementSize_empty           ElementSize = 0
		ElementSize_bit             ElementSize = 1
		ElementSize_byte            ElementSize = 2
		ElementSize_twoBytes        ElementSize = 3
		ElementSize_fourBytes       ElementSize = 4
		ElementSize_eightBytes      ElementSize = 5
		ElementSize_pointer         ElementSize = 6
		ElementSize_inlineComposite ElementSize = 7
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
		case ElementSize_empty:
			return "void"
		case ElementSize_bit:
			return "1 bit"
		case ElementSize_byte:
			return "8 bits"
		default:
			return ""
		}
	}


Interfaces

capnpc-go generates type-safe Client wrappers for interfaces. For parameter
lists and result lists, structs are generated as described above with the names
Interface_method_Params and Interface_method_Results, unless a single struct
type is used. For example, for this interface:

	interface Calculator {
		evaluate @0 (expression :Expression) -> (value :Value);
	}

capnpc-go generates the following Go code (along with the structs
Calculator_evaluate_Params and Calculator_evaluate_Results):

	// Calculator is a client to a Calculator interface.
	type Calculator struct { c capnp.Client }

	// NewCalculator creates a Calculator from a generic promise.
	func NewCalculator(c capnp.Client) Calculator

	// GenericClient returns the underlying generic client.
	func (c Calculator) GenericClient() capnp.Client

	// IsNull returns whether the underlying client is nil.
	func (c Calculator) IsNull() bool

	// Evaluate calls `evaluate` on the client.  params is called on a newly
	// allocated Calculator_evaluate_Params to fill in the parameters.
	func (c Calculator) Evaluate(
		ctx context.Context,
		params func(Calculator_evaluate_Params),
		opts ...capnp.CallOption) *Calculator_evaluate_Results_Promise

capnpc-go also generates code to implement the interface:

	// A Calculator_Server implements the Calculator interface.
	type Calculator_Server interface {
		Evaluate(Calculator_evaluate_Call) error
	}

	// Calculator_evaluate_Call holds the arguments for a Calculator.evaluate server call.
	type Calculator_evaluate_Call struct {
		Ctx     context.Context
		Options capnp.CallOptions
		Params  Calculator_evaluate_Params
		Results Calculator_evaluate_Results
	}

	// Calculator_ServerToClient is equivalent to calling:
	// NewCalculator(capnp.NewServer(Calculator_Methods(nil, s), s))
	// If s does not implement the Close method, then nil is used.
	func Calculator_ServerToClient(s Calculator_Server) Calculator

	// Calculator_Methods appends methods from Calculator that call to server and
	// returns the methods.  If methods is nil or the capacity of the underlying
	// slice is too small, a new slice is returned.
	func Calculator_Methods(methods []server.Method, s Calculator_Server) []server.Method

Since a single capability may want to implement many interfaces, you can
use multiple *_Methods functions to build a single slice to send to
NewServer.

An example of combining the client/server code to communicate with a locally
implemented Calculator:

	var srv Calculator_Server
	calc := Calculator_ServerToClient(srv)
	result := calc.Evaluate(ctx, func(params Calculator_evaluate_Params) {
		params.SetExpression(expr)
	})
	val := result.Value().Get()

A note about message ordering: when implementing a server method, you
are responsible for acknowledging delivery of a method call.  Failure to
do so can cause deadlocks.  See the server.Ack function for more details.
*/
package capnp // import "zombiezen.com/go/capnproto"
