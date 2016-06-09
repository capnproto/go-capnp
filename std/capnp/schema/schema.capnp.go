package schema

// AUTO GENERATED - DO NOT EDIT

import (
	math "math"
	strconv "strconv"
	capnp "zombiezen.com/go/capnproto2"
)

// Constants defined in schema.capnp.
const (
	Field_noDiscriminant = uint16(65535)
)

type Node struct{ capnp.Struct }
type Node_structNode Node
type Node_enum Node
type Node_interface Node
type Node_const Node
type Node_annotation Node
type Node_Which uint16

const (
	Node_Which_file       Node_Which = 0
	Node_Which_structNode Node_Which = 1
	Node_Which_enum       Node_Which = 2
	Node_Which_interface  Node_Which = 3
	Node_Which_const      Node_Which = 4
	Node_Which_annotation Node_Which = 5
)

func (w Node_Which) String() string {
	const s = "filestructNodeenuminterfaceconstannotation"
	switch w {
	case Node_Which_file:
		return s[0:4]
	case Node_Which_structNode:
		return s[4:14]
	case Node_Which_enum:
		return s[14:18]
	case Node_Which_interface:
		return s[18:27]
	case Node_Which_const:
		return s[27:32]
	case Node_Which_annotation:
		return s[32:42]

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
	return p.Text(), nil
}

func (s Node) HasDisplayName() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Node) DisplayNameBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	d := p.Data()
	if len(d) == 0 {
		return d, nil
	}
	return d[:len(d)-1], nil
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

func (s Node) HasParameters() bool {
	p, err := s.Struct.Ptr(5)
	return p.IsValid() || err != nil
}

func (s Node) SetParameters(v Node_Parameter_List) error {
	return s.Struct.SetPtr(5, v.List.ToPtr())
}

// NewParameters sets the parameters field to a newly
// allocated Node_Parameter_List, preferring placement in s's segment.
func (s Node) NewParameters(n int32) (Node_Parameter_List, error) {
	l, err := NewNode_Parameter_List(s.Struct.Segment(), n)
	if err != nil {
		return Node_Parameter_List{}, err
	}
	err = s.Struct.SetPtr(5, l.List.ToPtr())
	return l, err
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

func (s Node) HasNestedNodes() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s Node) SetNestedNodes(v Node_NestedNode_List) error {
	return s.Struct.SetPtr(1, v.List.ToPtr())
}

// NewNestedNodes sets the nestedNodes field to a newly
// allocated Node_NestedNode_List, preferring placement in s's segment.
func (s Node) NewNestedNodes(n int32) (Node_NestedNode_List, error) {
	l, err := NewNode_NestedNode_List(s.Struct.Segment(), n)
	if err != nil {
		return Node_NestedNode_List{}, err
	}
	err = s.Struct.SetPtr(1, l.List.ToPtr())
	return l, err
}

func (s Node) Annotations() (Annotation_List, error) {
	p, err := s.Struct.Ptr(2)
	if err != nil {
		return Annotation_List{}, err
	}
	return Annotation_List{List: p.List()}, nil
}

func (s Node) HasAnnotations() bool {
	p, err := s.Struct.Ptr(2)
	return p.IsValid() || err != nil
}

func (s Node) SetAnnotations(v Annotation_List) error {
	return s.Struct.SetPtr(2, v.List.ToPtr())
}

// NewAnnotations sets the annotations field to a newly
// allocated Annotation_List, preferring placement in s's segment.
func (s Node) NewAnnotations(n int32) (Annotation_List, error) {
	l, err := NewAnnotation_List(s.Struct.Segment(), n)
	if err != nil {
		return Annotation_List{}, err
	}
	err = s.Struct.SetPtr(2, l.List.ToPtr())
	return l, err
}

func (s Node) SetFile() {
	s.Struct.SetUint16(12, 0)

}

func (s Node) StructNode() Node_structNode { return Node_structNode(s) }
func (s Node) SetStructNode() {
	s.Struct.SetUint16(12, 1)
}
func (s Node_structNode) DataWordCount() uint16 {
	return s.Struct.Uint16(14)
}

func (s Node_structNode) SetDataWordCount(v uint16) {
	s.Struct.SetUint16(14, v)
}

func (s Node_structNode) PointerCount() uint16 {
	return s.Struct.Uint16(24)
}

func (s Node_structNode) SetPointerCount(v uint16) {
	s.Struct.SetUint16(24, v)
}

func (s Node_structNode) PreferredListEncoding() ElementSize {
	return ElementSize(s.Struct.Uint16(26))
}

func (s Node_structNode) SetPreferredListEncoding(v ElementSize) {
	s.Struct.SetUint16(26, uint16(v))
}

func (s Node_structNode) IsGroup() bool {
	return s.Struct.Bit(224)
}

func (s Node_structNode) SetIsGroup(v bool) {
	s.Struct.SetBit(224, v)
}

func (s Node_structNode) DiscriminantCount() uint16 {
	return s.Struct.Uint16(30)
}

func (s Node_structNode) SetDiscriminantCount(v uint16) {
	s.Struct.SetUint16(30, v)
}

func (s Node_structNode) DiscriminantOffset() uint32 {
	return s.Struct.Uint32(32)
}

func (s Node_structNode) SetDiscriminantOffset(v uint32) {
	s.Struct.SetUint32(32, v)
}

func (s Node_structNode) Fields() (Field_List, error) {
	p, err := s.Struct.Ptr(3)
	if err != nil {
		return Field_List{}, err
	}
	return Field_List{List: p.List()}, nil
}

func (s Node_structNode) HasFields() bool {
	p, err := s.Struct.Ptr(3)
	return p.IsValid() || err != nil
}

func (s Node_structNode) SetFields(v Field_List) error {
	return s.Struct.SetPtr(3, v.List.ToPtr())
}

// NewFields sets the fields field to a newly
// allocated Field_List, preferring placement in s's segment.
func (s Node_structNode) NewFields(n int32) (Field_List, error) {
	l, err := NewField_List(s.Struct.Segment(), n)
	if err != nil {
		return Field_List{}, err
	}
	err = s.Struct.SetPtr(3, l.List.ToPtr())
	return l, err
}

func (s Node) Enum() Node_enum { return Node_enum(s) }
func (s Node) SetEnum() {
	s.Struct.SetUint16(12, 2)
}
func (s Node_enum) Enumerants() (Enumerant_List, error) {
	p, err := s.Struct.Ptr(3)
	if err != nil {
		return Enumerant_List{}, err
	}
	return Enumerant_List{List: p.List()}, nil
}

func (s Node_enum) HasEnumerants() bool {
	p, err := s.Struct.Ptr(3)
	return p.IsValid() || err != nil
}

func (s Node_enum) SetEnumerants(v Enumerant_List) error {
	return s.Struct.SetPtr(3, v.List.ToPtr())
}

// NewEnumerants sets the enumerants field to a newly
// allocated Enumerant_List, preferring placement in s's segment.
func (s Node_enum) NewEnumerants(n int32) (Enumerant_List, error) {
	l, err := NewEnumerant_List(s.Struct.Segment(), n)
	if err != nil {
		return Enumerant_List{}, err
	}
	err = s.Struct.SetPtr(3, l.List.ToPtr())
	return l, err
}

func (s Node) Interface() Node_interface { return Node_interface(s) }
func (s Node) SetInterface() {
	s.Struct.SetUint16(12, 3)
}
func (s Node_interface) Methods() (Method_List, error) {
	p, err := s.Struct.Ptr(3)
	if err != nil {
		return Method_List{}, err
	}
	return Method_List{List: p.List()}, nil
}

func (s Node_interface) HasMethods() bool {
	p, err := s.Struct.Ptr(3)
	return p.IsValid() || err != nil
}

func (s Node_interface) SetMethods(v Method_List) error {
	return s.Struct.SetPtr(3, v.List.ToPtr())
}

// NewMethods sets the methods field to a newly
// allocated Method_List, preferring placement in s's segment.
func (s Node_interface) NewMethods(n int32) (Method_List, error) {
	l, err := NewMethod_List(s.Struct.Segment(), n)
	if err != nil {
		return Method_List{}, err
	}
	err = s.Struct.SetPtr(3, l.List.ToPtr())
	return l, err
}

func (s Node_interface) Superclasses() (Superclass_List, error) {
	p, err := s.Struct.Ptr(4)
	if err != nil {
		return Superclass_List{}, err
	}
	return Superclass_List{List: p.List()}, nil
}

func (s Node_interface) HasSuperclasses() bool {
	p, err := s.Struct.Ptr(4)
	return p.IsValid() || err != nil
}

func (s Node_interface) SetSuperclasses(v Superclass_List) error {
	return s.Struct.SetPtr(4, v.List.ToPtr())
}

// NewSuperclasses sets the superclasses field to a newly
// allocated Superclass_List, preferring placement in s's segment.
func (s Node_interface) NewSuperclasses(n int32) (Superclass_List, error) {
	l, err := NewSuperclass_List(s.Struct.Segment(), n)
	if err != nil {
		return Superclass_List{}, err
	}
	err = s.Struct.SetPtr(4, l.List.ToPtr())
	return l, err
}

func (s Node) Const() Node_const { return Node_const(s) }
func (s Node) SetConst() {
	s.Struct.SetUint16(12, 4)
}
func (s Node_const) Type() (Type, error) {
	p, err := s.Struct.Ptr(3)
	if err != nil {
		return Type{}, err
	}
	return Type{Struct: p.Struct()}, nil
}

