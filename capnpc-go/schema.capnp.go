package main

// AUTO GENERATED - DO NOT EDIT

import (
	math "math"
	strconv "strconv"
	capnp "zombiezen.com/go/capnproto2"
)

const (
	Field_noDiscriminant = uint16(65535)
)

type Node struct{ capnp.Struct }
type Node_structGroup Node
type Node_enum Node
type Node_interface Node
type Node_const Node
type Node_annotation Node
type Node_Which uint16

const (
	Node_Which_file        Node_Which = 0
	Node_Which_structGroup Node_Which = 1
	Node_Which_enum        Node_Which = 2
	Node_Which_interface   Node_Which = 3
	Node_Which_const       Node_Which = 4
	Node_Which_annotation  Node_Which = 5
)

func (w Node_Which) String() string {
	const s = "filestructGroupenuminterfaceconstannotation"
	switch w {
	case Node_Which_file:
		return s[0:4]
	case Node_Which_structGroup:
		return s[4:15]
	case Node_Which_enum:
		return s[15:19]
	case Node_Which_interface:
		return s[19:28]
	case Node_Which_const:
		return s[28:33]
	case Node_Which_annotation:
		return s[33:43]

	}
	return "Node_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
}

func NewNode(s *capnp.Segment) (Node, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 40, PointerCount: 6})
	if err != nil {
		return Node{}, err
	}
	return Node{st}, nil
}

func NewRootNode(s *capnp.Segment) (Node, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 40, PointerCount: 6})
	if err != nil {
		return Node{}, err
	}
	return Node{st}, nil
}

func ReadRootNode(msg *capnp.Message) (Node, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Node{}, err
	}
	return Node{root.Struct()}, nil
}

func (s Node) Which() Node_Which {
	return Node_Which(s.Struct.Uint16(12))
}

func (s Node) Id() uint64 {
	return s.Struct.Uint64(0)
}

func (s Node) SetId(v uint64) {

	s.Struct.SetUint64(0, v)
}

func (s Node) DisplayName() (string, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return "", err
	}

	return capnp.PtrToText(p), nil

}

func (s Node) DisplayNameBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	return capnp.PtrToData(p), nil
}

func (s Node) SetDisplayName(v string) error {

	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, t.List.ToPtr())
}

func (s Node) DisplayNamePrefixLength() uint32 {
	return s.Struct.Uint32(8)
}

func (s Node) SetDisplayNamePrefixLength(v uint32) {

	s.Struct.SetUint32(8, v)
}

func (s Node) ScopeId() uint64 {
	return s.Struct.Uint64(16)
}

func (s Node) SetScopeId(v uint64) {

	s.Struct.SetUint64(16, v)
}

func (s Node) Parameters() (Node_Parameter_List, error) {
	p, err := s.Struct.Ptr(5)
	if err != nil {
		return Node_Parameter_List{}, err
	}

	return Node_Parameter_List{List: p.List()}, nil

}

func (s Node) SetParameters(v Node_Parameter_List) error {

	return s.Struct.SetPtr(5, v.List.ToPtr())
}

func (s Node) IsGeneric() bool {
	return s.Struct.Bit(288)
}

func (s Node) SetIsGeneric(v bool) {

	s.Struct.SetBit(288, v)
}

func (s Node) NestedNodes() (Node_NestedNode_List, error) {
	p, err := s.Struct.Ptr(1)
	if err != nil {
		return Node_NestedNode_List{}, err
	}

	return Node_NestedNode_List{List: p.List()}, nil

}

func (s Node) SetNestedNodes(v Node_NestedNode_List) error {

	return s.Struct.SetPtr(1, v.List.ToPtr())
}

func (s Node) Annotations() (Annotation_List, error) {
	p, err := s.Struct.Ptr(2)
	if err != nil {
		return Annotation_List{}, err
	}

	return Annotation_List{List: p.List()}, nil

}

func (s Node) SetAnnotations(v Annotation_List) error {

	return s.Struct.SetPtr(2, v.List.ToPtr())
}

func (s Node) SetFile() {
	s.Struct.SetUint16(12, 0)
}
func (s Node) StructGroup() Node_structGroup { return Node_structGroup(s) }

func (s Node) SetStructGroup() { s.Struct.SetUint16(12, 1) }

func (s Node_structGroup) DataWordCount() uint16 {
	return s.Struct.Uint16(14)
}

func (s Node_structGroup) SetDataWordCount(v uint16) {

	s.Struct.SetUint16(14, v)
}

func (s Node_structGroup) PointerCount() uint16 {
	return s.Struct.Uint16(24)
}

func (s Node_structGroup) SetPointerCount(v uint16) {

	s.Struct.SetUint16(24, v)
}

func (s Node_structGroup) PreferredListEncoding() ElementSize {
	return ElementSize(s.Struct.Uint16(26))
}

func (s Node_structGroup) SetPreferredListEncoding(v ElementSize) {

	s.Struct.SetUint16(26, uint16(v))
}

func (s Node_structGroup) IsGroup() bool {
	return s.Struct.Bit(224)
}

func (s Node_structGroup) SetIsGroup(v bool) {

	s.Struct.SetBit(224, v)
}

func (s Node_structGroup) DiscriminantCount() uint16 {
	return s.Struct.Uint16(30)
}

func (s Node_structGroup) SetDiscriminantCount(v uint16) {

	s.Struct.SetUint16(30, v)
}

func (s Node_structGroup) DiscriminantOffset() uint32 {
	return s.Struct.Uint32(32)
}

func (s Node_structGroup) SetDiscriminantOffset(v uint32) {

	s.Struct.SetUint32(32, v)
}

func (s Node_structGroup) Fields() (Field_List, error) {
	p, err := s.Struct.Ptr(3)
	if err != nil {
		return Field_List{}, err
	}

	return Field_List{List: p.List()}, nil

}

func (s Node_structGroup) SetFields(v Field_List) error {

	return s.Struct.SetPtr(3, v.List.ToPtr())
}
func (s Node) Enum() Node_enum { return Node_enum(s) }

func (s Node) SetEnum() { s.Struct.SetUint16(12, 2) }

func (s Node_enum) Enumerants() (Enumerant_List, error) {
	p, err := s.Struct.Ptr(3)
	if err != nil {
		return Enumerant_List{}, err
	}

	return Enumerant_List{List: p.List()}, nil

}

func (s Node_enum) SetEnumerants(v Enumerant_List) error {

	return s.Struct.SetPtr(3, v.List.ToPtr())
}
func (s Node) Interface() Node_interface { return Node_interface(s) }

func (s Node) SetInterface() { s.Struct.SetUint16(12, 3) }

func (s Node_interface) Methods() (Method_List, error) {
	p, err := s.Struct.Ptr(3)
	if err != nil {
		return Method_List{}, err
	}

	return Method_List{List: p.List()}, nil

}

func (s Node_interface) SetMethods(v Method_List) error {

	return s.Struct.SetPtr(3, v.List.ToPtr())
}

func (s Node_interface) Superclasses() (Superclass_List, error) {
	p, err := s.Struct.Ptr(4)
	if err != nil {
		return Superclass_List{}, err
	}

	return Superclass_List{List: p.List()}, nil

}

func (s Node_interface) SetSuperclasses(v Superclass_List) error {

	return s.Struct.SetPtr(4, v.List.ToPtr())
}
func (s Node) Const() Node_const { return Node_const(s) }

func (s Node) SetConst() { s.Struct.SetUint16(12, 4) }

func (s Node_const) Type() (Type, error) {
	p, err := s.Struct.Ptr(3)
	if err != nil {
		return Type{}, err
	}

	return Type{Struct: p.Struct()}, nil

}

