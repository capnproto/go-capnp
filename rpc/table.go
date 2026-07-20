package rpc

// localTable stores entries whose IDs are allocated by this vat.  Its zero
// value is ready for use.  Removed IDs are reused in increasing order.
//
// localTable does not provide synchronization.  RPC callers protect it with
// Conn.lk.
type localTable[ID ~uint32, T any] struct {
	entries []T
	next    uint32
	free    uintSet
	used    uintSet
}

func (t *localTable[ID, T]) Add(value T) ID {
	var id uint32
	if first, ok := t.free.min(); ok {
		id = uint32(first)
		t.free.remove(first)
	} else {
		id = t.next
		if id == ^uint32(0) {
			// All locally allocated IDs are under application control, but a
			// peer can retain exports long enough to exhaust the space.
			panic("overflow ID")
		}
		t.next++
	}
	if int(id) == len(t.entries) {
		t.entries = append(t.entries, value)
	} else {
		t.entries[id] = value
	}
	t.used.add(uint(id))
	return ID(id)
}

func (t *localTable[ID, T]) Find(id ID) (T, bool) {
	if !t.used.has(uint(id)) {
		var zero T
		return zero, false
	}
	return t.entries[uint32(id)], true
}

// Remove removes id and returns its entry.  An ID becomes reusable only when
// Remove succeeds.
func (t *localTable[ID, T]) Remove(id ID) (T, bool) {
	value, ok := t.take(id)
	if ok {
		t.release(id)
	}
	return value, ok
}

// take removes an entry while reserving its ID.  It is used for questions
// whose ID cannot be reused until a Finish message is successfully queued.
func (t *localTable[ID, T]) take(id ID) (T, bool) {
	value, ok := t.Find(id)
	if !ok {
		return value, false
	}
	var zero T
	t.entries[uint32(id)] = zero
	t.used.remove(uint(id))
	return value, true
}

func (t *localTable[ID, T]) release(id ID) {
	t.free.add(uint(id))
}

func (t *localTable[ID, T]) Range(f func(ID, T) bool) {
	for i := uint32(0); i < uint32(len(t.entries)); i++ {
		if t.used.has(uint(i)) && !f(ID(i), t.entries[i]) {
			return
		}
	}
}

// Clear removes every entry without invoking any entry-owned cleanup.
func (t *localTable[ID, T]) Clear() (entries []T) {
	t.Range(func(_ ID, value T) bool {
		entries = append(entries, value)
		return true
	})
	*t = localTable[ID, T]{}
	return entries
}

// remoteTable stores entries whose IDs are selected by the peer.  Its zero
// value is ready for use and IDs are therefore allowed to be sparse.
//
// remoteTable does not provide synchronization.  RPC callers protect it with
// Conn.lk.
type remoteTable[ID ~uint32, T any] struct {
	entries map[ID]T
}

// Create inserts value unless id already exists.
func (t *remoteTable[ID, T]) Create(id ID, value T) bool {
	if _, ok := t.Find(id); ok {
		return false
	}
	if t.entries == nil {
		t.entries = make(map[ID]T)
	}
	t.entries[id] = value
	return true
}

func (t *remoteTable[ID, T]) FindOrCreate(id ID, create func() T) (T, bool) {
	if value, ok := t.Find(id); ok {
		return value, false
	}
	value := create()
	t.Create(id, value)
	return value, true
}

func (t *remoteTable[ID, T]) Find(id ID) (T, bool) {
	value, ok := t.entries[id]
	return value, ok
}

func (t *remoteTable[ID, T]) Remove(id ID) (T, bool) {
	value, ok := t.Find(id)
	if ok {
		delete(t.entries, id)
	}
	return value, ok
}

func (t *remoteTable[ID, T]) Range(f func(ID, T) bool) {
	for id, value := range t.entries {
		if !f(id, value) {
			return
		}
	}
}

// Clear removes every entry without invoking any entry-owned cleanup.
func (t *remoteTable[ID, T]) Clear() (entries []T) {
	t.Range(func(_ ID, value T) bool {
		entries = append(entries, value)
		return true
	})
	t.entries = nil
	return entries
}

// uintSet is a set of unsigned integers represented by a bit set.  It assumes
// that integers are packed closely to zero.
type uintSet []uint64

func (s uintSet) has(i uint) bool {
	j := i / 64
	mask := uint64(1) << (i % 64)
	return j < uint(len(s)) && s[j]&mask != 0
}

func (s *uintSet) add(i uint) {
	j := i / 64
	mask := uint64(1) << (i % 64)
	if j >= uint(len(*s)) {
		s2 := make(uintSet, j+1)
		copy(s2, *s)
		*s = s2
	}
	(*s)[j] |= mask
}

func (s uintSet) remove(i uint) {
	j := i / 64
	mask := uint64(1) << (i % 64)
	if j < uint(len(s)) {
		s[j] &^= mask
	}
}

func (s uintSet) min() (_ uint, ok bool) {
	for i, x := range s {
		if x == 0 {
			continue
		}
		for j := uint(0); j < 64; j++ {
			if x&(1<<j) != 0 {
				return uint(i)*64 + j, true
			}
		}
	}
	return 0, false
}
