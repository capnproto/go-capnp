package rpc

import (
	"context"

	capnp "capnproto.org/go/capnp/v3"
)

type PeerID struct {
	Value any
}

type ThirdPartyCapID capnp.Ptr
type RecipientID capnp.Ptr
type ProvisionID capnp.Ptr

type IntroductionInfo struct {
	SendToRecipient ThirdPartyCapID
	SendToProvider  RecipientID
}

type Network interface {
	MyID() PeerID
	Dial(PeerID, *Options) (*Conn, error)
	Accept(context.Context) (*Conn, error)
	Introduce(provider, recipient *Conn) (IntroductionInfo, error)
	DialIntroduced(capID ThirdPartyCapID) (*Conn, ProvisionID, error)
	AcceptIntroduced(recipientID RecipientID) (*Conn, error)
}
