/*
capnpc-go is the Cap'n proto code generator for Go.  It reads a
CodeGeneratorRequest from stdin and for a file foo.capnp it writes
foo.capnp.go.  This is usually invoked from `capnp compile -ogo`.

See https://capnproto.org/otherlang.html#how-to-write-compiler-plugins
for more details.
*/
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"zombiezen.com/go/capnproto2"
	"zombiezen.com/go/capnproto2/internal/schema"
)

var (
	genPromises = flag.Bool("promises", true, "generate code for promises")
)

// Non-stdlib import paths.
const (
	capnpImport   = "zombiezen.com/go/capnproto2"
	serverImport  = capnpImport + "/server"
	contextImport = "golang.org/x/net/context"
)

type nodeMap map[uint64]*node

func (nm nodeMap) mustFind(id uint64) (*node, error) {
	n := nm[id]
	if n == nil {
		return nil, fmt.Errorf("could not find node %#x in schema", id)
	}
	return n, nil
}

type generator struct {
	buf     bytes.Buffer
	fileID  uint64
	nodes   nodeMap
	imports imports
	data    staticData
}

func newGenerator(fileID uint64, nodes nodeMap) *generator {
	g := &generator{
		fileID: fileID,
		nodes:  nodes,
	}
	g.imports.init()
	g.data.init(fileID)
	return g
}

// Basename returns the name of the schema file with the directory name removed.
func (g *generator) Basename() (string, error) {
	f, err := g.nodes.mustFind(g.fileID)
	if err != nil {
		return "", err
	}
	dn, err := f.DisplayName()
	if err != nil {
		return "", err
	}
	return filepath.Base(dn), nil
}

func (g *generator) Imports() *imports {
	return &g.imports
}

func (g *generator) Capnp() string {
	return g.imports.Capnp()
}

// generate produces unformatted Go source code from the nodes defined in it.
func (g *generator) generate(pkg string) []byte {
	var out bytes.Buffer
	fmt.Fprintf(&out, "package %s\n\n", pkg)
	out.WriteString("// AUTO GENERATED - DO NOT EDIT\n\n")
	out.WriteString("import (\n")
	for _, imp := range g.imports.usedImports() {
		fmt.Fprintf(&out, "%v\n", imp)
	}
	out.WriteString(")\n")
	out.Write(g.buf.Bytes())
	if len(g.data.buf) > 0 {
		fmt.Fprintf(&out, "var %s = []byte{", g.data.name)
		for i, b := range g.data.buf {
			if i%8 == 0 {
				out.WriteByte('\n')
			}
			fmt.Fprintf(&out, "%d,", b)
		}
		fmt.Fprintf(&out, "\n}\n")
	}
	return out.Bytes()
}

func (g *generator) appendRemoteScope(b []byte, n, from *node) ([]byte, error) {
	if n.pkg == "" {
		return b, fmt.Errorf("internal error (bad schema?): missing package declaration for %s", n)
	}
	if n.imp == "" {
		return b, fmt.Errorf("internal error (bad schema?): missing import declaration for %s", n)
	}
	if from.imp == "" {
		return b, fmt.Errorf("internal error (bad schema?): missing import declaration for %s", from)
	}

	if n.imp == from.imp {
		return b, nil
	}
	name := g.imports.add(importSpec{path: n.imp, name: n.pkg})
	b = append(b, name...)
	b = append(b, '.')
	return b, nil
}

func (g *generator) RemoteNew(n, from *node) (string, error) {
	buf, err := g.appendRemoteScope(nil, n, from)
	if err != nil {
		return "", err
	}
	buf = append(buf, "New"...)
	buf = append(buf, n.Name...)
	return string(buf), nil
}

func (g *generator) RemoteName(n, from *node) (string, error) {
	buf, err := g.appendRemoteScope(nil, n, from)
	if err != nil {
		return "", err
	}
	buf = append(buf, n.Name...)
	return string(buf), nil
}

var templates = template.New("").Funcs(template.FuncMap{
	"title": strings.Title,
	"hasDiscriminant": func(f field) bool {
		return f.DiscriminantValue() != schema.Field_noDiscriminant
	},
	"discriminantOffset": func(n *node) uint32 {
		return n.StructGroup().DiscriminantOffset() * 2
	},
})

func displayName(n interface {
	DisplayName() (string, error)
}) string {
	dn, _ := n.DisplayName()
	return dn
}

type node struct {
	schema.Node
	pkg   string
	imp   string
	nodes []*node
	Name  string
}

func (n *node) shortDisplayName() string {
	dn, _ := n.DisplayName()
	return dn[n.DisplayNamePrefixLength():]
}

// String returns the node's display name.
func (n *node) String() string {
	return displayName(n)
}

type field struct {
	schema.Field
	Name string
}

// Tag types
const (
	defaultTag = iota
	noTag
	customTag
)

type annotations struct {
	Doc       string
	Package   string
	Import    string
	TagType   int
	CustomTag string
	Name      string
}

