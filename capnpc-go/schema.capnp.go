package main

// AUTO GENERATED - DO NOT EDIT

import (
	C "github.com/glycerine/go-capnproto"
	"math"
)

const (
	FieldNoDiscriminant = uint16(65535)
)

type Node C.Struct
type NodeStruct Node
type NodeEnum Node
type NodeInterface Node
type NodeConst Node
type NodeAnnotation Node
type Node_Which uint16

const (
	NODE_FILE       Node_Which = 0
	NODE_STRUCT     Node_Which = 1
	NODE_ENUM       Node_Which = 2
	NODE_INTERFACE  Node_Which = 3
	NODE_CONST      Node_Which = 4
	NODE_ANNOTATION Node_Which = 5
)

func NewNode(s *C.Segment) Node { return Node(s.NewStruct(C.ObjectSize{DataSize: 40, PointerCount: 6})) }
func NewRootNode(s *C.Segment) Node {
	return Node(s.NewRootStruct(C.ObjectSize{DataSize: 40, PointerCount: 6}))
}
func AutoNewNode(s *C.Segment) Node {
	return Node(s.NewStructAR(C.ObjectSize{DataSize: 40, PointerCount: 6}))
}
func ReadRootNode(s *C.Segment) Node                        { return Node(s.Root(0).ToStruct()) }
func (s Node) Which() Node_Which                            { return Node_Which(C.Struct(s).Get16(12)) }
func (s Node) Id() uint64                                   { return C.Struct(s).Get64(0) }
func (s Node) SetId(v uint64)                               { C.Struct(s).Set64(0, v) }
func (s Node) DisplayName() string                          { return C.Struct(s).GetObject(0).ToText() }
func (s Node) SetDisplayName(v string)                      { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s Node) DisplayNamePrefixLength() uint32              { return C.Struct(s).Get32(8) }
func (s Node) SetDisplayNamePrefixLength(v uint32)          { C.Struct(s).Set32(8, v) }
func (s Node) ScopeId() uint64                              { return C.Struct(s).Get64(16) }
func (s Node) SetScopeId(v uint64)                          { C.Struct(s).Set64(16, v) }
func (s Node) Parameters() NodeParameter_List               { return NodeParameter_List(C.Struct(s).GetObject(5)) }
func (s Node) SetParameters(v NodeParameter_List)           { C.Struct(s).SetObject(5, C.Object(v)) }
func (s Node) IsGeneric() bool                              { return C.Struct(s).Get1(288) }
func (s Node) SetIsGeneric(v bool)                          { C.Struct(s).Set1(288, v) }
func (s Node) NestedNodes() NodeNestedNode_List             { return NodeNestedNode_List(C.Struct(s).GetObject(1)) }
func (s Node) SetNestedNodes(v NodeNestedNode_List)         { C.Struct(s).SetObject(1, C.Object(v)) }
func (s Node) Annotations() Annotation_List                 { return Annotation_List(C.Struct(s).GetObject(2)) }
func (s Node) SetAnnotations(v Annotation_List)             { C.Struct(s).SetObject(2, C.Object(v)) }
func (s Node) SetFile()                                     { C.Struct(s).Set16(12, 0) }
func (s Node) Struct() NodeStruct                           { return NodeStruct(s) }
func (s Node) SetStruct()                                   { C.Struct(s).Set16(12, 1) }
func (s NodeStruct) DataWordCount() uint16                  { return C.Struct(s).Get16(14) }
func (s NodeStruct) SetDataWordCount(v uint16)              { C.Struct(s).Set16(14, v) }
func (s NodeStruct) PointerCount() uint16                   { return C.Struct(s).Get16(24) }
func (s NodeStruct) SetPointerCount(v uint16)               { C.Struct(s).Set16(24, v) }
func (s NodeStruct) PreferredListEncoding() ElementSize     { return ElementSize(C.Struct(s).Get16(26)) }
func (s NodeStruct) SetPreferredListEncoding(v ElementSize) { C.Struct(s).Set16(26, uint16(v)) }
func (s NodeStruct) IsGroup() bool                          { return C.Struct(s).Get1(224) }
func (s NodeStruct) SetIsGroup(v bool)                      { C.Struct(s).Set1(224, v) }
func (s NodeStruct) DiscriminantCount() uint16              { return C.Struct(s).Get16(30) }
func (s NodeStruct) SetDiscriminantCount(v uint16)          { C.Struct(s).Set16(30, v) }
func (s NodeStruct) DiscriminantOffset() uint32             { return C.Struct(s).Get32(32) }
func (s NodeStruct) SetDiscriminantOffset(v uint32)         { C.Struct(s).Set32(32, v) }
func (s NodeStruct) Fields() Field_List                     { return Field_List(C.Struct(s).GetObject(3)) }
func (s NodeStruct) SetFields(v Field_List)                 { C.Struct(s).SetObject(3, C.Object(v)) }
func (s Node) Enum() NodeEnum                               { return NodeEnum(s) }
func (s Node) SetEnum()                                     { C.Struct(s).Set16(12, 2) }
func (s NodeEnum) Enumerants() Enumerant_List               { return Enumerant_List(C.Struct(s).GetObject(3)) }
func (s NodeEnum) SetEnumerants(v Enumerant_List)           { C.Struct(s).SetObject(3, C.Object(v)) }
func (s Node) Interface() NodeInterface                     { return NodeInterface(s) }
func (s Node) SetInterface()                                { C.Struct(s).Set16(12, 3) }
func (s NodeInterface) Methods() Method_List                { return Method_List(C.Struct(s).GetObject(3)) }
func (s NodeInterface) SetMethods(v Method_List)            { C.Struct(s).SetObject(3, C.Object(v)) }
func (s NodeInterface) Superclasses() Superclass_List {
	return Superclass_List(C.Struct(s).GetObject(4))
}
func (s NodeInterface) SetSuperclasses(v Superclass_List) { C.Struct(s).SetObject(4, C.Object(v)) }
func (s Node) Const() NodeConst                           { return NodeConst(s) }
func (s Node) SetConst()                                  { C.Struct(s).Set16(12, 4) }
func (s NodeConst) Type() Type                            { return Type(C.Struct(s).GetObject(3).ToStruct()) }
func (s NodeConst) SetType(v Type)                        { C.Struct(s).SetObject(3, C.Object(v)) }
func (s NodeConst) Value() Value                          { return Value(C.Struct(s).GetObject(4).ToStruct()) }
func (s NodeConst) SetValue(v Value)                      { C.Struct(s).SetObject(4, C.Object(v)) }
func (s Node) Annotation() NodeAnnotation                 { return NodeAnnotation(s) }
func (s Node) SetAnnotation()                             { C.Struct(s).Set16(12, 5) }
func (s NodeAnnotation) Type() Type                       { return Type(C.Struct(s).GetObject(3).ToStruct()) }
func (s NodeAnnotation) SetType(v Type)                   { C.Struct(s).SetObject(3, C.Object(v)) }
func (s NodeAnnotation) TargetsFile() bool                { return C.Struct(s).Get1(112) }
func (s NodeAnnotation) SetTargetsFile(v bool)            { C.Struct(s).Set1(112, v) }
func (s NodeAnnotation) TargetsConst() bool               { return C.Struct(s).Get1(113) }
func (s NodeAnnotation) SetTargetsConst(v bool)           { C.Struct(s).Set1(113, v) }
func (s NodeAnnotation) TargetsEnum() bool                { return C.Struct(s).Get1(114) }
func (s NodeAnnotation) SetTargetsEnum(v bool)            { C.Struct(s).Set1(114, v) }
func (s NodeAnnotation) TargetsEnumerant() bool           { return C.Struct(s).Get1(115) }
func (s NodeAnnotation) SetTargetsEnumerant(v bool)       { C.Struct(s).Set1(115, v) }
func (s NodeAnnotation) TargetsStruct() bool              { return C.Struct(s).Get1(116) }
func (s NodeAnnotation) SetTargetsStruct(v bool)          { C.Struct(s).Set1(116, v) }
func (s NodeAnnotation) TargetsField() bool               { return C.Struct(s).Get1(117) }
func (s NodeAnnotation) SetTargetsField(v bool)           { C.Struct(s).Set1(117, v) }
func (s NodeAnnotation) TargetsUnion() bool               { return C.Struct(s).Get1(118) }
func (s NodeAnnotation) SetTargetsUnion(v bool)           { C.Struct(s).Set1(118, v) }
func (s NodeAnnotation) TargetsGroup() bool               { return C.Struct(s).Get1(119) }
func (s NodeAnnotation) SetTargetsGroup(v bool)           { C.Struct(s).Set1(119, v) }
func (s NodeAnnotation) TargetsInterface() bool           { return C.Struct(s).Get1(120) }
func (s NodeAnnotation) SetTargetsInterface(v bool)       { C.Struct(s).Set1(120, v) }
func (s NodeAnnotation) TargetsMethod() bool              { return C.Struct(s).Get1(121) }
func (s NodeAnnotation) SetTargetsMethod(v bool)          { C.Struct(s).Set1(121, v) }
func (s NodeAnnotation) TargetsParam() bool               { return C.Struct(s).Get1(122) }
func (s NodeAnnotation) SetTargetsParam(v bool)           { C.Struct(s).Set1(122, v) }
func (s NodeAnnotation) TargetsAnnotation() bool          { return C.Struct(s).Get1(123) }
func (s NodeAnnotation) SetTargetsAnnotation(v bool)      { C.Struct(s).Set1(123, v) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s Node) MarshalJSON() (bs []byte, err error) { return }