func (s Node_const) SetType(v Type) error {

	return s.Struct.SetPtr(3, v.Struct.ToPtr())
}

// NewType sets the type field to a newly
// allocated Type struct, preferring placement in s's segment.
func (s Node_const) NewType() (Type, error) {

	ss, err := NewType(s.Struct.Segment())
	if err != nil {
		return Type{}, err
	}
	err = s.Struct.SetPtr(3, ss.Struct.ToPtr())
	return ss, err
}

func (s Node_const) Value() (Value, error) {
	p, err := s.Struct.Ptr(4)
	if err != nil {
		return Value{}, err
	}

	return Value{Struct: p.Struct()}, nil

}

func (s Node_const) SetValue(v Value) error {

	return s.Struct.SetPtr(4, v.Struct.ToPtr())
}

// NewValue sets the value field to a newly
// allocated Value struct, preferring placement in s's segment.
func (s Node_const) NewValue() (Value, error) {

	ss, err := NewValue(s.Struct.Segment())
	if err != nil {
		return Value{}, err
	}
	err = s.Struct.SetPtr(4, ss.Struct.ToPtr())
	return ss, err
}
func (s Node) Annotation() Node_annotation { return Node_annotation(s) }

func (s Node) SetAnnotation() { s.Struct.SetUint16(12, 5) }

func (s Node_annotation) Type() (Type, error) {
	p, err := s.Struct.Ptr(3)
	if err != nil {
		return Type{}, err
	}

	return Type{Struct: p.Struct()}, nil

}

func (s Node_annotation) SetType(v Type) error {

	return s.Struct.SetPtr(3, v.Struct.ToPtr())
}

// NewType sets the type field to a newly
// allocated Type struct, preferring placement in s's segment.
func (s Node_annotation) NewType() (Type, error) {

	ss, err := NewType(s.Struct.Segment())
	if err != nil {
		return Type{}, err
	}
	err = s.Struct.SetPtr(3, ss.Struct.ToPtr())
	return ss, err
}

func (s Node_annotation) TargetsFile() bool {
	return s.Struct.Bit(112)
}

func (s Node_annotation) SetTargetsFile(v bool) {

	s.Struct.SetBit(112, v)
}

func (s Node_annotation) TargetsConst() bool {
	return s.Struct.Bit(113)
}

func (s Node_annotation) SetTargetsConst(v bool) {

	s.Struct.SetBit(113, v)
}

func (s Node_annotation) TargetsEnum() bool {
	return s.Struct.Bit(114)
}

func (s Node_annotation) SetTargetsEnum(v bool) {

	s.Struct.SetBit(114, v)
}

func (s Node_annotation) TargetsEnumerant() bool {
	return s.Struct.Bit(115)
}

func (s Node_annotation) SetTargetsEnumerant(v bool) {

	s.Struct.SetBit(115, v)
}

func (s Node_annotation) TargetsStruct() bool {
	return s.Struct.Bit(116)
}

func (s Node_annotation) SetTargetsStruct(v bool) {

	s.Struct.SetBit(116, v)
}

func (s Node_annotation) TargetsField() bool {
	return s.Struct.Bit(117)
}

func (s Node_annotation) SetTargetsField(v bool) {

	s.Struct.SetBit(117, v)
}

func (s Node_annotation) TargetsUnion() bool {
	return s.Struct.Bit(118)
}

func (s Node_annotation) SetTargetsUnion(v bool) {

	s.Struct.SetBit(118, v)
}

func (s Node_annotation) TargetsGroup() bool {
	return s.Struct.Bit(119)
}

func (s Node_annotation) SetTargetsGroup(v bool) {

	s.Struct.SetBit(119, v)
}

func (s Node_annotation) TargetsInterface() bool {
	return s.Struct.Bit(120)
}

func (s Node_annotation) SetTargetsInterface(v bool) {

	s.Struct.SetBit(120, v)
}

func (s Node_annotation) TargetsMethod() bool {
	return s.Struct.Bit(121)
}

func (s Node_annotation) SetTargetsMethod(v bool) {

	s.Struct.SetBit(121, v)
}

func (s Node_annotation) TargetsParam() bool {
	return s.Struct.Bit(122)
}

func (s Node_annotation) SetTargetsParam(v bool) {

	s.Struct.SetBit(122, v)
}

func (s Node_annotation) TargetsAnnotation() bool {
	return s.Struct.Bit(123)
}

func (s Node_annotation) SetTargetsAnnotation(v bool) {

	s.Struct.SetBit(123, v)
}

// Node_List is a list of Node.
type Node_List struct{ capnp.List }

// NewNode creates a new list of Node.
func NewNode_List(s *capnp.Segment, sz int32) (Node_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 40, PointerCount: 6}, sz)
	if err != nil {
		return Node_List{}, err
	}
	return Node_List{l}, nil
}

func (s Node_List) At(i int) Node           { return Node{s.List.Struct(i)} }
func (s Node_List) Set(i int, v Node) error { return s.List.SetStruct(i, v.Struct) }

type Node_Parameter struct{ capnp.Struct }

func NewNode_Parameter(s *capnp.Segment) (Node_Parameter, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Node_Parameter{}, err
	}
	return Node_Parameter{st}, nil
}

func NewRootNode_Parameter(s *capnp.Segment) (Node_Parameter, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Node_Parameter{}, err
	}
	return Node_Parameter{st}, nil
}

func ReadRootNode_Parameter(msg *capnp.Message) (Node_Parameter, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Node_Parameter{}, err
	}
	return Node_Parameter{root.Struct()}, nil
}

func (s Node_Parameter) Name() (string, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return "", err
	}

	return capnp.PtrToText(p), nil

}

func (s Node_Parameter) NameBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	return capnp.PtrToData(p), nil
}

func (s Node_Parameter) SetName(v string) error {

	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, t.List.ToPtr())
}

// Node_Parameter_List is a list of Node_Parameter.
type Node_Parameter_List struct{ capnp.List }

// NewNode_Parameter creates a new list of Node_Parameter.
func NewNode_Parameter_List(s *capnp.Segment, sz int32) (Node_Parameter_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return Node_Parameter_List{}, err
	}
	return Node_Parameter_List{l}, nil
}

func (s Node_Parameter_List) At(i int) Node_Parameter           { return Node_Parameter{s.List.Struct(i)} }
func (s Node_Parameter_List) Set(i int, v Node_Parameter) error { return s.List.SetStruct(i, v.Struct) }

type Node_NestedNode struct{ capnp.Struct }

func NewNode_NestedNode(s *capnp.Segment) (Node_NestedNode, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		return Node_NestedNode{}, err
	}
	return Node_NestedNode{st}, nil
}

func NewRootNode_NestedNode(s *capnp.Segment) (Node_NestedNode, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		return Node_NestedNode{}, err
	}
	return Node_NestedNode{st}, nil
}

func ReadRootNode_NestedNode(msg *capnp.Message) (Node_NestedNode, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Node_NestedNode{}, err
	}
	return Node_NestedNode{root.Struct()}, nil
}

func (s Node_NestedNode) Name() (string, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return "", err
	}

	return capnp.PtrToText(p), nil

}

func (s Node_NestedNode) NameBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	return capnp.PtrToData(p), nil
}

func (s Node_NestedNode) SetName(v string) error {

	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, t.List.ToPtr())
}

func (s Node_NestedNode) Id() uint64 {
	return s.Struct.Uint64(0)
}

func (s Node_NestedNode) SetId(v uint64) {

	s.Struct.SetUint64(0, v)
}

