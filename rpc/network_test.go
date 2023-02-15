package rpc_test

import (
	"context"
	"net"
	"sync"

	"capnproto.org/go/capnp/v3/exp/spsc"
	"capnproto.org/go/capnp/v3/rpc"
)

type inMemoryPeerId uint64

type inMemoryEdge struct {
	To, From inMemoryPeerId
}

type inMemoryNetworkRef struct {
	myId    inMemoryPeerId
	network *inMemoryNetwork
}

type inMemoryNetwork struct {
	mu          sync.Mutex
	nextId      inMemoryPeerId
	connections map[inMemoryEdge]*rpc.Conn
	incoming    map[inMemoryPeerId]spsc.Queue[inMemoryIncomingConn]
}

type inMemoryIncomingConn struct {
	Conn net.Conn
	Id   inMemoryPeerId
}

func newInMemoryNetwork() *inMemoryNetwork {
	return &inMemoryNetwork{
		connections: make(map[inMemoryEdge]*rpc.Conn),
		incoming:    make(map[inMemoryPeerId]spsc.Queue[inMemoryIncomingConn]),
	}
}

func (n *inMemoryNetwork) Join() rpc.Network {
	n.mu.Lock()
	defer n.mu.Unlock()
	ret := inMemoryNetworkRef{
		myId:    n.nextId,
		network: n,
	}
	n.nextId++
	return ret
}

func (n *inMemoryNetwork) getAcceptQueue(id inMemoryPeerId) spsc.Queue[inMemoryIncomingConn] {
	q, ok := n.incoming[id]
	if !ok {
		q = spsc.New[inMemoryIncomingConn]()
		n.incoming[id] = q
	}
	return q
}

func (n inMemoryNetworkRef) MyId() rpc.PeerId {
	return rpc.PeerId{n.myId}
}

func (n inMemoryNetworkRef) Dial(dst rpc.PeerId, opts *rpc.Options) (*rpc.Conn, error) {
	if opts == nil {
		opts = &rpc.Options{}
	}
	opts.Network = n
	opts.PeerId = dst
	dstId := dst.Value.(inMemoryPeerId)
	edge := inMemoryEdge{
		From: n.myId,
		To:   dstId,
	}

	n.network.mu.Lock()
	conn, ok := n.network.connections[edge]
	defer n.network.mu.Unlock()
	if ok {
		return conn, nil
	}
	q := n.network.getAcceptQueue(dstId)
	c1, c2 := net.Pipe()
	q.Send(inMemoryIncomingConn{
		Conn: c1,
		Id:   n.myId,
	})
	conn = rpc.NewConn(rpc.NewStreamTransport(c2), opts)
	n.network.connections[edge] = conn
	return conn, nil
}

func (n inMemoryNetworkRef) Accept(ctx context.Context) (*rpc.Conn, error) {
	n.network.mu.Lock()
	q := n.network.getAcceptQueue(n.myId)
	n.network.mu.Unlock()

	incoming, err := q.Recv(ctx)
	if err != nil {
		return nil, err
	}

	n.network.mu.Lock()
	defer n.network.mu.Unlock()

	conn := rpc.NewConn(rpc.NewStreamTransport(incoming.Conn), &rpc.Options{
		Network: n,
		PeerId:  rpc.PeerId{incoming.Id},
	})
	n.network.connections[inMemoryEdge{
		From: n.myId,
		To:   incoming.Id,
	}] = conn
	return conn, nil
}

func (n inMemoryNetworkRef) Introduce(provider, recipient *rpc.Conn) (rpc.IntroductionInfo, error) {
	panic("TODO")
}
func (n inMemoryNetworkRef) DialIntroduced(capId rpc.ThirdPartyCapId) (*rpc.Conn, rpc.ProvisionId, error) {
	panic("TODO")
}
func (n inMemoryNetworkRef) AcceptIntroduced(recipientId rpc.RecipientId) (*rpc.Conn, error) {
	panic("TODO")
}