func (s Node_const) HasType() bool {
	p, err := s.Struct.Ptr(3)
	return p.IsValid() || err != nil
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

func (s Node_const) HasValue() bool {
	p, err := s.Struct.Ptr(4)
	return p.IsValid() || err != nil
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
func (s Node) SetAnnotation() {
	s.Struct.SetUint16(12, 5)
}
func (s Node_annotation) Type() (Type, error) {
	p, err := s.Struct.Ptr(3)
	if err != nil {
		return Type{}, err
	}
	return Type{Struct: p.Struct()}, nil
}

func (s Node_annotation) HasType() bool {
	p, err := s.Struct.Ptr(3)
	return p.IsValid() || err != nil
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

// Node_Promise is a wrapper for a Node promised by a client call.
type Node_Promise struct{ *capnp.Pipeline }

func (p Node_Promise) Struct() (Node, error) {
	s, err := p.Pipeline.Struct()
	return Node{s}, err
}

func (p Node_Promise) StructNode() Node_structNode_Promise { return Node_structNode_Promise{p.Pipeline} }

// Node_structNode_Promise is a wrapper for a Node_structNode promised by a client call.
type Node_structNode_Promise struct{ *capnp.Pipeline }

func (p Node_structNode_Promise) Struct() (Node_structNode, error) {
	s, err := p.Pipeline.Struct()
	return Node_structNode{s}, err
}

func (p Node_Promise) Enum() Node_enum_Promise { return Node_enum_Promise{p.Pipeline} }

// Node_enum_Promise is a wrapper for a Node_enum promised by a client call.
type Node_enum_Promise struct{ *capnp.Pipeline }

func (p Node_enum_Promise) Struct() (Node_enum, error) {
	s, err := p.Pipeline.Struct()
	return Node_enum{s}, err
}

func (p Node_Promise) Interface() Node_interface_Promise { return Node_interface_Promise{p.Pipeline} }

// Node_interface_Promise is a wrapper for a Node_interface promised by a client call.
type Node_interface_Promise struct{ *capnp.Pipeline }

func (p Node_interface_Promise) Struct() (Node_interface, error) {
	s, err := p.Pipeline.Struct()
	return Node_interface{s}, err
}

func (p Node_Promise) Const() Node_const_Promise { return Node_const_Promise{p.Pipeline} }

// Node_const_Promise is a wrapper for a Node_const promised by a client call.
type Node_const_Promise struct{ *capnp.Pipeline }

func (p Node_const_Promise) Struct() (Node_const, error) {
	s, err := p.Pipeline.Struct()
	return Node_const{s}, err
}

func (p Node_const_Promise) Type() Type_Promise {
	return Type_Promise{Pipeline: p.Pipeline.GetPipeline(3)}
}

func (p Node_const_Promise) Value() Value_Promise {
	return Value_Promise{Pipeline: p.Pipeline.GetPipeline(4)}
}

func (p Node_Promise) Annotation() Node_annotation_Promise { return Node_annotation_Promise{p.Pipeline} }

// Node_annotation_Promise is a wrapper for a Node_annotation promised by a client call.
type Node_annotation_Promise struct{ *capnp.Pipeline }

func (p Node_annotation_Promise) Struct() (Node_annotation, error) {
	s, err := p.Pipeline.Struct()
	return Node_annotation{s}, err
}

func (p Node_annotation_Promise) Type() Type_Promise {
	return Type_Promise{Pipeline: p.Pipeline.GetPipeline(3)}
}

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
	return p.Text(), nil
}

func (s Node_Parameter) HasName() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Node_Parameter) NameBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	d := p.Data()
	if len(d) == 0 {
		return d, nil
	}
	return d[:len(d)-1], nil
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

// Node_Parameter_Promise is a wrapper for a Node_Parameter promised by a client call.
type Node_Parameter_Promise struct{ *capnp.Pipeline }

func (p Node_Parameter_Promise) Struct() (Node_Parameter, error) {
	s, err := p.Pipeline.Struct()
	return Node_Parameter{s}, err
}

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
	return p.Text(), nil
}

func (s Node_NestedNode) HasName() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Node_NestedNode) NameBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	d := p.Data()
	if len(d) == 0 {
		return d, nil
	}
	return d[:len(d)-1], nil
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

// Node_NestedNode_Promise is a wrapper for a Node_NestedNode promised by a client call.
type Node_NestedNode_Promise struct{ *capnp.Pipeline }

func (p Node_NestedNode_Promise) Struct() (Node_NestedNode, error) {
	s, err := p.Pipeline.Struct()
	return Node_NestedNode{s}, err
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
	return p.Text(), nil
}

func (s Field) HasName() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Field) NameBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	d := p.Data()
	if len(d) == 0 {
		return d, nil
	}
	return d[:len(d)-1], nil
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

func (s Field) HasAnnotations() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s Field) SetAnnotations(v Annotation_List) error {
	return s.Struct.SetPtr(1, v.List.ToPtr())
}

// NewAnnotations sets the annotations field to a newly
// allocated Annotation_List, preferring placement in s's segment.
func (s Field) NewAnnotations(n int32) (Annotation_List, error) {
	l, err := NewAnnotation_List(s.Struct.Segment(), n)
	if err != nil {
		return Annotation_List{}, err
	}
	err = s.Struct.SetPtr(1, l.List.ToPtr())
	return l, err
}

func (s Field) DiscriminantValue() uint16 {
	return s.Struct.Uint16(2) ^ 65535
}

func (s Field) SetDiscriminantValue(v uint16) {
	s.Struct.SetUint16(2, v^65535)
}

func (s Field) Slot() Field_slot { return Field_slot(s) }
func (s Field) SetSlot() {
	s.Struct.SetUint16(8, 0)
}
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

func (s Field_slot) HasType() bool {
	p, err := s.Struct.Ptr(2)
	return p.IsValid() || err != nil
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

func (s Field_slot) HasDefaultValue() bool {
	p, err := s.Struct.Ptr(3)
	return p.IsValid() || err != nil
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
func (s Field) SetGroup() {
	s.Struct.SetUint16(8, 1)
}
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

// Field_Promise is a wrapper for a Field promised by a client call.
type Field_Promise struct{ *capnp.Pipeline }

func (p Field_Promise) Struct() (Field, error) {
	s, err := p.Pipeline.Struct()
	return Field{s}, err
}

func (p Field_Promise) Slot() Field_slot_Promise { return Field_slot_Promise{p.Pipeline} }

// Field_slot_Promise is a wrapper for a Field_slot promised by a client call.
type Field_slot_Promise struct{ *capnp.Pipeline }

func (p Field_slot_Promise) Struct() (Field_slot, error) {
	s, err := p.Pipeline.Struct()
	return Field_slot{s}, err
}

func (p Field_slot_Promise) Type() Type_Promise {
	return Type_Promise{Pipeline: p.Pipeline.GetPipeline(2)}
}

func (p Field_slot_Promise) DefaultValue() Value_Promise {
	return Value_Promise{Pipeline: p.Pipeline.GetPipeline(3)}
}

func (p Field_Promise) Group() Field_group_Promise { return Field_group_Promise{p.Pipeline} }

// Field_group_Promise is a wrapper for a Field_group promised by a client call.
type Field_group_Promise struct{ *capnp.Pipeline }

func (p Field_group_Promise) Struct() (Field_group, error) {
	s, err := p.Pipeline.Struct()
	return Field_group{s}, err
}

func (p Field_Promise) Ordinal() Field_ordinal_Promise { return Field_ordinal_Promise{p.Pipeline} }

// Field_ordinal_Promise is a wrapper for a Field_ordinal promised by a client call.
type Field_ordinal_Promise struct{ *capnp.Pipeline }

func (p Field_ordinal_Promise) Struct() (Field_ordinal, error) {
	s, err := p.Pipeline.Struct()
	return Field_ordinal{s}, err
}

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
	return p.Text(), nil
}

func (s Enumerant) HasName() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Enumerant) NameBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	d := p.Data()
	if len(d) == 0 {
		return d, nil
	}
	return d[:len(d)-1], nil
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

func (s Enumerant) HasAnnotations() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s Enumerant) SetAnnotations(v Annotation_List) error {
	return s.Struct.SetPtr(1, v.List.ToPtr())
}

// NewAnnotations sets the annotations field to a newly
// allocated Annotation_List, preferring placement in s's segment.
func (s Enumerant) NewAnnotations(n int32) (Annotation_List, error) {
	l, err := NewAnnotation_List(s.Struct.Segment(), n)
	if err != nil {
		return Annotation_List{}, err
	}
	err = s.Struct.SetPtr(1, l.List.ToPtr())
	return l, err
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

// Enumerant_Promise is a wrapper for a Enumerant promised by a client call.
type Enumerant_Promise struct{ *capnp.Pipeline }

func (p Enumerant_Promise) Struct() (Enumerant, error) {
	s, err := p.Pipeline.Struct()
	return Enumerant{s}, err
}

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

func (s Superclass) HasBrand() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
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

// Superclass_Promise is a wrapper for a Superclass promised by a client call.
type Superclass_Promise struct{ *capnp.Pipeline }

func (p Superclass_Promise) Struct() (Superclass, error) {
	s, err := p.Pipeline.Struct()
	return Superclass{s}, err
}

func (p Superclass_Promise) Brand() Brand_Promise {
	return Brand_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

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
	return p.Text(), nil
}

func (s Method) HasName() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Method) NameBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	d := p.Data()
	if len(d) == 0 {
		return d, nil
	}
	return d[:len(d)-1], nil
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