type Node_List C.PointerList

func NewNodeList(s *C.Segment, sz int) Node_List {
	return Node_List(s.NewCompositeList(C.ObjectSize{DataSize: 40, PointerCount: 6}, sz))
}
func (s Node_List) Len() int             { return C.PointerList(s).Len() }
func (s Node_List) At(i int) Node        { return Node(C.PointerList(s).At(i).ToStruct()) }
func (s Node_List) Set(i int, item Node) { C.PointerList(s).Set(i, C.Object(item)) }

type NodeParameter C.Struct

func NewNodeParameter(s *C.Segment) NodeParameter {
	return NodeParameter(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func NewRootNodeParameter(s *C.Segment) NodeParameter {
	return NodeParameter(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewNodeParameter(s *C.Segment) NodeParameter {
	return NodeParameter(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootNodeParameter(s *C.Segment) NodeParameter { return NodeParameter(s.Root(0).ToStruct()) }
func (s NodeParameter) Name() string                   { return C.Struct(s).GetObject(0).ToText() }
func (s NodeParameter) SetName(v string)               { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s NodeParameter) MarshalJSON() (bs []byte, err error) { return }

type NodeParameter_List C.PointerList

func NewNodeParameterList(s *C.Segment, sz int) NodeParameter_List {
	return NodeParameter_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s NodeParameter_List) Len() int { return C.PointerList(s).Len() }
func (s NodeParameter_List) At(i int) NodeParameter {
	return NodeParameter(C.PointerList(s).At(i).ToStruct())
}
func (s NodeParameter_List) Set(i int, item NodeParameter) { C.PointerList(s).Set(i, C.Object(item)) }

type NodeNestedNode C.Struct

func NewNodeNestedNode(s *C.Segment) NodeNestedNode {
	return NodeNestedNode(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func NewRootNodeNestedNode(s *C.Segment) NodeNestedNode {
	return NodeNestedNode(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func AutoNewNodeNestedNode(s *C.Segment) NodeNestedNode {
	return NodeNestedNode(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func ReadRootNodeNestedNode(s *C.Segment) NodeNestedNode { return NodeNestedNode(s.Root(0).ToStruct()) }
func (s NodeNestedNode) Name() string                    { return C.Struct(s).GetObject(0).ToText() }
func (s NodeNestedNode) SetName(v string)                { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s NodeNestedNode) Id() uint64                      { return C.Struct(s).Get64(0) }
func (s NodeNestedNode) SetId(v uint64)                  { C.Struct(s).Set64(0, v) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s NodeNestedNode) MarshalJSON() (bs []byte, err error) { return }

type NodeNestedNode_List C.PointerList

func NewNodeNestedNodeList(s *C.Segment, sz int) NodeNestedNode_List {
	return NodeNestedNode_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 1}, sz))
}
func (s NodeNestedNode_List) Len() int { return C.PointerList(s).Len() }
func (s NodeNestedNode_List) At(i int) NodeNestedNode {
	return NodeNestedNode(C.PointerList(s).At(i).ToStruct())
}
func (s NodeNestedNode_List) Set(i int, item NodeNestedNode) { C.PointerList(s).Set(i, C.Object(item)) }

type Field C.Struct
type FieldSlot Field
type FieldGroup Field
type FieldOrdinal Field
type Field_Which uint16

const (
	FIELD_SLOT  Field_Which = 0
	FIELD_GROUP Field_Which = 1
)

type FieldOrdinal_Which uint16

const (
	FIELDORDINAL_IMPLICIT FieldOrdinal_Which = 0
	FIELDORDINAL_EXPLICIT FieldOrdinal_Which = 1
)

func NewField(s *C.Segment) Field {
	return Field(s.NewStruct(C.ObjectSize{DataSize: 24, PointerCount: 4}))
}
func NewRootField(s *C.Segment) Field {
	return Field(s.NewRootStruct(C.ObjectSize{DataSize: 24, PointerCount: 4}))
}
func AutoNewField(s *C.Segment) Field {
	return Field(s.NewStructAR(C.ObjectSize{DataSize: 24, PointerCount: 4}))
}
func ReadRootField(s *C.Segment) Field           { return Field(s.Root(0).ToStruct()) }
func (s Field) Which() Field_Which               { return Field_Which(C.Struct(s).Get16(8)) }
func (s Field) Name() string                     { return C.Struct(s).GetObject(0).ToText() }
func (s Field) SetName(v string)                 { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s Field) CodeOrder() uint16                { return C.Struct(s).Get16(0) }
func (s Field) SetCodeOrder(v uint16)            { C.Struct(s).Set16(0, v) }
func (s Field) Annotations() Annotation_List     { return Annotation_List(C.Struct(s).GetObject(1)) }
func (s Field) SetAnnotations(v Annotation_List) { C.Struct(s).SetObject(1, C.Object(v)) }
func (s Field) DiscriminantValue() uint16        { return C.Struct(s).Get16(2) ^ 65535 }
func (s Field) SetDiscriminantValue(v uint16)    { C.Struct(s).Set16(2, v^65535) }
func (s Field) Slot() FieldSlot                  { return FieldSlot(s) }
func (s Field) SetSlot()                         { C.Struct(s).Set16(8, 0) }
func (s FieldSlot) Offset() uint32               { return C.Struct(s).Get32(4) }
func (s FieldSlot) SetOffset(v uint32)           { C.Struct(s).Set32(4, v) }
func (s FieldSlot) Type() Type                   { return Type(C.Struct(s).GetObject(2).ToStruct()) }
func (s FieldSlot) SetType(v Type)               { C.Struct(s).SetObject(2, C.Object(v)) }
func (s FieldSlot) DefaultValue() Value          { return Value(C.Struct(s).GetObject(3).ToStruct()) }
func (s FieldSlot) SetDefaultValue(v Value)      { C.Struct(s).SetObject(3, C.Object(v)) }
func (s FieldSlot) HadExplicitDefault() bool     { return C.Struct(s).Get1(128) }
func (s FieldSlot) SetHadExplicitDefault(v bool) { C.Struct(s).Set1(128, v) }
func (s Field) Group() FieldGroup                { return FieldGroup(s) }
func (s Field) SetGroup()                        { C.Struct(s).Set16(8, 1) }
func (s FieldGroup) TypeId() uint64              { return C.Struct(s).Get64(16) }
func (s FieldGroup) SetTypeId(v uint64)          { C.Struct(s).Set64(16, v) }
func (s Field) Ordinal() FieldOrdinal            { return FieldOrdinal(s) }
func (s FieldOrdinal) Which() FieldOrdinal_Which { return FieldOrdinal_Which(C.Struct(s).Get16(10)) }
func (s FieldOrdinal) SetImplicit()              { C.Struct(s).Set16(10, 0) }
func (s FieldOrdinal) Explicit() uint16          { return C.Struct(s).Get16(12) }
func (s FieldOrdinal) SetExplicit(v uint16)      { C.Struct(s).Set16(10, 1); C.Struct(s).Set16(12, v) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s Field) MarshalJSON() (bs []byte, err error) { return }

type Field_List C.PointerList

func NewFieldList(s *C.Segment, sz int) Field_List {
	return Field_List(s.NewCompositeList(C.ObjectSize{DataSize: 24, PointerCount: 4}, sz))
}
func (s Field_List) Len() int              { return C.PointerList(s).Len() }
func (s Field_List) At(i int) Field        { return Field(C.PointerList(s).At(i).ToStruct()) }
func (s Field_List) Set(i int, item Field) { C.PointerList(s).Set(i, C.Object(item)) }

type Enumerant C.Struct

func NewEnumerant(s *C.Segment) Enumerant {
	return Enumerant(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 2}))
}
func NewRootEnumerant(s *C.Segment) Enumerant {
	return Enumerant(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 2}))
}
func AutoNewEnumerant(s *C.Segment) Enumerant {
	return Enumerant(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 2}))
}
func ReadRootEnumerant(s *C.Segment) Enumerant       { return Enumerant(s.Root(0).ToStruct()) }
func (s Enumerant) Name() string                     { return C.Struct(s).GetObject(0).ToText() }
func (s Enumerant) SetName(v string)                 { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s Enumerant) CodeOrder() uint16                { return C.Struct(s).Get16(0) }
func (s Enumerant) SetCodeOrder(v uint16)            { C.Struct(s).Set16(0, v) }
func (s Enumerant) Annotations() Annotation_List     { return Annotation_List(C.Struct(s).GetObject(1)) }
func (s Enumerant) SetAnnotations(v Annotation_List) { C.Struct(s).SetObject(1, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s Enumerant) MarshalJSON() (bs []byte, err error) { return }

type Enumerant_List C.PointerList

func NewEnumerantList(s *C.Segment, sz int) Enumerant_List {
	return Enumerant_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 2}, sz))
}
func (s Enumerant_List) Len() int                  { return C.PointerList(s).Len() }
func (s Enumerant_List) At(i int) Enumerant        { return Enumerant(C.PointerList(s).At(i).ToStruct()) }
func (s Enumerant_List) Set(i int, item Enumerant) { C.PointerList(s).Set(i, C.Object(item)) }

type Superclass C.Struct

func NewSuperclass(s *C.Segment) Superclass {
	return Superclass(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func NewRootSuperclass(s *C.Segment) Superclass {
	return Superclass(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func AutoNewSuperclass(s *C.Segment) Superclass {
	return Superclass(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func ReadRootSuperclass(s *C.Segment) Superclass { return Superclass(s.Root(0).ToStruct()) }
func (s Superclass) Id() uint64                  { return C.Struct(s).Get64(0) }
func (s Superclass) SetId(v uint64)              { C.Struct(s).Set64(0, v) }
func (s Superclass) Brand() Brand                { return Brand(C.Struct(s).GetObject(0).ToStruct()) }
func (s Superclass) SetBrand(v Brand)            { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s Superclass) MarshalJSON() (bs []byte, err error) { return }

type Superclass_List C.PointerList

func NewSuperclassList(s *C.Segment, sz int) Superclass_List {
	return Superclass_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 1}, sz))
}
func (s Superclass_List) Len() int                   { return C.PointerList(s).Len() }
func (s Superclass_List) At(i int) Superclass        { return Superclass(C.PointerList(s).At(i).ToStruct()) }
func (s Superclass_List) Set(i int, item Superclass) { C.PointerList(s).Set(i, C.Object(item)) }

type Method C.Struct

func NewMethod(s *C.Segment) Method {
	return Method(s.NewStruct(C.ObjectSize{DataSize: 24, PointerCount: 5}))
}
func NewRootMethod(s *C.Segment) Method {
	return Method(s.NewRootStruct(C.ObjectSize{DataSize: 24, PointerCount: 5}))
}
func AutoNewMethod(s *C.Segment) Method {
	return Method(s.NewStructAR(C.ObjectSize{DataSize: 24, PointerCount: 5}))
}
func ReadRootMethod(s *C.Segment) Method { return Method(s.Root(0).ToStruct()) }
func (s Method) Name() string            { return C.Struct(s).GetObject(0).ToText() }
func (s Method) SetName(v string)        { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s Method) CodeOrder() uint16       { return C.Struct(s).Get16(0) }
func (s Method) SetCodeOrder(v uint16)   { C.Struct(s).Set16(0, v) }
func (s Method) ImplicitParameters() NodeParameter_List {
	return NodeParameter_List(C.Struct(s).GetObject(4))
}
func (s Method) SetImplicitParameters(v NodeParameter_List) { C.Struct(s).SetObject(4, C.Object(v)) }
func (s Method) ParamStructType() uint64                    { return C.Struct(s).Get64(8) }
func (s Method) SetParamStructType(v uint64)                { C.Struct(s).Set64(8, v) }
func (s Method) ParamBrand() Brand                          { return Brand(C.Struct(s).GetObject(2).ToStruct()) }
func (s Method) SetParamBrand(v Brand)                      { C.Struct(s).SetObject(2, C.Object(v)) }
func (s Method) ResultStructType() uint64                   { return C.Struct(s).Get64(16) }
func (s Method) SetResultStructType(v uint64)               { C.Struct(s).Set64(16, v) }
func (s Method) ResultBrand() Brand                         { return Brand(C.Struct(s).GetObject(3).ToStruct()) }
func (s Method) SetResultBrand(v Brand)                     { C.Struct(s).SetObject(3, C.Object(v)) }
func (s Method) Annotations() Annotation_List               { return Annotation_List(C.Struct(s).GetObject(1)) }
func (s Method) SetAnnotations(v Annotation_List)           { C.Struct(s).SetObject(1, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s Method) MarshalJSON() (bs []byte, err error) { return }

type Method_List C.PointerList

func NewMethodList(s *C.Segment, sz int) Method_List {
	return Method_List(s.NewCompositeList(C.ObjectSize{DataSize: 24, PointerCount: 5}, sz))
}
func (s Method_List) Len() int               { return C.PointerList(s).Len() }
func (s Method_List) At(i int) Method        { return Method(C.PointerList(s).At(i).ToStruct()) }
func (s Method_List) Set(i int, item Method) { C.PointerList(s).Set(i, C.Object(item)) }

type Type C.Struct
type TypeList Type
type TypeEnum Type
type TypeStruct Type
type TypeInterface Type
type TypeAnyPointer Type
type TypeAnyPointerParameter Type
type TypeAnyPointerImplicitMethodParameter Type
type Type_Which uint16

const (
	TYPE_VOID       Type_Which = 0
	TYPE_BOOL       Type_Which = 1
	TYPE_INT8       Type_Which = 2
	TYPE_INT16      Type_Which = 3
	TYPE_INT32      Type_Which = 4
	TYPE_INT64      Type_Which = 5
	TYPE_UINT8      Type_Which = 6
	TYPE_UINT16     Type_Which = 7
	TYPE_UINT32     Type_Which = 8
	TYPE_UINT64     Type_Which = 9
	TYPE_FLOAT32    Type_Which = 10
	TYPE_FLOAT64    Type_Which = 11
	TYPE_TEXT       Type_Which = 12
	TYPE_DATA       Type_Which = 13
	TYPE_LIST       Type_Which = 14
	TYPE_ENUM       Type_Which = 15
	TYPE_STRUCT     Type_Which = 16
	TYPE_INTERFACE  Type_Which = 17
	TYPE_ANYPOINTER Type_Which = 18
)

type TypeAnyPointer_Which uint16

const (
	TYPEANYPOINTER_UNCONSTRAINED           TypeAnyPointer_Which = 0
	TYPEANYPOINTER_PARAMETER               TypeAnyPointer_Which = 1
	TYPEANYPOINTER_IMPLICITMETHODPARAMETER TypeAnyPointer_Which = 2
)

func NewType(s *C.Segment) Type { return Type(s.NewStruct(C.ObjectSize{DataSize: 24, PointerCount: 1})) }
func NewRootType(s *C.Segment) Type {
	return Type(s.NewRootStruct(C.ObjectSize{DataSize: 24, PointerCount: 1}))
}
func AutoNewType(s *C.Segment) Type {
	return Type(s.NewStructAR(C.ObjectSize{DataSize: 24, PointerCount: 1}))
}
func ReadRootType(s *C.Segment) Type       { return Type(s.Root(0).ToStruct()) }
func (s Type) Which() Type_Which           { return Type_Which(C.Struct(s).Get16(0)) }
func (s Type) SetVoid()                    { C.Struct(s).Set16(0, 0) }
func (s Type) SetBool()                    { C.Struct(s).Set16(0, 1) }
func (s Type) SetInt8()                    { C.Struct(s).Set16(0, 2) }
func (s Type) SetInt16()                   { C.Struct(s).Set16(0, 3) }
func (s Type) SetInt32()                   { C.Struct(s).Set16(0, 4) }
func (s Type) SetInt64()                   { C.Struct(s).Set16(0, 5) }
func (s Type) SetUint8()                   { C.Struct(s).Set16(0, 6) }
func (s Type) SetUint16()                  { C.Struct(s).Set16(0, 7) }
func (s Type) SetUint32()                  { C.Struct(s).Set16(0, 8) }
func (s Type) SetUint64()                  { C.Struct(s).Set16(0, 9) }
func (s Type) SetFloat32()                 { C.Struct(s).Set16(0, 10) }
func (s Type) SetFloat64()                 { C.Struct(s).Set16(0, 11) }
func (s Type) SetText()                    { C.Struct(s).Set16(0, 12) }
func (s Type) SetData()                    { C.Struct(s).Set16(0, 13) }
func (s Type) List() TypeList              { return TypeList(s) }
func (s Type) SetList()                    { C.Struct(s).Set16(0, 14) }
func (s TypeList) ElementType() Type       { return Type(C.Struct(s).GetObject(0).ToStruct()) }
func (s TypeList) SetElementType(v Type)   { C.Struct(s).SetObject(0, C.Object(v)) }
func (s Type) Enum() TypeEnum              { return TypeEnum(s) }
func (s Type) SetEnum()                    { C.Struct(s).Set16(0, 15) }
func (s TypeEnum) TypeId() uint64          { return C.Struct(s).Get64(8) }
func (s TypeEnum) SetTypeId(v uint64)      { C.Struct(s).Set64(8, v) }
func (s TypeEnum) Brand() Brand            { return Brand(C.Struct(s).GetObject(0).ToStruct()) }
func (s TypeEnum) SetBrand(v Brand)        { C.Struct(s).SetObject(0, C.Object(v)) }
func (s Type) Struct() TypeStruct          { return TypeStruct(s) }
func (s Type) SetStruct()                  { C.Struct(s).Set16(0, 16) }
func (s TypeStruct) TypeId() uint64        { return C.Struct(s).Get64(8) }
func (s TypeStruct) SetTypeId(v uint64)    { C.Struct(s).Set64(8, v) }
func (s TypeStruct) Brand() Brand          { return Brand(C.Struct(s).GetObject(0).ToStruct()) }
func (s TypeStruct) SetBrand(v Brand)      { C.Struct(s).SetObject(0, C.Object(v)) }
func (s Type) Interface() TypeInterface    { return TypeInterface(s) }
func (s Type) SetInterface()               { C.Struct(s).Set16(0, 17) }
func (s TypeInterface) TypeId() uint64     { return C.Struct(s).Get64(8) }
func (s TypeInterface) SetTypeId(v uint64) { C.Struct(s).Set64(8, v) }
func (s TypeInterface) Brand() Brand       { return Brand(C.Struct(s).GetObject(0).ToStruct()) }
func (s TypeInterface) SetBrand(v Brand)   { C.Struct(s).SetObject(0, C.Object(v)) }
func (s Type) AnyPointer() TypeAnyPointer  { return TypeAnyPointer(s) }
func (s Type) SetAnyPointer()              { C.Struct(s).Set16(0, 18) }
func (s TypeAnyPointer) Which() TypeAnyPointer_Which {
	return TypeAnyPointer_Which(C.Struct(s).Get16(8))
}
func (s TypeAnyPointer) SetUnconstrained()                   { C.Struct(s).Set16(8, 0) }
func (s TypeAnyPointer) Parameter() TypeAnyPointerParameter  { return TypeAnyPointerParameter(s) }
func (s TypeAnyPointer) SetParameter()                       { C.Struct(s).Set16(8, 1) }
func (s TypeAnyPointerParameter) ScopeId() uint64            { return C.Struct(s).Get64(16) }
func (s TypeAnyPointerParameter) SetScopeId(v uint64)        { C.Struct(s).Set64(16, v) }
func (s TypeAnyPointerParameter) ParameterIndex() uint16     { return C.Struct(s).Get16(10) }
func (s TypeAnyPointerParameter) SetParameterIndex(v uint16) { C.Struct(s).Set16(10, v) }
func (s TypeAnyPointer) ImplicitMethodParameter() TypeAnyPointerImplicitMethodParameter {
	return TypeAnyPointerImplicitMethodParameter(s)
}
func (s TypeAnyPointer) SetImplicitMethodParameter()                       { C.Struct(s).Set16(8, 2) }
func (s TypeAnyPointerImplicitMethodParameter) ParameterIndex() uint16     { return C.Struct(s).Get16(10) }
func (s TypeAnyPointerImplicitMethodParameter) SetParameterIndex(v uint16) { C.Struct(s).Set16(10, v) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s Type) MarshalJSON() (bs []byte, err error) { return }

type Type_List C.PointerList

func NewTypeList(s *C.Segment, sz int) Type_List {
	return Type_List(s.NewCompositeList(C.ObjectSize{DataSize: 24, PointerCount: 1}, sz))
}
func (s Type_List) Len() int             { return C.PointerList(s).Len() }
func (s Type_List) At(i int) Type        { return Type(C.PointerList(s).At(i).ToStruct()) }
func (s Type_List) Set(i int, item Type) { C.PointerList(s).Set(i, C.Object(item)) }

type Brand C.Struct

func NewBrand(s *C.Segment) Brand {
	return Brand(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func NewRootBrand(s *C.Segment) Brand {
	return Brand(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewBrand(s *C.Segment) Brand {
	return Brand(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootBrand(s *C.Segment) Brand      { return Brand(s.Root(0).ToStruct()) }
func (s Brand) Scopes() BrandScope_List     { return BrandScope_List(C.Struct(s).GetObject(0)) }
func (s Brand) SetScopes(v BrandScope_List) { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s Brand) MarshalJSON() (bs []byte, err error) { return }

type Brand_List C.PointerList

func NewBrandList(s *C.Segment, sz int) Brand_List {
	return Brand_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s Brand_List) Len() int              { return C.PointerList(s).Len() }
func (s Brand_List) At(i int) Brand        { return Brand(C.PointerList(s).At(i).ToStruct()) }
func (s Brand_List) Set(i int, item Brand) { C.PointerList(s).Set(i, C.Object(item)) }

type BrandScope C.Struct
type BrandScope_Which uint16

const (
	BRANDSCOPE_BIND    BrandScope_Which = 0
	BRANDSCOPE_INHERIT BrandScope_Which = 1
)

func NewBrandScope(s *C.Segment) BrandScope {
	return BrandScope(s.NewStruct(C.ObjectSize{DataSize: 16, PointerCount: 1}))
}
func NewRootBrandScope(s *C.Segment) BrandScope {
	return BrandScope(s.NewRootStruct(C.ObjectSize{DataSize: 16, PointerCount: 1}))
}
func AutoNewBrandScope(s *C.Segment) BrandScope {
	return BrandScope(s.NewStructAR(C.ObjectSize{DataSize: 16, PointerCount: 1}))
}
func ReadRootBrandScope(s *C.Segment) BrandScope { return BrandScope(s.Root(0).ToStruct()) }
func (s BrandScope) Which() BrandScope_Which     { return BrandScope_Which(C.Struct(s).Get16(8)) }
func (s BrandScope) ScopeId() uint64             { return C.Struct(s).Get64(0) }
func (s BrandScope) SetScopeId(v uint64)         { C.Struct(s).Set64(0, v) }
func (s BrandScope) Bind() BrandBinding_List     { return BrandBinding_List(C.Struct(s).GetObject(0)) }
func (s BrandScope) SetBind(v BrandBinding_List) {
	C.Struct(s).Set16(8, 0)
	C.Struct(s).SetObject(0, C.Object(v))
}
func (s BrandScope) SetInherit() { C.Struct(s).Set16(8, 1) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s BrandScope) MarshalJSON() (bs []byte, err error) { return }

type BrandScope_List C.PointerList

func NewBrandScopeList(s *C.Segment, sz int) BrandScope_List {
	return BrandScope_List(s.NewCompositeList(C.ObjectSize{DataSize: 16, PointerCount: 1}, sz))
}
func (s BrandScope_List) Len() int                   { return C.PointerList(s).Len() }
func (s BrandScope_List) At(i int) BrandScope        { return BrandScope(C.PointerList(s).At(i).ToStruct()) }
func (s BrandScope_List) Set(i int, item BrandScope) { C.PointerList(s).Set(i, C.Object(item)) }

type BrandBinding C.Struct
type BrandBinding_Which uint16

const (
	BRANDBINDING_UNBOUND BrandBinding_Which = 0
	BRANDBINDING_TYPE    BrandBinding_Which = 1
)

func NewBrandBinding(s *C.Segment) BrandBinding {
	return BrandBinding(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func NewRootBrandBinding(s *C.Segment) BrandBinding {
	return BrandBinding(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func AutoNewBrandBinding(s *C.Segment) BrandBinding {
	return BrandBinding(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func ReadRootBrandBinding(s *C.Segment) BrandBinding { return BrandBinding(s.Root(0).ToStruct()) }
func (s BrandBinding) Which() BrandBinding_Which     { return BrandBinding_Which(C.Struct(s).Get16(0)) }
func (s BrandBinding) SetUnbound()                   { C.Struct(s).Set16(0, 0) }
func (s BrandBinding) Type() Type                    { return Type(C.Struct(s).GetObject(0).ToStruct()) }
func (s BrandBinding) SetType(v Type)                { C.Struct(s).Set16(0, 1); C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s BrandBinding) MarshalJSON() (bs []byte, err error) { return }

type BrandBinding_List C.PointerList

func NewBrandBindingList(s *C.Segment, sz int) BrandBinding_List {
	return BrandBinding_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 1}, sz))
}
func (s BrandBinding_List) Len() int { return C.PointerList(s).Len() }
func (s BrandBinding_List) At(i int) BrandBinding {
	return BrandBinding(C.PointerList(s).At(i).ToStruct())
}
func (s BrandBinding_List) Set(i int, item BrandBinding) { C.PointerList(s).Set(i, C.Object(item)) }

type Value C.Struct
type Value_Which uint16

const (
	VALUE_VOID       Value_Which = 0
	VALUE_BOOL       Value_Which = 1
	VALUE_INT8       Value_Which = 2
	VALUE_INT16      Value_Which = 3
	VALUE_INT32      Value_Which = 4
	VALUE_INT64      Value_Which = 5
	VALUE_UINT8      Value_Which = 6
	VALUE_UINT16     Value_Which = 7
	VALUE_UINT32     Value_Which = 8
	VALUE_UINT64     Value_Which = 9
	VALUE_FLOAT32    Value_Which = 10
	VALUE_FLOAT64    Value_Which = 11
	VALUE_TEXT       Value_Which = 12
	VALUE_DATA       Value_Which = 13
	VALUE_LIST       Value_Which = 14
	VALUE_ENUM       Value_Which = 15
	VALUE_STRUCT     Value_Which = 16
	VALUE_INTERFACE  Value_Which = 17
	VALUE_ANYPOINTER Value_Which = 18
)

func NewValue(s *C.Segment) Value {
	return Value(s.NewStruct(C.ObjectSize{DataSize: 16, PointerCount: 1}))
}
func NewRootValue(s *C.Segment) Value {
	return Value(s.NewRootStruct(C.ObjectSize{DataSize: 16, PointerCount: 1}))
}
func AutoNewValue(s *C.Segment) Value {
	return Value(s.NewStructAR(C.ObjectSize{DataSize: 16, PointerCount: 1}))
}
func ReadRootValue(s *C.Segment) Value { return Value(s.Root(0).ToStruct()) }
func (s Value) Which() Value_Which     { return Value_Which(C.Struct(s).Get16(0)) }
func (s Value) SetVoid()               { C.Struct(s).Set16(0, 0) }
func (s Value) Bool() bool             { return C.Struct(s).Get1(16) }
func (s Value) SetBool(v bool)         { C.Struct(s).Set16(0, 1); C.Struct(s).Set1(16, v) }
func (s Value) Int8() int8             { return int8(C.Struct(s).Get8(2)) }
func (s Value) SetInt8(v int8)         { C.Struct(s).Set16(0, 2); C.Struct(s).Set8(2, uint8(v)) }
func (s Value) Int16() int16           { return int16(C.Struct(s).Get16(2)) }
func (s Value) SetInt16(v int16)       { C.Struct(s).Set16(0, 3); C.Struct(s).Set16(2, uint16(v)) }
func (s Value) Int32() int32           { return int32(C.Struct(s).Get32(4)) }
func (s Value) SetInt32(v int32)       { C.Struct(s).Set16(0, 4); C.Struct(s).Set32(4, uint32(v)) }
func (s Value) Int64() int64           { return int64(C.Struct(s).Get64(8)) }
func (s Value) SetInt64(v int64)       { C.Struct(s).Set16(0, 5); C.Struct(s).Set64(8, uint64(v)) }
func (s Value) Uint8() uint8           { return C.Struct(s).Get8(2) }
func (s Value) SetUint8(v uint8)       { C.Struct(s).Set16(0, 6); C.Struct(s).Set8(2, v) }
func (s Value) Uint16() uint16         { return C.Struct(s).Get16(2) }
func (s Value) SetUint16(v uint16)     { C.Struct(s).Set16(0, 7); C.Struct(s).Set16(2, v) }
func (s Value) Uint32() uint32         { return C.Struct(s).Get32(4) }
func (s Value) SetUint32(v uint32)     { C.Struct(s).Set16(0, 8); C.Struct(s).Set32(4, v) }
func (s Value) Uint64() uint64         { return C.Struct(s).Get64(8) }
func (s Value) SetUint64(v uint64)     { C.Struct(s).Set16(0, 9); C.Struct(s).Set64(8, v) }
func (s Value) Float32() float32       { return math.Float32frombits(C.Struct(s).Get32(4)) }
func (s Value) SetFloat32(v float32) {
	C.Struct(s).Set16(0, 10)
	C.Struct(s).Set32(4, math.Float32bits(v))
}
func (s Value) Float64() float64 { return math.Float64frombits(C.Struct(s).Get64(8)) }
func (s Value) SetFloat64(v float64) {
	C.Struct(s).Set16(0, 11)
	C.Struct(s).Set64(8, math.Float64bits(v))
}
func (s Value) Text() string { return C.Struct(s).GetObject(0).ToText() }
func (s Value) SetText(v string) {
	C.Struct(s).Set16(0, 12)
	C.Struct(s).SetObject(0, s.Segment.NewText(v))
}
func (s Value) Data() []byte { return C.Struct(s).GetObject(0).ToData() }
func (s Value) SetData(v []byte) {
	C.Struct(s).Set16(0, 13)
	C.Struct(s).SetObject(0, s.Segment.NewData(v))
}
func (s Value) List() C.Object           { return C.Struct(s).GetObject(0) }
func (s Value) SetList(v C.Object)       { C.Struct(s).Set16(0, 14); C.Struct(s).SetObject(0, v) }
func (s Value) Enum() uint16             { return C.Struct(s).Get16(2) }
func (s Value) SetEnum(v uint16)         { C.Struct(s).Set16(0, 15); C.Struct(s).Set16(2, v) }
func (s Value) Struct() C.Object         { return C.Struct(s).GetObject(0) }
func (s Value) SetStruct(v C.Object)     { C.Struct(s).Set16(0, 16); C.Struct(s).SetObject(0, v) }
func (s Value) SetInterface()            { C.Struct(s).Set16(0, 17) }
func (s Value) AnyPointer() C.Object     { return C.Struct(s).GetObject(0) }
func (s Value) SetAnyPointer(v C.Object) { C.Struct(s).Set16(0, 18); C.Struct(s).SetObject(0, v) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s Value) MarshalJSON() (bs []byte, err error) { return }

type Value_List C.PointerList

func NewValueList(s *C.Segment, sz int) Value_List {
	return Value_List(s.NewCompositeList(C.ObjectSize{DataSize: 16, PointerCount: 1}, sz))
}
func (s Value_List) Len() int              { return C.PointerList(s).Len() }
func (s Value_List) At(i int) Value        { return Value(C.PointerList(s).At(i).ToStruct()) }
func (s Value_List) Set(i int, item Value) { C.PointerList(s).Set(i, C.Object(item)) }

type Annotation C.Struct

func NewAnnotation(s *C.Segment) Annotation {
	return Annotation(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 2}))
}
func NewRootAnnotation(s *C.Segment) Annotation {
	return Annotation(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 2}))
}
func AutoNewAnnotation(s *C.Segment) Annotation {
	return Annotation(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 2}))
}
func ReadRootAnnotation(s *C.Segment) Annotation { return Annotation(s.Root(0).ToStruct()) }
func (s Annotation) Id() uint64                  { return C.Struct(s).Get64(0) }
func (s Annotation) SetId(v uint64)              { C.Struct(s).Set64(0, v) }
func (s Annotation) Brand() Brand                { return Brand(C.Struct(s).GetObject(1).ToStruct()) }
func (s Annotation) SetBrand(v Brand)            { C.Struct(s).SetObject(1, C.Object(v)) }
func (s Annotation) Value() Value                { return Value(C.Struct(s).GetObject(0).ToStruct()) }
func (s Annotation) SetValue(v Value)            { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s Annotation) MarshalJSON() (bs []byte, err error) { return }

type Annotation_List C.PointerList

func NewAnnotationList(s *C.Segment, sz int) Annotation_List {
	return Annotation_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 2}, sz))
}
func (s Annotation_List) Len() int                   { return C.PointerList(s).Len() }
func (s Annotation_List) At(i int) Annotation        { return Annotation(C.PointerList(s).At(i).ToStruct()) }
func (s Annotation_List) Set(i int, item Annotation) { C.PointerList(s).Set(i, C.Object(item)) }

type ElementSize uint16

const (
	ELEMENTSIZE_EMPTY           ElementSize = 0
	ELEMENTSIZE_BIT             ElementSize = 1
	ELEMENTSIZE_BYTE            ElementSize = 2
	ELEMENTSIZE_TWOBYTES        ElementSize = 3
	ELEMENTSIZE_FOURBYTES       ElementSize = 4
	ELEMENTSIZE_EIGHTBYTES      ElementSize = 5
	ELEMENTSIZE_POINTER         ElementSize = 6
	ELEMENTSIZE_INLINECOMPOSITE ElementSize = 7
)

func (c ElementSize) String() string {
	switch c {
	case ELEMENTSIZE_EMPTY:
		return "empty"
	case ELEMENTSIZE_BIT:
		return "bit"
	case ELEMENTSIZE_BYTE:
		return "byte"
	case ELEMENTSIZE_TWOBYTES:
		return "twoBytes"
	case ELEMENTSIZE_FOURBYTES:
		return "fourBytes"
	case ELEMENTSIZE_EIGHTBYTES:
		return "eightBytes"
	case ELEMENTSIZE_POINTER:
		return "pointer"
	case ELEMENTSIZE_INLINECOMPOSITE:
		return "inlineComposite"
	default:
		return ""
	}
}

func ElementSizeFromString(c string) ElementSize {
	switch c {
	case "empty":
		return ELEMENTSIZE_EMPTY
	case "bit":
		return ELEMENTSIZE_BIT
	case "byte":
		return ELEMENTSIZE_BYTE
	case "twoBytes":
		return ELEMENTSIZE_TWOBYTES
	case "fourBytes":
		return ELEMENTSIZE_FOURBYTES
	case "eightBytes":
		return ELEMENTSIZE_EIGHTBYTES
	case "pointer":
		return ELEMENTSIZE_POINTER
	case "inlineComposite":
		return ELEMENTSIZE_INLINECOMPOSITE
	default:
		return 0
	}
}

type ElementSize_List C.PointerList

func NewElementSizeList(s *C.Segment, sz int) ElementSize_List {
	return ElementSize_List(s.NewUInt16List(sz))
}
func (s ElementSize_List) Len() int             { return C.UInt16List(s).Len() }
func (s ElementSize_List) At(i int) ElementSize { return ElementSize(C.UInt16List(s).At(i)) }

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s ElementSize) MarshalJSON() (bs []byte, err error) { return }

type CodeGeneratorRequest C.Struct

func NewCodeGeneratorRequest(s *C.Segment) CodeGeneratorRequest {
	return CodeGeneratorRequest(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 2}))
}
func NewRootCodeGeneratorRequest(s *C.Segment) CodeGeneratorRequest {
	return CodeGeneratorRequest(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 2}))
}
func AutoNewCodeGeneratorRequest(s *C.Segment) CodeGeneratorRequest {
	return CodeGeneratorRequest(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 2}))
}
func ReadRootCodeGeneratorRequest(s *C.Segment) CodeGeneratorRequest {
	return CodeGeneratorRequest(s.Root(0).ToStruct())
}
func (s CodeGeneratorRequest) Nodes() Node_List     { return Node_List(C.Struct(s).GetObject(0)) }
func (s CodeGeneratorRequest) SetNodes(v Node_List) { C.Struct(s).SetObject(0, C.Object(v)) }
func (s CodeGeneratorRequest) RequestedFiles() CodeGeneratorRequestRequestedFile_List {
	return CodeGeneratorRequestRequestedFile_List(C.Struct(s).GetObject(1))
}
func (s CodeGeneratorRequest) SetRequestedFiles(v CodeGeneratorRequestRequestedFile_List) {
	C.Struct(s).SetObject(1, C.Object(v))
}

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s CodeGeneratorRequest) MarshalJSON() (bs []byte, err error) { return }

