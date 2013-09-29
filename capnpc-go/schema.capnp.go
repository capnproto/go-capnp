package main

import (
	C "github.com/jmckaskill/go-capnproto"
	"math"
)

var package_annotation = uint64(0xbea97f1023792be0)
var import_annotation = uint64(0xe130b601260e44b5)

type Node C.Struct
type NodeStruct Node
type NodeEnum Node
type NodeInterface Node
type NodeConst Node
type NodeAnnotation Node
type Node_which uint16

const (
	NODE_FILE       Node_which = 0
	NODE_STRUCT                = 1
	NODE_ENUM                  = 2
	NODE_INTERFACE             = 3
	NODE_CONST                 = 4
	NODE_ANNOTATION            = 5
)

func (s Node) which() Node_which                        { return Node_which(C.Struct(s).Get16(12)) }
func (s Node) Id() uint64                               { return C.Struct(s).Get64(0) }
func (s Node) DisplayName() string                      { return C.Struct(s).GetObject(0).ToString() }
func (s Node) DisplayNamePrefixLength() uint32          { return C.Struct(s).Get32(8) }
func (s Node) NestedNodes() NodeNestedNode_List         { return NodeNestedNode_List(C.Struct(s).GetObject(1)) }
func (s Node) Annotations() Annotation_List             { return Annotation_List(C.Struct(s).GetObject(2)) }
func (s Node) Struct() NodeStruct                       { return NodeStruct(s) }
func (s NodeStruct) DataWordCount() uint16              { return C.Struct(s).Get16(14) }
func (s NodeStruct) PointerCount() uint16               { return C.Struct(s).Get16(24) }
func (s NodeStruct) PreferredListEncoding() ElementSize { return ElementSize(C.Struct(s).Get16(26)) }
func (s NodeStruct) IsGroup() bool                      { return C.Struct(s).Get1(224) }
func (s NodeStruct) DiscriminantCount() uint16          { return C.Struct(s).Get16(30) }
func (s NodeStruct) DiscriminantOffset() uint32         { return C.Struct(s).Get32(32) }
func (s NodeStruct) Fields() Field_List                 { return Field_List(C.Struct(s).GetObject(3)) }
func (s Node) Enum() NodeEnum                           { return NodeEnum(s) }
func (s NodeEnum) Enumerants() Enumerant_List           { return Enumerant_List(C.Struct(s).GetObject(3)) }
func (s Node) Interface() NodeInterface                 { return NodeInterface(s) }
func (s Node) Const() NodeConst                         { return NodeConst(s) }
func (s NodeConst) Type() Type                          { return Type(C.Struct(s).GetObject(3).ToStruct()) }
func (s NodeConst) Value() Value                        { return Value(C.Struct(s).GetObject(4).ToStruct()) }
func (s Node) Annotation() NodeAnnotation               { return NodeAnnotation(s) }

type Node_List C.PointerList

func (s Node_List) Len() int      { return C.PointerList(s).Len() }
func (s Node_List) At(i int) Node { return Node(C.PointerList(s).At(i).ToStruct()) }
func (s Node_List) ToArray() []Node {
	v := make([]Node, s.Len())
	for i := 0; i < s.Len(); i++ {
		v[i] = s.At(i)
	}
	return v
}

type NodeNestedNode C.Struct

func (s NodeNestedNode) Name() string { return C.Struct(s).GetObject(0).ToString() }
func (s NodeNestedNode) Id() uint64   { return C.Struct(s).Get64(0) }

type NodeNestedNode_List C.PointerList

func (s NodeNestedNode_List) Len() int { return C.PointerList(s).Len() }
func (s NodeNestedNode_List) At(i int) NodeNestedNode {
	return NodeNestedNode(C.PointerList(s).At(i).ToStruct())
}
func (s NodeNestedNode_List) ToArray() []NodeNestedNode {
	v := make([]NodeNestedNode, s.Len())
	for i := 0; i < s.Len(); i++ {
		v[i] = s.At(i)
	}
	return v
}

type Field C.Struct
type FieldSlot Field
type FieldGroup Field
type FieldOrdinal Field
type Field_which uint16

const (
	FIELD_SLOT  Field_which = 0
	FIELD_GROUP             = 1
)

type FieldOrdinal_which uint16

const (
	FIELDORDINAL_IMPLICIT FieldOrdinal_which = 0
	FIELDORDINAL_EXPLICIT                    = 1
)

