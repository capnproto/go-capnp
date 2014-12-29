package main

import (
	"strings"
	"text/template"
)

var templates = template.Must(template.New("").Funcs(template.FuncMap{
	"capn":  g_imports.capn,
	"title": strings.Title,
}).Parse(`
{{define "promise"}}
type {{.Node.Name}}_Promise struct {
	p {{capn}}.Promise
}

func New{{.Node.Name}}_Promise(p {{capn}}.Promise) {{.Node.Name}}_Promise {
	return {{.Node.Name}}_Promise{p}
}

func (p {{.Node.Name}}_Promise) Get() ({{.Node.Name}}, error) {
	s, err := p.p.Get()
	return {{.Node.Name}}(s), err
}

func (p {{.Node.Name}}_Promise) GenericPromise() {{capn}}.Promise { return p.p }
{{end}}


{{define "promiseFieldStruct"}}
func (p {{.Node.Name}}_Promise) {{.Field.Name|title}}() {{.Struct.RemoteName .Node}}_Promise {
	return {{.Struct.RemoteScope .Node}}New{{.Struct.Name}}_Promise(
		{{if .BufName}}p.p.GetPromiseDefault({{.Field.Slot.Offset}}, {{.BufName}}, {{.DefaultOffset}}){{else}}p.p.GetPromise({{.Field.Slot.Offset}}){{end}})
}
{{end}}


{{define "promiseFieldAnyPointer"}}
func (p {{.Node.Name}}_Promise) {{.Field.Name|title}}() {{capn}}.Promise {
	return p.p.GetPromise({{.Field.Slot.Offset}})
}

func (p {{.Node.Name}}_Promise) {{.Field.Name|title}}_Client() {{capn}}.Client {
	return p.p.GetClient({{.Field.Slot.Offset}})
}
{{end}}


{{define "promiseGroup"}}func (p {{.Node.Name}}_Promise) {{.Field.Name|title}}() {{.Group.Name}}_Promise { return {{.Group.Name}}_Promise(p) }
{{end}}
`))

type promiseTemplateParams struct {
	Node   *node
	Fields []Field
}

type promiseGroupTemplateParams struct {
	Node  *node
	Field Field
	Group *node
}

type promiseFieldStructTemplateParams struct {
	Node          *node
	Field         Field
	Struct        *node
	BufName       string
	DefaultOffset int
}

type promiseFieldAnyPointerTemplateParams struct {
	Node  *node
	Field Field
}