type CodeGeneratorRequest_List C.PointerList

func NewCodeGeneratorRequestList(s *C.Segment, sz int) CodeGeneratorRequest_List {
	return CodeGeneratorRequest_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 2}, sz))
}
func (s CodeGeneratorRequest_List) Len() int { return C.PointerList(s).Len() }
func (s CodeGeneratorRequest_List) At(i int) CodeGeneratorRequest {
	return CodeGeneratorRequest(C.PointerList(s).At(i).ToStruct())
}
func (s CodeGeneratorRequest_List) Set(i int, item CodeGeneratorRequest) {
	C.PointerList(s).Set(i, C.Object(item))
}

type CodeGeneratorRequestRequestedFile C.Struct

func NewCodeGeneratorRequestRequestedFile(s *C.Segment) CodeGeneratorRequestRequestedFile {
	return CodeGeneratorRequestRequestedFile(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 2}))
}
func NewRootCodeGeneratorRequestRequestedFile(s *C.Segment) CodeGeneratorRequestRequestedFile {
	return CodeGeneratorRequestRequestedFile(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 2}))
}
func AutoNewCodeGeneratorRequestRequestedFile(s *C.Segment) CodeGeneratorRequestRequestedFile {
	return CodeGeneratorRequestRequestedFile(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 2}))
}
func ReadRootCodeGeneratorRequestRequestedFile(s *C.Segment) CodeGeneratorRequestRequestedFile {
	return CodeGeneratorRequestRequestedFile(s.Root(0).ToStruct())
}
func (s CodeGeneratorRequestRequestedFile) Id() uint64       { return C.Struct(s).Get64(0) }
func (s CodeGeneratorRequestRequestedFile) SetId(v uint64)   { C.Struct(s).Set64(0, v) }
func (s CodeGeneratorRequestRequestedFile) Filename() string { return C.Struct(s).GetObject(0).ToText() }
func (s CodeGeneratorRequestRequestedFile) SetFilename(v string) {
	C.Struct(s).SetObject(0, s.Segment.NewText(v))
}
func (s CodeGeneratorRequestRequestedFile) Imports() CodeGeneratorRequestRequestedFileImport_List {
	return CodeGeneratorRequestRequestedFileImport_List(C.Struct(s).GetObject(1))
}
func (s CodeGeneratorRequestRequestedFile) SetImports(v CodeGeneratorRequestRequestedFileImport_List) {
	C.Struct(s).SetObject(1, C.Object(v))
}

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s CodeGeneratorRequestRequestedFile) MarshalJSON() (bs []byte, err error) { return }