func (s Field) which() Field_which               { return Field_which(C.Struct(s).Get16(8)) }
func (s Field) Name() string                     { return C.Struct(s).GetObject(0).ToString() }
func (s Field) CodeOrder() uint16                { return C.Struct(s).Get16(0) }
func (s Field) DiscriminantValue() uint16        { return C.Struct(s).Get16(2) ^ 65535 }
func (s Field) Slot() FieldSlot                  { return FieldSlot(s) }
func (s FieldSlot) Offset() uint32               { return C.Struct(s).Get32(4) }
func (s FieldSlot) Type() Type                   { return Type(C.Struct(s).GetObject(2).ToStruct()) }
func (s FieldSlot) DefaultValue() Value          { return Value(C.Struct(s).GetObject(3).ToStruct()) }
func (s Field) Group() FieldGroup                { return FieldGroup(s) }
func (s FieldGroup) TypeId() uint64              { return C.Struct(s).Get64(16) }
func (s Field) Ordinal() FieldOrdinal            { return FieldOrdinal(s) }
func (s FieldOrdinal) which() FieldOrdinal_which { return FieldOrdinal_which(C.Struct(s).Get16(10)) }

type Field_List C.PointerList

func (s Field_List) Len() int       { return C.PointerList(s).Len() }
func (s Field_List) At(i int) Field { return Field(C.PointerList(s).At(i).ToStruct()) }
func (s Field_List) ToArray() []Field {
	v := make([]Field, s.Len())
	for i := 0; i < s.Len(); i++ {
		v[i] = s.At(i)
	}
	return v
}

type Enumerant C.Struct

func (s Enumerant) Name() string      { return C.Struct(s).GetObject(0).ToString() }
func (s Enumerant) CodeOrder() uint16 { return C.Struct(s).Get16(0) }

type Enumerant_List C.PointerList

func (s Enumerant_List) Len() int           { return C.PointerList(s).Len() }
func (s Enumerant_List) At(i int) Enumerant { return Enumerant(C.PointerList(s).At(i).ToStruct()) }
func (s Enumerant_List) ToArray() []Enumerant {
	v := make([]Enumerant, s.Len())
	for i := 0; i < s.Len(); i++ {
		v[i] = s.At(i)
	}
	return v
}

type Type C.Struct
type TypeList Type
type TypeEnum Type
type TypeStruct Type
type TypeInterface Type
type Type_which uint16

const (
	TYPE_VOID      Type_which = 0
	TYPE_BOOL                 = 1
	TYPE_INT8                 = 2
	TYPE_INT16                = 3
	TYPE_INT32                = 4
	TYPE_INT64                = 5
	TYPE_UINT8                = 6
	TYPE_UINT16               = 7
	TYPE_UINT32               = 8
	TYPE_UINT64               = 9
	TYPE_FLOAT32              = 10
	TYPE_FLOAT64              = 11
	TYPE_TEXT                 = 12
	TYPE_DATA                 = 13
	TYPE_LIST                 = 14
	TYPE_ENUM                 = 15
	TYPE_STRUCT               = 16
	TYPE_INTERFACE            = 17
	TYPE_OBJECT               = 18
)

func (s Type) which() Type_which        { return Type_which(C.Struct(s).Get16(0)) }
func (s Type) List() TypeList           { return TypeList(s) }
func (s TypeList) ElementType() Type    { return Type(C.Struct(s).GetObject(0).ToStruct()) }
func (s Type) Enum() TypeEnum           { return TypeEnum(s) }
func (s TypeEnum) TypeId() uint64       { return C.Struct(s).Get64(8) }
func (s Type) Struct() TypeStruct       { return TypeStruct(s) }
func (s TypeStruct) TypeId() uint64     { return C.Struct(s).Get64(8) }
func (s Type) Interface() TypeInterface { return TypeInterface(s) }

type Value C.Struct
type Value_which uint16

const (
	VALUE_VOID      Value_which = 0
	VALUE_BOOL                  = 1
	VALUE_INT8                  = 2
	VALUE_INT16                 = 3
	VALUE_INT32                 = 4
	VALUE_INT64                 = 5
	VALUE_UINT8                 = 6
	VALUE_UINT16                = 7
	VALUE_UINT32                = 8
	VALUE_UINT64                = 9
	VALUE_FLOAT32               = 10
	VALUE_FLOAT64               = 11
	VALUE_TEXT                  = 12
	VALUE_DATA                  = 13
	VALUE_LIST                  = 14
	VALUE_ENUM                  = 15
	VALUE_STRUCT                = 16
	VALUE_INTERFACE             = 17
	VALUE_OBJECT                = 18
)

