package rpc

import (
	"context"

	capnp "capnproto.org/go/capnp/v3"
)

type PeerId struct {
	Value any
}

type ThirdPartyCapId capnp.Ptr
type RecipientId capnp.Ptr
type ProvisionId capnp.Ptr

type IntroductionInfo struct {
	SendToRecipient ThirdPartyCapId
	SendToProvider  RecipientId
}

type Network interface {
	MyId() PeerId
	Dial(PeerId, *Options) (*Conn, error)
	Accept(context.Context) (*Conn, error)
	Introduce(provider, recipient *Conn) (IntroductionInfo, error)
	DialIntroduced(capId ThirdPartyCapId) (*Conn, ProvisionId, error)
	AcceptIntroduced(recipientId RecipientId) (*Conn, error)
}