// Node_NestedNode_List is a list of Node_NestedNode.
type Node_NestedNode_List struct{ capnp.List }

// NewNode_NestedNode creates a new list of Node_NestedNode.
func NewNode_NestedNode_List(s *capnp.Segment, sz int32) (Node_NestedNode_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1}, sz)
	if err != nil {
		return Node_NestedNode_List{}, err
	}
	return Node_NestedNode_List{l}, nil
}

func (s Node_NestedNode_List) At(i int) Node_NestedNode { return Node_NestedNode{s.List.Struct(i)} }
func (s Node_NestedNode_List) Set(i int, v Node_NestedNode) error {
	return s.List.SetStruct(i, v.Struct)
}

type Field struct{ capnp.Struct }
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

func NewField(s *capnp.Segment) (Field, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 24, PointerCount: 4})
	if err != nil {
		return Field{}, err
	}
	return Field{st}, nil
}

func NewRootField(s *capnp.Segment) (Field, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 24, PointerCount: 4})
	if err != nil {
		return Field{}, err
	}
	return Field{st}, nil
}

func ReadRootField(msg *capnp.Message) (Field, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Field{}, err
	}
	return Field{root.Struct()}, nil
}

func (s Field) Which() Field_Which {
	return Field_Which(s.Struct.Uint16(8))
}

func (s Field) Name() (string, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return "", err
	}

	return capnp.PtrToText(p), nil

}

func (s Field) NameBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	return capnp.PtrToData(p), nil
}

func (s Field) SetName(v string) error {

	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, t.List.ToPtr())
}

func (s Field) CodeOrder() uint16 {
	return s.Struct.Uint16(0)
}

func (s Field) SetCodeOrder(v uint16) {

	s.Struct.SetUint16(0, v)
}

func (s Field) Annotations() (Annotation_List, error) {
	p, err := s.Struct.Ptr(1)
	if err != nil {
		return Annotation_List{}, err
	}

	return Annotation_List{List: p.List()}, nil

}

func (s Field) SetAnnotations(v Annotation_List) error {

	return s.Struct.SetPtr(1, v.List.ToPtr())
}

func (s Field) DiscriminantValue() uint16 {
	return s.Struct.Uint16(2) ^ 65535
}

func (s Field) SetDiscriminantValue(v uint16) {

	s.Struct.SetUint16(2, v^65535)
}
func (s Field) Slot() Field_slot { return Field_slot(s) }

func (s Field) SetSlot() { s.Struct.SetUint16(8, 0) }

func (s Field_slot) Offset() uint32 {
	return s.Struct.Uint32(4)
}

func (s Field_slot) SetOffset(v uint32) {

	s.Struct.SetUint32(4, v)
}

func (s Field_slot) Type() (Type, error) {
	p, err := s.Struct.Ptr(2)
	if err != nil {
		return Type{}, err
	}

	return Type{Struct: p.Struct()}, nil

}

func (s Field_slot) SetType(v Type) error {

	return s.Struct.SetPtr(2, v.Struct.ToPtr())
}

// NewType sets the type field to a newly
// allocated Type struct, preferring placement in s's segment.
func (s Field_slot) NewType() (Type, error) {

	ss, err := NewType(s.Struct.Segment())
	if err != nil {
		return Type{}, err
	}
	err = s.Struct.SetPtr(2, ss.Struct.ToPtr())
	return ss, err
}

func (s Field_slot) DefaultValue() (Value, error) {
	p, err := s.Struct.Ptr(3)
	if err != nil {
		return Value{}, err
	}

	return Value{Struct: p.Struct()}, nil

}

func (s Field_slot) SetDefaultValue(v Value) error {

	return s.Struct.SetPtr(3, v.Struct.ToPtr())
}

// NewDefaultValue sets the defaultValue field to a newly
// allocated Value struct, preferring placement in s's segment.
func (s Field_slot) NewDefaultValue() (Value, error) {

	ss, err := NewValue(s.Struct.Segment())
	if err != nil {
		return Value{}, err
	}
	err = s.Struct.SetPtr(3, ss.Struct.ToPtr())
	return ss, err
}

func (s Field_slot) HadExplicitDefault() bool {
	return s.Struct.Bit(128)
}

func (s Field_slot) SetHadExplicitDefault(v bool) {

	s.Struct.SetBit(128, v)
}
func (s Field) Group() Field_group { return Field_group(s) }

func (s Field) SetGroup() { s.Struct.SetUint16(8, 1) }

func (s Field_group) TypeId() uint64 {
	return s.Struct.Uint64(16)
}

func (s Field_group) SetTypeId(v uint64) {

	s.Struct.SetUint64(16, v)
}
func (s Field) Ordinal() Field_ordinal { return Field_ordinal(s) }

func (s Field_ordinal) Which() Field_ordinal_Which {
	return Field_ordinal_Which(s.Struct.Uint16(10))
}

func (s Field_ordinal) SetImplicit() {
	s.Struct.SetUint16(10, 0)
}

func (s Field_ordinal) Explicit() uint16 {
	return s.Struct.Uint16(12)
}

func (s Field_ordinal) SetExplicit(v uint16) {
	s.Struct.SetUint16(10, 1)
	s.Struct.SetUint16(12, v)
}

// Field_List is a list of Field.
type Field_List struct{ capnp.List }

// NewField creates a new list of Field.
func NewField_List(s *capnp.Segment, sz int32) (Field_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 24, PointerCount: 4}, sz)
	if err != nil {
		return Field_List{}, err
	}
	return Field_List{l}, nil
}

func (s Field_List) At(i int) Field           { return Field{s.List.Struct(i)} }
func (s Field_List) Set(i int, v Field) error { return s.List.SetStruct(i, v.Struct) }

type Enumerant struct{ capnp.Struct }

func NewEnumerant(s *capnp.Segment) (Enumerant, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2})
	if err != nil {
		return Enumerant{}, err
	}
	return Enumerant{st}, nil
}

func NewRootEnumerant(s *capnp.Segment) (Enumerant, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2})
	if err != nil {
		return Enumerant{}, err
	}
	return Enumerant{st}, nil
}

func ReadRootEnumerant(msg *capnp.Message) (Enumerant, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Enumerant{}, err
	}
	return Enumerant{root.Struct()}, nil
}

func (s Enumerant) Name() (string, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return "", err
	}

	return capnp.PtrToText(p), nil

}

func (s Enumerant) NameBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	return capnp.PtrToData(p), nil
}

func (s Enumerant) SetName(v string) error {

	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, t.List.ToPtr())
}

func (s Enumerant) CodeOrder() uint16 {
	return s.Struct.Uint16(0)
}

func (s Enumerant) SetCodeOrder(v uint16) {

	s.Struct.SetUint16(0, v)
}

func (s Enumerant) Annotations() (Annotation_List, error) {
	p, err := s.Struct.Ptr(1)
	if err != nil {
		return Annotation_List{}, err
	}

	return Annotation_List{List: p.List()}, nil

}

func (s Enumerant) SetAnnotations(v Annotation_List) error {

	return s.Struct.SetPtr(1, v.List.ToPtr())
}

// Enumerant_List is a list of Enumerant.
type Enumerant_List struct{ capnp.List }

// NewEnumerant creates a new list of Enumerant.
func NewEnumerant_List(s *capnp.Segment, sz int32) (Enumerant_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2}, sz)
	if err != nil {
		return Enumerant_List{}, err
	}
	return Enumerant_List{l}, nil
}

func (s Enumerant_List) At(i int) Enumerant           { return Enumerant{s.List.Struct(i)} }
func (s Enumerant_List) Set(i int, v Enumerant) error { return s.List.SetStruct(i, v.Struct) }