func (s Method) HasImplicitParameters() bool {
	p, err := s.Struct.Ptr(4)
	return p.IsValid() || err != nil
}

func (s Method) SetImplicitParameters(v Node_Parameter_List) error {
	return s.Struct.SetPtr(4, v.List.ToPtr())
}

// NewImplicitParameters sets the implicitParameters field to a newly
// allocated Node_Parameter_List, preferring placement in s's segment.
func (s Method) NewImplicitParameters(n int32) (Node_Parameter_List, error) {
	l, err := NewNode_Parameter_List(s.Struct.Segment(), n)
	if err != nil {
		return Node_Parameter_List{}, err
	}
	err = s.Struct.SetPtr(4, l.List.ToPtr())
	return l, err
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

func (s Method) HasParamBrand() bool {
	p, err := s.Struct.Ptr(2)
	return p.IsValid() || err != nil
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

func (s Method) HasResultBrand() bool {
	p, err := s.Struct.Ptr(3)
	return p.IsValid() || err != nil
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

func (s Method) HasAnnotations() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s Method) SetAnnotations(v Annotation_List) error {
	return s.Struct.SetPtr(1, v.List.ToPtr())
}

// NewAnnotations sets the annotations field to a newly
// allocated Annotation_List, preferring placement in s's segment.
func (s Method) NewAnnotations(n int32) (Annotation_List, error) {
	l, err := NewAnnotation_List(s.Struct.Segment(), n)
	if err != nil {
		return Annotation_List{}, err
	}
	err = s.Struct.SetPtr(1, l.List.ToPtr())
	return l, err
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

// Method_Promise is a wrapper for a Method promised by a client call.
type Method_Promise struct{ *capnp.Pipeline }

func (p Method_Promise) Struct() (Method, error) {
	s, err := p.Pipeline.Struct()
	return Method{s}, err
}

func (p Method_Promise) ParamBrand() Brand_Promise {
	return Brand_Promise{Pipeline: p.Pipeline.GetPipeline(2)}
}

func (p Method_Promise) ResultBrand() Brand_Promise {
	return Brand_Promise{Pipeline: p.Pipeline.GetPipeline(3)}
}

type Type struct{ capnp.Struct }
type Type_list Type
type Type_enum Type
type Type_structType Type
type Type_interface Type
type Type_anyPointer Type
type Type_anyPointer_unconstrained Type
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
	Type_Which_structType Type_Which = 16
	Type_Which_interface  Type_Which = 17
	Type_Which_anyPointer Type_Which = 18
)

func (w Type_Which) String() string {
	const s = "voidboolint8int16int32int64uint8uint16uint32uint64float32float64textdatalistenumstructTypeinterfaceanyPointer"
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
	case Type_Which_structType:
		return s[80:90]
	case Type_Which_interface:
		return s[90:99]
	case Type_Which_anyPointer:
		return s[99:109]

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

type Type_anyPointer_unconstrained_Which uint16

const (
	Type_anyPointer_unconstrained_Which_anyKind    Type_anyPointer_unconstrained_Which = 0
	Type_anyPointer_unconstrained_Which_struct     Type_anyPointer_unconstrained_Which = 1
	Type_anyPointer_unconstrained_Which_list       Type_anyPointer_unconstrained_Which = 2
	Type_anyPointer_unconstrained_Which_capability Type_anyPointer_unconstrained_Which = 3
)

func (w Type_anyPointer_unconstrained_Which) String() string {
	const s = "anyKindstructlistcapability"
	switch w {
	case Type_anyPointer_unconstrained_Which_anyKind:
		return s[0:7]
	case Type_anyPointer_unconstrained_Which_struct:
		return s[7:13]
	case Type_anyPointer_unconstrained_Which_list:
		return s[13:17]
	case Type_anyPointer_unconstrained_Which_capability:
		return s[17:27]

	}
	return "Type_anyPointer_unconstrained_Which(" + strconv.FormatUint(uint64(w), 10) + ")"
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
func (s Type) SetList() {
	s.Struct.SetUint16(0, 14)
}
func (s Type_list) ElementType() (Type, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return Type{}, err
	}
	return Type{Struct: p.Struct()}, nil
}

func (s Type_list) HasElementType() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
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
func (s Type) SetEnum() {
	s.Struct.SetUint16(0, 15)
}
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

func (s Type_enum) HasBrand() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
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

func (s Type) StructType() Type_structType { return Type_structType(s) }
func (s Type) SetStructType() {
	s.Struct.SetUint16(0, 16)
}
func (s Type_structType) TypeId() uint64 {
	return s.Struct.Uint64(8)
}

func (s Type_structType) SetTypeId(v uint64) {
	s.Struct.SetUint64(8, v)
}

func (s Type_structType) Brand() (Brand, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return Brand{}, err
	}
	return Brand{Struct: p.Struct()}, nil
}

func (s Type_structType) HasBrand() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Type_structType) SetBrand(v Brand) error {
	return s.Struct.SetPtr(0, v.Struct.ToPtr())
}

// NewBrand sets the brand field to a newly
// allocated Brand struct, preferring placement in s's segment.
func (s Type_structType) NewBrand() (Brand, error) {
	ss, err := NewBrand(s.Struct.Segment())
	if err != nil {
		return Brand{}, err
	}
	err = s.Struct.SetPtr(0, ss.Struct.ToPtr())
	return ss, err
}

func (s Type) Interface() Type_interface { return Type_interface(s) }
func (s Type) SetInterface() {
	s.Struct.SetUint16(0, 17)
}
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

func (s Type_interface) HasBrand() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
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
func (s Type) SetAnyPointer() {
	s.Struct.SetUint16(0, 18)
}

func (s Type_anyPointer) Which() Type_anyPointer_Which {
	return Type_anyPointer_Which(s.Struct.Uint16(8))
}
func (s Type_anyPointer) Unconstrained() Type_anyPointer_unconstrained {
	return Type_anyPointer_unconstrained(s)
}
func (s Type_anyPointer) SetUnconstrained() {
	s.Struct.SetUint16(8, 0)
}

func (s Type_anyPointer_unconstrained) Which() Type_anyPointer_unconstrained_Which {
	return Type_anyPointer_unconstrained_Which(s.Struct.Uint16(10))
}
func (s Type_anyPointer_unconstrained) SetAnyKind() {
	s.Struct.SetUint16(10, 0)

}

func (s Type_anyPointer_unconstrained) SetStruct() {
	s.Struct.SetUint16(10, 1)

}

func (s Type_anyPointer_unconstrained) SetList() {
	s.Struct.SetUint16(10, 2)

}

func (s Type_anyPointer_unconstrained) SetCapability() {
	s.Struct.SetUint16(10, 3)

}

func (s Type_anyPointer) Parameter() Type_anyPointer_parameter { return Type_anyPointer_parameter(s) }
func (s Type_anyPointer) SetParameter() {
	s.Struct.SetUint16(8, 1)
}
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
func (s Type_anyPointer) SetImplicitMethodParameter() {
	s.Struct.SetUint16(8, 2)
}
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

// Type_Promise is a wrapper for a Type promised by a client call.
type Type_Promise struct{ *capnp.Pipeline }

func (p Type_Promise) Struct() (Type, error) {
	s, err := p.Pipeline.Struct()
	return Type{s}, err
}

func (p Type_Promise) List() Type_list_Promise { return Type_list_Promise{p.Pipeline} }

// Type_list_Promise is a wrapper for a Type_list promised by a client call.
type Type_list_Promise struct{ *capnp.Pipeline }

func (p Type_list_Promise) Struct() (Type_list, error) {
	s, err := p.Pipeline.Struct()
	return Type_list{s}, err
}

func (p Type_list_Promise) ElementType() Type_Promise {
	return Type_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Type_Promise) Enum() Type_enum_Promise { return Type_enum_Promise{p.Pipeline} }

// Type_enum_Promise is a wrapper for a Type_enum promised by a client call.
type Type_enum_Promise struct{ *capnp.Pipeline }

func (p Type_enum_Promise) Struct() (Type_enum, error) {
	s, err := p.Pipeline.Struct()
	return Type_enum{s}, err
}

func (p Type_enum_Promise) Brand() Brand_Promise {
	return Brand_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Type_Promise) StructType() Type_structType_Promise { return Type_structType_Promise{p.Pipeline} }

// Type_structType_Promise is a wrapper for a Type_structType promised by a client call.
type Type_structType_Promise struct{ *capnp.Pipeline }

func (p Type_structType_Promise) Struct() (Type_structType, error) {
	s, err := p.Pipeline.Struct()
	return Type_structType{s}, err
}

func (p Type_structType_Promise) Brand() Brand_Promise {
	return Brand_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Type_Promise) Interface() Type_interface_Promise { return Type_interface_Promise{p.Pipeline} }

// Type_interface_Promise is a wrapper for a Type_interface promised by a client call.
type Type_interface_Promise struct{ *capnp.Pipeline }

func (p Type_interface_Promise) Struct() (Type_interface, error) {
	s, err := p.Pipeline.Struct()
	return Type_interface{s}, err
}

func (p Type_interface_Promise) Brand() Brand_Promise {
	return Brand_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

func (p Type_Promise) AnyPointer() Type_anyPointer_Promise { return Type_anyPointer_Promise{p.Pipeline} }

// Type_anyPointer_Promise is a wrapper for a Type_anyPointer promised by a client call.
type Type_anyPointer_Promise struct{ *capnp.Pipeline }

func (p Type_anyPointer_Promise) Struct() (Type_anyPointer, error) {
	s, err := p.Pipeline.Struct()
	return Type_anyPointer{s}, err
}

func (p Type_anyPointer_Promise) Unconstrained() Type_anyPointer_unconstrained_Promise {
	return Type_anyPointer_unconstrained_Promise{p.Pipeline}
}

// Type_anyPointer_unconstrained_Promise is a wrapper for a Type_anyPointer_unconstrained promised by a client call.
type Type_anyPointer_unconstrained_Promise struct{ *capnp.Pipeline }

func (p Type_anyPointer_unconstrained_Promise) Struct() (Type_anyPointer_unconstrained, error) {
	s, err := p.Pipeline.Struct()
	return Type_anyPointer_unconstrained{s}, err
}

func (p Type_anyPointer_Promise) Parameter() Type_anyPointer_parameter_Promise {
	return Type_anyPointer_parameter_Promise{p.Pipeline}
}

// Type_anyPointer_parameter_Promise is a wrapper for a Type_anyPointer_parameter promised by a client call.
type Type_anyPointer_parameter_Promise struct{ *capnp.Pipeline }

func (p Type_anyPointer_parameter_Promise) Struct() (Type_anyPointer_parameter, error) {
	s, err := p.Pipeline.Struct()
	return Type_anyPointer_parameter{s}, err
}

func (p Type_anyPointer_Promise) ImplicitMethodParameter() Type_anyPointer_implicitMethodParameter_Promise {
	return Type_anyPointer_implicitMethodParameter_Promise{p.Pipeline}
}

// Type_anyPointer_implicitMethodParameter_Promise is a wrapper for a Type_anyPointer_implicitMethodParameter promised by a client call.
type Type_anyPointer_implicitMethodParameter_Promise struct{ *capnp.Pipeline }

func (p Type_anyPointer_implicitMethodParameter_Promise) Struct() (Type_anyPointer_implicitMethodParameter, error) {
	s, err := p.Pipeline.Struct()
	return Type_anyPointer_implicitMethodParameter{s}, err
}

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

func (s Brand) HasScopes() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Brand) SetScopes(v Brand_Scope_List) error {
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewScopes sets the scopes field to a newly
// allocated Brand_Scope_List, preferring placement in s's segment.
func (s Brand) NewScopes(n int32) (Brand_Scope_List, error) {
	l, err := NewBrand_Scope_List(s.Struct.Segment(), n)
	if err != nil {
		return Brand_Scope_List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
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

// Brand_Promise is a wrapper for a Brand promised by a client call.
type Brand_Promise struct{ *capnp.Pipeline }

func (p Brand_Promise) Struct() (Brand, error) {
	s, err := p.Pipeline.Struct()
	return Brand{s}, err
}

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

func (s Brand_Scope) HasBind() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Brand_Scope) SetBind(v Brand_Binding_List) error {
	s.Struct.SetUint16(8, 0)
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewBind sets the bind field to a newly
// allocated Brand_Binding_List, preferring placement in s's segment.
func (s Brand_Scope) NewBind(n int32) (Brand_Binding_List, error) {
	s.Struct.SetUint16(8, 0)
	l, err := NewBrand_Binding_List(s.Struct.Segment(), n)
	if err != nil {
		return Brand_Binding_List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
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

// Brand_Scope_Promise is a wrapper for a Brand_Scope promised by a client call.
type Brand_Scope_Promise struct{ *capnp.Pipeline }

func (p Brand_Scope_Promise) Struct() (Brand_Scope, error) {
	s, err := p.Pipeline.Struct()
	return Brand_Scope{s}, err
}

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

func (s Brand_Binding) HasType() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
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

// Brand_Binding_Promise is a wrapper for a Brand_Binding promised by a client call.
type Brand_Binding_Promise struct{ *capnp.Pipeline }

func (p Brand_Binding_Promise) Struct() (Brand_Binding, error) {
	s, err := p.Pipeline.Struct()
	return Brand_Binding{s}, err
}

func (p Brand_Binding_Promise) Type() Type_Promise {
	return Type_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

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
	Value_Which_structValue Value_Which = 16
	Value_Which_interface   Value_Which = 17
	Value_Which_anyPointer  Value_Which = 18
)

func (w Value_Which) String() string {
	const s = "voidboolint8int16int32int64uint8uint16uint32uint64float32float64textdatalistenumstructValueinterfaceanyPointer"
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
	case Value_Which_structValue:
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
	return p.Text(), nil
}

func (s Value) HasText() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Value) TextBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	d := p.Data()
	if len(d) == 0 {
		return d, nil
	}
	return d[:len(d)-1], nil
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
	return []byte(p.Data()), nil
}

func (s Value) HasData() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
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

func (s Value) HasList() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Value) ListPtr() (capnp.Ptr, error) {
	return s.Struct.Ptr(0)
}

func (s Value) SetList(v capnp.Pointer) error {
	s.Struct.SetUint16(0, 14)
	return s.Struct.SetPointer(0, v)
}

func (s Value) SetListPtr(v capnp.Ptr) error {
	s.Struct.SetUint16(0, 14)
	return s.Struct.SetPtr(0, v)
}

func (s Value) Enum() uint16 {
	return s.Struct.Uint16(2)
}

func (s Value) SetEnum(v uint16) {
	s.Struct.SetUint16(0, 15)
	s.Struct.SetUint16(2, v)
}

func (s Value) StructValue() (capnp.Pointer, error) {
	return s.Struct.Pointer(0)
}

func (s Value) HasStructValue() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Value) StructValuePtr() (capnp.Ptr, error) {
	return s.Struct.Ptr(0)
}

func (s Value) SetStructValue(v capnp.Pointer) error {
	s.Struct.SetUint16(0, 16)
	return s.Struct.SetPointer(0, v)
}

func (s Value) SetStructValuePtr(v capnp.Ptr) error {
	s.Struct.SetUint16(0, 16)
	return s.Struct.SetPtr(0, v)
}

func (s Value) SetInterface() {
	s.Struct.SetUint16(0, 17)

}

func (s Value) AnyPointer() (capnp.Pointer, error) {
	return s.Struct.Pointer(0)
}

func (s Value) HasAnyPointer() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s Value) AnyPointerPtr() (capnp.Ptr, error) {
	return s.Struct.Ptr(0)
}

func (s Value) SetAnyPointer(v capnp.Pointer) error {
	s.Struct.SetUint16(0, 18)
	return s.Struct.SetPointer(0, v)
}

func (s Value) SetAnyPointerPtr(v capnp.Ptr) error {
	s.Struct.SetUint16(0, 18)
	return s.Struct.SetPtr(0, v)
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

// Value_Promise is a wrapper for a Value promised by a client call.
type Value_Promise struct{ *capnp.Pipeline }

func (p Value_Promise) Struct() (Value, error) {
	s, err := p.Pipeline.Struct()
	return Value{s}, err
}

func (p Value_Promise) List() *capnp.Pipeline {
	return p.Pipeline.GetPipeline(0)
}

func (p Value_Promise) StructValue() *capnp.Pipeline {
	return p.Pipeline.GetPipeline(0)
}

func (p Value_Promise) AnyPointer() *capnp.Pipeline {
	return p.Pipeline.GetPipeline(0)
}

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

func (s Annotation) HasBrand() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
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

func (s Annotation) HasValue() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
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

// Annotation_Promise is a wrapper for a Annotation promised by a client call.
type Annotation_Promise struct{ *capnp.Pipeline }

func (p Annotation_Promise) Struct() (Annotation, error) {
	s, err := p.Pipeline.Struct()
	return Annotation{s}, err
}

func (p Annotation_Promise) Brand() Brand_Promise {
	return Brand_Promise{Pipeline: p.Pipeline.GetPipeline(1)}
}

func (p Annotation_Promise) Value() Value_Promise {
	return Value_Promise{Pipeline: p.Pipeline.GetPipeline(0)}
}

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

func (s CodeGeneratorRequest) HasNodes() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s CodeGeneratorRequest) SetNodes(v Node_List) error {
	return s.Struct.SetPtr(0, v.List.ToPtr())
}

// NewNodes sets the nodes field to a newly
// allocated Node_List, preferring placement in s's segment.
func (s CodeGeneratorRequest) NewNodes(n int32) (Node_List, error) {
	l, err := NewNode_List(s.Struct.Segment(), n)
	if err != nil {
		return Node_List{}, err
	}
	err = s.Struct.SetPtr(0, l.List.ToPtr())
	return l, err
}

func (s CodeGeneratorRequest) RequestedFiles() (CodeGeneratorRequest_RequestedFile_List, error) {
	p, err := s.Struct.Ptr(1)
	if err != nil {
		return CodeGeneratorRequest_RequestedFile_List{}, err
	}
	return CodeGeneratorRequest_RequestedFile_List{List: p.List()}, nil
}

func (s CodeGeneratorRequest) HasRequestedFiles() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s CodeGeneratorRequest) SetRequestedFiles(v CodeGeneratorRequest_RequestedFile_List) error {
	return s.Struct.SetPtr(1, v.List.ToPtr())
}

