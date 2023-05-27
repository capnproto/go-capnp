package capnp

// CapTable is the indexed list of the clients referenced in the
// message. Capability pointers inside the message will use this
// table to map pointers to Clients.   The table is populated by
// the RPC system.
//
// https://capnproto.org/encoding.html#capabilities-interfaces
type CapTable struct {
	// We maintain two parallel structurs of clients and corresponding
	// snapshots. We need to store both, so that Get*() can hand out
	// borrowed references in both cases.
	clients   []Client
	snapshots []ClientSnapshot
}

// Reset the cap table, releasing all capabilities and setting
// the length to zero.
func (ct *CapTable) Reset() {
	for _, c := range ct.clients {
		c.Release()
	}
	for _, s := range ct.snapshots {
		s.Release()
	}

	ct.clients = ct.clients[:0]
	ct.snapshots = ct.snapshots[:0]
}

// Len returns the number of capabilities in the table.
func (ct CapTable) Len() int {
	return len(ct.clients)
}

// ClientAt returns the client at the given index of the table.
func (ct CapTable) ClientAt(i int) Client {
	return ct.clients[i]
}

// SnapshotAt is like ClientAt, except that it returns a snapshot.
func (ct CapTable) SnapshotAt(i int) ClientSnapshot {
	return ct.snapshots[i]
}

// Contains returns true if the supplied interface corresponds
// to a client already present in the table.
func (ct CapTable) Contains(ifc Interface) bool {
	return ifc.IsValid() && ifc.Capability() < CapabilityID(ct.Len())
}

// GetClient gets the client corresponding to the supplied interface.
// It returns a null client if the interface's CapabilityID isn't
// in the table.
func (ct CapTable) GetClient(ifc Interface) (c Client) {
	if ct.Contains(ifc) {
		c = ct.clients[ifc.Capability()]
	}
	return
}

// GetSnapshot is like GetClient, except that it returns a snapshot
// instead of a Client.
func (ct CapTable) GetSnapshot(ifc Interface) (s ClientSnapshot) {
	if ct.Contains(ifc) {
		s = ct.snapshots[ifc.Capability()]
	}
	return
}

// SetClient sets the client for the supplied capability ID.  If a client
// for the given ID already exists, it will be replaced without
// releasing.
func (ct CapTable) SetClient(id CapabilityID, c Client) {
	ct.snapshots[id] = c.Snapshot()
	ct.clients[id] = c.Steal()
}

// SetSnapshot is like SetClient, but takes a snapshot.
func (ct CapTable) SetSnapshot(id CapabilityID, s ClientSnapshot) {
	ct.clients[id] = s.Client()
	ct.snapshots[id] = s.Steal()
}

// AddClient appends a capability to the message's capability table and
// returns its ID.  It "steals" c's reference: the Message will release
// the client when calling Reset.
func (ct *CapTable) AddClient(c Client) CapabilityID {
	ct.snapshots = append(ct.snapshots, c.Snapshot())
	ct.clients = append(ct.clients, c.Steal())
	return CapabilityID(ct.Len() - 1)
}

// AddSnapshot is like AddClient, except that it takes a snapshot rather
// than a Client.
func (ct *CapTable) AddSnapshot(s ClientSnapshot) CapabilityID {
	defer s.Release()
	return ct.AddClient(s.Client())
}
