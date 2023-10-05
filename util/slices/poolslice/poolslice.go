// Package poolslice supports allocating slots out of a slice.
package poolslice

// A PoolSlice wraps a slice, adding support for tracking free vs. in-use
// slots and inserting new elements into free slots.
type PoolSlice[T any] struct {
	Items  []T            // The slice
	idxGen IndexGenerator // Keeps track of free slots
}

// Add adds T to the slice, returning its index. Will use a free slot if
// available, otherwise will expand the slice.
func (s *PoolSlice[T]) Add(value T) int {
	index := s.idxGen.Alloc()
	if index == len(s.Items) {
		s.Items = append(s.Items, value)
	} else {
		s.Items[index] = value
	}
	return index
}

// Remove removes the item at the index, replacing it with the zero value and
// marking its slot for re-use.
func (s *PoolSlice[T]) Remove(index int) {
	var zero T
	s.Items[index] = zero
	s.idxGen.Release(index)
}

// An IndexGenerator allocates slice indexes, allowing re-use. This lets you insert
// remove items within a slice dynamically, re-using empty slots when available.
type IndexGenerator struct {
	free []int // list of free slots
	next int   // next index to use when no slots are free
}

// Alloc allocates a free index. Will use an empty slot if available.
func (g *IndexGenerator) Alloc() int {
	if len(g.free) == 0 {
		ret := g.next
		g.next++
		return ret
	}
	ret := g.free[len(g.free)-1]
	g.free = g.free[:len(g.free)-1]
	return ret
}

// Release marks an index as free.
func (g *IndexGenerator) Release(index int) {
	g.free = append(g.free, index)
}
