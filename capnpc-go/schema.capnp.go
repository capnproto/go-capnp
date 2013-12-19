package main

// AUTO GENERATED - DO NOT EDIT

import (
	"math"
	"unsafe"

	C "github.com/glycerine/go-capnproto"
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
	NODE_STRUCT                = 1
	NODE_ENUM                  = 2
	NODE_INTERFACE             = 3
	NODE_CONST                 = 4
	NODE_ANNOTATION            = 5
)

func NewNode(s *C.Segment) Node                             { return Node(s.NewStruct(40, 5)) }
func NewRootNode(s *C.Segment) Node                         { return Node(s.NewRootStruct(40, 5)) }
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
func (s Node) NestedNodes() NodeNestedNode_List             { return NodeNestedNode_List(C.Struct(s).GetObject(1)) }
func (s Node) SetNestedNodes(v NodeNestedNode_List)         { C.Struct(s).SetObject(1, C.Object(v)) }
func (s Node) Annotations() Annotation_List                 { return Annotation_List(C.Struct(s).GetObject(2)) }
func (s Node) SetAnnotations(v Annotation_List)             { C.Struct(s).SetObject(2, C.Object(v)) }
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
func (s Node) Const() NodeConst                             { return NodeConst(s) }
func (s Node) SetConst()                                    { C.Struct(s).Set16(12, 4) }
func (s NodeConst) Type() Type                              { return Type(C.Struct(s).GetObject(3).ToStruct()) }
func (s NodeConst) SetType(v Type)                          { C.Struct(s).SetObject(3, C.Object(v)) }
func (s NodeConst) Value() Value                            { return Value(C.Struct(s).GetObject(4).ToStruct()) }
func (s NodeConst) SetValue(v Value)                        { C.Struct(s).SetObject(4, C.Object(v)) }
func (s Node) Annotation() NodeAnnotation                   { return NodeAnnotation(s) }
func (s Node) SetAnnotation()                               { C.Struct(s).Set16(12, 5) }
func (s NodeAnnotation) Type() Type                         { return Type(C.Struct(s).GetObject(3).ToStruct()) }
func (s NodeAnnotation) SetType(v Type)                     { C.Struct(s).SetObject(3, C.Object(v)) }
func (s NodeAnnotation) TargetsFile() bool                  { return C.Struct(s).Get1(112) }
func (s NodeAnnotation) SetTargetsFile(v bool)              { C.Struct(s).Set1(112, v) }
func (s NodeAnnotation) TargetsConst() bool                 { return C.Struct(s).Get1(113) }
func (s NodeAnnotation) SetTargetsConst(v bool)             { C.Struct(s).Set1(113, v) }
func (s NodeAnnotation) TargetsEnum() bool                  { return C.Struct(s).Get1(114) }
func (s NodeAnnotation) SetTargetsEnum(v bool)              { C.Struct(s).Set1(114, v) }
func (s NodeAnnotation) TargetsEnumerant() bool             { return C.Struct(s).Get1(115) }
func (s NodeAnnotation) SetTargetsEnumerant(v bool)         { C.Struct(s).Set1(115, v) }
func (s NodeAnnotation) TargetsStruct() bool                { return C.Struct(s).Get1(116) }
func (s NodeAnnotation) SetTargetsStruct(v bool)            { C.Struct(s).Set1(116, v) }
func (s NodeAnnotation) TargetsField() bool                 { return C.Struct(s).Get1(117) }
func (s NodeAnnotation) SetTargetsField(v bool)             { C.Struct(s).Set1(117, v) }
func (s NodeAnnotation) TargetsUnion() bool                 { return C.Struct(s).Get1(118) }
func (s NodeAnnotation) SetTargetsUnion(v bool)             { C.Struct(s).Set1(118, v) }
func (s NodeAnnotation) TargetsGroup() bool                 { return C.Struct(s).Get1(119) }
func (s NodeAnnotation) SetTargetsGroup(v bool)             { C.Struct(s).Set1(119, v) }
func (s NodeAnnotation) TargetsInterface() bool             { return C.Struct(s).Get1(120) }
func (s NodeAnnotation) SetTargetsInterface(v bool)         { C.Struct(s).Set1(120, v) }
func (s NodeAnnotation) TargetsMethod() bool                { return C.Struct(s).Get1(121) }
func (s NodeAnnotation) SetTargetsMethod(v bool)            { C.Struct(s).Set1(121, v) }
func (s NodeAnnotation) TargetsParam() bool                 { return C.Struct(s).Get1(122) }
func (s NodeAnnotation) SetTargetsParam(v bool)             { C.Struct(s).Set1(122, v) }
func (s NodeAnnotation) TargetsAnnotation() bool            { return C.Struct(s).Get1(123) }
func (s NodeAnnotation) SetTargetsAnnotation(v bool)        { C.Struct(s).Set1(123, v) }

