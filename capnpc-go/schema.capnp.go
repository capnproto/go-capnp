package main

// AUTO GENERATED - DO NOT EDIT

import (
	C "github.com/glycerine/go-capnproto"
	"math"
)

const (
	Field_noDiscriminant = uint16(65535)
)

type Node C.Struct
type Node_struct Node
type Node_enum Node
type Node_interface Node
type Node_const Node
type Node_annotation Node
type Node_Which uint16

const (
	Node_Which_file       Node_Which = 0
	Node_Which_struct     Node_Which = 1
	Node_Which_enum       Node_Which = 2
	Node_Which_interface  Node_Which = 3
	Node_Which_const      Node_Which = 4
	Node_Which_annotation Node_Which = 5
)

func NewNode(s *C.Segment) Node { return Node(s.NewStruct(C.ObjectSize{DataSize: 40, PointerCount: 6})) }
func NewRootNode(s *C.Segment) Node {
	return Node(s.NewRootStruct(C.ObjectSize{DataSize: 40, PointerCount: 6}))
}
func AutoNewNode(s *C.Segment) Node {
	return Node(s.NewStructAR(C.ObjectSize{DataSize: 40, PointerCount: 6}))
}
func ReadRootNode(s *C.Segment) Node               { return Node(s.Root(0).ToStruct()) }
func (s Node) Which() Node_Which                   { return Node_Which(C.Struct(s).Get16(12)) }
func (s Node) Id() uint64                          { return C.Struct(s).Get64(0) }
func (s Node) SetId(v uint64)                      { C.Struct(s).Set64(0, v) }
func (s Node) DisplayName() string                 { return C.Struct(s).GetObject(0).ToText() }
func (s Node) SetDisplayName(v string)             { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s Node) DisplayNamePrefixLength() uint32     { return C.Struct(s).Get32(8) }
func (s Node) SetDisplayNamePrefixLength(v uint32) { C.Struct(s).Set32(8, v) }
func (s Node) ScopeId() uint64                     { return C.Struct(s).Get64(16) }
func (s Node) SetScopeId(v uint64)                 { C.Struct(s).Set64(16, v) }
func (s Node) Parameters() Node_Parameter_List     { return Node_Parameter_List(C.Struct(s).GetObject(5)) }
func (s Node) SetParameters(v Node_Parameter_List) { C.Struct(s).SetObject(5, C.Object(v)) }
func (s Node) IsGeneric() bool                     { return C.Struct(s).Get1(288) }
func (s Node) SetIsGeneric(v bool)                 { C.Struct(s).Set1(288, v) }
func (s Node) NestedNodes() Node_NestedNode_List {
	return Node_NestedNode_List(C.Struct(s).GetObject(1))
}
func (s Node) SetNestedNodes(v Node_NestedNode_List)         { C.Struct(s).SetObject(1, C.Object(v)) }
func (s Node) Annotations() Annotation_List                  { return Annotation_List(C.Struct(s).GetObject(2)) }
func (s Node) SetAnnotations(v Annotation_List)              { C.Struct(s).SetObject(2, C.Object(v)) }
func (s Node) SetFile()                                      { C.Struct(s).Set16(12, 0) }
func (s Node) Struct() Node_struct                           { return Node_struct(s) }
func (s Node) SetStruct()                                    { C.Struct(s).Set16(12, 1) }
func (s Node_struct) DataWordCount() uint16                  { return C.Struct(s).Get16(14) }
func (s Node_struct) SetDataWordCount(v uint16)              { C.Struct(s).Set16(14, v) }
func (s Node_struct) PointerCount() uint16                   { return C.Struct(s).Get16(24) }
func (s Node_struct) SetPointerCount(v uint16)               { C.Struct(s).Set16(24, v) }
func (s Node_struct) PreferredListEncoding() ElementSize     { return ElementSize(C.Struct(s).Get16(26)) }
func (s Node_struct) SetPreferredListEncoding(v ElementSize) { C.Struct(s).Set16(26, uint16(v)) }
func (s Node_struct) IsGroup() bool                          { return C.Struct(s).Get1(224) }
func (s Node_struct) SetIsGroup(v bool)                      { C.Struct(s).Set1(224, v) }
func (s Node_struct) DiscriminantCount() uint16              { return C.Struct(s).Get16(30) }
func (s Node_struct) SetDiscriminantCount(v uint16)          { C.Struct(s).Set16(30, v) }
func (s Node_struct) DiscriminantOffset() uint32             { return C.Struct(s).Get32(32) }
func (s Node_struct) SetDiscriminantOffset(v uint32)         { C.Struct(s).Set32(32, v) }
func (s Node_struct) Fields() Field_List                     { return Field_List(C.Struct(s).GetObject(3)) }
func (s Node_struct) SetFields(v Field_List)                 { C.Struct(s).SetObject(3, C.Object(v)) }
func (s Node) Enum() Node_enum                               { return Node_enum(s) }
func (s Node) SetEnum()                                      { C.Struct(s).Set16(12, 2) }
func (s Node_enum) Enumerants() Enumerant_List               { return Enumerant_List(C.Struct(s).GetObject(3)) }
func (s Node_enum) SetEnumerants(v Enumerant_List)           { C.Struct(s).SetObject(3, C.Object(v)) }
func (s Node) Interface() Node_interface                     { return Node_interface(s) }
func (s Node) SetInterface()                                 { C.Struct(s).Set16(12, 3) }
func (s Node_interface) Methods() Method_List                { return Method_List(C.Struct(s).GetObject(3)) }
func (s Node_interface) SetMethods(v Method_List)            { C.Struct(s).SetObject(3, C.Object(v)) }
func (s Node_interface) Superclasses() Superclass_List {
	return Superclass_List(C.Struct(s).GetObject(4))
}
func (s Node_interface) SetSuperclasses(v Superclass_List) { C.Struct(s).SetObject(4, C.Object(v)) }
func (s Node) Const() Node_const                           { return Node_const(s) }
func (s Node) SetConst()                                   { C.Struct(s).Set16(12, 4) }
func (s Node_const) Type() Type                            { return Type(C.Struct(s).GetObject(3).ToStruct()) }
func (s Node_const) SetType(v Type)                        { C.Struct(s).SetObject(3, C.Object(v)) }
func (s Node_const) Value() Value                          { return Value(C.Struct(s).GetObject(4).ToStruct()) }
func (s Node_const) SetValue(v Value)                      { C.Struct(s).SetObject(4, C.Object(v)) }
func (s Node) Annotation() Node_annotation                 { return Node_annotation(s) }
func (s Node) SetAnnotation()                              { C.Struct(s).Set16(12, 5) }
func (s Node_annotation) Type() Type                       { return Type(C.Struct(s).GetObject(3).ToStruct()) }
func (s Node_annotation) SetType(v Type)                   { C.Struct(s).SetObject(3, C.Object(v)) }
func (s Node_annotation) TargetsFile() bool                { return C.Struct(s).Get1(112) }
func (s Node_annotation) SetTargetsFile(v bool)            { C.Struct(s).Set1(112, v) }
func (s Node_annotation) TargetsConst() bool               { return C.Struct(s).Get1(113) }
func (s Node_annotation) SetTargetsConst(v bool)           { C.Struct(s).Set1(113, v) }
func (s Node_annotation) TargetsEnum() bool                { return C.Struct(s).Get1(114) }
func (s Node_annotation) SetTargetsEnum(v bool)            { C.Struct(s).Set1(114, v) }
func (s Node_annotation) TargetsEnumerant() bool           { return C.Struct(s).Get1(115) }
func (s Node_annotation) SetTargetsEnumerant(v bool)       { C.Struct(s).Set1(115, v) }
func (s Node_annotation) TargetsStruct() bool              { return C.Struct(s).Get1(116) }
func (s Node_annotation) SetTargetsStruct(v bool)          { C.Struct(s).Set1(116, v) }
func (s Node_annotation) TargetsField() bool               { return C.Struct(s).Get1(117) }
func (s Node_annotation) SetTargetsField(v bool)           { C.Struct(s).Set1(117, v) }
func (s Node_annotation) TargetsUnion() bool               { return C.Struct(s).Get1(118) }
func (s Node_annotation) SetTargetsUnion(v bool)           { C.Struct(s).Set1(118, v) }
func (s Node_annotation) TargetsGroup() bool               { return C.Struct(s).Get1(119) }
func (s Node_annotation) SetTargetsGroup(v bool)           { C.Struct(s).Set1(119, v) }
func (s Node_annotation) TargetsInterface() bool           { return C.Struct(s).Get1(120) }
func (s Node_annotation) SetTargetsInterface(v bool)       { C.Struct(s).Set1(120, v) }
func (s Node_annotation) TargetsMethod() bool              { return C.Struct(s).Get1(121) }
func (s Node_annotation) SetTargetsMethod(v bool)          { C.Struct(s).Set1(121, v) }
func (s Node_annotation) TargetsParam() bool               { return C.Struct(s).Get1(122) }
func (s Node_annotation) SetTargetsParam(v bool)           { C.Struct(s).Set1(122, v) }
func (s Node_annotation) TargetsAnnotation() bool          { return C.Struct(s).Get1(123) }
func (s Node_annotation) SetTargetsAnnotation(v bool)      { C.Struct(s).Set1(123, v) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Node) MarshalJSON() (bs []byte, err error) { return }

