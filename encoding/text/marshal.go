// Package text supports marshalling Cap'n Proto as text based on a schema.
package text

import (
	"bytes"
	"fmt"
	"strconv"

	"zombiezen.com/go/capnproto2"
	"zombiezen.com/go/capnproto2/internal/schema"
)

// AppendStruct appends the text representation of s to dst.
func AppendStruct(dst []byte, s capnp.Struct, typeID uint64, schemaData []byte) ([]byte, error) {
	msg, err := capnp.Unmarshal(schemaData)
	if err != nil {
		return nil, err
	}
	req, err := schema.ReadRootCodeGeneratorRequest(msg)
	if err != nil {
		return nil, err
	}
	nodes, err := req.Nodes()
	if err != nil {
		return nil, err
	}
	nodeMap := make(map[uint64]schema.Node, nodes.Len())
	for i := 0; i < nodes.Len(); i++ {
		n := nodes.At(i)
		nodeMap[n.Id()] = n
	}
	m := &marshaller{
		buf:   bytes.NewBuffer(dst),
		nodes: nodeMap,
	}
	err = m.marshalStruct(typeID, s)
	return m.buf.Bytes(), err
}

type marshaller struct {
	buf   *bytes.Buffer
	tmp   []byte
	nodes map[uint64]schema.Node
}

func (m *marshaller) marshalBool(v bool) {
	if v {
		m.buf.WriteString("true")
	} else {
		m.buf.WriteString("false")
	}
}

func (m *marshaller) marshalInt(i int64) {
	m.tmp = strconv.AppendInt(m.tmp[:0], i, 10)
	m.buf.Write(m.tmp)
}

func (m *marshaller) marshalUint(i uint64) {
	m.tmp = strconv.AppendUint(m.tmp[:0], i, 10)
	m.buf.Write(m.tmp)
}

func (m *marshaller) marshalText(t []byte) {
	m.buf.WriteByte('"')
	last := 0
	for i, b := range t {
		if !needsEscape(b) {
			continue
		}
		m.buf.Write(t[last:i])
		switch b {
		case '\a':
			m.buf.WriteString("\\a")
		case '\b':
			m.buf.WriteString("\\b")
		case '\f':
			m.buf.WriteString("\\f")
		case '\n':
			m.buf.WriteString("\\n")
		case '\r':
			m.buf.WriteString("\\r")
		case '\t':
			m.buf.WriteString("\\t")
		case '\v':
			m.buf.WriteString("\\v")
		case '\'':
			m.buf.WriteString("\\'")
		case '"':
			m.buf.WriteString("\\\"")
		case '\\':
			m.buf.WriteString("\\\\")
		default:
			m.buf.WriteString("\\x")
			m.buf.WriteByte(hexDigit(b / 16))
			m.buf.WriteByte(hexDigit(b % 16))
		}
		last = i + 1
	}
	m.buf.Write(t[last:])
	m.buf.WriteByte('"')
}

func needsEscape(b byte) bool {
	return b < 0x20 || b >= 0x7f
}

func hexDigit(b byte) byte {
	const digits = "0123456789abcdef"
	return digits[b]
}

func (m *marshaller) marshalStruct(typeID uint64, s capnp.Struct) error {
	n := m.nodes[typeID]
	if !n.IsValid() || n.Which() != schema.Node_Which_structNode {
		return fmt.Errorf("cannot find struct type %#x", typeID)
	}
	var discriminant uint16
	if n.StructNode().DiscriminantCount() > 0 {
		discriminant = s.Uint16(capnp.DataOffset(n.StructNode().DiscriminantOffset() * 2))
	}
	m.buf.WriteByte('(')
	fields := codeOrderFields(n.StructNode())
	first := true
	for _, f := range fields {
		if !(f.Which() == schema.Field_Which_slot || f.Which() == schema.Field_Which_group) {
			continue
		}
		if dv := f.DiscriminantValue(); !(dv == schema.Field_noDiscriminant || dv == discriminant) {
			continue
		}
		if !first {
			m.buf.WriteString(", ")
		}
		first = false
		name, err := f.NameBytes()
		if err != nil {
			return err
		}
		m.buf.Write(name)
		m.buf.WriteString(" = ")
		switch f.Which() {
		case schema.Field_Which_slot:
			if err := m.marshalFieldValue(s, f); err != nil {
				return err
			}
		case schema.Field_Which_group:
			if err := m.marshalStruct(f.Group().TypeId(), s); err != nil {
				return err
			}
		}
	}
	m.buf.WriteByte(')')
	return nil
}