type Node_List C.PointerList

func NewNodeList(s *C.Segment, sz int) Node_List { return Node_List(s.NewCompositeList(40, 5, sz)) }
func (s Node_List) Len() int                     { return C.PointerList(s).Len() }
func (s Node_List) At(i int) Node                { return Node(C.PointerList(s).At(i).ToStruct()) }
func (s Node_List) ToArray() []Node              { return *(*[]Node)(unsafe.Pointer(C.PointerList(s).ToArray())) }

type NodeNestedNode C.Struct

func NewNodeNestedNode(s *C.Segment) NodeNestedNode      { return NodeNestedNode(s.NewStruct(8, 1)) }
func NewRootNodeNestedNode(s *C.Segment) NodeNestedNode  { return NodeNestedNode(s.NewRootStruct(8, 1)) }
func ReadRootNodeNestedNode(s *C.Segment) NodeNestedNode { return NodeNestedNode(s.Root(0).ToStruct()) }
func (s NodeNestedNode) Name() string                    { return C.Struct(s).GetObject(0).ToText() }
func (s NodeNestedNode) SetName(v string)                { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s NodeNestedNode) Id() uint64                      { return C.Struct(s).Get64(0) }
func (s NodeNestedNode) SetId(v uint64)                  { C.Struct(s).Set64(0, v) }

type NodeNestedNode_List C.PointerList

