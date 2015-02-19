package main

// AUTO GENERATED - DO NOT EDIT

import (
	math "math"
	strconv "strconv"
	C "zombiezen.com/go/capnproto"
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

func (w Node_Which) String() string {
	const s = "filestructenuminterfaceconstannotation"
	switch w {
	case Node_Which_file:
		return s[0:4]
	case Node_Which_struct:
		return s[4:10]
	case Node_Which_enum:
		return s[10:14]
	case Node_Which_interface:
		return s[14:23]
	case Node_Which_const:
		return s[23:28]
	case Node_Which_annotation:
		return s[28:38]

	}
	return "Node_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
}

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

// capnp.JSON_enabled == false so we stub MarshalJSON().
func (s Node) MarshalJSON() (bs []byte, err error) { return }

type Node_List C.PointerList

func NewNode_List(s *C.Segment, sz int) Node_List {
	return Node_List(s.NewCompositeList(C.ObjectSize{DataSize: 40, PointerCount: 6}, sz))
}
func (s Node_List) Len() int             { return C.PointerList(s).Len() }
func (s Node_List) At(i int) Node        { return Node(C.PointerList(s).At(i).ToStruct()) }
func (s Node_List) Set(i int, item Node) { C.PointerList(s).Set(i, C.Object(item)) }

type Node_Promise C.Pipeline

func (p *Node_Promise) Get() (Node, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Node(s), err
}
func (p *Node_Promise) Struct() *Node_struct_Promise { return (*Node_struct_Promise)(p) }

type Node_struct_Promise C.Pipeline

func (p *Node_struct_Promise) Get() (Node_struct, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Node_struct(s), err
}
func (p *Node_Promise) Enum() *Node_enum_Promise { return (*Node_enum_Promise)(p) }

type Node_enum_Promise C.Pipeline

func (p *Node_enum_Promise) Get() (Node_enum, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Node_enum(s), err
}
func (p *Node_Promise) Interface() *Node_interface_Promise { return (*Node_interface_Promise)(p) }

type Node_interface_Promise C.Pipeline

func (p *Node_interface_Promise) Get() (Node_interface, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Node_interface(s), err
}
func (p *Node_Promise) Const() *Node_const_Promise { return (*Node_const_Promise)(p) }

type Node_const_Promise C.Pipeline

func (p *Node_const_Promise) Get() (Node_const, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Node_const(s), err
}

func (p *Node_const_Promise) Type() *Type_Promise {
	return (*Type_Promise)((*C.Pipeline)(p).GetPipeline(3))
}

func (p *Node_const_Promise) Value() *Value_Promise {
	return (*Value_Promise)((*C.Pipeline)(p).GetPipeline(4))
}
func (p *Node_Promise) Annotation() *Node_annotation_Promise { return (*Node_annotation_Promise)(p) }

type Node_annotation_Promise C.Pipeline

func (p *Node_annotation_Promise) Get() (Node_annotation, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Node_annotation(s), err
}

func (p *Node_annotation_Promise) Type() *Type_Promise {
	return (*Type_Promise)((*C.Pipeline)(p).GetPipeline(3))
}

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

// capnp.JSON_enabled == false so we stub MarshalJSON().
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

type Node_Parameter_Promise C.Pipeline

func (p *Node_Parameter_Promise) Get() (Node_Parameter, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Node_Parameter(s), err
}

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

// capnp.JSON_enabled == false so we stub MarshalJSON().
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

type Node_NestedNode_Promise C.Pipeline

func (p *Node_NestedNode_Promise) Get() (Node_NestedNode, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Node_NestedNode(s), err
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

func (w Field_Which) String() string {
	const s = "slotgroup"
	switch w {
	case Field_Which_slot:
		return s[0:4]
	case Field_Which_group:
		return s[4:9]

	}
	return "Field_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
}

type Field_ordinal_Which uint16

const (
	Field_ordinal_Which_implicit Field_ordinal_Which = 0
	Field_ordinal_Which_explicit Field_ordinal_Which = 1
)

func (w Field_ordinal_Which) String() string {
	const s = "implicitexplicit"
	switch w {
	case Field_ordinal_Which_implicit:
		return s[0:8]
	case Field_ordinal_Which_explicit:
		return s[8:16]

	}
	return "Field_ordinal_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
}

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

// capnp.JSON_enabled == false so we stub MarshalJSON().
func (s Field) MarshalJSON() (bs []byte, err error) { return }

type Field_List C.PointerList

func NewField_List(s *C.Segment, sz int) Field_List {
	return Field_List(s.NewCompositeList(C.ObjectSize{DataSize: 24, PointerCount: 4}, sz))
}
func (s Field_List) Len() int              { return C.PointerList(s).Len() }
func (s Field_List) At(i int) Field        { return Field(C.PointerList(s).At(i).ToStruct()) }
func (s Field_List) Set(i int, item Field) { C.PointerList(s).Set(i, C.Object(item)) }

type Field_Promise C.Pipeline

func (p *Field_Promise) Get() (Field, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Field(s), err
}
func (p *Field_Promise) Slot() *Field_slot_Promise { return (*Field_slot_Promise)(p) }

type Field_slot_Promise C.Pipeline

func (p *Field_slot_Promise) Get() (Field_slot, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Field_slot(s), err
}

func (p *Field_slot_Promise) Type() *Type_Promise {
	return (*Type_Promise)((*C.Pipeline)(p).GetPipeline(2))
}

func (p *Field_slot_Promise) DefaultValue() *Value_Promise {
	return (*Value_Promise)((*C.Pipeline)(p).GetPipeline(3))
}
func (p *Field_Promise) Group() *Field_group_Promise { return (*Field_group_Promise)(p) }

type Field_group_Promise C.Pipeline

func (p *Field_group_Promise) Get() (Field_group, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Field_group(s), err
}
func (p *Field_Promise) Ordinal() *Field_ordinal_Promise { return (*Field_ordinal_Promise)(p) }

type Field_ordinal_Promise C.Pipeline

func (p *Field_ordinal_Promise) Get() (Field_ordinal, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Field_ordinal(s), err
}

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

// capnp.JSON_enabled == false so we stub MarshalJSON().
func (s Enumerant) MarshalJSON() (bs []byte, err error) { return }

type Enumerant_List C.PointerList

func NewEnumerant_List(s *C.Segment, sz int) Enumerant_List {
	return Enumerant_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 2}, sz))
}
func (s Enumerant_List) Len() int                  { return C.PointerList(s).Len() }
func (s Enumerant_List) At(i int) Enumerant        { return Enumerant(C.PointerList(s).At(i).ToStruct()) }
func (s Enumerant_List) Set(i int, item Enumerant) { C.PointerList(s).Set(i, C.Object(item)) }

