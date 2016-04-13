package main

import "fmt"

type annotationParams struct {
	G    *generator
	Node *node
}

type enumParams struct {
	G           *generator
	Node        *node
	Annotations *annotations
	EnumValues  []enumval
}

type structTypesParams struct {
	G           *generator
	Node        *node
	Annotations *annotations
	BaseNode    *node
}

func (p structTypesParams) IsBase() bool {
	return p.Node == p.BaseNode
}

type newStructParams struct {
	G    *generator
	Node *node
}

type structFuncsParams struct {
	G    *generator
	Node *node
}

type structGroupParams struct {
	G     *generator
	Node  *node
	Group *node
	Field field
}

type structFieldParams struct {
	G           *generator
	Node        *node
	Field       field
	Annotations *annotations
	FieldType   string
}

type structBoolFieldParams struct {
	structFieldParams
	Default bool
}

type structUintFieldParams struct {
	structFieldParams
	Bits    uint
	Default uint64
}

func (p structUintFieldParams) Offset() uint32 {
	return p.Field.Slot().Offset() * uint32(p.Bits/8)
}

type structIntFieldParams struct {
	structUintFieldParams
	EnumName string
}

func (p structIntFieldParams) ReturnType() string {
	if p.EnumName != "" {
		return p.EnumName
	}
	return fmt.Sprintf("int%d", p.Bits)
}

type structTextFieldParams struct {
	structFieldParams
	Default string
}

type structDataFieldParams struct {
	structFieldParams
	Default []byte
}

type structObjectFieldParams struct {
	structFieldParams
	TypeNode *node
	Default  staticDataRef
}

type structListParams struct {
	G    *generator
	Node *node
}

type structEnumsParams struct {
	G          *generator
	Node       *node
	Fields     []field
	EnumString enumString
}

type promiseTemplateParams struct {
	G      *generator
	Node   *node
	Fields []field
}

type promiseGroupTemplateParams struct {
	G     *generator
	Node  *node
	Field field
	Group *node
}

type promiseFieldStructTemplateParams struct {
	G       *generator
	Node    *node
	Field   field
	Struct  *node
	Default staticDataRef
}

type promiseFieldAnyPointerTemplateParams struct {
	G     *generator
	Node  *node
	Field field
}

type promiseFieldInterfaceTemplateParams struct {
	G         *generator
	Node      *node
	Field     field
	Interface *node
}

type interfaceClientTemplateParams struct {
	G           *generator
	Node        *node
	Annotations *annotations
	Methods     []interfaceMethod
}

type interfaceServerTemplateParams struct {
	G           *generator
	Node        *node
	Annotations *annotations
	Methods     []interfaceMethod
}

type structValueTemplateParams struct {
	G     *generator
	Node  *node
	Typ   *node
	Value staticDataRef
}

type pointerValueTemplateParams struct {
	G     *generator
	Value staticDataRef
}

type listValueTemplateParams struct {
	G     *generator
	Typ   string
	Value staticDataRef
}