func parseAnnotations(list schema.Annotation_List) *annotations {
	ann := new(annotations)
	for i, n := 0, list.Len(); i < n; i++ {
		a := list.At(i)
		val, _ := a.Value()
		text, _ := val.Text()
		switch a.Id() {
		case capnp.Doc:
			ann.Doc = text
		case capnp.Package:
			ann.Package = text
		case capnp.Import:
			ann.Import = text
		case capnp.Tag:
			ann.TagType = customTag
			ann.CustomTag = text
		case capnp.Notag:
			ann.TagType = noTag
		case capnp.Name:
			ann.Name = text
		}
	}
	return ann
}

// Tag returns the string value that an enumerant value called name should have.
// An empty string indicates that this enumerant value has no tag.
func (ann *annotations) Tag(name string) string {
	switch ann.TagType {
	case noTag:
		return ""
	case customTag:
		return ann.CustomTag
	case defaultTag:
		fallthrough
	default:
		return name
	}
}

// Rename returns the overridden name from the annotations or the given name
// if no annotation was found.
func (ann *annotations) Rename(given string) string {
	if ann.Name == "" {
		return given
	}
	return ann.Name
}

func buildNodeMap(req schema.CodeGeneratorRequest) (nodeMap, error) {
	rnodes, err := req.Nodes()
	if err != nil {
		return nil, err
	}
	nodes := make(nodeMap, rnodes.Len())
	var allfiles []*node
	for i := 0; i < rnodes.Len(); i++ {
		ni := rnodes.At(i)
		n := &node{Node: ni}
		nodes[n.Id()] = n
		if n.Which() == schema.Node_Which_file {
			allfiles = append(allfiles, n)
		}
	}
	for _, f := range allfiles {
		fann, err := f.Annotations()
		if err != nil {
			return nil, fmt.Errorf("reading annotations for %v: %v", f, err)
		}
		ann := parseAnnotations(fann)
		f.pkg = ann.Package
		f.imp = ann.Import
		nnodes, _ := f.NestedNodes()
		for i := 0; i < nnodes.Len(); i++ {
			nn := nnodes.At(i)
			if ni := nodes[nn.Id()]; ni != nil {
				nname, _ := nn.Name()
				if err := resolveName(nodes, ni, "", nname, f); err != nil {
					return nil, err
				}
			}
		}
	}
	return nodes, nil
}

// resolveName is called as part of building up a node map to populate the name field of n.
func resolveName(nodes nodeMap, n *node, base, name string, file *node) error {
	na, err := n.Annotations()
	if err != nil {
		return fmt.Errorf("reading annotations for %s: %v", n, err)
	}
	name = parseAnnotations(na).Rename(name)
	if base == "" {
		n.Name = strings.Title(name)
	} else {
		n.Name = base + "_" + name
	}
	n.pkg = file.pkg
	n.imp = file.imp

	if n.Which() != schema.Node_Which_structGroup || !n.StructGroup().IsGroup() {
		file.nodes = append(file.nodes, n)
	}

	nnodes, err := n.NestedNodes()
	if err != nil {
		return fmt.Errorf("listing nested nodes for %s: %v", n, err)
	}
	for i := 0; i < nnodes.Len(); i++ {
		nn := nnodes.At(i)
		ni := nodes[nn.Id()]
		if ni == nil {
			continue
		}
		nname, err := nn.Name()
		if err != nil {
			return fmt.Errorf("reading name of nested node %d in %s: %v", i+1, n, err)
		}
		if err := resolveName(nodes, ni, n.Name, nname, file); err != nil {
			return err
		}
	}

	switch n.Which() {
	case schema.Node_Which_structGroup:
		fields, _ := n.StructGroup().Fields()
		for i := 0; i < fields.Len(); i++ {
			f := fields.At(i)
			if f.Which() != schema.Field_Which_group {
				continue
			}
			fa, _ := f.Annotations()
			fname, _ := f.Name()
			grp := nodes[f.Group().TypeId()]
			if grp == nil {
				return fmt.Errorf("could not find type information for group %s in %s", fname, n)
			}
			fname = parseAnnotations(fa).Rename(fname)
			if err := resolveName(nodes, grp, n.Name, fname, file); err != nil {
				return err
			}
		}
	case schema.Node_Which_interface:
		m, _ := n.Interface().Methods()
		methodResolve := func(id uint64, mname string, base string, name string) error {
			x := nodes[id]
			if x == nil {
				return fmt.Errorf("could not find type %#x for %s.%s", id, n, mname)
			}
			if x.ScopeId() != 0 {
				return nil
			}
			return resolveName(nodes, x, base, name, file)
		}
		for i := 0; i < m.Len(); i++ {
			mm := m.At(i)
			mname, _ := mm.Name()
			mann, _ := mm.Annotations()
			base := n.Name + "_" + parseAnnotations(mann).Rename(mname)
			if err := methodResolve(mm.ParamStructType(), mname, base, "Params"); err != nil {
				return err
			}
			if err := methodResolve(mm.ParamStructType(), mname, base, "Results"); err != nil {
				return err
			}
		}
	}
	return nil
}

type enumval struct {
	schema.Enumerant
	Name   string
	Val    int
	Tag    string
	parent *node
}

