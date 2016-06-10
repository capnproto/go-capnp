package schema

// AUTO GENERATED - DO NOT EDIT

import (
	math "math"
	strconv "strconv"
	capnp "zombiezen.com/go/capnproto2"
	schemas "zombiezen.com/go/capnproto2/schemas"
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

const schema_a93fc509624c72d9 = "x\xda\xacXkl\x1c\xd5\x15\xbewf\x1f\xde\xf5\x8e" +
	"\xd7\xe3qq\x1eDvH(\x90\x82Il\x93\x06\x17" +
	"\xea\xc4\xc4@\x90\x93z\xed\x04h\xa4\xaaY\xefN\xe2" +
	"\x11\xeb\x99\xcd\xec,\xc4Q\x90\x01\x15\x05\xd1\"\x1e%" +
	"<B\x85*\x1a~D@I\xd4\"5\xd06\xc4\x82" +
	"\x02Qx\xb5\xa04\x10\xaa\x16\x09Hi)qK\x80" +
	"\x90\xc7\xf4;wgg\xc7\xcb\x86\xa6\x12?V3\xf7" +
	"|g\xce=\xf7\x9c\xfb\x9d{\xee\xce\x7f+\xbaXZ" +
	"\x10nm`,\xb5-\x1cq\xbf\xd5\xab\xaf9\xbc\xaa" +
	"\xeb.\xa6*\xdc=`\xf7\x0f\xc7\x9e\xef\xd9\xce\xc2<" +
	"\xca\x98v$\xb2U;\x119\x07oj\xb4\x87\xf1\xdd" +
	"7-\xef\xfc\xce\xbb\xfb\xb6\xa4\x14.\x07T\xc3\xa4\xba" +
	"$\xfa\x88\xb6,\x8a\xb7\xce\xbe\xe8\x0b\x9cqw\xf9\xd6" +
	"\x95\x07\xff\xb1\xe9\xf6\xfb\x19\xb4\x83\x86%\xd2^\x1e{" +
	"Z[\x15\xa3\xb7T\xecI(7\xbdpr\xd3\xef\xfb" +
	"w\xde\xcfT-\xe4\xfe\xf0\x93\xabZ\xb7\x0d<\xb1\x95" +
	"1\xde\xc9\xe3M\\S\xe3\xd0\x1cJ\xc4e>4-" +
	".q\xc6**S]\xe9\x0bE%\x1e\xd2b\xf1\xad" +
	"\xf8\xa6\x05\xbe\xcc\x8a\xdfI\xbe|\xbe\x7f\xdbk\x0f\xcf" +
	"\xbbo{\xb5/b\x91\xeb\xeb'\xb4\xb1zz+\xd6" +
	"\xdf\x00\xe5\x97\xae\xcbl\xbf\xb5k\xffc,\xa5q\xa9" +
	"\x12\x9e>.l\xbf^\xbfW{Gh\x1f\xa8'\xcf" +
	"\xdf\xbbfN\xd3O\x97\x0c\xfc\x92\xa5\xe6s~b\xf8" +
	"\xe6\xf6\x1d\x0d\x1f\xbe*\x0cw\xde\x9cx\x9ak\x0f&" +
	"HwK\x82,\xef\xfcErr\xdf\xcc\xc1]Lm" +
	"\xe2n\xcb\xb69G\xfb\x1f\xbb\xe5}\xcf\x8b\xc9\xc4\xbb" +
	"\xda\x09\xa1{4\x81P\xbb\xaf\xec\x8f\xcf\xfe\xf4\xf2\xe7" +
	"wW\xa5\x85\xa2\xd79]AP\xceS\xb0@\xed\"" +
	"\x85\x0c\x1f\xff\xe0\xe1\xfb\xcex9\xf3\"\xb9\xcc\xa7\xba" +
	"\x0c\x9d\xa7\x94\xb7\xb4=\x0a\x99\xfe\x9d\xd0~(\xb1c" +
	"\xff\x1f\x0f\xcey\x85\xa2!\x05\x82\xc7\xa3\x1at\xa67" +
	"l\xd5f7\xd04\xb3\x1a.\x0c1\x7fM\xa9\xb3x" +
	"\xc0\xab\x92#\x8f6\xdd\xc2\xb5]M\xe4\xc8\x9e&\x8a" +
	"\xc6\x9a}G\xff\xf6\xd1\x86\xc2k\xd5y)\x99\x1e\xd5" +
	"\xee\xd1\x8a\x1a}\xb8^#\xd3\xee7\xa7\x0f>}\xe3" +
	"\x9d[^\xc7\"\x03\x8e0\xae\x8d\xb6\xec\xd5\xc6ZD" +
	"NZ^\x80b\xefD\xf2\x8bgV\xed\xfe\x0bK5" +
	"\xf1/\x85\xee\xef-\x1fjG\x84\xf2d\x0b\xad\xcf\x87" +
	"\xe1D8\xe0D$\x1a\xe1\x11m\xf9\xb4{\xb4U\xd3" +
	"\xce!'\xa6\xb5\xc8d\xbbx\xf6c\xca\x96\x17'k" +
	"n\xd4;fNh[f\xd2\xdb\xdd3\x9fd'\xdd" +
	"BfD\x1fM\xb7gx:o\xe6\xbb{\xed\xb4l" +
	"fSu<\xb0k\xd4XG%\x1fj\xb8\xb7u(" +
	"c\xe5\xf5\xf1^\xc3\xcc\x1a\xe6\xbaTH\x0e1\x16\xc2" +
	"\x1eV\x95n0\xb0N\xe6\xa9\xb9\x12\xef)\x90R\x81" +
	"70> s\xdeX1\xc78\x09\xab\xe6]\xae;" +
	"\xd1\x11+;\xc0y\xaa\xcd\xb7\xf7\xfa<\xd8\xdb\x07{" +
	"\x7f\x968\xe7\xcd\x9cdo\x0eB\xf6\x06d\x7f\x95\xb8" +
	"*C\x88\xfd\xa0\xbes\x0b\x84\x07!<\x04aXj" +
	"\xe62\x84\xef\xdd\x0e\xe1!\x08?\x810\x0aM\x98U" +
	"'\x87!<\x0c\xda%8\xa4!\xa8\x86\x11\x89\x18_" +
	"\x0d:\xd6q\xc8\x9bI\x1e\x91\x9by\x84j\x04\x87\xfa" +
	"P#\xc9\xcf$\xb9\x14j\x16\xf9\x99\xceAe\x88 " +
	"_\x04y\xd2L\x8f\xea<\xc1$\xfc\xb8\x9b\xb1\xb2\xfa" +
	"\xf7\xec\xac\xce\xb8\x0de\x09?\xee\xe6\xd3vzt\xc8" +
	"\xb1y1\xe3\xac\x1c\xcb\x03\x8b\x01\x89\x01\xb1\xf5B1" +
	"\xe7\x0c9\xdc.C\x15,m\x9a\x96\x93v\x0c\x16\xb5" +
	"\xcc@$\xfd\xf4z\x91\x14\xc6\x918\x86\xcc\x01\xf6\xc9" +
	"\x02\xb8\xd1\x9f\xa1\xd7f\xd1tM\xdc\x18\xcd\xe7\x8c\x8c" +
	"\xe1\xf0\x01\xb2\xa3;\xbal\x07&\xf3\x19^3m}" +
	"f\xb1gT\xb7\xd3\xa6C\x99K\xf8\x99\xeb\xa3\xcc-" +
	"F\xec\xfb+\x99[F\x99\xbb\x12\xb2\x95\x14I/s" +
	")\xca\x07fJ\xe5N'\x8c\xa7\x19\x90\xb2\x8f\xb2\xf0" +
	"\xf1rC\xcfe\xdbMk\xa9Q\xc8\xd8\xc6\xa8a\xa6" +
	"MN\xee\x92U%\xea\xba_Z\x14>\x90s\xd9T" +
	"\x88\x07\xab8\xdf\xe8\x96M\xb0\x1ea\xc4I\x9d\xe9\xaf" +
	"\xf7)Z\xef\x0e,\xe3\xb7\x95\xf5\xee\xa2\xf5\xfe\x06\xb2" +
	"\xe7\x02\xeb\xddC\xeb}\x16\xc2\x83\xde\xf6\xc5\x02\xd4\x03" +
	"\xf7T\xb6\xaf\x12r]Ll9\xbf\xbe\xee\x8a\xf0\x9c" +
	"\xe7\xb0\x8d\xe7\xc1\xcd\xf0I\x12\xaa\xb3\xde\x1ey\xfb\xe5" +
	"\xe3{\x99\xfar\x07\"\x12\x81\xe8\xfd\xfa\xcd\x9b'\xde" +
	"\xbc\xeb\x19\xb8\xd0\xcb\xa4\xaf-\x82Y/V\x9c\xd6y" +
	"u:W\xe4z%\\\xc9B\xcerZ\xd7\xd9V1" +
	"?n\xd9\xa8\x01\xe9\x9c\x1fAIDp\xa8\x98\xd7\xed" +
	"L.](0\xda\x17u~\x9c\xce\x9b\x81\x95\xce\xc5" +
	"J\xe7W\xe2tA\x07d\xe7B\xd6%q\xd9\xc8\x96" +
	"w\x7f\xeb\xb0]{\xbfN\x9d\xaa\x97\xb4\xdaEEb" +
	"\xacj\x13\xf6V6!\x8e\x1c\xd7\xdb\x86\x94\xaa\xa5\x90" +
	"\xae\x81T:\xe9\x96\xf2\xf2\x03\xd2\xbd\x16\xd2\xac\xc4\xc7" +
	"E\xe9Z\xe6{\x92\x1cF\xa5\xab\xc4\xca/\x85\xa5X" +
	"\x8d\x1b\xe6\x88n\x1b\x0e\x8b\xf8\x9e\x85\x85g\x97!\xfc" +
	"W\xe8&\xd8\xe1X\xf6\xa0\xbe\xbe\xa8\x17\x9cv\xef\xa9" +
	"g/7rz{\xcf\xb2\xd1\xbce;\xa7\x11\xa3y" +
	"5c45\xdfS\x03\xb3\x02\xd3\xb7\x97(\x1dut" +
	"\x9b\xe6\x08T\xeay^\xa5n\x96NaD>\xd5\x1a" +
	"\xb8C\xcc(\x9f\xa1*\xb7\xdd\xf2\x9aX\xabXUp" +
	")\x1d\xdeR\x16c\xb7\x97\xd7r\xe9F\x08/\x81\xf0" +
	"Z\x89\xb7\x9a\x98 \xb0\x0f\xfd\x83\xce\xdb\x87v\xd9t" +
	"\x8f0]\xd1,\xcf_\xc5\xf8\xe0\x9e\x10\xe7\x93l\xae" +
	"\xf3\xa2\x9bp\xdd\x92O\xbd\x95\xf0*\xfc\xa4\xfb\xe5\x00" +
	"\x8f\x17\xcda\xabhfY$\xe9\xa0 \xc3-\xbf\x09" +
	"\xa8\xda\x82\xa5z\x01~\xc8E\x9d\xa6\xe9\xf7\xa7\xd1\xc2" +
	"\x12,\x0eJt\xd2H\x12\x9f\x85\x89\x1a\xc5LZ\x8c" +
	"\x80\xa1\x10!\x8d\x84H'\\I\xecAM\x11H\x1d" +
	"!\xcd\x84\xc8\xc7].\x8e2M\x95:\xa8]$d" +
	"\x1a!\xa1cn\xe9<\xd3\xbe!\x90FB\xce$$" +
	"\xfc\x05!t\xa6M\x17H3!m\x84D\x8e\xd2<" +
	"t\xaa\xcd\x12\xc84B\xe6\x12\x12\xfd\x9c\xbe\xa1sm" +
	"\xb6\xd4M\xe7\x1a!\xe7\x12R\xf7\x19!u@\xce\x16" +
	"H\x1b!\xe7\x13\x12\xfb\x94\x90\x18\x90\xf3\x042\x97\x90" +
	"\xf9\x84\xc4\x8f\x10\x12\x07r\x81\x84P\xc3\x10\x90.B" +
	"\xea?!\xa4\x1e\xc8\x02\x81\x9cO\xc8\" J\xe2?" +
	"HC\x82\x9a?\x11\x82\xf9\x04\\B\x80\xf2o\x00\x0a" +
	"\x80\x8b\x05\xd0E\xc0b\x02\x1a&\x01\xa0\xe1\xd7.\x15" +
	"\xc0\"\x02\x96\xd2$\xc9\xc34I\x92\xfaw\x81\\B" +
	"\xc8\x95\xf4I\xe3\xc7\xf8\xa4\x11@\x9f\xd4\x8d_\xeb\xd0" +
	"\x08A\x0eA\xea\xbf\x00\xa9\xd49K(\xdaCy\x02" +
	"6\x11\xd0\xf4\x11\x80&\x00c\x12u\x08\x1b\x08\xf8\x11" +
	"\x80\xe4\xf5\x96A\xfbc\xd8\xb2r\xc8*\x88\x0a2\x1a" +
	"\xa6\x83n\x00\x03\x09\xd5\x0b\x83\x05\x0b\x91<\xf0\xb54" +
	"\xea\xec@\xc2\xd0o\x97F\x0b\xbb\x90$\x09?\xdeZ" +
	"\x14\xdfE\xa8\xa03\xdeS,}\xe8\xd5k1\xc4\x97" +
	"u\x18\xd6yC|\xea\xd1\x7f|m\xceJ\x13\x1c\xc7" +
	"8^\x1e\x03\xaf\xc7\xb8\x1e\x1e9\xfa\x06\xa7\xcc\xecd" +
	"6\xed\xa4\x11K\xac\x0a\x83\x9cQp\xb00\x09?\x9e" +
	"\xd4\xcd\xe2\xa8?c\xc1\xa1N\x04\x1c\x97\xdc\xc9;." +
	"liZ\xb3k\x82\xa5B\x12_\x82\xed\x8b\x14\xa1%" +
	"rK*W\xa7Y4W\xa4~\xc5\xb3\xe3\xc29\xdd" +
	"^\x9b\xce0\xae\xa3\x1a\xa6\xcd\xb1\x01\x0b\"&\xeb\xb6" +
	"\xaf\xf2\xff\x96H&\x8ea\xffJ\xa2\xf2n\xaff\x06" +
	"\xeb\xfc\x8c\x1a\xcd\xc6U\xb5\x9a\x0d\xa2~?\x84#S" +
	"\xaa\xa8\xbb\x16\x13Q\x11dX\x8b\x17\xadqC\xcc\x12" +
	"(M\xbe\x0f5\x1b!tn\xfc\xabK\x80_k\xaa" +
	"+\x80\x82\x0aP\xb3\x00((\x005\xf9\xaf\x80\xff5" +
	"\xe9\xaf\x80\xfe5\xd9\xaf\x80\xfd5\xc9\xaf\x80\xfc5\xb9" +
	"\xaf\x80\xfb5\xa9\xaf\x80\xfa5\x99\xaf\x80\xf95\x89\xaf" +
	"\x80\xf8_#\xef\xb1\x1f\xee\xffpM\xfc\xec\x8b?\xd8" +
	"\\\xa2\xbf\xa4\x10\xe9\xb9\xbb}\xf3\xf8\x8c\xf5\x1b\x1a\x1e" +
	"&\x97I\x0a\xc2C\xfa\xa7?\xf4\x7f\xbc\xc6\xea~\x9c" +
	"\xc2\xdb\x8d_+\x13|\x07\xb2;q\xec3c\xef\x8f" +
	"?\xa2\xa0\x0cB\x1fd\x87tr\xd9w\x1f9\xd4s" +
	"\xcd\x04\xe5i5Z\xaa\x00\xd3\xf1 \xae\xb2H\x89\xdd" +
	"\xa5ggG\xe9\xb9\xb0\x0b\xcfb\x09\xf7X\xec\xbd\x90" +
	"\x86\xc7\\\x16)s\xb6\xfcF2\xc1S<\x88\xa1x" +
	"\x107\x05'\xff7\x17W{\\\\9\xc6dj\x80" +
	"\x82\x14\x9cB\xc0\xaa\x13\xb2/\xa7\x8f\xea\xa63dl" +
	"\xf4\xba\xa66A\x90\xefw\xd0\xeeVSg\xe1!\x89" +
	"N\x89\xcb\xea\x12P\x89\x87\xd4KQ\x19yX\xbd\x18" +
	"u\x90G\xd4\x8b\x90H\x1eU/\xc0\xc5\xabU\x1f\xcd" +
	";c\xd1a\xc3I\x0e\x8f9\xba\xeb\xdc`\xf5\xe2Y" +
	"\x00\x9b\xdc\xb5V\xd1\xa6\x01\xe3\x05W7\xd6\x8d8\x18" +
	"\xc0\x9d\xc2x\xde\x12\x9e\xc2\xdf\x9ca\xea\x97Y!\xd0" +
	"\xad`@\xb1\xecg\xa5\x87Y!\x8aA\x92\xde\xab\x1a" +
	"\xa5y5\x1a\xa5\x19\x95s|Jc\x13\xe4\xfbT\xea" +
	"\xc20\xd7\xc5\x8d\xb7\xf2oFl\xb0rAWc\xab" +
	"\xdd\xf2\xd5\x08\xad\xb4[\xf2g\x85\xc5\xe4\xac\x9eZT" +
	"\xf6F\x1b\xe3\x98\x1a\xf78\xec\xd4\x9b\xb8\xef\x90v\xa3" +
	"\xb8Bn\"\xf1m\xe2\x0a\xe9\x1d\xf6\xb7\xf2\x09\xc8o" +
	"#\xf9\xbd$\x97K\xb7V\xednN\x14\xf9\x09\xc9\x1f" +
	"\x10WQ\xef\xa0\xdf\"\xec\xdcK\xf2\x1d$\x8fzW" +
	"\xd7'\x84\xfcq\x92\xbf\xc8\x89\xb6\xaeG\xf4\xe791" +
	"\xe7Y\x02\xf6\x11\x10\x13W\x88\x8b\x96\x1e\xfb\xf6C;" +
	"\x7f\x0e\x8e\xbc\xc4\xbb\xf1#6\xc4O\x10\xf2\xc0\x91\xce" +
	"\xcb\x16\xfe\xea\xaa\xa7\x98v\x80\x13{\xea\x8f\x93\xf4\xce" +
	"\xd9\x13\x87_\x1d:\xe7\x10\xd3\xf6p\xe2H\xe2\x18I" +
	"\xdb\xeeU6>\xfe\xe8\xed;iz\xdcD\x14\xe5\x0b" +
	"\x92\xde%\xc5\x17\xbf1\xfd\x8c\x7f\x92\xb3`\x8e\x1a\x0a" +
	"7\xf36xr\xb3\xb8]\xdfD\x9e\xfc\x0c\x9e,\x08" +
	"\xb7aM\xb3\x01<\xc8\xe9\xb4}\x80\x80m|jI" +
	"\xc6\x05$\x9fK\x8f\xad\xc0\xa55\xd0\x9c\x96\xa5\x12\x84" +
	"\x03\xb6\xbe\xd6\xd8\xd0\xaf\x9b\xeb\x9c\x11V> \xab[" +
	"w\xd7,'+:\xa5\xcf\xf4\x93\xeb\x15\xf3\xd3\xbb\x16" +
	"%\xe9\xa0\x00\x97O\x9b\x96\xa5M\xc2\x98`s\x90\x9b" +
	"\xad\x19\xcc\xe3T\xa6\x95-\xb3t\x8d\xa7=\xf6\x95\xf7" +
	"o\xa3 \x0eL\x83\xf1L\xb9\xf1\xa8\xe2\xf6\x12\xcf\xa8" +
	"eV\xdf\x87\x82\xe7\xa4\xe4\x9d\x93\x1d\xdemh\xc0\xeb" +
	"\xd0\xa9\x0c,\xef\xa8\x1c\x9e\xc1\x1b\xd9\xf5\xb8\x08R;" +
	"\xec\xff\xddVj\x87OuS\xfbo\x00\x00\x00\xff\xff" +
	"\x15\x015\xa8"

func init() {
	schemas.Register(schema_a93fc509624c72d9,
		0x903455f06065422b,
		0x9500cce23b334d80,
		0x978a7cebdc549a4d,
		0x97b14cbe7cfec712,
		0x9aad50a41f4af45f,
		0xa9962a9ed0a4d7f8,
		0xabd73485a9636bc9,
		0xae504193122357e5,
		0xb9521bccf10fa3b1,
		0xbfc546f6210ad7ce,
		0xc863cd16969ee7fc,
		0xce23dcd2d7b00c9b,
		0xcfea0eb02e810062,
		0xd07378ede1f9cc60,
		0xd1958f7dba521926,
		0xdebf55bbfa0fc242,
		0xe682ab4cf923a417,
		0xf1c8950dab257542)
}