type Node_List C.PointerList

func NewNode_List(s *C.Segment, sz int) Node_List {
	return Node_List(s.NewCompositeList(C.ObjectSize{DataSize: 40, PointerCount: 6}, sz))
}
func (s Node_List) Len() int             { return C.PointerList(s).Len() }
func (s Node_List) At(i int) Node        { return Node(C.PointerList(s).At(i).ToStruct()) }
func (s Node_List) Set(i int, item Node) { C.PointerList(s).Set(i, C.Object(item)) }

type Node_Parameter C.Struct

func NewNode_Parameter(s *C.Segment) Node_Parameter {
	return Node_Parameter(s.NewStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func NewRootNode_Parameter(s *C.Segment) Node_Parameter {
	return Node_Parameter(s.NewRootStruct(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func AutoNewNode_Parameter(s *C.Segment) Node_Parameter {
	return Node_Parameter(s.NewStructAR(C.ObjectSize{DataSize: 0, PointerCount: 1}))
}
func ReadRootNode_Parameter(s *C.Segment) Node_Parameter { return Node_Parameter(s.Root(0).ToStruct()) }
func (s Node_Parameter) Name() string                    { return C.Struct(s).GetObject(0).ToText() }
func (s Node_Parameter) SetName(v string)                { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Node_Parameter) MarshalJSON() (bs []byte, err error) { return }

type Node_Parameter_List C.PointerList

func NewNode_Parameter_List(s *C.Segment, sz int) Node_Parameter_List {
	return Node_Parameter_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s Node_Parameter_List) Len() int { return C.PointerList(s).Len() }
func (s Node_Parameter_List) At(i int) Node_Parameter {
	return Node_Parameter(C.PointerList(s).At(i).ToStruct())
}
func (s Node_Parameter_List) Set(i int, item Node_Parameter) { C.PointerList(s).Set(i, C.Object(item)) }

type Node_NestedNode C.Struct

func NewNode_NestedNode(s *C.Segment) Node_NestedNode {
	return Node_NestedNode(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func NewRootNode_NestedNode(s *C.Segment) Node_NestedNode {
	return Node_NestedNode(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func AutoNewNode_NestedNode(s *C.Segment) Node_NestedNode {
	return Node_NestedNode(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func ReadRootNode_NestedNode(s *C.Segment) Node_NestedNode {
	return Node_NestedNode(s.Root(0).ToStruct())
}
func (s Node_NestedNode) Name() string     { return C.Struct(s).GetObject(0).ToText() }
func (s Node_NestedNode) SetName(v string) { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s Node_NestedNode) Id() uint64       { return C.Struct(s).Get64(0) }
func (s Node_NestedNode) SetId(v uint64)   { C.Struct(s).Set64(0, v) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Node_NestedNode) MarshalJSON() (bs []byte, err error) { return }

type Node_NestedNode_List C.PointerList

func NewNode_NestedNode_List(s *C.Segment, sz int) Node_NestedNode_List {
	return Node_NestedNode_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 1}, sz))
}
func (s Node_NestedNode_List) Len() int { return C.PointerList(s).Len() }
func (s Node_NestedNode_List) At(i int) Node_NestedNode {
	return Node_NestedNode(C.PointerList(s).At(i).ToStruct())
}
func (s Node_NestedNode_List) Set(i int, item Node_NestedNode) {
	C.PointerList(s).Set(i, C.Object(item))
}

type Field C.Struct
type Field_slot Field
type Field_group Field
type Field_ordinal Field
type Field_Which uint16

const (
	Field_Which_slot  Field_Which = 0
	Field_Which_group Field_Which = 1
)

type Field_ordinal_Which uint16

const (
	Field_ordinal_Which_implicit Field_ordinal_Which = 0
	Field_ordinal_Which_explicit Field_ordinal_Which = 1
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
func ReadRootField(s *C.Segment) Field             { return Field(s.Root(0).ToStruct()) }
func (s Field) Which() Field_Which                 { return Field_Which(C.Struct(s).Get16(8)) }
func (s Field) Name() string                       { return C.Struct(s).GetObject(0).ToText() }
func (s Field) SetName(v string)                   { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s Field) CodeOrder() uint16                  { return C.Struct(s).Get16(0) }
func (s Field) SetCodeOrder(v uint16)              { C.Struct(s).Set16(0, v) }
func (s Field) Annotations() Annotation_List       { return Annotation_List(C.Struct(s).GetObject(1)) }
func (s Field) SetAnnotations(v Annotation_List)   { C.Struct(s).SetObject(1, C.Object(v)) }
func (s Field) DiscriminantValue() uint16          { return C.Struct(s).Get16(2) ^ 65535 }
func (s Field) SetDiscriminantValue(v uint16)      { C.Struct(s).Set16(2, v^65535) }
func (s Field) Slot() Field_slot                   { return Field_slot(s) }
func (s Field) SetSlot()                           { C.Struct(s).Set16(8, 0) }
func (s Field_slot) Offset() uint32                { return C.Struct(s).Get32(4) }
func (s Field_slot) SetOffset(v uint32)            { C.Struct(s).Set32(4, v) }
func (s Field_slot) Type() Type                    { return Type(C.Struct(s).GetObject(2).ToStruct()) }
func (s Field_slot) SetType(v Type)                { C.Struct(s).SetObject(2, C.Object(v)) }
func (s Field_slot) DefaultValue() Value           { return Value(C.Struct(s).GetObject(3).ToStruct()) }
func (s Field_slot) SetDefaultValue(v Value)       { C.Struct(s).SetObject(3, C.Object(v)) }
func (s Field_slot) HadExplicitDefault() bool      { return C.Struct(s).Get1(128) }
func (s Field_slot) SetHadExplicitDefault(v bool)  { C.Struct(s).Set1(128, v) }
func (s Field) Group() Field_group                 { return Field_group(s) }
func (s Field) SetGroup()                          { C.Struct(s).Set16(8, 1) }
func (s Field_group) TypeId() uint64               { return C.Struct(s).Get64(16) }
func (s Field_group) SetTypeId(v uint64)           { C.Struct(s).Set64(16, v) }
func (s Field) Ordinal() Field_ordinal             { return Field_ordinal(s) }
func (s Field_ordinal) Which() Field_ordinal_Which { return Field_ordinal_Which(C.Struct(s).Get16(10)) }
func (s Field_ordinal) SetImplicit()               { C.Struct(s).Set16(10, 0) }
func (s Field_ordinal) Explicit() uint16           { return C.Struct(s).Get16(12) }
func (s Field_ordinal) SetExplicit(v uint16)       { C.Struct(s).Set16(10, 1); C.Struct(s).Set16(12, v) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Field) MarshalJSON() (bs []byte, err error) { return }

type Field_List C.PointerList

func NewField_List(s *C.Segment, sz int) Field_List {
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

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Enumerant) MarshalJSON() (bs []byte, err error) { return }

type Enumerant_List C.PointerList

func NewEnumerant_List(s *C.Segment, sz int) Enumerant_List {
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

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Superclass) MarshalJSON() (bs []byte, err error) { return }

type Superclass_List C.PointerList

func NewSuperclass_List(s *C.Segment, sz int) Superclass_List {
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
func (s Method) ImplicitParameters() Node_Parameter_List {
	return Node_Parameter_List(C.Struct(s).GetObject(4))
}
func (s Method) SetImplicitParameters(v Node_Parameter_List) { C.Struct(s).SetObject(4, C.Object(v)) }
func (s Method) ParamStructType() uint64                     { return C.Struct(s).Get64(8) }
func (s Method) SetParamStructType(v uint64)                 { C.Struct(s).Set64(8, v) }
func (s Method) ParamBrand() Brand                           { return Brand(C.Struct(s).GetObject(2).ToStruct()) }
func (s Method) SetParamBrand(v Brand)                       { C.Struct(s).SetObject(2, C.Object(v)) }
func (s Method) ResultStructType() uint64                    { return C.Struct(s).Get64(16) }
func (s Method) SetResultStructType(v uint64)                { C.Struct(s).Set64(16, v) }
func (s Method) ResultBrand() Brand                          { return Brand(C.Struct(s).GetObject(3).ToStruct()) }
func (s Method) SetResultBrand(v Brand)                      { C.Struct(s).SetObject(3, C.Object(v)) }
func (s Method) Annotations() Annotation_List                { return Annotation_List(C.Struct(s).GetObject(1)) }
func (s Method) SetAnnotations(v Annotation_List)            { C.Struct(s).SetObject(1, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Method) MarshalJSON() (bs []byte, err error) { return }

type Method_List C.PointerList

func NewMethod_List(s *C.Segment, sz int) Method_List {
	return Method_List(s.NewCompositeList(C.ObjectSize{DataSize: 24, PointerCount: 5}, sz))
}
func (s Method_List) Len() int               { return C.PointerList(s).Len() }
func (s Method_List) At(i int) Method        { return Method(C.PointerList(s).At(i).ToStruct()) }
func (s Method_List) Set(i int, item Method) { C.PointerList(s).Set(i, C.Object(item)) }

type Type C.Struct
type Type_list Type
type Type_enum Type
type Type_struct Type
type Type_interface Type
type Type_anyPointer Type
type Type_anyPointer_parameter Type
type Type_anyPointer_implicitMethodParameter Type
type Type_Which uint16

const (
	Type_Which_void       Type_Which = 0
	Type_Which_bool       Type_Which = 1
	Type_Which_int8       Type_Which = 2
	Type_Which_int16      Type_Which = 3
	Type_Which_int32      Type_Which = 4
	Type_Which_int64      Type_Which = 5
	Type_Which_uint8      Type_Which = 6
	Type_Which_uint16     Type_Which = 7
	Type_Which_uint32     Type_Which = 8
	Type_Which_uint64     Type_Which = 9
	Type_Which_float32    Type_Which = 10
	Type_Which_float64    Type_Which = 11
	Type_Which_text       Type_Which = 12
	Type_Which_data       Type_Which = 13
	Type_Which_list       Type_Which = 14
	Type_Which_enum       Type_Which = 15
	Type_Which_struct     Type_Which = 16
	Type_Which_interface  Type_Which = 17
	Type_Which_anyPointer Type_Which = 18
)

type Type_anyPointer_Which uint16

const (
	Type_anyPointer_Which_unconstrained           Type_anyPointer_Which = 0
	Type_anyPointer_Which_parameter               Type_anyPointer_Which = 1
	Type_anyPointer_Which_implicitMethodParameter Type_anyPointer_Which = 2
)

func NewType(s *C.Segment) Type { return Type(s.NewStruct(C.ObjectSize{DataSize: 24, PointerCount: 1})) }
func NewRootType(s *C.Segment) Type {
	return Type(s.NewRootStruct(C.ObjectSize{DataSize: 24, PointerCount: 1}))
}
func AutoNewType(s *C.Segment) Type {
	return Type(s.NewStructAR(C.ObjectSize{DataSize: 24, PointerCount: 1}))
}
func ReadRootType(s *C.Segment) Type        { return Type(s.Root(0).ToStruct()) }
func (s Type) Which() Type_Which            { return Type_Which(C.Struct(s).Get16(0)) }
func (s Type) SetVoid()                     { C.Struct(s).Set16(0, 0) }
func (s Type) SetBool()                     { C.Struct(s).Set16(0, 1) }
func (s Type) SetInt8()                     { C.Struct(s).Set16(0, 2) }
func (s Type) SetInt16()                    { C.Struct(s).Set16(0, 3) }
func (s Type) SetInt32()                    { C.Struct(s).Set16(0, 4) }
func (s Type) SetInt64()                    { C.Struct(s).Set16(0, 5) }
func (s Type) SetUint8()                    { C.Struct(s).Set16(0, 6) }
func (s Type) SetUint16()                   { C.Struct(s).Set16(0, 7) }
func (s Type) SetUint32()                   { C.Struct(s).Set16(0, 8) }
func (s Type) SetUint64()                   { C.Struct(s).Set16(0, 9) }
func (s Type) SetFloat32()                  { C.Struct(s).Set16(0, 10) }
func (s Type) SetFloat64()                  { C.Struct(s).Set16(0, 11) }
func (s Type) SetText()                     { C.Struct(s).Set16(0, 12) }
func (s Type) SetData()                     { C.Struct(s).Set16(0, 13) }
func (s Type) List() Type_list              { return Type_list(s) }
func (s Type) SetList()                     { C.Struct(s).Set16(0, 14) }
func (s Type_list) ElementType() Type       { return Type(C.Struct(s).GetObject(0).ToStruct()) }
func (s Type_list) SetElementType(v Type)   { C.Struct(s).SetObject(0, C.Object(v)) }
func (s Type) Enum() Type_enum              { return Type_enum(s) }
func (s Type) SetEnum()                     { C.Struct(s).Set16(0, 15) }
func (s Type_enum) TypeId() uint64          { return C.Struct(s).Get64(8) }
func (s Type_enum) SetTypeId(v uint64)      { C.Struct(s).Set64(8, v) }
func (s Type_enum) Brand() Brand            { return Brand(C.Struct(s).GetObject(0).ToStruct()) }
func (s Type_enum) SetBrand(v Brand)        { C.Struct(s).SetObject(0, C.Object(v)) }
func (s Type) Struct() Type_struct          { return Type_struct(s) }
func (s Type) SetStruct()                   { C.Struct(s).Set16(0, 16) }
func (s Type_struct) TypeId() uint64        { return C.Struct(s).Get64(8) }
func (s Type_struct) SetTypeId(v uint64)    { C.Struct(s).Set64(8, v) }
func (s Type_struct) Brand() Brand          { return Brand(C.Struct(s).GetObject(0).ToStruct()) }
func (s Type_struct) SetBrand(v Brand)      { C.Struct(s).SetObject(0, C.Object(v)) }
func (s Type) Interface() Type_interface    { return Type_interface(s) }
func (s Type) SetInterface()                { C.Struct(s).Set16(0, 17) }
func (s Type_interface) TypeId() uint64     { return C.Struct(s).Get64(8) }
func (s Type_interface) SetTypeId(v uint64) { C.Struct(s).Set64(8, v) }
func (s Type_interface) Brand() Brand       { return Brand(C.Struct(s).GetObject(0).ToStruct()) }
func (s Type_interface) SetBrand(v Brand)   { C.Struct(s).SetObject(0, C.Object(v)) }
func (s Type) AnyPointer() Type_anyPointer  { return Type_anyPointer(s) }
func (s Type) SetAnyPointer()               { C.Struct(s).Set16(0, 18) }
func (s Type_anyPointer) Which() Type_anyPointer_Which {
	return Type_anyPointer_Which(C.Struct(s).Get16(8))
}
func (s Type_anyPointer) SetUnconstrained()                    { C.Struct(s).Set16(8, 0) }
func (s Type_anyPointer) Parameter() Type_anyPointer_parameter { return Type_anyPointer_parameter(s) }
func (s Type_anyPointer) SetParameter()                        { C.Struct(s).Set16(8, 1) }
func (s Type_anyPointer_parameter) ScopeId() uint64            { return C.Struct(s).Get64(16) }
func (s Type_anyPointer_parameter) SetScopeId(v uint64)        { C.Struct(s).Set64(16, v) }
func (s Type_anyPointer_parameter) ParameterIndex() uint16     { return C.Struct(s).Get16(10) }
func (s Type_anyPointer_parameter) SetParameterIndex(v uint16) { C.Struct(s).Set16(10, v) }
func (s Type_anyPointer) ImplicitMethodParameter() Type_anyPointer_implicitMethodParameter {
	return Type_anyPointer_implicitMethodParameter(s)
}
func (s Type_anyPointer) SetImplicitMethodParameter()                        { C.Struct(s).Set16(8, 2) }
func (s Type_anyPointer_implicitMethodParameter) ParameterIndex() uint16     { return C.Struct(s).Get16(10) }
func (s Type_anyPointer_implicitMethodParameter) SetParameterIndex(v uint16) { C.Struct(s).Set16(10, v) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Type) MarshalJSON() (bs []byte, err error) { return }

type Type_List C.PointerList

func NewType_List(s *C.Segment, sz int) Type_List {
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
func ReadRootBrand(s *C.Segment) Brand       { return Brand(s.Root(0).ToStruct()) }
func (s Brand) Scopes() Brand_Scope_List     { return Brand_Scope_List(C.Struct(s).GetObject(0)) }
func (s Brand) SetScopes(v Brand_Scope_List) { C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Brand) MarshalJSON() (bs []byte, err error) { return }

type Brand_List C.PointerList

func NewBrand_List(s *C.Segment, sz int) Brand_List {
	return Brand_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s Brand_List) Len() int              { return C.PointerList(s).Len() }
func (s Brand_List) At(i int) Brand        { return Brand(C.PointerList(s).At(i).ToStruct()) }
func (s Brand_List) Set(i int, item Brand) { C.PointerList(s).Set(i, C.Object(item)) }

type Brand_Scope C.Struct
type Brand_Scope_Which uint16

const (
	Brand_Scope_Which_bind    Brand_Scope_Which = 0
	Brand_Scope_Which_inherit Brand_Scope_Which = 1
)

func NewBrand_Scope(s *C.Segment) Brand_Scope {
	return Brand_Scope(s.NewStruct(C.ObjectSize{DataSize: 16, PointerCount: 1}))
}
func NewRootBrand_Scope(s *C.Segment) Brand_Scope {
	return Brand_Scope(s.NewRootStruct(C.ObjectSize{DataSize: 16, PointerCount: 1}))
}
func AutoNewBrand_Scope(s *C.Segment) Brand_Scope {
	return Brand_Scope(s.NewStructAR(C.ObjectSize{DataSize: 16, PointerCount: 1}))
}
func ReadRootBrand_Scope(s *C.Segment) Brand_Scope { return Brand_Scope(s.Root(0).ToStruct()) }
func (s Brand_Scope) Which() Brand_Scope_Which     { return Brand_Scope_Which(C.Struct(s).Get16(8)) }
func (s Brand_Scope) ScopeId() uint64              { return C.Struct(s).Get64(0) }
func (s Brand_Scope) SetScopeId(v uint64)          { C.Struct(s).Set64(0, v) }
func (s Brand_Scope) Bind() Brand_Binding_List     { return Brand_Binding_List(C.Struct(s).GetObject(0)) }
func (s Brand_Scope) SetBind(v Brand_Binding_List) {
	C.Struct(s).Set16(8, 0)
	C.Struct(s).SetObject(0, C.Object(v))
}
func (s Brand_Scope) SetInherit() { C.Struct(s).Set16(8, 1) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Brand_Scope) MarshalJSON() (bs []byte, err error) { return }

type Brand_Scope_List C.PointerList

func NewBrand_Scope_List(s *C.Segment, sz int) Brand_Scope_List {
	return Brand_Scope_List(s.NewCompositeList(C.ObjectSize{DataSize: 16, PointerCount: 1}, sz))
}
func (s Brand_Scope_List) Len() int                    { return C.PointerList(s).Len() }
func (s Brand_Scope_List) At(i int) Brand_Scope        { return Brand_Scope(C.PointerList(s).At(i).ToStruct()) }
func (s Brand_Scope_List) Set(i int, item Brand_Scope) { C.PointerList(s).Set(i, C.Object(item)) }

type Brand_Binding C.Struct
type Brand_Binding_Which uint16

const (
	Brand_Binding_Which_unbound Brand_Binding_Which = 0
	Brand_Binding_Which_type    Brand_Binding_Which = 1
)

func NewBrand_Binding(s *C.Segment) Brand_Binding {
	return Brand_Binding(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func NewRootBrand_Binding(s *C.Segment) Brand_Binding {
	return Brand_Binding(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func AutoNewBrand_Binding(s *C.Segment) Brand_Binding {
	return Brand_Binding(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func ReadRootBrand_Binding(s *C.Segment) Brand_Binding { return Brand_Binding(s.Root(0).ToStruct()) }
func (s Brand_Binding) Which() Brand_Binding_Which     { return Brand_Binding_Which(C.Struct(s).Get16(0)) }
func (s Brand_Binding) SetUnbound()                    { C.Struct(s).Set16(0, 0) }
func (s Brand_Binding) Type() Type                     { return Type(C.Struct(s).GetObject(0).ToStruct()) }
func (s Brand_Binding) SetType(v Type)                 { C.Struct(s).Set16(0, 1); C.Struct(s).SetObject(0, C.Object(v)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Brand_Binding) MarshalJSON() (bs []byte, err error) { return }

type Brand_Binding_List C.PointerList

func NewBrand_Binding_List(s *C.Segment, sz int) Brand_Binding_List {
	return Brand_Binding_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 1}, sz))
}
func (s Brand_Binding_List) Len() int { return C.PointerList(s).Len() }
func (s Brand_Binding_List) At(i int) Brand_Binding {
	return Brand_Binding(C.PointerList(s).At(i).ToStruct())
}
func (s Brand_Binding_List) Set(i int, item Brand_Binding) { C.PointerList(s).Set(i, C.Object(item)) }

type Value C.Struct
type Value_Which uint16

const (
	Value_Which_void       Value_Which = 0
	Value_Which_bool       Value_Which = 1
	Value_Which_int8       Value_Which = 2
	Value_Which_int16      Value_Which = 3
	Value_Which_int32      Value_Which = 4
	Value_Which_int64      Value_Which = 5
	Value_Which_uint8      Value_Which = 6
	Value_Which_uint16     Value_Which = 7
	Value_Which_uint32     Value_Which = 8
	Value_Which_uint64     Value_Which = 9
	Value_Which_float32    Value_Which = 10
	Value_Which_float64    Value_Which = 11
	Value_Which_text       Value_Which = 12
	Value_Which_data       Value_Which = 13
	Value_Which_list       Value_Which = 14
	Value_Which_enum       Value_Which = 15
	Value_Which_struct     Value_Which = 16
	Value_Which_interface  Value_Which = 17
	Value_Which_anyPointer Value_Which = 18
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

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Value) MarshalJSON() (bs []byte, err error) { return }

type Value_List C.PointerList

func NewValue_List(s *C.Segment, sz int) Value_List {
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

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s Annotation) MarshalJSON() (bs []byte, err error) { return }

type Annotation_List C.PointerList

func NewAnnotation_List(s *C.Segment, sz int) Annotation_List {
	return Annotation_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 2}, sz))
}
func (s Annotation_List) Len() int                   { return C.PointerList(s).Len() }
func (s Annotation_List) At(i int) Annotation        { return Annotation(C.PointerList(s).At(i).ToStruct()) }
func (s Annotation_List) Set(i int, item Annotation) { C.PointerList(s).Set(i, C.Object(item)) }

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

func (c ElementSize) String() string {
	switch c {
	case ElementSize_empty:
		return "empty"
	case ElementSize_bit:
		return "bit"
	case ElementSize_byte:
		return "byte"
	case ElementSize_twoBytes:
		return "twoBytes"
	case ElementSize_fourBytes:
		return "fourBytes"
	case ElementSize_eightBytes:
		return "eightBytes"
	case ElementSize_pointer:
		return "pointer"
	case ElementSize_inlineComposite:
		return "inlineComposite"
	default:
		return ""
	}
}

func ElementSizeFromString(c string) ElementSize {
	switch c {
	case "empty":
		return ElementSize_empty
	case "bit":
		return ElementSize_bit
	case "byte":
		return ElementSize_byte
	case "twoBytes":
		return ElementSize_twoBytes
	case "fourBytes":
		return ElementSize_fourBytes
	case "eightBytes":
		return ElementSize_eightBytes
	case "pointer":
		return ElementSize_pointer
	case "inlineComposite":
		return ElementSize_inlineComposite
	default:
		return 0
	}
}

type ElementSize_List C.PointerList

func NewElementSize_List(s *C.Segment, sz int) ElementSize_List {
	return ElementSize_List(s.NewUInt16List(sz))
}
func (s ElementSize_List) Len() int             { return C.UInt16List(s).Len() }
func (s ElementSize_List) At(i int) ElementSize { return ElementSize(C.UInt16List(s).At(i)) }

// capn.JSON_enabled == false so we stub MarshalJSON().
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
func (s CodeGeneratorRequest) RequestedFiles() CodeGeneratorRequest_RequestedFile_List {
	return CodeGeneratorRequest_RequestedFile_List(C.Struct(s).GetObject(1))
}
func (s CodeGeneratorRequest) SetRequestedFiles(v CodeGeneratorRequest_RequestedFile_List) {
	C.Struct(s).SetObject(1, C.Object(v))
}

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s CodeGeneratorRequest) MarshalJSON() (bs []byte, err error) { return }

type CodeGeneratorRequest_List C.PointerList

func NewCodeGeneratorRequest_List(s *C.Segment, sz int) CodeGeneratorRequest_List {
	return CodeGeneratorRequest_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 2}, sz))
}
func (s CodeGeneratorRequest_List) Len() int { return C.PointerList(s).Len() }
func (s CodeGeneratorRequest_List) At(i int) CodeGeneratorRequest {
	return CodeGeneratorRequest(C.PointerList(s).At(i).ToStruct())
}
func (s CodeGeneratorRequest_List) Set(i int, item CodeGeneratorRequest) {
	C.PointerList(s).Set(i, C.Object(item))
}

type CodeGeneratorRequest_RequestedFile C.Struct

func NewCodeGeneratorRequest_RequestedFile(s *C.Segment) CodeGeneratorRequest_RequestedFile {
	return CodeGeneratorRequest_RequestedFile(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 2}))
}
func NewRootCodeGeneratorRequest_RequestedFile(s *C.Segment) CodeGeneratorRequest_RequestedFile {
	return CodeGeneratorRequest_RequestedFile(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 2}))
}
func AutoNewCodeGeneratorRequest_RequestedFile(s *C.Segment) CodeGeneratorRequest_RequestedFile {
	return CodeGeneratorRequest_RequestedFile(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 2}))
}
func ReadRootCodeGeneratorRequest_RequestedFile(s *C.Segment) CodeGeneratorRequest_RequestedFile {
	return CodeGeneratorRequest_RequestedFile(s.Root(0).ToStruct())
}
func (s CodeGeneratorRequest_RequestedFile) Id() uint64     { return C.Struct(s).Get64(0) }
func (s CodeGeneratorRequest_RequestedFile) SetId(v uint64) { C.Struct(s).Set64(0, v) }
func (s CodeGeneratorRequest_RequestedFile) Filename() string {
	return C.Struct(s).GetObject(0).ToText()
}
func (s CodeGeneratorRequest_RequestedFile) SetFilename(v string) {
	C.Struct(s).SetObject(0, s.Segment.NewText(v))
}
func (s CodeGeneratorRequest_RequestedFile) Imports() CodeGeneratorRequest_RequestedFile_Import_List {
	return CodeGeneratorRequest_RequestedFile_Import_List(C.Struct(s).GetObject(1))
}
func (s CodeGeneratorRequest_RequestedFile) SetImports(v CodeGeneratorRequest_RequestedFile_Import_List) {
	C.Struct(s).SetObject(1, C.Object(v))
}

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s CodeGeneratorRequest_RequestedFile) MarshalJSON() (bs []byte, err error) { return }

type CodeGeneratorRequest_RequestedFile_List C.PointerList

func NewCodeGeneratorRequest_RequestedFile_List(s *C.Segment, sz int) CodeGeneratorRequest_RequestedFile_List {
	return CodeGeneratorRequest_RequestedFile_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 2}, sz))
}
func (s CodeGeneratorRequest_RequestedFile_List) Len() int { return C.PointerList(s).Len() }
func (s CodeGeneratorRequest_RequestedFile_List) At(i int) CodeGeneratorRequest_RequestedFile {
	return CodeGeneratorRequest_RequestedFile(C.PointerList(s).At(i).ToStruct())
}
func (s CodeGeneratorRequest_RequestedFile_List) Set(i int, item CodeGeneratorRequest_RequestedFile) {
	C.PointerList(s).Set(i, C.Object(item))
}