type Superclass struct{ capnp.Struct }

func NewSuperclass(s *capnp.Segment) (Superclass, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		return Superclass{}, err
	}
	return Superclass{st}, nil
}

func NewRootSuperclass(s *capnp.Segment) (Superclass, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		return Superclass{}, err
	}
	return Superclass{st}, nil
}

func ReadRootSuperclass(msg *capnp.Message) (Superclass, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Superclass{}, err
	}
	return Superclass{root.Struct()}, nil
}

func (s Superclass) Id() uint64 {
	return s.Struct.Uint64(0)
}

func (s Superclass) SetId(v uint64) {

	s.Struct.SetUint64(0, v)
}

func (s Superclass) Brand() (Brand, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return Brand{}, err
	}

	return Brand{Struct: p.Struct()}, nil

}

func (s Superclass) SetBrand(v Brand) error {

	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewBrand sets the brand field to a newly
// allocated Brand struct, preferring placement in s's segment.
func (s Superclass) NewBrand() (Brand, error) {

	ss, err := NewBrand(s.Struct.Segment())
	if err != nil {
		return Brand{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

// Superclass_List is a list of Superclass.
type Superclass_List struct{ capnp.List }

// NewSuperclass creates a new list of Superclass.
func NewSuperclass_List(s *capnp.Segment, sz int32) (Superclass_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1}, sz)
	if err != nil {
		return Superclass_List{}, err
	}
	return Superclass_List{l}, nil
}

func (s Superclass_List) At(i int) Superclass           { return Superclass{s.List.Struct(i)} }
func (s Superclass_List) Set(i int, v Superclass) error { return s.List.SetStruct(i, v.Struct) }

type Method struct{ capnp.Struct }

func NewMethod(s *capnp.Segment) (Method, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 24, PointerCount: 5})
	if err != nil {
		return Method{}, err
	}
	return Method{st}, nil
}

func NewRootMethod(s *capnp.Segment) (Method, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 24, PointerCount: 5})
	if err != nil {
		return Method{}, err
	}
	return Method{st}, nil
}

func ReadRootMethod(msg *capnp.Message) (Method, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Method{}, err
	}
	return Method{root.Struct()}, nil
}

func (s Method) Name() (string, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return "", err
	}

	return capnp.PtrToText(p), nil

}

func (s Method) NameBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	return capnp.PtrToData(p), nil
}

func (s Method) SetName(v string) error {

	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, t.List.ToPtr())
}

func (s Method) CodeOrder() uint16 {
	return s.Struct.Uint16(0)
}

func (s Method) SetCodeOrder(v uint16) {

	s.Struct.SetUint16(0, v)
}

func (s Method) ImplicitParameters() (Node_Parameter_List, error) {
	p, err := s.Struct.Ptr(4)
	if err != nil {
		return Node_Parameter_List{}, err
	}

	return Node_Parameter_List{List: p.List()}, nil

}

func (s Method) SetImplicitParameters(v Node_Parameter_List) error {

	return s.Struct.SetPtr(4, v.List.ToPtr())
}

func (s Method) ParamStructType() uint64 {
	return s.Struct.Uint64(8)
}

func (s Method) SetParamStructType(v uint64) {

	s.Struct.SetUint64(8, v)
}

func (s Method) ParamBrand() (Brand, error) {
	p, err := s.Struct.Ptr(2)
	if err != nil {
		return Brand{}, err
	}

	return Brand{Struct: p.Struct()}, nil

}

func (s Method) SetParamBrand(v Brand) error {

	return s.Struct.SetPtr(2, v.Struct.ToPtr())
}

// NewParamBrand sets the paramBrand field to a newly
// allocated Brand struct, preferring placement in s's segment.
func (s Method) NewParamBrand() (Brand, error) {

	ss, err := NewBrand(s.Struct.Segment())
	if err != nil {
		return Brand{}, err
	}
	err = s.Struct.SetPtr(2, ss.Struct.ToPtr())
	return ss, err
}

func (s Method) ResultStructType() uint64 {
	return s.Struct.Uint64(16)
}

func (s Method) SetResultStructType(v uint64) {

	s.Struct.SetUint64(16, v)
}

func (s Method) ResultBrand() (Brand, error) {
	p, err := s.Struct.Ptr(3)
	if err != nil {
		return Brand{}, err
	}

	return Brand{Struct: p.Struct()}, nil

}

func (s Method) SetResultBrand(v Brand) error {

	return s.Struct.SetPtr(3, v.Struct.ToPtr())
}

// NewResultBrand sets the resultBrand field to a newly
// allocated Brand struct, preferring placement in s's segment.
func (s Method) NewResultBrand() (Brand, error) {

	ss, err := NewBrand(s.Struct.Segment())
	if err != nil {
		return Brand{}, err
	}
	err = s.Struct.SetPtr(3, ss.Struct.ToPtr())
	return ss, err
}

func (s Method) Annotations() (Annotation_List, error) {
	p, err := s.Struct.Ptr(1)
	if err != nil {
		return Annotation_List{}, err
	}

	return Annotation_List{List: p.List()}, nil

}

func (s Method) SetAnnotations(v Annotation_List) error {

	return s.Struct.SetPtr(1, v.List.ToPtr())
}

// Method_List is a list of Method.
type Method_List struct{ capnp.List }

// NewMethod creates a new list of Method.
func NewMethod_List(s *capnp.Segment, sz int32) (Method_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 24, PointerCount: 5}, sz)
	if err != nil {
		return Method_List{}, err
	}
	return Method_List{l}, nil
}

func (s Method_List) At(i int) Method           { return Method{s.List.Struct(i)} }
func (s Method_List) Set(i int, v Method) error { return s.List.SetStruct(i, v.Struct) }

type Type struct{ capnp.Struct }
type Type_list Type
type Type_enum Type
type Type_structGroup Type
type Type_interface Type
type Type_anyPointer Type
type Type_anyPointer_parameter Type
type Type_anyPointer_implicitMethodParameter Type
type Type_Which uint16

const (
	Type_Which_void        Type_Which = 0
	Type_Which_bool        Type_Which = 1
	Type_Which_int8        Type_Which = 2
	Type_Which_int16       Type_Which = 3
	Type_Which_int32       Type_Which = 4
	Type_Which_int64       Type_Which = 5
	Type_Which_uint8       Type_Which = 6
	Type_Which_uint16      Type_Which = 7
	Type_Which_uint32      Type_Which = 8
	Type_Which_uint64      Type_Which = 9
	Type_Which_float32     Type_Which = 10
	Type_Which_float64     Type_Which = 11
	Type_Which_text        Type_Which = 12
	Type_Which_data        Type_Which = 13
	Type_Which_list        Type_Which = 14
	Type_Which_enum        Type_Which = 15
	Type_Which_structGroup Type_Which = 16
	Type_Which_interface   Type_Which = 17
	Type_Which_anyPointer  Type_Which = 18
)

func (w Type_Which) String() string {
	const s = "voidboolint8int16int32int64uint8uint16uint32uint64float32float64textdatalistenumstructGroupinterfaceanyPointer"
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
	case Type_Which_structGroup:
		return s[80:91]
	case Type_Which_interface:
		return s[91:100]
	case Type_Which_anyPointer:
		return s[100:110]

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

func NewType(s *capnp.Segment) (Type, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 24, PointerCount: 1})
	if err != nil {
		return Type{}, err
	}
	return Type{st}, nil
}

func NewRootType(s *capnp.Segment) (Type, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 24, PointerCount: 1})
	if err != nil {
		return Type{}, err
	}
	return Type{st}, nil
}

