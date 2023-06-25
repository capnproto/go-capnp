// Package testnetwork provides an in-memory implementation of rpc.Network for testing purposes.
package testnetwork

import (
	"context"
	"net"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/exp/spsc"
	"capnproto.org/go/capnp/v3/rpc"
	"zenhack.net/go/util"
	"zenhack.net/go/util/sync/mutex"
)

// PeerID is the implementation of peer ids used by a test network
type PeerID uint64

type edge struct {
	To, From PeerID
}

func (e edge) Flip() edge {
	return edge{
		To:   e.From,
		From: e.To,
	}
}

type TestNetwork struct {
	myID   PeerID
	global *Joiner
}

// A Joiner is a global view of a test network, which can be joined by a
// peer to acquire a TestNetwork.
type Joiner struct {
	state mutex.Mutex[joinerState]
}

type joinerState struct {
	nextID      PeerID
	nextNonce   uint64
	connections map[edge]*connectionEntry
	incoming    map[PeerID]spsc.Queue[PeerID]
}

type connectionEntry struct {
	Transport rpc.Transport
	Conn      *rpc.Conn // Might be nil, if we haven't initialized this yet.
}

func NewJoiner() *Joiner {
	return &Joiner{
		state: mutex.New(joinerState{
			connections: make(map[edge]*connectionEntry),
		}),
	}
}

func (j *Joiner) Join() TestNetwork {
	return mutex.With1(&j.state, func(js *joinerState) TestNetwork {
		ret := TestNetwork{
			myID:   js.nextID,
			global: j,
		}
		js.nextID++
		return ret
	})
}

func (j *joinerState) getAcceptQueue(id PeerID) spsc.Queue[PeerID] {
	q, ok := j.incoming[id]
	if !ok {
		q = spsc.New[PeerID]()
		j.incoming[id] = q
	}
	return q
}

func (n TestNetwork) LocalID() rpc.PeerID {
	return rpc.PeerID{n.myID}
}

func (n TestNetwork) Dial(dst rpc.PeerID, opts *rpc.Options) (*rpc.Conn, error) {
	conn, _, err := n.dial(dst, true, opts)
	return conn, err
}

// DialTransport is like Dial, except that a Conn is not created, and the raw Transport is
// returned instead.
func (n TestNetwork) DialTransport(dst rpc.PeerID) (rpc.Transport, error) {
	_, trans, err := n.dial(dst, false, nil)
	return trans, err
}

// Helper for Dial and DialTransport; setupConn indicates whether to create the Conn
// (if false it will be nil).
func (n TestNetwork) dial(dst rpc.PeerID, setupConn bool, opts *rpc.Options) (*rpc.Conn, rpc.Transport, error) {
	if opts == nil {
		opts = &rpc.Options{}
	}
	opts.Network = n
	opts.RemotePeerID = dst
	dstID := dst.Value.(PeerID)
	toEdge := edge{
		From: n.myID,
		To:   dstID,
	}
	fromEdge := toEdge.Flip()

	return mutex.With3(&n.global.state, func(state *joinerState) (*rpc.Conn, rpc.Transport, error) {
		ent, ok := state.connections[toEdge]
		if !ok {
			c1, c2 := net.Pipe()
			t1 := rpc.NewStreamTransport(c1)
			t2 := rpc.NewStreamTransport(c2)
			ent = &connectionEntry{Transport: t1}
			state.connections[toEdge] = ent
			state.connections[fromEdge] = &connectionEntry{Transport: t2}

		}
		if setupConn && ent.Conn == nil {
			ent.Conn = rpc.NewConn(ent.Transport, opts)
		} else {
			// There's already a connection, so we're not going to use this, but
			// we own it. So drop it:
			opts.BootstrapClient.Release()
		}
		return ent.Conn, ent.Transport, nil
	})
}

func (n TestNetwork) Accept(ctx context.Context, opts *rpc.Options) (*rpc.Conn, error) {
	q := mutex.With1(&n.global.state, func(js *joinerState) spsc.Queue[PeerID] {
		return js.getAcceptQueue(n.myID)
	})

	incoming, err := q.Recv(ctx)
	if err != nil {
		return nil, err
	}
	opts.Network = n
	opts.RemotePeerID = rpc.PeerID{incoming}
	return mutex.With2(&n.global.state, func(js *joinerState) (*rpc.Conn, error) {
		edge := edge{
			From: n.myID,
			To:   incoming,
		}
		ent := js.connections[edge]
		if ent.Conn == nil {
			ent.Conn = rpc.NewConn(ent.Transport, opts)
		} else {
			opts.BootstrapClient.Release()
		}
		return ent.Conn, nil
	})
}

func makePeerAndNonce(peerID, nonce uint64) PeerAndNonce {
	_, seg := capnp.NewSingleSegmentMessage(nil)
	ret, err := NewPeerAndNonce(seg)
	util.Chkfatal(err)
	ret.SetPeerId(peerID)
	ret.SetNonce(nonce)
	return ret
}

func (n TestNetwork) Introduce(provider, recipient *rpc.Conn) (rpc.IntroductionInfo, error) {
	providerPeerID := uint64(provider.RemotePeerID().Value.(PeerID))
	recipientPeerID := uint64(recipient.RemotePeerID().Value.(PeerID))
	return mutex.With2(&n.global.state, func(js *joinerState) (rpc.IntroductionInfo, error) {
		nonce := js.nextNonce
		js.nextNonce++
		return rpc.IntroductionInfo{
			SendToRecipient: rpc.ThirdPartyCapID(makePeerAndNonce(providerPeerID, nonce).ToPtr()),
			SendToProvider:  rpc.RecipientID(makePeerAndNonce(recipientPeerID, nonce).ToPtr()),
		}, nil
	})
}
func (n TestNetwork) DialIntroduced(capID rpc.ThirdPartyCapID, introducedBy *rpc.Conn) (*rpc.Conn, rpc.ProvisionID, error) {
	cid := PeerAndNonce(capnp.Ptr(capID).Struct())
	pid := makePeerAndNonce(
		uint64(introducedBy.RemotePeerID().Value.(PeerID)),
		cid.Nonce(),
	)
	conn, err := n.Dial(rpc.PeerID{PeerID(cid.PeerId())}, nil)
	return conn, rpc.ProvisionID(pid.ToPtr()), err
}
func (n TestNetwork) AcceptIntroduced(recipientID rpc.RecipientID, introducedBy *rpc.Conn) (*rpc.Conn, error) {
	panic("TODO")
}
