package main

import (
	"fmt"
	"strings"
	"text/template"

	C "zombiezen.com/go/capnproto"
)

var templates = template.Must(template.New("").Funcs(template.FuncMap{
	"capnp":   g_imports.capnp,
	"math":    g_imports.math,
	"server":  g_imports.server,
	"context": g_imports.context,
	"strconv": g_imports.strconv,
	"title":   strings.Title,
	"hasDiscriminant": func(f field) bool {
		return f.DiscriminantValue() != Field_noDiscriminant
	},
	"discriminantOffset": func(n *node) uint32 {
		return n.StructGroup().DiscriminantOffset() * 2
	},
}).Parse(`
{{define "enum"}}{{with .Annotations.Doc}}// {{.}}
{{end}}type {{.Node.Name}} uint16

{{with .EnumValues}}
// Values of {{$.Node.Name}}.
const (
{{range .}}{{.FullName}} {{$.Node.Name}} = {{.Val}}
{{end}}
)

// String returns the enum's constant name.
func (c {{$.Node.Name}}) String() string {
	switch c {
	{{range .}}{{if .Tag}}case {{.FullName}}: return {{printf "%q" .Tag}}
	{{end}}{{end}}
	default: return ""
	}
}

// {{$.Node.Name}}FromString returns the enum value with a name,
// or the zero value if there's no such value.
func {{$.Node.Name}}FromString(c string) {{$.Node.Name}} {
	switch c {
	{{range .}}{{if .Tag}}case {{printf "%q" .Tag}}: return {{.FullName}}
	{{end}}{{end}}
	default: return 0
	}
}
{{end}}

type {{.Node.Name}}_List struct { {{capnp}}.List }

func New{{.Node.Name}}_List(s *{{capnp}}.Segment, sz int32) ({{.Node.Name}}_List, error) {
	l, err := {{capnp}}.NewUInt16List(s, sz)
	if err != nil {
		return {{.Node.Name}}_List{}, err
	}
	return {{.Node.Name}}_List{l.List}, nil
}

func (l {{.Node.Name}}_List) At(i int) {{.Node.Name}} {
	ul := {{capnp}}.UInt16List{List: l.List}
	return {{.Node.Name}}(ul.At(i))
}

func (l {{.Node.Name}}_List) Set(i int, v {{.Node.Name}}) {
	ul := {{capnp}}.UInt16List{List: l.List}
	ul.Set(i, uint16(v))
}
{{end}}


{{define "structTypes"}}{{with .Annotations.Doc}}// {{.}}
{{end}}type {{.Node.Name}} {{if .IsBase}}struct{ {{capnp}}.Struct }{{else}}{{.BaseNode.Name}}{{end}}
{{end}}


{{define "newStructFunc"}}
func New{{.Node.Name}}(s *{{capnp}}.Segment) ({{.Node.Name}}, error) {
	st, err := {{capnp}}.NewStruct(s, {{.Node.ObjectSize}})
	if err != nil {
		return {{.Node.Name}}{}, err
	}
	return {{.Node.Name}}{st}, nil
}

func NewRoot{{.Node.Name}}(s *{{capnp}}.Segment) ({{.Node.Name}}, error) {
	st, err := {{capnp}}.NewRootStruct(s, {{.Node.ObjectSize}})
	if err != nil {
		return {{.Node.Name}}{}, err
	}
	return {{.Node.Name}}{st}, nil
}

func ReadRoot{{.Node.Name}}(msg *{{capnp}}.Message) ({{.Node.Name}}, error) {
	root, err := msg.Root()
	if err != nil {
		return {{.Node.Name}}{}, err
	}
	st := {{capnp}}.ToStruct(root)
	return {{.Node.Name}}{st}, nil
}
{{end}}


{{define "structFuncs"}}
{{if gt .Node.StructGroup.DiscriminantCount 0}}
func (s {{.Node.Name}}) Which() {{.Node.Name}}_Which {
	return {{.Node.Name}}_Which(s.Struct.Uint16({{discriminantOffset .Node}}))
}
{{end}}
{{end}}


{{define "settag"}}{{if hasDiscriminant .Field}}s.Struct.SetUint16({{discriminantOffset .Node}}, {{.Field.DiscriminantValue}}){{end}}{{end}}


{{define "structVoidField"}}{{if hasDiscriminant .Field}}
func (s {{.Node.Name}}) Set{{.Field.Name|title}}() {
	{{template "settag" .}}
}
{{end}}{{end}}


{{define "structBoolField"}}
func (s {{.Node.Name}}) {{.Field.Name|title}}() bool {
	return {{if .Default}}!{{end}}s.Struct.Bit({{.Field.Slot.Offset}})
}

func (s {{.Node.Name}}) Set{{.Field.Name|title}}(v bool) {
	{{template "settag" .}}
	s.Struct.SetBit({{.Field.Slot.Offset}}, {{if .Default}}!{{end}}v)
}
{{end}}


{{define "structUintField"}}
func (s {{.Node.Name}}) {{.Field.Name|title}}() uint{{.Bits}} {
	return s.Struct.Uint{{.Bits}}({{.Offset}}){{with .Default}} ^ {{.}}{{end}}
}

func (s {{.Node.Name}}) Set{{.Field.Name|title}}(v uint{{.Bits}}) {
	{{template "settag" .}}
	s.Struct.SetUint{{.Bits}}({{.Offset}}, v{{with .Default}}^{{.}}{{end}})
}
{{end}}


{{define "structIntField"}}
func (s {{.Node.Name}}) {{.Field.Name|title}}() {{.ReturnType}} {
	return {{.ReturnType}}(s.Struct.Uint{{.Bits}}({{.Offset}}){{with .Default}} ^ {{.}}{{end}})
}

func (s {{.Node.Name}}) Set{{.Field.Name|title}}(v {{.ReturnType}}) {
	{{template "settag" .}}
	s.Struct.SetUint{{.Bits}}({{.Offset}}, uint{{.Bits}}(v{{with .Default}}^{{.}}{{end}}))
}
{{end}}


{{define "structFloatField"}}
func (s {{.Node.Name}}) {{.Field.Name|title}}() float{{.Bits}} {
	return {{math}}.Float{{.Bits}}frombits(s.Struct.Uint{{.Bits}}({{.Offset}}){{with .Default}} ^ {{printf "%#x" .}}{{end}})
}

func (s {{.Node.Name}}) Set{{.Field.Name|title}}(v float{{.Bits}}) {
	{{template "settag" .}}
	s.Struct.SetUint{{.Bits}}({{.Offset}}, {{math}}.Float{{.Bits}}bits(v){{with .Default}}^{{printf "%#x" .}}{{end}})
}
{{end}}


{{define "structTextField"}}
func (s {{.Node.Name}}) {{.Field.Name|title}}() (string, error) {
	p, err := s.Struct.Pointer({{.Field.Slot.Offset}})
	if err != nil {
		return "", err
	}
	{{with .Default}}
	return {{capnp}}.ToTextDefault(p, {{printf "%q" .}})
	{{else}}
	return {{capnp}}.ToText(p), nil
	{{end}}
}

func (s {{.Node.Name}}) Set{{.Field.Name|title}}(v string) error {
	{{template "settag" .}}
	t, err := {{capnp}}.NewText(s.Struct.Segment(), v)
	if err != nil {
		return err
	}
	return s.Struct.SetPointer({{.Field.Slot.Offset}}, t)
}
{{end}}


{{define "structDataField"}}
func (s {{.Node.Name}}) {{.Field.Name|title}}() ({{.FieldType}}, error) {
	p, err := s.Struct.Pointer({{.Field.Slot.Offset}})
	if err != nil {
		return nil, err
	}
	{{with .Default}}
	v, err := {{capnp}}.ToDataDefault(p, {{printf "%#v" .}})
	return {{.FieldType}}(v), err
	{{else}}
	return {{.FieldType}}({{capnp}}.ToData(p)), nil
	{{end}}
}

func (s {{.Node.Name}}) Set{{.Field.Name|title}}(v {{.FieldType}}) error {
	{{template "settag" .}}
	d, err := {{capnp}}.NewData(s.Struct.Segment(), []byte(v))
	if err != nil {
		return err
	}
	return s.Struct.SetPointer({{.Field.Slot.Offset}}, d)
}
{{end}}


{{define "structStructField"}}
func (s {{.Node.Name}}) {{.Field.Name|title}}() ({{.FieldType}}, error) {
	p, err := s.Struct.Pointer({{.Field.Slot.Offset}})
	if err != nil {
		return {{.FieldType}}{}, err
	}
	{{if .Default.IsValid}}
	ss, err := {{capnp}}.ToStructDefault(p, {{.Default}})
	if err != nil {
		return {{.FieldType}}{}, err
	}
	{{else}}
	ss := {{capnp}}.ToStruct(p)
	{{end}}
	return {{.FieldType}}{Struct: ss}, nil
}

func (s {{.Node.Name}}) Set{{.Field.Name|title}}(v {{.FieldType}}) error {
	{{template "settag" .}}
	return s.Struct.SetPointer({{.Field.Slot.Offset}}, v.Struct)
}

// New{{.Field.Name|title}} sets the {{.Field.Name}} field to a newly
// allocated {{.FieldType}} struct, preferring placement in s's segment.
func (s {{.Node.Name}}) New{{.Field.Name|title}}() ({{.FieldType}}, error) {
	{{template "settag" .}}
	ss, err := {{.TypeNode.RemoteNew .Node}}(s.Struct.Segment())
	if err != nil {
		return {{.FieldType}}{}, err
	}
	err = s.Struct.SetPointer({{.Field.Slot.Offset}}, ss)
	return ss, err
}
{{end}}


{{define "structPointerField"}}
func (s {{.Node.Name}}) {{.Field.Name|title}}() ({{capnp}}.Pointer, error) {
	{{if .Default.IsValid}}
	p, err := s.Struct.Pointer({{.Field.Slot.Offset}})
	if err != nil {
		return nil, err
	}
	return {{capnp}}.PointerDefault(p, {{.Default}})
	{{else}}
	return s.Struct.Pointer({{.Field.Slot.Offset}})
	{{end}}
}

func (s {{.Node.Name}}) Set{{.Field.Name|title}}(v {{capnp}}.Pointer) error {
	{{template "settag" .}}
	return s.Struct.SetPointer({{.Field.Slot.Offset}}, v)
}
{{end}}


{{define "structListField"}}
func (s {{.Node.Name}}) {{.Field.Name|title}}() ({{.FieldType}}, error) {
	p, err := s.Struct.Pointer({{.Field.Slot.Offset}})
	if err != nil {
		return {{.FieldType}}{}, err
	}
	{{if .Default.IsValid}}
	l, err := {{capnp}}.ToListDefault(p, {{.Default}})
	if err != nil {
		return {{.FieldType}}{}, err
	}
	{{else}}
	l := {{capnp}}.ToList(p)
	{{end}}
	return {{.FieldType}}{List: l}, nil
}

func (s {{.Node.Name}}) Set{{.Field.Name|title}}(v {{.FieldType}}) error {
	{{template "settag" .}}
	return s.Struct.SetPointer({{.Field.Slot.Offset}}, v.List)
}
{{end}}


{{define "structInterfaceField"}}
func (s {{.Node.Name}}) {{.Field.Name|title}}() {{.FieldType}} {
	p, err := s.Struct.Pointer({{.Field.Slot.Offset}})
	if err != nil {
		{{/* Valid interface pointers never return errors. */}}
		return {{.FieldType}}{}
	}
	c := {{capnp}}.ToInterface(p).Client()
	return {{.FieldType}}{Client: c}
}

func (s {{.Node.Name}}) Set{{.Field.Name|title}}(v {{.FieldType}}) error {
	{{template "settag" .}}
	seg := s.Segment()
	if seg == nil {
		{{/* TODO(light): error? */}}
		return nil
	}
	ci := seg.Message().AddCap(v.Client)
	return s.Struct.SetPointer({{.Field.Slot.Offset}}, {{capnp}}.NewInterface(seg, ci))
}
{{end}}


{{define "structList"}}// {{.Node.Name}}_List is a list of {{.Node.Name}}.
type {{.Node.Name}}_List struct{ {{capnp}}.List }

// New{{.Node.Name}} creates a new list of {{.Node.Name}}.
func New{{.Node.Name}}_List(s *{{capnp}}.Segment, sz int32) ({{.Node.Name}}_List, error) {
	l, err := {{capnp}}.NewCompositeList(s, {{.Node.ObjectSize}}, sz)
	if err != nil  {
		return {{.Node.Name}}_List{}, err
	}
	return {{.Node.Name}}_List{l}, nil
}

func (s {{.Node.Name}}_List) At(i int) {{.Node.Name}} { return {{.Node.Name}}{ s.List.Struct(i) } }
func (s {{.Node.Name}}_List) Set(i int, v {{.Node.Name}}) error { return s.List.SetStruct(i, v.Struct) }
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


{{define "promise"}}// {{.Node.Name}}_Promise is a wrapper for a {{.Node.Name}} promised by a client call.
type {{.Node.Name}}_Promise struct { *{{capnp}}.Pipeline }

func (p {{.Node.Name}}_Promise) Struct() ({{.Node.Name}}, error) {
	s, err := p.Pipeline.Struct()
	return {{.Node.Name}}{s}, err
}
{{end}}


{{define "promiseFieldStruct"}}
func (p {{.Node.Name}}_Promise) {{.Field.Name|title}}() {{.Struct.RemoteName .Node}}_Promise {
	return {{.Struct.RemoteName .Node}}_Promise{Pipeline: p.Pipeline.{{if .Default.IsValid}}GetPipelineDefault({{.Field.Slot.Offset}}, {{.Default}}){{else}}GetPipeline({{.Field.Slot.Offset}}){{end}} }
}
{{end}}


{{define "promiseFieldAnyPointer"}}
func (p {{.Node.Name}}_Promise) {{.Field.Name|title}}() *{{capnp}}.Pipeline {
	return p.Pipeline.GetPipeline({{.Field.Slot.Offset}})
}
{{end}}


{{define "promiseFieldInterface"}}
func (p {{.Node.Name}}_Promise) {{.Field.Name|title}}() {{.Interface.RemoteName .Node}} {
	return {{.Interface.RemoteName .Node}}{Client: p.Pipeline.GetPipeline({{.Field.Slot.Offset}}).Client()}
}
{{end}}


{{define "promiseGroup"}}func (p {{.Node.Name}}_Promise) {{.Field.Name|title}}() {{.Group.Name}}_Promise { return {{.Group.Name}}_Promise{p.Pipeline} }
{{end}}


{{define "interfaceClient"}}{{with .Annotations.Doc}}// {{.}}
{{end}}type {{.Node.Name}} struct { {{capnp}}.Client }

{{range .Methods}}
func (c {{$.Node.Name}}) {{.Name|title}}(ctx {{context}}.Context, params func({{.Params.RemoteName $.Node}}) error, opts ...{{capnp}}.CallOption) {{.Results.RemoteName $.Node}}_Promise {
	if c.Client == nil {
		return {{.Results.RemoteName $.Node}}_Promise{Pipeline: {{capnp}}.NewPipeline({{capnp}}.ErrorAnswer({{capnp}}.ErrNullClient))}
	}
	return {{.Results.RemoteName $.Node}}_Promise{Pipeline: {{capnp}}.NewPipeline(c.Client.Call(&{{capnp}}.Call{
		Ctx: ctx,
		Method: {{capnp}}.Method{
			{{template "_interfaceMethod" .}}
		},
		ParamsSize: {{.Params.ObjectSize}},
		ParamsFunc: func(s {{capnp}}.Struct) error { return params({{.Params.RemoteName $.Node}}{Struct: s}) },
		Options: {{capnp}}.NewCallOptions(opts),
	}))}
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
	return {{.Node.Name}}{Client: {{server}}.New({{.Node.Name}}_Methods(nil, s), c)}
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
			call := {{.Interface.RemoteName $.Node}}_{{.Name}}{c, opts, {{.Params.RemoteName $.Node}}{Struct: p}, {{.Results.RemoteName $.Node}}{Struct: r} }
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
	Bits    int
	Default uint64
}

func (p structUintFieldParams) Offset() C.DataOffset {
	return C.DataOffset(p.Field.Slot().Offset()) * C.DataOffset(p.Bits/8)
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