func makeEnumval(enum *node, i int, e schema.Enumerant) enumval {
	eann, _ := e.Annotations()
	ann := parseAnnotations(eann)
	name, _ := e.Name()
	name = ann.Rename(name)
	t := ann.Tag(name)
	return enumval{e, name, i, t, enum}
}

func (e *enumval) FullName() string {
	return e.parent.Name + "_" + e.Name
}

func (g *generator) defineEnum(n *node) error {
	es, _ := n.Enum().Enumerants()
	ev := make([]enumval, es.Len())
	for i := 0; i < es.Len(); i++ {
		e := es.At(i)
		ev[e.CodeOrder()] = makeEnumval(n, i, e)
	}
	nann, _ := n.Annotations()
	err := templates.ExecuteTemplate(&g.buf, "enum", enumParams{
		G:           g,
		Node:        n,
		Annotations: parseAnnotations(nann),
		EnumValues:  ev,
	})
	if err != nil {
		return fmt.Errorf("enum %s: %v", n, err)
	}
	return nil
}

func isValueOfType(v schema.Value, t schema.Type) bool {
	// Ensure that the value is for the given type.  The schema ensures the union ordinals match.
	return !v.IsValid() || int(v.Which()) == int(t.Which())
}

// Value formats a value from a schema (like a field default) as Go source.
func (g *generator) Value(n *node, t schema.Type, v schema.Value) (string, error) {
	if !isValueOfType(v, t) {
		return "", fmt.Errorf("value type is %v, but found %v value", t.Which(), v.Which())
	}

	switch t.Which() {
	case schema.Type_Which_void:
		return "struct{}{}", nil

	case schema.Type_Which_interface:
		// The only statically representable interface value is null.
		return g.imports.Capnp() + ".Client(nil)", nil

	case schema.Type_Which_bool:
		if v.Bool() {
			return "true", nil
		} else {
			return "false", nil
		}

	case schema.Type_Which_uint8, schema.Type_Which_uint16, schema.Type_Which_uint32, schema.Type_Which_uint64:
		return fmt.Sprintf("uint%d(%d)", intbits(t.Which()), uintValue(v)), nil

	case schema.Type_Which_int8, schema.Type_Which_int16, schema.Type_Which_int32, schema.Type_Which_int64:
		return fmt.Sprintf("int%d(%d)", intbits(t.Which()), intValue(v)), nil

	case schema.Type_Which_float32:
		return fmt.Sprintf("%s.Float32frombits(0x%x)", g.imports.Math(), math.Float32bits(v.Float32())), nil

	case schema.Type_Which_float64:
		return fmt.Sprintf("%s.Float64frombits(0x%x)", g.imports.Math(), math.Float64bits(v.Float64())), nil

	case schema.Type_Which_text:
		text, _ := v.Text()
		return strconv.Quote(text), nil

	case schema.Type_Which_data:
		buf := make([]byte, 0, 1024)
		buf = append(buf, "[]byte{"...)
		data, _ := v.Data()
		for i, b := range data {
			if i > 0 {
				buf = append(buf, ',', ' ')
			}
			buf = strconv.AppendUint(buf, uint64(b), 10)
		}
		buf = append(buf, '}')
		return string(buf), nil

	case schema.Type_Which_enum:
		en := g.nodes[t.Enum().TypeId()]
		if en == nil || !en.IsValid() || en.Which() != schema.Node_Which_enum {
			return "", errors.New("expected enum type")
		}
		enums, _ := en.Enum().Enumerants()
		val := int(v.Enum())
		if val >= enums.Len() {
			rn, err := g.RemoteName(en, n)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("%s(%d)", rn, val), nil
		}
		ev := makeEnumval(en, val, enums.At(val))
		name, err := g.appendRemoteScope(nil, en, n)
		if err != nil {
			return "", err
		}
		name = append(name, ev.FullName()...)
		return string(name), nil

	case schema.Type_Which_structGroup:
		data, _ := v.StructFieldPtr()
		var buf bytes.Buffer
		tn, err := g.nodes.mustFind(t.StructGroup().TypeId())
		if err != nil {
			return "", err
		}
		sd, err := g.data.copyData(data)
		if err != nil {
			return "", err
		}
		err = templates.ExecuteTemplate(&buf, "structValue", structValueTemplateParams{
			G:     g,
			Node:  n,
			Typ:   tn,
			Value: sd,
		})
		return buf.String(), err

	case schema.Type_Which_anyPointer:
		data, _ := v.AnyPointerPtr()
		var buf bytes.Buffer
		sd, err := g.data.copyData(data)
		if err != nil {
			return "", err
		}
		err = templates.ExecuteTemplate(&buf, "pointerValue", structValueTemplateParams{
			G:     g,
			Value: sd,
		})
		return buf.String(), err

	case schema.Type_Which_list:
		data, _ := v.ListPtr()
		var buf bytes.Buffer
		ftyp, err := g.fieldType(n, t, new(annotations))
		if err != nil {
			return "", err
		}
		sd, err := g.data.copyData(data)
		if err != nil {
			return "", err
		}
		err = templates.ExecuteTemplate(&buf, "listValue", listValueTemplateParams{
			G:     g,
			Typ:   ftyp,
			Value: sd,
		})
		return buf.String(), err
	default:
		return "", fmt.Errorf("unhandled value type %v", t.Which())
	}
}