func NewNodeNestedNodeList(s *C.Segment, sz int) NodeNestedNode_List {
	return NodeNestedNode_List(s.NewCompositeList(8, 1, sz))
}
func (s NodeNestedNode_List) Len() int { return C.PointerList(s).Len() }
func (s NodeNestedNode_List) At(i int) NodeNestedNode {
	return NodeNestedNode(C.PointerList(s).At(i).ToStruct())
}
func (s NodeNestedNode_List) ToArray() []NodeNestedNode {
	return *(*[]NodeNestedNode)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type Field C.Struct
type FieldSlot Field
type FieldGroup Field
type FieldOrdinal Field
type Field_Which uint16

const (
	FIELD_SLOT  Field_Which = 0
	FIELD_GROUP             = 1
)

type FieldOrdinal_Which uint16

const (
	FIELDORDINAL_IMPLICIT FieldOrdinal_Which = 0
	FIELDORDINAL_EXPLICIT                    = 1
)

func NewField(s *C.Segment) Field                { return Field(s.NewStruct(24, 4)) }
func NewRootField(s *C.Segment) Field            { return Field(s.NewRootStruct(24, 4)) }
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
func (s Field) Group() FieldGroup                { return FieldGroup(s) }
func (s Field) SetGroup()                        { C.Struct(s).Set16(8, 1) }
func (s FieldGroup) TypeId() uint64              { return C.Struct(s).Get64(16) }
func (s FieldGroup) SetTypeId(v uint64)          { C.Struct(s).Set64(16, v) }
func (s Field) Ordinal() FieldOrdinal            { return FieldOrdinal(s) }
func (s FieldOrdinal) Which() FieldOrdinal_Which { return FieldOrdinal_Which(C.Struct(s).Get16(10)) }
func (s FieldOrdinal) Explicit() uint16          { return C.Struct(s).Get16(12) }
func (s FieldOrdinal) SetExplicit(v uint16)      { C.Struct(s).Set16(10, 1); C.Struct(s).Set16(12, v) }

type Field_List C.PointerList

func NewFieldList(s *C.Segment, sz int) Field_List { return Field_List(s.NewCompositeList(24, 4, sz)) }
func (s Field_List) Len() int                      { return C.PointerList(s).Len() }
func (s Field_List) At(i int) Field                { return Field(C.PointerList(s).At(i).ToStruct()) }
func (s Field_List) ToArray() []Field              { return *(*[]Field)(unsafe.Pointer(C.PointerList(s).ToArray())) }

type Enumerant C.Struct

func NewEnumerant(s *C.Segment) Enumerant            { return Enumerant(s.NewStruct(8, 2)) }
func NewRootEnumerant(s *C.Segment) Enumerant        { return Enumerant(s.NewRootStruct(8, 2)) }
func ReadRootEnumerant(s *C.Segment) Enumerant       { return Enumerant(s.Root(0).ToStruct()) }
func (s Enumerant) Name() string                     { return C.Struct(s).GetObject(0).ToText() }
func (s Enumerant) SetName(v string)                 { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s Enumerant) CodeOrder() uint16                { return C.Struct(s).Get16(0) }
func (s Enumerant) SetCodeOrder(v uint16)            { C.Struct(s).Set16(0, v) }
func (s Enumerant) Annotations() Annotation_List     { return Annotation_List(C.Struct(s).GetObject(1)) }
func (s Enumerant) SetAnnotations(v Annotation_List) { C.Struct(s).SetObject(1, C.Object(v)) }

type Enumerant_List C.PointerList

func NewEnumerantList(s *C.Segment, sz int) Enumerant_List {
	return Enumerant_List(s.NewCompositeList(8, 2, sz))
}
func (s Enumerant_List) Len() int           { return C.PointerList(s).Len() }
func (s Enumerant_List) At(i int) Enumerant { return Enumerant(C.PointerList(s).At(i).ToStruct()) }
func (s Enumerant_List) ToArray() []Enumerant {
	return *(*[]Enumerant)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type Method C.Struct

func NewMethod(s *C.Segment) Method               { return Method(s.NewStruct(8, 4)) }
func NewRootMethod(s *C.Segment) Method           { return Method(s.NewRootStruct(8, 4)) }
func ReadRootMethod(s *C.Segment) Method          { return Method(s.Root(0).ToStruct()) }
func (s Method) Name() string                     { return C.Struct(s).GetObject(0).ToText() }
func (s Method) SetName(v string)                 { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s Method) CodeOrder() uint16                { return C.Struct(s).Get16(0) }
func (s Method) SetCodeOrder(v uint16)            { C.Struct(s).Set16(0, v) }
func (s Method) Params() MethodParam_List         { return MethodParam_List(C.Struct(s).GetObject(1)) }
func (s Method) SetParams(v MethodParam_List)     { C.Struct(s).SetObject(1, C.Object(v)) }
func (s Method) RequiredParamCount() uint16       { return C.Struct(s).Get16(2) }
func (s Method) SetRequiredParamCount(v uint16)   { C.Struct(s).Set16(2, v) }
func (s Method) ReturnType() Type                 { return Type(C.Struct(s).GetObject(2).ToStruct()) }
func (s Method) SetReturnType(v Type)             { C.Struct(s).SetObject(2, C.Object(v)) }
func (s Method) Annotations() Annotation_List     { return Annotation_List(C.Struct(s).GetObject(3)) }
func (s Method) SetAnnotations(v Annotation_List) { C.Struct(s).SetObject(3, C.Object(v)) }

type Method_List C.PointerList

func NewMethodList(s *C.Segment, sz int) Method_List { return Method_List(s.NewCompositeList(8, 4, sz)) }
func (s Method_List) Len() int                       { return C.PointerList(s).Len() }
func (s Method_List) At(i int) Method                { return Method(C.PointerList(s).At(i).ToStruct()) }
func (s Method_List) ToArray() []Method {
	return *(*[]Method)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type MethodParam C.Struct

func NewMethodParam(s *C.Segment) MethodParam          { return MethodParam(s.NewStruct(0, 4)) }
func NewRootMethodParam(s *C.Segment) MethodParam      { return MethodParam(s.NewRootStruct(0, 4)) }
func ReadRootMethodParam(s *C.Segment) MethodParam     { return MethodParam(s.Root(0).ToStruct()) }
func (s MethodParam) Name() string                     { return C.Struct(s).GetObject(0).ToText() }
func (s MethodParam) SetName(v string)                 { C.Struct(s).SetObject(0, s.Segment.NewText(v)) }
func (s MethodParam) Type() Type                       { return Type(C.Struct(s).GetObject(1).ToStruct()) }
func (s MethodParam) SetType(v Type)                   { C.Struct(s).SetObject(1, C.Object(v)) }
func (s MethodParam) DefaultValue() Value              { return Value(C.Struct(s).GetObject(2).ToStruct()) }
func (s MethodParam) SetDefaultValue(v Value)          { C.Struct(s).SetObject(2, C.Object(v)) }
func (s MethodParam) Annotations() Annotation_List     { return Annotation_List(C.Struct(s).GetObject(3)) }
func (s MethodParam) SetAnnotations(v Annotation_List) { C.Struct(s).SetObject(3, C.Object(v)) }

type MethodParam_List C.PointerList

func NewMethodParamList(s *C.Segment, sz int) MethodParam_List {
	return MethodParam_List(s.NewCompositeList(0, 4, sz))
}
func (s MethodParam_List) Len() int             { return C.PointerList(s).Len() }
func (s MethodParam_List) At(i int) MethodParam { return MethodParam(C.PointerList(s).At(i).ToStruct()) }
func (s MethodParam_List) ToArray() []MethodParam {
	return *(*[]MethodParam)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type Type C.Struct
type TypeList Type
type TypeEnum Type
type TypeStruct Type
type TypeInterface Type
type Type_Which uint16

const (
	TYPE_VOID      Type_Which = 0
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

func NewType(s *C.Segment) Type            { return Type(s.NewStruct(16, 1)) }
func NewRootType(s *C.Segment) Type        { return Type(s.NewRootStruct(16, 1)) }
func ReadRootType(s *C.Segment) Type       { return Type(s.Root(0).ToStruct()) }
func (s Type) Which() Type_Which           { return Type_Which(C.Struct(s).Get16(0)) }
func (s Type) List() TypeList              { return TypeList(s) }
func (s Type) SetList()                    { C.Struct(s).Set16(0, 14) }
func (s TypeList) ElementType() Type       { return Type(C.Struct(s).GetObject(0).ToStruct()) }
func (s TypeList) SetElementType(v Type)   { C.Struct(s).SetObject(0, C.Object(v)) }
func (s Type) Enum() TypeEnum              { return TypeEnum(s) }
func (s Type) SetEnum()                    { C.Struct(s).Set16(0, 15) }
func (s TypeEnum) TypeId() uint64          { return C.Struct(s).Get64(8) }
func (s TypeEnum) SetTypeId(v uint64)      { C.Struct(s).Set64(8, v) }
func (s Type) Struct() TypeStruct          { return TypeStruct(s) }
func (s Type) SetStruct()                  { C.Struct(s).Set16(0, 16) }
func (s TypeStruct) TypeId() uint64        { return C.Struct(s).Get64(8) }
func (s TypeStruct) SetTypeId(v uint64)    { C.Struct(s).Set64(8, v) }
func (s Type) Interface() TypeInterface    { return TypeInterface(s) }
func (s Type) SetInterface()               { C.Struct(s).Set16(0, 17) }
func (s TypeInterface) TypeId() uint64     { return C.Struct(s).Get64(8) }
func (s TypeInterface) SetTypeId(v uint64) { C.Struct(s).Set64(8, v) }

type Type_List C.PointerList

func NewTypeList(s *C.Segment, sz int) Type_List { return Type_List(s.NewCompositeList(16, 1, sz)) }
func (s Type_List) Len() int                     { return C.PointerList(s).Len() }
func (s Type_List) At(i int) Type                { return Type(C.PointerList(s).At(i).ToStruct()) }
func (s Type_List) ToArray() []Type              { return *(*[]Type)(unsafe.Pointer(C.PointerList(s).ToArray())) }

type Value C.Struct
type Value_Which uint16

const (
	VALUE_VOID      Value_Which = 0
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

func NewValue(s *C.Segment) Value      { return Value(s.NewStruct(16, 1)) }
func NewRootValue(s *C.Segment) Value  { return Value(s.NewRootStruct(16, 1)) }
func ReadRootValue(s *C.Segment) Value { return Value(s.Root(0).ToStruct()) }
func (s Value) Which() Value_Which     { return Value_Which(C.Struct(s).Get16(0)) }
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
func (s Value) List() C.Object       { return C.Struct(s).GetObject(0) }
func (s Value) SetList(v C.Object)   { C.Struct(s).Set16(0, 14); C.Struct(s).SetObject(0, v) }
func (s Value) Enum() uint16         { return C.Struct(s).Get16(2) }
func (s Value) SetEnum(v uint16)     { C.Struct(s).Set16(0, 15); C.Struct(s).Set16(2, v) }
func (s Value) Struct() C.Object     { return C.Struct(s).GetObject(0) }
func (s Value) SetStruct(v C.Object) { C.Struct(s).Set16(0, 16); C.Struct(s).SetObject(0, v) }
func (s Value) Object() C.Object     { return C.Struct(s).GetObject(0) }
func (s Value) SetObject(v C.Object) { C.Struct(s).Set16(0, 18); C.Struct(s).SetObject(0, v) }

type Value_List C.PointerList

func NewValueList(s *C.Segment, sz int) Value_List { return Value_List(s.NewCompositeList(16, 1, sz)) }
func (s Value_List) Len() int                      { return C.PointerList(s).Len() }
func (s Value_List) At(i int) Value                { return Value(C.PointerList(s).At(i).ToStruct()) }
func (s Value_List) ToArray() []Value              { return *(*[]Value)(unsafe.Pointer(C.PointerList(s).ToArray())) }

type Annotation C.Struct

func NewAnnotation(s *C.Segment) Annotation      { return Annotation(s.NewStruct(8, 1)) }
func NewRootAnnotation(s *C.Segment) Annotation  { return Annotation(s.NewRootStruct(8, 1)) }
func ReadRootAnnotation(s *C.Segment) Annotation { return Annotation(s.Root(0).ToStruct()) }
func (s Annotation) Id() uint64                  { return C.Struct(s).Get64(0) }
func (s Annotation) SetId(v uint64)              { C.Struct(s).Set64(0, v) }
func (s Annotation) Value() Value                { return Value(C.Struct(s).GetObject(0).ToStruct()) }
func (s Annotation) SetValue(v Value)            { C.Struct(s).SetObject(0, C.Object(v)) }

type Annotation_List C.PointerList

func NewAnnotationList(s *C.Segment, sz int) Annotation_List {
	return Annotation_List(s.NewCompositeList(8, 1, sz))
}
func (s Annotation_List) Len() int            { return C.PointerList(s).Len() }
func (s Annotation_List) At(i int) Annotation { return Annotation(C.PointerList(s).At(i).ToStruct()) }
func (s Annotation_List) ToArray() []Annotation {
	return *(*[]Annotation)(unsafe.Pointer(C.PointerList(s).ToArray()))
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

type ElementSize_List C.PointerList

func NewElementSizeList(s *C.Segment, sz int) ElementSize_List {
	return ElementSize_List(s.NewUInt16List(sz))
}
func (s ElementSize_List) Len() int             { return C.UInt16List(s).Len() }
func (s ElementSize_List) At(i int) ElementSize { return ElementSize(C.UInt16List(s).At(i)) }
func (s ElementSize_List) ToArray() []ElementSize {
	return *(*[]ElementSize)(unsafe.Pointer(C.UInt16List(s).ToEnumArray()))
}

type CodeGeneratorRequest C.Struct

func NewCodeGeneratorRequest(s *C.Segment) CodeGeneratorRequest {
	return CodeGeneratorRequest(s.NewStruct(0, 2))
}
func NewRootCodeGeneratorRequest(s *C.Segment) CodeGeneratorRequest {
	return CodeGeneratorRequest(s.NewRootStruct(0, 2))
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

type CodeGeneratorRequest_List C.PointerList

func NewCodeGeneratorRequestList(s *C.Segment, sz int) CodeGeneratorRequest_List {
	return CodeGeneratorRequest_List(s.NewCompositeList(0, 2, sz))
}
func (s CodeGeneratorRequest_List) Len() int { return C.PointerList(s).Len() }
func (s CodeGeneratorRequest_List) At(i int) CodeGeneratorRequest {
	return CodeGeneratorRequest(C.PointerList(s).At(i).ToStruct())
}
func (s CodeGeneratorRequest_List) ToArray() []CodeGeneratorRequest {
	return *(*[]CodeGeneratorRequest)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type CodeGeneratorRequestRequestedFile C.Struct

func NewCodeGeneratorRequestRequestedFile(s *C.Segment) CodeGeneratorRequestRequestedFile {
	return CodeGeneratorRequestRequestedFile(s.NewStruct(8, 2))
}
func NewRootCodeGeneratorRequestRequestedFile(s *C.Segment) CodeGeneratorRequestRequestedFile {
	return CodeGeneratorRequestRequestedFile(s.NewRootStruct(8, 2))
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

type CodeGeneratorRequestRequestedFile_List C.PointerList

func NewCodeGeneratorRequestRequestedFileList(s *C.Segment, sz int) CodeGeneratorRequestRequestedFile_List {
	return CodeGeneratorRequestRequestedFile_List(s.NewCompositeList(8, 2, sz))
}
func (s CodeGeneratorRequestRequestedFile_List) Len() int { return C.PointerList(s).Len() }
func (s CodeGeneratorRequestRequestedFile_List) At(i int) CodeGeneratorRequestRequestedFile {
	return CodeGeneratorRequestRequestedFile(C.PointerList(s).At(i).ToStruct())
}
func (s CodeGeneratorRequestRequestedFile_List) ToArray() []CodeGeneratorRequestRequestedFile {
	return *(*[]CodeGeneratorRequestRequestedFile)(unsafe.Pointer(C.PointerList(s).ToArray()))
}

type CodeGeneratorRequestRequestedFileImport C.Struct

func NewCodeGeneratorRequestRequestedFileImport(s *C.Segment) CodeGeneratorRequestRequestedFileImport {
	return CodeGeneratorRequestRequestedFileImport(s.NewStruct(8, 1))
}
func NewRootCodeGeneratorRequestRequestedFileImport(s *C.Segment) CodeGeneratorRequestRequestedFileImport {
	return CodeGeneratorRequestRequestedFileImport(s.NewRootStruct(8, 1))
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

type CodeGeneratorRequestRequestedFileImport_List C.PointerList

func NewCodeGeneratorRequestRequestedFileImportList(s *C.Segment, sz int) CodeGeneratorRequestRequestedFileImport_List {
	return CodeGeneratorRequestRequestedFileImport_List(s.NewCompositeList(8, 1, sz))
}
func (s CodeGeneratorRequestRequestedFileImport_List) Len() int { return C.PointerList(s).Len() }
func (s CodeGeneratorRequestRequestedFileImport_List) At(i int) CodeGeneratorRequestRequestedFileImport {
	return CodeGeneratorRequestRequestedFileImport(C.PointerList(s).At(i).ToStruct())
}
func (s CodeGeneratorRequestRequestedFileImport_List) ToArray() []CodeGeneratorRequestRequestedFileImport {
	return *(*[]CodeGeneratorRequestRequestedFileImport)(unsafe.Pointer(C.PointerList(s).ToArray()))
}