func ReadRootType(msg *capnp.Message) (Type, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Type{}, err
	}
	return Type{root.Struct()}, nil
}

func (s Type) Which() Type_Which {
	return Type_Which(s.Struct.Uint16(0))
}

func (s Type) SetVoid() {
	s.Struct.SetUint16(0, 0)
}

func (s Type) SetBool() {
	s.Struct.SetUint16(0, 1)
}

func (s Type) SetInt8() {
	s.Struct.SetUint16(0, 2)
}

func (s Type) SetInt16() {
	s.Struct.SetUint16(0, 3)
}

func (s Type) SetInt32() {
	s.Struct.SetUint16(0, 4)
}

func (s Type) SetInt64() {
	s.Struct.SetUint16(0, 5)
}

func (s Type) SetUint8() {
	s.Struct.SetUint16(0, 6)
}

func (s Type) SetUint16() {
	s.Struct.SetUint16(0, 7)
}

func (s Type) SetUint32() {
	s.Struct.SetUint16(0, 8)
}

func (s Type) SetUint64() {
	s.Struct.SetUint16(0, 9)
}

func (s Type) SetFloat32() {
	s.Struct.SetUint16(0, 10)
}

func (s Type) SetFloat64() {
	s.Struct.SetUint16(0, 11)
}

func (s Type) SetText() {
	s.Struct.SetUint16(0, 12)
}

func (s Type) SetData() {
	s.Struct.SetUint16(0, 13)
}
func (s Type) List() Type_list { return Type_list(s) }

func (s Type) SetList() { s.Struct.SetUint16(0, 14) }

func (s Type_list) ElementType() (Type, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return Type{}, err
	}

	return Type{Struct: p.Struct()}, nil

}

func (s Type_list) SetElementType(v Type) error {

	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewElementType sets the elementType field to a newly
// allocated Type struct, preferring placement in s's segment.
func (s Type_list) NewElementType() (Type, error) {

	ss, err := NewType(s.Struct.Segment())
	if err != nil {
		return Type{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}
func (s Type) Enum() Type_enum { return Type_enum(s) }

func (s Type) SetEnum() { s.Struct.SetUint16(0, 15) }

func (s Type_enum) TypeId() uint64 {
	return s.Struct.Uint64(8)
}

func (s Type_enum) SetTypeId(v uint64) {

	s.Struct.SetUint64(8, v)
}

func (s Type_enum) Brand() (Brand, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return Brand{}, err
	}

	return Brand{Struct: p.Struct()}, nil

}

func (s Type_enum) SetBrand(v Brand) error {

	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewBrand sets the brand field to a newly
// allocated Brand struct, preferring placement in s's segment.
func (s Type_enum) NewBrand() (Brand, error) {

	ss, err := NewBrand(s.Struct.Segment())
	if err != nil {
		return Brand{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}
func (s Type) StructGroup() Type_structGroup { return Type_structGroup(s) }

func (s Type) SetStructGroup() { s.Struct.SetUint16(0, 16) }

func (s Type_structGroup) TypeId() uint64 {
	return s.Struct.Uint64(8)
}

func (s Type_structGroup) SetTypeId(v uint64) {

	s.Struct.SetUint64(8, v)
}

func (s Type_structGroup) Brand() (Brand, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return Brand{}, err
	}

	return Brand{Struct: p.Struct()}, nil

}

func (s Type_structGroup) SetBrand(v Brand) error {

	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewBrand sets the brand field to a newly
// allocated Brand struct, preferring placement in s's segment.
func (s Type_structGroup) NewBrand() (Brand, error) {

	ss, err := NewBrand(s.Struct.Segment())
	if err != nil {
		return Brand{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}
func (s Type) Interface() Type_interface { return Type_interface(s) }

func (s Type) SetInterface() { s.Struct.SetUint16(0, 17) }

func (s Type_interface) TypeId() uint64 {
	return s.Struct.Uint64(8)
}

func (s Type_interface) SetTypeId(v uint64) {

	s.Struct.SetUint64(8, v)
}

func (s Type_interface) Brand() (Brand, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return Brand{}, err
	}

	return Brand{Struct: p.Struct()}, nil

}

func (s Type_interface) SetBrand(v Brand) error {

	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewBrand sets the brand field to a newly
// allocated Brand struct, preferring placement in s's segment.
func (s Type_interface) NewBrand() (Brand, error) {

	ss, err := NewBrand(s.Struct.Segment())
	if err != nil {
		return Brand{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}
func (s Type) AnyPointer() Type_anyPointer { return Type_anyPointer(s) }

func (s Type) SetAnyPointer() { s.Struct.SetUint16(0, 18) }

func (s Type_anyPointer) Which() Type_anyPointer_Which {
	return Type_anyPointer_Which(s.Struct.Uint16(8))
}

func (s Type_anyPointer) SetUnconstrained() {
	s.Struct.SetUint16(8, 0)
}
func (s Type_anyPointer) Parameter() Type_anyPointer_parameter { return Type_anyPointer_parameter(s) }

func (s Type_anyPointer) SetParameter() { s.Struct.SetUint16(8, 1) }

func (s Type_anyPointer_parameter) ScopeId() uint64 {
	return s.Struct.Uint64(16)
}

func (s Type_anyPointer_parameter) SetScopeId(v uint64) {

	s.Struct.SetUint64(16, v)
}

func (s Type_anyPointer_parameter) ParameterIndex() uint16 {
	return s.Struct.Uint16(10)
}

func (s Type_anyPointer_parameter) SetParameterIndex(v uint16) {

	s.Struct.SetUint16(10, v)
}
func (s Type_anyPointer) ImplicitMethodParameter() Type_anyPointer_implicitMethodParameter {
	return Type_anyPointer_implicitMethodParameter(s)
}

func (s Type_anyPointer) SetImplicitMethodParameter() { s.Struct.SetUint16(8, 2) }

func (s Type_anyPointer_implicitMethodParameter) ParameterIndex() uint16 {
	return s.Struct.Uint16(10)
}

func (s Type_anyPointer_implicitMethodParameter) SetParameterIndex(v uint16) {

	s.Struct.SetUint16(10, v)
}

// Type_List is a list of Type.
type Type_List struct{ capnp.List }

// NewType creates a new list of Type.
func NewType_List(s *capnp.Segment, sz int32) (Type_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 24, PointerCount: 1}, sz)
	if err != nil {
		return Type_List{}, err
	}
	return Type_List{l}, nil
}

func (s Type_List) At(i int) Type           { return Type{s.List.Struct(i)} }
func (s Type_List) Set(i int, v Type) error { return s.List.SetStruct(i, v.Struct) }

type Brand struct{ capnp.Struct }

func NewBrand(s *capnp.Segment) (Brand, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Brand{}, err
	}
	return Brand{st}, nil
}

func NewRootBrand(s *capnp.Segment) (Brand, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	if err != nil {
		return Brand{}, err
	}
	return Brand{st}, nil
}

func ReadRootBrand(msg *capnp.Message) (Brand, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Brand{}, err
	}
	return Brand{root.Struct()}, nil
}

func (s Brand) Scopes() (Brand_Scope_List, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return Brand_Scope_List{}, err
	}

	return Brand_Scope_List{List: p.List()}, nil

}

func (s Brand) SetScopes(v Brand_Scope_List) error {

	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// Brand_List is a list of Brand.
type Brand_List struct{ capnp.List }

// NewBrand creates a new list of Brand.
func NewBrand_List(s *capnp.Segment, sz int32) (Brand_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	if err != nil {
		return Brand_List{}, err
	}
	return Brand_List{l}, nil
}

func (s Brand_List) At(i int) Brand           { return Brand{s.List.Struct(i)} }
func (s Brand_List) Set(i int, v Brand) error { return s.List.SetStruct(i, v.Struct) }

type Brand_Scope struct{ capnp.Struct }
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

func NewBrand_Scope(s *capnp.Segment) (Brand_Scope, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 1})
	if err != nil {
		return Brand_Scope{}, err
	}
	return Brand_Scope{st}, nil
}

func NewRootBrand_Scope(s *capnp.Segment) (Brand_Scope, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 1})
	if err != nil {
		return Brand_Scope{}, err
	}
	return Brand_Scope{st}, nil
}

func ReadRootBrand_Scope(msg *capnp.Message) (Brand_Scope, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Brand_Scope{}, err
	}
	return Brand_Scope{root.Struct()}, nil
}

