package capnp

type CapTable struct {
	cs []Client
}

func (ct *CapTable) Reset(cs ...Client) {
	for _, c := range ct.cs {
		c.Release()
	}

	ct.cs = append(ct.cs[:0], cs...)
}

func (ct CapTable) Len() int {
	return len(ct.cs)
}

func (ct CapTable) At(i int) Client {
	return ct.cs[i]
}

func (ct CapTable) Contains(ifc Interface) bool {
	return ifc.IsValid() && ifc.Capability() < CapabilityID(ct.Len())
}

func (ct CapTable) Get(ifc Interface) (c Client) {
	if ct.Contains(ifc) {
		c = ct.cs[ifc.Capability()]
	}

	return
}

func (ct CapTable) Set(id CapabilityID, c Client) {
	ct.cs[id] = c
}

// Add appends a capability to the message's capability table and
// returns its ID.  It "steals" c's reference: the Message will release
// the client when calling Reset.
func (ct *CapTable) Add(c Client) CapabilityID {
	ct.cs = append(ct.cs, c)
	return CapabilityID(ct.Len() - 1)
}