func (m *marshaller) marshalFieldValue(s capnp.Struct, f schema.Field) error {
	typ, err := f.Slot().Type()
	if err != nil {
		return err
	}
	dv, err := f.Slot().DefaultValue()
	if err != nil {
		return err
	}
	if dv.IsValid() && int(typ.Which()) != int(dv.Which()) {
		name, _ := f.Name()
		return fmt.Errorf("marshal field %s: default value is a %v, want %v", name, dv.Which(), typ.Which())
	}
	switch typ.Which() {
	case schema.Type_Which_void:
		m.buf.WriteString("void")
	case schema.Type_Which_bool:
		v := s.Bit(capnp.BitOffset(f.Slot().Offset()))
		d := dv.Bool()
		m.marshalBool(!d && v || d && !v)
	case schema.Type_Which_int8:
		v := s.Uint8(capnp.DataOffset(f.Slot().Offset()))
		d := uint8(dv.Int8())
		m.marshalInt(int64(int8(v ^ d)))
	case schema.Type_Which_int16:
		v := s.Uint16(capnp.DataOffset(f.Slot().Offset() * 2))
		d := uint16(dv.Int16())
		m.marshalInt(int64(int16(v ^ d)))
	case schema.Type_Which_int32:
		v := s.Uint32(capnp.DataOffset(f.Slot().Offset() * 4))
		d := uint32(dv.Int32())
		m.marshalInt(int64(int32(v ^ d)))
	case schema.Type_Which_int64:
		v := s.Uint64(capnp.DataOffset(f.Slot().Offset() * 8))
		d := uint64(dv.Int64())
		m.marshalInt(int64(v ^ d))
	case schema.Type_Which_uint8:
		v := s.Uint8(capnp.DataOffset(f.Slot().Offset()))
		d := dv.Uint8()
		m.marshalUint(uint64(v ^ d))
	case schema.Type_Which_uint16:
		v := s.Uint16(capnp.DataOffset(f.Slot().Offset() * 2))
		d := dv.Uint16()
		m.marshalUint(uint64(v ^ d))
	case schema.Type_Which_uint32:
		v := s.Uint32(capnp.DataOffset(f.Slot().Offset() * 4))
		d := dv.Uint32()
		m.marshalUint(uint64(v ^ d))
	case schema.Type_Which_uint64:
		v := s.Uint64(capnp.DataOffset(f.Slot().Offset() * 8))
		d := dv.Uint64()
		m.marshalUint(v ^ d)
	case schema.Type_Which_structType:
		p, err := s.Ptr(uint16(f.Slot().Offset()))
		if err != nil {
			return err
		}
		if !p.IsValid() {
			p, _ = dv.StructValuePtr()
		}
		return m.marshalStruct(typ.StructType().TypeId(), p.Struct())
	case schema.Type_Which_data:
		p, err := s.Ptr(uint16(f.Slot().Offset()))
		if err != nil {
			return err
		}
		if !p.IsValid() {
			b, _ := dv.Data()
			m.marshalText(b)
			return nil
		}
		m.marshalText(p.Data())
	case schema.Type_Which_text:
		p, err := s.Ptr(uint16(f.Slot().Offset()))
		if err != nil {
			return err
		}
		if !p.IsValid() {
			b, _ := dv.TextBytes()
			m.marshalText(b)
			return nil
		}
		b := p.Data()
		if len(b) > 0 {
			// Trim NUL byte
			b = b[:len(b)-1]
		}
		m.marshalText(b)
	}
	return nil
}

func codeOrderFields(s schema.Node_structNode) []schema.Field {
	list, _ := s.Fields()
	n := list.Len()
	fields := make([]schema.Field, n)
	for i := 0; i < n; i++ {
		f := list.At(i)
		fields[f.CodeOrder()] = f
	}
	return fields
}