type CodeGeneratorRequestRequestedFile_List C.PointerList

func NewCodeGeneratorRequestRequestedFileList(s *C.Segment, sz int) CodeGeneratorRequestRequestedFile_List {
	return CodeGeneratorRequestRequestedFile_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 2}, sz))
}
func (s CodeGeneratorRequestRequestedFile_List) Len() int { return C.PointerList(s).Len() }
func (s CodeGeneratorRequestRequestedFile_List) At(i int) CodeGeneratorRequestRequestedFile {
	return CodeGeneratorRequestRequestedFile(C.PointerList(s).At(i).ToStruct())
}
func (s CodeGeneratorRequestRequestedFile_List) Set(i int, item CodeGeneratorRequestRequestedFile) {
	C.PointerList(s).Set(i, C.Object(item))
}

type CodeGeneratorRequestRequestedFileImport C.Struct

func NewCodeGeneratorRequestRequestedFileImport(s *C.Segment) CodeGeneratorRequestRequestedFileImport {
	return CodeGeneratorRequestRequestedFileImport(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func NewRootCodeGeneratorRequestRequestedFileImport(s *C.Segment) CodeGeneratorRequestRequestedFileImport {
	return CodeGeneratorRequestRequestedFileImport(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func AutoNewCodeGeneratorRequestRequestedFileImport(s *C.Segment) CodeGeneratorRequestRequestedFileImport {
	return CodeGeneratorRequestRequestedFileImport(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func ReadRootCodeGeneratorRequestRequestedFileImport(s *C.Segment) CodeGeneratorRequestRequestedFileImport {
	return CodeGeneratorRequestRequestedFileImport(s.Root(0).ToStruct())
}
func (s CodeGeneratorRequestRequestedFileImport) Id() uint64     { return C.Struct(s).Get64(0) }
func (s CodeGeneratorRequestRequestedFileImport) SetId(v uint64) { C.Struct(s).Set64(0, v) }
func (s CodeGeneratorRequestRequestedFileImport) Name() string {
	return C.Struct(s).GetObject(0).ToText()
}
func (s CodeGeneratorRequestRequestedFileImport) SetName(v string) {
	C.Struct(s).SetObject(0, s.Segment.NewText(v))
}

// capn.JSON_enabled == false so we stub MarshallJSON().
func (s CodeGeneratorRequestRequestedFileImport) MarshalJSON() (bs []byte, err error) { return }

type CodeGeneratorRequestRequestedFileImport_List C.PointerList

func NewCodeGeneratorRequestRequestedFileImportList(s *C.Segment, sz int) CodeGeneratorRequestRequestedFileImport_List {
	return CodeGeneratorRequestRequestedFileImport_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 1}, sz))
}
func (s CodeGeneratorRequestRequestedFileImport_List) Len() int { return C.PointerList(s).Len() }
func (s CodeGeneratorRequestRequestedFileImport_List) At(i int) CodeGeneratorRequestRequestedFileImport {
	return CodeGeneratorRequestRequestedFileImport(C.PointerList(s).At(i).ToStruct())
}
func (s CodeGeneratorRequestRequestedFileImport_List) Set(i int, item CodeGeneratorRequestRequestedFileImport) {
	C.PointerList(s).Set(i, C.Object(item))
}