func (s Brand_Scope) Which() Brand_Scope_Which {
	return Brand_Scope_Which(s.Struct.Uint16(8))
}

func (s Brand_Scope) ScopeId() uint64 {
	return s.Struct.Uint64(0)
}

func (s Brand_Scope) SetScopeId(v uint64) {

	s.Struct.SetUint64(0, v)
}

func (s Brand_Scope) Bind() (Brand_Binding_List, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return Brand_Binding_List{}, err
	}

	return Brand_Binding_List{List: p.List()}, nil

}

func (s Brand_Scope) SetBind(v Brand_Binding_List) error {
	s.Struct.SetUint16(8, 0)
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

func (s Brand_Scope) SetInherit() {
	s.Struct.SetUint16(8, 1)
}

// Brand_Scope_List is a list of Brand_Scope.
type Brand_Scope_List struct{ capnp.List }

// NewBrand_Scope creates a new list of Brand_Scope.
func NewBrand_Scope_List(s *capnp.Segment, sz int32) (Brand_Scope_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 16, PointerCount: 1}, sz)
	if err != nil {
		return Brand_Scope_List{}, err
	}
	return Brand_Scope_List{l}, nil
}

func (s Brand_Scope_List) At(i int) Brand_Scope           { return Brand_Scope{s.List.Struct(i)} }
func (s Brand_Scope_List) Set(i int, v Brand_Scope) error { return s.List.SetStruct(i, v.Struct) }

type Brand_Binding struct{ capnp.Struct }
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

func NewBrand_Binding(s *capnp.Segment) (Brand_Binding, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		return Brand_Binding{}, err
	}
	return Brand_Binding{st}, nil
}

func NewRootBrand_Binding(s *capnp.Segment) (Brand_Binding, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		return Brand_Binding{}, err
	}
	return Brand_Binding{st}, nil
}

func ReadRootBrand_Binding(msg *capnp.Message) (Brand_Binding, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Brand_Binding{}, err
	}
	return Brand_Binding{root.Struct()}, nil
}

func (s Brand_Binding) Which() Brand_Binding_Which {
	return Brand_Binding_Which(s.Struct.Uint16(0))
}

func (s Brand_Binding) SetUnbound() {
	s.Struct.SetUint16(0, 0)
}

func (s Brand_Binding) Type() (Type, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return Type{}, err
	}

	return Type{Struct: p.Struct()}, nil

}

func (s Brand_Binding) SetType(v Type) error {
	s.Struct.SetUint16(0, 1)
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewType sets the type field to a newly
// allocated Type struct, preferring placement in s's segment.
func (s Brand_Binding) NewType() (Type, error) {
	s.Struct.SetUint16(0, 1)
	ss, err := NewType(s.Struct.Segment())
	if err != nil {
		return Type{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

// Brand_Binding_List is a list of Brand_Binding.
type Brand_Binding_List struct{ capnp.List }

// NewBrand_Binding creates a new list of Brand_Binding.
func NewBrand_Binding_List(s *capnp.Segment, sz int32) (Brand_Binding_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1}, sz)
	if err != nil {
		return Brand_Binding_List{}, err
	}
	return Brand_Binding_List{l}, nil
}

func (s Brand_Binding_List) At(i int) Brand_Binding           { return Brand_Binding{s.List.Struct(i)} }
func (s Brand_Binding_List) Set(i int, v Brand_Binding) error { return s.List.SetStruct(i, v.Struct) }

type Value struct{ capnp.Struct }
type Value_Which uint16

const (
	Value_Which_void        Value_Which = 0
	Value_Which_bool        Value_Which = 1
	Value_Which_int8        Value_Which = 2
	Value_Which_int16       Value_Which = 3
	Value_Which_int32       Value_Which = 4
	Value_Which_int64       Value_Which = 5
	Value_Which_uint8       Value_Which = 6
	Value_Which_uint16      Value_Which = 7
	Value_Which_uint32      Value_Which = 8
	Value_Which_uint64      Value_Which = 9
	Value_Which_float32     Value_Which = 10
	Value_Which_float64     Value_Which = 11
	Value_Which_text        Value_Which = 12
	Value_Which_data        Value_Which = 13
	Value_Which_list        Value_Which = 14
	Value_Which_enum        Value_Which = 15
	Value_Which_structField Value_Which = 16
	Value_Which_interface   Value_Which = 17
	Value_Which_anyPointer  Value_Which = 18
)

func (w Value_Which) String() string {
	const s = "voidboolint8int16int32int64uint8uint16uint32uint64float32float64textdatalistenumstructFieldinterfaceanyPointer"
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
	case Value_Which_structField:
		return s[80:91]
	case Value_Which_interface:
		return s[91:100]
	case Value_Which_anyPointer:
		return s[100:110]

	}
	return "Value_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
}

func NewValue(s *capnp.Segment) (Value, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 1})
	if err != nil {
		return Value{}, err
	}
	return Value{st}, nil
}

func NewRootValue(s *capnp.Segment) (Value, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 1})
	if err != nil {
		return Value{}, err
	}
	return Value{st}, nil
}

func ReadRootValue(msg *capnp.Message) (Value, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Value{}, err
	}
	return Value{root.Struct()}, nil
}

func (s Value) Which() Value_Which {
	return Value_Which(s.Struct.Uint16(0))
}

func (s Value) SetVoid() {
	s.Struct.SetUint16(0, 0)
}

func (s Value) Bool() bool {
	return s.Struct.Bit(16)
}

func (s Value) SetBool(v bool) {
	s.Struct.SetUint16(0, 1)
	s.Struct.SetBit(16, v)
}

func (s Value) Int8() int8 {
	return int8(s.Struct.Uint8(2))
}

func (s Value) SetInt8(v int8) {
	s.Struct.SetUint16(0, 2)
	s.Struct.SetUint8(2, uint8(v))
}

func (s Value) Int16() int16 {
	return int16(s.Struct.Uint16(2))
}

func (s Value) SetInt16(v int16) {
	s.Struct.SetUint16(0, 3)
	s.Struct.SetUint16(2, uint16(v))
}

func (s Value) Int32() int32 {
	return int32(s.Struct.Uint32(4))
}

func (s Value) SetInt32(v int32) {
	s.Struct.SetUint16(0, 4)
	s.Struct.SetUint32(4, uint32(v))
}

func (s Value) Int64() int64 {
	return int64(s.Struct.Uint64(8))
}

func (s Value) SetInt64(v int64) {
	s.Struct.SetUint16(0, 5)
	s.Struct.SetUint64(8, uint64(v))
}

func (s Value) Uint8() uint8 {
	return s.Struct.Uint8(2)
}

func (s Value) SetUint8(v uint8) {
	s.Struct.SetUint16(0, 6)
	s.Struct.SetUint8(2, v)
}

