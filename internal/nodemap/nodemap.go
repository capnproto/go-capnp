// Package nodemap provides a schema registry index type.
package nodemap

import (
	"sync"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/internal/schema"
	"capnproto.org/go/capnp/v3/schemas"
)

// Map is a lazy index of a registry.
// The zero value is an index of the default registry.
type Map struct {
	reg   *schemas.Registry
	index *nodeIndex
}

// defaultIndex is shared by all zero-value Maps.  Generated schemas are
// immutable after registration, so nodes are safe to share once an entire
// decoded request has been published under the lock.
var defaultIndex = new(nodeIndex)

type nodeIndex struct {
	mu    sync.RWMutex
	nodes map[uint64]schema.Node
}

func (m *Map) registry() *schemas.Registry {
	if m.reg != nil {
		return m.reg
	}
	return schemas.DefaultRegistry
}

// UseRegistry assigns reg to m and initializes an index for it.  The default
// registry uses the process-wide shared index; custom registries remain local
// to the Map.
func (m *Map) UseRegistry(reg *schemas.Registry) {
	m.reg = reg
	m.index = new(nodeIndex)
}

// Find returns the node for the given ID.
func (m *Map) Find(id uint64) (schema.Node, error) {
	reg := m.registry()
	index := m.index
	if reg == schemas.DefaultRegistry {
		index = defaultIndex
	}
	return index.find(reg, id)
}

func (index *nodeIndex) find(reg *schemas.Registry, id uint64) (schema.Node, error) {
	index.mu.RLock()
	n := index.nodes[id]
	index.mu.RUnlock()
	if n.IsValid() {
		return n, nil
	}
	data, err := reg.Find(id)
	if err != nil {
		return schema.Node{}, err
	}
	msg, err := capnp.Unmarshal(data)
	if err != nil {
		return schema.Node{}, err
	}
	req, err := schema.ReadRootCodeGeneratorRequest(msg)
	if err != nil {
		return schema.Node{}, err
	}
	nodes, err := req.Nodes()
	if err != nil {
		return schema.Node{}, err
	}
	decoded := make(map[uint64]schema.Node, nodes.Len())
	for i := 0; i < nodes.Len(); i++ {
		n := nodes.At(i)
		decoded[n.Id()] = n
	}
	// Cached nodes retain msg for the lifetime of the index.  Schema data is
	// trusted, immutable metadata, so a finite cumulative traversal budget
	// would eventually make an otherwise valid cache entry unreadable.
	msg.ResetReadLimit(^uint64(0))

	index.mu.Lock()
	if index.nodes == nil {
		index.nodes = make(map[uint64]schema.Node)
	}
	for nodeID, n := range decoded {
		if !index.nodes[nodeID].IsValid() {
			index.nodes[nodeID] = n
		}
	}
	n = index.nodes[id]
	index.mu.Unlock()
	return n, nil
}
