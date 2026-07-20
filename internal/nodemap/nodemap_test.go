package nodemap

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/internal/schema"
	"capnproto.org/go/capnp/v3/schemas"
)

func makeSchema(t *testing.T, nodes map[uint64]string) []byte {
	t.Helper()
	msg, seg := capnp.NewSingleSegmentMessage(nil)
	req, err := schema.NewRootCodeGeneratorRequest(seg)
	if err != nil {
		t.Fatal(err)
	}
	list, err := req.NewNodes(int32(len(nodes)))
	if err != nil {
		t.Fatal(err)
	}
	i := 0
	for id, name := range nodes {
		n := list.At(i)
		n.SetId(id)
		if err := n.SetDisplayName(name); err != nil {
			t.Fatal(err)
		}
		i++
	}
	data, err := msg.Marshal()
	if err != nil {
		t.Fatal(err)
	}
	return data
}

func registerSchema(t *testing.T, reg *schemas.Registry, nodes map[uint64]string) {
	t.Helper()
	ids := make([]uint64, 0, len(nodes))
	for id := range nodes {
		ids = append(ids, id)
	}
	if err := reg.Register(&schemas.Schema{Bytes: makeSchema(t, nodes), Nodes: ids}); err != nil {
		t.Fatal(err)
	}
}

func nodeName(t *testing.T, n schema.Node) string {
	t.Helper()
	name, err := n.DisplayName()
	if err != nil {
		t.Fatal(err)
	}
	return name
}

func useFreshDefaultRegistry(t *testing.T) {
	t.Helper()
	oldRegistry := schemas.DefaultRegistry
	schemas.DefaultRegistry = new(schemas.Registry)
	t.Cleanup(func() {
		schemas.DefaultRegistry = oldRegistry
	})
}

func TestDefaultIndexFindsLaterRegistration(t *testing.T) {
	useFreshDefaultRegistry(t)
	const firstID = 0x9c5c2bf81252e101
	const laterID = 0x9c5c2bf81252e102
	const siblingID = 0x9c5c2bf81252e103
	registerSchema(t, schemas.DefaultRegistry, map[uint64]string{
		firstID:   "first",
		siblingID: "sibling",
	})

	var first Map
	firstNode, err := first.Find(firstID)
	if err != nil {
		t.Fatal(err)
	}
	var sibling Map
	siblingNode, err := sibling.Find(siblingID)
	if err != nil {
		t.Fatal(err)
	}
	if firstNode.Message() != siblingNode.Message() {
		t.Fatal("default Maps did not share the decoded schema request")
	}

	// Populating the shared index must not snapshot the registry.  A schema
	// registered before its first lookup is still found on a later miss.
	registerSchema(t, schemas.DefaultRegistry, map[uint64]string{laterID: "later"})
	var later Map
	n, err := later.Find(laterID)
	if err != nil {
		t.Fatal(err)
	}
	if got := nodeName(t, n); got != "later" {
		t.Fatalf("display name = %q; want later", got)
	}
}

func TestDefaultIndexTracksRegistryReplacement(t *testing.T) {
	const id = 0x9c5c2bf81252e180
	first := new(schemas.Registry)
	second := new(schemas.Registry)
	registerSchema(t, first, map[uint64]string{id: "first registry"})
	registerSchema(t, second, map[uint64]string{id: "second registry"})

	oldRegistry := schemas.DefaultRegistry
	t.Cleanup(func() { schemas.DefaultRegistry = oldRegistry })
	schemas.DefaultRegistry = first

	var firstMap Map
	n, err := firstMap.Find(id)
	if err != nil {
		t.Fatal(err)
	}
	if got := nodeName(t, n); got != "first registry" {
		t.Fatalf("first registry returned %q", got)
	}

	schemas.DefaultRegistry = second
	var secondMap Map
	n, err = secondMap.Find(id)
	if err != nil {
		t.Fatal(err)
	}
	if got := nodeName(t, n); got != "second registry" {
		t.Fatalf("replacement registry returned %q; want second registry", got)
	}
}

func TestCustomRegistriesAreIsolated(t *testing.T) {
	const id = 0x9c5c2bf81252e201
	reg1, reg2 := new(schemas.Registry), new(schemas.Registry)
	registerSchema(t, reg1, map[uint64]string{id: "registry one"})
	registerSchema(t, reg2, map[uint64]string{id: "registry two"})

	var m1, m2 Map
	m1.UseRegistry(reg1)
	m2.UseRegistry(reg2)
	n1, err := m1.Find(id)
	if err != nil {
		t.Fatal(err)
	}
	n2, err := m2.Find(id)
	if err != nil {
		t.Fatal(err)
	}
	if got := nodeName(t, n1); got != "registry one" {
		t.Fatalf("first registry returned %q", got)
	}
	if got := nodeName(t, n2); got != "registry two" {
		t.Fatalf("second registry returned %q", got)
	}
}

func TestConcurrentDefaultFirstUse(t *testing.T) {
	useFreshDefaultRegistry(t)
	const baseID = 0x9c5c2bf81252e300
	nodes := make(map[uint64]string)
	for i := uint64(0); i < 4; i++ {
		nodes[baseID+i] = fmt.Sprintf("node %d", i)
	}
	registerSchema(t, schemas.DefaultRegistry, nodes)

	const goroutines = 32
	start := make(chan struct{})
	errs := make(chan error, goroutines)
	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			var m Map
			id := baseID + uint64(i%len(nodes))
			n, err := m.Find(id)
			if err == nil && !n.IsValid() {
				err = fmt.Errorf("node %#x is invalid", id)
			}
			errs <- err
		}(i)
	}
	close(start)
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Error(err)
		}
	}
}

func TestCachedNodeDoesNotExhaustTraversalLimit(t *testing.T) {
	useFreshDefaultRegistry(t)
	const id = 0x9c5c2bf81252e400
	name := strings.Repeat("x", 1<<20)
	registerSchema(t, schemas.DefaultRegistry, map[uint64]string{id: name})

	var m Map
	n, err := m.Find(id)
	if err != nil {
		t.Fatal(err)
	}
	// The default message traversal limit is 64 MiB.  Reading this field 80
	// times would exhaust a cached message unless publication makes its
	// trusted schema metadata effectively non-exhausting.
	for i := 0; i < 80; i++ {
		got, err := n.DisplayName()
		if err != nil {
			t.Fatalf("DisplayName iteration %d: %v", i, err)
		}
		if len(got) != len(name) {
			t.Fatalf("DisplayName iteration %d length = %d; want %d", i, len(got), len(name))
		}
	}
}