func (s Value) which() Value_which { return Value_which(C.Struct(s).Get16(0)) }
func (s Value) Bool() bool         { return C.Struct(s).Get1(16) }
func (s Value) Int8() int8         { return int8(C.Struct(s).Get8(2)) }
func (s Value) Int16() int16       { return int16(C.Struct(s).Get16(2)) }
func (s Value) Int32() int32       { return int32(C.Struct(s).Get32(4)) }
func (s Value) Int64() int64       { return int64(C.Struct(s).Get64(8)) }
func (s Value) Uint8() uint8       { return C.Struct(s).Get8(2) }
func (s Value) Uint16() uint16     { return C.Struct(s).Get16(2) }
func (s Value) Uint32() uint32     { return C.Struct(s).Get32(4) }
func (s Value) Uint64() uint64     { return C.Struct(s).Get64(8) }
func (s Value) Float32() float32   { return math.Float32frombits(C.Struct(s).Get32(4)) }
func (s Value) Float64() float64   { return math.Float64frombits(C.Struct(s).Get64(8)) }
func (s Value) Enum() uint16       { return C.Struct(s).Get16(2) }
func (s Value) Text() string       { return C.Struct(s).GetObject(0).ToString() }
func (s Value) Data() []byte       { return C.Struct(s).GetObject(0).ToData() }
func (s Value) List() C.Object     { return C.Struct(s).GetObject(0) }
func (s Value) Struct() C.Object   { return C.Struct(s).GetObject(0) }
func (s Value) Object() C.Object   { return C.Struct(s).GetObject(0) }

type Annotation C.Struct

func (s Annotation) Id() uint64   { return C.Struct(s).Get64(0) }
func (s Annotation) Value() Value { return Value(C.Struct(s).GetObject(0).ToStruct()) }

type Annotation_List C.PointerList

func (s Annotation_List) Len() int            { return C.PointerList(s).Len() }
func (s Annotation_List) At(i int) Annotation { return Annotation(C.PointerList(s).At(i).ToStruct()) }
func (s Annotation_List) ToArray() []Annotation {
	v := make([]Annotation, s.Len())
	for i := 0; i < s.Len(); i++ {
		v[i] = s.At(i)
	}
	return v
}

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

type CodeGeneratorRequest C.Struct

func (s CodeGeneratorRequest) Nodes() Node_List { return Node_List(C.Struct(s).GetObject(0)) }
func (s CodeGeneratorRequest) RequestedFiles() CodeGeneratorRequestRequestedFile_List {
	return CodeGeneratorRequestRequestedFile_List(C.Struct(s).GetObject(1))
}

type CodeGeneratorRequestRequestedFile C.Struct

func (s CodeGeneratorRequestRequestedFile) Id() uint64 { return C.Struct(s).Get64(0) }
func (s CodeGeneratorRequestRequestedFile) Filename() string {
	return C.Struct(s).GetObject(0).ToString()
}
func (s CodeGeneratorRequestRequestedFile) Imports() CodeGeneratorRequestRequestedFileImport_List {
	return CodeGeneratorRequestRequestedFileImport_List(C.Struct(s).GetObject(1))
}

type CodeGeneratorRequestRequestedFile_List C.PointerList

func (s CodeGeneratorRequestRequestedFile_List) Len() int { return C.PointerList(s).Len() }
func (s CodeGeneratorRequestRequestedFile_List) At(i int) CodeGeneratorRequestRequestedFile {
	return CodeGeneratorRequestRequestedFile(C.PointerList(s).At(i).ToStruct())
}
func (s CodeGeneratorRequestRequestedFile_List) ToArray() []CodeGeneratorRequestRequestedFile {
	v := make([]CodeGeneratorRequestRequestedFile, s.Len())
	for i := 0; i < s.Len(); i++ {
		v[i] = s.At(i)
	}
	return v
}

type CodeGeneratorRequestRequestedFileImport C.Struct

func (s CodeGeneratorRequestRequestedFileImport) Id() uint64 { return C.Struct(s).Get64(0) }
func (s CodeGeneratorRequestRequestedFileImport) Name() string {
	return C.Struct(s).GetObject(0).ToString()
}

type CodeGeneratorRequestRequestedFileImport_List C.PointerList

func (s CodeGeneratorRequestRequestedFileImport_List) Len() int { return C.PointerList(s).Len() }
func (s CodeGeneratorRequestRequestedFileImport_List) At(i int) CodeGeneratorRequestRequestedFileImport {
	return CodeGeneratorRequestRequestedFileImport(C.PointerList(s).At(i).ToStruct())
}