func (g *generator) defineAnnotation(n *node) error {
	err := templates.ExecuteTemplate(&g.buf, "annotation", annotationParams{
		G:    g,
		Node: n,
	})
	if err != nil {
		return fmt.Errorf("annotation %s: %v", n, err)
	}
	return nil
}

func isGoConstType(t schema.Type) bool {
	w := t.Which()
	return w == schema.Type_Which_bool ||
		w == schema.Type_Which_int8 ||
		w == schema.Type_Which_uint8 ||
		w == schema.Type_Which_int16 ||
		w == schema.Type_Which_uint16 ||
		w == schema.Type_Which_int32 ||
		w == schema.Type_Which_uint32 ||
		w == schema.Type_Which_int64 ||
		w == schema.Type_Which_uint64 ||
		w == schema.Type_Which_text ||
		w == schema.Type_Which_enum
}

func (g *generator) defineConstNodes(nodes []*node) error {
	constNodes := make([]*node, 0, len(nodes))
	for _, n := range nodes {
		if n.Which() != schema.Node_Which_const {
			continue
		}
		t, _ := n.Const().Type()
		if isGoConstType(t) {
			constNodes = append(constNodes, n)
		}
	}
	nc := len(constNodes)
	for _, n := range nodes {
		if n.Which() != schema.Node_Which_const {
			continue
		}
		t, _ := n.Const().Type()
		if !isGoConstType(t) {
			constNodes = append(constNodes, n)
		}
	}
	if len(constNodes) == 0 {
		// short path
		return nil
	}
	err := templates.ExecuteTemplate(&g.buf, "constants", constantsParams{
		G:      g,
		Consts: constNodes[:nc],
		Vars:   constNodes[nc:],
	})
	if err != nil {
		return fmt.Errorf("file constants: %v", err)
	}
	return nil
}

func (g *generator) defineField(n *node, f field) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("field %s.%s: %v", n.shortDisplayName(), f.Name, err)
		}
	}()

	fann, _ := f.Annotations()
	ann := parseAnnotations(fann)
	t, _ := f.Slot().Type()
	def, _ := f.Slot().DefaultValue()
	if !isValueOfType(def, t) {
		return fmt.Errorf("default value type is %v, but found %v value", t.Which(), def.Which())
	}
	ftyp, err := g.fieldType(n, t, ann)
	if err != nil {
		return err
	}
	params := structFieldParams{
		G:           g,
		Node:        n,
		Field:       f,
		Annotations: ann,
		FieldType:   ftyp,
	}
	switch t.Which() {
	case schema.Type_Which_void:
		return templates.ExecuteTemplate(&g.buf, "structVoidField", params)
	case schema.Type_Which_bool:
		return templates.ExecuteTemplate(&g.buf, "structBoolField", structBoolFieldParams{
			structFieldParams: params,
			Default:           def.Bool(),
		})

	case schema.Type_Which_uint8, schema.Type_Which_uint16, schema.Type_Which_uint32, schema.Type_Which_uint64:
		return templates.ExecuteTemplate(&g.buf, "structUintField", structUintFieldParams{
			structFieldParams: params,
			Bits:              intbits(t.Which()),
			Default:           uintValue(def),
		})

	case schema.Type_Which_int8, schema.Type_Which_int16, schema.Type_Which_int32, schema.Type_Which_int64:
		return templates.ExecuteTemplate(&g.buf, "structIntField", structIntFieldParams{
			structUintFieldParams: structUintFieldParams{
				structFieldParams: params,
				Bits:              intbits(t.Which()),
				Default:           uint64(intFieldDefaultMask(def)),
			},
		})

	case schema.Type_Which_enum:
		ni, err := g.nodes.mustFind(t.Enum().TypeId())
		if err != nil {
			return err
		}
		rn, err := g.RemoteName(ni, n)
		if err != nil {
			return err
		}
		return templates.ExecuteTemplate(&g.buf, "structIntField", structIntFieldParams{
			structUintFieldParams: structUintFieldParams{
				structFieldParams: params,
				Bits:              16,
				Default:           uint64(def.Enum()),
			},
			EnumName: rn,
		})
	case schema.Type_Which_float32:
		return templates.ExecuteTemplate(&g.buf, "structFloatField", structUintFieldParams{
			structFieldParams: params,
			Bits:              32,
			Default:           uint64(math.Float32bits(def.Float32())),
		})

	case schema.Type_Which_float64:
		return templates.ExecuteTemplate(&g.buf, "structFloatField", structUintFieldParams{
			structFieldParams: params,
			Bits:              64,
			Default:           math.Float64bits(def.Float64()),
		})

	case schema.Type_Which_text:
		d, err := def.Text()
		if err != nil {
			return err
		}
		return templates.ExecuteTemplate(&g.buf, "structTextField", structTextFieldParams{
			structFieldParams: params,
			Default:           d,
		})

	case schema.Type_Which_data:
		d, err := def.Data()
		if err != nil {
			return err
		}
		return templates.ExecuteTemplate(&g.buf, "structDataField", structDataFieldParams{
			structFieldParams: params,
			Default:           d,
		})

	case schema.Type_Which_structGroup:
		var defref staticDataRef
		if sf, err := def.StructFieldPtr(); err != nil {
			return err
		} else if sf.IsValid() {
			defref, err = g.data.copyData(sf)
			if err != nil {
				return err
			}
		}
		tn, err := g.nodes.mustFind(t.StructGroup().TypeId())
		if err != nil {
			return err
		}
		return templates.ExecuteTemplate(&g.buf, "structStructField", structObjectFieldParams{
			structFieldParams: params,
			TypeNode:          tn,
			Default:           defref,
		})

	case schema.Type_Which_anyPointer:
		var defref staticDataRef
		if p, err := def.AnyPointerPtr(); err != nil {
			return err
		} else if p.IsValid() {
			defref, err = g.data.copyData(p)
			if err != nil {
				return err
			}
		}
		return templates.ExecuteTemplate(&g.buf, "structPointerField", structObjectFieldParams{
			structFieldParams: params,
			Default:           defref,
		})

	case schema.Type_Which_list:
		var defref staticDataRef
		if l, err := def.ListPtr(); err != nil {
			return err
		} else if l.IsValid() {
			defref, err = g.data.copyData(l)
			if err != nil {
				return err
			}
		}
		return templates.ExecuteTemplate(&g.buf, "structListField", structObjectFieldParams{
			structFieldParams: params,
			Default:           defref,
		})

	case schema.Type_Which_interface:
		return templates.ExecuteTemplate(&g.buf, "structInterfaceField", params)
	default:
		return fmt.Errorf("defining unhandled field type %v", t.Which())
	}
}