// NewRequestedFiles sets the requestedFiles field to a newly
// allocated CodeGeneratorRequest_RequestedFile_List, preferring placement in s's segment.
func (s CodeGeneratorRequest) NewRequestedFiles(n int32) (CodeGeneratorRequest_RequestedFile_List, error) {
	l, err := NewCodeGeneratorRequest_RequestedFile_List(s.Struct.Segment(), n)
	if err != nil {
		return CodeGeneratorRequest_RequestedFile_List{}, err
	}
	err = s.Struct.SetPtr(1, l.List.ToPtr())
	return l, err
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

// CodeGeneratorRequest_Promise is a wrapper for a CodeGeneratorRequest promised by a client call.
type CodeGeneratorRequest_Promise struct{ *capnp.Pipeline }

func (p CodeGeneratorRequest_Promise) Struct() (CodeGeneratorRequest, error) {
	s, err := p.Pipeline.Struct()
	return CodeGeneratorRequest{s}, err
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
	return p.Text(), nil
}

func (s CodeGeneratorRequest_RequestedFile) HasFilename() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s CodeGeneratorRequest_RequestedFile) FilenameBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	d := p.Data()
	if len(d) == 0 {
		return d, nil
	}
	return d[:len(d)-1], nil
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

func (s CodeGeneratorRequest_RequestedFile) HasImports() bool {
	p, err := s.Struct.Ptr(1)
	return p.IsValid() || err != nil
}

