package main

import "fmt"

type annotationParams struct {
	Node *node
}

type enumParams struct {
	Node        *node
	Annotations *annotations
	EnumValues  []enumval
}

type structTypesParams struct {
	Node        *node
	Annotations *annotations
	BaseNode    *node
}

func (p structTypesParams) IsBase() bool {
	return p.Node == p.BaseNode
}

type newStructParams struct {
	Node *node
}

type structFuncsParams struct {
	Node *node
}

type structGroupParams struct {
	Node  *node
	Group *node
	Field field
}

type structFieldParams struct {
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
	Node *node
}

type structEnumsParams struct {
	Node       *node
	Fields     []field
	EnumString enumString
}

type promiseTemplateParams struct {
	Node   *node
	Fields []field
}

type promiseGroupTemplateParams struct {
	Node  *node
	Field field
	Group *node
}

type promiseFieldStructTemplateParams struct {
	Node    *node
	Field   field
	Struct  *node
	Default staticDataRef
}

type promiseFieldAnyPointerTemplateParams struct {
	Node  *node
	Field field
}

type promiseFieldInterfaceTemplateParams struct {
	Node      *node
	Field     field
	Interface *node
}

type interfaceClientTemplateParams struct {
	Node        *node
	Annotations *annotations
	Methods     []interfaceMethod
}

type interfaceServerTemplateParams struct {
	Node        *node
	Annotations *annotations
	Methods     []interfaceMethod
}

type structValueTemplateParams struct {
	Node  *node
	Typ   *node
	Value staticDataRef
}

type pointerValueTemplateParams struct {
	Value staticDataRef
}

type listValueTemplateParams struct {
	Typ   string
	Value staticDataRef
}