func (g *generator) fieldType(n *node, t schema.Type, ann *annotations) (string, error) {
	switch t.Which() {
	case schema.Type_Which_void:
		return "", nil
	case schema.Type_Which_bool:
		return "bool", nil
	case schema.Type_Which_int8:
		return "int8", nil
	case schema.Type_Which_int16:
		return "int16", nil
	case schema.Type_Which_int32:
		return "int32", nil
	case schema.Type_Which_int64:
		return "int64", nil
	case schema.Type_Which_uint8:
		return "uint8", nil
	case schema.Type_Which_uint16:
		return "uint16", nil
	case schema.Type_Which_uint32:
		return "uint32", nil
	case schema.Type_Which_uint64:
		return "uint64", nil
	case schema.Type_Which_float32:
		return "float32", nil
	case schema.Type_Which_float64:
		return "float64", nil
	case schema.Type_Which_text:
		return "string", nil
	case schema.Type_Which_data:
		return "[]byte", nil
	case schema.Type_Which_enum:
		ni, err := g.nodes.mustFind(t.Enum().TypeId())
		if err != nil {
			return "", err
		}
		return g.RemoteName(ni, n)
	case schema.Type_Which_structGroup:
		ni, err := g.nodes.mustFind(t.StructGroup().TypeId())
		if err != nil {
			return "", err
		}
		return g.RemoteName(ni, n)
	case schema.Type_Which_interface:
		ni, err := g.nodes.mustFind(t.Interface().TypeId())
		if err != nil {
			return "", err
		}
		return g.RemoteName(ni, n)
	case schema.Type_Which_anyPointer:
		return g.imports.Capnp() + ".Pointer", nil
	case schema.Type_Which_list:
		switch lt, _ := t.List().ElementType(); lt.Which() {
		case schema.Type_Which_void:
			return g.imports.Capnp() + ".VoidList", nil
		case schema.Type_Which_bool:
			return g.imports.Capnp() + ".BitList", nil
		case schema.Type_Which_int8:
			return g.imports.Capnp() + ".Int8List", nil
		case schema.Type_Which_uint8:
			return g.imports.Capnp() + ".UInt8List", nil
		case schema.Type_Which_int16:
			return g.imports.Capnp() + ".Int16List", nil
		case schema.Type_Which_uint16:
			return g.imports.Capnp() + ".UInt16List", nil
		case schema.Type_Which_int32:
			return g.imports.Capnp() + ".Int32List", nil
		case schema.Type_Which_uint32:
			return g.imports.Capnp() + ".UInt32List", nil
		case schema.Type_Which_int64:
			return g.imports.Capnp() + ".Int64List", nil
		case schema.Type_Which_uint64:
			return g.imports.Capnp() + ".UInt64List", nil
		case schema.Type_Which_float32:
			return g.imports.Capnp() + ".Float32List", nil
		case schema.Type_Which_float64:
			return g.imports.Capnp() + ".Float64List", nil
		case schema.Type_Which_text:
			return g.imports.Capnp() + ".TextList", nil
		case schema.Type_Which_data:
			return g.imports.Capnp() + ".DataList", nil
		case schema.Type_Which_enum:
			ni, err := g.nodes.mustFind(lt.Enum().TypeId())
			if err != nil {
				return "", err
			}
			rn, err := g.RemoteName(ni, n)
			if err != nil {
				return "", err
			}
			return rn + "_List", nil
		case schema.Type_Which_structGroup:
			ni, err := g.nodes.mustFind(lt.StructGroup().TypeId())
			if err != nil {
				return "", err
			}
			rn, err := g.RemoteName(ni, n)
			if err != nil {
				return "", err
			}
			return rn + "_List", nil
		case schema.Type_Which_anyPointer, schema.Type_Which_list, schema.Type_Which_interface:
			return g.imports.Capnp() + ".PointerList", nil
		}
	}
	return "", fmt.Errorf("unhandled field type %v", t.Which())
}

