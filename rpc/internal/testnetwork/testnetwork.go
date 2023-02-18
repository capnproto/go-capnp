// Package testnetwork provides an in-memory implementation of rpc.Network for testing purposes.
package testnetwork

import (
	"context"
	"net"
	"sync"

	"capnproto.org/go/capnp/v3/exp/spsc"
	"capnproto.org/go/capnp/v3/rpc"
)

// PeerID is the implementation of peer ids used by a test network
type PeerID uint64

type edge struct {
	To, From PeerID
}

type network struct {
	myID   PeerID
	global *Joiner
}

// A Joiner is a global view of a test network, which can be joined by a
// peer to acquire a Network.
type Joiner struct {
	mu          sync.Mutex
	nextID      PeerID
	connections map[edge]*rpc.Conn
	incoming    map[PeerID]spsc.Queue[incomingConn]
}

type incomingConn struct {
	Conn net.Conn
	ID   PeerID
}

func NewJoiner() *Joiner {
	return &Joiner{
		connections: make(map[edge]*rpc.Conn),
		incoming:    make(map[PeerID]spsc.Queue[incomingConn]),
	}
}

func (j *Joiner) Join() rpc.Network {
	j.mu.Lock()
	defer j.mu.Unlock()
	ret := network{
		myID:   j.nextID,
		global: j,
	}
	j.nextID++
	return ret
}

func (j *Joiner) getAcceptQueue(id PeerID) spsc.Queue[incomingConn] {
	q, ok := j.incoming[id]
	if !ok {
		q = spsc.New[incomingConn]()
		j.incoming[id] = q
	}
	return q
}

func (n network) MyID() rpc.PeerID {
	return rpc.PeerID{n.myID}
}

func (n network) Dial(dst rpc.PeerID, opts *rpc.Options) (*rpc.Conn, error) {
	if opts == nil {
		opts = &rpc.Options{}
	}
	opts.Network = n
	opts.RemotePeerID = dst
	dstID := dst.Value.(PeerID)
	edge := edge{
		From: n.myID,
		To:   dstID,
	}

	n.global.mu.Lock()
	defer n.global.mu.Unlock()
	conn, ok := n.global.connections[edge]
	if ok {
		return conn, nil
	}
	q := n.global.getAcceptQueue(dstID)
	c1, c2 := net.Pipe()
	q.Send(incomingConn{
		Conn: c1,
		ID:   n.myID,
	})
	conn = rpc.NewConn(rpc.NewStreamTransport(c2), opts)
	n.global.connections[edge] = conn
	return conn, nil
}

func (n network) Accept(ctx context.Context, opts *rpc.Options) (*rpc.Conn, error) {
	n.global.mu.Lock()
	q := n.global.getAcceptQueue(n.myID)
	n.global.mu.Unlock()

	incoming, err := q.Recv(ctx)
	if err != nil {
		return nil, err
	}
	opts.Network = n
	opts.RemotePeerID = rpc.PeerID{incoming.ID}
	n.global.mu.Lock()
	defer n.global.mu.Unlock()
	conn := rpc.NewConn(rpc.NewStreamTransport(incoming.Conn), opts)
	n.global.connections[edge{
		From: n.myID,
		To:   incoming.ID,
	}] = conn
	return conn, nil
}

func (n network) Introduce(provider, recipient *rpc.Conn) (rpc.IntroductionInfo, error) {
	panic("TODO")
}
func (n network) DialIntroduced(capID rpc.ThirdPartyCapID) (*rpc.Conn, rpc.ProvisionID, error) {
	panic("TODO")
}
func (n network) AcceptIntroduced(recipientID rpc.RecipientID) (*rpc.Conn, error) {
	panic("TODO")
}
