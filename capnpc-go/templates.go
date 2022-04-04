// Code generated from templates directory. DO NOT EDIT.

//go:generate ../internal/cmd/mktemplates/mktemplates templates.go templates

package main

import (
	"strings"
	"text/template"
)

var templates = template.Must(template.New("").Funcs(template.FuncMap{
	"title": strings.Title,
}).Parse(
	"{{define \"_checktag\"}}{{if .Field.HasDiscriminant}}if s.Struct.Uint16({{.Node.DiscriminantOffset}}) != {{.Field.DiscriminantValue}} {\n  panic({{printf \"Which() != %s\" .Field.Name | printf \"%q\"}})\n}\n{{end}}{{end}}{{define \"_hasfield\"}}func (s {{.Node.Name}}) Has{{.Field.Name | title}}() bool {\n\t{{if .Field.HasDiscriminant}}if s.Struct.Uint16({{.Node.DiscriminantOffset}}) != {{.Field.DiscriminantValue}} {\n\t\treturn false\n\t}\n\t{{end}}return s.Struct.HasPtr({{.Field.Slot.Offset}})\n}\n{{end}}{{define \"_interfaceMethod\"}}\t\t\tInterfaceID: {{.Interface.Id | printf \"%#x\"}},\n\t\t\tMethodID: {{.ID}},\n\t\t\tInterfaceName: {{.Interface.DisplayName | printf \"%q\"}},\n\t\t\tMethodName: {{.OriginalName | printf \"%q\"}},\n{{end}}{{define \"_settag\"}}{{if .Field.HasDiscriminant}}s.Struct.SetUint16({{.Node.DiscriminantOffset}}, {{.Field.DiscriminantValue}})\n{{end}}{{end}}{{define \"_typeid\"}}// {{.Name}}_TypeID is the unique identifier for the type {{.Name}}.\nconst {{.Name}}_TypeID = {{.Id | printf \"%#x\"}}\n{{end}}{{define \"annotation\"}}const {{.Node.Name}} = uint64({{.Node.Id | printf \"%#x\"}})\n{{end}}{{define \"baseStructFuncs\"}}{{template \"_typeid\" .Node}}\n\nfunc New{{.Node.Name}}(s *{{.G.Capnp}}.Segment) ({{.Node.Name}}, error) {\n\tst, err := {{$.G.Capnp}}.NewStruct(s, {{.G.ObjectSize .Node}})\n\treturn {{.Node.Name}}{st}, err\n}\n\nfunc NewRoot{{.Node.Name}}(s *{{.G.Capnp}}.Segment) ({{.Node.Name}}, error) {\n\tst, err := {{.G.Capnp}}.NewRootStruct(s, {{.G.ObjectSize .Node}})\n\treturn {{.Node.Name}}{st}, err\n}\n\nfunc ReadRoot{{.Node.Name}}(msg *{{.G.Capnp}}.Message) ({{.Node.Name}}, error) {\n\troot, err := msg.Root()\n\treturn {{.Node.Name}}{root.Struct()}, err\n}\n{{if .StringMethod}}\nfunc (s {{.Node.Name}}) String() string {\n\tstr, _ := {{.G.Imports.Text}}.Marshal({{.Node.Id | printf \"%#x\"}}, s.Struct)\n\treturn str\n}\n{{end}}\n\n{{end}}{{define \"constants\"}}{{with .Consts}}// Constants defined in {{$.G.Basename}}.\nconst (\n{{range .}}\t{{.Name}} = {{$.G.Value . .Const.Type .Const.Value}}\n{{end}}\n)\n{{end}}\n{{with .Vars}}// Constants defined in {{$.G.Basename}}.\nvar (\n{{range .}}\t{{.Name}} = {{$.G.Value . .Const.Type .Const.Value}}\n{{end}}\n)\n{{end}}\n{{end}}{{define \"enum\"}}{{with .Annotations.Doc}}// {{.}}\n{{end}}type {{.Node.Name}} uint16\n\n{{template \"_typeid\" .Node}}\n\n{{with .EnumValues}}// Values of {{$.Node.Name}}.\nconst (\n{{range .}}{{.FullName}} {{$.Node.Name}} = {{.Val}}\n{{end}}\n)\n\n// String returns the enum's constant name.\nfunc (c {{$.Node.Name}}) String() string {\n\tswitch c {\n\t{{range .}}{{if .Tag}}case {{.FullName}}: return {{printf \"%q\" .Tag}}\n\t{{end}}{{end}}\n\tdefault: return \"\"\n\t}\n}\n\n// {{$.Node.Name}}FromString returns the enum value with a name,\n// or the zero value if there's no such value.\nfunc {{$.Node.Name}}FromString(c string) {{$.Node.Name}} {\n\tswitch c {\n\t{{range .}}{{if .Tag}}case {{printf \"%q\" .Tag}}: return {{.FullName}}\n\t{{end}}{{end}}\n\tdefault: return 0\n\t}\n}\n{{end}}\n\ntype {{.Node.Name}}_List = {{$.G.Capnp}}.EnumList[{{.Node.Name}}]\n\nfunc New{{.Node.Name}}_List(s *{{$.G.Capnp}}.Segment, sz int32) ({{.Node.Name}}_List, error) {\n\treturn {{.G.Capnp}}.NewEnumList[{{.Node.Name}}](s, sz)\n}\n{{end}}{{define \"interfaceClient\"}}{{with .Annotations.Doc}}// {{.}}\n{{end}}type {{.Node.Name}} struct { Client *{{.G.Capnp}}.Client }\n\n{{template \"_typeid\" .Node}}\n\n{{range .Methods}}func (c {{$.Node.Name}}) {{.Name | title}}(ctx {{$.G.Imports.Context}}.Context, params func({{$.G.RemoteNodeName .Params $.Node}}) error) ({{$.G.RemoteNodeName .Results $.Node}}_Future, {{$.G.Capnp}}.ReleaseFunc) {\n\ts := {{$.G.Capnp}}.Send{\n\t\tMethod: {{$.G.Capnp}}.Method{\n\t\t\t{{template \"_interfaceMethod\" .}}\n\t\t},\n\t}\n\tif params != nil {\n\t\ts.ArgsSize = {{$.G.ObjectSize .Params}}\n\t\ts.PlaceArgs = func(s {{$.G.Capnp}}.Struct) error { return params({{$.G.RemoteNodeName .Params $.Node}}{Struct: s}) }\n\t}\n\tans, release := c.Client.SendCall(ctx, s)\n\treturn {{$.G.RemoteNodeName .Results $.Node}}_Future{Future: ans.Future()}, release\n}\n{{end}}\n\nfunc (c {{$.Node.Name}}) AddRef() {{$.Node.Name}} {\n\treturn {{$.Node.Name}} {\n\t\tClient: c.Client.AddRef(),\n\t}\n}\n\nfunc (c {{$.Node.Name}}) Release() {\n\tc.Client.Release()\n}\n{{end}}{{define \"interfaceServer\"}}// A {{.Node.Name}}_Server is a {{.Node.Name}} with a local implementation.\ntype {{.Node.Name}}_Server interface {\n\t{{range .Methods}}\n\t{{.Name | title}}({{$.G.Imports.Context}}.Context, {{$.G.RemoteNodeName .Interface $.Node}}_{{.Name}}) error\n\t{{end}}\n}\n\n// {{.Node.Name}}_NewServer creates a new Server from an implementation of {{.Node.Name}}_Server.\nfunc {{.Node.Name}}_NewServer(s {{.Node.Name}}_Server, policy *{{.G.Imports.Server}}.Policy) *{{.G.Imports.Server}}.Server {\n\tc, _ := s.({{.G.Imports.Server}}.Shutdowner)\n  return {{.G.Imports.Server}}.New({{.Node.Name}}_Methods(nil, s), s, c, policy)\n}\n\n// {{.Node.Name}}_ServerToClient creates a new Client from an implementation of {{.Node.Name}}_Server.\n// The caller is responsible for calling Release on the returned Client.\nfunc {{.Node.Name}}_ServerToClient(s {{.Node.Name}}_Server, policy *{{.G.Imports.Server}}.Policy) {{.Node.Name}} {\n\treturn {{.Node.Name}}{Client: {{.G.Capnp}}.NewClient({{.Node.Name}}_NewServer(s, policy))}\n}\n\n// {{.Node.Name}}_Methods appends Methods to a slice that invoke the methods on s.\n// This can be used to create a more complicated Server.\nfunc {{.Node.Name}}_Methods(methods []{{.G.Imports.Server}}.Method, s {{.Node.Name}}_Server) []{{.G.Imports.Server}}.Method {\n\tif cap(methods) == 0 {\n\t\tmethods = make([]{{.G.Imports.Server}}.Method, 0, {{len .Methods}})\n\t}\n\t{{range .Methods}}\n\tmethods = append(methods, {{$.G.Imports.Server}}.Method{\n\t\tMethod: {{$.G.Capnp}}.Method{\n\t\t\t{{template \"_interfaceMethod\" .}}\n\t\t},\n\t\tImpl: func(ctx {{$.G.Imports.Context}}.Context, call *{{$.G.Imports.Server}}.Call) error {\n\t\t\treturn s.{{.Name | title}}(ctx, {{$.G.RemoteNodeName .Interface $.Node}}_{{.Name}}{call})\n\t\t},\n\t})\n\t{{end}}\n\treturn methods\n}\n{{range .Methods}}{{if eq .Interface.Id $.Node.Id}}\n// {{$.Node.Name}}_{{.Name}} holds the state for a server call to {{$.Node.Name}}.{{.Name}}.\n// See server.Call for documentation.\ntype {{$.Node.Name}}_{{.Name}} struct {\n\t*{{$.G.Imports.Server}}.Call\n}\n\n// Args returns the call's arguments.\nfunc (c {{$.Node.Name}}_{{.Name}}) Args() {{$.G.RemoteNodeName .Params $.Node}} {\n\treturn {{$.G.RemoteNodeName .Params $.Node}}{Struct: c.Call.Args()}\n}\n\n// AllocResults allocates the results struct.\nfunc (c {{$.Node.Name}}_{{.Name}}) AllocResults() ({{$.G.RemoteNodeName .Results $.Node}}, error) {\n\tr, err := c.Call.AllocResults({{$.G.ObjectSize .Results}})\n\treturn {{$.G.RemoteNodeName .Results $.Node}}{Struct: r}, err\n}\n{{end}}{{end}}\n{{end}}{{define \"listValue\"}}{{.Typ}}{List: {{.G.Capnp}}.MustUnmarshalRoot({{.Value}}).List()}{{end}}{{define \"pointerValue\"}}{{.G.Capnp}}.MustUnmarshalRoot({{.Value}}){{end}}{{define \"promise\"}}// {{.Node.Name}}_Future is a wrapper for a {{.Node.Name}} promised by a client call.\ntype {{.Node.Name}}_Future struct { *{{.G.Capnp}}.Future }\n\nfunc (p {{.Node.Name}}_Future) Struct() ({{.Node.Name}}, error) {\n\ts, err := p.Future.Struct()\n\treturn {{.Node.Name}}{s}, err\n}\n\n{{end}}{{define \"promiseFieldAnyPointer\"}}func (p {{.Node.Name}}_Future) {{.Field.Name | title}}() *{{.G.Capnp}}.Future {\n\treturn p.Future.Field({{.Field.Slot.Offset}}, nil)\n}\n\n{{end}}{{define \"promiseFieldInterface\"}}func (p {{.Node.Name}}_Future) {{.Field.Name | title}}() {{.G.RemoteNodeName .Interface .Node}} {\n\treturn {{.G.RemoteNodeName .Interface .Node}}{Client: p.Future.Field({{.Field.Slot.Offset}}, nil).Client()}\n}\n\n{{end}}{{define \"promiseFieldStruct\"}}func (p {{.Node.Name}}_Future) {{.Field.Name | title}}() {{.G.RemoteNodeName .Struct .Node}}_Future {\n\treturn {{.G.RemoteNodeName .Struct .Node}}_Future{Future: p.Future.Field({{.Field.Slot.Offset}}, {{if .Default.IsValid}}{{.Default}}{{else}}nil{{end}})}\n}\n\n{{end}}{{define \"promiseGroup\"}}func (p {{.Node.Name}}_Future) {{.Field.Name | title}}() {{.Group.Name}}_Future { return {{.Group.Name}}_Future{p.Future} }\n{{end}}{{define \"schemaVar\"}}const schema_{{.FileID | printf \"%x\"}} = {{.SchemaLiteral}}\n\nfunc init() {\n  {{.G.Imports.Schemas}}.Register(schema_{{.FileID | printf \"%x\"}},{{range .NodeIDs}}\n\t{{. | printf \"%#x\"}},{{end}})\n}\n{{end}}{{define \"structBoolField\"}}func (s {{.Node.Name}}) {{.Field.Name | title}}() bool {\n\t{{template \"_checktag\" .}}return {{if .Default}}!{{end}}s.Struct.Bit({{.Field.Slot.Offset}})\n}\n\nfunc (s {{.Node.Name}}) Set{{.Field.Name | title}}(v bool) {\n\t{{template \"_settag\" .}}s.Struct.SetBit({{.Field.Slot.Offset}}, {{if .Default}}!{{end}}v)\n}\n\n{{end}}{{define \"structDataField\"}}func (s {{.Node.Name}}) {{.Field.Name | title}}() ({{.FieldType}}, error) {\n\t{{template \"_checktag\" .}}p, err := s.Struct.Ptr({{.Field.Slot.Offset}})\n\t{{with .Default}}return {{$.FieldType}}(p.DataDefault({{printf \"%#v\" .}})), err{{else}}return {{.FieldType}}(p.Data()), err{{end}}\n}\n\n{{template \"_hasfield\" .}}\n\nfunc (s {{.Node.Name}}) Set{{.Field.Name | title}}(v {{.FieldType}}) error {\n\t{{template \"_settag\" .}}{{if .Default}}if v == nil {\n\t\tv = []byte{}\n\t}\n\t{{end}}return s.Struct.SetData({{.Field.Slot.Offset}}, v)\n}\n\n{{end}}{{define \"structEnums\"}}type {{.Node.Name}}_Which uint16\n\nconst (\n{{range .Fields}}\t{{$.Node.Name}}_Which_{{.Name}} {{$.Node.Name}}_Which = {{.DiscriminantValue}}\n{{end}}\n)\n\nfunc (w {{.Node.Name}}_Which) String() string {\n\tconst s = {{.EnumString.ValueString | printf \"%q\"}}\n\tswitch w {\n\t{{range $i, $f := .Fields}}case {{$.Node.Name}}_Which_{{.Name}}:\n\t\treturn s{{$.EnumString.SliceFor $i}}\n\t{{end}}\n\t}\n\treturn \"{{.Node.Name}}_Which(\" + {{.G.Imports.Strconv}}.FormatUint(uint64(w), 10) + \")\"\n}\n\n{{end}}{{define \"structFloatField\"}}func (s {{.Node.Name}}) {{.Field.Name | title}}() float{{.Bits}} {\n\t{{template \"_checktag\" .}}return {{.G.Imports.Math}}.Float{{.Bits}}frombits(s.Struct.Uint{{.Bits}}({{.Offset}}){{with .Default}} ^ {{printf \"%#x\" .}}{{end}})\n}\n\nfunc (s {{.Node.Name}}) Set{{.Field.Name | title}}(v float{{.Bits}}) {\n\t{{template \"_settag\" .}}s.Struct.SetUint{{.Bits}}({{.Offset}}, {{.G.Imports.Math}}.Float{{.Bits}}bits(v){{with .Default}}^{{printf \"%#x\" .}}{{end}})\n}\n\n{{end}}{{define \"structFuncs\"}}{{if gt .Node.StructNode.DiscriminantCount 0}}\nfunc (s {{.Node.Name}}) Which() {{.Node.Name}}_Which {\n\treturn {{.Node.Name}}_Which(s.Struct.Uint16({{.Node.DiscriminantOffset}}))\n}\n{{end}}{{end}}{{define \"structGroup\"}}func (s {{.Node.Name}}) {{.Field.Name | title}}() {{.Group.Name}} { return {{.Group.Name}}(s) }\n{{if .Field.HasDiscriminant}}\nfunc (s {{.Node.Name}}) Set{{.Field.Name | title}}() { {{template \"_settag\" .}} }\n{{end}}\n{{end}}{{define \"structIntField\"}}func (s {{.Node.Name}}) {{.Field.Name | title}}() {{.ReturnType}} {\n\t{{template \"_checktag\" .}}return {{.ReturnType}}(s.Struct.Uint{{.Bits}}({{.Offset}}){{with .Default}} ^ {{.}}{{end}})\n}\n\nfunc (s {{.Node.Name}}) Set{{.Field.Name | title}}(v {{.ReturnType}}) {\n\t{{template \"_settag\" .}}s.Struct.SetUint{{.Bits}}({{.Offset}}, uint{{.Bits}}(v){{with .Default}}^{{.}}{{end}})\n}\n\n{{end}}{{define \"structInterfaceField\"}}func (s {{.Node.Name}}) {{.Field.Name | title}}() {{.FieldType}} {\n\t{{template \"_checktag\" .}}p, _ := s.Struct.Ptr({{.Field.Slot.Offset}})\n\treturn {{.FieldType}}{Client: p.Interface().Client()}\n}\n\n{{template \"_hasfield\" .}}\n\nfunc (s {{.Node.Name}}) Set{{.Field.Name | title}}(v {{.FieldType}}) error {\n\t{{template \"_settag\" .}}if !v.Client.IsValid() {\n\t\treturn s.Struct.SetPtr({{.Field.Slot.Offset}}, capnp.Ptr{})\n\t}\n\tseg := s.Segment()\n\tin := {{.G.Capnp}}.NewInterface(seg, seg.Message().AddCap(v.Client))\n\treturn s.Struct.SetPtr({{.Field.Slot.Offset}}, in.ToPtr())\n}\n\n{{end}}{{define \"structList\"}}// {{.Node.Name}}_List is a list of {{.Node.Name}}.\ntype {{.Node.Name}}_List struct{ {{.G.Capnp}}.List }\n\n// New{{.Node.Name}} creates a new list of {{.Node.Name}}.\nfunc New{{.Node.Name}}_List(s *{{.G.Capnp}}.Segment, sz int32) ({{.Node.Name}}_List, error) {\n\tl, err := {{.G.Capnp}}.NewCompositeList(s, {{.G.ObjectSize .Node}}, sz)\n\treturn {{.Node.Name}}_List{l}, err\n}\n\nfunc (s {{.Node.Name}}_List) At(i int) {{.Node.Name}} { return {{.Node.Name}}{ s.List.Struct(i) } }\n\nfunc (s {{.Node.Name}}_List) Set(i int, v {{.Node.Name}}) error { return s.List.SetStruct(i, v.Struct) }\n{{if .StringMethod}}\nfunc (s {{.Node.Name}}_List) String() string {\n\tstr, _ := {{.G.Imports.Text}}.MarshalList({{.Node.Id | printf \"%#x\"}}, s.List)\n\treturn str\n}\n{{end}}\n\n{{end}}{{define \"structListField\"}}func (s {{.Node.Name}}) {{.Field.Name | title}}() ({{.FieldType}}, error) {\n\t{{template \"_checktag\" .}}p, err := s.Struct.Ptr({{.Field.Slot.Offset}})\n\t{{if .Default.IsValid}}if err != nil {\n\t\treturn {{.FieldType}}{}, err\n\t}\n\tl, err := p.ListDefault({{.Default}})\n\treturn {{.FieldType}}{List: l}, err{{else}}return {{.FieldType}}{List: p.List()}, err{{end}}\n}\n\n{{template \"_hasfield\" .}}\n\nfunc (s {{.Node.Name}}) Set{{.Field.Name | title}}(v {{.FieldType}}) error {\n\t{{template \"_settag\" .}}return s.Struct.SetPtr({{.Field.Slot.Offset}}, v.List.ToPtr())\n}\n\n// New{{.Field.Name | title}} sets the {{.Field.Name}} field to a newly\n// allocated {{.FieldType}}, preferring placement in s's segment.\nfunc (s {{.Node.Name}}) New{{.Field.Name | title}}(n int32) ({{.FieldType}}, error) {\n\t{{template \"_settag\" .}}l, err := {{.G.RemoteTypeNew .Field.Slot.Type .Node}}(s.Struct.Segment(), n)\n\tif err != nil {\n\t\treturn {{.FieldType}}{}, err\n\t}\n\terr = s.Struct.SetPtr({{.Field.Slot.Offset}}, l.List.ToPtr())\n\treturn l, err\n}\n\n{{end}}{{define \"structPointerField\"}}func (s {{.Node.Name}}) {{.Field.Name | title}}() ({{.G.Capnp}}.Ptr, error) {\n\t{{template \"_checktag\" .}}{{if .Default.IsValid}}p, err := s.Struct.Ptr({{.Field.Slot.Offset}})\n\tif err != nil {\n\t\treturn nil, err\n\t}\n\treturn p.Default({{.Default}}){{else}}return s.Struct.Ptr({{.Field.Slot.Offset}}){{end}}\n}\n\n{{template \"_hasfield\" .}}\n\nfunc (s {{.Node.Name}}) Set{{.Field.Name | title}}(v {{.G.Capnp}}.Ptr) error {\n\t{{template \"_settag\" .}}return s.Struct.SetPtr({{.Field.Slot.Offset}}, v)\n}\n\n{{end}}{{define \"structStructField\"}}func (s {{.Node.Name}}) {{.Field.Name | title}}() ({{.FieldType}}, error) {\n\t{{template \"_checktag\" .}}p, err := s.Struct.Ptr({{.Field.Slot.Offset}})\n\t{{if .Default.IsValid}}if err != nil {\n\t\treturn {{.FieldType}}{}, err\n\t}\n\tss, err := p.StructDefault({{.Default}})\n\treturn {{.FieldType}}{Struct: ss}, err{{else}}return {{.FieldType}}{Struct: p.Struct()}, err{{end}}\n}\n\n{{template \"_hasfield\" .}}\n\nfunc (s {{.Node.Name}}) Set{{.Field.Name | title}}(v {{.FieldType}}) error {\n\t{{template \"_settag\" .}}return s.Struct.SetPtr({{.Field.Slot.Offset}}, v.Struct.ToPtr())\n}\n\n// New{{.Field.Name | title}} sets the {{.Field.Name}} field to a newly\n// allocated {{.FieldType}} struct, preferring placement in s's segment.\nfunc (s {{.Node.Name}}) New{{.Field.Name | title}}() ({{.FieldType}}, error) {\n\t{{template \"_settag\" .}}ss, err := {{.G.RemoteNodeNew .TypeNode .Node}}(s.Struct.Segment())\n\tif err != nil {\n\t\treturn {{.FieldType}}{}, err\n\t}\n\terr = s.Struct.SetPtr({{.Field.Slot.Offset}}, ss.Struct.ToPtr())\n\treturn ss, err\n}\n\n{{end}}{{define \"structTextField\"}}func (s {{.Node.Name}}) {{.Field.Name | title}}() (string, error) {\n\t{{template \"_checktag\" .}}p, err := s.Struct.Ptr({{.Field.Slot.Offset}})\n\t{{with .Default}}return p.TextDefault({{printf \"%q\" .}}), err{{else}}return p.Text(), err{{end}}\n}\n\n{{template \"_hasfield\" .}}\n\nfunc (s {{.Node.Name}}) {{.Field.Name | title}}Bytes() ([]byte, error) {\n\tp, err := s.Struct.Ptr({{.Field.Slot.Offset}})\n\t{{with .Default}}return p.TextBytesDefault({{printf \"%q\" .}}), err{{else}}return p.TextBytes(), err{{end}}\n}\n\nfunc (s {{.Node.Name}}) Set{{.Field.Name | title}}(v string) error {\n\t{{template \"_settag\" .}}{{if .Default}}return s.Struct.SetNewText({{.Field.Slot.Offset}}, v){{else}}return s.Struct.SetText({{.Field.Slot.Offset}}, v){{end}}\n}\n\n{{end}}{{define \"structTypes\"}}{{with .Annotations.Doc}}// {{.}}\n{{end}}type {{.Node.Name}} {{if .IsBase}}struct{ {{.G.Capnp}}.Struct }{{else}}{{.BaseNode.Name}}{{end}}\n{{end}}{{define \"structUintField\"}}func (s {{.Node.Name}}) {{.Field.Name | title}}() uint{{.Bits}} {\n\t{{template \"_checktag\" .}}return s.Struct.Uint{{.Bits}}({{.Offset}}){{with .Default}} ^ {{.}}{{end}}\n}\n\nfunc (s {{.Node.Name}}) Set{{.Field.Name | title}}(v uint{{.Bits}}) {\n\t{{template \"_settag\" .}}s.Struct.SetUint{{.Bits}}({{.Offset}}, v{{with .Default}}^{{.}}{{end}})\n}\n\n{{end}}{{define \"structValue\"}}{{.G.RemoteNodeName .Typ .Node}}{Struct: {{.G.Capnp}}.MustUnmarshalRoot({{.Value}}).Struct()}{{end}}{{define \"structVoidField\"}}{{if .Field.HasDiscriminant}}func (s {{.Node.Name}}) Set{{.Field.Name | title}}() {\n\t{{template \"_settag\" .}}\n}\n\n{{end}}{{end}}"))

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