type Enumerant_Promise C.Pipeline

func (p *Enumerant_Promise) Get() (Enumerant, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Enumerant(s), err
}

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

// capnp.JSON_enabled == false so we stub MarshalJSON().
func (s Superclass) MarshalJSON() (bs []byte, err error) { return }

type Superclass_List C.PointerList

func NewSuperclass_List(s *C.Segment, sz int) Superclass_List {
	return Superclass_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 1}, sz))
}
func (s Superclass_List) Len() int                   { return C.PointerList(s).Len() }
func (s Superclass_List) At(i int) Superclass        { return Superclass(C.PointerList(s).At(i).ToStruct()) }
func (s Superclass_List) Set(i int, item Superclass) { C.PointerList(s).Set(i, C.Object(item)) }

type Superclass_Promise C.Pipeline

func (p *Superclass_Promise) Get() (Superclass, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Superclass(s), err
}

func (p *Superclass_Promise) Brand() *Brand_Promise {
	return (*Brand_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

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

// capnp.JSON_enabled == false so we stub MarshalJSON().
func (s Method) MarshalJSON() (bs []byte, err error) { return }

type Method_List C.PointerList

func NewMethod_List(s *C.Segment, sz int) Method_List {
	return Method_List(s.NewCompositeList(C.ObjectSize{DataSize: 24, PointerCount: 5}, sz))
}
func (s Method_List) Len() int               { return C.PointerList(s).Len() }
func (s Method_List) At(i int) Method        { return Method(C.PointerList(s).At(i).ToStruct()) }
func (s Method_List) Set(i int, item Method) { C.PointerList(s).Set(i, C.Object(item)) }

type Method_Promise C.Pipeline

func (p *Method_Promise) Get() (Method, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Method(s), err
}

func (p *Method_Promise) ParamBrand() *Brand_Promise {
	return (*Brand_Promise)((*C.Pipeline)(p).GetPipeline(2))
}

func (p *Method_Promise) ResultBrand() *Brand_Promise {
	return (*Brand_Promise)((*C.Pipeline)(p).GetPipeline(3))
}

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

func (w Type_Which) String() string {
	const s = "voidboolint8int16int32int64uint8uint16uint32uint64float32float64textdatalistenumstructinterfaceanyPointer"
	switch w {
	case Type_Which_void:
		return s[0:4]
	case Type_Which_bool:
		return s[4:8]
	case Type_Which_int8:
		return s[8:12]
	case Type_Which_int16:
		return s[12:17]
	case Type_Which_int32:
		return s[17:22]
	case Type_Which_int64:
		return s[22:27]
	case Type_Which_uint8:
		return s[27:32]
	case Type_Which_uint16:
		return s[32:38]
	case Type_Which_uint32:
		return s[38:44]
	case Type_Which_uint64:
		return s[44:50]
	case Type_Which_float32:
		return s[50:57]
	case Type_Which_float64:
		return s[57:64]
	case Type_Which_text:
		return s[64:68]
	case Type_Which_data:
		return s[68:72]
	case Type_Which_list:
		return s[72:76]
	case Type_Which_enum:
		return s[76:80]
	case Type_Which_struct:
		return s[80:86]
	case Type_Which_interface:
		return s[86:95]
	case Type_Which_anyPointer:
		return s[95:105]

	}
	return "Type_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
}

type Type_anyPointer_Which uint16

const (
	Type_anyPointer_Which_unconstrained           Type_anyPointer_Which = 0
	Type_anyPointer_Which_parameter               Type_anyPointer_Which = 1
	Type_anyPointer_Which_implicitMethodParameter Type_anyPointer_Which = 2
)

func (w Type_anyPointer_Which) String() string {
	const s = "unconstrainedparameterimplicitMethodParameter"
	switch w {
	case Type_anyPointer_Which_unconstrained:
		return s[0:13]
	case Type_anyPointer_Which_parameter:
		return s[13:22]
	case Type_anyPointer_Which_implicitMethodParameter:
		return s[22:45]

	}
	return "Type_anyPointer_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
}

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

// capnp.JSON_enabled == false so we stub MarshalJSON().
func (s Type) MarshalJSON() (bs []byte, err error) { return }

type Type_List C.PointerList

func NewType_List(s *C.Segment, sz int) Type_List {
	return Type_List(s.NewCompositeList(C.ObjectSize{DataSize: 24, PointerCount: 1}, sz))
}
func (s Type_List) Len() int             { return C.PointerList(s).Len() }
func (s Type_List) At(i int) Type        { return Type(C.PointerList(s).At(i).ToStruct()) }
func (s Type_List) Set(i int, item Type) { C.PointerList(s).Set(i, C.Object(item)) }

type Type_Promise C.Pipeline

func (p *Type_Promise) Get() (Type, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Type(s), err
}
func (p *Type_Promise) List() *Type_list_Promise { return (*Type_list_Promise)(p) }

type Type_list_Promise C.Pipeline

func (p *Type_list_Promise) Get() (Type_list, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Type_list(s), err
}

func (p *Type_list_Promise) ElementType() *Type_Promise {
	return (*Type_Promise)((*C.Pipeline)(p).GetPipeline(0))
}
func (p *Type_Promise) Enum() *Type_enum_Promise { return (*Type_enum_Promise)(p) }

type Type_enum_Promise C.Pipeline

func (p *Type_enum_Promise) Get() (Type_enum, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Type_enum(s), err
}

func (p *Type_enum_Promise) Brand() *Brand_Promise {
	return (*Brand_Promise)((*C.Pipeline)(p).GetPipeline(0))
}
func (p *Type_Promise) Struct() *Type_struct_Promise { return (*Type_struct_Promise)(p) }

type Type_struct_Promise C.Pipeline

func (p *Type_struct_Promise) Get() (Type_struct, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Type_struct(s), err
}

func (p *Type_struct_Promise) Brand() *Brand_Promise {
	return (*Brand_Promise)((*C.Pipeline)(p).GetPipeline(0))
}
func (p *Type_Promise) Interface() *Type_interface_Promise { return (*Type_interface_Promise)(p) }

type Type_interface_Promise C.Pipeline

func (p *Type_interface_Promise) Get() (Type_interface, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Type_interface(s), err
}

func (p *Type_interface_Promise) Brand() *Brand_Promise {
	return (*Brand_Promise)((*C.Pipeline)(p).GetPipeline(0))
}
func (p *Type_Promise) AnyPointer() *Type_anyPointer_Promise { return (*Type_anyPointer_Promise)(p) }

type Type_anyPointer_Promise C.Pipeline

func (p *Type_anyPointer_Promise) Get() (Type_anyPointer, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Type_anyPointer(s), err
}
func (p *Type_anyPointer_Promise) Parameter() *Type_anyPointer_parameter_Promise {
	return (*Type_anyPointer_parameter_Promise)(p)
}

type Type_anyPointer_parameter_Promise C.Pipeline

func (p *Type_anyPointer_parameter_Promise) Get() (Type_anyPointer_parameter, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Type_anyPointer_parameter(s), err
}
func (p *Type_anyPointer_Promise) ImplicitMethodParameter() *Type_anyPointer_implicitMethodParameter_Promise {
	return (*Type_anyPointer_implicitMethodParameter_Promise)(p)
}

type Type_anyPointer_implicitMethodParameter_Promise C.Pipeline

func (p *Type_anyPointer_implicitMethodParameter_Promise) Get() (Type_anyPointer_implicitMethodParameter, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Type_anyPointer_implicitMethodParameter(s), err
}

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

// capnp.JSON_enabled == false so we stub MarshalJSON().
func (s Brand) MarshalJSON() (bs []byte, err error) { return }

type Brand_List C.PointerList

func NewBrand_List(s *C.Segment, sz int) Brand_List {
	return Brand_List(s.NewCompositeList(C.ObjectSize{DataSize: 0, PointerCount: 1}, sz))
}
func (s Brand_List) Len() int              { return C.PointerList(s).Len() }
func (s Brand_List) At(i int) Brand        { return Brand(C.PointerList(s).At(i).ToStruct()) }
func (s Brand_List) Set(i int, item Brand) { C.PointerList(s).Set(i, C.Object(item)) }

type Brand_Promise C.Pipeline

func (p *Brand_Promise) Get() (Brand, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Brand(s), err
}

type Brand_Scope C.Struct
type Brand_Scope_Which uint16

const (
	Brand_Scope_Which_bind    Brand_Scope_Which = 0
	Brand_Scope_Which_inherit Brand_Scope_Which = 1
)

func (w Brand_Scope_Which) String() string {
	const s = "bindinherit"
	switch w {
	case Brand_Scope_Which_bind:
		return s[0:4]
	case Brand_Scope_Which_inherit:
		return s[4:11]

	}
	return "Brand_Scope_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
}

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

// capnp.JSON_enabled == false so we stub MarshalJSON().
func (s Brand_Scope) MarshalJSON() (bs []byte, err error) { return }

type Brand_Scope_List C.PointerList

func NewBrand_Scope_List(s *C.Segment, sz int) Brand_Scope_List {
	return Brand_Scope_List(s.NewCompositeList(C.ObjectSize{DataSize: 16, PointerCount: 1}, sz))
}
func (s Brand_Scope_List) Len() int                    { return C.PointerList(s).Len() }
func (s Brand_Scope_List) At(i int) Brand_Scope        { return Brand_Scope(C.PointerList(s).At(i).ToStruct()) }
func (s Brand_Scope_List) Set(i int, item Brand_Scope) { C.PointerList(s).Set(i, C.Object(item)) }

type Brand_Scope_Promise C.Pipeline

func (p *Brand_Scope_Promise) Get() (Brand_Scope, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Brand_Scope(s), err
}

type Brand_Binding C.Struct
type Brand_Binding_Which uint16

const (
	Brand_Binding_Which_unbound Brand_Binding_Which = 0
	Brand_Binding_Which_type    Brand_Binding_Which = 1
)

func (w Brand_Binding_Which) String() string {
	const s = "unboundtype"
	switch w {
	case Brand_Binding_Which_unbound:
		return s[0:7]
	case Brand_Binding_Which_type:
		return s[7:11]

	}
	return "Brand_Binding_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
}

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

// capnp.JSON_enabled == false so we stub MarshalJSON().
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

type Brand_Binding_Promise C.Pipeline

func (p *Brand_Binding_Promise) Get() (Brand_Binding, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Brand_Binding(s), err
}

func (p *Brand_Binding_Promise) Type() *Type_Promise {
	return (*Type_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

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

func (w Value_Which) String() string {
	const s = "voidboolint8int16int32int64uint8uint16uint32uint64float32float64textdatalistenumstructinterfaceanyPointer"
	switch w {
	case Value_Which_void:
		return s[0:4]
	case Value_Which_bool:
		return s[4:8]
	case Value_Which_int8:
		return s[8:12]
	case Value_Which_int16:
		return s[12:17]
	case Value_Which_int32:
		return s[17:22]
	case Value_Which_int64:
		return s[22:27]
	case Value_Which_uint8:
		return s[27:32]
	case Value_Which_uint16:
		return s[32:38]
	case Value_Which_uint32:
		return s[38:44]
	case Value_Which_uint64:
		return s[44:50]
	case Value_Which_float32:
		return s[50:57]
	case Value_Which_float64:
		return s[57:64]
	case Value_Which_text:
		return s[64:68]
	case Value_Which_data:
		return s[68:72]
	case Value_Which_list:
		return s[72:76]
	case Value_Which_enum:
		return s[76:80]
	case Value_Which_struct:
		return s[80:86]
	case Value_Which_interface:
		return s[86:95]
	case Value_Which_anyPointer:
		return s[95:105]

	}
	return "Value_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
}

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

// capnp.JSON_enabled == false so we stub MarshalJSON().
func (s Value) MarshalJSON() (bs []byte, err error) { return }

type Value_List C.PointerList

func NewValue_List(s *C.Segment, sz int) Value_List {
	return Value_List(s.NewCompositeList(C.ObjectSize{DataSize: 16, PointerCount: 1}, sz))
}
func (s Value_List) Len() int              { return C.PointerList(s).Len() }
func (s Value_List) At(i int) Value        { return Value(C.PointerList(s).At(i).ToStruct()) }
func (s Value_List) Set(i int, item Value) { C.PointerList(s).Set(i, C.Object(item)) }

type Value_Promise C.Pipeline

func (p *Value_Promise) Get() (Value, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Value(s), err
}

func (p *Value_Promise) List() *C.Pipeline {
	return (*C.Pipeline)(p).GetPipeline(0)
}

func (p *Value_Promise) Struct() *C.Pipeline {
	return (*C.Pipeline)(p).GetPipeline(0)
}

func (p *Value_Promise) AnyPointer() *C.Pipeline {
	return (*C.Pipeline)(p).GetPipeline(0)
}

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

// capnp.JSON_enabled == false so we stub MarshalJSON().
func (s Annotation) MarshalJSON() (bs []byte, err error) { return }

type Annotation_List C.PointerList

func NewAnnotation_List(s *C.Segment, sz int) Annotation_List {
	return Annotation_List(s.NewCompositeList(C.ObjectSize{DataSize: 8, PointerCount: 2}, sz))
}
func (s Annotation_List) Len() int                   { return C.PointerList(s).Len() }
func (s Annotation_List) At(i int) Annotation        { return Annotation(C.PointerList(s).At(i).ToStruct()) }
func (s Annotation_List) Set(i int, item Annotation) { C.PointerList(s).Set(i, C.Object(item)) }

type Annotation_Promise C.Pipeline

func (p *Annotation_Promise) Get() (Annotation, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return Annotation(s), err
}

func (p *Annotation_Promise) Brand() *Brand_Promise {
	return (*Brand_Promise)((*C.Pipeline)(p).GetPipeline(1))
}

func (p *Annotation_Promise) Value() *Value_Promise {
	return (*Value_Promise)((*C.Pipeline)(p).GetPipeline(0))
}

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

// capnp.JSON_enabled == false so we stub MarshalJSON().
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

// capnp.JSON_enabled == false so we stub MarshalJSON().
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

type CodeGeneratorRequest_Promise C.Pipeline

func (p *CodeGeneratorRequest_Promise) Get() (CodeGeneratorRequest, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return CodeGeneratorRequest(s), err
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

// capnp.JSON_enabled == false so we stub MarshalJSON().
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

type CodeGeneratorRequest_RequestedFile_Promise C.Pipeline

func (p *CodeGeneratorRequest_RequestedFile_Promise) Get() (CodeGeneratorRequest_RequestedFile, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return CodeGeneratorRequest_RequestedFile(s), err
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

// capnp.JSON_enabled == false so we stub MarshalJSON().
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

type CodeGeneratorRequest_RequestedFile_Import_Promise C.Pipeline

func (p *CodeGeneratorRequest_RequestedFile_Import_Promise) Get() (CodeGeneratorRequest_RequestedFile_Import, error) {
	s, err := (*C.Pipeline)(p).Struct()
	return CodeGeneratorRequest_RequestedFile_Import(s), err
}