// intFieldDefaultMask returns the XOR mask used when getting or setting
// signed integer struct fields.
func intFieldDefaultMask(v schema.Value) uint64 {
	mask := uint64(1)<<intbits(schema.Type_Which(v.Which())) - 1
	return uint64(intValue(v)) & mask
}

// intValue returns the signed integer value of a schema value or zero
// if the value is invalid (the null pointer). Panics if the value is
// not a signed integer.
func intValue(v schema.Value) int64 {
	if !v.IsValid() {
		return 0
	}
	switch v.Which() {
	case schema.Value_Which_int8:
		return int64(v.Int8())
	case schema.Value_Which_int16:
		return int64(v.Int16())
	case schema.Value_Which_int32:
		return int64(v.Int32())
	case schema.Value_Which_int64:
		return v.Int64()
	}
	panic("unreachable")
}

// uintValue returns the unsigned integer value of a schema value or
// zero if the value is invalid (the null pointer). Panics if the value
// is not an unsigned integer.
func uintValue(v schema.Value) uint64 {
	if !v.IsValid() {
		return 0
	}
	switch v.Which() {
	case schema.Value_Which_uint8:
		return uint64(v.Uint8())
	case schema.Value_Which_uint16:
		return uint64(v.Uint16())
	case schema.Value_Which_uint32:
		return uint64(v.Uint32())
	case schema.Value_Which_uint64:
		return v.Uint64()
	}
	panic("unreachable")
}

// intbits returns the number of bits that an integer type requires.
func intbits(t schema.Type_Which) uint {
	switch t {
	case schema.Type_Which_uint8, schema.Type_Which_int8:
		return 8
	case schema.Type_Which_uint16, schema.Type_Which_int16:
		return 16
	case schema.Type_Which_uint32, schema.Type_Which_int32:
		return 32
	case schema.Type_Which_uint64, schema.Type_Which_int64:
		return 64
	default:
		panic("unreachable")
	}
}

func (n *node) codeOrderFields() []field {
	fields, _ := n.StructGroup().Fields()
	numFields := fields.Len()
	mbrs := make([]field, numFields)
	for i := 0; i < numFields; i++ {
		f := fields.At(i)
		fann, _ := f.Annotations()
		fname, _ := f.Name()
		fname = parseAnnotations(fann).Rename(fname)
		mbrs[f.CodeOrder()] = field{Field: f, Name: fname}
	}
	return mbrs
}

