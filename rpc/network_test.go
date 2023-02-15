package rpc_test

import (
	"context"
	"net"
	"sync"

	"capnproto.org/go/capnp/v3/exp/spsc"
	"capnproto.org/go/capnp/v3/rpc"
)

type inMemoryPeerID uint64

type inMemoryEdge struct {
	To, From inMemoryPeerID
}

type inMemoryNetworkRef struct {
	myID    inMemoryPeerID
	network *inMemoryNetwork
}

type inMemoryNetwork struct {
	mu          sync.Mutex
	nextID      inMemoryPeerID
	connections map[inMemoryEdge]*rpc.Conn
	incoming    map[inMemoryPeerID]spsc.Queue[inMemoryIncomingConn]
}

type inMemoryIncomingConn struct {
	Conn net.Conn
	ID   inMemoryPeerID
}

func newInMemoryNetwork() *inMemoryNetwork {
	return &inMemoryNetwork{
		connections: make(map[inMemoryEdge]*rpc.Conn),
		incoming:    make(map[inMemoryPeerID]spsc.Queue[inMemoryIncomingConn]),
	}
}

func (n *inMemoryNetwork) Join() rpc.Network {
	n.mu.Lock()
	defer n.mu.Unlock()
	ret := inMemoryNetworkRef{
		myID:    n.nextID,
		network: n,
	}
	n.nextID++
	return ret
}

func (n *inMemoryNetwork) getAcceptQueue(id inMemoryPeerID) spsc.Queue[inMemoryIncomingConn] {
	q, ok := n.incoming[id]
	if !ok {
		q = spsc.New[inMemoryIncomingConn]()
		n.incoming[id] = q
	}
	return q
}

func (n inMemoryNetworkRef) MyID() rpc.PeerID {
	return rpc.PeerID{n.myID}
}

func (n inMemoryNetworkRef) Dial(dst rpc.PeerID, opts *rpc.Options) (*rpc.Conn, error) {
	if opts == nil {
		opts = &rpc.Options{}
	}
	opts.Network = n
	opts.RemotePeerID = dst
	dstID := dst.Value.(inMemoryPeerID)
	edge := inMemoryEdge{
		From: n.myID,
		To:   dstID,
	}

	n.network.mu.Lock()
	conn, ok := n.network.connections[edge]
	defer n.network.mu.Unlock()
	if ok {
		return conn, nil
	}
	q := n.network.getAcceptQueue(dstID)
	c1, c2 := net.Pipe()
	q.Send(inMemoryIncomingConn{
		Conn: c1,
		ID:   n.myID,
	})
	conn = rpc.NewConn(rpc.NewStreamTransport(c2), opts)
	n.network.connections[edge] = conn
	return conn, nil
}

func (n inMemoryNetworkRef) Accept(ctx context.Context, opts *rpc.Options) (*rpc.Conn, error) {
	n.network.mu.Lock()
	q := n.network.getAcceptQueue(n.myID)
	n.network.mu.Unlock()

	incoming, err := q.Recv(ctx)
	if err != nil {
		return nil, err
	}
	opts.Network = n
	opts.RemotePeerID = rpc.PeerID{incoming.ID}
	n.network.mu.Lock()
	defer n.network.mu.Unlock()
	conn := rpc.NewConn(rpc.NewStreamTransport(incoming.Conn), opts)
	n.network.connections[inMemoryEdge{
		From: n.myID,
		To:   incoming.ID,
	}] = conn
	return conn, nil
}

func (n inMemoryNetworkRef) Introduce(provider, recipient *rpc.Conn) (rpc.IntroductionInfo, error) {
	panic("TODO")
}
func (n inMemoryNetworkRef) DialIntroduced(capID rpc.ThirdPartyCapID) (*rpc.Conn, rpc.ProvisionID, error) {
	panic("TODO")
}
func (n inMemoryNetworkRef) AcceptIntroduced(recipientID rpc.RecipientID) (*rpc.Conn, error) {
	panic("TODO")
}
