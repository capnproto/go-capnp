package main

import (
	"strings"
	"text/template"
)

var templates = template.Must(template.New("").Funcs(template.FuncMap{
	"capn":    g_imports.capn,
	"context": g_imports.context,
	"title":   strings.Title,
}).Parse(`
{{define "promise"}}
type {{.Node.Name}}_Promise {{capn}}.Promise

func (p *{{.Node.Name}}_Promise) Get() ({{.Node.Name}}, error) {
	s, err := (*{{capn}}.Promise)(p).Struct()
	return {{.Node.Name}}(s), err
}
{{end}}


{{define "promiseFieldStruct"}}
func (p *{{.Node.Name}}_Promise) {{.Field.Name|title}}() *{{.Struct.RemoteName .Node}}_Promise {
	return (*{{.Struct.RemoteName .Node}}_Promise)((*{{capn}}.Promise)(p).{{if .BufName}}GetPromiseDefault({{.Field.Slot.Offset}}, {{.BufName}}, {{.DefaultOffset}}){{else}}GetPromise({{.Field.Slot.Offset}}){{end}})
}
{{end}}


{{define "promiseFieldAnyPointer"}}
func (p *{{.Node.Name}}_Promise) {{.Field.Name|title}}() *{{capn}}.Promise {
	return (*{{capn}}.Promise)(p).GetPromise({{.Field.Slot.Offset}})
}
{{end}}


{{define "promiseFieldInterface"}}
func (p *{{.Node.Name}}_Promise) {{.Field.Name|title}}() {{.Interface.RemoteName .Node}} {
	return {{.Interface.RemoteNew .Node}}((*{{capn}}.Promise)(p).GetPromise({{.Field.Slot.Offset}}).Client())
}
{{end}}


{{define "promiseGroup"}}func (p *{{.Node.Name}}_Promise) {{.Field.Name|title}}() *{{.Group.Name}}_Promise { return (*{{.Group.Name}}_Promise)(p) }
{{end}}


{{define "interfaceClient"}}{{with .Annotations.Doc}}// {{.}}
{{end}}type {{.Node.Name}} struct { c {{capn}}.Client }

func New{{.Node.Name}}(c {{capn}}.Client) {{.Node.Name}} { return {{.Node.Name}}{c} }

func (c {{.Node.Name}}) GenericClient() {{capn}}.Client { return c.c }

func (c {{.Node.Name}}) IsNull() bool { return c.c == nil }

var clientMethods_{{.Node.Name}} = []{{capn}}.Method{
	{{range .Methods}}
	{
		{{template "_interfaceMethod" .}}
	},
	{{end}}
}

{{range $i, $m := .Methods}}
func (c {{$.Node.Name}}) {{.Name|title}}(ctx {{context}}.Context, params func({{.Params.RemoteName $.Node}})) *{{.Results.RemoteName $.Node}}_Promise {
	if c.c == nil {
		return (*{{.Results.RemoteName $.Node}}_Promise)({{capn}}.NewPromise({{capn}}.ErrorAnswer({{capn}}.ErrNullClient)))
	}
	return (*{{.Results.RemoteName $.Node}}_Promise)({{capn}}.NewPromise(c.c.NewCall(ctx,
		&clientMethods_{{$.Node.Name}}[{{$i}}],
		{{.Params.ObjectSize}},
		func(s {{capn}}.Struct) { params({{.Params.RemoteName $.Node}}(s)) })))
}
{{end}}
{{end}}


{{define "interfaceServer"}}type {{.Node.Name}}_Server interface {
	{{range .Methods}}
	{{.Name|title}}(ctx {{context}}.Context, params {{.Params.RemoteName $.Node}}, results {{.Results.RemoteName $.Node}}) error
	{{end}}
}

func {{.Node.Name}}_Methods(methods []{{capn}}.ServerMethod, server {{.Node.Name}}_Server) []{{capn}}.ServerMethod {
	if cap(methods) == 0 {
		methods = make([]{{capn}}.ServerMethod, 0, {{len .Methods}})
	}
	{{range .Methods}}
	methods = append(methods, {{capn}}.ServerMethod{
		Method: {{capn}}.Method{
			{{template "_interfaceMethod" .}}
		},
		Impl: func(c {{context}}.Context, p, r {{capn}}.Struct) error {
			return server.{{.Name|title}}(c, {{.Params.RemoteName $.Node}}(p), {{.Results.RemoteName $.Node}}(r))
		},
		ResultsSize: {{.Results.ObjectSize}},
	})
	{{end}}
	return methods
}
{{end}}


{{define "_interfaceMethod"}}
			InterfaceID: {{.Interface.Id|printf "%#x"}},
			MethodID: {{.ID}},
			InterfaceName: {{.Interface.DisplayName|printf "%q"}},
			MethodName: {{.Name|printf "%q"}},
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

type promiseFieldInterfaceTemplateParams struct {
	Node      *node
	Field     Field
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
