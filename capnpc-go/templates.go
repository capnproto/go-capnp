// Code generated from templates directory. DO NOT EDIT.

//go:generate ../internal/cmd/mktemplates/mktemplates templates.go templates

package main

func renderAnnotation(r renderer, p annotationParams) error {
	return r.Render("annotation", p)
}
func renderBaseStructFuncs(r renderer, p baseStructFuncsParams) error {
	return r.Render("baseStructFuncs", p)
}
func renderConstants(r renderer, p constantsParams) error {
	return r.Render("constants", p)
}
func renderEnum(r renderer, p enumParams) error {
	return r.Render("enum", p)
}
func renderInterfaceClient(r renderer, p interfaceClientParams) error {
	return r.Render("interfaceClient", p)
}
func renderInterfaceServer(r renderer, p interfaceServerParams) error {
	return r.Render("interfaceServer", p)
}
func renderListValue(r renderer, p listValueParams) error {
	return r.Render("listValue", p)
}
func renderPointerValue(r renderer, p pointerValueParams) error {
	return r.Render("pointerValue", p)
}
func renderPromise(r renderer, p promiseParams) error {
	return r.Render("promise", p)
}
func renderPromiseFieldAnyPointer(r renderer, p promiseFieldAnyPointerParams) error {
	return r.Render("promiseFieldAnyPointer", p)
}
func renderPromiseFieldInterface(r renderer, p promiseFieldInterfaceParams) error {
	return r.Render("promiseFieldInterface", p)
}
func renderPromiseFieldStruct(r renderer, p promiseFieldStructParams) error {
	return r.Render("promiseFieldStruct", p)
}
func renderPromiseGroup(r renderer, p promiseGroupParams) error {
	return r.Render("promiseGroup", p)
}
func renderSchemaVar(r renderer, p schemaVarParams) error {
	return r.Render("schemaVar", p)
}
func renderStructBoolField(r renderer, p structBoolFieldParams) error {
	return r.Render("structBoolField", p)
}
func renderStructDataField(r renderer, p structDataFieldParams) error {
	return r.Render("structDataField", p)
}
func renderStructEnums(r renderer, p structEnumsParams) error {
	return r.Render("structEnums", p)
}
func renderStructFloatField(r renderer, p structFloatFieldParams) error {
	return r.Render("structFloatField", p)
}
func renderStructFuncs(r renderer, p structFuncsParams) error {
	return r.Render("structFuncs", p)
}
func renderStructGroup(r renderer, p structGroupParams) error {
	return r.Render("structGroup", p)
}
func renderStructIntField(r renderer, p structIntFieldParams) error {
	return r.Render("structIntField", p)
}
func renderStructInterfaceField(r renderer, p structInterfaceFieldParams) error {
	return r.Render("structInterfaceField", p)
}
func renderStructList(r renderer, p structListParams) error {
	return r.Render("structList", p)
}
func renderStructListField(r renderer, p structListFieldParams) error {
	return r.Render("structListField", p)
}
func renderStructPointerField(r renderer, p structPointerFieldParams) error {
	return r.Render("structPointerField", p)
}
func renderStructStructField(r renderer, p structStructFieldParams) error {
	return r.Render("structStructField", p)
}
func renderStructTextField(r renderer, p structTextFieldParams) error {
	return r.Render("structTextField", p)
}
func renderStructTypes(r renderer, p structTypesParams) error {
	return r.Render("structTypes", p)
}
func renderStructUintField(r renderer, p structUintFieldParams) error {
	return r.Render("structUintField", p)
}
func renderStructValue(r renderer, p structValueParams) error {
	return r.Render("structValue", p)
}
func renderStructVoidField(r renderer, p structVoidFieldParams) error {
	return r.Render("structVoidField", p)
}