func (s Value) Uint16() uint16 {
	return s.Struct.Uint16(2)
}

func (s Value) SetUint16(v uint16) {
	s.Struct.SetUint16(0, 7)
	s.Struct.SetUint16(2, v)
}

func (s Value) Uint32() uint32 {
	return s.Struct.Uint32(4)
}

func (s Value) SetUint32(v uint32) {
	s.Struct.SetUint16(0, 8)
	s.Struct.SetUint32(4, v)
}

func (s Value) Uint64() uint64 {
	return s.Struct.Uint64(8)
}

func (s Value) SetUint64(v uint64) {
	s.Struct.SetUint16(0, 9)
	s.Struct.SetUint64(8, v)
}

func (s Value) Float32() float32 {
	return math.Float32frombits(s.Struct.Uint32(4))
}

func (s Value) SetFloat32(v float32) {
	s.Struct.SetUint16(0, 10)
	s.Struct.SetUint32(4, math.Float32bits(v))
}

func (s Value) Float64() float64 {
	return math.Float64frombits(s.Struct.Uint64(8))
}

func (s Value) SetFloat64(v float64) {
	s.Struct.SetUint16(0, 11)
	s.Struct.SetUint64(8, math.Float64bits(v))
}

func (s Value) Text() (string, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return "", err
	}

	return capnp.PtrToText(p), nil

}

func (s Value) TextBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	return capnp.PtrToData(p), nil
}

func (s Value) SetText(v string) error {
	s.Struct.SetUint16(0, 12)
	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, t.List.ToPtr())
}

func (s Value) Data() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}

	return []byte(capnp.PtrToData(p)), nil

}

func (s Value) SetData(v []byte) error {
	s.Struct.SetUint16(0, 13)
	d, err := capnp.NewData(s.Struct.Segment(), []byte(v))
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, d.List.ToPtr())
}

func (s Value) List() (capnp.Pointer, error) {

	return s.Struct.Pointer(0)

}

func (s Value) SetList(v capnp.Pointer) error {
	s.Struct.SetUint16(0, 14)
	return s.Struct.SetPointer(0, v)
}

func (s Value) Enum() uint16 {
	return s.Struct.Uint16(2)
}

func (s Value) SetEnum(v uint16) {
	s.Struct.SetUint16(0, 15)
	s.Struct.SetUint16(2, v)
}

func (s Value) StructField() (capnp.Pointer, error) {

	return s.Struct.Pointer(0)

}

func (s Value) SetStructField(v capnp.Pointer) error {
	s.Struct.SetUint16(0, 16)
	return s.Struct.SetPointer(0, v)
}

func (s Value) SetInterface() {
	s.Struct.SetUint16(0, 17)
}

func (s Value) AnyPointer() (capnp.Pointer, error) {

	return s.Struct.Pointer(0)

}

func (s Value) SetAnyPointer(v capnp.Pointer) error {
	s.Struct.SetUint16(0, 18)
	return s.Struct.SetPointer(0, v)
}

// Value_List is a list of Value.
type Value_List struct{ capnp.List }

// NewValue creates a new list of Value.
func NewValue_List(s *capnp.Segment, sz int32) (Value_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 16, PointerCount: 1}, sz)
	if err != nil {
		return Value_List{}, err
	}
	return Value_List{l}, nil
}

func (s Value_List) At(i int) Value           { return Value{s.List.Struct(i)} }
func (s Value_List) Set(i int, v Value) error { return s.List.SetStruct(i, v.Struct) }

type Annotation struct{ capnp.Struct }

func NewAnnotation(s *capnp.Segment) (Annotation, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2})
	if err != nil {
		return Annotation{}, err
	}
	return Annotation{st}, nil
}

func NewRootAnnotation(s *capnp.Segment) (Annotation, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2})
	if err != nil {
		return Annotation{}, err
	}
	return Annotation{st}, nil
}

func ReadRootAnnotation(msg *capnp.Message) (Annotation, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return Annotation{}, err
	}
	return Annotation{root.Struct()}, nil
}

func (s Annotation) Id() uint64 {
	return s.Struct.Uint64(0)
}

func (s Annotation) SetId(v uint64) {

	s.Struct.SetUint64(0, v)
}

func (s Annotation) Brand() (Brand, error) {
	p, err := s.Struct.Ptr(1)
	if err != nil {
		return Brand{}, err
	}

	return Brand{Struct: p.Struct()}, nil

}

func (s Annotation) SetBrand(v Brand) error {

	return s.Struct.SetPtr(1, v.Struct.ToPtr())
}

// NewBrand sets the brand field to a newly
// allocated Brand struct, preferring placement in s's segment.
func (s Annotation) NewBrand() (Brand, error) {

	ss, err := NewBrand(s.Struct.Segment())
	if err != nil {
		return Brand{}, err
	}
	err = s.Struct.SetPtr(1, ss.Struct.ToPtr())
	return ss, err
}

func (s Annotation) Value() (Value, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return Value{}, err
	}

	return Value{Struct: p.Struct()}, nil

}

func (s Annotation) SetValue(v Value) error {

	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewValue sets the value field to a newly
// allocated Value struct, preferring placement in s's segment.
func (s Annotation) NewValue() (Value, error) {

	ss, err := NewValue(s.Struct.Segment())
	if err != nil {
		return Value{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

// Annotation_List is a list of Annotation.
type Annotation_List struct{ capnp.List }

// NewAnnotation creates a new list of Annotation.
func NewAnnotation_List(s *capnp.Segment, sz int32) (Annotation_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2}, sz)
	if err != nil {
		return Annotation_List{}, err
	}
	return Annotation_List{l}, nil
}

func (s Annotation_List) At(i int) Annotation           { return Annotation{s.List.Struct(i)} }
func (s Annotation_List) Set(i int, v Annotation) error { return s.List.SetStruct(i, v.Struct) }

type ElementSize uint16

// Values of ElementSize.
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

// String returns the enum's constant name.
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

// ElementSizeFromString returns the enum value with a name,
// or the zero value if there's no such value.
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

type ElementSize_List struct{ capnp.List }

func NewElementSize_List(s *capnp.Segment, sz int32) (ElementSize_List, error) {
	l, err := capnp.NewUInt16List(s, sz)
	if err != nil {
		return ElementSize_List{}, err
	}
	return ElementSize_List{l.List}, nil
}

func (l ElementSize_List) At(i int) ElementSize {
	ul := capnp.UInt16List{List: l.List}
	return ElementSize(ul.At(i))
}

func (l ElementSize_List) Set(i int, v ElementSize) {
	ul := capnp.UInt16List{List: l.List}
	ul.Set(i, uint16(v))
}

type CodeGeneratorRequest struct{ capnp.Struct }

func NewCodeGeneratorRequest(s *capnp.Segment) (CodeGeneratorRequest, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	if err != nil {
		return CodeGeneratorRequest{}, err
	}
	return CodeGeneratorRequest{st}, nil
}

func NewRootCodeGeneratorRequest(s *capnp.Segment) (CodeGeneratorRequest, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2})
	if err != nil {
		return CodeGeneratorRequest{}, err
	}
	return CodeGeneratorRequest{st}, nil
}

func ReadRootCodeGeneratorRequest(msg *capnp.Message) (CodeGeneratorRequest, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return CodeGeneratorRequest{}, err
	}
	return CodeGeneratorRequest{root.Struct()}, nil
}

func (s CodeGeneratorRequest) Nodes() (Node_List, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return Node_List{}, err
	}

	return Node_List{List: p.List()}, nil

}

func (s CodeGeneratorRequest) SetNodes(v Node_List) error {

	return s.Struct.SetPtr(0, v.List.ToPtr())
}

