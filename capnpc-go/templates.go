package main

import (
	"strings"
	"text/template"

	C "zombiezen.com/go/capnproto"
)

var templates = template.Must(template.New("").Funcs(template.FuncMap{
	"capnp":   g_imports.capnp,
	"server":  g_imports.server,
	"context": g_imports.context,
	"strconv": g_imports.strconv,
	"title":   strings.Title,
}).Parse(`
{{define "structFuncs"}}
func (s {{.Node.Name}}) Segment() *{{capnp}}.Segment {
	return {{capnp}}.Struct(s).Segment()
}
{{if gt .Node.Struct.DiscriminantCount 0}}
func (s {{.Node.Name}}) Which() {{.Node.Name}}_Which {
	return {{.Node.Name}}_Which({{capnp}}.Struct(s).Uint16({{.DiscriminantOffset}}))
}
{{end}}
{{end}}


{{define "structEnums"}}type {{.Node.Name}}_Which uint16

const (
{{range .Fields}}	{{$.Node.Name}}_Which_{{.Name}} {{$.Node.Name}}_Which = {{.DiscriminantValue}}
{{end}}
)

func (w {{.Node.Name}}_Which) String() string {
	const s = {{.EnumString.ValueString|printf "%q"}}
	switch w {
	{{range $i, $f := .Fields}}case {{$.Node.Name}}_Which_{{.Name}}:
		return s{{$.EnumString.SliceFor $i}}
	{{end}}
	}
	return "{{.Node.Name}}_Which(" + {{strconv}}.FormatUint(uint64(w), 10) + ")"
}

{{end}}


{{define "annotation"}}const {{.Node.Name}} = uint64({{.Node.Id|printf "%#x"}})
{{end}}


{{define "promise"}}
type {{.Node.Name}}_Promise {{capnp}}.Pipeline

func (p *{{.Node.Name}}_Promise) Get() ({{.Node.Name}}, error) {
	s, err := (*{{capnp}}.Pipeline)(p).Struct()
	return {{.Node.Name}}(s), err
}
{{end}}


{{define "promiseFieldStruct"}}
func (p *{{.Node.Name}}_Promise) {{.Field.Name|title}}() *{{.Struct.RemoteName .Node}}_Promise {
	return (*{{.Struct.RemoteName .Node}}_Promise)((*{{capnp}}.Pipeline)(p).{{if .BufName}}GetPipelineDefault({{.Field.Slot.Offset}}, {{.BufName}}, {{.DefaultOffset}}){{else}}GetPipeline({{.Field.Slot.Offset}}){{end}})
}
{{end}}


{{define "promiseFieldAnyPointer"}}
func (p *{{.Node.Name}}_Promise) {{.Field.Name|title}}() *{{capnp}}.Pipeline {
	return (*{{capnp}}.Pipeline)(p).GetPipeline({{.Field.Slot.Offset}})
}
{{end}}


{{define "promiseFieldInterface"}}
func (p *{{.Node.Name}}_Promise) {{.Field.Name|title}}() {{.Interface.RemoteName .Node}} {
	return {{.Interface.RemoteNew .Node}}((*{{capnp}}.Pipeline)(p).GetPipeline({{.Field.Slot.Offset}}).Client())
}
{{end}}


{{define "promiseGroup"}}func (p *{{.Node.Name}}_Promise) {{.Field.Name|title}}() *{{.Group.Name}}_Promise { return (*{{.Group.Name}}_Promise)(p) }
{{end}}


{{define "interfaceClient"}}{{with .Annotations.Doc}}// {{.}}
{{end}}type {{.Node.Name}} struct { c {{capnp}}.Client }

func New{{.Node.Name}}(c {{capnp}}.Client) {{.Node.Name}} { return {{.Node.Name}}{c} }

func (c {{.Node.Name}}) GenericClient() {{capnp}}.Client { return c.c }

func (c {{.Node.Name}}) IsNull() bool { return c.c == nil }

{{range .Methods}}
func (c {{$.Node.Name}}) {{.Name|title}}(ctx {{context}}.Context, params func({{.Params.RemoteName $.Node}}), opts ...{{capnp}}.CallOption) *{{.Results.RemoteName $.Node}}_Promise {
	if c.c == nil {
		return (*{{.Results.RemoteName $.Node}}_Promise)({{capnp}}.NewPipeline({{capnp}}.ErrorAnswer({{capnp}}.ErrNullClient)))
	}
	return (*{{.Results.RemoteName $.Node}}_Promise)({{capnp}}.NewPipeline(c.c.Call(&{{capnp}}.Call{
		Ctx: ctx,
		Method: {{capnp}}.Method{
			{{template "_interfaceMethod" .}}
		},
		ParamsSize: {{.Params.ObjectSize}},
		ParamsFunc: func(s {{capnp}}.Struct) { params({{.Params.RemoteName $.Node}}(s)) },
		Options: {{capnp}}.NewCallOptions(opts),
	})))
}
{{end}}
{{end}}


{{define "interfaceServer"}}type {{.Node.Name}}_Server interface {
	{{range .Methods}}
	{{.Name|title}}({{.Interface.RemoteName $.Node}}_{{.Name}}) error
	{{end}}
}

func {{.Node.Name}}_ServerToClient(s {{.Node.Name}}_Server) {{.Node.Name}} {
	c, _ := s.({{server}}.Closer)
	return New{{.Node.Name}}({{server}}.New({{.Node.Name}}_Methods(nil, s), c))
}

func {{.Node.Name}}_Methods(methods []{{server}}.Method, s {{.Node.Name}}_Server) []{{server}}.Method {
	if cap(methods) == 0 {
		methods = make([]{{server}}.Method, 0, {{len .Methods}})
	}
	{{range .Methods}}
	methods = append(methods, {{server}}.Method{
		Method: {{capnp}}.Method{
			{{template "_interfaceMethod" .}}
		},
		Impl: func(c {{context}}.Context, opts {{capnp}}.CallOptions, p, r {{capnp}}.Struct) error {
			call := {{.Interface.RemoteName $.Node}}_{{.Name}}{c, opts, {{.Params.RemoteName $.Node}}(p), {{.Results.RemoteName $.Node}}(r)}
			return s.{{.Name|title}}(call)
		},
		ResultsSize: {{.Results.ObjectSize}},
	})
	{{end}}
	return methods
}
{{range .Methods}}{{if eq .Interface.Id $.Node.Id}}
// {{$.Node.Name}}_{{.Name}} holds the arguments for a server call to {{$.Node.Name}}.{{.Name}}.
type {{$.Node.Name}}_{{.Name}} struct {
	Ctx     {{context}}.Context
	Options {{capnp}}.CallOptions
	Params  {{.Params.RemoteName $.Node}}
	Results {{.Results.RemoteName $.Node}}
}
{{end}}{{end}}
{{end}}


{{define "_interfaceMethod"}}
			InterfaceID: {{.Interface.Id|printf "%#x"}},
			MethodID: {{.ID}},
			InterfaceName: {{.Interface.DisplayName|printf "%q"}},
			MethodName: {{.OriginalName|printf "%q"}},
{{end}}
`))

type annotationParams struct {
	Node *node
}

type structFuncsParams struct {
	Node               *node
	DiscriminantOffset C.DataOffset
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
	Node          *node
	Field         field
	Struct        *node
	BufName       string
	DefaultOffset C.Address
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