func (g *generator) defineStructTypes(n, baseNode *node) error {
	nann, _ := n.Annotations()
	ann := parseAnnotations(nann)
	err := templates.ExecuteTemplate(&g.buf, "structTypes", structTypesParams{
		G:           g,
		Node:        n,
		Annotations: ann,
		BaseNode:    baseNode,
	})
	if err != nil {
		dn, _ := n.DisplayName()
		return fmt.Errorf("struct type for %s: %v", dn, err)
	}

	for _, f := range n.codeOrderFields() {
		if f.Which() == schema.Field_Which_group {
			grp, err := g.nodes.mustFind(f.Group().TypeId())
			if err != nil {
				return err
			}
			if err := g.defineStructTypes(grp, baseNode); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *generator) defineStructEnums(n *node) error {
	fields := n.codeOrderFields()
	members := make([]field, 0, len(fields))
	es := make(enumString, 0, len(fields))
	for _, f := range fields {
		if f.DiscriminantValue() != schema.Field_noDiscriminant {
			members = append(members, f)
			es = append(es, f.Name)
		}
	}
	if n.StructGroup().DiscriminantCount() > 0 {
		err := templates.ExecuteTemplate(&g.buf, "structEnums", structEnumsParams{
			G:          g,
			Node:       n,
			Fields:     members,
			EnumString: es,
		})
		if err != nil {
			return fmt.Errorf("struct enums for %s: %v", n, err)
		}
	}
	for _, f := range fields {
		if f.Which() == schema.Field_Which_group {
			grp, err := g.nodes.mustFind(f.Group().TypeId())
			if err != nil {
				return err
			}
			if err := g.defineStructEnums(grp); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *generator) defineStructFuncs(n *node) error {
	err := templates.ExecuteTemplate(&g.buf, "structFuncs", structFuncsParams{
		G:    g,
		Node: n,
	})
	if err != nil {
		return fmt.Errorf("struct funcs for %s: %v", n, err)
	}

	for _, f := range n.codeOrderFields() {
		switch f.Which() {
		case schema.Field_Which_slot:
			return g.defineField(n, f)
		case schema.Field_Which_group:
			grp, err := g.nodes.mustFind(f.Group().TypeId())
			if err != nil {
				return err
			}
			err = templates.ExecuteTemplate(&g.buf, "structGroup", structGroupParams{
				G:     g,
				Node:  n,
				Group: grp,
				Field: f,
			})
			if err != nil {
				return fmt.Errorf("struct group for %s: %v", grp, err)
			}
			if err := g.defineStructFuncs(grp); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *generator) ObjectSize(n *node) (string, error) {
	if n.Which() != schema.Node_Which_structGroup {
		return "", fmt.Errorf("object size called for %v node", n.Which())
	}
	return fmt.Sprintf("%s.ObjectSize{DataSize: %d, PointerCount: %d}", g.imports.Capnp(), int(n.StructGroup().DataWordCount())*8, n.StructGroup().PointerCount()), nil
}

func (g *generator) defineNewStructFunc(n *node) error {
	err := templates.ExecuteTemplate(&g.buf, "newStructFunc", newStructParams{
		G:    g,
		Node: n,
	})
	if err != nil {
		return fmt.Errorf("new struct function for %s: %v", n, err)
	}
	return nil
}

func (g *generator) defineStructList(n *node) error {
	err := templates.ExecuteTemplate(&g.buf, "structList", structListParams{
		G:    g,
		Node: n,
	})
	if err != nil {
		return fmt.Errorf("new struct function for %s: %v", n, err)
	}
	return nil
}

func (g *generator) defineStructPromise(n *node) error {
	err := templates.ExecuteTemplate(&g.buf, "promise", promiseTemplateParams{
		G:      g,
		Node:   n,
		Fields: n.codeOrderFields(),
	})
	if err != nil {
		return fmt.Errorf("promise for struct %s: %v", n, err)
	}

	for _, f := range n.codeOrderFields() {
		switch f.Which() {
		case schema.Field_Which_slot:
			t, _ := f.Slot().Type()
			if tw := t.Which(); tw != schema.Type_Which_structGroup && tw != schema.Type_Which_interface && tw != schema.Type_Which_anyPointer {
				continue
			}
			if err := g.definePromiseField(n, f); err != nil {
				return fmt.Errorf("promise field %s.%s: %v", n.shortDisplayName(), f.Name, err)
			}
		case schema.Field_Which_group:
			grp, err := g.nodes.mustFind(f.Group().TypeId())
			if err != nil {
				return fmt.Errorf("promise group %s.%s: %v", n.shortDisplayName(), f.Name, err)
			}
			err = templates.ExecuteTemplate(&g.buf, "promiseGroup", promiseGroupTemplateParams{
				G:     g,
				Node:  n,
				Field: f,
				Group: grp,
			})
			if err != nil {
				return fmt.Errorf("promise for group %s: %v", grp, err)
			}
			if err := g.defineStructPromise(grp); err != nil {
				return err
			}
		}
	}
	return nil
}

func (g *generator) definePromiseField(n *node, f field) error {
	slot := f.Slot()
	switch t, _ := slot.Type(); t.Which() {
	case schema.Type_Which_structGroup:
		ni, err := g.nodes.mustFind(t.StructGroup().TypeId())
		if err != nil {
			return err
		}
		params := promiseFieldStructTemplateParams{
			G:      g,
			Node:   n,
			Field:  f,
			Struct: ni,
		}
		if def, _ := slot.DefaultValue(); def.IsValid() && def.Which() == schema.Value_Which_structField {
			if sf, _ := def.StructFieldPtr(); sf.IsValid() {
				params.Default, err = g.data.copyData(sf)
				if err != nil {
					return err
				}
			}
		}
		return templates.ExecuteTemplate(&g.buf, "promiseFieldStruct", params)
	case schema.Type_Which_anyPointer:
		return templates.ExecuteTemplate(&g.buf, "promiseFieldAnyPointer", promiseFieldAnyPointerTemplateParams{
			G:     g,
			Node:  n,
			Field: f,
		})
	case schema.Type_Which_interface:
		ni, err := g.nodes.mustFind(t.Interface().TypeId())
		if err != nil {
			return err
		}
		return templates.ExecuteTemplate(&g.buf, "promiseFieldInterface", promiseFieldInterfaceTemplateParams{
			G:         g,
			Node:      n,
			Field:     f,
			Interface: ni,
		})
	default:
		panic("unreachable")
	}
}

type interfaceMethod struct {
	schema.Method
	Interface    *node
	ID           int
	Name         string
	OriginalName string
	Params       *node
	Results      *node
}

func methodSet(methods []interfaceMethod, n *node, nodes nodeMap) ([]interfaceMethod, error) {
	ms, _ := n.Interface().Methods()
	for i := 0; i < ms.Len(); i++ {
		m := ms.At(i)
		mname, _ := m.Name()
		mann, _ := m.Annotations()
		pn, err := nodes.mustFind(m.ParamStructType())
		if err != nil {
			return methods, fmt.Errorf("could not find param type for %s.%s", n.shortDisplayName(), mname)
		}
		rn, err := nodes.mustFind(m.ResultStructType())
		if err != nil {
			return methods, fmt.Errorf("could not find result type for %s.%s", n.shortDisplayName(), mname)
		}
		methods = append(methods, interfaceMethod{
			Method:       m,
			Interface:    n,
			ID:           i,
			OriginalName: mname,
			Name:         parseAnnotations(mann).Rename(mname),
			Params:       pn,
			Results:      rn,
		})
	}
	// TODO(light): sort added methods by code order

	supers, _ := n.Interface().Superclasses()
	for i := 0; i < supers.Len(); i++ {
		s := supers.At(i)
		sn, err := nodes.mustFind(s.Id())
		if err != nil {
			return methods, fmt.Errorf("could not find superclass %#x of %s", s.Id(), n)
		}
		methods, err = methodSet(methods, sn, nodes)
		if err != nil {
			return methods, err
		}
	}
	return methods, nil
}

func (g *generator) defineInterface(n *node) error {
	m, err := methodSet(nil, n, g.nodes)
	if err != nil {
		return fmt.Errorf("building method set of interface %s: %v", n, err)
	}
	nann, _ := n.Annotations()
	err = templates.ExecuteTemplate(&g.buf, "interfaceClient", interfaceClientTemplateParams{
		G:           g,
		Node:        n,
		Annotations: parseAnnotations(nann),
		Methods:     m,
	})
	if err != nil {
		return fmt.Errorf("interface client %s: %v", n, err)
	}
	err = templates.ExecuteTemplate(&g.buf, "interfaceServer", interfaceServerTemplateParams{
		G:           g,
		Node:        n,
		Annotations: parseAnnotations(nann),
		Methods:     m,
	})
	if err != nil {
		return fmt.Errorf("interface server %s: %v", n, err)
	}
	return nil
}

type enumString []string

func (es enumString) ValueString() string {
	return strings.Join([]string(es), "")
}

func (es enumString) SliceFor(i int) string {
	n := 0
	for _, v := range es[:i] {
		n += len(v)
	}
	return fmt.Sprintf("[%d:%d]", n, n+len(es[i]))
}

type errCollector struct {
	err error
}

func (ec *errCollector) collect(e error) {
	if ec.err == nil {
		ec.err = e
	}
}

func generateFile(reqf schema.CodeGeneratorRequest_RequestedFile, nodes nodeMap) error {
	id := reqf.Id()
	fname, _ := reqf.Filename()
	g := newGenerator(id, nodes)
	f := nodes[id]
	if f == nil {
		return fmt.Errorf("no node in schema matches %#x", id)
	}
	if f.pkg == "" {
		return fmt.Errorf("missing package annotation for %s", fname)
	}

	var ec errCollector
	for _, n := range f.nodes {
		if n.Which() == schema.Node_Which_annotation {
			ec.collect(g.defineAnnotation(n))
		}
	}
	ec.collect(g.defineConstNodes(f.nodes))
	for _, n := range f.nodes {
		switch n.Which() {
		case schema.Node_Which_enum:
			ec.collect(g.defineEnum(n))
		case schema.Node_Which_structGroup:
			if !n.StructGroup().IsGroup() {
				ec.collect(g.defineStructTypes(n, n))
				ec.collect(g.defineStructEnums(n))
				ec.collect(g.defineNewStructFunc(n))
				ec.collect(g.defineStructFuncs(n))
				ec.collect(g.defineStructList(n))
				if *genPromises {
					ec.collect(g.defineStructPromise(n))
				}
			}
		case schema.Node_Which_interface:
			ec.collect(g.defineInterface(n))
		}
		if ec.err != nil {
			break
		}
	}
	if ec.err != nil {
		return ec.err
	}

	if dirPath, _ := filepath.Split(fname); dirPath != "" {
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	unformatted := g.generate(f.pkg)
	formatted, fmtErr := format.Source(unformatted)
	if fmtErr != nil {
		formatted = unformatted
	}

	file, err := os.Create(fname + ".go")
	if err != nil {
		return err
	}
	_, werr := file.Write(formatted)
	cerr := file.Close()
	if fmtErr != nil {
		return fmtErr
	}
	if werr != nil {
		return err
	}
	if cerr != nil {
		return err
	}
	return nil
}

func main() {
	flag.Parse()

	msg, err := capnp.NewDecoder(os.Stdin).Decode()
	if err != nil {
		fmt.Fprintln(os.Stderr, "capnpc-go: reading input:", err)
		os.Exit(1)
	}
	req, err := schema.ReadRootCodeGeneratorRequest(msg)
	if err != nil {
		fmt.Fprintln(os.Stderr, "capnpc-go: reading input:", err)
		os.Exit(1)
	}
	nodes, err := buildNodeMap(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "capnpc-go:", err)
		os.Exit(1)
	}
	success := true
	reqFiles, _ := req.RequestedFiles()
	for i := 0; i < reqFiles.Len(); i++ {
		reqf := reqFiles.At(i)
		err := generateFile(reqf, nodes)
		if err != nil {
			fname, _ := reqf.Filename()
			fmt.Fprintf(os.Stderr, "capnpc-go: generating %s: %v\n", fname, err)
			success = false
		}
	}
	if !success {
		os.Exit(1)
	}
}