func (s CodeGeneratorRequest) RequestedFiles() (CodeGeneratorRequest_RequestedFile_List, error) {
	p, err := s.Struct.Ptr(1)
	if err != nil {
		return CodeGeneratorRequest_RequestedFile_List{}, err
	}

	return CodeGeneratorRequest_RequestedFile_List{List: p.List()}, nil

}

func (s CodeGeneratorRequest) SetRequestedFiles(v CodeGeneratorRequest_RequestedFile_List) error {

	return s.Struct.SetPtr(1, v.List.ToPtr())
}

// CodeGeneratorRequest_List is a list of CodeGeneratorRequest.
type CodeGeneratorRequest_List struct{ capnp.List }

// NewCodeGeneratorRequest creates a new list of CodeGeneratorRequest.
func NewCodeGeneratorRequest_List(s *capnp.Segment, sz int32) (CodeGeneratorRequest_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 2}, sz)
	if err != nil {
		return CodeGeneratorRequest_List{}, err
	}
	return CodeGeneratorRequest_List{l}, nil
}

func (s CodeGeneratorRequest_List) At(i int) CodeGeneratorRequest {
	return CodeGeneratorRequest{s.List.Struct(i)}
}
func (s CodeGeneratorRequest_List) Set(i int, v CodeGeneratorRequest) error {
	return s.List.SetStruct(i, v.Struct)
}

type CodeGeneratorRequest_RequestedFile struct{ capnp.Struct }

func NewCodeGeneratorRequest_RequestedFile(s *capnp.Segment) (CodeGeneratorRequest_RequestedFile, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2})
	if err != nil {
		return CodeGeneratorRequest_RequestedFile{}, err
	}
	return CodeGeneratorRequest_RequestedFile{st}, nil
}

func NewRootCodeGeneratorRequest_RequestedFile(s *capnp.Segment) (CodeGeneratorRequest_RequestedFile, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2})
	if err != nil {
		return CodeGeneratorRequest_RequestedFile{}, err
	}
	return CodeGeneratorRequest_RequestedFile{st}, nil
}

func ReadRootCodeGeneratorRequest_RequestedFile(msg *capnp.Message) (CodeGeneratorRequest_RequestedFile, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return CodeGeneratorRequest_RequestedFile{}, err
	}
	return CodeGeneratorRequest_RequestedFile{root.Struct()}, nil
}

func (s CodeGeneratorRequest_RequestedFile) Id() uint64 {
	return s.Struct.Uint64(0)
}

func (s CodeGeneratorRequest_RequestedFile) SetId(v uint64) {

	s.Struct.SetUint64(0, v)
}

func (s CodeGeneratorRequest_RequestedFile) Filename() (string, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return "", err
	}

	return capnp.PtrToText(p), nil

}

func (s CodeGeneratorRequest_RequestedFile) FilenameBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	return capnp.PtrToData(p), nil
}

func (s CodeGeneratorRequest_RequestedFile) SetFilename(v string) error {

	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, t.List.ToPtr())
}

func (s CodeGeneratorRequest_RequestedFile) Imports() (CodeGeneratorRequest_RequestedFile_Import_List, error) {
	p, err := s.Struct.Ptr(1)
	if err != nil {
		return CodeGeneratorRequest_RequestedFile_Import_List{}, err
	}

	return CodeGeneratorRequest_RequestedFile_Import_List{List: p.List()}, nil

}

func (s CodeGeneratorRequest_RequestedFile) SetImports(v CodeGeneratorRequest_RequestedFile_Import_List) error {

	return s.Struct.SetPtr(1, v.List.ToPtr())
}

// CodeGeneratorRequest_RequestedFile_List is a list of CodeGeneratorRequest_RequestedFile.
type CodeGeneratorRequest_RequestedFile_List struct{ capnp.List }

// NewCodeGeneratorRequest_RequestedFile creates a new list of CodeGeneratorRequest_RequestedFile.
func NewCodeGeneratorRequest_RequestedFile_List(s *capnp.Segment, sz int32) (CodeGeneratorRequest_RequestedFile_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 2}, sz)
	if err != nil {
		return CodeGeneratorRequest_RequestedFile_List{}, err
	}
	return CodeGeneratorRequest_RequestedFile_List{l}, nil
}

func (s CodeGeneratorRequest_RequestedFile_List) At(i int) CodeGeneratorRequest_RequestedFile {
	return CodeGeneratorRequest_RequestedFile{s.List.Struct(i)}
}
func (s CodeGeneratorRequest_RequestedFile_List) Set(i int, v CodeGeneratorRequest_RequestedFile) error {
	return s.List.SetStruct(i, v.Struct)
}

type CodeGeneratorRequest_RequestedFile_Import struct{ capnp.Struct }

func NewCodeGeneratorRequest_RequestedFile_Import(s *capnp.Segment) (CodeGeneratorRequest_RequestedFile_Import, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		return CodeGeneratorRequest_RequestedFile_Import{}, err
	}
	return CodeGeneratorRequest_RequestedFile_Import{st}, nil
}

func NewRootCodeGeneratorRequest_RequestedFile_Import(s *capnp.Segment) (CodeGeneratorRequest_RequestedFile_Import, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1})
	if err != nil {
		return CodeGeneratorRequest_RequestedFile_Import{}, err
	}
	return CodeGeneratorRequest_RequestedFile_Import{st}, nil
}

func ReadRootCodeGeneratorRequest_RequestedFile_Import(msg *capnp.Message) (CodeGeneratorRequest_RequestedFile_Import, error) {
	root, err := msg.RootPtr()
	if err != nil {
		return CodeGeneratorRequest_RequestedFile_Import{}, err
	}
	return CodeGeneratorRequest_RequestedFile_Import{root.Struct()}, nil
}

func (s CodeGeneratorRequest_RequestedFile_Import) Id() uint64 {
	return s.Struct.Uint64(0)
}

func (s CodeGeneratorRequest_RequestedFile_Import) SetId(v uint64) {

	s.Struct.SetUint64(0, v)
}

func (s CodeGeneratorRequest_RequestedFile_Import) Name() (string, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return "", err
	}

	return capnp.PtrToText(p), nil

}

func (s CodeGeneratorRequest_RequestedFile_Import) NameBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	return capnp.PtrToData(p), nil
}

func (s CodeGeneratorRequest_RequestedFile_Import) SetName(v string) error {

	t, err := capnp.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPtr(0, t.List.ToPtr())
}

// CodeGeneratorRequest_RequestedFile_Import_List is a list of CodeGeneratorRequest_RequestedFile_Import.
type CodeGeneratorRequest_RequestedFile_Import_List struct{ capnp.List }

// NewCodeGeneratorRequest_RequestedFile_Import creates a new list of CodeGeneratorRequest_RequestedFile_Import.
func NewCodeGeneratorRequest_RequestedFile_Import_List(s *capnp.Segment, sz int32) (CodeGeneratorRequest_RequestedFile_Import_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 8, PointerCount: 1}, sz)
	if err != nil {
		return CodeGeneratorRequest_RequestedFile_Import_List{}, err
	}
	return CodeGeneratorRequest_RequestedFile_Import_List{l}, nil
}

func (s CodeGeneratorRequest_RequestedFile_Import_List) At(i int) CodeGeneratorRequest_RequestedFile_Import {
	return CodeGeneratorRequest_RequestedFile_Import{s.List.Struct(i)}
}
func (s CodeGeneratorRequest_RequestedFile_Import_List) Set(i int, v CodeGeneratorRequest_RequestedFile_Import) error {
	return s.List.SetStruct(i, v.Struct)
}