func (s CodeGeneratorRequest_RequestedFile) SetImports(v CodeGeneratorRequest_RequestedFile_Import_List) error {
	return s.Struct.SetPtr(1, v.List.ToPtr())
}

// NewImports sets the imports field to a newly
// allocated CodeGeneratorRequest_RequestedFile_Import_List, preferring placement in s's segment.
func (s CodeGeneratorRequest_RequestedFile) NewImports(n int32) (CodeGeneratorRequest_RequestedFile_Import_List, error) {
	l, err := NewCodeGeneratorRequest_RequestedFile_Import_List(s.Struct.Segment(), n)
	if err != nil {
		return CodeGeneratorRequest_RequestedFile_Import_List{}, err
	}
	err = s.Struct.SetPtr(1, l.List.ToPtr())
	return l, err
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

// CodeGeneratorRequest_RequestedFile_Promise is a wrapper for a CodeGeneratorRequest_RequestedFile promised by a client call.
type CodeGeneratorRequest_RequestedFile_Promise struct{ *capnp.Pipeline }

func (p CodeGeneratorRequest_RequestedFile_Promise) Struct() (CodeGeneratorRequest_RequestedFile, error) {
	s, err := p.Pipeline.Struct()
	return CodeGeneratorRequest_RequestedFile{s}, err
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
	return p.Text(), nil
}

func (s CodeGeneratorRequest_RequestedFile_Import) HasName() bool {
	p, err := s.Struct.Ptr(0)
	return p.IsValid() || err != nil
}

func (s CodeGeneratorRequest_RequestedFile_Import) NameBytes() ([]byte, error) {
	p, err := s.Struct.Ptr(0)
	if err != nil {
		return nil, err
	}
	d := p.Data()
	if len(d) == 0 {
		return d, nil
	}
	return d[:len(d)-1], nil
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

// CodeGeneratorRequest_RequestedFile_Import_Promise is a wrapper for a CodeGeneratorRequest_RequestedFile_Import promised by a client call.
type CodeGeneratorRequest_RequestedFile_Import_Promise struct{ *capnp.Pipeline }

func (p CodeGeneratorRequest_RequestedFile_Import_Promise) Struct() (CodeGeneratorRequest_RequestedFile_Import, error) {
	s, err := p.Pipeline.Struct()
	return CodeGeneratorRequest_RequestedFile_Import{s}, err
}

const schema_a93fc509624c72d9 = "x\xda\xacZ}\x94T\xe5y\x7f\xdf{gv\xf6k" +
	"\xb8s\xe7\x9d\x15A\xe8,_AVX\xd8]\xa4\xb0" +
	"\x91.\xbb\xb0\x18(\xe0\x0e\x03\xa8\x9cR\xb9;sw" +
	"\xf7\xe2\xec\xccp\xe7\x8e\xee\x12<\x0b\x9cX\x0d\xad\x15" +
	"?\x88J\xd4\xa6)\x9e\x86\xa8-\xb4\xd2\x0a\xd1`\xf6" +
	"HD\x0e\x1aIM\x8d\x896\xe0\x89\xd1cc\x85T" +
	"\xa3(r\xfb<\xef\xfd\xdc\xd9!\x92S\xff\xb8\xec\xdc" +
	"\xdf\xf3\xde\xf7\xe3\xf9~\x9e\x979\xbf\xaaY$4\x05" +
	"\xe3c\x08I\xec\x0dV\x98\x0f\xbc\xb7\xb1z\xda\x82w" +
	"\xee \x89(\x15\xcd\x8d'\xce\x9d~\x7f\xa0\xf0\x0a\xa9" +
	"\xa3!J\x08\xfb\xa8\xe20\xa1\xf0o\x1b\xa1\xe6\xbay" +
	"7^\x18\xbc\xe9\xab\x7fK\x12\x13a\xe4\xd9e\x7f\xf6" +
	"\xddw\xdb\xae\x1f&kad\x80\x06[\xc6\x85\xd6S" +
	"\x18;-\xf4\x0e\x8c\xbd\xaaC\xddxf\xed\xdc]D" +
	"\x0eS\xf3u}Ew\xd5\xd1\xb6}$HC0\xe7" +
	"\xe7\xa1=,X9\x1d~\x8d\xab\x84y\x8fl[\xd9" +
	"\xf2\xd5\xb7N\xecN\x84aVoh\x10\x87.\xab\xfc" +
	".KT\xc2\xaf\x96\x95\x95/\xc0\xec\xe6\xca=k\xde" +
	"\xf8\xef\xad;\x1f 0\xda?\xb1\x80\xa3\xd7V\x1ff" +
	"\x1b\xaa\xf1\xd7\x8d\xd5\xff\x0c\x83\xa3/\\\xd8\xfa\xc3\x15" +
	"\x07\x1e 2\x0b\x987}\xb8<\xbe\xb7\xeb\xc9=\x84" +
	"\xd0\x96\xaa\x9a(e\xe3j`d2V#\xd2d}" +
	"\x8d\x00\x87\xf5\x86\x8c\xdcJg $\xd0\x00\x93k\xf6" +
	"\xc07ca/\xd3j\xee\xc6\xbd\xdc\xbe\\\xf9p\xea" +
	"\xc7'\x1f)\xe1\x87\xc5\xb9\x96\xc1\xdaVd\xc7\xf6\xda" +
	"[a\xe8\xbe;\x86\xc6o\x1e\x18\xf3hy&\x9f\xae" +
	"E&\x9f\xe6#\xaf^\xf2\xd9\x9f~\xfb\xc0w\xf8\xc8" +
	"\xa09v\xef\x94s+\x1e\xdf\xf1\x1bRW\xc1G." +
	"\x0c\x1f\x87\xed\xb7\x87\xf9\xfa\xee\x06\x81\xab\x02;7f" +
	"\x13\xfb|\xcc\xbf\xb3\x87\xa4\xeb\x88`~\xf2\xda\xdeW" +
	"\x1em\xf8\xd6\xbeR>q\x01\x9c\x95\x86\xd99\x09\x7f" +
	"}$\xe1\x9a/\xde\x9c\xdaw\xfb\xdc\xd7\x1e'\x09F" +
	"\x05Ot\x9d\x94\x9f;\x119\xce6D8W#\xc8" +
	"\xd5\xff\xf8\xf1\x8a\x0f6\xe6Z\x9f(\x7f\x16*\xc3\x0e" +
	"YP\xc6y\xdf\xbe~J\xf4\xbe\xf6\xae\x7f\"\x899" +
	"\x94~\xde\xbd\xbdq\xff\x98\xf7~\xc2\xb7\xd0\xa2\xc8\x87" +
	")\x1b\x94q\xd6\"\x1f[\x7f\x7fx\xcb\x13\x8f\xed<" +
	"P\xfe\xdc'\xe5a\x98\xf5$\x1f\xf9\xe0G-\x8b\xe7" +
	"\xfd\xeb\xf2\x83\xe5G\xce\x8a\"/gEQa\x0f\xfc" +
	"\x83t\xf6\xc4\x15\xab\x0f\x119J\xbd\x81\x16\x0f\x94\xe8" +
	"[\xac?\x8a\xbf4>v\xe6M\xdb\x8e<z\xee\xc7" +
	"\x87H\"\x12\xa4\xe6\xf6Us\x9e\xfe\xfam\x9f<\x0b" +
	"\xca\xc2\xbe\x19\xfd)\xdb\x8d#\x93\xbb\xa2\"\xf2\xddx" +
	"k]m\xf4\xc5\x0f\x0e\x97\x97\xfb\xf6\xe8\xf7Q\xee\xf7" +
	"\xf0Y\x7fSs\xc7\x1d\xc3?\xdb\xf5\x03\xe4\xac\xe8\xe9" +
	"\xd7\xda@\x88\x0a4\xc8\x0eF\x7f\x01C\x0fE\xf1X" +
	"\xee\x92\xb4\x1a\xa49\x91\xed`\x93\xd8t\xb6\x90\xa14" +
	"O]58%2\xb4\xef\x87$Q\x05\x9b\xbb3\xf3" +
	"\xe6\x85\xc4\x84\x86\x93\xb8\xb9\x87\xd8N\xf6\xf7\x0c7\xf7" +
	"0\xe3\x9b{\xf9\xb5\xeaI\xbf_z\xf4H\x89\xe1\xa1" +
	"}\xb4\xdc\xc5@\xed\x1fa\xa0\xc2l\x1f\xc35\xdd\xbd" +
	"\x8f\x14%Z\xb4\x08\x92\x9f\x11{\x0fy\x19C\xa9\x9f" +
	"\xbdk\xf6\xd8\xe8\xc6C\xc3\xe4dU\xf0\x824b\x0f" +
	"\xdbc:\xbb=\x86{\xd8\x16\xe3{\xc8\x19O\xdd|" +
	"mp\xca\xf3%\xa7\xae\x0bp\x09m\x8e\xa1,7\xc7" +
	"\xd0M\x9c\x7f\xe7\xd1o]\xf6R\xea\x18\x8e\xa4#5" +
	"\x0fF&\xea~\xc16\xd4q\xcd\xab\xc3\xed\xca\x13\x7f" +
	"\xd9\xf7\xcb\x97\xce\x1f/?\xef\xc1:\xd4\xbcCu\xc8" +
	"\xf7o\xd7\xee\x7f\xed\xa7oLy\x19\xd5_\xf0Y2" +
	"\x0d1\x18\xf9z\xdd\x1ev\x1a\xe7my\xb3nv\x80" +
	"\xb8\xaa\x99\x98L}\x0c\xb4x\xa6\x8e\xdf\x01\x8a:\x1e" +
	"y\xb6}<2\xc2\xe5R\x89\x93\xb0\xa6\x9eu\xc5\xbd" +
	"\xec\xea+\xf0\xc3\xa6+pj\x8fQT\x04\xb1\xce\x8a" +
	"/gM\xf1[\xd97\xe3q\x10\xebW\xc6\xad>|" +
	"\xdb\xdd\xbbO\x82\xac|\x9bD\x95\x8b\x1fg\xbb\xe3x" +
	"\xf0{\xe2/\xc0$\x1d\xc3\xd2\xa7?X{\xe4\xbfP" +
	"N\xa349X\xff\x1e\x93\xeb\xf1W\xb8\x1e\xb9tp" +
	"\xc9\x98\xaf\xd0\x7f\x9bsz\xb4\xb2h\xf5;X?\x8e" +
	"L\xf6\xd5sA\xb93\xc1Y\x82\xbe\xb3T\x84*h" +
	"\x05\xbb\xb1\xfe^\xa6\xd4\x83\x9fn\xb9\xad~\xac\x08\xc3" +
	"\xef\x9e4|\xe6'\xc9\xe9\xef\x96\xb7\xbc\x87\xa6\xbc\x05" +
	"k<2\x05\xf7\xb0K\xa8^\xf4\xea\xb8\xcb~[~" +
	"\xe4\xc4\xa9\xa0W-\x93\xa6\xfeJ\x80\xa1Gj?\xfb" +
	"X;\xfe\xd7\xef\x97w'+\xa7\xe3\xa4\x89\xe98i" +
	"Gq\xda\xe3\xe1\xdd\xc7\xce\x96\xf5\xfdON\x1ff\x07" +
	"\xa7\xe3\xaf\x03\xd3QL\x85T\x9f\xda\xaf4\xa6\xa8\x92" +
	"\xcf\xe6[\xd7\x0c\xe6\xdb\xd4\xc6\x8cV0\x12\x011@" +
	"H\x8cB\xf8\x93\xc3\xdd\x10\x03kE\x9a\xb8\\\xa0\xa6" +
	"\x9aQ\xfb\xd5\xac\xb1\x86\x84\x06\xf3*\x8dx;!\x94" +
	"F|\x13\x06\x9c\x09\xd5F%;\xd8\x95\xd3\xb2\x86\xaa" +
	"7\x16\xb3\xa9\\\xb6`\xe8\x8a\x96\x15\xd5t\"\"\x06" +
	"jM3F\xa3\xb0\x8a\xd2\x01\xab\xfc\x05\xac\xd2'\xd0" +
	"0\xbd\x00\xe88@\xd5V@7\x02\x9a\x01T\xf8\x1c" +
	"\xd0\xf1\x80j\x0d\x80\xa6\x01\xcd\x03*\x9e\x07\xf4\x0a@" +
	"\xfb\xd7\x03\x9a\x01t@\xa0C\xb0\xe8\x9fk\xd94\xa9" +
	"h\x83\xe5\x8a)\x83THx.Ra\xa6\x94\xbc\xd2" +
	"\xade4\"\x1a\x83\xf0:\x92\x03\x1d\xba\"f\xd3\x89" +
	"J\xeas\xf7rU\xb3g\x81r\xb0#\x9eL\xe5\xf2" +
	"\xeaP\x07L\xafe{-N\x05(2\x0a7[\x09" +
	"\x1b\x98*\xd0\xb6\x02\x0e*\x00\x03i\x97\x08\xac\xf1\xa6" +
	"\x03F\x8d\x19\xc5\xf9\x95\xaa\x11\xea\xcb\xa5\xbb(M\xd4" +
	"\xbb\xf3\x9d\xc4c\x9e\x80\xf9~.PJc\x14\xb1\x9f" +
	"\xad\x06\xecU\xc0N\x09T\x16\x01\x04\x0f \xbf\xb9\x03" +
	"\xc07\x00|\x17\xc0\xa0\x10\x03#\"\xf2\xdb;\x01|" +
	"\x17\xc0\x0f\x01\x0c\xc1H\x98V>\x8b\xb2<\x03\xb1\xbc" +
	"\x16<\xab\x1c\x80\xa1A\xd0\x85*\x0a\xacKV\x827" +
	"K\xc6\x10\xaf\x10c\xb4\x02p\x99\xc2\xf0d\x04\xf1\x09" +
	"\x88\x0b\x81\x18\xb7\xa5q\x14\xf2\x03\x80\x00\x9f\x0f\xb8\x94" +
	"U\xfaUZK\x04x\xa8\x99\xca\xa5\xd5\xeb\xf4\xb4J" +
	"\xa8\x8e\x11\x17\x1ej\xe6\x15]\xe9O\x1a:\x05I\xa0" +
	"N\x10Z\x05\x94*\xa0\xe8j\xa1\x981\x92\x06\xd5\x1d" +
	"\x92GS\xb2\xd9\x9c\xa1\x18\x1a\x09\x81\xd2x\x9ct\x15" +
	"\xdc\xe6$\x9f\x1c\x04G@r@v\xdd\xa3\xad\x91\xd6" +
	"\x0a\x1d:\x09)e\xe9Z\x7f>\xa3\xa54\x83v\xe1" +
	"<\xaa\xa1\x8a\xbao178\x96\x15[g\xb6\xd8\xd6" +
	"\xaf\xeaJ\xd6@\xc9\xd5\xba\x92\xebD\xc9-\x02\xde\xaf" +
	"\xf0$\xb7\x0c%\xf75\xc0\xd6 'm\xc9%P\x1e" +
	"]\x96~\x7f1\x1b/\x91!\xce\x1eE\xbe\xc7\xa5\x9a" +
	"\x9aI7fsK\xb4BJ\xd7\xfa\xb5\xac\x92\xa5\xb8" +
	"]\x9c5\x1c2\xcdQ\x87\x82\x0f\xc4L:\x11\xa0\xfe" +
	"\xd4\x90n1\x9d)H\x1b\x9f\xc4HLp\xcf{\x10" +
	"\xcf\xbb\x1f\x8e\xf1\x8cw\xdeCx\xde\xa7\x01{\xdew" +
	"\xde\x1f\xe1y\x9f\x03\xf0\x0d[}\xe1\x00\xf2\xeb\xf7z" +
	"\xea\x1b\x0e\x98&\xf5\xc5F\xf9\xed\x06\xd8f\xf0\x02\x82" +
	"n`\x93_j\x06\x8eTP_\xe6 \x1f\xec \xc2" +
	"\x97\xc6\xc1\xb4\xcd+\x8a\xe7\\\xa7d\x8aT\xf5\xd8%" +
	"\x1529#\xde\xab\xe7\x8a\xf9\xa1\x9c\x0e>@\xc9\x94" +
	"\xb0\xbc\xd4\xed\x81\x86\xb6q\xd5\xd2\xc1=\x04\"`u" +
	"\x10\x04\xe5\x19\xe8\xf1\xa6\xc2\xb1\xe7\x00/h0Fc" +
	"\x00\xce\xda\x02\xe0L\x00\xe7\x83\x13\xe3>dY\xda5" +
	"\x88\xbc\xad\xa1\xa4M_\x96M\xab\x03\xee\xb1\xca\xb9q" +
	"5[\xec\xe7\xcb\x01\x97%\\\xae\xd5[\x0eeT\x87" +
	"\xab5\x03v%`s\xc1c\x19\x83\xfe\xc5\xe2\xddz" +
	"y{q\xd6\x12\xf8Z\xab\x80\xc9\x8d\xb6\x8b%\xa8\x11" +
	"\x91\x10\xf7\x11\xf2A\x1d\xde\x9f\x82\xa9\x9f\xc3\xd3\xd5\xc6" +
	"h%\x80\xcfn\x02\xf0\x19\x00\x8f\xa1N\x84c\xb0\x10" +
	"\x91\x8f\xfe\x0b\x80\xc7\x00|\x15u\xe2T\x0c2;\xf0" +
	"}\x1d\x9e\xef\x93\x03R\x8c\xd6\xa0\xf3CE\xf99\x80" +
	"\xbfF?W\x19\x031\x13\xf94\xf8\xa2\xc4\xaf\x01<" +
	"c{\xae0\x80\xef\xb7Z\xce/\x19\x00\xffd\xa6\x15" +
	"C\xb9\x1eDE\xe2\x8bs\xc5\xac\xe1\xb9%K>\x8b" +
	"\x894\x12\xd6\xd5\x1eU\xd7U\x9a^\x01!\xa33\x9b" +
	"\x8a\xe7\xd0\xd3S\xc9KE\x80\x17\x12\xa1CZ\xe1Z" +
	"T\x03\xd0w\xe0i\xa9\xde\xe0Z\xd4\x9bu\x04\xed\xba" +
	"\x9e\x9e\x82\xa8\x1a\xc0\x14\x01\x1e\xda\xd6\x83F\xeaSH" +
	"_\x196\xc2\xa4\x89\x84\\O\xcc\xf4\xe77r\xa2\xc1" +
	"\x1b/\xafl\xf6\xaa?\xf09^\x89#/[\xef\xd4" +
	"\x90\xf2\xb2V/p\x83\xb3\xf2\xd5\xa1\xed\xcd^V(" +
	"/l\xf6,C^\xb0\xde\x97\x89-\xe8\xf6\xa5\xd0\x0b" +
	"\xbe/\xa1\x1a\xc4\xb9\xa71\xc1'r\x97H\xa8a&" +
	"\x8byUOe\xc05\x17\x0am\x10\xe3 \xc4Ih" +
	"\x1c\xf1\x0eT\xae8Z\x96j\xb6;6)\xe6\xb2f" +
	"\xa7\x95_$IH\xdb\xa2\x9a\x8ba\xdak\xd5\xacJ" +
	"u\xc5\xc8\xe9\xab\xd5\xcdRQ\x85\xe4\x04\x82\x97\xaf\x12" +
	"\x99 \xd0\xf6\x99\xd4\x97\xfd_\x09\xc0\\\x7f\x86\x07\x1a" +
	"\xdf\x0e\x1e\x0f\xb5\x85\xea\x98\x00\x80\xe2\xb6\x16H\x9c3" +
	"\x15X\xcc\x09\xadm\x05\xdf{\x13\xddD\xcd-\xb9\xfe" +
	"nM\xdd\xa2\x06\xb2\x8d\xa9\\\xff\xec\xde\xdcl\xfe\xad" +
	"\x9e3r\xcd\xb3\x0bF\xdaz\x9d]HI\xd6\x87%" +
	"\xb6\xe1\x9c\xbeP \x18\x1e*]w9c\xfcHS" +
	"\xa4%\xa6(j\x7f\xac\x19r~6\xf2\xc4\x84\x90\x92" +
	"X\xd4\xe1\xc5\"(vL;\x1a\xa1\xc7^\x02\xe8F" +
	"L\xac0\xddB\xf7\xbc\x01\xc7\xde\x00hz\xb4\xf7\x91" +
	"\xba!\xe1\xf14\xd4\xcd\x88,\x0d\x1d\xd2\xb2}\xaa\xae" +
	"\x19\xbe\x8cJ\xf0|\xa1\xeb l\x7f\x14)\xe3\x8f." +
	"\xfb\x7f\xfa\xa3 _\xce\xd1\x18GaP_\x1a\xed\xbf" +
	"jz\xa9\x96Q\x1b\xdb\x96\xf5\xe7s\xbaq\x09\"i" +
	"(+\x92\x91Q\xa6\x8c;\xe49\xaeuV\xd1>k" +
	"\x83\xdf\xd5C\x16%\x97\x1cV2\xca\xa6\xd3\xf1[\xd0" +
	"B\x00wm\xb2\xe4\xd0\xd4Y\xd5q\xf8\x01kQ\x9e" +
	"\xb7\xaf\xb7\xf3\xf6+1ow\xacR4|n\xc6\xf5" +
	"\x14%n\xc6w\x18++\x0a\x81\x9bD\x86\xf9\x92\xdd" +
	"\x06;\xd9\x8d\x8d\xca\\\xb8Q\xa4\xae\x12\xaej\xb4L" +
	"\x0d\xa9\x85\xbc\x92\xa2*\xa6\x1d0\xacDf\xa5\xe1\xd2" +
	"\xc9\xc8,\x87\xc17 \xf1\xe0\x89\x87\x838y9." +
	"\xbf\xc5W\x94|a\\\x14|\x99\x10\x0f\xd9\xa2\x92A" +
	"\xf9\xf0\xea\x03\x03\xd3\x8c\xe5\x9e,&B\xf5Qa\x85" +
	"\xa6&\x84\xe7\x00|\x8d\xe0%\x8a\x84\x84Lu\xc0\xf9" +
	"M\xdc\xc5\x9csS\xfb\xdc\x84\x97\x11\xae\x9f\x92\xabV" +
	"\x9b/\xdfs\xf6\xc2\xe0?\xa6\x7f\x07/\x0d\xa6\xc3\x16" +
	"BU\xce\xc02nm\x8e\xdf\xad\xcd\x04`\xbe\xdf\xad" +
	"\xc1^\xdb\xd7\xd8nm\x8f\xe3\xd6\x14j\xbb\xd3\x9c\x98" +
	"-8\xbemr(50\xe08\xb6\xd5\x97\xec\xd8R" +
	"\x03\x14\xbe2{s\xd6\x81h+\xec\xf6f\xa5\xd7r" +
	"1#\x04)^\xcc\xf8\xa8\x81\x89\xa4\xd3:@\xf7\xeb" +
	"\x18#\x89ss\xf4\xdb`\xb3m%\x8b\xd0Jl#" +
	"\\\x88\x82\xbe\x06\xc0\x1b\x04\x1a\xcf\xc2\x02>\xf5uC" +
	"\xa0\xad\xbe\xba3u\x1b\x9f\xda\x1b\xe9\xac_V\xcd}" +
	"\xea'\xa1\xfe\xa1\xef\xac\xe5Y\xa8\xdb\xca\x95;u\xc8" +
	"\xfe(\xcfB\xdd~\xa6\xdc\xb4\x1a@(J\xa9\xaf\xd9" +
	"%O\x1c\x06\xb68u.\x89C\xa5\xab\xa6=\x05\x85" +
	"\xf3;j$\xfa\xd4\x1b\xb5\x9b8l\x8es\x83)\xe5" +
	"\xb0_\x851\x01%X=G\xec\xb2Ni\xf5\x8ag" +
	"\x99ZU\x9d\xac6x\xb5\xb3,X%\x9d\xaca\x06" +
	"\xd6\x07\xa0\x81\xc9\xd66+\xd9\xda\x8c)\x94\x01\xe06" +
	"p\xbb9HM\xbc\xcc\xe4\"\x8e\xc9L\xab=\x0a\x94" +
	"U\xeb !)\xef\xa1\xfa\x94t'\x1a\x09\x85\xa3." +
	"\xc1\xc1b\xc6p\xf3\xa4r\xc1\x8b\xd7\xd3b\xb6\xd7\xf6" +
	"\xcb \x01K)|Y\xb2\xd5\x17(u\xcdC\xc5l" +
	"7\xa4ZP\xeb_l\xb3\xe5x\xc8\xf3w\x98$`" +
	"%\xe4!\x7f\xfd\x1e\x1b\x15}J\x1c.$.b\x11" +
	"E\x94X\xe1\xee\x94\x05\x05\xd8\xd4j\x01\x8bk\xc1\xf2" +
	"\"\x11\xbeYV\x85\x84d\x00)\x11\xa4\x80\xca\x08<" +
	"\xde\xb20\xa7T\"%\x86\x14\xf1\xbcIy\xf5\xced" +
	"\x01\x8c\x01&\x02\xca\xe5H\x09|fZ\xb2fu\x9c" +
	"\x12A\xca\x04\xa4\x04?E\x0a\x96\xf1\xe38%\x86\x94" +
	"z\xa4T\x9c\xc3u\xb0\x90\x9f\xc8)\x97#e*R" +
	"B\x9f\xe07X\xcaO\x12Z\xb1\x94G\xca\x95H\xa9" +
	"\xfc\x18)\xe0\x12\xd94N\xa9G\xcaL\xa4T\xfd\x1e" +
	")\xe0\x15\xd9\x0cN\x99\x8a\x949H\xa9\xfe\x08)\xd5" +
	"\xd8\xd7\x13@Z0\x11P\xe6\"\xa5\xe6C\xa4@\xea" +
	"\xce\x9a8e&R\xe6\x03%\\\xfb\xbf&O\xdf\xd9" +
	"\xd5\x9c\x05s\x90p\x0d\x12\xc2\xbf3y\x0a\xcf\x16p" +
	"\xc2\\$,B\xc2\x98\xb3&\x8fkl!'\xccG" +
	"\xc2\x12\\D:cZ%\x0ek\xe7\x94k\x90\xf25" +
	"\xfc$\xf2\x81\xc9\xe3/\xeb\x14Z\xe1\x89'\xfb\x90d" +
	" I\xfe\x1f\x93Ga\xb6Y\x80:5\x99G\xc2V" +
	"$D\xdf\xb7\x1aRlP\xc0\xa6\xc8\x00\x12\xbe\x01\x04" +
	"\xe9\x96\x9c\x86*\xd6\x9d\xcbe\x1c=\x96 ^\xcd\x07" +
	"A\x0a\xf0\xd08\xbc4\xcd\xc3\xfe%<\xfc\xad\xa5\x19" +
	"\x04&\xc0\xc3\xdf\xe6\xcd\x05!\x09\xf0\xd0x\x91\x7fW" +
	"\x815,\xe4\xfeE\xebC;\x8c\xf0W\xf8\xd2\xa9\x0c" +
	"\x8a\xd6\xa7\xb6\x1e\x0e\xf5dr\x0a\x92\xab\xe1\xbd\xday" +
	"\x07z\x0d\xbc\xd7\xa0\xb9\xaa\x03\x86\x13\x89%\xac{\x80" +
	"\x97p*x\xc1\xf6\x17\x1cL\x80\x87J\x98\x0e\xb8+" +
	"Z\xd9\x19\xf8i\xc1\xeb_'\x02\x10\\bvp\xe9" +
	"6\xad!\xeb\x14\x12\x02k\x07\xb3\xb2\xe71y\xc0\xee" +
	"\xe11\x0c2?'\x86\x13Q\xd5\xdd!\x7fl~F" +
	"x\xe7\xc1\xbd\x14\xc1\xc4\xdcJ\xd8\xfc9\xed\xf82\xfd" +
	"\x95\xe5\xe5\xfa+\xe8=VX\x8e\xd1\x97\xc2\x99=\xb0" +
	"\x10zY\xc2\x83\"\xe7\xd6\x90\xc6W\xf1\x85\x17w\x0f" +
	"e{?\x106\xe8\x1fv\x01\xae\xbb*\xf5\x00V'" +
	"\xb3\x8c\x03\xb0\x9a\x99e\xec?\x0c\xf6_\xd6\xfc\xc3`" +
	"\xfee\xad?\x0c\xd6_\xd6\xf8\xc3`\xfcem?\x0c" +
	"\xb6_\xd6\xf4\xc3`\xfae-?\x0c\x96_\xd6\xf0\xc3" +
	"`\xf8_\xa2\xdd\x83>\xb8\xf7\xaf\xdc\xfc\x850\x1a\xbd" +
	"\xef\xc2\x10\xb6\x8c(\x18<\xf5]\xbd\x01{[\xe1\x89" +
	"\x13n\xef\xd4\xd7E\x07\xa6`\xf0\x06c\xa7\xbe\xfb\x1d" +
	"\x90\xd3z\"\xf8-\x1d\xfe\xa0\xad\x92\x0a\xcb\xba\xad\xbf" +
	"-\xcd\xd6\xdfys\xe1o\xd1\xa2\xdbVl\xff\xc0\x11" +
	"\xb6\xe5\x92\x0a\xc7f\x9d_\x88q;\x85?h\xa1v" +
	"k\x9a\xdb\xe4\x17\xdb\xe2z\xdb\x16\xd7\x0c\x12\x11\x8b=" +
	"\xbf\x09\x8e0@7m#\xd8\x98\xa1\xbeLRn\xea" +
	"\xf0\xb2HyV\xab\xf9\x97\xf7\xfd]\xe2\xd9\xff\xdcy" +
	"\x14\xa2\xedd\xf3\x85\x0f\x8eO\x1d\xf7\x94\xf1\x18\x91\xa7" +
	"M6\xa3\xa7V\xffv\xf0\xafn9F\xe4I\xcd\xe6" +
	"}\xf5\xb3O\xedQ#\x9fB\x92\xb3\xde\xdb\x9d<\xb1" +
	"a\xc8\xce\x0b\xdb,\x13\x0a\xa5s\xa9\x90\xa1\xf4\xc61" +
	"\x11\xed5S\xc5\x82\x91\xeb7\xf8~\xadD\x17O\xe7" +
	"\xe5\xb5x\xba\x88}\xba\xe6\xb8\xbd\xe5\x92\x98\xed\xb4\x05" +
	"\xb4-v}[\xcf\xcd\xfb\xc6f\xb4M91\x19\xfe" +
	"\x08\xbc\xa6\xa5\xa2\xdc\x0e\x8e\x80\x06\xe4\x85\xe0\xd7iP" +
	"^\x00^\x9cV\xc8Ww`\xa6.\xcf\xdaAH\\" +
	"\xed\xcf\x1b\x83\xa1n\xcd\x90\xba\x07\x0d\xd54n\xcdu" +
	"\xc0_H\x96\x89\xd9\x93+\xea\xf8Bh\xc1T\xb5\xde" +
	">\x03^\x80\x99\x85!\xbbU\x04\xdc\xce@F\xb78" +
	"\x17\x80\x93\x164\x18\xe8\xec\xd3\xab\x98VqW\xc6[" +
	"\"%5fC\x99\x1as\xbc\xaf\x08\xf4\x97Q~o" +
	"\xe5%\xe0\x16\x83G\xa5\xdfn\x19HU^txW" +
	"\xbaPt\xb8\xd7br\xd5z\xb3\xcb\x97\x89Z\x1b]" +
	"\x95#bZM\xccw\xb6\xc9\x06\xf1\x8e%i`w" +
	"\x7f\x1buw\xcan\xe3\x97\x01[\x11\xbe\x93_\x06\xd8" +
	"9\xcc\xedt\x18\xf0;\x11\xbf\x1fq\xd1\xba\x7f`\xf7" +
	"P\xb4\xfc\xbfA\xfcA~\xa9`\xe7/\xbb\xf9<\xf7" +
	"#\xbe\x1f\xf1\x90}\x09\xf1$\xc7\x9f@\xfc\x18Eo" +
	"d\xda\xfe\xeb(E\x87\xf0\x1c\x12N \xa1\x8a\xa7\xe1" +
	"\xee\xff\x00`/\xd2Vx\xd0\xc8\xaby.\xee\xde|" +
	"\xb3\xd7):\x85\x9a\xf3\x88\xba\xb7r\xecG\x14M\xbf" +
	"\xf63D\xdd\xfbtX\xbe\x19\xd0\xf0\xa7\x88\xba\xf7r" +
	"\xb0Yp\x08r\x00\x8a\xcez\xbc\xda\xe4\xf7$\xdbp" +
	"'\x0f\xc3N\x9a\x82\xf5p\xa6Ix\xa9G1\x89x" +
	"\x10\x09{\xe9\xc8H\x93\xd6\x0a\xf9\x8c2\xb8\x8a\x84\xfc" +
	"5\xb2\x83\x0a\x98\xe9\xebj\x8f6\xb0B\xcd\xf6\x1a}" +
	"\xc4\x89\xfb\xa3z\xbfYGX\xa1\x11%\x90+\\;" +
	"F]Z\x83[\xc2\xf8\xe7\xde\x8d]\x82\xb7\xb1\x94\x84" +
	"\x10\xee\xa4\xfc.'\xce\x0b\x1doY\xec\xe1\xb9\xd5\xce" +
	"\x1f\xbcI\xd1\x0a<\x0f\xd0\x08M]\xa4.\xe0\xd6d" +
	"\xad\x15RR\xaa\xd3P\x91\xfcU\xc1\"\xbb\xa1\x12\xc7" +
	"Rq\x93W*\x0e\xf5\xf3\xea\xca[\xdf\xe9}:\xb1" +
	"\xdc\xe9LJ\x85\x82\x9f\x9fn\xbb\xf4\xe2\x1d\x11\xeb\xac" +
	"\x12T\xd8\xd9\xc4\\kK\x98/n\xe0jz\x03\xea" +
	"@\x1a\xf5\x9a\xe6y\x9b\x9f)\\\xaf7\"\x9e\xe1v" +
	"\xb3\x99w\xfa\x99F7\xe1\x153\xe2\x06\xb7\x1b\x9d\xf7" +
	"\xe4\xd9f>>\x8f\xf8V~IW\xe0\xbd10\xcb" +
	"\x9d#\xec/h\xc4\xe8Xn\x7f:\xe0\xdf@|\x17" +
	"\xb7\xb3\"\xef\x90\xb0\xbb\xf8\xfc\x9e\xfd\x85n\xe1W\xaa" +
	"\xa0\xd2\x9b\x1c\xfb\xfb\x0e\xe2\x95\xb7\xf2KU\xf6\x08\xc7" +
	"\x1fF\xfc{\x88W\x0d\xf0kU\xf6\x18_\xf7{\x88" +
	"?\x85x\xf5`\x8cN\xc0\x0bd\xbe\xee~\xc4\x9fA" +
	"\xbcfK\x8cN\x04\xfc\x10\x9f\xe7i\xc4\x9fG\xbc\xf6" +
	"\xeb1\xfa'\x04\xcd\xee^\xc0\x9fG\xfc\x15z\xd1\x8e" +
	"\x97i(z\xafj\x14\x96\x92\x10\xa8\xa8\xab\x186\x8a" +
	"\x9dzP\xb7R\xb8\x93\x840\x9f-E\xa9\xdd\xf32" +
	"p\xee\x91\xb4$\x89s\xbd.\xc5\x97\x12\x09\xab\xc3R" +
	"x-\x91\xb2 \xeeR\xf8Z\"\x8dh\xff\xdb\xf02" +
	"j\xdb\x87:z\xe1\x95\x10~P1K\xf1.(\xa4" +
	"\xc1jJ\xe1v\xb7\x9bC\xb3\x17\xb1\x12\xde\xbe(\xb1" +
	"\x12ju\x18K[\xacc\xbf\x94+\x1f\xa7_\x9f\xcb" +
	"\x96v\x9a\xfdY\xb9`g\xe5\xcdv\x9f\xb9\xcb\xee\xe9" +
	"`\xd8^\xd9\xec\xa5\xea\xfe^\xf7E\xba\x9d\x17\xdb\xd7" +
	"\xff\x05\x00\x00\xff\xff\xf7\x86V\xb9"
