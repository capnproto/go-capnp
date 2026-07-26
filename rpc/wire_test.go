package rpc_test

import (
	"context"
	"fmt"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/pogs"
	"capnproto.org/go/capnp/v3/rpc"
	rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"
)

type rpcMessage struct {
	Which         rpccp.Message_Which
	Unimplemented *rpcMessage
	Abort         *rpcException
	Bootstrap     *rpcBootstrap
	Call          *rpcCall
	Return        *rpcReturn
	Finish        *rpcFinish
	Resolve       *rpcResolve
	Release       *rpcRelease
	Disembargo    *rpcDisembargo
}

func sendMessage(ctx context.Context, t rpc.Transport, msg *rpcMessage) error {
	outMsg, err := t.NewMessage()
	if err != nil {
		return fmt.Errorf("send message: %v", err)
	}
	defer outMsg.Release()
	if err := pogs.Insert(rpccp.Message_TypeID, capnp.Struct(outMsg.Message()), msg); err != nil {
		return fmt.Errorf("send message: %v", err)
	}
	if err := outMsg.Send(); err != nil {
		return fmt.Errorf("send message: %v", err)
	}
	return nil
}

func recvMessage(ctx context.Context, t rpc.Transport) (*rpcMessage, capnp.ReleaseFunc, error) {
	inMsg, err := t.RecvMessage()
	if err != nil {
		return nil, nil, err
	}
	r := new(rpcMessage)
	if err := pogs.Extract(r, rpccp.Message_TypeID, capnp.Struct(inMsg.Message())); err != nil {
		inMsg.Release()
		return nil, nil, fmt.Errorf("extract RPC message: %v", err)
	}
	if r.Which == rpccp.Message_Which_abort ||
		r.Which == rpccp.Message_Which_bootstrap ||
		r.Which == rpccp.Message_Which_finish ||
		r.Which == rpccp.Message_Which_resolve ||
		r.Which == rpccp.Message_Which_release ||
		r.Which == rpccp.Message_Which_disembargo {
		// These messages are guaranteed to not contain pointers back to
		// the original message, so we can release them early.
		inMsg.Release()
		return r, func() {}, nil
	}
	return r, inMsg.Release, nil
}

type rpcException struct {
	Reason string
	Type   rpccp.Exception_Type
}

type rpcBootstrap struct {
	QuestionID uint32 `capnp:"questionId"`
}

type rpcCall struct {
	QuestionID              uint32 `capnp:"questionId"`
	Target                  rpcMessageTarget
	InterfaceID             uint64 `capnp:"interfaceId"`
	MethodID                uint16 `capnp:"methodId"`
	AllowThirdPartyTailCall bool
	Params                  rpcPayload
	SendResultsTo           rpcCallSendResultsTo
}

type rpcCallSendResultsTo struct {
	Which rpccp.Call_sendResultsTo_Which
}

type rpcReturn struct {
	AnswerID         uint32 `capnp:"answerId"`
	ReleaseParamCaps bool
	NoFinishNeeded   bool

	Which                 rpccp.Return_Which
	Results               *rpcPayload
	Exception             *rpcException
	TakeFromOtherQuestion uint32
}

type rpcFinish struct {
	QuestionID        uint32 `capnp:"questionId"`
	ReleaseResultCaps bool
}

type rpcMessageTarget struct {
	Which          rpccp.MessageTarget_Which
	ImportedCap    uint32
	PromisedAnswer *rpcPromisedAnswer
}

type rpcPayload struct {
	Content  capnp.Ptr
	CapTable []rpcCapDescriptor
}

type rpcCapDescriptor struct {
	Which          rpccp.CapDescriptor_Which
	SenderHosted   uint32
	SenderPromise  uint32
	ReceiverHosted uint32
	ReceiverAnswer *rpcPromisedAnswer
}

type rpcPromisedAnswer struct {
	QuestionID uint32 `capnp:"questionId"`
	Transform  []rpcPromisedAnswerOp
}

func (pa *rpcPromisedAnswer) transformEquals(path ...uint16) bool {
	for _, op := range pa.Transform {
		switch op.Which {
		case rpccp.PromisedAnswer_Op_Which_noop:
			// Skip.
		case rpccp.PromisedAnswer_Op_Which_getPointerField:
			if len(path) == 0 || path[0] != op.GetPointerField {
				return false
			}
			path = path[1:]
		default:
			return false
		}
	}
	return len(path) == 0
}

type rpcPromisedAnswerOp struct {
	Which           rpccp.PromisedAnswer_Op_Which
	GetPointerField uint16
}

type rpcResolve struct {
	PromiseID uint32 `capnp:"promiseId"`
	Which     rpccp.Resolve_Which
	Cap       *rpcCapDescriptor
	Exception *rpcException
}

type rpcRelease struct {
	ID             uint32 `capnp:"id"`
	ReferenceCount uint32
}

type rpcDisembargo struct {
	Target  rpcMessageTarget
	Context rpcDisembargoContext
}

type rpcDisembargoContext struct {
	Which            rpccp.Disembargo_context_Which
	SenderLoopback   uint32
	ReceiverLoopback uint32
	Accept           capnp.Ptr
}