type CodeGeneratorRequest_RequestedFile_Import C.Struct

func NewCodeGeneratorRequest_RequestedFile_Import(s *C.Segment) CodeGeneratorRequest_RequestedFile_Import {
	return CodeGeneratorRequest_RequestedFile_Import(s.NewStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func NewRootCodeGeneratorRequest_RequestedFile_Import(s *C.Segment) CodeGeneratorRequest_RequestedFile_Import {
	return CodeGeneratorRequest_RequestedFile_Import(s.NewRootStruct(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func AutoNewCodeGeneratorRequest_RequestedFile_Import(s *C.Segment) CodeGeneratorRequest_RequestedFile_Import {
	return CodeGeneratorRequest_RequestedFile_Import(s.NewStructAR(C.ObjectSize{DataSize: 8, PointerCount: 1}))
}
func ReadRootCodeGeneratorRequest_RequestedFile_Import(s *C.Segment) CodeGeneratorRequest_RequestedFile_Import {
	return CodeGeneratorRequest_RequestedFile_Import(s.Root(0).ToStruct())
}
func (s CodeGeneratorRequest_RequestedFile_Import) Id() uint64     { return C.Struct(s).Get64(0) }
func (s CodeGeneratorRequest_RequestedFile_Import) SetId(v uint64) { C.Struct(s).Set64(0, v) }
func (s CodeGeneratorRequest_RequestedFile_Import) Name() string {
	return C.Struct(s).GetObject(0).ToText()
}
func (s CodeGeneratorRequest_RequestedFile_Import) SetName(v string) {
	C.Struct(s).SetObject(0, s.Segment.NewText(v))
}

// capn.JSON_enabled == false so we stub MarshalJSON().
func (s CodeGeneratorRequest_RequestedFile_Import) MarshalJSON() (bs []byte, err error) { return }

type CodeGeneratorRequest_RequestedFile_Import_List C.PointerList

func NewCodeGeneratorRequest_RequestedFile_Import_List(s *C.Segment, sz int) CodeGeneratorRequest_RequestedFile_Import_List {
	return CodeGeneratorRequest_RequestedFile_Import_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 1}, sz))
}
func (s CodeGeneratorRequest_RequestedFile_Import_List) Len() int { return C.PointerList(s).Len() }
func (s CodeGeneratorRequest_RequestedFile_Import_List) At(i int) CodeGeneratorRequest_RequestedFile_Import {
	return CodeGeneratorRequest_RequestedFile_Import(C.PointerList(s).At(i).ToStruct())
}
func (s CodeGeneratorRequest_RequestedFile_Import_List) Set(i int, item CodeGeneratorRequest_RequestedFile_Import) {
	C.PointerList(s).Set(i, C.Object(item))
}
