// Package text supports marshalling Cap'n Proto as text based on a schema.
package text

import (
	"bytes"
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

func (m *marshaller) marshalStruct(typeID uint64, s capnp.Struct) error {
	return nil
}
