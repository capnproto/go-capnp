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

var schema_a93fc509624c72d9 = []byte{
	120, 218, 172, 90, 125, 116, 148, 229,
	149, 127, 158, 119, 102, 50, 249, 152,
	97, 242, 230, 157, 160, 32, 227, 128,
	64, 13, 17, 66, 72, 192, 133, 84,
	54, 36, 16, 108, 88, 192, 76, 2,
	168, 156, 117, 101, 50, 243, 146, 140,
	78, 102, 194, 204, 27, 77, 40, 158,
	0, 167, 172, 46, 187, 93, 173, 74,
	117, 105, 203, 246, 116, 241, 108, 169,
	118, 23, 118, 101, 87, 168, 20, 225,
	72, 5, 14, 88, 233, 234, 90, 91,
	217, 162, 167, 86, 143, 91, 87, 232,
	106, 253, 230, 221, 223, 125, 222, 207,
	153, 12, 202, 217, 246, 143, 57, 51,
	115, 239, 51, 247, 185, 207, 125, 238,
	189, 191, 123, 239, 59, 141, 191, 170,
	90, 36, 205, 241, 69, 199, 49, 22,
	219, 237, 43, 211, 175, 107, 87, 215,
	157, 95, 61, 247, 1, 38, 7, 185,
	254, 74, 110, 121, 111, 197, 177, 214,
	61, 204, 199, 253, 140, 41, 239, 151,
	237, 84, 62, 43, 187, 22, 159, 100,
	127, 43, 227, 250, 27, 55, 79, 173,
	121, 168, 173, 235, 159, 88, 172, 145,
	243, 207, 122, 183, 52, 236, 29, 247,
	246, 79, 197, 210, 230, 54, 255, 65,
	174, 220, 234, 167, 95, 173, 246, 223,
	141, 181, 255, 241, 147, 229, 239, 174,
	203, 182, 60, 193, 98, 53, 220, 163,
	175, 59, 245, 209, 107, 239, 12, 231,
	95, 96, 181, 220, 207, 177, 102, 159,
	255, 36, 227, 202, 126, 177, 178, 230,
	185, 139, 155, 126, 188, 124, 223, 35,
	76, 86, 188, 250, 237, 239, 45, 139,
	238, 238, 250, 225, 78, 198, 120, 179,
	92, 94, 195, 149, 41, 229, 144, 217,
	51, 169, 220, 195, 123, 234, 202, 37,
	252, 214, 89, 18, 11, 66, 178, 173,
	114, 135, 215, 47, 113, 175, 50, 161,
	124, 39, 126, 115, 5, 84, 154, 85,
	126, 63, 135, 248, 21, 59, 87, 189,
	250, 223, 155, 182, 63, 66, 203, 221,
	39, 148, 72, 215, 45, 21, 7, 149,
	191, 170, 160, 79, 219, 42, 254, 25,
	139, 63, 124, 121, 247, 11, 187, 234,
	191, 185, 167, 120, 177, 48, 199, 188,
	202, 163, 202, 194, 74, 250, 180, 160,
	146, 20, 127, 244, 253, 230, 197, 215,
	255, 235, 178, 253, 116, 68, 159, 126,
	197, 238, 169, 31, 45, 127, 124, 235,
	111, 88, 109, 153, 56, 226, 142, 202,
	131, 56, 226, 142, 74, 50, 92, 251,
	209, 208, 199, 63, 90, 125, 248, 191,
	104, 37, 119, 86, 26, 98, 79, 84,
	190, 173, 188, 36, 196, 158, 17, 98,
	31, 121, 123, 93, 229, 244, 5, 111,
	222, 91, 218, 114, 179, 170, 72, 236,
	172, 42, 18, 123, 225, 235, 179, 175,
	168, 89, 119, 224, 40, 59, 83, 225,
	187, 24, 210, 239, 75, 159, 189, 24,
	155, 84, 127, 6, 182, 83, 98, 85,
	57, 101, 117, 21, 153, 174, 171, 202,
	67, 102, 248, 210, 132, 238, 131, 247,
	220, 191, 227, 12, 238, 89, 114, 14,
	134, 149, 109, 85, 39, 149, 21, 180,
	82, 233, 172, 122, 14, 11, 183, 172,
	108, 124, 234, 171, 247, 124, 120, 136,
	241, 74, 38, 41, 135, 170, 182, 42,
	71, 170, 174, 85, 206, 86, 221, 196,
	36, 103, 7, 238, 1, 111, 65, 96,
	153, 178, 48, 112, 183, 178, 35, 16,
	5, 79, 123, 125, 77, 160, 230, 196,
	187, 7, 89, 44, 2, 181, 47, 116,
	254, 233, 247, 222, 106, 189, 249, 168,
	169, 118, 243, 142, 192, 15, 160, 133,
	242, 88, 64, 232, 109, 49, 11, 79,
	184, 26, 75, 61, 184, 194, 19, 129,
	183, 177, 244, 116, 128, 46, 228, 112,
	224, 147, 15, 82, 39, 255, 250, 157,
	210, 198, 152, 18, 124, 29, 43, 167,
	7, 201, 108, 15, 72, 149, 139, 94,
	156, 48, 254, 183, 165, 111, 99, 75,
	16, 50, 155, 183, 5, 127, 37, 97,
	233, 182, 101, 241, 247, 166, 125, 112,
	230, 59, 165, 85, 221, 87, 221, 66,
	170, 30, 168, 38, 169, 182, 165, 112,
	85, 146, 50, 65, 190, 67, 137, 200,
	255, 174, 28, 144, 201, 26, 191, 169,
	186, 247, 222, 163, 47, 61, 240, 35,
	22, 83, 32, 198, 118, 203, 213, 94,
	63, 151, 184, 79, 145, 107, 126, 1,
	49, 181, 53, 36, 166, 125, 104, 250,
	227, 193, 29, 199, 47, 148, 116, 194,
	145, 154, 163, 202, 150, 26, 250, 116,
	79, 13, 157, 249, 249, 151, 43, 167,
	252, 126, 233, 177, 195, 69, 33, 73,
	107, 155, 107, 21, 196, 196, 116, 5,
	254, 173, 204, 81, 32, 217, 138, 194,
	216, 53, 220, 245, 67, 99, 237, 62,
	101, 43, 87, 142, 137, 181, 167, 21,
	18, 188, 230, 250, 91, 47, 142, 220,
	254, 229, 191, 45, 58, 55, 217, 221,
	203, 125, 205, 35, 225, 181, 116, 240,
	109, 225, 55, 177, 246, 196, 157, 137,
	61, 219, 230, 190, 252, 56, 29, 78,
	114, 178, 68, 7, 23, 81, 118, 91,
	237, 73, 37, 85, 75, 42, 171, 181,
	36, 121, 242, 195, 193, 141, 79, 60,
	182, 125, 95, 105, 227, 7, 199, 31,
	133, 220, 224, 120, 178, 196, 158, 123,
	71, 39, 110, 24, 30, 183, 171, 244,
	133, 166, 198, 147, 119, 167, 196, 74,
	155, 87, 20, 231, 220, 175, 96, 229,
	177, 241, 15, 42, 167, 199, 211, 49,
	79, 140, 159, 237, 197, 242, 79, 223,
	220, 245, 205, 241, 167, 19, 199, 73,
	97, 94, 168, 48, 45, 159, 248, 11,
	229, 204, 68, 82, 248, 244, 68, 8,
	63, 188, 121, 69, 243, 151, 95, 63,
	181, 163, 80, 180, 207, 71, 11, 102,
	92, 245, 61, 101, 206, 85, 36, 121,
	214, 85, 207, 81, 232, 200, 145, 95,
	246, 255, 242, 244, 167, 39, 139, 238,
	185, 214, 43, 84, 158, 23, 161, 84,
	182, 32, 66, 142, 125, 255, 148, 163,
	231, 127, 218, 115, 237, 91, 165, 205,
	112, 107, 132, 188, 245, 182, 8, 29,
	46, 171, 61, 121, 231, 141, 190, 169,
	207, 150, 150, 121, 54, 66, 6, 59,
	27, 161, 139, 248, 86, 96, 239, 203,
	63, 123, 117, 234, 243, 100, 6, 105,
	140, 25, 14, 93, 189, 83, 57, 118,
	53, 41, 123, 228, 106, 97, 134, 121,
	75, 62, 249, 147, 111, 237, 251, 238,
	174, 210, 42, 116, 78, 133, 178, 205,
	43, 166, 138, 204, 104, 51, 33, 217,
	231, 146, 92, 230, 47, 227, 101, 138,
	111, 250, 131, 74, 112, 58, 146, 127,
	243, 140, 233, 87, 120, 176, 124, 255,
	146, 113, 95, 226, 255, 214, 248, 26,
	139, 85, 248, 120, 65, 178, 233, 168,
	219, 170, 116, 214, 81, 178, 89, 82,
	39, 146, 205, 185, 235, 70, 166, 86,
	143, 238, 249, 241, 216, 165, 243, 234,
	182, 43, 11, 197, 210, 249, 198, 210,
	153, 183, 111, 62, 188, 235, 163, 159,
	28, 96, 177, 106, 159, 43, 249, 96,
	233, 140, 186, 159, 97, 57, 45, 109,
	52, 150, 238, 251, 135, 208, 133, 83,
	87, 117, 31, 96, 114, 205, 152, 36,
	58, 165, 238, 117, 101, 86, 157, 184,
	190, 58, 186, 137, 124, 162, 95, 29,
	136, 55, 36, 120, 124, 48, 51, 216,
	210, 158, 139, 123, 50, 201, 88, 57,
	119, 249, 181, 92, 209, 228, 248, 140,
	236, 107, 143, 246, 36, 178, 131, 234,
	104, 123, 42, 147, 76, 101, 250, 98,
	94, 143, 151, 49, 47, 108, 38, 7,
	91, 0, 149, 192, 159, 216, 52, 137,
	183, 230, 105, 81, 158, 143, 99, 188,
	203, 195, 121, 181, 35, 142, 113, 34,
	218, 251, 250, 196, 190, 139, 179, 73,
	245, 70, 53, 163, 230, 226, 90, 54,
	215, 173, 110, 24, 82, 243, 90, 131,
	249, 174, 38, 151, 166, 210, 106, 67,
	107, 231, 192, 96, 54, 167, 117, 113,
	142, 61, 172, 45, 103, 76, 196, 150,
	211, 176, 101, 163, 196, 57, 15, 115,
	162, 205, 170, 7, 13, 134, 136, 205,
	149, 184, 39, 149, 228, 21, 76, 194,
	139, 135, 50, 241, 1, 149, 7, 240,
	37, 224, 218, 94, 18, 219, 175, 26,
	25, 84, 27, 242, 90, 110, 40, 161,
	137, 35, 120, 171, 33, 171, 154, 228,
	183, 20, 202, 31, 79, 242, 155, 28,
	249, 173, 26, 126, 217, 105, 239, 17,
	237, 205, 197, 51, 73, 156, 214, 142,
	42, 156, 182, 218, 181, 157, 71, 108,
	183, 52, 165, 166, 147, 13, 153, 236,
	146, 84, 62, 145, 75, 13, 164, 50,
	241, 12, 167, 131, 81, 234, 12, 250,
	117, 125, 204, 181, 224, 7, 158, 116,
	50, 230, 229, 238, 34, 128, 111, 212,
	45, 17, 172, 85, 8, 209, 98, 147,
	108, 203, 236, 39, 43, 236, 133, 150,
	79, 59, 150, 57, 208, 13, 218, 83,
	160, 61, 43, 113, 89, 2, 17, 1,
	47, 31, 233, 5, 241, 25, 16, 95,
	5, 209, 3, 34, 174, 75, 126, 229,
	65, 16, 95, 5, 241, 45, 137, 7,
	189, 186, 206, 93, 129, 40, 191, 81,
	15, 53, 125, 23, 137, 104, 71, 188,
	124, 186, 9, 38, 40, 3, 201, 78,
	246, 242, 254, 118, 38, 21, 90, 61,
	129, 107, 190, 41, 151, 84, 25, 207,
	209, 89, 241, 226, 122, 60, 147, 201,
	106, 113, 45, 197, 252, 217, 140, 203,
	95, 108, 32, 48, 253, 37, 105, 218,
	138, 211, 57, 215, 196, 211, 67, 92,
	117, 204, 21, 202, 167, 179, 90, 180,
	47, 151, 29, 26, 28, 205, 230, 224,
	150, 241, 116, 145, 5, 59, 50, 67,
	173, 3, 112, 175, 140, 112, 160, 128,
	109, 166, 14, 50, 211, 34, 28, 116,
	185, 99, 166, 78, 50, 211, 87, 64,
	91, 229, 50, 83, 140, 204, 4, 213,
	98, 105, 137, 255, 209, 14, 85, 232,
	133, 61, 67, 131, 106, 46, 145, 142,
	231, 243, 236, 50, 188, 188, 169, 164,
	151, 127, 161, 7, 26, 230, 88, 153,
	77, 182, 170, 13, 106, 102, 104, 128,
	226, 183, 218, 19, 134, 58, 8, 224,
	181, 144, 25, 128, 204, 58, 137, 235,
	196, 36, 131, 49, 143, 230, 58, 129,
	93, 36, 150, 60, 1, 196, 170, 13,
	43, 69, 204, 134, 232, 115, 209, 49,
	234, 75, 28, 99, 162, 115, 140, 2,
	179, 186, 206, 84, 164, 59, 130, 21,
	186, 167, 83, 121, 205, 200, 61, 166,
	238, 189, 166, 238, 87, 146, 238, 105,
	117, 64, 205, 104, 171, 152, 31, 225,
	9, 173, 109, 132, 52, 141, 209, 151,
	109, 72, 144, 44, 22, 109, 161, 61,
	41, 244, 176, 109, 209, 97, 58, 12,
	33, 61, 169, 141, 42, 19, 23, 50,
	89, 56, 194, 173, 77, 36, 69, 142,
	93, 131, 55, 73, 238, 196, 153, 184,
	71, 110, 91, 134, 55, 175, 188, 16,
	158, 195, 125, 242, 2, 216, 145, 151,
	201, 243, 218, 241, 230, 151, 103, 109,
	101, 44, 170, 14, 12, 106, 35, 254,
	222, 148, 22, 234, 29, 209, 84, 93,
	187, 59, 219, 142, 247, 60, 67, 153,
	190, 62, 59, 148, 163, 47, 140, 231,
	117, 53, 213, 215, 175, 225, 11, 243,
	168, 249, 209, 193, 108, 42, 163, 169,
	57, 61, 149, 73, 167, 50, 234, 226,
	172, 7, 9, 48, 159, 194, 66, 161,
	251, 236, 196, 117, 215, 153, 167, 16,
	249, 218, 198, 7, 185, 162, 91, 127,
	254, 27, 23, 46, 142, 252, 99, 242,
	119, 248, 82, 175, 211, 17, 243, 131,
	241, 4, 227, 170, 48, 113, 44, 0,
	124, 116, 224, 100, 146, 196, 219, 26,
	221, 80, 52, 19, 132, 249, 110, 24,
	195, 221, 180, 173, 34, 19, 33, 233,
	236, 212, 197, 158, 45, 45, 113, 110,
	250, 120, 214, 147, 193, 57, 12, 238,
	53, 254, 196, 240, 176, 241, 101, 14,
	239, 230, 250, 198, 236, 64, 111, 74,
	221, 168, 122, 51, 13, 137, 236, 192,
	236, 190, 236, 108, 241, 235, 92, 86,
	203, 54, 205, 206, 107, 201, 217, 230,
	81, 134, 57, 126, 229, 92, 11, 174,
	114, 18, 119, 169, 36, 207, 105, 119,
	212, 145, 103, 181, 232, 127, 241, 208,
	223, 199, 14, 253, 231, 246, 99, 112,
	170, 107, 244, 231, 222, 61, 57, 109,
	194, 147, 218, 99, 76, 158, 126, 141,
	94, 115, 174, 251, 183, 35, 127, 121,
	23, 144, 106, 74, 147, 254, 208, 228,
	217, 231, 118, 170, 213, 31, 51, 57,
	178, 214, 105, 2, 228, 72, 253, 40,
	204, 113, 103, 188, 79, 109, 77, 9,
	80, 241, 39, 179, 9, 191, 22, 239,
	139, 210, 137, 250, 244, 196, 80, 94,
	203, 14, 104, 35, 204, 51, 104, 90,
	204, 11, 139, 57, 6, 242, 194, 30,
	213, 166, 61, 154, 162, 166, 202, 69,
	136, 38, 32, 37, 158, 25, 233, 50,
	46, 177, 1, 27, 165, 83, 137, 148,
	182, 66, 213, 250, 179, 201, 174, 120,
	46, 62, 16, 82, 193, 16, 241, 231,
	11, 243, 43, 201, 135, 55, 186, 124,
	120, 144, 150, 96, 5, 107, 205, 117,
	102, 146, 234, 176, 157, 91, 74, 64,
	151, 177, 79, 136, 54, 162, 220, 22,
	16, 9, 219, 46, 98, 229, 142, 28,
	18, 37, 23, 9, 219, 174, 232, 229,
	57, 221, 32, 74, 159, 17, 209, 238,
	72, 100, 148, 83, 146, 62, 148, 73,
	32, 107, 105, 57, 22, 141, 195, 239,
	146, 142, 38, 28, 190, 104, 30, 195,
	231, 58, 7, 29, 131, 89, 90, 57,
	74, 137, 115, 175, 247, 199, 19, 170,
	133, 168, 114, 9, 68, 189, 226, 15,
	68, 84, 87, 226, 49, 146, 110, 8,
	30, 153, 137, 205, 53, 178, 90, 13,
	234, 156, 219, 56, 98, 180, 231, 22,
	224, 90, 79, 18, 253, 134, 204, 7,
	195, 156, 42, 194, 56, 71, 202, 232,
	89, 71, 244, 52, 209, 165, 13, 97,
	30, 166, 210, 154, 223, 1, 122, 63,
	209, 53, 162, 123, 114, 97, 94, 11,
	250, 6, 177, 126, 144, 232, 155, 136,
	238, 205, 139, 122, 64, 25, 225, 219,
	65, 223, 68, 244, 251, 136, 238, 211,
	196, 169, 148, 109, 60, 7, 250, 215,
	136, 254, 0, 209, 203, 134, 196, 45,
	43, 95, 23, 242, 255, 134, 232, 143,
	18, 221, 127, 87, 152, 79, 160, 62,
	88, 208, 31, 38, 250, 119, 137, 94,
	126, 119, 152, 35, 59, 42, 223, 17,
	244, 111, 19, 253, 251, 68, 175, 24,
	14, 243, 171, 64, 127, 76, 236, 251,
	125, 162, 63, 73, 244, 202, 145, 48,
	159, 68, 35, 3, 177, 239, 94, 162,
	63, 77, 244, 170, 141, 97, 30, 1,
	253, 128, 144, 243, 20, 209, 159, 37,
	122, 224, 171, 97, 126, 53, 232, 71,
	56, 192, 30, 36, 208, 95, 0, 61,
	164, 149, 206, 153, 90, 60, 215, 167,
	106, 249, 165, 204, 143, 106, 12, 185,
	27, 215, 231, 80, 23, 179, 16, 185,
	76, 49, 185, 131, 249, 129, 34, 99,
	168, 220, 132, 22, 141, 100, 23, 242,
	122, 88, 84, 148, 97, 197, 244, 165,
	44, 68, 69, 83, 49, 121, 53, 11,
	101, 112, 221, 197, 228, 27, 89, 136,
	106, 129, 98, 114, 39, 23, 46, 9,
	143, 28, 187, 241, 10, 228, 104, 242,
	232, 98, 122, 23, 11, 145, 143, 23,
	147, 219, 236, 236, 199, 157, 221, 11,
	171, 188, 226, 12, 128, 48, 106, 53,
	130, 69, 196, 131, 36, 188, 80, 158,
	209, 238, 196, 131, 204, 125, 194, 5,
	229, 89, 148, 8, 102, 130, 56, 95,
	226, 163, 162, 146, 118, 34, 226, 139,
	19, 3, 52, 198, 254, 177, 153, 238,
	9, 138, 28, 171, 119, 154, 39, 121,
	69, 147, 3, 228, 40, 119, 156, 105,
	142, 220, 185, 214, 234, 248, 228, 206,
	22, 199, 7, 80, 39, 185, 6, 96,
	109, 77, 78, 179, 37, 47, 108, 114,
	202, 26, 32, 159, 107, 126, 178, 160,
	215, 213, 161, 47, 248, 129, 168, 8,
	162, 162, 246, 213, 59, 172, 226, 130,
	107, 186, 85, 246, 48, 79, 62, 223,
	106, 164, 149, 16, 217, 46, 218, 78,
	161, 31, 165, 90, 79, 213, 219, 172,
	130, 202, 147, 205, 232, 22, 50, 195,
	21, 55, 170, 186, 213, 59, 112, 171,
	121, 8, 81, 215, 80, 2, 228, 102,
	186, 65, 14, 37, 78, 219, 92, 55,
	200, 193, 254, 109, 93, 102, 82, 207,
	89, 32, 151, 103, 81, 97, 84, 11,
	223, 90, 90, 243, 174, 239, 115, 248,
	29, 151, 13, 113, 249, 68, 200, 248,
	97, 81, 230, 50, 122, 1, 81, 180,
	122, 226, 105, 242, 12, 100, 239, 48,
	47, 39, 215, 88, 230, 164, 197, 8,
	210, 119, 89, 24, 247, 143, 196, 77,
	228, 70, 144, 111, 0, 74, 88, 249,
	152, 49, 191, 174, 14, 91, 159, 217,
	37, 192, 194, 178, 98, 54, 195, 138,
	202, 224, 137, 174, 50, 88, 50, 203,
	96, 202, 202, 75, 64, 235, 34, 207,
	52, 203, 224, 21, 77, 78, 109, 236,
	46, 59, 239, 162, 91, 66, 214, 176,
	253, 194, 200, 26, 151, 217, 16, 149,
	104, 255, 184, 70, 237, 142, 53, 181,
	161, 43, 177, 218, 65, 22, 21, 13,
	161, 187, 176, 108, 50, 99, 104, 145,
	169, 41, 17, 23, 82, 12, 221, 0,
	226, 45, 18, 7, 174, 39, 221, 205,
	168, 29, 22, 102, 21, 155, 179, 68,
	183, 10, 209, 206, 74, 107, 255, 255,
	103, 211, 202, 68, 203, 102, 79, 131,
	201, 127, 140, 46, 246, 82, 166, 183,
	58, 144, 101, 165, 58, 16, 202, 20,
	203, 65, 236, 47, 48, 189, 190, 30,
	27, 81, 145, 194, 132, 83, 138, 226,
	121, 212, 40, 107, 92, 39, 182, 117,
	40, 58, 137, 183, 100, 170, 178, 42,
	0, 224, 191, 71, 77, 198, 170, 77,
	167, 4, 154, 202, 113, 210, 226, 207,
	13, 45, 168, 166, 16, 216, 37, 171,
	132, 234, 235, 140, 238, 136, 138, 10,
	129, 92, 114, 138, 10, 254, 36, 168,
	131, 160, 122, 62, 213, 5, 110, 201,
	3, 212, 101, 164, 65, 29, 70, 114,
	195, 166, 127, 150, 202, 36, 89, 89,
	171, 217, 126, 151, 133, 168, 178, 103,
	101, 20, 129, 241, 222, 84, 26, 17,
	143, 66, 172, 172, 200, 145, 69, 106,
	104, 16, 131, 8, 86, 236, 202, 237,
	142, 61, 131, 156, 212, 22, 22, 173,
	55, 157, 121, 29, 41, 72, 106, 147,
	73, 111, 163, 181, 183, 128, 154, 28,
	155, 103, 67, 189, 208, 203, 49, 160,
	61, 1, 49, 12, 56, 154, 202, 244,
	171, 185, 148, 54, 70, 51, 81, 137,
	8, 235, 25, 131, 4, 143, 57, 72,
	168, 119, 167, 121, 175, 81, 11, 185,
	235, 158, 75, 192, 238, 165, 2, 171,
	100, 79, 36, 250, 57, 179, 214, 10,
	149, 168, 181, 106, 255, 192, 90, 203,
	222, 139, 139, 206, 110, 185, 112, 11,
	50, 186, 226, 147, 112, 190, 110, 9,
	37, 68, 64, 178, 252, 130, 232, 21,
	68, 239, 241, 18, 163, 90, 178, 92,
	131, 198, 140, 65, 193, 40, 39, 70,
	88, 178, 188, 195, 67, 15, 80, 36,
	40, 8, 41, 96, 92, 73, 12, 239,
	39, 96, 224, 110, 149, 90, 193, 168,
	38, 198, 36, 98, 248, 62, 6, 195,
	7, 198, 4, 193, 8, 19, 99, 50,
	49, 202, 62, 2, 163, 12, 140, 136,
	96, 92, 73, 140, 105, 196, 240, 127,
	8, 134, 152, 129, 73, 45, 244, 168,
	132, 24, 117, 196, 40, 255, 192, 200,
	184, 202, 116, 193, 152, 76, 140, 153,
	196, 168, 248, 189, 46, 114, 174, 50,
	67, 48, 166, 17, 163, 145, 24, 149,
	239, 131, 81, 73, 79, 25, 36, 184,
	16, 164, 128, 49, 151, 24, 85, 239,
	129, 81, 69, 67, 103, 193, 152, 73,
	140, 249, 196, 8, 252, 47, 24, 1,
	26, 131, 138, 147, 55, 18, 227, 6,
	98, 4, 127, 7, 70, 144, 158, 149,
	8, 198, 92, 98, 44, 34, 198, 184,
	11, 84, 155, 219, 15, 57, 148, 133,
	18, 205, 93, 66, 231, 137, 106, 15,
	135, 161, 50, 81, 171, 223, 37, 170,
	253, 40, 9, 230, 109, 193, 43, 202,
	120, 80, 254, 31, 226, 216, 79, 7,
	96, 20, 170, 251, 107, 222, 33, 170,
	61, 216, 198, 61, 173, 101, 82, 232,
	174, 108, 10, 145, 24, 234, 205, 102,
	211, 120, 67, 42, 152, 207, 202, 162,
	120, 155, 115, 189, 241, 222, 220, 100,
	188, 95, 63, 23, 239, 67, 6, 191,
	117, 200, 92, 32, 62, 208, 10, 241,
	129, 150, 140, 174, 79, 103, 227, 130,
	100, 124, 34, 90, 72, 83, 135, 41,
	200, 147, 113, 45, 110, 198, 122, 136,
	220, 214, 12, 127, 209, 104, 217, 141,
	154, 104, 180, 194, 38, 38, 175, 213,
	141, 37, 171, 68, 91, 6, 28, 77,
	153, 197, 28, 218, 90, 221, 202, 93,
	232, 155, 115, 37, 83, 133, 152, 86,
	122, 50, 125, 230, 68, 194, 244, 219,
	130, 234, 203, 246, 218, 130, 9, 226,
	232, 80, 166, 55, 59, 68, 25, 234,
	82, 133, 113, 97, 116, 160, 128, 241,
	83, 99, 68, 243, 2, 59, 39, 157,
	33, 129, 167, 32, 240, 231, 78, 142,
	127, 137, 166, 76, 47, 130, 118, 206,
	156, 187, 81, 66, 58, 187, 213, 153,
	187, 201, 62, 73, 132, 132, 252, 6,
	234, 252, 216, 91, 32, 190, 71, 221,
	2, 23, 225, 32, 95, 160, 121, 199,
	121, 10, 56, 209, 139, 72, 70, 44,
	84, 240, 181, 20, 86, 84, 203, 135,
	69, 207, 225, 49, 66, 65, 22, 189,
	75, 53, 209, 39, 137, 94, 199, 107,
	68, 194, 4, 190, 147, 34, 129, 232,
	243, 249, 229, 12, 182, 68, 241, 217,
	163, 229, 184, 184, 9, 220, 131, 13,
	67, 57, 53, 63, 148, 214, 208, 50,
	229, 44, 150, 195, 187, 204, 113, 152,
	16, 142, 251, 98, 158, 146, 105, 200,
	216, 161, 61, 199, 252, 165, 211, 148,
	85, 11, 113, 171, 43, 245, 228, 92,
	155, 217, 67, 241, 146, 147, 43, 163,
	12, 19, 19, 68, 216, 213, 107, 212,
	231, 126, 247, 80, 59, 60, 38, 95,
	150, 202, 252, 69, 77, 175, 199, 76,
	196, 237, 238, 2, 5, 182, 143, 82,
	129, 114, 135, 83, 160, 140, 14, 136,
	186, 215, 81, 215, 170, 194, 45, 101,
	173, 26, 57, 148, 207, 187, 11, 25,
	187, 112, 255, 156, 67, 209, 108, 148,
	17, 134, 87, 155, 174, 19, 111, 113,
	32, 92, 230, 134, 231, 200, 106, 189,
	131, 224, 178, 100, 184, 141, 156, 34,
	21, 251, 65, 212, 200, 71, 55, 139,
	148, 39, 111, 128, 203, 196, 52, 16,
	55, 195, 34, 217, 245, 235, 243, 170,
	134, 228, 41, 225, 117, 201, 222, 49,
	169, 174, 143, 227, 234, 214, 160, 51,
	41, 141, 102, 253, 241, 100, 7, 85,
	175, 28, 215, 183, 132, 22, 123, 210,
	218, 152, 182, 202, 8, 48, 244, 3,
	158, 161, 207, 199, 31, 42, 150, 171,
	75, 2, 80, 4, 0, 36, 149, 68,
	160, 8, 16, 136, 151, 132, 160, 8,
	32, 136, 151, 196, 160, 8, 48, 136,
	151, 4, 161, 8, 64, 72, 42, 137,
	66, 17, 160, 16, 47, 9, 67, 17,
	192, 16, 47, 137, 67, 17, 224, 16,
	47, 9, 68, 17, 0, 17, 47, 137,
	68, 17, 32, 17, 255, 35, 66, 145,
	152, 178, 10, 12, 98, 144, 2, 198,
	18, 218, 132, 208, 72, 184, 184, 210,
	38, 56, 55, 16, 231, 43, 244, 19,
	32, 146, 168, 124, 148, 14, 64, 81,
	135, 20, 237, 233, 39, 150, 70, 44,
	64, 146, 168, 127, 148, 13, 192, 34,
	214, 51, 72, 140, 77, 196, 0, 42,
	25, 35, 155, 17, 137, 210, 216, 48,
	49, 190, 6, 134, 27, 152, 44, 175,
	16, 240, 132, 139, 148, 240, 226, 6,
	72, 209, 179, 118, 188, 184, 1, 85,
	184, 48, 9, 47, 110, 0, 22, 46,
	73, 194, 139, 27, 176, 133, 139, 145,
	240, 226, 38, 120, 89, 217, 205, 132,
	48, 203, 155, 77, 32, 179, 194, 221,
	130, 51, 122, 220, 143, 23, 183, 64,
	13, 54, 150, 240, 226, 2, 218, 172,
	228, 41, 0, 14, 182, 196, 169, 240,
	133, 96, 14, 7, 147, 240, 226, 2,
	238, 236, 29, 191, 16, 246, 122, 77,
	216, 91, 19, 103, 126, 196, 14, 2,
	198, 148, 227, 198, 63, 84, 161, 110,
	4, 180, 151, 148, 72, 80, 246, 83,
	174, 73, 72, 8, 126, 35, 199, 237,
	207, 225, 251, 147, 136, 232, 103, 40,
	33, 4, 140, 238, 243, 16, 197, 254,
	211, 32, 30, 167, 132, 16, 52, 122,
	207, 99, 255, 2, 226, 113, 16, 95,
	164, 132, 112, 206, 72, 8, 103, 218,
	29, 116, 147, 189, 33, 225, 116, 242,
	75, 244, 4, 233, 231, 32, 254, 154,
	144, 172, 92, 56, 156, 252, 26, 165,
	142, 95, 131, 120, 222, 196, 38, 56,
	155, 252, 78, 139, 1, 111, 61, 48,
	2, 82, 5, 236, 118, 51, 218, 97,
	22, 93, 12, 216, 213, 28, 224, 49,
	90, 19, 26, 52, 21, 144, 115, 234,
	122, 53, 151, 83, 121, 114, 57, 108,
	220, 145, 73, 68, 179, 244, 84, 146,
	135, 156, 57, 4, 82, 76, 136, 10,
	246, 252, 141, 5, 51, 161, 130, 7,
	74, 180, 23, 119, 164, 22, 240, 110,
	66, 138, 243, 56, 57, 174, 117, 61,
	101, 85, 87, 14, 118, 253, 19, 167,
	32, 7, 219, 79, 90, 184, 42, 6,
	244, 206, 83, 217, 138, 110, 231, 143,
	46, 114, 197, 90, 189, 203, 53, 96,
	53, 158, 158, 172, 204, 50, 79, 82,
	141, 205, 183, 42, 8, 101, 132, 26,
	42, 128, 43, 204, 180, 153, 219, 69,
	132, 114, 143, 192, 117, 103, 246, 40,
	153, 201, 109, 27, 63, 10, 250, 125,
	68, 127, 88, 204, 48, 141, 82, 66,
	249, 6, 111, 47, 152, 61, 150, 153,
	137, 109, 135, 144, 35, 102, 143, 123,
	197, 76, 210, 172, 39, 126, 40, 232,
	79, 16, 253, 56, 167, 74, 89, 55,
	107, 235, 99, 98, 168, 250, 12, 49,
	78, 17, 163, 66, 76, 151, 237, 103,
	234, 202, 9, 222, 130, 23, 21, 160,
	149, 98, 196, 108, 255, 5, 72, 121,
	133, 83, 193, 90, 245, 41, 81, 237,
	191, 1, 40, 71, 56, 149, 165, 129,
	79, 136, 106, 255, 71, 2, 219, 55,
	129, 26, 252, 152, 168, 246, 223, 86,
	160, 44, 138, 85, 217, 235, 11, 243,
	201, 244, 135, 21, 81, 242, 108, 38,
	77, 190, 13, 77, 230, 248, 38, 227,
	76, 83, 192, 248, 59, 78, 217, 229,
	81, 98, 236, 230, 133, 45, 51, 238,
	119, 48, 29, 31, 89, 137, 74, 194,
	85, 239, 88, 84, 137, 30, 9, 193,
	175, 82, 195, 203, 213, 76, 159, 214,
	207, 172, 171, 31, 51, 138, 203, 88,
	151, 229, 47, 24, 47, 216, 151, 107,
	122, 196, 229, 149, 63, 33, 106, 228,
	237, 70, 248, 50, 42, 97, 195, 73,
	24, 19, 25, 197, 157, 14, 162, 162,
	255, 116, 182, 165, 153, 153, 61, 53,
	252, 220, 162, 8, 65, 66, 3, 141,
	20, 227, 9, 59, 80, 172, 167, 50,
	188, 197, 152, 42, 48, 235, 113, 153,
	195, 48, 159, 166, 48, 155, 101, 61,
	161, 146, 204, 71, 84, 45, 214, 51,
	40, 126, 137, 135, 109, 34, 55, 25,
	113, 224, 199, 57, 8, 222, 93, 127,
	44, 168, 119, 106, 176, 130, 18, 245,
	255, 2, 0, 0, 255, 255, 75, 22,
	83, 33,
}
