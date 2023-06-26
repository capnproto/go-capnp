package rpc_test

import (
	"context"
	"testing"

	"capnproto.org/go/capnp/v3"
	"capnproto.org/go/capnp/v3/pogs"
	"capnproto.org/go/capnp/v3/rpc"
	"capnproto.org/go/capnp/v3/rpc/internal/testcapnp"
	"capnproto.org/go/capnp/v3/rpc/internal/testnetwork"
	rpccp "capnproto.org/go/capnp/v3/std/capnp/rpc"
	"github.com/stretchr/testify/require"
)

type rpcProvide struct {
	QuestionID uint32 `capnp:"questionId"`
	Target     rpcMessageTarget
	Recipient  capnp.Ptr
}

func TestSendProvide(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	j := testnetwork.NewJoiner()

	pp := testcapnp.PingPong_ServerToClient(&pingPonger{})
	defer pp.Release()

	cfgOpts := func(opts *rpc.Options) {
		opts.ErrorReporter = testErrorReporter{tb: t}
	}

	introducer := j.Join(cfgOpts)
	recipient := j.Join(cfgOpts)
	provider := j.Join(cfgOpts)

	go introducer.Serve(ctx)

	rConn, err := introducer.Dial(recipient.LocalID())
	require.NoError(t, err)

	pConn, err := introducer.Dial(provider.LocalID())
	require.NoError(t, err)

	rBs := rConn.Bootstrap(ctx)
	defer rBs.Release()
	pBs := pConn.Bootstrap(ctx)
	defer pBs.Release()

	rTrans, err := recipient.DialTransport(introducer.LocalID())
	require.NoError(t, err)

	pTrans, err := provider.DialTransport(introducer.LocalID())
	require.NoError(t, err)

	bootstrapExportID := uint32(10)
	doBootstrap := func(trans rpc.Transport) {
		// Receive bootstrap
		rmsg, release, err := recvMessage(ctx, trans)
		require.NoError(t, err)
		defer release()
		require.Equal(t, rpccp.Message_Which_bootstrap, rmsg.Which)
		qid := rmsg.Bootstrap.QuestionID

		// Write back return
		outMsg, err := trans.NewMessage()
		require.NoError(t, err, "trans.NewMessage()")
		iptr := capnp.NewInterface(outMsg.Message().Segment(), 0)
		require.NoError(t, pogs.Insert(rpccp.Message_TypeID, capnp.Struct(outMsg.Message()), &rpcMessage{
			Which: rpccp.Message_Which_return,
			Return: &rpcReturn{
				AnswerID: qid,
				Which:    rpccp.Return_Which_results,
				Results: &rpcPayload{
					Content: iptr.ToPtr(),
					CapTable: []rpcCapDescriptor{
						{
							Which:        rpccp.CapDescriptor_Which_senderHosted,
							SenderHosted: bootstrapExportID,
						},
					},
				},
			},
		}))
		require.NoError(t, outMsg.Send())

		// Receive finish
		rmsg, release, err = recvMessage(ctx, trans)
		require.NoError(t, err)
		defer release()
		require.Equal(t, rpccp.Message_Which_finish, rmsg.Which)
		require.Equal(t, qid, rmsg.Finish.QuestionID)
	}
	doBootstrap(rTrans)
	require.NoError(t, rBs.Resolve(ctx))
	doBootstrap(pTrans)
	require.NoError(t, pBs.Resolve(ctx))

	futEmpty, rel := testcapnp.EmptyProvider(pBs).GetEmpty(ctx, nil)
	defer rel()

	emptyExportID := uint32(30)
	{
		// Receive call
		rmsg, release, err := recvMessage(ctx, pTrans)
		require.NoError(t, err)
		defer release()
		require.Equal(t, rpccp.Message_Which_call, rmsg.Which)
		qid := rmsg.Call.QuestionID
		require.Equal(t, uint64(testcapnp.EmptyProvider_TypeID), rmsg.Call.InterfaceID)
		require.Equal(t, uint16(0), rmsg.Call.MethodID)

		// Send return
		outMsg, err := pTrans.NewMessage()
		require.NoError(t, err)
		seg := outMsg.Message().Segment()
		results, err := capnp.NewStruct(seg, capnp.ObjectSize{
			PointerCount: 1,
		})
		require.NoError(t, err)
		iptr := capnp.NewInterface(seg, 0)
		results.SetPtr(0, iptr.ToPtr())
		require.NoError(t, sendMessage(ctx, pTrans, &rpcMessage{
			Which: rpccp.Message_Which_return,
			Return: &rpcReturn{
				Which: rpccp.Return_Which_results,
				Results: &rpcPayload{
					Content: results.ToPtr(),
					CapTable: []rpcCapDescriptor{
						{
							Which:        rpccp.CapDescriptor_Which_senderHosted,
							SenderHosted: emptyExportID,
						},
					},
				},
			},
		}))

		// Receive finish
		rmsg, release, err = recvMessage(ctx, pTrans)
		require.NoError(t, err)
		defer release()
		require.Equal(t, rpccp.Message_Which_finish, rmsg.Which)
		require.Equal(t, qid, rmsg.Finish.QuestionID)
	}

	resEmpty, err := futEmpty.Struct()
	require.NoError(t, err)
	empty := resEmpty.Empty()

	_, rel = testcapnp.CapArgsTest(rBs).Call(ctx, func(p testcapnp.CapArgsTest_call_Params) error {
		return p.SetCap(capnp.Client(empty))
	})
	defer rel()

	//var provideQid uint32
	{
		// Provider should receive a provide message
		rmsg, release, err := recvMessage(ctx, pTrans)
		require.NoError(t, err)
		defer release()
		require.Equal(t, rpccp.Message_Which_provide, rmsg.Which)
		//provideQid = rmsg.Provide.QuestionID
		require.Equal(t, rpccp.MessageTarget_Which_importedCap, rmsg.Provide.Target.Which)
		require.Equal(t, emptyExportID, rmsg.Provide.Target.ImportedCap)
	}

	panic("TODO: check for messages on rTrans (call, resolve)")
}
